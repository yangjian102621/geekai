import {defineStore} from 'pinia';
import Storage from 'good-storage'

export const useSharedStore = defineStore('shared', {
    state: () => ({
        showLoginDialog: false,
        chatListStyle: Storage.get("chat_list_style","chat"),
        chatStream: Storage.get("chat_stream",true),
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
        }
    }
});
