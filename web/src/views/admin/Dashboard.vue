<template>
  <div class="dashboard">
    <!-- 顶部工具条：刷新 -->
    <div class="dash-toolbar">
      <div class="tool-actions">
        <button class="tool-btn" title="刷新" @click="reloadPage">
          <el-icon><Refresh /></el-icon>
        </button>
      </div>
    </div>

    <!-- 指标卡片 -->
    <div class="stat-grid">
      <div class="stat-card" v-for="card in cards" :key="card.title">
        <div class="stat-head">
          <el-icon class="stat-head-icon"><component :is="card.icon" /></el-icon>
          <span>{{ card.title }}</span>
        </div>
        <div class="stat-body">
          <div class="metric" v-for="m in card.metrics" :key="m.label">
            <span class="metric-badge" :style="{ background: m.color }">
              <el-icon><component :is="m.icon" /></el-icon>
            </span>
            <div class="metric-main">
              <div class="metric-label">{{ m.label }}</div>
              <div class="metric-value">{{ m.value }}</div>
              <div class="metric-sub" v-if="m.sub">{{ m.sub }}</div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 趋势 -->
    <div class="trend-card">
      <div class="trend-head">
        <div class="trend-head-title">
          <el-icon><Histogram /></el-icon>
          <h3>趋势</h3>
        </div>
        <div class="trend-tabs">
          <template v-for="(tab, i) in trendTabs" :key="tab.key">
            <span v-if="i > 0" class="tab-divider">/</span>
            <button
              class="trend-tab"
              :class="{ active: activeTab === tab.key }"
              @click="activeTab = tab.key"
            >
              {{ tab.label }}
            </button>
          </template>
        </div>
      </div>
      <div class="trend-body">
        <div class="trend-caption">
          <div class="caption-title">{{ currentTab.label }}</div>
          <div class="caption-total">{{ captionTotal }}</div>
        </div>
        <div ref="trendChartEl" class="trend-chart"></div>
      </div>
    </div>

    <!-- 最近订单 / 最近用户 -->
    <div class="bottom-grid">
      <div class="panel">
        <div class="panel-head">
          <el-icon><Tickets /></el-icon>
          <h3>最近订单</h3>
        </div>
        <div class="panel-list">
          <div class="list-item" v-for="order in recentOrders" :key="order.order_no">
            <div class="item-main">
              <span class="item-title mono">#{{ order.order_no }}</span>
            </div>
            <div class="item-meta">
              <span class="item-amount mono">￥{{ order.amount }}</span>
              <span class="item-time">{{ formatTime(order.created_at) }}</span>
            </div>
          </div>
          <div class="empty" v-if="!recentOrders.length">暂无数据</div>
        </div>
      </div>

      <div class="panel">
        <div class="panel-head">
          <el-icon><User /></el-icon>
          <h3>最近用户</h3>
        </div>
        <div class="panel-list">
          <div class="list-item" v-for="user in recentUsers" :key="user.nickname">
            <div class="item-main">
              <span class="item-avatar">{{ (user.nickname || '?').charAt(0).toUpperCase() }}</span>
              <span class="item-title">{{ user.nickname }}</span>
            </div>
            <div class="item-meta">
              <span class="item-time">{{ formatTime(user.last_active) }}</span>
            </div>
          </div>
          <div class="empty" v-if="!recentUsers.length">暂无数据</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { httpGet } from '@/utils/http'
import {
  ChatDotRound,
  Coin,
  DataAnalysis,
  Headset,
  Histogram,
  Lightning,
  Money,
  Picture,
  Refresh,
  ShoppingCart,
  Tickets,
  User,
  VideoPlay,
  Wallet,
} from '@element-plus/icons-vue'
import * as echarts from 'echarts'
import { ElMessage } from 'element-plus'
import { computed, nextTick, onActivated, onMounted, onUnmounted, ref, watch } from 'vue'

const stats = ref({
  users: 0,
  chats: 0,
  tokens: 0,
  power: 0,
  income: 0,
  orders: 0,
  activeUsers: 0,
  powerConsumption: 0,
  imageJobs: 0,
  videoJobs: 0,
  musicJobs: 0,
  todayUsers: 0,
  todayOrders: 0,
  todayIncome: 0,
  todayTokens: 0,
  todayPower: 0,
  todayImageJobs: 0,
  todayVideoJobs: 0,
  todayMusicJobs: 0,
})

// 数据列表
const recentOrders = ref([])
const recentUsers = ref([])

// 趋势数据源（用于卡片迷你折线）
const chartData = ref({ users: {}, historyMessage: {}, orders: {} })
const trendTabs = [
  { key: 'income', label: '收入趋势', series: 'orders', color: '#22c55e' },
  { key: 'users', label: '用户增长趋势', series: 'users', color: '#0ea5e9' },
]
const activeTab = ref('income')
const trendChartEl = ref(null)
let trendChart = null

// 数字格式化
const fmt = (num) => Number(num || 0).toLocaleString('en-US')
const fmtMoney = (num) =>
  Number(num || 0).toLocaleString('en-US', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
const fmtCompact = (num) => {
  const n = Number(num || 0)
  if (n >= 1e9) return (n / 1e9).toFixed(1) + 'B'
  if (n >= 1e6) return (n / 1e6).toFixed(1) + 'M'
  if (n >= 1e4) return (n / 1e4).toFixed(1) + '万'
  return fmt(n)
}

// 时间格式化工具
const formatTime = (dateStr) => {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  const diff = (new Date() - date) / 1000
  if (diff < 60) return '刚刚'
  if (diff < 60 * 60) return Math.floor(diff / 60) + '分钟前'
  if (diff < 60 * 60 * 24) return Math.floor(diff / 3600) + '小时前'
  return Math.floor(diff / 86400) + '天前'
}

// 折线 X 轴日期（升序）
const chartLabels = computed(() => {
  const keys = new Set()
  Object.values(chartData.value).forEach((series) => {
    if (series) Object.keys(series).forEach((k) => keys.add(k))
  })
  return Array.from(keys).sort()
})

// 指标卡片配置
const cards = computed(() => {
  const s = stats.value
  return [
    {
      title: '账户数据',
      icon: Wallet,
      metrics: [
        {
          label: '总收入',
          value: '￥' + fmtMoney(s.income),
          sub: '今日 ￥' + fmtMoney(s.todayIncome),
          color: '#3b82f6',
          icon: Money,
        },
        {
          label: '积分消耗',
          value: fmt(s.power || s.tokens),
          sub: '今日 ' + fmt(s.todayPower || s.todayTokens),
          color: '#9ca3af',
          icon: Coin,
        },
      ],
    },
    {
      title: '使用统计',
      icon: DataAnalysis,
      metrics: [
        {
          label: '用户总数',
          value: fmt(s.users),
          sub: '今日新增 ' + fmt(s.todayUsers),
          color: '#22c55e',
          icon: User,
        },
        {
          label: '对话总数',
          value: fmt(s.chats),
          sub: '今日 ' + fmt(s.todayChats),
          color: '#14b8a6',
          icon: ChatDotRound,
        },
      ],
    },
    {
      title: '生成任务',
      icon: Picture,
      metrics: [
        {
          label: '图片生成',
          value: fmt(s.imageJobs),
          sub: '今日 ' + fmt(s.todayImageJobs),
          color: '#f59e0b',
          icon: Picture,
        },
        {
          label: '视频生成',
          value: fmt(s.videoJobs),
          sub: '今日 ' + fmt(s.todayVideoJobs),
          color: '#eab308',
          icon: VideoPlay,
        },
      ],
    },
    {
      title: '业务指标',
      icon: Lightning,
      metrics: [
        {
          label: '音乐生成',
          value: fmt(s.musicJobs),
          sub: '今日 ' + fmt(s.todayMusicJobs),
          color: '#06b6d4',
          icon: Headset,
        },
        {
          label: '订单总数',
          value: fmt(s.orders),
          sub: '今日 ' + fmt(s.todayOrders),
          color: '#64748b',
          icon: ShoppingCart,
        },
      ],
    },
  ]
})

const currentTab = computed(() => trendTabs.find((t) => t.key === activeTab.value) || trendTabs[0])

const trendValues = computed(() => {
  const series = chartData.value[currentTab.value.series] || {}
  return chartLabels.value.map((label) => Number(series[label] || 0))
})

const trendTotal = computed(() => trendValues.value.reduce((a, b) => a + b, 0))

const captionTotal = computed(() => {
  const prefix = currentTab.value.key === 'income' ? '￥' : ''
  return '总计：' + prefix + fmtCompact(trendTotal.value)
})

// 渲染趋势图
const renderTrend = () => {
  if (!trendChartEl.value) return
  if (!trendChart) {
    trendChart = echarts.init(trendChartEl.value)
  }
  const tab = currentTab.value
  trendChart.setOption(
    {
      tooltip: {
        trigger: 'axis',
        backgroundColor: 'rgba(17, 20, 45, 0.92)',
        borderWidth: 0,
        padding: [8, 12],
        textStyle: { color: '#fff', fontSize: 12 },
        axisPointer: { type: 'shadow', shadowStyle: { color: 'rgba(99, 102, 241, 0.06)' } },
        valueFormatter: (v) => (tab.key === 'income' ? '￥' + fmtMoney(v) : fmt(v)),
      },
      grid: { left: 8, right: 16, top: 24, bottom: 6, containLabel: true },
      xAxis: {
        type: 'category',
        data: chartLabels.value,
        axisLine: { lineStyle: { color: 'rgba(128, 128, 128, 0.25)' } },
        axisTick: { show: false },
        axisLabel: { color: '#8c8c8c', fontSize: 11 },
      },
      yAxis: {
        type: 'value',
        axisLine: { show: false },
        axisTick: { show: false },
        axisLabel: {
          color: '#8c8c8c',
          fontSize: 11,
          formatter: (v) => fmtCompact(v),
        },
        splitLine: { lineStyle: { color: 'rgba(128, 128, 128, 0.14)' } },
      },
      series: [
        {
          type: 'bar',
          data: trendValues.value,
          barMaxWidth: 10,
          itemStyle: { color: tab.color, borderRadius: [3, 3, 0, 0] },
        },
      ],
    },
    true
  )
}

const fetchStats = () => {
  httpGet('/api/admin/dashboard/stats', { days: 7 })
    .then((res) => {
      Object.assign(stats.value, res.data)
      recentOrders.value = res.data.recentOrders || []
      recentUsers.value = res.data.recentUsers || []
      chartData.value = res.data.chart || { users: {}, historyMessage: {}, orders: {} }
      nextTick(renderTrend)
    })
    .catch((e) => {
      ElMessage.error('获取统计数据失败：' + e.message)
    })
}

// 刷新相当于浏览器刷新：整页重新加载
const reloadPage = () => window.location.reload()

const handleResize = () => {
  if (trendChart) trendChart.resize()
}

watch(activeTab, () => nextTick(renderTrend))

onMounted(() => {
  fetchStats()
  window.addEventListener('resize', handleResize)
})

onActivated(() => {
  nextTick(() => trendChart && trendChart.resize())
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  if (trendChart) {
    trendChart.dispose()
    trendChart = null
  }
})
</script>

<style scoped lang="scss">
.dashboard {
  --dash-surface: var(--admin-container-bg, #ffffff);
  --dash-border: rgba(0, 0, 0, 0.08);
  --dash-shadow: none;
  --dash-tool-bg: #ffffff;
  --dash-tool-border: rgba(0, 0, 0, 0.12);
  --dash-tool-color: #4b5563;
  --dash-hover: rgba(0, 0, 0, 0.02);
  --dash-divider: rgba(0, 0, 0, 0.06);

  padding: 20px 24px 28px;
  background: var(--admin-page-bg, #ffffff);
  min-height: 100%;
  box-sizing: border-box;

  .dash-toolbar {
    display: flex;
    justify-content: flex-end;
    align-items: center;
    gap: 12px;
    margin-bottom: 16px;
  }

  .tool-actions {
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .tool-btn {
    width: 34px;
    height: 34px;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    border-radius: 8px;
    border: 1px solid var(--dash-tool-border);
    background: var(--dash-tool-bg);
    color: var(--dash-tool-color);
    cursor: pointer;
    transition: all 0.2s ease;

    &:hover {
      color: var(--theme-text-color-primary);
      background: var(--dash-hover);
    }

    .el-icon {
      font-size: 16px;
    }
  }

  // 指标卡片
  .stat-grid {
    display: grid;
    grid-template-columns: repeat(4, minmax(0, 1fr));
    gap: 18px;
  }

  .stat-card {
    position: relative;
    background: var(--dash-surface);
    border: 1px solid var(--dash-border);
    border-radius: 4px;
    padding: 18px 22px 20px;
    overflow: hidden;
    min-height: 176px;
    box-sizing: border-box;
  }

  .stat-head {
    display: flex;
    align-items: center;
    gap: 8px;
    padding-bottom: 12px;
    margin-bottom: 16px;
    border-bottom: 1px solid var(--dash-divider);
    font-size: 15px;
    font-weight: 600;
    color: var(--theme-text-color-primary);

    .stat-head-icon {
      font-size: 17px;
      color: var(--theme-text-color-secondary);
    }
  }

  .stat-body {
    position: relative;
    display: flex;
    flex-direction: column;
    gap: 16px;
  }

  .metric {
    position: relative;
    z-index: 1;
    display: flex;
    align-items: flex-start;
    gap: 12px;
    min-height: 46px;
    box-sizing: border-box;
  }

  .metric-badge {
    width: 32px;
    height: 32px;
    flex-shrink: 0;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    border-radius: 4px;
    color: #fff;

    .el-icon {
      font-size: 16px;
    }
  }

  .metric-main {
    min-width: 0;
  }

  .metric-label {
    font-size: 13px;
    line-height: 1.2;
    color: var(--theme-text-color-secondary);
    margin-bottom: 4px;
  }

  .metric-value {
    font-size: 20px;
    font-weight: 600;
    line-height: 1.15;
    color: var(--theme-text-color-primary);
    white-space: nowrap;
  }

  .metric-sub {
    margin-top: 4px;
    font-size: 11px;
    line-height: 1.2;
    color: var(--theme-text-color-secondary);
    opacity: 0.85;
    white-space: nowrap;
  }

  // 趋势
  .trend-card {
    margin-top: 18px;
    background: var(--dash-surface);
    border: 1px solid var(--dash-border);
    border-radius: 4px;
    padding: 20px 24px 10px;
  }

  .trend-head {
    display: flex;
    justify-content: space-between;
    align-items: center;
    gap: 16px;
  }

  .trend-head-title {
    display: flex;
    align-items: center;
    gap: 10px;

    .el-icon {
      font-size: 20px;
      color: var(--theme-text-color-primary);
    }

    h3 {
      margin: 0;
      font-size: 18px;
      font-weight: 600;
      color: var(--theme-text-color-primary);
    }
  }

  .trend-tabs {
    display: flex;
    align-items: center;
  }

  .tab-divider {
    margin: 0 10px;
    color: var(--dash-divider);
    font-size: 13px;
  }

  .trend-tab {
    border: none;
    background: transparent;
    padding: 0;
    cursor: pointer;
    font-size: 14px;
    color: var(--theme-text-color-secondary);
    transition: color 0.2s ease;

    &:hover {
      color: var(--theme-text-color-primary);
    }

    &.active {
      color: var(--theme-text-color-primary);
      font-weight: 600;
    }
  }

  .trend-body {
    margin-top: 18px;
    padding-top: 16px;
    border-top: 1px solid var(--dash-divider);
  }

  .trend-caption {
    padding-left: 2px;

    .caption-title {
      font-size: 15px;
      font-weight: 600;
      color: var(--theme-text-color-primary);
    }

    .caption-total {
      margin-top: 5px;
      font-size: 12px;
      color: var(--theme-text-color-secondary);
    }
  }

  .trend-chart {
    width: 100%;
    height: 320px;
  }

  // 底部列表
  .bottom-grid {
    display: grid;
    grid-template-columns: repeat(2, minmax(0, 1fr));
    gap: 18px;
    margin-top: 18px;
  }

  .panel {
    background: var(--dash-surface);
    border: 1px solid var(--dash-border);
    border-radius: 4px;
    padding: 18px 24px 8px;
  }

  .panel-head {
    display: flex;
    align-items: center;
    gap: 8px;
    padding-bottom: 12px;
    border-bottom: 1px solid var(--dash-divider);

    .el-icon {
      font-size: 17px;
      color: var(--theme-text-color-secondary);
    }

    h3 {
      margin: 0;
      font-size: 16px;
      font-weight: 600;
      color: var(--theme-text-color-primary);
    }
  }

  .panel-list {
    max-height: 300px;
    overflow-y: auto;
  }

  .list-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
    gap: 12px;
    padding: 11px 0;
    border-bottom: 1px solid var(--dash-divider);

    &:last-child {
      border-bottom: none;
    }
  }

  .item-main {
    display: flex;
    align-items: center;
    gap: 10px;
    min-width: 0;
  }

  .item-avatar {
    width: 28px;
    height: 28px;
    flex-shrink: 0;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    border-radius: 4px;
    background: var(--el-color-primary, #010120);
    color: #fff;
    font-size: 13px;
    font-weight: 600;
    line-height: 1;
  }

  .item-title {
    font-size: 14px;
    color: var(--theme-text-color-primary);
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .item-meta {
    display: flex;
    align-items: center;
    gap: 14px;
    flex-shrink: 0;
  }

  .item-amount {
    font-size: 14px;
    font-weight: 600;
    color: var(--theme-text-color-primary);
  }

  .item-time {
    font-size: 12px;
    color: var(--theme-text-color-secondary);
  }

  .mono {
    font-family: 'SFMono-Regular', Consolas, 'Liberation Mono', monospace;
  }

  .empty {
    padding: 32px 0;
    text-align: center;
    font-size: 13px;
    color: var(--theme-text-color-secondary);
  }
}

:global([data-theme='dark'] .dashboard) {
  --dash-border: rgba(255, 255, 255, 0.08);
  --dash-tool-bg: transparent;
  --dash-tool-border: rgba(255, 255, 255, 0.14);
  --dash-tool-color: rgba(255, 255, 255, 0.7);
  --dash-hover: rgba(255, 255, 255, 0.06);
  --dash-divider: rgba(255, 255, 255, 0.08);
}

@media (max-width: 900px) {
  .dashboard {
    .stat-grid,
    .bottom-grid {
      grid-template-columns: minmax(0, 1fr);
    }
  }
}

</style>
