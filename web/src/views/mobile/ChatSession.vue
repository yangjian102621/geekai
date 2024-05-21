<template>
  <div class="app-background">
    <div class="mobile-chat" v-loading="loading" element-loading-text="正在连接会话...">
      <van-nav-bar ref="navBarRef">
        <template #title>
          <van-dropdown-menu>
            <van-dropdown-item :title="title">
              <van-cell center title="角色"> {{ role.name }}</van-cell>
              <van-cell center title="模型">{{ modelValue }}</van-cell>
            </van-dropdown-item>
          </van-dropdown-menu>
        </template>
        <template #left>
          <span class="setting">
            <van-icon name="setting-o" @click="showPicker = true"/>
          </span>
        </template>
        <template #right>
          <van-icon name="share-o" @click="showShare = true"/>
        </template>

      </van-nav-bar>

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
                  :icon="item.icon"/>
              <chat-reply v-else-if="item.type==='reply'"
                          :content="item.content"
                          :icon="item.icon"
                          :org-content="item.orgContent"/>
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
              <!--              <template #left-icon>-->
              <!--                <van-button round type="success" class="button-voice" @click="inputVoice">-->
              <!--                  <el-icon>-->
              <!--                    <Microphone/>-->
              <!--                  </el-icon>-->
              <!--                </van-button>-->
              <!--              </template>-->

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

    <!--    <van-overlay :show="showMic" z-index="100">-->
    <!--      <div class="mic-wrapper">-->
    <!--        <div class="image">-->
    <!--          <van-image-->
    <!--              width="100"-->
    <!--              height="100"-->
    <!--              src="/images/mic.gif"-->
    <!--          />-->
    <!--        </div>-->
    <!--        <van-button type="success" @click="stopVoice">说完了</van-button>-->
    <!--      </div>-->
    <!--    </van-overlay>-->
  </div>

  <van-popup v-model:show="showPicker" position="bottom" class="popup">
    <van-picker
        :columns="columns"
        title="选择模型和角色"
        @cancel="showPicker = false"
        @confirm="newChat"
    >
      <template #option="item">
        <div class="picker-option">
          <van-image
              v-if="item.icon"
              :src="item.icon"
              fit="cover"
              round
          />
          <span>{{ item.text }}</span>
        </div>
      </template>
    </van-picker>
  </van-popup>
</template>

<script setup>
import {nextTick, onMounted, onUnmounted, ref} from "vue";
import {showImagePreview, showNotify, showToast} from "vant";
import {onBeforeRouteLeave, useRouter} from "vue-router";
import {processContent, randString, renderInputText, UUID} from "@/utils/libs";
import {httpGet} from "@/utils/http";
import hl from "highlight.js";
import 'highlight.js/styles/a11y-dark.css'
import ChatPrompt from "@/components/mobile/ChatPrompt.vue";
import ChatReply from "@/components/mobile/ChatReply.vue";
import {getSessionId, getUserToken} from "@/store/session";
import {checkSession} from "@/action/session";
import Clipboard from "clipboard";
import {showLoginDialog} from "@/utils/dialog";

const winHeight = ref(0)
const navBarRef = ref(null)
const bottomBarRef = ref(null)
const router = useRouter()

const roles = ref([])
const roleId = ref(parseInt(router.currentRoute.value.query["role_id"]))
const role = ref({})
const models = ref([])
const modelId = ref(parseInt(router.currentRoute.value.query["model_id"]))
const modelValue = ref("")
const title = ref(router.currentRoute.value.query["title"])
const chatId = ref(router.currentRoute.value.query["chat_id"])
const loginUser = ref(null)
// const showMic = ref(false)
const showPicker = ref(false)
const columns = ref([roles.value, models.value])

checkSession().then(user => {
  loginUser.value = user
}).catch(() => {
  router.push('/login')
})

if (chatId.value) {
  httpGet(`/api/chat/detail?chat_id=${chatId.value}`).then(res => {
    title.value = res.data.title
    modelId.value = res.data.model_id
    roleId.value = res.data.role_id
  }).catch(() => {
  })
} else {
  title.value = "新建对话"
}

// 加载模型
httpGet('/api/model/list').then(res => {
  models.value = res.data
  if (!modelId.value) {
    modelId.value = models.value[0].id
  }
  for (let i = 0; i < models.value.length; i++) {
    models.value[i].text = models.value[i].name
    models.value[i].mValue = models.value[i].value
    models.value[i].value = models.value[i].id
  }
  modelValue.value = getModelValue(modelId.value)
  // 加载角色列表
  httpGet(`/api/role/list`).then((res) => {
    roles.value = res.data;
    if (!roleId.value) {
      roleId.value = roles.value[0]['id']
    }
    // build data for role picker
    for (let i = 0; i < roles.value.length; i++) {
      roles.value[i].text = roles.value[i].name
      roles.value[i].value = roles.value[i].id
      roles.value[i].helloMsg = roles.value[i].hello_msg
    }

    role.value = getRoleById(roleId.value)
    columns.value = [roles.value, models.value]
    // 新建对话
    if (!chatId.value) {
      connect(chatId.value, roleId.value, modelId.value)
    }
  }).catch((e) => {
    showNotify({type: "danger", message: '获取聊天角色失败: ' + e.messages})
  })
}).catch(e => {
  showNotify({type: "danger", message: "加载模型失败: " + e.message})
})

const url = ref(location.protocol + '//' + location.host + '/mobile/chat/export?chat_id=' + chatId.value)

onMounted(() => {
  winHeight.value = window.innerHeight - navBarRef.value.$el.offsetHeight - bottomBarRef.value.$el.offsetHeight - 70

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
  if (socket.value !== null) {
    socket.value.close()
    socket.value = null
  }
})

const newChat = (item) => {
  showPicker.value = false
  const options = item.selectedOptions
  roleId.value = options[0].value
  modelId.value = options[1].value
  chatId.value = ""
  chatData.value = []
  role.value = getRoleById(roleId.value)
  title.value = "新建对话"
  connect(chatId.value, roleId.value, modelId.value)
}

const chatData = ref([])
const loading = ref(false)
const finished = ref(false)
const error = ref(false)

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
  if (chatId.value) {
    checkSession().then(() => {
      httpGet('/api/chat/history?chat_id=' + chatId.value).then(res => {
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
        connect(chatId.value, roleId.value, modelId.value);
      }).catch(() => {
        error.value = true
      })
    }).catch(() => {
    })
  }
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
const connect = function (chat_id, role_id, model_id) {
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

  const _socket = new WebSocket(host + `/api/chat/new?session_id=${_sessionId}&role_id=${role_id}&chat_id=${chat_id}&model_id=${model_id}&token=${getUserToken()}`);
  _socket.addEventListener('open', () => {
    loading.value = false
    previousText.value = '';
    canSend.value = true;
    activelyClose.value = false;

    if (isNewChat) { // 加载打招呼信息
      chatData.value.push({
        type: "reply",
        id: randString(32),
        icon: role.value.icon,
        content: role.value.hello_msg,
        orgContent: role.value.hello_msg,
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
            icon: role.value.icon,
            content: ""
          });
          if (isNewChat) {
            title.value = previousText.value
          }
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
      connect(chat_id, role_id, model_id)
    }).catch(() => {
      showLoginDialog(router)
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
      query: {title: title, chat_id: chatId, role: role.value.name, model: modelValue}
    })
  } else if (option.icon === "link") {
    document.getElementById('copy-link-btn').click();
  }
}

const getRoleById = function (rid) {
  for (let i = 0; i < roles.value.length; i++) {
    if (roles.value[i]['id'] === rid) {
      return roles.value[i];
    }
  }
  return null;
}

const getModelValue = (model_id) => {
  for (let i = 0; i < models.value.length; i++) {
    if (models.value[i].id === model_id) {
      return models.value[i].mValue
    }
  }
  return ""
}

// // eslint-disable-next-line no-undef
// const recognition = new webkitSpeechRecognition() || SpeechRecognition();
// //recognition.lang = 'zh-CN' // 设置语音识别语言
// recognition.onresult = function (event) {
//   prompt.value = event.results[0][0].transcript
// };
//
// recognition.onerror = function (event) {
//   showMic.value = false
//   recognition.stop()
//   showNotify({type: 'danger', message: '语音识别错误:' + event.error})
// };
//
// recognition.onend = function () {
//   console.log('语音识别结束');
// };
// const inputVoice = () => {
//   showMic.value = true
//   recognition.start();
// }
//
// const stopVoice = () => {
//   showMic.value = false
//   recognition.stop()
// }
</script>

<style lang="stylus">
@import "@/assets/css/mobile/chat-session.styl"
</style>