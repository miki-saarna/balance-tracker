import { h, onMounted, computed } from "vue";
import type { Ref, VNode } from "vue";

export type LinkTokenResponse = {
  link_token: string;
};

interface LinkProps {
  linkToken: string | null;
  accessTokens: Ref<string[]>;
}

export type AccountsBalancesResponse = {
  accounts: any; // check Plaid docs for actual type
};

const generateLinkToken = async (): Promise<LinkTokenResponse | void> => {
  let data: LinkTokenResponse;
  try {
    const response = await fetch(
      "http://localhost:8000/api/create_link_token",
      {
        method: "POST",
      }
    );
    data = await response.json();
  } catch (err) {
    console.error("Error:", err);
    return;
  }
  return data;
};

const Link = (props: LinkProps): VNode => {
  // const onSuccess = React.useCallback(async (public_token, metadata) => {
  const onSuccess = async (public_token: string, metadata: any) => {
    const response = await fetch("http://localhost:8000/api/set_access_token", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ public_token }),
    });
    const data = await response.json(); // {access_token, item_id}
    props.accessTokens.value.push(data.access_token);
  };

  // @ts-ignore
  const config: Parameters<typeof Plaid.create>[0] = {
    token: props.linkToken,
    onSuccess,
  };

  // @ts-ignore
  const { open } = Plaid.create(config);

  const openLink = async () => {
    try {
      const res = await open();
    } catch (err) {
      console.error("Error:", err);
    }
  };

  return h(
    "button",
    {
      class:
        "py-1 px-2 rounded bg-green-300 disabled:bg-green-50 cursor-pointer",
      onClick: openLink,
    },
    "Connect a new account"
  );
};

const getAccountsBalances = async (
  accessToken: string
): Promise<AccountsBalancesResponse | void> => {
  let res;
  try {
    res = await fetch("http://localhost:8000/api/balance", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ access_token: accessToken }),
    });
  } catch (err) {
    console.log(
      "There was an error retrieving accounts with associated balances",
      err
    );
    return;
  }
  const data: AccountsBalancesResponse = await res.json();
  return data;
};

export { generateLinkToken, Link, getAccountsBalances };
