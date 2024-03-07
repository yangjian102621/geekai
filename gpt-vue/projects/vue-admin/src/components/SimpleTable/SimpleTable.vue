<script lang="ts" setup>
import { computed, ref, onActivated } from "vue";
import useAsyncTable from "./useAsyncTable";
import { useTableScroll } from "@/components/SearchTable/utils";
import { Message } from "@arco-design/web-vue";
import type { TableRequest, TableOriginalProps } from "./useAsyncTable";

export interface SimpleTable extends /* @vue-ignore */ TableOriginalProps {
  request: TableRequest<Record<string, unknown>>;
  params?: Record<string, unknown>;
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
  <div class="search-table">
    <div ref="tableContainerRef" class="search-table-container">
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
