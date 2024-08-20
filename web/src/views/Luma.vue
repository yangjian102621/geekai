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

    <h2 class="h-title">你的作品</h2>
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

const list = ref([
  {
    id: 1,
    name: '1',
    url: 'http://localhost/download/xmind.mp4'
  }
])


</script>

<style lang="stylus" scoped>
@import "@/assets/css/luma.styl"
</style>
