<template>
  <div>
    <div class="flex-center loginPage">
      <div class="left" v-if="enableRegister">
        <div class="login-box">
          <AccountTop title="注册" />
          <div class="input-form">
            <el-form :model="data" class="form">
              <el-tabs v-model="activeName">
                <el-tab-pane label="手机注册" name="mobile" v-if="enableMobile">
                  <el-form-item>
                    <div class="form-title">手机号码</div>
                    <el-input
                      placeholder="请输入手机号码"
                      size="large"
                      v-model="data.mobile"
                      maxlength="11"
                      autocomplete="off"
                    >
                    </el-input>
                  </el-form-item>
                  <el-form-item>
                    <div class="form-title">验证码</div>
                    <div class="flex w100">
                      <el-input
                        placeholder="请输入验证码"
                        size="large"
                        maxlength="30"
                        class="code-input"
                        v-model="data.code"
                        autocomplete="off"
                      >
                      </el-input>

                      <send-msg size="large" :receiver="data.mobile" type="mobile" />
                    </div>
                  </el-form-item>
                </el-tab-pane>
                <el-tab-pane label="邮箱注册" name="email" v-if="enableEmail">
                  <el-form-item class="block">
                    <div class="form-title">邮箱</div>
                    <el-input
                      placeholder="请输入邮箱地址"
                      size="large"
                      v-model="data.email"
                      autocomplete="off"
                    >
                    </el-input>
                  </el-form-item>
                  <el-form-item class="block">
                    <div class="form-title">验证码</div>
                    <div class="flex w100">
                      <el-input
                        placeholder="请输入验证码"
                        size="large"
                        maxlength="30"
                        class="code-input"
                        v-model="data.code"
                        autocomplete="off"
                      >
                      </el-input>

                      <send-msg size="large" :receiver="data.email" type="email" />
                    </div>
                  </el-form-item>
                </el-tab-pane>
                <el-tab-pane label="用户名注册" name="username" v-if="enableUser">
                  <el-form-item class="block">
                    <div class="form-title">用户名</div>

                    <el-input
                      placeholder="请输入用户名"
                      size="large"
                      v-model="data.username"
                      autocomplete="off"
                    >
                    </el-input>
                  </el-form-item>
                </el-tab-pane>
              </el-tabs>

              <el-form-item class="block">
                <div class="form-title">密码</div>

                <el-input
                  placeholder="请输入密码(8-16位)"
                  maxlength="16"
                  size="large"
                  v-model="data.password"
                  show-password
                  autocomplete="off"
                >
                </el-input>
              </el-form-item>

              <el-form-item class="block">
                <div class="form-title">重复密码</div>

                <el-input
                  placeholder="请再次输入密码(8-16位)"
                  size="large"
                  maxlength="16"
                  v-model="data.repass"
                  show-password
                  autocomplete="off"
                >
                </el-input>
              </el-form-item>

              <el-form-item class="block">
                <div class="form-title">邀请码</div>

                <el-input
                  placeholder="请输入邀请码(可选)"
                  size="large"
                  v-model="data.invite_code"
                  autocomplete="off"
                >
                </el-input>
              </el-form-item>

              <el-form-item
                label=""
                prop="agreement"
                :class="{ 'agreement-error': agreementError }"
              >
                <div class="agreement-box" :class="{ shake: isShaking }">
                  <el-checkbox v-model="data.agreement" @change="handleAgreementChange">
                    我已阅读并同意
                    <span class="agreement-link" @click.stop.prevent="openAgreement"
                      >《用户协议》</span
                    >
                    和
                    <span class="agreement-link" @click.stop.prevent="openPrivacy"
                      >《隐私政策》</span
                    >
                  </el-checkbox>
                </div>
              </el-form-item>

              <el-row class="btn-row" :gutter="20">
                <el-col :span="24">
                  <el-button class="login-btn" type="primary" size="large" @click="submitRegister"
                    >注册</el-button
                  >
                </el-col>
              </el-row>
            </el-form>
          </div>
        </div>
      </div>
      <div class="tip-result left" v-else>
        <el-result icon="error" title="注册功能已关闭">
          <template #sub-title>
            <p>抱歉，系统已关闭注册功能，请联系管理员添加账号！</p>
            <div class="wechat-card">
              <el-image :src="wxImg" />
            </div>

            <div class="mt-3">
              <el-button type="primary" @click="router.push('/')"
                ><i class="iconfont icon-home mr-1"></i> 返回首页</el-button
              >
            </div>
          </template>
        </el-result>
      </div>
      <captcha v-if="enableVerify" @success="doSubmitRegister" ref="captchaRef" />
      <account-bg />
    </div>
  </div>
</template>

<script setup>
import AccountBg from '@/components/AccountBg.vue'
import AccountTop from '@/components/AccountTop.vue'
import { ref } from 'vue'

import { httpGet, httpPost } from '@/utils/http'
import { ElMessage, ElMessageBox } from 'element-plus'
import MarkdownIt from 'markdown-it'
import { useRouter } from 'vue-router'

import Captcha from '@/components/Captcha.vue'
import SendMsg from '@/components/SendMsg.vue'
import { getLicenseInfo, getSystemInfo } from '@/store/cache'
import { setUserToken } from '@/store/session'
import { showMessageError, showMessageOK } from '@/utils/dialog'
import { arrayContains, isMobile } from '@/utils/libs'
import { validateEmail, validateMobile } from '@/utils/validate'

const router = useRouter()
const title = ref('')
const logo = ref('')
const data = ref({
  username: '',
  mobile: '',
  email: '',
  password: '',
  code: '',
  repass: '',
  invite_code: router.currentRoute.value.query['invite_code'],
  agreement: false,
})

const enableMobile = ref(false)
const enableEmail = ref(false)
const enableUser = ref(false)
const enableRegister = ref(true)
const activeName = ref('mobile')
const wxImg = ref('/images/wx.png')
const licenseConfig = ref({})
const enableVerify = ref(false)
const captchaRef = ref(null)
const agreementError = ref(false)
const isShaking = ref(false)

// 初始化markdown解析器
const md = new MarkdownIt({
  html: true,
  linkify: true,
  typographer: true,
})

const agreementContent = ref('')
const privacyContent = ref('')

// 记录邀请码点击次数
if (data.value.invite_code) {
  httpGet('/api/invite/hits', { code: data.value.invite_code })
}

getSystemInfo()
  .then((res) => {
    if (res.data) {
      title.value = res.data.title
      logo.value = res.data.logo
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
      enableVerify.value = res.data['enabled_verify']
    }
  })
  .catch((e) => {
    ElMessage.error('获取系统配置失败：' + e.message)
  })

// 获取用户协议
httpGet('/api/config/get?key=agreement')
  .then((res) => {
    if (res.data && res.data.content) {
      agreementContent.value = res.data.content
    } else {
      agreementContent.value =
        '# 用户协议\n\n用户在使用本服务前应当阅读并同意本协议。本协议内容包括协议正文及所有本平台已经发布的或将来可能发布的各类规则。所有规则为本协议不可分割的组成部分，与协议正文具有同等法律效力。'
    }
  })
  .catch((e) => {
    console.warn(e)
    agreementContent.value =
      '# 用户协议\n\n用户在使用本服务前应当阅读并同意本协议。本协议内容包括协议正文及所有本平台已经发布的或将来可能发布的各类规则。所有规则为本协议不可分割的组成部分，与协议正文具有同等法律效力。'
  })

// 获取隐私政策
httpGet('/api/config/get?key=privacy')
  .then((res) => {
    if (res.data && res.data.content) {
      privacyContent.value = res.data.content
    } else {
      privacyContent.value =
        '# 隐私政策\n\n我们非常重视用户的隐私和个人信息保护。您在使用我们的产品与服务时，我们可能会收集和使用您的相关信息。我们希望通过本《隐私政策》向您说明我们在收集和使用您相关信息时对应的处理规则。'
    }
  })
  .catch((e) => {
    console.warn(e)
    privacyContent.value =
      '# 隐私政策\n\n我们非常重视用户的隐私和个人信息保护。您在使用我们的产品与服务时，我们可能会收集和使用您的相关信息。我们希望通过本《隐私政策》向您说明我们在收集和使用您相关信息时对应的处理规则。'
  })

getLicenseInfo()
  .then((res) => {
    licenseConfig.value = res.data
  })
  .catch((e) => {
    showMessageError('获取 License 配置：' + e.message)
  })

// 注册操作
const submitRegister = () => {
  if (activeName.value === 'username' && data.value.username === '') {
    return showMessageError('请输入用户名')
  }

  if (activeName.value === 'mobile' && !validateMobile(data.value.mobile)) {
    return showMessageError('请输入合法的手机号')
  }

  if (activeName.value === 'email' && !validateEmail(data.value.email)) {
    return showMessageError('请输入合法的邮箱地址')
  }

  if (data.value.password.length < 8) {
    return showMessageError('密码的长度为8-16个字符')
  }
  if (data.value.repass !== data.value.password) {
    return showMessageError('两次输入密码不一致')
  }

  if ((activeName.value === 'mobile' || activeName.value === 'email') && data.value.code === '') {
    return showMessageError('请输入验证码')
  }

  if (!data.value.agreement) {
    agreementError.value = true
    isShaking.value = true
    setTimeout(() => {
      isShaking.value = false
    }, 500)
    showMessageError('请先阅读并同意用户协议和隐私政策')
    return
  }

  // 如果是用户名和密码登录，那么需要加载验证码
  if (enableVerify.value && activeName.value === 'username') {
    captchaRef.value.loadCaptcha()
  } else {
    doSubmitRegister({})
  }
}

const doSubmitRegister = (verifyData) => {
  data.value.key = verifyData.key
  data.value.dots = verifyData.dots
  data.value.x = verifyData.x
  data.value.reg_way = activeName.value
  httpPost('/api/user/register', data.value)
    .then((res) => {
      setUserToken(res.data.token)
      showMessageOK('注册成功，即将跳转到对话主界面...')
      if (isMobile()) {
        router.push('/mobile/index')
      } else {
        router.push('/chat')
      }
    })
    .catch((e) => {
      showMessageError('注册失败，' + e.message)
    })
}

const handleAgreementChange = () => {
  agreementError.value = !data.value.agreement
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

:deep(.back) {
  margin-bottom: 10px;
}

:deep(.orline) {
  margin-bottom: 10px;
}

.wechat-card {
  margin-top: 20px;
}

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
