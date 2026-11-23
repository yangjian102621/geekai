#!/usr/bin/env bash
set -euo pipefail

parse_installer() {
  if [[ $1 =~ ^ThinkBuddySetup-dt([0-9]+(\.[0-9]+)+)\.exe$ ]]; then
    platform=windows version=${BASH_REMATCH[1]}
  elif [[ $1 =~ ^ThinkBuddyDesktop-mac-(arm64|x86_64)-dt([0-9]+(\.[0-9]+)+)\.dmg$ ]]; then
    platform=${BASH_REMATCH[1]} version=${BASH_REMATCH[2]}
  else
    echo "无法识别安装包平台或版本，请保留原始安装包名称: $1" >&2
    return 2
  fi
}

validate_installers() {
  local file name platform version count=0 seen=' '
  shopt -s nullglob dotglob
  for file in "$1"/*; do
    name=${file##*/}
    [[ $name == .publish-sha256 && ${2:-} == with-checksums ]] && continue
    [[ -f $file && -s $file && ! -L $file && $name =~ ^[a-zA-Z0-9][a-zA-Z0-9._+-]*\.(exe|dmg)$ ]] || {
      echo "只允许非空 exe/dmg 普通文件，latest.json 由脚本自动更新: $name" >&2; exit 2;
    }
    parse_installer "$name"
    [[ $seen != *" $platform "* ]] || { echo "同一平台只允许一个安装包: $platform" >&2; exit 2; }
    seen+="$platform "
    count=$((count + 1))
  done
  shopt -u dotglob
  ((count > 0)) || { echo '安装包目录为空' >&2; exit 2; }
}

prepare_manifest() {
  local sha256 file platform version published_at
  published_at=$(TZ=UTC-8 date '+%Y-%m-%dT%H:%M:%S%:z')
  jq -e 'if type == "object" then . else error("latest.json 必须是 JSON 对象") end' "$1" > "$2"
  while read -r sha256 file; do
    file=${file#\*}
    parse_installer "$file"
    jq --arg platform "$platform" --arg version "$version" --arg file "$file" \
      --arg sha256 "$sha256" --arg published_at "$published_at" \
      --argjson size "$(wc -c < "$file")" '
      . as $root
      | (if $platform == "windows" then [] else ["mac", $platform] end) as $path
      | getpath($path) as $current
      | (($current.url // $root.url)
         | if type == "string" and test("^https?://[^/]+/")
           then sub("[?#].*$"; "") | sub("[^/]*$"; $file)
           else error("latest.json 缺少有效的下载 URL") end) as $url
      | setpath($path; (($current // {}) + {
          version: $version, url: $url, size: $size,
          sha256: $sha256, published_at: $published_at
        }))
    ' "$2" > "$2.tmp"
    mv -- "$2.tmp" "$2"
  done < .publish-sha256
}

if [[ ${1:-} == --remote ]]; then
  shift
  [[ $# == 5 ]] || exit 2
  remote_dir=$1 namespace=$2 selector=$3 container=$4 update_dir=$5
  [[ $remote_dir =~ ^/tmp/thinkbuddy-release\.[a-zA-Z0-9]+$ ]] || exit 2
  [[ $update_dir == /* && $update_dir != / && $update_dir != */ ]] || exit 2
  command -v kubectl >/dev/null
  command -v flock >/dev/null
  command -v jq >/dev/null || { echo 'ctl1 需要 jq，请先执行: sudo apt-get update && sudo apt-get install -y jq' >&2; exit 2; }
  # Serialize copies from different team members on ctl1.
  exec 9> /tmp/thinkbuddy-installers.lock
  flock -n 9 || { echo '另一个安装包发布正在运行，请稍后重试' >&2; exit 1; }
  validate_installers "$remote_dir/payload" with-checksums
  cd "$remote_dir/payload"
  [[ -f .publish-sha256 && ! -L .publish-sha256 ]] || exit 2
  sha256sum -c .publish-sha256
  echo "安装包来源: $remote_dir"
  while read -r sha256 file; do
    file=${file#\*}
    printf '文件: %s\n"size": %d,\n"sha256": "%s"\n\n' "$file" "$(wc -c < "$file")" "$sha256"
  done < .publish-sha256

  pods=$(kubectl -n "$namespace" get pods -l "$selector" \
    --field-selector=status.phase=Running \
    -o 'go-template={{range .items}}{{if not .metadata.deletionTimestamp}}{{.metadata.name}}{{"\n"}}{{end}}{{end}}')
  [[ -n "$pods" ]] || { echo '没有找到运行中的目标 Pod' >&2; exit 1; }
  mkdir -p "$remote_dir/manifests"
  while IFS= read -r pod; do
    kubectl -n "$namespace" wait --for=condition=Ready "pod/$pod" --timeout=120s
    kubectl -n "$namespace" exec "$pod" -c "$container" -- \
      sh -c 'command -v tar >/dev/null && command -v sha256sum >/dev/null && test -d "$1" && test -w "$1"' sh "$update_dir"
    kubectl -n "$namespace" exec "$pod" -c "$container" -- \
      cat "$update_dir/latest.json" > "$remote_dir/manifests/$pod.previous.json"
    prepare_manifest "$remote_dir/manifests/$pod.previous.json" "$remote_dir/manifests/$pod.json"
  done <<< "$pods"
  while IFS= read -r pod; do
    echo "复制安装包到 $namespace/$pod:$update_dir"
    stage=$(kubectl -n "$namespace" exec "$pod" -c "$container" -- \
      mktemp -d "$update_dir/.release-XXXXXXXXXX")
    kubectl -n "$namespace" cp "$remote_dir/payload/." "$pod:$stage" -c "$container"
    kubectl -n "$namespace" cp "$remote_dir/manifests/$pod.json" "$pod:$stage/.latest.json" -c "$container"
    kubectl -n "$namespace" exec "$pod" -c "$container" -- sh -eu -c '
      cd "$1"
      sha256sum -c .publish-sha256
      cp -p -- "$2/latest.json" "$2/latest.json.bak"
      chmod --reference="$2/latest.json" .latest.json
      for file in *; do mv -f -- "$file" "$2/$file"; done
      mv -f -- .latest.json "$2/latest.json"
      rm .publish-sha256
      cd "$2"
      rmdir "$1"
    ' sh "$stage" "$update_dir"
    echo "已更新 $namespace/$pod:$update_dir/latest.json（备份: latest.json.bak）"
    cat "$remote_dir/manifests/$pod.json"
  done <<< "$pods"
  echo '安装包复制完成，latest.json 已自动更新。'
  exit 0
fi

[[ $# == 1 && -d $1 ]] || {
  echo "用法: DEPLOY_HOST=<ctl1> DEPLOY_USER=<用户> bash $0 <本地安装包目录>" >&2; exit 2;
}
source_dir=$1
: "${DEPLOY_HOST:?缺少 DEPLOY_HOST（ctl1 地址）}"
: "${DEPLOY_USER:?缺少 DEPLOY_USER}"
DEPLOY_PORT=${DEPLOY_PORT:-22}
[[ $DEPLOY_HOST =~ ^[a-zA-Z0-9][a-zA-Z0-9.-]*$ ]] || exit 2
[[ $DEPLOY_USER =~ ^[a-zA-Z_][a-zA-Z0-9_-]*$ ]] || exit 2
[[ $DEPLOY_PORT =~ ^[0-9]{1,5}$ ]] && ((10#$DEPLOY_PORT > 0 && 10#$DEPLOY_PORT <= 65535)) || exit 2
target="$DEPLOY_USER@$DEPLOY_HOST"
ssh_options=(-o BatchMode=no -o StrictHostKeyChecking=ask -o ConnectTimeout=20 -o ServerAliveInterval=30 -o ServerAliveCountMax=3)
if [[ -n ${DEPLOY_SSH_KEY:-} ]]; then
  [[ -s $DEPLOY_SSH_KEY ]] || { echo 'SSH 私钥文件为空/不存在' >&2; exit 2; }
  ssh_options+=(-o IdentitiesOnly=yes -i "$DEPLOY_SSH_KEY")
fi
if [[ -n ${DEPLOY_KNOWN_HOSTS:-} ]]; then
  [[ -s $DEPLOY_KNOWN_HOSTS ]] || { echo 'known_hosts 文件为空/不存在' >&2; exit 2; }
  ssh_options+=(-o "UserKnownHostsFile=$DEPLOY_KNOWN_HOSTS")
fi

validate_installers "$source_dir"
checksum_command=(sha256sum)
if ! command -v sha256sum >/dev/null; then
  command -v shasum >/dev/null || { echo '需要 sha256sum 或 shasum' >&2; exit 2; }
  checksum_command=(shasum -a 256)
fi
work_dir=$(mktemp -d)
remote_dir=
upload_complete=0
cleanup() {
  local result=$?
  if [[ -n $remote_dir ]]; then
    if ((result == 0)); then
      ssh "${ssh_options[@]}" -p "$DEPLOY_PORT" "$target" "rm -rf -- '$remote_dir'" || \
        echo "ctl1 临时目录未能清理: $remote_dir" >&2
    else
      echo "上传或发布失败，ctl1 暂存文件已保留: $remote_dir" >&2
      if ((upload_complete)); then
        echo "可重新 SSH 登录 $target，在 ctl1 执行以下命令重试，无需再次上传:" >&2
        printf '%s\n' "$remote_command" >&2
      fi
    fi
  fi
  rm -rf -- "$work_dir"
}
trap cleanup EXIT
mkdir "$work_dir/payload"
cp -- "$source_dir"/* "$work_dir/payload/"
(cd "$work_dir/payload" && "${checksum_command[@]}" -- * > .publish-sha256)
remote_dir=$(ssh "${ssh_options[@]}" -p "$DEPLOY_PORT" "$target" \
  'command -v jq >/dev/null || { echo "ctl1 需要 jq，请先执行: sudo apt-get update && sudo apt-get install -y jq" >&2; exit 2; }; mktemp -d /tmp/thinkbuddy-release.XXXXXXXXXX')
if [[ ! $remote_dir =~ ^/tmp/thinkbuddy-release\.[a-zA-Z0-9]+$ ]]; then
  remote_dir=
  echo 'ctl1 未返回有效的上传目录' >&2
  exit 1
fi
echo "ctl1 暂存目录: $remote_dir"
scp "${ssh_options[@]}" -P "$DEPLOY_PORT" -r \
  "$work_dir/payload" "${BASH_SOURCE[0]}" "$target:$remote_dir/"
printf -v remote_command 'bash %q --remote %q %q %q %q %q' \
  "$remote_dir/${BASH_SOURCE[0]##*/}" "$remote_dir" \
  "${K8S_NAMESPACE:-thinkbuddy}" "${K8S_SELECTOR:-app=website}" \
  "${K8S_CONTAINER:-thinkbuddy}" "${POD_UPDATE_DIR:-/app/update}"
upload_complete=1
echo "发布目标: namespace=${K8S_NAMESPACE:-thinkbuddy}, selector=${K8S_SELECTOR:-app=website}, container=${K8S_CONTAINER:-thinkbuddy}, path=${POD_UPDATE_DIR:-/app/update}"
ssh "${ssh_options[@]}" -p "$DEPLOY_PORT" "$target" "$remote_command"
