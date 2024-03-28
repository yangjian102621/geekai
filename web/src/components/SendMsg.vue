<template>
  <el-container class="captcha-box">
    <el-button type="primary" class="send-btn" :size="props.size" :disabled="!canSend" @click="loadCaptcha" plain>
      {{ btnText }}
    </el-button>

    <el-dialog
        v-model="showCaptcha"
        :close-on-click-modal="true"
        :show-close="false"
        style="width:90%;max-width: 360px;"
    >
      <slide-captcha
          v-if="isIphone()"
          :bg-img="bgImg"
          :bk-img="bkImg"
          :result="result"
          @refresh="getSlideCaptcha"
          @confirm="handleSlideConfirm"
          @hide="showCaptcha = false"/>

      <captcha-plus
          v-else
          :max-dot="maxDot"
          :image-base64="imageBase64"
          :thumb-base64="thumbBase64"
          width="300"
          @close="showCaptcha = false"
          @refresh="handleRequestCaptCode"
          @confirm="handleConfirm"
      />
    </el-dialog>
  </el-container>
</template>

<script setup>
// 发送短信验证码组件
import {ref} from "vue";
import lodash from 'lodash'
import {validateEmail, validateMobile} from "@/utils/validate";
import {ElMessage} from "element-plus";
import {httpGet, httpPost} from "@/utils/http";
import CaptchaPlus from "@/components/CaptchaPlus.vue";
import SlideCaptcha from "@/components/SlideCaptcha.vue";
import {isIphone} from "@/utils/libs";

const props = defineProps({
  receiver: String,
  size: String,
});
const btnText = ref('发送验证码')
const canSend = ref(true)
const showCaptcha = ref(false)
const maxDot = ref(5)
const imageBase64 = ref('')
const thumbBase64 = ref('')
const captKey = ref('')
const dots = ref(null)

const handleRequestCaptCode = () => {

  httpGet('/api/captcha/get').then(res => {
    const data = res.data
    imageBase64.value = data.image
    thumbBase64.value = data.thumb
    captKey.value = data.key
  }).catch(e => {
    ElMessage.error('获取人机验证数据失败：' + e.message)
  })
}

const handleConfirm = (dots) => {
  if (lodash.size(dots) <= 0) {
    return ElMessage.error('请进行人机验证再操作')
  }

  let dotArr = []
  lodash.forEach(dots, (dot) => {
    dotArr.push(dot.x, dot.y)
  })
  dots.value = dotArr.join(',')
  httpPost('/api/captcha/check', {
    dots: dots.value,
    key: captKey.value
  }).then(() => {
    // ElMessage.success('人机验证成功')
    showCaptcha.value = false
    sendMsg()
  }).catch(() => {
    ElMessage.error('人机验证失败')
    handleRequestCaptCode()
  })
}

const loadCaptcha = () => {
  if (!validateMobile(props.receiver) && !validateEmail(props.receiver)) {
    return ElMessage.error("请输入合法的手机号/邮箱地址")
  }

  showCaptcha.value = true
  // iphone 手机用滑动验证码
  if (isIphone()) {
    getSlideCaptcha()
  } else {
    handleRequestCaptCode()
  }
}

const sendMsg = () => {
  if (!canSend.value) {
    return
  }

  canSend.value = false
  httpPost('/api/sms/code', {receiver: props.receiver, key: captKey.value, dots: dots.value}).then(() => {
    ElMessage.success('验证码发送成功')
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
    ElMessage.error('验证码发送失败：' + e.message)
  })
}

// 滑动验证码
const bgImg = ref('')
const bkImg = ref('')
const result = ref(0)

const getSlideCaptcha = () => {
  result.value = 0
  httpGet("/api/captcha/slide/get").then(res => {
    bkImg.value = res.data.bkImg
    bgImg.value = res.data.bgImg
    captKey.value = res.data.key
  }).catch(e => {
    ElMessage.error('获取人机验证数据失败：' + e.message)
  })
}

const handleSlideConfirm = (x) => {
  httpPost("/api/captcha/slide/check", {
    key: captKey.value,
    x: x
  }).then(() => {
    result.value = 1
    showCaptcha.value = false
    sendMsg()
  }).catch(() => {
    result.value = 2
  })
}
</script>

<style lang="stylus">

.captcha-box {
  .send-btn {
    width: 100%;
  }

  .el-dialog {
    .el-dialog__header {
      padding: 0;
    }

    .el-dialog__body {
      //padding 0
    }
  }
}
</style>