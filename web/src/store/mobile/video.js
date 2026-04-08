import { defineStore } from 'pinia'
import { ref, reactive, watch } from 'vue'
import { httpGet, httpPost } from '@/utils/http'
import { showMessageOK, showMessageError, showLoading, closeLoading } from '@/utils/dialog'
import { getSystemInfo } from '@/store/cache'

export const useVideoStore = defineStore('video', () => {
  // 状态
  const activeVideoType = ref('luma')
  const loading = ref(false)
  const generating = ref(false)
  const isGenerating = ref(false)
  const listLoading = ref(false)
  const listFinished = ref(false)
  const currentList = ref([])
  const showVideoDialog = ref(false)
  const currentVideoUrl = ref('')

  // Luma 参数
  const lumaParams = reactive({
    prompt: '',
    image: '',
    image_tail: '',
    loop: false,
    expand_prompt: false,
  })
  const lumaUseImageMode = ref(false)
  const lumaStartImage = ref([])
  const lumaEndImage = ref([])

  // KeLing 参数
  const kelingParams = reactive({
    aspect_ratio: '16:9',
    model: 'kling-v1-6',
    duration: '5',
    mode: 'std',
    cfg_scale: 0.5,
    prompt: '',
    negative_prompt: '',
    image: '',
    image_tail: '',
    camera_control: {
      type: '',
      config: {
        horizontal: 0,
        vertical: 0,
        pan: 0,
        tilt: 0,
        roll: 0,
        zoom: 0,
      },
    },
  })
  const kelingUseImageMode = ref(false)
  const kelingStartImage = ref([])
  const kelingEndImage = ref([])

  // 选项数据
  const aspectRatioOptions = ['16:9', '9:16', '1:1', '4:3']
  const modelOptions = [
    { label: '可灵 1.6', value: 'kling-v1-6' },
    { label: '可灵 1.5', value: 'kling-v1-5' },
    { label: '可灵 1.0', value: 'kling-v1' },
  ]
  const durationOptions = ['5', '10']
  const modeOptions = ['std', 'pro']
  const cameraControlOptions = [
    '',
    'simple',
    'down_back',
    'forward_up',
    'right_turn_forward',
    'left_turn_forward',
  ]
  const getCameraControlLabel = (option) => {
    const labelMap = {
      '': '请选择',
      simple: '简单运镜',
      down_back: '下移拉远',
      forward_up: '推进上移',
      right_turn_forward: '右旋推进',
      left_turn_forward: '左旋推进',
    }
    return labelMap[option] || option
  }

  // 页面数据
  const page = ref(1)
  const pageSize = ref(10)
  const total = ref(0)
  const lumaPowerCost = ref(0)
  const kelingPowerCost = ref(0)
  const taskPulling = ref(true)
  const keLingPowers = ref({})

  // 监听器：当可灵参数变化时更新算力
  watch(
    () => [kelingParams.model, kelingParams.mode, kelingParams.duration],
    () => {
      updateModelPower()
    },
    { deep: true }
  )

  // 方法
  const updateModelPower = () => {
    // 根据模型、模式、时长计算算力消耗
    const key = `${kelingParams.model}_${kelingParams.mode}_${kelingParams.duration}`
    kelingPowerCost.value = keLingPowers.value[key] || 10
  }
  watch(
    () => [kelingParams.model, kelingParams.mode, kelingParams.duration],
    () => {
      updateModelPower()
    },
    { deep: true }
  )

  // 监听器：当可灵参数变化时更新算力
  watch(
    () => [kelingParams.model, kelingParams.mode, kelingParams.duration],
    () => {
      updateModelPower()
    },
    { deep: true }
  )

  const switchVideoType = (type) => {
    activeVideoType.value = type
  }
  const handleLumaStartImageUpload = (e) => {
    if (e.target.files[0]) {
      uploadLumaStartImage({ file: e.target.files[0], name: e.target.files[0].name })
    }
  }
  const handleLumaEndImageUpload = (e) => {
    if (e.target.files[0]) {
      uploadLumaEndImage({ file: e.target.files[0], name: e.target.files[0].name })
    }
  }
  const handleKelingStartImageUpload = (e) => {
    if (e.target.files[0]) {
      uploadKelingStartImage({ file: e.target.files[0], name: e.target.files[0].name })
    }
  }
  const handleKelingEndImageUpload = (e) => {
    if (e.target.files[0]) {
      uploadKelingEndImage({ file: e.target.files[0], name: e.target.files[0].name })
    }
  }

  const generatePrompt = async () => {
    if (isGenerating.value) return

    const prompt = activeVideoType.value === 'luma' ? lumaParams.prompt : kelingParams.prompt
    if (!prompt) {
      return showMessageError('请输入原始提示词')
    }

    isGenerating.value = true
    try {
      const res = await httpPost('/api/prompt/video', { prompt })
      if (activeVideoType.value === 'luma') {
        lumaParams.prompt = res.data
      } else {
        kelingParams.prompt = res.data
      }
    } catch (error) {
      showMessageError('生成提示词失败：' + error.message)
    } finally {
      isGenerating.value = false
    }
  }
  const toggleLumaImageMode = () => {
    if (!lumaUseImageMode.value) {
      lumaParams.image = ''
      lumaParams.image_tail = ''
      lumaStartImage.value = []
      lumaEndImage.value = []
    }
  }
  const toggleKelingImageMode = () => {
    if (!kelingUseImageMode.value) {
      kelingParams.image = ''
      kelingParams.image_tail = ''
      kelingStartImage.value = []
      kelingEndImage.value = []
    }
  }
  const uploadLumaStartImage = (file) => {
    uploadImage(file, (url) => {
      lumaParams.image = url
    })
  }
  const uploadLumaEndImage = (file) => {
    uploadImage(file, (url) => {
      lumaParams.image_tail = url
    })
  }
  const uploadKelingStartImage = (file) => {
    uploadImage(file, (url) => {
      kelingParams.image = url
    })
  }
  const uploadKelingEndImage = (file) => {
    uploadImage(file, (url) => {
      kelingParams.image_tail = url
    })
  }
  const uploadImage = (file, callback) => {
    const formData = new FormData()
    formData.append('file', file.file, file.name)
    showLoading('正在上传图片...')
    httpPost('/api/upload', formData)
      .then((res) => {
        callback(res.data.url)
        showMessageOK('图片上传成功')
      })
      .catch((e) => {
        showMessageError('图片上传失败:' + e.message)
      })
      .finally(() => {
        closeLoading()
      })
  }
  const createLumaVideo = () => {
    if (!lumaParams.prompt.trim()) {
      showMessageError('请输入视频提示词')
      return
    }
    generating.value = true
    const params = {
      ...lumaParams,
      task_type: 'luma',
    }
    httpPost('/api/video/create', params)
      .then(() => {
        fetchData(1)
        taskPulling.value = true
        showMessageOK('创建任务成功')
      })
      .catch((e) => {
        showMessageError('创建任务失败：' + e.message)
      })
      .finally(() => {
        generating.value = false
      })
  }
  const createKelingVideo = () => {
    if (!kelingParams.prompt.trim()) {
      showMessageError('请输入视频提示词')
      return
    }
    generating.value = true
    const params = {
      ...kelingParams,
      task_type: 'keling',
    }
    httpPost('/api/video/create', params)
      .then(() => {
        fetchData(1)
        taskPulling.value = true
        showMessageOK('创建任务成功')
      })
      .catch((e) => {
        showMessageError('创建任务失败：' + e.message)
      })
      .finally(() => {
        generating.value = false
      })
  }
  const fetchData = (_page) => {
    if (_page) {
      page.value = _page
    }
    listLoading.value = true
    httpGet('/api/video/list', { page: page.value, page_size: pageSize.value })
      .then((res) => {
        total.value = res.data.total
        let needPull = false
        const items = []
        for (let v of res.data.items) {
          if (v.progress === 0 || v.progress === 102) {
            needPull = true
          }
          items.push(v)
        }
        listLoading.value = false
        taskPulling.value = needPull
        if (page.value === 1) {
          currentList.value = items
        } else {
          currentList.value.push(...items)
        }
        if (items.length < pageSize.value) {
          listFinished.value = true
        }
      })
      .catch((e) => {
        listLoading.value = false
        showMessageError('获取作品列表失败：' + e.message)
      })
  }
  const fetchUserPower = async () => {
    try {
      // 获取系统信息，更新算力配置
      const sysInfo = await getSystemInfo()
      lumaPowerCost.value = sysInfo.data.luma_power || 10
      keLingPowers.value = sysInfo.data.keling_powers || {}
      updateModelPower()
    } catch (error) {
      console.error('获取用户算力失败:', error)
      // 设置默认值
      lumaPowerCost.value = 10
      kelingPowerCost.value = 15
    }
  }
  const loadMore = () => {
    page.value++
    fetchData()
  }
  const playVideo = (item) => {
    currentVideoUrl.value = item.video_url
    showVideoDialog.value = true
  }
  const downloadVideo = (item) => {
    item.downloading = true
    const link = document.createElement('a')
    link.href = item.video_url
    link.download = item.title || 'video.mp4'
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    item.downloading = false
    showMessageOK('开始下载')
  }
  const removeJob = (item) => {
    // 建议在页面层处理弹窗，store 只负责数据和业务
  }

  return {
    // 状态
    activeVideoType,
    loading,
    generating,
    isGenerating,
    listLoading,
    listFinished,
    currentList,
    showVideoDialog,
    currentVideoUrl,
    lumaParams,
    lumaUseImageMode,
    lumaStartImage,
    lumaEndImage,
    kelingParams,
    kelingUseImageMode,
    kelingStartImage,
    kelingEndImage,
    aspectRatioOptions,
    modelOptions,
    durationOptions,
    modeOptions,
    cameraControlOptions,
    getCameraControlLabel,
    page,
    pageSize,
    total,
    lumaPowerCost,
    kelingPowerCost,
    taskPulling,
    keLingPowers,
    // 方法
    updateModelPower,
    switchVideoType,
    handleLumaStartImageUpload,
    handleLumaEndImageUpload,
    handleKelingStartImageUpload,
    handleKelingEndImageUpload,
    generatePrompt,
    toggleLumaImageMode,
    toggleKelingImageMode,
    uploadLumaStartImage,
    uploadLumaEndImage,
    uploadKelingStartImage,
    uploadKelingEndImage,
    uploadImage,
    createLumaVideo,
    createKelingVideo,
    fetchData,
    fetchUserPower,
    loadMore,
    playVideo,
    downloadVideo,
    removeJob,
  }
})
