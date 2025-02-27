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
                      :class="item.value === params.aspect_ratio ? 'grid-content active' : 'grid-content'"
                      @click="changeRate(item)"
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
                <el-slider v-model="params.cfg_scale" :min="0" :max="1" :step="0.1" />
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
                <el-select v-model="params.camera_control.type" placeholder="请选择运镜类型">
                  <el-option label="请选择" value="" />
                  <el-option label="简单运镜" value="simple" />
                  <el-option label="下移拉远" value="down_back" />
                  <el-option label="推进上移" value="forward_up" />
                  <el-option label="右旋推进" value="right_turn_forward" />
                  <el-option label="左旋推进" value="left_turn_forward" />
                </el-select>
              </el-form-item>

              <!-- 仅在simple模式下显示详细配置 -->
              <div class="camera-control" v-if="params.camera_control.type === 'simple'">
                <el-form-item label="水平移动">
                  <el-slider v-model="params.camera_control.config.horizontal" :min="-10" :max="10" />
                </el-form-item>
                <el-form-item label="垂直移动">
                  <el-slider v-model="params.camera_control.config.vertical" :min="-10" :max="10" />
                </el-form-item>
                <el-form-item label="左右旋转">
                  <el-slider v-model="params.camera_control.config.pan" :min="-10" :max="10" />
                </el-form-item>
                <el-form-item label="上下旋转">
                  <el-slider v-model="params.camera_control.config.tilt" :min="-10" :max="10" />
                </el-form-item>
                <el-form-item label="横向翻转">
                  <el-slider v-model="params.camera_control.config.roll" :min="-10" :max="10" />
                </el-form-item>
                <el-form-item label="镜头缩放">
                  <el-slider v-model="params.camera_control.config.zoom" :min="-10" :max="10" />
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
          <el-tabs v-model="params.task_type" @tab-change="tabChange" class="title-tabs">
            <el-tab-pane label="文生视频" name="text2video">
              <div class="text">使用文字描述想要生成视频的内容</div>
            </el-tab-pane>
            <el-tab-pane label="图生视频" name="image2video">
              <div class="text">以某张图片为底稿参考来创作视频，生成类似风格或类型视频，支持 PNG /JPG/JPEG 格式图片；</div>
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
              <el-button class="generate-btn" @click="generatePrompt" :loading="isGenerating" size="small" color="#5865f2">
                <i class="iconfont icon-chuangzuo"></i>
                生成专业视频提示词
              </el-button>
            </el-row>
          </div>

          <div v-else class="image2video">
            <div class="image-upload img-inline">
              <div class="upload-box img-uploader">
                <h4>起始帧</h4>
                <el-upload class="uploader img-uploader" :auto-upload="true" :show-file-list="false" :http-request="uploadStartImage" accept=".jpg,.png,.jpeg">
                  <img v-if="params.image" :src="params.image" class="preview" />
                  <el-icon v-else class="upload-icon"><Plus /></el-icon>
                </el-upload>
              </div>
              <div class="upload-box img-uploader">
                <h4>结束帧</h4>
                <el-upload class="uploader" :auto-upload="true" :show-file-list="false" :http-request="uploadEndImage" accept=".jpg,.png,.jpeg">
                  <img v-if="params.image_tail" :src="params.image_tail" class="preview" />
                  <el-icon v-else class="upload-icon"><Plus /></el-icon>
                </el-upload>
              </div>
            </div>
            <div class="param-line pt">
              <div class="flex-row justify-between items-center">
                <div class="flex-row justify-start items-center">
                  <span>提示词：</span>
                  <el-tooltip content="输入你想要的内容，用逗号分割" placement="right">
                    <el-icon>
                      <InfoFilled />
                    </el-icon>
                  </el-tooltip>
                </div>
              </div>
            </div>
            <div class="param-line pt">
              <el-input v-model="params.prompt" type="textarea" :autosize="{ minRows: 4, maxRows: 6 }" placeholder="描述视频画面细节" />
            </div>
          </div>

          <!-- 排除内容 -->
          <div class="param-line pt">
            <div class="flex-row justify-between items-center">
              <div class="flex-row justify-start items-center">
                <span>不希望出现的内容：（可选）</span>
                <el-tooltip content="不想出现在图片上的元素(例如：树，建筑)" placement="right">
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
              >每次生成视频消耗 <el-text type="warning">{{ powerCost }}算力;</el-text> </el-text
            >&nbsp;&nbsp;
            <el-text type="primary"
              >当前可用算力：<el-text type="warning">{{ availablePower }}</el-text></el-text
            >
          </el-row>

          <!-- 生成按钮 -->
          <div class="submit-btn">
            <el-button type="primary" :dark="false" @click="generate" round>立即生成</el-button>
          </div>
        </div>

        <!-- 任务列表区域 -->
        <div class="task-list job-list-box">
          <h2 class="text-xl param-line pt">任务列表</h2>

          <!-- 任务类型筛选 -->

          <el-radio-group v-model="taskFilter" @change="fetchTasks">
            <el-radio-button label="all">全部</el-radio-button>
            <el-radio-button label="text2video">文生视频</el-radio-button>
            <el-radio-button label="image2video">图生视频</el-radio-button>
          </el-radio-group>

          <!-- 运行中的任务 -->
          <div class="running-tasks" v-if="runningTasks.length > 0">
            <h3>运行中</h3>
            <div class="task-grid">
              <div v-for="task in runningTasks" :key="task.id" class="task-card">
                <div class="status">处理中...</div>
                <div class="prompt">{{ task.prompt }}</div>
              </div>
            </div>
          </div>

          <!-- 已完成的任务 -->
          <div class="finished-tasks">
            <h3>已完成</h3>
            <div class="task-grid">
              <div v-for="task in finishedTasks" :key="task.id" class="task-card">
                <video class="preview" :src="task.video_url" @click="previewVideo(task)" controls></video>
                <div class="tools">
                  <el-button @click="downloadVideo(task)">下载</el-button>
                  <el-button @click="deleteTask(task)">删除</el-button>
                </div>
              </div>
            </div>
          </div>

          <!-- 分页控制 -->
          <el-pagination
            v-if="total > pageSize"
            background
            layout="prev, pager, next"
            :total="total"
            :page-size="pageSize"
            @current-change="handlePageChange"
          />
        </div>
      </div>
    </div>

    <!-- 视频预览对话框 -->
    <el-dialog v-model="previewVisible" title="视频预览" width="80%">
      <video v-if="currentVideo" :src="currentVideo" controls style="width: 100%"></video>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from "vue";
import { Plus, InfoFilled } from "@element-plus/icons-vue";
import { httpGet, httpPost, httpDownload } from "@/utils/http";
import { ElMessage, ElMessageBox } from "element-plus";
import { getClientId } from "@/store/cache";
import {closeLoading, showLoading, showMessageError, showMessageOK} from "@/utils/dialog";

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
      zoom: 0,
    },
  },
  image: "",
  image_tail: "",
});
const rates = [
  { css: "square", value: "1:1", text: "1:1", img: "/images/mj/rate_1_1.png" },

  {
    css: "size16-9",
    value: "16:9",
    text: "16:9",
    img: "/images/mj/rate_16_9.png",
  },
  {
    css: "size9-16",
    value: "9:16",
    text: "9:16",
    img: "/images/mj/rate_9_16.png",
  },
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
const pageSize = ref(12);
const currentPage = ref(1);
const previewVisible = ref(false);
const currentVideo = ref("");

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
    closeLoading()
  } catch (e) {
    showMessageError("上传失败: " + e.message);
    closeLoading()
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
  if (!params.prompt) {
    return showMessageError("请输入视频描述");
  }

  generating.value = true;
  // 处理图片链接
  params.image = replaceImg(params.image);
  params.image_tail = replaceImg(params.image_tail);
  try {
    await httpPost("/api/video/keling/create", params);
    showMessageOK("任务创建成功");
    fetchTasks();
  } catch (e) {
    showMessageError("创建失败: " + e.message);
  } finally {
    generating.value = false;
  }
};

const fetchTasks = async () => {
  try {
    const res = await httpGet("/api/video/list", {
      page: currentPage.value,
      page_size: pageSize.value,
      type: "keling",
      task_type: taskFilter.value === "all" ? "" : taskFilter.value,
    });
    runningTasks.value = res.data.items.filter((task) => task.progress < 100);
    finishedTasks.value = res.data.items.filter((task) => task.progress === 100);
    total.value = res.data.total;
  } catch (e) {
    showMessageError("获取任务列表失败: " + e.message);
  }
};

const handlePageChange = (page) => {
  currentPage.value = page;
  fetchTasks();
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

// 生命周期钩子
onMounted(() => {
  fetchTasks();
});
</script>

<style lang="stylus" scoped>
@import "@/assets/css/image-keling.styl"
@import "@/assets/css/custom-scroll.styl"
</style>
