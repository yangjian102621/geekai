<template>
  <!-- 邀请页使用精简页脚 -->
  <footer v-if="isInvite" class="footer-shell text-on-dark-soft mt-auto">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6 text-center">
      <RouterLink to="/" class="text-xs text-on-dark-soft hover:text-on-dark transition">← 返回首页</RouterLink>
      <p class="text-xs text-on-dark-soft/70 mt-3">© 2026 ThinkBuddy. Powered by ThinkBuddy (hanyoai).</p>
    </div>
  </footer>

  <footer v-else class="footer-shell">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-14">
      <div class="footer-grid">
        <div>
          <div class="flex items-center gap-2 mb-3">
            <img src="/assets/logo.png" alt="ThinkBuddy" class="footer-logo h-7 w-7 rounded-lg" />
            <span class="text-lg font-semibold text-on-dark">ThinkBuddy</span>
          </div>
          <p class="text-sm text-on-dark-soft max-w-xs">{{ t('footer.tagline') }}</p>
        </div>
        <div>
          <h4 class="text-xs font-semibold text-on-dark uppercase tracking-wider mb-3">{{ t('footer.col_product') }}</h4>
          <ul class="space-y-2 text-sm">
            <li><a :href="downloadUrl ?? '/#download'" class="text-on-dark-soft hover:text-on-dark transition">{{ t('footer.download') }}</a></li>
            <li><a href="#versions" class="text-on-dark-soft hover:text-on-dark transition">{{ t('footer.versions') }}</a></li>
          </ul>
        </div>
        <div>
          <h4 class="text-xs font-semibold text-on-dark uppercase tracking-wider mb-3">{{ t('footer.col_resources') }}</h4>
          <ul class="space-y-2 text-sm">
            <li v-if="docsURL">
              <a :href="docsURL" target="_blank" rel="noopener noreferrer" class="text-on-dark-soft hover:text-on-dark transition">{{ t('footer.docs') }}</a>
            </li>
            <li><a href="#" class="text-on-dark-soft hover:text-on-dark transition">{{ t('footer.faq') }}</a></li>
          </ul>
        </div>
        <div>
          <h4 class="text-xs font-semibold text-on-dark uppercase tracking-wider mb-3">{{ t('footer.col_legal') }}</h4>
          <ul class="space-y-2 text-sm">
            <li><a href="#" class="text-on-dark-soft hover:text-on-dark transition">{{ t('footer.privacy') }}</a></li>
            <li><a href="#" class="text-on-dark-soft hover:text-on-dark transition">{{ t('footer.terms') }}</a></li>
          </ul>
        </div>
      </div>
      <div class="border-t border-white/10 mt-12 pt-6 text-center">
        <p class="text-xs text-on-dark-soft/70">{{ t('footer.copyright') }}</p>
      </div>
    </div>
  </footer>
</template>

<script setup>
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute } from 'vue-router'
import { downloadLinks, recommendedHref } from '@/composables/useDownloadLinks'
import { docsURL } from '@/services/site'

const { t } = useI18n()
const route = useRoute()
const isInvite = computed(() => route.path.startsWith('/invite'))
const downloadUrl = computed(() => recommendedHref(downloadLinks.value))
</script>
