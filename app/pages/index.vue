<script lang="ts" setup>
import LinkButton from "../composables/link_button.vue";
import { api } from "../../convex/_generated/api";

const {
  public: { DATASYNC_URL },
} = useRuntimeConfig();

const getAccounts = async () => {
  const res = await $fetch(`${DATASYNC_URL}/api/accounts`);
  console.log(res);
};

const getLiabilities = async () => {
  const res = await $fetch(`${DATASYNC_URL}/api/liabilities`);
  console.log(res);
};

const getTransactions = async () => {
  const res = await $fetch(`${DATASYNC_URL}/api/transactions`);
  console.log(res);
};

const getInvestments = async () => {
  const res = await $fetch(`${DATASYNC_URL}/api/investments`);
  console.log(res);
};

const { data: accountsData } = useConvexQuery(api.functions.queries.getAccounts, {
  user_id: 1,
});
const refreshAccounts = async () => {
  await $fetch(`${DATASYNC_URL}/api/accounts`);
};
</script>
<template>
  <div>
    <UHeader title="Finance Dashboard">
      <template #right>
        <LinkButton />
        <UButton @click="getAccounts">Get Accounts</UButton>
        <UButton @click="getLiabilities">Get Liabilities</UButton>
        <UButton @click="getTransactions">Get Transactions</UButton>
        <UButton @click="getInvestments">Get Investments</UButton>
        <UButton @click="refreshAccounts">Refresh Accounts</UButton>
      </template>
    </UHeader>
    <pre>{{ accountsData }}</pre>
  </div>
</template>
