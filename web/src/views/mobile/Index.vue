<template>
  <div class="index">
    <div class="header">
      <div class="user-greeting px-3">
        <div class="greeting-text">
          <h2 class="title">{{ getGreeting() }}</h2>
          <p class="subtitle">{{ title }}</p>
        </div>
        <div class="user-avatar" v-if="isLogin" @click="router.push('profile')">
          <van-image :src="userAvatar" round width="40" height="40" />
        </div>
        <div class="login-btn" v-else @click="router.push('login')">
          <van-button size="small" type="primary" round>登录</van-button>
        </div>
      </div>
    </div>

    <!-- 快捷操作区 -->
    <div class="quick-actions mb-6 px-3">
      <van-row :gutter="12">
        <van-col :span="12">
          <div class="action-card primary" @click="router.push('chat')">
            <div class="action-content">
              <i class="iconfont icon-chat action-icon"></i>
              <div class="action-text">
                <div class="action-title">AI 对话</div>
                <div class="action-desc">智能助手对话</div>
              </div>
            </div>
          </div>
        </van-col>
        <van-col :span="12">
          <div class="action-card secondary" @click="router.push('create')">
            <div class="action-content">
              <i class="iconfont icon-mj action-icon"></i>
              <div class="action-text">
                <div class="action-title">AI 创作</div>
                <div class="action-desc">图像音视频生成</div>
              </div>
            </div>
          </div>
        </van-col>
      </van-row>
    </div>

    <!-- 功能网格 -->
    <div class="feature-section mb-6 px-3">
      <div class="section-header">
        <h3 class="section-title">AI 功能</h3>
      </div>
      <van-grid :column-num="4" :border="false">
        <van-grid-item
          v-for="feature in features"
          :key="feature.key"
          @click="navigateToFeature(feature)"
          class="feature-item"
        >
          <template #icon>
            <div class="feature-icon" :style="{ backgroundColor: feature.color }">
              <i class="iconfont" :class="feature.icon"></i>
            </div>
          </template>
          <template #text>
            <div class="feature-text">{{ feature.name }}</div>
          </template>
        </van-grid-item>
      </van-grid>
    </div>

    <!-- 推荐应用 -->
    <div class="apps-section px-3">
      <div class="section-header">
        <h3 class="section-title">推荐应用</h3>
        <van-button
          class="more-btn"
          size="small"
          icon="arrow"
          type="primary"
          plain
          round
          @click="router.push('apps')"
        >
          更多
        </van-button>
      </div>

      <div class="app-list">
        <van-swipe :autoplay="3000" :show-indicators="false" class="app-swipe">
          <van-swipe-item v-for="chunk in appChunks" :key="chunk[0] && chunk[0].id">
            <div class="app-row px-3">
              <div v-for="item in chunk" :key="item.id" class="app-item" @click="useRole(item.id)">
                <div class="app-avatar">
                  <van-image :src="item.icon" round fit="cover" />
                </div>
                <div class="app-info">
                  <div class="app-name">{{ item.name }}</div>
                  <div class="app-desc">{{ item.intro }}</div>
                </div>
                <div class="app-action">
                  <!-- <van-button
                    size="mini"
                    type="primary"
                    plain
                    round
                    @click.stop="updateRole(item, hasRole(item.key) ? 'remove' : 'add')"
                  >
                    {{ hasRole(item.key) ? '已添加' : '添加' }}
                  </van-button> -->

                  <van-button size="small" type="primary" round> 开始对话 </van-button>
                </div>
              </div>
            </div>
          </van-swipe-item>
        </van-swipe>
      </div>
    </div>
  </div>
</template>

<script setup>
import { checkSession, getSystemInfo } from '@/store/cache'
import { httpGet, httpPost } from '@/utils/http'
import { arrayContains, removeArrayItem, showLoginDialog, substr } from '@/utils/libs'
import { ElMessage } from 'element-plus'
import { showNotify } from 'vant'
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'

const title = ref(import.meta.env.VITE_TITLE)
const router = useRouter()
const isLogin = ref(false)
const apps = ref([])
const loading = ref(false)
const roles = ref([])
const userAvatar = ref('/images/avatar/default.jpg')

// 功能配置
const features = ref([
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
  {
    key: 'imgWall',
    name: '作品展示',
    icon: 'icon-gallery',
    color: '#EC4899',
    url: '/mobile/imgWall',
  },
])

// 应用分组显示（每行2个）
const appChunks = computed(() => {
  const chunks = []
  const displayApps = apps.value.slice(0, 12) // 只显示前6个
  for (let i = 0; i < displayApps.length; i += 4) {
    chunks.push(displayApps.slice(i, i + 4))
  }
  return chunks
})

// 获取问候语
const getGreeting = () => {
  const hour = new Date().getHours()
  if (hour < 6) return '夜深了'
  if (hour < 12) return '早上好'
  if (hour < 18) return '下午好'
  return '晚上好'
}

// 导航到功能页面
const navigateToFeature = (feature) => {
  router.push(feature.url)
}

onMounted(() => {
  getSystemInfo()
    .then((res) => {
      title.value = res.data.title
    })
    .catch((e) => {
      ElMessage.error('获取系统配置失败：' + e.message)
    })

  checkSession()
    .then((user) => {
      isLogin.value = true
      roles.value = user.chat_roles
      userAvatar.value = user.avatar || '/images/avatar/default.jpg'
    })
    .catch(() => {})

  fetchApps()
})

const fetchApps = () => {
  httpGet('/api/app/list')
    .then((res) => {
      const items = res.data
      // 处理 hello message
      for (let i = 0; i < items.length; i++) {
        items[i].intro = substr(items[i].hello_msg, 30)
      }
      apps.value = items
    })
    .catch((e) => {
      showNotify({ type: 'danger', message: '获取应用失败：' + e.message })
    })
}

const updateRole = (row, opt) => {
  if (!isLogin.value) {
    return showLoginDialog(router)
  }

  let actionTitle = ''
  if (opt === 'add') {
    actionTitle = '添加应用'
    const exists = arrayContains(roles.value, row.key)
    if (exists) {
      return
    }
    roles.value.push(row.key)
  } else {
    actionTitle = '移除应用'
    const exists = arrayContains(roles.value, row.key)
    if (!exists) {
      return
    }
    roles.value = removeArrayItem(roles.value, row.key)
  }
  httpPost('/api/app/update', { keys: roles.value })
    .then(() => {
      showNotify({ type: 'success', message: actionTitle + '成功！', duration: 1000 })
    })
    .catch((e) => {
      showNotify({ type: 'danger', message: actionTitle + '失败：' + e.message })
    })
}

const hasRole = (roleKey) => {
  return arrayContains(roles.value, roleKey, (v1, v2) => v1 === v2)
}

const useRole = (roleId) => {
  if (!isLogin.value) {
    return showLoginDialog(router)
  }
  router.push(`/mobile/chat/session?role_id=${roleId}`)
}
</script>

<style scoped lang="scss">
.index {
  color: var(--van-text-color);
  background: linear-gradient(135deg, var(--van-background), var(--van-background-2));
  min-height: 100vh;
  padding: 0;
  .header {
    padding: 20px 0 16px;
    position: sticky;
    top: 0;
    z-index: 100;
    background: inherit;
    backdrop-filter: blur(10px);

    .user-greeting {
      display: flex;
      justify-content: space-between;
      align-items: center;

      .greeting-text {
        flex: 1;

        .title {
          font-size: 28px;
          font-weight: 700;
          color: var(--van-text-color);
          margin: 0 0 4px 0;
          background: linear-gradient(135deg, var(--van-primary-color), #8b5cf6);
          -webkit-background-clip: text;
          -webkit-text-fill-color: transparent;
          background-clip: text;
        }

        .subtitle {
          font-size: 14px;
          color: var(--van-gray-6);
          margin: 0;
        }
      }

      .user-avatar,
      .login-btn {
        flex-shrink: 0;
        margin-left: 16px;
      }
    }
  }

  .quick-actions {
    .action-card {
      border-radius: 16px;
      padding: 20px;
      background: var(--van-cell-background);
      box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
      transition: all 0.3s ease;
      cursor: pointer;

      &:active {
        transform: scale(0.98);
      }

      &.primary {
        background: linear-gradient(135deg, var(--van-primary-color), #8b5cf6);
        color: white;

        .action-icon,
        .action-title,
        .action-desc {
          color: white;
        }
      }

      &.secondary {
        background: linear-gradient(135deg, #06b6d4, #10b981);
        color: white;

        .action-icon,
        .action-title,
        .action-desc {
          color: white;
        }
      }

      .action-content {
        display: flex;
        align-items: center;

        .action-icon {
          font-size: 32px;
          margin-right: 16px;
          opacity: 0.9;
        }

        .action-text {
          .action-title {
            font-size: 18px;
            font-weight: 600;
            margin-bottom: 4px;
          }

          .action-desc {
            font-size: 13px;
            opacity: 0.8;
          }
        }
      }
    }
  }

  .feature-section,
  .apps-section {
    .section-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 16px;

      .section-title {
        font-size: 20px;
        font-weight: 700;
        color: var(--van-text-color);
        margin: 0;
      }

      .more-btn {
        padding: 6px 12px;
        font-size: 12px;
      }
    }
  }

  .feature-section {
    .feature-item {
      padding: 16px 8px;
      transition: all 0.3s ease;

      &:active {
        transform: scale(0.95);
      }

      .feature-icon {
        width: 48px;
        height: 48px;
        border-radius: 12px;
        display: flex;
        align-items: center;
        justify-content: center;
        margin: 0 auto 8px;
        box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);

        i {
          font-size: 24px;
          color: white;
        }
      }

      .feature-text {
        font-size: 12px;
        font-weight: 500;
        color: var(--van-text-color);
        text-align: center;
      }
    }
  }

  .apps-section {
    .app-swipe {
      margin: 0 -4px;

      .app-row {
        display: flex;
        flex-direction: column;
        gap: 12px;
        padding: 0 4px;

        .app-item {
          background: var(--van-cell-background);
          border-radius: 12px;
          padding: 16px;
          display: flex;
          align-items: center;
          box-shadow: 0 2px 12px rgba(0, 0, 0, 0.06);
          transition: all 0.3s ease;
          cursor: pointer;

          &:active {
            transform: scale(0.98);
          }

          .app-avatar {
            width: 44px;
            height: 44px;
            margin-right: 12px;
            flex-shrink: 0;

            :deep(.van-image) {
              width: 100%;
              height: 100%;
            }
          }

          .app-info {
            flex: 1;
            min-width: 0;

            .app-name {
              font-size: 15px;
              font-weight: 600;
              color: var(--van-text-color);
              margin-bottom: 2px;
              overflow: hidden;
              text-overflow: ellipsis;
              white-space: nowrap;
            }

            .app-desc {
              font-size: 12px;
              color: var(--van-gray-6);
              overflow: hidden;
              text-overflow: ellipsis;
              white-space: nowrap;
            }
          }

          .app-action {
            flex-shrink: 0;
            margin-left: 8px;
          }
        }
      }
    }
  }
}

// 响应式调整
@media (max-width: 375px) {
  .index {
    padding: 0 12px 60px;

    .header .user-greeting .greeting-text .title {
      font-size: 24px;
    }

    .quick-actions .action-card {
      padding: 16px;

      .action-content .action-icon {
        font-size: 28px;
        margin-right: 12px;
      }
    }
  }
}

// 深色主题优化
:deep(.van-theme-dark) {
  .index {
    .quick-actions .action-card {
      &.primary,
      &.secondary {
        box-shadow: 0 4px 20px rgba(0, 0, 0, 0.3);
      }
    }

    .feature-section .feature-item .feature-icon {
      box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
    }

    .apps-section .app-swipe .app-row .app-item {
      box-shadow: 0 2px 12px rgba(0, 0, 0, 0.2);
    }
  }
}
</style>
