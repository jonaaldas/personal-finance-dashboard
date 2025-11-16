<script setup lang="ts">
import LinkButton from "./composables/link_button.vue";
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
</script>

<template>
  <UPage>
    <UContainer>
      <NuxtPage />
    </UContainer>
  </UPage>
</template>
