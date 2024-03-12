<script lang="ts" setup>
import { computed } from "vue";
import { Message, type SwitchInstance } from "@arco-design/web-vue";
import type { BaseResponse } from "@gpt-vue/packages/type";

type OriginProps = SwitchInstance["$props"];

interface Props extends /* @vue-ignore */ OriginProps {
  modelValue: boolean | string | number;
  api: (params?: any) => Promise<BaseResponse<any>>;
}

const props = defineProps<Props>();

const emits = defineEmits(["update:modelValue"]);

const _value = computed({
  get: () => props.modelValue,
  set: (v) => {
    emits("update:modelValue", v);
  },
});

const onBeforeChange = async (params) => {
  try {
    await props.api({ ...params, value: !_value.value });
    Message.success("操作成功");
    return true;
  } catch (err) {
    console.log(err);
    return false;
  }
};
</script>
<template>
  <a-switch v-bind="{ ...props, ...$attrs }" v-model="_value" :before-change="onBeforeChange" />
</template>
