<script lang="ts" setup>
import http from "@/http/config";
import { ref, nextTick } from "vue";
import * as echarts from "echarts/core";
import { GridComponent, TitleComponent } from "echarts/components";
import { LineChart } from "echarts/charts";
import { UniversalTransition } from "echarts/features";
import { CanvasRenderer } from "echarts/renderers";

const icons = {
  users: "icon-user",
  chats: "icon-wechat",
  tokens: "icon-computer",
  income: "icon-wechatpay",
};
const dataSet = {
  users: "今日新增用户",
  chats: "今日新增对话",
  tokens: "今日消耗 Tokens",
  income: "今日入账",
};

const getData = () => {
  http({
    url: "api/admin/dashboard/stats",
    method: "get",
  }).then((res) => {
    data.value = res.data;
    handeChartData(res.data.chart);
  });
};
getData();
// 图表
const chartTitle = {
  historyMessage: "对话",
  orders: "订单",
  users: "用户数",
};
echarts.use([GridComponent, LineChart, CanvasRenderer, UniversalTransition, TitleComponent]);
const chartDomRefs = [];
const chartData = ref({});
const data = ref<Record<string, number>>({});
const handeChartData = (data) => {
  const _chartData = {};
  for (let key in data) {
    const type = data[key];
    _chartData[key] = {
      series: [],
      xAxis: [],
    };
    for (let date in type) {
      _chartData[key].series.push(type[date]);
      _chartData[key].xAxis.push(date);
    }
    nextTick(() => {
      const myChart = echarts.init(chartDomRefs.pop());
      myChart.setOption(createOption(_chartData[key], key));
    });
  }
  chartData.value = _chartData;
};
const createOption = (data, key) => {
  const { xAxis, series } = data;
  return {
    title: {
      left: "center",
      text: chartTitle[key],
    },
    xAxis: {
      type: "category",
      data: xAxis,
    },
    yAxis: {
      type: "value",
    },
    series: [
      {
        data: series,
        type: "line",
      },
    ],
  };
};
</script>
<template>
  <div class="dashboard">
    <a-grid :cols="{ xs: 1, sm: 1, md: 2, lg: 3, xl: 4 }" :colGap="12" :rowGap="16" class="grid">
      <a-grid-item v-for="(value, key) in dataSet" :key="key">
        <div class="data-card">
          <span :class="key" class="icon"><component :is="icons[key]" /> </span>
          <span class="count"
            ><a-statistic :extra="value" :value="data[key]" :precision="0"
          /></span>
        </div>
      </a-grid-item>
    </a-grid>
  </div>
  <div class="chart">
    <a-grid
      :cols="{ xs: 1, sm: 1, md: 1, lg: 2, xl: 2, xxl: 3 }"
      :colGap="12"
      :rowGap="16"
      class="grid"
    >
      <a-grid-item v-for="(value, key, index) in chartData" :key="key">
        <div
          :ref="
            (el) => {
              chartDomRefs[index] = el;
            }
          "
          class="chartDom"
        >
          {{ key }}
        </div>
      </a-grid-item>
    </a-grid>
  </div>
</template>
<style lang="less" scoped>
.dashboard {
  display: flex;
  text-align: center;
  .grid {
    width: 100%;
  }
  .data-card {
    width: 100%;
    display: flex;
    flex: 0 0 25%;
    padding: 0 10px;
    box-sizing: border-box;
    .icon {
      display: inline-block;
      font-size: 50px;
      width: 100px;
      height: 100px;
      text-align: center;
      line-height: 100px;
      color: #fff;
    }
    .users {
      background: #2d8cf0;
    }
    .chats {
      background: #64d572;
    }
    .tokens {
      background: #f25e43;
    }
    .income {
      background: #f25e43;
    }
    .count {
      flex: 1;
      display: flex;
      align-items: center;
      justify-content: center;
      background: #f3f3f3;
    }
  }
}
.chart {
  margin-top: 15px;
  .chartDom {
    width: 450px;
    height: 500px;
  }
}
</style>
