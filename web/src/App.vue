<template>
  <router-view></router-view>
</template>

<script>
import {defineComponent} from 'vue'

export default defineComponent({
  setup() {
    // TODO: 初始化操作
  },
})
</script>


<style lang="stylus">
html, body {
  height: 100%;
  margin: 0;
  padding: 0;
}

#app {
  margin: 0 !important;
  padding: 0 !important;
  font-family: 'Helvetica Neue', Helvetica, 'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei', '微软雅黑', Arial, sans-serif;
}
</style>
