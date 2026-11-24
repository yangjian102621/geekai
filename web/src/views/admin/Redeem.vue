<template>
  <div class="container redeem-page">
    <div class="panel-head">
      <div class="panel-title">
        <el-icon><Ticket /></el-icon>
        <span>兑换码管理</span>
      </div>
      <button class="compact-btn" :class="{ active: compact }" @click="compact = !compact">紧凑列表</button>
    </div>

    <div class="toolbar">
      <div class="toolbar-left">
        <button class="text-action blue" @click="add">添加兑换码</button>
        <button class="text-action" @click="copySelected">复制所选兑换码到剪贴板</button>
        <button class="text-action red" @click="clearInvalid">清除失效兑换码</button>
      </div>
      <div class="toolbar-right">
        <el-input
          v-model="query.code"
          placeholder="关键字(id或者名称)"
          class="search-input"
          :prefix-icon="Search"
          clearable
          @keyup.enter="search"
        />
        <button class="plain-btn" @click="search">查询</button>
        <button class="plain-btn" @click="reset">重置</button>
      </div>
    </div>

    <el-table
      v-if="items.length > 0"
      :data="items"
      :row-key="(row) => row.id"
      :class="['redeem-table', { compact }]"
      @selection-change="handleSelectionChange"
      table-layout="auto"
    >
      <el-table-column type="selection" width="40"></el-table-column>
      <el-table-column prop="name" label="名称" min-width="120" />
      <el-table-column label="兑换码" min-width="280">
        <template #default="scope">
          <span>{{ substr(scope.row.code, 24) }}</span>
          <el-icon class="copy-code" :data-clipboard-text="scope.row.code">
            <DocumentCopy />
          </el-icon>
        </template>
      </el-table-column>
      <el-table-column label="兑换人" min-width="130">
        <template #default="scope">
          <span v-if="scope.row['username'] !== ''">{{ scope.row['username'] }}</span>
          <span v-else class="pill pill-gray">未兑换</span>
        </template>
      </el-table-column>
      <el-table-column prop="power" label="积分" width="90" />

      <el-table-column label="生成时间" width="170">
        <template #default="scope">
          <span>{{ dateFormat(scope.row['created_at']) }}</span>
        </template>
      </el-table-column>

      <el-table-column label="兑换时间" width="170">
        <template #default="scope">
          <span v-if="scope.row['redeemed_at'] > 0">{{ dateFormat(scope.row['redeemed_at']) }}</span>
          <span v-else class="pill pill-gray">未兑换</span>
        </template>
      </el-table-column>

      <el-table-column label="启用状态" width="100">
        <template #default="scope">
          <el-switch
            v-model="scope.row['enabled']"
            @change="set('enabled', scope.row)"
            :disabled="scope.row['redeemed_at'] > 0"
          />
        </template>
      </el-table-column>

      <el-table-column label="操作" width="90" fixed="right">
        <template #default="scope">
          <el-popconfirm title="确定要删除当前记录吗?" @confirm="remove(scope.row)">
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
          @current-change="fetchData"
        />
        <span class="footer-text">每页条数:</span>
        <el-select v-model="pageSize" class="page-size-select" @change="changePageSize">
          <el-option :value="10" label="10" />
          <el-option :value="20" label="20" />
          <el-option :value="50" label="50" />
        </el-select>
      </div>
    </div>

    <el-empty :image-size="100" v-if="!loading && items.length === 0" />

    <el-dialog v-model="showDialog" title="生成兑换码">
      <template #default>
        <el-form :model="item" label-width="120px">
          <el-form-item label="名称：" prop="name">
            <el-input v-model="item.name" autocomplete="off" />
          </el-form-item>

          <el-form-item label="积分额度：" prop="power">
            <el-input v-model.number="item.power" autocomplete="off" />
          </el-form-item>

          <el-form-item label="生成数量：" prop="num">
            <el-input v-model.number="item.num" />
          </el-form-item>
        </el-form>
      </template>

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
import { showMessageError } from '@/utils/dialog'
import { httpGet, httpPost } from '@/utils/http'
import { dateFormat, substr } from '@/utils/libs'
import { DocumentCopy, Search, Ticket } from '@element-plus/icons-vue'
import ClipboardJS from 'clipboard'
import { ElMessage, ElMessageBox } from 'element-plus'
import { computed, onMounted, onUnmounted, ref } from 'vue'

// 变量定义
const items = ref([])
const loading = ref(true)
const query = ref({ code: '' })
const showDialog = ref(false)
const dialogLoading = ref(false)
const item = ref({ name: '', power: 0, num: 1 })
const itemIds = ref([])
const compact = ref(false)

const page = ref(1)
const pageSize = ref(10)
const total = ref(0)

const totalPages = computed(() => Math.max(1, Math.ceil((total.value || 0) / (pageSize.value || 10))))
const startIndex = computed(() => (total.value ? (page.value - 1) * pageSize.value + 1 : 0))
const endIndex = computed(() => Math.min(page.value * pageSize.value, total.value || 0))

const clipboard = ref(null)
onMounted(() => {
  clipboard.value = new ClipboardJS('.copy-code')
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

const add = () => {
  item.value = { name: '100积分点卡', power: 100, num: 1 }
  showDialog.value = true
  dialogLoading.value = false
}

const save = () => {
  if (item.value.name === '') {
    return showMessageError('请输入兑换码名称')
  }
  if (item.value.power === 0) {
    return showMessageError('请输入积分额度')
  }
  if (item.value.num <= 0) {
    return showMessageError('请输入生成数量')
  }
  dialogLoading.value = true
  httpPost('/api/admin/redeem/create', item.value)
    .then((res) => {
      ElMessage.success(`成功生成了${res.data.counter}个兑换码`)
      showDialog.value = false
      fetchData()
    })
    .catch((e) => {
      ElMessage.error('生成失败：' + e.message)
    })
}

const set = (filed, row) => {
  httpPost('/api/admin/redeem/set', { id: row.id, filed: filed, value: row[filed] })
    .then(() => {
      ElMessage.success('操作成功！')
    })
    .catch((e) => {
      ElMessage.error('操作失败：' + e.message)
    })
}

// 拉取全部兑换码（后端旧版本不支持关键字与批量清除，改为前端过滤）
const fetchAll = async () => {
  const size = 200
  let p = 1
  let all = []
  for (;;) {
    const res = await httpGet('/api/admin/redeem/list', { page: p, page_size: size })
    const rows = (res.data && res.data.items) || []
    all = all.concat(rows)
    const count = (res.data && res.data.total) || 0
    if (rows.length === 0 || all.length >= count) {
      break
    }
    p++
  }
  return all
}

const fetchData = () => {
  loading.value = true
  fetchAll()
    .then((all) => {
      const kw = (query.value.code || '').trim().toLowerCase()
      const filtered = kw
        ? all.filter(
            (it) =>
              String(it.id).includes(kw) ||
              String(it.name || '').toLowerCase().includes(kw) ||
              String(it.code || '').toLowerCase().includes(kw)
          )
        : all
      total.value = filtered.length
      const maxPage = Math.max(1, Math.ceil(filtered.length / (pageSize.value || 10)))
      if (page.value > maxPage) {
        page.value = maxPage
      }
      const start = (page.value - 1) * pageSize.value
      items.value = filtered.slice(start, start + pageSize.value)
      loading.value = false
    })
    .catch((e) => {
      loading.value = false
      ElMessage.error('获取数据失败：' + e.message)
    })
}

const search = () => {
  page.value = 1
  fetchData()
}

const reset = () => {
  query.value.code = ''
  page.value = 1
  fetchData()
}

const changePageSize = () => {
  page.value = 1
  fetchData()
}

const remove = function (row) {
  httpGet('/api/admin/redeem/remove?id=' + row.id)
    .then(() => {
      ElMessage.success('删除成功！')
      fetchData()
    })
    .catch((e) => {
      ElMessage.error('删除失败：' + e.message)
    })
}

const handleSelectionChange = (items) => {
  itemIds.value = items.map((item) => item.id)
}

// 复制所选兑换码
const copySelected = () => {
  if (!itemIds.value.length) {
    ElMessage.warning('请先选择要复制的兑换码')
    return
  }
  const codes = items.value
    .filter((it) => itemIds.value.includes(it.id))
    .map((it) => it.code)
    .join('\n')
  const done = () => ElMessage.success(`已复制 ${itemIds.value.length} 个兑换码`)
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

// 清除失效兑换码
const clearInvalid = () => {
  ElMessageBox.confirm('此操作将会删除所有已兑换或已停用的兑换码，继续操作吗?', '删除提示', {
    confirmButtonText: '确认',
    cancelButtonText: '取消',
    type: 'warning',
  }).then(async () => {
    try {
      const all = await fetchAll()
      const invalid = all.filter((it) => it.redeemed_at > 0 || !it.enabled)
      if (invalid.length === 0) {
        ElMessage.success('没有可清除的失效兑换码')
        return
      }
      let deleted = 0
      for (const row of invalid) {
        try {
          await httpGet('/api/admin/redeem/remove?id=' + row.id)
          deleted++
        } catch (e) {
          // 单条删除失败不影响其余兑换码
        }
      }
      ElMessage.success(`已清除 ${deleted} 个失效兑换码`)
      page.value = 1
      fetchData()
    } catch (e) {
      ElMessage.error('清除失败：' + e.message)
    }
  })
}
</script>

<style lang="scss" scoped>
.redeem-page {
  font-family: 'The Future', Arial, 'Microsoft YaHei', 'PingFang SC', sans-serif;
  font-weight: 400;
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

  .compact-btn {
    height: 30px;
    padding: 0 12px;
    border: none;
    border-radius: 5px;
    background: #f2f4f7;
    color: #1d2129;
    font-size: 13px;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s ease;

    &:hover {
      background: #e9ebf0;
    }

    &.active {
      background: #e8f2ff;
      color: #1677ff;
    }
  }

  .toolbar {
    display: flex;
    align-items: center;
    justify-content: space-between;
    flex-wrap: wrap;
    gap: 12px;
    padding: 12px 20px;
    border-bottom: 1px solid #f0f2f5;
  }

  .toolbar-left {
    display: flex;
    align-items: center;
    gap: 10px;
  }

  .text-action {
    height: 30px;
    padding: 0 12px;
    border: none;
    border-radius: 5px;
    background: #f2f4f7;
    color: #1d2129;
    font-size: 13px;
    font-weight: 500;
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

  .toolbar-right {
    display: flex;
    align-items: center;
    gap: 10px;
    margin-left: auto;
  }

  .search-input {
    width: 260px;

    :deep(.el-input__wrapper) {
      height: 30px;
      background: #f2f4f7;
      box-shadow: none;
      border: 1px solid transparent;
      border-radius: 5px;
      padding: 0 10px;
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
      font-size: 13px;
      font-weight: 400;
      color: #1d2129;
    }

    :deep(.el-input__inner::placeholder) {
      color: #98a2b3;
      font-weight: 400;
    }

    :deep(.el-input__prefix) {
      color: #98a2b3;
      font-size: 14px;
    }
  }

  .plain-btn {
    height: 30px;
    padding: 0 12px;
    border: none;
    border-radius: 5px;
    background: #f2f4f7;
    color: #1d2129;
    font-size: 13px;
    font-weight: 500;
    cursor: pointer;
    transition: background 0.2s ease;

    &:hover {
      background: #e9ebf0;
    }
  }

  .redeem-table {
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
      padding: 10px 0;
    }

    :deep(td.el-table__cell) {
      border-bottom: 1px solid #ebeef5;
      font-size: 14px;
      color: #303133;
      padding: 14px 0;
    }

    :deep(.el-table__body tr:hover > td.el-table__cell) {
      background: #fafbfc;
    }

    &.compact {
      :deep(td.el-table__cell) {
        padding: 7px 0;
      }
    }
  }

  .copy-code {
    cursor: pointer;
    margin-left: 6px;
    position: relative;
    top: 2px;
    font-size: 14px;
    color: #667085;
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
    color: #303133;
    font-size: 14px;
    cursor: pointer;

    &:hover {
      color: #1677ff;
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
    padding: 0 20px 16px;
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

:global([data-theme='dark'] .redeem-page) .panel-head {
  border-bottom-color: rgba(255, 255, 255, 0.08);
}

:global([data-theme='dark'] .redeem-page) :deep(th.el-table__cell) {
  background: rgba(255, 255, 255, 0.04);
  color: rgba(255, 255, 255, 0.55);
}

:global([data-theme='dark'] .redeem-page) :deep(td.el-table__cell) {
  border-bottom-color: rgba(255, 255, 255, 0.06);
  color: rgba(255, 255, 255, 0.85);
}

:global([data-theme='dark'] .redeem-page) .search-input :deep(.el-input__wrapper),
:global([data-theme='dark'] .redeem-page) .plain-btn,
:global([data-theme='dark'] .redeem-page) .compact-btn,
:global([data-theme='dark'] .redeem-page) .text-action {
  background: rgba(255, 255, 255, 0.06);
}
</style>
