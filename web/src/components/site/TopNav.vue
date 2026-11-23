<template>
  <header
    class="topnav fixed top-0 inset-x-0 z-50"
    :class="scrolled ? 'topnav-scrolled' : 'bg-canvas/80'"
  >
    <div class="max-w-shell mx-auto px-4 sm:px-6 lg:px-8">
      <div class="flex items-center justify-between h-16 relative">
        <RouterLink to="/" class="flex items-center gap-2 z-10">
          <img src="/assets/logo.png" alt="ThinkBuddy" class="h-8 w-8" />
          <span class="text-lg font-semibold tracking-tight text-ink">ThinkBuddy</span>
        </RouterLink>

        <nav class="hidden md:flex items-center gap-8 absolute left-1/2 -translate-x-1/2">
          <RouterLink to="/" :class="navClass(isHomeActive)">{{ t('nav.home') }}</RouterLink>
          <!-- 定价暂时隐藏 -->
          <RouterLink v-if="false" to="/pricing" :class="navClass(route.path.startsWith('/pricing'))">
            {{ t('nav.pricing') }}
          </RouterLink>
          <RouterLink to="/invite" :class="navClass(route.path.startsWith('/invite'))">{{ t('nav.invite') }}</RouterLink>
          <RouterLink v-if="user" to="/space" :class="navClass(route.path.startsWith('/space'))">{{ t('nav.space') }}</RouterLink>
          <a v-if="docsURL" :href="docsURL" target="_blank" rel="noopener noreferrer" class="nav-link text-sm font-medium transition-colors text-body-strong hover:text-ink">
            {{ t('nav.docs') }}
          </a>
        </nav>

        <div class="hidden md:flex items-center gap-3 z-10">
          <!-- 语言切换 -->
          <button
            class="text-xs font-semibold text-muted hover:text-brand-600 px-2 py-1 rounded-md transition"
            aria-label="toggle language"
            @click="toggleLanguage()"
          >
            {{ langLabel }}
          </button>

          <!-- 登录态 -->
          <div v-if="user" class="flex items-center gap-1">
            <span class="auth-user-name" :title="user.username">{{ user.username }}</span>
            <button class="auth-logout-btn" :title="t('auth.logout')" @click="handleLogout">
              <LogOut :size="13" />
            </button>
          </div>
          <button v-else class="auth-login-btn" @click="openLoginPage">{{ t('nav.login') }}</button>

          <DownloadMenu anchor="bottom">
            <a :href="downloadUrl ?? '/#download'" class="btn-primary inline-flex items-center gap-1.5 px-4 py-2 rounded-lg text-sm font-semibold">
              {{ t('nav.download') }}
              <ChevronDown :size="14" />
            </a>
          </DownloadMenu>
        </div>

        <button class="md:hidden p-2 -mr-2" :aria-expanded="open" aria-label="menu" @click="open = !open">
          <X v-if="open" :size="24" />
          <Menu v-else :size="24" />
        </button>
      </div>

      <div v-if="open" class="md:hidden border-t border-hairline py-4">
        <nav class="flex flex-col gap-1">
          <RouterLink to="/" class="block py-2 px-3 text-sm font-medium text-body-strong rounded hover:bg-surface-soft" @click="open = false">
            {{ t('nav.home') }}
          </RouterLink>
          <RouterLink to="/invite" class="block py-2 px-3 text-sm font-medium text-body-strong rounded hover:bg-surface-soft" @click="open = false">
            {{ t('nav.invite') }}
          </RouterLink>
          <a
            v-if="docsURL"
            :href="docsURL"
            target="_blank"
            rel="noopener noreferrer"
            class="block py-2 px-3 text-sm font-medium text-body-strong rounded hover:bg-surface-soft"
            @click="open = false"
          >
            {{ t('nav.docs') }}
          </a>
          <template v-if="user">
            <RouterLink to="/space" class="block py-2 px-3 text-sm font-medium text-body-strong rounded hover:bg-surface-soft" @click="open = false">
              {{ t('nav.space') }}
            </RouterLink>
            <button class="text-left py-2 px-3 text-sm font-medium text-body-strong rounded hover:bg-surface-soft" @click="handleLogout">
              {{ t('auth.logout') }} ({{ user.username }})
            </button>
          </template>
          <button
            v-else
            class="text-left py-2 px-3 text-sm font-medium text-body-strong rounded hover:bg-surface-soft"
            @click="open = false; openLoginPage()"
          >
            {{ t('nav.login') }}
          </button>
          <template v-if="options.length > 1">
            <template v-for="option in options" :key="option.key">
              <a
                v-if="option.href"
                :href="option.href"
                class="flex items-center gap-3 py-2 px-3 text-sm font-medium text-body-strong rounded hover:bg-surface-soft"
                @click="open = false"
              >
                <PlatformIcon :platform="option.key" :size="16" />
                <span class="flex-1">{{ t(option.labelKey) }}</span>
                <span v-if="option.badgeKey" class="text-[10px] font-semibold text-brand-600">{{ t(option.badgeKey) }}</span>
              </a>
              <div
                v-else
                class="flex items-center gap-3 py-2 px-3 text-sm font-medium text-muted-soft rounded cursor-not-allowed"
                :title="t('download.unavailable')"
              >
                <PlatformIcon :platform="option.key" :size="16" />
                <span class="flex-1">{{ t(option.labelKey) }}</span>
                <span class="text-[10px] font-semibold text-muted-soft">{{ t('download.unavailable') }}</span>
              </div>
            </template>
          </template>
          <a v-else :href="downloadUrl ?? '/#download'" class="btn-primary block text-center px-4 py-2 rounded-lg text-sm font-semibold mt-2">
            {{ t('nav.download') }}
          </a>
        </nav>
      </div>
    </div>
  </header>
</template>

<script setup>
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute } from 'vue-router'
import { ChevronDown, LogOut, Menu, X } from 'lucide-vue-next'
import DownloadMenu from './DownloadMenu.vue'
import PlatformIcon from './PlatformIcon.vue'
import { useAuth } from '@/composables/useAuth'
import { downloadLinks, downloadOptions, recommendedHref } from '@/composables/useDownloadLinks'
import { docsURL, logout, openLoginPage } from '@/services/site'
import { removeUserToken } from '@/store/session'
import { toggleLanguage } from '@/i18n'

const { t, locale } = useI18n()
const route = useRoute()
const open = ref(false)
const scrolled = ref(false)
const { user } = useAuth()

const options = computed(() => downloadOptions(downloadLinks.value))
const downloadUrl = computed(() => recommendedHref(downloadLinks.value))
const isHomeActive = computed(() => route.path === '/')

const langLabel = computed(() => (locale.value.startsWith('zh') ? 'EN' : '中'))

const navClass = (active) =>
  `nav-link text-sm font-medium transition-colors ${active ? 'text-brand-600' : 'text-body-strong hover:text-ink'}`

const handleLogout = async () => {
  await logout()
  removeUserToken()
  window.location.reload()
}

const onScroll = () => {
  scrolled.value = window.scrollY > 8
}
onMounted(() => {
  window.addEventListener('scroll', onScroll, { passive: true })
  onScroll()
})
onBeforeUnmount(() => {
  window.removeEventListener('scroll', onScroll)
})
</script>
