<template>
  <div class="profile-page">
    <div class="profile-header">
      <div class="header-bg"></div>
      <div class="header-content">
        <div class="user-info">
          <div class="avatar-container">
            <van-image :src="fileList[0].url" round width="80" height="80" />
          </div>
          <div class="user-details">
            <h2 class="username">{{ form.nickname || form.username }}</h2>
            <div class="user-meta">
              <van-tag type="info">剩余算力：{{ form.power || 0 }}</van-tag>
              <span class="user-id">ID: {{ form.id }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="profile-content">
      <!-- 快捷操作 -->
      <div class="quick-actions">
        <h3 class="section-title">快捷操作</h3>
        <van-row :gutter="12">
          <van-col :span="8">
            <div class="action-item" @click="router.push('/mobile/member')">
              <div class="action-icon recharge">
                <i class="iconfont icon-vip"></i>
              </div>
              <div class="action-label">会员中心</div>
            </div>
          </van-col>
          <van-col :span="8">
            <div class="action-item" @click="router.push('/mobile/invite')">
              <div class="action-icon share">
                <i class="iconfont icon-share"></i>
              </div>
              <div class="action-label">邀请</div>
            </div>
          </van-col>
          <van-col :span="8">
            <div class="action-item" @click="showSettings = true">
              <div class="action-icon settings">
                <i class="iconfont icon-config"></i>
              </div>
              <div class="action-label">设置</div>
            </div>
          </van-col>
        </van-row>
      </div>

      <!-- 我的服务 -->
      <div class="menu-section">
        <h3 class="section-title">我的服务</h3>
        <van-cell-group>
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
          <van-cell title="修改密码" is-link @click="showPasswordDialog = true">
            <template #icon>
              <i class="iconfont icon-password menu-icon"></i>
            </template>
          </van-cell>
          <van-cell
            title="消费记录"
            icon="notes-o"
            is-link
            @click="router.push('/mobile/power-log')"
          >
            <template #icon>
              <i class="iconfont icon-log menu-icon"></i>
            </template>
          </van-cell>
        </van-cell-group>
      </div>

      <!-- 退出登录 -->
      <div class="logout-section">
        <van-button size="large" block type="danger" plain @click="showLogoutConfirm = true">
          退出登录
        </van-button>
      </div>
    </div>

    <!-- 修改密码弹窗 -->
    <van-dialog
      v-model:show="showPasswordDialog"
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
    <van-action-sheet v-model:show="showSettings" title="设置">
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
        </van-cell-group>
      </div>
    </van-action-sheet>

    <!-- 绑定邮箱弹窗 -->
    <van-dialog
      v-model:show="showBindEmailDialog"
      title="绑定邮箱"
      :show-cancel-button="false"
      :show-confirm-button="false"
      width="90%"
      :close-on-click-overlay="true"
    >
      <div class="p-4">
        <bind-email @hide="showBindEmailDialog = false" />
      </div>
    </van-dialog>

    <!-- 绑定手机弹窗 -->
    <van-dialog
      v-model:show="showBindMobileDialog"
      title="绑定手机"
      :show-cancel-button="false"
      :show-confirm-button="false"
      width="90%"
      :close-on-click-overlay="true"
    >
      <div class="p-4">
        <bind-mobile @hide="showBindMobileDialog = false" />
      </div>
    </van-dialog>

    <!-- 退出登录确认 -->
    <van-dialog
      v-model:show="showLogoutConfirm"
      title="退出登录"
      message="确定要退出登录吗？"
      show-cancel-button
      @confirm="logout"
    />
  </div>
</template>

<script setup>
import BindEmail from '@/components/BindEmail.vue'
import BindMobile from '@/components/BindMobile.vue'
import { checkSession, getSystemInfo } from '@/store/cache'
import { removeUserToken } from '@/store/session'
import { useSharedStore } from '@/store/sharedata'
import { httpGet, httpPost } from '@/utils/http'
import { showFailToast, showLoadingToast, showSuccessToast } from 'vant'
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
const showSettings = ref(false)
const showPasswordDialog = ref(false)
const showBindEmailDialog = ref(false)
const showBindMobileDialog = ref(false)
const showLogoutConfirm = ref(false)
const store = useSharedStore()
const stream = ref(store.chatStream)
const dark = ref(store.theme === 'dark')
const title = ref(import.meta.env.VITE_TITLE)
const appVersion = ref('2.1.0')

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
      form.value = { ...form.value, ...user }
      fileList.value[0].url = user.avatar || '/images/avatar/default.jpg'

      // 获取用户详细信息
      fetchUserProfile()
    })
    .catch(() => {
      router.push('/login')
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
}

// 提交修改密码
const updatePass = () => {
  if (!pass.value.old) {
    return showFailToast('请输入旧密码')
  }
  if (!pass.value.new || pass.value.new.length < 8) {
    return showFailToast('密码长度为8-16个字符')
  }
  if (pass.value.renew !== pass.value.new) {
    return showFailToast('两次输入密码不一致')
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
      showSuccessToast('退出登录成功')
      router.push('/login')

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
  min-height: calc(100vh - 60px);
  background: var(--van-background);

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
}

// 深色主题优化
:deep(.van-theme-dark) {
  .profile-page {
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
