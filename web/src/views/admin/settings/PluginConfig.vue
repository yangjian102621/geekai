<template>
  <div class="settings container p-5">
    <el-tabs v-model="activeTab" type="border-card">
      <el-tab-pane name="captcha">
        <template #label>
          <div class="flex items-center">
            <i class="iconfont icon-yanzm"></i>
            <span class="ml-2">行为验证码</span>
          </div>
        </template>

        <Alert type="info">
          行为验证码服务，开启后用户登录的时候需要进行行为验证，可以有效防止恶意登录。<br />
          请联系作者获取令牌填入下面输入框开通验证服务。
        </Alert>

        <el-form :model="captchaConfig" label-position="top">
          <el-form-item label="服务令牌">
            <el-input v-model="captchaConfig.api_key" placeholder="请输入服务令牌" />
          </el-form-item>

          <el-form-item label="验证码类型">
            <el-radio-group v-model="captchaConfig.type" size="large">
              <el-radio value="dot" border>点选验证码</el-radio>
              <el-radio value="slide" border>滑动验证码</el-radio>
            </el-radio-group>
          </el-form-item>

          <el-form-item label="启用验证码">
            <el-switch size="large" v-model="captchaConfig.enabled" />
          </el-form-item>
        </el-form>

        <div class="flex justify-center mt-6">
          <el-button type="primary" @click="saveCaptchaConfig" :loading="loading"
            >提交保存</el-button
          >
        </div>
      </el-tab-pane>

      <el-tab-pane name="wxlogin">
        <template #label>
          <div class="flex items-center">
            <i class="iconfont icon-wechat"></i>
            <span class="ml-2">微信登录</span>
          </div>
        </template>

        <Alert type="info">
          微信登录服务，开启后用户可以使用微信扫码登录。<br />
          请联系作者获取令牌填入下面输入框开通验证服务。
        </Alert>

        <el-form :model="wxLoginConfig" label-position="top">
          <el-form-item label="服务令牌">
            <el-input v-model="wxLoginConfig.api_key" placeholder="请输入服务令牌" />
          </el-form-item>

          <el-form-item label="登录成功回调URL">
            <el-input v-model="wxLoginConfig.notify_url" placeholder="请输入回调URL" />
          </el-form-item>

          <el-form-item label="启用微信登录">
            <el-switch size="large" v-model="wxLoginConfig.enabled" />
          </el-form-item>
        </el-form>

        <div class="flex justify-center mt-6">
          <el-button type="primary" @click="saveWxloginConfig" :loading="loading"
            >提交保存</el-button
          >
        </div>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script setup>
import Alert from '@/components/ui/Alert.vue'
import { httpGet, httpPost } from '@/utils/http.js'
import { ElMessage } from 'element-plus'
import { onMounted, ref } from 'vue'

const activeTab = ref('captcha')
const loading = ref(false)

// 行为验证码配置
const captchaConfig = ref({
  api_key: '',
  type: 'slide',
  enabled: false,
})

// 微信登录配置
const wxLoginConfig = ref({
  api_key: '',
  notify_url: '',
  enabled: false,
})

onMounted(() => {
  fetchCaptchaConfig()
  fetchWxloginConfig()
})

const fetchCaptchaConfig = async () => {
  try {
    const res = await httpGet('/api/admin/config/get?key=captcha')
    captchaConfig.value = Object.assign(captchaConfig.value, res.data)
  } catch (e) {
    // 使用默认值
  }
}

const fetchWxloginConfig = async () => {
  try {
    const res = await httpGet('/api/admin/config/get?key=wx_login')
    wxLoginConfig.value = Object.assign(wxLoginConfig.value, res.data)
    wxLoginConfig.value.notify_url =
      wxLoginConfig.value.notify_url || window.location.origin + '/api/user/login/callback'
  } catch (e) {
    // 使用默认值
  }
}

// 保存行为验证码配置
const saveCaptchaConfig = async () => {
  loading.value = true
  try {
    await httpPost('/api/admin/config/update/captcha', captchaConfig.value)
    ElMessage.success('保存成功')
  } catch (e) {
    ElMessage.error('保存失败：' + (e.message || '未知错误'))
  } finally {
    loading.value = false
  }
}

// 保存微信登录配置
const saveWxloginConfig = async () => {
  loading.value = true
  try {
    await httpPost('/api/admin/config/update/wx_login', wxLoginConfig.value)
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
  .el-form-item__label {
    font-weight: 700;
  }
}
</style>
