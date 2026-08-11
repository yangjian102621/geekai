<template>
  <div class="page-video">
    <!-- 左侧参数设置面板 -->
    <div class="params-panel pt-2">
      <div class="provider-buttons">
        <div class="provider-grid">
          <button
            v-for="p in providerOrder"
            :key="p"
            class="provider-btn text-base"
            :class="{ active: store.activeProvider === p }"
            @click="store.switchProvider(p)"
            type="button"
          >
            <i class="iconfont mr-2 !text-xl" :class="getProviderIcon(p)"></i>
            {{ getProviderName(p) }}
          </button>
        </div>
      </div>

      <div class="function-params pt-3">
        <div class="mb-2" v-if="store.providerModels.length > 0">
          <label class="label text-left font-bold">模型选择</label>
        </div>
        <ParamBuilder
          v-model="store.formData"
          v-model:required-keys="store.requiredKeys"
          :items="store.providerModels"
          @price-params-change="handlePriceParamsChange"
        />

        <div
          class="power-info flex items-center justify-between mb-4 mt-3 p-3 rounded-lg bg-gradient-to-r from-blue-50 to-purple-50 border border-blue-200 shadow-sm"
        >
          <div class="flex items-center space-x-2">
            <el-icon color="#f59e42" size="20"><i class="iconfont icon-lightning"></i></el-icon>
            <span class="font-medium text-gray-700">当前可用算力：</span>
            <span class="font-bold text-lg text-yellow-500">{{ store.availablePower }}</span>
          </div>
          <el-tooltip content="算力用于生成视频，每次生成会消耗对应算力" placement="left">
            <el-icon color="#a78bfa" size="18"><InfoFilled /></el-icon>
          </el-tooltip>
        </div>

        <div class="flex justify-center" v-if="store.providerModels.length > 0">
          <button
            @click="store.createVideoTask"
            :disabled="store.submitting"
            class="w-full py-3 bg-gradient-to-r from-blue-500 to-purple-600 text-white rounded-xl disabled:from-gray-400 disabled:to-gray-400 disabled:cursor-not-allowed hover:from-blue-600 hover:to-purple-700 transition-all duration-200 flex items-center justify-center space-x-2 text-base"
            type="button"
          >
            <i v-if="store.submitting" class="iconfont icon-loading animate-spin"></i>
            <i v-else class="iconfont icon-chuangzuo"></i>
            <span>立即生成 ({{ store.currentPowerCost }}算力)</span>
          </button>
        </div>
      </div>
    </div>

    <!-- 右侧任务列表 -->
    <div
      class="main-content"
      v-loading="store.loading"
      element-loading-background="rgba(100,100,100,0.3)"
    >
      <div class="works-header">
        <h2 class="h-title text-2xl">你的作品</h2>
        <div class="filter-buttons">
          <el-button-group>
            <el-button
              :type="store.taskFilter === 'all' ? 'primary' : 'default'"
              @click="store.switchTaskFilter('all')"
              size="small"
            >
              全部
            </el-button>
            <el-button
              v-for="p in providerOrder"
              :key="`filter-${p}`"
              :type="store.taskFilter === p ? 'primary' : 'default'"
              @click="store.switchTaskFilter(p)"
              size="small"
            >
              {{ getProviderName(p) }}
            </el-button>
          </el-button-group>
        </div>
      </div>

      <div class="video-list">
        <div class="list-box" v-if="!store.noData">
          <div v-for="item in store.currentList" :key="item.id">
            <div class="item">
              <div class="left">
                <div class="container">
                  <div v-if="item.status === 'success'">
                    <video
                      class="video"
                      :src="store.replaceImg(item.video_url)"
                      preload="auto"
                      loop="loop"
                      muted="muted"
                    >
                      您的浏览器不支持视频播放
                    </video>
                    <button
                      class="play flex justify-center items-center"
                      @click="store.playVideo(item)"
                    >
                      <img src="/images/play.svg" alt="" />
                    </button>
                  </div>
                  <div
                    v-else-if="item.status === 'downloading'"
                    class="flex items-center justify-center"
                    style="height: 200px"
                  >
                    <div class="text-center">
                      <div
                        class="animate-spin rounded-full h-12 w-12 border-b-2 border-purple-600 mx-auto"
                      ></div>
                      <span class="text-sm text-purple-600 mt-2 block">视频下载中...</span>
                    </div>
                  </div>
                  <el-image
                    src="/images/failed.jpg"
                    class="border rounded-lg"
                    fit="cover"
                    v-else-if="item.status === 'failed'"
                  />
                  <div
                    v-else-if="(item.progress || 0) > 0 && (item.progress || 0) < 100"
                    class="flex h-[120px] items-center justify-center"
                  >
                    <el-progress
                      type="circle"
                      :percentage="item.progress || 0"
                      :width="80"
                      :stroke-width="6"
                      class="rounded-full bg-white/95 p-1 shadow-sm flex items-center justify-center"
                    >
                      <template #default="{ percentage }">
                        <span class="flex w-full justify-center text-base font-medium text-gray-700"
                          >{{ percentage }}%</span
                        >
                      </template>
                    </el-progress>
                  </div>
                  <div class="flex !items-end justify-center h-[120px]" v-else>
                    <Generating message="正在生成视频" />
                  </div>

                  <div
                    class="absolute top-0 right-0"
                    v-if="item.status === 'pending' && !((item.progress || 0) > 0)"
                  >
                    <!-- 非 in_progress 状态才显示 status 标签 -->
                    <el-tag type="info" class="mr-1"> 排队中 </el-tag>
                  </div>
                </div>
              </div>

              <div class="center">
                <div class="pb-2">
                  <el-tag class="mr-1">{{ item.type }}</el-tag>
                  <template v-if="item.params">
                    <el-tag class="mr-1" v-if="item.params.task_type">{{
                      item.params.task_type
                    }}</el-tag>
                    <el-tag class="mr-1" v-if="item.params.model">{{ item.params.model }}</el-tag>
                    <el-tag class="mr-1" v-if="item.params.duration"
                      >{{ item.params.duration }}秒</el-tag
                    >
                    <el-tag class="mr-1" v-if="item.params.mode">{{ item.params.mode }}</el-tag>
                    <el-tag class="mr-1" v-if="item.params.size">
                      分辨率：{{ item.params.size }}
                    </el-tag>
                  </template>
                  <el-tag class="mr-1" type="warning" v-if="item.power">
                    消耗算力：{{ item.power }}
                  </el-tag>
                </div>
                <div class="failed" v-if="item.status === 'failed'">
                  任务执行失败：{{ item.err_msg }}，任务提示词：{{ item.prompt }}
                </div>
                <div class="prompt" v-else>
                  {{ store.substr(item.prompt, 1000) }}
                </div>
              </div>

              <div class="right" v-if="item.progress === 100">
                <div class="tools">
                  <el-tooltip content="复制提示词" placement="top">
                    <button class="btn btn-icon copy-prompt" :data-clipboard-text="item.prompt">
                      <i class="iconfont icon-copy"></i>
                    </button>
                  </el-tooltip>

                  <el-tooltip content="下载视频" placement="top">
                    <button
                      class="btn btn-icon"
                      @click="store.downloadVideo(item)"
                      :disabled="item.downloading"
                    >
                      <i class="iconfont icon-download" v-if="!item.downloading"></i>
                      <el-image src="/images/loading.gif" class="downloading" fit="cover" v-else />
                    </button>
                  </el-tooltip>

                  <el-tooltip content="删除" placement="top">
                    <button class="btn btn-icon" @click="store.removeJob(item)">
                      <i class="iconfont icon-remove"></i>
                    </button>
                  </el-tooltip>
                </div>
              </div>

              <div class="right-error" v-else>
                <el-button type="danger" @click="store.removeJob(item)" circle>
                  <i class="iconfont icon-remove"></i>
                </el-button>
              </div>
            </div>
          </div>
        </div>

        <el-empty
          :image-size="100"
          :image="store.nodata"
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
            :total="store.total"
          />
        </div>
      </div>
    </div>

    <!-- 视频预览对话框 -->
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
import ParamBuilder from '@/components/ParamBuilder.vue'
import Generating from '@/components/ui/Generating.vue'
import { useVideoStore } from '@/store/video'
import { InfoFilled } from '@element-plus/icons-vue'
import { getProviderName } from '@/store/data/video_params'
import { computed, onMounted, onUnmounted, ref } from 'vue'

const store = useVideoStore()
const videoPlayerRef = ref(null)

const providerOrder = computed(() =>
  ['sora', 'veo', 'doubao', 'keling', 'minimax', 'wan'].filter((p) => store.providers.includes(p))
)

const getProviderIcon = (provider) => {
  const icons = {
    sora: 'icon-sora',
    doubao: 'icon-doubao',
    veo: 'icon-gemini',
    keling: 'icon-keling',
    minimax: 'icon-minimax',
    wan: 'icon-wan',
  }
  return icons[provider] || 'icon-video'
}

// 状态配置
const statusConfig = {
  pending: { label: '等待中', type: 'info' },
  in_progress: { label: '进行中', type: 'primary' },
  downloading: { label: '下载中', type: 'success' },
  success: { label: '成功', type: 'success' },
  failed: { label: '失败', type: 'danger' },
}

// 处理价格参数变化事件
const handlePriceParamsChange = () => {
  // 价格参数变化时，store 中的 watch 会自动触发 setCurrentPowerCost
  // setCurrentPowerCost 是异步的，会调用 API 获取最新算力值
  // 无需额外处理，watch 会自动更新 currentPowerCost
}

// 处理视频对话框关闭事件
const handleVideoDialogClose = () => {
  if (videoPlayerRef.value) {
    videoPlayerRef.value.pause()
    videoPlayerRef.value.currentTime = 0
  }
  store.showDialog = false
}

onMounted(() => {
  store.init()
})

onUnmounted(() => {
  store.cleanup()
})
</script>

<style lang="scss" scoped>
@use '../assets/css/video.scss' as *;

.provider-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 10px;
}

.provider-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 10px 12px;
  border-radius: 10px;
  border: 1px solid rgba(255, 255, 255, 0.08);
  background: rgba(0, 0, 0, 0.06);
  color: var(--text-theme-color);
  transition: all 0.15s ease;

  &.active {
    background: linear-gradient(90deg, #3b82f6, #a855f7);
    border-color: rgba(99, 102, 241, 0.6);
    color: #fff;
    box-shadow: 0 2px 8px rgba(59, 130, 246, 0.35);
  }

  &:hover:not(.active) {
    background: rgba(0, 0, 0, 0.1);
  }
}
</style>
