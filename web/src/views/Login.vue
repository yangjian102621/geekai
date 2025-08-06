<template>
  <div class="login-page">
    <router-link to="/" class="back-home-btn" title="返回首页">
      <i class="iconfont icon-home"></i>
    </router-link>
    <div class="login-container">
      <div class="login-card">
        <div class="login-header">
          <h1 class="login-title">欢迎登录</h1>
          <p class="login-subtitle">登录您的账户以继续使用服务</p>
        </div>

        <div class="login-content">
          <login-dialog
            :show="true"
            @hide="handleLoginHide"
            @success="handleLoginSuccess"
            ref="loginDialogRef"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import LoginDialog from '@/components/LoginDialog.vue'
import { getCurrentDeviceRedirectPath } from '@/utils/device'
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const loginDialogRef = ref(null)

// 处理登录弹窗隐藏
const handleLoginHide = () => {
  const redirectPath = getCurrentDeviceRedirectPath()
  router.push(redirectPath)
}

// 处理登录成功
const handleLoginSuccess = () => {
  const redirectPath = getCurrentDeviceRedirectPath()
  router.push(redirectPath)
}

onMounted(() => {
  // 确保默认显示登录状态
  if (loginDialogRef.value) {
    loginDialogRef.value.login = true
  }
})
</script>

<style lang="scss" scoped>
.login-page {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;

  .back-home-btn {
    position: absolute;
    top: 24px;
    left: 24px;
    z-index: 10;
    font-size: 22px;
    color: #fff;
    background: rgba(0, 0, 0, 0.15);
    border-radius: 50%;
    width: 40px;
    height: 40px;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: background 0.2s;
  }
  .back-home-btn:hover {
    background: rgba(0, 0, 0, 0.25);
  }
  @media (max-width: 768px) {
    .back-home-btn {
      top: 12px;
      left: 12px;
      font-size: 20px;
      width: 36px;
      height: 36px;
    }
  }
  :deep(.van-theme-dark) .back-home-btn {
    color: #fff;
    background: rgba(0, 0, 0, 0.35);
  }

  .login-container {
    width: 100%;
    max-width: 480px;

    .login-card {
      background: var(--el-bg-color);
      border-radius: 16px;
      box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1);
      overflow: hidden;

      .login-header {
        background: linear-gradient(135deg, var(--el-color-primary), #8b5cf6);
        color: white;
        padding: 40px 30px;
        text-align: center;

        .login-title {
          font-size: 28px;
          font-weight: 600;
          margin: 0 0 8px 0;
        }

        .login-subtitle {
          font-size: 16px;
          opacity: 0.9;
          margin: 0;
        }
      }

      .login-content {
        padding: 40px 30px;
      }
    }
  }
}

// 深色主题适配
:deep(.van-theme-dark) {
  .login-page {
    .login-card {
      background: var(--el-bg-color-overlay);
      box-shadow: 0 20px 40px rgba(0, 0, 0, 0.3);
    }
  }
}

// 移动端响应式设计
@media (max-width: 768px) {
  .login-page {
    padding: 16px;
    background: var(--van-background);

    .login-container {
      max-width: 100%;

      .login-card {
        border-radius: 20px;
        box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);

        .login-header {
          padding: 30px 20px;
          background: linear-gradient(135deg, var(--van-primary-color), #8b5cf6);

          .login-title {
            font-size: 24px;
          }

          .login-subtitle {
            font-size: 14px;
          }
        }

        .login-content {
          padding: 30px 20px;
        }
      }
    }
  }
}

// 小屏幕移动端优化
@media (max-width: 375px) {
  .login-page {
    padding: 12px;

    .login-container {
      .login-card {
        .login-header {
          padding: 24px 16px;

          .login-title {
            font-size: 22px;
          }

          .login-subtitle {
            font-size: 13px;
          }
        }

        .login-content {
          padding: 24px 16px;
        }
      }
    }
  }
}
</style>
