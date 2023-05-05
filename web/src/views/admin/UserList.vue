<template>
  <div class="user-list" v-loading="loading">
    <el-row class="opt-box">
      <el-button type="primary" @click="showUserDialog = true">
        <el-icon>
          <Plus/>
        </el-icon>
        新增用户
      </el-button>

      <el-button type="success" @click="showBatchAddUserDialog = true">
        <el-icon>
          <Plus/>
        </el-icon>
        批量新增
      </el-button>
    </el-row>

    <el-row>
      <el-table :data="users">
        <el-table-column prop="name" label="用户名"/>
        <el-table-column prop="max_calls" label="最大提问次数"/>
        <el-table-column prop="remaining_calls" label="剩余提问次数"/>
        <el-table-column label="激活时间" width="180">
          <template #default="scope">
            <el-tag type="info" v-if="scope.row.active_time === ''">未激活</el-tag>
            <span v-else>{{ scope.row.active_time }}</span>
          </template>
        </el-table-column>
        <el-table-column label="过期时间" width="180">
          <template #default="scope">
            <el-tag type="info" v-if="scope.row.expired_time === ''">未激活</el-tag>
            <span v-else>{{ scope.row.expired_time }}</span>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="180">
          <template #default="scope">
            <el-tag v-if="scope.row.status" type="success">正常</el-tag>
            <el-tag type="danger" v-else>停用</el-tag>
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
    </el-row>

    <el-dialog
        v-model="showUserDialog"
        title="新增用户"
        width="50%"
        :destroy-on-close="true"
    >
      <el-form :model="form1" label-width="100px" ref="userAddFormRef" :rules="rules">
        <el-form-item label="用户名：" prop="name">
          <el-input
              v-model="form1.name"
              autocomplete="off"
              placeholder="请输入用户名"
          />
        </el-form-item>

        <el-form-item label="提问次数：" prop="max_calls">
          <el-input
              v-model.number="form1.max_calls"
              autocomplete="off"
              placeholder="0 表示不限制提问次数"
          />
        </el-form-item>

        <el-form-item label="有效期：" prop="term">
          <el-input
              v-model.number="form1.term"
              autocomplete="off"
              placeholder="单位：天"
          />
        </el-form-item>

        <el-form-item label="聊天角色" prop="chat_roles">
          <el-select
              v-model="form1.chat_roles"
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
          <el-switch v-model="form1.enable_history"/>
        </el-form-item>
      </el-form>

      <template #footer>
            <span class="dialog-footer">
              <el-button @click="showUserDialog = false">取消</el-button>
              <el-button type="primary" @click="addUser">提交</el-button>
            </span>
      </template>
    </el-dialog>

    <el-dialog
        v-model="showBatchAddUserDialog"
        title="批量生成用户"
        width="50%"
        :destroy-on-close="true"

    >
      <el-form :model="form3" label-width="100px" ref="userEditFormRef" :rules="rules">
        <el-form-item label="提问次数：" prop="max_calls">
          <el-input
              v-model.number="form3.max_calls"
              autocomplete="off"
              placeholder="最大提问次数"
          />
        </el-form-item>

        <el-form-item label="用户数量：" prop="number">
          <el-input
              v-model.number="form3.number"
              autocomplete="off"
              placeholder="批量生成的用户数量"
          />
        </el-form-item>

        <el-form-item label="有效期：" prop="term">
          <el-input
              v-model.number="form3.term"
              autocomplete="off"
              placeholder="单位：天"
          />
        </el-form-item>

        <el-form-item label="聊天角色" prop="chat_roles">
          <el-select
              v-model="form3.chat_roles"
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
      </el-form>

      <template #footer>
            <span class="dialog-footer">
              <el-button @click="showBatchAddUserDialog = false">取消</el-button>
              <el-button type="success" @click="batchAddUser">提交</el-button>
            </span>
      </template>
    </el-dialog>

    <el-dialog
        v-model="showUserEditDialog"
        title="编辑用户"
        width="50%"
    >
      <el-form :model="form2" label-width="100px" ref="userEditFormRef" :rules="rules">
        <el-form-item label="用户名：" prop="name">
          <el-input
              v-model="form2.name"
              autocomplete="off"
              placeholder="请输入用户名"
              readonly
          />
        </el-form-item>

        <el-form-item label="提问次数：" prop="remaining_calls">
          <el-input
              v-model.number="form2.remaining_calls"
              autocomplete="off"
              placeholder="0"
          />
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
import {Plus} from "@element-plus/icons-vue";
import {httpPost} from "@/utils/http";
import {ElMessage} from "element-plus";
import {arrayContains, removeArrayItem} from "@/utils/libs";

// 变量定义
const users = ref([])
const form1 = ref({chat_roles: []})
const form2 = ref({})
const roles = ref([])
const showUserDialog = ref(false)
const showUserEditDialog = ref(false)
const rules = reactive({
  name: [{required: true, message: '请输入用户名', trigger: 'change',}],
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
    {required: true, message: '请输入用户数量', trigger: 'change'},
    {type: 'number', message: '请输入有效数字'},
  ],
  chat_roles: [{required: true, message: '请选择聊天角色', trigger: 'change'}],
})
const loading = ref(true)

const userAddFormRef = ref(null)
const userEditFormRef = ref(null)

onMounted(() => {
// 获取用户列表
  httpPost('/api/admin/user/list').then((res) => {
    users.value = res.data;
  }).catch(() => {
    ElMessage.error('获取系统配置失败')
  })

  // 获取角色列表
  httpPost('/api/admin/chat-roles/list').then((res) => {
    roles.value = res.data;
    roles.value.unshift({name: '全部', key: 'all'})
  }).catch(() => {
    ElMessage.error("获取聊天角色失败");
  })

  nextTick(() => {
    loading.value = false
  })
})

// 新增用户
const addUser = () => {
  userAddFormRef.value.validate((valid) => {
    if (valid) {
      showUserDialog.value = false
      form1.value.term = parseInt(form1.value.term)
      form1.value.max_calls = parseInt(form1.value.max_calls)
      httpPost('/api/admin/user/add', form1.value).then((res) => {
        ElMessage.success('添加用户成功')
        form1.value = {chat_roles: []}
        users.value.unshift(res.data)
      }).catch((e) => {
        ElMessage.error('添加用户失败，' + e.message)
      })
    } else {
      return false
    }
  })
}

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
    ElMessage.success('删除用户成功')
    users.value = removeArrayItem(users.value, user, function (v1, v2) {
      return v1.name === v2.name
    })
  }).catch((e) => {
    ElMessage.error('删除用户失败，' + e.message)
  })
}

const userEdit = function (user) {
  form2.value = user
  showUserEditDialog.value = true
}

// 更新用户
const updateUser = function () {
  userEditFormRef.value.validate((valid) => {
    if (valid) {
      showUserEditDialog.value = false
      form2.value.term = parseInt(form2.value.term)
      form2.value.remaining_calls = parseInt(form2.value.remaining_calls)
      httpPost('/api/admin/user/set', form2.value).then(() => {
        ElMessage.success('更新用户成功')
      }).catch((e) => {
        ElMessage.error('更新用户失败，' + e.message)
      })
    } else {
      return false
    }
  })
}

// 批量新增
const showBatchAddUserDialog = ref(false)
const form3 = ref({chat_roles: []})
const batchAddUser = function () {
  httpPost('api/admin/user/batch-add', form3.value).then((res) => {
    console.log(res.data)
    ElMessage.success('添加用户成功')
    users.value = [...res.data, ...users.value]
    showBatchAddUserDialog.value = false
  }).catch((e) => {
    console.log('添加用户失败，' + e.message)
  })
}
</script>

<style lang="stylus" scoped>
.user-list {

  .opt-box {
    padding-bottom: 10px;

    .el-icon {
      margin-right 5px;
    }
  }

  .el-select {
    width 100%
  }

}
</style>