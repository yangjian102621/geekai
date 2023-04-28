<template>
  <div class="system-config">
    {{ title }}
  </div>
</template>

<script>
import {defineComponent} from "vue";

export default defineComponent({
  name: 'RoleList',
  data() {
    return {
      title: "角色管理",

    }
  },
})
</script>

<style lang="stylus" scoped>
.system-config {

}
</style>