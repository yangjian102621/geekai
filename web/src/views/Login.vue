<template>
  <div>
    <div class="bg"></div>
    <div class="main">
      <div class="contain">
        <div class="logo">
          <el-image :src="logo" fit="cover" @click="router.push('/')"/>
        </div>
        <div class="header">{{ title }}</div>
        <div class="content">
          <div class="block">
            <el-input placeholder="账号" size="large" v-model="username" autocomplete="off" autofocus
                      @keyup="handleKeyup">
              <template #prefix>
                <el-icon>
                  <UserFilled/>
                </el-icon>
              </template>
            </el-input>
          </div>

          <div class="block">
            <el-input placeholder="密码" size="large" v-model="password" show-password autocomplete="off"
                      @keyup="handleKeyup">
              <template #prefix>
                <el-icon>
                  <Lock/>
                </el-icon>
              </template>
            </el-input>
          </div>

          <el-row class="btn-row">
            <el-button class="login-btn" size="large" type="primary" @click="login">登录</el-button>
          </el-row>

          <el-row class="opt" :gutter="24">
            <el-col :span="8"><el-link type="primary" @click="router.push('/register')">注册</el-link></el-col>
            <el-col :span="8">
              <el-link type="info" @click="showResetPass = true">重置密码</el-link>
            </el-col>
            <el-col :span="8">
              <el-link type="info" @click="router.push('/')">首页</el-link>
            </el-col>
          </el-row>

          <div v-if="wechatLoginURL !== ''">
            <el-divider class="divider">其他登录方式</el-divider>

            <div class="clogin">
              <a class="wechat-login" :href="wechatLoginURL"><i class="iconfont icon-wechat"></i></a>
            </div>
          </div>
        </div>
      </div>

      <reset-pass @hide="showResetPass = false" :show="showResetPass"/>
      
      <footer-bar/>
    </div>
  </div>
</template>

<script setup>

import {onMounted, ref} from "vue";
import {Lock, UserFilled} from "@element-plus/icons-vue";
import {httpGet, httpPost} from "@/utils/http";
import {useRouter} from "vue-router";
import FooterBar from "@/components/FooterBar.vue";
import {isMobile} from "@/utils/libs";
import {checkSession} from "@/action/session";
import {setUserToken} from "@/store/session";
import ResetPass from "@/components/ResetPass.vue";
import {showMessageError} from "@/utils/dialog";

const router = useRouter();
const title = ref('Geek-AI');
const username = ref(process.env.VUE_APP_USER);
const password = ref(process.env.VUE_APP_PASS);
const showResetPass = ref(false)
const logo = ref("")
const licenseConfig = ref({})
const wechatLoginURL = ref('')

onMounted(() => {
  // 获取系统配置
  httpGet("/api/config/get?key=system").then(res => {
    logo.value = res.data.logo
    title.value = res.data.title
  }).catch(e => {
    showMessageError("获取系统配置失败：" + e.message)
  })

  httpGet("/api/config/license").then(res => {
    licenseConfig.value = res.data
  }).catch(e => {
    showMessageError("获取 License 配置：" + e.message)
  })

  checkSession().then(() => {
    if (isMobile()) {
      router.push('/mobile')
    } else {
      router.push('/chat')
    }
  }).catch(() => {
  })

  const returnURL = `${location.protocol}//${location.host}/login/callback`
  httpGet("/api/user/clogin?return_url="+returnURL).then(res => {
    wechatLoginURL.value = res.data.url
  }).catch(e => {
    console.error(e)
  })
})

const handleKeyup = (e) => {
  if (e.key === 'Enter') {
    login();
  }
};

const login = function () {
  if (username.value.trim() === '') {
    return showMessageError("请输入用户民")
  }
  if (password.value.trim() === '') {
    return showMessageError('请输入密码');
  }

  httpPost('/api/user/login', {username: username.value.trim(), password: password.value.trim()}).then((res) => {
    setUserToken(res.data.token)
    if (isMobile()) {
      router.push('/mobile')
    } else {
      router.push('/chat')
    }

  }).catch((e) => {
    showMessageError('登录失败，' + e.message)
  })
}

</script>

<style lang="stylus" scoped>
@import "@/assets/css/login.styl"
</style>