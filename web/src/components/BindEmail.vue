<template>
  <div class="form">
    <div class="text-center" v-if="email !== ''">当前已绑定邮箱：{{ email }}</div>
    <el-form label-position="top">
      <el-form-item label="邮箱地址">
        <el-input v-model="form.email" />
      </el-form-item>
      <el-form-item label="验证码">
        <el-row :gutter="0">
          <el-col :span="16">
            <el-input v-model="form.code" maxlength="6" />
          </el-col>
          <el-col :span="8" style="padding-left: 10px">
            <send-msg :receiver="form.email" type="email" />
          </el-col>
        </el-row>
      </el-form-item>
    </el-form>
    <div class="dialog-footer text-center">
      <el-button type="primary" @click="save"> 保存 </el-button>
      <el-button @click="emits('hide')"> 取消 </el-button>
    </div>
  </div>
</template>

<script setup>
import SendMsg from '@/components/SendMsg.vue'
import { checkSession } from '@/store/cache'
import { httpPost } from '@/utils/http'
import { ElMessage } from 'element-plus'
import { onMounted, ref } from 'vue'

const email = ref('')
const form = ref({
  email: '',
  code: '',
})

onMounted(() => {
  checkSession().then((user) => {
    email.value = user.email
  })
})

const emits = defineEmits(['hide'])

const save = () => {
  if (form.value.code === '') {
    return ElMessage.error('请输入验证码')
  }

  httpPost('/api/user/bind/email', form.value)
    .then(() => {
      ElMessage.success('绑定成功')
      emits('hide')
    })
    .catch((e) => {
      ElMessage.error('绑定失败：' + e.message)
    })
}
</script>

<style lang="scss" scoped>
.form {
  .text-center {
    text-align: center;
    padding-bottom: 15px;
    font-size: 14px;
    color: #a1a1a1;
    font-weight: 700;
  }

  .el-form-item__content {
    .el-row {
      width: 100%;
    }
  }
}
</style>
