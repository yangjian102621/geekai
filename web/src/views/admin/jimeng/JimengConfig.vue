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
                  href="https://console.volcengine.com/ai/ability/detail/9"
                  target="_blank"
                  class="text-blue-500"
                  >智能绘图</a
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
            </Alert>
          </div>
          <el-form-item label="AccessKey" prop="access_key">
            <el-input v-model="jimengConfig.access_key" placeholder="请输入即梦AI的AccessKey" />
          </el-form-item>
          <el-form-item label="SecretKey" prop="secret_key">
            <el-input v-model="jimengConfig.secret_key" placeholder="请输入即梦AI的SecretKey" />
          </el-form-item>
        </div>
        <el-divider />
        <!-- 算力配置分组 -->
        <div class="mb-3">
          <h3 class="heading-3 mb-3">算力配置</h3>
          <el-form-item>
            <template #label>
              <div class="text-gray-500 text-sm">
                生成图片消耗的积分，包括：文生图、图生图、图片编辑、图片特效，<span
                  class="text-red-500"
                  >单位：积分/张</span
                >
              </div>
            </template>
            <el-input-number
              v-model="jimengConfig.power.image"
              :min="1"
              placeholder="请输入图片生成算力消耗"
            />
          </el-form-item>
          <el-form-item>
            <template #label>
              <div class="text-gray-500 text-sm">
                生成视频消耗的积分，包括：文生视频、图生视频，<span class="text-red-500"
                  >单位：积分/秒</span
                >
              </div>
            </template>
            <el-input-number
              v-model="jimengConfig.power.video"
              :min="1"
              placeholder="请输入视频生成算力消耗"
            />
          </el-form-item>
          <el-form-item>
            <template #label>
              <div class="text-gray-500 text-sm">
                生成数字人视频消耗的积分，<span class="text-red-500">单位：积分/秒</span>
              </div>
            </template>
            <el-input-number
              v-model="jimengConfig.power.virtual_human"
              :min="1"
              placeholder="请输入数字人视频生成算力消耗"
            />
          </el-form-item>
          <el-form-item>
            <template #label>
              <div class="text-gray-500 text-sm">
                生成视频动作迁移消耗的积分，<span class="text-red-500">单位：积分/秒</span>
              </div>
            </template>
            <el-input-number
              v-model="jimengConfig.power.action_transfer"
              :min="1"
              placeholder="请输入视频动作迁移算力消耗"
            />
          </el-form-item>
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
import { httpGet, httpPost } from '@/utils/http'
import { ElMessage } from 'element-plus'
import { onMounted, ref } from 'vue'

const jimengConfig = ref({
  access_key: '',
  secret_key: '',
  power: {
    text_to_image: 10,
    image_to_image: 15,
    image_edit: 20,
    image_effects: 25,
    text_to_video: 30,
    image_to_video: 35,
  },
})

const loading = ref(true)
const saving = ref(false)
const testing = ref(false)
const configFormRef = ref()

// 表单验证规则
const rules = {
  access_key: [{ required: true, message: '请输入AccessKey', trigger: 'blur' }],
  secret_key: [{ required: true, message: '请输入SecretKey', trigger: 'blur' }],
}

onMounted(() => {
  loadConfig()
})

// 加载配置
const loadConfig = async () => {
  try {
    const res = await httpGet('/api/admin/config/get?key=jimeng')
    jimengConfig.value = res.data
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
    power: {
      text_to_image: 10,
      image_to_image: 15,
      image_edit: 20,
      image_effects: 25,
      text_to_video: 30,
      image_to_video: 35,
    },
  }
  ElMessage.info('配置已重置')
}
</script>

<style lang="scss" scoped>
@use '../../../assets/css/admin/form.scss' as *;
@use '../../../assets/css/main.scss' as *;

.system-config {
  display: flex;
  justify-content: center;

  .container {
    width: 100%;
    max-width: 800px;
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
