<template>
  <el-dialog v-model="show" :fullscreen="true" @close="close" style="--el-dialog-border-radius: 0px">
    <template #header>
      <div class="header">
        <h3 style="color: var(--text-theme-color)">绘画任务详情</h3>
      </div>
    </template>
    <el-row :gutter="20">
      <el-col :span="16">
        <div class="img-container">
          <el-image :src="item['img_url']" fit="contain">
            <template #placeholder>
              <div class="image-slot">正在加载图片</div>
            </template>

            <template #error>
              <div class="image-slot">
                <el-icon>
                  <i class="iconfont icon-image"></i>
                </el-icon>
              </div>
            </template>
          </el-image>
        </div>
      </el-col>
      <el-col :span="8">
        <div class="task-info">
          <div class="info-line">
            <el-divider> 正向提示词 </el-divider>
            <div class="prompt">
              <span>{{ item.prompt }}</span>
              <el-icon class="copy-prompt-wall" :data-clipboard-text="item.prompt">
                <i class="iconfont icon-copy"></i>
              </el-icon>
            </div>
          </div>

          <div class="info-line">
            <el-divider> 反向提示词 </el-divider>
            <div class="prompt">
              <span>{{ item.params.negative_prompt }}</span>
              <el-icon class="copy-prompt-wall" :data-clipboard-text="item.params.negative_prompt">
                <i class="iconfont icon-copy"></i>
              </el-icon>
            </div>
          </div>

          <div class="info-line">
            <div class="wrapper">
              <label>采样方法：</label>
              <div class="item-value">{{ item.params.sampler }}</div>
            </div>
          </div>

          <div class="info-line">
            <div class="wrapper">
              <label>图片尺寸：</label>
              <div class="item-value">{{ item.params.width }} x {{ item.params.height }}</div>
            </div>
          </div>

          <div class="info-line">
            <div class="wrapper">
              <label>迭代步数：</label>
              <div class="item-value">{{ item.params.steps }}</div>
            </div>
          </div>

          <div class="info-line">
            <div class="wrapper">
              <label>引导系数：</label>
              <div class="item-value">{{ item.params.cfg_scale }}</div>
            </div>
          </div>

          <div class="info-line">
            <div class="wrapper">
              <label>随机因子：</label>
              <div class="item-value">{{ item.params.seed }}</div>
            </div>
          </div>

          <div v-if="item.params.hd_fix">
            <el-divider> 高清修复 </el-divider>
            <div class="info-line">
              <div class="wrapper">
                <label>重绘幅度：</label>
                <div class="item-value">{{ item.params.hd_redraw_rate }}</div>
              </div>
            </div>

            <div class="info-line">
              <div class="wrapper">
                <label>放大算法：</label>
                <div class="item-value">{{ item.params.hd_scale_alg }}</div>
              </div>
            </div>

            <div class="info-line">
              <div class="wrapper">
                <label>放大倍数：</label>
                <div class="item-value">{{ item.params.hd_scale }}</div>
              </div>
            </div>

            <div class="info-line">
              <div class="wrapper">
                <label>迭代步数：</label>
                <div class="item-value">{{ item.params.hd_steps }}</div>
              </div>
            </div>
          </div>

          <div class="copy-params">
            <el-button type="primary" round @click="drawSame(item)">画一张同款的</el-button>
          </div>
        </div>
      </el-col>
    </el-row>
  </el-dialog>
</template>

<script setup>
import { ref, watch, onMounted } from "vue";
import Clipboard from "clipboard";
import { showMessageOK, showMessageError } from "@/utils/dialog";

const props = defineProps({
  modelValue: Boolean,
  data: Object,
});

const item = ref(props.data);
const show = ref(props.modelValue);
const emit = defineEmits(["drawSame", "close"]);

const clipboard = ref(null);
onMounted(() => {
  clipboard.value = new Clipboard(".copy-prompt-wall");
  clipboard.value.on("success", () => {
    showMessageOK("复制成功！");
  });

  clipboard.value.on("error", () => {
    showMessageError("复制失败！");
  });
});

watch(
  () => props.modelValue,
  (newValue) => {
    show.value = newValue;
  }
);

watch(
  () => props.data,
  (newValue) => {
    item.value = newValue;
  }
);

const drawSame = (item) => {
  emit("drawSame", item);
};

const close = () => {
  emit("close");
};
</script>

<style lang="stylus" scoped></style>
