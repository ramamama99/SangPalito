<template>
  <div v-if="queueStore.verify" class="modal-overlay" @click="closeModal">
    <div
      class="modal"
      role="dialog"
      aria-modal="true"
      aria-label="Verifikasi pasien"
      @click.stop
    >
      <button class="modal-close" @click="closeModal" aria-label="Tutup">
        <X />
      </button>
      <span class="service-icon">
        <ShieldCheck :size="30" />
      </span>
      <h2>Verifikasi identitas pasien</h2>
      <p>Cocokkan data dengan KTP dan kartu BPJS fisik.</p>
      <div class="verify-data">
        <div>
          <span>Nomor antrean</span>
          <strong>{{ queueStore.verify.number }}</strong>
        </div>
        <div>
          <span>Nama lengkap</span>
          <strong>{{ queueStore.verify.name }}</strong>
        </div>
        <div>
          <span>NIK</span>
          <strong>{{ queueStore.verify.nik }}</strong>
        </div>
        <div>
          <span>Penjamin</span>
          <strong>{{ queueStore.verify.insurance }}</strong>
        </div>
        <div>
          <span>Nomor BPJS</span>
          <strong>{{ queueStore.verify.bpjs || "—" }}</strong>
        </div>
      </div>
      <button class="button primary full-width" @click="confirmAttendance">
        <Check :size="18" /> Konfirmasi Hadir
      </button>
      <button class="button danger full-width" @click="rejectPatient">
        Tolak · Data Tidak Valid
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ShieldCheck, Check, X } from "lucide-vue-next";
import { useQueueStore } from "../../stores/queue";

const queueStore = useQueueStore();

function closeModal() {
  queueStore.verify = null;
}

function confirmAttendance() {
  if (queueStore.verify) {
    queueStore.changeStatus(queueStore.verify.id, "Hadir");
    queueStore.verify = null;
    queueStore.setNotice("Identitas terverifikasi. Pasien siap dilayani.");
  }
}

function rejectPatient() {
  if (queueStore.verify) {
    queueStore.changeStatus(queueStore.verify.id, "Ditolak");
    queueStore.verify = null;
    queueStore.setNotice(
      "Data pasien ditolak. Silakan lakukan pendaftaran ulang.",
    );
  }
}
</script>
