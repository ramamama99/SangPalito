<template>
  <main class="medical-record-page">
    <div class="medical-record-wrapper">
      <!-- Back Button -->
      <div class="record-top-nav">
        <RouterLink to="/staff" class="back-button">
          <ChevronLeft :size="16" />
          <span>Kembali ke Pemeriksaan</span>
        </RouterLink>
      </div>

      <div v-if="chosenPatient" class="record-content">
        <!-- Patient Identity Header -->
        <header class="patient-header">
          <div class="patient-avatar" aria-hidden="true">
            {{ getInitials(chosenPatient.name) }}
          </div>
          <div class="patient-details">
            <h1 class="patient-name">{{ chosenPatient.name }}</h1>
            <div class="patient-meta-row">
              <span class="patient-rm">No. RM: belum tersedia</span>
              <span class="meta-dot">·</span>
              <span class="gender-badge">{{ chosenPatient.gender }}</span>
            </div>
            <div class="patient-meta-sub">
              {{ formatPatientAge(chosenPatient.birth) }}
            </div>
            <div
              v-if="chosenPatient.insurance === 'BPJS'"
              class="patient-meta-sub"
            >
              BPJS: {{ chosenPatient.bpjs || "Belum diisi" }}
            </div>
          </div>
        </header>

        <!-- Allergy Alert if any -->
        <div v-if="allergies.length" class="allergy-alert" role="alert">
          <strong>PERINGATAN ALERGI:</strong> {{ allergies.join(", ") }}
        </div>

        <!-- Main Medical History Card -->
        <section class="history-card" aria-label="Riwayat rekam medis">
          <h2 class="history-title">Riwayat rekam medis</h2>

          <!-- Filter Section -->
          <form
            class="filter-box"
            @submit.prevent
            @reset="handleResetFilter"
          >
            <div class="filter-date-row">
              <div class="filter-field">
                <label for="filter-from">Dari tanggal</label>
                <input
                  id="filter-from"
                  type="date"
                  v-model="filterFrom"
                  aria-label="Tanggal awal riwayat"
                />
              </div>
              <div class="filter-field">
                <label for="filter-to">Sampai tanggal</label>
                <input
                  id="filter-to"
                  type="date"
                  v-model="filterTo"
                  aria-label="Tanggal akhir riwayat"
                />
              </div>
              <div class="filter-actions">
                <button
                  type="reset"
                  class="reset-filter-btn"
                >
                  Reset filter
                </button>
              </div>
            </div>

            <!-- Poli Radios -->
            <div
              class="filter-poli-row"
              role="radiogroup"
              aria-label="Filter poli riwayat"
            >
              <label
                v-for="poli in availablePolis"
                :key="poli"
                class="poli-radio-label"
              >
                <input
                  type="radio"
                  name="history-poli"
                  :value="poli"
                  v-model="filterPoli"
                />
                <span class="radio-text">{{ poli }}</span>
              </label>
            </div>
          </form>

          <!-- Empty Filter Notice -->
          <p
            v-if="filteredVisits.length === 0"
            class="filter-empty-notice"
            role="status"
          >
            {{
              filterFrom && filterTo && filterFrom > filterTo
                ? "Tanggal akhir harus sama atau setelah tanggal awal."
                : "Tidak ada kunjungan sesuai filter. Coba ubah tanggal atau poli."
            }}
          </p>

          <!-- Visits List -->
          <div v-if="filteredVisits.length" class="visits-list">
            <article
              v-for="visit in filteredVisits"
              :key="visit.id"
              class="visit-card"
            >
              <!-- Visit Header -->
              <div class="visit-header">
                <h3 class="visit-title">
                  {{ formatVisitDate(visit.date) }} · {{ visit.poli }}
                </h3>
                <div class="visit-doctor">{{ getVisitDoctor(visit) }}</div>
                <div class="visit-diagnosis">
                  {{ getVisitDiagnosis(visit) }}
                </div>
              </div>

              <!-- Vitals Grid -->
              <div class="vitals-grid">
                <div
                  v-for="vital in getVitalsList(visit)"
                  :key="vital.label"
                  class="vital-box"
                >
                  <span class="vital-label">{{ vital.label }}</span>
                  <span
                    class="vital-value"
                    :class="{ 'vital-abnormal': vital.abnormal }"
                  >
                    {{ vital.value || "—" }}
                    <template v-if="vital.value"> {{ vital.unit }}</template>
                    <template v-if="vital.abnormal"> · Di luar acuan</template>
                  </span>
                </div>
              </div>

              <!-- Acuan Disclaimer Note -->
              <p class="vitals-footnote">
                Penanda memakai acuan umum dewasa; interpretasikan sesuai kondisi pasien.
              </p>

              <!-- Divider -->
              <hr class="section-divider" />

              <!-- SOAP Notes -->
              <div class="soap-section">
                <div class="soap-item">
                  <span class="soap-label">S · Subjective</span>
                  <p class="soap-text">{{ getVisitSubjective(visit) }}</p>
                </div>
                <hr class="soap-divider" />
                <div class="soap-item">
                  <span class="soap-label">O · Objective</span>
                  <p class="soap-text">{{ getVisitObjective(visit) }}</p>
                </div>
                <hr class="soap-divider" />
                <div class="soap-item">
                  <span class="soap-label">A · Assessment / ICD-10</span>
                  <p class="soap-text">{{ getVisitDiagnosis(visit) }}</p>
                </div>
                <hr class="soap-divider" />
                <div class="soap-item">
                  <span class="soap-label">P · Plan</span>
                  <p class="soap-text">{{ getVisitPlan(visit) }}</p>
                </div>
              </div>
            </article>
          </div>
        </section>
      </div>

      <!-- Empty State if no patient is selected -->
      <div v-else class="empty-state-card">
        <Stethoscope :size="48" class="empty-icon" />
        <h3>Pilih pasien terlebih dahulu</h3>
        <p>Buka dashboard staf dan pilih pasien untuk melihat rekam medis.</p>
        <button type="button" class="button-primary" @click="goToStaff">
          Ke Dashboard Staf
        </button>
      </div>
    </div>
  </main>
</template>

<script setup lang="ts">
import { ref, computed, watch } from "vue";
import { useRouter } from "vue-router";
import { ChevronLeft, Stethoscope } from "lucide-vue-next";
import { useQueueStore } from "../../stores/queue";
import type { Patient } from "../../types";
import { api } from "../../services/api";

const router = useRouter();
const queueStore = useQueueStore();

const filterFrom = ref("");
const filterTo = ref("");
const filterPoli = ref("Semua Poli");
const availablePolis = ["Semua Poli", "Poli Gigi", "Poli Umum"];

const chosenPatient = computed(() => {
  return (
    queueStore.chosen ||
    queueStore.patients.find((p) => p.name === "Budi Santoso") ||
    queueStore.patients.find((p) => p.id === "demo-1") ||
    queueStore.patients[0] ||
    null
  );
});

const allergies = computed(() => {
  if (!chosenPatient.value) return [];
  const p = chosenPatient.value;
  return [
    ...new Set(
      queueStore.patients
        .filter((item) => item.id === p.id || (p.nik && item.nik === p.nik))
        .map((item) => item.allergy?.trim())
        .filter((a) => a && !/^(tidak ada|tidak|nihil|none|-)$/i.test(a)),
    ),
  ];
});

const backendVisits = ref<Patient[]>([]);

watch(
  () => chosenPatient.value,
  async (p) => {
    if (!p || !p.id) {
      backendVisits.value = [];
      return;
    }
    if (
      /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i.test(
        p.id,
      )
    ) {
      try {
        const visits = await api.getPatientVisits(p.id);
        if (visits && visits.length > 0) {
          backendVisits.value = visits.map((v, index) => ({
            ...p,
            id: v.id || `backend-visit-${index}`,
            date: v.visit_date
              ? v.visit_date.slice(0, 10)
              : new Date().toISOString().slice(0, 10),
            poli: v.polyclinic_name || "Poli Umum",
            soap: JSON.stringify({
              doctor: v.doctor_name || "dr. Rama Palito",
              diagnosis:
                v.diagnosis ||
                v.diagnoses
                  ?.map((d) => `${d.icd10_code} — ${d.description || ""}`)
                  .join(", ") ||
                "Pemeriksaan Rawat Jalan",
              subjective: v.symptoms || "Pemeriksaan rutin",
              objective:
                v.physical_exam || "Pemeriksaan fisik dalam batas normal.",
              plan: v.treatment_plan || "Terapi sesuai indikasi.",
              bp_sys: String(v.vitals?.bp_sys || "120"),
              bp_dia: String(v.vitals?.bp_dia || "80"),
              hr: String(v.vitals?.hr || v.vitals?.heart_rate || "80"),
              temp: String(v.vitals?.temp || v.vitals?.temperature || "36.5"),
              spo2: String(v.vitals?.spo2 || "98"),
              weight: String(v.vitals?.weight || "65"),
              height: String(v.vitals?.height || "168"),
            }),
            source: "Rekam Medis Cloud",
            status: "Selesai",
          }));
        }
      } catch {
        // fallback
      }
    }
  },
  { immediate: true },
);

const simulatedHistory = computed<Patient[]>(() => {
  if (!chosenPatient.value) return [];
  const p = chosenPatient.value;
  const examples = [
    {
      days: 7,
      poli: "Poli Umum",
      doctor: "dr. Maya Putri",
      diagnosis: "J00 — Acute nasopharyngitis",
      subjective:
        "Pilek dan nyeri tenggorokan selama 2 hari, disertai demam ringan.",
      objective:
        "Kesadaran baik, faring hiperemis ringan. Suara napas vesikuler.",
      plan: "Edukasi istirahat, cukup cairan, serta kontrol bila keluhan memburuk.",
      bp_sys: "120",
      bp_dia: "80",
      hr: "84",
      temp: "38.2",
      spo2: "98",
      weight: "65",
      height: "168",
    },
    {
      days: 35,
      poli: "Poli Gigi",
      doctor: "drg. Rani Prameswari",
      diagnosis: "K02.9 — Karies gigi",
      subjective: "Gigi geraham kanan bawah terasa ngilu saat minum dingin.",
      objective: "Kavitas pada gigi 46, tanpa pembengkakan gingiva.",
      plan: "Edukasi kebersihan gigi dan jadwal perawatan restorasi.",
      bp_sys: "118",
      bp_dia: "78",
      hr: "76",
      temp: "36.6",
      spo2: "99",
      weight: "65",
      height: "168",
    },
    {
      days: 90,
      poli: "Poli Umum",
      doctor: "dr. Aditya Rahman",
      diagnosis: "I10 — Hipertensi esensial",
      subjective:
        "Kontrol tekanan darah, sesekali pusing. Tidak ada nyeri dada.",
      objective: "Kesadaran baik, tekanan darah meningkat pada pemeriksaan.",
      plan: "Evaluasi terapi, edukasi pembatasan garam, dan pemantauan tekanan darah.",
      bp_sys: "160",
      bp_dia: "100",
      hr: "88",
      temp: "36.7",
      spo2: "98",
      weight: "67",
      height: "168",
    },
  ].map((record, index) => {
    const d = new Date();
    d.setDate(d.getDate() - record.days);
    const localDate = `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, "0")}-${String(d.getDate()).padStart(2, "0")}`;
    return {
      ...p,
      id: `simulation-${p.id}-${index}`,
      date: localDate,
      poli: record.poli,
      soap: JSON.stringify(record),
      source: "Simulasi",
      status: "Selesai",
    };
  });

  const realVisits = queueStore.patients.filter(
    (item) =>
      item.status === "Selesai" &&
      (item.id === p.id || (p.nik && item.nik === p.nik)),
  );

  return [...backendVisits.value, ...realVisits, ...examples].sort((a, b) =>
    b.date.localeCompare(a.date),
  );
});

const filteredVisits = computed(() => {
  return simulatedHistory.value.filter((visit) => {
    if (filterFrom.value && visit.date < filterFrom.value) return false;
    if (filterTo.value && visit.date > filterTo.value) return false;
    if (filterPoli.value !== "Semua Poli" && visit.poli !== filterPoli.value)
      return false;
    return true;
  });
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
    const parts = birthStr.split("-").map(Number);
    const birth =
      parts.length === 3
        ? new Date(parts[0], parts[1] - 1, parts[2])
        : new Date(birthStr);
    const today = new Date();
    let age = today.getFullYear() - birth.getFullYear();
    const m = today.getMonth() - birth.getMonth();
    if (m < 0 || (m === 0 && today.getDate() < birth.getDate())) {
      age--;
    }
    return `${birth.getDate()}/${birth.getMonth() + 1}/${birth.getFullYear()} · ${Math.max(0, age)} tahun`;
  } catch {
    return "Tanggal lahir belum tersedia";
  }
}

function formatVisitDate(dateStr: string) {
  try {
    const parts = dateStr.split("-").map(Number);
    if (parts.length === 3) {
      const d = new Date(parts[0], parts[1] - 1, parts[2]);
      return d.toLocaleDateString("id-ID", {
        day: "numeric",
        month: "long",
        year: "numeric",
      });
    }
    return new Date(`${dateStr}T12:00:00`).toLocaleDateString("id-ID", {
      day: "numeric",
      month: "long",
      year: "numeric",
    });
  } catch {
    return dateStr;
  }
}

function parseSoap(visit: Patient): Record<string, string> {
  try {
    return JSON.parse(visit.soap || "{}") || {};
  } catch {
    return {};
  }
}

function getVisitDoctor(visit: Patient) {
  const rec = parseSoap(visit);
  return rec.doctor || "Dokter belum tercatat";
}

function getVisitDiagnosis(visit: Patient) {
  const rec = parseSoap(visit);
  return rec.diagnosis || "Diagnosis belum tercatat";
}

function getVisitSubjective(visit: Patient) {
  const rec = parseSoap(visit);
  return (
    rec.subjective ||
    [visit.complaint, visit.symptoms.join(", ")].filter(Boolean).join("\n") ||
    "Belum tercatat"
  );
}

function getVisitObjective(visit: Patient) {
  const rec = parseSoap(visit);
  return rec.objective || "Belum tercatat";
}

function getVisitPlan(visit: Patient) {
  const rec = parseSoap(visit);
  return rec.plan || "Belum tercatat";
}

function getVitalsList(visit: Patient) {
  const rec = parseSoap(visit);
  const pressure = String(rec.pressure || "").split("/");
  const sys = rec.bp_sys || pressure[0] || "";
  const dia = rec.bp_dia || pressure[1] || "";
  const hr = rec.hr || rec.pulse || "";
  const temp = rec.temp || rec.temperature || "";
  const spo2 = rec.spo2 || "";

  return [
    {
      label: "Tekanan darah",
      value: sys && dia ? `${sys}/${dia}` : sys || dia || "—",
      unit: "mmHg",
      abnormal:
        Number(sys) >= 140 ||
        Number(dia) >= 90 ||
        (Boolean(sys) && Number(sys) < 90),
    },
    {
      label: "Detak jantung",
      value: hr,
      unit: "/menit",
      abnormal: Boolean(hr) && (Number(hr) > 100 || Number(hr) < 60),
    },
    {
      label: "Suhu",
      value: temp,
      unit: "°C",
      abnormal: Boolean(temp) && (Number(temp) >= 38 || Number(temp) < 36),
    },
    {
      label: "Saturasi oksigen",
      value: spo2,
      unit: "%",
      abnormal: Boolean(spo2) && Number(spo2) < 95,
    },
    {
      label: "Berat badan",
      value: rec.weight,
      unit: "kg",
      abnormal: false,
    },
    {
      label: "Tinggi badan",
      value: rec.height,
      unit: "cm",
      abnormal: false,
    },
  ];
}

function handleResetFilter() {
  filterFrom.value = "";
  filterTo.value = "";
  filterPoli.value = "Semua Poli";
}

function goToStaff() {
  router.push("/staff");
  window.scrollTo(0, 0);
}
</script>

<style scoped>
.medical-record-page {
  min-height: 100vh;
  background-color: #f8fafc;
  padding: 32px 16px 48px;
  display: flex;
  justify-content: center;
  box-sizing: border-box;
  font-family: inherit;
  color: #1e293b;
}

.medical-record-wrapper {
  width: 100%;
  max-width: 860px;
  background-color: #ffffff;
  border: 1px solid #eef2f6;
  border-radius: 24px;
  padding: 28px 32px 36px;
  box-shadow: 0 4px 20px -2px rgba(0, 0, 0, 0.04);
  box-sizing: border-box;
}

/* Top Back Button */
.record-top-nav {
  margin-bottom: 24px;
}

.back-button {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 7px 16px;
  border-radius: 9999px;
  border: 1px solid #d1d5db;
  background-color: #ffffff;
  color: #374151;
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.15s ease;
  line-height: 1;
}

.back-button:hover {
  background-color: #f9fafb;
  border-color: #9ca3af;
}

/* Patient Header */
.patient-header {
  display: flex;
  align-items: flex-start;
  gap: 16px;
  margin-bottom: 24px;
}

.patient-avatar {
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

.patient-details {
  display: flex;
  flex-direction: column;
}

.patient-name {
  font-size: 18px;
  font-weight: 700;
  color: #111827;
  margin: 0 0 4px 0;
  line-height: 1.25;
}

.patient-meta-row {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  color: #64748b;
  margin-bottom: 3px;
}

.meta-dot {
  color: #94a3b8;
}

.gender-badge {
  background-color: #ecfdf5;
  color: #15803d;
  border: 1px solid #bbf7d0;
  font-size: 11px;
  font-weight: 600;
  padding: 1px 8px;
  border-radius: 6px;
  line-height: 1.4;
}

.patient-meta-sub {
  font-size: 13px;
  color: #64748b;
  line-height: 1.4;
}

/* Allergy Alert */
.allergy-alert {
  background-color: #fef2f2;
  border: 1px solid #fee2e2;
  color: #b91c1c;
  border-radius: 10px;
  padding: 10px 16px;
  font-size: 12px;
  margin-bottom: 20px;
}

/* History Card */
.history-card {
  background-color: #ffffff;
  border: 1px solid #e5e7eb;
  border-radius: 18px;
  padding: 24px 28px 32px;
  box-sizing: border-box;
}

.history-title {
  font-size: 18px;
  font-weight: 700;
  color: #111827;
  margin: 0 0 20px 0;
}

/* Filter Box */
.filter-box {
  border: 1px solid #e5e7eb;
  border-radius: 12px;
  padding: 18px 20px;
  background-color: #ffffff;
  margin-bottom: 24px;
}

.filter-date-row {
  display: grid;
  grid-template-columns: 1fr 1fr auto;
  gap: 16px;
  align-items: flex-end;
}

.filter-field {
  display: flex;
  flex-direction: column;
}

.filter-field label {
  font-size: 12px;
  font-weight: 600;
  color: #374151;
  margin-bottom: 6px;
}

.filter-field input[type="date"] {
  width: 100%;
  border: 1px solid #d1d5db;
  border-radius: 8px;
  padding: 7px 12px;
  font-size: 13px;
  color: #1f2937;
  height: 38px;
  outline: none;
  background-color: #ffffff;
  box-sizing: border-box;
  font-family: inherit;
  transition: border-color 0.15s ease, box-shadow 0.15s ease;
}

.filter-field input[type="date"]:focus {
  border-color: #0f766e;
  box-shadow: 0 0 0 2px rgba(15, 118, 110, 0.15);
}

.filter-actions {
  display: flex;
}

.reset-filter-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  height: 38px;
  padding: 0 20px;
  border: 1px solid #d1d5db;
  border-radius: 9999px;
  background-color: #ffffff;
  color: #374151;
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  white-space: nowrap;
  transition: all 0.15s ease;
  font-family: inherit;
}

.reset-filter-btn:hover {
  background-color: #f9fafb;
  border-color: #9ca3af;
}

/* Poli Radios */
.filter-poli-row {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 24px;
  margin-top: 16px;
}

.poli-radio-label {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  font-weight: 500;
  color: #374151;
  cursor: pointer;
}

.poli-radio-label input[type="radio"] {
  accent-color: #0f766e;
  width: 16px;
  height: 16px;
  cursor: pointer;
  margin: 0;
}

.radio-text {
  user-select: none;
}

/* Empty filter notice */
.filter-empty-notice {
  font-size: 13px;
  color: #64748b;
  padding: 16px 0;
  text-align: center;
  margin: 0;
}

/* Visits List */
.visits-list {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.visit-card {
  background-color: #ffffff;
  border: 1px solid #e5e7eb;
  border-left: 4px solid #0f766e;
  border-radius: 12px;
  padding: 22px 24px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.02);
  box-sizing: border-box;
}

/* Visit Header */
.visit-title {
  font-size: 15px;
  font-weight: 700;
  color: #111827;
  margin: 0 0 3px 0;
}

.visit-doctor {
  font-size: 13px;
  color: #64748b;
  margin: 0 0 8px 0;
}

.visit-diagnosis {
  font-size: 15px;
  font-weight: 700;
  color: #0f766e;
  margin: 0 0 18px 0;
}

/* Vitals Grid */
.vitals-grid {
  display: grid;
  grid-template-columns: repeat(5, minmax(0, 1fr));
  gap: 10px;
  margin-bottom: 8px;
}

.vital-box {
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  padding: 10px 12px;
  background-color: #ffffff;
  display: flex;
  flex-direction: column;
  justify-content: center;
  min-height: 60px;
  box-sizing: border-box;
}

.vital-label {
  font-size: 11px;
  color: #64748b;
  margin-bottom: 3px;
  display: block;
}

.vital-value {
  font-size: 13px;
  font-weight: 700;
  color: #111827;
  display: block;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.vital-abnormal {
  color: #dc2626 !important;
}

.vitals-footnote {
  font-size: 11px;
  color: #94a3b8;
  margin: 8px 0 16px 0;
}

/* Dividers */
.section-divider {
  border: 0;
  border-top: 1px solid #f1f5f9;
  margin: 16px 0;
}

.soap-divider {
  border: 0;
  border-top: 1px solid #f1f5f9;
  margin: 12px 0;
}

/* SOAP Section */
.soap-section {
  display: flex;
  flex-direction: column;
}

.soap-item {
  padding: 2px 0;
}

.soap-label {
  font-size: 13px;
  font-weight: 700;
  color: #0f766e;
  display: block;
  margin-bottom: 4px;
}

.soap-text {
  font-size: 13px;
  color: #374151;
  line-height: 1.5;
  margin: 0;
}

/* Empty State Card */
.empty-state-card {
  text-align: center;
  padding: 60px 20px;
  color: #64748b;
}

.empty-icon {
  margin: 0 auto 16px auto;
  color: #0f766e;
  opacity: 0.7;
}

.empty-state-card h3 {
  font-size: 18px;
  font-weight: 700;
  color: #111827;
  margin: 0 0 8px 0;
}

.empty-state-card p {
  font-size: 14px;
  color: #64748b;
  margin: 0 0 20px 0;
}

.button-primary {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 9px 20px;
  background-color: #0f766e;
  color: #ffffff;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: background-color 0.15s ease;
}

.button-primary:hover {
  background-color: #115e59;
}

/* Responsive adjustments */
@media (max-width: 768px) {
  .medical-record-wrapper {
    padding: 20px 16px 24px;
    border-radius: 16px;
  }

  .history-card {
    padding: 18px 16px 24px;
    border-radius: 14px;
  }

  .filter-date-row {
    grid-template-columns: 1fr;
  }

  .reset-filter-btn {
    width: 100%;
  }

  .vitals-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}
</style>
