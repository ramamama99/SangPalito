<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ChevronRight, CalendarDays } from 'lucide-vue-next'
import { useQueueStore } from '../stores/queue'

const router = useRouter()
const queueStore = useQueueStore()

const chosen = computed(() => {
  return (
    queueStore.chosenPatient ||
    queueStore.patients.find((p) => p.status === 'Hadir') ||
    queueStore.patients[0]
  )
})

const fromDate = ref('')
const toDate = ref('')
const selectedPoli = ref('Semua Poli')

const resetFilter = () => {
  fromDate.value = ''
  toDate.value = ''
  selectedPoli.value = 'Semua Poli'
}

const patientAllergies = computed(() => {
  if (!chosen.value) return []
  return [
    ...new Set(
      queueStore.patients
        .filter(
          (p) => p.id === chosen.value!.id || (chosen.value!.nik && p.nik === chosen.value!.nik),
        )
        .map((p) => p.allergy?.trim())
        .filter((a) => a && !/^(tidak ada|tidak|nihil|none|-)$/i.test(a)),
    ),
  ]
})

const visits = computed(() => {
  if (!chosen.value) return []

  const examples = [
    {
      days: 7,
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
      days: 35,
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
      days: 90,
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
  ].map((record, index) => {
    const date = new Date()
    date.setDate(date.getDate() - record.days)
    const localDate = `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')}`
    return {
      ...chosen.value!,
      id: `simulation-${chosen.value!.id}-${index}`,
      date: localDate,
      poli: record.poli,
      soap: JSON.stringify(record),
      source: 'Simulasi',
      status: 'Selesai',
    }
  })

  const realVisits = queueStore.patients.filter(
    (p) =>
      p.status === 'Selesai' &&
      (p.id === chosen.value!.id || (chosen.value!.nik && p.nik === chosen.value!.nik)),
  )

  const combined = [...realVisits, ...examples].sort((a, b) => b.date.localeCompare(a.date))

  return combined.filter((v) => {
    if (fromDate.value && v.date < fromDate.value) return false
    if (toDate.value && v.date > toDate.value) return false
    if (selectedPoli.value !== 'Semua Poli' && v.poli !== selectedPoli.value) return false
    return true
  })
})

const getParsedRecord = (soapStr?: string) => {
  try {
    return JSON.parse(soapStr || '{}') || {}
  } catch {
    return {}
  }
}
</script>

<template>
  <main class="dashboard medical-record-page">
    <div class="dashboard-body">
      <section v-if="chosen" class="dashboard-panel soap-panel">
        <div class="record-back">
          <button type="button" class="button outline small" @click="router.push('/staff')">
            <ChevronRight :size="17" class="rotate-180" /> Kembali ke Pemeriksaan
          </button>
        </div>

        <div class="record-banner">
          <div class="record-identity record-identity-compact">
            <span class="record-avatar" aria-hidden="true">
              {{
                chosen.name
                  .split(' ')
                  .slice(0, 2)
                  .map((n) => n[0])
                  .join('')
              }}
            </span>
            <div>
              <h3>{{ chosen.name }}</h3>
              <p>
                No. RM: belum tersedia ·
                <span class="badge">{{ chosen.gender }}</span>
              </p>
              <p>
                {{
                  chosen.birth
                    ? `${new Date(chosen.birth).toLocaleDateString('id-ID')} · ${Math.max(0, new Date().getFullYear() - new Date(chosen.birth).getFullYear())} tahun`
                    : 'Tanggal lahir belum tersedia'
                }}
              </p>
              <p v-if="chosen.insurance === 'BPJS'">BPJS: {{ chosen.bpjs || 'Belum diisi' }}</p>
            </div>
          </div>

          <div v-if="patientAllergies.length" class="record-alert" role="alert">
            PERINGATAN ALERGI: {{ patientAllergies.join(', ') }}
          </div>
        </div>

        <div class="record-layout">
          <div class="record-history record-history-responsive">
            <div class="record-scroll">
              <h3>Riwayat rekam medis</h3>

              <form class="record-filter-form" @submit.prevent>
                <div class="record-date-fields">
                  <label>
                    Dari tanggal
                    <input type="date" v-model="fromDate" aria-label="Tanggal awal riwayat" />
                  </label>
                  <label>
                    Sampai tanggal
                    <input type="date" v-model="toDate" aria-label="Tanggal akhir riwayat" />
                  </label>
                  <button class="button outline small" type="button" @click="resetFilter">
                    Reset filter
                  </button>
                </div>

                <div class="record-filters" role="group" aria-label="Filter poli riwayat">
                  <label v-for="poli in ['Semua Poli', 'Poli Gigi', 'Poli Umum']" :key="poli">
                    <input type="radio" :value="poli" v-model="selectedPoli" />
                    {{ poli }}
                  </label>
                </div>
              </form>

              <p v-if="!visits.length" class="record-filter-empty" role="status" aria-live="polite">
                Tidak ada kunjungan sesuai filter. Coba ubah tanggal atau poli.
              </p>

              <div v-else class="record-timeline">
                <details
                  v-for="visit in visits"
                  :key="visit.id"
                  class="record-visit"
                  :data-poli="visit.poli"
                  :data-date="visit.date"
                  open
                >
                  <summary>
                    <strong>
                      {{
                        new Date(`${visit.date}T12:00:00`).toLocaleDateString('id-ID', {
                          day: 'numeric',
                          month: 'long',
                          year: 'numeric',
                        })
                      }}
                      · {{ visit.poli }}
                    </strong>
                    <small>
                      {{ getParsedRecord(visit.soap).doctor || 'Dokter belum tercatat' }}
                    </small>
                    <strong>
                      {{ getParsedRecord(visit.soap).diagnosis || 'Diagnosis belum tercatat' }}
                    </strong>
                  </summary>

                  <div class="record-vitals">
                    <div>
                      <small>Tekanan darah</small>
                      <strong
                        :class="
                          Number(getParsedRecord(visit.soap).bp_sys) >= 140 ||
                          Number(getParsedRecord(visit.soap).bp_dia) >= 90 ||
                          Number(getParsedRecord(visit.soap).bp_sys) < 90
                            ? 'record-abnormal'
                            : ''
                        "
                      >
                        {{ getParsedRecord(visit.soap).bp_sys || '—' }}/{{
                          getParsedRecord(visit.soap).bp_dia || '—'
                        }}
                        mmHg
                        <template
                          v-if="
                            Number(getParsedRecord(visit.soap).bp_sys) >= 140 ||
                            Number(getParsedRecord(visit.soap).bp_dia) >= 90 ||
                            Number(getParsedRecord(visit.soap).bp_sys) < 90
                          "
                        >
                          · Di luar acuan
                        </template>
                      </strong>
                    </div>

                    <div>
                      <small>Detak jantung</small>
                      <strong
                        :class="
                          Number(getParsedRecord(visit.soap).hr) > 100 ||
                          Number(getParsedRecord(visit.soap).hr) < 60
                            ? 'record-abnormal'
                            : ''
                        "
                      >
                        {{ getParsedRecord(visit.soap).hr || '—' }} /menit
                        <template
                          v-if="
                            Number(getParsedRecord(visit.soap).hr) > 100 ||
                            Number(getParsedRecord(visit.soap).hr) < 60
                          "
                        >
                          · Di luar acuan
                        </template>
                      </strong>
                    </div>

                    <div>
                      <small>Suhu</small>
                      <strong
                        :class="
                          Number(getParsedRecord(visit.soap).temp) >= 38 ||
                          Number(getParsedRecord(visit.soap).temp) < 36
                            ? 'record-abnormal'
                            : ''
                        "
                      >
                        {{ getParsedRecord(visit.soap).temp || '—' }} °C
                        <template
                          v-if="
                            Number(getParsedRecord(visit.soap).temp) >= 38 ||
                            Number(getParsedRecord(visit.soap).temp) < 36
                          "
                        >
                          · Di luar acuan
                        </template>
                      </strong>
                    </div>

                    <div>
                      <small>Saturasi oksigen</small>
                      <strong
                        :class="
                          Number(getParsedRecord(visit.soap).spo2) < 95 ? 'record-abnormal' : ''
                        "
                      >
                        {{ getParsedRecord(visit.soap).spo2 || '—' }} %
                        <template v-if="Number(getParsedRecord(visit.soap).spo2) < 95">
                          · Di luar acuan
                        </template>
                      </strong>
                    </div>

                    <div>
                      <small>Berat badan</small>
                      <strong> {{ getParsedRecord(visit.soap).weight || '—' }} kg </strong>
                    </div>

                    <div>
                      <small>Tinggi badan</small>
                      <strong> {{ getParsedRecord(visit.soap).height || '—' }} cm </strong>
                    </div>
                  </div>

                  <p class="muted">
                    Penanda memakai acuan umum dewasa; interpretasikan sesuai kondisi pasien.
                  </p>

                  <div class="record-note">
                    <strong>S · Subjective</strong>
                    <p>
                      {{
                        getParsedRecord(visit.soap).subjective ||
                        [visit.complaint, visit.symptoms ? visit.symptoms.join(', ') : '']
                          .filter(Boolean)
                          .join('\n') ||
                        'Belum tercatat'
                      }}
                    </p>
                  </div>

                  <div class="record-note">
                    <strong>O · Objective</strong>
                    <p>
                      {{ getParsedRecord(visit.soap).objective || 'Belum tercatat' }}
                    </p>
                  </div>

                  <div class="record-note">
                    <strong>A · Assessment / ICD-10</strong>
                    <p>
                      {{ getParsedRecord(visit.soap).diagnosis || 'Belum tercatat' }}
                    </p>
                  </div>

                  <div class="record-note">
                    <strong>P · Plan</strong>
                    <p>
                      {{ getParsedRecord(visit.soap).plan || 'Belum tercatat' }}
                    </p>
                  </div>
                </details>
              </div>
            </div>
          </div>
        </div>
      </section>

      <div v-else class="empty-state">
        <CalendarDays :size="40" />
        <h2>Pilih pasien terlebih dahulu</h2>
        <button class="button primary" @click="router.push('/staff')">
          Kembali ke Dashboard Staf
        </button>
      </div>
    </div>
  </main>
</template>
