/* eslint-disable no-constant-condition */

/**
 * storage handler
 */

const SessionUserKey = 'LOGIN_USER';

export function getSessionId() {
    const user = getLoginUser();
    return user ? user['session_id'] : '';
}

export function removeLoginUser() {
    sessionStorage.removeItem(SessionUserKey)
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
