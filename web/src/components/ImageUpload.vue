<template>
  <div class="image-upload">
    <div class="upload-list" v-if="imageList.length > 0">
      <div v-for="(image, index) in imageList" :key="index" class="upload-item">
        <el-image
          :src="image"
          :preview-src-list="imageList"
          :initial-index="index"
          fit="cover"
          class="upload-image"
        />
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
          :auto-upload="true"
          :show-file-list="false"
          :http-request="handleUpload"
          accept="image/*"
          class="uploader"
        >
          <div class="upload-placeholder">
            <el-icon :size="20"><Plus /></el-icon>
            <span>上传图片</span>
          </div>
        </el-upload>
      </div>
    </div>

    <!-- 初始上传区域 -->
    <div v-else class="upload-area">
      <el-upload
        :auto-upload="true"
        :show-file-list="false"
        :http-request="handleUpload"
        accept="image/*"
        class="uploader"
      >
        <div class="upload-placeholder">
          <el-icon :size="40"><Plus /></el-icon>
          <div class="upload-text">
            <p>点击上传图片</p>
            <p class="upload-tip">支持 JPG、PNG 格式，最大 10MB</p>
          </div>
        </div>
      </el-upload>
    </div>

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
import { Delete, Plus } from '@element-plus/icons-vue'
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
    if (props.multiple) {
      return Array.isArray(props.modelValue) ? props.modelValue : []
    } else {
      return props.modelValue ? [props.modelValue] : []
    }
  },
  set(value) {
    if (props.multiple) {
      emit('update:modelValue', value)
    } else {
      emit('update:modelValue', value[0] || '')
    }
  },
})

// 处理上传
const handleUpload = async (uploadFile) => {
  const file = uploadFile.file

  // 检查文件类型
  if (!file.type.startsWith('image/')) {
    ElMessage.error('请选择图片文件')
    return
  }

  // 检查文件大小 (10MB)
  if (file.size > 10 * 1024 * 1024) {
    ElMessage.error('图片大小不能超过 10MB')
    return
  }

  // 检查数量限制
  if (props.multiple && imageList.value.length >= props.maxCount) {
    ElMessage.error(`最多只能上传 ${props.maxCount} 张图片`)
    return
  }

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

    const imageUrl = response.data.url

    // 更新图片列表
    if (props.multiple) {
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
}
</script>

<style lang="stylus" scoped>
.image-upload
  width 100%

.upload-list
  display flex
  flex-wrap wrap
  gap 10px

.upload-item
  position relative
  width 100px
  height 100px
  border-radius 6px
  overflow hidden
  border 1px solid #dcdfe6

  .upload-image
    width 100%
    height 100%

  .upload-overlay
    position absolute
    top 0
    left 0
    right 0
    bottom 0
    background rgba(0, 0, 0, 0.5)
    display flex
    align-items center
    justify-content center
    opacity 0
    transition opacity 0.3s

    .remove-btn
      background rgba(245, 108, 108, 0.8)
      border none
      color white

  &:hover .upload-overlay
    opacity 1

.upload-btn
  width 100px
  height 100px
  border 2px dashed #dcdfe6
  border-radius 6px
  display flex
  align-items center
  justify-content center
  cursor pointer
  transition all 0.3s

  &:hover
    border-color #409eff
    color #409eff

  .uploader
    width 100%
    height 100%

  .upload-placeholder
    display flex
    flex-direction column
    align-items center
    gap 5px
    font-size 12px
    color #8c939d

.upload-area
  border 2px dashed #dcdfe6
  border-radius 6px
  padding 40px
  text-align center
  cursor pointer
  transition all 0.3s

  &:hover
    border-color #409eff

  .uploader
    width 100%

  .upload-placeholder
    display flex
    flex-direction column
    align-items center
    gap 10px
    color #8c939d

  .upload-text
    p
      margin 5px 0

    .upload-tip
      font-size 12px
      color #c0c4cc

.upload-progress
  margin-top 10px

:deep(.el-upload)
  width 100%
  height 100%
  display flex
  align-items center
  justify-content center
</style>
