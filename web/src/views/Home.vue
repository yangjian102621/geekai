<template>
  <div class="home">
    <div class="navigator">
      <div class="logo">
        <el-link href="/">
          <el-image :src="logo"/>
        </el-link>
        <div class="divider"></div>
      </div>
      <ul class="nav-items">
        <li v-for="item in navs" :key="item.path">
          <!--          <el-tooltip effect="light" :content="item.title" placement="right">-->
          <!--            -->
          <!--          </el-tooltip>-->
          <a @click="changeNav(item)" :class="item.path === curPath ? 'active' : ''">
            <el-image :src="item.icon_path" :width="20" v-if="item.icon_path"/>
            <i :class="'iconfont icon-' + item.icon" v-else></i>
          </a>
          <div :class="item.path === curPath ? 'title active' : 'title'">{{ item.title }}</div>
        </li>
      </ul>
    </div>
    <div class="content">
      <router-view v-slot="{ Component }">
        <transition name="move" mode="out-in">
          <component :is="Component"></component>
        </transition>
      </router-view>
    </div>
  </div>
</template>

<script setup>

import {useRouter} from "vue-router";
import {checkSession} from "@/action/session";
import {isMobile} from "@/utils/libs";
import {ref} from "vue";

const router = useRouter();
const logo = '/images/logo.png';
const navs = ref([
  {path: "/chat", icon_path: "/images/chat.png", title: "对话聊天"},
  {path: "/mj", icon_path: "/images/mj.png", title: "MJ 绘画"},
  {path: "/sd", icon_path: "/images/sd.png", title: "SD 绘画"},
  {path: "/apps", icon: "menu", title: "应用中心"},
  {path: "/images-wall", icon: "image-list", title: "作品展示"},
  {path: "/knowledge", icon: "book", title: "知识库"},
  {path: "/member", icon: "vip-user", title: "会员计划"},
  {path: "/invite", icon: "share", title: "推广计划"},
])
const curPath = ref(router.currentRoute.value.path)

const changeNav = (item) => {
  curPath.value = item.path
  router.push(item.path)
}
</script>

<style lang="stylus" scoped>
@import '@/assets/iconfont/iconfont.css';
.home {
  display: flex;
  background-color: #25272D;
  height 100vh
  width 100%

  .navigator {
    display flex
    flex-flow column
    width 70px
    padding 10px 6px
    border-right: 1px solid #3c3c3c

    .logo {
      display flex
      flex-flow column
      align-items center


      .divider {
        border-bottom 1px solid #4A4A4A
        width 80%
        height 10px
      }
    }

    .nav-items {
      margin-top: 20px;
      padding-left: 10px;
      padding-right: 10px;

      li {
        margin-bottom 15px

        a {
          color #DADBDC
          background-color #40444A
          border-radius 10px
          width 48px
          height 48px
          display flex
          justify-content center
          align-items center
          cursor pointer

          .el-image {
            border-radius 10px
          }

          .iconfont {
            font-size 20px
          }
        }

        a:hover, a.active {
          color #47fff1
        }

        .title {
          font-size: 12px
          padding-top: 5px
          color: #e5e7eb;
          text-align: center;
        }

        .active {
          color #47fff1
        }
      }
    }
  }

  .content {
    width: 100%;
    height: 100vh;
    box-sizing: border-box;
  }

}
</style>
