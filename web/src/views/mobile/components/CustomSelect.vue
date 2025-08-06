<template>
  <div class="custom-select">
    <label v-if="label" class="block text-gray-700 font-medium mb-2">{{ label }}</label>
    <button
      @click="showPicker = true"
      class="w-full flex items-center justify-between px-4 py-3 bg-gray-50 rounded-lg border border-gray-200 hover:border-blue-300 transition-colors"
    >
      <span class="text-gray-900">{{ selectedLabel || placeholder || '请选择' }}</span>
      <i class="iconfont icon-down text-gray-400"></i>
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
import { ref, computed, watch } from 'vue'
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
  label: {
    type: String,
    default: '',
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
</style>
