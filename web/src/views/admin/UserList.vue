<template>
  <div class="user-list" v-loading="loading">
    <el-row>
      <el-table :data="users.items" :row-key="row => row.id" @selection-change="handleSelectionChange">
        <el-table-column type="selection" width="55"/>
        <el-table-column prop="username" label="用户名"/>
        <el-table-column prop="nickname" label="昵称"/>
        <el-table-column prop="calls" label="提问次数" width="100"/>
        <el-table-column label="状态" width="80">
          <template #default="scope">
            <el-tag v-if="scope.row.status" type="success">正常</el-tag>
            <el-tag type="danger" v-else>停用</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="过期时间">
          <template #default="scope">
            <span v-if="scope.row['expired_time'] > 0">{{ dateFormat(scope.row['expired_time']) }}</span>
            <el-tag v-else>长期有效</el-tag>
          </template>
        </el-table-column>

        <el-table-column label="注册时间">
          <template #default="scope">
            <span>{{ dateFormat(scope.row['created_at']) }}</span>
          </template>
        </el-table-column>

        <el-table-column label="操作" width="180">
          <template #default="scope">
            <el-button size="small" type="primary" @click="userEdit(scope.row)">编辑</el-button>
            <el-button size="small" type="danger" @click="removeUser(scope.row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination">
        <el-pagination v-if="users.total > 0" background
                       layout="prev, pager, next"
                       :hide-on-single-page="true"
                       v-model:current-page="users.page"
                       v-model:page-size="users.page_size"
                       @update:current-change="fetchUserList(users.page, users.page_size)"
                       :page-count="users.total_page"/>

      </div>
    </el-row>

    <el-dialog
        v-model="showUserEditDialog"
        title="编辑用户"
        width="50%"
    >
      <el-form :model="user" label-width="100px" ref="userEditFormRef" :rules="rules">
        <el-form-item label="昵称：" prop="nickname">
          <el-input v-model="user.nickname" autocomplete="off"/>
        </el-form-item>

        <el-form-item label="提问次数：" prop="calls">
          <el-input v-model.number="user.calls" autocomplete="off" placeholder="0"/>
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

        <el-form-item label="聊天角色" prop="chat_roles">
          <el-select
              v-model="user.chat_roles"
              multiple
              :filterable="true"
              placeholder="选择聊天角色，多选"
          >
            <el-option
                v-for="item in roles"
                :key="item.key"
                :label="item.name"
                :value="item.key"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="启用状态">
          <el-switch v-model="user.status"/>
        </el-form-item>
      </el-form>

      <template #footer>
            <span class="dialog-footer">
              <el-button @click="showUserEditDialog = false">取消</el-button>
              <el-button type="primary" @click="updateUser">提交</el-button>
            </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import {nextTick, onMounted, reactive, ref} from "vue";
import {httpGet, httpPost} from "@/utils/http";
import {ElMessage, ElMessageBox} from "element-plus";
import {dateFormat, removeArrayItem} from "@/utils/libs";

// 变量定义
const users = ref({})

const user = ref({chat_roles: []})
const roles = ref([])
const showUserEditDialog = ref(false)
const rules = reactive({
  nickname: [{required: true, message: '请输入昵称', trigger: 'change',}],
  calls: [
    {required: true, message: '请输入提问次数'},
    {type: 'number', message: '请输入有效数字'},
  ],
  chat_roles: [{required: true, message: '请选择聊天角色', trigger: 'change'}],
})
const loading = ref(true)

const userEditFormRef = ref(null)

onMounted(() => {
  fetchUserList(1, 20)
  // 获取角色列表
  httpGet('/api/admin/role/list').then((res) => {
    roles.value = res.data;
  }).catch(() => {
    ElMessage.error("获取聊天角色失败");
  })

  nextTick(() => {
    loading.value = false
  })
})

const fetchUserList = function (page, pageSize) {
  httpGet('/api/admin/user/list', {page: page, page_size: pageSize}).then((res) => {
    users.value = res.data;
  }).catch(() => {
    ElMessage.error('加载用户列表失败')
  })
}

const disabledDate = (time) => {
  return time.getTime() < Date.now()
}

// 删除用户
const removeUser = function (user) {
  ElMessageBox.confirm(
      '此操作将会永久删除用户信息和聊天记录，确认操作吗?',
      '警告',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
  )
      .then(() => {
        httpGet('/api/admin/user/remove', {id: user.id}).then(() => {
          ElMessage.success('操作成功！')
          users.value.items = removeArrayItem(users.value.items, user, function (v1, v2) {
            return v1.id === v2.id
          })
        }).catch((e) => {
          ElMessage.error('操作失败，' + e.message)
        })
      })
      .catch(() => {
        ElMessage.info('操作被取消')
      })

}

const userEdit = function (_user) {
  _user.expired_time = dateFormat(_user.expired_time)
  user.value = _user
  showUserEditDialog.value = true
}

// 更新口令
const updateUser = function () {
  userEditFormRef.value.validate((valid) => {
    if (valid) {
      showUserEditDialog.value = false
      httpPost('/api/admin/user/update', user.value).then(() => {
        ElMessage.success('操作成功！')
      }).catch((e) => {
        ElMessage.error('操作失败，' + e.message)
      })
    } else {
      return false
    }
  })
}

const handleSelectionChange = function (rows) {
  // TODO: 批量删除操作
  console.log(rows)
}
</script>

<style lang="stylus" scoped>
.user-list {

  .opt-box {
    padding-bottom: 10px;

    .el-icon {
      margin-right: 5px;
    }
  }

  .pagination {
    padding-top 20px;
    display flex
    justify-content center
    width 100%
  }

  .el-select {
    width: 100%
  }

}
</style>