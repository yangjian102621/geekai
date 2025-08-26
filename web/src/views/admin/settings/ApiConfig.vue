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
import { httpGet } from '@/utils/http'
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
  ElMessage.info('当前后端未提供 /api 配置的更新接口，已保留只读展示')
}

const test = () => {
  ElMessage.info('请在对应服务端手动测试 API 可用性')
}
</script>

<style scoped>
.form {
  padding: 10px 20px 40px 20px;
}
</style>
