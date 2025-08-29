<template>
  <div class="settings container p-5">
    <el-tabs v-model="active" type="border-card">
      <el-tab-pane label="支付宝" name="alipay">
        <template #label>
          <div class="d-flex align-items-center text-blue-600">
            <i class="iconfont icon-alipay"></i>
            <span class="ms-2">支付宝</span>
          </div>
        </template>
        <div class="rounded-md bg-blue-100 p-3 text-gray-500 border-blue-500 border-2 text-base">
          如果你不知道怎么获取这些配置信息，请参考文档：
          <a
            href="https://docs.geekai.me/plus/config/payment.html#%E6%94%AF%E4%BB%98%E5%AE%9D%E9%85%8D%E7%BD%AE"
            target="_blank"
            >支付宝配置</a
          >。
        </div>

        <el-form :model="alipay" class="mt-4" label-position="top">
          <el-form-item label="商户ID"><el-input v-model="alipay.app_id" /></el-form-item>
          <el-form-item label="商户私钥"
            ><el-input v-model="alipay.private_key" type="textarea" :rows="5"
          /></el-form-item>
          <el-form-item label="支付宝公钥"
            ><el-input v-model="alipay.alipay_public_key" type="textarea" :rows="3"
          /></el-form-item>
          <el-form-item>
            <template #label>
              <label class="form-label"
                >支付回调域名
                <el-tooltip
                  placement="right"
                  content="请确保回调域名已备案且在支付宝应用添加了白名单"
                >
                  <i class="iconfont icon-info"></i>
                </el-tooltip>
              </label>
            </template>
            <el-input v-model="alipay.domain" />
          </el-form-item>
          <el-form-item label="启用该支付通道"><el-switch v-model="alipay.enabled" /></el-form-item>
          <el-form-item label="启用沙盒模式"><el-switch v-model="alipay.sandbox" /></el-form-item>
        </el-form>
      </el-tab-pane>

      <el-tab-pane label="微信支付" name="wxpay">
        <template #label>
          <div class="d-flex align-items-center text-green-600">
            <i class="iconfont icon-wechat-pay"></i>
            <span class="ms-2">微信支付</span>
          </div>
        </template>

        <div class="rounded-md bg-blue-100 p-3 text-gray-500 border-blue-500 border-2 text-base">
          如果你不知道怎么获取这些配置信息，请参考文档：
          <a
            href="https://docs.geekai.me/plus/config/payment.html#%E5%BE%AE%E4%BF%A1%E6%94%AF%E4%BB%98%E9%85%8D%E7%BD%AE"
            target="_blank"
            >微信支付配置</a
          >。
        </div>

        <el-form :model="wxpay" class="mt-4" label-position="top">
          <el-form-item label="AppID"><el-input v-model="wxpay.app_id" /></el-form-item>
          <el-form-item label="商户号(MchId)"><el-input v-model="wxpay.mch_id" /></el-form-item>
          <el-form-item label="证书序列号"><el-input v-model="wxpay.serial_no" /></el-form-item>
          <el-form-item label="商户私钥"
            ><el-input v-model="wxpay.private_key" type="textarea" :rows="3"
          /></el-form-item>
          <el-form-item label="API V3 密钥"><el-input v-model="wxpay.api_v3_key" /></el-form-item>
          <el-form-item>
            <template #label>
              <label class="form-label">回调域名</label>
              <el-tooltip placement="right" content="请确保回调域名已备案且在微信应用添加了白名单">
                <i class="iconfont icon-info ml-2"></i>
              </el-tooltip>
            </template>
            <el-input v-model="wxpay.domain" />
          </el-form-item>
          <el-form-item label="启用该支付通道"><el-switch v-model="wxpay.enabled" /></el-form-item>
        </el-form>
      </el-tab-pane>

      <el-tab-pane label="易支付" name="epay">
        <template #label>
          <div class="d-flex align-items-center text-purple-600">
            <i class="iconfont icon-reward"></i>
            <span class="ms-2">易支付</span>
          </div>
        </template>
        <div class="rounded-md bg-blue-100 p-3 text-gray-500 border-blue-500 border-2 text-base">
          如果你不知道怎么获取这些配置信息，请参考文档：
          <a
            href="https://docs.geekai.me/plus/config/payment.html#%E6%98%93%E6%94%AF%E4%BB%98%E5%BC%80%E9%80%9A"
            target="_blank"
            >易支付配置</a
          >。
        </div>

        <el-form :model="epay" class="mt-4" label-position="top">
          <el-form-item label="商户ID"><el-input v-model="epay.app_id" /></el-form-item>
          <el-form-item label="商户私钥"><el-input v-model="epay.private_key" /></el-form-item>
          <el-form-item label="网关地址"><el-input v-model="epay.api_url" /></el-form-item>
          <el-form-item label="回调域名"><el-input v-model="epay.domain" /></el-form-item>
          <el-form-item label="启用该支付通道"><el-switch v-model="epay.enabled" /></el-form-item>
        </el-form>
      </el-tab-pane>
    </el-tabs>
    <div class="flex justify-center mt-6">
      <el-button type="primary" @click="save" :loading="loading">提交保存</el-button>
    </div>
  </div>
</template>

<script setup>
import { httpGet, httpPost } from '@/utils/http'
import { ElMessage } from 'element-plus'
import { onMounted, ref } from 'vue'

const loading = ref(true)
const active = ref('alipay')
const domain = computed(() => {
  return window.location.origin
})

const alipay = ref({
  enabled: false,
  sandbox: false,
  app_id: '',
  private_key: '',
  alipay_public_key: '',
  domain: '',
})
const wxpay = ref({
  enabled: false,
  app_id: '',
  mch_id: '',
  serial_no: '',
  private_key: '',
  api_v3_key: '',
  domain: '',
})
const epay = ref({ enabled: false, app_id: '', private_key: '', api_url: '', domain: '' })

onMounted(() => {
  httpGet('/api/admin/config/get?key=payment')
    .then((res) => {
      const data = res.data || {}
      alipay.value = { ...alipay.value, ...(data.alipay || {}) }
      wxpay.value = { ...wxpay.value, ...(data.wxpay || data.wechat || {}) }
      epay.value = { ...epay.value, ...(data.epay || {}) }

      // 如果 domain 为空，则设置为当前域名
      if (!alipay.value.domain) {
        alipay.value.domain = domain.value
      }
      if (!wxpay.value.domain) {
        wxpay.value.domain = domain.value
      }
      if (!epay.value.domain) {
        epay.value.domain = domain.value
      }
    })
    .catch(() => {})
    .finally(() => (loading.value = false))
})

const save = () => {
  loading.value = true
  const payload = { alipay: alipay.value, wxpay: wxpay.value, epay: epay.value }
  httpPost('/api/admin/config/update/payment', payload)
    .then(() => ElMessage.success('保存成功'))
    .catch((e) => ElMessage.error(e.message))
    .finally(() => (loading.value = false))
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
