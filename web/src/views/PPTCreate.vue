<template>
  <div class="grid gap-4 p-4 lg:grid-cols-[minmax(0,1.1fr)_minmax(0,1.4fr)]">
    <div>
      <el-card shadow="never" class="h-full">
        <template #header>
          <div class="flex flex-col gap-1">
            <span class="text-lg font-semibold">智能 PPT 生成</span>
            <span class="text-xs text-gray-400">粘贴你的大纲，或上传材料一键生成分镜和配图</span>
          </div>
        </template>

        <el-form label-position="top">
          <!-- 格式 -->
          <el-form-item label="格式">
            <div class="grid grid-cols-2 gap-3">
              <div
                class="flex cursor-pointer flex-col gap-1 rounded-lg border border-gray-200 p-3 transition-colors duration-200"
                :class="
                  form.mode === 'detailed'
                    ? 'border-[var(--el-color-primary)] ring-1 ring-[var(--el-color-primary)]'
                    : 'hover:border-gray-300'
                "
                @click="form.mode = 'detailed'"
              >
                <span class="text-sm font-semibold">详细演示文稿</span>
                <span class="text-xs text-gray-400"
                  >一整套包含全文和详情的演示文稿，适合邮件发送或单独阅读。</span
                >
              </div>
              <div
                class="flex cursor-pointer flex-col gap-1 rounded-lg border border-gray-200 p-3 transition-colors duration-200"
                :class="
                  form.mode === 'slides'
                    ? 'border-[var(--el-color-primary)] ring-1 ring-[var(--el-color-primary)]'
                    : 'hover:border-gray-300'
                "
                @click="form.mode = 'slides'"
              >
                <span class="text-sm font-semibold">演示用幻灯片</span>
                <span class="text-xs text-gray-400"
                  >简洁直观的幻灯片，附带要点，为演讲提供支持。</span
                >
              </div>
            </div>
          </el-form-item>

          <!-- 选择语言 -->
          <el-form-item label="选择语言">
            <el-select v-model="form.language" placeholder="请选择语言" class="w-full">
              <el-option label="中文" value="zh-CN" />
              <el-option label="English" value="en" />
            </el-select>
          </el-form-item>

          <div class="flex items-center gap-5">
            <!-- 时长 + PPT 页数 -->
            <el-form-item label="时长">
              <div class="flex gap-2">
                <div
                  class="cursor-pointer rounded-md border border-gray-200 px-5 py-0.5 transition-colors duration-200"
                  :class="
                    form.duration === 'short'
                      ? 'border-[var(--el-color-primary)] bg-[var(--el-color-primary)] text-white'
                      : 'hover:border-gray-300'
                  "
                  @click="setDuration('short')"
                >
                  短
                </div>
                <div
                  class="cursor-pointer rounded-md border border-gray-200 px-5 py-0.5 transition-colors duration-200"
                  :class="
                    form.duration === 'default'
                      ? 'border-[var(--el-color-primary)] bg-[var(--el-color-primary)] text-white'
                      : 'hover:border-gray-300'
                  "
                  @click="setDuration('default')"
                >
                  默认
                </div>
              </div>
            </el-form-item>
            <el-form-item label="PPT 页数">
              <el-input-number v-model="form.pages" :min="3" :max="30" class="w-full" />
            </el-form-item>
          </div>

          <!-- PPT素材来源（二选一） -->
          <CustomSlideTabs v-model="materialMode">
            <CustomTabPane name="paste" label="粘贴大纲" width="40%">
              <el-form-item label="把你的大纲或者内容粘贴在这里" required>
                <el-input
                  v-model="form.content"
                  type="textarea"
                  :rows="6"
                  placeholder="添加一份概略提纲，或指定受众、风格和重点，例如：为新手用户创建一套演示文稿，采用大胆活泼的风格，注重分步说明。"
                  @input="onPasteInput"
                />
              </el-form-item>
            </CustomTabPane>

            <CustomTabPane name="upload" label="上传材料" width="60%">
              <el-form-item label="上传 PPT 材料（支持 PDF/Word/TXT/Markdown 格式）">
                <FileUpload
                  v-model="materialFileUrl"
                  :max-count="1"
                  :max-size="10"
                  accept=".pdf,.doc,.docx,.txt,.md,.markdown"
                />
              </el-form-item>
            </CustomTabPane>
          </CustomSlideTabs>
          <el-form-item label="风格与设计要求（可选）">
            <el-input
              v-model="form.prompt"
              type="textarea"
              :rows="3"
              placeholder="例如：极简商务风、卡通手绘风"
            />
          </el-form-item>

          <div class="mt-3 flex justify-end">
            <el-button type="primary" :loading="creating" @click="createTask" size="large">
              <i class="iconfont icon-chuangzuo mr-1"></i>
              {{ creating ? '任务创建中...' : '生成' }}
            </el-button>
          </div>
        </el-form>

        <div class="mt-4 text-xs text-gray-500">
          <p>小提示：</p>
          <ul class="mt-1 list-disc pl-[18px]">
            <li class="mb-0.5">建议使用分层级的大纲（如 Markdown 标题）让分镜更加清晰</li>
            <li class="mb-0.5">页数会在生成分镜时作为约束，保证故事线完整</li>
          </ul>
        </div>
      </el-card>
    </div>

    <div v-loading="loadingList">
      <el-card shadow="never">
        <template #header>
          <div class="flex items-center justify-between">
            <span class="text-base font-semibold">生成任务列表</span>
            <el-button link type="primary" @click="fetchTaskList" :loading="loadingList"
              >刷新</el-button
            >
          </div>
        </template>
        <el-empty v-if="!taskList.length" description="暂无任务，在左侧创建任务吧～" />
        <div v-else class="flex flex-col gap-2">
          <div
            v-for="job in taskList"
            :key="job.task_id"
            class="group flex w-full gap-4 overflow-hidden rounded-2xl border border-slate-200/70 bg-white p-4 shadow-sm shadow-slate-900/[0.04] transition-all duration-200 hover:border-slate-300/90 hover:shadow-md hover:shadow-slate-900/[0.06]"
          >
            <div
              class="relative h-20 w-[8.75rem] max-w-[28vw] shrink-0 overflow-hidden rounded-xl bg-slate-100"
            >
              <el-image
                v-if="job.thumb"
                :src="pptListThumb(job.thumb)"
                fit="cover"
                class="h-full w-full cursor-pointer"
                loading="lazy"
                @click="openDetail(job.task_id)"
              />
              <div
                v-else-if="isJobGenerating(job)"
                class="ppt-list-gen-placeholder flex h-full w-full cursor-pointer flex-col items-center justify-center gap-0.5 px-1"
                @click="openDetail(job.task_id)"
              >
                <el-icon class="is-loading text-[var(--el-color-primary)]"><Loading /></el-icon>
                <span class="text-center text-[10px] leading-tight text-slate-500">正在生成中</span>
              </div>
              <div
                v-else
                class="flex h-full w-full items-center justify-center px-2 text-center text-[11px] leading-snug text-slate-400"
              >
                暂无预览
              </div>
            </div>
            <div class="flex min-w-0 w-full flex-1 flex-col gap-2">
              <!-- 第一行：标题 + 时间 -->
              <div class="flex min-w-0 w-full items-start justify-between gap-3">
                <h3
                  class="line-clamp-2 min-w-0 flex-1 text-[15px] font-medium leading-snug tracking-tight text-slate-900"
                >
                  {{ job.title || job.task_id }}
                </h3>
                <time
                  class="shrink-0 whitespace-nowrap text-xs tabular-nums text-slate-400"
                  :datetime="String(job.created_at)"
                >
                  {{ formatTime(job.created_at) }}
                </time>
              </div>
              <!-- 第二行：状态/页数 与 操作 同一行，100% 宽、两端对齐 -->
              <div
                class="flex w-full min-w-0 items-center justify-between gap-3 text-xs"
                @click.stop
              >
                <div class="flex min-w-0 shrink items-center gap-x-3 text-slate-500">
                  <el-tag :type="statusTagType(job.status)" size="small" effect="plain">
                    {{ statusText(job.status) }}
                  </el-tag>
                  <span class="tabular-nums text-slate-600"
                    >{{ job.completed_slides }} / {{ job.total_slides }} 页</span
                  >
                </div>
                <div class="flex flex-shrink-0 items-center justify-end gap-x-3 sm:gap-x-4">
                  <el-popover placement="bottom-start" :width="340" trigger="click">
                    <template #default>
                      <div
                        class="max-h-60 overflow-y-auto whitespace-pre-wrap break-words text-sm text-slate-700"
                      >
                        {{ job.prompt || '（无补充提示）' }}
                      </div>
                    </template>
                    <template #reference>
                      <el-button
                        size="small"
                        type="primary"
                        link
                        class="!px-0 !text-xs font-normal"
                      >
                        查看提示词
                      </el-button>
                    </template>
                  </el-popover>
                  <el-button
                    size="small"
                    link
                    class="!px-0 !text-xs font-normal text-slate-600"
                    @click="openContentDialog(job.task_id)"
                  >
                    查看大纲
                  </el-button>
                  <el-tooltip
                    v-if="canResumeJob(job)"
                    content="从中断处继续生成缺页配图（处理中不可重复点击）"
                    placement="top"
                  >
                    <span class="inline-flex">
                      <el-button
                        size="small"
                        link
                        type="primary"
                        class="!px-0 !text-xs font-normal"
                        :disabled="job.status === 'processing'"
                        :loading="resumeLoadingId === job.task_id"
                        @click="resumeTask(job.task_id)"
                      >
                        继续生成
                      </el-button>
                    </span>
                  </el-tooltip>
                  <el-dropdown
                    v-if="job.status === 'completed' && job.completed_slides > 0"
                    trigger="click"
                    @command="(cmd) => onListExportCommand(cmd, job)"
                    :disabled="exportLoadingKey === `${job.task_id}:list`"
                  >
                    <span class="inline-flex items-center">
                      <el-button
                        size="small"
                        link
                        type="primary"
                        class="!px-0 !text-xs font-normal"
                        :loading="exportLoadingKey === `${job.task_id}:list`"
                      >
                        导出
                        <el-icon class="ml-0.5 text-[10px]"><ArrowDown /></el-icon>
                      </el-button>
                    </span>
                    <template #dropdown>
                      <el-dropdown-menu>
                        <el-dropdown-item command="pdf">导出 PDF</el-dropdown-item>
                        <el-dropdown-item command="pptx">导出 PPTX</el-dropdown-item>
                      </el-dropdown-menu>
                    </template>
                  </el-dropdown>
                  <el-popconfirm
                    v-if="['completed', 'failed'].includes(job.status)"
                    confirm-button-text="确认"
                    cancel-button-text="取消"
                    title="确定删除该任务及生成的图片吗？"
                    @confirm="deleteTask(job.task_id)"
                  >
                    <template #reference>
                      <el-button type="danger" link size="small" class="!px-0 !text-xs font-normal">
                        删除
                      </el-button>
                    </template>
                  </el-popconfirm>
                </div>
              </div>
              <div
                v-if="job.status === 'failed' && job.error_message"
                class="rounded-lg bg-red-50/80 px-2.5 py-1.5 text-xs leading-relaxed text-red-700"
              >
                {{ job.error_message }}
              </div>
            </div>
          </div>
        </div>
        <div v-if="pagination.total > pagination.page_size" class="mt-4 flex justify-end">
          <el-pagination
            v-model:current-page="pagination.page"
            v-model:page-size="pagination.page_size"
            :page-sizes="[10, 20, 50]"
            :total="pagination.total"
            layout="total, sizes, prev, pager, next"
            @current-change="fetchTaskList"
            @size-change="
              () => {
                pagination.page = 1
                fetchTaskList()
              }
            "
          />
        </div>
      </el-card>
    </div>

    <PPTPreviewDialog
      v-model="previewVisible"
      :task-id="previewTaskId"
      endpoint-base="/api/v1/tasks"
      :editable="true"
      @updated="fetchTaskList"
      @closed="onPreviewClosed"
    />

    <!-- 大纲全文（懒加载） -->
    <el-dialog
      v-model="contentDialogVisible"
      title="PPT 大纲内容"
      width="min(640px, 92vw)"
      class="content-outline-dialog"
      destroy-on-close
      @closed="onContentDialogClosed"
    >
      <div v-if="contentDialogLoading" class="flex items-center justify-center py-12 text-gray-400">
        <el-icon class="is-loading mr-2"><Loading /></el-icon>
        加载中...
      </div>
      <el-input
        v-else
        :model-value="contentDialogText"
        type="textarea"
        :rows="18"
        readonly
        class="font-mono text-sm"
      />
    </el-dialog>
  </div>
</template>

<script setup>
import axios from 'axios'
import { getThumbURL, replaceImg } from '@/utils/libs'
import { httpDownload, httpGet, httpPost } from '@/utils/http'
import { ElMessage } from 'element-plus'
import { ArrowDown, Loading } from '@element-plus/icons-vue'
import { onBeforeUnmount, onMounted, ref, watch } from 'vue'
import FileUpload from '@/components/FileUpload.vue'
import CustomSlideTabs from '@/components/ui/CustomSlideTabs.vue'
import CustomTabPane from '@/components/ui/CustomTabPane.vue'
import PPTPreviewDialog from '@/components/ppt/PPTPreviewDialog.vue'

const form = ref({
  content: '',
  prompt: '',
  mode: 'detailed', // detailed | slides
  language: 'zh-CN',
  duration: 'default', // short | default
  pages: 10,
})

const materialFileUrl = ref('')
const materialMode = ref('paste')

watch(
  () => materialFileUrl.value,
  (val) => {
    // 当用户完成上传后，切换到“上传模式”，并清空粘贴内容，保证二选一。
    if (val) {
      materialMode.value = 'upload'
      form.value.content = ''
    } else if (materialMode.value === 'upload') {
      // 清空文件后回到粘贴模式，避免卡住。
      materialMode.value = 'paste'
    }
  }
)

watch(
  () => materialMode.value,
  (val) => {
    if (val === 'paste') {
      // 切到粘贴模式则清空已上传文件，确保二选一。
      if (materialFileUrl.value) materialFileUrl.value = ''
    } else if (val === 'upload') {
      // 切到上传模式则清空粘贴内容，确保二选一。
      if ((form.value.content || '').trim()) form.value.content = ''
    }
  }
)

const onPasteInput = () => {
  const v = (form.value.content || '').trim()
  // 当用户开始粘贴时，清空已上传文件，保证二选一。
  if (v) {
    materialMode.value = 'paste'
    materialFileUrl.value = ''
  }
}

const setDuration = (d) => {
  form.value.duration = d
  form.value.pages = d === 'short' ? 6 : 10
}

const creating = ref(false)
const taskList = ref([])
const loadingList = ref(false)
const pagination = ref({ page: 1, page_size: 20, total: 0 })
let listPollTimer = null

const previewVisible = ref(false)
const previewTaskId = ref('')

/** 与 Image.vue 一致：缩略图 URL（App.vue 已 initThumbTemplate） */
function pptThumb(url, width, height = 0) {
  if (!url) return ''
  return getThumbURL(replaceImg(url), width, height)
}
const pptListThumb = (url) => pptThumb(url, 280, 0)

/** 列表行 / 详情 导出中：`${taskId}:list` | `${taskId}:detail` */
const exportLoadingKey = ref('')
/** 继续生成请求中：当前 task_id */
const resumeLoadingId = ref('')

const canResumeJob = (job) => {
  if (!job || !job.task_id) return false
  const st = job.status
  if (st === 'completed' || st === 'processing') return false
  const total = Number(job.total_slides) || 0
  const done = Number(job.completed_slides) || 0
  return total > 0 && done < total
}

/** 列表缩略图：排队/生成中且无封面图时显示「正在生成」动效 */
const isJobGenerating = (job) => {
  if (!job) return false
  return ['pending', 'processing'].includes(job.status)
}

const safeExportFileBase = (title, taskId) => {
  const raw = (title && String(title).trim()) || taskId || 'export'
  return String(raw)
    .replace(/[/\\:*?"<>|]/g, '_')
    .slice(0, 120)
}

const runPptExportDownload = async (taskId, format, baseName, loadingSlot) => {
  const ext = format === 'pdf' ? '.pdf' : '.pptx'
  const url = `/api/v1/tasks/${encodeURIComponent(taskId)}/export?format=${encodeURIComponent(format)}`
  exportLoadingKey.value = `${taskId}:${loadingSlot}`
  try {
    const response = await httpDownload(url)
    if (response.status !== 200) {
      let msg = '导出失败'
      try {
        const text = await response.data.text()
        const j = JSON.parse(text)
        msg = j.message || j.msg || msg
      } catch {
        msg = '导出失败'
      }
      throw new Error(msg)
    }
    const blob = new Blob([response.data])
    const link = document.createElement('a')
    link.href = URL.createObjectURL(blob)
    link.download = `${safeExportFileBase(baseName, taskId)}${ext}`
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    URL.revokeObjectURL(link.href)
    ElMessage.success('已开始下载')
  } catch (e) {
    const msg = typeof e === 'string' ? e : e?.message || e?.response?.data?.message || '未知错误'
    ElMessage.error('导出失败：' + msg)
  } finally {
    exportLoadingKey.value = ''
  }
}

const onListExportCommand = (cmd, job) => {
  if (cmd !== 'pdf' && cmd !== 'pptx') return
  runPptExportDownload(job.task_id, cmd, job.title || job.task_id, 'list')
}

const resumeTask = async (taskId) => {
  if (!taskId || resumeLoadingId.value) return
  resumeLoadingId.value = taskId
  try {
    await httpPost(`/api/v1/tasks/${encodeURIComponent(taskId)}/resume`, {})
    ElMessage.success('已继续生成，请稍候…')
    fetchTaskList()
    startListPolling()
  } catch (e) {
    const msg = typeof e === 'string' ? e : e?.message || e?.response?.data?.message || '未知错误'
    ElMessage.error('继续生成失败：' + msg)
  } finally {
    resumeLoadingId.value = ''
  }
}

const statusText = (status) => {
  switch (status) {
    case 'pending':
      return '排队中'
    case 'processing':
      return '生成中'
    case 'completed':
      return '已完成'
    case 'failed':
      return '已失败'
    default:
      return status || '未知'
  }
}

const statusTagType = (status) => {
  switch (status) {
    case 'pending':
      return 'info'
    case 'processing':
      return 'warning'
    case 'completed':
      return 'success'
    case 'failed':
      return 'danger'
    default:
      return 'info'
  }
}

const createTask = () => {
  creating.value = true

  const hasFile = !!materialFileUrl.value
  const content = (form.value.content || '').trim()
  if (!hasFile && !content) {
    ElMessage.warning('请先填写演示文稿大纲，或上传 PPT 材料文件')
    creating.value = false
    return
  }

  if (hasFile) {
    httpPost('/api/v1/tasks/generate-slides/from-file', {
      file_url: materialFileUrl.value,
      prompt: form.value.prompt || '',
      language: form.value.language || 'zh-CN',
      pages: form.value.pages || 10,
      mode: form.value.mode || 'slides',
    })
      .then(() => {
        ElMessage.success('任务创建成功，正在生成 PPT...')
        fetchTaskList()
        startListPolling()
      })
      .catch((e) => {
        ElMessage.error('创建任务失败：' + (e?.message || 'unknown error'))
      })
      .finally(() => {
        creating.value = false
      })
    return
  }

  httpPost('/api/v1/tasks/generate-slides', {
    content: content,
    prompt: form.value.prompt,
    language: form.value.language || 'zh-CN',
    pages: form.value.pages || 10,
    mode: form.value.mode || 'slides',
  })
    .then(() => {
      ElMessage.success('任务创建成功，正在生成 PPT...')
      fetchTaskList()
      startListPolling()
    })
    .catch((e) => {
      ElMessage.error('创建任务失败：' + e.message)
    })
    .finally(() => {
      creating.value = false
    })
}

const fetchTaskList = () => {
  loadingList.value = true
  httpGet('/api/v1/tasks', {
    page: pagination.value.page,
    page_size: pagination.value.page_size,
  })
    .then((res) => {
      const data = res.data || {}
      taskList.value = data.jobs || []
      pagination.value.total = data.total ?? 0
      const hasActive = (data.jobs || []).some((j) => ['pending', 'processing'].includes(j.status))
      if (hasActive) {
        startListPolling()
      } else {
        stopListPolling()
      }
    })
    .catch((e) => {
      ElMessage.error('获取任务列表失败：' + e.message)
    })
    .finally(() => {
      loadingList.value = false
    })
}

const deleteTask = async (taskId) => {
  try {
    await axios.delete(`/api/v1/tasks/${taskId}`)
    ElMessage.success('删除成功')
    fetchTaskList()
    if (previewVisible.value && previewTaskId.value === taskId) {
      previewVisible.value = false
      onPreviewClosed()
    }
  } catch (e) {
    const msg = e?.message || e?.response?.data?.message || '删除失败'
    ElMessage.error('删除失败：' + msg)
  }
}

const formatTime = (ts) => {
  if (!ts) return ''
  const d = new Date(ts * 1000)
  return d.toLocaleString('zh-CN', {
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
  })
}

const startListPolling = () => {
  stopListPolling()
  listPollTimer = setInterval(() => fetchTaskList(), 3000)
}

const stopListPolling = () => {
  if (listPollTimer) {
    clearInterval(listPollTimer)
    listPollTimer = null
  }
}

const openDetail = (taskId) => {
  previewTaskId.value = taskId
  previewVisible.value = true
}

const onPreviewClosed = () => {
  previewTaskId.value = ''
}

const contentDialogVisible = ref(false)
const contentDialogLoading = ref(false)
const contentDialogText = ref('')

const openContentDialog = (taskId) => {
  contentDialogVisible.value = true
  contentDialogLoading.value = true
  contentDialogText.value = ''
  httpGet(`/api/v1/tasks/${taskId}`)
    .then((res) => {
      const data = res?.data || {}
      contentDialogText.value = data.content || ''
    })
    .catch((e) => {
      ElMessage.error('加载大纲失败：' + e.message)
      contentDialogVisible.value = false
    })
    .finally(() => {
      contentDialogLoading.value = false
    })
}

const onContentDialogClosed = () => {
  contentDialogText.value = ''
}

onMounted(() => {
  fetchTaskList()
})

onBeforeUnmount(() => {
  stopListPolling()
})
</script>

<style scoped>
/* 减少预览弹窗内容区顶部留白，使右侧缩略图与主图对齐更紧凑 */
.detail-dialog :deep(.el-dialog__body) {
  padding-top: 8px;
}

.content-outline-dialog :deep(.el-textarea__inner) {
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
}

/* 未出图占位：固定 16:9 区域 + 微光 + 与 el-icon Loading 的旋转动画一致 */
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

.ppt-list-gen-placeholder {
  background: linear-gradient(105deg, #e4e7ec 0%, #ebeef3 40%, #f4f6f9 55%, #e4e7ec 100%);
  background-size: 220% 100%;
  animation: ppt-shimmer 2.2s ease-in-out infinite;
}
</style>
