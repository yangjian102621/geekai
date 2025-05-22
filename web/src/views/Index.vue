<template>
  <div class="index-page">
    <!-- 主题切换 -->
    <ThemeChange />
    <div class="menu-box">
      <el-menu mode="horizontal" :ellipsis="false">
        <div class="menu-item">
          <img :src="logo" class="logo" alt="Geek-AI" />
        </div>
        <div class="menu-item">
          <span v-if="!license.de_copy">
            <el-tooltip v-if="!license.de_copy" class="box-item" content="部署文档" placement="bottom">
              <a :href="docsURL" class="link-button mr-2" target="_blank">
                <i class="iconfont icon-book"></i>
              </a>
            </el-tooltip>
            <el-tooltip v-if="!license.de_copy" class="box-item" content="项目源码" placement="bottom">
              <a :href="gitURL" class="link-button" target="_blank">
                <i class="iconfont icon-github"></i>
              </a>
            </el-tooltip>
          </span>

          <span v-if="!isLogin">
            <!-- <el-button @click="router.push('/login')" class="shadow" round
              >登录</el-button
            >
            <el-button @click="router.push('/register')" class="shadow" round
              >注册</el-button
            > -->
            <el-button @click="router.push('/login')" class="btn-go animate__animated animate__pulse animate__infinite" round>登录/注册</el-button>
          </span>
        </div>
      </el-menu>
    </div>
    <div class="content">
      <div style="height: 158px"></div>
      <h1 class="animate__animated animate__backInDown">
        {{ title }}
      </h1>
      <div class="msg-text cursor-ani">
        <span v-for="(char, index) in displayedChars" :key="index" :style="{ color: rainbowColor(index) }">
          {{ char }}
        </span>
      </div>

      <div class="navs animate__animated animate__backInDown">
        <el-space wrap :size="14">
          <div v-for="item in navs" :key="item.url" class="nav-item-box" @click="router.push(item.url)">
            <i :class="'iconfont ' + iconMap[item.url]"></i>
            <div>{{ item.name }}</div>
          </div>
        </el-space>
      </div>
    </div>

    <footer-bar />
  </div>
</template>

<script setup>
import { onMounted, ref } from "vue";
import { useRouter } from "vue-router";
import FooterBar from "@/components/FooterBar.vue";
import ThemeChange from "@/components/ThemeChange.vue";
import { httpGet } from "@/utils/http";
import { ElMessage } from "element-plus";
import { checkSession, getLicenseInfo, getSystemInfo } from "@/store/cache";
import { isMobile } from "@/utils/libs";

const router = useRouter();

if (isMobile()) {
  router.push("/mobile/index");
}

const title = ref("");
const logo = ref("");
const slogan = ref("");
const license = ref({ de_copy: true });

const isLogin = ref(false);
const docsURL = ref(import.meta.env.VITE_APP_DOCS_URL);
const gitURL = ref(import.meta.env.VITE_APP_GIT_URL);
const navs = ref([]);

const iconMap = ref({
  "/chat": "icon-chat",
  "/mj": "icon-mj",
  "/sd": "icon-sd",
  "/dalle": "icon-dalle",
  "/images-wall": "icon-image",
  "/suno": "icon-suno",
  "/xmind": "icon-xmind",
  "/apps": "icon-app",
  "/member": "icon-vip-user",
  "/invite": "icon-share",
  "/luma": "icon-luma",
});

const displayedChars = ref([]);
const initAnimation = ref("");
let timer = null; // 定时器句柄

// 初始化间隔时间和随机时间数组
const interTime = ref(50);
const interArr = [90, 100, 70, 88, 80, 110, 85, 400, 90, 99];

onMounted(() => {
  getSystemInfo()
    .then((res) => {
      title.value = res.data.title;
      logo.value = res.data.logo;
      slogan.value = res.data.slogan;
      if (timer) clearInterval(timer); // 清除定时器
      timer = setInterval(setContent, interTime.value);
    })
    .catch((e) => {
      ElMessage.error("获取系统配置失败：" + e.message);
    });

  getLicenseInfo()
    .then((res) => {
      license.value = res.data;
    })
    .catch((e) => {
      license.value = { de_copy: false };
      ElMessage.error("获取 License 配置失败：" + e.message);
    });

  httpGet("/api/menu/list?index=1")
    .then((res) => {
      navs.value = res.data;
    })
    .catch((e) => {
      ElMessage.error("获取导航菜单失败：" + e.message);
    });

  checkSession()
    .then(() => {
      isLogin.value = true;
    })
    .catch(() => {});
});
// 打字机内容逐字符显示
const setContent = () => {
  if (initAnimation.value.length >= slogan.value.length) {
    // 文本已全部输出
    initAnimation.value = "";
    displayedChars.value = [];
    if (timer) clearInterval(timer);
    timer = setInterval(setContent, interTime.value);
    return;
  } else {
    const nextChar = slogan.value.charAt(initAnimation.value.length);
    initAnimation.value += slogan.value.charAt(initAnimation.value.length); // 逐字符追加
    displayedChars.value.push(nextChar);
    interTime.value = interArr[Math.floor(Math.random() * interArr.length)]; // 设置随机间隔
    if (timer) clearInterval(timer);
    timer = setInterval(setContent, interTime.value);
  }
};
// 计算彩虹色
const rainbowColor = (index) => {
  const hue = (index * 40) % 360; // 每个字符间隔40度，形成彩虹色
  return `hsl(${hue}, 90%, 50%)`; // 色调(hue)，饱和度(70%)，亮度(50%)
};
</script>

<style lang="stylus" scoped>
@import "../assets/css/index.styl"
</style>
