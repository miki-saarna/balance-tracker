type AccessTokensResponse = {
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
  } catch (err) {
    console.error(err);
    return;
  }

  const data: AccessTokensResponse = await res.json();
  return data;
};

export { AccessTokensResponse, getAccessTokens };
