<template>
  <div class="mobile-suno-create">
    <!-- 页面头部 -->
    <div class="page-header">
      <van-nav-bar title="音乐创作" left-arrow @click-left="goBack" fixed placeholder />
    </div>

    <!-- 创作表单 -->
    <div class="create-form">
      <!-- 模式切换 -->
      <div class="mode-switch">
        <van-cell-group>
          <van-cell title="创作模式">
            <template #right-icon>
              <van-switch v-model="custom" size="24" @change="onModeChange" />
            </template>
          </van-cell>
          <van-cell title="自定义模式" :value="custom ? '开启' : '关闭'" />
        </van-cell-group>
      </div>

      <!-- 模型选择 -->
      <div class="model-select">
        <van-field
          v-model="selectedModelLabel"
          label="模型"
          readonly
          is-link
          @click="showModelPicker = true"
          placeholder="请选择模型"
        />
      </div>

      <!-- 纯音乐开关 -->
      <div class="pure-music">
        <van-cell title="纯音乐">
          <template #right-icon>
            <van-switch v-model="data.instrumental" size="24" />
          </template>
        </van-cell>
      </div>

      <!-- 自定义模式内容 -->
      <div v-if="custom">
        <!-- 歌词输入 -->
        <div v-if="!data.instrumental" class="lyrics-section">
          <van-field
            v-model="data.lyrics"
            label="歌词"
            type="textarea"
            placeholder="请在这里输入你自己写的歌词..."
            rows="6"
            maxlength="2000"
            show-word-limit
          />
          <van-button
            type="primary"
            size="small"
            @click="createLyric"
            :loading="isGenerating"
            block
            class="mt-2"
          >
            生成歌词
          </van-button>
        </div>

        <!-- 音乐风格 -->
        <div class="style-section">
          <van-field
            v-model="data.tags"
            label="音乐风格"
            type="textarea"
            placeholder="请输入音乐风格，多个风格之间用英文逗号隔开..."
            rows="3"
            maxlength="120"
            show-word-limit
          />
          <!-- 风格标签选择 -->
          <div class="style-tags">
            <van-tag
              v-for="tag in tags"
              :key="tag.value"
              type="primary"
              plain
              size="medium"
              @click="selectTag(tag)"
              class="mr-2 mb-2"
            >
              {{ tag.label }}
            </van-tag>
          </div>
        </div>

        <!-- 歌曲名称 -->
        <div class="title-section">
          <van-field
            v-model="data.title"
            label="歌曲名称"
            placeholder="请输入歌曲名称..."
            maxlength="100"
            show-word-limit
          />
        </div>
      </div>

      <!-- 简单模式内容 -->
      <div v-else>
        <van-field
          v-model="data.prompt"
          label="歌曲描述"
          type="textarea"
          placeholder="例如：一首关于爱情的摇滚歌曲..."
          rows="6"
          maxlength="1000"
          show-word-limit
        />
      </div>

      <!-- 续写歌曲 -->
      <div v-if="refSong" class="ref-song">
        <van-cell title="续写歌曲">
          <template #value>
            <van-button type="danger" size="small" @click="removeRefSong"> 移除 </van-button>
          </template>
        </van-cell>
        <van-cell title="歌曲名称" :value="refSong.title" />
        <van-field
          v-model="refSong.extend_secs"
          label="续写开始时间(秒)"
          type="number"
          placeholder="从第几秒开始续写"
        />
      </div>

      <!-- 上传音乐 -->
      <div class="upload-section">
        <div class="upload-area">
          <van-uploader
            v-model="uploadFiles"
            :max-count="1"
            :after-read="uploadAudio"
            accept=".wav,.mp3"
            :preview-size="80"
            :preview-image="false"
          >
            <div class="upload-placeholder">
              <van-icon name="plus" size="24" />
              <span>上传音乐文件</span>
              <small>支持 .wav, .mp3 格式</small>
            </div>
          </van-uploader>
        </div>
      </div>

      <!-- 生成按钮 -->
      <div class="submit-section">
        <van-button type="primary" size="large" @click="create" :loading="loading" block>
          {{ btnText }}
        </van-button>
      </div>
    </div>

    <!-- 作品列表 -->
    <div class="works-list">
      <van-list
        v-model:loading="listLoading"
        :finished="listFinished"
        finished-text="没有更多了"
        @load="loadMore"
      >
        <div v-for="item in list" :key="item.id" class="work-item">
          <van-card
            :title="item.title || '未命名歌曲'"
            :desc="item.tags || item.prompt"
            :thumb="item.cover_url"
          >
            <template #tags>
              <van-tag v-if="item.major_model_version" type="primary">
                {{ item.major_model_version }}
              </van-tag>
              <van-tag v-if="item.type === 4" type="success">用户上传</van-tag>
              <van-tag v-if="item.type === 3" type="warning">完整歌曲</van-tag>
            </template>
            <template #footer>
              <van-button v-if="item.progress === 100" size="small" @click="play(item)">
                播放
              </van-button>
              <van-button
                v-if="item.progress === 100"
                size="small"
                @click="download(item)"
                :loading="item.downloading"
              >
                下载
              </van-button>
              <van-button size="small" type="danger" @click="removeJob(item)"> 删除 </van-button>
            </template>
          </van-card>
        </div>
      </van-list>
    </div>

    <!-- 模型选择弹窗 -->
    <van-popup v-model:show="showModelPicker" position="bottom" round>
      <van-picker
        :columns="modelOptions"
        @confirm="onModelConfirm"
        @cancel="showModelPicker = false"
        title="选择模型"
      />
    </van-popup>

    <!-- 音乐播放器 -->
    <van-popup v-model:show="showPlayer" position="bottom" round :style="{ height: '40%' }">
      <div class="player-content">
        <div class="player-header">
          <h3>正在播放</h3>
          <van-icon name="cross" @click="showPlayer = false" />
        </div>
        <audio v-if="currentAudio" :src="currentAudio" controls autoplay class="w-full" />
      </div>
    </van-popup>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { showToast, showDialog } from 'vant'
import { httpGet, httpPost } from '@/utils/http'
import { checkSession } from '@/store/cache'

const router = useRouter()

// 响应式数据
const custom = ref(false)
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
const list = ref([])
const listLoading = ref(false)
const listFinished = ref(false)
const btnText = ref('开始创作')
const refSong = ref(null)
const showModelPicker = ref(false)
const showPlayer = ref(false)
const currentAudio = ref('')
const uploadFiles = ref([])
const isGenerating = ref(false)

// 模型选项
const models = ref([
  { label: 'v3.0', value: 'chirp-v3-0' },
  { label: 'v3.5', value: 'chirp-v3-5' },
  { label: 'v4.0', value: 'chirp-v4' },
  { label: 'v4.5', value: 'chirp-auk' },
])

const modelOptions = models.value.map((item) => item.label)

// 计算当前选中的模型标签
const selectedModelLabel = computed(() => {
  const selectedModel = models.value.find((item) => item.value === data.value.model)
  return selectedModel ? selectedModel.label : ''
})

// 风格标签
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

// 页面数据
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)
const taskPulling = ref(true)
const tastPullHandler = ref(null)

// 生命周期
onMounted(() => {
  checkSession()
    .then(() => {
      fetchData(1)
      tastPullHandler.value = setInterval(() => {
        if (taskPulling.value) {
          fetchData(1)
        }
      }, 5000)
    })
    .catch(() => {})
})

onUnmounted(() => {
  if (tastPullHandler.value) {
    clearInterval(tastPullHandler.value)
  }
})

// 方法
const goBack = () => {
  router.back()
}

const onModeChange = () => {
  if (!custom.value) {
    removeRefSong()
  }
}

const onModelConfirm = (value) => {
  const selectedModel = models.value.find((item) => item.label === value)
  if (selectedModel) {
    data.value.model = selectedModel.value
  }
  showModelPicker.value = false
}

const selectTag = (tag) => {
  if (data.value.tags.length + tag.value.length >= 119) {
    showToast('标签长度超出限制')
    return
  }
  const currentTags = data.value.tags.split(',').filter((t) => t.trim())
  if (!currentTags.includes(tag.value)) {
    currentTags.push(tag.value)
    data.value.tags = currentTags.join(',')
  }
}

const createLyric = () => {
  if (data.value.lyrics === '') {
    showToast('请输入歌词描述')
    return
  }
  isGenerating.value = true
  httpPost('/api/prompt/lyric', { prompt: data.value.lyrics })
    .then((res) => {
      const lines = res.data.split('\n')
      data.value.title = lines.shift().replace(/\*/g, '')
      lines.shift()
      data.value.lyrics = lines.join('\n')
      showToast('歌词生成成功')
    })
    .catch((e) => {
      showToast('歌词生成失败：' + e.message)
    })
    .finally(() => {
      isGenerating.value = false
    })
}

const uploadAudio = (file) => {
  const formData = new FormData()
  formData.append('file', file.file, file.name)
  showToast({ message: '正在上传文件...', duration: 0 })

  httpPost('/api/upload', formData)
    .then((res) => {
      httpPost('/api/suno/create', {
        audio_url: res.data.url,
        title: res.data.name,
        type: 4,
      })
        .then(() => {
          fetchData(1)
          showToast('歌曲上传成功')
          removeRefSong()
        })
        .catch((e) => {
          showToast('歌曲上传失败：' + e.message)
        })
    })
    .catch((e) => {
      showToast('文件上传失败:' + e.message)
    })
}

const create = () => {
  data.value.type = custom.value ? 2 : 1
  data.value.ref_task_id = refSong.value ? refSong.value.task_id : ''
  data.value.ref_song_id = refSong.value ? refSong.value.song_id : ''
  data.value.extend_secs = refSong.value ? refSong.value.extend_secs : 0

  if (refSong.value) {
    if (data.value.extend_secs > refSong.value.duration) {
      showToast('续写开始时间不能超过原歌曲长度')
      return
    }
  } else if (custom.value) {
    if (data.value.lyrics === '') {
      showToast('请输入歌词')
      return
    }
    if (data.value.title === '') {
      showToast('请输入歌曲标题')
      return
    }
  } else {
    if (data.value.prompt === '') {
      showToast('请输入歌曲描述')
      return
    }
  }

  loading.value = true
  httpPost('/api/suno/create', data.value)
    .then(() => {
      fetchData(1)
      taskPulling.value = true
      showToast('创建任务成功')
    })
    .catch((e) => {
      showToast('创建任务失败：' + e.message)
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
      showToast('获取作品列表失败：' + e.message)
    })
}

const loadMore = () => {
  page.value++
  fetchData()
}

const play = (item) => {
  currentAudio.value = item.audio_url
  showPlayer.value = true
}

const download = (item) => {
  item.downloading = true
  const link = document.createElement('a')
  link.href = item.audio_url
  link.download = item.title || 'song.mp3'
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
  item.downloading = false
  showToast('开始下载')
}

const removeJob = (item) => {
  showDialog({
    title: '确认删除',
    message: '此操作将会删除任务相关文件，继续操作吗?',
    showCancelButton: true,
  })
    .then(() => {
      httpGet('/api/suno/remove', { id: item.id })
        .then(() => {
          showToast('任务删除成功')
          fetchData(1)
        })
        .catch((e) => {
          showToast('任务删除失败：' + e.message)
        })
    })
    .catch(() => {})
}

const removeRefSong = () => {
  refSong.value = null
  btnText.value = '开始创作'
}
</script>

<style lang="scss" scoped>
.mobile-suno-create {
  min-height: 100vh;
  background: #f7f8fa;
  padding-bottom: 20px;

  .page-header {
    background: #fff;
  }

  .create-form {
    background: #fff;
    margin: 12px;
    border-radius: 12px;
    padding: 20px;
    box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);

    .mode-switch {
      margin-bottom: 20px;
    }

    .model-select {
      margin-bottom: 20px;
    }

    .pure-music {
      margin-bottom: 20px;
    }

    .lyrics-section,
    .style-section,
    .title-section {
      margin-bottom: 20px;
    }

    .style-tags {
      margin-top: 12px;
      display: flex;
      flex-wrap: wrap;
      gap: 8px;
    }

    .ref-song {
      margin-bottom: 20px;
      padding: 16px;
      background: #f8f9fa;
      border-radius: 8px;
      border: 1px solid #e9ecef;
    }

    .upload-section {
      margin-bottom: 20px;

      .upload-area {
        .upload-placeholder {
          display: flex;
          flex-direction: column;
          align-items: center;
          justify-content: center;
          height: 120px;
          background: #f8f9fa;
          border: 2px dashed #dee2e6;
          border-radius: 12px;
          color: #6c757d;
          transition: all 0.3s ease;

          &:hover {
            border-color: var(--van-primary-color);
            background: #f0f8ff;
          }

          .van-icon {
            margin-bottom: 8px;
            color: var(--van-primary-color);
          }

          span {
            font-size: 16px;
            font-weight: 500;
            margin-bottom: 4px;
          }

          small {
            font-size: 12px;
            opacity: 0.7;
          }
        }
      }
    }

    .submit-section {
      margin-top: 24px;
    }
  }

  .works-list {
    margin: 12px;

    .work-item {
      margin-bottom: 12px;
      background: #fff;
      border-radius: 12px;
      overflow: hidden;
      box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
    }
  }

  .player-content {
    padding: 20px;

    .player-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 16px;

      h3 {
        margin: 0;
        font-size: 18px;
        font-weight: 600;
      }

      .van-icon {
        font-size: 20px;
        cursor: pointer;
        color: #999;
      }
    }

    audio {
      width: 100%;
      border-radius: 8px;
    }
  }
}

// 深色主题适配
:deep(.van-theme-dark) {
  .mobile-suno-create {
    background: #1a1a1a;

    .create-form {
      background: #2a2a2a;
      box-shadow: 0 2px 12px rgba(0, 0, 0, 0.2);

      .ref-song {
        background: #333;
        border-color: #444;
      }

      .upload-area .upload-placeholder {
        background: #333;
        border-color: #555;

        &:hover {
          background: #2a2a2a;
          border-color: var(--van-primary-color);
        }
      }
    }

    .works-list .work-item {
      background: #2a2a2a;
      box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
    }
  }
}
</style>
