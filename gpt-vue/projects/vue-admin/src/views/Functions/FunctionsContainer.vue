<script lang="ts" setup>
import { Message, type TableColumnData } from "@arco-design/web-vue";
import SimpleTable from "@/components/SimpleTable/SimpleTable.vue";
import ConfirmSwitch from "@/components/ConfirmSwitch.vue";
import usePopup from "@/composables/usePopup";
import { getList, remove, setStatus, save } from "./api";
import FunctionsForm from "./FunctionsForm.vue";

const columns: TableColumnData[] = [
  {
    dataIndex: "name",
    title: "函数名称",
  },
  {
    dataIndex: "label",
    title: "函数别名",
  },
  {
    dataIndex: "description",
    title: "功能描述",
  },
  {
    dataIndex: "enabled",
    title: "启用状态",
    slotName: "switch",
  },
  {
    title: "操作",
    slotName: "actions",
    fixed: "right",
  },
];

const openFormModal = usePopup(FunctionsForm, {
  nodeProps: ([_, record]) => ({ record }),
  popupProps: ([reload, record], exposed) => ({
    title: `${record?.id ? "编辑" : "新增"}函数`,
    width: "800px",
    onBeforeOk: async (done) => {
      await exposed()?.handleSubmit(save, {
        id: record?.id,
        parameters: exposed()?.parameters(),
      });
      await reload();
      done(true);
    },
  }),
});

const handleRemove = async (id, reload) => {
  await remove({ id });
  Message.success("删除成功");
  await reload();
  return true;
};
</script>
<template>
  <a-button type="primary" @click="openFormModal">新增</a-button>
  <SimpleTable :request="getList" :columns="columns" :pagination="false">
    <template #switch="{ record, column }">
      <ConfirmSwitch
        v-model="record[column.dataIndex]"
        :api="(p) => setStatus({ ...p, id: record.id, filed: 'enabled' })"
      />
    </template>
    <template #exchange="{ record }">
      <a-tag v-if="record.exchange.calls > 0">聊天{{ record.exchange.calls }}次</a-tag>
      <a-tag v-else-if="record.exchange.img_calls > 0" color="green"
        >绘图{{ record.exchange.img_calls }}次</a-tag
      >
    </template>
    <template #actions="{ record, reload }">
      <a-link @click="openFormModal(reload, record)">编辑</a-link>
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
