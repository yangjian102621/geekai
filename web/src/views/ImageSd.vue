<template>
  <div>
    <div class="page-sd">
      <div class="inner custom-scroll">
        <div class="sd-box">
          <h2>Stable Diffusion 创作中心</h2>

          <div class="sd-params">
            <el-form :model="params" label-width="80px" label-position="left">
              <div class="param-line" style="padding-top: 10px">
                <el-form-item label="采样方法">
                  <template #default>
                    <div class="form-item-inner">
                      <el-select v-model="params.sampler" style="width: 150px">
                        <el-option
                          v-for="item in samplers"
                          :label="item"
                          :value="item"
                          :key="item"
                        />
                      </el-select>
                      <el-tooltip
                        content="出图效果比较好的一般是 Euler 和 DPM 系列算法"
                        raw-content
                        placement="right"
                      >
                        <el-icon class="info-icon">
                          <InfoFilled />
                        </el-icon>
                      </el-tooltip>
                    </div>
                  </template>
                </el-form-item>
              </div>

              <div class="param-line" style="padding-top: 10px">
                <el-form-item label="采样调度">
                  <template #default>
                    <div class="form-item-inner">
                      <el-select v-model="params.scheduler" style="width: 150px">
                        <el-option
                          v-for="item in schedulers"
                          :label="item"
                          :value="item"
                          :key="item"
                        />
                      </el-select>
                      <el-tooltip content="推荐自动或者 Karras" raw-content placement="right">
                        <el-icon class="info-icon">
                          <InfoFilled />
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
                          <el-input v-model.number="params.width" placeholder="图片宽度" />
                        </el-col>
                        <el-col :span="12">
                          <el-input v-model.number="params.height" placeholder="图片高度" />
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
                      <el-input v-model.number="params.steps" />
                      <el-tooltip
                        content="值越大则代表细节越多，同时也意味着出图速度越慢"
                        raw-content
                        placement="right"
                      >
                        <el-icon class="info-icon">
                          <InfoFilled />
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
                      <el-input v-model.number="params.cfg_scale" />
                      <el-tooltip
                        content="提示词引导系数，图像在多大程度上服从提示词<br/> 较低值会产生更有创意的结果"
                        raw-content
                        placement="right"
                      >
                        <el-icon class="info-icon">
                          <InfoFilled />
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
                      <el-input v-model.number="params.seed" />
                      <el-tooltip
                        content="随机数种子，相同的种子会得到相同的结果<br/> 设置为 -1 则每次随机生成种子"
                        raw-content
                        placement="right"
                      >
                        <el-icon class="info-icon">
                          <InfoFilled />
                        </el-icon>
                      </el-tooltip>

                      <el-tooltip content="使用随机数" raw-content placement="right">
                        <el-icon @click="params.seed = -1" class="info-icon">
                          <Orange />
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
                      <el-switch
                        v-model="params.hd_fix"
                        style="--el-switch-on-color: #47fff1"
                        size="large"
                      />
                      <el-tooltip
                        content="先以较小的分辨率生成图像，接着方法图像<br />然后在不更改构图的情况下再修改细节"
                        raw-content
                        placement="right"
                      >
                        <el-icon style="margin-left: 10px; top: 12px">
                          <InfoFilled />
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
                        <el-slider
                          v-model.number="params.hd_redraw_rate"
                          :max="1"
                          :step="0.1"
                          style="width: 180px; --el-slider-main-bg-color: #47fff1"
                        />
                        <el-tooltip
                          content="决定算法对图像内容的影响程度<br />较大的值将得到越有创意的图像"
                          raw-content
                          placement="right"
                        >
                          <el-icon class="info-icon">
                            <InfoFilled />
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
                        <el-select v-model="params.hd_scale_alg" style="width: 176px">
                          <el-option
                            v-for="item in scaleAlg"
                            :label="item"
                            :value="item"
                            :key="item"
                          />
                        </el-select>
                        <el-tooltip
                          content="高清修复放大算法，主流算法有Latent和ESRGAN_4x"
                          raw-content
                          placement="right"
                        >
                          <el-icon class="info-icon">
                            <InfoFilled />
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
                        <el-input v-model.number="params.hd_scale" />
                        <el-tooltip
                          content="随机数种子，相同的种子会得到相同的结果<br/> 设置为 -1 则每次随机生成种子"
                          raw-content
                          placement="right"
                        >
                          <el-icon class="info-icon">
                            <InfoFilled />
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
                        <el-input v-model.number="params.hd_steps" />
                        <el-tooltip
                          content="重绘迭代步数，如果设置为0，则设置跟原图相同的迭代步数"
                          raw-content
                          placement="right"
                        >
                          <el-icon class="info-icon">
                            <InfoFilled />
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
                  maxlength="2000"
                  placeholder="请在此输入绘画提示词，您也可以点击下面的提示词助手生成绘画提示词"
                  v-loading="isGenerating"
                />
              </div>

              <el-row class="text-info">
                <el-button
                  class="generate-btn"
                  size="small"
                  @click="generatePrompt"
                  color="#5865f2"
                  :disabled="isGenerating"
                >
                  <i class="iconfont icon-chuangzuo" style="margin-right: 5px"></i>
                  <span>生成专业绘画指令</span>
                </el-button>
              </el-row>

              <div class="param-line pt">
                <span>反向提示词：</span>
                <el-tooltip content="不希望出现的元素，下面给了默认的起手式" placement="right">
                  <el-icon class="info-icon">
                    <InfoFilled />
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
                <el-row :gutter="10">
                  <el-text type="primary"
                    >单次绘图消耗
                    <el-text type="warning">{{ sdPower }}算力，</el-text>
                  </el-text>
                  <el-text type="primary"
                    >当前可用 <el-text type="warning"> {{ power }}算力</el-text></el-text
                  >
                </el-row>
              </div>
            </el-form>
          </div>
          <div class="submit-btn">
            <el-button type="primary" :dark="false" round @click="generate">立即生成</el-button>
          </div>
        </div>
        <div class="task-list-box pl-6 pr-6 pb-4 pt-4 h-dvh">
          <div class="task-list-inner" :style="{ height: listBoxHeight + 'px' }">
            <div class="job-list-box">
              <h2 class="text-xl">任务列表</h2>
              <task-list :list="runningJobs" />
              <template v-if="finishedJobs.length > 0">
                <h2 class="text-xl">创作记录</h2>
                <div class="finish-job-list mt-3">
                  <div v-if="finishedJobs.length > 0">
                    <Waterfall
                      :list="finishedJobs"
                      :row-key="waterfallOptions.rowKey"
                      :gutter="waterfallOptions.gutter"
                      :has-around-gutter="waterfallOptions.hasAroundGutter"
                      :width="waterfallOptions.width"
                      :breakpoints="waterfallOptions.breakpoints"
                      :img-selector="waterfallOptions.imgSelector"
                      :background-color="waterfallOptions.backgroundColor"
                      :animation-effect="waterfallOptions.animationEffect"
                      :animation-duration="waterfallOptions.animationDuration"
                      :animation-delay="waterfallOptions.animationDelay"
                      :animation-cancel="waterfallOptions.animationCancel"
                      :lazyload="waterfallOptions.lazyload"
                      :load-props="waterfallOptions.loadProps"
                      :cross-origin="waterfallOptions.crossOrigin"
                      :align="waterfallOptions.align"
                      :is-loading="loading"
                      :is-over="isOver"
                      @afterRender="loading = false"
                    >
                      <template #default="{ item, url }">
                        <div
                          class="bg-gray-900 rounded-lg shadow-md overflow-hidden transition-all duration-300 ease-linear hover:shadow-md hover:shadow-purple-800 group"
                        >
                          <div class="overflow-hidden rounded-lg">
                            <LazyImg
                              :url="url"
                              v-if="item.progress === 100"
                              class="cursor-pointer transition-all duration-300 ease-linear group-hover:scale-105"
                              @click="showTask(item)"
                            />
                            <el-image v-else-if="item.progress === 101">
                              <template #error>
                                <div class="image-slot">
                                  <div class="err-msg-container">
                                    <div class="title">任务失败</div>
                                    <div class="opt">
                                      <el-popover
                                        title="错误详情"
                                        trigger="click"
                                        :width="250"
                                        :content="item['err_msg']"
                                        placement="top"
                                      >
                                        <template #reference>
                                          <el-button type="info">详情</el-button>
                                        </template>
                                      </el-popover>
                                      <el-button type="danger" @click="removeImage(item)"
                                        >删除</el-button
                                      >
                                    </div>
                                  </div>
                                </div>
                              </template>
                            </el-image>
                          </div>
                          <div
                            class="px-4 pt-2 pb-4 border-t border-t-gray-800"
                            v-if="item.progress === 100"
                          >
                            <div
                              class="pt-3 flex justify-center items-center border-t border-t-gray-600 border-opacity-50"
                            >
                              <div class="flex">
                                <el-tooltip content="取消分享" placement="top" v-if="item.publish">
                                  <el-button
                                    type="warning"
                                    @click="publishImage(item, false)"
                                    circle
                                  >
                                    <i class="iconfont icon-cancel-share"></i>
                                  </el-button>
                                </el-tooltip>
                                <el-tooltip content="分享" placement="top" v-else>
                                  <el-button
                                    type="success"
                                    @click="publishImage(item, true)"
                                    circle
                                  >
                                    <i class="iconfont icon-share-bold"></i>
                                  </el-button>
                                </el-tooltip>

                                <el-tooltip content="复制提示词" placement="top">
                                  <el-button
                                    type="info"
                                    circle
                                    class="copy-prompt"
                                    :data-clipboard-text="item.prompt"
                                  >
                                    <i class="iconfont icon-file"></i>
                                  </el-button>
                                </el-tooltip>
                                <el-tooltip content="删除" placement="top">
                                  <el-button
                                    type="danger"
                                    :icon="Delete"
                                    @click="removeImage(item)"
                                    circle
                                  />
                                </el-tooltip>
                              </div>
                            </div>
                          </div>
                        </div>
                      </template>
                    </Waterfall>

                    <div class="flex justify-center py-10">
                      <img
                        :src="waterfallOptions.loadProps.loading"
                        class="max-w-[50px] max-h-[50px]"
                        v-if="loading"
                      />
                      <div v-else>
                        <button
                          class="px-5 py-2 rounded-full bg-purple-700 text-md text-white cursor-pointer hover:bg-purple-800 transition-all duration-300"
                          @click="fetchFinishJobs"
                          v-if="!isOver"
                        >
                          加载更多
                        </button>
                        <div class="no-more-data" v-else>
                          <span class="text-gray-500 mr-2">没有更多数据了</span>
                          <i class="iconfont icon-face"></i>
                        </div>
                      </div>
                    </div>
                  </div>
                  <el-empty :image-size="100" v-else :image="nodata" description="暂无记录" />
                </div>
              </template>

              <!-- end finish job list-->
            </div>
          </div>
          <back-top :right="30" :bottom="30" />
        </div>
        <!-- end task list box -->
      </div>

      <!-- 任务详情弹框 -->
      <sd-task-view
        v-model="showTaskDialog"
        :data="item"
        @drawSame="copyParams"
        @close="showTaskDialog = false"
      />
    </div>
  </div>
</template>

<script setup>
import { nextTick, onMounted, onUnmounted, ref } from 'vue'
import { Delete, DocumentCopy, InfoFilled, Orange } from '@element-plus/icons-vue'
import nodata from '@/assets/img/no-data.png'

import { httpGet, httpPost } from '@/utils/http'
import { ElMessage, ElMessageBox } from 'element-plus'
import Clipboard from 'clipboard'
import { checkSession, getSystemInfo } from '@/store/cache'
import { useRouter } from 'vue-router'
import { getSessionId } from '@/store/session'
import { useSharedStore } from '@/store/sharedata'
import TaskList from '@/components/TaskList.vue'
import BackTop from '@/components/BackTop.vue'
import { showMessageError } from '@/utils/dialog'
import SdTaskView from '@/components/SdTaskView.vue'
import { LazyImg, Waterfall } from 'vue-waterfall-plugin-next'
import 'vue-waterfall-plugin-next/dist/style.css'

const listBoxHeight = ref(0)
// const paramBoxHeight = ref(0)
const showTaskDialog = ref(false)
const item = ref({})
const isLogin = ref(false)
const loading = ref(true)
const store = useSharedStore()
const waterfallOptions = store.waterfallOptions

const resizeElement = function () {
  listBoxHeight.value = window.innerHeight - 80
  // paramBoxHeight.value = window.innerHeight - 200
}
resizeElement()
window.onresize = () => {
  resizeElement()
}
const samplers = [
  'Euler a',
  'DPM++ 2S a',
  'DPM++ 2M',
  'DPM++ SDE',
  'DPM++ 2M SDE',
  'UniPC',
  'Restart',
]
const schedulers = ['Automatic', 'Karras', 'Exponential', 'Uniform']
const scaleAlg = ['Latent', 'ESRGAN_4x', 'R-ESRGAN 4x+', 'SwinIR_4x', 'LDSR']
const params = ref({
  width: 1024,
  height: 1024,
  sampler: samplers[0],
  scheduler: schedulers[0],
  seed: -1,
  steps: 20,
  cfg_scale: 7,
  hd_fix: false,
  hd_redraw_rate: 0.7,
  hd_scale: 2,
  hd_scale_alg: scaleAlg[0],
  hd_steps: 0,
  prompt: '',
  neg_prompt:
    'nsfw, paintings,low quality,easynegative,ng_deepnegative ,lowres,bad anatomy,bad hands,bad feet',
})

const runningJobs = ref([])
const finishedJobs = ref([])
const allowPulling = ref(true) // 是否允许轮询
const tastPullHandler = ref(null)
const router = useRouter()
// 检查是否有画同款的参数
const _params = router.currentRoute.value.params['copyParams']
if (_params) {
  params.value = JSON.parse(_params)
}
const power = ref(0)
const sdPower = ref(0) // 画一张 SD 图片消耗算力

const userId = ref(0)
const clipboard = ref(null)
onMounted(() => {
  initData()
  clipboard.value = new Clipboard('.copy-prompt-sd')
  clipboard.value.on('success', () => {
    ElMessage.success('复制成功！')
  })

  clipboard.value.on('error', () => {
    ElMessage.error('复制失败！')
  })

  getSystemInfo()
    .then((res) => {
      sdPower.value = res.data.sd_power
      params.value.neg_prompt = res.data.sd_neg_prompt
    })
    .catch((e) => {
      ElMessage.error('获取系统配置失败：' + e.message)
    })
})

onUnmounted(() => {
  clipboard.value.destroy()
  if (tastPullHandler.value) {
    clearInterval(tastPullHandler.value)
  }
})

const initData = () => {
  checkSession()
    .then((user) => {
      power.value = user['power']
      userId.value = user.id
      isLogin.value = true
      page.value = 0
      fetchRunningJobs()
      fetchFinishJobs()

      tastPullHandler.value = setInterval(() => {
        if (allowPulling.value) {
          fetchRunningJobs()
        }
      }, 5000)
    })
    .catch(() => {})
}

const fetchRunningJobs = () => {
  if (!isLogin.value) {
    return
  }

  // 获取运行中的任务
  httpGet(`/api/sd/jobs?finish=0`)
    .then((res) => {
      if (runningJobs.value.length !== res.data.items.length) {
        page.value = 0
        fetchFinishJobs()
      }
      if (runningJobs.value.length === 0) {
        allowPulling.value = false
      }
      runningJobs.value = res.data.items
    })
    .catch((e) => {
      ElMessage.error('获取任务失败：' + e.message)
    })
}

const page = ref(0)
const pageSize = ref(20)
const isOver = ref(false)
// 获取已完成的任务
const fetchFinishJobs = () => {
  if (!isLogin.value || isOver.value === true) {
    return
  }
  loading.value = true
  page.value = page.value + 1

  httpGet(`/api/sd/jobs?finish=1&page=${page.value}&page_size=${pageSize.value}`)
    .then((res) => {
      if (res.data.items.length < pageSize.value) {
        isOver.value = true
        loading.value = false
      }
      const imageList = res.data.items
      for (let i = 0; i < imageList.length; i++) {
        imageList[i]['img_thumb'] = imageList[i]['img_url'] + '?imageView2/4/w/300/h/0/q/75'
      }
      if (page.value === 1) {
        finishedJobs.value = imageList
      } else {
        finishedJobs.value = finishedJobs.value.concat(imageList)
      }
    })
    .catch((e) => {
      ElMessage.error('获取任务失败：' + e.message)
      loading.value = false
    })
}

// 创建绘图任务
const promptRef = ref(null)
const generate = () => {
  if (params.value.prompt === '') {
    promptRef.value.focus()
    return ElMessage.error('请输入绘画提示词！')
  }

  if (!isLogin.value) {
    store.setShowLoginDialog(true)
    return
  }

  if (!params.value.seed) {
    params.value.seed = -1
  }
  params.value.session_id = getSessionId()
  httpPost('/api/sd/image', params.value)
    .then(() => {
      ElMessage.success('绘画任务推送成功，请耐心等待任务执行...')
      power.value -= sdPower.value
      allowPulling.value = true
      runningJobs.value.push({
        progress: 0,
      })
      isOver.value = false
    })
    .catch((e) => {
      ElMessage.error('任务推送失败：' + e.message)
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

const removeImage = (item) => {
  ElMessageBox.confirm('此操作将会删除任务和图片，继续操作码?', '删除提示', {
    confirmButtonText: '确认',
    cancelButtonText: '取消',
    type: 'warning',
  })
    .then(() => {
      httpGet('/api/sd/remove', { id: item.id })
        .then(() => {
          ElMessage.success('任务删除成功')
          page.value = 0
          isOver.value = false
          fetchFinishJobs()
        })
        .catch((e) => {
          ElMessage.error('任务删除失败：' + e.message)
        })
    })
    .catch(() => {})
}

// 发布图片到作品墙
const publishImage = (item, action) => {
  let text = '图片发布'
  if (action === false) {
    text = '取消发布'
  }
  httpGet('/api/sd/publish', { id: item.id, action: action })
    .then(() => {
      ElMessage.success(text + '成功')
      item.publish = action
      page.value = 0
      isOver.value = false
      item.publish = action
    })
    .catch((e) => {
      ElMessage.error(text + '失败：' + e.message)
    })
}

const isGenerating = ref(false)
const generatePrompt = () => {
  if (params.value.prompt === '') {
    return showMessageError('请输入原始提示词')
  }
  isGenerating.value = true
  httpPost('/api/prompt/image', { prompt: params.value.prompt })
    .then((res) => {
      params.value.prompt = res.data
      isGenerating.value = false
    })
    .catch((e) => {
      showMessageError('生成提示词失败：' + e.message)
      isGenerating.value = false
    })
}
</script>

<style lang="stylus">
@import '@/assets/css/image-sd.styl';
@import '@/assets/css/custom-scroll.styl';
</style>
