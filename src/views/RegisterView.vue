<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import {
  ShieldCheck,
  CalendarDays,
  Stethoscope,
  Clock3,
  Check,
  CheckCircle2,
  ArrowRight,
} from 'lucide-vue-next'
import { services } from '../constants/services'
import { useQueueStore, today } from '../stores/queue'

const router = useRouter()
const route = useRoute()
const queueStore = useQueueStore()

const walkin = computed(() => route.query.walkin === 'true')
const step = ref(1)
const human = ref(false)

const form = reactive({
  name: '',
  nik: '',
  phone: '',
  birth: '',
  gender: '',
  insurance: 'Umum',
  bpjs: '',
  date: today,
  poli: 'Poli Umum',
  complaint: '',
  allergy: '',
  symptoms: [] as string[],
})

onMounted(() => {
  if (route.query.poli && typeof route.query.poli === 'string') {
    form.poli = route.query.poli
  }
})

const updateNik = (val: string) => {
  form.nik = val.replace(/\D/g, '')
}

const updateBpjs = (val: string) => {
  form.bpjs = val.replace(/\D/g, '')
}

const toggleSymptom = (s: string) => {
  if (form.symptoms.includes(s)) {
    form.symptoms = form.symptoms.filter((x) => x !== s)
  } else {
    form.symptoms.push(s)
  }
}

const handleNextOrSubmit = () => {
  if (step.value < 3) {
    step.value++
  } else {
    queueStore.registerPatient(
      {
        name: form.name,
        nik: form.nik,
        phone: form.phone,
        birth: form.birth,
        gender: form.gender,
        insurance: form.insurance,
        bpjs: form.bpjs,
        date: form.date,
        poli: form.poli,
        complaint: form.complaint,
        allergy: form.allergy,
        symptoms: form.symptoms,
      },
      walkin.value,
    )

    if (walkin.value) {
      queueStore.setNotice('Pasien walk-in berhasil didaftarkan.')
      router.push('/staff')
    } else {
      router.push('/ticket')
    }
  }
}

const goBack = () => {
  router.push(walkin.value ? '/staff' : '/')
}
</script>

<template>
  <main class="form-page container">
    <button class="back-link" @click="goBack">
      ← {{ walkin ? 'Kembali ke dashboard' : 'Kembali ke beranda' }}
    </button>

    <div class="center-heading">
      <div class="section-label">PENDAFTARAN {{ walkin ? 'WALK-IN' : 'ONLINE' }}</div>
      <h2>Selangkah lebih dekat dengan sehat.</h2>
      <p>Isi data dengan benar untuk pelayanan yang lebih cepat.</p>
    </div>

    <div class="wizard">
      <div class="stepper">
        <div
          v-for="(s, i) in ['Identitas diri', 'Pilih layanan', 'Keluhan awal']"
          :key="s"
          :class="['step', step >= i + 1 ? 'active' : '']"
        >
          <span>
            <Check v-if="step > i + 1" :size="17" />
            <template v-else>{{ i + 1 }}</template>
          </span>
          <b>{{ s }}</b>
        </div>
      </div>

      <form @submit.prevent="handleNextOrSubmit">
        <!-- Step 1: Identitas Diri -->
        <template v-if="step === 1">
          <div class="form-title">
            <ShieldCheck />
            <div>
              <h3>Data identitas diri</h3>
              <p>Gunakan data yang sesuai dengan KTP Anda.</p>
            </div>
          </div>

          <div class="fields">
            <label class="full">
              Nomor Induk Kependudukan (NIK) *
              <input
                required
                inputmode="numeric"
                pattern="[0-9]{16}"
                maxlength="17"
                :value="form.nik"
                @input="updateNik(($event.target as HTMLInputElement).value)"
                placeholder="Masukkan 16 digit NIK"
              />
              <small v-if="form.nik && form.nik.length !== 16" class="error">
                NIK wajib tepat 16 digit angka.
              </small>
            </label>

            <label class="full">
              Nama lengkap sesuai KTP *
              <input required v-model="form.name" placeholder="Nama lengkap Anda" />
            </label>

            <label>
              Tanggal lahir *
              <input type="date" required :max="today" v-model="form.birth" />
            </label>

            <label>
              Jenis kelamin *
              <select required v-model="form.gender">
                <option value="">Pilih jenis kelamin</option>
                <option value="Laki-laki">Laki-laki</option>
                <option value="Perempuan">Perempuan</option>
              </select>
            </label>

            <label class="full">
              Nomor WhatsApp aktif *
              <input
                type="tel"
                required
                pattern="(\+62|62|0)[0-9]{8,13}"
                v-model="form.phone"
                placeholder="Contoh: 081234567890"
              />
              <small>
                Wajib diisi dengan benar. Tiket antrean akan dikirim ke nomor ini saat layanan
                terintegrasi.
              </small>
            </label>

            <label class="full">
              Penjamin *
              <select v-model="form.insurance">
                <option value="Umum">Umum</option>
                <option value="BPJS">BPJS</option>
                <option value="Asuransi">Asuransi</option>
              </select>
            </label>

            <label v-if="form.insurance === 'BPJS'" class="full">
              Nomor BPJS *
              <input
                required
                pattern="[0-9]{13}"
                inputmode="numeric"
                maxlength="13"
                :value="form.bpjs"
                @input="updateBpjs(($event.target as HTMLInputElement).value)"
                placeholder="13 digit nomor kartu BPJS"
              />
            </label>

            <label v-if="form.insurance === 'Asuransi'" class="full">
              Nama dan nomor kartu asuransi *
              <input required placeholder="Penjamin dan nomor peserta" />
            </label>
          </div>

          <label class="captcha">
            <input type="checkbox" v-model="human" required />
            <span>
              Saya bukan robot
              <small>Simulasi verifikasi keamanan prototype</small>
            </span>
            <ShieldCheck :size="30" />
          </label>

          <div class="info-note">
            <ShieldCheck :size="17" /> Data pribadi Anda hanya digunakan untuk pelayanan kesehatan.
          </div>
        </template>

        <!-- Step 2: Pilih Layanan & Jadwal -->
        <template v-if="step === 2">
          <div class="form-title">
            <CalendarDays />
            <div>
              <h3>Pilih layanan & jadwal</h3>
              <p>Sesuaikan kunjungan dengan kebutuhan Anda.</p>
            </div>
          </div>

          <label>
            Tanggal kunjungan *
            <input type="date" required :min="today" v-model="form.date" />
          </label>

          <p class="field-heading">Pilih poli *</p>
          <div class="poli-grid">
            <button
              v-for="s in services"
              :key="s.name"
              type="button"
              :disabled="
                Math.max(
                  0,
                  s.quota -
                    queueStore.patients.filter(
                      (p) => p.date === form.date && p.poli === s.name && p.status !== 'Dibatalkan',
                    ).length,
                ) === 0
              "
              :class="['poli-choice', form.poli === s.name ? 'selected' : '']"
              @click="form.poli = s.name"
            >
              <component :is="s.icon" :size="25" />
              <strong>{{ s.name }}</strong>
              <small>08.00 – 14.00</small>
              <span>
                {{
                  Math.max(
                    0,
                    s.quota -
                      queueStore.patients.filter(
                        (p) =>
                          p.date === form.date && p.poli === s.name && p.status !== 'Dibatalkan',
                      ).length,
                  )
                }}
                kuota tersedia
              </span>
              <CheckCircle2 v-if="form.poli === s.name" :size="19" />
            </button>
          </div>

          <div class="info-note">
            <Clock3 :size="18" /> Datang 15 menit sebelum estimasi pelayanan. Jam pelayanan dapat
            berubah.
          </div>
        </template>

        <!-- Step 3: Keluhan Awal -->
        <template v-if="step === 3">
          <div class="form-title">
            <Stethoscope />
            <div>
              <h3>Ceritakan keluhan Anda</h3>
              <p>Informasi ini membantu tenaga medis mempersiapkan pelayanan.</p>
            </div>
          </div>

          <label>
            Apa keluhan utama Anda hari ini? *
            <textarea
              required
              rows="4"
              v-model="form.complaint"
              placeholder="Ceritakan keluhan dan sejak kapan Anda merasakannya..."
            />
          </label>

          <p class="field-heading">Gejala yang dirasakan</p>
          <div class="symptoms">
            <label
              v-for="s in ['Demam', 'Batuk', 'Pusing', 'Mual', 'Nyeri', 'Sesak napas']"
              :key="s"
            >
              <input
                type="checkbox"
                :checked="form.symptoms.includes(s)"
                @change="toggleSymptom(s)"
              />
              {{ s }}
            </label>
          </div>

          <label>
            Riwayat alergi obat / makanan
            <textarea
              rows="2"
              v-model="form.allergy"
              placeholder="Tuliskan alergi Anda, atau 'Tidak ada'"
            />
          </label>

          <div class="summary-box">
            <strong>Ringkasan kunjungan</strong>
            <p>{{ form.name }} · {{ form.poli }}</p>
            <p>
              {{
                new Date(form.date).toLocaleDateString('id-ID', {
                  dateStyle: 'long',
                })
              }}
              · {{ form.insurance }}
            </p>
          </div>
        </template>

        <div class="form-actions">
          <button v-if="step > 1" type="button" class="button outline" @click="step--">
            Kembali
          </button>
          <span>Langkah {{ step }} dari 3</span>
          <button class="button primary" type="submit">
            {{ step === 3 ? 'Daftar Sekarang' : 'Lanjutkan' }}
            <ArrowRight :size="18" />
          </button>
        </div>
      </form>
    </div>
  </main>
</template>
