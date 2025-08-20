<template>
  <div class="form" v-loading="loading">
    <el-form label-width="140px">
      <el-form-item label="存储引擎">
        <el-select v-model="active" style="width: 280px">
          <el-option label="本地" value="local" />
          <el-option label="MinIO" value="minio" />
          <el-option label="七牛云" value="qiniu" />
          <el-option label="阿里云OSS" value="aliyun" />
        </el-select>
      </el-form-item>

      <template v-if="active === 'local'">
        <el-form :model="local" label-width="140px">
          <el-form-item label="BasePath"><el-input v-model="local.base_path" /></el-form-item>
          <el-form-item label="BaseURL"><el-input v-model="local.base_url" /></el-form-item>
        </el-form>
      </template>

      <template v-else-if="active === 'minio'">
        <el-form :model="minio" label-width="140px">
          <el-form-item label="Endpoint"><el-input v-model="minio.endpoint" /></el-form-item>
          <el-form-item label="AccessKey"><el-input v-model="minio.access_key" /></el-form-item>
          <el-form-item label="AccessSecret"
            ><el-input v-model="minio.access_secret"
          /></el-form-item>
          <el-form-item label="Bucket"><el-input v-model="minio.bucket" /></el-form-item>
          <el-form-item label="UseSSL"><el-switch v-model="minio.use_ssl" /></el-form-item>
          <el-form-item label="Domain"><el-input v-model="minio.domain" /></el-form-item>
        </el-form>
      </template>

      <template v-else-if="active === 'qiniu'">
        <el-form :model="qiniu" label-width="140px">
          <el-form-item label="Zone"><el-input v-model="qiniu.zone" /></el-form-item>
          <el-form-item label="AccessKey"><el-input v-model="qiniu.access_key" /></el-form-item>
          <el-form-item label="AccessSecret"
            ><el-input v-model="qiniu.access_secret"
          /></el-form-item>
          <el-form-item label="Bucket"><el-input v-model="qiniu.bucket" /></el-form-item>
          <el-form-item label="Domain"><el-input v-model="qiniu.domain" /></el-form-item>
        </el-form>
      </template>

      <template v-else>
        <el-form :model="aliyun" label-width="140px">
          <el-form-item label="Endpoint"><el-input v-model="aliyun.endpoint" /></el-form-item>
          <el-form-item label="AccessKey"><el-input v-model="aliyun.access_key" /></el-form-item>
          <el-form-item label="AccessSecret"
            ><el-input v-model="aliyun.access_secret"
          /></el-form-item>
          <el-form-item label="Bucket"><el-input v-model="aliyun.bucket" /></el-form-item>
          <el-form-item label="SubDir"><el-input v-model="aliyun.sub_dir" /></el-form-item>
          <el-form-item label="Domain"><el-input v-model="aliyun.domain" /></el-form-item>
        </el-form>
      </template>

      <el-form-item>
        <el-button type="primary" @click="save">保存</el-button>
        <el-button @click="test">连接测试</el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script setup>
import { httpGet, httpPost } from '@/utils/http'
import { ElMessage } from 'element-plus'
import { onMounted, ref } from 'vue'

const loading = ref(true)
const active = ref('local')
const local = ref({ base_path: '', base_url: '' })
const minio = ref({
  endpoint: '',
  access_key: '',
  access_secret: '',
  bucket: '',
  use_ssl: false,
  domain: '',
})
const qiniu = ref({ zone: 'z2', access_key: '', access_secret: '', bucket: '', domain: '' })
const aliyun = ref({
  endpoint: '',
  access_key: '',
  access_secret: '',
  bucket: '',
  sub_dir: '',
  domain: '',
})

onMounted(() => {
  httpGet('/api/admin/config/get?key=oss')
    .then((res) => {
      const data = res.data || {}
      active.value = (data.active || 'local').toLowerCase()
      local.value = data.local || local.value
      minio.value = data.minio || minio.value
      qiniu.value = data.qiniu || qiniu.value
      aliyun.value = data.aliyun || aliyun.value
    })
    .catch(() => {})
    .finally(() => (loading.value = false))
})

const save = () => {
  httpPost('/api/admin/config/update', {
    key: 'oss',
    config: {
      active: active.value,
      local: local.value,
      minio: minio.value,
      qiniu: qiniu.value,
      aliyun: aliyun.value,
    },
  })
    .then(() => ElMessage.success('保存成功'))
    .catch((e) => ElMessage.error(e.message))
}

const test = () => {
  httpPost('/api/admin/config/test', { key: 'oss' })
    .then((res) => ElMessage.success(res.message || '连接成功'))
    .catch((e) => ElMessage.error(e.message || '连接失败'))
}
</script>

<style scoped>
.form {
  padding: 10px 20px 40px 20px;
}
</style>
