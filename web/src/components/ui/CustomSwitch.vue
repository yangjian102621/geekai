<template>
  <div class="custom-switch" :class="{ 'is-active': modelValue }" @click="toggleSwitch">
    <div class="switch-track" :style="trackStyle" ref="trackRef">
      <div class="switch-thumb" :style="thumbStyle"></div>
      <div class="switch-text inactive-text">
        <slot name="inactive-text"></slot>
      </div>
      <div class="switch-text active-text">
        <slot name="active-text"></slot>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, defineEmits, defineProps, nextTick, onMounted, ref } from 'vue'

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false,
  },
  disabled: {
    type: Boolean,
    default: false,
  },
  activeColor: {
    type: String,
    default: '#67c23a', // 默认绿色
  },
  inactiveColor: {
    type: String,
    default: '#f56c6c', // 默认红色
  },
  width: {
    type: [String, Number],
    default: null, // 默认不设置固定宽度，使用min-width
  },
  size: {
    type: String,
    default: 'default', // small, default, large
    validator: (value) => ['small', 'default', 'large'].includes(value),
  },
})

const emit = defineEmits(['update:modelValue', 'change'])

// 轨道宽度引用
const trackRef = ref(null)
const trackWidth = ref(120) // 默认宽度

// 获取轨道实际宽度
const updateTrackWidth = () => {
  if (trackRef.value) {
    trackWidth.value = trackRef.value.offsetWidth
  }
}

// 计算尺寸相关样式
const sizeConfig = computed(() => {
  const configs = {
    small: {
      height: 24,
      thumbSize: 20,
      thumbMargin: 2,
      fontSize: 12,
      padding: 12,
    },
    default: {
      height: 28,
      thumbSize: 24,
      thumbMargin: 2,
      fontSize: 13,
      padding: 14,
    },
    large: {
      height: 32,
      thumbSize: 28,
      thumbMargin: 2,
      fontSize: 14,
      padding: 16,
    },
  }
  return configs[props.size]
})

// 计算轨道样式
const trackStyle = computed(() => {
  const backgroundColor = props.modelValue ? props.activeColor : props.inactiveColor
  const config = sizeConfig.value
  const style = {
    backgroundColor: backgroundColor,
    height: `${config.height}px`,
    padding: `0 ${config.padding}px`,
  }

  // 如果传入了width属性，则设置固定宽度
  if (props.width !== null) {
    const widthValue = typeof props.width === 'number' ? `${props.width}px` : props.width
    style.width = widthValue
    style.minWidth = widthValue
  }

  return style
})

// 计算滑块样式 - 使用像素值而不是calc()
const thumbStyle = computed(() => {
  const config = sizeConfig.value
  const thumbWidth = config.thumbSize
  const thumbMargin = config.thumbMargin
  const maxTranslateX = trackWidth.value - thumbWidth - thumbMargin * 2

  const transform = props.modelValue ? `translateX(${maxTranslateX}px)` : 'translateX(0px)'
  console.log('Track width:', trackWidth.value, 'Transform:', transform)

  return {
    transform: transform,
    width: `${thumbWidth}px`,
    height: `${thumbWidth}px`,
    top: `${thumbMargin}px`,
    left: `${thumbMargin}px`,
  }
})

const toggleSwitch = () => {
  if (props.disabled) return

  const newValue = !props.modelValue
  emit('update:modelValue', newValue)
  emit('change', newValue)

  // 切换后更新宽度
  nextTick(() => {
    updateTrackWidth()
  })
}

// 组件挂载后获取轨道宽度
onMounted(() => {
  nextTick(() => {
    updateTrackWidth()
  })
})
</script>

<style scoped lang="scss">
.custom-switch {
  display: inline-block;
  cursor: pointer;
  user-select: none;

  &.is-active {
    .switch-track {
      .inactive-text {
        opacity: 0;
        transform: translateX(-20px);
      }

      .active-text {
        opacity: 1;
        transform: translateX(0);
      }
    }
  }

  &:not(.is-active) {
    .switch-track {
      .inactive-text {
        opacity: 1;
        transform: translateX(0);
      }

      .active-text {
        opacity: 0;
        transform: translateX(20px);
      }
    }
  }

  &:disabled {
    cursor: not-allowed;
    opacity: 0.6;
  }
}

.switch-track {
  position: relative;
  min-width: 120px; // 最小宽度
  border-radius: 16px;
  transition: background-color 0.3s ease;
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.switch-thumb {
  position: absolute;
  background-color: white;
  border-radius: 50%;
  transition: transform 0.3s ease;
  z-index: 2;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
}

.switch-text {
  position: relative; // 改为相对定位
  color: white;
  font-weight: 500;
  white-space: nowrap;
  transition: all 0.3s ease;
  z-index: 1;
  flex: 1; // 让文字占据剩余空间
  text-align: center; // 文字居中
}

.inactive-text {
  opacity: 1;
  transform: translateX(0);
}

.active-text {
  opacity: 0;
  transform: translateX(20px);
  position: absolute; // 绝对定位，避免影响布局
  left: 0;
  right: 0;
}

// 响应式设计
@media (max-width: 768px) {
  .switch-track {
    min-width: 100px;
  }
}
</style>
