<template>
  <el-dialog
      class="login-dialog"
      v-model="showDialog"
      :close-on-click-modal="true"
      :show-close="false"
      :before-close="close"
  >
    <template #header="{ close, titleId, titleClass }">
      <div class="header">
        <div class="title">用户登录</div>
        <div class="close-icon">
          <el-icon>
            <Close/>
          </el-icon>
        </div>
      </div>
    </template>

    <div class="login-box">
      <el-form label-width="75px">
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
            <el-input v-model="data.username" size="large" placeholder="账号"/>
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
            <el-input v-model="data.password" type="password" size="large" placeholder="密码"/>
          </template>
        </el-form-item>

        <div class="login-btn">
          <el-button type="primary" @click="submit" size="large" round>登录</el-button>
          <el-button plain @click="submit" size="large" round>注册</el-button>
        </div>
      </el-form>
    </div>

    <div class="register-box">
      <el-form :model="data" label-width="120px" ref="formRef">
        <div class="block">
          <el-input placeholder="账号"
                    size="large"
                    v-model="data.username"
                    autocomplete="off">
            <template #prefix>
              <el-icon>
                <Iphone/>
              </el-icon>
            </template>
          </el-input>
        </div>

        <div class="block">
          <el-input placeholder="请输入密码(8-16位)"
                    maxlength="16" size="large"
                    v-model="data.password" show-password
                    autocomplete="off">
            <template #prefix>
              <el-icon>
                <Lock/>
              </el-icon>
            </template>
          </el-input>
        </div>

        <div class="block">
          <el-input placeholder="重复密码(8-16位)"
                    size="large" maxlength="16" v-model="data.repass" show-password
                    autocomplete="off">
            <template #prefix>
              <el-icon>
                <Lock/>
              </el-icon>
            </template>
          </el-input>
        </div>

        <div class="block">
          <el-row :gutter="10">
            <el-col :span="12">
              <el-input placeholder="验证码"
                        size="large" maxlength="30"
                        v-model="data.code"
                        autocomplete="off">
                <template #prefix>
                  <el-icon>
                    <Checked/>
                  </el-icon>
                </template>
              </el-input>
            </el-col>
            <el-col :span="12">
              <send-msg size="large" :receiver="data.username"/>
            </el-col>
          </el-row>
        </div>

        <div class="block">
          <el-input placeholder="邀请码(可选)"
                    size="large"
                    v-model="data.invite_code"
                    autocomplete="off">
            <template #prefix>
              <el-icon>
                <Message/>
              </el-icon>
            </template>
          </el-input>
        </div>

        <el-row class="btn-row">
          <el-button class="login-btn" type="primary" @click="">注册</el-button>
        </el-row>

        <el-row class="text-line">
          已经有账号？
          <el-link type="primary" @click="">登录</el-link>
        </el-row>
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
import {Checked, Close, Iphone, Lock, Message, Position, User} from "@element-plus/icons-vue";
import SendMsg from "@/components/SendMsg.vue";

// eslint-disable-next-line no-undef
const props = defineProps({
  show: Boolean,
});
const showDialog = computed(() => {
  return props.show
})
const data = ref({
  username: "",
  password: ""
})
// eslint-disable-next-line no-undef
const emits = defineEmits(['hide']);
const submit = function () {
  httpPost('/api/user/login', data.value).then((res) => {
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
  border-radius 10px
  max-width 600px

  .header {
    position relative

    .title {
      padding 0
      font-size 18px
    }

    .close-icon {
      cursor pointer
      position absolute
      right 0
      top 0
      font-weight normal
      font-size 20px

      &:hover {
        color #20a0ff
      }
    }
  }

  .login-box {
    .label {
      padding-top 3px

      .el-icon {
        position relative
        font-size 20px
        margin-right 6px
        top 4px
      }

      span {
        font-size 16px
      }
    }

    .login-btn {
      display flex
      padding-top 10px
      justify-content center

      .el-button {
        font-size 16px
        width 100px
      }
    }
  }

}
</style>