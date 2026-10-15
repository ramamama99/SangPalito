<template>
  <header class="header">
    <div class="nav-wrap">
      <Brand />
      <nav ref="navRef" :class="['nav', { open: mobileMenuOpen }]">
        <button
          :class="{ 'nav-active': activeSection === 'home' && isHomeRoute }"
          data-section="home"
          @click="navigateSection('home')"
        >
          Beranda
        </button>
        <button
          :class="{ 'nav-active': activeSection === 'layanan' && isHomeRoute }"
          data-section="layanan"
          @click="navigateSection('layanan')"
        >
          Layanan
        </button>
        <button
          :class="{ 'nav-active': activeSection === 'cara' && isHomeRoute }"
          data-section="cara"
          @click="navigateSection('cara')"
        >
          Cara Daftar
        </button>
        <button
          :class="{ 'nav-active': activeSection === 'faq' && isHomeRoute }"
          data-section="faq"
          @click="navigateSection('faq')"
        >
          FAQ
        </button>
        <button
          :class="{ 'nav-active': route.path === '/live' }"
          @click="goToRoute('/live')"
        >
          <span class="live-dot" /> Antrean Live
        </button>
      </nav>
      <div class="nav-actions">
        <button class="staff-link" @click="goToRoute('/login')">
          <Users :size="16" /> Login Staf
        </button>
        <button class="button primary small" @click="handleRegister">
          Daftar Online <ArrowUpRight :size="17" />
        </button>
        <button
          class="menu-button"
          @click="mobileMenuOpen = !mobileMenuOpen"
          aria-label="Buka menu"
        >
          <X v-if="mobileMenuOpen" />
          <Menu v-else />
        </button>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from "vue";
import { useRoute, useRouter } from "vue-router";
import { Users, ArrowUpRight, Menu, X } from "lucide-vue-next";
import Brand from "../common/Brand.vue";
import { useQueueStore } from "../../stores/queue";

const route = useRoute();
const router = useRouter();
const queueStore = useQueueStore();

const mobileMenuOpen = ref(false);
const activeSection = ref("home");
const navRef = ref<HTMLElement | null>(null);

const isHomeRoute = computed(() => route.path === "/");

function goToRoute(path: string) {
  mobileMenuOpen.value = false;
  router.push(path);
  window.scrollTo(0, 0);
}

function handleRegister() {
  mobileMenuOpen.value = false;
  queueStore.startRegistration("Poli Umum", false);
  router.push("/register");
  window.scrollTo(0, 0);
}

function navigateSection(sectionId: string) {
  mobileMenuOpen.value = false;
  if (route.path !== "/") {
    router.push("/").then(() => {
      setTimeout(() => {
        scrollToSection(sectionId);
      }, 100);
    });
  } else {
    scrollToSection(sectionId);
  }
}

function scrollToSection(sectionId: string) {
  if (sectionId === "home") {
    window.scrollTo({ top: 0, behavior: "smooth" });
    activeSection.value = "home";
    return;
  }
  const el = document.getElementById(sectionId);
  if (el) {
    el.scrollIntoView({ behavior: "smooth" });
  }
}

function onScroll() {
  if (!isHomeRoute.value) return;
  const sectionIds = ["layanan", "cara", "faq"];
  let current = "home";

  for (const id of sectionIds) {
    const el = document.getElementById(id);
    if (el && el.getBoundingClientRect().top <= 140) {
      current = id;
    }
  }
  activeSection.value = current;
}

onMounted(() => {
  window.addEventListener("scroll", onScroll, { passive: true });
  onScroll();
});

onUnmounted(() => {
  window.removeEventListener("scroll", onScroll);
});
</script>
