<template>
  <div>
    <div class="page-image">
      <div class="inner custom-scroll">
        <div class="image-box">
          <h2 class="!text-[#252f76] py-3">AI图像生成</h2>

          <div class="sd-params">
            <el-form :model="params" label-width="80px" label-position="top">
              <div class="param-line pt-1">
                <el-form-item label="生图模型">
                  <template #default>
                    <el-select
                      v-model="selectedModel"
                      placeholder="请选择模型"
                      @change="changeModel"
                    >
                      <el-option v-for="v in models" :label="v.name" :value="v" :key="v.value" />
                    </el-select>
                  </template>
                </el-form-item>
              </div>

              <div class="param-line">
                <el-form-item>
                  <template #label>
                    <div class="flex items-center justify-between w-full">
                      <span>图片比例</span>
                    </div>
                  </template>
                  <template #default>
                    <div class="form-item-inner flex flex-wrap gap-2">
                      <div
                        v-for="item in radioAspectList"
                        :key="item.value"
                        class="w-[56px] h-[66px] rounded-lg border cursor-pointer flex flex-col items-center justify-center bg-white transition-all duration-200"
                        :class="
                          selectedAspect === item.value
                            ? 'border-[#6366f1] bg-[#eef2ff]'
                            : 'border-[#e5e7f5]'
                        "
                        @click="changeAspect(item.value)"
                      >
                        <div class="flex items-center justify-center h-[28px]">
                          <i class="iconfont mr-1" :class="item.icon"></i>
                        </div>
                        <div class="text-xs font-semibold text-[#252f76] leading-tight">
                          {{ item.value }}
                        </div>
                        <div class="text-[11px] text-gray-400 leading-tight line-clamp-1">
                          {{ item.label }}
                        </div>
                      </div>
                    </div>
                  </template>
                </el-form-item>
              </div>

              <div class="param-line">
                <el-form-item>
                  <template #label>
                    <div class="flex items-center justify-between w-full">
                      <span>图片尺寸</span>
                    </div>
                  </template>
                  <template #default>
                    <div class="form-item-inner flex flex-col gap-2">
                      <div class="flex flex-wrap gap-2">
                        <div
                          v-for="item in sizeLevels"
                          :key="item.value"
                          class="w-[72px] h-[50px] rounded-lg border cursor-pointer flex flex-col justify-center items-center bg-white transition-all duration-200"
                          :class="
                            sizeLevel === item.value
                              ? 'border-[#6366f1] bg-[#eef2ff]'
                              : 'border-[#e5e7f5]'
                          "
                          @click="changeSizeLevel(item.value)"
                        >
                          <template v-if="item.value === 'other'">
                            <div class="text-[16px] font-semibold !text-purple-600">
                              {{ item.label }}
                            </div>
                          </template>
                          <template v-else>
                            <div class="text-[16px] font-semibold !text-purple-600 h-[20px]">
                              {{ item.label }}
                            </div>
                            <div class="text-[11px] text-gray-400">{{ item.desc }}</div>
                          </template>
                        </div>
                      </div>
                      <el-input
                        v-if="sizeLevel === 'other'"
                        v-model="params.size"
                        placeholder="宽×高，如 1024x1024"
                        clearable
                      />
                    </div>
                  </template>
                </el-form-item>
              </div>

              <div class="mt-2 mb-2">
                <label class="text-gray-700 font-semibold">参考图(可选)</label>
                <div class="py-2">
                  <ImageUpload
                    v-model="params.image"
                    :max-count="5"
                    :max-size="20"
                    :multiple="true"
                  />
                </div>
              </div>

              <div class="param-line prompt-input-wrap">
                <el-form-item label="绘图提示词">
                  <ReferenceInput
                    v-model="params.prompt"
                    :image-list="params.image || []"
                    :autosize="{ minRows: 4, maxRows: 6 }"
                    :maxlength="1024"
                    :show-word-limit="true"
                    :loading="promptGenerating"
                    placeholder="请在此输入绘画提示词，输入 @ 可引用参考图（图1、图2...）"
                  />
                </el-form-item>
              </div>

              <div class="flex justify-end pt-2 pr-2">
                <el-button @click="generatePrompt" type="primary" :loading="promptGenerating">
                  <span v-if="!promptGenerating">
                    <i class="iconfont icon-chuangzuo"></i>
                    生成专业绘画指令
                  </span>
                  <span v-else>生成中...</span>
                </el-button>
              </div>
            </el-form>
          </div>
          <div class="py-4">
            <button
              class="w-full py-3 bg-gradient-to-r from-blue-500 to-purple-600 text-white rounded-xl disabled:from-gray-400 disabled:to-gray-400 disabled:cursor-not-allowed hover:from-blue-600 hover:to-purple-700 transition-all duration-200 flex items-center justify-center space-x-2 text-base"
              type="button"
              @click="generate"
            >
              <i v-if="isGenerating" class="iconfont icon-loading animate-spin"></i>
              <i v-else class="iconfont icon-chuangzuo"></i>
              <span v-if="isGenerating">创作中...</span>
              <span v-else>立即生成({{ imagePower }}算力)</span>
            </button>
          </div>
        </div>
        <div class="task-list-box pl-6 pr-6 pb-4 pt-4 h-dvh">
          <div class="task-list-inner">
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
                              @click="previewImg(item)"
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

                                <el-tooltip content="详情" placement="top">
                                  <el-button type="info" circle @click="showDetail(item)">
                                    <i class="iconfont icon-info"></i>
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

                                <el-tooltip content="下载" placement="top">
                                  <el-button
                                    type="primary"
                                    circle
                                    :icon="Download"
                                    @click="downloadImage(item)"
                                    :loading="item.downloading"
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
                  <el-empty :image-size="100" :image="nodata" description="暂无记录" v-else />
                </div>
              </template>
              <!-- end finish job list-->
            </div>
          </div>
          <back-top :right="30" :bottom="30" />
        </div>
        <!-- end task list box -->
      </div>
    </div>

    <el-image-viewer
      @close="
        () => {
          previewURL = ''
        }
      "
      v-if="previewURL !== ''"
      :url-list="[previewURL]"
    />

    <!-- 详情弹窗 -->
    <el-dialog
      v-model="detailDialogVisible"
      title="任务详情"
      width="600px"
      :close-on-click-modal="false"
    >
      <div class="detail-content" v-if="currentDetail">
        <el-descriptions :column="1" border>
          <el-descriptions-item label="提示词">
            <div>
              <span>{{ currentDetail.prompt }}</span>
              <el-tooltip content="复制提示词" placement="top">
                <i
                  class="iconfont icon-copy ml-2 cursor-pointer copy-prompt-detail"
                  :data-clipboard-text="currentDetail.prompt"
                />
              </el-tooltip>
            </div>
          </el-descriptions-item>

          <el-descriptions-item
            label="生成的图片"
            v-if="currentDetail.progress === 100 && currentDetail.img_url"
          >
            <el-image
              :src="getThumbURL(currentDetail.img_url, 200, 200)"
              :preview-src-list="[currentDetail.img_url]"
              fit="cover"
              style="width: 200px; height: 200px"
            />
          </el-descriptions-item>

          <el-descriptions-item
            label="参考图"
            v-if="currentDetail.params?.image && currentDetail.params.image.length > 0"
          >
            <div class="reference-images">
              <el-image
                v-for="(img, idx) in currentDetail.params.image"
                :key="idx"
                :src="getThumbURL(img, 100, 100)"
                :preview-src-list="currentDetail.params.image"
                :initial-index="idx"
                fit="cover"
                style="width: 100px; height: 100px; margin-right: 10px"
              />
            </div>
          </el-descriptions-item>

          <el-descriptions-item label="生图模型">
            {{ currentDetail.params?.model_name || '-' }}
          </el-descriptions-item>

          <el-descriptions-item label="消耗算力">
            {{ currentDetail.power || 0 }}
          </el-descriptions-item>

          <el-descriptions-item label="图片比例">
            {{ currentDetail.params?.aspect_ratio || '-' }}
          </el-descriptions-item>

          <el-descriptions-item label="图片尺寸">
            {{ currentDetail.params?.size || '-' }}
          </el-descriptions-item>

          <el-descriptions-item label="创建时间">
            {{ dateFormat(currentDetail.created_at) }}
          </el-descriptions-item>

          <el-descriptions-item
            label="错误信息"
            v-if="currentDetail.progress === 101 && currentDetail.err_msg"
          >
            <el-text type="danger">{{ currentDetail.err_msg }}</el-text>
          </el-descriptions-item>
        </el-descriptions>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import nodata from '@/assets/img/no-data.png'

import BackTop from '@/components/BackTop.vue'
import ImageUpload from '@/components/ImageUpload.vue'
import ReferenceInput from '@/components/ReferenceInput.vue'
import TaskList from '@/components/TaskList.vue'
import { checkSession } from '@/store/cache'
import { useSharedStore } from '@/store/sharedata'
import { showMessageError } from '@/utils/dialog'
import { downloadFile, httpGet, httpPost } from '@/utils/http'
import { dateFormat, getThumbURL } from '@/utils/libs'
import Clipboard from 'clipboard'
import { Delete, Download, List } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { LazyImg, Waterfall } from 'vue-waterfall-plugin-next'
import 'vue-waterfall-plugin-next/dist/style.css'

const listBoxHeight = ref(0)
const isLogin = ref(false)
const loading = ref(true)
const isOver = ref(false)
const previewURL = ref('')
const store = useSharedStore()
const models = ref([])
const waterfallOptions = store.waterfallOptions
const resizeElement = function () {
  listBoxHeight.value = window.innerHeight - 58
}

resizeElement()
window.onresize = () => {
  resizeElement()
}

// 为了兼容 size 和 aspect_ratio 参数，label 对应的是 aspect_ratio 的值，size 是预估值
const radioAspects = ref({
  '1:1': { value: '1:1', label: '方形', icon: 'icon-aspect_1_1' },
  '4:3': { value: '4:3', label: '横向', icon: 'icon-aspect_4_3' },
  '3:4': { value: '3:4', label: '竖向', icon: 'icon-aspect_3_4' },
  '16:9': { value: '16:9', label: '宽屏', icon: 'icon-aspect_16_9' },
  '9:16': { value: '9:16', label: '竖屏', icon: 'icon-aspect_9_16' },
  '3:2': { value: '3:2', label: '横向', icon: 'icon-aspect_3_2' },
  '2:3': { value: '2:3', label: '竖向', icon: 'icon-aspect_2_3' },
  '4:5': { value: '4:5', label: '竖向', icon: 'icon-aspect_4_5' },
  '5:4': { value: '5:4', label: '横向', icon: 'icon-aspect_5_4' },
  '21:9': { value: '21:9', label: '超宽', icon: 'icon-aspect_16_9' },
})
const radioAspectList = computed(() => Object.values(radioAspects.value))
const params = ref({
  aspect_ratio: '1:1',
  size: '1K',
  prompt: '',
})

const selectedAspect = ref(params.value.aspect_ratio || '1:1')
const sizeLevel = ref('1K')
const sizeLevels = [
  { value: '1K', label: '1K', desc: '≈1024px' },
  { value: '2K', label: '2K', desc: '≈2048px' },
  { value: '4K', label: '4K', desc: '≈4096px' },
  { value: 'other', label: '其他', desc: '' },
]

const finishedJobs = ref([])
const runningJobs = ref([])
const allowPulling = ref(true) // 是否允许轮询
const downloadPulling = ref(false) // 下载轮询
const tastPullHandler = ref(null)
const downloadPullHandler = ref(null)
const userPower = ref(0)
const imagePower = ref(0)
const userId = ref(0)
const selectedModel = ref(null)
const detailDialogVisible = ref(false)
const currentDetail = ref(null)
const clipboard = ref(null)

onMounted(() => {
  initData()

  // 获取模型列表
  httpGet('/api/image/models')
    .then((res) => {
      models.value = res.data
      selectedModel.value = models.value[0]
      params.value.model_id = selectedModel.value.id
      changeModel(selectedModel.value)
    })
    .catch((e) => {
      showMessageError('获取模型列表失败：' + e.message)
    })

  clipboard.value = new Clipboard('.copy-prompt-detail')
  clipboard.value.on('success', () => {
    ElMessage.success('复制成功！')
  })

  clipboard.value.on('error', () => {
    ElMessage.error('复制失败！')
  })
})

onUnmounted(() => {
  if (tastPullHandler.value) {
    clearInterval(tastPullHandler.value)
  }
  if (downloadPullHandler.value) {
    clearInterval(downloadPullHandler.value)
  }
  if (clipboard.value) {
    clipboard.value.destroy()
  }
})

const changeAspect = (value) => {
  selectedAspect.value = value
  params.value.aspect_ratio = value
}

const changeSizeLevel = (level) => {
  sizeLevel.value = level
  if (level === 'other') {
    if (['1K', '2K', '4K'].includes(params.value.size)) {
      params.value.size = ''
    }
    return
  }
  params.value.size = level
}

const initData = () => {
  checkSession()
    .then((user) => {
      userPower.value = user['power']
      userId.value = user.id
      isLogin.value = true
      page.value = 0
      fetchRunningJobs()
      fetchFinishJobs()

      // 轮询运行中任务
      tastPullHandler.value = setInterval(() => {
        if (allowPulling.value) {
          fetchRunningJobs()
        }
      }, 5000)

      // 图片下载轮询
      downloadPullHandler.value = setInterval(() => {
        if (downloadPulling.value) {
          page.value = 0
          fetchFinishJobs()
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
  httpGet(`/api/image/jobs?finish=false`)
    .then((res) => {
      // 如果任务有更新，则更新已完成任务列表
      if (res.data.items && res.data.items.length !== runningJobs.value.length) {
        page.value = 0
        fetchFinishJobs()
      }
      if (res.data.items.length > 0) {
        runningJobs.value = res.data.items
      } else {
        allowPulling.value = false
        runningJobs.value = []
      }
    })
    .catch((e) => {
      ElMessage.error('获取任务失败：' + e.message)
    })
}

const page = ref(1)
const pageSize = ref(15)
// 获取已完成的任务
const fetchFinishJobs = () => {
  if (!isLogin.value) {
    return
  }

  loading.value = true
  page.value = page.value + 1

  httpGet(`/api/image/jobs?finish=true&page=${page.value}&page_size=${pageSize.value}`)
    .then((res) => {
      if (res.data.items.length < pageSize.value) {
        isOver.value = true
        loading.value = false
      }
      const imageList = res.data.items
      let needPulling = false
      for (let i = 0; i < imageList.length; i++) {
        if (imageList[i]['img_url']) {
          imageList[i]['img_thumb'] = getThumbURL(imageList[i]['img_url'], 300, 0)
        } else if (imageList[i].progress === 100) {
          needPulling = true
          imageList[i]['img_thumb'] = waterfallOptions.loadProps.loading
        }
      }
      // 如果当前是第一页，则开启图片下载轮询
      if (page.value === 1) {
        downloadPulling.value = needPulling
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

const isGenerating = ref(false)
const generate = () => {
  if (isGenerating.value) {
    return
  }
  if (params.value.prompt === '') {
    return ElMessage.error('请输入绘画提示词！')
  }

  if (!isLogin.value) {
    store.setShowLoginDialog(true)
    return
  }

  if (!params.value.size) {
    return ElMessage.error('请选择或填写图片尺寸！')
  }

  if (sizeLevel.value === 'other') {
    const val = params.value.size.trim()
    const presetLevels = ['1K', '2K', '4K']
    if (!val) {
      return ElMessage.error('请输入自定义图片尺寸！')
    }
    if (!presetLevels.includes(val)) {
      const sizePattern = /^\d+x\d+$/i
      if (!sizePattern.test(val)) {
        return ElMessage.error('自定义尺寸格式错误，请输入如 1024x1024')
      }
    }
  }
  isGenerating.value = true
  httpPost('/api/image/image', params.value)
    .then(() => {
      ElMessage.success('任务执行成功！')
      userPower.value -= imagePower.value
      // 追加任务列表
      runningJobs.value.push({
        prompt: params.value.prompt,
        progress: 0,
      })
      allowPulling.value = true
      isOver.value = false
    })
    .catch((e) => {
      ElMessage.error('任务执行失败：' + e.message)
    })
    .finally(() => {
      isGenerating.value = false
    })
}

const removeImage = (item) => {
  ElMessageBox.confirm('此操作将会删除任务和图片，继续操作码?', '删除提示', {
    confirmButtonText: '确认',
    cancelButtonText: '取消',
    type: 'warning',
  })
    .then(() => {
      httpGet('/api/image/remove', { id: item.id })
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

const previewImg = (item) => {
  previewURL.value = item.img_url
}

// 发布图片到作品墙
const publishImage = (item, action) => {
  let text = '图片发布'
  if (action === false) {
    text = '取消发布'
  }
  httpGet('/api/image/publish', { id: item.id, action: action })
    .then(() => {
      ElMessage.success(text + '成功')
      item.publish = action
      page.value = 0
      isOver.value = false
    })
    .catch((e) => {
      ElMessage.error(text + '失败：' + e.message)
    })
}

const promptGenerating = ref(false)
const generatePrompt = () => {
  if (params.value.prompt === '') {
    return showMessageError('请输入原始提示词')
  }
  promptGenerating.value = true
  httpPost('/api/prompt/image', { prompt: params.value.prompt })
    .then((res) => {
      params.value.prompt = res.data
      promptGenerating.value = false
    })
    .catch((e) => {
      showMessageError('生成提示词失败：' + e.message)
      promptGenerating.value = false
    })
}

const changeModel = (model) => {
  imagePower.value = model.power
  params.value.model_id = selectedModel.value.id
}

// 下载图片
const downloadImage = (item) => {
  downloadFile(item, 'img_url')
}

// 显示详情
const showDetail = (item) => {
  // 解析 params 字段
  let params = {}
  try {
    if (item.params) {
      if (typeof item.params === 'string') {
        params = JSON.parse(item.params)
      } else {
        params = item.params
      }
    }
  } catch (e) {
    console.error('解析 params 失败:', e)
  }

  currentDetail.value = {
    ...item,
    params: params,
  }
  detailDialogVisible.value = true
}

</script>

<style lang="scss" scoped>
@use '../assets/css/image.scss' as *;
@use '../assets/css/custom-scroll.scss' as *;

.detail-content {
  :deep(.el-descriptions__label) {
    min-width: 110px;
  }

  .prompt-text {
    display: flex;
    align-items: center;
    word-break: break-all;
  }

  .reference-images {
    display: flex;
    flex-wrap: wrap;
  }
}
</style>

<style lang="scss"></style>
