<template>
  <div class="container list" v-loading="loading">

    <div class="handle-box">
      <el-button type="primary" :icon="Plus" @click="add">新增</el-button>
    </div>

    <el-row>
      <el-table :data="items" :row-key="row => row.id" table-layout="auto">
        <el-table-column prop="username" label="用户名"/>
        <el-table-column prop="last_login_ip" label="最后登录IP"/>

        <el-table-column label="最后登录时间">
          <template #default="scope">
            <span>{{ dateFormat(scope.row['last_login_at']) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="enabled" label="启用状态">
          <template #default="scope">
            <el-switch v-model="scope.row['status']" @change="enable(scope.row)"/>
          </template>
        </el-table-column>
        <el-table-column label="创建时间">
          <template #default="scope">
            <span>{{ dateFormat(scope.row['created_at']) }}</span>
          </template>
        </el-table-column>

        <el-table-column label="操作" width="180">
          <template #default="scope">
            <el-button size="small" type="primary" @click="resetPass(scope.row)">重置密码</el-button>
            <el-popconfirm title="确定要删除当前记录吗?" @confirm="remove(scope.row)" :width="200">
              <template #reference>
                <el-button size="small" type="danger">删除</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>
    </el-row>

    <el-dialog
        v-model="showDialog"
        title="添加用户"
        :close-on-click-modal="false"
    >
      <el-form :model="item" label-width="120px" ref="formRef" :rules="rules">
        <el-form-item label="用户名：" prop="username">
          <el-input v-model="item.username" autocomplete="off"/>
        </el-form-item>

        <el-form-item label="密码：" prop="password">
          <el-input v-model="item.password" type="password" autocomplete="off"/>
        </el-form-item>

        <el-form-item label="重复密码：" prop="repass">
          <el-input v-model="item.repass" type="password" autocomplete="off"/>
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
import {onMounted, reactive, ref} from "vue";
import {httpGet, httpPost} from "@/utils/http";
import {ElMessage, ElMessageBox} from "element-plus";
import {dateFormat, removeArrayItem} from "@/utils/libs";
import {Plus} from "@element-plus/icons-vue";
import {Sortable} from "sortablejs";

// 变量定义
const items = ref([])
const item = ref({})
const showDialog = ref(false)
const title = ref("")
const loading = ref(true)
const formRef = ref(null)

const rules = reactive({
  username: [{required: true, message: '请输入用户名', trigger: 'change',}],
  password: [{required: true, message: '请输入密码', trigger: 'change',}],
  repass: [{required: true, message: '请再次输入密码', trigger: 'change',}],
})

// 获取数据
const fetchData = () => {
  httpGet('/api/admin/list').then((res) => {
    items.value = res.data
    loading.value = false
  }).catch(() => {
    ElMessage.error("获取数据失败");
  })
}

onMounted(() => {
  fetchData()
})

const add = function () {
  showDialog.value = true
  item.value = {}
}

// 重置密码
const resetPass = function (row) {
  ElMessageBox.prompt('请输入新密码', '重置密码', {
    confirmButtonText: '确认',
    cancelButtonText: '取消',
  }).then(({value}) => {
    httpPost("/api/admin/resetPass", {
      id: row.id,
      password: value
    }).then(() => {
      ElMessage.success("操作成功")
    }).catch(e => {
      ElMessage.error("操作失败：" + e.message)
    })

  }).catch(() => {
  })
}

const save = function () {
  formRef.value.validate((valid) => {
    if (item.value.password !== item.value.repass) {
      return ElMessage.error("两次输入密码不一致！")
    }
    if (valid) {
      showDialog.value = false
      httpPost('/api/admin/save', item.value).then((res) => {
        ElMessage.success('操作成功！')
        fetchData()
      }).catch((e) => {
        ElMessage.error('操作失败，' + e.message)
      })
    } else {
      return false
    }
  })
}

const enable = (row) => {
  httpPost('/api/admin/enable', {id: row.id, enabled: row.status}).then(() => {
    ElMessage.success("操作成功！")
  }).catch(e => {
    ElMessage.error("操作失败：" + e.message)
  })
}

const remove = function (row) {
  httpGet('/api/admin/remove?id=' + row.id).then(() => {
    ElMessage.success("删除成功！")
    items.value = removeArrayItem(items.value, row, (v1, v2) => {
      return v1.id === v2.id
    })
  }).catch((e) => {
    ElMessage.error("删除失败：" + e.message)
  })
}
</script>

<style lang="stylus" scoped>
.list {

  .opt-box {
    padding-bottom: 10px;
    display flex;
    justify-content flex-end

    .el-icon {
      margin-right: 5px;
    }
  }

  .el-select {
    width: 100%
  }

}
</style>