<template>
  <el-container class="send-verify-code">
    <el-button type="primary" class="send-btn" :size="props.size" :disabled="!canSend" @click="showCaptcha" plain>
      {{ btnText }}
    </el-button>

    <captcha @success="sendMsg" ref="captchaRef"/>
  </el-container>
</template>

<script setup>
// 发送短信验证码组件
import {ref} from "vue";
import {validateEmail, validateMobile} from "@/utils/validate";
import {httpPost} from "@/utils/http";
import {showMessageError, showMessageOK} from "@/utils/dialog";
import Captcha from "@/components/Captcha.vue";

// eslint-disable-next-line no-undef
const props = defineProps({
  receiver: String,
  size: String,
});
const btnText = ref('发送验证码')
const canSend = ref(true)
const captchaRef = ref(null)

const showCaptcha = () => {
  if (!validateMobile(props.receiver) && !validateEmail(props.receiver)) {
    return showMessageError("请输入合法的手机号/邮箱地址")
  }
  captchaRef.value.loadCaptcha()
}

const sendMsg = (data) => {
  if (!canSend.value) {
    return
  }

  canSend.value = false
  httpPost('/api/sms/code', {receiver: props.receiver, key: data.key, dots: data.dots, x:data.x}).then(() => {
    showMessageOK('验证码发送成功')
    let time = 120
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