import {httpGet} from "@/utils/http";

export function checkSession() {
    return new Promise((resolve, reject) => {
        httpGet('/api/user/session').then(res => {
            resolve(res.data)
        }).catch(err => {
            reject(err)
        })
    })
}

export function checkAdminSession() {
    return new Promise((resolve, reject) => {
        httpGet('/api/admin/session').then(res => {
            resolve(res)
        }).catch(err => {
            reject(err)
        })
    })
}