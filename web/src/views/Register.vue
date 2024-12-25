<template>
  <div>
    <div class="flex-center loginPage">
      <div class="left" v-if="enableRegister">
        <div class="login-box">
          <AccountTop title="注册" />
          <div class="input-form">
            <el-form :model="data" class="form">
              <el-tabs v-model="activeName">
                <el-tab-pane label="手机注册" name="mobile" v-if="enableMobile">
                  <el-form-item>
                    <div class="form-title">手机号码</div>
                    <el-input placeholder="请输入手机号码" size="large" v-model="data.mobile" maxlength="11" autocomplete="off"> </el-input>
                  </el-form-item>
                  <el-form-item>
                    <div class="form-title">验证码</div>
                    <div class="flex w100">
                      <el-input placeholder="请输入验证码" size="large" maxlength="30" class="code-input" v-model="data.code" autocomplete="off"> </el-input>

                      <send-msg size="large" :receiver="data.mobile" type="mobile" />
                    </div>
                  </el-form-item>
                </el-tab-pane>
                <el-tab-pane label="邮箱注册" name="email" v-if="enableEmail">
                  <el-form-item class="block">
                    <div class="form-title">邮箱</div>
                    <el-input placeholder="请输入邮箱地址" size="large" v-model="data.email" autocomplete="off"> </el-input>
                  </el-form-item>
                  <el-form-item class="block">
                    <div class="form-title">验证码</div>
                    <div class="flex w100">
                      <el-input placeholder="请输入验证码" size="large" maxlength="30" class="code-input" v-model="data.code" autocomplete="off"> </el-input>

                      <send-msg size="large" :receiver="data.email" type="email" />
                    </div>
                  </el-form-item>
                </el-tab-pane>
                <el-tab-pane label="用户名注册" name="username" v-if="enableUser">
                  <el-form-item class="block">
                    <div class="form-title">用户名</div>

                    <el-input placeholder="请输入用户名" size="large" v-model="data.username" autocomplete="off"> </el-input>
                  </el-form-item>
                </el-tab-pane>
              </el-tabs>

              <el-form-item class="block">
                <div class="form-title">密码</div>

                <el-input placeholder="请输入密码(8-16位)" maxlength="16" size="large" v-model="data.password" show-password autocomplete="off"> </el-input>
              </el-form-item>

              <el-form-item class="block">
                <div class="form-title">重复密码</div>

                <el-input placeholder="请再次输入密码(8-16位)" size="large" maxlength="16" v-model="data.repass" show-password autocomplete="off"> </el-input>
              </el-form-item>

              <el-form-item class="block">
                <div class="form-title">邀请码</div>

                <el-input placeholder="请输入邀请码(可选)" size="large" v-model="data.invite_code" autocomplete="off"> </el-input>
              </el-form-item>

              <el-row class="btn-row" :gutter="20">
                <el-col :span="24">
                  <el-button class="login-btn" type="primary" size="large" @click="submitRegister">注册</el-button>
                </el-col>
              </el-row>
            </el-form>
          </div>
        </div>
      </div>
      <div class="tip-result left" v-else>
        <el-result icon="error" title="注册功能已关闭">
          <template #sub-title>
            <p>抱歉，系统已关闭注册功能，请联系管理员添加账号！</p>
            <div class="wechat-card">
              <el-image :src="wxImg" />
            </div>
          </template>
        </el-result>
      </div>
      <captcha v-if="enableVerify" @success="doSubmitRegister" ref="captchaRef" />
      <account-bg />
    </div>
  </div>
</template>

<script setup>
import { ref } from "vue";
import AccountTop from "@/components/AccountTop.vue";
import AccountBg from "@/components/AccountBg.vue";

import { httpGet, httpPost } from "@/utils/http";
import { ElMessage } from "element-plus";
import { useRouter } from "vue-router";

import SendMsg from "@/components/SendMsg.vue";
import { arrayContains, isMobile } from "@/utils/libs";
import { setUserToken } from "@/store/session";
import { validateEmail, validateMobile } from "@/utils/validate";
import { showMessageError, showMessageOK } from "@/utils/dialog";
import { getLicenseInfo, getSystemInfo } from "@/store/cache";
import Captcha from "@/components/Captcha.vue";

const router = useRouter();
const title = ref("");
const logo = ref("");
const data = ref({
  username: "",
  mobile: "",
  email: "",
  password: "",
  code: "",
  repass: "",
  invite_code: router.currentRoute.value.query["invite_code"],
});

const enableMobile = ref(false);
const enableEmail = ref(false);
const enableUser = ref(false);
const enableRegister = ref(true);
const activeName = ref("mobile");
const wxImg = ref("/images/wx.png");
const licenseConfig = ref({});
const enableVerify = ref(false);
const captchaRef = ref(null);

// 记录邀请码点击次数
if (data.value.invite_code) {
  httpGet("/api/invite/hits", { code: data.value.invite_code });
}

getSystemInfo()
  .then((res) => {
    if (res.data) {
      title.value = res.data.title;
      logo.value = res.data.logo;
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

getLicenseInfo()
  .then((res) => {
    licenseConfig.value = res.data;
  })
  .catch((e) => {
    showMessageError("获取 License 配置：" + e.message);
  });

// 注册操作
const submitRegister = () => {
  if (activeName.value === "username" && data.value.username === "") {
    return showMessageError("请输入用户名");
  }

  if (activeName.value === "mobile" && !validateMobile(data.value.mobile)) {
    return showMessageError("请输入合法的手机号");
  }

  if (activeName.value === "email" && !validateEmail(data.value.email)) {
    return showMessageError("请输入合法的邮箱地址");
  }

  if (data.value.password.length < 8) {
    return showMessageError("密码的长度为8-16个字符");
  }
  if (data.value.repass !== data.value.password) {
    return showMessageError("两次输入密码不一致");
  }

  if ((activeName.value === "mobile" || activeName.value === "email") && data.value.code === "") {
    return showMessageError("请输入验证码");
  }

  // 如果是用户名和密码登录，那么需要加载验证码
  if (enableVerify.value && activeName.value === "username") {
    captchaRef.value.loadCaptcha();
  } else {
    doSubmitRegister({});
  }
};

const doSubmitRegister = (verifyData) => {
  data.value.key = verifyData.key;
  data.value.dots = verifyData.dots;
  data.value.x = verifyData.x;
  data.value.reg_way = activeName.value;
  httpPost("/api/user/register", data.value)
    .then((res) => {
      setUserToken(res.data.token);
      showMessageOK("注册成功，即将跳转到对话主界面...");
      if (isMobile()) {
        router.push("/mobile/index");
      } else {
        router.push("/chat");
      }
    })
    .catch((e) => {
      showMessageError("注册失败，" + e.message);
    });
};
</script>

<style lang="stylus" scoped>
@import "@/assets/css/login.styl"
  :deep(.back){
    margin-bottom: 10px;
  }

  :deep(.orline){
    margin-bottom: 10px;
  }
  .wechat-card {
    margin-top: 20px

  }
</style>
