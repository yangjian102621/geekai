<template>
  <div class="form px-3">
    <div class="text-center" v-if="email !== ''">当前已绑定邮箱：{{ email }}</div>
    <el-form label-position="top">
      <el-form-item label="邮箱地址">
        <el-input v-model="form.email" />
      </el-form-item>
      <el-form-item label="验证码">
        <div class="flex w-full items-center gap-2">
          <el-input v-model="form.code" maxlength="6" class="flex-1" />
          <span class="flex-none">
            <send-msg :receiver="form.email" type="email" />
          </span>
        </div>
      </el-form-item>
    </el-form>
    <div class="pt-3 flex justify-end">
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
  .el-form-item__content {
    .el-row {
      width: 100%;
    }
  }
}
</style>
