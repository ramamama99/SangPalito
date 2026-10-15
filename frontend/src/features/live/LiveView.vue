<template>
  <main class="container py-8 sm:py-12">
    <div class="mx-auto max-w-5xl">
      <button
        class="back-link !inline-flex items-center gap-2 !text-[var(--deep)]"
        @click="goToHome"
      >
        <ChevronRight :size="18" class="rotate-180" aria-hidden="true" />
        Kembali ke beranda
      </button>
      <div class="mb-6 flex flex-wrap items-center justify-between gap-3">
        <div>
          <h2 class="!text-3xl !tracking-tight sm:!text-4xl">Antrean Live</h2>
        </div>
      </div>
      <div class="grid items-start gap-5 lg:grid-cols-2">
        <section
          class="dashboard-panel !overflow-hidden !p-0 text-center"
          aria-label="Antrean sedang dilayani"
          aria-live="polite"
        >
          <div
            class="flex items-center justify-between gap-3 border-b border-[var(--line)] bg-[var(--deep)] px-6 py-5 text-white"
          >
            <span class="inline-flex items-center gap-2 font-semibold">
              <Stethoscope :size="20" aria-hidden="true" /> Sedang dilayani
            </span>
          </div>
          <div v-if="calledPatients.length" class="grid gap-5 p-6 sm:p-8">
            <div
              v-for="p in calledPatients"
              :key="p.id"
              class="rounded-2xl bg-[var(--deep)]/5 px-4 py-8"
            >
              <span class="text-sm text-[var(--muted)]">Nomor antrean</span>
              <strong
                class="my-4 block text-6xl font-bold tracking-tight text-[var(--deep)] sm:text-7xl"
              >
                {{ p.number }}
              </strong>
              <span
                class="inline-flex items-center gap-2 rounded-full bg-white px-4 py-2 font-medium text-[var(--deep)]"
              >
                <MapPin :size="17" aria-hidden="true" />{{ p.poli }}
              </span>
            </div>
          </div>
          <div v-else class="px-6 py-12">
            <Clock3
              :size="32"
              class="mx-auto mb-3 text-[var(--deep)]"
              aria-hidden="true"
            />
            <h3>Belum ada antrean yang dipanggil</h3>
            <p>Nomor antrean akan tampil saat pelayanan dimulai.</p>
          </div>
          <p class="border-t border-[var(--line)] px-5 py-4 text-sm">
            Silakan menuju poli saat nomor Anda dipanggil.
          </p>
        </section>

        <section class="dashboard-panel" aria-label="Antrean berikutnya">
          <div class="mb-3 flex items-center justify-between gap-3">
            <h3 class="inline-flex items-center gap-2 text-lg">
              <Users
                :size="20"
                class="text-[var(--deep)]"
                aria-hidden="true"
              />Antrean berikutnya
            </h3>
            <span class="badge">{{ nextPatients.length }} menunggu</span>
          </div>
          <div
            class="flex justify-between gap-4 border-b border-[var(--line)] pb-3 pt-2 text-xs font-medium uppercase tracking-wider text-[var(--muted)]"
          >
            <span>Nomor antrean</span>
            <span>Poli tujuan</span>
          </div>
          <div
            v-for="p in nextPatients.slice(0, 6)"
            :key="p.id"
            class="flex items-center justify-between gap-4 border-b border-[var(--line)] py-4 last:border-0"
          >
            <strong
              class="rounded-lg bg-[var(--deep)]/5 px-3 py-2 text-xl tabular-nums text-[var(--deep)]"
            >
              {{ p.number }}
            </strong>
            <span class="text-right text-sm sm:text-base">{{ p.poli }}</span>
          </div>
          <p v-if="!nextPatients.length">Belum ada antrean berikutnya.</p>
        </section>
      </div>
      <p
        class="mt-6 flex items-center justify-center gap-2 rounded-xl border border-[var(--line)] px-4 py-4 text-sm"
      >
        <Clock3 :size="17" class="shrink-0" aria-hidden="true" /> Jam pelayanan
        08.00 – 14.00
      </p>
    </div>
  </main>
</template>

<script setup lang="ts">
import { computed } from "vue";
import { useRouter } from "vue-router";
import {
  ChevronRight,
  Stethoscope,
  MapPin,
  Clock3,
  Users,
} from "lucide-vue-next";
import { useQueueStore } from "../../stores/queue";
import { todayDate } from "../../constants/initialData";

const router = useRouter();
const queueStore = useQueueStore();

const calledPatients = computed(() =>
  queueStore.patients.filter(
    (p) => p.date === todayDate && p.status === "Dipanggil",
  ),
);

const nextPatients = computed(() =>
  queueStore.patients.filter(
    (p) => p.date === todayDate && p.status === "Hadir",
  ),
);

function goToHome() {
  router.push("/");
  window.scrollTo(0, 0);
}
</script>
