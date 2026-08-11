<template>
  <div class="chat-export-page" v-loading="loading">
    <div class="export-container">
      <!-- 页面头部 -->
      <header class="export-header">
        <div class="header-content">
          <h1 class="export-title">{{ chatTitle || '聊天记录导出' }}</h1>
          <div class="header-meta" v-if="chatData.length > 0">
            <span class="meta-item">
              <i class="iconfont icon-message"></i>
              {{ chatData.length }} 条消息
            </span>
            <span class="meta-item" v-if="exportDate">
              <i class="iconfont icon-clock"></i>
              {{ exportDate }}
            </span>
          </div>
        </div>
      </header>

      <!-- 聊天消息列表 -->
      <main class="export-content" v-if="chatData.length > 0">
        <div class="messages-list" id="messages-list">
          <div
            v-for="(message, index) in chatData"
            :key="message.id || index"
            class="message-item"
            :class="{
              'message-user': message.type === 'prompt',
              'message-assistant': message.type === 'reply' || message.type === 'mj',
            }"
          >
            <!-- 用户消息 -->
            <div v-if="message.type === 'prompt'" class="message-bubble message-bubble-user">
              <div class="message-header">
                <div class="message-avatar">
                  <img
                    :src="message.icon || '/images/user-icon.png'"
                    alt="用户"
                    class="avatar-img"
                  />
                </div>
                <div class="message-info">
                  <span class="message-role">用户</span>
                  <span class="message-time" v-if="message.created_at">
                    {{ formatTime(message.created_at) }}
                  </span>
                </div>
              </div>
              <div class="message-body">
                <!-- 文件列表 -->
                <div v-if="message.content?.files && message.content.files.length > 0" class="message-files">
                  <div
                    v-for="file in message.content.files"
                    :key="file.url"
                    class="file-item"
                  >
                    <img
                      v-if="isImageFile(file.ext)"
                      :src="file.url"
                      :alt="file.name"
                      class="file-image"
                    />
                    <div v-else class="file-card">
                      <i class="iconfont icon-file"></i>
                      <div class="file-info">
                        <div class="file-name">{{ file.name }}</div>
                        <div class="file-meta">{{ formatFileSize(file.size) }}</div>
                      </div>
                    </div>
                  </div>
                </div>
                <!-- 文本内容 -->
                <div
                  class="message-text"
                  v-html="renderMarkdown(message.content?.text || message.content || '')"
                ></div>
                <!-- 操作按钮 -->
                <div class="message-actions">
                  <button
                    class="action-btn"
                    @click="copyMessage(message.content?.text || message.content || '')"
                    title="复制"
                  >
                    <i class="iconfont icon-copy"></i>
                  </button>
                </div>
              </div>
            </div>

            <!-- AI 回复消息 -->
            <div
              v-else-if="message.type === 'reply' || message.type === 'mj'"
              class="message-bubble message-bubble-assistant"
            >
              <div class="message-header">
                <div class="message-avatar">
                  <img
                    :src="message.icon || '/images/gpt-icon.png'"
                    alt="AI助手"
                    class="avatar-img"
                  />
                </div>
                <div class="message-info">
                  <span class="message-role">AI 助手</span>
                  <span class="message-time" v-if="message.created_at">
                    {{ formatTime(message.created_at) }}
                  </span>
                </div>
              </div>
              <div class="message-body">
                <div
                  class="message-text"
                  v-html="renderMarkdown(getMessageContent(message))"
                ></div>

                <!-- 文件列表 -->
                <AttachmentList :files="message.content?.files || []" />
                <!-- 操作按钮 -->
                <div class="message-actions">
                  <button
                    class="action-btn"
                    @click="copyMessage(getMessageContent(message))"
                    title="复制"
                  >
                    <i class="iconfont icon-copy"></i>
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </main>

      <!-- 错误状态 -->
      <div class="error-state" v-if="error && !loading">
        <div class="error-content">
          <i class="iconfont icon-warning error-icon"></i>
          <h3 class="error-title">加载失败</h3>
          <p class="error-text">{{ error }}</p>
          <button class="retry-btn" @click="init">重试</button>
        </div>
      </div>

      <!-- 空状态 -->
      <div class="empty-state" v-else-if="!loading && !error && chatData.length === 0">
        <div class="empty-content">
          <i class="iconfont icon-chat empty-icon"></i>
          <h3 class="empty-title">暂无聊天记录</h3>
          <p class="empty-text">该聊天会话中没有消息记录</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { httpGet } from '@/utils/http'
import { ElMessage } from 'element-plus'
import hl from 'highlight.js'
import 'highlight.js/styles/a11y-dark.css'
import MarkdownIt from 'markdown-it'
import emoji from 'markdown-it-emoji'
import mathjaxPlugin from 'markdown-it-mathjax3'
import { nextTick, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { processContent, isImage } from '@/utils/libs'
import AttachmentList from '@/components/AttachmentList.vue'

const router = useRouter()
const chatId = router.currentRoute.value.query['chat_id']
const loading = ref(true)
const chatTitle = ref('')
const chatData = ref([])
const error = ref(null)
const exportDate = ref('')

// 初始化 Markdown 渲染器
const md = new MarkdownIt({
  breaks: true,
  html: true,
  linkify: true,
  typographer: true,
  highlight: function (str, lang) {
    const codeIndex = Date.now() + Math.floor(Math.random() * 10000000)
    const copyBtn = `<button class="code-copy-btn" onclick="navigator.clipboard.writeText(\`${str.replace(/`/g, '\\`').replace(/\$/g, '\\$')}\`); this.textContent='已复制'; setTimeout(()=>this.textContent='复制', 2000)">复制</button>`
    
    if (lang && hl.getLanguage(lang)) {
      try {
        const highlighted = hl.highlight(str, { language: lang, ignoreIllegals: true }).value
        return `<div class="code-block-wrapper">
          <div class="code-header">
            <span class="code-lang">${lang}</span>
            ${copyBtn}
          </div>
          <pre class="code-block"><code class="language-${lang} hljs">${highlighted}</code></pre>
        </div>`
      } catch (e) {
        console.warn('代码高亮失败:', e)
      }
    }
    
    const escaped = md.utils.escapeHtml(str)
    return `<div class="code-block-wrapper">
      <div class="code-header">
        ${copyBtn}
      </div>
      <pre class="code-block"><code>${escaped}</code></pre>
    </div>`
  },
})
md.use(mathjaxPlugin)
md.use(emoji)

// 渲染 Markdown
const renderMarkdown = (text) => {
  if (!text) return ''
  // 使用 processContent 处理内容（处理图片链接、think 标签等）
  const processed = processContent(String(text))
  return md.render(processed)
}

// 获取消息内容
const getMessageContent = (message) => {
  if (message.type === 'mj') {
    try {
      const content = typeof message.content === 'string' 
        ? JSON.parse(message.content) 
        : message.content
      return content?.content || content?.text || ''
    } catch (e) {
      return message.content?.text || message.content || ''
    }
  }
  return message.content?.text || message.content || ''
}

// 格式化时间
const formatTime = (timestamp) => {
  if (!timestamp) return ''
  const date = new Date(timestamp * 1000)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
  })
}

// 格式化文件大小
const formatFileSize = (bytes) => {
  if (!bytes) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i]
}

// 判断是否为图片文件
const isImageFile = (ext) => {
  if (!ext) return false
  // 使用项目中的 isImage 函数
  return isImage('.' + ext)
}

// 复制消息
const copyMessage = async (text) => {
  try {
    await navigator.clipboard.writeText(text)
    ElMessage.success('复制成功！')
  } catch (e) {
    ElMessage.error('复制失败：' + e.message)
  }
}

// 高亮代码块
const highlightCodeBlocks = () => {
  nextTick(() => {
    try {
      const messagesList = document.querySelector('#messages-list')
      if (!messagesList) return

      hl.configure({ ignoreUnescapedHTML: true })
      const blocks = messagesList.querySelectorAll('pre code')
      blocks.forEach((block) => {
        if (!block.classList.contains('hljs')) {
          try {
            hl.highlightElement(block)
          } catch (e) {
            console.warn('代码高亮失败:', e)
          }
        }
      })
    } catch (e) {
      console.warn('初始化代码高亮失败:', e)
    }
  })
}

// 加载聊天历史
const loadChatHistory = async () => {
  if (!chatId) {
    error.value = '缺少聊天ID参数'
    loading.value = false
    return
  }

  try {
    const res = await httpGet('/api/chat/history?chat_id=' + chatId)
    const data = res.data

    if (!data || !Array.isArray(data)) {
      loading.value = false
      return
    }

    chatData.value = data.map((item) => {
      // 处理 MJ 类型消息
      if (item.type === 'mj' && typeof item.content === 'string') {
        try {
          item.content = JSON.parse(item.content)
        } catch (e) {
          console.warn('解析 MJ 内容失败:', e)
        }
      }
      return item
    })

    highlightCodeBlocks()
    loading.value = false
  } catch (e) {
    console.error('加载聊天记录失败:', e)
    error.value = '加载聊天记录失败：' + (e.message || '未知错误')
    ElMessage.error(error.value)
    loading.value = false
  }
}

// 加载聊天详情
const loadChatDetail = async () => {
  if (!chatId) return

  try {
    const res = await httpGet('/api/chat/detail?chat_id=' + chatId)
    if (res.data) {
      chatTitle.value = res.data.title || '聊天记录导出'
    }
  } catch (e) {
    console.warn('加载聊天详情失败:', e)
  }
}

// 初始化
const init = async () => {
  loading.value = true
  error.value = null
  exportDate.value = new Date().toLocaleString('zh-CN')

  await Promise.all([loadChatHistory(), loadChatDetail()])
}

onMounted(() => {
  init()
})
</script>

<style lang="scss" scoped>
.chat-export-page {
  min-height: 100vh;
  background: var(--theme-bg-color);
  padding: 0;
  display: flex;
  justify-content: center;
  align-items: flex-start;

  .export-container {
    width: 100%;
    max-width: 900px;
    margin: 0 auto;
    background: var(--chat-bg, #ffffff);
    box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
    min-height: 100vh;

    // 页面头部
    .export-header {
      background: var(--theme-bg-color);
      border-bottom: 1px solid var(--el-border-color);
      padding: 2rem 2rem 1.5rem;
      position: sticky;
      top: 0;
      z-index: 10;
      backdrop-filter: blur(10px);
      background: var(--chat-bg, #ffffff);

      .header-content {
        .export-title {
          font-size: 1.75rem;
          font-weight: 600;
          color: var(--theme-text-color-primary);
          margin: 0 0 0.75rem 0;
          line-height: 1.3;
        }

        .header-meta {
          display: flex;
          gap: 1.5rem;
          flex-wrap: wrap;

          .meta-item {
            display: flex;
            align-items: center;
            gap: 0.5rem;
            font-size: 0.875rem;
            color: var(--theme-text-color-secondary);

            i {
              font-size: 1rem;
            }
          }
        }
      }
    }

    // 内容区域
    .export-content {
      padding: 2rem;

      .messages-list {
        display: flex;
        flex-direction: column;
        gap: 2rem;

        .message-item {
          display: flex;
          flex-direction: column;

          &.message-user {
            align-items: flex-end;
          }

          &.message-assistant {
            align-items: flex-start;
          }

          .message-bubble {
            max-width: 100%;
            min-width: 100%;
            border-radius: 12px;
            overflow: hidden;
            transition: transform 0.2s ease, box-shadow 0.2s ease;

            &:hover {
              transform: translateY(-1px);
              box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
            }

            .message-header {
              display: flex;
              align-items: center;
              gap: 0.75rem;
              padding: 0.75rem 1rem;
              background: var(--theme-bg-color);
              border-bottom: 1px solid var(--el-border-color);

              .message-avatar {
                width: 32px;
                height: 32px;
                border-radius: 50%;
                overflow: hidden;
                flex-shrink: 0;
                background: var(--el-border-color);

                .avatar-img {
                  width: 100%;
                  height: 100%;
                  object-fit: cover;
                }
              }

              .message-info {
                display: flex;
                align-items: center;
                gap: 0.75rem;
                flex: 1;

                .message-role {
                  font-weight: 600;
                  font-size: 0.875rem;
                  color: var(--theme-text-color-primary);
                }

                .message-time {
                  font-size: 0.75rem;
                  color: var(--theme-text-color-secondary);
                }
              }
            }

            .message-body {
              padding: 1rem 1.25rem;
              background: var(--chat-content-bg, #f5f7fc);
              position: relative;

              .message-files {
                display: flex;
                flex-direction: column;
                gap: 0.75rem;
                margin-bottom: 1rem;

                .file-item {
                  .file-image {
                    max-width: 100%;
                    max-height: 300px;
                    border-radius: 8px;
                    object-fit: contain;
                  }

                  .file-card {
                    display: flex;
                    align-items: center;
                    gap: 0.75rem;
                    padding: 0.75rem;
                    background: var(--theme-bg-color);
                    border: 1px solid var(--el-border-color);
                    border-radius: 8px;

                    i {
                      font-size: 1.5rem;
                      color: var(--el-color-primary);
                    }

                    .file-info {
                      flex: 1;
                      min-width: 0;

                      .file-name {
                        font-size: 0.875rem;
                        font-weight: 500;
                        color: var(--theme-text-color-primary);
                        word-break: break-all;
                      }

                      .file-meta {
                        font-size: 0.75rem;
                        color: var(--theme-text-color-secondary);
                        margin-top: 0.25rem;
                      }
                    }
                  }
                }
              }

              .message-text {
                font-size: 15px;
                line-height: 1.7;
                color: var(--theme-text-color-primary);
                word-break: break-word;

                :deep(p) {
                  margin: 0.75rem 0;
                  &:first-child {
                    margin-top: 0;
                  }
                  &:last-child {
                    margin-bottom: 0;
                  }
                }

                :deep(ul),
                :deep(ol) {
                  margin: 0.75rem 0;
                  padding-left: 1.5rem;
                }

                :deep(li) {
                  margin: 0.5rem 0;
                }

                :deep(blockquote) {
                  margin: 1rem 0;
                  padding: 0.75rem 1rem;
                  background: var(--quote-bg-color, #f0f0f0);
                  border-left: 3px solid var(--el-color-primary);
                  border-radius: 4px;
                }

                :deep(table) {
                  width: 100%;
                  border-collapse: collapse;
                  margin: 1rem 0;
                  font-size: 0.875rem;

                  th,
                  td {
                    padding: 0.5rem;
                    border: 1px solid var(--el-border-color);
                  }

                  th {
                    background: var(--theme-bg-color);
                    font-weight: 600;
                  }
                }

                :deep(img) {
                  max-width: 100%;
                  height: auto;
                  border-radius: 8px;
                  margin: 1rem 0;
                }

                :deep(a) {
                  color: var(--el-color-primary);
                  text-decoration: none;

                  &:hover {
                    text-decoration: underline;
                  }
                }

                :deep(.code-block-wrapper) {
                  margin: 1rem 0;
                  border-radius: 8px;
                  overflow: hidden;
                  background: #2b2b2b;
                  border: 1px solid var(--el-border-color);

                  .code-header {
                    display: flex;
                    justify-content: space-between;
                    align-items: center;
                    padding: 0.5rem 0.75rem;
                    background: #1e1e1e;
                    border-bottom: 1px solid #3a3a3a;

                    .code-lang {
                      font-size: 0.75rem;
                      color: #a0a0a0;
                      text-transform: uppercase;
                    }

                    .code-copy-btn {
                      background: transparent;
                      border: 1px solid #555;
                      color: #ccc;
                      padding: 0.25rem 0.5rem;
                      border-radius: 4px;
                      font-size: 0.75rem;
                      cursor: pointer;
                      transition: all 0.2s ease;

                      &:hover {
                        background: #3a3a3a;
                        border-color: #777;
                      }
                    }
                  }

                  .code-block {
                    margin: 0;
                    padding: 1rem;
                    overflow-x: auto;
                    background: #2b2b2b;

                    code {
                      font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
                      font-size: 0.875rem;
                      line-height: 1.6;
                      color: #d4d4d4;
                    }
                  }
                }
              }
            }

            .message-actions {
              display: flex;
              gap: 0.5rem;
              margin-top: 0.75rem;
              padding-top: 0.75rem;
              border-top: 1px solid var(--el-border-color);

              .action-btn {
                display: flex;
                align-items: center;
                justify-content: center;
                width: 32px;
                height: 32px;
                border: none;
                background: transparent;
                color: var(--theme-text-color-secondary);
                border-radius: 6px;
                cursor: pointer;
                transition: all 0.2s ease;

                &:hover {
                  background: var(--theme-bg-color);
                  color: var(--el-color-primary);
                }

                i {
                  font-size: 0.875rem;
                }
              }
            }
          }

          &.message-user .message-bubble {
            background: var(--chat-user-content-bg, #e0dfff);
            border: 1px solid rgba(107, 80, 225, 0.2);

            .message-body {
              background: var(--chat-user-content-bg, #e0dfff);
            }
          }

          &.message-assistant .message-bubble {
            background: var(--chat-content-bg, #f5f7fc);
            border: 1px solid var(--el-border-color);

            .message-body {
              background: var(--chat-content-bg, #f5f7fc);
            }
          }
        }
      }
    }

    // 错误状态
    .error-state {
      display: flex;
      justify-content: center;
      align-items: center;
      min-height: 400px;
      padding: 3rem 2rem;

      .error-content {
        text-align: center;
        max-width: 500px;

        .error-icon {
          font-size: 4rem;
          color: var(--el-color-error);
          opacity: 0.8;
          margin-bottom: 1rem;
          display: block;
        }

        .error-title {
          font-size: 1.25rem;
          font-weight: 600;
          color: var(--theme-text-color-primary);
          margin: 0 0 0.5rem 0;
        }

        .error-text {
          font-size: 0.9375rem;
          color: var(--el-color-error);
          margin: 0 0 1.5rem 0;
          line-height: 1.6;
        }

        .retry-btn {
          padding: 0.625rem 1.5rem;
          background: var(--el-color-primary);
          color: #fff;
          border: none;
          border-radius: 6px;
          font-size: 0.9375rem;
          cursor: pointer;
          transition: all 0.2s ease;

          &:hover {
            opacity: 0.9;
            transform: translateY(-1px);
          }
        }
      }
    }

    // 空状态
    .empty-state {
      display: flex;
      justify-content: center;
      align-items: center;
      min-height: 400px;
      padding: 3rem 2rem;

      .empty-content {
        text-align: center;

        .empty-icon {
          font-size: 4rem;
          color: var(--theme-text-color-secondary);
          opacity: 0.5;
          margin-bottom: 1rem;
          display: block;
        }

        .empty-title {
          font-size: 1.25rem;
          font-weight: 600;
          color: var(--theme-text-color-primary);
          margin: 0 0 0.5rem 0;
        }

        .empty-text {
          font-size: 0.9375rem;
          color: var(--theme-text-color-secondary);
          margin: 0;
          line-height: 1.6;
        }
      }
    }
  }
}

// 响应式设计
@media (max-width: 768px) {
  .chat-export-page {
    .export-container {
      .export-header {
        padding: 1.5rem 1rem 1rem;

        .header-content {
          .export-title {
            font-size: 1.5rem;
          }

          .header-meta {
            gap: 1rem;
            font-size: 0.8125rem;
          }
        }
      }

      .export-content {
        padding: 1.5rem 1rem;

        .messages-list {
          gap: 1.5rem;

          .message-item {
            .message-bubble {
              //max-width: 95%;

              .message-header {
                padding: 0.625rem 0.875rem;

                .message-avatar {
                  width: 28px;
                  height: 28px;
                }
              }

              .message-body {
                padding: 0.875rem 1rem;
                font-size: 14px;
              }
            }
          }
        }
      }
    }
  }
}

@media (max-width: 480px) {
  .chat-export-page {
    .export-container {
      .export-header {
        padding: 1rem 0.75rem 0.75rem;

        .header-content {
          .export-title {
            font-size: 1.25rem;
          }
        }
      }

      .export-content {
        padding: 1rem 0.75rem;

        .messages-list {
          gap: 1.25rem;
        }
      }
    }
  }
}
</style>
