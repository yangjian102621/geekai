import { closeLoading, showLoading, showMessageError, showMessageOK } from '@/utils/dialog'
import { httpGet, httpPost } from '@/utils/http'
import { defineStore } from 'pinia'
import { showConfirmDialog } from 'vant'
import { computed, ref } from 'vue'

export const useJimengStore = defineStore('mobile-jimeng', () => {
  // 响应式数据
  const activeCategory = ref('image_generation')
  const useImageInput = ref(false)
  const submitting = ref(false)
  const listLoading = ref(false)
  const listFinished = ref(false)
  const currentList = ref([])
  const showMediaDialog = ref(false)
  const currentMediaUrl = ref('')
  const currentPrompt = ref('')
  const page = ref(1)
  const pageSize = ref(10)
  const total = ref(0)
  const currentPowerCost = ref(0)
  const taskPulling = ref(true)
  const tastPullHandler = ref(null)

  // 功能分类
  const categories = ref([
    { key: 'image_generation', name: '图像生成' },
    { key: 'image_editing', name: '图像编辑' },
    { key: 'image_effects', name: '图像特效' },
    { key: 'video_generation', name: '视频生成' },
  ])

  // 选项数据
  const imageSizeOptions = [
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

  const videoAspectRatioOptions = [
    { label: '1:1 (正方形)', value: '1:1' },
    { label: '16:9 (横版)', value: '16:9' },
    { label: '9:16 (竖版)', value: '9:16' },
  ]

  const imageEffectsTemplateOptions = [
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

  // 功能参数
  const textToImageParams = ref({
    size: '1024x1024',
    scale: 7.5,
    use_pre_llm: false,
  })

  const imageToImageParams = ref({
    image_input: [],
    size: '1024x1024',
  })

  const imageEditParams = ref({
    image_urls: [],
    scale: 0.5,
  })

  const imageEffectsParams = ref({
    image_input1: [],
    template_id: '',
    size: '1024x1024',
  })

  const textToVideoParams = ref({
    aspect_ratio: '16:9',
  })

  const imageToVideoParams = ref({
    image_urls: [],
    aspect_ratio: '16:9',
  })

  // 计算属性
  const activeFunction = computed(() => {
    if (activeCategory.value === 'image_generation') {
      return useImageInput.value ? 'image_to_image' : 'text_to_image'
    } else if (activeCategory.value === 'image_editing') {
      return 'image_edit'
    } else if (activeCategory.value === 'video_generation') {
      return useImageInput.value ? 'image_to_video' : 'text_to_video'
    }
    return 'text_to_image'
  })

  // Actions
  const getCategoryIcon = (category) => {
    const iconMap = {
      image_generation: 'iconfont icon-image',
      image_editing: 'iconfont icon-edit',
      image_effects: 'iconfont icon-chuangzuo',
      video_generation: 'iconfont icon-video',
    }
    return iconMap[category] || 'iconfont icon-image'
  }

  const switchCategory = (key) => {
    activeCategory.value = key
    useImageInput.value = false
  }

  const switchInputMode = () => {
    currentPrompt.value = ''
  }

  const handleMultipleImageUpload = (event) => {
    const files = Array.from(event.target.files)
    files.forEach((file) => {
      if (imageToVideoParams.value.image_urls.length < 2) {
        onImageUpload({ file, name: file.name })
      }
    })
  }

  const removeImage = (index) => {
    imageToVideoParams.value.image_urls.splice(index, 1)
  }

  const onImageUpload = (file) => {
    const formData = new FormData()
    formData.append('file', file.file, file.name)
    showLoading('正在上传图片...')

    return httpPost('/api/upload', formData)
      .then((res) => {
        showMessageOK('图片上传成功')
        const imageData = { url: res.data.url, content: res.data.url }

        // 根据当前活动功能添加到相应的参数中
        if (activeFunction.value === 'image_to_image') {
          imageToImageParams.value.image_input = [imageData]
        } else if (activeFunction.value === 'image_edit') {
          imageEditParams.value.image_urls = [imageData]
        } else if (activeFunction.value === 'image_effects') {
          imageEffectsParams.value.image_input1 = [imageData]
        } else if (activeFunction.value === 'image_to_video') {
          imageToVideoParams.value.image_urls.push(imageData)
        }

        return res.data.url
      })
      .catch((e) => {
        showMessageError('图片上传失败:' + e.message)
      })
      .finally(() => {
        closeLoading()
      })
  }

  const submitTask = () => {
    if (!currentPrompt.value.trim()) {
      showMessageError('请输入提示词')
      return
    }

    submitting.value = true
    const params = {
      type: activeFunction.value,
      prompt: currentPrompt.value,
    }

    // 根据功能类型添加相应参数
    if (activeFunction.value === 'text_to_image') {
      Object.assign(params, textToImageParams.value)
    } else if (activeFunction.value === 'image_to_image') {
      Object.assign(params, imageToImageParams.value)
    } else if (activeFunction.value === 'image_edit') {
      Object.assign(params, imageEditParams.value)
    } else if (activeFunction.value === 'image_effects') {
      Object.assign(params, imageEffectsParams.value)
    } else if (activeFunction.value === 'text_to_video') {
      Object.assign(params, textToVideoParams.value)
    } else if (activeFunction.value === 'image_to_video') {
      Object.assign(params, imageToVideoParams.value)
    }

    return httpPost('/api/jimeng/create', params)
      .then(() => {
        fetchData(1)
        taskPulling.value = true
        showMessageOK('创建任务成功')
        currentPrompt.value = ''
      })
      .catch((e) => {
        showMessageError('创建任务失败：' + e.message)
      })
      .finally(() => {
        submitting.value = false
      })
  }

  const fetchData = (_page) => {
    if (_page) {
      page.value = _page
    }
    listLoading.value = true

    return httpPost('/api/jimeng/jobs', { page: page.value, page_size: pageSize.value })
      .then((res) => {
        total.value = res.data.total
        let needPull = false
        const items = []
        for (let v of res.data.items) {
          if (v.status === 'in_queue' || v.status === 'generating') {
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

  const loadMore = () => {
    page.value++
    fetchData()
  }

  const playMedia = (item) => {
    currentMediaUrl.value = item.img_url || item.video_url
    showMediaDialog.value = true
  }

  const downloadFile = (item) => {
    item.downloading = true
    const link = document.createElement('a')
    link.href = item.img_url || item.video_url
    link.download = item.title || 'file'
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    item.downloading = false
    showMessageSuccess('开始下载')
  }

  const retryTask = (id) => {
    return httpPost('/api/jimeng/retry', { id })
      .then(() => {
        showMessageOK('重试任务成功')
        fetchData(1)
      })
      .catch((e) => {
        showMessageError('重试任务失败：' + e.message)
      })
  }

  const removeJob = (item) => {
    return showConfirmDialog({
      title: '确认删除',
      message: '此操作将会删除任务相关文件，继续操作吗?',
      confirmButtonText: '确认删除',
      cancelButtonText: '取消',
    })
      .then(() => {
        return httpGet('/api/jimeng/remove', { id: item.id })
          .then(() => {
            showMessageOK('任务删除成功')
            fetchData(1)
          })
          .catch((e) => {
            showMessageError('任务删除失败：' + e.message)
          })
      })
      .catch(() => {})
  }

  const getFunctionName = (type) => {
    const nameMap = {
      text_to_image: '文生图',
      image_to_image: '图生图',
      image_edit: '图像编辑',
      image_effects: '图像特效',
      text_to_video: '文生视频',
      image_to_video: '图生视频',
    }
    return nameMap[type] || type
  }

  const getTaskType = (type) => {
    return type.includes('video') ? 'warning' : 'primary'
  }

  const startTaskPolling = () => {
    tastPullHandler.value = setInterval(() => {
      if (taskPulling.value) {
        fetchData(1)
      }
    }, 5000)
  }

  const stopTaskPolling = () => {
    if (tastPullHandler.value) {
      clearInterval(tastPullHandler.value)
    }
  }

  const resetParams = () => {
    textToImageParams.value = {
      size: '1024x1024',
      scale: 7.5,
      use_pre_llm: false,
    }
    imageToImageParams.value = {
      image_input: [],
      size: '1024x1024',
    }
    imageEditParams.value = {
      image_urls: [],
      scale: 0.5,
    }
    imageEffectsParams.value = {
      image_input1: [],
      template_id: '',
      size: '1024x1024',
    }
    textToVideoParams.value = {
      aspect_ratio: '16:9',
    }
    imageToVideoParams.value = {
      image_urls: [],
      aspect_ratio: '16:9',
    }
  }

  const closeMediaDialog = () => {
    showMediaDialog.value = false
    currentMediaUrl.value = ''
  }

  return {
    // State
    activeCategory,
    useImageInput,
    submitting,
    listLoading,
    listFinished,
    currentList,
    showMediaDialog,
    currentMediaUrl,
    currentPrompt,
    page,
    pageSize,
    total,
    currentPowerCost,
    taskPulling,
    tastPullHandler,
    categories,
    imageSizeOptions,
    videoAspectRatioOptions,
    imageEffectsTemplateOptions,
    textToImageParams,
    imageToImageParams,
    imageEditParams,
    imageEffectsParams,
    textToVideoParams,
    imageToVideoParams,

    // Computed
    activeFunction,

    // Actions
    getCategoryIcon,
    switchCategory,
    switchInputMode,
    handleMultipleImageUpload,
    removeImage,
    onImageUpload,
    submitTask,
    fetchData,
    loadMore,
    playMedia,
    downloadFile,
    retryTask,
    removeJob,
    getFunctionName,
    getTaskType,
    startTaskPolling,
    stopTaskPolling,
    resetParams,
    closeMediaDialog,
  }
})
