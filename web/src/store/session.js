/* eslint-disable no-constant-condition */

/**
 * storage handler
 */

const SessionUserKey = 'SESSION_ID';

export function getSessionId() {
    return sessionStorage.getItem(SessionUserKey)
}

export function removeLoginUser() {
    sessionStorage.removeItem(SessionUserKey)
}

export function setSessionId(sessionId) {
    sessionStorage.setItem(SessionUserKey, sessionId)
}
