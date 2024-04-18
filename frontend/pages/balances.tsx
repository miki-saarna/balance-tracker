// APP COMPONENT
// Upon rendering of App component, make a request to create and
// obtain a link token to be used in the Link component
import React, { useEffect, useState } from 'react';
import { usePlaidLink } from 'react-plaid-link';
import { AccessTokensResponse, getAccessTokens } from '../utils/db';

const App = () => {
  const [linkToken, setLinkToken] = useState(null);
  const [accessToken, setAccessToken] = useState(null);
  const [accessTokens, setAccessTokens] = useState<string[]>([]);

  useEffect(() => {
    (async () => {
      const data: AccessTokensResponse | void = await getAccessTokens();
      if (data) {
        setAccessTokens(data.access_tokens)
      } else {
        console.log("There was an error retrieving access_tokens")
      }
    })()
    generateToken();
  }, []);

  const generateToken = async () => {
    const response = await fetch('http://localhost:8000/api/create_link_token', {
      method: 'POST',
    });
    const data = await response.json();
    console.log(data)
    setLinkToken(data.link_token);
  };

  return (
    linkToken != null &&
      <>
        <Link linkToken={linkToken} setAccessToken={setAccessToken} />
        {
          accessToken != null &&
          <Balance accessToken={accessToken} />
        }
        {accessTokens.map((accessToken) => <Balance key={accessToken} accessToken={accessToken} />)}
      </>
  )
};

// LINK COMPONENT
// Use Plaid Link and pass link token and onSuccess function
// in configuration to initialize Plaid Link
interface LinkProps {
  linkToken: string | null;
  setAccessToken: Function
}

interface BalanceProps {
  accessToken: string | null
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

const Balance: React.FC<BalanceProps> = (props: BalanceProps) => {
  const [accounts, setAccounts] = useState<any[]>([]);

  useEffect(() => {
    getBalances();
  }, []);

  // pass in access_token
  const getBalances = async () => {
    const response = await fetch('http://localhost:8000/api/balance', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ accessToken: props.accessToken }),
    });
    const data = await response.json();
    console.log(data)
    setAccounts(accounts.concat(data.accounts))
  }

  return (
    <>
      {accounts.map((account) => {
        return (
          <div key={account.account_id} className="line">
            <div>{account.name}</div>
            <div>{account.subtype}</div>
            <div>{account.balances.current}</div>
          </div>
        )
      })}
      <div>
        Total: {accounts.reduce((acc, account) => acc + account.balances.current, 0)}
      </div>
    </>
  )
}

export default App;