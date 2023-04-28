<template>
  <div class="system-config">
    {{ title }}
  </div>
</template>

<script>
import {defineComponent} from "vue";

export default defineComponent({
  name: 'UserList',
  data() {
    return {
      title: "用户管理",
    }
  },
})
</script>

<style lang="stylus" scoped>
.system-config {

}
</style>