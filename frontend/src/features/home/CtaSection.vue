<template>
  <section class="cta container">
    <div>
      <div class="section-label">KESEHATAN TIDAK PERLU MENUNGGU</div>
      <h2>Mulai langkah sehat Anda hari ini.</h2>
      <p>Ambil antrean dari rumah. Kami menantikan Anda.</p>
    </div>
    <RouterLink
      to="/register"
      class="button white"
      @click="queueStore.startRegistration('Poli Umum', false)"
    >
      Daftar Sekarang <ArrowUpRight :size="20" />
    </RouterLink>
  </section>
</template>

<script setup lang="ts">
import { useRouter } from "vue-router";
import { ArrowUpRight } from "lucide-vue-next";
import { useQueueStore } from "../../stores/queue";

const router = useRouter();
const queueStore = useQueueStore();

function handleStart() {
  queueStore.startRegistration("Poli Umum", false);
  router.push("/register");
  window.scrollTo(0, 0);
}
</script>
