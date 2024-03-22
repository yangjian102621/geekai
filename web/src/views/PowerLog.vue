<template>
  <div class="power-log" v-loading="loading">
    <div class="inner">
      <h2>消费日志</h2>

      <div class="list-box" :style="{height: listBoxHeight + 'px'}">
        <div class="handle-box">
          <el-input v-model="query.model" placeholder="模型" class="handle-input mr10" clearable></el-input>
          <el-date-picker
              v-model="query.date"
              type="daterange"
              start-placeholder="开始日期"
              end-placeholder="结束日期"
              format="YYYY-MM-DD"
              value-format="YYYY-MM-DD"
              style="margin: 0 10px;width: 200px; position: relative;top:3px;"
          />
          <el-button type="primary" :icon="Search" @click="fetchData">搜索</el-button>
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
    </div>
  </div>
</template>

<script setup>
import {onMounted, ref} from "vue"
import {dateFormat} from "@/utils/libs";
import {Back, DocumentCopy, Search} from "@element-plus/icons-vue";
import Clipboard from "clipboard";
import {ElMessage} from "element-plus";
import {httpPost} from "@/utils/http";

const items = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(20)
const loading = ref(false)
const listBoxHeight = window.innerHeight - 117
const query = ref({
  model: "",
  date: []
})
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

// 获取数据
const fetchData = () => {
  loading.value = true
  httpPost('/api/powerLog/list', {
    model: query.value.model,
    date: query.value.date,
    page: page.value,
    page_size: pageSize.value
  }).then((res) => {
    if (res.data) {
      items.value = res.data.items
      total.value = res.data.total
      page.value = res.data.page
      pageSize.value = res.data.page_size
    }
    loading.value = false
  }).catch(e => {
    loading.value = false
    ElMessage.error("获取数据失败：" + e.message);
  })
}

</script>

<style lang="stylus" scoped>
.power-log {
  background-color #ffffff

  .inner {
    padding 0 20px 20px 20px

    ::-webkit-scrollbar {
      width: 8px; /* 滚动条宽度 */
    }

    /* 修改滚动条轨道的背景颜色 */

    ::-webkit-scrollbar-track {
      background-color: #ffffff;
    }

    /* 修改滚动条的滑块颜色 */

    ::-webkit-scrollbar-thumb {
      background-color: #cccccc;
      border-radius 8px
    }

    /* 修改滚动条的滑块的悬停颜色 */

    ::-webkit-scrollbar-thumb:hover {
      background-color: #999999;
    }

    .list-box {
      overflow-x hidden
      //overflow-y auto

      .handle-box {
        padding 20px 0

        .el-input {
          max-width 150px
        }

        .el-date-editor {
          max-width 200px;
        }
      }

      .pagination {
        display flex
        justify-content center
        width 100%
        padding 20px
      }
    }
  }
}
</style>
