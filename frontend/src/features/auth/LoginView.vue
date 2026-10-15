<template>
  <main class="login-page">
    <div class="login-card">
      <span class="service-icon">
        <ShieldCheck :size="30" />
      </span>
      <h2>Selamat datang kembali.</h2>
      <p>Masuk ke ruang kerja AntriSehat.</p>
      <form @submit.prevent="handleLogin">
        <label>
          Email staf
          <input type="email" required v-model="email" />
        </label>
        <label>
          Kata sandi
          <input type="password" required minlength="6" v-model="password" />
        </label>
        <button class="button primary full-width" type="submit">
          Masuk ke Dashboard <ArrowRight :size="18" />
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

const router = useRouter();
const email = ref("staf@antrisehat.id");
const password = ref("demo123");

function handleLogin() {
  router.push("/staff");
  window.scrollTo(0, 0);
}

function goToHome() {
  router.push("/");
  window.scrollTo(0, 0);
}
</script>
