<template>
  <div class="page-images-wall">
    <div class="inner custom-scroll">
      <div class="header">
        <h2 class="text-xl pt-4 pb-4">AI 绘画作品墙</h2>
        <div class="settings">
          <el-radio-group v-model="imgType" @change="changeImgType">
            <el-radio-button value="mj"
              ><i class="iconfont icon-mj mr-1"></i>MidJourney</el-radio-button
            >
            <el-radio-button value="image"
              ><i class="iconfont icon-dalle mr-1"></i>AI Image</el-radio-button
            >
          </el-radio-group>
        </div>
      </div>
      <div class="waterfall" :style="{ height: listBoxHeight + 'px' }" id="waterfall-box">
        <Waterfall
          v-if="imgType === 'mj'"
          id="waterfall-mj"
          v-bind="galleryWaterfallBind"
          :list="data['mj']"
          :is-loading="loading"
          :is-over="isOver"
          @afterRender="loading = false"
        >
          <template #default="{ item, url }">
            <div
              class="gallery-wall-card image-task-item bg-gray-900 shadow-md overflow-hidden transition-all duration-300 ease-linear hover:shadow-md hover:shadow-purple-800"
            >
              <div class="image-task-preview overflow-hidden">
                <LazyImg
                  :url="url"
                  class="image-task-image cursor-pointer transition-transform duration-300 ease-linear"
                  @click="previewImg(item)"
                />
              </div>
              <div class="image-task-overlay">
                <div class="image-task-overlay-time">
                  {{ dateFormat(item.created_at) }}
                </div>
                <div class="image-task-tools">
                  <el-tooltip content="任务详情" placement="top">
                    <button type="button" class="image-task-tool" @click.stop="openDetail(item)">
                      <i class="iconfont icon-info text-[#6366f1]"></i>
                    </button>
                  </el-tooltip>
                </div>
              </div>
            </div>
          </template>
        </Waterfall>

        <Waterfall
          v-if="imgType === 'image'"
          id="waterfall-image"
          v-bind="galleryWaterfallBind"
          :list="data['image']"
          :is-loading="loading"
          :is-over="isOver"
          @afterRender="loading = false"
        >
          <template #default="{ item, url }">
            <div
              class="gallery-wall-card image-task-item bg-gray-900 shadow-md overflow-hidden transition-all duration-300 ease-linear hover:shadow-md hover:shadow-purple-800"
            >
              <div class="image-task-preview overflow-hidden">
                <LazyImg
                  :url="url"
                  class="image-task-image cursor-pointer transition-transform duration-300 ease-linear"
                  @click="previewImg(item)"
                />
              </div>
              <div class="image-task-overlay">
                <div class="image-task-overlay-time">
                  {{ dateFormat(item.created_at) }}
                </div>
                <div class="image-task-tools">
                  <el-tooltip content="任务详情" placement="top">
                    <button type="button" class="image-task-tool" @click.stop="openDetail(item)">
                      <i class="iconfont icon-info text-[#6366f1]"></i>
                    </button>
                  </el-tooltip>
                </div>
              </div>
            </div>
          </template>
        </Waterfall>

        <div class="flex flex-col items-center justify-center py-10">
          <img
            :src="waterfallOptions.loadProps.loading"
            class="max-w-[50px] max-h-[50px]"
            v-if="loading"
          />
          <div v-else>
            <button
              class="px-5 py-2 rounded-full bg-purple-700 text-md text-white cursor-pointer hover:bg-purple-800 transition-all duration-300"
              @click="getNext"
              v-if="!isOver"
            >
              加载更多
            </button>
            <div class="no-more-data" v-else>
              <span class="text-gray-500 mr-2">没有更多数据了</span>
              <i class="iconfont icon-face"></i>
            </div>
          </div>
        </div>

        <back-top :right="30" :bottom="30" />
      </div>
    </div>

    <el-image-viewer @close="closePreview" v-if="previewURL !== ''" :url-list="[previewURL]" />

    <el-dialog
      v-model="detailDialogVisible"
      title="任务详情"
      :width="detailKind === 'mj' ? '680px' : '600px'"
      :class="['gallery-detail-dialog', { 'mj-detail-dialog': detailKind === 'mj' }]"
      :close-on-click-modal="false"
    >
      <div
        v-if="detailKind === 'mj' && currentDetail"
        class="gallery-mj-detail-body detail-content"
      >
        <el-descriptions :column="1" border size="small">
          <el-descriptions-item label="任务 ID">
            <div class="mj-detail-copy-row">
              <span class="break-all font-mono text-[13px]">{{
                currentDetail.task_id || '-'
              }}</span>
              <el-tooltip v-if="currentDetail.task_id" content="复制" placement="top">
                <i
                  class="iconfont icon-copy mj-detail-copy-ico"
                  @click="copyText(currentDetail.task_id)"
                />
              </el-tooltip>
            </div>
          </el-descriptions-item>
          <el-descriptions-item label="任务类型">
            {{ mjDetailTypeLabel(currentDetail.type) }}
          </el-descriptions-item>
          <el-descriptions-item label="状态">
            {{ currentDetail.status || '-' }}
            <span
              v-if="currentDetail.progress != null && currentDetail.status !== 'success'"
              class="text-gray-500 ml-1"
            >
              ({{ currentDetail.progress }}%)
            </span>
          </el-descriptions-item>
          <el-descriptions-item label="原始提示词" v-if="detailTaskPayload?.prompt">
            <div class="mj-detail-copy-row">
              <span class="break-all mj-detail-text">{{ detailTaskPayload.prompt }}</span>
              <el-tooltip content="复制" placement="top">
                <i
                  class="iconfont icon-copy mj-detail-copy-ico"
                  @click="copyText(detailTaskPayload.prompt)"
                />
              </el-tooltip>
            </div>
          </el-descriptions-item>
          <el-descriptions-item label="负面提示词" v-if="detailTaskPayload?.neg_prompt">
            <span class="break-all mj-detail-text">{{ detailTaskPayload.neg_prompt }}</span>
          </el-descriptions-item>
          <el-descriptions-item label="完整提示词">
            <div class="mj-detail-copy-row align-start">
              <span class="break-all mj-detail-text">{{ currentDetail.prompt || '—' }}</span>
              <el-tooltip v-if="currentDetail.prompt" content="复制" placement="top">
                <i
                  class="iconfont icon-copy mj-detail-copy-ico"
                  @click="copyText(currentDetail.prompt)"
                />
              </el-tooltip>
            </div>
          </el-descriptions-item>
          <el-descriptions-item label="引用图片" v-if="detailTaskImages.length > 0">
            <div class="mj-detail-refimgs">
              <el-image
                v-for="(u, idx) in detailTaskImages"
                :key="'ref-' + idx"
                :src="u"
                :preview-src-list="detailTaskImages"
                fit="cover"
                class="mj-detail-refimg"
                preview-teleported
              />
            </div>
          </el-descriptions-item>
          <el-descriptions-item label="局部重绘" v-if="detailHasMask">
            <el-text type="info">已提交蒙版（内容略）</el-text>
          </el-descriptions-item>
          <el-descriptions-item
            label="生成结果"
            v-if="currentDetail.status === 'success' && currentDetail.img_url"
          >
            <el-image
              :src="getThumbURL(currentDetail.img_url, 240, 240)"
              :preview-src-list="[currentDetail.img_url]"
              fit="cover"
              class="mj-detail-result"
              preview-teleported
            />
          </el-descriptions-item>
          <el-descriptions-item label="消耗积分">
            {{ currentDetail.power ?? 0 }}
          </el-descriptions-item>
          <el-descriptions-item label="创建时间">
            {{ dateFormat(currentDetail.created_at) }}
          </el-descriptions-item>
          <el-descriptions-item
            label="错误信息"
            v-if="currentDetail.status === 'failed' && currentDetail.err_msg"
          >
            <el-text type="danger">{{ currentDetail.err_msg }}</el-text>
          </el-descriptions-item>
          <el-descriptions-item
            label="原始载荷"
            v-if="currentDetail.task_info && !detailTaskPayload"
          >
            <span class="text-gray-500 text-xs">无法解析 JSON，以下为原始文本：</span>
            <pre class="mj-detail-raw">{{ currentDetail.task_info }}</pre>
          </el-descriptions-item>
        </el-descriptions>
      </div>

      <div v-else-if="detailKind === 'image' && currentDetail" class="detail-content">
        <el-descriptions :column="1" border>
          <el-descriptions-item label="提示词">
            <div>
              <span>{{ currentDetail.prompt }}</span>
              <el-tooltip content="复制提示词" placement="top">
                <i
                  class="iconfont icon-copy ml-2 cursor-pointer"
                  @click="copyText(currentDetail.prompt)"
                />
              </el-tooltip>
            </div>
          </el-descriptions-item>

          <el-descriptions-item
            label="生成的图片"
            v-if="currentDetail.status === 'success' && currentDetail.img_url"
          >
            <el-image
              :src="getThumbURL(currentDetail.img_url, 200, 200)"
              :preview-src-list="[currentDetail.img_url]"
              fit="cover"
              style="width: 200px; height: 200px"
              preview-teleported
            />
          </el-descriptions-item>

          <el-descriptions-item
            label="参考图"
            v-if="currentDetail.params?.image && currentDetail.params.image.length > 0"
          >
            <div class="reference-images">
              <el-image
                v-for="(img, idx) in currentDetail.params.image"
                :key="idx"
                :src="getThumbURL(img, 100, 100)"
                :preview-src-list="currentDetail.params.image"
                :initial-index="idx"
                fit="cover"
                style="width: 100px; height: 100px; margin-right: 10px"
                preview-teleported
              />
            </div>
          </el-descriptions-item>

          <el-descriptions-item label="生图模型">
            {{ currentDetail.params?.model_name || '-' }}
          </el-descriptions-item>

          <el-descriptions-item label="消耗积分">
            {{ currentDetail.power || 0 }}
          </el-descriptions-item>

          <el-descriptions-item label="图片比例">
            {{ currentDetail.params?.aspect_ratio || '-' }}
          </el-descriptions-item>

          <el-descriptions-item label="图片尺寸">
            {{ currentDetail.params?.size || '-' }}
          </el-descriptions-item>

          <el-descriptions-item label="创建时间">
            {{ dateFormat(currentDetail.created_at) }}
          </el-descriptions-item>

          <el-descriptions-item
            label="错误信息"
            v-if="currentDetail.status === 'failed' && currentDetail.err_msg"
          >
            <el-text type="danger">{{ currentDetail.err_msg }}</el-text>
          </el-descriptions-item>
        </el-descriptions>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import BackTop from '@/components/BackTop.vue'
import { useSharedStore } from '@/store/sharedata'
import { httpGet } from '@/utils/http'
import { dateFormat, getThumbURL } from '@/utils/libs'
import { ElMessage } from 'element-plus'
import { computed, nextTick, ref } from 'vue'
import { LazyImg, Waterfall } from 'vue-waterfall-plugin-next'
import 'vue-waterfall-plugin-next/dist/style.css'

const store = useSharedStore()
const waterfallOptions = store.waterfallOptions

/** 作品墙：间距约 2px、无外围 gutter，视觉接近小红书信息流 */
const galleryWaterfallBind = computed(() => ({
  ...store.waterfallOptions,
  gutter: 3,
  hasAroundGutter: false,
}))

const data = ref({
  mj: [],
  image: [],
})
const loading = ref(true)
const isOver = ref(false)
const imgType = ref('mj')
const listBoxHeight = window.innerHeight - 124
const previewURL = ref('')

const previewImg = (item) => {
  previewURL.value = item.img_url
}

const closePreview = () => {
  previewURL.value = ''
}

const page = ref(0)
const pageSize = ref(15)

const getNext = () => {
  if (isOver.value) {
    return
  }

  loading.value = true
  page.value = page.value + 1
  let url = ''
  switch (imgType.value) {
    case 'mj':
      url = '/api/mj/imgWall'
      break
    case 'image':
      url = '/api/image/imgWall'
      break
  }
  httpGet(`${url}?page=${page.value}&page_size=${pageSize.value}`)
    .then((res) => {
      if (!res.data.items || res.data.items.length === 0) {
        isOver.value = true
        loading.value = false
        return
      }

      const imageList = res.data.items
      for (let i = 0; i < imageList.length; i++) {
        imageList[i]['img_thumb'] = getThumbURL(imageList[i]['img_url'], 300, 0)
      }
      if (data.value[imgType.value].length === 0) {
        data.value[imgType.value] = imageList
        return
      }

      if (imageList.length < pageSize.value) {
        isOver.value = true
      }
      data.value[imgType.value] = data.value[imgType.value].concat(imageList)
    })
    .catch((e) => {
      ElMessage.error('获取图片失败：' + e.message)
      loading.value = false
    })
}

getNext()

const changeImgType = () => {
  document.getElementById('waterfall-box').scrollTo(0, 0)
  page.value = 0
  data.value = {
    mj: [],
    image: [],
  }
  loading.value = true
  isOver.value = false
  nextTick(() => getNext())
}

/** ---------- 任务详情（与 Image / ImageMj 对齐） ---------- */
const detailDialogVisible = ref(false)
const detailKind = ref('mj')
const currentDetail = ref(null)

function parseMjTaskInfo(raw) {
  if (raw == null || raw === '') return null
  if (typeof raw === 'object') return raw
  if (typeof raw !== 'string') return null
  try {
    return JSON.parse(raw)
  } catch {
    return null
  }
}

const detailTaskPayload = computed(() => {
  if (detailKind.value !== 'mj' || !currentDetail.value) return null
  return parseMjTaskInfo(currentDetail.value.task_info)
})

const detailTaskImages = computed(() => {
  const arr = detailTaskPayload.value?.img_arr
  return Array.isArray(arr) ? arr.filter((u) => u && String(u).trim()) : []
})

const detailHasMask = computed(() => {
  const m = detailTaskPayload.value?.mask_base64
  return typeof m === 'string' && m.length > 0
})

function mjDetailTypeLabel(type) {
  const map = {
    image: '绘图',
    upscale: '放大',
    variation: '变换',
    blend: '融图',
    swapFace: '换脸',
    modal: '局部重绘',
  }
  return map[type] || type || '-'
}

const openDetail = (item) => {
  detailKind.value = imgType.value
  if (imgType.value === 'image') {
    let params = {}
    try {
      if (item.params) {
        params = typeof item.params === 'string' ? JSON.parse(item.params) : item.params
      }
    } catch (e) {
      console.error('解析 params 失败:', e)
    }
    currentDetail.value = { ...item, params }
  } else {
    currentDetail.value = { ...item }
  }
  detailDialogVisible.value = true
}

const copyText = (text) => {
  if (!text) return
  navigator.clipboard
    .writeText(text)
    .then(() => {
      ElMessage.success('复制成功！')
    })
    .catch(() => {
      ElMessage.error('复制失败！')
    })
}
</script>

<style lang="scss">
@use '../assets/css/images-wall.scss' as *;
@use '../assets/css/custom-scroll.scss' as *;
@use '../assets/css/image.scss' as *;
</style>

<!-- Dialog  teleport 到 body，与 ImageMj 一致使用非 scoped 样式 -->
<style lang="scss">
.gallery-detail-dialog.mj-detail-dialog {
  :deep(.el-dialog__body) {
    padding-top: 6px;
  }
}

.gallery-mj-detail-body {
  max-height: min(72vh, 640px);
  overflow-y: auto;
}

.mj-detail-copy-row {
  display: flex;
  align-items: center;
  gap: 8px;

  &.align-start {
    align-items: flex-start;

    .mj-detail-copy-ico {
      margin-top: 4px;
    }
  }
}

.mj-detail-text {
  font-size: 13px;
  line-height: 1.55;
}

.mj-detail-copy-ico {
  flex-shrink: 0;
  cursor: pointer;
  color: var(--el-text-color-secondary);
  opacity: 0.85;
  transition:
    opacity 0.15s ease,
    color 0.15s ease;

  &:hover {
    opacity: 1;
    color: var(--el-color-primary);
  }
}

.mj-detail-refimgs {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.mj-detail-refimg {
  width: 88px;
  height: 88px;
  border-radius: 0;
  overflow: hidden;
  border: 1px solid var(--el-border-color-lighter);
}

.mj-detail-result {
  width: 240px;
  max-width: 100%;
  height: 240px;
  border-radius: 0;
  overflow: hidden;
  border: 1px solid var(--el-border-color-lighter);
}

.mj-detail-raw {
  display: block;
  margin-top: 8px;
  padding: 8px 10px;
  font-size: 11px;
  line-height: 1.45;
  max-height: 180px;
  overflow: auto;
  white-space: pre-wrap;
  word-break: break-all;
  border-radius: 0;
  background: var(--el-fill-color-light);
}

.page-images-wall .gallery-wall-card,
.page-images-wall .gallery-wall-card .image-task-preview {
  border-radius: 0;
}

.page-images-wall .gallery-wall-card :deep(img),
.page-images-wall .gallery-wall-card :deep(.lazy__img) {
  border-radius: 0 !important;
}

.gallery-detail-dialog :deep(.el-image__inner),
.gallery-detail-dialog :deep(.el-image__wrapper) {
  border-radius: 0 !important;
}

.gallery-wall-card:hover .image-task-image {
  transform: scale(1.05);
}
</style>
