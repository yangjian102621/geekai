<template>
  <div class="header admin-header">
    <div class="logo">
      <el-image :src="logo"/>
      <span class="text">{{ title }}</span>
    </div>
    <!-- 折叠按钮 -->
    <div class="collapse-btn" @click="collapseChange">
      <el-icon v-if="sidebar.collapse">
        <Expand/>
      </el-icon>
      <el-icon v-else>
        <Fold/>
      </el-icon>
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
        <!-- 用户头像 -->
        <el-avatar class="user-avatar" :size="30" :src="avatar"/>
        <!-- 用户名下拉菜单 -->
        <el-dropdown class="user-name" :hide-on-click="true" trigger="click">
					<span class="el-dropdown-link">
						{{ username }}
						<el-icon class="el-icon--right">
							<arrow-down/>
						</el-icon>
					</span>
          <template #dropdown>
            <el-dropdown-menu>

              <a href="https://github.com/yangjian102621/chatgpt-plus" target="_blank">
                <el-dropdown-item>
                  <i class="iconfont icon-github"></i>
                  <span>ChatGPT-Plus-V3</span>
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
        custom-class="donate-dialog"
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
import {useSidebarStore} from '@/store/sidebar';
import {useRouter} from 'vue-router';
import {ArrowDown, Expand, Fold} from "@element-plus/icons-vue";
import {httpGet} from "@/utils/http";
import {ElMessage} from "element-plus";

const message = ref(5);
const username = ref('极客学长')
const avatar = ref('/images/user-info.jpg')
const donateImg = ref('/images/wechat-pay.png')
const showDialog = ref(false)
const sidebar = useSidebarStore();
const title = ref('Chat-Plus 控制台')
const logo = ref('/images/logo.png')

// 加载系统配置
httpGet('/api/admin/config/get?key=system').then(res => {
  title.value = res.data['admin_title'];
}).catch(e => {
  ElMessage.error("加载系统配置失败: " + e.message)
})

// 侧边栏折叠
const collapseChange = () => {
  sidebar.handleCollapse();
};

onMounted(() => {
  if (document.body.clientWidth < 1024) {
    collapseChange();
  }
});

const router = useRouter();
const logout = function () {
  httpGet("/api/admin/logout").then(() => {
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
  width: 100%;
  height: 70px;
  font-size: 22px;
  color: #fff;

  .collapse-btn {
    display: flex;
    justify-content: center;
    align-items: center;
    height: 100%;
    float: left;
    padding: 0 21px;
    cursor: pointer;
  }

  .logo {
    float: left;
    padding-left 10px;
    display flex

    .text {
      line-height: 66px;
      margin-left 10px;
    }
  }

  .header-right {
    float: right;
    padding-right: 20px;

    .header-user-con {
      display: flex;
      height: 70px;
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
          color: #fff;
        }

        .icon-bell {
          font-size: 24px;
        }
      }

      .user-name {
        margin-left: 10px;
      }

      .user-avatar {
        margin-left: 20px;
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
  .logo {
    .el-image {
      padding-top 10px

      .el-image__inner {
        height 40px
      }
    }
  }
}

</style>
