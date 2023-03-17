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
                      :content="chat.content"/>
        </div>

      </div><!-- end chat box -->

      <div class="input-box" :style="{width: inputBoxWidth+'px'}" id="input-box">
        <div class="input-container">
          <textarea class="input-text" id="input-text" rows="1" :style="{minHeight:'24px', height: textHeight+'px'}"
                    v-on:keydown="inputKeyDown"
                    v-model="inputValue"
                    placeholder="Input any thing here..."
                    v-on:focus="focus"></textarea>
        </div>

        <div class="btn-container">
          <button type="button"
                  class="btn btn-success"
                  :disabled="sending"
                  v-on:click="sendMessage">发送
          </button>
        </div>

      </div><!-- end input box -->
    </div><!-- end container -->

  </div>
</template>

<script>
import {defineComponent, nextTick} from 'vue'
import ChatPrompt from "@/components/ChatPrompt.vue";
import ChatReply from "@/components/ChatReply.vue";
import {randString} from "@/utils/libs";

export default defineComponent({
  name: "XChat",
  components: {ChatPrompt, ChatReply},
  data() {
    return {
      title: "ChatGPT 控制台",
      chatData: [
        {
          id: "1",
          type: 'prompt',
          icon: 'images/user-icon.png',
          content: '请问棒球棒可以放进人的耳朵里面吗'
        },
        {
          id: "2",
          type: 'reply',
          icon: 'images/gpt-icon.png',
          content: '不可以。棒球棒的直径通常都比人的耳道大得多，而且人的耳朵是非常敏感和易受伤的，如果硬塞棒球棒可能会导致耳道损伤、出血和疼痛等问题。此外，塞入耳道的物体还可能引起耳屎的囤积和感染等问题，因此强烈建议不要将任何非耳朵医学用品的物品插入耳朵。如果您有耳道不适或者其他耳朵健康问题，应该咨询专业医生的建议。'
        }
      ],
      inputBoxHeight: 63,
      inputBoxWidth: 0,
      inputValue: '',
      textHeight: 24,
      textWidth: 0,
      chatBoxHeight: 0,
      isMobile: false,

      socket: null,
      sending: false
    }
  },

  computed: {},

  mounted: function () {
    nextTick(() => {
      this.inputBoxHeight = document.getElementById("input-box").offsetHeight;
      this.textWidth = document.getElementById("input-text").offsetWidth;
      this.chatBoxHeight = window.innerHeight - this.inputBoxHeight - 40;
    })

    //判断是否手机端访问
    const userAgentInfo = navigator.userAgent.toLowerCase();
    const Agents = ["android", "iphone", "windows phone", "ipad", "ipod"];
    for (let v = 0; v < Agents.length; v++) {
      if (userAgentInfo.toLowerCase().indexOf(Agents[v]) >= 0) {
        this.isMobile = true;
      }
    }

    this.inputBoxWidth = window.innerWidth;

    window.addEventListener('resize', this.windowResize);

    // 初始化 WebSocket 对象
    const socket = new WebSocket('ws://172.22.11.200:5678/api/chat');
    socket.addEventListener('open', () => {
      console.log('WebSocket 连接已打开');
    });
    socket.addEventListener('message', event => {
      if (event.data instanceof Blob) {
        const reader = new FileReader();
        reader.readAsText(event.data, "UTF-8");
        reader.onload = () => {
          // this.chatData.push({
          //   type: "reply",
          //   id: randString(32),
          //   icon: 'images/gpt-icon.png',
          //   content: reader.result
          // });
          this.chatData[this.chatData.length - 1]["content"] += reader.result
          this.sending = false;

          // 将聊天框的滚动条滑动到最底部
          nextTick(() => {
            document.getElementById('container').scrollTo(0, document.getElementById('container').scrollHeight)
          })
        };
      }

    });
    socket.addEventListener('close', event => {
      console.log('WebSocket 连接已关闭', event.reason);
    });
    socket.addEventListener('error', event => {
      console.error('WebSocket 连接发生错误', event);
    });

    this.socket = socket;

  },

  beforeUnmount() {
    window.removeEventListener("resize", this.windowResize);
  },

  methods: {
    inputKeyDown: function (e) {

      if (e.keyCode === 13) {
        if (!this.isMobile) { // PC 端按回车键直接提交数据
          e.preventDefault();
          return this.sendMessage();
        } else {
          return this.inputResize(true);
        }

      }
      this.inputResize(false);
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
      this.inputValue = '';
      this.inputResize(false);
      return true;
    },

    /**
     * 根据输入内容的多少动态调整输入框的大小
     * @param flag 是否输入回车键，如果输入了回车键则需要增加一行
     */
    inputResize: function (flag) {
      let line = 1;
      if (flag) {
        line++;
      }

      let textWidth = 0;
      for (let i in this.inputValue) {
        if (this.inputValue[i] === '\n') {
          line++;
          textWidth = 0; // 换行之后前面字数清零
          continue;
        }
        if (this.inputValue.charCodeAt(Number(i)) < 128) {
          textWidth += 8; // 英文字符宽度
        } else {
          textWidth += 16; // 中文字符宽度
        }

        if (textWidth >= this.textWidth) { // 另起一行
          textWidth = textWidth - this.textWidth;
          line++;
        }
      }

      this.inputBoxHeight = 63 + (line - 1) * 24;
      this.textHeight = line * 24;
    },

    windowResize: function () {
      this.inputResize(false);
      this.chatBoxHeight = window.innerHeight - this.inputBoxHeight - 40;
      this.inputBoxWidth = window.innerWidth;
    },

    // 获取焦点
    focus: function () {
      setTimeout(function () {
        document.getElementById('container').scrollTo(0, document.getElementById('container').scrollHeight)
      }, 200)
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
    justify-content center;
    align-items flex-start;
    height 100%;

    #container {
      overflow auto;

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

        position: absolute;
        bottom: 0
        display: flex;
        justify-content: center;
        align-items: flex-start;

        .input-container {
          overflow hidden
          width 100%
          margin: 0;
          border: none;
          border-radius: 6px;
          box-shadow: 0 2px 15px rgba(0, 0, 0, 0.1);
          background-color: rgba(255, 255, 255, 1);
          padding: 5px 10px;

          .input-text {
            font-size: 16px;
            padding 0
            margin 0
            outline: none;
            width 100%;
            border none
            background #ffffff
            resize none
            line-height 24px;
            color #333;
          }

          .input-text::-webkit-scrollbar {
            width: 0;
            height: 0;
          }
        }

        .btn-container {
          margin-left 10px;

          button {
            width 70px;
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


</style>
