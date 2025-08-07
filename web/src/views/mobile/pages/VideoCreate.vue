<template>
  <div class="min-h-screen bg-gray-50">
    <!-- 页面头部 -->
    <div class="sticky top-0 z-40 bg-white shadow-sm">
      <div class="flex items-center px-4 h-14">
        <button
          @click="goBack"
          class="flex items-center justify-center w-8 h-8 rounded-full hover:bg-gray-100 transition-colors"
        >
          <i class="iconfont icon-back text-gray-600"></i>
        </button>
        <h1 class="flex-1 text-center text-lg font-semibold text-gray-900">视频创作</h1>
        <div class="w-8"></div>
      </div>
    </div>

    <!-- 视频类型切换 -->
    <div class="p-4 space-y-6">
      <!-- 视频类型选择 -->
      <div class="bg-white rounded-xl p-4 shadow-sm">
        <div class="flex space-x-2">
          <button
            @click="switchVideoType('luma')"
            :class="[
              'flex-1 py-3 px-4 rounded-lg font-medium transition-colors',
              activeVideoType === 'luma'
                ? 'bg-blue-600 text-white'
                : 'bg-gray-100 text-gray-700 hover:bg-gray-200',
            ]"
          >
            Luma视频
          </button>
          <button
            @click="switchVideoType('keling')"
            :class="[
              'flex-1 py-3 px-4 rounded-lg font-medium transition-colors',
              activeVideoType === 'keling'
                ? 'bg-blue-600 text-white'
                : 'bg-gray-100 text-gray-700 hover:bg-gray-200',
            ]"
          >
            可灵视频
          </button>
        </div>
      </div>

      <!-- Luma 视频参数 -->
      <div v-if="activeVideoType === 'luma'" class="space-y-6">
        <!-- 提示词输入 -->
        <div class="bg-white rounded-xl p-4 shadow-sm">
          <label class="block text-gray-700 font-medium mb-3">提示词</label>
          <textarea
            v-model="lumaParams.prompt"
            placeholder="请在此输入视频提示词，用逗号分割"
            class="w-full px-4 py-3 border border-gray-200 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent resize-none"
            rows="4"
            maxlength="2000"
          />
          <div class="text-right mt-2">
            <span class="text-sm text-gray-500">{{ lumaParams.prompt.length }}/2000</span>
          </div>
        </div>

        <!-- 提示词生成按钮 -->
        <div class="bg-white rounded-xl p-4 shadow-sm">
          <button
            @click="generatePrompt"
            :disabled="isGenerating"
            class="w-full py-3 bg-blue-600 text-white rounded-lg font-medium disabled:bg-gray-300 disabled:cursor-not-allowed hover:bg-blue-700 transition-colors flex items-center justify-center space-x-2"
          >
            <i v-if="isGenerating" class="iconfont icon-loading animate-spin"></i>
            <i v-else class="iconfont icon-chuangzuo"></i>
            <span>{{ isGenerating ? '生成中...' : '生成AI视频提示词' }}</span>
          </button>
        </div>

        <!-- 图片辅助生成开关 -->
        <div class="bg-white rounded-xl p-4 shadow-sm">
          <div class="flex items-center justify-between">
            <div>
              <span class="text-gray-900 font-medium">使用图片辅助生成</span>
              <p class="text-sm text-gray-500 mt-1">上传起始帧和结束帧图片</p>
            </div>
            <el-switch v-model="lumaUseImageMode" @change="toggleLumaImageMode" size="default" />
          </div>
        </div>

        <!-- 图片上传区域 -->
        <div v-if="lumaUseImageMode" class="bg-white rounded-xl p-4 shadow-sm">
          <div class="grid grid-cols-2 gap-4">
            <div class="relative">
              <div
                class="border-2 border-dashed border-gray-300 rounded-lg p-6 text-center hover:border-blue-400 hover:bg-blue-50 transition-colors cursor-pointer h-32"
              >
                <input
                  ref="lumaStartInput"
                  type="file"
                  accept=".jpg,.png,.jpeg"
                  @change="handleLumaStartImageUpload"
                  class="hidden"
                />
                <div
                  @click="$refs.lumaStartInput.click()"
                  class="flex flex-col items-center space-y-2 h-full justify-center"
                >
                  <i
                    v-if="!lumaStartImage.length"
                    class="iconfont icon-upload text-blue-500 text-xl"
                  ></i>
                  <span v-if="!lumaStartImage.length" class="text-gray-700 text-sm">起始帧</span>
                  <div v-else class="w-full h-full relative">
                    <el-image
                      :src="lumaStartImage[0]?.url || lumaStartImage[0]?.content"
                      fit="cover"
                      class="w-full h-full rounded"
                    />
                    <button
                      @click.stop="lumaStartImage = []"
                      class="absolute top-1 right-1 w-6 h-6 bg-red-500 text-white rounded-full flex items-center justify-center text-xs hover:bg-red-600 transition-colors"
                    >
                      <i class="iconfont icon-close"></i>
                    </button>
                  </div>
                </div>
              </div>
            </div>
            <div class="relative">
              <div
                class="border-2 border-dashed border-gray-300 rounded-lg p-6 text-center hover:border-blue-400 hover:bg-blue-50 transition-colors cursor-pointer h-32"
              >
                <input
                  ref="lumaEndInput"
                  type="file"
                  accept=".jpg,.png,.jpeg"
                  @change="handleLumaEndImageUpload"
                  class="hidden"
                />
                <div
                  @click="$refs.lumaEndInput.click()"
                  class="flex flex-col items-center space-y-2 h-full justify-center"
                >
                  <i
                    v-if="!lumaEndImage.length"
                    class="iconfont icon-upload text-blue-500 text-xl"
                  ></i>
                  <span v-if="!lumaEndImage.length" class="text-gray-700 text-sm">结束帧</span>
                  <div v-else class="w-full h-full relative">
                    <el-image
                      :src="lumaEndImage[0]?.url || lumaEndImage[0]?.content"
                      fit="cover"
                      class="w-full h-full rounded"
                    />
                    <button
                      @click.stop="lumaEndImage = []"
                      class="absolute top-1 right-1 w-6 h-6 bg-red-500 text-white rounded-full flex items-center justify-center text-xs hover:bg-red-600 transition-colors"
                    >
                      <i class="iconfont icon-close"></i>
                    </button>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Luma 特有参数 -->
        <div class="bg-white rounded-xl p-4 shadow-sm">
          <div class="space-y-4">
            <div class="flex items-center justify-between">
              <span class="text-gray-900 font-medium">循环参考图</span>
              <el-switch v-model="lumaParams.loop" size="default" />
            </div>
            <div class="flex items-center justify-between">
              <span class="text-gray-900 font-medium">提示词优化</span>
              <el-switch v-model="lumaParams.expand_prompt" size="default" />
            </div>
          </div>
        </div>

        <!-- 算力显示和生成按钮 -->
        <div class="bg-white rounded-xl p-4 shadow-sm">
          <div class="flex items-center justify-between mb-4">
            <span class="text-gray-700">当前可用算力</span>
            <span class="text-blue-600 font-semibold">{{ availablePower }}</span>
          </div>
          <button
            @click="createLumaVideo"
            :disabled="generating"
            class="w-full py-4 bg-gradient-to-r from-blue-500 to-purple-600 text-white font-semibold rounded-xl disabled:from-gray-400 disabled:to-gray-400 disabled:cursor-not-allowed hover:from-blue-600 hover:to-purple-700 transition-all duration-200 flex items-center justify-center space-x-2"
          >
            <i v-if="generating" class="iconfont icon-loading animate-spin"></i>
            <span>{{ generating ? '创作中...' : `立即生成 (${lumaPowerCost}算力)` }}</span>
          </button>
        </div>
      </div>

      <!-- KeLing 视频参数 -->
      <div v-if="activeVideoType === 'keling'" class="space-y-6">
        <!-- 画面比例 -->
        <CustomSelect
          v-model="kelingParams.aspect_ratio"
          :options="aspectRatioOptions.map((ratio) => ({ label: ratio, value: ratio }))"
          label="画面比例"
          title="选择比例"
        />

        <!-- 模型选择 -->
        <CustomSelect
          v-model="kelingParams.model"
          :options="modelOptions.map((model) => ({ label: model, value: model }))"
          label="模型选择"
          title="选择模型"
        />

        <!-- 视频时长 -->
        <CustomSelect
          v-model="kelingParams.duration"
          :options="
            durationOptions.map((duration) => ({ label: `${duration}秒`, value: duration }))
          "
          label="视频时长"
          title="选择时长"
        />

        <!-- 生成模式 -->
        <CustomSelect
          v-model="kelingParams.mode"
          :options="
            modeOptions.map((mode) => ({
              label: mode === 'std' ? '标准模式' : '专业模式',
              value: mode,
            }))
          "
          label="生成模式"
          title="选择模式"
        />

        <!-- 创意程度 -->
        <div class="bg-white rounded-xl p-4 shadow-sm">
          <div class="space-y-4">
            <label class="block text-gray-700 font-medium">创意程度</label>
            <el-slider v-model="kelingParams.cfg_scale" :min="0" :max="1" :step="0.1" />
          </div>
        </div>

        <!-- 运镜控制 -->
        <CustomSelect
          v-model="kelingParams.camera_control.type"
          :options="
            cameraControlOptions.map((option) => ({
              label: getCameraControlLabel(option),
              value: option,
            }))
          "
          label="运镜控制"
          title="选择运镜类型"
        />

        <!-- 图片辅助生成开关 -->
        <div class="bg-white rounded-xl p-4 shadow-sm">
          <div class="flex items-center justify-between">
            <div>
              <span class="text-gray-900 font-medium">使用图片辅助生成</span>
              <p class="text-sm text-gray-500 mt-1">上传起始帧和结束帧图片</p>
            </div>
            <el-switch
              v-model="kelingUseImageMode"
              @change="toggleKelingImageMode"
              size="default"
            />
          </div>
        </div>

        <!-- 图片上传区域 -->
        <div v-if="kelingUseImageMode" class="bg-white rounded-xl p-4 shadow-sm">
          <div class="grid grid-cols-2 gap-4">
            <div class="relative">
              <div
                class="border-2 border-dashed border-gray-300 rounded-lg p-6 text-center hover:border-blue-400 hover:bg-blue-50 transition-colors cursor-pointer h-32"
              >
                <input
                  ref="kelingStartInput"
                  type="file"
                  accept=".jpg,.png,.jpeg"
                  @change="handleKelingStartImageUpload"
                  class="hidden"
                />
                <div
                  @click="$refs.kelingStartInput.click()"
                  class="flex flex-col items-center space-y-2 h-full justify-center"
                >
                  <i
                    v-if="!kelingStartImage.length"
                    class="iconfont icon-upload text-blue-500 text-xl"
                  ></i>
                  <span v-if="!kelingStartImage.length" class="text-gray-700 text-sm">起始帧</span>
                  <div v-else class="w-full h-full relative">
                    <el-image
                      :src="kelingStartImage[0]?.url || kelingStartImage[0]?.content"
                      fit="cover"
                      class="w-full h-full rounded"
                    />
                    <button
                      @click.stop="kelingStartImage = []"
                      class="absolute top-1 right-1 w-6 h-6 bg-red-500 text-white rounded-full flex items-center justify-center text-xs hover:bg-red-600 transition-colors"
                    >
                      <i class="iconfont icon-close"></i>
                    </button>
                  </div>
                </div>
              </div>
            </div>
            <div class="relative">
              <div
                class="border-2 border-dashed border-gray-300 rounded-lg p-6 text-center hover:border-blue-400 hover:bg-blue-50 transition-colors cursor-pointer h-32"
              >
                <input
                  ref="kelingEndInput"
                  type="file"
                  accept=".jpg,.png,.jpeg"
                  @change="handleKelingEndImageUpload"
                  class="hidden"
                />
                <div
                  @click="$refs.kelingEndInput.click()"
                  class="flex flex-col items-center space-y-2 h-full justify-center"
                >
                  <i
                    v-if="!kelingEndImage.length"
                    class="iconfont icon-upload text-blue-500 text-xl"
                  ></i>
                  <span v-if="!kelingEndImage.length" class="text-gray-700 text-sm">结束帧</span>
                  <div v-else class="w-full h-full relative">
                    <el-image
                      :src="kelingEndImage[0]?.url || kelingEndImage[0]?.content"
                      fit="cover"
                      class="w-full h-full rounded"
                    />
                    <button
                      @click.stop="kelingEndImage = []"
                      class="absolute top-1 right-1 w-6 h-6 bg-red-500 text-white rounded-full flex items-center justify-center text-xs hover:bg-red-600 transition-colors"
                    >
                      <i class="iconfont icon-close"></i>
                    </button>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 提示词输入 -->
        <div class="bg-white rounded-xl p-4 shadow-sm">
          <label class="block text-gray-700 font-medium mb-3">提示词</label>
          <textarea
            v-model="kelingParams.prompt"
            :placeholder="kelingUseImageMode ? '描述视频画面细节' : '请在此输入视频提示词'"
            class="w-full px-4 py-3 border border-gray-200 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent resize-none"
            rows="4"
            maxlength="500"
          />
          <div class="text-right mt-2">
            <span class="text-sm text-gray-500">{{ kelingParams.prompt.length }}/500</span>
          </div>
        </div>

        <!-- 提示词生成按钮 -->
        <div class="bg-white rounded-xl p-4 shadow-sm">
          <button
            @click="generatePrompt"
            :disabled="isGenerating"
            class="w-full py-3 bg-blue-600 text-white rounded-lg font-medium disabled:bg-gray-300 disabled:cursor-not-allowed hover:bg-blue-700 transition-colors flex items-center justify-center space-x-2"
          >
            <i v-if="isGenerating" class="iconfont icon-loading animate-spin"></i>
            <i v-else class="iconfont icon-chuangzuo"></i>
            <span>{{ isGenerating ? '生成中...' : '生成专业视频提示词' }}</span>
          </button>
        </div>

        <!-- 排除内容 -->
        <div class="bg-white rounded-xl p-4 shadow-sm">
          <label class="block text-gray-700 font-medium mb-3">不希望出现的内容</label>
          <textarea
            v-model="kelingParams.negative_prompt"
            placeholder="请在此输入你不希望出现在视频上的内容"
            class="w-full px-4 py-3 border border-gray-200 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent resize-none"
            rows="3"
            maxlength="500"
          />
          <div class="text-right mt-2">
            <span class="text-sm text-gray-500">{{ kelingParams.negative_prompt.length }}/500</span>
          </div>
        </div>

        <!-- 算力显示和生成按钮 -->
        <div class="bg-white rounded-xl p-4 shadow-sm">
          <div class="flex items-center justify-between mb-4">
            <span class="text-gray-700">当前可用算力</span>
            <span class="text-blue-600 font-semibold">{{ availablePower }}</span>
          </div>
          <button
            @click="createKelingVideo"
            :disabled="generating"
            class="w-full py-4 bg-gradient-to-r from-blue-500 to-purple-600 text-white font-semibold rounded-xl disabled:from-gray-400 disabled:to-gray-400 disabled:cursor-not-allowed hover:from-blue-600 hover:to-purple-700 transition-all duration-200 flex items-center justify-center space-x-2"
          >
            <i v-if="generating" class="iconfont icon-loading animate-spin"></i>
            <span>{{ generating ? '创作中...' : `立即生成 (${kelingPowerCost}算力)` }}</span>
          </button>
        </div>
      </div>
    </div>

    <!-- 作品列表 -->
    <div class="p-4">
      <h2 class="text-lg font-semibold text-gray-900 mb-4">我的作品</h2>
      <div class="space-y-4">
        <div v-for="item in currentList" :key="item.id" class="bg-white rounded-xl p-4 shadow-sm">
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
                      <i class="iconfont icon-video text-gray-400 text-xl"></i>
                    </div>
                  </template>
                </el-image>
                <!-- 视频播放按钮 -->
                <button
                  v-if="item.progress === 100"
                  @click="playVideo(item)"
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
                    {{ item.title || '未命名视频' }}
                  </h3>
                  <p class="text-gray-500 text-sm mt-1 line-clamp-2">
                    {{ item.prompt }}
                  </p>
                </div>
                <!-- 任务状态 -->
                <div v-if="item.progress < 100" class="flex items-center space-x-2 text-sm">
                  <div
                    v-if="item.progress === 101"
                    class="text-red-600 flex items-center space-x-1"
                  >
                    <i class="iconfont icon-warning"></i>
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
                  v-if="item.raw_data?.task_type"
                  class="px-2 py-1 text-xs bg-blue-100 text-blue-600 rounded-full"
                >
                  {{ item.raw_data.task_type }}
                </span>
                <span
                  v-if="item.raw_data?.model"
                  class="px-2 py-1 text-xs bg-green-100 text-green-600 rounded-full"
                >
                  {{ item.raw_data.model }}
                </span>
                <span
                  v-if="item.raw_data?.duration"
                  class="px-2 py-1 text-xs bg-yellow-100 text-yellow-600 rounded-full"
                >
                  {{ item.raw_data.duration }}秒
                </span>
              </div>
            </div>
          </div>

          <!-- 操作按钮 -->
          <div class="flex items-center justify-between mt-4">
            <div class="flex space-x-2">
              <button
                v-if="item.progress === 100"
                @click="playVideo(item)"
                class="px-3 py-1.5 bg-blue-600 text-white text-sm rounded-lg hover:bg-blue-700 transition-colors flex items-center space-x-1"
              >
                <i class="iconfont icon-play !text-xs"></i>
                <span>播放</span>
              </button>
              <button
                v-if="item.progress === 100"
                @click="downloadVideo(item)"
                :disabled="item.downloading"
                class="px-3 py-1.5 bg-green-600 text-white text-sm rounded-lg hover:bg-green-700 transition-colors disabled:bg-gray-400 flex items-center space-x-1"
              >
                <i v-if="item.downloading" class="iconfont icon-loading animate-spin !text-xs"></i>
                <i v-else class="iconfont icon-download !text-xs"></i>
                <span>{{ item.downloading ? '下载中...' : '下载' }}</span>
              </button>
            </div>
            <button
              @click="removeJob(item)"
              class="px-3 py-1.5 bg-red-100 text-red-600 text-sm rounded-lg hover:bg-red-200 transition-colors flex items-center space-x-1"
            >
              <i class="iconfont icon-remove !text-xs"></i>
              <span>删除</span>
            </button>
          </div>
        </div>

        <!-- 加载更多 -->
        <div v-if="listLoading" class="flex justify-center py-4">
          <i class="iconfont icon-loading animate-spin text-blue-500 text-xl"></i>
        </div>

        <!-- 没有更多了 -->
        <div v-if="listFinished && !listLoading" class="text-center py-4 text-gray-500">
          没有更多了
        </div>
      </div>
    </div>

    <!-- 视频预览弹窗 -->
    <div
      v-if="showVideoDialog"
      class="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50"
      @click="showVideoDialog = false"
    >
      <div @click.stop class="bg-white rounded-2xl w-full max-w-4xl max-h-[80vh] animate-scale-up">
        <div class="flex items-center justify-between p-4 border-b">
          <h3 class="text-lg font-semibold text-gray-900">视频预览</h3>
          <button @click="showVideoDialog = false" class="p-2 hover:bg-gray-100 rounded-full">
            <i class="iconfont icon-close text-gray-500"></i>
          </button>
        </div>
        <div class="p-6">
          <video
            v-if="currentVideoUrl"
            :src="currentVideoUrl"
            controls
            autoplay
            class="w-full max-h-[60vh] rounded-lg"
          >
            您的浏览器不支持视频播放
          </video>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { showConfirmDialog } from 'vant'
import { httpGet, httpPost } from '@/utils/http'
import { checkSession } from '@/store/cache'
import CustomSelect from '@/views/mobile/components/CustomSelect.vue'
import { showMessageOK, showMessageError, showLoading, closeLoading } from '@/utils/dialog'

const router = useRouter()

// 响应式数据
const activeVideoType = ref('luma')
const loading = ref(false)
const generating = ref(false)
const isGenerating = ref(false)
const listLoading = ref(false)
const listFinished = ref(false)
const currentList = ref([])
const showVideoDialog = ref(false)
const currentVideoUrl = ref('')

// Luma 参数
const lumaParams = ref({
  prompt: '',
  image: '',
  image_tail: '',
  loop: false,
  expand_prompt: false,
})
const lumaUseImageMode = ref(false)
const lumaStartImage = ref([])
const lumaEndImage = ref([])

// KeLing 参数
const kelingParams = ref({
  aspect_ratio: '16:9',
  model: 'v1.5',
  duration: '5',
  mode: 'std',
  cfg_scale: 0.5,
  prompt: '',
  negative_prompt: '',
  image: '',
  image_tail: '',
  camera_control: {
    type: '',
    config: {
      horizontal: 0,
      vertical: 0,
      pan: 0,
      tilt: 0,
      roll: 0,
      zoom: 0,
    },
  },
})
const kelingUseImageMode = ref(false)
const kelingStartImage = ref([])
const kelingEndImage = ref([])

// 选项数据
const aspectRatioOptions = ['16:9', '9:16', '1:1', '4:3']
const modelOptions = ['v1.0', 'v1.5']
const durationOptions = ['5', '10']
const modeOptions = ['std', 'pro']
const cameraControlOptions = [
  '',
  'simple',
  'down_back',
  'forward_up',
  'right_turn_forward',
  'left_turn_forward',
]

// 获取运镜控制标签
const getCameraControlLabel = (option) => {
  const labelMap = {
    '': '请选择',
    simple: '简单运镜',
    down_back: '下移拉远',
    forward_up: '推进上移',
    right_turn_forward: '右旋推进',
    left_turn_forward: '左旋推进',
  }
  return labelMap[option] || option
}

// 页面数据
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)
const availablePower = ref(0)
const lumaPowerCost = ref(0)
const kelingPowerCost = ref(0)
const taskPulling = ref(true)
const tastPullHandler = ref(null)

// 生命周期
onMounted(() => {
  checkSession()
    .then(() => {
      fetchData(1)
      fetchUserPower()
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

const switchVideoType = (type) => {
  activeVideoType.value = type
  onVideoTypeChange(type)
}

const handleLumaStartImageUpload = (e) => {
  if (e.target.files[0]) {
    uploadLumaStartImage({ file: e.target.files[0], name: e.target.files[0].name })
  }
}

const handleLumaEndImageUpload = (e) => {
  if (e.target.files[0]) {
    uploadLumaEndImage({ file: e.target.files[0], name: e.target.files[0].name })
  }
}

const handleKelingStartImageUpload = (e) => {
  if (e.target.files[0]) {
    uploadKelingStartImage({ file: e.target.files[0], name: e.target.files[0].name })
  }
}

const handleKelingEndImageUpload = (e) => {
  if (e.target.files[0]) {
    uploadKelingEndImage({ file: e.target.files[0], name: e.target.files[0].name })
  }
}

const generatePrompt = () => {
  isGenerating.value = true
  // TODO: 实现提示词生成逻辑
  setTimeout(() => {
    isGenerating.value = false
    showMessageSuccess('提示词生成功能开发中')
  }, 2000)
}

const toggleLumaImageMode = () => {
  if (!lumaUseImageMode.value) {
    lumaParams.value.image = ''
    lumaParams.value.image_tail = ''
    lumaStartImage.value = []
    lumaEndImage.value = []
  }
}

const toggleKelingImageMode = () => {
  if (!kelingUseImageMode.value) {
    kelingParams.value.image = ''
    kelingParams.value.image_tail = ''
    kelingStartImage.value = []
    kelingEndImage.value = []
  }
}

const uploadLumaStartImage = (file) => {
  uploadImage(file, (url) => {
    lumaParams.value.image = url
  })
}

const uploadLumaEndImage = (file) => {
  uploadImage(file, (url) => {
    lumaParams.value.image_tail = url
  })
}

const uploadKelingStartImage = (file) => {
  uploadImage(file, (url) => {
    kelingParams.value.image = url
  })
}

const uploadKelingEndImage = (file) => {
  uploadImage(file, (url) => {
    kelingParams.value.image_tail = url
  })
}

const uploadImage = (file, callback) => {
  const formData = new FormData()
  formData.append('file', file.file, file.name)
  showLoading('正在上传图片...')

  httpPost('/api/upload', formData)
    .then((res) => {
      callback(res.data.url)
      showMessageOK('图片上传成功')
    })
    .catch((e) => {
      showMessageError('图片上传失败:' + e.message)
    })
    .finally(() => {
      closeLoading()
    })
}

const createLumaVideo = () => {
  if (!lumaParams.value.prompt.trim()) {
    showMessageError('请输入视频提示词')
    return
  }

  generating.value = true
  const params = {
    ...lumaParams.value,
    task_type: 'luma',
  }

  httpPost('/api/video/create', params)
    .then(() => {
      fetchData(1)
      taskPulling.value = true
      showMessageOK('创建任务成功')
    })
    .catch((e) => {
      showMessageError('创建任务失败：' + e.message)
    })
    .finally(() => {
      generating.value = false
    })
}

const createKelingVideo = () => {
  if (!kelingParams.value.prompt.trim()) {
    showMessageError('请输入视频提示词')
    return
  }

  generating.value = true
  const params = {
    ...kelingParams.value,
    task_type: 'keling',
  }

  httpPost('/api/video/create', params)
    .then(() => {
      fetchData(1)
      taskPulling.value = true
      showMessageOK('创建任务成功')
    })
    .catch((e) => {
      showMessageError('创建任务失败：' + e.message)
    })
    .finally(() => {
      generating.value = false
    })
}

const fetchData = (_page) => {
  if (_page) {
    page.value = _page
  }
  listLoading.value = true
  httpGet('/api/video/list', { page: page.value, page_size: pageSize.value })
    .then((res) => {
      total.value = res.data.total
      let needPull = false
      const items = []
      for (let v of res.data.items) {
        if (v.progress === 0 || v.progress === 102) {
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

const fetchUserPower = () => {
  httpGet('/api/user/power')
    .then((res) => {
      availablePower.value = res.data.power || 0
      lumaPowerCost.value = 10 // 示例值
      kelingPowerCost.value = 15 // 示例值
    })
    .catch(() => {})
}

const loadMore = () => {
  page.value++
  fetchData()
}

const playVideo = (item) => {
  currentVideoUrl.value = item.video_url
  showVideoDialog.value = true
}

const downloadVideo = (item) => {
  item.downloading = true
  const link = document.createElement('a')
  link.href = item.video_url
  link.download = item.title || 'video.mp4'
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
  item.downloading = false
  showMessageOK('开始下载')
}

const removeJob = (item) => {
  showConfirmDialog({
    title: '确认删除',
    message: '此操作将会删除任务相关文件，继续操作吗?',
    confirmButtonText: '确认删除',
    cancelButtonText: '取消',
  })
    .then(() => {
      httpGet('/api/video/remove', { id: item.id })
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
</script>

<style scoped>
/* 自定义动画 */
@keyframes fade-in {
  from {
    opacity: 0;
    transform: translateY(-10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes fade-out {
  from {
    opacity: 1;
    transform: translateY(0);
  }
  to {
    opacity: 0;
    transform: translateY(-10px);
  }
}

@keyframes scale-up {
  from {
    opacity: 0;
    transform: scale(0.9);
  }
  to {
    opacity: 1;
    transform: scale(1);
  }
}

.animate-fade-in {
  animation: fade-in 0.3s ease-out;
}

.animate-fade-out {
  animation: fade-out 0.3s ease-out;
}

.animate-scale-up {
  animation: scale-up 0.3s ease-out;
}

/* 文本截断 */
.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

/* 深色模式适配 */
@media (prefers-color-scheme: dark) {
  .bg-gray-50 {
    background-color: #1f2937;
  }

  .bg-white {
    background-color: #374151;
  }

  .text-gray-900 {
    color: #f9fafb;
  }

  .text-gray-700 {
    color: #d1d5db;
  }

  .text-gray-600 {
    color: #9ca3af;
  }

  .text-gray-500 {
    color: #6b7280;
  }

  .border-gray-200 {
    border-color: #4b5563;
  }

  .bg-gray-100:hover {
    background-color: #4b5563;
  }
}
</style>
