import React from 'react'
import { usePlaidLink } from 'react-plaid-link'

type LinkTokenResponse = {
  "link_token": string
}

const generateLinkToken = async (): Promise<LinkTokenResponse | void> => {
  const response = await fetch('http://localhost:8000/api/create_link_token', {
    method: 'POST',
  });
  const data: LinkTokenResponse = await response.json();
  return data
};

interface LinkProps {
  linkToken: string | null;
  setAccessToken: Function
}

const Link: React.FC<LinkProps> = (props: LinkProps) => {
  const onSuccess = React.useCallback(async (public_token, metadata) => {
    const response = await fetch('http://localhost:8000/api/set_access_token', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ public_token }),
    });
    const data = await response.json();
    console.log(data) // {access_token, item_id}
    props.setAccessToken(data.access_token)
  }, []);

  const config: Parameters<typeof usePlaidLink>[0] = {
    token: props.linkToken!,
    onSuccess,
  };

  const { open, ready } = usePlaidLink(config);
  
  return (
    <button onClick={() => open()} disabled={!ready}>
      Link account
    </button>
  )
}

type AccountsBalancesResponse = {
  "accounts": any // check Plaid docs for actual type
}

const getAccountsBalances = async (accessToken: string): Promise<AccountsBalancesResponse | void> => {
  let res
  try {
    res = await fetch('http://localhost:8000/api/balance', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ access_token: accessToken }),
    });
  } catch (err) {
    console.log("There was an error retrieving accounts with associated balances", err)
    return
  }
  const data: AccountsBalancesResponse = await res.json();
  return data
}

export {
  LinkTokenResponse,
  generateLinkToken,
  Link,
  AccountsBalancesResponse,
  getAccountsBalances
}
