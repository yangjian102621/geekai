// * 邀请数据：邀请码 + 统计汇总、系统配置中的邀请奖励。

import { computed, ref, watch } from 'vue'
import { buildInviteURL, fetchInviteCode, fetchInviteReward, fetchInviteStats } from '@/services/site'

// enabled 支持传函数：登录态是异步恢复的，需要在 user 就绪后自动补拉数据。
export function useInviteSummary(enabled) {
  const summary = ref(null)
  const isLoading = ref(false)
  const isError = ref(false)

  async function load() {
    isLoading.value = true
    isError.value = false
    try {
      const code = await fetchInviteCode()
      const stats = await fetchInviteStats()
      summary.value = {
        invite_code: code,
        invite_url: buildInviteURL(code),
        invited_count: stats.invite_count ?? 0,
        total_invite_points: stats.reward_total ?? 0,
        today_invites: stats.today_invite ?? 0,
      }
    } catch (e) {
      isError.value = true
    } finally {
      isLoading.value = false
    }
  }

  const active = computed(() => (typeof enabled === 'function' ? enabled() : enabled))
  watch(
    active,
    (value) => {
      if (value) load()
    },
    { immediate: true },
  )
  return { summary, isLoading, isError, refetch: load }
}

export function useInviteReward() {
  const reward = ref(null)
  fetchInviteReward()
    .then((v) => (reward.value = v))
    .catch(() => {})
  return { reward }
}
