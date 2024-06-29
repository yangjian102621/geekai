<template>
  <div class="home">
    <div class="header">
      <div class="banner">
        <div class="logo">
          <el-image :src="logo" @click="router.push('/')"/>
        </div>
        <div class="title">
          <span>{{ title }}</span>
        </div>
      </div>

      <div class="navbar">

        <el-tooltip
            class="box-item"
            effect="light"
            content="部署文档"
            placement="bottom">
          <a href="https://ai.r9it.com/docs/install/" class="link-button" target="_blank">
            <i class="iconfont icon-book"></i>
          </a>
        </el-tooltip>

        <el-tooltip
            class="box-item"
            effect="light"
            content="项目源码"
            placement="bottom">
          <a href="https://github.com/yangjian102621/chatgpt-plus" class="link-button" target="_blank">
            <i class="iconfont icon-github"></i>
          </a>
        </el-tooltip>

        <el-dropdown :hide-on-click="true" class="user-info" trigger="click" v-if="loginUser.id">
                        <span class="el-dropdown-link">
                          <el-image :src="loginUser.avatar"/>
                        </span>
          <template #dropdown>
            <el-dropdown-menu class="user-info-menu">
              <el-dropdown-item @click="showConfigDialog = true">
                <el-icon>
                  <UserFilled/>
                </el-icon>
                <span class="username">{{ loginUser.nickname }}</span>
              </el-dropdown-item>

              <el-dropdown-item>
                <i class="iconfont icon-book"></i>
                <a href="https://github.com/yangjian102621/chatgpt-plus" target="_blank">
                  用户手册
                </a>
              </el-dropdown-item>

              <el-dropdown-item>
                <i class="iconfont icon-github"></i>
                <a href="https://ai.r9it.com/docs/" target="_blank">
                  Geek-AI {{ version }}
                </a>
              </el-dropdown-item>
              <el-divider style="margin: 2px 0"/>
              <el-dropdown-item @click="logout">
                <i class="iconfont icon-logout"></i>
                <span>退出登录</span>
              </el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>

        <div v-else>
          <el-button size="small" color="#21aa93" @click="store.setShowLoginDialog(true)" round>登录</el-button>
          <el-button size="small" @click="router.push('/register')" round>注册</el-button>
        </div>
      </div>
    </div>
    <div class="main">
      <div class="navigator">
        <ul class="nav-items">
          <li v-for="item in mainNavs" :key="item.url">
            <el-tooltip
                effect="light"
                :content="item.name"
                placement="right">
              <a @click="changeNav(item)" :class="item.url === curPath ? 'active' : ''">
                <el-image :src="item.icon" style="width: 30px;height: 30px"/>
              </a>
            </el-tooltip>
          </li>

          <el-popover
              v-if="moreNavs.length > 0"
              placement="right-end"
              trigger="hover"
          >
            <template #reference>
              <li>
                <a class="active">
                  <el-image src="/images/menu/more.png" style="width: 30px;height: 30px"/>
                </a>
              </li>
            </template>
            <template #default>
              <ul class="more-menus">
                <li v-for="item in moreNavs" :key="item.url" :class="item.url === curPath ? 'active' : ''">
                  <a @click="changeNav(item)">
                    <el-image :src="item.icon" style="width: 20px;height: 20px"/>
                    <span :class="item.url === curPath ? 'title active' : 'title'">{{ item.name }}</span>
                  </a>
                </li>
              </ul>
            </template>
          </el-popover>
        </ul>
      </div>

      <div class="content custom-scroll" :style="{height: mainWinHeight+'px'}">
        <router-view :key="routerViewKey" v-slot="{ Component }">
          <transition name="move" mode="out-in">
            <component :is="Component"></component>
          </transition>
        </router-view>
      </div>
    </div>

    <login-dialog :show="show" @hide="store.setShowLoginDialog(false)" @success="loginCallback"/>
    <config-dialog v-if="loginUser.id" :show="showConfigDialog" @hide="showConfigDialog = false"/>
  </div>
</template>

<script setup>

import {useRouter} from "vue-router";
import {onMounted, ref, watch} from "vue";
import {httpGet} from "@/utils/http";
import {ElMessage} from "element-plus";
import {UserFilled} from "@element-plus/icons-vue";
import {checkSession} from "@/action/session";
import {removeUserToken} from "@/store/session";
import LoginDialog from "@/components/LoginDialog.vue";
import {useSharedStore} from "@/store/sharedata";
import ConfigDialog from "@/components/ConfigDialog.vue";

const router = useRouter();
const logo = ref('/images/logo.png');
const mainNavs = ref([])
const moreNavs = ref([])
const curPath = ref(router.currentRoute.value.path)
const title = ref("")
const mainWinHeight = window.innerHeight - 50
const loginUser = ref({})
const version = ref(process.env.VUE_APP_VERSION)
const routerViewKey = ref(0)
const showConfigDialog = ref(false)

const store = useSharedStore();
const show = ref(false)
watch(() => store.showLoginDialog, (newValue) => {
  show.value = newValue
});

if (curPath.value === "/external") {
  curPath.value = router.currentRoute.value.query.url
}
const changeNav = (item) => {
  curPath.value = item.url
  if (item.url.indexOf("http") !== -1) { // 外部链接
    router.push({name: 'ExternalLink', query: {url: item.url}})
  } else {
    router.push(item.url)
  }
}

onMounted(() => {
  httpGet("/api/config/get?key=system").then(res => {
    logo.value = res.data.logo
    title.value = res.data.title
  }).catch(e => {
    ElMessage.error("获取系统配置失败：" + e.message)
  })
  // 获取菜单
  httpGet("/api/menu/list").then(res => {
    mainNavs.value = res.data
    // 根据窗口的高度计算应该显示多少菜单
    const rows = Math.floor((window.innerHeight - 90) / 60)
    if (res.data.length > rows) {
      mainNavs.value = res.data.slice(0, rows)
      moreNavs.value = res.data.slice(rows)
    }
  }).catch(e => {
    ElMessage.error("获取系统菜单失败：" + e.message)
  })

  init()
})

const init = () => {
  checkSession().then(user => {
    loginUser.value = user
  }).catch(() => {
  })
}

const logout = function () {
  httpGet('/api/user/logout').then(() => {
    removeUserToken()
    store.setShowLoginDialog(true)
    loginUser.value = {}
    // 刷新组件
    routerViewKey.value += 1
  }).catch(() => {
    ElMessage.error('注销失败！');
  })
}

const loginCallback = () => {
  init()
  // 刷新组件
  routerViewKey.value += 1
}
</script>

<style lang="stylus" scoped>
@import "@/assets/css/custom-scroll.styl"
@import "@/assets/css/home.styl"
</style>
