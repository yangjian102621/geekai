<template>
  <div class="header admin-header">
    <!-- 折叠按钮 -->
    <div class="collapse-btn" @click="collapseChange">
      <el-icon v-if="sidebar.collapse">
        <Expand/>
      </el-icon>
      <el-icon v-else>
        <Fold/>
      </el-icon>
    </div>

    <div class="breadcrumb">
      <el-breadcrumb :separator-icon="ArrowRight">
        <el-breadcrumb-item v-for="item in breadcrumb">{{ item.title }}</el-breadcrumb-item>
      </el-breadcrumb>
    </div>
    <div class="header-right">
      <div class="header-user-con">
        <!-- 消息中心 -->
        <div class="btn-bell">
          <el-tooltip
              effect="dark"
              :content="message ? `有${message}条未读消息` : `消息中心`"
              placement="bottom"
          >
            <i class="iconfont icon-bell"></i>
          </el-tooltip>
          <span class="btn-bell-badge" v-if="message"></span>
        </div>
        <!-- 用户名下拉菜单 -->
        <el-dropdown class="user-name" :hide-on-click="true" trigger="click">
					<span class="el-dropdown-link">
						<el-avatar class="user-avatar" :size="30" :src="avatar"/>
						<el-icon class="el-icon--right">
							<arrow-down/>
						</el-icon>
					</span>
          <template #dropdown>
            <el-dropdown-menu>

              <a href="https://github.com/yangjian102621/chatgpt-plus" target="_blank">
                <el-dropdown-item>
                  <i class="iconfont icon-github"></i>
                  <span>{{ sysTitle }}</span>
                </el-dropdown-item>
              </a>
              <el-dropdown-item @click="showDialog = true">
                <i class="iconfont icon-reward"></i>
                <span>打赏作者</span>
              </el-dropdown-item>
              <el-dropdown-item divided @click="logout">
                <i class="iconfont icon-logout"></i>
                <span>退出登录</span>
              </el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
    </div>

    <el-dialog
        v-model="showDialog"
        :show-close="true"
        class="donate-dialog"
        width="400px"
        title="请作者喝杯咖啡"
    >
      <el-alert type="info" :closable="false">
        <p>如果你觉得这个项目对你有帮助，并且情况允许的话，可以请作者喝杯咖啡，非常感谢你的支持～</p>
      </el-alert>
      <p>
        <el-image :src="donateImg"/>
      </p>
    </el-dialog>
  </div>
</template>
<script setup>
import {onMounted, ref} from 'vue';
import {getMenuItems, useSidebarStore} from '@/store/sidebar';
import {useRouter} from 'vue-router';
import {ArrowDown, ArrowRight, Expand, Fold} from "@element-plus/icons-vue";
import {httpGet} from "@/utils/http";
import {ElMessage} from "element-plus";
import {removeAdminToken} from "@/store/session";

const message = ref(5);
const sysTitle = ref(process.env.VUE_APP_TITLE)
const avatar = ref('/images/user-info.jpg')
const donateImg = ref('/images/wechat-pay.png')
const showDialog = ref(false)
const sidebar = useSidebarStore();
const router = useRouter();
const breadcrumb = ref([])


router.afterEach((to, from) => {
  initBreadCrumb(to.path)
});

onMounted(() => {
  initBreadCrumb(router.currentRoute.value.path)
})

// 初始化面包屑导航
const initBreadCrumb = (path) => {
  breadcrumb.value = [{title: "首页"}]
  const items = getMenuItems()
  if (items) {
    let bk = false
    for (let i = 0; i < items.length; i++) {
      if (items[i].index === path) {
        breadcrumb.value.push({
          title: items[i].title,
          path: items[i].index
        })
        break
      }
      if (bk) {
        break
      }

      if (items[i]['subs']) {
        const subs = items[i]['subs']
        for (let j = 0; j < subs.length; j++) {
          if (subs[j].index === path) {
            breadcrumb.value.push({
              title: items[i].title,
              path: items[i].index
            })
            breadcrumb.value.push({
              title: subs[j].title,
              path: subs[j].index
            })
            bk = true
            break
          }
        }
      }
    }
  }
}

// 侧边栏折叠
const collapseChange = () => {
  sidebar.handleCollapse();
};

onMounted(() => {
  if (document.body.clientWidth < 1024) {
    collapseChange();
  }
});

const logout = function () {
  httpGet("/api/admin/logout").then(() => {
    removeAdminToken()
    router.replace('/admin/login')
  }).catch((e) => {
    ElMessage.error("注销失败: " + e.message);
  })
}
</script>
<style scoped lang="stylus">
.header {
  position: relative;
  box-sizing: border-box;
  overflow hidden
  height: 50px;
  font-size: 22px;
  color: #303133;
  background-color #ffffff
  border-bottom 1px solid #eaecef

  .collapse-btn {
    display: flex;
    justify-content: center;
    align-items: center;
    height: 100%;
    float: left;
    padding: 0 10px;
    cursor: pointer;

    &:hover {
      background-color #eaecef
    }
  }

  .breadcrumb {
    float left
    display flex
    align-items center
    height 50px
  }

  .header-right {
    float: right;
    padding-right: 20px;

    .header-user-con {
      display: flex;
      height: 50px;
      align-items: center;

      .btn-bell {
        position: relative;
        width: 30px;
        height: 30px;
        text-align: center;
        border-radius: 15px;
        cursor: pointer;
        display: flex;
        align-items: center;

        .btn-bell-badge {
          position: absolute;
          right: 4px;
          top: 0;
          width: 8px;
          height: 8px;
          border-radius: 4px;
          background: #f56c6c;
          color: #303133;
        }

        .icon-bell {
          font-size: 24px;
        }
      }

      .user-name {
        margin-left: 10px;

        .el-icon {
          color: #303133;
        }
      }

      .user-avatar {

      }
    }
  }
}

.el-dropdown-link {
  color: #fff;
  cursor: pointer;
  display: flex;
  align-items: center;
}

.el-dropdown-menu__item {
  text-align: center;

  .icon-reward {
    font-size 18px;
    font-weight bold;
    color #F56C6C
  }
}
</style>

<style lang="stylus">
.donate-dialog {
  .el-dialog__body {
    text-align center;

    .el-alert__description {
      text-align left
      font-size 14px;
      line-height 1.5
    }
  }
}

.admin-header {

}

</style>
