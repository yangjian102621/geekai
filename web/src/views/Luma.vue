<template>
  <div class="page-luma">
    <div class="prompt-box">
      <div class="images">
        <template v-for="(img, index) in images" :key="img">
          <div class="item">
            <el-image :src="img" fit="cover"/>
            <el-icon @click="remove(img)"><CircleCloseFilled /></el-icon>
          </div>
          <div class="btn-swap" v-if="images.length == 2 && index == 0">
            <svg class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" width="32" height="32" @click="switchReverse">
              <path d="M96 416h832c0.32 0 0.576-0.192 0.896-0.192a30.656 30.656 0 0 0 30.976-30.976c-0.064-0.256 0.128-0.512 0.128-0.832a31.424 31.424 0 0 0-14.912-26.368l-188.48-188.48a30.72 30.72 0 1 0-43.456 43.456L852.544 352H96a32 32 0 0 0 0 64z m832 192H96c-0.32 0-0.576 0.192-0.896 0.192a30.528 30.528 0 0 0-30.976 30.976c0.064 0.256-0.128 0.512-0.128 0.832 0 11.264 6.144 20.672 14.912 26.368l188.48 188.48a30.72 30.72 0 1 0 43.456-43.456L171.456 672H928a32 32 0 0 0 0-64z"></path>
            </svg>
          </div>
        </template>
      </div>
      <div class="prompt-container">
        <div class="input-container">
          <div class="upload-icon" v-if="images.length < 2">
            <el-upload
                class="avatar-uploader"
                :auto-upload="true"
                :show-file-list="false"
                :http-request="upload"
                accept=".jpg,.png,.jpeg"
            >
              <i class="iconfont icon-image"></i>
            </el-upload>
          </div>
          <textarea
              class="prompt-input"
              :rows="row"
              v-model="formData.prompt"
              placeholder="请输入提示词或者上传图片"
              autofocus>
                      </textarea>
          <div class="send-icon" @click="create">
            <i class="iconfont icon-send"></i>
          </div>
        </div>

        <div class="params">
          <div class="item-group">
            <span class="label">循环参考图</span>
            <el-switch  v-model="formData.loop" size="small" style="--el-switch-on-color:#BF78BF;" />
          </div>
          <div class="item-group">
            <span class="label">提示词优化</span>
            <el-switch  v-model="formData.expand_prompt" size="small" style="--el-switch-on-color:#BF78BF;" />
          </div>
        </div>
      </div>
    </div>



    <el-container class="video-container" v-loading="loading" element-loading-background="rgba(100,100,100,0.3)">
      <h2 class="h-title">你的作品</h2>

      <!-- <el-row :gutter="20" class="videos">
        <el-col :span="8" class="item" :key="item.id" v-for="item in videos">
          <div class="video-box" @mouseover="item.playing = true" @mouseout="item.playing = false">
            <img :src="item.cover"  :alt="item.name" v-show="!item.playing"/>
            <video :src="item.url"  preload="auto" :autoplay="true" loop="loop" muted="muted" v-show="item.playing">
              您的浏览器不支持视频播放
            </video>
          </div>
          <div class="video-name">{{item.name}}</div>
          <div class="opts">
            <button class="btn" @click="download(item)" :disabled="item.downloading">
              <i class="iconfont icon-download" v-if="!item.downloading"></i>
              <el-image src="/images/loading.gif" fit="cover" v-else />
              <span>下载</span>
            </button>
          </div>
        </el-col>
      </el-row> -->

      <div class="list-box" v-if="!noData">
        <div v-for="item in list" :key="item.id">
          <div class="item" v-if="item.progress === 100">
            <div class="left">
              <div class="container">
                <el-image :src="item.cover_url" fit="cover" />
                <!-- <div class="duration">{{formatTime(item.duration)}}</div> -->
                <button class="play" @click="play(item)">
                  <img src="/images/play.svg" alt=""/>
                </button>
              </div>
            </div>
            <div class="center">
              <div class="title">
                <a>{{item.prompt}}</a>
                <!-- <span class="model" v-if="item.major_model_version">{{item.major_model_version}}</span>
                <span class="model" v-if="item.type === 4">用户上传</span>
                <span class="model" v-if="item.type === 3">
                  <i class="iconfont icon-mp3"></i>
                  完整歌曲
                </span>
                <span class="model" v-if="item.ref_song">
                    <i class="iconfont icon-link"></i>
                    {{item.ref_song.title}}
                  </span> -->
              </div>
              <div class="tags" v-if="item.prompt_ext">{{item.prompt_ext}}</div>
            </div>
            <div class="right">
              <div class="tools">
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
            <div style="width: 60px; flex-shrink: 0; display: flex; align-items: center;" v-if="item.params.start_img_url">
              <el-image :src="item.params.start_img_url" fit="cover" />
            </div>            
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
              <generating v-else>正在生成视频</generating>
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
    </el-container>
  </div>
</template>

<script setup>
import {onMounted, reactive, ref} from "vue";
import {CircleCloseFilled} from "@element-plus/icons-vue";
import {httpDownload, httpPost, httpGet} from "@/utils/http";
import {checkSession} from "@/store/cache";
import {showMessageError} from "@/utils/dialog";
import {ElMessage, ElMessageBox} from "element-plus";
import Clipboard from "clipboard";
import BlackSwitch from "@/components/ui/BlackSwitch.vue";
import Generating from "@/components/ui/Generating.vue";


const row = ref(1)
const images = ref([])

const formData = reactive({
  prompt: '',
  expand_prompt: true,
  loop: true,
  first_frame_img: '',
  end_frame_img: ''
})

const videos = ref([
  {
    id: 1,
    name: 'a dancing girl',
    url: 'https://storage.cdn-luma.com/dream_machine/d133794f-3124-4059-a9f2-e5fed79f0d5b/watermarked_video01944f69966f14d33b6c4486a8cfb8dde.mp4',
    cover: 'https://storage.cdn-luma.com/dream_machine/d133794f-3124-4059-a9f2-e5fed79f0d5b/video_0_thumb.jpg',
    playing: false
  },
  {
    id: 1,
    name: 'a dancing girl a dancing girl a dancing girl a dancing girl a dancing girl',
    url: 'https://storage.cdn-luma.com/dream_machine/92efa55a-f381-4161-a999-54f8fe460fca/watermarked_video0e5aad607a0644c66850d1d77022db847.mp4',
    cover: 'https://storage.cdn-luma.com/dream_machine/92efa55a-f381-4161-a999-54f8fe460fca/video_1_thumb.jpg',
    playing: false
  },
  {
    id: 1,
    name: 'a dancing girl',
    url: 'https://storage.cdn-luma.com/dream_machine/d133794f-3124-4059-a9f2-e5fed79f0d5b/watermarked_video01944f69966f14d33b6c4486a8cfb8dde.mp4',
    cover: 'https://storage.cdn-luma.com/dream_machine/d133794f-3124-4059-a9f2-e5fed79f0d5b/video_0_thumb.jpg',
    playing: false
  },
  {
    id: 1,
    name: 'a dancing girl',
    url: 'https://storage.cdn-luma.com/dream_machine/92efa55a-f381-4161-a999-54f8fe460fca/watermarked_video0e5aad607a0644c66850d1d77022db847.mp4',
    cover: 'https://storage.cdn-luma.com/dream_machine/92efa55a-f381-4161-a999-54f8fe460fca/video_1_thumb.jpg',
    playing: false
  },
  {
    id: 1,
    name: 'a dancing girl',
    url: 'https://storage.cdn-luma.com/dream_machine/d133794f-3124-4059-a9f2-e5fed79f0d5b/watermarked_video01944f69966f14d33b6c4486a8cfb8dde.mp4',
    cover: 'https://storage.cdn-luma.com/dream_machine/d133794f-3124-4059-a9f2-e5fed79f0d5b/video_0_thumb.jpg',
    playing: false
  },
])

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

  const _socket = new WebSocket(host + `/api/video/client?user_id=${userId.value}`);
  _socket.addEventListener('open', () => {
    socket.value = _socket;
  });

  _socket.addEventListener('message', event => {
    if (event.data instanceof Blob) {
      const reader = new FileReader();
      reader.readAsText(event.data, "UTF-8")
      reader.onload = () => {
        const message = String(reader.result)
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

onMounted(()=>{
  checkSession().then(user => {
    userId.value = user.id
    connect()
  })
  fetchData(1)
})

const download = (item) => {
  const downloadURL = `${process.env.VUE_APP_API_HOST}/api/download?url=${item.url}`
  // parse filename
  const urlObj = new URL(item.url);
  const fileName = urlObj.pathname.split('/').pop();
  item.downloading = true
  httpDownload(downloadURL).then(response  => {
    const blob = new Blob([response.data]);
    const link = document.createElement('a');
    link.href = URL.createObjectURL(blob);
    link.download = fileName;
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
    URL.revokeObjectURL(link.href);
    item.downloading = false
  }).catch(() => {
    showMessageError("下载失败")
    item.downloading = false
  })
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
    httpGet("/api/video/remove", {id: item.id}).then(() => {
      ElMessage.success("任务删除成功")
      fetchData()
    }).catch(e => {
      ElMessage.error("任务删除失败：" + e.message)
    })
  }).catch(() => {
  })
}

const publishJob = (item) => {
  httpGet("/api/video/publish", {id: item.id, publish:item.publish}).then(() => {
    ElMessage.success("操作成功")
  }).catch(e => {
    ElMessage.error("操作失败：" + e.message)
  })
}

const getShareURL = (item) => {
  return `${location.protocol}//${location.host}/song/${item.id}`
}

const upload = (file) => {
  const formData = new FormData();
  formData.append('file', file.file, file.name);
  // 执行上传操作
  httpPost('/api/upload', formData).then((res) => {
    images.value.push(res.data.url)
    ElMessage.success({message: "上传成功", duration: 500})
  }).catch((e) => {
    ElMessage.error('图片上传失败:' + e.message)
  })
};

const remove = (img) => {
  images.value = images.value.filter(item => item !== img)
}

const switchReverse = () => {
  images.value = images.value.reverse()
}
const loading = ref(false)
const list = ref([])
const noData = ref(true)
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)
const fetchData = (_page) => {
  if (_page) {
    page.value = _page
  }
  httpGet("/api/video/list",{page:page.value, page_size:pageSize.value, type: 'luma'}).then(res => {
    total.value = res.data.total
    const items = []
    for (let v of res.data.items) {
      if(v.prompt == '风雨交加的夜晚'){
        v.progress = 100
        v.video_url = 'https://storage.cdn-luma.com/dream_machine/92efa55a-f381-4161-a999-54f8fe460fca/watermarked_video0e5aad607a0644c66850d1d77022db847.mp4'
        v.cover_url = 'https://storage.cdn-luma.com/dream_machine/92efa55a-f381-4161-a999-54f8fe460fca/video_1_thumb.jpg'
      }

      if (v.progress === 100) {
        //v.major_model_version = v['raw_data']['major_model_version']
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

// 创建视频
const create = () => {

  const len =  images.value.length;
  if(len){
    formData.first_frame_img = images.value[0]
    if(len == 2){
      formData.end_frame_img = images.value[1]
    }
  }

  httpPost("/api/video/luma/create", formData).then(() => {
    fetchData(1)
    showMessageOK("创建任务成功")
  }).catch(e => {
    showMessageError("创建任务失败："+e.message)
  })
}


</script>

<style lang="stylus" scoped>
@import "@/assets/css/luma.styl"
</style>
