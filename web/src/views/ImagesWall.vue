<template>
  <div class="page-images-wall">
    <div class="inner custom-scroll">
      <div class="header">
        <h2>AI 绘画作品墙</h2>
        <div class="settings">
          <el-radio-group v-model="imgType">
            <el-radio label="mj" size="large">MidJourney</el-radio>
            <el-radio label="sd" size="large">Stable Diffusion</el-radio>
          </el-radio-group>
        </div>
      </div>
      <v3-waterfall class="waterfall" id="waterfall" :list="list" srcKey="img_url" :gap="12" :colWidth="colWidth"
                    :style="{ height:listBoxHeight + 'px' }"
                    :distanceToScroll="100" :isLoading="loading" :isOver="over" @scrollReachBottom="getNext">
        <template #default="slotProp">
          <div class="list-item">
            <el-image :src="slotProp.item['img_url']+'?imageView2/4/w/300/q/75'"
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
                  <el-icon v-if="slotProp.item['img'] !== ''">
                    <Picture/>
                  </el-icon>
                </div>
              </template>
            </el-image>
          </div>
        </template>
      </v3-waterfall>
    </div>
  </div>
</template>

<script setup>
import {ref} from "vue"
import {Picture} from "@element-plus/icons-vue";
import {httpGet} from "@/utils/http";
import {ElMessage} from "element-plus";

const list = ref([])
const loading = ref(true)
const over = ref(false)
const imgType = ref("mj") // 图片类别
const listBoxHeight = window.innerHeight - 100
const colWidth = ref(240)

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
  loading.value = true
  page.value = page.value + 1
  // 获取运行中的任务
  httpGet(`/api/mj/jobs?status=1&page=${page.value}&page_size=${pageSize.value}`).then(res => {
    loading.value = false
    if (list.value.length === 0) {
      list.value = res.data
      return
    } else if (res.data.length < pageSize.value) {
      over.value = true
    }

    list.value = list.value.concat(res.data)

  }).catch(e => {
    ElMessage.error("获取图片失败：" + e.message)
  })
}

getNext()

</script>

<style lang="stylus" scoped>
@import "@/assets/css/images-wall.styl"
</style>
