<template>
  <div>{{ title }}</div>
</template>

<script>
import { defineComponent } from "vue"

export default defineComponent({
  name: 'NotFound',
  data () {
    return {
      title: "404 Page",

    }
  },
})
</script>
