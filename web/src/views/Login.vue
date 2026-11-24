<template>
  <div class="login-page">
    <header class="login-header">
      <a href="/" class="login-brand" aria-label="ThinkBuddy 首页">
        <img :src="'/assets/logo.png'" width="32" height="32" alt="" />
        <span>ThinkBuddy</span>
      </a>
    </header>

    <main class="login-main">
      <section class="login-card" aria-labelledby="login-title">
        <a href="/" class="login-back">
          <el-icon><ArrowLeft /></el-icon>
          <span>返回首页</span>
        </a>
        <div class="login-intro">
          <h1 id="login-title" class="login-title">{{ title }}</h1>
          <p v-if="isRegister">{{ subtitle }}</p>
        </div>

        <div class="register-content">
          <login-dialog
            :show="true"
            :active="active"
            :inviteCode="inviteCode"
            @success="handleRegisterSuccess"
            @changeActive="handleChangeActive"
          />
        </div>
      </section>
    </main>

    <footer class="login-footer">
      © {{ currentYear }} ThinkBuddy. Powered by ThinkBuddy (hanyoai).
    </footer>
  </div>
</template>

<script setup>
import LoginDialog from '@/components/LoginDialog.vue'
import { ArrowLeft } from '@element-plus/icons-vue'
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { setUserToken } from '@/store/session'
import { useSharedStore } from '@/store/sharedata'
import { httpPost } from '@/utils/http'
import { ElMessage } from 'element-plus'
import { safeRedirect } from '@/utils/redirect'

const router = useRouter()
const store = useSharedStore()
const currentYear = new Date().getFullYear()
const inviteCode = router.currentRoute.value.query.invite_code || router.currentRoute.value.query.invite || ''
const isRegister = ref(router.currentRoute.value.path === '/register')
const active = isRegister.value ? 'register' : 'login'
const title = computed(() => (isRegister.value ? '用户注册' : '用户登录'))
const subtitle = computed(() =>
  isRegister.value ? '创建您的账户以开始使用服务' : '登录您的账户以继续使用服务'
)

// 处理注册成功
const handleRegisterSuccess = () => {
  window.location.assign(safeRedirect(router.currentRoute.value.query.redirect))
}

const handleChangeActive = (newValue) => {
  isRegister.value = !newValue
}

onMounted(async () => {
  const { tokengine, message, token, ...query } = router.currentRoute.value.query
  if (!tokengine) {
    if (token) {
      setUserToken(token)
      handleRegisterSuccess()
    }
    return
  }

  await router.replace({ path: '/login', query })
  if (tokengine !== 'success') {
    ElMessage.error(typeof message === 'string' ? message : 'Tokengine 登录失败，请重试')
    return
  }

  try {
    const res = await httpPost('/api/user/login/tokengine/complete')
    setUserToken(res.data.token)
    store.setIsLogin(true)
    ElMessage.success('登录成功！')
    handleRegisterSuccess()
  } catch (error) {
    ElMessage.error(error.message || 'Tokengine 登录失败，请重试')
  }
})
</script>

<style scoped>
.login-page,
:global(.login-policy-dialog),
:global(.reset-pass-dialog) {
  --el-color-primary: var(--login-link-color);
  --el-color-primary-light-3: var(--login-link-color);
  --el-color-primary-light-9: var(--login-link-hover-bg);
  --el-color-primary-dark-2: var(--login-primary-hover);
  --el-text-color-primary: var(--login-title-color);
  --el-text-color-regular: var(--login-text-color);
  --el-text-color-secondary: var(--login-subtitle-color);
  --el-text-color-placeholder: var(--login-placeholder-color);
  --el-border-color: var(--login-card-border);
  --el-border-color-light: var(--login-card-border);
  --el-fill-color-blank: var(--login-input-bg);
  --el-fill-color-light: var(--login-soft-bg);
  --el-bg-color: var(--login-card-bg);
  --el-bg-color-overlay: var(--login-card-bg);
  --el-box-shadow: 0 1px 3px rgba(26, 24, 21, 0.08);
  font-family: 'The Future', Arial, 'PingFang SC', 'Microsoft YaHei',
    'Noto Sans SC', system-ui, sans-serif;
  color: var(--login-title-color);
}

.login-page {
  min-height: 100dvh;
  display: flex;
  flex-direction: column;
  background-color: #fff;
  background-image: none;
  background-repeat: no-repeat;
}

.login-header {
  width: 100%;
  max-width: 1280px;
  height: 64px;
  margin: 0 auto;
  padding: 0 32px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-shrink: 0;
}

.login-brand,
.login-back {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  color: inherit;
  text-decoration: none;
}

.login-brand {
  font-size: 18px;
  font-weight: 600;
  letter-spacing: -0.025em;
}

:global(:root[data-theme='dark'] .login-brand img) {
  filter: invert(1);
}

.login-back {
  margin-bottom: 12px;
  color: var(--login-text-color);
  font-size: 13px;
}

.login-back:hover {
  color: var(--login-link-color);
}

.login-brand:focus-visible,
.login-back:focus-visible {
  outline: 2px solid var(--login-link-color);
  outline-offset: 4px;
  border-radius: 4px;
}

.login-main {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 32px 16px;
}

.login-card {
  width: 100%;
  max-width: 448px;
  margin-bottom: 32px;
  padding: 36px 40px;
  border: 1px solid var(--login-card-border);
  border-radius: 12px;
  background: var(--login-card-bg);
}

.login-intro {
  margin-bottom: 32px;
  text-align: center;
}

.login-page .login-title {
  margin: 0 0 12px;
  font-family: 'The Future', Arial, 'PingFang SC', 'Microsoft YaHei',
    'Noto Sans SC', sans-serif;
  font-size: 36px;
  font-weight: 400;
  letter-spacing: -0.03em;
  line-height: 1.2;
}

.login-page .login-intro p {
  margin: 0;
  font-family: inherit;
  font-size: 14px;
  line-height: 1.7;
  color: var(--login-subtitle-color);
}

.login-footer {
  padding: 16px 16px 24px;
  color: var(--login-subtitle-color);
  font-size: 12px;
  line-height: 1.8;
  text-align: center;
}

:global(.login-policy-dialog),
:global(.reset-pass-dialog) {
  max-width: calc(100vw - 32px);
  border: 1px solid var(--login-card-border);
  border-radius: 12px;
}

:global(.login-policy-dialog) {
  width: 640px;
}

@media (max-width: 480px) {
  .login-header {
    padding: 0 20px;
  }

  .login-main {
    padding: 20px 16px;
  }

  .login-card {
    padding: 28px 24px;
  }

  .login-page .login-title {
    font-size: 32px;
  }

}
</style>
