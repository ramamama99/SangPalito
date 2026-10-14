<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { Activity, CheckCircle2, Download, QrCode } from 'lucide-vue-next'
import { useQueueStore } from '../stores/queue'

const router = useRouter()
const route = useRoute()
const queueStore = useQueueStore()

onMounted(() => {
  const checkinId = route.query.checkin
  if (checkinId && typeof checkinId === 'string') {
    queueStore.ticket = checkinId
  }
})

const currentTicket = computed(() => queueStore.currentTicket)

const waitEstimate = computed(() => {
  if (!currentTicket.value) return '—'
  if (['Selesai', 'Dibatalkan', 'Ditolak'].includes(currentTicket.value.status)) {
    return '—'
  }
  const countBefore = queueStore.patients.filter(
    (p) =>
      p.poli === currentTicket.value!.poli &&
      p.date === currentTicket.value!.date &&
      !['Selesai', 'Dibatalkan', 'Ditolak'].includes(p.status) &&
      p.number < currentTicket.value!.number,
  ).length
  return `${countBefore * 8} menit`
})

const handleCheckInSimulation = () => {
  if (!currentTicket.value) return
  queueStore.changeStatus(currentTicket.value.id, 'Menunggu verifikasi')
  queueStore.setNotice('Check-in diterima. Silakan tunjukkan identitas kepada petugas.')
}

const handleWhatsAppShare = () => {
  if (!currentTicket.value) return
  const phone = currentTicket.value.phone.replace(/^0/, '62')
  const text = encodeURIComponent(
    `E-Tiket AntriSehat: ${currentTicket.value.number}, ${currentTicket.value.poli}, ${currentTicket.value.date}. Simpan tiket Anda.`,
  )
  window.open(`https://wa.me/${phone}?text=${text}`, '_blank')
}

const handleCancel = () => {
  if (!currentTicket.value) return
  if (window.confirm('Yakin ingin membatalkan antrean ini?')) {
    queueStore.changeStatus(currentTicket.value.id, 'Dibatalkan')
  }
}

const printTicket = () => {
  if (typeof window !== 'undefined') {
    window.print()
  }
}
</script>

<template>
  <main class="ticket-page container">
    <template v-if="currentTicket">
      <div class="success-banner">
        <CheckCircle2 :size="29" />
        <div>
          <strong>Pendaftaran berhasil. Jangan tutup halaman ini!</strong>
          <p>
            Wajib screenshot atau unduh tiket Anda. Pengiriman WhatsApp belum aktif pada prototype
            ini.
          </p>
        </div>
      </div>

      <div class="ticket-layout">
        <div class="ticket-card">
          <button class="brand" @click="router.push('/')" aria-label="AntriSehat beranda">
            <span class="brand-mark">
              <Activity :size="27" :stroke-width="2.6" />
            </span>
            <span>
              Antri<span class="teal">Sehat</span>
              <i />
            </span>
          </button>
          <p>E-TIKET ANTREAN</p>
          <h1>{{ currentTicket.number }}</h1>
          <span class="badge">{{ currentTicket.status }}</span>

          <img v-if="queueStore.qr" class="qr" :src="queueStore.qr" alt="QR check-in antrean" />

          <h3>{{ currentTicket.name }}</h3>
          <p>{{ currentTicket.poli }} · Puskesmas Sehat Bersama</p>

          <div class="ticket-details">
            <span>
              Tanggal kunjungan
              <strong>
                {{
                  new Date(currentTicket.date).toLocaleDateString('id-ID', {
                    dateStyle: 'long',
                  })
                }}
              </strong>
            </span>
            <span>
              Estimasi waktu tunggu
              <strong>{{ waitEstimate }}</strong>
            </span>
          </div>

          <button class="button primary full-width" @click="queueStore.downloadTicketCanvas">
            <Download :size="18" /> Unduh Tiket Gambar
          </button>
          <button class="button outline full-width" @click="printTicket">Cetak / Simpan PDF</button>
        </div>

        <div class="ticket-instructions">
          <h2>Selanjutnya, apa?</h2>
          <div
            v-for="([t, d], i) in [
              ['Simpan tiket Anda', 'Unduh atau screenshot tiket ini agar mudah ditemukan.'],
              ['Datang ke puskesmas', 'Bawa KTP dan kartu BPJS / asuransi sesuai penjamin.'],
              ['Scan QR di lokasi', 'Petugas akan mencocokkan identitas sebelum Anda dilayani.'],
            ]"
            :key="t"
            class="instruction"
          >
            <span>{{ i + 1 }}</span>
            <div>
              <h3>{{ t }}</h3>
              <p>{{ d }}</p>
            </div>
          </div>

          <button
            class="button primary full-width"
            :disabled="currentTicket.status !== 'Terdaftar'"
            @click="handleCheckInSimulation"
          >
            <QrCode :size="18" /> Simulasi Scan QR / Check-in
          </button>

          <button class="button outline full-width" @click="handleWhatsAppShare">
            Buka WhatsApp · kirim tiket
          </button>

          <p class="muted">
            QR dapat dipindai dengan kamera untuk membuka tiket. Check-in dan verifikasi adalah
            simulasi.
          </p>

          <button
            class="cancel-link"
            :disabled="['Selesai', 'Dibatalkan', 'Dipanggil'].includes(currentTicket.status)"
            @click="handleCancel"
          >
            Batalkan Antrean
          </button>

          <button class="back-link" @click="router.push('/')">Kembali ke beranda</button>
        </div>
      </div>
    </template>

    <div v-else class="empty-state">
      <QrCode :size="50" />
      <h2>Tiket tidak ditemukan</h2>
      <p>Periksa nomor antrean atau mulai pendaftaran baru.</p>
      <button class="button primary" @click="router.push('/register')">Daftar Online</button>
    </div>
  </main>
</template>
