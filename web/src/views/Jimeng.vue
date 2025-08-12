<template>
  <div class="page-jimeng">
    <!-- 左侧参数设置面板 -->
    <div class="params-panel">
      <!-- 功能分类按钮组 -->
      <div class="category-buttons">
        <div class="category-grid">
          <div
            v-for="category in store.categories"
            :key="category.key"
            :class="['category-btn', { active: store.activeCategory === category.key }]"
            @click="store.switchCategory(category.key)"
          >
            <div class="category-icon">
              <i :class="getCategoryIcon(category.key)"></i>
            </div>
            <div class="category-name">{{ category.name }}</div>
          </div>
        </div>
      </div>

      <!-- 功能开关 -->
      <div
        class="function-switch"
        v-if="
          store.activeCategory === 'image_generation' || store.activeCategory === 'video_generation'
        "
      >
        <div class="switch-label">
          <el-icon><Switch /></el-icon>
          生成模式
        </div>
        <div class="switch-container">
          <div class="switch-info">
            <div class="switch-title">
              {{ store.activeCategory === 'image_generation' ? '图生图人像写真' : '图生视频' }}
            </div>
          </div>
          <el-switch v-model="store.useImageInput" @change="store.switchInputMode" />
        </div>
      </div>

      <!-- 参数容器 -->
      <div class="params-container">
        <!-- 文生图 -->
        <div v-if="store.activeFunction === 'text_to_image'" class="function-panel">
          <div class="param-line pt">
            <span class="label">提示词：</span>
          </div>
          <div class="param-line">
            <el-input
              v-model="store.currentPrompt"
              type="textarea"
              :autosize="{ minRows: 3, maxRows: 5 }"
              placeholder="请输入图片描述，越详细越好"
              maxlength="2000"
              show-word-limit
            />
          </div>

          <div class="param-line pt">
            <span class="label">图片尺寸：</span>
          </div>
          <div class="param-line">
            <el-select v-model="store.textToImageParams.size" placeholder="选择尺寸">
              <el-option
                v-for="opt in imageSizeOptions"
                :key="opt.value"
                :label="opt.label"
                :value="opt.value"
              />
            </el-select>
          </div>

          <div class="param-line">
            <span class="label"
              >创意度
              <el-tooltip content="创意度越高，影响文本描述的程度越高" placement="top">
                <i class="iconfont icon-info cursor-pointer ml-1"></i> </el-tooltip
            ></span>
          </div>
          <div class="item-group">
            <el-slider v-model="store.textToImageParams.scale" :min="1" :max="10" :step="0.5" />
          </div>

          <div class="item-group flex justify-between">
            <span class="label">智能优化提示词</span>
            <el-switch v-model="store.textToImageParams.use_pre_llm" />
          </div>
        </div>

        <!-- 图生图 -->
        <div v-if="store.activeFunction === 'image_to_image'" class="function-panel">
          <div class="param-line pt">
            <span class="label">上传图片：</span>
          </div>
          <div class="param-line">
            <ImageUpload
              v-model="store.imageToImageParams.image_input"
              :max-count="1"
              :multiple="false"
            />
          </div>

          <div class="param-line pt">
            <span class="label">提示词：</span>
          </div>
          <div class="param-line">
            <el-input
              v-model="store.currentPrompt"
              type="textarea"
              :autosize="{ minRows: 3, maxRows: 5 }"
              placeholder="描述你想要的图片效果"
              maxlength="2000"
              show-word-limit
            />
          </div>

          <div class="param-line pt">
            <span class="label">图片尺寸：</span>
          </div>
          <div class="param-line">
            <el-select v-model="store.imageToImageParams.size" placeholder="选择尺寸">
              <el-option
                v-for="opt in imageSizeOptions"
                :key="opt.value"
                :label="opt.label"
                :value="opt.value"
              />
            </el-select>
          </div>
        </div>

        <!-- 图像编辑 -->
        <div v-if="store.activeFunction === 'image_edit'" class="function-panel">
          <div class="param-line pt">
            <span class="label">上传图片：</span>
          </div>
          <div class="param-line">
            <ImageUpload
              v-model="store.imageEditParams.image_input"
              :max-count="1"
              :multiple="false"
            />
          </div>

          <div class="param-line pt">
            <span class="label">编辑提示词：</span>
          </div>
          <div class="param-line">
            <el-input
              v-model="store.currentPrompt"
              type="textarea"
              :autosize="{ minRows: 3, maxRows: 5 }"
              placeholder="描述你想要的编辑效果"
              maxlength="2000"
              show-word-limit
            />
          </div>

          <div class="item-group">
            <span class="label">编辑强度：</span>
            <el-slider v-model="store.imageEditParams.scale" :min="0" :max="1" :step="0.1" />
          </div>
        </div>

        <!-- 图像特效 -->
        <div v-if="store.activeFunction === 'image_effects'" class="function-panel">
          <div class="param-line pt">
            <span class="label">上传图片：</span>
          </div>
          <div class="param-line">
            <ImageUpload
              v-model="store.imageEffectsParams.image_input"
              :max-count="1"
              :multiple="false"
            />
          </div>

          <div class="param-line pt">
            <span class="label">特效模板：</span>
          </div>
          <div class="param-line">
            <el-select
              v-model="store.imageEffectsParams.template_id"
              placeholder="选择特效模板"
              popper-class="jimeng-template-select"
              @change="handleTemplateChange($event)"
            >
              <template #prefix>
                <div class="flex items-center py-1">
                  <el-image
                    v-if="templatePreview"
                    :src="templatePreview"
                    class="w-[50px] h-[50px] object-cover rounded-md"
                    :preview-src-list="[templatePreview]"
                    :preview-teleported="true"
                    @click.stop
                  />
                </div>
              </template>
              <el-option
                v-for="opt in imageEffectsTemplateOptions"
                :key="opt.value"
                :label="opt.label"
                :value="opt.value"
              >
                <div class="flex flex-row justify-between">
                  <span class="template-label">{{ opt.label }}</span>
                  <img
                    v-if="opt.preview"
                    :src="opt.preview"
                    :alt="opt.label"
                    class="w-[50px] h-[50px] object-cover rounded-md"
                  />
                </div>
              </el-option>
            </el-select>
          </div>

          <div class="param-line pt">
            <span class="label">输出尺寸：</span>
          </div>
          <div class="param-line">
            <el-select v-model="store.imageEffectsParams.size" placeholder="选择尺寸">
              <el-option
                v-for="opt in imageSizeOptions"
                :key="opt.value"
                :label="opt.label"
                :value="opt.value"
              />
            </el-select>
          </div>
        </div>

        <!-- 文生视频 -->
        <div v-if="store.activeFunction === 'text_to_video'" class="function-panel">
          <div class="param-line pt">
            <span class="label">提示词：</span>
          </div>
          <div class="param-line">
            <el-input
              v-model="store.currentPrompt"
              type="textarea"
              :autosize="{ minRows: 3, maxRows: 5 }"
              placeholder="描述你想要的视频内容"
              maxlength="2000"
              show-word-limit
            />
          </div>

          <div class="param-line pt">
            <span class="label">视频比例：</span>
          </div>
          <div class="param-line">
            <el-select v-model="store.textToVideoParams.aspect_ratio" placeholder="选择比例">
              <el-option
                v-for="opt in videoAspectRatioOptions"
                :key="opt.value"
                :label="opt.label"
                :value="opt.value"
              />
            </el-select>
          </div>
        </div>

        <!-- 图生视频 -->
        <div v-if="store.activeFunction === 'image_to_video'" class="function-panel">
          <div class="param-line pt">
            <span class="label">上传图片：</span>
          </div>
          <div class="param-line">
            <ImageUpload
              v-model="store.imageToVideoParams.image_input"
              :max-count="2"
              :multiple="true"
            />
          </div>

          <div class="param-line pt">
            <span class="label">提示词：</span>
          </div>
          <div class="param-line">
            <el-input
              v-model="store.currentPrompt"
              type="textarea"
              :autosize="{ minRows: 3, maxRows: 5 }"
              placeholder="描述你想要的视频效果"
              maxlength="2000"
              show-word-limit
            />
          </div>

          <div class="param-line pt">
            <span class="label">视频比例：</span>
          </div>
          <div class="param-line">
            <el-select v-model="store.imageToVideoParams.aspect_ratio" placeholder="选择比例">
              <el-option
                v-for="opt in videoAspectRatioOptions"
                :key="opt.value"
                :label="opt.label"
                :value="opt.value"
              />
            </el-select>
          </div>
        </div>

        <!-- 提交按钮 -->
        <div class="submit-btn flex justify-center pt-4">
          <button
            @click="store.submitTask"
            :disabled="store.submitting"
            class="w-full py-3 bg-gradient-to-r from-blue-500 to-purple-600 text-white rounded-xl disabled:from-gray-400 disabled:to-gray-400 disabled:cursor-not-allowed hover:from-blue-600 hover:to-purple-700 transition-all duration-200 flex items-center justify-center space-x-2 text-base"
          >
            <i v-if="store.submitting" class="iconfont icon-loading animate-spin"></i>
            <i v-else class="iconfont icon-chuangzuo"></i>
            <span>立即生成 ({{ store.currentPowerCost }}算力)</span>
          </button>
        </div>
      </div>
    </div>

    <!-- 右侧任务列表 -->
    <div class="main-content">
      <div class="works-header">
        <h2 class="h-title">你的作品</h2>
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
              :type="store.taskFilter === 'image' ? 'primary' : 'default'"
              @click="store.switchTaskFilter('image')"
              size="small"
            >
              图片
            </el-button>
            <el-button
              :type="store.taskFilter === 'video' ? 'primary' : 'default'"
              @click="store.switchTaskFilter('video')"
              size="small"
            >
              视频
            </el-button>
          </el-button-group>
        </div>
      </div>

      <div class="task-list" v-loading="store.loading">
        <div v-if="store.currentList.length > 0">
          <Waterfall
            :list="store.currentList"
            v-bind="waterfallOptions"
            :is-loading="store.loading"
            :is-over="store.isOver"
            :lazyload="true"
            @afterRender="onWaterfallAfterRender"
          >
            <template #default="{ item }">
              <div class="task-item">
                <!-- 保持原有内容 -->
                <div class="task-left">
                  <div class="task-preview">
                    <el-image
                      v-if="item.img_url"
                      :src="item.img_url"
                      :preview-src-list="[item.img_url]"
                      :preview-teleported="true"
                      fit="cover"
                      class="preview-image"
                    >
                      <template #placeholder>
                        <div class="w-full h-full flex justify-center items-center">
                          <img :src="loadingIcon" class="max-w-[50px] max-h-[50px]" />
                        </div>
                      </template>
                    </el-image>
                    <div v-else-if="item.video_url" class="w-full h-full preview-video-wrapper">
                      <video
                        :src="item.video_url"
                        preload="auto"
                        loop="loop"
                        muted="muted"
                        class="w-full h-full object-cover"
                      >
                        您的浏览器不支持视频播放
                      </video>
                      <div class="video-mask" @click="store.playVideo(item)">
                        <div class="play-btn">
                          <img src="/images/play.svg" alt="播放" />
                        </div>
                      </div>
                    </div>

                    <div v-else class="preview-placeholder">
                      <div
                        v-if="item.status === 'in_queue'"
                        class="flex flex-col items-center gap-1"
                      >
                        <i class="iconfont icon-video" v-if="item.type.includes('video')"></i>
                        <i class="iconfont icon-dalle" v-else></i>
                        <span>
                          {{ store.getTaskStatusText(item.status) }}
                        </span>
                      </div>
                      <div
                        v-else-if="item.status === 'generating'"
                        class="flex flex-col items-center gap-1"
                      >
                        <span>
                          <Generating>
                            <div class="text-gray-400 text-base pt-3">
                              {{ store.getTaskStatusText(item.status) }}
                            </div></Generating
                          >
                        </span>
                      </div>
                      <div
                        v-else-if="item.status === 'failed'"
                        class="flex flex-col items-center gap-1"
                      >
                        <i class="iconfont icon-error text-red-500"></i>
                        <span class="text text-red-500">
                          {{ store.getTaskStatusText(item.status) }}
                        </span>

                        <span
                          class="text-sm text-red-400 err-msg-clip cursor-pointer mx-5"
                          @click="copyErrorMsg(item.err_msg)"
                        >
                          {{ item.err_msg }}
                        </span>
                      </div>
                    </div>
                  </div>
                </div>
                <div class="task-center">
                  <div class="task-info flex justify-between">
                    <div class="flex gap-2">
                      <el-tag size="small" :type="store.getTaskType(item.type)">
                        {{ store.getFunctionName(item.type) }}
                      </el-tag>
                    </div>
                    <div class="flex gap-2">
                      <span>
                        <el-tooltip content="复制提示词" placement="top">
                          <i
                            class="iconfont icon-copy cursor-pointer"
                            @click="copyPrompt(item.prompt)"
                          ></i>
                        </el-tooltip>
                      </span>

                      <template v-if="item.status === 'failed'">
                        <span class="ml-1" v-if="item.status === 'failed'">
                          <el-tooltip content="重试" placement="top">
                            <i
                              class="iconfont icon-refresh cursor-pointer"
                              @click="store.retryTask(item.id)"
                            ></i>
                          </el-tooltip>
                        </span>
                      </template>

                      <span class="ml-1">
                        <el-tooltip content="删除" placement="top">
                          <i
                            class="iconfont icon-remove cursor-pointer text-red-500"
                            @click="store.removeJob(item)"
                          ></i>
                        </el-tooltip>
                      </span>

                      <span class="ml-1" v-if="item.video_url || item.img_url">
                        <el-tooltip content="下载" placement="top">
                          <i
                            v-if="!item.downloading"
                            class="iconfont icon-download text-sm cursor-pointer"
                            @click="store.downloadFile(item)"
                          ></i>
                          <el-image src="/images/loading.gif" class="w-4 h-4" fit="cover" v-else />
                        </el-tooltip>
                      </span>
                    </div>
                  </div>
                  <div
                    class="task-prompt line-clamp-2 min-h-[40px] text-[14px] text-theme mb-2 leading-snug break-all"
                  >
                    {{ store.substr(item.prompt, 200) }}
                  </div>
                  <div class="task-meta">
                    <span>{{ dateFormat(item.created_at) }}</span>
                    <span v-if="item.power">{{ item.power }}算力</span>
                  </div>
                </div>
              </div>
            </template>
          </Waterfall>
          <div class="flex justify-center py-10">
            <img
              :src="waterfallOptions.loadProps.loading"
              class="max-w-[50px] max-h-[50px]"
              v-if="!waterfallRendered"
            />
            <div v-else>
              <div class="no-more-data" v-if="store.isOver">
                <span class="text-gray-500 mr-2">没有更多数据了</span>
                <i class="iconfont icon-face"></i>
              </div>
            </div>
          </div>
        </div>
        <el-empty v-else :image-size="100" description="暂无记录" />
      </div>
    </div>

    <!-- 视频预览对话框 -->
    <el-dialog v-model="store.showDialog" title="视频预览" center>
      <div class="flex justify-center items-center">
        <video
          :src="store.currentVideoUrl"
          autoplay
          controls
          preload="auto"
          loop
          muted
          style="max-height: calc(100vh - 100px); max-width: 100vw; object-fit: cover"
        >
          您的浏览器不支持视频播放
        </video>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import '@/assets/css/jimeng.scss'
import loadingIcon from '@/assets/img/loading.gif'
import ImageUpload from '@/components/ImageUpload.vue'

import Generating from '@/components/ui/Generating.vue'
import {
  imageEffectsTemplateOptions,
  imageSizeOptions,
  useJimengStore,
  videoAspectRatioOptions,
} from '@/store/jimeng'
import { useSharedStore } from '@/store/sharedata'
import { dateFormat } from '@/utils/libs'
import { Switch } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { onMounted, onUnmounted, ref, watch } from 'vue'
import { Waterfall } from 'vue-waterfall-plugin-next'
import 'vue-waterfall-plugin-next/dist/style.css'

const sharedStore = useSharedStore()
const waterfallOptions = sharedStore.waterfallOptions

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

const store = useJimengStore()

// 新增：瀑布流渲染完成状态
const waterfallRendered = ref(false)
// 新增：模板预览图
const templatePreview = ref('')

onMounted(() => {
  store.init()
})

onUnmounted(() => {
  store.cleanup()
})

// 监听 loading，每次 loading 变为 true 时重置渲染状态
watch(
  () => store.loading,
  (val) => {
    if (val) {
      waterfallRendered.value = false
    }
  }
)

watch(
  () => store.isOver,
  (val) => {
    if (val) {
      waterfallRendered.value = true
    }
  }
)

function handleTemplateChange(value) {
  templatePreview.value = imageEffectsTemplateOptions.find((opt) => opt.value === value)?.preview
  store.imageEffectsParams.prompt = imageEffectsTemplateOptions.find(
    (opt) => opt.value === value
  )?.label
}

function onWaterfallAfterRender() {
  waterfallRendered.value = true
  if (!store.loading && !store.isOver) {
    store.fetchData(store.page + 1)
  }
}

function copyPrompt(prompt) {
  navigator.clipboard
    .writeText(prompt)
    .then(() => {
      ElMessage.success('提示词已复制')
    })
    .catch(() => {
      ElMessage.error('复制失败')
    })
}

function copyErrorMsg(msg) {
  navigator.clipboard
    .writeText(msg)
    .then(() => {
      ElMessage.success('错误信息已复制')
    })
    .catch(() => {
      ElMessage.error('复制失败')
    })
}
</script>

<style lang="scss" scoped>
.task-list {
  .task-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
    gap: 20px;
    padding: 10px 0;
  }
  // 新增：增强任务项悬停动画
  .task-item {
    transition: box-shadow 3s cubic-bezier(0.4, 0, 0.2, 1),
      transform 0.5s cubic-bezier(0.4, 0, 0.2, 1), border-color 0.5s;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
    border: 1.5px solid transparent;
    border-radius: 12px;
    background: #fff;
    position: relative;
    z-index: 1;
  }
  .task-item:hover {
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.18), 0 1.5px 8px rgba(0, 0, 0, 0.1);
    border-color: #a259ff;
    transform: scale(1.025) translateY(-2px);
    z-index: 10;
  }
}
@media (max-width: 1200px) {
  .task-list .task-grid {
    grid-template-columns: repeat(auto-fill, minmax(260px, 1fr));
  }
}
@media (max-width: 768px) {
  .task-list .task-grid {
    grid-template-columns: 1fr;
  }
}
.preview-video-wrapper {
  position: relative;
  width: 100%;
  height: 100%;

  .video-mask {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background: rgba(0, 0, 0, 0.25);
    display: flex;
    justify-content: center;
    align-items: center;
    opacity: 0;
    transition: opacity 0.2s;
    z-index: 2;
  }

  &:hover .video-mask {
    opacity: 1;
  }

  .play-btn {
    width: 64px;
    height: 64px;
    background: rgba(255, 255, 255, 0.3);
    border-radius: 50%;
    display: flex;
    justify-content: center;
    align-items: center;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
    cursor: pointer;
    z-index: 3;
    transition: background 0.2s;

    &:hover {
      background: rgba(255, 255, 255, 0.4);
    }

    img {
      width: 36px;
      height: 36px;
    }
  }
}

.err-msg-clip {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  text-overflow: ellipsis;
  word-break: break-all;
  white-space: normal;
}

.jimeng-template-select {
  .el-select-dropdown__item {
    height: 60px;
    line-height: 60px;
  }
}
</style>
