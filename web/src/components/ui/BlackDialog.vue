<template>
  <div class="black-dialog">
    <el-dialog
        v-model="showDialog"
        style="--el-dialog-bg-color:#414141;
        --el-text-color-primary:#f1f1f1;
        --el-border-color:#414141;
        --el-color-primary:#21aa93;
        --el-color-primary-dark-2:#41555d;
        --el-color-white: #e1e1e1;
        --el-color-primary-light-3:#549688;
        --el-fill-color-blank:#616161;
         --el-color-primary-light-7:#717171;
        --el-color-primary-light-9:#717171;
        --el-text-color-regular:#e1e1e1"
        :title="title"
        :width="width"
        :before-close="cancel"
    >
      <div class="dialog-body">
        <slot></slot>
      </div>
      <template #footer v-if="!hideFooter">
        <div class="dialog-footer">
          <el-button @click="cancel">{{cancelText}}</el-button>
          <el-button type="primary" @click="$emit('confirm')" v-if="!hideConfirm">{{confirmText}}</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>

import {ref, watch} from "vue";
const props = defineProps({
  show : Boolean,
  title: {
    type: String,
    default: 'Tips',
  },
  width: {
    type: String,
    default: 'auto',
  },
  hideFooter:{
    type: Boolean,
    default: false
  },
  hideConfirm:{
    type: Boolean,
    default: false
  },
  confirmText: {
    type: String,
    default: '确定',
  },
  cancelText: {
    type: String,
    default: '取消',
  },
});
const emits = defineEmits(['confirm','cancal']);
const showDialog = ref(props.show)

watch(() => props.show, (newValue) => {
  showDialog.value = newValue
})
const cancel = () => {
  showDialog.value = false
  emits('cancal')
}
</script>

<style lang="stylus">
.black-dialog {
  .dialog-body {
    .form {
      .form-item {
        display flex
        flex-flow column
        font-family: "Neue Montreal";
        padding 10px 0

        .label {
          margin-bottom 0.6rem
          margin-inline-end 0.75rem
          color #ffffff
          font-size 1rem
          font-weight 500
        }

        .input {
          display flex
          padding 10px
          text-align left
          font-size 1rem
          background none
          border-radius 0.375rem
          border 1px solid #8f8f8f
          outline: none;
          transition: border-color 0.5s ease, box-shadow 0.5s ease;

          &:focus {
            border-color: #0F7A71;
            box-shadow: 0 0 5px #0F7A71;
          }
        }
      }
    }
  }
}

</style>
