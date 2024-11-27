<template>
  <el-dialog
      v-model="showDialog"
      :close-on-click-modal="true"
      :show-close="mobile !== ''"
      :before-close="close"
      :width="450"
      :title="title"
  >
    <div class="form" id="bind-mobile-form">
      <el-form :model="form">
        <el-form-item label="兑换码">
          <el-input v-model="form.code"/>
        </el-form-item>
      </el-form>
    </div>

    <template #footer>
      <span class="dialog-footer">
        <el-button type="primary" @click="save">
          立即兑换
        </el-button>
      </span>
    </template>
  </el-dialog>
</template>

<script setup>
import {computed, ref} from "vue";
import {ElMessage} from "element-plus";
import {httpPost} from "@/utils/http";
import {showMessageError, showMessageOK} from "@/utils/dialog";

const props = defineProps({
  show: Boolean,
});

const showDialog = computed(() => {
  return props.show
})

const title = ref('兑换码核销')
const form = ref({
  code: '',
})

const emits = defineEmits(['hide']);

const save = () => {
  if (form.value.code === '') {
    return ElMessage.error({message: "请输入兑换码"});
  }

  httpPost('/api/redeem/verify', form.value).then(() => {
    showMessageOK("兑换成功！")
    emits('hide', true)
  }).catch(e => {
    showMessageError("兑换失败：" + e.message)
  })
}

const close = function () {
  emits('hide', false);
}
</script>

<style scoped>

</style>