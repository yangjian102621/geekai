// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

/**
 * Util lib functions
 */
import { showConfirmDialog } from 'vant'
import { httpGet } from '@/utils/http'

// generate a random string
export function randString(length) {
  const str = '0123456789abcdefghijklmnopqrstuvwxyz'
  const size = str.length
  let buf = []
  for (let i = 0; i < length; i++) {
    const rand = Math.random() * size
    buf.push(str.charAt(rand))
  }
  return buf.join('')
}

export function UUID() {
  let d = new Date().getTime()
  return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function (c) {
    const r = (d + Math.random() * 16) % 16 | 0
    d = Math.floor(d / 16)
    return (c === 'x' ? r : (r & 0x3) | 0x8).toString(16)
  })
}

// 判断是否是移动设备
export function isMobile() {
  const userAgent = navigator.userAgent
  const mobileRegex =
    /Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini|Mobile|mobile|CriOS/i
  return mobileRegex.test(userAgent)
}

// 格式化日期
export function dateFormat(timestamp, format) {
  if (!timestamp) {
    return ''
  } else if (timestamp < 9680917502) {
    timestamp = timestamp * 1000
  }
  let year, month, day, HH, mm, ss
  let time = new Date(timestamp)
  let timeDate
  year = time.getFullYear() // 年
  month = time.getMonth() + 1 // 月
  day = time.getDate() // 日
  HH = time.getHours() // 时
  mm = time.getMinutes() // 分
  ss = time.getSeconds() // 秒

  month = month < 10 ? '0' + month : month
  day = day < 10 ? '0' + day : day
  HH = HH < 10 ? '0' + HH : HH // 时
  mm = mm < 10 ? '0' + mm : mm // 分
  ss = ss < 10 ? '0' + ss : ss // 秒

  switch (format) {
    case 'yyyy':
      timeDate = String(year)
      break
    case 'yyyy-MM':
      timeDate = year + '-' + month
      break
    case 'yyyy-MM-dd':
      timeDate = year + '-' + month + '-' + day
      break
    case 'yyyy/MM/dd':
      timeDate = year + '/' + month + '/' + day
      break
    case 'yyyy-MM-dd HH:mm:ss':
      timeDate = year + '-' + month + '-' + day + ' ' + HH + ':' + mm + ':' + ss
      break
    case 'HH:mm:ss':
      timeDate = HH + ':' + mm + ':' + ss
      break
    case 'MM':
      timeDate = String(month)
      break
    default:
      timeDate = year + '-' + month + '-' + day + ' ' + HH + ':' + mm + ':' + ss
      break
  }
  return timeDate
}

export function formatTime(time) {
  const minutes = Math.floor(time / 60)
  const seconds = Math.floor(time % 60)
  return `${minutes}:${seconds < 10 ? '0' : ''}${seconds}`
}

// 判断数组中是否包含某个元素
export function arrayContains(array, value, compare) {
  if (!array) {
    return false
  }

  if (typeof compare !== 'function') {
    compare = function (v1, v2) {
      return v1 === v2
    }
  }
  for (let i = 0; i < array.length; i++) {
    if (compare(array[i], value)) {
      return true
    }
  }
  return false
}

// 删除数组中指定的元素
export function removeArrayItem(array, value, compare) {
  if (typeof compare !== 'function') {
    compare = function (v1, v2) {
      return v1 === v2
    }
  }
  for (let i = 0; i < array.length; i++) {
    if (compare(array[i], value)) {
      array.splice(i, 1)
      break
    }
  }
  return array
}

// 渲染输入的换行符
export function renderInputText(text) {
  const replaceRegex = /(\n\r|\r\n|\r|\n)/g
  text = text || ''
  return text.replace(replaceRegex, '<br/>')
}

// 拷贝对象
export function copyObj(origin) {
  return JSON.parse(JSON.stringify(origin))
}

export function disabledDate(time) {
  return time.getTime() < Date.now()
}

// 字符串截取
export function substr(str, length) {
  let result = ''
  let count = 0

  for (let i = 0; i < str.length; i++) {
    const char = str.charAt(i)
    const charCode = str.charCodeAt(i)

    // 判断字符是否为中文字符
    if (charCode >= 0x4e00 && charCode <= 0x9fff) {
      // 中文字符算两个字符
      count += 2
    } else {
      count++
    }

    if (count <= length) {
      result += char
    } else {
      result += ' ...'
      break
    }
  }

  return result
}

export function isImage(url) {
  const expr = /\.(jpg|jpeg|png|gif|bmp|svg)$/i
  return expr.test(url)
}

export function processContent(content) {
  if (!content) {
    return ''
  }

  // 如果是图片链接地址，则直接替换成图片标签
  const linkRegex = /(https?:\/\/\S+)/g
  const links = content.match(linkRegex)
  if (links) {
    for (let link of links) {
      if (isImage(link)) {
        const index = content.indexOf(link)
        if (content.substring(index - 1, 2) !== ']') {
          content = content.replace(link, '\n![](' + link + ')\n')
        }
      }
    }
  }
  // 处理推理标签
  if (content.includes('<think>')) {
    content = content.replace(/<think>(.*?)<\/think>/gs, (match, content) => {
      if (content.length > 10) {
        return `<blockquote>\n\n${content}</blockquote>`
      }
      return ''
    })
    content = content.replace(/<think>(.*?)$/gs, (match, content) => {
      if (content.length > 10) {
        return `<blockquote>${content}</blockquote>`
      }
      return ''
    })
  }

  // 支持 \[ 公式标签
  content = content.replace(/\\\[/g, '$$').replace(/\\\]/g, '$$')
  content = content.replace(/\\\(\\boxed\{(\d+)\}\\\)/g, '<span class="boxed">$1</span>')
  return content
}

export function processPrompt(prompt) {
  return prompt.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;')
}

// 判断是否为微信浏览器
export function isWeChatBrowser() {
  return /MicroMessenger/i.test(navigator.userAgent)
}

export function showLoginDialog(router) {
  showConfirmDialog({
    title: '登录',
    message: '此操作需要登录才能进行，前往登录？',
  })
    .then(() => {
      router.push('/login')
    })
    .catch(() => {
      // on cancel
    })
}

export const replaceImg = (img) => {
  if (img.startsWith('http')) {
    return img
  }
  return `${location.protocol}//${location.host}${img}`
}

// 判断是否 google 浏览器
export function isChrome() {
  const userAgent = navigator.userAgent.toLowerCase()
  return /chrome/.test(userAgent) && !/edg/.test(userAgent)
}

// 格式化日期时间
export function formatDateTime(timestamp, format = 'yyyy-MM-dd HH:mm:ss') {
  return dateFormat(timestamp, format)
}

export function isWechat() {
  const ua = navigator.userAgent.toLowerCase()
  return ua.indexOf('micromessenger') !== -1
}

// 默认缩略图模板（本地存储格式）
const DEFAULT_THUMB_TEMPLATE = '?imageView2/4/w/{width}/h/{height}/q/75'
const THUMB_TEMPLATE_KEY = 'GEEKAI_THUMB_TEMPLATE'

/**
 * 初始化缩略图模板（从后端加载并缓存到localStorage）
 * @returns {Promise<void>}
 */
export async function initThumbTemplate() {
  try {
    const res = await httpGet('/api/config/oss/thumb')
    if (res.data && res.data.template) {
      const template = res.data.template || DEFAULT_THUMB_TEMPLATE
      localStorage.setItem(THUMB_TEMPLATE_KEY, template)
      return
    }
  } catch (error) {
    console.warn('获取缩略图模板失败，使用默认模板:', error)
  }

  // 获取失败时使用默认模板并缓存
  localStorage.setItem(THUMB_TEMPLATE_KEY, DEFAULT_THUMB_TEMPLATE)
}

/**
 * 获取缓存的缩略图模板
 * @returns {string} 缩略图模板
 */
function getCachedThumbTemplate() {
  const cached = localStorage.getItem(THUMB_TEMPLATE_KEY)
  return cached || DEFAULT_THUMB_TEMPLATE
}

/**
 * 生成缩略图URL
 * @param {string} originalURL - 原始图片URL
 * @param {number} width - 缩略图宽度
 * @param {number} height - 缩略图高度，0表示自适应
 * @param {string} template - 缩略图模板（可选），如果不提供则从localStorage获取
 * @returns {string} 缩略图URL
 */
export function getThumbURL(originalURL, width = 300, height = 0, template = null) {
  if (!originalURL) {
    return originalURL
  }

  // 如果没有提供模板，从localStorage获取
  if (!template) {
    template = getCachedThumbTemplate()
  }

  // 如果模板为空字符串，返回原图（表示不支持缩略图）
  if (template === '') {
    return originalURL
  }

  // 如果模板为空，使用默认模板（降级方案）
  if (!template) {
    template = DEFAULT_THUMB_TEMPLATE
  }

  // 替换变量
  const thumbURL = template.replace(/{width}/g, width).replace(/{height}/g, height)

  // 拼接原始URL和缩略图参数
  return originalURL + thumbURL
}
