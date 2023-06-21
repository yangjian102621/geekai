<template>
  <div class="sidebar">
    <el-menu
        class="sidebar-el-menu"
        :default-active="onRoutes"
        :collapse="sidebar.collapse"
        background-color="#324157"
        text-color="#bfcbd9"
        active-text-color="#20a0ff"
        unique-opened
        router
    >
      <template v-for="item in items">
        <template v-if="item.subs">
          <el-sub-menu :index="item.index" :key="item.index">
            <template #title>
              <i :class="'iconfont icon-'+item.icon"></i>
              <span>{{ item.title }}</span>
            </template>
            <template v-for="subItem in item.subs">
              <el-sub-menu
                  v-if="subItem.subs"
                  :index="subItem.index"
                  :key="subItem.index"
              >
                <template #title>{{ subItem.title }}</template>
                <el-menu-item v-for="(threeItem, i) in subItem.subs" :key="i" :index="threeItem.index">
                  {{ threeItem.title }}
                </el-menu-item>
              </el-sub-menu>
              <el-menu-item v-else :index="subItem.index">
                <i v-if="subItem.icon" :class="'iconfont icon-'+subItem.icon"></i>
                {{ subItem.title }}
              </el-menu-item>
            </template>
          </el-sub-menu>
        </template>
        <template v-else>
          <el-menu-item :index="item.index" :key="item.index">
            <i :class="'iconfont icon-'+item.icon"></i>
            <template #title>{{ item.title }}</template>
          </el-menu-item>
        </template>
      </template>
    </el-menu>
  </div>
</template>

<script setup>
import {computed} from 'vue';
import {useSidebarStore} from '@/store/sidebar';
import {useRoute} from 'vue-router';

const items = [
  {
    icon: 'home',
    index: '/admin/welcome',
    title: '系统首页',
  },
  {
    icon: 'config',
    index: '/admin/system',
    title: '系统设置',
  },

  {
    icon: 'user-fill',
    index: '/admin/user',
    title: '用户管理',
  },

  {
    icon: 'role',
    index: '/admin/role',
    title: '角色管理',
  },
  {
    icon: 'api-key',
    index: '/admin/apikey',
    title: 'API-KEY 管理',
  },
  {
    icon: 'log',
    index: '/admin/loginLog',
    title: '用户登录日志',
  },
  {
    icon: 'menu',
    index: '1',
    title: '常用模板页面',
    subs: [
      {
        index: '/admin/demo/form',
        title: '表单页面',
      },
      {
        index: '/admin/demo/table',
        title: '常用表格',
      },
      {
        index: '/admin/demo/import',
        title: '导入Excel',
      },
      {
        index: '/admin/demo/editor',
        title: '富文本编辑器',
      },
    ],
  },
];

const route = useRoute();
const onRoutes = computed(() => {
  return route.path;
});

const sidebar = useSidebarStore();
</script>

<style scoped lang="stylus">
.sidebar {
  display: block;
  position: absolute;
  left: 0;
  top: 70px;
  bottom: 0;
  overflow-y: scroll;

  ul {
    height: 100%;

    .el-menu-item, .el-sub-menu {
      .iconfont {
        font-size 16px;
        margin-right 5px;
      }
    }
  }

  .sidebar-el-menu:not(.el-menu--collapse) {
    width: 250px;
  }
}

.sidebar::-webkit-scrollbar {
  width: 0;
}

</style>
