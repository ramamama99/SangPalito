import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import RegisterView from '../views/RegisterView.vue'
import TicketView from '../views/TicketView.vue'
import LoginView from '../views/LoginView.vue'
import StaffView from '../views/StaffView.vue'
import MedicalRecordView from '../views/MedicalRecordView.vue'
import ScanView from '../views/ScanView.vue'
import LiveView from '../views/LiveView.vue'
import TvView from '../views/TvView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView,
    },
    {
      path: '/register',
      name: 'register',
      component: RegisterView,
    },
    {
      path: '/ticket',
      name: 'ticket',
      component: TicketView,
    },
    {
      path: '/login',
      name: 'login',
      component: LoginView,
    },
    {
      path: '/staff',
      name: 'staff',
      component: StaffView,
    },
    {
      path: '/medical-record',
      name: 'medical-record',
      component: MedicalRecordView,
    },
    {
      path: '/scan',
      name: 'scan',
      component: ScanView,
    },
    {
      path: '/live',
      name: 'live',
      component: LiveView,
    },
    {
      path: '/tv',
      name: 'tv',
      component: TvView,
    },
  ],
  scrollBehavior(_to, _from, savedPosition) {
    if (savedPosition) {
      return savedPosition
    }
    return { top: 0 }
  },
})

export default router
