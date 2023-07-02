<template>
  <el-button type="primary" class="sms-btn" :disabled="!canSend" :size="size" @click="sendMsg" plain>{{
      btnText
    }}
  </el-button>
</template>

<script setup>
// 发送短信验证码组件
import {ref} from "vue";
import {validateMobile} from "@/utils/validate";
import {ElMessage} from "element-plus";
import {httpGet, httpPost} from "@/utils/http";

const props = defineProps({
  mobile: String,
  size: String,
});
const btnText = ref('发送验证码')
const canSend = ref(true)

const sendMsg = () => {
  if (!canSend.value) {
    return
  }

  if (!validateMobile(props.mobile)) {
    return ElMessage.error("请输入合法的手机号")
  }
  canSend.value = false
  httpGet('/api/verify/token').then(res => {
    httpPost('/api/verify/sms', {token: res.data, mobile: props.mobile}).then(() => {
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
  }).catch(e => {
    console.log('failed to fetch token: ' + e.message)
  })
}

</script>

<style scoped>

</style>