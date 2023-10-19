<template>
  <div class="common-layout theme-white">
    <el-container>
      <el-aside>
        <div class="title-box">
          <span>{{ title }}</span>
        </div>
        <div class="chat-list">
          <div class="search-box">
            <el-input v-model="chatName" class="w-50 m-2" size="small" placeholder="搜索会话" @keyup="searchChat">
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
                        placeholder="请输入会话标题"/>
            </span>
                <span v-else class="chat-title">{{ chat.title }}</span>
                <span class="btn btn-check" v-if="chat.edit || chat.removing">
                <el-icon @click="confirm($event, chat)"><Check/></el-icon>
                <el-icon @click="cancel($event, chat)"><Close/></el-icon>
              </span>
                <span class="btn" v-else>
                <el-icon title="编辑" @click="editChatTitle($event, chat)"><Edit/></el-icon>
                <el-icon title="删除会话" @click="removeChat($event, chat)"><Delete/></el-icon>
              </span>
              </div>
            </el-row>
          </div>
        </div>

        <div class="tool-box">
          <el-dropdown :hide-on-click="true" class="user-info" trigger="click" v-if="isLogin">
                        <span class="el-dropdown-link">
                          <el-image :src="loginUser.avatar"/>
                          <span class="username">{{ '极客学长@' + loginUser.mobile }}</span>
                          <el-icon><ArrowDown/></el-icon>
                        </span>
            <template #dropdown>
              <el-dropdown-menu style="width: 296px;">
                <el-dropdown-item @click="showConfig">
                  <el-icon>
                    <Tools/>
                  </el-icon>
                  <span>聊天设置</span>
                </el-dropdown-item>

                <el-dropdown-item @click="showPasswordDialog=true">
                  <i class="iconfont icon-password"></i>
                  <span>修改密码</span>
                </el-dropdown-item>

                <el-dropdown-item @click="showBindMobileDialog = true">
                  <el-icon>
                    <Iphone/>
                  </el-icon>
                  <span>绑定手机号</span>
                </el-dropdown-item>

                <el-dropdown-item @click="showRewardDialog = true" v-if="enableReward">
                  <el-icon>
                    <Present/>
                  </el-icon>
                  <span>加入众筹</span>
                </el-dropdown-item>

                <el-dropdown-item @click="showRewardVerifyDialog = true" v-if="enableReward">
                  <el-icon>
                    <Checked/>
                  </el-icon>
                  <span>众筹核销</span>
                </el-dropdown-item>

                <el-dropdown-item @click="clearAllChats">
                  <el-icon>
                    <Delete/>
                  </el-icon>
                  <span>清除所有会话</span>
                </el-dropdown-item>

                <el-dropdown-item @click="logout">
                  <i class="iconfont icon-logout"></i>
                  <span>注销</span>
                </el-dropdown-item>

                <el-dropdown-item>
                  <i class="iconfont icon-github"></i>
                  <span>
                    powered by
                    <el-link type="primary" href="https://github.com/yangjian102621/chatgpt-plus" target="_blank">chatgpt-plus-v3</el-link>
                 </span>
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-aside>
      <el-main v-loading="loading" element-loading-background="rgba(122, 122, 122, 0.3)">
        <div class="chat-head">
          <div class="chat-config">
            <span class="role-select-label">聊天角色：</span>
            <el-select v-model="roleId" filterable placeholder="角色" class="role-select">
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

            <el-select v-model="modelID" placeholder="模型">
              <el-option
                  v-for="item in models"
                  :key="item.id"
                  :label="item.name"
                  :value="item.id"
              />
            </el-select>
            <el-button type="primary" @click="newChat">
              <el-icon>
                <Plus/>
              </el-icon>
              新建会话
            </el-button>

            <el-button type="success" @click="exportChat" plain>
              <i class="iconfont icon-export"></i>
              <span>导出会话</span>
            </el-button>
          </div>
        </div>

        <div class="right-box" :style="{height: mainWinHeight+'px'}">
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
                  <chat-reply v-else-if="item.type==='reply'"
                              :icon="item.icon"
                              :org-content="item.orgContent"
                              :created-at="dateFormat(item['created_at'])"
                              :tokens="item['tokens']"
                              :content="item.content"/>
                  <chat-mid-journey v-else-if="item.type==='mj'"
                                    :content="item.content"
                                    :role-id="item.role_id"
                                    :chat-id="item.chat_id"
                                    :icon="item.icon"
                                    @disable-input="disableInput(true)"
                                    @enable-input="enableInput"
                                    :created-at="dateFormat(item['created_at'])"/>
                </div>
              </div><!-- end chat box -->

              <div class="re-generate">
                <div class="btn-box">
                  <el-button type="info" v-if="showStopGenerate" @click="stopGenerate" plain>
                    <el-icon>
                      <VideoPause/>
                    </el-icon>
                    停止生成
                  </el-button>

                  <el-button type="primary" v-if="showReGenerate" @click="reGenerate" plain>
                    <el-icon>
                      <RefreshRight/>
                    </el-icon>
                    重新生成
                  </el-button>
                </div>
              </div>

              <div class="input-box">
                <div class="input-container">
                  <el-input
                      ref="textInput"
                      v-model="prompt"
                      v-on:keydown="inputKeyDown"
                      autofocus
                      type="textarea"
                      :rows="2"
                      placeholder="按 Enter 键发送消息，使用 Ctrl + Enter 换行"
                  />
                  <span class="send-btn">
                    <el-button @click="sendMessage">
                      <el-icon><Promotion/></el-icon>
                    </el-button>
                  </span>
                </div>
              </div><!-- end input box -->

            </div><!-- end container -->
          </div><!-- end loading -->
        </div>
      </el-main>
    </el-container>

    <config-dialog v-if="isLogin" :show="showConfigDialog" :models="models" @hide="showConfigDialog = false"
                   @update-user="updateUser"/>
    <password-dialog v-if="isLogin" :show="showPasswordDialog" @hide="showPasswordDialog = false"
                     @logout="logout"/>

    <bind-mobile v-if="isLogin" :show="showBindMobileDialog" :mobile="loginUser.mobile"
                 @hide="showBindMobileDialog = false"/>

    <reward-verify v-if="isLogin" :show="showRewardVerifyDialog" @hide="showRewardVerifyDialog = false"/>

    <el-dialog
        v-model="showRewardDialog"
        :show-close="true"
        width="400px"
        title="参与众筹"
    >
      <el-alert type="info" :closable="false">
        <div style="font-size: 14px">您好，众筹 9.9元，就可以兑换 100 次对话，以此来覆盖我们的 OpenAI
          账单和服务器的费用。<strong
              style="color: #f56c6c">由于本人没有开通微信支付，付款后请凭借转账单号进入核销【众筹核销】菜单手动核销。</strong>
        </div>
      </el-alert>
      <div style="text-align: center;padding-top: 10px;">
        <el-image v-if="enableReward" :src="rewardImg"/>
      </div>
    </el-dialog>
  </div>


</template>
<script setup>
import {nextTick, onMounted, ref} from 'vue'
import ChatPrompt from "@/components/ChatPrompt.vue";
import ChatReply from "@/components/ChatReply.vue";
import {
  ArrowDown,
  Check,
  Checked,
  Close,
  Delete,
  Edit,
  Iphone,
  Plus,
  Present,
  Promotion,
  RefreshRight,
  Search,
  Tools,
  VideoPause
} from '@element-plus/icons-vue'
import 'highlight.js/styles/a11y-dark.css'
import {dateFormat, isMobile, randString, removeArrayItem, renderInputText, UUID} from "@/utils/libs";
import {ElMessage, ElMessageBox} from "element-plus";
import hl from "highlight.js";
import {getSessionId, getUserToken, removeUserToken} from "@/store/session";
import {httpGet, httpPost} from "@/utils/http";
import {useRouter} from "vue-router";
import Clipboard from "clipboard";
import ConfigDialog from "@/components/ConfigDialog.vue";
import PasswordDialog from "@/components/PasswordDialog.vue";
import {checkSession} from "@/action/session";
import BindMobile from "@/components/BindMobile.vue";
import RewardVerify from "@/components/RewardVerify.vue";
import Welcome from "@/components/Welcome.vue";
import ChatMidJourney from "@/components/ChatMidJourney.vue";

const title = ref('ChatGPT-智能助手');
const enableReward = ref(false) // 是否启用众筹功能
const rewardImg = ref('/images/reward.png')
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
const roleId = ref(0)
const newChatItem = ref(null);
const router = useRouter();
const showConfigDialog = ref(false);
const showPasswordDialog = ref(false);
const showBindMobileDialog = ref(false);
const showRewardDialog = ref(false);
const showRewardVerifyDialog = ref(false);
const isLogin = ref(false)
const showHello = ref(true)
const textInput = ref(null)

if (isMobile()) {
  router.replace("/mobile")
}

onMounted(() => {
  resizeElement();
  checkSession().then((user) => {
    loginUser.value = user
    isLogin.value = true
    // 获取会话列表
    httpGet("/api/chat/list?user_id=" + loginUser.value.id).then((res) => {
      if (res.data) {
        chatList.value = res.data;
        allChats.value = res.data;
      }
      // 加载模型
      httpGet('/api/model/list?enable=1').then(res => {
        models.value = res.data
        modelID.value = models.value[0].id

        // 加载角色列表
        httpGet(`/api/role/list?user_id=${user.id}`).then((res) => {
          roles.value = res.data;
          roleId.value = roles.value[0]['id'];
          const chatId = localStorage.getItem("chat_id")
          const chat = getChatById(chatId)
          if (chat === null) {
            // 创建新的对话
            newChat();
          } else {
            // 加载对话
            loadChat(chat)
          }
        }).catch((e) => {
          ElMessage.error('获取聊天角色失败: ' + e.messages)
        })
      }).catch(e => {
        ElMessage.error("加载模型失败: " + e.message)
      })

    }).catch(() => {
      // TODO: 增加重试按钮
      ElMessage.error("加载会话列表失败！")
    })

    httpGet("/api/admin/config/get?key=system").then(res => {
      title.value = res.data.title
      rewardImg.value = res.data.reward_img
      enableReward.value = res.data.enabled_reward
    }).catch(e => {
      ElMessage.error("获取系统配置失败：" + e.message)
    })
  }).catch(() => {
    router.push('/login')
  });

  const clipboard = new Clipboard('.copy-reply');
  clipboard.on('success', () => {
    ElMessage.success('复制成功！');
  })

  clipboard.on('error', () => {
    ElMessage.error('复制失败！');
  })

  window.onresize = () => resizeElement();
});

const getRoleById = function (rid) {
  for (let i = 0; i < roles.value.length; i++) {
    if (roles.value[i]['id'] === rid) {
      return roles.value[i];
    }
  }
  return null;
}

const resizeElement = function () {
  chatBoxHeight.value = window.innerHeight - 51 - 82 - 38;
  mainWinHeight.value = window.innerHeight - 51;
  leftBoxHeight.value = window.innerHeight - 43 - 47 - 45;
};

// 新建会话
const newChat = function () {
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
  showReGenerate.value = false;
  connect(null, roleId.value)
}

// 切换会话
const changeChat = (chat) => {
  localStorage.setItem("chat_id", chat.chat_id)
  loadChat(chat)
}

const loadChat = function (chat) {
  if (activeChat.value['chat_id'] === chat.chat_id) {
    return;
  }
  activeChat.value = chat
  newChatItem.value = null;
  roleId.value = chat.role_id;
  modelID.value = chat.model_id;
  showStopGenerate.value = false;
  showReGenerate.value = false;
  connect(chat.chat_id, chat.role_id)
}

// 编辑会话标题
const curOpt = ref('')
const tmpChatTitle = ref('');
const editChatTitle = function (event, chat) {
  event.stopPropagation();
  chat.edit = true;
  curOpt.value = 'edit';
  tmpChatTitle.value = chat.title;
};


const titleKeydown = (e, chat) => {
  if (e.keyCode === 13) {
    e.stopPropagation();
    confirm(e, chat)
  }
}
// 确认修改
const confirm = function (event, chat) {
  event.stopPropagation();
  if (curOpt.value === 'edit') {
    if (tmpChatTitle.value === '') {
      return ElMessage.error("请输入会话标题！");
    }
    if (!chat.chat_id) {
      return ElMessage.error("对话 ID 为空，请刷新页面再试！");
    }
    httpPost('/api/chat/update', {chat_id: chat.chat_id, title: tmpChatTitle.value}).then(() => {
      chat.title = tmpChatTitle.value;
      chat.edit = false;
    }).catch(e => {
      ElMessage.error("操作失败：" + e.message);
    })
  } else if (curOpt.value === 'remove') {
    httpGet('/api/chat/remove?chat_id=' + chat.chat_id).then(() => {
      chatList.value = removeArrayItem(chatList.value, chat, function (e1, e2) {
        return e1.id === e2.id
      })
      // 重置会话
      newChat();
    }).catch(e => {
      ElMessage.error("操作失败：" + e.message);
    })

  }

}
// 取消修改
const cancel = function (event, chat) {
  event.stopPropagation();
  chat.edit = false;
  chat.removing = false;
}

// 删除会话
const removeChat = function (event, chat) {
  event.stopPropagation();
  chat.removing = true;
  curOpt.value = 'remove';
}

// 创建 socket 连接
const prompt = ref('');
const showStopGenerate = ref(false); // 停止生成
const showReGenerate = ref(false); // 重新生成
const previousText = ref(''); // 上一次提问
const lineBuffer = ref(''); // 输出缓冲行
const socket = ref(null);
const activelyClose = ref(false); // 主动关闭
const canSend = ref(true);
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
  const _socket = new WebSocket(host + `/api/chat/new?session_id=${_sessionId}&role_id=${role_id}&chat_id=${chat_id}&model_id=${modelID.value}&token=${getUserToken()}`);
  _socket.addEventListener('open', () => {
    chatData.value = []; // 初始化聊天数据
    previousText.value = '';
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

  });

  _socket.addEventListener('message', event => {
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
        } else if (data.type === "mj") {
          disableInput(true)
          const content = data.content;
          const md = require('markdown-it')({breaks: true});
          content.html = md.render(content.content)
          let key = content.key
          // fixed bug: 执行 Upscale 和 Variation 操作的时候覆盖之前的绘画
          if (content.status === "Finished") {
            key = randString(32)
            enableInput()
          }
          // console.log(content)
          // check if the message is in chatData
          let flag = false
          for (let i = 0; i < chatData.value.length; i++) {
            if (chatData.value[i].id === content.key) {
              flag = true
              chatData.value[i].content = content
              chatData.value[i].id = key
              break
            }
          }
          if (flag === false) {
            chatData.value.push({
              type: "mj",
              id: key,
              icon: "/images/avatar/mid_journey.png",
              content: content
            });
          }

        } else if (data.type === 'end') { // 消息接收完毕
          // 追加当前会话到会话列表
          if (isNewChat && newChatItem.value !== null) {
            newChatItem.value['title'] = previousText.value;
            newChatItem.value['chat_id'] = chat_id;
            chatList.value.unshift(newChatItem.value);
            activeChat.value = newChatItem.value;
            newChatItem.value = null; // 只追加一次
          }

          enableInput()
          lineBuffer.value = ''; // 清空缓冲

          // 获取 token
          const reply = chatData.value[chatData.value.length - 1]
          httpPost("/api/chat/tokens", {text: "", model: getModelValue(modelID.value)}).then(res => {
            reply['created_at'] = new Date().getTime();
            reply['tokens'] = res.data;
            // 将聊天框的滚动条滑动到最底部
            nextTick(() => {
              document.getElementById('chat-box').scrollTo(0, document.getElementById('chat-box').scrollHeight)
            })
          })

        } else {
          lineBuffer.value += data.content;
          const md = require('markdown-it')({breaks: true});
          const reply = chatData.value[chatData.value.length - 1]
          reply['orgContent'] = lineBuffer.value;
          reply['content'] = md.render(lineBuffer.value);

          nextTick(() => {
            hl.configure({ignoreUnescapedHTML: true})
            const lines = document.querySelectorAll('.chat-line');
            const blocks = lines[lines.length - 1].querySelectorAll('pre code');
            blocks.forEach((block) => {
              hl.highlightElement(block)
            })
          })
        }
        // 将聊天框的滚动条滑动到最底部
        nextTick(() => {
          document.getElementById('chat-box').scrollTo(0, document.getElementById('chat-box').scrollHeight)
        })
      };
    }

  });

  _socket.addEventListener('close', () => {
    if (activelyClose.value) { // 忽略主动关闭
      return;
    }
    // 停止发送消息
    disableInput(true)
    socket.value = null;
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
  showReGenerate.value = false;
  showStopGenerate.value = !force;
}

const enableInput = () => {
  canSend.value = true;
  showReGenerate.value = previousText.value !== "";
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
    content: renderInputText(prompt.value),
    created_at: new Date().getTime(),
  });

  nextTick(() => {
    document.getElementById('chat-box').scrollTo(0, document.getElementById('chat-box').scrollHeight)
  })

  showHello.value = false
  disableInput(false)
  socket.value.send(prompt.value);
  previousText.value = prompt.value;
  prompt.value = '';
  return true;
}

const showConfig = function () {
  showConfigDialog.value = true;
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
    removeUserToken();
    router.push('/login');
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

    const md = require('markdown-it')({breaks: true});
    // md.use(require('markdown-it-copy')); // 代码复制功能
    for (let i = 0; i < data.length; i++) {
      if (data[i].type === "prompt") {
        chatData.value.push(data[i]);
        continue;
      } else if (data[i].type === "mj") {
        data[i].content = JSON.parse(data[i].content)
        data[i].content.html = md.render(data[i].content?.content)
        chatData.value.push(data[i]);
        continue;
      }

      data[i].orgContent = data[i].content;
      data[i].content = md.render(data[i].content);
      chatData.value.push(data[i]);
    }

    nextTick(() => {
      hl.configure({ignoreUnescapedHTML: true})
      const blocks = document.querySelector("#chat-box").querySelectorAll('pre code');
      blocks.forEach((block) => {
        hl.highlightElement(block)
      })
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
const reGenerate = function () {
  disableInput(false)
  const text = '重新生成上述问题的答案：' + previousText.value;
  // 追加消息
  chatData.value.push({
    type: "prompt",
    id: randString(32),
    icon: loginUser.value.avatar,
    content: renderInputText(text)
  });
  socket.value.send(text);
}

const chatName = ref('')
// 搜索会话
const searchChat = function () {
  if (chatName.value === '') {
    chatList.value = allChats.value
    return
  }
  const items = [];
  for (let i = 0; i < allChats.value.length; i++) {
    if (allChats.value[i].title.toLowerCase().indexOf(chatName.value.toLowerCase()) !== -1) {
      items.push(allChats.value[i]);
    }
  }
  chatList.value = items;
}

const updateUser = function (data) {
  loginUser.value.avatar = data.avatar;
  loginUser.value.nickname = data.nickname;
}

// 导出会话
const exportChat = () => {
  if (!activeChat.value['chat_id']) {
    return ElMessage.error("请先选中一个会话")
  }

  const url = location.protocol + '//' + location.host + '/chat/export?chat_id=' + activeChat.value['chat_id']
  // console.log(url)
  window.open(url, '_blank');
}

const getChatById = (chatId) => {
  for (let index in chatList.value) {
    if (chatList.value[index].chat_id === chatId) {
      return chatList.value[index]
    }
  }
  return null
}

const getModelValue = (model_id) => {
  for (let i = 0; i < models.value.length; i++) {
    if (models.value[i].id === model_id) {
      return models.value[i].value
    }
  }
  return ""
}
</script>

<style scoped lang="stylus">
@import "@/assets/css/chat-plus.styl"
</style>