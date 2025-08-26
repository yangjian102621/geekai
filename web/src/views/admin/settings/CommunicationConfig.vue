<template>
  <div class="form" v-loading="loading">
    <el-tabs v-model="active">
      <el-tab-pane label="SMTP 邮件" name="smtp">
        <el-form :model="smtp" label-width="140px">
          <el-form-item label="启用TLS"><el-switch v-model="smtp.use_tls" /></el-form-item>
          <el-form-item label="SMTP服务器"><el-input v-model="smtp.host" /></el-form-item>
          <el-form-item label="端口"><el-input-number v-model="smtp.port" :min="1" /></el-form-item>
          <el-form-item label="应用名称"><el-input v-model="smtp.app_name" /></el-form-item>
          <el-form-item label="发件人地址"><el-input v-model="smtp.from" /></el-form-item>
          <el-form-item label="发件人密码"
            ><el-input v-model="smtp.password" type="password"
          /></el-form-item>
          <el-form-item>
            <el-button type="primary" @click="save('smtp')">保存</el-button>
            <el-button @click="test('smtp')">测试</el-button>
          </el-form-item>
        </el-form>
      </el-tab-pane>

      <el-tab-pane label="短信服务" name="sms">
        <el-form :model="sms" label-width="140px">
          <el-form-item label="服务商">
            <el-select v-model="sms.Active" style="width: 200px">
              <el-option label="阿里云" value="Ali" />
              <el-option label="短信宝" value="Bao" />
            </el-select>
          </el-form-item>
          <template v-if="sms.Active === 'Ali'">
            <el-form-item label="AccessKey"><el-input v-model="sms.Ali.AccessKey" /></el-form-item>
            <el-form-item label="AccessSecret"
              ><el-input v-model="sms.Ali.AccessSecret"
            /></el-form-item>
            <el-form-item label="签名"><el-input v-model="sms.Ali.Sign" /></el-form-item>
            <el-form-item label="模板ID"><el-input v-model="sms.Ali.CodeTempId" /></el-form-item>
          </template>
          <template v-else>
            <el-form-item label="用户名"><el-input v-model="sms.Bao.Username" /></el-form-item>
            <el-form-item label="密码"
              ><el-input v-model="sms.Bao.Password" type="password"
            /></el-form-item>
            <el-form-item label="签名"><el-input v-model="sms.Bao.Sign" /></el-form-item>
            <el-form-item label="模板"
              ><el-input v-model="sms.Bao.CodeTemplate" type="textarea" :rows="2"
            /></el-form-item>
          </template>
          <el-form-item>
            <el-button type="primary" @click="save('sms')">保存</el-button>
            <el-button @click="test('sms')">测试</el-button>
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
const active = ref('smtp')
const smtp = ref({ use_tls: false, host: '', port: 25, app_name: '', from: '', password: '' })
const sms = ref({
  Active: 'Ali',
  Ali: { AccessKey: '', AccessSecret: '', Sign: '', CodeTempId: '' },
  Bao: { Username: '', Password: '', Sign: '', CodeTemplate: '' },
})

onMounted(() => {
  Promise.all([httpGet('/api/admin/config/get?key=smtp'), httpGet('/api/admin/config/get?key=sms')])
    .then(([s1, s2]) => {
      const smtpData = s1?.data || {}
      smtp.value = { ...smtp.value, ...smtpData }

      const smsData = s2?.data || {}
      sms.value = {
        ...sms.value,
        ...smsData,
        Ali: { ...sms.value.Ali, ...(smsData.Ali || smsData.ali || {}) },
        Bao: { ...sms.value.Bao, ...(smsData.Bao || smsData.bao || {}) },
      }
    })
    .catch(() => {})
    .finally(() => (loading.value = false))
})

const save = (key) => {
  if (key === 'smtp') {
    httpPost('/api/admin/config/update/smtp', smtp.value)
      .then(() => ElMessage.success('保存成功'))
      .catch((e) => ElMessage.error(e.message))
  } else if (key === 'sms') {
    httpPost('/api/admin/config/update/sms', sms.value)
      .then(() => ElMessage.success('保存成功'))
      .catch((e) => ElMessage.error(e.message))
  }
}

const test = (key) => {
  ElMessage.info('请在对应服务侧进行验证')
}
</script>

<style scoped>
.form {
  padding: 10px 20px 40px 20px;
}
</style>
