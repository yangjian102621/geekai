<template>
  <div class="admin-login">
    <div class="main">
      <div class="contain">
        <div class="header">{{ title }}</div>
        <div class="content">
          <div class="block">
            <el-input placeholder="请输入用户名" size="large" v-model="username" autocomplete="off">
              <template #prefix>
                <el-icon>
                  <UserFilled/>
                </el-icon>
              </template>
            </el-input>
          </div>

          <div class="block">
            <el-input placeholder="请输入密码" size="large" v-model="password" show-password autocomplete="off">
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

        </div>
      </div>

      <footer class="footer">
        <footer-bar/>
      </footer>
    </div>
  </div>
</template>

<script setup>

import {onMounted, ref} from "vue";
import {Lock, UserFilled} from "@element-plus/icons-vue";
import {httpPost} from "@/utils/http";
import {ElMessage} from "element-plus";
import {useRouter} from "vue-router";
import FooterBar from "@/components/FooterBar.vue";
import {setAdminToken} from "@/store/session";
import {checkAdminSession} from "@/action/session";

const router = useRouter();
const title = ref('ChatGPT Plus Admin');
const username = ref(process.env.VUE_APP_ADMIN_USER);
const password = ref(process.env.VUE_APP_ADMIN_PASS);

checkAdminSession().then(() => {
  router.push("/admin")
}).catch(() => {
})
onMounted(() => {
  document.addEventListener('keyup', (e) => {
    if (e.key === 'Enter') {
      login();
    }
  });
})

const login = function () {
  if (username.value === '') {
    return ElMessage.error('请输入用户名');
  }
  if (password.value === '') {
    return ElMessage.error('请输入密码');
  }

  httpPost('/api/admin/login', {username: username.value.trim(), password: password.value.trim()}).then(res => {
    setAdminToken(res.data)
    router.push("/admin")
  }).catch((e) => {
    ElMessage.error('登录失败，' + e.message)
  })
}


</script>

<style lang="stylus" scoped>
.admin-login {
  display flex
  justify-content center
  width: 100%
  background #2D3A4B

  .main {
    width 400px;
    display flex
    justify-content center
    align-items center
    height 100vh

    .contain {
      width 100%
      padding 20px 40px;
      color #ffffff
      border-radius 10px;
      background rgba(255, 255, 255, 0.3)

      .header {
        width 100%
        margin-bottom 20px
        font-size 24px
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
}

</style>