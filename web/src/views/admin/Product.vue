<template>
  <div class="container list" v-loading="loading">

    <div class="handle-box">
      <el-button type="primary" :icon="Plus" @click="add">新增</el-button>
    </div>

    <el-row>
      <el-table :data="items" :row-key="row => row.id" table-layout="auto">
        <el-table-column prop="name" label="产品名称">
          <template #default="scope">
            <span class="sort" :data-id="scope.row.id">{{ scope.row.name }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="price" label="产品价格"/>
        <el-table-column prop="discount" label="优惠金额"/>
        <el-table-column prop="days" label="有效期(天)">
          <template #default="scope">
            <el-tag v-if="scope.row.days === 0">长期有效</el-tag>
            <span v-else>{{ scope.row.days }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="calls" label="对话次数"/>
        <el-table-column prop="img_calls" label="绘图次数"/>
        <el-table-column prop="sales" label="销量"/>
        <el-table-column prop="enabled" label="启用状态">
          <template #default="scope">
            <el-switch v-model="scope.row['enabled']" @change="enable(scope.row)"/>
          </template>
        </el-table-column>

        <el-table-column label="更新时间">
          <template #default="scope">
            <span>{{ dateFormat(scope.row['updated_at']) }}</span>
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
        style="width: 90%; max-width: 600px;"
    >
      <el-form :model="item" label-width="120px" ref="formRef" :rules="rules">
        <el-form-item label="产品名称：" prop="name">
          <el-input v-model="item.name" autocomplete="off"/>
        </el-form-item>

        <el-form-item label="产品价格：" prop="price">
          <el-input v-model="item.price" autocomplete="off"/>
        </el-form-item>

        <el-form-item label="优惠金额：" prop="discount">
          <el-input v-model="item.discount" autocomplete="off"/>
        </el-form-item>

        <el-form-item label="有效期：" prop="days">
          <el-input v-model.number="item.days" autocomplete="off" placeholder="会员有效期(天)"/>
        </el-form-item>

        <el-form-item label="对话次数：" prop="calls">
          <el-input v-model.number="item.calls" autocomplete="off" placeholder="增加对话次数"/>
        </el-form-item>

        <el-form-item label="绘图次数：" prop="img_calls">
          <el-input v-model.number="item.img_calls" autocomplete="off" placeholder="增加绘图次数"/>
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
import {Plus} from "@element-plus/icons-vue";
import {Sortable} from "sortablejs";

// 变量定义
const items = ref([])
const item = ref({})
const showDialog = ref(false)
const title = ref("")
const rules = reactive({
  name: [{required: true, message: '请输入产品名称', trigger: 'change',}],
  price: [{required: true, message: '请输产品价格', trigger: 'change',}],
  discount: [{required: true, message: '请输优惠金额', trigger: 'change',}],
  days: [{required: true, message: '请输入有效期', trigger: 'change',}],
})
const loading = ref(true)
const formRef = ref(null)

// 获取数据
httpGet('/api/admin/product/list').then((res) => {
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

      httpPost("/api/admin/product/sort", {ids: ids, sorts: sorts}).catch(e => {
        ElMessage.error("排序失败：" + e.message)
      })
    }
  })
})

const add = function () {
  title.value = "新增模型"
  showDialog.value = true
  item.value = {}
}

const edit = function (row) {
  title.value = "修改模型"
  showDialog.value = true
  item.value = row
}

const save = function () {
  formRef.value.validate((valid) => {
    if (valid) {
      showDialog.value = false
      item.value['price'] = parseFloat(item.value['price'])
      item.value['discount'] = parseFloat(item.value['discount'])
      httpPost('/api/admin/product/save', item.value).then((res) => {
        ElMessage.success('操作成功！')
        if (!item.value['id']) {
          const newItem = res.data
          items.value.push(newItem)
        }
      }).catch((e) => {
        ElMessage.error('操作失败，' + e.message)
      })
    } else {
      return false
    }
  })
}

const enable = (row) => {
  httpPost('/api/admin/product/enable', {id: row.id, enabled: row.enabled}).then(() => {
    ElMessage.success("操作成功！")
  }).catch(e => {
    ElMessage.error("操作失败：" + e.message)
  })
}

const remove = function (row) {
  httpGet('/api/admin/product/remove?id=' + row.id).then(() => {
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