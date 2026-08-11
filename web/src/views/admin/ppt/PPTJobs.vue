<template>
  <div class="p-5">
    <!-- 统计信息 -->
    <el-row :gutter="20" class="stats-row">
      <el-col :span="4">
        <el-card class="stat-card">
          <div class="stat-item">
            <div class="stat-number">{{ stats.totalTasks }}</div>
            <div class="stat-label">总任务数</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="4">
        <el-card class="stat-card">
          <div class="stat-item">
            <div class="stat-number info">{{ stats.pendingTasks }}</div>
            <div class="stat-label">排队中</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="4">
        <el-card class="stat-card">
          <div class="stat-item">
            <div class="stat-number warning">{{ stats.processingTasks }}</div>
            <div class="stat-label">生成中</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="4">
        <el-card class="stat-card">
          <div class="stat-item">
            <div class="stat-number success">{{ stats.completedTasks }}</div>
            <div class="stat-label">已完成</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="4">
        <el-card class="stat-card">
          <div class="stat-item">
            <div class="stat-number danger">{{ stats.failedTasks }}</div>
            <div class="stat-label">失败</div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 搜索筛选 -->
    <el-card class="filter-card" shadow="never">
      <div class="flex items-center gap-2">
        <el-input v-model="queryForm.user_id" placeholder="用户ID" clearable style="width: 150px" />
        <el-select
          v-model="queryForm.status"
          placeholder="任务状态"
          clearable
          style="width: 150px"
          @change="handleQuery"
        >
          <el-option label="排队中" value="pending" />
          <el-option label="生成中" value="processing" />
          <el-option label="已完成" value="completed" />
          <el-option label="失败" value="failed" />
        </el-select>

        <el-button type="primary" @click="handleQuery" :loading="loading">
          <i class="iconfont icon-search mr-1" />
          搜索
        </el-button>
      </div>
    </el-card>

    <!-- 任务列表 -->
    <el-card class="table-card">
      <el-table :data="taskList" v-loading="loading" border>
        <el-table-column prop="task_id" label="任务ID" min-width="200" />
        <el-table-column prop="user_id" label="用户ID" width="100" />
        <el-table-column prop="status" label="状态" width="110">
          <template #default="scope">
            <el-tag :type="getStatusColor(scope.row.status)" size="small">
              {{ getStatusName(scope.row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="total_slides" label="总页数" width="90" />
        <el-table-column prop="completed_slides" label="已完成页数" width="110" />
        <el-table-column
          prop="error_message"
          label="错误信息"
          min-width="200"
          show-overflow-tooltip
        />
        <el-table-column prop="created_at" label="创建时间" width="180">
          <template #default="scope">
            {{ formatDateTime(scope.row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column prop="updated_at" label="更新时间" width="180">
          <template #default="scope">
            {{ formatDateTime(scope.row.updated_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="180" fixed="right">
          <template #default="scope">
            <div class="flex items-center gap-3">
              <el-button link type="primary" @click="openPreview(scope.row.task_id)">预览</el-button>
              <el-dropdown
                v-if="scope.row.status === 'completed' && scope.row.completed_slides > 0"
                trigger="click"
                @command="(cmd) => onExportCommand(cmd, scope.row)"
                :disabled="exportLoadingKey === scope.row.task_id"
              >
                <span class="inline-flex items-center">
                  <el-button link type="primary" :loading="exportLoadingKey === scope.row.task_id">
                    导出
                  </el-button>
                </span>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item command="pdf">导出 PDF</el-dropdown-item>
                    <el-dropdown-item command="pptx">导出 PPTX</el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
            </div>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-container">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.size"
          :page-sizes="[10, 20, 50, 100]"
          :total="pagination.total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>

    <PPTPreviewDialog
      v-model="previewVisible"
      :task-id="previewTaskId"
      endpoint-base="/api/admin/ppt/jobs"
      :editable="false"
      :show-export="false"
      @closed="onPreviewClosed"
    />
  </div>
</template>

<script setup>
import { httpDownload, httpGet } from '@/utils/http'
import { formatDateTime } from '@/utils/libs'
import { ElMessage } from 'element-plus'
import { onMounted, reactive, ref } from 'vue'
import PPTPreviewDialog from '@/components/ppt/PPTPreviewDialog.vue'

// 查询条件
const queryForm = reactive({
  user_id: '',
  status: '',
})

// 分页
const pagination = reactive({
  page: 1,
  size: 20,
  total: 0,
})

// 列表数据与状态
const taskList = ref([])
const loading = ref(false)
const exportLoadingKey = ref('')
const previewVisible = ref(false)
const previewTaskId = ref('')

// 统计数据
const stats = reactive({
  totalTasks: 0,
  completedTasks: 0,
  processingTasks: 0,
  failedTasks: 0,
  pendingTasks: 0,
})

const getStatusName = (status) => {
  switch (status) {
    case 'pending':
      return '排队中'
    case 'processing':
      return '生成中'
    case 'completed':
      return '已完成'
    case 'failed':
      return '失败'
    default:
      return status || '未知'
  }
}

const getStatusColor = (status) => {
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

const fetchStats = () => {
  httpGet('/api/admin/ppt/stats')
    .then((res) => {
      Object.assign(stats, res.data || {})
    })
    .catch((e) => {
      ElMessage.error('获取统计信息失败：' + e.message)
    })
}

const fetchJobs = () => {
  loading.value = true
  const params = {
    page: pagination.page,
    page_size: pagination.size,
  }
  if (queryForm.user_id) {
    params.user_id = queryForm.user_id
  }
  if (queryForm.status) {
    params.status = queryForm.status
  }

  httpGet('/api/admin/ppt/jobs', params)
    .then((res) => {
      const data = res.data || {}
      taskList.value = data.jobs || []
      pagination.total = data.total || 0
    })
    .catch((e) => {
      ElMessage.error('获取任务列表失败：' + e.message)
    })
    .finally(() => {
      loading.value = false
    })
}

const safeExportFileBase = (title, taskId) => {
  const raw = (title && String(title).trim()) || taskId || 'export'
  return String(raw)
    .replace(/[/\\:*?"<>|]/g, '_')
    .slice(0, 120)
}

const runExportDownload = async (taskId, format, baseName) => {
  const ext = format === 'pdf' ? '.pdf' : '.pptx'
  exportLoadingKey.value = taskId
  try {
    const response = await httpDownload(
      `/api/admin/ppt/jobs/${encodeURIComponent(taskId)}/export?format=${encodeURIComponent(format)}`
    )
    if (response.status !== 200) {
      throw new Error('导出失败')
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

const onExportCommand = (cmd, job) => {
  if (cmd !== 'pdf' && cmd !== 'pptx') return
  runExportDownload(job.task_id, cmd, job.title || job.task_id)
}

const openPreview = (taskId) => {
  previewTaskId.value = taskId
  previewVisible.value = true
}

const onPreviewClosed = () => {
  previewTaskId.value = ''
}

const handleQuery = () => {
  pagination.page = 1
  fetchJobs()
  fetchStats()
}

const handleSizeChange = (size) => {
  pagination.size = size
  fetchJobs()
}

const handleCurrentChange = (page) => {
  pagination.page = page
  fetchJobs()
}

onMounted(() => {
  fetchStats()
  fetchJobs()
})
</script>

<style lang="scss" scoped>
@use '@/assets/css/main.scss' as *;

.stats-row {
  margin-bottom: 16px;
}

.stat-card {
  .stat-item {
    text-align: center;
    padding: 10px 0;
  }
  .stat-number {
    font-size: 20px;
    font-weight: 600;
    &.info {
      color: #409eff;
    }
    &.warning {
      color: #e6a23c;
    }
    &.success {
      color: #67c23a;
    }
    &.danger {
      color: #f56c6c;
    }
  }
  .stat-label {
    font-size: 12px;
    color: #909399;
  }
}

.filter-card {
  margin-bottom: 16px;
}

.table-card {
  .pagination-container {
    margin-top: 16px;
    display: flex;
    justify-content: flex-end;
  }
}
</style>
