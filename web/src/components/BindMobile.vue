<template>
  <el-dialog
      v-model="showDialog"
      :close-on-click-modal="true"
      style="max-width: 600px"
      :before-close="close"
      :title="title"
  >
    <div class="form" id="bind-mobile-form">
      <el-alert v-if="mobile !== ''" type="info" show-icon :closable="false" style="margin-bottom: 20px;">
        <p>当前用户已绑定手机号：{{ mobile }}, 绑定其他手机号之后自动解绑该手机号。</p>
      </el-alert>

      <el-form :model="form" label-width="120px">
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
import {validateMobile} from "@/utils/validate";

const props = defineProps({
  show: Boolean,
  mobile: String
});

const showDialog = computed(() => {
  return props.show
})

const title = ref('绑定手机号')
const form = ref({
  mobile: '',
  code: ''
})

const emits = defineEmits(['hide']);

const save = () => {
  if (!validateMobile(form.value.mobile)) {
    return ElMessage.error("请输入正确的手机号码");
  }
  if (form.value.code === '') {
    return ElMessage.error("请输入短信验证码");
  }

  httpPost('/api/user/bind/mobile', form.value).then(() => {
    ElMessage.success({
      message: '绑定成功',
      appendTo: '#bind-mobile-form',
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