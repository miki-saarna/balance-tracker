export type AccessTokensResponse = {
  access_tokens: string[];
};

const getAccessTokens = async (): Promise<AccessTokensResponse | void> => {
  let res;
  try {
    res = await fetch("http://localhost:8000/api/get_access_tokens", {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
      },
      body: null,
    });
    const data: AccessTokensResponse = await res.json();
    return data;
  } catch (err) {
    console.log("There was an error retrieving access_tokens:", err);
  }
};

async function removeAccount(
  account_id: string
  // persistent_id: string
): Promise<void> {
  // 2nd parameter necessary?

  let res;
  try {
    res = await fetch("http://localhost:8000/api/account/delete", {
      method: "DELETE",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ account_id }),
    });
  } catch (err) {
    console.log("There was an error deleting the account", err);
    return;
  }

  const data = await res.json();
  // console.log("data: ", data);
}

export { getAccessTokens, removeAccount };
