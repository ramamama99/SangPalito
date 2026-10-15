<template>
  <main class="min-h-screen bg-[var(--surface-subtle,#f8fafc)] text-slate-900 py-8 px-4 sm:px-6 lg:px-8 font-sans">
    <div class="mx-auto max-w-6xl space-y-6">
      <!-- TOP NAVIGATION & HEADER -->
      <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 border-b border-slate-200 pb-5">
        <div>
          <RouterLink
            to="/"
            class="inline-flex items-center gap-1.5 text-xs font-semibold text-slate-500 hover:text-emerald-700 transition-colors mb-2"
          >
            <ChevronLeft :size="16" /> Kembali ke Beranda
          </RouterLink>
          <div class="flex items-center gap-2.5">
            <span class="inline-flex p-2 rounded-xl bg-emerald-50 text-emerald-700 border border-emerald-200 shadow-xs">
              <Activity :size="24" class="text-emerald-600 animate-pulse" />
            </span>
            <div>
              <h1 class="text-2xl sm:text-3xl font-black text-slate-900 tracking-tight">
                Pemantauan Antrean Live
              </h1>
              <p class="text-xs sm:text-sm text-slate-500 mt-0.5">
                Pantau antrean berjalan di 4 Poliklinik Puskesmas Sehat Bersama secara real-time.
              </p>
            </div>
          </div>
        </div>

        <!-- Clock, Status & Action Controls -->
        <div class="flex flex-wrap items-center gap-2.5">
          <div class="px-3.5 py-1.5 rounded-xl bg-white border border-slate-200 shadow-xs text-right">
            <div class="text-sm font-mono font-bold text-emerald-700">
              {{ currentTimeString || '00:00:00 WIB' }}
            </div>
            <div class="text-[10px] font-semibold text-slate-400 uppercase tracking-wider">
              {{ currentDateString || 'Memuat Waktu...' }}
            </div>
          </div>

          <div class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-xl text-xs font-bold bg-emerald-50 text-emerald-800 border border-emerald-200 shadow-xs">
            <span class="relative flex h-2 w-2">
              <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75"></span>
              <span class="relative inline-flex rounded-full h-2 w-2 bg-emerald-500"></span>
            </span>
            Live Sinkron
          </div>

          <button
            type="button"
            class="p-2 sm:px-3 sm:py-2 rounded-xl text-xs font-semibold bg-white hover:bg-slate-50 border border-slate-200 text-slate-700 transition-colors shadow-xs flex items-center gap-1.5"
            @click="soundEnabled = !soundEnabled"
            :title="soundEnabled ? 'Matikan Suara Panggilan' : 'Aktifkan Suara Panggilan'"
          >
            <Volume2 v-if="soundEnabled" :size="15" class="text-emerald-600" />
            <VolumeX v-else :size="15" class="text-slate-400" />
            <span class="hidden sm:inline">{{ soundEnabled ? 'Audio Aktif' : 'Mute' }}</span>
          </button>

          <button
            type="button"
            class="p-2 sm:px-3 sm:py-2 rounded-xl text-xs font-semibold bg-white hover:bg-slate-50 border border-slate-200 text-slate-700 transition-colors shadow-xs flex items-center gap-1.5"
            @click="refreshData"
            :disabled="refreshing"
            title="Segarkan antrean sekarang"
          >
            <RefreshCw :size="14" :class="{ 'animate-spin': refreshing }" />
            <span class="hidden sm:inline">Segarkan</span>
          </button>
        </div>
      </div>

      <!-- TICKET FINDER SEARCH BAR (LACAK NOMOR TIKET SAYA) -->
      <div class="bg-white rounded-2xl border border-slate-200 p-4 sm:p-5 shadow-xs space-y-3">
        <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3">
          <div class="flex items-center gap-2">
            <div class="p-2 rounded-lg bg-emerald-50 text-emerald-700">
              <Search :size="18" />
            </div>
            <div>
              <h2 class="text-sm font-bold text-slate-900">Lacak Posisi Antrean Anda</h2>
              <p class="text-xs text-slate-500">Ketikkan nomor tiket Anda untuk melihat estimasi dan poli tujuan secara instan.</p>
            </div>
          </div>

          <div class="relative w-full sm:w-72">
            <input
              v-model="searchQuery"
              type="text"
              placeholder="Contoh: A-002, B-001, C-003..."
              class="w-full rounded-xl border border-slate-300 py-2 pl-3.5 pr-8 text-xs font-mono font-bold text-slate-900 placeholder:font-sans placeholder:font-normal focus:ring-2 focus:ring-emerald-500 focus:outline-none bg-slate-50/60 uppercase"
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

        <!-- SEARCH RESULT ALERT BANNER -->
        <div
          v-if="searchedPatientResult"
          class="p-3.5 rounded-xl border border-emerald-300 bg-gradient-to-r from-emerald-50 to-teal-50 text-xs text-slate-800 flex flex-col sm:flex-row sm:items-center justify-between gap-3 shadow-xs"
        >
          <div class="flex items-center gap-2.5">
            <CheckCircle2 :size="20" class="text-emerald-600 shrink-0" />
            <div>
              <p class="font-bold text-slate-900">
                Tiket <span class="font-mono text-emerald-700 text-sm">{{ searchedPatientResult.patient.number }}</span> terdaftar di
                <span class="underline text-emerald-800">{{ searchedPatientResult.poliName }}</span>
              </p>
              <p class="text-[11px] text-slate-600 mt-0.5">
                Atas nama: <strong>{{ searchedPatientResult.patient.name }}</strong> ·
                Status saat ini: <strong class="text-emerald-700">{{ searchedPatientResult.patient.status }}</strong>
              </p>
            </div>
          </div>

          <div class="flex items-center gap-2 self-start sm:self-auto">
            <span class="px-2.5 py-1 rounded-full text-xs font-bold bg-white text-emerald-800 border border-emerald-200 shadow-2xs font-mono">
              {{ searchedPatientResult.positionText }}
            </span>
          </div>
        </div>
      </div>

      <!-- 4-QUADRANT GRID: DIVIDED INTO 4 POLYCLINICS -->
      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <!-- POLI CARD TEMPLATE (REPEATED FOR 4 POLIS) -->
        <div
          v-for="poli in polyclinics"
          :key="poli.id"
          class="bg-white rounded-2xl border transition-all duration-300 p-5 sm:p-6 shadow-sm flex flex-col justify-between relative overflow-hidden"
          :class="[
            poli.borderClass,
            isPoliHighlighted(poli.code) ? 'ring-3 ring-emerald-500 shadow-md scale-[1.01]' : 'hover:shadow-md'
          ]"
        >
          <!-- Accent Top Bar -->
          <div class="absolute top-0 inset-x-0 h-1.5" :class="poli.topBarClass"></div>

          <div class="space-y-4">
            <!-- Poli Card Header -->
            <div class="flex items-start justify-between gap-2 border-b border-slate-100 pb-3.5">
              <div class="flex items-center gap-3">
                <div class="w-11 h-11 rounded-xl flex items-center justify-center shrink-0 border" :class="poli.iconContainerClass">
                  <component :is="poli.icon" :size="22" />
                </div>
                <div>
                  <div class="flex items-center gap-1.5">
                    <span class="px-2 py-0.5 rounded text-[10px] font-black uppercase tracking-wider" :class="poli.badgeClass">
                      Kode: {{ poli.code }}
                    </span>
                    <span class="text-xs font-semibold text-slate-500">· {{ poli.room }}</span>
                  </div>
                  <h3 class="text-lg font-bold text-slate-900 mt-0.5">{{ poli.name }}</h3>
                </div>
              </div>

              <!-- Status Tag -->
              <span class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-[11px] font-bold bg-emerald-50 text-emerald-700 border border-emerald-200">
                <span class="w-1.5 h-1.5 rounded-full bg-emerald-500 animate-pulse"></span>
                Melayani
              </span>
            </div>

            <!-- Doctor & Service Time Info -->
            <div class="flex items-center justify-between text-xs text-slate-600 bg-slate-50 px-3 py-2 rounded-xl border border-slate-100">
              <span class="flex items-center gap-1 font-medium truncate max-w-[200px]">
                <Users :size="13" class="text-slate-400 shrink-0" /> {{ poli.doctorName }}
              </span>
              <span class="text-slate-400 font-mono text-[11px]">
                Est. ~{{ poli.avgMinutes }} mnt/pasien
              </span>
            </div>

            <!-- 1. CURRENT SERVING BOX (SEDANG DILAYANI) -->
            <div
              class="rounded-xl border p-4 text-center space-y-2 relative overflow-hidden transition-all"
              :class="poli.currentServing ? poli.activeBoxClass : 'bg-slate-50/70 border-slate-200'"
            >
              <div class="flex items-center justify-between text-xs">
                <span class="font-bold uppercase tracking-wider text-[11px]" :class="poli.currentServing ? poli.activeSubtextClass : 'text-slate-500'">
                  Sedang Dilayani Sekarang
                </span>
                <span v-if="poli.currentServing" class="inline-flex items-center gap-1 text-[11px] font-bold px-2 py-0.5 rounded-full bg-white/80 shadow-2xs">
                  <span class="w-1.5 h-1.5 rounded-full bg-current animate-ping"></span>
                  Di Ruang Periksa
                </span>
              </div>

              <!-- Huge Ticket Number Display -->
              <div v-if="poli.currentServing" class="py-1">
                <div class="text-4xl sm:text-5xl font-mono font-black tracking-tight" :class="poli.activeNumberClass">
                  {{ poli.currentServing.number }}
                </div>
                <p class="text-xs font-semibold text-slate-700 mt-1 truncate">
                  {{ poli.currentServing.name }}
                </p>
              </div>

              <!-- Empty State when no one is called yet -->
              <div v-else class="py-4 text-slate-400 space-y-1">
                <Clock3 :size="28" class="mx-auto text-slate-300" />
                <p class="text-xs font-bold text-slate-600">Belum Ada Panggilan Aktif</p>
                <p class="text-[11px] text-slate-400">Menunggu dokter memulai panggilan antrean.</p>
              </div>
            </div>

            <!-- 2. NEXT IN LINE (ANTREAN BERIKUTNYA YANG MENUNGGU) -->
            <div class="space-y-2">
              <div class="flex items-center justify-between text-xs">
                <span class="font-bold text-slate-700 flex items-center gap-1.5">
                  <Clock :size="14" class="text-slate-400" /> Antrean Berikutnya
                </span>
                <span class="font-mono text-[11px] font-semibold text-slate-500">
                  {{ poli.nextQueues.length }} menunggu
                </span>
              </div>

              <div v-if="poli.nextQueues.length === 0" class="py-3 text-center text-xs text-slate-400 border border-dashed border-slate-200 rounded-xl">
                Tidak ada pasien menunggu saat ini.
              </div>

              <div v-else class="space-y-1.5">
                <div
                  v-for="(item, idx) in poli.nextQueues.slice(0, 3)"
                  :key="item.id"
                  class="flex items-center justify-between px-3 py-2 rounded-xl text-xs border transition-colors bg-white hover:bg-slate-50"
                  :class="item.number === searchQuery.toUpperCase() ? 'border-emerald-400 bg-emerald-50/50' : 'border-slate-200'"
                >
                  <div class="flex items-center gap-2.5 min-w-0">
                    <span class="w-5 h-5 rounded-md flex items-center justify-center font-bold text-[10px] bg-slate-100 text-slate-600 shrink-0 font-mono">
                      {{ idx + 1 }}
                    </span>
                    <div class="min-w-0">
                      <strong class="font-mono font-bold text-slate-900 text-sm mr-1.5">{{ item.number }}</strong>
                      <span class="text-slate-600 font-medium truncate">{{ item.name }}</span>
                    </div>
                  </div>

                  <div class="flex items-center gap-1.5 shrink-0">
                    <span class="text-[11px] font-mono text-slate-500">
                      ~{{ (idx + 1) * poli.avgMinutes }} mnt
                    </span>
                    <span class="w-1.5 h-1.5 rounded-full bg-emerald-500"></span>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- 3. FOOTER STATS PER POLI -->
          <div class="pt-3.5 mt-4 border-t border-slate-100 flex items-center justify-between text-[11px] text-slate-500">
            <span>Total Antrean Hari Ini: <strong class="text-slate-800 font-mono">{{ poli.totalToday }}</strong></span>
            <span class="flex items-center gap-1 font-semibold text-emerald-700">
              <Check :size="12" /> Ruang Siap
            </span>
          </div>
        </div>
      </div>

      <!-- BOTTOM SERVICE SCHEDULE & HELPFUL QUICK LINKS -->
      <div class="bg-white rounded-2xl border border-slate-200 p-5 shadow-xs flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div class="flex items-center gap-3 text-xs text-slate-600">
          <div class="p-2.5 rounded-xl bg-slate-100 text-slate-700 shrink-0">
            <Clock3 :size="20" />
          </div>
          <div>
            <strong class="text-slate-900 block font-bold">Jam Layanan Poliklinik: 08.00 – 14.00 WIB</strong>
            <p class="text-slate-500 mt-0.5">
              Pasien yang telah mendaftar online disarankan hadir 15 menit sebelum estimasi untuk scan tiket di Kios Mandiri.
            </p>
          </div>
        </div>

        <div class="flex flex-wrap items-center gap-2 shrink-0">
          <RouterLink
            to="/register"
            class="py-2 px-3.5 rounded-xl text-xs font-bold bg-emerald-600 hover:bg-emerald-700 text-white transition-colors flex items-center gap-1.5 shadow-xs shadow-emerald-600/20"
          >
            <Plus :size="14" /> Daftar Antrean Baru
          </RouterLink>
          <RouterLink
            to="/scan"
            class="py-2 px-3.5 rounded-xl text-xs font-semibold bg-white hover:bg-slate-50 border border-slate-300 text-slate-700 transition-colors flex items-center gap-1.5"
          >
            <QrCode :size="14" class="text-emerald-600" /> Kios Scan Mandiri
          </RouterLink>
        </div>
      </div>
    </div>
  </main>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from "vue";
import {
  ChevronLeft, Activity, Volume2, VolumeX, RefreshCw,
  Search, X, CheckCircle2, Stethoscope, Smile, Baby,
  HeartPulse, Users, Clock3, Clock, Check, Plus, QrCode
} from "lucide-vue-next";
import { useQueueStore } from "../../stores/queue";
import { todayDate } from "../../constants/initialData";
import { speakQueueCall } from "../../utils/speech";

const queueStore = useQueueStore();

// Real-time Digital Clock
const currentTimeString = ref("");
const currentDateString = ref("");
let clockIntervalId: number | null = null;
const soundEnabled = ref(true);
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

// ── THE 4 POLYCLINICS CONFIGURATION ──────────────────────────
const polyclinicsConfig = [
  {
    id: "11111111-1111-1111-1111-111111111111",
    code: "A",
    name: "Poli Umum",
    doctorName: "dr. Rama Palito",
    room: "Ruang Periksa 1 (Lt. 1)",
    avgMinutes: 10,
    icon: Stethoscope,
    borderClass: "border-emerald-200/90",
    topBarClass: "bg-emerald-600",
    badgeClass: "bg-emerald-100 text-emerald-800",
    iconContainerClass: "bg-emerald-50 border-emerald-200 text-emerald-700",
    activeBoxClass: "bg-gradient-to-br from-emerald-50 via-teal-50/40 to-white border-emerald-300 shadow-xs",
    activeNumberClass: "text-emerald-700",
    activeSubtextClass: "text-emerald-800",
  },
  {
    id: "22222222-2222-2222-2222-222222222222",
    code: "B",
    name: "Poli Gigi & Mulut",
    doctorName: "drg. Rani Prameswari",
    room: "Ruang Periksa 2 (Lt. 1)",
    avgMinutes: 15,
    icon: Smile,
    borderClass: "border-sky-200/90",
    topBarClass: "bg-sky-600",
    badgeClass: "bg-sky-100 text-sky-800",
    iconContainerClass: "bg-sky-50 border-sky-200 text-sky-700",
    activeBoxClass: "bg-gradient-to-br from-sky-50 via-blue-50/40 to-white border-sky-300 shadow-xs",
    activeNumberClass: "text-sky-700",
    activeSubtextClass: "text-sky-800",
  },
  {
    id: "33333333-3333-3333-3333-333333333333",
    code: "C",
    name: "Poli KIA / Anak",
    doctorName: "dr. Siti Rahmawati, Sp.A",
    room: "Ruang Periksa 3 (Lt. 2)",
    avgMinutes: 12,
    icon: Baby,
    borderClass: "border-rose-200/90",
    topBarClass: "bg-rose-500",
    badgeClass: "bg-rose-100 text-rose-800",
    iconContainerClass: "bg-rose-50 border-rose-200 text-rose-700",
    activeBoxClass: "bg-gradient-to-br from-rose-50 via-pink-50/40 to-white border-rose-300 shadow-xs",
    activeNumberClass: "text-rose-600",
    activeSubtextClass: "text-rose-800",
  },
  {
    id: "44444444-4444-4444-4444-444444444444",
    code: "D",
    name: "Poli Lansia / Geriatri",
    doctorName: "dr. Maya Putri",
    room: "Ruang Periksa 4 (Lt. 1 Khusus)",
    avgMinutes: 15,
    icon: HeartPulse,
    borderClass: "border-purple-200/90",
    topBarClass: "bg-purple-600",
    badgeClass: "bg-purple-100 text-purple-800",
    iconContainerClass: "bg-purple-50 border-purple-200 text-purple-700",
    activeBoxClass: "bg-gradient-to-br from-purple-50 via-indigo-50/40 to-white border-purple-300 shadow-xs",
    activeNumberClass: "text-purple-700",
    activeSubtextClass: "text-purple-800",
  },
];

// Reactive data mapped for all 4 polyclinics
const polyclinics = computed(() => {
  return polyclinicsConfig.map((config) => {
    // Filter patients belonging to this polyclinic by name or code prefix
    const poliPatients = queueStore.patients.filter((p) => {
      const matchPoliName = p.poli && p.poli.toLowerCase().includes(config.name.toLowerCase().slice(0, 7));
      const matchPrefix = p.number && p.number.startsWith(config.code + "-");
      return matchPoliName || matchPrefix;
    });

    // Currently serving patient (status 'Dipanggil' or 'Sedang Diperiksa')
    const currentServing = poliPatients.find(
      (p) => p.status === "Dipanggil" || p.status === "Sedang Diperiksa",
    ) || null;

    // Upcoming queue (status 'Hadir' or 'Menunggu verifikasi' or 'waiting')
    const nextQueues = poliPatients.filter(
      (p) => p.status === "Hadir" || p.status === "Menunggu verifikasi" || p.status === "waiting",
    );

    return {
      ...config,
      currentServing,
      nextQueues,
      totalToday: poliPatients.length,
    };
  });
});

// Search highlight check
function isPoliHighlighted(code: string) {
  if (!searchQuery.value.trim()) return false;
  const q = searchQuery.value.trim().toUpperCase();
  return q.startsWith(code + "-") || q === code;
}

// Searched Patient Result computation
const searchedPatientResult = computed(() => {
  if (!searchQuery.value.trim()) return null;
  const q = searchQuery.value.trim().toUpperCase();

  for (const poli of polyclinics.value) {
    if (poli.currentServing && poli.currentServing.number.toUpperCase() === q) {
      return {
        patient: poli.currentServing,
        poliName: poli.name,
        positionText: "Sedang Dipanggil di Poli!",
      };
    }

    const idx = poli.nextQueues.findIndex((item) => item.number.toUpperCase() === q);
    if (idx !== -1) {
      const aheadCount = idx;
      return {
        patient: poli.nextQueues[idx],
        poliName: poli.name,
        positionText: aheadCount === 0
          ? "Nomor Selanjutnya (Berikutnya)"
          : `${aheadCount} orang lagi di depan Anda`,
      };
    }
  }

  return null;
});

async function refreshData() {
  refreshing.value = true;
  try {
    const storedStr = localStorage.getItem("antrisehat-patients");
    if (storedStr) {
      queueStore.patients = JSON.parse(storedStr);
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
