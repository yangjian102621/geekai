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
                    :distanceToScroll="200" :isLoading="loading" :isOver="over" @scrollReachBottom="getNext">
        <template #default="slotProp">
          <div class="list-item">
            <div class="image">
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

            <div class="prompt">
              <span>{{ slotProp.item.prompt }}</span>
              <el-icon class="copy-prompt" :data-clipboard-text="slotProp.item.prompt">
                <DocumentCopy/>
              </el-icon>
            </div>
          </div>
        </template>

        <template #footer>
          <div class="footer">
            <span>没有更多数据了</span>
            <i class="iconfont icon-face"></i>
          </div>
        </template>
      </v3-waterfall>
    </div>
  </div>
</template>

<script setup>
import {onMounted, ref} from "vue"
import {DocumentCopy, Picture} from "@element-plus/icons-vue";
import {httpGet} from "@/utils/http";
import {ElMessage} from "element-plus";
import Clipboard from "clipboard";

const list = ref([])
const loading = ref(true)
const over = ref(false)
const imgType = ref("mj") // 图片类别
const listBoxHeight = window.innerHeight - 71
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

onMounted(() => {
  const clipboard = new Clipboard('.copy-prompt');
  clipboard.on('success', () => {
    ElMessage.success({message: "复制成功！", duration: 500});
  })

  clipboard.on('error', () => {
    ElMessage.error('复制失败！');
  })
})

</script>

<style lang="stylus" scoped>
@import "@/assets/css/images-wall.styl"
</style>
