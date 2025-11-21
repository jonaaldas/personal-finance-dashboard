import { httpRouter } from "convex/server";
import { httpAction } from "./_generated/server";
import { api } from "./_generated/api";
import type { PlaidAccountsGetResponse, SaveAccountsResponse } from "./types";

const http = httpRouter();

http.route({
  path: "/api/save",
  method: "POST",
  handler: httpAction(async (ctx, request) => {
    const body = (await request.json()) as PlaidAccountsGetResponse;
    await ctx.runMutation(api.functions.mutation.saveAccounts, {
      user_id: 1,
      accounts: body.accounts.map((account) => ({
        account_id: account.account_id,
        account_available_balance: account.balances.available ?? 0,
        account_available_current: account.balances.current,
        account_available_iso: parseFloat(account.balances.iso_currency_code || "0") || 0,
        account_available_limit: account.balances.limit ?? undefined,
        account_name: account.name,
        account_type: account.type,
        account_subtype: account.subtype,
        account_mask: account.mask,
        account_institution_name: body.item.institution_name,
        account_institution_id: body.item.institution_id,
        account_institution_logo: body.item.institution_logo ?? "",
        account_institution_url: body.item.institution_url ?? "",
      })),
    });

    const response: SaveAccountsResponse = {
      message: "Account saved successfully",
    };
    return new Response(JSON.stringify(response));
  }),
});

http.route({
  path: "/api/get",
  method: "GET",
  handler: httpAction(async (ctx, request) => {
    return new Response(JSON.stringify({ message: "Hello, world!" }));
  }),
});

http.route({
  path: "/api/access_tokens",
  method: "GET",
  handler: httpAction(async (ctx, request) => {
    let userId = 1; // this will be from JWT or something else
    const access_tokens = await ctx.runQuery(api.functions.queries.getAllTokens, {
      user_id: userId,
    });
    return new Response(JSON.stringify(access_tokens));
  }),
});
export default http;
