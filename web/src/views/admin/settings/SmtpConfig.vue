<template>
  <div class="settings container p-5">
    <el-tabs v-model="activeTab" type="border-card">
      <el-tab-pane label="邮件服务设置" name="smtp">
        <Alert type="info">
          如果你不知道怎么获取这些配置信息，请参考文档：
          <a
            href="https://docs.geekai.me/plus/config/sms.html#%E9%82%AE%E4%BB%B6%E6%9C%8D%E5%8A%A1%E9%85%8D%E7%BD%AE"
            target="_blank"
            >邮件服务配置</a
          >。
        </Alert>

        <el-form :model="smtpConfig" label-position="top">
          <el-form-item>
            <label class="form-label"
              >邮件服务器地址
              <el-tooltip placement="top">
                <template #content>
                  推荐使用网易邮箱，<br />
                  国外邮箱推荐使用 Google 邮箱
                </template>
                <i class="iconfont icon-info"></i>
              </el-tooltip>
            </label>
            <el-input
              v-model="smtpConfig.host"
              placeholder="请输入邮件服务器地址，推荐使用网易邮箱"
            />
          </el-form-item>

          <el-form-item>
            <label class="form-label"
              >邮件服务器端口
              <el-tooltip placement="top">
                <template #content> 线上推荐使用465端口，<br />本地测试推荐使用25端口 </template>
                <i class="iconfont icon-info"></i>
              </el-tooltip>
            </label>
            <el-input v-model="smtpConfig.port" type="number" placeholder="请输入端口号" />
          </el-form-item>

          <el-form-item label="是否使用TLS"
            ><el-switch v-model="smtpConfig.use_tls"
          /></el-form-item>

          <el-form-item>
            <label class="form-label"
              >应用名称
              <el-tooltip placement="top">
                <template #content> 应用名称会显示在邮件的抬头 </template>
                <i class="iconfont icon-info"></i>
              </el-tooltip>
            </label>
            <el-input v-model="smtpConfig.app_name" placeholder="请输入应用名称" />
          </el-form-item>

          <el-form-item>
            <label class="form-label">发件人邮箱地址</label>
            <el-input v-model="smtpConfig.from" type="email" placeholder="请输入发件人邮箱地址" />
          </el-form-item>

          <el-form-item>
            <label class="form-label"
              >发件人邮箱密码
              <el-tooltip placement="top">
                <template #content> 如果使用授权码，请输入授权码 </template>
                <i class="iconfont icon-info"></i>
              </el-tooltip>
            </label>
            <el-input
              v-model="smtpConfig.password"
              type="password"
              placeholder="请输入邮箱密码或授权码"
              show-password
            />
          </el-form-item>
        </el-form>
      </el-tab-pane>
    </el-tabs>

    <div class="flex justify-center mt-6">
      <el-button type="primary" @click="saveConfig" :loading="loading">提交保存</el-button>
    </div>
  </div>
</template>

<script setup>
import Alert from '@/components/ui/Alert.vue'
import { httpGet, httpPost } from '@/utils/http'
import { ElMessage } from 'element-plus'
import { onMounted, ref } from 'vue'

const loading = ref(false)
const activeTab = ref('smtp')
const smtpConfig = ref({
  use_tls: false,
  host: 'smtp.163.com',
  port: 465,
  app_name: 'GeekAI',
  from: '',
  password: '',
})

onMounted(async () => {
  try {
    const res = await httpGet('/api/admin/config/get?key=smtp')
    smtpConfig.value = Object.assign(smtpConfig.value, res.data || {})
  } catch (e) {
    // 使用默认值
  }
})

// 保存配置
const saveConfig = async () => {
  loading.value = true
  try {
    await httpPost('/api/admin/config/update/smtp', smtpConfig.value)
    ElMessage.success('保存成功')
  } catch (e) {
    ElMessage.error('保存失败：' + (e.message || '未知错误'))
  } finally {
    loading.value = false
  }
}
</script>

<style lang="scss">
.settings {
  a {
    color: #409eff;
    &:hover {
      text-decoration: underline;
    }
  }
  .el-form-item__label {
    font-weight: 700;
  }
}
</style>
