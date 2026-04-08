<template>
  <div class="form" id="redeem-form">
    <el-form :model="form">
      <el-form-item>
        <el-input v-model="form.code" placeholder="请输入兑换码" />
      </el-form-item>
    </el-form>
    <div class="dialog-footer">
      <el-button type="primary" @click="save">兑换</el-button>
      <el-button @click="emits('hide')">取消</el-button>
    </div>
  </div>
</template>

<script setup>
import { showMessageError, showMessageInfo, showMessageOK } from '@/utils/dialog'
import { httpPost } from '@/utils/http'
import { ref } from 'vue'

const form = ref({
  code: '',
})

const emits = defineEmits(['hide'])

const save = () => {
  if (form.value.code === '') {
    return showMessageInfo('请输入兑换码')
  }

  httpPost('/api/redeem/verify', form.value)
    .then(() => {
      showMessageOK('兑换成功！')
      emits('hide', true)
    })
    .catch((e) => {
      showMessageError('兑换失败：' + e.message)
    })
}
</script>

<style scoped lang="scss">
.form {
  padding: 20px;

  .form-title {
    font-size: 18px;
    font-weight: 600;
    color: var(--el-text-color-primary);
    margin-bottom: 20px;
    text-align: center;
  }

  .dialog-footer {
    text-align: center;
    margin-top: 20px;
  }
}
</style>
