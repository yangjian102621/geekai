<template>
  <!-- 倒计时组件 -->
  <div class="countdown">
    <el-tag size="large" :type="type">{{ timerStr }}</el-tag>
  </div>
</template>
<script setup>

import {onMounted, ref, watch} from "vue";

// eslint-disable-next-line no-undef
const props = defineProps({
  second: Number,
  type: {
    type: String,
    default: ""
  }
});

// eslint-disable-next-line no-undef
const emits = defineEmits(['timeout']);
const counter = ref(props.second)
const timerStr = ref("")
const handler = ref(null)

watch(() => props.second, (newVal) => {
  counter.value = newVal
  resetTimer()
});

onMounted(() => {
  resetTimer()
})

const resetTimer = () => {
  if (handler.value) {
    clearInterval(handler.value)
  }

  counter.value = props.second
  formatTimer(counter.value)
  handler.value = setInterval(() => {
    formatTimer(counter.value)
    if (counter.value === 0) {
      clearInterval(handler.value)
      emits("timeout")
    }
    counter.value--
  }, 1000)
}

const formatTimer = (secs) => {
  const timer = []
  let hour, min
  // 计算小时
  if (secs > 3600) {
    hour = Math.floor(secs / 3600)
    if (hour < 10) {
      hour = "0" + hour
    }
    secs = secs % 3600
    timer.push(hour + " 时 ")
  } else {
    timer.push("00 时 ")
  }
  // 计算分钟
  if (secs > 60) {
    min = Math.floor(secs / 60)
    if (min < 10) {
      min = "0" + min
    }
    secs = secs % 60
    timer.push(min + " 分 ")
  } else {
    timer.push("00 分 ")
  }
  // 计算秒数
  if (secs < 10) {
    secs = "0" + secs
  }
  timer.push(secs + " 秒")
  timerStr.value = timer.join("")
}

// eslint-disable-next-line no-undef
defineExpose({resetTimer})
</script>

<style lang="stylus">
.countdown {
  display flex

  .el-tag--large {
    .el-tag__content {
      font-size 14px
    }
  }
}

</style>