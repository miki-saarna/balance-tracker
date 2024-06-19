<template>
  <Link
    v-if="linkToken"
    :linkToken="linkToken"
    :accessTokens="toRef(accessTokens)"
    :refreshBalance="refreshBalance"
  />
  <!-- class="block" -->

  <AccountsFilter
    :filterOptions="filterOptions"
    :account-types="accountTypes"
  />

  <div class="mt-4 border-t border-gray-300">
    <AccountCard
      v-for="account of filteredAccounts"
      :key="account.account_id"
      :name="account.name"
      :type="account.subtype"
      :balance="account.balances.available"
      @refresh-balance="refreshBalance(account.access_token)"
      @remove-account="removeAccountHandler(account.access_token, account)"
    />
  </div>

  <div
    v-if="renderTotalBalance"
    class="mt-4 flex justify-between font-semibold"
  >
    <div>Total</div>
    <div class="mr-9">${{ renderTotalBalance }}</div>
  </div>
</template>

<script setup lang="ts">
import { ref, onBeforeMount, toRef, computed, reactive, watch } from "vue";
import type { Ref, ComputedRef } from "vue";
import { getAccessTokens, removeAccount } from "../utils/db";
import type { AccessTokensResponse } from "../utils/db";
import { Link, getAccountsBalances, genLinkToken } from "../utils/plaid_api";
import type { AccountsBalancesResponse } from "../utils/plaid_api";
import AccountCard from "../components/AccountCard.vue";
import AccountsFilter from "../components/balances/filters/filter.vue";
import type {
  FilterOptions,
  AccountTypesKeys,
  AccountTypes,
  PlaidAccountTypes,
} from "../components/balances/filters/types";

interface Account {
  account_id: string;
  balances: Record<"available", number>;
  name: string;
  subtype: PlaidAccountTypes;
  [key: string]: unknown;
}

interface AccountInList extends Account {
  access_token: string;
}

type AccountsByAccessToken = {
  [key: string]: Account[]; // update with correct type from Plaid
};

const props = defineProps({
  msg: String,
});

const linkToken: Ref<string> = ref("");
const accounts: AccountsByAccessToken = reactive({});
const accessTokens: Ref<AccessTokensResponse["access_tokens"]> = ref([]);

const filterOptions: FilterOptions = reactive({
  accountType: {},
  balance: {},
});

const renderTotalBalance: ComputedRef<number> = computed(() => {
  return filteredAccounts.value.reduce(
    (acc, account) => acc + account.balances.available,
    0
  );
});

const filteredAccounts: ComputedRef<AccountInList[]> = computed(() => {
  const accountsList: AccountInList[] = [];
  Object.entries(accounts).forEach(([accessToken, accounts]) => {
    accounts.forEach((account) => {
      for (const type of Object.keys(filterOptions.accountType)) {
        if (account.subtype === type) return;
      }

      for (const [item, amount] of Object.entries(filterOptions.balance)) {
        const accountBalance = account.balances.available;
        if (item === "equalTo") {
          if (accountBalance !== amount) return;
        }
        if (item === "greaterThan") {
          if (accountBalance <= amount) return;
        }
        if (item === "lessThan") {
          if (accountBalance >= amount) return;
        }
      }

      account.access_token = accessToken;
      accountsList.push(account as AccountInList);
    });
  });
  return accountsList;
});

const accountTypes: ComputedRef<AccountTypesKeys> = computed(() => {
  const accTypes: AccountTypes = {};
  Object.values(accounts)
    .flat()
    .forEach((acc) => {
      if (!accTypes[acc.subtype]) {
        accTypes[acc.subtype] = true;
      }
    });
  return Object.keys(accTypes) as AccountTypesKeys;
});

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

async function refreshBalance(accessToken: string): Promise<void> {
  // getAccountsBalances func - partially Plaid and partially db
  const data: AccountsBalancesResponse | void = await getAccountsBalances(
    accessToken
  );

  if (data) {
    accounts[accessToken] = data.accounts;
  }
}

async function removeAccountHandler(
  accessToken: string,
  account: Account
): Promise<void> {
  try {
    await removeAccount(
      account.account_id
      // account.persistent_account_id
    );
  } catch (err) {
    console.log(
      `There was an error trying to remove the account with ID ${account.account_id}:`,
      err
    );
    return;
  }

  const accountsBelongingToAccessToken = accounts[accessToken];
  if (accountsBelongingToAccessToken.length > 1) {
    const deletionIdx = accountsBelongingToAccessToken.findIndex(
      (account) => account.account_id === account.account_id
    );
    accounts[accessToken].splice(deletionIdx, 1);
  } else {
    delete accounts[accessToken];
  }
}
</script>
