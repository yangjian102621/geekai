<template>
  <div class="container image-page">
    <el-tabs v-model="activeName" @tab-change="handleChange">
      <el-tab-pane label="Midjourney" name="mj" v-loading="data.mj.loading">
        <div class="handle-box">
          <el-input
            v-model="data.mj.query.username"
            placeholder="用户名"
            class="handle-input mr10"
            @keyup="search($event, 'mj')"
            clearable
          />
          <el-input
            v-model="data.mj.query.prompt"
            placeholder="提示词"
            class="handle-input mr10"
            @keyup="search($event, 'mj')"
            clearable
          />
          <el-date-picker
            v-model="data.mj.query.created_at"
            type="daterange"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            format="YYYY-MM-DD"
            value-format="YYYY-MM-DD"
            style="margin-right: 10px; width: 200px; position: relative; top: 3px"
          />
          <el-button type="primary" :icon="Search" @click="fetchMjData">搜索</el-button>
        </div>

        <div v-if="data.mj.items.length > 0">
          <el-row>
            <el-table :data="data.mj.items" :row-key="(row) => row.id" table-layout="auto">
              <el-table-column prop="user_id" label="用户ID" />
              <el-table-column label="任务类型">
                <template #default="scope">
                  <el-button :color="taskTypeTheme[scope.row.type].color" size="small" plain>{{
                    taskTypeTheme[scope.row.type].text
                  }}</el-button>
                </template>
              </el-table-column>
              <el-table-column prop="progress" label="任务进度">
                <template #default="scope">
                  <span v-if="scope.row.progress <= 100">{{ scope.row.progress }}%</span>
                  <el-tag v-else type="danger">已失败</el-tag>
                </template>
              </el-table-column>
              <el-table-column prop="power" label="消耗算力" />
              <el-table-column label="结果图片">
                <template #default="scope">
                  <el-button
                    size="small"
                    type="success"
                    @click="showImage(scope.row.img_url)"
                    v-if="scope.row.img_url !== ''"
                    plain
                    >预览图片</el-button
                  >
                </template>
              </el-table-column>
              <el-table-column label="提示词">
                <template #default="scope">
                  <el-popover
                    placement="top-start"
                    title="绘画提示词"
                    :width="300"
                    trigger="hover"
                    :content="scope.row.prompt"
                  >
                    <template #reference>
                      <span>{{ substr(scope.row.prompt, 20) }}</span>
                    </template>
                  </el-popover>
                </template>
              </el-table-column>
              <el-table-column label="创建时间">
                <template #default="scope">
                  <span>{{ dateFormat(scope.row['created_at']) }}</span>
                </template>
              </el-table-column>
              <el-table-column label="失败原因">
                <template #default="scope">
                  <el-popover
                    placement="top-start"
                    title="失败原因"
                    :width="300"
                    trigger="hover"
                    :content="scope.row.err_msg"
                    v-if="scope.row.progress === 101"
                  >
                    <template #reference>
                      <el-text type="danger">{{ substr(scope.row.err_msg, 20) }}</el-text>
                    </template>
                  </el-popover>
                  <span v-else>无</span>
                </template>
              </el-table-column>
              <el-table-column label="操作" width="180">
                <template #default="scope">
                  <el-popconfirm title="确定要删除当前记录吗?" @confirm="remove(scope.row, 'mj')">
                    <template #reference>
                      <el-button size="small" type="danger">删除</el-button>
                    </template>
                  </el-popconfirm>
                </template>
              </el-table-column>
            </el-table>
          </el-row>

          <div class="pagination">
            <el-pagination
              v-if="data.mj.total > 0"
              background
              layout="total,prev, pager, next"
              :hide-on-single-page="true"
              v-model:current-page="data.mj.page"
              v-model:page-size="data.mj.pageSize"
              @current-change="fetchMjData()"
              :total="data.mj.total"
            />
          </div>
        </div>
        <el-empty v-else />
      </el-tab-pane>
      <el-tab-pane label="AI图像生成" name="image">
        <div class="handle-box">
          <el-input
            v-model="data.image.query.username"
            placeholder="用户名"
            class="handle-input mr10"
            @keyup="search($event, 'image')"
            clearable
          />
          <el-input
            v-model="data.image.query.prompt"
            placeholder="提示词"
            class="handle-input mr10"
            @keyup="search($event, 'image')"
            clearable
          />
          <el-date-picker
            v-model="data.image.query.created_at"
            type="daterange"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            format="YYYY-MM-DD"
            value-format="YYYY-MM-DD"
            style="margin-right: 10px; width: 200px; position: relative; top: 3px"
          />
          <el-button type="primary" :icon="Search" @click="fetchImageData">搜索</el-button>
        </div>

        <div v-if="data.image.items.length > 0">
          <el-row>
            <el-table :data="data.image.items" :row-key="(row) => row.id" table-layout="auto">
              <el-table-column prop="user_id" label="用户ID" />
              <el-table-column prop="progress" label="任务进度">
                <template #default="scope">
                  <span v-if="scope.row.progress <= 100">{{ scope.row.progress }}%</span>
                  <el-tag v-else type="danger">已失败</el-tag>
                </template>
              </el-table-column>
              <el-table-column prop="power" label="消耗算力" />
              <el-table-column label="结果图片">
                <template #default="scope">
                  <el-button
                    size="small"
                    type="success"
                    @click="showImage(scope.row.img_url)"
                    v-if="scope.row.img_url !== ''"
                    plain
                    >预览图片</el-button
                  >
                </template>
              </el-table-column>
              <el-table-column label="提示词">
                <template #default="scope">
                  <el-popover
                    placement="top-start"
                    title="绘画提示词"
                    :width="300"
                    trigger="hover"
                    :content="scope.row.prompt"
                  >
                    <template #reference>
                      <span>{{ substr(scope.row.prompt, 20) }}</span>
                    </template>
                  </el-popover>
                </template>
              </el-table-column>
              <el-table-column label="创建时间">
                <template #default="scope">
                  <span>{{ dateFormat(scope.row['created_at']) }}</span>
                </template>
              </el-table-column>
              <el-table-column label="失败原因">
                <template #default="scope">
                  <el-popover
                    placement="top-start"
                    title="失败原因"
                    :width="300"
                    trigger="hover"
                    :content="scope.row.err_msg"
                    v-if="scope.row.progress === 101"
                  >
                    <template #reference>
                      <el-text type="danger">{{ substr(scope.row.err_msg, 20) }}</el-text>
                    </template>
                  </el-popover>
                  <span v-else>无</span>
                </template>
              </el-table-column>
              <el-table-column label="操作" width="180">
                <template #default="scope">
                  <el-button size="small" type="primary" @click="showDetail(scope.row)" plain
                    >详情</el-button
                  >
                  <el-popconfirm
                    title="确定要删除当前记录吗?"
                    @confirm="remove(scope.row, 'image')"
                  >
                    <template #reference>
                      <el-button size="small" type="danger">删除</el-button>
                    </template>
                  </el-popconfirm>
                </template>
              </el-table-column>
            </el-table>
          </el-row>

          <div class="pagination">
            <el-pagination
              v-if="data.image.total > 0"
              background
              layout="total,prev, pager, next"
              :hide-on-single-page="true"
              v-model:current-page="data.image.page"
              v-model:page-size="data.image.pageSize"
              @current-change="fetchImageData()"
              :total="data.image.total"
            />
          </div>
        </div>
        <el-empty v-else />
      </el-tab-pane>
    </el-tabs>

    <el-dialog v-model="showImageDialog" title="图片预览" style="height: 95vh; overflow: auto">
      <el-image
        :src="imgURL"
        :zoom-rate="1.2"
        :max-scale="7"
        :min-scale="0.2"
        :preview-src-list="[imgURL]"
        :initial-index="0"
        fit="cover"
      />
    </el-dialog>

    <!-- 详情弹窗 -->
    <el-dialog
      v-model="showDetailDialog"
      title="任务详情"
      width="600px"
      :close-on-click-modal="false"
    >
      <div class="detail-content" v-if="currentDetail">
        <el-descriptions :column="1" border>
          <el-descriptions-item label="提示词">
            <div class="prompt-text">
              {{ currentDetail.prompt }}
              <el-button
                type="primary"
                size="small"
                @click="copyPrompt(currentDetail.prompt)"
                style="margin-left: 10px"
              >
                复制
              </el-button>
            </div>
          </el-descriptions-item>

          <el-descriptions-item
            label="生成的图片"
            v-if="currentDetail.progress === 100 && currentDetail.img_url"
          >
            <el-image
              :src="getThumbURL(currentDetail.img_url, 200, 200)"
              :preview-src-list="[currentDetail.img_url]"
              fit="cover"
              style="width: 200px; height: 200px"
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
              />
            </div>
          </el-descriptions-item>

          <el-descriptions-item label="生图模型">
            {{ currentDetail.params?.model_name || '-' }}
          </el-descriptions-item>

          <el-descriptions-item label="消耗算力">
            {{ currentDetail.power || 0 }}
          </el-descriptions-item>

          <el-descriptions-item label="图片比例">
            {{ currentDetail.params?.aspect_ratio || '-' }}
          </el-descriptions-item>

          <el-descriptions-item label="图片尺寸">
            {{ currentDetail.params?.size || '-' }}
          </el-descriptions-item>

          <el-descriptions-item label="创建时间">
            {{ formatTime(currentDetail.created_at) }}
          </el-descriptions-item>

          <el-descriptions-item
            label="错误信息"
            v-if="currentDetail.progress === 101 && currentDetail.err_msg"
          >
            <el-text type="danger">{{ currentDetail.err_msg }}</el-text>
          </el-descriptions-item>
        </el-descriptions>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { httpGet, httpPost } from '@/utils/http'
import { dateFormat, getThumbURL, substr } from '@/utils/libs'
import { Search } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { onMounted, ref } from 'vue'

// 变量定义
const data = ref({
  mj: {
    items: [],
    query: { prompt: '', username: '', created_at: [], page: 1, page_size: 15 },
    total: 0,
    page: 1,
    pageSize: 15,
    loading: true,
  },
  image: {
    items: [],
    query: { prompt: '', username: '', created_at: [], page: 1, page_size: 15 },
    total: 0,
    page: 1,
    pageSize: 15,
    loading: true,
  },
})
const activeName = ref('mj')
const taskTypeTheme = {
  image: { text: '绘图', color: '#2185d0' },
  upscale: { text: '放大', color: '#f2711c' },
  variation: { text: '变换', color: '#00b5ad' },
  blend: { text: '融图', color: '#21ba45' },
  swapFace: { text: '换脸', color: '#a333c8' },
}

onMounted(() => {
  fetchMjData()
})

const handleChange = (tab) => {
  switch (tab) {
    case 'mj':
      fetchMjData()
      break
    case 'image':
      fetchImageData()
      break
  }
}

// 搜索对话
const search = (evt, tab) => {
  if (evt.keyCode === 13) {
    handleChange(tab)
  }
}

// 获取数据
const fetchMjData = () => {
  const d = data.value.mj
  d.query.page = d.page
  d.query.page_size = d.pageSize
  httpPost('/api/admin/image/list/mj', d.query)
    .then((res) => {
      if (res.data) {
        d.items = res.data.items
        d.total = res.data.total
        d.page = res.data.page
        d.pageSize = res.data.page_size
      }
      d.loading = false
    })
    .catch((e) => {
      ElMessage.error('获取数据失败：' + e.message)
    })
}

const fetchImageData = () => {
  const d = data.value.image
  d.query.page = d.page
  d.query.page_size = d.pageSize
  httpPost('/api/admin/image/list/image', d.query)
    .then((res) => {
      if (res.data) {
        d.items = res.data.items
        d.total = res.data.total
        d.page = res.data.page
        d.pageSize = res.data.page_size
      }
      d.loading = false
    })
    .catch((e) => {
      ElMessage.error('获取数据失败：' + e.message)
    })
}

const remove = function (row, tab) {
  httpGet(`/api/admin/image/remove?id=${row.id}&tab=${tab}`)
    .then(() => {
      ElMessage.success('删除成功！')
      handleChange(tab)
    })
    .catch((e) => {
      ElMessage.error('删除失败：' + e.message)
    })
}

const showImageDialog = ref(false)
const imgURL = ref('')
const showImage = (url) => {
  showImageDialog.value = true
  imgURL.value = url
}

// 详情弹窗
const showDetailDialog = ref(false)
const currentDetail = ref(null)

const showDetail = (row) => {
  // 解析 params 字段
  let params = {}
  try {
    if (row.params) {
      if (typeof row.params === 'string') {
        params = JSON.parse(row.params)
      } else {
        params = row.params
      }
    }
  } catch (e) {
    console.error('解析 params 失败:', e)
  }

  currentDetail.value = {
    ...row,
    params: params,
  }
  showDetailDialog.value = true
}

// 格式化时间
const formatTime = (timestamp) => {
  if (!timestamp) return '-'
  const date = new Date(timestamp * 1000)
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  const hours = String(date.getHours()).padStart(2, '0')
  const minutes = String(date.getMinutes()).padStart(2, '0')
  const seconds = String(date.getSeconds()).padStart(2, '0')
  return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`
}

// 复制提示词
const copyPrompt = async (text) => {
  if (!text) return
  try {
    if (navigator.clipboard?.writeText) {
      await navigator.clipboard.writeText(text)
      ElMessage.success('复制成功！')
      return
    }
    const textarea = document.createElement('textarea')
    textarea.value = text
    document.body.appendChild(textarea)
    textarea.select()
    document.execCommand('copy')
    document.body.removeChild(textarea)
    ElMessage.success('复制成功！')
  } catch {
    ElMessage.error('复制失败！')
  }
}
</script>

<style lang="scss" scoped>
.image-page {
  .handle-box {
    margin-bottom: 20px;
    .handle-input {
      max-width: 150px;
      margin-right: 10px;
    }
  }

  .opt-box {
    padding-bottom: 10px;
    display: flex;
    justify-content: flex-end;

    .el-icon {
      margin-right: 5px;
    }
  }

  .el-select {
    width: 100%;
  }

  .pagination {
    margin-top: 20px;
    display: flex;
    justify-content: center;
  }
}

.detail-content {
  :deep(.el-descriptions__label) {
    min-width: 110px;
  }

  .prompt-text {
    display: flex;
    align-items: center;
    word-break: break-all;
  }

  .reference-images {
    display: flex;
    flex-wrap: wrap;
  }
}
</style>
