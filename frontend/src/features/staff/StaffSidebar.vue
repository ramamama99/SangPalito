<template>
  <aside
    class="staff-sidebar"
    :class="{ 'is-collapsed': collapsed }"
    :data-collapsed="collapsed"
  >
    <!-- 1. Header & Brand / Collapse Toggle -->
    <div class="sidebar-header">
      <div v-if="!collapsed" class="sidebar-brand">
        <div class="brand-badge">
          <Activity :size="18" class="text-[var(--deep)]" />
        </div>
        <div class="brand-info">
          <span class="brand-title">AntriSehat</span>
          <span class="brand-sub">Portal Faskes & EMR</span>
        </div>
      </div>
      <button
        type="button"
        class="toggle-btn"
        :aria-label="collapsed ? 'Perlebar Menu Sidebar' : 'Persempit Menu Sidebar'"
        :title="collapsed ? 'Perlebar Menu Sidebar' : 'Persempit Menu Sidebar'"
        @click="toggleCollapse"
      >
        <Menu :size="18" />
      </button>
    </div>

    <!-- 2. Nav Menu (Scrolls internally if viewport is small, never pushes logout down) -->
    <nav class="sidebar-nav">
      <!-- Group: Pelayanan Antrean -->
      <div v-if="!collapsed" class="nav-group-title">PELAYANAN KLINIK</div>

      <button
        type="button"
        class="nav-link"
        :class="{ active: queueStore.role === 'loket' }"
        :title="collapsed ? 'Dashboard Loket' : undefined"
        @click="selectRole('loket')"
      >
        <Users :size="19" class="nav-icon" />
        <span v-if="!collapsed" class="nav-text">Dashboard Loket</span>
        <span
          v-if="!collapsed && pendingVerificationCount > 0"
          class="badge-counter"
        >
          {{ pendingVerificationCount }}
        </span>
      </button>

      <button
        v-if="canAccessPoli"
        type="button"
        class="nav-link"
        :class="{ active: queueStore.role === 'medis' }"
        :title="collapsed ? 'Dashboard Poli Medis' : undefined"
        @click="selectRole('medis')"
      >
        <Stethoscope :size="19" class="nav-icon" />
        <span v-if="!collapsed" class="nav-text">Dashboard Poli</span>
        <span
          v-if="!collapsed && activePoliQueueCount > 0"
          class="badge-counter active-poli"
        >
          {{ activePoliQueueCount }}
        </span>
      </button>

      <!-- Group: Master Data Klinik -->
      <div v-if="!collapsed" class="nav-group-title">MASTER DATA KLINIK</div>

      <button
        type="button"
        class="nav-link"
        :class="{ active: queueStore.role === 'shift_dokter' }"
        :title="collapsed ? 'Shift & Jadwal Dokter' : undefined"
        @click="selectRole('shift_dokter')"
      >
        <CalendarClock :size="19" class="nav-icon" />
        <span v-if="!collapsed" class="nav-text">Shift Dokter (3 Shift)</span>
      </button>

      <button
        type="button"
        class="nav-link"
        :class="{ active: queueStore.role === 'database_pasien' }"
        :title="collapsed ? 'Database Master Pasien' : undefined"
        @click="selectRole('database_pasien')"
      >
        <Database :size="19" class="nav-icon" />
        <span v-if="!collapsed" class="nav-text">Database Pasien</span>
      </button>

      <button
        type="button"
        class="nav-link"
        :class="{ active: queueStore.role === 'database_obat' }"
        :title="collapsed ? 'Database Obat & e-FORNAS' : undefined"
        @click="selectRole('database_obat')"
      >
        <Pill :size="19" class="nav-icon" />
        <span v-if="!collapsed" class="nav-text">Database Obat</span>
      </button>

      <!-- Group: Administrator & Eksekutif -->
      <div v-if="!collapsed && (canAccessUserManagement || canAccessAuditLog || canAccessAnalytics)" class="nav-group-title">
        MANAJEMEN & AUDIT
      </div>

      <button
        v-if="canAccessUserManagement"
        type="button"
        class="nav-link"
        :class="{ active: queueStore.role === 'user_management' }"
        :title="collapsed ? 'Manajemen User' : undefined"
        @click="selectRole('user_management')"
      >
        <UserCog :size="19" class="nav-icon" />
        <span v-if="!collapsed" class="nav-text">Manajemen User</span>
      </button>

      <button
        v-if="canAccessAuditLog"
        type="button"
        class="nav-link"
        :class="{ active: queueStore.role === 'audit_log' }"
        :title="collapsed ? 'Log Audit Realtime' : undefined"
        @click="selectRole('audit_log')"
      >
        <ClipboardList :size="19" class="nav-icon" />
        <span v-if="!collapsed" class="nav-text">Log Audit</span>
        <span v-if="!collapsed" class="badge-counter badge-live">LIVE</span>
      </button>

      <button
        v-if="canAccessAnalytics"
        type="button"
        class="nav-link"
        :class="{ active: queueStore.role === 'analytics' }"
        :title="collapsed ? 'Dashboard Analitik' : undefined"
        @click="selectRole('analytics')"
      >
        <BarChart3 :size="19" class="nav-icon" />
        <span v-if="!collapsed" class="nav-text">Analitik Eksekutif</span>
      </button>

      <!-- Group: Display & Akses Publik -->
      <div v-if="!collapsed" class="nav-group-title">AKSES & DISPLAY</div>

      <RouterLink
        to="/tv"
        target="_blank"
        class="nav-link"
        :title="collapsed ? 'Layar TV Antrean Publik' : undefined"
      >
        <Monitor :size="19" class="nav-icon" />
        <span v-if="!collapsed" class="nav-text">Layar TV Publik</span>
      </RouterLink>

      <RouterLink
        to="/qr-test"
        class="nav-link"
        :title="collapsed ? 'Simulator QR & Antrean' : undefined"
      >
        <QrCode :size="19" class="nav-icon" />
        <span v-if="!collapsed" class="nav-text">Simulator QR</span>
      </RouterLink>
    </nav>

    <!-- 3. Pinned Bottom Area: Profil Petugas & Tombol Logout Langsung Terlihat -->
    <div class="sidebar-footer">
      <!-- Info Petugas / Dokter Aktif -->
      <div v-if="!collapsed" class="staff-profile-snippet">
        <div class="avatar-circle">
          {{ userInitials }}
        </div>
        <div class="profile-meta">
          <strong class="doctor-title">{{ staffDisplayName }}</strong>
          <span class="role-badge-label">{{ currentRoleBadge }}</span>
        </div>
      </div>

      <!-- Tombol Keluar (Logout) Pinned di Bawah, tidak perlu scroll -->
      <a
        href="#"
        class="logout-action-btn"
        :title="collapsed ? 'Keluar dari Sistem' : undefined"
        @click.prevent="$emit('logout')"
      >
        <LogOut :size="18" />
        <span v-if="!collapsed">Keluar</span>
      </a>
    </div>
  </aside>
</template>

<script setup lang="ts">
import { computed } from "vue";
import {
  Menu,
  Users,
  Stethoscope,
  Database,
  Pill,
  Monitor,
  QrCode,
  LogOut,
  Activity,
  CalendarClock,
  UserCog,
  ClipboardList,
  BarChart3,
} from "lucide-vue-next";
import { useQueueStore } from "../../stores/queue";
import { todayDate } from "../../constants/initialData";
import { isSameDate } from "../../utils/dateHelper";
import { api } from "../../services/api";
import { canAccessPoliDashboard } from "../../utils/auth";

const props = withDefaults(
  defineProps<{
    collapsed?: boolean;
    doctorName?: string;
  }>(),
  {
    collapsed: false,
    doctorName: "dr. Maya Putri",
  },
);

const emit = defineEmits<{
  (e: "update:collapsed", val: boolean): void;
  (e: "toggle"): void;
  (e: "logout"): void;
}>();

const queueStore = useQueueStore();

function toggleCollapse() {
  const nextVal = !props.collapsed;
  emit("update:collapsed", nextVal);
  emit("toggle");
}

const userRole = computed(() => {
  const currentUser = api.getCurrentUser();
  return currentUser?.role || 'admin';
});

const canAccessPoli = computed(() => canAccessPoliDashboard());

const canAccessUserManagement = computed(() => {
  return userRole.value === 'admin';
});

const canAccessAuditLog = computed(() => {
  return userRole.value === 'admin' || userRole.value === 'head';
});

const canAccessAnalytics = computed(() => {
  return userRole.value === 'admin' || userRole.value === 'head' || userRole.value === 'medic';
});

function selectRole(role: "loket" | "medis" | "database_pasien" | "database_obat" | "shift_dokter" | "user_management" | "audit_log" | "analytics") {
  if (role === "medis" && !canAccessPoli.value) {
    return;
  }
  queueStore.role = role;
  if (typeof window !== "undefined" && window.innerWidth < 768) {
    emit("update:collapsed", true);
  }
}

const staffDisplayName = computed(() => {
  const currentUser = api.getCurrentUser();
  if (currentUser?.full_name) {
    return currentUser.full_name;
  }
  switch (queueStore.role) {
    case "loket":
      return "Siti Rahmawati, A.Md.Kes";
    case "medis":
      return props.doctorName || "dr. Maya Putri, Sp.PD";
    case "database_pasien":
      return "Budi Pratama, A.Md.RMIK";
    case "database_obat":
      return "apt. Dimas Prasetyo, S.Farm";
    case "shift_dokter":
      return "dr. Hendra (Koordinator Medis)";
    case "user_management":
      return "Administrator Sistem";
    case "audit_log":
      return "Auditor Keamanan";
    case "analytics":
      return "dr. Hendra (Kepala Puskesmas)";
    default:
      return "Petugas Klinik";
  }
});

const currentRoleBadge = computed(() => {
  switch (queueStore.role) {
    case "loket":
      return "Petugas Loket";
    case "medis":
      return `Dokter · ${queueStore.staffPoli}`;
    case "database_pasien":
      return "Admin Rekam Medis";
    case "database_obat":
      return "Farmasi & Apoteker";
    case "shift_dokter":
      return "Admin Shift & SDM";
    case "user_management":
      return "Admin User";
    case "audit_log":
      return "Audit Trail";
    case "analytics":
      return "Kepala Faskes";
    default:
      return "Staff Klinik";
  }
});

const userInitials = computed(() => {
  const name = staffDisplayName.value;
  const parts = name.replace(/^(dr\.|apt\.)\s*/i, "").trim().split(" ");
  if (parts.length >= 2) {
    return (parts[0][0] + parts[1][0]).toUpperCase();
  }
  return name.slice(0, 2).toUpperCase();
});

function isMatchingToday(dateStr?: string): boolean {
  if (!dateStr) return false;
  return isSameDate(dateStr, todayDate);
}

const pendingVerificationCount = computed(() => {
  return queueStore.patients.filter(
    (p) =>
      isMatchingToday(p.date) &&
      ["Terdaftar", "Menunggu verifikasi"].includes(p.status),
  ).length;
});

const activePoliQueueCount = computed(() => {
  return queueStore.patients.filter(
    (p) =>
      p.poli === queueStore.staffPoli &&
      isMatchingToday(p.date) &&
      ["Hadir", "Dipanggil"].includes(p.status),
  ).length;
});
</script>

<style scoped>
.staff-sidebar {
  width: 250px;
  height: 100vh;
  position: sticky;
  top: 0;
  left: 0;
  background: #ffffff;
  border-right: 1px solid var(--line, #e2e8f0);
  display: flex;
  flex-direction: column;
  flex-shrink: 0;
  box-shadow: 2px 0 12px rgba(15, 23, 42, 0.03);
  z-index: 40;
  transition: width 0.22s cubic-bezier(0.4, 0, 0.2, 1);
  overflow: hidden;
}

.staff-sidebar.is-collapsed {
  width: 72px;
}

/* Header */
.sidebar-header {
  height: 72px;
  padding: 0 16px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  border-bottom: 1px solid var(--line, #f1f5f9);
  flex-shrink: 0;
  background: #ffffff;
}

.sidebar-brand {
  display: flex;
  align-items: center;
  gap: 10px;
  overflow: hidden;
}

.brand-badge {
  width: 36px;
  height: 36px;
  border-radius: 10px;
  background: #eaf6f2;
  border: 1px solid #c8e8df;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.brand-info {
  display: flex;
  flex-direction: column;
}

.brand-title {
  font-family: 'Manrope', sans-serif;
  font-size: 16px;
  font-weight: 700;
  color: var(--deep);
  letter-spacing: -0.2px;
  line-height: 1.1;
}

.brand-sub {
  font-size: 11px;
  color: var(--muted);
  font-weight: 500;
}

.toggle-btn {
  background: transparent;
  border: 1px solid transparent;
  border-radius: 50px;
  padding: 8px;
  cursor: pointer;
  color: var(--muted);
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.18s ease;
}

.toggle-btn:hover {
  background: #f6f8f5;
  color: var(--deep);
}

/* Navigation Menu */
.sidebar-nav {
  flex: 1;
  overflow-y: auto;
  overflow-x: hidden;
  padding: 16px 12px;
  display: flex;
  flex-direction: column;
  gap: 4px;
  scrollbar-width: thin;
  scrollbar-color: var(--line) transparent;
}

.sidebar-nav::-webkit-scrollbar {
  width: 4px;
}
.sidebar-nav::-webkit-scrollbar-track {
  background: transparent;
}
.sidebar-nav::-webkit-scrollbar-thumb {
  background: var(--line);
  border-radius: 4px;
}

.nav-group-title {
  font-size: 10px;
  font-weight: 700;
  color: var(--muted);
  letter-spacing: 0.8px;
  padding: 12px 10px 6px;
  text-transform: uppercase;
}

.nav-link {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 14px;
  border-radius: 10px;
  font-size: 13.5px;
  font-weight: 500;
  color: #243e40;
  background: transparent;
  border: none;
  width: 100%;
  cursor: pointer;
  text-decoration: none;
  transition: all 0.18s ease;
  position: relative;
  text-align: left;
}

.staff-sidebar.is-collapsed .nav-link {
  justify-content: center;
  padding: 12px 0;
}

.nav-link:hover {
  background: #f6f8f5;
  color: var(--deep);
}

.nav-link.active {
  background: #eaf6f2;
  color: var(--deep);
  font-weight: 700;
}

.nav-link.active .nav-icon {
  color: var(--teal);
}

.nav-icon {
  color: var(--muted);
  flex-shrink: 0;
  transition: color 0.18s ease;
}

.nav-text {
  flex: 1;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.badge-counter {
  font-size: 11px;
  font-weight: 700;
  padding: 2px 7px;
  border-radius: 999px;
  background: #f6f8f5;
  color: var(--deep);
  border: 1px solid var(--line);
}

.badge-counter.active-poli {
  background: var(--warm);
  color: var(--deep);
  border: 1px solid #e7ded4;
}

/* Pinned Bottom Footer */
.sidebar-footer {
  flex-shrink: 0;
  margin-top: auto;
  padding: 14px 12px;
  background: #ffffff;
  border-top: 1px solid var(--line);
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.staff-profile-snippet {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 8px 10px;
  border-radius: 12px;
  background: #f6f8f5;
  border: 1px solid var(--line);
}

.avatar-circle {
  width: 34px;
  height: 34px;
  border-radius: 10px;
  background: var(--deep);
  color: #ffffff;
  font-size: 12px;
  font-weight: 700;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  font-family: 'Manrope', sans-serif;
}

.profile-meta {
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.doctor-title {
  font-family: 'Manrope', sans-serif;
  font-size: 12.5px;
  font-weight: 700;
  color: var(--deep);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  line-height: 1.2;
}

.role-badge-label {
  font-size: 10.5px;
  color: var(--muted);
  font-weight: 500;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

/* Logout Button (Pinned, visually prominent, pill shape) */
.logout-action-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 9px 16px;
  border-radius: 50px;
  background: #fdf3f0;
  color: var(--coral);
  border: 1px solid #f7d6cf;
  font-size: 13px;
  font-weight: 600;
  text-decoration: none;
  cursor: pointer;
  transition: all 0.18s ease;
}

.staff-sidebar.is-collapsed .logout-action-btn {
  padding: 10px 0;
  border-radius: 12px;
  justify-content: center;
}

.logout-action-btn:hover {
  background: #fbeae5;
  color: #c96855;
  box-shadow: 0 2px 8px rgba(217, 127, 108, 0.15);
  transform: translateY(-1px);
}

@media (max-width: 768px) {
  .staff-sidebar {
    position: fixed;
    top: 0;
    left: 0;
    bottom: 0;
    z-index: 50;
    box-shadow: 4px 0 24px rgba(25, 90, 95, 0.12);
  }
  .staff-sidebar.is-collapsed {
    position: sticky;
    width: 60px;
    box-shadow: 1px 0 6px rgba(25, 90, 95, 0.04);
  }
}
</style>
