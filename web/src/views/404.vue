<template>
  <div class="page-404" :style="{ height: winHeight + 'px' }">
    <div class="inner">
      <h1>404</h1>
      <h2>饶了地球一圈，还是没有找到您要的页面</h2>
    </div>
  </div>
</template>

<script setup>
import {ref} from "vue"

const winHeight = ref(window.innerHeight)
</script>

<style lang="stylus" scoped>
.page-404 {
  display: flex;
  justify-content: center;
  background-color: #282c34;

  .inner {
    text-align center

    h1 {
      color: #202020;
      font-size: 120px;
      font-weight: bold;
      letter-spacing: 0.1em;
      text-shadow: -1px -1px 1px #111111, 2px 2px 1px #363636;
    }

    h2 {
      color #ffffff;
      font-weight: bold;
    }
  }

}
</style>
