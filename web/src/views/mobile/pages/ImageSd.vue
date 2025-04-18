<template>
  <div class="mobile-sd">
    <van-form @submit="generate">
      <van-cell-group inset>
        <div>
          <van-field v-model="params.sampler" is-link readonly label="采样方法" placeholder="选择采样方法" @click="showSamplerPicker = true" />
          <van-popup v-model:show="showSamplerPicker" position="bottom" teleport="#app">
            <van-picker :columns="samplers" @cancel="showSamplerPicker = false" @confirm="samplerConfirm" />
          </van-popup>
        </div>

        <van-field label="图片尺寸">
          <template #input>
            <van-row gutter="20">
              <van-col span="12">
                <el-input v-model="params.width" size="small" placeholder="宽" />
              </van-col>
              <van-col span="12">
                <el-input v-model="params.height" size="small" placeholder="高" />
              </van-col>
            </van-row>
          </template>
        </van-field>

        <van-field v-model.number="params.steps" label="迭代步数" placeholder="">
          <template #right-icon>
            <van-icon name="info-o" @click="showInfo('值越大则代表细节越多，同时也意味着出图速度越慢，一般推荐20-30')" />
          </template>
        </van-field>
        <van-field v-model.number="params.cfg_scale" label="引导系数" placeholder="">
          <template #right-icon>
            <van-icon name="info-o" @click="showInfo('提示词引导系数，图像在多大程度上服从提示词，较低值会产生更有创意的结果')" />
          </template>
        </van-field>
        <van-field v-model.number="params.seed" label="随机因子" placeholder="">
          <template #right-icon>
            <van-icon name="info-o" @click="showInfo('随机数种子，相同的种子会得到相同的结果，设置为 -1 则每次随机生成种子')" />
          </template>
        </van-field>

        <van-field label="高清修复">
          <template #input>
            <van-switch v-model="params.hd_fix" />
          </template>
        </van-field>

        <div v-if="params.hd_fix">
          <div>
            <van-field v-model="params.hd_scale_alg" is-link readonly label="放大算法" placeholder="选择放大算法" @click="showUpscalePicker = true" />
            <van-popup v-model:show="showUpscalePicker" position="bottom" teleport="#app">
              <van-picker :columns="upscaleAlgArr" @cancel="showUpscalePicker = false" @confirm="upscaleConfirm" />
            </van-popup>
          </div>

          <van-field v-model.number="params.hd_scale" label="放大倍数" />
          <van-field v-model.number="params.hd_steps" label="迭代步数" />

          <van-field label="重绘幅度">
            <template #input>
              <van-slider v-model.number="params.hd_redraw_rate" :max="1" :step="0.1" @update:model-value="showToast('当前值：' + params.hd_redraw_rate)" />
            </template>
            <template #right-icon>
              <van-icon name="info-o" @click="showInfo('决定算法对图像内容的影响程度，较大的值将得到越有创意的图像')" />
            </template>
          </van-field>
        </div>

        <van-field
          v-model="params.prompt"
          rows="3"
          autosize
          type="textarea"
          placeholder="请在此输入绘画提示词，系统会自动翻译中文提示词，高手请直接输入英文提示词"
        />

        <van-collapse v-model="activeColspan">
          <van-collapse-item title="反向提示词" name="neg_prompt">
            <van-field v-model="params.neg_prompt" rows="3" autosize type="textarea" placeholder="不想出现在图片上的元素(例如：树，建筑)" />
          </van-collapse-item>
        </van-collapse>

        <div class="text-line pt-6">
          <el-tag>绘图消耗{{ sdPower }}算力，当前算力：{{ power }}</el-tag>
        </div>

        <div class="text-line">
          <van-button round block type="primary" native-type="submit"> 立即生成 </van-button>
        </div>
      </van-cell-group>
    </van-form>

    <h3>任务列表</h3>
    <div class="running-job-list pt-3 pb-3">
      <van-empty
        v-if="runningJobs.length === 0"
        image="https://fastly.jsdelivr.net/npm/@vant/assets/custom-empty-image.png"
        image-size="80"
        description="暂无记录"
      />
      <van-grid :gutter="10" :column-num="3" v-else>
        <van-grid-item v-for="item in runningJobs" :key="item.id">
          <div v-if="item.progress > 0">
            <van-image src="/images/img-holder.png"></van-image>
            <div class="progress">
              <van-circle v-model:current-rate="item.progress" :rate="item.progress" :speed="100" :text="item.progress + '%'" :stroke-width="60" size="90px" />
            </div>
          </div>

          <div v-else class="task-in-queue">
            <span class="icon"><i class="iconfont icon-quick-start"></i></span>
            <span class="text">排队中</span>
          </div>
        </van-grid-item>
      </van-grid>
    </div>

    <h3>创作记录</h3>
    <div class="finish-job-list">
      <van-empty
        v-if="finishedJobs.length === 0"
        image="https://fastly.jsdelivr.net/npm/@vant/assets/custom-empty-image.png"
        image-size="80"
        description="暂无记录"
      />

      <van-list
        v-else
        v-model:error="error"
        v-model:loading="loading"
        :finished="finished"
        error-text="请求失败，点击重新加载"
        finished-text="没有更多了"
        @load="onLoad"
      >
        <van-grid :gutter="10" :column-num="2">
          <van-grid-item v-for="item in finishedJobs" :key="item.id">
            <div class="failed" v-if="item.progress === 101">
              <div class="title">任务失败</div>
              <div class="opt">
                <van-button size="small" @click="showErrMsg(item)">详情</van-button>
                <van-button type="danger" @click="removeImage($event, item)" size="small">删除</van-button>
              </div>
            </div>
            <div class="job-item" v-else>
              <van-image :src="item['img_url']" :class="item['can_opt'] ? '' : 'upscale'" lazy-load @click="imageView(item)" fit="cover">
                <template v-slot:loading>
                  <van-loading type="spinner" size="20" />
                </template>
              </van-image>

              <div class="remove">
                <el-button type="danger" :icon="Delete" @click="removeImage($event, item)" circle />
                <el-button type="warning" v-if="item.publish" @click="publishImage($event, item, false)" circle>
                  <i class="iconfont icon-cancel-share"></i>
                </el-button>
                <el-button type="success" v-else @click="publishImage($event, item, true)" circle>
                  <i class="iconfont icon-share-bold"></i>
                </el-button>
                <el-button type="primary" @click="showPrompt(item)" circle>
                  <i class="iconfont icon-prompt"></i>
                </el-button>
              </div>
            </div>
          </van-grid-item>
        </van-grid>
      </van-list>
    </div>
    <button style="display: none" class="copy-prompt-sd" :data-clipboard-text="prompt" id="copy-btn-sd">复制</button>
  </div>
</template>

<script setup>
import { onMounted, onUnmounted, ref } from "vue";
import { Delete } from "@element-plus/icons-vue";
import { httpGet, httpPost } from "@/utils/http";
import Clipboard from "clipboard";
import { checkSession, getClientId, getSystemInfo } from "@/store/cache";
import { useRouter } from "vue-router";
import { getSessionId } from "@/store/session";
import { showConfirmDialog, showDialog, showFailToast, showImagePreview, showNotify, showSuccessToast, showToast } from "vant";
import { showLoginDialog } from "@/utils/libs";
import { useSharedStore } from "@/store/sharedata";

const listBoxHeight = ref(window.innerHeight - 40);
const mjBoxHeight = ref(window.innerHeight - 150);
const isLogin = ref(false);
const activeColspan = ref([""]);

window.onresize = () => {
  listBoxHeight.value = window.innerHeight - 40;
  mjBoxHeight.value = window.innerHeight - 150;
};
const samplers = ref([
  { text: "Euler a", value: "Euler a" },
  { text: "DPM++ 2S a Karras", value: "DPM++ 2S a Karras" },
  { text: "DPM++ 2M Karras", value: "DPM++ 2M Karras" },
  { text: "DPM++ 2M SDE Karras", value: "DPM++ 2M SDE Karras" },
  { text: "DPM++ 2M Karras", value: "DPM++ 2M Karras" },
  { text: "DPM++ 3M SDE Karras", value: "DPM++ 3M SDE Karras" },
]);
const showSamplerPicker = ref(false);

const upscaleAlgArr = ref([
  { text: "Latent", value: "Latent" },
  { text: "ESRGAN_4x", value: "ESRGAN_4x" },
  { text: "ESRGAN 4x+", value: "ESRGAN 4x+" },
  { text: "SwinIR_4x", value: "SwinIR_4x" },
  { text: "LDSR", value: "LDSR" },
]);
const showUpscalePicker = ref(false);

const params = ref({
  client_id: getClientId(),
  width: 1024,
  height: 1024,
  sampler: samplers.value[0].value,
  seed: -1,
  steps: 20,
  cfg_scale: 7,
  hd_fix: false,
  hd_redraw_rate: 0.7,
  hd_scale: 2,
  hd_scale_alg: upscaleAlgArr.value[0].value,
  hd_steps: 0,
  prompt: "",
  neg_prompt: "nsfw, paintings,low quality,easynegative,ng_deepnegative ,lowres,bad anatomy,bad hands,bad feet",
});

const runningJobs = ref([]);
const finishedJobs = ref([]);
const router = useRouter();
// 检查是否有画同款的参数
const _params = router.currentRoute.value.params["copyParams"];
if (_params) {
  params.value = JSON.parse(_params);
}
const power = ref(0);
const sdPower = ref(0); // 画一张 SD 图片消耗算力

const userId = ref(0);
const store = useSharedStore();
const clipboard = ref(null);
const prompt = ref("");
onMounted(() => {
  initData();
  clipboard.value = new Clipboard(".copy-prompt-sd");
  clipboard.value.on("success", () => {
    showNotify({ type: "success", message: "复制成功", duration: 1000 });
  });
  clipboard.value.on("error", () => {
    showNotify({ type: "danger", message: "复制失败", duration: 2000 });
  });

  getSystemInfo()
    .then((res) => {
      sdPower.value = res.data.sd_power;
      params.value.neg_prompt = res.data.sd_neg_prompt;
    })
    .catch((e) => {
      showNotify({ type: "danger", message: "获取系统配置失败：" + e.message });
    });

  store.addMessageHandler("sd", (data) => {
    if (data.channel !== "sd" || data.clientId !== getClientId()) {
      return;
    }
    if (data.body === "FINISH" || data.body === "FAIL") {
      page.value = 1;
      fetchFinishJobs(1);
    }
    fetchRunningJobs();
  });
});

onUnmounted(() => {
  clipboard.value.destroy();
  store.removeMessageHandler("sd");
});

const initData = () => {
  checkSession()
    .then((user) => {
      power.value = user["power"];
      userId.value = user.id;
      isLogin.value = true;
      fetchRunningJobs();
      fetchFinishJobs(1);
    })
    .catch(() => {
      loading.value = false;
    });
};

const fetchRunningJobs = () => {
  // 获取运行中的任务
  httpGet(`/api/sd/jobs?finish=0`)
    .then((res) => {
      const jobs = res.data.items;
      const _jobs = [];
      for (let i = 0; i < jobs.length; i++) {
        if (jobs[i].progress === -1) {
          showNotify({
            message: `任务ID：${jobs[i]["task_id"]} 原因：${jobs[i]["err_msg"]}`,
            type: "danger",
          });
          power.value += sdPower.value;
          continue;
        }
        _jobs.push(jobs[i]);
      }
      runningJobs.value = _jobs;
    })
    .catch((e) => {
      showNotify({ type: "danger", message: "获取任务失败：" + e.message });
    });
};

const loading = ref(false);
const finished = ref(false);
const error = ref(false);
const page = ref(0);
const pageSize = ref(10);
// 获取已完成的任务
const fetchFinishJobs = (page) => {
  loading.value = true;
  httpGet(`/api/sd/jobs?finish=1&page=${page}&page_size=${pageSize.value}`)
    .then((res) => {
      const jobs = res.data.items;
      if (jobs.length < pageSize.value) {
        finished.value = true;
      }
      const _jobs = [];
      for (let i = 0; i < jobs.length; i++) {
        if (jobs[i].progress === -1) {
          jobs[i]["thumb_url"] = jobs[i]["img_url"] + "?imageView2/1/w/480/h/600/q/75";
        }
        _jobs.push(jobs[i]);
      }
      if (page === 1) {
        finishedJobs.value = _jobs;
      } else {
        finishedJobs.value = finishedJobs.value.concat(_jobs);
      }
      loading.value = false;
    })
    .catch((e) => {
      loading.value = false;
      showNotify({ type: "danger", message: "获取任务失败：" + e.message });
    });
};

const onLoad = () => {
  page.value += 1;
  fetchFinishJobs(page.value);
};

// 创建绘图任务
const promptRef = ref(null);
const generate = () => {
  if (!isLogin.value) {
    return showLoginDialog(router);
  }

  if (params.value.prompt === "") {
    promptRef.value.focus();
    return showToast("请输入绘画提示词！");
  }

  if (!params.value.seed) {
    params.value.seed = -1;
  }
  params.value.session_id = getSessionId();
  httpPost("/api/sd/image", params.value)
    .then(() => {
      showSuccessToast("绘画任务推送成功，请耐心等待任务执行...");
      power.value -= sdPower.value;
      fetchRunningJobs();
    })
    .catch((e) => {
      showFailToast("任务推送失败：" + e.message);
    });
};

const showPrompt = (item) => {
  prompt.value = item.prompt;
  showConfirmDialog({
    title: "绘画提示词",
    message: item.prompt,
    confirmButtonText: "复制",
    cancelButtonText: "关闭",
  })
    .then(() => {
      document.querySelector("#copy-btn-sd").click();
    })
    .catch(() => {});
};

const showErrMsg = (item) => {
  showDialog({
    title: "错误详情",
    message: item["err_msg"],
  }).then(() => {
    // on close
  });
};

const removeImage = (event, item) => {
  event.stopPropagation();
  showConfirmDialog({
    title: "标题",
    message: "此操作将会删除任务和图片，继续操作码?",
  })
    .then(() => {
      httpGet("/api/sd/remove", { id: item.id, user_id: item.user })
        .then(() => {
          showSuccessToast("任务删除成功");
          fetchFinishJobs(1);
        })
        .catch((e) => {
          showFailToast("任务删除失败：" + e.message);
        });
    })
    .catch(() => {
      showToast("您取消了操作");
    });
};

// 发布图片到作品墙
const publishImage = (event, item, action) => {
  event.stopPropagation();
  let text = "图片发布";
  if (action === false) {
    text = "取消发布";
  }
  httpGet("/api/sd/publish", { id: item.id, action: action, user_id: item.user })
    .then(() => {
      showSuccessToast(text + "成功");
      item.publish = action;
    })
    .catch((e) => {
      showFailToast(text + "失败：" + e.message);
    });
};

const imageView = (item) => {
  showImagePreview([item["img_url"]]);
};

const samplerConfirm = (item) => {
  params.value.sampler = item.selectedOptions[0].text;
  showSamplerPicker.value = false;
};

const upscaleConfirm = (item) => {
  params.value.hd_scale_alg = item.selectedOptions[0].text;
  showUpscalePicker.value = false;
};

const showInfo = (message) => {
  showDialog({
    title: "参数说明",
    message: message,
  }).then(() => {
    // on close
  });
};
</script>

<style lang="stylus">
@import "@/assets/css/mobile/image-sd.styl"
</style>
