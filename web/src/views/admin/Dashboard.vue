<template>
  <div class="dashboard" v-loading="loading">
    <!-- 统计卡片区域：Together AI Stats Bar —— 纯排版大数字，无 icon 盒 -->
    <el-row class="stats-row" :gutter="24">
      <el-col :span="6">
        <el-card class="stats-card" shadow="hover">
          <div class="card-label">用户总数</div>
          <div class="card-number">{{ stats.users }}</div>
          <div class="card-desc">今日新增 {{ stats.todayUsers || 0 }}</div>
        </el-card>
      </el-col>

      <el-col :span="6">
        <el-card class="stats-card" shadow="hover">
          <div class="card-label">对话总数</div>
          <div class="card-number">{{ stats.chats }}</div>
          <div class="card-desc">今日 {{ stats.todayChats }}</div>
        </el-card>
      </el-col>

      <el-col :span="6">
        <el-card class="stats-card" shadow="hover">
          <div class="card-label">积分消耗</div>
          <div class="card-number">{{ formatNumber(stats.power || stats.tokens) }}</div>
          <div class="card-desc">今日 {{ formatNumber(stats.todayPower || stats.todayTokens) }}</div>
        </el-card>
      </el-col>

      <el-col :span="6">
        <el-card class="stats-card" shadow="hover">
          <div class="card-label">总收入</div>
          <div class="card-number">￥{{ stats.income.toFixed(2) }}</div>
          <div class="card-desc">今日 ￥{{ stats.todayIncome?.toFixed(2) || '0.00' }}</div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 第二行统计卡片 -->
    <el-row class="stats-row" :gutter="24">
      <el-col :span="6">
        <el-card class="stats-card" shadow="hover">
          <div class="card-label">图片生成</div>
          <div class="card-number">{{ stats.imageJobs }}</div>
          <div class="card-desc">今日 {{ stats.todayImageJobs }}</div>
        </el-card>
      </el-col>

      <el-col :span="6">
        <el-card class="stats-card" shadow="hover">
          <div class="card-label">视频生成</div>
          <div class="card-number">{{ stats.videoJobs }}</div>
          <div class="card-desc">今日 {{ stats.todayVideoJobs }}</div>
        </el-card>
      </el-col>

      <el-col :span="6">
        <el-card class="stats-card" shadow="hover">
          <div class="card-label">音乐生成</div>
          <div class="card-number">{{ stats.musicJobs }}</div>
          <div class="card-desc">今日 {{ stats.todayMusicJobs }}</div>
        </el-card>
      </el-col>

      <el-col :span="6">
        <el-card class="stats-card" shadow="hover">
          <div class="card-label">订单总数</div>
          <div class="card-number">{{ stats.orders }}</div>
          <div class="card-desc">今日 {{ stats.todayOrders }}</div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 图表和列表区域 -->
    <el-row :gutter="24" class="content-row">
      <!-- 图表区域 -->
      <el-col :span="12">
        <el-card class="chart-card" shadow="hover">
          <div class="card-header">
            <h3>用户增长趋势</h3>
            <el-button type="text" size="small">30天</el-button>
          </div>
          <div id="chart-users" style="height: 280px"></div>
        </el-card>
      </el-col>

      <el-col :span="12">
        <el-card class="chart-card" shadow="hover">
          <div class="card-header">
            <h3>收入趋势</h3>
            <el-button type="text" size="small">7天</el-button>
          </div>
          <div id="chart-income" style="height: 280px"></div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 底部列表区域 -->
    <el-row :gutter="24" class="content-row">
      <!-- 最近订单 -->
      <el-col :span="12">
        <el-card class="list-card" shadow="hover">
          <div class="card-header">
            <h3>最近订单</h3>
          </div>
          <div class="order-list">
            <div class="order-item" v-for="order in recentOrders" :key="order.order_no">
              <div class="order-info">
                <div class="order-id">#{{ order.order_no }}</div>
                <div class="order-amount">￥{{ order.amount }}</div>
              </div>
              <div class="order-meta">
                <div class="order-date">{{ formatTime(order.created_at) }}</div>
              </div>
            </div>
          </div>
        </el-card>
      </el-col>

      <!-- 最近用户 -->
      <el-col :span="12">
        <el-card class="list-card" shadow="hover">
          <div class="card-header">
            <h3>最近用户</h3>
          </div>
          <div class="user-list">
            <div class="user-item" v-for="user in recentUsers" :key="user.nickname">
              <div class="user-avatar">
                <el-avatar :size="40" :src="user.avatar">{{ user.nickname.charAt(0) }}</el-avatar>
              </div>
              <div class="user-info">
                <div class="user-name">{{ user.nickname }}</div>
              </div>
              <div class="user-meta">
                <div class="user-time">{{ formatTime(user.last_active) }}</div>
              </div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { httpGet } from '@/utils/http'
import {
  ChatDotRound,
  Headset,
  Money,
  Picture,
  ShoppingCart,
  TrendCharts,
  User,
  VideoPlay,
} from '@element-plus/icons-vue'
import * as echarts from 'echarts'
import { ElMessage } from 'element-plus'
import { onMounted, ref } from 'vue'

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
const loading = ref(true)

// 数据列表
const recentOrders = ref([])
const recentUsers = ref([])
const popularJobs = ref([]) // 改为热门任务

const formatNumber = (num) => {
  if (num >= 10000) {
    return (num / 10000).toFixed(1) + 'w'
  }
  return num
}

// 时间格式化工具
const formatTime = (dateStr) => {
  const date = new Date(dateStr)
  const now = new Date()
  const diff = (now - date) / 1000
  if (diff < 60 * 60) {
    return Math.floor(diff / 60) + '分钟前'
  } else if (diff < 60 * 60 * 24) {
    return Math.floor(diff / 3600) + '小时前'
  } else {
    return Math.floor(diff / 86400) + '天前'
  }
}

onMounted(() => {
  const chartUsersEl = document.getElementById('chart-users')
  const chartIncomeEl = document.getElementById('chart-income')

  if (!chartUsersEl || !chartIncomeEl) {
    ElMessage.error('图表容器未找到')
    return
  }

  const chartUsers = echarts.init(chartUsersEl)
  const chartIncome = echarts.init(chartIncomeEl)
  httpGet('/api/admin/dashboard/stats')
    .then((res) => {
      // 更新统计数据
      Object.assign(stats.value, res.data)
      recentOrders.value = res.data.recentOrders || []
      recentUsers.value = res.data.recentUsers || []
      const chartData = res.data.chart || {}
      loading.value = false

      // 检查图表数据是否存在
      if (!chartData.users || !chartData.orders) {
        ElMessage.warning('图表数据不完整')
        return
      }

      const x = []
      const dataUsers = []
      for (let k in chartData.users) {
        x.push(k)
        dataUsers.push(chartData.users[k])
      }
      chartUsers.setOption({
        tooltip: {
          trigger: 'axis',
          backgroundColor: '#ffffff',
          borderColor: 'rgba(0, 0, 0, 0.08)',
          textStyle: { color: '#010120' },
          extraCssText: 'box-shadow: rgba(1, 1, 32, 0.1) 0px 4px 10px; border-radius: 4px;',
        },
        grid: {
          left: '3%',
          right: '4%',
          bottom: '3%',
          containLabel: true,
        },
        xAxis: {
          type: 'category',
          data: x,
          axisLine: { lineStyle: { color: 'rgba(0, 0, 0, 0.08)' } },
          axisTick: { show: false },
          axisLabel: { color: 'rgba(0, 0, 0, 0.6)' },
        },
        yAxis: {
          type: 'value',
          axisLine: { show: false },
          axisTick: { show: false },
          axisLabel: { color: 'rgba(0, 0, 0, 0.6)' },
          splitLine: { lineStyle: { color: 'rgba(0, 0, 0, 0.08)' } },
        },
        series: [
          {
            data: dataUsers,
            type: 'line',
            smooth: true,
            symbol: 'circle',
            symbolSize: 6,
            lineStyle: {
              color: '#010120',
              width: 2,
            },
            itemStyle: {
              color: '#010120',
            },
            areaStyle: {
              color: {
                type: 'linear',
                x: 0,
                y: 0,
                x2: 0,
                y2: 1,
                colorStops: [
                  {
                    offset: 0,
                    color: 'rgba(1, 1, 32, 0.08)',
                  },
                  {
                    offset: 1,
                    color: 'rgba(1, 1, 32, 0)',
                  },
                ],
              },
            },
          },
        ],
      })

      const dataIncome = []
      for (let k in chartData.orders) {
        dataIncome.push(chartData.orders[k])
      }
      chartIncome.setOption({
        tooltip: {
          trigger: 'axis',
          backgroundColor: '#ffffff',
          borderColor: 'rgba(0, 0, 0, 0.08)',
          textStyle: { color: '#010120' },
          extraCssText: 'box-shadow: rgba(1, 1, 32, 0.1) 0px 4px 10px; border-radius: 4px;',
        },
        grid: {
          left: '3%',
          right: '4%',
          bottom: '3%',
          containLabel: true,
        },
        xAxis: {
          type: 'category',
          data: x,
          axisLine: { lineStyle: { color: 'rgba(0, 0, 0, 0.08)' } },
          axisTick: { show: false },
          axisLabel: { color: 'rgba(0, 0, 0, 0.6)' },
        },
        yAxis: {
          type: 'value',
          axisLine: { show: false },
          axisTick: { show: false },
          axisLabel: { color: 'rgba(0, 0, 0, 0.6)' },
          splitLine: { lineStyle: { color: 'rgba(0, 0, 0, 0.08)' } },
        },
        series: [
          {
            data: dataIncome,
            type: 'line',
            smooth: true,
            symbol: 'circle',
            symbolSize: 6,
            lineStyle: {
              color: '#bdbbff',
              width: 2,
            },
            itemStyle: {
              color: '#bdbbff',
            },
            areaStyle: {
              color: {
                type: 'linear',
                x: 0,
                y: 0,
                x2: 0,
                y2: 1,
                colorStops: [
                  {
                    offset: 0,
                    color: 'rgba(189, 187, 255, 0.12)',
                  },
                  {
                    offset: 1,
                    color: 'rgba(189, 187, 255, 0)',
                  },
                ],
              },
            },
          },
        ],
      })
    })
    .catch((e) => {
      ElMessage.error('获取统计数据失败：' + e.message)
    })

  window.onresize = function () {
    // 自适应大小
    if (chartUsers) chartUsers.resize()
    if (chartIncome) chartIncome.resize()
  }
})
</script>

<style scoped lang="scss">
.dashboard {
  padding: 24px;
  background: var(--theme-bg-color);
  min-height: 100vh;

  .stats-row {
    margin-bottom: 24px;
  }

  .content-row {
    margin-bottom: 24px;
  }

  // Together AI contained card：hairline 边框 + 8px 圆角 + 深蓝着色阴影
  .stats-card {
    --el-card-border-color: rgba(0, 0, 0, 0.08);
    --el-card-border-radius: 8px;
    --el-card-bg-color: var(--card-bg);
    border: 1px solid var(--admin-container-border, rgba(0, 0, 0, 0.08));
    border-radius: 8px;
    overflow: hidden;
    transition: box-shadow 0.2s ease;
    background: var(--card-bg);

    &:hover {
      box-shadow: var(--ta-shadow-lg, rgba(1, 1, 32, 0.14) 0px 8px 24px);
    }

    :deep(.el-card__body) {
      padding: 24px 28px;
    }
  }

  // Together AI Stats Bar：Mono 大写小标签在上 → 大数字 → 小说明
  .card-label {
    font-family: var(--ta-font-mono, 'PP Neue Montreal Mono', Georgia, monospace);
    font-size: 11px;
    font-weight: 500;
    letter-spacing: 0.055px;
    text-transform: uppercase;
    color: var(--theme-text-color-secondary);
    margin-bottom: 12px;
  }

  // 大数字：weight 500 负字距，绝不加粗
  .card-number {
    font-size: 40px;
    font-weight: 500;
    color: var(--theme-text-color-primary);
    line-height: 1.1;
    letter-spacing: -1.2px;
    margin-bottom: 8px;
  }

  .card-desc {
    font-size: 13px;
    letter-spacing: -0.16px;
    color: var(--theme-text-color-secondary);
  }

  .chart-card,
  .list-card {
    --el-card-border-color: rgba(0, 0, 0, 0.08);
    --el-card-border-radius: 8px;
    --el-card-bg-color: var(--card-bg);
    border: 1px solid var(--admin-container-border, rgba(0, 0, 0, 0.08));
    border-radius: 8px;
    overflow: hidden;
    height: 100%;
    background: var(--card-bg);

    :deep(.el-card__body) {
      padding: 24px;
      height: 100%;
      display: flex;
      flex-direction: column;
    }
  }

  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;
    padding-bottom: 16px;
    border-bottom: 1px solid var(--admin-container-border, rgba(0, 0, 0, 0.08));

    h3 {
      margin: 0;
      font-size: 18px;
      font-weight: 500;
      letter-spacing: -0.36px;
      color: var(--theme-text-color-primary);
    }
  }

  .order-list,
  .user-list,
  .app-list {
    flex: 1;
    overflow-y: auto;
  }

  .order-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 12px 0;
    border-bottom: 1px solid var(--el-border-color);

    &:last-child {
      border-bottom: none;
    }

    .order-info {
      display: flex;
      align-items: center;
      gap: 12px;

      .order-id {
        font-family: var(--ta-font-mono, 'PP Neue Montreal Mono', Georgia, monospace);
        font-size: 13px;
        color: var(--theme-text-color-secondary);
      }

      .order-amount {
        font-family: var(--ta-font-mono, 'PP Neue Montreal Mono', Georgia, monospace);
        font-size: 14px;
        font-weight: 500;
        color: var(--theme-text-color-primary);
      }
    }

    .order-meta {
      display: flex;
      align-items: center;
      gap: 8px;

      .order-date {
        font-size: 12px;
        color: var(--theme-text-color-secondary);
      }
    }
  }

  .user-item {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 12px 0;
    border-bottom: 1px solid var(--el-border-color);

    &:last-child {
      border-bottom: none;
    }

    .user-info {
      flex: 1;
      min-width: 0;

      .user-name {
        font-size: 14px;
        font-weight: 500;
        color: var(--theme-text-color-primary);
        margin-bottom: 2px;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
      }

      .user-id {
        font-size: 12px;
        color: var(--theme-text-color-secondary);
      }
    }

    .user-meta {
      text-align: right;

      .user-time {
        font-size: 12px;
        color: var(--theme-text-color-secondary);
        margin-bottom: 4px;
      }
    }
  }
}
</style>
