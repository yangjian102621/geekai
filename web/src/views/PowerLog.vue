<template>
  <div class="power-log custom-scroll" v-loading="loading">
    <!-- :style="{ height: listBoxHeight + 'px' }" -->
    <div class="inner">
      <div class="list-box">
        <div class="handle-box">
          <el-input v-model="query.model" placeholder="模型" class="handle-input mr10" clearable></el-input>
          <el-date-picker
            v-model="query.date"
            type="daterange"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            format="YYYY-MM-DD"
            value-format="YYYY-MM-DD"
            style="margin: 0 10px; width: 200px"
          />
          <el-button type="primary" :icon="Search" @click="fetchData">搜索</el-button>
        </div>

        <el-row v-if="items.length > 0">
          <el-table :data="items" :row-key="(row) => row.id" table-layout="auto" border>
            <el-table-column prop="username" label="用户" width="130px" />
            <el-table-column prop="model" label="模型" width="130px" />
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
            <el-table-column prop="balance" label="余额" />
            <el-table-column label="发生时间" width="160px">
              <template #default="scope">
                <span>{{ dateFormat(scope.row["created_at"]) }}</span>
              </template>
            </el-table-column>
            <el-table-column prop="remark" label="备注" />
          </el-table>

          <div class="pagination">
            <el-pagination
              v-if="total > 0"
              background
              layout="total,prev, pager, next"
              style="--el-pagination-button-bg-color: rgba(86, 86, 95, 0.2)"
              :hide-on-single-page="true"
              v-model:current-page="page"
              v-model:page-size="pageSize"
              @current-change="fetchData()"
              :total="total"
            />
          </div>
        </el-row>
        <el-empty :image-size="100" v-else :image="nodata" description="暂无数据" />
      </div>
    </div>
  </div>
</template>

<script setup>
import nodata from "@/assets/img/no-data.png";

import { onMounted, ref } from "vue";
import { dateFormat } from "@/utils/libs";
import { Search } from "@element-plus/icons-vue";
import Clipboard from "clipboard";
import { ElMessage } from "element-plus";
import { httpPost } from "@/utils/http";
import { checkSession } from "@/store/cache";

const items = ref([]);
const total = ref(0);
const page = ref(1);
const pageSize = ref(20);
const loading = ref(false);
const listBoxHeight = window.innerHeight - 87;
const query = ref({
  model: "",
  date: [],
});
const tagColors = ref(["primary", "success", "primary", "danger", "info", "warning"]);

onMounted(() => {
  checkSession()
    .then(() => {
      fetchData();
    })
    .catch(() => {});
  const clipboard = new Clipboard(".copy-order-no");
  clipboard.on("success", () => {
    ElMessage.success("复制成功！");
  });

  clipboard.on("error", () => {
    ElMessage.error("复制失败！");
  });
});

// 获取数据
const fetchData = () => {
  loading.value = true;
  httpPost("/api/powerLog/list", {
    model: query.value.model,
    date: query.value.date,
    page: page.value,
    page_size: pageSize.value,
  })
    .then((res) => {
      if (res.data) {
        items.value = res.data.items;
        total.value = res.data.total;
        page.value = res.data.page;
        pageSize.value = res.data.page_size;
      }
      loading.value = false;
    })
    .catch((e) => {
      loading.value = false;
      ElMessage.error("获取数据失败：" + e.message);
    });
};
</script>

<style lang="stylus" scoped>
@import "@/assets/css/custom-scroll.styl"
.power-log {
  color #ffffff
  .inner {
    padding 0 20px 20px 20px
    overflow auto

    .list-box {
      overflow-x hidden
      //overflow-y auto
      background: var(--chat-bg);
      padding: 20px;
      margin-top: 20px;
      border-radius: 10px;
      .handle-box {
        padding 0 20px 20px 0

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
        margin-top 20px
      }
    }
  }
}
</style>
