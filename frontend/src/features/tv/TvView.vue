<template>
  <main class="w-screen h-screen overflow-hidden flex flex-col justify-between bg-[#0b2225] text-slate-100 p-3 sm:p-4 lg:p-5 select-none font-sans">
    <!-- 1. TV HEADER (Brand + Realtime UTC+7 Clock, No Dashboard Back Button) -->
    <header class="flex items-center justify-between gap-4 pb-3 border-b border-white/10 shrink-0">
      <div class="flex items-center gap-4">
        <Brand />
        <div class="hidden sm:block border-l border-white/20 pl-4">
          <h1 class="text-xs lg:text-sm font-black uppercase tracking-widest text-teal-300">
            Layar Monitor Antrean Poliklinik
          </h1>
          <p class="text-[11px] text-teal-100/70">
            Puskesmas Sehat Bersama · Sistem Terintegrasi
          </p>
        </div>
      </div>

      <!-- Real-time Clock Widget (UTC+7 / WIB) -->
      <div class="flex items-center gap-3 sm:gap-4 bg-white/10 backdrop-blur-md px-3 sm:px-4 py-2 rounded-2xl border border-white/15 shadow-inner">
        <!-- Live Pulse Indicator -->
        <div class="flex items-center gap-2 pr-3 border-r border-white/20">
          <span class="relative flex h-2.5 w-2.5 sm:h-3 sm:w-3">
            <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75"></span>
            <span class="relative inline-flex rounded-full h-2.5 w-2.5 sm:h-3 sm:w-3 bg-emerald-500"></span>
          </span>
          <span class="text-[10px] sm:text-xs font-black tracking-widest text-emerald-400 uppercase">LIVE</span>
        </div>

        <!-- Clock & Date Display in UTC+7 (Asia/Jakarta) -->
        <div class="text-right">
          <div class="font-mono text-xl sm:text-2xl lg:text-3xl font-black tracking-tight text-white leading-none tabular-nums">
            {{ currentTime || "00:00:00 WIB" }}
          </div>
          <div class="text-[10px] sm:text-xs font-semibold text-teal-200 mt-1">
            {{ currentDate }}
          </div>
        </div>
      </div>
    </header>

    <!-- 2. 4 COLUMNS SIDE BY SIDE (POLI UMUM, GIGI, KIA, LANSIA) -->
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-3 sm:gap-4 w-full h-full min-h-0 flex-1 items-stretch py-3">
      <section
        v-for="poli in polyclinics"
        :key="poli.id"
        class="bg-white rounded-2xl shadow-xl overflow-hidden flex flex-col justify-between h-full min-h-0 border transition-all duration-300"
        :class="isPoliCalling(poli.name, poli.code)
          ? 'border-amber-400 ring-4 ring-amber-400/50 shadow-2xl scale-[1.01]'
          : 'border-slate-200/90'"
        :aria-label="`Antrean ${poli.name}`"
      >
        <!-- Column Header with Poli Title & Badge -->
        <div
          class="flex items-center justify-between gap-2 px-4 py-3 text-white shrink-0 transition-colors"
          :class="isPoliCalling(poli.name, poli.code) ? 'bg-amber-600' : 'bg-[#195a5f]'"
        >
          <span class="inline-flex items-center gap-2 font-bold text-sm lg:text-base truncate">
            <component :is="poli.icon" :size="20" aria-hidden="true" class="shrink-0" />
            <span class="truncate">{{ poli.name }}</span>
          </span>
          <span class="rounded-lg bg-white/20 px-2.5 py-0.5 text-sm font-black font-mono tracking-wider shrink-0">
            {{ poli.code }}
          </span>
        </div>

        <!-- Main Body of Column -->
        <div class="p-3 sm:p-4 flex-1 flex flex-col justify-between space-y-3 min-h-0 overflow-hidden bg-slate-50/70">
          <!-- Sedang Dipanggil / Dilayani Stage -->
          <div
            class="rounded-xl px-3 py-4 text-center shrink-0 border transition-all"
            :class="poli.currentPatient?.status === 'Dipanggil'
              ? 'bg-amber-50 border-amber-300 ring-2 ring-amber-400/40'
              : 'bg-teal-50/80 border-teal-200/80'"
          >
            <div class="flex items-center justify-center gap-1.5 mb-1">
              <span
                v-if="poli.currentPatient?.status === 'Dipanggil'"
                class="inline-block w-2 h-2 rounded-full bg-amber-500 animate-ping"
              ></span>
              <span
                class="text-[11px] sm:text-xs uppercase font-extrabold tracking-widest"
                :class="poli.currentPatient?.status === 'Dipanggil' ? 'text-amber-800' : 'text-teal-800'"
              >
                {{ poli.currentPatient?.status === "Dipanggil" ? "SEDANG DIPANGGIL" : "SEDANG DILAYANI" }}
              </span>
            </div>

            <!-- Big Ticket Number -->
            <strong
              class="my-1 block text-5xl sm:text-6xl xl:text-7xl font-black tracking-tight tabular-nums leading-none"
              :class="poli.currentPatient?.status === 'Dipanggil' ? 'text-amber-700 animate-pulse' : 'text-[#195a5f]'"
            >
              {{ poli.currentPatient?.number || "—" }}
            </strong>

            <!-- Patient Name -->
            <div v-if="poli.currentPatient" class="mt-2">
              <span class="inline-flex items-center justify-center rounded-full bg-white px-3.5 py-1 text-xs sm:text-sm font-bold text-slate-800 shadow-xs truncate max-w-[240px] border border-slate-200">
                <span class="truncate">{{ poli.currentPatient.name }}</span>
              </span>
            </div>
            <p v-else class="text-xs text-slate-400 mt-2 font-medium">
              Menunggu panggilan staf
            </p>
          </div>

          <!-- Antrean Berikutnya List -->
          <div class="space-y-1.5 flex-1 min-h-0 flex flex-col justify-start overflow-hidden">
            <div class="flex items-center justify-between text-xs font-bold text-slate-700 border-b border-slate-200 pb-1.5 shrink-0">
              <span class="inline-flex items-center gap-1.5">
                <Users :size="14" aria-hidden="true" class="text-teal-700" /> Antrean berikutnya
              </span>
              <span class="text-[11px] font-bold text-teal-800 bg-teal-100 px-2 py-0.5 rounded-full">
                {{ poli.nextPatients.length }} antre
              </span>
            </div>

            <div v-if="poli.nextPatients.length" class="space-y-1.5 overflow-y-auto pr-0.5 min-h-0 flex-1">
              <div
                v-for="p in poli.nextPatients.slice(0, 5)"
                :key="p.id"
                class="flex items-center justify-between gap-2 border border-slate-200/80 bg-white py-1.5 px-2.5 text-xs rounded-lg shadow-2xs"
              >
                <strong class="rounded-md bg-teal-50 border border-teal-200 px-2 py-0.5 text-sm tabular-nums text-teal-900 font-black shrink-0 tracking-wide">
                  {{ p.number }}
                </strong>
                <span class="truncate text-slate-700 font-semibold text-right text-xs">
                  {{ p.name }}
                </span>
              </div>
            </div>

            <div v-else class="text-xs text-slate-400 py-6 text-center italic font-medium">
              Belum ada antrean berikutnya.
            </div>
          </div>

          <!-- Doctor & Room Info -->
          <div class="pt-2 border-t border-slate-200 flex items-center justify-between text-xs text-slate-500 shrink-0">
            <span class="truncate max-w-[150px] font-medium text-slate-700" :title="poli.doctorName">
              {{ poli.doctorName }}
            </span>
            <span class="font-mono text-slate-800 shrink-0 font-bold bg-slate-200 px-2 py-0.5 rounded text-[11px]">
              {{ poli.room }}
            </span>
          </div>
        </div>
      </section>
    </div>

    <!-- 3. TV RUNNING TICKER FOOTER -->
    <footer class="h-9 bg-white/10 backdrop-blur-md rounded-xl border border-white/15 overflow-hidden flex items-center px-3 shrink-0 text-xs">
      <div class="flex items-center gap-1.5 font-bold text-teal-300 pr-3 border-r border-white/20 shrink-0">
        <Clock3 :size="15" />
        <span class="tracking-wide uppercase text-[11px]">INFO</span>
      </div>
      <div class="overflow-hidden whitespace-nowrap flex-1 ml-3 text-slate-200">
        <div class="ticker-content">
          Selamat datang di Puskesmas Sehat Bersama · Jam Pelayanan: 08.00 – 14.00 WIB (UTC+7) · Mohon persiapkan KTP asli dan Kartu BPJS / Asuransi aktif · Perhatikan nomor tiket dan ruangan poli yang dipanggil · Demi kenyamanan bersama, mohon menjaga kebersihan dan ketertiban.
        </div>
      </div>
    </footer>
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
  Clock3,
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

// ── Polyclinics Configuration (4 Columns) ──
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
  return activePoli.includes(poliName.toLowerCase().slice(0, 6)) || activeNum.startsWith(code + "-");
}

onMounted(() => {
  updateUtc7Clock();
  clockTimer = setInterval(updateUtc7Clock, 1000);
});

onBeforeUnmount(() => {
  if (clockTimer) clearInterval(clockTimer);
});
</script>

<style scoped>
@keyframes marquee {
  0% {
    transform: translateX(100%);
  }
  100% {
    transform: translateX(-100%);
  }
}

.ticker-content {
  display: inline-block;
  white-space: nowrap;
  animation: marquee 38s linear infinite;
}

.ticker-content:hover {
  animation-play-state: paused;
}
</style>
