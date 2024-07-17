<template>
  <div class="black-input-wrapper">
    <el-input v-model="model" :type="type" :rows="rows"
              @input="$emit('update:value', $event)"
              style="--el-input-bg-color:#252020;
            --el-input-border-color:#414141;
            --el-input-focus-border-color:#414141;
            --el-text-color-regular: #f1f1f1;
            --el-input-border-radius: 10px;
            --el-border-color-hover:#616161"
              resize="none"
              :placeholder="placeholder"/>
  </div>

</template>

<script>
export default {
  name: 'BlackInput',
  props: {
    value : {
      type: String,
      default: '',
    },
    placeholder: {
      type: String,
      default: '',
    },
    type: {
      type: String,
      default: 'input',
    },
    rows: {
      type: Number,
      default: 5,
    }
  },
  data() {
    return {
      model: this.value
    }
  }
}
</script>

<style lang="stylus">
.black-input-wrapper {
  .el-textarea__inner {
    padding: 20px;
    font-size: 16px;
  }
}
</style>