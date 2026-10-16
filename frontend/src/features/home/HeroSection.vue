<template>
  <div>
    <main class="hero container">
      <div class="hero-copy">
        <h1>
          Lebih sedikit
          <br />
          menunggu.
          <br />
          <span>
            Lebih banyak
            <br class="mobile-break" />
            hidup sehat.
          </span>
        </h1>
        <p>
          Daftar dari rumah, datang saat giliran Anda.
          <br class="desktop-break" />
          Pelayanan Kesehatan kini lebih mudah, cepat,
          <br class="desktop-break" />
          dan nyaman untuk semua.
        </p>
        <div class="hero-actions">
          <RouterLink
            to="/register"
            class="button primary"
            @click="queueStore.startRegistration('Poli Umum', false)"
          >
            Ambil Antrean
            <span class="arrow-disc">
              <ArrowUpRight :size="20" />
            </span>
          </RouterLink>
          <button class="button outline" @click="queueStore.lookupOpen = true">
            <Search :size="17" /> Cek Antrean Saya
          </button>
        </div>
      </div>
      <div class="hero-visual">
        <div class="photo-frame">
          <img
            src="/images/doctor.jpg"
            alt="Dokter siap melayani kesehatan Anda"
          />
          <div class="photo-tint" />
          <span class="decor-plus">+</span>
        </div>
        <div class="floating-benefits">
          <div>
            <span>
              <ShieldCheck :size="20" />
            </span>
            <p>Data aman<small>& terlindungi</small></p>
          </div>
          <div>
            <span>
              <Clock3 :size="20" />
            </span>
            <p>Tanpa antre<small>lebih lama</small></p>
          </div>
        </div>
        <div class="floating-queue">
          <div class="queue-icon">
            <CalendarDays :size="24" />
          </div>
          <div>
            <small>Antrean Anda, dalam genggaman</small>
            <strong>Daftar kurang dari 2 menit</strong>
            <span> <span class="live-dot" /> Mudah. Cepat. Tanpa login. </span>
          </div>
          <CheckCircle2 :size="22" />
        </div>
      </div>
    </main>
    <div class="trust-strip container">
      <span> <ShieldCheck :size="18" /> Data pasien terlindungi </span>
      <i />
      <span> <Clock3 :size="18" /> Waktu tunggu 30% lebih singkat* </span>
      <i />
      <span>
        <CheckCircle2 :size="18" /> Terintegrasi dengan layanan BPJS
      </span>
      <i />
      <span> <Heart :size="18" /> Ramah untuk semua usia </span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from "vue-router";
import {
  ArrowUpRight,
  Search,
  ShieldCheck,
  Clock3,
  CalendarDays,
  CheckCircle2,
  Heart,
} from "lucide-vue-next";
import { useQueueStore } from "../../stores/queue";

const router = useRouter();
const queueStore = useQueueStore();

function handleStart() {
  queueStore.startRegistration("Poli Umum", false);
  router.push("/register");
  window.scrollTo(0, 0);
}
</script>
