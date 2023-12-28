<template>
  <div class="container role-list" v-loading="loading">
    <div class="handle-box">
      <el-button type="primary" :icon="Plus" @click="addRow">新增</el-button>
    </div>
    <el-row>
      <el-table :data="tableData" :border="parentBorder" style="width: 100%">
        <el-table-column label="函数名称" prop="name">
          <template #default="scope">
            <span class="sort" :data-id="scope.row.id">{{ scope.row.name }}</span>
          </template>
        </el-table-column>
        <el-table-column label="函数别名" prop="label"/>
        <el-table-column label="功能描述" prop="description"/>
        <el-table-column label="启用状态">
          <template #default="scope">
            <el-switch v-model="scope.row.enabled" @change="functionSet('enabled',scope.row)"/>
          </template>
        </el-table-column>

        <el-table-column label="操作" width="150" align="right">
          <template #default="scope">
            <el-button size="small" type="primary" @click="rowEdit(scope.$index, scope.row)">编辑</el-button>
            <el-popconfirm title="确定要删除当前函数吗?" @confirm="remove(scope.row)" :width="200">
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
              placeholder="函数名称最好为英文"
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

        <el-form-item label="API 地址：" prop="action">
          <el-input
              v-model="item.action"
              autocomplete="off"
              placeholder="该函数实现的API地址，可以是第三方服务API"
          />
        </el-form-item>

        <el-form-item label="API Token：" prop="token">
          <el-input
              v-model="item.token"
              autocomplete="off"
              placeholder="API授权Token"
          >
            <template #append>
              <el-tooltip
                  class="box-item"
                  effect="dark"
                  content="只有本地服务才可以使用自动生成Token<br/>第三方服务请填写第三方服务API Token"
                  placement="top-end"
                  raw-content
              >
                <el-button @click="generateToken">生成Token</el-button>
              </el-tooltip>
            </template>
          </el-input>
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

import {Delete, Plus} from "@element-plus/icons-vue";
import {onMounted, reactive, ref} from "vue";
import {httpGet, httpPost} from "@/utils/http";
import {ElMessage} from "element-plus";
import {arrayContains, copyObj} from "@/utils/libs";

const showDialog = ref(false)
const parentBorder = ref(true)
const childBorder = ref(true)
const tableData = ref([])
const item = ref({})
const params = ref([])
const formRef = ref(null)
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
    if (res.data) {
      tableData.value = res.data
    }
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
  // initialize parameters
  const props = item.value?.parameters?.properties
  const required = item.value?.parameters?.required
  const _params = []
  for (let key in props) {
    _params.push({
      name: key,
      type: props[key].type,
      desc: props[key].description,
      required: arrayContains(required, key)
    })
  }
  params.value = _params
  showDialog.value = true
}

const addRow = function () {
  item.value = {enabled: true}
  params.value = []
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
      item.value.parameters = {type: "object", "properties": properties, required: required}
      httpPost('/api/admin/function/save', item.value).then((res) => {
        ElMessage.success('操作成功')
        console.log(res.data)
        if (item.value.id > 0) {
          tableData.value[curIndex.value] = item.value
        } else {
          tableData.value.push(res.data)
        }
      }).catch((e) => {
        ElMessage.error('操作失败，' + e.message)
      })
    }
  })
}

const remove = function (row) {
  httpGet('/api/admin/function/remove?id=' + row.id).then(() => {
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
  params.value.push({name: "", type: "string", desc: "", required: false})
}

const removeParam = function (index) {
  params.value.splice(index, 1);
}

const functionSet = (filed, row) => {
  httpPost('/api/admin/function/set', {id: row.id, filed: filed, value: row[filed]}).then(() => {
    ElMessage.success("操作成功！")
  }).catch(e => {
    ElMessage.error("操作失败：" + e.message)
  })
}

const generateToken = () => {
  httpGet('/api/admin/function/token').then(res => {
    item.value.token = res.data
  }).catch(() => {
    ElMessage.error("生成 Token 失败")
  })
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