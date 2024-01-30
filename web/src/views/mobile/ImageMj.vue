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
                <div class="text">{{ item.text }}</div>
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
          <van-field label="垫图">
            <template #input>
              <van-uploader v-model="imgList" :after-read="afterRead"/>
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

        <div class="text-line">
          <van-button round block type="primary" native-type="submit">
            立即生成
          </van-button>
        </div>
      </van-form>
    </div>
  </div>
</template>

<script setup>
import {onMounted, ref} from "vue";
import {showFailToast, showNotify, showSuccessToast, showToast} from "vant";
import {httpGet, httpPost} from "@/utils/http";
import Compressor from 'compressorjs';

const title = ref('MidJourney 绘画')

const rates = [
  {css: "square", value: "1:1", text: "1:1", img: "/images/mj/rate_1_1.png"},
  {css: "size2-3", value: "2:3", text: "2:3", img: "/images/mj/rate_3_4.png"},
  {css: "size3-4", value: "3:4", text: "3:4", img: "/images/mj/rate_3_4.png"},
  {css: "size4-3", value: "4:3", text: "4:3", img: "/images/mj/rate_4_3.png"},
  {css: "size16-9", value: "16:9", text: "16:9", img: "/images/mj/rate_16_9.png"},
  {css: "size9-16", value: "9:16", text: "9:16", img: "/images/mj/rate_9_16.png"},
]
const models = [
  {text: "写实模式MJ-6.0", value: " --v 6", img: "/images/mj/mj-v6.png"},
  {text: "优质模式MJ-5.2", value: " --v 5.2", img: "/images/mj/mj-v5.2.png"},
  {text: "动漫风niji5 原始", value: " --niji 5", img: "/images/mj/mj-niji.png"},
  {text: "动漫风niji5 可爱", value: " --niji 5 --style cute", img: "/images/mj/nj1.jpg"},
  {text: "动漫风niji5 风景", value: " --niji 5 --style scenic", img: "/images/mj/nj2.jpg"},
  {text: "动漫风niji5 表现力", value: " --niji 5 --style expressive", img: "/images/mj/nj3.jpg"},
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

// 切换图片比例
const changeRate = (item) => {
  params.value.rate = item.value
}
// 切换模型
const changeModel = (item) => {
  params.value.model = item.value
}


const generate = () => {
  httpPost('/api/user/profile/update', form.value).then(() => {
    showSuccessToast('保存成功')
  }).catch(() => {
    showFailToast('保存失败')
  })
}
</script>

<style lang="stylus" scoped>
@import "@/assets/css/mobile/image-mj.styl"
</style>