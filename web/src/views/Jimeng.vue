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
            <ImageUpload v-model="store.imageToImageParams.image_input" />
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
            <ImageUpload v-model="store.imageEditParams.image_urls" :multiple="true" />
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
            <ImageUpload v-model="store.imageEffectsParams.image_input1" />
          </div>

          <div class="param-line pt">
            <span class="label">特效模板：</span>
          </div>
          <div class="param-line">
            <el-select v-model="store.imageEffectsParams.template_id" placeholder="选择特效模板">
              <el-option label="经典特效" value="classic" />
              <el-option label="艺术风格" value="artistic" />
              <el-option label="现代科技" value="modern" />
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
            <ImageUpload v-model="store.imageToVideoParams.image_urls" :multiple="true" />
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
          <el-button
            type="primary"
            @click="store.submitTask"
            :loading="store.submitting"
            :disabled="!store.isLogin || store.userPower < store.currentPowerCost"
            size="large"
          >
            立即生成 ({{ store.currentPowerCost }}<i class="iconfont icon-vip2"></i>)
          </el-button>
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
                      fit="cover"
                      class="preview-image"
                    />
                    <video
                      v-else-if="item.video_url"
                      :src="item.video_url"
                      class="preview-video"
                      preload="metadata"
                    />
                    <div v-else class="preview-placeholder">
                      <i
                        class="iconfont icon-video text-2xl"
                        v-if="item.type.includes('video')"
                      ></i>
                      <i class="iconfont icon-dalle text-2xl" v-else></i>
                      <span>{{ store.getTaskStatusText(item.status) }}</span>
                    </div>
                  </div>
                </div>
                <div class="task-center">
                  <div class="task-info flex justify-between">
                    <div class="flex gap-2">
                      <el-tag size="small" :type="store.getStatusType(item.status)">
                        {{ store.getTaskStatusText(item.status) }}
                      </el-tag>
                      <el-tag size="small">{{ store.getFunctionName(item.type) }}</el-tag>
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

                      <span class="ml-1">
                        <el-tooltip content="画同款" placement="top">
                          <i
                            class="iconfont icon-image-list cursor-pointer"
                            @click="store.drawSame(item)"
                          ></i>
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
                <div class="task-right">
                  <div class="task-actions">
                    <el-button
                      v-if="item.status === 'failed'"
                      type="primary"
                      size="small"
                      @click="store.retryTask(item.id)"
                    >
                      重试
                    </el-button>
                    <el-button
                      v-if="item.video_url || item.img_url"
                      type="default"
                      size="small"
                      @click="store.downloadFile(item)"
                    >
                      下载
                    </el-button>
                    <el-button
                      v-if="item.video_url"
                      type="default"
                      size="small"
                      @click="store.playVideo(item)"
                    >
                      播放
                    </el-button>
                    <el-button
                      type="danger"
                      v-if="item.status === 'failed'"
                      size="small"
                      @click="store.removeJob(item)"
                    >
                      删除
                    </el-button>
                  </div>
                </div>
              </div>
            </template>
          </Waterfall>
          <div class="flex justify-center py-10">
            <img
              :src="waterfallOptions.loadProps.loading"
              class="max-w-[50px] max-h-[50px]"
              v-if="store.loading"
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
    <el-dialog v-model="store.showDialog" title="视频预览" width="70%" center>
      <video :src="store.currentVideoUrl" controls style="width: 100%; max-height: 60vh">
        您的浏览器不支持视频播放
      </video>
    </el-dialog>
  </div>
</template>

<script setup>
import '@/assets/css/jimeng.styl'
import ImageUpload from '@/components/ImageUpload.vue'
import { imageSizeOptions, useJimengStore, videoAspectRatioOptions } from '@/store/jimeng'
import { useSharedStore } from '@/store/sharedata'
import { dateFormat } from '@/utils/libs'
import { Switch } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { onMounted, onUnmounted } from 'vue'
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

onMounted(() => {
  store.init()
})

onUnmounted(() => {
  store.cleanup()
})

function onWaterfallAfterRender() {
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
</script>

<style lang="stylus" scoped>
.task-list {
  .task-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
    gap: 20px;
    padding: 10px 0;
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
</style>
