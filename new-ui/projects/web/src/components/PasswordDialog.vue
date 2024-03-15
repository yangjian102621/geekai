<template>
  <el-dialog
      class="password-dialog"
      v-model="showDialog"
      :close-on-click-modal="true"
      :show-close="true"
      style="max-width: 600px"
      :before-close="close"
      title="修改密码"
  >
    <div class="form" id="password-form">
      <el-form :model="form" label-width="120px">
        <el-form-item label="原始密码">
          <el-input v-model="form['old_pass']" type="password"/>
        </el-form-item>
        <el-form-item label="新密码">
          <el-input v-model="form['password']" type="password"/>
        </el-form-item>
        <el-form-item label="确认密码">
          <el-input v-model="form['repass']" type="password"/>
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
import {computed, ref} from "vue"
import {httpPost} from "@/utils/http";
import {ElMessage} from "element-plus";

const props = defineProps({
  show: Boolean,
});
const showDialog = computed(() => {
  return props.show
})
const form = ref({})
const emits = defineEmits(['hide', 'logout']);
const save = function () {
  if (!form.value['password'] || form.value['password'].length < 8) {
    return ElMessage.error({message: "密码的长度为8-16个字符", appendTo: "#password-form"});
  }
  if (form.value['repass'] !== form.value['password']) {
    return ElMessage.error({message: '两次输入密码不一致', appendTo: '#password-form'});
  }
  httpPost('/api/user/password', form.value).then(() => {
    ElMessage.success({
      message: '更新成功',
      appendTo: '#password-form',
      duration: 1000,
      onClose: () => emits('logout', false)
    })
  }).catch((e) => {
    ElMessage.error({
      message: '更新失败，' + e.message,
      appendTo: '#password-form'
    })
  })
}
const close = function () {
  emits('hide', false);
}
</script>

<style lang="stylus">
.password-dialog {
  .el-dialog {
    --el-dialog-width 90%;
    max-width 650px;

    .el-dialog__body {
      overflow-y auto;

      .form {
        position relative;

        .el-message {
          position: absolute;
        }
      }
    }
  }
}
</style>