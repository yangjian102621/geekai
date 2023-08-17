/* eslint-disable no-constant-condition */

import {randString} from "@/utils/libs";

/**
 * storage handler
 */

const SessionIDKey = 'SESSION_ID';

export function getSessionId() {
    let sessionId = sessionStorage.getItem(SessionIDKey)
    if (!sessionId) {
        sessionId = randString(42)
        setSessionId(sessionId)
    }
    return sessionId
}

export function removeLoginUser() {
    sessionStorage.removeItem(SessionIDKey)
}

export function setSessionId(sessionId) {
    sessionStorage.setItem(SessionIDKey, sessionId)
}
