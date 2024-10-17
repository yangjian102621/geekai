<template>
  <div class="audio-chat-page">
    <el-button style="margin: 20px" type="primary" size="large" @click="connect()">开始语音对话</el-button>
  </div>
</template>

<script setup>
const connect = () => {

}
</script>

<style scoped lang="stylus">
.audio-chat-page {
  display flex
  flex-flow column
  justify-content center
  align-items center

}

canvas {
  background-color: transparent;
}


</style>