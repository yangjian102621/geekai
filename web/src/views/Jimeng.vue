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

      <!-- 提示词编写指南（可折叠） -->
      <div class="prompt-guide">
        <el-collapse v-model="guideActive">
          <el-collapse-item name="guide">
            <template #title>
              <div class="guide-title">
                <i class="iconfont icon-info mr-1"></i>
                Prompt建议
              </div>
            </template>
            <div class="guide-content">
              <!-- 创建图像 -->
              <div class="guide-section">
                <div class="guide-subtitle">创建图像（文生图）</div>
                <ul>
                  <li>
                    结构建议：<strong>主体描述 + 风格 + 美学</strong>（准确响应）；<strong
                      >风格 + 主体描述 + 美学 + 氛围</strong
                    >（更强美学）
                  </li>
                  <li>
                    用专业短词描述风格/镜头/构图；主体用自然语言完整描述（主体 + 行为 + 环境）
                  </li>
                  <li>关键信息靠前；用正向表达代替“不要xxx”类否定词</li>
                  <li>需要生成文字时，明确“生成文字”并补充位置/风格/材质</li>
                </ul>
                <blockquote class="quote">
                  <div>
                    <strong>示例</strong
                    >：新年主题海报，上方以手写涂鸦风格写着“新年快乐”，红金配色，纸张纹理，强对比光影，居中极简构图，留白用于标题。
                  </div>
                </blockquote>
                <blockquote class="quote">
                  <div><strong>Before</strong>：海报，“新年快乐”</div>
                  <div><strong>After</strong>：一张海报，上面文字写着：“新年快乐”</div>
                </blockquote>
                <blockquote class="quote">
                  <div><strong>Before</strong>：海报，“新年快乐”</div>
                  <div>
                    <strong>After</strong>：一张海报，画面上方有手写涂鸦风格的文字写着：“新年快乐”
                  </div>
                </blockquote>
                <p>特征与视角可反复强调：</p>
                <ul>
                  <li>御剑飞行 → 男人站在剑上，他踩在剑上，剑被他踩着，御剑飞仙</li>
                  <li>仰视视角 → 采用低角度，从下往上，仰视与广角构图</li>
                </ul>
                <blockquote class="quote">
                  <div>
                    <strong>示例</strong
                    >：百合南瓜羹特写，只展示半碗，米黄色糯米勾芡，橙色南瓜块与丝理清晰，点缀紫白色百合。
                  </div>
                </blockquote>
              </div>

              <!-- 编辑图像 -->
              <div class="guide-section">
                <div class="guide-subtitle">编辑图像（图生图/图像编辑）</div>
                <ul>
                  <li>应用场景可加：如“海报、平面设计”等词以增强对应风格</li>
                  <li>生成或保留的文字请用引号包裹，准确率更高</li>
                  <li>建议长度 ≤ 120 字，最多不超过 800 字，过长可能失效</li>
                  <li>编辑指令用自然语言；一次只做一件事更易生效</li>
                  <li>多实体时指明“对谁做什么”，局部编辑尽量精准</li>
                  <li>效果不明显可提高编辑强度 scale；底图越清晰效果越好</li>
                </ul>
                <blockquote class="quote">
                  <div><strong>示例（添加/删除实体）</strong>：删除图上的女孩；添加一道彩虹。</div>
                </blockquote>
                <blockquote class="quote">
                  <div>
                    <strong>示例（添加文字）</strong>：一张圣诞节海报，上面写着“Merry Christmas”。
                  </div>
                </blockquote>
                <blockquote class="quote">
                  <div><strong>示例（修改实体）</strong>：把手里的鸡腿改成汉堡。</div>
                </blockquote>
                <blockquote class="quote">
                  <div>
                    <strong>示例（修改风格/色彩/动作/背景）</strong
                    >：改成漫画风格；把外套改成粉色；让男孩微笑；背景换成海边日落。
                  </div>
                </blockquote>
              </div>

              <!-- 生成视频 -->
              <div class="guide-section">
                <div class="guide-subtitle">生成视频（文/图生视频）</div>
                <ul>
                  <li><strong>基础结构</strong>：主体 / 背景 / 镜头 + 动作</li>
                  <li>
                    <strong>多个镜头连贯叙事</strong>：镜头1 + 主体 + 动作1 + 镜头2 + 主体 + 动作2
                    ...
                  </li>
                  <li>
                    <strong>多个连续动作</strong>：时序：主体1 + 运动1 + 运动2；多主体：主体1 +
                    运动1 + 主体2 + 运动2 ...
                  </li>
                  <li>
                    <strong>运镜词典</strong
                    >：镜头切换；镜头向上/下/左/右移动；镜头拉近/拉远；镜头环绕/航拍/广角/360度旋转；镜头跟随；固定镜头；镜头特写；手持拍摄（晃动/抖动）
                  </li>
                  <li>
                    <strong>程度副词</strong
                    >：快速、缓缓、大幅度、高频率、剧烈等，突出动作强度与节奏
                  </li>
                </ul>
                <blockquote class="quote">
                  <div>
                    <strong>示例</strong
                    >：镜头1，城市夜景航拍，镜头环绕；镜头2，男主在屋顶奔跑，镜头跟随，快速；镜头3，男主停下特写，霓虹反光，缓慢拉近。
                  </div>
                </blockquote>
              </div>
            </div>
          </el-collapse-item>
        </el-collapse>
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

// 新增：提示词指南折叠面板状态（默认收起）
const guideActive = ref([])

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
  line-clamp: 2;
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

// 新增：提示词指南样式
.prompt-guide {
  margin: 12px 0 16px;

  .guide-title {
    display: flex;
    align-items: center;
    font-weight: 600;
    color: #666;
  }

  .guide-content {
    max-height: 220px;
    overflow: auto;
    line-height: 1.6;
    font-size: 12px;
    color: #555;
    padding-right: 4px;
  }

  .guide-section {
    margin-bottom: 10px;
  }

  .guide-subtitle {
    font-weight: 600;
    margin-bottom: 6px;
    color: #333;
  }

  ul {
    list-style: disc;
    padding-left: 18px;
    margin: 4px 0;
  }

  .quote {
    margin: 8px 0;
    padding: 8px 10px;
    border-left: 3px solid #a3a3a3;
    background: #f8f8f8;
    border-radius: 4px;
    color: #444;
  }
}
</style>
