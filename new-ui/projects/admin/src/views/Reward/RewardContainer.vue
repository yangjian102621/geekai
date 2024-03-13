<script lang="ts" setup>
import { Message, type TableColumnData } from "@arco-design/web-vue";
import { dateFormat } from "@chatgpt-plus/packages/utils";
import SimpleTable from "@/components/SimpleTable/SimpleTable.vue";
import { getList, remove } from "./api";

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
    width: 80,
  },
];

const handleRemove = async (id, reload) => {
  await remove({ id });
  Message.success("删除成功");
  await reload();
  return true;
};
</script>
<template>
  <SimpleTable :request="getList" :columns="columns">
    <template #updated_at="{ record }">
      <span v-if="record.status">{{ dateFormat(record.updated_at) }}</span>
      <a-tag v-else color="blue">未核销</a-tag>
    </template>
    <template #exchange="{ record }">
      <a-tag v-if="record.exchange.calls > 0">聊天{{ record.exchange.calls }}次</a-tag>
      <a-tag v-else-if="record.exchange.img_calls > 0" color="green"
        >绘图{{ record.exchange.img_calls }}次</a-tag
      >
    </template>
    <template #actions="{ record, reload }">
      <a-popconfirm
        content="是否删除？"
        position="left"
        type="warning"
        :on-before-ok="() => handleRemove(record.id, reload)"
      >
        <a-link status="danger">删除</a-link>
      </a-popconfirm>
    </template>
  </SimpleTable>
</template>
