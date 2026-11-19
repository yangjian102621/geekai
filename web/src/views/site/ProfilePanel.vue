<template>
  <div class="space-y-6">
    <div class="overflow-hidden rounded-xl border border-hairline bg-canvas">
      <div class="flex items-center gap-4 border-b border-hairline-soft px-6 py-5">
        <img
          :src="avatarSrc"
          :alt="user.username"
          class="h-14 w-14 rounded-xl object-cover"
          @error="onAvatarError"
        />
        <h2 class="min-w-0 break-all text-lg font-semibold text-ink">{{ user.username }}</h2>
      </div>
      <dl class="divide-y divide-hairline-soft px-6">
        <div v-for="row in rows" :key="row.key" class="grid grid-cols-[120px_minmax(0,1fr)] items-center gap-4 py-4 text-sm">
          <dt class="flex items-center gap-2 text-muted">
            <component :is="row.icon" :size="16" />
            {{ t(`space.${row.key}`) }}
          </dt>
          <dd class="min-w-0 break-all text-ink">{{ row.value || '—' }}</dd>
        </div>
      </dl>
    </div>
    <ModelTokenPanel :user-id="user.id" />
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { Coins, Phone, UserRound } from 'lucide-vue-next'
import ModelTokenPanel from './ModelTokenPanel.vue'
import { formatNumber } from '@/services/site'

const props = defineProps({
  user: { type: Object, required: true },
})

const { t, locale } = useI18n()

const FALLBACK_AVATAR = '/assets/logo.png'

const avatarSrc = computed(() =>
  props.user.avatar && props.user.avatar !== '/images/avatar/user.png' ? props.user.avatar : FALLBACK_AVATAR,
)

const onAvatarError = (event) => {
  if (event.currentTarget.getAttribute('src') !== FALLBACK_AVATAR) {
    event.currentTarget.src = FALLBACK_AVATAR
  }
}

const rows = computed(() => [
  { key: 'username', icon: UserRound, value: props.user.username },
  { key: 'mobile', icon: Phone, value: props.user.mobile },
  { key: 'balance', icon: Coins, value: formatNumber(props.user.power, locale.value) },
])
</script>
