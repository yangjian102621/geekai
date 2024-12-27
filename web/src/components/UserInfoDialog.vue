<template>
  <el-dialog class="config-dialog" v-model="showDialog" :close-on-click-modal="true" :before-close="close" style="max-width: 400px" title="账户信息">
    <div class="flex-center-col p-4 pt-0" id="user-info">
      <user-profile @hide="close" />
    </div>
  </el-dialog>
</template>

<script setup>
import { computed } from "vue";
import UserProfile from "@/components/UserProfile.vue";

// eslint-disable-next-line no-undef
const props = defineProps({
  show: Boolean,
});

const showDialog = computed(() => {
  return props.show;
});

// eslint-disable-next-line no-undef
const emits = defineEmits(["hide"]);
const close = function () {
  emits("hide", false);
};
</script>

<style lang="stylus">
.config-dialog {
  .el-dialog {
    --el-dialog-width 90%;
    max-width 800px;
  }
}
</style>
