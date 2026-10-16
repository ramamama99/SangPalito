<template>
  <div v-if="isOpen" class="modal-backdrop">
    <div class="modal-card first-login-card">
      <div class="first-login-header">
        <div class="shield-icon-badge">
          <ShieldAlert :size="28" class="text-teal-700" />
        </div>
        <div>
          <h3 class="first-login-title">Aktivasi Akun Dokter (Login Pertama Kali)</h3>
          <p class="first-login-subtitle">
            Demi keamanan data medis dan rekam medis pasien, Anda wajib mengganti kata sandi awal dan mengonfirmasi nomor kontak Anda.
          </p>
        </div>
      </div>

      <form novalidate @submit.prevent="handleSubmit">
        <div class="modal-body space-y-4">
          <div v-if="generalError" class="alert-box error">
            <AlertCircle :size="16" />
            <span>{{ generalError }}</span>
          </div>

          <!-- New Password -->
          <label class="form-field">
            <span class="field-label">Kata Sandi Baru <span class="req">*</span></span>
            <input
              type="password"
              v-model="newPassword"
              :class="{ 'input-error': errors.newPassword }"
              @input="errors.newPassword = ''"
              placeholder="Minimal 6 karakter"
            />
            <span v-if="errors.newPassword" class="field-error-text">
              {{ errors.newPassword }}
            </span>
          </label>

          <!-- Confirm New Password -->
          <label class="form-field">
            <span class="field-label">Konfirmasi Kata Sandi Baru <span class="req">*</span></span>
            <input
              type="password"
              v-model="confirmPassword"
              :class="{ 'input-error': errors.confirmPassword }"
              @input="errors.confirmPassword = ''"
              placeholder="Ulangi kata sandi baru"
            />
            <span v-if="errors.confirmPassword" class="field-error-text">
              {{ errors.confirmPassword }}
            </span>
          </label>

          <!-- Phone Number -->
          <label class="form-field">
            <span class="field-label">Nomor WhatsApp Aktif <span class="req">*</span></span>
            <input
              type="tel"
              v-model="phone"
              :class="{ 'input-error': errors.phone }"
              @input="errors.phone = ''"
              placeholder="Contoh: 081234567890"
            />
            <span v-if="errors.phone" class="field-error-text">
              {{ errors.phone }}
            </span>
            <span class="field-hint">Digunakan untuk menerima notifikasi sistem & pemulihan kata sandi via WhatsApp.</span>
          </label>

          <!-- STR No (Optional) -->
          <label class="form-field">
            <span class="field-label">Nomor Surat Tanda Registrasi (STR)</span>
            <input
              type="text"
              v-model="strNo"
              placeholder="Contoh: 31.1.1.100.1.22.123456"
            />
          </label>
        </div>

        <div class="first-login-footer">
          <button type="submit" class="button primary full-width" :disabled="loading">
            <span v-if="loading">Memperbarui Akun...</span>
            <span v-else style="display: inline-flex; align-items: center; justify-content: center; gap: 8px;">
              Simpan & Mulai Melayani Pasien <ArrowRight :size="16" />
            </span>
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from "vue";
import { ShieldAlert, AlertCircle, ArrowRight } from "lucide-vue-next";
import { api } from "../../services/api";

const emit = defineEmits(["completed"]);

const isOpen = ref(false);
const newPassword = ref("");
const confirmPassword = ref("");
const phone = ref("");
const strNo = ref("");
const loading = ref(false);
const generalError = ref("");

const errors = reactive({
  newPassword: "",
  confirmPassword: "",
  phone: "",
});

onMounted(() => {
  const user = api.getCurrentUser();
  if (user && user.is_first_login) {
    isOpen.value = true;
    if (user.phone) {
      phone.value = user.phone;
    }
  }
});

function validate(): boolean {
  let valid = true;
  errors.newPassword = "";
  errors.confirmPassword = "";
  errors.phone = "";
  generalError.value = "";

  if (!newPassword.value) {
    errors.newPassword = "Kata sandi baru wajib diisi";
    valid = false;
  } else if (newPassword.value.length < 6) {
    errors.newPassword = "Kata sandi minimal 6 karakter";
    valid = false;
  }

  if (!confirmPassword.value) {
    errors.confirmPassword = "Konfirmasi kata sandi wajib diisi";
    valid = false;
  } else if (confirmPassword.value !== newPassword.value) {
    errors.confirmPassword = "Konfirmasi kata sandi tidak cocok";
    valid = false;
  }

  if (!phone.value || !phone.value.trim()) {
    errors.phone = "Nomor WhatsApp aktif wajib diisi";
    valid = false;
  } else if (phone.value.trim().length < 9) {
    errors.phone = "Nomor WhatsApp minimal 9 digit";
    valid = false;
  }

  return valid;
}

async function handleSubmit() {
  if (!validate()) return;

  loading.value = true;
  generalError.value = "";

  try {
    const res = await api.completeFirstLogin({
      new_password: newPassword.value,
      phone: phone.value.trim(),
      str_no: strNo.value.trim() || undefined,
    });

    if (res.success || res.code === "OK" || (res.data && !res.error)) {
      // Update local storage user
      const user = api.getCurrentUser();
      if (user) {
        user.is_first_login = false;
        user.phone = phone.value.trim();
        localStorage.setItem("antrisehat-user", JSON.stringify(user));
      }
      isOpen.value = false;
      emit("completed");
    } else {
      generalError.value =
        typeof res.error === "string"
          ? res.error
          : res.error?.message || res.message || "Gagal memperbarui akun.";
    }
  } catch (err) {
    generalError.value = "Gagal terhubung ke server backend.";
  } finally {
    loading.value = false;
  }
}
</script>

<style scoped>
.modal-backdrop {
  position: fixed;
  inset: 0;
  background-color: rgba(15, 23, 42, 0.75);
  backdrop-filter: blur(6px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 2000;
  padding: 16px;
}

.first-login-card {
  background: white;
  border-radius: 20px;
  width: 100%;
  max-width: 480px;
  box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.25);
  overflow: hidden;
  animation: modalFadeIn 0.25s ease-out;
}

@keyframes modalFadeIn {
  from {
    opacity: 0;
    transform: scale(0.95);
  }
  to {
    opacity: 1;
    transform: scale(1);
  }
}

.first-login-header {
  display: flex;
  align-items: flex-start;
  gap: 14px;
  padding: 24px 24px 18px;
  border-bottom: 1px solid #f1f5f9;
}

.shield-icon-badge {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  background-color: #f0fdfa;
  border: 1px solid #ccfbf1;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.first-login-title {
  font-size: 16px;
  font-weight: 700;
  color: #0f172a;
  margin: 0 0 4px;
}

.first-login-subtitle {
  font-size: 12.5px;
  color: #64748b;
  margin: 0;
  line-height: 1.4;
}

.modal-body {
  padding: 20px 24px;
}

.first-login-footer {
  padding: 16px 24px 24px;
}

.form-field {
  display: flex;
  flex-direction: column;
  gap: 4px;
  text-align: left;
}

.field-label {
  font-size: 12.5px;
  font-weight: 600;
  color: #334155;
}

.form-field input {
  padding: 10px 12px;
  border-radius: 8px;
  border: 1px solid #cbd5e1;
  font-size: 13.5px;
  outline: none;
}

.form-field input:focus {
  border-color: #0f766e;
  box-shadow: 0 0 0 2px rgba(15, 118, 110, 0.15);
}

.field-hint {
  font-size: 11px;
  color: #64748b;
  margin-top: 2px;
}

.req {
  color: #dc2626;
  font-weight: 700;
  margin-left: 2px;
}

.input-error {
  border-color: #ef4444 !important;
  background-color: #fef2f2 !important;
}

.field-error-text {
  display: block;
  font-size: 11.5px;
  font-weight: 600;
  color: #dc2626;
  margin-top: 3px;
}

.alert-box {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 14px;
  border-radius: 8px;
  font-size: 13px;
  font-weight: 500;
  margin-bottom: 12px;
}

.alert-box.error {
  background-color: #fef2f2;
  border: 1px solid #fee2e2;
  color: #b91c1c;
}
</style>
