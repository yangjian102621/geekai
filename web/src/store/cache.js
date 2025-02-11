import {httpGet} from "@/utils/http";
import Storage from "good-storage";
import {randString} from "@/utils/libs";

const userDataKey = "USER_INFO_CACHE_KEY"
const adminDataKey = "ADMIN_INFO_CACHE_KEY"
const systemInfoKey = "SYSTEM_INFO_CACHE_KEY"
const licenseInfoKey = "LICENSE_INFO_CACHE_KEY"
export function checkSession() {
    return new Promise((resolve, reject) => {
        httpGet('/api/user/session').then(res => {
            resolve(res.data)
        }).catch(e => {
            Storage.remove(userDataKey)
            reject(e)
        })
    })
}
export function checkAdminSession() {
    const item = Storage.get(adminDataKey) ?? {expire:0, data:null}
    if (item.expire > Date.now()) {
        return Promise.resolve(item.data)
    }
    return new Promise((resolve, reject) => {
        httpGet('/api/admin/session').then(res => {
            item.data = res.data
            item.expire = Date.now() + 1000 * 30
            Storage.set(adminDataKey, item)
            resolve(item.data)
        }).catch(e => {
            Storage.remove(adminDataKey)
            reject(e)
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
            item.expire = Date.now() + 1000 * 30
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
            item.expire = Date.now() + 1000 * 30
            Storage.set(licenseInfoKey, item)
            resolve(item.data)
        }).catch(err => {
            resolve(err)
        })
    })
}

export function getClientId() {
    let clientId = Storage.get('client_id')
    if (clientId) {
        return clientId
    }
    clientId = randString(42)
    Storage.set('client_id', clientId)
    return clientId
}