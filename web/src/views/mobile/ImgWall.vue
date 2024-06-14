<template>
  <div class="img-wall container">
    <van-nav-bar :title="title"/>

    <div class="content">
      <van-tabs v-model:active="activeName">
        <van-tab title="MidJourney" name="mj">
          <van-list
              v-model:error="data['mj'].error"
              v-model:loading="data['mj'].loading"
              :finished="data['mj'].finished"
              error-text="请求失败，点击重新加载"
              finished-text="没有更多了"
              @load="onLoad"
              style="height: 100%;width: 100%;"
          >
            <van-cell v-for="item in data['mj'].data" :key="item.id">
              <van-image :src="item['img_thumb']" @click="showPrompt(item)" fit="cover"/>
            </van-cell>
          </van-list>
        </van-tab>
        <van-tab title="StableDiffusion" name="sd">
          <van-list
              v-model:error="data['sd'].error"
              v-model:loading="data['sd'].loading"
              :finished="data['sd'].finished"
              error-text="请求失败，点击重新加载"
              finished-text="没有更多了"
              @load="onLoad"
          >
            <van-cell v-for="item in data['sd'].data" :key="item.id">
              <van-image :src="item['img_thumb']" @click="showPrompt(item)" fit="cover"/>
            </van-cell>
          </van-list>
        </van-tab>
        <van-tab title="DALLE3" name="dalle3">
          <van-empty description="功能正在开发中"/>
        </van-tab>
      </van-tabs>
    </div>
  </div>
</template>

<script setup>
import {onMounted, ref} from "vue";
import {httpGet, httpPost} from "@/utils/http";
import {showDialog, showFailToast, showSuccessToast} from "vant";
import {ElMessage} from "element-plus";

const title = ref('图片创作广场')
const activeName = ref("mj")
const data = ref({
  "mj": {
    loading: false,
    finished: false,
    error: false,
    page: 1,
    pageSize: 12,
    url: "/api/mj/jobs",
    data: []
  },
  "sd": {
    loading: false,
    finished: false,
    error: false,
    page: 1,
    pageSize: 12,
    url: "/api/sd/jobs",
    data: []
  },
  "dalle3": {
    loading: false,
    finished: false,
    error: false,
    page: 1,
    pageSize: 12,
    url: "/api/dalle3/jobs",
    data: []
  }
})

const onLoad = () => {
  const d = data.value[activeName.value]
  httpGet(`${d.url}?status=1&page=${d.page}&page_size=${d.pageSize}&publish=true`).then(res => {
    d.loading = false
    if (res.data.length === 0) {
      d.finished = true
      return
    }

    // 生成缩略图
    const imageList = res.data
    for (let i = 0; i < imageList.length; i++) {
      imageList[i]["img_thumb"] = imageList[i]["img_url"] + "?imageView2/4/w/300/h/0/q/75"
    }
    if (imageList.length < d.pageSize) {
      d.finished = true
    }
    if (d.data.length === 0) {
      d.data = imageList
    } else {
      d.data = d.data.concat(imageList)
    }
    d.page += 1
  }).catch(() => {
    d.error = true
    showFailToast("加载图片数据失败")
  })
};

const showPrompt = (item) => {
  showDialog({
    title: "绘画提示词",
    message: item.prompt,
  }).then(() => {
    // on close
  });
}
</script>

<style lang="stylus">
.img-wall {
  .content {
    padding-top 60px

    .van-cell__value {
      .van-image {
        width 100%
      }
    }
  }
}
</style>
