<template>
  <div class="page-video">
    <!-- 左侧参数设置面板 -->
    <div class="params-panel">
      <!-- 视频类型切换标签页 -->
      <el-tabs
        v-model="store.activeVideoType"
        @tab-change="store.switchVideoType"
        class="video-type-tabs"
      >
        <!-- Luma 视频参数 -->
        <el-tab-pane label="Luma视频" name="luma">
          <div class="params-container">
            <div class="param-line">
              <el-input
                v-model="store.lumaParams.prompt"
                type="textarea"
                maxlength="2000"
                :autosize="{ minRows: 4, maxRows: 6 }"
                placeholder="请在此输入视频提示词，用逗号分割，您也可以点击下面的提示词助手生成视频提示词"
              />
            </div>

            <!-- 提示词生成按钮 -->
            <div class="flex justify-end pt-1">
              <el-button @click="store.generatePrompt" type="primary" :loading="store.isGenerating">
                <span v-if="!store.isGenerating">
                  <i class="iconfont icon-chuangzuo"></i>
                  生成提示词
                </span>
                <span v-else>生成中...</span>
              </el-button>
            </div>

            <!-- 图片辅助生成开关 -->
            <div class="param-line pt">
              <div class="image-mode-toggle">
                <span class="label">使用图片辅助生成</span>
                <el-switch
                  v-model="store.lumaUseImageMode"
                  @change="store.toggleLumaImageMode"
                  size="small"
                />
              </div>
            </div>

            <!-- 图片上传区域（可折叠） -->
            <div v-if="store.lumaUseImageMode" class="img-inline">
              <div class="img-uploader video-img-box mr-2">
                <el-icon
                  v-if="store.lumaParams.image"
                  @click="store.removeLumaImage('start')"
                  class="removeimg"
                >
                  <CircleCloseFilled />
                </el-icon>
                <el-upload
                  class="uploader img-uploader"
                  :auto-upload="true"
                  :show-file-list="false"
                  :http-request="store.uploadLumaStartImage"
                  accept=".jpg,.png,.jpeg"
                >
                  <el-image
                    v-if="store.lumaParams.image"
                    :src="store.lumaParams.image"
                    fit="cover"
                  />
                  <div class="flex flex-col" v-else>
                    <el-icon class="mb-1 text-base"><Plus /></el-icon>
                    <span>起始帧</span>
                  </div>
                </el-upload>
              </div>

              <div
                class="flex items-center h-[120px] cursor-pointer"
                v-if="store.lumaParams.image && store.lumaParams.image_tail"
              >
                <el-tooltip content="交换图片" placement="top">
                  <i class="iconfont icon-exchange" @click="store.switchLumaImages"></i>
                </el-tooltip>
              </div>

              <div class="img-uploader video-img-box ml-2">
                <el-icon
                  v-if="store.lumaParams.image_tail"
                  @click="store.removeLumaImage('end')"
                  class="removeimg"
                >
                  <CircleCloseFilled />
                </el-icon>
                <el-upload
                  class="uploader img-uploader"
                  :auto-upload="true"
                  :show-file-list="false"
                  :http-request="store.uploadLumaEndImage"
                  accept=".jpg,.png,.jpeg"
                >
                  <el-image
                    v-if="store.lumaParams.image_tail"
                    :src="store.lumaParams.image_tail"
                    fit="cover"
                  />
                  <div class="flex flex-col" v-else>
                    <el-icon class="mb-1 text-base"><Plus /></el-icon>
                    <span>结束帧</span>
                  </div>
                </el-upload>
              </div>
            </div>

            <!-- Luma 特有参数设置 -->
            <div class="item-group flex justify-between">
              <span class="label">循环参考图</span>
              <el-switch v-model="store.lumaParams.loop" size="small" />
            </div>

            <div class="item-group flex justify-between">
              <span class="label">提示词优化</span>
              <el-switch v-model="store.lumaParams.expand_prompt" size="small" />
            </div>

            <!-- 算力显示 -->
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
            <div class="flex justify-center">
              <button
                @click="store.createLumaVideo"
                :loading="store.generating"
                class="w-full py-3 bg-gradient-to-r from-blue-500 to-purple-600 text-white rounded-xl disabled:from-gray-400 disabled:to-gray-400 disabled:cursor-not-allowed hover:from-blue-600 hover:to-purple-700 transition-all duration-200 flex items-center justify-center space-x-2 text-base"
              >
                <i v-if="store.generating" class="iconfont icon-loading animate-spin"></i>
                <i v-else class="iconfont icon-chuangzuo"></i>
                <span>立即生成 ({{ store.lumaPowerCost }}算力)</span>
              </button>
            </div>
          </div>
        </el-tab-pane>

        <!-- KeLing 视频参数 -->
        <el-tab-pane label="可灵视频" name="keling">
          <div class="params-container">
            <el-form :model="store.kelingParams" label-width="80px" label-position="left">
              <!-- 画面比例 -->
              <div class="param-line">
                <div class="param-line pt">
                  <span>画面比例：</span>
                  <el-tooltip content="生成画面的尺寸比例" placement="right">
                    <el-icon><InfoFilled /></el-icon>
                  </el-tooltip>
                </div>

                <div class="param-line pt">
                  <el-row :gutter="10">
                    <el-col :span="8" v-for="item in store.rates" :key="item.value">
                      <div
                        class="flex-col items-center"
                        :class="
                          item.value === store.kelingParams.aspect_ratio
                            ? 'grid-content active'
                            : 'grid-content'
                        "
                        @click="store.changeRate(item)"
                      >
                        <el-image class="icon proportion" :src="item.img" fit="cover"></el-image>
                        <div class="texts">{{ item.text }}</div>
                      </div>
                    </el-col>
                  </el-row>
                </div>
              </div>

              <!-- 模型选择 -->
              <div class="param-line">
                <el-form-item label="模型选择">
                  <el-select
                    v-model="store.kelingParams.model"
                    placeholder="请选择模型"
                    @change="store.updateModelPower"
                  >
                    <el-option
                      v-for="item in store.models"
                      :key="item.value"
                      :label="item.text"
                      :value="item.value"
                    />
                  </el-select>
                </el-form-item>
              </div>

              <!-- 视频时长 -->
              <div class="param-line">
                <el-form-item label="视频时长">
                  <el-select
                    v-model="store.kelingParams.duration"
                    placeholder="请选择时长"
                    @change="store.updateModelPower"
                  >
                    <el-option label="5秒" value="5" />
                    <el-option label="10秒" value="10" />
                  </el-select>
                </el-form-item>
              </div>

              <!-- 生成模式 -->
              <div class="param-line">
                <el-form-item label="生成模式">
                  <el-select
                    v-model="store.kelingParams.mode"
                    placeholder="请选择模式"
                    @change="store.updateModelPower"
                  >
                    <el-option label="标准模式" value="std" />
                    <el-option label="专业模式" value="pro" />
                  </el-select>
                </el-form-item>
              </div>

              <!-- 创意程度 -->
              <div class="param-line">
                <el-form-item label="创意程度">
                  <el-slider v-model="store.kelingParams.cfg_scale" :min="0" :max="1" :step="0.1" />
                </el-form-item>
              </div>

              <!-- 运镜控制 -->
              <div class="param-line" v-if="store.showCameraControl">
                <div class="param-line pt">
                  <span>运镜控制：</span>
                  <el-tooltip content="生成画面的运镜效果，仅 1.5的高级模式可用" placement="right">
                    <el-icon><InfoFilled /></el-icon>
                  </el-tooltip>
                </div>

                <div class="param-line">
                  <el-select
                    v-model="store.kelingParams.camera_control.type"
                    placeholder="请选择运镜类型"
                  >
                    <el-option label="请选择" value="" />
                    <el-option label="简单运镜" value="simple" />
                    <el-option label="下移拉远" value="down_back" />
                    <el-option label="推进上移" value="forward_up" />
                    <el-option label="右旋推进" value="right_turn_forward" />
                    <el-option label="左旋推进" value="left_turn_forward" />
                  </el-select>
                </div>

                <!-- 仅在simple模式下显示详细配置 -->
                <div
                  class="camera-control mt-2"
                  v-if="store.kelingParams.camera_control.type === 'simple'"
                >
                  <el-form-item label="水平移动">
                    <el-slider
                      v-model="store.kelingParams.camera_control.config.horizontal"
                      :min="-10"
                      :max="10"
                    />
                  </el-form-item>
                  <el-form-item label="垂直移动">
                    <el-slider
                      v-model="store.kelingParams.camera_control.config.vertical"
                      :min="-10"
                      :max="10"
                    />
                  </el-form-item>
                  <el-form-item label="左右旋转">
                    <el-slider
                      v-model="store.kelingParams.camera_control.config.pan"
                      :min="-10"
                      :max="10"
                    />
                  </el-form-item>
                  <el-form-item label="上下旋转">
                    <el-slider
                      v-model="store.kelingParams.camera_control.config.tilt"
                      :min="-10"
                      :max="10"
                    />
                  </el-form-item>
                  <el-form-item label="横向翻转">
                    <el-slider
                      v-model="store.kelingParams.camera_control.config.roll"
                      :min="-10"
                      :max="10"
                    />
                  </el-form-item>
                  <el-form-item label="镜头缩放">
                    <el-slider
                      v-model="store.kelingParams.camera_control.config.zoom"
                      :min="-10"
                      :max="10"
                    />
                  </el-form-item>
                </div>
              </div>
            </el-form>

            <!-- 图片辅助生成开关 -->
            <div class="param-line pt">
              <div class="image-mode-toggle">
                <span class="label">使用图片辅助生成</span>
                <el-switch
                  v-model="store.kelingUseImageMode"
                  @change="store.toggleKelingImageMode"
                  size="small"
                />
              </div>
            </div>

            <!-- 图片上传区域（可折叠） -->
            <div v-if="store.kelingUseImageMode" class="img-inline">
              <div class="img-uploader video-img-box mr-2">
                <el-icon
                  v-if="store.kelingParams.image"
                  @click="store.removeKelingImage('start')"
                  class="removeimg"
                >
                  <CircleCloseFilled />
                </el-icon>
                <el-upload
                  class="uploader img-uploader"
                  :auto-upload="true"
                  :show-file-list="false"
                  :http-request="store.uploadKelingStartImage"
                  accept=".jpg,.png,.jpeg"
                >
                  <el-image
                    v-if="store.kelingParams.image"
                    :src="store.kelingParams.image"
                    fit="cover"
                  />
                  <div class="flex flex-col" v-else>
                    <el-icon class="mb-1 text-base"><Plus /></el-icon>
                    <span>起始帧</span>
                  </div>
                </el-upload>
              </div>
              <div
                class="flex items-center h-[120px] cursor-pointer"
                v-if="store.kelingParams.image && store.kelingParams.image_tail"
              >
                <el-tooltip content="交换图片" placement="top">
                  <i class="iconfont icon-exchange" @click="store.switchKelingImages"></i>
                </el-tooltip>
              </div>
              <div class="img-uploader video-img-box ml-2">
                <el-icon
                  v-if="store.kelingParams.image_tail"
                  @click="store.removeKelingImage('end')"
                  class="removeimg"
                >
                  <CircleCloseFilled />
                </el-icon>
                <el-upload
                  class="uploader img-uploader"
                  :auto-upload="true"
                  :show-file-list="false"
                  :http-request="store.uploadKelingEndImage"
                  accept=".jpg,.png,.jpeg"
                >
                  <el-image
                    v-if="store.kelingParams.image_tail"
                    :src="store.kelingParams.image_tail"
                    fit="cover"
                  />
                  <div class="flex flex-col" v-else>
                    <el-icon class="mb-1 text-base"><Plus /></el-icon>
                    <span>结束帧</span>
                  </div>
                </el-upload>
              </div>
            </div>

            <!-- 提示词输入 -->
            <div class="param-line pt">
              <span>提示词：</span>
              <el-tooltip content="输入你想要的内容，用逗号分割" placement="right">
                <el-icon><InfoFilled /></el-icon>
              </el-tooltip>
            </div>
            <div class="param-line pt">
              <el-input
                v-model="store.kelingParams.prompt"
                type="textarea"
                maxlength="500"
                :autosize="{ minRows: 4, maxRows: 6 }"
                :placeholder="
                  store.kelingUseImageMode
                    ? '描述视频画面细节'
                    : '请在此输入视频提示词，您也可以点击下面的提示词助手生成视频提示词'
                "
              />
            </div>

            <!-- 提示词生成按钮 -->
            <div class="flex justify-end">
              <el-button
                class="generate-btn"
                @click="store.generatePrompt"
                :loading="store.isGenerating"
                color="#5865f2"
              >
                <span v-if="!store.isGenerating">
                  <i class="iconfont icon-chuangzuo"></i> 生成提示词
                </span>
                <span v-else>生成中...</span>
              </el-button>
            </div>

            <!-- 排除内容 -->
            <div class="param-line pt">
              <span>不希望出现的内容：（可选）</span>
              <el-tooltip content="不想出现在图片上的元素(例如：树，建筑)" placement="right">
                <el-icon><InfoFilled /></el-icon>
              </el-tooltip>
            </div>
            <div class="param-line pt">
              <el-input
                v-model="store.kelingParams.negative_prompt"
                type="textarea"
                :autosize="{ minRows: 4, maxRows: 6 }"
                placeholder="请在此输入你不希望出现在视频上的内容"
              />
            </div>

            <!-- 算力显示 -->
            <!-- 算力显示 -->
            <div
              class="power-info flex items-center justify-between mb-4 mt-2 p-3 rounded-lg bg-gradient-to-r from-blue-50 to-purple-50 border border-blue-200 shadow-sm"
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

            <!-- 生成按钮 -->
            <div class="flex justify-center">
              <button
                @click="store.createKelingVideo"
                :loading="store.generating"
                class="w-full py-3 bg-gradient-to-r from-blue-500 to-purple-600 text-white rounded-xl disabled:from-gray-400 disabled:to-gray-400 disabled:cursor-not-allowed hover:from-blue-600 hover:to-purple-700 transition-all duration-200 flex items-center justify-center space-x-2 text-base"
              >
                <i v-if="store.generating" class="iconfont icon-loading animate-spin"></i>
                <i v-else class="iconfont icon-chuangzuo"></i>
                <span>立即生成 ({{ store.kelingPowerCost }}算力)</span>
              </button>
            </div>
          </div>
        </el-tab-pane>
      </el-tabs>
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
              :type="store.taskFilter === 'luma' ? 'primary' : 'default'"
              @click="store.switchTaskFilter('luma')"
              size="small"
            >
              Luma
            </el-button>
            <el-button
              :type="store.taskFilter === 'keling' ? 'primary' : 'default'"
              @click="store.switchTaskFilter('keling')"
              size="small"
            >
              可灵
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
                  <div v-if="item.progress === 100">
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
                  <el-image
                    :src="item.cover_url"
                    class="border rounded-lg"
                    fit="cover"
                    v-else-if="item.progress === 101"
                  />
                  <generating message="正在生成视频" v-else />
                </div>
              </div>

              <div class="center">
                <div class="pb-2" v-if="item.raw_data">
                  <el-tag class="mr-1">{{
                    item.raw_data.task_type || store.activeVideoType
                  }}</el-tag>
                  <el-tag class="mr-1" v-if="item.raw_data.model">{{ item.raw_data.model }}</el-tag>
                  <el-tag class="mr-1" v-if="item.raw_data.duration"
                    >{{ item.raw_data.duration }}秒</el-tag
                  >
                  <el-tag class="mr-1" v-if="item.raw_data.mode">{{ item.raw_data.mode }}</el-tag>
                </div>
                <div class="failed" v-if="item.progress === 101">
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
      @close="store.showDialog = false"
      width="auto"
    >
      <video
        style="max-width: 90vw; max-height: 90vh"
        :src="store.currentVideoUrl"
        preload="auto"
        :autoplay="true"
        loop="loop"
        muted="muted"
      >
        您的浏览器不支持视频播放
      </video>
    </el-dialog>
  </div>
</template>

<script setup>
import Generating from '@/components/ui/Generating.vue'
import { useVideoStore } from '@/store/video'
import { CircleCloseFilled, InfoFilled, Plus } from '@element-plus/icons-vue'
import { onMounted, onUnmounted } from 'vue'

const store = useVideoStore()

onMounted(() => {
  store.init()
})

onUnmounted(() => {
  store.cleanup()
})
</script>

<style lang="scss" scoped>
@use '../assets/css/video.scss' as *;
</style>
