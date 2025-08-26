<template>
  <div class="container py-3 px-10" v-loading="loading">
    <el-tabs v-model="active">
      <el-tab-pane label="支付宝" name="alipay">
        <el-form :model="alipay" label-width="140px">
          <el-form-item label="启用通道"><el-switch v-model="alipay.enabled" /></el-form-item>
          <el-form-item label="沙盒模式"><el-switch v-model="alipay.sandbox" /></el-form-item>
          <el-form-item label="AppId"><el-input v-model="alipay.app_id" /></el-form-item>
          <el-form-item label="商户私钥"
            ><el-input v-model="alipay.private_key" type="textarea" :rows="5"
          /></el-form-item>
          <el-form-item label="支付宝公钥"
            ><el-input v-model="alipay.alipay_public_key" type="textarea" :rows="3"
          /></el-form-item>
          <el-form-item label="回调域名"><el-input v-model="alipay.domain" /></el-form-item>
          <el-form-item>
            <el-button type="primary" @click="save('alipay')">保存</el-button>
            <el-button @click="test('alipay')">测试</el-button>
          </el-form-item>
        </el-form>
      </el-tab-pane>

      <el-tab-pane label="微信支付" name="wxpay">
        <el-form :model="wechat" label-width="140px">
          <el-form-item label="启用通道"><el-switch v-model="wechat.enabled" /></el-form-item>
          <el-form-item label="AppId"><el-input v-model="wechat.app_id" /></el-form-item>
          <el-form-item label="商户号(MchId)"><el-input v-model="wechat.mch_id" /></el-form-item>
          <el-form-item label="证书序列号"><el-input v-model="wechat.serial_no" /></el-form-item>
          <el-form-item label="商户私钥"
            ><el-input v-model="wechat.private_key" type="textarea" :rows="3"
          /></el-form-item>
          <el-form-item label="APIv3 Key"><el-input v-model="wechat.api_v3_key" /></el-form-item>
          <el-form-item label="回调域名"><el-input v-model="wechat.domain" /></el-form-item>
          <el-form-item>
            <el-button type="primary" @click="save('wxpay')">保存</el-button>
            <el-button @click="test('wxpay')">测试</el-button>
          </el-form-item>
        </el-form>
      </el-tab-pane>

      <el-tab-pane label="易支付" name="epay">
        <el-form :model="epay" label-width="140px">
          <el-form-item label="启用通道"><el-switch v-model="epay.enabled" /></el-form-item>
          <el-form-item label="商户ID"><el-input v-model="epay.app_id" /></el-form-item>
          <el-form-item label="商户私钥"
            ><el-input v-model="epay.private_key" type="textarea" :rows="3"
          /></el-form-item>
          <el-form-item label="网关地址"><el-input v-model="epay.api_url" /></el-form-item>
          <el-form-item label="回调域名"><el-input v-model="epay.domain" /></el-form-item>
          <el-form-item>
            <el-button type="primary" @click="save('epay')">保存</el-button>
            <el-button @click="test('epay')">测试</el-button>
          </el-form-item>
        </el-form>
      </el-tab-pane>

      <el-tab-pane label="GeekPay" name="geekpay">
        <el-form :model="geekpay" label-width="140px">
          <el-form-item label="启用通道"><el-switch v-model="geekpay.enabled" /></el-form-item>
          <el-form-item label="商户ID"><el-input v-model="geekpay.app_id" /></el-form-item>
          <el-form-item label="商户私钥"
            ><el-input v-model="geekpay.private_key" type="textarea" :rows="3"
          /></el-form-item>
          <el-form-item label="网关地址"><el-input v-model="geekpay.api_url" /></el-form-item>
          <el-form-item label="异步通知URL"><el-input v-model="geekpay.notify_url" /></el-form-item>
          <el-form-item label="同步回跳URL"><el-input v-model="geekpay.return_url" /></el-form-item>
          <el-form-item label="支付方式"
            ><items-input v-model:value="geekpay.methods"
          /></el-form-item>
          <el-form-item>
            <el-button type="primary" @click="save('geekpay')">保存</el-button>
            <el-button @click="test('geekpay')">测试</el-button>
          </el-form-item>
        </el-form>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script setup>
import ItemsInput from '@/components/ui/ItemsInput.vue'
import { httpGet, httpPost } from '@/utils/http'
import { ElMessage } from 'element-plus'
import { onMounted, ref } from 'vue'

const loading = ref(true)
const active = ref('alipay')

const alipay = ref({
  enabled: false,
  sandbox: false,
  app_id: '',
  private_key: '',
  alipay_public_key: '',
  domain: '',
})
const wechat = ref({
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
      wechat.value = { ...wechat.value, ...(data.wxpay || data.wechat || {}) }
      epay.value = { ...epay.value, ...(data.epay || {}) }
    })
    .catch(() => {})
    .finally(() => (loading.value = false))
})

const save = (key) => {
  const payload = { alipay: alipay.value, wxpay: wechat.value, epay: epay.value }
  httpPost('/api/admin/config/update/payment', payload)
    .then(() => ElMessage.success('保存成功'))
    .catch((e) => ElMessage.error(e.message))
}

const test = (key) => {
  // 后端未提供独立测试接口，保留占位提示
  ElMessage.info('请在支付端自行验证配置')
}
</script>

<style scoped>
.form {
  padding: 10px 20px 40px 20px;
}
</style>
