<script lang="ts" setup>
import { computed, ref, onActivated } from "vue";
import useAsyncTable from "./useAsyncTable";
import FormSection from "./FormSection.vue";
import type { SearchTableProps } from "./type";
import { useTableScroll, getDefaultFormData, useRequestParams } from "./utils";
import { Message } from "@arco-design/web-vue";

const props = defineProps<SearchTableProps>();
const formData = ref({ ...getDefaultFormData(props.columns) });
const tableContainerRef = ref<HTMLElement>();

// 表格请求参数
const requestParams = computed(() => ({
  ...useRequestParams(props.columns, formData.value),
  ...props.params,
}));

const [tableConfig, getList] = useAsyncTable(props.request, requestParams);

const _columns = computed(() => {
  return props.columns.map((item) => ({
    ellipsis: true,
    tooltip: true,
    ...item,
  }));
});

const handleSearch = async (tips?: boolean) => {
  tips && Message.success("操作成功");
  await getList();
};

onActivated(handleSearch);
</script>
<template>
  <div class="search-table">
    <div class="search-table-header">
      <div>
        <slot name="header-title">{{ props.headerTitle }}</slot>
      </div>
      <div class="header-option">
        <slot
          name="header-option"
          :formData="formData"
          :reload="handleSearch"
        />
      </div>
    </div>
    <FormSection
      v-model="formData"
      :columns="columns"
      :submitting="(tableConfig.loading as boolean)"
      @request="handleSearch"
    />
    <div ref="tableContainerRef" class="search-table-container">
      <ATable
        v-bind="{
          ...$attrs,
          ...tableConfig,
          ...props,
          scroll: useTableScroll(_columns, tableContainerRef as HTMLElement),
          columns: _columns
        }"
      >
        <template v-for="slot in Object.keys($slots)" #[slot]="config">
          <slot :name="slot" v-bind="{ ...config, reload: handleSearch }" />
        </template>
      </ATable>
    </div>
  </div>
</template>
<style scoped>
.search-table {
  display: flex;
  flex-direction: column;
  height: 100%;
}
.search-table-container {
  flex: 1;
}
.search-table-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}
</style>
