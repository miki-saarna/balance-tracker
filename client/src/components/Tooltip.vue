<template>
  <slot :toggleTooltip="toggleTooltip" :id="togglerId" />
  <Teleport to="#tooltip-container" v-if="isVisibleTeleport">
    <Transition>
      <div
        v-if="isVisibleTransition"
        ref="tooltipRef"
        class="fixed flex justify-center items-center p-2 border rounded translate-y-full -translate-x-full bg-white/75 backdrop-blur-[1px]"
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

const isVisibleTeleport = ref(false);
const isVisibleTransition = ref(false);
const togglerId = ref("");
const tooltipRef: Ref<HTMLElement | null> = ref(null);

onBeforeMount(() => {
  togglerId.value = uuid();
});

const toggleTooltip = (): void => {
  if (isVisibleTeleport.value) {
    isVisibleTransition.value = false;
    setTimeout(() => {
      isVisibleTeleport.value = false;
    });
  } else {
    isVisibleTeleport.value = true;
    setTimeout(() => {
      isVisibleTransition.value = true;
    });
  }
};

watch(isVisibleTransition, (v) => {
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

<style scoped>
.v-enter-active,
.v-leave-active {
  transition: opacity 200ms ease-out;
}

.v-enter-from,
.v-leave-to {
  opacity: 0;
}
</style>
