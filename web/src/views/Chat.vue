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
          >
            <div class="role-option">
              <el-image :src="item.icon"></el-image>
              <span>{{ item.name }}</span>
            </div>
          </el-option>
        </el-select>

        <el-button type="danger" class="clear-history" size="small" circle @click="clearChatHistory">
          <el-icon>
            <Delete/>
          </el-icon>
        </el-button>

        <el-button type="info" size="small" class="config" ref="send-btn" circle @click="configDialog">
          <el-icon>
            <Tools/>
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
                      :org-content="chat.orgContent"
                      :content="chat.content"/>
        </div>

      </div><!-- end chat box -->

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
            </el-row>
          </div>
        </div>

      </div><!-- end input box -->

    </div><!-- end container -->

    <config-dialog v-model:show="showConfigDialog" :user="userInfo"></config-dialog>

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
          <el-image :src="sysConfig['wechat_card']" fit="cover" style="width: 250px;"/>
        </el-row>

      </el-dialog>
    </div>
  </div>
</template>

<script>
import {defineComponent, nextTick} from 'vue'
import ChatPrompt from "@/components/ChatPrompt.vue";
import ChatReply from "@/components/ChatReply.vue";
import {isMobile, randString, renderInputText} from "@/utils/libs";
import {ElMessage, ElMessageBox} from 'element-plus'
import {Tools, Lock, Delete, VideoPause, RefreshRight} from '@element-plus/icons-vue'
import ConfigDialog from '@/components/ConfigDialog.vue'
import {httpPost, httpGet} from "@/utils/http";
import {getSessionId, getUserInfo, setLoginUser} from "@/utils/storage";
import hl from 'highlight.js'
import 'highlight.js/styles/a11y-dark.css'
import Clipboard from "clipboard";

export default defineComponent({
  name: "XChat",
  components: {RefreshRight, VideoPause, ChatPrompt, ChatReply, Tools, Lock, Delete, ConfigDialog},
  data() {
    return {
      logo: 'images/logo.png',
      chatData: [],
      chatRoles: [],
      role: 'gpt',
      inputValue: '', // 聊天内容
      chatBoxHeight: 0, // 聊天内容框高度

      showConfigDialog: false,
      userInfo: {},
      showLoginDialog: false,
      sysConfig: {}, // 系统配置
      hasHelloMsg: {}, // 是否发送过打招呼信息

      showStopGenerate: false,
      showReGenerate: false,
      canReGenerate: false, // 是否可以重新生
      previousText: '', // 上一次提问

      token: '', // 会话 token
      replyIcon: 'images/avatar/gpt.png', // 回复信息的头像

      lineBuffer: '', // 输出缓冲行
      connectingMessageBox: null, // 保存重连的消息框对象
      errorMessage: null, // 错误信息提示框
      socket: null,
      toolBoxHeight: 61 + 52, // 工具框的高度
      inputBoxWidth: window.innerWidth - 20,
      sending: true,
      loading: true
    }
  },

  mounted: function () {
    if (!isMobile()) {
      this.$router.push("plus");
      return;
    }

    const clipboard = new Clipboard('.copy-reply');
    clipboard.on('success', () => {
      ElMessage.success('复制成功！');
    })

    clipboard.on('error', () => {
      ElMessage.error('复制失败！');
    })

    nextTick(() => {
      this.chatBoxHeight = window.innerHeight - this.toolBoxHeight;
      ElMessage.warning("强烈建议使用PC浏览器访问获的更好的聊天体验！")
    })

    window.addEventListener("resize", () => {
      this.chatBoxHeight = window.innerHeight - this.toolBoxHeight;
      this.inputBoxWidth = window.innerWidth - 20;
    });

    // 获取系统配置
    httpGet('/api/config/get').then((res) => {
      this.sysConfig = res.data;
    }).catch(() => {
      ElMessage.error('获取系统配置失败')
    })

    this.connect();
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
        // 获取聊天角色
        if (this.chatRoles.length === 0) {
          httpGet("/api/chat-roles/list").then((res) => {
            // ElMessage.success('创建会话成功！');
            this.chatRoles = res.data;
            this.loading = false
          }).catch(() => {
            ElMessage.error("获取聊天角色失败");
          })
        } else {
          this.loading = false
        }

        this.sending = false; // 允许用户发送消息
        this.activelyClose = false;
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
            if (data['is_hello_msg'] && this.hasHelloMsg[this.role]) { // 一定发送过打招呼信息的
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
              } else {
                this.hasHelloMsg[this.role] = true
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
              document.getElementById('container').scrollTo(0, document.getElementById('container').scrollHeight)
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

    // 更换角色
    changeRole: function () {
      this.loading = true
      // 清空对话列表
      this.chatData = [];
      this.hasHelloMsg = {};
      this.showStopGenerate = false;
      this.showReGenerate = false;
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
        content: renderInputText(this.inputValue)
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
    }
  },

})
</script>

<style lang="stylus">
#app {
  height: 100%;

  .body {
    background-color: rgba(247, 247, 248, 1);

    background-image url("~@/assets/img/bg_01.jpeg")
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

        .clear-history, .config {
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
    }

    #container::-webkit-scrollbar {
      width: 0;
      height: 0;
    }

    .row-center {
      justify-content center
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

      .el-row {
        flex-wrap nowrap

        button {
          margin-left 5px;
        }
      }

      .tip-text {
        text-align left
        padding 20px;

        .el-alert {
          padding 5px;

          .el-alert__description {
            font-size 14px;
          }
        }

      }
    }


  }
}

.el-select-dropdown {
  .el-select-dropdown__item {
    padding 8px 5px;

    .role-option {
      display flex
      flex-flow row

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
