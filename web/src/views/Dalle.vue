<template>
  <div>
    <div class="page-dall">
      <div class="inner custom-scroll">
        <div class="sd-box">
          <h2>DALL-E 创作中心</h2>

          <div class="sd-params">
            <el-form :model="params" label-width="80px" label-position="left">
              <div class="param-line pt-1">
                <el-form-item label="生图模型">
                  <template #default>
                    <div class="form-item-inner">
                      <el-select
                        v-model="selectedModel"
                        style="width: 150px"
                        placeholder="请选择模型"
                        @change="changeModel"
                      >
                        <el-option
                          v-for="v in models"
                          :label="v.name"
                          :value="v"
                          :key="v.value"
                        />
                      </el-select>
                    </div>
                  </template>
                </el-form-item>
              </div>

              <div class="param-line">
                <el-form-item label="图片质量">
                  <template #default>
                    <div class="form-item-inner">
                      <el-select v-model="params.quality" style="width: 150px">
                        <el-option
                          v-for="v in qualities"
                          :label="v.name"
                          :value="v.value"
                          :key="v.value"
                        />
                      </el-select>
                    </div>
                  </template>
                </el-form-item>
              </div>

              <div class="param-line">
                <el-form-item label="图片尺寸">
                  <template #default>
                    <div class="form-item-inner">
                      <el-select v-model="params.size" style="width: 150px">
                        <el-option
                          v-for="v in sizes"
                          :label="v"
                          :value="v"
                          :key="v"
                        />
                      </el-select>
                    </div>
                  </template>
                </el-form-item>
              </div>

              <div class="param-line">
                <el-form-item label="图片样式">
                  <template #default>
                    <div class="form-item-inner">
                      <el-select v-model="params.style" style="width: 150px">
                        <el-option
                          v-for="v in styles"
                          :label="v.name"
                          :value="v.value"
                          :key="v.value"
                        />
                      </el-select>
                      <el-tooltip
                        content="生动使模型倾向于生成超真实和戏剧性的图像"
                        raw-content
                        placement="right"
                      >
                        <el-icon class="info-icon">
                          <InfoFilled />
                        </el-icon>
                      </el-tooltip>
                    </div>
                  </template>
                </el-form-item>
              </div>

              <div class="param-line">
                <el-input
                  v-model="params.prompt"
                  :autosize="{ minRows: 4, maxRows: 6 }"
                  type="textarea"
                  ref="promptRef"
                  maxlength="2000"
                  placeholder="请在此输入绘画提示词，您也可以点击下面的提示词助手生成绘画提示词"
                  v-loading="isGenerating"
                />
              </div>

              <el-row class="text-info">
                <el-button
                  class="generate-btn"
                  size="small"
                  @click="generatePrompt"
                  color="#5865f2"
                  :disabled="isGenerating"
                >
                  <i
                    class="iconfont icon-chuangzuo"
                    style="margin-right: 5px"
                  ></i>
                  <span>生成专业绘画指令</span>
                </el-button>
              </el-row>

              <div class="text-info">
                <el-row :gutter="10">
                  <el-text type="primary"
                    >每次绘图消耗
                    <el-text type="warning"
                      >{{ dallPower }}算力，</el-text
                    ></el-text
                  >
                  <el-text type="primary"
                    >当前可用
                    <el-text type="warning"> {{ power }}算力</el-text>
                  </el-text>
                </el-row>
              </div>
            </el-form>
          </div>
          <div class="submit-btn">
            <el-button type="primary" :dark="false" round @click="generate">
              立即生成
            </el-button>
          </div>
        </div>
        <div class="task-list-box pl-6 pr-6 pb-4 pt-4 h-dvh">
          <div class="task-list-inner">
            <div class="job-list-box">
              <h2 class="text-xl">任务列表</h2>
              <task-list :list="runningJobs" />
              <template v-if="finishedJobs.length > 0">
                <h2 class="text-xl">创作记录</h2>
                <div class="finish-job-list mt-3">
                  <div v-if="finishedJobs.length > 0">
                    <Waterfall
                      :list="finishedJobs"
                      :row-key="waterfallOptions.rowKey"
                      :gutter="waterfallOptions.gutter"
                      :has-around-gutter="waterfallOptions.hasAroundGutter"
                      :width="waterfallOptions.width"
                      :breakpoints="waterfallOptions.breakpoints"
                      :img-selector="waterfallOptions.imgSelector"
                      :background-color="waterfallOptions.backgroundColor"
                      :animation-effect="waterfallOptions.animationEffect"
                      :animation-duration="waterfallOptions.animationDuration"
                      :animation-delay="waterfallOptions.animationDelay"
                      :animation-cancel="waterfallOptions.animationCancel"
                      :lazyload="waterfallOptions.lazyload"
                      :load-props="waterfallOptions.loadProps"
                      :cross-origin="waterfallOptions.crossOrigin"
                      :align="waterfallOptions.align"
                      :is-loading="loading"
                      :is-over="isOver"
                      @afterRender="loading = false"
                    >
                      <template #default="{ item, url }">
                        <div
                          class="bg-gray-900 rounded-lg shadow-md overflow-hidden transition-all duration-300 ease-linear hover:shadow-md hover:shadow-purple-800 group"
                        >
                          <div class="overflow-hidden rounded-lg">
                            <LazyImg
                              :url="url"
                              v-if="item.progress === 100"
                              class="cursor-pointer transition-all duration-300 ease-linear group-hover:scale-105"
                              @click="previewImg(item)"
                            />
                            <el-image v-else-if="item.progress === 101">
                              <template #error>
                                <div class="image-slot">
                                  <div class="err-msg-container">
                                    <div class="title">任务失败</div>
                                    <div class="opt">
                                      <el-popover
                                        title="错误详情"
                                        trigger="click"
                                        :width="250"
                                        :content="item['err_msg']"
                                        placement="top"
                                      >
                                        <template #reference>
                                          <el-button type="info"
                                            >详情</el-button
                                          >
                                        </template>
                                      </el-popover>
                                      <el-button
                                        type="danger"
                                        @click="removeImage(item)"
                                        >删除</el-button
                                      >
                                    </div>
                                  </div>
                                </div>
                              </template>
                            </el-image>
                          </div>
                          <div
                            class="px-4 pt-2 pb-4 border-t border-t-gray-800"
                            v-if="item.progress === 100"
                          >
                            <div
                              class="pt-3 flex justify-center items-center border-t border-t-gray-600 border-opacity-50"
                            >
                              <div class="flex">
                                <el-tooltip
                                  content="取消分享"
                                  placement="top"
                                  v-if="item.publish"
                                >
                                  <el-button
                                    type="warning"
                                    @click="publishImage(item, false)"
                                    circle
                                  >
                                    <i class="iconfont icon-cancel-share"></i>
                                  </el-button>
                                </el-tooltip>
                                <el-tooltip
                                  content="分享"
                                  placement="top"
                                  v-else
                                >
                                  <el-button
                                    type="success"
                                    @click="publishImage(item, true)"
                                    circle
                                  >
                                    <i class="iconfont icon-share-bold"></i>
                                  </el-button>
                                </el-tooltip>

                                <el-tooltip
                                  content="复制提示词"
                                  placement="top"
                                >
                                  <el-button
                                    type="info"
                                    circle
                                    class="copy-prompt"
                                    :data-clipboard-text="item.prompt"
                                  >
                                    <i class="iconfont icon-file"></i>
                                  </el-button>
                                </el-tooltip>
                                <el-tooltip content="删除" placement="top">
                                  <el-button
                                    type="danger"
                                    :icon="Delete"
                                    @click="removeImage(item)"
                                    circle
                                  />
                                </el-tooltip>
                              </div>
                            </div>
                          </div>
                        </div>
                      </template>
                    </Waterfall>

                    <div class="flex justify-center py-10">
                      <img
                        :src="waterfallOptions.loadProps.loading"
                        class="max-w-[50px] max-h-[50px]"
                        v-if="loading"
                      />
                      <div v-else>
                        <button
                          class="px-5 py-2 rounded-full bg-purple-700 text-md text-white cursor-pointer hover:bg-purple-800 transition-all duration-300"
                          @click="fetchFinishJobs"
                          v-if="!isOver"
                        >
                          加载更多
                        </button>
                        <div class="no-more-data" v-else>
                          <span class="text-gray-500 mr-2">没有更多数据了</span>
                          <i class="iconfont icon-face"></i>
                        </div>
                      </div>
                    </div>
                  </div>
                  <el-empty
                    :image-size="100"
                    :image="nodata"
                    description="暂无记录"
                    v-else
                  />
                </div>
              </template>
              <!-- end finish job list-->
            </div>
          </div>
          <back-top :right="30" :bottom="30" />
        </div>
        <!-- end task list box -->
      </div>
    </div>

    <el-image-viewer
      @close="
        () => {
          previewURL = '';
        }
      "
      v-if="previewURL !== ''"
      :url-list="[previewURL]"
    />
  </div>
</template>

<script setup>
import nodata from "@/assets/img/no-data.png";

import { nextTick, onMounted, onUnmounted, ref } from "vue";
import { Delete, InfoFilled, Picture } from "@element-plus/icons-vue";
import { httpGet, httpPost } from "@/utils/http";
import { ElMessage, ElMessageBox } from "element-plus";
import Clipboard from "clipboard";
import { checkSession, getSystemInfo } from "@/store/cache";
import { useSharedStore } from "@/store/sharedata";
import TaskList from "@/components/TaskList.vue";
import BackTop from "@/components/BackTop.vue";
import { showMessageError, showMessageOK } from "@/utils/dialog";
import { LazyImg, Waterfall } from "vue-waterfall-plugin-next";
import "vue-waterfall-plugin-next/dist/style.css";

const listBoxHeight = ref(0);
// const paramBoxHeight = ref(0)
const isLogin = ref(false);
const loading = ref(true);
const isOver = ref(false);
const previewURL = ref("");
const store = useSharedStore();
const models = ref([]);
const waterfallOptions = store.waterfallOptions;
const resizeElement = function () {
  listBoxHeight.value = window.innerHeight - 58;
};

resizeElement();
window.onresize = () => {
  resizeElement();
};
const qualities = [
  { name: "标准", value: "standard" },
  { name: "高清", value: "hd" },
];
const dalleSizes = ["1024x1024", "1792x1024", "1024x1792"];
const fluxSizes = ["1024x1024", "1152x896", "896x1152", "1280x960", "1024x576"];
const sizes = ref(dalleSizes);
const styles = [
  { name: "生动", value: "vivid" },
  { name: "自然", value: "natural" },
];
const params = ref({
  quality: "standard",
  size: "1024x1024",
  style: "vivid",
  prompt: "",
});

const finishedJobs = ref([]);
const runningJobs = ref([]);
const allowPulling = ref(true); // 是否允许轮询
const downloadPulling = ref(false); // 下载轮询
const tastPullHandler = ref(null);
const downloadPullHandler = ref(null);
const power = ref(0);
const dallPower = ref(0); // 画一张 SD 图片消耗算力
const clipboard = ref(null);
const userId = ref(0);
const selectedModel = ref(null);

onMounted(() => {
  initData();
  clipboard.value = new Clipboard(".copy-prompt");
  clipboard.value.on("success", () => {
    showMessageOK("复制成功！");
  });

  clipboard.value.on("error", () => {
    showMessageError("复制失败！");
  });

  getSystemInfo()
    .then((res) => {
      dallPower.value = res.data["dall_power"];
    })
    .catch((e) => {
      showMessageError("获取系统配置失败：" + e.message);
    });

  // 获取模型列表
  httpGet("/api/dall/models")
    .then((res) => {
      models.value = res.data;
      selectedModel.value = models.value[0];
      params.value.model_id = selectedModel.value.id;
      changeModel(selectedModel.value);
    })
    .catch((e) => {
      showMessageError("获取模型列表失败：" + e.message);
    });
});

onUnmounted(() => {
  clipboard.value.destroy();
  if (tastPullHandler.value) {
    clearInterval(tastPullHandler.value);
  }
  if (downloadPullHandler.value) {
    clearInterval(downloadPullHandler.value);
  }
});

const initData = () => {
  checkSession()
    .then((user) => {
      power.value = user["power"];
      userId.value = user.id;
      isLogin.value = true;
      page.value = 0;
      fetchRunningJobs();
      fetchFinishJobs();

      // 轮询运行中任务
      tastPullHandler.value = setInterval(() => {
        if (allowPulling.value) {
          fetchRunningJobs();
        }
      }, 5000);

      // 图片下载轮询
      downloadPullHandler.value = setInterval(() => {
        if (downloadPulling.value) {
          page.value = 0;
          fetchFinishJobs();
        }
      }, 5000);
    })
    .catch(() => {});
};

const fetchRunningJobs = () => {
  if (!isLogin.value) {
    return;
  }
  // 获取运行中的任务
  httpGet(`/api/dall/jobs?finish=false`)
    .then((res) => {
      // 如果任务有更新，则更新已完成任务列表
      if (
        res.data.items &&
        res.data.items.length !== runningJobs.value.length
      ) {
        page.value = 0;
        fetchFinishJobs();
      }
      if (res.data.items.length > 0) {
        runningJobs.value = res.data.items;
      } else {
        allowPulling.value = false;
        runningJobs.value = [];
      }
    })
    .catch((e) => {
      ElMessage.error("获取任务失败：" + e.message);
    });
};

const page = ref(1);
const pageSize = ref(15);
// 获取已完成的任务
const fetchFinishJobs = () => {
  if (!isLogin.value) {
    return;
  }

  loading.value = true;
  page.value = page.value + 1;

  httpGet(
    `/api/dall/jobs?finish=true&page=${page.value}&page_size=${pageSize.value}`
  )
    .then((res) => {
      if (res.data.items.length < pageSize.value) {
        isOver.value = true;
        loading.value = false;
      }
      const imageList = res.data.items;
      let needPulling = false;
      for (let i = 0; i < imageList.length; i++) {
        if (imageList[i]["img_url"]) {
          imageList[i]["img_thumb"] =
            imageList[i]["img_url"] + "?imageView2/4/w/300/h/0/q/75";
        } else if (imageList[i].progress === 100) {
          needPulling = true;
          imageList[i]["img_thumb"] = waterfallOptions.loadProps.loading;
        }
      }
      // 如果当前是第一页，则开启图片下载轮询
      if (page.value === 1) {
        downloadPulling.value = needPulling;
      }

      if (page.value === 1) {
        finishedJobs.value = imageList;
      } else {
        finishedJobs.value = finishedJobs.value.concat(imageList);
      }
    })
    .catch((e) => {
      ElMessage.error("获取任务失败：" + e.message);
      loading.value = false;
    });
};

// 创建绘图任务
const promptRef = ref(null);
const generate = () => {
  if (params.value.prompt === "") {
    promptRef.value.focus();
    return ElMessage.error("请输入绘画提示词！");
  }

  if (!isLogin.value) {
    store.setShowLoginDialog(true);
    return;
  }
  httpPost("/api/dall/image", params.value)
    .then(() => {
      ElMessage.success("任务执行成功！");
      power.value -= dallPower.value;
      // 追加任务列表
      runningJobs.value.push({
        prompt: params.value.prompt,
        progress: 0,
      });
      allowPulling.value = true;
    })
    .catch((e) => {
      ElMessage.error("任务执行失败：" + e.message);
    });
};

const removeImage = (item) => {
  ElMessageBox.confirm("此操作将会删除任务和图片，继续操作码?", "删除提示", {
    confirmButtonText: "确认",
    cancelButtonText: "取消",
    type: "warning",
  })
    .then(() => {
      httpGet("/api/dall/remove", { id: item.id })
        .then(() => {
          ElMessage.success("任务删除成功");
          page.value = 0;
          isOver.value = false;
          fetchFinishJobs();
        })
        .catch((e) => {
          ElMessage.error("任务删除失败：" + e.message);
        });
    })
    .catch(() => {});
};

const previewImg = (item) => {
  previewURL.value = item.img_url;
};

// 发布图片到作品墙
const publishImage = (item, action) => {
  let text = "图片发布";
  if (action === false) {
    text = "取消发布";
  }
  httpGet("/api/dall/publish", { id: item.id, action: action })
    .then(() => {
      ElMessage.success(text + "成功");
      item.publish = action;
      page.value = 0;
      isOver.value = false;
    })
    .catch((e) => {
      ElMessage.error(text + "失败：" + e.message);
    });
};

const isGenerating = ref(false);
const generatePrompt = () => {
  if (params.value.prompt === "") {
    return showMessageError("请输入原始提示词");
  }
  isGenerating.value = true;
  httpPost("/api/prompt/image", { prompt: params.value.prompt })
    .then((res) => {
      params.value.prompt = res.data;
      isGenerating.value = false;
    })
    .catch((e) => {
      showMessageError("生成提示词失败：" + e.message);
      isGenerating.value = false;
    });
};

const changeModel = (model) => {
  if (model.value.startsWith("dall")) {
    sizes.value = dalleSizes;
  } else {
    sizes.value = fluxSizes;
  }
  params.value.model_id = selectedModel.value.id;
};
</script>

<style lang="stylus">
@import '@/assets/css/image-dall.styl';
@import '@/assets/css/custom-scroll.styl';
</style>
