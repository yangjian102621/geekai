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

    <!-- 参数设置区域 -->
    <div class="p-4 space-y-6">
      <!-- 提供商切换 -->
      <div class="bg-white rounded-xl p-3 shadow-sm">
        <div class="grid grid-cols-2 gap-2">
          <button
            v-for="p in providerOrder"
            :key="p"
            @click="store.switchProvider(p)"
            :class="[
              'flex items-center justify-center py-2.5 px-4 rounded-lg font-medium transition-colors',
              store.activeProvider === p
                ? 'bg-blue-600 text-white'
                : 'bg-gray-100 text-gray-700 hover:bg-gray-200',
            ]"
            type="button"
          >
            <i class="iconfont mr-2 !text-xl" :class="getProviderIcon(p)"></i>
            {{ getProviderName(p) }}
          </button>
        </div>
      </div>

      <!-- 参数构建器 -->
      <div v-if="store.providerModels.length > 0" class="space-y-4">
        <ParamBuilderMobile
          v-model="store.formData"
          v-model:required-keys="store.requiredKeys"
          :items="store.providerModels"
          @price-params-change="handlePriceParamsChange"
        />

        <!-- 算力信息 -->
        <div
          class="flex items-center justify-between p-3 rounded-lg bg-gradient-to-r from-blue-50 to-purple-50 border border-blue-200 shadow-sm"
        >
          <div class="flex items-center space-x-2">
            <i class="iconfont icon-lightning text-yellow-500 !text-xl"></i>
            <span class="font-medium text-gray-700">当前可用算力：</span>
            <span class="font-bold text-lg text-yellow-500">{{ store.availablePower }}</span>
          </div>
        </div>

        <!-- 生成按钮 -->
        <div class="bg-white rounded-xl p-4 shadow-sm">
          <button
            @click="store.createVideoTask"
            :disabled="store.submitting"
            type="button"
            class="w-full py-3 bg-gradient-to-r from-blue-500 to-purple-600 text-white rounded-xl disabled:from-gray-400 disabled:to-gray-400 disabled:cursor-not-allowed hover:from-blue-600 hover:to-purple-700 transition-all duration-200 flex items-center justify-center space-x-2"
          >
            <i v-if="store.submitting" class="iconfont icon-loading animate-spin"></i>
            <i v-else class="iconfont icon-chuangzuo"></i>
            <span>立即生成 ({{ store.currentPowerCost }}算力)</span>
          </button>
        </div>
      </div>
    </div>

    <!-- 作品列表 -->
    <div class="p-4">
      <div class="flex flex-col !items-start justify-between mb-4">
        <h2 class="text-lg font-semibold text-gray-900 mb-3">我的作品</h2>
        <!-- 过滤按钮 -->
        <CustomTabs
          v-model="store.taskFilter"
          @update:model-value="handleTaskFilterChange"
          class="w-full"
        >
          <CustomTabPane name="all" label="全部"></CustomTabPane>
          <CustomTabPane
            v-for="p in providerOrder"
            :key="`filter-${p}`"
            :name="p"
            :label="getProviderName(p)"
          ></CustomTabPane>
        </CustomTabs>
      </div>
      <div class="space-y-4" v-if="store.currentList.length > 0">
        <div
          v-for="item in store.currentList"
          :key="item.id"
          class="bg-white rounded-xl p-4 shadow-sm"
        >
          <div class="flex space-x-4">
            <div class="flex-shrink-0">
              <div class="relative w-16 h-16 rounded-lg overflow-hidden bg-gray-100">
                <el-image
                  v-if="item.status === 'failed'"
                  src="/images/failed.jpg"
                  fit="cover"
                  class="w-full h-full"
                  :preview-disabled="true"
                />
                <div
                  v-else-if="item.status === 'success'"
                  class="w-full h-full flex items-center justify-center bg-gray-100"
                >
                  <video
                    class="w-full h-full object-cover"
                    :src="store.replaceImg(item.video_url)"
                    preload="auto"
                    muted
                  ></video>
                </div>
                <div
                  v-else-if="item.status === 'in_progress'"
                  class="w-full h-full flex items-center justify-center bg-gray-100"
                >
                  <Generating message="正在生成视频" />
                </div>
                <div
                  v-else
                  class="w-full h-full flex items-center justify-center bg-gray-100"
                >
                  <i class="iconfont icon-video text-gray-400 text-xl"></i>
                </div>
                <!-- 视频播放按钮 -->
                <button
                  v-if="item.status === 'success'"
                  @click="store.playVideo(item)"
                  class="absolute inset-0 flex items-center justify-center bg-black bg-opacity-50 opacity-0 hover:opacity-100 transition-opacity"
                >
                  <i class="iconfont icon-play text-white text-xl"></i>
                </button>
                <!-- 下载中状态 -->
                <div
                  v-if="item.status === 'downloading'"
                  class="absolute inset-0 flex items-center justify-center bg-purple-500 bg-opacity-20"
                >
                  <div class="text-center">
                    <i class="iconfont icon-loading animate-spin text-purple-500 text-xl"></i>
                    <div class="text-xs text-purple-600 mt-1">下载中</div>
                  </div>
                </div>
                <!-- 进度动画 -->
                <div
                  v-if="item.status === 'in_progress' || item.status === 'pending'"
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
                  <p class="text-gray-500 text-sm mt-1 line-clamp-2">
                    {{ store.substr(item.prompt, 1000) }}
                  </p>
                </div>
                <!-- 任务状态 -->
                <div v-if="item.status !== 'success'" class="flex items-center space-x-2 text-sm">
                  <div
                    v-if="item.status === 'failed'"
                    class="text-red-600 flex items-center space-x-1"
                  >
                    <i class="iconfont icon-warning"></i>
                    <span>失败</span>
                  </div>
                  <div v-else-if="item.status === 'downloading'" class="text-purple-600 flex items-center space-x-1">
                    <div
                      class="w-3 h-3 border border-purple-600 border-t-transparent rounded-full animate-spin"
                    ></div>
                    <span>下载中</span>
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
              <div class="flex items-center gap-2 mt-2 flex-wrap">
                <span
                  v-if="item.type"
                  class="px-2 py-1 text-xs bg-blue-100 text-blue-600 rounded-full"
                >
                  {{ item.type }}
                </span>
                <span
                  v-if="item.params?.task_type"
                  class="px-2 py-1 text-xs bg-blue-100 text-blue-600 rounded-full"
                >
                  {{ item.params.task_type }}
                </span>
                <span
                  v-if="item.params?.model"
                  class="px-2 py-1 text-xs bg-green-100 text-green-600 rounded-full"
                >
                  {{ item.params.model }}
                </span>
                <span
                  v-if="item.params?.duration"
                  class="px-2 py-1 text-xs bg-yellow-100 text-yellow-600 rounded-full"
                >
                  {{ item.params.duration }}秒
                </span>
                <span
                  v-if="item.params?.size"
                  class="px-2 py-1 text-xs bg-purple-100 text-purple-600 rounded-full"
                >
                  分辨率：{{ item.params.size }}
                </span>
                <span
                  v-if="item.power"
                  class="px-2 py-1 text-xs bg-orange-100 text-orange-600 rounded-full"
                >
                  消耗算力：{{ item.power }}
                </span>
              </div>
              <div v-if="item.status === 'failed'" class="mt-2 text-sm text-red-600">
                任务执行失败：{{ item.err_msg }}
              </div>
            </div>
          </div>

          <!-- 操作按钮 -->
          <div class="flex items-center justify-between mt-4">
            <div class="flex space-x-2">
              <button
                v-if="item.status === 'success'"
                @click="store.playVideo(item)"
                class="px-3 py-1.5 bg-blue-600 text-white text-sm rounded-lg hover:bg-blue-700 transition-colors flex items-center space-x-1"
              >
                <i class="iconfont icon-play !text-xs"></i>
                <span>播放</span>
              </button>
              <button
                v-if="item.status === 'success'"
                @click="store.downloadVideo(item)"
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

        <!-- 加载中 -->
        <div v-if="store.loading" class="flex justify-center py-4">
          <i class="iconfont icon-loading animate-spin text-blue-500 text-xl"></i>
        </div>
      </div>

      <div class="px-4" v-if="store.noData">
        <van-empty description="没有任何作品，赶紧去创作吧！" image-size="120" />
      </div>

      <!-- 分页 -->
      <div class="flex justify-center p-4" v-if="store.total > store.pageSize">
        <el-pagination
          background
          layout="prev, pager, next"
          :hide-on-single-page="true"
          :current-page="store.page"
          :page-size="store.pageSize"
          @current-change="store.fetchData"
          :total="store.total"
        />
      </div>
    </div>

    <!-- 视频预览弹窗 -->
    <el-dialog
      v-model="store.showDialog"
      title="预览视频"
      hide-footer
      @close="handleVideoDialogClose"
      width="auto"
    >
      <video
        ref="videoPlayerRef"
        style="max-width: 90vw; max-height: 90vh"
        :src="store.currentVideoUrl"
        preload="auto"
        :autoplay="true"
        controls="controls"
      >
        您的浏览器不支持视频播放
      </video>
    </el-dialog>
  </div>
</template>

<script setup>
import '@/assets/css/mobile/video.scss'
import ParamBuilderMobile from '@/components/mobile/ParamBuilderMobile.vue'
import Generating from '@/components/ui/Generating.vue'
import CustomTabs from '@/components/ui/CustomTabs.vue'
import CustomTabPane from '@/components/ui/CustomTabPane.vue'
import { useVideoStore } from '@/store/video'
import { getProviderName } from '@/store/data/video_params'
import { showConfirmDialog } from 'vant'
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const store = useVideoStore()
const videoPlayerRef = ref(null)

// 页面专属方法
const goBack = () => {
  router.back()
}

// 提供商图标映射
const getProviderIcon = (provider) => {
  const icons = {
    veo: 'icon-gemini',
    sora: 'icon-sora',
    luma: 'icon-luma',
    keling: 'icon-keling',
    minimax: 'icon-minimax',
    wan: 'icon-wan',
  }
  return icons[provider] || 'icon-video'
}

// 提供商顺序
const providerOrder = computed(() =>
  ['sora', 'veo', 'luma', 'keling', 'minimax', 'wan'].filter((p) => store.providers.includes(p))
)

// 处理价格参数变化事件
const handlePriceParamsChange = () => {
  // 价格参数变化时，store 中的 watch 会自动触发 setCurrentPowerCost
  // setCurrentPowerCost 是异步的，会调用 API 获取最新算力值
  // 无需额外处理，watch 会自动更新 currentPowerCost
}

// 处理任务过滤变化
const handleTaskFilterChange = (filter) => {
  store.switchTaskFilter(filter)
}

// 处理视频对话框关闭事件
const handleVideoDialogClose = () => {
  if (videoPlayerRef.value) {
    videoPlayerRef.value.pause()
    videoPlayerRef.value.currentTime = 0
  }
  store.showDialog = false
}

// 删除弹窗（页面层处理）
const removeJob = (item) => {
  showConfirmDialog({
    title: '确认删除',
    message: '此操作将会删除任务相关文件，继续操作吗?',
    confirmButtonText: '确认删除',
    cancelButtonText: '取消',
  })
    .then(() => {
      store.removeJob(item)
    })
    .catch(() => {})
}

onMounted(() => {
  store.init()
})

onUnmounted(() => {
  store.cleanup()
})
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
