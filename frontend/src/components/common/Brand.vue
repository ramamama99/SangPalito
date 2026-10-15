<template>
  <button class="brand" @click="goToHome" aria-label="AntriSehat beranda">
    <span class="brand-mark">
      <Activity :size="27" :stroke-width="2.6" />
    </span>
    <span>
      Antri<span class="teal">Sehat</span>
      <i />
    </span>
  </button>
</template>

<script setup lang="ts">
import { useRouter } from "vue-router";
import { Activity } from "lucide-vue-next";

const router = useRouter();

function goToHome() {
  router.push("/");
  window.scrollTo(0, 0);
}
</script>
