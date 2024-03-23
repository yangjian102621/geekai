<template>
  <div class="container user-list" v-loading="loading">
    <div class="handle-box">
      <el-input v-model="query.username" placeholder="账号" class="handle-input mr10"></el-input>
      <el-button type="primary" :icon="Search" @click="handleSearch">搜索</el-button>

      <el-button type="success" :icon="Plus" @click="addUser">新增用户</el-button>
    </div>

    <el-row>
      <el-table :data="users.items" border class="table" :row-key="row => row.id"
                @selection-change="handleSelectionChange" table-layout="auto">
        <el-table-column type="selection" width="38"/>
        <el-table-column prop="mobile" label="账号">
          <template #default="scope">
            <span>{{ scope.row.username }}</span>
            <el-image v-if="scope.row.vip" :src="vipImg" style="height: 20px;position: relative; top:5px; left: 5px"/>
          </template>
        </el-table-column>
        <el-table-column prop="nickname" label="昵称"/>
        <el-table-column prop="power" label="剩余算力"/>
        <el-table-column label="状态" width="80">
          <template #default="scope">
            <el-tag v-if="scope.row.status" type="success">正常</el-tag>
            <el-tag type="danger" v-else>停用</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="过期时间">
          <template #default="scope">
            <span v-if="scope.row['expired_time']">{{ scope.row['expired_time'] }}</span>
            <el-tag v-else>长期有效</el-tag>
          </template>
        </el-table-column>

        <el-table-column label="注册时间">
          <template #default="scope">
            <span>{{ dateFormat(scope.row['created_at']) }}</span>
          </template>
        </el-table-column>

        <el-table-column fixed="right" label="操作" width="200">
          <template #default="scope">
            <el-button-group class="ml-4">
              <el-button size="small" type="primary" @click="userEdit(scope.row)">编辑</el-button>
              <el-button size="small" type="danger" @click="removeUser(scope.row)">删除</el-button>
              <el-button size="small" type="success" @click="resetPass(scope.row)">重置密码</el-button>
            </el-button-group>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination">
        <el-pagination v-if="users.total > 0"
                       background
                       layout="total, prev, pager, next"
                       v-model:current-page="users.page"
                       v-model:page-size="users.page_size"
                       :total="users.total"
                       @current-change="fetchUserList(users.page, users.page_size)"
        />

      </div>
    </el-row>

    <el-dialog
        v-model="showUserEditDialog"
        :title="title"
        :close-on-click-modal="false"
        width="50%"
    >
      <el-form :model="user" label-width="100px" ref="userEditFormRef" :rules="rules">
        <el-form-item label="账号：" prop="username">
          <el-input v-model="user.username" autocomplete="off"/>
        </el-form-item>
        <el-form-item v-if="add" label="密码：" prop="password">
          <el-input v-model="user.password" autocomplete="off"/>
        </el-form-item>
        <el-form-item label="剩余算力：" prop="power">
          <el-input v-model.number="user.power" autocomplete="off" placeholder="0"/>
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

        <el-form-item label="模型权限" prop="chat_models">
          <el-select
              v-model="user.chat_models"
              multiple
              :filterable="true"
              placeholder="选择AI模型，多选"
          >
            <el-option
                v-for="item in models"
                :key="item.id"
                :label="item.name"
                :value="item.id"
            />
          </el-select>
        </el-form-item>


        <el-form-item label="启用状态">
          <el-switch v-model="user.status"/>
        </el-form-item>

        <el-form-item label="开通VIP">
          <el-switch v-model="user.vip"/>
        </el-form-item>
      </el-form>

      <template #footer>
            <span class="dialog-footer">
              <el-button @click="showUserEditDialog = false">取消</el-button>
              <el-button type="primary" @click="saveUser">提交</el-button>
            </span>
      </template>
    </el-dialog>

    <el-dialog
        v-model="showResetPassDialog"
        title="重置密码"
        width="50%"
    >
      <el-form label-width="100px" ref="userEditFormRef">
        <el-form-item label="账户：">
          <el-input v-model="pass.username" autocomplete="off" readonly disabled/>
        </el-form-item>

        <el-form-item label="新密码：">
          <el-input v-model="pass.password" autocomplete="off"/>
        </el-form-item>
      </el-form>

      <template #footer>
            <span class="dialog-footer">
              <el-button type="primary" @click="doResetPass">提交</el-button>
            </span>
      </template>

    </el-dialog>
  </div>
</template>

<script setup>
import {onMounted, reactive, ref} from "vue";
import {httpGet, httpPost} from "@/utils/http";
import {ElMessage, ElMessageBox} from "element-plus";
import {dateFormat, disabledDate, removeArrayItem} from "@/utils/libs";
import {Plus, Search} from "@element-plus/icons-vue";

// 变量定义
const users = ref({page: 1, page_size: 15, items: []})
const query = ref({username: '', page: 1, page_size: 15})

const title = ref('添加用户')
const vipImg = ref("/images/vip.png")
const add = ref(true)
const user = ref({chat_roles: [], chat_models: []})
const pass = ref({username: '', password: '', id: 0})
const roles = ref([])
const models = ref([])
const showUserEditDialog = ref(false)
const showResetPassDialog = ref(false)
const rules = reactive({
  username: [{required: true, message: '请输入账号', trigger: 'change',}],
  password: [{required: true, message: '请输入密码', trigger: 'change',}],
  calls: [
    {required: true, message: '请输入提问次数'},
    {type: 'number', message: '请输入有效数字'},
  ],
  chat_roles: [{required: true, message: '请选择聊天角色', trigger: 'change'}],
  chat_models: [{required: true, message: '请选择AI模型', trigger: 'change'}],
})
const loading = ref(true)

const userEditFormRef = ref(null)

onMounted(() => {
  fetchUserList(users.value.page, users.value.page_size)
  // 获取角色列表
  httpGet('/api/admin/role/list').then((res) => {
    roles.value = res.data;
  }).catch(() => {
    ElMessage.error("获取聊天角色失败");
  })

  httpGet('/api/admin/model/list').then(res => {
    models.value = res.data
  }).catch(e => {
    ElMessage.error("获取模型失败：" + e.message)
  })
})

const fetchUserList = function (page, pageSize) {
  query.value.page = page
  query.value.page_size = pageSize
  httpGet('/api/admin/user/list', query.value).then((res) => {
    if (res.data) {
      // 初始化数据
      const arr = res.data.items;
      for (let i = 0; i < arr.length; i++) {
        arr[i].expired_time = dateFormat(arr[i].expired_time)
      }
      users.value.items = arr
      users.value.total = res.data.total
      users.value.page = res.data.page
      user.value.page_size = res.data.page_size
    }
    loading.value = false
  }).catch(() => {
    ElMessage.error('加载用户列表失败')
  })
}

const handleSearch = () => {
  fetchUserList(users.value.page, users.value.page_size)
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
  ).then(() => {
    httpGet('/api/admin/user/remove', {id: user.id}).then(() => {
      ElMessage.success('操作成功！')
      users.value.items = removeArrayItem(users.value.items, user, function (v1, v2) {
        return v1.id === v2.id
      })
    }).catch((e) => {
      ElMessage.error('操作失败，' + e.message)
    })
  }).catch(() => {
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
  user.value = {}
  title.value = '添加用户'
  showUserEditDialog.value = true
  add.value = true
}

const saveUser = function () {
  userEditFormRef.value.validate((valid) => {
    if (valid) {
      showUserEditDialog.value = false
      console.log(user.value)
      httpPost('/api/admin/user/save', user.value).then((res) => {
        ElMessage.success('操作成功！')
        if (add.value) {
          users.value.items.push(res.data)
        }
      }).catch((e) => {
        ElMessage.error('操作失败，' + e.message)
      })
    } else {
      return false
    }
  })
}

const handleSelectionChange = function (rows) {
  console.log(rows)
}

const resetPass = (row) => {
  showResetPassDialog.value = true
  pass.value.id = row.id
  pass.value.username = row.username
}

const doResetPass = () => {
  httpPost('/api/admin/user/resetPass', pass.value).then(() => {
    ElMessage.success('操作成功！')
    showResetPassDialog.value = false
  }).catch((e) => {
    ElMessage.error('操作失败，' + e.message)
  })
}

</script>

<style lang="stylus" scoped>
.user-list {

  .handle-box {
    .handle-input {
      max-width 150px;
      margin-right 10px;
    }
  }

  .table {
    width 100%
  }

  .opt-box {
    padding-bottom: 10px;

    .el-icon {
      margin-right: 5px;
    }
  }

  .el-select {
    width: 100%
  }

}
</style>