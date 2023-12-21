<template>
  <div class="container role-list" v-loading="loading">
    <div class="handle-box">
      <el-button type="primary" :icon="Plus" @click="addRole">新增</el-button>
    </div>
    <el-row>
      <el-table :data="tableData" :border="parentBorder" style="width: 100%">
        <el-table-column label="函数名称" prop="name">
          <template #default="scope">
            <span class="sort" :data-id="scope.row.id">{{ scope.row.name }}</span>
          </template>
        </el-table-column>
        <el-table-column label="功能描述" prop="key"/>
        <el-table-column label="">
          <template #default="scope">
            <el-tag v-if="scope.row.enable" type="success">启用</el-tag>
            <el-tag type="danger" v-else>禁用</el-tag>
          </template>
        </el-table-column>

        <el-table-column label="操作" width="150" align="right">
          <template #default="scope">
            <el-button size="small" type="primary" @click="rowEdit(scope.$index, scope.row)">编辑</el-button>
            <el-popconfirm title="确定要删除当前函数吗?" @confirm="remove(scope.row)">
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
        :title="title"
        width="50%"
    >
      <el-form :model="item" label-width="120px" ref="formRef" label-position="left" :rules="rules">
        <el-form-item label="函数名称：" prop="name">
          <el-input
              v-model="item.name"
              autocomplete="off"
          />
        </el-form-item>

        <el-form-item label="函数标签：" prop="label">
          <el-input
              v-model="item.label"
              placeholder="函数的中文名称"
              autocomplete="off"
          />
        </el-form-item>

        <el-form-item label="功能描述：" prop="description">
          <el-input
              v-model="item.description"
              autocomplete="off"
          />
        </el-form-item>

        <el-form-item label="函数参数：" prop="parameters">
          <template #default>
            <el-table :data="params" :border="childBorder" size="small">
              <el-table-column label="参数名称" width="120">
                <template #default="scope">
                  <el-input
                      v-model="scope.row.name"
                      autocomplete="off"
                  />
                </template>
              </el-table-column>
              <el-table-column label="参数类型" width="120">
                <template #default="scope">
                  <el-select v-model="scope.row.type" placeholder="参数类型">
                    <el-option v-for="pt in paramsType" :value="pt" :key="pt">{{ pt }}</el-option>
                  </el-select>
                </template>
              </el-table-column>
              <el-table-column label="参数描述">
                <template #default="scope">
                  <el-input
                      v-model="scope.row.desc"
                      autocomplete="off"
                  />
                </template>
              </el-table-column>

              <el-table-column label="必填参数" width="80">
                <template #default="scope">
                  <div class="param-opt">
                    <el-checkbox v-model="scope.row.required"/>
                  </div>
                </template>
              </el-table-column>

              <el-table-column label="操作" width="80">
                <template #default="scope">
                  <div class="param-opt">
                    <el-button type="danger" :icon="Delete" circle @click="removeParam(scope.$index)" size="small"/>
                  </div>
                </template>
              </el-table-column>
            </el-table>

            <div class="param-line">
              <el-button type="primary" @click="addParam" size="small">
                <el-icon>
                  <Plus/>
                </el-icon>
                增加参数
              </el-button>
            </div>
          </template>
        </el-form-item>

        <el-form-item label="启用状态">
          <el-switch v-model="item.enabled"/>
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

import {Delete, Plus, RemoveFilled} from "@element-plus/icons-vue";
import {onMounted, reactive, ref} from "vue";
import {httpGet, httpPost} from "@/utils/http";
import {ElMessage} from "element-plus";
import {copyObj, removeArrayItem} from "@/utils/libs";
import {Sortable} from "sortablejs"

const showDialog = ref(false)
const parentBorder = ref(true)
const childBorder = ref(true)
const tableData = ref([])
const item = ref({parameters: []})
const params = ref([])
const formRef = ref(null)
const editRow = ref({})
const loading = ref(true)
const title = ref("新增函数")

const rules = reactive({
  name: [{required: true, message: '请输入函数名称', trigger: 'blur',}],
  label: [{required: true, message: '请输入函数标签', trigger: 'blur',}],
  description: [{required: true, message: '请输入函数功能描述', trigger: 'blur',}],
})
const paramsType = ref(["string", "number"])


onMounted(() => {
  fetch()
})

const fetch = () => {
  httpGet('/api/admin/function/list').then((res) => {
    tableData.value = res.data
    loading.value = false
  }).catch(() => {
    ElMessage.error("获取数据失败");
  })
}

// 编辑
const curIndex = ref(0)
const rowEdit = function (index, row) {
  curIndex.value = index
  item.value = copyObj(row)
  showDialog.value = true
}

const addRole = function () {
  item.value = {parameters: []}
  showDialog.value = true
}

const save = function () {
  formRef.value.validate((valid) => {
    if (valid) {
      showDialog.value = false
      const properties = {}
      const required = []
      // process params
      for (let i = 0; i < params.value.length; i++) {
        properties[params.value[i].name] = {"type": params.value[i].type, "description": params.value[i].desc}
        if (params.value[i].required) {
          required.push(params.value[i].name)
        }
      }
      item.value.parameters = {type: "object", "properties": properties, "required": required}
      httpPost('/api/admin/function/save', item.value).then((res) => {
        ElMessage.success('操作成功')
        // 更新当前数据行
        // if (item.value.id) {
        //   tableData.value[curIndex.value] = item.value
        // } else {
        //   tableData.value.push(res.data)
        // }
      }).catch((e) => {
        ElMessage.error('操作失败，' + e.message)
      })
    }
  })
}

const remove = function (row) {
  httpGet('/api/admin/role/remove?id=' + row.id).then(() => {
    ElMessage.success("删除成功！")
    fetch()
  }).catch(() => {
    ElMessage.error("删除失败！")
  })
}

const addParam = function () {
  if (!params.value) {
    item.value = []
  }
  params.value.push({name: "", type: "", desc: "", required: false})
}

const removeParam = function (index) {
  params.value.splice(index, 1);
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

  .param-line {
    padding 5px 0

    .el-icon {
      margin-right 5px;
    }
  }

  .param-opt {
    display flex
    justify-content center

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