<template>
  <div class="layout">
    <div class="tab-box">
      <!-- <div class="flex-center-col pt-2 mb-2">
        <div class="flex flex-center-col">
          <div class="menuIcon" @click="store.setChatListExtend(!store.chatListExtend)">
            <el-tooltip content="隐藏对话列表" placement="right" v-if="store.chatListExtend">
              <i class="iconfont icon-colspan"></i>
            </el-tooltip>
            <el-tooltip content="展开对话列表" placement="right" v-else>
              <i class="iconfont icon-expand"></i>
            </el-tooltip>
          </div>
        </div>
      </div> -->

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
                <i class="iconfont icon-config" />
              </li>
            </template>
            <template #default>
              <ul class="more-menus setting-menus">
                <li>
                  <div @click="showConfigDialog = true" class="flex">
                    <el-icon>
                      <UserFilled />
                    </el-icon>
                    <span class="username title">账户信息</span>
                  </div>
                </li>
                <li v-if="!license.de_copy">
                  <a :href="githubURL" target="_blank" class="flex">
                    <i class="iconfont icon-github"></i>
                    <span class="title">项目源码</span>
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
      <div class="topheader" v-if="loginUser.id === undefined || !loginUser.id">
        <el-button
          @click="router.push('/login')"
          class="btn-go animate__animated animate__pulse animate__infinite"
          round
          >登录</el-button
        >
      </div>
      <div class="content custom-scroll">
        <router-view :key="routerViewKey" v-slot="{ Component }">
          <transition name="move" mode="out-in">
            <component :is="Component"></component>
          </transition>
        </router-view>
      </div>
      <!-- </div> -->
    </el-scrollbar>
    <config-dialog v-if="loginUser.id" :show="showConfigDialog" @hide="showConfigDialog = false" />

    <el-dialog v-model="showLoginDialog" width="500px" @close="store.setShowLoginDialog(false)">
      <template #header>
        <div class="text-center text-xl" style="color: var(--theme-text-color-primary)">
          登录后解锁功能
        </div>
      </template>
      <div class="p-4 pt-2 pb-2">
        <LoginDialog @success="loginSuccess" @hide="store.setShowLoginDialog(false)" />
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import LoginDialog from '@/components/LoginDialog.vue'
import ThemeChange from '@/components/ThemeChange.vue'
import ConfigDialog from '@/components/UserInfoDialog.vue'
import { checkSession, getLicenseInfo, getSystemInfo } from '@/store/cache'
import { removeUserToken } from '@/store/session'
import { useSharedStore } from '@/store/sharedata'
import { showMessageError } from '@/utils/dialog'
import { httpGet } from '@/utils/http'
import { UserFilled } from '@element-plus/icons-vue'
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
const showConfigDialog = ref(false)
const license = ref({ de_copy: true })
const showLoginDialog = ref(false)
const githubURL = ref(process.env.VUE_APP_GITHUB_URL)

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
  console.log(curPath.value)
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

  getLicenseInfo()
    .then((res) => {
      license.value = res.data
    })
    .catch((e) => {
      license.value = { de_copy: false }
      showMessageError('获取 License 配置：' + e.message)
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

<style lang="stylus" scoped>
@import "@/assets/css/custom-scroll.styl"
@import "@/assets/css/home.styl"
</style>
