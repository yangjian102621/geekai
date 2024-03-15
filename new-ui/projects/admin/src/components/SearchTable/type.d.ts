import type { Component } from "vue";
import type { JsxElement } from "typescript";
import {
  DatePicker,
  Input,
  InputNumber,
  RadioGroup,
  RangePicker,
  Select,
  Switch,
  type TableColumnData,
} from "@arco-design/web-vue";
import type { TableOriginalProps, TableRequest } from "./useAsyncTable";

type Object = Record<string, unknown>;

export enum ValueType {
  "input" = Input,
  "select" = Select,
  "number" = InputNumber,
  "date" = DatePicker,
  "range" = RangePicker,
  "radio" = RadioGroup,
  "switch" = Switch,
}

export type SearchConfig = {
  valueType?: keyof typeof ValueType;
  fieldProps?: Object;
  render?: Component | JsxElement;
  slotsName?: string;
  defaultValue?: any;
  transform?: (value) => Record<string, any>;
};

export interface SearchTableColumns extends TableColumnData {
  search?: SearchConfig;
  hideInTable?: boolean;
  [key: string]: any;
}

export type SearchColumns = SearchTableColumns & { search: SearchConfig };

export interface SearchTableProps extends /* @vue-ignore */ TableOriginalProps {
  request: TableRequest<Object>;
  params?: Object;
  columns: SearchTableColumns[];
  headerTitle?: string;
}
