<template>
  <RouterLink to="/" class="brand" aria-label="AntriSehat beranda">
    <span class="brand-mark">
      <Activity :size="27" :stroke-width="2.6" />
    </span>
    <span>
      Antri<span class="teal">Sehat</span>
      <i />
    </span>
  </RouterLink>
</template>

<script setup lang="ts">
import { RouterLink } from "vue-router";
import { Activity } from "lucide-vue-next";
</script>
