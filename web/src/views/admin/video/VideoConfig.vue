<template>
  <div class="system-config form" v-loading="loading">
    <div class="container">
      <el-form
        :model="videoConfig"
        label-width="150px"
        label-position="top"
        ref="configFormRef"
        :rules="rules"
        class="py-3 px-5"
      >
        <!-- API 配置分组 -->
        <div class="mb-3">
          <h3 class="heading-3 mb-2">API 配置</h3>
          <div class="py-3">
            <Alert type="info">
              <p class="mb-1">
                配置视频生成服务的 API 地址和 API Key。这些配置将用于所有视频生成模型的 API 调用。
              </p>
            </Alert>
          </div>
          <el-form-item label="API 地址" prop="api_url">
            <el-input
              v-model="videoConfig.api_url"
              placeholder="请输入视频生成服务的 API 地址，如：https://api.geekai.pro"
            />
          </el-form-item>
          <el-form-item label="API Key" prop="api_key">
            <el-input
              v-model="videoConfig.api_key"
              type="password"
              show-password
              placeholder="请输入视频生成服务的 API Key"
            />
          </el-form-item>
        </div>
        <el-divider />
        <!-- 算力配置分组 -->
        <div class="mb-3">
          <h3 class="heading-3 mb-3">模型算力配置</h3>
          <Alert type="info" class="mb-3">
            <div class="text-gray-500">
              根据模型的 priceParams
              配置自动生成价格组合。固定价格模型只需配置一个价格，多参数模型会生成所有参数组合的价格配置。
            </div>
          </Alert>

          <div v-for="provider in providers" :key="provider" class="mb-4">
            <h4 class="mb-2 text-base font-bold flex items-center gap-2">
              <span>{{ getProviderName(provider) }}</span>
            </h4>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
              <div
                v-for="model in getModelsByProvider(provider)"
                :key="model.key"
                class="p-3 rounded-md border border-gray-100"
              >
                <div class="text-sm mb-2">
                  <div class="font-bold">{{ model.name }}</div>
                  <div class="text-gray-500 line-clamp-2" :title="model.label">
                    {{ model.label }}
                  </div>
                </div>

                <!-- 模型基础信息 -->
                <el-form-item label="模型标识" class="mb-2">
                  <el-input
                    v-model="getModelConfig(model.key).model"
                    placeholder="模型标识，如：veo-2.0"
                    size="small"
                  />
                </el-form-item>

                <el-form-item label="服务提供商" class="mb-2">
                  <el-input
                    v-model="getModelConfig(model.key).provider"
                    placeholder="服务提供商，如：veo"
                    size="small"
                  />
                </el-form-item>

                <!-- 价格配置 -->
                <el-divider />
                <div class="text-xs text-gray-500 mb-2">
                  价格配置（基于 priceParams: {{ JSON.stringify(model.priceParams || []) }}）
                </div>
                <template v-for="priceKey in generatePriceKeys(model)" :key="priceKey">
                  <el-form-item :label="getPriceKeyLabel(priceKey, model)" class="mb-2">
                    <el-input-number
                      v-model="getModelConfig(model.key).power_config[priceKey]"
                      :min="1"
                      :placeholder="`请输入${getPriceKeyLabel(priceKey, model)}的算力消耗`"
                      size="small"
                      class="w-full"
                    />
                  </el-form-item>
                </template>
              </div>
            </div>
          </div>
        </div>
        <div style="padding: 10px">
          <el-form-item>
            <el-button type="primary" @click="saveConfig" :loading="saving">保存配置</el-button>
            <el-button @click="resetConfig">重置</el-button>
          </el-form-item>
        </div>
      </el-form>
    </div>
  </div>
</template>

<script setup>
import Alert from '@/components/ui/Alert.vue'
import { VideoParams, getVideoProviders } from '@/store/data/video_params'
import { httpGet, httpPost } from '@/utils/http'
import { ElMessage } from 'element-plus'
import { onMounted, ref } from 'vue'

const videoConfig = ref({
  api_url: '',
  api_key: '',
  video_powers: {},
})

const loading = ref(true)
const saving = ref(false)
const configFormRef = ref()
const providers = getVideoProviders()

// 表单验证规则
const rules = {
  api_url: [{ required: true, message: '请输入 API 地址', trigger: 'blur' }],
  api_key: [{ required: true, message: '请输入 API Key', trigger: 'blur' }],
}

// 获取服务提供商名称
const getProviderName = (provider) => {
  const names = {
    veo: 'Veo',
    sora: 'Sora',
    luma: 'Luma',
    keling: '可灵',
    minimax: 'MiniMax',
    wan: '通义万相',
  }
  return names[provider] || provider
}

// 获取指定 provider 的所有模型
const getModelsByProvider = (provider) => {
  return VideoParams[provider] || []
}

// 获取模型配置
const getModelConfig = (modelKey) => {
  if (!videoConfig.value.video_powers[modelKey]) {
    // 从 VideoParams 中获取默认配置
    const model = findModelByKey(modelKey)
    const priceKeys = generatePriceKeys(model)
    const defaultPowerConfig = {}
    priceKeys.forEach((key) => {
      defaultPowerConfig[key] = 10
    })

    videoConfig.value.video_powers[modelKey] = {
      provider: model?.provider || '',
      model: modelKey,
      power_config: defaultPowerConfig,
      api_key_type: '',
    }
  }
  // 确保 power_config 包含所有需要的 priceKeys
  const model = findModelByKey(modelKey)
  const priceKeys = generatePriceKeys(model)
  const currentConfig = videoConfig.value.video_powers[modelKey]
  if (!currentConfig.power_config) {
    currentConfig.power_config = {}
  }
  priceKeys.forEach((key) => {
    if (currentConfig.power_config[key] === undefined) {
      currentConfig.power_config[key] = 10
    }
  })
  return currentConfig
}

// 根据 key 查找模型
const findModelByKey = (key) => {
  for (const provider in VideoParams) {
    const models = VideoParams[provider]
    const model = models.find((m) => m.key === key)
    if (model) {
      return model
    }
  }
  return null
}

// 获取参数的所有可能值
const getParamValues = (model, paramName) => {
  const param = model.params?.find((p) => p.name === paramName)
  if (!param) return []

  if (param.type === 'select' && param.options) {
    return param.options.map((opt) => opt.value)
  }
  if (param.type === 'switch') {
    return [true, false]
  }
  return []
}

// 笛卡尔乘积生成函数
const cartesianProduct = (arrays) => {
  return arrays.reduce(
    (acc, arr) => {
      return acc.flatMap((x) => arr.map((y) => [...x, y]))
    },
    [[]]
  )
}

// 格式化参数值为字符串（用于生成 priceKey）
const formatParamValue = (paramName, value) => {
  if (paramName === 'sound') {
    // 处理 sound 参数：支持 'on'/'off' 字符串，也支持 true/false 布尔值
    return value === true || value === 'true' || value === 1 || value === 'on' ? 'sound' : 'silent'
  }
  if (typeof value === 'boolean') {
    return value ? 'true' : 'false'
  }
  return String(value)
}

// 根据模型的 priceParams 生成所有价格 key
const generatePriceKeys = (model) => {
  if (!model.priceParams || !Array.isArray(model.priceParams)) {
    return ['fixed']
  }

  // 如果是固定价格
  if (model.priceParams.length === 1 && model.priceParams[0] === 'fixed') {
    return ['fixed']
  }

  // 获取每个参数的所有可能值
  const paramValueArrays = model.priceParams.map((paramName) => {
    const values = getParamValues(model, paramName)
    if (values.length === 0) {
      // 如果参数没有预定义选项，返回参数名本身（这种情况应该很少）
      return [paramName]
    }
    return values
  })

  // 生成笛卡尔乘积
  const combinations = cartesianProduct(paramValueArrays)

  // 将每个组合转换为 priceKey（用下划线连接）
  return combinations.map((combination) => {
    return combination
      .map((value, index) => formatParamValue(model.priceParams[index], value))
      .join('_')
  })
}

// 生成价格 key 的显示标签
const getPriceKeyLabel = (priceKey, model) => {
  if (priceKey === 'fixed') {
    return '固定价格（算力）'
  }

  if (!model.priceParams || model.priceParams.length === 0) {
    return priceKey
  }

  // 解析 priceKey
  const parts = priceKey.split('_')
  const labels = []

  model.priceParams.forEach((paramName, index) => {
    const value = parts[index]
    const param = model.params?.find((p) => p.name === paramName)

    if (param) {
      let label = param.label || paramName

      // 查找对应的选项标签
      if (param.type === 'select' && param.options) {
        const option = param.options.find((opt) => {
          const optValue = formatParamValue(paramName, opt.value)
          return optValue === value
        })
        if (option) {
          label = `${label}: ${option.label}`
        } else {
          label = `${label}: ${value}`
        }
      } else if (param.type === 'switch') {
        label = `${label}: ${value === 'sound' ? '有声音' : value === 'silent' ? '无声音' : value}`
      } else {
        label = `${label}: ${value}`
      }

      labels.push(label)
    } else {
      labels.push(`${paramName}: ${value}`)
    }
  })

  return labels.join(' | ')
}

onMounted(() => {
  loadConfig()
})

// 加载配置
const loadConfig = async () => {
  try {
    const res = await httpGet('/api/admin/video/config')
    const cfg = res.data || {}
    cfg.video_powers = cfg.video_powers || {}
    videoConfig.value = cfg
  } catch (e) {
    ElMessage.error('加载配置失败: ' + e.message)
  } finally {
    loading.value = false
  }
}

// 保存配置
const saveConfig = async () => {
  try {
    await configFormRef.value.validate()
    saving.value = true
    await httpPost('/api/admin/video/config/update', videoConfig.value)
    ElMessage.success('配置保存成功！')
  } catch (e) {
    if (e.message) {
      ElMessage.error(e.message)
    } else {
      // 处理验证参数报错
      for (const key in e) {
        ElMessage.error(e[key][0]?.message)
        break
      }
    }
  } finally {
    saving.value = false
  }
}

// 重置配置
const resetConfig = () => {
  videoConfig.value = {
    api_url: '',
    api_key: '',
    video_powers: {},
  }
  ElMessage.info('配置已重置')
}
</script>

<style lang="scss" scoped>
@use '@/assets/css/admin/form.scss' as *;
@use '@/assets/css/main.scss' as *;

.system-config {
  display: flex;
  justify-content: center;

  .container {
    width: 100%;
    max-width: 1200px;
  }

  .heading-3 {
    color: var(--theme-text-color-primary);
  }

  .el-input-number {
    width: 100%;
  }
}
</style>
