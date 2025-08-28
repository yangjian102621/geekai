<template>
  <div class="container py-3 px-10" v-loading="loading">
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

        <el-form :model="alipay" label-width="140px" label-position="top">
          <el-form-item label="商户ID"><el-input v-model="alipay.app_id" /></el-form-item>
          <el-form-item label="商户私钥"
            ><el-input v-model="alipay.private_key" type="textarea" :rows="5"
          /></el-form-item>
          <el-form-item label="支付宝公钥"
            ><el-input v-model="alipay.alipay_public_key" type="textarea" :rows="3"
          /></el-form-item>
          <el-form-item
            label="回调域名（<span class='text-red-500'>请确保回调域名已备案且在支付宝应用添加了白名单</span>）"
          >
            <template #label>
              <label class="form-label"
                >支付回调域名
                <el-tooltip
                  placement="top"
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
          <el-form-item>
            <el-button type="primary" @click="save('alipay')">保存</el-button>
            <el-button @click="test('alipay')">测试</el-button>
          </el-form-item>
        </el-form>
      </el-tab-pane>

      <el-tab-pane label="微信支付" name="wxpay">
        <el-form :model="wxpay" label-width="140px">
          <el-form-item label="启用通道"><el-switch v-model="wxpay.enabled" /></el-form-item>
          <el-form-item label="AppId"><el-input v-model="wxpay.app_id" /></el-form-item>
          <el-form-item label="商户号(MchId)"><el-input v-model="wxpay.mch_id" /></el-form-item>
          <el-form-item label="证书序列号"><el-input v-model="wxpay.serial_no" /></el-form-item>
          <el-form-item label="商户私钥"
            ><el-input v-model="wxpay.private_key" type="textarea" :rows="3"
          /></el-form-item>
          <el-form-item label="APIv3 Key"><el-input v-model="wxpay.api_v3_key" /></el-form-item>
          <el-form-item label="回调域名"><el-input v-model="wxpay.domain" /></el-form-item>
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

      <el-tab-pane label="易支付" name="epay">
        <el-form :model="epay" label-width="140px">
          <el-form-item label="启用通道"><el-switch v-model="epay.enabled" /></el-form-item>
          <el-form-item label="商户ID"><el-input v-model="epay.app_id" /></el-form-item>
          <el-form-item label="商户私钥"
            ><el-input v-model="epay.private_key" type="textarea" :rows="3"
          /></el-form-item>
          <el-form-item label="网关地址"><el-input v-model="epay.api_url" /></el-form-item>
          <el-form-item label="回调域名（请确保回调域名已备案且在支付宝应用添加了白名单）"
            ><el-input v-model="epay.domain"
          /></el-form-item>
          <el-form-item>
            <el-button type="primary" @click="save('epay')">保存</el-button>
            <el-button @click="test('epay')">测试</el-button>
          </el-form-item>
        </el-form>
      </el-tab-pane>
    </el-tabs>
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
  const payload = { alipay: alipay.value, wxpay: wxpay.value, epay: epay.value }
  httpPost('/api/admin/config/update/payment', payload)
    .then(() => ElMessage.success('保存成功'))
    .catch((e) => ElMessage.error(e.message))
}
</script>

<style lang="scss" scoped>
.container {
  form {
    padding: 20px;
  }
  a {
    color: #409eff;
    &:hover {
      text-decoration: underline;
    }
  }
}
</style>
