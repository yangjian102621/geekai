<template>
  <div class="home">
    <div class="navigator">
      <div class="logo">
        <el-image :src="logo"/>
        <div class="divider"></div>
      </div>
      <ul class="nav-items">
        <li v-for="item in navs" :key="item.path">
          <!--          <el-tooltip effect="light" :content="item.title" placement="right">-->
          <!--            -->
          <!--          </el-tooltip>-->
          <a @click="changeNav(item)" :class="item.path === curPath ? 'active' : ''">
            <el-image :src="item.icon" :width="20"/>
          </a>
          <div :class="item.url === curPath ? 'title active' : 'title'">{{ item.name }}</div>
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
import {onMounted, ref} from "vue";
import {httpGet} from "@/utils/http";
import {ElMessage} from "element-plus";

const router = useRouter();
const logo = ref('/images/logo.png');
const navs = ref([])
const curPath = ref(router.currentRoute.value.path)

const changeNav = (item) => {
  curPath.value = item.url
  router.push(item.url)
}

onMounted(() => {
  httpGet("/api/config/get?key=system").then(res => {
    logo.value = res.data['logo']
  }).catch(e => {
    ElMessage.error("获取系统配置失败：" + e.message)
  })
  // 获取菜单
  httpGet("/api/menu/list").then(res => {
    navs.value = res.data
  }).catch(e => {
    ElMessage.error("获取系统菜单失败：" + e.message)
  })
})
</script>

<style lang="stylus" scoped>
@import '@/assets/iconfont/iconfont.css';
.home {
  display: flex;
  height 100vh
  width 100%

  .navigator {
    display flex
    flex-flow column
    width 70px
    padding 10px 6px
    border-right: 1px solid #3c3c3c
    background-color: #25272D

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
