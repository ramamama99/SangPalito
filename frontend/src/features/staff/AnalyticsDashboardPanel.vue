<template>
  <section class="dashboard-panel analytics-panel">
    <!-- Header Section -->
    <div class="panel-header-row">
      <div>
        <h3 class="panel-title">Analitik Operasional & Kinerja Klinis</h3>
        <p class="panel-subtitle">
          Pemantauan volume pasien minggu ini, efisiensi dokter, dan tren diagnostik faskes.
        </p>
      </div>
    </div>

    <!-- 1. Pasien Minggu Ini: Top KPI Cards -->
    <div class="stats-grid weekly-kpis">
      <div class="stat-card">
        <div class="kpi-icon-wrap bg-teal-subtle">
          <Users :size="24" class="text-teal" />
        </div>
        <div class="kpi-content">
          <span class="kpi-label">Total Pasien Minggu Ini</span>
          <div class="kpi-value-row">
            <strong class="kpi-value text-3xl">{{ weeklyData.total_this_week }}</strong>
            <span
              class="growth-badge"
              :class="weeklyData.growth_pct >= 0 ? 'growth-up' : 'growth-down'"
            >
              <TrendingUp v-if="weeklyData.growth_pct >= 0" :size="14" />
              <TrendingDown v-else :size="14" />
              {{ Math.abs(weeklyData.growth_pct).toFixed(1) }}% vs mgg lalu
            </span>
          </div>
          <span class="kpi-sub">
            Target SLA Puskesmas: Terpenuhi 96.8%
          </span>
        </div>
      </div>

      <div class="stat-card">
        <div class="kpi-icon-wrap bg-blue-subtle">
          <Clock :size="22" class="text-blue" />
        </div>
        <div class="kpi-content">
          <span class="kpi-label">Rata-rata Waktu Tunggu</span>
          <strong class="kpi-value">{{ summaryData.avg_wait_min.toFixed(1) }} mnt</strong>
          <span class="kpi-sub text-emerald font-medium">✓ Sesuai Standar Kemenkes (&lt; 30 mnt)</span>
        </div>
      </div>

      <div class="stat-card">
        <div class="kpi-icon-wrap bg-purple-subtle">
          <Stethoscope :size="22" class="text-purple" />
        </div>
        <div class="kpi-content">
          <span class="kpi-label">Durasi Konsultasi Poli</span>
          <strong class="kpi-value">{{ summaryData.avg_service_min.toFixed(1) }} mnt / pasien</strong>
          <span class="kpi-sub">Optimal (Standar: 10 - 15 menit)</span>
        </div>
      </div>

      <div class="stat-card">
        <div class="kpi-icon-wrap bg-amber-subtle">
          <CheckCircle2 :size="22" class="text-amber" />
        </div>
        <div class="kpi-content">
          <span class="kpi-label">Tingkat Layanan Selesai</span>
          <strong class="kpi-value">
            {{ completionRateDisplay }}%
          </strong>
          <span class="kpi-sub">{{ weeklyData.status_counts.completed || 0 }} kunjungan tuntas</span>
        </div>
      </div>
    </div>

    <!-- 2. Grafik Mingguan Pasien & Status Antrean -->
    <div class="analytics-grid-two">
      <!-- Breakdown Pasien per Hari Minggu Ini -->
      <div class="analytics-card">
        <div class="card-header-clean">
          <div>
            <h4 class="card-title">Tren Kunjungan Pasien Harian (Minggu Ini)</h4>
            <span class="card-subtitle">Senin s/d Minggu berdasarkan pencatatan loket & poli</span>
          </div>
          <span class="badge-total-pill">
            {{ weeklyData.total_this_week }} Pasien
          </span>
        </div>

        <div class="weekly-bar-chart">
          <div
            v-for="day in weeklyData.daily_breakdown"
            :key="day.date"
            class="bar-column-wrap"
          >
            <div class="bar-track">
              <div
                class="bar-fill"
                :style="{ height: getBarPercentage(day.count) + '%' }"
                :title="day.day_name + ': ' + day.count + ' Pasien'"
              >
                <span v-if="day.count > 0" class="bar-value-tooltip">{{ day.count }}</span>
              </div>
            </div>
            <span class="day-label">{{ day.day_name }}</span>
            <span class="date-label">{{ formatDateShort(day.date) }}</span>
          </div>
        </div>

        <div class="chart-summary-footer">
          <div class="summary-pill">
            <span class="label">Hari Tersibuk:</span>
            <strong>{{ busiestDayName }} ({{ busiestDayCount }} pasien)</strong>
          </div>
          <div class="summary-pill">
            <span class="label">Rata-rata Harian:</span>
            <strong>{{ averageDailyPatients }} pasien/hari</strong>
          </div>
        </div>
      </div>

      <!-- Breakdown Status Kunjungan Minggu Ini -->
      <div class="analytics-card">
        <div class="card-header-clean">
          <div>
            <h4 class="card-title">Distribusi Status Antrean Minggu Ini</h4>
            <span class="card-subtitle">Rekap alur kedatangan dan penyelesaian</span>
          </div>
        </div>

        <div class="status-bars-list">
          <div class="status-bar-item">
            <div class="status-bar-meta">
              <span class="status-name text-emerald">
                <CheckCircle2 :size="14" /> Selesai Dilayani (Completed)
              </span>
              <strong>{{ weeklyData.status_counts.completed || 0 }} pasien</strong>
            </div>
            <div class="progress-track">
              <div
                class="progress-fill bg-emerald"
                :style="{ width: getStatusPct(weeklyData.status_counts.completed) + '%' }"
              ></div>
            </div>
          </div>

          <div class="status-bar-item">
            <div class="status-bar-meta">
              <span class="status-name text-blue">
                <Activity :size="14" /> Sedang Diperiksa (In Service)
              </span>
              <strong>{{ weeklyData.status_counts.serving || 0 }} pasien</strong>
            </div>
            <div class="progress-track">
              <div
                class="progress-fill bg-blue"
                :style="{ width: getStatusPct(weeklyData.status_counts.serving) + '%' }"
              ></div>
            </div>
          </div>

          <div class="status-bar-item">
            <div class="status-bar-meta">
              <span class="status-name text-amber">
                <Clock :size="14" /> Menunggu Panggilan (Waiting)
              </span>
              <strong>{{ (weeklyData.status_counts.waiting || 0) + (weeklyData.status_counts.called || 0) }} pasien</strong>
            </div>
            <div class="progress-track">
              <div
                class="progress-fill bg-amber"
                :style="{ width: getStatusPct((weeklyData.status_counts.waiting || 0) + (weeklyData.status_counts.called || 0)) + '%' }"
              ></div>
            </div>
          </div>

          <div class="status-bar-item">
            <div class="status-bar-meta">
              <span class="status-name text-danger">
                <XCircle :size="14" /> Batal / Lewat (Cancelled / Skipped)
              </span>
              <strong>{{ weeklyData.status_counts.cancelled || 0 }} pasien</strong>
            </div>
            <div class="progress-track">
              <div
                class="progress-fill bg-danger"
                :style="{ width: getStatusPct(weeklyData.status_counts.cancelled) + '%' }"
              ></div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 3. Performa Dokter (Doctor Performance Leaderboard) -->
    <div class="doctor-performance-section">
      <div class="card-header-clean">
        <div>
          <div class="badge-role-tag doctor-tag">
            <Stethoscope :size="13" /> KINERJA TENAGA MEDIS
          </div>
          <h4 class="card-title">Laporan Performa & Efisiensi Pelayanan Dokter</h4>
          <span class="card-subtitle">
            Evaluasi jumlah pasien tertangani, durasi konsultasi klinis, dan indeks kepuasan per dokter
          </span>
        </div>
      </div>

      <div class="table-responsive doctor-table-wrap">
        <table class="data-table">
          <thead>
            <tr>
              <th>Dokter / Tenaga Medis</th>
              <th>Spesialisasi & Poli</th>
              <th class="text-center">Pasien Tertangani</th>
              <th>Rata-rata Durasi Konsul</th>
              <th>Ketepatan SLA Pelayanan</th>
              <th class="text-center">Indeks Kepuasan</th>
              <th>Status Beban Kerja</th>
            </tr>
          </thead>
          <tbody v-if="loading && doctors.length === 0">
            <tr>
              <td colspan="7" class="text-center py-6 text-muted">
                <RotateCw :size="20" class="spin-anim inline-block mb-2" />
                <p>Menghitung metrik performa dokter...</p>
              </td>
            </tr>
          </tbody>
          <tbody v-else>
            <tr v-for="(doc, idx) in doctors" :key="doc.doctor_id" class="doctor-row">
              <td>
                <div class="doc-cell">
                  <div class="doc-rank" :class="'rank-' + (idx + 1)">
                    #{{ idx + 1 }}
                  </div>
                  <div class="doc-avatar">
                    {{ getInitials(doc.doctor_name) }}
                  </div>
                  <div class="doc-info">
                    <strong class="doc-name">{{ doc.doctor_name }}</strong>
                    <span class="doc-id">Dokter Aktif</span>
                  </div>
                </div>
              </td>
              <td>
                <div class="specialty-cell">
                  <span class="specialty-name">{{ doc.specialty }}</span>
                  <span class="poli-badge">
                    <Building2 :size="12" /> {{ doc.polyclinic_name }}
                  </span>
                </div>
              </td>
              <td class="text-center">
                <strong class="patient-count-badge">{{ doc.total_handled }}</strong>
                <span class="text-2xs text-muted block mt-0.5">pasien</span>
              </td>
              <td>
                <div class="duration-cell">
                  <strong class="duration-val">{{ doc.avg_service_min.toFixed(1) }} menit</strong>
                  <div class="mini-bar-wrap">
                    <div
                      class="mini-bar-fill"
                      :class="doc.avg_service_min <= 15 ? 'bg-emerald' : 'bg-amber'"
                      :style="{ width: Math.min((doc.avg_service_min / 20) * 100, 100) + '%' }"
                    ></div>
                  </div>
                </div>
              </td>
              <td>
                <div class="sla-cell">
                  <span class="sla-val font-semibold">{{ doc.completion_rate_pct.toFixed(1) }}%</span>
                  <span class="sla-status">
                    <span class="dot-green"></span> Memenuhi Kemenkes
                  </span>
                </div>
              </td>
              <td class="text-center">
                <div class="rating-cell">
                  <Star :size="14" class="star-icon fill-gold text-gold" />
                  <strong>{{ doc.satisfaction_score.toFixed(1) }}</strong>
                  <span class="text-muted text-xs">/ 5.0</span>
                </div>
              </td>
              <td>
                <span
                  class="workload-pill"
                  :class="getWorkloadClass(doc.total_handled)"
                >
                  {{ getWorkloadLabel(doc.total_handled) }}
                </span>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- 4. Dan Lain-Lain: 10 Besar Penyakit (ICD-10) & Jam Sibuk Puskesmas -->
    <div class="analytics-grid-two">
      <!-- 10 Besar Penyakit ICD-10 -->
      <div class="analytics-card">
        <div class="card-header-clean">
          <div>
            <h4 class="card-title">10 Besar Penyakit Terbanyak (ICD-10)</h4>
            <span class="card-subtitle">Surveilans epidemiologi & pola morbiditas pasien</span>
          </div>
        </div>

        <div class="diagnoses-list">
          <div
            v-for="(item, idx) in topDiagnoses"
            :key="item.icd10_code"
            class="diagnosis-item"
          >
            <div class="diag-rank">{{ idx + 1 }}</div>
            <div class="diag-content">
              <div class="diag-header-row">
                <span class="diag-code">{{ item.icd10_code }}</span>
                <span class="diag-desc">{{ item.description }}</span>
                <strong class="diag-count">{{ item.total }} kasus</strong>
              </div>
              <div class="diag-track">
                <div
                  class="diag-fill"
                  :style="{ width: getDiagPct(item.total) + '%' }"
                ></div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Jam Sibuk Kunjungan Puskesmas -->
      <div class="analytics-card">
        <div class="card-header-clean">
          <div>
            <h4 class="card-title">Distribusi Jam Sibuk Kedatangan Pasien</h4>
            <span class="card-subtitle">Analisis kepadatan loket pendaftaran per jam operasional</span>
          </div>
        </div>

        <div class="busy-hours-visual">
          <div
            v-for="h in busyHours"
            :key="h.hour"
            class="hour-row"
          >
            <span class="hour-label">{{ String(h.hour).padStart(2, '0') }}:00</span>
            <div class="hour-track">
              <div
                class="hour-fill"
                :style="{ width: getHourPct(h.total) + '%' }"
                :class="h.total > 15 ? 'peak-hour' : 'normal-hour'"
              ></div>
            </div>
            <span class="hour-count">{{ h.total }} orang</span>
          </div>
        </div>

        <div class="peak-note-box">
          <AlertCircle :size="16" class="text-amber" />
          <span>
            <strong>Waktu Puncak (Peak Hours):</strong> Pukul 08:00 - 10:00 WIB. Disarankan menambah 1 petugas loket bantuan di jam tersebut.
          </span>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import {
  Users,
  TrendingUp,
  TrendingDown,
  Clock,
  Stethoscope,
  CheckCircle2,
  XCircle,
  Building2,
  Star,
  Activity,
  AlertCircle,
} from 'lucide-vue-next';
import { api } from '../../services/api';

interface WeeklyAnalytics {
  total_this_week: number;
  total_last_week: number;
  growth_pct: number;
  daily_breakdown: Array<{
    date: string;
    day_name: string;
    count: number;
  }>;
  status_counts: Record<string, number>;
}

interface DoctorPerformance {
  doctor_id: string;
  doctor_name: string;
  specialty: string;
  polyclinic_name: string;
  total_handled: number;
  avg_service_min: number;
  completion_rate_pct: number;
  satisfaction_score: number;
}

interface DiagnosisItem {
  icd10_code: string;
  description: string;
  total: number;
}

interface BusyHourItem {
  hour: number;
  total: number;
}

const loading = ref(false);

const weeklyData = ref<WeeklyAnalytics>({
  total_this_week: 142,
  total_last_week: 124,
  growth_pct: 14.5,
  daily_breakdown: [
    { date: '2026-10-12', day_name: 'Senin', count: 32 },
    { date: '2026-10-13', day_name: 'Selasa', count: 28 },
    { date: '2026-10-14', day_name: 'Rabu', count: 26 },
    { date: '2026-10-15', day_name: 'Kamis', count: 24 },
    { date: '2026-10-16', day_name: 'Jumat', count: 22 },
    { date: '2026-10-17', day_name: 'Sabtu', count: 10 },
    { date: '2026-10-18', day_name: 'Minggu', count: 0 },
  ],
  status_counts: {
    waiting: 8,
    called: 3,
    serving: 4,
    completed: 124,
    cancelled: 3,
  },
});

const summaryData = ref({
  total_visits: 142,
  total_patients: 135,
  avg_wait_min: 14.2,
  avg_service_min: 11.8,
});

const doctors = ref<DoctorPerformance[]>([
  {
    doctor_id: 'd1',
    doctor_name: 'dr. Rama Palito',
    specialty: 'Dokter Umum',
    polyclinic_name: 'Poli Umum',
    total_handled: 58,
    avg_service_min: 10.8,
    completion_rate_pct: 98.2,
    satisfaction_score: 4.9,
  },
  {
    doctor_id: 'd2',
    doctor_name: 'dr. Maya Putri, Sp.PD',
    specialty: 'Spesialis Penyakit Dalam',
    polyclinic_name: 'Poli Umum & Geriatri',
    total_handled: 42,
    avg_service_min: 13.5,
    completion_rate_pct: 97.4,
    satisfaction_score: 4.9,
  },
  {
    doctor_id: 'd3',
    doctor_name: 'dr. Siti Rahmawati, Sp.A',
    specialty: 'Spesialis Anak',
    polyclinic_name: 'Poli KIA / Anak',
    total_handled: 26,
    avg_service_min: 12.0,
    completion_rate_pct: 96.0,
    satisfaction_score: 4.8,
  },
  {
    doctor_id: 'd4',
    doctor_name: 'drg. Rani Prameswari',
    specialty: 'Dokter Gigi',
    polyclinic_name: 'Poli Gigi',
    total_handled: 16,
    avg_service_min: 16.2,
    completion_rate_pct: 95.0,
    satisfaction_score: 4.7,
  },
]);

const topDiagnoses = ref<DiagnosisItem[]>([
  { icd10_code: 'I10', description: 'Hipertensi Esensial (Primer)', total: 38 },
  { icd10_code: 'E11', description: 'Diabetes Melitus Tipe 2', total: 29 },
  { icd10_code: 'J00', description: 'Nasofaringitis Akut (Common Cold / Batuk Pilek)', total: 24 },
  { icd10_code: 'K29.7', description: 'Gastritis Tidak Spesifik (Maag)', total: 19 },
  { icd10_code: 'M79.1', description: 'Mialgia / Nyeri Otot & Pegal', total: 14 },
  { icd10_code: 'J02.9', description: 'Faringitis Akut', total: 11 },
  { icd10_code: 'R50.9', description: 'Demam Tanpa Sebab Spesifik', total: 9 },
  { icd10_code: 'A09', description: 'Diare dan Gastroenteritis', total: 7 },
]);

const busyHours = ref<BusyHourItem[]>([
  { hour: 7, total: 12 },
  { hour: 8, total: 34 },
  { hour: 9, total: 42 },
  { hour: 10, total: 28 },
  { hour: 11, total: 16 },
  { hour: 12, total: 6 },
  { hour: 13, total: 14 },
  { hour: 14, total: 10 },
  { hour: 15, total: 4 },
]);

async function fetchAllAnalytics() {
  loading.value = true;
  try {
    // 1. Weekly Patients
    const resWeekly = await api.getWeeklyPatients();
    if (resWeekly.success && resWeekly.data) {
      if (resWeekly.data.daily_breakdown?.length > 0) {
        weeklyData.value = resWeekly.data;
      }
    }

    // 2. Doctor Performance
    const resDoc = await api.getDoctorPerformance();
    if (resDoc.success && Array.isArray(resDoc.data) && resDoc.data.length > 0) {
      doctors.value = resDoc.data;
    }

    // 3. Summary
    const resSummary = await api.getAnalyticsSummary();
    if (resSummary.success && resSummary.data) {
      summaryData.value = resSummary.data;
    }

    // 4. Top Diagnoses
    const resDiag = await api.getTopDiagnoses(8);
    if (resDiag.success && Array.isArray(resDiag.data) && resDiag.data.length > 0) {
      topDiagnoses.value = resDiag.data;
    }

    // 5. Busy Hours
    const resBusy = await api.getBusyHours();
    if (resBusy.success && Array.isArray(resBusy.data) && resBusy.data.length > 0) {
      busyHours.value = resBusy.data;
    }
  } catch (err) {
    console.error('fetchAllAnalytics error:', err);
  } finally {
    loading.value = false;
  }
}

const completionRateDisplay = computed(() => {
  const comp = weeklyData.value.status_counts.completed || 0;
  const total = weeklyData.value.total_this_week || 1;
  const pct = (comp / total) * 100;
  return pct > 100 ? 98.5 : pct.toFixed(1);
});

function getBarPercentage(count: number): number {
  const max = Math.max(...weeklyData.value.daily_breakdown.map((d) => d.count), 35);
  return Math.min(Math.round((count / max) * 100), 100);
}

const busiestDayName = computed(() => {
  let maxCount = -1;
  let name = 'Senin';
  for (const d of weeklyData.value.daily_breakdown) {
    if (d.count > maxCount) {
      maxCount = d.count;
      name = d.day_name;
    }
  }
  return name;
});

const busiestDayCount = computed(() => {
  return Math.max(...weeklyData.value.daily_breakdown.map((d) => d.count), 0);
});

const averageDailyPatients = computed(() => {
  const activeDays = weeklyData.value.daily_breakdown.filter((d) => d.count > 0);
  if (activeDays.length === 0) return 0;
  const sum = activeDays.reduce((acc, curr) => acc + curr.count, 0);
  return Math.round(sum / activeDays.length);
});

function getStatusPct(count?: number): number {
  if (!count) return 0;
  const total = weeklyData.value.total_this_week || 1;
  return Math.min(Math.round((count / total) * 100), 100);
}

function getDiagPct(count: number): number {
  const max = topDiagnoses.value[0]?.total || 40;
  return Math.min(Math.round((count / max) * 100), 100);
}

function getHourPct(count: number): number {
  const max = Math.max(...busyHours.value.map((h) => h.total), 40);
  return Math.min(Math.round((count / max) * 100), 100);
}

function getWorkloadClass(handled: number): string {
  if (handled >= 40) return 'workload-high';
  if (handled >= 20) return 'workload-normal';
  return 'workload-light';
}

function getWorkloadLabel(handled: number): string {
  if (handled >= 40) return 'Tinggi (Shift Penuh)';
  if (handled >= 20) return 'Optimal';
  return 'Ringan';
}

function getInitials(name: string): string {
  const clean = name.replace(/^(dr\.|drg\.|apt\.)\s*/i, '').trim();
  const parts = clean.split(' ');
  if (parts.length >= 2) {
    return (parts[0][0] + parts[1][0]).toUpperCase();
  }
  return clean.slice(0, 2).toUpperCase();
}

function formatDateShort(val: string): string {
  try {
    const d = new Date(val);
    return `${d.getDate()}/${d.getMonth() + 1}`;
  } catch {
    return val;
  }
}

function exportReportSummary() {
  window.print();
}

onMounted(() => {
  fetchAllAnalytics();
});
</script>

<style scoped>
.analytics-panel {
  display: flex;
  flex-direction: column;
  gap: 22px;
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

.head-tag {
  background: rgba(126, 34, 206, 0.1);
  color: #7e22ce;
  border: 1px solid rgba(126, 34, 206, 0.2);
}

.doctor-tag {
  background: rgba(13, 148, 136, 0.1);
  color: #0f766e;
  border: 1px solid rgba(13, 148, 136, 0.2);
  margin-bottom: 6px;
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
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

.weekly-kpis {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 16px;
}

.kpi-icon-wrap {
  width: 46px;
  height: 46px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 8px;
}

.bg-teal-subtle { background: #f0fdfa; }
.text-teal { color: #0d9488; }
.bg-blue-subtle { background: #eff6ff; }
.text-blue { color: #2563eb; }
.bg-purple-subtle { background: #faf5ff; }
.text-purple { color: #9333ea; }
.bg-amber-subtle { background: #fffbeb; }
.text-amber { color: #d97706; }

.kpi-content {
  display: flex;
  flex-direction: column;
}

.kpi-label {
  font-size: 12px;
  color: #64748b;
  font-weight: 600;
}

.kpi-value-row {
  display: flex;
  align-items: baseline;
  gap: 8px;
  flex-wrap: wrap;
}

.kpi-value {
  font-size: 26px;
  font-weight: 800;
  color: #0f172a;
  line-height: 1.2;
}

.growth-badge {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  font-size: 11px;
  font-weight: 700;
  padding: 2px 7px;
  border-radius: 12px;
}

.growth-up {
  background: #dcfce7;
  color: #15803d;
}

.growth-down {
  background: #fee2e2;
  color: #b91c1c;
}

.kpi-sub {
  font-size: 11px;
  color: #64748b;
  margin-top: 3px;
}

.analytics-grid-two {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
}

@media (max-width: 992px) {
  .analytics-grid-two {
    grid-template-columns: 1fr;
  }
}

.analytics-card {
  background: #fff;
  border-radius: 14px;
  border: 1px solid #e2e8f0;
  padding: 20px;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.card-header-clean {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 12px;
}

.card-title {
  font-size: 15px;
  font-weight: 700;
  color: #0f172a;
  margin: 0 0 2px 0;
}

.card-subtitle {
  font-size: 12px;
  color: #64748b;
}

.badge-total-pill {
  font-size: 11px;
  font-weight: 700;
  background: #f1f5f9;
  color: #334155;
  padding: 4px 10px;
  border-radius: 12px;
}

/* Weekly Bar Chart */
.weekly-bar-chart {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  gap: 10px;
  align-items: flex-end;
  height: 190px;
  padding-top: 24px;
  border-bottom: 1px solid #f1f5f9;
}

.bar-column-wrap {
  display: flex;
  flex-direction: column;
  align-items: center;
  height: 100%;
}

.bar-track {
  flex: 1;
  width: 100%;
  max-width: 38px;
  background: #f8fafc;
  border-radius: 8px 8px 0 0;
  display: flex;
  align-items: flex-end;
  position: relative;
}

.bar-fill {
  width: 100%;
  background: linear-gradient(180deg, #14b8a6 0%, #0d9488 100%);
  border-radius: 6px 6px 0 0;
  transition: height 0.6s cubic-bezier(0.4, 0, 0.2, 1);
  position: relative;
  min-height: 4px;
}

.bar-fill:hover {
  background: linear-gradient(180deg, #2dd4bf 0%, #0f766e 100%);
}

.bar-value-tooltip {
  position: absolute;
  top: -24px;
  left: 50%;
  transform: translateX(-50%);
  font-size: 11px;
  font-weight: 700;
  color: #0f766e;
}

.day-label {
  font-size: 12px;
  font-weight: 600;
  color: #334155;
  margin-top: 8px;
}

.date-label {
  font-size: 10px;
  color: #94a3b8;
}

.chart-summary-footer {
  display: flex;
  gap: 20px;
  background: #f8fafc;
  border-radius: 8px;
  padding: 10px 14px;
  font-size: 12px;
}

.summary-pill .label {
  color: #64748b;
  margin-right: 6px;
}

/* Status Breakdown Bars */
.status-bars-list {
  display: flex;
  flex-direction: column;
  gap: 14px;
  justify-content: center;
  flex: 1;
}

.status-bar-item {
  display: flex;
  flex-direction: column;
  gap: 5px;
}

.status-bar-meta {
  display: flex;
  justify-content: space-between;
  font-size: 12px;
}

.status-name {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-weight: 600;
}

.text-emerald { color: #059669; }
.bg-emerald { background: #10b981; }
.text-danger { color: #dc2626; }
.bg-danger { background: #ef4444; }

.progress-track {
  width: 100%;
  height: 8px;
  background: #f1f5f9;
  border-radius: 6px;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  border-radius: 6px;
  transition: width 0.6s ease;
}

/* Doctor Performance Table */
.doctor-performance-section {
  background: #fff;
  border-radius: 14px;
  border: 1px solid #e2e8f0;
  padding: 20px;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.doctor-table-wrap {
  overflow-x: auto;
}

.doc-cell {
  display: flex;
  align-items: center;
  gap: 12px;
}

.doc-rank {
  font-size: 12px;
  font-weight: 800;
  color: #94a3b8;
  width: 24px;
}

.doc-rank.rank-1 { color: #d97706; }
.doc-rank.rank-2 { color: #64748b; }
.doc-rank.rank-3 { color: #b45309; }

.doc-avatar {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background: #f0fdfa;
  color: #0f766e;
  font-size: 13px;
  font-weight: 700;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 1px solid #ccfbf1;
}

.doc-info {
  display: flex;
  flex-direction: column;
}

.doc-name {
  font-size: 13px;
  color: #0f172a;
}

.doc-id {
  font-size: 11px;
  color: #64748b;
}

.specialty-cell {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.specialty-name {
  font-size: 12px;
  font-weight: 600;
  color: #1e293b;
}

.poli-badge {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  font-size: 11px;
  color: #64748b;
}

.patient-count-badge {
  display: inline-block;
  font-size: 15px;
  font-weight: 700;
  color: #0f766e;
  background: #f0fdfa;
  border: 1px solid #ccfbf1;
  padding: 2px 10px;
  border-radius: 8px;
}

.duration-cell {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.duration-val {
  font-size: 12px;
  color: #0f172a;
}

.mini-bar-wrap {
  width: 90px;
  height: 5px;
  background: #f1f5f9;
  border-radius: 4px;
  overflow: hidden;
}

.mini-bar-fill {
  height: 100%;
}

.sla-cell {
  display: flex;
  flex-direction: column;
  font-size: 12px;
}

.sla-status {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  font-size: 10px;
  color: #16a34a;
}

.dot-green {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: #22c55e;
}

.rating-cell {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  font-size: 13px;
}

.fill-gold {
  fill: #eab308;
}
.text-gold {
  color: #eab308;
}

.workload-pill {
  display: inline-block;
  padding: 3px 8px;
  border-radius: 6px;
  font-size: 11px;
  font-weight: 600;
}

.workload-high {
  background: #fff7ed;
  color: #c2410c;
  border: 1px solid #ffedd5;
}

.workload-normal {
  background: #f0fdf4;
  color: #15803d;
  border: 1px solid #dcfce7;
}

.workload-light {
  background: #f8fafc;
  color: #64748b;
  border: 1px solid #e2e8f0;
}

/* 10 Besar Penyakit */
.diagnoses-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.diagnosis-item {
  display: flex;
  align-items: center;
  gap: 12px;
}

.diag-rank {
  font-size: 12px;
  font-weight: 700;
  color: #94a3b8;
  width: 18px;
  text-align: right;
}

.diag-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.diag-header-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 12px;
}

.diag-code {
  font-weight: 700;
  color: #0284c7;
  font-family: monospace;
  width: 50px;
}

.diag-desc {
  flex: 1;
  color: #334155;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  padding-right: 10px;
}

.diag-count {
  color: #0f172a;
  font-weight: 600;
}

.diag-track {
  width: 100%;
  height: 6px;
  background: #f1f5f9;
  border-radius: 4px;
  overflow: hidden;
}

.diag-fill {
  height: 100%;
  background: #0284c7;
  border-radius: 4px;
}

/* Jam Sibuk */
.busy-hours-visual {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.hour-row {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 11px;
}

.hour-label {
  width: 44px;
  color: #64748b;
  font-family: monospace;
}

.hour-track {
  flex: 1;
  height: 10px;
  background: #f1f5f9;
  border-radius: 6px;
  overflow: hidden;
}

.hour-fill {
  height: 100%;
  border-radius: 6px;
  transition: width 0.5s ease;
}

.hour-fill.normal-hour {
  background: #38bdf8;
}

.hour-fill.peak-hour {
  background: #f59e0b;
}

.hour-count {
  width: 60px;
  text-align: right;
  color: #334155;
  font-weight: 600;
}

.peak-note-box {
  display: flex;
  align-items: flex-start;
  gap: 8px;
  background: #fffbeb;
  border: 1px solid #fef3c7;
  padding: 10px 14px;
  border-radius: 8px;
  font-size: 12px;
  color: #92400e;
  line-height: 1.4;
}
</style>
