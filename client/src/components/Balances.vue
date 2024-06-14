<template>
  <Link v-if="linkToken" :linkToken="linkToken" :accessTokens="accessTokens" />

  <div class="mt-4 border-t border-gray-300">
    <div v-for="[accessToken, accountsList] of Object.entries(accounts)">
      <div v-for="account of accountsList">
        <div
          :key="account.account_id"
          class="flex justify-between py-4 border-b border-gray-300 final:border-none"
        >
          <div class="flex flex-col">
            <div>{{ account.name }}</div>
            <div>{{ account.subtype }}</div>
            <div>${{ account.balances.available }}</div>
          </div>

          <div class="">
            <button @click="() => refreshBalance(accessToken)">
              <ArrowPathIcon class="w-5" />
            </button>

            <!-- currently not saving `persistent_account_id` within the db -->
            <button
              class="ml-4"
              @click="
                () => {
                  removeAccount(
                    accessToken,
                    account.account_id,
                    account.persistent_account_id
                  );
                }
              "
            >
              <TrashIcon class="w-5 text-red-500" />
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>

  <RenderTotalBalance class="mt-4" />
</template>

<script setup lang="ts">
import { shallowRef, ref, onBeforeMount, watch, h } from "vue";
import { usePlaidTokens } from "../composables/useAccessTokensRetrieval";
import { AccessTokensResponse, getAccessTokens } from "../utils/db";
import {
  Link,
  getAccountsBalances,
  AccountsBalancesResponse,
} from "../utils/plaid_api";
import { ArrowPathIcon, TrashIcon } from "@heroicons/vue/24/solid";

type AccountsByAccessToken = {
  [key: string]: any[]; // update with correct type from Plaid
};

const props = defineProps({
  msg: String,
});

const accounts = ref({});

// const accountsArray = computed(() => {
//   return Object.entries(accounts.value)
// })

const { accessTokens, genAccessTokens, linkToken, genLinkToken } =
  usePlaidTokens();

async function refreshBalance(accessToken: string): Promise<void> {
  // getAccountsBalances func - partially Plaid and partially db
  const data: AccountsBalancesResponse | void = await getAccountsBalances(
    accessToken
  );
  if (data) {
    accounts.value = {
      ...accounts.value,
      [accessToken]: data.accounts,
    };
  }
}

async function removeAccount(
  accessToken: string,
  account_id: string,
  persistent_id: string
): Promise<void> {
  // 2nd parameter necessary?

  let res;
  try {
    res = await fetch("http://localhost:8000/api/account/delete", {
      method: "DELETE",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ account_id }),
    });
  } catch (err) {
    console.log("There was an error deleting the account", err);
    return;
  }

  const data = await res.json();
  // console.log("data: ", data)

  const modifyAccounts = { ...accounts };

  const accountsBelongingToAccessToken = modifyAccounts[accessToken];
  if (accountsBelongingToAccessToken.length > 1) {
    const deletionIdx = accountsBelongingToAccessToken.findIndex(
      (account) => account.account_id === account_id
    );
    modifyAccounts[accessToken].splice(deletionIdx, 1);
    accounts.value = modifyAccounts;
  } else {
    delete modifyAccounts[accessToken];
    accounts.value = modifyAccounts;
  }
}

function RenderTotalBalance() {
  const accountsList = Object.values(accounts.value);
  if (!accountsList.length) return;
  const sum = accountsList
    .filter((account) => !!account)
    .flat()
    .reduce((acc, account) => acc + account.balances.available, 0);
  return h("div", `Total: $${sum}`);
}

onBeforeMount(async () => {
  await genAccessTokens();
  await genLinkToken();
});

watch(accessTokens, (tokens) => {
  for (const accessToken of tokens) {
    refreshBalance(accessToken);
  }

  setTimeout(() => {
    RenderTotalBalance();
  }, 1000);
});
</script>
