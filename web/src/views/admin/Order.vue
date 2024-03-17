<template>
  <div class="container order" v-loading="loading">
    <div class="handle-box">
      <el-input v-model="query.order_no" placeholder="订单号" class="handle-input mr10"></el-input>
      <el-select v-model="query.status" placeholder="订单状态" style="width: 100px">
        <el-option
            v-for="item in orderStatus"
            :key="item.value"
            :label="item.label"
            :value="item.value"
        />
      </el-select>
      <el-date-picker
          v-model="query.pay_time"
          type="daterange"
          start-placeholder="开始日期"
          end-placeholder="结束日期"
          format="YYYY-MM-DD"
          value-format="YYYY-MM-DD"
          style="margin: 0 10px;width: 200px; position: relative;top:3px;"
      />
      <el-button type="primary" :icon="Search" @click="fetchData">搜索</el-button>
    </div>

    <el-row>
      <el-table :data="items" :row-key="row => row.id" table-layout="auto">
        <el-table-column prop="order_no" label="订单号"/>
        <el-table-column prop="username" label="下单用户"/>
        <el-table-column prop="subject" label="产品名称"/>
        <el-table-column prop="amount" label="订单金额"/>
        <el-table-column label="充值算力">
          <template #default="scope">
            <span>{{ scope.row.remark?.power }}</span>
          </template>
        </el-table-column>

        <el-table-column label="下单时间">
          <template #default="scope">
            <span>{{ dateFormat(scope.row['created_at']) }}</span>
          </template>
        </el-table-column>

        <el-table-column label="支付时间">
          <template #default="scope">
            <span v-if="scope.row['pay_time']">{{ dateFormat(scope.row['pay_time']) }}</span>
            <el-tag v-else>未支付</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="pay_way" label="支付方式"/>

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
  </div>
</template>

<script setup>
import {onMounted, ref} from "vue";
import {httpGet, httpPost} from "@/utils/http";
import {ElMessage} from "element-plus";
import {dateFormat, removeArrayItem} from "@/utils/libs";
import {Search} from "@element-plus/icons-vue";

// 变量定义
const items = ref([])
const query = ref({order_no: "", pay_time: [], status: -1})
const total = ref(0)
const page = ref(1)
const pageSize = ref(15)
const loading = ref(true)
const orderStatus = ref([
  {value: -1, label: "全部"},
  {value: 0, label: "未支付"},
  {value: 2, label: "已支付"},
])

onMounted(() => {
  fetchData()
})
// 获取数据
const fetchData = () => {
  query.value.page = page.value
  query.value.page_size = pageSize.value
  httpPost('/api/admin/order/list', query.value).then((res) => {
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
  httpGet('/api/admin/order/remove?id=' + row.id).then(() => {
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
.order {

  .handle-box {
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

}
</style>