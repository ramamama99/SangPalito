<template>
  <main class="tv-screen">
    <div class="tv-header">
      <Brand />
      <RouterLink to="/staff" class="button outline small">
        <ChevronRight :size="18" class="rotate-180" /> Kembali ke Dashboard Staf
      </RouterLink>
    </div>
    <div class="tv-grid">
      <div class="tv-current">
        <span>SEDANG DIPANGGIL</span>
        <h1>{{ queueStore.active?.number || "—" }}</h1>
        <h2>{{ queueStore.active?.poli || "Menunggu panggilan" }}</h2>
        <p>Silakan menuju ruang pemeriksaan</p>
        <button
          class="button white"
          :disabled="!queueStore.active"
          @click="repeatCall"
        >
          <Volume2 :size="23" /> Ulangi Panggilan
        </button>
      </div>
      <div class="tv-next">
        <h2>Antrean berikutnya</h2>
        <div v-for="p in nextPatients" :key="p.id">
          <strong>{{ p.number }}</strong>
          <span>{{ p.poli }}</span>
        </div>
        <p v-if="!nextPatients.length">Belum ada antrean berikutnya.</p>
      </div>
    </div>
    <div class="tv-footer">
      <Clock3 :size="20" /> Jam pelayanan 08.00 – 14.00{" "}
      <span>
        Siapkan KTP dan kartu BPJS Anda · Jaga kebersihan dan kenyamanan bersama
      </span>
    </div>
  </main>
</template>

<script setup lang="ts">
import { computed } from "vue";
import { useRouter } from "vue-router";
import { ChevronRight, Volume2, Clock3 } from "lucide-vue-next";
import Brand from "../../components/common/Brand.vue";
import { useQueueStore } from "../../stores/queue";
import { todayDate } from "../../constants/initialData";

const router = useRouter();
const queueStore = useQueueStore();

const nextPatients = computed(() =>
  queueStore.patients
    .filter((p) => p.date === todayDate && p.status === "Hadir")
    .slice(0, 6),
);

function repeatCall() {
  if (queueStore.active) {
    queueStore.call(queueStore.active);
  }
}

function goToStaff() {
  router.push("/staff");
  window.scrollTo(0, 0);
}
</script>
