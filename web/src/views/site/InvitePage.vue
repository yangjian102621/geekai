<template>
  <!-- 邀请页：访客浏览奖励规则；登录用户查看邀请码/二维码/统计；
       带 aff 参数的访问者走注册引导视图。 -->
  <div class="invite-theme">
    <!-- 加载中（登录态） -->
    <section v-if="authLoading" class="invite-hero py-32">
      <div class="invite-hero-bg" />
      <div class="relative max-w-3xl mx-auto px-4 text-center" role="status">
        <div class="invite-loading-spinner mx-auto" />
        <p class="mt-4 text-body">{{ t('invite.loading') }}</p>
      </div>
    </section>

    <!-- 好友邀请视图（?aff=） -->
    <section v-else-if="aff" class="invite-hero pt-28 pb-20">
      <div class="invite-hero-bg" />
      <div class="relative max-w-3xl mx-auto px-4 text-center">
        <div class="invite-hero-badge">
          <span class="w-1.5 h-1.5 rounded-full bg-coral animate-pulse" />
          <span>{{ t('invite.invitee.badge') }}</span>
        </div>
        <h1 class="invite-hero-title">{{ t('invite.invitee.title') }}</h1>
        <p class="invite-hero-subtitle">{{ t('invite.invitee.desc') }}</p>
        <div class="invite-hero-cta-row">
          <template v-if="user">
            <template v-if="macUnknown">
              <a :href="macArm64.href ?? '/#download'" class="invite-cta-primary">
                <Download :size="20" />
                {{ t('download.mac_arm64') }}
              </a>
              <a :href="macX64.href ?? '/#download'" class="invite-cta-primary">
                <Download :size="20" />
                {{ t('download.mac_intel') }}
              </a>
            </template>
            <a v-else :href="downloadUrl" class="invite-cta-primary">
              <Download :size="20" />
              {{ t('invite.invitee.download_btn') }}
            </a>
          </template>
          <button v-else type="button" class="invite-cta-primary" @click="openLoginPage">
            <LinkIcon :size="20" />
            {{ t('invite.invitee.login_btn') }}
          </button>
        </div>
        <p class="text-sm text-body">
          {{ user ? t('invite.invitee.logged_in_hint') : t('invite.invitee.login_note') }}
        </p>
        <p class="text-xs text-muted mt-2">{{ t('invite.invitee.download_note') }}</p>
        <p v-if="authError" role="alert" class="mt-6 text-body-strong">{{ t('invite.auth_error') }}</p>
      </div>
    </section>

    <!-- 未登录视图 -->
    <template v-else-if="!user">
      <section class="invite-hero pt-28 pb-20">
        <div class="invite-hero-bg" />
        <div class="relative max-w-4xl mx-auto px-4 text-center">
          <div class="invite-hero-badge">
            <span class="w-1.5 h-1.5 rounded-full bg-coral animate-pulse" />
            <span>{{ t('invite.guest.badge') }}</span>
          </div>
          <h1 class="invite-hero-title">
            {{ t('invite.title') }}
            <br />
            <span class="invite-hero-title-accent">{{ t('invite.subtitle', { points: rewardAmount.toLocaleString() }) }}</span>
          </h1>
          <p class="invite-hero-subtitle">
            {{ t('invite.reward.per_invite', { points: rewardAmount.toLocaleString(), limit: inviteLimit, total: rewardTotal.toLocaleString() }) }}
          </p>
          <div class="invite-hero-cta-row">
            <button type="button" class="invite-cta-primary" @click="openLoginPage">
              <LinkIcon :size="20" />
              {{ t('invite.guest.login_btn') }}
            </button>
          </div>
          <div class="invite-hero-stats">
            <div class="text-center">
              <div class="invite-hero-stat-num">{{ rewardAmount.toLocaleString() }}</div>
              <div class="invite-hero-stat-label">{{ t('invite.reward.unit') }}</div>
            </div>
            <div class="invite-hero-stat-divider" />
            <div class="text-center">
              <div class="invite-hero-stat-num">{{ inviteLimit }}</div>
              <div class="invite-hero-stat-label">{{ t('invite.reward.limit_label') }}</div>
            </div>
            <div class="invite-hero-stat-divider" />
            <div class="text-center">
              <div class="invite-hero-stat-num">{{ rewardMax.toLocaleString() }}</div>
              <div class="invite-hero-stat-label">{{ t('invite.reward.max_label') }}</div>
            </div>
          </div>
        </div>
      </section>

      <StepsSection :reward="reward" />
      <RewardsSection :reward="reward" :limit="inviteLimit" />
      <FaqSection :reward="reward" />
      <BottomCta />
    </template>

    <!-- 数据加载中（已登录，尚无缓存数据） -->
    <section v-else-if="summaryLoading && !summary" class="invite-hero py-32">
      <div class="invite-hero-bg" />
      <div class="relative max-w-3xl mx-auto px-4 text-center" role="status">
        <div class="invite-loading-spinner mx-auto" />
        <p class="mt-4 text-body">{{ t('invite.loading') }}</p>
      </div>
    </section>

    <!-- 数据加载失败 -->
    <section v-else-if="summaryError || !summary" class="invite-hero py-32">
      <div class="invite-hero-bg" />
      <div class="relative max-w-3xl mx-auto px-4 text-center" role="alert">
        <p class="text-body-strong mb-5">{{ t('invite.error') }}</p>
        <button type="button" class="btn-primary px-5 py-2 rounded-lg" @click="refetch">
          {{ t('invite.retry') }}
        </button>
      </div>
    </section>

    <!-- 已登录视图 -->
    <template v-else>
      <section class="invite-hero pt-28 pb-20">
        <div class="invite-hero-bg" />
        <div class="relative max-w-4xl mx-auto px-4 text-center">
          <div class="invite-hero-badge">
            <span class="w-1.5 h-1.5 rounded-full bg-coral animate-pulse" />
            <span>{{ t('invite.guest.badge') }}</span>
          </div>
          <h1 class="invite-hero-title">
            {{ t('invite.title') }}
            <br />
            <span class="invite-hero-title-accent">{{ t('invite.subtitle', { points: rewardAmount.toLocaleString() }) }}</span>
          </h1>
          <div class="invite-hero-stats">
            <div class="text-center">
              <div class="invite-hero-stat-num">{{ summary.invited_count.toLocaleString() }}</div>
              <div class="invite-hero-stat-label">{{ t('invite.stats.invited') }}</div>
            </div>
            <div class="invite-hero-stat-divider" />
            <div class="text-center">
              <div class="invite-hero-stat-num">{{ summary.total_invite_points.toLocaleString() }}</div>
              <div class="invite-hero-stat-label">{{ t('invite.stats.reward') }}</div>
            </div>
            <div class="invite-hero-stat-divider" />
            <div class="text-center">
              <div class="invite-hero-stat-num">{{ summary.today_invites.toLocaleString() }}</div>
              <div class="invite-hero-stat-label">{{ t('invite.stats.today') }}</div>
            </div>
          </div>
        </div>
      </section>

      <section class="invite-share-section">
        <div class="max-w-4xl mx-auto px-4">
          <div class="invite-share-card">
            <div class="invite-share-qr-col">
              <div class="invite-qr-wrapper">
                <canvas ref="qrCanvas" role="img" :aria-label="summary.invite_url" />
              </div>
              <p class="text-xs text-body mt-3">{{ t('invite.qr_hint') }}</p>
            </div>
            <div class="invite-share-link-col">
              <label for="invite-link" class="invite-link-label">{{ t('invite.share_title') }}</label>
              <div class="invite-share-link-row">
                <input
                  id="invite-link"
                  :value="summary.invite_url"
                  readonly
                  class="invite-share-link-input"
                  @focus="$event.target.select()"
                />
                <button type="button" class="invite-copy-btn-big" @click="copyShareLink">
                  <Check v-if="copied" :size="16" />
                  <Copy v-else :size="16" />
                  {{ copied ? t('invite.copied') : t('invite.copy') }}
                </button>
              </div>
              <p v-if="shareCopyError" role="alert" class="text-sm text-body-strong mt-3">{{ t('invite.copy_manual') }}</p>
              <p class="text-sm text-body mt-3">{{ t('invite.share_hint') }}</p>
            </div>
          </div>
        </div>
      </section>

      <StepsSection :reward="reward" />
      <RewardsSection :reward="reward" :limit="inviteLimit" />
      <FaqSection :reward="reward" />
      <BottomCta :url="summary.invite_url" />
    </template>
  </div>
</template>

<script setup>
import { computed, defineComponent, h, ref, watch, nextTick, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import QRCode from 'qrcode'
import { Check, Copy, Download, Gift, Heart, Link as LinkIcon, Lock } from 'lucide-vue-next'
import { useAuth } from '@/composables/useAuth'
import { useInviteReward, useInviteSummary } from '@/composables/useInvite'
import {
  currentMacArch,
  currentOs,
  downloadLinks,
  downloadOptions,
  recommendedHref,
} from '@/composables/useDownloadLinks'
import { openLoginPage } from '@/services/site'

const { t } = useI18n()
const route = useRoute()
const aff = computed(() => route.query.aff?.trim() || '')

const { user, isLoading: authLoading, isError: authError } = useAuth()
const { reward } = useInviteReward()

/* 邀请活动规则：每位奖励积分、邀请上限 30 位 */
const inviteLimit = 30
const rewardAmount = computed(() =>
  typeof reward.value === 'number' && reward.value > 0 ? reward.value : 500,
)
const rewardMax = computed(() => rewardAmount.value * inviteLimit)
const rewardTotal = computed(() => rewardAmount.value * inviteLimit * 2)
const { summary, isLoading: summaryLoading, isError: summaryError, refetch } = useInviteSummary(!aff.value && !!user.value)

/* InviteeView 的下载选项（mac 架构未知时拆成两个按钮） */
const macUnknown = currentOs() === 'mac' && currentMacArch() === 'unknown'
const options = computed(() => downloadOptions(downloadLinks.value))
const macArm64 = computed(() => options.value[1])
const macX64 = computed(() => options.value[2])
const downloadUrl = computed(() => recommendedHref(downloadLinks.value) ?? '/#download')

async function writeClipboard(text) {
  if (navigator.clipboard?.writeText) {
    await navigator.clipboard.writeText(text)
    return
  }
  const input = document.createElement('input')
  input.value = text
  document.body.appendChild(input)
  input.select()
  document.execCommand('copy')
  document.body.removeChild(input)
}

/* 邀请链接复制（已登录分享卡片） */
const copied = ref(false)
const shareCopyError = ref(false)

async function copyShareLink() {
  if (!summary.value?.invite_url) return
  try {
    await writeClipboard(summary.value.invite_url)
    copied.value = true
    shareCopyError.value = false
  } catch (e) {
    shareCopyError.value = true
  }
}

/* ---------- 三步流程 / 奖励卡片 / FAQ / 底部 CTA 子区块 ---------- */

const StepsSection = defineComponent({
  name: 'StepsSection',
  props: { reward: { type: Number, default: null } },
  setup(props) {
    const points = props.reward != null && props.reward > 0 ? props.reward : 500
    const steps = [1, 2, 3].map((num) => ({
      num,
      title: t(`invite.steps.step${num}_title`),
      desc: t(`invite.steps.step${num}_desc`, { points: points.toLocaleString() }),
    }))
    return () =>
      h('section', { class: 'invite-steps-section py-20 lg:py-24' }, [
        h('div', { class: 'max-w-5xl mx-auto px-4 text-center' }, [
          h('h2', { class: 'invite-section-title' }, t('invite.steps.title')),
          h('p', { class: 'invite-section-subtitle mb-10' }, t('invite.steps.subtitle')),
          h('div', { class: 'grid grid-cols-1 md:grid-cols-3 gap-6' }, steps.map((step) =>
            h('div', { key: step.num, class: 'invite-step-card' }, [
              h('div', { class: 'invite-step-num' }, String(step.num)),
              h('div', { class: 'invite-step-title' }, step.title),
              h('div', { class: 'invite-step-desc' }, step.desc),
            ]),
          )),
        ]),
      ])
  },
})

const RewardsSection = defineComponent({
  name: 'RewardsSection',
  props: {
    reward: { type: Number, default: null },
    limit: { type: Number, default: 30 },
  },
  setup(props) {
    const points = props.reward != null && props.reward > 0 ? props.reward : 500
    const params = {
      points: points.toLocaleString(),
      limit: props.limit,
      total: (points * props.limit).toLocaleString(),
    }
    return () =>
      h('section', { class: 'invite-rewards-section py-20 lg:py-24' }, [
        h('div', { class: 'max-w-5xl mx-auto px-4' }, [
          h('div', { class: 'text-center mb-10' }, [
            h('h2', { class: 'invite-section-title' }, t('invite.rewards.title')),
            h('p', { class: 'invite-section-subtitle' }, t('invite.rewards.subtitle')),
          ]),
          h('div', { class: 'grid grid-cols-1 md:grid-cols-3 gap-6' }, [
            h('div', { class: 'invite-big-card invite-big-card-inviter' }, [
              h('div', { class: 'invite-big-card-top' }, [
                h('span', { class: 'invite-big-card-icon' }, [h(Gift, { size: 20 })]),
                h('h2', { class: 'invite-big-card-title' }, t('invite.rewards.inviter_title')),
              ]),
              h('div', { class: 'invite-big-card-metric' }, [
                h('div', { class: 'invite-big-card-num' }, points.toLocaleString()),
                h('div', { class: 'invite-big-card-unit' }, t('invite.rewards.unit_per_person', params)),
              ]),
              h('p', { class: 'invite-big-card-desc' }, t('invite.rewards.inviter_desc', params)),
            ]),
            h('div', { class: 'invite-big-card invite-big-card-invitee' }, [
              h('div', { class: 'invite-big-card-top' }, [
                h('span', { class: 'invite-big-card-icon' }, [h(Heart, { size: 20 })]),
                h('h2', { class: 'invite-big-card-title' }, t('invite.rewards.invitee_title')),
              ]),
              h('div', { class: 'invite-big-card-metric' }, [
                h('div', { class: 'invite-big-card-num' }, points.toLocaleString()),
                h('div', { class: 'invite-big-card-unit' }, t('invite.rewards.unit_new_user', params)),
              ]),
              h('p', { class: 'invite-big-card-desc' }, t('invite.rewards.invitee_desc', params)),
            ]),
            h('div', { class: 'invite-big-card-wrap' }, [
              h('span', { class: 'invite-big-card-badge' }, t('invite.rewards.both_badge')),
              h('div', { class: 'invite-big-card invite-big-card-both' }, [
                h('div', { class: 'invite-big-card-top' }, [
                  h('span', { class: 'invite-big-card-icon' }, [h(Lock, { size: 20 })]),
                  h('h2', { class: 'invite-big-card-title' }, t('invite.rewards.both_title')),
                ]),
                h('div', { class: 'invite-big-card-metric' }, [
                  h('div', { class: 'invite-big-card-num' }, points.toLocaleString()),
                  h('div', { class: 'invite-big-card-unit' }, t('invite.rewards.unit_per_both', params)),
                ]),
                h('p', { class: 'invite-big-card-desc' }, t('invite.rewards.both_desc', params)),
              ]),
            ]),
          ]),
        ]),
      ])
  },
})

const FaqSection = defineComponent({
  name: 'FaqSection',
  props: { reward: { type: Number, default: null } },
  setup(props) {
    const faqs = [
      { q: t('invite.faq.q1'), a: t('invite.faq.a1') },
      {
        q: t('invite.faq.q2'),
        a:
          props.reward != null && props.reward > 0
            ? t('invite.faq.a2_with_reward', { points: props.reward.toLocaleString() })
            : t('invite.faq.a2'),
      },
      { q: t('invite.faq.q3'), a: t('invite.faq.a3') },
      { q: t('invite.faq.q4'), a: t('invite.faq.a4') },
    ]
    return () =>
      h('section', { class: 'invite-faq-section py-20 lg:py-24' }, [
        h('div', { class: 'max-w-3xl mx-auto px-4' }, [
          h('h2', { class: 'invite-section-title text-center' }, t('invite.faq.title')),
          h('div', { class: 'invite-faq-list mt-8' }, faqs.map((faq, index) =>
            h('details', { key: faq.q, class: 'invite-faq-item', open: index === 0 }, [
              h('summary', { class: 'invite-faq-q' }, faq.q),
              h('p', { class: 'invite-faq-a' }, faq.a),
            ]),
          )),
        ]),
      ])
  },
})

const BottomCta = defineComponent({
  name: 'BottomCta',
  props: { url: { type: String, default: '' } },
  setup(props) {
    const copied = ref(false)
    const copy = async () => {
      if (!props.url) return
      try {
        await writeClipboard(props.url)
        copied.value = true
        window.setTimeout(() => (copied.value = false), 2000)
      } catch (e) {
        // 剪贴板被拒绝时保留原始页面，不打断用户操作
      }
    }
    return () =>
      h('section', { class: 'invite-bottom-cta py-20 lg:py-24 text-center' }, [
        h('div', { class: 'max-w-3xl mx-auto px-4' }, [
          h('h2', { class: 'invite-bottom-cta-title' }, t('invite.bottom_cta_title')),
          h('p', { class: 'invite-bottom-cta-sub' }, t('invite.bottom_cta_sub')),
          props.url
            ? h('button', { type: 'button', class: 'invite-cta-primary', onClick: () => void copy() }, [
                h(LinkIcon, { size: 20 }),
                copied.value ? t('invite.copied') : t('invite.bottom_cta_copy'),
              ])
            : h('button', { type: 'button', class: 'invite-cta-primary', onClick: openLoginPage }, [
                h(LinkIcon, { size: 20 }),
                t('invite.guest.login_btn'),
              ]),
        ]),
      ])
  },
})

/* 邀请二维码（qrcode 包 canvas 渲染，等价于原 QRCodeCanvas） */
const qrCanvas = ref(null)

async function renderQr() {
  if (!summary.value?.invite_url || !qrCanvas.value) return
  try {
    await QRCode.toCanvas(qrCanvas.value, summary.value.invite_url, {
      width: 180,
      margin: 0,
      errorCorrectionLevel: 'M',
    })
  } catch (e) {
    // 二维码渲染失败时静默忽略
  }
}

watch(summary, async (val) => {
  if (val) {
    await nextTick()
    renderQr()
  }
})

onMounted(async () => {
  await nextTick()
  renderQr()
})
</script>
