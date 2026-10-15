<template>
  <div :class="{ 'tv-app': isTvRoute }">
    <AppHeader v-if="showHeader" />
    <RouterView />
    <ToastNotification />
  </div>
</template>

<script setup lang="ts">
import { computed } from "vue";
import { useRoute } from "vue-router";
import AppHeader from "./components/layout/AppHeader.vue";
import ToastNotification from "./components/common/ToastNotification.vue";

const route = useRoute();

const isTvRoute = computed(() => route.path === "/tv");

const showHeader = computed(() => {
  const hiddenRoutes = ["/tv", "/staff", "/medical-record", "/scan", "/qr-test"];
  return !hiddenRoutes.includes(route.path);
});
</script>
