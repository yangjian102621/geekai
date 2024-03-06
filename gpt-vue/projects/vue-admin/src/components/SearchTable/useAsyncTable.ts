import { computed, onMounted, reactive, unref, type Ref } from "vue";
import type { TableInstance } from "@arco-design/web-vue";
import type { BaseResponse, ListResponse } from "@gpt-vue/packages/type";

export type TableOriginalProps = TableInstance["$props"];
export type TableRequest<T extends Record<string, unknown>> = (params?: any) => Promise<BaseResponse<ListResponse<T>>>
export type TableReturn = [TableOriginalProps, () => Promise<void>];
function useAsyncTable<T extends Record<string, unknown>>(
  request: TableRequest<T>,
  params?: Ref<Record<string, unknown>>
): TableReturn {
  const paginationState = reactive({
    current: 1,
    pageSize: 10,
    total: 0,
  });

  const tableState = reactive({
    loading: false,
    data: []
  })

  const tableConfig = computed<TableOriginalProps>(() => {
    return {
      ...tableState,
      rowKey: "id",
      pagination: {
        ...paginationState,
        showTotal: true,
        showPageSize: true,
      },
      onPageChange: (page) => {
        paginationState.current = page;
        getTableData();
      },
      onPageSizeChange(pageSize) {
        paginationState.pageSize = pageSize;
        getTableData();
      },
    };
  });

  const getTableData = async () => {
    tableState.loading = true
    try {
      const { data } = await request({
        ...unref(params ?? {}),
        page: paginationState.current,
        pageSize: paginationState.pageSize,
      });
      tableState.data = (data as any)?.items;
      paginationState.total = (data as any)?.total;
    } finally {
      tableState.loading = false
    }
  };

  onMounted(getTableData);

  return [tableConfig, getTableData] as TableReturn;
}

export default useAsyncTable;
