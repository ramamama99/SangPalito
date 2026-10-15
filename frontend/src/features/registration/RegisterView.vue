<template>
  <main class="form-page container">
    <RouterLink :to="queueStore.walkin ? '/staff' : '/'" class="back-link">
      ← {{ queueStore.walkin ? "Kembali ke dashboard" : "Kembali ke beranda" }}
    </RouterLink>
    <div class="center-heading">
      <div class="section-label">
        PENDAFTARAN {{ queueStore.walkin ? "WALK-IN" : "ONLINE" }}
      </div>
      <h2>Selangkah lebih dekat dengan sehat.</h2>
      <p>Isi data dengan benar untuk pelayanan yang lebih cepat.</p>
    </div>
    <div class="wizard">
      <div class="stepper">
        <div
          v-for="(s, i) in stepLabels"
          :key="s"
          :class="['step', { active: queueStore.step >= i + 1 }]"
        >
          <span>
            <Check v-if="queueStore.step > i + 1" :size="17" />
            <template v-else>{{ i + 1 }}</template>
          </span>
          <b>{{ s }}</b>
        </div>
      </div>
      <form @submit.prevent="handleSubmit">
        <!-- Step 1: Identitas Diri -->
        <template v-if="queueStore.step === 1">
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
                maxlength="16"
                :value="queueStore.form.nik"
                @input="handleNikInput"
                placeholder="Masukkan 16 digit NIK"
              />
              <small
                v-if="queueStore.form.nik && queueStore.form.nik.length !== 16"
                class="error"
              >
                NIK wajib tepat 16 digit angka.
              </small>
            </label>
            <label class="full">
              Nama lengkap sesuai KTP *
              <input
                required
                :value="queueStore.form.name"
                @input="
                  queueStore.updateForm(
                    'name',
                    ($event.target as HTMLInputElement).value,
                  )
                "
                placeholder="Nama lengkap Anda"
              />
            </label>
            <label>
              Tanggal lahir *
              <input
                type="date"
                required
                :max="todayDate"
                :value="queueStore.form.birth"
                @input="
                  queueStore.updateForm(
                    'birth',
                    ($event.target as HTMLInputElement).value,
                  )
                "
              />
            </label>
            <label>
              Jenis kelamin *
              <select
                required
                :value="queueStore.form.gender"
                @change="
                  queueStore.updateForm(
                    'gender',
                    ($event.target as HTMLSelectElement).value,
                  )
                "
              >
                <option value="">Pilih jenis kelamin</option>
                <option>Laki-laki</option>
                <option>Perempuan</option>
              </select>
            </label>
            <label class="full">
              Nomor WhatsApp aktif *
              <input
                type="tel"
                required
                pattern="(\+62|62|0)[0-9]{8,13}"
                :value="queueStore.form.phone"
                @input="
                  queueStore.updateForm(
                    'phone',
                    ($event.target as HTMLInputElement).value,
                  )
                "
                placeholder="Contoh: 081234567890"
              />
              <small>
                Wajib diisi dengan benar. Tiket antrean akan dikirim ke nomor
                ini saat layanan terintegrasi.
              </small>
            </label>
            <label class="full">
              Penjamin *
              <select
                :value="queueStore.form.insurance"
                @change="
                  queueStore.updateForm(
                    'insurance',
                    ($event.target as HTMLSelectElement).value,
                  )
                "
              >
                <option>Umum</option>
                <option>BPJS</option>
                <option>Asuransi</option>
              </select>
            </label>
            <label v-if="queueStore.form.insurance === 'BPJS'" class="full">
              Nomor BPJS *
              <input
                required
                pattern="[0-9]{13}"
                inputmode="numeric"
                maxlength="13"
                :value="queueStore.form.bpjs"
                @input="handleBpjsInput"
                placeholder="13 digit nomor kartu BPJS"
              />
            </label>
            <label v-if="queueStore.form.insurance === 'Asuransi'" class="full">
              Nama dan nomor kartu asuransi *
              <input required placeholder="Penjamin dan nomor peserta" />
            </label>
          </div>
          <label class="captcha">
            <input type="checkbox" v-model="queueStore.human" required />
            <span>
              Saya bukan robot
              <small>Simulasi verifikasi keamanan prototype</small>
            </span>
            <ShieldCheck :size="30" />
          </label>
          <div class="info-note">
            <ShieldCheck :size="17" /> Data pribadi Anda hanya digunakan untuk
            pelayanan kesehatan.
          </div>
        </template>

        <!-- Step 2: Layanan & Jadwal -->
        <template v-if="queueStore.step === 2">
          <div class="form-title">
            <CalendarDays />
            <div>
              <h3>Pilih layanan & jadwal</h3>
              <p>Sesuaikan kunjungan dengan kebutuhan Anda.</p>
            </div>
          </div>
          <label>
            Tanggal kunjungan *
            <input
              type="date"
              required
              :min="todayDate"
              :value="queueStore.form.date"
              @input="
                queueStore.updateForm(
                  'date',
                  ($event.target as HTMLInputElement).value,
                )
              "
            />
          </label>
          <p class="field-heading">Pilih poli *</p>
          <div class="poli-grid">
            <button
              v-for="s in services"
              :key="s.name"
              type="button"
              :disabled="getRemainingQuota(s) <= 0"
              :class="[
                'poli-choice',
                { selected: queueStore.form.poli === s.name },
              ]"
              @click="queueStore.updateForm('poli', s.name)"
            >
              <component :is="getPoliIcon(s.iconName)" :size="25" />
              <strong>{{ s.name }}</strong>
              <small>08.00 – 14.00</small>
              <span>{{ getRemainingQuota(s) }} kuota tersedia</span>
              <CheckCircle2 v-if="queueStore.form.poli === s.name" :size="19" />
            </button>
          </div>
          <div class="info-note">
            <Clock3 :size="18" /> Datang 15 menit sebelum estimasi pelayanan.
            Jam pelayanan dapat berubah.
          </div>
        </template>

        <!-- Step 3: Keluhan awal -->
        <template v-if="queueStore.step === 3">
          <div class="form-title">
            <Stethoscope />
            <div>
              <h3>Ceritakan keluhan Anda</h3>
              <p>
                Informasi ini membantu tenaga medis mempersiapkan pelayanan.
              </p>
            </div>
          </div>
          <label>
            Apa keluhan utama Anda hari ini? *
            <textarea
              required
              rows="4"
              :value="queueStore.form.complaint"
              @input="
                queueStore.updateForm(
                  'complaint',
                  ($event.target as HTMLTextAreaElement).value,
                )
              "
              placeholder="Ceritakan keluhan dan sejak kapan Anda merasakannya..."
            />
          </label>
          <p class="field-heading">Gejala yang dirasakan</p>
          <div class="symptoms">
            <label v-for="symptom in symptomList" :key="symptom">
              <input
                type="checkbox"
                :checked="queueStore.form.symptoms.includes(symptom)"
                @change="
                  toggleSymptom(
                    symptom,
                    ($event.target as HTMLInputElement).checked,
                  )
                "
              />
              {{ symptom }}
            </label>
          </div>
          <label>
            Riwayat alergi obat / makanan
            <textarea
              rows="2"
              :value="queueStore.form.allergy"
              @input="
                queueStore.updateForm(
                  'allergy',
                  ($event.target as HTMLTextAreaElement).value,
                )
              "
              placeholder="Tuliskan alergi Anda, atau 'Tidak ada'"
            />
          </label>
          <div class="summary-box">
            <strong>Ringkasan kunjungan</strong>
            <p>{{ queueStore.form.name }} · {{ queueStore.form.poli }}</p>
            <p>
              {{ formatDate(queueStore.form.date) }} ·
              {{ queueStore.form.insurance }}
            </p>
          </div>
        </template>

        <div class="form-actions">
          <button
            v-if="queueStore.step > 1"
            type="button"
            class="button outline"
            @click="queueStore.step--"
          >
            Kembali
          </button>
          <span>Langkah {{ queueStore.step }} dari 3</span>
          <button class="button primary" type="submit" :disabled="queueStore.isRegistering">
            <template v-if="queueStore.isRegistering">
              <span>Mendaftarkan & Membuat QR...</span>
            </template>
            <template v-else>
              {{ queueStore.step === 3 ? "Daftar Sekarang & Buat QR" : "Lanjutkan" }}
              <ArrowRight :size="18" />
            </template>
          </button>
        </div>
      </form>
    </div>
  </main>
</template>

<script setup lang="ts">
import { useRouter } from "vue-router";
import {
  Check,
  ShieldCheck,
  CalendarDays,
  Stethoscope,
  Clock3,
  CheckCircle2,
  ArrowRight,
  Smile,
  Baby,
  Heart,
} from "lucide-vue-next";
import { services } from "../../constants/services";
import { todayDate } from "../../constants/initialData";
import { useQueueStore } from "../../stores/queue";
import type { Service } from "../../types";

const router = useRouter();
const queueStore = useQueueStore();

const stepLabels = ["Identitas diri", "Pilih layanan", "Keluhan awal"];
const symptomList = [
  "Demam",
  "Batuk",
  "Pusing",
  "Mual",
  "Nyeri",
  "Sesak napas",
];

function getPoliIcon(name: string) {
  switch (name) {
    case "Smile":
      return Smile;
    case "Baby":
      return Baby;
    case "Heart":
      return Heart;
    case "Stethoscope":
    default:
      return Stethoscope;
  }
}

function handleBack() {
  if (queueStore.walkin) {
    router.push("/staff");
  } else {
    router.push("/");
  }
  window.scrollTo(0, 0);
}

function handleNikInput(e: Event) {
  const val = (e.target as HTMLInputElement).value.replace(/\D/g, "");
  queueStore.updateForm("nik", val);
}

function handleBpjsInput(e: Event) {
  const val = (e.target as HTMLInputElement).value.replace(/\D/g, "");
  queueStore.updateForm("bpjs", val);
}

function getRemainingQuota(s: Service) {
  const used = queueStore.patients.filter(
    (p) =>
      p.date === queueStore.form.date &&
      p.poli === s.name &&
      p.status !== "Dibatalkan",
  ).length;
  return Math.max(0, s.quota - used);
}

function toggleSymptom(symptom: string, checked: boolean) {
  if (checked) {
    if (!queueStore.form.symptoms.includes(symptom)) {
      queueStore.form.symptoms.push(symptom);
    }
  } else {
    queueStore.form.symptoms = queueStore.form.symptoms.filter(
      (x) => x !== symptom,
    );
  }
}

function formatDate(dateStr: string) {
  if (!dateStr) return "";
  try {
    return new Date(dateStr).toLocaleDateString("id-ID", { dateStyle: "long" });
  } catch {
    return dateStr;
  }
}

async function handleSubmit() {
  if (queueStore.step < 3) {
    queueStore.step++;
    window.scrollTo({ top: 150, behavior: "smooth" });
  } else {
    const result = await queueStore.registerPatient();
    if (result.isWalkin) {
      router.push("/staff");
    } else {
      router.push("/ticket");
    }
    window.scrollTo(0, 0);
  }
}
</script>
