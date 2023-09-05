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
      <captcha-plus
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
import {validateMobile} from "@/utils/validate";
import {ElMessage} from "element-plus";
import {httpGet, httpPost} from "@/utils/http";
import CaptchaPlus from "@/components/CaptchaPlus.vue";

const props = defineProps({
  mobile: String,
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
  if (!validateMobile(props.mobile)) {
    return ElMessage.error("请输入合法的手机号")
  }

  showCaptcha.value = true
  handleRequestCaptCode() // 每次点开都刷新验证码
}

const sendMsg = () => {
  if (!canSend.value) {
    return
  }

  canSend.value = false
  httpPost('/api/sms/code', {mobile: props.mobile, key: captKey.value, dots: dots.value}).then(() => {
    ElMessage.success('短信发送成功')
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
    ElMessage.error('短信发送失败：' + e.message)
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