<template>
  <el-dialog
      v-model="showDialog"
      :close-on-click-modal="true"
      style="max-width: 600px"
      :before-close="close"
      :title="title"
  >
    <div class="form" id="bind-mobile-form">
      <el-alert v-if="username !== ''" type="info" show-icon :closable="false" style="margin-bottom: 20px;">
        <p>当前绑定账号：{{ username }}，只允许使绑定有效的手机号或者邮箱地址作为登录账号。</p>
      </el-alert>

      <el-form :model="form" label-width="120px">
        <el-form-item label="新账号">
          <el-input v-model="form.username"/>
        </el-form-item>
        <el-form-item label="验证码">
          <el-row :gutter="20">
            <el-col :span="16">
              <el-input v-model="form.code" maxlength="6"/>
            </el-col>
            <el-col :span="8">
              <send-msg size="" :receiver="form.username"/>
            </el-col>
          </el-row>
        </el-form-item>
      </el-form>
    </div>

    <template #footer>
      <span class="dialog-footer">
        <el-button type="primary" @click="save">
          提交绑定
        </el-button>
      </span>
    </template>
  </el-dialog>
</template>

<script setup>
import {computed, ref} from "vue";
import SendMsg from "@/components/SendMsg.vue";
import {ElMessage} from "element-plus";
import {httpPost} from "@/utils/http";
import {validateEmail, validateMobile} from "@/utils/validate";

const props = defineProps({
  show: Boolean,
  username: String
});

const showDialog = computed(() => {
  return props.show
})

const title = ref('重置登录账号')
const form = ref({
  username: '',
  code: ''
})

const emits = defineEmits(['hide']);

const save = () => {
  if (!validateMobile(form.value.username) && !validateEmail(form.value.username)) {
    return ElMessage.error("请输入合法的手机号/邮箱地址")
  }
  if (form.value.code === '') {
    return ElMessage.error("请输入验证码");
  }

  httpPost('/api/user/bind/username', form.value).then(() => {
    ElMessage.success({
      message: '绑定成功',
      duration: 1000,
      onClose: () => emits('hide', false)
    })
  }).catch(e => {
    ElMessage.error("绑定失败：" + e.message);
  })
}

const close = function () {
  emits('hide', false);
}
</script>

<style lang="stylus" scoped>
#bind-mobile-form {
  .el-form-item__content {
    .el-row {
      width 100%
    }
  }
}
</style>