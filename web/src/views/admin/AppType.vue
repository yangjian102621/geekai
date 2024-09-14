<template>
  <div class="container app-type" v-loading="loading">
    <div class="handle-box">
      <el-button type="primary" :icon="Plus" @click="add">新增</el-button>
    </div>

    <el-row>
      <el-table :data="items" :row-key="row => row.id" table-layout="auto">
        <el-table-column type="selection" width="38"></el-table-column>
        <el-table-column prop="name" label="分类名称">
          <template #default="scope">
            <span class="sort" :data-id="scope.row.id">
              <i class="iconfont icon-drag"></i>
              {{ scope.row.name }}
            </span>
          </template>
        </el-table-column>
        <el-table-column label="图标" prop="icon">
          <template #default="scope">
            <el-image v-if="scope.row.icon" :src="scope.row.icon" style="width: 45px; height: 45px; border-radius: 50%"/>
            <el-tag type="info" v-else>无图标</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="enabled" label="启用状态">
          <template #default="scope">
            <el-switch v-model="scope.row['enabled']" @change="enableSet(scope.row)"/>
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
        style="width: 90%; max-width: 600px;"
    >
      <el-form :model="item" label-width="120px" ref="formRef" :rules="rules">
        <el-form-item label="分类名称：" prop="name">
          <el-input v-model="item.name" autocomplete="off"/>
        </el-form-item>
        <el-form-item label="应用图标：" prop="icon">
          <el-input v-model="item.icon">
            <template #append>
              <el-upload
                  :auto-upload="true"
                  :show-file-list="false"
                  :http-request="uploadImg"
              >
                上传
              </el-upload>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item label="启用状态：" prop="enable">
          <el-switch v-model="item.enabled" />
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
import {onMounted, onUnmounted, reactive, ref} from "vue";
import {httpGet, httpPost} from "@/utils/http";
import {ElMessage} from "element-plus";
import {removeArrayItem} from "@/utils/libs";
import {Sortable} from "sortablejs";
import Compressor from "compressorjs";

// 变量定义
const items = ref([])
const item = ref({})
const showDialog = ref(false)
const title = ref("")
const rules = reactive({
  name: [{required: true, message: '请输入分类名称', trigger: 'change',}],
})
const loading = ref(true)
const formRef = ref(null)

// 获取数据
const fetchData = () => {
  httpGet('/api/admin/app/type/list').then((res) => {
    if (res.data) {
      items.value = res.data
    }
    loading.value = false
  }).catch(() => {
    ElMessage.error("获取数据失败");
  })
}

onMounted(() => {
  fetchData()
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
        sorts.push(index + 1)
        items.value[index].sort_num = index + 1
      })

      httpPost("/api/admin/app/type/sort", {ids: ids, sorts: sorts}).then(() => {
      }).catch(e => {
        ElMessage.error("排序失败：" + e.message)
      })
    }
  })
})

const add = function () {
  title.value = "新增分类"
  showDialog.value = true
  item.value = { enabled: true, }
}

const edit = function (row) {
  title.value = "修改分类"
  showDialog.value = true
  item.value = row
}

const save = function () {
  formRef.value.validate((valid) => {
    if (!item.value.sort_num) {
      item.value.sort_num = items.value.length
    }
    if (valid) {
      showDialog.value = false
      httpPost('/api/admin/app/type/save', item.value).then(() => {
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

// 设置启用状态
const enableSet = (row) => {
  httpPost('/api/admin/app/type/enable', {id: row.id, enabled: row.enabled}).then(() => {
    ElMessage.success("操作成功！")
  }).catch(e => {
    ElMessage.error("操作失败：" + e.message)
  })
}

// 删除数据
const remove = function (row) {
  httpGet('/api/admin/app/type/remove?id=' + row.id).then(() => {
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