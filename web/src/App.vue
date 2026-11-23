<template>
  <el-config-provider>
    <router-view />
  </el-config-provider>
</template>

<script setup>
import { ElConfigProvider } from 'element-plus'
import { initThumbTemplate } from '@/utils/libs'

initThumbTemplate().catch(() => {})
</script>

<style lang="scss">
html,
body {
  margin: 0;
  padding: 0;
}

#app {
  font-family: 'The Future', Arial, 'PingFang SC', 'Microsoft YaHei', 'Noto Sans SC',
    sans-serif;
  -webkit-font-smoothing: antialiased;
}

.ellipsis {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
</style>
