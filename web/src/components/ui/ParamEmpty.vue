<template>
  <div class="flex flex-col justify-center items-center py-8 px-4">
    <!-- 开发中图标和动画 -->
    <div class="relative mb-4">
      <div
        class="w-16 h-16 bg-gradient-to-br from-blue-500 to-purple-600 rounded-full flex items-center justify-center shadow-lg animate-float animate-glow"
      >
        <svg
          class="w-8 h-8 text-white animate-pulse"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z"
          ></path>
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"
          ></path>
        </svg>
      </div>
      <!-- 旋转的装饰环 -->
      <div
        class="absolute inset-0 w-16 h-16 border-2 border-blue-200 rounded-full animate-spin"
        style="animation-duration: 3s"
      ></div>
      <div
        class="absolute inset-1 w-14 h-14 border border-purple-200 rounded-full animate-spin"
        style="animation-duration: 2s; animation-direction: reverse"
      ></div>
    </div>

    <!-- 主标题 -->
    <h3 class="text-lg font-semibold text-gray-700 mb-2">{{ title }}</h3>

    <!-- 开发中提示 -->
    <div class="text-center space-y-2">
      <p class="text-gray-500 text-sm">🚀 {{ statusText }}</p>
      <p class="text-xs text-gray-400 max-w-xs leading-relaxed">
        {{ description }}
      </p>
    </div>

    <!-- 进度条 -->
    <div class="w-48 mt-4">
      <div class="bg-gray-200 rounded-full h-2 overflow-hidden">
        <div
          class="bg-gradient-to-r from-blue-500 to-purple-600 h-full rounded-full progress-bar"
          :style="{ width: progress + '%' }"
        ></div>
      </div>
      <p class="text-xs text-gray-400 text-center mt-1">开发进度 {{ progress }}%</p>
    </div>

    <!-- 装饰性元素 -->
    <div class="flex space-x-1 mt-4">
      <div
        class="w-2 h-2 bg-blue-400 rounded-full animate-bounce"
        style="animation-delay: 0s"
      ></div>
      <div
        class="w-2 h-2 bg-purple-400 rounded-full animate-bounce"
        style="animation-delay: 0.1s"
      ></div>
      <div
        class="w-2 h-2 bg-pink-400 rounded-full animate-bounce"
        style="animation-delay: 0.2s"
      ></div>
    </div>
  </div>
</template>

<script setup>
const props = defineProps({
  // 进度百分比 (0-100)
  progress: {
    type: Number,
    default: 65,
    validator: (value) => value >= 0 && value <= 100,
  },
  // 主标题
  title: {
    type: String,
    default: '参数构建器',
  },
  // 状态文本
  statusText: {
    type: String,
    default: '功能正在开发中',
  },
  // 描述文本
  description: {
    type: String,
    default: '我们正在努力完善当前功能，敬请期待！',
  },
})
</script>

<style lang="scss">
/* 自定义动画效果 */
@keyframes float {
  0%,
  100% {
    transform: translateY(0px);
  }
  50% {
    transform: translateY(-10px);
  }
}

@keyframes glow {
  0%,
  100% {
    box-shadow: 0 0 5px rgba(59, 130, 246, 0.5);
  }
  50% {
    box-shadow: 0 0 20px rgba(59, 130, 246, 0.8), 0 0 30px rgba(147, 51, 234, 0.6);
  }
}

.animate-float {
  animation: float 3s ease-in-out infinite;
}

.animate-glow {
  animation: glow 2s ease-in-out infinite;
}

/* 进度条动画 */
@keyframes progress {
  0% {
    width: 0%;
  }
  100% {
    width: v-bind(progress + '%');
  }
}

.progress-bar {
  animation: progress 2s ease-out forwards;
}
</style>
