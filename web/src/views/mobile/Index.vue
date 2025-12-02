<template>
  <div class="index container">
    <div class="header">
      <h2 class="title">{{ title }}</h2>
    </div>

    <div class="content mb-8">
      <div class="feature-grid">
        <van-grid :column-num="3" :gutter="15" border>
          <van-grid-item @click="router.push('chat')" class="feature-item">
            <template #icon>
              <div class="feature-icon">
                <i class="iconfont icon-chat"></i>
              </div>
            </template>
            <template #text>
              <div class="text">AI 对话</div>
            </template>
          </van-grid-item>

          <van-grid-item @click="router.push('image')" class="feature-item">
            <template #icon>
              <div class="feature-icon">
                <i class="iconfont icon-mj"></i>
              </div>
            </template>
            <template #text>
              <div class="text">AI 绘画</div>
            </template>
          </van-grid-item>

          <van-grid-item @click="router.push('imgWall')" class="feature-item">
            <template #icon>
              <div class="feature-icon">
                <van-icon name="photo-o" />
              </div>
            </template>
            <template #text>
              <div class="text">AI 画廊</div>
            </template>
          </van-grid-item>
        </van-grid>
      </div>

      <div class="section-header">
        <h3 class="section-title">推荐应用</h3>
        <van-button class="more-btn" size="small" icon="arrow" @click="router.push('apps')"
          >更多</van-button
        >
      </div>

      <div class="app-list">
        <van-list v-model="loading" :finished="true" finished-text="" @load="fetchApps">
          <van-cell v-for="item in displayApps" :key="item.id" class="app-cell">
            <div class="app-card">
              <div class="app-info">
                <div class="app-image">
                  <van-image :src="item.icon" round />
                </div>
                <div class="app-detail">
                  <div class="app-title">{{ item.name }}</div>
                  <div class="app-desc">{{ item.hello_msg }}</div>
                </div>
              </div>

              <div class="app-actions">
                <van-button size="small" type="primary" class="action-btn" @click="useRole(item.id)"
                  >对话</van-button
                >
                <van-button
                  size="small"
                  :type="hasRole(item.key) ? 'danger' : 'success'"
                  class="action-btn"
                  @click="updateRole(item, hasRole(item.key) ? 'remove' : 'add')"
                >
                  {{ hasRole(item.key) ? '移除' : '添加' }}
                </van-button>
              </div>
            </div>
          </van-cell>
        </van-list>
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

const title = ref(process.env.VUE_APP_TITLE)
const router = useRouter()
const isLogin = ref(false)
const apps = ref([])
const loading = ref(false)
const roles = ref([])
const slogan = ref('你有多大想象力，AI就有多大创造力！')

// 只显示前5个应用
const displayApps = computed(() => {
  return apps.value.slice(0, 8)
})

onMounted(() => {
  getSystemInfo()
    .then((res) => {
      title.value = res.data.title
      if (res.data.slogan) {
        slogan.value = res.data.slogan
      }
    })
    .catch((e) => {
      ElMessage.error('获取系统配置失败：' + e.message)
    })

  checkSession()
    .then((user) => {
      isLogin.value = true
      roles.value = user.chat_roles
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
        items[i].intro = substr(items[i].hello_msg, 80)
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

  const title = ref('')
  if (opt === 'add') {
    title.value = '添加应用'
    const exists = arrayContains(roles.value, row.key)
    if (exists) {
      return
    }
    roles.value.push(row.key)
  } else {
    title.value = '移除应用'
    const exists = arrayContains(roles.value, row.key)
    if (!exists) {
      return
    }
    roles.value = removeArrayItem(roles.value, row.key)
  }
  httpPost('/api/app/update', { keys: roles.value })
    .then(() => {
      ElMessage.success({ message: title.value + '成功！', duration: 1000 })
    })
    .catch((e) => {
      ElMessage.error(title.value + '失败：' + e.message)
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

<style scoped lang="stylus">
.index {
  color var(--van-text-color)
  background-color var(--van-background)
  min-height 100vh
  display flex
  flex-direction column

  .header {
    flex-shrink 0
    padding 10px 15px
    text-align center
    background var(--van-background)
    position sticky
    top 0
    z-index 1

    .title {
      font-size 24px
      font-weight 600
      color var(--van-text-color)
    }

    .slogan {
      font-size 14px
      color var(--van-gray-6)
    }
  }

  .content {
    flex 1
    overflow-y auto
    padding 15px
    -webkit-overflow-scrolling touch

    .feature-grid {
      margin-bottom 30px

      .feature-item {
        padding 15px 0

        .feature-icon {
          width 50px
          height 50px
          border-radius 50%
          background var(--van-primary-color)
          display flex
          align-items center
          justify-content center
          margin-bottom 10px

          i, .van-icon {
            font-size 24px
            color white
          }
        }

        .text {
          font-size 14px
          font-weight 500
        }
      }
    }

    .section-header {
      display flex
      justify-content space-between
      align-items center
      margin-bottom 15px

      .section-title {
        font-size 18px
        font-weight 600
        color var(--van-text-color)
      }

      .more-btn {
        padding 0 10px
        font-size 12px
        border-radius 15px
      }
    }

    .app-list {
      .app-cell {
        padding 0
        margin-bottom 15px

        .app-card {
          background var(--van-cell-background)
          border-radius 12px
          padding 15px
          box-shadow 0 2px 12px rgba(0, 0, 0, 0.05)

          .app-info {
            display flex
            align-items center
            margin-bottom 15px

            .app-image {
              width 60px
              height 60px
              margin-right 15px

              :deep(.van-image) {
                width 100%
                height 100%
              }
            }

            .app-detail {
              flex 1

              .app-title {
                font-size 16px
                font-weight 600
                margin-bottom 5px
                color var(--van-text-color)
              }

              .app-desc {
                font-size 13px
                color var(--van-gray-6)
                display -webkit-box
                -webkit-box-orient vertical
                -webkit-line-clamp 2
                overflow hidden
              }
            }
          }

          .app-actions {
            display flex
            gap 10px

            .action-btn {
              flex 1
              border-radius 20px
              padding 0 10px
            }
          }
        }
      }
    }
  }
}
</style>
