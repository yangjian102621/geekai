<script lang="ts" setup>
import { reactive } from "vue";
import { useRoute } from "vue-router";
import { useAuthStore } from "@/stores/auth";
import useRequest from "@/composables/useRequest";

const route = useRoute();
const authStore = useAuthStore();

const [loginRequest, _, loading] = useRequest(authStore.login);
const formData = reactive({
  username: "",
  password: "",
});

async function handleSubmit({ errors, values }: any) {
  if (errors) return;
  await loginRequest({
    ...values,
    ...route.query,
  });
}
</script>
<template>
  <div class="bg">
    <div class="content">
      <!-- 左侧图片 -->
      <span class="left">
        <img src="/left-img.png" alt="" style="width: 468px" />
      </span>
      <!-- 表单 -->
      <div class="right-content">
        <div class="form-box">
          <div class="title">ChatGPT Plus Admin</div>
          <a-form
              ref="formRef"
              :model="formData"
              class="form"
              size="medium"
              auto-label-width
              @submit="handleSubmit"
          >
            <a-space direction="vertical" style="width: 100%">
              <a-form-item
                  field="username"
                  label="账号"
                  hide-label
                  hide-asterisk
                  :rules="[{ required: true, message: '请输入您的账号' }]"
              >
                <a-input
                    v-model="formData.username"
                    placeholder="请输入您的账号"
                    class="input"
                ></a-input>
              </a-form-item>
              <a-form-item
                  field="password"
                  label="密码"
                  hide-label
                  hide-asterisk
                  :rules="[{ required: true, message: '请输入您的密码' }]"
              >
                <a-input-password
                    v-model="formData.password"
                    placeholder="请输入您的密码"
                    class="input"
                >
                </a-input-password>
              </a-form-item>
            </a-space>
            <a-form-item hide-label>
              <a-button :disabled="loading" html-type="submit" long type="primary" class="sign-in-btn">
                登录
              </a-button>
            </a-form-item>
          </a-form>
        </div>
      </div>
    </div>
  </div>
</template>
<style lang="less" scoped>
.bg {
  width: 100%;
  height: 100vh;
  background: linear-gradient(133deg, #ffffff 0%, #dde8fe 100%);
  display: flex;
  justify-content: center;
  align-items: center;
  overflow: auto;
}

.content {
  width: 1080px;
  height: 557px;
  display: flex;
}

.left {
  width: 540px;
  height: 557px;
  background: #2670fe;
  border-radius: 16px 0 0 16px;
  display: flex;
  justify-content: center;
  align-items: center;
}

.right-content {
  width: 540px;
  height: 557px;
  background: #ffffff;
  border-radius: 0 16px 16px 0;
  display: flex;
  justify-content: center;
}

.title {
  font-size: 24px;
  font-weight: 800;
  color: #333333;
  line-height: 35px;
  letter-spacing: 1px;
  text-align: center;
  margin-bottom: 40px;
}

.form-box {
  display: flex;
  flex-direction: column;
  width: 438px;
  padding: 58px 0;
  height: 100%;
  box-sizing: border-box;
}

.form {
  flex: 1;
  justify-content: space-between;
}

.input {
  border-radius: 10px 10px 10px 10px;
  height: 60px;
  border: 2px solid #e5e6eb;
}

.captcha-image {
  //width: 100%;
  //height: 100%;
  cursor: pointer;
}

.sign-in-btn {
  height: 60px;
  font-weight: 500;
  line-height: 33px;
  font-size: 28px;
  border-radius: 10px 10px 10px 10px;
}
</style>
