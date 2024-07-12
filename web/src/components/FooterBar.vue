<template>
  <div class="foot-container">
    <div class="footer">
      Powered by {{ author }} @
      <el-link type="primary" :href="gitURL" target="_blank" style="--el-link-text-color:#ffffff">
        {{ title }} -
        {{ version }}
      </el-link>
    </div>
  </div>
</template>
<script setup>

import {ref} from "vue";

const title = ref(process.env.VUE_APP_TITLE)
const version = ref(process.env.VUE_APP_VERSION)
const gitURL = ref(process.env.VUE_APP_GIT_URL)
const author = ref('极客学长')
</script>

<style scoped lang="stylus">
.foot-container {
  position: fixed;
  left: 0;
  bottom: 0;
  width: 100%;
  display flex;
  justify-content center

  .footer {
    max-width 400px;
    text-align center;
    font-size 14px;
    padding 20px;
    width 100%

    .el-link {
      color #409eff
    }
  }
}

</style>