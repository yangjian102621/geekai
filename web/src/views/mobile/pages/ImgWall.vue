<template>
  <div class="img-wall container">
    <div class="content">
      <van-tabs v-model:active="activeName" animated sticky>
        <van-tab title="MJ" name="mj">
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
              <van-image :src="item['img_thumb']" @click="imageView(item)" fit="cover"/>

              <div class="opt-box">
                <el-button type="primary" @click="showPrompt(item)" circle>
                  <i class="iconfont icon-prompt"></i>
                </el-button>
              </div>
            </van-cell>
          </van-list>
        </van-tab>
        <van-tab title="SD" name="sd">
          <van-list
              v-model:error="data['sd'].error"
              v-model:loading="data['sd'].loading"
              :finished="data['sd'].finished"
              error-text="请求失败，点击重新加载"
              finished-text="没有更多了"
              @load="onLoad"
          >
            <van-cell v-for="item in data['sd'].data" :key="item.id">
              <van-image :src="item['img_thumb']" @click="imageView(item)" fit="cover"/>

              <div class="opt-box">
                <el-button type="primary" @click="showPrompt(item)" circle>
                  <i class="iconfont icon-prompt"></i>
                </el-button>
              </div>
            </van-cell>
          </van-list>
        </van-tab>
        <van-tab title="DALL" name="dall">
          <van-list
              v-model:error="data['dall'].error"
              v-model:loading="data['dall'].loading"
              :finished="data['dall'].finished"
              error-text="请求失败，点击重新加载"
              finished-text="没有更多了"
              @load="onLoad"
          >
            <van-cell v-for="item in data['dall'].data" :key="item.id">
              <van-image :src="item['img_thumb']" @click="imageView(item)" fit="cover"/>

              <div class="opt-box">
                <el-button type="primary" @click="showPrompt(item)" circle>
                  <i class="iconfont icon-prompt"></i>
                </el-button>
              </div>
            </van-cell>
          </van-list>
        </van-tab>
      </van-tabs>
    </div>

    <button style="display: none" class="copy-prompt-wall" :data-clipboard-text="prompt" id="copy-btn-wall">复制
    </button>
  </div>
</template>

<script setup>
import {onMounted, onUnmounted, ref} from "vue";
import {httpGet} from "@/utils/http";
import {showConfirmDialog, showFailToast, showImagePreview, showNotify} from "vant";
import Clipboard from "clipboard";
import {ElMessage} from "element-plus";

const activeName = ref("mj")
const data = ref({
  "mj": {
    loading: false,
    finished: false,
    error: false,
    page: 1,
    pageSize: 12,
    url: "/api/mj/imgWall",
    data: []
  },
  "sd": {
    loading: false,
    finished: false,
    error: false,
    page: 1,
    pageSize: 12,
    url: "/api/sd/imgWall",
    data: []
  },
  "dall": {
    loading: false,
    finished: false,
    error: false,
    page: 1,
    pageSize: 12,
    url: "/api/dall/imgWall",
    data: []
  }
})

const prompt = ref('')
const clipboard = ref(null)
onMounted(() => {
  clipboard.value = new Clipboard(".copy-prompt-wall");
  clipboard.value.on('success', () => {
    showNotify({type: 'success', message: '复制成功', duration: 1000})
  })
  clipboard.value.on('error', () => {
    showNotify({type: 'danger', message: '复制失败', duration: 2000})
  })


  clipboard.value.on('error', () => {
    ElMessage.error('复制失败！');
  })
})

onUnmounted(() => {
  clipboard.value.destroy()
})

const onLoad = () => {
  const d = data.value[activeName.value]
  httpGet(`${d.url}?status=1&page=${d.page}&page_size=${d.pageSize}&publish=true`).then(res => {
    d.loading = false
    if (res.data.items.length === 0) {
      d.finished = true
      return
    }

    // 生成缩略图
    const imageList = res.data.items
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
  prompt.value = item.prompt
  showConfirmDialog({
    title: "绘画提示词",
    message: item.prompt,
    confirmButtonText: "复制",
    cancelButtonText: "关闭",
  }).then(() => {
    document.querySelector('#copy-btn-wall').click()
  }).catch(() => {
  });
}

const imageView = (item) => {
  showImagePreview([item['img_url']]);
}
</script>

<style lang="stylus">
.img-wall {
  .content {
    .van-cell__value {
      min-height 80px

      .van-image {
        width 100%
      }

      .opt-box {
        position absolute
        right 0
        top 0
        padding 10px
      }
    }
  }
}
</style>
