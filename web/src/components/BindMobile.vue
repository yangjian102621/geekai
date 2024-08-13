<template>
  <el-dialog
      v-model="showDialog"
      :close-on-click-modal="true"
      style="max-width: 600px"
      :before-close="close"
      :title="title"
  >
    <div class="form">
      <div class="text-center">当前已绑手机号：{{ mobile }}</div>

      <el-form :model="form" label-width="120px">
        <el-form-item label="手机号">
          <el-input v-model="form.mobile"/>
        </el-form-item>
        <el-form-item label="验证码">
          <el-row :gutter="20">
            <el-col :span="16">
              <el-input v-model="form.code" maxlength="6"/>
            </el-col>
            <el-col :span="8">
              <send-msg size="" :receiver="form.username" type="mobile"/>
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
import {checkSession} from "@/store/cache";

const props = defineProps({
  show: Boolean,
});

const showDialog = computed(() => {
  return props.show
})

const title = ref('绑定手机')
const mobile = ref('')
const form = ref({
  mobile: '',
  code: ''
})

checkSession().then(user => {
  mobile.value = user.mobile
})

const emits = defineEmits(['hide']);

const save = () => {
  if (!validateMobile(form.value.mobile) && !validateEmail(form.value.mobile)) {
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
.form {
  .text-center {
    text-align center
    padding-bottom 15px
    font-size 14px
    color #a1a1a1
    font-weight 700
  }

  .el-form-item__content {
    .el-row {
      width 100%
    }
  }
}
</style>