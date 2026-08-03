<template>
  <div v-if="list.length" class="space-y-3 mb-2">
    <div
      v-for="(f, idx) in list"
      :key="f.url || idx"
      class="flex !items-start flex-col justify-center gap-3 p-3 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg"
    >
      <!-- Image -->
      <div
        v-if="isImageFile(f)"
        class="flex-shrink-0 bg-gray-100 max-w-[500px] max-h-[500px] dark:bg-gray-700 rounded-lg overflow-hidden"
      >
        <el-image
          :src="f.url"
          fit="cover"
          :preview-src-list="imageUrls"
          :initial-index="imageIndexMap.get(f.url) || 0"
          hide-on-click-modal
          :z-index="3000"
          class="w-full h-full"
        />
      </div>

      <!-- Video -->
      <div
        v-else-if="isVideoFile(f)"
        class="flex-shrink-0 bg-gray-100 max-w-[500px] dark:bg-gray-700 rounded-lg overflow-hidden"
      >
        <video :src="f.url" controls preload="metadata" class="w-full h-full object-cover"></video>
      </div>

      <!-- Audio -->
      <div
        v-else-if="isAudioFile(f)"
        class="flex-shrink-0 bg-gray-100 dark:bg-gray-700 rounded-lg flex items-center justify-center"
      >
        <audio :src="f.url" controls preload="metadata"></audio>
      </div>

      <!-- Other Files -->
      <div
        v-else
        class="flex-shrink-0 w-20 h-20 bg-gray-100 dark:bg-gray-700 rounded-lg flex items-center justify-center"
      >
        <img :src="GetFileIcon(extOf(f))" class="w-20 h-20" alt="file" />
      </div>

      <div class="flex w-full flex-row min-w-0 justify-between items-center">
        <!-- File Info -->
        <div class="text-xs text-gray-500 dark:text-gray-400 mt-1">
          {{ extOf(f).replace('.', '').toUpperCase() || 'FILE' }} ·
          {{ FormatFileSize(f.size || 0) }}
        </div>

        <!-- Download Button -->
        <div class="flex-shrink-0">
          <el-tooltip class="box-item" effect="dark" content="下载" placement="top">
            <i
              class="iconfont icon-download !text-sm cursor-pointer"
              v-if="!f.downloading"
              @click="downloadFile(f)"
            ></i>
            <el-image src="/images/loading.gif" class="w-4 h-4" fit="cover" v-else />
          </el-tooltip>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { FormatFileSize, GetFileIcon } from '@/store/system'
import { computed } from 'vue'
import { httpDownload } from '@/utils/http'
import { replaceImg } from '@/utils/libs'

const props = defineProps({
  files: {
    type: Array,
    default: () => [],
  },
})

const list = computed(() => props.files || [])

const normalizeExt = (ext) => (ext || '').toLowerCase()
const urlExt = (url) => {
  if (!url) return ''
  try {
    const path = url.split('?')[0]
    const dot = path.lastIndexOf('.')
    return dot >= 0 ? path.substring(dot) : ''
  } catch (_) {
    return ''
  }
}
const extOf = (f) => normalizeExt(f.ext || urlExt(f.url))

const IMAGE_EXTS = new Set(['.png', '.jpg', '.jpeg', '.gif', '.webp', '.bmp', '.svg'])
const VIDEO_EXTS = new Set(['.mp4', '.webm', '.ogg', '.mov', '.m4v'])
const AUDIO_EXTS = new Set(['.mp3', '.wav', '.ogg', '.aac', '.m4a'])

const isImageFile = (f) => IMAGE_EXTS.has(extOf(f))
const isVideoFile = (f) => VIDEO_EXTS.has(extOf(f))
const isAudioFile = (f) => AUDIO_EXTS.has(extOf(f))

// image preview urls and index mapping for el-image gallery
const imageUrls = computed(() => list.value.filter((f) => isImageFile(f)).map((f) => f.url))
const imageIndexMap = computed(() => {
  const map = new Map()
  imageUrls.value.forEach((u, i) => map.set(u, i))
  return map
})

const downloadFile = async (item) => {
  const url = replaceImg(item.url)
  const downloadURL = `/api/download?url=${url}`
  const urlObj = new URL(url)
  const fileName = urlObj.pathname.split('/').pop()

  item.downloading = true

  try {
    const response = await httpDownload(downloadURL)
    const blob = new Blob([response.data])
    const link = document.createElement('a')
    link.href = URL.createObjectURL(blob)
    link.download = fileName
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    URL.revokeObjectURL(link.href)
    item.downloading = false
  } catch (error) {
    showMessageError('下载失败')
    item.downloading = false
  }
}
</script>
