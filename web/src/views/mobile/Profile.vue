<template>
  <div class="profile-page">
    <div class="profile-header">
      <div class="header-bg"></div>
      <div class="header-content">
        <div class="user-info" v-if="isLogin">
          <div class="avatar-container" @click="showAvatarOptions = true">
            <van-image :src="fileList[0].url" round width="80" height="80" />
            <div class="avatar-badge">
              <van-icon name="photograph" />
            </div>
          </div>
          <div class="user-details">
            <h2 class="username">{{ form.nickname || form.username }}</h2>
            <div class="user-meta">
              <van-tag type="primary" v-if="isVip">VIP会员</van-tag>
              <van-tag type="default" v-else>普通用户</van-tag>
              <span class="user-id">ID: {{ form.id }}</span>
            </div>
          </div>
        </div>
        <div class="login-prompt" v-else>
          <van-button type="primary" size="large" round @click="showLoginDialog(router)">
            立即登录
          </van-button>
        </div>
      </div>
    </div>

    <div class="profile-content">
      <!-- 用户状态卡片 -->
      <div class="status-cards" v-if="isLogin">
        <van-row :gutter="12">
          <van-col :span="12">
            <div class="status-card" @click="router.push('/mobile/power-log')">
              <div class="card-icon power">
                <i class="iconfont icon-flash"></i>
              </div>
              <div class="card-value">{{ form.power || 0 }}</div>
              <div class="card-label">剩余算力</div>
            </div>
          </van-col>
          <van-col :span="12">
            <div class="status-card" @click="router.push('/mobile/invite')">
              <div class="card-icon invite">
                <i class="iconfont icon-user-plus"></i>
              </div>
              <div class="card-value">{{ inviteCount }}</div>
              <div class="card-label">邀请好友</div>
            </div>
          </van-col>
        </van-row>
      </div>

      <!-- 快捷操作 -->
      <div class="quick-actions" v-if="isLogin">
        <h3 class="section-title">快捷操作</h3>
        <van-row :gutter="12">
          <van-col :span="6">
            <div class="action-item" @click="router.push('/mobile/member')">
              <div class="action-icon recharge">
                <i class="iconfont icon-vip"></i>
              </div>
              <div class="action-label">会员中心</div>
            </div>
          </van-col>
          <van-col :span="6">
            <div class="action-item" @click="showPasswordDialog = true">
              <div class="action-icon password">
                <i class="iconfont icon-lock"></i>
              </div>
              <div class="action-label">改密码</div>
            </div>
          </van-col>
          <van-col :span="6">
            <div class="action-item" @click="router.push('/mobile/invite')">
              <div class="action-icon share">
                <i class="iconfont icon-share"></i>
              </div>
              <div class="action-label">邀请</div>
            </div>
          </van-col>
          <van-col :span="6">
            <div class="action-item" @click="showSettings = true">
              <div class="action-icon settings">
                <i class="iconfont icon-setting"></i>
              </div>
              <div class="action-label">设置</div>
            </div>
          </van-col>
        </van-row>
      </div>

      <!-- 账户管理 -->
      <div class="menu-section" v-if="isLogin">
        <h3 class="section-title">账户管理</h3>
        <van-cell-group inset>
          <van-cell title="绑定邮箱" is-link @click="showBindEmailDialog = true">
            <template #icon>
              <i class="iconfont icon-email menu-icon"></i>
            </template>
          </van-cell>
          <van-cell title="绑定手机" is-link @click="showBindMobileDialog = true">
            <template #icon>
              <i class="iconfont icon-mobile menu-icon"></i>
            </template>
          </van-cell>
          <van-cell title="第三方登录" is-link @click="showThirdLoginDialog = true">
            <template #icon>
              <i class="iconfont icon-login menu-icon"></i>
            </template>
          </van-cell>
        </van-cell-group>
      </div>

      <!-- 功能菜单 -->
      <div class="menu-section">
        <h3 class="section-title">我的服务</h3>
        <van-cell-group inset>
          <van-cell
            title="消费记录"
            icon="notes-o"
            is-link
            @click="router.push('/mobile/power-log')"
          >
            <template #icon>
              <i class="iconfont icon-history menu-icon"></i>
            </template>
          </van-cell>
          <van-cell
            title="邀请好友"
            icon="friends-o"
            is-link
            @click="router.push('/mobile/invite')"
          >
            <template #icon>
              <i class="iconfont icon-user-plus menu-icon"></i>
            </template>
          </van-cell>
          <van-cell title="聊天导出" icon="down" is-link @click="copyChatExportLink">
            <template #icon>
              <i class="iconfont icon-download menu-icon"></i>
            </template>
          </van-cell>
        </van-cell-group>

        <van-cell-group inset>
          <van-cell title="帮助中心" icon="question-o" is-link @click="router.push('/mobile/help')">
            <template #icon>
              <i class="iconfont icon-help menu-icon"></i>
            </template>
          </van-cell>
          <van-cell title="意见反馈" icon="chat-o" is-link @click="router.push('/mobile/feedback')">
            <template #icon>
              <i class="iconfont icon-message menu-icon"></i>
            </template>
          </van-cell>
          <van-cell title="关于我们" icon="info-o" is-link @click="showAbout = true">
            <template #icon>
              <i class="iconfont icon-info menu-icon"></i>
            </template>
          </van-cell>
        </van-cell-group>
      </div>

      <!-- 退出登录 -->
      <div class="logout-section" v-if="isLogin">
        <van-button size="large" block type="danger" plain @click="showLogoutConfirm = true">
          退出登录
        </van-button>
      </div>

      <!-- 版本信息 -->
      <div class="version-info">
        <p class="app-version">版本 v{{ appVersion }}</p>
        <p class="copyright">© 2024 {{ title }}. All rights reserved.</p>
      </div>

      <!-- 底部安全间距 -->
      <div class="bottom-safe-area"></div>
    </div>

    <!-- 修改密码弹窗 -->
    <van-dialog
      :model-value="showPasswordDialog"
      @update:model-value="showPasswordDialog = $event"
      title="修改密码"
      show-cancel-button
      @confirm="updatePass"
      @cancel="resetPasswordForm"
    >
      <van-form ref="passwordForm" @submit="updatePass">
        <van-cell-group inset>
          <van-field
            v-model="pass.old"
            type="password"
            label="旧密码"
            placeholder="请输入旧密码"
            required
            :rules="[{ required: true, message: '请输入旧密码' }]"
          />
          <van-field
            v-model="pass.new"
            type="password"
            label="新密码"
            placeholder="请输入新密码"
            required
            :rules="passwordRules"
          />
          <van-field
            v-model="pass.renew"
            type="password"
            label="确认密码"
            placeholder="请再次输入新密码"
            required
            :rules="[
              { required: true, message: '请再次输入新密码' },
              { validator: validateConfirmPassword },
            ]"
          />
        </van-cell-group>
      </van-form>
    </van-dialog>

    <!-- 设置弹窗 -->
    <van-action-sheet
      :model-value="showSettings"
      @update:model-value="showSettings = $event"
      title="设置"
    >
      <div class="settings-content">
        <van-cell-group>
          <van-cell title="暗黑主题">
            <template #right-icon>
              <van-switch
                v-model="dark"
                @change="(val) => store.setTheme(val ? 'dark' : 'light')"
              />
            </template>
          </van-cell>
          <van-cell title="流式输出">
            <template #right-icon>
              <van-switch v-model="stream" @change="(val) => store.setChatStream(val)" />
            </template>
          </van-cell>
          <van-cell title="消息通知">
            <template #right-icon>
              <van-switch v-model="notifications" />
            </template>
          </van-cell>
          <van-cell title="自动保存">
            <template #right-icon>
              <van-switch v-model="autoSave" />
            </template>
          </van-cell>
        </van-cell-group>
      </div>
    </van-action-sheet>

    <!-- 头像选择弹窗 -->
    <van-action-sheet
      :model-value="showAvatarOptions"
      @update:model-value="showAvatarOptions = $event"
      title="更换头像"
    >
      <div class="avatar-options">
        <van-cell title="拍照" icon="photograph" @click="selectAvatar('camera')" />
        <van-cell title="从相册选择" icon="photo-o" @click="selectAvatar('album')" />
        <van-cell title="默认头像" icon="user-o" @click="selectAvatar('default')" />
      </div>
    </van-action-sheet>

    <!-- 关于我们弹窗 -->
    <van-dialog
      :model-value="showAbout"
      @update:model-value="showAbout = $event"
      title="关于我们"
      :show-cancel-button="false"
    >
      <div class="about-content">
        <div class="about-logo">
          <img src="/images/logo.png" alt="Logo" />
        </div>
        <h3>{{ title }}</h3>
        <p class="about-desc">
          专业的AI创作平台，提供对话、绘画、音乐、视频等多种AI服务，让创作更简单、更高效。
        </p>
        <div class="about-info">
          <p>版本：v{{ appVersion }}</p>
          <p>更新时间：2024-01-01</p>
        </div>
      </div>
    </van-dialog>

    <!-- 组件弹窗 -->
    <bind-email v-if="isLogin" :show="showBindEmailDialog" @hide="showBindEmailDialog = false" />
    <bind-mobile v-if="isLogin" :show="showBindMobileDialog" @hide="showBindMobileDialog = false" />
    <third-login v-if="isLogin" :show="showThirdLoginDialog" @hide="showThirdLoginDialog = false" />

    <!-- 退出登录确认 -->
    <van-dialog
      :model-value="showLogoutConfirm"
      @update:model-value="showLogoutConfirm = $event"
      title="退出登录"
      message="确定要退出登录吗？"
      show-cancel-button
      @confirm="logout"
    />

    <!-- 隐藏的复制链接按钮 -->
    <button id="copy-chat-export-btn" style="display: none" :data-clipboard-text="chatExportUrl">
      复制聊天导出链接
    </button>
  </div>
</template>

<script setup>
import BindEmail from '@/components/BindEmail.vue'
import BindMobile from '@/components/BindMobile.vue'
import ThirdLogin from '@/components/ThirdLogin.vue'
import { checkSession, getSystemInfo } from '@/store/cache'
import { removeUserToken } from '@/store/session'
import { useSharedStore } from '@/store/sharedata'
import { httpGet, httpPost } from '@/utils/http'
import { showLoginDialog } from '@/utils/libs'
import Clipboard from 'clipboard'
import { showFailToast, showLoadingToast, showNotify, showSuccessToast } from 'vant'
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'

const form = ref({
  id: 0,
  username: 'GeekMaster',
  nickname: '极客学长@001',
  mobile: '1300000000',
  avatar: '',
  power: 0,
  expired_time: 0,
})

const fileList = ref([
  {
    url: '/images/avatar/default.jpg',
    message: '上传中...',
  },
])

const router = useRouter()
const isLogin = ref(false)
const showSettings = ref(false)
const showPasswordDialog = ref(false)
const showBindEmailDialog = ref(false)
const showBindMobileDialog = ref(false)
const showThirdLoginDialog = ref(false)
const showAvatarOptions = ref(false)
const showAbout = ref(false)
const showLogoutConfirm = ref(false)
const store = useSharedStore()
const stream = ref(store.chatStream)
const dark = ref(store.theme === 'dark')
const title = ref(import.meta.env.VITE_TITLE)
const appVersion = ref('2.1.0')

// 聊天导出链接
const chatExportUrl = ref(location.protocol + '//' + location.host + '/chat/export')

// 新增状态
const notifications = ref(true)
const autoSave = ref(true)
const inviteCount = ref(0)
const passwordForm = ref()

// 密码相关
const pass = ref({
  old: '',
  new: '',
  renew: '',
})

// 密码验证规则
const passwordRules = [
  { required: true, message: '请输入新密码' },
  { min: 8, max: 16, message: '密码长度为8-16个字符' },
]

// 计算属性
const isVip = computed(() => {
  const now = Date.now()
  const expiredTime = form.value.expired_time ? form.value.expired_time * 1000 : 0
  return expiredTime > now
})

onMounted(() => {
  getSystemInfo()
    .then((res) => {
      title.value = res.data.title
    })
    .catch((e) => {
      console.error('获取系统配置失败：', e.message)
    })

  checkSession()
    .then((user) => {
      isLogin.value = true
      form.value = { ...form.value, ...user }
      fileList.value[0].url = user.avatar || '/images/avatar/default.jpg'

      // 获取用户详细信息
      fetchUserProfile()
      fetchUserStats()
    })
    .catch(() => {
      isLogin.value = false
    })

  // 初始化复制功能
  const clipboard = new Clipboard('#copy-chat-export-btn')
  clipboard.on('success', (e) => {
    e.clearSelection()
    showNotify({ type: 'success', message: '链接已复制到剪贴板', duration: 2000 })
  })
  clipboard.on('error', () => {
    showNotify({ type: 'danger', message: '复制失败', duration: 2000 })
  })
})

// 获取用户详细信息
const fetchUserProfile = () => {
  httpGet('/api/user/profile')
    .then((res) => {
      form.value = { ...form.value, ...res.data }
      fileList.value[0].url = res.data.avatar || '/images/avatar/default.jpg'
    })
    .catch((e) => {
      console.error('获取用户信息失败:', e.message)
    })
}

// 获取用户统计信息
const fetchUserStats = () => {
  // 模拟数据，实际项目中应调用API
  inviteCount.value = Math.floor(Math.random() * 20)
}

// 复制聊天导出链接
const copyChatExportLink = () => {
  document.getElementById('copy-chat-export-btn').click()
}

// 确认密码验证
const validateConfirmPassword = (value) => {
  if (value !== pass.value.new) {
    return Promise.reject(new Error('两次输入的密码不一致'))
  }
  return Promise.resolve()
}

// 重置密码表单
const resetPasswordForm = () => {
  pass.value = {
    old: '',
    new: '',
    renew: '',
  }
  if (passwordForm.value) {
    passwordForm.value.resetValidation()
  }
}

// 提交修改密码
const updatePass = () => {
  if (!passwordForm.value) {
    updatePasswordAPI()
    return
  }

  passwordForm.value
    .validate()
    .then(() => {
      updatePasswordAPI()
    })
    .catch((errors) => {
      console.log('表单验证失败:', errors)
    })
}

const updatePasswordAPI = () => {
  if (!pass.value.old) {
    return showNotify({ type: 'danger', message: '请输入旧密码' })
  }
  if (!pass.value.new || pass.value.new.length < 8) {
    return showNotify({ type: 'danger', message: '密码长度为8-16个字符' })
  }
  if (pass.value.renew !== pass.value.new) {
    return showNotify({ type: 'danger', message: '两次输入密码不一致' })
  }

  showLoadingToast({
    message: '正在修改密码...',
    forbidClick: true,
  })

  httpPost('/api/user/password', {
    old_pass: pass.value.old,
    password: pass.value.new,
    repass: pass.value.renew,
  })
    .then(() => {
      showSuccessToast('密码修改成功！')
      showPasswordDialog.value = false
      resetPasswordForm()
    })
    .catch((e) => {
      showFailToast('密码修改失败：' + e.message)
    })
}

// 头像选择
const selectAvatar = (type) => {
  showAvatarOptions.value = false

  switch (type) {
    case 'camera':
      // 调用相机
      if (navigator.mediaDevices && navigator.mediaDevices.getUserMedia) {
        showNotify({ type: 'primary', message: '正在启动相机...' })
      } else {
        showNotify({ type: 'warning', message: '您的设备不支持相机功能' })
      }
      break
    case 'album':
      // 从相册选择
      const input = document.createElement('input')
      input.type = 'file'
      input.accept = 'image/*'
      input.onchange = (e) => {
        const file = e.target.files[0]
        if (file) {
          // 这里应该上传到服务器
          const reader = new FileReader()
          reader.onload = (e) => {
            fileList.value[0].url = e.target.result
            showSuccessToast('头像更新成功')
          }
          reader.readAsDataURL(file)
        }
      }
      input.click()
      break
    case 'default':
      // 使用默认头像
      fileList.value[0].url = '/images/avatar/default.jpg'
      showSuccessToast('已设置为默认头像')
      break
  }
}

// 退出登录
const logout = function () {
  showLoadingToast({
    message: '正在退出...',
    forbidClick: true,
  })

  httpGet('/api/user/logout')
    .then(() => {
      removeUserToken()
      store.setIsLogin(false)
      isLogin.value = false
      showSuccessToast('退出登录成功')
      showLogoutConfirm.value = false

      // 清除用户数据
      form.value = {
        id: 0,
        username: '',
        nickname: '',
        mobile: '',
        avatar: '',
        power: 0,
        expired_time: 0,
      }
      fileList.value[0].url = '/images/avatar/default.jpg'
    })
    .catch((e) => {
      showFailToast('退出登录失败：' + e.message)
    })
}
</script>

<style lang="scss" scoped>
.profile-page {
  min-height: 100vh;
  background: var(--van-background);
  padding-bottom: 60px;

  .profile-header {
    position: relative;
    height: 240px;
    overflow: hidden;
    background: linear-gradient(135deg, var(--van-primary-color), #8b5cf6);

    .header-bg {
      position: absolute;
      top: 0;
      left: 0;
      right: 0;
      bottom: 0;
      background: url('/images/profile-bg.png') center/cover;
      opacity: 0.3;
    }

    .header-content {
      position: relative;
      z-index: 2;
      display: flex;
      flex-direction: column;
      justify-content: center;
      align-items: center;
      height: 100%;
      padding: 20px;
      color: white;

      .user-info {
        text-align: center;

        .avatar-container {
          position: relative;
          display: inline-block;
          margin-bottom: 16px;
          cursor: pointer;

          .avatar-badge {
            position: absolute;
            bottom: 0;
            right: 0;
            width: 24px;
            height: 24px;
            background: rgba(255, 255, 255, 0.9);
            border-radius: 50%;
            display: flex;
            align-items: center;
            justify-content: center;
            color: var(--van-primary-color);
            font-size: 12px;
          }
        }

        .user-details {
          .username {
            font-size: 22px;
            font-weight: 700;
            margin: 0 0 8px 0;
            text-shadow: 0 2px 4px rgba(0, 0, 0, 0.3);
          }

          .user-meta {
            display: flex;
            align-items: center;
            justify-content: center;
            gap: 12px;

            .user-id {
              font-size: 12px;
              opacity: 0.8;
            }
          }
        }
      }

      .login-prompt {
        text-align: center;
      }
    }
  }

  .profile-content {
    margin-top: 20px;
    z-index: 3;
    padding: 0 16px 20px;

    .status-cards {
      margin-bottom: 24px;

      .status-card {
        background: var(--van-cell-background);
        border-radius: 12px;
        padding: 16px;
        text-align: center;
        box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
        cursor: pointer;
        transition: all 0.3s ease;

        &:active {
          transform: scale(0.95);
        }

        .card-icon {
          width: 48px;
          height: 48px;
          border-radius: 12px;
          display: flex;
          align-items: center;
          justify-content: center;
          margin: 0 auto 8px;

          &.power {
            background: linear-gradient(135deg, #ff9500, #ff6b35);
          }

          &.invite {
            background: linear-gradient(135deg, #1989fa, #0d7dff);
          }

          .iconfont {
            font-size: 24px;
            color: white;
          }
        }

        .card-value {
          font-size: 18px;
          font-weight: 700;
          color: var(--van-text-color);
          margin-bottom: 4px;
        }

        .card-label {
          font-size: 12px;
          color: var(--van-gray-6);
        }
      }
    }

    .quick-actions,
    .menu-section {
      margin-bottom: 24px;

      .section-title {
        font-size: 18px;
        font-weight: 600;
        color: var(--van-text-color);
        margin: 0 0 16px 4px;
      }
    }

    .quick-actions {
      .action-item {
        display: flex;
        flex-direction: column;
        align-items: center;
        padding: 16px 8px;
        background: var(--van-cell-background);
        border-radius: 12px;
        box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
        cursor: pointer;
        transition: all 0.3s ease;

        &:active {
          transform: scale(0.95);
        }

        .action-icon {
          width: 40px;
          height: 40px;
          border-radius: 10px;
          display: flex;
          align-items: center;
          justify-content: center;
          margin-bottom: 8px;

          &.recharge {
            background: linear-gradient(135deg, #ffd700, #ffb300);
          }

          &.password {
            background: linear-gradient(135deg, #ee0a24, #d60a21);
          }

          &.share {
            background: linear-gradient(135deg, #8b5cf6, #7c3aed);
          }

          &.settings {
            background: linear-gradient(135deg, #6b7280, #4b5563);
          }

          .iconfont {
            font-size: 18px;
            color: white;
          }
        }

        .action-label {
          font-size: 12px;
          color: var(--van-text-color);
          text-align: center;
        }
      }
    }

    .menu-section {
      :deep(.van-cell-group) {
        border-radius: 12px;
        overflow: hidden;
        box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
        margin-bottom: 12px;

        .van-cell {
          padding: 16px;

          .menu-icon {
            font-size: 18px;
            margin-right: 12px;
            color: var(--van-primary-color);
          }

          .van-cell__title {
            font-size: 15px;
            font-weight: 500;
          }
        }
      }
    }

    .logout-section {
      margin-bottom: 24px;
    }

    .version-info {
      text-align: center;
      padding: 20px 0;

      .app-version,
      .copyright {
        font-size: 12px;
        color: var(--van-gray-6);
        margin: 0 0 4px 0;
      }

      .copyright {
        margin: 0;
      }
    }

    .bottom-safe-area {
      height: 20px;
      width: 100%;
    }
  }

  // 弹窗样式
  .settings-content {
    padding: 16px;

    :deep(.van-cell-group) {
      .van-cell {
        padding: 16px;

        .van-cell__title {
          font-size: 15px;
          font-weight: 500;
        }
      }
    }
  }

  .avatar-options {
    padding: 0;

    :deep(.van-cell) {
      padding: 16px 20px;

      .van-cell__title {
        font-size: 15px;
        font-weight: 500;
      }
    }
  }

  .about-content {
    text-align: center;
    padding: 20px;

    .about-logo {
      margin-bottom: 16px;

      img {
        width: 60px;
        height: 60px;
        border-radius: 12px;
      }
    }

    h3 {
      font-size: 20px;
      font-weight: 600;
      color: var(--van-text-color);
      margin: 0 0 12px 0;
    }

    .about-desc {
      font-size: 14px;
      color: var(--van-gray-6);
      line-height: 1.5;
      margin: 0 0 20px 0;
    }

    .about-info {
      p {
        font-size: 13px;
        color: var(--van-gray-7);
        margin: 0 0 4px 0;

        &:last-child {
          margin: 0;
        }
      }
    }
  }
}

// 深色主题优化
:deep(.van-theme-dark) {
  .profile-page {
    .status-card,
    .action-item {
      box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
    }

    .van-cell-group {
      box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
    }
  }
}

// 响应式优化
@media (max-width: 375px) {
  .profile-page {
    .profile-header {
      height: 220px;

      .header-content .user-info .username {
        font-size: 20px;
      }
    }

    .profile-content {
      padding: 0 12px 20px;

      .status-cards .status-card {
        padding: 12px;

        .card-value {
          font-size: 16px;
        }
      }

      .quick-actions .action-item {
        padding: 12px 6px;

        .action-icon {
          width: 36px;
          height: 36px;

          .iconfont {
            font-size: 16px;
          }
        }
      }
    }
  }
}
</style>
