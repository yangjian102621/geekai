<template>
  <el-dialog
      v-model="$props.show"
      :close-on-click-modal="false"
      :show-close="false"
      title="聊天配置"
  >
    <span>正在努力开发中...</span>
    <template #footer>
      <span class="dialog-footer">
        <el-button @click="close">关闭</el-button>
        <el-button type="primary" @click="save">
          保存
        </el-button>
      </span>
    </template>
  </el-dialog>
</template>

<script>
import {defineComponent} from "vue"

export default defineComponent({
  name: 'ConfigDialog',
  props: {
    show: {
      type: Boolean,
      default: true
    }
  },
  data() {
    return {}
  },
  methods: {
    save: function () {
      this.$emit('update:show', false);
    },
    close: function () {
      this.$emit('update:show', false);
    }
  }
})
</script>

<style lang="stylus">
.el-dialog {
  --el-dialog-width 90%;
  max-width 800px;
}
</style>