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
        <h1 class="flex-1 text-center text-lg font-semibold text-gray-900">即梦AI</h1>
        <div class="w-8"></div>
      </div>
    </div>

    <!-- 功能分类选择 -->
    <div class="p-4 space-y-6">
      <!-- 功能类型选择 -->
      <div class="bg-white rounded-xl p-4 shadow-sm">
        <div class="grid grid-cols-2 gap-3">
          <button
            v-for="category in categories"
            :key="category.key"
            @click="switchCategory(category.key)"
            :class="[
              'flex flex-col items-center p-3 rounded-lg border-2 transition-colors',
              activeCategory === category.key
                ? 'border-blue-500 bg-blue-50 text-blue-700'
                : 'border-gray-200 bg-gray-50 text-gray-600 hover:border-gray-300 hover:bg-gray-100',
            ]"
          >
            <i :class="getCategoryIcon(category.key)" class="text-2xl mb-2"></i>
            <span class="text-sm font-medium">{{ category.name }}</span>
          </button>
        </div>
      </div>
      <!-- 生成模式切换 -->
      <div
        v-if="activeCategory === 'image_generation' || activeCategory === 'video_generation'"
        class="bg-white rounded-xl p-4 shadow-sm"
      >
        <div class="flex items-center justify-between">
          <div>
            <span class="text-gray-900 font-medium">生成模式</span>
            <p class="text-sm text-gray-500 mt-1">
              {{ activeCategory === 'image_generation' ? '图生图人像写真' : '图生视频' }}
            </p>
          </div>
          <el-switch v-model="useImageInput" @change="switchInputMode" size="default" />
        </div>
      </div>

      <!-- 文生图 -->
      <div v-if="activeFunction === 'text_to_image'" class="space-y-6">
        <!-- 提示词输入 -->
        <div class="bg-white rounded-xl p-4 shadow-sm">
          <label class="block text-gray-700 font-medium mb-3">提示词</label>
          <textarea
            v-model="currentPrompt"
            placeholder="请输入图片描述，越详细越好"
            class="w-full px-4 py-3 border border-gray-200 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent resize-none"
            rows="4"
            maxlength="2000"
          />
          <div class="text-right mt-2">
            <span class="text-sm text-gray-500">{{ currentPrompt.length }}/2000</span>
          </div>
        </div>

        <!-- 图片尺寸 -->
        <CustomSelect
          v-model="textToImageParams.size"
          :options="imageSizeOptions.map((opt) => ({ label: opt.label, value: opt.value }))"
          label="图片尺寸"
          title="选择尺寸"
        />

        <!-- 创意度 -->
        <div class="bg-white rounded-xl p-4 shadow-sm">
          <div class="space-y-4">
            <div class="flex items-center space-x-2">
              <label class="block text-gray-700 font-medium">创意度</label>
              <el-tooltip content="创意度越高，影响文本描述的程度越高" placement="top">
                <i class="iconfont icon-info text-gray-400 cursor-pointer"></i>
              </el-tooltip>
            </div>
            <el-slider v-model="textToImageParams.scale" :min="1" :max="10" :step="0.5" />
          </div>
        </div>

        <!-- 智能优化提示词 -->
        <div class="bg-white rounded-xl p-4 shadow-sm">
          <div class="flex items-center justify-between">
            <span class="text-gray-900 font-medium">智能优化提示词</span>
            <el-switch v-model="textToImageParams.use_pre_llm" size="default" />
          </div>
        </div>
      </div>

      <!-- 图生图 -->
      <div v-if="activeFunction === 'image_to_image'" class="space-y-6">
        <!-- 上传图片 -->
        <div class="bg-white rounded-xl p-4 shadow-sm">
          <label class="block text-gray-700 font-medium mb-3">上传图片</label>
          <div
            class="border-2 border-dashed border-gray-300 rounded-lg p-6 text-center hover:border-blue-400 hover:bg-blue-50 transition-colors cursor-pointer"
          >
            <input
              ref="imageToImageInput"
              type="file"
              accept=".jpg,.png,.jpeg"
              @change="
                (e) => onImageUpload({ file: e.target.files[0], name: e.target.files[0]?.name })
              "
              class="hidden"
            />
            <div
              @click="$refs.imageToImageInput?.click()"
              class="flex flex-col items-center space-y-2"
            >
              <i
                v-if="!imageToImageParams.image_input.length"
                class="iconfont icon-upload text-blue-500 text-2xl"
              ></i>
              <span v-if="!imageToImageParams.image_input.length" class="text-gray-700 font-medium"
                >上传图片</span
              >
              <div v-else class="relative">
                <el-image
                  :src="
                    imageToImageParams.image_input[0]?.url ||
                    imageToImageParams.image_input[0]?.content
                  "
                  fit="cover"
                  class="w-32 h-32 rounded"
                />
                <button
                  @click.stop="imageToImageParams.image_input = []"
                  class="absolute -top-2 -right-2 w-6 h-6 bg-red-500 text-white rounded-full flex items-center justify-center text-xs hover:bg-red-600 transition-colors"
                >
                  <i class="iconfont icon-close"></i>
                </button>
              </div>
            </div>
          </div>
        </div>

        <!-- 提示词输入 -->
        <div class="bg-white rounded-xl p-4 shadow-sm">
          <label class="block text-gray-700 font-medium mb-3">提示词</label>
          <textarea
            v-model="currentPrompt"
            placeholder="描述你想要的图片效果"
            class="w-full px-4 py-3 border border-gray-200 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent resize-none"
            rows="4"
            maxlength="2000"
          />
          <div class="text-right mt-2">
            <span class="text-sm text-gray-500">{{ currentPrompt.length }}/2000</span>
          </div>
        </div>

        <!-- 图片尺寸 -->
        <CustomSelect
          v-model="imageToImageParams.size"
          :options="imageSizeOptions.map((opt) => ({ label: opt.label, value: opt.value }))"
          label="图片尺寸"
          title="选择尺寸"
        />
      </div>

      <!-- 图像编辑 -->
      <div v-if="activeFunction === 'image_edit'" class="space-y-6">
        <!-- 上传图片 -->
        <div class="bg-white rounded-xl p-4 shadow-sm">
          <label class="block text-gray-700 font-medium mb-3">上传图片</label>
          <div
            class="border-2 border-dashed border-gray-300 rounded-lg p-6 text-center hover:border-blue-400 hover:bg-blue-50 transition-colors cursor-pointer"
          >
            <input
              ref="imageEditInput"
              type="file"
              accept=".jpg,.png,.jpeg"
              @change="
                (e) => onImageUpload({ file: e.target.files[0], name: e.target.files[0]?.name })
              "
              class="hidden"
            />
            <div
              @click="$refs.imageEditInput?.click()"
              class="flex flex-col items-center space-y-2"
            >
              <i
                v-if="!imageEditParams.image_urls.length"
                class="iconfont icon-upload text-blue-500 text-2xl"
              ></i>
              <span v-if="!imageEditParams.image_urls.length" class="text-gray-700 font-medium"
                >上传图片</span
              >
              <div v-else class="relative">
                <el-image
                  :src="
                    imageEditParams.image_urls[0]?.url || imageEditParams.image_urls[0]?.content
                  "
                  fit="cover"
                  class="w-32 h-32 rounded"
                />
                <button
                  @click.stop="imageEditParams.image_urls = []"
                  class="absolute -top-2 -right-2 w-6 h-6 bg-red-500 text-white rounded-full flex items-center justify-center text-xs hover:bg-red-600 transition-colors"
                >
                  <i class="iconfont icon-close"></i>
                </button>
              </div>
            </div>
          </div>
        </div>

        <!-- 编辑提示词 -->
        <div class="bg-white rounded-xl p-4 shadow-sm">
          <label class="block text-gray-700 font-medium mb-3">编辑提示词</label>
          <textarea
            v-model="currentPrompt"
            placeholder="描述你想要的编辑效果"
            class="w-full px-4 py-3 border border-gray-200 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent resize-none"
            rows="4"
            maxlength="2000"
          />
          <div class="text-right mt-2">
            <span class="text-sm text-gray-500">{{ currentPrompt.length }}/2000</span>
          </div>
        </div>

        <!-- 编辑强度 -->
        <div class="bg-white rounded-xl p-4 shadow-sm">
          <div class="space-y-4">
            <label class="block text-gray-700 font-medium">编辑强度</label>
            <el-slider v-model="imageEditParams.scale" :min="0" :max="1" :step="0.1" />
          </div>
        </div>
      </div>

      <!-- 图像特效 -->
      <div v-if="activeFunction === 'image_effects'" class="space-y-6">
        <!-- 上传图片 -->
        <div class="bg-white rounded-xl p-4 shadow-sm">
          <label class="block text-gray-700 font-medium mb-3">上传图片</label>
          <div
            class="border-2 border-dashed border-gray-300 rounded-lg p-6 text-center hover:border-blue-400 hover:bg-blue-50 transition-colors cursor-pointer"
          >
            <input
              ref="imageEffectsInput"
              type="file"
              accept=".jpg,.png,.jpeg"
              @change="
                (e) => onImageUpload({ file: e.target.files[0], name: e.target.files[0]?.name })
              "
              class="hidden"
            />
            <div
              @click="$refs.imageEffectsInput?.click()"
              class="flex flex-col items-center space-y-2"
            >
              <i
                v-if="!imageEffectsParams.image_input1.length"
                class="iconfont icon-upload text-blue-500 text-2xl"
              ></i>
              <span v-if="!imageEffectsParams.image_input1.length" class="text-gray-700 font-medium"
                >上传图片</span
              >
              <div v-else class="relative">
                <el-image
                  :src="
                    imageEffectsParams.image_input1[0]?.url ||
                    imageEffectsParams.image_input1[0]?.content
                  "
                  fit="cover"
                  class="w-32 h-32 rounded"
                />
                <button
                  @click.stop="imageEffectsParams.image_input1 = []"
                  class="absolute -top-2 -right-2 w-6 h-6 bg-red-500 text-white rounded-full flex items-center justify-center text-xs hover:bg-red-600 transition-colors"
                >
                  <i class="iconfont icon-close"></i>
                </button>
              </div>
            </div>
          </div>
        </div>

        <!-- 特效模板 -->
        <CustomSelect
          v-model="imageEffectsParams.template_id"
          :options="
            imageEffectsTemplateOptions.map((opt) => ({ label: opt.label, value: opt.value }))
          "
          label="特效模板"
          title="选择特效模板"
        />

        <!-- 输出尺寸 -->
        <CustomSelect
          v-model="imageEffectsParams.size"
          :options="imageSizeOptions.map((opt) => ({ label: opt.label, value: opt.value }))"
          label="输出尺寸"
          title="选择尺寸"
        />
      </div>

      <!-- 文生视频 -->
      <div v-if="activeFunction === 'text_to_video'" class="space-y-6">
        <!-- 提示词输入 -->
        <div class="bg-white rounded-xl p-4 shadow-sm">
          <label class="block text-gray-700 font-medium mb-3">提示词</label>
          <textarea
            v-model="currentPrompt"
            placeholder="描述你想要的视频内容"
            class="w-full px-4 py-3 border border-gray-200 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent resize-none"
            rows="4"
            maxlength="2000"
          />
          <div class="text-right mt-2">
            <span class="text-sm text-gray-500">{{ currentPrompt.length }}/2000</span>
          </div>
        </div>

        <!-- 视频比例 -->
        <CustomSelect
          v-model="textToVideoParams.aspect_ratio"
          :options="videoAspectRatioOptions.map((opt) => ({ label: opt.label, value: opt.value }))"
          label="视频比例"
          title="选择比例"
        />
      </div>

      <!-- 图生视频 -->
      <div v-if="activeFunction === 'image_to_video'" class="space-y-6">
        <!-- 上传图片 -->
        <div class="bg-white rounded-xl p-4 shadow-sm">
          <label class="block text-gray-700 font-medium mb-3">上传图片（最多2张）</label>
          <div
            class="border-2 border-dashed border-gray-300 rounded-lg p-6 text-center hover:border-blue-400 hover:bg-blue-50 transition-colors cursor-pointer"
          >
            <input
              ref="imageToVideoInput"
              type="file"
              accept=".jpg,.png,.jpeg"
              multiple
              @change="(e) => handleMultipleImageUpload(e)"
              class="hidden"
            />
            <div
              @click="$refs.imageToVideoInput?.click()"
              class="flex flex-col items-center space-y-2"
            >
              <i
                v-if="!imageToVideoParams.image_urls.length"
                class="iconfont icon-upload text-blue-500 text-2xl"
              ></i>
              <span v-if="!imageToVideoParams.image_urls.length" class="text-gray-700 font-medium"
                >上传图片</span
              >
              <div v-else class="flex space-x-3">
                <div
                  v-for="(image, index) in imageToVideoParams.image_urls"
                  :key="index"
                  class="relative"
                >
                  <el-image
                    :src="image?.url || image?.content"
                    fit="cover"
                    class="w-24 h-24 rounded"
                  />
                  <button
                    @click.stop="removeImage(index)"
                    class="absolute -top-2 -right-2 w-6 h-6 bg-red-500 text-white rounded-full flex items-center justify-center text-xs hover:bg-red-600 transition-colors"
                  >
                    <i class="iconfont icon-close"></i>
                  </button>
                </div>
                <div
                  v-if="imageToVideoParams.image_urls.length < 2"
                  @click.stop="$refs.imageToVideoInput?.click()"
                  class="w-24 h-24 border-2 border-dashed border-gray-300 rounded flex items-center justify-center cursor-pointer hover:border-blue-400"
                >
                  <i class="iconfont icon-plus text-gray-400 text-xl"></i>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 提示词输入 -->
        <div class="bg-white rounded-xl p-4 shadow-sm">
          <label class="block text-gray-700 font-medium mb-3">提示词</label>
          <textarea
            v-model="currentPrompt"
            placeholder="描述你想要的视频效果"
            class="w-full px-4 py-3 border border-gray-200 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent resize-none"
            rows="4"
            maxlength="2000"
          />
          <div class="text-right mt-2">
            <span class="text-sm text-gray-500">{{ currentPrompt.length }}/2000</span>
          </div>
        </div>

        <!-- 视频比例 -->
        <CustomSelect
          v-model="imageToVideoParams.aspect_ratio"
          :options="videoAspectRatioOptions.map((opt) => ({ label: opt.label, value: opt.value }))"
          label="视频比例"
          title="选择比例"
        />
      </div>

      <!-- 生成按钮 -->
      <div class="sticky bottom-4 bg-white rounded-xl p-4 shadow-lg">
        <button
          @click="submitTask"
          :disabled="submitting"
          class="w-full py-4 bg-gradient-to-r from-blue-500 to-purple-600 text-white font-semibold rounded-xl disabled:from-gray-400 disabled:to-gray-400 disabled:cursor-not-allowed hover:from-blue-600 hover:to-purple-700 transition-all duration-200 flex items-center justify-center space-x-2"
        >
          <i v-if="submitting" class="iconfont icon-loading animate-spin"></i>
          <span>{{ submitting ? '创作中...' : `立即生成 (${currentPowerCost}算力)` }}</span>
        </button>
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
                  v-if="item.img_url"
                  :src="item.img_url"
                  fit="cover"
                  class="w-full h-full"
                  :preview-disabled="true"
                >
                  <template #error>
                    <div class="w-full h-full flex items-center justify-center bg-gray-100">
                      <i class="iconfont icon-image text-gray-400 text-xl"></i>
                    </div>
                  </template>
                </el-image>
                <el-image
                  v-else-if="item.video_url"
                  :src="item.video_url"
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
                <div v-else class="w-full h-full flex items-center justify-center bg-gray-100">
                  <i
                    :class="
                      item.type.includes('video') ? 'iconfont icon-video' : 'iconfont icon-image'
                    "
                    class="text-gray-400 text-xl"
                  ></i>
                </div>
                <!-- 播放/查看按钮 -->
                <button
                  v-if="item.status === 'completed'"
                  @click="playMedia(item)"
                  class="absolute inset-0 flex items-center justify-center bg-black bg-opacity-50 opacity-0 hover:opacity-100 transition-opacity"
                >
                  <i
                    :class="
                      item.type.includes('video') ? 'iconfont icon-play' : 'iconfont icon-eye'
                    "
                    class="text-white text-xl"
                  ></i>
                </button>
                <!-- 进度动画 -->
                <div
                  v-if="item.status === 'in_queue' || item.status === 'generating'"
                  class="absolute inset-0 flex items-center justify-center bg-blue-500 bg-opacity-20"
                >
                  <i class="iconfont icon-loading animate-spin text-blue-500 text-xl"></i>
                </div>
                <!-- 失败状态 -->
                <div
                  v-if="item.status === 'failed'"
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
                    {{ getFunctionName(item.type) }}
                  </h3>
                  <p class="text-gray-500 text-sm mt-1 line-clamp-2">
                    {{ item.prompt }}
                  </p>
                </div>
                <!-- 任务状态 -->
                <div v-if="item.status !== 'completed'" class="flex items-center space-x-2 text-sm">
                  <div
                    v-if="item.status === 'failed'"
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
                  :class="[
                    'px-2 py-1 text-xs rounded-full',
                    getTaskType(item.type) === 'warning'
                      ? 'bg-yellow-100 text-yellow-600'
                      : 'bg-blue-100 text-blue-600',
                  ]"
                >
                  {{ getFunctionName(item.type) }}
                </span>
                <span
                  v-if="item.power"
                  class="px-2 py-1 text-xs bg-orange-100 text-orange-600 rounded-full"
                >
                  {{ item.power }}算力
                </span>
              </div>
            </div>
          </div>

          <!-- 操作按钮 -->
          <div class="flex items-center justify-between mt-4">
            <div class="flex space-x-2">
              <button
                v-if="item.status === 'completed'"
                @click="playMedia(item)"
                class="px-3 py-1.5 bg-blue-600 text-white text-sm rounded-lg hover:bg-blue-700 transition-colors flex items-center space-x-1"
              >
                <i
                  :class="item.type.includes('video') ? 'iconfont icon-play' : 'iconfont icon-eye'"
                  class="!text-xs"
                ></i>
                <span>{{ item.type.includes('video') ? '播放' : '查看' }}</span>
              </button>
              <button
                v-if="item.status === 'completed'"
                @click="downloadFile(item)"
                :disabled="item.downloading"
                class="px-3 py-1.5 bg-green-600 text-white text-sm rounded-lg hover:bg-green-700 transition-colors disabled:bg-gray-400 flex items-center space-x-1"
              >
                <i v-if="item.downloading" class="iconfont icon-loading animate-spin !text-xs"></i>
                <i v-else class="iconfont icon-download !text-xs"></i>
                <span>{{ item.downloading ? '下载中...' : '下载' }}</span>
              </button>
              <button
                v-if="item.status === 'failed'"
                @click="retryTask(item.id)"
                class="px-3 py-1.5 bg-orange-600 text-white text-sm rounded-lg hover:bg-orange-700 transition-colors flex items-center space-x-1"
              >
                <i class="iconfont icon-refresh !text-xs"></i>
                <span>重试</span>
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

    <!-- 媒体预览弹窗 -->
    <div
      v-if="showMediaDialog"
      class="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50"
      @click="showMediaDialog = false"
    >
      <div @click.stop class="bg-white rounded-2xl w-full max-w-4xl max-h-[80vh] animate-scale-up">
        <div class="flex items-center justify-between p-4 border-b">
          <h3 class="text-lg font-semibold text-gray-900">媒体预览</h3>
          <button @click="showMediaDialog = false" class="p-2 hover:bg-gray-100 rounded-full">
            <i class="iconfont icon-close text-gray-500"></i>
          </button>
        </div>
        <div class="p-6">
          <img
            v-if="currentMediaUrl && !currentMediaUrl.includes('video')"
            :src="currentMediaUrl"
            class="w-full max-h-[60vh] object-contain rounded-lg"
          />
          <video
            v-else-if="currentMediaUrl"
            :src="currentMediaUrl"
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
import CustomSelect from '@/components/ui/CustomSelect.vue'
import { checkSession } from '@/store/cache'
import { closeLoading, showLoading, showMessageError, showMessageSuccess } from '@/utils/dialog'
import { httpGet, httpPost } from '@/utils/http'
import { showConfirmDialog } from 'vant'
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

// 响应式数据
const activeCategory = ref('image_generation')
const useImageInput = ref(false)
const submitting = ref(false)
const listLoading = ref(false)
const listFinished = ref(false)
const currentList = ref([])
const showMediaDialog = ref(false)
const currentMediaUrl = ref('')

// 功能分类
const categories = ref([
  { key: 'image_generation', name: '图像生成' },
  { key: 'image_editing', name: '图像编辑' },
  { key: 'image_effects', name: '图像特效' },
  { key: 'video_generation', name: '视频生成' },
])

// 选项数据
const imageSizeOptions = [
  { label: '512x512', value: '512x512' },
  { label: '768x768', value: '768x768' },
  { label: '1024x1024', value: '1024x1024' },
  { label: '1024x1536', value: '1024x1536' },
  { label: '1536x1024', value: '1536x1024' },
]

const videoAspectRatioOptions = [
  { label: '16:9', value: '16:9' },
  { label: '9:16', value: '9:16' },
  { label: '1:1', value: '1:1' },
  { label: '4:3', value: '4:3' },
]

const imageEffectsTemplateOptions = [
  { label: '亚克力装饰', value: 'acrylic_ornaments' },
  { label: '天使小雕像', value: 'angel_figurine' },
  { label: '毛毫3D拍立得', value: 'felt_3d_polaroid' },
  { label: '水彩插图', value: 'watercolor_illustration' },
]

// 当前提示词
const currentPrompt = ref('')

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

// 选项数据
const imageSizeOptions = ['512x512', '768x768', '1024x1024', '1024x1536', '1536x1024']

const videoAspectRatioOptions = ['16:9', '9:16', '1:1', '4:3']

const imageEffectsTemplateOptions = [
  'acrylic_ornaments',
  'angel_figurine',
  'felt_3d_polaroid',
  'watercolor_illustration',
]

// 页面数据
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)
const currentPowerCost = ref(0)
const taskPulling = ref(true)
const tastPullHandler = ref(null)

// 计算属性
const activeFunction = computed(() => {
  if (activeCategory.value === 'image_generation') {
    return useImageInput.value ? 'image_to_image' : 'text_to_image'
  } else if (activeCategory.value === 'image_editing') {
    return 'image_edit' // 可以根据需要添加更多编辑功能
  } else if (activeCategory.value === 'video_generation') {
    return useImageInput.value ? 'image_to_video' : 'text_to_video'
  }
  return 'text_to_image'
})

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

// 获取分类图标
const getCategoryIcon = (category) => {
  const iconMap = {
    image_generation: 'iconfont icon-image',
    image_editing: 'iconfont icon-edit',
    image_effects: 'iconfont icon-chuangzuo',
    video_generation: 'iconfont icon-video',
  }
  return iconMap[category] || 'iconfont icon-image'
}

// 切换分类
const switchCategory = (key) => {
  activeCategory.value = key
  useImageInput.value = false
}

// 切换输入模式
const switchInputMode = () => {
  currentPrompt.value = ''
}

// 处理多图片上传
const handleMultipleImageUpload = (event) => {
  const files = Array.from(event.target.files)
  files.forEach((file) => {
    if (imageToVideoParams.value.image_urls.length < 2) {
      onImageUpload({ file, name: file.name })
    }
  })
}

// 移除图片
const removeImage = (index) => {
  imageToVideoParams.value.image_urls.splice(index, 1)
}

const onImageUpload = (file) => {
  const formData = new FormData()
  formData.append('file', file.file, file.name)
  showLoading('正在上传图片...')

  httpPost('/api/upload', formData)
    .then((res) => {
      showMessageSuccess('图片上传成功')
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

  httpPost('/api/jimeng/create', params)
    .then(() => {
      fetchData(1)
      taskPulling.value = true
      showMessageSuccess('创建任务成功')
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
  httpGet('/api/jimeng/list', { page: page.value, page_size: pageSize.value })
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
  httpPost('/api/jimeng/retry', { id })
    .then(() => {
      showMessageSuccess('重试任务成功')
      fetchData(1)
    })
    .catch((e) => {
      showMessageError('重试任务失败：' + e.message)
    })
}

const removeJob = (item) => {
  showConfirmDialog({
    title: '确认删除',
    message: '此操作将会删除任务相关文件，继续操作吗?',
    confirmButtonText: '确认删除',
    cancelButtonText: '取消',
  })
    .then(() => {
      httpGet('/api/jimeng/remove', { id: item.id })
        .then(() => {
          showMessageSuccess('任务删除成功')
          fetchData(1)
        })
        .catch((e) => {
          showMessageError('任务删除失败：' + e.message)
        })
    })
    .catch(() => {})
}

// 工具方法
const goBack = () => {
  router.back()
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
