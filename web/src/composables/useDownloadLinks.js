// * 下载链接共享状态：多组件（TopNav / Footer / HomePage / InvitePage）共用同一份
// * /api/update/latest 数据，模块级 ref 保证只请求一次。

import { computed, ref } from 'vue'
import { FALLBACK_LINKS } from '@/services/site'

export const downloadLinks = ref(FALLBACK_LINKS)
let loaded = false

async function loadLinks() {
  if (loaded) return
  loaded = true
  try {
    const res = await fetch('/api/update/latest', { cache: 'no-cache' })
    if (!res.ok) return
    const data = await res.json()
    // 接口成功但缺 key → 视为暂未提供（null），不伪造链接。
    downloadLinks.value = {
      version: data.version ?? '',
      windows: data.url || data.platforms?.windows || null,
      macArm64: data.platforms?.['mac-arm64'] || null,
      macX64: data.platforms?.['mac-x64'] || null,
    }
  } catch (e) {
    // 网络失败时保留 FALLBACK 与重试入口。
  }
}

/* 网络失败后手动重试：重新拉取并更新共享状态。 */
export async function refreshDownloadLinks() {
  loaded = false
  await loadLinks()
}

export function useDownloadLinks() {
  loadLinks()
  return downloadLinks
}

/* 决定默认下载哪个包：Mac 给 mac 安装包，其他平台给 Windows 包。 */
export function currentOs() {
  const ua = typeof navigator === 'undefined' ? '' : navigator.userAgent
  if (/Macintosh|Mac OS X/i.test(ua)) return 'mac'
  if (/Windows/i.test(ua)) return 'windows'
  return 'other'
}

/* 通过 userAgentData.architecture 识别 mac 架构；拿不到就返回 unknown，绝不默认 Intel 或 ARM。 */
export function currentMacArch() {
  if (typeof navigator === 'undefined') return 'unknown'
  const arch = navigator.userAgentData?.architecture
  if (arch === 'arm') return 'arm64'
  if (arch === 'x86') return 'x64'
  return 'unknown'
}

export function downloadOptions(links) {
  const os = currentOs()
  const macArch = currentMacArch()
  // 未知 mac 架构时不显示“当前设备”。
  const macBadgeKnown = os === 'mac' && macArch !== 'unknown'
  return [
    {
      key: 'windows',
      labelKey: 'download.win',
      noteKey: 'download.win_note',
      href: links.windows,
      badgeKey: os === 'windows' ? 'download.current' : '',
      available: !!links.windows,
    },
    {
      key: 'mac-arm64',
      labelKey: 'download.mac_arm64',
      noteKey: 'download.mac_arm64_note',
      href: links.macArm64,
      badgeKey: macBadgeKnown && macArch === 'arm64' ? 'download.current' : '',
      available: !!links.macArm64,
    },
    {
      key: 'mac-x64',
      labelKey: 'download.mac_intel',
      noteKey: 'download.mac_intel_note',
      href: links.macX64,
      badgeKey: macBadgeKnown && macArch === 'x64' ? 'download.current' : '',
      available: !!links.macX64,
    },
  ]
}

export function recommendedHref(links) {
  const os = currentOs()
  if (os === 'mac') {
    const arch = currentMacArch()
    if (arch === 'unknown') return null // 未知架构：不默认，让用户选
    const target = arch === 'arm64' ? links.macArm64 : links.macX64
    return target ?? null // 架构已知但包暂未提供 → null，不伪造链接
  }
  return links.windows || null
}

/* 兼容旧调用点（顶栏、页脚、定价页）：当前设备推荐的安装包地址。 */
export function useLatestSetup() {
  const links = useDownloadLinks()
  return computed(() => recommendedHref(links.value))
}
