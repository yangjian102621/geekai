<template>
  <div class="mobile-image container">
    <van-tabs v-model:active="activeName" class="my-tab" animated sticky>
      <van-tab title="MJ" name="mj" v-if="activeMenu.mj">
        <image-mj />
      </van-tab>
      <van-tab title="SD" name="sd" v-if="activeMenu.sd">
        <image-sd />
      </van-tab>
      <van-tab title="DALL" name="dall" v-if="activeMenu.dall">
        <image-dall />
      </van-tab>
    </van-tabs>
  </div>
</template>

<script setup>
import { onMounted, ref } from "vue";
import ImageMj from "@/views/mobile/pages/ImageMj.vue";
import ImageSd from "@/views/mobile/pages/ImageSd.vue";
import ImageDall from "@/views/mobile/pages/ImageDall.vue";
import { httpGet } from "@/utils/http";

const activeName = ref("");
const menus = ref([]);
const activeMenu = ref({
  mj: false,
  sd: false,
  dall: false,
});

onMounted(() => {
  httpGet("/api/menu/list").then((res) => {
    menus.value = res.data;
    activeMenu.value = {
      mj: menus.value.some((item) => item.url === "/mj"),
      sd: menus.value.some((item) => item.url === "/sd"),
      dall: menus.value.some((item) => item.url === "/dalle"),
    };
  });
});
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
