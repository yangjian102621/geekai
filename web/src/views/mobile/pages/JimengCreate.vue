<template>
  <div class="mobile-jimeng-create">
    <!-- 页面头部 -->
    <div class="page-header">
      <van-nav-bar title="即梦AI" left-arrow @click-left="goBack" fixed placeholder />
    </div>

    <!-- 功能分类选择 -->
    <div class="category-section">
      <van-tabs v-model="activeCategory" @change="onCategoryChange">
        <van-tab title="图像生成" name="image_generation">
          <div class="tab-content">
            <!-- 生成模式切换 -->
            <van-cell title="生成模式">
              <template #value>
                <van-switch v-model="useImageInput" size="24" @change="onInputModeChange" />
              </template>
            </van-cell>
            <van-cell title="图生图人像写真" :value="useImageInput ? '开启' : '关闭'" />

            <!-- 文生图 -->
            <div v-if="activeFunction === 'text_to_image'" class="function-panel">
              <van-field
                v-model="currentPrompt"
                label="提示词"
                type="textarea"
                placeholder="请输入图片描述，越详细越好"
                rows="4"
                maxlength="2000"
                show-word-limit
              />

              <van-field
                v-model="textToImageParams.size"
                label="图片尺寸"
                readonly
                is-link
                @click="showSizePicker = true"
              />

              <van-cell title="创意度">
                <template #value>
                  <van-slider
                    v-model="textToImageParams.scale"
                    :min="1"
                    :max="10"
                    :step="0.5"
                    style="width: 200px"
                  />
                </template>
              </van-cell>

              <van-cell title="智能优化提示词">
                <template #right-icon>
                  <van-switch v-model="textToImageParams.use_pre_llm" size="24" />
                </template>
              </van-cell>
            </div>

            <!-- 图生图 -->
            <div v-if="activeFunction === 'image_to_image'" class="function-panel">
              <van-uploader
                v-model="imageToImageParams.image_input"
                :max-count="1"
                :after-read="onImageUpload"
                accept=".jpg,.png,.jpeg"
              >
                <van-button icon="plus" type="primary" block> 上传图片 </van-button>
              </van-uploader>

              <van-field
                v-model="currentPrompt"
                label="提示词"
                type="textarea"
                placeholder="描述你想要的图片效果"
                rows="4"
                maxlength="2000"
                show-word-limit
              />

              <van-field
                v-model="imageToImageParams.size"
                label="图片尺寸"
                readonly
                is-link
                @click="showSizePicker = true"
              />
            </div>
          </div>
        </van-tab>

        <van-tab title="图像编辑" name="image_editing">
          <div class="tab-content">
            <!-- 图像编辑 -->
            <div v-if="activeFunction === 'image_edit'" class="function-panel">
              <van-uploader
                v-model="imageEditParams.image_urls"
                :max-count="1"
                :after-read="onImageUpload"
                accept=".jpg,.png,.jpeg"
              >
                <van-button icon="plus" type="primary" block> 上传图片 </van-button>
              </van-uploader>

              <van-field
                v-model="currentPrompt"
                label="编辑提示词"
                type="textarea"
                placeholder="描述你想要的编辑效果"
                rows="4"
                maxlength="2000"
                show-word-limit
              />

              <van-cell title="编辑强度">
                <template #value>
                  <van-slider
                    v-model="imageEditParams.scale"
                    :min="0"
                    :max="1"
                    :step="0.1"
                    style="width: 200px"
                  />
                </template>
              </van-cell>
            </div>

            <!-- 图像特效 -->
            <div v-if="activeFunction === 'image_effects'" class="function-panel">
              <van-uploader
                v-model="imageEffectsParams.image_input1"
                :max-count="1"
                :after-read="onImageUpload"
                accept=".jpg,.png,.jpeg"
              >
                <van-button icon="plus" type="primary" block> 上传图片 </van-button>
              </van-uploader>

              <van-field
                v-model="imageEffectsParams.template_id"
                label="特效模板"
                readonly
                is-link
                @click="showTemplatePicker = true"
              />

              <van-field
                v-model="imageEffectsParams.size"
                label="输出尺寸"
                readonly
                is-link
                @click="showSizePicker = true"
              />
            </div>
          </div>
        </van-tab>

        <van-tab title="视频生成" name="video_generation">
          <div class="tab-content">
            <!-- 生成模式切换 -->
            <van-cell title="生成模式">
              <template #value>
                <van-switch v-model="useImageInput" size="24" @change="onInputModeChange" />
              </template>
            </van-cell>
            <van-cell title="图生视频" :value="useImageInput ? '开启' : '关闭'" />

            <!-- 文生视频 -->
            <div v-if="activeFunction === 'text_to_video'" class="function-panel">
              <van-field
                v-model="currentPrompt"
                label="提示词"
                type="textarea"
                placeholder="描述你想要的视频内容"
                rows="4"
                maxlength="2000"
                show-word-limit
              />

              <van-field
                v-model="textToVideoParams.aspect_ratio"
                label="视频比例"
                readonly
                is-link
                @click="showAspectRatioPicker = true"
              />
            </div>

            <!-- 图生视频 -->
            <div v-if="activeFunction === 'image_to_video'" class="function-panel">
              <van-uploader
                v-model="imageToVideoParams.image_urls"
                :max-count="2"
                :after-read="onImageUpload"
                accept=".jpg,.png,.jpeg"
                multiple
              >
                <van-button icon="plus" type="primary" block> 上传图片 </van-button>
              </van-uploader>

              <van-field
                v-model="currentPrompt"
                label="提示词"
                type="textarea"
                placeholder="描述你想要的视频效果"
                rows="4"
                maxlength="2000"
                show-word-limit
              />

              <van-field
                v-model="imageToVideoParams.aspect_ratio"
                label="视频比例"
                readonly
                is-link
                @click="showAspectRatioPicker = true"
              />
            </div>
          </div>
        </van-tab>
      </van-tabs>
    </div>

    <!-- 生成按钮 -->
    <div class="submit-section">
      <van-button type="primary" size="large" @click="submitTask" :loading="submitting" block>
        立即生成 ({{ currentPowerCost }}算力)
      </van-button>
    </div>

    <!-- 作品列表 -->
    <div class="works-list">
      <van-list
        v-model:loading="listLoading"
        :finished="listFinished"
        finished-text="没有更多了"
        @load="loadMore"
      >
        <div v-for="item in currentList" :key="item.id" class="work-item">
          <van-card
            :title="getFunctionName(item.type)"
            :desc="item.prompt"
            :thumb="item.img_url || item.video_url"
          >
            <template #tags>
              <van-tag :type="getTaskType(item.type)" size="small">
                {{ getFunctionName(item.type) }}
              </van-tag>
              <van-tag v-if="item.power" type="warning" size="small">
                {{ item.power }}算力
              </van-tag>
            </template>
            <template #footer>
              <van-button v-if="item.status === 'completed'" size="small" @click="playMedia(item)">
                {{ item.type.includes('video') ? '播放' : '查看' }}
              </van-button>
              <van-button
                v-if="item.status === 'completed'"
                size="small"
                @click="downloadFile(item)"
                :loading="item.downloading"
              >
                下载
              </van-button>
              <van-button v-if="item.status === 'failed'" size="small" @click="retryTask(item.id)">
                重试
              </van-button>
              <van-button size="small" type="danger" @click="removeJob(item)"> 删除 </van-button>
            </template>
          </van-card>
        </div>
      </van-list>
    </div>

    <!-- 各种选择器弹窗 -->
    <van-popup v-model:show="showSizePicker" position="bottom">
      <van-picker
        :columns="imageSizeOptions"
        @confirm="onSizeConfirm"
        @cancel="showSizePicker = false"
      />
    </van-popup>

    <van-popup v-model:show="showAspectRatioPicker" position="bottom">
      <van-picker
        :columns="videoAspectRatioOptions"
        @confirm="onAspectRatioConfirm"
        @cancel="showAspectRatioPicker = false"
      />
    </van-popup>

    <van-popup v-model:show="showTemplatePicker" position="bottom">
      <van-picker
        :columns="imageEffectsTemplateOptions"
        @confirm="onTemplateConfirm"
        @cancel="showTemplatePicker = false"
      />
    </van-popup>

    <!-- 媒体预览弹窗 -->
    <van-popup
      v-model:show="showMediaDialog"
      position="center"
      :style="{ width: '90%', height: '60%' }"
    >
      <div class="media-preview">
        <img
          v-if="currentMediaUrl && !currentMediaUrl.includes('video')"
          :src="currentMediaUrl"
          style="width: 100%; height: 100%; object-fit: contain"
        />
        <video
          v-else-if="currentMediaUrl"
          :src="currentMediaUrl"
          controls
          autoplay
          style="width: 100%; height: 100%"
        >
          您的浏览器不支持视频播放
        </video>
      </div>
    </van-popup>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { showToast, showDialog } from 'vant'
import { httpGet, httpPost } from '@/utils/http'
import { checkSession } from '@/store/cache'

const router = useRouter()

// 响应式数据
const activeCategory = ref('image_generation')
const useImageInput = ref(false)
const submitting = ref(false)
const listLoading = ref(false)
const listFinished = ref(false)
const currentList = ref([])
const showMediaDialog = ref(false)
const currentMediaUrl = ref('')

// 选择器相关
const showSizePicker = ref(false)
const showAspectRatioPicker = ref(false)
const showTemplatePicker = ref(false)

// 当前提示词
const currentPrompt = ref('')

// 功能参数
const textToImageParams = ref({
  size: '1024x1024',
  scale: 7.5,
  use_pre_llm: false,
})

const imageToImageParams = ref({
  image_input: [],
  size: '1024x1024',
})

const imageEditParams = ref({
  image_urls: [],
  scale: 0.5,
})

const imageEffectsParams = ref({
  image_input1: [],
  template_id: '',
  size: '1024x1024',
})

const textToVideoParams = ref({
  aspect_ratio: '16:9',
})

const imageToVideoParams = ref({
  image_urls: [],
  aspect_ratio: '16:9',
})

// 选项数据
const imageSizeOptions = ['512x512', '768x768', '1024x1024', '1024x1536', '1536x1024']

const videoAspectRatioOptions = ['16:9', '9:16', '1:1', '4:3']

const imageEffectsTemplateOptions = [
  'acrylic_ornaments',
  'angel_figurine',
  'felt_3d_polaroid',
  'watercolor_illustration',
]

// 页面数据
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)
const currentPowerCost = ref(0)
const taskPulling = ref(true)
const tastPullHandler = ref(null)

// 计算属性
const activeFunction = computed(() => {
  if (activeCategory.value === 'image_generation') {
    return useImageInput.value ? 'image_to_image' : 'text_to_image'
  } else if (activeCategory.value === 'image_editing') {
    return 'image_edit' // 可以根据需要添加更多编辑功能
  } else if (activeCategory.value === 'video_generation') {
    return useImageInput.value ? 'image_to_video' : 'text_to_video'
  }
  return 'text_to_image'
})

// 生命周期
onMounted(() => {
  checkSession()
    .then(() => {
      fetchData(1)
      tastPullHandler.value = setInterval(() => {
        if (taskPulling.value) {
          fetchData(1)
        }
      }, 5000)
    })
    .catch(() => {})
})

onUnmounted(() => {
  if (tastPullHandler.value) {
    clearInterval(tastPullHandler.value)
  }
})

// 方法
const goBack = () => {
  router.back()
}

const onCategoryChange = (name) => {
  activeCategory.value = name
  useImageInput.value = false
}

const onInputModeChange = () => {
  // 重置相关参数
  currentPrompt.value = ''
}

const onImageUpload = (file) => {
  const formData = new FormData()
  formData.append('file', file.file, file.name)
  showToast({ message: '正在上传图片...', duration: 0 })

  httpPost('/api/upload', formData)
    .then((res) => {
      showToast('图片上传成功')
      return res.data.url
    })
    .catch((e) => {
      showToast('图片上传失败:' + e.message)
    })
}

const submitTask = () => {
  if (!currentPrompt.value.trim()) {
    showToast('请输入提示词')
    return
  }

  submitting.value = true
  const params = {
    type: activeFunction.value,
    prompt: currentPrompt.value,
  }

  // 根据功能类型添加相应参数
  if (activeFunction.value === 'text_to_image') {
    Object.assign(params, textToImageParams.value)
  } else if (activeFunction.value === 'image_to_image') {
    Object.assign(params, imageToImageParams.value)
  } else if (activeFunction.value === 'image_edit') {
    Object.assign(params, imageEditParams.value)
  } else if (activeFunction.value === 'image_effects') {
    Object.assign(params, imageEffectsParams.value)
  } else if (activeFunction.value === 'text_to_video') {
    Object.assign(params, textToVideoParams.value)
  } else if (activeFunction.value === 'image_to_video') {
    Object.assign(params, imageToVideoParams.value)
  }

  httpPost('/api/jimeng/create', params)
    .then(() => {
      fetchData(1)
      taskPulling.value = true
      showToast('创建任务成功')
      currentPrompt.value = ''
    })
    .catch((e) => {
      showToast('创建任务失败：' + e.message)
    })
    .finally(() => {
      submitting.value = false
    })
}

const fetchData = (_page) => {
  if (_page) {
    page.value = _page
  }
  listLoading.value = true
  httpGet('/api/jimeng/list', { page: page.value, page_size: pageSize.value })
    .then((res) => {
      total.value = res.data.total
      let needPull = false
      const items = []
      for (let v of res.data.items) {
        if (v.status === 'in_queue' || v.status === 'generating') {
          needPull = true
        }
        items.push(v)
      }
      listLoading.value = false
      taskPulling.value = needPull

      if (page.value === 1) {
        currentList.value = items
      } else {
        currentList.value.push(...items)
      }

      if (items.length < pageSize.value) {
        listFinished.value = true
      }
    })
    .catch((e) => {
      listLoading.value = false
      showToast('获取作品列表失败：' + e.message)
    })
}

const loadMore = () => {
  page.value++
  fetchData()
}

const playMedia = (item) => {
  currentMediaUrl.value = item.img_url || item.video_url
  showMediaDialog.value = true
}

const downloadFile = (item) => {
  item.downloading = true
  const link = document.createElement('a')
  link.href = item.img_url || item.video_url
  link.download = item.title || 'file'
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
  item.downloading = false
  showToast('开始下载')
}

const retryTask = (id) => {
  httpPost('/api/jimeng/retry', { id })
    .then(() => {
      showToast('重试任务成功')
      fetchData(1)
    })
    .catch((e) => {
      showToast('重试任务失败：' + e.message)
    })
}

const removeJob = (item) => {
  showDialog({
    title: '确认删除',
    message: '此操作将会删除任务相关文件，继续操作吗?',
    showCancelButton: true,
  })
    .then(() => {
      httpGet('/api/jimeng/remove', { id: item.id })
        .then(() => {
          showToast('任务删除成功')
          fetchData(1)
        })
        .catch((e) => {
          showToast('任务删除失败：' + e.message)
        })
    })
    .catch(() => {})
}

// 工具方法
const getFunctionName = (type) => {
  const nameMap = {
    text_to_image: '文生图',
    image_to_image: '图生图',
    image_edit: '图像编辑',
    image_effects: '图像特效',
    text_to_video: '文生视频',
    image_to_video: '图生视频',
  }
  return nameMap[type] || type
}

const getTaskType = (type) => {
  return type.includes('video') ? 'warning' : 'primary'
}

// 选择器确认方法
const onSizeConfirm = (value) => {
  if (activeFunction.value === 'text_to_image') {
    textToImageParams.value.size = value
  } else if (activeFunction.value === 'image_to_image') {
    imageToImageParams.value.size = value
  } else if (activeFunction.value === 'image_effects') {
    imageEffectsParams.value.size = value
  }
  showSizePicker.value = false
}

const onAspectRatioConfirm = (value) => {
  if (activeFunction.value === 'text_to_video') {
    textToVideoParams.value.aspect_ratio = value
  } else if (activeFunction.value === 'image_to_video') {
    imageToVideoParams.value.aspect_ratio = value
  }
  showAspectRatioPicker.value = false
}

const onTemplateConfirm = (value) => {
  imageEffectsParams.value.template_id = value
  showTemplatePicker.value = false
}
</script>

<style lang="scss" scoped>
.mobile-jimeng-create {
  min-height: 100vh;
  background: #f7f8fa;
  padding-bottom: 20px;

  .page-header {
    background: #fff;
  }

  .category-section {
    background: #fff;
    margin: 12px;
    border-radius: 8px;
    overflow: hidden;

    .tab-content {
      padding: 16px;

      .function-panel {
        .van-uploader {
          margin-bottom: 16px;
        }
      }
    }
  }

  .submit-section {
    margin: 12px;
    padding: 16px;
    background: #fff;
    border-radius: 8px;
  }

  .works-list {
    margin: 12px;

    .work-item {
      margin-bottom: 12px;
      background: #fff;
      border-radius: 8px;
      overflow: hidden;
    }
  }

  .media-preview {
    width: 100%;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
  }
}
</style>
