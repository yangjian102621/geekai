<template>
  <div class="container api-key-page">
    <div class="panel-head">
      <div class="panel-title">
        <el-icon><Key /></el-icon>
        <span>令牌管理</span>
      </div>
      <a class="purchase-link" href="https://api.geekai.pro" target="_blank">购买 API-KEY</a>
    </div>

    <div class="toolbar">
      <div class="action-row">
        <button class="pill-btn blue" @click="add">添加令牌</button>
        <button class="pill-btn" @click="copySelected">复制所选令牌</button>
        <button class="pill-btn red" @click="batchRemove">删除所选令牌</button>
      </div>

      <div class="search-row">
        <el-input
          v-model="query.keyword"
          placeholder="搜索关键字"
          class="search-input keyword-input"
          :prefix-icon="Search"
          clearable
          @keyup.enter="search"
        />
        <el-input
          v-model="query.key"
          placeholder="密钥"
          class="search-input key-input"
          :prefix-icon="Search"
          clearable
          @keyup.enter="search"
        />
        <button class="plain-btn" @click="search">查询</button>
        <button class="plain-btn" @click="reset">重置</button>
      </div>
    </div>

    <el-table
      v-if="pagedItems.length > 0"
      :data="pagedItems"
      :row-key="(row) => row.id"
      class="api-key-table"
      @selection-change="handleSelectionChange"
      table-layout="auto"
    >
      <el-table-column type="selection" width="40"></el-table-column>
      <el-table-column prop="name" label="名称" min-width="120" />
      <el-table-column label="API KEY" min-width="280">
        <template #default="scope">
          <span>{{ substr(scope.row.value, 20) }}</span>
          <el-icon class="copy-key" :data-clipboard-text="scope.row.value">
            <DocumentCopy />
          </el-icon>
        </template>
      </el-table-column>
      <el-table-column label="API URL" min-width="240">
        <template #default="scope">
          <span>{{ substr(scope.row.api_url, 30) }}</span>
          <el-icon class="copy-key" :data-clipboard-text="scope.row.api_url">
            <DocumentCopy />
          </el-icon>
        </template>
      </el-table-column>
      <el-table-column label="类型" width="110">
        <template #default="scope">
          {{ getTypeName(scope.row.type) }}
        </template>
      </el-table-column>
      <el-table-column prop="proxy_url" label="代理地址" min-width="140" />

      <el-table-column label="最后使用时间" width="170">
        <template #default="scope">
          <span v-if="scope.row['last_used_at']">{{ scope.row['last_used_at'] }}</span>
          <span v-else class="pill pill-gray">未使用</span>
        </template>
      </el-table-column>
      <el-table-column label="启用状态" width="100">
        <template #default="scope">
          <el-switch v-model="scope.row['enabled']" @change="set('enabled', scope.row)" />
        </template>
      </el-table-column>

      <el-table-column label="操作" width="110" fixed="right">
        <template #default="scope">
          <button class="link-btn primary" @click="edit(scope.row)">编辑</button>
          <el-popconfirm title="确定要删除当前记录吗?" @confirm="remove(scope.row)" :width="200">
            <template #reference>
              <button class="link-btn danger">删除</button>
            </template>
          </el-popconfirm>
        </template>
      </el-table-column>
    </el-table>

    <div class="table-footer" v-if="total > 0">
      <div class="footer-info">显示第 {{ startIndex }} - {{ endIndex }} 条，共 {{ total }} 条</div>
      <div class="footer-nav">
        <span class="footer-text">总页数: {{ totalPages }}</span>
        <el-pagination
          v-model:current-page="page"
          :page-size="pageSize"
          :total="total"
          layout="prev, pager, next"
        />
        <span class="footer-text">每页条数:</span>
        <el-select v-model="pageSize" class="page-size-select" @change="page = 1">
          <el-option :value="10" label="10" />
          <el-option :value="20" label="20" />
          <el-option :value="50" label="50" />
        </el-select>
      </div>
    </div>

    <el-empty :image-size="100" v-if="!loading && pagedItems.length === 0" />

    <el-dialog v-model="showDialog" :close-on-click-modal="false" :title="title">
      <el-form :model="item" label-width="120px" ref="formRef" :rules="rules">
        <el-form-item label="名称：" prop="name">
          <template #default>
            <el-row :gutter="10" class="w-full">
              <el-col :span="12">
                <el-input v-model="item.name" autocomplete="off" />
              </el-col>
              <el-col :span="12">
                <el-select v-model="preset" placeholder="请选择预设" @change="changePreset">
                  <el-option
                    v-for="item in presets"
                    :value="item"
                    :label="item.label"
                    :key="item.value"
                    >{{ item.label }}
                  </el-option>
                </el-select>
              </el-col>
            </el-row>
          </template>
        </el-form-item>
        <el-form-item label="类型：" prop="type">
          <el-select v-model="item.type" placeholder="请选择类型">
            <el-option
              v-for="item in types"
              :value="item.value"
              :label="item.label"
              :key="item.value"
              >{{ item.label }}
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="API KEY：" prop="value">
          <el-input v-model="item.value" type="password" show-password autocomplete="off" />
        </el-form-item>
        <el-form-item label="API URL：" prop="api_url">
          <el-input
            v-model="item.api_url"
            autocomplete="off"
            placeholder="只填 BASE URL 如：https://api.openai.com 或者完整 URL 如：https://api.openai.com/v1/chat/completions"
          />
        </el-form-item>

        <el-form-item label="启用状态：" prop="enable">
          <el-switch v-model="item.enabled" />
        </el-form-item>
      </el-form>

      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showDialog = false">取消</el-button>
          <el-button type="primary" @click="save">提交</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { httpGet, httpPost } from '@/utils/http'
import { dateFormat, substr } from '@/utils/libs'
import { DocumentCopy, Key, Search } from '@element-plus/icons-vue'
import ClipboardJS from 'clipboard'
import { ElMessage, ElMessageBox } from 'element-plus'
import { computed, onMounted, onUnmounted, reactive, ref } from 'vue'

// 变量定义
const items = ref([])
const query = ref({ keyword: '', key: '' })
const item = ref({})
const showDialog = ref(false)
const rules = reactive({
  name: [{ required: true, message: '请输入名称', trigger: 'change' }],
  type: [{ required: true, message: '请选择用途', trigger: 'change' }],
  value: [{ required: true, message: '请输入 API KEY 值', trigger: 'change' }],
})

const loading = ref(true)
const formRef = ref(null)
const title = ref('')
const types = ref([
  { label: '对话', value: 'chat' },
  { label: 'Midjourney', value: 'mj' },
  { label: 'AI图像生成', value: 'image' },
  { label: 'Suno文生歌', value: 'suno' },
  { label: 'Luma视频', value: 'luma' },
  { label: '可灵视频', value: 'keling' },
  { label: 'Realtime API', value: 'realtime' },
  { label: '语音合成', value: 'tts' },
  { label: '其他', value: 'other' },
])
const clipboard = ref(null)
const presets = ref([
  { label: 'GiteeAI', value: 'https://ai.gitee.com' },
  { label: '中转01', value: 'https://api.geekai.pro' },
  { label: '中转03', value: 'https://api.geekai.me' },
  { label: 'OpenAI', value: 'https://api.openai.com' },
])
const preset = ref(null)

const page = ref(1)
const pageSize = ref(10)
const selectedIds = ref([])

const total = computed(() => items.value.length)
const totalPages = computed(() => Math.max(1, Math.ceil(total.value / (pageSize.value || 10))))
const startIndex = computed(() => (total.value ? (page.value - 1) * pageSize.value + 1 : 0))
const endIndex = computed(() => Math.min(page.value * pageSize.value, total.value))
const pagedItems = computed(() =>
  items.value.slice((page.value - 1) * pageSize.value, page.value * pageSize.value)
)

onMounted(() => {
  clipboard.value = new ClipboardJS('.copy-key')
  clipboard.value.on('success', () => {
    ElMessage.success('复制成功！')
  })

  clipboard.value.on('error', () => {
    ElMessage.error('复制失败！')
  })

  fetchData()
})

onUnmounted(() => {
  clipboard.value.destroy()
})

const getTypeName = (type) => {
  for (let v of types.value) {
    if (v.value === type) {
      return v.label
    }
  }
  return ''
}

// 获取数据
const fetchData = () => {
  loading.value = true
  httpGet('/api/admin/apikey/list', query.value)
    .then((res) => {
      if (res.data) {
        // 初始化数据
        const arr = res.data
        for (let i = 0; i < arr.length; i++) {
          arr[i].last_used_at = dateFormat(arr[i].last_used_at)
        }
        items.value = arr
      } else {
        items.value = []
      }
      selectedIds.value = []
      loading.value = false
    })
    .catch(() => {
      loading.value = false
      ElMessage.error('获取数据失败')
    })
}

const search = () => {
  page.value = 1
  fetchData()
}

const reset = () => {
  query.value.keyword = ''
  query.value.key = ''
  page.value = 1
  fetchData()
}

const handleSelectionChange = (rows) => {
  selectedIds.value = rows.map((row) => row.id)
}

// 复制所选令牌
const copySelected = () => {
  const codes = items.value
    .filter((it) => selectedIds.value.includes(it.id))
    .map((it) => it.value)
    .join('\n')
  if (!codes) {
    ElMessage.warning('请先选择要复制的令牌')
    return
  }
  const done = () => ElMessage.success(`已复制 ${selectedIds.value.length} 个令牌`)
  if (navigator.clipboard && window.isSecureContext) {
    navigator.clipboard
      .writeText(codes)
      .then(done)
      .catch(() => fallbackCopy(codes, done))
  } else {
    fallbackCopy(codes, done)
  }
}

const fallbackCopy = (text, done) => {
  const textarea = document.createElement('textarea')
  textarea.value = text
  textarea.style.position = 'fixed'
  textarea.style.opacity = '0'
  document.body.appendChild(textarea)
  textarea.select()
  document.execCommand('copy')
  document.body.removeChild(textarea)
  done()
}

// 批量删除所选令牌
const batchRemove = () => {
  if (!selectedIds.value.length) {
    ElMessage.warning('请先选择要删除的令牌')
    return
  }
  ElMessageBox.confirm(`确定要删除所选的 ${selectedIds.value.length} 个令牌吗?`, '删除提示', {
    confirmButtonText: '确认',
    cancelButtonText: '取消',
    type: 'warning',
  }).then(() => {
    httpPost('/api/admin/apikey/batchRemove', { ids: selectedIds.value })
      .then((res) => {
        ElMessage.success(`已删除 ${res.data.deleted || 0} 个令牌`)
        page.value = 1
        fetchData()
      })
      .catch((e) => {
        ElMessage.error('删除失败：' + e.message)
      })
  })
}

const add = function () {
  showDialog.value = true
  title.value = '新增 API KEY'
  preset.value = null
  item.value = {
    enabled: true,
    api_url: '',
  }
}

const edit = function (row) {
  showDialog.value = true
  title.value = '修改 API KEY'
  item.value = row
}

const save = function () {
  formRef.value.validate((valid) => {
    if (valid) {
      showDialog.value = false
      httpPost('/api/admin/apikey/save', item.value)
        .then(() => {
          ElMessage.success('操作成功！')
          fetchData()
        })
        .catch((e) => {
          ElMessage.error('操作失败，' + e.message)
        })
    } else {
      return false
    }
  })
}

const remove = function (row) {
  httpGet('/api/admin/apikey/remove?id=' + row.id)
    .then(() => {
      ElMessage.success('删除成功！')
      fetchData()
    })
    .catch((e) => {
      ElMessage.error('删除失败：' + e.message)
    })
}

const set = (filed, row) => {
  httpPost('/api/admin/apikey/set', { id: row.id, filed: filed, value: row[filed] })
    .then(() => {
      ElMessage.success('操作成功！')
    })
    .catch((e) => {
      ElMessage.error('操作失败：' + e.message)
    })
}

const changePreset = (row) => {
  item.value.api_url = row.value
  if (!item.value.name) {
    item.value.name = row.label
  }
}
</script>

<style lang="scss" scoped>
.api-key-page {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'PingFang SC', 'Microsoft YaHei', sans-serif;
  padding: 0 !important;

  .panel-head {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 16px 20px;
    border-bottom: 1px solid #f0f2f5;
  }

  .panel-title {
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 17px;
    font-weight: 600;
    color: #1d2129;

    .el-icon {
      font-size: 18px;
      color: #1d2129;
    }
  }

  .purchase-link {
    font-size: 13px;
    color: #667085;
    text-decoration: none;

    &:hover {
      color: #1677ff;
    }
  }

  .toolbar {
    padding: 14px 20px;
    border-bottom: 1px solid #f0f2f5;
  }

  .action-row {
    display: flex;
    align-items: center;
    gap: 10px;
    margin-bottom: 12px;
  }

  .pill-btn {
    height: 36px;
    padding: 0 16px;
    border: none;
    border-radius: 6px;
    background: #f2f4f7;
    color: #1d2129;
    font-size: 14px;
    cursor: pointer;
    white-space: nowrap;
    transition: background 0.2s ease;

    &:hover {
      background: #e9ebf0;
    }

    &.blue {
      color: #1677ff;
    }

    &.red {
      color: #ef4444;
    }
  }

  .search-row {
    display: flex;
    align-items: center;
    gap: 10px;
    flex-wrap: wrap;
  }

  .search-input {
    :deep(.el-input__wrapper) {
      height: 36px;
      background: #f2f4f7;
      box-shadow: none;
      border: 1px solid transparent;
      border-radius: 6px;
      padding: 0 12px;
    }

    :deep(.el-input__wrapper:hover) {
      border-color: #e2e6ee;
    }

    :deep(.el-input__wrapper.is-focus) {
      background: #fff;
      border-color: #1677ff;
      box-shadow: 0 0 0 2px rgba(22, 119, 255, 0.08);
    }

    :deep(.el-input__inner) {
      font-size: 14px;
      color: #1d2129;
    }

    :deep(.el-input__inner::placeholder) {
      color: #98a2b3;
    }

    :deep(.el-input__prefix) {
      color: #98a2b3;
      font-size: 14px;
    }

    &.keyword-input {
      width: 340px;
    }

    &.key-input {
      width: 300px;
    }
  }

  .plain-btn {
    height: 36px;
    padding: 0 16px;
    border: none;
    border-radius: 6px;
    background: #f2f4f7;
    color: #1d2129;
    font-size: 14px;
    cursor: pointer;
    transition: background 0.2s ease;

    &:hover {
      background: #e9ebf0;
    }
  }

  .api-key-table {
    width: calc(100% - 40px);
    margin: 0 20px;

    :deep(.el-table__inner-wrapper::before) {
      display: none;
    }

    :deep(th.el-table__cell) {
      background: transparent;
      color: #8a8f99;
      font-size: 14px;
      font-weight: 400;
      border-bottom: 1px solid #ebeef5;
      padding: 12px 0;
    }

    :deep(td.el-table__cell) {
      border-bottom: 1px solid #ebeef5;
      font-size: 14px;
      color: #303133;
      padding: 13px 0;
    }

    :deep(.el-table__body tr:hover > td.el-table__cell) {
      background: #fafbfc;
    }
  }

  .copy-key {
    margin-left: 5px;
    cursor: pointer;
    color: #667085;
    font-size: 14px;
    position: relative;
    top: 2px;
  }

  .pill {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    height: 24px;
    padding: 0 10px;
    border-radius: 999px;
    font-size: 12px;
    line-height: 1;
    white-space: nowrap;

    &.pill-gray {
      background: #f2f4f7;
      color: #667085;
    }
  }

  .link-btn {
    border: none;
    background: transparent;
    padding: 0;
    margin-right: 12px;
    color: #303133;
    font-size: 14px;
    cursor: pointer;

    &:hover {
      color: #1677ff;
    }

    &.primary {
      color: #1677ff;

      &:hover {
        color: #4096ff;
      }
    }

    &.danger {
      color: #ef4444;

      &:hover {
        color: #dc2626;
      }
    }
  }

  .table-footer {
    display: flex;
    justify-content: space-between;
    align-items: center;
    flex-wrap: wrap;
    gap: 12px;
    padding: 16px 20px;
    font-size: 13px;
    color: #667085;
  }

  .footer-nav {
    display: flex;
    align-items: center;
    gap: 14px;
  }

  .footer-text {
    color: #667085;
  }

  :deep(.el-pagination) {
    padding: 0;
    --el-pagination-button-bg-color: transparent;
  }

  :deep(.el-pager li) {
    min-width: 30px;
    height: 30px;
    line-height: 30px;
    border-radius: 6px;
    background: transparent;
    color: #667085;
    font-weight: 500;
  }

  :deep(.el-pager li.is-active) {
    background: #e8f2ff;
    color: #1677ff;
  }

  :deep(.btn-prev),
  :deep(.btn-next) {
    background: transparent;
    color: #667085;
  }

  .page-size-select {
    width: 90px;
  }
}

:global([data-theme='dark'] .api-key-page) .panel-head,
:global([data-theme='dark'] .api-key-page) .toolbar {
  border-bottom-color: rgba(255, 255, 255, 0.08);
}

:global([data-theme='dark'] .api-key-page) :deep(th.el-table__cell) {
  background: transparent;
  color: rgba(255, 255, 255, 0.55);
  border-bottom-color: rgba(255, 255, 255, 0.08);
}

:global([data-theme='dark'] .api-key-page) :deep(td.el-table__cell) {
  border-bottom-color: rgba(255, 255, 255, 0.08);
  color: rgba(255, 255, 255, 0.85);
}

:global([data-theme='dark'] .api-key-page) .search-input :deep(.el-input__wrapper),
:global([data-theme='dark'] .api-key-page) .plain-btn,
:global([data-theme='dark'] .api-key-page) .pill-btn {
  background: rgba(255, 255, 255, 0.06);
}
</style>
