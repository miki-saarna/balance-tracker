import { h } from "vue";
// import { usePlaidLink } from "react-plaid-link";
// import { usePlaidLink } from "plaid-link";
import { create } from "plaid";

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
  // const onSuccess = React.useCallback(async (public_token, metadata) => {
  const onSuccess =
    (async (public_token) => {
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
      const data = await response.json(); // {access_token, item_id}
      props.setAccessTokens((prevTokens) => [...prevTokens, data.access_token]);
    },
    []);

  // const config: Parameters<typeof usePlaidLink>[0] = {
  const config = {
    token: props.linkToken,
    onSuccess,
  };

  const { open, ready } = create(config);
  // const { open, ready } = usePlaidLink(config);

  const vnode = h(
    "button",
    { onClick: open, disabled: !ready },
    "Link account"
  );

  return vnode;

  // return (
  //   <button onClick={() => open()} disabled={!ready}>
  //     Link account
  //   </button>
  // );
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
