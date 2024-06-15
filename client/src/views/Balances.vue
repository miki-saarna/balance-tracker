<template>
  <Link
    v-if="linkToken"
    :linkToken="linkToken"
    :accessTokens="toRef(accessTokens)"
    :refreshBalance="refreshBalance"
  />

  <div class="mt-4 border-t border-gray-300">
    <div v-for="[accessToken, accountsList] of accountsEntries">
      <div v-for="account of accountsList">
        <AccountCard
          :key="account.account_id"
          :name="account.name"
          :type="account.subtype"
          :balance="account.balances.available"
          @refresh-balance="refreshBalance(accessToken)"
          @remove-account="removeAccountHandler(accessToken, account)"
        />
      </div>
    </div>
  </div>

  <div v-if="renderTotalBalance" class="mt-4">
    Total: ${{ renderTotalBalance }}
  </div>
</template>

<script setup lang="ts">
import { ref, onBeforeMount, toRef, computed } from "vue";
import type { Ref, ComputedRef } from "vue";
import { getAccessTokens, removeAccount } from "../utils/db";
import type { AccessTokensResponse } from "../utils/db";
import { Link, getAccountsBalances, genLinkToken } from "../utils/plaid_api";
import type { AccountsBalancesResponse } from "../utils/plaid_api";
import AccountCard from "../components/AccountCard.vue";

type AccountsByAccessToken = {
  [key: string]: any[]; // update with correct type from Plaid
};

type Account = {
  [key: string]: any;
};

const props = defineProps({
  msg: String,
});

const linkToken: Ref<string> = ref("");
const accounts: Ref<{ [key: string]: Account[] }> = ref({});
const accessTokens: Ref<AccessTokensResponse["access_tokens"]> = ref([]);

const accountsEntries: ComputedRef<[string, Account[]][]> = computed(() => {
  return Object.entries(accounts.value);
});

const renderTotalBalance: ComputedRef<number | void> = computed(() => {
  const accountsList = Object.values(accounts.value);
  if (!accountsList.length) return;
  const sum = accountsList
    .filter((account) => !!account)
    .flat()
    .reduce((acc, account) => acc + account.balances.available, 0);
  return sum;
});

async function refreshBalance(accessToken: string): Promise<void> {
  // getAccountsBalances func - partially Plaid and partially db
  const data: AccountsBalancesResponse | void = await getAccountsBalances(
    accessToken
  );

  if (data) {
    accounts.value[accessToken] = data.accounts;
  }
}

async function removeAccountHandler(
  accessToken: string,
  account: Account
): Promise<void> {
  try {
    await removeAccount(account.account_id, account.persistent_account_id);
  } catch (err) {
    console.log(
      `There was an error trying to remove the account with ID ${account.account_id}:`,
      err
    );
    return;
  }

  const modifyAccounts = { ...accounts.value };

  const accountsBelongingToAccessToken = modifyAccounts[accessToken];
  if (accountsBelongingToAccessToken.length > 1) {
    const deletionIdx = accountsBelongingToAccessToken.findIndex(
      (account) => account.account_id === account.account_id
    );
    modifyAccounts[accessToken].splice(deletionIdx, 1);
    accounts.value = modifyAccounts;
  } else {
    delete modifyAccounts[accessToken];
    accounts.value = modifyAccounts;
  }
}

onBeforeMount(async () => {
  try {
    const accessTokensRes = await getAccessTokens();
    if (accessTokensRes) {
      accessTokens.value = accessTokensRes.access_tokens;
    }
    for (const accessToken of accessTokens.value) {
      refreshBalance(accessToken);
    }
  } catch (err) {
    console.error(`There was an error retrieving access_tokens from DB:`, err);
  }

  try {
    const linkTokenRes = await genLinkToken();
    if (linkTokenRes) {
      linkToken.value = linkTokenRes.link_token;
    }
  } catch (err) {
    console.error("There was an error generating a link_token:", err);
  }
});
</script>
