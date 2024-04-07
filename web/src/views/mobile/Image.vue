<template>
  <div class="mobile-image container">
    <van-tabs v-model:active="activeName" class="my-tab" animated sticky>
      <van-tab title="MidJourney" name="mj">
        <image-mj/>
      </van-tab>
      <van-tab title="Stable-Diffusion" name="sd">
        <image-sd/>
      </van-tab>
      <van-tab title="DALL-E" name="dall">
        <van-empty description="功能正在开发中"/>
      </van-tab>
    </van-tabs>
  </div>
</template>

<script setup>
import {ref} from "vue";
import ImageMj from "@/views/mobile/ImageMj.vue";
import ImageSd from "@/views/mobile/ImageSd.vue";

const activeName = ref("sd")
</script>

<style lang="stylus">
.mobile-image {
  .my-tab {
    .van-tab__panel {
      padding 10px
    }
  }
}
</style>