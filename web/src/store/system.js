// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import Storage from "good-storage";

const MOBILE_THEME = process.env.VUE_APP_KEY_PREFIX + "MOBILE_THEME"
const ADMIN_THEME = process.env.VUE_APP_KEY_PREFIX + "ADMIN_THEME"

export function getMobileTheme() {
    return Storage.get(MOBILE_THEME) ? Storage.get(MOBILE_THEME) : 'light'
}

export function setMobileTheme(theme) {
    Storage.set(MOBILE_THEME, theme)
}

export function getAdminTheme() {
    return Storage.get(ADMIN_THEME) ? Storage.get(ADMIN_THEME) : 'light'
}

export function setAdminTheme(theme) {
    Storage.set(ADMIN_THEME, theme)
}