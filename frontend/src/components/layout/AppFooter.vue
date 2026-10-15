<template>
  <footer class="container">
    <Brand />
    <p>Puskesmas Sehat Bersama · Sistem Antrean Digital Terpadu</p>
    <span
      >© {{ new Date().getFullYear() }} AntriSehat. Hak Cipta Dilindungi.</span
    >
  </footer>
</template>

<script setup lang="ts">
import Brand from "../common/Brand.vue";
</script>
