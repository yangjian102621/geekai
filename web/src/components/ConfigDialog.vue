<template>
  <div>
    <el-dialog
        v-show="show"
        title="聊天配置"
        width="30%"
        :before-close="beforeClose"
    >
      <span>正在努力开发中...</span>
      <template #footer>
      <span class="dialog-footer">
        <el-button @click="show = false">取消</el-button>
        <el-button type="primary" @click="show = false">
          保存
        </el-button>
      </span>
      </template>
    </el-dialog>
  </div>
</template>

<script>
import {defineComponent} from "vue"

export default defineComponent({
  name: 'ConfigDialog',
  data() {
    return {
      show: true
    }
  },
  methods: {
    beforeClose: function () {

    }
  }
})
</script>

<style lang="stylus">

</style>