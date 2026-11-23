<template>
  <div v-reveal class="tb-home">
    <!-- Hero -->
    <section class="tb-hero">
      <div class="tb-hero-grid" aria-hidden="true" />
      <div class="tb-hero-noise" aria-hidden="true" />
      <div class="tb-container tb-hero-inner">
        <div class="tb-pill" data-reveal>
          <span class="tb-pill-dot" />
          <span>{{ t('hero.badge') }}</span>
        </div>
        <h1 class="tb-title" data-reveal>
          <span class="tb-title-line">{{ t('hero.title_lead') }}</span>
          <span class="tb-grad-text">{{ t('hero.title_highlight') }}</span>
        </h1>
        <p class="tb-sub" data-reveal>{{ t('hero.subtitle') }}</p>
        <div class="tb-cta-row" data-reveal>
          <DownloadMenu anchor="right" class-name="tb-dm">
            <a :href="downloadUrl ?? '#download'" class="tb-btn tb-btn-primary">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4" />
                <polyline points="7 10 12 15 17 10" />
                <line x1="12" y1="15" x2="12" y2="3" />
              </svg>
              <span>{{ t('hero.cta_download') }}</span>
              <ChevronDown :size="14" />
            </a>
          </DownloadMenu>
        </div>
        <p class="tb-platforms" data-reveal>{{ t('hero.platforms') }}</p>
      </div>
    </section>

    <!-- 产品演示 -->
    <section class="tb-section tb-demo">
      <div class="tb-container">
        <div class="tb-head" data-reveal>
          <h2 class="tb-h2">{{ t('demo.title') }}</h2>
          <p class="tb-sub2">{{ t('demo.subtitle') }}</p>
        </div>
        <div class="tb-demo-grid">
          <article v-for="card in demoCards" :key="card.title" class="tb-card" data-reveal>
            <div class="tb-card-media" @click="lightbox = card.img">
              <img :src="card.img" :alt="card.title" />
            </div>
            <div class="tb-card-body">
              <h3>{{ card.title }}</h3>
              <p>{{ card.desc }}</p>
            </div>
          </article>
        </div>
      </div>
    </section>

    <!-- 核心能力矩阵 -->
    <section id="features" class="tb-section tb-features">
      <div class="tb-container">
        <div class="tb-head" data-reveal>
          <h2 class="tb-h2">{{ t('features.title') }}</h2>
          <p class="tb-sub2">{{ t('features.subtitle') }}</p>
          <div class="tb-stats" data-reveal>
            <div class="tb-stat">
              <b>6<span>+</span></b>
              <i>{{ t('stats.modules') }}</i>
            </div>
            <div class="tb-stat">
              <b>101<span>+</span></b>
              <i>{{ t('stats.apps') }}</i>
            </div>
            <div class="tb-stat">
              <b>30<span>+</span></b>
              <i>{{ t('stats.models') }}</i>
            </div>
            <div class="tb-stat">
              <b>3</b>
              <i>{{ t('stats.platforms') }}</i>
            </div>
          </div>
        </div>
        <div class="tb-feature-grid">
          <article v-for="(f, i) in features" :key="f.titleKey" class="tb-feature" data-reveal>
            <i class="tb-feature-icon" :style="{ background: tones[i % tones.length] }" v-html="f.icon" />
            <h3>{{ t(f.titleKey) }}</h3>
            <p>{{ t(f.descKey) }}</p>
          </article>
        </div>
      </div>
    </section>

    <!-- 大字错位排版区 -->
    <section class="tb-section tb-big">
      <div class="tb-container">
        <h2 class="tb-big-title tb-grad-text" data-reveal>{{ t('bigtext.title') }}</h2>
        <p class="tb-big-sub" data-reveal>{{ t('bigtext.subtitle') }}</p>
        <div class="tb-words" data-reveal>
          <span v-for="w in bigWords" :key="w" class="tb-word">{{ w }}</span>
        </div>
      </div>
    </section>

    <!-- 版本时间线 -->
    <section id="versions" class="tb-section tb-versions">
      <div class="tb-container tb-narrow">
        <div class="tb-head" data-reveal>
          <h2 class="tb-h2">{{ t('versions.title') }}</h2>
          <p class="tb-sub2">{{ t('versions.subtitle') }}</p>
        </div>
        <div class="tb-version-card" data-reveal>
          <div class="tb-timeline">
            <div v-for="item in versions" :key="item.v" class="tb-version-item">
              <span class="tb-dot" />
              <div class="tb-version-head">
                <b>{{ item.v }}</b>
                <span>{{ item.d }}</span>
              </div>
              <p>{{ item.h }}</p>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- 集成生态 -->
    <section id="ecosystem" class="tb-section tb-eco">
      <div class="tb-container">
        <div class="tb-head" data-reveal>
          <h2 class="tb-h2">{{ t('ecosystem.title') }}</h2>
          <p class="tb-sub2">{{ t('ecosystem.subtitle') }}</p>
        </div>
        <div class="tb-eco-groups">
          <div v-for="g in ecosystemGroups" :key="g.cat" class="tb-eco-group" data-reveal>
            <h3 class="tb-group-title">{{ g.cat }}</h3>
            <div class="tb-chips">
              <span v-for="it in g.items" :key="it" class="tb-chip">{{ it }}</span>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- 下载区 -->
    <section id="download" class="tb-section tb-download">
      <div class="tb-container tb-narrow-sm">
        <div class="tb-head" data-reveal>
          <h2 class="tb-h2">{{ t('download.title') }}</h2>
          <p v-if="links.version" class="tb-version">v{{ links.version }}</p>
        </div>
        <div v-if="links.unreachable" class="tb-note" data-reveal>
          <span>⚠ {{ t('download.unreachable') }}</span>
          <button type="button" @click="refreshDownloadLinks">{{ t('download.retry') }}</button>
        </div>
        <p
          v-if="!links.unreachable && currentOs() === 'mac' && currentMacArch() === 'unknown'"
          class="tb-note"
          data-reveal
        >
          {{ t('download.mac_unknown_hint') }}
        </p>
        <div class="tb-download-card" data-reveal>
          <template v-for="option in options" :key="option.key">
            <a v-if="option.available" :href="option.href ?? undefined" class="tb-download-row">
              <span class="tb-dl-icon"><PlatformIcon :platform="option.key" :size="22" /></span>
              <span class="tb-dl-text">
                <span class="tb-dl-title">{{ t(option.labelKey) }}</span>
                <span class="tb-dl-note">{{ t(option.noteKey) }}</span>
              </span>
              <span
                v-if="option.badgeKey"
                class="tb-dl-badge"
                :class="option.badgeKey === 'download.current' ? 'is-current' : 'is-new'"
              >
                {{ t(option.badgeKey) }}
              </span>
              <span class="tb-dl-arrow">
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4" /><polyline points="7 10 12 15 17 10" /><line x1="12" y1="15" x2="12" y2="3" /></svg>
              </span>
            </a>
            <div v-else class="tb-download-row is-disabled" :title="t('download.unavailable')">
              <span class="tb-dl-icon"><PlatformIcon :platform="option.key" :size="22" /></span>
              <span class="tb-dl-text">
                <span class="tb-dl-title">{{ t(option.labelKey) }}</span>
                <span class="tb-dl-note">{{ t(option.noteKey) }}</span>
              </span>
              <span class="tb-dl-badge">{{ t('download.unavailable') }}</span>
              <span class="tb-dl-arrow is-muted">
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4" /><polyline points="7 10 12 15 17 10" /><line x1="12" y1="15" x2="12" y2="3" /></svg>
              </span>
            </div>
          </template>
        </div>
      </div>
    </section>

    <!-- Lightbox -->
    <div v-if="lightbox" class="tb-lightbox" @click="lightbox = null">
      <img :src="lightbox" alt="放大预览" @click.stop />
    </div>
  </div>
</template>

<script setup>
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
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

// 能力卡图标底色：与参考站一致的一组低饱和渐变
const tones = [
  'linear-gradient(135deg, #5457e8, #8b5cf6)',
  'linear-gradient(135deg, #06b6d4, #3b82f6)',
  'linear-gradient(135deg, #8b5cf6, #ec4899)',
  'linear-gradient(135deg, #10b981, #0d9488)',
  'linear-gradient(135deg, #f59e0b, #f43f5e)',
  'linear-gradient(135deg, #22d3ee, #5457e8)',
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

<style scoped>
/* ============================================================
   首页皮肤：沿用原 ThinkBuddy 内容结构，替换为参考站配色与字体
   ============================================================ */
.tb-home {
  --tb-primary: #5457e8;
  --tb-primary-dark: #3f42c8;
  --tb-violet: #8b5cf6;
  --tb-cyan: #22d3ee;
  --tb-text: #14163a;
  --tb-text-2: #565b7e;
  --tb-text-3: #8a8eab;
  --tb-line: rgba(84, 87, 232, 0.14);
  --tb-bg: #f6f7fd;
  --tb-soft: #f2f3fc;
  --tb-grad: linear-gradient(120deg, #5457e8 0%, #8b5cf6 55%, #22d3ee 110%);
  --tb-grad-btn: linear-gradient(135deg, #5457e8 0%, #7c5cfc 100%);
  --tb-font: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'PingFang SC', 'Hiragino Sans GB',
    'Microsoft YaHei', sans-serif;
  --tb-font-num: 'Bahnschrift', 'DIN Alternate', 'Avenir Next Condensed', 'Segoe UI',
    'PingFang SC', 'Microsoft YaHei', sans-serif;

  font-family: var(--tb-font);
  color: var(--tb-text);
  background: var(--tb-bg);
  line-height: 1.65;
}

.tb-home *,
.tb-home *::before,
.tb-home *::after {
  box-sizing: border-box;
}

.tb-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 24px;
}
.tb-narrow {
  max-width: 960px;
}
.tb-narrow-sm {
  max-width: 720px;
}

.tb-section {
  padding: 72px 0;
  position: relative;
}
.tb-head {
  text-align: center;
  margin-bottom: 44px;
}
.tb-h2 {
  font-size: clamp(1.85rem, 3.4vw, 2.6rem);
  font-weight: 800;
  letter-spacing: -0.01em;
  line-height: 1.25;
  color: var(--tb-text);
  margin: 0 0 14px;
}
.tb-sub2 {
  font-size: 16px;
  color: var(--tb-text-2);
  max-width: 640px;
  margin: 0 auto;
}
.tb-grad-text {
  background: var(--tb-grad);
  -webkit-background-clip: text;
  background-clip: text;
  -webkit-text-fill-color: transparent;
  color: transparent;
}

/* ============================================================
   Hero
   ============================================================ */
.tb-hero {
  position: relative;
  overflow: hidden;
  padding: 148px 0 96px;
  background:
    radial-gradient(900px 500px at 15% 0%, rgba(139, 92, 246, 0.12), transparent 60%),
    radial-gradient(800px 480px at 85% 10%, rgba(34, 211, 238, 0.1), transparent 60%),
    linear-gradient(180deg, #f7f8fe 0%, #eef0fc 55%, #f6f7fd 100%);
}
.tb-hero-grid {
  position: absolute;
  inset: 0;
  background-image:
    linear-gradient(rgba(84, 87, 232, 0.055) 1px, transparent 1px),
    linear-gradient(90deg, rgba(84, 87, 232, 0.055) 1px, transparent 1px);
  background-size: 54px 54px;
  -webkit-mask-image: radial-gradient(ellipse 90% 70% at 50% 30%, #000 30%, transparent 75%);
  mask-image: radial-gradient(ellipse 90% 70% at 50% 30%, #000 30%, transparent 75%);
  pointer-events: none;
}
.tb-hero-noise {
  position: absolute;
  inset: 0;
  opacity: 0.5;
  pointer-events: none;
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='120' height='120'%3E%3Cfilter id='n'%3E%3CfeTurbulence type='fractalNoise' baseFrequency='0.9' numOctaves='2'/%3E%3CfeColorMatrix type='saturate' values='0'/%3E%3CfeComponentTransfer%3E%3CfeFuncA type='linear' slope='0.035'/%3E%3C/feComponentTransfer%3E%3C/filter%3E%3Crect width='120' height='120' filter='url(%23n)'/%3E%3C/svg%3E");
}
.tb-hero-inner {
  position: relative;
  text-align: center;
}
.tb-pill {
  display: inline-flex;
  align-items: center;
  gap: 10px;
  padding: 9px 20px;
  border-radius: 999px;
  font-size: 14px;
  color: var(--tb-text);
  background: rgba(255, 255, 255, 0.72);
  border: 1px solid rgba(84, 87, 232, 0.22);
  box-shadow: 0 4px 20px rgba(84, 87, 232, 0.12);
  backdrop-filter: blur(8px);
  -webkit-backdrop-filter: blur(8px);
  margin-bottom: 26px;
}
.tb-pill-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: var(--tb-grad);
  box-shadow: 0 0 10px rgba(84, 87, 232, 0.6);
  flex: none;
}
.tb-title {
  font-size: clamp(2.2rem, 5.6vw, 3.9rem);
  font-weight: 800;
  line-height: 1.16;
  letter-spacing: -0.02em;
  margin: 0 0 24px;
  color: var(--tb-text);
}
.tb-title-line {
  display: block;
}
.tb-sub {
  font-size: clamp(16px, 1.8vw, 19px);
  color: var(--tb-text-2);
  max-width: 680px;
  margin: 0 auto 34px;
}
.tb-cta-row {
  display: flex;
  justify-content: center;
  flex-wrap: wrap;
}
.tb-dm {
  width: 100%;
}
@media (min-width: 640px) {
  .tb-dm {
    width: auto;
  }
}
.tb-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  min-height: 52px;
  padding: 14px 32px;
  border-radius: 14px;
  font-size: 16px;
  font-weight: 600;
  border: 1px solid transparent;
  cursor: pointer;
  text-decoration: none;
  transition:
    transform 0.25s ease,
    box-shadow 0.25s ease,
    background 0.25s ease;
}
.tb-btn-primary {
  color: #fff;
  background: var(--tb-grad-btn);
  box-shadow: 0 8px 24px rgba(84, 87, 232, 0.38);
}
.tb-btn-primary:hover {
  transform: translateY(-2px);
  box-shadow: 0 14px 34px rgba(84, 87, 232, 0.48);
}
.tb-platforms {
  margin: 22px 0 0;
  font-size: 12px;
  color: var(--tb-text-3);
}

/* ============================================================
   产品演示
   ============================================================ */
.tb-demo {
  background: var(--tb-bg);
}
.tb-demo-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 24px;
}
.tb-card {
  background: #fff;
  border: 1px solid var(--tb-line);
  border-radius: 22px;
  overflow: hidden;
  box-shadow: 0 4px 18px rgba(20, 22, 58, 0.05);
  transition:
    transform 0.3s ease,
    box-shadow 0.3s ease,
    border-color 0.3s ease;
}
.tb-card:hover {
  transform: translateY(-6px);
  border-color: rgba(84, 87, 232, 0.4);
  box-shadow: 0 22px 48px rgba(84, 87, 232, 0.16);
}
.tb-card-media {
  position: relative;
  height: 176px;
  overflow: hidden;
  background: #eef0fc;
  border-bottom: 1px solid var(--tb-line);
  cursor: zoom-in;
}
.tb-card-media img {
  position: absolute;
  inset: 0;
  width: 100%;
  height: 100%;
  object-fit: cover;
  object-position: center top;
  transition: transform 600ms cubic-bezier(0.2, 0, 0.2, 1);
}
.tb-card:hover .tb-card-media img {
  transform: scale(1.04);
}
.tb-card-body {
  padding: 24px 26px 28px;
}
.tb-card-body h3 {
  font-size: 19px;
  font-weight: 700;
  margin: 0 0 10px;
  color: var(--tb-text);
}
.tb-card-body p {
  margin: 0;
  font-size: 14px;
  color: var(--tb-text-2);
  line-height: 1.75;
}

/* ============================================================
   核心能力
   ============================================================ */
.tb-features {
  background:
    radial-gradient(700px 400px at 90% 10%, rgba(139, 92, 246, 0.07), transparent 60%),
    var(--tb-soft);
}
.tb-stats {
  display: flex;
  align-items: center;
  justify-content: center;
  flex-wrap: wrap;
  margin: 36px auto 0;
  padding: 24px 28px;
  background: #fff;
  border: 1px solid var(--tb-line);
  border-radius: 18px;
  box-shadow: 0 4px 18px rgba(20, 22, 58, 0.05);
  max-width: 760px;
}
.tb-stat {
  padding: 0 34px;
  text-align: center;
  position: relative;
}
.tb-stat + .tb-stat::before {
  content: '';
  position: absolute;
  left: 0;
  top: 50%;
  transform: translateY(-50%);
  width: 1px;
  height: 36px;
  background: rgba(84, 87, 232, 0.16);
}
.tb-stat b {
  display: block;
  font-family: var(--tb-font-num);
  font-size: 34px;
  font-weight: 700;
  line-height: 1;
  letter-spacing: 0.01em;
  background: var(--tb-grad);
  -webkit-background-clip: text;
  background-clip: text;
  -webkit-text-fill-color: transparent;
}
.tb-stat b span {
  font-size: 22px;
}
.tb-stat i {
  display: block;
  margin-top: 8px;
  font-style: normal;
  font-size: 13px;
  color: var(--tb-text-3);
}
.tb-feature-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 22px;
}
.tb-feature {
  background: #fff;
  border: 1px solid var(--tb-line);
  border-radius: 20px;
  padding: 30px 28px;
  box-shadow: 0 4px 18px rgba(20, 22, 58, 0.05);
  transition:
    transform 0.3s ease,
    box-shadow 0.3s ease,
    border-color 0.3s ease;
}
.tb-feature:hover {
  transform: translateY(-6px);
  border-color: rgba(84, 87, 232, 0.4);
  box-shadow: 0 22px 48px rgba(84, 87, 232, 0.16);
}
.tb-feature-icon {
  width: 46px;
  height: 46px;
  border-radius: 13px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  margin-bottom: 18px;
  box-shadow: 0 8px 20px rgba(84, 87, 232, 0.22);
}
.tb-feature h3 {
  font-size: 17px;
  font-weight: 700;
  margin: 0 0 8px;
  color: var(--tb-text);
}
.tb-feature p {
  margin: 0;
  font-size: 14px;
  color: var(--tb-text-2);
  line-height: 1.75;
}

/* ============================================================
   大字区
   ============================================================ */
.tb-big {
  background: var(--tb-bg);
  text-align: center;
  overflow: hidden;
}
.tb-big-title {
  font-size: clamp(2rem, 5.4vw, 4.2rem);
  font-weight: 800;
  line-height: 1.16;
  letter-spacing: -0.02em;
  margin: 0 0 20px;
}
.tb-big-sub {
  font-size: clamp(14px, 2vw, 17px);
  color: var(--tb-text-2);
  max-width: 640px;
  margin: 0 auto 44px;
  line-height: 1.75;
}
.tb-words {
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
  gap: 12px;
}
.tb-word {
  font-size: clamp(13px, 2vw, 16px);
  font-weight: 600;
  color: #3a3f66;
  padding: 10px 22px;
  border-radius: 999px;
  background: #fff;
  border: 1px solid var(--tb-line);
  box-shadow: 0 3px 12px rgba(20, 22, 58, 0.05);
  transition: all 0.25s ease;
}
.tb-word:hover {
  color: var(--tb-primary);
  border-color: var(--tb-primary);
  transform: translateY(-2px);
}

/* ============================================================
   版本时间线
   ============================================================ */
.tb-versions {
  background: var(--tb-soft);
}
.tb-version-card {
  background: #fff;
  border: 1px solid var(--tb-line);
  border-radius: 22px;
  padding: 36px 40px;
  box-shadow: 0 4px 18px rgba(20, 22, 58, 0.05);
}
.tb-timeline {
  position: relative;
  padding-left: 28px;
}
.tb-timeline::before {
  content: '';
  position: absolute;
  left: 5px;
  top: 8px;
  bottom: 8px;
  width: 2px;
  background: linear-gradient(180deg, rgba(84, 87, 232, 0.35), rgba(34, 211, 238, 0.2));
}
.tb-version-item {
  position: relative;
  padding-bottom: 22px;
}
.tb-version-item:last-child {
  padding-bottom: 0;
}
.tb-dot {
  position: absolute;
  left: -28px;
  top: 7px;
  width: 12px;
  height: 12px;
  border-radius: 50%;
  background: #fff;
  border: 2px solid var(--tb-primary);
  box-shadow: 0 0 0 4px rgba(84, 87, 232, 0.12);
}
.tb-version-head {
  display: flex;
  align-items: baseline;
  gap: 10px;
  margin-bottom: 2px;
}
.tb-version-head b {
  font-size: 16px;
  font-weight: 700;
  color: var(--tb-text);
}
.tb-version-head span {
  font-size: 12px;
  color: var(--tb-text-3);
}
.tb-version-item p {
  margin: 0;
  font-size: 13px;
  color: var(--tb-text-2);
  line-height: 1.7;
}

/* ============================================================
   集成生态
   ============================================================ */
.tb-eco {
  background: var(--tb-bg);
}
.tb-eco-groups {
  display: flex;
  flex-direction: column;
  gap: 28px;
  max-width: 1000px;
  margin: 0 auto;
}
.tb-group-title {
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 13px;
  font-weight: 700;
  letter-spacing: 0.12em;
  text-transform: uppercase;
  color: var(--tb-text-3);
  margin: 0 0 14px;
}
.tb-group-title::before {
  content: '';
  width: 10px;
  height: 10px;
  border-radius: 3px;
  background: var(--tb-grad);
  flex: none;
}
.tb-group-title::after {
  content: '';
  flex: 1;
  height: 1px;
  background: linear-gradient(90deg, rgba(84, 87, 232, 0.18), transparent);
}
.tb-chips {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
}
.tb-chip {
  display: inline-flex;
  align-items: center;
  padding: 10px 20px;
  border-radius: 999px;
  font-size: 14px;
  font-weight: 600;
  color: #3a3f66;
  background: #fff;
  border: 1px solid var(--tb-line);
  box-shadow: 0 3px 12px rgba(20, 22, 58, 0.05);
  transition: all 0.25s ease;
}
.tb-chip:hover {
  color: var(--tb-primary);
  border-color: var(--tb-primary);
  transform: translateY(-2px);
}

/* ============================================================
   下载区
   ============================================================ */
.tb-download {
  background: var(--tb-soft);
}
.tb-version {
  margin: 0;
  font-size: 12px;
  color: var(--tb-text-3);
}
.tb-note {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  margin: -20px 0 24px;
  font-size: 13px;
  color: var(--tb-text-3);
}
.tb-note button {
  color: var(--tb-primary);
  font-weight: 600;
  background: none;
  border: none;
  cursor: pointer;
}
.tb-note button:hover {
  text-decoration: underline;
}
.tb-download-card {
  background: #fff;
  border: 1px solid var(--tb-line);
  border-radius: 22px;
  overflow: hidden;
  box-shadow: 0 4px 18px rgba(20, 22, 58, 0.05);
}
.tb-download-row {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 18px 22px;
  border-bottom: 1px solid rgba(84, 87, 232, 0.1);
  text-decoration: none;
  transition: background 0.2s ease;
}
.tb-download-row:last-child {
  border-bottom: none;
}
.tb-download-row:hover {
  background: #f7f8ff;
}
.tb-download-row.is-disabled {
  opacity: 0.55;
  cursor: not-allowed;
}
.tb-download-row.is-disabled:hover {
  background: transparent;
}
.tb-dl-icon {
  color: var(--tb-primary);
  display: inline-flex;
}
.tb-dl-text {
  flex: 1;
  text-align: left;
  min-width: 0;
}
.tb-dl-title {
  display: block;
  font-size: 14px;
  font-weight: 700;
  color: var(--tb-text);
}
.tb-dl-note {
  display: block;
  font-size: 12px;
  color: var(--tb-text-3);
}
.tb-dl-badge {
  padding: 3px 10px;
  border-radius: 999px;
  border: 1px solid var(--tb-line);
  font-size: 10px;
  font-weight: 700;
  color: var(--tb-text-3);
  white-space: nowrap;
}
.tb-dl-badge.is-new {
  color: var(--tb-primary);
  background: rgba(84, 87, 232, 0.1);
  border-color: transparent;
}
.tb-dl-arrow {
  color: var(--tb-text-3);
  display: inline-flex;
}
.tb-dl-arrow.is-muted {
  opacity: 0.4;
}

/* Lightbox */
.tb-lightbox {
  position: fixed;
  inset: 0;
  z-index: 100;
  background: rgba(10, 14, 39, 0.9);
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: zoom-out;
}
.tb-lightbox img {
  max-width: 92vw;
  max-height: 90vh;
  border-radius: 12px;
  box-shadow: 0 24px 70px rgba(0, 0, 0, 0.45);
}

/* ============================================================
   响应式
   ============================================================ */
@media (max-width: 900px) {
  .tb-section {
    padding: 56px 0;
  }
  .tb-hero {
    padding: 120px 0 72px;
  }
  .tb-demo-grid,
  .tb-feature-grid {
    grid-template-columns: repeat(2, 1fr);
  }
  .tb-version-card {
    padding: 28px;
  }
}

@media (max-width: 640px) {
  .tb-hero {
    padding: 112px 0 64px;
  }
  .tb-demo-grid,
  .tb-feature-grid {
    grid-template-columns: 1fr;
  }
  .tb-stats {
    gap: 18px 0;
    padding: 20px 16px;
  }
  .tb-stat {
    padding: 0 18px;
  }
  .tb-stat + .tb-stat::before {
    display: none;
  }
  .tb-stat b {
    font-size: 26px;
  }
  .tb-big-title {
    white-space: normal;
  }
  .tb-download-row {
    padding: 16px 18px;
    gap: 12px;
  }
  .tb-dl-badge {
    display: none;
  }
}
</style>

