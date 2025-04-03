<template>
  <div class="page-images-wall">
    <div class="inner custom-scroll">
      <div class="header">
        <h2 class="text-xl pt-4 pb-4">AI 绘画作品墙</h2>
        <div class="settings pr-14">
          <el-radio-group v-model="imgType" @change="changeImgType">
            <el-radio value="mj" size="large">MidJourney</el-radio>
            <el-radio value="sd" size="large">Stable Diffusion</el-radio>
            <el-radio value="dall" size="large">DALL-E</el-radio>
          </el-radio-group>
        </div>
      </div>
      <div
        class="waterfall"
        :style="{ height: listBoxHeight + 'px' }"
        id="waterfall-box"
      >
        <Waterfall
          v-if="imgType === 'mj'"
          id="waterfall-mj"
          :list="data['mj']"
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
                  class="cursor-pointer transition-all duration-300 ease-linear group-hover:scale-105"
                  @click="previewImg(item)"
                />
              </div>
              <div class="px-4 pt-2 pb-4 border-t border-t-gray-800">
                <div
                  class="pt-3 flex justify-center items-center border-t border-t-gray-600 border-opacity-50"
                >
                  <div class="opt">
                    <el-tooltip
                      class="box-item"
                      content="复制提示词"
                      placement="top"
                    >
                      <el-button
                        type="info"
                        circle
                        class="copy-prompt-wall"
                        :data-clipboard-text="item.prompt"
                      >
                        <i class="iconfont icon-file"></i>
                      </el-button>
                    </el-tooltip>

                    <el-tooltip
                      class="box-item"
                      content="画同款"
                      placement="top"
                    >
                      <el-button
                        type="primary"
                        circle
                        @click="drawSameMj(item)"
                      >
                        <i class="iconfont icon-palette"></i>
                      </el-button>
                    </el-tooltip>
                  </div>
                </div>
              </div>
            </div>
          </template>
        </Waterfall>

        <Waterfall
          v-if="imgType === 'sd'"
          id="waterfall-sd"
          :list="data['sd']"
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
                  class="cursor-pointer transition-all duration-300 ease-linear group-hover:scale-105"
                  @click="showTask(item)"
                />
              </div>
              <div class="px-4 pt-2 pb-4 border-t border-t-gray-800">
                <div
                  class="pt-3 flex justify-center items-center border-t border-t-gray-600 border-opacity-50"
                >
                  <div class="opt">
                    <el-tooltip
                      class="box-item"
                      content="复制提示词"
                      placement="top"
                    >
                      <el-button
                        type="info"
                        circle
                        class="copy-prompt-wall"
                        :data-clipboard-text="item.prompt"
                      >
                        <i class="iconfont icon-file"></i>
                      </el-button>
                    </el-tooltip>

                    <el-tooltip
                      class="box-item"
                      content="画同款"
                      placement="top"
                    >
                      <el-button
                        type="primary"
                        circle
                        @click="drawSameSd(item)"
                      >
                        <i class="iconfont icon-palette"></i>
                      </el-button>
                    </el-tooltip>
                  </div>
                </div>
              </div>
            </div>
          </template>
        </Waterfall>

        <Waterfall
          v-if="imgType === 'dall'"
          id="waterfall-dall"
          :list="data['dall']"
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
                  class="cursor-pointer transition-all duration-300 ease-linear group-hover:scale-105"
                  @click="previewImg(item)"
                />
              </div>
              <div class="px-4 pt-2 pb-4 border-t border-t-gray-800">
                <div
                  class="pt-3 flex justify-center items-center border-t border-t-gray-600 border-opacity-50"
                >
                  <div class="opt">
                    <el-tooltip
                      class="box-item"
                      content="复制提示词"
                      placement="top"
                    >
                      <el-button
                        type="info"
                        circle
                        class="copy-prompt-wall"
                        :data-clipboard-text="item.prompt"
                      >
                        <i class="iconfont icon-file"></i>
                      </el-button>
                    </el-tooltip>
                  </div>
                </div>
              </div>
            </div>
          </template>
        </Waterfall>

        <div class="flex flex-col items-center justify-center py-10">
          <img
            :src="waterfallOptions.loadProps.loading"
            class="max-w-[50px] max-h-[50px]"
            v-if="loading"
          />
          <div v-else>
            <button
              class="px-5 py-2 rounded-full bg-purple-700 text-md text-white cursor-pointer hover:bg-purple-800 transition-all duration-300"
              @click="getNext"
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

        <back-top :right="30" :bottom="30" />
      </div>
      <!-- end of waterfall -->
    </div>
    <!-- 任务详情弹框 -->
    <sd-task-view
      v-model="showTaskDialog"
      :data="item"
      @drawSame="drawSameSd"
      @close="showTaskDialog = false"
    />

    <!-- 图片预览 -->
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
import { DocumentCopy, Picture } from "@element-plus/icons-vue";
import { httpGet } from "@/utils/http";
import { ElMessage } from "element-plus";
import Clipboard from "clipboard";
import { useRouter } from "vue-router";
import BackTop from "@/components/BackTop.vue";
import SdTaskView from "@/components/SdTaskView.vue";
import { LazyImg, Waterfall } from "vue-waterfall-plugin-next";
import "vue-waterfall-plugin-next/dist/style.css";
import { useSharedStore } from "@/store/sharedata";

const store = useSharedStore();
const waterfallOptions = store.waterfallOptions;

const data = ref({
  mj: [],
  sd: [],
  dall: [],
});
const loading = ref(true);
const isOver = ref(false);
const imgType = ref("mj"); // 图片类别
const listBoxHeight = window.innerHeight - 124;
const showTaskDialog = ref(false);
const item = ref({});
const previewURL = ref("");

const previewImg = (item) => {
  previewURL.value = item.img_url;
};

const page = ref(0);
const pageSize = ref(15);
// 获取下一页数据
const getNext = () => {
  if (isOver.value) {
    return;
  }

  loading.value = true;
  page.value = page.value + 1;
  let url = "";
  switch (imgType.value) {
    case "mj":
      url = "/api/mj/imgWall";
      break;
    case "sd":
      url = "/api/sd/imgWall";
      break;
    case "dall":
      url = "/api/dall/imgWall";
      break;
  }
  httpGet(`${url}?page=${page.value}&page_size=${pageSize.value}`)
    .then((res) => {
      if (!res.data.items || res.data.items.length === 0) {
        isOver.value = true;
        loading.value = false;
        return;
      }

      // 生成缩略图
      const imageList = res.data.items;
      for (let i = 0; i < imageList.length; i++) {
        imageList[i]["img_thumb"] =
          imageList[i]["img_url"] + "?imageView2/4/w/300/h/0/q/75";
      }
      if (data.value[imgType.value].length === 0) {
        data.value[imgType.value] = imageList;
        return;
      }

      if (imageList.length < pageSize.value) {
        isOver.value = true;
      }
      data.value[imgType.value] = data.value[imgType.value].concat(imageList);
    })
    .catch((e) => {
      ElMessage.error("获取图片失败：" + e.message);
      loading.value = false;
    });
};

getNext();

const clipboard = ref(null);
onMounted(() => {
  clipboard.value = new Clipboard(".copy-prompt-wall");
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

const changeImgType = () => {
  console.log(imgType.value);
  document.getElementById("waterfall-box").scrollTo(0, 0);
  page.value = 0;
  data.value = {
    mj: [],
    sd: [],
    dall: [],
  };
  loading.value = true;
  isOver.value = false;
  nextTick(() => getNext());
};

const showTask = (row) => {
  item.value = row;
  showTaskDialog.value = true;
};

const router = useRouter();
const drawSameSd = (row) => {
  router.push({
    name: "image-sd",
    params: { copyParams: JSON.stringify(row.params) },
  });
};

const drawSameMj = (row) => {
  router.push({ name: "image-mj", params: { prompt: row.prompt } });
};
</script>

<style lang="stylus">
@import '@/assets/css/images-wall.styl';
@import '@/assets/css/custom-scroll.styl';
</style>
