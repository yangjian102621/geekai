<template>
  <div
    class="min-h-screen flex items-center justify-center p-5 relative overflow-auto"
    style="background: var(--login-bg)"
  >
    <router-link
      to="/"
      class="fixed top-5 left-5 z-50 flex items-center justify-center w-11 h-11 border border-transparent rounded-xl text-white no-underline shadow-lg backdrop-blur-sm transition-all duration-300 hover:-translate-y-0.5 hover:shadow-xl"
      style="background: var(--btnColor)"
      title="返回首页"
    >
      <i class="iconfont icon-home text-xl"></i>
    </router-link>
    <div class="w-full max-w-md mx-auto">
      <div
        class="rounded-3xl p-10 shadow-2xl backdrop-blur-sm relative overflow-hidden"
        style="background: var(--login-card-bg); border: 1px solid var(--login-card-border)"
      >
        <div class="absolute top-0 left-0 right-0 h-1" style="background: var(--btnColor)"></div>
        <div class="text-center mb-8">
          <h1
            class="text-3xl font-semibold m-0 mb-2 tracking-tight"
            style="color: var(--login-title-color)"
          >
            用户注册
          </h1>
          <p class="text-base m-0 leading-relaxed" style="color: var(--login-subtitle-color)">
            创建您的账户以开始使用服务
          </p>
        </div>

        <div class="register-content">
          <login-dialog
            :show="true"
            active="register"
            :inviteCode="inviteCode"
            @success="handleRegisterSuccess"
            ref="loginDialogRef"
          />
        </div>
      </div>

      <footer-bar />
    </div>
  </div>
</template>

<script setup>
import FooterBar from '@/components/FooterBar.vue'
import LoginDialog from '@/components/LoginDialog.vue'
import { isMobile } from '@/utils/libs'
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const loginDialogRef = ref(null)
const inviteCode = ref(router.currentRoute.value.query.invite_code || '')

// 处理注册成功
const handleRegisterSuccess = () => {
  if (isMobile()) {
    router.push('/mobile')
  } else {
    router.push('/chat')
  }
}

onMounted(() => {
  // 确保默认显示注册状态
  if (loginDialogRef.value) {
    loginDialogRef.value.login = false
  }
})
</script>

<style scoped>
/* 移动端适配 */
@media (max-width: 768px) {
  .min-h-screen {
    padding: 1rem;
  }

  .fixed.top-5.left-5 {
    top: 1rem;
    left: 1rem;
    width: 2.5rem;
    height: 2.5rem;
  }

  .fixed.top-5.left-5 .iconfont {
    font-size: 1.125rem;
  }

  .max-w-md {
    margin-top: 3.75rem;
  }

  .p-10 {
    padding: 2rem 1.5rem;
  }

  .rounded-3xl {
    border-radius: 1rem;
  }

  .text-3xl {
    font-size: 1.5rem;
  }

  .text-base {
    font-size: 0.875rem;
  }
}

/* 小屏幕手机适配 */
@media (max-width: 480px) {
  .p-10 {
    padding: 1.5rem 1.25rem;
  }

  .text-3xl {
    font-size: 1.25rem;
  }

  .text-base {
    font-size: 0.875rem;
  }
}
</style>
