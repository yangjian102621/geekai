// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import { getAdminToken, getUserToken, removeAdminToken, removeUserToken } from '@/store/session'
import axios from 'axios'

// Blob 数据读取和解析的辅助函数
export async function parseBlobResponse(blob) {
  try {
    // 检查 Blob 的类型
    if (blob.type && blob.type.includes('application/json')) {
      // 如果是 JSON 类型，直接解析
      const text = await blob.text()
      return JSON.parse(text)
    } else {
      // 如果不是 JSON 类型，尝试解析为文本
      const text = await blob.text()
      try {
        return JSON.parse(text)
      } catch (e) {
        // 如果解析 JSON 失败，返回文本内容
        return { message: text, rawData: text }
      }
    }
  } catch (error) {
    console.error('解析 Blob 响应失败:', error)
    return { message: '解析响应数据失败', error: error.message }
  }
}

axios.defaults.timeout = 180000
// axios.defaults.baseURL = process.env.VUE_APP_API_HOST
axios.defaults.withCredentials = true
//axios.defaults.headers.post['Content-Type'] = 'application/json'

// HTTP拦截器
axios.interceptors.request.use(
  (config) => {
    // set token
    config.headers['Authorization'] = getUserToken()
    config.headers['Admin-Authorization'] = getAdminToken()
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)
axios.interceptors.response.use(
  (response) => {
    return response
  },
  async (error) => {
    if (error.response.status === 401) {
      if (error.response.request.responseURL.indexOf('/api/admin') !== -1) {
        removeAdminToken()
      } else {
        removeUserToken()
      }
      error.response.data.message = '请先登录'
      return Promise.reject(error.response.data)
    }

    if (error.response.status === 400) {
      let errorMessage = error.response.data.message
      if (!errorMessage) {
        const parsedData = await parseBlobResponse(error.response.data)
        errorMessage = parsedData.message
      }
      return Promise.reject(new Error(errorMessage))
    } else {
      return Promise.reject(error)
    }
  }
)

// send a http get request
export function httpGet(url, params = {}) {
  return new Promise((resolve, reject) => {
    axios
      .get(url, {
        params: params,
      })
      .then((response) => {
        resolve(response.data)
      })
      .catch((err) => {
        reject(err)
      })
  })
}

// send a http post request
export function httpPost(url, data = {}, options = {}) {
  return new Promise((resolve, reject) => {
    axios
      .post(url, data, options)
      .then((response) => {
        resolve(response.data)
      })
      .catch((err) => {
        reject(err)
      })
  })
}

export function httpDownload(url) {
  return new Promise((resolve, reject) => {
    axios({
      method: 'GET',
      url: url,
      responseType: 'blob', // 将响应类型设置为 `blob`
    })
      .then((response) => {
        resolve(response)
      })
      .catch((err) => {
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
      responseType: 'blob', // 将响应类型设置为 `blob`
    })
      .then((response) => {
        resolve(response)
      })
      .catch((err) => {
        reject(err)
      })
  })
}
