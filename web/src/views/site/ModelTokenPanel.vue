<template>
  <section class="rounded-xl border border-hairline bg-canvas px-6 py-5" aria-labelledby="model-api-title">
    <h2 id="model-api-title" class="flex items-center gap-2 text-base font-semibold text-ink">
      <KeyRound :size="18" aria-hidden="true" />{{ t('model_api.title') }}
    </h2>
    <p class="mt-2 text-sm leading-6 text-muted">{{ t('model_api.description') }}</p>
    <dl class="mt-4 space-y-3 text-sm">
      <div>
        <dt class="text-muted">{{ t('model_api.base_url') }}</dt>
        <dd class="mt-1 select-all break-all font-mono text-ink">{{ origin }}/v1</dd>
      </div>
      <div>
        <dt class="text-muted">{{ t('model_api.protocol') }}</dt>
        <dd class="mt-1 text-ink">OpenAI Chat Completions</dd>
      </div>
    </dl>
    <p v-if="loading" role="status" class="mt-4 text-sm text-muted">{{ t('space.loading') }}</p>
    <div v-else-if="loadError" role="alert" class="mt-4 flex items-center gap-3 text-sm text-rose-700">
      {{ t('space.load_error') }}
      <button type="button" class="underline" @click="load">{{ t('space.retry') }}</button>
    </div>
    <template v-else>
      <p class="mt-4 text-sm text-muted">{{ t(enabled ? 'model_api.enabled' : 'model_api.disabled') }}</p>
      <div v-if="secret" class="mt-3">
        <label for="model-api-token" class="text-sm font-medium text-ink">{{ t('model_api.token') }}</label>
        <div class="mt-2 flex gap-2">
          <input
            id="model-api-token"
            :value="secret"
            readonly
            class="min-w-0 flex-1 rounded-lg border border-hairline bg-surface-soft px-3 py-2 font-mono text-sm text-ink"
            @focus="$event.currentTarget.select()"
          />
          <button
            type="button"
            class="inline-flex shrink-0 items-center gap-1.5 rounded-lg border border-hairline px-3 py-2 text-sm text-ink hover:bg-surface-soft"
            @click="copySecret"
          >
            <Copy :size="15" aria-hidden="true" />{{ t(copied ? 'model_api.copied' : 'model_api.copy') }}
          </button>
        </div>
        <p class="mt-2 text-xs leading-5 text-muted">{{ t('model_api.save_once') }}</p>
      </div>
      <div class="mt-4 flex flex-wrap gap-3">
        <button
          type="button"
          :disabled="busy"
          class="btn-primary rounded-lg px-4 py-2 text-sm font-medium disabled:opacity-50"
          @click="changeToken(false)"
        >
          {{ t(busy ? 'model_api.processing' : enabled ? 'model_api.reset' : 'model_api.create') }}
        </button>
        <button
          v-if="enabled"
          type="button"
          :disabled="busy"
          class="rounded-lg border border-hairline px-4 py-2 text-sm text-muted hover:bg-surface-soft disabled:opacity-50"
          @click="changeToken(true)"
        >
          {{ t('model_api.revoke') }}
        </button>
      </div>
    </template>
    <p v-if="error" role="alert" class="mt-3 text-sm text-rose-700">{{ error }}</p>
  </section>
</template>

<script setup>
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { Copy, KeyRound } from 'lucide-vue-next'
import { fetchModelToken, revokeModelToken, rotateModelToken } from '@/services/site'

const props = defineProps({
  userId: { type: Number, required: true },
})

const { t } = useI18n()

const loading = ref(true)
const loadError = ref(false)
const enabled = ref(false)
const secret = ref('')
const busy = ref(false)
const error = ref('')
const copied = ref(false)
const origin = window.location.origin

async function load() {
  loading.value = true
  loadError.value = false
  try {
    const status = await fetchModelToken()
    enabled.value = !!status.enabled
  } catch (e) {
    loadError.value = true
  } finally {
    loading.value = false
  }
}

async function changeToken(revoke) {
  if (enabled.value && !window.confirm(t(revoke ? 'model_api.confirm_revoke' : 'model_api.confirm_reset'))) return
  busy.value = true
  error.value = ''
  secret.value = ''
  copied.value = false
  try {
    if (revoke) {
      await revokeModelToken()
      enabled.value = false
    } else {
      const result = await rotateModelToken()
      // 明文只留在当前组件中，不写入持久存储或查询缓存。
      secret.value = result.token
      enabled.value = true
    }
  } catch (e) {
    error.value = t('model_api.action_error')
  } finally {
    busy.value = false
  }
}

async function copySecret() {
  try {
    await navigator.clipboard.writeText(secret.value)
    copied.value = true
  } catch (e) {
    error.value = t('model_api.copy_error')
  }
}

load()
</script>
