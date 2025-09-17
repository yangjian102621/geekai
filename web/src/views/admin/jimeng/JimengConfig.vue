<template>
  <div class="system-config form" v-loading="loading">
    <div class="container">
      <el-form
        :model="jimengConfig"
        label-width="150px"
        label-position="top"
        ref="configFormRef"
        :rules="rules"
        class="py-3 px-5"
      >
        <!-- 秘钥配置分组 -->
        <div class="mb-3">
          <h3 class="heading-3 mb-2">秘钥配置</h3>
          <div class="py-3">
            <Alert type="info">
              <p class="mb-1">
                1. 要使用即梦 AI 功能，需要先在火山引擎控制台开通
                <a
                  href="https://console.volcengine.com/ai/ability/detail/10"
                  target="_blank"
                  class="text-blue-500"
                  >即梦 AI</a
                >
                和
                <a
                  href="https://console.volcengine.com/ai/ability/detail/1"
                  target="_blank"
                  class="text-blue-500"
                  >智能绘图</a
                >以及<a
                  href="https://console.volcengine.com/ark/region:ark+cn-beijing/openManagement"
                  target="_blank"
                  class="text-blue-500"
                  >火山方舟</a
                >
                服务。
              </p>
              <p>
                2. AccessKey和SecretKey 请在火山引擎控制台 ->
                <a
                  href="https://console.volcengine.com/iam/keymanage/"
                  target="_blank"
                  class="text-blue-500"
                  >秘钥管理</a
                >
                获取。
              </p>
              <p>
                3. ApiKey 请在火山方舟控制台 ->
                <a
                  href="https://console.volcengine.com/ark/region:ark+cn-beijing/apiKey?apikey=%7B%7D"
                  target="_blank"
                  class="text-blue-500"
                >
                  API Key管理</a
                >
                获取。
              </p>
            </Alert>
          </div>
          <el-form-item label="AccessKey" prop="access_key">
            <el-input v-model="jimengConfig.access_key" placeholder="请输入即梦AI的AccessKey" />
          </el-form-item>
          <el-form-item label="SecretKey" prop="secret_key">
            <el-input v-model="jimengConfig.secret_key" placeholder="请输入即梦AI的SecretKey" />
          </el-form-item>
          <el-form-item prop="api_key">
            <template #label>
              <div class="text-sm">
                火山方舟服务API Key（<span class="text-red-400"
                  >目前火山方舟服务只支持API Key验证</span
                >）
              </div>
            </template>
            <el-input v-model="jimengConfig.api_key" placeholder="请输入火山方舟服务API Key" />
            <div class="text-sm mt-2 text-gray-500">
              目前豆包生图 4.0 模型在即梦API中不支持，需要使用火山方舟服务。
            </div>
          </el-form-item>
        </div>
        <el-divider />
        <!-- 算力配置分组 -->
        <div class="mb-3">
          <h3 class="heading-3 mb-3">任务积分配置</h3>
          <Alert type="info" class="mb-3">
            <div class="text-gray-500">
              图片类模型统一都是 0.2 元一张，假如你100积分售价1元，建议设置：20积分/张。
            </div>
            <div class="text-gray-500">
              视频/数字人/动作迁移单位：积分/秒，但是不同的模型的价格不一样，建议去火山方舟控制台查看，根据价格设置积分。
            </div>
          </Alert>

          <div v-for="func in functions" :key="func.key" class="mb-4">
            <h4 class="mb-2 text-base font-bold flex items-center gap-2">
              <i class="iconfont" :class="func.icon"></i>
              {{ func.name }}
              <el-tag size="small" type="info">{{ getUnit(func.key) }}</el-tag>
            </h4>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
              <div
                v-for="model in params[func.key]"
                :key="model.key"
                class="p-3 rounded-md border border-gray-100"
              >
                <div class="text-sm mb-2">
                  <div class="font-bold">{{ model.name }}</div>
                  <div class="text-gray-500 line-clamp-2" :title="model.label">
                    {{ model.label }}
                  </div>
                </div>
                <el-input-number
                  v-model="jimengConfig.powers[model.key]"
                  :min="1"
                  :placeholder="`对应模型：${model.key}（${getUnit(func.key)}）`"
                  class="w-full"
                />
                <div class="text-xs text-gray-400 mt-1">对应模型：{{ model.key }}</div>
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
import { JimengFunctions, JimengParams } from '@/store/data/jimeng_params'
import { httpGet, httpPost } from '@/utils/http'
import { ElMessage } from 'element-plus'
import { onMounted, ref } from 'vue'

const jimengConfig = ref({
  access_key: '',
  secret_key: '',
  api_key: '',
  powers: {},
})

const loading = ref(true)
const saving = ref(false)
const configFormRef = ref()
const functions = JimengFunctions
const params = JimengParams

// 表单验证规则
const rules = {
  access_key: [{ required: true, message: '请输入AccessKey', trigger: 'blur' }],
  secret_key: [{ required: true, message: '请输入SecretKey', trigger: 'blur' }],
}

const getUnit = (funcKey) => (funcKey === 'image' ? '积分/张' : '积分/秒')

onMounted(() => {
  loadConfig()
})

// 加载配置
const loadConfig = async () => {
  try {
    const res = await httpGet('/api/admin/config/get?key=jimeng')
    const cfg = res.data || {}
    cfg.powers = cfg.powers || {}
    jimengConfig.value = cfg
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
    await httpPost('/api/admin/jimeng/config/update', jimengConfig.value)
    ElMessage.success('配置保存成功！')
  } catch (e) {
    if (e.message) {
      ElMessage.error(e.message)
    }
  } finally {
    saving.value = false
  }
}

// 重置配置
const resetConfig = () => {
  jimengConfig.value = {
    access_key: '',
    secret_key: '',
    api_key: '',
    powers: {},
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
    max-width: 1000px;
  }

  .heading-3 {
    color: var(--theme-text-color-primary);
  }

  .label-title {
    display: flex;
    align-items: center;
    gap: 5px;
  }

  .el-input-number {
    width: 100%;
  }
}
</style>
