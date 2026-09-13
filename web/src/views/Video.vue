<template>
  <div class="page-video">
    <!-- 左侧参数设置面板 -->
    <div class="params-panel pt-2">
      <div class="provider-buttons">
        <div class="provider-grid">
          <button
            v-for="p in providerOrder"
            :key="p"
            class="provider-btn text-base"
            :class="{ active: store.activeProvider === p }"
            @click="store.switchProvider(p)"
            type="button"
          >
            <i class="iconfont mr-2 !text-xl" :class="getProviderIcon(p)"></i>
            {{ getProviderName(p) }}
          </button>
        </div>
      </div>

      <div class="function-params pt-3">
        <div class="mb-2" v-if="store.providerModels.length > 0">
          <label class="label text-left font-bold">模型选择</label>
        </div>
        <ParamBuilder
          v-model="store.formData"
          v-model:required-keys="store.requiredKeys"
          :items="store.providerModels"
          @price-params-change="handlePriceParamsChange"
        />

        <div
          class="power-info flex items-center justify-between mb-4 mt-3 p-3 rounded-lg bg-gradient-to-r from-blue-50 to-purple-50 border border-blue-200 shadow-sm"
        >
          <div class="flex items-center space-x-2">
            <el-icon color="#f59e42" size="20"><i class="iconfont icon-lightning"></i></el-icon>
            <span class="font-medium text-gray-700">当前可用积分：</span>
            <span class="font-bold text-lg text-yellow-500">{{ store.availablePower }}</span>
          </div>
          <el-tooltip content="积分用于生成视频，每次生成会消耗对应积分" placement="left">
            <el-icon color="#a78bfa" size="18"><InfoFilled /></el-icon>
          </el-tooltip>
        </div>

        <div class="flex justify-center" v-if="store.providerModels.length > 0">
          <button
            @click="store.createVideoTask"
            :disabled="store.submitting"
            class="w-full py-3 bg-gradient-to-r from-blue-500 to-purple-600 text-white rounded-xl disabled:from-gray-400 disabled:to-gray-400 disabled:cursor-not-allowed hover:from-blue-600 hover:to-purple-700 transition-all duration-200 flex items-center justify-center space-x-2 text-base"
            type="button"
          >
            <i v-if="store.submitting" class="iconfont icon-loading animate-spin"></i>
            <i v-else class="iconfont icon-chuangzuo"></i>
            <span>立即生成 ({{ store.currentPowerCost }}积分)</span>
          </button>
        </div>
      </div>
    </div>

    <!-- 右侧任务列表 -->
    <div
      class="main-content"
      v-loading="store.loading"
      element-loading-background="rgba(100,100,100,0.3)"
    >
      <div class="job-list-box px-2 pt-2 pb-2">
        <h2 class="text-xl mb-1">任务列表</h2>
        <task-list :list="videoRunningJobsForList" />
      </div>

      <div class="works-header">
        <h2 class="h-title text-2xl">你的作品</h2>
        <div class="filter-buttons">
          <el-button-group>
            <el-button
              :type="store.taskFilter === 'all' ? 'primary' : 'default'"
              @click="store.switchTaskFilter('all')"
              size="small"
            >
              全部
            </el-button>
            <el-button
              v-for="p in providerOrder"
              :key="`filter-${p}`"
              :type="store.taskFilter === p ? 'primary' : 'default'"
              @click="store.switchTaskFilter(p)"
              size="small"
            >
              {{ getProviderName(p) }}
            </el-button>
          </el-button-group>
        </div>
      </div>

      <div class="video-list">
        <div class="list-box" v-if="!store.noData">
          <Waterfall
            :list="worksListForWaterfall"
            v-bind="videoWaterfallOptions"
            :is-loading="store.loading"
            :is-over="store.isOver"
            :lazyload="true"
            @afterRender="onWaterfallAfterRender"
          >
            <template #default="{ item }">
              <div class="video-task-item">
                <div
                  class="video-task-preview"
                  :class="{
                    'video-task-preview--failed': item.status === 'failed',
                    'video-task-preview--busy':
                      item.status === 'downloading' ||
                      item.status === 'pending' ||
                      item.status === 'in_progress',
                  }"
                  @mouseenter="handleTaskPreviewEnter"
                  @mouseleave="handleTaskPreviewLeave"
                >
                  <div
                    v-if="item.status === 'success'"
                    class="video-task-preview-inner video-task-preview-inner--success"
                  >
                    <video
                      class="video-task-video"
                      :src="store.replaceImg(item.video_url)"
                      preload="metadata"
                      loop
                      muted
                      playsinline
                      @loadedmetadata="handleTaskMediaReady"
                      @loadeddata="handleTaskMediaReady"
                      @error="handleTaskMediaReady"
                      @click.stop="store.playVideo(item)"
                    >
                      您的浏览器不支持视频播放
                    </video>
                  </div>
                  <div
                    v-else-if="item.status === 'downloading'"
                    class="video-task-preview-inner video-task-preview-inner--busy"
                  >
                    <div class="text-center px-2">
                      <div
                        class="animate-spin rounded-full h-12 w-12 border-2 border-indigo-500 border-t-transparent mx-auto"
                      ></div>
                      <span class="text-sm text-indigo-600 dark:text-indigo-300 mt-3 block">
                        视频下载中…
                      </span>
                    </div>
                  </div>
                  <div
                    v-else-if="item.status === 'failed'"
                    class="video-task-preview-inner video-task-preview-inner--failed"
                  >
                    <img class="video-task-fail-img" :src="taskFailedImage" alt="" />
                  </div>
                  <div
                    v-else-if="(item.progress || 0) > 0 && (item.progress || 0) < 100"
                    class="video-task-preview-inner video-task-preview-inner--busy"
                  >
                    <el-progress
                      type="circle"
                      :percentage="item.progress || 0"
                      :width="96"
                      :stroke-width="6"
                    >
                      <template #default="{ percentage }">
                        <span class="text-base font-medium text-gray-700">{{ percentage }}%</span>
                      </template>
                    </el-progress>
                  </div>
                  <div v-else class="video-task-preview-inner video-task-preview-inner--busy">
                    <Generating message="正在生成视频" />
                  </div>

                  <div class="video-task-overlay">
                    <div class="video-task-overlay-time">{{ dateFormat(item.created_at) }}</div>
                    <div class="video-task-tools">
                      <el-tooltip content="复制提示词" placement="top">
                        <button
                          type="button"
                          class="video-task-tool copy-prompt"
                          :data-clipboard-text="item.prompt"
                        >
                          <i class="iconfont icon-copy"></i>
                        </button>
                      </el-tooltip>
                      <el-tooltip content="任务详情（含完整提示词）" placement="top">
                        <button type="button" class="video-task-tool" @click="showTaskDetail(item)">
                          <i class="iconfont icon-info text-[#6366f1]"></i>
                        </button>
                      </el-tooltip>
                      <el-tooltip
                        v-if="item.status === 'success' && item.video_url"
                        content="下载视频"
                        placement="top"
                      >
                        <button
                          type="button"
                          class="video-task-tool"
                          :disabled="item.downloading"
                          @click="store.downloadVideo(item)"
                        >
                          <i v-if="!item.downloading" class="iconfont icon-download"></i>
                          <img
                            v-else
                            src="/images/loading.gif"
                            class="video-task-tool-loading"
                            alt=""
                          />
                        </button>
                      </el-tooltip>
                      <el-tooltip content="删除任务" placement="top">
                        <button
                          type="button"
                          class="video-task-tool video-task-tool--danger"
                          @click="store.removeJob(item)"
                        >
                          <i class="iconfont icon-remove"></i>
                        </button>
                      </el-tooltip>
                    </div>
                  </div>
                </div>
              </div>
            </template>
          </Waterfall>
        </div>

        <el-empty
          :image-size="100"
          :image="store.nodata"
          description="没有任何作品，赶紧去创作吧！"
          v-else
        />

        <div class="waterfall-load-more" v-if="!store.noData">
          <img
            :src="videoWaterfallOptions.loadProps.loading"
            class="waterfall-loading-icon"
            v-if="!waterfallRendered"
            alt=""
          />
          <div
            v-if="waterfallRendered && store.isOver"
            class="waterfall-no-more bg-[#f5f5f5] text-gray-500 rounded-md px-3 py-2 mt-3"
          >
            <span class="mr-2 text-base">没有更多数据了</span>
            <i class="iconfont icon-face"></i>
          </div>
        </div>
      </div>
    </div>

    <!-- 视频预览对话框 -->
    <el-dialog
      v-model="store.showDialog"
      title="预览视频"
      hide-footer
      @close="handleVideoDialogClose"
      width="auto"
    >
      <video
        ref="videoPlayerRef"
        style="max-width: 90vw; max-height: 90vh"
        :src="store.currentVideoUrl"
        preload="auto"
        :autoplay="true"
        controls="controls"
      >
        您的浏览器不支持视频播放
      </video>
    </el-dialog>

    <!-- 任务详情 -->
    <el-dialog
      v-model="detailDialogVisible"
      title="任务详情"
      width="680px"
      :close-on-click-modal="false"
      destroy-on-close
    >
      <div class="detail-content" v-if="currentDetail">
        <el-descriptions :column="1" border>
          <el-descriptions-item label="提供商">
            {{ getProviderName(currentDetail.type) }}
          </el-descriptions-item>
          <el-descriptions-item v-if="currentDetail.channel" label="渠道">
            {{ currentDetail.channel }}
          </el-descriptions-item>
          <el-descriptions-item label="状态">
            {{ statusConfig[currentDetail.status]?.label || currentDetail.status }}
          </el-descriptions-item>
          <el-descriptions-item label="进度">
            {{ currentDetail.progress ?? 0 }}%
          </el-descriptions-item>
          <el-descriptions-item label="任务 ID">
            {{ currentDetail.task_id || '-' }}
          </el-descriptions-item>
          <el-descriptions-item v-if="detailFlatParams.model" label="模型">
            {{ detailFlatParams.model }}
          </el-descriptions-item>
          <el-descriptions-item v-if="detailFlatParams.task_type" label="生成模式">
            {{ detailFlatParams.task_type }}
          </el-descriptions-item>
          <el-descriptions-item
            v-if="detailFlatParams.duration != null && detailFlatParams.duration !== ''"
            label="时长"
          >
            {{ detailFlatParams.duration }} 秒
          </el-descriptions-item>
          <el-descriptions-item
            v-if="detailFlatParams.size || detailFlatParams.resolution"
            label="分辨率"
          >
            {{ detailFlatParams.size || detailFlatParams.resolution }}
          </el-descriptions-item>
          <el-descriptions-item v-if="detailFlatParams.mode" label="模式">
            {{ detailFlatParams.mode }}
          </el-descriptions-item>
          <el-descriptions-item
            v-if="detailFlatParams.sound != null && detailFlatParams.sound !== ''"
            label="声音"
          >
            {{ detailFlatParams.sound }}
          </el-descriptions-item>
          <el-descriptions-item label="消耗积分">
            {{ currentDetail.power ?? 0 }}
          </el-descriptions-item>
          <el-descriptions-item label="创建时间">
            {{ dateFormat(currentDetail.created_at) }}
          </el-descriptions-item>
          <el-descriptions-item label="提示词">
            <div class="prompt-with-copy">
              <span class="break-all">{{ currentDetail.prompt || '（无）' }}</span>
              <el-tooltip v-if="currentDetail.prompt" content="复制提示词" placement="top">
                <button
                  type="button"
                  class="inline-btn"
                  @click="copyPromptText(currentDetail.prompt)"
                >
                  <i class="iconfont icon-copy ml-2 cursor-pointer shrink-0" />
                </button>
              </el-tooltip>
            </div>
          </el-descriptions-item>
          <el-descriptions-item
            v-if="currentDetail.err_msg && String(currentDetail.err_msg).trim()"
            label="错误信息"
          >
            <el-text type="danger" class="break-all">{{ currentDetail.err_msg }}</el-text>
          </el-descriptions-item>
        </el-descriptions>

        <template v-if="detailUploadImages.length">
          <h4 class="detail-section-title">参考 / 上传的图片</h4>
          <div class="detail-media-grid detail-media-grid--image">
            <el-image
              v-for="(url, idx) in detailUploadImages"
              :key="'uimg-' + idx"
              :src="getThumbURL(url, 160, 160)"
              :preview-src-list="detailUploadImages"
              :initial-index="idx"
              fit="cover"
              class="detail-thumb"
            />
          </div>
        </template>

        <template v-if="detailUploadVideos.length">
          <h4 class="detail-section-title">参考 / 上传的视频</h4>
          <div class="detail-media-stack">
            <video
              v-for="(url, idx) in detailUploadVideos"
              :key="'uvid-' + idx"
              :src="url"
              controls
              preload="metadata"
              class="detail-video-preview"
            >
              您的浏览器不支持视频播放
            </video>
          </div>
        </template>

        <template v-if="detailUploadAudios.length">
          <h4 class="detail-section-title">参考 / 上传的音频</h4>
          <div class="detail-media-stack">
            <audio
              v-for="(url, idx) in detailUploadAudios"
              :key="'uaud-' + idx"
              :src="url"
              controls
              class="w-full"
            />
          </div>
        </template>

        <template v-if="detailResultVideo">
          <h4 class="detail-section-title">生成结果</h4>
          <div class="detail-media-stack">
            <video
              :src="detailResultVideo"
              controls
              preload="metadata"
              class="detail-video-preview detail-video-preview--large"
            >
              您的浏览器不支持视频播放
            </video>
          </div>
        </template>

        <template v-if="detailOtherParamsText">
          <el-collapse class="detail-params-collapse">
            <el-collapse-item title="其他请求参数（JSON）" name="params">
              <pre class="detail-json">{{ detailOtherParamsText }}</pre>
            </el-collapse-item>
          </el-collapse>
        </template>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import ParamBuilder from '@/components/ParamBuilder.vue'
import TaskList from '@/components/TaskList.vue'
import Generating from '@/components/ui/Generating.vue'
import { useSharedStore } from '@/store/sharedata'
import { useVideoStore } from '@/store/video'
import { InfoFilled } from '@element-plus/icons-vue'
import { getProviderName } from '@/store/data/video_params'
import { dateFormat, getThumbURL, replaceImg } from '@/utils/libs'
import { ElMessage } from 'element-plus'
import { computed, onMounted, onUnmounted, ref, watch } from 'vue'
import { Waterfall } from 'vue-waterfall-plugin-next'
import 'vue-waterfall-plugin-next/dist/style.css'

const store = useVideoStore()
const sharedStore = useSharedStore()
const videoPlayerRef = ref(null)
const waterfallRendered = ref(false)
const waterfallOptions = sharedStore.waterfallOptions
const taskFailedImage = sharedStore.taskFailedImage
const videoWaterfallOptions = computed(() => ({
  ...waterfallOptions,
  gutter: 2,
  hasAroundGutter: false,
}))
const worksListForWaterfall = computed(() => {
  return store.currentList
})

function isVideoTaskRunning(item) {
  return item.status === 'pending' || item.status === 'in_progress'
}

/** TaskList：队首为执行中，其余为排队；按创建时间升序 */
const videoRunningJobsForList = computed(() => {
  const rows = store.taskList.filter(isVideoTaskRunning)
  return [...rows]
    .sort((a, b) => new Date(a.created_at).getTime() - new Date(b.created_at).getTime())
    .map((row) => ({
      id: row.id,
      progress: row.progress ?? 0,
    }))
})

const VIDEO_IMAGE_PARAM_KEYS = [
  'image_urls',
  'images',
  'input_reference',
  'image',
  'image_tail',
  'first_frame_image',
]

function toUrlList(v) {
  if (v == null) {
    return []
  }
  if (Array.isArray(v)) {
    return v.map((x) => (typeof x === 'string' ? x.trim() : x)).filter(Boolean)
  }
  if (typeof v === 'string' && v.trim()) {
    return [v.trim()]
  }
  return []
}

function pickRefUrl(ref) {
  if (!ref || typeof ref !== 'object') {
    return ''
  }
  const url = ref.url
  return typeof url === 'string' && url.trim() ? url.trim() : ''
}

function collectContentUrls(params, type) {
  if (!params || !Array.isArray(params.content)) {
    return []
  }
  const out = []
  for (const row of params.content) {
    if (!row || typeof row !== 'object' || row.type !== type) {
      continue
    }
    if (type === 'image_url') {
      const u = pickRefUrl(row.image_url)
      if (u) {
        out.push(u)
      }
    } else if (type === 'video_url') {
      const u = pickRefUrl(row.video_url)
      if (u) {
        out.push(u)
      }
    } else if (type === 'audio_url') {
      const u = pickRefUrl(row.audio_url)
      if (u) {
        out.push(u)
      }
    }
  }
  return out
}

function uniqUrls(urls) {
  return [...new Set(urls.filter(Boolean))]
}

function parseJobParams(item) {
  let params = {}
  try {
    if (item.params) {
      if (typeof item.params === 'string') {
        params = JSON.parse(item.params)
      } else if (typeof item.params === 'object') {
        params = { ...item.params }
      }
    }
  } catch (e) {
    console.error('解析视频任务 params 失败:', e)
  }
  return params
}

// 数据库存的是 VideoTask JSON，模型表单在嵌套字段 params 里
function mergeVideoTaskParams(parsed) {
  if (!parsed || typeof parsed !== 'object') {
    return {}
  }
  const inner = parsed.params
  if (inner && typeof inner === 'object' && !Array.isArray(inner)) {
    const { params: _nested, ...rest } = parsed
    return { ...rest, ...inner }
  }
  return { ...parsed }
}

const detailDialogVisible = ref(false)
const currentDetail = ref(null)

const detailFlatParams = computed(() => {
  const row = currentDetail.value
  if (!row) {
    return {}
  }
  return mergeVideoTaskParams(parseJobParams(row))
})

const detailUploadImages = computed(() => {
  const p = detailFlatParams.value
  const acc = []
  for (const k of VIDEO_IMAGE_PARAM_KEYS) {
    if (p[k] != null) {
      acc.push(...toUrlList(p[k]))
    }
  }
  acc.push(...collectContentUrls(p, 'image_url'))
  return uniqUrls(acc.map((u) => replaceImg(u)))
})

const detailUploadVideos = computed(() => {
  const p = detailFlatParams.value
  const fromFields = toUrlList(p.video_url)
  const fromContent = collectContentUrls(p, 'video_url')
  return uniqUrls([...fromFields, ...fromContent].map((u) => replaceImg(u)))
})

const detailUploadAudios = computed(() => {
  const p = detailFlatParams.value
  const fromFields = toUrlList(p.audio_url)
  const fromContent = collectContentUrls(p, 'audio_url')
  return uniqUrls([...fromFields, ...fromContent].map((u) => replaceImg(u)))
})

const detailResultVideo = computed(() => {
  const row = currentDetail.value
  if (!row || row.status !== 'success' || !row.video_url) {
    return ''
  }
  return replaceImg(row.video_url)
})

const PARAM_SHOWN_IN_DESCRIPTIONS = new Set([
  'model',
  'task_type',
  'duration',
  'size',
  'mode',
  'sound',
  'resolution',
  ...VIDEO_IMAGE_PARAM_KEYS,
  'prompt',
  'video_url',
  'audio_url',
  'content',
  'image_urls',
])

const detailOtherParamsText = computed(() => {
  const p = detailFlatParams.value
  const rest = {}
  for (const k of Object.keys(p)) {
    if (PARAM_SHOWN_IN_DESCRIPTIONS.has(k)) {
      continue
    }
    rest[k] = p[k]
  }
  if (Object.keys(rest).length === 0) {
    return ''
  }
  try {
    return JSON.stringify(rest, null, 2)
  } catch {
    return ''
  }
})

function showTaskDetail(item) {
  currentDetail.value = { ...item }
  detailDialogVisible.value = true
}

function copyPromptText(text) {
  navigator.clipboard
    .writeText(text)
    .then(() => {
      ElMessage.success('提示词已复制')
    })
    .catch(() => {
      ElMessage.error('复制失败')
    })
}

const providerOrder = computed(() =>
  ['sora', 'veo', 'doubao', 'keling', 'minimax', 'wan'].filter((p) => store.providers.includes(p))
)

const getProviderIcon = (provider) => {
  const icons = {
    sora: 'icon-sora',
    doubao: 'icon-doubao',
    veo: 'icon-gemini',
    keling: 'icon-keling',
    minimax: 'icon-minimax',
    wan: 'icon-wan',
  }
  return icons[provider] || 'icon-video'
}

// 状态配置
const statusConfig = {
  pending: { label: '等待中', type: 'info' },
  in_progress: { label: '进行中', type: 'primary' },
  downloading: { label: '下载中', type: 'success' },
  success: { label: '成功', type: 'success' },
  failed: { label: '失败', type: 'danger' },
}

function videoStatusTagType(status) {
  return statusConfig[status]?.type || 'info'
}

// 处理价格参数变化事件
const handlePriceParamsChange = () => {
  // 价格参数变化时，store 中的 watch 会自动触发 setCurrentPowerCost
  // setCurrentPowerCost 是异步的，会调用 API 获取最新积分值
  // 无需额外处理，watch 会自动更新 currentPowerCost
}

// 处理视频对话框关闭事件
const handleVideoDialogClose = () => {
  if (videoPlayerRef.value) {
    videoPlayerRef.value.pause()
    videoPlayerRef.value.currentTime = 0
  }
  store.showDialog = false
}

function handleTaskPreviewEnter(event) {
  const target = event.currentTarget
  if (!(target instanceof HTMLElement)) {
    return
  }
  const videoElement = target.querySelector('.video-task-video')
  if (!(videoElement instanceof HTMLVideoElement)) {
    return
  }
  videoElement.muted = true
  const playPromise = videoElement.play()
  if (playPromise && typeof playPromise.catch === 'function') {
    playPromise.catch(() => {})
  }
}

function handleTaskPreviewLeave(event) {
  const target = event.currentTarget
  if (!(target instanceof HTMLElement)) {
    return
  }
  const videoElement = target.querySelector('.video-task-video')
  if (!(videoElement instanceof HTMLVideoElement)) {
    return
  }
  videoElement.pause()
}

let taskMediaReadyRafId = 0

function handleTaskMediaReady() {
  if (taskMediaReadyRafId) {
    cancelAnimationFrame(taskMediaReadyRafId)
  }
  taskMediaReadyRafId = requestAnimationFrame(() => {
    window.dispatchEvent(new Event('resize'))
    taskMediaReadyRafId = 0
  })
}

function onWaterfallAfterRender() {
  waterfallRendered.value = true
  if (!store.loading && !store.isOver) {
    store.fetchData(store.page + 1)
  }
}

watch(
  () => store.loading,
  (value) => {
    if (value) {
      waterfallRendered.value = false
    }
  }
)

watch(
  () => store.isOver,
  (value) => {
    if (value) {
      waterfallRendered.value = true
    }
  }
)

onMounted(() => {
  store.init()
})

onUnmounted(() => {
  if (taskMediaReadyRafId) {
    cancelAnimationFrame(taskMediaReadyRafId)
    taskMediaReadyRafId = 0
  }
  store.cleanup()
})
</script>

<style lang="scss" scoped>
@use '../assets/css/video.scss' as *;

.provider-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 10px;
}

.provider-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 10px 12px;
  border-radius: 10px;
  border: 1px solid rgba(255, 255, 255, 0.08);
  background: rgba(0, 0, 0, 0.06);
  color: var(--text-theme-color);
  transition: all 0.15s ease;

  &.active {
    background: linear-gradient(90deg, #3b82f6, #a855f7);
    border-color: rgba(99, 102, 241, 0.6);
    color: #fff;
    box-shadow: 0 2px 8px rgba(59, 130, 246, 0.35);
  }

  &:hover:not(.active) {
    background: rgba(0, 0, 0, 0.1);
  }
}

.detail-content {
  :deep(.el-descriptions__label) {
    min-width: 112px;
  }
}

.detail-section-title {
  margin: 16px 0 8px;
  font-size: 15px;
  font-weight: 600;
  color: var(--text-theme-color, #252f76);
}

.prompt-with-copy {
  display: flex;
  align-items: flex-start;
  gap: 4px;
}

.inline-btn {
  border: none;
  padding: 0;
  margin: 0;
  background: transparent;
  cursor: pointer;
  color: inherit;
  line-height: inherit;
}

.detail-media-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

.detail-media-grid--image .detail-thumb {
  width: 120px;
  height: 120px;
  border-radius: 8px;
  overflow: hidden;
}

.detail-media-stack {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.detail-video-preview {
  max-width: 100%;
  max-height: 220px;
  border-radius: 8px;
  background: #000;
}

.detail-video-preview--large {
  max-height: 360px;
}

.detail-params-collapse {
  margin-top: 16px;
}

.detail-json {
  margin: 0;
  padding: 10px 12px;
  font-size: 12px;
  line-height: 1.45;
  max-height: 240px;
  overflow: auto;
  border-radius: 8px;
  background: var(--el-fill-color-light, #f5f7fa);
  white-space: pre-wrap;
  word-break: break-all;
}
</style>
