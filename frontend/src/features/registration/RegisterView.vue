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
              <h3>Pilih layanan, jadwal & shift dokter</h3>
              <p>Pilih poli tujuan dan shift kedatangan yang sesuai dengan waktu Anda.</p>
            </div>
          </div>
          <div class="field-heading-row">
            <p class="field-heading" style="margin: 0;">1. Pilih tanggal kunjungan *</p>
            <span class="field-subtitle-note">Hingga {{ maxDateLabel }}</span>
          </div>

          <div class="date-quick-picker">
            <button
              v-for="d in allowedDates"
              :key="d.dateStr"
              type="button"
              :class="[
                'date-card',
                { selected: queueStore.form.date === d.dateStr },
              ]"
              @click="queueStore.updateForm('date', d.dateStr)"
            >
              <span class="date-card-title">{{ d.title }}</span>
              <strong class="date-card-day">{{ d.dayLabel }}</strong>
              <CheckCircle2
                v-if="queueStore.form.date === d.dateStr"
                :size="16"
                class="date-check-icon"
              />
            </button>
          </div>

          <div class="date-manual-input-wrap">
            <label class="date-manual-label">
              <span>Atau pilih via kalender:</span>
              <input
                type="date"
                required
                :min="minDate"
                :max="maxDate"
                :value="queueStore.form.date"
                @input="handleDateChange"
                class="date-input-control"
              />
            </label>
            <small class="date-restriction-note">
              * Pemilihan tanggal dibuka untuk hari saat mendaftar, besok, dan lusa.
            </small>
          </div>

          <p class="field-heading">2. Pilih poli tujuan *</p>
          <div class="poli-grid">
            <button
              v-for="s in services"
              :key="s.name"
              type="button"
              :class="[
                'poli-choice',
                { selected: queueStore.form.poli === s.name },
              ]"
              @click="queueStore.updateForm('poli', s.name)"
            >
              <component :is="getPoliIcon(s.iconName)" :size="25" />
              <strong>{{ s.name }}</strong>
              <small>3 Sesi Shift Dokter</small>
              <span>{{ getTotalPoliRemainingQuota(s.name) }} kuota tersedia</span>
              <CheckCircle2 v-if="queueStore.form.poli === s.name" :size="19" />
            </button>
          </div>

          <!-- 3. Pilih Shift Kedatangan -->
          <div class="shift-section">
            <div class="shift-heading-row">
              <p class="field-heading" style="margin: 0;">3. Pilih shift kedatangan Anda *</p>
            </div>
            <p class="shift-subtext">
              Pilih waktu pelayanan yang sesuai dengan jadwal kunjungan Anda.
            </p>

            <div class="shift-grid">
              <button
                v-for="shift in shiftList"
                :key="shift.id"
                type="button"
                :disabled="getShiftInfo(shift.id).isFull"
                :class="[
                  'shift-card',
                  { selected: queueStore.form.shift === shift.id, full: getShiftInfo(shift.id).isFull }
                ]"
                @click="queueStore.updateForm('shift', shift.id)"
              >
                <div class="shift-card-top">
                  <span class="shift-name-tag">{{ shift.label }}</span>
                  <span v-if="getShiftInfo(shift.id).isFull" class="status-tag-full">
                    Penuh
                  </span>
                  <span v-else class="status-tag-avail">
                    Sisa {{ getShiftInfo(shift.id).remainingQuota }} kuota
                  </span>
                </div>

                <div class="shift-time-main">
                  <Clock3 :size="16" class="time-icon" />
                  <strong>{{ shift.patientHours }} WIB</strong>
                </div>

                <div class="shift-doctor-line">
                  <UserCheck :size="14" class="text-slate-400 shrink-0" />
                  <span>{{ getShiftInfo(shift.id).doctorName }}</span>
                </div>

                <div class="shift-card-bottom">
                  <span class="est-tag">
                    Estimasi tiba: <strong>{{ getShiftInfo(shift.id).estimatedArrival }}</strong>
                  </span>
                  <CheckCircle2 v-if="queueStore.form.shift === shift.id" :size="17" class="check-icon" />
                </div>
              </button>
            </div>

            <!-- Detail Pilihan Shift & Ketentuan Waktu -->
            <div v-if="selectedShiftInfo" class="shift-summary-box">
              <div class="summary-grid">
                <div class="summary-item">
                  <span class="label">Dokter Bertugas</span>
                  <strong class="val">{{ selectedShiftInfo.doctorName }}</strong>
                  <small>{{ selectedShiftInfo.doctorSpecialty }} • {{ selectedShiftInfo.room }}</small>
                </div>
                <div class="summary-item">
                  <span class="label">Estimasi Pemeriksaan</span>
                  <strong class="val highlight">{{ selectedShiftInfo.estimatedArrival }}</strong>
                  <small>Antrean ke-{{ selectedShiftInfo.nextQueueOrder }} di sesi ini</small>
                </div>
                <div class="summary-item arrival-warning">
                  <span class="label">Disarankan Hadir</span>
                  <strong class="val early-val">{{ selectedShiftInfo.recommendedArrival }}</strong>
                  <small>15 menit sebelum estimasi</small>
                </div>
              </div>
              <div class="summary-notice-bar">
                <Clock3 :size="15" class="notice-icon shrink-0" />
                <span>
                  Harap hadir paling lambat pukul <strong>{{ selectedShiftInfo.recommendedArrival }}</strong> untuk verifikasi berkas dan skrining awal di loket klinik.
                </span>
              </div>
            </div>
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
import { ref, computed } from "vue";
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
  UserCheck,
  AlertCircle,
} from "lucide-vue-next";
import { services } from "../../constants/services";
import { todayDate } from "../../constants/initialData";
import {
  getRelativeDateStr,
  formatShortIndoDate,
  formatDayAndDate,
} from "@/utils/dateHelper";
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

const dateToday = getRelativeDateStr(0);
const dateTomorrow = getRelativeDateStr(1);
const dateDayAfter = getRelativeDateStr(2);

const minDate = computed(() => dateToday);
const maxDate = computed(() => dateDayAfter);
const maxDateLabel = computed(() => formatShortIndoDate(dateDayAfter));

const allowedDates = computed(() => [
  {
    dateStr: dateToday,
    title: "Hari Ini",
    dayLabel: formatDayAndDate(dateToday),
  },
  {
    dateStr: dateTomorrow,
    title: "Besok",
    dayLabel: formatDayAndDate(dateTomorrow),
  },
  {
    dateStr: dateDayAfter,
    title: "Lusa",
    dayLabel: formatDayAndDate(dateDayAfter),
  },
]);

function handleDateChange(e: Event) {
  const input = e.target as HTMLInputElement;
  let val = input.value;
  if (!val) return;
  if (val < minDate.value) val = minDate.value;
  if (val > maxDate.value) val = maxDate.value;
  input.value = val;
  queueStore.updateForm("date", val);
}

// Ensure form.date is within allowed range [minDate, maxDate]
if (
  !queueStore.form.date ||
  queueStore.form.date < dateToday ||
  queueStore.form.date > dateDayAfter
) {
  queueStore.updateForm("date", dateToday);
}

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

const shiftList = [
  {
    id: "pagi" as const,
    label: "Shift Pagi",
    patientHours: "07.00 – 11.00",
    doctorHours: "07.00 – 12.00",
  },
  {
    id: "siang" as const,
    label: "Shift Siang",
    patientHours: "12.00 – 16.00",
    doctorHours: "12.00 – 17.00",
  },
  {
    id: "sore" as const,
    label: "Shift Sore",
    patientHours: "17.00 – 21.00",
    doctorHours: "17.00 – 22.00",
  },
];

function getShiftInfo(shiftId: "pagi" | "siang" | "sore") {
  return queueStore.getShiftQuotaInfo(
    queueStore.form.poli,
    shiftId,
    queueStore.form.date
  );
}

const selectedShiftInfo = computed(() => {
  return getShiftInfo((queueStore.form.shift as "pagi" | "siang" | "sore") || "pagi");
});

function getTotalPoliRemainingQuota(poliName: string) {
  const pagi = queueStore.getShiftQuotaInfo(poliName, "pagi", queueStore.form.date).remainingQuota;
  const siang = queueStore.getShiftQuotaInfo(poliName, "siang", queueStore.form.date).remainingQuota;
  const sore = queueStore.getShiftQuotaInfo(poliName, "sore", queueStore.form.date).remainingQuota;
  return pagi + siang + sore;
}

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

/* Shift Section Styling - Clean, Professional Clinic UX */
.shift-section {
  margin-top: 24px;
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  border-radius: 16px;
  padding: 20px;
}
.shift-heading-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 6px;
}
.shift-subtext {
  font-size: 13px;
  color: #64748b;
  margin: 0 0 16px 0;
  line-height: 1.45;
}
.shift-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 14px;
  margin-bottom: 18px;
}
@media (max-width: 820px) {
  .shift-grid {
    grid-template-columns: 1fr;
  }
}
.shift-card {
  background: #ffffff;
  border: 1.5px solid #cbd5e1;
  border-radius: 14px;
  padding: 16px;
  text-align: left;
  cursor: pointer;
  transition: all 0.2s ease;
  display: flex;
  flex-direction: column;
  gap: 10px;
  position: relative;
}
.shift-card:hover:not(:disabled) {
  border-color: #0f766e;
  box-shadow: 0 4px 12px rgba(15, 118, 110, 0.08);
  transform: translateY(-2px);
}
.shift-card.selected {
  border-color: #0f766e;
  background: #ffffff;
  box-shadow: 0 0 0 2px #0f766e;
}
.shift-card.full,
.shift-card:disabled {
  opacity: 0.6;
  background: #f1f5f9;
  cursor: not-allowed;
  border-color: #e2e8f0;
}
.shift-card-top {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.shift-name-tag {
  font-size: 12px;
  font-weight: 700;
  color: #1e293b;
  background: #f1f5f9;
  padding: 3px 8px;
  border-radius: 6px;
  border: 1px solid #e2e8f0;
}
.status-tag-avail {
  font-size: 11px;
  font-weight: 600;
  color: #0f766e;
  background: #f0fdfa;
  padding: 2px 7px;
  border-radius: 4px;
}
.status-tag-full {
  font-size: 11px;
  font-weight: 700;
  color: #dc2626;
  background: #fee2e2;
  padding: 2px 7px;
  border-radius: 4px;
}
.shift-time-main {
  display: flex;
  align-items: center;
  gap: 7px;
  font-size: 15px;
  color: #0f172a;
}
.shift-time-main strong {
  font-weight: 800;
}
.time-icon {
  color: #0f766e;
}
.shift-doctor-line {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12.5px;
  color: #475569;
  font-weight: 500;
}
.shift-card-bottom {
  margin-top: 2px;
  padding-top: 8px;
  border-top: 1px dashed #e2e8f0;
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.est-tag {
  font-size: 11.5px;
  color: #64748b;
}
.est-tag strong {
  color: #0f766e;
  font-weight: 700;
}
.check-icon {
  color: #0f766e;
}

/* Appointment Summary Box */
.shift-summary-box {
  background: #ffffff;
  border: 1px solid #cbd5e1;
  border-radius: 14px;
  padding: 16px 18px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}
.summary-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16px;
}
@media (max-width: 680px) {
  .summary-grid {
    grid-template-columns: 1fr;
    gap: 12px;
  }
}
.summary-item {
  display: flex;
  flex-direction: column;
  gap: 2px;
}
.summary-item .label {
  font-size: 11px;
  color: #64748b;
  text-transform: uppercase;
  font-weight: 600;
  letter-spacing: 0.3px;
}
.summary-item .val {
  font-size: 15px;
  color: #0f172a;
  font-weight: 700;
}
.summary-item .val.highlight {
  color: #0f766e;
  font-size: 16px;
  font-weight: 800;
}
.summary-item.arrival-warning .early-val {
  color: #b45309;
  font-size: 16px;
  font-weight: 800;
}
.summary-item small {
  font-size: 11.5px;
  color: #64748b;
}
.summary-notice-bar {
  display: flex;
  align-items: center;
  gap: 8px;
  background: #fffbeb;
  border: 1px solid #fef3c7;
  border-radius: 8px;
  padding: 8px 12px;
  font-size: 12px;
  color: #92400e;
  line-height: 1.4;
}
.notice-icon {
  color: #d97706;
}

/* Date Quick Picker & Restriction Styles */
.field-heading-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-top: 10px;
  margin-bottom: 8px;
}
.field-subtitle-note {
  font-size: 12px;
  color: #0f766e;
  font-weight: 600;
  background: #f0fdf4;
  border: 1px solid #bbf7d0;
  padding: 3px 10px;
  border-radius: 999px;
}
.date-quick-picker {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 10px;
  margin-bottom: 12px;
}
@media (max-width: 640px) {
  .date-quick-picker {
    grid-template-columns: 1fr;
  }
}
.date-card {
  position: relative;
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  text-align: left;
  padding: 12px 14px;
  background: #ffffff;
  border: 1.5px solid #cbd5e1;
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.2s ease;
}
.date-card:hover {
  border-color: #0d9488;
  background: #f0fdfa;
}
.date-card.selected {
  border-color: #0f766e;
  background: #f0fdfa;
  box-shadow: 0 0 0 2px rgba(15, 118, 110, 0.15);
}
.date-card-title {
  font-size: 13px;
  font-weight: 700;
  color: #0f172a;
}
.date-card.selected .date-card-title {
  color: #0f766e;
}
.date-card-day {
  font-size: 12px;
  font-weight: 500;
  color: #64748b;
  margin-top: 2px;
}
.date-check-icon {
  position: absolute;
  top: 12px;
  right: 12px;
  color: #0f766e;
}
.date-manual-input-wrap {
  display: flex;
  flex-direction: column;
  gap: 4px;
  margin-bottom: 18px;
}
.date-manual-label {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 12.5px;
  color: #475569;
  font-weight: 500;
}
.date-input-control {
  padding: 6px 12px;
  border: 1px solid #cbd5e1;
  border-radius: 8px;
  font-size: 13px;
  outline: none;
  background: #ffffff;
}
.date-input-control:focus {
  border-color: #0f766e;
}
.date-restriction-note {
  font-size: 11.5px;
  color: #64748b;
  font-style: italic;
}
</style>
