import { defineSchema, defineTable } from "convex/server";
import { v } from "convex/values";

export default defineSchema({
  access_tokens: defineTable({
    user_id: v.number(),
    access_token: v.string(),
    item_id: v.string(),
  }).index("by_access_token", ["access_token"]),
});
