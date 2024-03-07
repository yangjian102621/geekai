<script lang="ts" setup>
import { computed, ref, type PropType } from "vue";
import { getDefaultFormData, useComponentConfig } from "./utils";
import { ValueType } from "./type.d";
import type { SearchTableColumns, SearchColumns } from "./type";

const props = defineProps({
  modelValue: {
    type: Object,
    default: () => ({}),
  },
  columns: {
    type: Array as PropType<SearchTableColumns[]>,
    default: () => [],
  },
  submitting: {
    type: Boolean,
    default: false,
  },
});

const emits = defineEmits(["update:modelValue", "request"]);

const size = "small";

const collapsed = ref(false);

const formData = computed({
  get: () => props.modelValue,
  set(value) {
    emits("update:modelValue", value);
  },
});

const searchColumns = computed(() => {
  return props.columns?.filter(
    (item) => item.dataIndex && item.search
  ) as (SearchColumns & { dataIndex: string })[];
});

const optionsEvent = {
  onReset: () => {
    formData.value = getDefaultFormData(props.columns);
    emits("request");
  },
  onSearch: () => emits("request"),
  onCollapse: (value: boolean) => {
    collapsed.value = value ?? !collapsed.value;
  },
};
</script>
<template>
  <AForm
    v-if="searchColumns?.length"
    class="search-form-conteiner"
    :model="formData"
    :size="size"
    :label-col-props="{ span: 0 }"
    :wrapper-col-props="{ span: 24 }"
    @submit="optionsEvent.onSearch"
  >
    <AGrid
      :cols="{ md: 1, lg: 2, xl: 3, xxl: 5 }"
      :row-gap="12"
      :col-gap="12"
      :collapsed="collapsed"
    >
      <AGridItem
        v-for="item in searchColumns"
        :key="item.dataIndex"
        style="transition: all 0.3s ease-in-out"
      >
        <AFormItem :field="item.dataIndex" :label="(item.title as string)">
          <slot :name="item.search.slotsName">
            <component
              v-model="formData[item.dataIndex]"
              :is="
                ValueType[item.search.valueType ?? 'input'] ??
                item.search.render
              "
              v-bind="useComponentConfig(size, item)"
            />
          </slot>
        </AFormItem>
      </AGridItem>
      <AGridItem suffix>
        <ASpace class="flex-end">
          <slot name="search-options" :option="optionsEvent">
            <AButton
              type="primary"
              html-type="submit"
              :size="size"
              :loading="submitting"
            >
              <icon-search />
              <span>查询</span>
            </AButton>
            <AButton
              :size="size"
              @click="optionsEvent.onReset"
              :loading="submitting"
            >
              <icon-refresh />
              <span>重置</span>
            </AButton>
          </slot>
          <slot name="search-extra" />
        </ASpace>
      </AGridItem>
    </AGrid>
  </AForm>
</template>
<style scoped>
.search-form-conteiner {
  padding: 16px 0;
}
.flex-end {
  display: flex;
  justify-content: end;
}
</style>
