<template>
  <el-switch
    v-model="model"
    :size="size"
    @change="$emit('update:value', $event)"
  />
</template>

<script setup>
import { ref, watch } from "vue";
const props = defineProps({
  value: Boolean,
  size: {
    type: String,
    default: "default"
  }
});
const model = ref(props.value);

watch(
  () => props.value,
  (newValue) => {
    model.value = newValue;
  }
);
</script>
