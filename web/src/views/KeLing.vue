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
                <el-select
                  v-model="params.model"
                  placeholder="请选择模型"
                  @change="updateModelPower"
                >
                  <el-option
                    v-for="item in models"
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
                  v-model="params.duration"
                  placeholder="请选择时长"
                  @change="updateModelPower"
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
                  v-model="params.mode"
                  placeholder="请选择模式"
                  @change="updateModelPower"
                >
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
            <div class="param-line" v-if="showCameraControl">
              <div class="param-line pt">
                <span>运镜控制：</span>
                <el-tooltip
                  content="生成画面的运镜效果，仅 1.5的高级模式可用"
                  placement="right"
                >
                  <el-icon>
                    <InfoFilled />
                  </el-icon>
                </el-tooltip>
              </div>

              <!-- 添加运镜类型选择 -->
              <div class="param-line">
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
              </div>

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
              maxlength="500"
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
              <div class="upload-box img-uploader video-img-box">
                <el-icon
                  v-if="params.image"
                  @click="remove('start', params.image)"
                  class="removeimg"
                  ><CircleCloseFilled
                /></el-icon>

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
              <div class="btn-swap" v-if="params.image && params.image_tail">
                <i class="iconfont icon-exchange" @click="switchReverse"></i>
              </div>
              <div class="upload-box img-uploader video-img-box">
                <el-icon
                  v-if="params.image_tail"
                  @click="remove('end', params.image_tail)"
                  class="removeimg"
                  ><CircleCloseFilled
                /></el-icon>
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
          <!-- <div class="type-btn-group">
            <el-radio-group v-model="taskFilter" @change="fetchTasks">
              <el-radio-button label="all">全部</el-radio-button>
              <el-radio-button label="text2video">文生视频</el-radio-button>
              <el-radio-button label="image2video">图生视频</el-radio-button>
            </el-radio-group>
          </div> -->

          <h2 class="text-xl pt">任务列表</h2>

          <!-- 运行中的任务 -->
          <task-list :list="runningTasks" />
          <template v-if="finishedTasks.length > 0">
            <h2 class="record-title pt">创作记录</h2>
            <!-- 已完成的任务 -->
            <v3-waterfall
              :key="waterfallKey"
              :list="finishedTasks"
              @scrollReachBottom="handleScroll"
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
                <div class="job-item-box">
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
                                :data-clipboard-text="slotProp.item.err_msg"
                              >
                                <DocumentCopy />
                              </el-icon>
                            </div>
                            <div class="mj-list-item-prompt">
                              <span>{{ slotProp.item.err_msg }}</span>
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
import { ref, reactive, onMounted, onUnmounted, watch } from "vue";
import {
  Plus,
  Delete,
  InfoFilled,
  ChromeFilled,
  DocumentCopy,
  Download,
  WarnTriangleFilled,
  CircleCloseFilled
} from "@element-plus/icons-vue";
import { httpGet, httpPost, httpDownload } from "@/utils/http";
import { ElMessage, ElMessageBox } from "element-plus";
import { checkSession, getSystemInfo } from "@/store/cache";
import Clipboard from "clipboard";

import {
  closeLoading,
  showLoading,
  showMessageError,
  showMessageOK
} from "@/utils/dialog";
import { replaceImg } from "@/utils/libs";

const models = ref([
  {
    text: "可灵 1.6",
    value: "kling-v1-6"
  },
  {
    text: "可灵 1.5",
    value: "kling-v1-5"
  },
  {
    text: "可灵 1.0",
    value: "kling-v1"
  }
]);
// 参数设置
const params = reactive({
  task_type: "text2video",
  model: models.value[0].value,
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
let pollTimer = null;
const generating = ref(false);
const isGenerating = ref(false);
const powerCost = ref(10);
const availablePower = ref(100);
const taskFilter = ref("all");
const runningTasks = ref([]);
const finishedTasks = ref([]);
const pageSize = ref(4);
const currentPage = ref(1);
const previewVisible = ref(false);
const currentVideo = ref("");
const isOver = ref(false);
const pullTask = ref(true);
const waterfallKey = ref(Date.now());
const showCameraControl = ref(false);
const keLingPowers = ref({});

// 动态更新模型消耗的算力
const updateModelPower = () => {
  showCameraControl.value =
    params.model === "kling-v1-5" && params.mode === "pro";
  powerCost.value =
    keLingPowers.value[`${params.model}_${params.mode}_${params.duration}`] ||
    {};
};

// tab切换
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

//移除图片
const remove = (type, img) => {
  if (type === "start") {
    params.image = "";
  } else if (type === "end") {
    params.image_tail = "";
  }
};

//图片交换方法
const switchReverse = () => {
  [params.image, params.image_tail] = [params.image_tail, params.image];
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
  if (isGenerating.value) return;
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
  //增加防抖
  if (generating.value) return;
  if (!params.prompt?.trim()) {
    return ElMessage.error("请输入视频描述");
  }
  // 提示词长度不能超过 500
  if (params.prompt.length > 500) {
    return ElMessage.error("视频描述不能超过 500 个字符");
  }
  // if (params.task_type === "image2video" && !params.image) {
  //   return ElMessage.error("请上传起始帧图片");
  // }
  generating.value = true;
  // 处理图片链接
  if (params.image) {
    params.image = replaceImg(params.image);
  }
  if (params.image_tail) {
    params.image_tail = replaceImg(params.image_tail);
  }
  try {
    await httpPost("/api/video/keling/create", params);
    showMessageOK("任务创建成功");
    // 新增重置
    currentPage.value = 1;
    runningTasks.value = [];
    finishedTasks.value = [];
    pullTask.value = true; //新增开始轮询-获取正在进行中最新数据
    isOver.value = false;
  } catch (e) {
    showMessageError("创建失败: " + e.message);
  } finally {
    generating.value = false;
  }
};
const loading = ref(false);
const isLogin = ref(false);
const runningPage = ref(1);
const taskfunction = async (running) => {
  const res = await httpGet("/api/video/list", {
    page: running ? runningPage.value : currentPage.value,
    page_size: running ? 100000 : pageSize.value,
    type: "keling",
    task_type: taskFilter.value === "all" ? "" : taskFilter.value
  });
  return res;
};
//滚动 瀑布流
const handleScroll = async () => {
  if (!isLogin.value || loading.value || isOver.value) return;
  try {
    loading.value = true;
    const res = await taskfunction(false);
    const data = res.data || {};

    // 去重合并新数据
    const newItems = data.items
      .filter(
        (task) =>
          task.progress >= 100 &&
          !finishedTasks.value.some((t) => t.id === task.id)
      )
      .map((item) => ({
        ...item,
        height: 300 * (Math.random() * 0.3 + 0.6) //生成300~420px随机高度
      }));
    finishedTasks.value = [...finishedTasks.value, ...newItems];

    if (
      currentPage.value >= data.total_page ||
      data.total <= finishedTasks.value.length
    ) {
      isOver.value = true;
    } else {
      currentPage.value++;
    }
  } catch (error) {
    console.error("Error:", error);
  } finally {
    loading.value = false;
    console.log("finishedTasks.value", finishedTasks.value);
  }
};

//正在生成任务中 轮询
const fetchTasks = async () => {
  if (!isLogin.value) return;
  try {
    const res = await taskfunction(true);

    const data = res.data || {};
    // 拿到任务中的数据
    runningTasks.value = data.items.filter((task) => task.progress < 100);
    if (runningPage.value >= data.total_page) {
      runningPage.value = 1;
    } else {
      runningPage.value++;
    }
  } catch (e) {
    showMessageError("获取任务列表失败: " + e.message);
  } finally {
  }
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
    await fetchTasks();
  } catch (e) {
    if (e !== "cancel") {
      showMessageError("删除失败: " + e.message);
    }
  }
};

const clipboard = ref(null);
// 生命周期钩子
onMounted(async () => {
  checkSession()
    .then(async (u) => {
      isLogin.value = true;
      availablePower.value = u.power;
      await fetchTasks();
      await handleScroll();
    })
    .catch(() => {});
  //复制
  clipboard.value = new Clipboard(".copy-prompt-kl");
  clipboard.value.on("success", () => {
    ElMessage.success("复制成功！");
  });

  clipboard.value.on("error", () => {
    ElMessage.error("复制失败！");
  });
  // 获取系统配置
  getSystemInfo().then((res) => {
    keLingPowers.value = res.data.keling_powers;
    updateModelPower();
  });
});
onUnmounted(() => {
  clearInterval(pollTimer);
  clipboard.value.destroy();
});

//监听 pullTask true 的时候定时器触发，否则清除定时器
watch(
  () => pullTask.value,
  (newVal) => {
    clearInterval(pollTimer);
    if (newVal) {
      pollTimer = setInterval(() => {
        fetchTasks();
      }, 5000);
    }
  },
  { immediate: true }
);

watch(
  () => runningTasks.value,
  async (newVal, oldVal) => {
    //正在排队的新旧数据不一致，说明有新的任务加入，或减少，重新获取历史记录
    if (
      newVal.length > 0 &&
      JSON.stringify(newVal) !== JSON.stringify(oldVal)
    ) {
      pullTask.value = true;
      // 立即获取最新历史记录数据
      currentPage.value = 1;
      finishedTasks.value = [];
      await handleScroll();
    }
    //当正在排队的长度0，停止轮询，历史记录重新请求一次
    if (newVal.length === 0) {
      pullTask.value = false;
      currentPage.value = 1;
      runningTasks.value = [];
      finishedTasks.value = [];
      isOver.value = false;
      // 立即获取最新历史记录数据
      await handleScroll();
    }
  },
  { deep: true }
);
</script>

<style lang="stylus" scoped>
@import "@/assets/css/image-keling.styl"
@import "@/assets/css/custom-scroll.styl"
  .video-img-box{
    position: relative;
    .removeimg{
      position: absolute;
      cursor: pointer;
      font-size: 20px;
      color: #545454;
      right: -6px;
      top: 13px;
      z-index: 1;

    }
    &:hover.img-mask{
      opacity: 1;
    }
  }
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
