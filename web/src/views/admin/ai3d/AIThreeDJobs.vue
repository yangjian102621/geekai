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
            <el-option label="已完成" value="completed" />
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
              <div class="stat-number">{{ stats.completed }}</div>
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
      <el-table :data="taskList" v-loading="loading" stripe border style="width: 100%">
        <el-table-column prop="user_id" label="用户ID" />
        <el-table-column prop="type" label="平台">
          <template #default="{ row }">
            <el-tag :type="row.type === 'gitee' ? 'success' : 'primary'">
              {{ row.type === 'gitee' ? '魔力方舟' : '腾讯混元' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="model" label="模型格式" />
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
            {{ formatTime(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column prop="updated_at" label="更新时间">
          <template #default="{ row }">
            {{ formatTime(row.updated_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="viewTask(row)">查看</el-button>
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
          <el-descriptions-item label="模型格式">{{ currentTask.model }}</el-descriptions-item>
          <el-descriptions-item label="算力消耗">{{ currentTask.power }}</el-descriptions-item>
          <el-descriptions-item label="任务状态">
            <el-tag :type="getStatusType(currentTask.status)">
              {{ getStatusText(currentTask.status) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="创建时间">{{
            formatTime(currentTask.created_at)
          }}</el-descriptions-item>
          <el-descriptions-item label="更新时间">{{
            formatTime(currentTask.updated_at)
          }}</el-descriptions-item>
        </el-descriptions>

        <div class="task-params">
          <h4>任务参数</h4>
          <el-input v-model="taskParamsDisplay" type="textarea" :rows="6" readonly />
        </div>

        <div v-if="currentTask.img_url" class="task-result">
          <h4>生成结果</h4>
          <div class="result-links">
            <el-button type="primary" @click="downloadModel(currentTask)"> 下载3D模型 </el-button>
            <el-button v-if="currentTask.preview_url" @click="viewPreview(currentTask.preview_url)">
              查看预览
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

    <!-- 预览图片弹窗 -->
    <el-dialog v-model="previewVisible" title="预览图片" width="50%">
      <div class="preview-container">
        <el-image :src="previewUrl" fit="contain" style="width: 100%; height: 400px" />
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { httpGet } from '@/utils/http'
import { ElMessage, ElMessageBox } from 'element-plus'
import { computed, onMounted, reactive, ref } from 'vue'

// 响应式数据
const loading = ref(false)
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)
const taskList = ref([])
const taskDetailVisible = ref(false)
const previewVisible = ref(false)
const currentTask = ref(null)
const previewUrl = ref('')

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

// 计算属性
const taskParamsDisplay = computed(() => {
  if (!currentTask.value?.params) return '无参数'

  try {
    const params = JSON.parse(currentTask.value.params)
    return JSON.stringify(params, null, 2)
  } catch {
    return currentTask.value.params
  }
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
      taskList.value = response.data.list
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

const downloadModel = (task) => {
  if (task.img_url) {
    const link = document.createElement('a')
    link.href = task.img_url
    link.download = `3d_model_${task.id}.${task.model}`
    link.style.display = 'none'
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    ElMessage.success('开始下载3D模型')
  } else {
    ElMessage.warning('模型文件不存在')
  }
}

const viewPreview = (url) => {
  previewUrl.value = url
  previewVisible.value = true
}

const getStatusType = (status) => {
  const typeMap = {
    pending: 'warning',
    processing: 'primary',
    completed: 'success',
    failed: 'danger',
  }
  return typeMap[status] || 'info'
}

const getStatusText = (status) => {
  const textMap = {
    pending: '等待中',
    processing: '处理中',
    completed: '已完成',
    failed: '失败',
  }
  return textMap[status] || status
}

const getProgressStatus = (status) => {
  if (status === 'failed') return 'exception'
  if (status === 'completed') return 'success'
  return ''
}

const formatTime = (timestamp) => {
  if (!timestamp) return '-'
  const date = new Date(timestamp * 1000)
  return date.toLocaleString()
}

// 生命周期
onMounted(() => {
  loadData()
  loadStats()
})
</script>

<style lang="scss" scoped>
.admin-threed-jobs {
  padding: 20px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;

  h2 {
    margin: 0;
    color: #333;
  }
}

.search-section {
  background: white;
  padding: 20px;
  border-radius: 8px;
  margin-bottom: 20px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);

  .el-form-item {
    margin-bottom: 0;
    .el-select__wrapper {
      height: 36px;
      line-height: 36px;
    }
  }
}

.stats-section {
  margin-bottom: 20px;

  .stat-card {
    background: white;
    padding: 20px;
    border-radius: 8px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    display: flex;
    align-items: center;
    gap: 16px;

    .stat-icon {
      width: 48px;
      height: 48px;
      border-radius: 50%;
      display: flex;
      align-items: center;
      justify-content: center;

      i {
        font-size: 24px;
        color: white;
      }

      &.pending {
        background: #e6a23c;
      }

      &.processing {
        background: #409eff;
      }

      &.completed {
        background: #67c23a;
      }

      &.failed {
        background: #f56c6c;
      }
    }

    .stat-content {
      .stat-number {
        font-size: 24px;
        font-weight: bold;
        color: #333;
        margin-bottom: 4px;
      }

      .stat-label {
        font-size: 14px;
        color: #666;
      }
    }
  }
}

.table-section {
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  overflow: hidden;
}

.pagination-section {
  padding: 20px;
  text-align: center;
}

.task-detail {
  .task-params,
  .task-result,
  .task-error {
    margin-top: 20px;

    h4 {
      margin: 0 0 12px 0;
      color: #333;
      font-size: 16px;
    }
  }

  .result-links {
    display: flex;
    gap: 12px;
  }
}

.preview-container {
  text-align: center;
}
</style>
