<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import {
  Activity,
  Users,
  Stethoscope,
  Monitor,
  LogOut,
  Bell,
  Search,
  QrCode,
  Menu,
  ChevronRight,
  Volume2,
  CheckCircle2,
} from 'lucide-vue-next'
import { services } from '../constants/services'
import { useQueueStore, today } from '../stores/queue'
import type { Patient } from '../types'

const router = useRouter()
const queueStore = useQueueStore()

const sidebarCollapsed = ref(false)
const searchInput = ref('')

const soap = ref({
  bp_sys: '120',
  bp_dia: '80',
  hr: '80',
  temp: '36.5',
  spo2: '98',
  weight: '60',
  height: '165',
  doctor: 'dr. Siti Rahmawati',
  subjective: '',
  objective: '',
  diagnosis: 'J00 — Acute nasopharyngitis',
  plan: 'Istirahat cukup, perbanyak minum air hangat, obat simtomatik 3x1.',
})

const filteredLoketPatients = computed(() => {
  const query = searchInput.value.toLowerCase().trim()
  return queueStore.patients.filter(
    (p) => p.date === today && `${p.name} ${p.nik} ${p.number}`.toLowerCase().includes(query),
  )
})

const waitingVerificationCount = computed(() => {
  return queueStore.patients.filter((p) => p.status === 'Menunggu verifikasi').length
})

const selectPatientForExam = (p: Patient) => {
  queueStore.selectedPatientId = p.id
  soap.value.subjective = [p.complaint, p.symptoms.length ? `Gejala: ${p.symptoms.join(', ')}` : '']
    .filter(Boolean)
    .join('\n')
  soap.value.objective = 'Kesadaran compos mentis, KU baik.'
}

const chosen = computed(() => queueStore.chosenPatient)

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

const handleSaveSoap = () => {
  if (!chosen.value) return
  const record = {
    ...soap.value,
    pressure: `${soap.value.bp_sys}/${soap.value.bp_dia}`,
    temperature: soap.value.temp,
    pulse: soap.value.hr,
    recordedAt: new Date().toISOString(),
  }

  queueStore.patients = queueStore.patients.map((p) =>
    p.id === chosen.value!.id
      ? {
          ...p,
          status: 'Selesai',
          soap: JSON.stringify(record),
        }
      : p,
  )

  queueStore.selectedPatientId = ''
  queueStore.setNotice('Pemeriksaan disimpan dan kunjungan ditutup.')
}

const notifyVerification = () => {
  queueStore.setNotice(`${waitingVerificationCount.value} pasien menunggu verifikasi`)
}

const todayStats = computed(() => [
  ['Pendaftaran hari ini', queueStore.patients.filter((p) => p.date === today).length],
  [
    'Menunggu verifikasi',
    queueStore.patients.filter(
      (p) => p.date === today && ['Terdaftar', 'Menunggu verifikasi'].includes(p.status),
    ).length,
  ],
  [
    'Siap dilayani',
    queueStore.patients.filter((p) => p.date === today && p.status === 'Hadir').length,
  ],
  [
    'Kunjungan selesai',
    queueStore.patients.filter((p) => p.date === today && p.status === 'Selesai').length,
  ],
])

const staffPoliPatients = computed(() =>
  queueStore.patients.filter(
    (p) =>
      p.poli === queueStore.staffPoli &&
      p.date === today &&
      ['Hadir', 'Dipanggil'].includes(p.status),
  ),
)
</script>

<template>
  <main class="dashboard">
    <!-- Sidebar -->
    <aside class="sidebar" :data-collapsed="sidebarCollapsed ? 'true' : 'false'">
      <button
        type="button"
        class="sidebar-toggle !inline-flex items-center gap-2.5 p-3 mb-4 self-stretch text-left"
        :aria-label="sidebarCollapsed ? 'Tampilkan menu sidebar' : 'Sembunyikan menu sidebar'"
        :aria-expanded="!sidebarCollapsed"
        :title="sidebarCollapsed ? 'Tampilkan menu sidebar' : 'Sembunyikan menu sidebar'"
        @click="sidebarCollapsed = !sidebarCollapsed"
      >
        <Menu :size="20" aria-hidden="true" />
        <span v-if="!sidebarCollapsed" class="sidebar-toggle-label text-sm font-medium">Menu</span>
      </button>

      <template v-if="!sidebarCollapsed">
        <button
          :class="queueStore.role === 'loket' ? 'selected' : ''"
          @click="queueStore.role = 'loket'"
        >
          <Users :size="20" /> Dashboard Loket
        </button>

        <button
          :class="queueStore.role === 'medis' ? 'selected' : ''"
          @click="queueStore.role = 'medis'"
        >
          <Stethoscope :size="20" /> Dashboard Poli
        </button>

        <button @click="router.push('/tv')"><Monitor :size="20" /> Layar TV Publik</button>

        <div class="sidebar-bottom">
          <button @click="router.push('/')"><LogOut :size="18" /> Keluar</button>
        </div>
      </template>
    </aside>

    <!-- Dashboard Body -->
    <div class="dashboard-body">
      <div class="dashboard-heading">
        <div>
          <div class="section-label">PUSKESMAS SEHAT BERSAMA</div>
          <h2>
            {{ queueStore.role === 'loket' ? 'Dashboard Loket' : 'Ruang Pemeriksaan' }}
          </h2>
          <p>
            {{
              new Date().toLocaleDateString('id-ID', {
                dateStyle: 'full',
              })
            }}
          </p>
        </div>

        <button class="button outline" @click="notifyVerification">
          <Bell :size="18" /> Notifikasi
        </button>
      </div>

      <!-- Stats Grid -->
      <div class="stats-grid">
        <div v-for="[t, n] in todayStats" :key="t" class="stat-card">
          <span>{{ t }}</span>
          <strong>{{ n }}</strong>
          <Activity :size="22" />
        </div>
      </div>

      <!-- Role: LOKET -->
      <section v-if="queueStore.role === 'loket'" class="dashboard-panel">
        <div class="panel-heading">
          <h3>Antrean hari ini</h3>
          <button class="button primary small" @click="router.push('/scan')">
            <QrCode :size="17" /> Scan QR Pasien
          </button>
        </div>

        <div class="search-field">
          <Search :size="18" />
          <input placeholder="Cari nama, NIK, atau nomor antrean..." v-model="searchInput" />
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
                  <strong class="queue-number">{{ p.number }}</strong>
                </td>
                <td>
                  <strong>{{ p.name }}</strong>
                  <small>{{ p.nik }}</small>
                </td>
                <td>{{ p.poli }}</td>
                <td>
                  <span :class="['badge', p.source === 'Online' ? '' : 'coral']">
                    {{ p.source }}
                  </span>
                </td>
                <td>
                  <span
                    :class="['status', `status-${p.status === 'Selesai' ? 'done' : 'pending'}`]"
                  >
                    {{ p.status }}
                  </span>
                </td>
                <td>
                  <button
                    v-if="['Terdaftar', 'Menunggu verifikasi'].includes(p.status)"
                    class="text-button"
                    @click="queueStore.verifyPatient = p"
                  >
                    Verifikasi <ChevronRight :size="15" />
                  </button>
                  <button
                    v-else-if="p.status === 'Hadir'"
                    class="text-button"
                    @click="queueStore.callPatient(p)"
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

      <!-- Role: MEDIS / POLI -->
      <template v-else>
        <div class="poli-toolbar">
          <label>
            Poli Anda
            <select v-model="queueStore.staffPoli" @change="queueStore.selectedPatientId = ''">
              <option v-for="s in services" :key="s.name">
                {{ s.name }}
              </option>
            </select>
          </label>
        </div>

        <div class="medical-grid">
          <!-- Patient List -->
          <section class="dashboard-panel patient-list">
            <h3>Antrean {{ queueStore.staffPoli }}</h3>
            <button
              v-for="p in staffPoliPatients"
              :key="p.id"
              :class="['patient', queueStore.selectedPatientId === p.id ? 'selected' : '']"
              @click="selectPatientForExam(p)"
            >
              <strong>
                {{ p.number }}
                <span>{{ p.status }}</span>
              </strong>
              <b>{{ p.name }}</b>
              <small>{{ p.insurance }} · {{ p.gender }}</small>
            </button>

            <p v-if="!staffPoliPatients.length" class="muted">
              Belum ada pasien yang terverifikasi.
            </p>
          </section>

          <!-- SOAP & Examination Form -->
          <section class="dashboard-panel soap-panel">
            <template
              v-if="
                chosen &&
                chosen.poli === queueStore.staffPoli &&
                ['Hadir', 'Dipanggil'].includes(chosen.status)
              "
            >
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
                    <p v-if="chosen.insurance === 'BPJS'">
                      BPJS: {{ chosen.bpjs || 'Belum diisi' }}
                    </p>
                  </div>

                  <button
                    type="button"
                    class="button primary small record-trigger"
                    aria-label="Lihat Riwayat Rekam Medis"
                    title="Lihat Riwayat Rekam Medis"
                    @click="router.push('/medical-record')"
                  >
                    <Activity :size="17" aria-hidden="true" />
                    <span class="record-trigger-label">MEDICAL RECORD</span>
                  </button>
                </div>

                <div v-if="patientAllergies.length" class="record-alert" role="alert">
                  PERINGATAN ALERGI: {{ patientAllergies.join(', ') }}
                </div>
              </div>

              <div class="record-layout">
                <div class="record-work">
                  <div class="summary-box">
                    <strong>Pra-anamnesis pasien</strong>
                    <p>{{ chosen.complaint }}</p>
                    <small>
                      Gejala:
                      {{ chosen.symptoms.join(', ') || 'Tidak diisi' }} · Alergi:
                      {{ chosen.allergy || 'Tidak ada' }}
                    </small>
                  </div>

                  <form @submit.prevent="handleSaveSoap">
                    <h3 class="field-heading">Tanda vital</h3>
                    <div class="vital-grid">
                      <label>
                        Sistolik (mmHg)
                        <input
                          required
                          type="number"
                          min="40"
                          max="300"
                          v-model="soap.bp_sys"
                          placeholder="120"
                        />
                      </label>
                      <label>
                        Diastolik (mmHg)
                        <input
                          required
                          type="number"
                          min="20"
                          max="200"
                          v-model="soap.bp_dia"
                          placeholder="80"
                        />
                      </label>
                      <label>
                        Detak jantung (/menit)
                        <input
                          required
                          type="number"
                          min="20"
                          max="300"
                          v-model="soap.hr"
                          placeholder="80"
                        />
                      </label>
                      <label>
                        Suhu (°C)
                        <input
                          required
                          type="number"
                          step="any"
                          min="25"
                          max="45"
                          v-model="soap.temp"
                          placeholder="36.5"
                        />
                      </label>
                      <label>
                        SpO₂ (%)
                        <input
                          required
                          type="number"
                          min="1"
                          max="100"
                          v-model="soap.spo2"
                          placeholder="98"
                        />
                      </label>
                      <label>
                        Berat badan (kg)
                        <input
                          required
                          type="number"
                          step="any"
                          min="0.1"
                          max="500"
                          v-model="soap.weight"
                          placeholder="60"
                        />
                      </label>
                      <label>
                        Tinggi badan (cm)
                        <input
                          required
                          type="number"
                          step="any"
                          min="20"
                          max="260"
                          v-model="soap.height"
                          placeholder="165"
                        />
                      </label>
                    </div>

                    <label>
                      Dokter pemeriksa
                      <input
                        required
                        v-model="soap.doctor"
                        placeholder="Nama lengkap dokter"
                        autocomplete="off"
                      />
                    </label>

                    <label>
                      S · Subjektif
                      <textarea required rows="2" v-model="soap.subjective" />
                    </label>

                    <label>
                      O · Objektif
                      <textarea required rows="2" v-model="soap.objective" />
                    </label>

                    <label>
                      A · Diagnosis ICD-10
                      <input
                        required
                        list="diagnoses-list"
                        v-model="soap.diagnosis"
                        placeholder="Cari kode atau diagnosis..."
                      />
                      <datalist id="diagnoses-list">
                        <option value="J00 — Acute nasopharyngitis" />
                        <option value="J06.9 — Infeksi saluran pernapasan atas" />
                        <option value="R50.9 — Demam" />
                        <option value="I10 — Hipertensi esensial" />
                        <option value="E11 — Diabetes melitus tipe 2" />
                        <option value="K02.9 — Karies gigi" />
                        <option value="R51 — Sakit kepala" />
                      </datalist>
                    </label>

                    <label>
                      P · Rencana tindakan / terapi
                      <textarea required rows="3" v-model="soap.plan" />
                    </label>

                    <button class="button primary full-width" type="submit">
                      <CheckCircle2 :size="18" /> Simpan & Tutup Kunjungan
                    </button>
                  </form>
                </div>
              </div>
            </template>

            <div v-else class="empty-state">
              <Stethoscope :size="44" />
              <h3>Pilih pasien untuk pemeriksaan</h3>
              <p>Pasien yang telah diverifikasi akan muncul di panel kiri.</p>
            </div>
          </section>
        </div>
      </template>
    </div>
  </main>
</template>
