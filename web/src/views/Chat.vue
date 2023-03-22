<template>
  <div class="body" v-loading="loading">
    <div id="container">
      <div class="tool-box">
        <el-image style="width: 24px; height: 24px" :src="logo"/>
        <el-button round>欢迎来到人工智能时代</el-button>
      </div>

      <div class="chat-box" :style="{height: chatBoxHeight+'px'}">
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

      <div class="input-box">
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
          <el-input v-model="token" placeholder="在此输入口令">
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
import {Tools, Lock} from '@element-plus/icons-vue'
import ConfigDialog from '@/components/ConfigDialog.vue'
import {httpPost} from "@/utils/http";
import {getSessionId, setSessionId} from "@/utils/storage";

export default defineComponent({
  name: "XChat",
  components: {ChatPrompt, ChatReply, Tools, Lock, ConfigDialog},
  data() {
    return {
      title: 'ChatGPT 控制台',
      logo: 'images/logo.png',
      chatData: [],
      inputValue: '',
      chatBoxHeight: 0,
      showConnectDialog: false,
      showLoginDialog: false,
      token: '',

      connectingMessageBox: null,
      socket: null,
      toolBoxHeight: 61 + 42,
      sending: false,
      loading: false
    }
  },

  mounted: function () {
    nextTick(() => {
      this.chatBoxHeight = window.innerHeight - this.toolBoxHeight;
    })

    this.checkSession();

    for (let i = 0; i < 10; i++) {
      this.chatData.push({
        type: "prompt",
        id: randString(32),
        icon: 'images/user-icon.png',
        content: "孙悟空为什么可以把金棍棒放进耳朵？",
      });
      this.chatData.push({
        type: "reply",
        id: randString(32),
        icon: 'images/gpt-icon.png',
        content: "孙悟空是中国神话中的人物，传说中他可以把金箍棒放进耳朵里，这是一种超自然能力，无法用现代科学解释。这种能力可能是象征孙悟空超人力量的古代文化传说。",
      });
    }

    window.addEventListener("resize", () => {
      this.chatBoxHeight = window.innerHeight - this.toolBoxHeight;
    });
  },

  methods: {
    // 检查会话
    checkSession: function () {
      httpPost("/api/session/get").then(() => {
        if (this.socket == null) {
          this.connect();
        }
        // 发送心跳
        setTimeout(() => this.checkSession(), 5000);
      }).catch((res) => {
        if (res.code === 400) {
          this.showLoginDialog = true;
        } else {
          this.connectingMessageBox = ElMessageBox.confirm(
              '^_^ 会话发生异常，您已经从服务器断开连接!',
              '注意：',
              {
                confirmButtonText: '重连会话',
                cancelButtonText: '不聊了',
                type: 'warning',
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
      })
    },

    connect: function () {
      // 初始化 WebSocket 对象
      const token = getSessionId();
      const socket = new WebSocket('ws://' + process.env.VUE_APP_API_HOST + '/api/chat?token=' + token);
      socket.addEventListener('open', () => {
        ElMessage.success('创建会话成功！');

        if (this.connectingMessageBox != null) {
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
                icon: 'images/gpt-icon.png',
                content: "",
                cursor: true
              });
            } else if (data.type === 'end') {
              this.sending = false;
              this.chatData[this.chatData.length - 1]["cursor"] = false;
            } else {
              let content = data.content;
              // 替换换行符
              if (content.indexOf("\n\n") >= 0) {
                if (this.chatData[this.chatData.length - 1]["content"].length === 0) {
                  return
                }
                content = content.replace("\n\n", "<br />");
              }
              this.chatData[this.chatData.length - 1]["content"] += content;
            }
            // 将聊天框的滚动条滑动到最底部
            nextTick(() => {
              document.getElementById('container').scrollTo(0, document.getElementById('container').scrollHeight)
            })
          };
        }

      });
      socket.addEventListener('close', () => {
        // 检查会话，自动登录
        this.checkSession();
      });

      this.socket = socket;
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
        icon: 'images/user-icon.png',
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

        .el-image {
          margin-right 5px;
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
