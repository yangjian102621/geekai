// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import { checkSession } from '@/store/cache'
import { showMessageError, showMessageOK } from '@/utils/dialog'
import { httpDownload, httpGet, httpPost } from '@/utils/http'
import { replaceImg, substr } from '@/utils/libs'
import { ElMessageBox } from 'element-plus'
import { defineStore } from 'pinia'
import { computed, nextTick, reactive, ref } from 'vue'

export const useJimengStore = defineStore('jimeng', () => {
  // 当前激活的功能分类和具体功能
  const activeCategory = ref('image_generation')
  const activeFunction = ref('text_to_image')
  const useImageInput = ref(false)

  // 新增：全局提示词
  const currentPrompt = ref('')

  // 共同状态
  const loading = ref(false)
  const submitting = ref(false)
  const page = ref(1)
  const pageSize = ref(10)
  const total = ref(0)
  const taskFilter = ref('all')
  const currentList = ref([])
  const isOver = ref(false)

  // 用户信息
  const isLogin = ref(false)
  const userPower = ref(100)

  // 视频预览
  const showDialog = ref(false)
  const currentVideoUrl = ref('')

  // 功能分类配置
  const categories = [
    { key: 'image_generation', name: '图片生成' },
    { key: 'image_editing', name: 'AI修图' },
    { key: 'image_effects', name: '图像特效' },
    { key: 'video_generation', name: '视频生成' },
  ]

  // 新增：动态获取算力消耗配置
  const powerConfig = reactive({})

  // 功能配置
  const functions = reactive([
    {
      key: 'text_to_image',
      name: '文生图',
      category: 'image_generation',
      needsPrompt: true,
      needsImage: false,
      power: 20,
    },
    {
      key: 'image_to_image',
      name: '图生图',
      category: 'image_generation',
      needsPrompt: true,
      needsImage: true,
      power: 30,
    },
    {
      key: 'image_edit',
      name: '图像编辑',
      category: 'image_editing',
      needsPrompt: true,
      needsImage: true,
      multiple: true,
      power: 25,
    },
    {
      key: 'image_effects',
      name: '图像特效',
      category: 'image_effects',
      needsPrompt: false,
      needsImage: true,
      power: 15,
    },
    {
      key: 'text_to_video',
      name: '文生视频',
      category: 'video_generation',
      needsPrompt: true,
      needsImage: false,
      power: 100,
    },
    {
      key: 'image_to_video',
      name: '图生视频',
      category: 'video_generation',
      needsPrompt: true,
      needsImage: true,
      multiple: true,
      power: 120,
    },
  ])

  // 动态设置算力消耗
  const setFunctionPowers = (config) => {
    functions.forEach((f) => {
      if (config[f.key] !== undefined) {
        f.power = config[f.key]
      }
    })
  }

  // 各功能的参数
  const textToImageParams = reactive({
    size: '1328x1328',
    scale: 2.5,
    seed: -1,
    use_pre_llm: true,
  })

  const imageToImageParams = reactive({
    image_input: '',
    size: '1328x1328',
    gpen: 0.4,
    skin: 0.3,
    skin_unifi: 0,
    gen_mode: 'creative',
    seed: -1,
  })

  const imageEditParams = reactive({
    image_urls: [],
    scale: 0.5,
    seed: -1,
  })

  const imageEffectsParams = reactive({
    image_input1: '',
    template_id: '',
    size: '1328x1328',
  })

  const textToVideoParams = reactive({
    aspect_ratio: '16:9',
    seed: -1,
  })

  const imageToVideoParams = reactive({
    image_urls: [],
    aspect_ratio: '16:9',
    seed: -1,
  })

  // 计算属性
  const currentFunction = computed(() => {
    return functions.find((f) => f.key === activeFunction.value) || functions[0]
  })

  const currentFunctions = computed(() => {
    return functions.filter((f) => f.category === activeCategory.value)
  })

  const needsPrompt = computed(() => currentFunction.value.needsPrompt)
  const needsImage = computed(() => currentFunction.value.needsImage)
  const needsMultipleImages = computed(() => currentFunction.value.multiple)
  const currentPowerCost = computed(() => currentFunction.value.power)

  // 初始化方法
  const init = async () => {
    try {
      // 获取算力消耗配置
      const powerRes = await httpGet('/api/jimeng/power-config')
      if (powerRes.data) {
        Object.assign(powerConfig, powerRes.data)
        setFunctionPowers(powerRes.data)
      }
      const user = await checkSession()
      isLogin.value = true
      userPower.value = user.power
      // 获取任务列表
      await fetchData(1)
      // 开始轮询
      startPolling()
    } catch (error) {
      console.error('初始化失败:', error)
    }
  }

  // 切换功能分类
  const switchCategory = (category) => {
    activeCategory.value = category
    const categoryFunctions = functions.filter((f) => f.category === category)
    if (categoryFunctions.length > 0) {
      if (category === 'image_generation') {
        activeFunction.value = useImageInput.value ? 'image_to_image' : 'text_to_image'
      } else if (category === 'video_generation') {
        activeFunction.value = useImageInput.value ? 'image_to_video' : 'text_to_video'
      } else {
        activeFunction.value = categoryFunctions[0].key
      }
    }
  }

  // 切换输入模式
  const switchInputMode = () => {
    if (activeCategory.value === 'image_generation') {
      activeFunction.value = useImageInput.value ? 'image_to_image' : 'text_to_image'
    } else if (activeCategory.value === 'video_generation') {
      activeFunction.value = useImageInput.value ? 'image_to_video' : 'text_to_video'
    }
  }

  // 切换功能
  const switchFunction = (functionKey) => {
    activeFunction.value = functionKey
  }

  // 获取当前算力消耗
  const getCurrentPowerCost = () => {
    return currentFunction.value.power
  }

  // 获取功能名称
  const getFunctionName = (type) => {
    const func = functions.find((f) => f.key === type)
    return func ? func.name : type
  }

  // 获取任务状态文本
  const getTaskStatusText = (status) => {
    const statusMap = {
      in_queue: '任务排队中',
      generating: '任务执行中',
      success: '任务成功',
      failed: '任务失败',
      canceled: '任务已取消',
    }
    return statusMap[status] || status
  }

  // 获取状态类型
  const getTaskType = (type) => {
    const typeMap = {
      text_to_image: 'primary',
      image_to_image: 'primary',
      image_edit: 'primary',
      image_effects: 'primary',
      text_to_video: 'success',
      image_to_video: 'success',
    }
    return typeMap[type] || 'primary'
  }

  // 切换任务筛选
  const switchTaskFilter = (filter) => {
    taskFilter.value = filter
    isOver.value = false
    fetchData(1)
  }

  // 轮询定时器
  let pollHandler = null
  // 获取任务列表
  const fetchData = async (pageNum = 1) => {
    try {
      loading.value = true
      page.value = pageNum

      const response = await httpPost('/api/jimeng/jobs', {
        page: pageNum,
        page_size: pageSize.value,
        filter: taskFilter.value,
      })

      const data = response.data
      if (!data.items || data.items.length === 0) {
        isOver.value = true
        if (pageNum === 1) {
          currentList.value = []
        }
        return
      }

      total.value = data.total || 0
      if (data.items.length < pageSize.value) {
        isOver.value = true
      }
      if (pageNum === 1) {
        currentList.value = data.items
      } else {
        currentList.value = currentList.value.concat(data.items)
      }
    } catch (error) {
      showMessageError('获取任务列表失败:' + error.message)
    } finally {
      loading.value = false
    }
  }

  // 简单轮询逻辑
  const startPolling = () => {
    if (pollHandler) {
      clearInterval(pollHandler)
    }
    pollHandler = setInterval(async () => {
      const response = await httpPost('/api/jimeng/jobs', {
        page: 1,
        page_size: 20,
      })
      const data = response.data
      if (data.items.length === 0) {
        stopPolling()
        return
      }

      const todoList = data.items.filter(
        (item) => item.status === 'in_queue' || item.status === 'generating'
      )
      // 更新当前列表
      currentList.value.forEach((item) => {
        const index = data.items.findIndex((i) => i.id === item.id)
        if (index !== -1) {
          Object.assign(item, data.items[index])
        }
      })
      if (todoList.length === 0) {
        stopPolling()
      }
    }, 3000)
  }

  const stopPolling = () => {
    if (pollHandler) {
      clearInterval(pollHandler)
      pollHandler = null
    }
  }

  // 提交任务
  const submitTask = async () => {
    if (!isLogin.value) {
      showMessageError('请先登录')
      return
    }
    if (userPower.value < currentPowerCost.value) {
      showMessageError('算力不足')
      return
    }
    // 新增：除图像特效外，其他任务类型必须有提示词
    if (activeFunction.value !== 'image_effects' && !currentPrompt.value) {
      showMessageError('提示词不能为空')
      return
    }
    try {
      submitting.value = true
      let requestData = { task_type: activeFunction.value, prompt: currentPrompt.value }
      switch (activeFunction.value) {
        case 'text_to_image':
          Object.assign(requestData, {
            width: parseInt(textToImageParams.size.split('x')[0]),
            height: parseInt(textToImageParams.size.split('x')[1]),
            scale: textToImageParams.scale,
            seed: textToImageParams.seed,
            use_pre_llm: textToImageParams.use_pre_llm,
          })
          break
        case 'image_to_image':
          Object.assign(requestData, {
            image_input: imageToImageParams.image_input,
            width: parseInt(imageToImageParams.size.split('x')[0]),
            height: parseInt(imageToImageParams.size.split('x')[1]),
            gpen: imageToImageParams.gpen,
            skin: imageToImageParams.skin,
            skin_unifi: imageToImageParams.skin_unifi,
            gen_mode: imageToImageParams.gen_mode,
            seed: imageToImageParams.seed,
          })
          break
        case 'image_edit':
          Object.assign(requestData, {
            image_urls: imageEditParams.image_urls,
            scale: imageEditParams.scale,
            seed: imageEditParams.seed,
          })
          break
        case 'image_effects':
          Object.assign(requestData, {
            image_input: imageEffectsParams.image_input1,
            template_id: imageEffectsParams.template_id,
            width: parseInt(imageEffectsParams.size.split('x')[0]),
            height: parseInt(imageEffectsParams.size.split('x')[1]),
          })
          break
        case 'text_to_video':
          Object.assign(requestData, {
            aspect_ratio: textToVideoParams.aspect_ratio,
            seed: textToVideoParams.seed,
          })
          break
        case 'image_to_video':
          Object.assign(requestData, {
            image_urls: imageToVideoParams.image_urls,
            aspect_ratio: imageToVideoParams.aspect_ratio,
            seed: imageToVideoParams.seed,
          })
          break
      }
      const response = await httpPost('/api/jimeng/task', requestData)
      if (response.data) {
        showMessageOK('任务提交成功')
        isOver.value = false
        await fetchData(1)
        startPolling()
      }
    } catch (error) {
      console.error('提交任务失败:', error)
      showMessageError(error.message || '提交任务失败')
    } finally {
      submitting.value = false
    }
  }

  const downloadFile = async (item) => {
    const url = replaceImg(item.video_url || item.img_url)
    const downloadURL = `${import.meta.env.VITE_API_HOST}/api/download?url=${url}`
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

  // 重试任务
  const retryTask = async (taskId) => {
    try {
      const response = await httpGet(`/api/jimeng/retry?id=${taskId}`)
      if (response.data) {
        showMessageOK('重试任务已提交')
        isOver.value = false
        await fetchData(1)
        startPolling()
      }
    } catch (error) {
      console.error('重试任务失败:', error)
      showMessageError(error.message || '重试任务失败')
    }
  }

  // 删除任务
  const removeJob = async (item) => {
    try {
      await ElMessageBox.confirm('确定要删除这个任务吗？', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      })

      const response = await httpGet('/api/jimeng/remove', { id: item.id })
      if (response.data) {
        showMessageOK('删除成功')
        await fetchData(1)
      }
    } catch (error) {
      if (error !== 'cancel') {
        console.error('删除任务失败:', error)
        showMessageError(error.message || '删除任务失败')
      }
    }
  }

  // 播放视频
  const playVideo = (item) => {
    currentVideoUrl.value = item.video_url
    showDialog.value = true
  }

  // 画同款功能
  const drawSame = (item) => {
    // 联动功能开关
    if (item.type === 'text_to_image' || item.type === 'image_to_image') {
      activeCategory.value = 'image_generation'
      useImageInput.value = item.type === 'image_to_image'
    } else if (item.type === 'text_to_video' || item.type === 'image_to_video') {
      activeCategory.value = 'video_generation'
      useImageInput.value = item.type === 'image_to_video'
    } else if (item.type === 'image_edit') {
      activeCategory.value = 'image_editing'
    } else if (item.type === 'image_effects') {
      activeCategory.value = 'image_effects'
    }
    switchFunction(item.type)
    nextTick(() => {
      currentPrompt.value = item.prompt
    })
    if (item.type === 'text_to_image') {
      if (item.width && item.height) {
        textToImageParams.size = `${item.width}x${item.height}`
      }
      if (item.scale) textToImageParams.scale = item.scale
      if (item.seed) textToImageParams.seed = item.seed
      if (item.use_pre_llm !== undefined) textToImageParams.use_pre_llm = item.use_pre_llm
    } else if (item.type === 'image_to_image') {
      if (item.image_input) imageToImageParams.image_input = item.image_input
      if (item.width && item.height) {
        imageToImageParams.size = `${item.width}x${item.height}`
      }
      if (item.gpen) imageToImageParams.gpen = item.gpen
      if (item.skin) imageToImageParams.skin = item.skin
      if (item.skin_unifi) imageToImageParams.skin_unifi = item.skin_unifi
      if (item.gen_mode) imageToImageParams.gen_mode = item.gen_mode
      if (item.seed) imageToImageParams.seed = item.seed
    } else if (item.type === 'image_edit') {
      if (item.image_urls) imageEditParams.image_urls = item.image_urls
      if (item.scale) imageEditParams.scale = item.scale
      if (item.seed) imageEditParams.seed = item.seed
    } else if (item.type === 'image_effects') {
      if (item.image_input1) imageEffectsParams.image_input1 = item.image_input1
      if (item.template_id) imageEffectsParams.template_id = item.template_id
      if (item.width && item.height) {
        imageEffectsParams.size = `${item.width}x${item.height}`
      }
    } else if (item.type === 'text_to_video') {
      if (item.aspect_ratio) textToVideoParams.aspect_ratio = item.aspect_ratio
      if (item.seed) textToVideoParams.seed = item.seed
    } else if (item.type === 'image_to_video') {
      if (item.image_urls) imageToVideoParams.image_urls = item.image_urls
      if (item.aspect_ratio) imageToVideoParams.aspect_ratio = item.aspect_ratio
      if (item.seed) imageToVideoParams.seed = item.seed
    }
    showMessageOK('已填入全部参数，可直接生成同款')
  }

  // 页面卸载时清理轮询
  const cleanup = () => {
    stopPolling()
  }

  // 返回所有状态和方法
  return {
    // 状态
    activeCategory,
    activeFunction,
    useImageInput,
    loading,
    submitting,
    page,
    pageSize,
    total,
    taskFilter,
    currentList,
    isOver,
    isLogin,
    userPower,
    showDialog,
    currentVideoUrl,

    // 配置
    categories,
    functions,
    currentFunctions,

    // 参数
    currentPrompt,
    textToImageParams,
    imageToImageParams,
    imageEditParams,
    imageEffectsParams,
    textToVideoParams,
    imageToVideoParams,

    // 计算属性
    currentFunction,
    needsPrompt,
    needsImage,
    needsMultipleImages,
    currentPowerCost,

    // 方法
    init,
    switchCategory,
    switchFunction,
    switchInputMode,
    getCurrentPowerCost,
    getFunctionName,
    getTaskStatusText,
    getTaskType,
    switchTaskFilter,
    fetchData,
    submitTask,
    downloadFile,
    retryTask,
    removeJob,
    playVideo,
    cleanup,
    drawSame,

    // 工具函数
    substr,
    replaceImg,
  }
})

export const imageSizeOptions = [
  { label: '1:1 (1328x1328)', value: '1328x1328' },
  { label: '3:2 (1584x1056)', value: '1584x1056' },
  { label: '2:3 (1056x1584)', value: '1056x1584' },
  { label: '4:3 (1472x1104)', value: '1472x1104' },
  { label: '3:4 (1104x1472)', value: '1104x1472' },
  { label: '16:9 (1664x936)', value: '1664x936' },
  { label: '9:16 (936x1664)', value: '936x1664' },
  { label: '21:9 (2016x864)', value: '2016x864' },
  { label: '9:21 (864x2016)', value: '864x2016' },
]

export const videoAspectRatioOptions = [
  { label: '1:1 (正方形)', value: '1:1' },
  { label: '16:9 (横版)', value: '16:9' },
  { label: '9:16 (竖版)', value: '9:16' },
]
