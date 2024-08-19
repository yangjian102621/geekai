<template>
  <div class="reset-pass">
    <el-dialog
        v-model="showDialog"
        :close-on-click-modal="true"
        width="540px"
        :before-close="close"
        :title="title"
        class="reset-pass-dialog"
    >
      <div class="form">
        <el-form :model="form" label-width="80px" label-position="left">
          <el-tabs v-model="form.type" class="demo-tabs">
            <el-tab-pane label="手机号验证" name="mobile">
              <el-form-item label="手机号">
                <el-input v-model="form.mobile" placeholder="请输入手机号"/>
              </el-form-item>
              <el-form-item label="验证码">
                <el-row class="code-row">
                  <el-col :span="16">
                    <el-input v-model="form.code" maxlength="6"/>
                  </el-col>
                  <el-col :span="8" class="send-button">
                    <send-msg size="" :receiver="form.mobile" type="mobile"/>
                  </el-col>
                </el-row>
              </el-form-item>

            </el-tab-pane>
            <el-tab-pane label="邮箱验证" name="email">
              <el-form-item label="邮箱地址">
                <el-input v-model="form.email" placeholder="请输入邮箱地址"/>
              </el-form-item>
              <el-form-item label="验证码">
                <el-row class="code-row">
                  <el-col :span="16">
                    <el-input v-model="form.code" maxlength="6"/>
                  </el-col>
                  <el-col :span="8" class="send-button">
                    <send-msg size="" :receiver="form.email" type="email"/>
                  </el-col>
                </el-row>
              </el-form-item>
            </el-tab-pane>
          </el-tabs>

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
import {validateEmail, validateMobile} from "@/utils/validate";

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
  email: '',
  type: 'mobile',
  code: '',
  password: '',
  repass: ''
})

const emits = defineEmits(['hide']);

const save = () => {
  if (form.value.code === '') {
    return ElMessage.error("请输入验证码");
  }
  if (form.value.password.length < 8) {
    return ElMessage.error("密码长度必须大于8位");
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
    padding 0 20px
  }

  .code-row {
    width 100%
    .send-button {
      padding-left 10px
    }
  }

  .reset-pass-dialog {
    .el-dialog__footer {
      text-align center
      padding-top 0
    }
    .el-dialog__body {
      padding 0
    }
  }
}

</style>