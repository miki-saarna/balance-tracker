<template>
  <p>Link token: {{ linkToken }}</p>
  <ol>
    <li v-for="(token, idx) of accessTokens">{{ idx + 1 }}. {{ token }}</li>
  </ol>
  <Link v-if="linkToken" :linkToken="linkToken" :accessTokens="accessTokens" />
</template>

<script setup lang="ts">
import { shallowRef, ref, onBeforeMount } from "vue";
import { usePlaidTokens } from "../composables/useAccessTokensRetrieval";
import { AccessTokensResponse, getAccessTokens } from "../utils/db";
import { Link } from "../utils/plaid_api";

const props = defineProps({
  msg: String,
});

const accounts = ref({});

const { accessTokens, genAccessTokens, linkToken, genLinkToken } =
  usePlaidTokens();

onBeforeMount(async () => {
  await genAccessTokens();
  await genLinkToken();
});
</script>
