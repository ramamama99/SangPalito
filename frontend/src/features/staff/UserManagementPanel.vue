<template>
  <section class="dashboard-panel user-management-panel">
    <!-- Header Section -->
    <div class="panel-header-row">
      <div>
        <h3 class="panel-title">Manajemen Akun Pengguna & Staf</h3>
        <p class="panel-subtitle">
          Kelola hak akses login, data staf medis, petugas loket, kepala puskesmas, dan penempatan poli.
        </p>
      </div>
      <div class="header-actions">
        <button type="button" class="button outline btn-sm" :disabled="loading" @click="fetchUsers">
          <RotateCw :size="16" :class="{ 'spin-anim': loading }" /> Muat Ulang
        </button>
        <button type="button" class="button primary btn-sm" @click="openCreateModal">
          <UserPlus :size="16" /> Tambah Staf Baru
        </button>
      </div>
    </div>

    <!-- Alert / Banner Notifikasi -->
    <div v-if="successMessage" class="alert-banner success">
      <CheckCircle2 :size="18" />
      <span>{{ successMessage }}</span>
      <button type="button" class="close-alert" @click="successMessage = ''">×</button>
    </div>
    <div v-if="errorMessage" class="alert-banner error">
      <AlertTriangle :size="18" />
      <span>{{ errorMessage }}</span>
      <button type="button" class="close-alert" @click="errorMessage = ''">×</button>
    </div>

    <!-- KPI Summary Cards -->
    <div class="stats-grid user-kpis">
      <div class="stat-card">
        <div class="kpi-icon-wrap bg-blue-subtle">
          <Users :size="20" class="text-blue" />
        </div>
        <div class="kpi-content">
          <span class="kpi-label">Total Akun Terdaftar</span>
          <strong class="kpi-value">{{ users.length }}</strong>
          <span class="kpi-sub">{{ activeUserCount }} akun aktif</span>
        </div>
      </div>

      <div class="stat-card">
        <div class="kpi-icon-wrap bg-teal-subtle">
          <Stethoscope :size="20" class="text-teal" />
        </div>
        <div class="kpi-content">
          <span class="kpi-label">Dokter & Medis</span>
          <strong class="kpi-value">{{ countByRole('medic') }}</strong>
          <span class="kpi-sub">Pemeriksa & EMR</span>
        </div>
      </div>

      <div class="stat-card">
        <div class="kpi-icon-wrap bg-amber-subtle">
          <UserCheck :size="20" class="text-amber" />
        </div>
        <div class="kpi-content">
          <span class="kpi-label">Petugas Loket</span>
          <strong class="kpi-value">{{ countByRole('registrar') }}</strong>
          <span class="kpi-sub">Pendaftaran & Antrean</span>
        </div>
      </div>

      <div class="stat-card">
        <div class="kpi-icon-wrap bg-purple-subtle">
          <Award :size="20" class="text-purple" />
        </div>
        <div class="kpi-content">
          <span class="kpi-label">Kepala Puskesmas & Admin</span>
          <strong class="kpi-value">{{ countByRole('head') + countByRole('admin') }}</strong>
          <span class="kpi-sub">{{ countByRole('head') }} Kepala · {{ countByRole('admin') }} Admin</span>
        </div>
      </div>
    </div>

    <!-- Filter & Search Bar -->
    <div class="filter-controls-row">
      <div class="search-input-wrap">
        <Search :size="16" class="search-icon" />
        <input
          v-model="searchQuery"
          type="text"
          placeholder="Cari berdasarkan nama, email, nomor HP, atau poli..."
          class="table-search-input"
        />
        <button v-if="searchQuery" type="button" class="clear-search-btn" @click="searchQuery = ''">×</button>
      </div>

      <div class="filter-selects">
        <select v-model="filterRole" class="filter-select">
          <option value="all">Semua Peran (Role)</option>
          <option value="admin">Administrator (admin)</option>
          <option value="registrar">Petugas Loket (registrar)</option>
          <option value="medic">Tenaga Medis / Dokter (medic)</option>
          <option value="head">Kepala Puskesmas (head)</option>
        </select>

        <select v-model="filterStatus" class="filter-select">
          <option value="all">Semua Status</option>
          <option value="active">Aktif</option>
          <option value="inactive">Nonaktif</option>
        </select>
      </div>
    </div>

    <!-- Tabel Daftar Pengguna -->
    <div class="table-responsive user-table-wrapper">
      <table class="data-table">
        <thead>
          <tr>
            <th>Pengguna / Staf</th>
            <th>Kontak & Email</th>
            <th>Peran (Role)</th>
            <th>Poli Bagian</th>
            <th>Status Akun</th>
            <th>Login Terakhir</th>
            <th class="text-center">Aksi</th>
          </tr>
        </thead>
        <tbody v-if="loading && users.length === 0">
          <tr>
            <td colspan="7" class="text-center py-8 text-muted">
              <RotateCw :size="24" class="spin-anim inline-block mb-2" />
              <p>Memuat daftar pengguna staf...</p>
            </td>
          </tr>
        </tbody>
        <tbody v-else-if="filteredUsers.length === 0">
          <tr>
            <td colspan="7" class="text-center py-8 text-muted">
              <UserX :size="32" class="inline-block mb-2 opacity-50" />
              <p>Tidak ada pengguna yang cocok dengan kriteria pencarian.</p>
            </td>
          </tr>
        </tbody>
        <tbody v-else>
          <tr v-for="u in filteredUsers" :key="u.id" class="user-row">
            <td>
              <div class="user-cell">
                <div class="user-avatar" :class="'avatar-' + u.role">
                  {{ getInitials(u.full_name) }}
                </div>
                <div class="user-info">
                  <strong class="user-name">{{ u.full_name }}</strong>
                  <span class="user-meta-id">ID: {{ u.id.slice(0, 8) }}...</span>
                </div>
              </div>
            </td>
            <td>
              <div class="contact-cell">
                <span class="user-email"><Mail :size="13" /> {{ u.email }}</span>
                <span class="user-phone">
                  <Phone :size="13" /> {{ u.phone || '-' }}
                </span>
              </div>
            </td>
            <td>
              <span class="role-badge" :class="'role-' + u.role">
                <Shield v-if="u.role === 'admin'" :size="12" />
                <Stethoscope v-else-if="u.role === 'medic'" :size="12" />
                <Award v-else-if="u.role === 'head'" :size="12" />
                <UserCheck v-else :size="12" />
                {{ formatRole(u.role) }}
              </span>
            </td>
            <td>
              <span v-if="u.polyclinic_name" class="poli-tag">
                <Building2 :size="12" /> {{ u.polyclinic_name }}
              </span>
              <span v-else class="text-muted text-xs">Semua Bagian / Loket</span>
            </td>
            <td>
              <span class="status-indicator" :class="u.is_active ? 'active' : 'inactive'">
                <span class="dot"></span>
                {{ u.is_active ? 'Aktif' : 'Nonaktif' }}
              </span>
            </td>
            <td>
              <span class="text-xs text-muted">
                {{ formatDateTime(u.last_login_at) }}
              </span>
            </td>
            <td class="text-center actions-col">
              <div class="action-buttons-group">
                <button
                  type="button"
                  class="action-btn edit"
                  title="Ubah Profil Pengguna"
                  @click="openEditModal(u)"
                >
                  <Pencil :size="15" />
                </button>
                <button
                  type="button"
                  class="action-btn key"
                  title="Reset Kata Sandi"
                  @click="openResetPasswordModal(u)"
                >
                  <KeyRound :size="15" />
                </button>
                <button
                  type="button"
                  class="action-btn delete"
                  title="Hapus Pengguna"
                  @click="confirmDeleteUser(u)"
                >
                  <Trash2 :size="15" />
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- MODAL: Tambah Pengguna Baru -->
    <div v-if="showCreateModal" class="modal-backdrop" @click.self="showCreateModal = false">
      <div class="modal-card modal-lg">
        <div class="modal-header">
          <div class="flex items-center gap-2">
            <div class="modal-icon-badge bg-blue-subtle text-blue">
              <UserPlus :size="18" />
            </div>
            <div>
              <h3 class="modal-title">Tambah Pengguna Staf Baru</h3>
              <p class="modal-subtitle">Buat akun staf untuk akses sistem AntriSehat dengan salt & hash enkripsi.</p>
            </div>
          </div>
          <button type="button" class="modal-close-btn" @click="showCreateModal = false">×</button>
        </div>

        <form novalidate @submit.prevent="submitCreateUser" class="modal-form">
          <div class="form-grid">
            <div class="form-group full-width">
              <label>Nama Lengkap Staf <span class="req">*</span></label>
              <input
                v-model="createForm.full_name"
                type="text"
                :class="['form-input', { 'input-error': createErrors.full_name }]"
                @input="delete createErrors.full_name"
                placeholder="Contoh: dr. Rama Palito, Sp.PD / Siti Rahmawati, A.Md"
              />
              <span v-if="createErrors.full_name" class="field-error-text">
                {{ createErrors.full_name }}
              </span>
            </div>

            <div class="form-group">
              <label>Alamat Email (Login ID) <span class="req">*</span></label>
              <input
                v-model="createForm.email"
                type="email"
                :class="['form-input', { 'input-error': createErrors.email }]"
                @input="delete createErrors.email"
                placeholder="staf@antrisehat.id"
              />
              <span v-if="createErrors.email" class="field-error-text">
                {{ createErrors.email }}
              </span>
            </div>

            <div class="form-group">
              <label>Nomor Telepon / WhatsApp <span class="req">*</span></label>
              <input
                v-model="createForm.phone"
                type="tel"
                :class="['form-input', { 'input-error': createErrors.phone }]"
                @input="delete createErrors.phone"
                placeholder="081234567890"
              />
              <span v-if="createErrors.phone" class="field-error-text">
                {{ createErrors.phone }}
              </span>
            </div>

            <div class="form-group">
              <label>Kata Sandi Baru <span class="req">*</span></label>
              <div class="password-input-wrap">
                <input
                  v-model="createForm.password"
                  :type="showCreatePassword ? 'text' : 'password'"
                  :class="['form-input', { 'input-error': createErrors.password }]"
                  @input="delete createErrors.password"
                  placeholder="Minimal 6 karakter"
                />
                <button
                  type="button"
                  class="toggle-pw-btn"
                  @click="showCreatePassword = !showCreatePassword"
                >
                  <Eye v-if="!showCreatePassword" :size="16" />
                  <EyeOff v-else :size="16" />
                </button>
              </div>
              <span v-if="createErrors.password" class="field-error-text">
                {{ createErrors.password }}
              </span>
              <span v-else class="input-hint">Password akan di-hash dan di-salt secara otomatis dengan bcrypt di backend.</span>
            </div>

            <div class="form-group">
              <label>Peran / Role Pengguna <span class="req">*</span></label>
              <select
                v-model="createForm.role"
                :class="['form-input', { 'input-error': createErrors.role }]"
                @change="delete createErrors.role"
              >
                <option value="registrar">Petugas Loket / Pendaftaran (registrar)</option>
                <option value="medic">Tenaga Medis / Dokter (medic)</option>
                <option value="head">Kepala Puskesmas / Faskes (head)</option>
                <option value="admin">Administrator Sistem (admin)</option>
              </select>
              <span v-if="createErrors.role" class="field-error-text">
                {{ createErrors.role }}
              </span>
            </div>

            <div class="form-group full-width">
              <label>Poli Bagian Apa (Penempatan Staf)</label>
              <select v-model="createForm.polyclinic_id" class="form-input">
                <option value="">-- Bebas / Semua Poliklinik (Default) --</option>
                <option v-for="poli in polyclinics" :key="poli.id" :value="poli.id">
                  {{ poli.name }} (Prefix: {{ poli.queue_prefix }})
                </option>
              </select>
              <span class="input-hint">Pilih poli jika pengguna bertugas khusus di poli tertentu (khusus dokter/perawat).</span>
            </div>
          </div>

          <div class="modal-footer">
            <button type="button" class="button outline" @click="showCreateModal = false">Batal</button>
            <button type="submit" class="button primary" :disabled="formSubmitting">
              <span v-if="formSubmitting">Menyimpan...</span>
              <span v-else>Simpan & Buat Akun</span>
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- MODAL: Edit Pengguna -->
    <div v-if="showEditModal && editingUser" class="modal-backdrop" @click.self="showEditModal = false">
      <div class="modal-card">
        <div class="modal-header">
          <div class="flex items-center gap-2">
            <div class="modal-icon-badge bg-amber-subtle text-amber">
              <Pencil :size="18" />
            </div>
            <div>
              <h3 class="modal-title">Perbarui Data Pengguna</h3>
              <p class="modal-subtitle">{{ editingUser.email }}</p>
            </div>
          </div>
          <button type="button" class="modal-close-btn" @click="showEditModal = false">×</button>
        </div>

        <form novalidate @submit.prevent="submitUpdateUser" class="modal-form">
          <div class="form-grid">
            <div class="form-group full-width">
              <label>Nama Lengkap Staf <span class="req">*</span></label>
              <input
                v-model="editForm.full_name"
                type="text"
                :class="['form-input', { 'input-error': editErrors.full_name }]"
                @input="delete editErrors.full_name"
              />
              <span v-if="editErrors.full_name" class="field-error-text">
                {{ editErrors.full_name }}
              </span>
            </div>

            <div class="form-group">
              <label>Nomor Telepon / WhatsApp</label>
              <input
                v-model="editForm.phone"
                type="tel"
                :class="['form-input', { 'input-error': editErrors.phone }]"
                @input="delete editErrors.phone"
                placeholder="081234567890"
              />
              <span v-if="editErrors.phone" class="field-error-text">
                {{ editErrors.phone }}
              </span>
            </div>

            <div class="form-group">
              <label>Peran / Role <span class="req">*</span></label>
              <select
                v-model="editForm.role"
                :class="['form-input', { 'input-error': editErrors.role }]"
                @change="delete editErrors.role"
              >
                <option value="registrar">Petugas Loket (registrar)</option>
                <option value="medic">Tenaga Medis / Dokter (medic)</option>
                <option value="head">Kepala Puskesmas (head)</option>
                <option value="admin">Administrator Sistem (admin)</option>
              </select>
              <span v-if="editErrors.role" class="field-error-text">
                {{ editErrors.role }}
              </span>
            </div>

            <div class="form-group full-width">
              <label>Poli Bagian</label>
              <select v-model="editForm.polyclinic_id" class="form-input">
                <option value="">-- Bebas / Semua Poliklinik --</option>
                <option v-for="poli in polyclinics" :key="poli.id" :value="poli.id">
                  {{ poli.name }} (Prefix: {{ poli.queue_prefix }})
                </option>
              </select>
            </div>

            <div class="form-group full-width">
              <label class="toggle-label-wrap">
                <input v-model="editForm.is_active" type="checkbox" class="toggle-checkbox" />
                <span>Akun Aktif (Bisa Login ke Sistem)</span>
              </label>
            </div>
          </div>

          <div class="modal-footer">
            <button type="button" class="button outline" @click="showEditModal = false">Batal</button>
            <button type="submit" class="button primary" :disabled="formSubmitting">
              <span v-if="formSubmitting">Memperbarui...</span>
              <span v-else>Simpan Perubahan</span>
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- MODAL: Reset Password -->
    <div v-if="showResetModal && targetUser" class="modal-backdrop" @click.self="showResetModal = false">
      <div class="modal-card">
        <div class="modal-header">
          <div class="flex items-center gap-2">
            <div class="modal-icon-badge bg-purple-subtle text-purple">
              <KeyRound :size="18" />
            </div>
            <div>
              <h3 class="modal-title">Reset Kata Sandi</h3>
              <p class="modal-subtitle">Ganti kata sandi untuk {{ targetUser.full_name }}</p>
            </div>
          </div>
          <button type="button" class="modal-close-btn" @click="showResetModal = false">×</button>
        </div>

        <form novalidate @submit.prevent="submitResetPassword" class="modal-form">
          <div class="form-group full-width">
            <label>Kata Sandi Baru <span class="req">*</span></label>
            <div class="password-input-wrap">
              <input
                v-model="newPassword"
                :type="showResetPw ? 'text' : 'password'"
                :class="['form-input', { 'input-error': resetPwError }]"
                @input="resetPwError = ''"
                placeholder="Masukkan kata sandi baru (min 6 karakter)"
              />
              <button
                type="button"
                class="toggle-pw-btn"
                @click="showResetPw = !showResetPw"
              >
                <Eye v-if="!showResetPw" :size="16" />
                <EyeOff v-else :size="16" />
              </button>
            </div>
            <span v-if="resetPwError" class="field-error-text">
              {{ resetPwError }}
            </span>
            <span v-else class="input-hint">Kata sandi akan di-hash ulang dengan salt baru sebelum disimpan ke database.</span>
          </div>

          <div class="modal-footer">
            <button type="button" class="button outline" @click="showResetModal = false">Batal</button>
            <button type="submit" class="button primary" :disabled="formSubmitting">
              <span v-if="formSubmitting">Menyimpan...</span>
              <span v-else>Update Kata Sandi</span>
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- MODAL: Konfirmasi Hapus Pengguna -->
    <div v-if="showDeleteModal && targetUser" class="modal-backdrop" @click.self="showDeleteModal = false">
      <div class="modal-card modal-confirm">
        <div class="confirm-icon-wrap text-danger">
          <AlertCircle :size="32" />
        </div>
        <h3 class="modal-title text-center">Hapus Akun Pengguna?</h3>
        <p class="text-center text-muted text-sm mb-4">
          Apakah Anda yakin ingin menghapus akun <strong>{{ targetUser.full_name }}</strong> ({{ targetUser.email }})?
          Pengguna ini tidak akan dapat login lagi ke sistem.
        </p>

        <div class="modal-footer justify-center">
          <button type="button" class="button outline" @click="showDeleteModal = false">Batal</button>
          <button type="button" class="button danger" :disabled="formSubmitting" @click="submitDeleteUser">
            <span v-if="formSubmitting">Menghapus...</span>
            <span v-else>Ya, Hapus Pengguna</span>
          </button>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import {
  Users,
  UserPlus,
  UserCheck,
  UserX,
  Stethoscope,
  Shield,
  Award,
  Search,
  RotateCw,
  Pencil,
  Trash2,
  KeyRound,
  Eye,
  EyeOff,
  Mail,
  Phone,
  Building2,
  AlertTriangle,
  CheckCircle2,
  AlertCircle,
} from 'lucide-vue-next';
import { api } from '../../services/api';

interface UserItem {
  id: string;
  facility_id: string;
  email: string;
  full_name: string;
  phone?: string;
  polyclinic_id?: string;
  polyclinic_name?: string;
  role: 'admin' | 'registrar' | 'medic' | 'head';
  is_active: boolean;
  last_login_at?: string;
  created_at: string;
}

interface PolyclinicItem {
  id: string;
  name: string;
  queue_prefix: string;
}

const users = ref<UserItem[]>([]);
const polyclinics = ref<PolyclinicItem[]>([]);
const loading = ref(false);
const formSubmitting = ref(false);

const searchQuery = ref('');
const filterRole = ref('all');
const filterStatus = ref('all');

const successMessage = ref('');
const errorMessage = ref('');

// Modals
const showCreateModal = ref(false);
const showEditModal = ref(false);
const showResetModal = ref(false);
const showDeleteModal = ref(false);

const showCreatePassword = ref(false);
const showResetPw = ref(false);

const editingUser = ref<UserItem | null>(null);
const targetUser = ref<UserItem | null>(null);
const newPassword = ref('');

const createErrors = ref<Record<string, string>>({});
const editErrors = ref<Record<string, string>>({});
const resetPwError = ref('');

const createForm = ref({
  full_name: '',
  email: '',
  phone: '',
  password: '',
  role: 'medic',
  polyclinic_id: '',
});

const editForm = ref({
  full_name: '',
  phone: '',
  role: 'medic',
  polyclinic_id: '',
  is_active: true,
});

async function fetchUsers() {
  loading.value = true;
  try {
    const res = await api.getUsers(1, 100);
    if (res.success && Array.isArray(res.data)) {
      users.value = res.data;
    } else {
      // Fallback demo accounts if server offline
      users.value = [
        {
          id: 'b1c2d3e4-f5a6-4b7c-8d9e-0f1a2b3c4d5e',
          facility_id: 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11',
          full_name: 'Administrator Sistem',
          email: 'admin@antrisehat.id',
          phone: '081234567800',
          role: 'admin',
          is_active: true,
          created_at: new Date().toISOString(),
        },
        {
          id: 'b2c3d4e5-f6a7-4b8c-9d0e-1f2a3b4c5d6e',
          facility_id: 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11',
          full_name: 'dr. Rama Palito',
          email: 'rama@antrisehat.id',
          phone: '081298765432',
          role: 'medic',
          polyclinic_name: 'Poli Umum',
          is_active: true,
          created_at: new Date().toISOString(),
        },
        {
          id: 'b3c4d5e6-f7a8-4b9c-0d1e-2f3a4b5c6d7e',
          facility_id: 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11',
          full_name: 'Petugas Loket 1',
          email: 'loket@antrisehat.id',
          phone: '081345678901',
          role: 'registrar',
          is_active: true,
          created_at: new Date().toISOString(),
        },
        {
          id: 'b4c5d6e7-f8a9-4b0c-1d2e-3f4a5b6c7d8e',
          facility_id: 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11',
          full_name: 'dr. Hendra (Kepala Puskesmas)',
          email: 'kepala@antrisehat.id',
          phone: '081122334455',
          role: 'head',
          is_active: true,
          created_at: new Date().toISOString(),
        },
      ];
    }
  } catch (err) {
    console.error('fetchUsers error:', err);
  } finally {
    loading.value = false;
  }
}

async function fetchPolyclinics() {
  try {
    const res = await api.getPolyclinics();
    if (res.success && Array.isArray(res.data)) {
      polyclinics.value = res.data;
    } else {
      polyclinics.value = [
        { id: '11111111-1111-1111-1111-111111111111', name: 'Poli Umum', queue_prefix: 'A' },
        { id: '22222222-2222-2222-2222-222222222222', name: 'Poli Gigi', queue_prefix: 'B' },
        { id: '33333333-3333-3333-3333-333333333333', name: 'Poli KIA / Anak', queue_prefix: 'C' },
        { id: '44444444-4444-4444-4444-444444444444', name: 'Poli Lansia / Geriatri', queue_prefix: 'D' },
      ];
    }
  } catch (err) {
    console.error('fetchPolyclinics error:', err);
  }
}

const activeUserCount = computed(() => users.value.filter((u) => u.is_active).length);

function countByRole(role: string) {
  return users.value.filter((u) => u.role === role).length;
}

const filteredUsers = computed(() => {
  return users.value.filter((u) => {
    // Role filter
    if (filterRole.value !== 'all' && u.role !== filterRole.value) return false;
    // Status filter
    if (filterStatus.value === 'active' && !u.is_active) return false;
    if (filterStatus.value === 'inactive' && u.is_active) return false;
    // Search query
    if (searchQuery.value) {
      const q = searchQuery.value.toLowerCase();
      const matchName = u.full_name.toLowerCase().includes(q);
      const matchEmail = u.email.toLowerCase().includes(q);
      const matchPhone = (u.phone || '').toLowerCase().includes(q);
      const matchPoli = (u.polyclinic_name || '').toLowerCase().includes(q);
      return matchName || matchEmail || matchPhone || matchPoli;
    }
    return true;
  });
});

function validateCreateUser(): boolean {
  createErrors.value = {};
  let valid = true;

  if (!createForm.value.full_name || !createForm.value.full_name.trim()) {
    createErrors.value.full_name = 'Nama lengkap staf wajib diisi';
    valid = false;
  }

  if (!createForm.value.email || !createForm.value.email.trim()) {
    createErrors.value.email = 'Alamat email wajib diisi';
    valid = false;
  } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(createForm.value.email.trim())) {
    createErrors.value.email = 'Format email tidak valid (contoh: staf@antrisehat.id)';
    valid = false;
  }

  if (!createForm.value.phone || !createForm.value.phone.trim()) {
    createErrors.value.phone = 'Nomor telepon / WhatsApp wajib diisi';
    valid = false;
  } else if (createForm.value.phone.trim().replace(/\D/g, '').length < 9) {
    createErrors.value.phone = 'Nomor telepon minimal 9 digit angka';
    valid = false;
  }

  if (!createForm.value.password) {
    createErrors.value.password = 'Kata sandi wajib diisi';
    valid = false;
  } else if (createForm.value.password.length < 6) {
    createErrors.value.password = 'Kata sandi minimal 6 karakter';
    valid = false;
  }

  if (!createForm.value.role) {
    createErrors.value.role = 'Peran pengguna wajib dipilih';
    valid = false;
  }

  return valid;
}

function validateEditUser(): boolean {
  editErrors.value = {};
  let valid = true;

  if (!editForm.value.full_name || !editForm.value.full_name.trim()) {
    editErrors.value.full_name = 'Nama lengkap staf wajib diisi';
    valid = false;
  }

  if (editForm.value.phone && editForm.value.phone.trim().replace(/\D/g, '').length < 9) {
    editErrors.value.phone = 'Nomor telepon minimal 9 digit angka';
    valid = false;
  }

  if (!editForm.value.role) {
    editErrors.value.role = 'Peran pengguna wajib dipilih';
    valid = false;
  }

  return valid;
}

function validateResetPassword(): boolean {
  resetPwError.value = '';
  if (!newPassword.value || !newPassword.value.trim()) {
    resetPwError.value = 'Kata sandi baru wajib diisi';
    return false;
  }
  if (newPassword.value.length < 6) {
    resetPwError.value = 'Kata sandi baru minimal 6 karakter';
    return false;
  }
  return true;
}

function openCreateModal() {
  createForm.value = {
    full_name: '',
    email: '',
    phone: '',
    password: '',
    role: 'medic',
    polyclinic_id: '',
  };
  createErrors.value = {};
  showCreatePassword.value = false;
  showCreateModal.value = true;
}

async function submitCreateUser() {
  if (!validateCreateUser()) {
    return;
  }

  formSubmitting.value = true;
  errorMessage.value = '';
  try {
    const payload: any = {
      full_name: createForm.value.full_name.trim(),
      email: createForm.value.email.trim(),
      phone: createForm.value.phone.trim() || undefined,
      password: createForm.value.password,
      role: createForm.value.role,
      polyclinic_id: createForm.value.polyclinic_id || undefined,
    };
    const res = await api.createUser(payload);
    if (res.success) {
      successMessage.value = `Pengguna "${payload.full_name}" berhasil dibuat dan kata sandi di-hash aman.`;
      showCreateModal.value = false;
      await fetchUsers();
    } else {
      errorMessage.value = res.error?.message || res.message || 'Gagal membuat pengguna baru';
    }
  } catch (err: any) {
    errorMessage.value = err.message || 'Gagal menyimpan ke server';
  } finally {
    formSubmitting.value = false;
  }
}

function openEditModal(u: UserItem) {
  editingUser.value = u;
  editForm.value = {
    full_name: u.full_name,
    phone: u.phone || '',
    role: u.role,
    polyclinic_id: u.polyclinic_id || '',
    is_active: u.is_active,
  };
  editErrors.value = {};
  showEditModal.value = true;
}

async function submitUpdateUser() {
  if (!editingUser.value) return;
  if (!validateEditUser()) {
    return;
  }

  formSubmitting.value = true;
  errorMessage.value = '';
  try {
    const payload: any = {
      full_name: editForm.value.full_name.trim(),
      phone: editForm.value.phone.trim() || undefined,
      role: editForm.value.role,
      polyclinic_id: editForm.value.polyclinic_id || undefined,
      is_active: editForm.value.is_active,
    };
    const res = await api.updateUser(editingUser.value.id, payload);
    if (res.success) {
      successMessage.value = `Data pengguna "${editForm.value.full_name}" berhasil diperbarui.`;
      showEditModal.value = false;
      await fetchUsers();
    } else {
      errorMessage.value = res.error?.message || res.message || 'Gagal memperbarui pengguna';
    }
  } catch (err: any) {
    errorMessage.value = err.message || 'Gagal menyimpan perubahan ke server';
  } finally {
    formSubmitting.value = false;
  }
}

function openResetPasswordModal(u: UserItem) {
  targetUser.value = u;
  newPassword.value = '';
  resetPwError.value = '';
  showResetPw.value = false;
  showResetModal.value = true;
}

async function submitResetPassword() {
  if (!targetUser.value) return;
  if (!validateResetPassword()) {
    return;
  }

  formSubmitting.value = true;
  errorMessage.value = '';
  try {
    const res = await api.resetUserPassword(targetUser.value.id, newPassword.value);
    if (res.success) {
      successMessage.value = `Kata sandi untuk "${targetUser.value.full_name}" berhasil direset.`;
      showResetModal.value = false;
    } else {
      errorMessage.value = res.error?.message || res.message || 'Gagal mereset kata sandi';
    }
  } catch (err: any) {
    errorMessage.value = err.message || 'Gagal mereset kata sandi';
  } finally {
    formSubmitting.value = false;
  }
}

function confirmDeleteUser(u: UserItem) {
  targetUser.value = u;
  showDeleteModal.value = true;
}

async function submitDeleteUser() {
  if (!targetUser.value) return;
  formSubmitting.value = true;
  errorMessage.value = '';
  try {
    const res = await api.deleteUser(targetUser.value.id);
    if (res.success) {
      successMessage.value = `Pengguna "${targetUser.value.full_name}" berhasil dihapus.`;
      showDeleteModal.value = false;
      await fetchUsers();
    } else {
      errorMessage.value = res.error?.message || res.message || 'Gagal menghapus pengguna';
    }
  } catch (err: any) {
    errorMessage.value = err.message || 'Gagal menghapus pengguna';
  } finally {
    formSubmitting.value = false;
  }
}

function formatRole(role: string): string {
  switch (role) {
    case 'admin':
      return 'Administrator';
    case 'medic':
      return 'Tenaga Medis';
    case 'registrar':
      return 'Petugas Loket';
    case 'head':
      return 'Kepala Puskesmas';
    default:
      return role;
  }
}

function getInitials(name: string): string {
  const clean = name.replace(/^(dr\.|drg\.|apt\.)\s*/i, '').trim();
  const parts = clean.split(' ');
  if (parts.length >= 2) {
    return (parts[0][0] + parts[1][0]).toUpperCase();
  }
  return clean.slice(0, 2).toUpperCase();
}

function formatDateTime(val?: string): string {
  if (!val) return 'Belum pernah';
  try {
    const d = new Date(val);
    return d.toLocaleString('id-ID', {
      day: 'numeric',
      month: 'short',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    });
  } catch {
    return val;
  }
}

onMounted(() => {
  fetchUsers();
  fetchPolyclinics();
});
</script>

<style scoped>
.user-management-panel {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.panel-header-row {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 16px;
  flex-wrap: wrap;
}

.badge-role-tag {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-size: 11px;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  padding: 4px 10px;
  border-radius: 6px;
  margin-bottom: 6px;
}

.admin-tag {
  background: rgba(220, 38, 38, 0.1);
  color: #b91c1c;
  border: 1px solid rgba(220, 38, 38, 0.2);
}

.panel-title {
  font-size: 20px;
  font-weight: 700;
  color: #0f172a;
  margin: 0 0 4px 0;
}

.panel-subtitle {
  font-size: 13px;
  color: #64748b;
  margin: 0;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 10px;
}

.btn-sm {
  padding: 8px 14px;
  font-size: 13px;
  display: inline-flex;
  align-items: center;
  gap: 6px;
}

.spin-anim {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}

.alert-banner {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 16px;
  border-radius: 8px;
  font-size: 13px;
  font-weight: 500;
}

.alert-banner.success {
  background: #f0fdf4;
  border: 1px solid #bbf7d0;
  color: #166534;
}

.alert-banner.error {
  background: #fef2f2;
  border: 1px solid #fecaca;
  color: #991b1b;
}

.close-alert {
  margin-left: auto;
  background: transparent;
  border: none;
  font-size: 18px;
  cursor: pointer;
  color: inherit;
}

.user-kpis {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(210px, 1fr));
  gap: 16px;
}

.kpi-icon-wrap {
  width: 44px;
  height: 44px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 8px;
}

.bg-blue-subtle {
  background: #eff6ff;
}
.text-blue {
  color: #2563eb;
}
.bg-teal-subtle {
  background: #f0fdfa;
}
.text-teal {
  color: #0d9488;
}
.bg-amber-subtle {
  background: #fffbeb;
}
.text-amber {
  color: #d97706;
}
.bg-purple-subtle {
  background: #faf5ff;
}
.text-purple {
  color: #9333ea;
}

.kpi-content {
  display: flex;
  flex-direction: column;
}

.kpi-label {
  font-size: 12px;
  color: #64748b;
  font-weight: 500;
}

.kpi-value {
  font-size: 24px;
  font-weight: 700;
  color: #0f172a;
  line-height: 1.2;
}

.kpi-sub {
  font-size: 12px;
  color: #94a3b8;
  margin-top: 2px;
}

.filter-controls-row {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
  align-items: center;
  justify-content: space-between;
}

.search-input-wrap {
  position: relative;
  flex: 1;
  min-width: 260px;
}

.search-icon {
  position: absolute;
  left: 12px;
  top: 50%;
  transform: translateY(-50%);
  color: #94a3b8;
}

.table-search-input {
  width: 100%;
  padding: 9px 36px 9px 36px;
  font-size: 13px;
  border: 1px solid #cbd5e1;
  border-radius: 8px;
  background: #fff;
  outline: none;
  transition: border-color 0.15s;
}

.table-search-input:focus {
  border-color: #0d9488;
}

.clear-search-btn {
  position: absolute;
  right: 10px;
  top: 50%;
  transform: translateY(-50%);
  background: transparent;
  border: none;
  font-size: 16px;
  color: #94a3b8;
  cursor: pointer;
}

.filter-selects {
  display: flex;
  gap: 10px;
}

.filter-select {
  padding: 8px 12px;
  font-size: 13px;
  border: 1px solid #cbd5e1;
  border-radius: 8px;
  background: #fff;
  color: #334155;
  outline: none;
}

.user-table-wrapper {
  background: #fff;
  border-radius: 12px;
  border: 1px solid #e2e8f0;
  overflow: hidden;
}

.data-table {
  width: 100%;
  border-collapse: collapse;
  text-align: left;
  font-size: 13px;
}

.data-table th {
  background: #f8fafc;
  padding: 12px 16px;
  font-weight: 600;
  color: #475569;
  border-bottom: 1px solid #e2e8f0;
}

.data-table td {
  padding: 14px 16px;
  border-bottom: 1px solid #f1f5f9;
  vertical-align: middle;
}

.user-row:hover {
  background: #f8fafc;
}

.user-cell {
  display: flex;
  align-items: center;
  gap: 12px;
}

.user-avatar {
  width: 38px;
  height: 38px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 13px;
  font-weight: 700;
  flex-shrink: 0;
}

.avatar-admin {
  background: #fee2e2;
  color: #dc2626;
}
.avatar-medic {
  background: #ccfbf1;
  color: #0f766e;
}
.avatar-registrar {
  background: #fef3c7;
  color: #b45309;
}
.avatar-head {
  background: #f3e8ff;
  color: #7e22ce;
}

.user-info {
  display: flex;
  flex-direction: column;
}

.user-name {
  font-size: 14px;
  color: #0f172a;
}

.user-meta-id {
  font-size: 11px;
  color: #94a3b8;
  font-family: monospace;
}

.contact-cell {
  display: flex;
  flex-direction: column;
  gap: 3px;
}

.user-email,
.user-phone {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  font-size: 12px;
  color: #475569;
}

.role-badge {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  padding: 4px 9px;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 600;
}

.role-admin {
  background: #fef2f2;
  color: #b91c1c;
  border: 1px solid #fecaca;
}
.role-medic {
  background: #f0fdfa;
  color: #0f766e;
  border: 1px solid #99f6e4;
}
.role-registrar {
  background: #fffbeb;
  color: #b45309;
  border: 1px solid #fde68a;
}
.role-head {
  background: #faf5ff;
  color: #7e22ce;
  border: 1px solid #e9d5ff;
}

.poli-tag {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 3px 8px;
  border-radius: 5px;
  background: #f1f5f9;
  color: #334155;
  font-size: 12px;
  font-weight: 500;
}

.status-indicator {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  font-weight: 500;
}

.status-indicator.active {
  color: #16a34a;
}
.status-indicator.inactive {
  color: #94a3b8;
}

.status-indicator .dot {
  width: 7px;
  height: 7px;
  border-radius: 50%;
  background: currentColor;
}

.action-buttons-group {
  display: inline-flex;
  align-items: center;
  gap: 6px;
}

.action-btn {
  width: 32px;
  height: 32px;
  border-radius: 6px;
  border: 1px solid #e2e8f0;
  background: #fff;
  color: #64748b;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.15s;
}

.action-btn:hover {
  background: #f8fafc;
}

.action-btn.edit:hover {
  border-color: #3b82f6;
  color: #2563eb;
  background: #eff6ff;
}

.action-btn.key:hover {
  border-color: #a855f7;
  color: #7e22ce;
  background: #faf5ff;
}

.action-btn.delete:hover {
  border-color: #ef4444;
  color: #dc2626;
  background: #fef2f2;
}

/* Modals */
.modal-backdrop {
  position: fixed;
  inset: 0;
  background: rgba(15, 23, 42, 0.6);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 16px;
}

.modal-card {
  background: #fff;
  border-radius: 16px;
  width: 100%;
  max-width: 500px;
  max-height: 90vh;
  overflow-y: auto;
  box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.15);
  display: flex;
  flex-direction: column;
}

.modal-card.modal-lg {
  max-width: 620px;
}

.modal-card.modal-confirm {
  max-width: 420px;
  padding: 24px;
}

.modal-header {
  padding: 18px 22px;
  border-bottom: 1px solid #f1f5f9;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.modal-icon-badge {
  width: 36px;
  height: 36px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.modal-title {
  font-size: 16px;
  font-weight: 700;
  color: #0f172a;
  margin: 0;
}

.modal-subtitle {
  font-size: 12px;
  color: #64748b;
  margin: 0;
}

.modal-close-btn {
  background: transparent;
  border: none;
  font-size: 24px;
  color: #94a3b8;
  cursor: pointer;
  line-height: 1;
}

.modal-form {
  padding: 20px 22px;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.form-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 14px;
}

.form-group.full-width {
  grid-column: span 2;
}

.form-group label {
  display: block;
  font-size: 12px;
  font-weight: 600;
  color: #334155;
  margin-bottom: 5px;
}

.req {
  color: #dc2626;
}

.form-input {
  width: 100%;
  padding: 9px 12px;
  font-size: 13px;
  border: 1px solid #cbd5e1;
  border-radius: 8px;
  background: #fff;
  color: #0f172a;
  outline: none;
  transition: border-color 0.15s, box-shadow 0.15s;
}

.form-input:focus {
  border-color: #0d9488;
}

.form-input.input-error {
  border-color: #ef4444 !important;
  background-color: #fef2f2 !important;
  box-shadow: 0 0 0 3px rgba(239, 68, 68, 0.12) !important;
}

.field-error-text {
  display: block;
  font-size: 11.5px;
  font-weight: 600;
  color: #dc2626;
  margin-top: 4px;
  text-align: left;
}

.password-input-wrap {
  position: relative;
  display: flex;
  align-items: center;
}

.password-input-wrap input {
  padding-right: 38px;
}

.toggle-pw-btn {
  position: absolute;
  right: 10px;
  background: transparent;
  border: none;
  color: #94a3b8;
  cursor: pointer;
  display: flex;
  align-items: center;
}

.input-hint {
  display: block;
  font-size: 11px;
  color: #64748b;
  margin-top: 4px;
}

.toggle-label-wrap {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  color: #1e293b;
  cursor: pointer;
}

.toggle-checkbox {
  width: 16px;
  height: 16px;
  accent-color: #0d9488;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  padding-top: 12px;
  border-top: 1px solid #f1f5f9;
}

.confirm-icon-wrap {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 56px;
  height: 56px;
  border-radius: 50%;
  background: #fef2f2;
  margin: 0 auto 16px auto;
}

.text-danger {
  color: #dc2626;
}

.button.danger {
  background: #dc2626;
  color: #fff;
  border: none;
}

.button.danger:hover {
  background: #b91c1c;
}
</style>
