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
      </div>
      <div class="service-grid">
        <button
          v-for="(s, i) in services"
          :key="s.name"
          class="service-card"
          @click="selectService(s.name)"
        >
          <div class="service-card-header">
            <span :class="['service-icon', `color-${i}`]">
              <component :is="getIcon(s.iconName)" :size="26" />
            </span>
            <span class="card-arrow">
              <ArrowUpRight :size="18" />
            </span>
          </div>
          <h3>{{ s.name }}</h3>
          <p>{{ s.desc }}</p>
        </button>
      </div>
    </div>
  </section>
</template>

<style scoped>
.service-card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
}

.card-arrow {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  border: 1px solid #d6e5de;
  display: grid;
  place-items: center;
  color: var(--deep);
  transition: all 0.2s ease;
}

.service-card:hover .card-arrow {
  background: var(--deep);
  color: white;
  border-color: var(--deep);
}
</style>

<script setup lang="ts">
import { useRouter } from "vue-router";
import {
  ArrowUpRight,
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
