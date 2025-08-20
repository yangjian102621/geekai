<template>
  <div class="markmap-config form" v-loading="loading">
    <div class="container">
      <h3>思维导图配置</h3>
      <el-form
        :model="system"
        label-width="150px"
        label-position="right"
        ref="systemFormRef"
        :rules="rules"
      >
        <el-form-item>
          <template #label>
            <div class="label-title">
              思维导图默认文本
              <el-tooltip
                effect="dark"
                content="用户访问思维导图页面时显示的默认文本内容，支持 Markdown 格式"
                raw-content
                placement="right"
              >
                <el-icon>
                  <InfoFilled />
                </el-icon>
              </el-tooltip>
            </div>
          </template>
          <md-editor
            class="mgb20"
            :theme="store.theme"
            v-model="system['mark_map_text']"
            @on-upload-img="onUploadImg"
            placeholder="请输入思维导图页面的默认文本内容，支持 Markdown 格式"
          />
        </el-form-item>

        <div style="padding: 10px">
          <el-form-item>
            <el-button type="primary" @click="save">保存</el-button>
          </el-form-item>
        </div>
      </el-form>
    </div>
  </div>
</template>

<script setup>
import { useSharedStore } from '@/store/sharedata'
import { httpGet, httpPost } from '@/utils/http'
import { InfoFilled } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import MdEditor from 'md-editor-v3'
import 'md-editor-v3/lib/style.css'
import { onMounted, reactive, ref } from 'vue'

const system = ref({})
const loading = ref(true)
const systemFormRef = ref(null)
const store = useSharedStore()

onMounted(() => {
  // 加载系统配置
  httpGet('/api/admin/config/get?key=system')
    .then((res) => {
      system.value = res.data
      loading.value = false
    })
    .catch((e) => {
      ElMessage.error('加载系统配置失败: ' + e.message)
      loading.value = false
    })
})

const rules = reactive({})

const save = function () {
  httpPost('/api/admin/config/update', {
    key: 'system',
    config: { mark_map_text: system.value.mark_map_text },
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
@use '../../../assets/css/admin/form.scss' as *;
@use '../../../assets/css/main.scss' as *;

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
