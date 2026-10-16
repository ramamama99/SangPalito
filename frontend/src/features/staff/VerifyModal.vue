<template>
  <Teleport to="body">
    <div
      v-if="queueStore.verify"
      class="verify-modal-backdrop"
      @click="closeModal"
      @keydown.esc="closeModal"
      tabindex="-1"
    >
      <div
        class="verify-modal-card"
        role="dialog"
        aria-modal="true"
        aria-labelledby="verify-modal-title"
        @click.stop
      >
        <!-- Modal Header -->
        <div class="verify-modal-header">
          <div class="header-left">
            <div class="shield-badge">
              <ShieldCheck :size="20" class="text-teal-700" />
            </div>
            <div>
              <h3 id="verify-modal-title" class="verify-modal-title">
                Verifikasi Identitas Pasien
              </h3>
              <p class="verify-modal-subtitle">
                Cocokkan data dengan KTP fisik & kartu asuransi/BPJS.
              </p>
            </div>
          </div>
          <button
            type="button"
            class="verify-modal-close"
            @click="closeModal"
            aria-label="Tutup modal"
          >
            <X :size="18" />
          </button>
        </div>

        <!-- Modal Body -->
        <div class="verify-modal-body">
          <!-- Queue Ticket Hero Banner -->
          <div class="queue-hero-card">
            <div class="queue-info">
              <span class="queue-label">NOMOR ANTREAN</span>
              <div class="queue-number">{{ queueStore.verify.number }}</div>
            </div>
            <div class="queue-meta">
              <span class="poli-pill">
                <Stethoscope :size="13" />
                {{ queueStore.verify.poli || "Poli Umum" }}
              </span>
              <span v-if="queueStore.verify.shift" class="shift-pill">
                Shift {{ queueStore.verify.shift }}
              </span>
            </div>
          </div>

          <!-- Patient Identity Card -->
          <div class="patient-id-card">
            <!-- Nama Pasien -->
            <div class="data-row">
              <span class="data-label">
                <User :size="15" /> Nama Lengkap
              </span>
              <span class="data-val patient-name">
                {{ queueStore.verify.name }}
              </span>
            </div>

            <!-- NIK -->
            <div class="data-row">
              <span class="data-label">
                <CreditCard :size="15" /> NIK (KTP)
              </span>
              <span class="nik-pill">
                {{ queueStore.verify.nik || "—" }}
              </span>
            </div>

            <!-- Penjamin / Asuransi -->
            <div class="data-row">
              <span class="data-label">
                <Building2 :size="15" /> Penjamin
              </span>
              <div>
                <span
                  v-if="isBpjs(queueStore.verify.insurance)"
                  class="badge-insurance bpjs"
                >
                  <CheckCircle2 :size="12" /> BPJS Kesehatan
                </span>
                <span
                  v-else-if="queueStore.verify.insurance === 'Asuransi Swasta'"
                  class="badge-insurance asuransi"
                >
                  Asuransi Swasta
                </span>
                <span v-else class="badge-insurance umum">
                  Umum (Mandiri)
                </span>
              </div>
            </div>

            <!-- Nomor Kartu BPJS -->
            <div
              v-if="isBpjs(queueStore.verify.insurance) || (queueStore.verify.bpjs && queueStore.verify.bpjs !== '—')"
              class="data-row"
            >
              <span class="data-label">
                <Hash :size="15" /> Nomor Kartu BPJS
              </span>
              <span class="bpjs-pill">
                {{ queueStore.verify.bpjs || "—" }}
              </span>
            </div>

            <!-- Dokter Pemeriksa (bila ada) -->
            <div v-if="queueStore.verify.doctorName" class="data-row">
              <span class="data-label">
                <Stethoscope :size="15" /> Dokter Pemeriksa
              </span>
              <span class="data-val text-slate-700">
                {{ queueStore.verify.doctorName }}
              </span>
            </div>
          </div>
        </div>

        <!-- Modal Actions Footer -->
        <div class="verify-modal-footer">
          <button
            type="button"
            class="btn-reject"
            @click="rejectPatient"
          >
            <XCircle :size="16" />
            <span>Tolak</span>
          </button>
          <button
            type="button"
            class="btn-confirm"
            @click="confirmAttendance"
          >
            <Check :size="16" :stroke-width="2.5" />
            <span>Konfirmasi Hadir</span>
          </button>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import {
  ShieldCheck,
  Check,
  X,
  XCircle,
  User,
  CreditCard,
  Building2,
  Stethoscope,
  Hash,
  CheckCircle2,
} from "lucide-vue-next";
import { useQueueStore } from "../../stores/queue";

const queueStore = useQueueStore();

function isBpjs(insurance?: string): boolean {
  if (!insurance) return false;
  return insurance.toUpperCase().includes("BPJS");
}

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

<style scoped>
.verify-modal-backdrop {
  position: fixed;
  inset: 0;
  z-index: 9999;
  background-color: rgba(15, 23, 42, 0.55);
  backdrop-filter: blur(5px);
  -webkit-backdrop-filter: blur(5px);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 16px;
  animation: backdropFade 0.2s ease-out;
}

@keyframes backdropFade {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

.verify-modal-card {
  background: #ffffff;
  border-radius: 20px;
  width: 100%;
  max-width: 480px;
  box-shadow:
    0 25px 50px -12px rgba(15, 23, 42, 0.25),
    0 0 0 1px rgba(226, 232, 240, 0.8);
  overflow: hidden;
  display: flex;
  flex-direction: column;
  animation: modalScaleUp 0.22s cubic-bezier(0.16, 1, 0.3, 1);
}

@keyframes modalScaleUp {
  from {
    opacity: 0;
    transform: scale(0.96) translateY(6px);
  }
  to {
    opacity: 1;
    transform: scale(1) translateY(0);
  }
}

/* Header */
.verify-modal-header {
  padding: 18px 22px;
  border-bottom: 1px solid #f1f5f9;
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: #ffffff;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.shield-badge {
  width: 40px;
  height: 40px;
  border-radius: 12px;
  background-color: #f0fdfa;
  border: 1px solid #ccfbf1;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #0f766e;
  flex-shrink: 0;
}

.verify-modal-title {
  font-family: "Manrope", sans-serif;
  font-size: 16px;
  font-weight: 700;
  color: #0f172a;
  margin: 0;
  line-height: 1.25;
}

.verify-modal-subtitle {
  font-size: 12px;
  color: #64748b;
  margin: 2px 0 0;
}

.verify-modal-close {
  background: transparent;
  border: none;
  width: 32px;
  height: 32px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #94a3b8;
  cursor: pointer;
  transition: all 0.15s ease;
}

.verify-modal-close:hover {
  background-color: #f1f5f9;
  color: #0f172a;
}

/* Body */
.verify-modal-body {
  padding: 20px 22px;
  display: flex;
  flex-direction: column;
  gap: 14px;
}

/* Queue Hero Card */
.queue-hero-card {
  background: linear-gradient(135deg, #f0fdf9 0%, #ecfdf5 55%, #f6f8f5 100%);
  border: 1px solid #c8e8df;
  border-radius: 14px;
  padding: 14px 18px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.queue-label {
  display: block;
  font-size: 10.5px;
  font-weight: 700;
  color: #0f766e;
  letter-spacing: 0.8px;
  text-transform: uppercase;
}

.queue-number {
  font-family: "JetBrains Mono", "Roboto Mono", monospace;
  font-size: 28px;
  font-weight: 800;
  color: #042f2e;
  letter-spacing: -0.5px;
  line-height: 1.1;
  margin-top: 2px;
}

.queue-meta {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 6px;
}

.poli-pill {
  background: #ffffff;
  color: #0f766e;
  border: 1px solid #99f6e4;
  padding: 4px 10px;
  border-radius: 999px;
  font-size: 11.5px;
  font-weight: 600;
  display: inline-flex;
  align-items: center;
  gap: 5px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.03);
}

.shift-pill {
  font-size: 11px;
  font-weight: 500;
  color: #64748b;
  text-transform: capitalize;
}

/* Patient Data Card */
.patient-id-card {
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  border-radius: 14px;
  padding: 6px 14px;
  display: flex;
  flex-direction: column;
}

.data-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 10px 0;
  border-bottom: 1px solid #edf2f7;
}

.data-row:last-child {
  border-bottom: none;
}

.data-label {
  display: inline-flex;
  align-items: center;
  gap: 7px;
  font-size: 12.5px;
  font-weight: 500;
  color: #64748b;
  flex-shrink: 0;
}

.data-val {
  font-size: 13.5px;
  font-weight: 600;
  text-align: right;
  word-break: break-word;
}

.patient-name {
  color: #0f172a;
  text-transform: uppercase;
  font-weight: 700;
  letter-spacing: 0.2px;
}

.nik-pill {
  font-family: "JetBrains Mono", "Roboto Mono", monospace;
  font-size: 13px;
  font-weight: 600;
  color: #1e293b;
  background: #ffffff;
  border: 1px solid #cbd5e1;
  padding: 2px 8px;
  border-radius: 6px;
  letter-spacing: 0.4px;
}

.bpjs-pill {
  font-family: "JetBrains Mono", "Roboto Mono", monospace;
  font-size: 12.5px;
  font-weight: 600;
  color: #0f766e;
  background: #ffffff;
  border: 1px solid #ccfbf1;
  padding: 2px 8px;
  border-radius: 6px;
}

.badge-insurance {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  font-size: 11.5px;
  font-weight: 700;
  padding: 3px 9px;
  border-radius: 6px;
}

.badge-insurance.bpjs {
  background-color: #ecfdf5;
  color: #047857;
  border: 1px solid #a7f3d0;
}

.badge-insurance.umum {
  background-color: #f1f5f9;
  color: #475569;
  border: 1px solid #cbd5e1;
}

.badge-insurance.asuransi {
  background-color: #eff6ff;
  color: #1d4ed8;
  border: 1px solid #bfdbfe;
}

/* Notice */
.verification-notice {
  background-color: #fefce8;
  border: 1px solid #fef08a;
  border-radius: 10px;
  padding: 10px 12px;
  display: flex;
  align-items: flex-start;
  gap: 9px;
}

.notice-icon {
  color: #d97706;
  flex-shrink: 0;
  margin-top: 1px;
}

.verification-notice p {
  font-size: 11.5px;
  color: #854d0e;
  line-height: 1.45;
  margin: 0;
}

/* Footer Actions */
.verify-modal-footer {
  padding: 14px 22px 18px;
  background: #ffffff;
  border-top: 1px solid #f1f5f9;
  display: flex;
  gap: 10px;
}

.btn-reject {
  flex: 1;
  height: 44px;
  border-radius: 12px;
  background-color: #fff1f2;
  border: 1px solid #fecdd3;
  color: #e11d48;
  font-size: 13px;
  font-weight: 600;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 7px;
  cursor: pointer;
  transition: all 0.15s ease;
}

.btn-reject:hover {
  background-color: #ffe4e6;
  border-color: #fda4af;
  color: #be123c;
}

.btn-reject:active {
  transform: scale(0.98);
}

.btn-confirm {
  flex: 1;
  height: 44px;
  border-radius: 12px;
  background-color: #0f766e;
  border: 1px solid #0f766e;
  color: #ffffff;
  font-size: 13px;
  font-weight: 700;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 7px;
  cursor: pointer;
  box-shadow: 0 4px 12px rgba(15, 118, 110, 0.22);
  transition: all 0.15s ease;
}

.btn-confirm:hover {
  background-color: #115e59;
  border-color: #115e59;
  box-shadow: 0 6px 16px rgba(15, 118, 110, 0.3);
  transform: translateY(-1px);
}

.btn-confirm:active {
  transform: scale(0.98);
}

@media (max-width: 480px) {
  .verify-modal-footer {
    flex-direction: column-reverse;
  }
  .btn-reject,
  .btn-confirm {
    width: 100%;
    flex: unset;
  }
}
</style>
