// * 门户站点登录态：模块级共享状态，/user/session 只拉一次（60s 内视为新鲜）。

import { computed, ref } from 'vue'
import { getUserToken, removeUserToken } from '@/store/session'
import { fetchSession } from '@/services/site'

const user = ref(null)
const isLoading = ref(false)
const isError = ref(false)
let loadedAt = 0
const STALE_MS = 60_000

export async function loadAuthUser(force = false) {
  if (!force && loadedAt > 0 && Date.now() - loadedAt < STALE_MS) return user.value
  isLoading.value = true
  isError.value = false
  try {
    if (!getUserToken()) {
      // 未登录：不发起会话请求（与原 React 版 useAuth 行为一致）
      user.value = null
    } else {
      user.value = await fetchSession()
      loadedAt = Date.now()
    }
  } catch (e) {
    // 401 在 fetchSession 内部已视为未登录返回 null；其余错误对外暴露。
    isError.value = true
    user.value = null
  } finally {
    isLoading.value = false
  }
  return user.value
}

export function clearAuthUser() {
  user.value = null
  loadedAt = 0
}

export function useAuth() {
  loadAuthUser()
  return {
    user,
    isLoading,
    isError,
  }
}

export const isLoggedIn = computed(() => !!user.value)
