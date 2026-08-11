<template>
  <div class="w-full">
    <div
      class="relative flex bg-gray-100 rounded-lg py-1 mb-2.5 overflow-hidden"
      ref="tabsHeader"
    >
      <div
        class="flex-1 text-center py-1.5 px-3 font-medium text-gray-700 cursor-pointer transition-colors duration-300 rounded-md relative z-20 hover:text-purple-600"
        v-for="(tab, index) in panes"
        :key="tab.name"
        :class="{ '!text-purple-600': modelValue === tab.name }"
        @click="handleTabClick(tab.name, index)"
        ref="tabItems"
      >
        <component
          v-if="tab.labelSlot"
          :is="{ render: () => tab.labelSlot() }"
        />
        <template v-else>
          {{ tab.label }}
        </template>
      </div>
      <div
        class="absolute top-1 bottom-1 bg-white rounded-md shadow-sm transition-all duration-300 ease-out z-10"
        :style="indicatorStyle"
        ref="indicator"
      ></div>
    </div>
    <div>
      <slot></slot>
    </div>
  </div>
</template>

<script setup>
  import { computed, nextTick, onMounted, provide, ref, watch } from 'vue'

  const props = defineProps({
    modelValue: {
      type: String,
      required: true,
    },
  })

  const emit = defineEmits(['update:modelValue', 'tab-click'])

  const tabsHeader = ref(null)
  const tabItems = ref([])
  const indicator = ref(null)
  const panes = ref([])

  const indicatorStyle = ref({
    transform: 'translateX(0px)',
    width: '0px',
  })

  // 提供当前激活的 tab 给子组件
  provide(
    'currentTab',
    computed(() => props.modelValue)
  )

  // 提供注册 pane 的方法给子组件
  provide('registerPane', (pane) => {
    // 检查是否已经存在相同 name 的 pane，避免重复注册
    const existingIndex = panes.value.findIndex((p) => p.name === pane.name)
    if (existingIndex === -1) {
      // 不存在则添加
      panes.value.push(pane)
    } else {
      // 存在则更新
      panes.value[existingIndex] = pane
    }
  })

  const handleTabClick = (tabName, index) => {
    emit('update:modelValue', tabName)
    emit('tab-click', tabName, index)
    updateIndicator(index)
  }

  const updateIndicator = async (activeIndex) => {
    await nextTick()
    if (tabItems.value && tabItems.value.length > 0 && tabsHeader.value) {
      const container = tabsHeader.value
      const containerWidth = container.offsetWidth
      const padding = 8 // p-1 左右内边距总和 (4px * 2)
      const availableWidth = containerWidth - padding
      const tabCount = tabItems.value.length
      const tabWidth = availableWidth / tabCount
      const leftPosition = activeIndex * tabWidth + 4 // 加上左内边距

      indicatorStyle.value = {
        transform: `translateX(${leftPosition}px)`,
        width: `${tabWidth}px`,
      }
    }
  }

  // 监听 modelValue 变化，更新指示器位置
  watch(
    () => props.modelValue,
    (newValue) => {
      const activeIndex = panes.value.findIndex(
        (pane) => pane.name === newValue
      )
      if (activeIndex !== -1) {
        updateIndicator(activeIndex)
      }
    }
  )

  onMounted(() => {
    // 初始化指示器位置
    nextTick(() => {
      const activeIndex = panes.value.findIndex(
        (pane) => pane.name === props.modelValue
      )
      if (activeIndex !== -1) {
        updateIndicator(activeIndex)
      }
    })
  })
</script>
