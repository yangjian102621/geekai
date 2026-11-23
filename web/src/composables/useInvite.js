// * 邀请数据：邀请码 + 统计汇总、系统配置中的邀请奖励。

import { ref, watch } from 'vue'
import { buildInviteURL, fetchInviteCode, fetchInviteReward, fetchInviteStats } from '@/services/site'

export function useInviteSummary(shouldLoad) {
  const summary = ref(null)
  const isLoading = ref(false)
  const isError = ref(false)

  async function load() {
    if (isLoading.value) return
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

  watch(shouldLoad, (enabled) => {
    if (enabled) load()
  }, { immediate: true })

  return { summary, isLoading, isError, refetch: load }
}

export function useInviteReward() {
  const reward = ref(null)
  fetchInviteReward()
    .then((v) => (reward.value = v))
    .catch(() => {})
  return { reward }
}
