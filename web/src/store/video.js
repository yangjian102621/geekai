// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import nodata from '@/assets/img/no-data.png'
import { checkSession } from '@/store/cache'
import { getVideoModelByKey, getVideoModels, getVideoProviders } from '@/store/data/video_params'
import { useSharedStore } from '@/store/sharedata'
import { closeLoading, showLoading, showMessageError, showMessageOK } from '@/utils/dialog'
import { httpDownload, httpGet, httpPost } from '@/utils/http'
import { replaceImg, substr } from '@/utils/libs'
import Clipboard from 'clipboard'
import { ElMessage, ElMessageBox } from 'element-plus'
import { defineStore } from 'pinia'
import { computed, ref, watch } from 'vue'

export const useVideoStore = defineStore('video', () => {
  const providers = getVideoProviders()
  const activeProvider = ref(providers.includes('sora') ? 'sora' : providers[0] || '')

  const loading = ref(false)
  const submitting = ref(false)
  const list = ref([])
  const noData = ref(true)
  const page = ref(1)
  const pageSize = ref(10)
  const total = ref(0)
  const taskPulling = ref(true)
  const pullHandler = ref(null)
  const clipboard = ref(null)

  const showDialog = ref(false)
  const currentVideoUrl = ref('')

  const isLogin = ref(false)
  const availablePower = ref(0)
  const shareStore = useSharedStore()

  const taskFilter = ref('all') // 'all' 或 provider

  const formData = ref({})
  const requiredKeys = ref({})

  const powerConfig = ref({})
  const currentPowerCost = ref(0)

  const currentList = computed(() => {
    return list.value.filter((item) => {
      if (taskFilter.value === 'all') return true
      return item.type === taskFilter.value
    })
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
    pullHandler.value = setInterval(() => {
      if (taskPulling.value) {
        fetchData(page.value)
      }
    }, 5000)
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
      isLogin.value = true
      availablePower.value = user.power

      initClipboard()
      await loadPowerConfig()
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

  const fetchData = async (_page) => {
    if (_page) {
      page.value = _page
    }

    try {
      loading.value = true
      const res = await httpGet('/api/video/list', {
        page: page.value,
        page_size: pageSize.value,
        type: taskFilter.value === 'all' ? '' : taskFilter.value,
      })

      total.value = res.data.total
      let needPull = false
      const items = []
      for (let v of res.data.items) {
        // 检查是否需要继续轮询：progress 为 0 或 102，或者状态为 pending/in_progress/downloading
        if (v.status === 'pending' || v.status === 'in_progress' || v.status === 'downloading') {
          needPull = true
        }
        items.push({
          ...v,
          downloading: false,
        })
      }
      taskPulling.value = needPull
      list.value = items
      noData.value = list.value.length === 0
    } catch (error) {
      noData.value = true
      console.error('获取任务列表失败:', error)
    } finally {
      loading.value = false
    }
  }

  const switchProvider = (provider) => {
    activeProvider.value = provider
  }

  const switchTaskFilter = (filter) => {
    taskFilter.value = filter
    page.value = 1
    fetchData(1)
  }

  const loadPowerConfig = async () => {
    try {
      const res = await httpGet('/api/video/power-config')
      powerConfig.value = res.data || {}
    } catch (error) {
      console.error('加载算力配置失败:', error)
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

  // 根据 priceKey 获取算力值
  const getPowerByPriceKey = async (modelKey, priceKey) => {
    try {
      const res = await httpGet('/api/video/power-by-key', {
        model_key: modelKey,
        price_key: priceKey,
      })
      return res.data?.power || 0
    } catch (error) {
      console.error('获取算力失败:', error)
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

      // 调用 API 获取算力
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
    if (!isLogin.value) {
      shareStore.setShowLoginDialog(true)
      return
    }

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
    currentList,
    noData,
    page,
    pageSize,
    total,
    showDialog,
    currentVideoUrl,
    isLogin,
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
