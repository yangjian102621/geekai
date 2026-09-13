// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import { checkSession } from '@/store/cache'
import { JimengFunctions, JimengParams } from '@/store/data/jimeng_params'
import { showMessageError, showMessageOK } from '@/utils/dialog'
import { httpDownload, httpGet, httpPost } from '@/utils/http'
import failedIcon from '@/assets/img/failed.png'
import loadingGif from '@/assets/img/loading.gif'
import { getThumbURL, replaceImg, substr } from '@/utils/libs'
import { ElMessageBox } from 'element-plus'
import { defineStore } from 'pinia'
import { nextTick, reactive, ref } from 'vue'

/** 与瀑布流 imgSelector: img_thumb 对齐（参考 Image.vue） */
function normalizeJimengWaterfallThumb(item) {
  const img = item.img_url ? replaceImg(item.img_url) : ''
  if (img) {
    item.img_thumb = getThumbURL(img, 300, 0)
    return
  }
  if (item.video_url) {
    item.img_thumb = loadingGif
    return
  }
  if (item.status === 'failed') {
    item.img_thumb = failedIcon
    return
  }
  item.img_thumb = loadingGif
}

function refreshJimengListThumbs(list) {
  if (!list || !list.length) {
    return
  }
  for (const row of list) {
    normalizeJimengWaterfallThumb(row)
  }
}

/** 轮询需跟进的非终态（含 submited：队列尚未改 in_queue 前） */
const POLLING_ACTIVE_STATUSES = new Set(['submited', 'in_queue', 'generating', 'done'])

export const useJimengStore = defineStore('jimeng', () => {
  // 共同状态
  const loading = ref(false)
  const submitting = ref(false)
  const page = ref(1)
  const pageSize = ref(10)
  const total = ref(0)
  const taskFilter = ref('all')
  const currentList = ref([])
  const isOver = ref(false)

  // 视频预览
  const showDialog = ref(false)
  const currentVideoUrl = ref('')

  // 积分消耗配置
  const powerConfig = reactive({ powers: {} })
  const currentPowerCost = ref('0积分')

  // 功能配置
  const functions = JimengFunctions
  // 当前激活的功能
  const activeFunction = ref('image')
  // 参数配置
  const functionParams = JimengParams
  // 表单数据
  const formData = ref({})
  // 必填参数
  const requiredKeys = ref({})
  // 进度
  const progress = ref({
    image: 100,
    video: 100,
    virtualHuman: 38,
    actionTransfer: 65,
  })
  // 切换功能
  const switchFunction = (f) => {
    activeFunction.value = f.key
    setFunctionPowers()
  }

  // 获取功能名称
  const getFunctionName = (type) => {
    const func = functions.find((f) => f.key === type)
    return func ? func.name : type
  }

  // 获取任务状态文本
  const getTaskStatusText = (status) => {
    const statusMap = {
      submited: '任务已提交',
      in_queue: '任务排队中',
      generating: '任务执行中',
      done: '处理完成',
      success: '任务成功',
      failed: '任务失败',
      canceled: '任务已取消',
    }
    return statusMap[status] || status
  }

  // 获取状态类型
  const getTaskType = (type) => {
    const typeMap = {
      image: 'info',
      video: 'primary',
      virtual_human: 'success',
      action_transfer: 'warning',
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
      if (!data.items || data.items.length < pageSize.value) {
        isOver.value = true
      }
      if (pageNum === 1) {
        currentList.value = data.items
      } else {
        currentList.value = currentList.value.concat(data.items)
      }
      refreshJimengListThumbs(currentList.value)
    } catch (error) {
      showMessageError('获取任务列表失败:' + error.message)
    } finally {
      loading.value = false
    }
  }

  /** 轮询合并：新行对象 + 新数组引用，避免瀑布流对同 id 行就地 mutate 不刷新 */
  const mergeJobsIntoCurrentList = (rows) => {
    if (!rows || !rows.length) {
      return
    }
    const next = currentList.value.map((item) => {
      const hit = rows.find((i) => i.id === item.id)
      if (hit) {
        const row = { ...item, ...hit }
        normalizeJimengWaterfallThumb(row)
        return row
      }
      return item
    })
    currentList.value = next
  }

  // 简单轮询逻辑
  const startPolling = () => {
    if (pollHandler) {
      clearInterval(pollHandler)
    }
    pollHandler = setInterval(async () => {
      try {
        const activeIds = currentList.value
          .filter((item) => POLLING_ACTIVE_STATUSES.has(item.status))
          .map((item) => item.id)

        const body = {
          page: 1,
          page_size: 100,
          filter: taskFilter.value,
        }
        if (activeIds.length > 0) {
          body.ids = activeIds
        }

        const response = await httpPost('/api/jimeng/jobs', body)
        const data = response.data || {}
        if (!data.items || data.items.length === 0) {
          if (activeIds.length === 0) {
            stopPolling()
          }
          return
        }

        mergeJobsIntoCurrentList(data.items)

        const stillActive = currentList.value.some((item) =>
          POLLING_ACTIVE_STATUSES.has(item.status)
        )

        if (!stillActive) {
          stopPolling()
        }
      } catch (e) {
        console.error('jimeng poll error', e)
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
    for (const key in requiredKeys.value) {
      if (!formData.value[key]) {
        showMessageError('缺少参数：' + requiredKeys.value[key].label)
        return
      }
    }

    try {
      submitting.value = true
      formData.value.type = activeFunction.value
      // 视频 duration 转成整数
      if (formData.value.duration) {
        formData.value.duration = parseInt(formData.value.duration)
      }

      const data = { ...formData.value }

      if (data.image_urls && !Array.isArray(data.image_urls)) {
        data.image_urls = [data.image_urls]
      }
      if (data.video_url && !Array.isArray(data.video_url)) {
        data.video_url = [data.video_url]
      }
      if (data.audio_url && !Array.isArray(data.audio_url)) {
        data.audio_url = [data.audio_url]
      }

      if (typeof data.req_key === 'string' && data.req_key.startsWith('doubao-seedance-')) {
        data.content = buildSeedanceContent(data)
        if (data.seedance_mode === 'multimodal' && data.content.length === 0) {
          throw new Error('多模态模式下，请至少上传图片、视频或音频中的一种素材')
        }
        // Seedance 统一使用 content[] 传多模态，避免旧字段类型与后端绑定冲突。
        delete data.video_url
        delete data.audio_url
        delete data.image_role
        delete data.seedance_mode
      }

      const response = await httpPost('/api/jimeng/task', data)
      showMessageOK('任务提交成功')
      isOver.value = false
      await fetchData(1)
      startPolling()
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
        isOver.value = false
        await fetchData(1)
      }
    } catch (error) {
      if (error !== 'cancel') {
        console.error('删除任务失败:', error)
        showMessageError(error.message || '删除任务失败')
      }
    }
  }

  const setFunctionPowers = () => {
    nextTick(() => {
      const key = formData.value.req_key
      const perUnit = key ? powerConfig.powers[key] : 0
      if (!perUnit) {
        currentPowerCost.value = '未配置积分'
        return
      }
      currentPowerCost.value =
        activeFunction.value === 'image' ? `${perUnit}积分/张` : `${perUnit}积分/秒`
    })
  }

  const normalizeMediaList = (value) => {
    if (!value) {
      return []
    }
    if (Array.isArray(value)) {
      return value.filter(Boolean)
    }
    return [value]
  }

  const buildSeedanceContent = (data) => {
    const content = []
    if (data.prompt) {
      content.push({
        type: 'text',
        text: data.prompt,
      })
    }

    const imageList = normalizeMediaList(data.image_urls)
    const videoList = normalizeMediaList(data.video_url)
    const audioList = normalizeMediaList(data.audio_url)
    const imageRole = data.image_role || 'reference_image'

    if (data.seedance_mode === 'image_first') {
      imageList.slice(0, 1).forEach((url) => {
        content.push({
          type: 'image_url',
          image_url: { url },
          role: 'first_frame',
        })
      })
    } else if (data.seedance_mode === 'image_first_last') {
      imageList.slice(0, 2).forEach((url, index) => {
        content.push({
          type: 'image_url',
          image_url: { url },
          role: index === 0 ? 'first_frame' : 'last_frame',
        })
      })
    } else {
      imageList.forEach((url) => {
        content.push({
          type: 'image_url',
          image_url: { url },
          role: imageRole,
        })
      })
    }

    videoList.forEach((url) => {
      content.push({
        type: 'video_url',
        video_url: { url },
        role: 'reference_video',
      })
    })
    audioList.forEach((url) => {
      content.push({
        type: 'audio_url',
        audio_url: { url },
        role: 'reference_audio',
      })
    })
    return content
  }

  watch(
    () => formData.value,
    () => {
      setFunctionPowers()
    }
  )

  // 初始化方法
  const init = async () => {
    try {
      // 获取积分消耗配置
      const powerRes = await httpGet('/api/jimeng/power-config')
      if (powerRes.data) {
        powerConfig.powers = powerRes.data.powers || {}
        setFunctionPowers()
      }
      await checkSession()
      // 获取任务列表
      await fetchData(1)
      // 开始轮询
      startPolling()
    } catch (error) {
      console.error('初始化失败:', error)
    }
  }

  // 页面卸载时清理轮询
  const cleanup = () => {
    page.value = 1
    pageSize.value = 10
    total.value = 0
    taskFilter.value = 'all'
    currentList.value = []
    isOver.value = false
    loading.value = false
    stopPolling()
  }

  // 返回所有状态和方法
  return {
    // 状态
    activeFunction,
    loading,
    submitting,
    page,
    pageSize,
    total,
    taskFilter,
    currentList,
    isOver,
    showDialog,
    currentVideoUrl,
    // 配置
    functions,
    activeFunction,
    functionParams,
    formData,
    requiredKeys,
    progress,
    currentPowerCost,

    // 方法
    init,
    switchFunction,
    getFunctionName,
    getTaskStatusText,
    getTaskType,
    switchTaskFilter,
    setFunctionPowers,
    fetchData,
    submitTask,
    downloadFile,
    retryTask,
    removeJob,
    cleanup,

    // 工具函数
    substr,
    replaceImg,
  }
})
