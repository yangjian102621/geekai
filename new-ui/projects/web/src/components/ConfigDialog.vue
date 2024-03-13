<template>
  <el-dialog
      class="config-dialog"
      v-model="showDialog"
      :close-on-click-modal="true"
      :before-close="close"
      style="max-width: 600px"
      title="账户信息"
  >
    <div class="user-info" id="user-info">
      <el-form v-if="user.id" :model="user" label-width="150px">
        <el-form-item label="账户">
          <span>{{ user.username }}</span>
        </el-form-item>
        <el-form-item label="剩余算力">
          <el-tag>{{ user['power'] }}</el-tag>
        </el-form-item>
        <el-form-item label="会员到期时间" v-if="user['expired_time']  > 0">
          <el-tag type="danger">{{ dateFormat(user['expired_time']) }}</el-tag>
        </el-form-item>
      </el-form>
    </div>
  </el-dialog>
</template>

<script setup>
import {computed, onMounted, ref} from "vue"
import {httpGet, httpPost} from "@/utils/http";
import {ElMessage} from "element-plus";
import {Plus} from "@element-plus/icons-vue";
import Compressor from "compressorjs";
import {dateFormat} from "@/utils/libs";

// eslint-disable-next-line no-undef
const props = defineProps({
  show: Boolean,
  user: Object,
  models: Array,
});

const showDialog = computed(() => {
  return props.show
})
const user = ref({
  username: '',
  nickname: '',
  avatar: '',
  calls: 0,
  tokens: 0,
})

onMounted(() => {
  // 获取最新用户信息
  httpGet('/api/user/profile').then(res => {
    user.value = res.data
  }).catch(e => {
    ElMessage.error("获取用户信息失败：" + e.message)
  });
})

// eslint-disable-next-line no-undef
const emits = defineEmits(['hide']);
const close = function () {
  emits('hide', false);
}
</script>

<style lang="stylus">
.config-dialog {
  .el-dialog {
    --el-dialog-width 90%;
    max-width 800px;

    .el-dialog__body {
      overflow-y auto;

      .user-info {
        position relative;

        .el-message {
          position: absolute;
        }
      }

      .tip {
        color #c1c1c1
        font-size 12px;
      }
    }
  }
}
</style>