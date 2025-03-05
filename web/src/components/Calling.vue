<template>
  <!--拨号组件-->
  <el-container class="calling-container" :style="{height: height}">
    <div class="phone-container">
      <div class="signal"></div>
      <div class="signal"></div>
      <div class="signal"></div>
      <div class="phone"></div>
    </div>
    <div class="status-text">{{ text }}</div>
  </el-container>
</template>

<script setup>
import {onMounted, ref} from "vue";

const fullText = "正在接通中...";
const text = ref("")
let index = 0;
const props = defineProps({
  height: {
    type: String,
    default: '100vh'
  }
})

function typeText() {
  if (index < fullText.length) {
    text.value += fullText[index];
    index++;
    setTimeout(typeText, 300); // 每300毫秒显示一个字
  } else {
    setTimeout(() => {
      text.value = '';
      index = 0;
      typeText();
    }, 1000); // 等待1秒后重新开始
  }
}

onMounted(() => {
  typeText()
})

</script>

<style scoped lang="stylus">

.calling-container {
  background-color: #000;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  margin: 0;
  overflow: hidden;
  font-family: Arial, sans-serif;
  width 100vw

  .phone-container {
    position: relative;
    width: 200px;
    height: 200px;
  }

  .phone {
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    width: 60px;
    height: 60px;
    background-color: #00ffcc;
    mask: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24'%3E%3Cpath d='M20 15.5c-1.25 0-2.45-.2-3.57-.57a1.02 1.02 0 0 0-1.02.24l-2.2 2.2a15.074 15.074 0 0 1-6.59-6.59l2.2-2.2c.27-.27.35-.68.24-1.02a11.36 11.36 0 0 1-.57-3.57c0-.55-.45-1-1-1H4c-.55 0-1 .45-1 1 0 9.39 7.61 17 17 17 .55 0 1-.45 1-1v-3.5c0-.55-.45-1-1-1zM5.03 5h1.5c.07.89.22 1.76.46 2.59l-1.2 1.2c-.41-1.2-.67-2.47-.76-3.79zM19 18.97c-1.32-.09-2.59-.35-3.8-.75l1.2-1.2c.85.24 1.72.39 2.6.45v1.5z'/%3E%3C/svg%3E") no-repeat 50% 50%;
    mask-size: cover;
    -webkit-mask: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24'%3E%3Cpath d='M20 15.5c-1.25 0-2.45-.2-3.57-.57a1.02 1.02 0 0 0-1.02.24l-2.2 2.2a15.074 15.074 0 0 1-6.59-6.59l2.2-2.2c.27-.27.35-.68.24-1.02a11.36 11.36 0 0 1-.57-3.57c0-.55-.45-1-1-1H4c-.55 0-1 .45-1 1 0 9.39 7.61 17 17 17 .55 0 1-.45 1-1v-3.5c0-.55-.45-1-1-1zM5.03 5h1.5c.07.89.22 1.76.46 2.59l-1.2 1.2c-.41-1.2-.67-2.47-.76-3.79zM19 18.97c-1.32-.09-2.59-.35-3.8-.75l1.2-1.2c.85.24 1.72.39 2.6.45v1.5z'/%3E%3C/svg%3E") no-repeat 50% 50%;
    -webkit-mask-size: cover;
    animation: shake 0.5s ease-in-out infinite;
  }

  .signal {
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    width: 100px;
    height: 100px;
    border: 2px dashed #00ffcc;
    border-radius: 50%;
    opacity: 0;
    animation: signal 2s linear infinite;
  }

  .signal:nth-child(2) {
    animation-delay: 0.5s;
  }

  .signal:nth-child(3) {
    animation-delay: 1s;
  }

  .status-text {
    color: #00ffcc;
    font-size: 18px;
    margin-top: 20px;
    height: 1.2em;
    overflow: hidden;
  }

  @keyframes shake {
    0%, 100% { transform: translate(-50%, -50%) rotate(0deg); }
    25% { transform: translate(-52%, -48%) rotate(-5deg); }
    75% { transform: translate(-48%, -52%) rotate(5deg); }
  }

  @keyframes signal {
    0% {
      width: 60px;
      height: 60px;
      opacity: 1;
    }
    100% {
      width: 200px;
      height: 200px;
      opacity: 0;
    }
  }
}

</style>