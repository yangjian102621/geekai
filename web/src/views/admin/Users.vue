<template>
  <div class="container user-list">
    <div class="handle-box">
      <button class="add-btn" @click="addUser">
        <el-icon><Plus /></el-icon>
        添加用户
      </button>

      <div class="filter-box">
        <el-input
          v-model="query.username"
          placeholder="搜索用户 ID、账号"
          class="filter-input search-input"
          :prefix-icon="Search"
          clearable
          @keyup.enter="handleSearch"
        />
        <el-input
          v-model="query.mobile"
          placeholder="手机"
          class="filter-input"
          clearable
          @keyup.enter="handleSearch"
        />
        <el-input
          v-model="query.email"
          placeholder="邮箱"
          class="filter-input"
          clearable
          @keyup.enter="handleSearch"
        />
        <button class="query-btn" @click="handleSearch">查询</button>
      </div>

      <div class="toolbar-right">
        <button class="text-btn danger" @click="multipleDelete">
          <el-icon><Delete /></el-icon>
          批量删除
        </button>
        <button class="text-btn" @click="showImportDialog = true">导入用户</button>
      </div>
    </div>

    <el-table
      :data="users.items"
      class="user-table"
      :row-key="(row) => row.id"
      @selection-change="handleSelectionChange"
    >
      <el-table-column type="selection" width="40"></el-table-column>
      <el-table-column prop="id" label="ID" width="72" />
      <el-table-column label="账号" min-width="150">
        <template #default="scope">
          <span class="account-name">{{ scope.row.username }}</span>
          <el-tooltip content="复制用户名" placement="top">
            <el-icon class="copy-icon" @click="copyUsername(scope.row.username)">
              <DocumentCopy />
            </el-icon>
          </el-tooltip>
          <el-image v-if="scope.row.vip" :src="vipImg" class="vip-icon" />
        </template>
      </el-table-column>
      <el-table-column prop="mobile" label="手机" min-width="120">
        <template #default="scope">{{ scope.row.mobile || '-' }}</template>
      </el-table-column>
      <el-table-column prop="email" label="邮箱" min-width="150">
        <template #default="scope">{{ scope.row.email || '-' }}</template>
      </el-table-column>
      <el-table-column prop="nickname" label="昵称" min-width="110">
        <template #default="scope">{{ scope.row.nickname || '-' }}</template>
      </el-table-column>
      <el-table-column label="剩余积分" width="110">
        <template #default="scope">
          <span class="pill pill-money">{{ scope.row.power }}</span>
        </template>
      </el-table-column>
      <el-table-column label="状态" width="90">
        <template #default="scope">
          <span class="pill" :class="scope.row.status ? 'pill-green' : 'pill-red'">
            {{ scope.row.status ? '正常' : '停用' }}
          </span>
        </template>
      </el-table-column>
      <el-table-column label="过期时间" width="150">
        <template #default="scope">
          <span v-if="scope.row.expired_time">{{ scope.row.expired_time }}</span>
          <span v-else class="pill pill-blue">长期有效</span>
        </template>
      </el-table-column>
      <el-table-column prop="platform" label="注册渠道" :formatter="registrationChannel" />
      <el-table-column label="注册时间" width="160">
        <template #default="scope">
          <span>{{ dateFormat(scope.row['created_at']) }}</span>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="180" fixed="right">
        <template #default="scope">
          <div class="row-actions">
            <button class="link-btn" @click="userEdit(scope.row)">编辑</button>
            <button class="link-btn blue" @click="resetPass(scope.row)">重置密码</button>
            <el-dropdown trigger="click">
              <button class="link-btn more">
                <i class="iconfont icon-more-horizontal"></i>
              </button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item @click="genLoginLink(scope.row)">生成登录链接</el-dropdown-item>
                  <el-dropdown-item divided @click="removeUser(scope.row)">删除</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </div>
        </template>
      </el-table-column>
    </el-table>

    <div class="table-footer" v-if="users.total > 0">
      <div class="footer-info">
        显示第 {{ startIndex }} - {{ endIndex }} 条，共 {{ users.total }} 条
      </div>
      <div class="footer-nav">
        <span class="footer-text">总页数: {{ totalPages }}</span>
        <el-pagination
          v-model:current-page="users.page"
          :page-size="users.page_size"
          :total="users.total"
          layout="prev, pager, next"
          @current-change="(p) => fetchUserList(p, users.page_size)"
        />
        <span class="footer-text">每页条数:</span>
        <el-select v-model="users.page_size" class="page-size-select" @change="changePageSize">
          <el-option :value="10" label="10" />
          <el-option :value="20" label="20" />
          <el-option :value="50" label="50" />
        </el-select>
      </div>
    </div>

    <el-dialog
      v-model="showUserEditDialog"
      :title="title"
      :close-on-click-modal="false"
      width="50%"
    >
      <el-form :model="user" label-width="100px" ref="userEditFormRef" :rules="rules">
        <el-form-item label="账号：" prop="username">
          <el-input v-model="user.username" autocomplete="off" />
        </el-form-item>
        <el-form-item label="手机：" prop="mobile">
          <el-input v-model="user.mobile" autocomplete="off" />
        </el-form-item>
        <el-form-item label="邮箱：" prop="email">
          <el-input v-model="user.email" autocomplete="off" />
        </el-form-item>
        <el-form-item v-if="add" label="密码：" prop="password">
          <el-input
            v-model="user.password"
            type="password"
            show-password
            autocomplete="off"
            placeholder="8-16位"
          />
        </el-form-item>
        <el-form-item label="剩余积分：" prop="power">
          <el-input v-model.number="user.power" autocomplete="off" placeholder="0" />
        </el-form-item>

        <el-form-item label="有效期：" prop="expired_time">
          <el-date-picker
            v-model="user.expired_time"
            type="datetime"
            placeholder="选择日期"
            format="YYYY-MM-DD HH:mm:ss"
            value-format="YYYY-MM-DD HH:mm:ss"
            :disabled-date="disabledDate"
          />
        </el-form-item>

        <el-form-item label="模型权限" prop="chat_models">
          <el-select
            v-model="user.chat_models"
            multiple
            :filterable="true"
            placeholder="选择AI模型，多选"
          >
            <el-option v-for="item in models" :key="item.id" :label="item.name" :value="item.id" />
          </el-select>
        </el-form-item>

        <el-form-item label="启用状态">
          <el-switch v-model="user.status" />
        </el-form-item>

        <el-form-item label="开通VIP">
          <el-switch v-model="user.vip" />
        </el-form-item>
      </el-form>

      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showUserEditDialog = false">取消</el-button>
          <el-button type="primary" @click="saveUser">提交</el-button>
        </span>
      </template>
    </el-dialog>

    <el-dialog v-model="showResetPassDialog" title="重置密码" width="50%">
      <el-form label-width="100px" ref="userEditFormRef">
        <el-form-item label="账户：">
          <el-input v-model="pass.username" autocomplete="off" readonly disabled />
        </el-form-item>

        <el-form-item label="新密码：">
          <el-input v-model="pass.password" type="password" show-password autocomplete="off" />
        </el-form-item>
      </el-form>

      <template #footer>
        <span class="dialog-footer">
          <el-button type="primary" @click="doResetPass">提交</el-button>
        </span>
      </template>
    </el-dialog>

    <el-dialog v-model="showGenLoginLinkDialog" title="自动登录链接" width="50%">
      <el-input v-model="loginLink" readonly disabled />
      <template #footer>
        <span class="dialog-footer">
          <el-button type="primary" @click="copyLoginLink">复制链接</el-button>
        </span>
      </template>
    </el-dialog>

    <el-dialog v-model="showImportDialog" title="批量导入用户" width="50%">
      <div class="mb-2 text-sm text-gray-500 flex items-center justify-between">
        <span>仅支持 .xlsx 格式的 Excel 文件。请先下载模板，按要求填写后再上传。</span>
        <el-button type="primary" link @click="downloadTemplate">下载模板</el-button>
      </div>
      <el-upload
        class="upload-demo"
        drag
        :auto-upload="true"
        :show-file-list="false"
        :before-upload="beforeImportUpload"
        :http-request="doImport"
        accept=".xlsx"
      >
        <el-icon class="el-icon--upload"><DocumentCopy /></el-icon>
        <div class="el-upload__text">将文件拖到此处，或 <em>点击上传</em></div>
        <template #tip>
          <div class="el-upload__tip text-gray-500 text-sm">
            只支持扩展名为 .xlsx 的 Excel 文件。
          </div>
        </template>
      </el-upload>

      <div v-if="importResult" class="mt-4 text-sm">
        <p>导入成功：{{ importResult.success }} 条，失败：{{ importResult.failed }} 条。</p>
        <ul v-if="importResult.errors && importResult.errors.length" class="mt-2 text-red-500">
          <li v-for="err in importResult.errors" :key="err.row">
            第 {{ err.row }} 行：{{ err.error }}
          </li>
        </ul>
      </div>

      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showImportDialog = false">关闭</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { closeLoading, showLoading, showMessageError, showMessageOK } from '@/utils/dialog'
import { httpDownload, httpGet, httpPost } from '@/utils/http'
import { dateFormat, disabledDate } from '@/utils/libs'
import { Delete, DocumentCopy, Plus, Search } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { computed, onMounted, reactive, ref } from 'vue'

// 变量定义
const users = ref({ page: 1, page_size: 10, items: [] })
const query = ref({ username: '', page: 1, page_size: 10 })

const totalPages = computed(() =>
  Math.max(1, Math.ceil((users.value.total || 0) / (users.value.page_size || 10)))
)
const startIndex = computed(() =>
  users.value.total ? (users.value.page - 1) * users.value.page_size + 1 : 0
)
const endIndex = computed(() =>
  Math.min(users.value.page * users.value.page_size, users.value.total || 0)
)

const changePageSize = (size) => {
  users.value.page = 1
  fetchUserList(1, size)
}

const title = ref('添加用户')
const vipImg = ref('/images/menu/member.png')
const add = ref(true)
const user = ref({ chat_models: [] })
const pass = ref({ username: '', password: '', id: 0 })
const models = ref([])
const showUserEditDialog = ref(false)
const showResetPassDialog = ref(false)
const showImportDialog = ref(false)
const importResult = ref(null)
const rules = reactive({
  username: [{ required: true, message: '请输入账号', trigger: 'blur' }],
  password: [
    {
      required: true,
      validator: (rule, value) => {
        return !(value.length > 16 || value.length < 8)
      },
      message: '密码必须为8-16',
      trigger: 'blur',
    },
  ],
  calls: [
    { required: true, message: '请输入提问次数' },
    { type: 'number', message: '请输入有效数字' },
  ],
})
const loading = ref(true)

const userEditFormRef = ref(null)

const registrationChannel = ({ platform, openid }) => {
  if (platform === 'tokengine') return 'Tokengine'
  if (platform === 'wechat' || (!platform && openid)) return '微信'
  return platform || '未记录'
}

onMounted(() => {
  fetchUserList(users.value.page, users.value.page_size)
  httpGet('/api/admin/model/list')
    .then((res) => {
      models.value = res.data
    })
    .catch((e) => {
      ElMessage.error('获取模型失败：' + e.message)
    })
})

const fetchUserList = function (page, pageSize) {
  query.value.page = page
  query.value.page_size = pageSize
  httpGet('/api/admin/user/list', query.value)
    .then((res) => {
      if (res.data) {
        // 初始化数据
        const arr = res.data.items
        for (let i = 0; i < arr.length; i++) {
          arr[i].expired_time = dateFormat(arr[i].expired_time)
        }
        users.value.items = arr
        users.value.total = res.data.total
        users.value.page = res.data.page
        user.value.page_size = res.data.page_size
      }
      loading.value = false
    })
    .catch(() => {
      ElMessage.error('加载用户列表失败')
    })
}

const handleSearch = () => {
  fetchUserList(users.value.page, users.value.page_size)
}

// 删除用户
const removeUser = function (user) {
  ElMessageBox.confirm('此操作将会永久删除用户信息和聊天记录，确认操作吗?', '警告', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning',
  })
    .then(() => {
      httpGet('/api/admin/user/remove', { id: user.id })
        .then(() => {
          ElMessage.success('操作成功！')
          fetchUserList(users.value.page, users.value.page_size)
        })
        .catch((e) => {
          ElMessage.error('操作失败，' + e.message)
        })
    })
    .catch(() => {
      ElMessage.info('操作被取消')
    })
}

const userEdit = function (row) {
  user.value = row
  title.value = '编辑用户'
  showUserEditDialog.value = true
  add.value = false
}

const addUser = () => {
  user.value = { chat_id: 0, chat_models: [] }
  title.value = '添加用户'
  showUserEditDialog.value = true
  add.value = true
}

const saveUser = function () {
  userEditFormRef.value.validate((valid) => {
    if (valid) {
      showUserEditDialog.value = false
      console.log(user.value)
      httpPost('/api/admin/user/save', user.value)
        .then((res) => {
          ElMessage.success('操作成功！')
          if (add.value) {
            users.value.items.push(res.data)
          }
        })
        .catch((e) => {
          ElMessage.error('操作失败，' + e.message)
        })
    } else {
      return false
    }
  })
}

const userIds = ref([])
const handleSelectionChange = function (rows) {
  userIds.value = []
  rows.forEach((row) => {
    userIds.value.push(row.id)
  })
}

const multipleDelete = function () {
  ElMessageBox.confirm('此操作将会永久删除用户信息和聊天记录，确认操作吗?', '警告', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning',
  })
    .then(() => {
      loading.value = true
      httpGet('/api/admin/user/remove', { ids: userIds.value })
        .then(() => {
          ElMessage.success('操作成功！')
          fetchUserList(users.value.page, users.value.page_size)
          loading.value = false
        })
        .catch((e) => {
          ElMessage.error('操作失败，' + e.message)
          loading.value = false
        })
    })
    .catch(() => {
      ElMessage.info('操作被取消')
    })
}

// 重置密码
const resetPass = (row) => {
  showResetPassDialog.value = true
  pass.value.id = row.id
  pass.value.username = row.username
}

const doResetPass = () => {
  httpPost('/api/admin/user/resetPass', pass.value)
    .then(() => {
      ElMessage.success('操作成功！')
      showResetPassDialog.value = false
    })
    .catch((e) => {
      ElMessage.error('操作失败，' + e.message)
    })
}

// 生成登录链接
const showGenLoginLinkDialog = ref(false)
const loginLink = ref('')
const genLoginLink = (row) => {
  showLoading()
  httpGet('/api/admin/user/genLoginLink', { id: row.id })
    .then((res) => {
      loginLink.value = `${window.location.origin}/login?token=${res.data}`
      showGenLoginLinkDialog.value = true
      closeLoading()
    })
    .catch((e) => {
      ElMessage.error('操作失败，' + e.message)
      closeLoading()
    })
}

const copyLoginLink = () => {
  try {
    navigator.clipboard.writeText(loginLink.value)
    showMessageOK('复制成功！')
  } catch (e) {
    showMessageError('复制失败')
  }
}

// 复制用户名
const copyUsername = (username) => {
  try {
    navigator.clipboard.writeText(username)
    ElMessage.success('用户名复制成功！')
  } catch (e) {
    ElMessage.error('复制失败')
  }
}

// 下载导入模板
const downloadTemplate = () => {
  httpDownload('/api/admin/user/import/template')
    .then((response) => {
      const blob = new Blob([response.data])
      const url = window.URL.createObjectURL(blob)
      const link = document.createElement('a')
      link.href = url
      link.setAttribute('download', 'user_import_template.xlsx')
      document.body.appendChild(link)
      link.click()
      document.body.removeChild(link)
      window.URL.revokeObjectURL(url)
    })
    .catch((e) => {
      ElMessage.error('下载模板失败：' + (e.message || '未知错误'))
    })
}

// 导入前校验
const beforeImportUpload = (file) => {
  const ext = file.name.slice(file.name.lastIndexOf('.')).toLowerCase()
  if (ext !== '.xlsx') {
    ElMessage.error('只支持 .xlsx 格式的 Excel 文件')
    return false
  }
  return true
}

// 执行导入
const doImport = (option) => {
  const formData = new FormData()
  formData.append('file', option.file)

  httpPost('/api/admin/user/import', formData, {
    headers: { 'Content-Type': 'multipart/form-data' },
  })
    .then((res) => {
      importResult.value = res.data
      ElMessage.success('导入完成')
      fetchUserList(users.value.page, users.value.page_size)
      option.onSuccess(res)
    })
    .catch((e) => {
      ElMessage.error('导入失败：' + (e.message || '未知错误'))
      option.onError(e)
    })
}
</script>

<style lang="stylus" scoped>
.user-list {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'PingFang SC', 'Microsoft YaHei', sans-serif;

  .el-select {
    width: 100%;
  }

  .handle-box {
    display: flex;
    align-items: center;
    flex-wrap: wrap;
    gap: 10px;
    margin-bottom: 16px;
  }

  .add-btn {
    display: inline-flex;
    align-items: center;
    gap: 6px;
    height: 34px;
    padding: 0 14px;
    border: 1px solid #e2e6ee;
    border-radius: 6px;
    background: #f2f4f7;
    color: #1677ff;
    font-size: 14px;
    font-weight: 500;
    cursor: pointer;
    white-space: nowrap;
    transition: all 0.2s ease;

    &:hover {
      border-color: #1677ff;
      background: #e9ebf0;
    }
  }

  .filter-box {
    display: flex;
    align-items: center;
    flex-wrap: wrap;
    gap: 10px;
  }

  .filter-input {
    width: 150px;
    font-size: 14px;

    :deep(.el-input__wrapper) {
      height: 34px;
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

    &.search-input {
      width: 300px;
    }
  }

  .query-btn {
    height: 34px;
    padding: 0 16px;
    border: 1px solid #e2e6ee;
    border-radius: 6px;
    background: #f2f4f7;
    color: #1d2129;
    font-size: 14px;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s ease;

    &:hover {
      border-color: #1677ff;
      background: #e9ebf0;
      color: #1677ff;
    }
  }

  .toolbar-right {
    display: flex;
    align-items: center;
    gap: 14px;
    margin-left: auto;
  }

  .text-btn {
    display: inline-flex;
    align-items: center;
    gap: 4px;
    border: none;
    background: transparent;
    padding: 4px 0;
    color: #667085;
    font-size: 13px;
    cursor: pointer;

    &:hover {
      color: #374151;
    }

    &.danger {
      color: #ef4444;

      &:hover {
        color: #dc2626;
      }
    }
  }

  .user-table {
    width: 100%;

    :deep(.el-table__inner-wrapper::before) {
      display: none;
    }

    :deep(th.el-table__cell) {
      background: #f7f8fa;
      color: #667085;
      font-size: 13px;
      font-weight: 500;
      border-bottom: none;
      padding: 10px 0;
    }

    :deep(td.el-table__cell) {
      border-bottom: 1px solid #f2f3f5;
      font-size: 13px;
      color: #1d2129;
      padding: 14px 0;
    }

    :deep(.el-table__body tr:hover > td.el-table__cell) {
      background: #fafbfc;
    }
  }

  .account-name {
    font-weight: 500;
  }

  .copy-icon {
    margin-left: 8px;
    color: #b0b3b8;
    cursor: pointer;
    vertical-align: middle;
    transition: all 0.2s ease;

    &:hover {
      color: #1677ff;
    }
  }

  .vip-icon {
    height: 18px;
    margin-left: 6px;
    vertical-align: middle;
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

    &.pill-green {
      background: #e9f9ef;
      color: #22c55e;
    }

    &.pill-red {
      background: #fdecec;
      color: #ef4444;
    }

    &.pill-blue {
      background: #e8f2ff;
      color: #1677ff;
    }

    &.pill-money {
      background: #fff;
      border: 1px solid #e5e7eb;
      color: #374151;
      font-variant-numeric: tabular-nums;
    }
  }

  .row-actions {
    display: flex;
    align-items: center;
    gap: 12px;
  }

  .link-btn {
    border: none;
    background: transparent;
    padding: 0;
    color: #667085;
    font-size: 13px;
    cursor: pointer;
    white-space: nowrap;

    &:hover {
      color: #1d2129;
    }

    &.blue {
      color: #1677ff;

      &:hover {
        color: #0958d9;
      }
    }

    &.more {
      display: inline-flex;
      align-items: center;
      color: #9ca3af;
      font-size: 16px;
    }
  }

  .table-footer {
    display: flex;
    justify-content: space-between;
    align-items: center;
    flex-wrap: wrap;
    gap: 12px;
    margin-top: 16px;
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

:global([data-theme='dark'] .user-list) :deep(th.el-table__cell) {
  background: rgba(255, 255, 255, 0.04);
  color: rgba(255, 255, 255, 0.55);
}

:global([data-theme='dark'] .user-list) :deep(td.el-table__cell) {
  border-bottom-color: rgba(255, 255, 255, 0.06);
  color: rgba(255, 255, 255, 0.85);
}

:global([data-theme='dark'] .user-list) :deep(.el-table__body tr:hover > td.el-table__cell) {
  background: rgba(255, 255, 255, 0.03);
}

:global([data-theme='dark'] .user-list) .filter-input :deep(.el-input__wrapper) {
  background: rgba(255, 255, 255, 0.06);
}

:global([data-theme='dark'] .user-list) .pill-money {
  background: transparent;
  border-color: rgba(255, 255, 255, 0.14);
  color: rgba(255, 255, 255, 0.85);
}
</style>
