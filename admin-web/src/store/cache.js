import { httpGet } from '@/utils/http'
import Storage from 'good-storage'
import { randString } from '@/utils/libs'

const userDataKey = 'USER_INFO_CACHE_KEY'
const adminDataKey = 'ADMIN_INFO_CACHE_KEY'
const systemInfoKey = 'SYSTEM_INFO_CACHE_KEY'
const menusKey = 'MENUS_CACHE_KEY'

export function getMenus() {
  const item = Storage.get(menusKey) ?? { expire: 0, data: null }
  if (item.expire > Date.now()) {
    return Promise.resolve(item.data)
  }
  return new Promise((resolve, reject) => {
    httpGet('/api/menu/list/all')
      .then((res) => {
        const menus = {}
        for (const menu of res.data) {
          menus[menu.url] = menu
        }
        item.data = menus
        item.expire = Date.now() + 1000 * 3
        Storage.set(menusKey, item)
        resolve(item.data)
      })
      .catch((err) => {
        reject(err)
      })
  })
}

export function checkSession() {
  const item = Storage.get(userDataKey) ?? { expire: 0, data: null }
  if (item.expire > Date.now()) {
    return Promise.resolve(item.data)
  }
  return new Promise((resolve, reject) => {
    httpGet('/api/user/session')
      .then((res) => {
        item.data = res.data
        item.expire = Date.now() + 1000 * 3
        Storage.set(userDataKey, item)
        resolve(item.data)
      })
      .catch((e) => {
        Storage.remove(userDataKey)
        reject(e)
      })
  })
}
export function checkAdminSession() {
  const item = Storage.get(adminDataKey) ?? { expire: 0, data: null }
  if (item.expire > Date.now()) {
    return Promise.resolve(item.data)
  }
  return new Promise((resolve, reject) => {
    httpGet('/api/admin/session')
      .then((res) => {
        item.data = res.data
        item.expire = Date.now() + 1000 * 3
        Storage.set(adminDataKey, item)
        resolve(item.data)
      })
      .catch((e) => {
        Storage.remove(adminDataKey)
        reject(e)
      })
  })
}

export function removeAdminInfo() {
  Storage.remove(adminDataKey)
}

export function getSystemInfo() {
  const item = Storage.get(systemInfoKey) ?? { expire: 0, data: null }
  if (item.expire > Date.now()) {
    return Promise.resolve(item.data)
  }
  return new Promise((resolve, reject) => {
    httpGet('/api/config/get?key=system')
      .then((res) => {
        item.data = res
        item.expire = Date.now() + 1000 * 3
        Storage.set(systemInfoKey, item)
        resolve(item.data)
      })
      .catch((err) => {
        reject(err)
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
