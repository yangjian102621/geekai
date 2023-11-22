<template>
  <div class="reset-pass">
    <el-dialog
        v-model="showDialog"
        :close-on-click-modal="true"
        width="540px"
        :before-close="close"
        :title="title"
    >
      <div class="form">

        <el-form :model="form" label-width="120px" label-position="left">
          <el-form-item label="手机号码">
            <el-input v-model="form.mobile"/>
          </el-form-item>
          <el-form-item label="手机验证码">
            <el-row :gutter="20">
              <el-col :span="16">
                <el-input v-model="form.code" maxlength="6"/>
              </el-col>
              <el-col :span="8">
                <send-msg size="" :mobile="form.mobile"/>
              </el-col>
            </el-row>
          </el-form-item>
          <el-form-item label="新密码">
            <el-input v-model="form.password" type="password"/>
          </el-form-item>
          <el-form-item label="重复密码">
            <el-input v-model="form.repass" type="password"/>
          </el-form-item>
        </el-form>
      </div>

      <template #footer>
        <div class="dialog-footer">
          <el-button type="primary" @click="save" round>
            重置密码
          </el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import {computed, ref} from "vue";
import SendMsg from "@/components/SendMsg.vue";
import {ElMessage} from "element-plus";
import {httpPost} from "@/utils/http";
import {validateMobile} from "@/utils/validate";

const props = defineProps({
  show: Boolean,
  mobile: String
});

const showDialog = computed(() => {
  return props.show
})

const title = ref('重置密码')
const form = ref({
  mobile: '',
  code: '',
  password: '',
  repass: ''
})

const emits = defineEmits(['hide']);

const save = () => {
  if (!validateMobile(form.value.mobile)) {
    return ElMessage.error("请输入正确的手机号码");
  }
  if (form.value.code === '') {
    return ElMessage.error("请输入短信验证码");
  }
  if (form.value.repass !== form.value.password) {
    return ElMessage.error("两次输入密码不一致");
  }

  httpPost('/api/user/resetPass', form.value).then(() => {
    ElMessage.success({
      message: '重置密码成功', duration: 1000, onClose: () => emits('hide', false)
    })
  }).catch(e => {
    ElMessage.error("重置密码失败：" + e.message);
  })
}

const close = function () {
  emits('hide', false);
}
</script>

<style lang="stylus">
.reset-pass {
  .form {
    padding 10px 40px
  }

  .el-dialog__footer {
    text-align center
  }
}

</style>