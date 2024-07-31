<template>
  <div class="page-suno" :style="{ height: winHeight + 'px' }">
    <div class="left-bar">
      <div class="bar-top">
        <el-tooltip effect="light" content="定义模式" placement="top">
          <black-switch  v-model:value="custom" size="large"  />
        </el-tooltip>
        <black-select v-model:value="data.model" :options="models" placeholder="请选择模型" style="width: 100px" />
      </div>

      <div class="params">
        <div class="pure-music">
          <span class="switch"><black-switch  v-model:value="data.instrumental" size="default"  /></span>
          <span class="text">纯音乐</span>
        </div>
        <div v-if="custom">
          <div class="item-group" v-if="!data.instrumental">
            <div class="label">
              <span class="text">歌词</span>
              <el-popover placement="right"
                          :width="200"
                          trigger="hover" content="自己写歌词或寻求 AI 的帮助。使用两节歌词（8 行）可获得最佳效果。">
                <template #reference>
                  <el-icon>
                    <InfoFilled/>
                  </el-icon>
                </template>
              </el-popover>
            </div>
            <div class="item"
                 v-loading="generating"
                 element-loading-text="正在生成歌词..."
                 element-loading-background="rgba(122, 122, 122, 0.8)">
              <black-input v-model:value="data.lyrics" type="textarea" :rows="10" placeholder="请在这里输入你自己写的歌词..."/>
              <button class="btn btn-lyric" @click="createLyric">生成歌词</button>
            </div>
          </div>

          <div class="item-group">
            <div class="label">
              <span class="text">音乐风格</span>
              <el-popover placement="right"
                          :width="200"
                          trigger="hover" content="描述您想要的音乐风格（例如“原声流行音乐”）。Sunos 模特无法识别艺术家的名字，但能够理解音乐流派和氛围。">
                <template #reference>
                  <el-icon>
                    <InfoFilled/>
                  </el-icon>
                </template>
              </el-popover>
            </div>
            <div class="item">
              <black-input v-model:value="data.tags" type="textarea" :maxlength="120" :rows="3" placeholder="请输入音乐风格，多个风格之间用英文逗号隔开..."/>
            </div>

            <div class="tag-select">
              <div class="inner">
                <span
                    class="tag"
                    @click="selectTag(tag)"
                    v-for="tag in tags"
                    :key="tag.value">{{ tag.label }}</span>
              </div>
              </div>
          </div>

          <div class="item-group">
            <div class="label">
              <span class="text">歌曲名称</span>
              <el-popover placement="right"
                          :width="200"
                          trigger="hover" content="给你的歌曲起一个标题，以便于分享、发现和组织。">
                <template #reference>
                  <el-icon>
                    <InfoFilled/>
                  </el-icon>
                </template>
              </el-popover>
            </div>
            <div class="item">
              <black-input v-model:value="data.title" type="textarea" :rows="1" placeholder="请输入歌曲名称..."/>
            </div>
          </div>
        </div>

        <div v-else>
          <div class="label">
            <span class="text">歌曲描述</span>
            <el-popover placement="right"
                        :width="200"
                        trigger="hover" content="描述您想要的音乐风格和主题（例如：关于假期的流行音乐）。请使用流派和氛围，而不是特定的艺术家和歌曲风格，AI无法识别。">
              <template #reference>
                <el-icon>
                  <InfoFilled/>
                </el-icon>
              </template>
            </el-popover>
          </div>
          <div class="item">
            <black-input v-model:value="data.prompt" type="textarea" :rows="10" placeholder="例如：一首关于鸟人的摇滚歌曲..."/>
          </div>
        </div>

        <div class="ref-song" v-if="refSong">
          <div class="label">
            <span class="text">续写</span>
            <el-popover placement="right"
                        :width="200"
                        trigger="hover" content="输入额外的歌词，根据您之前的歌词来扩展歌曲。">
              <template #reference>
                <el-icon>
                  <InfoFilled/>
                </el-icon>
              </template>
            </el-popover>
          </div>

          <div class="item">
            <div class="song">
              <el-image :src="refSong.cover_url" fit="cover" />
              <span class="title">{{refSong.title}}</span>
              <el-button type="info" @click="removeRefSong" size="small" :icon="Delete" circle />
            </div>
            <div class="extend-secs">
              从  <input v-model="refSong.extend_secs" type="text"/> 秒开始续写
            </div>
          </div>
        </div>

        <div class="item">
          <button class="create-btn" @click="create">
            <img src="/images/create-new.svg" alt=""/>
            <span>{{btnText}}</span>
          </button>
        </div>
      </div>
    </div>
    <div class="right-box" v-loading="loading" element-loading-background="rgba(100,100,100,0.3)">
      <div class="list-box" v-if="!noData">
        <div v-for="item in list">
          <div class="item" v-if="item.progress === 100">
            <div class="left">
              <div class="container">
                <el-image :src="item.cover_url" fit="cover" />
                <div class="duration">{{formatTime(item.duration)}}</div>
                <button class="play" @click="play(item)">
                  <img src="/images/play.svg" alt=""/>
                </button>
              </div>
            </div>
            <div class="center">
              <div class="title">
                <a :href="'/song/'+item.song_id"  target="_blank">{{item.title}}</a>
                <span class="model">{{item.major_model_version}}</span>
                <span class="model" v-if="item.ref_song">
                    <i class="iconfont icon-link"></i>
                    {{item.ref_song.title}}
                  </span>
              </div>
              <div class="tags">{{item.tags}}</div>
            </div>
            <div class="right">
              <div class="tools">
                <el-tooltip effect="light" content="以当前歌曲为素材继续创作" placement="top">
                  <button class="btn" @click="extend(item)">续写</button>
                </el-tooltip>

                <button class="btn btn-publish">
                  <span class="text">发布</span>
                  <black-switch v-model:value="item.publish" @change="publishJob(item)" size="small" />
                </button>

                <el-tooltip effect="light" content="下载歌曲" placement="top">
                  <a :href="item.audio_url" :download="item.title+'.mp3'" target="_blank">
                    <button class="btn btn-icon">
                      <i class="iconfont icon-download"></i>
                    </button>
                  </a>
                </el-tooltip>

                <el-tooltip effect="light" content="复制歌曲链接" placement="top">
                  <button class="btn btn-icon copy-link" :data-clipboard-text="getShareURL(item)" >
                    <i class="iconfont icon-share1"></i>
                  </button>
                </el-tooltip>

                <el-tooltip effect="light" content="编辑" placement="top">
                  <button class="btn btn-icon" @click="update(item)">
                    <i class="iconfont icon-edit"></i>
                  </button>
                </el-tooltip>

                <el-tooltip effect="light" content="删除" placement="top">
                  <button class="btn btn-icon" @click="removeJob(item)">
                    <i class="iconfont icon-remove"></i>
                  </button>
                </el-tooltip>
              </div>
            </div>
          </div>
          <div class="task" v-else>
            <div class="left">
              <div class="title">
                <span v-if="item.title">{{item.title}}</span>
                <span v-else>{{item.prompt}}</span>
              </div>
            </div>
            <div class="center">
              <div class="failed" v-if="item.progress === 101">
                {{item.err_msg}}
              </div>
              <generating v-else />
            </div>
            <div class="right">
              <el-button type="info" @click="removeJob(item)" circle>
                <i class="iconfont icon-remove"></i>
              </el-button>
            </div>
          </div>
        </div>
      </div>
      <el-empty :image-size="100" description="没有任何作品，赶紧去创作吧！" v-else/>

      <div class="pagination">
        <el-pagination v-if="total > pageSize" background
                       style="--el-pagination-button-bg-color:#414141;
                       --el-pagination-button-color:#d1d1d1;
                       --el-disabled-bg-color:#414141;
                       --el-color-primary:#666666;
                       --el-pagination-hover-color:#e1e1e1"
                       layout="total,prev, pager, next"
                       :hide-on-single-page="true"
                       v-model:current-page="page"
                       v-model:page-size="pageSize"
                       @current-change="fetchData(page)"
                       :total="total"/>

      </div>

      <div class="music-player" v-if="showPlayer">
        <music-player :songs="playList" ref="playerRef" :show-close="true" @close="showPlayer = false" />
      </div>
    </div>

    <black-dialog v-model:show="showDialog" title="修改歌曲" @cancal="showDialog = false" @confirm="updateSong" :width="500">
      <form class="form">
        <div class="form-item">
          <div class="label">歌曲名称</div>
          <input class="input" v-model="editData.title" type="text" />
        </div>

        <div class="form-item">
          <div class="label">封面图片</div>
          <el-upload
              class="avatar-uploader"
              :auto-upload="true"
              :show-file-list="false"
              :http-request="uploadCover"
              accept=".png,.jpg,.jpeg,.bmp"
          >
            <el-avatar :src="editData.cover" shape="square" :size="100"/>
          </el-upload>
        </div>
      </form>
    </black-dialog>
  </div>
</template>

<script setup>
import {nextTick, onMounted, onUnmounted, ref, watch} from "vue"
import {Delete, InfoFilled} from "@element-plus/icons-vue";
import BlackSelect from "@/components/ui/BlackSelect.vue";
import BlackSwitch from "@/components/ui/BlackSwitch.vue";
import BlackInput from "@/components/ui/BlackInput.vue";
import MusicPlayer from "@/components/MusicPlayer.vue";
import {compact} from "lodash";
import {httpGet, httpPost} from "@/utils/http";
import {showMessageError, showMessageOK} from "@/utils/dialog";
import Generating from "@/components/ui/Generating.vue";
import {checkSession} from "@/action/session";
import {ElMessage, ElMessageBox} from "element-plus";
import {formatTime} from "@/utils/libs";
import Clipboard from "clipboard";
import BlackDialog from "@/components/ui/BlackDialog.vue";
import Compressor from "compressorjs";

const winHeight = ref(window.innerHeight - 50)
const custom = ref(false)
const models = ref([
  {label: "v3.0", value: "chirp-v3-0"},
  {label: "v3.5", value:"chirp-v3-5"}
])
const tags = ref([
  {label: "女声", value: "female vocals"},
  {label: "男声", value: "male vocals"},
  {label: "流行", value: "pop"},
  {label: "摇滚", value: "rock"},
  {label: "硬摇滚", value: "hard rock"},
  {label: "电音", value: "electronic"},
  {label: "金属", value: "metal"},
  {label: "重金属", value: "heavy metal"},
  {label: "节拍", value: "beat"},
  {label: "弱拍", value: "upbeat"},
  {label: "合成器", value: "synth"},
  {label: "吉他", value: "guitar"},
  {label: "钢琴", value: "piano"},
  {label: "小提琴", value: "violin"},
  {label: "贝斯", value: "bass"},
  {label: "嘻哈", value: "hip hop"},
])
const data = ref({
  model: "chirp-v3-0",
  tags: "",
  lyrics: "",
  prompt: "",
  title: "",
  instrumental: false,
  ref_task_id: "",
  extend_secs: 0,
  ref_song_id: "",
})
const loading = ref(true)
const noData = ref(false)
const playList = ref([])
const playerRef = ref(null)
const showPlayer = ref(false)
const list = ref([])
const btnText = ref("开始创作")
const refSong = ref(null)
const showDialog = ref(false)
const editData = ref({title:"",cover:"",id:0})

const socket = ref(null)
const userId = ref(0)
const connect = () => {
  let host = process.env.VUE_APP_WS_HOST
  if (host === '') {
    if (location.protocol === 'https:') {
      host = 'wss://' + location.host;
    } else {
      host = 'ws://' + location.host;
    }
  }

  const _socket = new WebSocket(host + `/api/suno/client?user_id=${userId.value}`);
  _socket.addEventListener('open', () => {
    socket.value = _socket;
  });

  _socket.addEventListener('message', event => {
    if (event.data instanceof Blob) {
      const reader = new FileReader();
      reader.readAsText(event.data, "UTF-8")
      reader.onload = () => {
        const message = String(reader.result)
        console.log(message)
        if (message === "FINISH" || message === "FAIL") {
          fetchData()
        }
      }
    }
  });

  _socket.addEventListener('close', () => {
    if (socket.value !== null) {
      connect()
    }
  });
}

const clipboard = ref(null)
onMounted(() => {
  clipboard.value = new Clipboard('.copy-link');
  clipboard.value.on('success', () => {
    ElMessage.success("复制歌曲链接成功！");
  })

  clipboard.value.on('error', () => {
    ElMessage.error('复制失败！');
  })

  checkSession().then(user => {
    userId.value = user.id
    connect()
  })
  fetchData(1)
})

onUnmounted(() => {
  clipboard.value.destroy()
})

const page = ref(1)
const pageSize = ref(10)
const total = ref(0)
const fetchData = (_page) => {
  if (_page) {
    page.value = _page
  }
  httpGet("/api/suno/list",{page:page.value, page_size:pageSize.value}).then(res => {
    total.value = res.data.total
    const items = []
    for (let v of res.data.items) {
      if (v.progress === 100) {
        v.major_model_version = v['raw_data']['major_model_version']
      }
      items.push(v)
    }
    loading.value = false
    list.value = items
    noData.value = list.value.length === 0
  }).catch(e => {
    loading.value = false
    noData.value = true
    showMessageError("获取作品列表失败："+e.message)
  })
}

// 创建新的歌曲
const create = () => {
  data.value.type = custom.value ? 2 : 1
  data.value.ref_task_id = refSong.value ? refSong.value.task_id : ""
  data.value.ref_song_id = refSong.value ? refSong.value.song_id : ""
  data.value.extend_secs = refSong.value ? refSong.value.extend_secs : 0
  if (custom.value) {
    if (data.value.lyrics === "") {
      return showMessageError("请输入歌词")
    }
    if (data.value.title === "") {
      return showMessageError("请输入歌曲标题")
    }
  } else {
    if (data.value.prompt === "") {
      return showMessageError("请输入歌曲描述")
    }
  }
  if (refSong.value && data.value.extend_secs > refSong.value.duration) {
    return showMessageError("续写开始时间不能超过原歌曲长度")
  }

  httpPost("/api/suno/create", data.value).then(() => {
    fetchData(1)
    showMessageOK("创建任务成功")
  }).catch(e => {
    showMessageError("创建任务失败："+e.message)
  })
}

// 续写歌曲
const extend = (item) => {
  refSong.value = item
  refSong.value.extend_secs = item.duration
  data.value.title = item.title
  custom.value = true
  btnText.value = "续写歌曲"
}

// 更细歌曲
const update = (item) => {
  showDialog.value = true
  editData.value.title = item.title
  editData.value.cover = item.cover_url
  editData.value.id = item.id
}

const updateSong = ()  => {
  if (editData.value.title === "" || editData.value.cover === "") {
    return showMessageError("歌曲标题和封面不能为空")
  }
  httpPost("/api/suno/update", editData.value).then(() => {
    showMessageOK("更新歌曲成功")
    showDialog.value = false
    fetchData()
  }).catch(e => {
    showMessageError("更新歌曲失败："+e.message)
  })
}

watch(() => custom.value, (newValue) => {
  if (!newValue) {
    removeRefSong()
  }
})

const removeRefSong = () => {
  refSong.value = null
  btnText.value = "开始创作"
}

const play = (item) => {
  playList.value = [item]
  showPlayer.value = true
  nextTick(()=> playerRef.value.play())
}

const selectTag = (tag) => {
  if (data.value.tags.length + tag.value.length >= 119) {
    return
  }
  data.value.tags = compact([...data.value.tags.split(","), tag.value]).join(",")
}

const removeJob = (item) => {
  ElMessageBox.confirm(
      '此操作将会删除任务相关文件，继续操作码?',
      '删除提示',
      {
        confirmButtonText: '确认',
        cancelButtonText: '取消',
        type: 'warning',
      }
  ).then(() => {
    httpGet("/api/suno/remove", {id: item.id}).then(() => {
      ElMessage.success("任务删除成功")
      fetchData()
    }).catch(e => {
      ElMessage.error("任务删除失败：" + e.message)
    })
  }).catch(() => {
  })
}

const publishJob = (item) => {
  httpGet("/api/suno/publish", {id: item.id, publish:item.publish}).then(() => {
    ElMessage.success("操作成功")
  }).catch(e => {
    ElMessage.error("操作失败：" + e.message)
  })
}

const getShareURL = (item) => {
  return `${location.protocol}//${location.host}/song/${item.id}`
}

const uploadCover = (file) => {
  // 压缩图片并上传
  new Compressor(file.file, {
    quality: 0.6,
    success(result) {
      const formData = new FormData();
      formData.append('file', result, result.name);
      // 执行上传操作
      httpPost('/api/upload', formData).then((res) => {
        editData.value.cover = res.data.url
        ElMessage.success({message: "上传成功", duration: 500})
      }).catch((e) => {
        ElMessage.error('图片上传失败:' + e.message)
      })
    },
    error(err) {
      console.log(err.message);
    },
  });
}

const generating = ref(false)
const createLyric = () => {
  if (data.value.lyrics === "") {
    return showMessageError("请输入歌词描述")
  }
  generating.value = true
  httpPost("/api/suno/lyric", {prompt: data.value.lyrics}).then(res => {
    const lines = res.data.split('\n');
    data.value.title = lines.shift().replace(/\*/g,"")
    lines.shift()
    data.value.lyrics = lines.join('\n');
    generating.value = false
  }).catch(e => {
    showMessageError("歌词生成失败："+e.message)
    generating.value = false
  })
}

</script>

<style lang="stylus" scoped>
@import "@/assets/css/suno.styl"
</style>
