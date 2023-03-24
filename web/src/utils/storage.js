/* eslint-disable no-constant-condition */
/**
 * storage handler
 */

const SessionIdKey = 'ChatGPT_SESSION_ID';
export const Global = {}

export function getSessionId() {
    return sessionStorage.getItem(SessionIdKey)
}

export function setSessionId(value) {
    sessionStorage.setItem(SessionIdKey, value)
}