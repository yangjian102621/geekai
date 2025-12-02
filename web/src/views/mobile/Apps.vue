<template>
  <div class="apps-page">
    <van-nav-bar title="全部应用" left-arrow @click-left="router.back()" />

    <div class="apps-filter mb-8 pt-8" style="border: 1px solid #ccc">
      <van-tabs v-model="activeTab" animated swipeable>
        <van-tab title="全部分类">
          <div class="app-list">
            <van-list v-model="loading" :finished="true" finished-text="" @load="fetchApps()">
              <van-cell v-for="item in apps" :key="item.id" class="app-cell">
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
                    <van-button
                      size="small"
                      type="primary"
                      class="action-btn"
                      @click="useRole(item.id)"
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
        </van-tab>
        <van-tab v-for="type in appTypes" :key="type.id" :title="type.name">
          <div class="app-list">
            <van-list
              v-model="loading"
              :finished="true"
              finished-text=""
              @load="fetchApps(type.id)"
            >
              <van-cell v-for="item in typeApps" :key="item.id" class="app-cell">
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
                    <van-button
                      size="small"
                      type="primary"
                      class="action-btn"
                      @click="useRole(item.id)"
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
        </van-tab>
      </van-tabs>
    </div>
  </div>
</template>

<script setup>
import { checkSession } from '@/store/cache'
import { httpGet, httpPost } from '@/utils/http'
import { arrayContains, removeArrayItem, showLoginDialog, substr } from '@/utils/libs'
import { showNotify } from 'vant'
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const isLogin = ref(false)
const apps = ref([])
const typeApps = ref([])
const appTypes = ref([])
const loading = ref(false)
const roles = ref([])
const activeTab = ref(0)

onMounted(() => {
  checkSession()
    .then((user) => {
      isLogin.value = true
      roles.value = user.chat_roles
    })
    .catch(() => {})
  fetchAppTypes()
  fetchApps()
})

const fetchAppTypes = () => {
  httpGet('/api/app/type/list')
    .then((res) => {
      appTypes.value = res.data
    })
    .catch((e) => {
      showNotify({ type: 'danger', message: '获取应用分类失败：' + e.message })
    })
}

const fetchApps = (typeId = '') => {
  httpGet('/api/app/list', { tid: typeId })
    .then((res) => {
      const items = res.data
      // 处理 hello message
      for (let i = 0; i < items.length; i++) {
        items[i].intro = substr(items[i].hello_msg, 80)
      }

      if (typeId) {
        typeApps.value = items
      } else {
        apps.value = items
      }
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
      showNotify({ type: 'success', message: title.value + '成功！' })
    })
    .catch((e) => {
      showNotify({ type: 'danger', message: title.value + '失败：' + e.message })
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
.apps-page {
  min-height 100vh
  background-color var(--van-background)

  .apps-filter {
    padding 10px 0

    :deep(.van-tabs__nav) {
      background var(--van-background-2)
    }
  }

  .app-list {
    padding 0 15px

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
</style>
