<template>
  <div class="tags" v-if="tags.show">
    <ul>
      <li
          class="tags-li"
          v-for="(item, index) in tags.list"
          :class="{ active: isActive(item.path) }"
          :key="index"
      >
        <router-link :to="item.path" class="tags-li-title">{{ item.title }}</router-link>
        <el-icon @click="closeTags(index)">
          <Close/>
        </el-icon>
      </li>
    </ul>
    <div class="tags-close-box">
      <el-dropdown @command="handleTags">
        <el-button size="small" type="info">
          标签选项
          <el-icon class="el-icon--right">
            <arrow-down/>
          </el-icon>
        </el-button>
        <template #dropdown>
          <el-dropdown-menu size="small">
            <el-dropdown-item command="other">关闭其他</el-dropdown-item>
            <el-dropdown-item command="all">关闭所有</el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </div>
  </div>
</template>

<script setup>
import {useTagsStore} from '@/store/tags';
import {onBeforeRouteUpdate, useRoute, useRouter} from 'vue-router';
import {ArrowDown, Close} from "@element-plus/icons-vue";
import {checkAdminSession} from "@/action/session";
import {ElMessageBox} from "element-plus";

const router = useRouter();
checkAdminSession().catch(() => {
  ElMessageBox({
    title: '提示',
    message: "当前会话已经失效，请重新登录",
    confirmButtonText: 'OK',
    callback: () => router.replace('/admin/login')
  });
})
const isActive = (path) => {
  return path === route.fullPath;
};

const tags = useTagsStore();
const route = useRoute();
// 关闭单个标签
const closeTags = (index) => {
  const delItem = tags.list[index];
  tags.delTagsItem(index);
  const item = tags.list[index] ? tags.list[index] : tags.list[index - 1];
  if (item) {
    delItem.path === route.fullPath && router.push(item.path);
  } else {
    router.push('/admin');
  }
};

// 设置标签
const setTags = (route) => {
  const isExist = tags.list.some(item => {
    return item.path === route.fullPath;
  });
  if (!isExist) {
    if (tags.list.length >= 8) tags.delTagsItem(0);
    tags.setTagsItem({
      name: route.name,
      title: route.meta.title,
      path: route.fullPath
    });
  }
};
setTags(route);
onBeforeRouteUpdate(to => {
  setTags(to);
});

// 关闭全部标签
const closeAll = () => {
  tags.clearTags();
  router.push('/admin');
};
// 关闭其他标签
const closeOther = () => {
  const curItem = tags.list.filter(item => {
    return item.path === route.fullPath;
  });
  tags.closeTagsOther(curItem);
};
const handleTags = (command) => {
  command === 'other' ? closeOther() : closeAll();
};

// 关闭当前页面的标签页
// tags.closeCurrentTag({
//     $router: router,
//     $route: route
// });
</script>

<style>
.tags {
  position: relative;
  height: 30px;
  overflow: hidden;
  background: #fff;
  padding-right: 120px;
  -webkit-box-shadow: 0 1px 4px rgba(0, 21, 41, .08);
  box-shadow: 0 1px 4px rgba(0, 21, 41, .08);
}

.tags ul {
  box-sizing: border-box;
  width: 100%;
  height: 100%;
}

.tags-li {
  display: flex;
  align-items: center;
  float: left;
  margin: 3px 5px 2px 3px;
  border-radius: 3px;
  font-size: 12px;
  overflow: hidden;
  cursor: pointer;
  height: 23px;
  border: 1px solid #e9eaec;
  background: #fff;
  padding: 0 5px 0 12px;
  color: #666;
  -webkit-transition: all 0.3s ease-in;
  -moz-transition: all 0.3s ease-in;
  transition: all 0.3s ease-in;
}

.tags-li:not(.active):hover {
  background: #f8f8f8;
}

.tags-li.active {
  color: #fff;
}

.tags-li-title {
  float: left;
  max-width: 80px;
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
  margin-right: 5px;
  color: #666;
}

.tags-li.active .tags-li-title {
  color: #fff;
}

.tags-close-box {
  position: absolute;
  right: 0;
  top: 2px;
  box-sizing: border-box;
  padding-top: 1px;
  text-align: center;
  width: 110px;
  height: 30px;
  background: #fff;
//box-shadow: -3px 0 15px 3px rgba(0, 0, 0, 0.1); z-index: 10;
}
</style>
