<template>
  <div class="agreement-config container flex flex-col" v-loading="loading">
    <md-editor
      class="mgb20"
      v-model="agreement"
      :theme="store.theme"
      @on-upload-img="onUploadImg"
    />
    <div class="flex justify-center p-5">
      <el-button type="primary" @click="save">保存</el-button>
    </div>
  </div>
</template>

<script setup>
import { useSharedStore } from '@/store/sharedata'
import { httpGet, httpPost } from '@/utils/http'
import { ElMessage } from 'element-plus'
import MdEditor from 'md-editor-v3'
import 'md-editor-v3/lib/style.css'
import { onMounted, ref } from 'vue'

const store = useSharedStore()
const loading = ref(true)
const agreement = ref('')

onMounted(() => {
  httpGet('/api/admin/config/get?key=agreement')
    .then((res) => {
      agreement.value = res.data?.content || ''
    })
    .catch((e) => {
      console.warn('加载用户协议失败: ' + e.message)
    })
    .finally(() => {
      loading.value = false
    })
})

const save = () => {
  httpPost('/api/admin/config/update/agreement', { content: agreement.value })
    .then(() => {
      ElMessage.success('操作成功！')
    })
    .catch((e) => {
      ElMessage.error('操作失败：' + e.message)
    })
}

// 编辑期文件上传处理
const onUploadImg = (files, callback) => {
  Promise.all(
    files.map((file) => {
      return new Promise((rev, rej) => {
        const formData = new FormData()
        formData.append('file', file, file.name)
        httpPost('/api/admin/upload', formData)
          .then((res) => rev(res))
          .catch((error) => rej(error))
      })
    })
  )
    .then((res) => {
      ElMessage.success({ message: '上传成功', duration: 500 })
      callback(res.map((item) => item.data.url))
    })
    .catch((e) => {
      ElMessage.error('图片上传失败:' + e.message)
    })
}
</script>

<style scoped>
.agreement-config {
  display: flex;
  justify-content: center;
}

.container {
  width: 100%;
  background-color: var(--el-bg-color);
  padding: 10px 20px 40px 20px;
}

.mgb20 {
  margin-bottom: 20px;
}
</style>
