// APP COMPONENT
// Upon rendering of App component, make a request to create and
// obtain a link token to be used in the Link component
import React, { useEffect, useState } from 'react';
import { AccessTokensResponse, getAccessTokens } from '../utils/db';
import {
  LinkTokenResponse,
  generateLinkToken,
  Link,
  AccountsBalancesResponse,
  getAccountsBalances
} from '../utils/plaid_api';

const Balances = () => {
  const [linkToken, setLinkToken] = useState("");
  const [accessTokens, setAccessTokens] = useState<string[]>([]);

  useEffect(() => {
    (async () => {
      try {
        const data: AccessTokensResponse | void = await getAccessTokens();
        if (data) {
          setAccessTokens(data.access_tokens)
        }
      } catch (err) {
        console.log("There was an error retrieving access_tokens:", err)
      }

      try {
        const data: LinkTokenResponse | void = await generateLinkToken();
        if (data) {
          setLinkToken(data.link_token)
        }
      } catch (err) {
        console.log("There was an error retrieving the link_token:", err)
      }
    })()
  }, []);

  return (
    linkToken != null &&
      <>
        <Link linkToken={linkToken} setAccessTokens={setAccessTokens} />
        {accessTokens.map((accessToken) => <Balance key={accessToken} accessToken={accessToken} />)}
      </>
  )
};

interface BalanceProps {
  accessToken: string
}

type AccountsByAccessToken = {
  [key: string]: any[] // update with correct type from Plaid
}

const Balance: React.FC<BalanceProps> = React.memo((props: BalanceProps) => {
  const [accounts, setAccounts] = useState<AccountsByAccessToken>({});

  useEffect(() => {
    (async() => {
      await refreshBalance(props.accessToken)
    })()
  }, []);

  async function refreshBalance(accessToken: string): Promise<void> {
    const data: AccountsBalancesResponse | void = await getAccountsBalances(accessToken)
      if (data) {
        setAccounts({
          ...accounts,
          [props.accessToken]: data.accounts
        })
      }
  }

  return (
    <>
      {Object.values(accounts).flat().map((account) => {
        return (
          <div key={account.account_id} className="line">
            <div>{account.name}</div>
            <div>{account.subtype}</div>
            <div>{account.balances.current}</div>
            <button onClick={async () => await refreshBalance(props.accessToken)}>refresh</button>
          </div>
        )
      })}
      <div>
        Total: {Object.values(accounts).flat().reduce((acc, account) => acc + account.balances.current, 0)}
      </div>
    </>
  )
})

export default Balances;
