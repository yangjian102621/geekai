<template>
  <van-dialog v-model:show="showDialog"
              :title="title"
              :show-cancel-button="mobile !== ''"
              @confirm="save"
              @cancel="close">
    <van-cell-group inset>
      <van-field
          v-model="form.mobile"
          label="手机号"
          placeholder="请输入手机号"
      />
      <van-field
          v-model.number="form.code"
          center
          clearable
          label="短信验证码"
          placeholder="请输入短信验证码"
      >
        <template #button>
          <send-msg size="small" :mobile="form.mobile"/>
        </template>
      </van-field>
    </van-cell-group>
  </van-dialog>
</template>

<script setup>
import {computed, ref} from "vue";
import {httpPost} from "@/utils/http";
import {validateMobile} from "@/utils/validate";
import {showNotify} from "vant";
import SendMsg from "@/components/SendMsg.vue";

const props = defineProps({
  show: Boolean,
  mobile: String
});

const showDialog = computed(() => {
  return props.show
})

const title = ref('绑定手机号')
const form = ref({
  mobile: '',
  code: ''
})

const emits = defineEmits(['hide']);

const save = () => {
  if (!validateMobile(form.value.mobile)) {
    return showNotify({type: 'danger', message: '请输入正确的手机号码'});
  }
  if (form.value.code === '') {
    return showNotify({type: "danger", message: '请输入短信验证码'})
  }

  httpPost('/api/user/bind/mobile', form.value).then(() => {
    showNotify({type: 'success', message: '绑定成功', duration: 1000, onClose: emits('hide', false)});
  }).catch(e => {
    showNotify({type: 'danger', message: '绑定失败：' + e.message, duration: 2000});
  })
}

const close = function () {
  emits('hide', false);
}
</script>

<style scoped>

</style>