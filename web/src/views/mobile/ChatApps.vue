<template>
  <div class="mobile-setting container">
    <van-nav-bar :title="title"/>

    <div class="content">
      <van-empty description="功能正在开发中"/>
    </div>
  </div>
</template>

<script setup>
import {onMounted, ref} from "vue";
import {httpGet, httpPost} from "@/utils/http";
import {showFailToast, showSuccessToast} from "vant";
import {ElMessage} from "element-plus";

const title = ref('图片创作广场')
const form = ref({
  chat_config: {
    api_keys: {OpenAI: "", Azure: "", ChatGLM: ""}
  }
})
const showPicker = ref(false)
const models = ref([])

onMounted(() => {
  // 获取最新用户信息
  httpGet('/api/user/profile').then(res => {
    form.value = res.data
    form.value.chat_config.api_keys = res.data.chat_config.api_keys ?? {OpenAI: "", Azure: "", ChatGLM: ""}
  }).catch(() => {
    showFailToast('获取用户信息失败')
  });

  // 加载系统配置
  httpGet('/api/admin/config/get?key=system').then(res => {
    const mds = res.data.models;
    mds.forEach(item => {
      models.value.push({text: item, value: item})
    })
  }).catch(e => {
    ElMessage.error("加载系统配置失败: " + e.message)
  })
})

const selectModel = (item) => {
  showPicker.value = false
  form.value.chat_config.model = item.selectedValues[0]
}

const save = () => {
  httpPost('/api/user/profile/update', form.value).then(() => {
    showSuccessToast('保存成功')
  }).catch(() => {
    showFailToast('保存失败')
  })
}

</script>

<style lang="stylus">
.mobile-setting {
  .content {
    padding-top 60px

    .van-field__label {
      width 100px
      text-align right
    }
  }
}
</style>
