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
      <div class="filter-bar">
        <CustomTabs
          :model-value="activeType"
          @update:model-value="activeType = $event"
          @tab-click="onTypeChange"
        >
          <CustomTabPane name="all" label="全部" />
          <CustomTabPane name="chat" label="对话" />
          <CustomTabPane name="image" label="绘画" />
          <CustomTabPane name="music" label="音乐" />
          <CustomTabPane name="video" label="视频" />
        </CustomTabs>
      </div>

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
                  <div class="log-title">{{ item.title }}</div>
                  <div class="log-time">{{ formatTime(item.created_at) }}</div>
                </div>
                <div class="log-cost">
                  <span class="cost-value">-{{ item.cost }}</span>
                  <span class="cost-unit">算力</span>
                </div>
              </div>
              <div class="log-detail" v-if="item.remark">
                <van-text-ellipsis :content="item.remark" expand-text="展开" collapse-text="收起" />
              </div>
            </div>

            <!-- 空状态 -->
            <van-empty v-if="!loading && logList.length === 0" description="暂无消费记录" />
          </van-list>
        </van-pull-refresh>
      </div>
    </div>

    <!-- 筛选弹窗 -->
    <van-action-sheet
      :model-value="showFilter"
      @update:model-value="showFilter = $event"
      title="筛选条件"
    >
      <div class="filter-content">
        <van-form>
          <van-field label="时间范围">
            <template #input>
              <van-button size="small" @click="showDatePicker = true">
                {{
                  dateRange.start && dateRange.end
                    ? `${dateRange.start} 至 ${dateRange.end}`
                    : '选择时间'
                }}
              </van-button>
            </template>
          </van-field>
          <van-field label="消费类型">
            <template #input>
              <van-radio-group v-model="filterType" direction="horizontal">
                <van-radio name="all">全部</van-radio>
                <van-radio name="chat">对话</van-radio>
                <van-radio name="image">绘画</van-radio>
                <van-radio name="music">音乐</van-radio>
              </van-radio-group>
            </template>
          </van-field>
        </van-form>
        <div class="filter-actions">
          <van-button @click="resetFilter">重置</van-button>
          <van-button type="primary" @click="applyFilter">确定</van-button>
        </div>
      </div>
    </van-action-sheet>

    <!-- 日期选择器 -->
    <van-calendar
      :model-value="showDatePicker"
      @update:model-value="showDatePicker = $event"
      type="range"
      @confirm="onDateConfirm"
      :max-date="new Date()"
    />
  </div>
</template>

<script setup>
import CustomTabPane from '@/components/ui/CustomTabPane.vue'
import CustomTabs from '@/components/ui/CustomTabs.vue'
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const loading = ref(false)
const finished = ref(false)
const refreshing = ref(false)
const logList = ref([])
const activeType = ref('all')
const showFilter = ref(false)
const showDatePicker = ref(false)
const filterType = ref('all')
const dateRange = ref({
  start: '',
  end: '',
})

// 统计数据
const stats = ref({
  total: 0,
  today: 0,
  balance: 0,
})

// 分页参数
const pageParams = ref({
  page: 1,
  limit: 20,
  type: 'all',
})

onMounted(() => {
  fetchStats()
  onLoad()
})

// 获取统计数据
const fetchStats = () => {
  // 这里应该调用实际的API
  // httpGet('/api/user/power/stats').then(res => {
  //   stats.value = res.data
  // })

  // 临时使用模拟数据
  stats.value = {
    total: Math.floor(Math.random() * 10000),
    today: Math.floor(Math.random() * 100),
    balance: Math.floor(Math.random() * 1000),
  }
}

// 加载日志列表
const onLoad = () => {
  if (finished.value) return

  loading.value = true

  // 模拟API调用
  setTimeout(() => {
    const mockData = generateMockData(pageParams.value.page, pageParams.value.limit)

    if (pageParams.value.page === 1) {
      logList.value = mockData
    } else {
      logList.value.push(...mockData)
    }

    loading.value = false
    pageParams.value.page++

    // 模拟数据加载完成
    if (pageParams.value.page > 5) {
      finished.value = true
    }
  }, 1000)
}

// 下拉刷新
const onRefresh = () => {
  finished.value = false
  pageParams.value.page = 1
  refreshing.value = true

  setTimeout(() => {
    logList.value = generateMockData(1, pageParams.value.limit)
    refreshing.value = false
    pageParams.value.page = 2
  }, 1000)
}

// 类型切换
const onTypeChange = (type) => {
  pageParams.value.type = type
  pageParams.value.page = 1
  finished.value = false
  logList.value = []
  onLoad()
}

// 生成模拟数据
const generateMockData = (page, limit) => {
  const types = ['chat', 'image', 'music', 'video']
  const titles = {
    chat: ['GPT-4对话', 'Claude对话', '智能助手'],
    image: ['MidJourney生成', 'Stable Diffusion', 'DALL-E创作'],
    music: ['Suno音乐创作', '音频生成'],
    video: ['视频生成', 'Luma创作'],
  }

  const data = []
  const startIndex = (page - 1) * limit

  for (let i = 0; i < limit; i++) {
    const id = startIndex + i + 1
    const type = types[Math.floor(Math.random() * types.length)]
    const title = titles[type][Math.floor(Math.random() * titles[type].length)]

    // 如果有类型筛选且不匹配，跳过
    if (pageParams.value.type !== 'all' && type !== pageParams.value.type) {
      continue
    }

    data.push({
      id,
      type,
      title,
      cost: Math.floor(Math.random() * 50) + 1,
      remark: Math.random() > 0.5 ? '消费详情：使用高级模型进行AI创作，效果优质' : '',
      created_at: new Date(Date.now() - Math.random() * 7 * 24 * 60 * 60 * 1000).toISOString(),
    })
  }

  return data
}

// 获取类型图标
const getTypeIcon = (type) => {
  const icons = {
    chat: 'icon-chat',
    image: 'icon-mj',
    music: 'icon-music',
    video: 'icon-video',
  }
  return icons[type] || 'icon-chat'
}

// 获取类型颜色
const getTypeColor = (type) => {
  const colors = {
    chat: '#1989fa',
    image: '#8B5CF6',
    music: '#ee0a24',
    video: '#07c160',
  }
  return colors[type] || '#1989fa'
}

// 格式化时间
const formatTime = (timeStr) => {
  const date = new Date(timeStr)
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

// 日期选择确认
const onDateConfirm = (values) => {
  const [start, end] = values
  dateRange.value = {
    start: start.toLocaleDateString(),
    end: end.toLocaleDateString(),
  }
  showDatePicker.value = false
}

// 重置筛选
const resetFilter = () => {
  filterType.value = 'all'
  dateRange.value = { start: '', end: '' }
}

// 应用筛选
const applyFilter = () => {
  activeType.value = filterType.value
  pageParams.value.type = filterType.value
  pageParams.value.page = 1
  finished.value = false
  logList.value = []
  showFilter.value = false
  onLoad()
}
</script>

<style lang="scss" scoped>
.power-log {
  min-height: 100vh;
  background: var(--van-background);

  .power-content {
    padding-top: 46px;

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

      :deep(.van-tabs__nav) {
        padding: 0 16px;
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
