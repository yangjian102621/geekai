<template>
  <div class="common-layout">
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

    </div>

    <div class="input-box" :style="{width: inputBoxWidth+'px'}" id="input-box">
      <div class="input-container">
        <textarea class="input-text" id="input-text" rows="1" :style="{minHeight:'24px', height: textHeight+'px'}"
                  v-on:keydown="inputKeyDown"
                  v-model="inputValue"
                  placeholder="Input any thing here..."
                  autofocus></textarea>
      </div>
      <div class="btn-container">
        <button type="button"
                class="btn btn-success"
                :disabled="sending"
                v-on:click="sendMessage">发送
        </button>
      </div>
    </div>
  </div>
</template>

<script>
import {defineComponent} from 'vue'
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
      inputBoxWidth: window.innerWidth - 20,
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
    this.inputBoxHeight = document.getElementById("input-box").offsetHeight;
    this.textWidth = document.getElementById("input-text").offsetWidth;
    this.chatBoxHeight = window.innerHeight - this.inputBoxHeight - 40;

    //判断是否手机端访问
    const userAgentInfo = navigator.userAgent.toLowerCase();
    const Agents = ["android", "iphone", "windows phone", "ipad", "ipod"];
    for (let v = 0; v < Agents.length; v++) {
      if (userAgentInfo.indexOf(Agents[v]) >= 0) {
        this.isMobile = true;
      }
    }

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
          this.chatData.push({
            type: "reply",
            id: randString(32),
            icon: 'images/gpt-icon.png',
            content: reader.result
          });
          this.sending = false;
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
      // PC 端按回车键直接提交数据
      if (e.keyCode === 13 && !this.isMobile) {
        e.stopPropagation();
        e.preventDefault();
        return this.sendMessage();
      }

      this.inputResize();
    },

    // 发送消息
    sendMessage: function () {
      if (this.sending) {
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
      return true;
    },

    // 初始化
    inputResize: function () {
      // 根据输入的字数自动调整输入框的大小
      let line = 1;
      let textWidth = 0;
      for (let i in this.inputValue) {
        if (this.inputValue[i] === '\n') {
          line++;
          textWidth = 0; // 换行之后前面字数清零
        }
        if (this.inputValue.charCodeAt(Number(i)) < 128) {
          textWidth += 9.65; // 英文字符宽度
        } else {
          textWidth += 16.07; // 中文字符宽度
        }
      }
      line = line + (Math.ceil(textWidth / this.textWidth)) - 1;
      this.inputBoxHeight = 63 + (line - 1) * 24;
      this.textHeight = line * 24;
    },

    windowResize: function () {
      this.inputResize();
      this.chatBoxHeight = window.innerHeight - this.inputBoxHeight - 40;
      this.inputBoxWidth = window.innerWidth - 20;
    }
  },

})
</script>

<style lang="stylus">
.chat-box {
  // 变量定义
  --content-font-size: 16px;
  --content-color: #374151;

  background-color: rgba(247, 247, 248, 1);
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
  background-color: rgba(255, 255, 255, 1);
  padding 10px;

  position: absolute;
  bottom: 0
  display: flex;
  justify-content: center;
  align-items: center;

  .input-container {
    overflow hidden
    width 100%
    margin: 0;
    background #ffffff
    border: none;
    border-radius: 6px;
    box-shadow: 0 2px 15px rgba(0, 0, 0, 0.1);
    padding: 5px 10px;

    .input-text {
      font-size: 16px;
      padding 0
      margin 0
      outline: none;
      width 100%;
      border none
      background transparent
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
</style>
