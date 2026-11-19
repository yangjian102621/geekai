<template>
  <div class="relative group" :class="className">
    <!-- trigger 由父组件通过插槽传入 -->
    <slot />
    <div class="absolute hidden group-hover:block z-30" :class="panelPosition">
      <div class="w-64 rounded-xl border border-hairline bg-canvas p-1.5 shadow-lg text-left">
        <div class="relative group/mac">
          <div class="flex items-center gap-3 px-3 py-2 rounded-xl hover:bg-surface-soft cursor-default">
            <span class="text-muted"><PlatformIcon platform="mac-arm64" /></span>
            <span class="flex-1 text-sm text-body-strong">{{ t('download.mac') }}</span>
            <span v-if="isMac" class="text-[10px] font-semibold text-brand-600">{{ t('download.current') }}</span>
            <svg class="text-muted-soft" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="m9 18 6-6-6-6" />
            </svg>
          </div>
          <div class="absolute hidden group-hover/mac:block z-40" :class="macSubPosition">
            <div class="w-48 rounded-xl border border-hairline bg-canvas p-1.5 shadow-lg text-left">
              <template v-for="option in macOptions" :key="option.key">
                <a v-if="option.href" :href="option.href" class="flex items-center gap-3 px-3 py-2 rounded-xl hover:bg-surface-soft whitespace-nowrap">
                  <span class="flex-1 text-sm text-body-strong">{{ t(option.labelKey) }}</span>
                  <span v-if="option.badgeKey" class="text-[10px] font-semibold text-brand-600">{{ t(option.badgeKey) }}</span>
                </a>
                <div
                  v-else
                  class="flex items-center gap-3 px-3 py-2 rounded-xl opacity-55 cursor-not-allowed whitespace-nowrap"
                  :title="t('download.unavailable')"
                >
                  <span class="flex-1 text-sm text-muted">{{ t(option.labelKey) }}</span>
                  <span class="text-[10px] font-semibold text-muted-soft">{{ t('download.unavailable') }}</span>
                </div>
              </template>
            </div>
          </div>
        </div>
        <a
          v-if="windowsOption.href"
          :href="windowsOption.href"
          class="flex items-center gap-3 px-3 py-2 rounded-xl hover:bg-surface-soft"
        >
          <span class="text-muted"><PlatformIcon platform="windows" /></span>
          <span class="flex-1 text-sm text-body-strong">
            {{ t(windowsOption.labelKey) }}
            <span class="text-muted-soft">（{{ t(windowsOption.noteKey) }}）</span>
          </span>
          <span v-if="windowsOption.badgeKey" class="text-[10px] font-semibold text-brand-600">{{ t(windowsOption.badgeKey) }}</span>
        </a>
        <div
          v-else
          class="flex items-center gap-3 px-3 py-2 rounded-xl opacity-55 cursor-not-allowed"
          :title="t('download.unavailable')"
        >
          <span class="text-muted"><PlatformIcon platform="windows" /></span>
          <span class="flex-1 text-sm text-muted">{{ t(windowsOption.labelKey) }}</span>
          <span class="text-[10px] font-semibold text-muted-soft">{{ t('download.unavailable') }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import PlatformIcon from './PlatformIcon.vue'
import { currentOs, downloadLinks, downloadOptions } from '@/composables/useDownloadLinks'

const props = defineProps({
  // anchor="bottom"：面板在按钮下方（导航右对齐），macOS 二级菜单向左展开。
  // anchor="right"：面板展开在按钮右侧（hero 用），macOS 二级菜单向右展开。
  anchor: { type: String, default: 'bottom' },
  align: { type: String, default: 'right' },
  className: { type: String, default: '' },
})

const { t } = useI18n()

const options = computed(() => downloadOptions(downloadLinks.value))
const windowsOption = computed(() => options.value[0])
const macOptions = computed(() => options.value.slice(1))
const isMac = currentOs() === 'mac'

// 导航面板右对齐时二级菜单向左展开；hero 面板在按钮右侧时二级菜单向右展开。
const panelPosition = computed(() =>
  props.anchor === 'right'
    ? 'top-0 left-full pl-2'
    : `top-full pt-2 ${props.align === 'center' ? 'left-1/2 -translate-x-1/2' : 'right-0'}`,
)
const macSubPosition = computed(() => (props.anchor === 'right' ? 'left-full top-0 pl-2' : 'right-full top-0 pr-2'))
</script>
