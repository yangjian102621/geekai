import {defineStore} from 'pinia';
import Storage from 'good-storage'

export const useSharedStore = defineStore('shared', {
    state: () => ({
        showLoginDialog: false,
        chatListStyle: Storage.get("chat_list_style","chat"),
        chatStream: Storage.get("chat_stream",true),
        socket: WebSocket,
        messageHandlers:{},
        mobileTheme: Storage.get("mobile_theme", "light"),
        adminTheme: Storage.get("admin_theme", "light"),
        isLogin: false
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
            this.socket = value;
        },
        addMessageHandler(key, callback) {
            if (!this.messageHandlers[key]) {
                this.messageHandlers[key] = callback;
                this.setMessageHandler(callback)
            }
        },
        setMessageHandler(callback) {
            if (this.socket instanceof WebSocket && this.socket.readyState === WebSocket.OPEN) {
                this.socket.addEventListener('message', (event) => {
                    try {
                        if (event.data instanceof Blob) {
                            const reader = new FileReader();
                            reader.readAsText(event.data, "UTF-8");
                            reader.onload = () => {
                                callback(JSON.parse(String(reader.result)))
                            }
                        }
                    } catch (e) {
                        console.warn(e)
                    }
                })
            } else {
                setTimeout(() => {
                    this.setMessageHandler(callback)
                }, 1000)
            }
        },
        removeMessageHandler(key) {
            if (this.socket.readyState === WebSocket.OPEN) {
                this.socket.removeEventListener('message', this.messageHandlers[key])
            }
            delete this.messageHandlers[key]
        },
        setMobileTheme(theme) {
            this.mobileTheme = theme
            Storage.set("mobile_theme", theme)
        },
        setAdminTheme(theme) {
            this.adminTheme = theme
            Storage.set("admin_theme", theme)
        },
        setIsLogin(value) {
            this.isLogin = value
        }
    },
});
