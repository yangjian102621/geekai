<template>
  <div class="register-page">
    <router-link to="/" class="back-home-btn" title="返回首页">
      <i class="iconfont icon-home"></i>
    </router-link>
    <div class="register-container">
      <div class="register-card">
        <div class="register-header">
          <h1 class="register-title">用户注册</h1>
          <p class="register-subtitle">创建您的账户以开始使用服务</p>
        </div>

        <div class="register-content">
          <login-dialog
            :show="true"
            active="register"
            :inviteCode="inviteCode"
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
const inviteCode = ref(router.currentRoute.value.query.invite_code || '')

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
  background: var(--theme-bg-all);
  background-image: var(--panel-bg);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
  position: relative;
  overflow: auto;
}

.back-home-btn {
  position: fixed;
  top: 20px;
  left: 20px;
  z-index: 1000;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 44px;
  height: 44px;
  background: var(--card-bg);
  border: 1px solid var(--line-box);
  border-radius: 12px;
  color: var(--theme-text-color-primary);
  text-decoration: none;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
  backdrop-filter: blur(8px);
  transition: all 0.3s ease;

  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 6px 20px rgba(0, 0, 0, 0.15);
    background: var(--hover-deep-color);
  }

  .iconfont {
    font-size: 20px;
  }
}

.register-container {
  width: 100%;
  max-width: 480px;
  margin: 0 auto;
}

.register-card {
  background: var(--card-bg);
  border: 1px solid var(--line-box);
  border-radius: 20px;
  padding: 40px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  backdrop-filter: blur(10px);
  position: relative;
  overflow: hidden;

  &::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    height: 4px;
    background: var(--btnColor);
  }
}

.register-header {
  text-align: center;
  margin-bottom: 32px;
}

.register-title {
  font-size: 28px;
  font-weight: 600;
  color: var(--theme-text-color-primary);
  margin: 0 0 8px 0;
  letter-spacing: -0.5px;
}

.register-subtitle {
  font-size: 16px;
  color: var(--theme-text-color-secondary);
  margin: 0;
  line-height: 1.5;
}

.register-content {
  :deep(.login-dialog) {
    .form {
      .block {
        margin-bottom: 20px;

        .el-input {
          .el-input__wrapper {
            background: var(--el-fill-color-blank);
            border: 1px solid var(--line-box);
            border-radius: 12px;
            box-shadow: none;
            transition: all 0.3s ease;

            &:hover,
            &.is-focus {
              border-color: var(--border-active);
              box-shadow: 0 0 0 3px rgba(91, 98, 206, 0.1);
            }
          }

          .el-input__inner {
            color: var(--theme-text-color-primary);
            font-size: 16px;

            &::placeholder {
              color: var(--theme-text-color-secondary);
              opacity: 0.7;
            }
          }

          .el-input__prefix,
          .el-input__suffix {
            color: var(--theme-text-color-secondary);
          }
        }
      }

      .btn-row {
        margin-top: 32px;

        .el-button {
          &[type="primary"], 
          &.register-btn {
            width: 100%;
            height: 48px;
            border-radius: 12px;
            background: var(--btnColor);
            border: none;
            font-size: 16px;
            font-weight: 500;
            transition: all 0.3s ease;
            box-shadow: 0 4px 16px rgba(91, 98, 206, 0.3);

            &:hover {
              transform: translateY(-2px);
              box-shadow: 0 8px 24px rgba(91, 98, 206, 0.4);
            }

            &:active {
              transform: translateY(0);
            }
          }
        }
      }

      .text {
        margin-top: 24px;
        color: var(--theme-text-color-secondary);

        .el-button {
          color: var(--text-color-primary);
          background: transparent;
          border: none;
          padding: 0 8px;
          font-size: 14px;
          
          &:hover {
            background: var(--btn-bg);
            border-radius: 6px;
          }

          &.forget {
            color: var(--theme-text-color-secondary);
            
            &:hover {
              color: var(--text-color-primary);
            }
          }
        }
      }

      // 验证码输入框样式
      .verify-code {
        .el-row {
          .el-col:first-child {
            .el-input__wrapper {
              border-top-right-radius: 0;
              border-bottom-right-radius: 0;
            }
          }

          .el-col:last-child {
            .el-button {
              height: 40px;
              border-top-left-radius: 0;
              border-bottom-left-radius: 0;
              border-left: none;
              background: var(--btn-bg);
              color: var(--theme-text-color-primary);
              border: 1px solid var(--line-box);
              font-size: 14px;

              &:hover {
                background: var(--hover-deep-color);
                color: var(--text-color-primary);
              }

              &:disabled {
                opacity: 0.6;
                cursor: not-allowed;
              }
            }
          }
        }
      }

      // 邀请码输入框样式
      .invite-code {
        .el-input__wrapper {
          background: var(--quote-bg-color);
          border-color: var(--border-active);
        }
      }

      // 协议条款样式
      .agreement {
        display: flex;
        align-items: flex-start;
        margin-bottom: 20px;
        font-size: 14px;
        color: var(--theme-text-color-secondary);
        line-height: 1.5;

        .el-checkbox {
          margin-right: 8px;
          flex-shrink: 0;
          
          :deep(.el-checkbox__inner) {
            border-color: var(--line-box);
            background: var(--el-fill-color-blank);
            
            &:hover {
              border-color: var(--border-active);
            }
          }

          :deep(.el-checkbox__input.is-checked) {
            .el-checkbox__inner {
              background: var(--btnColor);
              border-color: var(--border-active);
            }
          }
        }

        a {
          color: var(--text-color-primary);
          text-decoration: none;
          
          &:hover {
            text-decoration: underline;
          }
        }
      }
    }
  }
}

// 移动端适配
@media (max-width: 768px) {
  .register-page {
    padding: 16px;
  }

  .back-home-btn {
    top: 16px;
    left: 16px;
    width: 40px;
    height: 40px;

    .iconfont {
      font-size: 18px;
    }
  }

  .register-card {
    padding: 32px 24px;
    border-radius: 16px;
    margin-top: 60px;
  }

  .register-title {
    font-size: 24px;
  }

  .register-subtitle {
    font-size: 15px;
  }

  .register-content {
    :deep(.login-dialog) {
      .form {
        .block {
          margin-bottom: 18px;

          .el-input {
            .el-input__wrapper {
              border-radius: 10px;
            }

            .el-input__inner {
              font-size: 16px;
            }
          }
        }

        .btn-row {
          margin-top: 28px;

          .el-button {
            &[type="primary"], 
            &.register-btn {
              height: 46px;
              border-radius: 10px;
              font-size: 15px;
            }
          }
        }

        .text {
          margin-top: 20px;
          font-size: 13px;

          .el-button {
            font-size: 13px;
            padding: 0 6px;
          }
        }

        .verify-code {
          .el-row {
            .el-col:last-child {
              .el-button {
                height: 38px;
                font-size: 13px;
              }
            }
          }
        }

        .agreement {
          font-size: 13px;
        }
      }
    }
  }
}

// 小屏幕手机适配
@media (max-width: 480px) {
  .register-card {
    padding: 24px 20px;
  }

  .register-title {
    font-size: 22px;
  }

  .register-subtitle {
    font-size: 14px;
  }

  .register-content {
    :deep(.login-dialog) {
      .form {
        .verify-code {
          .el-row {
            .el-col:first-child {
              padding-right: 0;
            }
            .el-col:last-child {
              padding-left: 0;
              margin-top: 10px;
              
              .el-button {
                width: 100%;
                border-radius: 10px;
                border-left: 1px solid var(--line-box);
              }
            }
          }
        }

        .agreement {
          font-size: 12px;
        }
      }
    }
  }
}
</style>
