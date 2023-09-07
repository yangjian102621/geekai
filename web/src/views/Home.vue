<template>
  <div class="home">
    <div class="inner">
      <h1>HI <br/> ChatGPT-PLUS</h1>
    </div>
  </div>
</template>

<script setup>
import {useRouter} from "vue-router";
import {checkSession} from "@/action/session";
import {isMobile} from "@/utils/libs";

const router = useRouter();
checkSession().then(() => {
  if (isMobile()) {
    router.push("/mobile")
  } else {
    router.push("/chat")
  }
}).catch(() => {
  router.push("/login")
})

</script>

<style lang="stylus" scoped>
.home {
  display: flex;
  justify-content: center;
  background-color: #282c34;
  height 100vh

  .inner {
    text-align center
    display flex
    justify-content center
    max-width 400px
    align-items center

    h1 {
      color: #202020;
      font-size: 55px;
      line-height 1.5
      font-weight: bold;
      text-shadow: -1px -1px 1px #111111, 2px 2px 1px #363636;
    }
  }

}
</style>
