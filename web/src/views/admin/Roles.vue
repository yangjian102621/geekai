<template>
  <div class="container role-list" v-loading="loading">
    <div class="handle-box">
      <el-button type="primary" :icon="Plus" @click="addRole">新增</el-button>
    </div>
    <el-row>
      <el-table :data="tableData" :border="parentBorder" style="width: 100%">
        <el-table-column type="expand">
          <template #default="props">
            <div>
              <el-table :data="props.row.context" :border="childBorder">
                <el-table-column label="对话角色" prop="role" width="120"/>
                <el-table-column label="对话内容" prop="content"/>
              </el-table>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="角色名称" prop="name">
          <template #default="scope">
            <span class="sort" :data-id="scope.row.id">{{ scope.row.name }}</span>
          </template>
        </el-table-column>
        <el-table-column label="角色标识" prop="key"/>
        <el-table-column label="启用状态">
          <template #default="scope">
            <el-switch v-model="scope.row['enable']" @change="roleSet('enable',scope.row)"/>
          </template>
        </el-table-column>
        <el-table-column label="角色图标" prop="icon">
          <template #default="scope">
            <el-image :src="scope.row.icon" style="width: 45px; height: 45px; border-radius: 50%"/>
          </template>
        </el-table-column>
        <el-table-column label="打招呼信息" prop="hello_msg"/>
        <el-table-column label="操作" width="150" align="right">
          <template #default="scope">
            <el-button size="small" type="primary" @click="rowEdit(scope.$index, scope.row)">编辑</el-button>
            <el-popconfirm title="确定要删除当前角色吗?" @confirm="removeRole(scope.row)" :width="200">
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
        title="编辑角色"
        :close-on-click-modal="false"
        width="50%"
    >
      <el-form :model="role" label-width="120px" ref="formRef" label-position="left" :rules="rules">
        <el-form-item label="角色名称：" prop="name">
          <el-input
              v-model="role.name"
              autocomplete="off"
          />
        </el-form-item>

        <el-form-item label="角色标志：" prop="key">
          <el-input
              v-model="role.key"
              autocomplete="off"
          />
        </el-form-item>

        <el-form-item label="角色图标：" prop="icon">
          <el-input
              v-model="role.icon"
              autocomplete="off"
          />
        </el-form-item>

        <el-form-item label="打招呼信息：" prop="hello_msg">
          <el-input
              v-model="role.hello_msg"
              autocomplete="off"
          />
        </el-form-item>

        <el-form-item label="上下文信息：" prop="context">
          <template #default>
            <el-table :data="role.context" :border="childBorder" size="small">
              <el-table-column label="对话角色" width="120">
                <template #default="scope">
                  <el-input
                      v-model="scope.row.role"
                      autocomplete="off"
                  />
                </template>
              </el-table-column>
              <el-table-column label="对话内容">
                <template #header>
                  <div class="context-msg-key">
                    <span>对话内容</span>
                    <span class="fr">
                      <el-button type="primary" @click="addContext" size="small">
                      <el-icon>
                        <Plus/>
                      </el-icon>
                      增加一行
                    </el-button>
                    </span>
                  </div>
                </template>

                <template #default="scope">
                  <div class="context-msg-content">
                    <el-input
                        v-model="scope.row.content"
                        autocomplete="off"
                    />
                    <span><el-icon @click="removeContext(scope.$index)"><RemoveFilled/></el-icon></span>
                  </div>
                </template>
              </el-table-column>
            </el-table>
          </template>
        </el-form-item>

        <el-form-item label="启用状态">
          <el-switch v-model="role.enable"/>
        </el-form-item>
      </el-form>

      <template #footer>
            <span class="dialog-footer">
              <el-button @click="showDialog = false">取消</el-button>
              <el-button type="primary" @click="save">保存</el-button>
            </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>

import {Plus, RemoveFilled} from "@element-plus/icons-vue";
import {onMounted, reactive, ref} from "vue";
import {httpGet, httpPost} from "@/utils/http";
import {ElMessage} from "element-plus";
import {copyObj, removeArrayItem} from "@/utils/libs";
import {Sortable} from "sortablejs"

const showDialog = ref(false)
const parentBorder = ref(true)
const childBorder = ref(true)
const tableData = ref([])
const sortedTableData = ref([])
const role = ref({context: []})
const formRef = ref(null)
const editRow = ref({})
const loading = ref(true)

const rules = reactive({
  name: [{required: true, message: '请输入用户名', trigger: 'blur',}],
  key: [{required: true, message: '请输入角色标识', trigger: 'blur',}],
  icon: [{required: true, message: '请输入角色图标', trigger: 'blur',}],
  sort: [
    {required: true, message: '请输入排序数字', trigger: 'blur'},
    {type: 'number', message: '请输入有效数字'},
  ],
  hello_msg: [{required: true, message: '请输入打招呼信息', trigger: 'change',}]
})

// 获取角色列表
httpGet('/api/admin/role/list').then((res) => {
  tableData.value = res.data
  sortedTableData.value = copyObj(tableData.value)
  loading.value = false
}).catch(() => {
  ElMessage.error("获取聊天角色失败");
})

onMounted(() => {
  const drawBodyWrapper = document.querySelector('.el-table__body tbody')

  // 初始化拖动排序插件
  Sortable.create(drawBodyWrapper, {
    sort: true,
    animation: 500,
    onEnd({newIndex, oldIndex, from}) {
      if (oldIndex === newIndex) {
        return
      }

      const sortedData = Array.from(from.children).map(row => row.querySelector('.sort').getAttribute('data-id'));
      const ids = []
      const sorts = []
      sortedData.forEach((id, index) => {
        ids.push(parseInt(id))
        sorts.push(index)
      })

      httpPost("/api/admin/role/sort", {ids: ids, sorts: sorts}).catch(e => {
        ElMessage.error("排序失败：" + e.message)
      })
    }
  })
})

const roleSet = (filed, row) => {
  httpPost('/api/admin/role/set', {id: row.id, filed: filed, value: row[filed]}).then(() => {
    ElMessage.success("操作成功！")
  }).catch(e => {
    ElMessage.error("操作失败：" + e.message)
  })
}

// 编辑
const curIndex = ref(0)
const rowEdit = function (index, row) {
  curIndex.value = index
  role.value = copyObj(row)
  showDialog.value = true
}

const addRole = function () {
  role.value = {context: []}
  showDialog.value = true
}

const save = function () {
  formRef.value.validate((valid) => {
    if (valid) {
      showDialog.value = false
      httpPost('/api/admin/role/save', role.value).then((res) => {
        ElMessage.success('操作成功')
        // 更新当前数据行
        if (role.value.id) {
          tableData.value[curIndex.value] = role.value
        } else {
          tableData.value.push(res.data)
        }
      }).catch((e) => {
        ElMessage.error('操作失败，' + e.message)
      })
    }
  })
}

const removeRole = function (row) {
  httpGet('/api/admin/role/remove?id=' + row.id).then(() => {
    ElMessage.success("删除成功！")
    tableData.value = removeArrayItem(tableData.value, row, (v1, v2) => {
      return v1.id === v2.id
    })
  }).catch(() => {
    ElMessage.error("删除失败！")
  })
}

const addContext = function () {
  if (!role.value.context) {
    role.value.context = []
  }
  role.value.context.push({role: '', content: ''})
}

const removeContext = function (index) {
  role.value.context.splice(index, 1);
}

</script>

<style lang="stylus" scoped>
.role-list {
  .opt-box {
    padding-bottom: 10px;
    display flex;
    justify-content flex-end

    .el-icon {
      margin-right 5px;
    }
  }

  .context-msg-key {
    .fr {
      float right

      .el-icon {
        margin-right 5px
      }
    }
  }

  .context-msg-content {
    display flex

    .el-icon {
      font-size: 20px;
      margin-top 5px;
      margin-left 5px;
      cursor pointer
    }
  }

  .el-input--small {
    width 30px;

    .el-input__inner {
      text-align center
    }
  }
}
</style>