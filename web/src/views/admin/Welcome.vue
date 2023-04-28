<template>
  <div class="welcome" :style="{ height: winHeight + 'px' }">
    <h1>Chat-Plus 控制台</h1>
  </div>
</template>

<script>
import {defineComponent} from "vue"

export default defineComponent({
  name: 'XWelcome',
  data() {
    return {
      winHeight: window.innerHeight,
    }
  },
  mounted: function () {
    const that = this
    window.addEventListener("resize", function () {
      that.winHeight = window.innerHeight
    })
  }
})
</script>

<style lang="stylus" scoped>
.welcome {
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
