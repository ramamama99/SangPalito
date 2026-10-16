<template>
  <header class="header" :class="{ 'is-scrolled': isScrolled }">
    <div class="nav-wrap" :class="{ 'is-capsule': isScrolled }">
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
        <RouterLink
          to="/live"
          :class="{ 'nav-active': route.path === '/live' }"
          @click="mobileMenuOpen = false"
        >
          <span class="live-dot" /> Antrean Live
        </RouterLink>
      </nav>
      <div class="nav-actions">
        <RouterLink to="/login" class="staff-link" @click="mobileMenuOpen = false">
          <Users :size="16" /> Login Staf
        </RouterLink>
        <RouterLink to="/register" class="button primary small" @click="mobileMenuOpen = false">
          Daftar Online <ArrowUpRight :size="17" />
        </RouterLink>
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

<style scoped>
.nav-wrap {
  top: 12px;
  height: 76px;
  width: calc(100% - 32px);
  max-width: 1180px;
  padding: 0 30px;
  border-radius: 16px;
  border: 1px solid var(--line);
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(16px);
  -webkit-backdrop-filter: blur(16px);
  box-shadow: 0 4px 20px -4px rgba(25, 90, 95, 0.07);
  transition: all 0.35s cubic-bezier(0.16, 1, 0.3, 1);
}

.nav-wrap.is-capsule {
  top: 10px;
  height: 64px;
  max-width: 1040px;
  width: calc(100% - 32px);
  padding: 0 24px;
  border-radius: 9999px;
  background: rgba(255, 255, 255, 0.88);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border: 1px solid rgba(220, 229, 223, 0.88);
  box-shadow: 0 10px 30px -8px rgba(25, 90, 95, 0.12), 0 3px 8px rgba(0, 0, 0, 0.02);
}

.nav-wrap.is-capsule :deep(.brand) {
  font-size: 22px;
  gap: 8px;
  transition: all 0.3s ease;
}

.nav-wrap.is-capsule :deep(.brand-mark) {
  width: 32px;
  height: 32px;
  border-radius: 9px;
  transition: all 0.3s ease;
}

.nav-wrap.is-capsule :deep(.brand-mark svg) {
  width: 22px;
  height: 22px;
}

.nav-wrap.is-capsule .nav {
  gap: 22px;
  transition: gap 0.3s ease;
}

.nav-wrap.is-capsule .nav button,
.nav-wrap.is-capsule .nav a {
  padding: 8px 0;
  transition: padding 0.3s ease;
}

.nav-wrap.is-capsule .nav-actions {
  gap: 16px;
  transition: gap 0.3s ease;
}

.nav-wrap.is-capsule .nav-actions .button.small {
  padding: 8px 16px;
  font-size: 13.5px;
}

@media (max-width: 900px) {
  .nav-wrap {
    top: 8px;
    height: 64px;
    width: calc(100% - 20px);
    padding: 0 16px;
    border-radius: 14px;
  }

  .nav-wrap .nav.open {
    top: 72px;
    border-radius: 16px;
    border: 1px solid var(--line);
    background: rgba(255, 255, 255, 0.98);
    backdrop-filter: blur(20px);
    -webkit-backdrop-filter: blur(20px);
    box-shadow: 0 16px 36px rgba(25, 90, 95, 0.14);
  }

  .nav-wrap.is-capsule {
    top: 8px;
    height: 58px;
    width: calc(100% - 20px);
    padding: 0 16px;
    border-radius: 50px;
  }

  .nav-wrap.is-capsule .nav.open {
    top: 64px;
    border-radius: 20px;
    border: 1px solid var(--line);
    background: rgba(255, 255, 255, 0.98);
    backdrop-filter: blur(20px);
    -webkit-backdrop-filter: blur(20px);
    box-shadow: 0 16px 36px rgba(25, 90, 95, 0.14);
  }
}
</style>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from "vue";
import { useRoute, useRouter, RouterLink } from "vue-router";
import { Users, ArrowUpRight, Menu, X } from "lucide-vue-next";
import Brand from "../common/Brand.vue";
import { useQueueStore } from "../../stores/queue";

const route = useRoute();
const router = useRouter();
const queueStore = useQueueStore();

const mobileMenuOpen = ref(false);
const activeSection = ref("home");
const navRef = ref<HTMLElement | null>(null);
const isScrolled = ref(false);

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
  activeSection.value = sectionId;
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
  isScrolled.value = window.scrollY > 24;

  if (!isHomeRoute.value) return;

  // Jika scroll sudah mentok di bagian paling bawah halaman
  const scrollPosition = window.innerHeight + Math.ceil(window.scrollY);
  const isAtBottom =
    scrollPosition >= document.documentElement.scrollHeight - 80;

  if (isAtBottom) {
    activeSection.value = "faq";
    return;
  }

  const sectionIds = ["layanan", "cara", "faq"];
  let current = "home";

  for (const id of sectionIds) {
    const el = document.getElementById(id);
    if (el && el.getBoundingClientRect().top <= 180) {
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
