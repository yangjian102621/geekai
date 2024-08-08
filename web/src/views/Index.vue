<template>
  <div class="index-page" :style="{height: winHeight+'px'}">
    <div :class="theme.imageBg?'color-bg image-bg':'color-bg'" :style="{backgroundImage:'url('+bgStyle.backgroundImage+')', backgroundColor:bgStyle.backgroundColor}"></div>
    <div class="menu-box">
      <el-menu
          mode="horizontal"
          :ellipsis="false"
      >
        <div class="menu-item">
          <el-image :src="logo" alt="Geek-AI"/>
          <div class="title" :style="{color:theme.textColor}">{{ title }}</div>
        </div>
        <div class="menu-item">
          <span v-if="!license.de_copy">
            <a :href="docsURL" target="_blank">
            <el-button :color="theme.btnBgColor" :style="{color: theme.btnTextColor}" class="shadow" round>
              <i class="iconfont icon-book"></i>
              <span>文档</span>
            </el-button>
          </a>

          <a :href="gitURL" target="_blank">
            <el-button :color="theme.btnBgColor" :style="{color: theme.btnTextColor}" class="shadow" round>
              <i class="iconfont icon-github"></i>
              <span>源码</span>
            </el-button>
          </a>
          </span>

          <span v-if="!isLogin">
            <el-button :color="theme.btnBgColor" :style="{color: theme.btnTextColor}" @click="router.push('/login')" class="shadow" round>登录</el-button>
            <el-button :color="theme.btnBgColor" :style="{color: theme.btnTextColor}" @click="router.push('/register')" class="shadow" round>注册</el-button>
          </span>
        </div>
      </el-menu>
    </div>
    <div class="content">
      <h1 :style="{color:theme.textColor}">欢迎使用 {{ title }}</h1>
      <p :style="{color:theme.textColor}">{{ slogan }}</p>

      <div class="navs">
        <el-space wrap>
          <div v-for="item in navs" class="nav-item">
            <el-button @click="router.push(item.url)" :color="theme.btnBgColor" :style="{color: theme.btnTextColor}" class="shadow" :dark="false">
              <i :class="'iconfont '+iconMap[item.url]"></i>
              <span>{{item.name}}</span>
            </el-button>
          </div>
        </el-space>
      </div>
    </div>

    <footer-bar :text-color="theme.textColor" />
  </div>
</template>

<script setup>

import {onMounted, ref} from "vue";
import {useRouter} from "vue-router";
import FooterBar from "@/components/FooterBar.vue";
import {httpGet} from "@/utils/http";
import {ElMessage} from "element-plus";
import {isMobile} from "@/utils/libs";
import {checkSession, getLicenseInfo, getSystemInfo} from "@/store/cache";

const router = useRouter()

if (isMobile()) {
  router.push("/mobile")
}

const title = ref("")
const logo = ref("")
const slogan = ref("")
const license = ref({de_copy: true})
const winHeight = window.innerHeight - 150
const isLogin = ref(false)
const docsURL = ref(process.env.VUE_APP_DOCS_URL)
const gitURL = ref(process.env.VUE_APP_GIT_URL)
const navs  = ref([])
const btnColors = ref([
  {bgColor: "#fff143", textColor: "#50616D"},
  {bgColor: "#eaff56", textColor: "#50616D"},
  {bgColor: "#bddd22", textColor: "#50616D"},
  {bgColor: "#1bd1a5", textColor: "#50616D"},
  {bgColor: "#e0eee8", textColor: "#50616D"},
  {bgColor: "#7bcfa6", textColor: "#50616D"},
  {bgColor: "#bce672", textColor: "#50616D"},
  {bgColor: "#44cef6", textColor: "#ffffff"},
  {bgColor: "#70f3ff", textColor: "#50616D"},
  {bgColor: "#fffbf0", textColor: "#50616D"},
  {bgColor: "#d6ecf0", textColor: "#50616D"},
  {bgColor: "#88ada6", textColor: "#50616D"},
  {bgColor: "#30dff3", textColor: "#50616D"},
  {bgColor: "#d3e0f3", textColor: "#50616D"},
  {bgColor: "#e9e7ef", textColor: "#50616D"},
  {bgColor: "#eacd76", textColor: "#50616D"},
  {bgColor: "#f2be45", textColor: "#50616D"},
  {bgColor: "#549688", textColor: "#ffffff"},
  {bgColor: "#758a99", textColor: "#ffffff"},
  {bgColor: "#41555d", textColor: "#ffffff"},
  {bgColor: "#21aa93", textColor: "#ffffff"},
  {bgColor: "#0aa344", textColor: "#ffffff"},
  {bgColor: "#f05654", textColor: "#ffffff"},
  {bgColor: "#db5a6b", textColor: "#ffffff"},
  {bgColor: "#db5a6b", textColor: "#ffffff"},
  {bgColor: "#8d4bbb", textColor: "#ffffff"},
  {bgColor: "#426666", textColor: "#ffffff"},
  {bgColor: "#177cb0", textColor: "#ffffff"},
  {bgColor: "#395260", textColor: "#ffffff"},
  {bgColor: "#519a73", textColor: "#ffffff"},
  {bgColor: "#75878a", textColor: "#ffffff"},
])
const iconMap =ref(
    {
      "/chat": "icon-chat",
      "/mj": "icon-mj",
      "/sd": "icon-sd",
      "/dalle": "icon-dalle",
      "/images-wall": "icon-image",
      "/suno": "icon-suno",
      "/xmind": "icon-xmind",
      "/apps": "icon-app",
      "/member": "icon-vip-user",
      "/invite": "icon-share",
    }
)
const bgStyle = {}
const color = btnColors.value[Math.floor(Math.random() * btnColors.value.length)]
const theme = ref({bgColor: "#ffffff", btnBgColor: color.bgColor, btnTextColor: color.textColor, textColor: "#ffffff", imageBg:true})

onMounted(() => {
  getSystemInfo().then(res => {
    title.value = res.data.title
    logo.value = res.data.logo
    if (res.data.index_bg_url === 'color') {
      // 随机选取一种颜色
      theme.value.bgColor = color.bgColor
      theme.value.btnBgColor = color.bgColor
      theme.value.textColor = color.textColor
      theme.value.btnTextColor = color.textColor
      // 设置背景颜色
      bgStyle.backgroundColor = theme.value.bgColor
      bgStyle.backgroundImage = "/images/transparent-bg.png"
      theme.value.imageBg = false
    } else if (res.data.index_bg_url) {
      bgStyle.backgroundImage = res.data.index_bg_url
    } else {
      bgStyle.backgroundImage = "/images/index-bg.jpg"
    }

    slogan.value = res.data.slogan
  }).catch(e => {
    ElMessage.error("获取系统配置失败：" + e.message)
  })

  getLicenseInfo().then(res => {
    license.value = res.data
  }).catch(e => {
    license.value = {de_copy: false}
    ElMessage.error("获取 License 配置失败：" + e.message)
  })

  httpGet("/api/menu/list?index=1").then(res => {
    navs.value = res.data
  }).catch(e => {
    ElMessage.error("获取导航菜单失败：" + e.message)
  })

  checkSession().then(() => {
    isLogin.value = true
  }).catch(()=>{})
})
</script>

<style lang="stylus" scoped>
@import '@/assets/iconfont/iconfont.css'
@import "@/assets/css/index.styl"
</style>
