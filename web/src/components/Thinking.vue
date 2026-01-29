<template>
  <div class="loading-container" :style="containerStyle">
    <div
      v-for="i in 3"
      :key="i"
      class="dot"
      :style="[dotStyle, { animationDelay: `${(i - 1) * 0.2}s` }]"
    ></div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  // 点的大小(px)
  size: {
    type: Number,
    default: 8,
  },
  // 点之间的间距(px)
  spacing: {
    type: Number,
    default: 4,
  },
  // 动画持续时间(秒)
  duration: {
    type: Number,
    default: 1.4,
  },
})

// 计算样式
const containerStyle = computed(() => ({
  height: `${props.size * 2}px`,
}))

const dotStyle = computed(() => ({
  width: `${props.size}px`,
  height: `${props.size}px`,
  margin: `0 ${props.spacing}px`,
  background: `var(--el-text-color-regular)`,
  animationDuration: `${props.duration}s`,
}))
</script>
<style scoped>
.loading-container {
  display: flex;
  align-items: center;
  justify-content: center;
}

.dot {
  border-radius: 50%;
  display: inline-block;
  animation: loading ease-in-out infinite;
}

@keyframes loading {
  0% {
    opacity: 0.2;
    transform: scale(0.8);
  }
  20% {
    opacity: 1;
    transform: scale(1.2);
  }
  40% {
    opacity: 0.2;
    transform: scale(0.8);
  }
  100% {
    opacity: 0.2;
    transform: scale(0.8);
  }
}
</style>
