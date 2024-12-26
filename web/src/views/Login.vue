<template>
  <div class="flex-center loginPage">
    <div class="left">
      <div class="login-box">
        <AccountTop>
          <template #default>
            <div class="wechatLog flex-center" v-if="wechatLoginURL !== ''">
              <a :href="wechatLoginURL" @click="setRoute(router.currentRoute.value.path)"> <i class="iconfont icon-wechat"></i>使用微信登录 </a>
            </div>
          </template>
        </AccountTop>

        <div class="input-form">
          <el-form ref="ruleFormRef" :model="ruleForm" :rules="rules">
            <el-form-item label="" prop="username">
              <div class="form-title">账号</div>
              <el-input v-model="ruleForm.username" size="large" placeholder="请输入账号" @keyup="handleKeyup" />
            </el-form-item>
            <el-form-item label="" prop="password">
              <div class="flex-between w100">
                <div class="form-title">密码</div>
                <div class="form-forget text-color-primary" @click="router.push('/resetpassword')">忘记密码？</div>
              </div>

              <el-input size="large" v-model="ruleForm.password" placeholder="请输入密码" show-password autocomplete="off" @keyup="handleKeyup" />
            </el-form-item>
            <el-form-item>
              <el-button class="login-btn" size="large" type="primary" @click="login">登录</el-button>
            </el-form-item>
          </el-form>
        </div>
      </div>
    </div>
    <account-bg />

    <captcha v-if="enableVerify" @success="doLogin" ref="captchaRef" />
  </div>
</template>

<script setup>
import { onMounted, ref, reactive } from "vue";
import { httpGet, httpPost } from "@/utils/http";
import { useRouter } from "vue-router";
import AccountBg from "@/components/AccountBg.vue";
import { isMobile } from "@/utils/libs";
import { checkSession, getLicenseInfo, getSystemInfo } from "@/store/cache";
import { setUserToken } from "@/store/session";
import { showMessageError } from "@/utils/dialog";
import { setRoute } from "@/store/system";
import { useSharedStore } from "@/store/sharedata";

import AccountTop from "@/components/AccountTop.vue";
import Captcha from "@/components/Captcha.vue";

const router = useRouter();
const title = ref("Geek-AI");
const username = ref(process.env.VUE_APP_USER);
const password = ref(process.env.VUE_APP_PASS);

const logo = ref("");
const licenseConfig = ref({});
const wechatLoginURL = ref("");
const enableVerify = ref(false);
const captchaRef = ref(null);
const ruleFormRef = ref(null);
const ruleForm = reactive({
  username: process.env.VUE_APP_USER,
  password: process.env.VUE_APP_PASS,
});
const rules = {
  username: [{ required: true, trigger: "blur", message: "请输入账号" }],
  password: [{ required: true, trigger: "blur", message: "请输入密码" }],
};
onMounted(() => {
  // 获取系统配置
  getSystemInfo()
    .then((res) => {
      logo.value = res.data.logo;
      title.value = res.data.title;
      enableVerify.value = res.data["enabled_verify"];
    })
    .catch((e) => {
      showMessageError("获取系统配置失败：" + e.message);
    });

  getLicenseInfo()
    .then((res) => {
      licenseConfig.value = res.data;
    })
    .catch((e) => {
      showMessageError("获取 License 配置：" + e.message);
    });

  checkSession()
    .then(() => {
      if (isMobile()) {
        router.push("/mobile");
      } else {
        router.push("/chat");
      }
    })
    .catch(() => {});

  const returnURL = `${location.protocol}//${location.host}/login/callback?action=login`;
  httpGet("/api/user/clogin?return_url=" + returnURL)
    .then((res) => {
      wechatLoginURL.value = res.data.url;
    })
    .catch((e) => {
      console.error(e);
    });
});

const handleKeyup = (e) => {
  if (e.key === "Enter") {
    login();
  }
};

const login = async function () {
  await ruleFormRef.value.validate(async (valid) => {
    if (valid) {
      if (enableVerify.value) {
        captchaRef.value.loadCaptcha();
      } else {
        doLogin({});
      }
    }
  });
};

const store = useSharedStore();
const doLogin = (verifyData) => {
  httpPost("/api/user/login", {
    username: username.value.trim(),
    password: password.value.trim(),
    key: verifyData.key,
    dots: verifyData.dots,
    x: verifyData.x,
  })
    .then((res) => {
      setUserToken(res.data.token);
      store.setIsLogin(true);
      if (isMobile()) {
        router.push("/mobile");
      } else {
        router.push("/chat");
      }
    })
    .catch((e) => {
      showMessageError("登录失败，" + e.message);
    });
};
</script>

<style lang="stylus" scoped>
@import "@/assets/css/login.styl"
</style>
