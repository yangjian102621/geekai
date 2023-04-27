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
                @click="changeMenu(nav)"
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
              id="main-container"
              :style="{ height: winHeight + 'px' }"
          >
            <XWelcome v-if="curPage==='welcome'"/>
            <TestPage v-if="curPage==='config'"/>
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
import XWelcome from "@/views/Welcome.vue";
import TestPage from "@/views/Test.vue";


export default defineComponent({
  name: "XAdmin",
  components: {TestPage, XWelcome, Fold, Menu},
  data() {
    return {
      title: "Chat-Plus 控制台",
      logo: 'images/logo.png',
      user: {},
      navs: [
        {
          id: 1,
          title: '系统配置',
          page: 'config',
          active: false,
        },
        {
          id: 2,
          title: '用户管理',
          page: 'user',
          active: false,
        },
        {
          id: 3,
          title: '角色管理',
          page: 'role',
          active: false,
        }
      ],
      curNav: null,
      curPage: 'welcome',

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
    hideSidebar: function () {
      this.showSidebar = !this.showSidebar
    },
    changeMenu: function (nav) {
      if (this.curNav) {
        this.curNav.active = false
      }
      nav.active = true;
      this.curNav = nav;
      this.curPage = nav.page;
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
  overflow: hidden;
  margin: 0;
  background-image url("~@/assets/img/bg_01.jpeg")

  .main-container {
    display: flex;
    flex-flow: column;

    .console-wrapper {
      position: absolute;
      width: 100%;
    }
  }
}
</style>
