<template>
  <div class="body-plus">
    <el-row>
      <div class="chat-head">
        <el-row class="row-center">
          <el-col :span="12">
            <div class="title-box">
              <el-image :src="logo" class="logo"/>
              <span>ChatGPT-Plus</span>
            </div>
          </el-col>
          <el-col :span="12">
            <div class="tool-box">

              <el-dropdown :hide-on-click="true" class="user-info" trigger="click">
                <span class="el-dropdown-link">
                  <el-image src="images/user-info.jpg"/>
                  <el-icon><ArrowDown/></el-icon>
                </span>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item @click="showConnectDialog = true">
                      <el-icon>
                        <Tools/>
                      </el-icon>
                      <span>聊天设置</span>
                    </el-dropdown-item>

                    <el-dropdown-item @click="clearChatHistory">
                      <el-icon>
                        <Delete/>
                      </el-icon>
                      <span>删除记录</span>
                    </el-dropdown-item>

                    <el-dropdown-item @click="logout">
                      <el-icon>
                        <Monitor/>
                      </el-icon>
                      <span>注销</span>
                    </el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
            </div>
          </el-col>
        </el-row>
      </div>
    </el-row>
    <el-row>
      <div class="left-box">
        <div class="search-box">
          <el-input v-model="roleName" class="w-50 m-2" size="small" placeholder="搜索聊天角色" @keyup="searchRole">
            <template #prefix>
              <el-icon class="el-input__icon">
                <Search/>
              </el-icon>
            </template>
          </el-input>
        </div>

        <div class="content" :style="{height: leftBoxHeight+'px'}">
          <el-row v-for="item in chatRoles" :key="item.key">
            <div :class="item.key === this.role?'chat-role-item active':'chat-role-item'" @click="changeRole(item)">
              <el-image :src="item.icon" class="avatar"/>
              <span>{{ item.name }}</span>
            </div>
          </el-row>
        </div>
      </div>
      <div class="right-box" :style="{height: mainWinHeight+'px'}">
        <div v-loading="loading" element-loading-background="rgba(122, 122, 122, 0.8)">
          <div id="container">
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

            <div class="re-generate">
              <div class="btn-box">
                <el-button type="info" v-if="showStopGenerate" @click="stopGenerate" plain>
                  <el-icon>
                    <VideoPause/>
                  </el-icon>
                  停止生成
                </el-button>

                <el-button type="info" v-if="showReGenerate" @click="reGenerate" plain>
                  <el-icon>
                    <RefreshRight/>
                  </el-icon>
                  重新生成
                </el-button>
              </div>
            </div>

            <el-row class="chat-tool-box">
              <el-tooltip
                  class="box-item"
                  effect="dark"
                  content="进入AI绘画模式"
                  placement="top"
              >
                <el-icon @click="drawImage">
                  <Picture/>
                </el-icon>
              </el-tooltip>

            </el-row>

            <div class="input-box">
              <div class="input-container">
                <el-input
                    ref="text-input"
                    v-model="inputValue"
                    :autosize="{ minRows: 5, maxRows: 10 }"
                    v-on:keydown="inputKeyDown"
                    v-on:focus="focus"
                    autofocus
                    type="textarea"
                    placeholder="先聊五毛钱吧..."
                />
              </div>
            </div><!-- end input box -->

          </div><!-- end container -->
        </div><!-- end loading -->
      </div>
    </el-row>

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

        <el-row class="row-center">
          <p>打开微信扫下面二维码免费领取口令</p>
        </el-row>

        <el-row class="row-center">
          <el-image src="images/wx.png" fit="cover"/>
        </el-row>

      </el-dialog>
    </div> <!--end token dialog-->
  </div>


</template>

<script>
import {defineComponent, nextTick} from 'vue'
import ChatPrompt from "@/components/plus/ChatPrompt.vue";
import ChatReply from "@/components/plus/ChatReply.vue";
import {isMobile, randString} from "@/utils/libs";
import {ElMessage, ElMessageBox} from 'element-plus'
import {
  Tools,
  Lock,
  Delete,
  Picture,
  Search,
  ArrowDown,
  Monitor,
  VideoPause,
  RefreshRight
} from '@element-plus/icons-vue'
import ConfigDialog from '@/components/ConfigDialog.vue'
import {httpPost, httpGet} from "@/utils/http";
import {getSessionId, setLoginUser} from "@/utils/storage";
import hl from 'highlight.js'
import 'highlight.js/styles/a11y-dark.css'

export default defineComponent({
  name: "ChatPlus",
  components: {
    RefreshRight,
    VideoPause,
    ArrowDown,
    Search,
    ChatPrompt,
    ChatReply,
    Tools,
    Lock,
    Delete,
    Picture,
    Monitor,
    ConfigDialog
  },
  data() {
    return {
      title: 'ChatGPT 控制台',
      logo: 'images/logo.png',
      chatData: [],
      chatRoles: [], // 当前显示的角色集合
      allChatRoles: [], // 所有角色集合
      role: 'gpt',
      inputValue: '', // 聊天内容
      showConnectDialog: false,
      showLoginDialog: false,
      token: '', // 会话 token
      replyIcon: 'images/avatar/gpt.png', // 回复信息的头像
      roleName: "", // 搜索角色名称

      showStopGenerate: false, // 停止生成
      showReGenerate: false, // 重新生成
      canReGenerate: false, // 是否可以重新生
      previousText: '', // 上一次提问

      lineBuffer: '', // 输出缓冲行
      connectingMessageBox: null, // 保存重连的消息框对象
      errorMessage: null, // 错误信息提示框
      socket: null,
      mainWinHeight: 0, // 主窗口高度
      chatBoxHeight: 0, // 聊天内容框高度
      leftBoxHeight: 0,
      sending: true,
      loading: true
    }
  },

  mounted: function () {
    if (isMobile()) {
      this.$router.push("mobile");
      return;
    }

    nextTick(() => {
      this.resizeElement();
    })
    window.addEventListener("resize", () => {
      this.resizeElement();
    });

    this.connect();
  },

  methods: {
    resizeElement: function () {
      this.chatBoxHeight = window.innerHeight - 61 - 115 - 38;
      this.mainWinHeight = window.innerHeight - 61;
      this.leftBoxHeight = window.innerHeight - 61 - 100;
    },
    // 创建 socket 会话连接
    connect: function () {
      // 初始化 WebSocket 对象
      const sessionId = getSessionId();
      const socket = new WebSocket(process.env.VUE_APP_WS_HOST + `/api/chat?sessionId=${sessionId}&role=${this.role}`);
      socket.addEventListener('open', () => {
        // 获取聊天角色
        if (this.chatRoles.length === 0) {
          httpGet("/api/config/chat-roles/get").then((res) => {
            // ElMessage.success('创建会话成功！');
            this.chatRoles = res.data;
            this.allChatRoles = res.data;
            this.loading = false
          }).catch(() => {
            ElMessage.error("获取聊天角色失败");
          })
        } else {
          this.loading = false
        }

        this.sending = false; // 允许用户发送消息
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
            } else if (data.type === 'end') { // 消息接收完毕
              this.sending = false;
              if (data['is_hello_mgs'] !== true) {
                this.showReGenerate = true;
              }
              this.showStopGenerate = false;
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

    drawImage: function () {
      ElMessage({
        message: '客观别急，AI 绘画服服务正在紧锣密鼓搭建中...',
        type: 'info',
      })
    },

    // 更换角色
    changeRole: function (item) {
      this.loading = true
      this.role = item.key;
      this.replyIcon = item.icon;
      // 清空对话列表
      this.chatData = [];
      this.connect();
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
          ElMessage.warning("AI 正在作答中，请稍后...");
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

      if (this.inputValue.trim().length === 0 || this.sending) {
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
        this.token = '';
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

    // 搜索聊天角色
    searchRole: function () {
      if (this.roleName === '') {
        this.chatRoles = this.allChatRoles;
        return;
      }
      const roles = [];
      for (let i = 0; i < this.allChatRoles.length; i++) {
        if (this.allChatRoles[i].name.indexOf(this.roleName) !== -1) {
          roles.push(this.allChatRoles[i]);
        }
      }
      this.chatRoles = roles;
    },

    // 退出登录
    logout: function () {
      httpPost("/api/logout", {opt: "logout"}).then(() => {
        this.checkSession();
      }).catch(() => {
        ElMessage.error("注销失败");
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

  .body-plus {
    height 100%;

    .chat-head {
      width 100%;
      height 60px;
      background-color: #28292A
      border-bottom 1px solid #4f4f4f;

      .title-box {
        padding-top 6px;
        display flex
        color #ffffff;
        font-size 20px;

        .logo {
          background-color #ffffff
          border-radius 50%;
          width 45px;
          height 45px;
        }

        span {
          padding-top: 12px;
          padding-left: 10px;
        }
      }

      .tool-box {
        padding-top 16px;
        padding-right 20px;
        display flex;
        justify-content flex-end;
        align-items center;

        .user-info {
          margin-left 20px;

          .el-dropdown-link {
            cursor pointer

            img {
              width 30px;
              height 30px;
              border-radius 50%;
            }

            .el-icon {
              bottom 8px
              color #cccccc
              margin-left 5px;
            }
          }
        }
      }

    }

    .el-row {
      overflow hidden;
      display: flex;

      .left-box {
        display flex
        flex-flow column
        min-width 220px;
        max-width 250px;
        background-color: #28292A
        border-top: 1px solid #2F3032
        border-right: 1px solid #2F3032

        .search-box {
          flex-wrap wrap
          padding 10px 15px;

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
          width 100%
          overflow-y scroll

          .chat-role-item {
            display flex
            width 100%
            justify-content flex-start
            padding 10px 18px
            border-bottom: 1px solid #3c3c3c
            cursor pointer


            .avatar {
              width 36px;
              height 36px;
              border-radius 50%;
            }

            span {
              color #c1c1c1
              padding 8px 10px;
            }
          }


          .chat-role-item.active {
            background-color: #363535;
          }
        }
      }

      .right-box {
        min-width: 0;
        flex: 1;
        background-color #232425
        border-left 1px solid #4f4f4f
      }
    }


    #container {
      overflow hidden;
      width 100%;

      ::-webkit-scrollbar {
        width: 0;
        height: 0;
        background-color: transparent;
      }

      .chat-box {
        overflow-y: scroll;
        border-bottom 1px solid #4f4f4f

        // 变量定义
        --content-font-size: 16px;
        --content-color: #c1c1c1;

        font-family 'Microsoft YaHei', '微软雅黑', Arial, sans-serif;
        padding: 10px;

        .chat-line {
          padding 10px 5px;
          font-size 14px;
          display: flex;
          align-items: flex-start;

          .chat-icon {
            img {
              width 45px;
              height 45px;
              border 1px solid #666;
              border-radius 50%;
              padding 1px;
            }
          }
        }
      }

      .re-generate {
        position: relative;
        display: flex;
        justify-content: center;

        .btn-box {
          position absolute
          bottom 10px;

          .el-button {
            .el-icon {
              margin-right 5px;
            }
          }

        }
      }

      .chat-tool-box {
        padding 10px;
        border-top: 1px solid #2F3032

        .el-icon svg {
          color #cccccc
          width 1em;
          background-color #232425
          cursor pointer
        }
      }

      .input-box {
        background-color #232425
        display: flex;
        justify-content: flex-start;
        align-items: center;

        .input-container {
          width: 100%
          margin: 0;
          border: none;
          border-radius: 6px;
          box-shadow: 0 2px 15px rgba(0, 0, 0, 0.1);
          background-color #232425
          padding: 5px 10px;

          .el-textarea__inner {
            box-shadow: none
            padding 5px 0
            background-color #232425
            color #B5B7B8
          }

          .el-textarea__inner::-webkit-scrollbar {
            width: 0;
            height: 0;
          }
        }
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
}

</style>
