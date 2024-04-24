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

type AccountsByAccessToken = {
  [key: string]: any[] // update with correct type from Plaid
}

const Balances = () => {
  const [linkToken, setLinkToken] = useState("");
  const [accessTokens, setAccessTokens] = useState<string[]>([]);
  const [accounts, setAccounts] = useState<AccountsByAccessToken>({});

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

  useEffect(() => {
    (async() => {
      for (const accessToken of accessTokens) {
        await refreshBalance(accessToken)
      }
    })()
  }, [accessTokens])

  async function refreshBalance(accessToken: string): Promise<void> {
    // getAccountsBalances func - partially Plaid and partially db
    const data: AccountsBalancesResponse | void = await getAccountsBalances(accessToken)
    if (data) {
      setAccounts((prevAccounts) => ({
        ...prevAccounts,
        [accessToken]: data.accounts
      }))
    }
  }

  async function removeAccount(accessToken: string, account_id: string, p): Promise<void> { // 2nd parameter necessary?
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

    const accountsBelongingToAccessToken = modifyAccounts[accessToken]
    if (accountsBelongingToAccessToken.length > 1) {
      const deletionIdx = accountsBelongingToAccessToken.findIndex((account) => account.account_id === account_id)
      modifyAccounts[accessToken].splice(deletionIdx, 1)
      setAccounts(modifyAccounts)
    } else {
      delete modifyAccounts[accessToken]
      setAccounts(modifyAccounts)
    }
  }

  function RenderTotalBalance() {
    const accountsList = Object.values(accounts)
    if (!accountsList.length) return
    const sum = Object.values(accounts).flat().reduce((acc, account) => acc + account.balances.current, 0)
    return <div>Total: ${sum}</div>
  }

  return (
    <>
      <Link linkToken={linkToken} setAccessTokens={setAccessTokens} />

      {Object.entries(accounts).map(([accessToken, accountsList]) => 
        accountsList.map((account) => {
          return (
            <div key={account.account_id} className="line">
              <div>{account.name}</div>
              <div>{account.subtype}</div>
              <div>${account.balances.current}</div>
              <button onClick={async () => await refreshBalance(accessToken)}>refresh</button>
              {/* currently not saving `persistent_account_id` within the db */}
              <button onClick={() => removeAccount(accessToken, account.account_id, account.persistent_account_id)}>remove</button>
            </div>
          )
        })
      )}

      <RenderTotalBalance />
    </>
  )
};

export default Balances;
