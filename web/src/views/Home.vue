<template>
  <div class="layout">
    <div class="tab-box">
      <div class="flex-center-col big-top-title xxx">
        <div class="flex-center-col" @click="isCollapse = !isCollapse">
          <el-icon v-if="isCollapse" class="openicon">
            <svg
              t="1733138242826"
              class="icon"
              viewBox="0 0 1024 1024"
              version="1.1"
              xmlns="http://www.w3.org/2000/svg"
              p-id="1853"
              id="mx_n_1733138242827"
              width="200"
              height="200"
            >
              <path
                d="M715 267c135.31 0 245 109.69 245 245S850.31 757 715 757H309C173.69 757 64 647.31 64 512s109.69-245 245-245h406zM309 367c-80.081 0-145 64.919-145 145s64.919 145 145 145 145-64.919 145-145-64.919-145-145-145z"
                fill="#754ff6"
                p-id="1854"
              ></path>
            </svg>
          </el-icon>
        </div>
        <div class="flex" :class="{ 'top-collapse': !isCollapse }">
          <div class="top-avatar flex">
            <span class="title" v-if="!isCollapse">GeekAI</span>
            <img
              v-if="loginUser.id"
              :src="!!loginUser.avatar ? loginUser.avatar : avatarImg"
              alt=""
              :class="{ marr: !isCollapse }"
            />
          </div>
          <div class="menuIcon xxx" @click="isCollapse = !isCollapse">
            <el-icon v-if="!isCollapse" class="openicon">
              <svg
                t="1733138405307"
                class="icon"
                viewBox="0 0 1024 1024"
                version="1.1"
                xmlns="http://www.w3.org/2000/svg"
                p-id="2064"
                width="200"
                height="200"
              >
                <path
                  d="M715 267c135.31 0 245 109.69 245 245S850.31 757 715 757H309C173.69 757 64 647.31 64 512s109.69-245 245-245h406z m0 100c-80.081 0-145 64.919-145 145s64.919 145 145 145 145-64.919 145-145-64.919-145-145-145z"
                  fill="#754ff6"
                  p-id="2065"
                ></path>
              </svg>
            </el-icon>
          </div>
        </div>
      </div>

      <div
        class="menu-list"
        :style="{ width: isCollapse ? '65px' : '170px' }"
        :class="{ 'menu-list-collapse': !isCollapse }"
      >
        <ul>
          <li
            class="menu-list-item flex-center-col"
            v-for="item in mainNavs"
            :key="item.url"
            @click="changeNav(item)"
            :class="item.url === curPath ? 'active' : ''"
          >
            <el-image :src="item.icon" class="el-icon" />
            <div
              class="menu-title"
              :class="{ 'menu-title-collapse': !isCollapse }"
            >
              {{ item.name }}
            </div>
          </li>

          <!-- <li
            class="menu-list-item flex-center-col"
            v-for="item in 5"
            :key="item"
          >
            <el-icon><Location /></el-icon>
            <div>首页</div>
          </li> -->

          <!-- 更多 -->
          <div class="bot" :style="{ width: isCollapse ? '65px' : '170px' }">
            <div class="bot-line"></div>

            <el-popover
              v-if="moreNavs.length > 0"
              placement="right-end"
              trigger="hover"
            >
              <template #reference>
                <li class="menu-list-item flex-center-col">
                  <el-icon><CirclePlus /></el-icon>
                  <div class="menu-title">更多</div>
                </li>
              </template>
              <template #default>
                <ul class="more-menus">
                  <li
                    v-for="(item, index) in moreNavs"
                    :key="item.url"
                    :class="{
                      active: item.url === curPath,
                      moreTitle: index !== 3 && index !== 4,
                      twoTittle: index === 3 || index === 4
                    }"
                  >
                    <a @click="changeNav(item)">
                      <el-image
                        :src="item.icon"
                        style="width: 20px; height: 20px"
                      />
                      <span
                        :class="item.url === curPath ? 'title active' : 'title'"
                        >{{ item.name }}</span
                      >
                    </a>
                  </li>
                </ul>
              </template>
            </el-popover>
            <el-popover
              placement="right-end"
              trigger="hover"
              v-if="loginUser.id"
            >
              <template #reference>
                <li class="menu-list-item flex-center-col">
                  <el-icon><Setting /></el-icon>
                  <div v-if="!isCollapse">设置</div>
                </li>
              </template>
              <template #default>
                <ul class="more-menus setting-menus">
                  <li>
                    <div @click="showConfigDialog = true" class="flex">
                      <el-icon>
                        <UserFilled />
                      </el-icon>
                      <span class="username title">{{
                        loginUser.nickname
                      }}</span>
                    </div>
                  </li>
                  <li>
                    <a @click="logout" class="flex">
                      <i class="iconfont icon-logout"></i>
                      <span class="title">退出登录</span>
                    </a>
                  </li>
                </ul>
              </template>
            </el-popover>
            <li class="menu-bot-item">
              <a
                :href="gitURL"
                class="link-button"
                target="_blank"
                v-if="!license.de_copy && !isCollapse"
              >
                <i class="iconfont icon-github"></i>
              </a>

              <a @click="router.push('/')" class="link-button">
                <i class="iconfont icon-house"></i>
              </a>

              <ThemeChange />
              <!-- <div v-if="!isCollapse">会员</div> -->
            </li>
          </div>
        </ul>
      </div>
    </div>
    <!-- :style="{ 'padding-left': isCollapse ? '65px' : '170px' }" -->
    <div class="right-main">
      <div
        v-if="loginUser.id === undefined || !loginUser.id"
        class="loginMask"
        :style="{ left: isCollapse ? '65px' : '170px' }"
        @click="showNoticeLogin = true"
      ></div>
      <div class="topheader" v-if="loginUser.id === undefined || !loginUser.id">
        <el-button
          @click="router.push('/login')"
          class="btn-go animate__animated animate__pulse animate__infinite"
          round
          >登录</el-button
        >
      </div>
      <!-- <div class="content custom-scroll"> -->
      <div class="content custom-scroll">
        <router-view :key="routerViewKey" v-slot="{ Component }">
          <transition name="move" mode="out-in">
            <component :is="Component"></component>
          </transition>
        </router-view>
      </div>
      <!-- </div> -->
    </div>
    <config-dialog
      v-if="loginUser.id"
      :show="showConfigDialog"
      @hide="showConfigDialog = false"
    />
  </div>
  <el-dialog v-model="showNoticeLogin">
    <el-result icon="warning" title="未登录" sub-title="登录后解锁功能">
      <template #extra>
        <el-button type="primary" @click="router.push('/login')"
          >登录</el-button
        >
      </template>
    </el-result>
  </el-dialog>
</template>

<script setup>
import { CirclePlus, Setting } from "@element-plus/icons-vue";
import ThemeChange from "@/components/ThemeChange.vue";
import avatarImg from "@/assets/img/avatar.jpg";
import { useRouter } from "vue-router";
import { onMounted, ref, watch } from "vue";
import { httpGet } from "@/utils/http";
import { ElMessage } from "element-plus";
import { UserFilled } from "@element-plus/icons-vue";
import { checkSession, getLicenseInfo, getSystemInfo } from "@/store/cache";
import { removeUserToken } from "@/store/session";
import LoginDialog from "@/components/LoginDialog.vue";
import { useSharedStore } from "@/store/sharedata";
import ConfigDialog from "@/components/UserInfoDialog.vue";
import { showMessageError } from "@/utils/dialog";

const isCollapse = ref(true);
const router = useRouter();
const logo = ref("");
const mainNavs = ref([]);
const moreNavs = ref([]);
// const curPath = ref(router.currentRoute.value.path);
const curPath = ref();

const title = ref("");
const showNoticeLogin = ref(false);
// const mainWinHeight = window.innerHeight - 50;

/**
 * 从路径名中提取第一个路径段
 * @param pathname - URL 的路径名部分，例如 '/chat/12345'
 * @returns 第一个路径段（不含斜杠），例如 'chat'，如果不存在则返回 null
 */
const extractFirstSegment = (pathname) => {
  const segments = pathname.split("/").filter((segment) => segment.length > 0);
  return segments.length > 0 ? segments[0] : null;
};
const getFirstPathSegment = (url) => {
  try {
    // 尝试使用 URL 构造函数解析完整的 URL
    const parsedUrl = new URL(url);
    return extractFirstSegment(parsedUrl.pathname);
  } catch (error) {
    // 如果解析失败，假设是相对路径，使用当前窗口的位置作为基准
    if (typeof window !== "undefined") {
      const parsedUrl = new URL(url, window.location.origin);
      return extractFirstSegment(parsedUrl.pathname);
    }
    // 如果无法解析，返回 null
    return null;
  }
};
const loginUser = ref({});
const mainWinHeight = loginUser.value.id
  ? window.innerHeight
  : window.innerHeight;

const version = ref(process.env.VUE_APP_VERSION);
const routerViewKey = ref(0);
const showConfigDialog = ref(false);
const license = ref({ de_copy: true });
const docsURL = ref(process.env.VUE_APP_DOCS_URL);
const gitURL = ref(process.env.VUE_APP_GIT_URL);

const store = useSharedStore();
const show = ref(false);
watch(
  () => store.showLoginDialog,
  (newValue) => {
    show.value = newValue;
  }
);

// 监听路由变化
// router.beforeEach((to, from, next) => {
//   curPath.value = to.path;
//   next();
// });

if (curPath.value === "/external") {
  curPath.value = router.currentRoute.value.query.url;
}
const changeNav = (item) => {
  curPath.value = item.url;
  if (item.url.indexOf("http") !== -1) {
    // 外部链接
    router.push({ name: "ExternalLink", query: { url: item.url } });
  } else {
    router.push(item.url);
  }
};

onMounted(() => {
  getSystemInfo()
    .then((res) => {
      logo.value = res.data.logo;
      title.value = res.data.title;
    })
    .catch((e) => {
      ElMessage.error("获取系统配置失败：" + e.message);
    });
  // 获取菜单
  httpGet("/api/menu/list")
    .then((res) => {
      mainNavs.value = res.data;
      // 根据窗口的高度计算应该显示多少菜单
      const rows = Math.floor((window.innerHeight - 100) / 90);
      if (res.data.length > rows) {
        mainNavs.value = res.data.slice(0, rows);
        moreNavs.value = res.data.slice(rows);
      }
    })
    .catch((e) => {
      ElMessage.error("获取系统菜单失败：" + e.message);
    });

  getLicenseInfo()
    .then((res) => {
      license.value = res.data;
    })
    .catch((e) => {
      license.value = { de_copy: false };
      showMessageError("获取 License 配置：" + e.message);
    });
  curPath.value = "/" + getFirstPathSegment(window.location.href);
  init();
});

const init = () => {
  checkSession()
    .then((user) => {
      loginUser.value = user;
    })
    .catch(() => {});
};

const logout = function () {
  httpGet("/api/user/logout")
    .then(() => {
      removeUserToken();
      store.setShowLoginDialog(true);
      store.setIsLogin(false);
      loginUser.value = {};
      // 刷新组件
      routerViewKey.value += 1;
    })
    .catch(() => {
      ElMessage.error("注销失败！");
    });
};

const loginCallback = () => {
  init();
  // 刷新组件
  routerViewKey.value += 1;
};
</script>

<style lang="stylus" scoped>
@import "@/assets/css/custom-scroll.styl"
@import "@/assets/css/home.styl"
</style>
