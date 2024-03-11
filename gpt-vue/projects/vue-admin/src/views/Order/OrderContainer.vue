<script lang="ts" setup>
import SearchTable from "@/components/SearchTable/SearchTable.vue";
import type { SearchTableColumns } from "@/components/SearchTable/type";
import { dateFormat } from "@gpt-vue/packages/utils";
import { getList } from "./api";

const columns: SearchTableColumns[] = [
  {
    dataIndex: "order_no",
    title: "订单号",
    search: {
      valueType: "input",
    },
    width: 280,
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
    render: ({ record }) => dateFormat(record.created_at),
    width: 200,
  },
  {
    dataIndex: "status",
    title: "订单状态",
    hideInTable: true,
    search: {
      valueType: "select",
      defaultValue: -1,
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
    slotName: "pay_time",
    width: 200,
  },
  {
    dataIndex: "pay_way",
    title: "支付方式",
  },
  {
    title: "操作",
    slotName: "actions",
    fixed: "right",
    width: 80,
  },
];
</script>
<template>
  <SearchTable :request="getList" :columns="columns">
    <template #pay_time="{ record }">
      <a-tag v-if="!record.pay_time" color="blue">未支付</a-tag>
      <span v-else>{{ dateFormat(record.pay_time) }}</span>
    </template>
    <template #actions="{ record }">
      <a-link :key="record.id" status="danger">删除</a-link>
    </template>
  </SearchTable>
</template>
