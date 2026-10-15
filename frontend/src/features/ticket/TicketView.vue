<template>
  <main class="ticket-page container">
    <template v-if="queueStore.currentTicket">
      <Transition name="popup">
    <div v-if="isOpen" class="success-popup">
      <div class="success-banner">
        <CheckCircle2 :size="29" />
        <div>
          <strong> Pendaftaran berhasil. Jangan tutup halaman ini! </strong>
          <p>
            Wajib screenshot atau unduh tiket Anda. Pengiriman WhatsApp belum
            aktif pada prototype ini.
          </p>
        </div>
      </div>
      </div>
      </Transition>
      <div class="ticket-layout">
        <div class="ticket-card">
          <Brand />
          <p>E-TIKET ANTREAN</p>
          <h1>{{ queueStore.currentTicket.number }}</h1>
          <span class="badge">{{ queueStore.currentTicket.status }}</span>
          <img
            v-if="queueStore.qr"
            class="qr"
            :src="queueStore.qr"
            alt="QR check-in antrean"
          />
          <!-- QR Validity Badge: QR is strictly active on visit date only -->
          <div
            v-if="isTodayTicket"
            class="my-2 px-3 py-1.5 rounded-full text-xs font-semibold bg-emerald-50 text-emerald-800 border border-emerald-200 flex items-center justify-center gap-1.5"
          >
            <span class="w-2 h-2 rounded-full bg-emerald-500 animate-pulse"></span>
            QR Aktif Hari Ini · Siap Scan di Lokasi
          </div>
          <div
            v-else-if="isFutureTicket"
            class="my-2 px-3 py-1.5 rounded-full text-xs font-semibold bg-amber-50 text-amber-800 border border-amber-200 flex items-center justify-center gap-1.5"
          >
            <Clock :size="13" class="text-amber-600" />
            QR Hanya Aktif pada Hari Kunjungan ({{ formatDate(queueStore.currentTicket.date) }})
          </div>
          <div
            v-else-if="isPastTicket"
            class="my-2 px-3 py-1.5 rounded-full text-xs font-semibold bg-rose-50 text-rose-800 border border-rose-200 flex items-center justify-center gap-1.5"
          >
            <AlertCircle :size="13" class="text-rose-600" />
            QR Kedaluwarsa (Jadwal: {{ formatDate(queueStore.currentTicket.date) }})
          </div>

          <div v-if="queueStore.currentTicket.qr_token" class="text-center my-2 p-2 bg-slate-50 border border-slate-200 rounded-lg text-xs">
            <span class="block text-[10px] text-slate-500 uppercase tracking-wider font-semibold">Token Check-In</span>
            <code class="text-emerald-700 font-mono font-bold text-[11px] break-all select-all">{{ queueStore.currentTicket.qr_token }}</code>
          </div>
          <h3>{{ queueStore.currentTicket.name }}</h3>
          <p>{{ queueStore.currentTicket.poli }} · Puskesmas Sehat Bersama</p>
          <div class="ticket-details">
            <span>
              Tanggal kunjungan
              <strong>
                {{ formatDate(queueStore.currentTicket.date) }}
              </strong>
            </span>
            <span>
              Estimasi waktu tunggu
              <strong>
                {{ estimatedWaitTime }}
              </strong>
            </span>
          </div>
        </div>
        <div class="ticket-instructions">
          <!-- E-RESEP PASIEN DIGITAL -->
          <div v-if="ticketPrescription || queueStore.currentTicket.medication_plan || queueStore.currentTicket.prescription_no" class="patient-rx-card">
            <div class="rx-card-top">
              <div class="flex items-center gap-2">
                <Pill :size="20" class="text-teal-700" />
                <h3 class="font-bold text-sm text-slate-800">e-Resep Dokter Anda</h3>
              </div>
              <span
                class="rx-status-pill"
                :class="(ticketPrescription?.status === 'dispensed' || queueStore.currentTicket.prescription?.status === 'dispensed') ? 'bg-emerald-100 text-emerald-800 border-emerald-300' : 'bg-amber-100 text-amber-800 border-amber-300'"
              >
                {{ (ticketPrescription?.status === 'dispensed' || queueStore.currentTicket.prescription?.status === 'dispensed') ? '✓ Obat Sudah Ditebus' : '⏳ Menunggu di Farmasi' }}
              </span>
            </div>

            <div class="rx-card-meta">
              <span class="text-xs text-slate-500 font-medium">Nomor e-Resep:</span>
              <strong class="font-mono text-xs text-teal-800 bg-teal-50 px-2 py-0.5 rounded border border-teal-200">
                {{ ticketPrescription?.prescription_no || queueStore.currentTicket.prescription_no || queueStore.currentTicket.prescription?.prescription_no || 'RX-TERTULIS' }}
              </strong>
            </div>

            <!-- List of medicines with detailed usage instructions -->
            <div v-if="rxItemsList.length" class="rx-card-items">
              <div
                v-for="(item, itIdx) in rxItemsList"
                :key="itIdx"
                class="rx-card-item"
              >
                <div class="flex items-center justify-between">
                  <div class="flex items-center gap-2 flex-wrap">
                    <span class="route-badge" :data-route="item.route">
                      {{ formatRouteLabel(item.route) }}
                    </span>
                    <strong class="text-xs text-slate-900">{{ item.medication_name }}</strong>
                  </div>
                  <span class="text-xs font-mono font-bold text-teal-900">{{ item.quantity }} {{ item.unit }}</span>
                </div>
                <div class="text-xs text-slate-600 mt-1 space-y-0.5">
                  <p><strong>Aturan:</strong> {{ item.frequency }} · {{ item.timing }}</p>
                  <p v-if="item.instruction" class="text-teal-800 font-medium">
                    <strong>Petunjuk:</strong> {{ item.instruction }}
                  </p>
                </div>
              </div>
            </div>
            <div v-else-if="queueStore.currentTicket.medication_plan" class="text-xs text-slate-700 whitespace-pre-line p-2.5 bg-slate-50 rounded-lg border border-slate-200">
              {{ queueStore.currentTicket.medication_plan }}
            </div>

            <p class="text-[11px] text-slate-500 mt-2">
              Tunjukkan nomor e-resep ini ke loket Farmasi / Pengambilan Obat puskesmas.
            </p>
          </div>

          <h2>Selanjutnya, apa?</h2>
          <div
            v-for="(inst, i) in instructions"
            :key="inst.title"
            class="instruction"
          >
            <span>{{ i + 1 }}</span>
            <div>
              <h3>{{ inst.title }}</h3>
              <p>{{ inst.desc }}</p>
            </div>
          </div>
          <button
            class="button primary full-width"
            @click="printTicket"
          >
            <Download :size="18" /> Unduh Tiket Gambar
          </button>
          <p class="muted">
            QR dapat dipindai dengan kamera untuk membuka tiket. Check-in dan
            verifikasi tersambung langsung ke database backend.
          </p>
          <button
            class="cancel-link"
            :disabled="isCancelDisabled"
            @click="handleCancel"
          >
            Batalkan Antrean
          </button>
          <RouterLink to="/" class="back-link">
            Kembali ke beranda
          </RouterLink>
        </div>
      </div>
    </template>
    <template v-else>
      <div class="empty-state">
        <QrCode :size="50" />
        <h2>Tiket tidak ditemukan</h2>
        <p>Periksa nomor antrean atau mulai pendaftaran baru.</p>
        <RouterLink to="/register" class="button primary">
          Daftar Online
        </RouterLink>
      </div>
    </template>
  </main>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted, onBeforeUnmount } from "vue";
import { useRouter } from "vue-router";
import {
  CheckCircle2,
  Download,
  QrCode,
  Sparkles,
  Clock,
  AlertCircle,
  Pill,
} from "lucide-vue-next";
import Brand from "../../components/common/Brand.vue";
import { useQueueStore } from "../../stores/queue";
import { api } from "../../services/api";
import type { Prescription, PrescriptionItem } from "../../types";

const router = useRouter();
const queueStore = useQueueStore();

const ticketPrescription = ref<Prescription | null>(null);

const rxItemsList = computed<PrescriptionItem[]>(() => {
  if (ticketPrescription.value?.items?.length) {
    return ticketPrescription.value.items;
  }
  if (queueStore.currentTicket?.prescription?.items?.length) {
    return queueStore.currentTicket.prescription.items;
  }
  return [];
});

watch(
  () => queueStore.currentTicket,
  async (t) => {
    if (!t) {
      ticketPrescription.value = null;
      return;
    }
    if (t.prescription) {
      ticketPrescription.value = t.prescription;
      return;
    }
    try {
      const soap = JSON.parse(t.soap || "{}");
      if (soap.prescription) {
        ticketPrescription.value = soap.prescription;
        return;
      }
    } catch {}

    const token = t.qr_token || t.id;
    if (token) {
      try {
        const rx = await api.getPrescriptionByAppointmentToken(token);
        if (rx) {
          ticketPrescription.value = rx;
        }
      } catch {}
    }
  },
  { immediate: true },
);

function formatRouteLabel(route?: string) {
  if (!route) return "Diminum";
  const routeLabels: Record<string, string> = {
    oral_minum: "Diminum",
    oral_kunyah: "Dikunyah",
    topikal_oles: "Dioleskan",
    tetes_mata: "Tetes Mata",
    tetes_telinga: "Tetes Telinga",
    kumur: "Dikumur",
    inhalasi: "Dihirup",
    sirup_minum: "Sirup",
  };
  return routeLabels[route] || route;
}

const instructions = [
  {
    title: "Simpan tiket Anda",
    desc: "Unduh atau screenshot tiket ini agar mudah ditemukan.",
  },
  {
    title: "Datang ke puskesmas",
    desc: "Bawa KTP dan kartu BPJS / asuransi sesuai penjamin.",
  },
  {
    title: "Scan QR di lokasi",
    desc: "Petugas akan mencocokkan identitas sebelum Anda dilayani.",
  },
];

const isOpen = ref(false);
let successTimer: ReturnType<typeof setTimeout>;

onMounted(() => {
  isOpen.value = true;
  successTimer = setTimeout(() => {
    isOpen.value = false;
  }, 5000); // durasi popup dalam ms (5000 = 5 detik)
});

onBeforeUnmount(() => clearTimeout(successTimer));


const todayDateStr = computed(() => {
  const d = new Date();
  const year = d.getFullYear();
  const month = String(d.getMonth() + 1).padStart(2, "0");
  const day = String(d.getDate()).padStart(2, "0");
  return `${year}-${month}-${day}`;
});

const isTodayTicket = computed(() => {
  const t = queueStore.currentTicket;
  if (!t || !t.date) return true;
  return t.date === todayDateStr.value;
});

const isFutureTicket = computed(() => {
  const t = queueStore.currentTicket;
  if (!t || !t.date) return false;
  return t.date > todayDateStr.value;
});

const isPastTicket = computed(() => {
  const t = queueStore.currentTicket;
  if (!t || !t.date) return false;
  return t.date < todayDateStr.value;
});

const estimatedWaitTime = computed(() => {
  const t = queueStore.currentTicket;
  if (!t) return "—";
  if (["Selesai", "Dibatalkan", "Ditolak"].includes(t.status)) {
    return "—";
  }
  const aheadCount = queueStore.patients.filter(
    (p) =>
      p.poli === t.poli &&
      p.date === t.date &&
      !["Selesai", "Dibatalkan", "Ditolak"].includes(p.status) &&
      p.number < t.number,
  ).length;
  return `${aheadCount * 8} menit`;
});

const isCancelDisabled = computed(() => {
  const t = queueStore.currentTicket;
  if (!t) return true;
  return ["Selesai", "Dibatalkan", "Dipanggil"].includes(t.status);
});

function formatDate(dateStr: string) {
  try {
    return new Date(dateStr).toLocaleDateString("id-ID", { dateStyle: "long" });
  } catch {
    return dateStr;
  }
}

function printTicket() {
  window.print();
  goToHome();
}


function handleCancel() {
  if (queueStore.currentTicket) {
    if (window.confirm("Yakin ingin membatalkan antrean ini?")) {
      queueStore.changeStatus(queueStore.currentTicket.id, "Dibatalkan");
    }
  }
}

function goToHome() {
  router.push("/");
  window.scrollTo(0, 0);
}

function handleNewRegistration() {
  queueStore.startRegistration("Poli Umum", false);
  router.push("/register");
  window.scrollTo(0, 0);
}
</script>

<style scoped>
.patient-rx-card {
  margin-bottom: 24px;
  padding: 18px;
  background: #f0fdf4;
  border: 1.5px solid #86efac;
  border-radius: 16px;
  box-shadow: 0 4px 12px rgba(16, 185, 129, 0.08);
}

.rx-card-top {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 8px;
  flex-wrap: wrap;
  gap: 8px;
}

.rx-status-pill {
  font-size: 11px;
  font-weight: 700;
  padding: 3px 10px;
  border-radius: 9999px;
  border-width: 1px;
}

.rx-card-meta {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 12px;
}

.rx-card-items {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.rx-card-item {
  padding: 10px 14px;
  background: #ffffff;
  border: 1px solid #d1fae5;
  border-radius: 10px;
}

.route-badge {
  display: inline-flex;
  align-items: center;
  padding: 2px 7px;
  border-radius: 6px;
  font-size: 10px;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.3px;
}

.route-badge[data-route="oral_minum"],
.route-badge[data-route="sirup_minum"] {
  background: #d1fae5;
  color: #065f46;
  border: 1px solid #a7f3d0;
}

.route-badge[data-route="oral_kunyah"] {
  background: #fef3c7;
  color: #92400e;
  border: 1px solid #fcd34d;
}

.route-badge[data-route="topikal_oles"] {
  background: #dbeafe;
  color: #1e40af;
  border: 1px solid #bfdbfe;
}

.route-badge[data-route="tetes_mata"],
.route-badge[data-route="tetes_telinga"] {
  background: #f3e8ff;
  color: #6b21a8;
  border: 1px solid #e9d5ff;
}

.route-badge[data-route="kumur"] {
  background: #ffedd5;
  color: #9a3412;
  border: 1px solid #fed7aa;
}

.route-badge[data-route="inhalasi"] {
  background: #e0f2fe;
  color: #0369a1;
  border: 1px solid #bae6fd;
}

.success-popup {
  position: fixed;
  top: 100px;
  left: 50%;
  transform: translateX(-50%);
  z-index: 1000;
  width: min(92vw, 520px);
}

.popup-enter-active,
.popup-leave-active {
  transition: opacity 0.3s ease, transform 0.3s ease;
}
.popup-enter-from,
.popup-leave-to {
  opacity: 0;
  transform: translateX(-50%) translateY(-80px);
}
</style>
