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
    },
    {
      path: "/medical-record",
      name: "medical-record",
      component: MedicalRecordView,
    },
    {
      path: "/scan",
      name: "scan",
      component: PatientScanner,
    },
    {
      path: "/qr-test",
      name: "qr-test",
      component: PatientScanner,
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

export default router;
