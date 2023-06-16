<template>
  <div>
    <div class="bg"></div>
    <div class="main">
      <div class="contain">
        <div class="logo">
          <el-image src="images/logo.png" fit="cover" />
        </div>

        <div class="header">{{ title }}</div>
        <div class="content">
          <el-form :model="formData" label-width="120px" ref="formRef" :rules="rules">
            <div class="block">
              <el-input placeholder="手机号/邮箱(4-30位)"
                        size="large" maxlength="30"
                        v-model="formData.username"
                        autocomplete="off">
                <template #prefix>
                  <el-icon>
                    <UserFilled/>
                  </el-icon>
                </template>
              </el-input>
            </div>

            <div class="block">
              <el-input placeholder="请输入密码(8-16位)"
                        maxlength="16" size="large"
                        v-model="formData.password" show-password
                        autocomplete="off">
                <template #prefix>
                  <el-icon>
                    <Lock/>
                  </el-icon>
                </template>
              </el-input>
            </div>

            <div class="block">
              <el-input placeholder="重复密码(8-16位)"
                        size="large" maxlength="16" v-model="formData.repass" show-password
                        autocomplete="off">
                <template #prefix>
                  <el-icon>
                    <Lock/>
                  </el-icon>
                </template>
              </el-input>
            </div>

            <el-row class="btn-row">
              <el-button class="login-btn" size="large" type="primary" @click="register">注册</el-button>
            </el-row>

            <el-row class="text-line">
              已经有账号？
              <el-link type="primary" @click="router.push('login')">登录</el-link>
            </el-row>
          </el-form>
        </div>
      </div>

      <footer class="footer">
        <footer-bar />
      </footer>
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

const router = useRouter();
const title = ref('ChatGPT-PLUS 用户注册');
const formData = ref({
  username: '',
  password: '',
  repass: '',
})
const formRef = ref(null)

const register = function () {
  if (formData.value.username.length < 4) {
    return ElMessage.error('用户名的长度为4-30个字符');
  }
  if (!validateEmail(formData.value.username) && !validateMobile(formData.value.username)) {
    return ElMessage.error('用户名不合法，请输入手机号码或者邮箱地址');
  }
  if (formData.value.password.length < 8) {
    return ElMessage.error('密码的长度为8-16个字符');
  }
  if (formData.value.repass !== formData.value.password) {
    return ElMessage.error('两次输入密码不一致');
  }

  httpPost('/api/user/register', formData.value).then(() => {
    ElMessage.success({"message": "注册成功，即将跳转到登录页...", onClose: () => router.push("login")})
  }).catch((e) => {
    ElMessage.error('注册失败，' + e.message)
  })
}

const validateEmail = function (email) {
  const regex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  return regex.test(email);
}
const validateMobile = function (mobile) {
  const regex = /^1[345789]\d{9}$/;
  return regex.test(mobile);
}

</script>

<style lang="stylus" scoped>
.bg {
  position fixed
  left 0
  right 0
  top 0
  bottom 0
  background-color #091519
  background-image url("~@/assets/img/reg-bg.jpg")
  background-size cover
  background-position center
  background-repeat no-repeat
  //filter: blur(10px); /* 调整模糊程度，可以根据需要修改值 */
}

.main {
  .contain {
    position fixed
    left 50%
    top 40%
    width 90%
    max-width 400px
    transform translate(-50%, -50%)
    padding 20px;
    color #ffffff
    border-radius 10px;
    background rgba(255, 255, 255, 0.3)

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