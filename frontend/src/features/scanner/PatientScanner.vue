<template>
  <main class="min-h-screen bg-[var(--surface-subtle,#f8fafc)] py-8 px-4 sm:px-6 lg:px-8 text-[var(--ink,#0f172a)]">
    <div class="mx-auto max-w-4xl space-y-6">
      <!-- Header & Top Navigation -->
      <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 border-b border-[var(--line,#e2e8f0)] pb-5">
        <div>
          <div class="flex items-center gap-2 mb-1">
            <RouterLink
              to="/staff"
              class="inline-flex items-center gap-1.5 text-sm font-medium text-[var(--ink-muted,#64748b)] hover:text-[var(--primary,#059669)] transition-colors"
            >
              <ChevronLeft :size="16" /> Dashboard Staf
            </RouterLink>
            <span class="text-xs text-slate-300">|</span>
            <RouterLink
              to="/"
              class="inline-flex items-center gap-1.5 text-sm font-medium text-[var(--ink-muted,#64748b)] hover:text-[var(--primary,#059669)] transition-colors"
            >
              Beranda Pasien
            </RouterLink>
          </div>
          <h1 class="text-2xl sm:text-3xl font-bold tracking-tight text-[var(--ink,#0f172a)] flex items-center gap-2.5">
            <span class="inline-flex p-2 rounded-xl bg-emerald-500/10 text-emerald-600">
              <QrCode :size="28" />
            </span>
            Uji Coba & Verifikasi Scan QR Antrean
          </h1>
          <p class="mt-1 text-sm text-[var(--ink-muted,#64748b)]">
            Halaman pengujian pemindaian QR tiket pasien secara live dengan verifikasi database Supabase PostgreSQL.
          </p>
        </div>

        <!-- System Status Pill -->
        <div class="flex flex-wrap items-center gap-2">
          <div class="inline-flex items-center gap-2 px-3 py-1.5 rounded-full text-xs font-semibold bg-emerald-50 text-emerald-700 border border-emerald-200">
            <span class="relative flex h-2 w-2">
              <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75"></span>
              <span class="relative inline-flex rounded-full h-2 w-2 bg-emerald-500"></span>
            </span>
            API v1 Online
          </div>
          <button
            class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-xs font-medium border border-slate-200 bg-white hover:bg-slate-50 transition-colors shadow-sm"
            @click="soundEnabled = !soundEnabled"
            :title="soundEnabled ? 'Matikan Suara Beep' : 'Aktifkan Suara Beep'"
          >
            <Volume2 v-if="soundEnabled" :size="15" class="text-emerald-600" />
            <VolumeX v-else :size="15" class="text-slate-400" />
            {{ soundEnabled ? 'Audio Aktif' : 'Audio Mute' }}
          </button>
        </div>
      </div>

      <!-- Navigation Tabs -->
      <div class="flex flex-wrap gap-2 border-b border-[var(--line,#e2e8f0)] pb-1">
        <button
          v-for="tab in tabs"
          :key="tab.id"
          class="flex items-center gap-2 px-4 py-2.5 text-sm font-semibold rounded-xl transition-all"
          :class="activeTab === tab.id
            ? 'bg-emerald-600 text-white shadow-md shadow-emerald-600/20'
            : 'bg-white text-slate-600 hover:bg-slate-100 border border-slate-200'"
          @click="switchTab(tab.id)"
        >
          <component :is="tab.icon" :size="18" />
          {{ tab.label }}
          <span
            v-if="tab.badge"
            class="ml-1 px-1.5 py-0.5 text-xs rounded-full bg-emerald-100 text-emerald-800"
            :class="{ '!bg-white/20 !text-white': activeTab === tab.id }"
          >
            {{ tab.badge }}
          </span>
        </button>
      </div>

      <!-- MAIN TESTING SECTION -->
      <div class="grid grid-cols-1 lg:grid-cols-12 gap-6">
        <!-- LEFT COLUMN: Testing Method Controls (7 cols) -->
        <div class="lg:col-span-7 space-y-6">
          <!-- TAB 1: Real Live Camera Scanner -->
          <div
            v-show="activeTab === 'camera'"
            class="bg-white rounded-2xl border border-slate-200 p-5 shadow-sm space-y-4"
          >
            <div class="flex items-center justify-between">
              <div>
                <h2 class="text-lg font-bold text-slate-900 flex items-center gap-2">
                  <Camera :size="20" class="text-emerald-600" />
                  Kamera Scanner Langsung
                </h2>
                <p class="text-xs text-slate-500 mt-0.5">
                  Mendukung webcam laptop/PC dan kamera ponsel smartphone secara otomatis.
                </p>
              </div>
              <span
                class="px-2.5 py-1 text-xs font-semibold rounded-full"
                :class="isCameraActive ? 'bg-emerald-100 text-emerald-800' : 'bg-slate-100 text-slate-600'"
              >
                {{ isCameraActive ? 'Kamera Aktif' : 'Kamera Nonaktif' }}
              </span>
            </div>

            <!-- Camera Viewport with Scanner Overlay -->
            <div class="relative overflow-hidden rounded-xl bg-slate-950 aspect-video flex items-center justify-center border-2 border-slate-800">
              <div id="qr-camera-reader" class="w-full h-full object-cover"></div>

              <!-- Animated Laser Reticle overlay when active -->
              <div
                v-if="isCameraActive"
                class="pointer-events-none absolute inset-0 flex items-center justify-center"
              >
                <div class="relative w-64 h-64 border-2 border-emerald-400/70 rounded-2xl shadow-[0_0_20px_rgba(16,185,129,0.3)]">
                  <!-- Reticle Corners -->
                  <div class="absolute -top-1 -left-1 w-6 h-6 border-t-4 border-l-4 border-emerald-400 rounded-tl-lg"></div>
                  <div class="absolute -top-1 -right-1 w-6 h-6 border-t-4 border-r-4 border-emerald-400 rounded-tr-lg"></div>
                  <div class="absolute -bottom-1 -left-1 w-6 h-6 border-b-4 border-l-4 border-emerald-400 rounded-bl-lg"></div>
                  <div class="absolute -bottom-1 -right-1 w-6 h-6 border-b-4 border-r-4 border-emerald-400 rounded-br-lg"></div>

                  <!-- Scanning Laser Line -->
                  <div class="absolute inset-x-0 h-1 bg-gradient-to-r from-transparent via-emerald-400 to-transparent shadow-[0_0_12px_#34d399] animate-scanner-laser"></div>
                </div>
              </div>

              <!-- Idle placeholder when camera is not running -->
              <div
                v-if="!isCameraActive"
                class="absolute inset-0 flex flex-col items-center justify-center p-6 text-center text-slate-400 gap-3"
              >
                <div class="p-4 rounded-2xl bg-slate-900 border border-slate-800 text-slate-400">
                  <CameraOff :size="36" />
                </div>
                <div>
                  <p class="font-medium text-slate-200">Kamera Belum Dinyalakan</p>
                  <p class="text-xs text-slate-400 max-w-sm mt-1">
                    Klik tombol di bawah untuk meminta izin webcam browser Anda dan memulai pemindaian otomatis.
                  </p>
                </div>
              </div>
            </div>

            <!-- Camera Controls & Device Picker -->
            <div class="space-y-3">
              <div v-if="cameraDevices.length > 1" class="flex items-center gap-2">
                <label for="camera-select" class="text-xs font-semibold text-slate-700 whitespace-nowrap">
                  Pilih Perangkat:
                </label>
                <select
                  id="camera-select"
                  v-model="selectedCameraId"
                  :disabled="isCameraActive"
                  class="flex-1 rounded-lg border border-slate-300 py-1.5 px-3 text-xs focus:ring-2 focus:ring-emerald-500 focus:outline-none"
                >
                  <option v-for="dev in cameraDevices" :key="dev.id" :value="dev.id">
                    {{ dev.label || `Kamera ${dev.id}` }}
                  </option>
                </select>
              </div>

              <div class="flex items-center gap-3">
                <button
                  type="button"
                  class="flex-1 inline-flex items-center justify-center gap-2 py-3 px-4 rounded-xl font-semibold text-sm transition-all"
                  :class="isCameraActive
                    ? 'bg-rose-600 hover:bg-rose-700 text-white shadow-lg shadow-rose-600/20'
                    : 'bg-emerald-600 hover:bg-emerald-700 text-white shadow-lg shadow-emerald-600/20'"
                  :disabled="cameraLoading"
                  @click="toggleCamera"
                >
                  <Loader2 v-if="cameraLoading" :size="18" class="animate-spin" />
                  <Square v-else-if="isCameraActive" :size="18" />
                  <Play v-else :size="18" />
                  {{ isCameraActive ? 'Hentikan Kamera' : 'Nyalakan Kamera Scanner' }}
                </button>
              </div>

              <p v-if="cameraError" class="text-xs text-rose-600 bg-rose-50 p-2.5 rounded-lg border border-rose-200 flex items-center gap-1.5">
                <AlertCircle :size="16" class="shrink-0" /> {{ cameraError }}
              </p>
            </div>
          </div>

          <!-- TAB 2: File Upload QR Image Reader -->
          <div
            v-show="activeTab === 'upload'"
            class="bg-white rounded-2xl border border-slate-200 p-5 shadow-sm space-y-4"
          >
            <div>
              <h2 class="text-lg font-bold text-slate-900 flex items-center gap-2">
                <FileImage :size="20" class="text-emerald-600" />
                Upload Gambar Kode QR
              </h2>
              <p class="text-xs text-slate-500 mt-0.5">
                Unggah file gambar struk / tiket e-tiket pasien (PNG, JPG, JPEG, WEBP) untuk mendeteksi QR secara instan.
              </p>
            </div>

            <!-- Drag and Drop Area -->
            <div
              class="border-2 border-dashed rounded-2xl p-8 text-center transition-colors flex flex-col items-center justify-center gap-3 cursor-pointer"
              :class="isDragging ? 'border-emerald-500 bg-emerald-50/50' : 'border-slate-300 hover:border-emerald-400 bg-slate-50/50'"
              @dragover.prevent="isDragging = true"
              @dragleave.prevent="isDragging = false"
              @drop.prevent="handleFileDrop"
              @click="triggerFileInput"
            >
              <input
                ref="fileInputRef"
                type="file"
                accept="image/*"
                class="hidden"
                @change="handleFileChange"
              />

              <div class="p-3.5 rounded-full bg-emerald-100 text-emerald-600">
                <UploadCloud :size="28" />
              </div>

              <div>
                <p class="text-sm font-semibold text-slate-800">
                  Tarik file gambar ke sini atau <span class="text-emerald-600 underline">pilih dari file</span>
                </p>
                <p class="text-xs text-slate-400 mt-1">Mendukung format PNG, JPG, JPEG, WEBP (maks 5MB)</p>
              </div>
            </div>

            <!-- Upload File Status -->
            <div v-if="fileScanning" class="flex items-center justify-center gap-2 py-4 text-emerald-600 text-sm font-medium">
              <Loader2 :size="18" class="animate-spin" /> Membaca dan menganalisis QR dari gambar...
            </div>
            <div v-else-if="fileError" class="text-xs text-rose-600 bg-rose-50 p-3 rounded-lg border border-rose-200 flex items-center gap-2">
              <AlertCircle :size="16" class="shrink-0" /> {{ fileError }}
            </div>
          </div>

          <!-- TAB 3: 1-Click Demo Patient Suite -->
          <div
            v-show="activeTab === 'demo'"
            class="bg-white rounded-2xl border border-slate-200 p-5 shadow-sm space-y-4"
          >
            <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3">
              <div>
                <h2 class="text-lg font-bold text-slate-900 flex items-center gap-2">
                  <Sparkles :size="20" class="text-amber-500" />
                  Uji Coba Pasien Demo (1-Click)
                </h2>
                <p class="text-xs text-slate-500 mt-0.5">
                  Data antrean realistis yang tersimpan di Supabase PostgreSQL. Klik tombol uji coba untuk menguji respon sistem.
                </p>
              </div>

              <button
                type="button"
                class="inline-flex items-center gap-1.5 px-3 py-2 rounded-xl text-xs font-semibold bg-slate-100 hover:bg-slate-200 text-slate-700 transition-colors self-start sm:self-auto border border-slate-300"
                :disabled="resettingLoading"
                @click="resetAllDemoQueues"
              >
                <RefreshCw :size="14" :class="{ 'animate-spin': resettingLoading }" />
                Reset Semua ke Status 'Waiting'
              </button>
            </div>

            <!-- Demo Patients Cards Grid -->
            <div class="grid grid-cols-1 sm:grid-cols-2 gap-3.5">
              <div
                v-for="item in demoQueues"
                :key="item.qr_token"
                class="border rounded-xl p-3.5 flex flex-col justify-between transition-all hover:shadow-md bg-white"
                :class="{
                  'border-emerald-300 bg-emerald-50/30': item.status === 'checked_in',
                  'border-blue-300 bg-blue-50/30': item.status === 'waiting',
                  'border-amber-300 bg-amber-50/30': item.status === 'called',
                  'border-purple-300 bg-purple-50/30': item.status === 'in_service'
                }"
              >
                <div class="flex items-start gap-3">
                  <!-- Generated QR Code Preview -->
                  <div class="w-16 h-16 shrink-0 rounded-lg bg-white p-1 border border-slate-200 shadow-sm flex items-center justify-center">
                    <img
                      v-if="qrCodeDataUrls[item.qr_token]"
                      :src="qrCodeDataUrls[item.qr_token]"
                      :alt="`QR ${item.patient_name}`"
                      class="w-full h-full object-contain"
                    />
                    <QrCode v-else :size="28" class="text-slate-400" />
                  </div>

                  <!-- Patient Info -->
                  <div class="min-w-0 flex-1">
                    <div class="flex items-center justify-between gap-1">
                      <span class="font-bold text-sm text-slate-900 truncate">{{ item.patient_name }}</span>
                      <span class="text-xs font-black px-1.5 py-0.5 rounded bg-slate-900 text-white">
                        {{ item.queue_code }}
                      </span>
                    </div>
                    <p class="text-xs text-slate-500 mt-0.5">{{ item.polyclinic_name }}</p>

                    <!-- Status Badge -->
                    <div class="mt-2 flex items-center gap-1.5">
                      <span
                        class="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-[11px] font-semibold"
                        :class="getStatusBadgeClass(item.status)"
                      >
                        <span class="w-1.5 h-1.5 rounded-full bg-current"></span>
                        {{ formatStatus(item.status) }}
                      </span>
                    </div>
                  </div>
                </div>

                <!-- Actions for each card -->
                <div class="mt-3.5 pt-2.5 border-t border-slate-100 flex items-center gap-2">
                  <button
                    type="button"
                    class="flex-1 py-1.5 px-2.5 rounded-lg text-xs font-semibold bg-emerald-600 hover:bg-emerald-700 text-white transition-colors flex items-center justify-center gap-1 shadow-sm"
                    :disabled="verifyingToken === item.qr_token"
                    @click="handleCheckIn(item.qr_token, 'Simulasi ' + item.patient_name)"
                  >
                    <Loader2 v-if="verifyingToken === item.qr_token" :size="13" class="animate-spin" />
                    <Check v-else :size="13" />
                    Uji Scan
                  </button>

                  <button
                    type="button"
                    class="py-1.5 px-2 rounded-lg text-xs font-medium text-slate-600 hover:text-slate-900 hover:bg-slate-100 border border-slate-200 transition-colors"
                    title="Kembalikan status pasien ini ke 'waiting'"
                    @click="resetSingleQueue(item.qr_token)"
                  >
                    <RotateCcw :size="13" />
                  </button>
                </div>
              </div>
            </div>
          </div>

          <!-- TAB 4: Manual Code Input -->
          <div
            v-show="activeTab === 'manual'"
            class="bg-white rounded-2xl border border-slate-200 p-5 shadow-sm space-y-4"
          >
            <div>
              <h2 class="text-lg font-bold text-slate-900 flex items-center gap-2">
                <Keyboard :size="20" class="text-emerald-600" />
                Input Manual Token QR
              </h2>
              <p class="text-xs text-slate-500 mt-0.5">
                Ketikkan kode token QR string atau paste URL checkin untuk pengujian langsung.
              </p>
            </div>

            <form @submit.prevent="handleManualSubmit" class="space-y-3">
              <div>
                <label for="manual-token" class="block text-xs font-semibold text-slate-700 mb-1">
                  String Token QR:
                </label>
                <div class="relative">
                  <input
                    id="manual-token"
                    v-model="manualToken"
                    type="text"
                    required
                    placeholder="Contoh: qr-ahmad-fauzi-004 atau token UUID"
                    class="w-full rounded-xl border border-slate-300 py-2.5 pl-3.5 pr-10 text-sm focus:ring-2 focus:ring-emerald-500 focus:outline-none"
                  />
                  <button
                    v-if="manualToken"
                    type="button"
                    class="absolute right-3 top-1/2 -translate-y-1/2 text-slate-400 hover:text-slate-600"
                    @click="manualToken = ''"
                  >
                    <X :size="16" />
                  </button>
                </div>
              </div>

              <!-- Quick Chips -->
              <div class="flex flex-wrap items-center gap-1.5">
                <span class="text-xs text-slate-400">Pilih cepat:</span>
                <button
                  v-for="token in quickChips"
                  :key="token.code"
                  type="button"
                  class="text-[11px] font-mono px-2 py-1 rounded bg-slate-100 hover:bg-slate-200 text-slate-700 transition-colors border border-slate-200"
                  @click="manualToken = token.code"
                >
                  {{ token.label }}
                </button>
              </div>

              <button
                type="submit"
                class="w-full py-2.5 px-4 rounded-xl font-semibold text-sm bg-emerald-600 hover:bg-emerald-700 text-white transition-all shadow-md shadow-emerald-600/20 flex items-center justify-center gap-2"
                :disabled="verifyingToken === manualToken || !manualToken.trim()"
              >
                <Loader2 v-if="verifyingToken === manualToken" :size="16" class="animate-spin" />
                <CheckCircle2 v-else :size="16" />
                Verifikasi Check-In Token
              </button>
            </form>
          </div>
        </div>

        <!-- RIGHT COLUMN: Verification Outcome & Live Inspector (5 cols) -->
        <div class="lg:col-span-5 space-y-6">
          <!-- RESULT CARD: Verification Outcome -->
          <div class="bg-white rounded-2xl border border-slate-200 p-5 shadow-sm space-y-4">
            <div class="flex items-center justify-between border-b border-slate-100 pb-3">
              <h3 class="font-bold text-slate-900 flex items-center gap-2 text-base">
                <Activity :size="18" class="text-emerald-600" />
                Hasil Verifikasi Scan
              </h3>
              <span v-if="lastResult" class="text-[11px] text-slate-400">
                {{ lastResultTime }}
              </span>
            </div>

            <!-- 1. State: Waiting for Scan -->
            <div
              v-if="!lastResult && !verifyingToken"
              class="py-12 px-4 text-center text-slate-400 flex flex-col items-center justify-center gap-2"
            >
              <div class="p-3.5 rounded-full bg-slate-100 text-slate-400">
                <ScanLine :size="32" />
              </div>
              <p class="font-medium text-slate-700 text-sm">Belum Ada QR yang Dipindai</p>
              <p class="text-xs text-slate-400 max-w-xs">
                Arahkan kamera ke QR code, upload file, atau klik salah satu tombol uji coba di sebelah kiri.
              </p>
            </div>

            <!-- 2. State: Scanning / Verifying in Progress -->
            <div
              v-else-if="verifyingToken"
              class="py-12 px-4 text-center flex flex-col items-center justify-center gap-3 text-emerald-600"
            >
              <div class="relative">
                <div class="w-14 h-14 rounded-full border-4 border-emerald-200 border-t-emerald-600 animate-spin"></div>
                <div class="absolute inset-0 flex items-center justify-center">
                  <QrCode :size="20" class="text-emerald-600" />
                </div>
              </div>
              <p class="font-bold text-slate-800 text-sm">Menghubungi Server SangPalito...</p>
              <p class="text-xs text-slate-500">Memvalidasi token QR terhadap database PostgreSQL...</p>
            </div>

            <!-- 3. State: SUCCESS -->
            <div
              v-else-if="lastResult && lastResult.success"
              class="rounded-xl border border-emerald-200 bg-gradient-to-br from-emerald-50 to-teal-50/30 p-5 space-y-4"
            >
              <div class="flex items-start gap-3">
                <div class="p-2 rounded-xl bg-emerald-600 text-white shadow-md shadow-emerald-600/30">
                  <CheckCircle2 :size="26" />
                </div>
                <div>
                  <div class="flex items-center gap-2">
                    <span class="px-2 py-0.5 rounded text-xs font-bold bg-emerald-600 text-white">
                      BERHASIL CHECK-IN
                    </span>
                    <span class="text-xs text-emerald-800 font-medium">Valid & Terverifikasi</span>
                  </div>
                  <h4 class="text-lg font-bold text-slate-900 mt-1">
                    {{ lastResultPatientName || 'Pasien AntriSehat' }}
                  </h4>
                </div>
              </div>

              <!-- Ticket Code & Details Badge -->
              <div class="p-4 rounded-xl bg-white border border-emerald-200/80 shadow-sm space-y-3">
                <div class="flex items-center justify-between">
                  <div>
                    <p class="text-xs text-slate-500">Nomor Antrean</p>
                    <p class="text-3xl font-black text-emerald-700 tracking-tight">
                      {{ lastResult.data?.queue_code || 'A-001' }}
                    </p>
                  </div>
                  <div class="text-right">
                    <p class="text-xs text-slate-500">Status Terbaru</p>
                    <span class="inline-flex items-center gap-1 px-2.5 py-1 rounded-full text-xs font-bold bg-emerald-100 text-emerald-800">
                      <span class="w-1.5 h-1.5 rounded-full bg-emerald-600"></span>
                      checked_in
                    </span>
                  </div>
                </div>

                <div class="grid grid-cols-2 gap-2 pt-2 border-t border-slate-100 text-xs">
                  <div>
                    <span class="text-slate-400">Poliklinik:</span>
                    <p class="font-semibold text-slate-800">{{ lastResultPolyclinic || 'Poli Umum' }}</p>
                  </div>
                  <div>
                    <span class="text-slate-400">Waktu Check-In:</span>
                    <p class="font-semibold text-slate-800">{{ lastResultCheckInTime }}</p>
                  </div>
                </div>

                <!-- Wait Time Estimate Info (EMA) -->
                <div v-if="queueEstimate" class="mt-2 p-2.5 rounded-lg bg-emerald-50 border border-emerald-100 text-xs text-emerald-900 space-y-1">
                  <div class="flex items-center justify-between font-semibold">
                    <span>Estimasi Menunggu:</span>
                    <span class="text-emerald-700 font-bold">~{{ queueEstimate.estimate_minutes }} Menit</span>
                  </div>
                  <p class="text-[11px] text-emerald-700">
                    {{ queueEstimate.ahead_count }} pasien menunggu di depan antrean ini.
                  </p>
                </div>
              </div>

              <div class="flex gap-2">
                <RouterLink
                  to="/staff"
                  class="flex-1 py-2 px-3 rounded-lg text-xs font-semibold bg-emerald-600 hover:bg-emerald-700 text-white transition-colors flex items-center justify-center gap-1.5"
                >
                  <Users :size="14" /> Buka di Dashboard Staf
                </RouterLink>
                <button
                  type="button"
                  class="py-2 px-3 rounded-lg text-xs font-semibold bg-white hover:bg-slate-50 text-slate-700 border border-slate-200 transition-colors"
                  @click="resetLastResult"
                >
                  Uji Lagi
                </button>
              </div>
            </div>

            <!-- 4. State: FAILED / ERROR -->
            <div
              v-else-if="lastResult && !lastResult.success"
              class="rounded-xl border border-rose-200 bg-gradient-to-br from-rose-50 to-orange-50/30 p-5 space-y-4"
            >
              <div class="flex items-start gap-3">
                <div class="p-2 rounded-xl bg-rose-600 text-white shadow-md shadow-rose-600/30">
                  <XCircle :size="26" />
                </div>
                <div>
                  <span class="px-2 py-0.5 rounded text-xs font-bold bg-rose-600 text-white">
                    CHECK-IN GAGAL
                  </span>
                  <h4 class="text-base font-bold text-slate-900 mt-1">
                    {{ lastResult.message || 'Verifikasi Tidak Berhasil' }}
                  </h4>
                </div>
              </div>

              <div class="p-3.5 rounded-xl bg-white border border-rose-200/80 text-xs space-y-2">
                <div class="flex items-center justify-between text-slate-600">
                  <span>Kode Error:</span>
                  <span class="font-mono font-bold text-rose-700">{{ lastResult.code || 'BAD_REQUEST' }}</span>
                </div>
                <div class="flex items-center justify-between text-slate-600">
                  <span>Token yang Diuji:</span>
                  <span class="font-mono text-slate-900 truncate max-w-[200px]">{{ lastScannedToken }}</span>
                </div>
                <div class="pt-2 border-t border-slate-100 text-slate-500 text-[11px] leading-relaxed">
                  💡 <strong>Solusi:</strong> Jika error menyebut status antrean bukan <em>waiting</em>, klik tombol <strong>Reset</strong> pada tab Simulasi untuk mengembalikan antrean ke status <em>waiting</em> agar bisa di-check-in kembali.
                </div>
              </div>

              <button
                type="button"
                class="w-full py-2 px-3 rounded-lg text-xs font-semibold bg-white hover:bg-slate-50 text-slate-700 border border-slate-200 transition-colors"
                @click="resetLastResult"
              >
                Coba Token Lain
              </button>
            </div>
          </div>

          <!-- RAW JSON INSPECTOR -->
          <div class="bg-white rounded-2xl border border-slate-200 p-5 shadow-sm space-y-3">
            <div class="flex items-center justify-between cursor-pointer" @click="showJsonInspector = !showJsonInspector">
              <h3 class="font-bold text-slate-900 flex items-center gap-2 text-sm">
                <Code :size="16" class="text-slate-500" />
                Developer JSON Inspector
              </h3>
              <span class="text-xs text-slate-400 flex items-center gap-1">
                {{ showJsonInspector ? 'Sembunyikan' : 'Tampilkan' }}
                <ChevronDown :size="14" :class="{ 'rotate-180': showJsonInspector }" />
              </span>
            </div>

            <div v-show="showJsonInspector" class="space-y-2">
              <div class="text-[11px] text-slate-500 flex justify-between">
                <span>Payload Respon HTTP Backend:</span>
                <span class="font-mono">POST /api/v1/public/queues/checkin</span>
              </div>
              <pre class="bg-slate-950 text-emerald-400 p-3.5 rounded-xl text-[11px] font-mono overflow-x-auto max-h-56 leading-relaxed border border-slate-800">{{ formattedJsonResult }}</pre>
            </div>
          </div>

          <!-- SCAN SESSION AUDIT LOG -->
          <div class="bg-white rounded-2xl border border-slate-200 p-5 shadow-sm space-y-3">
            <div class="flex items-center justify-between border-b border-slate-100 pb-2">
              <h3 class="font-bold text-slate-900 flex items-center gap-2 text-sm">
                <History :size="16" class="text-slate-500" />
                Riwayat Uji Coba Sesi Ini
              </h3>
              <span class="text-xs text-slate-400">{{ scanHistory.length }} riwayat</span>
            </div>

            <div v-if="scanHistory.length === 0" class="text-xs text-slate-400 text-center py-4">
              Belum ada riwayat uji coba pada sesi ini.
            </div>
            <div v-else class="space-y-2 max-h-48 overflow-y-auto pr-1">
              <div
                v-for="(hist, idx) in scanHistory"
                :key="idx"
                class="flex items-center justify-between p-2 rounded-lg text-xs border"
                :class="hist.success ? 'bg-emerald-50/50 border-emerald-100' : 'bg-rose-50/50 border-rose-100'"
              >
                <div class="min-w-0 flex-1 pr-2">
                  <div class="flex items-center gap-1.5">
                    <span
                      class="px-1.5 py-0.2 rounded text-[10px] font-bold"
                      :class="hist.success ? 'bg-emerald-600 text-white' : 'bg-rose-600 text-white'"
                    >
                      {{ hist.success ? 'SUKSES' : 'GAGAL' }}
                    </span>
                    <span class="font-medium text-slate-800 truncate">{{ hist.source }}</span>
                  </div>
                  <p class="font-mono text-[10px] text-slate-500 truncate mt-0.5">{{ hist.token }}</p>
                </div>
                <div class="text-right text-[10px] text-slate-400 shrink-0">
                  <p>{{ hist.time }}</p>
                  <p>{{ hist.durationMs }}ms</p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </main>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from "vue";
import { useRouter } from "vue-router";
import { Html5Qrcode } from "html5-qrcode";
import QRCode from "qrcode";
import {
  QrCode, Camera, CameraOff, UploadCloud, FileImage, Sparkles,
  Keyboard, Activity, CheckCircle2, XCircle, Code, History,
  Volume2, VolumeX, ChevronLeft, ChevronDown, Check, X,
  RotateCcw, RefreshCw, Play, Square, Loader2, Users,
  AlertCircle, ScanLine
} from "lucide-vue-next";
import { api } from "../../services/api";

const router = useRouter();

// Tabs Configuration
const tabs = [
  { id: "camera", label: "Kamera Scanner", icon: Camera },
  { id: "demo", label: "Simulasi Pasien Demo", icon: Sparkles, badge: "Populer" },
  { id: "upload", label: "Upload File QR", icon: FileImage },
  { id: "manual", label: "Input Manual", icon: Keyboard },
];
const activeTab = ref("demo");

// Camera Scanner State
const isCameraActive = ref(false);
const cameraLoading = ref(false);
const cameraError = ref("");
const selectedCameraId = ref("");
const cameraDevices = ref<Array<{ id: string; label: string }>>([]);
let html5QrCode: Html5Qrcode | null = null;

// File Upload State
const fileInputRef = ref<HTMLInputElement | null>(null);
const isDragging = ref(false);
const fileScanning = ref(false);
const fileError = ref("");

// Demo Patients State
interface DemoQueueItem {
  qr_token: string;
  patient_name: string;
  queue_code: string;
  queue_number: number;
  polyclinic_name: string;
  status: string;
  checked_in_at?: string;
}
const demoQueues = ref<DemoQueueItem[]>([
  { qr_token: "qr-ahmad-fauzi-004", patient_name: "Ahmad Fauzi", queue_code: "A-004", queue_number: 4, polyclinic_name: "Poli Umum", status: "waiting" },
  { qr_token: "qr-dewi-lestari-003", patient_name: "Dewi Lestari", queue_code: "A-003", queue_number: 3, polyclinic_name: "Poli Umum", status: "checked_in" },
  { qr_token: "qr-budi-santoso-002", patient_name: "Budi Santoso", queue_code: "A-002", queue_number: 2, polyclinic_name: "Poli Umum", status: "called" },
  { qr_token: "qr-siti-aminah-001", patient_name: "Siti Aminah", queue_code: "A-001", queue_number: 1, polyclinic_name: "Poli Umum", status: "in_service" },
]);
const qrCodeDataUrls = ref<Record<string, string>>({});
const resettingLoading = ref(false);

// Manual Input State
const manualToken = ref("");
const quickChips = [
  { label: "Ahmad Fauzi", code: "qr-ahmad-fauzi-004" },
  { label: "Dewi Lestari", code: "qr-dewi-lestari-003" },
  { label: "Budi Santoso", code: "qr-budi-santoso-002" },
  { label: "Siti Aminah", code: "qr-siti-aminah-001" },
  { label: "Token Salah (Test 404)", code: "qr-invalid-token-test-404" },
];

// Result & Audio Feedback State
const verifyingToken = ref("");
const lastScannedToken = ref("");
const lastResult = ref<any>(null);
const lastResultTime = ref("");
const lastResultPatientName = ref("");
const lastResultPolyclinic = ref("");
const lastResultCheckInTime = ref("");
const queueEstimate = ref<any>(null);
const soundEnabled = ref(true);
const showJsonInspector = ref(true);

// Session History Log
interface ScanHistoryItem {
  source: string;
  token: string;
  success: boolean;
  time: string;
  durationMs: number;
}
const scanHistory = ref<ScanHistoryItem[]>([]);

const formattedJsonResult = computed(() => {
  if (!lastResult.value) return "// Belum ada respon scan.";
  return JSON.stringify(lastResult.value, null, 2);
});

// Sound Generator using Web Audio API
function playSound(type: "success" | "error") {
  if (!soundEnabled.value) return;
  try {
    const ctx = new (window.AudioContext || (window as any).webkitAudioContext)();
    const osc = ctx.createOscillator();
    const gain = ctx.createGain();
    osc.connect(gain);
    gain.connect(ctx.destination);

    if (type === "success") {
      // Pleasant high-pitch double chime
      osc.type = "sine";
      osc.frequency.setValueAtTime(587.33, ctx.currentTime); // D5
      osc.frequency.setValueAtTime(880, ctx.currentTime + 0.08); // A5
      gain.gain.setValueAtTime(0.2, ctx.currentTime);
      gain.gain.exponentialRampToValueAtTime(0.01, ctx.currentTime + 0.3);
      osc.start(ctx.currentTime);
      osc.stop(ctx.currentTime + 0.3);
    } else {
      // Low alert buzz
      osc.type = "sawtooth";
      osc.frequency.setValueAtTime(220, ctx.currentTime);
      osc.frequency.setValueAtTime(160, ctx.currentTime + 0.1);
      gain.gain.setValueAtTime(0.25, ctx.currentTime);
      gain.gain.exponentialRampToValueAtTime(0.01, ctx.currentTime + 0.25);
      osc.start(ctx.currentTime);
      osc.stop(ctx.currentTime + 0.25);
    }
  } catch {
    // Audio context may be restricted before user gesture
  }
}

// Formatters
function formatStatus(status: string) {
  switch (status) {
    case "waiting": return "Menunggu (waiting)";
    case "checked_in": return "Sudah Check-In";
    case "called": return "Dipanggil";
    case "in_service": return "Sedang Dilayani";
    case "done": return "Selesai";
    default: return status;
  }
}

function getStatusBadgeClass(status: string) {
  switch (status) {
    case "waiting": return "bg-blue-100 text-blue-800";
    case "checked_in": return "bg-emerald-100 text-emerald-800";
    case "called": return "bg-amber-100 text-amber-800";
    case "in_service": return "bg-purple-100 text-purple-800";
    default: return "bg-slate-100 text-slate-700";
  }
}

// ── Check-In Verification Logic ──────────────────────────────
async function handleCheckIn(rawToken: string, source: string) {
  let token = rawToken.trim();
  // Support checkin url extraction
  try {
    const url = new URL(token);
    token = url.searchParams.get("checkin") || url.searchParams.get("token") || token;
  } catch {
    // raw string
  }

  verifyingToken.value = token;
  lastScannedToken.value = token;
  const startTime = performance.now();

  try {
    const res = await api.checkInByQR(token);
    const duration = Math.round(performance.now() - startTime);

    lastResult.value = res;
    lastResultTime.value = new Date().toLocaleTimeString("id-ID");

    if (res && res.success) {
      playSound("success");
      const matched = demoQueues.value.find(q => q.qr_token === token);
      lastResultPatientName.value = matched?.patient_name || "Pasien Terverifikasi";
      lastResultPolyclinic.value = matched?.polyclinic_name || "Poli Umum";
      lastResultCheckInTime.value = new Date().toLocaleTimeString("id-ID", { hour: '2-digit', minute: '2-digit' });

      // Fetch queue estimate
      if (res.data?.id) {
        api.getQueueEstimate(res.data.id).then(est => {
          queueEstimate.value = est;
        });
      }

      scanHistory.value.unshift({
        source,
        token,
        success: true,
        time: lastResultTime.value,
        durationMs: duration,
      });

      // Refresh demo queues list
      loadDemoQueues();
    } else {
      playSound("error");
      scanHistory.value.unshift({
        source,
        token,
        success: false,
        time: lastResultTime.value,
        durationMs: duration,
      });
    }
  } catch (err: any) {
    playSound("error");
    const duration = Math.round(performance.now() - startTime);
    lastResult.value = {
      success: false,
      code: "NETWORK_ERROR",
      message: err.message || "Gagal menghubungi server",
    };
    scanHistory.value.unshift({
      source,
      token,
      success: false,
      time: new Date().toLocaleTimeString("id-ID"),
      durationMs: duration,
    });
  } finally {
    verifyingToken.value = "";
  }
}

function resetLastResult() {
  lastResult.value = null;
  queueEstimate.value = null;
}

// ── Tab 1: Live Camera Scanner ──────────────────────────────
async function initCamera() {
  try {
    cameraLoading.value = true;
    cameraError.value = "";

    const devices = await Html5Qrcode.getCameras();
    if (devices && devices.length) {
      cameraDevices.value = devices;
      if (!selectedCameraId.value) {
        selectedCameraId.value = devices[0].id;
      }
    }

    if (!html5QrCode) {
      html5QrCode = new Html5Qrcode("qr-camera-reader");
    }

    const cameraId = selectedCameraId.value || { facingMode: "environment" };

    await html5QrCode.start(
      cameraId,
      {
        fps: 10,
        qrbox: { width: 250, height: 250 },
      },
      (decodedText) => {
        if (!verifyingToken.value) {
          handleCheckIn(decodedText, "Kamera Scanner");
        }
      },
      () => {
        // Ignored frame error
      }
    );

    isCameraActive.value = true;
  } catch (err: any) {
    cameraError.value = err?.message || "Tidak dapat mengakses kamera. Pastikan izin kamera telah diberikan.";
    isCameraActive.value = false;
  } finally {
    cameraLoading.value = false;
  }
}

async function stopCamera() {
  if (html5QrCode && isCameraActive.value) {
    try {
      await html5QrCode.stop();
    } catch {
      // ignore
    }
    isCameraActive.value = false;
  }
}

async function toggleCamera() {
  if (isCameraActive.value) {
    await stopCamera();
  } else {
    await initCamera();
  }
}

// ── Tab 2: Upload File QR ────────────────────────────────────
function triggerFileInput() {
  fileInputRef.value?.click();
}

async function handleFileChange(event: Event) {
  const target = event.target as HTMLInputElement;
  const file = target.files?.[0];
  if (file) {
    await processImageFile(file);
  }
}

async function handleFileDrop(event: DragEvent) {
  isDragging.value = false;
  const file = event.dataTransfer?.files?.[0];
  if (file && file.type.startsWith("image/")) {
    await processImageFile(file);
  }
}

async function processImageFile(file: File) {
  fileScanning.value = true;
  fileError.value = "";

  try {
    if (!html5QrCode) {
      html5QrCode = new Html5Qrcode("qr-camera-reader");
    }
    const decodedText = await html5QrCode.scanFile(file, true);
    await handleCheckIn(decodedText, "Upload File (" + file.name + ")");
  } catch (err: any) {
    fileError.value = "Tidak menemukan kode QR yang jelas dalam gambar tersebut. Coba gunakan gambar lain.";
    playSound("error");
  } finally {
    fileScanning.value = false;
  }
}

// ── Tab 3: Demo Patients ────────────────────────────────────
async function loadDemoQueues() {
  const list = await api.getDemoQueues();
  if (list && list.length) {
    demoQueues.value = list;
    // Generate QR images
    for (const q of list) {
      if (!qrCodeDataUrls.value[q.qr_token]) {
        try {
          qrCodeDataUrls.value[q.qr_token] = await QRCode.toDataURL(q.qr_token, {
            width: 140,
            margin: 1,
            color: { dark: "#0f172a", light: "#ffffff" },
          });
        } catch (e) {
          console.error(e);
        }
      }
    }
  }
}

async function resetAllDemoQueues() {
  resettingLoading.value = true;
  await api.resetDemoQueue("all");
  await loadDemoQueues();
  resettingLoading.value = false;
}

async function resetSingleQueue(token: string) {
  await api.resetDemoQueue(token);
  await loadDemoQueues();
}

// ── Tab 4: Manual Input ─────────────────────────────────────
function handleManualSubmit() {
  if (manualToken.value.trim()) {
    handleCheckIn(manualToken.value, "Input Manual");
  }
}

function switchTab(tabId: string) {
  if (activeTab.value === "camera" && tabId !== "camera") {
    stopCamera();
  }
  activeTab.value = tabId;
  if (tabId === "camera" && !isCameraActive.value) {
    initCamera();
  }
}

onMounted(async () => {
  await loadDemoQueues();
});

onUnmounted(() => {
  stopCamera();
});
</script>

<style scoped>
@keyframes laserSweep {
  0% { top: 0%; opacity: 0.8; }
  50% { top: 96%; opacity: 1; }
  100% { top: 0%; opacity: 0.8; }
}

.animate-scanner-laser {
  animation: laserSweep 2s ease-in-out infinite;
}
</style>
