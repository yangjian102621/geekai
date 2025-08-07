import { defineStore } from 'pinia'
import { ref, reactive } from 'vue'
import { checkSession } from '@/store/cache'
import { closeLoading, showLoading, showToastMessage } from '@/utils/dialog'
import { httpDownload, httpGet, httpPost } from '@/utils/http'
import { replaceImg } from '@/utils/libs'
import { getSystemInfo } from '@/store/cache'

export const useSunoStore = defineStore('suno', () => {
  // 状态
  const custom = ref(false)
  const data = reactive({
    model: 'chirp-auk',
    tags: '',
    lyrics: '',
    prompt: '',
    title: '',
    instrumental: false,
    ref_task_id: '',
    extend_secs: 0,
    ref_song_id: '',
    type: 1,
  })
  const loading = ref(false)
  const list = ref([])
  const listLoading = ref(false)
  const listFinished = ref(false)
  const btnText = ref('开始创作')
  const refSong = ref(null)
  const showModelPicker = ref(false)
  const showPlayer = ref(false)
  const showDeleteModal = ref(false)
  const currentAudio = ref('')
  const uploadFiles = ref([])
  const uploadRef = ref(null)
  const isGenerating = ref(false)
  const deleting = ref(false)
  const deleteItem = ref(null)
  const models = ref([
    { label: 'v3.0', value: 'chirp-v3-0' },
    { label: 'v3.5', value: 'chirp-v3-5' },
    { label: 'v4.0', value: 'chirp-v4' },
    { label: 'v4.5', value: 'chirp-auk' },
  ])
  const tags = ref([
    { label: '女声', value: 'female vocals' },
    { label: '男声', value: 'male vocals' },
    { label: '流行', value: 'pop' },
    { label: '摇滚', value: 'rock' },
    { label: '电音', value: 'electronic' },
    { label: '钢琴', value: 'piano' },
    { label: '吉他', value: 'guitar' },
    { label: '嘻哈', value: 'hip hop' },
  ])
  const page = ref(1)
  const pageSize = ref(10)
  const total = ref(0)
  const taskPulling = ref(true)
  const tastPullHandler = ref(null)
  const sunoPowerCost = ref(0)

  onMounted(() => {
    getSystemInfo().then((res) => {
      sunoPowerCost.value = res.data.suno_power
    })
  })

  // 方法
  const onModelSelect = (selectedModel) => {
    data.model = selectedModel.value
  }
  const selectTag = (tag) => {
    if (data.tags.length + tag.value.length >= 119) {
      showToastMessage('标签长度超出限制', 'error')
      return
    }
    const currentTags = data.tags.split(',').filter((t) => t.trim())
    if (!currentTags.includes(tag.value)) {
      currentTags.push(tag.value)
      data.tags = currentTags.join(',')
    }
  }
  const createLyric = () => {
    if (data.lyrics === '') {
      showToastMessage('请输入歌词描述', 'error')
      return
    }
    isGenerating.value = true
    httpPost('/api/prompt/lyric', { prompt: data.lyrics })
      .then((res) => {
        const lines = res.data.split('\n')
        data.title = lines.shift().replace(/\*/g, '')
        lines.shift()
        data.lyrics = lines.join('\n')
        showToastMessage('歌词生成成功', 'success')
      })
      .catch((e) => {
        showToastMessage('歌词生成失败：' + e.message, 'error')
      })
      .finally(() => {
        isGenerating.value = false
      })
  }
  const handleFileChange = (file) => {
    uploadFiles.value = [file]
    if (file.status === 'ready') {
      uploadAudio(file)
    }
  }
  const beforeUpload = (file) => {
    const isLt10M = file.size / 1024 / 1024 < 10
    if (!isLt10M) {
      showToastMessage('文件大小不能超过 10MB!', 'error')
      return false
    }
    return true
  }
  const uploadAudio = (file) => {
    const formData = new FormData()
    formData.append('file', file.raw, file.name)
    showLoading('正在上传文件...')
    httpPost('/api/upload', formData)
      .then((res) => {
        httpPost('/api/suno/create', {
          audio_url: res.data.url,
          title: res.data.name,
          type: 4,
        })
          .then(() => {
            fetchData(1)
            showToastMessage('歌曲上传成功', 'success')
            removeRefSong()
            uploadFiles.value = []
            if (uploadRef.value) {
              uploadRef.value.clearFiles()
            }
          })
          .catch((e) => {
            showToastMessage('歌曲上传失败：' + e.message, 'error')
          })
          .finally(() => {
            closeLoading()
          })
      })
      .catch((e) => {
        showToastMessage('文件上传失败:' + e.message, 'error')
      })
      .finally(() => {
        closeLoading()
      })
  }
  const create = () => {
    data.type = custom.value ? 2 : 1
    data.ref_task_id = refSong.value ? refSong.value.task_id : ''
    data.ref_song_id = refSong.value ? refSong.value.song_id : ''
    data.extend_secs = refSong.value ? refSong.value.extend_secs : 0
    if (refSong.value) {
      if (data.extend_secs > refSong.value.duration) {
        showToastMessage('续写开始时间不能超过原歌曲长度', 'error')
        return
      }
    } else if (custom.value) {
      if (data.lyrics === '') {
        showToastMessage('请输入歌词', 'error')
        return
      }
      if (data.title === '') {
        showToastMessage('请输入歌曲标题', 'error')
        return
      }
    } else {
      if (data.prompt === '') {
        showToastMessage('请输入歌曲描述', 'error')
        return
      }
    }
    loading.value = true
    httpPost('/api/suno/create', data)
      .then(() => {
        fetchData(1)
        taskPulling.value = true
        showToastMessage('创建任务成功', 'success')
      })
      .catch((e) => {
        showToastMessage('创建任务失败：' + e.message, 'error')
      })
      .finally(() => {
        loading.value = false
      })
  }
  const fetchData = (_page) => {
    if (_page) {
      page.value = _page
    }
    listLoading.value = true
    httpGet('/api/suno/list', { page: page.value, page_size: pageSize.value })
      .then((res) => {
        total.value = res.data.total
        let needPull = false
        const items = []
        for (let v of res.data.items) {
          if (v.progress === 100) {
            v.major_model_version = v['raw_data']['major_model_version']
          }
          if (v.progress === 0 || v.progress === 102) {
            needPull = true
          }
          items.push(v)
        }
        listLoading.value = false
        taskPulling.value = needPull
        if (page.value === 1) {
          list.value = items
        } else {
          list.value.push(...items)
        }
        if (items.length < pageSize.value) {
          listFinished.value = true
        }
      })
      .catch((e) => {
        listLoading.value = false
        showToastMessage('获取作品列表失败：' + e.message, 'error')
      })
  }
  const loadMore = () => {
    if (!listFinished.value && !listLoading.value) {
      page.value++
      fetchData()
    }
  }
  const refreshFirstPage = () => {
    const currentPage = page.value
    const currentList = [...list.value]
    httpGet('/api/suno/list', { page: 1, page_size: pageSize.value })
      .then((res) => {
        let needPull = false
        const firstPageItems = []
        for (let v of res.data.items) {
          if (v.progress === 100) {
            v.major_model_version = v['raw_data']['major_model_version']
          }
          if (v.progress === 0 || v.progress === 102) {
            needPull = true
          }
          firstPageItems.push(v)
        }
        taskPulling.value = needPull
        if (currentPage === 1) {
          list.value = firstPageItems
        } else {
          const otherPagesData = currentList.slice(pageSize.value)
          list.value = [...firstPageItems, ...otherPagesData]
        }
      })
      .catch((e) => {
        console.error('刷新第一页数据失败：', e)
      })
  }
  const play = (item) => {
    currentAudio.value = item.audio_url
    showPlayer.value = true
  }
  const download = (item) => {
    const url = replaceImg(item.audio_url)
    const downloadURL = `${import.meta.env.VITE_API_HOST}/api/download?url=${url}`
    const urlObj = new URL(url)
    const fileName = urlObj.pathname.split('/').pop()
    item.downloading = true
    httpDownload(downloadURL)
      .then((response) => {
        const blob = new Blob([response.data])
        const link = document.createElement('a')
        link.href = URL.createObjectURL(blob)
        link.download = fileName
        document.body.appendChild(link)
        link.click()
        document.body.removeChild(link)
        URL.revokeObjectURL(link.href)
        item.downloading = false
      })
      .catch(() => {
        showToastMessage('下载失败', 'error')
        item.downloading = false
      })
      .finally(() => {
        item.downloading = false
      })
  }
  const showDeleteDialog = (item) => {
    deleteItem.value = item
    // 这里建议在页面层处理弹窗，store 只负责数据和业务
  }
  const extend = (item) => {
    refSong.value = item
    refSong.value.extend_secs = item.duration
    data.title = item.title
    custom.value = true
    btnText.value = '续写歌曲'
    window.scrollTo({ top: 0, behavior: 'smooth' })
  }
  const removeRefSong = () => {
    refSong.value = null
    btnText.value = '开始创作'
  }

  // 副作用（定时轮询、滚动监听）建议在页面层处理，store 只暴露方法

  return {
    // 状态
    custom,
    data,
    loading,
    list,
    listLoading,
    listFinished,
    btnText,
    refSong,
    showModelPicker,
    showPlayer,
    showDeleteModal,
    currentAudio,
    uploadFiles,
    uploadRef,
    isGenerating,
    deleting,
    deleteItem,
    models,
    tags,
    page,
    pageSize,
    total,
    taskPulling,
    tastPullHandler,
    sunoPowerCost,
    // 方法
    onModelSelect,
    selectTag,
    createLyric,
    handleFileChange,
    beforeUpload,
    uploadAudio,
    create,
    fetchData,
    loadMore,
    refreshFirstPage,
    play,
    download,
    showDeleteDialog,
    extend,
    removeRefSong,
  }
})
