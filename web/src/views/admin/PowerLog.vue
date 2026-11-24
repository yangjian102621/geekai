<template>
  <div class="container power-log">
    <div class="handle-box">
      <div class="filter-box">
        <el-input v-model="query.model" placeholder="模型" class="filter-input" clearable @keyup.enter="search" />
        <el-input v-model="query.username" placeholder="用户名" class="filter-input" clearable @keyup.enter="search" />
        <el-input
          v-model.number="query.userid"
          placeholder="用户ID"
          class="filter-input input-uid"
          clearable
          @keyup.enter="search"
        />
        <el-select v-model="query.type" class="filter-select select-type">
          <el-option label="全部" :value="0" />
          <el-option label="充值" :value="1" />
          <el-option label="消费" :value="2" />
          <el-option label="退款" :value="3" />
          <el-option label="奖励" :value="4" />
          <el-option label="众筹" :value="5" />
          <el-option label="签到" :value="7" />
        </el-select>
        <el-select v-model="query.quota_status" placeholder="额度状态" class="filter-select select-quota" clearable>
          <el-option v-for="(label, value) in quotaStatuses" :key="value" :label="label" :value="value" />
        </el-select>
        <el-date-picker
          v-model="query.date"
          type="daterange"
          start-placeholder="开始日期"
          end-placeholder="结束日期"
          format="YYYY-MM-DD"
          value-format="YYYY-MM-DD"
          class="date-range"
        />
        <button class="query-btn" @click="search">
          <el-icon><Search /></el-icon>
          搜索
        </button>
      </div>

      <div class="toolbar-right" v-if="totalPower > 0">
        <span class="total-info">积分总额：{{ totalPower }}</span>
      </div>
    </div>

    <el-table
      v-if="items.length > 0"
      :data="items"
      class="log-table"
      :row-key="(row) => row.id"
    >
      <el-table-column prop="username" label="用户" min-width="140" />
      <el-table-column prop="model" label="模型" min-width="140" />
      <el-table-column label="类型" width="96">
        <template #default="scope">
          <span class="pill" :class="typePillClass(scope.row.type)">{{ scope.row.type_str }}</span>
        </template>
      </el-table-column>
      <el-table-column label="数额" width="110">
        <template #default="scope">
          <span v-if="scope.row.mark === 1" class="amount-in">+{{ scope.row.amount }}</span>
          <span v-else class="amount-out">-{{ scope.row.amount }}</span>
        </template>
      </el-table-column>
      <el-table-column prop="balance" label="余额" width="100" />
      <el-table-column label="Token 额度" width="130">
        <template #default="scope">
          {{ scope.row.quota_amount > 0 ? scope.row.quota_amount.toLocaleString() : '—' }}
        </template>
      </el-table-column>
      <el-table-column label="到账状态" width="120">
        <template #default="scope">
          <el-tooltip :content="scope.row.quota_error" :disabled="!scope.row.quota_error">
            <span
              v-if="scope.row.quota_status"
              class="pill"
              :class="quotaPillClass(scope.row.quota_status)"
            >
              {{ quotaStatuses[scope.row.quota_status] || '—' }}
            </span>
            <span v-else>—</span>
          </el-tooltip>
        </template>
      </el-table-column>
      <el-table-column label="发生时间" width="160">
        <template #default="scope">
          <span>{{ dateFormat(scope.row['created_at']) }}</span>
        </template>
      </el-table-column>
      <el-table-column prop="remark" label="备注" min-width="160">
        <template #default="scope">{{ scope.row.remark || '-' }}</template>
      </el-table-column>
    </el-table>

    <div class="table-footer" v-if="total > 0">
      <div class="footer-info">显示第 {{ startIndex }} - {{ endIndex }} 条，共 {{ total }} 条</div>
      <div class="footer-nav">
        <span class="footer-text">总页数: {{ totalPages }}</span>
        <el-pagination
          v-model:current-page="page"
          :page-size="pageSize"
          :total="total"
          layout="prev, pager, next"
          @current-change="fetchData"
        />
        <span class="footer-text">每页条数:</span>
        <el-select v-model="pageSize" class="page-size-select" @change="changePageSize">
          <el-option :value="10" label="10" />
          <el-option :value="20" label="20" />
          <el-option :value="50" label="50" />
        </el-select>
      </div>
    </div>

    <el-empty :image-size="100" v-if="!loading && items.length === 0" />
  </div>
</template>

<script setup>
import { httpPost } from '@/utils/http'
import { dateFormat } from '@/utils/libs'
import { Search } from '@element-plus/icons-vue'
import Clipboard from 'clipboard'
import { ElMessage } from 'element-plus'
import { computed, onMounted, ref } from 'vue'

const items = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)
const loading = ref(false)
const quotaStatuses = {
  waiting_binding: '待绑定 Token 账号',
  pending: '待发放',
  processing: '发放中',
  succeeded: '已到账',
  failed: '发放失败',
  unknown: '需核对到账结果',
}
const query = ref({
  model: '',
  date: [],
  type: 0,
  quota_status: '',
  page: 1,
  page_size: 10,
})
const totalPower = ref(0)

const totalPages = computed(() => Math.max(1, Math.ceil((total.value || 0) / (pageSize.value || 10))))
const startIndex = computed(() => (total.value ? (page.value - 1) * pageSize.value + 1 : 0))
const endIndex = computed(() => Math.min(page.value * pageSize.value, total.value || 0))

const typePillClass = (type) => {
  const map = {
    1: 'pill-green',
    2: 'pill-blue',
    3: 'pill-red',
    4: 'pill-orange',
    5: 'pill-purple',
    7: 'pill-gray',
  }
  return map[type] || 'pill-gray'
}

const quotaPillClass = (status) => {
  const map = {
    succeeded: 'pill-green',
    failed: 'pill-red',
    pending: 'pill-orange',
    processing: 'pill-blue',
    waiting_binding: 'pill-gray',
    unknown: 'pill-purple',
  }
  return map[status] || 'pill-gray'
}

onMounted(() => {
  fetchData()
  const clipboard = new Clipboard('.copy-order-no')
  clipboard.on('success', () => {
    ElMessage.success('复制成功！')
  })

  clipboard.on('error', () => {
    ElMessage.error('复制失败！')
  })
})

// 搜索
const search = () => {
  page.value = 1
  fetchData()
}

const changePageSize = (size) => {
  page.value = 1
  fetchData()
}

// 获取数据
const fetchData = () => {
  loading.value = true
  query.value.page = page.value
  query.value.page_size = pageSize.value
  httpPost('/api/admin/powerLog/list', query.value)
    .then((res) => {
      const data = res.data.data
      if (data) {
        items.value = data.items
        total.value = data.total
        page.value = data.page
        pageSize.value = data.page_size
      }
      totalPower.value = res.data.stat
      loading.value = false
    })
    .catch((e) => {
      loading.value = false
      ElMessage.error('获取数据失败：' + e.message)
    })
}
</script>

<style lang="scss" scoped>
.power-log {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'PingFang SC', 'Microsoft YaHei', sans-serif;

  .handle-box {
    display: flex;
    align-items: center;
    flex-wrap: wrap;
    gap: 10px;
    margin-bottom: 16px;
  }

  .filter-box {
    display: flex;
    align-items: center;
    flex-wrap: wrap;
    gap: 10px;
  }

  .filter-input {
    width: 140px;
    font-size: 14px;

    :deep(.el-input__wrapper) {
      height: 34px;
      background: #f2f4f7;
      box-shadow: none;
      border: 1px solid transparent;
      border-radius: 6px;
      padding: 0 12px;
    }

    :deep(.el-input__wrapper:hover) {
      border-color: #e2e6ee;
    }

    :deep(.el-input__wrapper.is-focus) {
      background: #fff;
      border-color: #1677ff;
      box-shadow: 0 0 0 2px rgba(22, 119, 255, 0.08);
    }

    :deep(.el-input__inner) {
      font-size: 14px;
      color: #1d2129;
    }

    :deep(.el-input__inner::placeholder) {
      color: #98a2b3;
    }

    &.input-uid {
      width: 120px;
    }
  }

  .filter-select {
    :deep(.el-select__wrapper) {
      min-height: 34px;
      background: #f2f4f7;
      box-shadow: none;
      border: 1px solid transparent;
      border-radius: 6px;
      padding: 0 12px;
      font-size: 14px;
    }

    :deep(.el-select__wrapper:hover) {
      border-color: #e2e6ee;
    }

    :deep(.el-select__wrapper.is-focused) {
      background: #fff;
      border-color: #1677ff;
      box-shadow: 0 0 0 2px rgba(22, 119, 255, 0.08);
    }

    :deep(.el-select__placeholder) {
      color: #98a2b3;
      font-size: 14px;
    }

    &.select-type {
      width: 120px;
    }

    &.select-quota {
      width: 150px;
    }
  }

  .date-range {
    width: 250px;
    height: 34px;
    padding: 0 12px;
    background: #f2f4f7;
    border: 1px solid transparent;
    border-radius: 6px;
    box-shadow: none;

    &:hover {
      border-color: #e2e6ee;
    }

    &.is-active {
      background: #fff;
      border-color: #1677ff;
      box-shadow: 0 0 0 2px rgba(22, 119, 255, 0.08);
    }

    :deep(.el-range-input) {
      font-size: 14px;
      color: #1d2129;
    }

    :deep(.el-range-input::placeholder) {
      color: #98a2b3;
    }

    :deep(.el-range-separator) {
      color: #98a2b3;
    }
  }

  .query-btn {
    display: inline-flex;
    align-items: center;
    gap: 4px;
    height: 34px;
    padding: 0 16px;
    border: 1px solid #e2e6ee;
    border-radius: 6px;
    background: #f2f4f7;
    color: #1d2129;
    font-size: 14px;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s ease;

    &:hover {
      border-color: #1677ff;
      background: #e9ebf0;
      color: #1677ff;
    }
  }

  .toolbar-right {
    display: flex;
    align-items: center;
    margin-left: auto;
  }

  .total-info {
    font-size: 13px;
    color: #667085;
  }

  .log-table {
    width: 100%;

    :deep(.el-table__inner-wrapper::before) {
      display: none;
    }

    :deep(th.el-table__cell) {
      background: #f7f8fa;
      color: #667085;
      font-size: 13px;
      font-weight: 500;
      border-bottom: none;
      padding: 10px 0;
    }

    :deep(td.el-table__cell) {
      border-bottom: 1px solid #f2f3f5;
      font-size: 13px;
      color: #1d2129;
      padding: 14px 0;
    }

    :deep(.el-table__body tr:hover > td.el-table__cell) {
      background: #fafbfc;
    }
  }

  .amount-in {
    color: #22c55e;
    font-weight: 500;
  }

  .amount-out {
    color: #ef4444;
    font-weight: 500;
  }

  .pill {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    height: 24px;
    padding: 0 10px;
    border-radius: 999px;
    font-size: 12px;
    line-height: 1;
    white-space: nowrap;

    &.pill-green {
      background: #e9f9ef;
      color: #22c55e;
    }

    &.pill-red {
      background: #fdecec;
      color: #ef4444;
    }

    &.pill-blue {
      background: #e8f2ff;
      color: #1677ff;
    }

    &.pill-orange {
      background: #fff4e5;
      color: #f59e0b;
    }

    &.pill-purple {
      background: #f3eeff;
      color: #8b5cf6;
    }

    &.pill-gray {
      background: #f2f4f7;
      color: #667085;
    }
  }

  .table-footer {
    display: flex;
    justify-content: space-between;
    align-items: center;
    flex-wrap: wrap;
    gap: 12px;
    margin-top: 16px;
    font-size: 13px;
    color: #667085;
  }

  .footer-nav {
    display: flex;
    align-items: center;
    gap: 14px;
  }

  .footer-text {
    color: #667085;
  }

  :deep(.el-pagination) {
    padding: 0;
    --el-pagination-button-bg-color: transparent;
  }

  :deep(.el-pager li) {
    min-width: 30px;
    height: 30px;
    line-height: 30px;
    border-radius: 6px;
    background: transparent;
    color: #667085;
    font-weight: 500;
  }

  :deep(.el-pager li.is-active) {
    background: #e8f2ff;
    color: #1677ff;
  }

  :deep(.btn-prev),
  :deep(.btn-next) {
    background: transparent;
    color: #667085;
  }

  .page-size-select {
    width: 90px;
  }
}

:global([data-theme='dark'] .power-log) :deep(th.el-table__cell) {
  background: rgba(255, 255, 255, 0.04);
  color: rgba(255, 255, 255, 0.55);
}

:global([data-theme='dark'] .power-log) :deep(td.el-table__cell) {
  border-bottom-color: rgba(255, 255, 255, 0.06);
  color: rgba(255, 255, 255, 0.85);
}

:global([data-theme='dark'] .power-log) .filter-input :deep(.el-input__wrapper),
:global([data-theme='dark'] .power-log) .filter-select :deep(.el-select__wrapper),
:global([data-theme='dark'] .power-log) .date-range {
  background: rgba(255, 255, 255, 0.06);
}
</style>
