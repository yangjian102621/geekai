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
      <el-alert v-if="mobile !== ''" type="info" show-icon :closable="false" style="margin-bottom: 20px;">
        <p>请输入您参与众筹的 <strong style="color:#F56C6C">微信支付转账单号</strong> 兑换相应的对话次数。</p>
      </el-alert>

      <el-form :model="form">
        <el-form-item label="转账单号">
          <el-input v-model="form.tx_id"/>
        </el-form-item>
      </el-form>

      <el-form :model="form">
        <el-form-item label="兑换类别">
          <el-radio-group v-model="form.type">
            <el-radio label="chat" border>对话聊天</el-radio>
            <el-radio label="img" border>AI绘图</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
    </div>

    <template #footer>
      <span class="dialog-footer">
        <el-button type="primary" @click="save">
          确认核销
        </el-button>
      </span>
    </template>
  </el-dialog>
</template>

<script setup>
import {computed, ref} from "vue";
import {ElMessage} from "element-plus";
import {httpPost} from "@/utils/http";

const props = defineProps({
  show: Boolean,
  mobile: String
});

const showDialog = computed(() => {
  return props.show
})

const title = ref('众筹码核销')
const form = ref({
  tx_id: '',
  type: 'chat'
})

const emits = defineEmits(['hide']);

const save = () => {
  if (form.value.tx_id === '') {
    return ElMessage.error({message: "请输入微信支付转账单号"});
  }

  httpPost('/api/reward/verify', form.value).then(() => {
    ElMessage.success({
      message: '核销成功',
      duration: 1000,
      onClose: () => location.reload()
    })
  }).catch(e => {
    ElMessage.error({message: "核销失败：" + e.message});
  })
}

const close = function () {
  emits('hide', false);
}
</script>

<style scoped>

</style>