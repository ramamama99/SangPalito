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
          <!-- Header Tiket: Brand & Badge E-Tiket -->
          <div class="ticket-header-bar">
            <Brand />
            <div class="ticket-type-pill">
              <Ticket :size="13" />
              <span>E-Tiket Antrean</span>
            </div>
          </div>

          <div class="ticket-queue-section">
            <span class="ticket-queue-label">Nomor Antrean Anda</span>
            <h1 class="ticket-queue-number">{{ queueStore.currentTicket.number }}</h1>
            <span class="badge">{{ queueStore.currentTicket.status }}</span>
          </div>
          <img
            v-if="queueStore.qr"
            class="qr"
            :src="queueStore.qr"
            alt="QR check-in antrean"
          />
          <!-- QR Validity Badge: QR is strictly active on visit date only -->
          <div
            v-if="isTodayTicket"
            class="ticket-validity-badge today"
          >
            <span class="live-dot"></span>
            QR Aktif Hari Ini · Siap Scan di Lokasi
          </div>
          <div
            v-else-if="isFutureTicket"
            class="ticket-validity-badge future"
          >
            <Clock :size="13" />
            QR Hanya Aktif pada Hari Kunjungan ({{ formatDate(queueStore.currentTicket.date) }})
          </div>
          <div
            v-else-if="isPastTicket"
            class="ticket-validity-badge past"
          >
            <AlertCircle :size="13" />
            QR Kedaluwarsa (Jadwal: {{ formatDate(queueStore.currentTicket.date) }})
          </div>

          <div v-if="queueStore.currentTicket.qr_token" class="ticket-token-box">
            <span class="ticket-token-label">Token Check-In</span>
            <code class="ticket-token-code">{{ queueStore.currentTicket.qr_token }}</code>
          </div>
          <h3>{{ queueStore.currentTicket.name }}</h3>
          <p>{{ queueStore.currentTicket.poli }} · Klinik Buah Batu</p>
          <div class="ticket-details">
            <span>
              Tanggal kunjungan
              <strong>
                {{ formatDate(queueStore.currentTicket.date) }}
              </strong>
            </span>
            <span>
              Sesi Pelayanan
              <strong>
                {{ ticketShiftLabel }}
              </strong>
            </span>
            <span class="full-detail-span">
              Estimasi Kedatangan
              <strong class="highlight-time">
                {{ ticketEstimatedArrival }}
              </strong>
            </span>
          </div>

          <!-- Peringatan Datang 15 Menit Sebelum Estimasi Kedatangan -->
          <div class="ticket-warning-box">
            <div class="warning-title">
              <Clock :size="16" />
              <span>Wajib Hadir 15 Menit Lebih Awal</span>
            </div>
            <p class="warning-text">
              Harap hadir paling lambat pukul <strong class="time-highlight">{{ recommendedArrivalTime }}</strong> untuk verifikasi berkas dan skrining awal di loket klinik.
            </p>
          </div>
        </div>
        <div class="ticket-instructions">
          <!-- E-RESEP PASIEN DIGITAL -->
          <div v-if="ticketPrescription || queueStore.currentTicket.medication_plan || queueStore.currentTicket.prescription_no" class="patient-rx-card">
            <div class="rx-card-top">
              <div class="rx-card-title-wrap">
                <Pill :size="18" class="rx-card-icon" />
                <h3 class="rx-card-title">e-Resep Dokter Anda</h3>
              </div>
              <span
                class="rx-status-pill"
                :class="(ticketPrescription?.status === 'dispensed' || queueStore.currentTicket.prescription?.status === 'dispensed') ? 'dispensed' : 'pending'"
              >
                {{ (ticketPrescription?.status === 'dispensed' || queueStore.currentTicket.prescription?.status === 'dispensed') ? '✓ Obat Sudah Ditebus' : '⏳ Menunggu di Farmasi' }}
              </span>
            </div>

            <div class="rx-card-meta">
              <span class="rx-meta-label">Nomor e-Resep:</span>
              <strong class="rx-meta-no">
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
                <div class="rx-item-header">
                  <div class="rx-item-med">
                    <span class="route-badge" :data-route="item.route">
                      {{ formatRouteLabel(item.route) }}
                    </span>
                    <strong class="rx-item-name">{{ item.medication_name }}</strong>
                  </div>
                  <span class="rx-item-qty">{{ item.quantity }} {{ item.unit }}</span>
                </div>
                <div class="rx-item-details">
                  <p><strong>Aturan:</strong> {{ item.frequency }} · {{ item.timing }}</p>
                  <p v-if="item.instruction" class="rx-item-instruction">
                    <strong>Petunjuk:</strong> {{ item.instruction }}
                  </p>
                </div>
              </div>
            </div>
            <div v-else-if="queueStore.currentTicket.medication_plan" class="rx-card-plan">
              {{ queueStore.currentTicket.medication_plan }}
            </div>

            <p class="rx-card-footer">
              Tunjukkan nomor e-resep ini ke loket Farmasi / Pengambilan Obat klinik.
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
  Ticket,
} from "lucide-vue-next";
import Brand from "../../components/common/Brand.vue";
import { useQueueStore } from "../../stores/queue";
import { api } from "../../services/api";
import type { Prescription, PrescriptionItem } from "../../types";
import {
  calculateEstimatedArrival,
  calculateRecommendedArrivalTime,
  SHIFT_TIMES,
} from "../../constants/doctorShifts";

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
    title: "Datang ke klinik",
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

const ticketShiftLabel = computed(() => {
  const t = queueStore.currentTicket;
  if (!t) return "Shift Pagi";
  const shiftKey = (t.shift as "pagi" | "siang" | "sore") || "pagi";
  const info = SHIFT_TIMES[shiftKey] || SHIFT_TIMES.pagi;
  return `Shift ${info.label}`;
});

const ticketEstimatedArrival = computed(() => {
  const t = queueStore.currentTicket;
  if (!t) return "07.00 WIB";
  if (t.estimatedArrival) return t.estimatedArrival;
  const shiftKey = (t.shift as "pagi" | "siang" | "sore") || "pagi";
  const numPart = parseInt(t.number.replace(/\D/g, ""), 10) || 1;
  const smartAvg = t.poli.includes("Gigi") ? 15 : t.poli.includes("KIA") ? 12 : t.poli.includes("Lansia") ? 15 : 10;
  return calculateEstimatedArrival(shiftKey, smartAvg, numPart);
});

const recommendedArrivalTime = computed(() => {
  return calculateRecommendedArrivalTime(ticketEstimatedArrival.value);
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
/* Validity Badge (Soft Design System) */
.ticket-validity-badge {
  margin: 8px 0;
  padding: 6px 14px;
  border-radius: 50px;
  font-size: 12px;
  font-weight: 600;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 7px;
}
.ticket-validity-badge.today {
  background: #eaf6f2;
  color: var(--deep);
  border: 1px solid #c8e8df;
}
.ticket-validity-badge.future {
  background: var(--warm);
  color: var(--deep);
  border: 1px solid #e7ded4;
}
.ticket-validity-badge.past {
  background: #fdf0ed;
  color: var(--coral);
  border: 1px solid #f7d6cf;
}

/* Check-In Token */
.ticket-token-box {
  text-align: center;
  margin: 10px 0;
  padding: 10px 14px;
  background: #f6f8f5;
  border: 1px solid var(--line);
  border-radius: 9px;
  display: flex;
  flex-direction: column;
  gap: 2px;
}
.ticket-token-label {
  font-size: 10.5px;
  color: var(--muted);
  text-transform: uppercase;
  letter-spacing: 0.8px;
  font-weight: 600;
}
.ticket-token-code {
  color: var(--deep);
  font-family: monospace;
  font-weight: 750;
  font-size: 12px;
  word-break: break-all;
  user-select: all;
}

/* e-Resep Pasien Card */
.patient-rx-card {
  margin-bottom: 24px;
  padding: 20px;
  background: #ffffff;
  border: 1px solid var(--line);
  border-radius: 19px;
  box-shadow: 0 4px 15px rgba(25, 90, 95, 0.04);
}
.rx-card-top {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 10px;
  flex-wrap: wrap;
  gap: 8px;
}
.rx-card-title-wrap {
  display: flex;
  align-items: center;
  gap: 8px;
}
.rx-card-icon {
  color: var(--teal);
}
.rx-card-title {
  font-family: "Manrope", sans-serif;
  font-weight: 700;
  font-size: 15px;
  color: var(--deep);
  margin: 0;
}
.rx-status-pill {
  font-size: 11px;
  font-weight: 600;
  padding: 3px 10px;
  border-radius: 50px;
}
.rx-status-pill.dispensed {
  background: #eaf6f2;
  color: var(--deep);
  border: 1px solid #c8e8df;
}
.rx-status-pill.pending {
  background: var(--warm);
  color: var(--deep);
  border: 1px solid #e7ded4;
}
.rx-card-meta {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 14px;
}
.rx-meta-label {
  font-size: 12px;
  color: var(--muted);
  font-weight: 500;
}
.rx-meta-no {
  font-family: monospace;
  font-size: 12px;
  color: var(--deep);
  background: #f0f4f2;
  padding: 2px 8px;
  border-radius: 6px;
  border: 1px solid var(--line);
}
.rx-card-items {
  display: flex;
  flex-direction: column;
  gap: 10px;
}
.rx-card-item {
  padding: 12px 14px;
  background: #f6f8f5;
  border: 1px solid var(--line);
  border-radius: 12px;
}
.rx-item-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
}
.rx-item-med {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}
.rx-item-name {
  font-size: 13px;
  color: #243e40;
  font-weight: 700;
}
.rx-item-qty {
  font-size: 12px;
  font-family: monospace;
  font-weight: 700;
  color: var(--deep);
}
.rx-item-details {
  font-size: 12px;
  color: var(--muted);
  margin-top: 4px;
  line-height: 1.5;
}
.rx-item-details p {
  margin: 0;
  color: var(--muted);
}
.rx-item-instruction {
  color: var(--deep);
  font-weight: 550;
  margin-top: 2px !important;
}
.rx-card-plan {
  font-size: 12px;
  color: #243e40;
  white-space: pre-line;
  padding: 12px 14px;
  background: #f6f8f5;
  border-radius: 9px;
  border: 1px solid var(--line);
}
.rx-card-footer {
  font-size: 11.5px;
  color: var(--muted);
  margin-top: 10px;
}

/* Route Badges (Soft Palette) */
.route-badge {
  display: inline-flex;
  align-items: center;
  padding: 2px 7px;
  border-radius: 5px;
  font-size: 10px;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.3px;
  background: #eaf6f2;
  color: var(--deep);
  border: 1px solid #c8e8df;
}
.route-badge[data-route="oral_kunyah"],
.route-badge[data-route="sirup_minum"] {
  background: var(--warm);
  color: var(--deep);
  border: 1px solid #e7ded4;
}

/* Success Banner Popup */
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

.full-detail-span {
  grid-column: 1 / -1;
  background: #f6f8f5;
  padding: 10px 14px;
  border-radius: 9px;
  border: 1px solid var(--line);
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.highlight-time {
  font-size: 16px !important;
  color: var(--deep) !important;
  font-weight: 750 !important;
}

.ticket-warning-box {
  margin-top: 16px;
  background: var(--warm);
  border: 1px solid #e7ded4;
  border-radius: 14px;
  padding: 14px 16px;
  text-align: left;
}
.warning-title {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 11.5px;
  font-weight: 650;
  color: var(--deep);
  text-transform: uppercase;
  letter-spacing: 0.5px;
  margin-bottom: 5px;
}
.warning-text {
  font-size: 12.5px !important;
  color: var(--deep) !important;
  line-height: 1.5 !important;
  margin: 0 !important;
}
.time-highlight {
  color: var(--coral);
  font-weight: 750;
}

/* Header Tiket Modern & Rapi */
.ticket-header-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding-bottom: 16px;
  margin-bottom: 20px;
  border-bottom: 1px dashed var(--line);
}
.ticket-type-pill {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  background: #eaf6f2;
  color: var(--deep);
  border: 1px solid #c8e8df;
  padding: 4px 12px;
  border-radius: 50px;
  font-size: 11.5px;
  font-weight: 600;
  letter-spacing: 0.3px;
  white-space: nowrap;
}
.ticket-queue-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  margin-bottom: 14px;
}
.ticket-queue-label {
  font-size: 11px;
  font-weight: 650;
  text-transform: uppercase;
  letter-spacing: 1.5px;
  color: var(--muted);
  margin-bottom: 4px;
}
.ticket-queue-number {
  font-family: "Manrope", sans-serif !important;
  font-size: 84px !important;
  letter-spacing: -3px !important;
  margin: 0 0 8px 0 !important;
  color: var(--deep) !important;
  font-weight: 800 !important;
  line-height: 1 !important;
}
@media (max-width: 400px) {
  .ticket-header-bar {
    flex-direction: column;
    gap: 10px;
    align-items: center;
  }
}
</style>
