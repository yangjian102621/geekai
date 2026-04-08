import { closeLoading, showLoading, showMessageError, showMessageOK } from '@/utils/dialog'
import { httpDownload, httpGet, httpPost } from '@/utils/http'
import { replaceImg } from '@/utils/libs'
import Compressor from 'compressorjs'
import { ElMessage, ElMessageBox } from 'element-plus'
import { compact } from 'lodash'
import { defineStore } from 'pinia'
import { computed, onMounted, ref } from 'vue'
import { checkSession, getSystemInfo } from './cache'
import { useSharedStore } from './sharedata'

export const useSunoStore = defineStore('suno', () => {
  // 响应式数据
  const custom = ref(false)
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
    { label: '硬摇滚', value: 'hard rock' },
    { label: '电音', value: 'electronic' },
    { label: '金属', value: 'metal' },
    { label: '重金属', value: 'heavy metal' },
    { label: '节拍', value: 'beat' },
    { label: '弱拍', value: 'upbeat' },
    { label: '合成器', value: 'synth' },
    { label: '吉他', value: 'guitar' },
    { label: '钢琴', value: 'piano' },
    { label: '小提琴', value: 'violin' },
    { label: '贝斯', value: 'bass' },
    { label: '嘻哈', value: 'hip hop' },
  ])

  const data = ref({
    model: 'chirp-auk',
    tags: '',
    lyrics: '',
    prompt: '',
    title: '',
    instrumental: false,
    ref_task_id: '',
    extend_secs: 0,
    ref_song_id: '',
  })

  const loading = ref(false)
  const noData = ref(true)
  const playList = ref([])
  const showPlayer = ref(false)
  const list = ref([])
  const taskPulling = ref(true)
  const btnText = ref('开始创作')
  const refSong = ref(null)
  const showDialog = ref(false)
  const editData = ref({ title: '', cover: '', id: 0 })
  const promptPlaceholder = ref('请在这里输入你自己写的歌词...')
  const isGenerating = ref(false)
  const sunoPower = ref(0)
  const isLogin = ref(false)
  const shareStore = useSharedStore()

  // 分页相关
  const page = ref(1)
  const pageSize = ref(10)
  const total = ref(0)

  // 定时器引用
  let tastPullHandler = null

  // 计算属性
  const hasRefSong = computed(() => refSong.value !== null)

  onMounted(() => {
    getSystemInfo().then((res) => {
      sunoPower.value = res.data.suno_power
    })
    checkSession().then((res) => {
      isLogin.value = true
    })
  })

  // 方法
  const fetchData = async (_page) => {
    if (_page) {
      page.value = _page
    }
    loading.value = true

    try {
      const res = await httpGet('/api/suno/list', {
        page: page.value,
        page_size: pageSize.value,
      })

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

      loading.value = false
      taskPulling.value = needPull

      // 如果任务有变化，则刷新任务列表
      if (JSON.stringify(list.value) !== JSON.stringify(items)) {
        list.value = items
      }
      noData.value = list.value.length === 0
    } catch (e) {
      loading.value = false
      noData.value = true
      showMessageError('获取作品列表失败：' + e.message)
    }
  }

  const create = async () => {
    if (!isLogin.value) {
      return shareStore.setShowLoginDialog(true)
    }

    data.value.type = custom.value ? 2 : 1
    data.value.ref_task_id = refSong.value ? refSong.value.task_id : ''
    data.value.ref_song_id = refSong.value ? refSong.value.song_id : ''
    data.value.extend_secs = refSong.value ? refSong.value.extend_secs : 0

    // 验证输入
    if (refSong.value) {
      if (data.value.extend_secs > refSong.value.duration) {
        return showMessageError('续写开始时间不能超过原歌曲长度')
      }
    } else if (custom.value) {
      if (data.value.lyrics === '') {
        return showMessageError('请输入歌词')
      }
      if (data.value.title === '') {
        return showMessageError('请输入歌曲标题')
      }
    } else {
      if (data.value.prompt === '') {
        return showMessageError('请输入歌曲描述')
      }
    }

    try {
      await httpPost('/api/suno/create', data.value)
      await fetchData(1)
      taskPulling.value = true
      showMessageOK('创建任务成功')
    } catch (e) {
      showMessageError('创建任务失败：' + e.message)
    }
  }

  const merge = async (item) => {
    try {
      await httpPost('/api/suno/create', { song_id: item.song_id, type: 3 })
      await fetchData(1)
      taskPulling.value = true
      showMessageOK('创建任务成功')
    } catch (e) {
      showMessageError('合并歌曲失败：' + e.message)
    }
  }

  const download = async (item) => {
    const url = replaceImg(item.audio_url)
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

  const uploadAudio = async (file) => {
    // 判断是否登录
    if (!isLogin.value) {
      return shareStore.setShowLoginDialog(true)
    }

    const formData = new FormData()
    formData.append('file', file.file, file.name)
    showLoading('正在上传文件...')

    try {
      const res = await httpPost('/api/upload', formData)
      await httpPost('/api/suno/create', {
        audio_url: res.data.url,
        title: res.data.name,
        type: 4,
      })
      await fetchData(1)
      showMessageOK('歌曲上传成功')
      closeLoading()
      removeRefSong()
      ElMessage.success({ message: '上传成功', duration: 500 })
    } catch (e) {
      showMessageError('歌曲上传失败：' + e.message)
      closeLoading()
    }
  }

  const extend = (item) => {
    refSong.value = item
    refSong.value.extend_secs = item.duration
    data.value.title = item.title
    custom.value = true
    btnText.value = '续写歌曲'
    promptPlaceholder.value = '输入额外的歌词，根据您之前的歌词来扩展歌曲...'
  }

  const update = (item) => {
    showDialog.value = true
    editData.value.title = item.title
    editData.value.cover = item.cover_url
    editData.value.id = item.id
  }

  const updateSong = async () => {
    if (editData.value.title === '' || editData.value.cover === '') {
      return showMessageError('歌曲标题和封面不能为空')
    }

    try {
      await httpPost('/api/suno/update', editData.value)
      showMessageOK('更新歌曲成功')
      showDialog.value = false
      await fetchData()
    } catch (e) {
      showMessageError('更新歌曲失败：' + e.message)
    }
  }

  const removeRefSong = () => {
    refSong.value = null
    btnText.value = '开始创作'
    promptPlaceholder.value = '请在这里输入你自己写的歌词...'
  }

  const selectTag = (tag) => {
    const currentTags = data.value.tags.trim()
    const newTagLength = tag.value.length

    if (currentTags.length + newTagLength >= 119) {
      return
    }

    const tagArray = currentTags
      ? currentTags
          .split(',')
          .map((t) => t.trim())
          .filter((t) => t)
      : []
    const newTags = compact([...tagArray, tag.value])
    data.value.tags = newTags.join(',')
  }

  const removeJob = async (item) => {
    try {
      await ElMessageBox.confirm('此操作将会删除任务相关文件，继续操作码?', '删除提示', {
        confirmButtonText: '确认',
        cancelButtonText: '取消',
        type: 'warning',
      })

      await httpGet('/api/suno/remove', { id: item.id })
      ElMessage.success('任务删除成功')
      await fetchData()
    } catch (e) {
      if (e !== 'cancel') {
        ElMessage.error('任务删除失败：' + e.message)
      }
    }
  }

  const publishJob = async (item) => {
    try {
      await httpGet('/api/suno/publish', { id: item.id, publish: item.publish })
      ElMessage.success('操作成功')
    } catch (e) {
      ElMessage.error('操作失败：' + e.message)
    }
  }

  const getShareURL = (item) => {
    return `${location.protocol}//${location.host}/song/${item.song_id}`
  }

  const uploadCover = (file) => {
    new Compressor(file.file, {
      quality: 0.6,
      success(result) {
        const formData = new FormData()
        formData.append('file', result, result.name)
        showLoading('图片上传中...')

        httpPost('/api/upload', formData)
          .then((res) => {
            editData.value.cover = res.data.url
            ElMessage.success({ message: '上传成功', duration: 500 })
            closeLoading()
          })
          .catch((e) => {
            ElMessage.error('图片上传失败:' + e.message)
            closeLoading()
          })
      },
      error(err) {
        console.log(err.message)
      },
    })
  }

  const createLyric = async () => {
    if (data.value.lyrics === '') {
      return showMessageError('请输入歌词描述')
    }

    isGenerating.value = true

    try {
      const res = await httpPost('/api/prompt/lyric', { prompt: data.value.lyrics })
      const lines = res.data.split('\n')
      data.value.title = lines.shift().replace(/\*/g, '')
      lines.shift()
      data.value.lyrics = lines.join('\n')
      isGenerating.value = false
    } catch (e) {
      showMessageError('歌词生成失败：' + e.message)
      isGenerating.value = false
    }
  }

  const startTaskPolling = () => {
    tastPullHandler = setInterval(() => {
      if (taskPulling.value) {
        fetchData(1)
      }
    }, 5000)
  }

  const stopTaskPolling = () => {
    if (tastPullHandler) {
      clearInterval(tastPullHandler)
      tastPullHandler = null
    }
  }

  const resetData = () => {
    data.value = {
      model: 'chirp-auk',
      tags: '',
      lyrics: '',
      prompt: '',
      title: '',
      instrumental: false,
      ref_task_id: '',
      extend_secs: 0,
      ref_song_id: '',
    }
    custom.value = false
    refSong.value = null
    btnText.value = '开始创作'
    promptPlaceholder.value = '请在这里输入你自己写的歌词...'
  }

  return {
    // 状态
    custom,
    models,
    tags,
    data,
    loading,
    noData,
    playList,
    showPlayer,
    list,
    taskPulling,
    btnText,
    refSong,
    showDialog,
    editData,
    promptPlaceholder,
    isGenerating,
    page,
    pageSize,
    total,
    hasRefSong,
    sunoPower,
    // 方法
    fetchData,
    create,
    merge,
    download,
    uploadAudio,
    extend,
    update,
    updateSong,
    removeRefSong,
    selectTag,
    removeJob,
    publishJob,
    getShareURL,
    uploadCover,
    createLyric,
    startTaskPolling,
    stopTaskPolling,
    resetData,
  }
})
