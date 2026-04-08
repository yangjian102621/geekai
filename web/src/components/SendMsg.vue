<template>
  <el-container class="send-verify-code">
    <el-button type="success" :size="props.size" :disabled="!canSend" @click="sendMsg">
      {{ btnText }}
    </el-button>

    <captcha @success="doSendMsg" ref="captchaRef" :type="captchaType" />
  </el-container>
</template>

<script setup>
// 发送短信验证码组件
import Captcha from '@/components/Captcha.vue'
import { httpGet, httpPost } from '@/utils/http'
import { validateEmail, validateMobile } from '@/utils/validate'
import { ElMessage } from 'element-plus'
import { ref } from 'vue'

// eslint-disable-next-line no-undef
const props = defineProps({
  receiver: String,
  size: String,
  type: {
    type: String,
    default: 'mobile',
  },
})
const btnText = ref('发送验证码')
const canSend = ref(true)
const captchaRef = ref(null)
const enableCaptcha = ref(false)
const captchaType = ref('')

httpGet('/api/captcha/config').then((res) => {
  enableCaptcha.value = res.data['enabled']
  captchaType.value = res.data['type']
})

const sendMsg = () => {
  if (!validateMobile(props.receiver) && props.type === 'mobile') {
    return ElMessage.error('请输入合法的手机号')
  }
  if (!validateEmail(props.receiver) && props.type === 'email') {
    return ElMessage.error('请输入合法的邮箱地址')
  }

  if (enableCaptcha.value) {
    captchaRef.value.loadCaptcha()
  } else {
    doSendMsg({})
  }
}

const doSendMsg = (data) => {
  if (!canSend.value) {
    return
  }

  canSend.value = false
  httpPost('/api/sms/code', {
    receiver: props.receiver,
    key: data.key,
    dots: data.dots,
    x: data.x,
  })
    .then(() => {
      if (props.type === 'mobile') {
        ElMessage.success('验证码发送成功')
      } else if (props.type === 'email') {
        ElMessage.success('验证码已发送至邮箱，如果长时间未收到，请检查是否在垃圾邮件中！')
      }
      let time = 60
      btnText.value = time
      const handler = setInterval(() => {
        time = time - 1
        if (time <= 0) {
          clearInterval(handler)
          btnText.value = '重新发送'
          canSend.value = true
        } else {
          btnText.value = time
        }
      }, 1000)
    })
    .catch((e) => {
      canSend.value = true
      ElMessage.error('验证码发送失败：' + e.message)
    })
}
</script>

<style lang="scss" scoped>
.send-verify-code {
  .send-btn {
    width: 100%;
  }
}
</style>
