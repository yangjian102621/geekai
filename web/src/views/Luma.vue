<template>
  <div class="page-luma">
    <div class="prompt-box">
      <div class="images">
        <div v-for="img in images" class="item">
          <el-image :src="img" fit="cover"/>
          <el-icon><CircleCloseFilled /></el-icon>
        </div>
      </div>
      <div class="prompt-container">
        <div class="input-container">
          <div class="upload-icon">
            <i class="iconfont icon-image"></i>
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

      </div>
    </div>

    <el-container class="video-container">
      <h2 class="h-title">你的作品</h2>

      <el-row :gutter="20" class="videos">
        <el-col :span="8" class="item" :key="item.id" v-for="item in videos">
          <div class="video-box" @mouseover="item.playing = true" @mouseout="item.playing = false">
            <img :src="item.cover"  :alt="item.name" v-if="!item.playing"/>
            <video :src="item.url"  preload="auto" :autoplay="true" loop="loop" muted="muted" v-else>
              您的浏览器不支持视频播放
            </video>
          </div>
          <div class="video-name">{{item.name}}</div>
        </el-col>
      </el-row>
    </el-container>
  </div>
</template>

<script setup>
import {ref} from "vue";
import {CircleCloseFilled} from "@element-plus/icons-vue";

const row = ref(1)
const prompt = ref('')
const images = ref([
    "http://nk.img.r9it.com/chatgpt-plus/1719371605709871.jpg",
    "http://nk.img.r9it.com/chatgpt-plus/1719371605709871.jpg"
])

const videos = ref([
  {
    id: 1,
    name: 'a dancing girl',
    url: 'http://localhost/download/xmind.mp4',
    cover: 'https://storage.cdn-luma.com/dream_machine/d133794f-3124-4059-a9f2-e5fed79f0d5b/video_0_thumb.jpg',
    playing: false
  },
  {
    id: 1,
    name: 'a dancing girl',
    url: 'http://localhost/download/dancing.mp4',
    cover: 'https://storage.cdn-luma.com/dream_machine/92efa55a-f381-4161-a999-54f8fe460fca/video_1_thumb.jpg',
    playing: false
  },
  {
    id: 1,
    name: 'a dancing girl',
    url: 'http://localhost/download/xmind.mp4',
    cover: 'https://storage.cdn-luma.com/dream_machine/d133794f-3124-4059-a9f2-e5fed79f0d5b/video_0_thumb.jpg',
    playing: false
  }
])


</script>

<style lang="stylus" scoped>
@import "@/assets/css/luma.styl"
</style>
