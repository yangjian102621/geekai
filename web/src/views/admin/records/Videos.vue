<template>
  <div class="min-h-screen bg-gray-50 p-6">
    <!-- 搜索筛选区域 -->
    <div class="bg-white rounded-lg shadow-sm p-4 mb-6">
      <div class="flex items-center gap-3 flex-wrap">
        <el-input
          v-model="query.prompt"
          placeholder="提示词"
          style="width: 200px"
          clearable
          @keyup="search"
        />
        <el-select v-model="query.type" placeholder="视频类型" style="width: 120px" clearable>
          <el-option v-for="provider in orderedProviders" :key="provider" :label="getProviderName(provider)" :value="provider" />
        </el-select>
        <el-select v-model="query.status" placeholder="任务状态" style="width: 120px" clearable>
          <el-option label="全部" value="" />
          <el-option label="等待中" value="pending" />
          <el-option label="进行中" value="in_progress" />
          <el-option label="下载中" value="downloading" />
          <el-option label="成功" value="success" />
          <el-option label="失败" value="failed" />
        </el-select>
        <el-date-picker
          v-model="query.created_at"
          type="daterange"
          start-placeholder="开始日期"
          end-placeholder="结束日期"
          format="YYYY-MM-DD"
          value-format="YYYY-MM-DD"
          style="width: 240px"
        />
        <el-button type="primary" @click="fetchData">搜索</el-button>
      </div>
    </div>

    <!-- 视频列表 - 卡片式布局 -->
    <div v-if="items.length > 0" v-loading="loading" class="space-y-4">
      <div
        v-for="item in items"
        :key="item.id"
        class="bg-white rounded-lg shadow-sm p-4 hover:shadow-md transition-shadow cursor-pointer"
      >
        <div class="flex gap-4">
          <!-- 视频预览 -->
          <div class="w-40 h-24 flex-shrink-0 bg-gray-100 rounded-lg overflow-hidden relative">
            <video
              v-if="item.status === 'success'"
              :src="replaceImg(item.video_url)"
              class="w-full h-full object-cover"
              muted
              loop
              @click="playVideo(item)"
            />
            <div
              v-else-if="item.status === 'downloading'"
              class="flex items-center justify-center h-full"
            >
              <div class="text-center">
                <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-purple-600 mx-auto"></div>
                <span class="text-xs text-purple-600 mt-2 block">下载中...</span>
              </div>
            </div>
            <div
              v-else-if="item.status === 'in_progress'"
              class="flex items-center justify-center h-full"
            >
              <div class="text-center">
                <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 mx-auto"></div>
                <span class="text-xs text-blue-600 mt-2 block">{{ item.progress }}%</span>
              </div>
            </div>
            <div
              v-else-if="item.status === 'failed'"
              class="flex items-center justify-center h-full bg-red-50"
            >
              <i class="text-red-500 text-2xl">✕</i>
            </div>
            <div
              v-else
              class="flex items-center justify-center h-full"
            >
              <i class="text-gray-400 text-2xl">⏱</i>
            </div>
            
            <!-- 播放按钮 -->
            <button
              v-if="item.status === 'success'"
              @click="playVideo(item)"
              class="absolute inset-0 flex items-center justify-center bg-black bg-opacity-50 opacity-0 hover:opacity-100 transition-opacity"
            >
              <svg class="w-12 h-12 text-white" fill="currentColor" viewBox="0 0 20 20">
                <path d="M6.3 2.841A1.5 1.5 0 004 4.11V15.89a1.5 1.5 0 002.3 1.269l9.344-5.89a1.5 1.5 0 000-2.538L6.3 2.84z" />
              </svg>
            </button>
          </div>

          <!-- 任务信息 -->
          <div class="flex-1 min-w-0">
            <div class="flex items-start justify-between mb-2">
              <div class="flex items-center gap-2 flex-wrap">
                <span :class="getStatusClass(item.status)" class="px-2 py-1 text-xs rounded-full font-medium">
                  {{ getStatusText(item.status) }}
                </span>
                <span class="text-sm text-gray-600">用户ID: {{ item.user_id }}</span>
                <span class="text-xs text-gray-500 px-2 py-1 bg-gray-100 rounded">{{ item.type }}</span>
              </div>
              <span class="text-xs text-gray-500 whitespace-nowrap">{{ dateFormat(item.created_at) }}</span>
            </div>

            <p class="text-sm text-gray-700 mb-2 line-clamp-2">{{ item.prompt }}</p>

            <div class="flex items-center gap-4 text-xs text-gray-600">
              <span>渠道: {{ item.channel }}</span>
              <span>算力: {{ item.power }}</span>
              <span v-if="item.status === 'in_progress'">进度: {{ item.progress }}%</span>
              <span v-if="item.err_msg" class="text-red-600">错误: {{ substr(item.err_msg, 30) }}</span>
            </div>
          </div>

          <!-- 操作按钮 -->
          <div class="flex flex-col gap-2">
            <span>
              <el-button size="small" @click="showDetailDialog(item)">详情</el-button>
            </span>
            <el-popconfirm title="确定要删除当前记录吗?" @confirm="remove(item)">
              <template #reference>
                <el-button size="small" type="danger">删除</el-button>
              </template>
            </el-popconfirm>
          </div>
        </div>
      </div>
    </div>
    <el-empty v-else description="暂无数据" />

    <!-- 分页 -->
    <div v-if="total > 0" class="mt-6 flex justify-center">
      <el-pagination
        background
        layout="total, prev, pager, next"
        v-model:current-page="page"
        v-model:page-size="pageSize"
        :total="total"
        @current-change="fetchData"
      />
    </div>

    <!-- 视频预览对话框 -->
    <el-dialog v-model="showVideoDialog" title="视频预览" width="800px" @close="handleVideoDialogClose">
      <video
        ref="videoPlayerRef"
        style="width: 100%; max-height: 70vh"
        :src="currentVideoUrl"
        controls
        autoplay
        loop
      >
        您的浏览器不支持视频播放
      </video>
    </el-dialog>

    <!-- 详情对话框 -->
    <el-dialog v-model="detailDialogVisible" title="任务详情" width="800px" @close="handleDetailDialogClose">
      <div v-if="currentDetail" class="space-y-4">
        <div class="grid grid-cols-2 gap-4">
          <div class="p-3 bg-gray-50 rounded">
            <div class="text-xs text-gray-500 mb-1">任务ID</div>
            <div class="text-sm font-medium">{{ currentDetail.task_id }}</div>
          </div>
          <div class="p-3 bg-gray-50 rounded">
            <div class="text-xs text-gray-500 mb-1">状态</div>
            <span :class="getStatusClass(currentDetail.status)" class="px-2 py-1 text-xs rounded-full">
              {{ getStatusText(currentDetail.status) }}
            </span>
          </div>
          <div class="p-3 bg-gray-50 rounded">
            <div class="text-xs text-gray-500 mb-1">渠道</div>
            <div class="text-sm font-medium">{{ currentDetail.channel }}</div>
          </div>
          <div class="p-3 bg-gray-50 rounded">
            <div class="text-xs text-gray-500 mb-1">类型</div>
            <div class="text-sm font-medium">{{ currentDetail.type }}</div>
          </div>
          <div class="p-3 bg-gray-50 rounded">
            <div class="text-xs text-gray-500 mb-1">算力消耗</div>
            <div class="text-sm font-medium">{{ currentDetail.power }}</div>
          </div>
          <div class="p-3 bg-gray-50 rounded">
            <div class="text-xs text-gray-500 mb-1">创建时间</div>
            <div class="text-sm font-medium">{{ dateFormat(currentDetail.created_at) }}</div>
          </div>
        </div>

        <div class="p-3 bg-gray-50 rounded">
          <div class="text-xs text-gray-500 mb-2">提示词</div>
          <div class="text-sm">{{ currentDetail.prompt }}</div>
        </div>

        <div v-if="currentDetail.params" class="p-3 bg-gray-50 rounded">
          <div class="text-xs text-gray-500 mb-2">任务参数</div>
          <pre class="text-xs overflow-x-auto">{{ JSON.stringify(currentDetail.params, null, 2) }}</pre>
        </div>

        <div v-if="currentDetail.output" class="p-3 bg-gray-50 rounded">
          <div class="text-xs text-gray-500 mb-2">输出信息</div>
          <pre class="text-xs overflow-x-auto">{{ JSON.stringify(currentDetail.output, null, 2) }}</pre>
        </div>

        <div v-if="currentDetail.err_msg" class="p-3 bg-red-50 rounded">
          <div class="text-xs text-red-500 mb-2">错误信息</div>
          <div class="text-sm text-red-700">{{ currentDetail.err_msg }}</div>
        </div>

          <!-- 视频预览区域 - 仅成功任务显示 -->
        <div v-if="currentDetail.status === 'success' && currentDetail.video_url" class="rounded-lg overflow-hidden">
          <div class="text-base text-blue-500 mb-2">视频预览</div>
          <video
            ref="detailVideoPlayerRef"
            :src="replaceImg(currentDetail.video_url)"
            class="w-full"
            style="max-height: 400px"
            controls
            loop
          >
            您的浏览器不支持视频播放
          </video>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { httpGet, httpPost } from '@/utils/http'
import { dateFormat, replaceImg, substr } from '@/utils/libs'
import { ElMessage } from 'element-plus'
import { onMounted, ref } from 'vue'
import { getVideoProviders } from '@/store/data/video_params'
import { computed } from 'vue'
import { getProviderName } from '@/store/data/video_params'

const items = ref([])
const query = ref({
  prompt: '',
  type: '',
  status: '',
  created_at: [],
  page: 1,
  page_size: 10,
})
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)
const loading = ref(false)
const showVideoDialog = ref(false)
const currentVideoUrl = ref('')
const videoPlayerRef = ref(null)
const detailDialogVisible = ref(false)
const detailVideoPlayerRef = ref(null)
const currentDetail = ref(null)
const providers = getVideoProviders()
const orderedProviders = computed(() => {
  return ['sora', 'veo', 'luma', 'keling', 'minimax', 'wan'].filter((p) => providers.includes(p))
})


// 状态配置
const statusConfig = {
  pending: { label: '等待中', class: 'bg-gray-100 text-gray-700' },
  in_progress: { label: '进行中', class: 'bg-blue-100 text-blue-700' },
  downloading: { label: '下载中', class: 'bg-purple-100 text-purple-700' },
  success: { label: '成功', class: 'bg-green-100 text-green-700' },
  failed: { label: '失败', class: 'bg-red-100 text-red-700' },
}

const getStatusText = (status) => {
  return statusConfig[status]?.label || status
}

const getStatusClass = (status) => {
  return statusConfig[status]?.class || 'bg-gray-100 text-gray-700'
}

onMounted(() => {
  fetchData()
})

const search = (evt) => {
  if (evt.keyCode === 13) {
    fetchData()
  }
}

const fetchData = () => {
  loading.value = true
  query.value.page = page.value
  query.value.page_size = pageSize.value
  httpPost('/api/admin/video/list', query.value)
    .then((res) => {
      if (res.data) {
        items.value = res.data.items || []
        total.value = res.data.total
        page.value = res.data.page
        pageSize.value = res.data.page_size
      }
    })
    .catch((e) => {
      ElMessage.error('获取数据失败：' + e.message)
    })
    .finally(() => {
      loading.value = false
    })
}

const remove = (item) => {
  httpGet(`/api/admin/video/remove?id=${item.id}&tab=${item.type}`)
    .then(() => {
      ElMessage.success('删除成功！')
      fetchData()
    })
    .catch((e) => {
      ElMessage.error('删除失败：' + e.message)
    })
}

const playVideo = (item) => {
  currentVideoUrl.value = replaceImg(item.video_url)
  showVideoDialog.value = true
}

// 处理视频对话框关闭事件
const handleVideoDialogClose = () => {
  if (videoPlayerRef.value) {
    videoPlayerRef.value.pause()
    videoPlayerRef.value.currentTime = 0
  }
  showVideoDialog.value = false
}

// 处理详情对话框关闭事件
const handleDetailDialogClose = () => {
  if (detailVideoPlayerRef.value) {
    detailVideoPlayerRef.value.pause()
    detailVideoPlayerRef.value.currentTime = 0
  }
  detailDialogVisible.value = false
}

const showDetailDialog = (item) => {
  currentDetail.value = item
  detailDialogVisible.value = true
}
</script>

<style scoped>
.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>
