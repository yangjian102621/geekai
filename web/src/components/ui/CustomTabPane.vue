<template>
  <div v-if="active" class="custom-tab-pane">
    <slot></slot>
  </div>
</template>

<script setup>
  import { computed, inject, useSlots } from 'vue'

  const props = defineProps({
    label: {
      type: String,
      required: false,
    },
    name: {
      type: String,
      required: true,
    },
  })

  const slots = useSlots()

  // 从父组件注入当前激活的 tab
  const currentTab = inject('currentTab', '')

  const active = computed(() => {
    return currentTab.value === props.name
  })

  // 向父组件提供当前 pane 的信息，优先使用 labelSlot
  const parentRegisterPane = inject('registerPane', () => {})

  // 立即注册，不要等到 onMounted
  parentRegisterPane({
    name: props.name,
    label: props.label || '', // 如果没有传 label 则使用空字符串
    labelSlot: slots.label,
  })
</script>

<style scoped>
  .custom-tab-pane {
    width: 100%;
  }
</style>
