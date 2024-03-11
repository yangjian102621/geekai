import { computed, onMounted, reactive, unref } from "vue";
import type { TableInstance } from "@arco-design/web-vue";
import type { BaseResponse } from "@gpt-vue/packages/type";

export type TableOriginalProps = TableInstance["$props"];
export type TableRequest<T extends Record<string, unknown>> = (
  params?: any
) => Promise<BaseResponse<T[]>>;
export type TableReturn = [TableOriginalProps, () => Promise<void>];
function useAsyncTable<T extends Record<string, unknown>>(
  request: TableRequest<T>,
  params?: Record<string, unknown>
): TableReturn {
  const tableState = reactive<{ loading: Boolean; data: T[] }>({
    loading: false,
    data: [],
  });

  const tableConfig = computed<TableOriginalProps>(() => {
    return {
      ...tableState,
      rowKey: "id",
    };
  });

  const getTableData = async () => {
    tableState.loading = true;
    try {
      const { data } = await request({
        ...unref(params ?? {}),
      });
      tableState.data = data as any;
    } finally {
      tableState.loading = false;
    }
  };

  onMounted(getTableData);

  return [tableConfig, getTableData] as TableReturn;
}

export default useAsyncTable;
