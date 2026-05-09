<template>
  <div class="chat-line chat-line-prompt-chat">
    <div class="chat-line-inner">
      <div class="chat-icon">
        <img :src="data.icon" alt="User" />
      </div>

      <div class="chat-item">
        <div v-if="files && files.length > 0" class="file-list-box">
          <div v-for="file in files" :key="file.url">
            <div class="image" v-if="isImage(file.ext)">
              <el-image :src="file.url" fit="cover" />
            </div>
            <div class="item" v-else>
              <div class="icon">
                <el-image :src="GetFileIcon(file.ext)" fit="cover" />
              </div>
              <div class="body">
                <div class="title">
                  <el-link :href="file.url" target="_blank" style="--el-font-weight-primary: bold"
                    >{{ file.name }}
                  </el-link>
                </div>
                <div class="info">
                  <span>{{ GetFileType(file.ext) }}</span>
                  <span>{{ FormatFileSize(file.size) }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 编辑模式 -->
        <div v-if="isEditing" class="edit-mode">
          <div class="flex flex-row space-x-2 w-full">
            <div>
              <el-tooltip class="box-item" effect="dark" content="取消" placement="top">
                <i
                  class="iconfont icon-error-line !text-lg mr-1 cursor-pointer"
                  @click="cancelEdit"
                ></i>
              </el-tooltip>
            </div>
            <div class="w-full">
              <textarea
                v-model="editText"
                :rows="3"
                placeholder="请输入修改后的内容..."
                class="w-full p-3 border-2 border-purple-500 rounded-md text-sm"
                resize="vertical"
                style="background-color: var(--chat-content-bg); color: var(--theme-text-primary)"
              ></textarea>
            </div>
            <div>
              <el-tooltip class="box-item" effect="dark" content="提交" placement="top">
                <i
                  class="iconfont icon-back-circle cursor-pointer !text-3xl text-purple-500 mr-1 hover:text-purple-700 mb-2"
                  style="transform: rotate(90deg); display: inline-block"
                  @click="submitEdit"
                ></i>
              </el-tooltip>
            </div>
          </div>
        </div>

        <!-- 显示模式 -->
        <div v-else>
          <div class="content-wrapper">
            <div class="content position-relative">
              <div v-html="content"></div>
            </div>
          </div>
          <div
            class="flex text-gray-500 text-sm py-2 justify-end items-center space-x-2"
            v-if="data.created_at > 0"
          >
            <span class="flex items-center"
              ><i class="iconfont icon-clock mr-1"></i> {{ dateFormat(data.created_at) }}</span
            >
            <span class="flex items-center">
              <el-tooltip class="box-item" effect="dark" content="复制" placement="top">
                <i
                  class="iconfont icon-copy cursor-pointer !text-sm"
                  @click="copyContent(data.content.text)"
                ></i>
              </el-tooltip>
            </span>
            <span class="flex items-center">
              <el-tooltip class="box-item" effect="dark" content="修改提问" placement="top">
                <i class="iconfont icon-edit cursor-pointer !text-sm" @click="startEdit"></i>
              </el-tooltip>
            </span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { FormatFileSize, GetFileIcon, GetFileType } from '@/store/system'
import { showMessageSuccess } from '@/utils/dialog'
import { dateFormat, isImage, processPrompt } from '@/utils/libs'
import { ElMessage } from 'element-plus'
import hl from 'highlight.js'
import MarkdownIt from 'markdown-it'
import emoji from 'markdown-it-emoji'
import mathjaxPlugin from 'markdown-it-mathjax3'
import { onMounted, ref } from 'vue'

const md = new MarkdownIt({
  breaks: true,
  html: true,
  linkify: true,
  typographer: true,
  highlight: function (str, lang) {
    const codeIndex = parseInt(Date.now()) + Math.floor(Math.random() * 10000000)
    // 显示复制代码按钮
    const copyBtn = `<span class="copy-code-btn" data-clipboard-action="copy" data-clipboard-target="#copy-target-${codeIndex}">复制</span>
<textarea style="position: absolute;top: -9999px;left: -9999px;z-index: -9999;" id="copy-target-${codeIndex}">${str.replace(
      /<\/textarea>/g,
      '&lt;/textarea>'
    )}</textarea>`
    if (lang && hl.getLanguage(lang)) {
      const langHtml = `<span class="lang-name">${lang}</span>`
      // 处理代码高亮
      const preCode = hl.highlight(str, { language: lang, ignoreIllegals: true }).value
      // 将代码包裹在 pre 中
      return `<pre class="code-container"><code class="language-${lang} hljs">${preCode}</code>${copyBtn} ${langHtml}</pre>`
    }

    // 处理代码高亮
    const preCode = md.utils.escapeHtml(str)
    // 将代码包裹在 pre 中
    return `<pre class="code-container"><code class="language-${lang} hljs">${preCode}</code>${copyBtn}</pre>`
  },
})
md.use(mathjaxPlugin)
md.use(emoji)

const props = defineProps({
  data: {
    type: Object,
    default: {
      content: {
        text: '',
        files: [],
      },
      created_at: '',
      tokens: 0,
      model: '',
      icon: '',
    },
  },
  messageIndex: {
    type: Number,
    default: -1,
  },
})
const finalTokens = ref(props.data.tokens)
const content = ref(processPrompt(props.data.content.text))
const files = ref(props.data.content.files)

// 编辑相关状态
const isEditing = ref(false)
const editText = ref('')

// 定义emit事件
const emit = defineEmits(['edit'])

onMounted(() => {
  processFiles()
})

const processFiles = () => {
  if (!props.data.content) {
    return
  }
  content.value = md.render(content.value.trim())
}
const isExternalImg = (link, files) => {
  return isImage(link) && !files.find((file) => file.url === link)
}
const copyContent = (text) => {
  navigator.clipboard.writeText(text)
  showMessageSuccess('复制成功')
}

// 开始编辑
const startEdit = () => {
  isEditing.value = true
  editText.value = props.data.content.text
}

// 取消编辑
const cancelEdit = () => {
  isEditing.value = false
  editText.value = ''
}

// 提交编辑
const submitEdit = () => {
  if (!editText.value.trim()) {
    ElMessage.warning('内容不能为空')
    return
  }
  // 发送重新提交事件，传递修改后的内容
  emit('edit', {
    messageIndex: props.messageIndex,
    newContent: editText.value.trim(),
  })

  // 退出编辑模式
  isEditing.value = false
  editText.value = ''
}
</script>

<style lang="scss">
@use '@/assets/css/markdown/vue.css' as *;
.chat-page {
  .chat-line-prompt-chat {
    background: var(--chat-bg);
    justify-content: center;
    width: 100%;
    padding-bottom: 1.5rem;
    padding-top: 1.5rem;

    .chat-line-inner {
      display: flex;
      width: 100%;
      padding: 0 25px;
      flex-flow: row-reverse;

      .chat-icon {
        margin-left: 20px;

        img {
          width: 36px;
          height: 36px;
          border-radius: 50%;
          padding: 1px;
        }
      }

      .chat-item {
        padding: 0;
        overflow: hidden;
        max-width: calc(100% - 110px);
        width: 100%;

        .file-list-box {
          display: flex;
          flex-flow: column;

          .image {
            display: flex;
            flex-flow: row;
            margin-right: 10px;
            position: relative;
            justify-content: end;

            .el-image {
              border: 1px solid #e3e3e3;
              border-radius: 10px;
              margin-bottom: 10px;
              max-width: 150px;
              max-height: 150px;
            }
          }

          .item {
            display: flex;
            flex-flow: row;
            border-radius: 10px;
            background-color: var(--chat-content-bg);
            color: var(--theme-text-color-primary);
            border: 1px solid #e3e3e3;
            padding: 6px;
            margin-bottom: 10px;

            .icon {
              .el-image {
                width: 40px;
                height: 40px;
              }
            }

            .body {
              margin-left: 8px;
              font-size: 14px;

              .title {
                font-weight: bold;
                line-height: 24px;
                color: #0d0d0d;
              }

              .info {
                color: #b4b4b4;

                span {
                  margin-right: 10px;
                }
              }
            }
          }
        }

        .content-wrapper {
          display: flex;
          flex-flow: row-reverse;

          .content {
            word-break: break-word;
            padding: 1rem;
            color: var(--theme-text-primary);
            font-size: var(--content-font-size);
            overflow: auto;
            background-color: var(--chat-user-content-bg);
            border-radius: 10px 0 10px 10px;

            img {
              max-width: 600px;
              border-radius: 10px;
              margin: 10px 0;
            }

            p {
              line-height: 1.5;
            }

            p:last-child {
              margin-bottom: 0;
            }

            p:first-child {
              margin-top: 0;
            }
          }
        }

        .edit-mode {
          width: 100%;
          margin-top: 15px;
        }
      }
    }
  }
}
</style>
