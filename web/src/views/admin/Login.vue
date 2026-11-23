<template>
  <div>
    <div class="bg"></div>
    <div class="admin-login">
      <div class="main">
        <div class="contain">
          <a class="logo" href="/">
            <el-image :src="logo" fit="cover" />
          </a>

          <h1 class="header">登录 {{ title }}</h1>
          <div class="content">
            <el-input
              v-model="username"
              placeholder="请输入用户名"
              size="large"
              autocomplete="off"
              autofocus
              @keyup.enter="login"
            >
              <template #prefix>
                <el-icon>
                  <UserFilled />
                </el-icon>
              </template>
            </el-input>

            <el-input
              v-model="password"
              placeholder="请输入密码"
              size="large"
              show-password
              autocomplete="off"
              @keyup.enter="login"
            >
              <template #prefix>
                <el-icon>
                  <Lock />
                </el-icon>
              </template>
            </el-input>

            <el-row class="btn-row">
              <el-button
                class="login-btn"
                size="large"
                type="primary"
                @click="login"
                :loading="loading"
                >{{ loading ? '登录中...' : '登录' }}</el-button
              >
            </el-row>
          </div>
        </div>

        <captcha v-if="enableVerify" @success="doLogin" ref="captchaRef" />
        <footer-bar class="footer" />
      </div>
    </div>
  </div>
</template>

<script setup>
import Captcha from '@/components/Captcha.vue'
import FooterBar from '@/components/FooterBar.vue'
import { checkAdminSession, getSystemInfo } from '@/store/cache'
import { setAdminToken } from '@/store/session'
import { httpPost } from '@/utils/http'
import { Lock, UserFilled } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const title = ref('Geek-AI Console')
const username = ref(import.meta.env.VITE_ADMIN_USER)
const password = ref(import.meta.env.VITE_ADMIN_PASS)
const logo = ref('')
const enableVerify = ref(false)
const captchaRef = ref(null)
const loading = ref(false)

onMounted(() => {
  // 判断是否登录
  checkAdminSession()
    .then(() => {
      router.replace({ name: 'admin-dashboard' })
    })
    .catch(() => {})

  // 加载系统配置
  getSystemInfo()
    .then((res) => {
      title.value = res.data.admin_title
      logo.value = res.data.logo
    })
    .catch((e) => {
      ElMessage.error('加载系统配置失败: ' + e.message)
    })
})

const login = function () {
  if (username.value === '') {
    return ElMessage.error('请输入用户名')
  }
  if (password.value === '') {
    return ElMessage.error('请输入密码')
  }
  if (enableVerify.value) {
    captchaRef.value.loadCaptcha()
  } else {
    doLogin({})
  }
}

const doLogin = function (verifyData) {
  loading.value = true
  httpPost('/api/admin/login', {
    username: username.value.trim(),
    password: password.value.trim(),
    key: verifyData.key,
    dots: verifyData.dots,
    x: verifyData.x,
  })
    .then((res) => {
      setAdminToken(res.data.token)
      router.replace({ name: 'admin-dashboard' })
    })
    .catch((e) => {
      ElMessage.error('登录失败，' + e.message)
    })
    .finally(() => {
      loading.value = false
    })
}
</script>

<style lang="scss" scoped>
.bg {
  // Together AI：浅色画布 + 洋红/橙渐变插画背景（canvas 兜底，品牌色仅用于渐变）
  position: absolute;
  left: 0;
  right: 0;
  top: 0;
  bottom: 0;
  background: var(--ta-canvas, #ffffff);
  background-image: linear-gradient(
    135deg,
    rgba(239, 44, 193, 0.08) 0%,
    rgba(252, 76, 2, 0.05) 50%,
    rgba(189, 187, 255, 0.1) 100%
  );
  z-index: 0;
}

.admin-login {
  position: absolute;
  left: 0;
  top: 0;
  z-index: 10;
  display: flex;
  justify-content: center;
  width: 100%;
  height: 100vh;

  .main {
    max-width: 400px;
    display: flex;
    justify-content: center;
    align-items: center;
    height: 100vh;

    .contain {
      width: 100%;
      padding: 40px;
      color: var(--ta-ink, #000000);
      // Together AI contained card：hairline 边框 + 8px 圆角 + 深蓝着色阴影
      border: 1px solid var(--ta-hairline-light, rgba(0, 0, 0, 0.08));
      border-radius: 8px;
      background: var(--ta-elevated-light, #ffffff);
      box-shadow: var(--ta-shadow, rgba(1, 1, 32, 0.1) 0px 4px 10px);

      .logo {
        text-align: center;

        .el-image {
          width: 96px;
          cursor: pointer;
          border-radius: 4px;
        }
      }

      .header {
        width: 100%;
        padding: 10px 0 16px;
        font-size: 26px;
        font-weight: 500;
        letter-spacing: -0.58px;
        text-align: center;
        color: var(--ta-ink, #000000);
      }

      .content {
        width: 100%;
        height: auto;
        border-radius: 4px;

        .el-input {
          margin: 10px 0;
        }

        .block {
          margin-bottom: 16px;

          .el-input__inner {
            border: 1px solid var(--ta-hairline-light, rgba(0, 0, 0, 0.08)) !important;
            border-radius: 4px !important;

            .el-icon-user,
            .el-icon-lock {
              font-size: 20px;
            }
          }
        }

        .btn-row {
          padding-top: 10px;

          .login-btn {
            width: 100%;
            font-size: 16px;
            font-weight: 500;
            letter-spacing: -0.32px;
          }
        }

        .text-line {
          justify-content: center;
          padding-top: 10px;
          font-size: 14px;
        }
      }
    }

    .foot-container {
      background: transparent;
      --text-color: rgba(0, 0, 0, 0.6);
    }
  }
}
</style>
