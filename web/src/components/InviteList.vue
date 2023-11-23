<template>
  <div class="invite-list" v-loading="loading">
    <el-row v-if="items.length > 0">
      <el-table :data="items" :row-key="row => row.id" table-layout="auto" border
                style="--el-table-border-color:#373C47;
                --el-table-tr-bg-color:#2D323B;
                --el-table-row-hover-bg-color:#373C47;
                --el-table-header-bg-color:#474E5C;
                --el-table-text-color:#d1d1d1">
        <el-table-column prop="username" label="用户"/>
        <el-table-column prop="invite_code" label="邀请码"/>
        <el-table-column label="邀请奖励">
          <template #default="scope">
            <span>对话：{{ scope.row['reward']['chat_calls'] }}次</span>，
            <span>绘图：{{ scope.row['reward']['chat_calls'] }}次</span>
          </template>
        </el-table-column>

        <el-table-column label="注册时间">
          <template #default="scope">
            <span>{{ dateFormat(scope.row['created_at']) }}</span>
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
import {httpGet, httpPost} from "@/utils/http";
import {ElMessage} from "element-plus";
import {dateFormat} from "@/utils/libs";
import {DocumentCopy} from "@element-plus/icons-vue";
import Clipboard from "clipboard";

const items = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)
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
  httpPost('/api/invite/list', {page: page.value, page_size: pageSize.value}).then((res) => {
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
.invite-list {
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