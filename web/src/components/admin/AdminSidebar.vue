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
      <template v-for="item in items" :key="item.path">
        <template v-if="item.subs">
          <el-sub-menu :index="item.index" :key="item.index">
            <template #title>
              <i :class="'iconfont icon-' + item.icon"></i>
              <span>{{ item.title }}</span>
            </template>
            <template v-for="subItem in item.subs" :key="subItem.index">
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
// 菜单信息架构：仪表盘 + 8 个语义分组，统一「分组 → 页面」两级心智模型。
// 分组按使用频率自上而下：日常运营（用户/交易）→ 能力供给（模型/AI/应用）→ 质量与安全 → 平台运维。
// 注意：叶子 index 必须与 router.js 中 /admin/* 路由路径一致，勿改。
const items = [
  {
    icon: 'home',
    index: '/admin/dashboard',
    title: '仪表盘',
  },
  {
    icon: 'user-fill',
    index: 'users',
    title: '用户管理',
    subs: [
      {
        icon: 'user-fill',
        index: '/admin/user',
        title: '用户列表',
      },
      {
        icon: 'power',
        index: '/admin/powerLog',
        title: '积分日志',
      },
      {
        icon: 'log',
        index: '/admin/loginLog',
        title: '登录日志',
      },
    ],
  },
  {
    icon: 'recharge',
    index: 'trading',
    title: '交易中心',
    subs: [
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
    ],
  },
  {
    icon: 'model',
    index: 'ai-center',
    title: '模型服务',
    subs: [
      {
        icon: 'model',
        index: '/admin/chat/model',
        title: '语言模型',
      },
      {
        icon: 'api-key',
        index: '/admin/apikey',
        title: 'API-KEY',
      },
      {
        icon: 'control',
        index: '/admin/functions',
        title: '函数管理',
      },
    ],
  },
  {
    icon: 'video',
    index: 'media-center',
    title: '多模态任务',
    subs: [
      {
        icon: 'jimeng',
        index: '/admin/jimeng/jobs',
        title: '即梦绘图任务',
      },
      {
        icon: 'jimeng',
        index: '/admin/jimeng/config',
        title: '即梦设置',
      },
      {
        icon: 'video',
        index: '/admin/video/config',
        title: '视频生成配置',
      },
      {
        icon: 'list',
        index: '/admin/records/videos',
        title: '视频生成记录',
      },
    ],
  },
  {
    icon: 'menu',
    index: 'content-center',
    title: '应用与内容',
    subs: [
      {
        icon: 'menu',
        index: '/admin/app',
        title: '应用列表',
      },
      {
        icon: 'chuangzuo',
        index: '/admin/app/type',
        title: '应用分类',
      },
      {
        icon: 'linggan',
        index: '/admin/library',
        title: '公共资料库',
      },
    ],
  },
  {
    icon: 'list',
    index: 'records',
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
        title: '音乐记录',
      },
    ],
  },
  {
    icon: 'moderation',
    index: 'moderation-center',
    title: '内容安全',
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
    icon: 'config',
    index: 'platform-center',
    title: '平台管理',
    subs: [
      {
        icon: 'role',
        index: '/admin/manger',
        title: '管理员',
      },
      {
        icon: 'config',
        index: '/admin/config/basic',
        title: '基础配置',
      },
      {
        icon: 'power',
        index: '/admin/config/power',
        title: '积分配置',
      },
      {
        icon: 'menu',
        index: '/admin/config/menu',
        title: '门户菜单',
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
        icon: 'wechat',
        index: '/admin/config/wechat',
        title: '微信配置',
      },
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
        title: '思维导图页',
      },
    ],
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
      border-radius: 4px;
      background: #fff;
    }

    .text {
      color: var(--admin-sidebar-text);
      font-weight: 500;
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
      margin: 1px 12px;
      border-radius: 4px;
      height: 40px;
      line-height: 40px;

      &:hover {
        background: var(--admin-header-hover);
        color: var(--admin-sidebar-text);
      }

      &.is-active {
        background: var(--admin-sidebar-active-bg);
        color: var(--admin-sidebar-text);
        font-weight: 500;
      }
    }

    .el-sub-menu__title {
      transition: background 0.2s ease;
      margin: 1px 12px;
      border-radius: 4px;
      height: 40px;
      line-height: 40px;

      &:hover {
        background: var(--admin-header-hover);
        color: var(--admin-sidebar-text);
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
