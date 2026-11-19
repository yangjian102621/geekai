<template>
  <div class="container role-list" v-loading="loading">
    <div class="handle-box">
      <el-button type="primary" :icon="Plus" @click="addRole">新增</el-button>
    </div>
    <el-row>
      <el-table :data="tableData" border>
        <el-table-column label="应用名称" prop="name">
          <template #default="scope">
            <span class="sort" :data-id="scope.row.id">
              <i class="iconfont icon-drag"></i>
              {{ scope.row.name }}
            </span>
          </template>
        </el-table-column>
        <el-table-column label="应用类型" prop="type_name" />
        <el-table-column label="绑定模型" prop="model_name" />
        <el-table-column label="启用状态">
          <template #default="scope">
            <el-switch v-model="scope.row['enable']" @change="roleSet('enable', scope.row)" />
          </template>
        </el-table-column>
        <el-table-column label="应用图标" prop="icon">
          <template #default="scope">
            <el-image :src="scope.row.icon" style="width: 45px; height: 45px; border-radius: 50%" />
          </template>
        </el-table-column>
        <el-table-column label="打招呼信息" prop="hello_msg" />
        <el-table-column label="操作" width="150">
          <template #default="scope">
            <el-button size="small" type="primary" @click="rowEdit(scope.$index, scope.row)"
              >编辑</el-button
            >
            <el-popconfirm
              title="确定要删除当前应用吗?"
              @confirm="removeRole(scope.row)"
              :width="200"
            >
              <template #reference>
                <el-button size="small" type="danger">删除</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>
    </el-row>

    <el-dialog v-model="showDialog" :title="optTitle" :close-on-click-modal="false" width="50%">
      <el-form :model="role" label-width="120px" ref="formRef" label-position="left" :rules="rules">
        <el-form-item label="应用名称：" prop="name">
          <el-input v-model="role.name" autocomplete="off" />
        </el-form-item>
        <el-form-item label="应用分类：" prop="tid">
          <el-select v-model="role.tid" filterable placeholder="请选择分类" clearable>
            <el-option
              v-for="item in appTypes"
              :key="item.id"
              :label="item.name"
              :value="item.id"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="应用图标：" prop="icon">
          <el-input v-model="role.icon">
            <template #append>
              <el-upload :auto-upload="true" :show-file-list="false" :http-request="uploadImg">
                上传
              </el-upload>
            </template>
          </el-input>
        </el-form-item>

        <el-form-item label="绑定模型：" prop="model_id">
          <el-select v-model="role.model_id" filterable placeholder="请选择模型" clearable>
            <el-option v-for="item in models" :key="item.id" :label="item.name" :value="item.id" />
          </el-select>
        </el-form-item>

        <el-form-item label="打招呼信息：" prop="hello_msg">
          <el-input v-model="role.hello_msg" autocomplete="off" />
        </el-form-item>

        <el-form-item label="预设提示词：" prop="system_prompt">
          <el-input
            v-model="role.system_prompt"
            type="textarea"
            :rows="6"
            placeholder="请输入系统预设提示词，将作为该应用的系统指令"
          />
        </el-form-item>

        <el-form-item label="启用状态">
          <el-switch v-model="role.enable" />
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
import { httpGet, httpPost } from '@/utils/http'
import { copyObj, removeArrayItem } from '@/utils/libs'
import { Plus } from '@element-plus/icons-vue'
import Compressor from 'compressorjs'
import { ElMessage } from 'element-plus'
import { Sortable } from 'sortablejs'
import { onMounted, reactive, ref } from 'vue'

const showDialog = ref(false)
const parentBorder = ref(true)
const childBorder = ref(true)
const tableData = ref([])
const sortedTableData = ref([])
const role = ref({ system_prompt: '' })
const formRef = ref(null)
const optTitle = ref('')
const loading = ref(true)

const rules = reactive({
  name: [{ required: true, message: '请输入应用名称', trigger: 'blur' }],
  icon: [{ required: true, message: '请输入应用图标', trigger: 'blur' }],
  sort: [
    { required: true, message: '请输入排序数字', trigger: 'blur' },
    { type: 'number', message: '请输入有效数字' },
  ],
  hello_msg: [{ required: true, message: '请输入打招呼信息', trigger: 'change' }],
})

const appTypes = ref([])
const models = ref([])
onMounted(() => {
  fetchData()

  // get chat models
  httpGet('/api/admin/model/list?enable=1')
    .then((res) => {
      models.value = res.data
    })
    .catch(() => {
      ElMessage.error('获取AI模型数据失败')
    })

  // get app type
  httpGet('/api/admin/app/type/list?enable=1')
    .then((res) => {
      appTypes.value = res.data
    })
    .catch(() => {
      ElMessage.error('获取应用分类数据失败')
    })
})

const fetchData = () => {
  // 获取应用列表
  httpGet('/api/admin/role/list')
    .then((res) => {
      // 初始化数据
      // const arr = res.data;
      // for (let i = 0; i < arr.length; i++) {
      //   if(arr[i].model_id == 0){
      //     arr[i].model_id = ''
      //   }
      // }
      tableData.value = res.data
      sortedTableData.value = copyObj(tableData.value)
      loading.value = false
    })
    .catch(() => {
      ElMessage.error('获取聊天应用失败')
    })

  const drawBodyWrapper = document.querySelector('.el-table__body tbody')
  // 初始化拖动排序插件
  Sortable.create(drawBodyWrapper, {
    sort: true,
    animation: 500,
    onEnd({ newIndex, oldIndex, from }) {
      if (oldIndex === newIndex) {
        return
      }

      const sortedData = Array.from(from.children).map((row) =>
        row.querySelector('.sort').getAttribute('data-id')
      )
      const ids = []
      const sorts = []
      sortedData.forEach((id, index) => {
        ids.push(parseInt(id))
        sorts.push(index + 1)
        tableData.value[index].sort_num = index + 1
      })

      httpPost('/api/admin/role/sort', { ids: ids, sorts: sorts }).catch((e) => {
        ElMessage.error('排序失败：' + e.message)
      })
    },
  })
}

const roleSet = (filed, row) => {
  httpPost('/api/admin/role/set', {
    id: row.id,
    filed: filed,
    value: row[filed],
  })
    .then(() => {
      ElMessage.success('操作成功！')
    })
    .catch((e) => {
      ElMessage.error('操作失败：' + e.message)
    })
}

// 编辑
const curIndex = ref(0)
const rowEdit = function (index, row) {
  optTitle.value = '修改应用'
  curIndex.value = index
  role.value = copyObj(row)
  showDialog.value = true
}

const addRole = function () {
  optTitle.value = '添加新应用'
  role.value = { system_prompt: '', enable: true }
  showDialog.value = true
}

const save = function () {
  formRef.value.validate((valid) => {
    if (valid) {
      showDialog.value = false
      httpPost('/api/admin/role/save', role.value)
        .then(() => {
          ElMessage.success('操作成功')
          fetchData()
        })
        .catch((e) => {
          ElMessage.error('操作失败，' + e.message)
        })
    }
  })
}

const removeRole = function (row) {
  httpGet('/api/admin/role/remove?id=' + row.id)
    .then(() => {
      ElMessage.success('删除成功！')
      tableData.value = removeArrayItem(tableData.value, row, (v1, v2) => {
        return v1.id === v2.id
      })
    })
    .catch(() => {
      ElMessage.error('删除失败！')
    })
}

// 图片上传
const uploadImg = (file) => {
  // 压缩图片并上传
  new Compressor(file.file, {
    quality: 0.6,
    success(result) {
      const formData = new FormData()
      formData.append('file', result, result.name)
      // 执行上传操作
      httpPost('/api/admin/upload', formData)
        .then((res) => {
          role.value.icon = res.data.url
          ElMessage.success('上传成功')
        })
        .catch((e) => {
          ElMessage.error('上传失败:' + e.message)
        })
    },
    error(e) {
      ElMessage.error('上传失败:' + e.message)
    },
  })
}
</script>

<style lang="scss" scoped>
.role-list {
  .handle-box {
    margin-bottom: 20px;

    .handle-input {
      max-width: 150px;
      margin-right: 10px;
    }
  }

  .opt-box {
    padding-bottom: 10px;
    display: flex;
    justify-content: flex-end;

    .el-icon {
      margin-right: 5px;
    }
  }

  .el-input--small {
    width: 30px;

    .el-input__inner {
      text-align: center;
    }
  }

  .sort {
    cursor: move;
    .iconfont {
      position: relative;
      top: 1px;
    }
  }

  .pagination {
    padding: 20px 0;
    display: flex;
    justify-content: right;
  }
}
</style>
