<template>
  <div class="container suno-page">
    <div class="handle-box">
      <el-input
        v-model="query.title"
        placeholder="标题"
        class="handle-input mr10"
        @keyup="search"
        clearable
      />
      <el-input
        v-model="query.prompt"
        placeholder="提示词"
        class="handle-input mr10"
        @keyup="search"
        clearable
      />
      <el-date-picker
        v-model="query.created_at"
        type="daterange"
        start-placeholder="开始日期"
        end-placeholder="结束日期"
        format="YYYY-MM-DD"
        value-format="YYYY-MM-DD"
       
      />
      <el-button type="primary" class="ml-2" :icon="Search" @click="fetchData">搜索</el-button>
    </div>

    <div v-if="items.length > 0" v-loading="loading">
      <el-row>
        <el-table :data="items" :row-key="(row) => row.id" table-layout="auto">
          <el-table-column prop="user_id" label="用户ID" width="100" />
          <el-table-column prop="title" label="标题" width="120">
            <template #default="{ row }">
              <span class="text-ellipsis">{{ substr(row.title, 20) }}</span>
            </template>
          </el-table-column>
          <el-table-column label="预览" width="180">
            <template #default="scope">
              <div class="container" v-if="scope.row.cover_url">
                <el-image :src="scope.row.cover_url" fit="cover" />
                <div class="duration">
                  {{ formatTime(scope.row.duration) }}
                </div>
                <button class="play flex justify-center items-center" @click="playMusic(scope.row)">
                  <img src="/images/play.svg" alt="" />
                </button>
              </div>
              <el-image v-else-if="scope.row.progress === 101" src="/images/failed.jpg" style="height: 90px" fit="cover" />
              <div class="flex flex-col items-center justify-center h-[100px]" v-else>
                <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-purple-600 mx-auto"></div>
                <span class="text-xs text-purple-600 mt-2 block">生成中...</span>
              </div>
            </template>
          </el-table-column>
          <el-table-column prop="progress" label="任务进度" width="100">
            <template #default="scope">
              <span v-if="scope.row.progress <= 100">{{ scope.row.progress }}%</span>
              <el-tag v-else type="danger">已失败</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="power" label="消耗算力" width="100" />
          <el-table-column prop="play_times" label="播放次数" width="100" />
          <el-table-column label="歌词" width="110">
            <template #default="scope">
              <el-button size="small" type="primary" plain @click="showLyric(scope.row)"
                >查看歌词</el-button
              >
            </template>
          </el-table-column>
          <el-table-column label="创建时间" width="160">
            <template #default="scope">
              <span>{{ dateFormat(scope.row['created_at']) }}</span>
            </template>
          </el-table-column>
          <el-table-column label="失败原因" width="180">
            <template #default="scope">
              <el-popover
                v-if="scope.row.progress === 101"
                placement="top-start"
                title="失败原因"
                :width="300"
                trigger="hover"
                :content="scope.row.err_msg"
              >
                <template #reference>
                  <el-text type="danger">{{ substr(scope.row.err_msg, 20) }}</el-text>
                </template>
              </el-popover>
              <span v-else>无</span>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="150">
            <template #default="scope">
              <el-button size="small" @click="showDetailDialog(scope.row)">详情</el-button>
              <el-popconfirm title="确定要删除当前记录吗?" @confirm="remove(scope.row)">
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
          v-if="total > 0"
          background
          layout="total,prev, pager, next"
          :hide-on-single-page="true"
          v-model:current-page="page"
          v-model:page-size="pageSize"
          @current-change="fetchData"
          :total="total"
        />
      </div>
    </div>
    <el-empty v-else />

    <div class="music-player" v-if="showPlayer">
      <music-player
        :songs="playList"
        ref="playerRef"
        :show-close="true"
        @close="showPlayer = false"
      />
    </div>

    <el-dialog v-model="showLyricDialog" title="歌词">
      <div class="chat-line" v-html="lyrics"></div>
    </el-dialog>

    <!-- 详情对话框 -->
    <el-dialog v-model="detailDialogVisible" title="任务详情" width="800px">
      <div v-if="currentDetail" class="space-y-4">
        <div class="detail-grid">
          <div class="detail-item">
            <div class="detail-label">任务ID</div>
            <div class="detail-value">{{ currentDetail.task_id }}</div>
          </div>
          <div class="detail-item">
            <div class="detail-label">用户ID</div>
            <div class="detail-value">{{ currentDetail.user_id }}</div>
          </div>
          <div class="detail-item">
            <div class="detail-label">风格</div>
            <div class="detail-value">{{ currentDetail.params.tags }}</div>
          </div>
          <div class="detail-item">
            <div class="detail-label">任务进度</div>
            <div class="detail-value">
              <span v-if="currentDetail.progress <= 100">{{ currentDetail.progress }}%</span>
              <el-tag v-else type="danger">已失败</el-tag>
            </div>
          </div>
          <div class="detail-item">
            <div class="detail-label">算力消耗</div>
            <div class="detail-value">{{ currentDetail.power }}</div>
          </div>
          <div class="detail-item">
            <div class="detail-label">播放次数</div>
            <div class="detail-value">{{ currentDetail.play_times }}</div>
          </div>
          <div class="detail-item">
            <div class="detail-label">时长</div>
            <div class="detail-value">{{ formatTime(currentDetail.duration) }}</div>
          </div>
          <div class="detail-item">
            <div class="detail-label">创建时间</div>
            <div class="detail-value">{{ dateFormat(currentDetail.created_at) }}</div>
          </div>
        </div>

        <div class="detail-full">
          <div class="detail-label">歌词</div>
          <div class="detail-value" v-html="md.render(currentDetail.prompt)" />
        </div>

        <div v-if="currentDetail.progress === 101 && currentDetail.err_msg" class="detail-full error">
          <div class="detail-label">错误信息</div>
          <div class="detail-value">{{ currentDetail.err_msg }}</div>
        </div>

        <!-- 音乐预览区域 - 仅成功任务显示 -->
        <div v-if="currentDetail.progress === 100 && currentDetail.cover_url" class="detail-preview">
          <div class="detail-label mb-2">音乐预览</div>
          <div class="preview-container">
            <el-image :src="currentDetail.cover_url" fit="cover" class="preview-image" />
            <div class="preview-info">
              <div class="preview-duration">{{ formatTime(currentDetail.duration) }}</div>
              <el-button type="primary" @click="playMusic(currentDetail)">播放</el-button>
            </div>
          </div>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import MusicPlayer from '@/components/MusicPlayer.vue'
import { httpGet, httpPost } from '@/utils/http'
import { dateFormat, formatTime, substr } from '@/utils/libs'
import { Search } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import MarkdownIt from 'markdown-it'
import { nextTick, onMounted, ref } from 'vue'

const items = ref([])
const query = ref({ prompt: '', title: '', created_at: [], page: 1, page_size: 10 })
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)
const loading = ref(true)

const playList = ref([])
const playerRef = ref(null)
const showPlayer = ref(false)
const showLyricDialog = ref(false)
const lyrics = ref('')
const detailDialogVisible = ref(false)
const currentDetail = ref(null)

onMounted(() => {
  fetchData()
})

// 搜索
const search = (evt) => {
  if (evt.keyCode === 13) {
    fetchData()
  }
}

// 获取数据
const fetchData = () => {
  query.value.page = page.value
  query.value.page_size = pageSize.value
  httpPost('/api/admin/suno/list', query.value)
    .then((res) => {
      if (res.data) {
        items.value = res.data.items
        total.value = res.data.total
        page.value = res.data.page
        pageSize.value = res.data.page_size
      }
      loading.value = false
    })
    .catch((e) => {
      ElMessage.error('获取数据失败：' + e.message)
    })
}

const remove = function (row) {
  httpGet(`/api/admin/suno/remove?id=${row.id}`)
    .then(() => {
      ElMessage.success('删除成功！')
      fetchData()
    })
    .catch((e) => {
      ElMessage.error('删除失败：' + e.message)
    })
}

const playMusic = (item) => {
  playList.value = [item]
  showPlayer.value = true
  nextTick(() => playerRef.value.play())
}

const md = MarkdownIt({
  breaks: true,
  html: true,
  linkify: true,
})

const showLyric = (item) => {
  showLyricDialog.value = true
  lyrics.value = md.render(item.prompt)
}

const showDetailDialog = (item) => {
  currentDetail.value = item
  detailDialogVisible.value = true
}
</script>

<style lang="scss" scoped>
.suno-page {
  .handle-box {
    margin-bottom: 20px;
    .handle-input {
      max-width: 150px;
      margin-right: 10px;
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

  :deep(.el-table__header) {
    th {
      white-space: nowrap;
    }
  }

  .container {
    width: 160px;
    position: relative;

    .el-image {
      width: 160px;
      height: 90px;
      border-radius: 5px;
    }

    .duration {
      position: absolute;
      bottom: 6px;
      right: 0;
      background-color: rgba(255, 255, 255, 0.7);
      padding: 0 3px;
      font-family: 'Input Sans';
      font-size: 14px;
      font-weight: 700;
      border-radius: 0.125rem;
    }

    .play {
      position: absolute;
      width: 100%;
      height: 100%;
      top: 0;
      left: 50%;
      border: none;
      border-radius: 5px;
      background: rgba(100, 100, 100, 0.3);
      cursor: pointer;
      color: #ffffff;
      opacity: 0;
      transform: translate(-50%, 0px);
      transition: opacity 0.3s ease 0s;
    }

    &:hover {
      .play {
        opacity: 1;
      }
    }
  }

  .music-player {
    position: absolute;
    bottom: 20px;
    z-index: 99999;
    width: 100%;
  }

  .text-ellipsis {
    display: block;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .space-y-4 > * + * {
    margin-top: 1rem;
  }

  .detail-grid {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: 1rem;
  }

  .detail-item {
    padding: 0.75rem;
    background-color: #f9fafb;
    border-radius: 0.375rem;
  }

  .detail-full {
    padding: 0.75rem;
    background-color: #f9fafb;
    border-radius: 0.375rem;

    &.error {
      background-color: #fef2f2;

      .detail-label {
        color: #dc2626;
      }

      .detail-value {
        color: #991b1b;
      }
    }
  }

  .detail-label {
    font-size: 0.75rem;
    color: #6b7280;
    margin-bottom: 0.25rem;
  }

  .detail-value {
    font-size: 0.875rem;
    font-weight: 500;
  }

  .detail-preview {
    .mb-2 {
      margin-bottom: 0.5rem;
    }

    .preview-container {
      display: flex;
      gap: 1rem;
      align-items: center;

      .preview-image {
        width: 160px;
        height: 90px;
        border-radius: 0.5rem;
      }

      .preview-info {
        flex: 1;
        display: flex;
        align-items: center;
        gap: 1rem;

        .preview-duration {
          font-size: 0.875rem;
          color: #6b7280;
        }
      }
    }
  }
}
</style>
