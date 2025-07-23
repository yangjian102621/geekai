<template>
  <div class="image-upload">
    <!-- 单图模式 -->
    <template v-if="props.maxCount === 1">
      <div class="single-upload">
        <div v-if="imageList.length === 0" class="upload-btn">
          <el-upload
            drag
            :auto-upload="true"
            :show-file-list="false"
            :http-request="handleUpload"
            :multiple="false"
            accept="image/*"
            class="uploader"
          >
            <div class="upload-placeholder">
              <el-icon :size="20"><UploadFilled /></el-icon>
              <span>上传图片</span>
            </div>
          </el-upload>
        </div>
        <div v-else class="upload-item single-image-item">
          <el-image :src="imageList[0]" fit="cover" class="upload-image" />
          <div class="upload-overlay" style="opacity: 1">
            <el-button
              type="danger"
              :icon="Delete"
              size="small"
              circle
              @click="removeImage(0)"
              class="remove-btn"
            />
          </div>
        </div>
      </div>
    </template>

    <!-- 多图模式 -->
    <template v-else>
      <div class="upload-list" v-if="imageList.length > 0">
        <div v-for="(image, index) in imageList" :key="index" class="upload-item">
          <el-image :src="image" fit="cover" class="upload-image" />
          <div class="upload-overlay">
            <el-button
              type="danger"
              :icon="Delete"
              size="small"
              circle
              @click="removeImage(index)"
              class="remove-btn"
            />
          </div>
        </div>
        <!-- 上传按钮 -->
        <div v-if="!multiple || imageList.length < maxCount" class="upload-btn">
          <el-upload
            drag
            :auto-upload="true"
            :show-file-list="false"
            :http-request="handleUpload"
            :multiple="multiple"
            accept="image/*"
            class="uploader"
            :limit="maxCount"
          >
            <div class="upload-placeholder">
              <el-icon :size="20"><UploadFilled /></el-icon>
              <span>上传图片</span>
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
          accept="image/*"
          class="uploader"
          :limit="maxCount"
        >
          <el-icon :size="40" class="el-icon--upload"><UploadFilled /></el-icon>
          <div class="el-upload__text">拖拽图片到此处，或 <em>点击上传</em></div>
          <template #tip>
            <div class="el-upload__tip text-center">
              支持 JPG、PNG 格式，最多上传 {{ maxCount }} 张，单张最大 5MB
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
import { httpPost } from '@/utils/http'
import { replaceImg } from '@/utils/libs'
import { Delete, UploadFilled } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { computed, ref } from 'vue'

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
})

const emit = defineEmits(['update:modelValue', 'upload-success'])

// 上传状态
const uploading = ref(false)
const uploadProgress = ref(0)

// 图片列表
const imageList = computed({
  get() {
    if (props.multiple || props.maxCount > 1) {
      return Array.isArray(props.modelValue) ? props.modelValue : []
    } else {
      return props.modelValue ? [props.modelValue] : []
    }
  },
  set(value) {
    if (props.multiple || props.maxCount > 1) {
      emit('update:modelValue', value)
    } else {
      emit('update:modelValue', value[0] || '')
    }
  },
})

const uploadCount = ref(1)
// 处理上传
const handleUpload = async (uploadFile) => {
  const file = uploadFile.file

  // 检查文件类型
  if (!file.type.startsWith('image/')) {
    ElMessage.error('请选择图片文件')
    return
  }

  // 检查文件大小 (5MB)
  if (file.size > 5 * 1024 * 1024) {
    ElMessage.error('图片大小不能超过 5MB')
    return
  }

  // 检查数量限制
  if (uploadCount.value > props.maxCount) {
    ElMessage.error(`最多只能上传 ${props.maxCount} 张图片`)
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

    const imageUrl = replaceImg(response.data.url)

    // 更新图片列表
    if (props.multiple || props.maxCount > 1) {
      const newList = [...imageList.value, imageUrl]
      imageList.value = newList
    } else {
      imageList.value = [imageUrl]
    }

    emit('upload-success', imageUrl)
    ElMessage.success('上传成功')
  } catch (error) {
    ElMessage.error('上传失败: ' + (error.message || '网络错误'))
  } finally {
    uploading.value = false
    uploadProgress.value = 0
  }
}

// 移除图片
const removeImage = (index) => {
  const newList = [...imageList.value]
  newList.splice(index, 1)
  imageList.value = newList
  uploadCount.value--
}
</script>

<style lang="stylus">
.image-upload {
  width: 100%;
}

.single-upload {
  width: 100px;
  height: 100px;
  position: relative;
}

.single-image-item {
  width: 100px;
  height: 100px;
  position: relative;
  border-radius: 6px;
  overflow: hidden;
  border: 1px solid #dcdfe6;
}

.upload-list {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

.upload-item {
  position: relative;
  width: 100px;
  height: 100px;
  border-radius: 6px;
  overflow: hidden;
  border: 1px solid #dcdfe6;

  .upload-image {
    width: 100%;
    height: 100%;
  }

  .upload-overlay {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.5);
    display: flex;
    align-items: center;
    justify-content: center;
    opacity: 0;
    transition: opacity 0.3s;

    .remove-btn {
      background: rgba(245, 108, 108, 0.8);
      border: none;
      color: white;
    }
  }

  &:hover .upload-overlay {
    opacity: 1;
  }
}

.upload-btn {
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
</style>
