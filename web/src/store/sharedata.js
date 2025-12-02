import {defineStore} from "pinia";
import Storage from "good-storage";
import errorIcon from "@/assets/img/failed.png";
import loadingIcon from "@/assets/img/loading.gif";

let waterfallOptions = {
  // 唯一key值
  rowKey: "id",
  // 卡片之间的间隙
  gutter: 10,
  // 是否有周围的gutter
  hasAroundGutter: true,
  // 卡片在PC上的宽度
  width: 200,
  // 自定义行显示个数，主要用于对移动端的适配
  breakpoints: {
    3840: {
      // 4K下
      rowPerView: 8,
    },
    2560: {
      // 2K下
      rowPerView: 7,
    },
    1920: {
      // 2K下
      rowPerView: 6,
    },
    1600: {
      // 2K下
      rowPerView: 5,
    },
    1366: {
      // 2K下
      rowPerView: 4,
    },
    800: {
      // 当屏幕宽度小于等于800
      rowPerView: 3,
    },
    500: {
      // 当屏幕宽度小于等于500
      rowPerView: 2,
    },
  },
  // 动画效果
  animationEffect: "animate__fadeInUp",
  // 动画时间
  animationDuration: 1000,
  // 动画延迟
  animationDelay: 300,
  animationCancel: false,
  // 背景色
  backgroundColor: "",
  // imgSelector
  imgSelector: "img_thumb",
  // 是否跨域
  crossOrigin: true,
  // 加载配置
  loadProps: {
    loading: loadingIcon,
    error: errorIcon,
    ratioCalculator: (width, height) => {
      const minRatio = 3 / 4;
      const maxRatio = 4 / 3;
      const curRatio = height / width;
      if (curRatio < minRatio) {
        return minRatio;
      } else if (curRatio > maxRatio) {
        return maxRatio;
      } else {
        return curRatio;
      }
    },
  },
  // 是否懒加载
  lazyload: true,
  align: "center",
}

export const useSharedStore = defineStore("shared", {
  state: () => ({
    showLoginDialog: false,
    chatListStyle: Storage.get("chat_list_style", "chat"),
    chatStream: Storage.get("chat_stream", true),
    socket: { conn: null, handlers: {} },
    theme: Storage.get("theme", "light"),
    isLogin: false,
    chatListExtend: Storage.get("chat_list_extend", true),
    ttsModel: Storage.get("tts_model", ""),
    waterfallOptions,
  }),
  getters: {},
  actions: {
    setShowLoginDialog(value) {
      this.showLoginDialog = value;
    },
    setChatListStyle(value) {
      this.chatListStyle = value;
      Storage.set("chat_list_style", value);
    },
    setChatStream(value) {
      this.chatStream = value;
      Storage.set("chat_stream", value);
    },
    setSocket(value) {
      for (const key in this.socket.handlers) {
        this.setMessageHandler(value, this.socket.handlers[key]);
      }
      this.socket.conn = value;
    },
    setChatListExtend(value) {
      this.chatListExtend = value;
      Storage.set("chat_list_extend", value);
    },
    addMessageHandler(key, callback) {
      if (!this.socket.handlers[key]) {
        this.socket.handlers[key] = callback;
      }
      this.setMessageHandler(this.socket.conn, callback);
    },

    setMessageHandler(conn, callback) {
      if (!conn) {
        return;
      }
      conn.addEventListener("message", (event) => {
        try {
          if (event.data instanceof Blob) {
            const reader = new FileReader();
            reader.readAsText(event.data, "UTF-8");
            reader.onload = () => {
              callback(JSON.parse(String(reader.result)));
            };
          }
        } catch (e) {
          console.warn(e);
        }
      });
    },

    removeMessageHandler(key) {
      if (this.socket.conn && this.socket.conn.readyState === WebSocket.OPEN) {
        this.socket.conn.removeEventListener("message", this.socket.handlers[key]);
      }
      delete this.socket.handlers[key];
    },
    setTheme(theme) {
      this.theme = theme;
      document.documentElement.setAttribute("data-theme", theme); // 设置 HTML 的 data-theme 属性
      Storage.set("theme", theme);
    },
    setIsLogin(value) {
      this.isLogin = value;
    },

    setTtsModel(value) {
      this.ttsModel = value;
      Storage.set("tts_model", value);
    },
  },
});
