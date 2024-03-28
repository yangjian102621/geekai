<template>
  <el-config-provider :locale="zhCn">
    <router-view/>
  </el-config-provider>
</template>

<script setup>
import {ElConfigProvider} from 'element-plus';
import zhCn from 'element-plus/es/locale/lang/zh-cn';

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
  font-family: Helvetica Neue, Helvetica, PingFang SC, Hiragino Sans GB, Microsoft YaHei, Arial, sans-serif
  -webkit-font-smoothing: antialiased;
  text-rendering: optimizeLegibility;
}

.el-overlay-dialog {
  display flex
  justify-content center
  align-items center
  overflow hidden

  .el-dialog {
    margin 0;

    .el-dialog__body {
      max-height 90vh
    }
  }
}

/* 省略显示 */
.ellipsis {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

</style>
