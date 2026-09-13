// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import nodata from '@/assets/img/no-data.png'
import failedIcon from '@/assets/img/failed.png'
import loadingGif from '@/assets/img/loading.gif'
import { checkSession } from '@/store/cache'
import { getVideoModelByKey, getVideoModels, getVideoProviders } from '@/store/data/video_params'
import { closeLoading, showLoading, showMessageError, showMessageOK } from '@/utils/dialog'
import { httpDownload, httpGet, httpPost } from '@/utils/http'
import { getThumbURL, replaceImg, substr } from '@/utils/libs'
import Clipboard from 'clipboard'
import { ElMessage, ElMessageBox } from 'element-plus'
import { defineStore } from 'pinia'
import { computed, ref, watch } from 'vue'

export const useVideoStore = defineStore('video', () => {
  const normalizeWorkItemThumb = (item) => {
    if (!item) {
      return item
    }
    if (item.status === 'success' && item.video_url) {
      item.img_thumb = getThumbURL(replaceImg(item.video_url), 300, 0)
      return item
    }
    if (item.status === 'failed') {
      item.img_thumb = failedIcon
      return item
    }
    if (item.status === 'downloading') {
      item.img_thumb = loadingGif
      return item
    }
    item.img_thumb = loadingGif
    return item
  }

  const normalizeWorkListThumbs = (rows) => {
    if (!rows || !rows.length) {
      return rows
    }
    for (const row of rows) {
      normalizeWorkItemThumb(row)
    }
    return rows
  }

  const providers = getVideoProviders()
  const activeProvider = ref(providers.includes('sora') ? 'sora' : providers[0] || '')

  const loading = ref(false)
  const submitting = ref(false)
  const taskList = ref([])
  const list = ref([])
  const noData = ref(true)
  const page = ref(1)
  const pageSize = ref(10)
  const total = ref(0)
  const isOver = ref(false)
  const taskPulling = ref(true)
  const pullHandler = ref(null)
  const clipboard = ref(null)

  const showDialog = ref(false)
  const currentVideoUrl = ref('')

  const availablePower = ref(0)

  const taskFilter = ref('all') // 'all' 或 provider

  const formData = ref({})
  const requiredKeys = ref({})

  const powerConfig = ref({})
  const currentPowerCost = ref(0)

  const currentList = computed(() => {
    const filtered = list.value.filter((item) => {
      if (taskFilter.value === 'all') return true
      return item.type === taskFilter.value
    })
    return filtered
  })

  const providerModels = computed(() => {
    if (!activeProvider.value) return []
    return getVideoModels(activeProvider.value)
  })

  const initClipboard = () => {
    if (clipboard.value) {
      clipboard.value.destroy()
    }
    clipboard.value = new Clipboard('.copy-prompt')
    clipboard.value.on('success', () => ElMessage.success('复制成功！'))
    clipboard.value.on('error', () => ElMessage.error('复制失败！'))
  }

  const startPolling = () => {
    if (pullHandler.value) {
      clearInterval(pullHandler.value)
    }
    pollLatest()
    pullHandler.value = setInterval(() => {
        pollLatest()
    }, 2000)
  }

  const stopPolling = () => {
    if (pullHandler.value) {
      clearInterval(pullHandler.value)
      pullHandler.value = null
    }
  }

  const init = async () => {
    try {
      const user = await checkSession()
      availablePower.value = user.power

      initClipboard()
      await loadPowerConfig()
      await fetchTaskList()
      await fetchData(1)
      startPolling()
    } catch (error) {
      console.error('初始化失败:', error)
    }
  }

  const cleanup = () => {
    if (clipboard.value) {
      clipboard.value.destroy()
    }
    stopPolling()
  }

  const fetchData = async (pageNum = 1) => {
    try {
      loading.value = true
      page.value = pageNum
      if (pageNum === 1) {
        isOver.value = false
      }

      const res = await httpGet('/api/video/works', {
        page: pageNum,
        page_size: pageSize.value,
        type: taskFilter.value === 'all' ? '' : taskFilter.value,
      })
  
      total.value = res.data.total
      const items = (res.data.items || []).map((v) => ({
        ...v,
        downloading: false,
      }))
      normalizeWorkListThumbs(items)
  
      if (items.length === 0) {
        isOver.value = true
        if (pageNum === 1) {
          list.value = []
          noData.value = true
        }
        taskPulling.value = false
        return
      }

      if (items.length < pageSize.value || pageNum * pageSize.value >= total.value) {
        isOver.value = true
      }

      const needPull = items.some((item) => item.status === 'downloading')
      taskPulling.value = needPull
  
      if (pageNum === 1) {
        list.value = items
      } else {
        const merged = [...list.value]
        const exists = new Set(merged.map((it) => it.id))
        for (const item of items) {
          if (!exists.has(item.id)) {
            merged.push(item)
          }
        }
        list.value = merged
      }
      noData.value = list.value.length === 0
    } catch (error) {
      if (pageNum === 1) {
        noData.value = true
      }
      console.error('获取任务列表失败:', error)
    } finally {
      loading.value = false
    }
  }

  const pollLatest = async () => {
    try {
      const [taskRes, workRes] = await Promise.all([
        httpGet('/api/video/tasks', {
          type: taskFilter.value === 'all' ? '' : taskFilter.value,
        }),
        httpGet('/api/video/works', {
          page: 1,
          page_size: pageSize.value,
          type: taskFilter.value === 'all' ? '' : taskFilter.value,
        }),
      ])
      const nextTasks = (taskRes.data || []).map((v) => ({ ...v, downloading: false }))
      taskList.value = nextTasks
        const res = workRes
      const items = (res.data.items || []).map((v) => ({ ...v, downloading: false }))
      normalizeWorkListThumbs(items)
        let needPull = false
      const latestMap = new Map(items.map((it) => [it.id, it]))
      const matchedIds = []
      const missedIds = []
      for (const fetchedItem of items) {
        if (list.value.some((localItem) => localItem.id === fetchedItem.id)) {
          matchedIds.push(fetchedItem.id)
        } else {
          missedIds.push(fetchedItem.id)
        }
      }
        list.value = list.value.map((row) => {
        const latest = latestMap.get(row.id)
        if (!latest) {
          return row
        }
        if (latest.status === 'downloading') {
          needPull = true
        }
        return normalizeWorkItemThumb({ ...row, ...latest })
      })
      for (const row of items) {
        if (!list.value.some((item) => item.id === row.id)) {
          list.value.unshift(row)
        }
      }
      normalizeWorkListThumbs(list.value)
      taskPulling.value = needPull || nextTasks.length > 0
      } catch (error) {
        console.error('轮询视频任务失败:', error)
    }
  }

  const fetchTaskList = async () => {
    try {
      const res = await httpGet('/api/video/tasks', {
        type: taskFilter.value === 'all' ? '' : taskFilter.value,
      })
      const items = (res.data || []).map((v) => ({ ...v, downloading: false }))
      taskList.value = items
      } catch (error) {
      console.error('获取任务列表失败:', error)
    }
  }

  const switchProvider = (provider) => {
    activeProvider.value = provider
  }

  const switchTaskFilter = (filter) => {
    taskFilter.value = filter
    page.value = 1
    isOver.value = false
    fetchTaskList()
    fetchData(1)
  }

  const loadPowerConfig = async () => {
    try {
      const res = await httpGet('/api/video/power-config')
      powerConfig.value = res.data || {}
    } catch (error) {
      console.error('加载积分配置失败:', error)
      powerConfig.value = {}
    }
  }

  // 格式化参数值为字符串（用于生成 priceKey）
  const formatParamValue = (paramName, value) => {
    if (paramName === 'sound') {
      return value === true || value === 'true' || value === 1 ? 'sound' : 'silent'
    }
    if (typeof value === 'boolean') {
      return value ? 'true' : 'false'
    }
    return String(value)
  }

  // 根据模型配置和表单数据生成价格 key
  const generatePriceKey = (model, formData) => {
    if (!model.priceParams || !Array.isArray(model.priceParams)) {
      return 'fixed'
    }

    // 如果是固定价格
    if (model.priceParams.length === 1 && model.priceParams[0] === 'fixed') {
      return 'fixed'
    }

    // 从 formData 中提取 priceParams 指定的参数值
    const values = model.priceParams.map((paramName) => {
      let value = formData[paramName]

      // 如果值缺失，尝试从模型参数配置中获取默认值
      if (value === undefined || value === null || value === '') {
        const param = model.params?.find((p) => p.name === paramName)
        if (param) {
          if (param.type === 'select' && param.options && param.options.length > 0) {
            value = param.value || param.options[0].value
          } else {
            value = param.value
          }
        }
      }

      // 如果仍然没有值，使用空字符串（这种情况应该很少）
      if (value === undefined || value === null) {
        value = ''
      }

      return formatParamValue(paramName, value)
    })

    // 用下划线连接生成 priceKey
    return values.join('_')
  }

  // 防抖定时器
  let powerDebounceTimer = null

  // 根据 priceKey 获取积分值
  const getPowerByPriceKey = async (modelKey, priceKey) => {
    try {
      const res = await httpGet('/api/video/power-by-key', {
        model_key: modelKey,
        price_key: priceKey,
      })
      return res.data?.power || 0
    } catch (error) {
      console.error('获取积分失败:', error)
      return 0
    }
  }

  const setCurrentPowerCost = async () => {
    // 清除之前的定时器
    if (powerDebounceTimer) {
      clearTimeout(powerDebounceTimer)
    }

    // 设置新的定时器（防抖）
    powerDebounceTimer = setTimeout(async () => {
      const modelKey = formData.value?.req_key
      if (!modelKey) {
        currentPowerCost.value = 0
        return
      }

      const model = getVideoModelByKey(modelKey)
      if (!model) {
        currentPowerCost.value = 0
        return
      }

      const priceKey = generatePriceKey(model, formData.value)
      if (!priceKey) {
        currentPowerCost.value = 0
        return
      }

      // 调用 API 获取积分
      const power = await getPowerByPriceKey(modelKey, priceKey)
      currentPowerCost.value = power
    }, 300) // 300ms 防抖
  }

  watch(
    () => formData.value,
    () => setCurrentPowerCost(),
    { deep: true }
  )

  const isEmptyValue = (v) => {
    if (v === undefined || v === null) return true
    if (typeof v === 'string' && v.trim() === '') return true
    if (Array.isArray(v) && v.length === 0) return true
    return false
  }

  const normalizeImageParam = (paramName, value) => {
    if (!value) return value
    if (Array.isArray(value)) {
      return value.filter(Boolean).map((u) => replaceImg(u))
    }
    if (typeof value === 'string') {
      if (paramName === 'images') return [replaceImg(value)]
      return replaceImg(value)
    }
    return value
  }

  const createVideoTask = async () => {
    const modelKey = formData.value?.req_key
    if (!modelKey) {
      return ElMessage.error('请选择模型')
    }

    const model = getVideoModelByKey(modelKey)
    if (!model) {
      return ElMessage.error('模型配置不存在')
    }

    for (const key in requiredKeys.value) {
      if (isEmptyValue(formData.value?.[key])) {
        return showMessageError('缺少参数：' + requiredKeys.value[key].label)
      }
    }

    const raw = { ...(formData.value || {}) }
    delete raw.req_key
    delete raw.action

    const prompt = raw.prompt
    if (!prompt || !prompt.trim()) {
      return ElMessage.error('请输入视频描述')
    }

    // prompt 作为顶层字段提交
    delete raw.prompt

    // 图片字段归一化（string/array + replaceImg）
    if (Array.isArray(model.params)) {
      model.params.forEach((p) => {
        if (p.type === 'image') {
          raw[p.name] = normalizeImageParam(p.name, raw[p.name])
        }
      })
    }

    // 生成 priceKey
    const priceKey = generatePriceKey(model, formData.value)

    const requestData = {
      provider: model.provider,
      model: model.key,
      prompt,
      params: raw,
      price_key: priceKey,
    }

    try {
      submitting.value = true
      showLoading('创建任务中...')
      await httpPost('/api/video/create', requestData)
      showMessageOK('任务创建成功')
      closeLoading()
      isOver.value = false
      await fetchTaskList()
      await fetchData(1)
      taskPulling.value = true
    } catch (error) {
      closeLoading()
      showMessageError('创建任务失败：' + error.message)
    } finally {
      submitting.value = false
    }
  }

  const playVideo = (item) => {
    currentVideoUrl.value = replaceImg(item.video_url)
    showDialog.value = true
  }

  const downloadVideo = async (item) => {
    const url = replaceImg(item.video_url)
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
    } catch (error) {
      showMessageError('下载失败')
    } finally {
      item.downloading = false
    }
  }

  const removeJob = async (item) => {
    try {
      await ElMessageBox.confirm('此操作将会删除任务相关文件，继续操作吗?', '删除提示', {
        confirmButtonText: '确认',
        cancelButtonText: '取消',
        type: 'warning',
      })

      await httpGet('/api/video/remove', { id: item.id })
      ElMessage.success('任务删除成功')
      isOver.value = false
      await fetchTaskList()
      await fetchData(1)
    } catch (error) {
      if (error !== 'cancel') {
        ElMessage.error('任务删除失败：' + error.message)
      }
    }
  }

  return {
    activeProvider,
    providerModels,
    providers,

    loading,
    submitting,
    list,
    taskList,
    currentList,
    noData,
    page,
    pageSize,
    total,
    isOver,
    showDialog,
    currentVideoUrl,
    availablePower,
    nodata,
    taskFilter,

    formData,
    requiredKeys,
    powerConfig,
    currentPowerCost,

    init,
    cleanup,
    fetchData,
    fetchTaskList,
    switchProvider,
    switchTaskFilter,

    loadPowerConfig,
    createVideoTask,

    playVideo,
    downloadVideo,
    removeJob,
    substr,
    replaceImg,
  }
})
