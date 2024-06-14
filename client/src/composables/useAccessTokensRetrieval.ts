import { ref, Ref } from "vue";
import { AccessTokensResponse, getAccessTokens } from "../utils/db";
import { LinkTokenResponse, generateLinkToken } from "../utils/plaid_api";

export function usePlaidTokens() {
  const accessTokens: Ref<string[]> = ref([]);
  const linkToken: Ref<string> = ref("");

  const genAccessTokens = async () => {
    try {
      const data: AccessTokensResponse | void = await getAccessTokens();
      if (data?.access_tokens) {
        accessTokens.value = data.access_tokens;
      }
    } catch (err) {
      console.log("There was an error retrieving access_tokens:", err);
    }
  };

  const genLinkToken = async () => {
    try {
      const data: LinkTokenResponse | void = await generateLinkToken();
      if (data) {
        linkToken.value = data.link_token;
      }
    } catch (err) {
      console.log("There was an error retrieving the link_token:", err);
    }
  };

  return {
    accessTokens,
    genAccessTokens,
    linkToken,
    genLinkToken,
  };
}
