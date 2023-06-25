import Storage from 'good-storage'

const CHAT_CONFIG_KEY = "chat_config"

export function getChatConfig() {
    return Storage.get(CHAT_CONFIG_KEY)
}

export function setChatConfig(chatConfig) {
    Storage.set(CHAT_CONFIG_KEY, chatConfig)
}