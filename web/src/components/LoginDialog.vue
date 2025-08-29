<template>
  <div class="login-dialog w-full">
    <div class="login-box" v-if="login">
      <el-form :model="data" class="form space-y-5">
        <div class="block">
          <el-input placeholder="账号" size="large" v-model="data.username" autocomplete="off">
            <template #prefix>
              <el-icon>
                <Iphone />
              </el-icon>
            </template>
          </el-input>
        </div>

        <div class="block">
          <el-input
            placeholder="请输入密码(8-16位)"
            maxlength="16"
            size="large"
            v-model="data.password"
            show-password
            autocomplete="off"
          >
            <template #prefix>
              <el-icon>
                <Lock />
              </el-icon>
            </template>
          </el-input>
        </div>

        <el-row class="btn-row mt-8" :gutter="20">
          <el-col :span="24">
            <button
              class="w-full h-12 rounded-xl text-base font-medium text-white bg-gradient-to-r from-blue-500 to-purple-600 hover:from-blue-600 hover:to-purple-700 transition-all duration-300 hover:-translate-y-0.5 hover:shadow-lg active:translate-y-0 shadow-md"
              @click="submitLogin"
              type="button"
            >
              {{ loading ? '登录中...' : '登 录' }}
            </button>
          </el-col>
        </el-row>

        <div class="w-full">
          <div
            class="text flex justify-center items-center pt-3 text-sm"
            style="color: var(--login-text-color)"
          >
            还没有账号？
            <el-button
              size="small"
              class="ml-2 rounded-md px-2 py-1 transition-colors duration-200"
              style="color: var(--login-link-color)"
              @click="login = false"
              @mouseenter="$event.target.style.background = 'var(--login-link-hover-bg)'"
              @mouseleave="$event.target.style.background = 'transparent'"
              >注册</el-button
            >

            <el-button type="info" class="forget ml-4" size="small" @click="showResetPass = true"
              >忘记密码？</el-button
            >
          </div>
        </div>
      </el-form>
    </div>

    <div class="register-box w-full" v-else>
      <el-form :model="data" class="form space-y-5" v-if="enableRegister">
        <el-tabs v-model="activeName" class="demo-tabs dark:text-white">
          <el-tab-pane label="手机注册" name="mobile" v-if="enableMobile">
            <div class="block">
              <el-input
                placeholder="手机号码"
                size="large"
                v-model="data.mobile"
                maxlength="11"
                autocomplete="off"
              >
                <template #prefix>
                  <el-icon>
                    <Iphone />
                  </el-icon>
                </template>
              </el-input>
            </div>
            <div class="block mt-4">
              <el-row :gutter="10">
                <el-col :span="12">
                  <el-input
                    placeholder="验证码"
                    size="large"
                    maxlength="30"
                    v-model="data.code"
                    autocomplete="off"
                  >
                    <template #prefix>
                      <el-icon>
                        <Checked />
                      </el-icon>
                    </template>
                  </el-input>
                </el-col>
                <el-col :span="12">
                  <send-msg size="large" :receiver="data.mobile" type="mobile" />
                </el-col>
              </el-row>
            </div>
          </el-tab-pane>
          <el-tab-pane label="邮箱注册" name="email" v-if="enableEmail">
            <div class="block">
              <el-input placeholder="邮箱地址" size="large" v-model="data.email" autocomplete="off">
                <template #prefix>
                  <el-icon>
                    <Message />
                  </el-icon>
                </template>
              </el-input>
            </div>
            <div class="block mt-4">
              <el-row :gutter="10">
                <el-col :span="12">
                  <el-input
                    placeholder="验证码"
                    size="large"
                    maxlength="30"
                    v-model="data.code"
                    autocomplete="off"
                  >
                    <template #prefix>
                      <el-icon>
                        <Checked />
                      </el-icon>
                    </template>
                  </el-input>
                </el-col>
                <el-col :span="12">
                  <send-msg size="large" :receiver="data.email" type="email" />
                </el-col>
              </el-row>
            </div>
          </el-tab-pane>
          <el-tab-pane label="用户名注册" name="username" v-if="enableUser">
            <div class="block">
              <el-input
                placeholder="用户名"
                size="large"
                v-model="data.username"
                autocomplete="off"
              >
                <template #prefix>
                  <el-icon>
                    <Iphone />
                  </el-icon>
                </template>
              </el-input>
            </div>
          </el-tab-pane>
        </el-tabs>

        <div class="block">
          <el-input
            placeholder="请输入密码(8-16位)"
            maxlength="16"
            size="large"
            v-model="data.password"
            show-password
            autocomplete="off"
          >
            <template #prefix>
              <el-icon>
                <Lock />
              </el-icon>
            </template>
          </el-input>
        </div>

        <div class="block">
          <el-input
            placeholder="重复密码(8-16位)"
            size="large"
            maxlength="16"
            v-model="data.repass"
            show-password
            autocomplete="off"
          >
            <template #prefix>
              <el-icon>
                <Lock />
              </el-icon>
            </template>
          </el-input>
        </div>

        <div class="block">
          <el-input
            placeholder="邀请码(可选)"
            size="large"
            v-model="data.invite_code"
            autocomplete="off"
          >
            <template #prefix>
              <el-icon>
                <Message />
              </el-icon>
            </template>
          </el-input>
        </div>

        <div class="block text-sm">
          <el-checkbox v-model="agreeChecked">
            我已阅读并同意
            <a href="javascript:void(0)" class="text-blue-500" @click="openAgreement"
              >《用户协议》</a
            >
            和
            <a href="javascript:void(0)" class="text-blue-500" @click="openPrivacy">《隐私政策》</a>
          </el-checkbox>
        </div>

        <div class="w-full">
          <button
            class="w-full h-12 rounded-xl text-base font-medium text-white bg-gradient-to-r from-blue-500 to-purple-600 hover:from-blue-600 hover:to-purple-700 transition-all duration-300 hover:-translate-y-0.5 hover:shadow-lg active:translate-y-0 shadow-md"
            @click="submitRegister"
            type="button"
          >
            {{ loading ? '注册中...' : '注 册' }}
          </button>
        </div>

        <div
          class="text text-sm flex justify-center items-center w-full pt-3"
          style="color: var(--login-text-color)"
        >
          已有账号？
          <el-button
            size="small"
            class="ml-2 rounded-md px-2 py-1 transition-colors duration-200"
            style="color: var(--login-link-color)"
            @click="login = true"
            @mouseenter="$event.target.style.background = 'var(--login-link-hover-bg)'"
            @mouseleave="$event.target.style.background = 'transparent'"
            >登录</el-button
          >
        </div>
      </el-form>

      <div class="tip-result" v-else>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-result icon="error" title="注册功能已关闭">
              <template #sub-title>
                <p>抱歉，系统已关闭注册功能，请联系管理员添加账号！</p>
              </template>
            </el-result>
          </el-col>

          <el-col :span="12">
            <div class="wechat-card">
              <el-image :src="wxImg" />
            </div>
          </el-col>
        </el-row>
      </div>
    </div>
    <captcha v-if="enableCaptcha" :type="captchaType" @success="submit" ref="captchaRef" />

    <reset-pass @hide="showResetPass = false" :show="showResetPass" />

    <el-dialog v-model="showAgreement" title="用户协议" :append-to-body="true">
      <div class="prose" v-html="agreementHtml"></div>
    </el-dialog>

    <el-dialog v-model="showPrivacy" title="隐私政策" :append-to-body="true">
      <div class="prose" v-html="privacyHtml"></div>
    </el-dialog>
  </div>
</template>

<script setup>
import Captcha from '@/components/Captcha.vue'
import ResetPass from '@/components/ResetPass.vue'
import SendMsg from '@/components/SendMsg.vue'
import { getSystemInfo } from '@/store/cache'
import { setUserToken } from '@/store/session'
import { useSharedStore } from '@/store/sharedata'
import { httpGet, httpPost } from '@/utils/http'
import { arrayContains } from '@/utils/libs'
import { validateEmail, validateMobile } from '@/utils/validate'
import { Checked, Iphone, Lock, Message } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { onMounted, ref, watch } from 'vue'
import { marked } from 'marked'

// eslint-disable-next-line no-undef
const props = defineProps({
  show: Boolean,
  active: {
    type: String,
    default: 'login',
  },
  inviteCode: {
    type: String,
    default: '',
  },
})
const showDialog = ref(false)
watch(
  () => props.show,
  (newValue) => {
    showDialog.value = newValue
  }
)

const login = ref(props.active === 'login')
const data = ref({
  username: import.meta.env.VITE_USER,
  password: import.meta.env.VITE_PASS,
  mobile: '',
  email: '',
  repass: '',
  code: '',
  invite_code: props.inviteCode,
})
const enableMobile = ref(false)
const enableEmail = ref(false)
const enableUser = ref(false)
const enableRegister = ref(true)

const activeName = ref('')
const wxImg = ref('/images/wx.png')
const captchaRef = ref(null)
// eslint-disable-next-line no-undef
const emits = defineEmits(['hide', 'success'])
const action = ref('login')
const enableCaptcha = ref(false)
const captchaType = ref('')
const showResetPass = ref(false)
const store = useSharedStore()
const loading = ref(false)
const agreeChecked = ref(false)
const showAgreement = ref(false)
const showPrivacy = ref(false)
const agreementHtml = ref('')
const privacyHtml = ref('')

onMounted(() => {
  getSystemInfo()
    .then((res) => {
      if (res.data) {
        const registerWays = res.data['register_ways']
        if (arrayContains(registerWays, 'username')) {
          enableUser.value = true
          activeName.value = 'username'
        }
        if (arrayContains(registerWays, 'email')) {
          enableEmail.value = true
          activeName.value = 'email'
        }
        if (arrayContains(registerWays, 'mobile')) {
          enableMobile.value = true
          activeName.value = 'mobile'
        }
        // 是否启用注册
        enableRegister.value = res.data['enabled_register']
        // 使用后台上传的客服微信二维码
        if (res.data['wechat_card_url'] !== '') {
          wxImg.value = res.data['wechat_card_url']
        }
      }
    })
    .catch((e) => {
      ElMessage.error('获取系统配置失败：' + e.message)
    })

  httpGet('/api/captcha/config').then((res) => {
    enableCaptcha.value = res.data['enabled']
    captchaType.value = res.data['type']
  })
})

const submit = (verifyData) => {
  if (action.value === 'login') {
    doLogin(verifyData)
  } else if (action.value === 'register') {
    doRegister(verifyData)
  }
}

// 登录操作
const submitLogin = () => {
  if (!data.value.username) {
    return ElMessage.error('请输入用户名')
  }
  if (!data.value.password) {
    return ElMessage.error('请输入密码')
  }
  if (enableCaptcha.value) {
    captchaRef.value.loadCaptcha()
    action.value = 'login'
  } else {
    doLogin({})
  }
}

const doLogin = (verifyData) => {
  data.value.key = verifyData.key
  data.value.dots = verifyData.dots
  data.value.x = verifyData.x
  loading.value = true
  httpPost('/api/user/login', data.value)
    .then((res) => {
      setUserToken(res.data.token)
      store.setIsLogin(true)
      ElMessage.success('登录成功！')
      emits('hide')
      emits('success')
    })
    .catch((e) => {
      ElMessage.error('登录失败，' + e.message)
    })
    .finally(() => {
      loading.value = false
    })
}

// 注册操作
const submitRegister = () => {
  if (activeName.value === 'username' && data.value.username === '') {
    return ElMessage.error('请输入用户名')
  }

  if (activeName.value === 'mobile' && !validateMobile(data.value.mobile)) {
    return ElMessage.error('请输入合法的手机号')
  }

  if (activeName.value === 'email' && !validateEmail(data.value.email)) {
    return ElMessage.error('请输入合法的邮箱地址')
  }

  if (data.value.password.length < 8) {
    return ElMessage.error('密码的长度为8-16个字符')
  }
  if (data.value.repass !== data.value.password) {
    return ElMessage.error('两次输入密码不一致')
  }

  if ((activeName.value === 'mobile' || activeName.value === 'email') && data.value.code === '') {
    return ElMessage.error('请输入验证码')
  }
  if (!agreeChecked.value) {
    return ElMessage.error('请先阅读并同意《用户协议》和《隐私政策》')
  }
  if (enableCaptcha.value && activeName.value === 'username') {
    captchaRef.value.loadCaptcha()
    action.value = 'register'
  } else {
    doRegister({})
  }
}

const doRegister = (verifyData) => {
  data.value.key = verifyData.key
  data.value.dots = verifyData.dots
  data.value.x = verifyData.x
  data.value.reg_way = activeName.value
  loading.value = true
  httpPost('/api/user/register', data.value)
    .then((res) => {
      setUserToken(res.data.token)
      ElMessage.success({
        message: '注册成功!',
        onClose: () => {
          emits('hide')
          emits('success')
        },
        duration: 1000,
      })
    })
    .catch((e) => {
      ElMessage.error('注册失败，' + e.message)
    })
    .finally(() => {
      loading.value = false
    })
}

// 打开并加载协议
const openAgreement = () => {
  if (!agreementHtml.value) {
    httpGet('/api/config/get?key=agreement')
      .then((res) => {
        agreementHtml.value = marked.parse(res.data?.content || '')
        showAgreement.value = true
      })
      .catch((e) => ElMessage.error('加载用户协议失败：' + e.message))
  } else {
    showAgreement.value = true
  }
}

// 打开并加载隐私政策
const openPrivacy = () => {
  if (!privacyHtml.value) {
    httpGet('/api/config/get?key=privacy')
      .then((res) => {
        privacyHtml.value = marked.parse(res.data?.content || '')
        showPrivacy.value = true
      })
      .catch((e) => ElMessage.error('加载隐私政策失败：' + e.message))
  } else {
    showPrivacy.value = true
  }
}
</script>

<style lang="scss">
.login-dialog {
  border-radius: 10px;

  // Dark theme support for Element Plus components
  :deep(.el-tabs) {
    .el-tabs__header {
      .el-tabs__nav-wrap {
        .el-tabs__nav {
          .el-tabs__item {
            color: var(--el-text-color-primary);

            &.is-active {
              color: var(--el-color-primary);
            }
          }
        }
      }
    }
  }

  :deep(.el-input) {
    .el-input__wrapper {
      background: var(--el-fill-color-blank);
      border-color: var(--el-border-color);

      &.is-focus {
        border-color: var(--el-color-primary);
      }
    }

    .el-input__inner {
      color: var(--el-text-color-primary);

      &::placeholder {
        color: var(--el-text-color-placeholder);
      }
    }

    .el-input__prefix,
    .el-input__suffix {
      color: var(--el-text-color-regular);
    }
  }

  :deep(.el-button) {
    &.el-button--info {
      color: var(--el-text-color-regular);
      background: transparent;
      border: none;

      &:hover {
        background: var(--el-fill-color-light);
      }
    }
  }
}
</style>
