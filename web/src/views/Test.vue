<template>
  <div class="audio-chat-page">
    <el-button style="margin: 20px" type="primary" size="large" @click="connect()">开始语音对话</el-button>

    <el-dialog v-model="showDialog" title="语音通话" :before-close="close">
      <realtime-conversation  @close="showDialog = false" ref="conversationRef" :height="dialogHeight+'px'" />
    </el-dialog>
  </div>
</template>

<script setup>
import {nextTick, ref} from 'vue';
import RealtimeConversation from "@/components/RealtimeConversation.vue";

const showDialog = ref(false);
const dialogHeight = ref(window.innerHeight - 75);
const conversationRef = ref(null);
const connect = () => {
  showDialog.value = true;
  nextTick(() => {
    conversationRef.value.connect()
  })
}
const close = () => {
  showDialog.value = false;
  conversationRef.value.hangUp()
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