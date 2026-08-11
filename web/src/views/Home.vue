<template>
  <div class="layout">
    <div class="tab-box">
      <div class="menu-list pt-2">
        <ul>
          <li
            class="menu-list-item flex-center-col"
            v-for="item in mainNavs"
            :key="item.url"
            @click="changeNav(item)"
            :class="curPath.startsWith(item.url) ? 'active' : ''"
          >
            <span v-if="item.icon.startsWith('icon')">
              <i class="iconfont" :class="item.icon"></i>
            </span>
            <el-image :src="item.icon" class="el-icon ml-1" v-else />
            <div class="menu-title">
              {{ item.name }}
            </div>
          </li>
        </ul>

        <!-- 更多 -->
        <div class="bot p-2">
          <div class="bot-line"></div>
          <el-popover v-if="moreNavs.length > 0" placement="right-end" trigger="hover">
            <template #reference>
              <li class="menu-list-item flex-center-col">
                <i class="iconfont icon-more" />
              </li>
            </template>
            <template #default>
              <ul class="more-menus">
                <li
                  v-for="(item, index) in moreNavs"
                  :key="item.url"
                  :class="{
                    active: curPath.startsWith(item.url),
                    moreTitle: index !== 3 && index !== 4,
                    twoTittle: index === 3 || index === 4,
                  }"
                >
                  <a @click="changeNav(item)">
                    <span v-if="item.icon.startsWith('icon')" class="mr-2">
                      <i class="iconfont" :class="item.icon"></i>
                    </span>
                    <el-image :src="item.icon" style="width: 20px; height: 20px" v-else />
                    <span class="title" :class="curPath.startsWith(item.url) ? 'active' : ''">{{
                      item.name
                    }}</span>
                  </a>
                </li>
              </ul>
            </template>
          </el-popover>

          <el-popover placement="right-end" trigger="hover" v-if="loginUser.id">
            <template #reference>
              <li class="menu-list-item flex-center-col">
                <el-avatar
                  v-if="loginUser.avatar"
                  :src="loginUser.avatar"
                  shape="circle"
                  :size="32"
                  class="user-avatar"
                />
                <i v-else class="iconfont icon-user-circle" />
              </li>
            </template>
            <template #default>
              <ul class="more-menus setting-menus">
                <li>
                  <div @click="showPowerLogDialog = true" class="flex">
                    <i class="iconfont icon-list"></i>
                    <span class="title">算力日志</span>
                  </div>
                </li>
                <li>
                  <div @click="showMemberDialog = true" class="flex">
                    <i class="iconfont icon-config"></i>
                    <span class="title">用户设置</span>
                  </div>
                </li>
                <li>
                  <div @click="showInvitationDialog = true" class="flex">
                    <i class="iconfont icon-share"></i>
                    <span class="title">推广计划</span>
                  </div>
                </li>
                <li>
                    <a :href="githubURL" target="_blank" class="flex">
                      <i class="iconfont icon-github"></i>
                      <span class="title">项目源码</span>
                    </a>
                </li>
                <li>
                    <a href="https://docs.geekai.me" target="_blank" class="flex">
                      <i class="iconfont icon-book"></i>
                      <span class="title">项目文档</span>
                    </a>
                </li>
                <li>
                  <a @click="logout" class="flex">
                    <i class="iconfont icon-logout"></i>
                    <span class="title">退出登录</span>
                  </a>
                </li>
              </ul>
            </template>
          </el-popover>
          <div v-else class="mb-2 flex justify-center">
            <el-button @click="store.setShowLoginDialog(true)" type="primary" size="small">
              登录
            </el-button>
          </div>
          <div class="menu-bot-item">
            <a @click="router.push('/')" class="link-button">
              <i class="iconfont icon-house"></i>
            </a>
            <div class="pl-1">
              <ThemeChange size="small" />
            </div>
          </div>
        </div>
      </div>
    </div>
    <el-scrollbar class="right-main">
      <div class="content custom-scroll">
        <router-view :key="routerViewKey" v-slot="{ Component }">
          <transition name="move" mode="out-in">
            <component :is="Component"></component>
          </transition>
        </router-view>
      </div>
    </el-scrollbar>
    <!-- 算力日志弹窗 -->
    <el-dialog
      v-model="showPowerLogDialog"
      title="算力日志"
      width="90%"
      :close-on-click-modal="true"
      :close-on-press-escape="true"
      style="max-width: 1200px"
      @close="showPowerLogDialog = false"
    >
      <div class="powerlog-dialog-content">
        <PowerLog />
      </div>
    </el-dialog>

    <!-- 用户设置弹窗 -->
    <el-dialog
      v-model="showMemberDialog"
      title="用户设置"
      width="90%"
      :close-on-click-modal="true"
      :close-on-press-escape="true"
      style="max-width: 1400px"
      @close="showMemberDialog = false"
    >
      <div class="member-dialog-content">
        <Member />
      </div>
    </el-dialog>

    <!-- 推广计划弹窗 -->
    <el-dialog
      v-model="showInvitationDialog"
      title="推广计划"
      width="90%"
      :close-on-click-modal="true"
      :close-on-press-escape="true"
      style="max-width: 1200px"
      @close="showInvitationDialog = false"
    >
      <div class="invitation-dialog-content">
        <Invitation />
      </div>
    </el-dialog>

    <el-dialog v-model="showLoginDialog" width="500px" @close="store.setShowLoginDialog(false)">
      <template #header>
        <div class="text-center text-xl" style="color: var(--theme-text-color-primary)">
          登录后解锁功能
        </div>
      </template>
      <div class="w-full p-4 pt-2 pb-2">
        <LoginDialog @success="loginSuccess" @hide="store.setShowLoginDialog(false)" />
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import LoginDialog from '@/components/LoginDialog.vue'
import ThemeChange from '@/components/ThemeChange.vue'
import PowerLog from '@/views/PowerLog.vue'
import Member from '@/views/Member.vue'
import Invitation from '@/views/Invitation.vue'
import { checkSession, getSystemInfo } from '@/store/cache'
import { removeUserToken } from '@/store/session'
import { useSharedStore } from '@/store/sharedata'
import { httpGet } from '@/utils/http'
import { ElMessage } from 'element-plus'
import { computed, onMounted, ref, watch } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const logo = ref('')
const mainNavs = ref([])
const moreNavs = ref([])
const curPath = ref()

const title = ref('')
const store = useSharedStore()
const loginUser = ref({})
const routerViewKey = ref(0)
const showPowerLogDialog = ref(false)
const showMemberDialog = ref(false)
const showInvitationDialog = ref(false)
const showLoginDialog = ref(false)
const githubURL = ref(import.meta.env.VITE_GITHUB_URL)

/**
 * 从路径名中提取第一个路径段
 * @param pathname - URL 的路径名部分，例如 '/chat/12345'
 * @returns 第一个路径段（不含斜杠），例如 'chat'，如果不存在则返回 null
 */
const extractFirstSegment = (pathname) => {
  const segments = pathname.split('/').filter((segment) => segment.length > 0)
  return segments.length > 0 ? segments[0] : null
}
const getFirstPathSegment = (url) => {
  try {
    // 尝试使用 URL 构造函数解析完整的 URL
    const parsedUrl = new URL(url)
    return extractFirstSegment(parsedUrl.pathname)
  } catch (error) {
    // 如果解析失败，假设是相对路径，使用当前窗口的位置作为基准
    if (typeof window !== 'undefined') {
      const parsedUrl = new URL(url, window.location.origin)
      return extractFirstSegment(parsedUrl.pathname)
    }
    // 如果无法解析，返回 null
    return null
  }
}

const stars = computed(() => {
  return 1000
})

watch(
  () => store.showLoginDialog,
  (newValue) => {
    showLoginDialog.value = newValue
  }
)

// 监听路由变化;
router.beforeEach((to, from, next) => {
  curPath.value = to.path
  next()
})

if (curPath.value === '/external') {
  curPath.value = router.currentRoute.value.query.url
}
const changeNav = (item) => {
  curPath.value = item.url
  if (item.url.indexOf('http') !== -1) {
    // 外部链接
    router.push({ path: '/external', query: { url: item.url, title: item.name } })
  } else {
    // 路由切换，确保路径变化
    if (router.currentRoute.value.path !== item.url) {
      router.push(item.url).then(() => {
        // 刷新 `routerViewKey` 触发视图重新渲染
        routerViewKey.value += 1
      })
    }
  }
}

onMounted(() => {
  curPath.value = router.currentRoute.value.path
  getSystemInfo()
    .then((res) => {
      logo.value = res.data.logo
      title.value = res.data.title
    })
    .catch((e) => {
      ElMessage.error('获取系统配置失败：' + e.message)
    })
  // 获取菜单
  httpGet('/api/menu/list')
    .then((res) => {
      mainNavs.value = res.data
      // 根据窗口的高度计算应该显示多少菜单
      const rows = Math.floor((window.innerHeight - 100) / 75)
      if (res.data.length > rows) {
        mainNavs.value = res.data.slice(0, rows)
        moreNavs.value = res.data.slice(rows)
      }
    })
    .catch((e) => {
      ElMessage.error('获取系统菜单失败：' + e.message)
    })

  curPath.value = '/' + getFirstPathSegment(window.location.href)
  init()
})

const init = () => {
  checkSession()
    .then((user) => {
      loginUser.value = user
    })
    .catch(() => {})
}

const logout = function () {
  httpGet('/api/user/logout')
    .then(() => {
      removeUserToken()
      router.push('/login')
    })
    .catch(() => {
      ElMessage.error('注销失败！')
    })
}

const loginSuccess = () => {
  init()
  store.setShowLoginDialog(false)
  // 刷新组件
  routerViewKey.value += 1
}
</script>

<style lang="scss" scoped>
@use '../assets/css/custom-scroll.scss' as *;
@use '../assets/css/home.scss' as *;
</style>

<style lang="scss">
.powerlog-dialog-content,
.member-dialog-content,
.invitation-dialog-content {
  max-height: calc(100vh - 150px);
  overflow-y: auto;
  padding: 0;
}

.powerlog-dialog-content {
  .power-log {
    .inner {
      padding: 0;
    }
  }
}

.member-dialog-content {
  .member-page {
    min-height: auto;
    padding: 0;
  }
}

.invitation-dialog-content {
  .page-invitation {
    .inner {
      padding: 20px;
    }
  }
}
</style>
