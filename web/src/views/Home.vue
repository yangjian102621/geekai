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
        <el-dropdown :hide-on-click="true" class="user-info" trigger="click" v-if="loginUser.id">
                        <span class="el-dropdown-link">
                          <el-image :src="loginUser.avatar"/>
                        </span>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item>
                <el-icon>
                  <UserFilled/>
                </el-icon>
                <span class="username">{{ loginUser.nickname }}</span>
              </el-dropdown-item>

              <el-dropdown-item>
                <i class="iconfont icon-book"></i>
                <span>
                    <el-link type="primary" href="https://github.com/yangjian102621/chatgpt-plus" target="_blank">
                      用户手册
                    </el-link>
                 </span>
              </el-dropdown-item>

              <el-dropdown-item>
                <i class="iconfont icon-github"></i>
                <span>
                    <el-link type="primary" href="https://ai.r9it.com/docs/" target="_blank">
                      Geek-AI {{ version }}
                    </el-link>
                 </span>
              </el-dropdown-item>

              <el-dropdown-item @click="logout">
                <i class="iconfont icon-logout"></i>
                <span>退出登录</span>
              </el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
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

      <div class="content" :style="{height: mainWinHeight+'px'}">
        <router-view v-slot="{ Component }">
          <transition name="move" mode="out-in">
            <component :is="Component"></component>
          </transition>
        </router-view>
      </div>
    </div>
  </div>
</template>

<script setup>

import {useRouter} from "vue-router";
import {onMounted, ref} from "vue";
import {httpGet} from "@/utils/http";
import {ElMessage} from "element-plus";
import {UserFilled} from "@element-plus/icons-vue";
import {checkSession} from "@/action/session";
import {removeUserToken} from "@/store/session";

const router = useRouter();
const logo = ref('/images/logo.png');
const mainNavs = ref([])
const moreNavs = ref([])
const curPath = ref(router.currentRoute.value.path)
const title = ref("")
const mainWinHeight = window.innerHeight - 50
const loginUser = ref({})
const version = ref(process.env.VUE_APP_VERSION)

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

  checkSession().then(user => {
    loginUser.value = user
  }).catch(() => {
  })
})

const logout = function () {
  httpGet('/api/user/logout').then(() => {
    removeUserToken()
    router.push("/login")
  }).catch(() => {
    ElMessage.error('注销失败！');
  })
}
</script>

<style lang="stylus" scoped>
@import '@/assets/iconfont/iconfont.css';
@import "@/assets/css/home.styl"
</style>
