<template>
  <div class="ppt-config form" v-loading="loading">
    <div class="container">
      <el-form :model="config" label-position="top" ref="formRef" class="px-3 py-5" :rules="rules">
        <!-- 分镜 LLM 配置 -->
        <div class="mb-4">
          <h3 class="heading-3 mb-2">分镜 LLM 配置</h3>
          <el-form-item label="分镜 LLM API 地址" prop="outline_llm_api_url">
            <el-input
              v-model="config.outline_llm_api_url"
              placeholder="例如：https://api.xxx.com/v1/chat/completions"
            />
          </el-form-item>
          <el-form-item label="分镜 LLM 模型名称" prop="outline_llm_model">
            <el-input
              v-model="config.outline_llm_model"
              placeholder="例如：gemini-3.1-pro-preview"
            />
          </el-form-item>
          <el-form-item label="分镜 LLM API Key" prop="outline_llm_api_key">
            <el-input
              v-model="config.outline_llm_api_key"
              type="password"
              show-password
              placeholder="用于生成 PPT 分镜的 LLM API Key"
            />
          </el-form-item>
        </div>

        <el-divider />

        <!-- 通用生成与算力配置 -->
        <div class="mb-4">
          <h3 class="heading-3 mb-2">生成策略与算力配置</h3>
          <el-form-item label="单个任务最大 PPT 页数" prop="max_slides_per_task">
            <el-input-number
              v-model="config.max_slides_per_task"
              :min="1"
              :max="50"
              class="w-full"
            />
          </el-form-item>
          <el-form-item label="每张 PPT 图片消耗算力" prop="power_cost_per_slide">
            <el-input-number v-model="config.power_cost_per_slide" :min="0" class="w-full" />
          </el-form-item>
          <el-form-item label="图片生成最大并发数" prop="max_concurrent_requests">
            <el-input-number
              v-model="config.max_concurrent_requests"
              :min="1"
              :max="10"
              class="w-full"
            />
          </el-form-item>
          <el-form-item label="外部图片 API QPS 限制" prop="qps_limit">
            <el-input-number v-model="config.qps_limit" :min="1" :max="10" class="w-full" />
          </el-form-item>
        </div>

        <el-divider />

        <!-- 图片模型配置 -->
        <div class="mb-4">
          <h3 class="heading-3 mb-2">图片生成模型配置</h3>
          <el-form-item label="当前使用的图片模型" prop="active_image_provider">
            <el-select v-model="config.active_image_provider" placeholder="请选择图片生成模型">
              <el-option label="Nano Banana" value="nano_banana" />
              <el-option label="Doubao Seedream" value="seedream" />
            </el-select>
          </el-form-item>

          <!-- Nano Banana 配置 -->
          <el-card
            v-if="config.active_image_provider === 'nano_banana'"
            class="mb-3"
            shadow="never"
          >
            <template #header>
              <span class="font-bold">Nano Banana 配置</span>
            </template>
            <el-form-item label="Nano Banana API 地址">
              <el-input
                v-model="config.nano_banana_api_url"
                placeholder="例如：https://xxx/v1/images/generations"
              />
            </el-form-item>
            <el-form-item label="Nano Banana 模型">
              <el-input
                v-model="config.nano_banana_model"
                placeholder="例如：gemini-3.1-flash-image-preview"
              />
            </el-form-item>
            <el-form-item label="响应格式">
              <el-select
                v-model="config.nano_banana_response_format"
                placeholder="默认 url"
                class="w-full"
                clearable
              >
                <el-option label="url" value="url" />
                <el-option label="b64_json" value="b64_json" />
              </el-select>
            </el-form-item>
            <el-form-item label="宽高比">
              <el-select
                v-model="config.nano_banana_aspect_ratio"
                placeholder="请选择宽高比"
                class="w-full"
              >
                <el-option label="16:9" value="16:9" />
                <el-option label="21:9" value="21:9" />
              </el-select>
            </el-form-item>
            <el-form-item label="Nano Banana API Key">
              <el-input
                v-model="config.nano_banana_api_key"
                type="password"
                show-password
                placeholder="Nano Banana API Key"
              />
            </el-form-item>
          </el-card>

          <!-- Seedream 配置 -->
          <el-card v-if="config.active_image_provider === 'seedream'" class="mb-3" shadow="never">
            <template #header>
              <span class="font-bold">Doubao Seedream 配置</span>
            </template>
            <el-form-item label="Seedream Base URL">
              <el-input
                v-model="config.seedream_base_url"
                placeholder="例如：https://ark.cn-beijing.volces.com/api/v3（API 根地址，SDK 会拼路径）"
              />
            </el-form-item>
            <el-form-item label="Seedream API Key (ARK_API_KEY)">
              <el-input
                v-model="config.seedream_api_key"
                type="password"
                show-password
                placeholder="从火山控制台获取的 ARK_API_KEY"
              />
            </el-form-item>
            <el-form-item label="模型 ID">
              <el-input
                v-model="config.seedream_model"
                placeholder="例如：doubao-seedream-5-0-260128"
              />
            </el-form-item>
            <el-form-item label="图片尺寸">
              <el-input v-model="config.seedream_size" placeholder="例如：2K" />
            </el-form-item>
            <el-form-item label="输出格式">
              <el-input v-model="config.seedream_output_format" placeholder="例如：png" />
            </el-form-item>
            <el-form-item label="响应格式">
              <el-input v-model="config.seedream_response_format" placeholder="例如：url" />
            </el-form-item>
            <el-form-item label="是否开启水印">
              <el-switch v-model="config.seedream_watermark" />
            </el-form-item>
          </el-card>
        </div>

        <div style="padding: 10px">
          <el-form-item>
            <el-button type="primary" @click="save" :loading="saving">保存配置</el-button>
            <el-button @click="loadConfig">重置</el-button>
          </el-form-item>
        </div>
      </el-form>
    </div>
  </div>
</template>

<script setup>
import { httpGet, httpPost } from '@/utils/http'
import { ElMessage } from 'element-plus'
import { onMounted, reactive, ref } from 'vue'

const config = ref({
  outline_llm_api_url: '',
  outline_llm_api_key: '',
  outline_llm_model: 'gemini-3.1-pro-preview',
  active_image_provider: '',
  max_slides_per_task: 10,
  power_cost_per_slide: 0,
  max_concurrent_requests: 3,
  qps_limit: 1,
  nano_banana_api_url: '',
  nano_banana_api_key: '',
  nano_banana_model: 'gemini-3.1-flash-image-preview',
  nano_banana_response_format: 'url',
  nano_banana_aspect_ratio: '16:9',
  seedream_base_url: '',
  seedream_api_key: '',
  seedream_model: '',
  seedream_size: '2K',
  seedream_output_format: 'png',
  seedream_response_format: 'url',
  seedream_watermark: false,
})

const loading = ref(false)
const saving = ref(false)
const formRef = ref(null)

const rules = reactive({
  max_slides_per_task: [{ required: true, message: '请输入最大 PPT 页数', trigger: 'blur' }],
  power_cost_per_slide: [{ type: 'number', message: '请输入合法的算力值', trigger: 'blur' }],
})

const applyDefaultConfig = () => {
  if (!config.value.outline_llm_model) {
    config.value.outline_llm_model = 'gemini-3.1-pro-preview'
  }
  if (!config.value.nano_banana_model) {
    config.value.nano_banana_model = 'gemini-3.1-flash-image-preview'
  }
  if (!config.value.nano_banana_response_format) {
    config.value.nano_banana_response_format = 'url'
  }
  if (!config.value.nano_banana_aspect_ratio) {
    config.value.nano_banana_aspect_ratio = '16:9'
  }
}

const loadConfig = () => {
  loading.value = true
  httpGet('/api/admin/ppt/config')
    .then((res) => {
      if (res.data) {
        config.value = Object.assign(config.value, res.data)
      }
      applyDefaultConfig()
    })
    .catch((e) => {
      ElMessage.error('加载 PPT 配置失败: ' + e.message)
    })
    .finally(() => {
      loading.value = false
    })
}

const save = () => {
  formRef.value.validate((valid) => {
    if (!valid) {
      return
    }
    saving.value = true
    httpPost('/api/admin/ppt/config/update', config.value)
      .then(() => {
        ElMessage.success('保存成功！')
        loadConfig()
      })
      .catch((e) => {
        ElMessage.error('保存失败：' + e.message)
      })
      .finally(() => {
        saving.value = false
      })
  })
}

onMounted(() => {
  applyDefaultConfig()
  loadConfig()
})
</script>

<style lang="scss" scoped>
@use '@/assets/css/admin/form.scss' as *;
@use '@/assets/css/main.scss' as *;

.ppt-config {
  display: flex;
  justify-content: center;
  padding: 20px;
}
</style>
