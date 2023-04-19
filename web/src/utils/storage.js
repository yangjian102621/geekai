/* eslint-disable no-constant-condition */
import {dateFormat} from "@/utils/libs";
import Storage from 'good-storage'

/**
 * storage handler
 */

const SessionUserKey = 'LOGIN_USER';
const ChatHistoryKey = 'CHAT_HISTORY';
const ChatListKey = 'CHAT_LIST';

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

// 追加历史记录
export function appendChatHistory(chatId, message) {
    let history = Storage.get(ChatHistoryKey);
    if (!history) {
        history = {};
    }
    if (!history[chatId]) {
        history[chatId] = [message];
    } else {
        history[chatId].push(message);
    }
    Storage.set(ChatHistoryKey, history);
}

export function clearChatHistory() {
    Storage.remove(ChatHistoryKey);
    Storage.remove(ChatListKey);
}

// 获取指定会话的历史记录
export function getChatHistory(chatId) {
    const history = Storage.get(ChatHistoryKey);
    if (!history) {
        return null;
    }

    return history[chatId] ? history[chatId] : null;
}

export function getChatList() {
    return Storage.get(ChatListKey);
}

export function getChat(chatId) {
    let chatList = Storage.get(ChatListKey);
    if (!chatList) {
        return null;
    }

    return chatList[chatId] ? chatList[chatId] : null;
}

export function setChat(chat) {
    let chatList = Storage.get(ChatListKey);
    if (!chatList) {
        chatList = {};
    }

    chatList[chat.id] = chat;
    Storage.set(ChatListKey, chatList);
}

export function removeChat(chatId) {
    const chatList = Storage.get(ChatListKey);
    if (chatList) {
        delete chatList[chatId];
        Storage.set(ChatListKey, chatList);
    }
}
