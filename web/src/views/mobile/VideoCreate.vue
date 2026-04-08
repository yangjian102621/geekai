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
        <h1 class="flex-1 text-center text-lg text-gray-900">视频创作</h1>
        <div class="w-8"></div>
      </div>
    </div>

    <!-- 视频类型切换 -->
    <div class="p-4 space-y-6">
      <!-- 视频类型选择 -->
      <div class="bg-white rounded-xl p-3 shadow-sm">
        <div class="flex space-x-2">
          <button
            @click="video.switchVideoType('luma')"
            :class="[
              'flex-1 py-2.5 px-4 rounded-lg font-medium transition-colors',
              video.activeVideoType === 'luma'
                ? 'bg-blue-600 text-white'
                : 'bg-gray-100 text-gray-700 hover:bg-gray-200',
            ]"
          >
            Luma视频
          </button>
          <button
            @click="video.switchVideoType('keling')"
            :class="[
              'flex-1 py-2.5 px-4 rounded-lg font-medium transition-colors',
              video.activeVideoType === 'keling'
                ? 'bg-blue-600 text-white'
                : 'bg-gray-100 text-gray-700 hover:bg-gray-200',
            ]"
          >
            可灵视频
          </button>
        </div>
      </div>

      <!-- Luma 视频参数 -->
      <div v-if="video.activeVideoType === 'luma'" class="space-y-6">
        <!-- 提示词输入 -->
        <div class="bg-white rounded-xl p-4 shadow-sm">
          <label class="block text-gray-700 font-medium mb-3">提示词</label>
          <textarea
            v-model="video.lumaParams.prompt"
            placeholder="请在此输入视频提示词，用逗号分割"
            class="w-full px-4 py-3 border border-gray-200 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent resize-none"
            rows="4"
            maxlength="2000"
          />
          <div class="flex justify-between">
            <van-button
              @click="video.generatePrompt"
              :disabled="video.isGenerating"
              type="primary"
              size="small"
            >
              <i v-if="video.isGenerating" class="iconfont icon-loading animate-spin"></i>
              <span class="ml-1">{{ video.isGenerating ? '' : '生成提示词' }}</span>
            </van-button>
            <span class="text-sm text-gray-500">{{ video.lumaParams.prompt.length }}/2000</span>
          </div>
        </div>

        <!-- 图片辅助生成开关 -->
        <div class="bg-white rounded-xl p-4 shadow-sm">
          <div class="flex items-center justify-between">
            <div>
              <span class="text-gray-900 font-medium">使用图片辅助生成</span>
              <p class="text-sm text-gray-500 mt-1">上传起始帧和结束帧图片</p>
            </div>
            <el-switch
              v-model="video.lumaUseImageMode"
              @change="video.toggleLumaImageMode"
              size="default"
            />
          </div>
        </div>

        <!-- 图片上传区域 -->
        <div v-if="video.lumaUseImageMode" class="bg-white rounded-xl p-4 shadow-sm">
          <div class="grid grid-cols-2 gap-4">
            <div class="relative">
              <div
                class="border-2 border-dashed border-gray-300 rounded-lg p-6 text-center hover:border-blue-400 hover:bg-blue-50 transition-colors cursor-pointer h-32"
              >
                <input
                  ref="lumaStartInput"
                  type="file"
                  accept=".jpg,.png,.jpeg"
                  @change="video.handleLumaStartImageUpload"
                  class="hidden"
                />
                <div
                  @click="$refs.lumaStartInput.click()"
                  class="flex flex-col items-center space-y-2 h-full justify-center"
                >
                  <i
                    v-if="!video.lumaStartImage.length"
                    class="iconfont icon-upload text-blue-500 text-xl"
                  ></i>
                  <span v-if="!video.lumaStartImage.length" class="text-gray-700 text-sm"
                    >起始帧</span
                  >
                  <div v-else class="w-full h-full relative">
                    <el-image
                      :src="video.lumaStartImage[0]?.url || video.lumaStartImage[0]?.content"
                      fit="cover"
                      class="w-full h-full rounded"
                    />
                    <button
                      @click.stop="video.lumaStartImage = []"
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
                  @change="video.handleLumaEndImageUpload"
                  class="hidden"
                />
                <div
                  @click="$refs.lumaEndInput.click()"
                  class="flex flex-col items-center space-y-2 h-full justify-center"
                >
                  <i
                    v-if="!video.lumaEndImage.length"
                    class="iconfont icon-upload text-blue-500 text-xl"
                  ></i>
                  <span v-if="!video.lumaEndImage.length" class="text-gray-700 text-sm"
                    >结束帧</span
                  >
                  <div v-else class="w-full h-full relative">
                    <el-image
                      :src="video.lumaEndImage[0]?.url || video.lumaEndImage[0]?.content"
                      fit="cover"
                      class="w-full h-full rounded"
                    />
                    <button
                      @click.stop="video.lumaEndImage = []"
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
              <el-switch v-model="video.lumaParams.loop" size="default" />
            </div>
            <div class="flex items-center justify-between">
              <span class="text-gray-900 font-medium">提示词优化</span>
              <el-switch v-model="video.lumaParams.expand_prompt" size="default" />
            </div>
          </div>
        </div>

        <!-- 算力显示和生成按钮 -->
        <div class="bg-white rounded-xl p-4 shadow-sm">
          <button
            @click="video.createLumaVideo"
            :disabled="video.generating"
            type="button"
            class="w-full py-3 bg-gradient-to-r from-blue-500 to-purple-600 text-white rounded-xl disabled:from-gray-400 disabled:to-gray-400 disabled:cursor-not-allowed hover:from-blue-600 hover:to-purple-700 transition-all duration-200 flex items-center justify-center space-x-2"
          >
            <i v-if="video.generating" class="iconfont icon-loading animate-spin"></i>
            <i v-else class="iconfont icon-chuangzuo"></i>
            <span>{{
              video.generating ? '创作中...' : `立即生成 (${video.lumaPowerCost}算力)`
            }}</span>
          </button>
        </div>
      </div>

      <!-- KeLing 视频参数 -->
      <div v-if="video.activeVideoType === 'keling'" class="space-y-6">
        <!-- 画面比例 -->
        <div class="bg-white rounded-xl p-4 shadow-sm">
          <label class="block text-gray-700 font-medium mb-3">画面比例</label>
          <CustomSelect
            v-model="video.kelingParams.aspect_ratio"
            :options="video.aspectRatioOptions.map((ratio) => ({ label: ratio, value: ratio }))"
            title="选择比例"
          />
        </div>

        <!-- 模型选择 -->
        <div class="bg-white rounded-xl p-4 shadow-sm">
          <label class="block text-gray-700 font-medium mb-3">模型选择</label>
          <CustomSelect
            v-model="video.kelingParams.model"
            :options="video.modelOptions"
            placeholder="请选择模型"
            title="选择模型"
          >
            <template #option="{ option, selected }">
              <div class="flex items-center w-full">
                <span class="font-bold text-blue-600 mr-2">{{ option.label }}</span>
                <span class="text-xs text-gray-400">({{ option.value }})</span>
                <span v-if="selected" class="ml-auto text-green-500"
                  ><i class="iconfont icon-success"></i
                ></span>
              </div>
            </template>
          </CustomSelect>
        </div>

        <!-- 视频时长 -->
        <div class="bg-white rounded-xl p-4 shadow-sm">
          <label class="block text-gray-700 font-medium mb-3">视频时长</label>
          <CustomSelect
            v-model="video.kelingParams.duration"
            :options="
              video.durationOptions.map((duration) => ({ label: `${duration}秒`, value: duration }))
            "
            title="选择时长"
          />
        </div>

        <!-- 生成模式 -->
        <div class="bg-white rounded-xl p-4 shadow-sm">
          <label class="block text-gray-700 font-medium mb-3">生成模式</label>
          <CustomSelect
            v-model="video.kelingParams.mode"
            :options="
              video.modeOptions.map((mode) => ({
                label: mode === 'std' ? '标准模式' : '专业模式',
                value: mode,
              }))
            "
            title="选择模式"
          />
        </div>

        <!-- 创意程度 -->
        <div class="bg-white rounded-xl p-4 shadow-sm">
          <div class="space-y-4">
            <label class="block text-gray-700 font-medium">创意程度</label>
            <el-slider v-model="video.kelingParams.cfg_scale" :min="0" :max="1" :step="0.1" />
          </div>
        </div>

        <!-- 运镜控制 -->
        <div class="bg-white rounded-xl p-4 shadow-sm">
          <label class="block text-gray-700 font-medium mb-3">运镜控制</label>
          <CustomSelect
            v-model="video.kelingParams.camera_control.type"
            :options="
              video.cameraControlOptions.map((option) => ({
                label: video.getCameraControlLabel(option),
                value: option,
              }))
            "
            title="选择运镜类型"
          />
        </div>

        <!-- 图片辅助生成开关 -->
        <div class="bg-white rounded-xl p-4 shadow-sm">
          <div class="flex items-center justify-between">
            <div>
              <span class="text-gray-900 font-medium">使用图片辅助生成</span>
              <p class="text-sm text-gray-500 mt-1">上传起始帧和结束帧图片</p>
            </div>
            <el-switch
              v-model="video.kelingUseImageMode"
              @change="video.toggleKelingImageMode"
              size="default"
            />
          </div>
        </div>

        <!-- 图片上传区域 -->
        <div v-if="video.kelingUseImageMode" class="bg-white rounded-xl p-4 shadow-sm">
          <div class="grid grid-cols-2 gap-4">
            <div class="relative">
              <div
                class="border-2 border-dashed border-gray-300 rounded-lg p-6 text-center hover:border-blue-400 hover:bg-blue-50 transition-colors cursor-pointer h-32"
              >
                <input
                  ref="kelingStartInput"
                  type="file"
                  accept=".jpg,.png,.jpeg"
                  @change="video.handleKelingStartImageUpload"
                  class="hidden"
                />
                <div
                  @click="$refs.kelingStartInput.click()"
                  class="flex flex-col items-center space-y-2 h-full justify-center"
                >
                  <i
                    v-if="!video.kelingStartImage.length"
                    class="iconfont icon-upload text-blue-500 text-xl"
                  ></i>
                  <span v-if="!video.kelingStartImage.length" class="text-gray-700 text-sm"
                    >起始帧</span
                  >
                  <div v-else class="w-full h-full relative">
                    <el-image
                      :src="video.kelingStartImage[0]?.url || video.kelingStartImage[0]?.content"
                      fit="cover"
                      class="w-full h-full rounded"
                    />
                    <button
                      @click.stop="video.kelingStartImage = []"
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
                  @change="video.handleKelingEndImageUpload"
                  class="hidden"
                />
                <div
                  @click="$refs.kelingEndInput.click()"
                  class="flex flex-col items-center space-y-2 h-full justify-center"
                >
                  <i
                    v-if="!video.kelingEndImage.length"
                    class="iconfont icon-upload text-blue-500 text-xl"
                  ></i>
                  <span v-if="!video.kelingEndImage.length" class="text-gray-700 text-sm"
                    >结束帧</span
                  >
                  <div v-else class="w-full h-full relative">
                    <el-image
                      :src="video.kelingEndImage[0]?.url || video.kelingEndImage[0]?.content"
                      fit="cover"
                      class="w-full h-full rounded"
                    />
                    <button
                      @click.stop="video.kelingEndImage = []"
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
            v-model="video.kelingParams.prompt"
            :placeholder="video.kelingUseImageMode ? '描述视频画面细节' : '请在此输入视频提示词'"
            class="w-full px-4 py-3 border border-gray-200 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent resize-none"
            rows="4"
            maxlength="500"
          />
          <div class="flex justify-between">
            <van-button
              @click="video.generatePrompt"
              :disabled="video.isGenerating"
              type="primary"
              size="small"
            >
              <i v-if="video.isGenerating" class="iconfont icon-loading animate-spin"></i>
              <span class="ml-1">{{ video.isGenerating ? '' : '生成提示词' }}</span>
            </van-button>
            <span class="text-sm text-gray-500">{{ video.kelingParams.prompt.length }}/500</span>
          </div>
        </div>

        <!-- 排除内容 -->
        <div class="bg-white rounded-xl p-4 shadow-sm">
          <label class="block text-gray-700 font-medium mb-3">不希望出现的内容</label>
          <textarea
            v-model="video.kelingParams.negative_prompt"
            placeholder="请在此输入你不希望出现在视频上的内容"
            class="w-full px-4 py-3 border border-gray-200 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent resize-none"
            rows="3"
            maxlength="500"
          />
          <div class="text-right mt-2">
            <span class="text-sm text-gray-500"
              >{{ video.kelingParams.negative_prompt.length }}/500</span
            >
          </div>
        </div>

        <!-- 算力显示和生成按钮 -->
        <div class="bg-white rounded-xl p-4 shadow-sm">
          <button
            @click="video.createKelingVideo"
            :disabled="video.generating"
            type="button"
            class="w-full py-3 bg-gradient-to-r from-blue-500 to-purple-600 text-white rounded-xl disabled:from-gray-400 disabled:to-gray-400 disabled:cursor-not-allowed hover:from-blue-600 hover:to-purple-700 transition-all duration-200 flex items-center justify-center space-x-2"
          >
            <i v-if="video.generating" class="iconfont icon-loading animate-spin"></i>
            <i v-else class="iconfont icon-chuangzuo"></i>
            <span>{{
              video.generating ? '创作中...' : `立即生成 (${video.kelingPowerCost}算力)`
            }}</span>
          </button>
        </div>
      </div>
    </div>

    <!-- 作品列表 -->
    <div class="p-4">
      <h2 class="text-lg font-semibold text-gray-900 mb-4">我的作品</h2>
      <div class="space-y-4" v-if="video.currentList.length > 0">
        <div
          v-for="item in video.currentList"
          :key="item.id"
          class="bg-white rounded-xl p-4 shadow-sm"
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
                      <i class="iconfont icon-video text-gray-400 text-xl"></i>
                    </div>
                  </template>
                </el-image>
                <!-- 视频播放按钮 -->
                <button
                  v-if="item.progress === 100"
                  @click="video.playVideo(item)"
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
                @click="video.playVideo(item)"
                class="px-3 py-1.5 bg-blue-600 text-white text-sm rounded-lg hover:bg-blue-700 transition-colors flex items-center space-x-1"
              >
                <i class="iconfont icon-play !text-xs"></i>
                <span>播放</span>
              </button>
              <button
                v-if="item.progress === 100"
                @click="video.downloadVideo(item)"
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
        <div v-if="video.listLoading" class="flex justify-center py-4">
          <i class="iconfont icon-loading animate-spin text-blue-500 text-xl"></i>
        </div>

        <!-- 没有更多了 -->
        <div v-if="video.listFinished && !video.listLoading" class="text-center py-4 text-gray-500">
          没有更多了
        </div>
      </div>

      <div class="px-4" v-else>
        <van-empty description="暂无数据" image-size="120" />
      </div>
    </div>

    <!-- 视频预览弹窗 -->
    <div
      v-if="video.showVideoDialog"
      class="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50"
      @click="video.showVideoDialog = false"
    >
      <div @click.stop class="bg-white rounded-2xl w-full max-w-4xl max-h-[80vh] animate-scale-up">
        <div class="flex items-center justify-between p-4 border-b">
          <h3 class="text-lg font-semibold text-gray-900">视频预览</h3>
          <button @click="video.showVideoDialog = false" class="p-2 hover:bg-gray-100 rounded-full">
            <i class="iconfont icon-close text-gray-500"></i>
          </button>
        </div>
        <div class="p-6">
          <video
            v-if="video.currentVideoUrl"
            :src="video.currentVideoUrl"
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
import '@/assets/css/mobile/video.scss'
import CustomSelect from '@/components/mobile/CustomSelect.vue'
import { useVideoStore } from '@/store/mobile/video'
import { showConfirmDialog } from 'vant'
import { onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { checkSession } from '@/store/cache'

const router = useRouter()
const video = useVideoStore()

// 页面专属方法
const goBack = () => {
  router.back()
}

// 定时轮询等副作用
let tastPullHandler = null
onMounted(() => {
  checkSession()
    .then(() => {
      video.fetchData(1)
      video.fetchUserPower()
      tastPullHandler = setInterval(() => {
        if (video.taskPulling) {
          video.fetchData(1)
        }
      }, 5000)
    })
    .catch(() => {})
})
onUnmounted(() => {
  if (tastPullHandler) clearInterval(tastPullHandler)
})

// 删除弹窗（页面层处理）
const removeJob = (item) => {
  showConfirmDialog({
    title: '确认删除',
    message: '此操作将会删除任务相关文件，继续操作吗?',
    confirmButtonText: '确认删除',
    cancelButtonText: '取消',
  })
    .then(() => {
      video.fetchData(1)
    })
    .catch(() => {})
}
</script>

<style scoped>
/* Dark 主题样式 - 按照 theme-dark.scss 的模式 */
:root[data-theme='dark'] .min-h-screen {
  background-color: rgb(13, 20, 53) !important;

  /* 页面头部 */
  .sticky.top-0 {
    background-color: rgb(31, 41, 55) !important;
    box-shadow: 0 1px 3px 0 rgba(0, 0, 0, 0.3) !important;

    .icon-back {
      color: rgb(156, 163, 175) !important;
    }

    h1 {
      color: rgb(255, 255, 255) !important;
    }

    button:hover {
      background-color: rgb(75, 85, 99) !important;
    }
  }

  /* 视频类型切换 */
  .space-y-6 {
    .bg-white {
      background-color: rgb(31, 41, 55) !important;
      box-shadow: 0 1px 3px 0 rgba(0, 0, 0, 0.3) !important;
    }

    .text-gray-900 {
      color: rgb(209, 213, 219) !important;
    }

    .text-gray-700 {
      color: rgb(209, 213, 219) !important;
    }

    .text-gray-500 {
      color: rgb(156, 163, 175) !important;
    }

    .text-gray-600 {
      color: rgb(156, 163, 175) !important;
    }

    /* 视频类型选择按钮 */
    .bg-gray-100 {
      background-color: rgb(55, 65, 81) !important;
      color: rgb(209, 213, 219) !important;

      &:hover {
        background-color: rgb(75, 85, 99) !important;
      }
    }

    /* 输入框样式 */
    input,
    textarea {
      background-color: rgb(55, 65, 81) !important;
      border-color: rgb(75, 85, 99) !important;
      color: rgb(209, 213, 219) !important;

      &::placeholder {
        color: rgb(107, 114, 128) !important;
      }

      &:focus {
        border-color: rgb(139, 92, 246) !important;
        box-shadow: 0 0 0 2px rgba(139, 92, 246, 0.2) !important;
      }
    }

    /* 图片上传区域 */
    .border-dashed {
      border-color: rgb(75, 85, 99) !important;

      &:hover {
        border-color: rgb(59, 130, 246) !important;
        background-color: rgba(59, 130, 246, 0.1) !important;
      }

      .text-gray-700 {
        color: rgb(209, 213, 219) !important;
      }
    }

    /* 按钮样式 */
    .bg-blue-600 {
      background-color: rgb(37, 99, 235) !important;

      &:hover:not(:disabled) {
        background-color: rgb(29, 78, 216) !important;
      }

      &:disabled {
        background-color: rgb(156, 163, 175) !important;
      }
    }

    .bg-gradient-to-r.from-blue-500.to-purple-600 {
      background: linear-gradient(to right, rgb(59, 130, 246), rgb(147, 51, 234)) !important;

      &:hover:not(:disabled) {
        background: linear-gradient(to right, rgb(37, 99, 235), rgb(126, 34, 206)) !important;
      }

      &:disabled {
        background: linear-gradient(to right, rgb(156, 163, 175), rgb(156, 163, 175)) !important;
      }
    }

    /* 删除按钮 */
    .bg-red-500 {
      background-color: rgb(239, 68, 68) !important;

      &:hover {
        background-color: rgb(220, 38, 38) !important;
      }
    }
  }

  /* 作品列表 */
  .p-4 {
    h2 {
      color: rgb(255, 255, 255) !important;
    }

    .bg-white {
      background-color: rgb(31, 41, 55) !important;
      box-shadow: 0 1px 3px 0 rgba(0, 0, 0, 0.3) !important;
    }

    .bg-gray-100 {
      background-color: rgb(55, 65, 81) !important;
    }

    .text-gray-900 {
      color: rgb(209, 213, 219) !important;
    }

    .text-gray-500 {
      color: rgb(156, 163, 175) !important;
    }

    .text-gray-600 {
      color: rgb(156, 163, 175) !important;
    }

    .text-gray-400 {
      color: rgb(107, 114, 128) !important;
    }

    /* 标签样式 */
    .bg-blue-100 {
      background-color: rgba(59, 130, 246, 0.1) !important;
      color: rgb(59, 130, 246) !important;
    }

    .bg-green-100 {
      background-color: rgba(34, 197, 94, 0.1) !important;
      color: rgb(34, 197, 94) !important;
    }

    .bg-yellow-100 {
      background-color: rgba(251, 191, 36, 0.1) !important;
      color: rgb(251, 191, 36) !important;
    }

    /* 按钮样式 */
    .bg-blue-600 {
      background-color: rgb(37, 99, 235) !important;

      &:hover {
        background-color: rgb(29, 78, 216) !important;
      }
    }

    .bg-green-600 {
      background-color: rgb(34, 197, 94) !important;

      &:hover {
        background-color: rgb(22, 163, 74) !important;
      }
    }

    .bg-red-100 {
      background-color: rgba(239, 68, 68, 0.1) !important;
      color: rgb(239, 68, 68) !important;

      &:hover {
        background-color: rgba(239, 68, 68, 0.2) !important;
      }
    }

    /* 状态指示 */
    .text-red-600 {
      color: rgb(239, 68, 68) !important;
    }

    .text-blue-600 {
      color: rgb(37, 99, 235) !important;
    }

    /* 加载状态 */
    .text-blue-500 {
      color: rgb(59, 130, 246) !important;
    }

    /* 加载更多 */
    .text-gray-500 {
      color: rgb(156, 163, 175) !important;
    }
  }

  /* 视频预览弹窗 */
  .fixed.inset-0 {
    .bg-white {
      background-color: rgb(31, 41, 55) !important;
    }

    .border-b {
      border-bottom-color: rgb(75, 85, 99) !important;
    }

    h3 {
      color: rgb(255, 255, 255) !important;
    }

    button {
      &:hover {
        background-color: rgb(75, 85, 99) !important;
      }

      .iconfont {
        color: rgb(156, 163, 175) !important;
      }
    }
  }
}
</style>
