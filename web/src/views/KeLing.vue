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
                <el-select v-model="params.model" placeholder="请选择模型" @change="updateModelPower">
                  <el-option v-for="item in models" :key="item.value" :label="item.text" :value="item.value" />
                </el-select>
              </el-form-item>
            </div>

            <!-- 视频时长 -->
            <div class="param-line">
              <el-form-item label="视频时长">
                <el-select v-model="params.duration" placeholder="请选择时长" @change="updateModelPower">
                  <el-option label="5秒" value="5" />
                  <el-option label="10秒" value="10" />
                </el-select>
              </el-form-item>
            </div>

            <!-- 生成模式 -->
            <div class="param-line">
              <el-form-item label="生成模式">
                <el-select v-model="params.mode" placeholder="请选择模式" @change="updateModelPower">
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
            <div class="param-line" v-if="showCameraControl">
              <div class="param-line pt">
                <span>运镜控制：</span>
                <el-tooltip content="生成画面的运镜效果，仅 1.5的高级模式可用" placement="right">
                  <el-icon>
                    <InfoFilled />
                  </el-icon>
                </el-tooltip>
              </div>

              <!-- 添加运镜类型选择 -->
              <div class="param-line">
                <el-select v-model="params.camera_control.type" placeholder="请选择运镜类型">
                  <el-option label="请选择" value="" />
                  <el-option label="简单运镜" value="simple" />
                  <el-option label="下移拉远" value="down_back" />
                  <el-option label="推进上移" value="forward_up" />
                  <el-option label="右旋推进" value="right_turn_forward" />
                  <el-option label="左旋推进" value="left_turn_forward" />
                </el-select>
              </div>

              <!-- 仅在simple模式下显示详细配置 -->
              <div class="camera-control mt-2" v-if="params.camera_control.type === 'simple'">
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
              maxlength="500"
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
              <div class="upload-box img-uploader video-img-box">
                <el-icon v-if="params.image" @click="removeImage('start')" class="removeimg"><CircleCloseFilled /></el-icon>

                <h4>起始帧</h4>
                <el-upload class="uploader img-uploader" :auto-upload="true" :show-file-list="false" :http-request="uploadStartImage" accept=".jpg,.png,.jpeg">
                  <img v-if="params.image" :src="params.image" class="preview" />

                  <el-icon v-else class="upload-icon"><Plus /></el-icon>
                </el-upload>
              </div>
              <div class="btn-swap" v-if="params.image && params.image_tail">
                <i class="iconfont icon-exchange" @click="switchReverse"></i>
              </div>
              <div class="upload-box img-uploader video-img-box">
                <el-icon v-if="params.image_tail" @click="removeImage('end')" class="removeimg"><CircleCloseFilled /></el-icon>
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
        <div class="video-list">
          <h2 class="text-xl p-3">你的作品</h2>

          <div class="list-box" v-if="!noData">
            <div v-for="item in list" :key="item.id">
              <div class="item">
                <div class="left">
                  <div class="container">
                    <div v-if="item.progress === 100">
                      <video class="video" :src="item.video_url" preload="auto" loop="loop" muted="muted">您的浏览器不支持视频播放</video>
                      <button class="play flex justify-center items-center" @click="previewVideo(item)">
                        <img src="/images/play.svg" alt="" />
                      </button>
                    </div>
                    <el-image v-else-if="item.progress === 101" :src="item.cover_url" fit="cover" />
                    <generating message="正在生成视频" v-else />
                  </div>
                </div>

                <div class="center">
                  <div class="pb-2">
                    <el-tag class="mr-1">{{ item.raw_data.task_type }}</el-tag>
                    <el-tag class="mr-1">{{ item.raw_data.model }}</el-tag>
                    <el-tag class="mr-1">{{ item.raw_data.duration }}秒</el-tag>
                    <el-tag class="mr-1">{{ item.raw_data.mode }}</el-tag>
                  </div>
                  <div class="failed" v-if="item.progress === 101">任务执行失败：{{ item.err_msg }}，任务提示词：{{ item.prompt }}</div>
                  <div class="prompt" v-else>
                    {{ substr(item.prompt, 1000) }}
                  </div>
                </div>

                <div class="right" v-if="item.progress === 100">
                  <div class="tools">
                    <el-tooltip content="复制提示词" placement="top">
                      <button class="btn btn-icon copy-prompt" :data-clipboard-text="item.prompt">
                        <i class="iconfont icon-copy"></i>
                      </button>
                    </el-tooltip>

                    <!-- <button class="btn btn-publish">
                      <span class="text">发布</span>
                      <black-switch v-model:value="item.publish" @change="publishJob(item)" size="small" />
                    </button> -->

                    <el-tooltip content="下载视频" placement="top">
                      <button class="btn btn-icon" @click="downloadVideo(item)" :disabled="item.downloading">
                        <i class="iconfont icon-download" v-if="!item.downloading"></i>
                        <el-image src="/images/loading.gif" class="downloading" fit="cover" v-else />
                      </button>
                    </el-tooltip>

                    <el-tooltip content="删除" placement="top">
                      <button class="btn btn-icon" @click="removeJob(item)">
                        <i class="iconfont icon-remove"></i>
                      </button>
                    </el-tooltip>
                  </div>
                </div>

                <div class="right-error" v-else>
                  <el-button type="danger" @click="removeJob(item)" circle>
                    <i class="iconfont icon-remove"></i>
                  </el-button>
                </div>
              </div>
            </div>
          </div>

          <el-empty :image-size="100" :image="nodata" description="没有任何作品，赶紧去创作吧！" v-else />

          <div class="pagination">
            <el-pagination
              v-if="total > pageSize"
              background
              style="--el-pagination-button-bg-color: rgba(86, 86, 95, 0.2)"
              layout="total,prev, pager, next"
              :hide-on-single-page="true"
              v-model:current-page="page"
              v-model:page-size="pageSize"
              @current-change="fetchData"
              :total="total"
            />
          </div>
        </div>
      </div>
    </div>

    <!-- 视频预览对话框 -->
    <black-dialog v-model:show="previewVisible" title="视频预览" hide-footer @cancal="previewVisible = false" width="auto">
      <video
        v-if="currentVideo"
        :src="currentVideo"
        controls
        style="max-width: 90vw; max-height: 90vh"
        :autoplay="true"
        loop="loop"
        muted="muted"
        preload="auto"
      >
        您的浏览器不支持视频播放
      </video>
    </black-dialog>
  </div>
</template>

<script setup>
import failed from "@/assets/img/failed.png";
import TaskList from "@/components/TaskList.vue";
import { ref, reactive, onMounted, onUnmounted, watch } from "vue";
import { Plus, Delete, InfoFilled, ChromeFilled, DocumentCopy, Download, WarnTriangleFilled, CircleCloseFilled } from "@element-plus/icons-vue";
import { httpGet, httpPost, httpDownload } from "@/utils/http";
import { ElMessage, ElMessageBox } from "element-plus";
import { checkSession, getSystemInfo } from "@/store/cache";
import Clipboard from "clipboard";
import BlackDialog from "@/components/ui/BlackDialog.vue";
import BlackSwitch from "@/components/ui/BlackSwitch.vue";
import { closeLoading, showLoading, showMessageError, showMessageOK } from "@/utils/dialog";
import { replaceImg, substr } from "@/utils/libs";
import Generating from "@/components/ui/Generating.vue";
import nodata from "@/assets/img/no-data.png";

const models = ref([
  {
    text: "可灵 1.6",
    value: "kling-v1-6",
  },
  {
    text: "可灵 1.5",
    value: "kling-v1-5",
  },
  {
    text: "可灵 1.0",
    value: "kling-v1",
  },
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

const generating = ref(false);
const isGenerating = ref(false);
const powerCost = ref(10);
const availablePower = ref(100);
const taskFilter = ref("all");
const loading = ref(false);
const list = ref([]);
const noData = ref(true);
const page = ref(1);
const pageSize = ref(10);
const total = ref(0);
const taskPulling = ref(true);
const pullHandler = ref(null);
const previewVisible = ref(false);
const currentVideo = ref("");
const showCameraControl = ref(false);
const keLingPowers = ref({});
const isLogin = ref(false);
// 动态更新模型消耗的算力
const updateModelPower = () => {
  showCameraControl.value = params.model === "kling-v1-5" && params.mode === "pro";
  powerCost.value = keLingPowers.value[`${params.model}_${params.mode}_${params.duration}`] || {};
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
const removeImage = (type) => {
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
  if (params.task_type === "image2video" && !params.image) {
    return ElMessage.error("请上传起始帧图片");
  }
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
    page.value = 1;
    list.value.unshift({
      progress: 0,
      prompt: params.prompt,
      raw_data: {
        task_type: params.task_type,
        model: params.model,
        duration: params.duration,
        mode: params.mode,
      },
    });
    taskPulling.value = true;
  } catch (e) {
    showMessageError("创建失败: " + e.message);
  } finally {
    generating.value = false;
  }
};

const fetchData = (_page) => {
  if (_page) {
    page.value = _page;
  }

  httpGet("/api/video/list", {
    page: page.value,
    page_size: pageSize.value,
    type: "keling",
    task_type: taskFilter.value === "all" ? "" : taskFilter.value,
  })
    .then((res) => {
      total.value = res.data.total;
      let needPull = false;
      const items = [];
      for (let v of res.data.items) {
        if (v.progress === 0 || v.progress === 102) {
          needPull = true;
        }
        items.push({
          ...v,
          downloading: false,
        });
      }
      loading.value = false;
      taskPulling.value = needPull;
      if (JSON.stringify(list.value) !== JSON.stringify(items)) {
        list.value = items;
      }
      noData.value = list.value.length === 0;
    })
    .catch(() => {
      loading.value = false;
      noData.value = true;
    });
};

const previewVideo = (task) => {
  currentVideo.value = task.video_url;
  previewVisible.value = true;
};

const downloadVideo = async (task) => {
  try {
    const res = await httpDownload(`/api/download?url=${replaceImg(task.video_url)}`);
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

// 删除任务
const removeJob = (item) => {
  ElMessageBox.confirm("此操作将会删除任务相关文件，继续操作码?", "删除提示", {
    confirmButtonText: "确认",
    cancelButtonText: "取消",
    type: "warning",
  })
    .then(() => {
      httpGet("/api/video/remove", { id: item.id })
        .then(() => {
          ElMessage.success("任务删除成功");
          fetchData(page.value);
        })
        .catch((e) => {
          ElMessage.error("任务删除失败：" + e.message);
        });
    })
    .catch(() => {});
};

const clipboard = ref(null);
// 生命周期钩子
onMounted(() => {
  checkSession()
    .then((u) => {
      isLogin.value = true;
      availablePower.value = u.power;
      fetchData(1);
      // 设置轮询
      pullHandler.value = setInterval(() => {
        if (taskPulling.value) {
          fetchData(page.value);
        }
      }, 5000);
    })
    .catch((e) => {
      console.log(e);
    });

  clipboard.value = new Clipboard(".copy-prompt");
  clipboard.value.on("success", () => {
    ElMessage.success("复制成功！");
  });
  clipboard.value.on("error", () => {
    ElMessage.error("复制失败！");
  });

  getSystemInfo().then((res) => {
    keLingPowers.value = res.data.keling_powers;
    updateModelPower();
  });
});

onUnmounted(() => {
  clipboard.value.destroy();
  if (pullHandler.value) {
    clearInterval(pullHandler.value);
  }
});
</script>

<style lang="stylus" scoped>
@import "@/assets/css/keling.styl"
</style>
