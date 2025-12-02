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
        <el-form-item label="语音音色：">
          <el-select v-model="data.ttsModel" placeholder="请选择语音音色" @change="changeTTSModel">
            <el-option v-for="v in models" :value="v.id" :label="v.name" :key="v.id">
              {{ v.name }}
            </el-option>
          </el-select>
        </el-form-item>
      </el-form>
    </div>
  </el-dialog>
</template>

<script setup>
import {computed, ref, onMounted} from "vue"
import {useSharedStore} from "@/store/sharedata";
import {httpGet} from "@/utils/http";
const store = useSharedStore();

const data = ref({
  style: store.chatListStyle,
  stream: store.chatStream,
  ttsModel: store.ttsModel,
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
const models = ref([]);
onMounted(() => {
  // 获取模型列表
  httpGet("/api/model/list?type=tts").then((res) => {
    models.value = res.data;
    if (!data.ttsModel) {
      store.setTtsModel(models.value[0].id);
    }
  })
})

const changeTTSModel = (item) => {
  store.setTtsModel(item);
}
</script>

<style lang="stylus" scoped>
.chat-setting {

}
</style>