<template>
  <div class="page-suno">
    <div class="left-bar">
      <!-- 参数设置区域 -->
      <div class="space-y-6">
        <!-- 自定义开关 -->
        <div class="setting-card">
          <div class="flex items-center justify-between">
            <div>
              <span class="card-title">自定义模式</span>
              <p class="card-description">可以更精确地控制生成内容</p>
            </div>
            <el-switch v-model="store.custom" size="large" />
          </div>
        </div>

        <div class="setting-card">
          <div class="flex items-center justify-between mb-3">
            <label class="card-label">选择模型</label>
            <el-popover
              placement="right"
              :width="200"
              trigger="hover"
              content="选择不同的模型，可以获得不同的生成效果"
            >
              <template #reference>
                <el-icon class="help-icon">
                  <InfoFilled />
                </el-icon>
              </template>
            </el-popover>
          </div>
          <el-select v-model="store.data.model" placeholder="请选择模型" class="w-full">
            <el-option
              v-for="model in store.models"
              :key="model.value"
              :label="model.label"
              :value="model.value"
            />
          </el-select>
        </div>

        <!-- 纯音乐开关 -->
        <div class="setting-card">
          <div class="flex items-center justify-between">
            <div>
              <span class="card-title">纯音乐</span>
              <p class="card-description">生成不包含人声的音乐</p>
            </div>
            <el-switch v-model="store.data.instrumental" size="large" />
          </div>
        </div>

        <!-- 自定义模式内容 -->
        <div v-if="store.custom" class="space-y-6">
          <!-- 歌词输入 -->
          <div v-if="!store.data.instrumental" class="setting-card">
            <div class="flex items-center justify-between mb-3">
              <label class="card-label">歌词</label>
              <el-popover
                placement="right"
                :width="200"
                trigger="hover"
                content="自己写歌词或寻求 AI 的帮助。使用两节歌词（8 行）可获得最佳效果。"
              >
                <template #reference>
                  <el-icon class="help-icon">
                    <InfoFilled />
                  </el-icon>
                </template>
              </el-popover>
            </div>
            <div class="relative">
              <el-input
                v-model="store.data.lyrics"
                type="textarea"
                :rows="8"
                :placeholder="store.promptPlaceholder"
                :maxlength="1024"
                :show-word-limit="true"
                resize="none"
                class="mb-3"
              />
              <div class="flex justify-end">
                <button @click="store.createLyric" class="lyric-btn" :disabled="store.isGenerating">
                  <i class="iconfont icon-magic text-xs"></i>
                  <span v-if="!store.isGenerating">生成歌词</span>
                  <span v-else class="flex items-center space-x-1">
                    <i class="iconfont icon-loading animate-spin text-xs"></i>
                    <span>生成中...</span>
                  </span>
                </button>
              </div>
            </div>
          </div>

          <!-- 音乐风格 -->
          <div class="setting-card">
            <div class="flex items-center justify-between mb-3">
              <label class="card-label">音乐风格</label>
              <el-popover
                placement="right"
                :width="200"
                trigger="hover"
                content="描述您想要的音乐风格（例如：原声流行音乐）。Sunos 模特无法识别艺术家的名字，但能够理解音乐流派和氛围。"
              >
                <template #reference>
                  <el-icon class="help-icon">
                    <InfoFilled />
                  </el-icon>
                </template>
              </el-popover>
            </div>
            <el-input
              v-model="store.data.tags"
              type="textarea"
              :rows="3"
              :maxlength="120"
              :show-word-limit="true"
              resize="none"
              placeholder="请输入音乐风格，多个风格之间用英文逗号隔开..."
              class="mb-4"
            />
            <!-- 风格标签选择 -->
            <div class="flex flex-wrap gap-2">
              <button
                v-for="tag in store.tags"
                :key="tag.value"
                @click="store.selectTag(tag)"
                class="tag-btn"
              >
                {{ tag.label }}
              </button>
            </div>
          </div>

          <!-- 歌曲名称 -->
          <div class="setting-card">
            <div class="flex items-center justify-between mb-3">
              <label class="card-label">歌曲名称</label>
              <el-popover
                placement="right"
                :width="200"
                trigger="hover"
                content="给你的歌曲起一个标题，以便于分享、发现和组织。"
              >
                <template #reference>
                  <el-icon class="help-icon">
                    <InfoFilled />
                  </el-icon>
                </template>
              </el-popover>
            </div>
            <el-input
              v-model="store.data.title"
              placeholder="请输入歌曲名称..."
              maxlength="100"
              show-word-limit
            />
          </div>
        </div>

        <!-- 简单模式内容 -->
        <div v-else class="setting-card">
          <div class="flex items-center justify-between mb-3">
            <label class="card-label">歌曲描述</label>
          </div>
          <el-input
            v-model="store.data.prompt"
            type="textarea"
            :rows="8"
            :maxlength="1024"
            :show-word-limit="true"
            resize="none"
            placeholder="描述您想要的音乐风格和主题（例如：关于假期的流行音乐）。请使用流派和氛围，而不是特定的艺术家和歌曲风格，AI无法识别。"
          />
        </div>

        <!-- 续写歌曲 -->
        <div v-if="store.refSong" class="setting-card extend-song-card">
          <div class="flex items-center justify-between mb-3">
            <h3 class="card-title">续写歌曲</h3>
            <button @click="store.removeRefSong" class="remove-btn">移除</button>
          </div>
          <div class="space-y-3">
            <div class="flex justify-between">
              <span class="text-secondary">歌曲名称：</span>
              <span class="text-primary font-medium">{{ store.refSong.title }}</span>
            </div>
            <div>
              <label class="block text-secondary font-medium mb-2">续写开始时间(秒)</label>
              <input
                v-model="store.refSong.extend_secs"
                type="number"
                placeholder="从第几秒开始续写"
                class="extend-input"
              />
            </div>
          </div>
        </div>

        <!-- 生成按钮 -->
        <div class="setting-card">
          <button @click="store.create" :disabled="store.loading" class="create-btn">
            <i v-if="store.loading" class="iconfont icon-loading animate-spin"></i>
            <i v-else class="iconfont icon-chuangzuo"></i>
            <span
              >{{ store.loading ? '创作中...' : store.btnText }} ({{ store.sunoPower }}
              <i class="iconfont icon-vip2 !text-xs"></i>
              )</span
            >
          </button>
        </div>

        <!-- 上传音乐 -->
        <div class="setting-card">
          <div class="flex items-center justify-between mb-3">
            <label class="card-label">上传音乐文件</label>
            <el-popover
              placement="right"
              :width="200"
              trigger="hover"
              content="上传你自己的音乐文件，然后进行二次创作"
            >
              <template #reference>
                <el-icon class="help-icon">
                  <InfoFilled />
                </el-icon>
              </template>
            </el-popover>
          </div>

          <!-- 上传区域 -->
          <el-upload
            class="custom-upload"
            drag
            :auto-upload="true"
            :show-file-list="false"
            :http-request="store.uploadAudio"
            accept=".wav,.mp3"
            :limit="1"
          >
            <template #trigger>
              <div class="p-2">
                <el-button class="upload-btn" size="large" type="primary">
                  <i class="iconfont icon-upload mr-2"></i>
                  <span>上传音乐</span>
                </el-button>
              </div>
            </template>
          </el-upload>

          <!-- 上传提示 -->
          <div class="upload-tips">
            <p>• 上传你自己的音乐文件，然后进行二次创作</p>
            <p>• 请上传6-60秒的原始音频</p>
            <p>• 检测到人声的音频将仅设为私人音频</p>
          </div>
        </div>
      </div>
    </div>

    <!-- 右侧作品列表 -->
    <div
      class="right-box"
      v-loading="store.loading"
      element-loading-background="rgba(100,100,100,0.3)"
    >
      <div class="list-box" v-if="!store.noData">
        <div v-for="item in store.list" :key="item.id" class="song-card">
          <div class="flex space-x-4">
            <div class="flex-shrink-0">
              <div class="song-cover">
                <el-image
                  :src="item.cover_url"
                  fit="cover"
                  class="w-full h-full"
                  :preview-disabled="true"
                >
                  <template #error>
                    <div class="cover-placeholder">
                      <i class="iconfont icon-mp3 text-gray-400 text-xl"></i>
                    </div>
                  </template>
                </el-image>
                <!-- 音乐播放按钮 -->
                <button v-if="item.progress === 100" @click="play(item)" class="play-overlay">
                  <i class="iconfont icon-play text-white text-xl"></i>
                </button>
                <!-- 进度动画 -->
                <div v-if="item.progress < 100 && item.progress !== 101" class="progress-overlay">
                  <i class="iconfont icon-loading animate-spin text-blue-500 text-xl"></i>
                </div>
                <!-- 失败状态 -->
                <div v-if="item.progress === 101" class="error-overlay">
                  <i class="iconfont icon-warning text-red-500 text-xl"></i>
                </div>
              </div>
            </div>
            <div class="flex-1 min-w-0">
              <div class="flex items-start justify-between">
                <div class="flex-1">
                  <h3 class="song-title">
                    <a :href="'/song/' + item.song_id" target="_blank" class="song-link">
                      {{ item.title || '未命名歌曲' }}
                    </a>
                  </h3>
                  <p class="song-description">
                    {{ item.tags || item.prompt }}
                  </p>
                </div>
                <!-- 任务状态 -->
                <div v-if="item.progress < 100" class="task-status">
                  <div v-if="item.progress === 101" class="status-error">
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M12 9v2m0 4h.01"
                      />
                    </svg>
                    <span>失败</span>
                  </div>
                  <div v-else class="status-loading">
                    <div class="loading-spinner"></div>
                    <span>生成中</span>
                  </div>
                </div>
              </div>
              <!-- 标签 -->
              <div class="song-tags">
                <span v-if="item.major_model_version" class="model-tag">
                  {{ item.major_model_version }}
                </span>
                <span v-if="item.type === 4" class="upload-tag">
                  <i class="iconfont icon-upload mr-1"></i>用户上传
                </span>
                <span v-if="item.type === 3" class="full-song-tag">
                  <i class="iconfont icon-mp3 mr-1"></i>完整歌曲
                </span>
                <span v-if="item.ref_song" class="extend-tag">
                  <i class="iconfont icon-link mr-1"></i>续写
                </span>
              </div>
            </div>
          </div>

          <!-- 操作按钮 -->
          <div class="song-actions">
            <div class="action-buttons">
              <button v-if="item.progress === 100" @click="play(item)" class="action-btn play-btn">
                <i class="iconfont icon-play text-xs"></i>
                <span>播放</span>
              </button>
              <button
                v-if="item.progress === 100"
                @click="store.download(item)"
                :disabled="item.downloading"
                class="action-btn download-btn"
              >
                <svg
                  v-if="item.downloading"
                  class="w-3 h-3 animate-spin"
                  fill="none"
                  viewBox="0 0 24 24"
                >
                  <circle
                    class="opacity-25"
                    cx="12"
                    cy="12"
                    r="10"
                    stroke="currentColor"
                    stroke-width="4"
                  />
                  <path
                    class="opacity-75"
                    fill="currentColor"
                    d="m4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                  />
                </svg>
                <i v-else class="iconfont icon-download text-xs"></i>
                <span>{{ item.downloading ? '下载中...' : '下载' }}</span>
              </button>
              <button
                v-if="item.progress === 100 && item.ref_song"
                @click="store.merge(item)"
                class="action-btn merge-btn"
              >
                <i class="iconfont icon-concat text-xs"></i>
                <span>合并</span>
              </button>
              <button
                v-if="item.progress !== 101"
                @click="store.extend(item)"
                class="action-btn extend-btn"
              >
                <i class="iconfont icon-edit text-xs"></i>
                <span>续写</span>
              </button>
            </div>
            <div class="action-buttons">
              <button
                v-if="item.progress !== 101"
                @click="store.update(item)"
                class="action-btn edit-btn"
              >
                <i class="iconfont icon-edit text-xs"></i>
                <span>编辑</span>
              </button>
              <button @click="store.removeJob(item)" class="action-btn delete-btn">
                <i class="iconfont icon-remove text-xs"></i>
                <span>删除</span>
              </button>
            </div>
          </div>

          <!-- 进度条 -->
          <div v-if="item.progress < 100 && item.progress !== 101" class="progress-bar">
            <div class="progress-info">
              <span>生成进度</span>
              <span>{{ item.progress }}%</span>
            </div>
            <div class="progress-track">
              <div class="progress-fill" :style="{ width: item.progress + '%' }"></div>
            </div>
          </div>

          <!-- 错误信息 -->
          <div v-if="item.progress === 101" class="error-message">
            <div class="flex items-start space-x-2">
              <div>
                <p class="error-text">{{ item.err_msg || '未知错误' }}</p>
              </div>
            </div>
          </div>
        </div>
      </div>

      <el-empty
        :image-size="100"
        :image="nodata"
        description="没有任何作品，赶紧去创作吧！"
        v-else
      />

      <div class="pagination">
        <el-pagination
          v-if="store.total > store.pageSize"
          background
          layout="total,prev, pager, next"
          :hide-on-single-page="true"
          :current-page="store.page"
          :page-size="store.pageSize"
          @current-change="store.fetchData"
          @size-change="
            (size) => {
              store.pageSize = size
              store.fetchData(1)
            }
          "
          :total="store.total"
        />
      </div>

      <div class="music-player" v-if="store.showPlayer">
        <music-player
          :songs="store.playList"
          ref="playerRef"
          :show-close="true"
          @close="store.showPlayer = false"
        />
      </div>
    </div>

    <!-- 编辑对话框 -->
    <el-dialog
      v-model="store.showDialog"
      title="修改歌曲"
      width="500px"
      :before-close="
        () => {
          store.showDialog = false
        }
      "
    >
      <form class="form">
        <div class="form-item">
          <div class="label">歌曲名称</div>
          <el-input v-model="store.editData.title" type="text" />
        </div>

        <div class="form-item">
          <div class="label">封面图片</div>
          <el-upload
            class="avatar-uploader"
            :auto-upload="true"
            :show-file-list="false"
            :http-request="store.uploadCover"
            accept=".png,.jpg,.jpeg,.bmp"
          >
            <el-avatar :src="store.editData.cover" shape="square" :size="100" />
          </el-upload>
        </div>
      </form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="store.showDialog = false">取消</el-button>
          <el-button type="primary" @click="store.updateSong">确认</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import nodata from '@/assets/img/no-data.png'

import MusicPlayer from '@/components/MusicPlayer.vue'
import { checkSession } from '@/store/cache'
import { useSunoStore } from '@/store/suno'
import { InfoFilled } from '@element-plus/icons-vue'
import { nextTick, onMounted, onUnmounted, ref, watch } from 'vue'

// 使用 Pinia store
const store = useSunoStore()

// 组件内部状态
const playerRef = ref(null)

// 播放音乐
const play = (item) => {
  store.playList = [item]
  store.showPlayer = true
  nextTick(() => playerRef.value.play())
}

// 监听器
watch(
  () => store.custom,
  (newValue) => {
    if (!newValue) {
      store.removeRefSong()
    }
  }
)

// 生命周期
onMounted(() => {
  // 检查会话并初始化数据
  checkSession()
    .then(() => {
      store.fetchData(1)
      store.startTaskPolling()
    })
    .catch(() => {})
})

onUnmounted(() => {
  // 清理资源
  store.stopTaskPolling()
})
</script>

<style lang="scss" scoped>
@import '../assets/css/suno.scss';
</style>
