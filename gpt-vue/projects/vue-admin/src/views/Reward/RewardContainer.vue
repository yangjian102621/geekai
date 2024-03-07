<script lang="ts" setup>
import type { TableColumnData } from "@arco-design/web-vue";
import { dateFormat } from "@gpt-vue/packages/utils";
import SimpleTable from "@/components/SimpleTable/SimpleTable.vue";
import { getList } from "./api";

const columns: TableColumnData[] = [
  {
    dataIndex: "username",
    title: "用户",
  },
  {
    dataIndex: "tx_id",
    title: "转账单号",
  },
  {
    dataIndex: "amount",
    title: "转账金额",
  },
  {
    dataIndex: "remark",
    title: "备注",
  },
  {
    dataIndex: "created_at",
    title: "转账时间",
    render: ({ record }) => dateFormat(record.created_at),
  },
  {
    title: "核销时间",
    slotName: "updated_at",
  },
  {
    title: "兑换详情",
    slotName: "exchange",
  },
  {
    title: "操作",
    slotName: "actions",
    fixed: "right",
  },
];
</script>
<template>
  <SimpleTable :request="getList" :columns="columns">
    <template #updated_at="{ record }">
      <span v-if="record.status">{{ dateFormat(record.updated_at) }}</span>
      <a-tag v-else>未核销</a-tag>
    </template>
    <template #exchange="{ record }">
      <a-tag v-if="record.exchange.calls > 0"
        >聊天{{ record.exchange.calls }}次</a-tag
      >
      <a-tag v-else-if="record.exchange.img_calls > 0"
        >绘图{{ record.exchange.img_calls }}次</a-tag
      >
    </template>
  </SimpleTable>
</template>
