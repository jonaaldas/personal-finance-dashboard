<script lang="ts" setup>
import LinkButton from "../composables/link_button.vue";
import { api } from "../../convex/_generated/api";

const { data: accountsData } = useConvexQuery(api.functions.queries.getAccounts, {
  user_id: 1,
});

const {
  public: { DATASYNC_URL },
} = useRuntimeConfig();
const refreshAccounts = async () => {
  await $fetch(`${DATASYNC_URL}/api/accounts`);
};
</script>
<template>
  <div>
    <UHeader title="Finance Dashboard">
      <template #right>
        <LinkButton />
        <UButton @click="refreshAccounts">Refresh Accounts</UButton>
      </template>
    </UHeader>
    <pre>{{ accountsData }}</pre>
  </div>
</template>
