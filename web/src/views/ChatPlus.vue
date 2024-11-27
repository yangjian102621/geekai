<template>
  <div class="chat-page">
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
                   @click="loadChat(chat)">
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
        <div class="chat-container">
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

            <span class="setting" @click="showChatSetting = true">
                    <el-tooltip class="box-item" effect="dark" content="对话设置">
                      <i class="iconfont icon-config"></i>
                    </el-tooltip>
                  </span>
          </div>

          <div>
            <div id="container" :style="{height: mainWinHeight+'px'}">
              <div class="chat-box" id="chat-box" :style="{height: chatBoxHeight+'px'}">
                <div v-if="showHello">
                  <welcome @send="autofillPrompt"/>
                </div>
                <div v-for="item in chatData" :key="item.id" v-else>
                  <chat-prompt
                      v-if="item.type==='prompt'" :data="item" :list-style="listStyle"/>
                  <chat-reply v-else-if="item.type==='reply'" :data="item" @regen="reGenerate" :read-only="false" :list-style="listStyle"/>
                </div>

                <back-top :right="30" :bottom="100" bg-color="#19C27D"/>
              </div><!-- end chat box -->

              <div class="input-box">
                <div class="input-box-inner">
                  <span class="tool-item" @click="ElMessage.info('暂时不支持语音输入')">
                    <el-tooltip class="box-item" effect="dark" content="语音输入">
                      <i class="iconfont icon-mic-bold"></i>
                    </el-tooltip>
                  </span>

                  <span class="tool-item" v-if="isLogin">
                    <el-tooltip class="box-item" effect="dark" content="上传附件">
                      <file-select v-if="isLogin" :user-id="loginUser.id" @selected="insertFile"/>
                    </el-tooltip>
                  </span>

                  <div class="input-body">
                    <div ref="textHeightRef" class="hide-div">{{prompt}}</div>
                    <div class="input-border">
                      <div class="input-inner">
                        <div class="file-list" v-if="files.length > 0">
                          <file-list :files="files" @remove-file="removeFile" />
                        </div>
                        <textarea
                            ref="inputRef"
                            class="prompt-input"
                            :rows="row"
                            v-model="prompt"
                            @keydown="onInput"
                            @input="onInput"
                            placeholder="按 Enter 键发送消息，使用 Ctrl + Enter 换行"
                            autofocus>
                      </textarea>
                      </div>

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

                  </div>
                </div><!-- end input box -->
              </div>
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

    <ChatSetting :show="showChatSetting" @hide="showChatSetting = false"/>
  </div>


</template>
<script setup>
import {nextTick, onMounted, onUnmounted, ref, watch} from 'vue'
import ChatPrompt from "@/components/ChatPrompt.vue";
import ChatReply from "@/components/ChatReply.vue";
import {Delete, Edit, More, Plus, Promotion, Search, Share, VideoPause} from '@element-plus/icons-vue'
import 'highlight.js/styles/a11y-dark.css'
import {
  isMobile,
  randString,
  removeArrayItem,
  UUID
} from "@/utils/libs";
import {ElMessage, ElMessageBox} from "element-plus";
import {getSessionId, getUserToken, removeUserToken} from "@/store/session";
import {httpGet, httpPost} from "@/utils/http";
import {useRouter} from "vue-router";
import Clipboard from "clipboard";
import {checkSession, getSystemInfo} from "@/store/cache";
import Welcome from "@/components/Welcome.vue";
import {useSharedStore} from "@/store/sharedata";
import FileSelect from "@/components/FileSelect.vue";
import FileList from "@/components/FileList.vue";
import ChatSetting from "@/components/ChatSetting.vue";
import BackTop from "@/components/BackTop.vue";
import {showMessageError} from "@/utils/dialog";
import hl from "highlight.js";

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
const inputRef = ref(null)
const textHeightRef = ref(null)
const showNotice = ref(false)
const notice = ref("")
const noticeKey = ref("SYSTEM_NOTICE")
const store = useSharedStore();
const row = ref(1)
const showChatSetting = ref(false)
const listStyle = ref(store.chatListStyle)
watch(() => store.chatListStyle, (newValue) => {
  listStyle.value = newValue
});


if (isMobile()) {
  router.replace("/mobile/chat")
}

// 获取系统配置
getSystemInfo().then(res => {
  title.value = res.data.title
}).catch(e => {
  ElMessage.error("获取系统配置失败：" + e.message)
})

const md = require('markdown-it')({
  breaks: true,
  html: true,
  linkify: true,
  typographer: true
});
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
      if (router.currentRoute.value.query.role_id) {
        roleId.value = parseInt(router.currentRoute.value.query.role_id)
      }
      // 加载模型
      httpGet('/api/model/list').then(res => {
        models.value = res.data
        modelID.value = models.value[0].id
        // 加载角色列表
        httpGet(`/api/role/list`,{id:roleId.value}).then((res) => {
          roles.value = res.data;
          if (!roleId.value) {
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
    // 加载模型
    httpGet('/api/model/list',{id:roleId.value}).then(res => {
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

  inputRef.value.addEventListener('paste', (event) => {
    const items = (event.clipboardData || window.clipboardData).items;
    let fileFound = false;

    for (let item of items) {
      if (item.kind === 'file') {
        const file = item.getAsFile();
        fileFound = true;

        const formData = new FormData();
        formData.append('file', file);
        loading.value = true
        // 执行上传操作
        httpPost('/api/upload', formData).then((res) => {
          files.value.push(res.data)
          ElMessage.success({message: "上传成功", duration: 500})
          loading.value = false
        }).catch((e) => {
          ElMessage.error('文件上传失败:' + e.message)
          loading.value = false
        })

        break;
      }
    }
    
    if (!fileFound) {
      document.getElementById('status').innerText = 'No file found in paste data.';
    }
  });
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
  chatBoxHeight.value = window.innerHeight - 101 - 82 - 38;
  mainWinHeight.value = window.innerHeight - 101;
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

// 创建 socket 连接
const prompt = ref('');
const showStopGenerate = ref(false); // 停止生成
const lineBuffer = ref(''); // 输出缓冲行
const socket = ref(null);
const activelyClose = ref(false); // 主动关闭
const canSend = ref(true);
const heartbeatHandle = ref(null)
const sessionId = ref("")
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
  sessionId.value = getSessionId();
  let host = process.env.VUE_APP_WS_HOST
  if (host === '') {
    if (location.protocol === 'https:') {
      host = 'wss://' + location.host;
    } else {
      host = 'ws://' + location.host;
    }
  }

  const _socket = new WebSocket(host + `/api/chat/new?session_id=${sessionId.value}&role_id=${role_id}&chat_id=${chat_id}&model_id=${modelID.value}&token=${getUserToken()}`);
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
      })
      ElMessage.success({message: "对话连接成功！", duration: 1000})
    } else { // 加载聊天记录
      loadChatHistory(chat_id);
    }
  });

  _socket.addEventListener('message', event => {
    try {
      if (event.data instanceof Blob) {
        const reader = new FileReader();
        reader.readAsText(event.data, "UTF-8");
        reader.onload = () => {
          const data = JSON.parse(String(reader.result));
          if (data.type === 'start') {
            const prePrompt = chatData.value[chatData.value.length-1]?.content
            chatData.value.push({
              type: "reply",
              id: randString(32),
              icon: _role['icon'],
              prompt:prePrompt,
              content: "",
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
            if (reply) {
              reply['content'] = lineBuffer.value;
            }
          }
          // 将聊天框的滚动条滑动到最底部
          nextTick(() => {
            document.getElementById('chat-box').scrollTo(0, document.getElementById('chat-box').scrollHeight)
            localStorage.setItem("chat_id", chat_id)
          })
        };
      }
    } catch (e) {
      console.warn(e)
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
      showMessageError("会话已断开，刷新页面...")
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

const onInput = (e) => {
  // 根据输入的内容自动计算输入框的行数
  const lineHeight = parseFloat(window.getComputedStyle(inputRef.value).lineHeight)
  textHeightRef.value.style.width = inputRef.value.clientWidth + 'px'; // 设定宽度和 textarea 相同
  const lines = Math.floor(textHeightRef.value.clientHeight / lineHeight);
  inputRef.value.scrollTo(0, inputRef.value.scrollHeight)
  if (prompt.value.length < 10) {
    row.value = 1
  } else if (lines <= 7){
    row.value = lines
  } else {
    row.value = 7
  }

  // 输入回车自动提交
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
  inputRef.value.focus()
  sendMessage()
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
  // 如果携带了文件，则串上文件地址
  let content = prompt.value
  if (files.value.length > 0) {
    content += files.value.map(file => file.url).join(" ")
  }
  // 追加消息
  chatData.value.push({
    type: "prompt",
    id: randString(32),
    icon: loginUser.value.avatar,
    content: content,
    model: getModelValue(modelID.value),
    created_at: new Date().getTime() / 1000,
  });

  nextTick(() => {
    document.getElementById('chat-box').scrollTo(0, document.getElementById('chat-box').scrollHeight)
  })

  showHello.value = false
  disableInput(false)
  socket.value.send(JSON.stringify({type: "chat", content: content}));
  tmpChatTitle.value = content
  prompt.value = ''
  files.value = []
  row.value = 1
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
      if (data[i].type === 'reply' && i > 0) {
        data[i].prompt = data[i - 1].content
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
  httpGet("/api/chat/stop?session_id=" + sessionId.value).then(() => {
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
    content: text
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

const files = ref([])
// 插入文件
const insertFile = (file) => {
  files.value.push(file)
}
const removeFile = (file) => {
  files.value = removeArrayItem(files.value, file, (v1,v2) => v1.url===v2.url)
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
</style>