import {randString} from "@/utils/libs";
import Storage from "good-storage";
import {removeAdminInfo} from "@/store/cache";

/**
 * storage handler
 */

const UserTokenKey = import.meta.env.VITE_APP_KEY_PREFIX + "Authorization";
const AdminTokenKey = import.meta.env.VITE_APP_KEY_PREFIX + "Admin-Authorization"

export function getSessionId() {
    return randString(42)
}

export function getUserToken() {
    return Storage.get(UserTokenKey) ?? ""
}
export function setUserToken(token) {
    // 刷新 session 缓存
    Storage.set(UserTokenKey, token)
}

export function removeUserToken() {
    Storage.remove(UserTokenKey)
}

export function getAdminToken() {
    return Storage.get(AdminTokenKey) ?? ""
}

export function setAdminToken(token) {
    Storage.set(AdminTokenKey, token)
}

export function removeAdminToken() {
    Storage.remove(AdminTokenKey)
    removeAdminInfo()
}
