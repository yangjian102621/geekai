<template>
  <div class="reset-pass"></div>
  <div class="flex-center loginPage">
    <div class="left">
      <div class="login-box">
        <AccountTop title="重置密码" />
        <div class="input-form">
          <el-form :model="form">
            <el-tabs v-model="form.type">
              <el-tab-pane label="手机号验证" name="mobile">
                <el-form-item>
                  <div class="form-title">手机号码</div>
                  <el-input
                    v-model="form.mobile"
                    size="large"
                    placeholder="请输入手机号"
                  />
                </el-form-item>
                <el-form-item>
                  <div class="form-title">验证码</div>
                  <div class="flex w100">
                    <el-input
                      v-model="form.code"
                      maxlength="6"
                      size="large"
                      placeholder="请输入验证码"
                      class="code-input"
                    />
                    <send-msg
                      size="large"
                      :receiver="form.mobile"
                      type="mobile"
                    />
                  </div>
                </el-form-item>
              </el-tab-pane>
              <el-tab-pane label="邮箱验证" name="email">
                <el-form-item>
                  <div class="form-title">邮箱</div>

                  <el-input
                    v-model="form.email"
                    placeholder="请输入邮箱"
                    size="large"
                  />
                </el-form-item>
                <el-form-item>
                  <div class="form-title">验证码</div>
                  <div class="flex w100">
                    <el-input
                      v-model="form.code"
                      maxlength="6"
                      size="large"
                      placeholder="请输入验证码"
                      class="code-input"
                    />
                    <send-msg
                      size="large"
                      :receiver="form.email"
                      type="email"
                    />
                  </div>
                </el-form-item>
              </el-tab-pane>
            </el-tabs>

            <el-form-item>
              <div class="form-title">新密码</div>

              <el-input
                v-model="form.password"
                type="password"
                placeholder="请输入新密码(8-16位)"
                size="large"
              />
            </el-form-item>
            <el-form-item>
              <div class="form-title">重复密码</div>

              <el-input
                v-model="form.repass"
                type="password"
                placeholder="请再次输入密码(8-16位)"
                size="large"
              />
            </el-form-item>
            <el-form-item>
              <el-button
                class="login-btn"
                size="large"
                type="primary"
                @click="save"
              >
                重置密码
              </el-button>
            </el-form-item>
          </el-form>
        </div>
      </div>
    </div>
    <account-bg />
  </div>
</template>

<script setup>
import { ref } from "vue";
import SendMsg from "@/components/SendMsg.vue";
import AccountTop from "@/components/AccountTop.vue";

import { ElMessage } from "element-plus";
import { httpPost } from "@/utils/http";
import AccountBg from "@/components/AccountBg.vue";
import { validateEmail, validateMobile } from "@/utils/validate";

const form = ref({
  mobile: "",
  email: "",
  type: "mobile",
  code: "",
  password: "",
  repass: ""
});

const save = () => {
  if (form.value.code === "") {
    return ElMessage.error("请输入验证码");
  }
  if (form.value.password.length < 8) {
    return ElMessage.error("密码长度必须大于8位");
  }
  if (form.value.repass !== form.value.password) {
    return ElMessage.error("两次输入密码不一致");
  }

  httpPost("/api/user/resetPass", form.value)
    .then(() => {
      ElMessage.success({
        message: "重置密码成功",
        duration: 1000
      });
    })
    .catch((e) => {
      ElMessage.error("重置密码失败：" + e.message);
    });
};
</script>

<style lang="stylus">
@import "../assets/css/login.styl"
::v-deep(.el-tabs__item.is-active, .el-tabs__item:hover){
    color: var(--common-text-color) !important;
  }
</style>
