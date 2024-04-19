<template>
  <div>
    <div class="page-mark-map">
      <div class="inner custom-scroll">
        <div class="mark-map-box">
          <h2>思维导图创作中心</h2>

          <div class="mark-map-params" :style="{ height: leftBoxHeight + 'px' }">
            <el-form label-width="80px" label-position="left">
              <div class="param-line">
                你的需求？
              </div>
              <div class="param-line">
                <el-input
                    v-model="prompt"
                    :autosize="{ minRows: 4, maxRows: 6 }"
                    type="textarea"
                    placeholder="请给AI输入提示词，让AI帮你完善"
                />
              </div>

              <div class="param-line">
                请选择生成思维导图的AI模型
              </div>
              <div class="param-line">
                <el-select v-model="modelID" placeholder="请选择模型" @change="changeModel">
                  <el-option
                      v-for="item in models"
                      :key="item.id"
                      :label="item.name"
                      :value="item.id"
                  >
                    <span>{{ item.name }}</span>
                    <el-tag style="margin-left: 5px; position: relative; top:-2px" type="info" size="small">{{
                        item.power
                      }}算力
                    </el-tag>
                  </el-option>
                </el-select>
              </div>

              <div class="text-info">
                <el-tag type="success">当前可用算力：{{ loginUser.power }}</el-tag>
              </div>

              <div class="param-line">
                <el-button color="#47fff1" :dark="false" round @click="generateAI" :loading="loading">
                  智能生成思维导图
                </el-button>
              </div>

              <div class="param-line">
                使用已有内容生成？
              </div>
              <div class="param-line">
                <el-input
                    v-model="content"
                    :autosize="{ minRows: 4, maxRows: 6 }"
                    type="textarea"
                    placeholder="请用markdown语法输入您想要生成思维导图的内容！"
                />
              </div>

              <div class="param-line">
                <el-button color="#C5F9AE" :dark="false" round @click="generate">直接生成（免费）</el-button>
              </div>

            </el-form>
          </div>
        </div>

        <div class="right-box">
          <h2>思维导图</h2>
          <div class="markdown" v-if="loading">
            <div v-html="html"></div>
          </div>
          <div class="body" v-show="!loading">
            <svg ref="svgRef" :style="{ height: rightBoxHeight + 'px' }"/>
          </div>
        </div><!-- end task list box -->
      </div>

    </div>

    <login-dialog :show="showLoginDialog" @hide="showLoginDialog =  false" @success="initData"/>
  </div>
</template>

<script setup>
import LoginDialog from "@/components/LoginDialog.vue";
import {nextTick, onMounted, onUnmounted, ref} from 'vue';
import {Markmap} from 'markmap-view';
import {loadCSS, loadJS} from 'markmap-common';
import {Transformer} from 'markmap-lib';
import {checkSession} from "@/action/session";
import {httpGet} from "@/utils/http";
import {ElMessage, ElMessageBox} from "element-plus";

const leftBoxHeight = ref(window.innerHeight - 105)
const rightBoxHeight = ref(window.innerHeight - 85)

const prompt = ref("")
const text = ref(`# Geek-AI 助手

- 完整的开源系统，前端应用和后台管理系统皆可开箱即用。
- 基于 Websocket 实现，完美的打字机体验。
- 内置了各种预训练好的角色应用,轻松满足你的各种聊天和应用需求。
- 支持 OPenAI，Azure，文心一言，讯飞星火，清华 ChatGLM等多个大语言模型。
- 支持 MidJourney / Stable Diffusion AI 绘画集成，开箱即用。
- 支持使用个人微信二维码作为充值收费的支付渠道，无需企业支付通道。
- 已集成支付宝支付功能，微信支付，支持多种会员套餐和点卡购买功能。
- 集成插件 API 功能，可结合大语言模型的 function 功能开发各种强大的插件。
`)
const md = require('markdown-it')({breaks: true});
const content = ref(text.value)
const html = ref("")

const showLoginDialog = ref(false)
const isLogin = ref(false)
const loginUser = ref({power: 0})
const transformer = new Transformer();
const {scripts, styles} = transformer.getAssets()
loadCSS(styles);
loadJS(scripts);


const svgRef = ref(null)
const markMap = ref(null)
const models = ref([])
const modelID = ref(0)
const loading = ref(false)

onMounted(() => {
  initData()
  markMap.value = Markmap.create(svgRef.value)
  update()
});

const initData = () => {
  checkSession().then(user => {
    loginUser.value = user
    isLogin.value = true

    httpGet("/api/model/list").then(res => {
      for (let v of res.data) {
        if (v.platform === "OpenAI") {
          models.value.push(v)
        }
      }
      modelID.value = models.value[0].id
      connect(user.id)
    }).catch(e => {
      ElMessage.error("获取模型失败：" + e.message)
    })
  }).catch(() => {
  });
}

const update = () => {

  const {root} = transformer.transform(processContent(text.value))
  markMap.value.setData(root)
  markMap.value.fit()
}

const processContent = (text) => {
  const arr = []
  const lines = text.split("\n")
  for (let line of lines) {
    if (line.indexOf("```") !== -1) {
      continue
    }
    line = line.replace(/([*_~`>])|(\d+\.)\s/g, '')
    arr.push(line)
  }
  return arr.join("\n")
}

onUnmounted(() => {
  if (socket.value !== null) {
    socket.value.close()
  }
  socket.value = null
})

window.onresize = () => {
  leftBoxHeight.value = window.innerHeight - 145
  rightBoxHeight.value = window.innerHeight - 85
}

const socket = ref(null)
const heartbeatHandle = ref(null)
const connect = (userId) => {
  if (socket.value !== null) {
    socket.value.close()
  }

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
    clearTimeout(heartbeatHandle.value)
    new Promise((resolve, reject) => {
      if (socket.value !== null) {
        socket.value.send(JSON.stringify({type: "heartbeat", content: "ping"}))
      }
      resolve("success")
    }).then(() => {
      heartbeatHandle.value = setTimeout(() => sendHeartbeat(), 5000)
    });
  }

  const _socket = new WebSocket(host + `/api/markMap/client?user_id=${userId}&model_id=${modelID.value}`);
  _socket.addEventListener('open', () => {
    socket.value = _socket;

    // 发送心跳消息
    sendHeartbeat()
  });

  _socket.addEventListener('message', event => {
    if (event.data instanceof Blob) {
      const reader = new FileReader();
      reader.readAsText(event.data, "UTF-8")
      reader.onload = () => {
        const data = JSON.parse(String(reader.result))
        switch (data.type) {
          case "start":
            text.value = ""
            break
          case "middle":
            text.value += data.content
            html.value = md.render(processContent(text.value))
            break
          case "end":
            loading.value = false
            content.value = processContent(text.value)
            nextTick(() => update())
            break
          case "error":
            loading.value = false
            ElMessage.error(data.content)
            break
        }
      }
    }
  })

  _socket.addEventListener('close', () => {
    loading.value = false
    if (socket.value !== null) {
      connect(userId)
    }
  });
}

const generate = () => {
  text.value = content.value
  update()
}

// 使用 AI 智能生成
const generateAI = () => {
  html.value = ''
  text.value = ''
  if (prompt.value === '') {
    return ElMessage.error("请输入你的需求")
  }
  if (!isLogin.value) {
    showLoginDialog.value = true
    return
  }
  loading.value = true
  socket.value.send(JSON.stringify({type: "message", content: prompt.value}))
}

const changeModel = () => {
  if (socket.value !== null) {
    socket.value.send(JSON.stringify({type: "model_id", content: modelID.value}))
  }
}

</script>

<style lang="stylus">
@import "@/assets/css/mark-map.styl"
@import "@/assets/css/custom-scroll.styl"
</style>
