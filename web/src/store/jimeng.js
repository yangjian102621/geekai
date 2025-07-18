// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import nodata from '@/assets/img/no-data.png'
import { checkSession } from '@/store/cache'
import { closeLoading, showLoading, showMessageError, showMessageOK } from '@/utils/dialog'
import { httpGet, httpPost } from '@/utils/http'
import { replaceImg, substr, dateFormat } from '@/utils/libs'
import { ElMessage, ElMessageBox } from 'element-plus'
import { defineStore } from 'pinia'
import { computed, reactive, ref } from 'vue'

export const useJimengStore = defineStore('jimeng', () => {
  // 当前激活的功能分类和具体功能
  const activeCategory = ref('image_generation')
  const activeFunction = ref('text_to_image')
  const useImageInput = ref(false)
  
  // 共同状态
  const loading = ref(false)
  const submitting = ref(false)
  const list = ref([])
  const noData = ref(true)
  const page = ref(1)
  const pageSize = ref(20)
  const total = ref(0)
  const taskPulling = ref(false)
  const pullHandler = ref(null)
  const taskFilter = ref('all')
  const currentList = ref([])
  
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
  
  // 功能配置
  const functions = [
    { key: 'text_to_image', name: '文生图', category: 'image_generation', needsPrompt: true, needsImage: false, power: 20 },
    { key: 'image_to_image_portrait', name: '图生图', category: 'image_generation', needsPrompt: true, needsImage: true, power: 30 },
    { key: 'image_edit', name: '图像编辑', category: 'image_editing', needsPrompt: true, needsImage: true, multiple: true, power: 25 },
    { key: 'image_effects', name: '图像特效', category: 'image_effects', needsPrompt: false, needsImage: true, power: 15 },
    { key: 'text_to_video', name: '文生视频', category: 'video_generation', needsPrompt: true, needsImage: false, power: 100 },
    { key: 'image_to_video', name: '图生视频', category: 'video_generation', needsPrompt: true, needsImage: true, multiple: true, power: 120 },
  ]
  
  // 各功能的参数
  const textToImageParams = reactive({
    prompt: '',
    size: '1328x1328',
    scale: 2.5,
    seed: -1,
    use_pre_llm: false,
  })
  
  const imageToImageParams = reactive({
    image_input: '',
    prompt: '演唱会现场的合照，闪光灯拍摄',
    size: '1328x1328',
    gpen: 0.4,
    skin: 0.3,
    skin_unifi: 0,
    gen_mode: 'creative',
    seed: -1,
  })
  
  const imageEditParams = reactive({
    image_urls: [],
    prompt: '',
    scale: 0.5,
    seed: -1,
  })
  
  const imageEffectsParams = reactive({
    image_input1: '',
    template_id: '',
    size: '1328x1328',
  })
  
  const textToVideoParams = reactive({
    prompt: '',
    aspect_ratio: '16:9',
    seed: -1,
  })
  
  const imageToVideoParams = reactive({
    image_urls: [],
    prompt: '',
    aspect_ratio: '16:9',
    seed: -1,
  })
  
  // 计算属性
  const currentFunction = computed(() => {
    return functions.find(f => f.key === activeFunction.value) || functions[0]
  })
  
  const currentFunctions = computed(() => {
    return functions.filter(f => f.category === activeCategory.value)
  })
  
  const needsPrompt = computed(() => currentFunction.value.needsPrompt)
  const needsImage = computed(() => currentFunction.value.needsImage)
  const needsMultipleImages = computed(() => currentFunction.value.multiple)
  const currentPowerCost = computed(() => currentFunction.value.power)
  
  // 初始化方法
  const init = async () => {
    try {
      const user = await checkSession()
      isLogin.value = true
      userPower.value = user.power
      
      // 获取任务列表
      await fetchData(1)
      
      // 检查是否需要开始轮询
      const pendingCount = await getPendingCount()
      if (pendingCount > 0) {
        startPolling()
      }
    } catch (error) {
      console.error('初始化失败:', error)
    }
  }
  
  // 切换功能分类
  const switchCategory = (category) => {
    activeCategory.value = category
    const categoryFunctions = functions.filter(f => f.category === category)
    if (categoryFunctions.length > 0) {
      if (category === 'image_generation') {
        activeFunction.value = useImageInput.value ? 'image_to_image_portrait' : 'text_to_image'
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
      activeFunction.value = useImageInput.value ? 'image_to_image_portrait' : 'text_to_image'
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
    const func = functions.find(f => f.key === type)
    return func ? func.name : type
  }
  
  // 获取任务状态文本
  const getTaskStatusText = (status) => {
    const statusMap = {
      'pending': '等待中',
      'processing': '处理中',
      'completed': '已完成',
      'failed': '失败'
    }
    return statusMap[status] || status
  }
  
  // 获取状态类型
  const getStatusType = (status) => {
    const typeMap = {
      'pending': 'info',
      'processing': 'warning',
      'completed': 'success',
      'failed': 'danger'
    }
    return typeMap[status] || 'info'
  }
  
  // 切换任务筛选
  const switchTaskFilter = (filter) => {
    taskFilter.value = filter
    updateCurrentList()
  }
  
  // 更新当前列表
  const updateCurrentList = () => {
    if (taskFilter.value === 'all') {
      currentList.value = list.value
    } else if (taskFilter.value === 'image') {
      currentList.value = list.value.filter(item => 
        ['text_to_image', 'image_to_image_portrait', 'image_edit', 'image_effects'].includes(item.type)
      )
    } else if (taskFilter.value === 'video') {
      currentList.value = list.value.filter(item => 
        ['text_to_video', 'image_to_video'].includes(item.type)
      )
    }
  }
  
  // 获取任务列表
  const fetchData = async (pageNum = 1) => {
    try {
      loading.value = true
      page.value = pageNum
      
      const response = await httpGet('/api/jimeng/jobs', {
        page: pageNum,
        page_size: pageSize.value
      })
      
      if (response.data) {
        list.value = response.data.jobs || []
        total.value = response.data.total || 0
        noData.value = list.value.length === 0
        updateCurrentList()
      }
    } catch (error) {
      console.error('获取任务列表失败:', error)
      showMessageError('获取任务列表失败')
    } finally {
      loading.value = false
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
    
    try {
      submitting.value = true
      let apiUrl = ''
      let requestData = {}
      
      switch (activeFunction.value) {
        case 'text_to_image':
          apiUrl = '/api/jimeng/text-to-image'
          requestData = {
            prompt: textToImageParams.prompt,
            width: parseInt(textToImageParams.size.split('x')[0]),
            height: parseInt(textToImageParams.size.split('x')[1]),
            scale: textToImageParams.scale,
            seed: textToImageParams.seed,
            use_pre_llm: textToImageParams.use_pre_llm,
          }
          break
          
        case 'image_to_image_portrait':
          apiUrl = '/api/jimeng/image-to-image-portrait'
          requestData = {
            image_input: imageToImageParams.image_input,
            prompt: imageToImageParams.prompt,
            width: parseInt(imageToImageParams.size.split('x')[0]),
            height: parseInt(imageToImageParams.size.split('x')[1]),
            gpen: imageToImageParams.gpen,
            skin: imageToImageParams.skin,
            skin_unifi: imageToImageParams.skin_unifi,
            gen_mode: imageToImageParams.gen_mode,
            seed: imageToImageParams.seed,
          }
          break
          
        case 'image_edit':
          apiUrl = '/api/jimeng/image-edit'
          requestData = {
            image_urls: imageEditParams.image_urls,
            prompt: imageEditParams.prompt,
            scale: imageEditParams.scale,
            seed: imageEditParams.seed,
          }
          break
          
        case 'image_effects':
          apiUrl = '/api/jimeng/image-effects'
          requestData = {
            image_input1: imageEffectsParams.image_input1,
            template_id: imageEffectsParams.template_id,
            width: parseInt(imageEffectsParams.size.split('x')[0]),
            height: parseInt(imageEffectsParams.size.split('x')[1]),
          }
          break
          
        case 'text_to_video':
          apiUrl = '/api/jimeng/text-to-video'
          requestData = {
            prompt: textToVideoParams.prompt,
            aspect_ratio: textToVideoParams.aspect_ratio,
            seed: textToVideoParams.seed,
          }
          break
          
        case 'image_to_video':
          apiUrl = '/api/jimeng/image-to-video'
          requestData = {
            image_urls: imageToVideoParams.image_urls,
            prompt: imageToVideoParams.prompt,
            aspect_ratio: imageToVideoParams.aspect_ratio,
            seed: imageToVideoParams.seed,
          }
          break
      }
      
      const response = await httpPost(apiUrl, requestData)
      
      if (response.data) {
        showMessageOK('任务提交成功')
        // 重新获取任务列表
        await fetchData(1)
        // 开始轮询
        startPolling()
      }
    } catch (error) {
      console.error('提交任务失败:', error)
      showMessageError(error.message || '提交任务失败')
    } finally {
      submitting.value = false
    }
  }
  
  // 获取待处理任务数量
  const getPendingCount = async () => {
    try {
      const response = await httpGet('/api/jimeng/pending-count')
      return response.data?.count || 0
    } catch (error) {
      console.error('获取待处理任务数量失败:', error)
      return 0
    }
  }
  
  // 开始轮询
  const startPolling = () => {
    if (taskPulling.value) return
    
    taskPulling.value = true
    pullHandler.value = setInterval(async () => {
      const pendingCount = await getPendingCount()
      if (pendingCount > 0) {
        await fetchData(page.value)
      } else {
        stopPolling()
      }
    }, 3000)
  }
  
  // 停止轮询
  const stopPolling = () => {
    if (pullHandler.value) {
      clearInterval(pullHandler.value)
      pullHandler.value = null
    }
    taskPulling.value = false
  }
  
  // 重试任务
  const retryTask = async (taskId) => {
    try {
      const response = await httpPost(`/api/jimeng/retry/${taskId}`)
      if (response.data) {
        showMessageOK('重试任务已提交')
        await fetchData(page.value)
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
        await fetchData(page.value)
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
  
  // 下载文件
  const downloadFile = (item) => {
    const url = item.video_url || item.img_url
    if (url) {
      const link = document.createElement('a')
      link.href = url
      link.download = `jimeng_${item.id}.${item.video_url ? 'mp4' : 'jpg'}`
      link.click()
    }
  }
  
  // 清理
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
    list,
    noData,
    page,
    pageSize,
    total,
    taskFilter,
    currentList,
    isLogin,
    userPower,
    showDialog,
    currentVideoUrl,
    nodata,
    
    // 配置
    categories,
    functions,
    currentFunctions,
    
    // 参数
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
    getStatusType,
    switchTaskFilter,
    updateCurrentList,
    fetchData,
    submitTask,
    getPendingCount,
    startPolling,
    stopPolling,
    retryTask,
    removeJob,
    playVideo,
    downloadFile,
    cleanup,
    
    // 工具函数
    substr,
    replaceImg,
  }
})