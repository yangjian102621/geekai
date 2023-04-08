<template>
  <div>{{ title }}</div>
</template>

<script>
import { defineComponent } from "vue"
import {isMobile} from "@/utils/libs";

export default defineComponent({
  name: 'HomePage',
  data () {
    return {
      title: "Loading page...",

    }
  },
  mounted() {
    if (isMobile()) {
      this.$router.push("mobile");
    } else {
      this.$router.push("plus");
    }

  }
})
</script>
