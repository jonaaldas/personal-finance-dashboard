<template>
  <UButton @mouseover="getPublicToken" @click="open">Connect a bank account</UButton>
</template>

<script lang="ts" setup>
import type {
  PlaidLinkOnEvent,
  PlaidLinkOnExit,
  PlaidLinkOnSuccess,
  PlaidLinkOptions,
} from "@jcss/vue-plaid-link";
import { usePlaidLink } from "@jcss/vue-plaid-link";
import { api } from "../../convex/_generated/api";
import { v } from "convex/values";
const { mutate: saveAccessTokens, isPending: isSavingAccessTokens } = useConvexMutation(
  api.functions.mutation.saveAccessTokens
);

const {
  public: { DATASYNC_URL: DATA_ENGINE_URL },
} = useRuntimeConfig();

let linkToken: Ref<string> = ref("");
const getPublicToken = async () => {
  try {
    const data = await $fetch<{ link_token: string }>(`${DATA_ENGINE_URL}/api/create_link_token`, {
      method: "POST",
    });
    linkToken = computed(() => data.link_token);
  } catch (error) {
    console.error("Error creating link token", error);
  }
};

const onSuccess: PlaidLinkOnSuccess = async (publicToken, metadata) => {
  try {
    const res = await $fetch<{ access_token: string; item_id: string }>(
      `${DATA_ENGINE_URL}/api/set_access_token`,
      {
        method: "POST",
        body: {
          publicToken,
        },
      }
    );
    saveAccessTokens({ access_token: res.access_token, item_id: res.item_id });
  } catch (error) {
    console.error("Error setting access token", error);
  }
};

const onEvent: PlaidLinkOnEvent = (eventName, metadata) => {
  console.log(eventName, metadata);
};

const onExit: PlaidLinkOnExit = (error, metadata) => {
  linkToken.value = "";
};

const config = computed(() => {
  const config: PlaidLinkOptions = {
    token: linkToken.value,
    onSuccess,
    onEvent,
    onExit,
  };
  return config;
});

const { open, ready } = usePlaidLink(config);
</script>
