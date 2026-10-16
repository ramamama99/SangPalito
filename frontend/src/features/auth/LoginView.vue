<template>
  <main class="login-page">
    <div class="login-card">
      <span class="service-icon">
        <ShieldCheck :size="30" />
      </span>
      <h2>Selamat datang kembali.</h2>
      <p>Masuk ke ruang kerja AntriSehat.</p>
      <form @submit.prevent="handleLogin">
        <div
          v-if="errorMessage"
          class="login-error-alert"
        >
          {{ errorMessage }}
        </div>
        <label>
          Email staf
          <input type="email" required v-model="email" />
        </label>
        <label>
          Kata sandi
          <input type="password" required minlength="6" v-model="password" />
        </label>
        <button class="button primary full-width" type="submit" :disabled="loading">
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
import { ref } from "vue";
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

async function handleLogin() {
  loading.value = true;
  errorMessage.value = "";

  try {
    const res = await api.login(email.value, password.value);
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
