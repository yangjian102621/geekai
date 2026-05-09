<template>
  <div class="param-builder-mobile">
    <ParamEmpty
      v-if="items.length === 0"
      :progress="progress"
      :title="title"
      :status-text="statusText"
      :description="description"
    />

    <div v-else class="flex flex-col w-full space-y-5">
      <!-- 模型选择（移动端样式，使用 CustomSelect） -->
      <div class="bg-white rounded-xl shadow-sm w-full p-4">
        <label class="block text-gray-700 mb-2 font-semibold">选择模型</label>
        <CustomSelect v-model="selectedModelKey" :options="modelOptions" title="选择模型">
          <template #label>
            <div class="flex items-center w-full">
              <i class="iconfont icon-model !text-xl mr-2"></i>
              <span class="text-gray-700 font-semibold">{{
                selectedModel.name || '请选择模型'
              }}</span>
            </div>
          </template>
          <template #option="{ option, selected }">
            <div class="flex items-center w-full">
              <span
                class="flex items-center justify-center text-white model-version mr-2 w-[40px] h-[40px] rounded-lg"
                :class="option.iconSize ? option.iconSize : '!text-xl'"
                >{{ option.iconText }}</span
              >
              <div class="flex !items-start flex-col py-1">
                <span
                  class="font-semibold text-gray-900"
                  :class="{ '!text-purple-600': selected }"
                  >{{ option.label }}</span
                >
                <span
                  class="text-xs text-gray-500 line-clamp-1 max-w-[250px]"
                  :title="option.subLabel"
                  >{{ option.subLabel }}</span
                >
              </div>
            </div>
          </template>
        </CustomSelect>
      </div>

      <!-- 参数渲染（移动端卡片样式） -->
      <template v-for="param in selectedModel.params">
        <div
          class="bg-white rounded-xl shadow-sm w-full p-4"
          :key="param.name"
          v-if="param.type !== 'hidden'"
        >
          <!-- switch 类型单独处理 -->
          <div class="w-full flex flex-col !items-start space-y-2" v-if="param.type === 'switch'">
            <div class="w-full flex justify-between items-center">
              <label class="text-gray-700 font-semibold">{{ param.label }}</label>
              <van-switch v-model="modelValue[param.name]" size="default" />
            </div>
            <p v-if="param.info" class="text-xs text-gray-500">{{ param.info }}</p>
          </div>

          <div class="w-full flex flex-col !items-start space-y-2" v-else>
            <label class="text-gray-700 font-semibold">
              {{ param.label }}
              <span v-if="param.required" class="text-red-500 ml-1">*</span>
            </label>
            <p v-if="param.info" class="text-xs text-gray-500">{{ param.info }}</p>
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

              <!-- 使用 CustomSelect 替换 el-select -->
              <CustomSelect
                v-if="param.type === 'select'"
                v-model="modelValue[param.name]"
                :options="formatParamOptions(param)"
                :title="param.placeholder || '请选择' + param.label"
                class="w-full"
              >
                <template #label>
                  <div class="flex items-center w-full">
                    <i class="iconfont !text-xl mr-2" :class="param.prefix ? param.prefix : ''"></i>
                    <span class="text-gray-700 font-semibold">{{
                      param.placeholder || '请选择' + param.label
                    }}</span>
                  </div>
                </template>
                <template v-if="hasImageOption(param)" #option="{ option, selected }">
                  <div class="flex items-center w-full">
                    <el-image
                      v-if="option.image"
                      :src="option.image"
                      fit="cover"
                      class="w-10 h-10 rounded-lg mr-2"
                    />
                    <div class="flex !items-start flex-col py-1">
                      <span
                        class="font-bold text-gray-900 mr-2"
                        :class="{ '!text-purple-600': selected }"
                        >{{ option.label }}</span
                      >
                      <span
                        class="text-xs text-gray-500 line-clamp-1 max-w-[200px]"
                        :title="option.value"
                        >{{ option.value }}</span
                      >
                    </div>
                  </div>
                </template>
                <template #option="{ option, selected }" v-else>
                  <div class="flex items-center w-full">
                    <span class="mr-2" :class="{ 'font-bold !text-purple-600': selected }">{{
                      option.label
                    }}</span>
                  </div>
                </template>
              </CustomSelect>

              <el-input
                type="textarea"
                v-if="param.type === 'textarea'"
                v-model="modelValue[param.name]"
                :autosize="param.autosize || { minRows: 3, maxRows: 6 }"
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
import FileUpload from '@/components/FileUpload.vue'
import ImageUpload from '@/components/ImageUpload.vue'
import CustomSelect from '@/components/mobile/CustomSelect.vue'
import ParamEmpty from '@/components/ui/ParamEmpty.vue'
import { computed, onMounted, ref, watch } from 'vue'

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
const selectedModelKey = ref(props.items[0] ? props.items[0].key : '')
const requiredKeys = ref(props.requiredKeys)

const emit = defineEmits(['update:modelValue', 'update:requiredKeys'])

const initModelValue = (model) => {
  if (!props.items || props.items.length === 0) {
    return {}
  }
  const defaultValues = {}
  requiredKeys.value = {}
  if (model && model.params) {
    model.params.forEach((param) => {
      if (param.required) {
        requiredKeys.value[param.name] = { required: true, label: param.label }
      }
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
  defaultValues.req_key = selectedModel.value.key
  defaultValues.action = selectedModel.value.action
    ? selectedModel.value.action
    : 'CVSync2AsyncSubmitTask'
  return defaultValues
}

const modelValue = ref(initModelValue(selectedModel.value))

const modelOptions = computed(() => {
  return (props.items || []).map((m) => ({
    label: m.name,
    value: m.key,
    subLabel: m.label,
    iconText: m.icon?.text || '',
    iconSize: m.icon?.size || '!text-xl',
  }))
})

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
    selectedModelKey.value = newValue[0]?.key || ''
    modelValue.value = initModelValue(selectedModel.value)
  },
  { deep: true }
)

watch(selectedModelKey, (key) => {
  if (!key) return
  const found = (props.items || []).find((m) => m.key === key)
  if (found) {
    selectedModel.value = found
    modelValue.value = initModelValue(found)
  }
})

onMounted(() => {
  if (props.modelValue && Object.keys(props.modelValue).length > 0) {
    modelValue.value = { ...props.modelValue }
  } else {
    modelValue.value = initModelValue(selectedModel.value)
  }
})

const hasImageOption = (param) => {
  return Array.isArray(param.options) && param.options.some((o) => !!o.image)
}

const formatParamOptions = (param) => {
  return (param.options || []).map((o) => ({
    label: o.label,
    value: o.value,
    image: o.image,
  }))
}
</script>

<style lang="scss" scoped>
@use '@/assets/css/mobile/jimeng.scss';

.param-builder-mobile {
  .model-version {
    background: url('@/assets/img/model-version.png') no-repeat center center;
    background-size: cover;
  }
}

/* 采用 JimengCreate.vue 的卡片/表单视觉（该文件已引入 mobile/jimeng.scss） */
:deep(.custom-select) {
  .select-trigger {
    background-color: rgb(31, 41, 55);
    border-color: rgb(75, 85, 99);
    color: rgb(209, 213, 219);
  }
  .select-dropdown {
    background-color: rgb(55, 65, 81);
    border-color: rgb(75, 85, 99);
    box-shadow: 0 0 15px rgba(107, 80, 225, 0.8);
  }
  .select-option {
    color: rgb(209, 213, 219);
    &:hover {
      background-color: rgb(75, 85, 99);
    }
    &.selected {
      background-color: rgb(139, 92, 246);
      color: rgb(255, 255, 255);
    }
  }
}
</style>
