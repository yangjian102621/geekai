<template>
  <!-- 个人空间：积分明细和账号信息。管理功能由独立管理后台（/admin）提供，
      这里只承载官网登录用户自己的数据。 -->
  <section class="min-h-[calc(100dvh-4rem)] page-wash-soft pt-8 pb-16 sm:pt-10">
    <div class="mx-auto w-full max-w-5xl px-4 sm:px-6 lg:px-8">
      <header class="mb-8">
        <h1 class="font-serif text-3xl font-normal tracking-[-0.03em] text-ink">{{ current?.label }}</h1>
        <p class="mt-2 text-muted">{{ description }}</p>
      </header>

      <div v-if="isLoading" class="flex flex-col items-center justify-center gap-3 py-24 text-muted" role="status">
        <span class="h-8 w-8 animate-spin rounded-full border-2 border-hairline border-t-brand-500" aria-hidden="true" />
        <span>{{ t('space.loading') }}</span>
      </div>

      <div v-else-if="!user" class="rounded-xl border border-hairline bg-canvas px-6 py-16 text-center">
        <UserRound class="mx-auto text-muted-soft opacity-40" :size="40" aria-hidden="true" />
        <h2 class="mt-4 text-lg font-semibold text-ink">{{ t('space.login_required') }}</h2>
        <p class="mt-2 text-sm text-muted">{{ t('space.login_hint') }}</p>
        <button type="button" class="btn-primary mt-6 rounded-lg px-6 py-2.5 text-sm font-medium" @click="openLoginPage">
          {{ t('nav.login') }}
        </button>
        <p v-if="isError" role="alert" class="mt-4 text-sm text-rose-600">{{ t('space.load_error') }}</p>
      </div>

      <template v-else>
        <nav :aria-label="t('space.navigate')" class="mb-6 flex gap-1 overflow-x-auto rounded-xl border border-hairline bg-canvas p-1">
          <RouterLink
            v-for="tab in tabs"
            :key="tab.key"
            :to="tab.key"
            :aria-current="tab.key === current?.key ? 'page' : undefined"
            class="inline-flex flex-1 items-center justify-center gap-2 whitespace-nowrap rounded-lg px-4 py-2.5 text-sm font-medium transition-colors"
            :class="tab.key === current?.key ? 'bg-surface-card text-ink' : 'text-muted hover:bg-surface-soft hover:text-ink'"
          >
            <component :is="tab.icon" :size="18" stroke-width="1.75" aria-hidden="true" />
            {{ tab.label }}
          </RouterLink>
        </nav>

        <router-view v-if="user" :user-id="user.id" />
      </template>
    </div>
  </section>
</template>

<script setup>
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute } from 'vue-router'
import { Coins, UserRound } from 'lucide-vue-next'
import { useAuth } from '@/composables/useAuth'
import { openLoginPage } from '@/services/site'

const { t } = useI18n()
const route = useRoute()
const { user, isLoading, isError } = useAuth()

const tabs = computed(() => [
  { key: '/space', label: t('space.points'), icon: Coins },
  { key: '/space/profile', label: t('space.profile'), icon: UserRound },
])

const current = computed(() => {
  const selected = route.path.replace(/\/$/, '')
  return tabs.value.find((tab) => tab.key === selected) ?? tabs.value[0]
})

const description = computed(() =>
  current.value?.key === '/space/profile' ? t('space.profile_hint') : t('space.points_description'),
)
</script>
