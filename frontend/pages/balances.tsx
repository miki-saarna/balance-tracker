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
        const data: AccessTokensResponse | void = await getAccessTokens(); // void might be unused
        if (data?.access_tokens) {
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
    // getAccountsBalances func - partially Plaid and partially db
    const data: AccountsBalancesResponse | void = await getAccountsBalances(accessToken)
      if (data) {
        setAccounts({
          ...accounts,
          [props.accessToken]: data.accounts
        })
      }
  }

  async function removeAccount(account_id: string, p): Promise<void> { // 2nd parameter necessary?
    // console.log("p", p)

    let res
    try {
      res = await fetch("http://localhost:8000/api/account/delete", {
        method: "DELETE",
        headers: {
          "Content-Type": "application/json"
        },
        body: JSON.stringify({ account_id })
      })
    } catch (err) {
      console.log("There was an error deleting the account", err)
      return
    }

    const data = await res.json()
    // console.log("data: ", data)

    const modifyAccounts = { ...accounts }

    const accountsBelongingToAccessToken = modifyAccounts[props.accessToken]
    if (accountsBelongingToAccessToken.length > 1) {
      const deletionIdx = accountsBelongingToAccessToken.findIndex((account) => account.account_id === account_id)
      modifyAccounts[props.accessToken].splice(deletionIdx, 1)
      setAccounts(modifyAccounts)
    } else {
      delete modifyAccounts[props.accessToken]
      setAccounts(modifyAccounts)
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
            {/* currently not saving `persistent_account_id` within the db */}
            <button onClick={() => removeAccount(account.account_id, account.persistent_account_id)}>remove</button>
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
