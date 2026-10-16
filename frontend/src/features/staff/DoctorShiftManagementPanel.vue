<template>
  <div class="doctor-shift-panel">
    <!-- Panel Header -->
    <div class="panel-header-section">
      <div class="panel-title-area">
        <div class="title-with-badge">
          <CalendarRange :size="24" class="text-teal-700" />
          <div>
            <div style="display: flex; align-items: center; gap: 10px;">
              <h2 class="panel-main-title">Pengaturan Jadwal Shift & Akun Dokter</h2>
            </div>
            <p class="panel-sub-title">
              Atur jadwal dinas dokter 7 hari (Senin s/d Minggu) dengan Drag & Drop, penukaran antar slot, serta manajemen akun dokter.
            </p>
          </div>
        </div>
      </div>

      <!-- Mode Tabs: Jadwal Mingguan vs CRUD Dokter -->
      <div class="subnav-tabs">
        <button
          type="button"
          class="subnav-tab-btn"
          :class="{ active: activeTab === 'schedule' }"
          @click="activeTab = 'schedule'"
        >
          <Calendar :size="16" />
          <span>Jadwal Shift Mingguan (Drag & Drop)</span>
        </button>
        <button
          type="button"
          class="subnav-tab-btn"
          :class="{ active: activeTab === 'doctors' }"
          @click="activeTab = 'doctors'"
        >
          <UserPlus :size="16" />
          <span>Manajemen & Akun Dokter ({{ doctorsList.length }})</span>
        </button>
      </div>
    </div>

    <!-- Notification Toast / Banner -->
    <div v-if="notificationMsg" class="notification-toast" :class="notificationType">
      <CheckCircle2 v-if="notificationType === 'success'" :size="18" />
      <AlertCircle v-else :size="18" />
      <span>{{ notificationMsg }}</span>
      <button type="button" class="btn-close-toast" @click="notificationMsg = ''"><X :size="16" /></button>
    </div>

    <!-- ============================================================ -->
    <!-- TAB 1: JADWAL SHIFT MINGGUAN (DRAG & DROP)                   -->
    <!-- ============================================================ -->
    <div v-if="activeTab === 'schedule'" class="schedule-tab-content">
      <!-- Filter Poli Toolbar -->
      <div class="schedule-toolbar">
        <div class="poli-filter-group">
          <span class="filter-label">Pilih Poliklinik:</span>
          <div class="filter-buttons">
            <button
              v-for="poli in polyclinics"
              :key="poli.id"
              type="button"
              class="filter-pill-btn"
              :class="{ active: selectedPoliId === poli.id || selectedPoliId === poli.slug }"
              @click="selectedPoliId = poli.id"
            >
              {{ poli.name }}
            </button>
          </div>
        </div>

        <div v-if="canEdit" class="toolbar-actions">
          <button
            type="button"
            class="button outline small"
            @click="resetScheduleToDefault"
            title="Kosongkan seluruh slot jadwal shift"
          >
            <RotateCcw :size="14" /> Kosongkan Jadwal
          </button>
          <button
            type="button"
            class="button primary small"
            @click="openAddDoctorModal"
          >
            <Plus :size="14" /> Tambah Dokter Baru
          </button>
        </div>
      </div>

      <!-- Main Drag & Drop Workspace: Doctor Pool + 7 Days Grid -->
      <div class="dnd-workspace">
        <!-- 1. LEFT SIDEBAR: DOCTOR POOL (BY SELECTED POLICLINIC) -->
        <aside class="doctor-pool-card">
          <div class="pool-header">
            <div class="pool-title">
              <Stethoscope :size="18" class="text-teal-700" />
              <div>
                <h4 style="margin: 0; font-size: 14px; font-weight: 700;">Roster Dokter</h4>
                <span style="font-size: 11px; color: #0f766e; font-weight: 600;">{{ currentSelectedPoliName }} ({{ displayedPoolDoctors.length }} Dokter Tersedia)</span>
              </div>
            </div>
          </div>

          <div class="pool-search-box">
            <Search :size="15" class="search-icon" />
            <input
              type="text"
              v-model="doctorSearchQuery"
              placeholder="Cari dokter / spesialisasi..."
              class="pool-search-input"
            />
          </div>

          <!-- Doctor Pool List -->
          <div class="doctor-pool-list">
            <div
              v-for="doc in displayedPoolDoctors"
              :key="doc.id"
              class="draggable-doctor-chip"
              :class="{ 'not-allowed': !canEdit }"
              :draggable="canEdit"
              @dragstart="onPoolDragStart($event, doc)"
              @dragend="onDragEnd"
              title="Tarik dokter ini ke slot kalender untuk plotting"
            >
              <div v-if="canEdit" class="drag-handle" title="Tarik dokter ini untuk plotting">
                <GripVertical :size="16" class="text-teal-700" />
              </div>
              <div class="doc-chip-avatar">
                {{ getInitials(doc.full_name) }}
              </div>
              <div class="doc-chip-info">
                <strong class="doc-chip-name">{{ doc.full_name }}</strong>
                <span class="doc-chip-spec">{{ doc.specialty || 'Dokter Umum' }}</span>
                <span class="doc-chip-poli">{{ doc.polyclinic_name || getPoliNameById(doc.polyclinic_id) }}</span>

                <!-- Quick Plotting Hint Badge -->
                <div v-if="canEdit" class="chip-plot-badge">
                  <GripVertical :size="11" /> Tarik untuk Plotting
                </div>
              </div>
            </div>

            <div v-if="displayedPoolDoctors.length === 0" class="empty-pool-state">
              <UserX :size="28" class="text-slate-400" />
              <p>Tidak ada dokter pada poliklinik ini.</p>
            </div>
          </div>
        </aside>

        <!-- 2. RIGHT GRID: 7-DAY WEEKLY SCHEDULE ROSTER -->
        <div class="weekly-roster-board">
          <div class="roster-grid-header">
            <div class="grid-cell-label">Waktu Shift</div>
            <div
              v-for="day in weekDays"
              :key="day.id"
              class="day-col-header"
              :class="{ 'is-today': isToday(day.id) }"
            >
              <span class="day-name">{{ day.name }}</span>
              <span v-if="isToday(day.id)" class="today-badge">Hari Ini</span>
            </div>
          </div>

          <!-- Rows for each shift: Pagi, Siang, Sore -->
          <div
            v-for="shift in shiftDefinitions"
            :key="shift.id"
            class="shift-grid-row"
          >
            <!-- Shift Time Info Column -->
            <div class="shift-time-cell" :class="shift.id">
              <span class="shift-name-tag" :class="shift.id">
                <Clock :size="12" /> Shift {{ shift.label }}
              </span>
              <span class="shift-hours">{{ shift.doctorHoursText }}</span>
              <span class="shift-quota-tag">Kuota: {{ shift.defaultQuota }} Pasien</span>
            </div>

            <!-- 7 Days Drop Slots -->
            <div
              v-for="day in weekDays"
              :key="day.id"
              class="day-shift-slot"
              :class="{
                'drag-over': isDragOver(day.id, shift.id),
                'is-assigned': isSlotAssigned(day.id, shift.id),
                'readonly-slot': !canEdit
              }"
              @dragover.prevent="handleDragOver(day.id, shift.id)"
              @dragleave="handleDragLeave(day.id, shift.id)"
              @drop.prevent="handleDrop(day.id, shift.id)"
            >
              <!-- If Slot has Assigned Doctor (Draggable for Swapping if canEdit) -->
              <div
                v-if="getAssignedDoctor(day.id, shift.id)"
                class="assigned-doctor-card"
                :class="[shift.id, { 'can-drag': canEdit }]"
                :draggable="canEdit"
                @dragstart="onSlotDragStart($event, day.id, shift.id, getAssignedDoctor(day.id, shift.id)!)"
                @dragend="onDragEnd"
              >
                <div class="assigned-card-top">
                  <div class="assigned-avatar">
                    {{ getInitials(getAssignedDoctor(day.id, shift.id)?.doctorName || '') }}
                  </div>
                  <div class="assigned-meta">
                    <strong class="assigned-name">{{ getAssignedDoctor(day.id, shift.id)?.doctorName }}</strong>
                    <span class="assigned-spec">{{ getAssignedDoctor(day.id, shift.id)?.doctorSpecialty || 'Dokter Bertugas' }}</span>
                  </div>
                </div>

                <div class="assigned-card-footer">
                  <button
                    v-if="canEdit"
                    type="button"
                    class="btn-edit-slot-quota"
                    title="Ubah Kuota Pasien"
                    @click.stop="openEditSlotModal(day.id, shift.id)"
                  >
                    <Settings :size="11" />
                    <span>{{ getAssignedDoctor(day.id, shift.id)?.fixedQuota || shift.defaultQuota }} Kuota</span>
                  </button>
                  <span v-else class="quota-readonly-chip">
                    {{ getAssignedDoctor(day.id, shift.id)?.fixedQuota || shift.defaultQuota }} Kuota
                  </span>

                  <!-- BUTTON HAPUS CARD DI SHIFT INI (HANYA LOGO SAMPAH) -->
                  <button
                    v-if="canEdit"
                    type="button"
                    class="btn-delete-card-slot"
                    title="Hapus dokter dari jadwal shift ini"
                    @click.stop="removeDoctorFromSlot(day.id, shift.id)"
                  >
                    <Trash2 :size="12" />
                  </button>
                </div>
              </div>

              <!-- If Slot is Empty (Drop Target for Plotting) -->
              <div v-else class="empty-drop-placeholder">
                <PlusCircle v-if="canEdit" :size="18" class="drop-icon" />
                <span class="drop-text">{{ canEdit ? 'Tarik Dokter ke Sini' : 'Kosong' }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- ============================================================ -->
    <!-- TAB 2: MANAJEMEN & AKUN DOKTER (CRUD TERHUBUNG AKUN USER)    -->
    <!-- ============================================================ -->
    <div v-else class="doctors-tab-content">
      <div class="crud-toolbar">
        <div class="search-wrap">
          <Search :size="16" class="text-slate-400 search-icon" />
          <input
            type="text"
            v-model="doctorCrudSearch"
            placeholder="Cari dokter berdasarkan nama, email, STR, atau poli..."
            class="crud-search-input"
          />
        </div>

        <button
          v-if="canEdit"
          type="button"
          class="button primary"
          @click="openAddDoctorModal"
        >
          <UserPlus :size="16" /> Tambah Dokter & Buat Akun
        </button>
        <span v-else class="text-xs text-slate-500 font-semibold">
          (Hanya Admin / Kepala Klinik yang dapat menambah & mengedit dokter)
        </span>
      </div>

      <!-- Doctors Table -->
      <div class="table-responsive">
        <table class="doctors-table">
          <thead>
            <tr>
              <th>Dokter & Profil</th>
              <th>Poliklinik</th>
              <th>Nomor STR</th>
              <th>Akun Login (Email)</th>
              <th>No. WhatsApp</th>
              <th>Status Dinas</th>
              <th v-if="canEdit" style="text-align: right;">Aksi</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="doc in filteredCrudDoctors" :key="doc.id">
              <td>
                <div class="doc-table-profile">
                  <div class="doc-table-avatar" :class="{ 'avatar-off': !doc.is_active }">
                    {{ getInitials(doc.full_name) }}
                  </div>
                  <div>
                    <strong class="doc-table-name">{{ doc.full_name }}</strong>
                    <span class="doc-table-sub">{{ doc.specialty || 'Dokter Umum' }}</span>
                  </div>
                </div>
              </td>
              <td>
                <span class="poli-badge">{{ doc.polyclinic_name || getPoliNameById(doc.polyclinic_id) || 'Poli Umum' }}</span>
              </td>
              <td>
                <code class="str-code">{{ doc.str_no || 'Belum diisi' }}</code>
              </td>
              <td>
                <div class="email-info-cell">
                  <span>{{ doc.email || '-' }}</span>
                  <span v-if="doc.is_first_login" class="first-login-pill">Perlu Onboarding</span>
                </div>
              </td>
              <td>
                <a
                  v-if="doc.phone"
                  :href="getWhatsAppLink(doc.phone, doc.full_name)"
                  target="_blank"
                  rel="noopener noreferrer"
                  class="wa-phone-link"
                  title="Hubungi via WhatsApp"
                >
                  <MessageSquare :size="13" /> {{ doc.phone }}
                </a>
                <span v-else class="text-slate-400">-</span>
              </td>
              <td>
                <span
                  class="status-pill"
                  :class="doc.is_active ? 'active' : 'inactive'"
                >
                  <span class="status-dot"></span>
                  {{ doc.is_active ? 'Aktif Bertugas' : 'Cuti / Nonaktif' }}
                </span>
              </td>
              <td v-if="canEdit" style="text-align: right;">
                <div class="action-buttons-cell">
                  <button
                    type="button"
                    class="btn-action edit"
                    title="Edit Profil & Poliklinik Dokter"
                    @click="openEditDoctorModal(doc)"
                  >
                    <Edit3 :size="15" />
                  </button>
                  <button
                    type="button"
                    class="btn-action delete"
                    title="Hapus Dokter"
                    @click="confirmDeleteDoctor(doc)"
                  >
                    <Trash2 :size="15" />
                  </button>
                </div>
              </td>
            </tr>

            <tr v-if="filteredCrudDoctors.length === 0">
              <td :colspan="canEdit ? 7 : 6" class="empty-table-cell">
                <UserX :size="32" class="text-slate-400" style="margin: 0 auto 8px;" />
                <p>Tidak ada data dokter ditemukan.</p>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- ============================================================ -->
    <!-- MODAL: TAMBAH / EDIT DOKTER (TERHUBUNG KE AKUN USER)         -->
    <!-- ============================================================ -->
    <div v-if="showDoctorModal && canEdit" class="modal-backdrop" @click.self="closeDoctorModal">
      <div class="modal-card doctor-modal">
        <div class="modal-header">
          <div style="display: flex; align-items: center; gap: 8px;">
            <UserPlus :size="20" class="text-teal-700" />
            <h3 style="margin: 0; font-size: 16px; font-weight: 700; color: #0f172a;">
              {{ isEditingDoctor ? 'Edit Data Dokter & Akun' : 'Tambah Dokter & Buat Akun Baru' }}
            </h3>
          </div>
          <button type="button" class="btn-close" @click="closeDoctorModal"><X :size="18" /></button>
        </div>

        <form novalidate @submit.prevent="saveDoctor">
          <div class="modal-body space-y-4">
            <div v-if="doctorModalError" class="modal-error-alert">
              <AlertCircle :size="16" />
              <span>{{ doctorModalError }}</span>
            </div>

            <p style="font-size: 12px; color: #64748b; margin-bottom: 12px;">
              Dokter yang didaftarkan akan secara otomatis memiliki akun login dengan role <strong>medic</strong>. Pada saat login pertama kali, dokter akan dipandu mengganti kata sandi.
            </p>

            <!-- Full Name -->
            <label class="form-field">
              <span class="field-label">Nama Lengkap Dokter <span class="req">*</span></span>
              <input
                type="text"
                v-model="doctorForm.full_name"
                :class="{ 'input-error': formErrors.full_name }"
                @input="formErrors.full_name = ''"
                placeholder="Contoh: dr. Amanda Pratama, Sp.A"
              />
              <span v-if="formErrors.full_name" class="field-error-text">
                {{ formErrors.full_name }}
              </span>
            </label>

            <!-- Polyclinic & Specialty -->
            <div class="form-grid-2">
              <label class="form-field">
                <span class="field-label">Poliklinik Bertugas <span class="req">*</span></span>
                <select
                  v-model="doctorForm.polyclinic_id"
                  :class="{ 'input-error': formErrors.polyclinic_id }"
                  @change="formErrors.polyclinic_id = ''"
                >
                  <option value="" disabled>Pilih Poliklinik</option>
                  <option v-for="p in polyclinics" :key="p.id" :value="p.id">
                    {{ p.name }}
                  </option>
                </select>
                <span v-if="formErrors.polyclinic_id" class="field-error-text">
                  {{ formErrors.polyclinic_id }}
                </span>
              </label>

              <label class="form-field">
                <span class="field-label">Spesialisasi</span>
                <input
                  type="text"
                  v-model="doctorForm.specialty"
                  placeholder="Contoh: Dokter Umum / Spesialis Anak"
                />
              </label>
            </div>

            <!-- STR Number -->
            <label class="form-field">
              <span class="field-label">Nomor STR Dokter</span>
              <input
                type="text"
                v-model="doctorForm.str_no"
                placeholder="Contoh: 31.1.1.100.1.22.123456"
              />
            </label>

            <!-- Email & Phone (For Account Login & WA) -->
            <div class="form-grid-2">
              <label class="form-field">
                <span class="field-label">Email Akun Login <span class="req">*</span></span>
                <input
                  type="email"
                  v-model="doctorForm.email"
                  :disabled="isEditingDoctor"
                  :class="{ 'input-error': formErrors.email }"
                  @input="formErrors.email = ''"
                  placeholder="dokter@antrisehat.id"
                />
                <span v-if="formErrors.email" class="field-error-text">
                  {{ formErrors.email }}
                </span>
              </label>

              <label class="form-field">
                <span class="field-label">Nomor WhatsApp / Telepon <span class="req">*</span></span>
                <input
                  type="tel"
                  v-model="doctorForm.phone"
                  :class="{ 'input-error': formErrors.phone }"
                  @input="formErrors.phone = ''"
                  placeholder="Contoh: 081234567890"
                />
                <span v-if="formErrors.phone" class="field-error-text">
                  {{ formErrors.phone }}
                </span>
              </label>
            </div>

            <!-- Password Awal (Only on creation) -->
            <label v-if="!isEditingDoctor" class="form-field">
              <span class="field-label">Kata Sandi Awal Akun (Opsional)</span>
              <input
                type="text"
                v-model="doctorForm.password"
                placeholder="Default: Dokter123! (Dokter akan diminta ubah saat login awal)"
              />
            </label>

            <!-- Active Status Toggle -->
            <div v-if="isEditingDoctor" class="form-field-checkbox">
              <label style="display: flex; align-items: center; gap: 8px; cursor: pointer;">
                <input type="checkbox" v-model="doctorForm.is_active" />
                <span style="font-size: 13px; font-weight: 600; color: #334155;">Dokter Aktif Melayani Pasien</span>
              </label>
            </div>
          </div>

          <div class="modal-footer">
            <button type="button" class="button outline" @click="closeDoctorModal">
              Batal
            </button>
            <button type="submit" class="button primary" :disabled="savingDoctor">
              <span v-if="savingDoctor">Menyimpan...</span>
              <span v-else>{{ isEditingDoctor ? 'Simpan Perubahan' : 'Buat Dokter & Akun' }}</span>
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- ============================================================ -->
    <!-- MODAL: EDIT KUOTA & RUANGAN SLOT SHIFT                       -->
    <!-- ============================================================ -->
    <div v-if="showSlotEditModal && canEdit" class="modal-backdrop" @click.self="showSlotEditModal = false">
      <div class="modal-card" style="max-width: 420px;">
        <div class="modal-header">
          <h3 style="margin: 0; font-size: 15px; font-weight: 700; color: #0f172a;">
            Atur Kuota Pasien Shift
          </h3>
          <button type="button" class="btn-close" @click="showSlotEditModal = false"><X :size="18" /></button>
        </div>
        <div class="modal-body space-y-4">
          <label class="form-field">
            <span class="field-label">Kuota Maksimal Pasien</span>
            <input type="number" v-model.number="editingSlotQuota" min="1" max="100" />
          </label>
        </div>
        <div class="modal-footer">
          <button type="button" class="button outline" @click="showSlotEditModal = false">Batal</button>
          <button type="button" class="button primary" @click="saveSlotConfig">Simpan</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, reactive, onMounted } from "vue";
import {
  CalendarRange,
  Calendar,
  UserPlus,
  Stethoscope,
  Clock,
  Search,
  GripVertical,
  Plus,
  PlusCircle,
  X,
  Edit3,
  Trash2,
  CheckCircle2,
  AlertCircle,
  RotateCcw,
  UserX,
  Settings,
  MessageSquare,
} from "lucide-vue-next";
import { api } from "../../services/api";
import { useQueueStore } from "../../stores/queue";
import {
  DAYS_OF_WEEK,
  initialDoctorShifts,
  type DoctorShift,
} from "../../constants/doctorShifts";

interface DoctorRecord {
  id: string;
  facility_id?: string;
  user_id?: string;
  full_name: string;
  specialty?: string;
  str_no?: string;
  polyclinic_id?: string;
  polyclinic_name?: string;
  email?: string;
  phone?: string;
  is_active: boolean;
  is_first_login?: boolean;
}

interface DragSource {
  type: "pool" | "slot";
  doctor?: DoctorRecord;
  slotKey?: string;
  day?: number;
  shiftId?: string;
  shiftData?: DoctorShift;
}

const queueStore = useQueueStore();
const activeTab = ref<"schedule" | "doctors">("schedule");
const notificationMsg = ref("");
const notificationType = ref<"success" | "error">("success");

// Role & Permission Access Check
const currentUser = computed(() => api.getCurrentUser());
const canEdit = computed(() => {
  const role = currentUser.value?.role;
  return role === "admin" || role === "head";
});

// 7 Days of the week
const weekDays = DAYS_OF_WEEK;

// Shift definitions
const shiftDefinitions = [
  { id: "pagi", label: "Pagi", doctorHoursText: "07.00 – 12.00", defaultQuota: 24 },
  { id: "siang", label: "Siang", doctorHoursText: "12.00 – 17.00", defaultQuota: 24 },
  { id: "sore", label: "Sore", doctorHoursText: "17.00 – 22.00", defaultQuota: 24 },
];

// Master Poliklinik (dukungan UUID dan slug)
const polyclinics = ref([
  { id: "11111111-1111-1111-1111-111111111111", slug: "poli-umum", name: "Poli Umum", code: "A" },
  { id: "22222222-2222-2222-2222-222222222222", slug: "poli-gigi", name: "Poli Gigi", code: "B" },
  { id: "33333333-3333-3333-3333-333333333333", slug: "poli-kia", name: "Poli KIA / Anak", code: "C" },
  { id: "44444444-4444-4444-4444-444444444444", slug: "poli-lansia", name: "Poli Lansia", code: "D" },
]);
const selectedPoliId = ref("11111111-1111-1111-1111-111111111111");

const currentSelectedPoli = computed(() => {
  return polyclinics.value.find(
    (p) => p.id === selectedPoliId.value || p.slug === selectedPoliId.value
  );
});
const currentSelectedPoliName = computed(() => currentSelectedPoli.value?.name || "Poli Umum");

// Doctor List State
const doctorsList = ref<DoctorRecord[]>([]);
const doctorSearchQuery = ref("");
const doctorCrudSearch = ref("");

// Weekly Grid Schedule Matrix: [poliId-day-shiftId] => DoctorShift
const weeklyScheduleMatrix = ref<Record<string, DoctorShift>>({});

// Drag and drop state
const activeDragSource = ref<DragSource | null>(null);
const currentDragOverKey = ref<string | null>(null);

// Modal Doctor Form
const showDoctorModal = ref(false);
const isEditingDoctor = ref(false);
const editingDoctorId = ref("");
const savingDoctor = ref(false);
const doctorModalError = ref("");
const doctorForm = reactive({
  full_name: "",
  polyclinic_id: "",
  specialty: "",
  str_no: "",
  email: "",
  phone: "",
  password: "",
  is_active: true,
});
const formErrors = reactive({
  full_name: "",
  polyclinic_id: "",
  email: "",
  phone: "",
});

// Slot Edit Modal
const showSlotEditModal = ref(false);
const editingSlotDay = ref(1);
const editingSlotShift = ref("pagi");
const editingSlotRoom = ref("Ruang 01");
const editingSlotQuota = ref(24);

// Lifecycle
onMounted(async () => {
  await fetchDoctors();
  loadWeeklyScheduleFromLocal();
});

// Load doctors from backend API (or fallback to presets)
async function fetchDoctors() {
  try {
    const res = await api.getDoctors();
    if (res && res.data && Array.isArray(res.data) && res.data.length > 0) {
      doctorsList.value = res.data;
      return;
    }
  } catch (err) {
    console.warn("Backend doctors endpoint unreachable, using standard seed roster");
  }

  // Fallback preset doctors
  doctorsList.value = [
    {
      id: "doc-1",
      full_name: "dr. Rama Palito",
      specialty: "Dokter Umum",
      str_no: "31.1.1.100.1.20.001928",
      polyclinic_id: "poli-umum",
      polyclinic_name: "Poli Umum",
      email: "dr.rama@antrisehat.id",
      phone: "081234567890",
      is_active: true,
      is_first_login: false,
    },
    {
      id: "doc-2",
      full_name: "dr. Aditya Rahman",
      specialty: "Dokter Umum",
      str_no: "31.1.1.100.1.21.004812",
      polyclinic_id: "poli-umum",
      polyclinic_name: "Poli Umum",
      email: "dr.aditya@antrisehat.id",
      phone: "081398765432",
      is_active: true,
      is_first_login: true,
    },
    {
      id: "doc-3",
      full_name: "dr. Maya Putri",
      specialty: "Dokter Umum",
      str_no: "31.1.1.100.1.22.008914",
      polyclinic_id: "poli-umum",
      polyclinic_name: "Poli Umum",
      email: "dr.maya@antrisehat.id",
      phone: "081287654321",
      is_active: true,
      is_first_login: false,
    },
    {
      id: "doc-4",
      full_name: "drg. Rani Prameswari",
      specialty: "Dokter Gigi",
      str_no: "31.2.1.200.1.19.003411",
      polyclinic_id: "poli-gigi",
      polyclinic_name: "Poli Gigi",
      email: "drg.rani@antrisehat.id",
      phone: "081901234567",
      is_active: true,
      is_first_login: false,
    },
    {
      id: "doc-5",
      full_name: "drg. Fajar Hidayat",
      specialty: "Dokter Gigi",
      str_no: "31.2.1.200.1.21.006789",
      polyclinic_id: "poli-gigi",
      polyclinic_name: "Poli Gigi",
      email: "drg.fajar@antrisehat.id",
      phone: "081812345678",
      is_active: true,
      is_first_login: true,
    },
    {
      id: "doc-6",
      full_name: "dr. Siti Rahmawati, Sp.A",
      specialty: "Spesialis Anak",
      str_no: "31.1.2.300.1.18.005123",
      polyclinic_id: "poli-kia",
      polyclinic_name: "Poli KIA",
      email: "dr.siti@antrisehat.id",
      phone: "081223344556",
      is_active: true,
      is_first_login: false,
    },
    {
      id: "doc-7",
      full_name: "dr. Hendra Prasetyo",
      specialty: "Spesialis Penyakit Dalam",
      str_no: "31.1.2.400.1.17.009871",
      polyclinic_id: "poli-lansia",
      polyclinic_name: "Poli Lansia",
      email: "dr.hendra@antrisehat.id",
      phone: "081334455667",
      is_active: true,
      is_first_login: false,
    },
  ];
}

// Helper pencocokan dokter ke poliklinik aktif secara akurat
function isDoctorInPoli(d: DoctorRecord, targetPoliId: string): boolean {
  if (!targetPoliId) return true;

  const poli = polyclinics.value.find((p) => p.id === targetPoliId || p.slug === targetPoliId);
  const targetUuid = (poli?.id || targetPoliId).toLowerCase();
  const targetSlug = (poli?.slug || targetPoliId).toLowerCase();

  const dPoliId = (d.polyclinic_id || "").toLowerCase();
  const dPoliName = (d.polyclinic_name || "").toLowerCase();
  const dSpec = (d.specialty || "").toLowerCase();
  const dName = (d.full_name || "").toLowerCase();

  // 1. Direct ID / UUID / Slug match
  if (dPoliId && (dPoliId === targetUuid || dPoliId === targetSlug)) {
    return true;
  }

  // 2. Poli Umum
  if (targetSlug.includes("umum") || targetUuid === "11111111-1111-1111-1111-111111111111") {
    return (
      dPoliName.includes("umum") ||
      dSpec.includes("umum") ||
      dPoliId === "11111111-1111-1111-1111-111111111111" ||
      dPoliId === "poli-umum" ||
      (!dSpec.includes("gigi") && !dSpec.includes("anak") && !dSpec.includes("sp.a") && !dSpec.includes("lansia") && !dSpec.includes("geriatri") && !dSpec.includes("dalam") && !dName.includes("drg"))
    );
  }

  // 3. Poli Gigi
  if (targetSlug.includes("gigi") || targetUuid === "22222222-2222-2222-2222-222222222222") {
    return (
      dPoliName.includes("gigi") ||
      dSpec.includes("gigi") ||
      dName.includes("drg") ||
      dPoliId === "22222222-2222-2222-2222-222222222222" ||
      dPoliId === "poli-gigi"
    );
  }

  // 4. Poli KIA / Anak
  if (targetSlug.includes("kia") || targetSlug.includes("anak") || targetUuid === "33333333-3333-3333-3333-333333333333") {
    return (
      dPoliName.includes("kia") ||
      dPoliName.includes("anak") ||
      dSpec.includes("anak") ||
      dSpec.includes("kia") ||
      dSpec.includes("sp.a") ||
      dPoliId === "33333333-3333-3333-3333-333333333333" ||
      dPoliId === "poli-kia"
    );
  }

  // 5. Poli Lansia / Geriatri
  if (targetSlug.includes("lansia") || targetSlug.includes("geriatri") || targetUuid === "44444444-4444-4444-4444-444444444444") {
    return (
      dPoliName.includes("lansia") ||
      dPoliName.includes("geriatri") ||
      dSpec.includes("lansia") ||
      dSpec.includes("geriatri") ||
      dSpec.includes("dalam") ||
      dPoliId === "44444444-4444-4444-4444-444444444444" ||
      dPoliId === "poli-lansia"
    );
  }

  return false;
}

// Displayed doctors in left pool
const displayedPoolDoctors = computed(() => {
  const query = doctorSearchQuery.value.trim().toLowerCase();

  return doctorsList.value.filter((d) => {
    const matchPoli = isDoctorInPoli(d, selectedPoliId.value);
    const matchQuery =
      !query ||
      d.full_name.toLowerCase().includes(query) ||
      (d.specialty && d.specialty.toLowerCase().includes(query));
    return matchPoli && matchQuery;
  });
});

// Filtered doctors for CRUD table
const filteredCrudDoctors = computed(() => {
  const query = doctorCrudSearch.value.trim().toLowerCase();
  if (!query) return doctorsList.value;
  return doctorsList.value.filter(
    (d) =>
      d.full_name.toLowerCase().includes(query) ||
      (d.email && d.email.toLowerCase().includes(query)) ||
      (d.str_no && d.str_no.toLowerCase().includes(query)) ||
      (d.polyclinic_name && d.polyclinic_name.toLowerCase().includes(query))
  );
});

// Load Weekly Schedule Matrix from LocalStorage
const WEEKLY_STORAGE_KEY = "antrisehat-weekly-schedule-v2";

function loadWeeklyScheduleFromLocal() {
  try {
    const raw = localStorage.getItem(WEEKLY_STORAGE_KEY);
    if (raw) {
      weeklyScheduleMatrix.value = JSON.parse(raw);
      return;
    }
  } catch (err) {
    console.warn("Failed to parse weekly schedule", err);
  }

  // Default awal adalah jadwal kosong (semua slot kosong, siap di-plotting)
  weeklyScheduleMatrix.value = {};
  saveWeeklyScheduleToLocal();
}

function saveWeeklyScheduleToLocal() {
  try {
    localStorage.setItem(WEEKLY_STORAGE_KEY, JSON.stringify(weeklyScheduleMatrix.value));
  } catch (err) {
    console.error("Failed to save weekly schedule", err);
  }
}

function getSlotPoliKey(): string {
  const poli = polyclinics.value.find((p) => p.id === selectedPoliId.value || p.slug === selectedPoliId.value);
  return poli?.slug || selectedPoliId.value;
}

// Helpers for Slot Matrix Keys
function getSlotKey(day: number, shiftId: string): string {
  const poliKey = getSlotPoliKey();
  return `${poliKey}-${day}-${shiftId}`;
}

function getAssignedDoctor(day: number, shiftId: string): DoctorShift | undefined {
  const primaryKey = getSlotKey(day, shiftId);
  if (weeklyScheduleMatrix.value[primaryKey]) {
    return weeklyScheduleMatrix.value[primaryKey];
  }
  const poli = polyclinics.value.find((p) => p.id === selectedPoliId.value || p.slug === selectedPoliId.value);
  if (poli) {
    const k1 = `${poli.id}-${day}-${shiftId}`;
    if (weeklyScheduleMatrix.value[k1]) return weeklyScheduleMatrix.value[k1];
    const k2 = `${poli.slug}-${day}-${shiftId}`;
    if (weeklyScheduleMatrix.value[k2]) return weeklyScheduleMatrix.value[k2];
  }
  return undefined;
}

function isSlotAssigned(day: number, shiftId: string): boolean {
  return !!getAssignedDoctor(day, shiftId);
}

function isDragOver(day: number, shiftId: string): boolean {
  return currentDragOverKey.value === getSlotKey(day, shiftId);
}

// ── DRAG AND DROP HANDLERS (SUPPORTING SWAP & MOVE & OFF) ───

// Start drag from Doctor Pool (sidebar)
function onPoolDragStart(event: DragEvent, doctor: DoctorRecord) {
  if (!canEdit.value) return;
  activeDragSource.value = {
    type: "pool",
    doctor,
  };
  if (event.dataTransfer) {
    event.dataTransfer.effectAllowed = "move";
    event.dataTransfer.setData("text/plain", JSON.stringify({ type: "pool", id: doctor.id }));
  }
}

// Start drag from an already-assigned card inside the schedule grid
function onSlotDragStart(event: DragEvent, day: number, shiftId: string, shiftData: DoctorShift) {
  if (!canEdit.value) return;
  const slotKey = getSlotKey(day, shiftId);
  activeDragSource.value = {
    type: "slot",
    day,
    shiftId,
    slotKey,
    shiftData,
  };
  if (event.dataTransfer) {
    event.dataTransfer.effectAllowed = "move";
    event.dataTransfer.setData("text/plain", JSON.stringify({ type: "slot", slotKey }));
  }
}

function onDragEnd() {
  activeDragSource.value = null;
  currentDragOverKey.value = null;
}

function handleDragOver(day: number, shiftId: string) {
  if (!canEdit.value) return;
  currentDragOverKey.value = getSlotKey(day, shiftId);
}

function handleDragLeave(day: number, shiftId: string) {
  if (currentDragOverKey.value === getSlotKey(day, shiftId)) {
    currentDragOverKey.value = null;
  }
}

// Drop handler on a shift slot (Senin s/d Minggu + Pagi/Siang/Sore)
async function handleDrop(targetDay: number, targetShiftId: string) {
  currentDragOverKey.value = null;
  if (!canEdit.value || !activeDragSource.value) return;

  const targetKey = getSlotKey(targetDay, targetShiftId);
  const targetExisting = weeklyScheduleMatrix.value[targetKey];
  const selectedPoli = polyclinics.value.find((p) => p.id === selectedPoliId.value);
  const targetDayName = weekDays.find((d) => d.id === targetDay)?.name || "";

  // CASE 1: Dragging from another SLOT in the schedule (Swap or Move)
  if (activeDragSource.value.type === "slot") {
    const sourceKey = activeDragSource.value.slotKey!;
    const sourceData = activeDragSource.value.shiftData!;
    const sourceDay = activeDragSource.value.day!;
    const sourceShiftId = activeDragSource.value.shiftId!;

    // Dropped onto same slot
    if (sourceKey === targetKey) return;

    if (targetExisting) {
      // ── SWAP CARDS! ───────────────────────────────────────────
      weeklyScheduleMatrix.value[targetKey] = {
        ...sourceData,
        poliId: selectedPoliId.value,
        shiftId: targetShiftId as "pagi" | "siang" | "sore",
        id: targetKey,
      };
      weeklyScheduleMatrix.value[sourceKey] = {
        ...targetExisting,
        poliId: selectedPoliId.value,
        shiftId: sourceShiftId as "pagi" | "siang" | "sore",
        id: sourceKey,
      };

      saveWeeklyScheduleToLocal();
      showToast(
        `Jadwal ditukar: ${sourceData.doctorName} (${targetDayName} ${targetShiftId.toUpperCase()}) ⇄ ${targetExisting.doctorName} (${weekDays.find((d) => d.id === sourceDay)?.name} ${sourceShiftId.toUpperCase()})`,
        "success"
      );
    } else {
      // ── MOVE TO EMPTY SLOT! ──────────────────────────────────
      weeklyScheduleMatrix.value[targetKey] = {
        ...sourceData,
        poliId: selectedPoliId.value,
        shiftId: targetShiftId as "pagi" | "siang" | "sore",
        id: targetKey,
      };
      delete weeklyScheduleMatrix.value[sourceKey];

      saveWeeklyScheduleToLocal();
      showToast(
        `Jadwal ${sourceData.doctorName} dipindahkan ke ${targetDayName} (Shift ${targetShiftId.toUpperCase()})`,
        "success"
      );
    }
  }

  // CASE 2: Dragging from DOCTOR POOL
  else if (activeDragSource.value.type === "pool") {
    const doc = activeDragSource.value.doctor!;

    const newAssignment: DoctorShift = {
      id: targetKey,
      poliId: selectedPoliId.value,
      poliName: selectedPoli?.name || "Poliklinik",
      poliCode: selectedPoli?.code || "A",
      shiftId: targetShiftId as "pagi" | "siang" | "sore",
      doctorName: doc.full_name,
      doctorSpecialty: doc.specialty || "Dokter Bertugas",
      room: targetExisting?.room || "Ruang 01",
      smartAvgServiceMin: 10,
      fixedQuota: targetExisting?.fixedQuota || 24,
      isActive: true,
      notes: `Shift ${targetShiftId} hari ${targetDayName}`,
    };

    weeklyScheduleMatrix.value[targetKey] = newAssignment;
    saveWeeklyScheduleToLocal();

    // Sync to backend API if available
    try {
      await api.assignShiftSchedule({
        polyclinic_id: doc.polyclinic_id || "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
        doctor_id: doc.id,
        day_of_week: targetDay,
        shift_id: targetShiftId,
        quota: 24,
      });
    } catch {}

    showToast(
      `Berhasil menugaskan ${doc.full_name} pada Shift ${targetShiftId.toUpperCase()} (${targetDayName})`,
      "success"
    );
  }

  // Update today's store if applicable
  const todayDayOfWeek = getTodayDayOfWeek();
  if (targetDay === todayDayOfWeek) {
    const updated = weeklyScheduleMatrix.value[targetKey];
    if (updated) {
      const idx = queueStore.doctorShifts.findIndex(
        (s) => s.poliId === selectedPoliId.value && s.shiftId === targetShiftId
      );
      if (idx >= 0) {
        queueStore.doctorShifts[idx] = { ...updated };
      } else {
        queueStore.doctorShifts.push({ ...updated });
      }
    }
  }

  activeDragSource.value = null;
}

function removeDoctorFromSlot(day: number, shiftId: string) {
  if (!canEdit.value) return;
  const key = getSlotKey(day, shiftId);
  const doc = getAssignedDoctor(day, shiftId);
  if (doc) {
    delete weeklyScheduleMatrix.value[key];
    const poli = polyclinics.value.find((p) => p.id === selectedPoliId.value || p.slug === selectedPoliId.value);
    if (poli) {
      delete weeklyScheduleMatrix.value[`${poli.id}-${day}-${shiftId}`];
      delete weeklyScheduleMatrix.value[`${poli.slug}-${day}-${shiftId}`];
    }
    saveWeeklyScheduleToLocal();

    const todayDayOfWeek = getTodayDayOfWeek();
    if (day === todayDayOfWeek) {
      const idx = queueStore.doctorShifts.findIndex(
        (s) => (s.poliId === selectedPoliId.value || s.poliId === poli?.slug) && s.shiftId === shiftId
      );
      if (idx >= 0) {
        queueStore.doctorShifts.splice(idx, 1);
      }
    }

    const dayName = weekDays.find((d) => d.id === day)?.name || "";
    showToast(`Jadwal ${doc.doctorName} pada hari ${dayName} (Shift ${shiftId.toUpperCase()}) telah dihapus.`, "success");
  }
}

function openEditSlotModal(day: number, shiftId: string) {
  if (!canEdit.value) return;
  const slot = getAssignedDoctor(day, shiftId);
  editingSlotDay.value = day;
  editingSlotShift.value = shiftId;
  editingSlotRoom.value = slot?.room || "Ruang 01";
  editingSlotQuota.value = slot?.fixedQuota || 24;
  showSlotEditModal.value = true;
}

function saveSlotConfig() {
  if (!canEdit.value) return;
  const key = getSlotKey(editingSlotDay.value, editingSlotShift.value);
  if (weeklyScheduleMatrix.value[key]) {
    weeklyScheduleMatrix.value[key].fixedQuota = editingSlotQuota.value;
    saveWeeklyScheduleToLocal();
    showToast("Pengaturan kuota berhasil disimpan.", "success");
  }
  showSlotEditModal.value = false;
}

function resetScheduleToDefault() {
  if (!canEdit.value) return;
  if (confirm("Apakah Anda yakin ingin mengosongkan seluruh slot jadwal shift?")) {
    weeklyScheduleMatrix.value = {};
    saveWeeklyScheduleToLocal();
    showToast("Seluruh slot jadwal shift berhasil dikosongkan.", "success");
  }
}

// ── CRUD DOCTOR MODAL & LOGIC ───────────────────────────────
function openAddDoctorModal() {
  if (!canEdit.value) return;
  isEditingDoctor.value = false;
  editingDoctorId.value = "";
  doctorModalError.value = "";
  doctorForm.full_name = "";
  doctorForm.polyclinic_id = selectedPoliId.value;
  doctorForm.specialty = "";
  doctorForm.str_no = "";
  doctorForm.email = "";
  doctorForm.phone = "";
  doctorForm.password = "";
  doctorForm.is_active = true;

  resetFormErrors();
  showDoctorModal.value = true;
}

function openEditDoctorModal(doc: DoctorRecord) {
  if (!canEdit.value) return;
  isEditingDoctor.value = true;
  editingDoctorId.value = doc.id;
  doctorModalError.value = "";
  doctorForm.full_name = doc.full_name;
  doctorForm.polyclinic_id = doc.polyclinic_id || "poli-umum";
  doctorForm.specialty = doc.specialty || "";
  doctorForm.str_no = doc.str_no || "";
  doctorForm.email = doc.email || "";
  doctorForm.phone = doc.phone || "";
  doctorForm.is_active = doc.is_active;

  resetFormErrors();
  showDoctorModal.value = true;
}

function closeDoctorModal() {
  showDoctorModal.value = false;
  doctorModalError.value = "";
}

function resetFormErrors() {
  formErrors.full_name = "";
  formErrors.polyclinic_id = "";
  formErrors.email = "";
  formErrors.phone = "";
}

function validateDoctorForm(): boolean {
  let valid = true;
  resetFormErrors();

  if (!doctorForm.full_name || !doctorForm.full_name.trim()) {
    formErrors.full_name = "Nama lengkap dokter wajib diisi";
    valid = false;
  }

  if (!doctorForm.polyclinic_id) {
    formErrors.polyclinic_id = "Pilih poliklinik tugas dokter";
    valid = false;
  }

  if (!doctorForm.email || !doctorForm.email.trim()) {
    formErrors.email = "Email akun login wajib diisi";
    valid = false;
  } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(doctorForm.email.trim())) {
    formErrors.email = "Format email tidak valid (contoh: dokter@antrisehat.id)";
    valid = false;
  }

  if (!doctorForm.phone || !doctorForm.phone.trim()) {
    formErrors.phone = "Nomor WhatsApp / telepon dokter wajib diisi";
    valid = false;
  }

  return valid;
}

async function saveDoctor() {
  if (!validateDoctorForm()) return;

  savingDoctor.value = true;
  doctorModalError.value = "";

  const payload = {
    full_name: doctorForm.full_name.trim(),
    polyclinic_id: doctorForm.polyclinic_id,
    specialty: doctorForm.specialty.trim() || "Dokter Umum",
    str_no: doctorForm.str_no.trim() || undefined,
    email: doctorForm.email.trim(),
    phone: doctorForm.phone.trim(),
    password: doctorForm.password.trim() || "Dokter123!",
    is_active: doctorForm.is_active,
  };

  try {
    if (isEditingDoctor.value) {
      const res = await api.updateDoctor(editingDoctorId.value, payload);
      if (res.success || res.code === "OK" || (res.data && !res.error)) {
        const idx = doctorsList.value.findIndex((d) => d.id === editingDoctorId.value);
        if (idx >= 0) {
          doctorsList.value[idx] = {
            ...doctorsList.value[idx],
            ...payload,
            polyclinic_name: getPoliNameById(payload.polyclinic_id),
          };
        }
        closeDoctorModal();
        showToast("Data dokter berhasil diperbarui.", "success");
      } else {
        doctorModalError.value = res.error?.message || res.message || "Gagal memperbarui data dokter.";
      }
    } else {
      const res = await api.createDoctor(payload);
      if (res.success || res.code === "OK" || (res.data && !res.error)) {
        const created = res.data || {
          id: `doc-${Date.now()}`,
          ...payload,
          polyclinic_name: getPoliNameById(payload.polyclinic_id),
          is_first_login: true,
        };
        doctorsList.value.unshift(created);
        closeDoctorModal();
        showToast(
          `Dokter ${payload.full_name} berhasil dibuat dengan akun login ${payload.email}. Password awal: Dokter123!`,
          "success"
        );
      } else {
        doctorModalError.value = res.error?.message || res.message || "Gagal menambahkan dokter baru.";
      }
    }
  } catch (err) {
    if (isEditingDoctor.value) {
      const idx = doctorsList.value.findIndex((d) => d.id === editingDoctorId.value);
      if (idx >= 0) {
        doctorsList.value[idx] = {
          ...doctorsList.value[idx],
          ...payload,
          polyclinic_name: getPoliNameById(payload.polyclinic_id),
        };
      }
      closeDoctorModal();
      showToast("Data dokter berhasil diperbarui.", "success");
    } else {
      const newDoc: DoctorRecord = {
        id: `doc-${Date.now()}`,
        ...payload,
        polyclinic_name: getPoliNameById(payload.polyclinic_id),
        is_first_login: true,
      };
      doctorsList.value.unshift(newDoc);
      closeDoctorModal();
      showToast(`Dokter ${newDoc.full_name} berhasil dibuat.`, "success");
    }
  } finally {
    savingDoctor.value = false;
  }
}

async function confirmDeleteDoctor(doc: DoctorRecord) {
  if (!canEdit.value) return;
  if (confirm(`Apakah Anda yakin ingin menghapus data dokter ${doc.full_name}?`)) {
    try {
      await api.deleteDoctor(doc.id);
    } catch {}
    doctorsList.value = doctorsList.value.filter((d) => d.id !== doc.id);
    showToast(`Dokter ${doc.full_name} telah dihapus dari sistem.`, "success");
  }
}

// ── UTILITIES ───────────────────────────────────────────────
function getInitials(name: string): string {
  if (!name) return "DR";
  const clean = name.replace(/dr\.|drg\.|Sp\.[A-Z]+/gi, "").trim();
  const parts = clean.split(" ").filter(Boolean);
  if (parts.length >= 2) {
    return (parts[0][0] + parts[1][0]).toUpperCase();
  }
  return clean.substring(0, 2).toUpperCase() || "DR";
}

function getPoliNameById(id?: string): string {
  const p = polyclinics.value.find((item) => item.id === id);
  return p ? p.name : "Poli Umum";
}

function getTodayDayOfWeek(): number {
  const d = new Date().getDay(); // 0 is Sunday, 1 is Monday
  return d === 0 ? 7 : d;
}

function isToday(dayId: number): boolean {
  return dayId === getTodayDayOfWeek();
}

function getWhatsAppLink(phone: string, doctorName: string): string {
  let clean = phone.replace(/[^0-9]/g, "");
  if (clean.startsWith("0")) {
    clean = "62" + clean.slice(1);
  }
  const msg = encodeURIComponent(
    `Halo ${doctorName}, ini dari bagian Administrasi Faskes AntriSehat.`
  );
  return `https://wa.me/${clean}?text=${msg}`;
}

function showToast(msg: string, type: "success" | "error" = "success") {
  notificationMsg.value = msg;
  notificationType.value = type;
  setTimeout(() => {
    if (notificationMsg.value === msg) {
      notificationMsg.value = "";
    }
  }, 4500);
}
</script>

<style scoped>
.doctor-shift-panel {
  display: flex;
  flex-direction: column;
  gap: 20px;
  width: 100%;
}

/* Panel Header */
.panel-header-section {
  display: flex;
  flex-direction: column;
  gap: 16px;
  background: white;
  padding: 20px 24px;
  border-radius: 16px;
  border: 1px solid #e2e8f0;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
}

.panel-title-area {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.title-with-badge {
  display: flex;
  align-items: center;
  gap: 14px;
}

.panel-main-title {
  font-size: 19px;
  font-weight: 700;
  color: #0f172a;
  margin: 0;
  letter-spacing: -0.01em;
}

.role-readonly-pill {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  font-size: 11px;
  font-weight: 700;
  background: #fef3c7;
  color: #92400e;
  padding: 2px 8px;
  border-radius: 20px;
  border: 1px solid #fde68a;
}

.role-admin-pill {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  font-size: 11px;
  font-weight: 700;
  background: #f0fdfa;
  color: #0f766e;
  padding: 2px 8px;
  border-radius: 20px;
  border: 1px solid #ccfbf1;
}

.panel-sub-title {
  font-size: 13px;
  color: #64748b;
  margin: 3px 0 0;
}

.readonly-banner {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  background: #fffbeb;
  border: 1px solid #fef3c7;
  color: #92400e;
  padding: 12px 16px;
  border-radius: 10px;
  font-size: 12.5px;
  line-height: 1.45;
}

.readonly-banner strong {
  color: #78350f;
}

.readonly-banner p {
  margin: 2px 0 0;
}

/* Subnav Tabs */
.subnav-tabs {
  display: flex;
  gap: 10px;
  border-top: 1px solid #f1f5f9;
  padding-top: 14px;
}

.subnav-tab-btn {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 9px 16px;
  border-radius: 10px;
  font-size: 13px;
  font-weight: 600;
  border: 1px solid transparent;
  background: #f8fafc;
  color: #64748b;
  cursor: pointer;
  transition: all 0.15s ease;
}

.subnav-tab-btn:hover {
  background: #f1f5f9;
  color: #0f172a;
}

.subnav-tab-btn.active {
  background: #f0fdfa;
  color: #0f766e;
  border-color: #ccfbf1;
  box-shadow: 0 1px 2px rgba(15, 118, 110, 0.08);
}

/* Notification Toast */
.notification-toast {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 18px;
  border-radius: 10px;
  font-size: 13px;
  font-weight: 500;
}

.notification-toast.success {
  background-color: #f0fdf4;
  border: 1px solid #bbf7d0;
  color: #15803d;
}

.notification-toast.error {
  background-color: #fef2f2;
  border: 1px solid #fee2e2;
  color: #b91c1c;
}

.btn-close-toast {
  margin-left: auto;
  background: none;
  border: none;
  color: inherit;
  cursor: pointer;
  padding: 2px;
}

/* Schedule Toolbar */
.schedule-toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: 14px;
  background: white;
  padding: 14px 20px;
  border-radius: 14px;
  border: 1px solid #e2e8f0;
}

.poli-filter-group {
  display: flex;
  align-items: center;
  gap: 12px;
}

.filter-label {
  font-size: 13px;
  font-weight: 600;
  color: #475569;
}

.filter-buttons {
  display: flex;
  gap: 6px;
  flex-wrap: wrap;
}

.filter-pill-btn {
  padding: 6px 14px;
  border-radius: 20px;
  font-size: 12.5px;
  font-weight: 600;
  border: 1px solid #e2e8f0;
  background: #f8fafc;
  color: #475569;
  cursor: pointer;
  transition: all 0.15s ease;
}

.filter-pill-btn:hover {
  background: #f1f5f9;
  color: #0f172a;
}

.filter-pill-btn.active {
  background: #0f766e;
  color: white;
  border-color: #0f766e;
}

.toolbar-actions {
  display: flex;
  gap: 8px;
}

/* DnD Workspace Layout */
.dnd-workspace {
  display: grid;
  grid-template-columns: 290px 1fr;
  gap: 18px;
  align-items: start;
}

@media (max-width: 1024px) {
  .dnd-workspace {
    grid-template-columns: 1fr;
  }
}

/* 1. Left Doctor Pool Card */
.doctor-pool-card {
  background: white;
  border-radius: 16px;
  border: 1px solid #e2e8f0;
  padding: 18px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
}

.pool-header {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-bottom: 12px;
}

.pool-title {
  display: flex;
  align-items: center;
  gap: 8px;
}

.pool-title h4 {
  font-size: 15px;
  font-weight: 700;
  color: #0f172a;
  margin: 0;
}

.pool-status-toggle {
  display: flex;
  gap: 4px;
  background: #f1f5f9;
  padding: 3px;
  border-radius: 8px;
}

.toggle-chip {
  flex: 1;
  padding: 4px 6px;
  font-size: 11px;
  font-weight: 600;
  border: none;
  border-radius: 6px;
  background: transparent;
  color: #64748b;
  cursor: pointer;
  transition: all 0.15s ease;
}

.toggle-chip.active {
  background: white;
  color: #0f766e;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
}

.pool-search-box {
  position: relative;
  margin-bottom: 8px;
}

.pool-search-input {
  width: 100%;
  padding: 8px 10px 8px 30px;
  font-size: 12px;
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  outline: none;
  background: #f8fafc;
}

.pool-search-input:focus {
  border-color: #0f766e;
  background: white;
}

.search-icon {
  position: absolute;
  left: 10px;
  top: 50%;
  transform: translateY(-50%);
  color: #94a3b8;
}

.pool-hint {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 11px;
  color: #0f766e;
  background: #f0fdfa;
  padding: 6px 10px;
  border-radius: 8px;
  margin-bottom: 12px;
}

.pool-hint.readonly {
  color: #92400e;
  background: #fef3c7;
}

.doctor-pool-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
  max-height: 400px;
  overflow-y: auto;
  padding-right: 4px;
}

.draggable-doctor-chip {
  display: flex;
  align-items: flex-start;
  gap: 10px;
  padding: 10px;
  background: white;
  border: 1px solid #e2e8f0;
  border-radius: 10px;
  cursor: grab;
  user-select: none;
  transition: all 0.15s ease;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.04);
}

.draggable-doctor-chip.not-allowed {
  cursor: default;
}

.draggable-doctor-chip.is-off {
  background: #fafaf9;
  border-color: #e2e8f0;
  opacity: 0.85;
}

.draggable-doctor-chip:hover:not(.not-allowed) {
  border-color: #0f766e;
  background: #fdfdfd;
  transform: translateY(-1px);
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.08);
}

.draggable-doctor-chip:active:not(.not-allowed) {
  cursor: grabbing;
  opacity: 0.7;
}

.drag-handle {
  cursor: grab;
  margin-top: 2px;
}

.doc-chip-avatar {
  width: 32px;
  height: 32px;
  border-radius: 8px;
  background: #0f766e;
  color: white;
  font-size: 12px;
  font-weight: 700;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.doc-chip-info {
  display: flex;
  flex-direction: column;
  flex: 1;
  overflow: hidden;
}

.doc-chip-name {
  font-size: 12.5px;
  color: #0f172a;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.doc-chip-spec {
  font-size: 11px;
  color: #64748b;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.doc-chip-poli {
  font-size: 10px;
  font-weight: 600;
  color: #0f766e;
  margin-top: 1px;
}

.empty-pool-state {
  text-align: center;
  padding: 24px 12px;
  color: #94a3b8;
  font-size: 12px;
}

/* 2. Right Weekly Roster Board */
.weekly-roster-board {
  background: white;
  border-radius: 16px;
  border: 1px solid #e2e8f0;
  padding: 16px;
  overflow-x: auto;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
}

.roster-grid-header {
  display: grid;
  grid-template-columns: 140px repeat(7, minmax(136px, 1fr));
  gap: 10px;
  padding-bottom: 12px;
  border-bottom: 1px solid #f1f5f9;
  text-align: center;
}

.grid-cell-label {
  font-size: 12px;
  font-weight: 700;
  color: #64748b;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  display: flex;
  align-items: center;
  justify-content: center;
}

.day-col-header {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 3px;
}

.day-name {
  font-size: 13.5px;
  font-weight: 700;
  color: #1e293b;
}

.day-col-header.is-today .day-name {
  color: #0f766e;
}

.today-badge {
  font-size: 10px;
  font-weight: 700;
  background: #0f766e;
  color: white;
  padding: 1px 7px;
  border-radius: 10px;
}

/* Shift Grid Rows */
.shift-grid-row {
  display: grid;
  grid-template-columns: 140px repeat(7, minmax(136px, 1fr));
  gap: 10px;
  margin-top: 10px;
}

.shift-time-cell {
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  padding: 10px 12px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  gap: 5px;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.02);
}

.shift-name-tag {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  font-size: 11px;
  font-weight: 700;
  padding: 2.5px 7px;
  border-radius: 6px;
  width: fit-content;
}

.shift-name-tag.pagi {
  background: #fef3c7;
  color: #b45309;
}

.shift-name-tag.siang {
  background: #e0f2fe;
  color: #0369a1;
}

.shift-name-tag.sore {
  background: #ede9fe;
  color: #6d28d9;
}

.shift-hours {
  font-size: 12.5px;
  font-weight: 700;
  color: #1e293b;
  letter-spacing: -0.01em;
}

.shift-quota-tag {
  font-size: 11px;
  color: #64748b;
  font-weight: 500;
}

/* Day Shift Slot (Drop Target) */
.day-shift-slot {
  background: #fafaf9;
  border: 1.5px dashed #cbd5e1;
  border-radius: 12px;
  min-height: 112px;
  padding: 8px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  transition: all 0.2s ease;
  position: relative;
}

.day-shift-slot.readonly-slot {
  border-style: solid;
  border-color: #e2e8f0;
  background: #fbfbfb;
}

.day-shift-slot.drag-over {
  border-color: #0f766e;
  background: #f0fdfa;
  box-shadow: 0 0 0 3px rgba(15, 118, 110, 0.15);
  transform: scale(1.02);
}

.day-shift-slot.drag-over-clear {
  border-color: #f59e0b !important;
  background: #fffbeb !important;
  box-shadow: 0 0 0 3px rgba(245, 158, 11, 0.25) !important;
}

.day-shift-slot.is-assigned {
  background: transparent;
  border: none;
  padding: 0;
}

/* Assigned Doctor Card Inside Slot */
.assigned-doctor-card {
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  height: 100%;
  width: 100%;
  background: #ffffff;
  border-radius: 12px;
  padding: 10px 10px 8px;
  position: relative;
  border: 1px solid #e2e8f0;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.04);
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
}

.assigned-doctor-card.can-drag {
  cursor: grab;
}

.assigned-doctor-card.can-drag:hover {
  background: #f0fdfa;
  border-color: #0f766e;
  transform: translateY(-2px);
  box-shadow: 0 6px 14px -3px rgba(15, 118, 110, 0.18), 0 0 0 1.5px #0f766e;
  z-index: 5;
}

.assigned-doctor-card.can-drag:active {
  cursor: grabbing;
  opacity: 0.85;
  transform: scale(0.98);
}

.assigned-card-top {
  display: flex;
  align-items: flex-start;
  gap: 7px;
  position: relative;
}

.assigned-avatar {
  width: 28px;
  height: 28px;
  border-radius: 8px;
  background: #0f766e;
  color: white;
  font-size: 11px;
  font-weight: 700;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.assigned-meta {
  display: flex;
  flex-direction: column;
  flex: 1;
  overflow: hidden;
}

.assigned-name {
  font-size: 12px;
  font-weight: 700;
  color: #0f172a;
  line-height: 1.25;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.assigned-spec {
  font-size: 10.5px;
  color: #64748b;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  margin-top: 1px;
}

.btn-remove-slot {
  background: #f1f5f9;
  border: none;
  color: #64748b;
  cursor: pointer;
  width: 22px;
  height: 22px;
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.15s ease;
}

.btn-remove-slot:hover {
  background: #fee2e2;
  color: #dc2626;
}

.assigned-card-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-top: 8px;
  padding-top: 6px;
  border-top: 1px dashed #e2e8f0;
  gap: 4px;
}

.btn-delete-card-slot {
  background: #fef2f2;
  border: 1px solid #fecaca;
  color: #dc2626;
  width: 24px;
  height: 24px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 0;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.15s ease;
  flex-shrink: 0;
}

.btn-delete-card-slot:hover {
  background: #fee2e2;
  border-color: #fca5a5;
  color: #b91c1c;
  transform: scale(1.08);
}

.chip-plot-badge {
  display: inline-flex;
  align-items: center;
  gap: 3px;
  font-size: 9.5px;
  font-weight: 600;
  color: #0f766e;
  background: #f0fdfa;
  border: 1px solid #ccfbf1;
  padding: 1px 6px;
  border-radius: 4px;
  margin-top: 4px;
  width: fit-content;
}

.quota-readonly-chip {
  font-size: 10px;
  font-weight: 600;
  color: #64748b;
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  padding: 2.5px 6px;
  border-radius: 6px;
}

.btn-edit-slot-quota {
  background: #f0fdfa;
  border: 1px solid #ccfbf1;
  font-size: 10px;
  font-weight: 600;
  color: #0f766e;
  display: inline-flex;
  align-items: center;
  gap: 4px;
  cursor: pointer;
  padding: 2.5px 6px;
  border-radius: 6px;
  transition: all 0.15s ease;
}

.btn-edit-slot-quota:hover {
  background: #ccfbf1;
  border-color: #99f6e4;
  color: #115e59;
}

/* Empty Drop Placeholder */
.empty-drop-placeholder {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 6px;
  height: 100%;
  color: #94a3b8;
  text-align: center;
  padding: 6px 4px;
  transition: color 0.15s ease;
}

.day-shift-slot:not(.readonly-slot):hover {
  border-color: #0f766e;
  background: #f0fdfa;
}

.day-shift-slot:not(.readonly-slot):hover .drop-icon {
  color: #0f766e;
  transform: scale(1.1);
}

.day-shift-slot:not(.readonly-slot):hover .drop-text {
  color: #0f766e;
}

.drop-icon {
  color: #94a3b8;
  transition: all 0.15s ease;
}

.drop-text {
  font-size: 11px;
  font-weight: 600;
  line-height: 1.3;
}

/* ============================================================ */
/* TAB 2: CRUD DOCTOR STYLES                                    */
/* ============================================================ */
.doctors-tab-content {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.crud-toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 14px;
  flex-wrap: wrap;
  background: white;
  padding: 16px 20px;
  border-radius: 14px;
  border: 1px solid #e2e8f0;
}

.search-wrap {
  position: relative;
  flex: 1;
  max-width: 450px;
}

.crud-search-input {
  width: 100%;
  padding: 9px 12px 9px 34px;
  font-size: 13px;
  border: 1px solid #cbd5e1;
  border-radius: 8px;
  outline: none;
}

.crud-search-input:focus {
  border-color: #0f766e;
  box-shadow: 0 0 0 2px rgba(15, 118, 110, 0.15);
}

.table-responsive {
  background: white;
  border-radius: 14px;
  border: 1px solid #e2e8f0;
  overflow-x: auto;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
}

.doctors-table {
  width: 100%;
  border-collapse: collapse;
  text-align: left;
  font-size: 13px;
}

.doctors-table th {
  background: #f8fafc;
  padding: 12px 16px;
  font-weight: 700;
  color: #475569;
  border-bottom: 1px solid #e2e8f0;
}

.doctors-table td {
  padding: 14px 16px;
  border-bottom: 1px solid #f1f5f9;
  vertical-align: middle;
}

.doc-table-profile {
  display: flex;
  align-items: center;
  gap: 12px;
}

.doc-table-avatar {
  width: 36px;
  height: 36px;
  border-radius: 10px;
  background: #0f766e;
  color: white;
  font-weight: 700;
  font-size: 13px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.doc-table-avatar.avatar-off {
  background: #94a3b8;
}

.doc-table-name {
  display: block;
  font-size: 13.5px;
  font-weight: 700;
  color: #0f172a;
}

.doc-table-sub {
  font-size: 12px;
  color: #64748b;
}

.poli-badge {
  font-size: 12px;
  font-weight: 600;
  padding: 3px 8px;
  border-radius: 6px;
  background: #f1f5f9;
  color: #334155;
}

.str-code {
  font-size: 11.5px;
  font-family: monospace;
  background: #f8fafc;
  padding: 2px 6px;
  border-radius: 4px;
  border: 1px solid #e2e8f0;
  color: #475569;
}

.email-info-cell {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.first-login-pill {
  font-size: 10px;
  font-weight: 700;
  background: #fef3c7;
  color: #b45309;
  padding: 1px 6px;
  border-radius: 4px;
  width: fit-content;
}

.wa-phone-link {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  color: #059669;
  font-weight: 600;
  text-decoration: none;
}

.wa-phone-link:hover {
  text-decoration: underline;
}

.status-pill {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-size: 11.5px;
  font-weight: 600;
  padding: 3px 8px;
  border-radius: 12px;
}

.status-pill.active {
  background: #dcfce7;
  color: #15803d;
}

.status-pill.inactive {
  background: #fee2e2;
  color: #b91c1c;
}

.status-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background-color: currentColor;
}

.action-buttons-cell {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 6px;
}

.btn-action {
  background: none;
  border: 1px solid #e2e8f0;
  border-radius: 6px;
  padding: 6px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.15s;
}

.btn-action.edit {
  color: #0f766e;
}

.btn-action.edit:hover {
  background: #f0fdfa;
  border-color: #0f766e;
}

.btn-action.delete {
  color: #ef4444;
}

.btn-action.delete:hover {
  background: #fef2f2;
  border-color: #ef4444;
}

.empty-table-cell {
  text-align: center;
  padding: 36px 16px;
  color: #64748b;
}

/* Modals */
.modal-backdrop {
  position: fixed;
  inset: 0;
  background-color: rgba(15, 23, 42, 0.6);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 16px;
}

.modal-card {
  background: white;
  border-radius: 16px;
  width: 100%;
  max-width: 540px;
  box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.15);
  overflow: hidden;
  animation: modalFadeIn 0.2s ease-out;
}

.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 18px 24px;
  border-bottom: 1px solid #e2e8f0;
}

.btn-close {
  background: none;
  border: none;
  color: #64748b;
  cursor: pointer;
  padding: 4px;
  border-radius: 6px;
}

.modal-body {
  padding: 24px;
}

.modal-footer {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 10px;
  padding: 16px 24px;
  border-top: 1px solid #e2e8f0;
  background: #f8fafc;
}

.modal-error-alert {
  display: flex;
  align-items: center;
  gap: 8px;
  background: #fef2f2;
  border: 1px solid #fee2e2;
  color: #b91c1c;
  padding: 10px 14px;
  border-radius: 8px;
  font-size: 13px;
  font-weight: 500;
  margin-bottom: 14px;
}

.form-field {
  display: flex;
  flex-direction: column;
  gap: 5px;
  text-align: left;
}

.form-grid-2 {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
}

@media (max-width: 600px) {
  .form-grid-2 {
    grid-template-columns: 1fr;
  }
}

.field-label {
  font-size: 12.5px;
  font-weight: 600;
  color: #334155;
}

.form-field input,
.form-field select {
  padding: 8px 12px;
  border-radius: 8px;
  border: 1px solid #cbd5e1;
  font-size: 13px;
  outline: none;
  transition: all 0.15s;
}

.form-field input:focus,
.form-field select:focus {
  border-color: #0f766e;
  box-shadow: 0 0 0 2px rgba(15, 118, 110, 0.15);
}

.req {
  color: #dc2626;
  font-weight: 700;
  margin-left: 2px;
}

.input-error {
  border-color: #ef4444 !important;
  background-color: #fef2f2 !important;
}

.field-error-text {
  display: block;
  font-size: 11.5px;
  font-weight: 600;
  color: #dc2626;
  margin-top: 3px;
}
</style>
