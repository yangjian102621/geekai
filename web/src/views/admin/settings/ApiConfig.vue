<template>
  <div class="form" v-loading="loading">
    <el-form :model="api" label-width="140px">
      <el-form-item label="API 网关"><el-input v-model="api.api_url" /></el-form-item>
      <el-form-item label="AppId"><el-input v-model="api.app_id" /></el-form-item>
      <el-form-item label="Token"><el-input v-model="api.token" type="password" /></el-form-item>

      <el-divider>即梦 AI</el-divider>
      <el-form-item label="AccessKey"
        ><el-input v-model="api.jimeng_config.access_key"
      /></el-form-item>
      <el-form-item label="SecretKey"
        ><el-input v-model="api.jimeng_config.secret_key"
      /></el-form-item>
      <el-form-item>
        <el-button type="primary" @click="save">保存</el-button>
        <el-button @click="test">测试</el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script setup>
import { httpGet, httpPost } from '@/utils/http'
import { ElMessage } from 'element-plus'
import { onMounted, ref } from 'vue'

const loading = ref(true)
const api = ref({
  api_url: '',
  app_id: '',
  token: '',
  jimeng_config: { access_key: '', secret_key: '' },
})

onMounted(() => {
  httpGet('/api/admin/config/get?key=api')
    .then((res) => (api.value = res.data || api.value))
    .catch(() => {})
    .finally(() => (loading.value = false))
})

const save = () => {
  httpPost('/api/admin/config/update', { key: 'api', config: api.value })
    .then(() => ElMessage.success('保存成功'))
    .catch((e) => ElMessage.error(e.message))
}

const test = () => {
  httpPost('/api/admin/config/test', { key: 'api' })
    .then((res) => ElMessage.success(res.message || '测试成功'))
    .catch((e) => ElMessage.error(e.message || '测试失败'))
}
</script>

<style scoped>
.form {
  padding: 10px 20px 40px 20px;
}
</style>
