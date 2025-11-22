import { defineSchema, defineTable } from "convex/server";
import { v } from "convex/values";

export default defineSchema({
  users: defineTable({
    user_id: v.number(),
    email: v.string(),
  }).index("by_user_id", ["user_id"]),
  access_tokens: defineTable({
    user_id: v.number(),
    access_token: v.string(),
    item_id: v.string(),
  }).index("by_access_token_user_id", ["user_id", "access_token"]),
  accounts: defineTable({
    user_id: v.number(),
    account_id: v.string(),
    account_available_balance: v.number(),
    account_available_current: v.number(),
    account_available_iso: v.number(),
    account_available_limit: v.optional(v.number()),
    account_name: v.string(),
    account_type: v.string(),
    account_subtype: v.string(),
    account_mask: v.string(),
    account_institution_name: v.string(),
    account_institution_id: v.string(),
    account_institution_logo: v.string(),
    account_institution_url: v.string(),
  }).index("by_user_id", ["user_id"]),
});
