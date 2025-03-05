<template>
  <el-config-provider>
    <router-view/>
  </el-config-provider>
</template>

<script setup>
import {ElConfigProvider} from 'element-plus';
import {onMounted, ref, watch} from "vue";
import {checkSession, getClientId, getSystemInfo} from "@/store/cache";
import {isChrome, isMobile} from "@/utils/libs";
import {showMessageInfo} from "@/utils/dialog";
import {useSharedStore} from "@/store/sharedata";
import {getUserToken} from "@/store/session";

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

const store = useSharedStore()
onMounted(() => {
  // 获取系统参数
  getSystemInfo().then((res) => {
    const link = document.createElement('link')
    link.rel = 'shortcut icon'
    link.href = res.data.logo
    document.head.appendChild(link)
  })
  if (!isChrome() && !isMobile()) {
    showMessageInfo("建议使用 Chrome 浏览器以获得最佳体验。")
  }

  checkSession().then(() => {
    store.setIsLogin(true)
  }).catch(()=>{})
})

watch(() => store.isLogin, (val) => {
  if (val) {
    connect()
  }
})

const handler = ref(0)
// 初始化 websocket 连接
const connect = () => {
  let host = process.env.VUE_APP_WS_HOST
  if (host === '') {
    if (location.protocol === 'https:') {
      host = 'wss://' + location.host;
    } else {
      host = 'ws://' + location.host;
    }
  }
  const clientId = getClientId()
  const _socket = new WebSocket(host + `/api/ws?client_id=${clientId}`,["token",getUserToken()]);
  _socket.addEventListener('open', () => {
    console.log('WebSocket 已连接')
    handler.value = setInterval(() => {
      if (_socket.readyState === WebSocket.OPEN) {
        _socket.send(JSON.stringify({"type":"ping"}))
      }
    },5000)
  })
  _socket.addEventListener('close', () => {
    clearInterval(handler.value)
    connect()
  });
  store.setSocket(_socket)
}

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
      //max-height 80vh
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
