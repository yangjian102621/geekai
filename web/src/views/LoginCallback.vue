<template>
  <div class="login-callback"
       v-loading="loading"
       element-loading-text="正在同步登录信息..."
       :style="{ height: winHeight + 'px' }">

    <el-dialog
        v-model="show"
        :close-on-click-modal="false"
        :show-close="false"
        style="width: 360px;"
    >
      <el-result
          icon="success"
          title="登录成功"
          style="--el-result-padding:10px"
      >
          <template #sub-title>
            <div class="user-info">
              <div class="line">您的初始账户信息如下：</div>
              <div class="line"><span>用户名：</span>{{username}}</div>
              <div class="line"><span>密码：</span>{{password}}</div>
              <div class="line">您后期也可以通过此账号和密码登录</div>
            </div>
          </template>
          <template #extra>
            <el-button type="primary" @click="finishLogin">我知道了</el-button>
          </template>
      </el-result>
    </el-dialog>

  </div>
</template>

<script setup>
import {ref} from "vue"
import {useRouter} from "vue-router"
import {ElMessage, ElMessageBox} from "element-plus";
import {httpGet} from "@/utils/http";
import {setUserToken} from "@/store/session";
import {isMobile} from "@/utils/libs";

const winHeight = ref(window.innerHeight)
const loading = ref(true)
const router = useRouter()
const show = ref(false)
const username = ref('')
const password = ref('')


const code = router.currentRoute.value.query.code
if (code === "") {
  ElMessage.error({message: "登录失败：code 参数不能为空",duration: 2000, onClose: () => router.push("/")})
} else {
  // 发送请求获取用户信息
  httpGet("/api/user/clogin/callback",{login_type: "wx",code: code}).then(res => {
    setUserToken(res.data.token)
    if (res.data.username) {
      username.value = res.data.username
      password.value = res.data.password
      show.value = true
      loading.value = false
    } else {
      finishLogin()
    }
  }).catch(e => {
    ElMessageBox.alert(e.message, {
      confirmButtonText: '重新登录',
      type:"error",
      title:"登录失败",
      callback: () => {
        router.push("/login")
      },
    })
  })
}
const finishLogin = () => {
  if (isMobile()) {
    router.push('/mobile')
  } else {
    router.push('/chat')
  }
}
</script>

<style lang="stylus" scoped>
.login-callback {
  .user-info {
    display flex
    flex-direction column
    padding 10px
    border 1px dashed #e1e1e1
    border-radius 10px
    .line {
      text-align left
      font-size 14px
      line-height 1.5

      span{
        font-weight bold
      }
    }
  }
}
</style>
