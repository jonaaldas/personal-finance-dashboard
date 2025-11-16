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
