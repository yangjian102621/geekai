<template>
  <div v-if="message !== null" class="chat-detail">
    <div class="chat-line" v-html="md.render(processContent(message))"></div>
  </div>
  <div v-else class="chat-box chat-page p-2">
    <div v-for="item in messages" :key="item.id">
      <chat-prompt v-if="item.type === 'prompt'" :data="item" />
      <chat-reply v-else-if="item.type === 'reply'" :read-only="true" :data="item" />
    </div>
  </div>
</template>

<script setup>
import ChatPrompt from '@/components/ChatPrompt.vue'
import ChatReply from '@/components/ChatReply.vue'
import { processContent } from '@/utils/libs'
import hl from 'highlight.js'
import 'highlight.js/styles/a11y-dark.css'
import MarkdownIt from 'markdown-it'
import mathjaxPlugin from 'markdown-it-mathjax3'

defineProps({
  message: { type: String, default: null },
  messages: { type: Array, default: () => [] },
})

const md = MarkdownIt({
  breaks: true,
  html: true,
  linkify: true,
  typographer: true,
  highlight(str, lang) {
    const code = lang && hl.getLanguage(lang)
      ? hl.highlight(str, { language: lang }).value
      : md.utils.escapeHtml(str)
    return `<pre class="code-container"><code class="language-${lang} hljs">${code}</code></pre>`
  },
})
md.use(mathjaxPlugin)
</script>

<style lang="scss" scoped>
.chat-detail,
.chat-box {
  max-height: 90vh;
  overflow: auto;
}

.chat-box {
  --content-font-size: 16px;
  --content-color: #c1c1c1;
  font-family: 'Microsoft YaHei', '微软雅黑', Arial, sans-serif;

  .chat-line {
    ::-webkit-scrollbar {
      width: 0;
      height: 0;
      background-color: transparent;
    }

    font-size: 14px;
    display: flex;
    justify-content: flex-start;
  }
}
</style>
