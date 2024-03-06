<script lang="ts" setup>
import { reactive } from "vue";
import { useRoute } from "vue-router";
import { IconUser, IconLock } from "@arco-design/web-vue/es/icon";
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
  <div class="login-wrapper public-bg">
    <div class="login-container">
      <div class="login-box">
        <AForm :model="formData" size="large" @submit="handleSubmit">
          <h1 class="title">ChatGPT Plus Admin</h1>
          <AFormItem
            hide-label
            field="username"
            :rules="[{ required: true, message: '请输入用户名' }]"
          >
            <AInput v-model="formData.username" placeholder="用户名">
              <template #prefix>
                <IconUser />
              </template>
            </AInput>
          </AFormItem>
          <AFormItem
            hide-label
            field="password"
            :rules="[{ required: true, message: '请输入密码' }]"
          >
            <AInputPassword v-model="formData.password" placeholder="密码">
              <template #prefix>
                <IconLock />
              </template>
            </AInputPassword>
          </AFormItem>
          <AFormItem hide-label>
            <AButton type="primary" long html-type="submit" :loading="loading">
              登录
            </AButton>
          </AFormItem>
        </AForm>
      </div>
    </div>
  </div>
</template>
<style lang="less" scoped>
.login-wrapper {
  .login-container {
    display: flex;
    width: 1000px;
    height: 500px;
    align-items: center;
    justify-content: right;
    background-color: #ffffff;
    background-image: url("/login-content.png");
    background-size: contain;
    background-repeat: no-repeat;
    background-position: left;
    border-radius: 40px;
    .login-box {
      display: flex;
      width: 50%;
      height: 100%;
      padding: 20px;
      border-radius: 40px;
      align-items: center;
      justify-content: center;
      box-sizing: border-box;
      .title {
        text-align: center;
      }
    }
  }
}
</style>
