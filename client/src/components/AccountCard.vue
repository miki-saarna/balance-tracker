<template>
  <div
    class="flex justify-between py-4 border-b border-gray-300 final:border-none"
  >
    <div class="flex flex-col">
      <div>{{ name }}</div>
      <div>{{ type }}</div>
      <div>${{ balance }}</div>
    </div>

    <div>
      <button @click="$emit('refresh-balance')">
        <ArrowPathIcon class="w-5" />
      </button>

      <!-- currently not saving `persistent_account_id` within the db -->
      <button class="ml-4" @click="$emit('remove-account')">
        <TrashIcon class="w-5 text-red-500" />
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ArrowPathIcon, TrashIcon } from "@heroicons/vue/24/solid";

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
