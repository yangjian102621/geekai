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
            <i
              class="iconfont !text-2xl mr-2"
              :class="item.icon.iconfont"
              v-if="item.icon.iconfont"
            ></i>
            <span
              class="flex items-center justify-center text-white model-version mr-2 w-[40px] h-[40px] rounded-lg"
              v-else
              :class="item.icon.size ? item.icon.size : '!text-xl'"
              >{{ item.icon.text }}</span
            >
            <div class="flex !items-start flex-col py-2 space-y-1">
              <span class="label text-sm">{{ item.name }}</span>
              <div class="whitespace-pre-line">
                <span
                  class="text-xs text-gray-500 break-words line-clamp-1 max-w-[250px]"
                  :title="item.label"
                  >{{ item.label }}</span
                >
              </div>
            </div>
          </div>
        </el-option>
      </el-select>

      <!-- 模型介绍 -->
      <div
        v-if="selectedModel && selectedModel.label"
        class="w-full p-4 bg-gray-50 rounded-lg border border-gray-200"
      >
        <div class="flex-1 min-w-0">
          <h4 class="text-sm font-semibold text-blue-600 mb-1">
            <i
              class="iconfont !text-base"
              :class="selectedModel.icon.iconfont"
              v-if="selectedModel.icon.iconfont"
            ></i>

            {{ selectedModel.name }}
          </h4>
          <p class="text-sm text-gray-600 whitespace-pre-line leading-relaxed">
            {{ selectedModel.label }}
          </p>
        </div>
      </div>

      <template v-for="param in selectedModel.params">
        <div class="w-full" :key="param.name" v-if="isParamVisible(param)">
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
            <div class="flex w-full flex-col">
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
              <ReferenceInput
                v-if="param.type === 'textarea' && param.name === 'prompt'"
                v-model="modelValue[param.name]"
                :autosize="param.autosize"
                :maxlength="param.maxlength"
                :show-word-limit="param.showWordLimit"
                :placeholder="param.placeholder"
                :image-list="
                  Array.isArray(modelValue[param.reference])
                    ? modelValue[param.reference]
                    : modelValue[param.reference]
                      ? [modelValue[param.reference]]
                      : []
                "
              />
              <el-input
                v-else-if="param.type === 'textarea'"
                type="textarea"
                v-model="modelValue[param.name]"
                :autosize="param.autosize"
                :maxlength="param.maxlength"
                :show-word-limit="param.showWordLimit"
                :placeholder="param.placeholder"
              />
              <div class="flex justify-end pt-2 w-full" v-if="param.enablePromptOptimizer">
                <el-button
                  @click="generatePrompt"
                  type="primary"
                  size="small"
                  :loading="submitting"
                >
                  <i class="iconfont icon-chuangzuo mr-1"></i>
                  优化提示词
                </el-button>
              </div>
              <ImageUpload
                v-if="param.type === 'image'"
                v-model="modelValue[param.name]"
                :max-count="param.maxCount"
                :multiple="param.multiple"
                :max-size="param.maxSize"
                :accept="param.accept"
              />
              <FileUpload
                v-if="param.type === 'file'"
                v-model="modelValue[param.name]"
                :max-count="param.maxCount"
                :multiple="param.multiple"
                :max-size="param.maxSize"
                :accept="param.accept"
              />
            </div>
          </div>
        </div>
      </template>
    </div>
  </div>
</template>

<script setup>
import FileUpload from './FileUpload.vue'
import ImageUpload from './ImageUpload.vue'
import ReferenceInput from './ReferenceInput.vue'
import ParamEmpty from './ui/ParamEmpty.vue'
import { httpPost } from '@/utils/http'
import { showMessageError } from '@/utils/dialog'
import { onMounted, ref, watch } from 'vue'

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
const submitting = ref(false)

const emit = defineEmits(['update:modelValue', 'update:requiredKeys', 'price-params-change'])

// 判断单个 showWhen 条件是否匹配
const matchSingleShowWhen = (condition, values) => {
  if (!condition || !values) return false
  const { field, value } = condition
  if (!field) return false
  return values[field] === value
}

// 统一的参数可见性判断，支持：
// - 无 showWhen：始终可见
// - showWhen 为对象：单条件
// - showWhen 为数组：所有条件满足（AND）
const isParamVisible = (param, values = modelValue.value) => {
  if (!param) return false
  if (param.type === 'hidden') return false

  const { showWhen } = param
  if (!showWhen) {
    return true
  }

  // 对象写法：单条件
  if (!Array.isArray(showWhen)) {
    return matchSingleShowWhen(showWhen, values)
  }

  // 数组写法：全部条件满足才显示（AND）
  return showWhen.every((condition) => matchSingleShowWhen(condition, values))
}

const recalcRequiredKeys = (values, model) => {
  const next = {}
  if (model && model.params) {
    model.params.forEach((param) => {
      if (param.required && isParamVisible(param, values)) {
        next[param.name] = { required: true, label: param.label }
      }
    })
  }
  requiredKeys.value = next
}

const generatePrompt = async () => {
  const prompt = modelValue.value?.prompt
  if (!prompt) {
    return showMessageError('请输入原始提示词')
  }
  try {
    submitting.value = true
    const res = await httpPost('/api/prompt/video', { prompt })
    modelValue.value = { ...modelValue.value, prompt: res.data }
  } catch (error) {
    showMessageError('生成提示词失败：' + error.message)
  } finally {
    submitting.value = false
  }
}

// 初始化 modelValue 默认值
const initModelValue = (model) => {
  if (props.items.length === 0) {
    return {}
  }
  const defaultValues = {}
  if (model && model.params) {
    model.params.forEach((param) => {
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
          defaultValues[param.name] =
            param.value || (param.multiple || (param.maxCount && param.maxCount > 1) ? [] : '')
          break
        case 'file':
          defaultValues[param.name] =
            param.value || (param.multiple || (param.maxCount && param.maxCount > 1) ? [] : '')
          break
        default:
          defaultValues[param.name] = param.value || ''
      }
    })
  }
  // 初始化 req_key 和 action
  defaultValues.req_key = selectedModel.value.key
  defaultValues.action = selectedModel.value.action
    ? selectedModel.value.action
    : 'CVSync2AsyncSubmitTask'
  recalcRequiredKeys(defaultValues, model)
  return defaultValues
}

// 初始化默认值
const modelValue = ref(initModelValue(selectedModel.value))

// 检查字段是否在 priceParams 中
const isPriceParam = (fieldName) => {
  if (!selectedModel.value || !selectedModel.value.priceParams) {
    return false
  }
  return selectedModel.value.priceParams.includes(fieldName)
}

// 监听 modelValue 变化，通知父组件
watch(
  modelValue,
  (newValue, oldValue) => {
    emit('update:modelValue', newValue)
    recalcRequiredKeys(newValue, selectedModel.value)

    // 检查是否有 priceParams 相关字段变化
    if (selectedModel.value && selectedModel.value.priceParams) {
      const priceParamsChanged = selectedModel.value.priceParams.some((paramName) => {
        return newValue[paramName] !== oldValue?.[paramName]
      })

      if (priceParamsChanged) {
        emit('price-params-change', newValue)
      }
    }
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
    // 模型切换时也触发价格参数变化事件
    if (item.priceParams && item.priceParams.length > 0) {
      emit('price-params-change', modelValue.value)
    }
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
