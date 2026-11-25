<template>
  <div class="login-dialog w-full">
    <div v-if="!configReady" class="login-loading">
      <el-icon class="is-loading"><Loading /></el-icon>
      <span>正在加载登录方式...</span>
    </div>
    <div class="login-box" v-else-if="login">
      <div v-if="hasLocalLogin" class="primary-login-methods" role="group" aria-label="主要登录方式">
        <button
          v-if="mobileLoginEnabled"
          type="button"
          :class="{ active: loginActiveName === 'mobile' }"
          :aria-pressed="loginActiveName === 'mobile'"
          @click="selectedLoginMethod = 'mobile'"
        >
          手机号登录
        </button>
        <button
          v-if="accountLoginEnabled"
          type="button"
          :class="{ active: loginActiveName === 'account' }"
          :aria-pressed="loginActiveName === 'account'"
          @click="selectedLoginMethod = 'account'"
        >
          账号登录
        </button>
      </div>

      <div>
        <!-- 账号密码登录 -->
        <div v-if="loginActiveName === 'account'">
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
                  class="login-submit"
                  :disabled="!accountAgreeChecked"
                  @click="submitLogin"
                  type="button"
                >
                  {{ loading ? '登录中...' : '登 录' }}
                </button>
              </el-col>
            </el-row>

            <div class="flex items-center justify-between pt-1 text-sm" style="color: var(--login-text-color)">
              <div v-if="enableRegister" class="flex items-center gap-1">
                <span>没有账户？</span>
                <el-button
                  size="small"
                  class="px-2 py-1 transition-colors duration-200"
                  style="color: var(--login-link-color)"
                  @click="login = false"
                  @mouseenter="$event.target.style.background = 'var(--login-link-hover-bg)'"
                  @mouseleave="$event.target.style.background = 'transparent'"
                  >注册</el-button
                >
              </div>
              <el-button
                type="info"
                class="forget ml-auto"
                size="small"
                @click="showResetPass = true"
                >忘记密码？</el-button
              >
            </div>

            <div class="text-sm leading-6" style="color: var(--login-text-color)">
              <el-checkbox v-model="accountAgreeChecked">我已阅读并同意</el-checkbox>
              <a href="javascript:void(0)" class="login-link whitespace-nowrap" @click="openAgreement">《用户协议》</a>
              和
              <a href="javascript:void(0)" class="login-link whitespace-nowrap" @click="openPrivacy">《隐私政策》</a>
            </div>
          </el-form>
        </div>

        <div v-if="loginActiveName === 'mobile'">
          <el-form
            :model="mobileLoginData"
            class="form space-y-5"
            @submit.prevent="submitMobileLogin"
          >
            <el-input
              v-model="mobileLoginData.mobile"
              placeholder="请输入手机号"
              aria-label="登录手机号"
              size="large"
              maxlength="11"
              inputmode="tel"
              autocomplete="tel"
            >
              <template #prefix>
                <el-icon><Iphone /></el-icon>
              </template>
            </el-input>
            <div class="flex gap-2.5">
              <el-input
                class="flex-1 min-w-0"
                v-model="mobileLoginData.code"
                placeholder="短信验证码"
                aria-label="登录短信验证码"
                size="large"
                maxlength="6"
                inputmode="numeric"
                autocomplete="one-time-code"
              >
                <template #prefix>
                  <el-icon><Checked /></el-icon>
                </template>
              </el-input>
              <send-msg
                class="w-[120px] !flex-none"
                size="large"
                :receiver="mobileLoginData.mobile.trim()"
                purpose="login"
              />
            </div>
            <div class="text-sm leading-6" style="color: var(--login-text-color)">
              <el-checkbox v-model="mobileAgreeChecked">我已阅读并同意</el-checkbox>
              <a href="javascript:void(0)" class="login-link whitespace-nowrap" @click="openAgreement">《用户协议》</a>
              和
              <a href="javascript:void(0)" class="login-link whitespace-nowrap" @click="openPrivacy">《隐私政策》</a>
            </div>
            <button
              class="login-submit"
              :disabled="loading"
              type="submit"
            >
              {{ mobileLoginButtonText }}
            </button>
            <div v-if="enableRegister && !accountLoginEnabled" class="text-center text-sm">
              <button type="button" class="login-link" @click="login = false">注册</button>
            </div>
          </el-form>
        </div>

        <!-- 微信登录 -->
        <div v-if="loginActiveName === 'wechat'">
          <div class="wechat-login pt-3">
            <div class="qr-code-container">
              <div class="qr-code-wrapper w-[200px] h-[200px] mx-auto" v-loading="qrcodeLoading">
                <img :src="wechatLoginQRCode" class="qr-frame" v-if="wechatLoginQRCode" />
                <div
                  v-else
                  class="w-[200px] h-[200px] flex justify-center items-center text-green-600"
                >
                  <i class="iconfont icon-wechat !text-3xl"></i>
                </div>
                <!-- 二维码过期蒙版 -->
                <div v-if="qrcodeExpired" class="qr-expired-mask">
                  <div class="expired-content">
                    <i class="iconfont icon-refresh-ccw expired-icon"></i>
                    <p class="expired-text">二维码已过期</p>
                    <button
                      @click="getWxLoginURL"
                      class="bg-gray-200 text-gray-600 px-2.5 py-1 rounded-md hover:bg-gray-300"
                    >
                      <i class="iconfont icon-refresh text-lg"></i>
                    </button>
                  </div>
                </div>
              </div>
              <p class="text-center mt-4 text-gray-600 dark:text-gray-400">
                请使用微信扫描二维码登录
              </p>
            </div>
          </div>
        </div>

        <div v-if="loginActiveName === 'tokengine'" class="py-4 text-center">
          <el-icon class="!text-4xl mb-4" style="color: var(--login-link-color)"><Key /></el-icon>
          <p class="mb-6 text-sm" style="color: var(--login-text-color)">使用 Tokengine 账号登录</p>
          <button
            type="button"
            class="login-submit"
            :disabled="!tokengineEnabled"
            @click="startTokengineLogin"
          >
            Tokengine
          </button>
          <p v-if="!tokengineEnabled" class="mt-3 text-sm" style="color: var(--login-text-color)">Tokengine 登录暂未开放</p>
        </div>

        <p v-if="!hasLocalLogin && !hasThirdPartyLogin" class="py-4 text-center text-sm">
          暂无可用的登录方式，请联系管理员。
        </p>
      </div>

      <div v-if="hasLocalLogin && showThirdPartySwitcher" class="other-login-divider"><span>其他登录方式</span></div>
      <div v-if="showThirdPartySwitcher" class="secondary-login-methods" role="group" aria-label="第三方登录方式">
        <button
          v-if="showTokengineLogin"
          type="button"
          aria-label="Tokengine"
          :aria-pressed="loginActiveName === 'tokengine'"
          :class="{ active: loginActiveName === 'tokengine' }"
          @click="selectedLoginMethod = 'tokengine'"
        >
          <span class="secondary-login-icon"><el-icon><Key /></el-icon></span>
          <span>Tokengine</span>
        </button>
        <button
          v-if="showWechatLogin"
          type="button"
          aria-label="微信登录"
          :aria-pressed="loginActiveName === 'wechat'"
          :class="{ active: loginActiveName === 'wechat' }"
          @click="selectedLoginMethod = 'wechat'"
        >
          <span class="secondary-login-icon text-green-600"><i class="iconfont icon-wechat" aria-hidden="true"></i></span>
          <span>微信登录</span>
        </button>
      </div>
    </div>

    <div class="register-box w-full" v-else>
      <el-form :model="registerData" class="form space-y-5" v-if="enableRegister">
        <div class="block">
          <el-input
            placeholder="用户名（必填）"
            aria-label="用户名"
            size="large"
            v-model="registerData.username"
            maxlength="30"
            autocomplete="username"
          >
            <template #prefix>
              <el-icon><User /></el-icon>
            </template>
          </el-input>
        </div>

        <div class="block">
          <el-input
            placeholder="密码（选填，8-16位）"
            aria-label="密码（选填）"
            maxlength="16"
            size="large"
            v-model="registerData.password"
            show-password
            autocomplete="new-password"
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
            placeholder="确认密码"
            aria-label="确认密码"
            size="large"
            maxlength="16"
            v-model="registerData.repass"
            show-password
            autocomplete="new-password"
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
            placeholder="手机号（必填）"
            aria-label="手机号"
            size="large"
            v-model="registerData.mobile"
            maxlength="11"
            inputmode="tel"
            autocomplete="tel"
          >
            <template #prefix>
              <el-icon><Iphone /></el-icon>
            </template>
          </el-input>
        </div>

        <div class="block">
          <div class="flex gap-2.5">
            <el-input
              class="flex-1 min-w-0"
              placeholder="短信验证码"
              aria-label="短信验证码"
              size="large"
              maxlength="6"
              inputmode="numeric"
              v-model="registerData.code"
              autocomplete="one-time-code"
            >
              <template #prefix>
                <el-icon><Checked /></el-icon>
              </template>
            </el-input>
            <send-msg
              class="w-[120px] !flex-none"
              size="large"
              :receiver="registerData.mobile.trim()"
              type="mobile"
            />
          </div>
        </div>

        <div class="block">
          <el-input
            placeholder="邀请码(可选)"
            size="large"
            v-model="registerData.invite_code"
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
            <a href="javascript:void(0)" class="login-link" @click="openAgreement"
              >《用户协议》</a
            >
            和
            <a href="javascript:void(0)" class="login-link" @click="openPrivacy">《隐私政策》</a>
          </el-checkbox>
        </div>

        <div class="w-full">
          <button
            class="login-submit"
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
            <el-result icon="error" title="表单注册暂不可用">
              <template #sub-title>
                <p>如需开通账号，请联系管理员。</p>
              </template>
              <template #extra>
                <el-button @click="login = true">返回登录</el-button>
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

    <el-dialog v-model="showAgreement" title="用户协议" class="login-policy-dialog" :append-to-body="true">
      <div class="prose" v-html="agreementHtml"></div>
    </el-dialog>

    <el-dialog v-model="showPrivacy" title="隐私政策" class="login-policy-dialog" :append-to-body="true">
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
import { safeRedirect } from '@/utils/redirect'
import { validateMobile } from '@/utils/validate'
import { Checked, Iphone, Key, Loading, Lock, Message, User } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { marked } from 'marked'
import QRCode from 'qrcode'
import { isWechat } from '@/utils/libs'
import { computed, onMounted, onUnmounted, ref, watch } from 'vue'

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
const login = ref(props.active === 'login')
const registerWays = ref([])
const mobileLoginEnabled = computed(() => registerWays.value.includes('mobile'))
const accountLoginEnabled = computed(() => registerWays.value.includes('username'))
const hasLocalLogin = computed(() => mobileLoginEnabled.value || accountLoginEnabled.value)
const showTokengineLogin = ref(false)
const showWechatLogin = ref(false)
const hasThirdPartyLogin = computed(() => showTokengineLogin.value || showWechatLogin.value)
const enableRegister = ref(false)
const configReady = ref(false)
const selectedLoginMethod = ref('')
const loginActiveName = computed(() => {
  if (selectedLoginMethod.value) return selectedLoginMethod.value
  if (mobileLoginEnabled.value) return 'mobile'
  if (accountLoginEnabled.value) return 'account'
  if (showTokengineLogin.value) return 'tokengine'
  return showWechatLogin.value ? 'wechat' : ''
})
const showThirdPartySwitcher = computed(() => {
  const onLocalLogin = loginActiveName.value === 'mobile' || loginActiveName.value === 'account'
  return (onLocalLogin && hasThirdPartyLogin.value) || (showTokengineLogin.value && showWechatLogin.value)
})
const data = ref({
  username: import.meta.env.VITE_USER,
  password: import.meta.env.VITE_PASS,
})
const registerData = ref({
  username: '',
  mobile: '',
  password: '',
  repass: '',
  code: '',
  invite_code: props.inviteCode,
})
const mobileLoginData = ref({ mobile: '', code: '' })
const mobileAgreeChecked = ref(false)
const accountAgreeChecked = ref(false)
const tokengineEnabled = ref(false)

const startTokengineLogin = () => {
  const query = new URLSearchParams({
    invite_code: props.inviteCode,
    redirect: safeRedirect(new URLSearchParams(window.location.search).get('redirect')),
  })
  window.location.assign('/api/user/login/tokengine/start?' + query)
}

// 微信登录相关变量
const wechatLoginQRCode = ref('')
const wechatLoginState = ref('')
const qrcodeLoading = ref(false)
const pollingTimer = ref(null)
const qrcodeExpired = ref(false)
const qrcodeTimer = ref(null)

const wxImg = ref('/images/wx.png')
const captchaRef = ref(null)
// eslint-disable-next-line no-undef
const emits = defineEmits(['hide', 'success', 'changeActive'])
const action = ref('login')
const enableCaptcha = ref(false)
const captchaType = ref('')
const showResetPass = ref(false)
const store = useSharedStore()
const loading = ref(false)
const mobileLoginButtonText = computed(() => {
  if (loading.value) return '登录中...'
  return enableRegister.value ? '登录 / 注册' : '登录'
})
const agreeChecked = ref(false)
const showAgreement = ref(false)
const showPrivacy = ref(false)
const agreementHtml = ref('')
const privacyHtml = ref('')

const handleLoginSuccess = (res) => {
  setUserToken(res.data.token)
  store.setIsLogin(true)
  ElMessage.success('登录成功！')
  emits('hide')
  emits('success')
}

const clearWechatTimers = () => {
  clearInterval(pollingTimer.value)
  clearTimeout(qrcodeTimer.value)
}

watch(
  () => login.value,
  (newValue) => {
    emits('changeActive', newValue)
  }
)

const loginWithWechatBrowser = () => {
  const urlParams = new URLSearchParams(window.location.search)
  if (!isWechat() || urlParams.has('tokengine')) return

  const code = urlParams.get('code')
  if (code) {
    httpPost('/api/user/wxAuthLogin', { code })
      .then(handleLoginSuccess)
      .catch((e) => {
        ElMessage.error('登录失败，' + e.message)
      })
      .finally(() => {
        loading.value = false
      })
    return
  }

  httpGet('/api/config/get?key=wx_gzh')
    .then((res) => {
      const redirectUri = encodeURIComponent(`${location.protocol}//${location.host}/login`)
      window.location.href = `https://open.weixin.qq.com/connect/oauth2/authorize?appid=${res.data.app_id}&redirect_uri=${redirectUri}&response_type=code&scope=snsapi_userinfo&state=STATE#wechat_redirect`
    })
    .catch((e) => {
      console.error('获取微信公众号配置失败：', e.message)
    })
}

onMounted(() => {
  getSystemInfo()
    .then((res) => {
      const config = res.data
      registerWays.value = Array.isArray(config?.register_ways) ? config.register_ways : []
      enableRegister.value = config?.enabled_register === true && hasLocalLogin.value
      showTokengineLogin.value = config?.show_tokengine_login !== false
      showWechatLogin.value = config?.show_wechat_login !== false
      if (config?.wechat_card_url) {
        wxImg.value = config.wechat_card_url
      }
      if (showTokengineLogin.value) {
        httpGet('/api/user/login/tokengine/config')
          .then((result) => { tokengineEnabled.value = result.data?.enabled === true })
          .catch(() => { tokengineEnabled.value = false })
      }
      if (showWechatLogin.value) loginWithWechatBrowser()
    })
    .catch((e) => {
      enableRegister.value = false
      ElMessage.error('获取系统配置失败：' + e.message)
    })
    .finally(() => {
      configReady.value = true
    })

  // 获取验证码配置
  httpGet('/api/captcha/config').then((res) => {
    enableCaptcha.value = res.data['enabled']
    captchaType.value = res.data['type']
  })
})

// 监听登录标签页切换
watch(loginActiveName, (newValue) => {
  if (newValue === 'wechat') {
    getWxLoginURL()
  } else {
    clearWechatTimers()
  }
})

const submit = (verifyData) => {
  if (action.value === 'login') {
    doLogin(verifyData)
  } else if (action.value === 'register') {
    doRegister(verifyData)
  } else if (action.value === 'mobileLogin') {
    doLogin(verifyData, true)
  }
}

// 获取微信登录 URL
const getWxLoginURL = () => {
  wechatLoginQRCode.value = ''
  qrcodeLoading.value = true
  qrcodeExpired.value = false

  clearWechatTimers()

  httpGet('/api/user/login/qrcode')
    .then((res) => {
      // 生成二维码
      QRCode.toDataURL(res.data.url, { width: 200, height: 200, margin: 2 }, (error, url) => {
        if (error) {
          console.error(error)
        } else {
          wechatLoginQRCode.value = url
        }
      })
      wechatLoginState.value = res.data.state
      // 开始轮询状态
      startPolling()

      // 设置1分钟后二维码过期
      qrcodeTimer.value = setTimeout(() => {
        qrcodeExpired.value = true
        clearWechatTimers()
      }, 60 * 1000) // 1分钟过期
    })
    .catch((e) => {
      ElMessage.error('获取微信登录 URL 失败，' + e.message)
    })
    .finally(() => {
      qrcodeLoading.value = false
    })
}

// 开始轮询
const startPolling = () => {
  clearInterval(pollingTimer.value)
  pollingTimer.value = setInterval(checkLoginStatus, 1000)
}

// 检查登录状态
const checkLoginStatus = () => {
  if (!wechatLoginState.value) return

  httpGet(`/api/user/login/status?state=${wechatLoginState.value}`)
    .then((res) => {
      const status = res.data.status
      if (status === 'pending') return
      clearWechatTimers()

      switch (status) {
        case 'success':
          // 登录成功
          handleLoginSuccess(res)
          break

        case 'expired':
          // 二维码过期
          qrcodeExpired.value = true
          break

        default:
          // 其他错误情况
          ElMessage.error('登录失败，请重试')
          break
      }
    })
    .catch(() => {
      // 发生错误时显示过期状态
      clearWechatTimers()
      qrcodeExpired.value = true
    })
}

// 登录操作
const submitLogin = () => {
  if (!data.value.username) {
    return ElMessage.error('请输入用户名')
  }
  if (!data.value.password) {
    return ElMessage.error('请输入密码')
  }
  if (!accountAgreeChecked.value) {
    return ElMessage.error('请先阅读并同意《用户协议》和《隐私政策》')
  }
  if (enableCaptcha.value) {
    captchaRef.value.loadCaptcha()
    action.value = 'login'
  } else {
    doLogin({})
  }
}

const submitMobileLogin = () => {
  if (loading.value) return
  mobileLoginData.value.mobile = mobileLoginData.value.mobile.trim()
  mobileLoginData.value.code = mobileLoginData.value.code.trim()
  if (!validateMobile(mobileLoginData.value.mobile)) {
    return ElMessage.error('请输入合法的手机号')
  }
  if (!/^\d{6}$/.test(mobileLoginData.value.code)) {
    return ElMessage.error('请输入6位短信验证码')
  }
  if (!mobileAgreeChecked.value) {
    return ElMessage.error('请先阅读并同意《用户协议》和《隐私政策》')
  }
  if (enableCaptcha.value) {
    action.value = 'mobileLogin'
    captchaRef.value.loadCaptcha()
  } else {
    doLogin({}, true)
  }
}

const doLogin = (verifyData, mobileLogin = false) => {
  const payload = mobileLogin
    ? { ...mobileLoginData.value, invite_code: props.inviteCode }
    : data.value
  loading.value = true
  httpPost(mobileLogin ? '/api/user/login/mobile' : '/api/user/login', {
    ...payload,
    key: verifyData.key,
    dots: verifyData.dots,
    x: verifyData.x,
  })
    .then(handleLoginSuccess)
    .catch((e) => {
      ElMessage.error('登录失败，' + e.message)
    })
    .finally(() => {
      loading.value = false
    })
}

// 注册操作
const submitRegister = () => {
  registerData.value.username = registerData.value.username.trim()
  registerData.value.mobile = registerData.value.mobile.trim()
  registerData.value.code = registerData.value.code.trim()
  if (!registerData.value.username) {
    return ElMessage.error('请输入用户名')
  }

  if (registerData.value.password && registerData.value.password.length < 8) {
    return ElMessage.error('密码的长度为8-16个字符')
  }
  if (registerData.value.repass !== registerData.value.password) {
    return ElMessage.error('两次输入密码不一致')
  }

  if (!validateMobile(registerData.value.mobile)) {
    return ElMessage.error('请输入合法的手机号')
  }

  if (!/^\d{6}$/.test(registerData.value.code)) {
    return ElMessage.error('请输入6位短信验证码')
  }

  if (!agreeChecked.value) {
    return ElMessage.error('请先阅读并同意《用户协议》和《隐私政策》')
  }
  if (enableCaptcha.value) {
    captchaRef.value.loadCaptcha()
    action.value = 'register'
  } else {
    doRegister({})
  }
}

const doRegister = (verifyData) => {
  loading.value = true
  httpPost('/api/user/register', {
    ...registerData.value,
    reg_way: 'mobile',
    key: verifyData.key,
    dots: verifyData.dots,
    x: verifyData.x,
  })
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

// 组件卸载时清除定时器
onUnmounted(clearWechatTimers)
</script>

<style lang="scss">
.login-dialog .send-btn,
.reset-pass-dialog .send-btn {
  --el-button-bg-color: var(--login-soft-bg);
  --el-button-border-color: var(--login-card-border);
  --el-button-text-color: var(--login-link-color);
  --el-button-hover-bg-color: var(--login-link-hover-bg);
  --el-button-hover-border-color: var(--login-link-color);
  --el-button-hover-text-color: var(--login-link-color);
  --el-button-active-bg-color: var(--login-link-hover-bg);
  --el-button-active-border-color: var(--login-link-color);
  --el-button-active-text-color: var(--login-link-color);
  --el-button-disabled-bg-color: var(--login-soft-bg);
  --el-button-disabled-border-color: var(--login-card-border);
  --el-button-disabled-text-color: var(--login-subtitle-color);
  height: 44px;
  padding: 0 12px;
  border-radius: 8px;
  font-size: 13px;
}

.reset-pass-dialog .send-btn {
  height: 32px;
}

.login-dialog {
  .login-loading {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 12px;
    padding: 48px 0;
    color: var(--login-text-color);
    font-size: 14px;

    .el-icon {
      font-size: 24px;
    }
  }

  .login-submit {
    width: 100%;
    height: 48px;
    border: 1px solid var(--login-primary-color);
    border-radius: 8px;
    background: var(--login-primary-color);
    color: #fff;
    font-size: 15px;
    font-weight: 500;
    transition: background-color 160ms ease, border-color 160ms ease;

    &:hover:not(:disabled) {
      background: var(--login-primary-hover);
      border-color: var(--login-primary-hover);
    }

    &:disabled {
      cursor: not-allowed;
      opacity: 0.6;
    }
  }

  .login-link {
    color: var(--login-link-color);
    white-space: nowrap;

    &:hover {
      text-decoration: underline;
      text-underline-offset: 3px;
    }
  }

  .el-checkbox {
    height: auto;
    align-items: flex-start;
    white-space: normal;

    .el-checkbox__input {
      margin-top: 4px;
    }

    .el-checkbox__label {
      white-space: normal;
      line-height: 1.7;
    }
  }

  .login-submit:focus-visible,
  .login-link:focus-visible,
  .primary-login-methods button:focus-visible,
  .secondary-login-methods button:focus-visible {
    outline: 2px solid var(--login-link-color);
    outline-offset: 4px;
    border-radius: 4px;
  }

  // 微信登录样式
  .wechat-login {
    display: flex;
    justify-content: center;
    align-items: center;
    min-height: 240px;

    .qr-code-container {
      text-align: center;

      .qr-code-wrapper {
        display: inline-block;
        border: 1px solid var(--el-border-color);
        border-radius: 8px;
        overflow: hidden;
        position: relative;

        .qr-frame {
          display: block;
          width: 100%;
          height: 100%;
        }

        .qr-expired-mask {
          position: absolute;
          top: 0;
          left: 0;
          right: 0;
          bottom: 0;
          background: rgba(0, 0, 0, 0.7);
          display: flex;
          align-items: center;
          justify-content: center;
          border-radius: 8px;

          .expired-content {
            text-align: center;
            color: white;

            .expired-icon {
              font-size: 48px;
              color: #f56565;
              margin-bottom: 12px;
              display: block;
            }

            .expired-text {
              font-size: 16px;
              margin: 0 0 16px 0;
              font-weight: 500;
            }
          }
        }
      }
    }
  }

  .primary-login-methods {
    display: flex;
    gap: 24px;
    margin-bottom: 24px;
    border-bottom: 1px solid var(--login-card-border);

    button {
      flex: 1;
      padding: 0 0 12px;
      border-bottom: 2px solid transparent;
      color: var(--login-text-color);
      font-size: 16px;
      font-weight: 500;

      &.active {
        color: var(--login-link-color);
        border-bottom-color: currentColor;
      }
    }
  }

  .other-login-divider {
    display: flex;
    align-items: center;
    gap: 16px;
    margin: 28px 0 20px;
    color: var(--login-text-color);
    font-size: 13px;

    &::before,
    &::after {
      content: '';
      flex: 1;
      height: 1px;
      background: var(--login-card-border);
    }
  }

  .secondary-login-methods {
    display: flex;
    justify-content: center;
    gap: 48px;

    button {
      display: flex;
      flex-direction: column;
      align-items: center;
      gap: 8px;
      color: var(--login-text-color);
      font-size: 12px;

      &:hover,
      &.active {
        color: var(--login-link-color);

        .secondary-login-icon {
          border-color: var(--login-link-color);
        }
      }
    }
  }

  .secondary-login-icon {
    display: flex;
    justify-content: center;
    align-items: center;
    width: 44px;
    height: 44px;
    border: 1px solid var(--login-card-border);
    border-radius: 50%;
    font-size: 22px;

    .iconfont {
      font-size: 24px;
    }
  }

  .el-input {
    .el-input__wrapper {
      min-height: 44px;
      border-radius: 8px;
      background: var(--login-input-bg);
      box-shadow: 0 0 0 1px var(--login-card-border) inset;
      transition: box-shadow 160ms ease;

      &:hover {
        box-shadow: 0 0 0 1px var(--login-placeholder-color) inset;
      }

      &.is-focus {
        box-shadow: 0 0 0 1px var(--login-link-color) inset;
      }
    }

    .el-input__inner {
      color: var(--login-title-color);

      &::placeholder {
        color: var(--login-placeholder-color);
      }
    }

    .el-input__prefix,
    .el-input__suffix {
      color: var(--login-subtitle-color);
    }
  }

  .el-button {
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

// 响应式设计
@media (max-width: 576px) {
  .login-dialog {
    .wechat-login {
      .qr-code-wrapper {
        width: 240px !important;
        height: 240px !important;

        .qr-expired-mask {
          .expired-content {
            .expired-icon {
              font-size: 36px;
              margin-bottom: 8px;
            }

            .expired-text {
              font-size: 14px;
              margin: 0 0 12px 0;
            }
          }
        }
      }
    }
  }
}
</style>
