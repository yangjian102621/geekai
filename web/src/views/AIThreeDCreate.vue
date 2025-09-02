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
                <span class="label">上传图片：</span>
              </div>
              <div class="param-line">
                <ImageUpload v-model="giteeForm.image_url" :max-count="1" :multiple="false" />
              </div>

              <!-- 文本提示词 -->
              <div class="param-line pt">
                <span class="label">提示词：</span>
              </div>
              <div class="param-line">
                <el-input
                  v-model="giteeForm.prompt"
                  type="textarea"
                  :autosize="{ minRows: 3, maxRows: 5 }"
                  placeholder="请输入3D模型描述，越详细越好"
                  maxlength="2000"
                  show-word-limit
                />
              </div>

              <!-- 模型选择 -->
              <div class="param-line pt">
                <span class="label">输出格式：</span>
              </div>
              <div class="param-line">
                <el-select
                  v-model="giteeForm.model"
                  placeholder="选择输出格式"
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

              <!-- 高级参数 -->
              <div class="param-line pt">
                <el-button
                  @click="giteeAdvancedVisible = !giteeAdvancedVisible"
                  class="advanced-toggle-btn"
                >
                  <i
                    :class="
                      giteeAdvancedVisible
                        ? 'iconfont icon-arrow-down'
                        : 'iconfont icon-arrow-right'
                    "
                  ></i>
                  <span>高级参数设置</span>
                </el-button>
              </div>

              <!-- 高级参数内容 -->
              <div v-show="giteeAdvancedVisible" class="advanced-params">
                <!-- 纹理开关 -->
                <div class="param-line">
                  <el-checkbox v-model="giteeForm.texture">启用纹理</el-checkbox>
                </div>

                <!-- 随机种子 -->
                <div class="param-line">
                  <span class="label">随机种子：</span>
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
                  <span class="label">迭代次数：</span>
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
                  <span class="label">引导系数：</span>
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
                  <span class="label">3D渲染精度：</span>
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
              <!-- 图片上传区域 -->
              <div class="param-line pt">
                <span class="label">上传图片：</span>
              </div>
              <div class="param-line">
                <ImageUpload v-model="tencentForm.image_url" :max-count="1" :multiple="false" />
              </div>

              <!-- 文本提示词 -->
              <div class="param-line pt">
                <span class="label">提示词：</span>
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

              <!-- 模型选择 -->
              <div class="param-line pt">
                <span class="label">输出格式：</span>
              </div>
              <div class="param-line">
                <el-select
                  v-model="tencentForm.model"
                  @change="handleModelChange"
                  placeholder="选择输出格式"
                >
                  <el-option
                    v-for="model in configs.tencent.models"
                    :key="model.name"
                    :label="model.name"
                    :value="model.name"
                  />
                </el-select>
              </div>

              <!-- 高级参数 -->
              <div class="param-line pt">
                <span class="label">高级参数：</span>
              </div>

              <!-- PBR材质开关 -->
              <div class="param-line">
                <el-checkbox v-model="tencentForm.enable_pbr">启用PBR材质</el-checkbox>
              </div>

              <!-- 文件格式选择 -->
              <div class="param-line">
                <span class="label">文件格式：</span>
                <el-select v-model="tencentForm.file_format" style="width: 100%">
                  <el-option label="GLB" value="glb" />
                  <el-option label="GLTF" value="gltf" />
                  <el-option label="OBJ" value="obj" />
                  <el-option label="FBX" value="fbx" />
                </el-select>
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
            class="task-item"
            :class="{ completed: task.status === 'completed' }"
          >
            <div class="task-header">
              <span class="task-id">#{{ task.id }}</span>
              <span class="task-status" :class="task.status">
                {{ getStatusText(task.status) }}
              </span>
            </div>

            <div class="task-content">
              <div class="task-prompt">
                {{ task.params?.prompt }}
              </div>
              <div class="task-progress" v-if="task.status === 'processing'">
                <el-progress :percentage="task.progress" :stroke-width="4" />
              </div>
            </div>

            <div class="task-actions" v-if="task.status === 'completed'">
              <el-button size="small" @click="preview3D(task)">预览</el-button>
              <el-button size="small" type="primary" @click="download(task)">下载</el-button>
            </div>

            <div class="task-actions" v-else>
              <el-button size="small" @click="deleteTask(task.id)">删除</el-button>
            </div>
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
          v-if="currentPreviewTask && currentPreviewTask.img_url"
          :model-url="currentPreviewTask.img_url"
          :model-type="currentPreviewTask.model"
        />
        <div v-else class="preview-placeholder">
          <i class="iconfont icon-3d"></i>
          <p>暂无3D模型</p>
        </div>
      </div>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="closePreview">关闭</el-button>
          <el-button type="primary" @click="downloadCurrentModel">下载模型</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import ImageUpload from '@/components/ImageUpload.vue'
import ThreeDPreview from '@/components/ThreeDPreview.vue'
import CustomTabPane from '@/components/ui/CustomTabPane.vue'
import CustomTabs from '@/components/ui/CustomTabs.vue'
import { httpGet, httpPost } from '@/utils/http'
import { ElMessage, ElMessageBox } from 'element-plus'
import { computed, onMounted, ref } from 'vue'
import { checkSession } from '@/store/cache'

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
const tencentForm = ref({
  prompt: '',
  image_url: '',
  model: '',
  power: 0,
  file_format: '', // 输出文件格式
  enable_pbr: false, // 是否开启PBR材质
})
const giteeForm = ref({
  prompt: '',
  image_url: '',
  model: '',
  power: 0,
  file_format: '', // 输出文件格式
  texture: false, // 是否开启纹理
  seed: 1234, // 随机种子
  num_inference_steps: 5, //迭代次数
  guidance_scale: 7.5, //引导系数
  octree_resolution: 128, // 3D 渲染精度，越高3D 细节越丰富
})
const currentPower = ref(0)

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
  } else {
    const model = configs.value.gitee.models.find((model) => model.name === value)
    currentPower.value = model.power
    giteeForm.value.power = model.power
  }
}

const handlePlatformChange = (value) => {
  currentPower.value = value === 'tencent' ? tencentForm.value.power : giteeForm.value.power
}

const generate3D = async () => {
  if (currentPower.value === 0) {
    ElMessage.warning('请完善生成参数')
    return
  }

  if (!currentPrompt.value.trim()) {
    ElMessage.warning('请输入提示词')
    return
  }

  if (!selectedModel.value) {
    ElMessage.warning('请选择输出格式')
    return
  }

  try {
    loading.value = true

    const requestData = {
      type: activePlatform.value,
      model: selectedModel.value,
      prompt: currentPrompt.value,
      image_url: currentImage.value[0]?.url || '',
      power: currentPower.value,
      ...currentForm.value, // 包含所有表单参数
    }

    const response = await httpPost('/api/ai3d/generate', requestData)

    if (response.code === 0) {
      ElMessage.success('任务创建成功')
      // 清空表单
      tencentForm.value = {
        prompt: '',
        image_url: '',
        model: '',
        power: 0,
        file_format: '',
        enable_pbr: false,
      }
      giteeForm.value = {
        prompt: '',
        image_url: '',
        model: '',
        power: 0,
        file_format: '',
        texture: false,
        seed: 1234,
        num_inference_steps: 5,
        guidance_scale: 7.5,
        octree_resolution: 128,
      }
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
    const response = await httpGet('/api/ai3d/jobs', {
      page: currentPage.value,
      page_size: pageSize.value,
    })

    if (response.code === 0) {
      taskList.value = response.data.list
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

    const response = await httpGet(`/api/ai3d/job/${taskId}/delete`)
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
  currentPreviewTask.value = null
}

const download = async (task) => {
  if (!task.img_url) {
    ElMessage.warning('模型文件不存在')
    return
  }

  try {
    // 创建一个隐藏的a标签来下载文件
    const link = document.createElement('a')
    link.href = task.img_url
    link.download = `3d_model_${task.id}.${task.model}`
    link.style.display = 'none'
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)

    ElMessage.success('开始下载3D模型')
  } catch (error) {
    console.error('下载失败:', error)
    ElMessage.error('下载失败，请重试')
  }
}

const downloadCurrentModel = () => {
  if (currentPreviewTask.value) {
    download(currentPreviewTask.value)
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
.page-threed {
  display: flex;
  height: 100vh;
  background: #f5f5f5;
}

.params-panel {
  width: 400px;
  background: white;
  border-right: 1px solid #e4e7ed;
  padding: 20px;
  overflow-y: auto;
}

.platform-tabs {
  margin-bottom: 20px;
}

.platform-info {
  display: flex;
  align-items: center;
  gap: 8px;

  i {
    font-size: 18px;
    color: #409eff;
  }
}

.params-container {
  .param-line {
    margin-bottom: 16px;

    &.pt {
      margin-top: 20px;
    }

    .label {
      display: block;
      margin-bottom: 8px;
      font-weight: 500;
      color: #333;
    }
  }

  .advanced-toggle-btn {
    padding: 0;
    font-size: 14px;
    color: #409eff;
    border: none;
    background: none;
    display: flex;
    align-items: center;
    gap: 4px;
    transition: all 0.3s ease;

    &:hover {
      color: #66b1ff;
      background: #f0f9ff;
      border-radius: 4px;
      padding: 4px 8px;
    }

    i {
      font-size: 12px;
      transition: transform 0.3s ease;
    }
  }

  .advanced-params {
    margin-left: 16px;
    padding: 10px 16px;
    border-left: 3px solid #e4e7ed;
    margin-top: 8px;
  }
}

.power-display {
  display: flex;
  align-items: center;
  gap: 8px;

  .power-value {
    font-size: 24px;
    font-weight: bold;
    color: #409eff;
  }

  .power-unit {
    color: #666;
  }
}

.generate-section {
  margin-top: 30px;

  .generate-btn {
    width: 100%;
    height: 44px;
    font-size: 16px;
  }
}

.content-panel {
  flex: 1;
  padding: 20px;
  overflow-y: auto;
}

.task-list {
  .list-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;

    h3 {
      margin: 0;
      color: #333;
    }
  }
}

.task-items {
  .task-item {
    background: white;
    border-radius: 8px;
    padding: 16px;
    margin-bottom: 16px;
    border: 1px solid #e4e7ed;

    &.completed {
      border-color: #67c23a;
      background: #f0f9ff;
    }

    .task-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 12px;

      .task-id {
        font-weight: 500;
        color: #666;
      }

      .task-status {
        padding: 4px 8px;
        border-radius: 4px;
        font-size: 12px;

        &.pending {
          background: #fdf6ec;
          color: #e6a23c;
        }

        &.processing {
          background: #ecf5ff;
          color: #409eff;
        }

        &.completed {
          background: #f0f9ff;
          color: #67c23a;
        }

        &.failed {
          background: #fef0f0;
          color: #f56c6c;
        }
      }
    }

    .task-content {
      margin-bottom: 12px;

      .task-prompt {
        color: #666;
        margin-bottom: 8px;
        line-height: 1.4;
      }
    }

    .task-actions {
      display: flex;
      gap: 8px;
    }
  }
}

.pagination {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}

.preview-container {
  .three-container {
    width: 100%;
    height: 500px;
    background: #f0f0f0;
    border-radius: 8px;
    display: flex;
    align-items: center;
    justify-content: center;
    color: #666;
  }

  .preview-placeholder {
    width: 100%;
    height: 500px;
    background: #f0f0f0;
    border-radius: 8px;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    color: #666;

    i {
      font-size: 48px;
      margin-bottom: 12px;
      color: #999;
    }

    p {
      margin: 0;
      font-size: 14px;
    }
  }
}

// 响应式设计
@media (max-width: 768px) {
  .page-threed {
    flex-direction: column;
  }

  .params-panel {
    width: 100%;
    border-right: none;
    border-bottom: 1px solid #e4e7ed;
  }
}
</style>
