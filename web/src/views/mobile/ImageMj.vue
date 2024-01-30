<template>
  <div class="mobile-mj container">
    <van-nav-bar :title="title"/>

    <div class="content">
      <van-form @submit="generate">
        <div class="text-line">图片比例</div>
        <div class="text-line">
          <van-row :gutter="10">
            <van-col :span="4" v-for="item in rates" :key="item.value">
              <div
                  :class="item.value === params.rate ? 'rate active' : 'rate'"
                  @click="changeRate(item)">
                <div class="icon">
                  <van-image :src="item.img" fit="cover"></van-image>
                </div>
                <div class="text">{{ item.text }}</div>
              </div>
            </van-col>
          </van-row>
        </div>
        <div class="text-line">模型选择</div>
        <div class="text-line">
          <van-row :gutter="10">
            <van-col :span="8" v-for="item in models" :key="item.value">
              <div :class="item.value === params.model ? 'model active' : 'model'"
                   @click="changeModel(item)">
                <div class="icon">
                  <van-image :src="item.img" fit="cover"></van-image>
                </div>
                <div class="text">
                  <van-text-ellipsis :content="item.text"/>
                </div>
              </div>
            </van-col>
          </van-row>
        </div>
        <div class="text-line">
          <van-field label="创意度">
            <template #input>
              <van-slider v-model.number="params.chaos" :max="100" :step="1"
                          @update:model-value="showToast('当前值：' + params.chaos)"/>
            </template>
          </van-field>
        </div>

        <div class="text-line">
          <van-field label="风格化">
            <template #input>
              <van-slider v-model.number="params.stylize" :max="1000" :step="1"
                          @update:model-value="showToast('当前值：' + params.stylize)"/>
            </template>
          </van-field>
        </div>

        <div class="text-line">
          <van-field label="原始模式">
            <template #input>
              <van-switch v-model="params.raw"/>
            </template>
          </van-field>
        </div>

        <div class="text-line">
          <van-field
              v-model="params.prompt"
              rows="3"
              autosize
              label="提示词"
              type="textarea"
              placeholder="如：一个美丽的中国女孩站在电影院门口，手上拿着爆米花，微笑，写实风格，电影灯光效果，半身像"
          />
        </div>

        <van-collapse v-model="activeColspan">
          <van-collapse-item title="垫图" name="img">
            <van-field>
              <template #input>
                <van-uploader v-model="imgList" :after-read="uploadImg"/>
              </template>
            </van-field>
          </van-collapse-item>
          <van-collapse-item title="反向提示词" name="neg_prompt">
            <van-field
                v-model="params.prompt"
                rows="3"
                autosize
                type="textarea"
                placeholder="不想出现在图片上的元素(例如：树，建筑)"
            />
          </van-collapse-item>
        </van-collapse>

        <div class="text-line">
          <van-button round block type="primary" native-type="submit">
            <van-tag type="success">可用额度:{{ imgCalls }}</van-tag>
            立即生成
          </van-button>
        </div>
      </van-form>

      <h2>任务列表</h2>

    </div>
  </div>
</template>

<script setup>
import {onMounted, ref} from "vue";
import {showFailToast, showToast} from "vant";
import {httpPost} from "@/utils/http";
import Compressor from "compressorjs";
import {ElMessage} from "element-plus";
import {getSessionId} from "@/store/session";
import {checkSession} from "@/action/session";
import Clipboard from "clipboard";
import {useRouter} from "vue-router";

const title = ref('MidJourney 绘画')
const activeColspan = ref([""])

const rates = [
  {css: "square", value: "1:1", text: "1:1", img: "/images/mj/rate_1_1.png"},
  {css: "size2-3", value: "2:3", text: "2:3", img: "/images/mj/rate_3_4.png"},
  {css: "size3-4", value: "3:4", text: "3:4", img: "/images/mj/rate_3_4.png"},
  {css: "size4-3", value: "4:3", text: "4:3", img: "/images/mj/rate_4_3.png"},
  {css: "size16-9", value: "16:9", text: "16:9", img: "/images/mj/rate_16_9.png"},
  {css: "size9-16", value: "9:16", text: "9:16", img: "/images/mj/rate_9_16.png"},
]
const models = [
  {text: "MJ-6.0", value: " --v 6", img: "/images/mj/mj-v6.png"},
  {text: "MJ-5.2", value: " --v 5.2", img: "/images/mj/mj-v5.2.png"},
  {text: "Niji5 原始", value: " --niji 5", img: "/images/mj/mj-niji.png"},
  {text: "Niji5 可爱", value: " --niji 5 --style cute", img: "/images/mj/nj1.jpg"},
  {text: "Niji5 风景", value: " --niji 5 --style scenic", img: "/images/mj/nj2.jpg"},
  {text: "Niji5 表现力", value: " --niji 5 --style expressive", img: "/images/mj/nj3.jpg"},
]
const imgList = ref([])
const params = ref({
  task_type: "image",
  rate: rates[0].value,
  model: models[0].value,
  chaos: 0,
  stylize: 100,
  seed: 0,
  img_arr: [],
  raw: false,
  weight: 0.25,
  prompt: "",
  neg_prompt: "",
  tile: false,
  quality: 0
})
const imgCalls = ref(0)
const userId = ref(0)
const router = useRouter()
onMounted(() => {
  checkSession().then(user => {
    imgCalls.value = user['img_calls']
    userId.value = user.id

    // fetchRunningJobs(userId.value)
    // fetchFinishJobs(userId.value)
    // connect()

  }).catch(() => {
    router.push('/login')
  });

  const clipboard = new Clipboard('.copy-prompt');
  clipboard.on('success', () => {
    ElMessage.success("复制成功！");
  })

  clipboard.on('error', () => {
    ElMessage.error('复制失败！');
  })
})
// 切换图片比例
const changeRate = (item) => {
  params.value.rate = item.value
}
// 切换模型
const changeModel = (item) => {
  params.value.model = item.value
}


// 图片上传
const uploadImg = (file) => {
  file.status = "uploading"
  // 压缩图片并上传
  new Compressor(file.file, {
    quality: 0.6,
    success(result) {
      const formData = new FormData();
      formData.append('file', result, result.name);
      // 执行上传操作
      httpPost('/api/upload', formData).then(res => {
        file.url = res.data.url
        file.status = "done"
      }).catch(e => {
        file.status = 'failed'
        file.message = '上传失败'
        showFailToast("图片上传失败：" + e.message)
      })
    },
    error(err) {
      console.log(err.message);
    },
  });
};

const generate = () => {
  if (params.value.prompt === '' && params.value.task_type === "image") {
    return showFailToast("请输入绘画提示词！")
  }
  if (params.value.model.indexOf("niji") !== -1 && params.value.raw) {
    return showFailToast("动漫模型不允许启用原始模式")
  }
  params.value.session_id = getSessionId()
  params.value.img_arr = imgList.value.map(img => img.url)
  httpPost("/api/mj/image", params.value).then(() => {
    ElMessage.success("绘画任务推送成功，请耐心等待任务执行...")
    imgCalls.value -= 1
  }).catch(e => {
    ElMessage.error("任务推送失败：" + e.message)
  })
}
</script>

<style lang="stylus" scoped>
@import "@/assets/css/mobile/image-mj.styl"
</style>