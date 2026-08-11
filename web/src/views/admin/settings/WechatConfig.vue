<template>
  <div class="system-config form">
    <div class="container">
      <el-tabs type="border-card" v-model="activeName">
        <el-tab-pane label="微信公众号配置" name="gzh">
          <el-form
            label-width="150px"
            label-position="top"
            :model="gzhConfig"
            ref="wechatFormRef"
            :rules="rules"
          >
            <el-form-item label="微信公众号 AppID" prop="app_id">
              <el-input v-model="gzhConfig['app_id']" />
            </el-form-item>
            <el-form-item label="微信公众号 AppSecret" prop="secret">
              <el-input v-model="gzhConfig['secret']" />
            </el-form-item>
            <el-form-item label="微信公众号 Token" prop="token">
              <el-input v-model="gzhConfig['token']" />
            </el-form-item>
            <el-form-item label="微信公众号 EncodingAESKey" prop="encoding_aes_key">
              <el-input v-model="gzhConfig['encoding_aes_key']" />
            </el-form-item>
            <el-form-item label="启用微信公众号登录" prop="enabled">
              <el-switch v-model="gzhConfig['enabled']" />
            </el-form-item>
          </el-form>
        </el-tab-pane>
      </el-tabs>

      <div style="padding: 10px">
        <el-button type="primary" @click="save">保存</el-button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { onMounted, reactive, ref } from 'vue'
import { httpGet, httpPost } from '@/utils/http'
import { ElMessage } from 'element-plus'

const activeName = ref('gzh')
const gzhConfig = ref({ app_id: '', secret: '' })
const wechatFormRef = ref(null)

onMounted(() => {
  // 加载系统配置
  httpGet('/api/admin/config/get?key=wx_gzh')
    .then((res) => {
      gzhConfig.value = res.data || {}
    })
    .catch((e) => {
      ElMessage.error('加载系统配置失败: ' + e.message)
    })
})

const rules = reactive({
  app_id: [{ required: true, message: '请输入公众号 AppID', trigger: 'blur' }],
  secret: [{ required: true, message: '请输入公众号 AppSecret', trigger: 'blur' }],
})
const save = function () {
  wechatFormRef.value.validate((valid) => {
    if (valid) {
      httpPost('/api/admin/config/update/wx_gzh', {
        app_id: gzhConfig.value.app_id,
        secret: gzhConfig.value.secret,
        token: gzhConfig.value.token,
        encoding_aes_key: gzhConfig.value.encoding_aes_key,
        enabled: gzhConfig.value.enabled,
      })
        .then(() => {
          ElMessage.success('操作成功！')
        })
        .catch((e) => {
          ElMessage.error('操作失败：' + e.message)
        })
    }
  })
}
</script>

<style lang="css" scoped>
@import '@/assets/css/admin/form.css';
@import '@/assets/css/main.css';

.system-config {
  display: flex;
  justify-content: center;
  .sys-tabs {
    width: 100%;
    background-color: var(--el-bg-color);
    padding: 10px 20px 40px 20px;
    /* border: 1px solid var(--el-border-color); */
  }
}
</style>
