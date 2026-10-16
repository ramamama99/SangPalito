<template>
  <Transition name="lookup-fade">
    <div
      v-if="queueStore.lookupOpen"
      class="lookup-modal-overlay"
      @click="closeModal"
      @keydown.esc="closeModal"
    >
      <div
        class="lookup-modal-card"
        role="dialog"
        aria-modal="true"
        aria-label="Cek Status Antrean"
        @click.stop
      >
        <!-- Close Button -->
        <button
          type="button"
          class="lookup-close-btn"
          @click="closeModal"
          aria-label="Tutup modal"
        >
          <X :size="20" />
        </button>

        <!-- Header -->
        <div class="lookup-header">
          <div class="lookup-icon-badge">
            <Search :size="26" class="text-teal-700" />
          </div>
          <h2 class="lookup-title">Cek Status & E-Tiket Antrean</h2>
          <p class="lookup-desc">
            Masukkan Nomor Antrean atau NIK Anda untuk melacak posisi antrean atau mencetak ulang tiket.
          </p>
        </div>

        <!-- Quick Access to Active Saved Ticket (if any) -->
        <div v-if="activeSavedTicket" class="active-ticket-shortcut" @click="goToSavedTicket">
          <div class="shortcut-info">
            <Ticket :size="18" class="shortcut-icon" />
            <div>
              <span class="shortcut-label">Tiket Anda Tersimpan:</span>
              <strong class="shortcut-val">{{ activeSavedTicket.number }} · {{ activeSavedTicket.poli }}</strong>
            </div>
          </div>
          <button type="button" class="btn-shortcut-action">
            Buka Tiket <ArrowRight :size="14" />
          </button>
        </div>

        <!-- Inline Error Alert when not found -->
        <div v-if="searchError" class="lookup-error-alert">
          <AlertCircle :size="20" class="error-alert-icon" />
          <div class="error-alert-content">
            <strong>Tiket Tidak Ditemukan</strong>
            <p>{{ searchError }}</p>
          </div>
        </div>

        <!-- Search Form (novalidate to prevent native tooltips) -->
        <form novalidate @submit.prevent="handleSearch" class="lookup-form">
          <!-- Field 1: Nomor Antrean -->
          <div class="form-group">
            <label class="input-label" for="lookup-number">
              <span>Nomor Antrean / Kode Booking</span>
              <span class="label-hint">Contoh: A-001</span>
            </label>
            <div class="input-wrap" :class="{ 'has-error': formErrors.number }">
              <Ticket :size="17" class="input-icon" />
              <input
                id="lookup-number"
                type="text"
                v-model="numberInput"
                placeholder="Masukkan nomor antrean (misal: A-001)"
                class="lookup-input uppercase-input"
                @input="clearErrors('number')"
                autocomplete="off"
              />
              <button
                v-if="numberInput"
                type="button"
                class="btn-clear-input"
                @click="numberInput = ''; clearErrors('number')"
                title="Hapus"
              >
                <X :size="14" />
              </button>
            </div>
            <!-- Red error text below box -->
            <span v-if="formErrors.number" class="field-error-text">
              {{ formErrors.number }}
            </span>
          </div>

          <!-- Divider -->
          <div class="form-divider">
            <span class="divider-line"></span>
            <span class="divider-text">ATAU</span>
            <span class="divider-line"></span>
          </div>

          <!-- Field 2: NIK -->
          <div class="form-group">
            <label class="input-label" for="lookup-nik">
              <span>Nomor Induk Kependudukan (NIK)</span>
              <span class="label-hint">16 Digit KTP / KK</span>
            </label>
            <div class="input-wrap" :class="{ 'has-error': formErrors.nik }">
              <CreditCard :size="17" class="input-icon" />
              <input
                id="lookup-nik"
                type="text"
                inputmode="numeric"
                maxlength="16"
                v-model="nikInput"
                placeholder="Contoh: 3273011201900001"
                class="lookup-input font-mono"
                @input="handleNikInput"
                autocomplete="off"
              />
              <button
                v-if="nikInput"
                type="button"
                class="btn-clear-input"
                @click="nikInput = ''; clearErrors('nik')"
                title="Hapus"
              >
                <X :size="14" />
              </button>
            </div>
            <!-- Red error text below box -->
            <span v-if="formErrors.nik" class="field-error-text">
              {{ formErrors.nik }}
            </span>
          </div>

          <!-- General Form Error Message -->
          <span v-if="formErrors.general" class="field-error-text text-center block" style="margin-top: 4px;">
            {{ formErrors.general }}
          </span>

          <!-- Submit Button -->
          <button
            type="submit"
            class="button primary full-width lookup-submit-btn"
            :disabled="isSearching"
          >
            <Loader2 v-if="isSearching" :size="18" class="animate-spin" />
            <span v-if="isSearching">Mencari Tiket...</span>
            <template v-else>
              <span>Temukan Tiket Saya</span>
              <ArrowRight :size="18" />
            </template>
          </button>

          <!-- Help Note -->
          <div class="lookup-help-note">
            <Clock :size="13" class="text-teal-700 flex-shrink-0" />
            <span>Pencarian otomatis memeriksa antrean aktif hari ini untuk semua poliklinik.</span>
          </div>
        </form>
      </div>
    </div>
  </Transition>
</template>

<script setup lang="ts">
import { ref, reactive, computed } from "vue";
import { useRouter } from "vue-router";
import {
  Search,
  ArrowRight,
  X,
  Ticket,
  CreditCard,
  AlertCircle,
  Clock,
  Loader2,
} from "lucide-vue-next";
import { useQueueStore } from "../../stores/queue";
import { api } from "../../services/api";
import type { Patient } from "../../types";

const router = useRouter();
const queueStore = useQueueStore();

const nikInput = ref("");
const numberInput = ref("");
const isSearching = ref(false);
const searchError = ref("");

const formErrors = reactive({
  number: "",
  nik: "",
  general: "",
});

// Check if user has an active ticket saved in store
const activeSavedTicket = computed(() => {
  const tId = queueStore.ticket;
  if (!tId) return null;
  return queueStore.patients.find((p) => p.id === tId || p.qr_token === tId) || null;
});

function goToSavedTicket() {
  if (activeSavedTicket.value) {
    queueStore.ticket = activeSavedTicket.value.id;
    queueStore.lookupOpen = false;
    router.push("/ticket");
    window.scrollTo(0, 0);
  }
}

function closeModal() {
  queueStore.lookupOpen = false;
  clearAllErrors();
  searchError.value = "";
}

function clearErrors(field?: "number" | "nik") {
  if (field) {
    formErrors[field] = "";
  }
  formErrors.general = "";
  searchError.value = "";
}

function clearAllErrors() {
  formErrors.number = "";
  formErrors.nik = "";
  formErrors.general = "";
}

function handleNikInput(e: Event) {
  const target = e.target as HTMLInputElement;
  nikInput.value = target.value.replace(/[^0-9]/g, "");
  clearErrors("nik");
}

async function handleSearch() {
  clearAllErrors();
  searchError.value = "";

  const numVal = numberInput.value.trim().toUpperCase();
  const nikVal = nikInput.value.trim();

  // 1. Validation: At least one field must be filled
  if (!numVal && !nikVal) {
    formErrors.general = "Silakan masukkan Nomor Antrean atau NIK Anda.";
    formErrors.number = "Harap masukkan nomor antrean (contoh: A-001)";
    formErrors.nik = "Atau masukkan 16 digit NIK Anda";
    return;
  }

  // 2. Validation: If NIK is filled, must be 16 digits
  if (nikVal && nikVal.length !== 16) {
    formErrors.nik = `NIK harus terdiri dari 16 digit angka (saat ini ${nikVal.length} digit).`;
    return;
  }

  isSearching.value = true;

  try {
    const norm = (str: string) => str.replace(/[^A-Z0-9]/gi, "").toUpperCase();

    // 1. Search in local store patients
    let patient: Patient | undefined;

    if (numVal && nikVal) {
      // Both provided: verify both match
      patient = queueStore.patients.find((p) => {
        const matchNik = p.nik === nikVal;
        const matchNum =
          norm(p.number) === norm(numVal) ||
          p.number.toUpperCase() === numVal ||
          p.qr_token === numVal ||
          p.id === numVal;
        return matchNik && matchNum;
      });
    } else if (numVal) {
      // Only number provided
      patient = queueStore.patients.find((p) => {
        return (
          norm(p.number) === norm(numVal) ||
          p.number.toUpperCase() === numVal ||
          p.qr_token === numVal ||
          p.id === numVal
        );
      });
    } else if (nikVal) {
      // Only NIK provided
      patient = queueStore.patients.find((p) => p.nik === nikVal);
    }

    // 2. If not found locally, try fetching fresh queues from backend API
    if (!patient) {
      try {
        const backendQueues = await api.getQueuesToday();
        if (backendQueues && backendQueues.length > 0) {
          const matchApi = backendQueues.find((q) => {
            const matchNik = nikVal ? q.patient?.medical_record_no === nikVal : true;
            const matchNum = numVal ? norm(q.queue_code) === norm(numVal) : true;
            return matchNik && matchNum;
          });

          if (matchApi) {
            patient = {
              id: matchApi.id,
              number: matchApi.queue_code,
              name: matchApi.patient?.full_name || "Pasien",
              nik: matchApi.patient?.medical_record_no || nikVal,
              phone: "",
              birth: "",
              gender: "Laki-laki",
              insurance: "Umum",
              bpjs: "",
              date: matchApi.queue_date || new Date().toISOString().split("T")[0],
              poli: matchApi.polyclinic?.name || "Poli Umum",
              complaint: "",
              allergy: "",
              symptoms: [],
              status: matchApi.status,
              source: "Online",
            };
            // Add to store so TicketView can display it
            queueStore.patients.push(patient);
          }
        }
      } catch (err) {
        console.warn("Backend queue search failed:", err);
      }
    }

    if (patient) {
      queueStore.ticket = patient.id;
      queueStore.lookupOpen = false;
      router.push("/ticket");
      window.scrollTo(0, 0);
      return;
    }

    // Not found
    searchError.value =
      "Tiket antrean tidak ditemukan. Mohon periksa kembali nomor antrean atau NIK yang Anda masukkan, dan pastikan sudah terdaftar.";
  } catch (err) {
    searchError.value = "Terjadi kendala saat memeriksa tiket antrean. Silakan coba lagi.";
  } finally {
    isSearching.value = false;
  }
}
</script>

<style scoped>
/* Overlay with backdrop blur */
.lookup-modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(15, 23, 42, 0.55);
  backdrop-filter: blur(8px);
  -webkit-backdrop-filter: blur(8px);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 16px;
  z-index: 999;
}

/* Modal Card */
.lookup-modal-card {
  position: relative;
  background: #ffffff;
  border-radius: 20px;
  max-width: 480px;
  width: 100%;
  padding: 32px 28px;
  box-shadow: 0 25px 50px -12px rgba(15, 23, 42, 0.25), 0 0 0 1px rgba(226, 232, 240, 0.8);
  border: 1px solid #e2e8f0;
  max-height: 92vh;
  overflow-y: auto;
}

/* Close button */
.lookup-close-btn {
  position: absolute;
  top: 18px;
  right: 18px;
  width: 32px;
  height: 32px;
  border-radius: 8px;
  background: #f1f5f9;
  border: none;
  color: #64748b;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.15s ease;
}

.lookup-close-btn:hover {
  background: #fee2e2;
  color: #ef4444;
}

/* Header */
.lookup-header {
  text-align: center;
  margin-bottom: 22px;
}

.lookup-icon-badge {
  width: 56px;
  height: 56px;
  border-radius: 16px;
  background: #f0fdfa;
  border: 1px solid #ccfbf1;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto 14px auto;
  box-shadow: 0 4px 12px rgba(15, 118, 110, 0.12);
}

.lookup-title {
  font-size: 20px;
  font-weight: 800;
  color: #0f172a;
  margin: 0 0 6px 0;
  line-height: 1.3;
}

.lookup-desc {
  font-size: 13px;
  color: #64748b;
  margin: 0;
  line-height: 1.45;
}

/* Active Saved Ticket Shortcut */
.active-ticket-shortcut {
  background: #f0fdfa;
  border: 1px solid #99f6e4;
  border-radius: 12px;
  padding: 10px 14px;
  margin-bottom: 18px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.active-ticket-shortcut:hover {
  background: #ccfbf1;
  transform: translateY(-1px);
  box-shadow: 0 4px 10px rgba(15, 118, 110, 0.12);
}

.shortcut-info {
  display: flex;
  align-items: center;
  gap: 10px;
}

.shortcut-icon {
  color: #0f766e;
  flex-shrink: 0;
}

.shortcut-label {
  display: block;
  font-size: 11px;
  font-weight: 600;
  color: #0f766e;
}

.shortcut-val {
  display: block;
  font-size: 13.5px;
  font-weight: 800;
  color: #134e4a;
}

.btn-shortcut-action {
  background: #0f766e;
  color: white;
  border: none;
  font-size: 11.5px;
  font-weight: 700;
  padding: 6px 12px;
  border-radius: 8px;
  display: inline-flex;
  align-items: center;
  gap: 5px;
  cursor: pointer;
  flex-shrink: 0;
  transition: background 0.15s ease;
}

.btn-shortcut-action:hover {
  background: #115e59;
}

/* Inline Error Alert */
.lookup-error-alert {
  background: #fff1f2;
  border: 1px solid #fecaca;
  border-radius: 12px;
  padding: 12px 14px;
  margin-bottom: 18px;
  display: flex;
  align-items: flex-start;
  gap: 10px;
}

.error-alert-icon {
  color: #e11d48;
  margin-top: 2px;
}

.error-alert-content strong {
  display: block;
  font-size: 13px;
  color: #9f1239;
  font-weight: 700;
  margin-bottom: 2px;
}

.error-alert-content p {
  font-size: 12px;
  color: #be123c;
  margin: 0;
  line-height: 1.4;
}

/* Form Styles */
.lookup-form {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.form-group {
  display: flex;
  flex-direction: column;
}

.input-label {
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-size: 12.5px;
  font-weight: 700;
  color: #334155;
  margin-bottom: 6px;
}

.label-hint {
  font-size: 11px;
  font-weight: 500;
  color: #64748b;
}

.input-wrap {
  position: relative;
  display: flex;
  align-items: center;
  background: #f8fafc;
  border: 1.5px solid #cbd5e1;
  border-radius: 10px;
  transition: all 0.15s ease;
}

.input-wrap:focus-within {
  background: #ffffff;
  border-color: #0f766e;
  box-shadow: 0 0 0 3px rgba(15, 118, 110, 0.15);
}

.input-wrap.has-error {
  border-color: #ef4444;
  background: #fffafa;
}

.input-wrap.has-error:focus-within {
  box-shadow: 0 0 0 3px rgba(239, 68, 68, 0.15);
}

.input-icon {
  position: absolute;
  left: 12px;
  color: #94a3b8;
  pointer-events: none;
}

.input-wrap:focus-within .input-icon {
  color: #0f766e;
}

.lookup-input {
  width: 100%;
  border: none;
  background: transparent;
  padding: 11px 36px 11px 38px;
  font-size: 13.5px;
  color: #0f172a;
  outline: none;
}

.lookup-input::placeholder {
  color: #94a3b8;
  font-size: 13px;
}

.uppercase-input {
  text-transform: uppercase;
  font-weight: 600;
  letter-spacing: 0.03em;
}

.btn-clear-input {
  position: absolute;
  right: 10px;
  width: 20px;
  height: 20px;
  border-radius: 50%;
  background: #e2e8f0;
  border: none;
  color: #64748b;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
}

.btn-clear-input:hover {
  background: #cbd5e1;
  color: #1e293b;
}

/* Red error text below box */
.field-error-text {
  font-size: 11.5px;
  font-weight: 600;
  color: #dc2626;
  margin-top: 5px;
  display: flex;
  align-items: center;
  gap: 4px;
}

/* Divider */
.form-divider {
  display: flex;
  align-items: center;
  margin: 4px 0;
  gap: 12px;
}

.divider-line {
  flex: 1;
  height: 1px;
  background: #e2e8f0;
}

.divider-text {
  font-size: 11px;
  font-weight: 700;
  color: #94a3b8;
  letter-spacing: 0.05em;
}

/* Submit button */
.lookup-submit-btn {
  margin-top: 6px;
  padding: 12px 20px;
  font-size: 14px;
  font-weight: 700;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  background: #0f766e;
  border: none;
  color: white;
  cursor: pointer;
  box-shadow: 0 4px 12px rgba(15, 118, 110, 0.25);
  transition: all 0.2s ease;
}

.lookup-submit-btn:hover:not(:disabled) {
  background: #115e59;
  transform: translateY(-1px);
  box-shadow: 0 6px 16px rgba(15, 118, 110, 0.3);
}

.lookup-submit-btn:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

/* Help note */
.lookup-help-note {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 11.5px;
  color: #64748b;
  line-height: 1.4;
  margin-top: 2px;
}

/* Transitions */
.lookup-fade-enter-active,
.lookup-fade-leave-active {
  transition: all 0.2s cubic-bezier(0.16, 1, 0.3, 1);
}

.lookup-fade-enter-from,
.lookup-fade-leave-to {
  opacity: 0;
  transform: scale(0.97);
}
</style>
