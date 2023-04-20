<template>
  <div class="chat-free-page">
    <div class="sidebar" id="sidebar">
      <div class="block-btn new-chat" @click="newChat">
        <div class="block-btn-container text-center">
          <span class="icon"><el-icon><Plus/></el-icon></span>
          <span class="text">新建会话</span>
          <span class="btn" @click="toggleSidebar"><el-button size="small" type="info" circle><el-icon><CloseBold/></el-icon></el-button></span>
        </div>
      </div>

      <nav>
        <ul>
          <li v-for="chat in chatList" :key="chat.id" @click="changeChat(chat)"
              :class="chat.id === curChat.id ? 'active' : ''"><a>
            <span class="icon"><el-icon><ChatRound/></el-icon></span>

            <span class="text" v-if="chat.edit">
              <el-input v-model="tmpChatTitle" size="small" placeholder="请输入会话标题"/>
            </span>
            <span class="text" v-else>{{ chat.title }}</span>

            <span class="btn btn-check" v-if="chat.edit || chat.removing">
              <el-icon @click="confirm(chat)"><Check/></el-icon>
              <el-icon @click="cancel(chat)"><Close/></el-icon>
            </span>
            <span class="btn" v-else>
              <el-icon title="编辑" @click="editChatTitle(chat)"><Edit/></el-icon>
              <el-icon title="删除会话" @click="removeChat(chat)"><Delete/></el-icon>
            </span>

          </a></li>

        </ul>
      </nav>

      <div class="block-btn no-border mt-10">
        <div class="block-btn-container" @click="clearChatHistory">
          <span class="icon"><el-icon><DeleteFilled/></el-icon></span>
          <span class="text">清空聊天记录</span>
        </div>
      </div>

      <div class="block-btn no-border">
        <div class="block-btn-container" @click="logout">
          <span class="icon"><el-icon><Monitor/></el-icon></span>
          <span class="text">退出登录</span>
        </div>
      </div>
    </div>

    <div class="main-content" v-loading="loading" element-loading-background="rgba(122, 122, 122, 0.8)">
      <div class="title">
        <span class="icon" @click="toggleSidebar">
          <el-icon><Fold/></el-icon>
        </span>
        <span class="text">{{ curChat ? curChat.title : '新建会话' }}</span>
        <span class="plus-icon" @click="newChat">
          <el-icon><Plus/></el-icon>
        </span>
      </div>

      <div class="chat-container">
        <div class="chat-box" id="chat-box">
          <div v-for="chat in chatData" :key="chat.id">
            <chat-prompt
                v-if="chat.type==='prompt'"
                :icon="chat.icon"
                :content="chat.content"/>
            <chat-reply v-else-if="chat.type==='reply'"
                        :icon="replyIcon"
                        :org-content="chat.orgContent"
                        :content="chat.content"/>
          </div>

        </div><!-- end chat box -->
      </div>

      <div class="input-box" :style="{width: inputBoxWidth+'px'}">
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

        <div class="input-wrapper">
          <div class="input-container">
            <el-input
                ref="text-input"
                v-model="inputValue"
                :autosize="{ minRows: 1, maxRows: 5 }"
                v-on:keydown="inputKeyDown"
                v-on:focus="focus"
                autofocus
                type="textarea"
                placeholder="开始你的提问"
            />
          </div>

          <div class="btn-container">
            <el-row>
              <el-button type="success" class="send" :disabled="sending" v-on:click="sendMessage">发送
              </el-button>
            </el-row>
          </div>
        </div>

      </div><!-- end input box -->

    </div>

    <div class="token-dialog">
      <el-dialog
          v-model="showLoginDialog"
          :show-close="false"
          :close-on-click-modal="false"
          top="5vh"
          title="请输入口令继续访问"
      >
        <el-row>
          <el-input v-model="token" placeholder="在此输入口令" type="password" @keyup="loginInputKeyup">
            <template #prefix>
              <el-icon class="el-input__icon">
                <Lock/>
              </el-icon>
            </template>
          </el-input>
          <el-button type="primary" @click="submitToken">提交</el-button>
        </el-row>

        <div class="tip-text">
          <p>扫码加入群聊，在群公告获取免费体验账号</p>
          <el-alert type="warning" :closable="false">
            <strong>特别声明：</strong> 我们充分尊重用户隐私，因此所有用户的聊天记录均只保存在本地设备，所以请尽量用同一设备访问，以便能查阅所有的聊天记录。
          </el-alert>

        </div>

        <el-row class="row-center">
          <el-image src="https://img.r9it.com/chatgpt/wechat-group.jpeg" fit="cover"/>
        </el-row>

      </el-dialog>
    </div>
  </div>
</template>

<script>
import {defineComponent, nextTick} from "vue"
import {
  ChatRound, Check, Close,
  CloseBold,
  Delete, DeleteFilled, Edit,
  Fold,
  Lock, Monitor,
  Plus,
  RefreshRight,
  VideoPause
} from "@element-plus/icons-vue";
import {httpGet, httpPost} from "@/utils/http";
import hl from "highlight.js";
import ChatReply from "@/components/ChatReply.vue";
import ChatPrompt from "@/components/ChatPrompt.vue";
import {
  setChat,
  appendChatHistory,
  getChatHistory,
  getChatList,
  getSessionId,
  getUserInfo,
  setLoginUser, removeChat, clearChatHistory
} from "@/utils/storage";
import {ElMessage, ElMessageBox} from "element-plus";
import {isMobile, randString} from "@/utils/libs";
import Clipboard from "clipboard";

// 免费版 ChatGPT
export default defineComponent({
  name: 'ChatFree',
  components: {
    Monitor,
    DeleteFilled,
    Close,
    Check,
    Edit,
    Delete,
    CloseBold,
    Lock,
    VideoPause,
    RefreshRight,
    ChatPrompt,
    ChatReply,
    ChatRound,
    Plus,
    Fold
  },
  data() {
    return {
      chatData: [],
      inputValue: '', // 聊天内容

      userInfo: {},
      showLoginDialog: false,
      role: 'gpt',
      replyIcon: 'images/avatar/yi_yan.png', // 回复信息的头像

      chatList: [], // 会话列表
      tmpChatTitle: '',
      curOpt: '', // 当前操作
      curChat: null, // 当前会话
      curPrompt: null, // 当前用户输入

      showStopGenerate: false,
      showReGenerate: false,
      canReGenerate: false, // 是否可以重新生
      previousText: '', // 上一次提问

      token: '', // 会话 token
      lineBuffer: '', // 输出缓冲行
      errorMessage: null, // 错误信息提示框
      socket: null,
      activelyClose: false, // 主动关闭
      inputBoxWidth: 0,
      sending: true,
      loading: true,

    }
  },

  mounted() {
    const clipboard = new Clipboard('.copy-reply');
    clipboard.on('success', () => {
      ElMessage.success('复制成功！');
    })

    clipboard.on('error', () => {
      ElMessage.error('复制失败！');
    })

    nextTick(() => {
      if (isMobile()) {
        this.inputBoxWidth = window.innerWidth - 20;
      } else {
        this.inputBoxWidth = window.innerWidth - document.getElementById('sidebar').offsetWidth - 20;
      }
    })

    window.addEventListener("resize", () => {
      if (isMobile()) {
        this.inputBoxWidth = window.innerWidth - 20;
      } else {
        this.inputBoxWidth = window.innerWidth - document.getElementById('sidebar').offsetWidth - 20;
      }
    });

    this.newChat();

  },

  methods: {
    configDialog: function () {
      this.showConfigDialog = true;
      this.userInfo = getUserInfo();
    },
    // 创建 socket 会话连接
    connect: function () {
      // 先关闭已有连接
      if (this.socket !== null) {
        this.activelyClose = true;
        this.socket.close();
      }

      // 初始化 WebSocket 对象
      const sessionId = getSessionId();
      const socket = new WebSocket(process.env.VUE_APP_WS_HOST + `/api/chat?sessionId=${sessionId}&role=${this.role}`);
      socket.addEventListener('open', () => {
        this.sending = false; // 允许用户发送消息
        this.loading = false; // 隐藏加载层
        if (this.errorMessage !== null) {
          this.errorMessage.close(); // 关闭错误提示信息
        }
        this.activelyClose = false;
        // 加载聊天列表
        const chatList = getChatList();
        if (chatList) {
          this.chatList = chatList;
        }
      });

      socket.addEventListener('message', event => {
        if (event.data instanceof Blob) {
          const reader = new FileReader();
          reader.readAsText(event.data, "UTF-8");
          reader.onload = () => {
            const data = JSON.parse(String(reader.result));
            // 有聊天记录就不输出打招呼消息
            if (data['is_hello_msg'] && this.chatData.length > 1) {
              return
            }

            if (data.type === 'start') {
              this.chatData.push({
                type: "reply",
                id: randString(32),
                icon: this.replyIcon,
                content: "",
              });
              if (data['is_hello_msg'] !== true) {
                this.canReGenerate = true;
              }
            } else if (data.type === 'end') {
              this.sending = false;
              if (data['is_hello_msg'] !== true) {
                this.showReGenerate = true;
                // 保存聊天记录
                appendChatHistory(this.curChat.id, this.curPrompt);
                appendChatHistory(this.curChat.id, {
                  type: "reply",
                  id: randString(32),
                  icon: this.replyIcon,
                  content: this.lineBuffer,
                })
              }
              this.showStopGenerate = false;

              this.lineBuffer = ''; // 清空缓冲

            } else {
              this.lineBuffer += data.content;
              this.chatData[this.chatData.length - 1]['orgContent'] = this.lineBuffer;
              let md = require('markdown-it')();
              this.chatData[this.chatData.length - 1]['content'] = md.render(this.lineBuffer);

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

      socket.addEventListener('close', () => {
        if (this.activelyClose) { // 忽略主动关闭
          return;
        }

        // 停止送消息
        this.sending = true;
        this.checkSession();
      });

      this.socket = socket;
    },

    checkSession: function () {
      // 检查会话
      httpGet("/api/session/get").then(() => {
        // 自动重新连接
        this.connect();
      }).catch((res) => {
        if (res.code === 400) {
          this.showLoginDialog = true;
          if (this.errorMessage !== null) {
            this.errorMessage.close();
          }
        } else {
          if (this.errorMessage === null) {
            this.errorMessage = ElMessage({
              message: '当前无法连接服务器，可检查网络设置是否正常',
              type: 'error',
              duration: 0,
              showClose: false
            });
          }
          // 3 秒后继续重连
          setTimeout(() => this.checkSession(), 3000)
        }
      })
    },

    // 从后端获取聊天历史记录
    fetchChatHistory: function (chatId) {
      const list = getChatHistory(chatId);
      if (list) {
        const md = require('markdown-it')();
        console.log(list)
        for (let i = 0; i < list.length; i++) {
          if (list[i].type === "prompt") {
            this.chatData.push(list[i]);
            continue;
          }
          list[i].orgContent = list[i].content;
          list[i].content = md.render(list[i].content);
          this.chatData.push(list[i]);
        }

        nextTick(() => {
          hl.configure({ignoreUnescapedHTML: true})
          const blocks = document.querySelector("#chat-box").querySelectorAll('pre code');
          blocks.forEach((block) => {
            hl.highlightElement(block)
          })
        })
      }
    },

    inputKeyDown: function (e) {
      if (e.keyCode === 13) {
        if (this.sending) {
          e.preventDefault();
        } else {
          this.sendMessage();
        }
      }
    },

    // 发送消息
    sendMessage: function (e) {
      // 强制按钮失去焦点
      if (e) {
        let target = e.target;
        if (target.nodeName === "SPAN") {
          target = e.target.parentNode;
        }
        target.blur();
      }

      if (this.sending || this.inputValue.trim().length === 0) {
        return false;
      }

      // 追加消息
      this.curPrompt = {
        type: "prompt",
        id: randString(32),
        icon: 'images/avatar/user.png',
        content: this.inputValue
      };
      this.chatData.push(this.curPrompt);

      this.sending = true;
      this.showStopGenerate = true;
      this.showReGenerate = false;
      this.socket.send(this.inputValue);
      this.$refs["text-input"].blur();
      this.previousText = this.inputValue;
      this.inputValue = '';
      // 等待 textarea 重新调整尺寸之后再自动获取焦点
      setTimeout(() => this.$refs["text-input"].focus(), 100);
      return true;
    },

    // 获取焦点
    focus: function () {
      setTimeout(function () {
        document.getElementById('chat-box').scrollTo(0, document.getElementById('chat-box').scrollHeight)
      }, 200)
    },

    // 提交 Token
    submitToken: function () {
      this.showLoginDialog = false;
      this.loading = true

      // 获取会话
      httpPost("/api/login", {
        token: this.token
      }).then((res) => {
        setLoginUser(res.data)
        this.connect();
        this.loading = false;
      }).catch(() => {
        ElMessage.error("口令错误");
        this.token = '';
        this.showLoginDialog = true;
        this.loading = false;
      })
    },

    // 登录输入框输入事件处理
    loginInputKeyup: function (e) {
      if (e.keyCode === 13) {
        this.submitToken();
      }
    },

    // 清空聊天记录
    clearChatHistory: function () {
      ElMessageBox.confirm(
          '确认要清空所有聊天历史记录吗?<br/>此操作不可以撤销！',
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
        clearChatHistory();
        this.chatData = [];
        this.chatList = [];
        ElMessage.success("当前角色会话已清空");
      }).catch(() => {
      })
    },

    // 停止生成
    stopGenerate: function () {
      this.showStopGenerate = false;
      httpPost("/api/chat/stop").then(() => {
        console.log("stopped generate.")
        this.sending = false;
        if (this.canReGenerate) {
          this.showReGenerate = true;
        }
      })
    },

    // 重新生成
    reGenerate: function () {
      this.sending = true;
      this.showStopGenerate = true;
      this.showReGenerate = false;
      this.socket.send('重新生成上述问题的答案：' + this.previousText);
    },

    // 显示/关闭侧边栏
    toggleSidebar: function () {
      document.getElementById("sidebar").classList.toggle('show');
    },

    // 新建会话
    newChat: function () {
      // 判断当前会话是否已经有聊天记录
      if (this.curChat !== null) {
        const chatHistory = getChatHistory(this.curChat.id);
        if (chatHistory === null) {
          return;
        }
        this.curChat.title = chatHistory[0].content;
        // 追加会话
        setChat(this.curChat);
      }

      this.curChat = {
        id: randString(32),
        edit: false, // 是否处于编辑模式
        removing: false, // 是否处于删除模式
        title: '新会话 - 0'
      };

      this.chatData = [];
      this.showReGenerate = false;
      this.showStopGenerate = false;
      this.loading = true;
      this.sending = true;
      this.connect();
    },

    // 编辑会话标题
    editChatTitle: function (chat) {
      chat.edit = true;
      this.curOpt = 'edit';
      this.tmpChatTitle = chat.title;
    },

    // 确认修改
    confirm: function (chat) {
      if (this.curOpt === 'edit') {
        chat.title = this.tmpChatTitle;
        chat.edit = false;
        setChat(chat)
      } else if (this.curOpt === 'remove') {
        delete this.chatList[chat.id];
        if (this.curChat.id === chat.id) {
          this.chatData = [];
        }
        removeChat(chat.id);
        chat.removing = false;
      }

    },
    // 取消修改
    cancel: function (chat) {
      chat.edit = false;
      chat.removing = false;
    },

    // 删除会话
    removeChat: function (chat) {
      chat.removing = true;
      this.curOpt = 'remove';
    },

    // 切换会话
    changeChat: function (chat) {
      if (this.curChat.id === chat.id) {
        return;
      }

      this.curChat = chat;
      this.fetchChatHistory(chat.id);
    },

    // 退出登录
    logout: function () {
      httpPost("/api/logout", {opt: "logout"}).then(() => {
        this.checkSession();
      }).catch(() => {
        ElMessage.error("注销失败");
      })
    },

  }
})
</script>

<style lang="stylus">
.chat-free-page {
  display: flex;
  flex-direction: row;
  height 100%;
  background-color: rgba(247, 247, 248, 1);

  .sidebar {
    background-color: #1D1E20
    height: 100%;
    width 350px;

    .block-btn {
      font-size 14px;

      .block-btn-container {
        border: 1px solid #4A4B4D;
        padding: 8px 10px
        color: #ffffff
        cursor pointer
        box-sizing: border-box;
        border-radius 5px;
        position relative;

        &:hover {
          background-color #3E3F49
        }

        .text {
          margin-left 5px;
        }

        .btn {
          position absolute;
          display none
          right 8px;
          top 8px;

          .el-icon {
            margin-left 0;
            color #ffffff
          }
        }
      }

      .text-center {
        text-align center;
      }
    }


    .new-chat {
      padding 10px 10px 0 10px;
    }

    .mt-10 {
      margin-top 10px;
    }

    .no-border {
      .block-btn-container {
        border none
        padding 10px;

        .text {
          margin-left 12px;
        }
      }
    }

    nav {
      margin 0
      padding 0
      width 100%
      height calc(100vh - 250px)
      overflow-y auto

      ul {
        list-style-type: none
        padding 8px
        margin 0

        li {
          padding: 10px
          color: #ffffff
          cursor pointer
          margin-bottom 10px;
          box-sizing: border-box;
          border-radius 5px;

          &:hover {
            background-color #3E3F49

            a {
              .btn {
                display block
                background-color: #3E3F49;
              }
            }
          }

          a {
            display flex
            text-decoration: none;
            position relative

            .icon {
              font-size 16px;
              margin-top 3px;
              margin-right 8px;
            }

            .text {
              font-size 14px;
              padding-top 2px;
              overflow hidden
              white-space: nowrap;
              text-overflow: ellipsis;
              max-width 200px;
            }

            .btn {
              display none
              position absolute
              right 0;
              top 2px;

              .el-icon {
                margin-left 5px;
                color #9f9f9f
              }

              .el-icon:hover {
                color #ffffff
              }

            }

            .btn-check {
              .el-icon {
                margin-left 10px;
                font-size 18px;
              }
            }
          }

        }

        li.active {
          background-color #3E3F49
        }

        li.block-btn {
          border: 1px solid #4A4B4D;

          a {
            .btn {
              display none
              right -2px;
              top -2px;

              .el-icon {
                margin-left 0;
                color #ffffff
              }
            }
          }
        }

      }
    }

    nav::-webkit-scrollbar {
      width: 0;
      height: 0;
    }
  }

  .main-content {
    width 100%;
    display: flex;
    flex-direction: column;
    justify-content: space-between;
    position: relative;
    overflow hidden
    background-image url("~@/assets/img/bg_01.jpeg")

    .title {
      height 30px
      width 100%
      padding 5px 0
      font-size 16px
      color: rgba(217, 217, 227, 1)
      background-color: rgba(52, 53, 65, 1)
      display none

      .el-icon {
        font-size 24px;
        cursor pointer
        padding-left 10px
      }

      .text {
        width: 100%
        text-align: center;
      }

    }

    .chat-container {
      display flex
      justify-content flex-start

      .chat-box {
        width 100%;
        padding: 10px;
        overflow-y: auto;
        height: calc(100vh - 80px);
        // 变量定义
        --content-font-size: 16px;
        --content-color: #374151;

        font-family 'Microsoft YaHei', '微软雅黑', Arial, sans-serif;

        .chat-line {
          padding 10px 5px;
          font-size 14px;
          display: flex;
          align-items: flex-start;

          .chat-icon {
            img {
              width 32px;
              height 32px;
            }
          }
        }

      }
    }


    .input-box {
      padding 10px;
      background #ffffff;
      position: absolute;
      bottom: 0
      display: flex;
      justify-content: flex-start;
      align-items: center;
      flex-flow: column;

      .re-generate {
        position relative
        display flex
        justify-content center

        .btn-box {
          position absolute
          bottom 20px

          .el-icon {
            margin-right 5px;
          }
        }
      }

      .input-wrapper {
        width 100%;
        display flex;

        .input-container {
          overflow hidden
          width 100%
          margin: 0;
          border: none;
          border-radius: 6px;
          box-shadow: 0 2px 15px rgba(0, 0, 0, 0.1);
          background-color: rgba(255, 255, 255, 1);
          padding: 5px 10px;

          .el-textarea__inner {
            box-shadow: none
            padding 5px 0
            min-height 24px !important;
          }

          .el-textarea__inner::-webkit-scrollbar {
            width: 0;
            height: 0;
          }
        }

        .btn-container {
          margin-left 10px;

          .el-row {
            flex-wrap nowrap
            //width 106px;
            align-items center
          }

          .send {
            width 60px;
            height 40px;
            background-color: var(--el-color-success)
          }

          .is-disabled {
            background-color: var(--el-button-disabled-bg-color);
            border-color: var(--el-button-disabled-border-color);
          }
        }
      }

      // end of input wrapper

    }

    // end of input box

  }
}

.row-center {
  justify-content center
}

/* 移动端适配 */
@media (max-width: 768px) {
  .chat-free-page {
    flex-direction: column;

    .sidebar {
      width: 90%;
      position absolute;
      z-index 9999;
      top: 0;
      left: -350px;
      transition: transform 0.3s ease-in-out;

      .block-btn {
        .block-btn-container {
          .btn {
            display inline
          }
        }
      }
    }

    .sidebar.show {
      transform: translateX(350px);
    }

    .main-content {
      width: 100%;

      .title {
        display flex

        .plus-icon {
          padding-right 10px;
        }
      }

      .chat-container {
        height: calc(100vh - 40px);

        .chat-box {
          height: calc(100vh - 120px);
        }
      }

    }
  }

}
</style>
