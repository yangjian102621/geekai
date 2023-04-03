<template>
  <el-dialog
      v-model="$props.show"
      :close-on-click-modal="false"
      :show-close="false"
      title="聊天配置"
  >
    <el-row class="row-center">
      <span>功能正在开发中，有什么使用建议可以通过下面的方式联系作者。</span>
    </el-row>
    <el-row>
      <el-col :span="12">
        <el-image :src="wechatGroup"></el-image>
      </el-col>

      <el-col :span="12">
        <el-image :src="wechatCard"></el-image>
      </el-col>
    </el-row>

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
    return {
      wechatGroup: "https://img.r9it.com/chatgpt/wechat-group.jpeg",
      wechatCard: "https://img.r9it.com/chatgpt/wechat-card.jpeg"
    }
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