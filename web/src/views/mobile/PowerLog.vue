<template>
  <div class="power-log">
    <div class="power-content">
      <!-- 统计概览 -->
      <div class="stats-overview">
        <div class="stats-card">
          <van-row :gutter="12">
            <van-col :span="8">
              <div class="stat-item">
                <div class="stat-value">{{ stats.total }}</div>
                <div class="stat-label">总消费</div>
              </div>
            </van-col>
            <van-col :span="8">
              <div class="stat-item">
                <div class="stat-value">{{ stats.today }}</div>
                <div class="stat-label">今日消费</div>
              </div>
            </van-col>
            <van-col :span="8">
              <div class="stat-item">
                <div class="stat-value">{{ stats.balance }}</div>
                <div class="stat-label">剩余算力</div>
              </div>
            </van-col>
          </van-row>
        </div>
      </div>

      <!-- 筛选栏 -->
      <div class="filter-bar" style="display: none"></div>

      <!-- 日志列表 -->
      <div class="log-list">
        <van-pull-refresh
          :model-value="refreshing"
          @update:model-value="refreshing = $event"
          @refresh="onRefresh"
        >
          <van-list
            :model-value="loading"
            @update:model-value="loading = $event"
            :finished="finished"
            finished-text="没有更多了"
            @load="onLoad"
          >
            <div v-for="item in logList" :key="item.id" class="log-item">
              <div class="log-header">
                <div class="log-icon" :style="{ backgroundColor: getTypeColor(item.type) }">
                  <i class="iconfont" :class="getTypeIcon(item.type)"></i>
                </div>
                <div class="log-info">
                  <div class="log-title">
                    {{ item.model || getTypeTitle(item.type) }}
                    <van-tag type="primary" class="ml-2">{{ item.type_str }}</van-tag>
                  </div>
                  <div class="log-time">{{ formatTime(item.created_at) }}</div>
                </div>
                <div class="log-cost">
                  <span class="cost-value" :class="{ income: item.mark === 1 }">
                    {{ item.mark === 1 ? '+' : '-' }}{{ item.amount }}
                  </span>
                  <span class="cost-unit">算力</span>
                </div>
              </div>
              <div class="log-detail" v-if="item.remark">
                <van-text-ellipsis :content="item.remark" expand-text="展开" collapse-text="收起" />
              </div>
              <div class="log-balance" v-if="item.balance !== undefined">
                <span class="balance-label">余额：</span>
                <span class="balance-value">{{ item.balance }}</span>
              </div>
            </div>

            <!-- 空状态 -->
            <van-empty v-if="!loading && logList.length === 0" description="暂无消费记录" />
          </van-list>
        </van-pull-refresh>
      </div>
    </div>
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { httpPost, httpGet } from '@/utils/http'
import { checkSession } from '@/store/cache'
import { ElMessage } from 'element-plus'

const router = useRouter()
const loading = ref(false)
const finished = ref(false)
const refreshing = ref(false)
const logList = ref([])

// 统计数据
const stats = ref({
  total: 0,
  today: 0,
  balance: 0,
})

// 分页参数
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)

onMounted(() => {
  checkSession()
    .then(() => {
      fetchStats()
      fetchData()
    })
    .catch(() => {})
})

// 获取统计数据
const fetchStats = () => {
  // 调用后端统计API
  httpGet('/api/powerLog/stats')
    .then((res) => {
      if (res.data) {
        stats.value = {
          total: res.data.total || 0,
          today: res.data.today || 0,
          balance: res.data.balance || 0,
        }
      }
    })
    .catch((e) => {
      console.error('获取统计数据失败:', e)
      // 使用默认值
      stats.value = {
        total: 0,
        today: 0,
        balance: 0,
      }
    })
}

// 获取数据
const fetchData = () => {
  loading.value = true
  httpPost('/api/powerLog/list', {
    model: '', // 移除筛选参数
    date: [], // 移除筛选参数
    page: page.value,
    page_size: pageSize.value,
  })
    .then((res) => {
      const items = res.data.items || []
      if (items.length === 0) {
        finished.value = true
        return
      }
      if (page.value === 1) {
        logList.value = items
      } else {
        logList.value.push(...res.data.items)
      }
      total.value = res.data.total
      // 判断是否加载完成
      if (logList.value.length >= total.value) {
        finished.value = true
      }
    })
    .catch((e) => {
      loading.value = false
      refreshing.value = false
      ElMessage.error('获取数据失败：' + e.message)
    })
    .finally(() => {
      loading.value = false
      refreshing.value = false
    })
}

// 加载更多
const onLoad = () => {
  if (finished.value) return
  page.value++
  fetchData()
}

// 下拉刷新
const onRefresh = () => {
  finished.value = false
  page.value = 1
  refreshing.value = true
  fetchData()
}

// 获取类型图标
const getTypeIcon = (type) => {
  const icons = {
    1: 'icon-recharge', // 充值
    2: 'icon-chat', // 消费
    3: 'icon-withdraw-log', // 退款
    4: 'icon-yaoqm', // 邀请
    5: 'icon-redeem', // 兑换
    6: 'icon-present', // 赠送
    7: 'icon-linggan', // 签到
  }
  return icons[type] || 'icon-chat'
}

// 获取类型颜色
const getTypeColor = (type) => {
  const colors = {
    1: '#07c160', // 充值 - 绿色
    2: '#1989fa', // 消费 - 蓝色
    3: '#ff976a', // 退款 - 橙色
    4: '#8B5CF6', // 邀请 - 紫色
    5: '#ee0a24', // 兑换 - 红色
    6: '#07c160', // 赠送 - 绿色
    7: '#1989fa', // 签到 - 蓝色
  }
  return colors[type] || '#1989fa'
}

// 获取类型标题
const getTypeTitle = (type) => {
  const titles = {
    1: '充值',
    2: '消费',
    3: '退款',
    4: '邀请奖励',
    5: '兑换',
    6: '系统赠送',
    7: '每日签到',
  }
  return titles[type] || '其他'
}

// 格式化时间
const formatTime = (timestamp) => {
  const date = new Date(timestamp * 1000)
  const now = new Date()
  const diff = now - date

  if (diff < 60000) {
    // 1分钟内
    return '刚刚'
  } else if (diff < 3600000) {
    // 1小时内
    return `${Math.floor(diff / 60000)}分钟前`
  } else if (diff < 86400000) {
    // 24小时内
    return `${Math.floor(diff / 3600000)}小时前`
  } else {
    return date.toLocaleDateString()
  }
}

// 移除 showFilter, showDatePicker, query.date, onDateButtonClick, onDateConfirm, resetFilter, applyFilter 相关逻辑
</script>

<style lang="scss" scoped>
.power-log {
  min-height: 100vh;
  background: var(--van-background);

  .power-content {
    .stats-overview {
      padding: 16px;
      background: linear-gradient(135deg, var(--van-primary-color), #8b5cf6);

      .stats-card {
        background: rgba(255, 255, 255, 0.1);
        border-radius: 12px;
        padding: 16px;
        backdrop-filter: blur(10px);

        .stat-item {
          text-align: center;

          .stat-value {
            font-size: 20px;
            font-weight: 700;
            color: white;
            margin-bottom: 4px;
          }

          .stat-label {
            font-size: 12px;
            color: rgba(255, 255, 255, 0.8);
          }
        }
      }
    }

    .filter-bar {
      background: var(--van-background);
      border-bottom: 1px solid var(--van-border-color);
      display: flex;
      align-items: center;
      padding: 0 16px;

      .filter-actions {
        margin-left: 12px;
      }

      :deep(.van-tabs__nav) {
        padding: 0;
      }

      :deep(.van-tab) {
        font-size: 14px;
      }
    }

    .log-list {
      padding: 8px 16px 60px;

      .log-item {
        background: var(--van-cell-background);
        border-radius: 12px;
        margin-bottom: 12px;
        padding: 16px;
        box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);

        .log-header {
          display: flex;
          align-items: center;

          .log-icon {
            width: 40px;
            height: 40px;
            border-radius: 10px;
            display: flex;
            align-items: center;
            justify-content: center;
            margin-right: 12px;

            .iconfont {
              font-size: 20px;
              color: white;
            }
          }

          .log-info {
            flex: 1;

            .log-title {
              font-size: 15px;
              font-weight: 500;
              color: var(--van-text-color);
              margin-bottom: 4px;
            }

            .log-time {
              font-size: 12px;
              color: var(--van-gray-6);
            }
          }

          .log-cost {
            text-align: right;

            .cost-value {
              font-size: 16px;
              font-weight: 600;
              color: #ee0a24;

              &.income {
                color: #07c160;
              }
            }

            .cost-unit {
              font-size: 12px;
              color: var(--van-gray-6);
              margin-left: 2px;
            }
          }
        }

        .log-detail {
          margin-top: 12px;
          padding-top: 12px;
          border-top: 1px solid var(--van-border-color);

          :deep(.van-text-ellipsis__text) {
            font-size: 13px;
            color: var(--van-gray-6);
            line-height: 1.4;
          }
        }

        .log-balance {
          margin-top: 8px;
          font-size: 12px;
          color: var(--van-gray-6);

          .balance-label {
            margin-right: 4px;
          }

          .balance-value {
            font-weight: 500;
            color: var(--van-text-color);
          }
        }
      }
    }
  }

  .filter-content {
    padding: 20px;

    .filter-actions {
      display: flex;
      gap: 12px;
      margin-top: 20px;

      .van-button {
        flex: 1;
      }
    }
  }
}

// 深色主题优化
:deep(.van-theme-dark) {
  .power-log {
    .stats-card {
      background: rgba(255, 255, 255, 0.05) !important;
    }

    .log-item {
      box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
    }
  }
}
</style>
