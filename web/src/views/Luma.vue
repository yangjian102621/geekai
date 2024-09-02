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
              v-model="prompt"
              placeholder="请输入提示词或者上传图片"
              autofocus>
                      </textarea>
          <div class="send-icon">
            <i class="iconfont icon-send"></i>
          </div>
        </div>

        <div class="params">
          <div class="item-group">
            <span class="label">循环参考图</span>
            <el-switch  v-model="loop" size="small" style="--el-switch-on-color:#BF78BF;" />
          </div>
          <div class="item-group">
            <span class="label">提示词优化</span>
            <el-switch  v-model="promptExtend" size="small" style="--el-switch-on-color:#BF78BF;" />
          </div>
        </div>
      </div>
    </div>

    <el-container class="video-container">
      <h2 class="h-title">你的作品</h2>

      <el-row :gutter="20" class="videos">
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
      </el-row>
    </el-container>
  </div>
</template>

<script setup>
import {ref} from "vue";
import {CircleCloseFilled} from "@element-plus/icons-vue";
import {httpDownload, httpPost} from "@/utils/http";
import {showMessageError} from "@/utils/dialog";
import {ElMessage} from "element-plus";

const row = ref(1)
const prompt = ref('')
const loop = ref(false)
const promptExtend = ref(false)
const images = ref([])

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


</script>

<style lang="stylus" scoped>
@import "@/assets/css/luma.styl"
</style>
