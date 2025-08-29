<template>
  <div class="rounded-md p-3 border-2 text-base mb-4 flex items-center" :class="typeClass">
    <div :class="textColor">
      <slot name="icon">
        <i class="iconfont !text-xl" :class="typeIcon"></i>
      </slot>
    </div>
    <div class="ml-3 text-sm leading-relaxed">
      <slot>{{ message }}</slot>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  type: {
    type: String,
    default: 'info',
  },
})

const typeClass = computed(() => {
  return {
    info: 'bg-blue-100 text-blue-500 border-blue-500',
    success: 'bg-green-100 text-green-500 border-green-500',
    warning: 'bg-yellow-100 text-yellow-500 border-yellow-500',
  }[props.type]
})
const typeIcon = computed(() => {
  return {
    info: 'icon-info',
    success: 'icon-success',
    warning: 'icon-warning',
  }[props.type]
})
const textColor = computed(() => {
  return {
    info: 'text-blue-500',
    success: 'text-green-500',
    warning: 'text-yellow-500',
  }[props.type]
})
</script>

<style scoped lang="scss"></style>
