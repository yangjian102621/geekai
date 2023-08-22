<template>
  <div class="admin-home" v-if="isLogin">
    <admin-sidebar/>
    <div class="content-box" :class="{ 'content-collapse': sidebar.collapse }">
      <admin-header/>
      <admin-tags/>
      <div class="content">
        <router-view v-slot="{ Component }">
          <transition name="move" mode="out-in">
            <keep-alive :include="tags.nameList">
              <component :is="Component"></component>
            </keep-alive>
          </transition>
        </router-view>
      </div>
    </div>
  </div>
</template>
<script setup>
import {useSidebarStore} from '@/store/sidebar';
import {useTagsStore} from '@/store/tags';
import AdminHeader from "@/components/admin/AdminHeader.vue";
import AdminSidebar from "@/components/admin/AdminSidebar.vue";
import AdminTags from "@/components/admin/AdminTags.vue";
import {useRouter} from "vue-router";
import {checkAdminSession} from "@/action/session";
import {ref} from "vue";

const sidebar = useSidebarStore();
const tags = useTagsStore();
const isLogin = ref(false)

// 获取会话信息
const router = useRouter();
checkAdminSession().then(() => {
  isLogin.value = true
}).catch(() => {
  router.replace('/admin/login')
})
</script>

<style scoped lang="stylus">
@import '@/assets/css/main.css';
@import '@/assets/css/color-dark.css';
@import '@/assets/iconfont/iconfont.css';
</style>
