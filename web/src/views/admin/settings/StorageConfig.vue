<template>
  <div class="settings container p-5">
    <el-tabs v-model="activeTab" type="border-card">
      <el-tab-pane label="本地" name="local">
        <el-form :model="local" label-position="top">
          <el-form-item>
            <label class="form-label"
              >文件存储根目录
              <el-tooltip placement="top">
                <template #content>
                  可以是绝对路径，如：/data/static/upload<br />也可以是相对路径，如：./static/upload
                </template>
                <i class="iconfont icon-info"></i>
              </el-tooltip>
            </label>
            <el-input
              v-model="local.base_path"
              placeholder="请输入文件存储根目录，如：./static/upload"
            />
          </el-form-item>
          <el-form-item>
            <label class="form-label"
              >文件访问根 URL
              <el-tooltip placement="top">
                <template #content>
                  可以是绝对路径，如：https://oss.geekai.me/static/upload
                  <br />也可以是相对路径，如：/static/upload
                </template>
                <i class="iconfont icon-info"></i>
              </el-tooltip>
            </label>
            <el-input
              v-model="local.base_url"
              placeholder="请输入文件存储URL，如：/static/upload"
            />
          </el-form-item>
        </el-form>
      </el-tab-pane>

      <el-tab-pane label="MinIO" name="minio">
        <div class="rounded-md bg-blue-100 p-3 text-gray-500 border-blue-500 border-2 text-base">
          如果你不知道怎么获取这些配置信息，请参考文档：
          <a
            href="https://docs.geekai.me/plus/config/oss.html#%E6%90%AD%E5%BB%BA-minio-%E5%AD%98%E5%82%A8%E6%9C%8D%E5%8A%A1"
            target="_blank"
            >Minio 配置</a
          >。
        </div>
        <el-form :model="minio" class="mt-4" label-position="top">
          <el-form-item label="Endpoint"><el-input v-model="minio.endpoint" /></el-form-item>
          <el-form-item label="AccessKey"><el-input v-model="minio.access_key" /></el-form-item>
          <el-form-item label="AccessSecret"
            ><el-input v-model="minio.access_secret"
          /></el-form-item>
          <el-form-item label="Bucket"><el-input v-model="minio.bucket" /></el-form-item>
          <el-form-item label="UseSSL"><el-switch v-model="minio.use_ssl" /></el-form-item>
          <el-form-item label="Domain"><el-input v-model="minio.domain" /></el-form-item>
        </el-form>
      </el-tab-pane>

      <el-tab-pane label="七牛云" name="qiniu">
        <div class="rounded-md bg-blue-100 p-3 text-gray-500 border-blue-500 border-2 text-base">
          如果你不知道怎么获取这些配置信息，请参考文档：
          <a
            href="https://docs.geekai.me/plus/config/oss.html#%E4%B8%83%E7%89%9B%E4%BA%91-oss-%E9%85%8D%E7%BD%AE"
            target="_blank"
            >七牛云配置</a
          >。
        </div>
        <el-form :model="qiniu" class="mt-4" label-position="top">
          <el-form-item label="Zone">
            <template #label>
              <label
                >区域（Zone）
                <el-tooltip
                  placement="right"
                  content="华南：z2，华东：z0，华北：z1，北美：na0，新加坡：as0"
                >
                  <i class="iconfont icon-info"></i>
                </el-tooltip>
              </label>
            </template>
            <el-input
              v-model="qiniu.zone"
              placeholder="华南：z2，华东：z0，华北：z1，北美：na0，新加坡：as0"
          /></el-form-item>
          <el-form-item label="AccessKey"><el-input v-model="qiniu.access_key" /></el-form-item>
          <el-form-item label="AccessSecret"
            ><el-input v-model="qiniu.access_secret"
          /></el-form-item>
          <el-form-item label="Bucket"><el-input v-model="qiniu.bucket" /></el-form-item>
          <el-form-item label="Domain"
            ><el-input v-model="qiniu.domain" placeholder="请输入七牛云Bucket绑定的域名"
          /></el-form-item>
        </el-form>
      </el-tab-pane>

      <el-tab-pane label="阿里云OSS" name="aliyun">
        <el-form :model="aliyun" class="mt-4" label-position="top">
          <el-form-item label="Endpoint"><el-input v-model="aliyun.endpoint" /></el-form-item>
          <el-form-item label="AccessKey"><el-input v-model="aliyun.access_key" /></el-form-item>
          <el-form-item label="AccessSecret"
            ><el-input v-model="aliyun.access_secret"
          /></el-form-item>
          <el-form-item label="Bucket"><el-input v-model="aliyun.bucket" /></el-form-item>
          <el-form-item label="Domain"><el-input v-model="aliyun.domain" /></el-form-item>
        </el-form>
      </el-tab-pane>
    </el-tabs>

    <div class="mt-3">
      <label class="form-label mr-2">存储引擎</label>
      <el-radio-group v-model="active" size="large">
        <el-radio value="local" border>本地存储</el-radio>
        <el-radio value="aliyun" border>阿里云</el-radio>
        <el-radio value="qiniu" border>七牛云</el-radio>
        <el-radio value="minio" border>Minio</el-radio>
      </el-radio-group>
    </div>

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
const activeTab = ref('local')
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
const qiniu = ref({
  zone: 'z2',
  access_key: '',
  access_secret: '',
  bucket: '',
  domain: '',
})
const aliyun = ref({
  endpoint: '',
  access_key: '',
  access_secret: '',
  bucket: '',
  domain: '',
})

onMounted(() => {
  httpGet('/api/admin/config/get?key=oss')
    .then((res) => {
      const data = res.data || {}
      active.value = data.active.toLowerCase() || active.value
      local.value = data.local || local.value
      minio.value = data.minio || minio.value
      qiniu.value = data.qiniu || qiniu.value
      aliyun.value = data.aliyun || aliyun.value

      minio.value.bucket = minio.value.bucket || 'geekai'
      qiniu.value.bucket = qiniu.value.bucket || 'geekai'
      aliyun.value.bucket = aliyun.value.bucket || 'geekai'
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
