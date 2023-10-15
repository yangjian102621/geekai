<template>
  <el-dialog
      class="login-dialog"
      v-model="showDialog"
      :close-on-click-modal="true"
      :show-close="true"
      :before-close="close"
      :width="400"
      title="用户登录"
  >
    <div class="form">
      <el-form label-width="65px">
        <el-form-item>
          <template #label>
            <div class="label">
              <el-icon>
                <User/>
              </el-icon>
              <span>账号</span>
            </div>
          </template>
          <template #default>
            <el-input v-model="username" placeholder="手机号码"/>
          </template>
        </el-form-item>
        <el-form-item>
          <template #label>
            <div class="label">
              <el-icon>
                <Lock/>
              </el-icon>
              <span>密码</span>
            </div>
          </template>
          <template #default>
            <el-input v-model="password" type="password" placeholder="密码"/>
          </template>
        </el-form-item>

        <div class="login-btn">
          <el-button type="primary" @click="submit" round>登录</el-button>
        </div>
      </el-form>
    </div>
  </el-dialog>
</template>

<script setup>
import {computed, ref} from "vue"
import {httpPost} from "@/utils/http";
import {ElMessage} from "element-plus";
import {setUserToken} from "@/store/session";
import {validateMobile} from "@/utils/validate";
import {Lock, User} from "@element-plus/icons-vue";

// eslint-disable-next-line no-undef
const props = defineProps({
  show: Boolean,
});
const showDialog = computed(() => {
  return props.show
})
const username = ref("")
const password = ref("")
// eslint-disable-next-line no-undef
const emits = defineEmits(['hide']);
const submit = function () {
  if (!validateMobile(username.value)) {
    return ElMessage.error('请输入合法的手机号');
  }
  if (password.value.trim() === '') {
    return ElMessage.error('请输入密码');
  }

  httpPost('/api/user/login', {username: username.value.trim(), password: password.value.trim()}).then((res) => {
    setUserToken(res.data)
    ElMessage.success("登录成功！")
    emits("hide")
  }).catch((e) => {
    ElMessage.error('登录失败，' + e.message)
  })
}
const close = function () {
  emits('hide', false);
}
</script>

<style lang="stylus">
.login-dialog {
  border-radius 20px

  .label {
    .el-icon {
      font-size 16px
      margin-right 6px
    }
  }

  .login-btn {
    text-align center
    padding-top 10px

    .el-button {
      width 50%
    }
  }
}
</style>