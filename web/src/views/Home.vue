<template>
  <div class="home">
    <div class="navigator">
      <div class="logo">
        <el-image :src="logo" @click="router.push('/')"/>
        <div class="divider"></div>
      </div>
      <ul class="nav-items">
        <li v-for="item in mainNavs" :key="item.url">
          <a @click="changeNav(item)" :class="item.url === curPath ? 'active' : ''">
            <el-image :src="item.icon" style="width: 30px;height: 30px"/>
          </a>
          <div :class="item.url === curPath ? 'title active' : 'title'">{{ item.name }}</div>
        </li>

        <el-popover
            v-if="moreNavs.length > 0"
            placement="right-end"
            trigger="hover"
        >
          <template #reference>
            <li>
              <a class="active">
                <el-image src="/images/menu/more.png" style="width: 30px;height: 30px"/>
              </a>
            </li>
          </template>
          <template #default>
            <ul class="more-menus">
              <li v-for="item in moreNavs" :key="item.url" :class="item.url === curPath ? 'active' : ''">
                <a @click="changeNav(item)">
                  <el-image :src="item.icon" style="width: 20px;height: 20px"/>
                  <span :class="item.url === curPath ? 'title active' : 'title'">{{ item.name }}</span>
                </a>
              </li>
            </ul>
          </template>
        </el-popover>
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
const mainNavs = ref([])
const moreNavs = ref([])
const curPath = ref(router.currentRoute.value.path)

if (curPath.value === "/external") {
  curPath.value = router.currentRoute.value.query.url
}
const changeNav = (item) => {
  curPath.value = item.url
  if (item.url.indexOf("http") !== -1) { // 外部链接
    router.push({name: 'ExternalLink', query: {url: item.url}})
  } else {
    router.push(item.url)
  }
}

onMounted(() => {
  httpGet("/api/config/get?key=system").then(res => {
    logo.value = res.data['logo']
  }).catch(e => {
    ElMessage.error("获取系统配置失败：" + e.message)
  })
  // 获取菜单
  httpGet("/api/menu/list").then(res => {
    mainNavs.value = res.data
    // 根据窗口的高度计算应该显示多少菜单
    const rows = Math.floor((window.innerHeight - 90) / 90)
    if (res.data.length > rows) {
      mainNavs.value = res.data.slice(0, rows)
      moreNavs.value = res.data.slice(rows)
    }
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
    width 60px
    padding 10px 6px
    border-right: 1px solid #3c3c3c
    background-color: #25272D

    .logo {
      display flex
      flex-flow column
      align-items center
      cursor pointer

      .el-image {
        width 50px
        height 50px
      }

      .divider {
        border-bottom 1px solid #4A4A4A
        width 80%
        height 10px
      }
    }

    .nav-items {
      margin-top: 10px;
      padding 0 5px

      li {
        margin-bottom 15px

        a {
          color #DADBDC
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
          background-color #0F7A71
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
    width: 100%
    height: 100vh
    box-sizing: border-box
    background-color #282c34
  }

}

.el-popper {
  .more-menus {
    li {
      padding 10px 15px
      cursor pointer
      border-radius 5px
      margin 5px 0

      .el-image {
        position: relative
        top 5px
        right 5px
      }

      &:hover {
        background-color #f1f1f1
      }
    }

    li.active {
      background-color #f1f1f1
    }
  }
}
</style>
