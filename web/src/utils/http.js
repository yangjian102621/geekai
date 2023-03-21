import axios from 'axios'
import {getSessionId} from "@/utils/storage";

axios.defaults.timeout = 5000
axios.defaults.baseURL = process.env.VUE_APP_API_SECURE === true ? 'https://' + process.env.VUE_APP_API_HOST : 'http://' + process.env.VUE_APP_API_HOST
axios.defaults.withCredentials = true;
axios.defaults.headers.post['Content-Type'] = 'application/json'

// HTTP拦截器
axios.interceptors.request.use(
    config => {
        // set token
        config.headers['ChatGPT-Token'] = getSessionId();
        return config
    }, error => {
        return Promise.reject(error)
    })
axios.interceptors.response.use(
    response => {
        let data = response.data;
        if (data.code === 0) {
            return response
        } else {
            return Promise.reject(response.data)
        }
    }, error => {
        return Promise.reject(error)
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
