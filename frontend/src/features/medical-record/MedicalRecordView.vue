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

              <!-- Dedicated Section for Medication & e-Prescription -->
              <div v-if="getVisitMedicationPlan(visit) || getVisitPrescription(visit)" class="rx-history-card">
                <div class="rx-history-header">
                  <div class="rx-history-title">
                    <Pill :size="16" class="rx-pill-icon" />
                    <strong>Perencanaan Terapi Obat & e-Resep</strong>
                  </div>
                  <div v-if="getVisitPrescription(visit)?.prescription_no" class="rx-status-meta">
                    <span class="rx-no-badge">{{ getVisitPrescription(visit)?.prescription_no }}</span>
                    <span
                      class="rx-status-badge"
                      :class="getVisitPrescription(visit)?.status === 'dispensed' ? 'dispensed' : 'pending'"
                    >
                      {{ getVisitPrescription(visit)?.status === 'dispensed' ? 'Sudah Ditebus' : 'Menunggu Penebusan' }}
                    </span>
                  </div>
                </div>

                <div v-if="getVisitPrescription(visit)?.items?.length" class="rx-meds-list">
                  <div
                    v-for="(item, itIdx) in getVisitPrescription(visit)!.items"
                    :key="itIdx"
                    class="rx-med-item"
                  >
                    <div class="rx-med-top">
                      <div class="rx-med-name-row">
                        <span class="route-badge" :data-route="item.route">
                          {{ formatRouteLabel(item.route) }}
                        </span>
                        <strong class="rx-med-name">{{ item.medication_name }}</strong>
                      </div>
                      <span class="rx-med-qty">{{ item.quantity }} {{ item.unit }}</span>
                    </div>
                    <div class="rx-med-details">
                      <p><strong>Aturan:</strong> {{ item.frequency }} · {{ item.timing }}</p>
                      <p v-if="item.instruction" class="rx-med-instruction">
                        <strong>Petunjuk Cara Pakai:</strong> {{ item.instruction }}
                      </p>
                    </div>
                  </div>
                </div>
                <div v-else class="rx-plain-plan">
                  <p>{{ getVisitMedicationPlan(visit) }}</p>
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
import { ChevronLeft, Stethoscope, Pill } from "lucide-vue-next";
import { useQueueStore } from "../../stores/queue";
import type { Patient, Prescription } from "../../types";
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
            medication_plan: v.medication_plan,
            prescription_no: v.prescription_no,
            prescription: v.prescription,
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
              medication_plan: v.medication_plan,
              prescription_no: v.prescription_no,
              prescription: v.prescription,
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
      medication_plan: "1. Paracetamol 500mg (Diminum - Oral) - 3 x 1 tablet sehari sesudah makan.\n2. Antasida DOEN (Dikunyah - Chewable) - 3 x 1 tablet kunyah 1 jam sebelum makan.\n3. Vitamin C 500mg (Dikunyah - Chewable) - 1 x 1 tablet kunyah sesudah makan.",
      prescription_no: "RX-20261008-U012",
      prescription: {
        id: "rx-sim-1",
        prescription_no: "RX-20261008-U012",
        status: "dispensed",
        doctor_name: "dr. Maya Putri",
        polyclinic: "Poli Umum",
        items: [
          {
            medication_name: "Paracetamol 500mg",
            dosage_form: "Tablet",
            route: "oral_minum",
            frequency: "3 x 1 tablet sehari",
            timing: "Sesudah makan",
            quantity: 10,
            unit: "Tablet",
            instruction: "Diminum utuh dengan 1 gelas air putih sesudah makan bila demam atau nyeri.",
          },
          {
            medication_name: "Antasida DOEN Kunyah",
            dosage_form: "Tablet Kunyah",
            route: "oral_kunyah",
            frequency: "3 x 1 tablet kunyah sehari",
            timing: "1 jam sebelum makan atau saat maag kambuh",
            quantity: 12,
            unit: "Tablet",
            instruction: "Wajib dikunyah hingga lumat di mulut sebelum ditelan, lalu minum sedikit air hangat.",
          },
        ],
      },
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
      plan: "Edukasi kebersihan gigi dan penambalan kavitas.",
      medication_plan: "1. Asam Mefenamat 500mg (Diminum - Oral) - 3 x 1 kapsul sesudah makan jika nyeri.",
      prescription_no: "RX-20260910-G005",
      prescription: {
        id: "rx-sim-2",
        prescription_no: "RX-20260910-G005",
        status: "dispensed",
        doctor_name: "drg. Rani Prameswari",
        polyclinic: "Poli Gigi",
        items: [
          {
            medication_name: "Asam Mefenamat 500mg",
            dosage_form: "Kapsul",
            route: "oral_minum",
            frequency: "3 x 1 kapsul sehari",
            timing: "Sesudah makan",
            quantity: 10,
            unit: "Kapsul",
            instruction: "Diminum segera sesudah makan untuk mencegah iritasi lambung. Hentikan jika nyeri reda.",
          },
        ],
      },
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
        "Kontrol tekanan darah rutin, sesekali pusing tegang di tengkuk. Tidak ada nyeri dada.",
      objective: "Kesadaran baik, tekanan darah meningkat pada pemeriksaan.",
      plan: "Evaluasi terapi, edukasi pembatasan garam, dan pemantauan tekanan darah.",
      medication_plan: "1. Amlodipine 5mg (Diminum - Oral) - 1 x 1 tablet malam sebelum tidur.",
      prescription_no: "RX-20260717-U022",
      prescription: {
        id: "rx-sim-3",
        prescription_no: "RX-20260717-U022",
        status: "dispensed",
        doctor_name: "dr. Aditya Rahman",
        polyclinic: "Poli Umum",
        items: [
          {
            medication_name: "Amlodipine 5mg",
            dosage_form: "Tablet",
            route: "oral_minum",
            frequency: "1 x 1 tablet sehari",
            timing: "Malam hari sebelum tidur",
            quantity: 30,
            unit: "Tablet",
            instruction: "Diminum teratur tiap malam pada jam yang sama.",
          },
        ],
      },
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
      medication_plan: record.medication_plan,
      prescription_no: record.prescription_no,
      prescription: record.prescription,
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

function getVisitPrescription(visit: Patient): Prescription | undefined {
  if (visit.prescription) return visit.prescription;
  const rec = parseSoap(visit) as any;
  return rec.prescription;
}

function getVisitMedicationPlan(visit: Patient): string {
  if (visit.medication_plan) return visit.medication_plan;
  const rec = parseSoap(visit) as any;
  return rec.medication_plan || "";
}

function formatRouteLabel(route?: string) {
  if (!route) return "Diminum";
  const routeLabels: Record<string, string> = {
    oral_minum: "Diminum",
    oral_kunyah: "Dikunyah",
    topikal_oles: "Dioleskan",
    tetes_mata: "Tetes Mata",
    tetes_telinga: "Tetes Telinga",
    kumur: "Dikumur",
    inhalasi: "Dihirup",
    sirup_minum: "Sirup",
  };
  return routeLabels[route] || route;
}

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
  width: 100%;
  background-color: #f6f8f5;
  padding: clamp(14px, 2.5vw, 36px);
  display: flex;
  justify-content: center;
  box-sizing: border-box;
  font-family: inherit;
  color: #243e40;
}

.medical-record-wrapper {
  width: 100%;
  max-width: 100%;
  background-color: #ffffff;
  border: 1px solid var(--line);
  border-radius: 19px;
  padding: clamp(16px, 2.5vw, 36px);
  box-shadow: 0 4px 20px -2px rgba(25, 90, 95, 0.04);
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
  padding: 8px 18px;
  border-radius: 50px;
  border: 1px solid var(--line);
  background-color: #ffffff;
  color: var(--deep);
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
  line-height: 1;
}

.back-button:hover {
  background-color: #f6f8f5;
  border-color: #c8e8df;
  color: var(--deep);
}

/* Patient Header */
.patient-header {
  display: flex;
  align-items: flex-start;
  gap: 16px;
  margin-bottom: 24px;
}

.patient-avatar {
  width: 48px;
  height: 48px;
  border-radius: 14px;
  background-color: #eaf6f2;
  color: var(--deep);
  border: 1px solid #c8e8df;
  font-weight: 700;
  font-size: 15px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  letter-spacing: 0.5px;
  font-family: 'Manrope', sans-serif;
}

.patient-details {
  display: flex;
  flex-direction: column;
}

.patient-name {
  font-family: 'Manrope', sans-serif;
  font-size: 20px;
  font-weight: 700;
  color: var(--deep);
  margin: 0 0 4px 0;
  line-height: 1.25;
}

.patient-meta-row {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  color: var(--muted);
  margin-bottom: 3px;
}

.meta-dot {
  color: var(--muted);
}

.gender-badge {
  background-color: #eaf6f2;
  color: var(--deep);
  border: 1px solid #c8e8df;
  font-size: 11px;
  font-weight: 600;
  padding: 1px 8px;
  border-radius: 6px;
  line-height: 1.4;
}

.patient-meta-sub {
  font-size: 13px;
  color: var(--muted);
  line-height: 1.4;
}

/* Allergy Alert */
.allergy-alert {
  background-color: #fdf3f0;
  border: 1px solid #f7d6cf;
  color: var(--coral);
  border-radius: 12px;
  padding: 12px 18px;
  font-size: 13px;
  margin-bottom: 20px;
}

/* History Card */
.history-card {
  background-color: #ffffff;
  border: 1px solid var(--line);
  border-radius: 19px;
  padding: 24px 28px 32px;
  box-sizing: border-box;
}

.history-title {
  font-family: 'Manrope', sans-serif;
  font-size: 18px;
  font-weight: 700;
  color: var(--deep);
  margin: 0 0 20px 0;
}

/* Filter Box */
.filter-box {
  border: 1px solid var(--line);
  border-radius: 14px;
  padding: 18px 20px;
  background-color: #f6f8f5;
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
  color: var(--deep);
  margin-bottom: 6px;
}

.filter-field input[type="date"] {
  width: 100%;
  border: 1px solid #dce5df;
  border-radius: 9px;
  padding: 8px 12px;
  font-size: 13px;
  color: #243e40;
  height: 38px;
  outline: none;
  background-color: #ffffff;
  box-sizing: border-box;
  font-family: inherit;
  transition: border-color 0.15s ease, box-shadow 0.15s ease;
}

.filter-field input[type="date"]:focus {
  border-color: var(--teal);
  box-shadow: 0 0 0 2px rgba(18, 184, 176, 0.2);
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
  border: 1px solid var(--line);
  border-radius: 50px;
  background-color: #ffffff;
  color: var(--deep);
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  white-space: nowrap;
  transition: all 0.2s ease;
  font-family: inherit;
}

.reset-filter-btn:hover {
  background-color: #eaf6f2;
  border-color: #c8e8df;
}

/* Poli Radios */
.filter-poli-row {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 20px;
  margin-top: 16px;
}

.poli-radio-label {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  font-weight: 500;
  color: #243e40;
  cursor: pointer;
}

.poli-radio-label input[type="radio"] {
  accent-color: var(--teal);
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
  color: var(--muted);
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
  border: 1px solid var(--line);
  border-left: 4px solid var(--teal);
  border-radius: 14px;
  padding: 22px 24px;
  box-shadow: 0 2px 8px -2px rgba(25, 90, 95, 0.04);
  box-sizing: border-box;
}

/* Visit Header */
.visit-title {
  font-family: 'Manrope', sans-serif;
  font-size: 16px;
  font-weight: 700;
  color: var(--deep);
  margin: 0 0 4px 0;
}

.visit-doctor {
  font-size: 13px;
  color: var(--muted);
  margin: 0 0 8px 0;
}

.visit-diagnosis {
  font-family: 'Manrope', sans-serif;
  font-size: 15px;
  font-weight: 700;
  color: var(--teal);
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
  border: 1px solid var(--line);
  border-radius: 9px;
  padding: 10px 12px;
  background-color: #f6f8f5;
  display: flex;
  flex-direction: column;
  justify-content: center;
  min-height: 60px;
  box-sizing: border-box;
}

.vital-label {
  font-size: 11px;
  color: var(--muted);
  margin-bottom: 3px;
  display: block;
}

.vital-value {
  font-family: 'Manrope', sans-serif;
  font-size: 13px;
  font-weight: 700;
  color: var(--deep);
  display: block;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.vital-abnormal {
  color: var(--coral) !important;
}

.vitals-footnote {
  font-size: 11px;
  color: var(--muted);
  margin: 8px 0 16px 0;
}

/* Dividers */
.section-divider {
  border: 0;
  border-top: 1px solid var(--line);
  margin: 16px 0;
}

.soap-divider {
  border: 0;
  border-top: 1px solid var(--line);
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
  font-family: 'Manrope', sans-serif;
  font-size: 13px;
  font-weight: 700;
  color: var(--deep);
  display: block;
  margin-bottom: 4px;
}

.soap-text {
  font-size: 13px;
  color: #243e40;
  line-height: 1.6;
  margin: 0;
}

/* Empty State Card */
.empty-state-card {
  text-align: center;
  padding: 60px 20px;
  color: var(--muted);
}

.empty-icon {
  margin: 0 auto 16px auto;
  color: var(--teal);
  opacity: 0.8;
}

.empty-state-card h3 {
  font-family: 'Manrope', sans-serif;
  font-size: 18px;
  font-weight: 700;
  color: var(--deep);
  margin: 0 0 8px 0;
}

.empty-state-card p {
  font-size: 14px;
  color: var(--muted);
  margin: 0 0 20px 0;
}

.button-primary {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 10px 24px;
  background-color: var(--teal);
  color: #ffffff;
  border: none;
  border-radius: 50px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: background-color 0.2s ease;
  box-shadow: 0 4px 12px rgba(18, 184, 176, 0.25);
}

.button-primary:hover {
  background-color: #0ea59e;
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

/* ── Prescription & Medication History Styles ── */
.rx-history-card {
  margin-top: 16px;
  padding: 16px 18px;
  background: #f6f8f5;
  border: 1px solid var(--line);
  border-radius: 14px;
}

.rx-history-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
  flex-wrap: wrap;
  gap: 8px;
}

.rx-history-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-family: 'Manrope', sans-serif;
  font-size: 13px;
  font-weight: 700;
  color: var(--deep);
}

.rx-pill-icon {
  color: var(--teal);
}

.rx-status-meta {
  display: flex;
  align-items: center;
  gap: 8px;
}

.rx-no-badge {
  font-family: monospace;
  font-weight: 700;
  font-size: 11px;
  padding: 2px 8px;
  background: #eaf6f2;
  border: 1px solid #c8e8df;
  border-radius: 6px;
  color: var(--deep);
}

.rx-status-badge {
  font-size: 11px;
  font-weight: 700;
  padding: 2px 10px;
  border-radius: 50px;
}

.rx-status-badge.dispensed {
  background: #eaf6f2;
  color: var(--deep);
  border: 1px solid #c8e8df;
}

.rx-status-badge.pending {
  background: var(--warm);
  color: var(--deep);
  border: 1px solid #e7ded4;
}

.rx-meds-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.rx-med-item {
  padding: 12px 14px;
  background: #ffffff;
  border: 1px solid var(--line);
  border-radius: 9px;
  font-size: 12px;
}

.rx-med-top {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 6px;
}

.rx-med-name-row {
  display: flex;
  align-items: center;
  gap: 8px;
}

.rx-med-name {
  color: var(--deep);
  font-size: 13px;
  font-weight: 700;
}

.rx-med-qty {
  font-weight: 700;
  font-family: monospace;
  color: var(--deep);
}

.rx-med-details {
  color: #243e40;
  line-height: 1.5;
}

.rx-med-instruction {
  color: var(--teal);
  margin-top: 2px;
  font-weight: 500;
}

.rx-plain-plan {
  font-size: 12px;
  color: #243e40;
  white-space: pre-line;
  background: #ffffff;
  padding: 10px 12px;
  border-radius: 9px;
  border: 1px solid var(--line);
}

/* Route Badges */
.route-badge {
  display: inline-flex;
  align-items: center;
  padding: 2px 8px;
  border-radius: 6px;
  font-size: 10px;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.3px;
}

.route-badge[data-route="oral_minum"],
.route-badge[data-route="sirup_minum"] {
  background: #eaf6f2;
  color: var(--deep);
  border: 1px solid #c8e8df;
}

.route-badge[data-route="oral_kunyah"] {
  background: var(--warm);
  color: var(--deep);
  border: 1px solid #e7ded4;
}

.route-badge[data-route="topikal_oles"] {
  background: #eaf3f6;
  color: var(--deep);
  border: 1px solid #c5dbe3;
}

.route-badge[data-route="tetes_mata"],
.route-badge[data-route="tetes_telinga"] {
  background: #f4ecf8;
  color: #5a3e75;
  border: 1px solid #dfcce6;
}

.route-badge[data-route="kumur"] {
  background: #fdf3f0;
  color: var(--coral);
  border: 1px solid #f7d6cf;
}

.route-badge[data-route="inhalasi"] {
  background: #eaf6f2;
  color: var(--teal);
  border: 1px solid #c8e8df;
}
</style>
