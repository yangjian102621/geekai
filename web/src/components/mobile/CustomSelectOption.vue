<template>
  <!--
    CustomSelectOption 组件
    Props:
      - option: 选项对象，必需，包含 label/desc/value 等属性
      - selected: 是否为当前选中项
    Emits:
      - select(option): 选中该项时触发
    Slots:
      - 默认插槽（default）：用于自定义 option 内容，slotProps: { option, selected }
        示例：
        <template #option="{ option, selected }">
          <div>{{ option.label }}</div>
          <div v-if="selected">✔</div>
        </template>
  -->
  <div
    class="flex items-center justify-between p-4 hover:bg-gray-50 cursor-pointer transition-colors border-b last:border-b-0"
    @click="$emit('select', option)"
  >
    <slot :option="option" :selected="selected">
      <div>
        <span class="text-gray-900 font-medium">{{ option.label }}</span>
        <p v-if="option.desc" class="text-sm text-gray-500 mt-1">{{ option.desc }}</p>
      </div>
      <div v-if="selected" class="text-blue-600">
        <i class="iconfont icon-success"></i>
      </div>
    </slot>
  </div>
</template>

<script setup>
import { defineProps, defineEmits } from 'vue'

const props = defineProps({
  option: {
    type: Object,
    required: true,
  },
  selected: {
    type: Boolean,
    default: false,
  },
})

const emit = defineEmits(['select'])
</script>

<style scoped>
/* Dark 主题样式 - 按照 theme-dark.scss 的模式 */
:root[data-theme='dark'] .flex {
  background-color: transparent !important;
  border-bottom-color: rgb(75, 85, 99) !important;

  &:hover {
    background-color: rgb(75, 85, 99) !important;
  }

  &:last-child {
    border-bottom-color: transparent !important;
  }

  .text-gray-900 {
    color: rgb(209, 213, 219) !important;
  }

  .text-gray-500 {
    color: rgb(156, 163, 175) !important;
  }

  .text-blue-600 {
    color: rgb(139, 92, 246) !important;
  }
}
</style>
