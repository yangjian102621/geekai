import type { TableColumnData } from "@arco-design/web-vue";
import type { SearchTableColumns, SearchColumns } from "./type";

export function useTableXScroll(columns: TableColumnData[]) {
  return columns.reduce((prev, curr) => {
    const width = curr.width ?? 150;
    return prev + width;
  }, 0);
}

export function useTableScroll(columns: SearchTableColumns[], container?: HTMLElement) {
  const x = columns.reduce((prev, curr) => {
    const width = curr.hideInTable ? 0 : curr.width ?? 150;
    return prev + width;
  }, 0);
  const y = container?.clientHeight ?? undefined;
  return { x, y };
}

export function getDefaultFormData(columns: SearchTableColumns[]) {
  return columns?.reduce((field, curr) => {
    if (curr.dataIndex && curr?.search?.defaultValue) {
      field[curr.dataIndex] = curr.search.defaultValue;
    }
    return field;
  }, {});
}

export function useRequestParams(
  columns: SearchTableColumns[],
  originFormData: Record<string, any>
) {
  const filterFormData = columns?.reduce((prev, curr) => {
    if (!curr.dataIndex || !curr.search) {
      return prev;
    }
    if (curr?.search?.transform) {
      const filters = curr.search.transform(originFormData[curr.dataIndex]);
      return Object.assign(prev, filters);
    }
    return Object.assign(prev, { [curr.dataIndex]: originFormData[curr.dataIndex] });
  }, {});
  return filterFormData as Record<string, any>;
}

export function useComponentConfig(size: string, item: SearchColumns) {
  return {
    size,
    placeholder: item.search.valueType === "range" ? ["开始时间", "结束时间"] : item.title,
    allowClear: true,
    ...(item.search.fieldProps ?? {}),
  };
}
