<template>

  <el-switch v-model="model" :size="size"
             @change="$emit('update:value', $event)"
             style="--el-switch-on-color:#555555;--el-color-white:#0E0808"/>

</template>

<script setup>

import {ref, watch} from "vue";
const props = defineProps({
  value : Boolean,
  size: {
    type: String,
    default: 'default',
  }
});
const model = ref(props.value)

watch(() => props.value, (newValue) => {
  model.value = newValue
})
</script>
