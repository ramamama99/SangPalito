<template>
  <main class="container py-8 sm:py-12 !max-w-[1400px]">
    <div class="mx-auto">
      <!-- Back Link -->
      <RouterLink
        to="/"
        class="back-link !inline-flex items-center gap-2 !text-[var(--deep)]"
      >
        <ChevronRight :size="18" class="rotate-180" aria-hidden="true" />
        Kembali ke beranda
      </RouterLink>

      <!-- Page Header -->
      <div class="mb-6 flex flex-wrap items-center justify-between gap-4">
        <div>
          <h2 class="!text-3xl !tracking-tight sm:!text-4xl text-[var(--deep)] font-bold">
            Antrean Live
          </h2>
          <p class="text-sm text-[var(--muted)] mt-1">
            Status panggilan aktif dan antrean berikutnya pada 4 Poliklinik Puskesmas Sehat Bersama.
          </p>
        </div>

        <!-- Clock, Status & Actions -->
        <div class="flex flex-wrap items-center gap-2.5">
          <!-- Digital Clock -->
          <div class="rounded-xl border border-[var(--line)] bg-white px-3.5 py-1.5 text-right shadow-2xs">
            <div class="font-mono text-sm font-bold text-[var(--deep)]">
              {{ currentTimeString || '00:00:00 WIB' }}
            </div>
            <div class="text-[10px] font-medium text-[var(--muted)] uppercase tracking-wider">
              {{ currentDateString || 'Hari ini' }}
            </div>
          </div>

          <!-- Live Indicator Badge -->
          <div class="inline-flex items-center gap-1.5 rounded-xl border border-[var(--line)] bg-[var(--deep)]/5 px-3 py-2 text-xs font-semibold text-[var(--deep)]">
            <span class="relative flex h-2 w-2">
              <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75"></span>
              <span class="relative inline-flex rounded-full h-2 w-2 bg-emerald-500"></span>
            </span>
            Live Real-Time
          </div>

          <!-- Refresh Button -->
          <button
            type="button"
            class="button outline small !inline-flex items-center gap-1.5 !py-2 !px-3"
            @click="refreshData"
            :disabled="refreshing"
            title="Segarkan data antrean"
          >
            <RefreshCw :size="14" :class="{ 'animate-spin': refreshing }" />
            <span class="hidden sm:inline">Segarkan</span>
          </button>
        </div>
      </div>

      <!-- Quick Ticket Finder (Lacak Nomor Saya) -->
      <div class="mb-6 rounded-2xl border border-[var(--line)] bg-white p-4 shadow-2xs flex flex-col sm:flex-row sm:items-center justify-between gap-3">
        <div class="flex items-center gap-2 text-xs text-[var(--muted)]">
          <Search :size="16" class="text-[var(--deep)] shrink-0" />
          <span>Ketik nomor tiket Anda untuk memantau posisi antrean:</span>
        </div>
        <div class="relative w-full sm:w-64">
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Contoh: A-002, B-001..."
            class="w-full rounded-xl border border-[var(--line)] bg-[var(--warm)]/30 py-2 pl-3 pr-8 text-xs font-mono font-bold text-[var(--deep)] uppercase placeholder:font-sans placeholder:font-normal focus:outline-none focus:ring-2 focus:ring-[var(--deep)]"
          />
          <button
            v-if="searchQuery"
            type="button"
            class="absolute right-2.5 top-1/2 -translate-y-1/2 text-slate-400 hover:text-slate-600"
            @click="searchQuery = ''"
          >
            <X :size="14" />
          </button>
        </div>
      </div>

      <!-- Ticket Match Notice if searched -->
      <div
        v-if="searchedResult"
        class="mb-6 rounded-2xl border border-emerald-300 bg-emerald-50/70 p-4 text-xs text-slate-800 flex items-center justify-between gap-3 shadow-2xs"
      >
        <div class="flex items-center gap-2.5">
          <CheckCircle2 :size="18" class="text-emerald-700 shrink-0" />
          <div>
            <strong>Tiket {{ searchedResult.patient.number }}</strong> terdaftar di
            <span class="font-bold text-emerald-800">{{ searchedResult.poliName }}</span>
            (Atas nama: {{ searchedResult.patient.name }})
          </div>
        </div>
        <span class="rounded-full bg-white px-3 py-1 font-mono font-bold text-emerald-800 border border-emerald-200 shadow-2xs">
          {{ searchedResult.positionText }}
        </span>
      </div>

      <!-- 4 COLUMNS SIDE BY SIDE (4 KESAMPING) -->
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-5 items-stretch">
        <section
          v-for="poli in polyclinics"
          :key="poli.id"
          class="dashboard-panel !overflow-hidden !p-0 flex flex-col justify-between transition-all duration-200"
          :class="{ '!border-[var(--deep)] !shadow-md !ring-2 !ring-[var(--deep)]/30': isPoliHighlighted(poli.code) }"
          :aria-label="`Antrean ${poli.name}`"
        >
          <!-- Column Header with Poli Title & Badge -->
          <div class="flex items-center justify-between gap-2 border-b border-[var(--line)] bg-[var(--deep)] px-5 py-4 text-white">
            <span class="inline-flex items-center gap-2 font-semibold text-sm">
              <component :is="poli.icon" :size="18" aria-hidden="true" />
              {{ poli.name }}
            </span>
            <span class="rounded-md bg-white/20 px-2 py-0.5 text-xs font-bold font-mono tracking-wider">
              {{ poli.code }}
            </span>
          </div>

          <!-- Main Body of the Column -->
          <div class="p-5 flex-1 flex flex-col justify-between space-y-5">
            <!-- 1. Sedang Dilayani Box (Center Stage) -->
            <div class="rounded-2xl bg-[var(--deep)]/5 px-3 py-7 text-center">
              <span class="text-xs uppercase font-medium tracking-wider text-[var(--muted)]">
                Sedang dilayani
              </span>

              <strong class="my-2 block text-5xl sm:text-6xl font-bold tracking-tight text-[var(--deep)] tabular-nums">
                {{ poli.currentPatient?.number || "—" }}
              </strong>

              <div v-if="poli.currentPatient" class="mt-2">
                <span class="inline-flex items-center gap-1.5 rounded-full bg-white px-3 py-1 text-xs font-medium text-[var(--deep)] shadow-xs">
                  <MapPin :size="13" aria-hidden="true" /> {{ poli.currentPatient.name }}
                </span>
              </div>
              <p v-else class="text-xs text-[var(--muted)] mt-1">
                Belum ada antrean dipanggil
              </p>
            </div>

            <!-- 2. Antrean Berikutnya List -->
            <div class="space-y-2 flex-1">
              <div class="flex items-center justify-between text-xs font-semibold text-[var(--deep)] border-b border-[var(--line)] pb-2">
                <span class="inline-flex items-center gap-1.5">
                  <Users :size="14" aria-hidden="true" /> Antrean berikutnya
                </span>
                <span class="badge">{{ poli.nextPatients.length }} menunggu</span>
              </div>

              <div v-if="poli.nextPatients.length" class="space-y-1">
                <div
                  v-for="p in poli.nextPatients.slice(0, 4)"
                  :key="p.id"
                  class="flex items-center justify-between gap-2 border-b border-[var(--line)] py-2 text-xs last:border-0"
                  :class="{ 'bg-emerald-50/60 font-bold': p.number === searchQuery.toUpperCase() }"
                >
                  <strong class="rounded-lg bg-[var(--deep)]/5 px-2.5 py-1 text-sm tabular-nums text-[var(--deep)] font-bold">
                    {{ p.number }}
                  </strong>
                  <span class="truncate text-slate-700 font-medium text-right max-w-[130px]">
                    {{ p.name }}
                  </span>
                </div>
              </div>

              <p v-else class="text-xs text-[var(--muted)] py-4 text-center">
                Belum ada antrean berikutnya.
              </p>
            </div>

            <!-- 3. Doctor, Room & Operational Info -->
            <div class="pt-3 border-t border-[var(--line)] flex items-center justify-between text-[11px] text-[var(--muted)]">
              <span class="truncate max-w-[140px] font-medium" :title="poli.doctorName">
                {{ poli.doctorName }}
              </span>
              <span class="font-mono text-slate-500">{{ poli.room }}</span>
            </div>
          </div>
        </section>
      </div>

      <!-- Footer Info -->
      <p class="mt-8 flex items-center justify-center gap-2 rounded-xl border border-[var(--line)] bg-white px-4 py-4 text-sm text-[var(--muted)] shadow-2xs">
        <Clock3 :size="17" class="shrink-0 text-[var(--deep)]" aria-hidden="true" />
        Jam pelayanan 08.00 – 14.00 WIB · Pasien disarankan hadir 15 menit sebelum estimasi giliran.
      </p>
    </div>
  </main>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from "vue";
import {
  ChevronRight, Stethoscope, Smile, Baby, HeartPulse,
  MapPin, Clock3, Users, Search, X, CheckCircle2, RefreshCw
} from "lucide-vue-next";
import { useQueueStore } from "../../stores/queue";
import { todayDate } from "../../constants/initialData";

const queueStore = useQueueStore();

// Real-time Digital Clock
const currentTimeString = ref("");
const currentDateString = ref("");
let clockIntervalId: number | null = null;
const refreshing = ref(false);
const searchQuery = ref("");

function updateClock() {
  const now = new Date();
  currentTimeString.value = now.toLocaleTimeString("id-ID", {
    hour: "2-digit",
    minute: "2-digit",
    second: "2-digit",
  }) + " WIB";

  currentDateString.value = now.toLocaleDateString("id-ID", {
    weekday: "long",
    day: "numeric",
    month: "long",
    year: "numeric",
  });
}

// 4 Polyclinics Configuration
const polyclinicsConfig = [
  {
    id: "poli-umum",
    code: "A",
    name: "Poli Umum",
    doctorName: "dr. Rama Palito",
    room: "Ruang 01",
    icon: Stethoscope,
  },
  {
    id: "poli-gigi",
    code: "B",
    name: "Poli Gigi",
    doctorName: "drg. Rani Prameswari",
    room: "Ruang 02",
    icon: Smile,
  },
  {
    id: "poli-kia",
    code: "C",
    name: "Poli KIA / Anak",
    doctorName: "dr. Siti Rahmawati, Sp.A",
    room: "Ruang 03",
    icon: Baby,
  },
  {
    id: "poli-lansia",
    code: "D",
    name: "Poli Lansia",
    doctorName: "dr. Maya Putri",
    room: "Ruang 04",
    icon: HeartPulse,
  },
];

// Reactive mapping for the 4 polyclinics side by side
const polyclinics = computed(() => {
  return polyclinicsConfig.map((config) => {
    // Find patients belonging to this polyclinic by poli name or code prefix
    const poliPatients = queueStore.patients.filter((p) => {
      const matchDate = p.date === todayDate;
      const matchPoli = p.poli && p.poli.toLowerCase().includes(config.name.toLowerCase().slice(0, 6));
      const matchPrefix = p.number && p.number.startsWith(config.code + "-");
      return matchDate && (matchPoli || matchPrefix);
    });

    // Currently called patient
    const currentPatient = poliPatients.find(
      (p) => p.status === "Dipanggil" || p.status === "Sedang Dilayani",
    ) || null;

    // Upcoming patients
    const nextPatients = poliPatients.filter(
      (p) => p.status === "Hadir" || p.status === "Menunggu verifikasi" || p.status === "waiting",
    );

    return {
      ...config,
      currentPatient,
      nextPatients,
      totalToday: poliPatients.length,
    };
  });
});

function isPoliHighlighted(code: string) {
  if (!searchQuery.value.trim()) return false;
  const q = searchQuery.value.trim().toUpperCase();
  return q.startsWith(code + "-") || q === code;
}

const searchedResult = computed(() => {
  if (!searchQuery.value.trim()) return null;
  const q = searchQuery.value.trim().toUpperCase();

  for (const poli of polyclinics.value) {
    if (poli.currentPatient && poli.currentPatient.number.toUpperCase() === q) {
      return {
        patient: poli.currentPatient,
        poliName: poli.name,
        positionText: "Sedang Dipanggil di Poli!",
      };
    }

    const idx = poli.nextPatients.findIndex((item) => item.number.toUpperCase() === q);
    if (idx !== -1) {
      return {
        patient: poli.nextPatients[idx],
        poliName: poli.name,
        positionText: idx === 0
          ? "Nomor Selanjutnya (Berikutnya)"
          : `${idx} antrean lagi di depan Anda`,
      };
    }
  }

  return null;
});

function refreshData() {
  refreshing.value = true;
  try {
    const stored = localStorage.getItem("antrisehat-patients");
    if (stored) {
      queueStore.patients = JSON.parse(stored);
    }
  } catch (_) {}

  setTimeout(() => {
    refreshing.value = false;
  }, 400);
}

onMounted(() => {
  updateClock();
  clockIntervalId = window.setInterval(updateClock, 1000);
});

onUnmounted(() => {
  if (clockIntervalId) {
    clearInterval(clockIntervalId);
  }
});
</script>
