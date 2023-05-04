<template>
  <div class="common-layout">
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
                @click="hideSidebar()"
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
              </el-tabs>
            </div>
          </div>
        </el-main>

        <el-dialog
            v-model="showDialog"
            title="管理后台登录"
            width="30%"
            :destroy-on-close="true"
        >
          <el-form :model="user" label-width="80px">
            <el-form-item label="用户名：">
              <el-input
                  v-model="user.username"
                  autocomplete="off"
                  placeholder="请输入用户名"
              />
            </el-form-item>

            <el-form-item label="密码：">
              <el-input
                  v-model="user.password"
                  autocomplete="off"
                  type="password"
                  placeholder="请输入密码"
              />
            </el-form-item>
          </el-form>

          <template #footer>
            <span class="dialog-footer">
              <el-button @click="showDialog = false">取消</el-button>
              <el-button type="primary" @click="saveHost">提交</el-button>
            </span>
          </template>
        </el-dialog>
      </el-container>
    </el-container>
  </div>
</template>

<script>
import {defineComponent} from 'vue'
import {Fold, Menu} from "@element-plus/icons-vue"
import XWelcome from "@/views/admin/Welcome.vue";
import SysConfig from "@/views/admin/SysConfig.vue";
import {arrayContains, removeArrayItem} from "@/utils/libs";
import UserList from "@/views/admin/UserList.vue";
import RoleList from "@/views/admin/RoleList.vue";


export default defineComponent({
  name: "XAdmin",
  components: {RoleList, UserList, SysConfig, XWelcome, Fold, Menu},
  data() {
    return {
      title: "Chat-Plus 控制台",
      logo: 'images/logo.png',
      user: {},
      navs: [
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
        }
      ],

      curNav: null,
      curTab: 'welcome',
      tabs: [],

      showDialog: false,

      // window height
      winHeight: window.innerHeight,
      showSidebar: true
    }
  },

  computed: {
    sideWidth: function () {
      return this.showSidebar ? 250 : 30
    },

    foldIconRight: function () {
      return this.showSidebar ? 3 : 0
    },

    nodeListPaddingLeft: function () {
      return this.showSidebar ? 20 : 5
    }
  },

  mounted: function () {

    // bind window resize event
    window.addEventListener("resize", function () {
      this.winHeight = window.innerHeight
    })
  },

  methods: {
    arrayContains(array, value, compare) {
      return arrayContains(array, value, compare);
    },

    hideSidebar: function () {
      this.showSidebar = !this.showSidebar
    },

    // 添加 tab 窗口
    addTab: function (nav) {
      if (this.curNav) {
        this.curNav.active = false
      }
      this.curNav = nav;
      this.curNav.active = true;
      this.curTab = nav.tab;
      if (!arrayContains(this.tabs, nav.tab)) {
        this.tabs.push(nav.tab);
      }
    },

    changeTab: function (name) {
      for (let i = 0; i < this.navs.length; i++) {
        if (this.navs[i].tab === name) {
          this.curNav.active = false
          this.curNav = this.navs[i];
          this.curNav.active = true;
          break;
        }
      }
    },

    // 删除 tab 窗口
    removeTab: function (name) {
      this.tabs = removeArrayItem(this.tabs, name);
      if (this.tabs.length === 0) {
        this.curTab = 'welcome';
        return;
      }

      for (let i = 0; i < this.navs.length; i++) {
        if (this.navs[i].tab === this.tabs[this.tabs.length - 1]) {
          this.addTab(this.navs[i]);
        }
      }

    }
  },

})
</script>

<style lang="stylus">
$sideBgColor = #252526;
$borderColor = #4676d0;
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
      color: #aaa;
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
}

.el-main {
  --el-main-padding: 0;
  margin: 0;
  background-image url("~@/assets/img/bg_01.jpeg")

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
    }

  }
}
</style>
