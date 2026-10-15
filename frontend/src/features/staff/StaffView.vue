<template>
  <main class="dashboard">
    <aside class="sidebar" :data-collapsed="sidebarCollapsed">
      <button
        type="button"
        class="sidebar-toggle"
        :aria-label="
          sidebarCollapsed
            ? 'Tampilkan menu sidebar'
            : 'Sembunyikan menu sidebar'
        "
        :aria-expanded="!sidebarCollapsed"
        :title="
          sidebarCollapsed
            ? 'Tampilkan menu sidebar'
            : 'Sembunyikan menu sidebar'
        "
        @click="sidebarCollapsed = !sidebarCollapsed"
      >
        <Menu :size="20" aria-hidden="true" />
      </button>
      <button
        :class="{ selected: queueStore.role === 'loket' }"
        @click="queueStore.role = 'loket'"
      >
        <Users :size="20" /> Dashboard Loket
      </button>
      <button
        :class="{ selected: queueStore.role === 'medis' }"
        @click="queueStore.role = 'medis'"
      >
        <Stethoscope :size="20" /> Dashboard Poli
      </button>
      <button @click="goToTv"><Monitor :size="20" /> Layar TV Publik</button>
      <div class="sidebar-bottom">
        <button @click="goToHome"><LogOut :size="18" /> Keluar</button>
      </div>
    </aside>

    <div class="dashboard-body">
      <div class="dashboard-heading">
        <div>
          <div class="section-label">PUSKESMAS SEHAT BERSAMA</div>
          <h2>
            {{
              queueStore.role === "loket"
                ? "Dashboard Loket"
                : "Ruang Pemeriksaan"
            }}
          </h2>
          <p>
            {{ currentFormattedDate }}
          </p>
        </div>
        <button class="button outline" @click="showPendingNotice">
          <Bell :size="18" /> Notifikasi
        </button>
      </div>

      <!-- Stats Grid -->
      <div class="stats-grid">
        <div class="stat-card">
          <span>Pendaftaran hari ini</span>
          <strong>{{ todayRegistrationsCount }}</strong>
          <Activity :size="22" />
        </div>
        <div class="stat-card">
          <span>Menunggu verifikasi</span>
          <strong>{{ pendingVerificationCount }}</strong>
          <Activity :size="22" />
        </div>
        <div class="stat-card">
          <span>Siap dilayani</span>
          <strong>{{ readyToServeCount }}</strong>
          <Activity :size="22" />
        </div>
        <div class="stat-card">
          <span>Kunjungan selesai</span>
          <strong>{{ finishedCount }}</strong>
          <Activity :size="22" />
        </div>
      </div>

      <!-- Role: Loket View -->
      <section v-if="queueStore.role === 'loket'" class="dashboard-panel">
        <div class="panel-heading">
          <h3>Antrean hari ini</h3>
          <div style="display: flex; gap: 8px">
            <button
              class="button outline small"
              @click="handleWalkinRegistration"
            >
              + Pasien Walk-in
            </button>
            <button class="button primary small" @click="goToScan">
              <QrCode :size="17" /> Scan QR Pasien
            </button>
          </div>
        </div>
        <div class="search-field">
          <Search :size="18" />
          <input
            placeholder="Cari nama, NIK, atau nomor antrean..."
            v-model="queueStore.search"
          />
        </div>
        <div class="table-scroll">
          <table>
            <thead>
              <tr>
                <th>Antrean</th>
                <th>Pasien</th>
                <th>Poli</th>
                <th>Pendaftaran</th>
                <th>Status</th>
                <th>Tindakan</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="p in filteredLoketPatients" :key="p.id">
                <td>
                  <strong class="queue-number">
                    {{ p.number }}
                  </strong>
                </td>
                <td>
                  <strong>{{ p.name }}</strong>
                  <small>{{ p.nik }}</small>
                </td>
                <td>{{ p.poli }}</td>
                <td>
                  <span :class="['badge', { coral: p.source !== 'Online' }]">
                    {{ p.source }}
                  </span>
                </td>
                <td>
                  <span
                    :class="[
                      'status',
                      `status-${p.status === 'Selesai' ? 'done' : 'pending'}`,
                    ]"
                  >
                    {{ p.status }}
                  </span>
                </td>
                <td>
                  <button
                    v-if="
                      ['Terdaftar', 'Menunggu verifikasi'].includes(p.status)
                    "
                    class="text-button"
                    @click="queueStore.verify = p"
                  >
                    Verifikasi <ChevronRight :size="15" />
                  </button>
                  <button
                    v-else-if="p.status === 'Hadir'"
                    class="text-button"
                    @click="queueStore.call(p)"
                  >
                    Panggil <Volume2 :size="16" />
                  </button>
                  <span v-else class="muted">—</span>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </section>

      <!-- Role: Medis / Poli View -->
      <template v-else>
        <div class="poli-toolbar">
          <label>
            Poli Anda
            <select
              v-model="queueStore.staffPoli"
              @change="queueStore.selected = ''"
            >
              <option v-for="s in services" :key="s.name">
                {{ s.name }}
              </option>
            </select>
          </label>
        </div>
        <div class="medical-grid">
          <section class="dashboard-panel patient-list">
            <h3>Antrean {{ queueStore.staffPoli }}</h3>
            <button
              v-for="p in poliPatients"
              :key="p.id"
              :class="['patient', { selected: queueStore.selected === p.id }]"
              @click="selectPatient(p)"
            >
              <strong>
                {{ p.number }}
                <span>{{ p.status }}</span>
              </strong>
              <b>{{ p.name }}</b>
              <small> {{ p.insurance }} · {{ p.gender }} </small>
            </button>
            <p v-if="!poliPatients.length" class="muted">
              Belum ada pasien yang terverifikasi.
            </p>
          </section>

          <section class="dashboard-panel soap-panel">
            <div v-if="activePatient" :key="activePatient.id">
              <div class="record-banner">
                <div class="record-identity record-identity-compact">
                  <div class="record-patient-info">
                    <span class="record-avatar" aria-hidden="true">
                      {{ getInitials(activePatient.name) }}
                    </span>
                    <div class="record-patient-details">
                      <h3>{{ activePatient.name }}</h3>
                      <p>
                        No. RM: belum tersedia ·
                        <span class="badge">{{ activePatient.gender }}</span>
                      </p>
                      <p>
                        {{ formatPatientAge(activePatient.birth) }}
                      </p>
                      <p v-if="activePatient.insurance === 'BPJS'">
                        BPJS: {{ activePatient.bpjs || "Belum diisi" }}
                      </p>
                    </div>
                  </div>
                  <button
                    type="button"
                    class="button primary small record-trigger"
                    aria-label="Lihat Riwayat Rekam Medis"
                    title="Lihat Riwayat Rekam Medis"
                    @click="goToMedicalRecord"
                  >
                    <Activity :size="17" aria-hidden="true" />
                    <span class="record-trigger-label">MEDICAL RECORD</span>
                  </button>
                </div>
                <div
                  v-if="patientAllergies.length"
                  class="record-alert"
                  role="alert"
                >
                  PERINGATAN ALERGI: {{ patientAllergies.join(", ") }}
                </div>
              </div>

              <div class="record-layout">
                <div class="record-work">
                  <div class="summary-box">
                    <strong>Pra-anamnesis pasien</strong>
                    <p>{{ activePatient.complaint }}</p>
                    <small>
                      Gejala:
                      {{ activePatient.symptoms.join(", ") || "Tidak diisi" }} ·
                      Alergi: {{ activePatient.allergy || "Tidak ada" }}
                    </small>
                  </div>

                  <form @submit.prevent="saveMedicalRecord">
                    <h3 class="field-heading">Tanda vital</h3>
                    <div class="vital-grid">
                      <label v-for="field in vitalFields" :key="field.name">
                        {{ field.label }}
                        <input
                          required
                          type="number"
                          v-model="vitals[field.name]"
                          :min="field.min"
                          :max="field.max"
                          step="any"
                          :placeholder="field.placeholder"
                        />
                      </label>
                    </div>

                    <label>
                      Dokter pemeriksa
                      <input
                        required
                        v-model="doctorName"
                        placeholder="Nama lengkap dokter"
                        autocomplete="off"
                      />
                    </label>

                    <label>
                      S · Subjektif
                      <textarea required rows="2" v-model="subjectiveText" />
                    </label>

                    <label>
                      O · Objektif
                      <textarea required rows="2" v-model="objectiveText" />
                    </label>

                    <label>
                      A · Diagnosis ICD-10
                      <input
                        required
                        list="diagnoses-list"
                        v-model="queueStore.soap.diagnosis"
                        placeholder="Cari kode atau diagnosis..."
                      />
                      <datalist id="diagnoses-list">
                        <option value="J00 — Acute nasopharyngitis" />
                        <option
                          value="J06.9 — Infeksi saluran pernapasan atas"
                        />
                        <option value="R50.9 — Demam" />
                        <option value="I10 — Hipertensi esensial" />
                        <option value="E11 — Diabetes melitus tipe 2" />
                        <option value="K02.9 — Karies gigi" />
                        <option value="R51 — Sakit kepala" />
                      </datalist>
                    </label>

                    <label>
                      P · Rencana tindakan / terapi
                      <textarea
                        required
                        rows="3"
                        v-model="queueStore.soap.plan"
                      />
                    </label>

                    <button class="button primary full-width" type="submit">
                      <CheckCircle2 :size="18" /> Simpan & Tutup Kunjungan
                    </button>
                  </form>
                </div>
              </div>
            </div>
            <div v-else class="empty-state">
              <Stethoscope :size="44" />
              <h3>Pilih pasien untuk pemeriksaan</h3>
              <p>Pasien yang telah diverifikasi akan muncul di panel kiri.</p>
            </div>
          </section>
        </div>
      </template>
    </div>

    <!-- Modals -->
    <VerifyModal />
  </main>
</template>

<script setup lang="ts">
import { ref, reactive, computed } from "vue";
import { useRouter } from "vue-router";
import {
  Menu,
  Users,
  Stethoscope,
  Monitor,
  LogOut,
  Bell,
  Activity,
  QrCode,
  Search,
  ChevronRight,
  Volume2,
  CheckCircle2,
} from "lucide-vue-next";
import { services } from "../../constants/services";
import { todayDate } from "../../constants/initialData";
import { useQueueStore } from "../../stores/queue";
import type { Patient } from "../../types";
import VerifyModal from "./VerifyModal.vue";

const router = useRouter();
const queueStore = useQueueStore();

const sidebarCollapsed = ref(false);
const doctorName = ref("dr. Maya Putri");
const subjectiveText = ref("");
const objectiveText = ref("");

const vitals = reactive<Record<string, string>>({
  bp_sys: "120",
  bp_dia: "80",
  hr: "80",
  temp: "36.5",
  spo2: "98",
  weight: "60",
  height: "165",
});

const vitalFields = [
  {
    name: "bp_sys",
    label: "Sistolik (mmHg)",
    placeholder: "120",
    min: "40",
    max: "300",
  },
  {
    name: "bp_dia",
    label: "Diastolik (mmHg)",
    placeholder: "80",
    min: "20",
    max: "200",
  },
  {
    name: "hr",
    label: "Detak jantung (/menit)",
    placeholder: "80",
    min: "20",
    max: "300",
  },
  {
    name: "temp",
    label: "Suhu (°C)",
    placeholder: "36.5",
    min: "25",
    max: "45",
  },
  { name: "spo2", label: "SpO₂ (%)", placeholder: "98", min: "1", max: "100" },
  {
    name: "weight",
    label: "Berat badan (kg)",
    placeholder: "60",
    min: "0.1",
    max: "500",
  },
  {
    name: "height",
    label: "Tinggi badan (cm)",
    placeholder: "165",
    min: "20",
    max: "260",
  },
];

const currentFormattedDate = computed(() =>
  new Date().toLocaleDateString("id-ID", { dateStyle: "full" }),
);

const todayRegistrationsCount = computed(
  () => queueStore.patients.filter((p) => p.date === todayDate).length,
);

const pendingVerificationCount = computed(
  () =>
    queueStore.patients.filter(
      (p) =>
        p.date === todayDate &&
        ["Terdaftar", "Menunggu verifikasi"].includes(p.status),
    ).length,
);

const readyToServeCount = computed(
  () =>
    queueStore.patients.filter(
      (p) => p.date === todayDate && p.status === "Hadir",
    ).length,
);

const finishedCount = computed(
  () =>
    queueStore.patients.filter(
      (p) => p.date === todayDate && p.status === "Selesai",
    ).length,
);

const filteredLoketPatients = computed(() =>
  queueStore.patients.filter(
    (p) =>
      p.date === todayDate &&
      `${p.name} ${p.nik} ${p.number}`
        .toLowerCase()
        .includes(queueStore.search.toLowerCase()),
  ),
);

const poliPatients = computed(() =>
  queueStore.patients.filter(
    (p) =>
      p.poli === queueStore.staffPoli &&
      p.date === todayDate &&
      ["Hadir", "Dipanggil"].includes(p.status),
  ),
);

const activePatient = computed(() => {
  const p = queueStore.chosen;
  if (
    p &&
    p.poli === queueStore.staffPoli &&
    ["Hadir", "Dipanggil"].includes(p.status)
  ) {
    return p;
  }
  return null;
});

const patientAllergies = computed(() => {
  if (!activePatient.value) return [];
  const chosen = activePatient.value;
  const allergies = [
    ...new Set(
      queueStore.patients
        .filter(
          (p) => p.id === chosen.id || (chosen.nik && p.nik === chosen.nik),
        )
        .map((p) => p.allergy?.trim())
        .filter((a) => a && !/^(tidak ada|tidak|nihil|none|-)$/i.test(a)),
    ),
  ];
  return allergies;
});

function getInitials(name: string) {
  return name
    .split(" ")
    .slice(0, 2)
    .map((n) => n[0])
    .join("");
}

function formatPatientAge(birthStr: string) {
  if (!birthStr) return "Tanggal lahir belum tersedia";
  try {
    const birth = new Date(birthStr);
    const today = new Date();
    let age = today.getFullYear() - birth.getFullYear();
    const m = today.getMonth() - birth.getMonth();
    if (m < 0 || (m === 0 && today.getDate() < birth.getDate())) {
      age--;
    }
    return `${birth.toLocaleDateString("id-ID")} · ${Math.max(0, age)} tahun`;
  } catch {
    return "Tanggal lahir belum tersedia";
  }
}

function selectPatient(p: Patient) {
  queueStore.selected = p.id;
  subjectiveText.value = [
    p.complaint,
    p.symptoms.length ? `Gejala: ${p.symptoms.join(", ")}` : "",
  ]
    .filter(Boolean)
    .join("\n");
  objectiveText.value = "";
  queueStore.soap.diagnosis = "";
  queueStore.soap.plan = "";
}

function saveMedicalRecord() {
  if (!activePatient.value) return;
  const record = {
    ...vitals,
    pressure: `${vitals.bp_sys}/${vitals.bp_dia}`,
    temperature: vitals.temp,
    pulse: vitals.hr,
    doctor: doctorName.value,
    subjective: subjectiveText.value,
    objective: objectiveText.value,
    diagnosis: queueStore.soap.diagnosis,
    plan: queueStore.soap.plan,
    recordedAt: new Date().toISOString(),
  };

  queueStore.patients = queueStore.patients.map((p) =>
    p.id === activePatient.value?.id
      ? {
          ...p,
          status: "Selesai",
          soap: JSON.stringify(record),
        }
      : p,
  );

  queueStore.selected = "";
  queueStore.setNotice("Pemeriksaan disimpan dan kunjungan ditutup.");
}

function showPendingNotice() {
  queueStore.setNotice(
    `${pendingVerificationCount.value} pasien menunggu verifikasi`,
  );
}

function handleWalkinRegistration() {
  queueStore.startRegistration(queueStore.staffPoli || "Poli Umum", true);
  router.push("/register");
  window.scrollTo(0, 0);
}

function goToScan() {
  router.push("/scan");
  window.scrollTo(0, 0);
}

function goToTv() {
  router.push("/tv");
  window.scrollTo(0, 0);
}

function goToHome() {
  router.push("/");
  window.scrollTo(0, 0);
}

function goToMedicalRecord() {
  router.push("/medical-record");
  window.scrollTo(0, 0);
}
</script>

<style scoped>
.record-banner {
  margin-bottom: 24px;
}

.record-identity {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  width: 100%;
}

.record-patient-info {
  display: flex;
  align-items: flex-start;
  gap: 14px;
}

.record-avatar {
  width: 44px;
  height: 44px;
  border-radius: 50%;
  background-color: #1e464c;
  color: #ffffff;
  font-weight: 700;
  font-size: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  letter-spacing: 0.5px;
}

.record-patient-details h3 {
  margin: 0 0 4px 0;
  font-size: 18px;
  font-weight: 700;
  color: #111827;
  line-height: 1.25;
}

.record-patient-details p {
  margin: 2px 0 0 0;
  font-size: 13px;
  color: #64748b;
  line-height: 1.4;
}

.record-trigger {
  margin-left: auto;
  align-self: center;
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 8px 18px;
  border-radius: 9999px;
  font-size: 12px;
  font-weight: 700;
  letter-spacing: 0.5px;
  background-color: #1e464c;
  color: #ffffff;
  border: none;
  cursor: pointer;
  white-space: nowrap;
  transition: all 0.2s ease;
  box-shadow: 0 2px 6px rgba(30, 70, 76, 0.2);
}

.record-trigger:hover {
  background-color: #13353a;
  transform: translateY(-1px);
}

.record-alert {
  background-color: #fef2f2;
  border: 1px solid #fee2e2;
  color: #b91c1c;
  border-radius: 8px;
  padding: 10px 14px;
  font-size: 12px;
  font-weight: 600;
  margin-top: 14px;
}

@media (max-width: 640px) {
  .record-identity {
    flex-direction: column;
    align-items: stretch;
  }

  .record-trigger {
    margin-left: 0;
    align-self: stretch;
    justify-content: center;
    margin-top: 8px;
  }
}
</style>
