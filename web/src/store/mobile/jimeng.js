import { showMessageError, showMessageOK } from '@/utils/dialog'
import { httpDownload, httpGet, httpPost } from '@/utils/http'
import { replaceImg } from '@/utils/libs'
import { defineStore } from 'pinia'
import { showConfirmDialog } from 'vant'
import { computed, reactive, ref, watch } from 'vue'

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

  // 新增：算力配置
  const powerConfig = ref({
    text_to_image: 20,
    image_to_image: 30,
    image_edit: 25,
    image_effects: 15,
    text_to_video: 100,
    image_to_video: 120,
  })

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
    image_urls: [],
    aspect_ratio: '16:9',
    seed: -1,
  })

  // 计算属性
  const activeFunction = computed(() => {
    if (activeCategory.value === 'image_generation') {
      return useImageInput.value ? 'image_to_image' : 'text_to_image'
    } else if (activeCategory.value === 'image_editing') {
      return 'image_edit'
    } else if (activeCategory.value === 'image_effects') {
      return 'image_effects'
    } else if (activeCategory.value === 'video_generation') {
      return useImageInput.value ? 'image_to_video' : 'text_to_video'
    }
    return 'text_to_image'
  })

  // 新增：动态计算当前算力消耗
  const updateCurrentPowerCost = () => {
    const functionKey = activeFunction.value
    currentPowerCost.value = powerConfig.value[functionKey] || 10
  }

  // 监听任务类型变化，自动更新算力
  watch(
    [activeCategory, useImageInput],
    () => {
      updateCurrentPowerCost()
    },
    { immediate: true }
  )

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

  // 新增：获取算力配置
  const fetchPowerConfig = async () => {
    try {
      const res = await httpGet('/api/jimeng/power-config')
      if (res.data) {
        powerConfig.value = res.data
        updateCurrentPowerCost() // 更新当前算力消耗
      }
    } catch (error) {
      console.error('获取算力配置失败:', error)
    }
  }

  const submitTask = () => {
    if (!currentPrompt.value.trim()) {
      showMessageError('请输入提示词')
      return
    }

    submitting.value = true
    let requestData = { task_type: activeFunction.value, prompt: currentPrompt.value }
    // 根据功能类型添加相应参数
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

    return httpPost('/api/jimeng/task', requestData)
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
        if (res.data.items) {
          for (let v of res.data.items) {
            if (v.status === 'in_queue' || v.status === 'generating') {
              needPull = true
            }
            items.push(v)
          }
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

  const retryTask = (id) => {
    return httpGet('/api/jimeng/retry', { id })
      .then(() => {
        showMessageOK('重试任务成功')
        fetchData(1)
      })
      .catch((e) => {
        showMessageError('重试任务失败：' + e.message)
      })
  }

  const removeJob = async (item) => {
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

  const closeMediaDialog = () => {
    showMediaDialog.value = false
    currentMediaUrl.value = ''
  }

  // 新增：复制提示词功能
  const copyPrompt = (prompt) => {
    navigator.clipboard
      .writeText(prompt)
      .then(() => {
        showMessageOK('提示词已复制')
      })
      .catch(() => {
        showMessageError('复制失败')
      })
  }

  // 新增：复制错误信息功能
  const copyErrorMsg = (msg) => {
    navigator.clipboard
      .writeText(msg)
      .then(() => {
        showMessageOK('错误信息已复制')
      })
      .catch(() => {
        showMessageError('复制失败')
      })
  }

  // 新增：初始化方法
  const init = async () => {
    await fetchPowerConfig()
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
    powerConfig,

    // Computed
    activeFunction,

    // Actions
    getCategoryIcon,
    switchCategory,
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
    closeMediaDialog,
    fetchPowerConfig,
    copyPrompt,
    copyErrorMsg,
    init,
  }
})
