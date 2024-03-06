<script lang="ts" setup>
import SearchTable from "@/components/SearchTable/SearchTable.vue";
import type { SearchTableColumns } from "@/components/SearchTable/type";
import { getList } from "./api";

const columns: SearchTableColumns[] = [
  {
    dataIndex: "order_no",
    title: "订单号",
    search: {
      valueType: "input",
    },
  },
  {
    dataIndex: "username",
    title: "下单用户",
  },
  {
    dataIndex: "subject",
    title: "产品名称",
  },
  {
    dataIndex: "amount",
    title: "订单金额",
  },
  {
    dataIndex: "remark.calls",
    title: "调用次数",
  },
  {
    dataIndex: "created_at",
    title: "下单时间",
  },
  {
    dataIndex: "status",
    title: "订单状态",
    hideInTable: true,
    search: {
      valueType: "select",
      fieldProps: {
        options: [
          { label: "全部", value: -1 },
          { label: "未支付", value: 0 },
          { label: "已支付", value: 2 },
        ],
      },
    },
  },
  {
    dataIndex: "pay_time",
    title: "支付时间",
    search: {
      valueType: "range",
    },
  },
  {
    dataIndex: "pay_way",
    title: "支付方式",
  },
  {
    title: "操作",
    slotName: "actions",
  },
];
</script>
<template>
  <SearchTable :request="getList" :columns="columns">
    <template #actions="{ record }">
      <a-link :key="record.id" status="danger">删除</a-link>
    </template>
  </SearchTable>
</template>
