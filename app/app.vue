<script setup lang="ts">
import LinkButton from "./composables/link_button.vue";
const DATA_ENGINE_URL = "http://localhost:5050";

const getAccounts = async () => {
  const res = await $fetch(`${DATA_ENGINE_URL}/api/accounts`);
  console.log(res);
};

const getLiabilities = async () => {
  const res = await $fetch(`${DATA_ENGINE_URL}/api/liabilities`);
  console.log(res);
};
</script>

<template>
  <UPage>
    <UContainer>
      <LinkButton />
      <UButton @click="getAccounts">Get Accounts</UButton>
      <UButton @click="getLiabilities">Get Liabilities</UButton>
    </UContainer>
  </UPage>
</template>
