<template>
  <div class="admin-body common-layout">
    <el-container>
      <el-container>
        <el-aside
            :style="{
            height: winHeight + 'px',
            width: sideWidth + 'px',
          }"
        >
          <div :class="showSidebar?'title':'title hide'">
            <el-image :src="logo" class="logo"/>
            <span class="text">{{ title }}</span>
            <span
                class="fold"
                @click="showSidebar = !showSidebar"
                :style="{ right: foldIconRight + 'px' }"
            >
              <el-icon>
                <Fold/>
              </el-icon>
            </span>
          </div>
          <ul class="nav-list">
            <li
                v-for="nav in navs"
                :key="nav.id"
                :style="{ paddingLeft: nodeListPaddingLeft + 'px' }"
                :class="nav.active?'active':''"
                @click="addTab(nav)"
            >
              <el-tooltip
                  class="box-item"
                  effect="light"
                  :content="nav.title"
                  placement="right"
              >
                <el-icon>
                  <Menu/>
                </el-icon>
              </el-tooltip>

              <span v-if="showSidebar">{{ nav.title }}</span>
            </li>
          </ul>

          <el-row class="nav-footer">
            <div class="source">
              <i class="iconfont icon-github"></i>
              <el-link href="https://github.com/yangjian102621/chatgpt-plus" target="_blank">ChatGPT-Plus-V3</el-link>
            </div>

            <div class="logout" @click="logout">
              <i class="iconfont icon-logout"></i>
              <span>退出登录</span>
            </div>
          </el-row>
        </el-aside>

        <el-main>
          <div
              class="main-container"
              :style="{ height: winHeight + 'px' }"
          >
            <x-welcome v-if="curTab==='welcome'"/>

            <div v-else>
              <el-tabs
                  v-model="curTab"
                  class="content-tabs"
                  type="card"
                  closable
                  @tab-remove="removeTab"
                  @tab-change="changeTab"
              >
                <el-tab-pane label="系统配置" name="config" v-if="arrayContains(tabs, 'config')">
                  <sys-config v-if="curTab==='config'"/>
                </el-tab-pane>

                <el-tab-pane label="用户管理" name="user" v-if="arrayContains(tabs, 'user')">
                  <user-list v-if="curTab==='user'"/>
                </el-tab-pane>

                <el-tab-pane label="角色管理" name="role" v-if="arrayContains(tabs, 'role')">
                  <role-list v-if="curTab==='role'"/>
                </el-tab-pane>
                <el-tab-pane label="API KEY" name="apikey" v-if="arrayContains(tabs, 'apikey')">
                  <api-key v-if="curTab==='apikey'"/>
                </el-tab-pane>
              </el-tabs>
            </div>
          </div>
        </el-main>
      </el-container>
    </el-container>
  </div>
</template>

<script setup>
import {computed, onMounted, ref} from 'vue'
import {Fold, Menu} from "@element-plus/icons-vue"
import XWelcome from "@/views/admin/Welcome.vue";
import SysConfig from "@/views/admin/SysConfig.vue";
import {arrayContains, removeArrayItem} from "@/utils/libs";
import UserList from "@/views/admin/UserList.vue";
import RoleList from "@/views/admin/RoleList.vue";
import {httpGet} from "@/utils/http";
import {ElMessage} from "element-plus";
import {useRouter} from "vue-router";
import ApiKey from "@/views/admin/ApiKey.vue";

const title = ref('Chat-Plus 控制台')
const logo = ref('images/logo.png')
const user = ref({})
const navs = ref([
  {
    id: 1,
    title: '系统配置',
    tab: 'config',
    active: false,
  },
  {
    id: 2,
    title: '用户管理',
    tab: 'user',
    active: false,
  },
  {
    id: 3,
    title: '角色管理',
    tab: 'role',
    active: false,
  },
  {
    id: 4,
    title: 'API KEY',
    tab: 'apikey',
    active: false,
  }
])
const tabs = ref([])
const curNav = ref(null)
const curTab = ref('welcome')
const winHeight = ref(window.innerHeight)
const showSidebar = ref(true)

const sideWidth = computed(() => {
  return showSidebar.value ? 250 : 30
})
const foldIconRight = computed(() => {
  return showSidebar.value ? 3 : 0
})
const nodeListPaddingLeft = computed(() => {
  return showSidebar.value ? 20 : 5
})
const router = useRouter()

onMounted(() => {
  window.addEventListener("resize", function () {
    winHeight.value = window.innerHeight
  })

  // 获取会话信息
  httpGet("/api/admin/session").catch(() => {
    router.push('/admin/login')
  })

  // 加载系统配置
  httpGet('/api/admin/config/get?key=system').then(res => {
    title.value = res.data['admin_title'];
  }).catch(e => {
    ElMessage.error("加载系统配置失败: " + e.message)
  })
})

const logout = function () {
  httpGet("/api/admin/logout").then(() => {
    router.push('/admin/login')
  }).catch((e) => {
    ElMessage.error("注销失败: " + e.message);
  })
}

// 添加 tab 窗口
const addTab = function (nav) {
  if (curNav.value) {
    curNav.value.active = false
  }
  nav.active = true
  curNav.value = nav;
  curTab.value = nav.tab;
  if (!arrayContains(tabs.value, nav.tab)) {
    this.tabs.push(nav.tab);
  }
}

// 切换 tab 窗口
const changeTab = function (name) {
  for (let i = 0; i < navs.value.length; i++) {
    let _nav = navs.value[i]
    if (_nav.tab === name) {
      curNav.value.active = false // 取消上一个 active 窗口的激活状态
      _nav.active = true
      curNav.value = _nav;
      break;
    }
  }
}

// 删除 tab 窗口
const removeTab = function (name) {
  tabs.value = removeArrayItem(tabs.value, name);
  if (tabs.value.length === 0) {
    curTab.value = 'welcome';
    return;
  }

  for (let i = 0; i < navs.value.length; i++) {
    if (navs.value[i].tab === tabs.value[tabs.value.length - 1]) {
      addTab(navs.value[i]);
    }
  }

}
</script>

<style scoped lang="stylus">
$sideBgColor = #252526;
$borderColor = #4676d0;
.admin-body {
  .el-aside {
    background-color: $sideBgColor;

    .title {
      text-align: center;
      line-height: 60px;
      color: #fff;
      font-size: 20px;
      border-bottom: 2px solid #333841;
      display flex
      flex-direction row

      .logo {
        background-color #ffffff
        border-radius 50%;
        width 32px;
        height 32px;
        margin: 12px 5px 0 5px;
      }

      .fold {
        cursor: pointer;
        position: relative;
        top: 2px;
        margin-left 10px;
      }
    }

    .title.hide {
      .text {
        display none
      }

      .logo {
        display none
      }

      .fold {
        margin-left 5px;
      }
    }

    .nav-list {
      list-style: none;
      position: relative;
      margin: 0;
      padding-left: 0;
      text-align: left;

      li {
        line-height: 40px;
        color: #ffffff;
        font-size: 14px;
        cursor: pointer;
        padding: 0 10px 0 10px;
        border-bottom: 1px dashed #333841;

        i {
          margin-right: 6px;
          position: relative;
          top: 1px;
        }

        .delete {
          float: right;
        }
      }

      li.active {
        background-color: #363535
      }
    }

    .nav-footer {
      flex-direction column

      div {
        padding 10px 20px;
        font-size 14px;
        color #aaaaaa

        .el-link {
          color #aaaaaa
        }

        .iconfont {
          margin-right 5px
          position relative
          top 1px
        }
      }

      .logout {
        cursor pointer
      }
    }
  }

  .el-main {
    --el-main-padding: 0;
    margin: 0;

    .main-container {
      display: flex;
      flex-flow: column;

      .content-tabs {
        background: #ffffff;
        padding 10px 20px;

        .el-tabs__item {
          height 35px
          line-height 35px
        }

        .el-tabs__content {
          padding 10px 20px 20px 20px;
        }
      }

    }
  }
}


</style>
