<template>
  <div class="min-h-screen bg-gray-50">
    <!-- 页面头部 -->
    <div class="sticky top-0 z-40 bg-white shadow-sm">
      <div class="flex items-center px-4 h-14">
        <button
          @click="goBack"
          class="flex items-center justify-center w-8 h-8 rounded-full hover:bg-gray-100 transition-colors"
        >
          <i class="iconfont icon-back text-gray-600"></i>
        </button>
        <h1 class="flex-1 text-center text-lg font-semibold text-gray-900">音乐创作</h1>
        <div class="w-8"></div>
      </div>
    </div>

    <!-- 创作表单 -->
    <div class="p-4 space-y-6">
      <!-- 模式切换 -->
      <div class="bg-white rounded-xl p-4 shadow-sm">
        <div class="flex items-center justify-between mb-3">
          <span class="text-gray-900 font-medium">创作模式</span>
          <van-switch v-model="custom" @change="onModeChange" size="24px" />
        </div>
        <p class="text-sm text-gray-500">
          {{ custom ? '自定义模式：可设置歌词、风格等详细参数' : '简单模式：通过描述快速生成' }}
        </p>
      </div>

      <!-- 模型选择 -->
      <CustomSelect
        v-model="data.model"
        :options="models"
        label="模型版本"
        title="选择模型"
        @change="onModelSelect"
      >
        <template #option="{ option, selected }">
          <div class="flex items-center w-full">
            <span class="font-bold text-blue-600 mr-2">{{ option.label }}</span>
            <span class="text-xs text-gray-400">({{ option.value }})</span>
            <span v-if="selected" class="ml-auto text-green-500"
              ><i class="iconfont icon-success"></i
            ></span>
          </div>
        </template>
      </CustomSelect>

      <!-- 纯音乐开关 -->
      <div class="bg-white rounded-xl p-4 shadow-sm">
        <div class="flex items-center justify-between">
          <div>
            <span class="text-gray-900 font-medium">纯音乐</span>
            <p class="text-sm text-gray-500 mt-1">生成不包含人声的音乐</p>
          </div>
          <van-switch v-model="data.instrumental" size="24px" />
        </div>
      </div>

      <!-- 自定义模式内容 -->
      <div v-if="custom" class="space-y-6">
        <!-- 歌词输入 -->
        <div v-if="!data.instrumental" class="bg-white rounded-xl p-4 shadow-sm">
          <label class="block text-gray-700 font-medium mb-3">歌词</label>
          <textarea
            v-model="data.lyrics"
            placeholder="请在这里输入你自己写的歌词..."
            class="w-full px-4 py-3 border border-gray-200 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent resize-none"
            rows="6"
            maxlength="2000"
          />
          <div class="flex items-center justify-between mt-3">
            <span class="text-sm text-gray-500">{{ data.lyrics.length }}/2000</span>
            <button
              @click="createLyric"
              :disabled="isGenerating || !data.lyrics"
              class="px-4 py-2 bg-blue-600 text-white rounded-lg font-medium disabled:bg-gray-300 disabled:cursor-not-allowed hover:bg-blue-700 transition-colors flex items-center space-x-2"
            >
              <i v-if="isGenerating" class="iconfont icon-loading animate-spin"></i>
              <span>{{ isGenerating ? '生成中...' : '生成歌词' }}</span>
            </button>
          </div>
        </div>

        <!-- 音乐风格 -->
        <div class="bg-white rounded-xl p-4 shadow-sm">
          <label class="block text-gray-700 font-medium mb-3">音乐风格</label>
          <textarea
            v-model="data.tags"
            placeholder="请输入音乐风格，多个风格之间用英文逗号隔开..."
            class="w-full px-4 py-3 border border-gray-200 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent resize-none"
            rows="3"
            maxlength="120"
          />
          <div class="flex justify-between items-center mt-2 mb-3">
            <span class="text-sm text-gray-500">{{ data.tags.length }}/120</span>
          </div>
          <!-- 风格标签选择 -->
          <div class="flex flex-wrap gap-2">
            <button
              v-for="tag in tags"
              :key="tag.value"
              @click="selectTag(tag)"
              class="px-3 py-1 text-sm border border-blue-200 text-blue-600 rounded-full hover:bg-blue-50 transition-colors"
            >
              {{ tag.label }}
            </button>
          </div>
        </div>

        <!-- 歌曲名称 -->
        <div class="bg-white rounded-xl p-4 shadow-sm">
          <label class="block text-gray-700 font-medium mb-3">歌曲名称</label>
          <input
            v-model="data.title"
            placeholder="请输入歌曲名称..."
            class="w-full px-4 py-3 border border-gray-200 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            maxlength="100"
          />
          <div class="text-right mt-2">
            <span class="text-sm text-gray-500">{{ data.title.length }}/100</span>
          </div>
        </div>
      </div>

      <!-- 简单模式内容 -->
      <div v-else class="bg-white rounded-xl p-4 shadow-sm">
        <label class="block text-gray-700 font-medium mb-3">歌曲描述</label>
        <textarea
          v-model="data.prompt"
          placeholder="例如：一首关于爱情的摇滚歌曲..."
          class="w-full px-4 py-3 border border-gray-200 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent resize-none"
          rows="6"
          maxlength="1000"
        />
        <div class="text-right mt-2">
          <span class="text-sm text-gray-500">{{ data.prompt.length }}/1000</span>
        </div>
      </div>

      <!-- 续写歌曲 -->
      <div v-if="refSong" class="bg-white rounded-xl p-4 shadow-sm border-l-4 border-orange-400">
        <div class="flex items-center justify-between mb-3">
          <h3 class="text-gray-900 font-medium">续写歌曲</h3>
          <button
            @click="removeRefSong"
            class="px-3 py-1 text-sm bg-red-100 text-red-600 rounded-lg hover:bg-red-200 transition-colors"
          >
            移除
          </button>
        </div>
        <div class="space-y-3">
          <div class="flex justify-between">
            <span class="text-gray-600">歌曲名称：</span>
            <span class="text-gray-900 font-medium">{{ refSong.title }}</span>
          </div>
          <div>
            <label class="block text-gray-700 font-medium mb-2">续写开始时间(秒)</label>
            <input
              v-model="refSong.extend_secs"
              type="number"
              placeholder="从第几秒开始续写"
              class="w-full px-4 py-3 border border-gray-200 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            />
          </div>
        </div>
      </div>

      <!-- 上传音乐 -->
      <div class="bg-white rounded-xl p-4 shadow-sm">
        <label class="block text-gray-700 font-medium mb-3">上传音乐文件(可选)</label>
        <div
          class="border-2 border-dashed border-gray-300 rounded-lg p-6 text-center hover:border-blue-400 hover:bg-blue-50 transition-colors cursor-pointer"
        >
          <input
            ref="fileInput"
            type="file"
            accept=".wav,.mp3"
            @change="handleFileSelect"
            class="hidden"
          />
          <div @click="$refs.fileInput.click()" class="flex flex-col items-center space-y-2">
            <i class="iconfont icon-upload text-blue-500 text-2xl"></i>
            <span class="text-gray-700 font-medium">上传音乐文件</span>
            <small class="text-gray-500">支持 .wav, .mp3 格式</small>
          </div>
        </div>
        <div v-if="uploadFiles.length > 0" class="mt-3 p-3 bg-gray-50 rounded-lg">
          <div class="flex items-center space-x-2">
            <i class="iconfont icon-success text-green-500"></i>
            <span class="text-sm text-gray-700">{{ uploadFiles[0].name }}</span>
          </div>
        </div>
      </div>

      <!-- 生成按钮 -->
      <div class="sticky bottom-4 bg-white rounded-xl p-4 shadow-lg">
        <button
          @click="create"
          :disabled="loading"
          class="w-full py-4 bg-gradient-to-r from-blue-500 to-purple-600 text-white font-semibold rounded-xl disabled:from-gray-400 disabled:to-gray-400 disabled:cursor-not-allowed hover:from-blue-600 hover:to-purple-700 transition-all duration-200 flex items-center justify-center space-x-2"
        >
          <i v-if="loading" class="iconfont icon-loading animate-spin"></i>
          <span>{{ loading ? '创作中...' : btnText }}</span>
        </button>
      </div>
    </div>

    <!-- 作品列表 -->
    <div class="p-4">
      <h2 class="text-lg font-semibold text-gray-900 mb-4">我的作品</h2>
      <div class="space-y-4">
        <div v-for="item in list" :key="item.id" class="bg-white rounded-xl p-4 shadow-sm">
          <div class="flex space-x-4">
            <div class="flex-shrink-0">
              <div class="relative w-16 h-16 rounded-lg overflow-hidden bg-gray-100">
                <el-image
                  :src="item.cover_url"
                  fit="cover"
                  class="w-full h-full"
                  :preview-disabled="true"
                >
                  <template #error>
                    <div class="w-full h-full flex items-center justify-center bg-gray-100">
                      <i class="iconfont icon-mp3 text-gray-400 text-xl"></i>
                    </div>
                  </template>
                </el-image>
                <!-- 音乐播放按钮 -->
                <button
                  v-if="item.progress === 100"
                  @click="play(item)"
                  class="absolute inset-0 flex items-center justify-center bg-black bg-opacity-50 opacity-0 hover:opacity-100 transition-opacity"
                >
                  <i class="iconfont icon-play text-white text-xl"></i>
                </button>
                <!-- 进度动画 -->
                <div
                  v-if="item.progress < 100 && item.progress !== 101"
                  class="absolute inset-0 flex items-center justify-center bg-blue-500 bg-opacity-20"
                >
                  <i class="iconfont icon-loading animate-spin text-blue-500 text-xl"></i>
                </div>
                <!-- 失败状态 -->
                <div
                  v-if="item.progress === 101"
                  class="absolute inset-0 flex items-center justify-center bg-red-500 bg-opacity-20"
                >
                  <i class="iconfont icon-warning text-red-500 text-xl"></i>
                </div>
              </div>
            </div>
            <div class="flex-1 min-w-0">
              <div class="flex items-start justify-between">
                <div class="flex-1">
                  <h3 class="text-gray-900 font-medium truncate">
                    {{ item.title || '未命名歌曲' }}
                  </h3>
                  <p class="text-gray-500 text-sm mt-1 line-clamp-2">
                    {{ item.tags || item.prompt }}
                  </p>
                </div>
                <!-- 任务状态 -->
                <div v-if="item.progress < 100" class="flex items-center space-x-2 text-sm">
                  <div
                    v-if="item.progress === 101"
                    class="text-red-600 flex items-center space-x-1"
                  >
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M12 9v2m0 4h.01"
                      />
                    </svg>
                    <span>失败</span>
                  </div>
                  <div v-else class="text-blue-600 flex items-center space-x-1">
                    <div
                      class="w-3 h-3 border border-blue-600 border-t-transparent rounded-full animate-spin"
                    ></div>
                    <span>生成中</span>
                  </div>
                </div>
              </div>
              <!-- 标签 -->
              <div class="flex items-center space-x-2 mt-2">
                <span
                  v-if="item.major_model_version"
                  class="px-2 py-1 text-xs bg-blue-100 text-blue-600 rounded-full"
                >
                  {{ item.major_model_version }}
                </span>
                <span
                  v-if="item.type === 4"
                  class="px-2 py-1 text-xs bg-green-100 text-green-600 rounded-full"
                >
                  <i class="iconfont icon-upload mr-1"></i>用户上传
                </span>
                <span
                  v-if="item.type === 3"
                  class="px-2 py-1 text-xs bg-yellow-100 text-yellow-600 rounded-full"
                >
                  <i class="iconfont icon-mp3 mr-1"></i>完整歌曲
                </span>
                <span
                  v-if="item.ref_song"
                  class="px-2 py-1 text-xs bg-purple-100 text-purple-600 rounded-full"
                >
                  <i class="iconfont icon-link mr-1"></i>续写
                </span>
              </div>
            </div>
          </div>

          <!-- 操作按钮 -->
          <div class="flex items-center justify-between mt-4">
            <div class="flex space-x-2">
              <button
                v-if="item.progress === 100"
                @click="play(item)"
                class="px-3 py-1.5 bg-blue-600 text-white text-sm rounded-lg hover:bg-blue-700 transition-colors flex items-center space-x-1"
              >
                <i class="iconfont icon-play !text-xs"></i>
                <span>播放</span>
              </button>
              <button
                v-if="item.progress === 100"
                @click="download(item)"
                :disabled="item.downloading"
                class="px-3 py-1.5 bg-green-600 text-white text-sm rounded-lg hover:bg-green-700 transition-colors disabled:bg-gray-400 flex items-center space-x-1"
              >
                <svg
                  v-if="item.downloading"
                  class="w-3 h-3 animate-spin"
                  fill="none"
                  viewBox="0 0 24 24"
                >
                  <circle
                    class="opacity-25"
                    cx="12"
                    cy="12"
                    r="10"
                    stroke="currentColor"
                    stroke-width="4"
                  />
                  <path
                    class="opacity-75"
                    fill="currentColor"
                    d="m4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                  />
                </svg>
                <i v-else class="iconfont icon-download !text-xs"></i>
                <span>{{ item.downloading ? '下载中...' : '下载' }}</span>
              </button>
            </div>
            <button
              @click="showDeleteDialog(item)"
              class="px-3 py-1.5 bg-red-100 text-red-600 text-sm rounded-lg hover:bg-red-200 transition-colors flex items-center space-x-1"
            >
              <i class="iconfont icon-remove !text-xs"></i>
              <span>删除</span>
            </button>
          </div>

          <!-- 进度条 -->
          <div v-if="item.progress < 100 && item.progress !== 101" class="mt-4">
            <div class="flex justify-between text-sm text-gray-600 mb-1">
              <span>生成进度</span>
              <span>{{ item.progress }}%</span>
            </div>
            <div class="w-full bg-gray-200 rounded-full h-2">
              <div
                class="bg-blue-600 h-2 rounded-full transition-all duration-300"
                :style="{ width: item.progress + '%' }"
              ></div>
            </div>
          </div>

          <!-- 错误信息 -->
          <div
            v-if="item.progress === 101"
            class="mt-4 p-3 bg-red-50 border border-red-200 rounded-lg"
          >
            <div class="flex items-start space-x-2">
              <div>
                <p class="text-red-600 text-sm">{{ item.err_msg || '未知错误' }}</p>
              </div>
            </div>
          </div>
        </div>

        <!-- 加载更多 -->
        <div v-if="listLoading" class="flex justify-center py-4">
          <svg class="w-6 h-6 animate-spin text-gray-400" fill="none" viewBox="0 0 24 24">
            <circle
              class="opacity-25"
              cx="12"
              cy="12"
              r="10"
              stroke="currentColor"
              stroke-width="4"
            />
            <path
              class="opacity-75"
              fill="currentColor"
              d="m4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
            />
          </svg>
        </div>

        <!-- 没有更多了 -->
        <div v-if="listFinished && !listLoading" class="text-center py-4 text-gray-500">
          没有更多了
        </div>
      </div>
    </div>

    <!-- 音乐播放器 -->
    <div
      v-if="showPlayer"
      class="fixed inset-0 z-50 flex items-end justify-center bg-black bg-opacity-50"
      @click="showPlayer = false"
    >
      <div @click.stop class="bg-white rounded-t-2xl w-full max-w-md animate-slide-up">
        <div class="flex items-center justify-between p-4 border-b">
          <h3 class="text-lg font-semibold text-gray-900">正在播放</h3>
          <button @click="showPlayer = false" class="p-2 hover:bg-gray-100 rounded-full">
            <svg
              class="w-5 h-5 text-gray-500"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M6 18L18 6M6 6l12 12"
              />
            </svg>
          </button>
        </div>
        <div class="p-6">
          <audio
            v-if="currentAudio"
            :src="currentAudio"
            controls
            autoplay
            class="w-full rounded-lg"
          >
            您的浏览器不支持音频播放
          </audio>
        </div>
      </div>
    </div>

    <!-- 删除确认对话框 -->
    <!-- 已移除，改为 van showConfirmDialog 方式 -->
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { httpGet, httpPost, httpDownload } from '@/utils/http'
import { checkSession } from '@/store/cache'
import CustomSelect from '@/views/mobile/components/CustomSelect.vue'
import { showToastMessage, showLoading, closeLoading } from '@/utils/dialog'
import { replaceImg } from '@/utils/libs'
import { showConfirmDialog } from 'vant'

const router = useRouter()

// 响应式数据
const custom = ref(false)
const data = ref({
  model: 'chirp-auk',
  tags: '',
  lyrics: '',
  prompt: '',
  title: '',
  instrumental: false,
  ref_task_id: '',
  extend_secs: 0,
  ref_song_id: '',
})
const loading = ref(false)
const list = ref([])
const listLoading = ref(false)
const listFinished = ref(false)
const btnText = ref('开始创作')
const refSong = ref(null)
const showModelPicker = ref(false)
const showPlayer = ref(false)
const showDeleteModal = ref(false)
const currentAudio = ref('')
const uploadFiles = ref([])
const isGenerating = ref(false)
const deleting = ref(false)
const deleteItem = ref(null)

// 模型选项
const models = ref([
  { label: 'v3.0', value: 'chirp-v3-0' },
  { label: 'v3.5', value: 'chirp-v3-5' },
  { label: 'v4.0', value: 'chirp-v4' },
  { label: 'v4.5', value: 'chirp-auk' },
])

const onModelSelect = (selectedModel) => {
  data.value.model = selectedModel.value
}

// 风格标签
const tags = ref([
  { label: '女声', value: 'female vocals' },
  { label: '男声', value: 'male vocals' },
  { label: '流行', value: 'pop' },
  { label: '摇滚', value: 'rock' },
  { label: '电音', value: 'electronic' },
  { label: '钢琴', value: 'piano' },
  { label: '吉他', value: 'guitar' },
  { label: '嘻哈', value: 'hip hop' },
])

// 页面数据
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)
const taskPulling = ref(true)
const tastPullHandler = ref(null)

// 滚动监听，自动加载更多
const handleScroll = () => {
  const scrollTop = window.pageYOffset || document.documentElement.scrollTop
  const windowHeight = window.innerHeight
  const documentHeight = document.documentElement.scrollHeight

  // 当滚动到底部附近时加载更多
  if (scrollTop + windowHeight >= documentHeight - 100) {
    loadMore()
  }
}

// 生命周期
onMounted(() => {
  checkSession()
    .then(() => {
      fetchData(1)
      // 启动定时轮询，检查任务状态
      tastPullHandler.value = setInterval(() => {
        if (taskPulling.value) {
          fetchData(1) // 只刷新第一页数据
        }
      }, 5000)

      // 添加滚动监听
      window.addEventListener('scroll', handleScroll)
    })
    .catch(() => {})
})

onUnmounted(() => {
  if (tastPullHandler.value) {
    clearInterval(tastPullHandler.value)
  }
  // 移除滚动监听
  window.removeEventListener('scroll', handleScroll)
})

// 方法
const goBack = () => {
  router.back()
}

const onModeChange = () => {
  if (!custom.value) {
    removeRefSong()
  }
}

const onModelConfirm = (value) => {
  const selectedModel = models.value.find((item) => item.label === value)
  if (selectedModel) {
    data.value.model = selectedModel.value
  }
  showModelPicker.value = false
}

const selectTag = (tag) => {
  if (data.value.tags.length + tag.value.length >= 119) {
    showToastMessage('标签长度超出限制', 'error')
    return
  }
  const currentTags = data.value.tags.split(',').filter((t) => t.trim())
  if (!currentTags.includes(tag.value)) {
    currentTags.push(tag.value)
    data.value.tags = currentTags.join(',')
  }
}

const createLyric = () => {
  if (data.value.lyrics === '') {
    showToastMessage('请输入歌词描述', 'error')
    return
  }
  isGenerating.value = true
  httpPost('/api/prompt/lyric', { prompt: data.value.lyrics })
    .then((res) => {
      const lines = res.data.split('\n')
      data.value.title = lines.shift().replace(/\*/g, '')
      lines.shift()
      data.value.lyrics = lines.join('\n')
      showToastMessage('歌词生成成功', 'success')
    })
    .catch((e) => {
      showToastMessage('歌词生成失败：' + e.message, 'error')
    })
    .finally(() => {
      isGenerating.value = false
    })
}

const handleFileSelect = (event) => {
  const file = event.target.files[0]
  if (!file) return

  uploadFiles.value = [{ file, name: file.name }]
  uploadAudio({ file, name: file.name })
}

const uploadAudio = (file) => {
  const formData = new FormData()
  formData.append('file', file.file, file.name)
  showLoading('正在上传文件...')

  httpPost('/api/upload', formData)
    .then((res) => {
      httpPost('/api/suno/create', {
        audio_url: res.data.url,
        title: res.data.name,
        type: 4,
      })
        .then(() => {
          fetchData(1)
          showToastMessage('歌曲上传成功', 'success')
          removeRefSong()
        })
        .catch((e) => {
          showToastMessage('歌曲上传失败：' + e.message, 'error')
        })
        .finally(() => {
          closeLoading()
        })
    })
    .catch((e) => {
      showToastMessage('文件上传失败:' + e.message, 'error')
    })
    .finally(() => {
      closeLoading()
    })
}

const create = () => {
  data.value.type = custom.value ? 2 : 1
  data.value.ref_task_id = refSong.value ? refSong.value.task_id : ''
  data.value.ref_song_id = refSong.value ? refSong.value.song_id : ''
  data.value.extend_secs = refSong.value ? refSong.value.extend_secs : 0

  if (refSong.value) {
    if (data.value.extend_secs > refSong.value.duration) {
      showToastMessage('续写开始时间不能超过原歌曲长度', 'error')
      return
    }
  } else if (custom.value) {
    if (data.value.lyrics === '') {
      showToastMessage('请输入歌词', 'error')
      return
    }
    if (data.value.title === '') {
      showToastMessage('请输入歌曲标题', 'error')
      return
    }
  } else {
    if (data.value.prompt === '') {
      showToastMessage('请输入歌曲描述', 'error')
      return
    }
  }

  loading.value = true
  httpPost('/api/suno/create', data.value)
    .then(() => {
      fetchData(1)
      taskPulling.value = true
      showToastMessage('创建任务成功', 'success')
    })
    .catch((e) => {
      showToastMessage('创建任务失败：' + e.message, 'error')
    })
    .finally(() => {
      loading.value = false
    })
}

const fetchData = (_page) => {
  if (_page) {
    page.value = _page
  }
  listLoading.value = true
  httpGet('/api/suno/list', { page: page.value, page_size: pageSize.value })
    .then((res) => {
      total.value = res.data.total
      let needPull = false
      const items = []
      for (let v of res.data.items) {
        if (v.progress === 100) {
          v.major_model_version = v['raw_data']['major_model_version']
        }
        // 检查是否有未完成的任务（进度为 0 或 102）
        if (v.progress === 0 || v.progress === 102) {
          needPull = true
        }
        items.push(v)
      }
      listLoading.value = false
      taskPulling.value = needPull

      // 如果任务有变化，则刷新任务列表
      if (JSON.stringify(list.value) !== JSON.stringify(items)) {
        list.value = items
      }

      if (page.value === 1) {
        list.value = items
      } else {
        list.value.push(...items)
      }

      if (items.length < pageSize.value) {
        listFinished.value = true
      }
    })
    .catch((e) => {
      listLoading.value = false
      showToastMessage('获取作品列表失败：' + e.message, 'error')
    })
}

const loadMore = () => {
  if (!listFinished.value && !listLoading.value) {
    page.value++
    fetchData()
  }
}

const play = (item) => {
  currentAudio.value = item.audio_url
  showPlayer.value = true
}

const download = (item) => {
  const url = replaceImg(item.audio_url)
  const downloadURL = `${import.meta.env.VITE_API_HOST}/api/download?url=${url}`
  // parse filename
  const urlObj = new URL(url)
  const fileName = urlObj.pathname.split('/').pop()
  item.downloading = true
  httpDownload(downloadURL)
    .then((response) => {
      const blob = new Blob([response.data])
      const link = document.createElement('a')
      link.href = URL.createObjectURL(blob)
      link.download = fileName
      document.body.appendChild(link)
      link.click()
      document.body.removeChild(link)
      URL.revokeObjectURL(link.href)
      item.downloading = false
    })
    .catch(() => {
      showToastMessage('下载失败', 'error')
      item.downloading = false
    })
    .finally(() => {
      item.downloading = false
    })
}

const showDeleteDialog = (item) => {
  deleteItem.value = item
  showConfirmDialog({
    title: '确认删除',
    message: '此操作将会删除任务相关文件，继续操作吗？',
    confirmButtonText: '确认删除',
    cancelButtonText: '取消',
  })
    .then(() => {
      // on confirm
      if (!deleteItem.value) return
      deleting.value = true
      httpGet('/api/suno/remove', { id: deleteItem.value.id })
        .then(() => {
          showToastMessage('任务删除成功', 'success')
          fetchData(1)
          deleteItem.value = null
        })
        .catch((e) => {
          showToastMessage('任务删除失败：' + e.message, 'error')
        })
        .finally(() => {
          deleting.value = false
        })
    })
    .catch(() => {
      // on cancel
      deleteItem.value = null
    })
}

const removeRefSong = () => {
  refSong.value = null
  btnText.value = '开始创作'
}
</script>

<style scoped>
/* 自定义动画 */
@keyframes fade-in {
  from {
    opacity: 0;
    transform: translateY(-10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes fade-out {
  from {
    opacity: 1;
    transform: translateY(0);
  }
  to {
    opacity: 0;
    transform: translateY(-10px);
  }
}

@keyframes slide-up {
  from {
    opacity: 0;
    transform: translateY(100%);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes scale-up {
  from {
    opacity: 0;
    transform: scale(0.9);
  }
  to {
    opacity: 1;
    transform: scale(1);
  }
}

.animate-fade-in {
  animation: fade-in 0.3s ease-out;
}

.animate-fade-out {
  animation: fade-out 0.3s ease-out;
}

.animate-slide-up {
  animation: slide-up 0.3s ease-out;
}

.animate-scale-up {
  animation: scale-up 0.3s ease-out;
}

/* 文本截断 */
.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

/* 滚动监听自动加载更多 */
.scroll-container {
  height: 100vh;
  overflow-y: auto;
}

/* 深色模式适配 */
@media (prefers-color-scheme: dark) {
  .bg-gray-50 {
    background-color: #1f2937;
  }

  .bg-white {
    background-color: #374151;
  }

  .text-gray-900 {
    color: #f9fafb;
  }

  .text-gray-700 {
    color: #d1d5db;
  }

  .text-gray-600 {
    color: #9ca3af;
  }

  .text-gray-500 {
    color: #6b7280;
  }

  .border-gray-200 {
    border-color: #4b5563;
  }

  .bg-gray-100:hover {
    background-color: #4b5563;
  }
}
</style>
