<template>
  <div v-reveal class="section-soft min-h-screen">
    <section class="pt-28 pb-16">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 text-center">
        <p class="eyebrow mb-4" data-reveal>{{ t('pricing.badge') }}</p>
        <h1 class="font-serif text-3xl sm:text-[3.25rem] font-normal tracking-[-0.03em] text-ink mb-4" data-reveal>{{ t('pricing.title') }}</h1>
        <p class="max-w-2xl mx-auto text-body" data-reveal>{{ t('pricing.subtitle') }}</p>

        <div class="inline-flex items-center gap-1 mt-8 p-1 bg-cream-strong/60 rounded-full" data-reveal>
          <button class="billing-toggle px-5 py-2 rounded-full text-sm font-medium" :class="{ 'billing-toggle-active': !yearly }" @click="yearly = false">
            {{ t('pricing.monthly') }}
          </button>
          <button class="billing-toggle px-5 py-2 rounded-full text-sm font-medium" :class="{ 'billing-toggle-active': yearly }" @click="yearly = true">
            {{ t('pricing.yearly') }}
          </button>
        </div>
      </div>
    </section>

    <!-- 个人套餐 -->
    <section class="pb-20">
      <div class="max-w-5xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <!-- Free -->
          <div class="pricing-card bg-canvas border border-hairline rounded-xl p-8 flex flex-col" data-reveal>
            <h3 class="font-serif text-xl font-normal text-ink mb-1">{{ t('pricing.free.name') }}</h3>
            <p class="text-sm text-muted mb-4">{{ t('pricing.free.tagline') }}</p>
            <div class="mb-1">
              <span class="font-serif text-[2rem] font-normal tracking-[-0.02em] text-ink">¥0</span>
              <span class="text-sm text-muted">{{ t('pricing.per_month') }}</span>
            </div>
            <p class="text-xs text-brand-600 font-medium mb-6">{{ t('pricing.free.trial_note') }}</p>
            <ul class="space-y-3 text-sm text-body mb-8">
              <li>{{ t('pricing.free.f1') }}{{ t('pricing.free.f1u') }}</li>
              <li>{{ t('pricing.free.f2') }}</li>
              <li>{{ t('pricing.free.f3') }}</li>
              <li>{{ t('pricing.free.f4') }}</li>
              <li>{{ t('pricing.free.f5') }}</li>
            </ul>
            <a :href="downloadUrl ?? '/#download'" class="btn-secondary mt-auto text-center px-5 py-2.5 rounded-lg text-sm font-medium">{{ t('pricing.free.cta') }}</a>
          </div>

          <!-- Pro -->
          <div class="pricing-card pricing-card-featured bg-canvas border-2 rounded-xl p-8 flex flex-col relative" data-reveal>
            <span class="absolute -top-3 left-1/2 -translate-x-1/2 px-3 py-1 rounded-full bg-ink text-on-dark text-xs font-medium">{{ t('pricing.pro.badge') }}</span>
            <h3 class="font-serif text-xl font-normal text-ink mb-1">{{ t('pricing.pro.name') }}</h3>
            <p class="text-sm text-muted mb-4">{{ t('pricing.pro.tagline') }}</p>
            <div class="mb-1">
              <span class="font-serif text-[2rem] font-normal tracking-[-0.02em] text-ink">¥{{ yearly ? proYearly : proMonthly }}</span>
              <span class="text-sm text-muted">{{ yearly ? '/年' : t('pricing.per_month') }}</span>
            </div>
            <p class="text-xs text-brand-600 font-medium mb-6">{{ t('pricing.pro.trial_note') }}</p>
            <ul class="space-y-3 text-sm text-body mb-8">
              <li>{{ t('pricing.pro.f1') }}</li>
              <li class="text-brand-600 font-medium">{{ t('pricing.pro.f1b') }}</li>
              <li>{{ t('pricing.pro.f2') }}</li>
              <li>{{ t('pricing.pro.f3') }}</li>
              <li>{{ t('pricing.pro.f4') }}</li>
              <li>{{ t('pricing.pro.f5') }}</li>
              <li>{{ t('pricing.pro.f6') }}</li>
            </ul>
            <a :href="downloadUrl ?? '/#download'" class="btn-primary mt-auto text-center px-5 py-2.5 rounded-lg text-sm font-medium">{{ t('pricing.pro.cta') }}</a>
          </div>
        </div>
      </div>
    </section>

    <!-- FAQ -->
    <section class="py-16 bg-canvas">
      <div class="max-w-3xl mx-auto px-4 sm:px-6 lg:px-8">
        <h2 class="text-center font-serif text-2xl sm:text-[2rem] font-normal tracking-[-0.025em] text-ink mb-8" data-reveal>{{ t('pricing.faq.title') }}</h2>
        <div class="space-y-3">
          <details v-for="q in faqKeys" :key="q" class="faq-item bg-canvas border border-hairline rounded-xl p-5" data-reveal>
            <summary class="font-semibold text-ink cursor-pointer">{{ t(`pricing.faq.${q}`) }}</summary>
            <p class="mt-3 text-sm text-body leading-relaxed">{{ t(`pricing.faq.a${q.slice(1)}`) }}</p>
          </details>
        </div>
        <div class="text-center mt-8">
          <RouterLink
            to="/"
            class="inline-flex items-center justify-center px-6 py-3 rounded-lg text-sm font-medium text-white bg-primary hover:bg-primary-hover transition"
          >
            返回首页
          </RouterLink>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup>
import { computed, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useReveal } from '@/composables/useReveal'
import { useLatestSetup } from '@/composables/useDownloadLinks'

const vReveal = useReveal()
const { t } = useI18n()
const yearly = ref(false)
const downloadUrl = useLatestSetup()

const proMonthly = 10
const proYearly = 96

const faqKeys = ['q1', 'q2', 'q3', 'q4', 'q5']
</script>
