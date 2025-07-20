<template>
  <div class="system-config form" v-loading="loading">
    <div class="container">
      <el-form
        :model="jimengConfig"
        label-width="150px"
        label-position="right"
        ref="configFormRef"
        :rules="rules"
      >
        <el-tabs type="border-card">
          <el-tab-pane>
            <template #label>
              <i class="iconfont icon-token mr-1"></i>
              <span>秘钥配置</span>
            </template>
            <el-form-item label="AccessKey" prop="access_key">
              <el-input
                v-model="jimengConfig.access_key"
                placeholder="请输入即梦AI的AccessKey"
                show-password
              />
            </el-form-item>
            <el-form-item label="SecretKey" prop="secret_key">
              <el-input
                v-model="jimengConfig.secret_key"
                placeholder="请输入即梦AI的SecretKey"
                show-password
              />
            </el-form-item>
          </el-tab-pane>

          <el-tab-pane>
            <template #label>
              <i class="iconfont icon-logout mr-1"></i>
              <span>算力配置</span>
            </template>
            <el-form-item>
              <template #label>
                <div class="label-title">
                  文生图算力
                  <el-tooltip
                    effect="dark"
                    content="用户使用文生图功能时消耗的算力"
                    raw-content
                    placement="right"
                  >
                    <el-icon>
                      <InfoFilled />
                    </el-icon>
                  </el-tooltip>
                </div>
              </template>
              <el-input-number
                v-model="jimengConfig.power.text_to_image"
                :min="1"
                :max="100"
                placeholder="请输入文生图算力消耗"
              />
            </el-form-item>

            <el-form-item>
              <template #label>
                <div class="label-title">
                  图生图算力
                  <el-tooltip
                    effect="dark"
                    content="用户使用图生图功能时消耗的算力"
                    raw-content
                    placement="right"
                  >
                    <el-icon>
                      <InfoFilled />
                    </el-icon>
                  </el-tooltip>
                </div>
              </template>
              <el-input-number
                v-model="jimengConfig.power.image_to_image"
                :min="1"
                :max="100"
                placeholder="请输入图生图算力消耗"
              />
            </el-form-item>

            <el-form-item>
              <template #label>
                <div class="label-title">
                  图片编辑算力
                  <el-tooltip
                    effect="dark"
                    content="用户使用图片编辑功能时消耗的算力"
                    raw-content
                    placement="right"
                  >
                    <el-icon>
                      <InfoFilled />
                    </el-icon>
                  </el-tooltip>
                </div>
              </template>
              <el-input-number
                v-model="jimengConfig.power.image_edit"
                :min="1"
                :max="100"
                placeholder="请输入图片编辑算力消耗"
              />
            </el-form-item>

            <el-form-item>
              <template #label>
                <div class="label-title">
                  图片特效算力
                  <el-tooltip
                    effect="dark"
                    content="用户使用图片特效功能时消耗的算力"
                    raw-content
                    placement="right"
                  >
                    <el-icon>
                      <InfoFilled />
                    </el-icon>
                  </el-tooltip>
                </div>
              </template>
              <el-input-number
                v-model="jimengConfig.power.image_effects"
                :min="1"
                :max="100"
                placeholder="请输入图片特效算力消耗"
              />
            </el-form-item>

            <el-form-item>
              <template #label>
                <div class="label-title">
                  文生视频算力
                  <el-tooltip
                    effect="dark"
                    content="用户使用文生视频功能时消耗的算力"
                    raw-content
                    placement="right"
                  >
                    <el-icon>
                      <InfoFilled />
                    </el-icon>
                  </el-tooltip>
                </div>
              </template>
              <el-input-number
                v-model="jimengConfig.power.text_to_video"
                :min="1"
                :max="100"
                placeholder="请输入文生视频算力消耗"
              />
            </el-form-item>

            <el-form-item>
              <template #label>
                <div class="label-title">
                  图生视频算力
                  <el-tooltip
                    effect="dark"
                    content="用户使用图生视频功能时消耗的算力"
                    raw-content
                    placement="right"
                  >
                    <el-icon>
                      <InfoFilled />
                    </el-icon>
                  </el-tooltip>
                </div>
              </template>
              <el-input-number
                v-model="jimengConfig.power.image_to_video"
                :min="1"
                :max="100"
                placeholder="请输入图生视频算力消耗"
              />
            </el-form-item>
          </el-tab-pane>
        </el-tabs>

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
import { httpGet, httpPost } from '@/utils/http'
import { InfoFilled } from '@element-plus/icons-vue'
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
    const res = await httpGet('/api/admin/jimeng/config')
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

    await httpPost('/api/admin/jimeng/config', {
      config: jimengConfig.value,
    })

    ElMessage.success('配置保存成功！')
  } catch (e) {
    if (e.message) {
      ElMessage.error('保存失败：' + e.message)
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

<style lang="stylus" scoped>
@import '../../../assets/css/admin/form.styl'
@import '../../../assets/css/main.styl'

.system-config {
  display flex
  justify-content center

  .container {
    width 100%
    max-width 800px
  }

  .label-title {
    display flex
    align-items center
    gap 5px
  }

  .el-input-number {
    width 100%
  }
}
</style>
