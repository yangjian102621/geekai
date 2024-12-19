<template>
  <el-select
    v-model="model"
    :placeholder="placeholder"
    :value="value"
    @change="$emit('update:value', $event)"
    style="--el-border-radius-base: 20px"
  >
    <el-option
      v-for="item in options"
      :key="item.value"
      :label="item.label"
      :value="item.value"
    >
    </el-option>
  </el-select>
</template>

<script>
export default {
  name: "BlackSelect",
  props: {
    value: {
      type: String,
      default: ""
    },
    placeholder: {
      type: String,
      default: "请选择"
    },
    options: {
      type: Array,
      default: []
    }
  },
  data() {
    return {
      model: this.value
    };
  }
};
</script>
