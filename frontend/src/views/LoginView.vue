<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { ShieldCheck, ArrowRight } from 'lucide-vue-next'

const router = useRouter()
const email = ref('staf@antrisehat.id')
const password = ref('demo123')

const handleLogin = () => {
  router.push('/staff')
}
</script>

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

      <button class="back-link" @click="router.push('/')">← Kembali ke beranda</button>
    </div>
  </main>
</template>
