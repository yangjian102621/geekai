<template>
  <div class="mobile-setting container">
    <van-nav-bar :title="title"/>

    <div class="content">
      <van-form @submit="save" v-model="form">
        <van-cell-group inset>
          <van-field
              v-model="form.chat_config.model"
              is-link
              readonly
              label="默认模型"
              placeholder=""
              @click="showPicker = true"
          />
          <van-field
              v-model.number="form.chat_config.max_tokens"
              name="MaxTokens"
              type="number"
              label="MaxTokens"
              placeholder="每次请求最大 token 数量"
              :rules="[{ required: true, message: '请填写 MaxTokens' }]"
          />
          <van-field
              v-model.number="form.chat_config.temperature"
              type="number"
              name="Temperature"
              label="Temperature"
              placeholder="模型温度"
              :rules="[{ required: true, message: '请填写 Temperature' }]"
          />

          <van-field name="switch" label="聊天记录">
            <template #input>
              <van-switch v-model="form.chat_config.enable_history"/>
            </template>
          </van-field>

          <van-field name="switch" label="聊天上下文">
            <template #input>
              <van-switch v-model="form.chat_config.enable_context"/>
            </template>
          </van-field>
          <van-field
              v-model="form.chat_config.api_key"
              name="API KEY"
              label="API KEY"
              placeholder="配置自己的 api key"
              :rules="[{ required: true, message: '请填写 API KEY' }]"
          />
        </van-cell-group>
        <div style="margin: 16px;">
          <van-button round block type="primary" native-type="submit">
            保存
          </van-button>
        </div>
      </van-form>
    </div>

    <van-popup v-model:show="showPicker" round position="bottom">
      <van-picker
          :columns="models"
          @cancel="showPicker = false"
          @confirm="selectModel"
      />
    </van-popup>
  </div>
</template>

<script setup>
import {onMounted, ref} from "vue";
import {httpGet, httpPost} from "@/utils/http";
import {showFailToast, showSuccessToast} from "vant";
import {ElMessage} from "element-plus";

const title = ref('聊天设置')
const form = ref({
  chat_config: {
    model: '',
    max_tokens: 0,
    enable_context: false,
    enable_history: false,
    temperature: false,
    api_key: ''
  }
})
const showPicker = ref(false)
const models = ref([])

onMounted(() => {
  // 获取最新用户信息
  httpGet('/api/user/profile').then(res => {
    console.log(res.data)
    form.value = res.data
  }).catch(() => {
    showFailToast('获取用户信息失败')
  });

  // 加载系统配置
  httpGet('/api/admin/config/get?key=system').then(res => {
    const mds = res.data.models;
    console.log(mds)
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

<style scoped lang="stylus">
.mobile-setting {
  .content {
    padding-top 60px
  }
}
</style>