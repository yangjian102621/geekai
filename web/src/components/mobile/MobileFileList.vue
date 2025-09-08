<template>
  <div v-if="items.length" :class="[containerClass, direction === 'col' ? '!items-end' : '']">
    <div
      v-for="(f, idx) in items"
      :key="f.url || idx"
      :class="[
        'relative inline-flex items-center border border-gray-200 rounded-xl bg-white dark:bg-[#2b2b2b] dark:border-gray-700',
        isImage(f.ext) ? 'p-0' : 'p-2',
      ]"
    >
      <div v-if="isImage(f.ext)" class="relative w-[56px] h-[56px] overflow-hidden rounded-lg">
        <img :src="f.url" alt="img" class="w-full h-full object-cover" />
      </div>
      <div v-else :class="['flex items-center', direction === 'col' ? 'w-full' : 'max-w-[240px]']">
        <img :src="GetFileIcon(f.ext)" class="w-10 h-10 mr-2" />
        <div class="min-w-0 flex flex-col items-center gap-1 text-sm">
          <a
            :href="f.url"
            target="_blank"
            class="truncate block"
            :class="direction === 'col' ? 'text-base max-w-full' : 'max-w-[180px]'"
            >{{ f.name }}</a
          >
          <div class="text-xs flex w-full justify-start text-gray-500">
            {{ extUpper(f.ext) }} · {{ FormatFileSize(f.size || 0) }}
          </div>
        </div>
      </div>
      <button
        v-if="removable"
        class="absolute -right-2 -top-2 w-5 h-5 rounded-full bg-rose-600 text-white flex items-center justify-center text-[10px]"
        @click="onRemove(f, idx)"
      >
        ×
      </button>
    </div>
  </div>
</template>

<script setup>
import { FormatFileSize, GetFileIcon } from '@/store/system'
import { isImage } from '@/utils/libs'

const props = defineProps({
  files: {
    type: Array,
    default: () => [],
  },
  removable: {
    type: Boolean,
    default: false,
  },
  direction: {
    type: String,
    default: 'row',
  },
})
const emits = defineEmits(['remove'])
const items = computed(() => props.files || [])
const onRemove = (f, idx) => emits('remove', { file: f, index: idx })
const direction = computed(() => props.direction)
const containerClass = computed(() =>
  props.direction === 'col' ? 'flex flex-col gap-2' : 'flex flex-wrap gap-2 px-2 pt-2'
)
const extUpper = (ext) => (ext || '').replace('.', '').toUpperCase() || 'FILE'
</script>

<style scoped></style>
