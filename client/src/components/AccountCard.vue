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
        <template #default="{ toggleTooltip }">
          <EllipsisVerticalIcon
            v-slot="{ toggleTooltip }"
            @click="toggleTooltip"
            class="ml-4 w-5 cursor-pointer"
          />
        </template>
        <template #tooltip>
          <div
            class="fixed flex justify-center items-center p-2 border rounded"
          >
            <button @click="$emit('refresh-balance')">
              <ArrowPathIcon class="w-5" />
            </button>

            <!-- currently not saving `persistent_account_id` within the db -->
            <button class="ml-4" @click="$emit('remove-account')">
              <TrashIcon class="w-5 text-red-500" />
            </button>
          </div>
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
