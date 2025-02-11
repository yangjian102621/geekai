<template>
  <el-dialog
      class="config-dialog"
      v-model="showDialog"
      :close-on-click-modal="true"
      :before-close="close"
      style="max-width: 600px"
      title="聊天配置"
  >
    <div class="chat-setting">
      <el-form :model="data" label-width="100px" label-position="left">
        <el-form-item label="聊天样式：">
          <el-radio-group v-model="data.style" @change="(val) => {store.setChatListStyle(val)}">
            <el-radio value="list">列表样式</el-radio>
            <el-radio value="chat">对话样式</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="流式输出：">
          <el-switch v-model="data.stream" @change="(val) => {store.setChatStream(val)}" />
        </el-form-item>
      </el-form>
    </div>
  </el-dialog>
</template>

<script setup>
import {computed, ref} from "vue"
import {useSharedStore} from "@/store/sharedata";
const store = useSharedStore();

const data = ref({
  style: store.chatListStyle,
  stream: store.chatStream,
})
// eslint-disable-next-line no-undef
const props = defineProps({
  show: Boolean,
});

const showDialog = computed(() => {
  return props.show
})
const emits = defineEmits(['hide']);
const close = function () {
  emits('hide', false);
}
</script>

<style lang="stylus" scoped>
.chat-setting {

}
</style>