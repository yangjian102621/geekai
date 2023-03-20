<template>
  <el-dialog
      v-model="$props.show"
      title="聊天配置"
      width="30%"
      :before-close="beforeClose"
  >
    <span>正在努力开发中...</span>
    <template #footer>
      <span class="dialog-footer">
        <el-button @click="this.$emit('update:show', false)">取消</el-button>
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
    beforeClose: function () {
    },
    save: function () {
      this.$emit('update:show', false);
    }
  }
})
</script>

<style lang="stylus">

</style>