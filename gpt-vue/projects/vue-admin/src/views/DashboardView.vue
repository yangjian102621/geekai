<script lang="ts" setup>
import http from "@/http/config";
import { ref } from "vue";

const dataSet = {
  users: "今日新增用户",
  chats: "今日新增对话",
  tokens: "今日消耗 Tokens",
  income: "今日入账",
};

const icons = {
  users: "icon-user",
  chats: "icon-wechat",
  tokens: "icon-computer",
  income: "icon-wechatpay",
};

const data = ref<Record<string, number>>({});
const getData = () => {
  http({
    url: "api/admin/dashboard/stats",
    method: "get",
  }).then((res) => {
    data.value = res.data;
  });
};
getData();
</script>
<template>
  <div class="dashboard">
    <a-grid :cols="{ xs: 1, sm: 1, md: 2, lg: 3, xl: 4 }" :colGap="12" :rowGap="16" class="grid">
      <a-grid-item v-for="(value, key) in dataSet" :key="key">
        <div class="data-card">
          <span :class="key" class="icon"><icon-user /></span>
          <span class="count"
            ><a-statistic :extra="value" :value="data[key]" :precision="0"
          /></span>
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
    }
  }
}
</style>
