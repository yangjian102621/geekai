<template>
  <div class="login flex w-full flex-col place-content-center h-lvh">
    <el-image src="/images/logo.png" class="w-1/2 mx-auto logo" />
    <div class="title text-center text-3xl font-bold mt-8">{{ title }}</div>
    <div class="w-full p-8">
      <login-dialog @success="loginSuccess" />
    </div>
  </div>
</template>

<script setup>
import LoginDialog from "@/components/LoginDialog.vue";
import { getSystemInfo } from "@/store/cache";
import { useRouter } from "vue-router";
import { ref, onMounted } from "vue";

const router = useRouter();
const title = ref("登录");

const loginSuccess = () => {
  router.back();
};

onMounted(() => {
  getSystemInfo().then((res) => {
    title.value = res.data.title;
  });
});
</script>

<style scoped lang="stylus">
.login {
  background: var(--theme-bg);
  transition: all 0.3s ease;

  .logo {
    background: #ffffff;
    border-radius: 50%;
  }

  .title {
    color: var(--text-theme-color);
  }
}
</style>
