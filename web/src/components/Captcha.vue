<template>
  <el-container class="captcha-box">
    <el-dialog
        v-model="show"
        :close-on-click-modal="true"
        :show-close="false"
        style="width: 360px;"
    >
      <slide-captcha
          v-if="isMobile()"
          :bg-img="bgImg"
          :bk-img="bkImg"
          :result="result"
          @refresh="getSlideCaptcha"
          @confirm="handleSlideConfirm"
          @hide="show = false"/>

      <captcha-plus
          v-else
          :max-dot="maxDot"
          :image-base64="imageBase64"
          :thumb-base64="thumbBase64"
          width="300"
          @close="show = false"
          @refresh="handleRequestCaptCode"
          @confirm="handleConfirm"
      />
    </el-dialog>

  </el-container>
</template>

<script setup>
import {ref} from "vue";
import lodash from 'lodash'
import {validateEmail, validateMobile} from "@/utils/validate";
import {httpGet, httpPost} from "@/utils/http";
import CaptchaPlus from "@/components/CaptchaPlus.vue";
import SlideCaptcha from "@/components/SlideCaptcha.vue";
import {isMobile} from "@/utils/libs";
import {showMessageError, showMessageOK} from "@/utils/dialog";

const show = ref(false)
const maxDot = ref(5)
const imageBase64 = ref('')
const thumbBase64 = ref('')
const captKey = ref('')
const dots = ref(null)

const emits = defineEmits(['success']);
const handleRequestCaptCode = () => {

  httpGet('/api/captcha/get').then(res => {
    const data = res.data
    imageBase64.value = data.image
    thumbBase64.value = data.thumb
    captKey.value = data.key
  }).catch(e => {
    showMessageError('获取人机验证数据失败：' + e.message)
  })
}

const handleConfirm = (dts) => {
  if (lodash.size(dts) <= 0) {
    return showMessageError('请进行人机验证再操作')
  }

  let dotArr = []
  lodash.forEach(dts, (dot) => {
    dotArr.push(dot.x, dot.y)
  })
  dots.value = dotArr.join(',')
  httpPost('/api/captcha/check', {
    dots: dots.value,
    key: captKey.value
  }).then(() => {
    // ElMessage.success('人机验证成功')
    show.value = false
    emits('success', {key:captKey.value, dots:dots.value})
  }).catch(() => {
    showMessageError('人机验证失败')
    handleRequestCaptCode()
  })
}

const loadCaptcha = () => {
  show.value = true
  // 手机用滑动验证码
  if (isMobile()) {
    getSlideCaptcha()
  } else {
    handleRequestCaptCode()
  }
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
    showMessageError('获取人机验证数据失败：' + e.message)
  })
}

const handleSlideConfirm = (x) => {
  httpPost("/api/captcha/slide/check", {
    key: captKey.value,
    x: x
  }).then(() => {
    result.value = 1
    show.value = false
    emits('success',{key:captKey.value, x:x})
  }).catch(() => {
    result.value = 2
  })
}

// 导出方法以便父组件调用
defineExpose({
  loadCaptcha
});
</script>

<style lang="stylus">

.captcha-box {
  .el-dialog {
    .el-dialog__header {
      padding: 0;
    }

    .el-dialog__body {
      padding 0
    }
  }
}
</style>