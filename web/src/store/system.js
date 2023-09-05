import Storage from "good-storage";

const MOBILE_THEME = process.env.VUE_APP_KEY_PREFIX + "MOBILE_THEME"

export function getMobileTheme() {
    return Storage.get(MOBILE_THEME) ? Storage.get(MOBILE_THEME) : 'light'
}

export function setMobileTheme(theme) {
    Storage.set(MOBILE_THEME, theme)
}