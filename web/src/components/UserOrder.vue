<template>
  <div class="user-bill" v-loading="loading" element-loading-background="rgba(255,255,255,.3)">
    <!-- PC端表格 -->
    <div class="desktop-table" v-if="items.length > 0">
      <el-table :data="items" :row-key="(row) => row.id" table-layout="auto" border>
        <el-table-column prop="order_no" label="订单号">
          <template #default="scope">
            <span>{{ scope.row.order_no }}</span>
            <el-icon class="copy-order-no" :data-clipboard-text="scope.row.order_no">
              <DocumentCopy />
            </el-icon>
          </template>
        </el-table-column>
        <el-table-column prop="subject" label="产品名称" />
        <el-table-column prop="amount" label="订单金额" />
        <el-table-column label="订单算力">
          <template #default="scope">
            <span>{{ scope.row.remark && scope.row.remark.power }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="channel_name" label="支付渠道" />
        <el-table-column prop="pay_name" label="支付名称" />
        <el-table-column label="支付时间">
          <template #default="scope">
            <span v-if="scope.row['pay_time']">{{ dateFormat(scope.row['pay_time']) }}</span>
            <el-tag v-else>未支付</el-tag>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <!-- 移动端卡片列表 -->
    <div class="mobile-cards" v-if="items.length > 0">
      <div v-for="item in items" :key="item.id" class="order-card">
        <div class="card-header">
          <div class="order-no">
            <span class="label">订单号：</span>
            <span class="value">{{ item.order_no }}</span>
            <el-icon class="copy-icon" :data-clipboard-text="item.order_no">
              <DocumentCopy />
            </el-icon>
          </div>
        </div>

        <div class="card-content">
          <div class="info-row">
            <span class="label">产品名称：</span>
            <span class="value">{{ item.subject }}</span>
          </div>
          <div class="info-row">
            <span class="label">订单金额：</span>
            <span class="value amount">￥{{ item.amount }}</span>
          </div>
          <div class="info-row">
            <span class="label">订单算力：</span>
            <span class="value">{{ (item.remark && item.remark.power) || '-' }}</span>
          </div>
          <div class="info-row">
            <span class="label">支付渠道：</span>
            <span class="value">{{ item.pay_method }}</span>
          </div>
          <div class="info-row">
            <span class="label">支付名称：</span>
            <span class="value">{{ item.pay_name }}</span>
          </div>
          <div class="info-row">
            <span class="label">支付时间：</span>
            <span class="value">{{ item['pay_time'] ? dateFormat(item['pay_time']) : '-' }}</span>
          </div>
        </div>
      </div>
    </div>

    <el-empty :image-size="100" v-else :image="nodata" description="暂无数据" />

    <div class="pagination pb-5">
      <el-pagination
        v-if="total > 0"
        background
        layout="total,prev, pager, next"
        :hide-on-single-page="true"
        :current-page="page"
        :page-size="pageSize"
        @current-change="handlePageChange"
        style="--el-pagination-button-bg-color: rgba(86, 86, 95, 0.2)"
        :total="total"
      />
    </div>
  </div>
</template>

<script setup>
import nodata from '@/assets/img/no-data.png'
import { httpGet } from '@/utils/http'
import { dateFormat } from '@/utils/libs'
import { DocumentCopy } from '@element-plus/icons-vue'
import Clipboard from 'clipboard'
import { ElMessage } from 'element-plus'
import { onMounted, ref } from 'vue'

const items = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(12)
const loading = ref(true)

onMounted(() => {
  fetchData()
  const clipboard = new Clipboard('.copy-order-no, .copy-icon')
  clipboard.on('success', () => {
    ElMessage.success('复制成功！')
  })

  clipboard.on('error', () => {
    ElMessage.error('复制失败！')
  })
})

// 获取数据
const fetchData = () => {
  httpGet('/api/order/list', { page: page.value, page_size: pageSize.value })
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
      ElMessage.error('获取数据失败：' + e.message)
    })
}

// 处理分页变化
const handlePageChange = (newPage) => {
  page.value = newPage
  fetchData()
}
</script>

<style scoped lang="scss">
.user-bill {
  // background-color: var(--el-bg-color);

  .pagination {
    margin: 20px 0 0 0;
    display: flex;
    justify-content: center;
    width: 100%;
  }

  .copy-order-no {
    cursor: pointer;
    position: relative;
    left: 6px;
    top: 2px;
    color: #20a0ff;
  }

  // 移动端卡片样式
  .mobile-cards {
    display: none;

    .order-card {
      background: transparent;
      border-radius: 12px;
      padding: 16px;
      margin-bottom: 16px;
      box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
      border: 1px solid var(--el-border-color-light);
      background: var(--van-cell-background);

      .card-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 12px;
        padding-bottom: 12px;
        border-bottom: 1px solid var(--el-border-color-lighter);

        .order-no {
          display: flex;
          align-items: center;
          flex: 1;
          margin-right: 12px;

          .label {
            font-size: 14px;
            color: var(--el-text-color-regular);
            margin-right: 4px;
          }

          .value {
            font-size: 14px;
            color: var(--el-text-color-primary);
            font-weight: 500;
            margin-right: 8px;
          }

          .copy-icon {
            cursor: pointer;
            color: #20a0ff;
            font-size: 16px;
          }
        }

        .order-status {
          flex-shrink: 0;
        }
      }

      .card-content {
        .info-row {
          display: flex;
          justify-content: space-between;
          align-items: center;
          margin-bottom: 8px;

          &:last-child {
            margin-bottom: 0;
          }

          .label {
            font-size: 14px;
            color: var(--el-text-color-regular);
            flex-shrink: 0;
            margin-right: 8px;
          }

          .value {
            font-size: 14px;
            color: var(--el-text-color-primary);
            text-align: right;
            word-break: break-all;

            &.amount {
              font-weight: 600;
              color: #ff6b35;
            }
          }
        }
      }
    }
  }
}

// 响应式设计
@media (max-width: 768px) {
  .user-bill {
    .desktop-table {
      display: none;
    }

    .mobile-cards {
      display: block;
    }

    .pagination {
      :deep(.el-pagination) {
        .el-pagination__total {
          display: none;
        }

        .el-pager {
          li {
            min-width: 32px;
            height: 32px;
            line-height: 32px;
            font-size: 12px;
          }
        }
      }
    }
  }
}

// 深色主题适配
:deep(.van-theme-dark) {
  .user-bill {
    .mobile-cards .order-card {
      background: transparent;
      border-color: var(--el-border-color);
      box-shadow: none;
    }
  }
}
</style>
