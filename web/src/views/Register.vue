<template>
  <div>
    <div class="bg"></div>
    <div class="register-page">
      <div class="page-inner">
        <div class="contain" v-if="enableRegister">
          <div class="logo">
            <el-image :src="logo" fit="cover" @click="router.push('/')"/>
          </div>

          <div class="header">{{ title }}</div>
          <div class="content">
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
                          <Message/>
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
                <el-col :span="24">
                  <el-button class="login-btn" type="primary" size="large" @click="submitRegister">注册</el-button>
                </el-col>
              </el-row>

              <el-row class="text-line">
                已经有账号？
                <el-link type="primary" @click="router.push('/login')">登录</el-link>
              </el-row>
            </el-form>


          </div>
        </div>

        <div class="tip-result" v-else>
          <el-result icon="error" title="注册功能已关闭">
            <template #sub-title>
              <p>抱歉，系统已关闭注册功能，请联系管理员添加账号！</p>
              <div class="wechat-card">
                <el-image :src="wxImg"/>
              </div>
            </template>
          </el-result>
        </div>

        <footer class="footer"  v-if="!licenseConfig.de_copy">
          <footer-bar/>
        </footer>
      </div>
    </div>
  </div>
</template>

<script setup>

import {ref} from "vue";
import {Checked, Iphone, Lock, Message} from "@element-plus/icons-vue";
import {httpGet, httpPost} from "@/utils/http";
import {ElMessage} from "element-plus";
import {useRouter} from "vue-router";
import FooterBar from "@/components/FooterBar.vue";
import SendMsg from "@/components/SendMsg.vue";
import {arrayContains} from "@/utils/libs";
import {setUserToken} from "@/store/session";
import {validateEmail, validateMobile} from "@/utils/validate";
import {showMessageError, showMessageOK} from "@/utils/dialog";
import {getLicenseInfo, getSystemInfo} from "@/store/cache";

const router = useRouter();
const title = ref('');
const logo = ref("")
const data = ref({
  username: '',
  password: '',
  code: '',
  repass: '',
  invite_code: router.currentRoute.value.query['invite_code'],
})

const enableMobile = ref(false)
const enableEmail = ref(false)
const enableUser = ref(false)
const enableRegister = ref(true)
const activeName = ref("mobile")
const wxImg = ref("/images/wx.png")
const licenseConfig = ref({})

getSystemInfo().then(res => {
  if (res.data) {
    title.value = res.data.title
    logo.value = res.data.logo
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

getLicenseInfo().then(res => {
  licenseConfig.value = res.data
}).catch(e => {
  showMessageError("获取 License 配置：" + e.message)
})

// 注册操作
const submitRegister = () => {
  if (data.value.username === '') {
    return showMessageError('请输入用户名');
  }

  if (activeName.value === 'mobile' && !validateMobile(data.value.username)) {
    return showMessageError('请输入合法的手机号');
  }

  if (activeName.value === 'email' && !validateEmail(data.value.username)) {
    return showMessageError('请输入合法的邮箱地址');
  }

  if (data.value.password.length < 8) {
    return showMessageError('密码的长度为8-16个字符');
  }
  if (data.value.repass !== data.value.password) {
    return showMessageError('两次输入密码不一致');
  }

  if ((activeName.value === 'mobile' || activeName.value === 'email') && data.value.code === '') {
    return showMessageError('请输入验证码');
  }
  data.value.reg_way = activeName.value
  httpPost('/api/user/register', data.value).then((res) => {
    setUserToken(res.data)
    showMessageOK({
      "message": "注册成功，即将跳转到对话主界面...",
      onClose: () => router.push("/chat"),
      duration: 1000
    })
  }).catch((e) => {
    showMessageError('注册失败，' + e.message)
  })
}

</script>

<style lang="stylus" scoped>
.bg {
  position fixed
  left 0
  right 0
  top 0
  bottom 0
  background-color #091519
  background-image url("~@/assets/img/reg-bg.jpg")
  background-size cover
  background-position center
  background-repeat no-repeat
  //filter: blur(10px); /* 调整模糊程度，可以根据需要修改值 */
}

.register-page {
  display flex
  justify-content center

  .page-inner {
    max-width 450px
    width 100%
    height 100vh
    display flex
    justify-content center
    align-items center

    .contain {
      padding 20px 40px 20px 40px;
      width 100%
      color #ffffff
      border-radius 10px;
      z-index 10
      background-color rgba(255, 255, 255, 0.2)

      .logo {
        text-align center

        .el-image {
          width 120px;
          cursor pointer
        }
      }

      .header {
        width 100%
        margin-bottom 24px
        font-size 24px
        color $white_v1
        letter-space 2px
        text-align center
        padding-top 10px
      }

      .content {
        width 100%
        height: auto
        border-radius 3px

        .block {
          margin-bottom 16px

          .el-input__inner {
            border 1px solid $gray-v6 !important

            .el-icon-user, .el-icon-lock {
              font-size 20px
            }
          }
        }

        .btn-row {
          padding-top 10px;

          .login-btn {
            width 100%
            font-size 16px
            letter-spacing 2px
          }
        }

        .text-line {
          justify-content center
          padding-top 10px;
          font-size 14px;
        }
      }
    }


    .tip-result {
      z-index 10

      .wechat-card {
        padding 20px
      }
    }

    .footer {
      color #ffffff;

      .container {
        padding 20px;
      }
    }

  }

}
</style>

<style lang="stylus">
.register-page {
  .el-result {

    border-radius 10px;
    background-color rgba(14, 25, 30, 0.6)
    border 1px solid #666

    .el-result__title p {
      color #ffffff
    }

    .el-result__subtitle p {
      color #c1c1c1
    }
  }

  .el-tabs__item {
    color #ffffff
  }
}
</style>