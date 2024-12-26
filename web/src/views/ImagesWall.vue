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
      <div class="waterfall" :style="{ height: listBoxHeight + 'px' }" id="waterfall-box">
        <v3-waterfall
          v-if="imgType === 'mj'"
          id="waterfall"
          :list="data['mj']"
          srcKey="img_thumb"
          :gap="12"
          :bottomGap="-5"
          :colWidth="colWidth"
          :distanceToScroll="100"
          :isLoading="loading"
          :isOver="isOver"
          @scrollReachBottom="getNext"
        >
          <template #default="slotProp">
            <div class="list-item">
              <div class="image">
                <el-image
                  :src="slotProp.item['img_thumb']"
                  :zoom-rate="1.2"
                  :preview-src-list="[slotProp.item['img_url']]"
                  :preview-teleported="true"
                  :initial-index="10"
                  loading="lazy"
                >
                  <template #placeholder>
                    <div class="image-slot">正在加载图片</div>
                  </template>

                  <template #error>
                    <div class="image-slot">
                      <el-icon>
                        <Picture />
                      </el-icon>
                    </div>
                  </template>
                </el-image>
              </div>
              <div class="opt">
                <el-tooltip class="box-item" content="复制提示词" placement="top">
                  <el-icon class="copy-prompt-wall" :data-clipboard-text="slotProp.item.prompt">
                    <DocumentCopy />
                  </el-icon>
                </el-tooltip>

                <el-tooltip class="box-item" content="画同款" placement="top">
                  <i class="iconfont icon-palette-pen" @click="drawSameMj(slotProp.item)"></i>
                </el-tooltip>
              </div>
            </div>
          </template>
        </v3-waterfall>

        <v3-waterfall
          v-else-if="imgType === 'dall'"
          id="waterfall"
          :list="data['dall']"
          srcKey="img_thumb"
          :gap="12"
          :bottomGap="-5"
          :colWidth="colWidth"
          :distanceToScroll="100"
          :isLoading="loading"
          :isOver="isOver"
          @scrollReachBottom="getNext"
        >
          <template #default="slotProp">
            <div class="list-item">
              <div class="image">
                <el-image
                  :src="slotProp.item['img_thumb']"
                  :zoom-rate="1.2"
                  :preview-src-list="[slotProp.item['img_url']]"
                  :preview-teleported="true"
                  :initial-index="10"
                  loading="lazy"
                >
                  <template #placeholder>
                    <div class="image-slot">正在加载图片</div>
                  </template>

                  <template #error>
                    <div class="image-slot">
                      <el-icon>
                        <Picture />
                      </el-icon>
                    </div>
                  </template>
                </el-image>
              </div>
              <div class="opt">
                <el-tooltip class="box-item" content="复制提示词" placement="top">
                  <el-icon class="copy-prompt-wall" :data-clipboard-text="slotProp.item.prompt">
                    <DocumentCopy />
                  </el-icon>
                </el-tooltip>
              </div>
            </div>
          </template>
        </v3-waterfall>

        <v3-waterfall
          v-else
          id="waterfall"
          :list="data['sd']"
          srcKey="img_thumb"
          :gap="12"
          :bottomGap="-5"
          :colWidth="colWidth"
          :distanceToScroll="100"
          :isLoading="loading"
          :isOver="isOver"
          @scrollReachBottom="getNext"
        >
          <template #default="slotProp">
            <div class="list-item">
              <div class="image">
                <el-image :src="slotProp.item['img_thumb']" loading="lazy" @click="showTask(slotProp.item)">
                  <template #placeholder>
                    <div class="image-slot">正在加载图片</div>
                  </template>

                  <template #error>
                    <div class="image-slot">
                      <el-icon>
                        <Picture />
                      </el-icon>
                    </div>
                  </template>
                </el-image>
              </div>
            </div>
          </template>
        </v3-waterfall>

        <div class="footer" v-if="isOver">
          <!-- <el-empty
            :image-size="100"
            :image="nodata"
            description="没有更多数据了"
          /> -->
          <span>没有更多数据了</span>
          <i class="iconfont icon-face"></i>
        </div>

        <back-top :right="30" :bottom="30" />
      </div>
      <!-- end of waterfall -->
    </div>
    <!-- 任务详情弹框 -->
    <sd-task-view v-model="showTaskDialog" :data="item" @drawSame="drawSameSd" @close="showTaskDialog = false" />
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
const data = ref({
  mj: [],
  sd: [],
  dall: [],
});
const loading = ref(true);
const isOver = ref(false);
const imgType = ref("mj"); // 图片类别
const listBoxHeight = window.innerHeight - 124;
const colWidth = ref(220);
const showTaskDialog = ref(false);
const item = ref({});

// 计算瀑布流列宽度
const calcColWidth = () => {
  const listBoxWidth = window.innerWidth - 60 - 80;
  const rows = Math.floor(listBoxWidth / colWidth.value);
  colWidth.value = Math.floor((listBoxWidth - (rows - 1) * 12) / rows);
};
calcColWidth();
window.onresize = () => {
  calcColWidth();
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
      loading.value = false;
      if (!res.data.items || res.data.items.length === 0) {
        isOver.value = true;
        return;
      }

      // 生成缩略图
      const imageList = res.data.items;
      for (let i = 0; i < imageList.length; i++) {
        imageList[i]["img_thumb"] = imageList[i]["img_url"] + "?imageView2/4/w/300/h/0/q/75";
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
@import "@/assets/css/images-wall.styl"
@import "@/assets/css/custom-scroll.styl"
</style>
