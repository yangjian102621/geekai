<template>
  <div class="settings container p-5">
    <el-tabs v-model="activeTab" type="border-card">
      <el-tab-pane label="阿里云短信" name="aliyun">
        <template #label>
          <div class="flex items-center">
            <i class="iconfont icon-aliyun"></i>
            <span class="ml-2">阿里云短信</span>
          </div>
        </template>

        <Alert type="info">
          如果你不知道怎么获取这些配置信息，请参考文档：
          <a
            href="https://docs.geekai.me/plus/config/sms.html#%E9%98%BF%E9%87%8C%E4%BA%91"
            target="_blank"
            >阿里云短信配置</a
          >。
        </Alert>

        <el-form :model="configs.aliyun" label-position="top">
          <el-form-item label="AccessKey">
            <el-input v-model="configs.aliyun.access_key" placeholder="请输入AccessKey" />
          </el-form-item>
          <el-form-item label="AccessSecret">
            <el-input v-model="configs.aliyun.access_secret" placeholder="请输入AccessSecret" />
          </el-form-item>
          <el-form-item label="短信签名">
            <el-input v-model="configs.aliyun.sign" placeholder="请输入短信签名" />
          </el-form-item>
          <el-form-item label="验证码模板ID">
            <el-input v-model="configs.aliyun.code_temp_id" placeholder="请输入验证码模板ID" />
          </el-form-item>
        </el-form>
      </el-tab-pane>

      <el-tab-pane label="短信宝" name="bao">
        <template #label>
          <div class="flex items-center">
            <i class="iconfont icon-sms"></i>
            <span class="ml-2">短信宝</span>
          </div>
        </template>

        <Alert type="info">
          如果你不知道怎么获取这些配置信息，请参考文档：
          <a
            href="https://docs.geekai.me/plus/config/sms.html#%E7%9F%AD%E4%BF%A1%E5%AE%9D"
            target="_blank"
            >短信宝配置</a
          >。
        </Alert>

        <el-form :model="configs.bao" label-position="top">
          <el-form-item label="用户名">
            <el-input v-model="configs.bao.username" placeholder="请输入用户名" />
          </el-form-item>
          <el-form-item label="密码">
            <el-input v-model="configs.bao.password" type="password" show-password placeholder="请输入密码" />
          </el-form-item>
          <el-form-item label="短信签名">
            <el-input v-model="configs.bao.sign" placeholder="请输入短信签名" />
          </el-form-item>
          <el-form-item label="验证码模板">
            <el-input v-model="configs.bao.code_template" placeholder="请输入验证码模板" />
          </el-form-item>
        </el-form>
      </el-tab-pane>

      <el-tab-pane label="腾讯云短信" name="tencent">
        <template #label>
          <div class="flex items-center">
            <i class="iconfont icon-tencent"></i>
            <span class="ml-2">腾讯云短信</span>
          </div>
        </template>

        <Alert type="info">
          如果你不知道怎么获取这些配置信息，请参考文档：
          <a
            href="https://docs.geekai.me/plus/config/sms.html#%E8%85%BE%E8%AE%AF%E4%BA%91"
            target="_blank"
            >腾讯云短信配置</a
          >。
        </Alert>

        <el-form :model="configs.tencent" label-position="top">
          <el-form-item label="SecretId">
            <el-input v-model="configs.tencent.secret_id" placeholder="请输入SecretId" />
          </el-form-item>
          <el-form-item label="SecretKey">
            <el-input
              v-model="configs.tencent.secret_key"
              type="password"
              show-password
              placeholder="请输入SecretKey"
            />
          </el-form-item>
          <el-form-item label="应用ID (SmsSdkAppId)">
            <el-input v-model="configs.tencent.sms_sdk_app_id" placeholder="请输入短信应用ID" />
          </el-form-item>
          <el-form-item label="短信签名">
            <el-input v-model="configs.tencent.sign" placeholder="请输入短信签名" />
          </el-form-item>
          <el-form-item label="验证码模板ID">
            <el-input v-model="configs.tencent.code_temp_id" placeholder="请输入验证码模板ID" />
            <div class="text-sm text-gray-500 mt-1">在腾讯云短信控制台创建的模板ID</div>
          </el-form-item>
          <el-form-item label="地区 (Region)">
            <el-input
              v-model="configs.tencent.region"
              placeholder="请输入地区，默认 ap-guangzhou"
            />
            <div class="text-sm text-gray-500 mt-1">
              常见地区：ap-guangzhou（广州）、ap-beijing（北京）、ap-shanghai（上海）
            </div>
          </el-form-item>
        </el-form>
      </el-tab-pane>
    </el-tabs>

    <div class="mt-3">
      <label class="form-label mr-2">默认使用</label>
      <el-radio-group v-model="configs.active" size="large">
        <el-radio value="aliyun" border>阿里云</el-radio>
        <el-radio value="bao" border>短信宝</el-radio>
        <el-radio value="tencent" border>腾讯云</el-radio>
      </el-radio-group>
    </div>

    <div class="flex justify-center mt-6">
      <el-button type="primary" @click="saveSmsConfig" :loading="loading">提交保存</el-button>
    </div>
  </div>
</template>

<script setup>
import Alert from '@/components/ui/Alert.vue'
import { httpGet, httpPost } from '@/utils/http'
import { ElMessage } from 'element-plus'
import { onMounted, ref } from 'vue'

const loading = ref(false)
const activeTab = ref('aliyun')
const configs = ref({
  active: 'aliyun',
  aliyun: {
    access_key: '',
    access_secret: '',
    sign: '',
    code_temp_id: '',
  },
  bao: {
    username: '',
    password: '',
    sign: '',
    code_template: '',
  },
  tencent: {
    secret_id: '',
    secret_key: '',
    sms_sdk_app_id: '',
    sign: '',
    code_temp_id: '',
    code_template: '',
    region: 'ap-guangzhou',
  },
})

onMounted(async () => {
  try {
    const res = await httpGet('/api/admin/config/get?key=sms')
    configs.value = Object.assign(configs.value, res.data)
  } catch (e) {
    // 使用默认值
  }
})

const saveSmsConfig = async () => {
  loading.value = true
  try {
    await httpPost('/api/admin/config/update/sms', configs.value)
    ElMessage.success('保存成功')
  } catch (e) {
    ElMessage.error('保存失败：' + (e.message || '未知错误'))
  } finally {
    loading.value = false
  }
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
