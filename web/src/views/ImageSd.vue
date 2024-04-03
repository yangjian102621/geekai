<template>
  <div>
    <div class="page-sd">
      <div class="inner custom-scroll">
        <div class="sd-box">
          <h2>Stable Diffusion 创作中心</h2>

          <div class="sd-params" :style="{ height: mjBoxHeight + 'px' }">
            <el-form :model="params" label-width="80px" label-position="left">
              <div class="param-line" style="padding-top: 10px">
                <el-form-item label="采样方法">
                  <template #default>
                    <div class="form-item-inner">
                      <el-select v-model="params.sampler" size="small">
                        <el-option v-for="item in samplers" :label="item" :value="item" :key="item"/>
                      </el-select>
                      <el-tooltip
                          effect="light"
                          content="出图效果比较好的一般是 Euler 和 DPM 系列算法"
                          raw-content
                          placement="right"
                      >
                        <el-icon>
                          <InfoFilled/>
                        </el-icon>
                      </el-tooltip>
                    </div>
                  </template>
                </el-form-item>
              </div>

              <div class="param-line">
                <el-form-item label="图片尺寸">
                  <template #default>
                    <div class="form-item-inner">
                      <el-row :gutter="20">
                        <el-col :span="12">
                          <el-input v-model.number="params.width" size="small" placeholder="图片宽度"/>
                        </el-col>
                        <el-col :span="12">
                          <el-input v-model.number="params.height" size="small" placeholder="图片高度"/>
                        </el-col>
                      </el-row>
                    </div>
                  </template>
                </el-form-item>
              </div>

              <div class="param-line">
                <el-form-item label="迭代步数">
                  <template #default>
                    <div class="form-item-inner">
                      <el-input v-model.number="params.steps" size="small"/>
                      <el-tooltip
                          effect="light"
                          content="值越大则代表细节越多，同时也意味着出图速度越慢"
                          raw-content
                          placement="right"
                      >
                        <el-icon>
                          <InfoFilled/>
                        </el-icon>
                      </el-tooltip>
                    </div>
                  </template>
                </el-form-item>
              </div>

              <div class="param-line">
                <el-form-item label="引导系数">
                  <template #default>
                    <div class="form-item-inner">
                      <el-input v-model.number="params.cfg_scale" size="small"/>
                      <el-tooltip
                          effect="light"
                          content="提示词引导系数，图像在多大程度上服从提示词<br/> 较低值会产生更有创意的结果"
                          raw-content
                          placement="right"
                      >
                        <el-icon>
                          <InfoFilled/>
                        </el-icon>
                      </el-tooltip>
                    </div>
                  </template>
                </el-form-item>
              </div>

              <div class="param-line">
                <el-form-item label="随机因子">
                  <template #default>
                    <div class="form-item-inner">
                      <el-input v-model.number="params.seed" size="small"/>
                      <el-tooltip
                          effect="light"
                          content="随机数种子，相同的种子会得到相同的结果<br/> 设置为 -1 则每次随机生成种子"
                          raw-content
                          placement="right"
                      >
                        <el-icon>
                          <InfoFilled/>
                        </el-icon>
                      </el-tooltip>

                      <el-tooltip
                          effect="light"
                          content="使用随机数"
                          raw-content
                          placement="right"
                      >
                        <el-icon @click="params.seed = -1">
                          <Orange/>
                        </el-icon>
                      </el-tooltip>
                    </div>
                  </template>
                </el-form-item>
              </div>

              <div class="param-line">
                <el-form-item label="高清修复">
                  <template #default>
                    <div class="form-item-inner">
                      <el-switch v-model="params.hd_fix" style="--el-switch-on-color: #47fff1;"/>
                      <el-tooltip
                          effect="light"
                          content="先以较小的分辨率生成图像，接着方法图像<br />然后在不更改构图的情况下再修改细节"
                          raw-content
                          placement="right"
                      >
                        <el-icon style="margin-top: 6px">
                          <InfoFilled/>
                        </el-icon>
                      </el-tooltip>
                    </div>
                  </template>
                </el-form-item>
              </div>

              <div v-show="params.hd_fix">
                <div class="param-line">
                  <el-form-item label="重绘幅度">
                    <template #default>
                      <div class="form-item-inner">
                        <el-slider v-model.number="params.hd_redraw_rate" :max="1" :step="0.1"
                                   style="width: 180px;--el-slider-main-bg-color:#47fff1"/>
                        <el-tooltip
                            effect="light"
                            content="决定算法对图像内容的影响程度<br />较大的值将得到越有创意的图像"
                            raw-content
                            placement="right"
                        >
                          <el-icon style="margin-top: 6px">
                            <InfoFilled/>
                          </el-icon>
                        </el-tooltip>
                      </div>
                    </template>
                  </el-form-item>
                </div>

                <div class="param-line">
                  <el-form-item label="放大算法">
                    <template #default>
                      <div class="form-item-inner">
                        <el-select v-model="params.hd_scale_alg" size="small">
                          <el-option v-for="item in scaleAlg" :label="item" :value="item" :key="item"/>
                        </el-select>
                        <el-tooltip
                            effect="light"
                            content="高清修复放大算法，主流算法有Latent和ESRGAN_4x"
                            raw-content
                            placement="right"
                        >
                          <el-icon>
                            <InfoFilled/>
                          </el-icon>
                        </el-tooltip>
                      </div>
                    </template>
                  </el-form-item>
                </div>

                <div class="param-line">
                  <el-form-item label="放大倍数">
                    <template #default>
                      <div class="form-item-inner">
                        <el-input v-model.number="params.hd_scale" size="small"/>
                        <el-tooltip
                            effect="light"
                            content="随机数种子，相同的种子会得到相同的结果<br/> 设置为 -1 则每次随机生成种子"
                            raw-content
                            placement="right"
                        >
                          <el-icon>
                            <InfoFilled/>
                          </el-icon>
                        </el-tooltip>
                      </div>
                    </template>
                  </el-form-item>
                </div>

                <div class="param-line">
                  <el-form-item label="迭代步数">
                    <template #default>
                      <div class="form-item-inner">
                        <el-input v-model.number="params.hd_steps" size="small"/>
                        <el-tooltip
                            effect="light"
                            content="重绘迭代步数，如果设置为0，则设置跟原图相同的迭代步数"
                            raw-content
                            placement="right"
                        >
                          <el-icon>
                            <InfoFilled/>
                          </el-icon>
                        </el-tooltip>
                      </div>
                    </template>
                  </el-form-item>
                </div>
              </div>

              <div class="param-line">
                <el-input
                    v-model="params.prompt"
                    :autosize="{ minRows: 4, maxRows: 6 }"
                    type="textarea"
                    ref="promptRef"
                    placeholder="请在此输入绘画提示词，系统会自动翻译中文提示词，高手请直接输入英文提示词"
                />
              </div>

              <div class="param-line pt">
                <span>反向提示词：</span>
                <el-tooltip
                    effect="light"
                    content="不希望出现的元素，下面给了默认的起手式"
                    placement="right"
                >
                  <el-icon>
                    <InfoFilled/>
                  </el-icon>
                </el-tooltip>
              </div>
              <div class="param-line">
                <el-input
                    v-model="params.neg_prompt"
                    :autosize="{ minRows: 4, maxRows: 6 }"
                    type="textarea"
                    placeholder="反向提示词"
                />
              </div>

              <div class="text-info">
                <el-tag>每次绘图消耗{{ sdPower }}算力</el-tag>
                <el-tag type="success">当前可用算力：{{ power }}</el-tag>
              </div>

            </el-form>
          </div>
          <div class="submit-btn">
            <el-button color="#47fff1" :dark="false" round @click="generate">立即生成</el-button>
          </div>
        </div>
        <div class="task-list-box" @scrollend="handleScrollEnd">
          <div class="task-list-inner" :style="{ height: listBoxHeight + 'px' }">
            <div class="job-list-box">
              <h2>任务列表</h2>
              <div class="running-job-list">
                <ItemList :items="runningJobs" v-if="runningJobs.length > 0" :width="240">
                  <template #default="scope">
                    <div class="job-item">
                      <div v-if="scope.item.progress > 0" class="job-item-inner">
                        <el-image :src="scope.item['img_url']"
                                  fit="cover"
                                  loading="lazy">
                          <template #placeholder>
                            <div class="image-slot">
                              正在加载图片
                            </div>
                          </template>

                          <template #error>
                            <div class="image-slot">
                              <el-icon v-if="scope.item['img_url'] !== ''">
                                <Picture/>
                              </el-icon>
                            </div>
                          </template>
                        </el-image>

                        <div class="progress">
                          <el-progress type="circle" :percentage="scope.item.progress" :width="100" color="#47fff1"/>
                        </div>
                      </div>
                      <el-image fit="cover" v-else>
                        <template #error>
                          <div class="image-slot">
                            <i class="iconfont icon-quick-start"></i>
                            <span>任务正在排队中</span>
                          </div>
                        </template>
                      </el-image>
                    </div>
                  </template>
                </ItemList>
                <el-empty :image-size="100" v-else/>
              </div>
              <h2>创作记录</h2>
              <div class="finish-job-list" v-loading="loading" element-loading-background="rgba(0, 0, 0, 0.5)">
                <div v-if="finishedJobs.length > 0">
                  <ItemList :items="finishedJobs" :width="240" :gap="16">
                    <template #default="scope">
                      <div class="job-item animate" @click="showTask(scope.item)">
                        <el-image
                            :src="scope.item['img_url']+'?imageView2/1/w/240/h/240/q/75'"
                            fit="cover"
                            loading="lazy">
                          <template #placeholder>
                            <div class="image-slot">
                              正在加载图片
                            </div>
                          </template>

                          <template #error>
                            <div class="image-slot">
                              <el-icon>
                                <Picture/>
                              </el-icon>
                            </div>
                          </template>
                        </el-image>

                        <div class="remove">
                          <el-button type="danger" :icon="Delete" @click="removeImage($event,scope.item)" circle/>
                          <el-button type="warning" v-if="scope.item.publish"
                                     @click="publishImage($event,scope.item, false)"
                                     circle>
                            <i class="iconfont icon-cancel-share"></i>
                          </el-button>
                          <el-button type="success" v-else @click="publishImage($event,scope.item, true)" circle>
                            <i class="iconfont icon-share-bold"></i>
                          </el-button>
                        </div>
                      </div>
                    </template>
                  </ItemList>

                  <div class="no-more-data" v-if="isOver">
                    <span>没有更多数据了</span>
                    <i class="iconfont icon-face"></i>
                  </div>
                </div>
                <el-empty :image-size="100" v-else/>
              </div> <!-- end finish job list-->
            </div>
          </div>

        </div><!-- end task list box -->
      </div>

      <!-- 任务详情弹框 -->
      <el-dialog v-model="showTaskDialog" title="绘画任务详情" :fullscreen="true">
        <el-row :gutter="20">
          <el-col :span="16">
            <div class="img-container" :style="{maxHeight: fullImgHeight+'px'}">
              <el-image :src="item['img_url']" fit="contain"/>
            </div>
          </el-col>
          <el-col :span="8">
            <div class="task-info">
              <div class="info-line">
                <el-divider>
                  正向提示词
                </el-divider>
                <div class="prompt">
                  <span>{{ item.prompt }}</span>
                  <el-icon class="copy-prompt-sd" :data-clipboard-text="item.prompt">
                    <DocumentCopy/>
                  </el-icon>
                </div>

              </div>

              <div class="info-line">
                <el-divider>
                  反向提示词
                </el-divider>
                <div class="prompt">
                  <span>{{ item.params.neg_prompt }}</span>
                  <el-icon class="copy-prompt-sd" :data-clipboard-text="item.params.neg_prompt">
                    <DocumentCopy/>
                  </el-icon>
                </div>
              </div>

              <div class="info-line">
                <div class="wrapper">
                  <label>采样方法：</label>
                  <div class="item-value">{{ item.params.sampler }}</div>
                </div>
              </div>

              <div class="info-line">
                <div class="wrapper">
                  <label>图片尺寸：</label>
                  <div class="item-value">{{ item.params.width }} x {{ item.params.height }}</div>
                </div>
              </div>

              <div class="info-line">
                <div class="wrapper">
                  <label>迭代步数：</label>
                  <div class="item-value">{{ item.params.steps }}</div>
                </div>
              </div>

              <div class="info-line">
                <div class="wrapper">
                  <label>引导系数：</label>
                  <div class="item-value">{{ item.params.cfg_scale }}</div>
                </div>
              </div>

              <div class="info-line">
                <div class="wrapper">
                  <label>随机因子：</label>
                  <div class="item-value">{{ item.params.seed }}</div>
                </div>
              </div>

              <div v-if="item.params.hd_fix">
                <el-divider>
                  高清修复
                </el-divider>
                <div class="info-line">
                  <div class="wrapper">
                    <label>重绘幅度：</label>
                    <div class="item-value">{{ item.params.hd_redraw_rate }}</div>
                  </div>
                </div>

                <div class="info-line">
                  <div class="wrapper">
                    <label>放大算法：</label>
                    <div class="item-value">{{ item.params.hd_scale_alg }}</div>
                  </div>
                </div>

                <div class="info-line">
                  <div class="wrapper">
                    <label>放大倍数：</label>
                    <div class="item-value">{{ item.params.hd_scale }}</div>
                  </div>
                </div>

                <div class="info-line">
                  <div class="wrapper">
                    <label>迭代步数：</label>
                    <div class="item-value">{{ item.params.hd_steps }}</div>
                  </div>
                </div>
              </div>

              <div class="copy-params">
                <el-button type="primary" round @click="copyParams(item)">画一张同款的</el-button>
              </div>

            </div>
          </el-col>
        </el-row>

      </el-dialog>
    </div>

    <login-dialog :show="showLoginDialog" @hide="showLoginDialog =  false" @success="initData"/>
  </div>
</template>

<script setup>
import {onMounted, onUnmounted, ref} from "vue"
import {Delete, DocumentCopy, InfoFilled, Orange, Picture} from "@element-plus/icons-vue";
import {httpGet, httpPost} from "@/utils/http";
import {ElMessage, ElMessageBox, ElNotification} from "element-plus";
import ItemList from "@/components/ItemList.vue";
import Clipboard from "clipboard";
import {checkSession} from "@/action/session";
import {useRouter} from "vue-router";
import {getSessionId} from "@/store/session";
import LoginDialog from "@/components/LoginDialog.vue";

const listBoxHeight = ref(window.innerHeight - 40)
const mjBoxHeight = ref(window.innerHeight - 150)
const fullImgHeight = ref(window.innerHeight - 60)
const showTaskDialog = ref(false)
const item = ref({})
const showLoginDialog = ref(false)
const isLogin = ref(false)

window.onresize = () => {
  listBoxHeight.value = window.innerHeight - 40
  mjBoxHeight.value = window.innerHeight - 150
}
const samplers = ["Euler a", "DPM++ 2S a Karras", "DPM++ 2M Karras", "DPM++ SDE Karras", "DPM++ 2M SDE Karras"]
const scaleAlg = ["Latent", "ESRGAN_4x", "R-ESRGAN 4x+", "SwinIR_4x", "LDSR"]
const params = ref({
  width: 1024,
  height: 1024,
  sampler: samplers[0],
  seed: -1,
  steps: 20,
  cfg_scale: 7,
  hd_fix: false,
  hd_redraw_rate: 0.7,
  hd_scale: 2,
  hd_scale_alg: scaleAlg[0],
  hd_steps: 0,
  prompt: "",
  neg_prompt: "nsfw, paintings,low quality,easynegative,ng_deepnegative ,lowres,bad anatomy,bad hands,bad feet",
})

const runningJobs = ref([])
const finishedJobs = ref([])
const router = useRouter()
// 检查是否有画同款的参数
const _params = router.currentRoute.value.params["copyParams"]
if (_params) {
  params.value = JSON.parse(_params)
}
const power = ref(0)
const sdPower = ref(0) // 画一张 SD 图片消耗算力

const socket = ref(null)
const userId = ref(0)
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

  const _socket = new WebSocket(host + `/api/sd/client?user_id=${userId.value}`);
  _socket.addEventListener('open', () => {
    socket.value = _socket;

    // 发送心跳消息
    sendHeartbeat()
  });

  _socket.addEventListener('message', event => {
    if (event.data instanceof Blob) {
      fetchRunningJobs()
      isOver.value = false
      page.value = 1
      fetchFinishJobs(page.value)
    }
  });

  _socket.addEventListener('close', () => {
    if (socket.value !== null) {
      connect()
    }
  });
}

const clipboard = ref(null)
onMounted(() => {
  initData()
  clipboard.value = new Clipboard('.copy-prompt-sd');
  clipboard.value.on('success', () => {
    ElMessage.success("复制成功！");
  })

  clipboard.value.on('error', () => {
    ElMessage.error('复制失败！');
  })

  httpGet("/api/config/get?key=system").then(res => {
    sdPower.value = res.data["sd_power"]
  }).catch(e => {
    ElMessage.error("获取系统配置失败：" + e.message)
  })
})

onUnmounted(() => {
  clipboard.value.destroy()
  socket.value = null
})


const initData = () => {
  checkSession().then(user => {
    power.value = user['power']
    userId.value = user.id
    isLogin.value = true
    fetchRunningJobs()
    fetchFinishJobs()
    connect()
  }).catch(() => {
    loading.value = false
  });
}

const fetchRunningJobs = () => {
  // 获取运行中的任务
  httpGet(`/api/sd/jobs?status=0`).then(res => {
    const jobs = res.data
    const _jobs = []
    for (let i = 0; i < jobs.length; i++) {
      if (jobs[i].progress === -1) {
        ElNotification({
          title: '任务执行失败',
          dangerouslyUseHTMLString: true,
          message: `任务ID：${jobs[i]['task_id']}<br />原因：${jobs[i]['err_msg']}`,
          type: 'error',
        })
        power.value += sdPower.value
        continue
      }
      _jobs.push(jobs[i])
    }
    runningJobs.value = _jobs
  }).catch(e => {
    ElMessage.error("获取任务失败：" + e.message)
  })
}

const handleScrollEnd = () => {
  if (isOver.value === true) {
    return
  }
  page.value += 1
  fetchFinishJobs(page.value)
}

const page = ref(1)
const pageSize = ref(15)
const isOver = ref(false)
const loading = ref(false)
// 获取已完成的任务
const fetchFinishJobs = (page) => {
  loading.value = true
  httpGet(`/api/sd/jobs?status=1&page=${page}&page_size=${pageSize.value}`).then(res => {
    if (res.data.length < pageSize.value) {
      isOver.value = true
    }
    if (page === 1) {
      finishedJobs.value = res.data
    } else {
      finishedJobs.value = finishedJobs.value.concat(res.data)
    }
    loading.value = false
  }).catch(e => {
    loading.value = false
    ElMessage.error("获取任务失败：" + e.message)
  })
}


// 创建绘图任务
const promptRef = ref(null)
const generate = () => {
  if (params.value.prompt === '') {
    promptRef.value.focus()
    return ElMessage.error("请输入绘画提示词！")
  }

  if (!isLogin.value) {
    showLoginDialog.value = true
    return
  }

  if (params.value.seed === '') {
    params.value.seed = -1
  }
  params.value.session_id = getSessionId()
  httpPost("/api/sd/image", params.value).then(() => {
    ElMessage.success("绘画任务推送成功，请耐心等待任务执行...")
    power.value -= sdPower.value
  }).catch(e => {
    ElMessage.error("任务推送失败：" + e.message)
  })
}

const showTask = (row) => {
  item.value = row
  showTaskDialog.value = true
}

const copyParams = (row) => {
  params.value = row.params
  showTaskDialog.value = false
}

const removeImage = (event, item) => {
  event.stopPropagation()
  ElMessageBox.confirm(
      '此操作将会删除任务和图片，继续操作码?',
      '删除提示',
      {
        confirmButtonText: '确认',
        cancelButtonText: '取消',
        type: 'warning',
      }
  ).then(() => {
    httpPost("/api/sd/remove", {id: item.id, img_url: item.img_url, user_id: userId.value}).then(() => {
      ElMessage.success("任务删除成功")
    }).catch(e => {
      ElMessage.error("任务删除失败：" + e.message)
    })
  }).catch(() => {
  })
}

// 发布图片到作品墙
const publishImage = (event, item, action) => {
  event.stopPropagation()
  let text = "图片发布"
  if (action === false) {
    text = "取消发布"
  }
  httpPost("/api/sd/publish", {id: item.id, action: action}).then(() => {
    ElMessage.success(text + "成功")
    item.publish = action
  }).catch(e => {
    ElMessage.error(text + "失败：" + e.message)
  })
}

</script>

<style lang="stylus">
@import "@/assets/css/image-sd.styl"
@import "@/assets/css/custom-scroll.styl"
</style>
