<template>
  <div class="container menu" v-loading="loading">

    <div class="handle-box">
      <el-button type="primary" :icon="Plus" @click="add">新增</el-button>
    </div>

    <el-row>
      <el-table :data="items" :row-key="row => row.id" table-layout="auto">
        <el-table-column prop="name" label="菜单名称">
          <template #default="scope">
            <span class="sort" :data-id="scope.row.id">{{ scope.row.name }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="icon" label="菜单图标">
          <template #default="scope">
            <el-image class="menu-icon" :src="scope.row.icon"/>
          </template>
        </el-table-column>
        <el-table-column prop="url" label="菜单URL"/>
        <el-table-column prop="enabled" label="启用状态">
          <template #default="scope">
            <el-switch v-model="scope.row['enabled']" @change="enable(scope.row)"/>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="180">
          <template #default="scope">
            <el-button size="small" type="primary" @click="edit(scope.row)">编辑</el-button>
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
        :title="title"
        :close-on-click-modal="false"
    >
      <el-form :model="item" label-width="120px" ref="formRef" :rules="rules">
        <el-form-item label="菜单名称：" prop="name">
          <el-input v-model="item.name" autocomplete="off"/>
        </el-form-item>

        <el-form-item label="菜单图标：" prop="icon">
          <el-input v-model="item.icon" placeholder="菜单图标地址">
            <template #append>
              <el-upload
                  :auto-upload="true"
                  :show-file-list="false"
                  :http-request="uploadImg"
              >
                <el-icon class="uploader-icon">
                  <UploadFilled/>
                </el-icon>
              </el-upload>
            </template>
          </el-input>
        </el-form-item>

        <el-form-item label="菜单URL：" prop="url">
          <el-input v-model="item.url" autocomplete="off"/>
        </el-form-item>

        <el-form-item label="启用状态：" prop="enable">
          <el-switch v-model="item.enabled"/>
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
import {ElMessage} from "element-plus";
import {dateFormat, removeArrayItem} from "@/utils/libs";
import {Plus, UploadFilled} from "@element-plus/icons-vue";
import {Sortable} from "sortablejs";
import Compressor from "compressorjs";

// 变量定义
const items = ref([])
const item = ref({})
const showDialog = ref(false)
const title = ref("")
const rules = reactive({
  name: [{required: true, message: '请输入菜单名称', trigger: 'change',}],
  icon: [{required: true, message: '请上传菜单图标', trigger: 'change',}],
  url: [{required: true, message: '请输入菜单地址', trigger: 'change',}],
})
const loading = ref(true)
const formRef = ref(null)

const fetchData = () => {
  // 获取数据
  httpGet('/api/admin/menu/list').then((res) => {
    if (res.data) {
      // 初始化数据
      const arr = res.data;
      for (let i = 0; i < arr.length; i++) {
        arr[i].last_used_at = dateFormat(arr[i].last_used_at)
      }
      items.value = arr
    }
    loading.value = false
  }).catch(() => {
    ElMessage.error("获取数据失败");
  })
}

onMounted(() => {
  const drawBodyWrapper = document.querySelector('.el-table__body tbody')
  fetchData()

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

      httpPost("/api/admin/menu/sort", {ids: ids, sorts: sorts}).catch(e => {
        ElMessage.error("排序失败：" + e.message)
      })
    }
  })
})

const add = function () {
  title.value = "新增菜单"
  showDialog.value = true
  item.value = {}
}

const edit = function (row) {
  title.value = "修改菜单"
  showDialog.value = true
  item.value = row
}

const save = function () {
  formRef.value.validate((valid) => {
    if (valid) {
      showDialog.value = false
      if (!item.value.id) {
        item.value.sort_num = items.value.length + 1
      }
      httpPost('/api/admin/menu/save', item.value).then(() => {
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
  httpPost('/api/admin/menu/enable', {id: row.id, enabled: row.enabled}).then(() => {
    ElMessage.success("操作成功！")
  }).catch(e => {
    ElMessage.error("操作失败：" + e.message)
  })
}

const remove = function (row) {
  httpGet('/api/admin/menu/remove?id=' + row.id).then(() => {
    ElMessage.success("删除成功！")
    items.value = removeArrayItem(items.value, row, (v1, v2) => {
      return v1.id === v2.id
    })
  }).catch((e) => {
    ElMessage.error("删除失败：" + e.message)
  })
}

// 图片上传
const uploadImg = (file) => {
  // 压缩图片并上传
  new Compressor(file.file, {
    quality: 0.6,
    success(result) {
      const formData = new FormData();
      formData.append('file', result, result.name);
      // 执行上传操作
      httpPost('/api/admin/upload', formData).then((res) => {
        item.value.icon = res.data.url
        ElMessage.success('上传成功')
      }).catch((e) => {
        ElMessage.error('上传失败:' + e.message)
      })
    },
    error(e) {
      ElMessage.error('上传失败:' + e.message)
    },
  });
};
</script>

<style lang="stylus" scoped>
.menu {

  .opt-box {
    padding-bottom: 10px;
    display flex;
    justify-content flex-end

    .el-icon {
      margin-right: 5px;
    }
  }

  .menu-icon {
    width 36px
    height 36px
  }

  .el-select {
    width: 100%
  }

}
</style>