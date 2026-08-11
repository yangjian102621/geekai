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
          <el-form-item>
            <label class="form-label"
              >缩略图模板
              <el-tooltip placement="top">
                <template #content>
                  使用{width}和{height}作为变量占位符<br />默认：?imageView2/4/w/{width}/h/{height}/q/75<br />留空则使用默认模板
                </template>
                <i class="iconfont icon-info"></i>
              </el-tooltip>
            </label>
            <el-input
              v-model="local.thumb_template"
              placeholder="?imageView2/4/w/{width}/h/{height}/q/75"
            />
          </el-form-item>
        </el-form>
      </el-tab-pane>

      <el-tab-pane label="MinIO" name="minio">
        <Alert type="info">
          如果你不知道怎么获取这些配置信息，请参考文档：
          <a
            href="https://docs.geekai.me/plus/config/oss.html#%E6%90%AD%E5%BB%BA-minio-%E5%AD%98%E5%82%A8%E6%9C%8D%E5%8A%A1"
            target="_blank"
            >Minio 配置</a
          >。
        </Alert>
        <el-form :model="minio" class="mt-4" label-position="top">
          <el-form-item label="Endpoint"><el-input v-model="minio.endpoint" /></el-form-item>
          <el-form-item label="AccessKey"><el-input v-model="minio.access_key" /></el-form-item>
          <el-form-item label="AccessSecret"
            ><el-input v-model="minio.access_secret"
          /></el-form-item>
          <el-form-item label="Bucket"><el-input v-model="minio.bucket" /></el-form-item>
          <el-form-item label="UseSSL"><el-switch v-model="minio.use_ssl" /></el-form-item>
          <el-form-item label="Domain"><el-input v-model="minio.domain" /></el-form-item>
          <el-form-item>
            <template #label>
              <label
                >缩略图模板
                <el-tooltip
                  placement="top"
                  content="使用{width}和{height}作为变量占位符。MinIO不支持缩略图，留空则返回原图"
                >
                  <i class="iconfont icon-info"></i>
                </el-tooltip>
              </label>
            </template>
            <el-input
              v-model="minio.thumb_template"
              placeholder="留空则返回原图（MinIO不支持缩略图）"
            />
          </el-form-item>
        </el-form>
      </el-tab-pane>

      <el-tab-pane label="七牛云" name="qiniu">
        <Alert type="info">
          如果你不知道怎么获取这些配置信息，请参考文档：
          <a
            href="https://docs.geekai.me/plus/config/oss.html#%E4%B8%83%E7%89%9B%E4%BA%91-oss-%E9%85%8D%E7%BD%AE"
            target="_blank"
            >七牛云配置</a
          >。
        </Alert>
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
          <el-form-item>
            <template #label>
              <label
                >缩略图模板
                <el-tooltip
                  placement="top"
                  content="使用{width}和{height}作为变量占位符。默认：?imageView2/4/w/{width}/h/{height}/q/75"
                >
                  <i class="iconfont icon-info"></i>
                </el-tooltip>
              </label>
            </template>
            <el-input
              v-model="qiniu.thumb_template"
              placeholder="?imageView2/4/w/{width}/h/{height}/q/75"
            />
          </el-form-item>
        </el-form>
      </el-tab-pane>

      <el-tab-pane label="阿里云OSS" name="aliyun">
        <Alert type="info">
          如果你不知道怎么获取这些配置信息，请参考文档：
          <a
            href="https://docs.geekai.me/plus/config/oss.html#%E9%98%BF%E9%87%8C%E4%BA%91-oss-%E9%85%8D%E7%BD%AE"
            target="_blank"
            >阿里云OSS配置</a
          >。
        </Alert>
        <el-form :model="aliyun" class="mt-4" label-position="top">
          <el-form-item label="Endpoint"><el-input v-model="aliyun.endpoint" /></el-form-item>
          <el-form-item label="AccessKey"><el-input v-model="aliyun.access_key" /></el-form-item>
          <el-form-item label="AccessSecret"
            ><el-input v-model="aliyun.access_secret"
          /></el-form-item>
          <el-form-item label="Bucket"><el-input v-model="aliyun.bucket" /></el-form-item>
          <el-form-item label="Domain"><el-input v-model="aliyun.domain" /></el-form-item>
          <el-form-item>
            <template #label>
              <label
                >缩略图模板
                <el-tooltip
                  placement="top"
                  content="使用{width}和{height}作为变量占位符。默认：?x-oss-process=image/resize,m_lfit,w_{width},h_{height}"
                >
                  <i class="iconfont icon-info"></i>
                </el-tooltip>
              </label>
            </template>
            <el-input
              v-model="aliyun.thumb_template"
              placeholder="?x-oss-process=image/resize,m_lfit,w_{width},h_{height}"
            />
          </el-form-item>
        </el-form>
      </el-tab-pane>

      <el-tab-pane label="腾讯云COS" name="tencent">
        <Alert type="info">
          如果你不知道怎么获取这些配置信息，请参考文档：
          <a href="https://cloud.tencent.com/document/product/436" target="_blank">腾讯云COS配置</a
          >。
        </Alert>
        <el-form :model="tencent" class="mt-4" label-position="top">
          <el-form-item label="Region">
            <el-select
              v-model="tencent.region"
              filterable
              allow-create
              placeholder="请选择或输入区域"
              style="width: 100%"
            >
              <el-option
                v-for="region in tencentRegions"
                :key="region.value"
                :label="region.label"
                :value="region.value"
              />
            </el-select>
          </el-form-item>
          <el-form-item label="SecretId"><el-input v-model="tencent.secret_id" /></el-form-item>
          <el-form-item label="SecretKey"
            ><el-input v-model="tencent.secret_key" type="password" show-password
          /></el-form-item>
          <el-form-item label="Bucket"><el-input v-model="tencent.bucket" /></el-form-item>
          <el-form-item label="Domain"
            ><el-input
              v-model="tencent.domain"
              placeholder="请输入自定义域名，留空则使用COS默认域名"
          /></el-form-item>
          <el-form-item>
            <template #label>
              <label
                >缩略图模板
                <el-tooltip
                  placement="top"
                  content="使用{width}和{height}作为变量占位符。默认：?imageView2/1/w/{width}/h/{height}/format/jpg"
                >
                  <i class="iconfont icon-info"></i>
                </el-tooltip>
              </label>
            </template>
            <el-input
              v-model="tencent.thumb_template"
              placeholder="?imageView2/1/w/{width}/h/{height}/format/jpg"
            />
          </el-form-item>
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
        <el-radio value="tencent" border>腾讯云</el-radio>
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
import Alert from '@/components/ui/Alert.vue'

const loading = ref(true)
const activeTab = ref('local')
const active = ref('local')
// 默认缩略图模板
const DEFAULT_LOCAL_THUMB = '?imageView2/4/w/{width}/h/{height}/q/75'
const DEFAULT_QINIU_THUMB = '?imageView2/4/w/{width}/h/{height}/q/75'
const DEFAULT_ALIYUN_THUMB = '?x-oss-process=image/resize,m_lfit,w_{width},h_{height}'
const DEFAULT_TENCENT_THUMB = '?imageView2/1/w/{width}/h/{height}/format/jpg'

const local = ref({
  base_path: './static/upload',
  base_url: '/static/upload',
  thumb_template: DEFAULT_LOCAL_THUMB,
})
const minio = ref({
  endpoint: '',
  access_key: '',
  access_secret: '',
  bucket: '',
  use_ssl: false,
  domain: '',
  thumb_template: '', // MinIO不支持缩略图，留空
})
const qiniu = ref({
  zone: 'z2',
  access_key: '',
  access_secret: '',
  bucket: '',
  domain: '',
  thumb_template: DEFAULT_QINIU_THUMB,
})
const aliyun = ref({
  endpoint: '',
  access_key: '',
  access_secret: '',
  bucket: '',
  domain: '',
  thumb_template: DEFAULT_ALIYUN_THUMB,
})
// 腾讯云 COS 常见地区选项
const tencentRegions = [
  { label: '北京 (ap-beijing)', value: 'ap-beijing' },
  { label: '南京 (ap-nanjing)', value: 'ap-nanjing' },
  { label: '上海 (ap-shanghai)', value: 'ap-shanghai' },
  { label: '广州 (ap-guangzhou)', value: 'ap-guangzhou' },
  { label: '成都 (ap-chengdu)', value: 'ap-chengdu' },
  { label: '重庆 (ap-chongqing)', value: 'ap-chongqing' },
  { label: '中国香港 (ap-hongkong)', value: 'ap-hongkong' },
  { label: '新加坡 (ap-singapore)', value: 'ap-singapore' },
  { label: '雅加达 (ap-jakarta)', value: 'ap-jakarta' },
  { label: '首尔 (ap-seoul)', value: 'ap-seoul' },
  { label: '曼谷 (ap-bangkok)', value: 'ap-bangkok' },
  { label: '东京 (ap-tokyo)', value: 'ap-tokyo' },
]

const tencent = ref({
  region: 'ap-beijing',
  secret_id: '',
  secret_key: '',
  bucket: '',
  domain: '',
  thumb_template: DEFAULT_TENCENT_THUMB,
})

onMounted(() => {
  httpGet('/api/admin/config/get?key=oss')
    .then((res) => {
      const data = res.data || {}
      active.value = data.active.toLowerCase() || active.value

      // 合并配置，如果服务器返回的配置中没有thumb_template或为空，则使用默认值
      if (data.local) {
        local.value = { ...local.value, ...data.local }
        if (!data.local.thumb_template) {
          local.value.thumb_template = DEFAULT_LOCAL_THUMB
        }
      }

      if (data.minio) {
        minio.value = { ...minio.value, ...data.minio }
        // MinIO不支持缩略图，保持为空
      }

      if (data.qiniu) {
        qiniu.value = { ...qiniu.value, ...data.qiniu }
        if (!data.qiniu.thumb_template) {
          qiniu.value.thumb_template = DEFAULT_QINIU_THUMB
        }
      }

      if (data.aliyun) {
        aliyun.value = { ...aliyun.value, ...data.aliyun }
        if (!data.aliyun.thumb_template) {
          aliyun.value.thumb_template = DEFAULT_ALIYUN_THUMB
        }
      }

      if (data.tencent) {
        tencent.value = { ...tencent.value, ...data.tencent }
        if (!data.tencent.thumb_template) {
          tencent.value.thumb_template = DEFAULT_TENCENT_THUMB
        }
      }

      minio.value.bucket = minio.value.bucket || 'geekai'
      qiniu.value.bucket = qiniu.value.bucket || 'geekai'
      aliyun.value.bucket = aliyun.value.bucket || 'geekai'
      tencent.value.bucket = tencent.value.bucket || 'geekai'
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
    tencent: tencent.value,
  })
    .then(() => {
      ElMessage.success('保存成功')
    })
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
