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

export {
  LinkTokenResponse,
  generateLinkToken
}
