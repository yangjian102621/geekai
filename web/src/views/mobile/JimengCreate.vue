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

    <!-- 功能与参数（复用 PC 端逻辑） -->
    <div class="jimeng-create__content">
      <!-- 功能分类按钮（来源于 PC 端 store.functions） -->
      <div class="bg-white rounded-xl p-4 shadow-sm mb-3">
        <CustomTabs v-model="store.activeFunction" @tab-click="store.setFunctionPowers">
          <CustomTabPane v-for="f in store.functions" :key="f.key" :name="f.key" :label="f.name">
            <template #label>
              <i class="iconfont mr-1" :class="f.icon"></i>
              {{ f.name }}
            </template>
          </CustomTabPane>
        </CustomTabs>
      </div>

      <!-- 参数构建器（移动端组件） -->
      <div class="mb-3">
        <ParamBuilderMobile
          v-model="store.formData"
          :required-keys="store.requiredKeys"
          @update:required-keys="(v) => (store.requiredKeys = v)"
          :items="store.functionParams[store.activeFunction]"
          :progress="store.progress[store.activeFunction]"
        />
      </div>

      <!-- 提交按钮 -->
      <div
        class="bg-white rounded-xl p-4 shadow-sm mb-3"
        v-if="store.functionParams[store.activeFunction].length > 0"
      >
        <button
          class="w-full py-3 bg-gradient-to-r from-blue-500 to-purple-600 text-white rounded-xl disabled:from-gray-400 disabled:to-gray-400 disabled:cursor-not-allowed hover:from-blue-600 hover:to-purple-700 transition-all duration-200 flex items-center justify-center space-x-2 text-base"
          type="button"
          @click="store.submitTask"
          :disabled="store.submitting"
        >
          <i v-if="store.submitting" class="iconfont icon-loading animate-spin"></i>
          <i v-else class="iconfont icon-chuangzuo"></i>
          <span>立即生成 ({{ store.currentPowerCost }})</span>
        </button>
      </div>
    </div>

    <!-- 作品列表 -->
    <div class="jimeng-create__works">
      <h2 class="jimeng-create__works-title">我的作品</h2>
      <van-list
        :loading="store.loading"
        @update:loading="store.loading = $event"
        :finished="store.isOver"
        finished-text="没有更多了"
        @load="onLoadMore"
      >
        <div class="flex flex-col space-y-4">
          <div
            v-for="item in store.currentList"
            :key="item.id"
            class="jimeng-create__works-item w-full"
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
                      @click="playMedia(item)"
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
                      {{ store.getFunctionName(item.type) }}
                    </h3>
                    <p class="jimeng-create__works-item-info-prompt line-clamp-2">
                      {{ item.prompt }}
                    </p>
                  </div>
                  <!-- 任务状态 -->
                  <div
                    v-if="item.status !== 'success'"
                    class="jimeng-create__works-item-info-status"
                  >
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
                      store.getTaskType(item.type) === 'warning'
                        ? 'jimeng-create__works-item-info-tags-item--warning'
                        : 'jimeng-create__works-item-info-tags-item--primary',
                    ]"
                  >
                    {{ store.getFunctionName(item.type) }}
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
                  @click="store.copyPrompt(item.prompt)"
                  class="jimeng-create__works-item-quick-action-btn"
                  title="复制提示词"
                >
                  <i class="iconfont icon-copy"></i>
                </button>

                <!-- 下载 -->
                <button
                  v-if="item.status === 'success' && (item.img_url || item.video_url)"
                  @click="store.downloadFile(item)"
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
                @click="store.retryTask(item.id)"
                class="p-2 text-green-500"
              >
                <i class="iconfont icon-refresh"></i>
                <span class="ml-1">重试</span>
              </button>

              <!-- 删除 -->
              <button @click="store.removeJob(item)" class="p-2 text-red-500">
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
                  @click="store.copyErrorMsg(item.err_msg)"
                  class="jimeng-create__works-item-error-copy-btn"
                  title="复制错误信息"
                >
                  <i class="iconfont icon-copy"></i>
                </button>
              </div>
            </div>
          </div>
        </div>
      </van-list>

      <div class="px-4" v-if="store.currentList.length === 0 && !store.loading">
        <van-empty description="暂无数据" image-size="120" />
      </div>
    </div>

    <!-- 媒体预览弹窗 -->
    <div v-if="store.showDialog" class="jimeng-create__media-dialog" @click="closeMediaDialog">
      <div @click.stop class="jimeng-create__media-dialog-content animate-scale-up">
        <div class="jimeng-create__media-dialog-header">
          <h3>媒体预览</h3>
          <button @click="closeMediaDialog">
            <i class="iconfont icon-error"></i>
          </button>
        </div>
        <div class="jimeng-create__media-dialog-body">
          <video
            :src="store.currentVideoUrl"
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
import ParamBuilderMobile from '@/components/mobile/ParamBuilderMobile.vue'
import CustomTabPane from '@/components/ui/CustomTabPane.vue'
import CustomTabs from '@/components/ui/CustomTabs.vue'
import { useJimengStore } from '@/store/jimeng'
import { onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const store = useJimengStore()

function goBack() {
  router.back()
}

function playMedia(item) {
  store.currentVideoUrl = item.video_url
  store.showDialog = true
}

function closeMediaDialog() {
  store.showDialog = false
  store.currentVideoUrl = ''
}

onMounted(() => {
  store.init()
})

onUnmounted(() => {
  if (store.cleanup) {
    store.cleanup()
  }
})

function onLoadMore() {
  store.fetchData(store.page + 1)
}
</script>

<style lang="scss" scoped>
@use '@/assets/css/mobile/jimeng.scss';
</style>
