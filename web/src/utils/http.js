// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import axios from 'axios'
import {getAdminToken, getUserToken, removeAdminToken, removeUserToken} from "@/store/session";

axios.defaults.timeout = 180000
axios.defaults.baseURL = import.meta.env.VITE_APP_API_HOST
axios.defaults.withCredentials = true;
//axios.defaults.headers.post['Content-Type'] = 'application/json'

// HTTP拦截器
axios.interceptors.request.use(
    config => {
        // set token
        config.headers['Authorization'] = getUserToken();
        config.headers['Admin-Authorization'] = getAdminToken();
        return config
    }, error => {
        return Promise.reject(error)
    })
axios.interceptors.response.use(
    response => {
        return response
    }, error => {
        if (error.response.status === 401) {
            if (error.response.request.responseURL.indexOf("/api/admin") !== -1) {
                removeAdminToken()
            } else {
                console.log("FUCK")
                removeUserToken()
            }
            error.response.data.message = "请先登录"
            return Promise.reject(error.response.data)
        }
        if (error.response.status === 400) {
            return Promise.reject(new Error(error.response.data.message))
        } else {
            return Promise.reject(error)
        }
    })


// send a http get request
export function httpGet(url, params = {}) {
    return new Promise((resolve, reject) => {
        axios.get(url, {
            params: params
        }).then(response => {
            resolve(response.data)
        }).catch(err => {
            reject(err)
        })
    })
}


// send a http post request
export function httpPost(url, data = {}, options = {}) {
    return new Promise((resolve, reject) => {
        axios.post(url, data, options).then(response => {
            resolve(response.data)
        }).catch(err => {
            reject(err)
        })
    })
}

export function httpDownload(url) {
    return new Promise((resolve, reject) => {
        axios({
            method: 'GET',
            url: url,
            responseType: 'blob' // 将响应类型设置为 `blob`
        }).then(response => {
            resolve(response)
        }).catch(err => {
            reject(err)
        })
    })
}

export function httpPostDownload(url, data) {
    return new Promise((resolve, reject) => {
        axios({
            method: 'POST',
            url: url,
            data: data,
            responseType: 'blob' // 将响应类型设置为 `blob`
        }).then(response => {
            resolve(response)
        }).catch(err => {
            reject(err)
        })
    })
}