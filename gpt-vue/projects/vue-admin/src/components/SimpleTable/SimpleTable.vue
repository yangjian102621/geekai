<script lang="ts" setup>
import { computed, ref, onActivated } from "vue";
import useAsyncTable from "./useAsyncTable";
import { useTableScroll } from "@/components/SearchTable/utils";
import { Message } from "@arco-design/web-vue";
import type { TableRequest, TableOriginalProps } from "./useAsyncTable";

interface SimpleTable extends /* @vue-ignore */ TableOriginalProps {
  request: TableRequest<Record<string, unknown>>;
  params?: Record<string, unknown>;
  columns?: TableOriginalProps["columns"];
}

const props = defineProps<SimpleTable>();
const tableContainerRef = ref<HTMLElement>();

// 表格请求参数
const [tableConfig, getList] = useAsyncTable(props.request, props.params);

const _columns = computed(() => {
  return props.columns?.map((item) => ({
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
  <div class="simple-header">
    <a-space>
      <slot name="header" v-bind="{ reload: handleSearch }" />
    </a-space>
  </div>
  <div class="simple-table">
    <div ref="tableContainerRef" class="simple-table-container">
      <ATable
        v-bind="{
          ...$attrs,
          ...tableConfig,
          ...props,
          scroll: useTableScroll(_columns || [], tableContainerRef as HTMLElement),
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
.simple-table {
  display: flex;
  flex-direction: column;
  height: 100%;
}
.simple-table-container {
  flex: 1;
}
.simple-table-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.simple-header {
  padding: 16px 0;
}
</style>
