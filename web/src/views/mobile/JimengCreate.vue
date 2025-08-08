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
      <div class="jimeng-create__category-section">
        <CustomTabs
          v-model="jimengStore.activeCategory"
          @update:modelValue="jimengStore.switchCategory"
        >
          <CustomTabPane
            v-for="category in jimengStore.categories"
            :key="category.key"
            :label="category.name"
            :name="category.key"
          >
            <template #label>
              <span>{{ category.name }}</span>
            </template>

            <!-- 功能开关 -->
            <div
              class="jimeng-create__mode-section"
              v-if="category.key === 'image_generation' || category.key === 'video_generation'"
            >
              <div class="jimeng-create__mode-section-content">
                <div>
                  <span class="jimeng-create__mode-section-title">生成模式</span>
                  <p class="jimeng-create__mode-section-description">
                    {{ category.key === 'image_generation' ? '图生图人像写真' : '图生视频' }}
                  </p>
                </div>
                <el-switch
                  v-model="jimengStore.useImageInput"
                  @change="jimengStore.switchInputMode"
                  size="default"
                />
              </div>
            </div>

            <!-- 参数容器 -->
            <div class="jimeng-create__params-container">
              <!-- 文生图 -->
              <div
                v-if="jimengStore.activeFunction === 'text_to_image'"
                class="jimeng-create__function-panel"
              >
                <div class="jimeng-create__param-line">
                  <span class="jimeng-create__param-label">提示词：</span>
                </div>
                <div class="jimeng-create__param-line">
                  <textarea
                    v-model="jimengStore.currentPrompt"
                    placeholder="请输入图片描述，越详细越好"
                    class="jimeng-create__form-section-textarea"
                    rows="4"
                    maxlength="2000"
                  />
                  <div class="jimeng-create__form-section-counter">
                    <span>{{ jimengStore.currentPrompt.length }}/2000</span>
                  </div>
                </div>

                <div class="jimeng-create__param-line">
                  <span class="jimeng-create__param-label">图片尺寸：</span>
                </div>
                <div class="jimeng-create__param-line">
                  <CustomSelect
                    v-model="jimengStore.textToImageParams.size"
                    :options="
                      jimengStore.imageSizeOptions.map((opt) => ({
                        label: opt.label,
                        value: opt.value,
                      }))
                    "
                    label="图片尺寸"
                    title="选择尺寸"
                  />
                </div>

                <div class="jimeng-create__param-line">
                  <span class="jimeng-create__param-label">
                    创意度
                    <el-tooltip content="创意度越高，影响文本描述的程度越高" placement="top">
                      <i class="iconfont icon-info cursor-pointer ml-1"></i>
                    </el-tooltip>
                  </span>
                </div>
                <div class="jimeng-create__param-line">
                  <el-slider
                    v-model="jimengStore.textToImageParams.scale"
                    :min="1"
                    :max="10"
                    :step="0.5"
                  />
                </div>

                <div class="jimeng-create__param-line">
                  <div class="jimeng-create__switch-section">
                    <span class="jimeng-create__param-label">智能优化提示词</span>
                    <el-switch v-model="jimengStore.textToImageParams.use_pre_llm" size="default" />
                  </div>
                </div>
              </div>

              <!-- 图生图 -->
              <div
                v-if="jimengStore.activeFunction === 'image_to_image'"
                class="jimeng-create__function-panel"
              >
                <div class="jimeng-create__param-line">
                  <span class="jimeng-create__param-label">上传图片：</span>
                </div>
                <div class="jimeng-create__param-line">
                  <div class="jimeng-create__upload">
                    <input
                      ref="imageToImageInput"
                      type="file"
                      accept=".jpg,.png,.jpeg"
                      @change="
                        (e) =>
                          jimengStore.onImageUpload({
                            file: e.target.files[0],
                            name: e.target.files[0]?.name,
                          })
                      "
                      class="hidden"
                    />
                    <div
                      @click="$refs.imageToImageInput?.click()"
                      class="jimeng-create__upload-content"
                    >
                      <i
                        v-if="!jimengStore.imageToImageParams.image_input.length"
                        class="jimeng-create__upload-icon iconfont icon-upload"
                      ></i>
                      <span
                        v-if="!jimengStore.imageToImageParams.image_input.length"
                        class="jimeng-create__upload-text"
                        >上传图片</span
                      >
                      <div v-else class="jimeng-create__upload-preview">
                        <el-image
                          :src="
                            jimengStore.imageToImageParams.image_input[0]?.url ||
                            jimengStore.imageToImageParams.image_input[0]?.content
                          "
                          fit="cover"
                          class="w-32 h-32 rounded"
                        />
                        <button
                          @click.stop="jimengStore.imageToImageParams.image_input = []"
                          class="jimeng-create__upload-remove-btn"
                        >
                          <i class="iconfont icon-close"></i>
                        </button>
                      </div>
                    </div>
                  </div>
                </div>

                <div class="jimeng-create__param-line">
                  <span class="jimeng-create__param-label">提示词：</span>
                </div>
                <div class="jimeng-create__param-line">
                  <textarea
                    v-model="jimengStore.currentPrompt"
                    placeholder="描述你想要的图片效果"
                    class="jimeng-create__form-section-textarea"
                    rows="4"
                    maxlength="2000"
                  />
                  <div class="jimeng-create__form-section-counter">
                    <span>{{ jimengStore.currentPrompt.length }}/2000</span>
                  </div>
                </div>

                <div class="jimeng-create__param-line">
                  <span class="jimeng-create__param-label">图片尺寸：</span>
                </div>
                <div class="jimeng-create__param-line">
                  <CustomSelect
                    v-model="jimengStore.imageToImageParams.size"
                    :options="
                      jimengStore.imageSizeOptions.map((opt) => ({
                        label: opt.label,
                        value: opt.value,
                      }))
                    "
                    label="图片尺寸"
                    title="选择尺寸"
                  />
                </div>
              </div>

              <!-- 图像编辑 -->
              <div
                v-if="jimengStore.activeFunction === 'image_edit'"
                class="jimeng-create__function-panel"
              >
                <div class="jimeng-create__param-line">
                  <span class="jimeng-create__param-label">上传图片：</span>
                </div>
                <div class="jimeng-create__param-line">
                  <div class="jimeng-create__upload">
                    <input
                      ref="imageEditInput"
                      type="file"
                      accept=".jpg,.png,.jpeg"
                      @change="
                        (e) =>
                          jimengStore.onImageUpload({
                            file: e.target.files[0],
                            name: e.target.files[0]?.name,
                          })
                      "
                      class="hidden"
                    />
                    <div
                      @click="$refs.imageEditInput?.click()"
                      class="jimeng-create__upload-content"
                    >
                      <i
                        v-if="!jimengStore.imageEditParams.image_urls.length"
                        class="jimeng-create__upload-icon iconfont icon-upload"
                      ></i>
                      <span
                        v-if="!jimengStore.imageEditParams.image_urls.length"
                        class="jimeng-create__upload-text"
                        >上传图片</span
                      >
                      <div v-else class="jimeng-create__upload-preview">
                        <el-image
                          :src="
                            jimengStore.imageEditParams.image_urls[0]?.url ||
                            jimengStore.imageEditParams.image_urls[0]?.content
                          "
                          fit="cover"
                          class="w-32 h-32 rounded"
                        />
                        <button
                          @click.stop="jimengStore.imageEditParams.image_urls = []"
                          class="jimeng-create__upload-remove-btn"
                        >
                          <i class="iconfont icon-close"></i>
                        </button>
                      </div>
                    </div>
                  </div>
                </div>

                <div class="jimeng-create__param-line">
                  <span class="jimeng-create__param-label">编辑提示词：</span>
                </div>
                <div class="jimeng-create__param-line">
                  <textarea
                    v-model="jimengStore.currentPrompt"
                    placeholder="描述你想要的编辑效果"
                    class="jimeng-create__form-section-textarea"
                    rows="4"
                    maxlength="2000"
                  />
                  <div class="jimeng-create__form-section-counter">
                    <span>{{ jimengStore.currentPrompt.length }}/2000</span>
                  </div>
                </div>

                <div class="jimeng-create__param-line">
                  <span class="jimeng-create__param-label">编辑强度：</span>
                </div>
                <div class="jimeng-create__param-line">
                  <el-slider
                    v-model="jimengStore.imageEditParams.scale"
                    :min="0"
                    :max="1"
                    :step="0.1"
                  />
                </div>
              </div>

              <!-- 图像特效 -->
              <div
                v-if="jimengStore.activeFunction === 'image_effects'"
                class="jimeng-create__function-panel"
              >
                <div class="jimeng-create__param-line">
                  <span class="jimeng-create__param-label">上传图片：</span>
                </div>
                <div class="jimeng-create__param-line">
                  <div class="jimeng-create__upload">
                    <input
                      ref="imageEffectsInput"
                      type="file"
                      accept=".jpg,.png,.jpeg"
                      @change="
                        (e) =>
                          jimengStore.onImageUpload({
                            file: e.target.files[0],
                            name: e.target.files[0]?.name,
                          })
                      "
                      class="hidden"
                    />
                    <div
                      @click="$refs.imageEffectsInput?.click()"
                      class="jimeng-create__upload-content"
                    >
                      <i
                        v-if="!jimengStore.imageEffectsParams.image_input1.length"
                        class="jimeng-create__upload-icon iconfont icon-upload"
                      ></i>
                      <span
                        v-if="!jimengStore.imageEffectsParams.image_input1.length"
                        class="jimeng-create__upload-text"
                        >上传图片</span
                      >
                      <div v-else class="jimeng-create__upload-preview">
                        <el-image
                          :src="
                            jimengStore.imageEffectsParams.image_input1[0]?.url ||
                            jimengStore.imageEffectsParams.image_input1[0]?.content
                          "
                          fit="cover"
                          class="w-32 h-32 rounded"
                        />
                        <button
                          @click.stop="jimengStore.imageEffectsParams.image_input1 = []"
                          class="jimeng-create__upload-remove-btn"
                        >
                          <i class="iconfont icon-close"></i>
                        </button>
                      </div>
                    </div>
                  </div>
                </div>

                <div class="jimeng-create__param-line">
                  <span class="jimeng-create__param-label">特效模板：</span>
                </div>
                <div class="jimeng-create__param-line">
                  <CustomSelect
                    v-model="jimengStore.imageEffectsParams.template_id"
                    :options="
                      jimengStore.imageEffectsTemplateOptions.map((opt) => ({
                        label: opt.label,
                        value: opt.value,
                      }))
                    "
                    label="特效模板"
                    title="选择特效模板"
                  />
                </div>

                <div class="jimeng-create__param-line">
                  <span class="jimeng-create__param-label">输出尺寸：</span>
                </div>
                <div class="jimeng-create__param-line">
                  <CustomSelect
                    v-model="jimengStore.imageEffectsParams.size"
                    :options="
                      jimengStore.imageSizeOptions.map((opt) => ({
                        label: opt.label,
                        value: opt.value,
                      }))
                    "
                    label="输出尺寸"
                    title="选择尺寸"
                  />
                </div>
              </div>

              <!-- 文生视频 -->
              <div
                v-if="jimengStore.activeFunction === 'text_to_video'"
                class="jimeng-create__function-panel"
              >
                <div class="jimeng-create__param-line">
                  <span class="jimeng-create__param-label">提示词：</span>
                </div>
                <div class="jimeng-create__param-line">
                  <textarea
                    v-model="jimengStore.currentPrompt"
                    placeholder="描述你想要的视频内容"
                    class="jimeng-create__form-section-textarea"
                    rows="4"
                    maxlength="2000"
                  />
                  <div class="jimeng-create__form-section-counter">
                    <span>{{ jimengStore.currentPrompt.length }}/2000</span>
                  </div>
                </div>

                <div class="jimeng-create__param-line">
                  <span class="jimeng-create__param-label">视频比例：</span>
                </div>
                <div class="jimeng-create__param-line">
                  <CustomSelect
                    v-model="jimengStore.textToVideoParams.aspect_ratio"
                    :options="
                      jimengStore.videoAspectRatioOptions.map((opt) => ({
                        label: opt.label,
                        value: opt.value,
                      }))
                    "
                    label="视频比例"
                    title="选择比例"
                  />
                </div>
              </div>

              <!-- 图生视频 -->
              <div
                v-if="jimengStore.activeFunction === 'image_to_video'"
                class="jimeng-create__function-panel"
              >
                <div class="jimeng-create__param-line">
                  <span class="jimeng-create__param-label">上传图片：</span>
                </div>
                <div class="jimeng-create__param-line">
                  <div class="jimeng-create__upload">
                    <input
                      ref="imageToVideoInput"
                      type="file"
                      accept=".jpg,.png,.jpeg"
                      multiple
                      @change="(e) => jimengStore.handleMultipleImageUpload(e)"
                      class="hidden"
                    />
                    <div
                      @click="$refs.imageToVideoInput?.click()"
                      class="jimeng-create__upload-content"
                    >
                      <i
                        v-if="!jimengStore.imageToVideoParams.image_urls.length"
                        class="jimeng-create__upload-icon iconfont icon-upload"
                      ></i>
                      <span
                        v-if="!jimengStore.imageToVideoParams.image_urls.length"
                        class="jimeng-create__upload-text"
                        >上传图片</span
                      >
                      <div v-else class="jimeng-create__upload-multiple">
                        <div
                          v-for="(image, index) in jimengStore.imageToVideoParams.image_urls"
                          :key="index"
                          class="jimeng-create__upload-multiple-item"
                        >
                          <el-image
                            :src="image?.url || image?.content"
                            fit="cover"
                            class="w-24 h-24 rounded"
                          />
                          <button
                            @click.stop="jimengStore.removeImage(index)"
                            class="jimeng-create__upload-remove-btn"
                          >
                            <i class="iconfont icon-close"></i>
                          </button>
                        </div>
                        <div
                          v-if="jimengStore.imageToVideoParams.image_urls.length < 2"
                          @click.stop="$refs.imageToVideoInput?.click()"
                          class="jimeng-create__upload-multiple-add"
                        >
                          <i class="iconfont icon-plus"></i>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>

                <div class="jimeng-create__param-line">
                  <span class="jimeng-create__param-label">提示词：</span>
                </div>
                <div class="jimeng-create__param-line">
                  <textarea
                    v-model="jimengStore.currentPrompt"
                    placeholder="描述你想要的视频效果"
                    class="jimeng-create__form-section-textarea"
                    rows="4"
                    maxlength="2000"
                  />
                  <div class="jimeng-create__form-section-counter">
                    <span>{{ jimengStore.currentPrompt.length }}/2000</span>
                  </div>
                </div>

                <div class="jimeng-create__param-line">
                  <span class="jimeng-create__param-label">视频比例：</span>
                </div>
                <div class="jimeng-create__param-line">
                  <CustomSelect
                    v-model="jimengStore.imageToVideoParams.aspect_ratio"
                    :options="
                      jimengStore.videoAspectRatioOptions.map((opt) => ({
                        label: opt.label,
                        value: opt.value,
                      }))
                    "
                    label="视频比例"
                    title="选择比例"
                  />
                </div>
              </div>

              <!-- 提交按钮 -->
              <div class="jimeng-create__submit-btn">
                <button @click="jimengStore.submitTask" :disabled="jimengStore.submitting">
                  <i v-if="jimengStore.submitting" class="iconfont icon-loading animate-spin"></i>
                  <span>{{
                    jimengStore.submitting
                      ? '创作中...'
                      : `立即生成 (${jimengStore.currentPowerCost}算力)`
                  }}</span>
                </button>
              </div>
            </div>
          </CustomTabPane>
        </CustomTabs>
      </div>
    </div>

    <!-- 作品列表 -->
    <div class="jimeng-create__works">
      <h2 class="jimeng-create__works-title">我的作品</h2>
      <div class="jimeng-create__works-list space-y-4">
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
                  fit="cover"
                  class="w-full h-full"
                  :preview-disabled="true"
                >
                  <template #error>
                    <div class="jimeng-create__works-item-thumb-placeholder">
                      <i class="iconfont icon-image"></i>
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
                    <div class="jimeng-create__works-item-thumb-placeholder">
                      <i class="iconfont icon-video"></i>
                    </div>
                  </template>
                </el-image>
                <div v-else class="jimeng-create__works-item-thumb-placeholder">
                  <i
                    :class="
                      item.type.includes('video') ? 'iconfont icon-video' : 'iconfont icon-image'
                    "
                  ></i>
                </div>
                <!-- 播放/查看按钮 -->
                <button
                  v-if="item.status === 'completed'"
                  @click="jimengStore.playMedia(item)"
                  class="jimeng-create__works-item-thumb-overlay"
                >
                  <i
                    :class="
                      item.type.includes('video') ? 'iconfont icon-play' : 'iconfont icon-eye'
                    "
                  ></i>
                </button>
                <!-- 进度动画 -->
                <div
                  v-if="item.status === 'in_queue' || item.status === 'generating'"
                  class="jimeng-create__works-item-thumb-status jimeng-create__works-item-thumb-status--loading"
                >
                  <i class="iconfont icon-loading animate-spin"></i>
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

          <!-- 操作按钮 -->
          <div class="jimeng-create__works-item-actions">
            <div class="jimeng-create__works-item-actions-left">
              <button
                v-if="item.status === 'completed'"
                @click="jimengStore.playMedia(item)"
                class="jimeng-create__works-item-actions-btn jimeng-create__works-item-actions-btn--primary"
              >
                <i
                  :class="item.type.includes('video') ? 'iconfont icon-play' : 'iconfont icon-eye'"
                ></i>
                <span>{{ item.type.includes('video') ? '播放' : '查看' }}</span>
              </button>
              <button
                v-if="item.status === 'completed'"
                @click="jimengStore.downloadFile(item)"
                :disabled="item.downloading"
                class="jimeng-create__works-item-actions-btn jimeng-create__works-item-actions-btn--success"
              >
                <i v-if="item.downloading" class="iconfont icon-loading animate-spin"></i>
                <i v-else class="iconfont icon-download"></i>
                <span>{{ item.downloading ? '下载中...' : '下载' }}</span>
              </button>
              <button
                v-if="item.status === 'failed'"
                @click="jimengStore.retryTask(item.id)"
                class="jimeng-create__works-item-actions-btn jimeng-create__works-item-actions-btn--warning"
              >
                <i class="iconfont icon-refresh"></i>
                <span>重试</span>
              </button>
            </div>
            <button
              @click="jimengStore.removeJob(item)"
              class="jimeng-create__works-item-actions-btn jimeng-create__works-item-actions-btn--danger"
            >
              <i class="iconfont icon-remove"></i>
              <span>删除</span>
            </button>
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
            <i class="iconfont icon-close"></i>
          </button>
        </div>
        <div class="jimeng-create__media-dialog-body">
          <img
            v-if="jimengStore.currentMediaUrl && !jimengStore.currentMediaUrl.includes('video')"
            :src="jimengStore.currentMediaUrl"
            class="w-full max-h-[60vh] object-contain rounded-lg"
          />
          <video
            v-else-if="jimengStore.currentMediaUrl"
            :src="jimengStore.currentMediaUrl"
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
import CustomSelect from '@/components/mobile/CustomSelect.vue'
import CustomTabPane from '@/components/ui/CustomTabPane.vue'
import CustomTabs from '@/components/ui/CustomTabs.vue'
import { checkSession } from '@/store/cache'
import { useJimengStore } from '@/store/mobile/jimeng'
import { onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const jimengStore = useJimengStore()

// 生命周期
onMounted(() => {
  checkSession()
    .then(() => {
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
@import '@/assets/css/mobile/jimeng.scss';
</style>
