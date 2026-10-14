<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { Activity, ChevronRight, Volume2, Clock3 } from 'lucide-vue-next'
import { useQueueStore, today } from '../stores/queue'

const router = useRouter()
const queueStore = useQueueStore()

const active = computed(() => queueStore.activePatient)

const nextPatients = computed(() =>
  queueStore.patients.filter((p) => p.date === today && p.status === 'Hadir').slice(0, 6),
)

const repeatCall = () => {
  if (active.value) {
    queueStore.callPatient(active.value)
  }
}
</script>

<template>
  <div class="tv-app">
    <main class="tv-screen">
      <div class="tv-header">
        <button class="brand" @click="router.push('/')" aria-label="AntriSehat beranda">
          <span class="brand-mark">
            <Activity :size="27" :stroke-width="2.6" />
          </span>
          <span>
            Antri<span class="teal">Sehat</span>
            <i />
          </span>
        </button>

        <button class="button outline small" @click="router.push('/staff')">
          <ChevronRight :size="18" class="rotate-180" /> Kembali ke Dashboard Staf
        </button>
      </div>

      <div class="tv-grid">
        <div class="tv-current">
          <span>SEDANG DIPANGGIL</span>
          <h1>{{ active?.number || '—' }}</h1>
          <h2>{{ active?.poli || 'Menunggu panggilan' }}</h2>
          <p>Silakan menuju ruang pemeriksaan</p>
          <button class="button white" :disabled="!active" @click="repeatCall">
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
        <Clock3 :size="20" /> Jam pelayanan 08.00 – 14.00
        <span> Siapkan KTP dan kartu BPJS Anda · Jaga kebersihan dan kenyamanan bersama </span>
      </div>
    </main>
  </div>
</template>
