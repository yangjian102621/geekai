<template>
  <div class="mobile-mj container">
    <van-nav-bar :title="title"/>

    <div class="content">
      <van-form @submit="generate">
        <div class="text-line">图片比例</div>
        <div class="text-line">
          <van-row :gutter="10">
            <van-col :span="4" v-for="item in rates" :key="item.value">
              <div
                  :class="item.value === params.rate ? 'rate active' : 'rate'"
                  @click="changeRate(item)">
                <div class="icon">
                  <van-image :src="item.img" fit="cover"></van-image>
                </div>
                <div class="text">{{ item.text }}</div>
              </div>
            </van-col>
          </van-row>
        </div>
        <div class="text-line">模型选择</div>
        <div class="text-line">
          <van-row :gutter="10">
            <van-col :span="8" v-for="item in models" :key="item.value">
              <div :class="item.value === params.model ? 'model active' : 'model'"
                   @click="changeModel(item)">
                <div class="icon">
                  <van-image :src="item.img" fit="cover"></van-image>
                </div>
                <div class="text">
                  <van-text-ellipsis :content="item.text"/>
                </div>
              </div>
            </van-col>
          </van-row>
        </div>
        <div class="text-line">
          <van-field label="创意度">
            <template #input>
              <van-slider v-model.number="params.chaos" :max="100" :step="1"
                          @update:model-value="showToast('当前值：' + params.chaos)"/>
            </template>
          </van-field>
        </div>

        <div class="text-line">
          <van-field label="风格化">
            <template #input>
              <van-slider v-model.number="params.stylize" :max="1000" :step="1"
                          @update:model-value="showToast('当前值：' + params.stylize)"/>
            </template>
          </van-field>
        </div>

        <div class="text-line">
          <van-field label="原始模式">
            <template #input>
              <van-switch v-model="params.raw"/>
            </template>
          </van-field>
        </div>

        <div class="text-line">
          <van-field
              v-model="params.prompt"
              rows="3"
              autosize
              label="提示词"
              type="textarea"
              placeholder="如：一个美丽的中国女孩站在电影院门口，手上拿着爆米花，微笑，写实风格，电影灯光效果，半身像"
          />
        </div>

        <van-collapse v-model="activeColspan">
          <van-collapse-item title="垫图" name="img">
            <van-field>
              <template #input>
                <van-uploader v-model="imgList" :after-read="uploadImg"/>
              </template>
            </van-field>
          </van-collapse-item>
          <van-collapse-item title="反向提示词" name="neg_prompt">
            <van-field
                v-model="params.prompt"
                rows="3"
                autosize
                type="textarea"
                placeholder="不想出现在图片上的元素(例如：树，建筑)"
            />
          </van-collapse-item>
        </van-collapse>

        <div class="text-line">
          <van-button round block type="primary" native-type="submit">
            <van-tag type="success">可用额度:{{ imgCalls }}</van-tag>
            立即生成
          </van-button>
        </div>
      </van-form>

      <h3>任务列表</h3>
      <div class="running-job-list">
        <van-empty v-if="runningJobs.length ===0"
                   image="https://fastly.jsdelivr.net/npm/@vant/assets/custom-empty-image.png"
                   image-size="80"
                   description="暂无记录"
        />
        <van-grid :gutter="10" :column-num="3" v-else>
          <van-grid-item v-for="item in runningJobs">
            <div v-if="item.progress > 0">
              <van-image :src="item['img_url']">
                <template v-slot:error>加载失败</template>
              </van-image>
              <div class="progress">
                <van-circle
                    v-model:current-rate="item.progress"
                    :rate="item.progress"
                    :speed="100"
                    :text="item.progress+'%'"
                    :stroke-width="60"
                    size="90px"
                />
              </div>
            </div>

            <div v-else class="task-in-queue">
              <span class="icon"><i class="iconfont icon-quick-start"></i></span>
              <span class="text">排队中</span>
            </div>

          </van-grid-item>
        </van-grid>
      </div>

      <h3>创作记录</h3>
      <div class="finish-job-list">
        <van-empty v-if="finishedJobs.length ===0"
                   image="https://fastly.jsdelivr.net/npm/@vant/assets/custom-empty-image.png"
                   image-size="80"
                   description="暂无记录"
        />
        <van-grid :gutter="10" :column-num="2" v-else>
          <van-grid-item v-for="item in finishedJobs">
            <div class="job-item">
              <el-image
                  :src="item['thumb_url']"
                  :class="item['can_opt'] ? '' : 'upscale'" :zoom-rate="1.2"
                  :preview-src-list="[item['img_url']]" fit="cover" :initial-index="0"
                  loading="lazy" v-if="item.progress > 0">
                <template #placeholder>
                  <div class="image-slot">
                    正在加载图片
                  </div>
                </template>

                <template #error>
                  <div class="image-slot" v-if="item['img_url'] === ''">
                    <i class="iconfont icon-loading"></i>
                    <span>正在下载图片</span>
                  </div>
                  <div class="image-slot" v-else>
                    <el-icon>
                      <Picture/>
                    </el-icon>
                  </div>
                </template>
              </el-image>

              <div class="opt" v-if="item['can_opt']">

                <van-grid :gutter="0" :column-num="4">
                  <van-grid-item><a @click="upscale(1, item)" class="opt-btn">U1</a></van-grid-item>
                  <van-grid-item><a @click="upscale(2, item)" class="opt-btn">U2</a></van-grid-item>
                  <van-grid-item><a @click="upscale(3, item)" class="opt-btn">U3</a></van-grid-item>
                  <van-grid-item><a @click="upscale(4, item)" class="opt-btn">U4</a></van-grid-item>
                  <van-grid-item><a @click="variation(1, item)" class="opt-btn">V1</a></van-grid-item>
                  <van-grid-item><a @click="variation(2, item)" class="opt-btn">V2</a></van-grid-item>
                  <van-grid-item><a @click="variation(3, item)" class="opt-btn">V3</a></van-grid-item>
                  <van-grid-item><a @click="variation(4, item)" class="opt-btn">V4</a></van-grid-item>
                </van-grid>
              </div>

              <div class="remove">
                <el-button type="danger" :icon="Delete" @click="removeImage(item)" circle/>
                <el-button type="warning" v-if="item.publish" @click="publishImage(item, false)"
                           circle>
                  <i class="iconfont icon-cancel-share"></i>
                </el-button>
                <el-button type="success" v-else @click="publishImage(item, true)" circle>
                  <i class="iconfont icon-share-bold"></i>
                </el-button>
                <el-button type="primary" @click="showPrompt(item)" circle>
                  <i class="iconfont icon-prompt"></i>
                </el-button>
              </div>
            </div>
          </van-grid-item>
        </van-grid>
      </div>

    </div>
  </div>
</template>

<script setup>
import {onMounted, ref} from "vue";
import {showConfirmDialog, showFailToast, showNotify, showToast, showDialog} from "vant";
import {httpGet, httpPost} from "@/utils/http";
import Compressor from "compressorjs";
import {ElMessage} from "element-plus";
import {getSessionId} from "@/store/session";
import {checkSession} from "@/action/session";
import Clipboard from "clipboard";
import {useRouter} from "vue-router";
import {Delete, Picture} from "@element-plus/icons-vue";

const title = ref('MidJourney 绘画')
const activeColspan = ref([""])

const rates = [
  {css: "square", value: "1:1", text: "1:1", img: "/images/mj/rate_1_1.png"},
  {css: "size2-3", value: "2:3", text: "2:3", img: "/images/mj/rate_3_4.png"},
  {css: "size3-4", value: "3:4", text: "3:4", img: "/images/mj/rate_3_4.png"},
  {css: "size4-3", value: "4:3", text: "4:3", img: "/images/mj/rate_4_3.png"},
  {css: "size16-9", value: "16:9", text: "16:9", img: "/images/mj/rate_16_9.png"},
  {css: "size9-16", value: "9:16", text: "9:16", img: "/images/mj/rate_9_16.png"},
]
const models = [
  {text: "MJ-6.0", value: " --v 6", img: "/images/mj/mj-v6.png"},
  {text: "MJ-5.2", value: " --v 5.2", img: "/images/mj/mj-v5.2.png"},
  {text: "Niji5 原始", value: " --niji 5", img: "/images/mj/mj-niji.png"},
  {text: "Niji5 可爱", value: " --niji 5 --style cute", img: "/images/mj/nj1.jpg"},
  {text: "Niji5 风景", value: " --niji 5 --style scenic", img: "/images/mj/nj2.jpg"},
  {text: "Niji5 表现力", value: " --niji 5 --style expressive", img: "/images/mj/nj3.jpg"},
]
const imgList = ref([])
const params = ref({
  task_type: "image",
  rate: rates[0].value,
  model: models[0].value,
  chaos: 0,
  stylize: 0,
  seed: 0,
  img_arr: [],
  raw: false,
  weight: 0.25,
  prompt: "",
  neg_prompt: "",
  tile: false,
  quality: 0
})
const imgCalls = ref(0)
const userId = ref(0)
const router = useRouter()
const runningJobs = ref([])
const finishedJobs = ref([])
const socket = ref(null)

onMounted(() => {
  checkSession().then(user => {
    imgCalls.value = user['img_calls']
    userId.value = user.id

    fetchRunningJobs(userId.value)
    fetchFinishJobs(userId.value)
    connect()

  }).catch(() => {
    router.push('/login')
  });
})

const heartbeatHandle = ref(null)
const connect = () => {
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

  const _socket = new WebSocket(host + `/api/mj/client?user_id=${userId.value}`);
  _socket.addEventListener('open', () => {
    socket.value = _socket;

    // 发送心跳消息
    sendHeartbeat()
  });

  _socket.addEventListener('message', event => {
    if (event.data instanceof Blob) {
      fetchRunningJobs(userId.value)
      fetchFinishJobs(userId.value)
    }
  });

  _socket.addEventListener('close', () => {
    connect()
  });
}

// 获取运行中的任务
const fetchRunningJobs = (userId) => {
  httpGet(`/api/mj/jobs?status=0&user_id=${userId}`).then(res => {
    const jobs = res.data
    const _jobs = []
    for (let i = 0; i < jobs.length; i++) {
      if (jobs[i].progress === -1) {
        showNotify({
          message: `任务执行失败：${jobs[i]['err_msg']}`,
          type: 'danger',
        })
        imgCalls.value += 1
        continue
      }
      _jobs.push(jobs[i])
    }
    runningJobs.value = _jobs
  }).catch(e => {
    ElMessage.error("获取任务失败：" + e.message)
  })
}

const fetchFinishJobs = (userId) => {
  // 获取已完成的任务
  httpGet(`/api/mj/jobs?status=1&user_id=${userId}`).then(res => {
    const jobs = res.data
    for (let i = 0; i < jobs.length; i++) {
      if (jobs[i]['use_proxy']) {
        jobs[i]['thumb_url'] = jobs[i]['img_url'] + '?x-oss-process=image/quality,q_60&format=webp'
      } else {
        if (jobs[i].type === 'upscale' || jobs[i].type === 'swapFace') {
          jobs[i]['thumb_url'] = jobs[i]['img_url'] + '?imageView2/1/w/480/h/600/q/75'
        } else {
          jobs[i]['thumb_url'] = jobs[i]['img_url'] + '?imageView2/1/w/480/h/480/q/75'
        }
      }

      if (jobs[i].type === 'image' || jobs[i].type === 'variation') {
        jobs[i]['can_opt'] = true
      }
    }
    finishedJobs.value = jobs
  }).catch(e => {
    ElMessage.error("获取任务失败：" + e.message)
  })
}

// 切换图片比例
const changeRate = (item) => {
  params.value.rate = item.value
}
// 切换模型
const changeModel = (item) => {
  params.value.model = item.value
}


// 图片上传
const uploadImg = (file) => {
  file.status = "uploading"
  // 压缩图片并上传
  new Compressor(file.file, {
    quality: 0.6,
    success(result) {
      const formData = new FormData();
      formData.append('file', result, result.name);
      // 执行上传操作
      httpPost('/api/upload', formData).then(res => {
        file.url = res.data.url
        file.status = "done"
      }).catch(e => {
        file.status = 'failed'
        file.message = '上传失败'
        showFailToast("图片上传失败：" + e.message)
      })
    },
    error(err) {
      console.log(err.message);
    },
  });
};

const send = (url, index, item) => {
  httpPost(url, {
    index: index,
    channel_id: item.channel_id,
    message_id: item.message_id,
    message_hash: item.hash,
    session_id: getSessionId(),
    prompt: item.prompt,
  }).then(() => {
    ElMessage.success("任务推送成功，请耐心等待任务执行...")
    imgCalls.value -= 1
  }).catch(e => {
    ElMessage.error("任务推送失败：" + e.message)
  })
}

// 图片放大任务
const upscale = (index, item) => {
  send('/api/mj/upscale', index, item)
}

// 图片变换任务
const variation = (index, item) => {
  send('/api/mj/variation', index, item)
}

const generate = () => {
  if (params.value.prompt === '' && params.value.task_type === "image") {
    return showFailToast("请输入绘画提示词！")
  }
  if (params.value.model.indexOf("niji") !== -1 && params.value.raw) {
    return showFailToast("动漫模型不允许启用原始模式")
  }
  params.value.session_id = getSessionId()
  params.value.img_arr = imgList.value.map(img => img.url)
  httpPost("/api/mj/image", params.value).then(() => {
    showToast("绘画任务推送成功，请耐心等待任务执行")
    imgCalls.value -= 1
  }).catch(e => {
    showFailToast("任务推送失败：" + e.message)
  })
}

const removeImage = (item) => {
  showConfirmDialog({
    title: '标题',
    message:
        '此操作将会删除任务和图片，继续操作码?',
  }).then(() => {
    httpPost("/api/mj/remove", {id: item.id, img_url: item.img_url, user_id: userId.value}).then(() => {
      ElMessage.success("任务删除成功")
    }).catch(e => {
      ElMessage.error("任务删除失败：" + e.message)
    })
  }).catch(() => {
    showToast("您取消了操作")
  });
}
// 发布图片到作品墙
const publishImage = (item, action) => {
  let text = "图片发布"
  if (action === false) {
    text = "取消发布"
  }
  httpPost("/api/mj/publish", {id: item.id, action: action}).then(() => {
    ElMessage.success(text + "成功")
    item.publish = action
  }).catch(e => {
    ElMessage.error(text + "失败：" + e.message)
  })
}

const showPrompt = (item) => {
  showDialog({
    title: "绘画提示词",
    message: item.prompt,
  }).then(() => {
    // on close
  });
}

</script>

<style lang="stylus">
@import "@/assets/css/mobile/image-mj.styl"
</style>