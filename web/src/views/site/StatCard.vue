<template>
  <div class="rounded-xl border border-hairline bg-canvas p-5">
    <div class="flex items-center gap-2 text-muted">
      <span class="flex h-9 w-9 items-center justify-center rounded-xl bg-brand-50 text-brand-600">
        <component :is="icon" :size="20" />
      </span>
      <span class="text-sm">{{ label }}</span>
    </div>
    <div class="mt-4 font-serif text-[1.75rem] font-normal tracking-[-0.02em] text-ink">
      <span v-if="loading" class="inline-block h-7 w-24 animate-pulse rounded bg-surface-card" />
      <template v-else>{{ value }}</template>
    </div>
  </div>
</template>

<script setup>
defineProps({
  icon: { type: [Object, Function], required: true },
  label: { type: String, required: true },
  value: { type: String, required: true },
  loading: { type: Boolean, default: false },
})
</script>
