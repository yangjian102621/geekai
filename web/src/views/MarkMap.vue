<template>
  <div>
    <div class="page-mark-map">
      <div class="inner custom-scroll">
        <div class="mark-map-box" :style="{ height: leftBoxHeight + 'px' }">
          <h2>思维导图创作中心</h2>

          <div class="mark-map-params">
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
                <el-select v-model="modelID" placeholder="请选择模型" @change="changeModel" style="width:100%">
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

        <div class="chat-box">
          <div class="top-bar">
            <el-button @click="downloadImage" type="primary">
              <el-icon>
                <Download/>
              </el-icon>
              <span>下载图片</span>
            </el-button>
          </div>

          <div class="markdown" v-if="loading">
            <div :style="{ height: rightBoxHeight + 'px', overflow:'auto',width:'80%' }" v-html="html"></div>
          </div>
          <div class="body" id="markmap" v-show="!loading">
            <svg ref="svgRef" :style="{ height: rightBoxHeight + 'px' }"/>
            <div id="toolbar"></div>
          </div>
        </div><!-- end task list box -->
      </div>

    </div>
  </div>
</template>

<script setup>
import {nextTick, onUnmounted, ref} from 'vue';
import {Markmap} from 'markmap-view';
import {Transformer} from 'markmap-lib';
import {checkSession, getSystemInfo} from "@/store/cache";
import {httpGet} from "@/utils/http";
import {ElMessage} from "element-plus";
import {Download} from "@element-plus/icons-vue";
import {Toolbar} from 'markmap-toolbar';
import {useSharedStore} from "@/store/sharedata";

const leftBoxHeight = ref(window.innerHeight - 105)
const rightBoxHeight = ref(window.innerHeight - 115)
const title = ref("")

const prompt = ref("")
const text = ref("")
const md = require('markdown-it')({breaks: true});
const content = ref(text.value)
const html = ref("")

const isLogin = ref(false)
const loginUser = ref({power: 0})
const transformer = new Transformer();
const store = useSharedStore();


const svgRef = ref(null)
const markMap = ref(null)
const models = ref([])
const modelID = ref(0)
const loading = ref(false)

getSystemInfo().then(res => {
  text.value = res.data['mark_map_text']
  content.value = text.value
  initData()
  nextTick(() => {
    try {
      markMap.value = Markmap.create(svgRef.value)
      const {el} = Toolbar.create(markMap.value);
      document.getElementById('toolbar').append(el);
      update()
    } catch (e) {
      console.error(e)
    }
  })
}).catch(e => {
  ElMessage.error("获取系统配置失败：" + e.message)
})

const initData = () => {
  httpGet("/api/model/list").then(res => {
    for (let v of res.data) {
      if (v.value.indexOf("gpt-4-gizmo") === -1) {
        models.value.push(v)
      }
    }
    modelID.value = models.value[0].id
  }).catch(e => {
    ElMessage.error("获取模型失败：" + e.message)
  })
  
  checkSession().then(user => {
    loginUser.value = user
    isLogin.value = true
    connect(user.id)
  }).catch(() => {
  });
}

const update = () => {
  try {
    const {root} = transformer.transform(processContent(text.value))
    markMap.value.setData(root)
    markMap.value.fit()
  } catch (e) {
    console.error(e)
  }
}

const processContent = (text) => {
  if (!text) {
    return text
  }

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
const heartbeatHandle = ref(0)
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
      const model = getModelById(modelID.value)
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
            loginUser.value.power -= model.power
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
    checkSession().then(() => {
      connect(userId)
    }).catch(() => {
    })
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
    store.setShowLoginDialog(true)
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

const getModelById = (modelId) => {
  for (let e of models.value) {
    if (e.id === modelId) {
      return e
    }
  }
}

// download SVG to png file
const downloadImage = () => {
  const svgElement = document.getElementById("markmap");
  // 将 SVG 渲染到图片对象
  const serializer = new XMLSerializer()
  const source = '<?xml version="1.0" standalone="no"?>\r\n' + serializer.serializeToString(svgRef.value)
  const image = new Image()
  image.src = 'data:image/svg+xml;charset=utf-8,' + encodeURIComponent(source)

  // 将图片对象渲染
  const canvas = document.createElement('canvas')
  canvas.width = svgElement.offsetWidth
  canvas.height = svgElement.offsetHeight
  let context = canvas.getContext('2d')
  context.clearRect(0, 0, canvas.width, canvas.height);

  image.onload = function () {
    context.drawImage(image, 0, 0)
    const a = document.createElement('a')
    a.download = "geek-ai-xmind.png"
    a.href = canvas.toDataURL(`image/png`)
    a.click()
  }
}

</script>

<style lang="stylus">
@import "@/assets/css/mark-map.styl"
@import "@/assets/css/custom-scroll.styl"
</style>
