<template>
  <el-dialog
      v-model="showDialog"
      :close-on-click-modal="true"
      style="max-width: 400px"
      @close="close"
      :title="title"
  >
    <div class="third-login" v-loading="loading">
      <div class="item" v-if="wechatBindURL !== ''">
        <a class="link" :href="wechatBindURL"><i class="iconfont icon-wechat"></i></a>
        <span class="text ok" v-if="openid !== ''">已绑定</span>
        <span class="text" v-else>未绑定</span>
      </div>
    </div>
  </el-dialog>
</template>

<script setup>
import {computed, ref, watch} from "vue";
import {httpGet} from "@/utils/http";
import {checkSession} from "@/store/cache";
import {showMessageError} from "@/utils/dialog";

const props = defineProps({
  show: Boolean,
});
const emits = defineEmits(['hide']);

const showDialog = computed(() => {
  return props.show
})


const title = ref('绑定第三方登录')
const openid = ref('')
const wechatBindURL = ref('')
const loading = ref(true)

watch(showDialog, (val) => {
  if (val) {
    checkSession().then(user => {
      openid.value = user.openid
    })
    const returnURL = `${location.protocol}//${location.host}/login/callback?action=bind`
    httpGet("/api/user/clogin?return_url="+returnURL).then(res => {
      wechatBindURL.value = res.data.url
      loading.value = false
    }).catch(e => {
      showMessageError(e.message)
    })
  }
})

const close = function () {
  emits('hide');
}
</script>

<style lang="stylus" scoped>
.third-login {
  display flex
  justify-content center
  min-height 100px

  .item {
    display flex
    flex-flow column
    align-items center

    .link {
      display flex
      .iconfont {
        font-size 30px
        cursor pointer
        background #e9f1f6
        padding 10px
        border-radius 50%
      }
      margin-bottom 10px
    }


    .text {
      font-size 14px
    }

    .icon-wechat,.ok {
      color: #0bc15f;
    }
  }
}
</style>