// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import nodata from '@/assets/img/no-data.png'
import { checkSession, getSystemInfo } from '@/store/cache'
import { closeLoading, showLoading, showMessageError, showMessageOK } from '@/utils/dialog'
import { httpDownload, httpGet, httpPost } from '@/utils/http'
import { replaceImg, substr } from '@/utils/libs'
import Clipboard from 'clipboard'
import { ElMessage, ElMessageBox } from 'element-plus'
import { defineStore } from 'pinia'
import { computed, reactive, ref } from 'vue'

export const useVideoStore = defineStore('video', () => {
  // 当前活跃的视频类型
  const activeVideoType = ref('luma')

  // 共同状态
  const loading = ref(false)
  const list = ref([])
  const noData = ref(true)
  const page = ref(1)
  const pageSize = ref(10)
  const total = ref(0)
  const taskPulling = ref(true)
  const pullHandler = ref(null)
  const clipboard = ref(null)

  // 视频预览
  const showDialog = ref(false)
  const currentVideoUrl = ref('')

  // 用户信息
  const isLogin = ref(false)
  const availablePower = ref(100)

  // 任务筛选
  const taskFilter = ref('all') // 'all', 'luma', 'keling'

  // Luma 相关状态
  const lumaUseImageMode = ref(false) // 是否使用图片辅助生成
  const lumaParams = reactive({
    prompt: '',
    expand_prompt: false,
    loop: false,
    image: '', // 起始帧
    image_tail: '', // 结束帧
  })

  // KeLing 相关状态
  const isGenerating = ref(false)
  const generating = ref(false)
  const kelingPowerCost = ref(10)
  const lumaPowerCost = ref(10)
  const showCameraControl = ref(false)
  const keLingPowers = ref({})

  const models = ref([
    { text: '可灵 1.6', value: 'kling-v1-6' },
    { text: '可灵 1.5', value: 'kling-v1-5' },
    { text: '可灵 1.0', value: 'kling-v1' },
  ])

  const rates = [
    { css: 'square', value: '1:1', text: '1:1', img: '/images/mj/rate_1_1.png' },
    { css: 'size16-9', value: '16:9', text: '16:9', img: '/images/mj/rate_16_9.png' },
    { css: 'size9-16', value: '9:16', text: '9:16', img: '/images/mj/rate_9_16.png' },
  ]

  // KeLing 相关状态
  const kelingUseImageMode = ref(false) // 是否使用图片辅助生成
  const kelingParams = reactive({
    model: 'kling-v1-6',
    prompt: '',
    negative_prompt: '',
    cfg_scale: 0.7,
    mode: 'std',
    aspect_ratio: '16:9',
    duration: '5',
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
    image: '',
    image_tail: '',
  })

  // 计算属性
  const currentList = computed(() => {
    return list.value.filter((item) => {
      if (taskFilter.value === 'all') {
        return true
      } else if (taskFilter.value === 'luma') {
        return item.type === 'luma' || !item.type // 兼容旧数据
      } else if (taskFilter.value === 'keling') {
        return item.type === 'keling'
      }
      return true
    })
  })

  // 初始化方法
  const init = async () => {
    try {
      const user = await checkSession()
      isLogin.value = true
      availablePower.value = user.power

      // 初始化剪贴板
      if (clipboard.value) {
        clipboard.value.destroy()
      }
      clipboard.value = new Clipboard('.copy-prompt')
      clipboard.value.on('success', () => {
        ElMessage.success('复制成功！')
      })
      clipboard.value.on('error', () => {
        ElMessage.error('复制失败！')
      })

      // 获取系统信息
      const sysInfo = await getSystemInfo()
      lumaPowerCost.value = sysInfo.data.luma_power
      keLingPowers.value = sysInfo.data.keling_powers
      updateModelPower()

      // 获取数据并开始轮询
      await fetchData(1)
      startPolling()
    } catch (error) {
      console.error('初始化失败:', error)
    }
  }

  // 清理方法
  const cleanup = () => {
    if (clipboard.value) {
      clipboard.value.destroy()
    }
    stopPolling()
  }

  // 开始轮询
  const startPolling = () => {
    if (pullHandler.value) {
      clearInterval(pullHandler.value)
    }
    pullHandler.value = setInterval(() => {
      if (taskPulling.value) {
        fetchData(page.value)
      }
    }, 5000)
  }

  // 停止轮询
  const stopPolling = () => {
    if (pullHandler.value) {
      clearInterval(pullHandler.value)
      pullHandler.value = null
    }
  }

  // 获取任务列表
  const fetchData = async (_page) => {
    if (_page) {
      page.value = _page
    }

    try {
      const res = await httpGet('/api/video/list', {
        page: page.value,
        page_size: pageSize.value,
        type: taskFilter.value === 'all' ? '' : taskFilter.value,
      })

      total.value = res.data.total
      let needPull = false
      const items = []

      for (let v of res.data.items) {
        if (v.progress === 0 || v.progress === 102) {
          needPull = true
        }
        items.push({
          ...v,
          downloading: false,
        })
      }

      loading.value = false
      taskPulling.value = needPull

      if (JSON.stringify(list.value) !== JSON.stringify(items)) {
        list.value = items
      }
      noData.value = list.value.length === 0
    } catch (error) {
      loading.value = false
      noData.value = true
      console.error('获取任务列表失败:', error)
    }
  }

  // Luma 相关方法
  const uploadLumaStartImage = async (file) => {
    const formData = new FormData()
    formData.append('file', file.file)

    try {
      showLoading('图片上传中...')
      const res = await httpPost('/api/upload', formData)
      lumaParams.image = res.data.url
      ElMessage.success('上传成功')
      closeLoading()
    } catch (error) {
      showMessageError('上传失败: ' + error.message)
      closeLoading()
    }
  }

  const uploadLumaEndImage = async (file) => {
    const formData = new FormData()
    formData.append('file', file.file)

    try {
      showLoading('图片上传中...')
      const res = await httpPost('/api/upload', formData)
      lumaParams.image_tail = res.data.url
      ElMessage.success('上传成功')
    } catch (error) {
      showMessageError('上传失败: ' + error.message)
    } finally {
      closeLoading()
    }
  }

  const removeLumaImage = (type) => {
    if (type === 'start') {
      lumaParams.image = ''
    } else if (type === 'end') {
      lumaParams.image_tail = ''
    }
  }

  const switchLumaImages = () => {
    ;[lumaParams.image, lumaParams.image_tail] = [lumaParams.image_tail, lumaParams.image]
  }

  const toggleLumaImageMode = (enabled) => {
    lumaUseImageMode.value = enabled
    // 关闭时清空图片
    if (!enabled) {
      lumaParams.image = ''
      lumaParams.image_tail = ''
    }
  }

  const createLumaVideo = async () => {
    if (!lumaParams.prompt?.trim()) {
      return ElMessage.error('请输入视频描述')
    }

    if (lumaUseImageMode.value && !lumaParams.image) {
      return ElMessage.error('请上传起始帧图片')
    }

    // 处理参数
    const requestData = {
      ...lumaParams,
      task_type: lumaUseImageMode.value ? 'image2video' : 'text2video',
    }

    // 处理图片链接
    if (requestData.image) {
      requestData.first_frame_img = replaceImg(requestData.image)
    }
    if (requestData.image_tail) {
      requestData.end_frame_img = replaceImg(requestData.image_tail)
    }

    try {
      await httpPost('/api/video/luma/create', requestData)
      await fetchData(1)
      taskPulling.value = true
      showMessageOK('创建任务成功')
    } catch (error) {
      showMessageError('创建任务失败：' + error.message)
    }
  }

  // KeLing 相关方法
  const changeRate = (item) => {
    kelingParams.aspect_ratio = item.value
  }

  const updateModelPower = () => {
    showCameraControl.value = kelingParams.model === 'kling-v1-5' && kelingParams.mode === 'pro'
    kelingPowerCost.value =
      keLingPowers.value[`${kelingParams.model}_${kelingParams.mode}_${kelingParams.duration}`] ||
      10
  }

  const toggleKelingImageMode = (enabled) => {
    kelingUseImageMode.value = enabled
    // 关闭时清空图片
    if (!enabled) {
      kelingParams.image = ''
      kelingParams.image_tail = ''
    }
  }

  const uploadKelingStartImage = async (file) => {
    const formData = new FormData()
    formData.append('file', file.file)

    try {
      showLoading('图片上传中...')
      const res = await httpPost('/api/upload', formData)
      kelingParams.image = res.data.url
      ElMessage.success('上传成功')
      closeLoading()
    } catch (error) {
      showMessageError('上传失败: ' + error.message)
      closeLoading()
    }
  }

  const uploadKelingEndImage = async (file) => {
    const formData = new FormData()
    formData.append('file', file.file)

    try {
      showLoading('图片上传中...')
      const res = await httpPost('/api/upload', formData)
      kelingParams.image_tail = res.data.url
      ElMessage.success('上传成功')
    } catch (error) {
      showMessageError('上传失败: ' + error.message)
    } finally {
      closeLoading()
    }
  }

  const removeKelingImage = (type) => {
    if (type === 'start') {
      kelingParams.image = ''
    } else if (type === 'end') {
      kelingParams.image_tail = ''
    }
  }

  const switchKelingImages = () => {
    ;[kelingParams.image, kelingParams.image_tail] = [kelingParams.image_tail, kelingParams.image]
  }

  const createKelingVideo = async () => {
    if (generating.value) return

    if (!kelingParams.prompt?.trim()) {
      return ElMessage.error('请输入视频描述')
    }

    if (kelingParams.prompt.length > 500) {
      return ElMessage.error('视频描述不能超过 500 个字符')
    }

    if (kelingUseImageMode.value && !kelingParams.image) {
      return ElMessage.error('请上传起始帧图片')
    }

    generating.value = true

    // 处理参数
    const requestData = {
      ...kelingParams,
      task_type: kelingUseImageMode.value ? 'image2video' : 'text2video',
    }

    // 处理图片链接
    if (requestData.image) {
      requestData.image = replaceImg(requestData.image)
    }
    if (requestData.image_tail) {
      requestData.image_tail = replaceImg(requestData.image_tail)
    }

    try {
      await httpPost('/api/video/keling/create', requestData)
      showMessageOK('任务创建成功')

      // 新增重置
      page.value = 1
      list.value.unshift({
        progress: 0,
        prompt: requestData.prompt,
        raw_data: {
          task_type: requestData.task_type,
          model: requestData.model,
          duration: requestData.duration,
          mode: requestData.mode,
        },
      })
      taskPulling.value = true
    } catch (error) {
      showMessageError('创建失败: ' + error.message)
    } finally {
      generating.value = false
    }
  }

  // 提示词生成
  const generatePrompt = async () => {
    if (isGenerating.value) return

    const prompt = activeVideoType.value === 'luma' ? lumaParams.prompt : kelingParams.prompt
    if (!prompt) {
      return showMessageError('请输入原始提示词')
    }

    isGenerating.value = true
    showLoading('正在生成视频脚本...')

    try {
      const res = await httpPost('/api/prompt/video', { prompt })
      if (activeVideoType.value === 'luma') {
        lumaParams.prompt = res.data
      } else {
        kelingParams.prompt = res.data
      }
      closeLoading()
    } catch (error) {
      showMessageError('生成提示词失败：' + error.message)
      closeLoading()
    } finally {
      isGenerating.value = false
    }
  }

  // 视频预览
  const playVideo = (item) => {
    currentVideoUrl.value = replaceImg(item.video_url)
    showDialog.value = true
  }

  // 视频下载
  const downloadVideo = async (item) => {
    const url = replaceImg(item.video_url)
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

  // 删除任务
  const removeJob = async (item) => {
    try {
      await ElMessageBox.confirm('此操作将会删除任务相关文件，继续操作码?', '删除提示', {
        confirmButtonText: '确认',
        cancelButtonText: '取消',
        type: 'warning',
      })

      await httpGet('/api/video/remove', { id: item.id })
      ElMessage.success('任务删除成功')
      await fetchData()
    } catch (error) {
      if (error !== 'cancel') {
        ElMessage.error('任务删除失败：' + error.message)
      }
    }
  }

  // 发布任务
  const publishJob = async (item) => {
    try {
      await httpGet('/api/video/publish', { id: item.id, publish: item.publish })
      ElMessage.success('操作成功')
    } catch (error) {
      ElMessage.error('操作失败：' + error.message)
    }
  }

  // 切换视频类型
  const switchVideoType = (type) => {
    activeVideoType.value = type
  }

  // 切换任务筛选
  const switchTaskFilter = (filter) => {
    taskFilter.value = filter
    page.value = 1
    fetchData(1)
  }

  return {
    // 状态
    activeVideoType,
    loading,
    list,
    currentList,
    noData,
    page,
    pageSize,
    total,
    taskPulling,
    showDialog,
    currentVideoUrl,
    isLogin,
    availablePower,
    nodata,
    taskFilter,

    // Luma 状态
    lumaUseImageMode,
    lumaParams,
    lumaPowerCost,
    // KeLing 状态
    kelingUseImageMode,
    isGenerating,
    generating,
    kelingPowerCost,
    showCameraControl,
    keLingPowers,
    models,
    rates,
    kelingParams,

    // 方法
    init,
    cleanup,
    fetchData,
    switchVideoType,
    switchTaskFilter,

    // Luma 方法
    toggleLumaImageMode,
    uploadLumaStartImage,
    uploadLumaEndImage,
    removeLumaImage,
    switchLumaImages,
    createLumaVideo,

    // KeLing 方法
    toggleKelingImageMode,
    changeRate,
    updateModelPower,
    uploadKelingStartImage,
    uploadKelingEndImage,
    removeKelingImage,
    switchKelingImages,
    createKelingVideo,

    // 共同方法
    generatePrompt,
    playVideo,
    downloadVideo,
    removeJob,
    publishJob,
    substr,
    replaceImg,
  }
})
