<template>
  <div class="page-images-wall">
    <div class="inner custom-scroll">
      <div class="header">
        <h2>AI 绘画作品墙</h2>
        <div class="settings">
          <el-radio-group v-model="imgType" @change="changeImgType">
            <el-radio label="mj" size="large">MidJourney</el-radio>
            <el-radio label="sd" size="large">Stable Diffusion</el-radio>
          </el-radio-group>
        </div>
      </div>
      <div class="waterfall" :style="{ height:listBoxHeight + 'px' }" id="waterfall-box">
        <v3-waterfall v-if="imgType === 'mj'"
                      id="waterfall"
                      :list="list"
                      srcKey="img_thumb"
                      :gap="12"
                      :bottomGap="-5"
                      :colWidth="colWidth"
                      :distanceToScroll="100"
                      :isLoading="loading"
                      :isOver="false"
                      @scrollReachBottom="getNext">
          <template #default="slotProp">
            <div class="list-item">
              <div class="image">
                <el-image :src="slotProp.item['img_thumb']"
                          :zoom-rate="1.2"
                          :preview-src-list="[slotProp.item['img_url']]"
                          :preview-teleported="true"
                          :initial-index="10"
                          loading="lazy">
                  <template #placeholder>
                    <div class="image-slot">
                      正在加载图片
                    </div>
                  </template>

                  <template #error>
                    <div class="image-slot">
                      <el-icon>
                        <Picture/>
                      </el-icon>
                    </div>
                  </template>
                </el-image>
              </div>
              <div class="prompt">
                <span>{{ slotProp.item.prompt }}</span>
                <el-icon class="copy-prompt" :data-clipboard-text="slotProp.item.prompt">
                  <DocumentCopy/>
                </el-icon>
              </div>
            </div>
          </template>
        </v3-waterfall>

        <v3-waterfall v-else
                      id="waterfall"
                      :list="list"
                      srcKey="img_thumb"
                      :gap="12"
                      :bottomGap="-5"
                      :colWidth="colWidth"
                      :distanceToScroll="100"
                      :isLoading="loading"
                      :isOver="false"
                      @scrollReachBottom="getNext">
          <template #default="slotProp">
            <div class="list-item">
              <div class="image">
                <el-image :src="slotProp.item['img_thumb']" loading="lazy"
                          @click="showTask(slotProp.item)">
                  <template #placeholder>
                    <div class="image-slot">
                      正在加载图片
                    </div>
                  </template>

                  <template #error>
                    <div class="image-slot">
                      <el-icon>
                        <Picture/>
                      </el-icon>
                    </div>
                  </template>
                </el-image>
              </div>

              <div class="prompt">
                <span>{{ slotProp.item.prompt }}</span>
                <el-icon class="copy-prompt" :data-clipboard-text="slotProp.item.prompt">
                  <DocumentCopy/>
                </el-icon>
              </div>
            </div>
          </template>
        </v3-waterfall>

        <div class="footer" v-if="isOver">
          <span>没有更多数据了</span>
          <i class="iconfont icon-face"></i>
        </div>

      </div>
    </div>
    <!-- 任务详情弹框 -->
    <el-dialog v-model="showTaskDialog" title="绘画任务详情" :fullscreen="true">
      <el-row :gutter="20">
        <el-col :span="16">
          <div class="img-container" :style="{maxHeight: fullImgHeight+'px'}">
            <el-image :src="item['img_url']" fit="contain">
              <template #placeholder>
                <div class="image-slot">
                  正在加载图片
                </div>
              </template>

              <template #error>
                <div class="image-slot">
                  <el-icon>
                    <Picture/>
                  </el-icon>
                </div>
              </template>
            </el-image>
          </div>
        </el-col>
        <el-col :span="8">
          <div class="task-info">
            <div class="info-line">
              <el-divider>
                正向提示词
              </el-divider>
              <div class="prompt">
                <span>{{ item.prompt }}</span>
                <el-icon class="copy-prompt" :data-clipboard-text="item.prompt">
                  <DocumentCopy/>
                </el-icon>
              </div>

            </div>

            <div class="info-line">
              <el-divider>
                反向提示词
              </el-divider>
              <div class="prompt">
                <span>{{ item.params.negative_prompt }}</span>
                <el-icon class="copy-prompt" :data-clipboard-text="item.params.negative_prompt">
                  <DocumentCopy/>
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
              <el-divider>
                高清修复
              </el-divider>
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
              <el-button type="primary" round @click="copyParams(item)">画一张同款的</el-button>
            </div>

          </div>
        </el-col>
      </el-row>

    </el-dialog>
  </div>
</template>

<script setup>
import {nextTick, onMounted, ref} from "vue"
import {DocumentCopy, Picture} from "@element-plus/icons-vue";
import {httpGet} from "@/utils/http";
import {ElMessage} from "element-plus";
import Clipboard from "clipboard";
import {useRouter} from "vue-router";

const list = ref([])
const loading = ref(true)
const isOver = ref(false)
const imgType = ref("mj") // 图片类别
const listBoxHeight = window.innerHeight - 74
const colWidth = ref(240)
const fullImgHeight = ref(window.innerHeight - 60)
const showTaskDialog = ref(false)
const item = ref({})

// 计算瀑布流列宽度
const calcColWidth = () => {
  const listBoxWidth = window.innerWidth - 60 - 80
  const rows = Math.floor(listBoxWidth / colWidth.value)
  colWidth.value = Math.floor((listBoxWidth - (rows - 1) * 12) / rows)
}
calcColWidth()
window.onresize = () => {
  calcColWidth()
}

const page = ref(0)
const pageSize = ref(20)
// 获取下一页数据
const getNext = () => {
  if (isOver.value) {
    return
  }

  loading.value = true
  page.value = page.value + 1
  const url = imgType.value === "mj" ? "/api/mj/jobs" : "/api/sd/jobs"
  // 获取运行中的任务
  httpGet(`${url}?status=1&page=${page.value}&page_size=${pageSize.value}`).then(res => {
    loading.value = false
    if (res.data.length === 0) {
      isOver.value = true
      return
    }

    // 生成缩略图
    const imageList = res.data
    for (let i = 0; i < imageList.length; i++) {
      imageList[i]["img_thumb"] = imageList[i]["img_url"] + "?imageView2/4/w/300/h/0/q/75"
    }
    if (list.value.length === 0) {
      list.value = imageList
      return
    }

    if (imageList.length < pageSize.value) {
      isOver.value = true
    }
    list.value = list.value.concat(imageList)

  }).catch(e => {
    ElMessage.error("获取图片失败：" + e.message)
  })
}

getNext()

onMounted(() => {
  const clipboard = new Clipboard('.copy-prompt');
  clipboard.on('success', () => {
    ElMessage.success("复制成功！");
  })

  clipboard.on('error', () => {
    ElMessage.error('复制失败！');
  })
})

const changeImgType = () => {
  document.getElementById('waterfall-box').scrollTo(0, 0)
  page.value = 0
  list.value = []
  loading.value = true
  isOver.value = false
  nextTick(() => getNext())
}

const showTask = (row) => {
  item.value = row
  showTaskDialog.value = true
}


const router = useRouter()
const copyParams = (row) => {
  router.push({name: "image-sd", params: {copyParams: JSON.stringify(row.params)}})
}

</script>

<style lang="stylus">
@import "@/assets/css/images-wall.styl"
@import "@/assets/css/custom-scroll.styl"
</style>
