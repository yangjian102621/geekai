<template>
  <div class="banner">
    <h1 class="banner-title">ChatPuls-V3 智能助理</h1>
  </div>
  <van-config-provider theme="dark">
    <div class="mobile-home">
      <router-view/>
      <van-tabbar
          fixed
          route
          v-model="active"
          active-color="#F5F5F5"
          inactive-color="#2CC995"
          class="my-tabbar"
      >
        <van-tabbar-item to="/mobile/chat/list" name="home" icon="chat-o">对话</van-tabbar-item>
        <van-tabbar-item to="/mobile/imageSd" name="imageSd" icon="photo-o">绘图</van-tabbar-item>
        <van-tabbar-item to="/mobile/apps" name="apps" icon="apps-o">应用</van-tabbar-item>
        <van-tabbar-item to="/mobile/profile" name="profile" icon="user-o">我的</van-tabbar-item>
        <van-tabbar-item to="/mobile/invitation" name="invitation" icon="share-o">分享</van-tabbar-item>
      </van-tabbar>
    </div>
  </van-config-provider>
</template>

<script setup>
import {ref} from "vue";
import {getMobileTheme} from "@/store/system";
import {useRouter} from "vue-router";
import {isMobile} from "@/utils/libs";
import {checkSession} from "@/action/session";

const router = useRouter()
checkSession().then(() => {
  if (!isMobile()) {
    router.replace('/chat')
  }
}).catch(() => {
  router.push('/login')
})

const active = ref('home')
const onChange = (index) => {
  console.log(index)
  // showToast(`标签 ${index}`);
}

</script>

<style lang="stylus" scoped>
@import "@/assets/css/mobile/home.css"
</style>
