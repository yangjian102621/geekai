<template>
  <div class="mobile-chat">
    <van-sticky :offset-top="0" position="top" ref="navBarRef">
      <van-nav-bar left-text="返回" left-arrow @click-left="router.back()">
        <template #title>
          <van-dropdown-menu>
            <van-dropdown-item :title="title">
              <van-cell center title="角色"> {{ role.name }}</van-cell>
              <van-cell center title="模型">{{ model }}</van-cell>
            </van-dropdown-item>
          </van-dropdown-menu>
        </template>

        <template #right>
          <van-icon name="delete-o" @click="clearChatHistory"/>
        </template>
      </van-nav-bar>
    </van-sticky>


    <div class="message-list-box" id="message-list-box" :style="{height: winHeight+'px'}">
      <van-list
          v-model:loading="loading"
          :finished="finished"
          v-model:error="error"
          error-text="请求失败，点击重新加载"
          @load="onLoad"
      >
        <van-cell v-for="item in chatData" :key="item">
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
        </van-cell>
      </van-list>
    </div>

    <van-sticky :offset-bottom="0" position="bottom" ref="bottomBarRef">
      <div class="chat-box">
        <van-cell-group>
          <van-field
              v-model="prompt"
              center
              clearable
              placeholder="输入你的问题"
          >
            <template #button>
              <van-button size="small" type="primary" @click="sendMessage">发送</van-button>
            </template>
            <template #extra>
              <div class="icon-box">
                <van-icon v-if="showStopGenerate" name="stop-circle-o" @click="stopGenerate"/>
                <van-icon v-if="showReGenerate" name="play-circle-o" @click="reGenerate"/>
              </div>
            </template>
          </van-field>
        </van-cell-group>
      </div>
    </van-sticky>
  </div>

</template>

<script setup>
import {nextTick, onMounted, ref} from "vue";
import {showToast} from "vant";
import {useRouter} from "vue-router";
import {dateFormat, UUID} from "@/utils/libs";
import {getChatConfig} from "@/store/chat";
import {httpGet} from "@/utils/http";
import hl from "highlight.js";
import 'highlight.js/styles/a11y-dark.css'
import ChatPrompt from "@/components/mobile/ChatPrompt.vue";
import ChatReply from "@/components/mobile/ChatReply.vue";

const winHeight = ref(0)
const navBarRef = ref(null)
const bottomBarRef = ref(null)
const router = useRouter()

const chatConfig = getChatConfig()
const role = chatConfig.role
const model = chatConfig.model
const title = chatConfig.title
const chatId = chatConfig.chatId

onMounted(() => {
  winHeight.value = document.body.offsetHeight - navBarRef.value.$el.offsetHeight - bottomBarRef.value.$el.offsetHeight
})

const chatData = ref([])
const loading = ref(false)
const finished = ref(false)
const error = ref(false)
const onLoad = () => {
  httpGet('/api/chat/history?chat_id=' + chatId).then(res => {
    // 加载状态结束
    loading.value = false;
    finished.value = true;
    const data = res.data
    if (!data || data.length === 0) {
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
      const blocks = document.querySelector("#message-list-box").querySelectorAll('pre code');
      blocks.forEach((block) => {
        hl.highlightElement(block)
      })
    })
  }).catch(() => {
    error.value = true
  })

};

// 创建 socket 连接
const prompt = ref('');
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
      host = 'wss://' + location.host;
    } else {
      host = 'ws://' + location.host;
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

const clearChatHistory = () => {
  showToast('清空聊记录')
}

const sendMessage = () => {
  showToast("发送成功")
}

const stopGenerate = () => {
  showToast("停止生成")
}

const reGenerate = () => {
  showToast('重新生成')
}
</script>

<style scoped lang="stylus">
.mobile-chat {
  .message-list-box {
    overflow-x auto
  }

  .chat-box {
    .icon-box {
      .van-icon {
        font-size 24px
        margin-left 10px;
      }
    }
  }
}
</style>

<style lang="stylus">
.mobile-chat {
  .van-nav-bar__title {
    .van-dropdown-menu__title {
      margin-right 10px
    }

    .van-cell__title {
      text-align left
    }
  }

  .van-nav-bar__right {
    .van-icon {
      font-size 20px
    }
  }
}
</style>