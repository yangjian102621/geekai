<template>
  <div class="mobile-threed-create">
    <!-- 顶部导航 -->
    <div class="top-nav">
      <div class="nav-left" @click="$router.go(-1)">
        <i class="iconfont icon-arrow-left"></i>
      </div>
      <div class="nav-title">3D模型生成</div>
      <div class="nav-right"></div>
    </div>

    <!-- 平台选择 -->
    <div class="platform-selector">
      <div class="selector-tabs">
        <div
          v-for="platform in platforms"
          :key="platform.key"
          :class="['selector-tab', { active: activePlatform === platform.key }]"
          @click="activePlatform = platform.key"
        >
          <div class="tab-icon">
            <i :class="platform.icon"></i>
          </div>
          <div class="tab-name">{{ platform.name }}</div>
        </div>
      </div>
    </div>

    <!-- 参数设置 -->
    <div class="params-section">
      <!-- 图片上传 -->
      <div class="param-group">
        <div class="param-label">上传图片</div>
        <div class="image-upload-area">
          <ImageUpload
            v-model="currentImage"
            :max-count="1"
            :multiple="false"
            @change="handleImageChange"
          />
        </div>
      </div>

      <!-- 提示词输入 -->
      <div class="param-group">
        <div class="param-label">提示词描述</div>
        <div class="prompt-input">
          <el-input
            v-model="currentPrompt"
            type="textarea"
            :rows="4"
            placeholder="请输入3D模型描述，越详细越好"
            maxlength="2000"
            show-word-limit
          />
        </div>
      </div>

      <!-- 模型选择 -->
      <div class="param-group">
        <div class="param-label">输出格式</div>
        <div class="model-selector">
          <div
            v-for="(model, key) in availableModels"
            :key="key"
            :class="['model-option', { active: selectedModel === key }]"
            @click="selectedModel = key"
          >
            <div class="model-name">{{ model.name }}</div>
            <div class="model-power">{{ model.power }}点</div>
          </div>
        </div>
      </div>

      <!-- 算力消耗 -->
      <div class="power-info">
        <div class="power-label">算力消耗</div>
        <div class="power-value">{{ currentPower }} 点</div>
      </div>
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

    <!-- 任务列表 -->
    <div class="task-section">
      <div class="section-header">
        <h3>生成任务</h3>
        <el-button size="small" @click="refreshTasks">刷新</el-button>
      </div>

      <div class="task-list">
        <div
          v-for="task in taskList"
          :key="task.id"
          class="task-item"
          :class="{ completed: task.status === 'completed' }"
        >
          <div class="task-main">
            <div class="task-info">
              <div class="task-id">#{{ task.id }}</div>
              <div class="task-status" :class="task.status">
                {{ getStatusText(task.status) }}
              </div>
            </div>

            <div class="task-prompt">
              {{ task.params ? getPromptFromParams(task.params) : '' }}
            </div>

            <div class="task-progress" v-if="task.status === 'processing'">
              <el-progress :percentage="task.progress" :stroke-width="6" />
            </div>
          </div>

          <div class="task-actions">
            <template v-if="task.status === 'completed'">
              <el-button size="small" @click="preview3D(task)">预览</el-button>
              <el-button size="small" type="primary" @click="download3D(task)">下载</el-button>
            </template>
            <template v-else>
              <el-button size="small" @click="deleteTask(task.id)">删除</el-button>
            </template>
          </div>
        </div>
      </div>

      <!-- 加载更多 -->
      <div class="load-more" v-if="hasMore">
        <el-button size="small" @click="loadMoreTasks">加载更多</el-button>
      </div>
    </div>

    <!-- 3D预览弹窗 -->
    <el-dialog
      v-model="previewVisible"
      title="3D模型预览"
      width="90%"
      :before-close="closePreview"
      class="mobile-dialog"
    >
      <div class="preview-container">
        <div id="three-container" class="three-container">
          <div class="preview-placeholder">
            <i class="iconfont icon-3d"></i>
            <p>3D模型预览</p>
          </div>
        </div>
      </div>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="closePreview">关闭</el-button>
          <el-button type="primary" @click="downloadCurrentModel">下载模型</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import ImageUpload from '@/components/ImageUpload.vue'
import { httpGet, httpPost } from '@/utils/http'
import { ElMessage, ElMessageBox } from 'element-plus'
import { computed, nextTick, onMounted, ref } from 'vue'

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
const hasMore = ref(true)

// 平台配置
const platforms = [
  {
    key: 'gitee',
    name: '魔力方舟',
    icon: 'icon-gitee',
  },
  {
    key: 'tencent',
    name: '腾讯混元',
    icon: 'icon-tencent',
  },
]

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
const handleImageChange = (files) => {
  currentImage.value = files
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
      loadTasks(true)
    } else {
      ElMessage.error(response.message || '创建任务失败')
    }
  } catch (error) {
    ElMessage.error('创建任务失败：' + error.message)
  } finally {
    generating.value = false
  }
}

const loadTasks = async (reset = false) => {
  try {
    if (reset) {
      currentPage.value = 1
      taskList.value = []
    }

    const response = await httpGet('/api/3d/jobs', {
      page: currentPage.value,
      page_size: pageSize.value,
    })

    if (response.code === 0) {
      if (reset) {
        taskList.value = response.data.list
      } else {
        taskList.value.push(...response.data.list)
      }
      total.value = response.data.total
      hasMore.value = taskList.value.length < total.value
    }
  } catch (error) {
    ElMessage.error('加载任务列表失败：' + error.message)
  }
}

const refreshTasks = () => {
  loadTasks(true)
}

const loadMoreTasks = () => {
  if (hasMore.value) {
    currentPage.value++
    loadTasks()
  }
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
      loadTasks(true)
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

  nextTick(() => {
    initThreeJS(task)
  })
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

// Three.js 初始化
const initThreeJS = (task) => {
  // TODO: 实现Three.js 3D模型预览
  console.log('初始化Three.js预览:', task)
}

// 生命周期
onMounted(() => {
  loadTasks(true)
})
</script>

<style lang="scss" scoped>
.mobile-threed-create {
  min-height: 100vh;
  background: #f5f5f5;
  padding-bottom: 20px;
}

.top-nav {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px;
  background: white;
  border-bottom: 1px solid #e4e7ed;
  position: sticky;
  top: 0;
  z-index: 100;

  .nav-left {
    width: 24px;
    height: 24px;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;

    i {
      font-size: 20px;
      color: #333;
    }
  }

  .nav-title {
    font-size: 18px;
    font-weight: 500;
    color: #333;
  }

  .nav-right {
    width: 24px;
  }
}

.platform-selector {
  background: white;
  margin: 16px 20px;
  border-radius: 12px;
  padding: 20px;

  .selector-tabs {
    display: flex;
    gap: 12px;

    .selector-tab {
      flex: 1;
      display: flex;
      flex-direction: column;
      align-items: center;
      gap: 8px;
      padding: 16px 12px;
      border-radius: 8px;
      border: 2px solid #e4e7ed;
      cursor: pointer;
      transition: all 0.3s;

      &.active {
        border-color: #409eff;
        background: #ecf5ff;
      }

      .tab-icon {
        width: 32px;
        height: 32px;
        display: flex;
        align-items: center;
        justify-content: center;

        i {
          font-size: 24px;
          color: #409eff;
        }
      }

      .tab-name {
        font-size: 14px;
        color: #333;
        text-align: center;
      }
    }
  }
}

.params-section {
  background: white;
  margin: 16px 20px;
  border-radius: 12px;
  padding: 20px;

  .param-group {
    margin-bottom: 24px;

    .param-label {
      font-size: 16px;
      font-weight: 500;
      color: #333;
      margin-bottom: 12px;
    }

    .image-upload-area {
      border: 2px dashed #d9d9d9;
      border-radius: 8px;
      padding: 20px;
      text-align: center;
      transition: border-color 0.3s;

      &:hover {
        border-color: #409eff;
      }
    }

    .prompt-input {
      .el-textarea {
        .el-textarea__inner {
          border-radius: 8px;
          border: 1px solid #d9d9d9;
        }
      }
    }

    .model-selector {
      display: grid;
      grid-template-columns: repeat(2, 1fr);
      gap: 12px;

      .model-option {
        padding: 16px 12px;
        border: 2px solid #e4e7ed;
        border-radius: 8px;
        text-align: center;
        cursor: pointer;
        transition: all 0.3s;

        &.active {
          border-color: #409eff;
          background: #ecf5ff;
        }

        .model-name {
          font-size: 14px;
          color: #333;
          margin-bottom: 4px;
        }

        .model-power {
          font-size: 12px;
          color: #409eff;
          font-weight: 500;
        }
      }
    }
  }
}

.power-info {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px;
  background: #f0f9ff;
  border-radius: 8px;
  border: 1px solid #b3d8ff;

  .power-label {
    font-size: 16px;
    color: #333;
  }

  .power-value {
    font-size: 20px;
    font-weight: bold;
    color: #409eff;
  }
}

.generate-section {
  margin: 16px 20px;

  .generate-btn {
    width: 100%;
    height: 48px;
    font-size: 16px;
    border-radius: 12px;
  }
}

.task-section {
  background: white;
  margin: 16px 20px;
  border-radius: 12px;
  padding: 20px;

  .section-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;

    h3 {
      margin: 0;
      font-size: 18px;
      color: #333;
    }
  }
}

.task-list {
  .task-item {
    border: 1px solid #e4e7ed;
    border-radius: 8px;
    padding: 16px;
    margin-bottom: 16px;

    &.completed {
      border-color: #67c23a;
      background: #f0f9ff;
    }

    .task-main {
      margin-bottom: 16px;

      .task-info {
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

      .task-prompt {
        color: #666;
        margin-bottom: 12px;
        line-height: 1.4;
        font-size: 14px;
      }
    }

    .task-actions {
      display: flex;
      gap: 8px;

      .el-button {
        flex: 1;
      }
    }
  }
}

.load-more {
  text-align: center;
  margin-top: 20px;
}

.preview-container {
  .three-container {
    width: 100%;
    height: 300px;
    background: #f0f0f0;
    border-radius: 8px;
    display: flex;
    align-items: center;
    justify-content: center;

    .preview-placeholder {
      text-align: center;
      color: #666;

      i {
        font-size: 48px;
        margin-bottom: 12px;
        display: block;
      }

      p {
        margin: 0;
        font-size: 14px;
      }
    }
  }
}

// 移动端弹窗样式
.mobile-dialog {
  :deep(.el-dialog) {
    margin: 20px;
    border-radius: 12px;
  }

  :deep(.el-dialog__header) {
    padding: 20px 20px 0;
  }

  :deep(.el-dialog__body) {
    padding: 20px;
  }

  :deep(.el-dialog__footer) {
    padding: 0 20px 20px;

    .dialog-footer {
      display: flex;
      gap: 12px;

      .el-button {
        flex: 1;
      }
    }
  }
}
</style>
