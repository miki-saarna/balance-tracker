import { h, ref, onMounted, computed } from "vue";
import type { Ref, VNode } from "vue";

export type LinkTokenResponse = {
  link_token: string;
};

interface LinkProps {
  linkToken: string | null;
  accessTokens: Ref<string[]>;
  refreshBalance: (accessToken: string) => Promise<void>;
}

export type AccountsBalancesResponse = {
  accounts: any; // check Plaid docs for actual type
};

const genLinkToken = async (): Promise<LinkTokenResponse | void> => {
  let data: LinkTokenResponse;
  try {
    const res = await fetch("http://localhost:8000/api/create_link_token", {
      method: "POST",
    });
    data = await res.json();
    return data;
  } catch (err) {
    throw err;
  }
};

const Link = (props: LinkProps): VNode => {
  // const onSuccess = React.useCallback(async (public_token, metadata) => {
  const onSuccess = async (public_token: string, metadata: any) => {
    let data;
    try {
      const response = await fetch(
        "http://localhost:8000/api/set_access_token",
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({ public_token }),
        }
      );

      data = await response.json(); // {access_token, item_id}
    } catch (err) {
      console.error("There was an error accessing an account:", err);
      return;
    }

    props.accessTokens.value.push(data.access_token);

    try {
      await props.refreshBalance(data.access_token);
    } catch (err) {
      console.error(
        `There was an error retrieving balance of account with access_token ${data.access_token}:`,
        err
      );
    }
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
    const data: AccountsBalancesResponse = await res.json();
    return data;
  } catch (err) {
    console.log(
      "There was an error retrieving accounts with associated balances",
      err
    );
  }
};

export { genLinkToken, Link, getAccountsBalances };
