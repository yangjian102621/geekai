<template>
  <el-container class="send-verify-code">
    <el-button type="primary" class="send-btn" :size="props.size" :disabled="!canSend" @click="sendMsg" plain>
      {{ btnText }}
    </el-button>

    <captcha @success="doSendMsg" ref="captchaRef"/>
  </el-container>
</template>

<script setup>
// 发送短信验证码组件
import {ref} from "vue";
import {validateEmail, validateMobile} from "@/utils/validate";
import {httpPost} from "@/utils/http";
import {showMessageError, showMessageOK} from "@/utils/dialog";
import Captcha from "@/components/Captcha.vue";
import {getSystemInfo} from "@/store/cache";

// eslint-disable-next-line no-undef
const props = defineProps({
  receiver: String,
  size: String,
  type: {
    type: String,
    default: 'mobile'
  }
});
const btnText = ref('发送验证码')
const canSend = ref(true)
const captchaRef = ref(null)
const enableVerify = ref(false)

getSystemInfo().then(res => {
  enableVerify.value = res.data['enabled_verify']
})

const sendMsg = () => {
  if (!validateMobile(props.receiver) && props.type === 'mobile') {
    return showMessageError("请输入合法的手机号")
  }
  if (!validateEmail(props.receiver) && props.type === 'email') {
    return showMessageError("请输入合法的邮箱地址")
  }
  
  if (enableVerify.value) {
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
  httpPost('/api/sms/code', {receiver: props.receiver, key: data.key, dots: data.dots, x:data.x}).then(() => {
    showMessageOK('验证码发送成功')
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
  }).catch(e => {
    canSend.value = true
    showMessageError('验证码发送失败：' + e.message)
  })
}
</script>

<style lang="stylus" scoped>

.send-verify-code {
  .send-btn {
    width: 100%;
  }
}
</style>