<template>
  <div class="login flex w-full flex-col place-content-center h-lvh">
    <!-- 返回首页链接 -->
    <div class="back-home">
      <el-button @click="goHome" type="text" class="back-btn">
        <i class="iconfont icon-home"></i>
        返回首页
      </el-button>
    </div>

    <el-image :src="logo" class="w-1/2 mx-auto logo" />
    <div class="title text-center text-3xl font-bold mt-8">{{ title }}</div>
    <div class="w-full p-8">
      <login-dialog @success="loginSuccess" />
    </div>
  </div>
</template>

<script setup>
import LoginDialog from '@/components/LoginDialog.vue'
import { getSystemInfo } from '@/store/cache'
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const title = ref('登录')
const logo = ref('')

const loginSuccess = () => {
  router.back()
}

const goHome = () => {
  router.push('/mobile/index')
}

onMounted(() => {
  getSystemInfo().then((res) => {
    title.value = res.data.title
    logo.value = res.data.logo
  })
})
</script>

<style scoped lang="scss">
.login {
  background: var(--theme-bg);
  transition: all 0.3s ease;
  position: relative;

  .back-home {
    position: absolute;
    top: 20px;
    left: 20px;
    z-index: 10;

    .back-btn {
      color: var(--text-theme-color);
      font-size: 14px;

      .iconfont {
        margin-right: 4px;
      }
    }
  }

  .logo {
    background: #ffffff;
    border-radius: 50%;
  }

  .title {
    color: var(--text-theme-color);
  }
}
</style>
