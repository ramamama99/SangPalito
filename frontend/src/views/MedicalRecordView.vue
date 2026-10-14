<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ChevronLeft } from 'lucide-vue-next'
import { useQueueStore } from '../stores/queue'

const router = useRouter()
const queueStore = useQueueStore()

const fromDate = ref('')
const toDate = ref('')
const selectedPoli = ref('Semua Poli')

const resetFilter = () => {
  fromDate.value = ''
  toDate.value = ''
  selectedPoli.value = 'Semua Poli'
}

interface VisitRecord {
  id: string
  date: string
  poli: string
  doctor: string
  diagnosis: string
  subjective: string
  objective: string
  plan: string
  bp_sys: string
  bp_dia: string
  hr: string
  temp: string
  spo2: string
  weight: string
  height: string
}

const defaultVisits: VisitRecord[] = [
  {
    id: 'visit-1',
    date: '2026-10-08',
    poli: 'Poli Umum',
    doctor: 'dr. Maya Putri',
    diagnosis: 'J00 — Acute nasopharyngitis',
    subjective: 'Pilek dan nyeri tenggorokan selama 2 hari, disertai demam ringan.',
    objective: 'Kesadaran baik, faring hiperemis ringan. Suara napas vesikuler.',
    plan: 'Edukasi istirahat, cukup cairan, serta kontrol bila keluhan memburuk.',
    bp_sys: '120',
    bp_dia: '80',
    hr: '84',
    temp: '38.2',
    spo2: '98',
    weight: '65',
    height: '168',
  },
  {
    id: 'visit-2',
    date: '2026-09-10',
    poli: 'Poli Gigi',
    doctor: 'drg. Rani Prameswari',
    diagnosis: 'K02.9 — Karies gigi',
    subjective: 'Gigi geraham kanan bawah terasa ngilu saat minum dingin.',
    objective: 'Kavitas pada gigi 46, tanpa pembengkakan gingiva.',
    plan: 'Edukasi kebersihan gigi dan jadwal perawatan restorasi.',
    bp_sys: '118',
    bp_dia: '78',
    hr: '76',
    temp: '36.6',
    spo2: '99',
    weight: '65',
    height: '168',
  },
  {
    id: 'visit-3',
    date: '2026-07-17',
    poli: 'Poli Umum',
    doctor: 'dr. Aditya Rahman',
    diagnosis: 'I10 — Hipertensi esensial',
    subjective: 'Kontrol tekanan darah, sesekali pusing. Tidak ada nyeri dada.',
    objective: 'Kesadaran baik, tekanan darah meningkat pada pemeriksaan.',
    plan: 'Evaluasi terapi, edukasi pembatasan garam, dan pemantauan tekanan darah.',
    bp_sys: '160',
    bp_dia: '100',
    hr: '88',
    temp: '36.7',
    spo2: '98',
    weight: '67',
    height: '168',
  },
]

// Keep track of which visit cards are opened
const openVisits = ref<Record<string, boolean>>({})

const toggleVisit = (id: string) => {
  openVisits.value[id] = !openVisits.value[id]
}

const isOpen = (id: string) => {
  return !!openVisits.value[id]
}

const formatDate = (dateStr: string) => {
  const [y, m, d] = dateStr.split('-')
  const months = [
    'Januari',
    'Februari',
    'Maret',
    'April',
    'Mei',
    'Juni',
    'Juli',
    'Agustus',
    'September',
    'Oktober',
    'November',
    'Desember',
  ]
  const day = parseInt(d, 10)
  const month = months[parseInt(m, 10) - 1]
  return `${day} ${month} ${y}`
}

const isBpAbnormal = (sys?: string, dia?: string) => {
  const s = Number(sys)
  const d = Number(dia)
  if (!s || !d) return false
  return s >= 140 || s < 90 || d >= 90 || d < 60
}

const isHrAbnormal = (hr?: string) => {
  const h = Number(hr)
  if (!h) return false
  return h > 100 || h < 60
}

const isTempAbnormal = (temp?: string) => {
  const t = Number(temp)
  if (!t) return false
  return t >= 38 || t < 36
}

const isSpo2Abnormal = (spo2?: string) => {
  const s = Number(spo2)
  if (!s) return false
  return s < 95
}

const allVisits = computed<VisitRecord[]>(() => {
  const realVisits: VisitRecord[] = []
  const patient = queueStore.chosenPatient || queueStore.patients[0]

  if (patient) {
    const finished = queueStore.patients.filter(
      (p) =>
        p.status === 'Selesai' &&
        (p.id === patient.id || (patient.nik && p.nik === patient.nik)),
    )
    for (const f of finished) {
      if (f.soap) {
        try {
          const parsed = JSON.parse(f.soap)
          realVisits.push({
            id: f.id,
            date: f.date || '2026-10-15',
            poli: f.poli,
            doctor: parsed.doctor || (f.poli === 'Poli Gigi' ? 'drg. Rani Prameswari' : 'dr. Maya Putri'),
            diagnosis: parsed.diagnosis || f.complaint,
            subjective: parsed.subjective || f.complaint,
            objective: parsed.objective || 'Pemeriksaan fisik dalam batas normal.',
            plan: parsed.plan || 'Terapi sesuai indikasi medis.',
            bp_sys: parsed.bp_sys || '120',
            bp_dia: parsed.bp_dia || '80',
            hr: parsed.hr || '80',
            temp: parsed.temp || '36.5',
            spo2: parsed.spo2 || '98',
            weight: parsed.weight || '65',
            height: parsed.height || '168',
          })
        } catch {
          // parse fallback
        }
      }
    }
  }

  // Combine and sort by date descending
  return [...realVisits, ...defaultVisits].sort((a, b) => b.date.localeCompare(a.date))
})

const filteredVisits = computed(() => {
  return allVisits.value.filter((v) => {
    if (fromDate.value && v.date < fromDate.value) return false
    if (toDate.value && v.date > toDate.value) return false
    if (selectedPoli.value !== 'Semua Poli' && v.poli !== selectedPoli.value) return false
    return true
  })
})
</script>

<template>
  <div class="medical-record-wrapper min-h-screen bg-[#f3f6f5] py-6 px-4 sm:px-6">
    <div class="max-w-4xl mx-auto">
      <!-- Back button -->
      <div class="mb-3">
        <button
          type="button"
          class="back-btn inline-flex items-center gap-1 text-xs text-gray-500 hover:text-[#12b8b0] transition-colors"
          @click="router.push('/staff')"
        >
          <ChevronLeft :size="15" />
          <span>Kembali ke Pemeriksaan</span>
        </button>
      </div>

      <!-- Main card container matching the screenshots -->
      <div class="mr-container bg-[#f8faf9] rounded-2xl border border-gray-200/80 p-5 sm:p-7 shadow-xs">
        <h2 class="mr-title text-[20px] font-bold text-[#1f2937] mb-4">
          Riwayat rekam medis
        </h2>

        <!-- Filter Card -->
        <div class="filter-card bg-white rounded-xl border border-gray-200 p-4 sm:p-5 mb-6 shadow-xs">
          <div class="filter-dates-row flex flex-wrap items-end gap-3 sm:gap-4">
            <!-- Dari tanggal -->
            <div class="filter-field flex-1 min-w-[180px]">
              <label class="block text-xs font-semibold text-gray-700 mb-1">
                Dari tanggal
              </label>
              <input
                type="date"
                v-model="fromDate"
                class="date-input w-full border border-gray-300 rounded-lg px-3 py-2 text-xs sm:text-sm text-gray-700 bg-white"
              />
            </div>

            <!-- Sampai tanggal -->
            <div class="filter-field flex-1 min-w-[180px]">
              <label class="block text-xs font-semibold text-gray-700 mb-1">
                Sampai tanggal
              </label>
              <input
                type="date"
                v-model="toDate"
                class="date-input w-full border border-gray-300 rounded-lg px-3 py-2 text-xs sm:text-sm text-gray-700 bg-white"
              />
            </div>

            <!-- Reset filter button -->
            <div class="reset-btn-wrap">
              <button
                type="button"
                @click="resetFilter"
                class="reset-btn border border-gray-700 text-gray-800 text-xs font-medium rounded-full px-5 py-2 hover:bg-gray-50 transition-colors"
              >
                Reset filter
              </button>
            </div>
          </div>

          <!-- Radio Filter Pills -->
          <div class="radio-pills flex flex-wrap items-center gap-2.5 mt-3.5">
            <label
              v-for="poli in ['Semua Poli', 'Poli Gigi', 'Poli Umum']"
              :key="poli"
              class="radio-pill border border-gray-200 rounded-full px-3 py-1.5 inline-flex items-center gap-2 text-xs font-medium cursor-pointer transition-colors bg-white hover:border-gray-300"
              :class="{ 'radio-pill-active': selectedPoli === poli }"
            >
              <input
                type="radio"
                :value="poli"
                v-model="selectedPoli"
                class="accent-gray-800 w-3.5 h-3.5"
              />
              <span>{{ poli }}</span>
            </label>
          </div>
        </div>

        <!-- Visits List with Left Teal Line -->
        <div v-if="filteredVisits.length" class="timeline-container border-l-2 border-[#12b8b0] pl-4 sm:pl-5 space-y-3.5">
          <div
            v-for="visit in filteredVisits"
            :key="visit.id"
            class="visit-card bg-white rounded-xl border border-gray-200 p-4 sm:p-5 shadow-xs transition-shadow"
          >
            <!-- Card Header / Toggle -->
            <div
              class="visit-header cursor-pointer select-none"
              @click="toggleVisit(visit.id)"
            >
              <!-- Arrow Triangle -->
              <div class="arrow-wrap mb-1.5">
                <svg
                  v-if="!isOpen(visit.id)"
                  viewBox="0 0 10 10"
                  class="triangle-arrow w-2.5 h-2.5 fill-gray-800"
                >
                  <polygon points="0,0 10,5 0,10" />
                </svg>
                <svg
                  v-else
                  viewBox="0 0 10 10"
                  class="triangle-arrow w-2.5 h-2.5 fill-gray-800"
                >
                  <polygon points="0,0 10,0 5,10" />
                </svg>
              </div>

              <!-- Date & Poli -->
              <div class="visit-title text-[15px] font-bold text-gray-900 leading-snug">
                {{ formatDate(visit.date) }} · {{ visit.poli }}
              </div>

              <!-- Doctor -->
              <div class="visit-doctor text-xs text-gray-500 mt-0.5">
                {{ visit.doctor }}
              </div>

              <!-- Diagnosis -->
              <div class="visit-diagnosis text-[13.5px] font-bold text-[#0f766e] mt-1.5">
                {{ visit.diagnosis }}
              </div>
            </div>

            <!-- Card Body (Shown when open) -->
            <div v-if="isOpen(visit.id)" class="visit-body mt-4">
              <!-- Vitals Grid -->
              <div class="vitals-grid grid grid-cols-2 sm:grid-cols-3 md:grid-cols-5 gap-2.5 sm:gap-3">
                <!-- Tekanan Darah -->
                <div class="vital-item border border-gray-200 rounded-lg p-3 bg-white">
                  <div class="vital-label text-xs text-gray-500 font-normal">Tekanan darah</div>
                  <div
                    class="vital-val text-[13.5px] sm:text-[14px] font-bold mt-1"
                    :class="isBpAbnormal(visit.bp_sys, visit.bp_dia) ? 'text-[#991b1b]' : 'text-gray-900'"
                  >
                    {{ visit.bp_sys }}/{{ visit.bp_dia }} mmHg
                    <template v-if="isBpAbnormal(visit.bp_sys, visit.bp_dia)">
                      · Di luar acuan
                    </template>
                  </div>
                </div>

                <!-- Detak Jantung -->
                <div class="vital-item border border-gray-200 rounded-lg p-3 bg-white">
                  <div class="vital-label text-xs text-gray-500 font-normal">Detak jantung</div>
                  <div
                    class="vital-val text-[13.5px] sm:text-[14px] font-bold mt-1"
                    :class="isHrAbnormal(visit.hr) ? 'text-[#991b1b]' : 'text-gray-900'"
                  >
                    {{ visit.hr }} /menit
                    <template v-if="isHrAbnormal(visit.hr)">
                      · Di luar acuan
                    </template>
                  </div>
                </div>

                <!-- Suhu -->
                <div class="vital-item border border-gray-200 rounded-lg p-3 bg-white">
                  <div class="vital-label text-xs text-gray-500 font-normal">Suhu</div>
                  <div
                    class="vital-val text-[13.5px] sm:text-[14px] font-bold mt-1"
                    :class="isTempAbnormal(visit.temp) ? 'text-[#991b1b]' : 'text-gray-900'"
                  >
                    {{ visit.temp }} °C
                    <template v-if="isTempAbnormal(visit.temp)">
                      · Di luar acuan
                    </template>
                  </div>
                </div>

                <!-- Saturasi Oksigen -->
                <div class="vital-item border border-gray-200 rounded-lg p-3 bg-white">
                  <div class="vital-label text-xs text-gray-500 font-normal">Saturasi oksigen</div>
                  <div
                    class="vital-val text-[13.5px] sm:text-[14px] font-bold mt-1"
                    :class="isSpo2Abnormal(visit.spo2) ? 'text-[#991b1b]' : 'text-gray-900'"
                  >
                    {{ visit.spo2 }} %
                    <template v-if="isSpo2Abnormal(visit.spo2)">
                      · Di luar acuan
                    </template>
                  </div>
                </div>

                <!-- Berat Badan -->
                <div class="vital-item border border-gray-200 rounded-lg p-3 bg-white">
                  <div class="vital-label text-xs text-gray-500 font-normal">Berat badan</div>
                  <div class="vital-val text-[13.5px] sm:text-[14px] font-bold text-gray-900 mt-1">
                    {{ visit.weight }} kg
                  </div>
                </div>

                <!-- Tinggi Badan (Row 2) -->
                <div class="vital-item border border-gray-200 rounded-lg p-3 bg-white">
                  <div class="vital-label text-xs text-gray-500 font-normal">Tinggi badan</div>
                  <div class="vital-val text-[13.5px] sm:text-[14px] font-bold text-gray-900 mt-1">
                    {{ visit.height }} cm
                  </div>
                </div>
              </div>

              <!-- Note below vitals -->
              <p class="vitals-note text-xs text-gray-500 mt-3 mb-3.5">
                Penanda memakai acuan umum dewasa; interpretasikan sesuai kondisi pasien.
              </p>

              <!-- SOAP Sections with divider lines -->
              <div class="soap-section border-t border-gray-200 pt-3 pb-3">
                <div class="soap-label text-xs font-bold text-[#0f766e]">S · Subjective</div>
                <p class="soap-content text-[13.5px] text-gray-700 mt-1 leading-relaxed">
                  {{ visit.subjective }}
                </p>
              </div>

              <div class="soap-section border-t border-gray-200 pt-3 pb-3">
                <div class="soap-label text-xs font-bold text-[#0f766e]">O · Objective</div>
                <p class="soap-content text-[13.5px] text-gray-700 mt-1 leading-relaxed">
                  {{ visit.objective }}
                </p>
              </div>

              <div class="soap-section border-t border-gray-200 pt-3 pb-3">
                <div class="soap-label text-xs font-bold text-[#0f766e]">A · Assessment / ICD-10</div>
                <p class="soap-content text-[13.5px] text-gray-700 mt-1 leading-relaxed">
                  {{ visit.diagnosis }}
                </p>
              </div>

              <div class="soap-section border-t border-gray-200 pt-3 pb-1">
                <div class="soap-label text-xs font-bold text-[#0f766e]">P · Plan</div>
                <p class="soap-content text-[13.5px] text-gray-700 mt-1 leading-relaxed">
                  {{ visit.plan }}
                </p>
              </div>
            </div>
          </div>
        </div>

        <!-- Filter Empty State -->
        <p v-else class="text-sm text-gray-500 py-8 text-center bg-white rounded-xl border border-gray-200">
          Tidak ada kunjungan sesuai filter. Coba ubah tanggal atau poli.
        </p>
      </div>
    </div>
  </div>
</template>

<style scoped>
.date-input {
  min-height: 38px;
  border: 1px solid #d1d5db;
  border-radius: 8px;
  padding: 6px 12px;
  font-family: inherit;
  outline: none;
}
.date-input:focus {
  border-color: #0f766e;
  box-shadow: 0 0 0 1px #0f766e;
}
.reset-btn {
  height: 38px;
  border: 1px solid #4b5563;
  color: #1f2937;
  display: inline-flex;
  align-items: center;
  justify-content: center;
}
.radio-pill {
  user-select: none;
}
.radio-pill:hover {
  background-color: #f9fafb;
}
.timeline-container {
  border-left-width: 2px;
  border-left-color: #12b8b0;
}
.triangle-arrow {
  display: inline-block;
  vertical-align: middle;
}
</style>
