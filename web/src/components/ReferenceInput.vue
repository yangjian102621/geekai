<template>
  <div class="reference-input w-full">
    <el-input
      ref="inputRef"
      v-model="innerValue"
      type="textarea"
      :autosize="autosize"
      :maxlength="maxlength"
      :show-word-limit="showWordLimit"
      :placeholder="placeholder"
      :disabled="disabled"
      :rows="rows"
      v-loading="loading"
      @input="handleInput"
      @keydown.esc="showImageMentionPopup = false"
    />

    <!-- @ 触发的参考图列表弹窗，定位到光标旁 -->
    <Teleport to="body">
      <div
        v-if="showImageMentionPopup && imageListForPrompt.length > 0"
        ref="imageMentionPopupRef"
        class="image-mention-popup"
        :style="{ left: mentionPopupLeft + 'px', top: mentionPopupTop + 'px' }"
      >
        <div
          v-for="(item, idx) in imageListForPrompt"
          :key="idx"
          class="image-mention-item"
          @click="onSelectImageMention(item)"
        >
          <img :src="item.url" class="image-mention-thumb" alt="" @error="onThumbError" />
          <span class="image-mention-name">{{ item.name }}</span>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup>
import { Teleport, computed, nextTick, onUnmounted, ref, watch } from 'vue'

const props = defineProps({
  modelValue: {
    type: String,
    default: '',
  },
  imageList: {
    type: Array,
    default: () => [],
  },
  placeholder: {
    type: String,
    default: '请在此输入绘画提示词，输入 @ 可引用参考图（图1、图2...）',
  },
  autosize: {
    type: Object,
    default: () => ({ minRows: 4, maxRows: 6 }),
  },
  maxlength: {
    type: Number,
    default: 1024,
  },
  showWordLimit: {
    type: Boolean,
    default: true,
  },
  loading: {
    type: Boolean,
    default: false,
  },
  disabled: {
    type: Boolean,
    default: false,
  },
  rows: {
    type: Number,
    default: 4,
  },
})

const emit = defineEmits(['update:modelValue'])

// 受控 v-model
const innerValue = computed({
  get() {
    return props.modelValue
  },
  set(val) {
    emit('update:modelValue', val)
  },
})

const inputRef = ref(null)

// @ 引用参考图：弹窗显隐与光标定位
const showImageMentionPopup = ref(false)
const mentionPopupLeft = ref(0)
const mentionPopupTop = ref(0)
const triggerIndex = ref(0)
const imageMentionPopupRef = ref(null)

// 参考图列表 for 提示词（图 1、图 2... + url）
const imageListForPrompt = computed(() => {
  const imgs = Array.isArray(props.imageList) ? props.imageList.filter(Boolean) : []
  return imgs.map((url, i) => ({ name: `图${i + 1}`, url }))
})

// 获取 textarea 光标处的视口像素坐标（用于弹窗定位）
function getCaretPixelPosition(textarea, caretIndex) {
  if (!textarea || caretIndex == null) return { left: 0, top: 0 }
  const style = getComputedStyle(textarea)
  const mirror = document.createElement('div')
  const mirrorStyle = mirror.style
  mirrorStyle.position = 'fixed'
  mirrorStyle.left = '-9999px'
  mirrorStyle.top = '0'
  mirrorStyle.visibility = 'hidden'
  mirrorStyle.whiteSpace = 'pre-wrap'
  mirrorStyle.wordWrap = 'break-word'
  mirrorStyle.width = style.width
  mirrorStyle.font = style.font
  mirrorStyle.fontSize = style.fontSize
  mirrorStyle.fontFamily = style.fontFamily
  mirrorStyle.lineHeight = style.lineHeight
  mirrorStyle.padding = style.padding
  mirrorStyle.border = style.border
  mirrorStyle.boxSizing = style.boxSizing
  const text = textarea.value.substring(0, caretIndex)
  mirror.textContent = text
  const span = document.createElement('span')
  span.textContent = '\u200b'
  mirror.appendChild(span)
  document.body.appendChild(mirror)
  const mirrorRect = mirror.getBoundingClientRect()
  const spanRect = span.getBoundingClientRect()
  document.body.removeChild(mirror)
  const taRect = textarea.getBoundingClientRect()
  const left = taRect.left + (spanRect.left - mirrorRect.left) - textarea.scrollLeft
  const top = taRect.top + (spanRect.top - mirrorRect.top) - textarea.scrollTop
  return { left, top }
}

function handleInput() {
  const el = inputRef.value?.$el
  const textarea = el?.querySelector?.('textarea') ?? el
  if (!textarea || typeof textarea.selectionStart !== 'number') return

  const pos = textarea.selectionStart
  const text = textarea.value || ''

  // 光标不在有效位置时直接关闭弹窗
  if (pos < 1 || pos > text.length) {
    showImageMentionPopup.value = false
    return
  }

  const char = text[pos - 1]
  if (char !== '@') {
    showImageMentionPopup.value = false
    return
  }

  if (imageListForPrompt.value.length === 0) {
    showImageMentionPopup.value = false
    return
  }

  triggerIndex.value = pos - 1
  const { left, top } = getCaretPixelPosition(textarea, pos)
  mentionPopupLeft.value = Math.round(left)
  mentionPopupTop.value = Math.round(top) + 20
  showImageMentionPopup.value = true
}

function onSelectImageMention(item) {
  const text = innerValue.value || ''
  const idx = triggerIndex.value
  const before = text.slice(0, idx)
  const after = text.slice(idx + 1)
  const nextText = `${before} ${item.name} ${after}`
  innerValue.value = nextText
  showImageMentionPopup.value = false

  nextTick(() => {
    const el = inputRef.value?.$el
    const textarea = el?.querySelector?.('textarea') ?? el
    if (textarea) {
      const newPos = idx + item.name.length + 2
      textarea.focus()
      textarea.setSelectionRange(newPos, newPos)
    }
  })
}

function onThumbError(e) {
  if (e?.target) {
    e.target.style.display = 'none'
  }
}

// 点击弹窗外部关闭
let clickOutsideCleanup = null
watch(showImageMentionPopup, (visible) => {
  if (clickOutsideCleanup) {
    clickOutsideCleanup()
    clickOutsideCleanup = null
  }
  if (!visible) return
  const onMouseDown = (e) => {
    const popup = imageMentionPopupRef.value
    const inputEl = inputRef.value?.$el
    const target = e.target
    if (popup?.contains(target) || inputEl?.contains(target)) return
    showImageMentionPopup.value = false
  }
  nextTick(() => {
    document.addEventListener('mousedown', onMouseDown)
    clickOutsideCleanup = () => document.removeEventListener('mousedown', onMouseDown)
  })
})

onUnmounted(() => {
  if (clickOutsideCleanup) clickOutsideCleanup()
})
</script>

<style lang="scss">
/* 参考图 @ 提及弹窗样式，全局使用 */
.image-mention-popup {
  position: fixed;
  z-index: 2000;
  min-width: 160px;
  max-height: 240px;
  overflow-y: auto;
  padding: 4px 0;
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  border: 1px solid #e4e7ed;
}

.image-mention-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 6px 12px;
  cursor: pointer;
  font-size: 14px;
  color: #303133;

  &:hover {
    background: #f5f7fa;
  }
}

.image-mention-thumb {
  width: 32px;
  height: 32px;
  object-fit: cover;
  border-radius: 4px;
  flex-shrink: 0;
}

.image-mention-name {
  flex: 1;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
</style>
