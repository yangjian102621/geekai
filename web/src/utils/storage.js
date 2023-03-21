/* eslint-disable no-constant-condition */
/**
 * storage handler
 */
import Storage from 'good-storage'

const SessionIdKey = 'ChatGPT_SESSION_ID';
export const Global = {}

export function getSessionId() {
    return Storage.get(SessionIdKey)
}

export function setSessionId(value) {
    Storage.set(SessionIdKey, value)
}