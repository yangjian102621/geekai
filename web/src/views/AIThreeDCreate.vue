<template>
  <div class="page-threed">
    <!-- 左侧参数设置面板 -->
    <div class="params-panel">
      <!-- 平台选择Tab -->
      <div class="platform-tabs">
        <CustomTabs v-model="activePlatform" @tab-click="handlePlatformChange">
          <CustomTabPane name="gitee" width="48%">
            <template #label>
              <div class="flex items-center justify-center">
                <i class="iconfont icon-gitee mr-1"></i>
                <span>Gitee 模力方舟</span>
              </div>
            </template>
            <!-- 参数容器 -->
            <div class="params-container">
              <!-- 图片上传区域 -->
              <div class="param-line pt">
                <span class="label"><span class="text-red-500 mr-1">*</span>上传图片：</span>
              </div>
              <div class="param-line">
                <ImageUpload v-model="giteeForm.image_url" :max-count="1" :multiple="false" />
              </div>

              <!-- 模型选择 -->
              <div class="param-line pt">
                <span class="label"><span class="text-red-500 mr-1">*</span>模型选择：</span>
              </div>
              <div class="param-line">
                <el-select
                  v-model="giteeForm.model"
                  placeholder="选择模型"
                  @change="handleModelChange"
                >
                  <el-option
                    v-for="model in configs.gitee.models"
                    :key="model.name"
                    :label="model.name"
                    :value="model.name"
                  />
                </el-select>
              </div>
              <div class="param-line">
                <el-alert v-if="giteeForm.model_desc" type="info" :closable="false">
                  {{ giteeForm.model_desc }}
                </el-alert>
              </div>

              <!-- 文件格式选择 -->
              <div class="param-line">
                <span class="label mb-3"><span class="text-red-500 mr-1">*</span>输出格式：</span>
                <el-select v-model="giteeForm.file_format" style="width: 100%">
                  <el-option
                    v-for="format in giteeSupportedFormats"
                    :key="format"
                    :label="format"
                    :value="format"
                  />
                </el-select>
              </div>

              <!-- 纹理开关 -->
              <div class="flex justify-between param-line">
                <span class="label">生成纹理：</span>
                <el-switch v-model="giteeForm.texture" size="large" />
              </div>

              <!-- 高级参数 -->
              <div class="param-line pt">
                <el-button
                  @click="giteeAdvancedVisible = !giteeAdvancedVisible"
                  class="advanced-toggle-btn"
                >
                  <i
                    :class="
                      giteeAdvancedVisible ? 'iconfont icon-arrow-up' : 'iconfont icon-arrow-down'
                    "
                  ></i>
                  <span>高级参数设置</span>
                </el-button>
              </div>

              <!-- 高级参数内容 -->
              <div v-show="giteeAdvancedVisible" class="advanced-params">
                <!-- 随机种子 -->
                <div class="param-line">
                  <span class="label mb-3">随机种子：</span>
                  <el-input-number
                    v-model="giteeForm.seed"
                    :min="1"
                    :max="999999"
                    controls-position="right"
                    style="width: 100%"
                  />
                </div>

                <!-- 迭代次数 -->
                <div class="param-line">
                  <span class="label mb-3">迭代次数：</span>
                  <el-input-number
                    v-model="giteeForm.num_inference_steps"
                    :min="1"
                    :max="50"
                    controls-position="right"
                    style="width: 100%"
                  />
                </div>

                <!-- 引导系数 -->
                <div class="param-line">
                  <span class="label mb-3">引导系数：</span>
                  <el-input-number
                    v-model="giteeForm.guidance_scale"
                    :min="1"
                    :max="20"
                    :step="0.5"
                    controls-position="right"
                    style="width: 100%"
                  />
                </div>

                <!-- 3D渲染精度 -->
                <div class="param-line">
                  <span class="label mb-3">3D渲染精度：</span>
                  <el-select v-model="giteeForm.octree_resolution" style="width: 100%">
                    <el-option label="64 (低精度)" :value="64" />
                    <el-option label="128 (中精度)" :value="128" />
                    <el-option label="256 (高精度)" :value="256" />
                  </el-select>
                </div>
              </div>
            </div>
          </CustomTabPane>
          <CustomTabPane name="tencent" width="48%">
            <template #label>
              <div class="flex items-center justify-center">
                <i class="iconfont icon-tencent mr-1"></i>
                <span>腾讯云混元3D</span>
              </div>
            </template>
            <!-- 参数容器 -->
            <div class="params-container">
              <div class="param-line pt flex justify-between items-center">
                <span class="label">生成模式：</span>
                <custom-switch
                  v-model="tencentForm.text3d"
                  active-color="#9c27b0"
                  inactive-color="#409eff"
                  :width="120"
                  size="large"
                >
                  <template #active-text>
                    <div class="flex items-center justify-start pl-4 text-sm">
                      <i class="iconfont icon-image mr-1"></i> <span>文生3D</span>
                    </div>
                  </template>
                  <template #inactive-text>
                    <div class="flex items-center justify-end pl-4 text-sm">
                      <i class="iconfont icon-doc mr-1"></i> <span>图生3D</span>
                    </div>
                  </template>
                </custom-switch>
              </div>

              <!-- 文本提示词 -->
              <div v-if="tencentForm.text3d">
                <div class="param-line pt">
                  <span class="label"><span class="text-red-500 mr-1">*</span>提示词：</span>
                </div>
                <div class="param-line">
                  <el-input
                    v-model="tencentForm.prompt"
                    type="textarea"
                    :autosize="{ minRows: 3, maxRows: 5 }"
                    placeholder="请输入3D模型描述，越详细越好"
                    maxlength="2000"
                    show-word-limit
                  />
                </div>
              </div>
              <div v-else>
                <!-- 图片上传区域 -->
                <div class="param-line pt">
                  <span class="label"><span class="text-red-500 mr-1">*</span>上传图片：</span>
                </div>
                <div class="param-line">
                  <ImageUpload v-model="tencentForm.image_url" :max-count="1" :multiple="false" />
                </div>
              </div>

              <!-- 模型选择 -->
              <div class="param-line pt">
                <span class="label mb-2"><span class="text-red-500 mr-1">*</span>模型选择：</span>
              </div>
              <div class="param-line">
                <el-select
                  v-model="tencentForm.model"
                  @change="handleModelChange"
                  placeholder="选择模型"
                >
                  <el-option
                    v-for="model in configs.tencent.models"
                    :key="model.name"
                    :label="model.name"
                    :value="model.name"
                  />
                </el-select>
              </div>
              <div class="param-line">
                <el-alert v-if="tencentForm.model_desc" type="info" :closable="false">
                  {{ tencentForm.model_desc }}
                </el-alert>
              </div>

              <!-- 文件格式选择 -->
              <div class="param-line">
                <span class="label mb-3"><span class="text-red-500 mr-1">*</span>输出格式：</span>
                <el-select v-model="tencentForm.file_format" style="width: 100%">
                  <el-option
                    v-for="format in tencentSupportedFormats"
                    :key="format"
                    :label="format"
                    :value="format"
                  />
                </el-select>
              </div>

              <!-- PBR材质开关 -->
              <div class="flex justify-between param-line">
                <span class="label">启用PBR材质：</span>
                <el-switch v-model="tencentForm.enable_pbr" size="large" />
              </div>
            </div>
          </CustomTabPane>
          <!-- 生成按钮 -->
          <div class="generate-section">
            <button
              @click="generate3D"
              :disabled="loading"
              type="button"
              class="w-full py-3 bg-gradient-to-r from-blue-500 to-purple-600 text-white rounded-xl disabled:from-gray-400 disabled:to-gray-400 disabled:cursor-not-allowed hover:from-blue-600 hover:to-purple-700 transition-all duration-200 flex items-center justify-center space-x-2"
            >
              <i v-if="loading" class="iconfont icon-loading animate-spin"></i>
              <i v-else class="iconfont icon-chuangzuo"></i>
              <span>{{ loading ? '创作中...' : `立即生成 (${currentPower}算力)` }}</span>
            </button>
          </div>
        </CustomTabs>
      </div>
    </div>

    <!-- 右侧内容区域 -->
    <div class="content-panel">
      <!-- 任务列表 -->
      <div class="task-list">
        <div class="list-header">
          <h3>生成任务</h3>
          <el-button size="small" @click="refreshTasks">刷新</el-button>
        </div>

        <div class="task-items">
          <div
            v-for="task in taskList"
            :key="task.id"
            class="task-card"
            :class="getTaskCardClass(task.status)"
          >
            <!-- 任务卡片头部 -->
            <div class="task-card-header">
              <div class="task-info">
                <div class="task-id">
                  <i class="iconfont icon-renwu mr-2"></i>
                  #{{ task.id }}
                </div>
                <div class="task-platform">
                  <i :class="getPlatformIcon(task.type)" class="mr-1"></i>
                  {{ getPlatformName(task.type) }}
                </div>
              </div>
              <div class="task-status-wrapper">
                <div class="task-status" :class="task.status">
                  <i :class="getStatusIcon(task.status)" class="mr-1"></i>
                  {{ getStatusText(task.status) }}
                </div>
                <div class="task-power">
                  <i class="iconfont icon-suanli mr-1"></i>
                  {{ task.power }}
                </div>
              </div>
            </div>

            <!-- 任务卡片内容 -->
            <div class="task-card-content">
              <!-- 左侧预览图 -->
              <div class="task-preview">
                <div v-if="task.status === 'completed' && task.preview_url" class="preview-image">
                  <img :src="task.preview_url" :alt="getTaskPrompt(task)" />
                  <div class="preview-overlay">
                    <i class="iconfont icon-yulan"></i>
                  </div>
                </div>
                <div v-else-if="getTaskImageUrl(task)" class="input-image">
                  <img :src="getTaskImageUrl(task)" :alt="getTaskPrompt(task)" />
                  <div class="input-overlay">
                    <i class="iconfont icon-tupian"></i>
                  </div>
                </div>
                <div v-else class="prompt-placeholder">
                  <i class="iconfont icon-wenzi"></i>
                  <span>{{ getTaskPrompt(task) }}</span>
                </div>
              </div>

              <!-- 右侧任务详情 -->
              <div class="task-details">
                <div class="task-model">
                  <i class="iconfont icon-moxing mr-1"></i>
                  {{ task.model }}
                </div>

                <div class="task-prompt" v-if="getTaskPrompt(task)">
                  <i class="iconfont icon-tishi mr-1"></i>
                  <span>{{ getTaskPrompt(task) }}</span>
                </div>

                <div class="task-params" v-if="getTaskParams(task)">
                  <i class="iconfont icon-shezhi mr-1"></i>
                  <span>{{ getTaskParams(task) }}</span>
                </div>

                <div class="task-time">
                  <i class="iconfont icon-shijian mr-1"></i>
                  {{ formatTime(task.created_at) }}
                </div>

                <div class="task-error" v-if="task.status === 'failed' && task.err_msg">
                  <i class="iconfont icon-cuowu mr-1"></i>
                  <span>{{ task.err_msg }}</span>
                </div>
              </div>
            </div>

            <!-- 任务卡片底部操作 -->
            <div class="task-card-footer">
              <div class="task-actions">
                <el-button
                  v-if="task.status === 'completed'"
                  size="small"
                  type="primary"
                  @click="preview3D(task)"
                  class="action-btn preview-btn"
                >
                  <i class="iconfont icon-eye-open mr-1"></i>
                  预览
                </el-button>

                <el-button
                  v-if="task.status === 'completed'"
                  size="small"
                  type="success"
                  @click="downloadFile(task)"
                  :loading="task.downloading"
                  class="action-btn download-btn"
                >
                  <i class="iconfont icon-download mr-1" v-if="!task.downloading"></i>
                  <span v-if="task.downloading">下载中...</span>
                  <span v-else>下载</span>
                </el-button>

                <el-button
                  size="small"
                  type="danger"
                  @click="deleteTask(task.id)"
                  class="action-btn delete-btn"
                >
                  <i class="iconfont icon-remove mr-1"></i>
                  删除
                </el-button>

                <el-button
                  v-if="task.status === 'processing'"
                  size="small"
                  type="info"
                  disabled
                  class="action-btn processing-btn"
                >
                  <i class="iconfont icon-loading animate-spin mr-1"></i>
                  处理中...
                </el-button>
              </div>
            </div>
          </div>

          <!-- 空状态 -->
          <div v-if="taskList.length === 0" class="empty-state">
            <i class="iconfont icon-kong"></i>
            <p>暂无任务，开始创建你的第一个3D模型吧！</p>
          </div>
        </div>

        <!-- 分页 -->
        <div class="pagination" v-if="total > 0">
          <el-pagination
            :current-page="currentPage"
            :page-size="pageSize"
            :page-sizes="[10, 20, 50]"
            :total="total"
            layout="total, sizes, prev, pager, next, jumper"
            @size-change="handlePageSizeChange"
            @current-change="handleCurrentPageChange"
          />
        </div>
      </div>
    </div>

    <!-- 3D预览弹窗 -->
    <el-dialog v-model="previewVisible" title="3D模型预览" width="80%" :before-close="closePreview">
      <div class="preview-container">
        <ThreeDPreview
          v-if="currentPreviewTask && currentPreviewTask.file_url"
          :model-url="currentPreviewTask.file_url"
        />
        <div v-else class="preview-placeholder">
          <i class="iconfont icon-3d"></i>
          <p>暂无3D模型</p>
        </div>
      </div>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="closePreview">关闭</el-button>
          <el-button
            type="primary"
            @click="downloadCurrentModel"
            :loading="currentPreviewTask.downloading"
          >
            <span v-if="!currentPreviewTask.downloading">下载模型</span>
            <span v-else>下载中...</span>
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import ImageUpload from '@/components/ImageUpload.vue'
import ThreeDPreview from '@/components/ThreeDPreview.vue'
import CustomSwitch from '@/components/ui/CustomSwitch.vue'
import CustomTabPane from '@/components/ui/CustomTabPane.vue'
import CustomTabs from '@/components/ui/CustomTabs.vue'
import { checkSession } from '@/store/cache'
import { httpDownload, httpGet, httpPost } from '@/utils/http'
import { ElMessage, ElMessageBox } from 'element-plus'
import { computed, onMounted, ref } from 'vue'
import { showMessageError } from '../utils/dialog'
import { replaceImg } from '../utils/libs'

// 响应式数据
const activePlatform = ref('gitee')
const loading = ref(false)
const previewVisible = ref(false)
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)
const taskList = ref([])
const currentPreviewTask = ref(null)
const giteeAdvancedVisible = ref(false) // 控制Gitee高级参数显示状态
const tencentDefaultForm = {
  text3d: false,
  prompt: '',
  image_url: '',
  model: '',
  file_format: '', // 输出文件格式
  enable_pbr: false, // 是否开启PBR材质
  model_desc: '', // 模型描述
  power: 0, // 算力消耗
}
const giteeDefaultForm = {
  prompt: '',
  image_url: '',
  model: '',
  file_format: '', // 输出文件格式
  texture: false, // 是否开启纹理
  seed: 1234, // 随机种子
  num_inference_steps: 5, //迭代次数
  guidance_scale: 7.5, //引导系数
  octree_resolution: 128, // 3D 渲染精度，越高3D 细节越丰富
  model_desc: '', // 模型描述
  power: 0, // 算力消耗
}
const tencentForm = ref(tencentDefaultForm)
const giteeForm = ref(giteeDefaultForm)
const currentPower = ref(0)
const tencentSupportedFormats = ref([])
const giteeSupportedFormats = ref([])

// 计算属性：获取当前活跃平台的表单数据
const currentForm = computed(() => {
  return activePlatform.value === 'tencent' ? tencentForm.value : giteeForm.value
})

const selectedModel = computed(() => {
  return currentForm.value.model
})

const currentPrompt = computed(() => {
  return currentForm.value.prompt
})

const currentImage = computed(() => {
  return currentForm.value.image_url ? [{ url: currentForm.value.image_url }] : []
})

const configs = ref({
  gitee: { models: [] },
  tencent: { models: [] },
})

const loadConfigs = async () => {
  const response = await httpGet('/api/ai3d/configs')
  configs.value = response.data
}

const handleModelChange = (value) => {
  if (activePlatform.value === 'tencent') {
    const model = configs.value.tencent.models.find((model) => model.name === value)
    currentPower.value = model.power
    tencentForm.value.power = model.power
    tencentForm.value.model_desc = model.desc
    tencentForm.value.file_format = model.formats[0]
    tencentSupportedFormats.value = model.formats
  } else {
    const model = configs.value.gitee.models.find((model) => model.name === value)
    currentPower.value = model.power
    giteeForm.value.power = model.power
    giteeForm.value.model_desc = model.desc
    giteeForm.value.file_format = model.formats[0]
    giteeSupportedFormats.value = model.formats
  }
}

const handlePlatformChange = (value) => {
  currentPower.value = value === 'tencent' ? tencentForm.value.power : giteeForm.value.power
}

const generate3D = async () => {
  const requestData = {
    ...(activePlatform.value === 'tencent' ? tencentForm.value : giteeForm.value),
  }
  if (requestData.model === '') {
    ElMessage.warning('请选择模型')
    return
  }
  if (requestData.file_format === '') {
    ElMessage.warning('请选择输出格式')
    return
  }

  try {
    loading.value = true

    requestData.type = activePlatform.value
    if (requestData.image_url !== '') {
      requestData.image_url = replaceImg(requestData.image_url[0].url)
    }

    const response = await httpPost('/api/ai3d/generate', requestData)

    if (response.code === 0) {
      ElMessage.success('任务创建成功')
      // 清空表单
      tencentForm.value = tencentDefaultForm
      giteeForm.value = giteeDefaultForm
      currentPower.value = 0
      // 刷新任务列表
      loadTasks()
    } else {
      ElMessage.error(response.message || '创建任务失败')
    }
  } catch (error) {
    ElMessage.error('创建任务失败：' + error.message)
  } finally {
    loading.value = false
  }
}

const loadTasks = async () => {
  try {
    const response = await httpGet('/api/ai3d/jobs/mock', {
      page: currentPage.value,
      page_size: pageSize.value,
    })

    if (response.code === 0) {
      taskList.value = response.data.items
      total.value = response.data.total
    }
  } catch (error) {
    ElMessage.error('加载任务列表失败：' + error.message)
  }
}

const refreshTasks = () => {
  loadTasks()
}

const handlePageSizeChange = (size) => {
  pageSize.value = size
  currentPage.value = 1
  loadTasks()
}

const handleCurrentPageChange = (page) => {
  currentPage.value = page
  loadTasks()
}

const deleteTask = async (taskId) => {
  try {
    await ElMessageBox.confirm('确定要删除这个任务吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    })

    const response = await httpGet(`/api/ai3d/job/delete?id=${taskId}`)
    if (response.code === 0) {
      ElMessage.success('删除成功')
      loadTasks()
    } else {
      ElMessage.error(response.message || '删除失败')
    }
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败：' + error.message)
    }
  }
}

const preview3D = (task) => {
  currentPreviewTask.value = task
  previewVisible.value = true
}

const closePreview = () => {
  previewVisible.value = false
}

const downloadFile = async (item) => {
  const url = replaceImg(item.file_url)
  const downloadURL = `/api/download?url=${url}`
  const urlObj = new URL(url)
  const fileName = urlObj.pathname.split('/').pop()

  item.downloading = true

  try {
    const response = await httpDownload(downloadURL)
    const blob = new Blob([response.data])
    const link = document.createElement('a')
    link.href = URL.createObjectURL(blob)
    link.download = fileName
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    URL.revokeObjectURL(link.href)
    item.downloading = false
  } catch (error) {
    showMessageError('下载失败')
    item.downloading = false
  }
}

const downloadCurrentModel = () => {
  if (currentPreviewTask.value) {
    downloadFile(currentPreviewTask.value)
  }
}

const getStatusText = (status) => {
  const statusMap = {
    pending: '等待中',
    processing: '处理中',
    completed: '已完成',
    failed: '失败',
  }
  return statusMap[status] || status
}

const getTaskCardClass = (status) => {
  if (status === 'completed') {
    return 'task-card-completed'
  } else if (status === 'processing') {
    return 'task-card-processing'
  } else if (status === 'failed') {
    return 'task-card-failed'
  } else {
    return 'task-card-default'
  }
}

const getPlatformIcon = (type) => {
  if (type === 'gitee') {
    return 'iconfont icon-gitee'
  } else if (type === 'tencent') {
    return 'iconfont icon-tencent'
  }
  return 'iconfont icon-question'
}

const getPlatformName = (type) => {
  if (type === 'gitee') {
    return 'Gitee 模力方舟'
  } else if (type === 'tencent') {
    return '腾讯云混元3D'
  }
  return '未知平台'
}

const getStatusIcon = (status) => {
  if (status === 'pending') {
    return 'iconfont icon-pending'
  } else if (status === 'processing') {
    return 'iconfont icon-processing'
  } else if (status === 'completed') {
    return 'iconfont icon-completed'
  } else if (status === 'failed') {
    return 'iconfont icon-failed'
  }
  return 'iconfont icon-question'
}

const getTaskPrompt = (task) => {
  try {
    if (task.params) {
      const parsedParams = JSON.parse(task.params)
      return parsedParams.prompt || '文生3D任务'
    }
    return '文生3D任务'
  } catch (e) {
    return '文生3D任务'
  }
}

const getTaskImageUrl = (task) => {
  try {
    if (task.params) {
      const parsedParams = JSON.parse(task.params)
      return parsedParams.image_url || null
    }
    return null
  } catch (e) {
    return null
  }
}

const getTaskParams = (task) => {
  try {
    if (task.params) {
      const parsedParams = JSON.parse(task.params)
      const params = []

      if (parsedParams.texture) {
        params.push('纹理')
      }
      if (parsedParams.enable_pbr) {
        params.push('PBR材质')
      }
      if (parsedParams.num_inference_steps && parsedParams.num_inference_steps !== 5) {
        params.push(`迭代次数: ${parsedParams.num_inference_steps}`)
      }
      if (parsedParams.guidance_scale && parsedParams.guidance_scale !== 7.5) {
        params.push(`引导系数: ${parsedParams.guidance_scale}`)
      }
      if (parsedParams.octree_resolution && parsedParams.octree_resolution !== 128) {
        params.push(`精度: ${parsedParams.octree_resolution}`)
      }
      if (parsedParams.seed && parsedParams.seed !== 1234) {
        params.push(`种子: ${parsedParams.seed}`)
      }

      return params.join('，')
    }
    return ''
  } catch (e) {
    return ''
  }
}

const formatTime = (timestamp) => {
  const date = new Date(timestamp)
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  const hours = String(date.getHours()).padStart(2, '0')
  const minutes = String(date.getMinutes()).padStart(2, '0')
  return `${year}-${month}-${day} ${hours}:${minutes}`
}

// 生命周期
onMounted(() => {
  loadConfigs()
  checkSession()
    .then(() => {
      loadTasks()
    })
    .catch(() => {})
})
</script>

<style lang="scss" scoped>
@use '@/assets/css/ai3d.scss' as ai3d;
</style>
