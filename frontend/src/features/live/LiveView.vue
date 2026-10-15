<template>
  <main class="w-full h-[calc(100vh-94px)] max-h-[calc(100vh-94px)] overflow-hidden px-4 sm:px-6 lg:px-8 pt-3 pb-14 flex flex-col justify-between select-none font-sans">
    <!-- 4 COLUMNS SIDE BY SIDE (4 KESAMPING, FULL WIDTH & FIT TO SCREEN) -->
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 w-full h-full min-h-0 flex-1 items-stretch">
      <section
        v-for="poli in polyclinics"
        :key="poli.id"
        class="dashboard-panel !overflow-hidden !p-0 flex flex-col justify-between h-full min-h-0 transition-all duration-200"
        :class="{ '!border-[var(--deep)] !shadow-xl !ring-2 !ring-[var(--deep)]/50': isPoliHighlighted(poli.code) }"
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
                class="flex items-center justify-between gap-3 border-b border-[var(--line)] py-2.5 px-2 text-sm last:border-0 rounded-lg transition-colors"
                :class="p.number === searchQuery.toUpperCase() ? 'bg-emerald-100 font-bold ring-2 ring-emerald-400' : 'hover:bg-slate-50'"
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

    <!-- FLOATING SEARCH BUTTON AT BOTTOM RIGHT -->
    <div class="fixed bottom-3 right-6 z-50">
      <!-- Search Trigger Button -->
      <button
        type="button"
        class="flex items-center gap-2 bg-[var(--deep)] hover:bg-[#124246] text-white shadow-xl shadow-[var(--deep)]/25 px-4 py-2.5 rounded-full font-bold text-xs transition-all transform hover:scale-105 active:scale-95 border-2 border-white"
        @click="toggleSearchModal"
        title="Lacak Posisi Antrean Saya"
      >
        <Search :size="15" />
        <span>Lacak Antrean Saya</span>
        <span
          v-if="searchedResult"
          class="w-2 h-2 rounded-full bg-emerald-300 animate-ping"
        ></span>
      </button>

      <!-- Floating Search Popover -->
      <div
        v-if="isSearchModalOpen"
        class="absolute bottom-12 right-0 w-80 sm:w-96 bg-white rounded-2xl shadow-2xl border border-[var(--line)] p-4 space-y-3"
      >
        <div class="flex items-center justify-between border-b border-[var(--line)] pb-2">
          <span class="text-xs font-bold text-[var(--deep)] flex items-center gap-1.5">
            <Search :size="14" /> Lacak Posisi Antrean
          </span>
          <button
            type="button"
            class="text-slate-400 hover:text-slate-700 p-1 rounded-md transition-colors"
            @click="isSearchModalOpen = false"
          >
            <X :size="15" />
          </button>
        </div>

        <div class="relative">
          <input
            ref="searchInputRef"
            v-model="searchQuery"
            type="text"
            placeholder="Ketik nomor tiket (contoh: A-002, B-001)..."
            class="w-full rounded-xl border border-[var(--line)] bg-[var(--warm)]/40 py-2.5 pl-3 pr-8 text-xs font-mono font-bold text-[var(--deep)] uppercase focus:outline-none focus:ring-2 focus:ring-[var(--deep)]"
          />
          <button
            v-if="searchQuery"
            type="button"
            class="absolute right-2.5 top-1/2 -translate-y-1/2 text-slate-400 hover:text-slate-600"
            @click="searchQuery = ''"
          >
            <X :size="13" />
          </button>
        </div>

        <!-- Matched Result Display -->
        <div
          v-if="searchedResult"
          class="rounded-xl border border-emerald-300 bg-emerald-50/80 p-3 text-xs text-slate-800 space-y-1 shadow-2xs"
        >
          <div class="flex items-center justify-between">
            <span class="font-bold text-emerald-800">{{ searchedResult.poliName }}</span>
            <span class="px-2 py-0.5 rounded-full text-[10px] font-bold bg-white text-emerald-800 border border-emerald-200">
              {{ searchedResult.positionText }}
            </span>
          </div>
          <div class="text-sm font-mono font-black text-[var(--deep)]">
            {{ searchedResult.patient.number }}
          </div>
          <p class="text-[11px] text-slate-600 truncate">
            Atas nama: <strong>{{ searchedResult.patient.name }}</strong>
          </p>
        </div>

        <p v-else-if="searchQuery.trim()" class="text-center text-xs text-[var(--muted)] py-2">
          Nomor antrean tidak ditemukan.
        </p>

        <p class="text-[10px] text-[var(--muted)] text-center">
          Poli terkait otomatis disorot dengan border bersinar pada layar.
        </p>
      </div>
    </div>
  </main>
</template>

<script setup lang="ts">
import { ref, computed, nextTick } from "vue";
import {
  Stethoscope, Smile, Baby, HeartPulse,
  Users, Search, X
} from "lucide-vue-next";
import { useQueueStore } from "../../stores/queue";
import { todayDate } from "../../constants/initialData";
import { getCurrentShiftId } from "../../constants/doctorShifts";

const queueStore = useQueueStore();

const searchQuery = ref("");
const isSearchModalOpen = ref(false);
const searchInputRef = ref<HTMLInputElement | null>(null);

function toggleSearchModal() {
  isSearchModalOpen.value = !isSearchModalOpen.value;
  if (isSearchModalOpen.value) {
    nextTick(() => {
      searchInputRef.value?.focus();
    });
  }
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
  const currentShiftId = getCurrentShiftId();
  return polyclinicsConfig.map((config) => {
    const activeDoctorShift = queueStore.doctorShifts.find(
      (s) => s.poliCode === config.code && s.shiftId === currentShiftId,
    );
    const resolvedDoctorName = activeDoctorShift?.doctorName || config.doctorName;
    const resolvedRoom = activeDoctorShift?.room || config.room;

    const poliPatients = queueStore.patients.filter((p) => {
      const matchDate = p.date === todayDate;
      const matchPoli = p.poli && p.poli.toLowerCase().includes(config.name.toLowerCase().slice(0, 6));
      const matchPrefix = p.number && p.number.startsWith(config.code + "-");
      return matchDate && (matchPoli || matchPrefix);
    });

    const currentPatient = poliPatients.find(
      (p) => p.status === "Dipanggil" || p.status === "Sedang Dilayani",
    ) || null;

    const nextPatients = poliPatients.filter(
      (p) => p.status === "Hadir" || p.status === "Menunggu verifikasi" || p.status === "waiting",
    );

    return {
      ...config,
      doctorName: resolvedDoctorName,
      room: resolvedRoom,
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
          ? "Nomor Selanjutnya"
          : `${idx} antrean lagi di depan Anda`,
      };
    }
  }

  return null;
});
</script>
