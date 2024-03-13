<script lang="ts" setup>
import { ref, h } from "vue";
import { Message, Modal } from "@arco-design/web-vue";
import { dateFormat } from "@gpt-vue/packages/utils";
import SearchTable from "@/components/SearchTable/SearchTable.vue";
import type { SearchTableColumns } from "@/components/SearchTable/type";
import app from "@/main";
import { getList, message, remove } from "./api";
import ChatsLogs from "./ChatsLogs.vue";

const columns: SearchTableColumns[] = [
  {
    dataIndex: "user_id",
    title: "账户ID",
    search: {
      valueType: "input",
    },
  },
  {
    dataIndex: "username",
    title: "账户",
  },
  {
    dataIndex: "title",
    title: "标题",
    search: {
      valueType: "input",
    },
  },
  {
    dataIndex: "model",
    title: "模型",
    search: {
      valueType: "input",
    },
  },
  {
    dataIndex: "msg_num",
    title: "消息数量",
  },
  {
    dataIndex: "token",
    title: "消耗算力",
  },
  {
    dataIndex: "username",
    title: "账户",
  },
  {
    dataIndex: "created_at",
    title: "创建时间",
    search: {
      valueType: "range",
    },
    render: ({ record }) => dateFormat(record.created_at),
  },
  {
    title: "操作",
    fixed: "right",
    slotName: "actions",
  },
];

const tabsList = [
  { key: "1", title: "对话列表", api: getList, columns },
  { key: "2", title: "消息记录", api: message, columns },
];

const activeKey = ref(tabsList[0].key);

const handleRemove = async (chat_id, reload) => {
  await remove({ chat_id });
  Message.success("删除成功");
  await reload();
  return true;
};

const handleCheck = (record) => {
  if (activeKey.value === "1") {
    Modal._context = app._context;
    Modal.info({
      title: "对话详情",
      width: 800,
      content: () => h(ChatsLogs, { id: record.chat_id }),
    });
    return;
  }
  Modal.info({
    title: "消息详情",
    content: record.content,
  });
};
</script>
<template>
  <a-tabs v-model:active-key="activeKey" lazy-load justify>
    <a-tab-pane v-for="item in tabsList" :key="item.key" :title="item.title">
      <SearchTable :request="item.api" :columns="item.columns">
        <template #actions="{ record, reload }">
          <a-link @click="handleCheck(record)">查看</a-link>
          <a-popconfirm
            content="是否删除？"
            position="left"
            type="warning"
            :on-before-ok="() => handleRemove(record.id, reload)"
          >
            <a-link status="danger">删除</a-link>
          </a-popconfirm>
        </template>
      </SearchTable>
    </a-tab-pane>
  </a-tabs>
</template>
