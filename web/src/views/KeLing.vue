<template>
  <div class="page-keling">
    <div class="inner custom-scroll">
      <!-- 左侧参数设置面板 -->
      <el-scrollbar max-height="100vh">
        <div class="mj-box">
          <h2>视频参数设置</h2>
          <el-form :model="params" label-width="80px" label-position="left">
            <!-- 画面比例 -->
            <div class="param-line">
              <div class="param-line pt">
                <span>画面比例：</span>
                <el-tooltip content="生成画面的尺寸比例" placement="right">
                  <el-icon>
                    <InfoFilled />
                  </el-icon>
                </el-tooltip>
              </div>

              <div class="param-line pt">
                <el-row :gutter="10">
                  <el-col :span="8" v-for="item in rates" :key="item.value">
                    <div
                      class="flex-col items-center"
                      :class="
                        item.value === params.aspect_ratio
                          ? 'grid-content active'
                          : 'grid-content'
                      "
                      @click="changeRate(item)"
                    >
                      <el-image
                        class="icon proportion"
                        :src="item.img"
                        fit="cover"
                      ></el-image>
                      <div class="texts">{{ item.text }}</div>
                    </div>
                  </el-col>
                </el-row>
              </div>
            </div>
            <!-- 模型选择 -->
            <div class="param-line">
              <el-form-item label="模型选择">
                <el-select v-model="params.model" placeholder="请选择模型">
                  <el-option label="默认模型" value="default" />
                  <el-option label="动漫风格" value="anime" />
                </el-select>
              </el-form-item>
            </div>

            <!-- 视频时长 -->
            <div class="param-line">
              <el-form-item label="视频时长">
                <el-select v-model="params.duration" placeholder="请选择时长">
                  <el-option label="5秒" value="5" />
                  <el-option label="10秒" value="10" />
                </el-select>
              </el-form-item>
            </div>

            <!-- 生成模式 -->
            <div class="param-line">
              <el-form-item label="生成模式">
                <el-select v-model="params.mode" placeholder="请选择模式">
                  <el-option label="标准模式" value="std" />
                  <el-option label="专业模式" value="pro" />
                </el-select>
              </el-form-item>
            </div>

            <!-- 创意程度 -->
            <div class="param-line">
              <el-form-item label="创意程度">
                <el-slider
                  v-model="params.cfg_scale"
                  :min="0"
                  :max="1"
                  :step="0.1"
                />
              </el-form-item>
            </div>

            <!-- 运镜控制 -->
            <div class="param-line">
              <div class="param-line pt">
                <span>运镜控制：</span>
                <el-tooltip content="生成画面的运镜效果" placement="right">
                  <el-icon>
                    <InfoFilled />
                  </el-icon>
                </el-tooltip>
              </div>

              <!-- 添加运镜类型选择 -->
              <el-form-item label="运镜类型">
                <el-select
                  v-model="params.camera_control.type"
                  placeholder="请选择运镜类型"
                >
                  <el-option label="请选择" value="" />
                  <el-option label="简单运镜" value="simple" />
                  <el-option label="下移拉远" value="down_back" />
                  <el-option label="推进上移" value="forward_up" />
                  <el-option label="右旋推进" value="right_turn_forward" />
                  <el-option label="左旋推进" value="left_turn_forward" />
                </el-select>
              </el-form-item>

              <!-- 仅在simple模式下显示详细配置 -->
              <div
                class="camera-control"
                v-if="params.camera_control.type === 'simple'"
              >
                <el-form-item label="水平移动">
                  <el-slider
                    v-model="params.camera_control.config.horizontal"
                    :min="-10"
                    :max="10"
                  />
                </el-form-item>
                <el-form-item label="垂直移动">
                  <el-slider
                    v-model="params.camera_control.config.vertical"
                    :min="-10"
                    :max="10"
                  />
                </el-form-item>
                <el-form-item label="左右旋转">
                  <el-slider
                    v-model="params.camera_control.config.pan"
                    :min="-10"
                    :max="10"
                  />
                </el-form-item>
                <el-form-item label="上下旋转">
                  <el-slider
                    v-model="params.camera_control.config.tilt"
                    :min="-10"
                    :max="10"
                  />
                </el-form-item>
                <el-form-item label="横向翻转">
                  <el-slider
                    v-model="params.camera_control.config.roll"
                    :min="-10"
                    :max="10"
                  />
                </el-form-item>
                <el-form-item label="镜头缩放">
                  <el-slider
                    v-model="params.camera_control.config.zoom"
                    :min="-10"
                    :max="10"
                  />
                </el-form-item>
              </div>
            </div>
          </el-form>
        </div>
      </el-scrollbar>

      <!-- 右侧主内容区 -->
      <div class="main-content task-list-inner">
        <!-- 任务类型选择 -->
        <div class="param-line">
          <el-tabs
            v-model="params.task_type"
            @tab-change="tabChange"
            class="title-tabs"
          >
            <el-tab-pane label="文生视频" name="text2video">
              <div class="text">使用文字描述想要生成视频的内容</div>
            </el-tab-pane>
            <el-tab-pane label="图生视频" name="image2video">
              <div class="text">
                以某张图片为底稿参考来创作视频，生成类似风格或类型视频，支持 PNG
                /JPG/JPEG 格式图片；
              </div>
            </el-tab-pane>
          </el-tabs>
        </div>

        <!-- 生成操作区 -->
        <div class="generation-area">
          <div v-if="params.task_type === 'text2video'" class="text2video">
            <el-input
              v-model="params.prompt"
              type="textarea"
              :autosize="{ minRows: 4, maxRows: 6 }"
              placeholder="请在此输入视频提示词，您也可以点击下面的提示词助手生成视频提示词"
            />
            <el-row class="text-info">
              <el-button
                class="generate-btn"
                @click="generatePrompt"
                :loading="isGenerating"
                size="small"
                color="#5865f2"
              >
                <i class="iconfont icon-chuangzuo"></i>
                生成专业视频提示词
              </el-button>
            </el-row>
          </div>

          <div v-else class="image2video">
            <div class="image-upload img-inline">
              <div class="upload-box img-uploader">
                <h4>起始帧</h4>
                <el-upload
                  class="uploader img-uploader"
                  :auto-upload="true"
                  :show-file-list="false"
                  :http-request="uploadStartImage"
                  accept=".jpg,.png,.jpeg"
                >
                  <img
                    v-if="params.image"
                    :src="params.image"
                    class="preview"
                  />
                  <el-icon v-else class="upload-icon"><Plus /></el-icon>
                </el-upload>
              </div>
              <div class="upload-box img-uploader">
                <h4>结束帧</h4>
                <el-upload
                  class="uploader"
                  :auto-upload="true"
                  :show-file-list="false"
                  :http-request="uploadEndImage"
                  accept=".jpg,.png,.jpeg"
                >
                  <img
                    v-if="params.image_tail"
                    :src="params.image_tail"
                    class="preview"
                  />
                  <el-icon v-else class="upload-icon"><Plus /></el-icon>
                </el-upload>
              </div>
            </div>
            <div class="param-line pt">
              <div class="flex-row justify-between items-center">
                <div class="flex-row justify-start items-center">
                  <span>提示词：</span>
                  <el-tooltip
                    content="输入你想要的内容，用逗号分割"
                    placement="right"
                  >
                    <el-icon>
                      <InfoFilled />
                    </el-icon>
                  </el-tooltip>
                </div>
              </div>
            </div>
            <div class="param-line pt">
              <el-input
                v-model="params.prompt"
                type="textarea"
                :autosize="{ minRows: 4, maxRows: 6 }"
                placeholder="描述视频画面细节"
              />
            </div>
          </div>

          <!-- 排除内容 -->
          <div class="param-line pt">
            <div class="flex-row justify-between items-center">
              <div class="flex-row justify-start items-center">
                <span>不希望出现的内容：（可选）</span>
                <el-tooltip
                  content="不想出现在图片上的元素(例如：树，建筑)"
                  placement="right"
                >
                  <el-icon>
                    <InfoFilled />
                  </el-icon>
                </el-tooltip>
              </div>
            </div>
          </div>
          <div class="param-line pt">
            <el-input
              v-model="params.negative_prompt"
              type="textarea"
              :autosize="{ minRows: 4, maxRows: 6 }"
              placeholder="请在此输入你不希望出现在视频上的内容"
            />
          </div>

          <!-- 算力显示 -->
          <el-row class="text-info">
            <el-text type="primary"
              >每次生成视频消耗
              <el-text type="warning">{{ powerCost }}算力;</el-text> </el-text
            >&nbsp;&nbsp;
            <el-text type="primary"
              >当前可用算力：<el-text type="warning">{{
                availablePower
              }}</el-text></el-text
            >
          </el-row>

          <!-- 生成按钮 -->
          <div class="submit-btn">
            <el-button type="primary" :dark="false" @click="generate" round
              >立即生成</el-button
            >
          </div>
        </div>

        <!-- 任务列表区域 -->
        <div class="task-list job-list-box">
          <!-- 任务类型筛选 -->
          <div class="type-btn-group">
            <el-radio-group v-model="taskFilter" @change="fetchTasks">
              <el-radio-button label="all">全部</el-radio-button>
              <el-radio-button label="text2video">文生视频</el-radio-button>
              <el-radio-button label="image2video">图生视频</el-radio-button>
            </el-radio-group>
          </div>

          <h2 class="text-xl pt">任务列表</h2>

          <!-- 运行中的任务 -->
          <task-list :list="runningTasks" />
          <template v-if="finishedTasks.length > 0">
            <h2 class="record-title pt">创作记录</h2>
            <!-- 已完成的任务 -->
            <v3-waterfall
              :virtual-time="200"
              :distance-to-scroll="150"
              :key="waterfallKey"
              :list="finishedTasks"
              @scrollReachBottom="fetchTasks"
              :gap="8"
              :bottomGap="8"
              :colWidth="300"
              :distanceToScroll="100"
              :isLoading="loading"
              :isOver="isOver"
              class="task-waterfall"
            >
              <template #default="slotProp">
                <!-- 视频成功渲染部分 -->
                <div
                  class="job-item-box"
                  :class="{
                    processing: slotProp.item.progress < 100,
                    error: slotProp.item.progress === 101
                  }"
                >
                  <video
                    v-if="
                      slotProp.item.progress >= 100 && slotProp.item.video_url
                    "
                    class="preview"
                    :src="slotProp.item.video_url"
                    @click="previewVideo(slotProp.item)"
                    controls
                    :style="{
                      width: '100%',
                      height: `${slotProp.item.height || 400}px`
                    }"
                  ></video>

                  <!-- 失败/无图状态 -->
                  <div
                    v-else
                    class="error-container"
                    :style="{
                      width: '100%',
                      height: `${slotProp.item.height || 300}px`,
                      objectFit: 'cover'
                    }"
                  >
                    <div
                      v-if="
                        slotProp.item.progress >= 100 &&
                        !slotProp.item.video_url
                      "
                      class="error-status"
                    >
                      <img :src="failed" />
                      生成失败
                    </div>
                    <div v-else class="processing-status">
                      <el-progress
                        :percentage="slotProp.item.progress"
                        :stroke-width="12"
                        status="success"
                      />
                    </div>
                  </div>
                  <div class="tools-box">
                    <div class="tools">
                      <el-button
                        type="primary"
                        v-if="
                          slotProp.item.progress >= 100 &&
                          slotProp.item.video_url
                        "
                        @click="downloadVideo(slotProp.item)"
                      >
                        <el-icon><Download /></el-icon>
                      </el-button>

                      <div
                        class="show-prompt"
                        v-if="
                          slotProp.item.progress >= 100 &&
                          !slotProp.item.video_url &&
                          slotProp.item.err_msg
                        "
                      >
                        <el-popover
                          placement="left"
                          :width="240"
                          trigger="hover"
                        >
                          <template #reference>
                            <el-icon class="chromefilled error-txt"
                              ><WarnTriangleFilled
                            /></el-icon>
                          </template>

                          <template #default>
                            <div class="top-tips">
                              <span>错误详细信息</span
                              ><el-icon
                                class="copy-prompt-kl"
                                :data-clipboard-text="slotProp.item.prompt"
                              >
                                <DocumentCopy />
                              </el-icon>
                            </div>
                            <div class="mj-list-item-prompt">
                              <span>{{ slotProp.item.prompt }}</span>
                            </div>
                          </template>
                        </el-popover>
                      </div>

                      <el-button
                        type="danger"
                        @click="deleteTask(slotProp.item)"
                      >
                        <el-icon><Delete /></el-icon>
                      </el-button>
                      <div class="show-prompt">
                        <el-popover
                          placement="left"
                          :width="240"
                          trigger="hover"
                        >
                          <template #reference>
                            <el-icon class="chromefilled">
                              <ChromeFilled />
                            </el-icon>
                          </template>

                          <template #default>
                            <div class="top-tips">
                              <span>提示词</span
                              ><el-icon
                                class="copy-prompt-kl"
                                :data-clipboard-text="slotProp.item.prompt"
                              >
                                <DocumentCopy />
                              </el-icon>
                            </div>
                            <div class="mj-list-item-prompt">
                              <span>{{ slotProp.item.prompt }}</span>
                            </div>
                          </template>
                        </el-popover>
                      </div>
                    </div>
                  </div>
                </div>
              </template>
              <template #footer>
                <div class="no-more-data">
                  <span>没有更多数据了</span>
                  <i class="iconfont icon-face"></i>
                </div>
              </template>
            </v3-waterfall>
          </template>
        </div>
      </div>
    </div>

    <!-- 视频预览对话框 -->
    <el-dialog v-model="previewVisible" title="视频预览" width="80%">
      <video
        v-if="currentVideo"
        :src="currentVideo"
        controls
        style="width: 100%"
      ></video>
    </el-dialog>
  </div>
</template>

<script setup>
import failed from "@/assets/img/failed.png";
import TaskList from "@/components/TaskList.vue";
import { ref, reactive, onMounted, onUnmounted, watch, computed } from "vue";
import {
  Plus,
  Delete,
  InfoFilled,
  ChromeFilled,
  DocumentCopy,
  Download,
  WarnTriangleFilled
} from "@element-plus/icons-vue";
import { httpGet, httpPost, httpDownload } from "@/utils/http";
import { ElMessage, ElMessageBox } from "element-plus";
import { getClientId, checkSession } from "@/store/cache";
import Clipboard from "clipboard";

import {
  closeLoading,
  showLoading,
  showMessageError,
  showMessageOK
} from "@/utils/dialog";
import { replaceImg } from "@/utils/libs";

// 参数设置
const params = reactive({
  client_id: getClientId(),
  task_type: "text2video",
  model: "default",
  prompt: "",
  negative_prompt: "",
  cfg_scale: 0.7,
  mode: "std",
  aspect_ratio: "16:9",
  duration: "5",
  camera_control: {
    type: "",
    config: {
      horizontal: 0,
      vertical: 0,
      pan: 0,
      tilt: 0,
      roll: 0,
      zoom: 0
    }
  },
  image: "",
  image_tail: ""
});
const rates = [
  { css: "square", value: "1:1", text: "1:1", img: "/images/mj/rate_1_1.png" },

  {
    css: "size16-9",
    value: "16:9",
    text: "16:9",
    img: "/images/mj/rate_16_9.png"
  },
  {
    css: "size9-16",
    value: "9:16",
    text: "9:16",
    img: "/images/mj/rate_9_16.png"
  }
];

// 切换图片比例
const changeRate = (item) => {
  params.aspect_ratio = item.value;
};

// 状态变量
const generating = ref(false);
const isGenerating = ref(false);
const powerCost = ref(10);
const availablePower = ref(100);
const taskFilter = ref("all");
const runningTasks = ref([]);
const finishedTasks = ref([]);
const total = ref(0);
const pageSize = ref(15);
const currentPage = ref(1);
const previewVisible = ref(false);
const currentVideo = ref("");
const isOver = ref(false);

// 方法定义

const tabChange = (tab) => {
  params.task_type = tab;
};

const uploadStartImage = async (file) => {
  const formData = new FormData();
  formData.append("file", file.file);
  try {
    showLoading("图片上传中...");
    const res = await httpPost("/api/upload", formData);
    params.image = res.data.url;
    ElMessage.success("上传成功");
    closeLoading();
  } catch (e) {
    showMessageError("上传失败: " + e.message);
    closeLoading();
  }
};

const uploadEndImage = async (file) => {
  const formData = new FormData();
  formData.append("file", file.file);
  try {
    const res = await httpPost("/api/upload", formData);
    params.image_tail = res.data.url;
    ElMessage.success("上传成功");
  } catch (e) {
    showMessageError("上传失败: " + e.message);
  }
};

const generatePrompt = async () => {
  if (!params.prompt) {
    return showMessageError("请输入视频描述");
  }
  isGenerating.value = true;
  try {
    const res = await httpPost("/api/prompt/video", { prompt: params.prompt });
    params.prompt = res.data;
  } catch (e) {
    showMessageError("生成失败: " + e.message);
  } finally {
    isGenerating.value = false;
  }
};

const generate = async () => {
  if (!params.prompt?.trim()) {
    return ElMessage.error("请输入视频描述");
  }
  // if (params.task_type === "image2video" && !params.image) {
  //   return ElMessage.error("请上传起始帧图片");
  // }

  generating.value = true;
  // 处理图片链接
  params.image = replaceImg(params.image);
  params.image_tail = replaceImg(params.image_tail);
  try {
    await httpPost("/api/video/keling/create", params);
    showMessageOK("任务创建成功");
    // 立即获取最新数据
    fetchTasks();
  } catch (e) {
    showMessageError("创建失败: " + e.message);
  } finally {
    generating.value = false;
  }
};
const loading = ref(false);
const isLogin = ref(false);

const fetchTasks = async () => {
  console.log("fetchTasks", !isLogin.value || isOver.value || loading.value);

  if (!isLogin.value || isOver.value || loading.value) return;

  loading.value = true;
  try {
    const res = await httpGet("/api/video/list", {
      page: currentPage.value,
      page_size: pageSize.value,
      type: "keling",
      task_type: taskFilter.value === "all" ? "" : taskFilter.value
    });

    // 精确任务过滤逻辑
    const data = res.data || {};
    const newRunning = data.items.filter(
      (task) => task.progress < 100 && task.progress !== 101
    );
    runningTasks.value = [...runningTasks.value, ...newRunning];

    const newfinished = data.items.filter((task) => task.progress >= 100);
    const finishedList = [...finishedTasks.value, ...newfinished];
    finishedTasks.value = finishedList.map((item) => ({
      ...item,
      height: 300 * (Math.random() * 0.4 + 0.6) // 生成300~420px随机高度
    }));
    console.log("finishedTasks: " + finishedList);

    // // 强制刷新瀑布流
    waterfallKey.value = Date.now();
    total.value = data.total;

    const shouldLoadNextPage =
      runningTasks.value.length > 0 ||
      (runningTasks.value.length === 0 &&
        finishedTasks.value.length < total.value);

    if (shouldLoadNextPage) {
      currentPage.value++;
    }

    // 优化加载完成判断
    const loadedCount = runningTasks.value.length + finishedTasks.value.length;
    isOver.value = loadedCount >= total.value;

    console.log("[Pagination] isOver:", isOver.value, {
      running: runningTasks.value.length,
      finished: finishedTasks.value.length,
      total: total.value
    });

    waterfallKey.value = Date.now();
  } catch (e) {
    showMessageError("获取任务列表失败: " + e.message);
  } finally {
    loading.value = false;
  }
};

const waterfallKey = ref(Date.now());

// 检测任务状态
const checkAllCompleted = () => {
  return runningTasks.value.length === 0;
};

const previewVideo = (task) => {
  currentVideo.value = task.video_url;
  previewVisible.value = true;
};

const downloadVideo = async (task) => {
  try {
    const res = await httpDownload(`/api/download?url=${task.video_url}`);
    const blob = new Blob([res.data]);
    const link = document.createElement("a");
    link.href = URL.createObjectURL(blob);
    link.download = `video_${task.id}.mp4`;
    link.click();
    URL.revokeObjectURL(link.href);
  } catch (e) {
    showMessageError("下载失败: " + e.message);
  }
};

const deleteTask = async (task) => {
  try {
    await ElMessageBox.confirm("确定要删除该任务吗？");
    await httpGet("/api/video/remove", { id: task.id });
    showMessageOK("删除成功");
    fetchTasks();
  } catch (e) {
    if (e !== "cancel") {
      showMessageError("删除失败: " + e.message);
    }
  }
};

fetchTasks();
const clipboard = ref(null);

// 生命周期钩子
onMounted(async () => {
  checkSession()
    .then(async () => {
      isLogin.value = true;
      console.log("mounted-isLogin-可以继续", isLogin.value);
      await fetchTasks();
    })
    .catch(() => {});

  // fetchTasks();
  clipboard.value = new Clipboard(".copy-prompt-kl");
  clipboard.value.on("success", () => {
    ElMessage.success("复制成功！");
  });

  clipboard.value.on("error", () => {
    ElMessage.error("复制失败！");
  });
});
onUnmounted(() => {
  clipboard.value.destroy();
});
// 监听任务状态变化
watch([runningTasks, finishedTasks], () => {
  if (checkAllCompleted()) {
    console.log("所有任务已完成");
  }
});
</script>

<style lang="stylus" scoped>
@import "@/assets/css/image-keling.styl"
@import "@/assets/css/custom-scroll.styl"
.copy-prompt-kl{
  cursor pointer
}
.top-tips{
  height: 30px
  font-size: 18px
  line-height: 30px
  display: flex
  align-items: center;
  span{
    margin-right: 10px
    color:#000
  }
}
.mj-list-item-prompt{
  max-height: 600px;
  overflow: auto;
}
:deep(.running-job-box .image-slot){
  display: flex
  align-items: center
  flex-direction: column;
  justify-content: center
  width: 100%
  height: 100%
  font-size: 30px
  color: #909399
  background: #f5f7fa
  text-align: center
  width: 200px;
  height: 200px;
  .iconfont{
    font-size: 45px;

  }
  span{
    font-size: 15px
  }
}


.record-title
  padding:1rem 0
.type-btn-group
  margin-bottom: 20px
.task-waterfall
  margin: 0 -10px
  transition: opacity 0.3s ease

.job-item-box
  position: relative
  background: #f5f5f5;
  transition: height 0.3s ease;
  overflow: hidden
  // margin: 10px
  // border: 1px solid #666;
  // padding: 6px;
  border-radius: 6px;
  break-inside: avoid
  video
    min-height: 200px;
    width: 100%;
    object-fit: cover;

  .chromefilled
    font-size: 24px;
    color: #fff;
    &.error-txt{
      color: #ffff54;
      cursor:pointer;
    }
  .show-prompt
    display: flex;
    align-items: center;
  &:hover
    // transform: translateY(-3px)
    .tools-box{
      display:block
      background:rgba(0, 0, 0, 0.3)
      width : 100%;
    }

  .error-container
    position: relative
    background: var(--bg-deep-color)
    display: flex
    align-items: center
    justify-content: center
    img{
      width: 66%;
      height: 66%;
      object-fit: cover;
      margin: 0 auto;
    }

    .error-status
      color: #c2c6cc
      text-align: center
      font-size: 24px


    .processing-status
      width: 80%
      .el-progress
        margin: 0 auto
  .tools-box{
    display:none
    position:absolute;
    top: 0;
    right: 0;
  }
  .tools
    align-items: center;
    justify-content: flex-end;
    display: flex
    gap: 5px
    margin: 5px 5px 5px 0;


    .el-button+.el-button
      margin-left: 0px;

    .el-button
      padding: 3px
      border-radius: 50%
</style>
