<template>
  <div class="space-y-6">
    <div
      v-if="hasError"
      role="alert"
      class="flex flex-wrap items-center justify-between gap-3 rounded-xl border border-amber-200 bg-amber-50 px-5 py-4 text-sm text-amber-900"
    >
      <span>{{ t('space.load_error') }}</span>
      <button
        type="button"
        class="rounded-lg bg-amber-600 px-4 py-1.5 font-semibold text-white hover:bg-amber-700"
        @click="retry"
      >
        {{ t('space.retry') }}
      </button>
    </div>

    <div class="grid gap-4 sm:grid-cols-2 xl:grid-cols-5">
      <StatCard :icon="Coins" :label="t('space.balance')" :value="formatNumber(stats?.balance, language)" :loading="statsLoading" />
      <StatCard :icon="TrendingDown" :label="t('space.total_consumed')" :value="formatNumber(stats?.total, language)" :loading="statsLoading" />
      <StatCard :icon="TrendingUp" :label="t('space.today_consumed')" :value="formatNumber(stats?.today, language)" :loading="statsLoading" />
      <StatCard :icon="Gift" :label="t('space.invite_rewards')" :value="formatNumber(invite?.reward_total, language)" :loading="inviteLoading" />
      <StatCard :icon="Coins" :label="t('space.quota_granted')" :value="formatNumber(stats?.quota_granted, language)" :loading="statsLoading" />
    </div>

    <div class="rounded-xl border border-hairline bg-canvas">
      <div class="flex flex-wrap items-center justify-between gap-3 border-b border-hairline-soft px-5 py-4">
        <div>
          <h2 class="text-base font-semibold text-ink">{{ t('space.points') }}</h2>
          <p class="mt-1 text-sm text-muted">{{ t('space.points_description') }}</p>
        </div>
        <RouterLink
          to="/invite"
          class="inline-flex items-center gap-2 rounded-lg border border-hairline px-3 py-1.5 text-sm font-medium text-body-strong hover:border-brand-500 hover:text-brand-600"
        >
          <Gift :size="15" />
          {{ t('space.invite') }}
          <span v-if="invite" class="text-xs text-muted-soft">{{ t('space.invite_count', { count: invite.invite_count }) }}</span>
        </RouterLink>
      </div>

      <div class="flex flex-wrap items-end gap-3 border-b border-hairline-soft px-5 py-4">
        <label class="flex flex-col gap-1 text-xs font-medium text-muted">
          <span>{{ t('space.start_date') }}</span>
          <input
            v-model="from"
            type="date"
            :max="to || todayString()"
            class="rounded-lg border border-hairline px-3 py-2 text-sm text-body-strong outline-none focus:border-brand-500"
            @change="page = 1; loadLogs()"
          />
        </label>
        <label class="flex flex-col gap-1 text-xs font-medium text-muted">
          <span>{{ t('space.end_date') }}</span>
          <input
            v-model="to"
            type="date"
            :min="from || undefined"
            :max="todayString()"
            class="rounded-lg border border-hairline px-3 py-2 text-sm text-body-strong outline-none focus:border-brand-500"
            @change="page = 1; loadLogs()"
          />
        </label>
        <button
          type="button"
          class="rounded-lg border border-hairline px-3 py-2 text-sm text-body hover:bg-surface-soft"
          @click="clearFilter"
        >
          {{ t('space.clear') }}
        </button>
        <div class="ml-auto flex items-center gap-2 text-xs text-muted-soft">
          <CalendarDays :size="14" />
          {{ from || to ? t('space.date_range') : t('space.all') }}
        </div>
      </div>

      <div class="overflow-x-auto">
        <table class="min-w-full divide-y divide-hairline-soft text-sm">
          <thead class="bg-surface-soft text-left text-xs font-semibold uppercase tracking-wide text-muted">
            <tr>
              <th class="px-5 py-3">{{ t('space.time') }}</th>
              <th class="px-5 py-3">{{ t('space.type') }}</th>
              <th class="px-5 py-3 text-right">{{ t('space.amount') }}</th>
              <th class="px-5 py-3 text-right">{{ t('space.balance_after') }}</th>
              <th class="px-5 py-3 text-right">{{ t('space.quota_amount') }}</th>
              <th class="px-5 py-3">{{ t('space.quota_status') }}</th>
              <th class="px-5 py-3">{{ t('space.detail') }}</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-hairline-soft">
            <tr v-if="logsLoading">
              <td colspan="7" class="px-5 py-10 text-center text-muted">{{ t('space.loading') }}</td>
            </tr>
            <template v-else-if="logs?.items?.length">
              <tr v-for="log in logs.items" :key="log.id" class="hover:bg-surface-soft">
                <td class="whitespace-nowrap px-5 py-3 text-body">{{ formatTime(log.created_at, language) }}</td>
                <td class="whitespace-nowrap px-5 py-3">
                  <span class="rounded-full bg-surface-card px-2.5 py-1 text-xs font-medium text-body-strong">{{ typeLabel(log) }}</span>
                </td>
                <td
                  class="whitespace-nowrap px-5 py-3 text-right font-semibold"
                  :class="log.mark === 1 ? 'text-emerald-600' : 'text-rose-600'"
                >
                  {{ log.mark === 1 ? '+' : '-' }}{{ formatNumber(Math.abs(log.amount), language) }}
                </td>
                <td class="whitespace-nowrap px-5 py-3 text-right text-body">{{ formatNumber(log.balance, language) }}</td>
                <td class="whitespace-nowrap px-5 py-3 text-right text-body">
                  {{ log.quota_amount ? `+${formatNumber(log.quota_amount, language)}` : '—' }}
                </td>
                <td class="whitespace-nowrap px-5 py-3 text-body" :title="log.quota_error ? t('space.quota_error_hint') : undefined">
                  <span
                    v-if="log.quota_status"
                    class="rounded-full px-2.5 py-1 text-xs font-medium"
                    :class="quotaPillClass(log.quota_status)"
                  >{{ t(`space.quota_states.${log.quota_status}`, t('space.quota_states.unknown')) }}</span>
                  <template v-else>—</template>
                </td>
                <td class="min-w-[180px] px-5 py-3 text-muted">{{ log.remark || log.model || '—' }}</td>
              </tr>
            </template>
            <tr v-else>
              <td colspan="7" class="px-5 py-10 text-center text-muted">{{ t('space.no_records') }}</td>
            </tr>
          </tbody>
        </table>
      </div>

      <div class="flex flex-wrap items-center justify-between gap-3 border-t border-hairline-soft px-5 py-4 text-sm text-muted">
        <span>{{ t('space.page', { page, total: totalPages }) }}</span>
        <div class="flex items-center gap-2">
          <button
            type="button"
            :disabled="page <= 1"
            class="rounded-lg border border-hairline px-3 py-1.5 disabled:cursor-not-allowed disabled:opacity-40"
            @click="page = Math.max(1, page - 1); loadLogs()"
          >
            {{ t('space.previous') }}
          </button>
          <button
            type="button"
            :disabled="page >= totalPages"
            class="rounded-lg border border-hairline px-3 py-1.5 disabled:cursor-not-allowed disabled:opacity-40"
            @click="page = Math.min(totalPages, page + 1); loadLogs()"
          >
            {{ t('space.next') }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { CalendarDays, Coins, Gift, TrendingDown, TrendingUp } from 'lucide-vue-next'
import StatCard from './StatCard.vue'
import { fetchInviteStats, fetchPowerLogs, fetchPowerStats, formatNumber, formatTime } from '@/services/site'

const { t, locale } = useI18n()
const language = computed(() => locale.value)

defineProps({
  userId: { type: Number, required: true },
})

const page = ref(1)
const pageSize = 10
const from = ref('')
const to = ref('')

const stats = ref(null)
const statsLoading = ref(false)
const statsError = ref(false)
const invite = ref(null)
const inviteLoading = ref(false)
const logs = ref(null)
const logsLoading = ref(false)
const logsError = ref(false)
let statsTimer
let logsTimer
let logsRequest = 0

const total = computed(() => logs.value?.total ?? 0)
const totalPages = computed(() => Math.max(1, Math.ceil(total.value / pageSize)))
const hasError = computed(() => statsError.value || logsError.value)

function todayString(offsetDays = 0) {
  const date = new Date()
  date.setDate(date.getDate() + offsetDays)
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  return `${date.getFullYear()}-${month}-${day}`
}

function typeLabel(log) {
  const key = `space.types.${log.type}`
  return t(key) !== key ? t(key) : log.type_str || t('space.other')
}

function quotaPillClass(status) {
  if (status === 'succeeded') return 'bg-emerald-50 text-emerald-700'
  if (status === 'failed' || status === 'unknown') return 'bg-amber-50 text-amber-800'
  return 'bg-surface-card text-body-strong'
}

async function loadStats(background = false) {
  if (statsLoading.value) return
  statsLoading.value = !background
  statsError.value = false
  try {
    stats.value = await fetchPowerStats()
  } catch (e) {
    statsError.value = true
  } finally {
    statsLoading.value = false
  }
}

async function loadInvite() {
  inviteLoading.value = true
  try {
    invite.value = await fetchInviteStats()
  } catch (e) {
    // 邀请统计失败不阻塞积分面板
  } finally {
    inviteLoading.value = false
  }
}

async function loadLogs(background = false) {
  clearTimeout(logsTimer)
  const request = ++logsRequest
  logsLoading.value = !background
  logsError.value = false
  try {
    const result = await fetchPowerLogs({
      page: page.value,
      page_size: pageSize,
      date: from.value && to.value ? [from.value, to.value] : undefined,
    })
    if (request !== logsRequest) return
    logs.value = result
    if (result.items?.some(({ quota_status }) => quota_status === 'pending' || quota_status === 'processing')) {
      logsTimer = setTimeout(() => loadLogs(true), 5000)
    }
  } catch (e) {
    if (request === logsRequest) logsError.value = true
  } finally {
    if (request === logsRequest) logsLoading.value = false
  }
}

function clearFilter() {
  from.value = ''
  to.value = ''
  page.value = 1
  loadLogs()
}

function retry() {
  loadStats()
  loadLogs()
}

onMounted(() => {
  loadStats()
  loadInvite()
  loadLogs()
  statsTimer = setInterval(() => loadStats(true), 15000)
})

onBeforeUnmount(() => {
  ++logsRequest
  clearInterval(statsTimer)
  clearTimeout(logsTimer)
})
</script>
