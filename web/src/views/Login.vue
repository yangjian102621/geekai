<template>
  <div>
    <div class="bg"></div>
    <div class="main">
      <div class="contain">
        <div class="logo">
          <el-image :src="logo" fit="cover"/>
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

          <el-row class="text-line" gutter="20">
            <el-button type="primary" @click="router.push('/register')" size="small" plain>注册新账号</el-button>
            <el-button type="success" @click="showResetPass = true" size="small" plain>重置密码</el-button>
          </el-row>
        </div>
      </div>

      <reset-pass @hide="showResetPass = false" :show="showResetPass"/>

      <footer class="footer">
        <footer-bar/>
      </footer>
    </div>
  </div>
</template>

<script setup>

import {ref} from "vue";
import {Lock, UserFilled} from "@element-plus/icons-vue";
import {httpGet, httpPost} from "@/utils/http";
import {ElMessage} from "element-plus";
import {useRouter} from "vue-router";
import FooterBar from "@/components/FooterBar.vue";
import {isMobile} from "@/utils/libs";
import {checkSession} from "@/action/session";
import {setUserToken} from "@/store/session";
import {prevRoute} from "@/router";
import ResetPass from "@/components/ResetPass.vue";

const router = useRouter();
const title = ref('ChatPlus 用户登录');
const username = ref(process.env.VUE_APP_USER);
const password = ref(process.env.VUE_APP_PASS);
const showResetPass = ref(false)
const logo = ref("/images/logo.png")

// 获取系统配置
httpGet("/api/config/get?key=system").then(res => {
  logo.value = res.data.logo
}).catch(e => {
  ElMessage.error("获取系统配置失败：" + e.message)
})


checkSession().then(() => {
  if (isMobile()) {
    router.push('/mobile')
  } else {
    router.push('/chat')
  }
}).catch(() => {
})

const handleKeyup = (e) => {
  if (e.key === 'Enter') {
    login();
  }
};

const login = function () {
  if (username.value.trim() === '') {
    return ElMessage.error("请输入用户民")
  }
  if (password.value.trim() === '') {
    return ElMessage.error('请输入密码');
  }

  httpPost('/api/user/login', {username: username.value.trim(), password: password.value.trim()}).then((res) => {
    setUserToken(res.data)
    if (prevRoute.path === '' || prevRoute.path === '/register') {
      if (isMobile()) {
        router.push('/mobile')
      } else {
        router.push('/chat')
      }
    } else {
      router.push(prevRoute.path)
    }

  }).catch((e) => {
    ElMessage.error('登录失败，' + e.message)
  })
}


</script>

<style lang="stylus" scoped>
.bg {
  position fixed
  left 0
  right 0
  top 0
  bottom 0
  background-color #313237
  background-image url("~@/assets/img/login-bg.jpg")
  background-size cover
  background-position center
  background-repeat repeat-y
  //filter: blur(10px); /* 调整模糊程度，可以根据需要修改值 */
}

.main {
  .contain {
    position fixed
    left 50%
    top 40%
    width 90%
    max-width 400px;
    transform translate(-50%, -50%)
    padding 20px 10px;
    color #ffffff
    border-radius 10px;

    .logo {
      text-align center

      .el-image {
        width 120px;
      }
    }

    .header {
      width 100%
      margin-bottom 24px
      font-size 24px
      color $white_v1
      letter-space 2px
      text-align center
    }

    .content {
      width 100%
      height: auto
      border-radius 3px

      .block {
        margin-bottom 16px

        .el-input__inner {
          border 1px solid $gray-v6 !important

          .el-icon-user, .el-icon-lock {
            font-size 20px
          }
        }
      }

      .btn-row {
        padding-top 10px;

        .login-btn {
          width 100%
          font-size 16px
          letter-spacing 2px
        }
      }

      .text-line {
        justify-content center
        padding-top 10px;
        font-size 14px;
      }
    }
  }

  .footer {
    color #ffffff;

    .container {
      padding 20px;
    }
  }
}
</style>