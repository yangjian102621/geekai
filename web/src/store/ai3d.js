import { checkSession } from '@/store/cache'
import { showMessageError } from '@/utils/dialog'
import { httpDownload, httpGet, httpPost } from '@/utils/http'
import { replaceImg } from '@/utils/libs'
import { ElMessage, ElMessageBox } from 'element-plus'
import { defineStore } from 'pinia'
import { computed, onMounted, ref } from 'vue'

export const useAI3DStore = defineStore('ai3d', () => {
  // 响应式数据
  const activePlatform = ref('gitee')
  const loading = ref(false)
  const previewVisible = ref(false)
  const currentPage = ref(1)
  const pageSize = ref(10)
  const total = ref(0)
  const taskList = ref([])
  const currentPreviewTask = ref({
    downloading: false,
  })
  const giteeAdvancedVisible = ref(false)
  const taskPulling = ref(false)

  const tencentDefaultForm = {
    text3d: false,
    prompt: '',
    image_url: '',
    model: '',
    file_format: '',
    enable_pbr: false,
    model_desc: '',
    power: 0,
  }

  const giteeDefaultForm = {
    prompt: '',
    image_url: '',
    model: '',
    file_format: '',
    texture: false,
    seed: 1234,
    num_inference_steps: 5,
    guidance_scale: 7.5,
    octree_resolution: 128,
    model_desc: '',
    power: 0,
  }

  const tencentForm = ref({ ...tencentDefaultForm })
  const giteeForm = ref({ ...giteeDefaultForm })
  const currentPower = ref(0)
  const tencentSupportedFormats = ref([])
  const giteeSupportedFormats = ref([])

  // 定时器引用
  let taskPullHandler = null

  const configs = ref({
    gitee: { models: [] },
    tencent: { models: [] },
  })

  // 计算属性
  const currentForm = computed(() =>
    activePlatform.value === 'tencent' ? tencentForm.value : giteeForm.value
  )
  const selectedModel = computed(() => currentForm.value.model)
  const currentPrompt = computed(() => currentForm.value.prompt)
  const currentImage = computed(() =>
    currentForm.value.image_url ? [{ url: currentForm.value.image_url }] : []
  )

  // 加载配置
  const loadConfigs = async () => {
    const response = await httpGet('/api/ai3d/configs')
    configs.value = response.data
  }

  const handleModelChange = (value) => {
    if (activePlatform.value === 'tencent') {
      const model = configs.value.tencent.models.find((m) => m.name === value)
      if (!model) return
      currentPower.value = model.power
      tencentForm.value.power = model.power
      tencentForm.value.model_desc = model.desc
      tencentForm.value.file_format = model.formats[0]
      tencentSupportedFormats.value = model.formats
    } else {
      const model = configs.value.gitee.models.find((m) => m.name === value)
      if (!model) return
      currentPower.value = model.power
      giteeForm.value.power = model.power
      giteeForm.value.model_desc = model.desc
      giteeForm.value.file_format = model.formats[0]
      giteeSupportedFormats.value = model.formats
    }
  }

  const handlePlatformChange = (value) => {
    activePlatform.value = value
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
      const response = await httpPost('/api/ai3d/generate', requestData)
      ElMessage.success('任务创建成功')
      await loadTasks()
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
        let needPull = false
        const items = response.data.items

        // 检查是否有进行中的任务
        for (let item of items) {
          if (item.status === 'pending' || item.status === 'processing') {
            needPull = true
            break
          }
        }

        taskPulling.value = needPull
        taskList.value = items
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
      pending: { text: '等待中', type: 'warning' },
      processing: { text: '处理中', type: 'primary' },
      success: { text: '已完成', type: 'success' },
      failed: { text: '失败', type: 'danger' },
    }
    return statusMap[status] || status
  }

  const getTaskCardClass = (status) => {
    if (status === 'success') return 'task-card-completed'
    if (status === 'processing') return 'task-card-processing'
    if (status === 'failed') return 'task-card-failed'
    return 'task-card-default'
  }

  const getPlatformIcon = (type) => {
    if (type === 'gitee') return 'iconfont icon-gitee'
    if (type === 'tencent') return 'iconfont icon-tencent'
    return 'iconfont icon-question'
  }

  const getPlatformName = (type) => {
    if (type === 'gitee') return 'Gitee 模力方舟'
    if (type === 'tencent') return '腾讯云混元3D'
    return '未知平台'
  }

  const getTaskPrompt = (task) => {
    return task.params.prompt ? task.params.prompt : '图生3D任务'
  }

  const getTaskImageUrl = (task) => {
    try {
      if (task.params) {
        const parsedParams = task.params
        return parsedParams.image_url || null
      }
      return null
    } catch (e) {
      return null
    }
  }

  const getTaskParams = (task) => {
    const parsedParams = task.params
    const params = []
    if (parsedParams.texture) params.push('纹理')
    if (parsedParams.enable_pbr) params.push('PBR材质')
    if (parsedParams.num_inference_steps)
      params.push(`迭代次数: ${parsedParams.num_inference_steps}`)
    if (parsedParams.guidance_scale) params.push(`引导系数: ${parsedParams.guidance_scale}`)
    if (parsedParams.octree_resolution) params.push(`精度: ${parsedParams.octree_resolution}`)
    if (parsedParams.seed) params.push(`Seed: ${parsedParams.seed}`)
    return params.join('，')
  }

  const startTaskPolling = () => {
    taskPullHandler = setInterval(() => {
      if (taskPulling.value) {
        loadTasks()
      }
    }, 5000)
  }

  const stopTaskPolling = () => {
    if (taskPullHandler) {
      clearInterval(taskPullHandler)
      taskPullHandler = null
    }
  }

  // 生命周期：加载配置与任务
  onMounted(() => {
    loadConfigs()
    checkSession()
      .then(() => {
        loadTasks()
        startTaskPolling()
      })
      .catch(() => {})
  })

  return {
    // 状态
    activePlatform,
    loading,
    previewVisible,
    currentPage,
    pageSize,
    total,
    taskList,
    currentPreviewTask,
    giteeAdvancedVisible,
    taskPulling,
    tencentForm,
    giteeForm,
    currentPower,
    tencentSupportedFormats,
    giteeSupportedFormats,
    configs,
    currentForm,
    selectedModel,
    currentPrompt,
    currentImage,
    // 方法
    loadConfigs,
    handleModelChange,
    handlePlatformChange,
    generate3D,
    loadTasks,
    refreshTasks,
    handlePageSizeChange,
    handleCurrentPageChange,
    deleteTask,
    preview3D,
    closePreview,
    downloadFile,
    downloadCurrentModel,
    getStatusText,
    getTaskCardClass,
    getPlatformIcon,
    getPlatformName,
    getTaskPrompt,
    getTaskImageUrl,
    getTaskParams,
    startTaskPolling,
    stopTaskPolling,
  }
})
