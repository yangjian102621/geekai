<template>
  <div class="page-luma">
    <div class="prompt-box">
      <div class="images">
        <template v-for="(img, index) in images" :key="img">
          <div class="item">
            <el-image :src="replaceImg(img)" fit="cover"/>
            <el-icon @click="remove(img)">
              <CircleCloseFilled/>
            </el-icon>
          </div>
          <div class="btn-swap" v-if="images.length === 2 && index === 0">
            <i class="iconfont icon-exchange" @click="switchReverse"></i>
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
            <el-button class="generate-btn" size="small" @click="generatePrompt" color="#5865f2"
                       :disabled="isGenerating">
              <i class="iconfont icon-chuangzuo" style="margin-right: 5px"></i>
              <span>生成AI视频提示词</span>
            </el-button>
          </div>
          <div class="item-group">
            <span class="label">循环参考图</span>
            <el-switch v-model="formData.loop" size="small" style="--el-switch-on-color:#BF78BF;"/>
          </div>
          <div class="item-group">
            <span class="label">提示词优化</span>
            <el-switch v-model="formData.expand_prompt" size="small" style="--el-switch-on-color:#BF78BF;"/>
          </div>
        </div>
      </div>
    </div>


    <el-container class="video-container" v-loading="loading" element-loading-background="rgba(100,100,100,0.3)">
      <h2 class="h-title">你的作品</h2>

      <div class="list-box" v-if="!noData">
        <div v-for="item in list" :key="item.id">
          <div class="item">
            <div class="left">
              <div class="container">
                <div v-if="item.progress === 100">
                  <video class="video" :src="replaceImg(item.video_url)" preload="auto" loop="loop" muted="muted">
                    您的浏览器不支持视频播放
                  </video>
                  <button class="play" @click="play(item)">
                    <img src="/images/play.svg" alt=""/>
                  </button>
                </div>
                <el-image :src="item.cover_url" fit="cover" v-else-if="item.progress > 100"/>
                <generating message="正在生成视频" v-else/>

              </div>
            </div>
            <div class="center">
              <div class="failed" v-if="item.progress === 101">
                任务执行失败：{{ item.err_msg }}，任务提示词：{{ item.prompt }}
              </div>
              <div class="prompt" v-else>{{ item.prompt }}</div>
            </div>
            <div class="right" v-if="item.progress === 100">
              <div class="tools">
                <button class="btn btn-publish">
                  <span class="text">发布</span>
                  <black-switch v-model:value="item.publish" @change="publishJob(item)" size="small"/>
                </button>

                <el-tooltip effect="light" content="下载视频" placement="top">
                  <button class="btn btn-icon" @click="download(item)" :disabled="item.downloading">
                    <i class="iconfont icon-download" v-if="!item.downloading"></i>
                    <el-image src="/images/loading.gif" class="downloading" fit="cover" v-else/>
                  </button>
                </el-tooltip>
                <el-tooltip effect="light" content="删除" placement="top">
                  <button class="btn btn-icon" @click="removeJob(item)">
                    <i class="iconfont icon-remove"></i>
                  </button>
                </el-tooltip>
              </div>
            </div>
            <div class="right-error" v-else>
              <el-button type="danger" @click="removeJob(item)" circle>
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
    <black-dialog v-model:show="showDialog" title="预览视频" hide-footer @cancal="showDialog = false" width="auto">
      <video style="width: 100%; max-height: 90vh;" :src="currentVideoUrl" preload="auto" :autoplay="true" loop="loop"
             muted="muted" v-show="showDialog">
        您的浏览器不支持视频播放
      </video>
    </black-dialog>
  </div>
</template>

<script setup>
import {onMounted, onUnmounted, reactive, ref} from "vue";
import {CircleCloseFilled} from "@element-plus/icons-vue";
import {httpDownload, httpGet, httpPost} from "@/utils/http";
import {checkSession, getClientId} from "@/store/cache";
import {showMessageError, showMessageOK} from "@/utils/dialog";
import {replaceImg} from "@/utils/libs"
import {ElMessage, ElMessageBox} from "element-plus";
import BlackSwitch from "@/components/ui/BlackSwitch.vue";
import Generating from "@/components/ui/Generating.vue";
import BlackDialog from "@/components/ui/BlackDialog.vue";
import {useSharedStore} from "@/store/sharedata";

const showDialog = ref(false)
const currentVideoUrl = ref('')
const row = ref(1)
const images = ref([])

const formData = reactive({
  client_id: getClientId(),
  prompt: '',
  expand_prompt: false,
  loop: false,
  first_frame_img: '',
  end_frame_img: ''
})

const store = useSharedStore()
onMounted(() => {
  checkSession().then(() => {
    fetchData(1)
  })

  store.addMessageHandler("luma", (data) => {
    // 丢弃无关消息
    if (data.channel !== "luma" || data.clientId !== getClientId()) {
      return
    }

    if (data.body === "FINISH" || data.body === "FAIL") {
      fetchData(1)
    }
  })
})

onUnmounted(() => {
  store.removeMessageHandler("luma")
})

const download = (item) => {
  const url = replaceImg(item.video_url)
  const downloadURL = `${process.env.VUE_APP_API_HOST}/api/download?url=${url}`
  // parse filename
  const urlObj = new URL(url);
  const fileName = urlObj.pathname.split('/').pop();
  item.downloading = true
  httpDownload(downloadURL).then(response => {
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

const play = (item) => {
  currentVideoUrl.value = replaceImg(item.video_url)
  showDialog.value = true
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
  httpGet("/api/video/publish", {id: item.id, publish: item.publish}).then(() => {
    ElMessage.success("操作成功")
  }).catch(e => {
    ElMessage.error("操作失败：" + e.message)
  })
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
  httpGet("/api/video/list", {page: page.value, page_size: pageSize.value, type: 'luma'}).then(res => {
    total.value = res.data.total
    loading.value = false
    list.value = res.data.items
    noData.value = list.value.length === 0
  }).catch(() => {
    loading.value = false
    noData.value = true
  })
}

// 创建视频
const create = () => {

  const len = images.value.length;
  if (len) {
    formData.first_frame_img = images.value[0]
    if (len === 2) {
      formData.end_frame_img = images.value[1]
    }
  }

  httpPost("/api/video/luma/create", formData).then(() => {
    fetchData(1)
    showMessageOK("创建任务成功")
  }).catch(e => {
    showMessageError("创建任务失败：" + e.message)
  })
}

const isGenerating = ref(false)
const generatePrompt = () => {
  if (formData.prompt === "") {
    return showMessageError("请输入原始提示词")
  }
  isGenerating.value = true
  httpPost("/api/prompt/image", {prompt: formData.prompt}).then(res => {
    formData.prompt = res.data
    isGenerating.value = false
  }).catch(e => {
    showMessageError("生成提示词失败：" + e.message)
    isGenerating.value = false
  })
}

</script>

<style lang="stylus" scoped>
@import "@/assets/css/luma.styl"
</style>
