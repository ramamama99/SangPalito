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
          class="login-error-alert"
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

        <button class="button primary full-width" type="submit" :disabled="loading" style="margin-top: 20px;">
          <span v-if="loading">Memverifikasi...</span>
          <span v-else class="login-btn-content">
            Masuk ke Dashboard <ArrowRight :size="18" />
          </span>
        </button>
      </form>
      <RouterLink to="/" class="back-link">← Kembali ke beranda</RouterLink>
    </div>
  </main>
</template>

<style scoped>
.login-error-alert {
  margin-bottom: 16px;
  padding: 10px 14px;
  border-radius: 10px;
  background-color: #fdf3f0;
  border: 1px solid #f7d6cf;
  color: var(--coral);
  font-size: 13px;
  font-weight: 600;
  text-align: left;
}

.login-btn-content {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
}

.back-link {
  display: block;
  text-align: center;
  margin: 22px auto 0;
  font-size: 14px;
  color: var(--muted);
  transition: color 0.15s ease;
  text-decoration: none;
}

.back-link:hover {
  color: var(--deep);
}
</style>

<script setup lang="ts">
import { ref, reactive } from "vue";
import { useRouter } from "vue-router";
import { ShieldCheck, ArrowRight } from "lucide-vue-next";
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
      const redirectPath = (router.currentRoute.value.query.redirect as string) || "/staff";
      router.push(redirectPath);
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
</style>

