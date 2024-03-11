<script lang="ts" setup>
import { getList, save, deleting, setStatus } from "./api";
import { ref } from "vue";
import ApiKeyForm from "./ApiKeyForm.vue";
import useCustomFormPopup from "@/composables/useCustomFormPopup";
import { Message } from "@arco-design/web-vue";
import SimpleTable from "@/components/SimpleTable/SimpleTable.vue";
import { dateFormat } from "@gpt-vue/packages/utils";
// table 配置
const columns = [
  {
    title: "所属平台",
    dataIndex: "platform",
  },
  {
    title: "名称",
    dataIndex: "name",
  },
  {
    title: "key",
    dataIndex: "value",
    slotName: "value",
  },
  {
    title: "用途",
    dataIndex: "type",
  },
  {
    title: "使用代理",
    dataIndex: "use_proxy",
    slotName: "proxy",
  },
  {
    title: "最后使用时间",
    dataIndex: "last_used_at",
    render: ({ record }) => {
      return dateFormat(record.last_used_at);
    },
  },
  {
    title: "启用状态",
    dataIndex: "enabled",
    slotName: "status",
  },
  {
    title: "操作",
    slotName: "action",
  },
];

// 数据
const tableData = ref([]);
const getData = () => {
  getList().then(({ code, data }) => {
    if (code === 0) {
      tableData.value = data;
    }
  });
};
getData();

//  新增编辑
const popup = useCustomFormPopup(ApiKeyForm, save, {
  popupProps: (arg) => ({ title: arg[0].record ? "编辑ApiKey" : "新增ApiKey" }),
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
        ><template #icon> <icon-plus /> </template>新增
      </a-button>
    </template>
    <template #value="{ record, column }">
      <a-typography-text copyable ellipsis style="margin: 0">
        {{ record[column.dataIndex] }}
      </a-typography-text>
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
    <template #proxy="{ record, reload }">
      <a-switch
        v-model="record.use_proxy"
        @change="
          (value) => {
            handleStatusChange({ filed: 'use_proxy', value, record, reload });
          }
        "
      />
    </template>
  </SimpleTable>
</template>
