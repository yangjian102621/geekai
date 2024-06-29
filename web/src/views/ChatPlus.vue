<template>
  <div class="common-layout">
    <el-container>
      <el-aside>
        <div class="chat-list">
          <el-button @click="newChat" color="#21aa93">
            <el-icon style="margin-right: 5px">
              <Plus/>
            </el-icon>
            新建对话
          </el-button>

          <div class="search-box">
            <el-input v-model="chatName" placeholder="搜索会话" @keyup="searchChat($event)" style=""
                      class="search-input">
              <template #prefix>
                <el-icon class="el-input__icon">
                  <Search/>
                </el-icon>
              </template>
            </el-input>
          </div>

          <div class="content" :style="{height: leftBoxHeight+'px'}">
            <el-row v-for="chat in chatList" :key="chat.chat_id">
              <div :class="chat.chat_id === activeChat.chat_id?'chat-list-item active':'chat-list-item'"
                   @click="changeChat(chat)">
                <el-image :src="chat.icon" class="avatar"/>
                <span class="chat-title-input" v-if="chat.edit">
                  <el-input v-model="tmpChatTitle" size="small" @keydown="titleKeydown($event, chat)"
                            :id="'chat-'+chat.chat_id"
                            @blur="editConfirm(chat)"
                            @click="stopPropagation($event)"
                            placeholder="请输入标题"/>
                </span>
                <span v-else class="chat-title">{{ chat.title }}</span>

                <span class="chat-opt">
                  <el-dropdown trigger="click">
                    <span class="el-dropdown-link" @click="stopPropagation($event)">
                      <el-icon><More/></el-icon>
                    </span>
                    <template #dropdown>
                      <el-dropdown-menu>
                        <el-dropdown-item :icon="Edit" @click="editChatTitle(chat)">重命名</el-dropdown-item>
                        <el-dropdown-item :icon="Delete"
                                          style="--el-text-color-regular: var(--el-color-danger);
                                          --el-dropdown-menuItem-hover-fill:#F8E1DE;
                                          --el-dropdown-menuItem-hover-color: var(--el-color-danger)"
                                          @click="removeChat(chat)">删除</el-dropdown-item>
                        <el-dropdown-item :icon="Share" @click="shareChat(chat)">分享</el-dropdown-item>
                      </el-dropdown-menu>
                    </template>
                  </el-dropdown>
                </span>
              </div>
            </el-row>
          </div>
        </div>

        <div class="tool-box">
          <el-button type="danger" size="small" @click="clearAllChats">
            <i class="iconfont icon-clear"></i> 清空聊天记录
          </el-button>
        </div>
      </el-aside>
      <el-main v-loading="loading" element-loading-background="rgba(122, 122, 122, 0.3)">
        <div class="chat-box" :style="{height: mainWinHeight+'px'}">
          <div>
            <div id="container">
              <div class="chat-box" id="chat-box" :style="{height: chatBoxHeight+'px'}">
                <div v-if="showHello">
                  <welcome @send="autofillPrompt"/>
                </div>
                <div v-for="item in chatData" :key="item.id" v-else>
                  <chat-prompt
                      v-if="item.type==='prompt'"
                      :icon="item.icon"
                      :created-at="dateFormat(item['created_at'])"
                      :tokens="item['tokens']"
                      :model="getModelValue(modelID)"
                      :content="item.content"/>
                  <chat-reply v-else-if="item.type==='reply'" :data="item" @regen="reGenerate" :read-only="false"/>
                </div>
              </div><!-- end chat box -->

              <el-affix position="bottom" :offset="0">
                <div class="input-box">
                  <span class="tool-item">
                      <el-popover
                          :width="300"
                          trigger="click"
                          placement="top-start"
                      >
                        <template #reference>
                          <div>
                            <el-tooltip effect="dark" content="模型选择">
                              <i class="iconfont icon-model"></i>
                            </el-tooltip>
                          </div>
                        </template>

                        <template #default>
                          <div class="chat-config">
                            <el-select v-model="roleId" filterable placeholder="角色" @change="_newChat"
                                       class="role-select"
                                       style="width:150px">
                              <el-option
                                  v-for="item in roles"
                                  :key="item.id"
                                  :label="item.name"
                                  :value="item.id"
                              >
                                <div class="role-option">
                                  <el-image :src="item.icon"></el-image>
                                  <span>{{ item.name }}</span>
                                </div>
                              </el-option>
                            </el-select>

                            <el-select v-model="modelID" filterable placeholder="模型" @change="_newChat"
                                       :disabled="disableModel"
                                       style="width:150px">
                              <el-option
                                  v-for="item in models"
                                  :key="item.id"
                                  :label="item.name"
                                  :value="item.id"
                              >
                                <span>{{ item.name }}</span>
                                <el-tag style="margin-left: 5px; position: relative; top:-2px" type="info" size="small">{{
                                    item.power
                                  }}算力
                                </el-tag>
                              </el-option>
                            </el-select>
                          </div>
                        </template>
                      </el-popover>
                  </span>

                  <span class="tool-item" @click="ElMessage.info('暂时不支持语音输入')">
                    <el-tooltip class="box-item" effect="dark" content="语音输入">
                      <i class="iconfont icon-mic-bold"></i>
                    </el-tooltip>
                  </span>

                  <span class="tool-item" v-if="isLogin">
                    <el-tooltip class="box-item" effect="dark" content="上传附件">
                      <file-select v-if="isLogin" :user-id="loginUser.id" @selected="insertURL"/>
                    </el-tooltip>
                  </span>

                  <div class="input-container">
                    <el-input
                        ref="textInput"
                        v-model="prompt"
                        v-on:keydown="inputKeyDown"
                        autofocus
                        type="textarea"
                        :rows="2"
                        style="--el-input-focus-border-color:#21AA93;
                        border: 1px solid #21AA93;--el-input-border-color:#21AA93;
                        border-radius: 5px; --el-input-hover-border-color:#21AA93;"
                        placeholder="按 Enter 键发送消息，使用 Ctrl + Enter 换行"
                    />
                    <span class="send-btn">
                    <el-button type="info" v-if="showStopGenerate" @click="stopGenerate" plain>
                      <el-icon>
                        <VideoPause/>
                      </el-icon>
                    </el-button>
                    <el-button @click="sendMessage" color="#19c37d" style="color:#ffffff" v-else>
                      <el-icon><Promotion/></el-icon>
                    </el-button>
                  </span>
                  </div>
                </div><!-- end input box -->
              </el-affix>
            </div><!-- end container -->
          </div><!-- end loading -->
        </div>
      </el-main>
    </el-container>

    <el-dialog
        v-model="showNotice"
        :show-close="true"
        class="notice-dialog"
        title="网站公告"
    >
      <div class="notice">
        <div v-html="notice"></div>

        <p style="text-align: right">
          <el-button @click="notShow" type="success" plain>我知道了，不再显示</el-button>
        </p>
      </div>
    </el-dialog>
  </div>


</template>
<script setup>
import {nextTick, onMounted, onUnmounted, ref} from 'vue'
import ChatPrompt from "@/components/ChatPrompt.vue";
import ChatReply from "@/components/ChatReply.vue";
import {Delete, Edit, More, Plus, Promotion, Search, Share, VideoPause} from '@element-plus/icons-vue'
import 'highlight.js/styles/a11y-dark.css'
import {dateFormat, escapeHTML, isMobile, processContent, randString, removeArrayItem, UUID} from "@/utils/libs";
import {ElMessage, ElMessageBox} from "element-plus";
import hl from "highlight.js";
import {getSessionId, getUserToken, removeUserToken} from "@/store/session";
import {httpGet, httpPost} from "@/utils/http";
import {useRouter} from "vue-router";
import Clipboard from "clipboard";
import {checkSession} from "@/action/session";
import Welcome from "@/components/Welcome.vue";
import {useSharedStore} from "@/store/sharedata";
import FileSelect from "@/components/FileSelect.vue";

const title = ref('ChatGPT-智能助手');
const models = ref([])
const modelID = ref(0)
const chatData = ref([]);
const allChats = ref([]); // 会话列表
const chatList = ref(allChats.value);
const activeChat = ref({});
const mainWinHeight = ref(0); // 主窗口高度
const chatBoxHeight = ref(0); // 聊天内容框高度
const leftBoxHeight = ref(0);
const loading = ref(true);
const loginUser = ref(null);
const roles = ref([]);
const router = useRouter();
const roleId = ref(0)
const newChatItem = ref(null);
const isLogin = ref(false)
const showHello = ref(true)
const textInput = ref(null)
const showNotice = ref(false)
const notice = ref("")
const noticeKey = ref("SYSTEM_NOTICE")
const store = useSharedStore();

if (isMobile()) {
  router.replace("/mobile/chat")
}

// 获取系统配置
httpGet("/api/config/get?key=system").then(res => {
  title.value = res.data.title
}).catch(e => {
  ElMessage.error("获取系统配置失败：" + e.message)
})

// 获取系统公告
httpGet("/api/config/get?key=notice").then(res => {
  try {
    notice.value = md.render(res.data['content'])
    const oldNotice = localStorage.getItem(noticeKey.value);
    // 如果公告有更新，则显示公告
    if (oldNotice !== notice.value && notice.value.length > 10) {
      showNotice.value = true
    }
  } catch (e) {
    console.warn(e)
  }

}).catch(e => {
  ElMessage.error("获取系统配置失败：" + e.message)
})

onMounted(() => {
  resizeElement();
  initData()

  const clipboard = new Clipboard('.copy-reply, .copy-code-btn');
  clipboard.on('success', () => {
    ElMessage.success('复制成功！');
  })

  clipboard.on('error', () => {
    ElMessage.error('复制失败！');
  })

  window.onresize = () => resizeElement();
});

onUnmounted(() => {
  if (socket.value !== null) {
    socket.value.close()
    socket.value = null
  }
})

// 初始化数据
const initData = () => {
  // 检查会话
  checkSession().then((user) => {
    loginUser.value = user
    isLogin.value = true

    // 获取会话列表
    httpGet("/api/chat/list").then((res) => {
      if (res.data) {
        chatList.value = res.data;
        allChats.value = res.data;
      }

      // 加载模型
      httpGet('/api/model/list').then(res => {
        models.value = res.data
        modelID.value = models.value[0].id

        // 加载角色列表
        httpGet(`/api/role/list`).then((res) => {
          roles.value = res.data;
          if (router.currentRoute.value.query.role_id) {
            roleId.value = parseInt(router.currentRoute.value.query.role_id)
          } else {
            roleId.value = roles.value[0]['id']
          }

          newChat();
        }).catch((e) => {
          ElMessage.error('获取聊天角色失败: ' + e.messages)
        })
      }).catch(e => {
        ElMessage.error("加载模型失败: " + e.message)
      })
    }).catch(() => {
      ElMessage.error("加载会话列表失败！")
    })
  }).catch(() => {
    loading.value = false
    // 加载会话
    httpGet("/api/chat/list").then((res) => {
      if (res.data) {
        chatList.value = res.data;
        allChats.value = res.data;
      }
    }).catch(() => {
      ElMessage.error("加载会话列表失败！")
    })

    // 加载模型
    httpGet('/api/model/list').then(res => {
      models.value = res.data
      modelID.value = models.value[0].id
    }).catch(e => {
      ElMessage.error("加载模型失败: " + e.message)
    })

    // 加载角色列表
    httpGet(`/api/role/list`).then((res) => {
      roles.value = res.data;
      roleId.value = roles.value[0]['id'];
    }).catch((e) => {
      ElMessage.error('获取聊天角色失败: ' + e.messages)
    })
  })
}

const getRoleById = function (rid) {
  for (let i = 0; i < roles.value.length; i++) {
    if (roles.value[i]['id'] === rid) {
      return roles.value[i];
    }
  }
  return null;
}

const resizeElement = function () {
  chatBoxHeight.value = window.innerHeight - 50 - 82 - 38;
  mainWinHeight.value = window.innerHeight - 50;
  leftBoxHeight.value = window.innerHeight - 90 - 45 - 82;
};

const _newChat = () => {
  if (isLogin.value) {
    newChat()
  }
}
const disableModel = ref(false)
// 新建会话
const newChat = () => {
  if (!isLogin.value) {
    store.setShowLoginDialog(true)
    return;
  }
  const role = getRoleById(roleId.value)
  showHello.value = role.key === 'gpt';
  // if the role bind a model, disable model change
  disableModel.value = false
  if (role.model_id > 0) {
    modelID.value = role.model_id
    disableModel.value = true
  }
  // 已有新开的会话
  if (newChatItem.value !== null && newChatItem.value['role_id'] === roles.value[0]['role_id']) {
    return;
  }

  // 获取当前聊天角色图标
  let icon = '';
  roles.value.forEach(item => {
    if (item['id'] === roleId.value) {
      icon = item['icon']
    }
  })
  newChatItem.value = {
    chat_id: "",
    icon: icon,
    role_id: roleId.value,
    model_id: modelID.value,
    title: '',
    edit: false,
    removing: false,
  };
  activeChat.value = {} //取消激活的会话高亮
  showStopGenerate.value = false;
  connect(null, roleId.value)
}

// 切换会话
const changeChat = (chat) => {
  localStorage.setItem("chat_id", chat.chat_id)
  loadChat(chat)
}

const loadChat = function (chat) {
  if (!isLogin.value) {
    store.setShowLoginDialog(true)
    return;
  }

  if (activeChat.value['chat_id'] === chat.chat_id) {
    return;
  }

  activeChat.value = chat
  newChatItem.value = null;
  roleId.value = chat.role_id;
  modelID.value = chat.model_id;
  showStopGenerate.value = false;
  connect(chat.chat_id, chat.role_id)
}

// 编辑会话标题
const tmpChatTitle = ref('');
const editChatTitle = (chat) => {
  chat.edit = true;
  tmpChatTitle.value = chat.title;
  console.log(chat.chat_id)
  nextTick(() => {
    document.getElementById('chat-' + chat.chat_id).focus()
  })
};


const titleKeydown = (e, chat) => {
  if (e.keyCode === 13) {
    e.stopPropagation();
    editConfirm(chat)
  }
}

const stopPropagation = (e) => {
  e.stopPropagation();
}
// 确认修改
const editConfirm = function (chat) {
  if (tmpChatTitle.value === '') {
    return ElMessage.error("请输入会话标题！");
  }
  if (!chat.chat_id) {
    return ElMessage.error("对话 ID 为空，请刷新页面再试！");
  }
  if (tmpChatTitle.value === chat.title) {
    chat.edit = false;
    return
  }

  httpPost('/api/chat/update', {chat_id: chat.chat_id, title: tmpChatTitle.value}).then(() => {
    chat.title = tmpChatTitle.value;
    chat.edit = false;
  }).catch(e => {
    ElMessage.error("操作失败：" + e.message);
  })

}
// 删除会话
const removeChat = function (chat) {
  ElMessageBox.confirm(
      `该操作会删除"${chat.title}"`,
      '删除聊天',
      {
        confirmButtonText: '删除',
        cancelButtonText: '取消',
        type: 'warning',
      }
  )
      .then(() => {
        httpGet('/api/chat/remove?chat_id=' + chat.chat_id).then(() => {
          chatList.value = removeArrayItem(chatList.value, chat, function (e1, e2) {
            return e1.id === e2.id
          })
          // 重置会话
          newChat();
        }).catch(e => {
          ElMessage.error("操作失败：" + e.message);
        })
      })
      .catch(() => {
      })

}

const latexPlugin = require('markdown-it-latex2img')
const mathjaxPlugin = require('markdown-it-mathjax')
const md = require('markdown-it')({
  breaks: true,
  html: true,
  linkify: true,
  typographer: true,
  highlight: function (str, lang) {
    const codeIndex = parseInt(Date.now()) + Math.floor(Math.random() * 10000000)
    // 显示复制代码按钮
    const copyBtn = `<span class="copy-code-btn" data-clipboard-action="copy" data-clipboard-target="#copy-target-${codeIndex}">复制</span>
<textarea style="position: absolute;top: -9999px;left: -9999px;z-index: -9999;" id="copy-target-${codeIndex}">${str.replace(/<\/textarea>/g, '&lt;/textarea>')}</textarea>`
    if (lang && hl.getLanguage(lang)) {
      const langHtml = `<span class="lang-name">${lang}</span>`
      // 处理代码高亮
      const preCode = hl.highlight(lang, str, true).value
      // 将代码包裹在 pre 中
      return `<pre class="code-container"><code class="language-${lang} hljs">${preCode}</code>${copyBtn} ${langHtml}</pre>`
    }

    // 处理代码高亮
    const preCode = md.utils.escapeHtml(str)
    // 将代码包裹在 pre 中
    return `<pre class="code-container"><code class="language-${lang} hljs">${preCode}</code>${copyBtn}</pre>`
  }
});
md.use(latexPlugin)
md.use(mathjaxPlugin)

// 创建 socket 连接
const prompt = ref('');
const showStopGenerate = ref(false); // 停止生成
const lineBuffer = ref(''); // 输出缓冲行
const socket = ref(null);
const activelyClose = ref(false); // 主动关闭
const canSend = ref(true);
const heartbeatHandle = ref(null)
const connect = function (chat_id, role_id) {
  let isNewChat = false;
  if (!chat_id) {
    isNewChat = true;
    chat_id = UUID();
  }
  // 先关闭已有连接
  if (socket.value !== null) {
    activelyClose.value = true;
    socket.value.close();
  }

  const _role = getRoleById(role_id);
  // 初始化 WebSocket 对象
  const _sessionId = getSessionId();
  let host = process.env.VUE_APP_WS_HOST
  if (host === '') {
    if (location.protocol === 'https:') {
      host = 'wss://' + location.host;
    } else {
      host = 'ws://' + location.host;
    }
  }

  // 心跳函数
  const sendHeartbeat = () => {
    clearTimeout(heartbeatHandle.value)
    new Promise((resolve, reject) => {
      if (socket.value !== null) {
        socket.value.send(JSON.stringify({type: "heartbeat", content: "ping"}))
      }
      resolve("success")
    }).then(() => {
      heartbeatHandle.value = setTimeout(() => sendHeartbeat(), 5000)
    });
  }
  const _socket = new WebSocket(host + `/api/chat/new?session_id=${_sessionId}&role_id=${role_id}&chat_id=${chat_id}&model_id=${modelID.value}&token=${getUserToken()}`);
  _socket.addEventListener('open', () => {
    chatData.value = []; // 初始化聊天数据
    enableInput()
    activelyClose.value = false;

    if (isNewChat) { // 加载打招呼信息
      loading.value = false;
      chatData.value.push({
        chat_id: chat_id,
        role_id: role_id,
        type: "reply",
        id: randString(32),
        icon: _role['icon'],
        content: _role['hello_msg'],
        orgContent: _role['hello_msg'],
      })
      ElMessage.success({message: "对话连接成功！", duration: 1000})
    } else { // 加载聊天记录
      loadChatHistory(chat_id);
    }
    // 发送心跳消息
    sendHeartbeat()
  });

  _socket.addEventListener('message', event => {
    try {
      if (event.data instanceof Blob) {
        const reader = new FileReader();
        reader.readAsText(event.data, "UTF-8");
        reader.onload = () => {
          const data = JSON.parse(String(reader.result));
          if (data.type === 'start') {
            chatData.value.push({
              type: "reply",
              id: randString(32),
              icon: _role['icon'],
              content: ""
            });
          } else if (data.type === 'end') { // 消息接收完毕
            // 追加当前会话到会话列表
            if (isNewChat && newChatItem.value !== null) {
              newChatItem.value['title'] = tmpChatTitle.value;
              newChatItem.value['chat_id'] = chat_id;
              chatList.value.unshift(newChatItem.value);
              activeChat.value = newChatItem.value;
              newChatItem.value = null; // 只追加一次
            }

            enableInput()
            lineBuffer.value = ''; // 清空缓冲

            // 获取 token
            const reply = chatData.value[chatData.value.length - 1]
            httpPost("/api/chat/tokens", {
              text: "",
              model: getModelValue(modelID.value),
              chat_id: chat_id
            }).then(res => {
              reply['created_at'] = new Date().getTime();
              reply['tokens'] = res.data;
              // 将聊天框的滚动条滑动到最底部
              nextTick(() => {
                document.getElementById('chat-box').scrollTo(0, document.getElementById('chat-box').scrollHeight)
              })
            }).catch(() => {
            })

          } else {
            lineBuffer.value += data.content;
            const reply = chatData.value[chatData.value.length - 1]
            reply['orgContent'] = lineBuffer.value;
            reply['content'] = md.render(processContent(lineBuffer.value));
          }
          // 将聊天框的滚动条滑动到最底部
          nextTick(() => {
            document.getElementById('chat-box').scrollTo(0, document.getElementById('chat-box').scrollHeight)
            localStorage.setItem("chat_id", chat_id)
          })
        };
      }
    } catch (e) {
      console.error(e)
    }

  });

  _socket.addEventListener('close', () => {
    if (activelyClose.value || socket.value === null) { // 忽略主动关闭
      return;
    }
    // 停止发送消息
    disableInput(true)
    loading.value = true;
    checkSession().then(() => {
      connect(chat_id, role_id)
    }).catch(() => {
      loading.value = true
      setTimeout(() => connect(chat_id, role_id), 3000)
    });
  });

  socket.value = _socket;
}

const disableInput = (force) => {
  canSend.value = false;
  showStopGenerate.value = !force;
}

const enableInput = () => {
  canSend.value = true;
  showStopGenerate.value = false;
}

// 登录输入框输入事件处理
const inputKeyDown = function (e) {
  if (e.keyCode === 13) {
    if (e.ctrlKey) { // Ctrl + Enter 换行
      prompt.value += "\n";
      return;
    }
    e.preventDefault();
    sendMessage();
  }
}

// 自动填充 prompt
const autofillPrompt = (text) => {
  prompt.value = text
  textInput.value.focus()
  // sendMessage()
}
// 发送消息
const sendMessage = function () {
  if (!isLogin.value) {
    store.setShowLoginDialog(true)
    return;
  }

  if (canSend.value === false) {
    ElMessage.warning("AI 正在作答中，请稍后...");
    return
  }

  if (prompt.value.trim().length === 0 || canSend.value === false) {
    return false;
  }
  // 追加消息
  chatData.value.push({
    type: "prompt",
    id: randString(32),
    icon: loginUser.value.avatar,
    content: md.render(escapeHTML(processContent(prompt.value))),
    created_at: new Date().getTime() / 1000,
  });

  nextTick(() => {
    document.getElementById('chat-box').scrollTo(0, document.getElementById('chat-box').scrollHeight)
  })

  showHello.value = false
  disableInput(false)
  socket.value.send(JSON.stringify({type: "chat", content: prompt.value}));
  tmpChatTitle.value = prompt.value
  prompt.value = '';
  return true;
}

const clearAllChats = function () {
  ElMessageBox.confirm(
      '确认要清空所有的会话历史记录吗?<br/>此操作不可以撤销！',
      '操作提示：',
      {
        confirmButtonText: '确认',
        cancelButtonText: '取消',
        type: 'warning',
        dangerouslyUseHTMLString: true,
        showClose: true,
        closeOnClickModal: false,
        center: true,
      }
  ).then(() => {
    httpGet("/api/chat/clear").then(() => {
      ElMessage.success("操作成功！");
      chatData.value = [];
      chatList.value = [];
      newChat();
    }).catch(e => {
      ElMessage.error("操作失败：" + e.message)
    })
  }).catch(() => {
  })
}

const logout = function () {
  activelyClose.value = true;
  httpGet('/api/user/logout').then(() => {
    removeUserToken()
    router.push("/login")
  }).catch(() => {
    ElMessage.error('注销失败！');
  })
}

const loadChatHistory = function (chatId) {
  httpGet('/api/chat/history?chat_id=' + chatId).then(res => {
    const data = res.data
    if (!data) {
      loading.value = false
      return
    }
    showHello.value = false
    for (let i = 0; i < data.length; i++) {
      data[i].orgContent = data[i].content;
      data[i].content = md.render(processContent(data[i].content))
      if (i > 0 && data[i].type === 'reply') {
        data[i].prompt = data[i - 1].orgContent
      }
      chatData.value.push(data[i]);
    }

    nextTick(() => {
      document.getElementById('chat-box').scrollTo(0, document.getElementById('chat-box').scrollHeight)
    })
    loading.value = false
  }).catch(e => {
    // TODO: 显示重新加载按钮
    ElMessage.error('加载聊天记录失败：' + e.message);
  })
}

const stopGenerate = function () {
  showStopGenerate.value = false;
  httpGet("/api/chat/stop?session_id=" + getSessionId()).then(() => {
    enableInput()
  })
}

// 重新生成
const reGenerate = function (prompt) {
  disableInput(false)
  const text = '重新生成下面问题的答案：' + prompt;
  // 追加消息
  chatData.value.push({
    type: "prompt",
    id: randString(32),
    icon: loginUser.value.avatar,
    content: md.render(text)
  });
  socket.value.send(JSON.stringify({type: "chat", content: prompt}));
}

const chatName = ref('')
// 搜索会话
const searchChat = function (e) {
  if (chatName.value === '') {
    chatList.value = allChats.value
    return
  }
  if (e.keyCode === 13) {
    const items = [];
    for (let i = 0; i < allChats.value.length; i++) {
      if (allChats.value[i].title.toLowerCase().indexOf(chatName.value.toLowerCase()) !== -1) {
        items.push(allChats.value[i]);
      }
    }
    chatList.value = items;
  }
}

// 导出会话
const shareChat = (chat) => {
  if (!chat.chat_id) {
    return ElMessage.error("请先选中一个会话")
  }

  const url = location.protocol + '//' + location.host + '/chat/export?chat_id=' + chat.chat_id
  // console.log(url)
  window.open(url, '_blank');
}


const getModelValue = (model_id) => {
  for (let i = 0; i < models.value.length; i++) {
    if (models.value[i].id === model_id) {
      return models.value[i].value
    }
  }
  return ""
}


const notShow = () => {
  localStorage.setItem(noticeKey.value, notice.value)
  showNotice.value = false
}

// 插入文件路径
const insertURL = (url) => {
  prompt.value += " " + url + " "
}
</script>

<style scoped lang="stylus">
@import "@/assets/css/chat-plus.styl"
</style>

<style lang="stylus">
.notice-dialog {
  .el-dialog__header {
    padding-bottom 0
  }

  .el-dialog__body {
    padding 0 20px

    ol, ul {
      padding-left 10px
    }

    ol {
      list-style decimal-leading-zero
      padding-left 20px
    }

    ul {
      list-style disc
    }
  }
}

.input-container {
  .el-textarea {
    .el-textarea__inner {
      padding-right 40px
    }
  }
}

.chat-config {
  display flex
  flex-direction row
  padding-top 10px;

  .role-select-label {
    color #ffffff
  }

  .el-select {
    max-width 150px;
    margin-right 10px;
  }

  .role-select {
    max-width 130px;
  }

  .el-button {
    .el-icon {
      margin-right 5px;
    }
  }
}
</style>