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
          style="margin-bottom: 14px; padding: 10px 14px; border-radius: 8px; background-color: #fef2f2; border: 1px solid #fee2e2; color: #b91c1c; font-size: 13px; font-weight: 500; text-align: left;"
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
          <span v-else style="display: inline-flex; align-items: center; justify-content: center; gap: 8px;">
            Masuk ke Dashboard <ArrowRight :size="18" />
          </span>
        </button>
      </form>
      <RouterLink to="/" class="back-link">← Kembali ke beranda</RouterLink>
    </div>
  </main>
</template>

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
      } else {
        queueStore.role = "loket";
      }
      router.push("/staff");
      window.scrollTo(0, 0);
    } else {
      errorMessage.value = res.error?.message || res.message || "Email atau kata sandi tidak valid";
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
