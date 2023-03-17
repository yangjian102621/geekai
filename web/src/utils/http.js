import axios from 'axios'
import JSONBigInt from 'json-bigint'
import qs from 'qs'

import { ElMessageBox } from 'element-plus'

axios.defaults.timeout = 5000
axios.defaults.baseURL = process.env.VUE_APP_API_SECURE === true ? 'https://' + process.env.VUE_APP_API_HOST : 'http://' + process.env.VUE_APP_API_HOST
axios.defaults.withCredentials = true
axios.defaults.headers.post['Content-Type'] = 'application/x-www-form-urlencoded'
axios.defaults.transformResponse = [(data, headers) => {
  if (headers['content-type'].indexOf('application/json') !== -1) {
    try {
      data = JSONBigInt.parse(data)
    } catch (e) { /* Ignore */ }
  }
  return data
}]


// HTTP拦截器
axios.interceptors.request.use(
  config => {
    // set session-name
    config.headers['Session-Name'] = "xwebssh-sess-token"
    return config
  }, error => {
    return Promise.reject(error)
  })
axios.interceptors.response.use(
  response => {
    if (response.data.code == 0) {
      return response
    } else {
      return Promise.reject(response.data)
    }
  }, error => {
    if (error.response.status === 401) {
      ElMessageBox.alert('您未登录或者登录已退出，请先登录再操作。', '登录提醒', {
        confirmButtonText: '确定',
        callback: () => {
          // TODO: goto login page
        },
      })
    }

    if (error.response.status === 400) {
      return Promise.reject(error.response.data)
    }
    return Promise.reject(error)
  })


// send a http get request
export function httpGet (url, params = {}) {
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
export function httpPost (url, data = {}, options = {}) {
  return new Promise((resolve, reject) => {
    axios.post(url, qs.stringify(data), options).then(response => {
      resolve(response.data)
    }).catch(err => {
      reject(err)
    })
  })
}
