<template>
  <div class="page-threed">
    <!-- 左侧参数设置面板 -->
    <div class="params-panel">
      <!-- 平台选择Tab -->
      <div class="platform-tabs">
        <CustomTabs v-model="activePlatform" @change="handlePlatformChange">
          <CustomTabPane label="魔力方舟" name="gitee">
            <div class="platform-info">
              <i class="iconfont icon-gitee"></i>
              <span>Gitee AI 3D生成</span>
            </div>
          </CustomTabPane>
          <CustomTabPane label="腾讯混元" name="tencent">
            <div class="platform-info">
              <i class="iconfont icon-tencent"></i>
              <span>腾讯云混元3D生成</span>
            </div>
          </CustomTabPane>
        </CustomTabs>
      </div>

      <!-- 参数容器 -->
      <div class="params-container">
        <!-- 图片上传区域 -->
        <div class="param-line pt">
          <span class="label">上传图片：</span>
        </div>
        <div class="param-line">
          <ImageUpload
            v-model="currentImage"
            :max-count="1"
            :multiple="false"
            @change="handleImageChange"
          />
        </div>

        <!-- 文本提示词 -->
        <div class="param-line pt">
          <span class="label">提示词：</span>
        </div>
        <div class="param-line">
          <el-input
            v-model="currentPrompt"
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
          <el-select v-model="selectedModel" placeholder="选择输出格式" @change="handleModelChange">
            <el-option
              v-for="(model, key) in availableModels"
              :key="key"
              :label="model.name"
              :value="key"
            />
          </el-select>
        </div>

        <!-- 算力消耗显示 -->
        <div class="param-line pt">
          <span class="label">算力消耗：</span>
        </div>
        <div class="power-display">
          <span class="power-value">{{ currentPower }}</span>
          <span class="power-unit">点</span>
        </div>

        <!-- 生成按钮 -->
        <div class="generate-section">
          <el-button
            type="primary"
            size="large"
            :loading="generating"
            :disabled="!canGenerate"
            @click="generate3D"
            class="generate-btn"
          >
            {{ generating ? '生成中...' : '开始生成' }}
          </el-button>
        </div>
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
                {{ task.params ? getPromptFromParams(task.params) : '' }}
              </div>
              <div class="task-progress" v-if="task.status === 'processing'">
                <el-progress :percentage="task.progress" :stroke-width="4" />
              </div>
            </div>

            <div class="task-actions" v-if="task.status === 'completed'">
              <el-button size="small" @click="preview3D(task)">预览</el-button>
              <el-button size="small" type="primary" @click="download3D(task)">下载</el-button>
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
            @size-change="handleSizeChange"
            @current-change="handleCurrentChange"
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

// 响应式数据
const activePlatform = ref('gitee')
const currentImage = ref([])
const currentPrompt = ref('')
const selectedModel = ref('obj')
const generating = ref(false)
const previewVisible = ref(false)
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)
const taskList = ref([])
const currentPreviewTask = ref(null)

// 平台配置
const platformConfig = {
  gitee: {
    name: '魔力方舟',
    models: {
      obj: { name: 'OBJ格式', power: 45 },
      glb: { name: 'GLB格式', power: 55 },
      stl: { name: 'STL格式', power: 35 },
      usdz: { name: 'USDZ格式', power: 65 },
      fbx: { name: 'FBX格式', power: 75 },
      mp4: { name: 'MP4格式', power: 85 },
    },
  },
  tencent: {
    name: '腾讯混元',
    models: {
      obj: { name: 'OBJ格式', power: 50 },
      glb: { name: 'GLB格式', power: 60 },
      stl: { name: 'STL格式', power: 40 },
      usdz: { name: 'USDZ格式', power: 70 },
      fbx: { name: 'FBX格式', power: 80 },
      mp4: { name: 'MP4格式', power: 90 },
    },
  },
}

// 计算属性
const availableModels = computed(() => {
  return platformConfig[activePlatform.value]?.models || {}
})

const currentPower = computed(() => {
  return availableModels.value[selectedModel.value]?.power || 0
})

const canGenerate = computed(() => {
  return currentPrompt.value.trim() && currentImage.value.length > 0 && selectedModel.value
})

// 方法
const handlePlatformChange = (platform) => {
  // 切换平台时重置模型选择
  if (!availableModels.value[selectedModel.value]) {
    selectedModel.value = Object.keys(availableModels.value)[0]
  }
}

const handleImageChange = (files) => {
  currentImage.value = files
}

const handleModelChange = () => {
  // 模型改变时的处理逻辑
}

const generate3D = async () => {
  if (!canGenerate.value) {
    ElMessage.warning('请完善生成参数')
    return
  }

  try {
    generating.value = true

    const requestData = {
      type: activePlatform.value,
      model: selectedModel.value,
      prompt: currentPrompt.value,
      image_url: currentImage.value[0]?.url || '',
      power: currentPower.value,
    }

    const response = await httpPost('/api/3d/generate', requestData)

    if (response.code === 0) {
      ElMessage.success('任务创建成功')
      // 清空表单
      currentImage.value = []
      currentPrompt.value = ''
      // 刷新任务列表
      loadTasks()
    } else {
      ElMessage.error(response.message || '创建任务失败')
    }
  } catch (error) {
    ElMessage.error('创建任务失败：' + error.message)
  } finally {
    generating.value = false
  }
}

const loadTasks = async () => {
  try {
    const response = await httpGet('/api/3d/jobs', {
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

const handleSizeChange = (size) => {
  pageSize.value = size
  currentPage.value = 1
  loadTasks()
}

const handleCurrentChange = (page) => {
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

    const response = await httpGet(`/api/3d/job/${taskId}/delete`)
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

const download3D = async (task) => {
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
    download3D(currentPreviewTask.value)
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

const getPromptFromParams = (paramsStr) => {
  try {
    const params = JSON.parse(paramsStr)
    return params.prompt || ''
  } catch {
    return ''
  }
}

// 生命周期
onMounted(() => {
  loadTasks()
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
