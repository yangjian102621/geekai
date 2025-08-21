<template>
  <div class="container py-3 px-10" v-loading="loading">
    <el-tabs v-model="active">
      <el-tab-pane label="支付宝" name="alipay">
        <el-form :model="alipay" label-width="140px">
          <el-form-item label="启用通道"><el-switch v-model="alipay.enabled" /></el-form-item>
          <el-form-item label="沙盒模式"><el-switch v-model="alipay.sand_box" /></el-form-item>
          <el-form-item label="AppId"><el-input v-model="alipay.app_id" /></el-form-item>
          <el-form-item label="商户私钥"
            ><el-input v-model="alipay.private_key" type="textarea" :rows="5"
          /></el-form-item>
          <el-form-item label="支付宝公钥"
            ><el-input v-model="alipay.alipay_public_key" type="textarea" :rows="3"
          /></el-form-item>
          <el-form-item label="异步通知URL"><el-input v-model="alipay.notify_url" /></el-form-item>
          <el-form-item label="同步回跳URL"><el-input v-model="alipay.return_url" /></el-form-item>
          <el-form-item>
            <el-button type="primary" @click="save('alipay')">保存</el-button>
            <el-button @click="test('alipay')">测试</el-button>
          </el-form-item>
        </el-form>
      </el-tab-pane>

      <el-tab-pane label="微信支付" name="wechat">
        <el-form :model="wechat" label-width="140px">
          <el-form-item label="启用通道"><el-switch v-model="wechat.enabled" /></el-form-item>
          <el-form-item label="AppId"><el-input v-model="wechat.app_id" /></el-form-item>
          <el-form-item label="商户号(MchId)"><el-input v-model="wechat.mch_id" /></el-form-item>
          <el-form-item label="证书序列号"><el-input v-model="wechat.serial_no" /></el-form-item>
          <el-form-item label="商户私钥"
            ><el-input v-model="wechat.private_key" type="textarea" :rows="3"
          /></el-form-item>
          <el-form-item label="APIv3 Key"><el-input v-model="wechat.api_v3_key" /></el-form-item>
          <el-form-item label="异步通知URL"><el-input v-model="wechat.notify_url" /></el-form-item>
          <el-form-item>
            <el-button type="primary" @click="save('wechat')">保存</el-button>
            <el-button @click="test('wechat')">测试</el-button>
          </el-form-item>
        </el-form>
      </el-tab-pane>

      <el-tab-pane label="虎皮椒" name="hupi">
        <el-form :model="hupi" label-width="140px">
          <el-form-item label="启用通道"><el-switch v-model="hupi.enabled" /></el-form-item>
          <el-form-item label="AppId"><el-input v-model="hupi.app_id" /></el-form-item>
          <el-form-item label="AppSecret"><el-input v-model="hupi.app_secret" /></el-form-item>
          <el-form-item label="网关地址"><el-input v-model="hupi.api_url" /></el-form-item>
          <el-form-item label="异步通知URL"><el-input v-model="hupi.notify_url" /></el-form-item>
          <el-form-item label="同步回跳URL"><el-input v-model="hupi.return_url" /></el-form-item>
          <el-form-item>
            <el-button type="primary" @click="save('hupi')">保存</el-button>
            <el-button @click="test('hupi')">测试</el-button>
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

const alipay = ref({ enabled: false, sand_box: false })
const wechat = ref({ enabled: false })
const hupi = ref({ enabled: false })
const geekpay = ref({ enabled: false, methods: [] })

onMounted(() => {
  Promise.all([
    httpGet('/api/admin/config/get?key=alipay'),
    httpGet('/api/admin/config/get?key=wechat'),
    httpGet('/api/admin/config/get?key=hupi'),
    httpGet('/api/admin/config/get?key=geekpay'),
  ])
    .then(([a, w, h, g]) => {
      alipay.value = a.data || alipay.value
      wechat.value = w.data || wechat.value
      hupi.value = h.data || hupi.value
      geekpay.value = g.data || geekpay.value
    })
    .catch(() => {})
    .finally(() => (loading.value = false))
})

const save = (key) => {
  const map = { alipay, wechat, hupi, geekpay }
  httpPost('/api/admin/config/update', { key, config: map[key].value })
    .then(() => ElMessage.success('保存成功'))
    .catch((e) => ElMessage.error(e.message))
}

const test = (key) => {
  httpPost('/api/admin/config/test', { key })
    .then((res) => ElMessage.success(res.message || '测试成功'))
    .catch((e) => ElMessage.error(e.message || '测试失败'))
}
</script>

<style scoped>
.form {
  padding: 10px 20px 40px 20px;
}
</style>
