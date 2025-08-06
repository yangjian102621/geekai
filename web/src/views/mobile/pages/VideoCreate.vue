<template>
  <div class="mobile-video-create">
    <!-- 页面头部 -->
    <div class="page-header">
      <van-nav-bar title="视频生成" left-arrow @click-left="goBack" fixed placeholder />
    </div>

    <!-- 视频类型切换 -->
    <div class="video-type-tabs">
      <van-tabs v-model="activeVideoType" @change="onVideoTypeChange">
        <van-tab title="Luma视频" name="luma">
          <div class="tab-content">
            <!-- Luma 视频参数 -->
            <div class="params-container">
              <!-- 提示词输入 -->
              <van-field
                v-model="lumaParams.prompt"
                label="提示词"
                type="textarea"
                placeholder="请在此输入视频提示词，用逗号分割"
                rows="4"
                maxlength="2000"
                show-word-limit
              />

              <!-- 提示词生成按钮 -->
              <van-button
                type="primary"
                size="small"
                @click="generatePrompt"
                :loading="isGenerating"
                block
                class="mt-2"
              >
                生成AI视频提示词
              </van-button>

              <!-- 图片辅助生成开关 -->
              <van-cell title="使用图片辅助生成">
                <template #right-icon>
                  <van-switch v-model="lumaUseImageMode" size="24" @change="toggleLumaImageMode" />
                </template>
              </van-cell>

              <!-- 图片上传区域 -->
              <div v-if="lumaUseImageMode" class="image-upload-section">
                <div class="image-upload-row">
                  <div class="image-upload-item">
                    <van-uploader
                      v-model="lumaStartImage"
                      :max-count="1"
                      :after-read="uploadLumaStartImage"
                      accept=".jpg,.png,.jpeg"
                    >
                      <div class="upload-placeholder">
                        <van-icon name="plus" />
                        <span>起始帧</span>
                      </div>
                    </van-uploader>
                  </div>
                  <div class="image-upload-item">
                    <van-uploader
                      v-model="lumaEndImage"
                      :max-count="1"
                      :after-read="uploadLumaEndImage"
                      accept=".jpg,.png,.jpeg"
                    >
                      <div class="upload-placeholder">
                        <van-icon name="plus" />
                        <span>结束帧</span>
                      </div>
                    </van-uploader>
                  </div>
                </div>
              </div>

              <!-- Luma 特有参数 -->
              <van-cell title="循环参考图">
                <template #right-icon>
                  <van-switch v-model="lumaParams.loop" size="24" />
                </template>
              </van-cell>

              <van-cell title="提示词优化">
                <template #right-icon>
                  <van-switch v-model="lumaParams.expand_prompt" size="24" />
                </template>
              </van-cell>

              <!-- 算力显示 -->
              <van-cell title="当前可用算力" :value="`${availablePower}`" />

              <!-- 生成按钮 -->
              <van-button
                type="primary"
                size="large"
                @click="createLumaVideo"
                :loading="generating"
                block
                class="mt-4"
              >
                立即生成 ({{ lumaPowerCost }}算力)
              </van-button>
            </div>
          </div>
        </van-tab>

        <van-tab title="可灵视频" name="keling">
          <div class="tab-content">
            <!-- KeLing 视频参数 -->
            <div class="params-container">
              <!-- 画面比例 -->
              <van-field
                v-model="kelingParams.aspect_ratio"
                label="画面比例"
                readonly
                is-link
                @click="showAspectRatioPicker = true"
              />

              <!-- 模型选择 -->
              <van-field
                v-model="kelingParams.model"
                label="模型选择"
                readonly
                is-link
                @click="showModelPicker = true"
              />

              <!-- 视频时长 -->
              <van-field
                v-model="kelingParams.duration"
                label="视频时长"
                readonly
                is-link
                @click="showDurationPicker = true"
              />

              <!-- 生成模式 -->
              <van-field
                v-model="kelingParams.mode"
                label="生成模式"
                readonly
                is-link
                @click="showModePicker = true"
              />

              <!-- 创意程度 -->
              <van-cell title="创意程度">
                <template #value>
                  <van-slider
                    v-model="kelingParams.cfg_scale"
                    :min="0"
                    :max="1"
                    :step="0.1"
                    style="width: 200px"
                  />
                </template>
              </van-cell>

              <!-- 运镜控制 -->
              <van-field
                v-model="kelingParams.camera_control.type"
                label="运镜控制"
                readonly
                is-link
                @click="showCameraControlPicker = true"
              />

              <!-- 图片辅助生成开关 -->
              <van-cell title="使用图片辅助生成">
                <template #right-icon>
                  <van-switch
                    v-model="kelingUseImageMode"
                    size="24"
                    @change="toggleKelingImageMode"
                  />
                </template>
              </van-cell>

              <!-- 图片上传区域 -->
              <div v-if="kelingUseImageMode" class="image-upload-section">
                <div class="image-upload-row">
                  <div class="image-upload-item">
                    <van-uploader
                      v-model="kelingStartImage"
                      :max-count="1"
                      :after-read="uploadKelingStartImage"
                      accept=".jpg,.png,.jpeg"
                    >
                      <div class="upload-placeholder">
                        <van-icon name="plus" />
                        <span>起始帧</span>
                      </div>
                    </van-uploader>
                  </div>
                  <div class="image-upload-item">
                    <van-uploader
                      v-model="kelingEndImage"
                      :max-count="1"
                      :after-read="uploadKelingEndImage"
                      accept=".jpg,.png,.jpeg"
                    >
                      <div class="upload-placeholder">
                        <van-icon name="plus" />
                        <span>结束帧</span>
                      </div>
                    </van-uploader>
                  </div>
                </div>
              </div>

              <!-- 提示词输入 -->
              <van-field
                v-model="kelingParams.prompt"
                label="提示词"
                type="textarea"
                :placeholder="kelingUseImageMode ? '描述视频画面细节' : '请在此输入视频提示词'"
                rows="4"
                maxlength="500"
                show-word-limit
              />

              <!-- 提示词生成按钮 -->
              <van-button
                type="primary"
                size="small"
                @click="generatePrompt"
                :loading="isGenerating"
                block
                class="mt-2"
              >
                生成专业视频提示词
              </van-button>

              <!-- 排除内容 -->
              <van-field
                v-model="kelingParams.negative_prompt"
                label="不希望出现的内容"
                type="textarea"
                placeholder="请在此输入你不希望出现在视频上的内容"
                rows="3"
                maxlength="500"
                show-word-limit
              />

              <!-- 算力显示 -->
              <van-cell title="当前可用算力" :value="`${availablePower}`" />

              <!-- 生成按钮 -->
              <van-button
                type="primary"
                size="large"
                @click="createKelingVideo"
                :loading="generating"
                block
                class="mt-4"
              >
                立即生成 ({{ kelingPowerCost }}算力)
              </van-button>
            </div>
          </div>
        </van-tab>
      </van-tabs>
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
          <van-card :title="item.title || '未命名视频'" :desc="item.prompt" :thumb="item.cover_url">
            <template #tags>
              <van-tag v-if="item.raw_data?.task_type" type="primary">
                {{ item.raw_data.task_type }}
              </van-tag>
              <van-tag v-if="item.raw_data?.model" type="success">
                {{ item.raw_data.model }}
              </van-tag>
              <van-tag v-if="item.raw_data?.duration" type="warning">
                {{ item.raw_data.duration }}秒
              </van-tag>
            </template>
            <template #footer>
              <van-button v-if="item.progress === 100" size="small" @click="playVideo(item)">
                播放
              </van-button>
              <van-button
                v-if="item.progress === 100"
                size="small"
                @click="downloadVideo(item)"
                :loading="item.downloading"
              >
                下载
              </van-button>
              <van-button size="small" type="danger" @click="removeJob(item)"> 删除 </van-button>
            </template>
          </van-card>
        </div>
      </van-list>
    </div>

    <!-- 各种选择器弹窗 -->
    <van-popup v-model:show="showAspectRatioPicker" position="bottom">
      <van-picker
        :columns="aspectRatioOptions"
        @confirm="onAspectRatioConfirm"
        @cancel="showAspectRatioPicker = false"
      />
    </van-popup>

    <van-popup v-model:show="showModelPicker" position="bottom">
      <van-picker
        :columns="modelOptions"
        @confirm="onModelConfirm"
        @cancel="showModelPicker = false"
      />
    </van-popup>

    <van-popup v-model:show="showDurationPicker" position="bottom">
      <van-picker
        :columns="durationOptions"
        @confirm="onDurationConfirm"
        @cancel="showDurationPicker = false"
      />
    </van-popup>

    <van-popup v-model:show="showModePicker" position="bottom">
      <van-picker
        :columns="modeOptions"
        @confirm="onModeConfirm"
        @cancel="showModePicker = false"
      />
    </van-popup>

    <van-popup v-model:show="showCameraControlPicker" position="bottom">
      <van-picker
        :columns="cameraControlOptions"
        @confirm="onCameraControlConfirm"
        @cancel="showCameraControlPicker = false"
      />
    </van-popup>

    <!-- 视频预览弹窗 -->
    <van-popup
      v-model:show="showVideoDialog"
      position="center"
      :style="{ width: '90%', height: '60%' }"
    >
      <div class="video-preview">
        <video
          v-if="currentVideoUrl"
          :src="currentVideoUrl"
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
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { showToast, showDialog } from 'vant'
import { httpGet, httpPost } from '@/utils/http'
import { checkSession } from '@/store/cache'

const router = useRouter()

// 响应式数据
const activeVideoType = ref('luma')
const loading = ref(false)
const generating = ref(false)
const isGenerating = ref(false)
const listLoading = ref(false)
const listFinished = ref(false)
const currentList = ref([])
const showVideoDialog = ref(false)
const currentVideoUrl = ref('')

// Luma 参数
const lumaParams = ref({
  prompt: '',
  image: '',
  image_tail: '',
  loop: false,
  expand_prompt: false,
})
const lumaUseImageMode = ref(false)
const lumaStartImage = ref([])
const lumaEndImage = ref([])

// KeLing 参数
const kelingParams = ref({
  aspect_ratio: '16:9',
  model: 'v1.5',
  duration: '5',
  mode: 'std',
  cfg_scale: 0.5,
  prompt: '',
  negative_prompt: '',
  image: '',
  image_tail: '',
  camera_control: {
    type: '',
    config: {
      horizontal: 0,
      vertical: 0,
      pan: 0,
      tilt: 0,
      roll: 0,
      zoom: 0,
    },
  },
})
const kelingUseImageMode = ref(false)
const kelingStartImage = ref([])
const kelingEndImage = ref([])

// 选择器相关
const showAspectRatioPicker = ref(false)
const showModelPicker = ref(false)
const showDurationPicker = ref(false)
const showModePicker = ref(false)
const showCameraControlPicker = ref(false)

// 选项数据
const aspectRatioOptions = ['16:9', '9:16', '1:1', '4:3']
const modelOptions = ['v1.0', 'v1.5']
const durationOptions = ['5', '10']
const modeOptions = ['std', 'pro']
const cameraControlOptions = [
  '',
  'simple',
  'down_back',
  'forward_up',
  'right_turn_forward',
  'left_turn_forward',
]

// 页面数据
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)
const availablePower = ref(0)
const lumaPowerCost = ref(0)
const kelingPowerCost = ref(0)
const taskPulling = ref(true)
const tastPullHandler = ref(null)

// 生命周期
onMounted(() => {
  checkSession()
    .then(() => {
      fetchData(1)
      fetchUserPower()
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

const onVideoTypeChange = (name) => {
  activeVideoType.value = name
}

const generatePrompt = () => {
  isGenerating.value = true
  // TODO: 实现提示词生成逻辑
  setTimeout(() => {
    isGenerating.value = false
    showToast('提示词生成功能开发中')
  }, 2000)
}

const toggleLumaImageMode = () => {
  if (!lumaUseImageMode.value) {
    lumaParams.value.image = ''
    lumaParams.value.image_tail = ''
    lumaStartImage.value = []
    lumaEndImage.value = []
  }
}

const toggleKelingImageMode = () => {
  if (!kelingUseImageMode.value) {
    kelingParams.value.image = ''
    kelingParams.value.image_tail = ''
    kelingStartImage.value = []
    kelingEndImage.value = []
  }
}

const uploadLumaStartImage = (file) => {
  uploadImage(file, (url) => {
    lumaParams.value.image = url
  })
}

const uploadLumaEndImage = (file) => {
  uploadImage(file, (url) => {
    lumaParams.value.image_tail = url
  })
}

const uploadKelingStartImage = (file) => {
  uploadImage(file, (url) => {
    kelingParams.value.image = url
  })
}

const uploadKelingEndImage = (file) => {
  uploadImage(file, (url) => {
    kelingParams.value.image_tail = url
  })
}

const uploadImage = (file, callback) => {
  const formData = new FormData()
  formData.append('file', file.file, file.name)
  showToast({ message: '正在上传图片...', duration: 0 })

  httpPost('/api/upload', formData)
    .then((res) => {
      callback(res.data.url)
      showToast('图片上传成功')
    })
    .catch((e) => {
      showToast('图片上传失败:' + e.message)
    })
}

const createLumaVideo = () => {
  if (!lumaParams.value.prompt.trim()) {
    showToast('请输入视频提示词')
    return
  }

  generating.value = true
  const params = {
    ...lumaParams.value,
    task_type: 'luma',
  }

  httpPost('/api/video/create', params)
    .then(() => {
      fetchData(1)
      taskPulling.value = true
      showToast('创建任务成功')
    })
    .catch((e) => {
      showToast('创建任务失败：' + e.message)
    })
    .finally(() => {
      generating.value = false
    })
}

const createKelingVideo = () => {
  if (!kelingParams.value.prompt.trim()) {
    showToast('请输入视频提示词')
    return
  }

  generating.value = true
  const params = {
    ...kelingParams.value,
    task_type: 'keling',
  }

  httpPost('/api/video/create', params)
    .then(() => {
      fetchData(1)
      taskPulling.value = true
      showToast('创建任务成功')
    })
    .catch((e) => {
      showToast('创建任务失败：' + e.message)
    })
    .finally(() => {
      generating.value = false
    })
}

const fetchData = (_page) => {
  if (_page) {
    page.value = _page
  }
  listLoading.value = true
  httpGet('/api/video/list', { page: page.value, page_size: pageSize.value })
    .then((res) => {
      total.value = res.data.total
      let needPull = false
      const items = []
      for (let v of res.data.items) {
        if (v.progress === 0 || v.progress === 102) {
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

const fetchUserPower = () => {
  httpGet('/api/user/power')
    .then((res) => {
      availablePower.value = res.data.power || 0
      lumaPowerCost.value = 10 // 示例值
      kelingPowerCost.value = 15 // 示例值
    })
    .catch(() => {})
}

const loadMore = () => {
  page.value++
  fetchData()
}

const playVideo = (item) => {
  currentVideoUrl.value = item.video_url
  showVideoDialog.value = true
}

const downloadVideo = (item) => {
  item.downloading = true
  const link = document.createElement('a')
  link.href = item.video_url
  link.download = item.title || 'video.mp4'
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
  item.downloading = false
  showToast('开始下载')
}

const removeJob = (item) => {
  showDialog({
    title: '确认删除',
    message: '此操作将会删除任务相关文件，继续操作吗?',
    showCancelButton: true,
  })
    .then(() => {
      httpGet('/api/video/remove', { id: item.id })
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

// 选择器确认方法
const onAspectRatioConfirm = (value) => {
  kelingParams.value.aspect_ratio = value
  showAspectRatioPicker.value = false
}

const onModelConfirm = (value) => {
  kelingParams.value.model = value
  showModelPicker.value = false
}

const onDurationConfirm = (value) => {
  kelingParams.value.duration = value
  showDurationPicker.value = false
}

const onModeConfirm = (value) => {
  kelingParams.value.mode = value
  showModePicker.value = false
}

const onCameraControlConfirm = (value) => {
  kelingParams.value.camera_control.type = value
  showCameraControlPicker.value = false
}
</script>

<style lang="scss" scoped>
.mobile-video-create {
  min-height: 100vh;
  background: #f7f8fa;
  padding-bottom: 20px;

  .page-header {
    background: #fff;
  }

  .video-type-tabs {
    background: #fff;
    margin: 12px;
    border-radius: 8px;
    overflow: hidden;

    .tab-content {
      padding: 16px;
    }

    .params-container {
      .image-upload-section {
        margin: 16px 0;

        .image-upload-row {
          display: flex;
          gap: 12px;

          .image-upload-item {
            flex: 1;

            .upload-placeholder {
              display: flex;
              flex-direction: column;
              align-items: center;
              justify-content: center;
              height: 100px;
              background: #f5f5f5;
              border: 2px dashed #ddd;
              border-radius: 8px;
              color: #999;

              .van-icon {
                font-size: 24px;
                margin-bottom: 8px;
              }
            }
          }
        }
      }
    }
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

  .video-preview {
    width: 100%;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
  }
}
</style>
