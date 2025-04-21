<template>
  <div>
    <div class="page-mark-map">
      <div class="inner custom-scroll">
        <div class="mark-map-box" :style="{ height: leftBoxHeight + 'px' }">
          <h2>思维导图创作中心</h2>

          <div class="mark-map-params">
            <el-form label-width="80px" label-position="left">
              <div class="param-line">你的需求？</div>
              <div class="param-line">
                <el-input
                  v-model="prompt"
                  :autosize="{ minRows: 4, maxRows: 6 }"
                  type="textarea"
                  placeholder="请给AI输入提示词，让AI帮你完善"
                />
              </div>

              <div class="param-line">请选择生成思维导图的AI模型</div>
              <div class="param-line">
                <el-select v-model="modelID" placeholder="请选择模型" style="width: 100%">
                  <el-option
                    v-for="item in models"
                    :key="item.id"
                    :label="item.name"
                    :value="item.id"
                  >
                    <span>{{ item.name }}</span>
                    <el-tag
                      style="margin-left: 5px; position: relative; top: -2px"
                      type="info"
                      size="small"
                      >{{ item.power }}算力
                    </el-tag>
                  </el-option>
                </el-select>
              </div>

              <div class="text-info">
                <el-text type="primary"
                  >当前可用算力：<el-text type="warning">{{ loginUser.power }}</el-text></el-text
                >
              </div>

              <div class="submit-btn flex-center">
                <el-button
                  style="width: 200px"
                  type="primary"
                  :dark="false"
                  round
                  @click="generateAI"
                  :loading="loading"
                >
                  生成思维导图
                </el-button>
              </div>

              <div class="param-line">使用已有内容生成？</div>
              <div class="param-line">
                <el-input
                  v-model="content"
                  :autosize="{ minRows: 4, maxRows: 6 }"
                  type="textarea"
                  placeholder="请用markdown语法输入您想要生成思维导图的内容！"
                />
              </div>

              <div class="param-line flex-center">
                <el-button
                  color="rgb(78, 51, 254)"
                  style="width: 200px"
                  :dark="false"
                  round
                  @click="generate"
                  >直接生成（免费）</el-button
                >
              </div>
            </el-form>
          </div>
        </div>

        <div class="chat-box">
          <!-- <div class="top-bar">
            <el-button @click="downloadImage" type="primary">
              <el-icon>
                <Download />
              </el-icon>
              <span>下载图片-</span>
            </el-button>
          </div> -->

          <div class="body" id="markmap">
            <svg ref="svgRef" :style="{ height: rightBoxHeight + 'px' }" />
            <div id="toolbar">
              <el-button @click="downloadImage" type="primary">
                <el-icon>
                  <Download />
                </el-icon>
                <span>下载图片</span>
              </el-button>
            </div>
          </div>
        </div>
        <!-- end task list box -->
      </div>
    </div>
  </div>
</template>

<script setup>
import { nextTick, onMounted, ref } from 'vue'
import { Markmap } from 'markmap-view'
import { Transformer } from 'markmap-lib'
import { checkSession, getSystemInfo } from '@/store/cache'
import { httpGet, httpPost } from '@/utils/http'
import { ElMessage } from 'element-plus'
import { Download } from '@element-plus/icons-vue'
import { Toolbar } from 'markmap-toolbar'
import { useSharedStore } from '@/store/sharedata'

const leftBoxHeight = ref(window.innerHeight - 105)
//const rightBoxHeight = ref(window.innerHeight - 115);
const rightBoxHeight = ref(window.innerHeight)

const prompt = ref('')
const text = ref('')
const content = ref(text.value)
const html = ref('')

const isLogin = ref(false)
const loginUser = ref({ power: 0 })
const transformer = new Transformer()
const store = useSharedStore()
const loading = ref(false)

const svgRef = ref(null)
const markMap = ref(null)
const models = ref([])
const modelID = ref(0)
const cacheKey = ref('MarkMapCache')

onMounted(async () => {
  const cache = localStorage.getItem(cacheKey.value)
  if (cache) {
    text.value = cache
  } else {
    const res = await getSystemInfo().catch((e) => {
      ElMessage.error('获取系统配置失败：' + e.message)
    })
    text.value = res.data['mark_map_text']
    content.value = text.value
  }

  initData()
  nextTick(() => {
    try {
      markMap.value = Markmap.create(svgRef.value)
      const { el } = Toolbar.create(markMap.value)
      document.getElementById('toolbar').append(el)
      update()
    } catch (e) {
      console.error(e)
    }
  })
})

const initData = () => {
  httpGet('/api/model/list')
    .then((res) => {
      for (let v of res.data) {
        models.value.push(v)
      }
      modelID.value = models.value[0].id
    })
    .catch((e) => {
      ElMessage.error('获取模型失败：' + e.message)
    })

  checkSession()
    .then((user) => {
      loginUser.value = user
      isLogin.value = true
    })
    .catch(() => {})
}

const update = () => {
  try {
    const { root } = transformer.transform(processContent(text.value))
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
  const lines = text.split('\n')
  for (let line of lines) {
    if (line.indexOf('```') !== -1) {
      continue
    }
    line = line.replace(/([*_~`>])|(\d+\.)\s/g, '')
    arr.push(line)
  }
  return arr.join('\n')
}

window.onresize = () => {
  leftBoxHeight.value = window.innerHeight - 145
  rightBoxHeight.value = window.innerHeight - 85
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
    return ElMessage.error('请输入你的需求')
  }
  if (!isLogin.value) {
    store.setShowLoginDialog(true)
    return
  }
  loading.value = true
  httpPost('/api/markMap/gen', {
    prompt: prompt.value,
    model_id: modelID.value,
  })
    .then((res) => {
      text.value = res.data
      content.value = processContent(text.value)
      const model = getModelById(modelID.value)
      loginUser.value.power -= model.power
      nextTick(() => update())
      loading.value = false
      // 缓存结果
      localStorage.setItem(cacheKey.value, text.value)
    })
    .catch((e) => {
      ElMessage.error('生成思维导图失败：' + e.message)
      loading.value = false
    })
}

const getModelById = (modelId) => {
  for (let m of models.value) {
    if (m.id === modelId) {
      return m
    }
  }
}

// download SVG to png file
const downloadImage = async() => {
  // 先自适应思维导图到可视化区域
  await markMap.value.fit() 

  const svgElement = document.getElementById('markmap')
  const serializer = new XMLSerializer()
  const source =
    '<?xml version="1.0" standalone="no"?>\r\n' + serializer.serializeToString(svgRef.value)
  const image = new Image()
  image.src = 'data:image/svg+xml;charset=utf-8,' + encodeURIComponent(source)

  // 分辨率倍数，越高图片越清晰，但文件越大
  const scale = 4 
  const canvas = document.createElement('canvas')
  canvas.width = svgElement.offsetWidth * scale
  canvas.height = svgElement.offsetHeight * scale

  let context = canvas.getContext('2d')
  context.clearRect(0, 0, canvas.width, canvas.height)
  context.fillStyle = 'white'
  context.fillRect(0, 0, canvas.width, canvas.height)

  image.onload = function () {
    // 关键：缩放 context
    context.setTransform(scale, 0, 0, scale, 0, 0)
    context.drawImage(image, 0, 0)
    const a = document.createElement('a')
    a.download = 'geek-ai-xmind.png'
    a.href = canvas.toDataURL('image/png')
    a.click()
  }
}
</script>

<style lang="stylus">
@import "@/assets/css/mark-map.styl"
@import "@/assets/css/custom-scroll.styl"
</style>
