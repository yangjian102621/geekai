<template>
  <div :class="'sidebar ' + theme">
    <a class="logo w-full flex items-center" href="/" target="_blank">
      <img :src="logo" />
      <span class="text" v-show="!sidebar.collapse">{{ title }}</span>
    </a>

    <el-menu
      class="sidebar-el-menu"
      :default-active="onRoutes"
      :collapse="sidebar.collapse"
      background-color="transparent"
      text-color="var(--admin-sidebar-text-muted)"
      active-text-color="var(--admin-sidebar-text)"
      unique-opened
      router
    >
      <template v-for="item in items">
        <template v-if="item.subs">
          <el-sub-menu :index="item.index" :key="item.index">
            <template #title>
              <i :class="'iconfont icon-' + item.icon"></i>
              <span>{{ item.title }}</span>
            </template>
            <template v-for="subItem in item.subs">
              <el-sub-menu v-if="subItem.subs" :index="subItem.index" :key="subItem.index">
                <template #title>{{ subItem.title }}</template>
                <el-menu-item
                  v-for="(threeItem, i) in subItem.subs"
                  :key="i"
                  :index="threeItem.index"
                >
                  {{ threeItem.title }}
                </el-menu-item>
              </el-sub-menu>
              <el-menu-item v-else :index="subItem.index">
                <i v-if="subItem.icon" :class="'iconfont icon-' + subItem.icon"></i>
                {{ subItem.title }}
              </el-menu-item>
            </template>
          </el-sub-menu>
        </template>
        <template v-else>
          <el-menu-item :index="item.index" :key="item.index">
            <i :class="'iconfont icon-' + item.icon"></i>
            <template #title>{{ item.title }}</template>
          </el-menu-item>
        </template>
      </template>
    </el-menu>
  </div>
</template>

<script setup>
import { useSharedStore } from '@/store/sharedata'
import { setMenuItems, useSidebarStore } from '@/store/sidebar'
import { httpGet } from '@/utils/http'
import { ElMessage } from 'element-plus'
import { computed, ref, watch } from 'vue'
import { useRoute } from 'vue-router'

const title = ref('')
const logo = ref('')

// 加载系统配置
httpGet('/api/admin/config/get?key=system')
  .then((res) => {
    title.value = res.data.admin_title
    logo.value = res.data.logo
  })
  .catch((e) => {
    ElMessage.error('加载系统配置失败: ' + e.message)
  })
const store = useSharedStore()
const theme = ref(store.theme)
watch(
  () => store.theme,
  (val) => {
    theme.value = val
  }
)
const items = [
  {
    icon: 'home',
    index: '/admin/dashboard',
    title: '仪表盘',
  },

  {
    icon: 'user-fill',
    index: '/admin/user',
    title: '用户管理',
  },
  {
    icon: 'menu',
    index: '1',
    title: '应用管理',
    subs: [
      {
        index: '/admin/app',
        title: '应用列表',
        icon: 'sub-menu',
      },
      {
        index: '/admin/app/type',
        title: '应用分类',
        icon: 'chuangzuo',
      },
    ],
  },

  {
    icon: 'api-key',
    index: '/admin/apikey',
    title: 'API-KEY',
  },
  {
    icon: 'model',
    index: '/admin/chat/model',
    title: '模型管理',
  },
  {
    icon: 'recharge',
    index: '/admin/product',
    title: '充值产品',
  },
  {
    icon: 'order',
    index: '/admin/order',
    title: '充值订单',
  },
  {
    icon: 'reward',
    index: '/admin/redeem',
    title: '兑换码',
  },
  {
    icon: 'control',
    index: '/admin/functions',
    title: '函数管理',
  },
  {
    icon: 'menu',
    index: '2',
    title: '创作记录',
    subs: [
      {
        icon: 'prompt',
        index: '/admin/chats',
        title: '对话记录',
      },
      {
        icon: 'image',
        index: '/admin/images',
        title: '绘图记录',
      },
      {
        icon: 'mp3',
        index: '/admin/records/suno',
        title: 'Suno音乐',
      },
      
    ],
  },
  {
    icon: 'jimeng',
    index: '/admin/jimeng',
    title: '即梦AI',
    subs: [
      {
        icon: 'list',
        index: '/admin/jimeng/jobs',
        title: '任务列表',
      },

      {
        icon: 'config',
        index: '/admin/jimeng/config',
        title: '即梦设置',
      },
    ],
  },
  {
    icon: 'video',
    title: '视频生成',
    index: '/admin/video',
    subs: [
      {
        icon: 'config',
        index: '/admin/video/config',
        title: '生成配置',
      },

      {
        icon: 'list',
        index: '/admin/records/videos',
        title: '生成记录',
      },
      
    ],
  },
  {
    icon: 'ppt',
    title: 'AIPPT',
    index: '/admin/ppt',
    subs: [
      {
        icon: 'config',
        index: '/admin/ppt/config',
        title: 'PPT 生成配置',
      },
      {
        icon: 'list',
        index: '/admin/ppt/jobs',
        title: 'PPT 任务列表',
      },
    ],
  },

  {
    icon: 'moderation',
    index: '/admin/config/moderation',
    title: '文本审查',
    subs: [
      {
        icon: 'list',
        index: '/admin/moderation/list',
        title: '审核记录',
      },
      {
        icon: 'config',
        index: '/admin/moderation/config',
        title: '审查配置',
      },
    ],
  },
  {
    icon: 'role',
    index: '/admin/manger',
    title: '管理员',
  },

  {
    icon: 'config',
    index: 'config-center',
    title: '系统设置',
    subs: [
      {
        icon: 'config',
        index: '/admin/config/basic',
        title: '基础配置',
      },
      {
        icon: 'power',
        index: '/admin/config/power',
        title: '算力配置',
      },

      {
        icon: 'menu',
        index: '/admin/config/menu',
        title: '菜单配置',
      },
      {
        icon: 'recharge',
        index: '/admin/config/payment',
        title: '支付配置',
      },
      {
        icon: 'menu',
        index: '/admin/config/storage',
        title: '存储配置',
      },
      {
        icon: 'sms',
        index: '/admin/config/sms',
        title: '短信配置',
      },
      {
        icon: 'email',
        index: '/admin/config/smtp',
        title: '邮件配置',
      },
      {
        icon: 'plugin',
        index: '/admin/config/plugin',
        title: '插件配置',
      },
      {
        index: '/admin/config/wechat',
        title: '微信配置',
        icon: 'wechat',
      },
    ],
  },
  {
    icon: 'linggan',
    index: 'content-config',
    title: '文案配置',
    subs: [
      {
        icon: 'speaker',
        index: '/admin/config/notice',
        title: '公告配置',
      },
      {
        icon: 'info',
        index: '/admin/config/agreement',
        title: '用户协议',
      },
      {
        icon: 'info',
        index: '/admin/config/privacy',
        title: '隐私声明',
      },
      {
        icon: 'xmind',
        index: '/admin/config/markmap',
        title: '思维导图',
      },
    ],
  },
  {
    icon: 'log',
    index: '/admin/powerLog',
    title: '用户算力日志',
  },
  {
    icon: 'log',
    index: '/admin/loginLog',
    title: '用户登录日志',
  },
]

const route = useRoute()
const onRoutes = computed(() => {
  return route.path
})

const sidebar = useSidebarStore()
setMenuItems(items)
</script>

<style scoped lang="scss">
.sidebar {
  display: block;
  position: absolute;
  left: 0;
  top: 0;
  bottom: 0;
  overflow-y: scroll;
  background: var(--admin-sidebar-bg);
  backdrop-filter: blur(16px);
  -webkit-backdrop-filter: blur(16px);
  border-right: 1px solid var(--admin-sidebar-border);
  transition: background 0.2s ease, border-color 0.2s ease;

  .logo {
    display: flex;
    padding: 6px 15px;
    cursor: pointer;
    background: transparent;
    transition: background 0.2s ease;

    &:hover {
      background: var(--admin-header-hover);
    }

    img {
      height: 36px;
      padding: 2px;
      border-radius: 100%;
      background: #fff;
      border: 2px solid #754ff6;
    }

    .text {
      color: var(--admin-sidebar-text);
      font-weight: bold;
      padding: 12px 0 12px 10px;
      transition: width 2s ease, color 0.2s ease;
    }
  }

  ul {
    height: auto;
    min-height: 100%;
    background: transparent;

    .el-menu-item,
    .el-sub-menu {
      .iconfont {
        font-size: 16px;
        margin-right: 5px;
      }
    }

    .el-menu-item {
      transition: background 0.2s ease, color 0.2s ease;
      margin: 2px 8px;
      border-radius: 8px;

      &:hover {
        background: var(--admin-header-hover);
      }

      &.is-active {
        background: var(--admin-sidebar-active-bg);
        color: var(--admin-sidebar-text);
        box-shadow: inset 3px 0 0 var(--admin-sidebar-active-border);
      }
    }

    .el-sub-menu__title {
      transition: background 0.2s ease;
      margin: 2px 8px;
      border-radius: 8px;

      &:hover {
        background: var(--admin-header-hover);
      }
    }

  }

  .sidebar-el-menu {
    border-right: none;
    background: transparent;
  }

  .sidebar-el-menu:not(.el-menu--collapse) {
    width: 250px;
  }

  .el-menu {
    border-color: transparent;
  }
}

.sidebar::-webkit-scrollbar {
  width: 0;
}
</style>
