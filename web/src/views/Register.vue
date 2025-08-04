<template>
  <div class="register-page">
    <div class="register-container">
      <div class="register-card">
        <div class="register-header">
          <h1 class="register-title">用户注册</h1>
          <p class="register-subtitle">创建您的账户以开始使用服务</p>
        </div>

        <div class="register-content">
          <login-dialog
            :show="true"
            @hide="handleRegisterHide"
            @success="handleRegisterSuccess"
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

// 处理注册弹窗隐藏
const handleRegisterHide = () => {
  const redirectPath = getCurrentDeviceRedirectPath()
  router.push(redirectPath)
}

// 处理注册成功
const handleRegisterSuccess = () => {
  const redirectPath = getCurrentDeviceRedirectPath()
  router.push(redirectPath)
}

onMounted(() => {
  // 确保默认显示注册状态
  if (loginDialogRef.value) {
    loginDialogRef.value.login = false
  }
})
</script>

<style lang="scss" scoped>
.register-page {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;

  .register-container {
    width: 100%;
    max-width: 480px;

    .register-card {
      background: var(--el-bg-color);
      border-radius: 16px;
      box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1);
      overflow: hidden;

      .register-header {
        background: linear-gradient(135deg, #10b981, #059669);
        color: white;
        padding: 40px 30px;
        text-align: center;

        .register-title {
          font-size: 28px;
          font-weight: 600;
          margin: 0 0 8px 0;
        }

        .register-subtitle {
          font-size: 16px;
          opacity: 0.9;
          margin: 0;
        }
      }

      .register-content {
        padding: 40px 30px;
      }
    }
  }
}

// 深色主题适配
:deep(.van-theme-dark) {
  .register-page {
    .register-card {
      background: var(--el-bg-color-overlay);
      box-shadow: 0 20px 40px rgba(0, 0, 0, 0.3);
    }
  }
}

// 移动端响应式设计
@media (max-width: 768px) {
  .register-page {
    padding: 16px;
    background: var(--van-background);

    .register-container {
      max-width: 100%;

      .register-card {
        border-radius: 20px;
        box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);

        .register-header {
          padding: 30px 20px;
          background: linear-gradient(135deg, #10b981, #059669);

          .register-title {
            font-size: 24px;
          }

          .register-subtitle {
            font-size: 14px;
          }
        }

        .register-content {
          padding: 30px 20px;
        }
      }
    }
  }
}

// 小屏幕移动端优化
@media (max-width: 375px) {
  .register-page {
    padding: 12px;

    .register-container {
      .register-card {
        .register-header {
          padding: 24px 16px;

          .register-title {
            font-size: 22px;
          }

          .register-subtitle {
            font-size: 13px;
          }
        }

        .register-content {
          padding: 24px 16px;
        }
      }
    }
  }
}
</style>
