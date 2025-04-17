import { defineStore } from "pinia";
import Storage from "good-storage";

export const useSharedStore = defineStore("shared", {
  state: () => ({
    showLoginDialog: false,
    chatListStyle: Storage.get("chat_list_style", "chat"),
    chatStream: Storage.get("chat_stream", true),
    socket: { conn: null, handlers: {} },
    theme: Storage.get("theme", "light"),
    isLogin: false,
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
  },
});
