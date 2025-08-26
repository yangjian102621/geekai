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
          <el-form-item label="BasePath"><el-input v-model="local.BasePath" /></el-form-item>
          <el-form-item label="BaseURL"><el-input v-model="local.BaseURL" /></el-form-item>
        </el-form>
      </template>

      <template v-else-if="active === 'minio'">
        <el-form :model="minio" label-width="140px">
          <el-form-item label="Endpoint"><el-input v-model="minio.Endpoint" /></el-form-item>
          <el-form-item label="AccessKey"><el-input v-model="minio.AccessKey" /></el-form-item>
          <el-form-item label="AccessSecret"
            ><el-input v-model="minio.AccessSecret"
          /></el-form-item>
          <el-form-item label="Bucket"><el-input v-model="minio.Bucket" /></el-form-item>
          <el-form-item label="UseSSL"><el-switch v-model="minio.UseSSL" /></el-form-item>
          <el-form-item label="Domain"><el-input v-model="minio.Domain" /></el-form-item>
        </el-form>
      </template>

      <template v-else-if="active === 'qiniu'">
        <el-form :model="qiniu" label-width="140px">
          <el-form-item label="Zone"><el-input v-model="qiniu.Zone" /></el-form-item>
          <el-form-item label="AccessKey"><el-input v-model="qiniu.AccessKey" /></el-form-item>
          <el-form-item label="AccessSecret"
            ><el-input v-model="qiniu.AccessSecret"
          /></el-form-item>
          <el-form-item label="Bucket"><el-input v-model="qiniu.Bucket" /></el-form-item>
          <el-form-item label="Domain"><el-input v-model="qiniu.Domain" /></el-form-item>
        </el-form>
      </template>

      <template v-else>
        <el-form :model="aliyun" label-width="140px">
          <el-form-item label="Endpoint"><el-input v-model="aliyun.Endpoint" /></el-form-item>
          <el-form-item label="AccessKey"><el-input v-model="aliyun.AccessKey" /></el-form-item>
          <el-form-item label="AccessSecret"
            ><el-input v-model="aliyun.AccessSecret"
          /></el-form-item>
          <el-form-item label="Bucket"><el-input v-model="aliyun.Bucket" /></el-form-item>
          <el-form-item label="SubDir"><el-input v-model="aliyun.SubDir" /></el-form-item>
          <el-form-item label="Domain"><el-input v-model="aliyun.Domain" /></el-form-item>
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
const local = ref({ BasePath: '', BaseURL: '' })
const minio = ref({
  Endpoint: '',
  AccessKey: '',
  AccessSecret: '',
  Bucket: '',
  SubDir: '',
  UseSSL: false,
  Domain: '',
})
const qiniu = ref({
  Zone: 'z2',
  AccessKey: '',
  AccessSecret: '',
  Bucket: '',
  SubDir: '',
  Domain: '',
})
const aliyun = ref({
  Endpoint: '',
  AccessKey: '',
  AccessSecret: '',
  Bucket: '',
  SubDir: '',
  Domain: '',
})

onMounted(() => {
  httpGet('/api/admin/config/get?key=oss')
    .then((res) => {
      const data = res.data || {}
      const Active = data.Active || data.active || 'local'
      active.value = String(Active).toLowerCase()
      local.value = data.Local || data.local || local.value
      minio.value = data.Minio || data.minio || minio.value
      qiniu.value = data.QiNiu || data.qiniu || qiniu.value
      aliyun.value = data.AliYun || data.aliyun || aliyun.value
    })
    .catch(() => {})
    .finally(() => (loading.value = false))
})

const save = () => {
  httpPost('/api/admin/config/update/oss', {
    active: active.value,
    local: local.value,
    minio: minio.value,
    qiniu: qiniu.value,
    aliyun: aliyun.value,
  })
    .then(() => ElMessage.success('保存成功'))
    .catch((e) => ElMessage.error(e.message))
}

const test = () => {
  ElMessage.info('请在对象存储端验证配置')
}
</script>

<style scoped>
.form {
  padding: 10px 20px 40px 20px;
}
</style>
