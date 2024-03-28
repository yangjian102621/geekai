<template>
  <div class="container power-log" v-loading="loading">
    <div class="handle-box">
      <el-input v-model="query.model" placeholder="模型" class="handle-input mr10" clearable></el-input>
      <el-select v-model="query.type" placeholder="类别" style="width: 100px">
        <el-option label="全部" :value="0"/>
        <el-option label="充值" :value="1"/>
        <el-option label="消费" :value="2"/>
        <el-option label="退款" :value="3"/>
        <el-option label="奖励" :value="4"/>
        <el-option label="众筹" :value="5"/>
      </el-select>
      <el-date-picker
          v-model="query.date"
          type="daterange"
          start-placeholder="开始日期"
          end-placeholder="结束日期"
          format="YYYY-MM-DD"
          value-format="YYYY-MM-DD"
          style="margin: 0 10px;width: 200px; position: relative;top:3px;"
      />
      <el-button type="primary" :icon="Search" @click="search">搜索</el-button>

      <el-button v-if="totalPower > 0">算力总额：{{ totalPower }}</el-button>
    </div>

    <el-row v-if="items.length > 0">
      <el-table :data="items" :row-key="row => row.id" table-layout="auto" border>
        <el-table-column prop="username" label="用户"/>
        <el-table-column prop="model" label="模型"/>
        <el-table-column prop="type" label="类型">
          <template #default="scope">
            <el-tag size="small" :type="tagColors[scope.row.type]">{{ scope.row.type_str }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="数额">
          <template #default="scope">
            <div>
              <el-text type="success" v-if="scope.row.mark === 1">+{{ scope.row.amount }}</el-text>
              <el-text type="danger" v-if="scope.row.mark === 0">-{{ scope.row.amount }}</el-text>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="balance" label="余额"/>
        <el-table-column label="发生时间">
          <template #default="scope">
            <span>{{ dateFormat(scope.row['created_at']) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="remark" label="备注"/>
      </el-table>

      <div class="pagination">
        <el-pagination v-if="total > 0" background
                       layout="total,prev, pager, next"
                       :hide-on-single-page="true"
                       v-model:current-page="page"
                       v-model:page-size="pageSize"
                       @current-change="fetchData()"
                       :total="total"/>

      </div>
    </el-row>
    <el-empty :image-size="100" v-else/>

  </div>
</template>

<script setup>
import {onMounted, ref} from "vue";
import {httpGet, httpPost} from "@/utils/http";
import {ElMessage} from "element-plus";
import {dateFormat} from "@/utils/libs";
import {Search} from "@element-plus/icons-vue";
import Clipboard from "clipboard";

const items = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(20)
const loading = ref(false)
const query = ref({
  model: "",
  date: [],
  type: 0
})
const totalPower = ref(0)

const tagColors = ref(["", "success", "primary", "danger", "info", "warning"])

onMounted(() => {
  fetchData()
  const clipboard = new Clipboard('.copy-order-no');
  clipboard.on('success', () => {
    ElMessage.success("复制成功！");
  })

  clipboard.on('error', () => {
    ElMessage.error('复制失败！');
  })
})

// 搜索
const search = () => {
  page.value = 1
  fetchData()
}

// 获取数据
const fetchData = () => {
  loading.value = true
  httpPost('/api/admin/powerLog/list', {
    model: query.value.model,
    date: query.value.date,
    type: query.value.type,
    page: page.value,
    page_size: pageSize.value
  }).then((res) => {
    const data = res.data.data
    if (data) {
      items.value = data.items
      total.value = data.total
      page.value = data.page
      pageSize.value = data.page_size
    }
    totalPower.value = res.data.stat
    loading.value = false
  }).catch(e => {
    loading.value = false
    ElMessage.error("获取数据失败：" + e.message);
  })
}
</script>

<style lang="stylus" scoped>
.power-log {
  .handle-box {
    .handle-input {
      max-width 150px;
      margin-right 10px;
    }
  }

  .opt-box {
    padding-bottom: 10px;
    display flex;
    justify-content flex-start

    .el-icon {
      margin-right: 5px;
    }
  }

  .el-select {
    width: 100%
  }

}
</style>