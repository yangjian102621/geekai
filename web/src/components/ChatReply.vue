<template>
  <div class="chat-reply">
    <div class="chat-line chat-line-reply-chat">
      <div class="chat-line-inner">
        <div class="chat-icon">
          <img :src="data.icon" alt="ChatGPT" />
        </div>
        <div class="chat-item">
          <div class="content-wrapper">
            <div
              class="content"
              v-html="md.render(processContent(data.content.text))"
              v-if="data.content.text"
            ></div>
            <div class="content flex justify-start items-center" v-else>
              <span class="mr-2">AI 思考中</span> <Thinking :duration="1.5" />
            </div>
          </div>
          <div
            class="flex text-gray-500 text-sm py-2 items-center space-x-2"
            v-if="data.created_at"
          >
            <span class="flex items-center"
              ><i class="iconfont icon-clock mr-1"></i> {{ dateFormat(data.created_at) }}</span
            >
            <span class="flex items-center">
              <el-tooltip class="box-item" effect="dark" content="复制回答" placement="top">
                <el-icon class="copy-reply" :data-clipboard-text="data.content.text">
                  <DocumentCopy />
                </el-icon>
              </el-tooltip>
            </span>
            <span class="flex items-center" @click="reGenerate">
              <el-tooltip class="box-item" effect="dark" content="重新生成" placement="bottom">
                <i class="iconfont icon-refresh cursor-pointer !text-sm"></i>
              </el-tooltip>
            </span>

            <span class="flex items-center">
              <el-tooltip
                class="box-item"
                effect="dark"
                content="生成语音朗读"
                placement="bottom"
                v-if="!isPlaying"
              >
                <i
                  class="iconfont icon-speaker !text-sm cursor-pointer"
                  @click="synthesis(data.content.text)"
                ></i>
              </el-tooltip>
              <el-tooltip
                class="box-item"
                effect="dark"
                content="暂停播放"
                placement="bottom"
                v-else
              >
                <el-image class="voice-icon" :src="playIcon" @click="stopSynthesis()" />
              </el-tooltip>
            </span>
          </div>
        </div>
      </div>
    </div>
    <audio ref="audio" @ended="isPlaying = false" />
  </div>
</template>

<script setup>
import { useSharedStore } from '@/store/sharedata'
import { httpPost } from '@/utils/http'
import { dateFormat, processContent } from '@/utils/libs'
import { DocumentCopy } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import hl from 'highlight.js'
import MarkdownIt from 'markdown-it'
import emoji from 'markdown-it-emoji'
import mathjaxPlugin from 'markdown-it-mathjax3'
import { nextTick, onMounted, reactive, ref, watchEffect } from 'vue'
import Thinking from './Thinking.vue'
// eslint-disable-next-line no-undef,no-unused-vars
const props = defineProps({
  data: {
    type: Object,
    default: {
      type: 'text',
      icon: '',
      content: {
        text: '',
        files: [],
      },
      created_at: '',
      tokens: 0,
    },
  },
  messageIndex: {
    type: Number,
    default: -1,
  },
})

const audio = ref(null)
const isPlaying = ref(false)
const playIcon = ref('/images/voice.gif')
const store = useSharedStore()

// 添加代码块展开/收起状态管理
const codeBlockStates = reactive({})

const md = new MarkdownIt({
  breaks: true,
  html: true,
  linkify: true,
  typographer: true,
  highlight: function (str, lang) {
    const codeIndex = parseInt(Date.now()) + Math.floor(Math.random() * 10000000)
    // 显示复制代码按钮和展开/收起按钮
    const copyBtn = `<div class="flex">
      <span class="text-[12px] mr-2 text-[#00e0e0] cursor-pointer expand-btn" data-code-id="${codeIndex}" onclick="window.toggleCodeBlock('${codeIndex}')">收起</span>
      <span class="copy-code-btn" data-clipboard-action="copy" data-clipboard-target="#copy-target-${codeIndex}">复制</span>
      </div><textarea style="position: absolute;top: -9999px;left: -9999px;z-index: -9999;" id="copy-target-${codeIndex}">${str.replace(
      /<\/textarea>/g,
      '&lt;/textarea>'
    )}</textarea>`
    let langHtml = ''
    let preCode = ''
    // 处理代码高亮
    if (lang && hl.getLanguage(lang)) {
      langHtml = `<span class="lang-name">${lang}</span>`
      preCode = hl.highlight(str, { language: lang }).value
    } else {
      preCode = md.utils.escapeHtml(str)
    }

    // 将代码包裹在 pre 中，添加展开状态的类（默认展开）
    return `<pre class="code-container flex flex-col code-expanded" data-code-id="${codeIndex}">
      <div class="flex justify-between bg-[#50505a] w-full rounded-tl-[10px] rounded-tr-[10px] px-3 py-1">${langHtml}${copyBtn}</div>
      <code class="language-${lang} hljs">${preCode}</code> 
      <span class="copy-code-btn absolute right-3 bottom-3" data-clipboard-action="copy" data-clipboard-target="#copy-target-${codeIndex}">复制</span></pre>`
  },
})
md.use(mathjaxPlugin)
md.use(emoji)
const emits = defineEmits(['regen'])

if (!props.data.icon) {
  props.data.icon = 'images/gpt-icon.png'
}

const synthesis = (text) => {
  isPlaying.value = true
  httpPost('/api/chat/tts', { text: text, model_id: store.ttsModel }, { responseType: 'blob' })
    .then((response) => {
      // 创建 Blob 对象，明确指定 MIME 类型
      const blob = new Blob([response], { type: 'audio/mpeg' }) // 假设音频格式为 MP3
      const audioUrl = URL.createObjectURL(blob)
      // 播放音频
      audio.value.src = audioUrl
      audio.value
        .play()
        .then(() => {
          // 播放完成后释放 URL
          URL.revokeObjectURL(audioUrl)
        })
        .catch(() => {
          ElMessage.error('音频播放失败，请检查浏览器是否支持该音频格式')
          isPlaying.value = false
        })
    })
    .catch((e) => {
      ElMessage.error('语音合成失败：' + e.message)
      isPlaying.value = false
    })
}

const stopSynthesis = () => {
  isPlaying.value = false
  audio.value.pause()
  audio.value.currentTime = 0
}

// 重新生成
const reGenerate = () => {
  emits('regen', props.messageIndex)
}

// 添加代码块展开/收起功能
const toggleCodeBlock = (codeId) => {
  const codeContainer = document.querySelector(`pre[data-code-id="${codeId}"]`)
  const expandBtn = document.querySelector(`.expand-btn[data-code-id="${codeId}"]`)

  if (codeContainer && expandBtn) {
    if (codeContainer.classList.contains('code-collapsed')) {
      codeContainer.classList.remove('code-collapsed')
      codeContainer.classList.add('code-expanded')
      expandBtn.textContent = '收起'
    } else {
      codeContainer.classList.remove('code-expanded')
      codeContainer.classList.add('code-collapsed')
      expandBtn.textContent = '展开'
    }
  }
}

// 将函数暴露到全局作用域
window.toggleCodeBlock = toggleCodeBlock

// 添加事件监听
onMounted(() => {
  nextTick(() => {
    setupCodeBlockEvents()
  })
})

// 监听内容变化，重新绑定事件
watchEffect(() => {
  if (props.data.content.text) {
    nextTick(() => {
      // 延迟一点时间确保DOM完全渲染
      setTimeout(() => {
        setupCodeBlockEvents()
      }, 100)
    })
  }
})

const setupCodeBlockEvents = () => {
  // 检查所有代码块并设置展开按钮的显示状态
  const expandBtns = document.querySelectorAll('.expand-btn')

  expandBtns.forEach((btn) => {
    const codeId = btn.getAttribute('data-code-id')
    const codeContainer = document.querySelector(`pre[data-code-id="${codeId}"]`)
    const codeElement = codeContainer?.querySelector('.hljs')

    if (codeElement) {
      // 临时移除高度限制来获取真实高度
      const originalMaxHeight = codeElement.style.maxHeight
      codeElement.style.maxHeight = 'none'
      const realHeight = codeElement.scrollHeight
      codeElement.style.maxHeight = originalMaxHeight

      // 如果代码块高度小于等于200px，隐藏展开按钮
      if (realHeight <= 200) {
        btn.style.display = 'none'
        // 移除收起状态的类，让短代码块完全展示
        codeContainer.classList.remove('code-collapsed')
        codeContainer.classList.add('code-expanded')
      } else {
        btn.style.display = 'inline'
        // 确保长代码块默认展开
        if (
          !codeContainer.classList.contains('code-expanded') &&
          !codeContainer.classList.contains('code-collapsed')
        ) {
          codeContainer.classList.add('code-expanded')
        }
      }
    }
  })
}
</script>

<style lang="scss">
@use '@/assets/css/markdown/vue.css' as *;

.chat-page,
.chat-export {
  --font-family: Menlo, '微软雅黑', 'Roboto Mono', 'Courier New', Courier, monospace, 'Inter',
    sans-serif;
  font-family: var(--font-family);

  .chat-line-reply-chat {
    justify-content: center;
    padding: 1.5rem;

    .chat-line-inner {
      display: flex;
      width: 100%;
      flex-flow: row;

      .chat-icon {
        margin-right: 20px;

        img {
          width: 36px;
          height: 36px;
          border-radius: 50%;
          padding: 1px;
        }
      }

      .chat-item {
        position: relative;
        padding: 0;
        overflow: hidden;
        width: 100%;
        max-width: calc(100% - 110px);

        .content-wrapper {
          display: flex;
          .content {
            min-height: 20px;
            word-break: break-word;
            padding: 1rem;
            color: var(--theme-text-primary);

            font-size: var(--content-font-size);
            overflow: auto;
            // background-color #F5F5F5
            background-color: var(--chat-content-bg);
            border-radius: 0 10px 10px 10px;
            width: 100%;

            p:first-child {
              margin-top: 0;
            }

            p:last-child {
              margin-bottom: 0;
            }
          }
        }
      }

      .tool-box {
        font-size: 16px;

        .el-button {
          height: 20px;
          padding: 5px 2px;
        }
      }
    }
  }
}
</style>
