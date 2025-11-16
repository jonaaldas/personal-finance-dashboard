import tailwindcss from "@tailwindcss/vite";
import { env } from "./env";
console.log(env);
export default defineNuxtConfig({
  compatibilityDate: "2025-07-15",
  devtools: { enabled: true },
  modules: ["@nuxt/image", "@nuxt/ui", "convex-nuxt"],
  vite: {
    plugins: [tailwindcss()],
  },
  css: ["./app/styles/index.css"],
  devServer: {
    port: 9250,
  },
  convex: {
    url: env.CONVEX_URL,
  },
  runtimeConfig: {
    public: {
      DATASYNC_URL: env.DATASYNC_URL,
    },
  },
});
