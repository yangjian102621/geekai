<template>
  <div class="file__upload-container">
    <!-- 单文件模式 -->
    <template v-if="props.maxCount === 1">
      <div class="single-upload">
        <div v-if="fileList.length === 0" class="upload-btn upload-btn-single">
          <el-upload
            drag
            :auto-upload="true"
            :show-file-list="false"
            :http-request="handleUpload"
            :multiple="false"
            :accept="accept"
            class="uploader"
          >
            <div class="upload-placeholder">
              <el-icon :size="20"><UploadFilled /></el-icon>
              <span>上传文件</span>
            </div>
          </el-upload>
        </div>
        <div
          v-else
          class="relative inline-flex items-center border border-gray-200 rounded-xl bg-white dark:bg-[#2b2b2b] dark:border-gray-700 p-2 w-full"
        >
          <img :src="getFileImage(fileList[0].url)" class="w-10 h-10 mr-2" />
          <div class="min-w-0 flex flex-col items-center gap-1 text-sm">
            <a
              :href="fileList[0].url"
              target="_blank"
              class="truncate block text-[var(--theme-text-color-primary,#0d0d0d)] max-w-[220px]"
            >
              {{ fileList[0].name }}
            </a>
            <div class="text-xs flex w-full justify-start text-gray-500">
              {{ GetFileType(getFileExt(fileList[0].name)) }} ·
              {{ FormatFileSize(fileList[0].size || 0) }}
            </div>
          </div>
          <button
            class="absolute -right-2 -top-2 w-5 h-5 rounded-full bg-rose-600 text-white flex items-center justify-center text-[10px]"
            @click="removeFile(0)"
            aria-label="remove"
          >
            ×
          </button>
        </div>
      </div>
    </template>

    <!-- 多文件模式 -->
    <template v-else>
      <div class="flex flex-col gap-2 px-2 pt-2 !items-start" v-if="fileList.length > 0">
        <div
          v-for="(file, index) in fileList"
          :key="file.url || index"
          class="relative inline-flex items-center border border-gray-200 rounded-xl bg-white dark:bg-[#2b2b2b] dark:border-gray-700 p-2 w-full"
        >
          <img :src="getFileImage(file.url)" class="w-10 h-10 mr-2" />
          <div class="min-w-0 flex flex-col items-center gap-1 text-sm">
            <a :href="file.url" target="_blank" class="truncate block max-w-[180px]">{{
              file.name
            }}</a>
            <div class="text-xs flex w-full justify-start text-gray-500">
              {{ GetFileType(getFileExt(file.name)) }} · {{ FormatFileSize(file.size || 0) }}
            </div>
          </div>
          <button
            class="absolute -right-2 -top-2 w-5 h-5 rounded-full bg-rose-600 text-white flex items-center justify-center text-[10px]"
            @click="removeFile(index)"
            aria-label="remove"
          >
            ×
          </button>
        </div>
        <!-- 上传按钮 -->
        <div v-if="!multiple || fileList.length < maxCount" class="upload-btn">
          <el-upload
            drag
            :auto-upload="true"
            :show-file-list="false"
            :http-request="handleUpload"
            :multiple="multiple"
            :accept="accept"
            class="uploader"
            :limit="maxCount"
          >
            <div class="upload-placeholder">
              <el-icon :size="20"><UploadFilled /></el-icon>
              <span>上传文件</span>
            </div>
          </el-upload>
        </div>
      </div>
      <!-- 初始上传区域 -->
      <div v-else class="upload-area">
        <el-upload
          drag
          :auto-upload="true"
          :show-file-list="false"
          :http-request="handleUpload"
          :multiple="multiple"
          :accept="accept"
          class="uploader"
          :limit="maxCount"
        >
          <el-icon :size="40" class="el-icon--upload"><UploadFilled /></el-icon>
          <div class="el-upload__text">拖拽文件到此处，或 <em>点击上传</em></div>
          <template #tip>
            <div class="el-upload__tip text-gray-500 text-sm">
              支持 {{ accept }} 格式，最多上传 {{ maxCount }} 个，单个最大 {{ maxSize }}MB
            </div>
          </template>
        </el-upload>
      </div>
    </template>

    <!-- 上传进度 -->
    <el-progress
      v-if="uploading"
      :percentage="uploadProgress"
      :stroke-width="4"
      class="upload-progress"
    />
  </div>
</template>

<script setup>
import { FormatFileSize, GetFileIcon, GetFileType } from '@/store/system'
import { httpPost } from '@/utils/http'
import { isImage, replaceImg } from '@/utils/libs'
import { UploadFilled } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { computed, ref, watch } from 'vue'

const props = defineProps({
  modelValue: {
    type: [String, Array],
    default: '',
  },
  multiple: {
    type: Boolean,
    default: false,
  },
  maxCount: {
    type: Number,
    default: 1,
  },
  maxSize: {
    type: Number,
    default: 10,
  },
  accept: {
    type: String,
    default: '.pdf,.doc,.docx,.xls,.xlsx,.ppt,.pptx,.txt,.md,.zip,.rar,.7z',
  },
})

const emit = defineEmits(['update:modelValue', 'upload-success', 'remove-file'])

// 上传状态
const uploading = ref(false)
const uploadProgress = ref(0)
const FileInfoList = ref([])

// 文件列表
const fileList = computed({
  get() {
    if (props.multiple || props.maxCount > 1) {
      return FileInfoList.value
    } else {
      return FileInfoList.value && FileInfoList.value.length > 0 ? FileInfoList.value : []
    }
  },
  set(value) {
    const isMulti = props.multiple || props.maxCount > 1
    const normalized = Array.isArray(value) ? value : value ? [value] : []
    FileInfoList.value = normalized
    if (isMulti) {
      const urls = normalized.map((v) => v && v.url).filter((u) => !!u)
      emit('update:modelValue', urls)
    } else {
      const url =
        normalized.length > 0 && normalized[0] && normalized[0].url ? normalized[0].url : ''
      emit('update:modelValue', url)
    }
  },
})

const uploadCount = ref(1)

// 获取文件扩展名
const getFileExt = (filename) => {
  return '.' + filename.split('.').pop().toLowerCase()
}

const getFileName = (url) => {
  return url.split('/').pop()
}

// 获取文件
const getFileImage = (url) => {
  return isImage(url) ? url : GetFileIcon(getFileExt(url))
}

// 将外部 modelValue 同步为内部文件对象列表
const urlToFileInfo = (url) => ({
  url,
  name: getFileName(url),
  size: 0,
  ext: getFileExt(url),
})

// 通过 HEAD 请求尝试获取远程资源大小
const fetchRemoteFileSize = async (url) => {
  try {
    const res = await fetch(url, { method: 'HEAD' })
    const len = res.headers.get('content-length')
    return len ? parseInt(len, 10) : 0
  } catch (e) {
    return 0
  }
}

// 对 size 为空或 0 的项进行补充
const updateUnknownSizes = async (items) => {
  const tasks = items.map(async (it) => {
    if (!it || !it.url) return it
    if (!it.size || it.size === 0) {
      const s = await fetchRemoteFileSize(it.url)
      if (s > 0) {
        it.size = s
      }
    }
    return it
  })
  await Promise.all(tasks)
}

watch(
  () => props.modelValue,
  (newVal) => {
    const isMulti = props.multiple || props.maxCount > 1
    if (isMulti) {
      const urls = Array.isArray(newVal) ? newVal : []
      FileInfoList.value = urls.map((u) => urlToFileInfo(u))
      // 异步补齐大小
      updateUnknownSizes(FileInfoList.value)
    } else {
      const url = typeof newVal === 'string' ? newVal : ''
      FileInfoList.value = url ? [urlToFileInfo(url)] : []
      // 异步补齐大小
      updateUnknownSizes(FileInfoList.value)
    }
  },
  { immediate: true }
)

// 处理上传
const handleUpload = async (uploadFile) => {
  const file = uploadFile.file
  // 检查文件大小
  if (file.size > props.maxSize * 1024 * 1024) {
    ElMessage.error(`文件大小不能超过 ${props.maxSize}MB`)
    return
  }

  // 检查数量限制
  if (uploadCount.value > props.maxCount) {
    ElMessage.error(`最多只能上传 ${props.maxCount} 个文件`)
    return
  }
  uploadCount.value++

  uploading.value = true
  uploadProgress.value = 0

  try {
    const formData = new FormData()
    formData.append('file', file)

    // 模拟上传进度
    const progressTimer = setInterval(() => {
      if (uploadProgress.value < 90) {
        uploadProgress.value += 10
      }
    }, 100)

    const response = await httpPost('/api/upload', formData)

    clearInterval(progressTimer)
    uploadProgress.value = 100

    const fileUrl = replaceImg(response.data.url)
    const fileInfo = {
      name: file.name,
      url: fileUrl,
      size: file.size,
      ext: getFileExt(file.name),
    }

    // 更新文件列表
    if (props.multiple || props.maxCount > 1) {
      const newList = [...fileList.value, fileInfo]
      fileList.value = newList
    } else {
      fileList.value = [fileInfo]
    }

    emit('upload-success', fileInfo)
    ElMessage.success('上传成功')
  } catch (error) {
    ElMessage.error('上传失败: ' + (error.message || '网络错误'))
  } finally {
    uploading.value = false
    uploadProgress.value = 0
  }
}

// 移除文件
const removeFile = (index) => {
  const file = fileList.value[index]
  const newList = [...fileList.value]
  newList.splice(index, 1)
  fileList.value = newList
  uploadCount.value--
  emit('remove-file', file)
}
</script>

<style lang="scss">
.file__upload-container {
  width: 100%;

  .single-upload {
    width: 100%;
    position: relative;

    .upload-btn-single {
      .uploader {
        width: 100%;

        .el-upload-dragger {
          width: 100px;
          height: 100px;
          display: flex;
          align-items: center;
          justify-content: center;
        }
      }
    }
  }

  // .single-file-item {
  //   width: 100%;
  //   position: relative;
  //   border-radius: 6px;
  //   overflow: hidden;
  //   border: 1px solid #dcdfe6;
  //   background-color: var(--chat-content-bg, #f5f5f5);
  //   padding: 8px;
  // }

  .upload-list {
    display: flex;
    flex-wrap: wrap;
    gap: 10px;
  }

  .upload-item {
    position: relative;
    width: 100%;
    max-width: 300px;
    border-radius: 6px;
    overflow: hidden;
    border: 1px solid #dcdfe6;
    background-color: var(--chat-content-bg, #f5f5f5);
    padding: 8px;

    .file-info {
      display: flex;
      flex-direction: row;
      align-items: center;

      .icon {
        .el-image {
          width: 40px;
          height: 40px;
        }
      }

      .body {
        margin-left: 8px;
        flex: 1;
        font-size: 12px;

        .title {
          color: #0d0d0d;
          margin-bottom: 4px;
        }

        .info {
          color: #b4b4b4;

          span {
            margin-right: 10px;
          }
        }
      }
    }

    .upload-overlay {
      position: absolute;
      top: -8px;
      right: -8px;
      opacity: 1;

      .remove-btn {
        background: rgba(245, 108, 108, 0.8);
        border: none;
        color: white;
      }
    }
  }

  .upload-btn {
    .uploader {
      width: 100%;

      .el-upload-dragger {
        display: flex;
        align-items: center;
        justify-content: center;
      }
    }

    .upload-placeholder {
      display: flex;
      flex-direction: column;
      align-items: center;
      gap: 5px;
      font-size: 12px;
      color: #8c939d;
    }
  }

  .upload-area {
    .el-upload-dragger {
      width: 100%;
    }

    .uploader {
      width: 100%;
    }
  }

  .upload-progress {
    margin-top: 10px;
  }

  :deep(.el-upload) {
    width: 100%;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
  }
}
</style>
