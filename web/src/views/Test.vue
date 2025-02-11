<template>
  <div>
   {{data}}
  </div>
</template>

<script setup>
import {onMounted, ref} from "vue";

const data = ref('abc')

onMounted(() => {
  
})

</script>
