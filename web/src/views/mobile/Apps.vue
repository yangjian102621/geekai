<template>
  <div class="apps-page p-3">
    <div class="apps-filter mb-8">
      <CustomTabs :model-value="activeTab" @update:model-value="handleTabChange">
        <CustomTabPane name="all" label="全部分类">
          <div class="app-list">
            <van-list v-model="loading" :finished="!loading" finished-text="" @load="() => {}">
              <template v-if="!loading && currentApps.length > 0">
                <AppCard
                  v-for="item in currentApps"
                  :key="item.id"
                  :app="item"
                  :has-role="hasRole(item.key)"
                  @use-role="useRole"
                  @update-role="updateRole"
                />
              </template>
              <template v-else-if="!loading && currentApps.length === 0">
                <EmptyState
                  type="search"
                  description="暂无应用"
                  :show-action="true"
                  action-text="刷新"
                  @action="refreshData"
                />
              </template>
              <template v-else>
                <div class="loading-state">
                  <van-loading type="spinner" size="24px">加载中...</van-loading>
                </div>
              </template>
            </van-list>
          </div>
        </CustomTabPane>
        <CustomTabPane
          v-for="type in appTypes"
          :key="type.id"
          :name="type.id.toString()"
          :label="type.name"
        >
          <div class="app-list">
            <van-list v-model="loading" :finished="!loading" finished-text="" @load="() => {}">
              <template v-if="!loading && getAppsByType(type.id).length > 0">
                <AppCard
                  v-for="item in getAppsByType(type.id)"
                  :key="item.id"
                  :app="item"
                  :has-role="hasRole(item.key)"
                  @use-role="useRole"
                  @update-role="updateRole"
                />
              </template>
              <template v-else-if="!loading && getAppsByType(type.id).length === 0">
                <EmptyState
                  type="search"
                  :description="`${type.name}分类暂无应用`"
                  :show-action="true"
                  action-text="刷新"
                  @action="refreshData"
                />
              </template>
              <template v-else>
                <div class="loading-state">
                  <van-loading type="spinner" size="24px">加载中...</van-loading>
                </div>
              </template>
            </van-list>
          </div>
        </CustomTabPane>
      </CustomTabs>
    </div>
  </div>
</template>

<script setup>
import AppCard from '@/components/mobile/AppCard.vue'
import EmptyState from '@/components/mobile/EmptyState.vue'
import CustomTabPane from '@/components/ui/CustomTabPane.vue'
import CustomTabs from '@/components/ui/CustomTabs.vue'
import { checkSession } from '@/store/cache'
import { httpGet, httpPost } from '@/utils/http'
import { arrayContains, removeArrayItem, showLoginDialog, substr } from '@/utils/libs'
import { showNotify } from 'vant'
import { computed, nextTick, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const isLogin = ref(false)
const allApps = ref([]) // 存储所有应用数据
const appTypes = ref([])
const loading = ref(false)
const roles = ref([])
const activeTab = ref('all')
const initialized = ref(false)

// 按分类分组的应用数据
const appsByType = computed(() => {
  const grouped = {}
  allApps.value.forEach((app) => {
    const tid = app.tid || 0
    if (!grouped[tid]) {
      grouped[tid] = []
    }
    grouped[tid].push(app)
  })
  return grouped
})

// 获取当前tab显示的应用列表
const currentApps = computed(() => {
  if (activeTab.value === 'all') {
    return allApps.value
  }
  return getAppsByType(parseInt(activeTab.value)) || []
})

onMounted(async () => {
  try {
    const user = await checkSession()
    isLogin.value = true
    roles.value = user.chat_roles
  } catch (error) {
    // 用户未登录，继续执行
  }

  await Promise.all([fetchAppTypes(), fetchAllApps()])

  initialized.value = true
})

const fetchAppTypes = async () => {
  try {
    const res = await httpGet('/api/app/type/list')
    appTypes.value = res.data
  } catch (e) {
    showNotify({ type: 'danger', message: '获取应用分类失败：' + e.message })
  }
}

// 一次性获取所有应用数据
const fetchAllApps = async () => {
  loading.value = true
  try {
    const res = await httpGet('/api/app/list')
    const items = res.data
    // 处理 hello message
    for (let i = 0; i < items.length; i++) {
      items[i].intro = substr(items[i].hello_msg, 80)
    }
    allApps.value = items
  } catch (e) {
    showNotify({ type: 'danger', message: '获取应用失败：' + e.message })
  } finally {
    loading.value = false
  }
}

// 刷新数据
const refreshData = async () => {
  await Promise.all([fetchAppTypes(), fetchAllApps()])
  showNotify({ type: 'success', message: '数据已刷新' })
}

// 根据分类ID获取对应的应用列表
const getAppsByType = (typeId) => {
  return appsByType.value[typeId] || []
}

// 处理tab切换
const handleTabChange = async (tabName) => {
  activeTab.value = tabName
  // 等待DOM更新完成
  await nextTick()
}

const updateRole = async (app, opt) => {
  if (!isLogin.value) {
    return showLoginDialog(router)
  }

  let actionTitle = ''
  if (opt === 'add') {
    actionTitle = '添加应用'
    const exists = arrayContains(roles.value, app.key)
    if (exists) {
      return
    }
    roles.value.push(app.key)
  } else {
    actionTitle = '移除应用'
    const exists = arrayContains(roles.value, app.key)
    if (!exists) {
      return
    }
    roles.value = removeArrayItem(roles.value, app.key)
  }

  try {
    await httpPost('/api/app/update', { keys: roles.value })
    showNotify({ type: 'success', message: actionTitle + '成功！' })
  } catch (e) {
    showNotify({ type: 'danger', message: actionTitle + '失败：' + e.message })
  }
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
.apps-page {
  min-height: 100vh;
  background-color: var(--van-background);

  .apps-filter {
    :deep(.van-tabs__nav) {
      background: var(--van-background-2);
    }
  }

  .app-list {
    padding: 0;

    .loading-state {
      padding: 40px 0;
      text-align: center;
      color: var(--van-gray-6);
    }
  }
}
</style>
