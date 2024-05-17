import {defineStore} from 'pinia';

export const useSharedStore = defineStore('shared', {
    state: () => ({
        showLoginDialog: false
    }),
    getters: {},
    actions: {
        setShowLoginDialog(value) {
            this.showLoginDialog = value;
        }
    }
});
