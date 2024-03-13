import {randString} from "@/utils/libs";
import Storage from "good-storage";

/**
 * storage handler
 */

const SessionIDKey = process.env.VUE_APP_KEY_PREFIX + 'SESSION_ID';
const UserTokenKey = process.env.VUE_APP_KEY_PREFIX + "Authorization";
const AdminTokenKey = process.env.VUE_APP_KEY_PREFIX + "Admin-Authorization"

export function getSessionId() {
    let sessionId = Storage.get(SessionIDKey)
    if (!sessionId) {
        sessionId = randString(42)
        setSessionId(sessionId)
    }
    return sessionId
}

export function removeSessionId() {
    Storage.remove(SessionIDKey)
}

export function setSessionId(sessionId) {
    Storage.set(SessionIDKey, sessionId)
}

export function getUserToken() {
    return Storage.get(UserTokenKey) ?? ""
}

export function setUserToken(token) {
    Storage.set(UserTokenKey, token)
}

export function removeUserToken() {
    Storage.remove(UserTokenKey)
}

export function getAdminToken() {
    return Storage.get(AdminTokenKey)
}

export function setAdminToken(token) {
    Storage.set(AdminTokenKey, token)
}

export function removeAdminToken() {
    Storage.remove(AdminTokenKey)
}
