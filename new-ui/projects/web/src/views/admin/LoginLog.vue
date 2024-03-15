<template>
  <div class="container list" v-loading="loading">
    <el-row>
      <el-table :data="items" border :row-key="row => row.id">
        <el-table-column label="用户名" prop="username"/>
        <el-table-column label="登录IP" prop="login_ip"/>
        <el-table-column label="登录地址" prop="login_address"/>
        <el-table-column label="登录时间">
          <template #default="scope">
            <span>{{ dateFormat(scope.row['created_at']) }}</span>
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
                     @current-change="fetchList(page, pageSize)"
                     :total="total"/>

    </div>

  </div>
</template>

<script setup>
import {onMounted, ref} from "vue";
import {httpGet} from "@/utils/http";
import {ElMessage} from "element-plus";
import {dateFormat} from "@/utils/libs";

// 用户登录日志
const items = ref([])
const loading = ref(true)
const total = ref(0)
const page = ref(0)
const pageSize = ref(0)

onMounted(() => {
  fetchList(1, 15)
})

// 获取数据
const fetchList = function (_page, _pageSize) {
  httpGet(`/api/admin/user/loginLog?page=${_page}&page_size=${_pageSize}`).then((res) => {
    if (res.data) {
      items.value = res.data.items
      total.value = res.data.total
      page.value = res.data.page
      pageSize.value = res.data.page_size
    }
    loading.value = false
  }).catch(() => {
    ElMessage.error("获取数据失败");
  })
}
</script>

<style lang="stylus" scoped>
.list {

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