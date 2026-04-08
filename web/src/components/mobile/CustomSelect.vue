<template>
  <div class="custom-select w-full">
    <button
      @click="showPicker = true"
      class="w-full flex items-center justify-between px-4 py-3 bg-gray-50 rounded-lg border border-gray-200 hover:border-blue-300 transition-colors"
    >
      <span>{{ selectedLabel || placeholder || '请选择' }}</span>
      <i class="iconfont icon-arrow-down text-gray-400"></i>
    </button>

    <!-- 选择器弹窗 -->
    <div
      v-if="showPicker"
      class="fixed inset-0 z-50 flex items-end justify-center bg-black bg-opacity-50"
      @click="showPicker = false"
    >
      <div @click.stop class="bg-white rounded-t-2xl w-full max-w-md animate-slide-up">
        <div class="flex items-center justify-between p-4 border-b">
          <h3 class="text-lg font-semibold text-gray-900">{{ title || '请选择' }}</h3>
          <button @click="showPicker = false" class="p-2 hover:bg-gray-100 rounded-full">
            <i class="iconfont icon-close text-gray-500"></i>
          </button>
        </div>
        <div class="max-h-80 overflow-y-auto">
          <CustomSelectOption
            v-for="option in options"
            :key="option.value"
            :option="option"
            :selected="modelValue === option.value"
            @select="onSelect"
          >
            <template #default="slotProps">
              <slot name="option" v-bind="slotProps">
                <div>
                  <span class="text-gray-900 font-medium">{{ slotProps.option.label }}</span>
                  <p v-if="slotProps.option.desc" class="text-sm text-gray-500 mt-1">
                    {{ slotProps.option.desc }}
                  </p>
                </div>
                <div v-if="slotProps.selected" class="text-blue-600">
                  <i class="iconfont icon-success"></i>
                </div>
              </slot>
            </template>
          </CustomSelectOption>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, ref } from 'vue'
import CustomSelectOption from './CustomSelectOption.vue'

// Props
const props = defineProps({
  modelValue: {
    type: [String, Number],
    default: '',
  },
  options: {
    type: Array,
    default: () => [],
  },
  placeholder: {
    type: String,
    default: '请选择',
  },
  title: {
    type: String,
    default: '请选择',
  },
})

// Emits
const emit = defineEmits(['update:modelValue', 'change'])

// Data
const showPicker = ref(false)

// Computed
const selectedLabel = computed(() => {
  const selected = props.options.find((option) => option.value === props.modelValue)
  return selected ? selected.label : ''
})

// Methods
const onSelect = (option) => {
  emit('update:modelValue', option.value)
  emit('change', option)
  showPicker.value = false
}
</script>

<style scoped>
@keyframes slide-up {
  from {
    opacity: 0;
    transform: translateY(100%);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.animate-slide-up {
  animation: slide-up 0.3s ease-out;
}

/* Dark 主题样式 - 按照 theme-dark.scss 的模式 */
:root[data-theme='dark'] .custom-select {
  /* 选择器触发器 */
  button {
    background-color: rgb(31, 41, 55) !important;
    border-color: rgb(75, 85, 99) !important;
    color: rgb(209, 213, 219) !important;

    &:hover {
      border-color: rgb(139, 92, 246) !important;
    }

    .iconfont {
      color: rgb(156, 163, 175) !important;
    }
  }

  /* 选择器弹窗 */
  .fixed {
    .bg-white {
      background-color: rgb(55, 65, 81) !important;
    }

    .border-b {
      border-bottom-color: rgb(75, 85, 99) !important;
    }

    h3 {
      color: rgb(255, 255, 255) !important;
    }

    button {
      background-color: transparent !important;
      color: rgb(156, 163, 175) !important;

      &:hover {
        background-color: rgb(75, 85, 99) !important;
        color: rgb(209, 213, 219) !important;
      }

      .iconfont {
        color: inherit !important;
      }
    }

    .max-h-80 {
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
  }
}
</style>
