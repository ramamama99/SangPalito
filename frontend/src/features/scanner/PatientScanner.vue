<template>
  <main
    ref="scannerContainerRef"
    class="min-h-screen bg-[var(--surface-subtle,#f8fafc)] text-slate-900 flex flex-col font-sans transition-colors duration-200"
    :class="{ 'p-4 sm:p-6 lg:p-8': !isFullscreen, 'p-6 sm:p-10': isFullscreen }"
  >
    <!-- TOP KIOSK HEADER -->
    <header class="w-full max-w-6xl mx-auto mb-6 bg-white rounded-2xl border border-slate-200 p-4 sm:p-5 shadow-sm flex flex-col md:flex-row md:items-center md:justify-between gap-4">
      <!-- Hospital Branding -->
      <div class="flex items-center gap-3.5">
        <div class="w-12 h-12 rounded-2xl bg-emerald-50 text-emerald-600 border border-emerald-200 flex items-center justify-center shrink-0 shadow-xs">
          <Activity :size="26" class="text-emerald-600" />
        </div>
        <div>
          <div class="flex items-center gap-2">
            <span class="text-[11px] font-bold tracking-wider uppercase px-2 py-0.5 rounded-full bg-emerald-100 text-emerald-800">
              Anjungan Check-In Mandiri
            </span>
            <span class="text-xs text-slate-400 font-mono">SangPalito</span>
          </div>
          <h1 class="text-xl sm:text-2xl font-black text-slate-900 tracking-tight mt-0.5">
            Puskesmas Sehat Bersama
          </h1>
          <p class="text-xs text-slate-500 hidden sm:block">
            Pindai tiket antrean pasien menggunakan mesin scanner barcode fisik atau kamera kios.
          </p>
        </div>
      </div>

      <!-- Real-time Digital Clock -->
      <div class="flex items-center justify-between md:justify-center px-4 py-2.5 rounded-xl bg-slate-50 border border-slate-200 text-center">
        <div>
          <div class="text-2xl sm:text-3xl font-mono font-bold tracking-wider text-emerald-700">
            {{ currentTimeString || '00:00:00 WIB' }}
          </div>
          <div class="text-[11px] font-medium text-slate-500 uppercase tracking-wider">
            {{ currentDateString || 'Memuat Waktu...' }}
          </div>
        </div>
      </div>

      <!-- Controls & Navigation -->
      <div class="flex flex-wrap items-center gap-2.5">
        <!-- Audio Toggle Button -->
        <button
          type="button"
          class="p-2 sm:px-3 sm:py-2 rounded-xl text-xs font-semibold border border-slate-200 bg-white hover:bg-slate-50 text-slate-700 transition-colors shadow-xs flex items-center gap-1.5"
          @click="soundEnabled = !soundEnabled"
          :title="soundEnabled ? 'Matikan Suara Konfirmasi' : 'Aktifkan Suara Konfirmasi'"
        >
          <Volume2 v-if="soundEnabled" :size="16" class="text-emerald-600" />
          <VolumeX v-else :size="16" class="text-slate-400" />
          <span class="hidden sm:inline">{{ soundEnabled ? 'Audio Aktif' : 'Mute' }}</span>
        </button>

        <!-- Fullscreen Kiosk Mode Toggle -->
        <button
          type="button"
          class="p-2 sm:px-3 sm:py-2 rounded-xl text-xs font-semibold border transition-all flex items-center gap-1.5 shadow-xs"
          :class="isFullscreen
            ? 'bg-emerald-600 border-emerald-600 text-white shadow-emerald-600/20'
            : 'bg-white border-slate-200 hover:bg-slate-50 text-slate-700'"
          @click="toggleFullscreen"
          title="Mode Layar Penuh"
        >
          <Minimize v-if="isFullscreen" :size="16" />
          <Maximize v-else :size="16" />
          <span class="hidden sm:inline">{{ isFullscreen ? 'Keluar Fullscreen' : 'Layar Penuh' }}</span>
        </button>

        <!-- Back to Staff Dashboard -->
        <RouterLink
          to="/staff"
          class="p-2 sm:px-3 sm:py-2 rounded-xl text-xs font-semibold bg-white hover:bg-slate-50 border border-slate-200 text-slate-700 transition-colors flex items-center gap-1 shadow-xs"
        >
          <ChevronLeft :size="15" />
          <span class="hidden sm:inline">Staf</span>
        </RouterLink>
      </div>
    </header>

    <!-- COOLDOWN & SENSOR ACTIVE STATUS BANNER -->
    <section class="w-full max-w-6xl mx-auto mb-6">
      <div
        class="relative overflow-hidden rounded-2xl border transition-all duration-300 p-4 shadow-xs flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3"
        :class="isCooldown
          ? 'bg-amber-50 border-amber-300 text-amber-950'
          : 'bg-emerald-50 border-emerald-200 text-emerald-950'"
      >
        <!-- Animated Cooldown Progress Bar -->
        <div
          v-if="isCooldown"
          class="absolute bottom-0 left-0 h-1.5 bg-gradient-to-r from-amber-400 via-orange-400 to-amber-500 transition-all duration-75 ease-linear"
          :style="{ width: `${cooldownProgress}%` }"
        ></div>

        <div class="flex items-center gap-3.5 min-w-0">
          <div
            class="w-10 h-10 rounded-xl flex items-center justify-center shrink-0 shadow-xs transition-colors"
            :class="isCooldown
              ? 'bg-amber-100 text-amber-700 border border-amber-300 animate-pulse'
              : 'bg-emerald-100 text-emerald-700 border border-emerald-300'"
          >
            <Hourglass v-if="isCooldown" :size="20" class="animate-spin" />
            <Zap v-else :size="20" class="text-emerald-600" />
          </div>

          <div class="min-w-0">
            <div class="flex items-center gap-2">
              <span
                class="text-xs font-black uppercase tracking-wider px-2 py-0.5 rounded-md"
                :class="isCooldown
                  ? 'bg-amber-200 text-amber-900 border border-amber-300'
                  : 'bg-emerald-200 text-emerald-900 border border-emerald-300'"
              >
                {{ isCooldown ? 'Jeda Pemindaian Aktif' : 'Sensor Siap Memindai' }}
              </span>
              <span class="text-xs text-slate-500 font-mono hidden sm:inline">
                Jeda 3 Detik Anti-Duplikasi
              </span>
            </div>

            <p class="text-sm font-semibold text-slate-800 truncate mt-0.5">
              <template v-if="isCooldown">
                Menyiapkan mesin untuk pasien berikutnya... Mohon tunggu
                <span class="font-mono font-bold text-amber-800 text-base">({{ cooldownRemaining }}s)</span>
              </template>
              <template v-else>
                Silakan tempelkan QR tiket ke mesin scan fisik atau hadapkan ke kamera di bawah.
              </template>
            </p>
          </div>
        </div>

        <!-- Sensor Status Indicator Pill -->
        <div class="flex items-center gap-2 shrink-0 self-end sm:self-auto">
          <div
            class="px-3 py-1.5 rounded-xl text-xs font-mono font-bold flex items-center gap-2 border bg-white shadow-xs"
            :class="isCooldown
              ? 'border-amber-300 text-amber-800'
              : 'border-emerald-300 text-emerald-800'"
          >
            <span class="w-2.5 h-2.5 rounded-full" :class="isCooldown ? 'bg-amber-500 animate-ping' : 'bg-emerald-500'"></span>
            {{ isCooldown ? `JEDA (${cooldownRemaining}s)` : 'SIAP MEMINDAI' }}
          </div>
        </div>
      </div>
    </section>

    <!-- MAIN SINGLE-PAGE SCANNER WORKSPACE -->
    <div class="w-full max-w-6xl mx-auto grid grid-cols-1 lg:grid-cols-12 gap-6 flex-1 pb-10">
      <!-- LEFT COLUMN: LIVE SCANNER VIEWPORT (7 COLS) -->
      <section class="lg:col-span-7 space-y-4">
        <!-- Main Scanner Frame Card -->
        <div class="bg-white rounded-2xl border border-slate-200 p-5 sm:p-6 shadow-sm space-y-4">
          <!-- Viewport Header -->
          <div class="flex items-center justify-between border-b border-slate-100 pb-3">
            <div>
              <h2 class="text-lg font-bold text-slate-900 flex items-center gap-2">
                <ScanLine :size="22" class="text-emerald-600" />
                Area Pemindaian QR Tiket
              </h2>
              <p class="text-xs text-slate-500 mt-0.5">
                Kamera kios aktif membaca kode QR secara real-time.
              </p>
            </div>

            <!-- Mode Switcher (Kamera vs Input Manual) -->
            <div class="flex items-center gap-1.5 bg-slate-100 p-1 rounded-xl">
              <button
                type="button"
                class="px-3 py-1 text-xs font-semibold rounded-lg transition-all flex items-center gap-1.5"
                :class="scannerMode === 'camera' ? 'bg-white text-emerald-700 shadow-xs' : 'text-slate-600 hover:text-slate-900'"
                @click="scannerMode = 'camera'"
              >
                <Camera :size="13" /> Kamera Kios
              </button>
              <button
                type="button"
                class="px-3 py-1 text-xs font-semibold rounded-lg transition-all flex items-center gap-1.5"
                :class="scannerMode === 'manual' ? 'bg-white text-emerald-700 shadow-xs' : 'text-slate-600 hover:text-slate-900'"
                @click="scannerMode = 'manual'"
              >
                <Keyboard :size="13" /> Input Kode
              </button>
            </div>
          </div>

          <!-- MODE A: CAMERA SCANNER VIEWPORT -->
          <div v-show="scannerMode === 'camera'" class="space-y-4">
            <div class="relative overflow-hidden rounded-2xl bg-slate-950 aspect-video flex items-center justify-center border-2 border-slate-800 shadow-inner">
              <video
                ref="nativeVideoRef"
                class="w-full h-full object-cover"
                autoplay
                playsinline
                muted
              ></video>
              <canvas ref="scanCanvasRef" class="hidden"></canvas>
              <div id="qr-camera-reader" class="hidden"></div>

              <!-- Reticle & Laser Sweep Overlay -->
              <div
                v-if="isCameraActive"
                class="pointer-events-none absolute inset-0 flex items-center justify-center"
              >
                <div
                  class="relative w-64 h-64 border-2 rounded-2xl transition-all duration-300"
                  :class="isCooldown
                    ? 'border-amber-400/80 shadow-[0_0_20px_rgba(251,191,36,0.3)]'
                    : 'border-emerald-400/80 shadow-[0_0_20px_rgba(16,185,129,0.3)]'"
                >
                  <!-- Reticle Corners -->
                  <div class="absolute -top-1 -left-1 w-6 h-6 border-t-4 border-l-4 rounded-tl-lg" :class="isCooldown ? 'border-amber-400' : 'border-emerald-400'"></div>
                  <div class="absolute -top-1 -right-1 w-6 h-6 border-t-4 border-r-4 rounded-tr-lg" :class="isCooldown ? 'border-amber-400' : 'border-emerald-400'"></div>
                  <div class="absolute -bottom-1 -left-1 w-6 h-6 border-b-4 border-l-4 rounded-bl-lg" :class="isCooldown ? 'border-amber-400' : 'border-emerald-400'"></div>
                  <div class="absolute -bottom-1 -right-1 w-6 h-6 border-b-4 border-r-4 rounded-br-lg" :class="isCooldown ? 'border-amber-400' : 'border-emerald-400'"></div>

                  <!-- Laser Sweep -->
                  <div
                    class="absolute inset-x-0 h-1 animate-scanner-laser shadow-[0_0_12px]"
                    :class="isCooldown
                      ? 'bg-gradient-to-r from-transparent via-amber-400 to-transparent shadow-amber-400'
                      : 'bg-gradient-to-r from-transparent via-emerald-400 to-transparent shadow-emerald-400'"
                  ></div>

                  <!-- Cooldown Overlay -->
                  <div
                    v-if="isCooldown"
                    class="absolute inset-0 bg-slate-950/70 backdrop-blur-xs rounded-xl flex flex-col items-center justify-center p-3 text-center text-amber-300 gap-1"
                  >
                    <Hourglass :size="24" class="animate-spin text-amber-400" />
                    <span class="text-xs font-bold uppercase tracking-wider">Jeda Sensor</span>
                    <span class="text-xl font-mono font-black text-white">({{ cooldownRemaining }}s)</span>
                  </div>
                </div>
              </div>

              <!-- Camera Inactive Placeholder -->
              <div
                v-if="!isCameraActive"
                class="absolute inset-0 flex flex-col items-center justify-center p-6 text-center text-slate-400 gap-3 bg-slate-900/90"
              >
                <div class="p-4 rounded-2xl bg-slate-800 border border-slate-700 text-slate-400">
                  <CameraOff :size="36" />
                </div>
                <div class="max-w-xs space-y-1">
                  <p class="font-bold text-slate-200 text-sm">Kamera Belum Dinyalakan</p>
                  <p class="text-xs text-slate-400">
                    Nyalakan kamera untuk mengaktifkan pemindaian visual QR secara otomatis.
                  </p>
                </div>
              </div>
            </div>

            <!-- Camera Controls & Device Picker -->
            <div class="flex flex-wrap items-center gap-3">
              <button
                type="button"
                class="flex-1 py-2.5 px-4 rounded-xl font-bold text-xs sm:text-sm transition-all shadow-xs flex items-center justify-center gap-2"
                :class="isCameraActive
                  ? 'bg-rose-600 hover:bg-rose-700 text-white'
                  : 'bg-emerald-600 hover:bg-emerald-700 text-white shadow-emerald-600/20'"
                :disabled="cameraLoading"
                @click="toggleCamera"
              >
                <Loader2 v-if="cameraLoading" :size="16" class="animate-spin" />
                <Square v-else-if="isCameraActive" :size="16" />
                <Play v-else :size="16" />
                {{ isCameraActive ? 'Matikan Kamera' : 'Nyalakan Kamera Kios' }}
              </button>

              <div v-if="cameraDevices.length > 1" class="flex items-center gap-1.5">
                <select
                  id="camera-select"
                  v-model="selectedCameraId"
                  :disabled="isCameraActive || cameraLoading"
                  class="rounded-xl border border-slate-300 py-2 px-3 text-xs bg-white text-slate-800 focus:ring-2 focus:ring-emerald-500 focus:outline-none"
                  @change="onCameraChange"
                >
                  <option v-for="dev in cameraDevices" :key="dev.id" :value="dev.id">
                    {{ dev.label || `Kamera ${dev.id}` }}
                  </option>
                </select>
                <button
                  type="button"
                  title="Muat ulang kamera"
                  class="p-2 rounded-xl border border-slate-200 bg-white hover:bg-slate-50 text-slate-600 transition-colors"
                  @click="refreshCameras"
                >
                  <RefreshCw :size="14" />
                </button>
              </div>
            </div>

            <!-- Camera Error Card -->
            <div
              v-if="cameraError"
              class="rounded-xl border p-3.5 text-xs space-y-2"
              :class="isDeviceInUseError
                ? 'bg-amber-50 border-amber-300 text-amber-900'
                : 'bg-rose-50 border-rose-200 text-rose-900'"
            >
              <div class="flex items-start gap-2">
                <AlertTriangle v-if="isDeviceInUseError" :size="16" class="text-amber-600 shrink-0 mt-0.5" />
                <AlertCircle v-else :size="16" class="text-rose-600 shrink-0 mt-0.5" />
                <p class="leading-relaxed font-medium">
                  {{ cameraError }}
                </p>
              </div>
              <button
                type="button"
                class="py-1 px-2.5 rounded-lg font-semibold text-xs bg-white border border-slate-300 hover:bg-slate-50 text-slate-700 transition-colors"
                @click="forceResetCamera"
              >
                Coba Reset Kamera
              </button>
            </div>
          </div>

          <!-- MODE B: MANUAL CODE INPUT VIEWPORT -->
          <div v-show="scannerMode === 'manual'" class="space-y-4">
            <form @submit.prevent="handleManualSubmit" class="space-y-3">
              <div>
                <label for="manual-token" class="block text-xs font-semibold text-slate-700 mb-1">
                  Masukkan Kode Token QR Tiket:
                </label>
                <div class="relative">
                  <input
                    id="manual-token"
                    v-model="manualToken"
                    type="text"
                    required
                    placeholder="Contoh: qr-ahmad-fauzi-004 atau token registrasi..."
                    class="w-full rounded-xl border border-slate-300 py-2.5 pl-3.5 pr-10 text-sm font-mono text-slate-900 focus:ring-2 focus:ring-emerald-500 focus:outline-none bg-white"
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

              <!-- Quick Demo Chips for Testing -->
              <div class="flex flex-wrap items-center gap-1.5 pt-1">
                <span class="text-xs text-slate-400">Pilihan Cepat:</span>
                <button
                  v-for="chip in quickChips"
                  :key="chip.code"
                  type="button"
                  class="text-[11px] font-mono px-2 py-1 rounded bg-slate-100 hover:bg-slate-200 text-slate-700 transition-colors border border-slate-200"
                  @click="manualToken = chip.code"
                >
                  {{ chip.label }}
                </button>
              </div>

              <button
                type="submit"
                class="w-full py-2.5 px-4 rounded-xl font-bold text-sm transition-all shadow-xs flex items-center justify-center gap-2"
                :class="isCooldown
                  ? 'bg-slate-200 text-slate-400 cursor-not-allowed'
                  : 'bg-emerald-600 hover:bg-emerald-700 text-white shadow-emerald-600/20'"
                :disabled="verifyingToken === manualToken || !manualToken.trim() || isCooldown"
              >
                <Loader2 v-if="verifyingToken === manualToken" :size="16" class="animate-spin" />
                <CheckCircle2 v-else :size="16" />
                {{ isCooldown ? `Jeda (${cooldownRemaining}s)` : 'Verifikasi Check-In' }}
              </button>
            </form>
          </div>
        </div>

        <!-- HARDWARE SCANNER PERPETUAL LISTENER STATUS CARD -->
        <div class="bg-white rounded-2xl border border-slate-200 p-4 sm:p-5 shadow-xs flex items-center justify-between gap-4">
          <div class="flex items-center gap-3 min-w-0">
            <div
              class="w-10 h-10 rounded-xl flex items-center justify-center shrink-0 border transition-all"
              :class="scannerSignalActive
                ? 'bg-amber-100 border-amber-300 text-amber-700 shadow-sm animate-bounce'
                : 'bg-emerald-50 border-emerald-200 text-emerald-700'"
            >
              <Barcode :size="20" />
            </div>
            <div class="min-w-0">
              <div class="flex items-center gap-2">
                <span class="text-xs font-bold text-slate-900">Mesin Scan Fisik (USB / Bluetooth HID)</span>
                <span
                  class="text-[10px] font-bold px-1.5 py-0.2 rounded-full uppercase"
                  :class="scannerSignalActive ? 'bg-amber-100 text-amber-800' : 'bg-emerald-100 text-emerald-800'"
                >
                  {{ scannerSignalActive ? 'Data Terbaca!' : 'Selalu Siaga' }}
                </span>
              </div>
              <p class="text-xs text-slate-500 truncate mt-0.5">
                <template v-if="scannerSignalActive">
                  Terdeteksi scan barcode: <strong class="text-slate-800 font-mono">{{ lastDetectedScanString }}</strong>
                </template>
                <template v-else>
                  Tembakkan barcode scanner langsung ke QR tiket kapan saja tanpa perlu klik kursor.
                </template>
              </p>
            </div>
          </div>

          <div class="flex items-center gap-1.5 text-xs font-semibold text-emerald-700 shrink-0">
            <span class="w-2 h-2 rounded-full bg-emerald-500 animate-pulse"></span>
            <span class="hidden sm:inline">Plug & Play Ready</span>
          </div>
        </div>
      </section>

      <!-- RIGHT COLUMN: VERIFICATION RESULT CARD (5 COLS) -->
      <section class="lg:col-span-5 space-y-4">
        <div class="bg-white rounded-2xl border border-slate-200 p-5 sm:p-6 shadow-sm space-y-4 min-h-[420px] flex flex-col justify-between">
          <div>
            <div class="flex items-center justify-between border-b border-slate-100 pb-3 mb-4">
              <h3 class="font-bold text-slate-900 flex items-center gap-2 text-base">
                <Activity :size="18" class="text-emerald-600" />
                Hasil Verifikasi Pasien
              </h3>
              <span v-if="lastResult" class="text-xs font-mono text-slate-400">
                {{ lastResultTime }} WIB
              </span>
            </div>

            <!-- STATE 1: STANDBY / READY TO SCAN -->
            <div
              v-if="!lastResult && !verifyingToken"
              class="py-16 px-4 text-center text-slate-400 flex flex-col items-center justify-center gap-3"
            >
              <div class="p-4 rounded-2xl bg-slate-50 border border-slate-200 text-slate-400">
                <ScanLine :size="42" class="text-emerald-600 animate-pulse" />
              </div>
              <div class="space-y-1">
                <p class="font-bold text-slate-800 text-base">Silakan Pindai Tiket Anda</p>
                <p class="text-xs text-slate-500 max-w-xs leading-relaxed">
                  Dekatkan QR tiket antrean ke kamera atau sorotkan laser mesin scanner fisik.
                </p>
              </div>
            </div>

            <!-- STATE 2: LOADING / VERIFYING -->
            <div
              v-else-if="verifyingToken"
              class="py-16 px-4 text-center flex flex-col items-center justify-center gap-3 text-emerald-600"
            >
              <div class="relative">
                <div class="w-14 h-14 rounded-full border-4 border-emerald-100 border-t-emerald-600 animate-spin"></div>
                <div class="absolute inset-0 flex items-center justify-center">
                  <QrCode :size="20" class="text-emerald-600" />
                </div>
              </div>
              <p class="font-bold text-slate-900 text-sm">Memvalidasi Tiket Pasien...</p>
              <p class="text-xs text-slate-500">Menghubungi database Supabase PostgreSQL...</p>
            </div>

            <!-- STATE 3A: SUCCESS (NEW CHECK-IN) -->
            <div
              v-else-if="lastResult && lastResult.success && !lastResult.data?.already_checked_in"
              class="rounded-2xl border border-emerald-200 bg-gradient-to-br from-emerald-50 via-teal-50/30 to-white p-5 space-y-4 shadow-sm"
            >
              <div class="flex items-start gap-3">
                <div class="p-2.5 rounded-xl bg-emerald-600 text-white shadow-sm shadow-emerald-600/30 shrink-0">
                  <CheckCircle2 :size="26" />
                </div>
                <div class="min-w-0">
                  <span class="px-2 py-0.5 rounded text-[11px] font-bold bg-emerald-600 text-white uppercase tracking-wider">
                    Check-In Berhasil
                  </span>
                  <h4 class="text-xl font-bold text-slate-900 mt-1 truncate">
                    {{ lastResultPatientName || 'Pasien AntriSehat' }}
                  </h4>
                </div>
              </div>

              <!-- Patient ID & Ticket Code Box -->
              <div class="p-4 rounded-xl bg-white border border-emerald-200/80 shadow-xs space-y-3">
                <div class="flex items-center justify-between">
                  <div>
                    <span class="text-xs text-slate-500 font-medium">Nomor Antrean Anda</span>
                    <p class="text-4xl font-black text-emerald-700 tracking-tight">
                      {{ lastResult.data?.queue_code || 'A-001' }}
                    </p>
                  </div>
                  <div class="text-right">
                    <span class="text-xs text-slate-500 font-medium">Status</span>
                    <div>
                      <span class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-bold bg-emerald-100 text-emerald-800">
                        <span class="w-1.5 h-1.5 rounded-full bg-emerald-600"></span>
                        Checked-In
                      </span>
                    </div>
                  </div>
                </div>

                <!-- Patient ID and Polyclinic details -->
                <div class="pt-3 border-t border-slate-100 space-y-2 text-xs">
                  <div class="flex items-center justify-between">
                    <span class="text-slate-500">ID Pasien:</span>
                    <span class="font-mono font-bold text-slate-900 bg-slate-50 px-2 py-0.5 rounded border border-slate-200 truncate max-w-[210px]">
                      {{ lastResultPatientId || lastResult.data?.patient_id || '-' }}
                    </span>
                  </div>
                  <div class="flex items-center justify-between">
                    <span class="text-slate-500">Nama Lengkap:</span>
                    <span class="font-bold text-slate-900">{{ lastResultPatientName || '-' }}</span>
                  </div>
                  <div class="flex items-center justify-between">
                    <span class="text-slate-500">Poli Tujuan:</span>
                    <span class="font-semibold text-slate-800">{{ lastResultPolyclinic || 'Poli Umum' }}</span>
                  </div>
                  <div class="flex items-center justify-between">
                    <span class="text-slate-500">Waktu Check-In:</span>
                    <span class="font-semibold text-slate-800">{{ lastResultCheckInTime }} WIB</span>
                  </div>
                </div>

                <!-- Wait time estimate info -->
                <div v-if="queueEstimate" class="p-2.5 rounded-lg bg-emerald-50 border border-emerald-100 text-xs text-emerald-900 flex items-center justify-between font-semibold">
                  <span>Estimasi Menunggu:</span>
                  <span class="text-emerald-700 font-bold">~{{ queueEstimate.estimate_minutes }} Menit ({{ queueEstimate.ahead_count }} di depan)</span>
                </div>
              </div>

              <div class="p-2.5 rounded-xl bg-emerald-100/60 text-xs text-emerald-900 flex items-start gap-2">
                <Info :size="15" class="text-emerald-700 shrink-0 mt-0.5" />
                <p>Silakan simpan struk Anda dan langsung menuju ruang tunggu <strong>{{ lastResultPolyclinic }}</strong>.</p>
              </div>
            </div>

            <!-- STATE 3B: DUPLICATE CHECK-IN PREVENTED -->
            <div
              v-else-if="lastResult && lastResult.success && lastResult.data?.already_checked_in"
              class="rounded-2xl border border-sky-300 bg-gradient-to-br from-sky-50 via-blue-50/40 to-white p-5 space-y-4 shadow-sm"
            >
              <div class="flex items-start gap-3">
                <div class="p-2.5 rounded-xl bg-sky-600 text-white shadow-sm shadow-sky-600/30 shrink-0">
                  <CheckCheck :size="26" />
                </div>
                <div class="min-w-0">
                  <span class="px-2 py-0.5 rounded text-[11px] font-bold bg-sky-600 text-white uppercase tracking-wider">
                    Sudah Check-In Sebelumnya
                  </span>
                  <h4 class="text-xl font-bold text-slate-900 mt-1 truncate">
                    {{ lastResultPatientName || 'Pasien AntriSehat' }}
                  </h4>
                </div>
              </div>

              <!-- Ticket Info -->
              <div class="p-4 rounded-xl bg-white border border-sky-200 shadow-xs space-y-3">
                <div class="flex items-center justify-between">
                  <div>
                    <span class="text-xs text-slate-500 font-medium">Nomor Antrean Anda</span>
                    <p class="text-4xl font-black text-sky-700 tracking-tight">
                      {{ lastResult.data?.queue_code || 'A-001' }}
                    </p>
                  </div>
                  <div class="text-right">
                    <span class="text-xs text-slate-500 font-medium">Status</span>
                    <div>
                      <span class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-bold bg-sky-100 text-sky-800">
                        <span class="w-1.5 h-1.5 rounded-full bg-sky-600"></span>
                        {{ formatStatus(lastResult.data?.status) }}
                      </span>
                    </div>
                  </div>
                </div>

                <div class="pt-3 border-t border-slate-100 space-y-2 text-xs">
                  <div class="flex items-center justify-between">
                    <span class="text-slate-500">ID Pasien:</span>
                    <span class="font-mono font-bold text-slate-900 bg-slate-50 px-2 py-0.5 rounded border border-slate-200 truncate max-w-[210px]">
                      {{ lastResultPatientId || lastResult.data?.patient_id || '-' }}
                    </span>
                  </div>
                  <div class="flex items-center justify-between">
                    <span class="text-slate-500">Nama Lengkap:</span>
                    <span class="font-bold text-slate-900">{{ lastResultPatientName || '-' }}</span>
                  </div>
                  <div class="flex items-center justify-between">
                    <span class="text-slate-500">Poli Tujuan:</span>
                    <span class="font-semibold text-slate-800">{{ lastResultPolyclinic || 'Poli Umum' }}</span>
                  </div>
                  <div class="flex items-center justify-between">
                    <span class="text-slate-500">Waktu Check-In Awal:</span>
                    <span class="font-semibold text-slate-800">{{ lastResultCheckInTime }} WIB</span>
                  </div>
                </div>

                <div class="p-2.5 rounded-xl bg-sky-100/60 text-xs text-sky-900 flex items-start gap-2">
                  <Info :size="15" class="text-sky-700 shrink-0 mt-0.5" />
                  <p><strong>Tidak Duplikat:</strong> Anda sudah terdaftar dan nomor antrean Anda tetap aktif. Silakan langsung menuju ruang tunggu.</p>
                </div>
              </div>
            </div>

            <!-- STATE 4: ERROR / FAILED -->
            <div
              v-else-if="lastResult && !lastResult.success"
              class="rounded-2xl border border-rose-200 bg-gradient-to-br from-rose-50 via-orange-50/20 to-white p-5 space-y-4 shadow-sm"
            >
              <div class="flex items-start gap-3">
                <div class="p-2.5 rounded-xl bg-rose-600 text-white shadow-sm shrink-0">
                  <XCircle :size="26" />
                </div>
                <div>
                  <span class="px-2 py-0.5 rounded text-[11px] font-bold bg-rose-600 text-white uppercase tracking-wider">
                    Check-In Gagal
                  </span>
                  <h4 class="text-base font-bold text-slate-900 mt-1">
                    {{ lastResult.message || 'Verifikasi Tiket Tidak Berhasil' }}
                  </h4>
                </div>
              </div>

              <div class="p-3.5 rounded-xl bg-white border border-rose-200 text-xs space-y-1.5">
                <div class="flex items-center justify-between text-slate-600">
                  <span>Kode Masalah:</span>
                  <span class="font-mono font-bold text-rose-700">{{ lastResult.code || 'BAD_REQUEST' }}</span>
                </div>
                <div class="flex items-center justify-between text-slate-600">
                  <span>Token Terpindai:</span>
                  <span class="font-mono text-slate-800 truncate max-w-[200px]">{{ lastScannedToken }}</span>
                </div>
              </div>
            </div>
          </div>

          <!-- Bottom Dismiss Button -->
          <div v-if="lastResult" class="pt-3 border-t border-slate-100 flex gap-2">
            <button
              type="button"
              class="w-full py-2.5 px-4 rounded-xl text-xs font-bold bg-slate-100 hover:bg-slate-200 text-slate-700 transition-colors"
              @click="resetLastResult"
            >
              Tutup & Siapkan Pindai Berikutnya
            </button>
          </div>
        </div>
      </section>
    </div>
  </main>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from "vue";
import { Html5Qrcode } from "html5-qrcode";
import {
  QrCode, Camera, CameraOff, Keyboard, Activity, CheckCircle2,
  XCircle, Volume2, VolumeX, ChevronLeft, CheckCheck, Info,
  Zap, Hourglass, Barcode, ScanLine, Maximize, Minimize,
  Loader2, Square, Play, RefreshCw, X, AlertCircle, AlertTriangle
} from "lucide-vue-next";
import { api } from "../../services/api";
import { useQueueStore } from "../../stores/queue";
import { extractCheckInToken } from "../../utils/qrReader";
import jsQR from "jsqr";

const queueStore = useQueueStore();

// Scanner Display Mode: Camera or Manual
const scannerMode = ref<"camera" | "manual">("camera");

// Real-time Digital Clock
const currentTimeString = ref("");
const currentDateString = ref("");
let clockIntervalId: number | null = null;

function updateClock() {
  const now = new Date();
  currentTimeString.value = now.toLocaleTimeString("id-ID", {
    hour: "2-digit",
    minute: "2-digit",
    second: "2-digit",
  }) + " WIB";

  currentDateString.value = now.toLocaleDateString("id-ID", {
    weekday: "long",
    day: "numeric",
    month: "long",
    year: "numeric",
  });
}

// Fullscreen Kiosk Mode
const scannerContainerRef = ref<HTMLElement | null>(null);
const isFullscreen = ref(false);

function toggleFullscreen() {
  if (!document.fullscreenElement) {
    document.documentElement.requestFullscreen().then(() => {
      isFullscreen.value = true;
    }).catch((err) => {
      console.warn("Fullscreen error:", err);
    });
  } else {
    document.exitFullscreen().then(() => {
      isFullscreen.value = false;
    }).catch((err) => {
      console.warn("Exit fullscreen error:", err);
    });
  }
}

function onFullscreenChange() {
  isFullscreen.value = !!document.fullscreenElement;
}

// ── COOLDOWN MECHANISM (2-3 Detik Jeda) ──────────────────────
const isCooldown = ref(false);
const cooldownDurationSeconds = 3;
const cooldownRemaining = ref(0);
const cooldownProgress = ref(0);
let cooldownTimerId: number | null = null;

function triggerCooldown(seconds = 3) {
  isCooldown.value = true;
  cooldownRemaining.value = seconds;
  cooldownProgress.value = 100;

  if (cooldownTimerId) {
    clearInterval(cooldownTimerId);
  }

  const intervalMs = 50;
  const totalSteps = (seconds * 1000) / intervalMs;
  let currentStep = 0;

  cooldownTimerId = window.setInterval(() => {
    currentStep++;
    cooldownProgress.value = Math.max(0, 100 - (currentStep / totalSteps) * 100);
    cooldownRemaining.value = Math.max(0, Math.ceil((totalSteps - currentStep) * (intervalMs / 1000)));

    if (currentStep >= totalSteps) {
      clearInterval(cooldownTimerId!);
      cooldownTimerId = null;
      isCooldown.value = false;
      cooldownRemaining.value = 0;
      cooldownProgress.value = 0;
    }
  }, intervalMs);
}

// ── HARDWARE SCANNER INTEGRATION (USB / BLUETOOTH HID) ───────
let barcodeKeystrokeBuffer = "";
let lastKeystrokeTime = 0;
const scannerSignalActive = ref(false);
const lastDetectedScanString = ref("");
let signalResetTimer: number | null = null;

function flashScannerSignal(code: string) {
  scannerSignalActive.value = true;
  lastDetectedScanString.value = code;
  playSound("scanner_beep");

  if (signalResetTimer) {
    clearTimeout(signalResetTimer);
  }
  signalResetTimer = window.setTimeout(() => {
    scannerSignalActive.value = false;
  }, 2500);
}

function handleGlobalKeyDown(event: KeyboardEvent) {
  const now = Date.now();
  const timeSinceLastKey = now - lastKeystrokeTime;
  lastKeystrokeTime = now;

  // Barcode scanners type rapidly (< 100ms per char).
  if (timeSinceLastKey > 200) {
    barcodeKeystrokeBuffer = "";
  }

  const activeEl = document.activeElement;
  const isExcludedInput = activeEl && (activeEl.id === "manual-token");

  if (event.key === "Enter") {
    if (barcodeKeystrokeBuffer.trim().length >= 4) {
      const scannedCode = barcodeKeystrokeBuffer.trim();
      barcodeKeystrokeBuffer = "";

      if (isExcludedInput) {
        return;
      }

      event.preventDefault();
      flashScannerSignal(scannedCode);

      if (isCooldown.value) {
        playSound("error");
        return;
      }

      handleCheckIn(scannedCode, "Mesin Pemindai USB / Barcode (HID)");
    }
    return;
  }

  if (event.key.length === 1 && !event.ctrlKey && !event.altKey && !event.metaKey) {
    barcodeKeystrokeBuffer += event.key;
  }
}

// ── CAMERA SCANNER STATE ────────────────────────────────────
const nativeVideoRef = ref<HTMLVideoElement | null>(null);
const scanCanvasRef = ref<HTMLCanvasElement | null>(null);
let activeMediaStream: MediaStream | null = null;
let scanAnimationId: number | null = null;

const isCameraActive = ref(false);
const cameraLoading = ref(false);
const cameraError = ref("");
const isDeviceInUseError = ref(false);
const selectedCameraId = ref("");
const cameraDevices = ref<Array<{ id: string; label: string }>>([]);
let html5QrCode: Html5Qrcode | null = null;

// Manual Input State
const manualToken = ref("");
const quickChips = [
  { label: "Ahmad Fauzi", code: "qr-ahmad-fauzi-004" },
  { label: "Dewi Lestari", code: "qr-dewi-lestari-003" },
  { label: "Budi Santoso", code: "qr-budi-santoso-002" },
  { label: "Siti Aminah", code: "qr-siti-aminah-001" },
];

// Result & Audio Feedback State
const verifyingToken = ref("");
const lastScannedToken = ref("");
const lastResult = ref<any>(null);
const lastResultTime = ref("");
const lastResultPatientId = ref("");
const lastResultPatientName = ref("");
const lastResultPolyclinic = ref("");
const lastResultCheckInTime = ref("");
const queueEstimate = ref<any>(null);
const soundEnabled = ref(true);

// Web Audio API feedback
function playSound(type: "success" | "already_checked_in" | "error" | "scanner_beep") {
  if (!soundEnabled.value) return;
  try {
    const ctx = new (window.AudioContext || (window as any).webkitAudioContext)();
    const osc = ctx.createOscillator();
    const gain = ctx.createGain();
    osc.connect(gain);
    gain.connect(ctx.destination);

    if (type === "success") {
      osc.type = "sine";
      osc.frequency.setValueAtTime(587.33, ctx.currentTime); // D5
      osc.frequency.setValueAtTime(880, ctx.currentTime + 0.09); // A5
      gain.gain.setValueAtTime(0.25, ctx.currentTime);
      gain.gain.exponentialRampToValueAtTime(0.01, ctx.currentTime + 0.35);
      osc.start(ctx.currentTime);
      osc.stop(ctx.currentTime + 0.35);
    } else if (type === "already_checked_in") {
      osc.type = "triangle";
      osc.frequency.setValueAtTime(440, ctx.currentTime); // A4
      osc.frequency.setValueAtTime(659.25, ctx.currentTime + 0.1); // E5
      gain.gain.setValueAtTime(0.2, ctx.currentTime);
      gain.gain.exponentialRampToValueAtTime(0.01, ctx.currentTime + 0.35);
      osc.start(ctx.currentTime);
      osc.stop(ctx.currentTime + 0.35);
    } else if (type === "scanner_beep") {
      osc.type = "sine";
      osc.frequency.setValueAtTime(1200, ctx.currentTime);
      gain.gain.setValueAtTime(0.15, ctx.currentTime);
      gain.gain.exponentialRampToValueAtTime(0.01, ctx.currentTime + 0.06);
      osc.start(ctx.currentTime);
      osc.stop(ctx.currentTime + 0.06);
    } else {
      osc.type = "sawtooth";
      osc.frequency.setValueAtTime(220, ctx.currentTime);
      osc.frequency.setValueAtTime(160, ctx.currentTime + 0.1);
      gain.gain.setValueAtTime(0.25, ctx.currentTime);
      gain.gain.exponentialRampToValueAtTime(0.01, ctx.currentTime + 0.25);
      osc.start(ctx.currentTime);
      osc.stop(ctx.currentTime + 0.25);
    }
  } catch {}
}

function formatStatus(status: string) {
  switch (status) {
    case "waiting": return "Menunggu";
    case "checked_in": return "Sudah Check-In";
    case "called": return "Dipanggil ke Ruang Poli";
    case "in_service": return "Sedang Dilayani Dokter";
    case "done": return "Pemeriksaan Selesai";
    default: return status || "Checked-In";
  }
}

// ── Check-In Verification Logic ──────────────────────────────
async function handleCheckIn(rawToken: string, source: string) {
  if (isCooldown.value) {
    return;
  }

  let token = extractCheckInToken(rawToken);
  if (!token) return;

  verifyingToken.value = token;
  lastScannedToken.value = token;

  try {
    let res = await api.checkInByQR(token);

    lastResult.value = res;
    lastResultTime.value = new Date().toLocaleTimeString("id-ID", {
      hour: "2-digit",
      minute: "2-digit",
      second: "2-digit",
    });

    if (res && res.success) {
      const isDuplicate = !!res.data?.already_checked_in;

      if (isDuplicate) {
        playSound("already_checked_in");
      } else {
        playSound("success");
      }

      // Check localStorage or queueStore for fallback names
      let localMatched = queueStore.patients.find(
        (p) => p.qr_token === token || p.id === token || p.number === token,
      );
      let storedPatient: any = null;
      try {
        const storedStr = localStorage.getItem("antrisehat-patients");
        if (storedStr) {
          const list = JSON.parse(storedStr);
          storedPatient = list.find((p: any) => p.qr_token === token || p.id === token || p.number === token);
        }
      } catch (_) {}

      // Explicitly return and display ID and Name from backend or stored match
      lastResultPatientId.value =
        res.data?.patient_id ||
        localMatched?.id ||
        storedPatient?.id ||
        "ID-" + (token.slice(0, 8).toUpperCase());

      lastResultPatientName.value =
        res.data?.patient_name ||
        localMatched?.name ||
        storedPatient?.name ||
        "Pasien Terdaftar";

      lastResultPolyclinic.value =
        res.data?.polyclinic_name ||
        localMatched?.poli ||
        storedPatient?.poli ||
        "Poli Umum";

      lastResultCheckInTime.value = new Date().toLocaleTimeString("id-ID", {
        hour: "2-digit",
        minute: "2-digit",
      });

      if (localMatched) {
        localMatched.status = "Hadir";
        queueStore.changeStatus(localMatched.id, "Hadir");
      }
      if (storedPatient && storedPatient.id) {
        queueStore.changeStatus(storedPatient.id, "Hadir");
      }

      if (res.data?.id) {
        api.getQueueEstimate(res.data.id).then((est) => {
          queueEstimate.value = est;
        });
      }
    } else {
      playSound("error");
    }
  } catch (err: any) {
    playSound("error");
    lastResult.value = {
      success: false,
      code: "NETWORK_ERROR",
      message: err.message || "Gagal menghubungi server",
    };
  } finally {
    verifyingToken.value = "";
    triggerCooldown(cooldownDurationSeconds);
  }
}

function resetLastResult() {
  lastResult.value = null;
  queueEstimate.value = null;
}

// ── Tab 2: Live Camera Scanner ──────────────────────────────
async function refreshCameras() {
  try {
    const devices = await Html5Qrcode.getCameras();
    if (devices && devices.length > 0) {
      cameraDevices.value = devices;
      if (!selectedCameraId.value) {
        const preferred = devices.find(
          (d) =>
            !d.label.toLowerCase().includes("ir") &&
            !d.label.toLowerCase().includes("virtual") &&
            !d.label.toLowerCase().includes("hello"),
        );
        selectedCameraId.value = preferred ? preferred.id : devices[0].id;
      }
    }
  } catch (err) {
    console.warn("Could not list camera devices:", err);
  }
}

async function onCameraChange() {
  if (isCameraActive.value) {
    await forceResetCamera();
  }
}

async function forceResetCamera() {
  await stopCamera();
  await initCamera();
}

async function stopCamera() {
  cameraLoading.value = false;
  if (scanAnimationId) {
    cancelAnimationFrame(scanAnimationId);
    scanAnimationId = null;
  }

  if (activeMediaStream) {
    activeMediaStream.getTracks().forEach((track) => {
      try {
        track.stop();
      } catch (_) {}
    });
    activeMediaStream = null;
  }

  if (nativeVideoRef.value) {
    nativeVideoRef.value.srcObject = null;
  }

  if (html5QrCode) {
    try {
      if (html5QrCode.isScanning) {
        await html5QrCode.stop();
      }
      await html5QrCode.clear();
    } catch (_) {}
    html5QrCode = null;
  }

  try {
    const videos = document.querySelectorAll("video");
    videos.forEach((video) => {
      if (video.srcObject && (video.srcObject as MediaStream).getTracks) {
        (video.srcObject as MediaStream).getTracks().forEach((track) => {
          track.stop();
        });
        video.srcObject = null;
      }
    });
  } catch (_) {}

  isCameraActive.value = false;
}

function startNativeScanLoop() {
  const video = nativeVideoRef.value;
  if (!video) return;

  const canvas = scanCanvasRef.value || document.createElement("canvas");
  const ctx = canvas.getContext("2d", { willReadFrequently: true });
  if (!ctx) return;
  const context: CanvasRenderingContext2D = ctx;

  function tick() {
    if (!isCameraActive.value) return;

    if (video && video.readyState === video.HAVE_ENOUGH_DATA) {
      canvas.width = video.videoWidth;
      canvas.height = video.videoHeight;
      context.drawImage(video, 0, 0, canvas.width, canvas.height);

      if (!isCooldown.value && !verifyingToken.value) {
        const imageData = context.getImageData(0, 0, canvas.width, canvas.height);
        const code = jsQR(imageData.data, imageData.width, imageData.height, {
          inversionAttempts: "dontInvert",
        });

        if (code && code.data && !verifyingToken.value) {
          handleCheckIn(code.data, "Kamera Kios Live");
        }
      }
    }

    scanAnimationId = requestAnimationFrame(tick);
  }

  scanAnimationId = requestAnimationFrame(tick);
}

async function initCamera() {
  try {
    cameraLoading.value = true;
    cameraError.value = "";
    isDeviceInUseError.value = false;

    await stopCamera();
    await refreshCameras();

    const constraintsList: MediaStreamConstraints[] = [];

    if (selectedCameraId.value) {
      constraintsList.push({
        video: { deviceId: { exact: selectedCameraId.value } },
        audio: false,
      });
    }

    for (const dev of cameraDevices.value) {
      if (dev.id !== selectedCameraId.value) {
        constraintsList.push({
          video: { deviceId: { exact: dev.id } },
          audio: false,
        });
      }
    }

    constraintsList.push({ video: { facingMode: "user" }, audio: false });
    constraintsList.push({ video: { facingMode: "environment" }, audio: false });
    constraintsList.push({ video: true, audio: false });

    let stream: MediaStream | null = null;
    let lastErr: any = null;

    for (const constraints of constraintsList) {
      try {
        stream = await navigator.mediaDevices.getUserMedia(constraints);
        if (stream) {
          activeMediaStream = stream;
          break;
        }
      } catch (err: any) {
        lastErr = err;
      }
    }

    if (!stream && lastErr) {
      throw lastErr;
    }

    if (stream && nativeVideoRef.value) {
      nativeVideoRef.value.srcObject = stream;
      await nativeVideoRef.value.play();
      isCameraActive.value = true;
      startNativeScanLoop();
    }
  } catch (err: any) {
    isCameraActive.value = false;
    const msg = String(err?.message || err || "");
    const name = String(err?.name || "");

    if (
      name === "NotReadableError" ||
      msg.toLowerCase().includes("in use") ||
      msg.toLowerCase().includes("could not start video source") ||
      msg.toLowerCase().includes("not readable")
    ) {
      isDeviceInUseError.value = true;
      cameraError.value =
        "Perangkat kamera sedang dikunci atau digunakan oleh aplikasi lain di komputer ini.";
    } else if (
      name === "NotAllowedError" ||
      msg.toLowerCase().includes("permission") ||
      msg.toLowerCase().includes("denied")
    ) {
      cameraError.value =
        "Izin akses kamera ditolak. Silakan klik ikon gembok di address bar browser dan izinkan kamera.";
    } else if (
      name === "NotFoundError" ||
      msg.toLowerCase().includes("no camera") ||
      msg.toLowerCase().includes("not found")
    ) {
      cameraError.value =
        "Tidak ditemukan perangkat kamera yang terpasang pada komputer ini.";
    } else {
      cameraError.value =
        err?.message || "Tidak dapat mengaktifkan kamera scanner.";
    }
  } finally {
    cameraLoading.value = false;
  }
}

async function toggleCamera() {
  if (isCameraActive.value) {
    await stopCamera();
  } else {
    await initCamera();
  }
}

function handleManualSubmit() {
  if (manualToken.value.trim() && !isCooldown.value) {
    handleCheckIn(manualToken.value, "Input Manual");
  }
}

onMounted(async () => {
  updateClock();
  clockIntervalId = window.setInterval(updateClock, 1000);

  // Global Hardware Barcode Scanner listener
  window.addEventListener("keydown", handleGlobalKeyDown);
  document.addEventListener("fullscreenchange", onFullscreenChange);

  await refreshCameras();
  await initCamera();

  const params = new URLSearchParams(window.location.search);
  const urlToken = params.get("token") || params.get("checkin");
  if (urlToken) {
    scannerMode.value = "manual";
    manualToken.value = urlToken;
    handleCheckIn(urlToken, "Tautan URL / Parameter");
  }
});

onUnmounted(() => {
  stopCamera();
  if (clockIntervalId) {
    clearInterval(clockIntervalId);
  }
  if (cooldownTimerId) {
    clearInterval(cooldownTimerId);
  }
  if (signalResetTimer) {
    clearTimeout(signalResetTimer);
  }
  window.removeEventListener("keydown", handleGlobalKeyDown);
  document.removeEventListener("fullscreenchange", onFullscreenChange);
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
