<template>
  <div class="body" v-loading="loading">
    <div id="container">
      <div class="tool-box">
        <el-image style="width: 24px; height: 24px" :src="logo"/>
        <!--        <el-button round>WeChatGPT</el-button>-->
        <el-select v-model="role" class="chat-role"
                   v-on:change="changeRole"
                   placeholder="请选择对话角色">
          <el-option
              v-for="item in chatRoles"
              :key="item.key"
              :label="item.name"
              :value="item.key"
          />
        </el-select>

        <el-button type="danger" class="clear-history" size="small" circle @click="clearChatHistory">
          <el-icon>
            <Delete/>
          </el-icon>
        </el-button>
      </div>

      <div class="chat-box" id="chat-box" :style="{height: chatBoxHeight+'px'}">
        <div v-for="chat in chatData" :key="chat.id">
          <chat-prompt
              v-if="chat.type==='prompt'"
              :icon="chat.icon"
              :content="chat.content"/>
          <chat-reply v-else-if="chat.type==='reply'"
                      :icon="chat.icon"
                      :content="chat.content"/>
        </div>

      </div><!-- end chat box -->

      <div class="input-box" :style="{width: inputBoxWidth+'px'}">
        <div class="input-container">
          <el-input
              ref="text-input"
              v-model="inputValue"
              :autosize="{ minRows: 1, maxRows: 10 }"
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
            <el-button type="info" class="config" ref="send-btn" circle @click="showConnectDialog = true">
              <el-icon>
                <Tools/>
              </el-icon>
            </el-button>
          </el-row>
        </div>

      </div><!-- end input box -->

    </div><!-- end container -->

    <config-dialog v-model:show="showConnectDialog"></config-dialog>

    <div class="token-dialog">
      <el-dialog
          v-model="showLoginDialog"
          :show-close="false"
          :close-on-click-modal="false"
          title="请输入口令继续访问"
      >
        <el-row>
          <el-input v-model="token" placeholder="在此输入口令" @keyup="loginInputKeyup">
            <template #prefix>
              <el-icon class="el-input__icon">
                <Lock/>
              </el-icon>
            </template>
          </el-input>
          <el-button type="primary" @click="submitToken">提交</el-button>
        </el-row>

      </el-dialog>
    </div>
  </div>
</template>

<script>
import {defineComponent, nextTick} from 'vue'
import ChatPrompt from "@/components/ChatPrompt.vue";
import ChatReply from "@/components/ChatReply.vue";
import {randString} from "@/utils/libs";
import {ElMessage, ElMessageBox} from 'element-plus'
import {Tools, Lock, Delete} from '@element-plus/icons-vue'
import ConfigDialog from '@/components/ConfigDialog.vue'
import {httpPost, httpGet} from "@/utils/http";
import {getSessionId, setSessionId} from "@/utils/storage";
import hl from 'highlight.js'
import 'highlight.js/styles/a11y-dark.css'

export default defineComponent({
  name: "XChat",
  components: {ChatPrompt, ChatReply, Tools, Lock, Delete, ConfigDialog},
  data() {
    return {
      title: 'ChatGPT 控制台',
      logo: 'images/logo.png',
      chatData: [],
      chatRoles: [],
      role: 'gpt',
      inputValue: '', // 聊天内容
      chatBoxHeight: 0, // 聊天内容框高度
      showConnectDialog: false,
      showLoginDialog: false,
      token: '', // 会话 token
      replyIcon: 'images/avatar/gpt.png', // 回复信息的头像

      lineBuffer: '', // 输出缓冲行
      connectingMessageBox: null, // 保存重连的消息框对象load
      socket: null,
      toolBoxHeight: 61 + 42, // 工具框的高度
      inputBoxWidth: window.innerWidth - 20,
      sending: false,
      loading: false
    }
  },

  mounted: function () {
    nextTick(() => {
      this.chatBoxHeight = window.innerHeight - this.toolBoxHeight;
    })

    // for (let i = 0; i < 10; i++) {
    //   this.chatData.push({
    //     type: "prompt",
    //     id: randString(32),
    //     icon: 'images/user-icon.png',
    //     content: "孙悟空为什么可以把金棍棒放进耳朵？",
    //   });
    //   this.chatData.push({
    //     type: "reply",
    //     id: randString(32),
    //     icon: 'images/gpt-icon.png',
    //     content: "以下是一个使用 WebSocket API 建立 WebSocket 连接并发送消息的 JavaScript 示例代码：\n" +
    //         "\n" +
    //         "```js\n" +
    //         "const socket = new WebSocket('ws://localhost:8080');\n" +
    //         "\n" +
    //         "// 监听 WebSocket 连接打开事件\n" +
    //         "socket.addEventListener('open', (event) => {\n" +
    //         "  console.log('WebSocket 连接已打开');\n" +
    //         "\n" +
    //         "  // 发送消息\n" +
    //         "  socket.send('Hello WebSocket!');\n" +
    //         "});\n" +
    //         "\n" +
    //         "// 监听 WebSocket 接收到消息事件\n" +
    //         "socket.addEventListener('message', (event) => {\n" +
    //         "  console.log('接收到消息：' + event.data);\n" +
    //         "});\n" +
    //         "\n" +
    //         "// 监听 WebSocket 连接关闭事件\n" +
    //         "socket.addEventListener('close', (event) => {\n" +
    //         "   console.log('WebSocket 连接已关闭');\n" +
    //         "});\n" +
    //         "\n" +
    //         "// 监听 WebSocket 出错事件\n" +
    //         "socket.addEventListener('error', (event) => {\n" +
    //         "   console.log('WebSocket 连接出错');\n" +
    //         "});\n" +
    //         "```\n" +
    //         "\n" +
    //         "在实际使用时，需要替换上述代码中的 WebSocket 连接地址和端口号。此外，根据后端的实现，可能需要在客户端发送的消息中携带一些特定信息，以便后端能够正确地处理这些消息。",
    //   });
    // }
    //
    // let md = require('markdown-it')();
    // this.chatData[this.chatData.length - 1]["content"] = md.render(this.chatData[this.chatData.length - 1]["content"]);
    //
    // nextTick(() => {
    //   const lines = document.querySelectorAll('.chat-line');
    //   const blocks = lines[lines.length - 1].querySelectorAll('pre code');
    //   blocks.forEach((block) => {
    //     hl.highlightElement(block)
    //   })
    // })

    window.addEventListener("resize", () => {
      this.chatBoxHeight = window.innerHeight - this.toolBoxHeight;
      this.inputBoxWidth = window.innerWidth - 20;
    });

    this.connect();
  },

  methods: {
    // 创建 socket 会话连接
    connect: function () {
      // 初始化 WebSocket 对象
      const sessionId = getSessionId();
      const socket = new WebSocket(process.env.VUE_APP_WS_HOST + `/api/chat?sessionId=${sessionId}&role=${this.role}`);
      socket.addEventListener('open', () => {
        // 获取聊天角色
        httpGet("/api/config/chat-roles/get").then((res) => {
          // ElMessage.success('创建会话成功！');
          this.chatRoles = res.data;
          this.loading = false
        }).catch(() => {
          ElMessage.error("获取聊天角色失败");
        })

        // 加载聊天记录
        this.fetchChatHistory();

        if (this.connectingMessageBox && typeof this.connectingMessageBox.close === 'function') {
          this.connectingMessageBox.close();
          this.connectingMessageBox = null;
        }
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
            } else if (data.type === 'end') {
              this.sending = false;
              this.lineBuffer = ''; // 清空缓冲
            } else {
              this.lineBuffer += data.content;
              let md = require('markdown-it')();
              this.chatData[this.chatData.length - 1]["content"] = md.render(this.lineBuffer);

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
              document.getElementById('container').scrollTo(0, document.getElementById('container').scrollHeight)
            })
          };
        }

      });
      socket.addEventListener('close', () => {
        // 检查会话
        httpGet("/api/session/get").then(() => {
          if (this.connectingMessageBox === null) {
            this.connectingMessageBox = ElMessageBox.confirm(
                '^_^ 会话发生异常，您已经从服务器断开连接!',
                '注意：',
                {
                  confirmButtonText: '重连会话',
                  cancelButtonText: '不聊了',
                  type: 'warning',
                  showClose: false,
                  closeOnClickModal: false
                }
            ).then(() => {
              this.connect();
            }).catch(() => {
              ElMessage({
                type: 'info',
                message: '您关闭了会话',
              })
            })
          }
        }).catch((res) => {
          if (res.code === 400) {
            this.showLoginDialog = true;
          } else {
            ElMessage.error(res.message)
          }
        })

      });

      this.socket = socket;
    },

    // 更换角色
    changeRole: function () {
      this.loading = true
      // 清空对话列表
      this.chatData = [];
      this.connect();
      for (const key in this.chatRoles) {
        if (this.chatRoles[key].key === this.role) {
          this.replyIcon = this.chatRoles[key].icon;
          break;
        }
      }
    },

    // 从后端获取聊天历史记录
    fetchChatHistory: function () {
      httpPost("/api/chat/history", {role: this.role}).then((res) => {
        const data = res.data
        const md = require('markdown-it')();
        for (let i = 0; i < data.length; i++) {
          if (data[i].type === "prompt") {
            this.chatData.push(data[i]);
            continue;
          }

          data[i].content = md.render(data[i].content);
          this.chatData.push(data[i]);
        }

        nextTick(() => {
          hl.configure({ignoreUnescapedHTML: true})
          const lines = document.querySelectorAll('.chat-line');
          const blocks = lines[lines.length - 1].querySelectorAll('pre code');
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
      this.socket.send(this.inputValue);
      this.$refs["text-input"].blur();
      this.inputValue = '';
      // 等待 textarea 重新调整尺寸之后再自动获取焦点
      setTimeout(() => this.$refs["text-input"].focus(), 100);
      return true;
    },

    // 获取焦点
    focus: function () {
      setTimeout(function () {
        document.getElementById('container').scrollTo(0, document.getElementById('container').scrollHeight)
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
        setSessionId(res.data)
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
    }
  },

})
</script>

<style lang="stylus">
#app {
  height: 100%;

  .body {
    background-color: rgba(247, 247, 248, 1);
    display flex;
    //justify-content center;
    align-items flex-start;
    height 100%;

    #container {
      overflow auto;
      width 100%;

      .tool-box {
        padding-top 10px;
        display flex;
        justify-content center;
        align-items center;

        .el-select {
          max-width 120px;
        }

        .chat-role {
          margin-left 5px;
        }

        .el-image {
          margin-right 5px;
        }

        .clear-history {
          margin-left 5px;
        }
      }

      .chat-box {
        // 变量定义
        --content-font-size: 16px;
        --content-color: #374151;

        font-family 'Microsoft YaHei', '微软雅黑', Arial, sans-serif;
        padding: 0 10px 10px 10px;

        .chat-line {
          padding 10px;
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

      .input-box {
        padding 10px;
        width 100%;

        position: absolute;
        bottom: 0
        display: flex;
        justify-content: start;
        align-items: center;

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
            width 106px;
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
    }

    #container::-webkit-scrollbar {
      width: 0;
      height: 0;
    }
  }
}

.el-message-box {
  width 90%;
  max-width 420px;
}

.el-message {
  min-width: 100px;
  max-width 600px;
}

.token-dialog {
  .el-dialog {
    --el-dialog-width 90%;
    max-width 400px;

    .el-dialog__body {
      padding 10px 10px 20px 10px;
    }

    .el-row {
      flex-wrap nowrap

      button {
        margin-left 5px;
      }
    }
  }
}

</style>
