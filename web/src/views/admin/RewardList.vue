<template>
  <div class="container list" v-loading="loading">

    <el-row>
      <el-table :data="items" :row-key="row => row.id">
        <el-table-column prop="username" label="用户"/>
        <el-table-column prop="tx_id" label="转账单号"/>
        <el-table-column prop="amount" label="转账金额"/>
        <el-table-column prop="remark" label="备注"/>

        <el-table-column label="转账时间">
          <template #default="scope">
            <span>{{ dateFormat(scope.row['created_at']) }}</span>
          </template>
        </el-table-column>

        <el-table-column label="核销时间">
          <template #default="scope">
            <span v-if="scope.row['status']">{{ dateFormat(scope.row['updated_at']) }}</span>
            <el-tag v-else>未核销</el-tag>
          </template>
        </el-table-column>

      </el-table>
    </el-row>

  </div>
</template>

<script setup>
import {ref} from "vue";
import {httpGet} from "@/utils/http";
import {ElMessage} from "element-plus";
import {dateFormat} from "@/utils/libs";

// 变量定义
const items = ref([])
const loading = ref(true)

// 获取数据
httpGet('/api/admin/reward/list').then((res) => {
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