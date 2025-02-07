<template>
  <el-config-provider>
    <router-view />
  </el-config-provider>
</template>

<script setup>
import { ElConfigProvider } from "element-plus";
import { onMounted, ref, watch } from "vue";
import { checkSession, getClientId, getSystemInfo } from "@/store/cache";
import { isChrome, isMobile } from "@/utils/libs";
import { showMessageInfo } from "@/utils/dialog";
import { useSharedStore } from "@/store/sharedata";
import { getUserToken } from "@/store/session";

const debounce = (fn, delay) => {
  let timer;
  return (...args) => {
    if (timer) {
      clearTimeout(timer);
    }
    timer = setTimeout(() => {
      fn(...args);
    }, delay);
  };
};

const _ResizeObserver = window.ResizeObserver;
window.ResizeObserver = class ResizeObserver extends _ResizeObserver {
  constructor(callback) {
    callback = debounce(callback, 200);
    super(callback);
  }
};

const store = useSharedStore();
onMounted(() => {
  // 获取系统参数
  getSystemInfo().then((res) => {
    const link = document.createElement("link");
    link.rel = "shortcut icon";
    link.href = res.data.logo;
    document.head.appendChild(link);
  });
  if (!isChrome() && !isMobile()) {
    showMessageInfo("建议使用 Chrome 浏览器以获得最佳体验。");
  }

  checkSession()
    .then(() => {
      store.setIsLogin(true);
    })
    .catch(() => {});

  // 设置主题
  document.documentElement.setAttribute("data-theme", store.theme);
});

watch(
  () => store.isLogin,
  (val) => {
    if (val) {
      connect();
    }
  }
);

const handler = ref(0);
// 初始化 websocket 连接
const connect = () => {
  let host = process.env.VUE_APP_WS_HOST;
  if (host === "") {
    if (location.protocol === "https:") {
      host = "wss://" + location.host;
    } else {
      host = "ws://" + location.host;
    }
  }
  const clientId = getClientId();
  const _socket = new WebSocket(host + `/api/ws?client_id=${clientId}`, ["token", getUserToken()]);
  _socket.addEventListener("open", () => {
    console.log("WebSocket 已连接");
    handler.value = setInterval(() => {
      if (_socket.readyState === WebSocket.OPEN) {
        _socket.send(JSON.stringify({ type: "ping" }));
      }
    }, 5000);
  });
  _socket.addEventListener("close", () => {
    clearInterval(handler.value);
    connect();
  });
  store.setSocket(_socket);
};

// 打印 banner
const banner = `

  .oooooo.                        oooo              .o.       ooooo 
 d8P'  'Y8b                        888             .888.       888
888            .ooooo.   .ooooo.   888  oooo      .8"888.      888  
888           d88'  88b d88'  88b  888 .8P'      .8'  888.     888  
888     ooooo 888ooo888 888ooo888  888888.      .88ooo8888.    888  
'88.    .88'  888    .o 888    .o  888  88b.   .8'      888.   888  
  Y8bood8P'    Y8bod8P'  Y8bod8P' o888o o888o o88o     o8888o o888o

  `;
console.log("%c" + banner + "", "color: purple;font-size: 18px;");

console.log("%c感谢大家为 GeekAI 做出的卓越贡献！", "color: green;font-size: 40px;font-family: '微软雅黑';");
console.log(
  "%c项目源码：https://github.com/yangjian102621/geekai %c 您的 star 对我们非常重要！",
  "color: green;font-size: 20px;font-family: '微软雅黑';",
  "color: red;font-size: 20px;font-family: '微软雅黑';"
);

console.log("%c 愿你出走半生，归来仍是少年！大奉武夫许七安，前来凿阵！", "color: #7c39ed;font-size: 18px;font-family: '微软雅黑';");
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

  h1 { font-size: 2em; } /* 通常是 2em */
  h2 { font-size: 1.5em; } /* 通常是 1.5em */
  h3 { font-size: 1.17em; } /* 通常是 1.17em */
  h4 { font-size: 1em; } /* 通常是 1em */
  h5 { font-size: 0.83em; } /* 通常是 0.83em */
  h6 { font-size: 0.67em; } /* 通常是 0.67em */

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

@import '@/assets/iconfont/iconfont.css'
</style>
