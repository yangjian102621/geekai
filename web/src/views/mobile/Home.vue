<template>
  <van-config-provider :theme="theme">
    <div class="mobile-home">
      <router-view/>

      <van-tabbar route v-model="active">
        <van-tabbar-item to="/mobile/index" name="home" icon="home-o">首页</van-tabbar-item>
        <van-tabbar-item to="/mobile/chat" name="chat" icon="chat-o">对话</van-tabbar-item>
        <van-tabbar-item to="/mobile/image" name="image" icon="photo-o">绘图</van-tabbar-item>
        <van-tabbar-item to="/mobile/profile" name="profile" icon="user-o">我的
        </van-tabbar-item>
      </van-tabbar>

    </div>
  </van-config-provider>

</template>

<script setup>
import {ref, watch} from "vue";
import {useSharedStore} from "@/store/sharedata";

const active = ref('home')
const store = useSharedStore()
const theme = ref(store.mobileTheme)

watch(() => store.mobileTheme, (val) => {
  theme.value = val
})

</script>

<style lang="stylus">
@import '@/assets/iconfont/iconfont.css';
.mobile-home {
  .container {
    .van-nav-bar {
      position fixed
      width 100%
    }

    padding 0 10px
  }

}

// 黑色主题
.van-theme-dark body {
  background #1c1c1e
}

.van-nav-bar {
  position fixed
  width 100%
}
</style>