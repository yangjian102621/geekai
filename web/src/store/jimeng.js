// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import { checkSession } from '@/store/cache'
import { JimengFunctions, JimengParams } from '@/store/data/jimeng_params'
import { useSharedStore } from '@/store/sharedata'
import { showMessageError, showMessageOK } from '@/utils/dialog'
import { httpDownload, httpGet, httpPost } from '@/utils/http'
import { replaceImg, substr } from '@/utils/libs'
import { ElMessageBox } from 'element-plus'
import { defineStore } from 'pinia'
import { reactive, ref } from 'vue'

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

  // 用户信息
  const isLogin = ref(false)
  // 视频预览
  const showDialog = ref(false)
  const currentVideoUrl = ref('')

  // 登录弹窗
  const shareStore = useSharedStore()

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
      if (!data.items || data.items.length === 0) {
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
      const user = await checkSession()
      isLogin.value = true
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
    isLogin,
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
