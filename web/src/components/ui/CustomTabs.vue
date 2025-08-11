<template>
  <div class="w-full">
    <div
      class="relative bg-gray-100 rounded-lg py-1.5 mb-3 px-2 overflow-hidden custom-tabs-header"
      ref="tabsHeader"
    >
      <!-- 左滑动指示器 -->
      <div
        v-show="canScrollLeft"
        class="absolute left-1 top-1/2 -translate-y-1/2 z-30 w-6 h-6 bg-white/95 backdrop-blur-sm rounded-full shadow-sm border border-gray-200/50 flex items-center justify-center cursor-pointer hover:bg-white hover:shadow-md hover:scale-105 transition-all duration-200 group scroll-indicator scroll-indicator-left"
        @click="scrollLeft"
      >
        <svg
          class="w-3 h-3 text-gray-500 group-hover:text-purple-600 transition-colors duration-200"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2.5"
            d="M15 19l-7-7 7-7"
          ></path>
        </svg>
      </div>

      <!-- 右滑动指示器 -->
      <div
        v-show="canScrollRight"
        class="absolute right-1 top-1/2 -translate-y-1/2 z-30 w-6 h-6 bg-white/95 backdrop-blur-sm rounded-full shadow-sm border border-gray-200/50 flex items-center justify-center cursor-pointer hover:bg-white hover:shadow-md hover:shadow-md hover:scale-105 transition-all duration-200 group scroll-indicator scroll-indicator-right"
        @click="scrollRight"
      >
        <svg
          class="w-3 h-3 text-gray-500 group-hover:text-purple-600 transition-colors duration-200"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2.5"
            d="M9 5l7 7-7 7"
          ></path>
        </svg>
      </div>

      <div
        class="flex whitespace-nowrap overflow-x-auto scrollbar-hide"
        ref="tabsContainer"
        @scroll="checkScrollPosition"
      >
        <div
          class="flex-shrink-0 text-center py-1 px-2 font-medium text-gray-700 cursor-pointer transition-all duration-300 rounded-md relative z-20 hover:text-purple-600 custom-tab-item"
          v-for="(tab, index) in panes"
          :key="tab.name"
          :class="{
            '!text-purple-600 bg-white shadow-sm custom-tab-active': modelValue === tab.name,
            'hover:bg-gray-50': modelValue !== tab.name,
          }"
          @click="handleTabClick(tab.name, index)"
          ref="tabItems"
        >
          <component v-if="tab.labelSlot" :is="{ render: () => tab.labelSlot() }" />
          <template v-else>
            {{ tab.label }}
          </template>
        </div>
      </div>
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
const tabsContainer = ref(null)
const tabItems = ref([])
const panes = ref([])

// 滑动状态
const canScrollLeft = ref(false)
const canScrollRight = ref(false)

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

// 检查滑动位置状态
const checkScrollPosition = () => {
  if (!tabsContainer.value) return

  const container = tabsContainer.value
  canScrollLeft.value = container.scrollLeft > 0
  canScrollRight.value = container.scrollLeft < container.scrollWidth - container.clientWidth
}

// 向左滑动
const scrollLeft = () => {
  if (!tabsContainer.value) return

  const container = tabsContainer.value
  const scrollAmount = Math.min(200, container.scrollLeft) // 每次滑动200px或剩余距离

  container.scrollTo({
    left: container.scrollLeft - scrollAmount,
    behavior: 'smooth',
  })
}

// 向右滑动
const scrollRight = () => {
  if (!tabsContainer.value) return

  const container = tabsContainer.value
  const maxScroll = container.scrollWidth - container.clientWidth
  const scrollAmount = Math.min(200, maxScroll - container.scrollLeft) // 每次滑动200px或剩余距离

  container.scrollTo({
    left: container.scrollLeft + scrollAmount,
    behavior: 'smooth',
  })
}

const handleTabClick = (tabName, index) => {
  emit('update:modelValue', tabName)
  emit('tab-click', tabName, index)
  scrollToTab(index)
}

// 简化后的滚动到指定tab的函数
const scrollToTab = async (activeIndex) => {
  await nextTick()
  if (!tabsContainer.value || !tabItems.value || tabItems.value.length === 0) {
    return
  }

  const activeTab = tabItems.value[activeIndex]
  if (!activeTab) {
    return
  }

  const container = tabsContainer.value
  const tabRect = activeTab.getBoundingClientRect()
  const containerRect = container.getBoundingClientRect()

  // 计算tab相对于容器的位置
  const tabLeft = tabRect.left - containerRect.left
  const tabRight = tabLeft + tabRect.width
  const containerWidth = containerRect.width

  // 检查tab是否在可视区域内（增加一些容错空间）
  const tolerance = 4 // 4px的容错空间
  const isVisible = tabLeft >= -tolerance && tabRight <= containerWidth + tolerance

  if (!isVisible) {
    let scrollLeft = container.scrollLeft

    if (tabLeft < -tolerance) {
      // tab在左侧不可见，滚动到tab的起始位置
      scrollLeft += tabLeft - 12 // 留出12px的边距
    } else if (tabRight > containerWidth + tolerance) {
      // tab在右侧不可见，滚动到tab的结束位置
      scrollLeft += tabRight - containerWidth + 12 // 留出12px的边距
    }

    // 确保滚动位置不超出边界
    scrollLeft = Math.max(0, Math.min(scrollLeft, container.scrollWidth - containerWidth))

    // 平滑滚动到目标位置
    container.scrollTo({
      left: scrollLeft,
      behavior: 'smooth',
    })
  }

  // 更新滑动状态
  setTimeout(checkScrollPosition, 300)
}

// 监听 modelValue 变化，滚动到tab
watch(
  () => props.modelValue,
  (newValue) => {
    const activeIndex = panes.value.findIndex((pane) => pane.name === newValue)
    if (activeIndex !== -1) {
      scrollToTab(activeIndex)
    }
  }
)

// 监听 panes 变化，当tab数量变化时重新计算
watch(
  () => panes.value.length,
  () => {
    nextTick(() => {
      const activeIndex = panes.value.findIndex((pane) => pane.name === props.modelValue)
      if (activeIndex !== -1) {
        scrollToTab(activeIndex)
      }
      // 检查滑动状态
      setTimeout(checkScrollPosition, 100)
    })
  }
)

onMounted(() => {
  // 初始化时滚动到选中tab
  nextTick(() => {
    const activeIndex = panes.value.findIndex((pane) => pane.name === props.modelValue)
    if (activeIndex !== -1) {
      scrollToTab(activeIndex)
    }
    // 检查初始滑动状态
    setTimeout(checkScrollPosition, 100)
  })

  // 监听窗口大小变化，重新计算滚动
  const handleResize = () => {
    const activeIndex = panes.value.findIndex((pane) => pane.name === props.modelValue)
    if (activeIndex !== -1) {
      scrollToTab(activeIndex)
    }
    // 检查滑动状态
    setTimeout(checkScrollPosition, 100)
  }

  window.addEventListener('resize', handleResize)

  // 清理事件监听器
  return () => {
    window.removeEventListener('resize', handleResize)
  }
})
</script>

<style scoped>
.scrollbar-hide {
  -ms-overflow-style: none; /* IE and Edge */
  scrollbar-width: none; /* Firefox */
}

.scrollbar-hide::-webkit-scrollbar {
  display: none; /* Chrome, Safari and Opera */
}

/* 确保标签页容器有足够的内边距 */
.overflow-x-auto {
  padding: 0 4px;
}

/* 优化标签页间距 */
.flex-shrink-0 {
  margin-right: 4px;
}

.flex-shrink-0:last-child {
  margin-right: 0;
}

/* 添加平滑滚动效果 */
.overflow-x-auto {
  scroll-behavior: smooth;
}

/* 滑动指示器样式 */
.absolute {
  transition: opacity 0.2s ease-in-out;
}

/* Dark 主题样式 - 按照 theme-dark.scss 的模式 */
:root[data-theme='dark'] .custom-tabs-header {
  background-color: rgb(31, 41, 55);
}

:root[data-theme='dark'] .custom-tab-item {
  color: rgb(209, 213, 219);
}

:root[data-theme='dark'] .custom-tab-item:hover {
  color: rgb(196, 181, 253);
}

:root[data-theme='dark'] .custom-tab-active {
  background-color: rgb(55, 65, 81);
  color: rgb(196, 181, 253);
}

:root[data-theme='dark'] .custom-tab-item:hover:not(.custom-tab-active) {
  background-color: rgb(75, 85, 99);
}

:root[data-theme='dark'] .scroll-indicator {
  background-color: rgba(55, 65, 81, 0.95);
  border-color: rgb(75, 85, 99);
}

:root[data-theme='dark'] .scroll-indicator:hover {
  background-color: rgb(75, 85, 99);
}

:root[data-theme='dark'] .scroll-indicator svg {
  color: rgb(156, 163, 175);
}

:root[data-theme='dark'] .scroll-indicator:hover svg {
  color: rgb(196, 181, 253);
}
</style>
