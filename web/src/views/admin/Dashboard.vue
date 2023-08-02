<template>
  <div class="dashboard">
    <el-row class="mgb20" :gutter="20">
      <el-col :span="8">
        <el-card shadow="hover" :body-style="{ padding: '0px' }">
          <div class="grid-content grid-con-1">
            <el-icon class="grid-con-icon">
              <User/>
            </el-icon>
            <div class="grid-cont-right">
              <div class="grid-num">{{ stats.users }}</div>
              <div>今日新增用户</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card shadow="hover" :body-style="{ padding: '0px' }">
          <div class="grid-content grid-con-2">
            <el-icon class="grid-con-icon">
              <ChatDotRound/>
            </el-icon>
            <div class="grid-cont-right">
              <div class="grid-num">{{ stats.chats }}</div>
              <div>今日新增对话</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card shadow="hover" :body-style="{ padding: '0px' }">
          <div class="grid-content grid-con-3">
            <el-icon class="grid-con-icon">
              <TrendCharts/>
            </el-icon>
            <div class="grid-cont-right">
              <div class="grid-num">{{ stats.tokens }}</div>
              <div>今日消耗 Tokens</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import {ref} from 'vue';
import {ChatDotRound, TrendCharts, User} from "@element-plus/icons-vue";
import {httpGet} from "@/utils/http";
import {ElMessage} from "element-plus";

const stats = ref({users: 0, chats: 0, tokens: 0})

httpGet('/api/admin/dashboard/stats').then((res) => {
  stats.value.users = res.data.users
  stats.value.chats = res.data.chats
  stats.value.tokens = res.data.tokens
}).catch((e) => {
  ElMessage.error("获取统计数据失败：" + e.message)
})

</script>

<style scoped lang="stylus">
.dashboard {
  .grid-content {
    display: flex;
    align-items: center;
    height: 100px;
  }

  .grid-cont-right {
    flex: 1;
    text-align: center;
    font-size: 14px;
    color: #999;
  }

  .grid-num {
    font-size: 30px;
    font-weight: bold;
  }

  .grid-con-icon {
    font-size: 50px;
    width: 100px;
    height: 100px;
    text-align: center;
    line-height: 100px;
    color: #fff;
  }

  .grid-con-1 .grid-con-icon {
    background: rgb(45, 140, 240);
  }

  .grid-con-1 .grid-num {
    color: rgb(45, 140, 240);
  }

  .grid-con-2 .grid-con-icon {
    background: rgb(100, 213, 114);
  }

  .grid-con-2 .grid-num {
    color: rgb(100, 213, 114);
  }

  .grid-con-3 .grid-con-icon {
    background: rgb(242, 94, 67);
  }

  .grid-con-3 .grid-num {
    color: rgb(242, 94, 67);
  }
}


</style>
