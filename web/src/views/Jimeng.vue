<template>
  <div class="page-jimeng">
    <!-- 左侧参数设置面板 -->
    <div class="params-panel">
      <!-- 功能分类按钮组 -->
      <div class="category-buttons">
        <div class="category-grid">
          <button
            v-for="f in store.functions"
            :key="f.key"
            class="category-btn text-base"
            :class="{ active: store.activeFunction === f.key }"
            @click="store.switchFunction(f)"
          >
            <i class="iconfont mr-2 !text-xl" :class="f.icon"></i>
            {{ f.name }}
          </button>
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
      <div class="function-params">
        <div class="mb-3">
          <div class="mb-2" v-if="store.functionParams[store.activeFunction].length > 0">
            <label class="label text-left font-bold">模型选择</label>
          </div>
          <param-builder
            v-model="store.formData"
            v-model:required-keys="store.requiredKeys"
            :items="store.functionParams[store.activeFunction]"
            :progress="store.progress[store.activeFunction]"
          />
        </div>

        <!-- 提交按钮 -->
        <div
          class="submit-btn flex justify-center pt-4"
          v-if="store.functionParams[store.activeFunction].length > 0"
        >
          <button
            @click="store.submitTask"
            :disabled="store.submitting"
            class="w-full py-3 bg-gradient-to-r from-blue-500 to-purple-600 text-white rounded-xl disabled:from-gray-400 disabled:to-gray-400 disabled:cursor-not-allowed hover:from-blue-600 hover:to-purple-700 transition-all duration-200 flex items-center justify-center space-x-2 text-base"
            type="button"
          >
            <i v-if="store.submitting" class="iconfont icon-loading animate-spin"></i>
            <i v-else class="iconfont icon-chuangzuo"></i>
            <span>立即生成 ({{ store.currentPowerCost }})</span>
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
                    <span v-if="item.power">{{ item.power }}积分</span>
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
import loadingIcon from '@/assets/img/loading.gif'

import ParamBuilder from '@/components/ParamBuilder.vue'
import Generating from '@/components/ui/Generating.vue'
import { useJimengStore } from '@/store/jimeng'
import { useSharedStore } from '@/store/sharedata'
import { dateFormat } from '@/utils/libs'
import { ElMessage } from 'element-plus'
import { onMounted, onUnmounted, ref, watch } from 'vue'
import { Waterfall } from 'vue-waterfall-plugin-next'
import 'vue-waterfall-plugin-next/dist/style.css'

const sharedStore = useSharedStore()
const waterfallOptions = sharedStore.waterfallOptions

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
@use '@/assets/css/jimeng.scss' as *;
</style>
