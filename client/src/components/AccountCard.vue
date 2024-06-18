<template>
  <div
    class="flex justify-between py-4 border-b border-gray-300 final:border-none"
  >
    <div class="flex flex-col">
      <div class="font-semibold">{{ name }}</div>
      <div class="text-sm">{{ type }}</div>
    </div>

    <div class="self-start flex items-center">
      <div class="font-semibold">${{ balance }}</div>
      <Tooltip>
        <template #default="{ toggleTooltip, id }">
          <button @click="toggleTooltip" :id="id" class="ml-4">
            <EllipsisVerticalIcon class="w-5" />
          </button>
        </template>
        <template #tooltip>
          <button @click="$emit('refresh-balance')">
            <ArrowPathIcon class="w-5" />
          </button>

          <!-- currently not saving `persistent_account_id` within the db -->
          <button class="ml-4" @click="$emit('remove-account')">
            <TrashIcon class="w-5 text-red-500" />
          </button>
        </template>
      </Tooltip>
    </div>
  </div>
</template>

<script setup lang="ts">
import Tooltip from "./Tooltip.vue";
import {
  ArrowPathIcon,
  TrashIcon,
  EllipsisVerticalIcon,
} from "@heroicons/vue/24/solid";

defineProps({
  name: {
    type: String,
    required: true,
  },
  type: {
    default: "",
    type: String,
  },
  balance: {
    type: Number,
    required: true,
  },
});
</script>
