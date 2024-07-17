<template>
  <div class="index-page" :style="{height: winHeight+'px'}">
    <div class="index-bg" :style="{backgroundImage: 'url('+bgImgUrl+')'}"></div>
    <div class="menu-box">
      <el-menu
          mode="horizontal"
          :ellipsis="false"
      >
        <div class="menu-item">
          <el-image :src="logo" alt="Geek-AI"/>
          <div class="title">{{ title }}</div>
        </div>
        <div class="menu-item">
          <span v-if="!license.de_copy">
            <a :href="docsURL" target="_blank">
            <el-button type="primary" round>
              <i class="iconfont icon-book"></i>
              <span>文档</span>
            </el-button>
          </a>

          <a :href="gitURL" target="_blank">
            <el-button type="success" round>
              <i class="iconfont icon-github"></i>
              <span>源码</span>
            </el-button>
          </a>
          </span>

          <span v-if="!isLogin">
            <el-button @click="router.push('/login')" round>登录</el-button>
            <el-button @click="router.push('/register')" round>注册</el-button>
          </span>
        </div>
      </el-menu>
    </div>
    <div class="content">
      <h1>欢迎使用 {{ title }}</h1>
      <p>{{ slogan }}</p>
      <el-button @click="router.push('/chat')" color="#ffffff" style="color:#007bff" :dark="false">
        <i class="iconfont icon-chat"></i>
        <span>AI 对话</span>
      </el-button>
      <el-button @click="router.push('/mj')" color="#C4CCFD" style="color:#424282" :dark="false">
        <i class="iconfont icon-mj"></i>
        <span>MJ 绘画</span>
      </el-button>

      <el-button @click="router.push('/sd')" color="#4AE6DF" style="color:#424282" :dark="false">
        <i class="iconfont icon-sd"></i>
        <span>SD 绘画</span>
      </el-button>
      <el-button @click="router.push('/xmind')" color="#FFFD55" style="color:#424282" :dark="false">
        <i class="iconfont icon-xmind"></i>
        <span>思维导图</span>
      </el-button>
      <!--      <div id="animation-container"></div>-->
    </div>

    <footer-bar />
  </div>
</template>

<script setup>

import {onMounted, ref} from "vue";
import {useRouter} from "vue-router";
import FooterBar from "@/components/FooterBar.vue";
import {httpGet} from "@/utils/http";
import {ElMessage} from "element-plus";
import {isMobile} from "@/utils/libs";
import {checkSession} from "@/action/session";

const router = useRouter()

if (isMobile()) {
  router.push("/mobile")
}

const title = ref("")
const logo = ref("")
const slogan = ref("")
const license = ref({})
const winHeight = window.innerHeight - 150
const bgImgUrl = ref('')
const isLogin = ref(false)
const docsURL = ref(process.env.VUE_APP_DOCS_URL)
const gitURL = ref(process.env.VUE_APP_GIT_URL)

onMounted(() => {
  httpGet("/api/config/get?key=system").then(res => {
    title.value = res.data.title
    logo.value = res.data.logo
    if (res.data.index_bg_url) {
      bgImgUrl.value = res.data.index_bg_url
    } else {
      bgImgUrl.value = "/images/index-bg.jpg"
    }
    if (res.data.slogan) {
      slogan.value = res.data.slogan
    }
  }).catch(e => {
    ElMessage.error("获取系统配置失败：" + e.message)
  })

  httpGet("/api/config/license").then(res => {
    license.value = res.data
  }).catch(e => {
    ElMessage.error("获取 License 配置：" + e.message)
  })

  checkSession().then(() => {
    isLogin.value = true
  }).catch(()=>{})
})
</script>

<style lang="stylus" scoped>
@import '@/assets/iconfont/iconfont.css'
.index-page {
  margin: 0
  overflow hidden
  color #ffffff
  display flex
  justify-content center
  align-items baseline
  padding-top 150px

  .index-bg {
    position absolute
    top 0
    left 0
    width 100vw
    height 100vh
    filter: blur(8px);
    background-size: cover;
    background-position: center;
  }

  .menu-box {
    position absolute
    top 0
    width 100%
    display flex

    .el-menu {
      padding 0 30px
      width 100%
      display flex
      justify-content space-between
      background none
      border none

      .menu-item {
        display flex
        padding 20px 0

        color #ffffff

        .title {
          font-size 24px
          padding 10px 10px 0 10px
        }

        .el-image {
          height 50px
        }

        .el-button {
          margin-left 10px

          span {
            margin-left 5px
          }
        }
      }
    }
  }

  .content {
    text-align: center;
    position relative

    h1 {
      font-size: 5rem;
      margin-bottom: 1rem;
    }

    p {
      font-size: 1.5rem;
      margin-bottom: 2rem;
    }

    .el-button {
      padding: 25px 20px;
      font-size: 1.3rem;
      transition: all 0.3s ease;

      .iconfont {
        font-size 1.6rem
        margin-right 10px
      }
    }

    #animation-container {
      display flex
      justify-content center
      width 100%
      height: 300px;
      position: absolute;
      top: 350px

    }
  }

  .footer {
    .el-link__inner {
      color #ffffff
    }
  }

}
</style>
