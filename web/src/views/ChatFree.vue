<template>
  <div class="chat-free-page">
    <div class="sidebar" id="sidebar">
      <nav>
        <ul>
          <li class="new-chat"><a>
            <span class="icon"><el-icon><Plus/></el-icon></span>
            <span class="text">新建会话</span>
            <span class="btn" @click="toggleSidebar"><el-button size="small" type="info" circle><el-icon><CloseBold/></el-icon></el-button></span>
          </a></li>
          <li><a>
            <span class="icon"><el-icon><ChatRound/></el-icon></span>
            <span class="text">会话一</span>
          </a></li>
          <li class="active"><a>
            <span class="icon"><el-icon><ChatRound/></el-icon></span>
            <span class="text">会话二</span>
          </a></li>
        </ul>
      </nav>
    </div>

    <div class="main-content" v-loading="loading" element-loading-background="rgba(122, 122, 122, 0.8)">
      <div class="title">
        <span class="icon" @click="toggleSidebar">
          <el-icon><Fold/></el-icon>
        </span>
        <span class="text">响应式页面布局代码</span>
      </div>

      <div class="chat-container">
        <div class="chat-box" id="chat-box">
          <div v-for="chat in chatData" :key="chat.id">
            <chat-prompt
                v-if="chat.type==='prompt'"
                :icon="chat.icon"
                :content="chat.content"/>
            <chat-reply v-else-if="chat.type==='reply'"
                        :icon="chat.icon"
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
              <el-button type="success" class="send" :disabled="sending" v-on:click="sendMessage">发送</el-button>
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
          打开微信扫下面二维码免费领取口令, <strong>强烈建议你使用 PC 浏览器访问获得更好的聊天体验。</strong>
        </div>

        <el-row class="row-center">
          <el-image src="images/wx.png" fit="cover"/>
        </el-row>

      </el-dialog>
    </div>
  </div>
</template>

<script>
import {defineComponent, nextTick} from "vue"
import {ChatRound, CloseBold, Fold, Lock, Plus, RefreshRight, VideoPause} from "@element-plus/icons-vue";
import {httpGet, httpPost} from "@/utils/http";
import hl from "highlight.js";
import ChatReply from "@/components/ChatReply.vue";
import ChatPrompt from "@/components/ChatPrompt.vue";
import {getSessionId, getUserInfo, setLoginUser} from "@/utils/storage";
import {ElMessage, ElMessageBox} from "element-plus";
import {isMobile, randString} from "@/utils/libs";
import Clipboard from "clipboard";

// 免费版 ChatGPT
export default defineComponent({
  name: 'ChatFree',
  components: {CloseBold, Lock, VideoPause, RefreshRight, ChatPrompt, ChatReply, ChatRound, Plus, Fold},
  data() {
    return {
      chatData: [],
      inputValue: '', // 聊天内容

      userInfo: {},
      showLoginDialog: false,
      role: 'gpt',
      replyIcon: 'images/avatar/gpt.png', // 回复信息的头像

      showStopGenerate: false,
      showReGenerate: false,
      canReGenerate: false, // 是否可以重新生
      previousText: '', // 上一次提问

      token: '', // 会话 token
      lineBuffer: '', // 输出缓冲行
      errorMessage: null, // 错误信息提示框
      socket: null,
      inputBoxWidth: 0,
      sending: true,
      loading: true,

    }
  },

  mounted() {
    this.fetchChatHistory();

    const clipboard = new Clipboard('.reply-content');
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

    this.connect();

  },

  methods: {
    configDialog: function () {
      this.showConfigDialog = true;
      this.userInfo = getUserInfo();
    },
    // 创建 socket 会话连接
    connect: function () {
      // 初始化 WebSocket 对象
      const sessionId = getSessionId();
      const socket = new WebSocket(process.env.VUE_APP_WS_HOST + `/api/chat?sessionId=${sessionId}&role=${this.role}`);
      socket.addEventListener('open', () => {
        this.sending = false; // 允许用户发送消息
        this.loading = false; // 隐藏加载层
        if (this.errorMessage !== null) {
          this.errorMessage.close(); // 关闭错误提示信息
        }
        // 加载聊天记录
        this.fetchChatHistory();
      });

      socket.addEventListener('message', event => {
        if (event.data instanceof Blob) {
          const reader = new FileReader();
          reader.readAsText(event.data, "UTF-8");
          reader.onload = () => {
            const data = JSON.parse(String(reader.result));
            if (data.type === 'start') {
              this.chatData.push({
                type: "reply",
                id: randString(32),
                icon: this.replyIcon,
                content: "",
                cursor: true
              });
              if (data['is_hello_msg'] !== true) {
                this.canReGenerate = true;
              }
            } else if (data.type === 'end') {
              this.sending = false;
              if (data['is_hello_msg'] !== true) {
                this.showReGenerate = true;
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
    fetchChatHistory: function () {
      httpPost("/api/chat/history", {role: this.role}).then((res) => {
        if (this.chatData.length > 0) { // 如果已经有聊天记录了，就不追加了
          return
        }

        const data = res.data
        const md = require('markdown-it')();
        for (let i = 0; i < data.length; i++) {
          if (data[i].type === "prompt") {
            this.chatData.push(data[i]);
            continue;
          }
          data[i].orgContent = data[i].content;
          data[i].content = md.render(data[i].content);
          this.chatData.push(data[i]);
        }

        nextTick(() => {
          hl.configure({ignoreUnescapedHTML: true})
          const blocks = document.querySelector("#chat-box").querySelectorAll('pre code');
          blocks.forEach((block) => {
            hl.highlightElement(block)
          })
        })
      }).catch(() => {
        // console.error(e.message)
      })
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
      this.chatData.push({
        type: "prompt",
        id: randString(32),
        icon: 'images/avatar/user.png',
        content: this.inputValue
      });

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
          '确认要清空当前角色聊天历史记录吗?<br/>此操作不可以撤销！',
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
        httpPost("/api/chat/history/clear", {role: this.role}).then(() => {
          ElMessage.success("当前角色会话已清空");
          this.chatData = [];
        }).catch(() => {
          ElMessage.error("删除失败")
        })
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

    nav {
      margin 0
      padding 0
      width 100%

      ul {
        list-style-type: none
        padding 5px
        margin 0

        li {
          padding: 10px
          color: #ffffff
          cursor pointer
          margin-bottom 10px;
          box-sizing: border-box;
          border-radius 5px;

          &:hover {
            background-color #2E2F39
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
            }

          }

        }

        li.active {
          background-color #2E2F39
        }

        li.new-chat {
          border: 1px solid #4A4B4D;

          .btn {
            display none
            position absolute
            right -2px;
            top -2px;
          }
        }

      }
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
      justify-content: start;
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

      nav {
        ul {
          li.new-chat {
            .btn {
              display inline
            }
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
