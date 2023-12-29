<template>
  <div class="page-sd" theme="dark">
    <div class="inner">
      <div class="sd-box" style="padding-top: 35px">
        <h2>Ai绘图</h2>
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
            <!-- 这是一个示例注释
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
             -->
            <!-- 这是一个示例注释
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
            -->
            <div class="param-line" v-loading="loading" element-loading-background="rgba(122, 122, 122, 0.8)">
              <el-input
                  v-model="params.prompt"
                  :autosize="{ minRows: 4, maxRows: 6 }"
                  type="textarea"
                  ref="promptRef"
                  placeholder="正向提示词，例如：“一个中国女孩，黑色的长发，穿着职业女装走在深圳夜晚的街头，看着街上的车流”，输入后点击翻译或翻译并重写按钮后点击生成。"
              />
            </div>

            <div style="padding: 10px">
              <el-button type="primary" @click="translatePrompt" size="small">
                <el-icon style="margin-right: 6px;font-size: 18px;">
                  <Refresh/>
                </el-icon>
                翻译
              </el-button>

              <el-tooltip
                  class="box-item"
                  effect="dark"
                  raw-content
                  content="使用 AI 翻译并重写提示词，<br/>增加更多细节，风格等描述"
                  placement="top-end"
              >
                <el-button type="success" @click="rewritePrompt" size="small">
                  <el-icon style="margin-right: 6px;font-size: 18px;">
                    <Refresh/>
                  </el-icon>
                  翻译并重写
                </el-button>
              </el-tooltip>
              <el-tag type="success" style="margin-left: 12px;">绘图可用额度：{{ imgCalls }}</el-tag>
            </div>
            <!--
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
            -->
          </el-form>
        </div>
        <div class="submit-btn">
          <el-button color="#47fff1" :dark="false" round @click="generate">立即生成</el-button>
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
              <ItemList :items="finishedJobs" v-if="finishedJobs.length > 0" width="480" :gap="20">
                <template #default="scope">
                  <div class="job-item animate" @click="showTask(scope.item)">
                    <el-image
                        :src="scope.item['img_url']+'?imageView2/1/w/480/h/480/q/75'"
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
    </div>
    <el-dialog v-model="showTaskDialog" title="绘画任务详情" :fullscreen="true">
      <el-row class="mobile-form-row">
        <el-col :span="24" class="mobile-form-col">
          <div class="img-container" :style="{maxHeight: fullImgHeight+'px'}">
            <el-image :src="item['img_url']" fit="contain"/>
          </div>
          <el-col :span="50">
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
        </el-col>
      </el-row>

    </el-dialog>
  </div>
</template>

<script setup>
import {onMounted, ref} from "vue"
import {DocumentCopy, InfoFilled, Orange, Picture, Refresh} from "@element-plus/icons-vue";
import {httpGet, httpPost} from "@/utils/http";
import {ElMessage, ElNotification} from "element-plus";
import ItemList from "@/components/ItemList.vue";
import Clipboard from "clipboard";
import {checkSession} from "@/action/session";
import {useRouter} from "vue-router";
import {getSessionId} from "@/store/session";

const listBoxHeight = ref(window.innerHeight - 40)
// const mjBoxHeight = ref(window.innerHeight - 150)
const fullImgHeight = ref(window.innerHeight - 60)
const showTaskDialog = ref(false)
const item = ref({})
const loading = ref(false)

window.onresize = () => {
  listBoxHeight.value = window.innerHeight - 40
//   mjBoxHeight.value = window.innerHeight - 150
}
const samplers = ["Euler a", "Euler", "DPM2 a Karras", "DPM++ 2S a Karras", "DPM++ 2M Karras", "DPM++ SDE Karras", "DPM2", "DPM2 a", "DPM++ 2S a", "DPM++ 2M", "DPM++ SDE", "DPM fast", "DPM adaptive",
  "LMS Karras", "DPM2 Karras", "DDIM", "PLMS", "UniPC", "LMS", "Heun",]
const scaleAlg = ["Latent", "ESRGAN_4x", "R-ESRGAN 4x+", "SwinIR_4x", "LDSR"]
const params = ref({
  width: 768,
  height: 1024,
  sampler: samplers[0],
  seed: -1,
  steps: 20,
  cfg_scale: 7,
  face_fix: false,
  hd_fix: false,
  hd_redraw_rate: 0.7,
  hd_scale: 2,
  hd_scale_alg: scaleAlg[0],
  hd_steps: 10,
  prompt: "",
  negative_prompt: "nsfw, paintings,low quality,easynegative,ng_deepnegative ,lowres,bad anatomy,bad hands,bad feet",
})

const runningJobs = ref([])
const finishedJobs = ref([])
const router = useRouter()
// 检查是否有画同款的参数
const _params = router.currentRoute.value.params["copyParams"]
if (_params) {
  params.value = JSON.parse(_params)
}

const imgCalls = ref(0)

const rewritePrompt = () => {
  loading.value = true
  httpPost("/api/prompt/rewrite", {"prompt": params.value.prompt}).then(res => {
    params.value.prompt = res.data
    loading.value = false
  }).catch(e => {
    loading.value = false
    ElMessage.error("翻译失败：" + e.message)
  })
}

const translatePrompt = () => {
  loading.value = true
  httpPost("/api/prompt/translate", {"prompt": params.value.prompt}).then(res => {
    params.value.prompt = res.data
    loading.value = false
  }).catch(e => {
    loading.value = false
    ElMessage.error("翻译失败：" + e.message)
  })
}

onMounted(() => {
  checkSession().then(user => {
    imgCalls.value = user['img_calls']

    fetchRunningJobs(user.id)
    fetchFinishJobs(user.id)

  }).catch(() => {
    router.push('/login')
  });

  const fetchRunningJobs = (userId) => {
    // 获取运行中的任务
    httpGet(`/api/sd/jobs?status=0&user_id=${userId}`).then(res => {
      const jobs = res.data
      const _jobs = []
      for (let i = 0; i < jobs.length; i++) {
        if (jobs[i].progress === -1) {
          ElNotification({
            title: '任务执行失败',
            message: "任务ID：" + jobs[i]['task_id'],
            type: 'error',
          })
          continue
        }
        _jobs.push(jobs[i])
      }
      runningJobs.value = _jobs

      setTimeout(() => fetchRunningJobs(userId), 5000)
    }).catch(e => {
      ElMessage.error("获取任务失败：" + e.message)
    })
  }

  // 获取已完成的任务
  const fetchFinishJobs = (userId) => {
    httpGet(`/api/sd/jobs?status=1&user_id=${userId}`).then(res => {
      finishedJobs.value = res.data
      setTimeout(() => fetchFinishJobs(userId), 5000)
    }).catch(e => {
      ElMessage.error("获取任务失败：" + e.message)
    })
  }

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
@import "@/assets/css/mobile/image-sd.css"
</style>
