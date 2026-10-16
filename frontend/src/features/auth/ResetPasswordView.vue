<template>
  <main class="login-page">
    <div class="login-card">
      <span class="service-icon">
        <KeyRound :size="30" />
      </span>
      <h2>Atur Ulang Kata Sandi</h2>
      <p>Masukkan kata sandi baru untuk akun AntriSehat Anda.</p>

      <!-- Token Missing or Invalid State -->
      <div v-if="!token" class="alert-box error">
        <AlertCircle :size="18" />
        <div>
          <strong>Tautan Tidak Valid</strong>
          <p>Tautan reset kata sandi tidak memiliki token atau sudah kedaluwarsa.</p>
        </div>
      </div>

      <!-- Success State -->
      <div v-else-if="successMessage" class="alert-box success">
        <CheckCircle2 :size="20" />
        <div>
          <strong>Berhasil Diperbarui!</strong>
          <p>{{ successMessage }}</p>
        </div>
        <button
          type="button"
          class="button primary full-width"
          style="margin-top: 14px;"
          @click="router.push('/login')"
        >
          Masuk dengan Kata Sandi Baru
        </button>
      </div>

      <!-- Reset Password Form -->
      <form v-else novalidate @submit.prevent="handleReset">
        <div v-if="errorMessage" class="alert-box error">
          <AlertCircle :size="18" />
          <span>{{ errorMessage }}</span>
        </div>

        <label>
          <span>Kata Sandi Baru <span class="req">*</span></span>
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

        <label>
          <span>Konfirmasi Kata Sandi Baru <span class="req">*</span></span>
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

        <button class="button primary full-width" type="submit" :disabled="loading">
          <span v-if="loading">Menyimpan...</span>
          <span v-else style="display: inline-flex; align-items: center; justify-content: center; gap: 8px;">
            Simpan Kata Sandi Baru <ArrowRight :size="18" />
          </span>
        </button>
      </form>

      <RouterLink to="/login" class="back-link">← Kembali ke halaman login</RouterLink>
    </div>
  </main>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from "vue";
import { useRoute, useRouter } from "vue-router";
import { KeyRound, ArrowRight, AlertCircle, CheckCircle2 } from "lucide-vue-next";
import { api } from "../../services/api";

const route = useRoute();
const router = useRouter();

const token = ref("");
const newPassword = ref("");
const confirmPassword = ref("");
const loading = ref(false);
const errorMessage = ref("");
const successMessage = ref("");

const errors = reactive({
  newPassword: "",
  confirmPassword: "",
});

onMounted(() => {
  const queryToken = route.query.token as string;
  if (queryToken) {
    token.value = queryToken;
  }
});

function validate(): boolean {
  let valid = true;
  errors.newPassword = "";
  errors.confirmPassword = "";

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

  return valid;
}

async function handleReset() {
  if (!validate()) return;

  loading.value = true;
  errorMessage.value = "";

  try {
    const res = await api.resetPasswordWithToken(token.value, newPassword.value);
    if (res.success || res.code === "OK" || (res.data && !res.error)) {
      successMessage.value = "Kata sandi akun Anda telah berhasil diperbarui. Silakan masuk kembali.";
    } else {
      errorMessage.value =
        typeof res.error === "string"
          ? res.error
          : res.error?.message || res.message || "Gagal mengatur ulang kata sandi.";
    }
  } catch (err) {
    errorMessage.value = "Gagal terhubung ke server backend.";
  } finally {
    loading.value = false;
  }
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

.alert-box {
  display: flex;
  align-items: flex-start;
  gap: 10px;
  padding: 12px 14px;
  border-radius: 8px;
  font-size: 13px;
  text-align: left;
  margin-bottom: 14px;
}

.alert-box.error {
  background-color: #fef2f2;
  border: 1px solid #fee2e2;
  color: #b91c1c;
}

.alert-box.success {
  background-color: #f0fdf4;
  border: 1px solid #dcfce7;
  color: #15803d;
  flex-direction: column;
}

.back-link {
  display: inline-block;
  margin-top: 18px;
  font-size: 13px;
  color: #64748b;
  text-decoration: none;
  font-weight: 500;
}

.back-link:hover {
  color: #0f766e;
}
</style>
