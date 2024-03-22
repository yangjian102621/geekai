<template>
  <div class="dashboard" v-loading="loading">
    <el-row class="mgb20" :gutter="20">
      <el-col :span="6">
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
      <el-col :span="6">
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
      <el-col :span="6">
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
      <el-col :span="6">
        <el-card shadow="hover" :body-style="{ padding: '0px' }">
          <div class="grid-content grid-con-3">
            <el-icon class="grid-con-icon">
              <i class="iconfont icon-reward"></i>
            </el-icon>
            <div class="grid-cont-right">
              <div class="grid-num">￥{{ stats.income }}</div>
              <div>今日入账</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row class="mgb20" :gutter="20">
      <el-col :span="8">
        <div class="e-chart">
          <div id="chart-users" style="height: 400px"></div>
          <div class="title">最近7日用户注册</div>
        </div>
      </el-col>

      <el-col :span="8">
        <div class="e-chart">
          <div id="chart-tokens" style="height: 400px"></div>
          <div class="title">最近7日Token消耗</div>
        </div>
      </el-col>

      <el-col :span="8">
        <div class="e-chart">
          <div id="chart-income" style="height: 400px"></div>
          <div class="title">最近7日收入</div>
        </div>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import {onMounted, ref} from 'vue';
import {ChatDotRound, TrendCharts, User} from "@element-plus/icons-vue";
import {httpGet} from "@/utils/http";
import {ElMessage} from "element-plus";
import * as echarts from 'echarts';

const stats = ref({users: 0, chats: 0, tokens: 0, rewards: 0})
const loading = ref(true)

onMounted(() => {
  const chartUsers = echarts.init(document.getElementById("chart-users"))
  const chartTokens = echarts.init(document.getElementById("chart-tokens"))
  const chartIncome = echarts.init(document.getElementById("chart-income"))
  httpGet('/api/admin/dashboard/stats').then((res) => {
    stats.value.users = res.data.users
    stats.value.chats = res.data.chats
    stats.value.tokens = res.data.tokens
    stats.value.income = res.data.income
    const chartData = res.data.chart
    loading.value = false

    const x = []
    const dataUsers = []
    for (let k in chartData.users) {
      x.push(k)
      dataUsers.push(chartData.users[k])
    }
    chartUsers.setOption({
      xAxis: {
        data: x
      },
      yAxis: {},
      series: [
        {
          data: dataUsers,
          type: 'line',
          label: {
            show: true,
            position: 'bottom',
            textStyle: {
              fontSize: 18
            }
          }
        }
      ]
    })
    const dataTokens = []
    for (let k in chartData.historyMessage) {
      dataTokens.push(chartData.historyMessage[k])
    }
    chartTokens.setOption({
      xAxis: {
        data: x
      },
      yAxis: {},
      series: [
        {
          data: dataTokens,
          type: 'line',
          label: {
            show: true,
            position: 'bottom',
            textStyle: {
              fontSize: 18
            }
          }
        }
      ]
    })

    const dataIncome = []
    for (let k in chartData.orders) {
      dataIncome.push(chartData.orders[k])
    }
    chartIncome.setOption({
      xAxis: {
        data: x
      },
      yAxis: {},
      series: [
        {
          data: dataIncome,
          type: 'line',
          label: {
            show: true,
            position: 'bottom',
            textStyle: {
              fontSize: 18
            }
          }
        }
      ]
    })

  }).catch((e) => {
    ElMessage.error("获取统计数据失败：" + e.message)
  })

  window.onresize = function () { // 自适应大小
    chartUsers.resize()
    chartTokens.resize()
    chartIncome.resize()
  };
})

</script>

<style scoped lang="stylus">
.dashboard {
  padding 20px

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

    .iconfont {
      font-size: 50px;
    }
  }

  .e-chart {
    .title {
      text-align center
      font-size 16px
      color #444444
    }
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
