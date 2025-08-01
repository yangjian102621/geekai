<template>
  <div class="flex-center loginPage">
    <div class="left">
      <div class="login-box">
        <AccountTop>
          <template #default>
            <div class="wechatLog flex-center" v-if="wechatLoginURL !== ''">
              <a :href="wechatLoginURL" @click="setRoute(router.currentRoute.value.path)">
                <i class="iconfont icon-wechat"></i>使用微信登录
              </a>
            </div>
          </template>
        </AccountTop>

        <div class="input-form">
          <el-form ref="ruleFormRef" :model="ruleForm" :rules="rules">
            <el-form-item label="" prop="username">
              <div class="form-title">账号</div>
              <el-input
                v-model="ruleForm.username"
                size="large"
                placeholder="请输入账号"
                @keyup="handleKeyup"
              />
            </el-form-item>
            <el-form-item label="" prop="password">
              <div class="flex-between w100">
                <div class="form-title">密码</div>
                <div class="form-forget text-color-primary" @click="router.push('/resetpassword')">
                  忘记密码？
                </div>
              </div>

              <el-input
                size="large"
                v-model="ruleForm.password"
                placeholder="请输入密码"
                show-password
                autocomplete="off"
                @keyup="handleKeyup"
              />
            </el-form-item>
            <el-form-item label="" prop="agreement" :class="{ 'agreement-error': agreementError }">
              <div class="agreement-box" :class="{ shake: isShaking }">
                <el-checkbox v-model="ruleForm.agreement" @change="handleAgreementChange">
                  我已阅读并同意
                  <span class="agreement-link" @click.stop.prevent="openAgreement"
                    >《用户协议》</span
                  >
                  和
                  <span class="agreement-link" @click.stop.prevent="openPrivacy">《隐私政策》</span>
                </el-checkbox>
              </div>
            </el-form-item>
            <el-form-item>
              <el-button class="login-btn" size="large" type="primary" @click="login"
                >登录</el-button
              >
            </el-form-item>
          </el-form>
        </div>
      </div>
    </div>
    <account-bg />

    <captcha v-if="enableVerify" @success="doLogin" ref="captchaRef" />
  </div>
</template>

<script setup>
import AccountBg from '@/components/AccountBg.vue'
import { checkSession, getLicenseInfo, getSystemInfo } from '@/store/cache'
import { setUserToken } from '@/store/session'
import { useSharedStore } from '@/store/sharedata'
import { setRoute } from '@/store/system'
import { showMessageError } from '@/utils/dialog'
import { httpGet, httpPost } from '@/utils/http'
import { ElMessageBox } from 'element-plus'
import MarkdownIt from 'markdown-it'
import { onMounted, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'

import AccountTop from '@/components/AccountTop.vue'
import Captcha from '@/components/Captcha.vue'

const router = useRouter()
const title = ref('')
const logo = ref('')
const licenseConfig = ref({})
const wechatLoginURL = ref('')
const enableVerify = ref(false)
const captchaRef = ref(null)
const ruleFormRef = ref(null)
const ruleForm = reactive({
  username: import.meta.env.VITE_USER,
  password: import.meta.env.VITE_PASS,
  agreement: false,
})
const rules = {
  username: [{ required: true, trigger: 'blur', message: '请输入账号' }],
  password: [{ required: true, trigger: 'blur', message: '请输入密码' }],
  agreement: [{ required: true, trigger: 'change', message: '请同意用户协议' }],
}
const agreementContent = ref('')
const privacyContent = ref('')

// 初始化markdown解析器
const md = new MarkdownIt({
  html: true,
  linkify: true,
  typographer: true,
})

onMounted(() => {
  // 检查URL中是否存在token参数
  const urlParams = new URLSearchParams(window.location.search)
  const token = urlParams.get('token')
  if (token) {
    setUserToken(token)
    store.setIsLogin(true)
    router.push('/chat')
    return
  }

  // 获取系统配置
  getSystemInfo()
    .then((res) => {
      logo.value = res.data.logo
      title.value = res.data.title
      enableVerify.value = res.data['enabled_verify']
    })
    .catch((e) => {
      showMessageError('获取系统配置失败：' + e.message)
      title.value = 'Geek-AI'
    })

  // 获取用户协议
  httpGet('/api/config/get?key=agreement')
    .then((res) => {
      if (res.data && res.data.content) {
        agreementContent.value = res.data.content
      } else {
        agreementContent.value =
          '用户在使用本服务前应当阅读并同意本协议。本协议内容包括协议正文及所有本平台已经发布的或将来可能发布的各类规则。所有规则为本协议不可分割的组成部分，与协议正文具有同等法律效力。'
      }
    })
    .catch((e) => {
      agreementContent.value =
        '用户在使用本服务前应当阅读并同意本协议。本协议内容包括协议正文及所有本平台已经发布的或将来可能发布的各类规则。所有规则为本协议不可分割的组成部分，与协议正文具有同等法律效力。'
    })

  // 获取隐私政策
  httpGet('/api/config/get?key=privacy')
    .then((res) => {
      if (res.data && res.data.content) {
        privacyContent.value = res.data.content
      } else {
        privacyContent.value =
          '我们非常重视用户的隐私和个人信息保护。您在使用我们的产品与服务时，我们可能会收集和使用您的相关信息。我们希望通过本《隐私政策》向您说明我们在收集和使用您相关信息时对应的处理规则。'
      }
    })
    .catch((e) => {
      privacyContent.value =
        '我们非常重视用户的隐私和个人信息保护。您在使用我们的产品与服务时，我们可能会收集和使用您的相关信息。我们希望通过本《隐私政策》向您说明我们在收集和使用您相关信息时对应的处理规则。'
    })

  getLicenseInfo()
    .then((res) => {
      licenseConfig.value = res.data
    })
    .catch((e) => {
      showMessageError('获取 License 配置：' + e.message)
    })

  checkSession()
    .then(() => {
      router.back()
    })
    .catch(() => {})

  const returnURL = `${location.protocol}//${location.host}/login/callback?action=login`
  httpGet('/api/user/clogin?return_url=' + returnURL)
    .then((res) => {
      wechatLoginURL.value = res.data.url
    })
    .catch((e) => {
      console.error(e)
    })
})

const handleKeyup = (e) => {
  if (e.key === 'Enter') {
    login()
  }
}

const login = async function () {
  if (!ruleForm.agreement) {
    agreementError.value = true
    isShaking.value = true
    setTimeout(() => {
      isShaking.value = false
    }, 500)
    showMessageError('请先阅读并同意用户协议')
    return
  }

  await ruleFormRef.value.validate(async (valid) => {
    if (valid) {
      if (enableVerify.value) {
        captchaRef.value.loadCaptcha()
      } else {
        doLogin({})
      }
    }
  })
}

const store = useSharedStore()
const doLogin = (verifyData) => {
  httpPost('/api/user/login', {
    username: ruleForm.username,
    password: ruleForm.password,
    key: verifyData.key,
    dots: verifyData.dots,
    x: verifyData.x,
  })
    .then((res) => {
      setUserToken(res.data.token)
      store.setIsLogin(true)
      router.back()
    })
    .catch((e) => {
      showMessageError('登录失败，' + e.message)
    })
}

const agreementError = ref(false)
const isShaking = ref(false)

const handleAgreementChange = () => {
  agreementError.value = !ruleForm.agreement
  if (agreementError.value) {
    isShaking.value = true
    setTimeout(() => {
      isShaking.value = false
    }, 500)
  }
}

const openAgreement = () => {
  // 使用弹窗显示用户协议内容，支持Markdown格式
  ElMessageBox.alert(
    `<div class="markdown-content">${md.render(agreementContent.value)}</div>`,
    '用户协议',
    {
      confirmButtonText: '我已阅读',
      dangerouslyUseHTMLString: true,
      callback: () => {},
    }
  )
}

const openPrivacy = () => {
  // 使用弹窗显示隐私政策内容，支持Markdown格式
  ElMessageBox.alert(
    `<div class="markdown-content">${md.render(privacyContent.value)}</div>`,
    '隐私政策',
    {
      confirmButtonText: '我已阅读',
      dangerouslyUseHTMLString: true,
      callback: () => {},
    }
  )
}
</script>

<style lang="scss" scoped>
@use '../assets/css/login.scss' as *;

.agreement-box {
  margin-bottom: 10px;
  transition: all 0.3s;
}

.agreement-link {
  color: var(--el-color-primary);
  cursor: pointer;
}

.agreement-error {
  .el-checkbox {
    .el-checkbox__input {
      .el-checkbox__inner {
        border-color: #f56c6c !important;
      }
    }
  }
}

.shake {
  animation: shake 0.5s cubic-bezier(0.36, 0.07, 0.19, 0.97) both;
}

@keyframes shake {
  10%,
  90% {
    transform: translate3d(-1px, 0, 0);
  }
  20%,
  80% {
    transform: translate3d(2px, 0, 0);
  }
  30%,
  50%,
  70% {
    transform: translate3d(-4px, 0, 0);
  }
  40%,
  60% {
    transform: translate3d(4px, 0, 0);
  }
}
</style>

<style>
/* 全局样式，用于Markdown内容显示 */
.markdown-content {
  text-align: left;
  max-height: 60vh;
  overflow-y: auto;
  padding: 10px;
}

.markdown-content h1 {
  font-size: 1.5em;
  margin-bottom: 15px;
}

.markdown-content h2 {
  font-size: 1.3em;
  margin: 15px 0 10px;
}

.markdown-content p {
  margin-bottom: 10px;
  line-height: 1.5;
}

.markdown-content ul,
.markdown-content ol {
  padding-left: 20px;
  margin-bottom: 10px;
}

.markdown-content blockquote {
  border-left: 4px solid #ccc;
  padding-left: 10px;
  color: #666;
  margin: 10px 0;
}

.markdown-content code {
  background-color: #f0f0f0;
  padding: 2px 4px;
  border-radius: 3px;
  font-family: monospace;
}

.markdown-content pre {
  background-color: #f0f0f0;
  padding: 10px;
  border-radius: 5px;
  overflow-x: auto;
  margin: 10px 0;
}
</style>
