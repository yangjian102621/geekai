<template>
  <div class="license-config form" v-loading="loading">
    <div class="container">
      <el-descriptions
        v-if="license.is_active"
        class="margin-top"
        title="已授权信息"
        :column="1"
        border
      >
        <el-descriptions-item>
          <template #label>
            <div class="cell-item">License Key</div>
          </template>
          {{ license.key }}
        </el-descriptions-item>
        <el-descriptions-item>
          <template #label>
            <div class="cell-item">机器码</div>
          </template>
          {{ license.machine_id }}
        </el-descriptions-item>
        <el-descriptions-item>
          <template #label>
            <div class="cell-item">到期时间</div>
          </template>
          {{ dateFormat(license.expired_at) }}
        </el-descriptions-item>
        <el-descriptions-item>
          <template #label>
            <div class="cell-item">用户人数</div>
          </template>
          {{ license.configs?.user_num }}
        </el-descriptions-item>
        <el-descriptions-item>
          <template #label>
            <div class="cell-item">去版权</div>
          </template>
          <el-icon class="selected" v-if="license.configs?.de_copy"><Select /></el-icon>
          <el-icon class="closed" v-else><CloseBold /></el-icon>
          <span class="text">去版权之后前端页面将不会显示版权信息和源码地址</span>
        </el-descriptions-item>
      </el-descriptions>

      <h3>激活后可获得以下权限：</h3>
      <ol class="active-info">
        <li>1、使用任意第三方中转 API KEY，而不用局限于 GeekAI 推荐的白名单列表</li>
        <li>2、可以在相关页面去除 GeekAI 的版权信息，或者修改为自己的版权信息</li>
      </ol>

      <el-form :model="system" label-width="150px" label-position="right">
        <el-form-item label="许可授权码" prop="license">
          <el-input v-model="licenseKey" />
        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="active">立即激活</el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script setup>
import { httpGet, httpPost } from '@/utils/http'
import { dateFormat } from '@/utils/libs'
import { ElMessage } from 'element-plus'
import { CloseBold, Select } from '@element-plus/icons-vue'
import { onMounted, ref } from 'vue'

const loading = ref(true)
const license = ref({ is_active: false })
const licenseKey = ref('')

onMounted(() => {
  fetchLicense()
})

const fetchLicense = () => {
  httpGet('/api/admin/config/license')
    .then((res) => {
      license.value = res.data
    })
    .catch((e) => {
      ElMessage.error('获取 License 失败：' + e.message)
    })
    .finally(() => {
      loading.value = false
    })
}

// 激活授权
const active = () => {
  if (licenseKey.value === '') {
    return ElMessage.error('请输入授权码')
  }
  httpPost('/api/admin/config/active', { license: licenseKey.value })
    .then((res) => {
      ElMessage.success('授权成功，机器编码为：' + res.data)
      fetchLicense()
    })
    .catch((e) => {
      ElMessage.error(e.message)
    })
}
</script>

<style scoped>
.license-config {
  display: flex;
  justify-content: center;
}

.container {
  width: 100%;
  background-color: var(--el-bg-color);
  padding: 10px 20px 40px 20px;
}

.margin-top {
  margin-top: 20px;
}

.cell-item {
  font-weight: bold;
}

.selected {
  color: #67c23a;
}

.closed {
  color: #f56c6c;
}

.text {
  margin-left: 10px;
}

.active-info {
  margin: 20px 0;
  padding-left: 20px;
}

.active-info li {
  margin: 10px 0;
  line-height: 1.6;
}
</style>
