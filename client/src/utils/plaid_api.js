import { h, onMounted, computed } from "vue";

// type LinkTokenResponse = {
//   link_token: string;
// };

// const generateLinkToken = async (): Promise<LinkTokenResponse | void> => {
const generateLinkToken = async () => {
  const response = await fetch("http://localhost:8000/api/create_link_token", {
    method: "POST",
  });
  // const data: LinkTokenResponse = await response.json();
  const data = await response.json();
  return data;
};

// interface LinkProps {
//   linkToken: string | null;
//   setAccessTokens: Function;
// }

// const Link: React.FC<LinkProps> = (props: LinkProps) => {
const Link = (props) => {
  // const linkButton = ref(null);
  // let plaidLinkHandler = null;

  // const onSuccess = React.useCallback(async (public_token, metadata) => {
  const onSuccess = async (public_token, metadata) => {
    const response = await fetch("http://localhost:8000/api/set_access_token", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ public_token }),
    });
    const data = await response.json(); // {access_token, item_id}
    props.accessTokens.value = [...props.accessTokens.value, data.access_token];
  };

  // const config: Parameters<typeof usePlaidLink>[0] = {
  const config = {
    token: props.linkToken,
    onSuccess,
  };

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

// type AccountsBalancesResponse = {
//   accounts: any; // check Plaid docs for actual type
// };

const getAccountsBalances = async (accessToken) => {
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
  const data = await res.json();
  return data;
};

export {
  // LinkTokenResponse,
  generateLinkToken,
  Link,
  // AccountsBalancesResponse,
  getAccountsBalances,
};
