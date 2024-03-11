<script lang="ts" setup>
const props = defineProps({
  modelValue: {
    type: Array,
    default: () => [],
  },
});

const emits = defineEmits(["update:modelValue"]);

const handleCreateRow = () => {
  emits("update:modelValue", [
    ...(props.modelValue ?? []),
    { name: "", type: "", description: "", required: false },
  ]);
};

const handleRemoveRow = (index) => {
  emits(
    "update:modelValue",
    props.modelValue.filter((_, i) => i !== index)
  );
};
</script>
<template>
  <a-space direction="vertical" style="width: 100%">
    <a-table :data="modelValue" size="small" style="width: 100%" :pagination="false">
      <template #columns>
        <a-table-column title="参数名称" :width="120">
          <template #cell="scope">
            <a-input v-model="scope.record.name" />
          </template>
        </a-table-column>
        <a-table-column title="参数类型" :width="150">
          <template #cell="scope">
            <a-select
              v-model="scope.record.type"
              placeholder="参数类型"
              :options="['string', 'number']"
            />
          </template>
        </a-table-column>
        <a-table-column title="参数描述">
          <template #cell="scope">
            <a-input v-model="scope.record.description" />
          </template>
        </a-table-column>

        <a-table-column title="必填参数" :width="100" align="center">
          <template #cell="scope">
            <a-checkbox v-model="scope.record.required" />
          </template>
        </a-table-column>

        <a-table-column title="操作" :width="80">
          <template #cell="scope">
            <a-button
              status="danger"
              shape="circle"
              @click="handleRemoveRow(scope.rowIndex)"
              size="small"
            >
              <icon-delete />
            </a-button>
          </template>
        </a-table-column>
      </template>
    </a-table>
    <a-button type="primary" @click="handleCreateRow">
      <template #icon><icon-plus /></template>
      <span>新增参数</span>
    </a-button>
  </a-space>
</template>
