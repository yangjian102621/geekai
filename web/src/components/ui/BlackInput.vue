<template>
  <div class="black-input-wrapper">
    <el-input v-model="model" :type="type" :rows="rows"
              @input="onInput"
              style="--el-input-bg-color:#252020;
            --el-input-border-color:#414141;
            --el-input-focus-border-color:#414141;
            --el-text-color-regular: #f1f1f1;
            --el-input-border-radius: 10px;
            --el-border-color-hover:#616161"
              resize="none"
              :placeholder="placeholder" :maxlength="maxlength"/>
    <div class="word-stat" v-if="rows > 1">
      <span>{{value.length}}</span>/<span>{{maxlength}}</span>
    </div>
  </div>

</template>

<script setup>

import {ref, watch} from "vue";

const props = defineProps({
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
  },
  maxlength: {
    type: Number,
    default: 1024
  }
});
watch(() => props.value, (newValue) => {
  model.value = newValue
})
const model = ref(props.value)
const emits = defineEmits(['update:value']);
const onInput = (value) => {
  emits('update:value',value)
}
</script>

<style lang="stylus">
.black-input-wrapper {
  position relative

  .el-textarea__inner {
    padding: 20px;
    font-size: 16px;
  }

  .word-stat {
    position: absolute;
    bottom 10px
    right 10px
    color rgb(209 203 199)
    font-family: Neue Montreal, ui-sans-serif, system-ui, sans-serif, Apple Color Emoji, Segoe UI Emoji, Segoe UI Symbol, Noto Color Emoji;
    font-size .875rem
    line-height 1.25rem

    span {
      margin 0 1px
    }
  }
}
</style>