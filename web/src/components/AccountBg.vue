<template>
  <div class="right flex-center">
    <div class="logo">
      <el-image :src="logo" alt="" style="max-width: 300px; max-height: 300px" class="rounded-full" />
    </div>
    <div>welcome</div>
    <footer-bar />
  </div>
</template>

<script setup>
import FooterBar from "@/components/FooterBar.vue";
import { getSystemInfo } from "@/store/cache";
import { ref } from "vue";

const logo = ref("");
const title = ref("");

getSystemInfo()
  .then((res) => {
    logo.value = res.data.logo;
    title.value = res.data.title;
  })
  .catch((err) => {
    console.log(err);
    logo.value = "/images/logo.png";
    title.value = "Geek-AI";
  });
</script>

<style lang="stylus" scoped>
.right{
  font-size: 40px
  font-weight: bold
  color:#fff
  flex-direction: column
  background-image url("~@/assets/img/login-bg.png")
  background-size cover
  background-position center
  width: 50%;
  min-height: 100vh
  max-height: 100vh
  background-repeat: no-repeat;
  position: relative;
  overflow: hidden;
  z-index: 1;
  :deep(.foot-container){
    position: absolute;
    bottom: 20px;
    width: 100%;
    background: none;
    color: var(--sm-txt);
    font-size: 12px;
    text-align: center;

    .footer{
      a,
      span{
        color: var(--text-fff)
      }
    }
  }
}
.logo{
  margin-bottom: 26px;
  width: 200px
  height: 200px
  background: #fff
  border-radius: 50%
  img{
    width: 100%;
    object-fit: cover;
    height: 100%;
  }
}
</style>
