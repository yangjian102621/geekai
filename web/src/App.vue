<template>
  <el-config-provider>
    <router-view/>
  </el-config-provider>
</template>

<script setup>
import {ElConfigProvider} from 'element-plus';
import {onMounted} from "vue";
import {getSystemInfo} from "@/store/cache";
import {isChrome, isMobile} from "@/utils/libs";
import {showMessageInfo} from "@/utils/dialog";

const debounce = (fn, delay) => {
  let timer
  return (...args) => {
    if (timer) {
      clearTimeout(timer)
    }
    timer = setTimeout(() => {
      fn(...args)
    }, delay)
  }
}

const _ResizeObserver = window.ResizeObserver;
window.ResizeObserver = class ResizeObserver extends _ResizeObserver {
  constructor(callback) {
    callback = debounce(callback, 200);
    super(callback);
  }
}

onMounted(() => {
  getSystemInfo().then((res) => {
    const link = document.createElement('link')
    link.rel = 'shortcut icon'
    link.href = res.data.logo
    document.head.appendChild(link)
  })
  if (!isChrome() && !isMobile()) {
    showMessageInfo("检测到您使用的浏览器不是 Chrome，可能会导致部分功能无法正常使用，建议使用 Chrome 浏览器。")
  }
})
</script>


<style lang="stylus">
html, body {
  margin: 0;
  padding: 0;
}

#app {
  margin: 0 !important;
  padding: 0 !important;
  font-family: Helvetica Neue, Helvetica, PingFang SC, Hiragino Sans GB, Microsoft YaHei, Arial, sans-serif
  -webkit-font-smoothing: antialiased;
  text-rendering: optimizeLegibility;

  --primary-color: #21aa93
}

.el-overlay-dialog {
  display flex
  justify-content center
  align-items center
  overflow hidden

  .el-dialog {
    margin 0;

    .el-dialog__body {
      max-height 80vh
      overflow-y auto
    }
  }
}

/* 省略显示 */
.ellipsis {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.van-toast--fail {
  background #fef0f0
  color #f56c6c
}

.van-toast--success {
  background #D6FBCC
  color #07C160
}

</style>
