import { v } from "convex/values";
import { mutation } from "../_generated/server";

export const saveAccessTokens = mutation({
  args: { access_token: v.string(), item_id: v.string() },
  handler: async (ctx, args) => {
    return await ctx.db.insert("access_tokens", {
      user_id: 1,
      access_token: args.access_token,
      item_id: args.item_id,
    });
  },
});

export const saveAccounts = mutation({
  args: {
    user_id: v.number(),
    accounts: v.array(
      v.object({
        account_id: v.string(),
        account_name: v.string(),
        account_type: v.string(),
        account_subtype: v.string(),
        account_mask: v.string(),
        account_institution_name: v.string(),
        account_institution_id: v.string(),
        account_institution_logo: v.string(),
        account_institution_url: v.string(),
        account_available_balance: v.number(),
        account_available_current: v.number(),
        account_available_iso: v.number(),
        account_available_limit: v.optional(v.number()),
      })
    ),
  },
  handler: async (ctx, args) => {
    const uniqueAccounts = args.accounts.filter(
      (account, index, self) => index === self.findIndex((t) => t.account_id === account.account_id)
    );
    for (const account of uniqueAccounts) {
      await ctx.db.insert("accounts", {
        user_id: args.user_id,
        account_id: account.account_id,
        account_name: account.account_name,
        account_type: account.account_type,
        account_subtype: account.account_subtype,
        account_mask: account.account_mask,
        account_institution_name: account.account_institution_name,
        account_institution_id: account.account_institution_id,
        account_institution_logo: account.account_institution_logo,
        account_institution_url: account.account_institution_url,
        account_available_balance: account.account_available_balance,
        account_available_current: account.account_available_current,
        account_available_iso: account.account_available_iso,
        account_available_limit: account.account_available_limit,
      });
    }
    return {
      message: "Accounts saved successfully",
    };
  },
});
