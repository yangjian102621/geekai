// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import { checkSession } from '@/store/cache'
import { useSharedStore } from '@/store/sharedata'
import { showMessageError, showMessageOK } from '@/utils/dialog'
import { httpDownload, httpGet, httpPost } from '@/utils/http'
import { replaceImg, substr } from '@/utils/libs'
import { ElMessageBox } from 'element-plus'
import { defineStore } from 'pinia'
import { computed, reactive, ref } from 'vue'

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

  // 登录弹窗
  const shareStore = useSharedStore()

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
    image_input: '',
    scale: 0.5,
    seed: -1,
  })

  const imageEffectsParams = reactive({
    image_input: '',
    template_id: '',
    size: '1328x1328',
  })

  const textToVideoParams = reactive({
    aspect_ratio: '16:9',
    seed: -1,
  })

  const imageToVideoParams = reactive({
    image_input: [],
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
      shareStore.setShowLoginDialog(true)
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
            image_input: imageToImageParams.image_input[0],
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
            image_input: imageEditParams.image_input[0],
            scale: imageEditParams.scale,
            seed: imageEditParams.seed,
          })
          break
        case 'image_effects':
          Object.assign(requestData, {
            image_input: imageEffectsParams.image_input[0],
            template_id: imageEffectsParams.template_id,
            width: parseInt(imageEffectsParams.size.split('x')[0]),
            height: parseInt(imageEffectsParams.size.split('x')[1]),
            prompt: imageEffectsParams.prompt,
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
            image_urls: imageToVideoParams.image_input,
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

export const imageEffectsTemplateOptions = [
  {
    label: '毛毡3D拍立得风格',
    value: 'felt_3d_polaroid',
    preview: '/images/jimeng/templates/felt_3d_polaroid.png',
  },
  { label: '像素世界风', value: 'my_world', preview: '/images/jimeng/templates/my_world.png' },
  {
    label: '像素世界-万物通用版',
    value: 'my_world_universal',
    preview: '/images/jimeng/templates/my_world_universal.png',
  },
  {
    label: '盲盒玩偶风',
    value: 'plastic_bubble_figure',
    preview: '/images/jimeng/templates/plastic_bubble_figure.png',
  },
  {
    label: '塑料泡罩人偶-文字卡头版',
    value: 'plastic_bubble_figure_cartoon_text',
    preview: '/images/jimeng/templates/plastic_bubble_figure_cartoon_text.png',
  },
  {
    label: '毛绒玩偶风',
    value: 'furry_dream_doll',
    preview: '/images/jimeng/templates/furry_dream_doll.png',
  },
  {
    label: '迷你世界玩偶风',
    value: 'micro_landscape_mini_world',
    preview: '/images/jimeng/templates/micro_landscape_mini_world.png',
  },
  {
    label: '微型景观小世界-职业版',
    value: 'micro_landscape_mini_world_professional',
    preview: '/images/jimeng/templates/micro_landscape_mini_world_professional.png',
  },
  {
    label: '亚克力挂饰',
    value: 'acrylic_ornaments',
    preview: '/images/jimeng/templates/acrylic_ornaments.png',
  },
  {
    label: '毛毡钥匙扣',
    value: 'felt_keychain',
    preview: '/images/jimeng/templates/felt_keychain.png',
  },
  {
    label: 'Lofi 像素人物小卡',
    value: 'lofi_pixel_character_mini_card',
    preview: '/images/jimeng/templates/lofi_pixel_character_mini_card.png',
  },
  {
    label: '天使形象手办',
    value: 'angel_figurine',
    preview: '/images/jimeng/templates/angel_figurine.png',
  },
  {
    label: '躺在毛茸茸肚皮里',
    value: 'lying_in_fluffy_belly',
    preview: '/images/jimeng/templates/lying_in_fluffy_belly.png',
  },
  { label: '玻璃球', value: 'glass_ball', preview: '/images/jimeng/templates/glass_ball.png' },
]
