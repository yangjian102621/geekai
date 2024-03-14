<script lang="ts" setup>
import { hasPermission } from "@/directives/permission";
defineProps<{
  permission: string | string[] | true;
}>();
</script>
<template>
  <slot v-if="hasPermission(permission)" />
  <slot v-else name="none" />
</template>
