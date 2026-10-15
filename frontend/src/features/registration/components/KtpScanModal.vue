<template>
  <Teleport to="body">
    <Transition name="modal-fade">
      <div
        v-if="modelValue"
        class="fixed inset-0 z-[999] flex items-center justify-center p-3 sm:p-5 bg-black/65 backdrop-blur-sm overflow-y-auto"
        @click.self="closeModal"
      >
        <div
          class="relative w-full max-w-2xl bg-white rounded-3xl shadow-2xl border border-slate-100 overflow-hidden flex flex-col max-h-[92vh] my-auto transition-all"
        >
          <!-- Header Modal -->
          <div class="px-5 sm:px-7 pt-5 pb-4 border-b border-slate-100 flex items-start justify-between bg-gradient-to-r from-slate-50 to-white">
            <div class="space-y-1">
              <h2 class="text-xl sm:text-2xl font-bold text-[#195a5f] tracking-tight">
                {{ parsedData ? 'Verifikasi Data Hasil Scan e-KTP' : 'Scan Kartu e-KTP Otomatis' }}
              </h2>
              <p class="text-xs sm:text-sm text-slate-500 max-w-md">
                {{ parsedData
                  ? 'Periksa dan ubah data hasil pembacaan OCR di bawah ini secara manual bila diperlukan.'
                  : 'Arahkan kartu ke kamera atau unggah berkas foto untuk pengisian formulir instan.' }}
              </p>
            </div>

            <!-- Tombol Tutup (X) -->
            <button
              type="button"
              @click="closeModal"
              class="w-9 h-9 rounded-full bg-slate-100 hover:bg-slate-200 text-slate-500 hover:text-slate-800 flex items-center justify-center transition cursor-pointer shrink-0 -mr-1 -mt-1"
              title="Tutup Modal"
            >
              <X class="w-5 h-5" />
            </button>
          </div>

          <!-- Body Modal -->
          <div class="p-5 sm:p-7 overflow-y-auto space-y-5">
            <!-- 1. Kondisi Sedang Memproses Ekstraksi OCR -->
            <div v-if="isScanning" class="py-10 flex flex-col items-center justify-center text-center space-y-4">
              <div class="relative w-16 h-16 flex items-center justify-center">
                <Loader2 class="w-16 h-16 text-[var(--teal,#12b8b0)] animate-spin" />
                <Sparkles class="w-6 h-6 text-[var(--deep,#195a5f)] absolute" />
              </div>
              <div class="space-y-1 max-w-sm">
                <h3 class="text-base font-bold text-slate-800">
                  Menganalisis Teks e-KTP...
                </h3>
                <p class="text-xs text-slate-500 leading-relaxed">
                  {{ statusText || 'Mengevaluasi 4 metode citra (Optimal, Otsu, Kontras, Asli) untuk akurasi terbaik...' }}
                </p>
              </div>

              <!-- Progress Bar -->
              <div class="w-full max-w-xs bg-slate-100 rounded-full h-2 overflow-hidden border border-slate-200">
                <div
                  class="bg-gradient-to-r from-[var(--teal,#12b8b0)] to-[var(--deep,#195a5f)] h-2 transition-all duration-300"
                  :style="{ width: `${progress}%` }"
                ></div>
              </div>
              <span class="text-xs font-semibold text-slate-400">{{ progress }}%</span>
            </div>

            <!-- 2. Kondisi Hasil Ekstraksi Sudah Didapat (Langkah Verifikasi & Konfirmasi Editable) -->
            <div v-else-if="parsedData" class="space-y-5">
              <!-- Banner Sukses & Panduan Verifikasi -->
              <div class="p-4 rounded-2xl bg-emerald-50/90 border border-emerald-200 flex items-start gap-3">
                <div class="w-9 h-9 rounded-xl bg-emerald-600 text-white flex items-center justify-center shrink-0 shadow-xs mt-0.5">
                  <CheckCircle2 class="w-5 h-5" />
                </div>
                <div class="space-y-0.5 flex-1">
                  <h4 class="text-sm font-bold text-emerald-950">
                    Data e-KTP Berhasil Ditangkap!
                  </h4>
                  <p class="text-xs text-emerald-800 leading-relaxed">
                    Filter optimal: <strong>{{ winningMethodLabel || 'Otomatis' }}</strong> (Tingkat akurasi: {{ parsedData.confidence }}%).
                    Seluruh kolom di bawah bersifat <strong>editable</strong>. Anda dapat mengedit dan memperbaiki data manual jika ada kesalahan baca OCR.
                  </p>
                </div>
              </div>

              <!-- Thumbnail Tangkapan e-KTP Asli untuk Perbandingan -->
              <div v-if="previewUrl" class="p-3 bg-slate-50 rounded-2xl border border-slate-200/80 flex items-center justify-between gap-3">
                <div class="flex items-center gap-3 min-w-0">
                  <img
                    :src="previewUrl"
                    alt="Pratinjau KTP"
                    class="w-16 h-11 object-cover rounded-xl border border-slate-300 shadow-2xs shrink-0"
                  />
                  <div class="truncate">
                    <span class="text-[11px] font-bold text-slate-700 block truncate">Foto e-KTP yang Diproses</span>
                    <span class="text-[11px] text-slate-500 block truncate">Bandingkan data form dengan fisik kartu di foto</span>
                  </div>
                </div>
                <button
                  type="button"
                  @click="resetScan"
                  class="text-xs text-[var(--deep,#195a5f)] hover:text-teal-700 font-bold px-3 py-1.5 rounded-lg bg-white border border-slate-200 hover:bg-slate-100 transition shrink-0 flex items-center gap-1 cursor-pointer"
                >
                  <RotateCcw class="w-3.5 h-3.5" />
                  <span>Foto Ulang</span>
                </button>
              </div>

              <!-- Formulir Verifikasi & Koreksi Data (Editable Form) -->
              <div class="bg-white rounded-2xl border border-slate-200/90 shadow-xs p-4 sm:p-5 space-y-4">
                <div class="flex items-center justify-between border-b border-slate-100 pb-3">
                  <div class="flex items-center gap-2">
                    <Edit3 class="w-4 h-4 text-[var(--teal,#12b8b0)]" />
                    <h3 class="text-sm font-bold text-slate-800">
                      Pengecekan & Koreksi Data e-KTP
                    </h3>
                  </div>
                  <span class="text-[11px] font-semibold text-slate-400">
                    Formulir Editable
                  </span>
                </div>

                <!-- Petunjuk Filter Pembersihan Teks -->
                <div class="text-[11px] text-slate-600 bg-slate-50 rounded-xl p-2.5 border border-slate-200/70 flex items-center gap-2">
                  <Info class="w-4 h-4 text-slate-400 shrink-0" />
                  <span>
                    <strong>Filter aktif:</strong> Hanya huruf, angka, tanda titik (.), dan tanda garis miring (/) yang diizinkan masuk ke form.
                  </span>
                </div>

                <!-- Grid Input Data Utama -->
                <div class="grid grid-cols-1 sm:grid-cols-2 gap-3.5">
                  <!-- 1. NIK (16 Digit) -->
                  <div class="sm:col-span-2 space-y-1">
                    <div class="flex items-center justify-between">
                      <label class="text-xs font-bold text-slate-700">
                        Nomor Induk Kependudukan (NIK) *
                      </label>
                      <span
                        v-if="editableForm.nik.length === 16"
                        class="text-[11px] font-bold text-emerald-700 bg-emerald-50 px-2 py-0.5 rounded-full border border-emerald-200 flex items-center gap-1"
                      >
                        <CheckCircle2 class="w-3 h-3 text-emerald-600" />
                        16 Digit Valid
                      </span>
                      <span
                        v-else
                        class="text-[11px] font-bold text-amber-700 bg-amber-50 px-2 py-0.5 rounded-full border border-amber-200"
                      >
                        {{ editableForm.nik.length }}/16 Digit (Wajib 16 Digit)
                      </span>
                    </div>
                    <input
                      type="text"
                      inputmode="numeric"
                      maxlength="16"
                      :value="editableForm.nik"
                      @input="handleNikInput"
                      class="w-full px-3.5 py-2.5 text-sm font-mono font-bold bg-slate-50/70 border rounded-xl transition focus:bg-white focus:outline-none focus:ring-2 focus:ring-[var(--teal,#12b8b0)]"
                      :class="editableForm.nik.length === 16 ? 'border-slate-300 text-slate-900' : 'border-amber-300 text-amber-900 bg-amber-50/30'"
                      placeholder="Masukkan 16 digit angka NIK"
                    />
                    <small class="text-[11px] text-slate-400 block">
                      Karakter huruf dan simbol otomatis dikoreksi menjadi digit angka.
                    </small>
                  </div>

                  <!-- 2. Nama Lengkap Sesuai KTP -->
                  <div class="sm:col-span-2 space-y-1">
                    <label class="text-xs font-bold text-slate-700">
                      Nama Lengkap Sesuai KTP *
                    </label>
                    <input
                      type="text"
                      :value="editableForm.nama"
                      @input="handleTextInput('nama', ($event.target as HTMLInputElement).value)"
                      class="w-full px-3.5 py-2.5 text-sm font-semibold text-slate-900 bg-slate-50/70 border border-slate-300 rounded-xl transition focus:bg-white focus:outline-none focus:ring-2 focus:ring-[var(--teal,#12b8b0)]"
                      placeholder="Nama lengkap sesuai KTP"
                    />
                  </div>

                  <!-- 3. Tempat Lahir (Kota) -->
                  <div class="space-y-1">
                    <label class="text-xs font-bold text-slate-700">
                      Tempat Lahir (Kota) *
                    </label>
                    <input
                      type="text"
                      :value="editableForm.tempatLahir"
                      @input="handleTextInput('tempatLahir', ($event.target as HTMLInputElement).value)"
                      class="w-full px-3.5 py-2.5 text-sm text-slate-900 bg-slate-50/70 border border-slate-300 rounded-xl transition focus:bg-white focus:outline-none focus:ring-2 focus:ring-[var(--teal,#12b8b0)]"
                      placeholder="Contoh: Jakarta"
                    />
                  </div>

                  <!-- 4. Tanggal Lahir (Format: DD/MM/YYYY) -->
                  <div class="space-y-1">
                    <label class="text-xs font-bold text-slate-700">
                      Tanggal Lahir (DD/MM/YYYY) *
                    </label>
                    <input
                      type="text"
                      :value="editableForm.tglLahir"
                      @input="handleTextInput('tglLahir', ($event.target as HTMLInputElement).value)"
                      class="w-full px-3.5 py-2.5 text-sm text-slate-900 bg-slate-50/70 border border-slate-300 rounded-xl transition focus:bg-white focus:outline-none focus:ring-2 focus:ring-[var(--teal,#12b8b0)]"
                      placeholder="Contoh: 19/05/2006"
                    />
                  </div>

                  <!-- 5. Jenis Kelamin -->
                  <div class="space-y-1">
                    <label class="text-xs font-bold text-slate-700">
                      Jenis Kelamin *
                    </label>
                    <select
                      v-model="editableForm.jenisKelamin"
                      class="w-full px-3.5 py-2.5 text-sm text-slate-900 bg-slate-50/70 border border-slate-300 rounded-xl transition focus:bg-white focus:outline-none focus:ring-2 focus:ring-[var(--teal,#12b8b0)] cursor-pointer"
                    >
                      <option value="LAKI-LAKI">LAKI-LAKI</option>
                      <option value="PEREMPUAN">PEREMPUAN</option>
                    </select>
                  </div>

                  <!-- 6. Golongan Darah -->
                  <div class="space-y-1">
                    <label class="text-xs font-bold text-slate-700">
                      Golongan Darah
                    </label>
                    <select
                      v-model="editableForm.golDarah"
                      class="w-full px-3.5 py-2.5 text-sm text-slate-900 bg-slate-50/70 border border-slate-300 rounded-xl transition focus:bg-white focus:outline-none focus:ring-2 focus:ring-[var(--teal,#12b8b0)] cursor-pointer"
                    >
                      <option value="-">- (Tidak Tahu)</option>
                      <option value="A">A</option>
                      <option value="B">B</option>
                      <option value="AB">AB</option>
                      <option value="O">O</option>
                    </select>
                  </div>
                </div>

                <!-- BAGIAN PENGECEKAN ALAMAT SESUAI KTP (TERPISAH SESUAI SYARAT 1.c) -->
                <div class="pt-3 border-t border-slate-200/80 space-y-3">
                  <div class="flex items-center gap-2">
                    <MapPin class="w-4 h-4 text-[var(--teal,#12b8b0)]" />
                    <div>
                      <h4 class="text-xs font-bold uppercase tracking-wider text-slate-800">
                        Pengecekan Alamat Sesuai KTP (Kolom Terpisah)
                      </h4>
                      <p class="text-[11px] text-slate-500">
                        Data alamat dipecah ke kolom Jalan, RT, RW, Kelurahan/Desa, dan Kecamatan.
                      </p>
                    </div>
                  </div>

                  <div class="grid grid-cols-1 sm:grid-cols-2 gap-3 bg-slate-50/80 p-3.5 rounded-2xl border border-slate-200">
                    <!-- 1. Jalan / No. / Blok (Full width) -->
                    <div class="sm:col-span-2 space-y-1">
                      <label class="text-xs font-bold text-slate-700">
                        Jalan / No. / Blok *
                      </label>
                      <input
                        type="text"
                        :value="editableForm.alamatJalan"
                        @input="handleTextInput('alamatJalan', ($event.target as HTMLInputElement).value)"
                        class="w-full px-3.5 py-2.5 text-sm text-slate-900 bg-white border border-slate-300 rounded-xl transition focus:outline-none focus:ring-2 focus:ring-[var(--teal,#12b8b0)]"
                        placeholder="Contoh: Jl. Pasirhejo IV No. 16"
                      />
                    </div>

                    <!-- 2. RT (Kolom Terpisah) -->
                    <div class="space-y-1">
                      <label class="text-xs font-bold text-slate-700">
                        RT *
                      </label>
                      <input
                        type="text"
                        inputmode="numeric"
                        maxlength="3"
                        :value="editableForm.rt"
                        @input="handleDigitInput('rt', ($event.target as HTMLInputElement).value, 3)"
                        class="w-full px-3.5 py-2.5 text-sm text-slate-900 bg-white border border-slate-300 rounded-xl transition focus:outline-none focus:ring-2 focus:ring-[var(--teal,#12b8b0)]"
                        placeholder="007"
                      />
                    </div>

                    <!-- 3. RW (Kolom Terpisah) -->
                    <div class="space-y-1">
                      <label class="text-xs font-bold text-slate-700">
                        RW *
                      </label>
                      <input
                        type="text"
                        inputmode="numeric"
                        maxlength="3"
                        :value="editableForm.rw"
                        @input="handleDigitInput('rw', ($event.target as HTMLInputElement).value, 3)"
                        class="w-full px-3.5 py-2.5 text-sm text-slate-900 bg-white border border-slate-300 rounded-xl transition focus:outline-none focus:ring-2 focus:ring-[var(--teal,#12b8b0)]"
                        placeholder="001"
                      />
                    </div>

                    <!-- 4. Kelurahan atau Desa (Kolom Terpisah) -->
                    <div class="space-y-1">
                      <label class="text-xs font-bold text-slate-700">
                        Kelurahan / Desa *
                      </label>
                      <input
                        type="text"
                        :value="editableForm.kelDesa"
                        @input="handleTextInput('kelDesa', ($event.target as HTMLInputElement).value)"
                        class="w-full px-3.5 py-2.5 text-sm text-slate-900 bg-white border border-slate-300 rounded-xl transition focus:outline-none focus:ring-2 focus:ring-[var(--teal,#12b8b0)]"
                        placeholder="Contoh: Pasirhejo"
                      />
                    </div>

                    <!-- 5. Kecamatan (Kolom Terpisah) -->
                    <div class="space-y-1">
                      <label class="text-xs font-bold text-slate-700">
                        Kecamatan *
                      </label>
                      <input
                        type="text"
                        :value="editableForm.kecamatan"
                        @input="handleTextInput('kecamatan', ($event.target as HTMLInputElement).value)"
                        class="w-full px-3.5 py-2.5 text-sm text-slate-900 bg-white border border-slate-300 rounded-xl transition focus:outline-none focus:ring-2 focus:ring-[var(--teal,#12b8b0)]"
                        placeholder="Contoh: Cibeunying"
                      />
                    </div>
                  </div>
                </div>
              </div>

              <!-- Tombol Aksi Konfirmasi & Verifikasi -->
              <div class="flex items-center gap-3 pt-2">
                <button
                  type="button"
                  @click="applyDataToForm"
                  class="flex-1 py-3 px-5 bg-[var(--teal,#12b8b0)] hover:brightness-105 active:scale-[0.99] text-white font-bold text-sm rounded-xl shadow-md transition flex items-center justify-center gap-2 cursor-pointer"
                >
                  <Check class="w-4 h-4" />
                  <span>Gunakan Data KTP Ini</span>
                </button>
                <button
                  type="button"
                  @click="resetScan"
                  class="py-3 px-4 bg-slate-100 hover:bg-slate-200 text-slate-700 font-semibold text-sm rounded-xl transition cursor-pointer flex items-center gap-1.5"
                >
                  <RotateCcw class="w-4 h-4" />
                  <span>Scan Ulang</span>
                </button>
              </div>
            </div>

            <!-- 3. Kondisi Memilih Opsi Input: Kamera vs Upload Berkas -->
            <div v-else class="space-y-4">
              <!-- Tab Navigasi & Tombol Contoh KTP -->
              <div class="flex items-center justify-between gap-2 flex-wrap">
                <!-- Tab Selector -->
                <div class="inline-flex p-1 bg-slate-100 rounded-xl text-xs font-bold text-slate-600">
                  <button
                    type="button"
                    @click="activeTab = 'camera'"
                    class="flex items-center gap-1.5 py-2 px-3 sm:px-4 rounded-lg transition cursor-pointer"
                    :class="activeTab === 'camera' ? 'bg-white text-[var(--deep,#195a5f)] shadow-xs' : 'hover:text-slate-900'"
                  >
                    <Camera class="w-4 h-4 text-[var(--teal,#12b8b0)]" />
                    <span>Scan Kamera</span>
                  </button>
                  <button
                    type="button"
                    @click="activeTab = 'upload'"
                    class="flex items-center gap-1.5 py-2 px-3 sm:px-4 rounded-lg transition cursor-pointer"
                    :class="activeTab === 'upload' ? 'bg-white text-[var(--deep,#195a5f)] shadow-xs' : 'hover:text-slate-900'"
                  >
                    <UploadCloud class="w-4 h-4 text-[var(--teal,#12b8b0)]" />
                    <span>Upload Foto</span>
                  </button>
                </div>

                <!-- Tombol Contoh KTP untuk Demo Cepat -->
                <button
                  type="button"
                  @click="generateSampleKtp"
                  class="text-xs font-semibold bg-emerald-50 hover:bg-emerald-100 text-emerald-700 border border-emerald-200 px-3 py-1.5 rounded-xl transition cursor-pointer flex items-center gap-1.5 shadow-2xs"
                  title="Simulasikan kartu e-KTP untuk pengujian otomatis langsung"
                >
                  <Sparkles class="w-3.5 h-3.5 text-emerald-600" />
                  <span>Uji Contoh KTP</span>
                </button>
              </div>

              <!-- Pilihan A: Kamera Live Auto-Read -->
              <div v-if="activeTab === 'camera'">
                <KtpCameraScanner
                  @auto-captured="handleCaptured"
                  @switch-to-upload="activeTab = 'upload'"
                />
              </div>

              <!-- Pilihan B: Upload Berkas Foto KTP -->
              <div v-else-if="activeTab === 'upload'">
                <label
                  class="flex flex-col items-center justify-center w-full min-h-[300px] sm:min-h-[340px] border-2 border-dashed rounded-2xl cursor-pointer transition p-6 text-center"
                  :class="isScanning ? 'border-slate-300 bg-slate-50 cursor-not-allowed' : 'border-teal-300 bg-teal-50/20 hover:bg-teal-50/40'"
                >
                  <div class="w-16 h-16 rounded-2xl bg-teal-100 text-[var(--deep,#195a5f)] flex items-center justify-center mb-4 shadow-xs">
                    <UploadCloud class="w-8 h-8 text-[var(--teal,#12b8b0)]" />
                  </div>
                  <h4 class="text-sm font-bold text-slate-800 mb-1">
                    Pilih atau Tarik Berkas Foto e-KTP ke Sini
                  </h4>
                  <p class="text-xs text-slate-500 max-w-sm mb-4 leading-relaxed">
                    Mendukung semua format gambar standar (JPG, JPEG, PNG, WEBP). Pastikan foto e-KTP tegak, pencahayaan merata, dan nomor NIK terbaca jelas.
                  </p>
                  <span class="text-xs font-semibold px-4 py-2 bg-[var(--teal,#12b8b0)] text-white rounded-xl shadow-xs hover:brightness-105 transition">
                    Jelajahi Berkas Foto
                  </span>
                  <input
                    type="file"
                    accept="image/*"
                    class="hidden"
                    :disabled="isScanning"
                    @change="handleFileUpload"
                  />
                </label>
              </div>
            </div>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import {
  X,
  Camera,
  UploadCloud,
  Check,
  CheckCircle2,
  Sparkles,
  Loader2,
  RotateCcw,
  Edit3,
  MapPin,
  Info
} from 'lucide-vue-next'
import KtpCameraScanner from './KtpCameraScanner.vue'
import { processMultiEnhancementOcr } from '@/utils/ocrService'
import {
  type KtpData,
  sanitizeCleanText,
  sanitizeNik,
  cleanCharForDigit
} from '@/utils/ktpParser'
import type { PreprocessOptions } from '@/utils/imageProcessor'

const props = defineProps<{
  modelValue: boolean
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: boolean): void
  (e: 'apply-ktp', data: KtpData): void
}>()

const activeTab = ref<'camera' | 'upload'>('camera')
const isScanning = ref(false)
const progress = ref(0)
const statusText = ref('')
const previewUrl = ref<string | null>(null)
const parsedData = ref<KtpData | null>(null)
const winningMethodLabel = ref<string>('Otomatis')
let rawImageSource: HTMLCanvasElement | File | Blob | string | null = null

// Form State yang bersifat editable untuk verifikasi dan koreksi manual data scan
const editableForm = ref<KtpData>({
  nik: '',
  nama: '',
  tempatLahir: '',
  tglLahir: '',
  tempatTglLahir: '',
  jenisKelamin: 'LAKI-LAKI',
  golDarah: '-',
  alamatJalan: '',
  rt: '',
  rw: '',
  kelDesa: '',
  kecamatan: '',
  confidence: 0,
  isValidNik: false,
})

const closeModal = () => {
  emit('update:modelValue', false)
}

/**
 * Mengisi state editableForm dari hasil pemindaian OCR awal
 */
const populateEditableForm = (data: KtpData) => {
  editableForm.value = {
    nik: sanitizeNik(data.nik),
    nama: sanitizeCleanText(data.nama),
    tempatLahir: sanitizeCleanText(data.tempatLahir),
    tglLahir: sanitizeCleanText(data.tglLahir),
    tempatTglLahir: data.tempatTglLahir || '',
    jenisKelamin: data.jenisKelamin || 'LAKI-LAKI',
    golDarah: data.golDarah || '-',
    alamatJalan: sanitizeCleanText(data.alamatJalan),
    rt: cleanCharForDigit(data.rt || '').replace(/\D/g, '').slice(0, 3),
    rw: cleanCharForDigit(data.rw || '').replace(/\D/g, '').slice(0, 3),
    kelDesa: sanitizeCleanText(data.kelDesa),
    kecamatan: sanitizeCleanText(data.kecamatan),
    confidence: data.confidence,
    isValidNik: Boolean(data.nik && data.nik.length === 16),
  }
}

/**
 * Filter Pembersihan Teks pada Kolom Teks:
 * Hanya mengizinkan huruf, angka, tanda titik (.), tanda garis miring (/), dan spasi.
 */
const handleTextInput = (field: 'nama' | 'tempatLahir' | 'tglLahir' | 'alamatJalan' | 'kelDesa' | 'kecamatan', val: string) => {
  editableForm.value[field] = sanitizeCleanText(val)
}

/**
 * Filter Pembersihan NIK:
 * Hanya mengizinkan digit angka numerik tepat hingga 16 digit.
 */
const handleNikInput = (event: Event) => {
  const target = event.target as HTMLInputElement
  const cleaned = sanitizeNik(target.value)
  editableForm.value.nik = cleaned
  editableForm.value.isValidNik = cleaned.length === 16
}

/**
 * Filter Pembersihan Digit untuk RT dan RW (Maksimal 3 Digit)
 */
const handleDigitInput = (field: 'rt' | 'rw', val: string, maxLen = 3) => {
  editableForm.value[field] = cleanCharForDigit(val).replace(/\D/g, '').slice(0, maxLen)
}

/**
 * Menerima tangkapan langsung dari KtpCameraScanner
 */
const handleCaptured = (payload: {
  canvas: HTMLCanvasElement
  parsed: KtpData
  rawText: string
  winningMethod?: string
  winningMethodLabel?: string
  processedDataUrl?: string
}) => {
  rawImageSource = payload.canvas
  previewUrl.value = payload.processedDataUrl || payload.canvas.toDataURL('image/png')
  parsedData.value = payload.parsed
  winningMethodLabel.value = payload.winningMethodLabel || 'Otomatis (Optimal)'
  populateEditableForm(payload.parsed)
}

/**
 * Handle upload berkas foto KTP dari disk
 */
const handleFileUpload = async (event: Event) => {
  const target = event.target as HTMLInputElement
  if (!target.files || target.files.length === 0) return

  const file = target.files[0]
  rawImageSource = file
  await processAndRunOcr()
}

/**
 * Eksekusi Multi-Enhancement OCR saat upload atau simulasi
 */
const processAndRunOcr = async (forceFilter?: PreprocessOptions['mode']) => {
  if (!rawImageSource) return

  isScanning.value = true
  progress.value = 10
  statusText.value = 'Mempersiapkan mesin pemroses citra AI...'

  try {
    const multiResult = await processMultiEnhancementOcr(
      rawImageSource,
      (pct, status) => {
        progress.value = pct
        statusText.value = status
      },
      {
        targetWidth: 1400,
        fastAutoScan: false
      }
    )

    const winner = multiResult.winner
    previewUrl.value = winner.dataUrl
    parsedData.value = winner.parsed
    winningMethodLabel.value = winner.methodLabel
    populateEditableForm(winner.parsed)
  } catch (error) {
    console.error('OCR Error:', error)
    alert('Terjadi kendala saat membaca foto e-KTP. Pastikan gambar tidak buram dan coba lagi.')
  } finally {
    isScanning.value = false
  }
}

/**
 * Buat contoh e-KTP presisi tinggi untuk pengujian langsung (fitur demo juri)
 */
const generateSampleKtp = async () => {
  const canvas = document.createElement('canvas')
  canvas.width = 1000
  canvas.height = 630
  const ctx = canvas.getContext('2d')
  if (!ctx) return

  // 1. Background KTP (Gradien Cyan Khas e-KTP Indonesia)
  const bgGrad = ctx.createLinearGradient(0, 0, canvas.width, canvas.height)
  bgGrad.addColorStop(0, '#bae6fd')
  bgGrad.addColorStop(0.5, '#e0f2fe')
  bgGrad.addColorStop(1, '#7dd3fc')
  ctx.fillStyle = bgGrad
  ctx.fillRect(0, 0, canvas.width, canvas.height)

  // 2. Pola Guilloche (Simulasi Batik KTP)
  ctx.strokeStyle = 'rgba(56, 189, 248, 0.25)'
  ctx.lineWidth = 1
  for (let y = 10; y < canvas.height; y += 12) {
    ctx.beginPath()
    ctx.moveTo(0, y)
    for (let x = 0; x < canvas.width; x += 30) {
      ctx.quadraticCurveTo(x + 15, y + (x % 20 === 0 ? 8 : -8), x + 30, y)
    }
    ctx.stroke()
  }

  // 3. Header KTP
  ctx.fillStyle = '#0f172a'
  ctx.textAlign = 'center'
  ctx.font = 'bold 22px Arial'
  ctx.fillText('PROVINSI DKI JAKARTA', canvas.width / 2, 48)
  ctx.font = 'bold 20px Arial'
  ctx.fillText('KOTA JAKARTA SELATAN', canvas.width / 2, 78)

  // 4. Nomor NIK
  ctx.textAlign = 'left'
  ctx.font = 'bold 30px "Courier New", monospace'
  ctx.fillText('NIK : 3174051905060001', 50, 140)

  // 5. Teks Rincian Data KTP
  ctx.font = 'bold 19px Arial'
  const lines = [
    'Nama                 : BUDI SANTOSO',
    'Tempat/Tgl Lahir  : JAKARTA, 19/05/2006',
    'Jenis Kelamin     : LAKI-LAKI        Gol. Darah : B',
    'Alamat               : JL. PASIRHEJO IV NO. 16',
    '    RT/RW            : 007/001',
    '    Kel/Desa         : PASIRHEJO',
    '    Kecamatan        : CIBEUNYING',
    'Agama                : ISLAM',
    'Status Perkawinan : KAWIN',
    'Pekerjaan           : KARYAWAN SWASTA',
    'Kewarganegaraan   : WNI',
    'Berlaku Hingga    : SEUMUR HIDUP'
  ]

  let y = 190
  for (const line of lines) {
    ctx.fillText(line, 50, y)
    y += 34
  }

  // 6. Kotak Foto KTP
  ctx.fillStyle = '#cbd5e1'
  ctx.fillRect(740, 160, 210, 270)
  ctx.strokeStyle = '#94a3b8'
  ctx.lineWidth = 2
  ctx.strokeRect(740, 160, 210, 270)
  ctx.fillStyle = '#64748b'
  ctx.font = '14px Arial'
  ctx.textAlign = 'center'
  ctx.fillText('FOTO PASIEN', 845, 300)

  rawImageSource = canvas
  activeTab.value = 'upload'
  await processAndRunOcr()
}

/**
 * Reset pemindaian saat ini
 */
const resetScan = () => {
  parsedData.value = null
  previewUrl.value = null
  rawImageSource = null
  progress.value = 0
  statusText.value = ''
  editableForm.value = {
    nik: '',
    nama: '',
    tempatLahir: '',
    tglLahir: '',
    tempatTglLahir: '',
    jenisKelamin: 'LAKI-LAKI',
    golDarah: '-',
    alamatJalan: '',
    rt: '',
    rw: '',
    kelDesa: '',
    kecamatan: '',
    confidence: 0,
    isValidNik: false,
  }
}

/**
 * Terapkan data yang telah diverifikasi dan dikoreksi ke formulir registrasi
 */
const applyDataToForm = () => {
  const form = editableForm.value
  const cleanPayload: KtpData = {
    nik: sanitizeNik(form.nik),
    nama: sanitizeCleanText(form.nama),
    tempatLahir: sanitizeCleanText(form.tempatLahir),
    tglLahir: sanitizeCleanText(form.tglLahir),
    tempatTglLahir: `${sanitizeCleanText(form.tempatLahir)} / ${sanitizeCleanText(form.tglLahir)}`,
    jenisKelamin: form.jenisKelamin,
    golDarah: form.golDarah || '-',
    alamatJalan: sanitizeCleanText(form.alamatJalan),
    rt: cleanCharForDigit(form.rt).replace(/\D/g, '').slice(0, 3),
    rw: cleanCharForDigit(form.rw).replace(/\D/g, '').slice(0, 3),
    kelDesa: sanitizeCleanText(form.kelDesa),
    kecamatan: sanitizeCleanText(form.kecamatan),
    confidence: form.confidence,
    isValidNik: form.nik.length === 16,
  }
  emit('apply-ktp', cleanPayload)
  closeModal()
}
</script>

<style scoped>
.modal-fade-enter-active,
.modal-fade-leave-active {
  transition: opacity 0.25s ease, transform 0.25s ease;
}

.modal-fade-enter-from,
.modal-fade-leave-to {
  opacity: 0;
  transform: scale(0.96);
}
</style>
