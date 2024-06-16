<template>
  <slot :toggleTooltip="() => (isVisible = !isVisible)" :id="togglerId" />
  <Teleport to="#tooltip-container">
    <Transition>
      <div
        v-if="isVisible"
        ref="tooltipRef"
        class="fixed flex justify-center items-center p-2 border rounded translate-y-full -translate-x-full bg-white/50 backdrop-blur-[1px]"
      >
        <slot name="tooltip" />
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, onBeforeMount, watch, nextTick } from "vue";
import type { Ref } from "vue";
import { v4 as uuid } from "uuid";
import { useTrackCoordinates } from "../composables/useCoordinatesTracker";

const isVisible = ref(false);
const togglerId = ref("");
const tooltipRef: Ref<HTMLElement | null> = ref(null);

onBeforeMount(() => {
  togglerId.value = uuid();
});

watch(isVisible, (v) => {
  if (v) {
    nextTick(() => {
      const foundTogglerEl = document.getElementById(togglerId.value);
      if (foundTogglerEl && tooltipRef.value) {
        const { x, y } = useTrackCoordinates(foundTogglerEl);

        tooltipRef.value.style.left = x.value + "px";
        tooltipRef.value.style.top = y.value + "px";
      }
    });
  }
});
</script>

<style>
.v-enter-active,
.v-leave-active {
  transition: opacity 200ms ease-out;
}

.v-enter-from,
.v-leave-to {
  opacity: 0;
}
</style>
