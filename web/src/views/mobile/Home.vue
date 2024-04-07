<template>
  <van-config-provider :theme="getMobileTheme()">
    <div class="mobile-home">
      <router-view/>

      <van-tabbar route v-model="active" @change="onChange">
        <van-tabbar-item to="/mobile/chat" name="home" icon="chat-o">对话</van-tabbar-item>
        <van-tabbar-item to="/mobile/image" name="image" icon="photo-o">绘图</van-tabbar-item>
        <van-tabbar-item to="/mobile/img-wall" name="apps" icon="apps-o">广场</van-tabbar-item>
        <van-tabbar-item to="/mobile/profile" name="profile" icon="user-o">我的</van-tabbar-item>
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
}

</script>

<style lang="stylus">
@import '@/assets/iconfont/iconfont.css';
.mobile-home {
  .container {
    .van-nav-bar {
      position fixed
      width 100%
    }

    .content {
      padding 46px 10px 60px 10px
    }
  }

}

// 黑色主题
.van-theme-dark body {
  background #1c1c1e
}

.van-toast--fail {
  background #fef0f0
  color #f56c6c
}

.van-toast--success {
  background #D6FBCC
  color #07C160
}

.van-nav-bar {
  position fixed
  width 100%
}
</style>