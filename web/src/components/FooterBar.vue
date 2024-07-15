<template>
  <div class="foot-container">
    <div class="footer">
      <div v-if="license.de_copy">{{copyRight}}</div>
      <div v-else>
        <span>{{copyRight}}</span>
        <el-link type="primary" :href="gitURL" target="_blank" style="--el-link-text-color:#ffffff">
          {{ title }} -
          {{ version }}
        </el-link>
      </div>
    </div>
  </div>
</template>
<script setup>

import {ref} from "vue";
import {httpGet} from "@/utils/http";
import {showMessageError} from "@/utils/dialog";

const title = ref("")
const version = ref(process.env.VUE_APP_VERSION)
const gitURL = ref(process.env.VUE_APP_GIT_URL)
const copyRight = ref('')
const license = ref({})


// 获取系统配置
httpGet("/api/config/get?key=system").then(res => {
  title.value = res.data.title??process.env.VUE_APP_TITLE
  copyRight.value = res.data.copyright??'极客学长 © 2023 - '+new Date().getFullYear()+' All rights reserved.'
}).catch(e => {
  showMessageError("获取系统配置失败：" + e.message)
})

httpGet("/api/config/license").then(res => {
  license.value = res.data
}).catch(e => {
  showMessageError("获取 License 失败：" + e.message)
})
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