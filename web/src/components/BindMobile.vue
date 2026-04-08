<template>
  <div class="form px-3">
    <div class="text-center" v-if="mobile !== ''">当前已绑手机号：{{ mobile }}</div>
    <el-form label-position="top">
      <el-form-item label="手机号">
        <el-input v-model="form.mobile" />
      </el-form-item>
      <el-form-item label="验证码">
        <el-row :gutter="0">
          <el-col :span="16">
            <el-input v-model="form.code" maxlength="6" />
          </el-col>
          <el-col :span="8" style="padding-left: 10px">
            <send-msg :receiver="form.mobile" size="default" type="mobile" />
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

const mobile = ref('')
const form = ref({
  mobile: '',
  code: '',
})

onMounted(() => {
  checkSession().then((user) => {
    mobile.value = user.mobile
  })
})

const emits = defineEmits(['hide'])

const save = () => {
  if (form.value.code === '') {
    return ElMessage.error('请输入验证码')
  }

  httpPost('/api/user/bind/mobile', form.value)
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
