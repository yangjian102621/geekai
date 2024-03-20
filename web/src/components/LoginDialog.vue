<template>
  <el-dialog
      class="login-dialog"
      v-model="showDialog"
      :close-on-click-modal="true"
      :show-close="false"
      :before-close="close"
  >
    <template #header="{titleId, titleClass }">
      <div class="header">
        <div class="title" v-if="login">用户登录</div>
        <div class="title" v-else>用户注册</div>
        <div class="close-icon">
          <el-icon @click="close">
            <Close/>
          </el-icon>
        </div>
      </div>
    </template>

    <div class="login-box" v-if="login">
      <el-form :model="data" label-width="120px" class="form">
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

        <el-row class="btn-row" :gutter="20">
          <el-col :span="12">
            <el-button class="login-btn" type="primary" size="large" @click="submitLogin">登录</el-button>
          </el-col>
          <el-col :span="12">
            <div class="text">
              还没有账号？
              <el-tag @click="login = false">注册</el-tag>
            </div>
          </el-col>

        </el-row>
      </el-form>
    </div>

    <div class="register-box" v-else>
      <el-form :model="data" class="form" v-if="enableRegister">
        <el-tabs v-model="activeName" class="demo-tabs">
          <el-tab-pane label="手机注册" name="mobile" v-if="enableMobile">
            <div class="block">
              <el-input placeholder="手机号码"
                        size="large"
                        v-model="data.username"
                        maxlength="11"
                        autocomplete="off">
                <template #prefix>
                  <el-icon>
                    <Iphone/>
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
          </el-tab-pane>
          <el-tab-pane label="邮箱注册" name="email" v-if="enableEmail">
            <div class="block">
              <el-input placeholder="邮箱地址"
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
          </el-tab-pane>
          <el-tab-pane label="用户名注册" name="username" v-if="enableUser">
            <div class="block">
              <el-input placeholder="用户名"
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
          </el-tab-pane>
        </el-tabs>

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

        <el-row class="btn-row" :gutter="20">
          <el-col :span="12">
            <el-button class="login-btn" type="primary" size="large" @click="submitRegister">注册</el-button>
          </el-col>
          <el-col :span="12">
            <div class="text">
              已有账号？
              <el-tag @click="login = true">登录</el-tag>
            </div>
          </el-col>

        </el-row>
      </el-form>

      <div class="tip-result" v-else>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-result icon="error" title="注册功能已关闭">
              <template #sub-title>
                <p>抱歉，系统已关闭注册功能，请联系管理员添加账号！</p>
              </template>
            </el-result>
          </el-col>

          <el-col :span="12">
            <div class="wechat-card">
              <el-image :src="wxImg"/>
            </div>
          </el-col>
        </el-row>
      </div>
    </div>
  </el-dialog>
</template>

<script setup>
import {computed, ref} from "vue"
import {httpGet, httpPost} from "@/utils/http";
import {ElMessage} from "element-plus";
import {setUserToken} from "@/store/session";
import {validateEmail, validateMobile} from "@/utils/validate";
import {Checked, Close, Iphone, Lock, Message, Position, User} from "@element-plus/icons-vue";
import SendMsg from "@/components/SendMsg.vue";
import {arrayContains} from "@/utils/libs";
import {useRouter} from "vue-router";

// eslint-disable-next-line no-undef
const props = defineProps({
  show: Boolean,
});
const showDialog = computed(() => {
  return props.show
})

const login = ref(true)
const data = ref({
  username: "",
  password: "",
  repass: "",
  code: "",
  invite_code: ""
})
const enableMobile = ref(false)
const enableEmail = ref(false)
const enableUser = ref(false)
const enableRegister = ref(false)
const activeName = ref("mobile")
const wxImg = ref("/images/wx.png")
// eslint-disable-next-line no-undef
const emits = defineEmits(['hide', 'success']);

httpGet("/api/config/get?key=system").then(res => {
  if (res.data) {
    const registerWays = res.data['register_ways']
    if (arrayContains(registerWays, "mobile")) {
      enableMobile.value = true
    }
    if (arrayContains(registerWays, "email")) {
      enableEmail.value = true
    }
    if (arrayContains(registerWays, "username")) {
      enableUser.value = true
    }
    // 是否启用注册
    enableRegister.value = res.data['enabled_register']
    // 使用后台上传的客服微信二维码
    if (res.data['wechat_card_url'] !== '') {
      wxImg.value = res.data['wechat_card_url']
    }
  }
}).catch(e => {
  ElMessage.error("获取系统配置失败：" + e.message)
})

// 登录操作
const submitLogin = () => {
  if (data.value.username === '') {
    return ElMessage.error('请输入用户名');
  }
  if (data.value.password === '') {
    return ElMessage.error('请输入密码');
  }

  httpPost('/api/user/login', data.value).then((res) => {
    setUserToken(res.data)
    ElMessage.success("登录成功！")
    emits("hide")
    emits('success')
  }).catch((e) => {
    ElMessage.error('登录失败，' + e.message)
  })
}

// 注册操作
const submitRegister = () => {
  if (data.value.username === '') {
    return ElMessage.error('请输入用户名');
  }

  if (activeName.value === 'mobile' && !validateMobile(data.value.username)) {
    return ElMessage.error('请输入合法的手机号');
  }

  if (activeName.value === 'email' && !validateEmail(data.value.username)) {
    return ElMessage.error('请输入合法的邮箱地址');
  }

  if (data.value.password.length < 8) {
    return ElMessage.error('密码的长度为8-16个字符');
  }
  if (data.value.repass !== data.value.password) {
    return ElMessage.error('两次输入密码不一致');
  }

  if ((activeName.value === 'mobile' || activeName.value === 'email') && data.value.code === '') {
    return ElMessage.error('请输入验证码');
  }
  data.value.reg_way = activeName.value
  httpPost('/api/user/register', data.value).then((res) => {
    setUserToken(res.data)
    ElMessage.success({
      "message": "注册成功!",
      onClose: () => {
        emits("hide")
        emits('success')
      },
      duration: 1000
    })
  }).catch((e) => {
    ElMessage.error('注册失败，' + e.message)
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
      right -10px
      top 0
      font-weight normal
      font-size 20px

      &:hover {
        color #20a0ff
      }
    }
  }


  .el-dialog__body {
    padding 10px 20px 20px 20px
  }

  .form {
    .block {
      margin-bottom 10px
    }

    .btn-row {
      display flex

      .el-button {
        width 100%
      }

      .text {
        line-height 40px

        .el-tag {
          cursor pointer
        }
      }

    }
  }

  .register-box {
    .wechat-card {
      text-align center
    }
  }

}
</style>