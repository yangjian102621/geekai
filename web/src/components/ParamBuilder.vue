<template>
  <div class="param-builder flex flex-col">
    <ParamEmpty
      v-if="items.length === 0"
      :progress="progress"
      :title="title"
      :status-text="statusText"
      :description="description"
    />
    <div v-else class="flex flex-col w-full space-y-5">
      <el-select
        v-model="selectedModel"
        placeholder="请选择模型"
        @change="changeModel"
        popper-class="model-select"
        value-key="name"
      >
        <template #prefix>
          <i class="iconfont icon-model"></i>
        </template>

        <el-option v-for="item in items" :key="item.name" :label="item.name" :value="item">
          <div class="flex justify-start">
            <span
              class="flex items-center justify-center text-white !text-xl model-version mr-2 w-[40px] h-[40px] rounded-lg"
              >{{ item.version }}</span
            >
            <div class="flex !items-start flex-col py-2 space-y-1">
              <span class="label text-sm">{{ item.name }}</span>
              <div class="whitespace-pre-line">
                <span
                  class="text-xs text-gray-500 break-words line-clamp-1 max-w-[200px]"
                  :title="item.label"
                  >{{ item.label }}</span
                >
              </div>
            </div>
          </div>
        </el-option>
      </el-select>

      <div v-for="param in selectedModel.params" :key="param.name" class="w-full">
        <div class="w-full flex flex-col !items-start space-y-2" v-if="param.type === 'switch'">
          <div class="w-full flex justify-between">
            <label class="label font-bold">{{ param.label }}</label>
            <el-switch v-model="modelValue[param.name]" size="large" />
          </div>
          <p v-if="param.info" class="text-xs text-gray-500 mb-1">{{ param.info }}</p>
        </div>
        <div class="w-full flex flex-col !items-start space-y-2" v-else>
          <label class="label font-bold">
            {{ param.label }}
            <span v-if="param.required" class="text-red-500 ml-1">*</span>
          </label>
          <p v-if="param.info" class="text-xs text-gray-500 mb-1">{{ param.info }}</p>
          <div class="flex w-full">
            <el-input
              v-if="param.type === 'text'"
              v-model="modelValue[param.name]"
              :placeholder="param.placeholder"
            />
            <el-input-number
              v-if="param.type === 'number'"
              v-model="modelValue[param.name]"
              class="!w-full"
              :placeholder="param.placeholder"
              :min="param.min"
              :max="param.max"
              :step="param.step"
            />
            <el-slider
              v-if="param.type === 'slider'"
              v-model="modelValue[param.name]"
              :min="param.min"
              :max="param.max"
              :step="param.step"
            />
            <el-date-picker
              v-if="param.type === 'date'"
              v-model="modelValue[param.name]"
              :placeholder="param.placeholder"
            />
            <el-time-picker
              v-if="param.type === 'time'"
              v-model="modelValue[param.name]"
              :placeholder="param.placeholder"
            />
            <el-select
              v-if="param.type === 'select'"
              v-model="modelValue[param.name]"
              :placeholder="param.placeholder"
              :popper-class="param.popperClass"
              filterable
            >
              <template #prefix v-if="param.prefix">
                <i class="iconfont !text-lg" :class="param.prefix"></i>
              </template>
              <el-option
                v-for="option in param.options"
                :key="option.value"
                :label="option.label"
                :value="option.value"
              >
                <div class="flex justify-start" v-if="option.image">
                  <span class="flex py-3 mr-2">
                    <img
                      :src="option.image"
                      class="rounded-lg"
                      :style="{ width: param.imgSize, height: param.imgSize }"
                  /></span>
                  <div class="flex !items-start flex-col py-2 space-y-1">
                    <span class="label text-sm">{{ option.label }}</span>
                    <span
                      class="text-xs text-gray-500 break-words line-clamp-1 max-w-[200px]"
                      :title="option.value"
                      >{{ option.value }}</span
                    >
                  </div>
                </div>
                <div class="flex justify-start items-center h-full" v-else>
                  <span class="label text-sm">{{ option.label }}</span>
                </div>
              </el-option>
            </el-select>
            <el-input
              type="textarea"
              v-if="param.type === 'textarea'"
              v-model="modelValue[param.name]"
              :autosize="param.autosize"
              :maxlength="param.maxlength"
              :show-word-limit="param.showWordLimit"
              :placeholder="param.placeholder"
            />
            <ImageUpload
              v-if="param.type === 'image'"
              v-model="modelValue[param.name]"
              :max-count="param.maxCount"
              :multiple="param.multiple"
              :max-size="param.maxSize"
              :accept="param.accept"
            />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import ImageUpload from './ImageUpload.vue'
import ParamEmpty from './ui/ParamEmpty.vue'

const title = ref('参数构建器')
const statusText = ref('功能正在开发中')
const description = ref('我们正在努力完善当前功能，敬请期待！')

const props = defineProps({
  modelValue: {
    type: Object,
    required: true,
  },
  requiredKeys: {
    type: Object,
    default: {},
    required: false,
  },
  items: {
    type: Array,
    required: true,
  },
  progress: {
    type: Number,
    default: 65,
    validator: (value) => value >= 0 && value <= 100,
  },
})

const selectedModel = ref(props.items[0])
const requiredKeys = ref(props.requiredKeys)

const emit = defineEmits(['update:modelValue', 'update:requiredKeys'])

// 初始化 modelValue 默认值
const initModelValue = (model) => {
  if (props.items.length === 0) {
    return {}
  }
  const defaultValues = {}
  requiredKeys.value = {}
  if (model && model.params) {
    model.params.forEach((param) => {
      if (param.required) {
        requiredKeys.value[param.name] = { required: true, label: param.label }
      }
      // 根据参数类型设置默认值
      switch (param.type) {
        case 'text':
        case 'textarea':
          defaultValues[param.name] = param.value || ''
          break
        case 'number':
          defaultValues[param.name] = param.value || 0
          break
        case 'slider':
          defaultValues[param.name] = param.value || param.min || 0
          break
        case 'select':
          // 如果有选项，选择第一个选项作为默认值
          defaultValues[param.name] =
            param.value || (param.options && param.options[0] ? param.options[0].value : '')
          break
        case 'checkbox':
        case 'switch':
          defaultValues[param.name] = param.value || false
          break
        case 'date':
        case 'time':
          defaultValues[param.name] = param.value || null
          break
        case 'image':
          defaultValues[param.name] = param.value || []
          break
        default:
          defaultValues[param.name] = param.value || ''
      }
    })
  }
  defaultValues.model = selectedModel.value.key
  return defaultValues
}

// 初始化默认值
const modelValue = ref(initModelValue(selectedModel.value))

// 监听 modelValue 变化，通知父组件
watch(
  modelValue,
  (newValue) => {
    emit('update:modelValue', newValue)
  },
  { deep: true }
)

watch(
  requiredKeys,
  (newValue) => {
    emit('update:requiredKeys', newValue)
  },
  { deep: true }
)

watch(
  () => props.items,
  (newValue) => {
    selectedModel.value = newValue[0]
    modelValue.value = initModelValue(selectedModel.value)
  },
  { deep: true }
)

// 组件挂载时初始化
onMounted(() => {
  // 确保初始值被正确设置
  if (props.modelValue && Object.keys(props.modelValue).length > 0) {
    modelValue.value = { ...props.modelValue }
  } else {
    modelValue.value = initModelValue(selectedModel.value)
  }
})

const changeModel = (item) => {
  if (item) {
    selectedModel.value = item
    // 更新 modelValue 为选中模型的默认值
    modelValue.value = initModelValue(item)
  }
}
</script>

<style lang="scss">
.param-builder {
  .model-version {
    background: url('@/assets/img/model-version.png') no-repeat center center;
    background-size: cover;
  }
  .el-select__wrapper {
    min-height: 34px;
    line-height: 25px;
  }
}
.model-select {
  .el-select-dropdown__item {
    height: auto !important;
  }
  .model-version {
    background: url('@/assets/img/model-version.png') no-repeat center center;
    background-size: cover;
  }
}
</style>
