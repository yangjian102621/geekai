/* eslint-disable no-constant-condition */
import {dateFormat} from "@/utils/libs";

/**
 * storage handler
 */

const SessionUserKey = 'LOGIN_USER';
export const Global = {}

export function getSessionId() {
    const user = getLoginUser();
    return user ? user['session_id'] : '';
}

export function getLoginUser() {
    const value = sessionStorage.getItem(SessionUserKey);
    if (value) {
        return JSON.parse(value);
    } else {
        return null;
    }
}

export function setLoginUser(user) {
    sessionStorage.setItem(SessionUserKey, JSON.stringify(user))
}

export function getUserInfo() {
    const data = getLoginUser();
    if (data !== null) {
        const user = data["user"];
        user['active_time'] = dateFormat(user['active_time']);
        user['expired_time'] = dateFormat(user['expired_time']);
        return user;
    }
    return {}
}