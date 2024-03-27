<template>
  <div class="app-background">
    <div class="mobile-chat" v-loading="loading" element-loading-text="正在连接会话...">
      <van-sticky ref="navBarRef" :offset-top="0" position="top">
        <van-nav-bar left-arrow left-text="返回" @click-left="router.back()">
          <template #title>
            <van-dropdown-menu>
              <van-dropdown-item :title="title">
                <van-cell center title="角色"> {{ role.name }}</van-cell>
                <van-cell center title="模型">{{ modelValue }}</van-cell>
              </van-dropdown-item>
            </van-dropdown-menu>
          </template>

          <template #right>
            <van-icon name="share-o" @click="showShare = true"/>
          </template>

        </van-nav-bar>
      </van-sticky>

      <van-share-sheet
          v-model:show="showShare"
          title="立即分享给好友"
          :options="shareOptions"
          @select="shareChat"
      />

      <div class="chat-list-wrapper">
        <div id="message-list-box" :style="{height: winHeight + 'px'}" class="message-list-box">
          <van-list
              v-model:error="error"
              :finished="finished"
              error-text="请求失败，点击重新加载"
              @load="onLoad"
          >
            <van-cell v-for="item in chatData" :key="item" :border="false" class="message-line">
              <chat-prompt
                  v-if="item.type==='prompt'"
                  :content="item.content"
                  :created-at="dateFormat(item['created_at'])"
                  :icon="item.icon"
                  :model="model"
                  :tokens="item['tokens']"/>
              <chat-reply v-else-if="item.type==='reply'"
                          :content="item.content"
                          :created-at="dateFormat(item['created_at'])"
                          :icon="item.icon"
                          :org-content="item.orgContent"
                          :tokens="item['tokens']"/>
            </van-cell>
          </van-list>
        </div>
      </div>
      <div class="chat-box-wrapper">
        <van-sticky ref="bottomBarRef" :offset-bottom="0" position="bottom">

          <van-cell-group inset>
            <van-field
                v-model="prompt"
                center
                clearable
                placeholder="输入你的问题"
            >
              <template #left-icon>
                <van-button round type="success" class="button-voice" @click="inputVoice">
                  <el-icon>
                    <Microphone/>
                  </el-icon>
                </van-button>
              </template>

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
        </van-sticky>
      </div>
    </div>

    <button id="copy-link-btn" style="display: none;" :data-clipboard-text="url">复制链接地址</button>

    <van-overlay :show="showMic" z-index="100">
      <div class="mic-wrapper">
        <div class="image">
          <van-image
              width="100"
              height="100"
              src="/images/mic.gif"
          />
        </div>
        <van-button type="success" @click="stopVoice">说完了</van-button>
      </div>
    </van-overlay>
  </div>
</template>

<script setup>
import {nextTick, onMounted, onUnmounted, ref} from "vue";
import {showImagePreview, showNotify, showToast} from "vant";
import {onBeforeRouteLeave, useRouter} from "vue-router";
import {dateFormat, processContent, randString, renderInputText, UUID} from "@/utils/libs";
import {getChatConfig} from "@/store/chat";
import {httpGet} from "@/utils/http";
import hl from "highlight.js";
import 'highlight.js/styles/a11y-dark.css'
import ChatPrompt from "@/components/mobile/ChatPrompt.vue";
import ChatReply from "@/components/mobile/ChatReply.vue";
import {getSessionId, getUserToken} from "@/store/session";
import {checkSession} from "@/action/session";
import Clipboard from "clipboard";
import {Microphone} from "@element-plus/icons-vue";

const winHeight = ref(0)
const navBarRef = ref(null)
const bottomBarRef = ref(null)
const router = useRouter()

const chatConfig = getChatConfig()
const role = chatConfig.role
const model = chatConfig.model
const modelValue = chatConfig.modelValue
const title = chatConfig.title
const chatId = chatConfig.chatId
const loginUser = ref(null)
const showMic = ref(false)

const url = location.protocol + '//' + location.host + '/mobile/chat/export?chat_id=' + chatId

onMounted(() => {
  winHeight.value = document.body.offsetHeight - navBarRef.value.$el.offsetHeight - bottomBarRef.value.$el.offsetHeight

  const clipboard = new Clipboard(".content-mobile,.copy-code-mobile,#copy-link-btn");
  clipboard.on('success', (e) => {
    e.clearSelection()
    showNotify({type: 'success', message: '复制成功', duration: 1000})
  })
  clipboard.on('error', () => {
    showNotify({type: 'danger', message: '复制失败', duration: 2000})
  })
})

onUnmounted(() => {
  socket.value = null
})

const chatData = ref([])
const loading = ref(false)
const finished = ref(false)
const error = ref(false)

checkSession().then(user => {
  loginUser.value = user
}).catch(() => {
  router.push('/login')
})

const latexPlugin = require('markdown-it-latex2img')
const mathjaxPlugin = require('markdown-it-mathjax')
const md = require('markdown-it')({
  breaks: true,
  html: true,
  linkify: true,
  typographer: true,
  highlight: function (str, lang) {
    const codeIndex = parseInt(Date.now()) + Math.floor(Math.random() * 10000000)
    // 显示复制代码按钮
    const copyBtn = `<span class="copy-code-mobile" data-clipboard-action="copy" data-clipboard-target="#copy-target-${codeIndex}">复制</span>
<textarea style="position: absolute;top: -9999px;left: -9999px;z-index: -9999;" id="copy-target-${codeIndex}">${str.replace(/<\/textarea>/g, '&lt;/textarea>')}</textarea>`
    if (lang && hl.getLanguage(lang)) {
      const langHtml = `<span class="lang-name">${lang}</span>`
      // 处理代码高亮
      const preCode = hl.highlight(lang, str, true).value
      // 将代码包裹在 pre 中
      return `<pre class="code-container"><code class="language-${lang} hljs">${preCode}</code>${copyBtn} ${langHtml}</pre>`
    }

    // 处理代码高亮
    const preCode = md.utils.escapeHtml(str)
    // 将代码包裹在 pre 中
    return `<pre class="code-container"><code class="language-${lang} hljs">${preCode}</code>${copyBtn}</pre>`
  }
});
md.use(latexPlugin)
md.use(mathjaxPlugin)


const onLoad = () => {
  httpGet('/api/chat/history?chat_id=' + chatId).then(res => {
    // 加载状态结束
    finished.value = true;
    const data = res.data
    if (data && data.length > 0) {
      for (let i = 0; i < data.length; i++) {
        if (data[i].type === "prompt") {
          chatData.value.push(data[i]);
          continue;
        }

        data[i].orgContent = data[i].content;
        data[i].content = md.render(processContent(data[i].content))
        chatData.value.push(data[i]);
      }

      nextTick(() => {
        hl.configure({ignoreUnescapedHTML: true})
        const blocks = document.querySelector("#message-list-box").querySelectorAll('pre code');
        blocks.forEach((block) => {
          hl.highlightElement(block)
        })

        scrollListBox()
      })
    }

    // 连接会话
    connect(chatId, role.id);
  }).catch(() => {
    error.value = true
  })

};

// 离开页面时主动关闭 websocket 连接，节省网络资源
onBeforeRouteLeave(() => {
  if (socket.value !== null) {
    activelyClose.value = true;
    clearTimeout(heartbeatHandle.value)
    socket.value.close();
  }

})

// 创建 socket 连接
const prompt = ref('');
const showStopGenerate = ref(false); // 停止生成
const showReGenerate = ref(false); // 重新生成
const previousText = ref(''); // 上一次提问
const lineBuffer = ref(''); // 输出缓冲行
const socket = ref(null);
const activelyClose = ref(false); // 主动关闭
const canSend = ref(true);
const heartbeatHandle = ref(null)
const connect = function (chat_id, role_id) {
  let isNewChat = false;
  if (!chat_id) {
    isNewChat = true;
    chat_id = UUID();
  }

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

  // 心跳函数
  const sendHeartbeat = () => {
    if (socket.value !== null) {
      new Promise((resolve) => {
        socket.value.send(JSON.stringify({type: "heartbeat", content: "ping"}))
        resolve("success")
      }).then(() => {
        heartbeatHandle.value = setTimeout(() => sendHeartbeat(), 5000)
      });
    }
  }

  const _socket = new WebSocket(host + `/api/chat/new?session_id=${_sessionId}&role_id=${role_id}&chat_id=${chat_id}&model_id=${model}&token=${getUserToken()}`);
  _socket.addEventListener('open', () => {
    loading.value = false
    previousText.value = '';
    canSend.value = true;
    activelyClose.value = false;

    if (isNewChat) { // 加载打招呼信息
      chatData.value.push({
        type: "reply",
        id: randString(32),
        icon: role.icon,
        content: role.helloMsg,
        orgContent: role.helloMsg,
      })
    }

    // 发送心跳消息
    sendHeartbeat()
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
            icon: role.icon,
            content: ""
          });
        } else if (data.type === 'end') { // 消息接收完毕
          enableInput()
          lineBuffer.value = ''; // 清空缓冲

        } else {
          lineBuffer.value += data.content;
          const reply = chatData.value[chatData.value.length - 1]
          reply['orgContent'] = lineBuffer.value;
          reply['content'] = md.render(processContent(lineBuffer.value));

          nextTick(() => {
            hl.configure({ignoreUnescapedHTML: true})
            const lines = document.querySelectorAll('.message-line');
            const blocks = lines[lines.length - 1].querySelectorAll('pre code');
            blocks.forEach((block) => {
              hl.highlightElement(block)
            })
            scrollListBox()

            const items = document.querySelectorAll('.message-line')
            const imgs = items[items.length - 1].querySelectorAll('img')
            for (let i = 0; i < imgs.length; i++) {
              if (!imgs[i].src) {
                continue
              }
              imgs[i].addEventListener('click', (e) => {
                e.stopPropagation()
                showImagePreview([imgs[i].src]);
              })
            }
          })
        }

      };
    }

  });

  _socket.addEventListener('close', () => {
    if (activelyClose.value || socket.value === null) { // 忽略主动关闭
      return;
    }
    // 停止发送消息
    canSend.value = true;
    // 重连
    checkSession().then(() => {
      connect(chat_id, role_id)
    }).catch(() => {
      loading.value = true
      setTimeout(() => connect(chat_id, role_id), 3000)
    });
  });

  socket.value = _socket;
}

const disableInput = (force) => {
  canSend.value = false;
  showReGenerate.value = false;
  showStopGenerate.value = !force;
}

const enableInput = () => {
  canSend.value = true;
  showReGenerate.value = previousText.value !== "";
  showStopGenerate.value = false;
}

// 将聊天框的滚动条滑动到最底部
const scrollListBox = () => {
  document.getElementById('message-list-box').scrollTo(0, document.getElementById('message-list-box').scrollHeight + 46)
}

const sendMessage = () => {
  if (canSend.value === false) {
    showToast("AI 正在作答中，请稍后...");
    return
  }

  if (prompt.value.trim().length === 0) {
    showToast("请输入需要 AI 回答的问题")
    return false;
  }

  // 追加消息
  chatData.value.push({
    type: "prompt",
    id: randString(32),
    icon: loginUser.value.avatar,
    content: renderInputText(prompt.value),
    created_at: new Date().getTime(),
  });

  nextTick(() => {
    scrollListBox()
  })

  disableInput(false)
  socket.value.send(JSON.stringify({type: "chat", content: prompt.value}));
  previousText.value = prompt.value;
  prompt.value = '';
  return true;
}

const stopGenerate = () => {
  showStopGenerate.value = false;
  httpGet("/api/chat/stop?session_id=" + getSessionId()).then(() => {
    enableInput()
  })
}

const reGenerate = () => {
  disableInput(false)
  const text = '重新生成上述问题的答案：' + previousText.value;
  // 追加消息
  chatData.value.push({
    type: "prompt",
    id: randString(32),
    icon: loginUser.value.avatar,
    content: renderInputText(text)
  });
  socket.value.send(JSON.stringify({type: "chat", content: previousText.value}));
}

const showShare = ref(false)
const shareOptions = [
  {name: '微信', icon: 'wechat'},
  {name: '复制链接', icon: 'link'},
]
const shareChat = (option) => {
  showShare.value = false
  if (option.icon === "wechat") {
    showToast({message: "当前会话已经导出，请通过浏览器或者微信的自带分享功能分享给好友", duration: 5000})
    router.push({
      path: "/mobile/chat/export",
      query: {title: title, chat_id: chatId, role: role.name, model: modelValue}
    })
  } else if (option.icon === "link") {
    document.getElementById('copy-link-btn').click();
  }
}

// eslint-disable-next-line no-undef
const recognition = new webkitSpeechRecognition() || SpeechRecognition();
//recognition.lang = 'zh-CN' // 设置语音识别语言
recognition.onresult = function (event) {
  prompt.value = event.results[0][0].transcript
};

recognition.onerror = function (event) {
  showMic.value = false
  recognition.stop()
  showNotify({type: 'danger', message: '语音识别错误:' + event.error})
};

recognition.onend = function () {
  console.log('语音识别结束');
};
const inputVoice = () => {
  showMic.value = true
  recognition.start();
}

const stopVoice = () => {
  showMic.value = false
  recognition.stop()
}
</script>

<style lang="stylus">
@import "@/assets/css/mobile/chat-session.styl"
</style>