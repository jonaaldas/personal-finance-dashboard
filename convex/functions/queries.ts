import { query } from "../_generated/server";
import { v } from "convex/values";

export const getAccounts = query({
  args: {
    user_id: v.number(),
  },
  handler: async (ctx, args) => {
    return await ctx.db
      .query("accounts")
      .withIndex("by_user_id", (q) => q.eq("user_id", args.user_id))
      .collect();
  },
});

export const getAllTokens = query({
  args: {
    user_id: v.number(),
  },
  handler: async (ctx, args) => {
    return await ctx.db
      .query("access_tokens")
      .withIndex("by_access_token_user_id", (q) => q.eq("user_id", args.user_id))
      .collect();
  },
});
