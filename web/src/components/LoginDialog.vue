<template>
  <div class="login-dialog w-full">
    <div class="login-box" v-if="login">
      <el-form :model="data" class="form">
        <div class="block">
          <el-input placeholder="账号" size="large" v-model="data.username" autocomplete="off">
            <template #prefix>
              <el-icon>
                <Iphone />
              </el-icon>
            </template>
          </el-input>
        </div>

        <div class="block">
          <el-input placeholder="请输入密码(8-16位)" maxlength="16" size="large" v-model="data.password" show-password autocomplete="off">
            <template #prefix>
              <el-icon>
                <Lock />
              </el-icon>
            </template>
          </el-input>
        </div>

        <el-row class="btn-row" :gutter="20">
          <el-col :span="24">
            <el-button class="login-btn" type="primary" size="large" @click="submitLogin">登 录</el-button>
          </el-col>
        </el-row>

        <div class="w-full">
          <div class="text flex justify-center items-center pt-3 text-sm">
            还没有账号？
            <el-button size="small" @click="login = false">注册</el-button>

            <el-button type="info" class="forget" size="small" @click="showResetPass = true">忘记密码？</el-button>
          </div>
          <div v-if="wechatLoginURL !== ''">
            <el-divider>
              <div class="text-center">其他登录方式</div>
            </el-divider>
            <div class="c-login flex justify-center">
              <!-- <div class="login-type mr-2">
            <a class="wechat-login" :href="wechatLoginURL" @click="setRoute(router.currentRoute.value.path)"><i class="iconfont icon-wechat"></i></a>
          </div> -->
              <div class="p-2 w-full">
                <el-button type="success" class="w-full" size="large" :href="wechatLoginURL" @click="setRoute(router.currentRoute.value.path)"
                  ><i class="iconfont icon-wechat mr-2"></i> 微信登录
                </el-button>
              </div>
            </div>
          </div>
        </div>
      </el-form>
    </div>

    <div class="register-box w-full" v-else>
      <el-form :model="data" class="form" v-if="enableRegister">
        <el-tabs v-model="activeName" class="demo-tabs">
          <el-tab-pane label="手机注册" name="mobile" v-if="enableMobile">
            <div class="block">
              <el-input placeholder="手机号码" size="large" v-model="data.mobile" maxlength="11" autocomplete="off">
                <template #prefix>
                  <el-icon>
                    <Iphone />
                  </el-icon>
                </template>
              </el-input>
            </div>
            <div class="block">
              <el-row :gutter="10">
                <el-col :span="12">
                  <el-input placeholder="验证码" size="large" maxlength="30" v-model="data.code" autocomplete="off">
                    <template #prefix>
                      <el-icon>
                        <Checked />
                      </el-icon>
                    </template>
                  </el-input>
                </el-col>
                <el-col :span="12">
                  <send-msg size="large" :receiver="data.mobile" type="mobile" />
                </el-col>
              </el-row>
            </div>
          </el-tab-pane>
          <el-tab-pane label="邮箱注册" name="email" v-if="enableEmail">
            <div class="block">
              <el-input placeholder="邮箱地址" size="large" v-model="data.email" autocomplete="off">
                <template #prefix>
                  <el-icon>
                    <Message />
                  </el-icon>
                </template>
              </el-input>
            </div>
            <div class="block">
              <el-row :gutter="10">
                <el-col :span="12">
                  <el-input placeholder="验证码" size="large" maxlength="30" v-model="data.code" autocomplete="off">
                    <template #prefix>
                      <el-icon>
                        <Checked />
                      </el-icon>
                    </template>
                  </el-input>
                </el-col>
                <el-col :span="12">
                  <send-msg size="large" :receiver="data.email" type="email" />
                </el-col>
              </el-row>
            </div>
          </el-tab-pane>
          <el-tab-pane label="用户名注册" name="username" v-if="enableUser">
            <div class="block">
              <el-input placeholder="用户名" size="large" v-model="data.username" autocomplete="off">
                <template #prefix>
                  <el-icon>
                    <Iphone />
                  </el-icon>
                </template>
              </el-input>
            </div>
          </el-tab-pane>
        </el-tabs>

        <div class="block">
          <el-input placeholder="请输入密码(8-16位)" maxlength="16" size="large" v-model="data.password" show-password autocomplete="off">
            <template #prefix>
              <el-icon>
                <Lock />
              </el-icon>
            </template>
          </el-input>
        </div>

        <div class="block">
          <el-input placeholder="重复密码(8-16位)" size="large" maxlength="16" v-model="data.repass" show-password autocomplete="off">
            <template #prefix>
              <el-icon>
                <Lock />
              </el-icon>
            </template>
          </el-input>
        </div>

        <div class="block">
          <el-input placeholder="邀请码(可选)" size="large" v-model="data.invite_code" autocomplete="off">
            <template #prefix>
              <el-icon>
                <Message />
              </el-icon>
            </template>
          </el-input>
        </div>

        <div class="w-full">
          <el-button class="login-btn w-full" type="primary" size="large" @click="submitRegister">注 册</el-button>
        </div>

        <div class="text text-sm flex justify-center items-center w-full pt-3">
          已有账号？
          <el-button size="small" @click="login = true">登录</el-button>
        </div>
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
              <el-image :src="wxImg" />
            </div>
          </el-col>
        </el-row>
      </div>
    </div>
    <captcha v-if="enableVerify" @success="submit" ref="captchaRef" />

    <reset-pass @hide="showResetPass = false" :show="showResetPass" />
  </div>
</template>

<script setup>
import { onMounted, ref, watch } from "vue";
import { httpGet, httpPost } from "@/utils/http";
import { ElMessage } from "element-plus";
import { setUserToken } from "@/store/session";
import { validateEmail, validateMobile } from "@/utils/validate";
import { Checked, Close, Iphone, Lock, Message } from "@element-plus/icons-vue";
import SendMsg from "@/components/SendMsg.vue";
import { arrayContains } from "@/utils/libs";
import { getSystemInfo } from "@/store/cache";
import Captcha from "@/components/Captcha.vue";
import ResetPass from "@/components/ResetPass.vue";
import { setRoute } from "@/store/system";
import { useRouter } from "vue-router";
import { useSharedStore } from "@/store/sharedata";

// eslint-disable-next-line no-undef
const props = defineProps({
  show: Boolean,
});
const showDialog = ref(false);
watch(
  () => props.show,
  (newValue) => {
    showDialog.value = newValue;
  }
);

const login = ref(true);
const data = ref({
  username: process.env.VUE_APP_USER,
  password: process.env.VUE_APP_PASS,
  mobile: "",
  email: "",
  repass: "",
  code: "",
  invite_code: "",
});
const enableMobile = ref(false);
const enableEmail = ref(false);
const enableUser = ref(false);
const enableRegister = ref(true);
const wechatLoginURL = ref("");
const activeName = ref("");
const wxImg = ref("/images/wx.png");
const captchaRef = ref(null);
// eslint-disable-next-line no-undef
const emits = defineEmits(["hide", "success"]);
const action = ref("login");
const enableVerify = ref(false);
const showResetPass = ref(false);
const router = useRouter();
const store = useSharedStore();

onMounted(() => {
  const returnURL = `${location.protocol}//${location.host}/login/callback?action=login`;
  httpGet("/api/user/clogin?return_url=" + returnURL)
    .then((res) => {
      wechatLoginURL.value = res.data.url;
    })
    .catch((e) => {
      console.log(e.message);
    });

  getSystemInfo()
    .then((res) => {
      if (res.data) {
        const registerWays = res.data["register_ways"];
        if (arrayContains(registerWays, "username")) {
          enableUser.value = true;
          activeName.value = "username";
        }
        if (arrayContains(registerWays, "email")) {
          enableEmail.value = true;
          activeName.value = "email";
        }
        if (arrayContains(registerWays, "mobile")) {
          enableMobile.value = true;
          activeName.value = "mobile";
        }
        // 是否启用注册
        enableRegister.value = res.data["enabled_register"];
        // 使用后台上传的客服微信二维码
        if (res.data["wechat_card_url"] !== "") {
          wxImg.value = res.data["wechat_card_url"];
        }
        enableVerify.value = res.data["enabled_verify"];
      }
    })
    .catch((e) => {
      ElMessage.error("获取系统配置失败：" + e.message);
    });
});

const submit = (verifyData) => {
  if (action.value === "login") {
    doLogin(verifyData);
  } else if (action.value === "register") {
    doRegister(verifyData);
  }
};

// 登录操作
const submitLogin = () => {
  if (!data.value.username) {
    return ElMessage.error("请输入用户名");
  }
  if (!data.value.password) {
    return ElMessage.error("请输入密码");
  }
  if (enableVerify.value) {
    captchaRef.value.loadCaptcha();
    action.value = "login";
  } else {
    doLogin({});
  }
};

const doLogin = (verifyData) => {
  data.value.key = verifyData.key;
  data.value.dots = verifyData.dots;
  data.value.x = verifyData.x;
  httpPost("/api/user/login", data.value)
    .then((res) => {
      setUserToken(res.data.token);
      store.setIsLogin(true);
      ElMessage.success("登录成功！");
      emits("hide");
      emits("success");
    })
    .catch((e) => {
      ElMessage.error("登录失败，" + e.message);
    });
};

// 注册操作
const submitRegister = () => {
  if (activeName.value === "username" && data.value.username === "") {
    return ElMessage.error("请输入用户名");
  }

  if (activeName.value === "mobile" && !validateMobile(data.value.mobile)) {
    return ElMessage.error("请输入合法的手机号");
  }

  if (activeName.value === "email" && !validateEmail(data.value.email)) {
    return ElMessage.error("请输入合法的邮箱地址");
  }

  if (data.value.password.length < 8) {
    return ElMessage.error("密码的长度为8-16个字符");
  }
  if (data.value.repass !== data.value.password) {
    return ElMessage.error("两次输入密码不一致");
  }

  if ((activeName.value === "mobile" || activeName.value === "email") && data.value.code === "") {
    return ElMessage.error("请输入验证码");
  }
  if (enableVerify.value && activeName.value === "username") {
    captchaRef.value.loadCaptcha();
    action.value = "register";
  } else {
    doRegister({});
  }
};

const doRegister = (verifyData) => {
  data.value.key = verifyData.key;
  data.value.dots = verifyData.dots;
  data.value.x = verifyData.x;
  data.value.reg_way = activeName.value;
  httpPost("/api/user/register", data.value)
    .then((res) => {
      setUserToken(res.data.token);
      ElMessage.success({
        message: "注册成功!",
        onClose: () => {
          emits("hide");
          emits("success");
        },
        duration: 1000,
      });
    })
    .catch((e) => {
      ElMessage.error("注册失败，" + e.message);
    });
};
</script>

<style lang="stylus">
.login-dialog {
  border-radius 10px

  .el-tabs__nav {
    display flex
    width 100%
    justify-content space-between
  }


  .form {
    .block {
      margin-bottom 10px
    }

    .btn-row {
      display flex

      .login-btn {
        font-size 16px
        width 100%
      }

      .text {
        line-height 40px

        .el-tag {
          cursor pointer
        }
      }

      .forget {
        margin-left 10px
      }
    }

    .c-login {
      display flex
      .text {
        font-size 16px
        color #a1a1a1
        display: flex;
        align-items: center;
      }
      .login-type {
        display flex
        justify-content center

        .iconfont {
          font-size 18px
          background: #E9F1F6;
          padding: 8px;
          border-radius: 50%;
        }
        .iconfont.icon-wechat {
          color #0bc15f
        }
      }
    }

    .text {
      color var(--el-text-color-primary)
    }
  }

  .register-box {
    .wechat-card {
      text-align center
    }
  }

}
</style>
