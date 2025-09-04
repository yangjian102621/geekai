<template>
  <div class="admin-threed-jobs">
    <!-- 搜索和筛选 -->
    <div class="search-section">
      <el-form :model="searchForm" inline>
        <el-form-item label="任务状态">
          <el-select
            v-model="searchForm.status"
            placeholder="选择状态"
            style="width: 120px"
            clearable
          >
            <el-option label="全部" value="" />
            <el-option label="等待中" value="pending" />
            <el-option label="处理中" value="processing" />
            <el-option label="已完成" value="success" />
            <el-option label="失败" value="failed" />
          </el-select>
        </el-form-item>

        <el-form-item label="平台类型">
          <el-select
            v-model="searchForm.type"
            placeholder="选择平台"
            style="width: 120px"
            clearable
          >
            <el-option label="全部" value="" />
            <el-option label="魔力方舟" value="gitee" />
            <el-option label="腾讯混元" value="tencent" />
          </el-select>
        </el-form-item>

        <el-form-item label="用户ID">
          <el-input v-model="searchForm.userId" placeholder="输入用户ID" clearable />
        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="resetSearch">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <!-- 数据统计 -->
    <div class="stats-section">
      <el-row :gutter="20">
        <el-col :span="6">
          <div class="stat-card">
            <div class="stat-icon pending">
              <i class="iconfont icon-clock"></i>
            </div>
            <div class="stat-content">
              <div class="stat-number">{{ stats.pending }}</div>
              <div class="stat-label">等待中</div>
            </div>
          </div>
        </el-col>
        <el-col :span="6">
          <div class="stat-card">
            <div class="stat-icon processing">
              <i class="iconfont icon-loading"></i>
            </div>
            <div class="stat-content">
              <div class="stat-number">{{ stats.processing }}</div>
              <div class="stat-label">处理中</div>
            </div>
          </div>
        </el-col>
        <el-col :span="6">
          <div class="stat-card">
            <div class="stat-icon completed">
              <i class="iconfont icon-check"></i>
            </div>
            <div class="stat-content">
              <div class="stat-number">{{ stats.success }}</div>
              <div class="stat-label">已完成</div>
            </div>
          </div>
        </el-col>
        <el-col :span="6">
          <div class="stat-card">
            <div class="stat-icon failed">
              <i class="iconfont icon-error"></i>
            </div>
            <div class="stat-content">
              <div class="stat-number">{{ stats.failed }}</div>
              <div class="stat-label">失败</div>
            </div>
          </div>
        </el-col>
      </el-row>
    </div>

    <!-- 任务列表 -->
    <div class="table-section w-full">
      <el-table :data="taskList" v-loading="loading" border style="width: 100%">
        <el-table-column prop="user_id" label="用户ID" width="80" />
        <el-table-column prop="type" label="平台">
          <template #default="{ row }">
            <el-tag :type="row.type === 'gitee' ? 'success' : 'primary'">
              {{ row.type === 'gitee' ? '魔力方舟' : '腾讯混元' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="model" label="模型名称" />
        <el-table-column label="模型格式">
          <template #default="{ row }">
            {{ row.params.file_format }}
          </template>
        </el-table-column>
        <el-table-column prop="power" label="算力消耗" />
        <el-table-column prop="status" label="状态">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间">
          <template #default="{ row }">
            {{ dateFormat(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column prop="updated_at" label="更新时间">
          <template #default="{ row }">
            {{ dateFormat(row.updated_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="300" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="viewTask(row)">查看</el-button>
            <el-button
              size="small"
              type="primary"
              plain
              v-if="row.status === 'success'"
              @click="openModelPreview(row)"
            >
              预览模型
            </el-button>
            <el-button size="small" type="danger" @click="deleteTask(row.id)"> 删除 </el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-section">
        <el-pagination
          :current-page="currentPage"
          :page-size="pageSize"
          :page-sizes="[20, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </div>

    <!-- 任务详情弹窗 -->
    <el-dialog
      v-model="taskDetailVisible"
      title="任务详情"
      width="60%"
      :before-close="closeTaskDetail"
    >
      <div v-if="currentTask" class="task-detail">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="任务ID">{{ currentTask.id }}</el-descriptions-item>
          <el-descriptions-item label="用户ID">{{ currentTask.user_id }}</el-descriptions-item>
          <el-descriptions-item label="平台类型">
            <el-tag :type="currentTask.type === 'gitee' ? 'success' : 'primary'">
              {{ currentTask.type === 'gitee' ? '魔力方舟' : '腾讯混元' }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="模型名称">{{ currentTask.model }}</el-descriptions-item>
          <el-descriptions-item label="算力消耗">{{ currentTask.power }}</el-descriptions-item>
          <el-descriptions-item label="任务状态">
            <el-tag :type="getStatusType(currentTask.status)">
              {{ getStatusText(currentTask.status) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="创建时间">{{
            dateFormat(currentTask.created_at)
          }}</el-descriptions-item>
          <el-descriptions-item label="更新时间">{{
            dateFormat(currentTask.updated_at)
          }}</el-descriptions-item>
        </el-descriptions>

        <div class="task-params">
          <h4>任务参数</h4>
          <div class="params-content">
            <pre>{{ JSON.stringify(currentTask.params, null, 2) }}</pre>
          </div>
        </div>

        <div v-if="currentTask.img_url || currentTask.file_url" class="task-result">
          <h4>生成结果</h4>
          <div class="result-links">
            <el-button type="primary" @click="downloadModel(currentTask)"> 下载3D模型 </el-button>
            <el-button
              v-if="currentTask.file_url"
              type="success"
              plain
              @click="openModelPreview(currentTask)"
            >
              预览模型
            </el-button>
          </div>
        </div>

        <div v-if="currentTask.err_msg" class="task-error">
          <h4>错误信息</h4>
          <el-alert :title="currentTask.err_msg" type="error" :closable="false" show-icon />
        </div>
      </div>

      <template #footer>
        <span class="dialog-footer">
          <el-button @click="closeTaskDetail">关闭</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 3D 模型预览弹窗 -->
    <el-dialog
      v-model="modelPreviewVisible"
      :class="['model-preview-dialog', { dark: isDarkTheme }]"
      title="模型预览"
      fullscreen
      destroy-on-close
    >
      <div class="model-preview-wrapper">
        <ThreeDPreview :model-url="modelPreviewUrl" />
      </div>
      <template #footer>
        <span class="dialog-footer">
          <el-button
            type="primary"
            @click="downloadModel(currentTask)"
            :loading="currentTask.downloading"
          >
            下载3D模型
          </el-button>
          <el-button @click="modelPreviewVisible = false">关闭</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import ThreeDPreview from '@/components/ThreeDPreview.vue'
import { showMessageError } from '@/utils/dialog'
import { httpDownload, httpGet } from '@/utils/http'
import { dateFormat, replaceImg } from '@/utils/libs'
import { ElMessage, ElMessageBox } from 'element-plus'
import { onMounted, reactive, ref } from 'vue'

// 响应式数据
const loading = ref(false)
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)
const taskList = ref([])
const taskDetailVisible = ref(false)
const currentTask = ref({
  downloading: false,
})
const previewUrl = ref('')
// 3D 预览
const modelPreviewVisible = ref(false)
const modelPreviewUrl = ref('')
// 简单检测暗色主题（若全局有主题管理可替换）
const isDarkTheme = ref(
  document.documentElement.classList.contains('dark') || document.body.classList.contains('dark')
)

// 搜索表单
const searchForm = reactive({
  status: '',
  type: '',
  userId: '',
})

// 统计数据
const stats = reactive({
  pending: 0,
  processing: 0,
  completed: 0,
  failed: 0,
})

// 方法
const loadData = async () => {
  try {
    loading.value = true

    const params = {
      page: currentPage.value,
      page_size: pageSize.value,
      ...searchForm,
    }

    // 移除空值
    Object.keys(params).forEach((key) => {
      if (params[key] === '') {
        delete params[key]
      }
    })

    const response = await httpGet('/api/admin/ai3d/jobs', params)

    if (response.code === 0) {
      taskList.value = response.data.items
      total.value = response.data.total
    } else {
      ElMessage.error(response.message || '加载数据失败')
    }
  } catch (error) {
    ElMessage.error('加载数据失败：' + error.message)
  } finally {
    loading.value = false
  }
}

const loadStats = async () => {
  try {
    const response = await httpGet('/api/admin/ai3d/stats')
    if (response.code === 0) {
      Object.assign(stats, response.data)
    }
  } catch (error) {
    console.error('加载统计数据失败:', error)
  }
}

const handleSearch = () => {
  currentPage.value = 1
  loadData()
}

const resetSearch = () => {
  Object.assign(searchForm, {
    status: '',
    type: '',
    userId: '',
  })
  currentPage.value = 1
  loadData()
}

const handleSizeChange = (size) => {
  pageSize.value = size
  currentPage.value = 1
  loadData()
}

const handleCurrentChange = (page) => {
  currentPage.value = page
  loadData()
}

const refreshData = () => {
  loadData()
  loadStats()
}

const viewTask = (task) => {
  currentTask.value = task
  taskDetailVisible.value = true
}

const closeTaskDetail = () => {
  taskDetailVisible.value = false
  currentTask.value = null
}

const deleteTask = async (taskId) => {
  try {
    await ElMessageBox.confirm('确定要删除这个任务吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    })

    const response = await httpGet(`/api/admin/ai3d/jobs/${taskId}/delete`)

    if (response.code === 0) {
      ElMessage.success('删除成功')
      loadData()
      loadStats()
    } else {
      ElMessage.error(response.message || '删除失败')
    }
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败：' + error.message)
    }
  }
}

const downloadModel = async (task) => {
  const url = replaceImg(task.file_url)
  const downloadURL = `/api/download?url=${url}`
  const urlObj = new URL(url)
  const fileName = urlObj.pathname.split('/').pop()
  task.downloading = true
  try {
    const response = await httpDownload(downloadURL)
    const blob = new Blob([response.data])
    const link = document.createElement('a')
    link.href = URL.createObjectURL(blob)
    link.download = fileName
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    URL.revokeObjectURL(link.href)
    task.downloading = false
  } catch (error) {
    showMessageError('下载失败:' + error.message)
    task.downloading = false
  }
}

const openModelPreview = (task) => {
  // 优先使用文件直链，后端下载代理也可拼接
  const url = task.file_url
  if (!url) {
    ElMessage.warning('暂无可预览的模型文件')
    return
  }
  currentTask.value = task
  modelPreviewUrl.value = url
  modelPreviewVisible.value = true
}

const getStatusType = (status) => {
  const typeMap = {
    pending: 'warning',
    processing: 'primary',
    success: 'success',
    failed: 'danger',
  }
  return typeMap[status] || 'info'
}

const getStatusText = (status) => {
  const textMap = {
    pending: '等待中',
    processing: '处理中',
    success: '已完成',
    failed: '失败',
  }
  return textMap[status] || status
}

// 生命周期
onMounted(() => {
  loadData()
  loadStats()
})
</script>

<style lang="scss" scoped>
@use '@/assets/css/admin/ai3d.scss' as *;
</style>
