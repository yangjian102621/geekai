<template>

  <el-select v-model="model" :placeholder="placeholder"
             :value="value" @change="$emit('update:value', $event)"
             style="--el-fill-color-blank:#252020;
                   --el-text-color-regular: #a1a1a1;
                   --el-select-disabled-color:#0E0808;
                   --el-color-primary-light-9:#0E0808;
                   --el-border-radius-base:20px;
                   --el-border-color:#0E0808;">
    <el-option v-for="item in options"
               :key="item.value"
               :label="item.label"
               :value="item.value">
    </el-option>
  </el-select>

</template>

<script>
export default {
  name: 'BlackSelect',
  props: {
    value : {
      type: String,
      default: '',
    },
    placeholder: {
      type: String,
      default: '请选择',
    },
    options: {
      type: Array,
      default: []
    }
  },
  data() {
    return {
      model: this.value
    }
  }
}
</script>
