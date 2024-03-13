import {defineStore} from 'pinia';
import Storage from "good-storage";

export const useSidebarStore = defineStore('sidebar', {
    state: () => {
        return {
            collapse: false
        };
    },
    getters: {},
    actions: {
        handleCollapse() {
            this.collapse = !this.collapse;
        }
    }
});

const MENU_STORE_KEY = "admin_menu_items"

export function getMenuItems() {
    return Storage.get(MENU_STORE_KEY)
}

export function setMenuItems(items) {
    return Storage.set(MENU_STORE_KEY, items)
}
