<template>
  <div class="common-layout theme-white">
    <el-container>
      <el-aside width="320px">
        <div class="title-box">
          <el-image :src="logo" class="logo"/>
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
              <el-input v-model="tmpChatTitle" size="small" placeholder="请输入会话标题"/>
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
          <el-dropdown :hide-on-click="true" class="user-info" trigger="click" v-if="user">
                        <span class="el-dropdown-link">
                          <el-image :src="user['avatar']"/>
                          <span class="username">{{ user ? user['nickname'] : 'Chat-Plus-User' }}</span>
                          <el-icon><ArrowDown/></el-icon>
                        </span>
            <template #dropdown>
              <el-dropdown-menu style="width: 315px;">
                <el-dropdown-item @click="showConfig">
                  <el-icon>
                    <Tools/>
                  </el-icon>
                  <span>聊天设置</span>
                </el-dropdown-item>

                <el-dropdown-item @click="clearAllChats">
                  <el-icon>
                    <Delete/>
                  </el-icon>
                  <span>清除所有会话</span>
                </el-dropdown-item>

                <el-dropdown-item @click="logout">
                  <el-icon>
                    <Monitor/>
                  </el-icon>
                  <span>注销</span>
                </el-dropdown-item>

<!--                <el-dropdown-item>-->
<!--                  <i class="icon-font icon-github"></i>-->
<!--                 <span>-->
<!--                    <el-link type="primary" href="https://github.com/yangjian102621/chatgpt-plus" target="_blank">ChatGPT-Plus-V3</el-link>-->
<!--                 </span>-->
<!--                </el-dropdown-item>-->
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-aside>
      <el-main v-loading="loading" element-loading-background="rgba(122, 122, 122, 0.3)">
        <div class="chat-head">
          <div class="chat-config">
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

            <el-select v-model="model" placeholder="模型">
              <el-option
                  v-for="item in models"
                  :key="item"
                  :label="item.toUpperCase()"
                  :value="item"
              />
            </el-select>
            <el-button type="primary" @click="newChat">
              <el-icon>
                <Plus/>
              </el-icon>
              新建会话
            </el-button>

          </div>
        </div>

        <div class="right-box" :style="{height: mainWinHeight+'px'}">
          <div>
            <div id="container">
              <div class="chat-box" id="chat-box" :style="{height: chatBoxHeight+'px'}">
                <div v-for="item in chatData" :key="item.id">
                  <chat-prompt
                      v-if="item.type==='prompt'"
                      :icon="item.icon"
                      :created-at="dateFormat(item['created_at'])"
                      :tokens="item['tokens']"
                      :model="model"
                      :content="item.content"/>
                  <chat-reply v-else-if="item.type==='reply'"
                              :icon="item.icon"
                              :org-content="item.orgContent"
                              :created-at="dateFormat(item['created_at'])"
                              :tokens="item['tokens']"
                              :content="item.content"/>
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
                      ref="text-input"
                      v-model="prompt"
                      v-on:keydown="inputKeyDown"
                      autofocus
                      type="textarea"
                      :rows="2"
                      placeholder="按 Enter 键发送消息，使用 Ctrl + Enter 换行"
                  />
                  <span class="send-btn">
                    <el-button @click="sendMessage">
                      <el-icon><Promotion /></el-icon>
                    </el-button>
                  </span>
                </div>
              </div><!-- end input box -->

            </div><!-- end container -->
          </div><!-- end loading -->
        </div>
      </el-main>
    </el-container>

    <config-dialog :show="showConfigDialog" :models="models" @hide="showConfigDialog = false" @update-user="updateUser" />
  </div>


</template>
<script setup>
import {nextTick, onMounted, ref} from 'vue'
import ChatPrompt from "@/components/ChatPrompt.vue";
import ChatReply from "@/components/ChatReply.vue";
import {
  ArrowDown,
  Check,
  Close,
  Delete,
  Edit,
  Monitor,
  Plus, Promotion,
  RefreshRight,
  Search,
  Tools,
  VideoPause
} from '@element-plus/icons-vue'
import 'highlight.js/styles/a11y-dark.css'
import {dateFormat, randString, removeArrayItem, renderInputText, UUID} from "@/utils/libs";
import {ElMessage, ElMessageBox} from "element-plus";
import hl from "highlight.js";
import {getLoginUser, getSessionId} from "@/utils/storage";
import {httpGet, httpPost} from "@/utils/http";
import {useRouter} from "vue-router";
import Clipboard from "clipboard";
import ConfigDialog from "@/components/ConfigDialog.vue";

const title = ref('ChatGPT-智能助手');
const logo = 'images/logo.png';
const models = ref([])
const model = ref('gpt-3.5-turbo')
const chatData = ref([]);
const allChats = ref([]); // 会话列表
const chatList = ref(allChats.value);
const activeChat = ref({});
const mainWinHeight = ref(0); // 主窗口高度
const chatBoxHeight = ref(0); // 聊天内容框高度
const leftBoxHeight = ref(0);
const loading = ref(true);
const user = ref(getLoginUser());
const roles = ref([]);
const roleId = ref(0)
const newChatItem = ref(null);
const router = useRouter();
const showConfigDialog = ref(false);

if (!user.value) {
  router.push("login");
} else {
  onMounted(() => {
    resizeElement();
    checkSession().then(() => {
      // 加载角色列表
      httpGet(`/api/chat/role/list?user_id=${user.value.id}`).then((res) => {
        roles.value = res.data;
        roleId.value = roles.value[0]['id'];
        // 获取会话列表
        loadChats();
        // 创建新的会话
        newChat();
      }).catch((e) => {
        console.log(e)
        ElMessage.error('获取聊天角色失败')
      })
    }).catch(() => {
      router.push('login')
    });

    const clipboard = new Clipboard('.copy-reply');
    clipboard.on('success', () => {
      ElMessage.success('复制成功！');
    })

    clipboard.on('error', () => {
      ElMessage.error('复制失败！');
    })
  });
}

// 加载系统配置
httpGet('/api/config/get?key=system').then(res => {
  title.value = res.data.title;
  models.value = res.data.models;
}).catch(e => {
  ElMessage.error("加载系统配置失败: " + e.message)
})

const checkSession = function () {
  return new Promise((resolve, reject) => {
    httpGet('/api/user/session').then(res => {
      resolve(res)
    }).catch(err => {
      reject(err)
    })
  })
}

// 加载会话
const loadChats = function () {
  httpGet("/api/chat/list?user_id=" + user.value.id).then((res) => {
    if (res.data) {
      chatList.value = res.data;
      allChats.value = res.data;
    }
  }).catch(() => {
    // TODO: 增加重试按钮
    ElMessage.error("加载会话列表失败！")
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
  chatBoxHeight.value = window.innerHeight - 51 - 82 - 38;
  mainWinHeight.value = window.innerHeight - 51;
  leftBoxHeight.value = window.innerHeight - 51 - 100;
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
    model: model.value,
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
const changeChat = function (chat) {
  if (activeChat.value['chat_id'] === chat.chat_id) {
    return;
  }
  activeChat.value = chat
  newChatItem.value = null;
  roleId.value = chat.role_id;
  model.value = chat.model;
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

// 确认修改
const confirm = function (event, chat) {
  event.stopPropagation();
  if (curOpt.value === 'edit') {
    if (tmpChatTitle.value === '') {
      ElMessage.error("请输入会话标题！");
      return;
    }

    httpPost('/api/chat/update', {id: chat.id, title: tmpChatTitle.value}).then(() => {
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
// const replyIcon = 'images/avatar/gpt.png';// 回复信息的头像
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
      host = 'wss://'+location.host;
    } else {
      host = 'ws://'+location.host;
    }
  }
  const _socket = new WebSocket(host + `/api/chat/new?session_id=${_sessionId}&role_id=${role_id}&chat_id=${chat_id}&model=${model.value}`);
  _socket.addEventListener('open', () => {
    chatData.value = []; // 初始化聊天数据
    previousText.value = '';
    canSend.value = true;
    activelyClose.value = false;

    if (isNewChat) { // 加载打招呼信息
      loading.value = false;
      chatData.value.push({
        type: "reply",
        id: randString(32),
        icon: _role['icon'],
        content: _role['hello_msg'],
        orgContent: _role['hello_msg'],
      })
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
        } else if (data.type === 'end') { // 消息接收完毕
          canSend.value = true;
          showReGenerate.value = true;
          showStopGenerate.value = false;
          lineBuffer.value = ''; // 清空缓冲

          // 追加当前会话到会话列表
          if (isNewChat && newChatItem.value !== null) {
            newChatItem.value['title'] = previousText.value;
            newChatItem.value['chat_id'] = chat_id;
            chatList.value.unshift(newChatItem.value);
            activeChat.value = newChatItem.value;
            newChatItem.value = null; // 只追加一次
          }

          // 获取 token
          const reply = chatData.value[chatData.value.length - 1]
          httpGet(`/api/chat/tokens?text=${reply.orgContent}&model=${model.value}`).then(res => {
            reply['created_at'] = new Date().getTime();
            reply['tokens'] = res.data;
            // 将聊天框的滚动条滑动到最底部
            nextTick(() => {
              document.getElementById('chat-box').scrollTo(0, document.getElementById('chat-box').scrollHeight)
            })
          })

        } else {
          lineBuffer.value += data.content;
          let md = require('markdown-it')();
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
    canSend.value = true;
    socket.value = null;
    loading.value = true;
    checkSession().then(() => {
      connect(chat_id, role_id)
    }).catch(() => {
      ElMessageBox({
        title: '会话提示',
        message: "当前会话已经失效，请重新登录",
        confirmButtonText: 'OK',
        callback: () => router.push('login')
      });
    });
  });

  socket.value = _socket;
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
    icon: user.value.avatar,
    content: renderInputText(prompt.value),
    created_at: new Date().getTime(),
  });

  nextTick(() => {
    document.getElementById('chat-box').scrollTo(0, document.getElementById('chat-box').scrollHeight)
  })

  canSend.value = false;
  showStopGenerate.value = true;
  showReGenerate.value = false;
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
    router.push('login');
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

    const md = require('markdown-it')();
    for (let i = 0; i < data.length; i++) {
      if (data[i].type === "prompt") {
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
  httpGet("/api/chat/stop?session_id=" + getSessionId() ).then(() => {
    canSend.value = true;
    if (previousText.value !== '') {
      showReGenerate.value = true;
    }
  })
}

// 重新生成
const reGenerate = function () {
  canSend.value = false;
  showStopGenerate.value = true;
  showReGenerate.value = false;
  const text = '重新生成上述问题的答案：' + previousText.value;
  // 追加消息
  chatData.value.push({
    type: "prompt",
    id: randString(32),
    icon: 'images/avatar/user.png',
    content: renderInputText(text)
  });
  socket.value.send(text);
}

const chatName = ref('')
// 搜索会话
const searchChat = function () {
  if (chatName.value === '') {
    return
  }
  const roles = [];
  for (let i = 0; i < allChats.value.length; i++) {
    if (allChats.value[i].title.indexOf(chatName.value) !== -1) {
      roles.push(allChats.value[i]);
    }
  }
  chatList.value = roles;
}

const updateUser = function (data) {
  user.value.avatar = data.avatar;
  user.value.nickname = data.nickname;
}
</script>

<style lang="stylus">
#app {
  height: 100%;

  .common-layout {
    height: 100%;

    // left side

    .el-aside {
      .title-box {
        padding: 6px 10px;
        display: flex;
        color: #ffffff;
        font-size: 20px;

        .logo {
          background-color: #ffffff
          border-radius: 50%;
          width: 40px;
          height: 40px;
        }

        span {
          padding-top: 8px;
          padding-left: 10px;
        }
      }

      .chat-list {
        display: flex
        flex-flow: column
        background-color: #28292A
        border-top: 1px solid #2F3032
        border-right: 1px solid #2F3032

        .search-box {
          flex-wrap: wrap
          padding: 10px 15px;

          .el-input__wrapper {
            background-color: #363535;
            box-shadow: none
          }
        }

        // 隐藏滚动条

        ::-webkit-scrollbar {
          width: 0;
          height: 0;
          background-color: transparent;
        }

        .content {
          //display flex
          //flex-wrap: wrap;
          //flex-direction column
          width: 100%
          overflow-y: scroll

          .chat-list-item {
            display: flex
            width: 100%
            justify-content: flex-start
            padding: 8px 12px
            border-bottom: 1px solid #3c3c3c
            cursor: pointer


            .avatar {
              width: 28px;
              height: 28px;
              border-radius: 50%;
            }

            .chat-title-input {
              font-size: 14px;
              margin-top: 4px;
              margin-left: 10px;
              overflow: hidden;
              white-space: nowrap;
              text-overflow: ellipsis;
              width: 210px;
            }

            .chat-title {
              color: #c1c1c1
              padding: 5px 10px;
              max-width 220px;
              font-size 14px;
              overflow: hidden;
              white-space: nowrap;
              text-overflow: ellipsis;
            }

            .btn {
              display none
              position: absolute;
              right: 2px;
              top: 16px;
              color #ffffff

              .el-icon {
                margin-right 8px;
              }
            }
          }

          .chat-list-item.active {
            background-color: #363535;

            .btn {
              display inline
            }
          }
        }
      }


      .tool-box {
        display: flex;
        justify-content: flex-end;
        align-items: center;
        padding 5px 20px;

        .user-info {
          width 100%
          padding-top 10px;

          .el-tooltip__trigger {
            width 100%;

            .el-dropdown-link {
              cursor: pointer
              display flex

              .el-image {
                display flex

                img {
                  width: 20px;
                  height: 20px;
                  border-radius: 5px;
                }
              }


              .username {
                display flex
                line-height 24px;
                width 230px;
                padding-left 10px;

              }

              .el-icon {
                color: #cccccc;
                line-height 24px;
              }
            }
          }

        }
      }
    }

    .el-main {
      overflow: hidden;

      .chat-head {
        width: 100%;
        height: 50px;
        background-color: #28292A

        .chat-config {
          display flex
          flex-direction row
          align-items: center;
          justify-content center;
          padding-top 10px;

          .el-select {
            //max-width 150px;
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
      }

      .right-box {
        min-width: 0;
        flex: 1;
        background-color: #ffffff
        border-left: 1px solid #4f4f4f

        #container {
          overflow: hidden;
          width: 100%;

          ::-webkit-scrollbar {
            width: 0;
            height: 0;
            background-color: transparent;
          }

          .chat-box {
            overflow-y: scroll;
            //border-bottom: 1px solid #4f4f4f

            // 变量定义
            --content-font-size: 16px;
            --content-color: #c1c1c1;

            font-family: 'Microsoft YaHei', '微软雅黑', Arial, sans-serif;
            padding: 0 0 50px 0;

            .chat-line {
              font-size: 14px;
              display: flex;
              align-items: flex-start;

            }
          }

          .re-generate {
            position: relative;
            display: flex;
            justify-content: center;

            .btn-box {
              position: absolute
              bottom: 10px;

              .el-button {
                .el-icon {
                  margin-right 5px;
                }
              }

            }
          }

          .input-box {
            background-color: #ffffff
            display: flex;
            justify-content: center;
            align-items: center;
            box-shadow: 0 2px 15px rgba(0, 0, 0, 0.1);

            .input-container {
              max-width 600px;
              width 100%
              margin: 0;
              border: none;
              padding: 10px 0;
              display flex
              justify-content center
              position relative

              .el-textarea {
                max-width 600px;
                .el-textarea__inner::-webkit-scrollbar {
                  width: 0;
                  height: 0;
                }
              }

              .send-btn {
                position absolute;
                right 12px;
                top 20px;

                .el-button {
                  padding 8px 5px;
                  border-radius 6px;
                  background:rgb(25,195,125)
                  color #ffffff;
                  font-size 20px;
                }
              }

            }
          }
        }

        #container::-webkit-scrollbar {
          width: 0;
          height: 0;
        }
      }
    }
  }

  .el-message-box {
    width: 90%;
    max-width: 420px;
  }

  .el-message {
    min-width: 100px;
    max-width: 600px;
  }

  .token-dialog {
    .el-dialog {
      --el-dialog-width: 90%;
      max-width: 400px;

      .el-dialog__body {
        padding: 10px 10px 20px 10px;
      }

      .el-row {
        flex-wrap: nowrap

        button {
          margin-left: 5px;
        }
      }
    }
  }
}

.el-select-dropdown__wrap {
  .el-select-dropdown__item {
    .role-option {
      display flex
      flex-flow row
      margin-top 8px;

      .el-image {
        width 20px
        height 20px
        border-radius 50%
      }

      span {
        margin-left 5px;
        height 20px;
        line-height 20px;
      }
    }
  }
}
</style>
