<template>
  <section class="services-section" id="layanan">
    <div class="container">
      <div class="section-heading">
        <div>
          <div class="section-label"><span /> LAYANAN KESEHATAN</div>
          <h2>
            Perawatan yang tepat,
            <br />
            untuk setiap kebutuhan.
          </h2>
        </div>
        <p>
          Pilih layanan yang Anda butuhkan.
          <br />
          Kami siap mendampingi perjalanan sehat Anda.
        </p>
      </div>
      <div class="service-grid">
        <button
          v-for="(s, i) in services"
          :key="s.name"
          class="service-card"
          @click="selectService(s.name)"
        >
          <span :class="['service-icon', `color-${i}`]">
            <component :is="getIcon(s.iconName)" :size="26" />
          </span>
          <span class="availability"> <span class="live-dot" /> Tersedia </span>
          <h3>{{ s.name }}</h3>
          <p>{{ s.desc }}</p>
          <div class="card-bottom">
            <span> <Clock3 :size="14" /> 08.00 – 14.00 </span>
            <span class="card-arrow">
              <ArrowUpRight :size="21" />
            </span>
          </div>
        </button>
      </div>
      <p class="section-footnote">
        <MapPin :size="14" /> Klinik Buah Batu · Senin – Sabtu{" "}
        <span>*Target efisiensi pelayanan AntriSehat</span>
      </p>
    </div>
  </section>
</template>

<script setup lang="ts">
import { useRouter } from "vue-router";
import {
  Clock3,
  ArrowUpRight,
  MapPin,
  Stethoscope,
  Smile,
  Baby,
  Heart,
} from "lucide-vue-next";
import { services } from "../../constants/services";
import { useQueueStore } from "../../stores/queue";

const router = useRouter();
const queueStore = useQueueStore();

function getIcon(name: string) {
  switch (name) {
    case "Stethoscope":
      return Stethoscope;
    case "Smile":
      return Smile;
    case "Baby":
      return Baby;
    case "Heart":
      return Heart;
    default:
      return Stethoscope;
  }
}

function selectService(name: string) {
  queueStore.startRegistration(name, false);
  router.push("/register");
  window.scrollTo(0, 0);
}
</script>
