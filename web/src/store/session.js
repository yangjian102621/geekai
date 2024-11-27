import {randString} from "@/utils/libs";
import Storage from "good-storage";
import {removeAdminInfo, removeUserInfo} from "@/store/cache";

/**
 * storage handler
 */

const UserTokenKey = process.env.VUE_APP_KEY_PREFIX + "Authorization";
const AdminTokenKey = process.env.VUE_APP_KEY_PREFIX + "Admin-Authorization"

export function getSessionId() {
    return randString(42)
}

export function getUserToken() {
    return Storage.get(UserTokenKey) ?? ""
}

export function setUserToken(token) {
    Storage.set(UserTokenKey, token)
}

export function removeUserToken() {
    Storage.remove(UserTokenKey)
    removeUserInfo()
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
