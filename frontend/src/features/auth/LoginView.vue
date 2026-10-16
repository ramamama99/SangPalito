<template>
  <main class="login-page">
    <div class="login-card">
      <span class="service-icon">
        <ShieldCheck :size="30" />
      </span>
      <h2>Selamat datang kembali.</h2>
      <p>Masuk ke ruang kerja AntriSehat.</p>
      <form novalidate @submit.prevent="handleLogin">
        <div
          v-if="errorMessage"
          style="margin-bottom: 14px; padding: 10px 14px; border-radius: 8px; background-color: #fef2f2; border: 1px solid #fee2e2; color: #b91c1c; font-size: 13px; font-weight: 500; text-align: left;"
        >
          {{ errorMessage }}
        </div>
        <label>
          <span>Email staf <span class="req">*</span></span>
          <input
            type="email"
            v-model="email"
            :class="{ 'input-error': errors.email }"
            @input="errors.email = ''"
            placeholder="staf@antrisehat.id"
          />
          <span v-if="errors.email" class="field-error-text">
            {{ errors.email }}
          </span>
        </label>
        <label>
          <span>Kata sandi <span class="req">*</span></span>
          <input
            type="password"
            v-model="password"
            :class="{ 'input-error': errors.password }"
            @input="errors.password = ''"
            placeholder="Minimal 6 karakter"
          />
          <span v-if="errors.password" class="field-error-text">
            {{ errors.password }}
          </span>
        </label>

        <div style="display: flex; justify-content: flex-end; margin-top: -6px; margin-bottom: 16px;">
          <button
            type="button"
            class="forgot-pwd-btn"
            @click="openForgotModal"
          >
            Lupa kata sandi?
          </button>
        </div>

        <button class="button primary full-width" type="submit" :disabled="loading">
          <span v-if="loading">Memverifikasi...</span>
          <span v-else style="display: inline-flex; align-items: center; justify-content: center; gap: 8px;">
            Masuk ke Dashboard <ArrowRight :size="18" />
          </span>
        </button>
      </form>
      <RouterLink to="/" class="back-link">← Kembali ke beranda</RouterLink>
    </div>

    <!-- Modal Lupa Kata Sandi -->
    <div v-if="showForgotModal" class="modal-backdrop" @click.self="closeForgotModal">
      <div class="modal-card">
        <div class="modal-header">
          <div style="display: flex; align-items: center; gap: 8px;">
            <KeyRound :size="20" style="color: #0f766e;" />
            <h3 style="font-size: 16px; font-weight: 700; color: #0f172a; margin: 0;">Pemulihan Kata Sandi</h3>
          </div>
          <button type="button" class="btn-close" @click="closeForgotModal">
            <X :size="18" />
          </button>
        </div>

        <div class="modal-body">
          <p style="font-size: 13px; color: #475569; margin-bottom: 16px;">
            Masukkan email akun Anda. Kami akan membuatkan tautan pemulihan kata sandi yang dapat dikirimkan langsung melalui <strong>WhatsApp</strong> atau <strong>Email</strong>.
          </p>

          <!-- Result State -->
          <div v-if="forgotResult" class="forgot-result-box">
            <div style="display: flex; align-items: center; gap: 8px; color: #15803d; font-weight: 600; font-size: 14px; margin-bottom: 8px;">
              <CheckCircle2 :size="18" />
              <span>Tautan Pemulihan Siap Digunakan!</span>
            </div>
            <p style="font-size: 12px; color: #475569; margin-bottom: 14px;">
              Tautan telah dibuat untuk email <strong>{{ forgotResult.email }}</strong>. Pilih metode pengiriman di bawah:
            </p>

            <div style="display: flex; flex-direction: column; gap: 10px;">
              <a
                v-if="forgotResult.whatsapp_url"
                :href="forgotResult.whatsapp_url"
                target="_blank"
                rel="noopener noreferrer"
                class="button whatsapp-btn"
              >
                <MessageSquare :size="16" /> Kirim / Buka di WhatsApp
              </a>

              <button
                type="button"
                class="button secondary"
                @click="copyResetLink(forgotResult.reset_url)"
              >
                <Copy :size="16" /> {{ copySuccess ? 'Tersalin ke Clipboard!' : 'Salin Tautan Email' }}
              </button>

              <button
                type="button"
                class="button outline"
                @click="goToResetPage(forgotResult.reset_token)"
              >
                <ExternalLink :size="16" /> Langsung Buka Halaman Reset
              </button>
            </div>
          </div>

          <!-- Input State -->
          <form v-else novalidate @submit.prevent="handleForgotPassword">
            <div
              v-if="forgotErrorMessage"
              style="margin-bottom: 12px; padding: 10px 12px; border-radius: 8px; background-color: #fef2f2; border: 1px solid #fee2e2; color: #b91c1c; font-size: 13px; font-weight: 500; text-align: left;"
            >
              {{ forgotErrorMessage }}
            </div>

            <label>
              <span>Email Akun Terdaftar <span class="req">*</span></span>
              <input
                type="email"
                v-model="forgotEmail"
                :class="{ 'input-error': forgotEmailError }"
                @input="forgotEmailError = ''"
                placeholder="nama@antrisehat.id"
              />
              <span v-if="forgotEmailError" class="field-error-text">
                {{ forgotEmailError }}
              </span>
            </label>

            <div style="display: flex; justify-content: flex-end; gap: 10px; margin-top: 20px;">
              <button type="button" class="button outline" @click="closeForgotModal">
                Batal
              </button>
              <button type="submit" class="button primary" :disabled="forgotLoading">
                <span v-if="forgotLoading">Memproses...</span>
                <span v-else>Buat Tautan Reset</span>
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  </main>
</template>

<script setup lang="ts">
import { ref, reactive } from "vue";
import { useRouter } from "vue-router";
import {
  ShieldCheck,
  ArrowRight,
  KeyRound,
  X,
  CheckCircle2,
  MessageSquare,
  Copy,
  ExternalLink,
} from "lucide-vue-next";
import { api } from "../../services/api";
import { useQueueStore } from "../../stores/queue";

const router = useRouter();
const queueStore = useQueueStore();

const email = ref("loket@antrisehat.id");
const password = ref("password123");
const loading = ref(false);
const errorMessage = ref("");
const errors = reactive({
  email: "",
  password: "",
});

// Forgot Password State
const showForgotModal = ref(false);
const forgotEmail = ref("");
const forgotEmailError = ref("");
const forgotLoading = ref(false);
const forgotErrorMessage = ref("");
const forgotResult = ref<{
  email: string;
  phone?: string;
  reset_token: string;
  reset_url: string;
  whatsapp_url?: string;
} | null>(null);
const copySuccess = ref(false);

function openForgotModal() {
  showForgotModal.value = true;
  forgotEmail.value = email.value || "";
  forgotEmailError.value = "";
  forgotErrorMessage.value = "";
  forgotResult.value = null;
  copySuccess.value = false;
}

function closeForgotModal() {
  showForgotModal.value = false;
  forgotResult.value = null;
}

async function handleForgotPassword() {
  forgotEmailError.value = "";
  forgotErrorMessage.value = "";

  if (!forgotEmail.value || !forgotEmail.value.trim()) {
    forgotEmailError.value = "Email akun wajib diisi";
    return;
  }
  if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(forgotEmail.value.trim())) {
    forgotEmailError.value = "Format email tidak valid (contoh: dokter@antrisehat.id)";
    return;
  }

  forgotLoading.value = true;
  try {
    const res = await api.forgotPassword(forgotEmail.value.trim());
    if (res.success && res.data) {
      forgotResult.value = res.data;
    } else {
      forgotErrorMessage.value =
        typeof res.error === "string"
          ? res.error
          : res.error?.message || res.message || "Email tidak ditemukan dalam sistem faskes.";
    }
  } catch (err) {
    forgotErrorMessage.value = "Gagal memproses permintaan reset kata sandi.";
  } finally {
    forgotLoading.value = false;
  }
}

function copyResetLink(link: string) {
  if (navigator.clipboard) {
    navigator.clipboard.writeText(link);
    copySuccess.value = true;
    setTimeout(() => {
      copySuccess.value = false;
    }, 2500);
  }
}

function goToResetPage(resetToken: string) {
  closeForgotModal();
  router.push(`/reset-password?token=${resetToken}`);
}

function validate(): boolean {
  let valid = true;
  errors.email = "";
  errors.password = "";

  if (!email.value || !email.value.trim()) {
    errors.email = "Email staf wajib diisi";
    valid = false;
  } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email.value.trim())) {
    errors.email = "Format email tidak valid (contoh: staf@antrisehat.id)";
    valid = false;
  }

  if (!password.value) {
    errors.password = "Kata sandi wajib diisi";
    valid = false;
  } else if (password.value.length < 6) {
    errors.password = "Kata sandi minimal 6 karakter";
    valid = false;
  }

  return valid;
}

async function handleLogin() {
  if (!validate()) {
    return;
  }

  loading.value = true;
  errorMessage.value = "";

  try {
    const res = await api.login(email.value.trim(), password.value);
    if (res.success && res.data?.access_token) {
      const user = res.data.user;
      if (user?.role === "medic") {
        queueStore.role = "medis";
      } else if (user?.role === "head") {
        queueStore.role = "analytics";
      } else if (user?.role === "admin") {
        queueStore.role = "user_management";
      } else {
        queueStore.role = "loket";
      }
      router.push("/staff");
      window.scrollTo(0, 0);
    } else {
      errorMessage.value =
        typeof res.error === "string"
          ? res.error
          : res.error?.message || res.message || "Email atau kata sandi tidak valid";
    }
  } catch (err: any) {
    errorMessage.value = "Gagal terhubung ke server backend";
  } finally {
    loading.value = false;
  }
}

function goToHome() {
  router.push("/");
  window.scrollTo(0, 0);
}
</script>

<style scoped>
.req {
  color: #dc2626;
  font-weight: 700;
  margin-left: 2px;
}

.input-error {
  border-color: #ef4444 !important;
  background-color: #fef2f2 !important;
  box-shadow: 0 0 0 3px rgba(239, 68, 68, 0.15) !important;
}

.field-error-text {
  display: block;
  font-size: 12px;
  font-weight: 600;
  color: #dc2626;
  margin-top: 5px;
  text-align: left;
}

.forgot-pwd-btn {
  background: none;
  border: none;
  color: #0f766e;
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  padding: 0;
  transition: color 0.15s;
}

.forgot-pwd-btn:hover {
  color: #115e59;
  text-decoration: underline;
}

.modal-backdrop {
  position: fixed;
  inset: 0;
  background-color: rgba(15, 23, 42, 0.6);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 16px;
}

.modal-card {
  background: white;
  border-radius: 16px;
  width: 100%;
  max-width: 480px;
  box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.15), 0 10px 10px -5px rgba(0, 0, 0, 0.04);
  overflow: hidden;
  animation: modalFadeIn 0.2s ease-out;
}

@keyframes modalFadeIn {
  from {
    opacity: 0;
    transform: scale(0.96);
  }
  to {
    opacity: 1;
    transform: scale(1);
  }
}

.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 18px 24px;
  border-bottom: 1px solid #e2e8f0;
}

.btn-close {
  background: none;
  border: none;
  color: #64748b;
  cursor: pointer;
  padding: 4px;
  border-radius: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.15s;
}

.btn-close:hover {
  background-color: #f1f5f9;
  color: #0f172a;
}

.modal-body {
  padding: 24px;
}

.forgot-result-box {
  background-color: #f0fdf4;
  border: 1px solid #bbf7d0;
  border-radius: 12px;
  padding: 18px;
  text-align: left;
}

.button.whatsapp-btn {
  background-color: #25d366;
  color: white;
  border: none;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  font-weight: 600;
  padding: 10px 16px;
  border-radius: 8px;
  text-decoration: none;
  transition: background-color 0.15s;
}

.button.whatsapp-btn:hover {
  background-color: #1ebc59;
}
</style>

