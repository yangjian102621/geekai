import {defineStore} from 'pinia';
import Storage from 'good-storage'

export const useSharedStore = defineStore('shared', {
    state: () => ({
        showLoginDialog: false,
        chatListStyle: Storage.get("chat_list_style","chat")
    }),
    getters: {},
    actions: {
        setShowLoginDialog(value) {
            this.showLoginDialog = value;
        },
        setChatListStyle(value) {
            this.chatListStyle = value;
            Storage.set("chat_list_style", value);
        }
    }
});
