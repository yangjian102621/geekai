<template>
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
              <el-form-item label="面部修复">
                <template #default>
                  <div class="form-item-inner">
                    <el-switch v-model="params.face_fix" style="--el-switch-on-color: #47fff1;"/>
                    <el-tooltip
                        effect="light"
                        content="仅对绘制人物图像有效果。"
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
                  placeholder="正向提示词，例如：A chinese girl walking in the middle of a cobblestone street"
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
                  v-model="params.negative_prompt"
                  :autosize="{ minRows: 4, maxRows: 6 }"
                  type="textarea"
                  placeholder="反向提示词"
              />
            </div>

            <div class="param-line pt">
              <el-form-item label="剩余次数">
                <template #default>
                  <el-tag type="info">{{ imgCalls }}</el-tag>
                </template>
              </el-form-item>
            </div>
          </el-form>
        </div>
        <div class="submit-btn">
          <el-button color="#47fff1" :dark="false" round @click="generate">立即生成</el-button>
        </div>
      </div>
      <div class="task-list-box">
        <div class="task-list-inner" :style="{ height: listBoxHeight + 'px' }">
          <h2>任务列表</h2>
          <div class="running-job-list">
            <ItemList :items="runningJobs" v-if="runningJobs.length > 0" width="240">
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
          <div class="finish-job-list">
            <ItemList :items="finishedJobs" v-if="finishedJobs.length > 0" width="240" :gap="16">
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
                </div>
              </template>
            </ItemList>
            <el-empty :image-size="100" v-else/>
          </div> <!-- end finish job list-->
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
                <el-icon class="copy-prompt" :data-clipboard-text="item.prompt">
                  <DocumentCopy/>
                </el-icon>
              </div>

            </div>

            <div class="info-line">
              <el-divider>
                反向提示词
              </el-divider>
              <div class="prompt">
                <span>{{ item.params.negative_prompt }}</span>
                <el-icon class="copy-prompt" :data-clipboard-text="item.params.negative_prompt">
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
</template>

<script setup>
import {onMounted, ref} from "vue"
import {DocumentCopy, InfoFilled, Orange, Picture} from "@element-plus/icons-vue";
import {httpGet, httpPost} from "@/utils/http";
import {ElMessage, ElNotification} from "element-plus";
import ItemList from "@/components/ItemList.vue";
import Clipboard from "clipboard";
import {checkSession} from "@/action/session";
import {useRouter} from "vue-router";
import {getSessionId, getUserToken} from "@/store/session";
import {removeArrayItem} from "@/utils/libs";

const listBoxHeight = ref(window.innerHeight - 40)
const mjBoxHeight = ref(window.innerHeight - 150)
const fullImgHeight = ref(window.innerHeight - 60)
const showTaskDialog = ref(false)
const item = ref({})

window.onresize = () => {
  listBoxHeight.value = window.innerHeight - 40
  mjBoxHeight.value = window.innerHeight - 150
}
const samplers = ["Euler a", "Euler", "DPM2 a Karras", "DPM++ 2S a Karras", "DPM++ 2M Karras", "DPM++ SDE Karras", "DPM2", "DPM2 a", "DPM++ 2S a", "DPM++ 2M", "DPM++ SDE", "DPM fast", "DPM adaptive",
  "LMS Karras", "DPM2 Karras", "DDIM", "PLMS", "UniPC", "LMS", "Heun",]
const scaleAlg = ["ESRGAN_4x", "R-ESRGAN 4x+", "SwinIR_4x", "LDSR"]
const params = ref({
  width: 1024,
  height: 1024,
  sampler: samplers[0],
  seed: -1,
  steps: 20,
  cfg_scale: 7,
  face_fix: false,
  hd_fix: false,
  hd_redraw_rate: 0.3,
  hd_scale: 2,
  hd_scale_alg: scaleAlg[0],
  hd_steps: 0,
  prompt: "A beautiful Chinese girl riding on a tiger",
  negative_prompt: "nsfw, paintings, cartoon, anime, sketches, low quality,easynegative,ng_deepnegative _v1 75t,(worst quality:2),(low quality:2),(normalquality:2),lowres,bad anatomy,bad hands,normal quality,((monochrome)),((grayscale)),((watermark))",
})

const runningJobs = ref([])
const finishedJobs = ref([])
const previewImgList = ref([])
const router = useRouter()
// 检查是否有画同款的参数
const _params = router.currentRoute.value.params["copyParams"]
if (_params) {
  params.value = JSON.parse(_params)
}

const socket = ref(null)
const imgCalls = ref(0)

const connect = () => {
  let host = process.env.VUE_APP_WS_HOST
  if (host === '') {
    if (location.protocol === 'https:') {
      host = 'wss://' + location.host;
    } else {
      host = 'ws://' + location.host;
    }
  }
  const _socket = new WebSocket(host + `/api/sd/client?session_id=${getSessionId()}&token=${getUserToken()}`);
  _socket.addEventListener('open', () => {
    socket.value = _socket;
  });

  _socket.addEventListener('message', event => {
    if (event.data instanceof Blob) {
      const reader = new FileReader();
      reader.readAsText(event.data, "UTF-8");
      reader.onload = () => {
        const data = JSON.parse(String(reader.result));
        let append = true
        if (data.progress === 100) { // 任务已完成
          for (let i = 0; i < finishedJobs.value.length; i++) {
            if (finishedJobs.value[i].id === data.id) {
              append = false
              break
            }
          }
          for (let i = 0; i < runningJobs.value.length; i++) {
            if (runningJobs.value[i].id === data.id) {
              runningJobs.value.splice(i, 1)
              break
            }
          }
          if (append) {
            finishedJobs.value.unshift(data)
          }
          previewImgList.value.unshift(data["img_url"])
        } else if (data.progress === -1) { // 任务执行失败
          ElNotification({
            title: '任务执行失败',
            message: "任务ID：" + data['task_id'],
            type: 'error',
          })
          runningJobs.value = removeArrayItem(runningJobs.value, data, (v1, v2) => v1.id === v2.id)

        } else { // 启动新的任务
          for (let i = 0; i < runningJobs.value.length; i++) {
            if (runningJobs.value[i].id === data.id) {
              append = false
              runningJobs.value[i] = data
              break
            }
          }
          if (append) {
            runningJobs.value.push(data)
          }
        }
      }
    }
  });

  _socket.addEventListener('close', () => {
    connect()
  });
}

onMounted(() => {
  checkSession().then(user => {
    imgCalls.value = user['img_calls']
    // 获取运行中的任务
    httpGet(`/api/sd/jobs?status=0&user_id=${user['id']}`).then(res => {
      runningJobs.value = res.data
    }).catch(e => {
      ElMessage.error("获取任务失败：" + e.message)
    })

    // 获取运行中的任务
    httpGet(`/api/sd/jobs?status=1&user_id=${user['id']}`).then(res => {
      finishedJobs.value = res.data
      previewImgList.value = []
      for (let index in finishedJobs.value) {
        previewImgList.value.push(finishedJobs.value[index]["img_url"])
      }
    }).catch(e => {
      ElMessage.error("获取任务失败：" + e.message)
    })

    // 连接 socket
    connect();
  }).catch(() => {
    router.push('/login')
  });

  const clipboard = new Clipboard('.copy-prompt');
  clipboard.on('success', () => {
    ElMessage.success("复制成功！");
  })

  clipboard.on('error', () => {
    ElMessage.error('复制失败！');
  })
})


// 创建绘图任务
const promptRef = ref(null)
const generate = () => {
  if (params.value.prompt === '') {
    promptRef.value.focus()
    return ElMessage.error("请输入绘画提示词！")
  }
  if (params.value.seed === '') {
    params.value.seed = -1
  }
  params.value.session_id = getSessionId()
  httpPost("/api/sd/image", params.value).then(() => {
    ElMessage.success("绘画任务推送成功，请耐心等待任务执行...")
    imgCalls.value -= 1
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

</script>

<style lang="stylus">
@import "@/assets/css/image-sd.styl"
@import "@/assets/css/custom-scroll.styl"
</style>
