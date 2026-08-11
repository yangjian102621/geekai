<template>
  <el-dialog
    v-model="dialogVisible"
    :title="dialogTitle"
    class="ppt-preview-dialog"
    width="90%"
    top="5vh"
    destroy-on-close
    @closed="onDialogClosed"
  >
    <template v-if="loading">
      <div class="p-10 text-center text-gray-400">
        <el-icon class="is-loading"><Loading /></el-icon> 加载中...
      </div>
    </template>
    <template v-else-if="errorText">
      <el-alert type="error" :title="errorText" show-icon />
    </template>
    <template v-else>
      <el-alert
        v-if="taskStatus === 'failed' && taskErrorMessage"
        type="error"
        :title="taskErrorMessage"
        show-icon
        class="mb-3"
      />
      <div
        v-if="currentSlide"
        class="mb-3 flex flex-wrap items-center justify-end gap-2 border-b border-slate-100 pb-3"
      >
        <el-popover
          placement="bottom-end"
          :width="360"
          trigger="click"
          :title="currentSlide.title || '—'"
          :content="currentSlide.image_prompt || '—'"
        >
          <template #reference>
            <el-button type="primary" class="!m-0">显示提示词</el-button>
          </template>
        </el-popover>
        <el-button v-if="editable" class="!m-0" @click="openSlideEditDrawer">修改</el-button>
        <el-dropdown
          v-if="showExport && taskStatus === 'completed' && hasSlideImages"
          trigger="click"
          @command="onExportCommand"
        >
          <el-button class="!m-0" :loading="exportLoading">
            导出
            <el-icon class="ml-0.5"><ArrowDown /></el-icon>
          </el-button>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="pdf">导出 PDF</el-dropdown-item>
              <el-dropdown-item command="pptx">导出 PPTX</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
      <div class="ppt-preview-layout flex min-h-0 items-stretch gap-4">
        <div class="flex min-h-0 min-w-0 flex-1 flex-col gap-3">
          <template v-if="currentSlide">
            <div class="flex min-w-0 flex-col gap-2">
              <div class="relative aspect-video w-full overflow-hidden rounded-lg bg-[#f5f5f5]">
                <el-image
                  v-if="currentSlide.image_url"
                  :src="replaceImg(currentSlide.image_url)"
                  fit="contain"
                  class="absolute inset-0 h-full w-full"
                />
                <div
                  v-else-if="taskStatus === 'failed'"
                  class="absolute inset-0 flex min-h-0 flex-col items-center justify-center gap-2 bg-slate-100/95"
                >
                  <span class="text-sm font-medium text-slate-600">本页未生成配图</span>
                </div>
                <div
                  v-else
                  class="ppt-gen-surface absolute inset-0 flex min-h-0 flex-col items-center justify-center gap-2"
                >
                  <el-icon class="is-loading text-2xl text-[var(--el-color-primary)]">
                    <Loading />
                  </el-icon>
                  <span class="text-sm font-medium text-slate-600">正在生成中</span>
                  <span class="max-w-[90%] text-center text-xs text-slate-400"
                    >配图生成完成后将自动显示</span
                  >
                </div>
              </div>
            </div>
          </template>
          <el-empty v-else description="暂无幻灯片" />
        </div>
        <div class="ppt-preview-thumbs flex w-[120px] shrink-0 flex-col gap-2">
          <div
            v-for="(slide, idx) in slides"
            :key="slide.slide_index"
            class="relative w-full shrink-0 cursor-pointer overflow-hidden rounded-md border-2 border-transparent transition-colors duration-200"
            :class="
              slide.slide_index === currentSlide?.slide_index
                ? 'border-[var(--el-color-primary)]'
                : 'hover:border-gray-300'
            "
            @click="currentSlideIndex = idx"
          >
            <div class="relative aspect-video w-full overflow-hidden rounded-sm bg-slate-200/80">
              <el-image
                v-if="slide.image_url"
                :src="slideThumb(slide.image_url)"
                fit="cover"
                class="absolute inset-0 h-full w-full"
              />
              <div
                v-else-if="taskStatus === 'failed'"
                class="absolute inset-0 flex flex-col items-center justify-center bg-slate-200/90 px-0.5"
              >
                <span class="text-center text-[9px] leading-tight text-slate-600">未生成</span>
              </div>
              <div
                v-else
                class="ppt-gen-surface ppt-gen-surface--thumb absolute inset-0 flex flex-col items-center justify-center gap-0.5 px-0.5"
              >
                <el-icon class="is-loading text-[var(--el-color-primary)]"><Loading /></el-icon>
                <span class="text-center text-[9px] leading-tight text-slate-500">生成中</span>
              </div>
            </div>
            <span
              class="pointer-events-none absolute left-1 top-1 z-10 min-w-[1.125rem] rounded bg-black/65 px-1.5 py-1 text-center text-[14px] font-medium leading-none text-white shadow-sm backdrop-blur-[2px]"
            >
              {{ idx + 1 }}
            </span>
          </div>
        </div>
      </div>
    </template>
  </el-dialog>

  <el-drawer
    v-if="editable"
    v-model="slideEditDrawerVisible"
    title="修改幻灯片"
    direction="rtl"
    size="min(420px, 92vw)"
    destroy-on-close
    @closed="onSlideEditDrawerClosed"
  >
    <div v-if="currentSlide" class="flex flex-col gap-4">
      <div class="aspect-video w-full overflow-hidden rounded-lg bg-slate-100">
        <el-image
          v-if="currentSlide.image_url"
          :src="replaceImg(currentSlide.image_url)"
          fit="contain"
          class="h-full w-full"
        />
      </div>
      <div class="w-full">
        <div class="mb-2 text-sm font-medium text-slate-700">生成新版本</div>
        <el-input
          v-model="slideEditPrompt"
          type="textarea"
          :rows="4"
          placeholder="描述你想做的修改..."
          :disabled="slideEditGenerating"
        />
        <el-button
          type="primary"
          class="mt-3 w-full"
          :loading="slideEditGenerating"
          :disabled="!taskId || !currentSlide?.image_url"
          @click="submitSlideEdit"
        >
          生成
        </el-button>
      </div>
      <div>
        <div class="mb-2 text-sm font-medium text-slate-700">历史版本</div>
        <div class="mb-2 text-xs text-slate-400">点击下方历史图片，可切换当前幻灯片版本</div>
        <div
          v-if="!slideImageHistory.length"
          class="rounded-lg border border-dashed border-slate-200 py-6 text-center text-xs text-slate-400"
        >
          暂无历史版本
        </div>
        <div v-else class="grid max-h-[40vh] grid-cols-1 gap-3 overflow-y-auto pr-0.5 sm:grid-cols-2">
          <div
            v-for="(ver, vidx) in slideImageHistory"
            :key="vidx + '-' + (ver.image_url || '')"
            class="cursor-pointer rounded-md border-2 transition-colors"
            :class="
              activeSlideVersionIndex === vidx
                ? 'border-[var(--el-color-primary)] bg-[var(--el-fill-color-light)]'
                : 'border-transparent hover:border-slate-300'
            "
            @click="activateSlideVersion(vidx)"
          >
            <div class="relative aspect-video w-full overflow-hidden rounded-sm bg-slate-100">
              <el-image :src="historyThumb(ver.image_url)" fit="cover" class="h-full w-full" />
              <span
                v-if="activeSlideVersionIndex === vidx"
                class="absolute right-1 top-1 rounded bg-[var(--el-color-primary)] px-1.5 py-0.5 text-[10px] font-medium text-white shadow"
                >当前</span
              >
            </div>
            <p
              class="mt-1.5 line-clamp-3 px-0.5 text-left text-[11px] leading-snug text-slate-600"
              :title="slideVersionPromptText(ver, vidx)"
            >
              {{ slideVersionPromptText(ver, vidx) }}
            </p>
          </div>
        </div>
      </div>
    </div>
  </el-drawer>
</template>

<script setup>
import { computed, ref, watch } from 'vue'
import { ArrowDown, Loading } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { getThumbURL, replaceImg } from '@/utils/libs'
import { httpDownload, httpGet, httpPatch, httpPost } from '@/utils/http'

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false,
  },
  taskId: {
    type: String,
    default: '',
  },
  endpointBase: {
    type: String,
    required: true,
  },
  editable: {
    type: Boolean,
    default: false,
  },
  /** 预览弹窗内是否展示导出（管理端可在列表导出，预览仅浏览） */
  showExport: {
    type: Boolean,
    default: true,
  },
})

const emit = defineEmits(['update:modelValue', 'updated', 'closed'])

const dialogVisible = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value),
})

const loading = ref(false)
const errorText = ref('')
const slides = ref([])
const taskStatus = ref('')
const taskErrorMessage = ref('')
const taskTitle = ref('')
const currentSlideIndex = ref(0)
const exportLoading = ref(false)

const slideEditDrawerVisible = ref(false)
const slideEditPrompt = ref('')
const slideEditGenerating = ref(false)

const currentSlide = computed(() => slides.value[currentSlideIndex.value] || null)
const hasSlideImages = computed(() =>
  slides.value.some((slide) => slide.image_url && String(slide.image_url).trim())
)

const dialogTitle = computed(() => {
  if (loading.value) return '加载中...'
  if (errorText.value) return '幻灯片详情'
  const title = currentSlide.value?.title || taskTitle.value
  return title && String(title).trim() ? title : '幻灯片详情'
})

const slideImageHistory = computed(() => {
  const slide = currentSlide.value
  if (!slide) return []
  const history = slide.image_history
  if (Array.isArray(history) && history.length) {
    return history.map((item) => {
      if (typeof item === 'string') {
        return { image_url: item, prompt: '' }
      }
      return {
        image_url: item.image_url || '',
        prompt: item.prompt != null ? String(item.prompt) : '',
      }
    })
  }
  if (slide.image_url) {
    return [
      {
        image_url: slide.image_url,
        prompt: slide.image_prompt != null ? String(slide.image_prompt) : '',
      },
    ]
  }
  return []
})

const activeSlideVersionIndex = computed(() => {
  const slide = currentSlide.value
  if (!slide?.image_url) return -1
  const index = slideImageHistory.value.findIndex((item) => item.image_url === slide.image_url)
  return index >= 0 ? index : 0
})

const pptThumb = (url, width, height = 0) => {
  if (!url) return ''
  return getThumbURL(replaceImg(url), width, height)
}

const slideThumb = (url) => pptThumb(url, 120, 0)
const historyThumb = (url) => pptThumb(url, 200, 0)

const detailUrl = computed(() => {
  if (!props.taskId) return ''
  return `${props.endpointBase}/${encodeURIComponent(props.taskId)}`
})

const exportUrl = (format) =>
  `${props.endpointBase}/${encodeURIComponent(props.taskId)}/export?format=${encodeURIComponent(format)}`

const resetState = () => {
  slides.value = []
  taskStatus.value = ''
  taskErrorMessage.value = ''
  taskTitle.value = ''
  currentSlideIndex.value = 0
  errorText.value = ''
  slideEditDrawerVisible.value = false
  slideEditPrompt.value = ''
}

const loadDetail = async () => {
  if (!props.taskId || !dialogVisible.value) return
  loading.value = true
  errorText.value = ''
  try {
    const res = await httpGet(detailUrl.value)
    const data = res.data || {}
    taskTitle.value = data.title || ''
    taskStatus.value = data.status || ''
    taskErrorMessage.value = data.error_message || ''
    slides.value = data.slides || []
    currentSlideIndex.value = 0
  } catch (e) {
    errorText.value = '加载失败：' + e.message
  } finally {
    loading.value = false
  }
}

watch(
  () => [props.modelValue, props.taskId],
  ([visible, taskId]) => {
    if (!visible) return
    if (!taskId) {
      resetState()
      return
    }
    loadDetail()
  },
  { immediate: true }
)

const safeExportFileBase = (title, taskId) => {
  const raw = (title && String(title).trim()) || taskId || 'export'
  return String(raw)
    .replace(/[/\\:*?"<>|]/g, '_')
    .slice(0, 120)
}

const onExportCommand = async (cmd) => {
  if (cmd !== 'pdf' && cmd !== 'pptx') return
  exportLoading.value = true
  const ext = cmd === 'pdf' ? '.pdf' : '.pptx'
  try {
    const response = await httpDownload(exportUrl(cmd))
    if (response.status !== 200) {
      throw new Error('导出失败')
    }
    const blob = new Blob([response.data])
    const link = document.createElement('a')
    link.href = URL.createObjectURL(blob)
    link.download = `${safeExportFileBase(taskTitle.value, props.taskId)}${ext}`
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    URL.revokeObjectURL(link.href)
    ElMessage.success('已开始下载')
  } catch (e) {
    const msg = typeof e === 'string' ? e : e?.message || e?.response?.data?.message || '未知错误'
    ElMessage.error('导出失败：' + msg)
  } finally {
    exportLoading.value = false
  }
}

const openSlideEditDrawer = () => {
  if (!currentSlide.value?.image_url) {
    ElMessage.warning('该页暂无配图，无法编辑')
    return
  }
  slideEditPrompt.value = ''
  slideEditDrawerVisible.value = true
}

const onSlideEditDrawerClosed = () => {
  slideEditPrompt.value = ''
}

const mergeSlides = (nextSlides) => {
  if (!Array.isArray(nextSlides) || !nextSlides.length) return
  const currentSlideID = currentSlide.value?.slide_index
  slides.value = nextSlides
  if (currentSlideID != null) {
    const nextIndex = nextSlides.findIndex((slide) => slide.slide_index === currentSlideID)
    if (nextIndex >= 0) currentSlideIndex.value = nextIndex
  }
}

const submitSlideEdit = async () => {
  const prompt = (slideEditPrompt.value || '').trim()
  const slide = currentSlide.value
  if (!props.taskId || !slide?.slide_index) return
  if (!prompt) {
    ElMessage.warning('请填写修改说明')
    return
  }
  slideEditGenerating.value = true
  try {
    const res = await httpPost(
      `${props.endpointBase}/${encodeURIComponent(props.taskId)}/slides/${slide.slide_index}/edit-image`,
      { prompt }
    )
    const payload = res?.data || {}
    if (payload.slides?.length) mergeSlides(payload.slides)
    else await loadDetail()
    slideEditPrompt.value = ''
    ElMessage.success('已生成新版本')
    emit('updated')
  } catch (e) {
    const msg = typeof e === 'string' ? e : e?.message || e?.response?.data?.message || '生成失败'
    ElMessage.error(msg)
  } finally {
    slideEditGenerating.value = false
  }
}

function slideVersionPromptText(ver, index) {
  const prompt = (ver.prompt || '').trim()
  if (index === 0) {
    return prompt ? `初始：${prompt}` : '初始版本（文生图）'
  }
  return prompt ? `修改：${prompt}` : `版本 ${index + 1}（无说明）`
}

const activateSlideVersion = async (versionIndex) => {
  if (activeSlideVersionIndex.value === versionIndex) return
  const slide = currentSlide.value
  if (!props.taskId || !slide?.slide_index) return
  try {
    const res = await httpPatch(
      `${props.endpointBase}/${encodeURIComponent(props.taskId)}/slides/${slide.slide_index}/active-image`,
      { version_index: versionIndex }
    )
    const payload = res?.data || {}
    if (payload.slides?.length) mergeSlides(payload.slides)
    else await loadDetail()
    emit('updated')
  } catch (e) {
    const msg = typeof e === 'string' ? e : e?.message || e?.response?.data?.message || '切换失败'
    ElMessage.error(msg)
  }
}

const onDialogClosed = () => {
  resetState()
  emit('closed')
}
</script>

<style scoped>
.ppt-preview-dialog :deep(.el-dialog__body) {
  padding-top: 8px;
  overflow: hidden;
}

.ppt-preview-layout {
  height: calc(90vh - 120px);
  max-height: calc(90vh - 120px);
  min-height: 0;
  overflow: hidden;
}

.ppt-preview-thumbs {
  height: 100%;
  max-height: 100%;
  overflow-y: auto;
  overscroll-behavior: contain;
}

@keyframes ppt-shimmer {
  0% {
    background-position: 200% 0;
  }
  100% {
    background-position: -200% 0;
  }
}

.ppt-gen-surface {
  background: linear-gradient(105deg, #e4e7ec 0%, #ebeef3 35%, #f4f6f9 50%, #e4e7ec 100%);
  background-size: 220% 100%;
  animation: ppt-shimmer 2.2s ease-in-out infinite;
}

.ppt-gen-surface--thumb {
  animation-duration: 1.6s;
}
</style>
