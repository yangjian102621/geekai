<template>
  <div class="power-log" :style="{ height: winHeight + 'px' }">
    <div class="inner">
      <h2>消费日志</h2>
    </div>
  </div>
</template>

<script setup>
import {ref} from "vue"

const winHeight = ref(window.innerHeight)
</script>

<style lang="stylus" scoped>
.power-log {
  display: flex;
  justify-content: center;
  background-color: #282c34;
  color #ffffff

  .inner {
    display flex
    justify-content left
    width 100%
    padding 0 20px 20px 20px

  }

}
</style>
