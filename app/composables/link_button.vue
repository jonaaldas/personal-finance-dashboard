<template>
  <UButton :disabled="!ready" @click="open">Connect a bank account</UButton>
</template>

<script lang="ts" setup>
import type {
  PlaidLinkOnEvent,
  PlaidLinkOnExit,
  PlaidLinkOnSuccess,
  PlaidLinkOptions,
} from "@jcss/vue-plaid-link";
import { usePlaidLink } from "@jcss/vue-plaid-link";
const DATA_ENGINE_URL = "http://localhost:5050";
let linkToken: Ref<string> = ref("");
try {
  const { data } = useFetch<{ link_token: string }>(
    `${DATA_ENGINE_URL}/api/create_link_token`,
    {
      method: "POST",
    }
  );
  linkToken = computed(() => {
    console.log(linkToken);
    if (!data.value) {
      return "";
    }
    return data.value.link_token;
  });
} catch (error) {
  console.error("Error creating link token", error);
}

const onSuccess: PlaidLinkOnSuccess = async (publicToken, metadata) => {
  try {
    const res = await $fetch(`${DATA_ENGINE_URL}/api/set_access_token`, {
      method: "POST",
      body: {
        publicToken,
      },
    });
    console.log("this is the token", res);
  } catch (error) {
    console.error("Error setting access token", error);
  }
};

const onEvent: PlaidLinkOnEvent = (eventName, metadata) => {
  // log onEvent callbacks from Link
  // https://plaid.com/docs/link/web/#onevent
  console.log(eventName, metadata);
};

const onExit: PlaidLinkOnExit = (error, metadata) => {
  // log onExit callbacks from Link, handle errors
  // https://plaid.com/docs/link/web/#onexit
  console.log(error, metadata);
};

const config = computed(() => {
  console.log("this is the token", linkToken.value);
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
