import {httpGet} from "@/utils/http";
import Storage from "good-storage";
import {showMessageError} from "@/utils/dialog";

const userDataKey = "USER_INFO_CACHE_KEY"
const adminDataKey = "ADMIN_INFO_CACHE_KEY"
const systemInfoKey = "SYSTEM_INFO_CACHE_KEY"
const licenseInfoKey = "LICENSE_INFO_CACHE_KEY"
export function checkSession() {
    const item = Storage.get(userDataKey) ?? {expire:0, data:null}
    if (item.expire > Date.now()) {
        return Promise.resolve(item.data)
    }

    return new Promise((resolve, reject) => {
        httpGet('/api/user/session').then(res => {
            item.data = res.data
            // cache expires after 5 minutes
            item.expire = Date.now() + 1000 * 60 * 5
            Storage.set(userDataKey, item)
            resolve(item.data)
        }).catch(err => {
            Storage.remove(userDataKey)
            reject(err)
        })
    })
}

export function removeUserInfo() {
    Storage.remove(userDataKey)
}

export function checkAdminSession() {
    const item = Storage.get(adminDataKey) ?? {expire:0, data:null}
    if (item.expire > Date.now()) {
        return Promise.resolve(item.data)
    }
    return new Promise((resolve, reject) => {
        httpGet('/api/admin/session').then(res => {
            item.data = res.data
            // cache expires after 10 minutes
            item.expire = Date.now() + 1000 * 60 * 10
            Storage.set(adminDataKey, item)
            resolve(item.data)
        }).catch(err => {
            reject(err)
        })
    })
}

export function removeAdminInfo() {
    Storage.remove(adminDataKey)
}

export function getSystemInfo() {
    const item = Storage.get(systemInfoKey) ?? {expire:0, data:null}
    if (item.expire > Date.now()) {
        return Promise.resolve(item.data)
    }
    return new Promise((resolve, reject) => {
        httpGet('/api/config/get?key=system').then(res => {
            item.data = res
            // cache expires after 10 minutes
            item.expire = Date.now() + 1000 * 60 * 10
            Storage.set(systemInfoKey, item)
            resolve(item.data)
        }).catch(err => {
            reject(err)
        })
    })
}

export function getLicenseInfo() {
    const item = Storage.get(licenseInfoKey) ?? {expire:0, data:null}
    if (item.expire > Date.now()) {
        return Promise.resolve(item.data)
    }

    return new Promise((resolve, reject) => {
        httpGet('/api/config/license').then(res => {
            item.data = res
            // cache expires after 10 minutes
            item.expire = Date.now() + 1000 * 60 * 10
            Storage.set(licenseInfoKey, item)
            resolve(item.data)
        }).catch(err => {
            reject(err)
        })
    })
}