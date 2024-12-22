<template>
  <div class="container list" v-loading="loading">
    <div class="handle-box">
      <el-input v-model="query.code" placeholder="兑换码" class="handle-input mr10"></el-input>
      <el-select v-model="query.status" placeholder="核销状态" style="width: 100px" class="handle-input mr10">
        <el-option
            v-for="item in redeemStatus"
            :key="item.value"
            :label="item.label"
            :value="item.value"
        />
      </el-select>
      <el-button type="primary" :icon="Search" @click="fetchData">搜索</el-button>
      <el-button type="success" :icon="Plus" @click="add">添加兑换码</el-button>
      <el-button type="primary" @click="exportItems" :loading="exporting"><i class="iconfont icon-export mr-1"></i> 导出
      </el-button>
    </div>

    <el-row>
      <el-table :data="items" :row-key="row => row.id"
                @selection-change="handleSelectionChange" table-layout="auto">
        <el-table-column type="selection" width="38"></el-table-column>
        <el-table-column prop="name" label="名称"/>
        <el-table-column prop="code" label="兑换码">
          <template #default="scope">
            <span>{{ substr(scope.row.code, 24) }}</span>
            <el-icon class="copy-code" :data-clipboard-text="scope.row.code">
              <DocumentCopy/>
            </el-icon>
          </template>
        </el-table-column>
        <el-table-column label="兑换人">
          <template #default="scope">
            <span v-if="scope.row['username'] !== ''">{{ scope.row['username'] }}</span>
            <el-tag v-else>未兑换</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="power" label="算力"/>

        <el-table-column label="生成时间">
          <template #default="scope">
            <span>{{ dateFormat(scope.row['created_at']) }}</span>
          </template>
        </el-table-column>

        <el-table-column label="兑换时间">
          <template #default="scope">
            <span v-if="scope.row['redeemed_at'] > 0">{{ dateFormat(scope.row['redeemed_at']) }}</span>
            <el-tag v-else>未兑换</el-tag>
          </template>
        </el-table-column>

        <el-table-column prop="enabled" label="启用状态">
          <template #default="scope">
            <el-switch v-model="scope.row['enabled']" @change="set('enabled',scope.row)"
                       :disabled="scope.row['redeemed_at']>0"/>
          </template>
        </el-table-column>

        <el-table-column label="操作" width="180">
          <template #default="scope">
            <el-popconfirm title="确定要删除当前记录吗?" @confirm="remove(scope.row)">
              <template #reference>
                <el-button size="small" type="danger">删除</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>

      </el-table>
    </el-row>

    <div class="pagination">
      <el-pagination v-if="total > 0" background
                     layout="total,prev, pager, next"
                     :hide-on-single-page="true"
                     v-model:current-page="page"
                     v-model:page-size="pageSize"
                     @current-change="fetchData()"
                     :total="total"/>

    </div>

    <el-dialog
        v-model="showDialog"
        title="生成兑换码">
      <template #default>
        <el-form :model="item" label-width="120px" v-loading="dialogLoading">
          <el-form-item label="名称：" prop="name">
            <el-input v-model="item.name" autocomplete="off"/>
          </el-form-item>

          <el-form-item label="算力额度：" prop="power">
            <el-input v-model.number="item.power" autocomplete="off"/>
          </el-form-item>

          <el-form-item label="生成数量：" prop="num">
            <el-input v-model.number="item.num"/>
          </el-form-item>
        </el-form>
      </template>

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
import {onMounted, onUnmounted, ref} from "vue";
import {httpGet, httpPost, httpPostDownload} from "@/utils/http";
import {ElMessage} from "element-plus";
import {dateFormat, substr, UUID} from "@/utils/libs";
import {DocumentCopy, Plus, Search} from "@element-plus/icons-vue";
import {showMessageError} from "@/utils/dialog";
import ClipboardJS from "clipboard";

// 变量定义
const items = ref([])
const loading = ref(true)
const query = ref({code: "", status: -1})
const redeemStatus = ref([
  {value: -1, label: "全部"},
  {value: 0, label: "未核销"},
  {value: 1, label: "已核销"},
])
const showDialog = ref(false)
const dialogLoading = ref(false)
const item = ref({name: "", power: 0, num: 1})
const itemIds = ref([])
const exporting = ref(false)

const clipboard = ref(null)
onMounted(() => {
  clipboard.value = new ClipboardJS('.copy-code');
  clipboard.value.on('success', () => {
    ElMessage.success('复制成功！');
  })

  clipboard.value.on('error', () => {
    ElMessage.error('复制失败！');
  })

  fetchData()
})

onUnmounted(() => {
  clipboard.value.destroy()
})

const add = () => {
  item.value = {name: "100算力点卡", power: 100, num: 1}
  showDialog.value = true
  dialogLoading.value = false
}

const save = () => {
  if (item.value.name === "") {
    return showMessageError("请输入兑换码名称")
  }
  if (item.value.power === 0) {
    return showMessageError("请输入算力额度")
  }
  if (item.value.num <= 0) {
    return showMessageError("请输入生成数量")
  }
  dialogLoading.value = true
  httpPost('/api/admin/redeem/create', item.value).then((res) => {
    ElMessage.success(`成功生成了${res.data.counter}个兑换码`)
    showDialog.value = false
    fetchData()
  }).catch((e) => {
    ElMessage.error("生成失败：" + e.message)
  })
}

const set = (filed, row) => {
  httpPost('/api/admin/redeem/set', {id: row.id, filed: filed, value: row[filed]}).then(() => {
    ElMessage.success("操作成功！")
  }).catch(e => {
    ElMessage.error("操作失败：" + e.message)
  })
}


const page = ref(1)
const pageSize = ref(12)
const total = ref(0)
const fetchData = () => {
  query.value.page = page.value
  query.value.page_size = pageSize.value
  httpGet('/api/admin/redeem/list', query.value).then((res) => {
    if (res.data) {
      items.value = res.data.items
      total.value = res.data.total
      page.value = res.data.page
      pageSize.value = res.data.page_size
    }
    loading.value = false
  }).catch(e => {
    ElMessage.error("获取数据失败：" + e.message);
  })
}

const remove = function (row) {
  httpGet('/api/admin/redeem/remove?id=' + row.id).then(() => {
    ElMessage.success("删除成功！")
    fetchData()
  }).catch((e) => {
    ElMessage.error("删除失败：" + e.message)
  })
}

const handleSelectionChange = (items) => {
  itemIds.value = items.map(item => item.id)
}

const exportItems = () => {
  query.value.ids = itemIds.value
  exporting.value = true
  httpPostDownload("/api/admin/redeem/export", query.value).then(response => {
    const url = window.URL.createObjectURL(new Blob([response.data]));
    const link = document.createElement('a');
    link.href = url;
    link.setAttribute('download', UUID() + ".csv"); // 设置下载文件的名称
    document.body.appendChild(link);
    link.click();

    // 移除 <a> 标签
    document.body.removeChild(link);
    window.URL.revokeObjectURL(url);
    exporting.value = false
  }).catch(() => {
    exporting.value = false
    showMessageError("下载失败")
  })
}

</script>

<style lang="stylus" scoped>
.list {
  .handle-box {
    margin-bottom 20px

    .handle-input {
      max-width 150px;
      margin-right 10px;
    }
  }

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

  .copy-code {
    cursor pointer
    margin-left 6px
    position relative
    top 2px
    font-size 14px
  }

  .pagination {
    padding 20px 0
    display flex
    justify-content right
  }

}
</style>