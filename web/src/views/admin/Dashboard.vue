<template>
  <div class="dashboard" v-loading="loading">
    <!-- 统计卡片区域 -->
    <el-row class="stats-row" :gutter="24">
      <el-col :span="6">
        <el-card class="stats-card" shadow="hover">
          <div class="card-content">
            <div class="card-icon user-icon">
              <el-icon><User /></el-icon>
            </div>
            <div class="card-info">
              <div class="card-number">{{ stats.users }}</div>
              <div class="card-label">用户总数</div>
              <div class="card-desc">今日新增: {{ stats.todayUsers || 1 }}</div>
            </div>
          </div>
        </el-card>
      </el-col>

      <el-col :span="6">
        <el-card class="stats-card" shadow="hover">
          <div class="card-content">
            <div class="card-icon chat-icon">
              <el-icon><ChatDotRound /></el-icon>
            </div>
            <div class="card-info">
              <div class="card-number">{{ stats.chats }}</div>
              <div class="card-label">对话总数</div>
              <div class="card-desc">今日: {{ stats.todayChats }}</div>
            </div>
          </div>
        </el-card>
      </el-col>

      <el-col :span="6">
        <el-card class="stats-card" shadow="hover">
          <div class="card-content">
            <div class="card-icon token-icon">
              <el-icon><TrendCharts /></el-icon>
            </div>
            <div class="card-info">
              <div class="card-number">{{ formatNumber(stats.tokens) }}</div>
              <div class="card-label">Token消耗</div>
              <div class="card-desc">今日: {{ formatNumber(stats.todayTokens) }}</div>
            </div>
          </div>
        </el-card>
      </el-col>

      <el-col :span="6">
        <el-card class="stats-card" shadow="hover">
          <div class="card-content">
            <div class="card-icon income-icon">
              <el-icon><Money /></el-icon>
            </div>
            <div class="card-info">
              <div class="card-number">￥{{ stats.income.toFixed(2) }}</div>
              <div class="card-label">总收入</div>
              <div class="card-desc">今日: ￥{{ stats.todayIncome?.toFixed(2) || '0.00' }}</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 第二行统计卡片 -->
    <el-row class="stats-row" :gutter="24">
      <el-col :span="6">
        <el-card class="stats-card" shadow="hover">
          <div class="card-content">
            <div class="card-icon image-icon">
              <el-icon><Picture /></el-icon>
            </div>
            <div class="card-info">
              <div class="card-number">{{ stats.imageJobs }}</div>
              <div class="card-label">图片生成</div>
              <div class="card-desc">今日: {{ stats.todayImageJobs }}</div>
            </div>
          </div>
        </el-card>
      </el-col>

      <el-col :span="6">
        <el-card class="stats-card" shadow="hover">
          <div class="card-content">
            <div class="card-icon video-icon">
              <el-icon><VideoPlay /></el-icon>
            </div>
            <div class="card-info">
              <div class="card-number">{{ stats.videoJobs }}</div>
              <div class="card-label">视频生成</div>
              <div class="card-desc">今日: {{ stats.todayVideoJobs }}</div>
            </div>
          </div>
        </el-card>
      </el-col>

      <el-col :span="6">
        <el-card class="stats-card" shadow="hover">
          <div class="card-content">
            <div class="card-icon music-icon">
              <el-icon><Headset /></el-icon>
            </div>
            <div class="card-info">
              <div class="card-number">{{ stats.musicJobs }}</div>
              <div class="card-label">音乐生成</div>
              <div class="card-desc">今日: {{ stats.todayMusicJobs }}</div>
            </div>
          </div>
        </el-card>
      </el-col>

      <el-col :span="6">
        <el-card class="stats-card" shadow="hover">
          <div class="card-content">
            <div class="card-icon order-icon">
              <el-icon><ShoppingCart /></el-icon>
            </div>
            <div class="card-info">
              <div class="card-number">{{ stats.orders }}</div>
              <div class="card-label">订单总数</div>
              <div class="card-desc">今日: {{ stats.todayOrders }}</div>
            </div>
          </div>
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
      <el-col :span="8">
        <el-card class="list-card" shadow="hover">
          <div class="card-header">
            <h3>最近订单</h3>
          </div>
          <div class="order-list">
            <div class="order-item" v-for="order in recentOrders" :key="order.id">
              <div class="order-info">
                <div class="order-id">#{{ order.id }}</div>
                <div class="order-amount">{{ order.amount }}</div>
              </div>
              <div class="order-meta">
                <div class="order-date">{{ order.date }}</div>
                <el-tag :type="order.status === '已支付' ? 'success' : 'warning'" size="small">
                  {{ order.status }}
                </el-tag>
              </div>
            </div>
          </div>
        </el-card>
      </el-col>

      <!-- 最近用户 -->
      <el-col :span="8">
        <el-card class="list-card" shadow="hover">
          <div class="card-header">
            <h3>最近用户</h3>
          </div>
          <div class="user-list">
            <div class="user-item" v-for="user in recentUsers" :key="user.id">
              <div class="user-avatar">
                <el-avatar :size="40" :src="user.avatar">{{ user.name.charAt(0) }}</el-avatar>
              </div>
              <div class="user-info">
                <div class="user-name">{{ user.name }}</div>
                <div class="user-id">{{ user.userId }}</div>
              </div>
              <div class="user-meta">
                <div class="user-time">{{ user.time }}</div>
                <el-tag type="info" size="small">{{ user.status }}</el-tag>
              </div>
            </div>
          </div>
        </el-card>
      </el-col>

      <!-- 任务统计 -->
      <el-col :span="8">
        <el-card class="list-card" shadow="hover">
          <div class="card-header">
            <h3>AI任务统计</h3>
          </div>
          <div class="job-stats">
            <div class="job-stat-item">
              <div class="stat-icon image-stat">
                <el-icon><Picture /></el-icon>
              </div>
              <div class="stat-info">
                <div class="stat-label">图片生成</div>
                <div class="stat-number">{{ stats.imageJobs }}</div>
              </div>
            </div>
            <div class="job-stat-item">
              <div class="stat-icon video-stat">
                <el-icon><VideoPlay /></el-icon>
              </div>
              <div class="stat-info">
                <div class="stat-label">视频生成</div>
                <div class="stat-number">{{ stats.videoJobs }}</div>
              </div>
            </div>
            <div class="job-stat-item">
              <div class="stat-icon music-stat">
                <el-icon><Headset /></el-icon>
              </div>
              <div class="stat-info">
                <div class="stat-label">音乐生成</div>
                <div class="stat-number">{{ stats.musicJobs }}</div>
              </div>
            </div>
            <div class="job-stat-item">
              <div class="stat-icon order-stat">
                <el-icon><ShoppingCart /></el-icon>
              </div>
              <div class="stat-info">
                <div class="stat-label">订单总数</div>
                <div class="stat-number">{{ stats.orders }}</div>
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

onMounted(() => {
  const chartUsers = echarts.init(document.getElementById('chart-users'))
  const chartTokens = echarts.init(document.getElementById('chart-tokens'))
  const chartIncome = echarts.init(document.getElementById('chart-income'))
  httpGet('/api/admin/dashboard/stats')
    .then((res) => {
      // 更新统计数据
      Object.assign(stats.value, res.data)
      const chartData = res.data.chart
      loading.value = false

      const x = []
      const dataUsers = []
      for (let k in chartData.users) {
        x.push(k)
        dataUsers.push(chartData.users[k])
      }
      chartUsers.setOption({
        tooltip: {
          trigger: 'axis',
          backgroundColor: 'rgba(255, 255, 255, 0.9)',
          borderColor: '#ddd',
          textStyle: { color: '#666' },
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
          axisLine: { lineStyle: { color: '#ddd' } },
          axisTick: { show: false },
          axisLabel: { color: '#999' },
        },
        yAxis: {
          type: 'value',
          axisLine: { show: false },
          axisTick: { show: false },
          axisLabel: { color: '#999' },
          splitLine: { lineStyle: { color: '#f0f0f0' } },
        },
        series: [
          {
            data: dataUsers,
            type: 'line',
            smooth: true,
            symbol: 'circle',
            symbolSize: 6,
            lineStyle: {
              color: '#8B5CF6',
              width: 3,
            },
            itemStyle: {
              color: '#8B5CF6',
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
                    color: 'rgba(139, 92, 246, 0.3)',
                  },
                  {
                    offset: 1,
                    color: 'rgba(139, 92, 246, 0.05)',
                  },
                ],
              },
            },
          },
        ],
      })
      const dataTokens = []
      for (let k in chartData.historyMessage) {
        dataTokens.push(chartData.historyMessage[k])
      }
      chartTokens.setOption({
        xAxis: {
          data: x,
        },
        yAxis: {},
        series: [
          {
            data: dataTokens,
            type: 'line',
            label: {
              show: true,
              position: 'bottom',
              textStyle: {
                fontSize: 18,
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
          backgroundColor: 'rgba(255, 255, 255, 0.9)',
          borderColor: '#ddd',
          textStyle: { color: '#666' },
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
          axisLine: { lineStyle: { color: '#ddd' } },
          axisTick: { show: false },
          axisLabel: { color: '#999' },
        },
        yAxis: {
          type: 'value',
          axisLine: { show: false },
          axisTick: { show: false },
          axisLabel: { color: '#999' },
          splitLine: { lineStyle: { color: '#f0f0f0' } },
        },
        series: [
          {
            data: dataIncome,
            type: 'line',
            smooth: true,
            symbol: 'circle',
            symbolSize: 6,
            lineStyle: {
              color: '#10B981',
              width: 3,
            },
            itemStyle: {
              color: '#10B981',
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
                    color: 'rgba(16, 185, 129, 0.3)',
                  },
                  {
                    offset: 1,
                    color: 'rgba(16, 185, 129, 0.05)',
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
    chartUsers.resize()
    chartIncome.resize()
  }
})
</script>

<style scoped lang="stylus">
.dashboard {
  padding: 24px;
  background: #f8fafc;
  min-height: 100vh;

  .stats-row {
    margin-bottom: 24px;
  }

  .content-row {
    margin-bottom: 24px;
  }

  .stats-card {
    border: none;
    border-radius: 12px;
    overflow: hidden;
    transition: all 0.3s ease;

    &:hover {
      transform: translateY(-2px);
      box-shadow: 0 8px 25px rgba(0, 0, 0, 0.1);
    }

    :deep(.el-card__body) {
      padding: 24px;
    }
  }

  .card-content {
    display: flex;
    align-items: center;
    gap: 16px;
  }

  .card-icon {
    width: 64px;
    height: 64px;
    border-radius: 16px;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 28px;
    color: white;
    flex-shrink: 0;

    &.user-icon {
      background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    }

    &.chat-icon {
      background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
    }

    &.token-icon {
      background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    }

    &.income-icon {
      background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%);
    }

    &.order-icon {
      background: linear-gradient(135deg, #fa709a 0%, #fee140 100%);
    }

    &.image-icon {
      background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    }

    &.video-icon {
      background: linear-gradient(135deg, #fa709a 0%, #fee140 100%);
    }

    &.music-icon {
      background: linear-gradient(135deg, #a8edea 0%, #fed6e3 100%);
    }
  }

  .card-info {
    flex: 1;
    min-width: 0;
  }

  .card-number {
    font-size: 32px;
    font-weight: 700;
    color: #1f2937;
    margin-bottom: 4px;
    line-height: 1;
  }

  .card-label {
    font-size: 16px;
    color: #6b7280;
    margin-bottom: 4px;
    font-weight: 500;
  }

  .card-desc {
    font-size: 14px;
    color: #9ca3af;
  }

  .chart-card, .list-card {
    border: none;
    border-radius: 12px;
    overflow: hidden;
    height: 100%;

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
    border-bottom: 1px solid #f3f4f6;

    h3 {
      margin: 0;
      font-size: 18px;
      font-weight: 600;
      color: #1f2937;
    }
  }

  .order-list, .user-list, .app-list {
    flex: 1;
    overflow-y: auto;
  }

  .order-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 12px 0;
    border-bottom: 1px solid #f3f4f6;

    &:last-child {
      border-bottom: none;
    }

    .order-info {
      display: flex;
      align-items: center;
      gap: 12px;

      .order-id {
        font-size: 14px;
        color: #6b7280;
      }

      .order-amount {
        font-size: 16px;
        font-weight: 600;
        color: #10b981;
      }
    }

    .order-meta {
      display: flex;
      align-items: center;
      gap: 8px;

      .order-date {
        font-size: 12px;
        color: #9ca3af;
      }
    }
  }

  .user-item {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 12px 0;
    border-bottom: 1px solid #f3f4f6;

    &:last-child {
      border-bottom: none;
    }

    .user-info {
      flex: 1;
      min-width: 0;

      .user-name {
        font-size: 14px;
        font-weight: 500;
        color: #1f2937;
        margin-bottom: 2px;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
      }

      .user-id {
        font-size: 12px;
        color: #9ca3af;
      }
    }

    .user-meta {
      text-align: right;

      .user-time {
        font-size: 12px;
        color: #6b7280;
        margin-bottom: 4px;
      }
    }
  }

  .job-stats {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 16px;

    .job-stat-item {
      display: flex;
      align-items: center;
      gap: 12px;
      padding: 16px;
      background: #f8fafc;
      border-radius: 8px;
      transition: all 0.3s ease;

      &:hover {
        background: #f1f5f9;
      }

      .stat-icon {
        width: 40px;
        height: 40px;
        border-radius: 8px;
        display: flex;
        align-items: center;
        justify-content: center;
        font-size: 20px;
        color: white;
        flex-shrink: 0;

        &.image-stat {
          background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
        }

        &.video-stat {
          background: linear-gradient(135deg, #fa709a 0%, #fee140 100%);
        }

        &.music-stat {
          background: linear-gradient(135deg, #a8edea 0%, #fed6e3 100%);
        }

        &.order-stat {
          background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%);
        }
      }

      .stat-info {
        flex: 1;

        .stat-label {
          font-size: 12px;
          color: #6b7280;
          margin-bottom: 4px;
        }

        .stat-number {
          font-size: 18px;
          font-weight: 600;
          color: #1f2937;
        }
      }
    }
  }
}
</style>
