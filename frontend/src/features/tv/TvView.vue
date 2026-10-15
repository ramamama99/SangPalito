<template>
  <main class="w-full h-screen max-h-screen overflow-hidden px-4 sm:px-6 lg:px-8 py-3 flex flex-col justify-between select-none font-sans bg-[var(--bg)]">
    <!-- Header TV: Brand + Jam Real-Time UTC+7 (Tanpa tombol kembali ke dashboard) -->
    <header class="flex items-center justify-between gap-4 pb-2.5 border-b border-[var(--line)] shrink-0">
      <div class="flex items-center gap-3">
        <Brand />
        <span class="badge !py-1 !px-2.5 !text-xs !font-bold flex items-center gap-1.5">
          <span class="live-dot" /> Layar TV Antrean
        </span>
      </div>

      <!-- Real-time Clock Widget (UTC+7 / WIB) - Diperbesar untuk Layar TV -->
      <div class="flex items-center gap-3 sm:gap-5 bg-white px-4 sm:px-6 py-2 sm:py-2.5 rounded-2xl border border-[var(--line)] shadow-sm">
        <div class="flex items-center gap-2 text-sm sm:text-base font-semibold text-[var(--muted)]">
          <Calendar :size="18" class="text-[var(--deep)] shrink-0" />
          <span class="capitalize">{{ currentDate }}</span>
        </div>
        <div class="h-6 w-px bg-[var(--line)]" />
        <div class="flex items-center gap-2.5 font-mono font-black text-xl sm:text-2xl lg:text-3xl text-[var(--deep)] tabular-nums tracking-tight">
          <Clock :size="22" class="text-[var(--deep)] shrink-0" />
          <span>{{ currentTime }}</span>
        </div>
      </div>
    </header>

    <!-- 4 COLUMNS SIDE BY SIDE (Persis seperti di Antrean Live) -->
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 w-full h-full min-h-0 flex-1 items-stretch py-3">
      <section
        v-for="poli in polyclinics"
        :key="poli.id"
        class="dashboard-panel !overflow-hidden !p-0 flex flex-col justify-between h-full min-h-0 transition-all duration-200"
        :class="{ '!border-[var(--deep)] !shadow-xl !ring-2 !ring-[var(--deep)]/50': isPoliCalling(poli.name, poli.code) }"
        :aria-label="`Antrean ${poli.name}`"
      >
        <!-- Column Header with Poli Title & Badge -->
        <div class="flex items-center justify-between gap-2 border-b border-[var(--line)] bg-[var(--deep)] px-4 py-3 text-white shrink-0">
          <span class="inline-flex items-center gap-2.5 font-bold text-base truncate">
            <component :is="poli.icon" :size="20" aria-hidden="true" class="shrink-0" />
            <span class="truncate">{{ poli.name }}</span>
          </span>
          <span class="rounded-lg bg-white/20 px-2.5 py-1 text-sm font-black font-mono tracking-wider shrink-0">
            {{ poli.code }}
          </span>
        </div>

        <!-- Main Body of the Column -->
        <div class="p-4 sm:p-5 flex-1 flex flex-col justify-between space-y-4 min-h-0 overflow-hidden">
          <!-- 1. Sedang Dilayani Box (Center Stage) -->
          <div class="rounded-2xl bg-[var(--deep)]/5 px-4 py-5 text-center shrink-0 border border-[var(--deep)]/10 shadow-2xs">
            <span class="text-xs sm:text-sm uppercase font-bold tracking-widest text-[var(--muted)]">
              Sedang dilayani
            </span>

            <strong class="my-2 block text-6xl sm:text-7xl xl:text-8xl font-black tracking-tight text-[var(--deep)] tabular-nums leading-none">
              {{ poli.currentPatient?.number || "—" }}
            </strong>

            <div v-if="poli.currentPatient" class="mt-2.5">
              <span class="inline-flex items-center justify-center rounded-full bg-white px-4 py-1.5 text-sm sm:text-base font-bold text-[var(--deep)] shadow-sm truncate max-w-[280px] border border-[var(--line)]">
                <span class="truncate">{{ poli.currentPatient.name }}</span>
              </span>
            </div>
            <p v-else class="text-sm text-[var(--muted)] mt-2 font-medium">
              Belum ada panggilan
            </p>
          </div>

          <!-- 2. Antrean Berikutnya List -->
          <div class="space-y-2 flex-1 min-h-0 flex flex-col justify-start overflow-hidden">
            <div class="flex items-center justify-between text-sm font-bold text-[var(--deep)] border-b border-[var(--line)] pb-2 shrink-0">
              <span class="inline-flex items-center gap-2">
                <Users :size="16" aria-hidden="true" /> Antrean berikutnya
              </span>
              <span class="badge !py-1 !px-2.5 !text-xs !font-bold">{{ poli.nextPatients.length }} menunggu</span>
            </div>

            <div v-if="poli.nextPatients.length" class="space-y-2 overflow-y-auto pr-0.5 min-h-0 flex-1">
              <div
                v-for="p in poli.nextPatients.slice(0, 5)"
                :key="p.id"
                class="flex items-center justify-between gap-3 border-b border-[var(--line)] py-2.5 px-2 text-sm last:border-0 rounded-lg transition-colors hover:bg-slate-50"
              >
                <strong class="rounded-lg bg-[var(--deep)]/10 px-3 py-1.5 text-base sm:text-lg tabular-nums text-[var(--deep)] font-black shrink-0 tracking-wide">
                  {{ p.number }}
                </strong>
                <span class="truncate text-slate-800 font-semibold text-right text-xs sm:text-sm">
                  {{ p.name }}
                </span>
              </div>
            </div>

            <p v-else class="text-sm text-[var(--muted)] py-4 text-center font-medium">
              Belum ada antrean berikutnya.
            </p>
          </div>

          <!-- 3. Doctor & Room Info -->
          <div class="pt-2.5 border-t border-[var(--line)] flex items-center justify-between text-xs sm:text-sm text-[var(--muted)] shrink-0">
            <span class="truncate max-w-[150px] font-semibold text-slate-600" :title="poli.doctorName">
              {{ poli.doctorName }}
            </span>
            <span class="font-mono text-slate-600 shrink-0 font-bold bg-slate-100 px-2 py-0.5 rounded">{{ poli.room }}</span>
          </div>
        </div>
      </section>
    </div>
  </main>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount } from "vue";
import {
  Stethoscope,
  Smile,
  Baby,
  HeartPulse,
  Users,
  Clock,
  Calendar,
} from "lucide-vue-next";
import Brand from "../../components/common/Brand.vue";
import { useQueueStore } from "../../stores/queue";
import { todayDate } from "../../constants/initialData";

const queueStore = useQueueStore();

// ── Realtime Clock with UTC+7 (Asia/Jakarta / WIB) ──
const currentTime = ref("");
const currentDate = ref("");
let clockTimer: ReturnType<typeof setInterval>;

function updateUtc7Clock() {
  const now = new Date();

  // Waktu format 24 jam dengan zona waktu Asia/Jakarta (UTC+7)
  const timeFormatter = new Intl.DateTimeFormat("id-ID", {
    timeZone: "Asia/Jakarta",
    hour: "2-digit",
    minute: "2-digit",
    second: "2-digit",
    hour12: false,
  });

  // Tanggal format Indonesia dengan zona waktu Asia/Jakarta (UTC+7)
  const dateFormatter = new Intl.DateTimeFormat("id-ID", {
    timeZone: "Asia/Jakarta",
    weekday: "long",
    day: "numeric",
    month: "long",
    year: "numeric",
  });

  const timeStr = timeFormatter.format(now).replace(/\./g, ":");
  currentTime.value = `${timeStr} WIB`;
  currentDate.value = dateFormatter.format(now);
}

// ── 4 Polyclinics Configuration (Identik dengan Antrean Live) ──
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

// Reactive mapping for the 4 polyclinics
const polyclinics = computed(() => {
  return polyclinicsConfig.map((config) => {
    const poliPatients = queueStore.patients.filter((p) => {
      const matchDate = p.date === todayDate;
      const matchPoli =
        p.poli &&
        p.poli.toLowerCase().includes(config.name.toLowerCase().slice(0, 6));
      const matchPrefix = p.number && p.number.startsWith(config.code + "-");
      return matchDate && (matchPoli || matchPrefix);
    });

    const currentPatient =
      poliPatients.find(
        (p) => p.status === "Dipanggil" || p.status === "Sedang Dilayani",
      ) || null;

    const nextPatients = poliPatients.filter(
      (p) =>
        p.status === "Hadir" ||
        p.status === "Menunggu verifikasi" ||
        p.status === "waiting",
    );

    return {
      ...config,
      currentPatient,
      nextPatients,
      totalToday: poliPatients.length,
    };
  });
});

function isPoliCalling(poliName: string, code: string) {
  if (!queueStore.active) return false;
  const activePoli = (queueStore.active.poli || "").toLowerCase();
  const activeNum = (queueStore.active.number || "").toUpperCase();
  return (
    activePoli.includes(poliName.toLowerCase().slice(0, 6)) ||
    activeNum.startsWith(code + "-")
  );
}

onMounted(() => {
  updateUtc7Clock();
  clockTimer = setInterval(updateUtc7Clock, 1000);
});

onBeforeUnmount(() => {
  if (clockTimer) clearInterval(clockTimer);
});
</script>
