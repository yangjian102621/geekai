// * 门户站点（原 React web）的 API 与下载链接服务，移植到 Vue 后复用全局 axios 拦截器。

import axios from 'axios'
import { httpGet, httpPost } from '@/utils/http'
import { getUserToken, removeUserToken } from '@/store/session'

/* ---------------- 认证 ----------------
   会话请求使用独立实例：401 时仅清除本地 token 并视为未登录，
   不走全局拦截器的统一报错逻辑（与原 React 站点行为一致）。 */
export const siteApi = axios.create({
  baseURL: '/api',
  timeout: 10000,
  withCredentials: true,
  headers: { Accept: 'application/json' },
})

siteApi.interceptors.request.use((config) => {
  const token = getUserToken()
  if (token) config.headers.Authorization = token
  return config
})

export async function fetchSession() {
  try {
    const response = await siteApi.get('/user/session')
    return response.data?.data ?? null
  } catch (error) {
    if (error?.response?.status === 401) {
      removeUserToken()
      return null
    }
    throw error
  }
}

export async function logout() {
  try {
    await siteApi.get('/user/logout')
  } catch (e) {
    return
  }
}

/* ---------------- 邀请 ---------------- */

export async function fetchInviteCode() {
  const body = await httpGet('/api/invite/code')
  return body?.data?.code ?? ''
}

export async function fetchInviteStats() {
  const body = await httpGet('/api/invite/stats')
  return body?.data ?? {}
}

export async function fetchInviteReward() {
  const body = await httpGet('/api/config/get?key=system')
  return body?.data?.invite_power ?? null
}

export function buildInviteURL(code) {
  return `${window.location.origin}/invite?aff=${encodeURIComponent(code)}`
}

/* ---------------- 积分 ---------------- */

export async function fetchPowerStats() {
  const body = await httpGet('/api/powerLog/stats')
  return body?.data ?? {}
}

export async function fetchPowerLogs(params) {
  const body = await httpPost('/api/powerLog/list', params)
  return body?.data ?? { items: [], total: 0 }
}

/* ---------------- 登录跳转 ---------------- */

export function loginURL(location = window.location) {
  const redirect = location.pathname + location.search + location.hash
  const params = new URLSearchParams({ redirect })
  const invitation = new URLSearchParams(location.search).get('aff')
  if (invitation) params.set('invite_code', invitation)
  return (invitation ? '/register?' : '/login?') + params
}

export function openLoginPage() {
  window.location.assign(loginURL())
}

export const docsURL = import.meta.env.VITE_DOCS_URL?.trim()

/* ---------------- 格式化 ---------------- */

export function formatTime(value, language) {
  if (value === undefined || value === null || value === '') return '—'
  const date = typeof value === 'number' ? new Date(value * 1000) : new Date(value)
  if (Number.isNaN(date.getTime())) return '—'
  return date.toLocaleString(language)
}

export function formatNumber(value, language) {
  return (value ?? 0).toLocaleString(language)
}

/* ---------------- 下载链接 ----------------
   后端 /api/update/latest 返回安装包的同源相对路径（如 /update/ThinkBuddySetup-dt1.6.14.exe），
   mac 包按清单 mac.<arch>.url 输出（缺分支/缺包不输出该 key，即“暂未提供”）。
   接口不可用时退回平台下载接口（Windows），mac 不做无依据的推导。 */

export const FALLBACK_LINKS = {
  version: '',
  windows: '/api/update/download',
  macArm64: null,
  macX64: null,
  unreachable: true,
}
