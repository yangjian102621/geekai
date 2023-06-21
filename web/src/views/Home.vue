<template>
  <div class="home" :style="{ height: winHeight + 'px' }">
    <h1>{{ title }}</h1>
  </div>
</template>

<script setup>
import {ref} from "vue"
import {useRouter} from "vue-router"; // 导入useRouter函数

const title = ref("HI, ChatGPT PLUS!");
const winHeight = ref(window.innerHeight)
const router = useRouter();
// onMounted(() => {
//   if (isMobile()) {
//     router.push("mobile");
//   } else {
//     router.push("chat");
//   }
// })


</script>

<style lang="stylus" scoped>
.home {
  display: flex;
  justify-content: center;
  align-items: center;
  color: #202020;
  background-color: #282c34;

  h1 {
    font-size: 300%;
    font-weight: bold;
    letter-spacing: 0.1em;
    text-shadow: -1px -1px 1px #111111, 2px 2px 1px #363636;
  }
}
</style>
