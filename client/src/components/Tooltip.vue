<template>
  <slot :toggleTooltip="() => (isVisible = !isVisible)" />
  <Teleport to="body" v-if="isVisible">
    <slot name="tooltip" />
  </Teleport>
</template>

<script setup lang="ts">
import { shallowRef } from "vue";

const isVisible = shallowRef(false);
</script>
