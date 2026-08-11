<template>
  <div class="image__upload-container">
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
              <el-button size="small" @click.stop="openPasteZone" v-if="!showPasteZone"
                >粘贴截图</el-button
              >
            </div>
          </el-upload>
        </div>
        <div v-else class="upload-item single-image-item">
          <el-image :src="imageList[0]" fit="cover" class="upload-image" />
          <div class="upload-overlay flex items-center justify-center space-x-2">
            <el-tooltip content="删除" placement="top">
              <i
                class="iconfont icon-remove text-base text-red-500 cursor-pointer"
                @click="removeImage(0)"
              ></i>
            </el-tooltip>
            <el-tooltip content="预览" placement="top">
              <i
                class="iconfont icon-eye-open text-lg text-white cursor-pointer"
                @click="previewImage(0)"
              ></i>
            </el-tooltip>
          </div>
        </div>
      </div>
    </template>

    <!-- 多图模式 -->
    <template v-else>
      <div class="upload-list" v-if="imageList.length > 0">
        <div ref="uploadListRef" class="upload-list-inner">
          <div v-for="(image, index) in imageList" :key="image" class="upload-item">
            <el-image :src="image" fit="cover" class="upload-image" />
            <div class="upload-overlay flex items-center justify-center space-x-2">
              <el-tooltip content="删除" placement="top">
                <i
                  class="iconfont icon-remove text-base text-red-500 cursor-pointer"
                  @click="removeImage(index)"
                ></i>
              </el-tooltip>
              <el-tooltip content="预览" placement="top">
                <i
                  class="iconfont icon-eye-open text-lg text-white cursor-pointer"
                  @click="previewImage(index)"
                ></i>
              </el-tooltip>
            </div>
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
            :accept="accept"
            class="uploader"
            :limit="maxCount"
          >
            <div class="upload-placeholder">
              <i class="iconfont icon-plus"></i>
            </div>
          </el-upload>
          <div class="paste-actions">
            <el-button size="small" @click="openPasteZone" v-if="!showPasteZone"
              >粘贴截图</el-button
            >
          </div>
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
          <div class="el-upload__text">拖拽图片到此处，或 <em>点击上传</em></div>
          <template #tip>
            <div class="el-upload__tip text-gray-500 text-sm">
              支持 {{ accept }} 格式，最多上传 {{ maxCount }} 张，单张最大 {{ maxSize }}MB
            </div>
          </template>
        </el-upload>
        <div class="paste-actions">
          <el-button size="small" @click="openPasteZone" v-if="!showPasteZone">粘贴截图</el-button>
        </div>
      </div>
    </template>

    <!-- 粘贴区域（多图模式与单图模式共用） -->
    <div
      v-if="showPasteZone && (props.multiple || props.maxCount > 1 || imageList.length === 0)"
      ref="pasteZoneRef"
      class="paste-zone paste-zone-global"
      tabindex="0"
      @paste.prevent="onPaste"
    >
      <div class="paste-zone-text">在此区域按 Ctrl+V 粘贴截图</div>
      <el-button text size="small" class="paste-zone-close" @click="closePasteZone">关闭</el-button>
    </div>

    <!-- 上传进度 -->
    <el-progress
      v-if="uploading"
      :percentage="uploadProgress"
      :stroke-width="4"
      class="upload-progress"
    />

    <!-- 图片预览弹窗 -->
    <el-image-viewer
      v-if="previewVisible"
      :url-list="[previewImageSrc]"
      @close="previewVisible = false"
    />
  </div>
</template>

<script setup>
import { httpPost } from '@/utils/http'
import { replaceImg } from '@/utils/libs'
import { UploadFilled } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { Sortable } from 'sortablejs'
import { computed, nextTick, ref, watch, onBeforeUnmount } from 'vue'

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
    default: 5,
  },
  accept: {
    type: String,
    default: '.png,.jpg,.jpeg',
  },
})

const emit = defineEmits(['update:modelValue', 'upload-success'])

// 上传状态
const uploading = ref(false)
const uploadProgress = ref(0)
const previewVisible = ref(false)
const previewImageSrc = ref('')

// 粘贴截图区域
const showPasteZone = ref(false)
const pasteZoneRef = ref(null)

// 拖拽排序
const uploadListRef = ref(null)
let sortableInstance = null

// 图片列表
const imageList = computed({
  get() {
    if (props.multiple || props.maxCount > 1) {
      return Array.isArray(props.modelValue) ? props.modelValue : []
    } else {
      if (Array.isArray(props.modelValue)) {
        return props.modelValue.length > 0 && props.modelValue[0] ? [props.modelValue[0]] : []
      }
      return props.modelValue && props.modelValue.length > 0 ? [props.modelValue] : []
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

// 使用已选图片数量进行限制，不再使用全局计数
// 处理上传
const handleUpload = async (uploadFile) => {
  const file = uploadFile.file

  // 检查文件类型
  if (!file.type.startsWith('image/')) {
    ElMessage.error('请选择图片文件')
    return
  }

  // 检查文件大小 (5MB)
  if (file.size > props.maxSize * 1024 * 1024) {
    ElMessage.error(`图片大小不能超过 ${props.maxSize}MB`)
    return
  }

  // 检查数量限制（单图或多图）
  if ((props.multiple || props.maxCount > 1) && imageList.value.length >= props.maxCount) {
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
}

const previewImage = (index) => {
  previewImageSrc.value = imageList.value[index]
  previewVisible.value = true
}

const openPasteZone = () => {
  showPasteZone.value = true
  nextTick(() => {
    if (pasteZoneRef.value) {
      pasteZoneRef.value.focus()
    }
  })
}

const closePasteZone = () => {
  showPasteZone.value = false
}

const onPaste = (e) => {
  const items = e.clipboardData && e.clipboardData.items
  if (!items || !items.length) return

  for (let i = 0; i < items.length; i++) {
    const item = items[i]
    if (item.kind === 'file' && item.type && item.type.startsWith('image/')) {
      const file = item.getAsFile()
      if (file) {
        handleUpload({ file })
      }
      break
    }
  }
}

// 拖拽排序：同步新顺序到 v-model
const handleSortEnd = (evt) => {
  const { oldIndex, newIndex } = evt
  if (oldIndex === newIndex) return
  const list = [...imageList.value]
  const [removed] = list.splice(oldIndex, 1)
  list.splice(newIndex, 0, removed)
  imageList.value = list
  // 调试：拖拽后打印图片列表
  console.log('[ImageUpload] 拖拽排序', { oldIndex, newIndex, imageList: list })
}

const initSortable = () => {
  if (!uploadListRef.value || sortableInstance) return
  sortableInstance = Sortable.create(uploadListRef.value, {
    animation: 150,
    ghostClass: 'upload-item--sort-ghost',
    onEnd: handleSortEnd,
  })
}

const destroySortable = () => {
  if (sortableInstance) {
    sortableInstance.destroy()
    sortableInstance = null
  }
}

watch(
  () => [props.maxCount, imageList.value.length],
  () => {
    if (props.maxCount > 1 && imageList.value.length > 0) {
      nextTick(() => initSortable())
    } else {
      destroySortable()
    }
  },
  { immediate: true }
)

onBeforeUnmount(destroySortable)
</script>

<style lang="scss">
.image__upload-container {
  --upload-size: 100px;
  width: 100%;

  .el-upload-dragger {
    --el-upload-dragger-padding-horizontal: 20px;
  }

  .single-upload {
    width: var(--upload-size);
    height: var(--upload-size);
    position: relative;
  }

  .single-image-item {
    width: var(--upload-size);
    height: var(--upload-size);
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

  .upload-list-inner {
    display: flex;
    flex-wrap: wrap;
    gap: 10px;
  }

  .upload-item {
    position: relative;
    width: var(--upload-size);
    height: var(--upload-size);
    border-radius: 6px;
    overflow: hidden;
    border: 1px solid #dcdfe6;
    cursor: move;

    &.upload-item--sort-ghost {
      opacity: 0.5;
    }

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
        width: var(--upload-size);
        height: var(--upload-size);
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

  .paste-actions {
    margin-top: 8px;
  }

  .paste-zone {
    margin-top: 8px;
    padding: 8px 10px;
    border: 1px dashed #c0c4cc;
    border-radius: 6px;
    min-height: 40px;
    display: flex;
    align-items: center;
    justify-content: space-between;
    background-color: #f5f7fa;
    outline: none;

    .paste-zone-text {
      font-size: 12px;
      color: #606266;
    }

    .paste-zone-close {
      margin-left: 8px;
    }
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
