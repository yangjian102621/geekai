<template>
  <div class="chat-page">
    <el-container>
      <el-aside v-show="store.chatListExtend" class="chat-aside">
        <div class="sidebar-inner">
          <div class="sidebar-header">
            <img v-if="logo" :src="logo" class="sidebar-logo" :alt="title" />
            <h2 v-else class="sidebar-title">{{ title }}</h2>
          </div>

          <div class="sidebar-actions">
            <button type="button" class="sidebar-action-item" @click="_newChat">
              <el-icon class="sidebar-action-icon"><Edit /></el-icon>
              <span>发起新对话</span>
              <el-icon class="sidebar-action-extra"><DocumentCopy /></el-icon>
            </button>
            <button type="button" class="sidebar-action-item">
              <el-icon class="sidebar-action-icon"><Star /></el-icon>
              <span>我的内容</span>
            </button>
          </div>

          <el-scrollbar :height="sidebarScrollHeight" class="sidebar-body-scroll">
            <div class="sidebar-section">
              <div
                class="sidebar-section-head sidebar-section-head-clickable"
                @click="showAgentManage = true"
              >
                <span class="sidebar-section-title">智能体</span>
                <el-icon class="sidebar-section-arrow"><ArrowRight /></el-icon>
              </div>
              <GemList
                :active-id="roleId"
                :list="pinnedList"
                mode="sidebar"
                @select="onSidebarSelectGem"
                @unpin="onSidebarUnpin"
              />
            </div>

            <div class="sidebar-section">
              <div class="sidebar-section-head">
                <span class="sidebar-section-title">对话</span>
              </div>
              <div class="search-box">
                <el-input
                  v-model="chatName"
                  placeholder="搜索会话"
                  @input="searchChat"
                  class="search-input"
                  clearable
                >
                  <template #prefix>
                    <el-icon><Search /></el-icon>
                  </template>
                </el-input>
              </div>
              <div class="content">
                <div
                  v-for="chat in chatList"
                  :key="chat.chat_id"
                  :class="[
                    'sidebar-list-item',
                    'chat-list-item',
                    { active: chat.chat_id === chatId },
                  ]"
                  @click="loadChat(chat)"
                >
                  <el-image :src="chat.icon" class="avatar" />
                  <span v-if="chat.edit" class="chat-title-input">
                    <el-input
                      v-model="tmpChatTitle"
                      size="small"
                      @keydown="titleKeydown($event, chat)"
                      :id="'chat-' + chat.chat_id"
                      @blur="editConfirm(chat)"
                      @click.stop
                      placeholder="请输入标题"
                    />
                  </span>
                  <span v-else class="sidebar-item-text">{{ chat.title }}</span>
                  <span class="sidebar-item-pin" @click.stop>
                    <el-dropdown trigger="click" @click.stop>
                      <span class="pin-trigger">
                        <i class="iconfont icon-more-vertical"></i>
                      </span>
                      <template #dropdown>
                        <el-dropdown-menu>
                          <el-dropdown-item :icon="Edit" @click="editChatTitle(chat)"
                            >重命名</el-dropdown-item
                          >
                          <el-dropdown-item
                            :icon="Delete"
                            @click="removeChat(chat)"
                            style="color: var(--el-color-danger)"
                            >删除</el-dropdown-item
                          >
                          <el-dropdown-item :icon="Share" @click="shareChat(chat)"
                            >分享</el-dropdown-item
                          >
                        </el-dropdown-menu>
                      </template>
                    </el-dropdown>
                  </span>
                </div>
              </div>
            </div>
          </el-scrollbar>

          <div class="sidebar-footer">
            <el-button
              type="primary"
              size="small"
              plain
              class="clear-chats-btn"
              @click="clearAllChats"
            >
              <el-icon><Delete /></el-icon>
              <span>清除所有对话</span>
            </el-button>
          </div>
        </div>
      </el-aside>

      <el-main
        v-loading="loading"
        element-loading-background="rgba(122, 122, 122, 0.3)"
        class="relative"
      >
        <div class="absolute top-2 left-2 cursor-pointer">
          <div @click="store.setChatListExtend(!store.chatListExtend)">
            <el-tooltip content="隐藏对话列表" placement="right" v-if="store.chatListExtend">
              <i class="iconfont icon-colspan text-xl"></i>
            </el-tooltip>
            <el-tooltip content="展开对话列表" placement="right" v-else>
              <i class="iconfont icon-expand text-xl"></i>
            </el-tooltip>
          </div>
        </div>

        <div class="chat-container">
          <div class="flex justify-center">
            <div id="container" :style="{ height: mainWinHeight + 'px' }">
              <!-- 智能体管理视图：仅点击侧栏「智能体」时展示 -->
              <div
                v-if="showAgentManage"
                class="agent-manage-view"
                :style="{ minHeight: chatBoxHeight + 'px' }"
              >
                <AgentManagePage
                  :profile="userProfile"
                  :refresh-trigger="agentListRefreshTrigger"
                  @select="onManageSelectGem"
                  @new="openAgentForm"
                  @edit="openAgentFormForEdit"
                  @profile-update="onAgentManageProfileUpdate"
                />
              </div>
              <!-- 对话首页：欢迎语 + 居中输入框 -->
              <div
                v-else-if="showHello"
                class="chat-box welcome-wrap"
                :style="{ height: chatBoxHeight + 'px' }"
              >
                <div class="welcome-home">
                  <!-- 智能体名片：仅当 hello_msg 不为空时显示 -->
                  <div v-if="currentAgentWithHello" class="w-full max-w-[640px] px-5">
                    <div class="agent-welcome-card">
                      <div class="agent-welcome-card__avatar-wrap">
                        <el-image
                          v-if="currentAgentWithHello.icon"
                          :src="currentAgentWithHello.icon"
                          class="agent-welcome-card__avatar"
                          fit="cover"
                        />
                        <span v-else class="agent-welcome-card__avatar-ph" />
                      </div>
                      <div class="agent-welcome-card__body">
                        <div class="agent-welcome-card__name">{{ currentAgentWithHello.name }}</div>
                        <p class="agent-welcome-card__hello">
                          {{ currentAgentWithHello.hello_msg }}
                        </p>
                      </div>
                    </div>
                  </div>
                  <div class="flex flex-col items-start gap-2 mb-4 w-full max-w-[640px] px-5">
                    <div class="text-2xl flex items-center justify-start w-full gap-2">
                      <i class="iconfont icon-rebot !text-3xl text-purple-500"></i>你好！
                    </div>
                    <div class="w-full">
                      <p class="text-4xl">今天需要我为你做些什么？</p>
                    </div>
                  </div>
                  <div class="welcome-input-wrap">
                    <ChatInputBox
                      v-model="prompt"
                      :files="files"
                      :model-id="modelID"
                      :models="displayedModels"
                      :tools="tools"
                      :tool-selected="toolSelected"
                      :is-generating="isGenerating"
                      :disable-model="disableModel"
                      :user-id="loginUser?.id"
                      :placeholder="'给 ' + (selectedModel?.name || 'AI') + ' 发送消息'"
                      :agent-name="currentAgentName"
                      @update:model-value="prompt = $event"
                      @update:model-id="modelID = $event"
                      @update:tool-selected="toolSelected = $event"
                      @send="sendMessage()"
                      @stop="stopGenerate"
                      @paste="onPaste"
                      @input="onInput"
                      @file-selected="insertFile"
                      @file-removed="removeFile"
                      @show-setting="showChatSetting = true"
                    />
                  </div>
                </div>
              </div>
              <div v-else class="chat-box" id="chat-box" :style="{ height: chatBoxHeight + 'px' }">
                <div v-for="(item, index) in chatData" :key="item.id">
                  <chat-prompt
                    v-if="item.type === 'prompt'"
                    :data="item"
                    :message-index="index"
                    @edit="editUserPrompt"
                  />
                  <chat-reply
                    v-else-if="item.type === 'reply'"
                    :data="item"
                    @regen="reGenerate"
                    :message-index="index"
                  />
                </div>

                <back-top :right="30" :bottom="170" />
              </div>
              <!-- end chat box -->

              <div v-if="!(showAgentManage || showHello)" class="input-box">
                <ChatInputBox
                  v-model="prompt"
                  :files="files"
                  :model-id="modelID"
                  :models="models"
                  :tools="tools"
                  :tool-selected="toolSelected"
                  :is-generating="isGenerating"
                  :disable-model="disableModel"
                  :user-id="loginUser?.id"
                  :agent-name="currentAgentName"
                  placeholder="按 Enter 键发送消息，使用 Shift + Enter 换行"
                  @update:model-value="prompt = $event"
                  @update:model-id="modelID = $event"
                  @update:tool-selected="toolSelected = $event"
                  @send="sendMessage()"
                  @stop="stopGenerate"
                  @paste="onPaste"
                  @input="onInput"
                  @file-selected="insertFile"
                  @file-removed="removeFile"
                  @show-setting="showChatSetting = true"
                />
                <!-- end input box -->
              </div>
            </div>
            <!-- end container -->
          </div>
          <!-- end loading -->
        </div>
      </el-main>
    </el-container>

    <el-drawer
      v-model="showAgentFormDrawer"
      :title="agentFormDrawerTitle"
      direction="rtl"
      size="400px"
    >
      <el-form :model="agentForm" label-width="80px" label-position="left">
        <el-form-item label="名称" required>
          <el-input
            v-model="agentForm.name"
            placeholder="为你的智能体命名"
            maxlength="30"
            show-word-limit
          />
        </el-form-item>
        <el-form-item label="指令">
          <el-input
            v-model="agentForm.system_prompt"
            type="textarea"
            :rows="5"
            placeholder="系统提示词 / 角色设定"
          />
        </el-form-item>
        <el-form-item label="打招呼">
          <el-input
            v-model="agentForm.hello_msg"
            placeholder="首次进入时的问候语"
            maxlength="255"
            show-word-limit
          />
        </el-form-item>
        <el-form-item label="图标">
          <ImageUpload v-model="agentForm.icon" :max-count="1" />
          <div class="text-gray-500 text-xs mt-1">上传智能体头像，不传则使用默认图标</div>
        </el-form-item>
        <el-form-item label="绑定模型">
          <el-select v-model="agentForm.model_id" placeholder="可选" clearable style="width: 100%">
            <el-option v-for="m in models" :key="m.id" :label="m.name" :value="m.id" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showAgentFormDrawer = false">取消</el-button>
        <el-button type="primary" @click="submitAgentForm">保存</el-button>
      </template>
    </el-drawer>

    <ChatSetting :show="showChatSetting" @hide="showChatSetting = false" />

    <el-dialog v-model="showConversationDialog" title="实时语音通话" :fullscreen="true">
      <div v-loading="!frameLoaded">
        <iframe
          style="width: 100%; height: calc(100vh - 100px); border: none"
          :src="voiceChatUrl"
          @load="frameLoaded = true"
          allow="microphone *; camera *"
        ></iframe>
      </div>
    </el-dialog>
  </div>
</template>
<script setup>
import BackTop from '@/components/BackTop.vue'
import ChatPrompt from '@/components/ChatPrompt.vue'
import ChatReply from '@/components/ChatReply.vue'
import ChatSetting from '@/components/ChatSetting.vue'
import GemList from '@/components/GemList.vue'
import FileList from '@/components/FileList.vue'
import FileSelect from '@/components/FileSelect.vue'
import AgentManagePage from '@/components/AgentManagePage.vue'
import ImageUpload from '@/components/ImageUpload.vue'
import ChatInputBox from '@/components/ChatInputBox.vue'
import { checkSession, getClientId, getSystemInfo } from '@/store/cache'
import { useSharedStore } from '@/store/sharedata'
import { closeLoading, showLoading, showMessageError, showMessageInfo } from '@/utils/dialog'
import { httpGet, httpPost } from '@/utils/http'
import { isMobile, randString, removeArrayItem, UUID } from '@/utils/libs'
import {
  ArrowRight,
  Cpu,
  Delete,
  DocumentCopy,
  Edit,
  InfoFilled,
  More,
  Promotion,
  Rank,
  Search,
  Share,
  Star,
  VideoPause,
} from '@element-plus/icons-vue'
import { fetchEventSource } from '@microsoft/fetch-event-source'
import Clipboard from 'clipboard'
import { ElMessage, ElMessageBox } from 'element-plus'
import 'highlight.js/styles/a11y-dark.css'
import { computed, nextTick, onMounted, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { getUserToken } from '../store/session'
import { substr } from '../utils/libs'

const title = ref('GeekAI-智能助手')
const logo = ref('')
const models = ref([])
const modelID = ref(0)
const chatData = ref([])
const allChats = ref([]) // 会话列表
const chatList = ref(allChats.value)
const mainWinHeight = ref(0) // 主窗口高度
const chatBoxHeight = ref(0) // 聊天内容框高度
const sidebarScrollHeight = ref('400px')
const loading = ref(false)
const loginUser = ref(null)
const roles = ref([])
const router = useRouter()
const roleId = ref(0)
const chatId = ref()
const newChatItem = ref(null)
const isLogin = ref(false)
const showHello = ref(true)
const showAgentManage = ref(false)
const welcomeAgents = ref([])
const userProfile = ref({ gem_ids: [] })
const allAgentsList = ref([])
const agentListRefreshTrigger = ref(0)
const showAgentFormDrawer = ref(false)
const agentForm = ref({
  id: 0,
  name: '',
  system_prompt: '',
  hello_msg: '',
  icon: '/images/avatar/gpt.png',
  model_id: 0,
})
const agentFormDrawerTitle = computed(() => (agentForm.value.id ? '编辑智能体' : '新建智能体'))
const currentAgentName = computed(() => {
  const role = getRoleById(roleId.value)
  return role && role.name ? role.name : ''
})

// 当前智能体且存在 hello_msg 时用于欢迎页名片展示
const currentAgentWithHello = computed(() => {
  const role = getRoleById(roleId.value)
  if (!role) return null
  const msg = typeof role.hello_msg === 'string' ? role.hello_msg.trim() : ''
  return msg ? { ...role, hello_msg: msg } : null
})
const inputRef = ref(null)
const textHeightRef = ref(null)

const store = useSharedStore()
const row = ref(3)
const showChatSetting = ref(false)
const listStyle = ref(store.chatListStyle)
const config = ref({ advance_voice_power: 0 })
const voiceChatUrl = ref('')
const modelSearchKeyword = ref('') // 模型搜索关键词
const selectedCategory = ref('')
const modelCategories = ref([])
const groupedModels = ref([])
const activeCategory = ref('') // 当前激活的分类标签
const showFreeModelsOnly = ref(false) // 是否只显示免费模型

const tools = ref([])
const toolSelected = ref([])
const stream = ref(store.chatStream)
const modelSelectorRef = ref(null)

// 侧栏固定智能体列表：按 gem_ids 顺序，最多 8 个（id 统一按数字比较）
const pinnedList = computed(() => {
  const rawIds = userProfile.value?.gem_ids || []
  const ids = Array.from(
    new Set(
      rawIds
        .map((v) => (typeof v === 'number' ? v : Number(v)))
        .filter((n) => !Number.isNaN(n) && n > 0)
    )
  )
  const list = allAgentsList.value || []
  const byId = new Map(list.map((a) => [a.id, a]))
  const ordered = []
  for (const id of ids) {
    if (byId.has(id)) ordered.push(byId.get(id))
  }
  return ordered.slice(0, 8)
})

// 过滤后的模型列表
const filteredModels = computed(() => {
  if (!modelSearchKeyword.value && !showFreeModelsOnly.value && !activeCategory.value) {
    return models.value
  }

  return models.value.filter((model) => {
    // 搜索关键词匹配
    const matchesSearch =
      !modelSearchKeyword.value ||
      model.name.toLowerCase().includes(modelSearchKeyword.value.toLowerCase()) ||
      (model.description &&
        model.description.toLowerCase().includes(modelSearchKeyword.value.toLowerCase()))

    // 分类匹配
    const matchesCategory = !activeCategory.value || model.tag === activeCategory.value

    // 免费模型匹配
    const matchesFree = !showFreeModelsOnly.value || model.power <= 0

    return matchesSearch && matchesCategory && matchesFree
  })
})

// 最终展示的模型列表
const displayedModels = computed(() => {
  return filteredModels.value
})

// 切换是否只显示免费模型
const toggleFreeModels = () => {
  showFreeModelsOnly.value = !showFreeModelsOnly.value
  if (showFreeModelsOnly.value) {
    activeCategory.value = ''
  }
}

// 提取所有模型分类
const updateModelCategories = () => {
  const categories = new Set()
  models.value.forEach((model) => {
    if (model.tag) {
      categories.add(model.tag)
    }
  })
  modelCategories.value = Array.from(categories)
}

// 按分类对模型进行分组
const updateGroupedModels = () => {
  const filtered = filteredModels.value

  // 如果已经指定分类，则只显示该分类
  if (selectedCategory.value) {
    groupedModels.value = [
      {
        category: selectedCategory.value,
        models: filtered,
      },
    ]
    return
  }

  // 否则按分类分组展示
  const groups = {}
  filtered.forEach((model) => {
    const category = model.tag || '未分类'
    if (!groups[category]) {
      groups[category] = []
    }
    groups[category].push(model)
  })

  groupedModels.value = Object.keys(groups).map((category) => ({
    category,
    models: groups[category],
  }))

  // 对分组进行排序（未分类放最后）
  groupedModels.value.sort((a, b) => {
    if (a.category === '未分类') return 1
    if (b.category === '未分类') return -1
    return a.category.localeCompare(b.category)
  })
}

// 当筛选条件变化时更新分组
watch([filteredModels, selectedCategory], () => {
  updateGroupedModels()
})

// 监听模型数据变化，更新分类列表
watch(
  () => models.value,
  () => {
    updateModelCategories()
    updateGroupedModels()
  },
  { deep: true }
)

// 获取选中的模型名称
const selectedModel = computed(() => {
  const model = getSelectedModel()
  return model ? model : { name: '选择模型', power: 0 }
})

// 获取选中的模型
const getSelectedModel = () => {
  return models.value.find((model) => model.id === modelID.value)
}

// 选择模型
const selectModel = (model) => {
  modelID.value = model.id
  modelSelectorRef.value.hide()
  _newChat()
}

// 根据算力获取标签类型
const getTagType = (power) => {
  const powerNum = Number(power)
  if (powerNum <= 5) return 'info'
  if (powerNum <= 15) return 'warning'
  return 'danger'
}

watch(
  () => store.chatListStyle,
  (newValue) => {
    listStyle.value = newValue
  }
)

watch(
  () => store.chatStream,
  (newValue) => {
    stream.value = newValue
  }
)

// 初始化角色ID参数
if (router.currentRoute.value.query.role_id) {
  roleId.value = parseInt(router.currentRoute.value.query.role_id)
}

// 初始化 ChatID
chatId.value = router.currentRoute.value.params.id
if (!chatId.value) {
  chatId.value = UUID()
} else {
  // 查询对话信息
  httpGet('/api/chat/detail', { chat_id: chatId.value })
    .then((res) => {
      document.title = res.data.title
      roleId.value = res.data.role_id
      modelID.value = res.data.model_id
      const role = getRoleById(res.data.role_id)
      disableModel.value = !!(role && role.model_id > 0)
    })
    .catch((e) => {
      console.error('获取对话信息失败：' + e.message)
    })
}

// 获取系统配置
getSystemInfo()
  .then((res) => {
    config.value = res.data
    title.value = config.value.title
    logo.value = res.data.bar_logo
  })
  .catch((e) => {
    ElMessage.error('获取系统配置失败：' + e.message)
  })

// 获取工具函数
httpGet('/api/function/list')
  .then((res) => {
    tools.value = res.data
  })
  .catch((e) => {
    showMessageError('获取工具函数失败：' + e.message)
  })

const prompt = ref('')
const isGenerating = ref(false)
const lineBuffer = ref('') // 输出缓冲行
const isNewMsg = ref(true)
const abortController = ref(null)

onMounted(() => {
  resizeElement()
  initData()

  const clipboard = new Clipboard('.copy-reply, .copy-code-btn')
  clipboard.on('success', () => {
    ElMessage.success('复制成功！')
  })

  clipboard.on('error', () => {
    ElMessage.error('复制失败！')
  })

  window.onresize = () => resizeElement()
})

// 初始化数据
const initData = async () => {
  try {
    // 获取角色列表
    const roleRes = await httpGet('/api/app/list')
    roles.value = roleRes.data
    if (roles.value.length > 0 && !roleId.value) {
      roleId.value = roles.value[0].id
    }

    // 获取模型列表
    const modelRes = await httpGet('/api/model/list')
    models.value = modelRes.data
    if (models.value.length > 0 && !modelID.value) {
      modelID.value = models.value[0].id
    }
    if (roleId.value) {
      const r = getRoleById(roleId.value)
      disableModel.value = !!(r && r.model_id > 0)
    }

    // 获取用户信息
    const user = await checkSession()
    loginUser.value = user
    isLogin.value = true

    // 用户 profile（含 gem_ids 用于侧栏固定列表）
    try {
      const profileRes = await httpGet('/api/user/profile')
      userProfile.value = profileRes.data || { gem_ids: [] }
      if (!Array.isArray(userProfile.value.gem_ids)) {
        userProfile.value.gem_ids = []
      }
    } catch (_) {
      userProfile.value = { gem_ids: [] }
    }

    // 全量智能体列表（用于侧栏 pinned 与 Welcome）
    try {
      const listRes = await httpGet('/api/app/list/user')
      const list = listRes.data || []
      allAgentsList.value = list
      welcomeAgents.value = list.length > 9 ? list.slice(0, 9) : list
    } catch (_) {
      allAgentsList.value = roles.value || []
      welcomeAgents.value = (roles.value || []).slice(0, 9)
    }

    // 获取聊天列表
    const chatRes = await httpGet('/api/chat/list')
    allChats.value = chatRes.data
    chatList.value = allChats.value
    if (chatId.value) {
      loadChatHistory(chatId.value)
    }
  } catch (error) {
    if (error.response?.status === 401) {
      isLogin.value = false
      userProfile.value = { gem_ids: [] }
      allAgentsList.value = roles.value || []
      welcomeAgents.value = (roles.value || []).slice(0, 9)
    } else {
      console.warn('初始化数据失败：' + error.message)
    }
  }
}
abortController.value = new AbortController()
// 发送 SSE 请求
const sendSSERequest = async (message) => {
  isGenerating.value = true
  try {
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
      signal: abortController.value.signal,
      onopen(response) {
        if (response.ok && response.status === 200) {
          console.log('SSE connection opened')
        } else {
          console.error('SSE connection failed', response)
          isGenerating.value = false
        }
      },
      onmessage(msg) {
        try {
          const data = JSON.parse(msg.data)
          if (data.type === 'error') {
            const reply = chatData.value[chatData.value.length - 1]
            if (reply) {
              reply['content'].text = `<div class="text-red-500 rounded-md">${data.body}</div>`
            }
            isGenerating.value = false
            return
          }

          if (data.type === 'end') {
            isGenerating.value = false
            lineBuffer.value = '' // 清空缓冲

            // 获取 token
            const reply = chatData.value[chatData.value.length - 1]
            httpPost('/api/chat/tokens', {
              text: '',
              model: getModelValue(modelID.value),
              chat_id: chatId.value,
            })
              .then((res) => {
                reply['created_at'] = new Date().getTime()
                reply['tokens'] = res.data
                // 将聊天框的滚动条滑动到最底部
                nextTick(() => {
                  document
                    .getElementById('chat-box')
                    .scrollTo(0, document.getElementById('chat-box').scrollHeight)
                })
              })
              .catch(() => {})
            isNewMsg.value = true
            tmpChatTitle.value = message.prompt
            console.log('chatData.value', chatData.value)
            // 判断 chatlist 中指定的 chat_id 是否存在
            const chat = chatList.value.find((chat) => chat.chat_id === chatId.value)
            if (!chat) {
              const _role = getRoleById(roleId.value)
              chatList.value.unshift({
                chat_id: chatId.value,
                title: substr(message.prompt, 15),
                role_id: roleId.value,
                model_id: modelID.value,
                icon: _role.icon,
                created_at: new Date().getTime(),
                updated_at: new Date().getTime(),
              })
            }
            return
          }

          if (data.type === 'text') {
            if (isNewMsg.value) {
              isNewMsg.value = false
              lineBuffer.value = data.body
              const reply = chatData.value[chatData.value.length - 1]
              if (reply) {
                reply['content'].text = lineBuffer.value
              }
            } else {
              lineBuffer.value += data.body
              const reply = chatData.value[chatData.value.length - 1]
              if (reply) {
                reply['content'].text = lineBuffer.value
              }
            }
          }

          // 回答完毕，更新完整的消息内容
          if (data.type === 'complete') {
            chatData.value[chatData.value.length - 1] = data.body
          }

          // 将聊天框的滚动条滑动到最底部
          nextTick(() => {
            document
              .getElementById('chat-box')
              .scrollTo(0, document.getElementById('chat-box').scrollHeight)
            localStorage.setItem('chat_id', chatId.value)
          })
        } catch (error) {
          console.error('Error processing message:', error)
          isGenerating.value = false
          ElMessage.error('消息处理出错，请重试')
        }
      },
      onerror(err) {
        console.error('SSE Error:', err)
        try {
          abortController.value && abortController.value.abort()
        } catch (e) {
          console.error('AbortController abort error:', e)
        }
        isGenerating.value = false
        // ElMessage.error('连接已断开，发生错误：' + err.message)
        const reply = chatData.value[chatData.value.length - 1]
        if (reply) {
          reply['content'].text = `<div class="text-red-500 rounded-md">${err.message}</div>`
        }
      },
      onclose() {
        console.log('SSE connection closed')
        isGenerating.value = false
      },
    })
  } catch (error) {
    console.error('Failed to send message:', error)
    isGenerating.value = false
    ElMessage.error('发送消息失败，请重试')
  }
}

// 发送消息
const sendMessage = (messageId = 0) => {
  if (!isLogin.value) {
    console.log('未登录')
    store.setShowLoginDialog(true)
    return
  }

  if (isGenerating.value) {
    ElMessage.warning('AI 正在作答中，请稍后...')
    return
  }

  if (prompt.value === '') {
    showMessageError('请输入要发送的消息！')
    return false
  }

  // 追加消息
  chatData.value.push({
    type: 'prompt',
    id: 0,
    icon: loginUser.value.avatar,
    content: {
      text: prompt.value,
      files: files.value,
    },
    model: getModelValue(modelID.value),
    created_at: new Date().getTime() / 1000,
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
    document
      .getElementById('chat-box')
      .scrollTo(0, document.getElementById('chat-box').scrollHeight)
  })

  showHello.value = false

  // 异步发送 SSE 请求
  sendSSERequest({
    user_id: loginUser.value.id,
    role_id: roleId.value,
    model_id: modelID.value,
    chat_id: chatId.value,
    prompt: prompt.value,
    tools: toolSelected.value,
    stream: stream.value,
    files: files.value,
    last_msg_id: messageId || 0,
  })

  prompt.value = ''
  files.value = []
  row.value = 1
}

const getRoleById = function (rid) {
  for (let i = 0; i < roles.value.length; i++) {
    if (roles.value[i]['id'] === rid) {
      return roles.value[i]
    }
  }
  const fromAll = (allAgentsList.value || []).find((a) => a.id === rid)
  if (fromAll) return fromAll
  return null
}

const resizeElement = function () {
  const h = window.innerHeight
  sidebarScrollHeight.value = `${Math.max(280, h - 230)}px`
  mainWinHeight.value = h
  chatBoxHeight.value = h - 101 - 82
}

const _newChat = () => {
  if (!isLogin.value) {
    store.setShowLoginDialog(true)
    return
  }

  // Generate new chat ID
  chatId.value = UUID()

  // Reset chat state
  chatData.value = []
  prompt.value = ''
  files.value = []
  isGenerating.value = false

  // Close agent manage view if open
  showAgentManage.value = false

  // Call newChat to set up the new conversation
  newChat()

  // 发起新对话时进入欢迎页
  showHello.value = true
}
const disableModel = ref(false)
// 新建会话
const newChat = () => {
  if (!isLogin.value) {
    store.setShowLoginDialog(true)
    return
  }
  const role = getRoleById(roleId.value)
  showHello.value = role.name === 'GPT'
  // if the role bind a model, disable model change
  disableModel.value = false
  if (role.model_id > 0) {
    modelID.value = role.model_id
    disableModel.value = true
  }
  // 已有新开的会话
  if (newChatItem.value !== null && newChatItem.value['role_id'] === roles.value[0]['role_id']) {
    return
  }

  // 获取当前聊天角色图标
  let icon = ''
  roles.value.forEach((item) => {
    if (item['id'] === roleId.value) {
      icon = item['icon']
    }
  })
  newChatItem.value = {
    chat_id: '',
    icon: icon,
    role_id: roleId.value,
    model_id: modelID.value,
    title: '',
    edit: false,
    removing: false,
  }
  isGenerating.value = false
  loadChatHistory(chatId.value)
  router.push(`/chat/${chatId.value}`)
}

const handleSelectGem = (gem) => {
  roleId.value = gem.id
  disableModel.value = false
  if (gem.model_id > 0) {
    modelID.value = gem.model_id
    disableModel.value = true
  }
  _newChat()
}

// 侧栏点击智能体：进入对话，关闭管理视图
const onSidebarSelectGem = (gem) => {
  showAgentManage.value = false
  handleSelectGem(gem)
}

// 刷新用户 profile（用于 pin/unpin 后更新侧栏）
const refreshProfile = async () => {
  try {
    const res = await httpGet('/api/user/profile')
    userProfile.value = res.data || { gem_ids: [] }
    if (!Array.isArray(userProfile.value.gem_ids)) {
      userProfile.value.gem_ids = []
    }
  } catch (_) {
    userProfile.value = { gem_ids: [] }
  }
}

// 统一转为数字并去重，避免字符串/数字混用及重复 ID
const toNum = (v) => (typeof v === 'number' && !Number.isNaN(v) ? v : Number(v))
const normalizeGemIds = (arr) => {
  const nums = (Array.isArray(arr) ? arr : []).map(toNum).filter((n) => !Number.isNaN(n) && n > 0)
  return Array.from(new Set(nums))
}

// 侧栏取消固定
const onSidebarUnpin = async (id) => {
  const raw = userProfile.value?.gem_ids || []
  const target = toNum(id)
  const ids = normalizeGemIds(raw).filter((x) => x !== target)
  try {
    await httpPost('/api/user/profile/update', { ...userProfile.value, gem_ids: ids })
    userProfile.value = { ...userProfile.value, gem_ids: ids }
  } catch (e) {
    ElMessage.error('取消固定失败：' + (e.message || ''))
  }
}

// 管理页 pin/unpin 后更新本地 profile（与侧栏联动）
const onAgentManageProfileUpdate = (gemIds) => {
  userProfile.value = { ...userProfile.value, gem_ids: normalizeGemIds(gemIds || []) }
}

// 管理页点击智能体：发起新对话，关闭管理视图
const onManageSelectGem = (gem) => {
  showAgentManage.value = false
  handleSelectGem(gem)
}

// 新建智能体：打开抽屉
const openAgentForm = () => {
  agentForm.value = {
    id: 0,
    name: '',
    system_prompt: '',
    hello_msg: '',
    icon: '/images/avatar/gpt.png',
    model_id: 0,
  }
  showAgentFormDrawer.value = true
}

// 编辑智能体：填充表单并打开抽屉
const openAgentFormForEdit = (agent) => {
  agentForm.value = {
    id: agent.id,
    name: agent.name || '',
    system_prompt: agent.system_prompt || '',
    hello_msg: agent.hello_msg || '',
    icon: agent.icon || '/images/avatar/gpt.png',
    model_id: agent.model_id || 0,
  }
  showAgentFormDrawer.value = true
}

// 提交智能体表单（新建或更新）
const submitAgentForm = async () => {
  if (!agentForm.value.name?.trim()) {
    ElMessage.warning('请输入智能体名称')
    return
  }
  const isEdit = !!agentForm.value.id
  try {
    if (isEdit) {
      await httpPost('/api/app/update', {
        id: agentForm.value.id,
        name: agentForm.value.name.trim(),
        system_prompt: agentForm.value.system_prompt || '',
        hello_msg: agentForm.value.hello_msg || '',
        icon: agentForm.value.icon || '/images/avatar/gpt.png',
        model_id: agentForm.value.model_id || 0,
      })
      ElMessage.success('保存成功')
    } else {
      await httpPost('/api/app/create', {
        name: agentForm.value.name.trim(),
        system_prompt: agentForm.value.system_prompt || '',
        hello_msg: agentForm.value.hello_msg || '',
        icon: agentForm.value.icon || '/images/avatar/gpt.png',
        model_id: agentForm.value.model_id || 0,
        tid: 0,
        sort: 0,
      })
      ElMessage.success('创建成功')
    }
    showAgentFormDrawer.value = false
    agentListRefreshTrigger.value++
    const listRes = await httpGet('/api/app/list/user')
    allAgentsList.value = listRes.data || []
  } catch (e) {
    ElMessage.error(isEdit ? '保存失败：' + (e.message || '') : '创建失败：' + (e.message || ''))
  }
}

// 切换会话
const loadChat = function (chat) {
  if (!isLogin.value) {
    store.setShowLoginDialog(true)
    return
  }

  if (chatId.value === chat.chat_id) {
    return
  }
  showAgentManage.value = false
  newChatItem.value = null
  roleId.value = chat.role_id
  modelID.value = chat.model_id
  chatId.value = chat.chat_id
  const role = getRoleById(chat.role_id)
  disableModel.value = !!(role && role.model_id > 0)
  isGenerating.value = false
  loadChatHistory(chatId.value)
  router.push(`/chat/${chatId.value}`)
}

// 编辑会话标题
const tmpChatTitle = ref('')
const editChatTitle = (chat) => {
  chat.edit = true
  tmpChatTitle.value = chat.title
  nextTick(() => {
    document.getElementById('chat-' + chat.chat_id).focus()
  })
}

const titleKeydown = (e, chat) => {
  if (e.keyCode === 13) {
    e.stopPropagation()
    editConfirm(chat)
  }
}

const stopPropagation = (e) => {
  e.stopPropagation()
}
// 确认修改
const editConfirm = function (chat) {
  if (tmpChatTitle.value === '') {
    return ElMessage.error('请输入会话标题！')
  }
  if (!chat.chat_id) {
    return ElMessage.error('对话 ID 为空，请刷新页面再试！')
  }
  if (tmpChatTitle.value === chat.title) {
    chat.edit = false
    return
  }

  httpPost('/api/chat/update', {
    chat_id: chat.chat_id,
    title: tmpChatTitle.value,
  })
    .then(() => {
      chat.title = tmpChatTitle.value
      chat.edit = false
    })
    .catch((e) => {
      ElMessage.error('操作失败：' + e.message)
    })
}
// 删除会话
const removeChat = function (chat) {
  ElMessageBox.confirm(`该操作会删除"${chat.title}"`, '删除聊天', {
    confirmButtonText: '删除',
    cancelButtonText: '取消',
    type: 'warning',
  })
    .then(() => {
      httpGet('/api/chat/remove?chat_id=' + chat.chat_id)
        .then(() => {
          chatList.value = removeArrayItem(chatList.value, chat, function (e1, e2) {
            return e1.id === e2.id
          })
          // 重置会话
          _newChat()
        })
        .catch((e) => {
          ElMessage.error('操作失败：' + e.message)
        })
    })
    .catch(() => {})
}

const onInput = (e) => {
  const inputEl = inputRef.value?.$el?.querySelector?.('textarea') ?? inputRef.value
  if (!inputEl) return
  const lineHeight = parseFloat(window.getComputedStyle(inputEl).lineHeight)
  textHeightRef.value.style.width = inputEl.clientWidth + 'px'
  const lines = Math.floor(textHeightRef.value.clientHeight / lineHeight)
  //inputRef.value.scrollTo(0, inputRef.value.scrollHeight)
  if (prompt.value.length < 10) {
    row.value = 1
  } else if (lines <= 7) {
    row.value = lines
  } else {
    row.value = 7
  }

  // 输入回车自动提交
  if (e.keyCode === 13) {
    // Shift + Enter 换行
    if (e.shiftKey) {
      return
    }
    e.preventDefault()
    sendMessage()
  }
}

// 自动填充 prompt
const autofillPrompt = (text) => {
  prompt.value = text
  inputRef.value.focus()
  sendMessage()
}

const clearAllChats = function () {
  ElMessageBox.confirm('清除所有对话?此操作不可撤销！', '警告', {
    confirmButtonText: '删除对话',
    cancelButtonText: '取消',

    dangerouslyUseHTMLString: true,
    showClose: true,
    closeOnClickModal: false,
    center: false,
  })
    .then(() => {
      httpGet('/api/chat/clear')
        .then(() => {
          ElMessage.success('操作成功！')
          chatData.value = []
          chatList.value = []
          newChat()
        })
        .catch((e) => {
          ElMessage.error('操作失败：' + e.message)
        })
    })
    .catch(() => {})
}

const loadChatHistory = function (chatId) {
  chatData.value = []
  loading.value = true
  httpGet('/api/chat/history?chat_id=' + chatId)
    .then((res) => {
      loading.value = false
      const data = res.data
      if ((!data || data.length === 0) && chatData.value.length === 0) {
        // 加载打招呼信息
        const _role = getRoleById(roleId.value)
        chatData.value.push({
          chat_id: chatId,
          role_id: roleId.value,
          type: 'reply',
          id: 0,
          icon: _role['icon'],
          isHello: true,
          content: {
            text: _role['hello_msg'],
            files: [],
          },
        })
        return
      }
      showHello.value = false
      for (let i = 0; i < data.length; i++) {
        if (data[i].type === 'reply' && i > 0) {
          data[i].prompt = data[i - 1].content
        }
        chatData.value.push(data[i])
      }

      nextTick(() => {
        document
          .getElementById('chat-box')
          .scrollTo(0, document.getElementById('chat-box').scrollHeight)
      })
    })
    .catch((e) => {
      // TODO: 显示重新加载按钮
      ElMessage.error('加载聊天记录失败：' + e.message)
    })
}

// 停止生成
const stopGenerate = function () {
  if (abortController.value) {
    abortController.value.abort()
    isGenerating.value = false
    httpGet('/api/chat/stop?session_id=' + getClientId())
      .then(() => {
        showMessageInfo('会话已中断')
      })
      .catch((e) => {
        showMessageError('中断对话失败:' + e.message)
      })
  }
}

// 重新生成
const reGenerate = function (messageIndex) {
  // 恢复发送按钮状态
  if (isGenerating.value) {
    ElMessage.warning('AI 正在作答中，请稍后...')
    return
  }

  if (messageIndex === -1 || isNaN(messageIndex)) {
    ElMessage.error('找不到要重新生成消息')
    return
  }

  // 找到该消息的ID
  const messageId = chatData.value[messageIndex].id
  // 移除该消息之后的所有消息
  chatData.value = chatData.value.slice(0, messageIndex)
  const userPrompt = chatData.value.pop()
  prompt.value = userPrompt.content.text
  // 将光标定位到输入框并聚焦
  nextTick(() => {
    sendMessage(messageId)
    if (inputRef.value) {
      inputRef.value.focus()
      // 触发输入事件以更新文本高度
      onInput({ keyCode: null })
    }
  })
}

// 编辑用户消息提交
const editUserPrompt = function (editData) {
  const { messageIndex, newContent } = editData

  if (messageIndex === -1 || isNaN(messageIndex)) {
    ElMessage.error('找不到要编辑的消息')
    return
  }

  // 找到该消息下一条消息的ID
  const messageId = chatData.value[messageIndex + 1].id
  // 移除该消息之后的所有消息
  chatData.value = chatData.value.slice(0, messageIndex)
  console.log('chatData.value', chatData.value)
  // 设置该消息的内容
  prompt.value = newContent
  console.log('messageId', messageId)
  // 发送消息
  nextTick(() => {
    sendMessage(messageId)
    if (inputRef.value) {
      inputRef.value.focus()
      // 触发输入事件以更新文本高度
      onInput({ keyCode: null })
    }
  })
}

const chatName = ref('')
// 搜索会话（实时过滤）
const searchChat = function () {
  const keyword = chatName.value.trim().toLowerCase()
  if (!keyword) {
    chatList.value = allChats.value
    return
  }
  const items = []
  for (let i = 0; i < allChats.value.length; i++) {
    if (allChats.value[i].title.toLowerCase().indexOf(keyword) !== -1) {
      items.push(allChats.value[i])
    }
  }
  chatList.value = items
}

// 导出会话
const shareChat = (chat) => {
  if (!chat.chat_id) {
    return ElMessage.error('请先选中一个会话')
  }

  const url = location.protocol + '//' + location.host + '/chat/export?chat_id=' + chat.chat_id
  window.open(url, '_blank')
}

const getModelValue = (model_id) => {
  for (let i = 0; i < models.value.length; i++) {
    if (models.value[i].id === model_id) {
      return models.value[i].value
    }
  }
  return ''
}

const files = ref([])
// 插入文件
const insertFile = (file) => {
  files.value.push(file)
}
const removeFile = (file) => {
  files.value = removeArrayItem(files.value, file, (v1, v2) => v1.url === v2.url)
}

// 处理输入框粘贴图片自动上传
const onPaste = async (e) => {
  const items = e.clipboardData?.items || []
  const images = Array.from(items).filter((it) => it.type && it.type.startsWith('image/'))
  if (images.length === 0) {
    return
  }

  // 阻止将图片 base64 粘贴进文本域
  e.preventDefault()

  showLoading('图片上传中...')
  try {
    for (const it of images) {
      const file = it.getAsFile()
      if (!file) continue
      // 10MB 限制
      if (file.size > 10 * 1024 * 1024) {
        ElMessage.error('图片超过10MB限制')
        continue
      }

      const formData = new FormData()
      formData.append('file', file, file.name || 'clipboard.png')
      const res = await httpPost('/api/upload', formData)
      const f = res.data || {}
      if (f.url && f.url.indexOf('http') === -1) {
        f.url = location.protocol + '//' + location.host + f.url
      }
      files.value.push(f)
    }
    ElMessage.success('粘贴图片已上传')
  } catch (err) {
    ElMessage.error('图片上传失败：' + (err?.message || ''))
  } finally {
    closeLoading()
  }
}

// 实时语音对话
const showConversationDialog = ref(false)
// const conversationRef = ref(null);
// const dialogHeight = ref(window.innerHeight - 75);
const frameLoaded = ref(false)
const realtimeChat = () => {
  if (!isLogin.value) {
    store.setShowLoginDialog(true)
    return
  }
  showLoading('正在连接...')
  httpPost('/api/realtime/voice')
    .then((res) => {
      voiceChatUrl.value = res.data
      showConversationDialog.value = true
      closeLoading()
    })
    .catch((e) => {
      showMessageError('连接失败：' + e.message)
      closeLoading()
    })
}

// const hangUp = () => {
//   showConversationDialog.value = false;
//   conversationRef.value.hangUp();
// };
</script>

<style scoped lang="scss">
@use '@/assets/css/chat-plus.scss' as *;
</style>

<style lang="scss">
@use '@/assets/css/markdown/vue.css' as *;
@use 'sass:color';

.input-container {
  .el-textarea {
    .el-textarea__inner {
      padding-right: 40px;
    }
  }
}

.model-selector-popover {
  max-width: 820px !important;
}

.el-popper.model-selector-popover {
  left: 50% !important;
  transform: translateX(-50%) !important;
}

.model-selector-container {
  padding: 16px;

  .model-search {
    margin-bottom: 15px;
    display: flex;
    align-items: center;
  }

  .category-tabs {
    display: flex;
    flex-wrap: wrap;
    border-bottom: 1px solid #e4e7ed;
    margin-bottom: 16px;

    .category-tab {
      padding: 8px 16px;
      cursor: pointer;
      margin-right: 8px;
      margin-bottom: -1px;
      font-size: 14px;
      color: #606266;
      transition: all 0.2s;
      border-bottom: 2px solid transparent;

      &:hover {
        color: #409eff;
      }

      &.active {
        color: #409eff;
        border-bottom-color: #409eff;
        font-weight: 500;
      }

      &.reset-filter {
        color: #f56c6c;
        margin-left: auto;

        &:hover {
          color: color.adjust(#f56c6c, $lightness: -10%);
        }
      }
    }
  }

  .no-results {
    padding: 30px;
    text-align: center;
  }

  .models-grid {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 16px;
    max-height: 450px;
    overflow-y: auto;
    padding: 4px 4px 16px 4px;
  }

  .model-card {
    border: 1px solid #dcdfe6;
    border-radius: 6px;
    padding: 14px;
    cursor: pointer;
    transition: all 0.25s ease;
    height: 100%;
    display: flex;
    flex-direction: column;
    min-width: 0; /* 防止内容溢出 */

    &:hover {
      border-color: #409eff;
      transform: translateY(-2px);
      box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
    }

    &.selected {
      border-color: #409eff;
      background-color: #ecf5ff;
    }

    .model-card-header {
      display: flex;
      justify-content: space-between;
      align-items: flex-start;
      margin-bottom: 8px;

      .model-name {
        font-weight: bold;
        word-break: break-word;
        display: -webkit-box;
        -webkit-line-clamp: 3;
        -webkit-box-orient: vertical;
        overflow: hidden;
        line-height: 1.3;
        max-width: 170px;
        margin-right: 8px;
      }
    }

    .model-description {
      font-size: 12px;
      color: #606266;
      margin-bottom: 10px;
      display: -webkit-box;
      -webkit-line-clamp: 3;
      -webkit-box-orient: vertical;
      overflow: hidden;
      text-overflow: ellipsis;
      line-height: 1.4;
      flex-grow: 1;
    }

    //.model-metadata {
    //  display: flex;
    //  flex-direction: column;
    //  margin-top: auto;
    //}

    .model-detail {
      display: flex;
      justify-content: space-between;
      font-size: 12px;
      color: #909399;
    }
  }
}

.adaptive-width-button {
  min-width: 180px;
  max-width: 350px;
  width: auto !important;
  padding-left: 15px;
  padding-right: 15px;
}

.selected-model-display {
  display: flex;
  align-items: center;
  justify-content: center;

  .model-name-text {
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    max-width: 280px;
  }
}

.customer-service-content {
  text-align: center;
  padding: 10px 0;

  .service-tip {
    font-size: 16px;
    color: #303133;
    margin-bottom: 15px;
  }

  .qrcode-image {
    width: 200px;
    height: 200px;
    margin: 0 auto;
  }

  .service-note {
    font-size: 14px;
    color: #909399;
    margin-top: 15px;
  }
}

.customer-service-btn {
  margin-left: 8px;
}
</style>
