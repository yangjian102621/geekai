<template>
  <div class="flex flex-col justify-center items-center py-10 px-4 text-center">
    <div class="mb-4 flex items-center justify-center">
      <div class="w-14 h-14 rounded-full bg-gray-100 flex items-center justify-center shadow">
        <svg class="w-7 h-7 text-gray-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M12 11c1.657 0 3-1.567 3-3.5S13.657 4 12 4 9 5.567 9 7.5 10.343 11 12 11z"
          />
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M6 20v-1a5 5 0 015-5h2a5 5 0 015 5v1"
          />
        </svg>
      </div>
    </div>

    <h3 class="text-lg font-semibold text-gray-800 mb-1">{{ title }}</h3>
    <p class="text-gray-500 text-sm">{{ statusText }}</p>
    <p v-if="description" class="text-xs text-gray-400 mt-2 max-w-sm leading-relaxed">
      {{ description }}
    </p>
  </div>

  <!-- 允许使用插槽扩展，比如放置返回按钮或客服入口 -->
  <slot />
</template>

<script setup>
const props = defineProps({
  title: {
    type: String,
    default: '功能暂未开放',
  },
  statusText: {
    type: String,
    default: '该功能暂不可用，请稍后重试或联系管理员',
  },
  description: {
    type: String,
    default: '',
  },
})
</script>

<style lang="scss">
/* 保留最小样式，去除无关动画与进度条 */
</style>
