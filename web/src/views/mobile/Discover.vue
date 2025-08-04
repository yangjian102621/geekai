<template>
  <div class="discover-page">
    <div class="discover-content">
      <!-- 功能分类 -->
      <div class="category-section">
        <h3 class="category-title">AI 工具</h3>
        <van-row :gutter="12">
          <van-col :span="6" v-for="tool in aiTools" :key="tool.key">
            <div class="tool-card" @click="navigateTo(tool.url)">
              <div class="tool-icon" :style="{ backgroundColor: tool.color }">
                <i class="iconfont" :class="tool.icon"></i>
              </div>
              <div class="tool-name">{{ tool.name }}</div>
            </div>
          </van-col>
        </van-row>
      </div>

      <!-- 用户服务 -->
      <div class="category-section">
        <h3 class="category-title">我的服务</h3>
        <van-cell-group inset>
          <van-cell
            v-for="service in userServices"
            :key="service.key"
            :title="service.name"
            :value="service.desc"
            :icon="service.icon"
            is-link
            @click="navigateTo(service.url)"
          >
            <template #icon>
              <i class="iconfont" :class="service.icon" :style="{ color: service.color }"></i>
            </template>
          </van-cell>
        </van-cell-group>
      </div>

      <!-- 实用功能 -->
      <div class="category-section">
        <h3 class="category-title">实用功能</h3>
        <van-cell-group inset>
          <van-cell
            v-for="utility in utilities"
            :key="utility.key"
            :title="utility.name"
            :value="utility.desc"
            is-link
            @click="navigateTo(utility.url)"
          >
            <template #icon>
              <i class="iconfont" :class="utility.icon" :style="{ color: utility.color }"></i>
            </template>
          </van-cell>
        </van-cell-group>
      </div>

      <!-- 推荐内容 -->
      <div class="category-section">
        <h3 class="category-title">精选推荐</h3>
        <van-grid :column-num="2" :gutter="12" :border="false">
          <van-grid-item
            v-for="item in recommendations"
            :key="item.key"
            @click="navigateTo(item.url)"
            class="recommend-item"
          >
            <div class="recommend-card">
              <div class="recommend-image">
                <van-image :src="item.image" fit="cover" />
              </div>
              <div class="recommend-info">
                <div class="recommend-title">{{ item.title }}</div>
                <div class="recommend-desc">{{ item.desc }}</div>
              </div>
            </div>
          </van-grid-item>
        </van-grid>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

// AI工具配置
const aiTools = ref([
  { key: 'mj', name: 'MJ绘画', icon: 'icon-mj', color: '#8B5CF6', url: '/mobile/create?tab=mj' },
  { key: 'sd', name: 'SD绘画', icon: 'icon-sd', color: '#06B6D4', url: '/mobile/create?tab=sd' },
  {
    key: 'dalle',
    name: 'DALL·E',
    icon: 'icon-dalle',
    color: '#F59E0B',
    url: '/mobile/create?tab=dalle',
  },
  {
    key: 'suno',
    name: '音乐创作',
    icon: 'icon-music',
    color: '#EF4444',
    url: '/mobile/create?tab=suno',
  },
  {
    key: 'video',
    name: '视频生成',
    icon: 'icon-video',
    color: '#10B981',
    url: '/mobile/create?tab=video',
  },
  {
    key: 'jimeng',
    name: '即梦AI',
    icon: 'icon-jimeng',
    color: '#F97316',
    url: '/mobile/create?tab=jimeng',
  },
  {
    key: 'xmind',
    name: '思维导图',
    icon: 'icon-mind',
    color: '#3B82F6',
    url: '/mobile/tools?tab=xmind',
  },
  { key: 'apps', name: '应用中心', icon: 'icon-apps', color: '#EC4899', url: '/mobile/apps' },
])

// 用户服务
const userServices = ref([
  {
    key: 'member',
    name: '会员中心',
    desc: '充值升级享受更多权益',
    icon: 'icon-vip',
    color: '#FFD700',
    url: '/mobile/member',
  },
  {
    key: 'powerLog',
    name: '消费记录',
    desc: '查看算力使用详情',
    icon: 'icon-history',
    color: '#10B981',
    url: '/mobile/power-log',
  },
  {
    key: 'invite',
    name: '邀请好友',
    desc: '推广获取奖励',
    icon: 'icon-user-plus',
    color: '#F59E0B',
    url: '/mobile/invite',
  },
  {
    key: 'export',
    name: '导出对话',
    desc: '保存聊天记录',
    icon: 'icon-download',
    color: '#06B6D4',
    url: '/mobile/chat/export',
  },
])

// 实用功能
const utilities = ref([
  {
    key: 'imgWall',
    name: '作品展示',
    desc: '浏览精美AI作品',
    icon: 'icon-gallery',
    color: '#EC4899',
    url: '/mobile/imgWall',
  },
  {
    key: 'settings',
    name: '设置中心',
    desc: '个性化配置',
    icon: 'icon-setting',
    color: '#6B7280',
    url: '/mobile/settings',
  },
  {
    key: 'help',
    name: '帮助中心',
    desc: '使用指南和常见问题',
    icon: 'icon-help',
    color: '#8B5CF6',
    url: '/mobile/help',
  },
  {
    key: 'feedback',
    name: '意见反馈',
    desc: '提出建议和问题',
    icon: 'icon-message',
    color: '#EF4444',
    url: '/mobile/feedback',
  },
])

// 推荐内容
const recommendations = ref([
  {
    key: 'new-features',
    title: '新功能发布',
    desc: '体验最新AI创作工具',
    image: '/images/recommend/new-features.jpg',
    url: '/mobile/news',
  },
  {
    key: 'tutorials',
    title: '使用教程',
    desc: '快速上手AI创作',
    image: '/images/recommend/tutorials.jpg',
    url: '/mobile/tutorials',
  },
  {
    key: 'gallery',
    title: '精选作品',
    desc: '欣赏优秀AI作品',
    image: '/images/recommend/gallery.jpg',
    url: '/mobile/imgWall',
  },
  {
    key: 'community',
    title: '用户社区',
    desc: '交流创作心得',
    image: '/images/recommend/community.jpg',
    url: '/mobile/community',
  },
])

// 导航处理
const navigateTo = (url) => {
  if (url.startsWith('http')) {
    window.open(url, '_blank')
  } else {
    router.push(url)
  }
}
</script>

<style lang="scss" scoped>
.discover-page {
  min-height: 100vh;
  background: var(--van-background);

  .nav-left {
    display: flex;
    align-items: center;

    .iconfont {
      font-size: 20px;
      color: var(--van-primary-color);
    }
  }

  .discover-content {
    .category-section {
      margin-bottom: 24px;

      .category-title {
        font-size: 18px;
        font-weight: 600;
        color: var(--van-text-color);
        margin: 0 0 16px 0;
        padding-left: 4px;
      }
    }

    // AI工具卡片
    .tool-card {
      display: flex;
      flex-direction: column;
      align-items: center;
      padding: 16px 8px;
      background: var(--van-cell-background);
      border-radius: 12px;
      box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
      transition: all 0.3s ease;
      cursor: pointer;

      &:active {
        transform: scale(0.95);
      }

      .tool-icon {
        width: 44px;
        height: 44px;
        border-radius: 12px;
        display: flex;
        align-items: center;
        justify-content: center;
        margin-bottom: 8px;
        box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);

        .iconfont {
          font-size: 22px;
          color: white;
        }
      }

      .tool-name {
        font-size: 12px;
        font-weight: 500;
        color: var(--van-text-color);
        text-align: center;
      }
    }

    // 服务列表
    :deep(.van-cell-group) {
      border-radius: 12px;
      overflow: hidden;
      box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);

      .van-cell {
        padding: 16px;

        .van-cell__title {
          font-weight: 500;
        }

        .van-cell__value {
          color: var(--van-gray-6);
          font-size: 13px;
        }

        .iconfont {
          font-size: 20px;
          margin-right: 12px;
        }
      }
    }

    // 推荐内容
    .recommend-item {
      .recommend-card {
        background: var(--van-cell-background);
        border-radius: 12px;
        overflow: hidden;
        box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
        transition: all 0.3s ease;
        cursor: pointer;

        &:active {
          transform: scale(0.98);
        }

        .recommend-image {
          height: 120px;
          overflow: hidden;

          :deep(.van-image) {
            width: 100%;
            height: 100%;
          }
        }

        .recommend-info {
          padding: 12px;

          .recommend-title {
            font-size: 14px;
            font-weight: 600;
            color: var(--van-text-color);
            margin-bottom: 4px;
          }

          .recommend-desc {
            font-size: 12px;
            color: var(--van-gray-6);
            line-height: 1.4;
          }
        }
      }
    }
  }
}

// 深色主题优化
:deep(.van-theme-dark) {
  .discover-page {
    .tool-card {
      box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);

      .tool-icon {
        box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
      }
    }

    .van-cell-group {
      box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
    }

    .recommend-item .recommend-card {
      box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
    }
  }
}
</style>
