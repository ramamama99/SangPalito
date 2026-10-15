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
          <div class="form-title-wrap">
            <div class="form-title">
              <ShieldCheck />
              <div>
                <h3>Data identitas diri</h3>
                <p>Gunakan data yang sesuai dengan KTP Anda.</p>
              </div>
            </div>
            <button
              type="button"
              class="scan-ktp-btn"
              @click="isKtpModalOpen = true"
              title="Pindai fisik e-KTP langsung lewat kamera atau unggah foto"
            >
              <Scan :size="16" />
              <span>Scan KTP Otomatis</span>
            </button>
          </div>

          <!-- Alert Notifikasi Saat Scan Sukses -->
          <div v-if="scanSuccessNotice" class="scan-success-alert">
            <CheckCircle2 :size="18" />
            <span>{{ scanSuccessNotice }}</span>
          </div>

          <div class="fields">
            <!-- 3.a Nomor Induk Kependudukan (NIK) -->
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
                NIK wajib tepat 16 digit angka (saat ini: {{ queueStore.form.nik.length }} digit).
              </small>
            </label>

            <!-- 3.b Nama Lengkap -->
            <label class="full">
              Nama lengkap sesuai KTP *
              <input
                required
                :value="queueStore.form.name"
                @input="
                  queueStore.updateForm(
                    'name',
                    sanitizeCleanText(($event.target as HTMLInputElement).value),
                  )
                "
                placeholder="Nama lengkap Anda sesuai KTP"
              />
            </label>

            <!-- 3.c Tempat Lahir (Kota) -->
            <label>
              Tempat lahir (Kota) *
              <input
                required
                :value="queueStore.form.birthPlace"
                @input="
                  queueStore.updateForm(
                    'birthPlace',
                    sanitizeCleanText(($event.target as HTMLInputElement).value),
                  )
                "
                placeholder="Contoh: Jakarta"
              />
            </label>

            <!-- 3.d Tanggal Lahir -->
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

            <!-- 3.e Jenis Kelamin -->
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
                <option value="Laki-laki">Laki-laki</option>
                <option value="Perempuan">Perempuan</option>
              </select>
            </label>

            <!-- 3.f Golongan Darah -->
            <label>
              Golongan darah
              <select
                :value="queueStore.form.bloodType"
                @change="
                  queueStore.updateForm(
                    'bloodType',
                    ($event.target as HTMLSelectElement).value,
                  )
                "
              >
                <option value="">Pilih golongan darah</option>
                <option value="A">A</option>
                <option value="B">B</option>
                <option value="AB">AB</option>
                <option value="O">O</option>
                <option value="-">- (Tidak tahu)</option>
              </select>
            </label>

            <!-- 3.g Alamat (Jalan / No. / Blok) -->
            <label class="full">
              Alamat (Jalan / No. / Blok) *
              <input
                required
                :value="queueStore.form.address"
                @input="
                  queueStore.updateForm(
                    'address',
                    sanitizeCleanText(($event.target as HTMLInputElement).value),
                  )
                "
                placeholder="Contoh: Jl. Pasirhejo IV No. 16"
              />
            </label>

            <!-- 3.h RT -->
            <label>
              RT *
              <input
                required
                inputmode="numeric"
                maxlength="3"
                :value="queueStore.form.rt"
                @input="
                  queueStore.updateForm(
                    'rt',
                    cleanCharForDigit(($event.target as HTMLInputElement).value).replace(/\D/g, '').slice(0, 3),
                  )
                "
                placeholder="Contoh: 007"
              />
            </label>

            <!-- 3.i RW -->
            <label>
              RW *
              <input
                required
                inputmode="numeric"
                maxlength="3"
                :value="queueStore.form.rw"
                @input="
                  queueStore.updateForm(
                    'rw',
                    cleanCharForDigit(($event.target as HTMLInputElement).value).replace(/\D/g, '').slice(0, 3),
                  )
                "
                placeholder="Contoh: 001"
              />
            </label>

            <!-- 3.j Kelurahan / Desa -->
            <label>
              Kelurahan / Desa *
              <input
                required
                :value="queueStore.form.kelurahan"
                @input="
                  queueStore.updateForm(
                    'kelurahan',
                    sanitizeCleanText(($event.target as HTMLInputElement).value),
                  )
                "
                placeholder="Contoh: Pasirhejo"
              />
            </label>

            <!-- 3.k Kecamatan -->
            <label>
              Kecamatan *
              <input
                required
                :value="queueStore.form.kecamatan"
                @input="
                  queueStore.updateForm(
                    'kecamatan',
                    sanitizeCleanText(($event.target as HTMLInputElement).value),
                  )
                "
                placeholder="Contoh: Cibeunying"
              />
            </label>

            <!-- 4.a Nomor WhatsApp aktif -->
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
                Wajib diisi dengan benar. Tiket antrean dan pengingat panggilan akan dikirim ke nomor ini.
              </small>
            </label>

            <!-- 4.b Penjamin -->
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
                <option value="Umum">Umum (Biaya Mandiri)</option>
                <option value="BPJS">BPJS Kesehatan</option>
                <option value="Asuransi">Asuransi Swasta</option>
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
              <input
                required
                :value="queueStore.form.insuranceName"
                @input="
                  queueStore.updateForm(
                    'insuranceName',
                    sanitizeCleanText(($event.target as HTMLInputElement).value),
                  )
                "
                placeholder="Penjamin dan nomor peserta asuransi"
              />
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

    <!-- Modal Pop-up Scan e-KTP -->
    <KtpScanModal
      v-model="isKtpModalOpen"
      @apply-ktp="handleKtpApplied"
    />
  </main>
</template>

<script setup lang="ts">
import { ref } from "vue";
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
  Scan,
} from "lucide-vue-next";
import { services } from "../../constants/services";
import { todayDate } from "../../constants/initialData";
import { useQueueStore } from "../../stores/queue";
import type { Service } from "../../types";
import KtpScanModal from "./components/KtpScanModal.vue";
import {
  type KtpData,
  sanitizeCleanText,
  sanitizeNik,
  cleanCharForDigit,
} from "@/utils/ktpParser";

const router = useRouter();
const queueStore = useQueueStore();

const isKtpModalOpen = ref(false);
const scanSuccessNotice = ref("");

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
  const val = sanitizeNik((e.target as HTMLInputElement).value);
  queueStore.updateForm("nik", val);
}

function handleBpjsInput(e: Event) {
  const val = sanitizeNik((e.target as HTMLInputElement).value).slice(0, 13);
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

function convertKtpDateToIso(val: string): string {
  if (!val) return "";
  const cleaned = val.trim();
  const m = cleaned.match(/^(\d{1,2})[-/.](\d{1,2})[-/.](\d{4})$/);
  if (m) {
    const day = m[1].padStart(2, "0");
    const month = m[2].padStart(2, "0");
    const year = m[3];
    return `${year}-${month}-${day}`;
  }
  if (/^\d{4}-\d{2}-\d{2}$/.test(cleaned)) return cleaned;
  return cleaned;
}

function handleKtpApplied(data: KtpData) {
  if (data.nik) {
    queueStore.updateForm("nik", sanitizeNik(data.nik));
  }
  if (data.nama) {
    queueStore.updateForm("name", sanitizeCleanText(data.nama));
  }
  if (data.tempatLahir) {
    queueStore.updateForm("birthPlace", sanitizeCleanText(data.tempatLahir));
  }
  if (data.tglLahir) {
    const isoDate = convertKtpDateToIso(data.tglLahir);
    queueStore.updateForm("birth", isoDate);
  }
  if (data.jenisKelamin) {
    const g = data.jenisKelamin.toUpperCase();
    if (g.includes("PEREMPUAN") || g.includes("WANITA")) {
      queueStore.updateForm("gender", "Perempuan");
    } else {
      queueStore.updateForm("gender", "Laki-laki");
    }
  }
  if (data.golDarah) {
    queueStore.updateForm("bloodType", data.golDarah);
  }
  if (data.alamatJalan) {
    queueStore.updateForm("address", sanitizeCleanText(data.alamatJalan));
  }
  if (data.rt) {
    queueStore.updateForm("rt", cleanCharForDigit(data.rt).replace(/\D/g, '').slice(0, 3));
  }
  if (data.rw) {
    queueStore.updateForm("rw", cleanCharForDigit(data.rw).replace(/\D/g, '').slice(0, 3));
  }
  if (data.kelDesa) {
    queueStore.updateForm("kelurahan", sanitizeCleanText(data.kelDesa));
  }
  if (data.kecamatan) {
    queueStore.updateForm("kecamatan", sanitizeCleanText(data.kecamatan));
  }

  scanSuccessNotice.value = `Data e-KTP a.n. ${data.nama || "Pasien"} berhasil diverifikasi dan dimasukkan ke formulir!`;
  setTimeout(() => {
    scanSuccessNotice.value = "";
  }, 7000);
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

<style scoped>
.form-title-wrap {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 16px;
  margin-bottom: 22px;
}
.form-title {
  margin-bottom: 0;
}
.scan-ktp-btn {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  background: linear-gradient(135deg, #e6f7f6 0%, #f0faf9 100%);
  color: var(--deep, #195a5f);
  border: 1.5px solid var(--teal, #12b8b0);
  border-radius: 12px;
  padding: 9px 18px;
  font-size: 13.5px;
  font-weight: 650;
  cursor: pointer;
  transition: all 0.2s cubic-bezier(0.16, 1, 0.3, 1);
  box-shadow: 0 2px 8px rgba(18, 184, 176, 0.15);
  white-space: nowrap;
  shrink: 0;
}
.scan-ktp-btn:hover {
  background: var(--teal, #12b8b0);
  color: #fff;
  transform: translateY(-1.5px);
  box-shadow: 0 4px 14px rgba(18, 184, 176, 0.35);
}
.scan-ktp-btn:active {
  transform: translateY(0);
}
.scan-success-alert {
  display: flex;
  align-items: center;
  gap: 10px;
  background: #ebfbf9;
  border: 1.5px solid #a7ebd8;
  color: #0b6860;
  border-radius: 12px;
  padding: 12px 18px;
  font-size: 13.5px;
  font-weight: 600;
  margin-bottom: 22px;
  box-shadow: 0 2px 10px rgba(18, 184, 176, 0.08);
  animation: slideDown 0.3s ease-out;
}
@keyframes slideDown {
  from {
    opacity: 0;
    transform: translateY(-8px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
@media (max-width: 640px) {
  .form-title-wrap {
    flex-direction: column;
    align-items: flex-start;
    gap: 14px;
  }
  .scan-ktp-btn {
    width: 100%;
    justify-content: center;
    padding: 11px;
  }
}
</style>
