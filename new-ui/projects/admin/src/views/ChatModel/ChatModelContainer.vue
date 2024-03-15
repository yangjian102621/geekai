<script lang="ts" setup>
import { getList, save, deleting, setStatus } from "./api";
import ChatModelForm from "./ChatModelForm.vue";
import useCustomFormPopup from "@/composables/useCustomFormPopup";
import { Message, type TableColumnData } from "@arco-design/web-vue";
import SimpleTable from "@/components/SimpleTable/SimpleTable.vue";
import { dateFormat } from "@chatgpt-plus/packages/utils";
// table 配置
const columns: TableColumnData[] = [
  {
    title: "所属平台",
    dataIndex: "platform",
  },
  {
    title: "模型名称",
    dataIndex: "name",
  },
  {
    title: "模型值",
    dataIndex: "value",
  },
  {
    title: "对话权重",
    dataIndex: "weight",
  },
  {
    title: "启用状态",
    dataIndex: "enabled",
    slotName: "status",
  },
  {
    title: "开放状态",
    dataIndex: "open",
    slotName: "open",
  },
  {
    title: "创建时间",
    dataIndex: "created_at",
    render: ({ record }) => {
      return dateFormat(record.created_at);
    },
  },
  {
    title: "操作",
    slotName: "action",
    width: 120,
    fixed: "right",
  },
];

//  新增编辑
const popup = useCustomFormPopup(ChatModelForm, save, {
  popupProps: (arg) => ({ title: arg[0].record ? "编辑模型" : "新增模型" }),
});

// 删除
const handleDelete = ({ id }, reload) => {
  deleting(id).then(({ code }) => {
    if (code === 0) {
      Message.success("操作成功");
      reload();
    }
  });
};

// 状态
const handleStatusChange = ({ filed, value, record, reload }) => {
  setStatus({
    id: record.id,
    value,
    filed,
  }).then(({ code }) => {
    if (code === 0) {
      Message.success("操作成功");
      reload();
    }
  });
};
</script>
<template>
  <SimpleTable :columns="columns" :request="getList">
    <template #action="{ record, reload }">
      <a-link @click="popup({ record, reload })">编辑</a-link>
      <a-popconfirm content="确定删除？" @ok="handleDelete(record, reload)">
        <a-link status="danger">删除</a-link>
      </a-popconfirm>
    </template>
    <template #header="{ reload }">
      <a-button @click="popup({ reload })" size="small" type="primary"
        ><template #icon> <icon-plus /> </template>新增</a-button
      >
    </template>
    <template #status="{ record, reload }">
      <a-switch
        v-model="record.enabled"
        @change="
          (value) => {
            handleStatusChange({ filed: 'enabled', value, record, reload });
          }
        "
      />
    </template>
    <template #open="{ record, reload }">
      <a-switch
        v-model="record.open"
        @change="
          (value) => {
            handleStatusChange({ filed: 'open', value, record, reload });
          }
        "
      />
    </template>
  </SimpleTable>
</template>
