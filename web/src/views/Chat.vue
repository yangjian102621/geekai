<template>
  <div class="body">
    <div id="container">
      <div class="chat-box" :style="{height: chatBoxHeight+'px'}">
        <div v-for="chat in chatData" :key="chat.id">
          <chat-prompt
              v-if="chat.type==='prompt'"
              :icon="chat.icon"
              :content="chat.content"/>
          <chat-reply v-else-if="chat.type==='reply'"
                      :icon="chat.icon"
                      :cursor="chat.cursor"
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
              placeholder="Input any thing here..."
          />
        </div>

        <div class="btn-container">
          <el-row>
            <el-button type="success" class="send" :disabled="sending" v-on:click="sendMessage">发送</el-button>
            <el-button type="info" class="config" circle @click="showDialog = true">
              <el-icon>
                <Tools/>
              </el-icon>
            </el-button>
          </el-row>
        </div>

      </div><!-- end input box -->

    </div><!-- end container -->

    <config-dialog v-model:show="showDialog"></config-dialog>
  </div>
</template>

<script>
import {defineComponent, nextTick} from 'vue'
import ChatPrompt from "@/components/ChatPrompt.vue";
import ChatReply from "@/components/ChatReply.vue";
import {randString} from "@/utils/libs";
import {ElMessage, ElMessageBox} from 'element-plus'
import {Tools} from '@element-plus/icons-vue'
import ConfigDialog from '@/components/ConfigDialog.vue'

export default defineComponent({
  name: "XChat",
  components: {ChatPrompt, ChatReply, Tools, ConfigDialog},
  data() {
    return {
      title: "ChatGPT 控制台",
      chatData: [],
      inputValue: '',
      chatBoxHeight: 0,
      showDialog: false,

      connectingMessageBox: null,
      socket: null,
      sending: false
    }
  },

  computed: {},

  mounted: function () {
    nextTick(() => {
      this.chatBoxHeight = window.innerHeight - 61;
    })
    this.connect();

  },

  methods: {
    connect: function () {
      if (this.online) {
        return
      }

      // 初始化 WebSocket 对象
      const socket = new WebSocket(process.env.VUE_APP_WS_HOST + '/api/chat');
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
        ElMessageBox.confirm(
            '^_^ 会话发生异常，您已经从服务器断开连接!',
            '注意：',
            {
              confirmButtonText: '重连会话',
              cancelButtonText: '不聊了',
              type: 'warning',
            }
        )
            .then(() => {
              this.connect();
            })
            .catch(() => {
              ElMessage({
                type: 'info',
                message: '您关闭了会话',
              })
            })
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
    sendMessage: function () {
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

      // TODO: 使用 websocket 提交数据到后端
      this.sending = true;
      this.socket.send(this.inputValue);
      this.$refs["text-input"].blur();
      this.inputValue = '';
      // 等待 textarea 重新调整尺寸之后再自动获取焦点
      setTimeout(() => this.$refs["text-input"].focus(), 100)
      return true;
    },

    // 获取焦点
    focus: function () {
      setTimeout(function () {
        document.getElementById('container').scrollTo(0, document.getElementById('container').scrollHeight)
      }, 200)
    },

  },

})
</script>

<style lang="stylus">
#app {
  height: 100%;

  .body {
    background-color: rgba(247, 247, 248, 1);
    display flex;
    justify-content center;
    align-items flex-start;
    height 100%;

    #container {
      overflow auto;
      width 100%;

      .chat-box {
        // 变量定义
        --content-font-size: 16px;
        --content-color: #374151;

        font-family 'Microsoft YaHei', '微软雅黑', Arial, sans-serif;
        padding: 20px 10px;

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
  width 90%;
  min-width: 300px;
  max-width 600px;
}


</style>
