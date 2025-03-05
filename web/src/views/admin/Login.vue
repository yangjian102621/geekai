<template>
  <div>
    <div class="bg"></div>
    <div class="admin-login">
      <div class="main">
        <div class="contain">
          <div class="logo" @click="router.push('/')">
            <el-image :src="logo" fit="cover"/>
          </div>

          <h1 class="header">登录 {{ title }}</h1>
          <div class="content">
            <el-input v-model="username" placeholder="请输入用户名" size="large"
                      autocomplete="off" autofocus @keyup.enter="login">
              <template #prefix>
                <el-icon>
                  <UserFilled/>
                </el-icon>
              </template>
            </el-input>

            <el-input v-model="password" placeholder="请输入密码" size="large"
                      show-password autocomplete="off" @keyup.enter="login">
              <template #prefix>
                <el-icon>
                  <Lock/>
                </el-icon>
              </template>
            </el-input>

            <el-row class="btn-row">
              <el-button class="login-btn" size="large" type="primary" @click="login">登录</el-button>
            </el-row>
          </div>
        </div>

        <captcha v-if="enableVerify" @success="doLogin" ref="captchaRef"/>
        <footer-bar class="footer"/>
      </div>
    </div>
  </div>
</template>

<script setup>

import {ref} from "vue";
import {Lock, UserFilled} from "@element-plus/icons-vue";
import {httpPost} from "@/utils/http";
import {ElMessage} from "element-plus";
import {useRouter} from "vue-router";
import FooterBar from "@/components/FooterBar.vue";
import {setAdminToken} from "@/store/session";
import {checkAdminSession, getSystemInfo} from "@/store/cache";
import Captcha from "@/components/Captcha.vue";

const router = useRouter();
const title = ref('Geek-AI Console');
const username = ref(process.env.VUE_APP_ADMIN_USER);
const password = ref(process.env.VUE_APP_ADMIN_PASS);
const logo = ref("")
const enableVerify = ref(false)
const captchaRef = ref(null)

checkAdminSession().then(() => {
  router.push("/admin")
}).catch(() => {
})

// 加载系统配置
getSystemInfo().then(res => {
  title.value = res.data.admin_title
  logo.value = res.data.logo
  enableVerify.value = res.data['enabled_verify']
}).catch(e => {
  ElMessage.error("加载系统配置失败: " + e.message)
})

const login = function () {
  if (username.value === '') {
    return ElMessage.error('请输入用户名');
  }
  if (password.value === '') {
    return ElMessage.error('请输入密码');
  }
  if (enableVerify.value) {
    captchaRef.value.loadCaptcha()
  } else {
    doLogin({})
  }
}

const doLogin = function (verifyData) {
  httpPost('/api/admin/login', {
    username: username.value.trim(),
    password: password.value.trim(),
    key: verifyData.key,
    dots: verifyData.dots,
    x: verifyData.x
  }).then(res => {
    setAdminToken(res.data.token)
    router.push("/admin")
  }).catch((e) => {
    ElMessage.error('登录失败，' + e.message)
  })
}


</script>

<style lang="stylus" scoped>
.bg {
  position absolute
  left 0
  right 0
  top 0
  bottom 0
  background-color #091519
  background-image url("~@/assets/img/admin-login-bg.jpg")
  background-size cover
  background-position center
  background-repeat no-repeat
  filter: blur(10px); /* 调整模糊程度，可以根据需要修改值 */
  z-index 0
}

.admin-login {
  position absolute
  left 0
  top 0
  z-index 10
  display flex
  justify-content center
  width: 100%
  height: 100vh

  .main {
    max-width 400px
    display flex
    justify-content center
    align-items center
    height 100vh

    .contain {
      width 100%
      padding 40px;
      color #ffffff
      border-radius 10px;
      background rgba(0, 0, 0, 0.3)

      .logo {
        text-align center

        .el-image {
          width 120px;
          cursor pointer
          border-radius 50%
        }
      }

      .header {
        width 100%
        //margin-bottom 20px
        padding 10px
        font-size 26px
        text-align center
      }

      .content {
        width 100%
        height: auto
        border-radius 3px

        .el-input {
          margin 10px 0
        }

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
}

</style>