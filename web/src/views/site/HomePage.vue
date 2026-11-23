<template>
  <div v-reveal class="page-wash-all">
    <!-- Hero -->
    <section class="hero-bg relative overflow-hidden pt-28 pb-24 lg:pt-36 lg:pb-32">
      <div class="hero-aurora" aria-hidden="true" />
      <div class="hero-grid-overlay absolute inset-0 pointer-events-none" />
      <div class="relative max-w-shell mx-auto px-4 sm:px-6 lg:px-8 text-center">
        <div
          class="inline-flex items-center gap-2 px-3 py-1.5 rounded-full bg-canvas border border-hairline text-body-strong text-xs font-medium tracking-wide mb-7"
          data-reveal
        >
          <span class="w-1.5 h-1.5 rounded-full bg-coral" />
          <span>{{ t('hero.badge') }}</span>
        </div>
        <h1
          class="font-serif text-[clamp(2rem,7.5vw,2.6rem)] sm:text-5xl lg:text-[3.75rem] font-normal tracking-[-0.03em] text-ink leading-[1.06] mb-7"
          data-reveal
        >
          <span class="block text-ink">{{ t('hero.title_lead') }}</span>
          <!-- 强调行：墨色文字 + 极淡赤陶底色，不用渐变字 -->
          <span class="hero-highlight inline-block">{{ t('hero.title_highlight') }}</span>
        </h1>
        <p class="max-w-2xl mx-auto text-base sm:text-lg text-body leading-relaxed mb-10" data-reveal>
          {{ t('hero.subtitle') }}
        </p>
        <!-- relative z-30：hero 每行都带 data-reveal（transform 会生成层叠上下文），
            否则下方的支持文案行会盖住下载下拉面板。 -->
        <div class="relative z-30 flex flex-col sm:flex-row items-center justify-center gap-3 mb-8" data-reveal>
          <DownloadMenu anchor="right" class-name="w-full sm:w-auto">
            <a :href="downloadUrl ?? '#download'" class="btn-primary inline-flex items-center gap-2 px-6 py-3 rounded-lg text-sm font-medium w-full sm:w-auto justify-center">
              <svg class="download-icon" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4" />
                <polyline points="7 10 12 15 17 10" />
                <line x1="12" y1="15" x2="12" y2="3" />
              </svg>
              <span>{{ t('hero.cta_download') }}</span>
              <ChevronDown :size="14" />
            </a>
          </DownloadMenu>
        </div>
        <p class="text-xs text-muted-soft mt-6" data-reveal>{{ t('hero.platforms') }}</p>
      </div>
    </section>

    <!-- 产品演示 -->
    <section class="py-20 lg:py-24 bg-canvas">
      <div class="max-w-shell mx-auto px-4 sm:px-6 lg:px-8">
        <div class="text-center mb-14" data-reveal>
          <h2 class="font-serif text-3xl sm:text-[2.6rem] font-normal tracking-[-0.03em] text-ink mb-3">{{ t('demo.title') }}</h2>
          <p class="max-w-2xl mx-auto text-body">{{ t('demo.subtitle') }}</p>
        </div>
        <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
          <div
            v-for="card in demoCards"
            :key="card.title"
            class="bg-canvas border border-hairline rounded-xl overflow-hidden"
            data-reveal
          >
            <div class="demo-card-visual h-44 relative" @click="lightbox = card.img">
              <img :src="card.img" :alt="card.title" class="demo-card-img" />
            </div>
            <div class="p-6">
              <h3 class="text-lg font-semibold text-ink mb-2">{{ card.title }}</h3>
              <p class="text-sm text-body leading-relaxed">{{ card.desc }}</p>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- 核心能力矩阵 -->
    <section id="features" class="section-soft py-20 lg:py-24">
      <div class="max-w-shell mx-auto px-4 sm:px-6 lg:px-8">
        <div class="text-center mb-14" data-reveal>
          <h2 class="font-serif text-3xl sm:text-[2.6rem] font-normal tracking-[-0.03em] text-ink mb-3">{{ t('features.title') }}</h2>
          <p class="text-body">{{ t('features.subtitle') }}</p>
          <div class="stats-bar" data-reveal>
            <div class="stat-item">
              <div class="stat-number">6<span class="stat-plus">+</span></div>
              <div class="stat-label">{{ t('stats.modules') }}</div>
            </div>
            <div class="stat-divider" />
            <div class="stat-item">
              <div class="stat-number">101<span class="stat-plus">+</span></div>
              <div class="stat-label">{{ t('stats.apps') }}</div>
            </div>
            <div class="stat-divider" />
            <div class="stat-item">
              <div class="stat-number">30<span class="stat-plus">+</span></div>
              <div class="stat-label">{{ t('stats.models') }}</div>
            </div>
            <div class="stat-divider" />
            <div class="stat-item">
              <div class="stat-number">3</div>
              <div class="stat-label">{{ t('stats.platforms') }}</div>
            </div>
          </div>
        </div>
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
          <div
            v-for="(f, i) in features"
            :key="f.titleKey"
            class="feature-card rounded-xl p-7"
            :style="{ '--tone': tones[i % tones.length].tone, '--tone-soft': tones[i % tones.length].soft }"
            data-reveal
          >
            <div class="feature-icon-wrap w-11 h-11 rounded-lg flex items-center justify-center mb-5" v-html="f.icon" />
            <h3 class="text-lg font-semibold text-ink mb-2">{{ t(f.titleKey) }}</h3>
            <p class="text-sm text-body leading-relaxed">{{ t(f.descKey) }}</p>
          </div>
        </div>
      </div>
    </section>

    <!-- 大字错位排版区 -->
    <section class="bigtext-section py-20 lg:py-24 bg-canvas overflow-hidden">
      <div class="max-w-shell mx-auto px-4 sm:px-6 lg:px-8">
        <h2 class="bigtext-title" data-reveal>{{ t('bigtext.title') }}</h2>
        <p class="bigtext-subtitle" data-reveal>{{ t('bigtext.subtitle') }}</p>
        <div class="bigtext-words" data-reveal>
          <span v-for="(w, i) in bigWords" :key="w" class="big-word" :style="{ '--d': `${i * 120}ms` }">{{ w }}</span>
        </div>
      </div>
    </section>

    <!-- 版本时间线 -->
    <section id="versions" class="py-12 lg:py-16 section-soft">
      <div class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="text-center mb-8" data-reveal>
          <h2 class="font-serif text-2xl sm:text-[2rem] font-normal tracking-[-0.025em] text-ink mb-2">{{ t('versions.title') }}</h2>
          <p class="text-sm text-body">{{ t('versions.subtitle') }}</p>
        </div>
        <div class="relative pl-7" data-reveal>
          <div class="timeline-line absolute left-2 top-2 bottom-2 w-0.5" />
          <div v-for="item in versions" :key="item.v" class="relative mb-5">
            <span class="timeline-dot absolute -left-7 top-1 w-3 h-3 rounded-full" />
            <div class="flex flex-wrap items-baseline gap-2 mb-0.5">
              <span class="text-base font-semibold text-ink">{{ item.v }}</span>
              <span class="text-xs text-muted">{{ item.d }}</span>
            </div>
            <p class="text-xs text-body">{{ item.h }}</p>
          </div>
        </div>
      </div>
    </section>

    <!-- 集成生态 -->
    <section id="ecosystem" class="bg-canvas py-20 lg:py-24">
      <div class="max-w-shell mx-auto px-4 sm:px-6 lg:px-8">
        <div class="text-center mb-14" data-reveal>
          <h2 class="font-serif text-3xl sm:text-[2.6rem] font-normal tracking-[-0.03em] text-ink mb-3">{{ t('ecosystem.title') }}</h2>
          <p class="max-w-2xl mx-auto text-body">{{ t('ecosystem.subtitle') }}</p>
        </div>
        <div class="space-y-10">
          <div v-for="g in ecosystemGroups" :key="g.cat" data-reveal>
            <h3 class="text-sm font-semibold text-muted uppercase tracking-wider mb-3">{{ g.cat }}</h3>
            <div class="flex flex-wrap gap-2">
              <span
                v-for="it in g.items"
                :key="it"
                class="ecosystem-chip bg-canvas border border-hairline rounded-full px-4 py-2 text-sm font-medium text-body-strong"
              >
                {{ it }}
              </span>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- 下载区 -->
    <section id="download" class="pb-12 lg:pb-16 section-soft">
      <div class="max-w-2xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="text-center mb-6" data-reveal>
          <h2 class="font-serif text-2xl sm:text-[2rem] font-normal tracking-[-0.025em] text-ink mb-2">{{ t('download.title') }}</h2>
          <p v-if="links.version" class="text-xs text-muted">v{{ links.version }}</p>
        </div>
        <div v-if="links.unreachable" class="flex items-center justify-center gap-2 mb-3 text-xs text-muted" data-reveal>
          <span>⚠ {{ t('download.unreachable') }}</span>
          <button type="button" class="text-brand-600 font-semibold hover:underline" @click="refreshDownloadLinks">
            {{ t('download.retry') }}
          </button>
        </div>
        <p v-if="!links.unreachable && currentOs() === 'mac' && currentMacArch() === 'unknown'" class="text-xs text-muted text-center mb-3" data-reveal>
          {{ t('download.mac_unknown_hint') }}
        </p>
        <div class="bg-canvas border border-hairline rounded-xl overflow-hidden divide-y divide-hairline-soft" data-reveal>
          <template v-for="option in options" :key="option.key">
            <a v-if="option.available" :href="option.href ?? undefined" class="flex items-center gap-4 px-5 py-4 hover:bg-surface-soft transition">
              <span class="text-body"><PlatformIcon :platform="option.key" :size="22" /></span>
              <span class="flex-1 text-left">
                <span class="block text-sm font-semibold text-ink">{{ t(option.labelKey) }}</span>
                <span class="block text-xs text-muted">{{ t(option.noteKey) }}</span>
              </span>
              <span
                v-if="option.badgeKey"
                class="rounded-full border px-2 py-0.5 text-[10px] font-semibold"
                :class="option.badgeKey === 'download.current' ? 'border-hairline text-muted' : 'border-brand-500/30 bg-brand-50 text-brand-600'"
              >
                {{ t(option.badgeKey) }}
              </span>
              <span class="text-muted-soft">
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4" /><polyline points="7 10 12 15 17 10" /><line x1="12" y1="15" x2="12" y2="3" /></svg>
              </span>
            </a>
            <div v-else class="flex items-center gap-4 px-5 py-4 opacity-55 cursor-not-allowed" :title="t('download.unavailable')">
              <span class="text-muted"><PlatformIcon :platform="option.key" :size="22" /></span>
              <span class="flex-1 text-left">
                <span class="block text-sm font-semibold text-muted">{{ t(option.labelKey) }}</span>
                <span class="block text-xs text-muted-soft">{{ t(option.noteKey) }}</span>
              </span>
              <span class="rounded-full border border-hairline px-2 py-0.5 text-[10px] font-semibold text-muted-soft">
                {{ t('download.unavailable') }}
              </span>
              <span class="text-muted-soft opacity-40">
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4" /><polyline points="7 10 12 15 17 10" /><line x1="12" y1="15" x2="12" y2="3" /></svg>
              </span>
            </div>
          </template>
        </div>
        <!-- 定价入口暂时隐藏 -->
        <p class="hidden text-xs text-muted-soft mt-1 text-center">
          {{ t('credits.more_pricing') }}
          <RouterLink to="/pricing" class="text-brand-600 font-medium hover:underline">定价页</RouterLink>
        </p>
      </div>
    </section>

    <!-- Lightbox -->
    <div v-if="lightbox" class="lightbox-overlay active" @click="lightbox = null">
      <img :src="lightbox" alt="放大预览" @click.stop />
    </div>
  </div>
</template>

<script setup>
import { onBeforeUnmount, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { ChevronDown } from 'lucide-vue-next'
import DownloadMenu from '@/components/site/DownloadMenu.vue'
import PlatformIcon from '@/components/site/PlatformIcon.vue'
import { useReveal } from '@/composables/useReveal'
import {
  currentMacArch,
  currentOs,
  downloadLinks,
  downloadOptions,
  recommendedHref,
  refreshDownloadLinks,
} from '@/composables/useDownloadLinks'
import { computed } from 'vue'

const vReveal = useReveal()
const { t } = useI18n()
const lightbox = ref(null)
const links = downloadLinks
const options = computed(() => downloadOptions(links.value))
const downloadUrl = computed(() => recommendedHref(links.value))

const demoCards = computed(() => [
  { img: '/assets/demo-reading.png', title: t('demo.card1.title'), desc: t('demo.card1.desc') },
  { img: '/assets/demo-knowledge-base.png', title: t('demo.card2.title'), desc: t('demo.card2.desc') },
  { img: '/assets/demo-smart-writing.png', title: t('demo.card3.title'), desc: t('demo.card3.desc') },
])

const bigWords = ['阅读', '笔记', '研究', '写作', '问答', '视频', '课程']

const versions = computed(() => [
  { v: 'v1.7.3', d: '2026.11.20', h: t('versions.v173.highlight') },
  { v: 'v1.7.2', d: '2026.11.13', h: t('versions.v172.highlight') },
  { v: 'v1.7.1', d: '2026.11.6', h: t('versions.v171.highlight') },
  { v: 'v1.7.0', d: '2026.10.30', h: t('versions.v170.highlight') },
  { v: 'v1.6.15', d: '2026.10.23', h: t('versions.v1615.highlight') },
])

const ecosystemGroups = computed(() => [
  { cat: t('ecosystem.cat.im'), items: ['飞书 / Feishu', '微信 / WeChat', 'Telegram'] },
  { cat: t('ecosystem.cat.tools'), items: ['MCP Server', '101+ CLI Apps', '第三方 Capability 插件'] },
  { cat: t('ecosystem.cat.kb'), items: ['LightRAG', 'GraphRAG', 'PageIndex', 'MarginNote 4', 'GitHub 仓库同步'] },
  { cat: t('ecosystem.cat.models'), items: ['OpenAI Codex', 'Gemini', '通义', '豆包', 'Zhipu', '…数十家'] },
])

// 能力卡图标只用三种低饱和冷色轮换：靛蓝 / 青蓝 / 紫灰。
// 全站强调色保持稀缺，这里的色彩是“点缀”而不是主题色。
const tones = [
  { tone: '#3b5bdb', soft: '#e9ebf7' },
  { tone: '#2f7d8c', soft: '#e4edf1' },
  { tone: '#6b5bd6', soft: '#eeebf5' },
]

const features = [
  { titleKey: 'features.reading.title', descKey: 'features.reading.desc', icon: '<svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M2 3h6a4 4 0 0 1 4 4v14a3 3 0 0 0-3-3H2z"/><path d="M22 3h-6a4 4 0 0 0-4 4v14a3 3 0 0 1 3-3h7z"/></svg>' },
  { titleKey: 'features.knowledge.title', descKey: 'features.knowledge.desc', icon: '<svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/></svg>' },
  { titleKey: 'features.learning.title', descKey: 'features.learning.desc', icon: '<svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M4 19.5A2.5 2.5 0 0 1 6.5 17H20"/><path d="M6.5 2H20v20H6.5A2.5 2.5 0 0 1 4 19.5v-15A2.5 2.5 0 0 1 6.5 2z"/></svg>' },
  { titleKey: 'features.video.title', descKey: 'features.video.desc', icon: '<svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polygon points="23 7 16 12 23 17 23 7"/><rect x="1" y="5" width="15" height="14" rx="2" ry="2"/></svg>' },
  { titleKey: 'features.cowriter.title', descKey: 'features.cowriter.desc', icon: '<svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 20h9"/><path d="M16.5 3.5a2.121 2.121 0 0 1 3 3L7 19l-4 1 1-4L16.5 3.5z"/></svg>' },
  { titleKey: 'features.agent.title', descKey: 'features.agent.desc', icon: '<svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="3"/><path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1-2.83 2.83l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-4 0v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1 0-4h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 2.83-2.83l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 4 0v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 0 4h-.09a1.65 1.65 0 0 0-1.51 1z"/></svg>' },
]

const onKey = (e) => {
  if (e.key === 'Escape') lightbox.value = null
}
onMounted(() => window.addEventListener('keydown', onKey))
onBeforeUnmount(() => window.removeEventListener('keydown', onKey))
</script>
