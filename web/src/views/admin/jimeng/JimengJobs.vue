<template>
  <div class="app-container">
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
            <div class="stat-number !text-blue-500">{{ stats.pendingTasks }}</div>
            <div class="stat-label">排队中</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="4">
        <el-card class="stat-card">
          <div class="stat-item">
            <div class="stat-number warning">{{ stats.processingTasks }}</div>
            <div class="stat-label">处理中</div>
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
      <el-form :model="queryForm" ref="queryFormRef" :inline="true" label-width="80px">
        <el-form-item label="用户ID">
          <el-input
            v-model="queryForm.user_id"
            placeholder="请输入用户ID"
            clearable
            style="width: 150px"
          />
        </el-form-item>
        <el-form-item label="任务类型">
          <el-select
            v-model="queryForm.type"
            placeholder="请选择任务类型"
            clearable
            style="width: 150px"
            @change="handleQuery"
          >
            <el-option label="文生图" value="text_to_image" />
            <el-option label="图生图" value="image_to_image" />
            <el-option label="图像编辑" value="image_edit" />
            <el-option label="图像特效" value="image_effects" />
            <el-option label="文生视频" value="text_to_video" />
            <el-option label="图生视频" value="image_to_video" />
          </el-select>
        </el-form-item>
        <el-form-item label="任务状态">
          <el-select
            v-model="queryForm.status"
            placeholder="请选择状态"
            clearable
            style="width: 120px"
            @change="handleQuery"
          >
            <el-option label="等待中" value="in_queue" />
            <el-option label="处理中" value="generating" />
            <el-option label="已完成" value="success" />
            <el-option label="失败" value="failed" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleQuery" :loading="loading">
            <i class="iconfont icon-search mr-1" />
            搜索
          </el-button>
          <el-button type="danger" @click="handleBatchDelete" :disabled="!multipleSelection.length">
            <i class="iconfont icon-remove mr-1" />
            批量删除
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 任务列表 -->
    <el-card class="table-card">
      <el-table
        :data="taskList"
        v-loading="loading"
        @selection-change="handleSelectionChange"
        border
      >
        <el-table-column type="selection" width="55" />
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="user_id" label="用户ID" width="80" />
        <el-table-column prop="type" label="任务类型" width="120">
          <template #default="scope">
            <el-tag size="small">{{ getTaskTypeName(scope.row.type) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="prompt" label="提示词" min-width="200" show-overflow-tooltip />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="scope">
            <el-tag :type="getStatusColor(scope.row.status)" size="small">
              {{ getStatusName(scope.row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="progress" label="进度" width="100">
          <template #default="scope">
            <el-progress :percentage="scope.row.progress" :stroke-width="4" />
          </template>
        </el-table-column>
        <el-table-column prop="power" label="算力" width="80" />
        <el-table-column prop="created_at" label="创建时间" width="180">
          <template #default="scope">
            {{ formatDateTime(scope.row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="scope">
            <el-button type="primary" size="small" text @click="handleViewDetail(scope.row)">
              详情
            </el-button>
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

    <!-- 任务详情对话框 -->
    <el-dialog
      v-model="detailDialog.visible"
      :title="`任务详情 - ${detailDialog.data.id}`"
      width="800px"
      :close-on-click-modal="false"
    >
      <div class="detail-content" v-if="detailDialog.data">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="任务ID">{{ detailDialog.data.id }}</el-descriptions-item>
          <el-descriptions-item label="用户ID">{{
            detailDialog.data.user_id
          }}</el-descriptions-item>
          <el-descriptions-item label="任务类型">{{
            getTaskTypeName(detailDialog.data.type)
          }}</el-descriptions-item>
          <el-descriptions-item label="状态">
            <el-tag :type="getStatusColor(detailDialog.data.status)">
              {{ getStatusName(detailDialog.data.status) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="进度"
            >{{ detailDialog.data.progress }}%</el-descriptions-item
          >
          <el-descriptions-item label="算力消耗">{{
            detailDialog.data.power
          }}</el-descriptions-item>
          <el-descriptions-item label="创建时间">{{
            formatDateTime(detailDialog.data.created_at)
          }}</el-descriptions-item>
          <el-descriptions-item label="更新时间">{{
            formatDateTime(detailDialog.data.updated_at)
          }}</el-descriptions-item>
        </el-descriptions>

        <div class="detail-section">
          <h4 class="text-base pt-2 font-bold">提示词</h4>
          <div class="prompt-content">{{ detailDialog.data.prompt || '无' }}</div>
        </div>

        <div class="detail-section" v-if="detailDialog.data.task_params">
          <h4>任务参数</h4>
          <el-input
            v-model="detailDialog.data.task_params"
            type="textarea"
            :rows="5"
            readonly
            class="params-content"
          />
        </div>

        <div class="detail-section" v-if="detailDialog.data.err_msg">
          <h4>错误信息</h4>
          <el-alert :title="detailDialog.data.err_msg" type="error" :closable="false" />
        </div>

        <div class="detail-section" v-if="detailDialog.data.img_url || detailDialog.data.video_url">
          <h4>生成结果</h4>
          <div class="result-content">
            <div v-if="detailDialog.data.img_url" class="result-item">
              <label>图片：</label>
              <el-image
                :src="detailDialog.data.img_url"
                :preview-src-list="[detailDialog.data.img_url]"
                fit="cover"
                style="width: 100px; height: 100px; border-radius: 4px"
              />
            </div>
            <div v-if="detailDialog.data.video_url" class="result-item">
              <label>视频：</label>
              <video
                :src="detailDialog.data.video_url"
                controls
                style="width: 200px; height: 150px; border-radius: 4px"
              />
            </div>
          </div>
        </div>

        <div class="detail-section" v-if="detailDialog.data.raw_data">
          <h4>原始响应数据</h4>
          <el-input
            v-model="formattedRawData"
            type="textarea"
            :rows="10"
            readonly
            class="raw-data-content"
          />
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { httpGet, httpPost } from '@/utils/http'
import { formatDateTime } from '@/utils/libs'
import { ElMessage, ElMessageBox } from 'element-plus'
import { computed, onMounted, reactive, ref } from 'vue'

// 查询表单
const queryForm = reactive({
  user_id: '',
  type: '',
  status: '',
})

// 分页信息
const pagination = reactive({
  page: 1,
  size: 20,
  total: 0,
})

// 数据
const taskList = ref([])
const loading = ref(false)
const multipleSelection = ref([])
const queryFormRef = ref(null)

// 统计信息
const stats = reactive({
  totalTasks: 0,
  completedTasks: 0,
  processingTasks: 0,
  failedTasks: 0,
})

// 详情对话框
const detailDialog = reactive({
  visible: false,
  data: {},
})

// 格式化原始数据
const formattedRawData = computed(() => {
  if (!detailDialog.data.raw_data) return ''
  try {
    return JSON.stringify(JSON.parse(detailDialog.data.raw_data), null, 2)
  } catch (error) {
    return detailDialog.data.raw_data
  }
})

// 获取任务类型名称
const getTaskTypeName = (type) => {
  const typeMap = {
    text_to_image: '文生图',
    image_to_image: '图生图',
    image_edit: '图像编辑',
    image_effects: '图像特效',
    text_to_video: '文生视频',
    image_to_video: '图生视频',
  }
  return typeMap[type] || type
}

// 获取状态名称
const getStatusName = (status) => {
  const statusMap = {
    in_queue: '等待中',
    generating: '处理中',
    success: '已完成',
    failed: '失败',
  }
  return statusMap[status] || status
}

// 获取状态颜色
const getStatusColor = (status) => {
  const colorMap = {
    in_queue: '',
    generating: 'warning',
    success: 'success',
    failed: 'danger',
  }
  return colorMap[status] || ''
}

// 获取任务列表
const getTaskList = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.page,
      page_size: pagination.size,
      ...queryForm,
    }

    const response = await httpGet('/api/admin/jimeng/jobs', params)
    taskList.value = response.data.jobs || []
    pagination.total = response.data.total || 0
  } catch (error) {
    ElMessage.error('获取任务列表失败')
  } finally {
    loading.value = false
  }
}

// 获取统计信息
const getStats = async () => {
  try {
    const response = await httpGet('/api/admin/jimeng/stats')
    Object.assign(stats, response.data)
  } catch (error) {
    console.error('获取统计信息失败:', error)
  }
}

// 查询
const handleQuery = () => {
  pagination.page = 1
  getTaskList()
}

// 选择变化
const handleSelectionChange = (selection) => {
  multipleSelection.value = selection
}

// 查看详情
const handleViewDetail = async (row) => {
  try {
    const response = await httpGet(`/api/admin/jimeng/jobs/${row.id}`)
    detailDialog.data = response.data
    detailDialog.visible = true
  } catch (error) {
    ElMessage.error('获取任务详情失败')
  }
}

// 批量删除
const handleBatchDelete = async () => {
  if (!multipleSelection.value.length) {
    ElMessage.warning('请选择要删除的任务')
    return
  }

  try {
    await ElMessageBox.confirm(
      `确定要删除选中的 ${multipleSelection.value.length} 个任务吗？`,
      '提示',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )

    const jobIds = multipleSelection.value.map((item) => item.id)
    await httpPost('/api/admin/jimeng/jobs/remove', { job_ids: jobIds })
    ElMessage.success('批量删除成功')
    getTaskList()
    getStats()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('批量删除失败')
    }
  }
}

// 分页大小变化
const handleSizeChange = (size) => {
  pagination.size = size
  pagination.page = 1
  getTaskList()
}

// 当前页变化
const handleCurrentChange = (page) => {
  pagination.page = page
  getTaskList()
}

// 初始化
onMounted(() => {
  getTaskList()
  getStats()
})
</script>

<style lang="scss">
.app-container {
  padding: 20px;

  .el-form-item {
    margin-bottom: 0;
  }
}

.page-header {
  margin-bottom: 20px;

  h2 {
    margin: 0 0 8px 0;
    color: #303133;
  }

  p {
    margin: 0;
    color: #606266;
    font-size: 14px;
  }
}

.filter-card {
  margin-bottom: 20px;
}

.stats-row {
  margin-bottom: 20px;
}

.stat-card {
  .stat-item {
    text-align: center;
    padding: 20px;

    .stat-number {
      font-size: 28px;
      font-weight: bold;
      color: #303133;
      margin-bottom: 8px;

      &.success {
        color: #67c23a;
      }

      &.warning {
        color: #e6a23c;
      }

      &.danger {
        color: #f56c6c;
      }
    }

    .stat-label {
      font-size: 14px;
      color: #909399;
    }
  }
}

.table-card {
  .pagination-container {
    margin-top: 20px;
    display: flex;
    justify-content: center;
  }
}

.detail-content {
  .detail-section {
    margin-bottom: 20px;

    h4 {
      margin: 0 0 10px 0;
      color: #303133;
      font-size: 16px;
    }

    .prompt-content {
      background: #f5f7fa;
      padding: 12px;
      border-radius: 4px;
      color: #606266;
      line-height: 1.6;
    }

    .params-content,
    .raw-data-content {
      font-family: monospace;
    }

    .result-content {
      .result-item {
        margin-bottom: 10px;
        display: flex;
        align-items: center;
        gap: 10px;

        label {
          font-weight: bold;
          color: #303133;
          min-width: 50px;
        }
      }
    }
  }
}
</style>
