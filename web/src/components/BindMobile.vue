<template>
  <el-dialog
      v-model="showDialog"
      :close-on-click-modal="false"
      :show-close="false"
      :title="title"
  >
    <div class="form" id="password-form">
      <el-form :model="form" label-width="120px">
        <el-form-item label="手机号码">
          <el-input v-model="form.mobile" type="password"/>
        </el-form-item>
        <el-form-item label="手机验证码">
          <el-row :gutter="10">
            <el-col :span="12">
              <el-input v-model="form.code" type="password"/>
            </el-col>
            <el-col :span="12">
              <send-msg size="large" v-model:mobile="form.mobile"/>
            </el-col>
          </el-row>

        </el-form-item>
      </el-form>
    </div>

    <template #footer>
      <span class="dialog-footer">
        <el-button type="primary" @click="save">
          提交绑定
        </el-button>
      </span>
    </template>
  </el-dialog>
</template>

<script setup>
import {ref} from "vue";
import {checkSession} from "@/action/session";
import SendMsg from "@/components/SendMsg.vue";

const title = ref('绑定手机号')
const showDialog = ref(false)
const form = ref({
  mobile: '',
  code: ''
})

checkSession().then(user => {
  if (user.mobile === '') {
    showDialog.value = true
  }
}).catch(e => {
  console.log(e.message)
})

const save = () => {

}

const sendMsg = () => {

}
</script>

<style scoped>

</style>