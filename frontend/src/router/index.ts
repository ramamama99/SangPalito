import { createRouter, createWebHistory } from "vue-router";
import HomeView from "../features/home/HomeView.vue";
import RegisterView from "../features/registration/RegisterView.vue";
import TicketView from "../features/ticket/TicketView.vue";
import LiveView from "../features/live/LiveView.vue";
import TvView from "../features/tv/TvView.vue";
import LoginView from "../features/auth/LoginView.vue";
import StaffView from "../features/staff/StaffView.vue";
import MedicalRecordView from "../features/medical-record/MedicalRecordView.vue";
import PatientScanner from "../features/scanner/PatientScanner.vue";

import { verifyStaffAuth, isJwtValid } from "../utils/auth";

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: "/",
      name: "home",
      component: HomeView,
    },
    {
      path: "/register",
      name: "register",
      component: RegisterView,
    },
    {
      path: "/ticket",
      name: "ticket",
      component: TicketView,
    },
    {
      path: "/live",
      name: "live",
      component: LiveView,
    },
    {
      path: "/tv",
      name: "tv",
      component: TvView,
      meta: { requiresAuth: true },
    },
    {
      path: "/login",
      name: "login",
      component: LoginView,
    },
    {
      path: "/staff",
      name: "staff",
      component: StaffView,
      meta: { requiresAuth: true },
    },
    {
      path: "/database-obat",
      name: "database-obat",
      component: StaffView,
      meta: { requiresAuth: true },
    },
    {
      path: "/database-pasien",
      name: "database-pasien",
      component: StaffView,
      meta: { requiresAuth: true },
    },
    {
      path: "/shift-dokter",
      name: "shift-dokter",
      component: StaffView,
      meta: { requiresAuth: true },
    },
    {
      path: "/manajemen-user",
      name: "manajemen-user",
      component: StaffView,
    },
    {
      path: "/audit-log",
      name: "audit-log",
      component: StaffView,
    },
    {
      path: "/analytics",
      name: "analytics",
      component: StaffView,
    },
    {
      path: "/medical-record",
      name: "medical-record",
      component: MedicalRecordView,
      meta: { requiresAuth: true },
    },
    {
      path: "/scan",
      name: "scan",
      component: PatientScanner,
      meta: { requiresAuth: true },
    },
    {
      path: "/qr-test",
      name: "qr-test",
      component: PatientScanner,
      meta: { requiresAuth: true },
    },
    {
      path: "/:pathMatch(.*)*",
      redirect: "/",
    },
  ],
  scrollBehavior(_to, _from, savedPosition) {
    if (savedPosition) {
      return savedPosition;
    } else {
      return { top: 0 };
    }
  },
});

router.beforeEach(async (to, _from, next) => {
  const requiresAuth = to.matched.some((record) => record.meta?.requiresAuth);

  if (requiresAuth) {
    const isValid = await verifyStaffAuth();
    if (!isValid) {
      next({
        path: "/login",
        query: { redirect: to.fullPath },
      });
      return;
    }
  }

  // Jika sudah terotentikasi dan mencoba buka /login, arahkan langsung ke /staff
  if (to.path === "/login") {
    const token = localStorage.getItem("antrisehat-token");
    if (isJwtValid(token)) {
      next({ path: "/staff" });
      return;
    }
  }

  next();
});

export default router;
