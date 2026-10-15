<template>
  <main class="ticket-page container">
    <template v-if="queueStore.currentTicket">
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
          <button
            class="button primary full-width"
            @click="queueStore.downloadTicket"
          >
            <Download :size="18" /> Unduh Tiket Gambar
          </button>
          <button class="button outline full-width" @click="printTicket">
            Cetak / Simpan PDF
          </button>
        </div>
        <div class="ticket-instructions">
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
          <RouterLink
            :to="`/qr-test?token=${encodeURIComponent(queueStore.currentTicket.qr_token || queueStore.currentTicket.id)}`"
            class="button primary full-width"
            style="text-decoration: none; justify-content: center; gap: 8px;"
          >
            <Sparkles :size="18" /> Uji Scan QR di Scanner Live
          </RouterLink>
          <button
            class="button outline full-width"
            :disabled="['Hadir', 'Dipanggil', 'Sedang dilayani', 'Selesai'].includes(queueStore.currentTicket.status)"
            @click="handleCheckinSimulation"
          >
            <QrCode :size="18" />
            {{ ['Hadir', 'Dipanggil', 'Sedang dilayani', 'Selesai'].includes(queueStore.currentTicket.status) ? "Sudah Terverifikasi / Hadir" : "Check-in Langsung Sekarang" }}
          </button>
          <button class="button outline full-width" @click="openWhatsApp">
            Buka WhatsApp · kirim tiket
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
import { computed } from "vue";
import { useRouter } from "vue-router";
import { CheckCircle2, Download, QrCode, Sparkles } from "lucide-vue-next";
import Brand from "../../components/common/Brand.vue";
import { useQueueStore } from "../../stores/queue";
import { api } from "../../services/api";

const router = useRouter();
const queueStore = useQueueStore();

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
}

async function handleCheckinSimulation() {
  if (queueStore.currentTicket) {
    const token = queueStore.currentTicket.qr_token || queueStore.currentTicket.id;
    try {
      const res = await api.checkInByQR(token);
      if (res && res.success) {
        queueStore.changeStatus(queueStore.currentTicket.id, "Hadir");
        queueStore.setNotice(
          "✅ Check-in berhasil terverifikasi di server database! Antrean Anda siap dipanggil.",
        );
      } else {
        queueStore.changeStatus(queueStore.currentTicket.id, "Hadir");
        queueStore.setNotice(res?.message || "Check-in berhasil diproses.");
      }
    } catch {
      queueStore.changeStatus(queueStore.currentTicket.id, "Hadir");
      queueStore.setNotice("Check-in diterima secara lokal.");
    }
  }
}

function openWhatsApp() {
  const t = queueStore.currentTicket;
  if (!t) return;
  const phone = t.phone.replace(/^0/, "62");
  const msg = encodeURIComponent(
    `E-Tiket AntriSehat: ${t.number}, ${t.poli}, ${t.date}. Simpan tiket Anda.`,
  );
  window.open(`https://wa.me/${phone}?text=${msg}`, "_blank");
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
