<template>
  <div class="mobile-mj">
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
        <van-tabs v-model:active="activeName" @change="tabChange" animated>
          <van-tab title="文生图" name="txt2img">
            <div class="text-line">
              <van-field v-model="params.prompt"
                         rows="3"
                         autosize
                         type="textarea"
                         placeholder="请在此输入绘画提示词，系统会自动翻译中文提示词，高手请直接输入英文提示词"/>
            </div>
          </van-tab>
          <van-tab title="图生图" name="img2img">
            <div class="text-line">
              <van-field v-model="params.prompt"
                         rows="3"
                         autosize
                         type="textarea"
                         placeholder="请在此输入绘画提示词，系统会自动翻译中文提示词，高手请直接输入英文提示词"/>
            </div>

            <div class="text-line">
              <van-uploader v-model="imgList" :after-read="uploadImg"/>
            </div>
            <div class="text-line">
              <van-field label="垫图权重">
                <template #input>
                  <van-slider v-model.number="params.iw" :max="1" :step="0.01"
                              @update:model-value="showToast('当前值：' + params.iw)"/>
                </template>
              </van-field>
            </div>

            <div class="tip-text">提示：只有于 niji6 和 v6 模型支持一致性功能，如果选择其他模型此功能将会生成失败。</div>
            <van-cell-group>
              <van-field
                  v-model="params.cref"
                  center
                  clearable
                  label="角色一致性"
                  placeholder="请输入图片URL或者上传图片"
              >
                <template #button>
                  <van-uploader @click="beforeUpload('cref')" :after-read="uploadImg">
                    <van-button size="mini" type="primary" icon="plus"/>
                  </van-uploader>
                </template>
              </van-field>
            </van-cell-group>

            <van-cell-group>
              <van-field
                  v-model="params.sref"
                  center
                  clearable
                  label="风格一致性"
                  placeholder="请输入图片URL或者上传图片"
              >
                <template #button>
                  <van-uploader @click="beforeUpload('sref')" :after-read="uploadImg">
                    <van-button size="mini" type="primary" icon="plus"/>
                  </van-uploader>
                </template>
              </van-field>
            </van-cell-group>

            <div class="text-line">
              <van-field label="一致性权重">
                <template #input>
                  <van-slider v-model.number="params.cw" :max="100" :step="1"
                              @update:model-value="showToast('当前值：' + params.cw)"/>
                </template>
              </van-field>
            </div>
          </van-tab>
          <van-tab title="融图" name="blend">
            <div class="tip-text">请上传两张以上的图片，最多不超过五张，超过五张图片请使用图生图功能。</div>
            <div class="text-line">
              <van-uploader v-model="imgList" :after-read="uploadImg"/>
            </div>
          </van-tab>
          <van-tab title="换脸" name="swapFace">
            <div class="tip-text">请上传两张有脸部的图片，用左边图片的脸替换右边图片的脸。</div>
            <div class="text-line">
              <van-uploader v-model="imgList" :after-read="uploadImg"/>
            </div>
          </van-tab>
        </van-tabs>
      </div>

      <div class="text-line">
        <van-collapse v-model="activeColspan">
          <van-collapse-item title="反向提示词" name="neg_prompt">
            <van-field
                v-model="params.neg_prompt"
                rows="3"
                autosize
                type="textarea"
                placeholder="不想出现在图片上的元素(例如：树，建筑)"
            />
          </van-collapse-item>
        </van-collapse>
      </div>

      <div class="text-line">
        <el-tag>绘图消耗{{ mjPower }}算力，U/V 操作消耗{{ mjActionPower }}算力，当前算力：{{ power }}</el-tag>
      </div>

      <div class="text-line">
        <van-button round block type="primary" native-type="submit">
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

      <van-list v-else
                v-model:error="error"
                v-model:loading="loading"
                :finished="finished"
                error-text="请求失败，点击重新加载"
                finished-text="没有更多了"
                @load="onLoad"
      >
        <van-grid :gutter="10" :column-num="2">
          <van-grid-item v-for="item in finishedJobs">
            <div class="job-item">
              <van-image
                  :src="item['thumb_url']"
                  :class="item['can_opt'] ? '' : 'upscale'"
                  lazy-load
                  @click="imageView(item)"
                  fit="cover">
                <template v-slot:loading>
                  <van-loading type="spinner" size="20"/>
                </template>
              </van-image>

              <div class="opt" v-if="item['can_opt']">

                <van-grid :gutter="3" :column-num="4">
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
      </van-list>

    </div>

  </div>
</template>

<script setup>
import {nextTick, onMounted, onUnmounted, ref} from "vue";
import {
  showConfirmDialog,
  showFailToast,
  showNotify,
  showToast,
  showDialog,
  showImagePreview,
  showSuccessToast
} from "vant";
import {httpGet, httpPost} from "@/utils/http";
import Compressor from "compressorjs";
import {getSessionId} from "@/store/session";
import {checkSession} from "@/action/session";
import {useRouter} from "vue-router";
import {Delete} from "@element-plus/icons-vue";

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
  {text: "Niji5", value: " --niji 5", img: "/images/mj/mj-niji.png"},
  {text: "Niji5 可爱", value: " --niji 5 --style cute", img: "/images/mj/nj1.jpg"},
  {text: "Niji5 风景", value: " --niji 5 --style scenic", img: "/images/mj/nj2.jpg"},
  {text: "Niji6", value: " --niji 6", img: "/images/mj/nj3.jpg"},
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
  iw: 0,
  prompt: "",
  neg_prompt: "",
  tile: false,
  quality: 0,
  cref: "",
  sref: "",
  cw: 0,
})
const userId = ref(0)
const router = useRouter()
const runningJobs = ref([])
const finishedJobs = ref([])
const socket = ref(null)
const power = ref(0)
const activeName = ref("txt2img")

onMounted(() => {
  checkSession().then(user => {
    power.value = user['power']
    userId.value = user.id

    fetchRunningJobs()
    fetchFinishJobs(1)
    connect()

  }).catch(() => {
    router.push('/login')
  });
})

onUnmounted(() => {
  socket.value = null
})

const mjPower = ref(1)
const mjActionPower = ref(1)
httpGet("/api/config/get?key=system").then(res => {
  mjPower.value = res.data["mj_power"]
  mjActionPower.value = res.data["mj_action_power"]
}).catch(e => {
  showNotify({type: "danger", message: "获取系统配置失败：" + e.message})
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
      fetchRunningJobs()
      fetchFinishJobs(1)
    }
  });

  _socket.addEventListener('close', () => {
    if (socket.value !== null) {
      connect()
    }
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
        if (jobs[i].type === 'image') {
          power.value += mjPower.value
        } else {
          power.value += mjActionPower.value
        }
        continue
      }
      _jobs.push(jobs[i])
    }
    runningJobs.value = _jobs
  }).catch(e => {
    showNotify({type: "danger", message: "获取任务失败：" + e.message})
  })
}

const loading = ref(false)
const finished = ref(false)
const error = ref(false)
const page = ref(0)
const pageSize = ref(10)
const fetchFinishJobs = (page) => {
  loading.value = true
  // 获取已完成的任务
  httpGet(`/api/mj/jobs?status=1&page=${page}&page_size=${pageSize.value}`).then(res => {
    const jobs = res.data
    for (let i = 0; i < jobs.length; i++) {
      if (jobs[i].progress === -1) {
        showNotify({
          message: `任务ID：${jobs[i]['task_id']} 原因：${jobs[i]['err_msg']}`,
          type: 'danger',
        })
        if (jobs[i].type === 'image') {
          power.value += mjPower.value
        } else {
          power.value += mjActionPower.value
        }
        continue
      }

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
    if (jobs.length < pageSize.value) {
      finished.value = true
    }
    if (page === 1) {
      finishedJobs.value = jobs
    } else {
      finishedJobs.value = finishedJobs.value.concat(jobs)
    }
    nextTick(() => loading.value = false)
  }).catch(e => {
    loading.value = false
    error.value = true
    showFailToast("获取任务失败：" + e.message)
  })
}

const onLoad = () => {
  page.value += 1
  fetchFinishJobs(page.value)
};

// 切换图片比例
const changeRate = (item) => {
  params.value.rate = item.value
}
// 切换模型
const changeModel = (item) => {
  params.value.model = item.value
}

const imgKey = ref("")
const beforeUpload = (key) => {
  imgKey.value = key
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
        if (imgKey.value !== "") { // 单张图片上传
          params.value[imgKey.value] = res.data.url
          imgKey.value = ''
        }
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
    showSuccessToast("任务推送成功，请耐心等待任务执行...")
    power.value -= mjActionPower.value
  }).catch(e => {
    showFailToast("任务推送失败：" + e.message)
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
    power.value -= mjPower.value
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
      showSuccessToast("任务删除成功")
    }).catch(e => {
      showFailToast("任务删除失败：" + e.message)
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
    showSuccessToast(text + "成功")
    item.publish = action
  }).catch(e => {
    showFailToast(text + "失败：" + e.message)
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

const imageView = (item) => {
  showImagePreview([item['img_url']]);
}

// 切换菜单
const tabChange = (tab) => {
  if (tab === "txt2img" || tab === "img2img") {
    params.value.task_type = "image"
  } else {
    params.value.task_type = tab
  }
}
</script>

<style lang="stylus">
@import "@/assets/css/mobile/image-mj.styl"
</style>