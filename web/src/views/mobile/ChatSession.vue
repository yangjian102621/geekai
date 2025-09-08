<template>
  <div class="app-background">
    <div class="mobile-chat" v-loading="loading" element-loading-text="正在连接会话...">
      <van-nav-bar ref="navBarRef">
        <template #title>
          <van-dropdown-menu>
            <van-dropdown-item :title="title">
              <van-cell center title="角色"> {{ role.name }}</van-cell>
              <van-cell center title="模型">{{ modelValue }}</van-cell>
            </van-dropdown-item>
          </van-dropdown-menu>
        </template>
        <template #left>
          <span class="setting">
            <van-icon name="add-o" @click="showPicker = true" />
          </span>
        </template>
        <template #right>
          <van-icon name="share-o" @click="copyShareUrl" />
        </template>
      </van-nav-bar>

      <div class="chat-list-wrapper">
        <div id="message-list-box" :style="{ height: winHeight + 'px' }" class="message-list-box">
          <van-list
            v-model:error="error"
            :finished="finished"
            error-text="请求失败，点击重新加载"
            @load="onLoad"
          >
            <van-cell v-for="item in chatData" :key="item" :border="false" class="message-line">
              <chat-prompt
                v-if="item.type === 'prompt'"
                :content="item.content"
                :icon="item.icon"
              />
              <chat-reply
                v-else-if="item.type === 'reply'"
                :content="item.content"
                :icon="item.icon"
                :org-content="item.orgContent"
                :message-id="item.id"
                :is-generating="isGenerating"
                :show-actions="item.showAction"
                :error="item.error"
                @regenerate="handleRegenerate"
              />
            </van-cell>
          </van-list>
        </div>
      </div>

      <div class="chat-box-wrapper">
        <van-sticky ref="bottomBarRef" :offset-bottom="0" position="bottom">
          <van-cell-group inset style="--van-cell-background: var(--van-cell-background-light)">
            <van-field
              v-model="prompt"
              center
              clearable
              placeholder="输入你的问题"
              type="textarea"
              rows="1"
              :autosize="{ maxHeight: 100, minHeight: 20 }"
              show-word-limit
              @keyup.enter="sendMessage"
            >
              <template #left-icon>
                <van-uploader
                  :after-read="afterRead"
                  :max-count="6"
                  :multiple="true"
                  :preview-image="false"
                  accept=".doc,.docx,.jpg,.jpeg,.png,.gif,.bmp,.webp,.svg,.ico,.xls,.xlsx,.ppt,.pptx,.pdf,.mp4,.mp3,.txt,.md,.csv,.html"
                >
                  <van-icon name="photo" />
                </van-uploader>
              </template>
              <template #button>
                <van-button size="small" type="primary" v-if="!isGenerating" @click="sendMessage"
                  >发送</van-button
                >
              </template>
              <template #extra>
                <div class="icon-box">
                  <van-icon v-if="isGenerating" name="stop-circle-o" @click="stopGenerate" />
                </div>
              </template>
            </van-field>
          </van-cell-group>
        </van-sticky>
        <MobileFileList :files="files" removable @remove="onRemovePreview" />
      </div>
    </div>
  </div>

  <van-popup v-model:show="showPicker" position="bottom" class="popup">
    <van-picker
      :columns="columns"
      v-model="selectedValues"
      title="选择模型和角色"
      @change="onChange"
      @cancel="showPicker = false"
      @confirm="newChat"
    >
      <template #option="item">
        <div class="picker-option">
          <van-image v-if="item.icon" :src="item.icon" fit="cover" round />
          <span>{{ item.text }}</span>
        </div>
      </template>
    </van-picker>
  </van-popup>
</template>

<script setup>
import ChatPrompt from '@/components/mobile/ChatPrompt.vue'
import ChatReply from '@/components/mobile/ChatReply.vue'
import { checkSession } from '@/store/cache'
import { getUserToken } from '@/store/session'
import { useSharedStore } from '@/store/sharedata'
import { showMessageError, showLoading, closeLoading } from '@/utils/dialog'
import { httpGet, httpPost } from '@/utils/http'
import MobileFileList from '@/components/mobile/MobileFileList.vue'
import { processContent, randString, renderInputText, UUID } from '@/utils/libs'
import { fetchEventSource } from '@microsoft/fetch-event-source'
import hl from 'highlight.js'
import 'highlight.js/styles/a11y-dark.css'
import { showImagePreview, showNotify, showToast } from 'vant'
import { nextTick, onMounted, onUnmounted, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { getClientId } from '@/store/cache'

const winHeight = ref(0)
const navBarRef = ref(null)
const bottomBarRef = ref(null)
const router = useRouter()

const roles = ref([])
const roleId = ref(parseInt(router.currentRoute.value.query['role_id']))
const role = ref({})
const models = ref([])
const modelId = ref(parseInt(router.currentRoute.value.query['model_id']))
const modelValue = ref('')
const title = ref(router.currentRoute.value.query['title'])
const chatId = ref(router.currentRoute.value.query['chat_id'])
const loginUser = ref(null)
const showPicker = ref(false)
const columns = ref([roles.value, models.value])
const selectedValues = ref([roleId.value, modelId.value])

checkSession()
  .then((user) => {
    loginUser.value = user
  })
  .catch(() => {
    router.push('/mobile/login')
  })

const loadModels = () => {
  // 加载模型
  httpGet('/api/model/list')
    .then((res) => {
      models.value = res.data
      if (!modelId.value) {
        modelId.value = models.value[0].id
      }
      for (let i = 0; i < models.value.length; i++) {
        models.value[i].text = models.value[i].name
        models.value[i].mValue = models.value[i].value
        models.value[i].value = models.value[i].id
      }
      modelValue.value = getModelName(modelId.value)
      // 加载角色列表
      httpGet(`/api/app/list/user`, { id: roleId.value })
        .then((res) => {
          roles.value = res.data
          if (!roleId.value) {
            roleId.value = roles.value[0]['id']
          }
          // build data for role picker
          for (let i = 0; i < roles.value.length; i++) {
            roles.value[i].text = roles.value[i].name
            roles.value[i].value = roles.value[i].id
            roles.value[i].helloMsg = roles.value[i].hello_msg
          }
          role.value = getRoleById(roleId.value)
          columns.value = [roles.value, models.value]
          selectedValues.value = [roleId.value, modelId.value]
          loadChatHistory()
        })
        .catch((e) => {
          showNotify({ type: 'danger', message: '获取聊天角色失败: ' + e.messages })
        })
    })
    .catch((e) => {
      showNotify({ type: 'danger', message: '加载模型失败: ' + e.message })
    })
}
if (chatId.value) {
  httpGet(`/api/chat/detail?chat_id=${chatId.value}`)
    .then((res) => {
      title.value = res.data.title
      modelId.value = res.data.model_id
      roleId.value = res.data.role_id
      loadModels()
    })
    .catch(() => {
      loadModels()
    })
} else {
  title.value = '新建对话'
  chatId.value = UUID()
  loadModels()
}

const chatData = ref([])
const loading = ref(false)
const finished = ref(false)
const error = ref(false)
const store = useSharedStore()
const url = ref(location.protocol + '//' + location.host + '/chat/export?chat_id=' + chatId.value)

onMounted(() => {
  winHeight.value =
    window.innerHeight - navBarRef.value.$el.offsetHeight - bottomBarRef.value.$el.offsetHeight - 70
})

const newChat = (item) => {
  showPicker.value = false
  const options = item.selectedOptions
  roleId.value = options[0].value
  modelId.value = options[1].value
  modelValue.value = getModelName(modelId.value)
  chatId.value = UUID()
  chatData.value = []
  role.value = getRoleById(roleId.value)
  title.value = '新建对话'
  loadChatHistory()
}

const onLoad = () => {
  // 加载更多消息的逻辑可以在这里实现
}

const loadChatHistory = () => {
  httpGet('/api/chat/history?chat_id=' + chatId.value)
    .then((res) => {
      const role = getRoleById(roleId.value)
      // 加载状态结束
      finished.value = true
      const data = res.data
      if (data.length === 0) {
        chatData.value.push({
          type: 'reply',
          id: randString(32),
          icon: role.icon,
          content: {
            text: role.hello_msg,
          },
          orgContent: role.hello_msg,
          showAction: false,
        })
        return
      }

      for (let i = 0; i < data.length; i++) {
        if (data[i].type === 'prompt') {
          chatData.value.push(data[i])
          continue
        }
        data[i].showAction = true
        data[i].orgContent = data[i].content.text
        chatData.value.push(data[i])
      }

      nextTick(() => {
        scrollListBox()
      })
    })
    .catch(() => {
      error.value = true
    })
}

// 创建 socket 连接
const prompt = ref('')
const previousText = ref('') // 上一次提问
const lineBuffer = ref('') // 输出缓冲行
const isGenerating = ref(false)
const isNewMsg = ref(true)
const stream = ref(store.chatStream)
const abortController = new AbortController()
watch(
  () => store.chatStream,
  (newValue) => {
    stream.value = newValue
  }
)

// 将聊天框的滚动条滑动到最底部
const scrollListBox = () => {
  document
    .getElementById('message-list-box')
    .scrollTo(0, document.getElementById('message-list-box').scrollHeight + 46)
}

// 滚动到输入区域，确保预览文件可见
const scrollToBottomBar = () => {
  try {
    // 优先让底部输入区域进入视野
    bottomBarRef.value &&
      bottomBarRef.value.$el &&
      bottomBarRef.value.$el.scrollIntoView({
        behavior: 'smooth',
        block: 'end',
      })
  } catch (e) {}
  // 再兜底滚动到页面底部
  try {
    window.scrollTo({ top: document.documentElement.scrollHeight, behavior: 'smooth' })
  } catch (e) {}
}

// 发送 SSE 请求
const sendSSERequest = async (message) => {
  try {
    isGenerating.value = true
    await fetchEventSource('/api/chat/message', {
      method: 'POST',
      headers: {
        Authorization: getUserToken(),
      },
      body: JSON.stringify(message),
      openWhenHidden: true,
      // 重试机制，避免连接断开后一直重试
      retry: 3000,
      // 设置重试延迟为0，确保不重试
      retryDelay: 3000,
      // 设置最大重试次数为0
      maxRetries: 3,
      signal: abortController.signal,
      onopen(response) {
        if (response.ok && response.status === 200) {
          console.log('SSE connection opened')
        } else {
          throw new Error(`Failed to open SSE connection: ${response.status}`)
        }
      },
      onmessage(msg) {
        try {
          const data = JSON.parse(msg.data)
          if (data.type === 'error') {
            chatData.value[chatData.value.length - 1].error = data.body
            isGenerating.value = false
            return
          }

          if (data.type === 'end') {
            isGenerating.value = false
            lineBuffer.value = '' // 清空缓冲
            isNewMsg.value = true
            return
          }

          if (data.type === 'text') {
            if (isNewMsg.value) {
              isNewMsg.value = false
              lineBuffer.value = data.body
              const reply = chatData.value[chatData.value.length - 1]
              if (reply) {
                reply['content']['text'] = lineBuffer.value
              }
            } else {
              lineBuffer.value += data.body
              const reply = chatData.value[chatData.value.length - 1]
              reply['orgContent'] = lineBuffer.value
              reply['content']['text'] = lineBuffer.value

              nextTick(() => {
                scrollListBox()

                const items = document.querySelectorAll('.message-line')
                const imgs = items[items.length - 1].querySelectorAll('img')
                for (let i = 0; i < imgs.length; i++) {
                  if (!imgs[i].src) {
                    continue
                  }
                  imgs[i].addEventListener('click', (e) => {
                    e.stopPropagation()
                    showImagePreview([imgs[i].src])
                  })
                }
              })
            }
          }

          // 回答完毕，更新完整的消息内容
          if (data.type === 'complete') {
            data.body.showAction = true
            data.body.orgContent = data.body.content.text
            chatData.value[chatData.value.length - 1] = data.body
          }
        } catch (error) {
          console.error('Error processing message:', error)
          isGenerating.value = false
          showMessageError('消息处理出错，请重试')
        }
      },
      onerror(err) {
        console.error('SSE Error:', err)
        try {
          abortController && abortController.abort()
        } catch (e) {
          console.error('AbortController abort error:', e)
        }
        isGenerating.value = false
        showMessageError('连接已断开，请重试')
      },
      onclose() {
        console.log('SSE connection closed')
        isGenerating.value = false
      },
    })
  } catch (error) {
    try {
      abortController && abortController.abort()
    } catch (e) {
      console.error('AbortController abort error:', e)
    }
    console.error('Failed to send message:', error)
    isGenerating.value = false
    showMessageError('发送消息失败，请重试')
  }
}

// 发送消息
const sendMessage = () => {
  if (isGenerating.value) {
    showToast('AI 正在作答中，请稍后...')
    return
  }

  if (prompt.value.trim().length === 0) {
    showToast('请输入需要 AI 回答的问题')
    return false
  }

  // 追加消息
  chatData.value.push({
    type: 'prompt',
    id: randString(32),
    icon: loginUser.value.avatar,
    content: { text: renderInputText(prompt.value), files: files.value },
    created_at: new Date().getTime(),
  })
  // 添加空回复消息
  const _role = getRoleById(roleId.value)
  chatData.value.push({
    chat_id: chatId,
    role_id: roleId.value,
    type: 'reply',
    id: randString(32),
    icon: _role['icon'],
    content: {
      text: '',
      files: [],
    },
  })

  nextTick(() => {
    scrollListBox()
  })

  // 发送 SSE 请求
  sendSSERequest({
    user_id: loginUser.value.id,
    role_id: roleId.value,
    model_id: modelId.value,
    chat_id: chatId.value,
    prompt: prompt.value,
    stream: stream.value,
    files: files.value,
  })

  previousText.value = prompt.value
  prompt.value = ''
  files.value = []
  return true
}

// 停止生成
const stopGenerate = function () {
  if (abortController) {
    abortController.abort()
    isGenerating.value = false
    httpGet('/api/chat/stop?session_id=' + getClientId())
      .then(() => {
        showToast('会话已中断')
      })
      .catch((e) => {
        showMessageError('中断对话失败:' + e.message)
      })
  }
}

// 处理从ChatReply组件触发的重新生成
const handleRegenerate = (messageId) => {
  if (isGenerating.value) {
    showToast('AI 正在作答中，请稍后...')
    return
  }

  console.log('messageId', messageId)
  console.log('chatData.value', chatData.value)

  // 判断 messageId 是整数
  if (messageId !== '' && isNaN(messageId)) {
    showToast('消息 ID 不合法，无法重新生成')
    return
  }

  chatData.value = chatData.value.filter((item) => item.id < messageId && !item.isHello)
  const userPrompt = chatData.value.pop()

  // 添加空回复消息
  const _role = getRoleById(roleId.value)
  chatData.value.push({
    chat_id: chatId,
    role_id: roleId.value,
    type: 'reply',
    id: randString(32),
    icon: _role['icon'],
    content: {
      text: '',
    },
  })

  // 发送 SSE 请求
  sendSSERequest({
    user_id: loginUser.value.id,
    role_id: roleId.value,
    model_id: modelId.value,
    chat_id: chatId.value,
    last_msg_id: messageId,
    prompt: userPrompt.content.text,
    stream: stream.value,
    files: [],
  })
}

const getRoleById = function (rid) {
  for (let i = 0; i < roles.value.length; i++) {
    if (roles.value[i]['id'] === rid) {
      return roles.value[i]
    }
  }
  return null
}

const getModelName = (model_id) => {
  for (let i = 0; i < models.value.length; i++) {
    if (models.value[i].id === model_id) {
      return models.value[i].text
    }
  }
  return ''
}

const onChange = (item) => {
  const selectedValues = item.selectedOptions
  if (selectedValues[0].model_id) {
    for (let i = 0; i < columns.value[1].length; i++) {
      columns.value[1][i].disabled = columns.value[1][i].value !== selectedValues[0].model_id
    }
  } else {
    for (let i = 0; i < columns.value[1].length; i++) {
      columns.value[1][i].disabled = false
    }
  }
}

// 新增复制分享链接方法
const copyShareUrl = async () => {
  try {
    await navigator.clipboard.writeText(url.value)
    showNotify({ type: 'success', message: '复制成功，请把链接发送给好友', duration: 3000 })
  } catch (e) {
    showNotify({ type: 'danger', message: '复制失败', duration: 2000 })
  }
}

// 文件上传与管理
const files = ref([])
const isHttpUrl = (url) => url.startsWith('http://') || url.startsWith('https://')
const toAbsUrl = (url) => (isHttpUrl(url) ? url : location.protocol + '//' + location.host + url)
const afterRead = async (fileItem) => {
  showLoading('文件上传中...')
  try {
    const file = Array.isArray(fileItem) ? fileItem[0].file : fileItem.file || fileItem
    const formData = new FormData()
    formData.append('file', file, file.name)
    const res = await httpPost('/api/upload', formData)
    const f = res.data || {}
    f.url = toAbsUrl(f.url || '')
    files.value = [f, ...files.value]
    // 确保上传后文件预览立即可见
    nextTick(() => {
      scrollToBottomBar()
    })
  } catch (e) {
    showNotify({ type: 'danger', message: '文件上传失败：' + (e.message || '网络错误') })
  } finally {
    closeLoading()
  }
}
const removeFile = (f, idx) => {
  files.value.splice(idx, 1)
}

const onRemovePreview = ({ file, index }) => {
  files.value.splice(index, 1)
}
</script>

<style lang="scss" scoped>
@use '../../assets/css/mobile/chat-session.scss' as *;
</style>
