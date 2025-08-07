<template>
  <div class="page-suno">
    <div class="left-bar">
      <!-- 顶部工具栏 -->
      <div class="bg-white rounded-xl p-4 shadow-sm mb-6">
        <div class="flex items-center justify-between">
          <div class="flex items-center space-x-4">
            <el-tooltip content="定义模式" placement="top">
              <div class="flex items-center space-x-2">
                <span class="text-gray-700 font-medium">自定义模式</span>
                <el-switch v-model="store.custom" size="large" />
              </div>
            </el-tooltip>
          </div>

          <div class="flex items-center space-x-3">
            <el-tooltip
              content="请上传6-60秒的原始音频，检测到人声的音频将仅设为私人音频。"
              placement="bottom-end"
            >
              <el-upload
                class="avatar-uploader"
                :auto-upload="true"
                :show-file-list="false"
                :http-request="store.uploadAudio"
                accept=".wav,.mp3"
              >
                <el-button
                  class="bg-blue-600 hover:bg-blue-700 text-white border-0"
                  round
                  type="primary"
                >
                  <i class="iconfont icon-upload mr-2"></i>
                  <span>上传音乐</span>
                </el-button>
              </el-upload>
            </el-tooltip>

            <el-select v-model="store.data.model" placeholder="请选择模型" class="w-32">
              <el-option
                v-for="model in store.models"
                :key="model.value"
                :label="model.label"
                :value="model.value"
              />
            </el-select>
          </div>
        </div>
      </div>

      <!-- 参数设置区域 -->
      <div class="space-y-6">
        <!-- 自定义开关 -->
        <div class="bg-white rounded-xl p-4 shadow-sm">
          <div class="flex items-center justify-between">
            <div>
              <span class="text-gray-900 font-medium">自定义模式</span>
              <p class="text-sm text-gray-500 mt-1">使用自定义模式，可以更精确地控制生成内容</p>
            </div>
            <el-switch v-model="store.data.instrumental" size="default" />
          </div>
        </div>

        <!-- 纯音乐开关 -->
        <div class="bg-white rounded-xl p-4 shadow-sm">
          <div class="flex items-center justify-between">
            <div>
              <span class="text-gray-900 font-medium">纯音乐</span>
              <p class="text-sm text-gray-500 mt-1">生成不包含人声的音乐</p>
            </div>
            <el-switch v-model="store.data.instrumental" size="default" />
          </div>
        </div>

        <!-- 自定义模式内容 -->
        <div v-if="store.custom" class="space-y-6">
          <!-- 歌词输入 -->
          <div v-if="!store.data.instrumental" class="bg-white rounded-xl p-4 shadow-sm">
            <div class="flex items-center justify-between mb-3">
              <label class="text-gray-700 font-medium">歌词</label>
              <el-popover
                placement="right"
                :width="200"
                trigger="hover"
                content="自己写歌词或寻求 AI 的帮助。使用两节歌词（8 行）可获得最佳效果。"
              >
                <template #reference>
                  <el-icon class="text-gray-400 cursor-help">
                    <InfoFilled />
                  </el-icon>
                </template>
              </el-popover>
            </div>
            <div
              class="relative"
              v-loading="store.isGenerating"
              element-loading-text="正在生成歌词..."
            >
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
              <button
                @click="store.createLyric"
                class="absolute bottom-3 right-3 px-3 py-1.5 bg-blue-600 text-white text-sm rounded-lg hover:bg-blue-700 transition-colors flex items-center space-x-1"
              >
                <i class="iconfont icon-magic text-xs"></i>
                <span>生成歌词</span>
              </button>
            </div>
          </div>

          <!-- 音乐风格 -->
          <div class="bg-white rounded-xl p-4 shadow-sm">
            <div class="flex items-center justify-between mb-3">
              <label class="text-gray-700 font-medium">音乐风格</label>
              <el-popover
                placement="right"
                :width="200"
                trigger="hover"
                content="描述您想要的音乐风格（例如：原声流行音乐）。Sunos 模特无法识别艺术家的名字，但能够理解音乐流派和氛围。"
              >
                <template #reference>
                  <el-icon class="text-gray-400 cursor-help">
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
                class="px-3 py-1.5 text-sm border border-blue-200 text-blue-600 rounded-full hover:bg-blue-50 transition-colors hover:border-blue-300"
              >
                {{ tag.label }}
              </button>
            </div>
          </div>

          <!-- 歌曲名称 -->
          <div class="bg-white rounded-xl p-4 shadow-sm">
            <div class="flex items-center justify-between mb-3">
              <label class="text-gray-700 font-medium">歌曲名称</label>
              <el-popover
                placement="right"
                :width="200"
                trigger="hover"
                content="给你的歌曲起一个标题，以便于分享、发现和组织。"
              >
                <template #reference>
                  <el-icon class="text-gray-400 cursor-help">
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
        <div v-else class="bg-white rounded-xl p-4 shadow-sm">
          <div class="flex items-center justify-between mb-3">
            <label class="text-gray-700 font-medium">歌曲描述</label>
            <el-popover
              placement="right"
              :width="200"
              trigger="hover"
              content="描述您想要的音乐风格和主题（例如：关于假期的流行音乐）。请使用流派和氛围，而不是特定的艺术家和歌曲风格，AI无法识别。"
            >
              <template #reference>
                <el-icon class="text-gray-400 cursor-help">
                  <InfoFilled />
                </el-icon>
              </template>
            </el-popover>
          </div>
          <el-input
            v-model="store.data.prompt"
            type="textarea"
            :rows="8"
            :maxlength="1024"
            :show-word-limit="true"
            resize="none"
            placeholder="例如：一首关于爱情的摇滚歌曲..."
          />
        </div>

        <!-- 续写歌曲 -->
        <div
          v-if="store.refSong"
          class="bg-white rounded-xl p-4 shadow-sm border-l-4 border-orange-400"
        >
          <div class="flex items-center justify-between mb-3">
            <h3 class="text-gray-900 font-medium">续写歌曲</h3>
            <button
              @click="store.removeRefSong"
              class="px-3 py-1.5 text-sm bg-red-100 text-red-600 rounded-lg hover:bg-red-200 transition-colors"
            >
              移除
            </button>
          </div>
          <div class="space-y-3">
            <div class="flex justify-between">
              <span class="text-gray-600">歌曲名称：</span>
              <span class="text-gray-900 font-medium">{{ store.refSong.title }}</span>
            </div>
            <div>
              <label class="block text-gray-700 font-medium mb-2">续写开始时间(秒)</label>
              <input
                v-model="store.refSong.extend_secs"
                type="number"
                placeholder="从第几秒开始续写"
                class="w-full px-4 py-3 border border-gray-200 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              />
            </div>
          </div>
        </div>

        <!-- 上传音乐 -->
        <div class="bg-white rounded-xl p-4 shadow-sm">
          <div class="flex items-center justify-between mb-3">
            <label class="text-gray-700 font-medium">上传音乐文件</label>
            <el-popover
              placement="right"
              :width="200"
              trigger="hover"
              content="支持 .wav, .mp3 格式，文件大小不超过 50MB"
            >
              <template #reference>
                <el-icon class="text-gray-400 cursor-help">
                  <InfoFilled />
                </el-icon>
              </template>
            </el-popover>
          </div>

          <!-- 上传区域 -->
          <div
            class="border-2 border-dashed border-gray-300 rounded-lg p-6 text-center hover:border-blue-400 hover:bg-blue-50 transition-colors cursor-pointer relative"
            @click="triggerFileInput"
            @drop="handleDrop"
            @dragover.prevent
            @dragenter.prevent
          >
            <input
              ref="fileInput"
              type="file"
              accept=".wav,.mp3"
              @change="handleFileSelect"
              class="hidden"
            />

            <!-- 上传中状态 -->
            <div v-if="uploading" class="flex flex-col items-center space-y-3">
              <div
                class="w-12 h-12 border-4 border-blue-200 border-t-blue-600 rounded-full animate-spin"
              ></div>
              <div class="text-center">
                <p class="text-gray-700 font-medium">正在上传...</p>
                <p class="text-sm text-gray-500">{{ uploadProgress }}%</p>
              </div>
            </div>

            <!-- 默认状态 -->
            <div
              v-else-if="!uploadedFile"
              @click="triggerFileInput"
              class="flex flex-col items-center space-y-3"
            >
              <div class="w-16 h-16 bg-blue-100 rounded-full flex items-center justify-center">
                <i class="iconfont icon-upload text-blue-500 text-2xl"></i>
              </div>
              <div class="text-center">
                <p class="text-gray-700 font-medium">点击或拖拽上传音乐文件</p>
                <p class="text-sm text-gray-500 mt-1">支持 .wav, .mp3 格式，最大 10MB</p>
              </div>
            </div>

            <!-- 已上传文件预览 -->
            <div v-else class="flex items-center space-x-4 w-full">
              <div class="flex-shrink-0">
                <div class="w-12 h-12 bg-blue-100 rounded-lg flex items-center justify-center">
                  <i class="iconfont icon-mp3 text-blue-500 text-xl"></i>
                </div>
              </div>
              <div class="flex-1 text-left">
                <p class="text-gray-700 font-medium truncate">{{ uploadedFile.name }}</p>
                <p class="text-sm text-gray-500">{{ formatFileSize(uploadedFile.size) }}</p>
              </div>
              <div class="flex space-x-2">
                <button
                  @click.stop="playUploadedFile"
                  class="p-2 bg-blue-100 text-blue-600 rounded-lg hover:bg-blue-200 transition-colors"
                  title="播放"
                >
                  <i class="iconfont icon-play text-sm"></i>
                </button>
                <button
                  @click.stop="removeUploadedFile"
                  class="p-2 bg-red-100 text-red-600 rounded-lg hover:bg-red-200 transition-colors"
                  title="移除"
                >
                  <i class="iconfont icon-remove text-sm"></i>
                </button>
              </div>
            </div>
          </div>

          <!-- 上传提示 -->
          <div class="mt-3 text-sm text-gray-500">
            <p>• 请上传6-60秒的原始音频</p>
            <p>• 支持格式：WAV, MP3, 最大 10MB</p>
            <p>• 检测到人声的音频将仅设为私人音频</p>
          </div>
        </div>

        <!-- 生成按钮 -->
        <div class="bg-white rounded-xl p-4 shadow-sm">
          <button
            @click="store.create"
            :disabled="store.loading"
            class="w-full py-4 bg-gradient-to-r from-blue-500 to-purple-600 text-white font-semibold rounded-xl disabled:from-gray-400 disabled:to-gray-400 disabled:cursor-not-allowed hover:from-blue-600 hover:to-purple-700 transition-all duration-200 flex items-center justify-center space-x-2"
          >
            <i v-if="store.loading" class="iconfont icon-loading animate-spin"></i>
            <span>{{ store.loading ? '创作中...' : store.btnText }}</span>
          </button>
        </div>
      </div>
    </div>

    <!-- 右侧作品列表 -->
    <div
      class="right-box h-dvh"
      v-loading="store.loading"
      element-loading-background="rgba(100,100,100,0.3)"
    >
      <div class="list-box" v-if="!store.noData">
        <div
          v-for="item in store.list"
          :key="item.id"
          class="bg-white rounded-xl p-4 shadow-sm mb-4"
        >
          <div class="flex space-x-4">
            <div class="flex-shrink-0">
              <div class="relative w-16 h-16 rounded-lg overflow-hidden bg-gray-100">
                <el-image
                  :src="item.cover_url"
                  fit="cover"
                  class="w-full h-full"
                  :preview-disabled="true"
                >
                  <template #error>
                    <div class="w-full h-full flex items-center justify-center bg-gray-100">
                      <i class="iconfont icon-mp3 text-gray-400 text-xl"></i>
                    </div>
                  </template>
                </el-image>
                <!-- 音乐播放按钮 -->
                <button
                  v-if="item.progress === 100"
                  @click="play(item)"
                  class="absolute inset-0 flex items-center justify-center bg-black bg-opacity-50 opacity-0 hover:opacity-100 transition-opacity"
                >
                  <i class="iconfont icon-play text-white text-xl"></i>
                </button>
                <!-- 进度动画 -->
                <div
                  v-if="item.progress < 100 && item.progress !== 101"
                  class="absolute inset-0 flex items-center justify-center bg-blue-500 bg-opacity-20"
                >
                  <i class="iconfont icon-loading animate-spin text-blue-500 text-xl"></i>
                </div>
                <!-- 失败状态 -->
                <div
                  v-if="item.progress === 101"
                  class="absolute inset-0 flex items-center justify-center bg-red-500 bg-opacity-20"
                >
                  <i class="iconfont icon-warning text-red-500 text-xl"></i>
                </div>
              </div>
            </div>
            <div class="flex-1 min-w-0">
              <div class="flex items-start justify-between">
                <div class="flex-1">
                  <h3 class="text-gray-900 font-medium truncate">
                    <a
                      :href="'/song/' + item.song_id"
                      target="_blank"
                      class="hover:text-blue-600 transition-colors"
                    >
                      {{ item.title || '未命名歌曲' }}
                    </a>
                  </h3>
                  <p class="text-gray-500 text-sm mt-1 line-clamp-2">
                    {{ item.tags || item.prompt }}
                  </p>
                </div>
                <!-- 任务状态 -->
                <div v-if="item.progress < 100" class="flex items-center space-x-2 text-sm">
                  <div
                    v-if="item.progress === 101"
                    class="text-red-600 flex items-center space-x-1"
                  >
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
                  <div v-else class="text-blue-600 flex items-center space-x-1">
                    <div
                      class="w-3 h-3 border border-blue-600 border-t-transparent rounded-full animate-spin"
                    ></div>
                    <span>生成中</span>
                  </div>
                </div>
              </div>
              <!-- 标签 -->
              <div class="flex items-center space-x-2 mt-2">
                <span
                  v-if="item.major_model_version"
                  class="px-2 py-1 text-xs bg-blue-100 text-blue-600 rounded-full"
                >
                  {{ item.major_model_version }}
                </span>
                <span
                  v-if="item.type === 4"
                  class="px-2 py-1 text-xs bg-green-100 text-green-600 rounded-full"
                >
                  <i class="iconfont icon-upload mr-1"></i>用户上传
                </span>
                <span
                  v-if="item.type === 3"
                  class="px-2 py-1 text-xs bg-yellow-100 text-yellow-600 rounded-full"
                >
                  <i class="iconfont icon-mp3 mr-1"></i>完整歌曲
                </span>
                <span
                  v-if="item.ref_song"
                  class="px-2 py-1 text-xs bg-purple-100 text-purple-600 rounded-full"
                >
                  <i class="iconfont icon-link mr-1"></i>续写
                </span>
              </div>
            </div>
          </div>

          <!-- 操作按钮 -->
          <div class="flex items-center justify-between mt-4">
            <div class="flex space-x-2">
              <button
                v-if="item.progress === 100"
                @click="play(item)"
                class="px-3 py-1.5 bg-blue-600 text-white text-sm rounded-lg hover:bg-blue-700 transition-colors flex items-center space-x-1"
              >
                <i class="iconfont icon-play text-xs"></i>
                <span>播放</span>
              </button>
              <button
                v-if="item.progress === 100"
                @click="store.download(item)"
                :disabled="item.downloading"
                class="px-3 py-1.5 bg-green-600 text-white text-sm rounded-lg hover:bg-green-700 transition-colors disabled:bg-gray-400 flex items-center space-x-1"
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
                class="px-3 py-1.5 bg-purple-600 text-white text-sm rounded-lg hover:bg-purple-700 transition-colors flex items-center space-x-1"
              >
                <i class="iconfont icon-concat text-xs"></i>
                <span>合并</span>
              </button>
              <button
                @click="store.extend(item)"
                class="px-3 py-1.5 bg-orange-600 text-white text-sm rounded-lg hover:bg-orange-700 transition-colors flex items-center space-x-1"
              >
                <i class="iconfont icon-edit text-xs"></i>
                <span>续写</span>
              </button>
            </div>
            <div class="flex space-x-2">
              <button
                @click="store.update(item)"
                class="px-3 py-1.5 bg-gray-600 text-white text-sm rounded-lg hover:bg-gray-700 transition-colors flex items-center space-x-1"
              >
                <i class="iconfont icon-edit text-xs"></i>
                <span>编辑</span>
              </button>
              <button
                @click="store.removeJob(item)"
                class="px-3 py-1.5 bg-red-100 text-red-600 text-sm rounded-lg hover:bg-red-200 transition-colors flex items-center space-x-1"
              >
                <i class="iconfont icon-remove text-xs"></i>
                <span>删除</span>
              </button>
            </div>
          </div>

          <!-- 进度条 -->
          <div v-if="item.progress < 100 && item.progress !== 101" class="mt-4">
            <div class="flex justify-between text-sm text-gray-600 mb-1">
              <span>生成进度</span>
              <span>{{ item.progress }}%</span>
            </div>
            <div class="w-full bg-gray-200 rounded-full h-2">
              <div
                class="bg-blue-600 h-2 rounded-full transition-all duration-300"
                :style="{ width: item.progress + '%' }"
              ></div>
            </div>
          </div>

          <!-- 错误信息 -->
          <div
            v-if="item.progress === 101"
            class="mt-4 p-3 bg-red-50 border border-red-200 rounded-lg"
          >
            <div class="flex items-start space-x-2">
              <div>
                <p class="text-red-600 text-sm">{{ item.err_msg || '未知错误' }}</p>
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
          style="--el-pagination-button-bg-color: rgba(86, 86, 95, 0.2)"
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
import Clipboard from 'clipboard'
import { ElMessage } from 'element-plus'
import { nextTick, onMounted, onUnmounted, ref, watch } from 'vue'

// 使用 Pinia store
const store = useSunoStore()

// 组件内部状态
const playerRef = ref(null)
const clipboard = ref(null)
const fileInput = ref(null)
const uploading = ref(false)
const uploadProgress = ref(0)
const uploadedFile = ref(null)
const uploadedAudioUrl = ref('')

// 播放音乐
const play = (item) => {
  store.playList = [item]
  store.showPlayer = true
  nextTick(() => playerRef.value.play())
}

// 文件上传相关方法
const triggerFileInput = () => {
  fileInput.value?.click()
}

const formatFileSize = (bytes) => {
  if (bytes === 0) return '0 Bytes'
  const k = 1024
  const sizes = ['Bytes', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

const handleFileSelect = (event) => {
  const file = event.target.files[0]
  if (file) {
    validateAndUploadFile(file)
  }
}

const handleDrop = (event) => {
  event.preventDefault()
  const files = event.dataTransfer.files
  if (files.length > 0) {
    validateAndUploadFile(files[0])
  }
}

const validateAndUploadFile = (file) => {
  // 验证文件类型
  const allowedTypes = ['audio/wav', 'audio/mp3', 'audio/mpeg']
  if (!allowedTypes.includes(file.type)) {
    ElMessage.error('只支持 WAV 和 MP3 格式的音频文件')
    return
  }

  // 验证文件大小 (50MB)
  const maxSize = 50 * 1024 * 1024
  if (file.size > maxSize) {
    ElMessage.error('文件大小不能超过 50MB')
    return
  }

  uploadedFile.value = file
  uploadFile(file)
}

const uploadFile = async (file) => {
  uploading.value = true
  uploadProgress.value = 0

  try {
    const formData = new FormData()
    formData.append('file', file)

    // 模拟上传进度
    const progressInterval = setInterval(() => {
      if (uploadProgress.value < 90) {
        uploadProgress.value += Math.random() * 10
      }
    }, 200)

    // 只上传文件，不创建任务
    const { httpPost } = await import('@/utils/http')
    const res = await httpPost('/api/upload', formData)

    clearInterval(progressInterval)
    uploadProgress.value = 100

    // 保存上传的音频URL
    uploadedAudioUrl.value = res.data.url

    ElMessage.success('文件上传成功')

    // 延迟重置状态
    setTimeout(() => {
      uploading.value = false
      uploadProgress.value = 0
    }, 1000)
  } catch (error) {
    uploading.value = false
    uploadProgress.value = 0
    uploadedFile.value = null
    ElMessage.error('文件上传失败：' + error.message)
  }
}

const playUploadedFile = () => {
  if (uploadedAudioUrl.value) {
    const audio = new Audio(uploadedAudioUrl.value)
    audio.play()
  }
}

const removeUploadedFile = () => {
  uploadedFile.value = null
  uploadedAudioUrl.value = ''
  if (fileInput.value) {
    fileInput.value.value = ''
  }
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
  // 初始化剪贴板
  clipboard.value = new Clipboard('.copy-link')
  clipboard.value.on('success', () => {
    ElMessage.success('复制歌曲链接成功！')
  })
  clipboard.value.on('error', () => {
    ElMessage.error('复制失败！')
  })

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
  if (clipboard.value) {
    clipboard.value.destroy()
  }
  store.stopTaskPolling()
})
</script>

<style lang="scss" scoped>
@use '../assets/css/suno.scss' as *;

/* 文本截断 */
.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

/* 动画 */
@keyframes spin {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}

.animate-spin {
  animation: spin 1s linear infinite;
}

/* 自定义滚动条 */
::-webkit-scrollbar {
  width: 6px;
}

::-webkit-scrollbar-track {
  background: #f1f1f1;
  border-radius: 3px;
}

::-webkit-scrollbar-thumb {
  background: #c1c1c1;
  border-radius: 3px;
}

::-webkit-scrollbar-thumb:hover {
  background: #a8a8a8;
}
</style>
