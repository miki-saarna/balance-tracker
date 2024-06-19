<template>
  <Tooltip>
    <template #default="{ toggleTooltip, id }">
      <button @click="toggleTooltip" :id="id" class="mt-4">
        <FunnelIcon class="w-5" />
      </button>
    </template>
    <template #tooltip>
      <ul>
        <li v-for="type of accountTypes">
          <button
            class="group flex items-center cursor-pointer"
            @click="toggleAccountTypeFilter(type)"
          >
            <div class="w-3.5">
              <CheckIcon
                v-if="checkIfFilterIsApplied(type)"
                class="group-hover:text-blue-400"
              />
            </div>
            <span class="group-hover:text-blue-400 ml-2 text-xs">{{
              type
            }}</span>
          </button>
        </li>
      </ul>
    </template>
  </Tooltip>
</template>

<script setup lang="ts">
import { ref, onBeforeMount, toRef, computed, reactive, watch } from "vue";
import type { Ref, ComputedRef, PropType } from "vue";
import Tooltip from "../../Tooltip.vue";
import { FunnelIcon, CheckIcon } from "@heroicons/vue/24/outline";
import type { FilterOptions, AccountTypesKeys } from "./types";

const props = defineProps({
  filterOptions: {
    type: Object as PropType<FilterOptions>,
    default: {},
  },
  accountTypes: {
    type: Array as PropType<AccountTypesKeys>,
    default: [],
  },
});

const toggleAccountTypeFilter = (type: AccountTypesKeys[number]) => {
  const selected = !!props.filterOptions.accountType[type];
  if (selected) {
    delete props.filterOptions.accountType[type];
  } else {
    props.filterOptions.accountType[type] = true;
  }
};

const checkIfFilterIsApplied = (type: AccountTypesKeys[number]) => {
  return !props.filterOptions.accountType[type];
};
</script>
