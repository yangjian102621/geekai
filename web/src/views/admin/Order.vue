<template>
  <div class="container order">
    <div class="filter-panel">
      <div class="filter-row">
        <el-date-picker
          v-model="query.pay_time"
          type="daterange"
          range-separator="~"
          start-placeholder="开始日期"
          end-placeholder="结束日期"
          format="YYYY-MM-DD"
          value-format="YYYY-MM-DD"
          class="date-range"
        />
        <el-select v-model="query.status" class="filter-select status-select">
          <el-option
            v-for="item in orderStatus"
            :key="item.value"
            :label="item.label"
            :value="item.value"
          />
        </el-select>
      </div>

      <div class="filter-row">
        <el-input
          v-model="query.order_no"
          placeholder="订单号"
          class="order-search"
          :prefix-icon="Search"
          clearable
          @keyup.enter="search"
        />
        <button class="primary-btn" @click="search">查询</button>
        <button class="primary-btn" @click="exportOrders">
          <el-icon><Download /></el-icon>
          导出
        </button>
        <button class="clear-btn" @click="clearOrders">清空未支付订单</button>
      </div>
    </div>

    <el-table
      v-if="items.length > 0"
      :data="items"
      class="order-table"
      :row-key="(row) => row.id"
    >
      <el-table-column prop="order_no" label="订单号" min-width="190" />
      <el-table-column prop="trade_no" label="交易号" min-width="190" />
      <el-table-column prop="username" label="下单用户" min-width="120" />
      <el-table-column prop="subject" label="产品名称" min-width="150" />
      <el-table-column label="订单金额" width="100">
        <template #default="scope">
          <span class="amount">{{ scope.row.amount }}</span>
        </template>
      </el-table-column>
      <el-table-column label="充值积分" width="100">
        <template #default="scope">
          <span>{{ scope.row.remark?.power }}</span>
        </template>
      </el-table-column>
      <el-table-column label="下单时间" width="160">
        <template #default="scope">
          <span>{{ dateFormat(scope.row['created_at']) }}</span>
        </template>
      </el-table-column>
      <el-table-column label="订单状态" width="100">
        <template #default="scope">
          <span class="pill" :class="statusPillClass(scope.row.status)">
            {{ statusText(scope.row.status) }}
          </span>
        </template>
      </el-table-column>
      <el-table-column label="支付时间" width="160">
        <template #default="scope">
          <span>{{ dateFormat(scope.row.pay_time) || '--' }}</span>
        </template>
      </el-table-column>
      <el-table-column prop="pay_name" label="支付方式" width="110" />
      <el-table-column label="支付渠道" min-width="110">
        <template #default="scope">{{ scope.row.channel_name || '-' }}</template>
      </el-table-column>
      <el-table-column label="操作" width="90" fixed="right">
        <template #default="scope">
          <el-popconfirm title="确定要删除当前记录吗?" @confirm="remove(scope.row)">
            <template #reference>
              <button class="link-btn danger">删除</button>
            </template>
          </el-popconfirm>
        </template>
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
import { httpGet, httpPost } from '@/utils/http'
import { dateFormat } from '@/utils/libs'
import { Download, Search } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { computed, onMounted, ref } from 'vue'

// 变量定义
const items = ref([])
const query = ref({ order_no: '', pay_time: [], status: -1 })
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)
const loading = ref(true)
const orderStatus = ref([
  { value: -1, label: '全部' },
  { value: 0, label: '未支付' },
  { value: 2, label: '已支付' },
])

const totalPages = computed(() => Math.max(1, Math.ceil((total.value || 0) / (pageSize.value || 10))))
const startIndex = computed(() => (total.value ? (page.value - 1) * pageSize.value + 1 : 0))
const endIndex = computed(() => Math.min(page.value * pageSize.value, total.value || 0))

const statusText = (status) => {
  if (status === 2) return '已支付'
  if (status === 3) return '已关闭'
  return '未支付'
}

const statusPillClass = (status) => {
  if (status === 2) return 'pill-green'
  if (status === 3) return 'pill-red'
  return 'pill-orange'
}

onMounted(() => {
  fetchData()
})

const search = () => {
  page.value = 1
  fetchData()
}

const changePageSize = () => {
  page.value = 1
  fetchData()
}

// 获取数据
const fetchData = () => {
  loading.value = true
  query.value.page = page.value
  query.value.page_size = pageSize.value
  httpPost('/api/admin/order/list', query.value)
    .then((res) => {
      if (res.data) {
        items.value = res.data.items
        total.value = res.data.total
        page.value = res.data.page
        pageSize.value = res.data.page_size
      }
      loading.value = false
    })
    .catch((e) => {
      loading.value = false
      ElMessage.error('获取数据失败：' + e.message)
    })
}

const remove = function (row) {
  httpGet('/api/admin/order/remove?id=' + row.id)
    .then(() => {
      ElMessage.success('删除成功！')
      fetchData()
    })
    .catch((e) => {
      ElMessage.error('删除失败：' + e.message)
    })
}

const clearOrders = () => {
  ElMessageBox.confirm('此操作将会删除所有未支付订单，继续操作吗?', '删除提示', {
    confirmButtonText: '确认',
    cancelButtonText: '取消',
    type: 'warning',
  }).then(() => {
    httpGet('/api/admin/order/clear').then(() => {
      ElMessage.success('订单删除成功')
      page.value = 1
      fetchData()
    })
  })
}

// 导出当前页账单
const exportOrders = () => {
  if (!items.value.length) {
    ElMessage.warning('暂无可导出的数据')
    return
  }
  const headers = [
    '订单号',
    '交易号',
    '下单用户',
    '产品名称',
    '订单金额',
    '充值积分',
    '下单时间',
    '订单状态',
    '支付时间',
    '支付方式',
    '支付渠道',
  ]
  const rows = items.value.map((row) => [
    row.order_no,
    row.trade_no,
    row.username,
    row.subject,
    row.amount,
    row.remark?.power,
    dateFormat(row['created_at']),
    statusText(row.status),
    dateFormat(row.pay_time) || '--',
    row.pay_name,
    row.channel_name,
  ])
  const escape = (val) => `"${String(val ?? '').replace(/"/g, '""')}"`
  const csv = '\ufeff' + [headers, ...rows].map((row) => row.map(escape).join(',')).join('\n')
  const blob = new Blob([csv], { type: 'text/csv;charset=utf-8;' })
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  const now = new Date()
  const stamp = `${now.getFullYear()}${String(now.getMonth() + 1).padStart(2, '0')}${String(
    now.getDate()
  ).padStart(2, '0')}`
  link.href = url
  link.download = `充值账单_${stamp}.csv`
  link.click()
  URL.revokeObjectURL(url)
}
</script>

<style lang="scss" scoped>
.order {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'PingFang SC', 'Microsoft YaHei', sans-serif;

  .filter-panel {
    margin-bottom: 16px;
  }

  .filter-row {
    display: flex;
    align-items: center;
    flex-wrap: wrap;
    gap: 10px;

    & + .filter-row {
      margin-top: 12px;
    }
  }

  .date-range {
    width: 320px;
    height: 40px;
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

  .filter-select {
    :deep(.el-select__wrapper) {
      min-height: 40px;
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

    :deep(.el-select__selected-item) {
      font-size: 14px;
      color: #1d2129;
    }

    &.status-select {
      width: 160px;
    }
  }

  .order-search {
    flex: 1;
    min-width: 240px;
    font-size: 14px;

    :deep(.el-input__wrapper) {
      height: 40px;
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

    :deep(.el-input__prefix) {
      color: #98a2b3;
      font-size: 14px;
    }
  }

  .primary-btn {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    gap: 5px;
    height: 40px;
    padding: 0 18px;
    border: none;
    border-radius: 6px;
    background: #1677ff;
    color: #fff;
    font-size: 14px;
    font-weight: 500;
    cursor: pointer;
    white-space: nowrap;
    transition: background 0.2s ease;

    &:hover {
      background: #4096ff;
    }
  }

  .clear-btn {
    margin-left: auto;
    border: none;
    background: transparent;
    padding: 4px 0;
    color: #ef4444;
    font-size: 13px;
    cursor: pointer;

    &:hover {
      color: #dc2626;
    }
  }

  .order-table {
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

  .amount {
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

    &.pill-orange {
      background: #fff4e5;
      color: #f59e0b;
    }
  }

  .link-btn {
    border: none;
    background: transparent;
    padding: 0;
    color: #667085;
    font-size: 13px;
    cursor: pointer;

    &:hover {
      color: #1d2129;
    }

    &.danger {
      color: #ef4444;

      &:hover {
        color: #dc2626;
      }
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

:global([data-theme='dark'] .order) :deep(th.el-table__cell) {
  background: rgba(255, 255, 255, 0.04);
  color: rgba(255, 255, 255, 0.55);
}

:global([data-theme='dark'] .order) :deep(td.el-table__cell) {
  border-bottom-color: rgba(255, 255, 255, 0.06);
  color: rgba(255, 255, 255, 0.85);
}

:global([data-theme='dark'] .order) .order-search :deep(.el-input__wrapper),
:global([data-theme='dark'] .order) .filter-select :deep(.el-select__wrapper),
:global([data-theme='dark'] .order) .date-range {
  background: rgba(255, 255, 255, 0.06);
}
</style>
