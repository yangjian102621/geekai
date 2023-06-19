<template>
  <div class="user-list" v-loading="loading">
    <el-row>
      <el-table :data="users.items">
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
            <el-popconfirm
                width="220"
                confirm-button-text="确定"
                cancel-button-text="取消"
                title="确定删除该记录吗?"
                :hide-after="0"
                @confirm="removeUser(scope.row)"
            >
              <template #reference>
                <el-button size="small" type="danger">删除</el-button>
              </template>
            </el-popconfirm>

          </template>
        </el-table-column>
      </el-table>

      <div class="pagination">
        <el-pagination background
                       layout="prev, pager, next"
                       :hide-on-single-page="true"
                       v-model:current-page="users.page"
                       v-model:page-size="users['page_size']"
                       :total="1000"/>
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

        <el-form-item label="有效期：" prop="term">
          <el-date-picker
              v-model="form2.expired_time"
              type="datetime"
              placeholder="选择时间"
              format="YYYY-MM-DD HH:mm:ss"
              value-format="YYYY-MM-DD HH:mm:ss"
          />
        </el-form-item>

        <el-form-item label="聊天角色" prop="chat_roles">
          <el-select
              v-model="form2.chat_roles"
              multiple
              :filterable="true"
              placeholder="选择聊天角色，多选"
              @change="selectRole"
          >
            <el-option
                v-for="item in roles"
                :key="item.key"
                :label="item.name"
                :value="item.key"
                :disabled="item.disabled"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="聊天记录">
          <el-switch v-model="form2.enable_history"/>
        </el-form-item>

        <el-form-item label="启用状态">
          <el-switch v-model="form2.status"/>
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
import {ElMessage} from "element-plus";
import {arrayContains, dateFormat, removeArrayItem} from "@/utils/libs";

// 变量定义
const users = ref({})

const user = ref({chat_roles: []})
const roles = ref([])
const showUserDialog = ref(false)
const showUserEditDialog = ref(false)
const rules = reactive({
  name: [{required: true, message: '请输入口令名称', trigger: 'change',}],
  max_calls: [
    {required: true, message: '请输入提问次数'},
    {type: 'number', message: '请输入有效数字'},
  ],
  remaining_calls: [
    {required: true, message: '请输入提问次数'},
    {type: 'number', message: '请输入有效数字'},
  ],
  term: [
    {required: true, message: '请输入有效期', trigger: 'change'},
    {type: 'number', message: '请输入有效数字'},
  ],
  number: [
    {required: true, message: '请输入口令数量', trigger: 'change'},
    {type: 'number', message: '请输入有效数字'},
  ],
  chat_roles: [{required: true, message: '请选择聊天角色', trigger: 'change'}],
})
const loading = ref(true)

const userAddFormRef = ref(null)
const userEditFormRef = ref(null)

onMounted(() => {
// 获取口令列表
  httpGet('/api/admin/user/list', {page: 1, page_size: 20}).then((res) => {
    users.value = res.data;
  }).catch(() => {
    ElMessage.error('加载用户列表失败')
  })

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

// 选择角色事件
const selectRole = function (items) {
  if (arrayContains(items, 'all')) {
    for (let i = 0; i < roles.value.length; i++) {
      if (roles.value[i].key === 'all') {
        continue
      }
      roles.value[i].disabled = true
      form1.value.chat_roles = ['all']
      form2.value.chat_roles = ['all']
      form3.value.chat_roles = ['all']
    }
  } else {
    for (let i = 0; i < roles.value.length; i++) {
      if (roles.value[i].key === 'all') {
        continue
      }
      roles.value[i].disabled = false
    }
  }

}

// 删除用户
const removeUser = function (user) {
  httpPost('/api/admin/user/remove', {name: user.name}).then(() => {
    ElMessage.success('操作成功！')
    users.value = removeArrayItem(users.value, user, function (v1, v2) {
      return v1.name === v2.name
    })
  }).catch((e) => {
    ElMessage.error('操作失败，' + e.message)
  })
}

const userEdit = function (_user) {
  user.value = _user
  showUserEditDialog.value = true
}

// 更新口令
const updateUser = function () {
  userEditFormRef.value.validate((valid) => {
    if (valid) {
      showUserEditDialog.value = false
      form2.value.term = parseInt(form2.value.term)
      form2.value.remaining_calls = parseInt(form2.value.remaining_calls)
      httpPost('/api/admin/user/update', form2.value).then(() => {
        ElMessage.success('操作成功！')
      }).catch((e) => {
        ElMessage.error('操作失败，' + e.message)
      })
    } else {
      return false
    }
  })
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