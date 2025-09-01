<template>
  <div class="moderation-list container p-5">
    <!-- 搜索筛选区域 -->
    <el-card shadow="never" class="mb-6">
      <el-form :model="queryForm" label-position="top" class="search-form">
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          <el-form-item label="用户名">
            <el-input
              v-model="queryForm.username"
              placeholder="请输入用户名"
              clearable
              @keyup.enter="handleSearch"
            />
          </el-form-item>

          <el-form-item label="来源">
            <el-select
              v-model="queryForm.source"
              placeholder="请选择来源"
              clearable
              style="width: 100%"
            >
              <el-option
                v-for="item in sourceList"
                :key="item.id"
                :label="item.name"
                :value="item.id"
              />
            </el-select>
          </el-form-item>

          <el-form-item label="时间范围">
            <el-date-picker
              v-model="queryForm.dateRange"
              type="daterange"
              range-separator="至"
              start-placeholder="开始日期"
              end-placeholder="结束日期"
              format="YYYY-MM-DD"
              value-format="YYYY-MM-DD"
              style="width: 100%"
            />
          </el-form-item>
        </div>

        <div class="flex justify-between items-center mt-4">
          <div class="flex space-x-2">
            <el-button type="primary" @click="handleSearch" :loading="loading">
              <i class="iconfont icon-search mr-1"></i>
              搜索
            </el-button>
            <el-button @click="handleReset">
              <i class="iconfont icon-refresh mr-1"></i>
              重置
            </el-button>
          </div>

          <div class="text-sm text-gray-500">
            共找到 <span class="font-semibold text-blue-600">{{ total }}</span> 条记录
          </div>
        </div>
      </el-form>
    </el-card>

    <!-- 数据列表 -->
    <el-card shadow="never">
      <div class="table-header flex justify-between items-center mb-4">
        <h3 class="text-lg font-semibold text-gray-700">审核记录列表</h3>
        <div class="flex space-x-2">
          <el-button
            type="danger"
            size="small"
            @click="handleBatchDelete"
            :disabled="selectedRows.length === 0"
          >
            <i class="iconfont icon-delete mr-1"></i>
            批量删除 ({{ selectedRows.length }})
          </el-button>
        </div>
      </div>

      <el-table
        :data="tableData"
        v-loading="loading"
        @selection-change="handleSelectionChange"
        stripe
        border
        style="width: 100%"
      >
        <el-table-column type="selection" width="55" />

        <el-table-column prop="id" label="ID" width="80" align="center" />

        <el-table-column prop="username" label="用户名" width="120">
          <template #default="{ row }">
            <span class="font-medium text-gray-700">{{ row.username || '未知用户' }}</span>
          </template>
        </el-table-column>

        <el-table-column prop="source" label="来源" width="140" align="center">
          <template #default="{ row }">
            <el-tag type="primary" size="small">
              {{ getSourceLabel(row.source) }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column prop="input" label="用户输入" max-width="300">
          <template #default="{ row }">
            <div class="text-content">
              <span class="text-display">
                {{
                  row.input
                    ? row.input.length > 30
                      ? row.input.substring(0, 30) + '...'
                      : row.input
                    : '-'
                }}
              </span>
            </div>
          </template>
        </el-table-column>

        <el-table-column prop="output" label="AI 输出" max-width="300">
          <template #default="{ row }">
            <div class="text-content">
              <span class="text-display">
                {{
                  row.output
                    ? row.output.length > 30
                      ? row.output.substring(0, 30) + '...'
                      : row.output
                    : '-'
                }}
              </span>
            </div>
          </template>
        </el-table-column>

        <el-table-column prop="created_at" label="创建时间" width="180" align="center">
          <template #default="{ row }">
            <span class="text-gray-600">{{ dateFormat(row.created_at) }}</span>
          </template>
        </el-table-column>

        <el-table-column label="操作" width="150" align="center" fixed="right">
          <template #default="{ row }">
            <div class="flex space-x-2">
              <el-button type="primary" size="small" @click="handleView(row)">
                <i class="iconfont icon-view mr-1"></i>
                查看
              </el-button>
              <el-button type="danger" size="small" @click="handleDelete(row)">
                <i class="iconfont icon-delete mr-1"></i>
                删除
              </el-button>
            </div>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-wrapper flex justify-center mt-6">
        <el-pagination
          :current-page="currentPage"
          :page-size="pageSize"
          :page-sizes="[15, 30, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>

    <!-- 查看详情弹窗 -->
    <el-dialog
      v-model="detailDialogVisible"
      title="审核记录详情"
      width="800px"
      :close-on-click-modal="false"
    >
      <div v-if="currentRecord" class="record-detail">
        <div class="grid grid-cols-1 gap-4">
          <div class="detail-item">
            <label class="block text-sm font-medium text-gray-700 mb-2">用户信息</label>
            <div class="bg-gray-50 p-3 rounded">
              <span class="text-gray-900">{{ currentRecord.username || '未知用户' }}</span>
              <span class="text-gray-500 ml-4">ID: {{ currentRecord.user_id }}</span>
            </div>
          </div>

          <div class="detail-item">
            <label class="block text-sm font-medium text-gray-700 mb-2">来源</label>
            <div class="bg-gray-50 p-3 rounded">
              <el-tag type="primary" size="small">
                {{ getSourceLabel(currentRecord.source) }}
              </el-tag>
            </div>
          </div>

          <div class="detail-item">
            <label class="block text-sm font-medium text-gray-700 mb-2">用户输入</label>
            <div class="bg-gray-50 p-3 rounded">
              <pre class="whitespace-pre-wrap text-sm text-gray-900">{{ currentRecord.input }}</pre>
            </div>
          </div>

          <div class="detail-item">
            <label class="block text-sm font-medium text-gray-700 mb-2">AI 输出</label>
            <div class="bg-gray-50 p-3 rounded">
              <pre class="whitespace-pre-wrap text-sm text-gray-900">{{
                currentRecord.output
              }}</pre>
            </div>
          </div>

          <div class="detail-item">
            <label class="block text-sm font-medium text-gray-700 mb-2">审核结果</label>
            <div class="bg-gray-50 p-3 rounded">
              <div class="flex flex-col space-y-2">
                <el-tag
                  type="primary"
                  size="small"
                  v-for="item in currentRecord.result"
                  :key="item"
                >
                  {{ item }}
                </el-tag>
              </div>
            </div>
          </div>

          <div class="detail-item">
            <label class="block text-sm font-medium text-gray-700 mb-2">创建时间</label>
            <div class="bg-gray-50 p-3 rounded">
              <span class="text-gray-900">{{ dateFormat(currentRecord.created_at) }}</span>
            </div>
          </div>
        </div>
      </div>

      <template #footer>
        <div class="flex justify-end space-x-2">
          <el-button @click="detailDialogVisible = false">关闭</el-button>
          <el-button type="danger" @click="handleDelete(currentRecord)"> 删除记录 </el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { httpGet, httpPost } from '@/utils/http'
import { dateFormat } from '@/utils/libs'
import { ElMessage, ElMessageBox } from 'element-plus'
import { computed, onMounted, ref } from 'vue'

// 响应式数据
const loading = ref(false)
const tableData = ref([])
const selectedRows = ref([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(15)
const detailDialogVisible = ref(false)
const currentRecord = ref(null)
const sourceList = ref([])

// 查询表单
const queryForm = ref({
  username: '',
  source: '',
  dateRange: [],
})

// 计算属性
const hasFilters = computed(() => {
  return queryForm.value.username || queryForm.value.source || queryForm.value.dateRange.length > 0
})

// 生命周期
onMounted(() => {
  fetchData()
  fetchSourceList()
})

// 获取来源列表
const fetchSourceList = async () => {
  try {
    const response = await httpGet('/api/admin/moderation/source-list')
    sourceList.value = response.data
  } catch (error) {
    console.error('获取来源列表失败:', error)
  }
}

// 获取数据
const fetchData = async () => {
  loading.value = true
  try {
    const params = {
      page: currentPage.value,
      page_size: pageSize.value,
      ...queryForm.value,
    }

    // 处理日期范围
    if (params.dateRange && params.dateRange.length === 2) {
      params.start_date = params.dateRange[0]
      params.end_date = params.dateRange[1]
      delete params.dateRange
    }

    const response = await httpPost('/api/admin/moderation/list', params)

    if (response.data) {
      tableData.value = response.data.items || []
      total.value = response.data.total || 0
      currentPage.value = response.data.page || 1
      pageSize.value = response.data.page_size || 15
    }
  } catch (error) {
    ElMessage.error('获取数据失败：' + (error.message || '未知错误'))
  } finally {
    loading.value = false
  }
}

// 搜索
const handleSearch = () => {
  currentPage.value = 1
  fetchData()
}

// 重置
const handleReset = () => {
  queryForm.value = {
    username: '',
    source: '',
    dateRange: [],
  }
  currentPage.value = 1

  fetchData()
}

// 分页处理
const handleSizeChange = (size) => {
  pageSize.value = size
  currentPage.value = 1
  fetchData()
}

const handleCurrentChange = (page) => {
  currentPage.value = page
  fetchData()
}

// 选择处理
const handleSelectionChange = (selection) => {
  selectedRows.value = selection
}

// 查看详情
const handleView = (row) => {
  currentRecord.value = row
  detailDialogVisible.value = true
}

// 删除单条记录
const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm('确定要删除这条审核记录吗？删除后无法恢复。', '确认删除', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    })

    await httpGet(`/api/admin/moderation/remove?id=${row.id}`)
    ElMessage.success('删除成功')
    fetchData()

    // 如果当前查看的记录被删除，关闭弹窗
    if (detailDialogVisible.value && currentRecord.value?.id === row.id) {
      detailDialogVisible.value = false
      currentRecord.value = null
    }
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败：' + (error.message || '未知错误'))
    }
  }
}

// 批量删除
const handleBatchDelete = async () => {
  if (selectedRows.value.length === 0) {
    ElMessage.warning('请选择要删除的记录')
    return
  }

  try {
    await ElMessageBox.confirm(
      `确定要删除选中的 ${selectedRows.value.length} 条审核记录吗？删除后无法恢复。`,
      '确认批量删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )

    const ids = selectedRows.value.map((row) => row.id)
    await httpPost('/api/admin/moderation/batch-remove', { ids })

    ElMessage.success('批量删除成功')
    selectedRows.value = []
    fetchData()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('批量删除失败：' + (error.message || '未知错误'))
    }
  }
}

const getSourceLabel = (source) => {
  const sourceMap = sourceList.value.find((item) => item.id === source)
  return sourceMap.name || source || '未知'
}
</script>

<style lang="scss">
.moderation-list {
  .page-header {
    border-bottom: 1px solid #e5e7eb;
    padding-bottom: 1rem;
  }

  .search-form {
    .el-form-item {
      margin-bottom: 0;
      .el-select__wrapper {
        height: 36px;
        line-height: 36px;
      }
    }
  }

  .table-header {
    border-bottom: 1px solid #f3f4f6;
    padding-bottom: 1rem;
  }

  .text-content {
    max-width: 300px;
    word-break: break-all;
  }

  .record-detail {
    .detail-item {
      label {
        color: #374151;
        font-weight: 500;
      }

      .bg-gray-50 {
        background-color: #f9fafb;
        border: 1px solid #e5e7eb;
      }
    }
  }

  .pagination-wrapper {
    padding: 1rem 0;
  }
}

// 响应式设计
@media (max-width: 768px) {
  .moderation-list {
    .container {
      padding: 1rem;
    }

    .search-form {
      .grid {
        grid-template-columns: 1fr;
      }
    }

    .table-header {
      flex-direction: column;
      align-items: flex-start;
      gap: 1rem;
    }
  }
}
</style>
