<template>
  <div class="user-bill" v-loading="loading">
    <el-row v-if="items.length > 0">
      <el-table :data="items" :row-key="row => row.id" table-layout="auto" border
                style="--el-table-border-color:#373C47;
                --el-table-tr-bg-color:#2D323B;
                --el-table-row-hover-bg-color:#373C47;
                --el-table-header-bg-color:#474E5C;
                --el-table-text-color:#d1d1d1">
        <el-table-column prop="order_no" label="订单号">
          <template #default="scope">
            <span>{{ scope.row.order_no }}</span>
            <el-icon class="copy-order-no" :data-clipboard-text="scope.row.order_no">
              <DocumentCopy/>
            </el-icon>
          </template>
        </el-table-column>
        <el-table-column prop="subject" label="产品名称"/>
        <el-table-column prop="amount" label="订单金额"/>
        <el-table-column label="订单算力">
          <template #default="scope">
            <span>{{ scope.row.remark?.power }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="pay_way" label="支付方式"/>
        <el-table-column label="支付时间">
          <template #default="scope">
            <span v-if="scope.row['pay_time']">{{ dateFormat(scope.row['pay_time']) }}</span>
            <el-tag v-else>未支付</el-tag>
          </template>
        </el-table-column>
      </el-table>
    </el-row>
    <el-empty :image-size="100" v-else/>
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
import {onMounted, ref, watch} from "vue";
import {httpPost} from "@/utils/http";
import {ElMessage} from "element-plus";
import {dateFormat} from "@/utils/libs";
import {DocumentCopy} from "@element-plus/icons-vue";
import Clipboard from "clipboard";

const items = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(12)
const loading = ref(true)

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

// 获取数据
const fetchData = () => {
  httpPost('/api/order/list', {page: page.value, page_size: pageSize.value}).then((res) => {
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
</script>

<style scoped lang="stylus">
.user-bill {
  .pagination {
    margin: 20px 0 0 0;
    display: flex;
    justify-content: center;
    width: 100%;
  }

  .copy-order-no {
    cursor pointer
    position relative
    left 6px
    top 2px
    color #20a0ff
  }
}
</style>