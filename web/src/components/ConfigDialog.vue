<template>
  <el-dialog
      v-model="props.show"
      :close-on-click-modal="false"
      :show-close="true"
      :before-close="close"
      :top="top"
      title="用户设置"
  >
    <div class="user-info" id="user-info">
      <el-form :model="form" label-width="120px">
        <el-form-item label="昵称">
          <el-input v-model="form['nickname']"/>
        </el-form-item>
        <el-form-item label="头像">
          <el-input v-model="form['avatar']"/>
        </el-form-item>
        <el-form-item label="用户名">
          <el-input v-model="form['username']" disabled/>
        </el-form-item>

        <el-form-item label="聊天上下文">
          <el-switch v-model="form['chat_config']['enable_context']"/>
        </el-form-item>
        <el-form-item label="聊天记录">
          <el-switch v-model="form['chat_config']['enable_history']"/>
        </el-form-item>
        <el-form-item label="Model">
          <el-select v-model="form['chat_config']['model']" placeholder="默认会话模型">
            <el-option
                v-for="item in props.models"
                :key="item"
                :label="item.toUpperCase()"
                :value="item"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="MaxTokens">
          <el-input v-model.number="form['chat_config']['max_tokens']"/>
        </el-form-item>
        <el-form-item label="Temperature">
          <el-input v-model.number="form['chat_config']['temperature']"/>
        </el-form-item>
        <el-form-item label="剩余调用次数">
          <el-tag>{{ form['calls'] }}</el-tag>
        </el-form-item>
        <el-form-item label="剩余 Tokens">
          <el-tag type="info">{{ form['tokens'] }}</el-tag>
        </el-form-item>
        <el-form-item label="API KEY">
          <el-input v-model="form['chat_config']['api_key']"/>
        </el-form-item>
      </el-form>
    </div>

    <template #footer>
      <span class="dialog-footer">
        <el-button @click="close">关闭</el-button>
        <el-button type="primary" @click="save">
          保存
        </el-button>
      </span>
    </template>
  </el-dialog>
</template>

<script setup>


import {computed, defineEmits, defineProps, onMounted, ref} from "vue"
import {httpGet, httpPost} from "@/utils/http";
import {ElMessage} from "element-plus";

const props = defineProps({
  show: Boolean,
  models: Array,
});

const form = ref({})
const top = computed(() => {
  if (window.innerHeight < 768) {
    return '1vh';
  } else {
    return '15vh';
  }
})

onMounted(() => {
  // 获取最新用户信息
  httpGet('/api/user/profile').then(res => {
    form.value = res.data
  }).catch(() => {
    ElMessage.error('获取用户信息失败')
  });
})

const emits = defineEmits(['update:show']);
const save = function () {
  httpPost('/api/user/profile/update', form.value).then(() => {
    ElMessage.success({
      message: '更新成功',
      appendTo: document.getElementById('user-info'),
      onClose: () => emits('update:show', false)
    })
  }).catch(() => {
    ElMessage.error({
      message: '更新失败',
      appendTo: document.getElementById('user-info')
    })
  })
}
const close = function () {
  emits('update:show', false);
}
</script>

<style lang="stylus">
.el-dialog {
  --el-dialog-width 90%;
  max-width 800px;

  .el-dialog__body {
    padding-top 10px;
    max-height 600px;
    overflow-y auto;

    .user-info {
      position relative;

      .el-message {
        position: absolute;
      }
    }
  }
}
</style>