<template>
  <div class="jimeng-create">
    <!-- 页面头部 -->
    <div class="sticky top-0 z-40 bg-white shadow-sm">
      <div class="flex items-center px-4 h-14">
        <button
          @click="goBack"
          class="flex items-center justify-center w-8 h-8 rounded-full hover:bg-gray-100 transition-colors"
        >
          <i class="iconfont icon-back text-gray-600"></i>
        </button>
        <h1 class="flex-1 text-center text-lg text-gray-900">即梦AI</h1>
        <div class="w-8"></div>
      </div>
    </div>

    <!-- 功能分类选择 -->
    <div class="jimeng-create__content">
      <CustomTabs
        v-model="jimengStore.activeCategory"
        @update:modelValue="jimengStore.switchCategory"
      >
        <CustomTabPane
          :label="jimengStore.categories[0].name"
          :name="jimengStore.categories[0].key"
        >
          <template #label>
            <span>{{ jimengStore.categories[0].name }}</span>
          </template>

          <!-- 参数容器 -->
          <div class="py-3">
            <!-- 文生图 -->
            <div class="bg-white rounded-xl p-4 shadow-sm mb-3">
              <div class="mb-3">
                <label class="text-gray-700 font-semibold">提示词：</label>
              </div>
              <el-input
                v-model="jimengStore.currentPrompt"
                type="textarea"
                placeholder="请输入图片描述，越详细越好"
                :rows="4"
                maxlength="2000"
                show-word-limit
              />
            </div>

            <!-- 功能开关 -->
            <div class="bg-white rounded-xl p-4 shadow-sm mb-3">
              <div class="flex justify-between items-center w-full">
                <span class="text-gray-700 font-semibold">图生图人像写真</span>
                <el-switch v-model="jimengStore.useImageInput" size="default" />
              </div>
            </div>

            <!-- 图生图参数 -->
            <div class="bg-white rounded-xl p-4 shadow-sm mb-3" v-if="jimengStore.useImageInput">
              <ImageUpload
                v-model="jimengStore.imageToImageParams.image_input"
                :max-count="1"
                :multiple="false"
              />
            </div>

            <div class="bg-white rounded-xl p-4 shadow-sm mb-3">
              <label class="block text-gray-700 mb-3 font-semibold">图片尺寸：</label>
              <CustomSelect
                v-model="jimengStore.textToImageParams.size"
                :options="jimengStore.imageSizeOptions"
                title="选择尺寸"
              />
            </div>

            <div class="bg-white rounded-xl p-4 shadow-sm mb-3">
              <span class="flex justify-between items-center mb-3">
                <span class="text-gray-700 font-semibold">创意度：</span>
                <el-tooltip content="创意度越高，影响文本描述的程度越高" placement="top">
                  <i class="iconfont icon-info cursor-pointer ml-1"></i>
                </el-tooltip>
              </span>

              <div class="mt-3">
                <el-slider
                  v-model="jimengStore.textToImageParams.scale"
                  :min="1"
                  :max="10"
                  :step="0.5"
                />
              </div>
            </div>

            <div class="bg-white rounded-xl p-4 shadow-sm mb-3">
              <div class="flex justify-between items-center w-full">
                <label class="text-gray-700 font-semibold">智能优化提示词</label>
                <el-switch v-model="jimengStore.textToImageParams.use_pre_llm" size="default" />
              </div>
            </div>
          </div>
        </CustomTabPane>

        <CustomTabPane
          :name="jimengStore.categories[1].key"
          :label="jimengStore.categories[1].name"
        >
          <template #label>
            <span>{{ jimengStore.categories[1].name }}</span>
          </template>

          <div class="py-3">
            <div class="bg-white rounded-xl p-4 shadow-sm mb-3">
              <div class="mb-3">
                <label class="text-gray-700 font-semibold">编辑提示词：</label>
              </div>
              <el-input
                v-model="jimengStore.currentPrompt"
                type="textarea"
                placeholder="描述你想要的编辑效果"
                :rows="4"
                maxlength="2000"
                show-word-limit
              />
            </div>

            <div class="bg-white rounded-xl p-4 shadow-sm mb-3">
              <ImageUpload
                v-model="jimengStore.imageEditParams.image_input"
                :max-count="1"
                :multiple="true"
              />
            </div>

            <div class="bg-white rounded-xl p-4 shadow-sm mb-3">
              <div class="mb-3">
                <label class="text-gray-700 font-semibold">编辑强度：</label>
              </div>
              <el-slider
                v-model="jimengStore.imageEditParams.scale"
                :min="0"
                :max="1"
                :step="0.1"
              />
            </div>
          </div>
        </CustomTabPane>

        <CustomTabPane
          :name="jimengStore.categories[2].key"
          :label="jimengStore.categories[2].name"
        >
          <template #label>
            <span>{{ jimengStore.categories[2].name }}</span>
          </template>

          <div class="py-3">
            <div class="bg-white rounded-xl p-4 shadow-sm mb-3">
              <ImageUpload
                v-model="jimengStore.imageEffectsParams.image_input"
                :max-count="1"
                :multiple="true"
              />
            </div>

            <div class="bg-white rounded-xl p-4 shadow-sm mb-3">
              <div class="mb-3">
                <label class="text-gray-700 font-semibold">特效模板：</label>
              </div>
              <CustomSelect
                v-model="jimengStore.imageEffectsParams.template_id"
                :options="jimengStore.imageEffectsTemplateOptions"
                title="选择模板"
              >
                <template #option="{ option, selected }">
                  <div class="flex items-center w-full">
                    <el-image :src="option.preview" fit="cover" class="w-10 h-10 rounded-lg mr-2" />
                    <span class="font-bold text-blue-600 mr-2">{{ option.label }}</span>
                    <span v-if="selected" class="ml-auto text-green-500"
                      ><i class="iconfont icon-success"></i
                    ></span>
                  </div>
                </template>
              </CustomSelect>
            </div>

            <div class="bg-white rounded-xl p-4 shadow-sm mb-3">
              <div class="mb-3">
                <label class="text-gray-700 font-semibold">输出尺寸：</label>
              </div>
              <CustomSelect
                v-model="jimengStore.imageEffectsParams.size"
                :options="jimengStore.imageSizeOptions"
                title="选择尺寸"
              />
            </div>
          </div>
        </CustomTabPane>

        <CustomTabPane
          :name="jimengStore.categories[3].key"
          :label="jimengStore.categories[3].name"
        >
          <template #label>
            <span>{{ jimengStore.categories[3].name }}</span>
          </template>

          <div class="py-3">
            <div class="bg-white rounded-xl p-4 shadow-sm mb-3">
              <div class="mb-3">
                <label class="text-gray-700 font-semibold">提示词：</label>
              </div>
              <el-input
                v-model="jimengStore.currentPrompt"
                type="textarea"
                placeholder="请输入你想要的视频效果"
                :rows="4"
                maxlength="2000"
                show-word-limit
              />
            </div>

            <div class="bg-white rounded-xl p-4 shadow-sm mb-3">
              <div class="flex justify-between items-center w-full">
                <label class="text-gray-700 font-semibold">使用图片辅助生成：</label>
                <el-switch v-model="jimengStore.useImageInput" size="default" />
              </div>
            </div>

            <div class="bg-white rounded-xl p-4 shadow-sm mb-3" v-if="jimengStore.useImageInput">
              <ImageUpload
                v-model="jimengStore.imageToVideoParams.image_input"
                :max-count="2"
                :multiple="true"
              />
            </div>

            <div class="bg-white rounded-xl p-4 shadow-sm mb-3">
              <div class="mb-3">
                <label class="text-gray-700 font-semibold">视频比例：</label>
              </div>
              <CustomSelect
                v-model="jimengStore.textToVideoParams.aspect_ratio"
                :options="jimengStore.videoAspectRatioOptions"
                title="选择比例"
              />
            </div>
          </div>
        </CustomTabPane>
      </CustomTabs>

      <!-- 提交按钮 -->
      <div class="bg-white rounded-xl p-4 shadow-sm mb-3">
        <button
          class="w-full py-3 bg-gradient-to-r from-blue-500 to-purple-600 text-white rounded-xl disabled:from-gray-400 disabled:to-gray-400 disabled:cursor-not-allowed hover:from-blue-600 hover:to-purple-700 transition-all duration-200 flex items-center justify-center space-x-2 text-base"
          @click="jimengStore.submitTask"
          :disabled="jimengStore.submitting"
        >
          <i v-if="jimengStore.submitting" class="iconfont icon-loading animate-spin"></i>
          <i v-else class="iconfont icon-chuangzuo"></i>
          <span>{{
            jimengStore.submitting ? '创作中...' : `立即生成 (${jimengStore.currentPowerCost}算力)`
          }}</span>
        </button>
      </div>
    </div>

    <!-- 作品列表 -->
    <div class="jimeng-create__works">
      <h2 class="jimeng-create__works-title">我的作品</h2>
      <div class="jimeng-create__works-list space-y-4" v-if="jimengStore.currentList.length > 0">
        <div
          v-for="item in jimengStore.currentList"
          :key="item.id"
          class="jimeng-create__works-item"
        >
          <div class="jimeng-create__works-item-content">
            <div class="jimeng-create__works-item-thumb">
              <div class="jimeng-create__works-item-thumb-container">
                <el-image
                  v-if="item.img_url"
                  :src="item.img_url"
                  :preview-src-list="[item.img_url]"
                  fit="cover"
                  class="w-full h-full"
                >
                  <template #error>
                    <div class="jimeng-create__works-item-thumb-placeholder">
                      <i class="iconfont icon-image"></i>
                    </div>
                  </template>
                </el-image>
                <div
                  v-else-if="item.video_url"
                  class="jimeng-create__works-item-thumb-placeholder relative"
                >
                  <video
                    :src="item.video_url"
                    preload="auto"
                    loop="loop"
                    muted="muted"
                    class="w-full h-full object-cover"
                  >
                    您的浏览器不支持视频播放
                  </video>
                  <div
                    class="video-mask absolute top-0 left-0 w-full h-full flex justify-center items-center"
                    @click="jimengStore.playMedia(item)"
                  >
                    <div class="play-btn">
                      <img src="/images/play.svg" alt="播放" />
                    </div>
                  </div>
                </div>
                <div v-else class="jimeng-create__works-item-thumb-placeholder">
                  <i
                    :class="
                      item.type.includes('video') ? 'iconfont icon-video' : 'iconfont icon-image'
                    "
                  ></i>
                </div>

                <!-- 失败状态 -->
                <div
                  v-if="item.status === 'failed'"
                  class="jimeng-create__works-item-thumb-status jimeng-create__works-item-thumb-status--failed"
                >
                  <i class="iconfont icon-warning"></i>
                </div>
              </div>
            </div>
            <div class="jimeng-create__works-item-info">
              <div class="jimeng-create__works-item-info-header">
                <div class="flex-1">
                  <h3 class="jimeng-create__works-item-info-title">
                    {{ jimengStore.getFunctionName(item.type) }}
                  </h3>
                  <p class="jimeng-create__works-item-info-prompt line-clamp-2">
                    {{ item.prompt }}
                  </p>
                </div>
                <!-- 任务状态 -->
                <div v-if="item.status !== 'success'" class="jimeng-create__works-item-info-status">
                  <div
                    v-if="item.status === 'failed'"
                    class="jimeng-create__works-item-info-status--failed"
                  >
                    <i class="iconfont icon-warning"></i>
                    <el-tag type="danger">任务失败</el-tag>
                  </div>
                  <div
                    v-else
                    class="flex items-center jimeng-create__works-item-info-status--processing"
                  >
                    <div class="loading-spinner mr-1"></div>
                    <span>生成中</span>
                  </div>
                </div>
              </div>
              <!-- 标签 -->
              <div class="jimeng-create__works-item-info-tags">
                <span
                  :class="[
                    'jimeng-create__works-item-info-tags-item',
                    jimengStore.getTaskType(item.type) === 'warning'
                      ? 'jimeng-create__works-item-info-tags-item--warning'
                      : 'jimeng-create__works-item-info-tags-item--primary',
                  ]"
                >
                  {{ jimengStore.getFunctionName(item.type) }}
                </span>
                <span
                  v-if="item.power"
                  class="jimeng-create__works-item-info-tags-item jimeng-create__works-item-info-tags-item--power"
                >
                  {{ item.power }}算力
                </span>
              </div>
            </div>
          </div>

          <!-- 快捷操作按钮 -->
          <div class="jimeng-create__works-item-quick-actions">
            <span v-if="item.status === 'success'" class="flex">
              <!-- 复制提示词 -->
              <button
                v-if="item.prompt"
                @click="jimengStore.copyPrompt(item.prompt)"
                class="jimeng-create__works-item-quick-action-btn"
                title="复制提示词"
              >
                <i class="iconfont icon-copy"></i>
              </button>

              <!-- 下载 -->
              <button
                v-if="item.status === 'success' && (item.img_url || item.video_url)"
                @click="jimengStore.downloadFile(item)"
                :disabled="item.downloading"
                class="p-2 text-blue-500"
              >
                <i v-if="item.downloading" class="iconfont icon-loading animate-spin"></i>
                <i v-else class="iconfont icon-download"></i>
                <span class="ml-1">下载</span>
              </button>
            </span>

            <!-- 重试 -->
            <button
              v-if="item.status === 'failed'"
              @click="jimengStore.retryTask(item.id)"
              class="p-2 text-green-500"
            >
              <i class="iconfont icon-refresh"></i>
              <span class="ml-1">重试</span>
            </button>

            <!-- 删除 -->
            <button @click="jimengStore.removeJob(item)" class="p-2 text-red-500">
              <i class="iconfont icon-remove"></i>
              <span class="ml-1">删除</span>
            </button>
          </div>

          <!-- 错误信息复制 -->
          <div
            v-if="item.status === 'failed' && item.err_msg"
            class="jimeng-create__works-item-error"
          >
            <div class="jimeng-create__works-item-error-content">
              <span class="jimeng-create__works-item-error-text line-clamp-3">{{
                item.err_msg
              }}</span>
              <button
                @click="jimengStore.copyErrorMsg(item.err_msg)"
                class="jimeng-create__works-item-error-copy-btn"
                title="复制错误信息"
              >
                <i class="iconfont icon-copy"></i>
              </button>
            </div>
          </div>
        </div>

        <!-- 加载更多 -->
        <div v-if="jimengStore.listLoading" class="jimeng-create__works-loading">
          <i class="iconfont icon-loading animate-spin"></i>
        </div>

        <!-- 没有更多了 -->
        <div
          v-if="jimengStore.listFinished && !jimengStore.listLoading"
          class="jimeng-create__works-finished"
        >
          没有更多了
        </div>
      </div>

      <div class="px-4" v-else>
        <van-empty description="暂无数据" image-size="120" />
      </div>
    </div>

    <!-- 媒体预览弹窗 -->
    <div
      v-if="jimengStore.showMediaDialog"
      class="jimeng-create__media-dialog"
      @click="jimengStore.closeMediaDialog"
    >
      <div @click.stop class="jimeng-create__media-dialog-content animate-scale-up">
        <div class="jimeng-create__media-dialog-header">
          <h3>媒体预览</h3>
          <button @click="jimengStore.closeMediaDialog">
            <i class="iconfont icon-error"></i>
          </button>
        </div>
        <div class="jimeng-create__media-dialog-body">
          <video
            :src="jimengStore.currentMediaUrl"
            controls
            autoplay
            class="w-full max-h-[60vh] rounded-lg object-cover"
          >
            您的浏览器不支持视频播放
          </video>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import ImageUpload from '@/components/ImageUpload.vue'
import CustomSelect from '@/components/mobile/CustomSelect.vue'
import CustomTabPane from '@/components/ui/CustomTabPane.vue'
import CustomTabs from '@/components/ui/CustomTabs.vue'
import { checkSession } from '@/store/cache'
import { useJimengStore } from '@/store/mobile/jimeng'
import { onMounted, onUnmounted, ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const jimengStore = useJimengStore()

// 模板预览相关
const templatePreview = ref('')

// 处理模板变更
const handleTemplateChange = (value) => {
  const selectedTemplate = jimengStore.imageEffectsTemplateOptions.find(
    (opt) => opt.value === value
  )
  if (selectedTemplate) {
    templatePreview.value = selectedTemplate.preview || ''
    // 自动设置提示词为模板名称
    jimengStore.currentPrompt = selectedTemplate.label
  }
}

// 生命周期
onMounted(() => {
  checkSession()
    .then(() => {
      jimengStore.init() // 初始化算力配置
      jimengStore.fetchData(1)
      jimengStore.startTaskPolling()
    })
    .catch(() => {})
})

onUnmounted(() => {
  jimengStore.stopTaskPolling()
})

// 工具方法
const goBack = () => {
  router.back()
}
</script>

<style lang="scss" scoped>
@use '@/assets/css/mobile/jimeng.scss';
</style>
