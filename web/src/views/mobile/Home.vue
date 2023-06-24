<template>
  <div>
    <router-view/>

    <van-tabbar route v-model="active" @change="onChange">
      <van-tabbar-item to="/mobile/chat" name="home" icon="chat-o"></van-tabbar-item>
      <van-tabbar-item to="/mobile/setting" name="setting" icon="setting-o"></van-tabbar-item>
      <van-tabbar-item to="/mobile/profile" name="profile" icon="user-o"></van-tabbar-item>
    </van-tabbar>

  </div>
</template>

<script setup>
import {ref} from "vue";

const active = ref('home');
const onChange = (index) => {
  console.log(index)
  // showToast(`标签 ${index}`);
}
</script>

<style lang="stylus">
.van-toast--fail {
  background #fef0f0
  color #f56c6c
}
</style>