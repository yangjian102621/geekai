<template>
  <div class="mobile-image container">
    <van-tabs v-model:active="activeName" class="my-tab" animated sticky>
      <van-tab title="MJ" name="mj">
        <image-mj/>
      </van-tab>
      <van-tab title="SD" name="sd">
        <image-sd/>
      </van-tab>
      <van-tab title="DALL" name="dall">
        <image-dall />
      </van-tab>
    </van-tabs>
  </div>
</template>

<script setup>
import {ref} from "vue";
import ImageMj from "@/views/mobile/pages/ImageMj.vue";
import ImageSd from "@/views/mobile/pages/ImageSd.vue";
import ImageDall from "@/views/mobile/pages/ImageDall.vue";

const activeName = ref("mj")
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