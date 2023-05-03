<template>
  <div class="system-config" v-loading="loading">
    <el-row class="opt-box">
      <el-button type="primary" @click="showUserDialog = true">
        <el-icon><Plus /></el-icon> 新增
      </el-button>

      <el-button type="success">
        <el-icon><Plus /></el-icon> 批量新增
      </el-button>
    </el-row>

    <el-row>
      <el-table :data="users">
        <el-table-column prop="name" label="用户名" />
        <el-table-column prop="max_calls" label="最大调用次数" />
        <el-table-column prop="remaining_calls" label="剩余调用次数" />
        <el-table-column label="激活时间">
          <template #default="scope">
            <el-tag type="info" v-if="scope.row.active_time === 0">未激活</el-tag>
            <span v-else>{{dateFormat(scope.row.active_time)}}</span>
          </template>
        </el-table-column>
        <el-table-column label="过期时间">
          <template #default="scope">
            <el-tag type="info" v-if="scope.row.active_time === 0">未激活</el-tag>
            <span v-else>{{dateFormat(scope.row.expired_time)}}</span>
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
            <el-button size="small" type="primary" @click="removeApiKey(scope.row.value)">编辑</el-button>
            <el-button size="small" type="danger" @click="removeApiKey(scope.row.value)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-row>

    <el-dialog
        v-model="showUserDialog"
        title="新增用户"
        width="30%"
        :destroy-on-close="true"
    >
      <el-form :model="user" label-width="100px" ref="userAddFormRef" :rules="rules">
        <el-form-item label="用户名：" prop="name">
          <el-input
              v-model="user.name"
              autocomplete="off"
              placeholder="请输入用户名"
          />
        </el-form-item>

        <el-form-item label="调用次数：">
          <el-input
              v-model="user.max_calls"
              autocomplete="off"
              placeholder="0 表示不限制调用次数"
          />
        </el-form-item>

        <el-form-item label="有效期：">
          <el-input
              v-model="user.term"
              autocomplete="off"
              placeholder="单位：天"
          />
        </el-form-item>

        <el-form-item label="聊天角色">
          <el-select
              v-model="user.chat_roles"
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
          <el-switch v-model="user.enable_history" />
        </el-form-item>
      </el-form>

      <template #footer>
            <span class="dialog-footer">
              <el-button @click="showUserDialog = false">取消</el-button>
              <el-button type="primary" @click="saveUser">提交</el-button>
            </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import {onMounted} from "vue";

const userAddFormRef = ref(null);
onMounted(() => {
  alert('xxxx')
})

const saveUser = ()=> {
  userAddFormRef.value.validate((valid) => {
    if (valid) {

        this.showUserDialog = false
        this.user.term = parseInt(this.user.term)
        this.user.max_calls = parseInt(this.user.max_calls)
        httpPost('/api/admin/user/add', this.user).then((res) => {
          ElMessage.success('添加用户成功')
          this.user = {}
          this.users.push(res.data)
        }).catch((e) => {
          ElMessage.error('添加用户失败，'+e.message)
        })

    } else {
      ElMessage.error('error submit !')
      return false
    }
  })
}
</script>

<script>
import {defineComponent, nextTick, reactive, ref} from "vue";
import {Plus} from "@element-plus/icons-vue";
import {httpPost} from "@/utils/http";
import {ElMessage} from "element-plus";
import {arrayContains, dateFormat} from "@/utils/libs";

export default defineComponent({
  name: 'UserList',
  components: {Plus},
  data() {
    return {
      title: "用户管理",
      users: [],
      user: {
        term: 30,
        enable_history: true
      },
      roles: [],
      showUserDialog: false,
      rules: reactive({
        name: [
          {
            required: true,
            message: 'Please select Activity zone',
            trigger: 'change',
          },
        ]
      }),
      loading: true
    }
  },
  mounted() {
    // 获取用户列表
    httpPost('/api/admin/user/list').then((res) => {
      this.users = res.data;
    }).catch(() => {
      ElMessage.error('获取系统配置失败')
    })

    // 获取角色列表
    httpPost("/api/admin/chat-roles/get").then((res) => {
      this.roles = res.data;
      this.roles.unshift({name:'全部',key:'all'})
    }).catch(() => {
      ElMessage.error("获取聊天角色失败");
    })

    nextTick(() => {
      this.loading = false
    })
  },
  computed: {
    dateFormat() {
      return dateFormat
    },
  },
  methods: {
    selectRole: function (items) {
      if (arrayContains(items, 'all')) {
        for (let i =0;i < this.roles.length;i++) {
          if (this.roles[i].key === 'all') {
            continue
          }
          this.roles[i].disabled = true
          this.user.chat_roles = ['all']
        }
      } else {
        for (let i =0;i < this.roles.length;i++) {
          if (this.roles[i].key === 'all') {
            continue
          }
          this.roles[i].disabled = false
        }
      }

    },

    removeUser: function (user) {
      console.log(user)
    },

    updateUser:function (user) {
      console.log(user)
    }
  }
})
</script>

<style lang="stylus" scoped>
.system-config {

  .opt-box {
    padding-bottom: 10px;

    .el-icon {
      margin-right 5px;
    }
  }

}
</style>