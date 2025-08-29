<template>
  <div class="container" v-loading="loading">
    <md-editor
      :theme="store.theme"
      v-model="content"
      @on-upload-img="onUploadImg"
      placeholder="请输入思维导图页面的默认文本内容，支持 Markdown 格式"
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
import { onMounted, reactive, ref } from 'vue'

const content = ref('')
const loading = ref(true)
const store = useSharedStore()

onMounted(() => {
  // 加载系统配置
  httpGet('/api/admin/config/get?key=mark_map')
    .then((res) => {
      content.value = res.data?.content || ''
    })
    .catch((e) => {
      ElMessage.error('加载系统配置失败: ' + e.message)
    })
    .finally(() => {
      loading.value = false
    })
})

const rules = reactive({})

const save = function () {
  httpPost('/api/admin/config/update/mark_map', {
    content: content.value,
  })
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
        // 执行上传操作
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

<style lang="scss" scoped>
@use '@/assets/css/admin/form.scss' as *;
@use '@/assets/css/main.scss' as *;

.markmap-config {
  display: flex;
  justify-content: center;
  padding: 20px;
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
