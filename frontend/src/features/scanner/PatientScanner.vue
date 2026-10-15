<template>
  <main
    ref="kioskContainerRef"
    class="min-h-screen bg-slate-900 text-slate-100 flex flex-col font-sans transition-colors duration-300 select-none"
    :class="{ 'p-4 sm:p-6 lg:p-8': !isFullscreen, 'p-6 sm:p-10': isFullscreen }"
  >
    <!-- TOP KIOSK HEADER BAR -->
    <header class="w-full max-w-7xl mx-auto mb-6 bg-slate-800/80 backdrop-blur-md rounded-2xl border border-slate-700/80 p-4 sm:p-5 shadow-xl shadow-black/20 flex flex-col md:flex-row md:items-center md:justify-between gap-4">
      <!-- Hospital Branding -->
      <div class="flex items-center gap-3.5">
        <div class="w-12 h-12 rounded-2xl bg-gradient-to-tr from-emerald-600 to-teal-400 flex items-center justify-center text-white shadow-lg shadow-emerald-500/25 shrink-0 ring-2 ring-emerald-400/30">
          <Activity :size="26" class="animate-pulse" />
        </div>
        <div>
          <div class="flex items-center gap-2">
            <span class="text-[11px] font-bold tracking-widest uppercase px-2 py-0.5 rounded-full bg-emerald-500/20 text-emerald-300 border border-emerald-500/30">
              Kios Check-In Mandiri Pasien
            </span>
            <span class="text-xs text-slate-400 font-mono">v2.5</span>
          </div>
          <h1 class="text-xl sm:text-2xl font-black text-white tracking-tight mt-0.5 flex items-center gap-2">
            Puskesmas Sehat Bersama
          </h1>
          <p class="text-xs text-slate-400 hidden sm:block">
            Anjungan Pelayanan Pendaftaran Terpadu · Terhubung Langsung dengan Supabase Cloud
          </p>
        </div>
      </div>

      <!-- Real-time Digital Monospace Clock -->
      <div class="flex items-center justify-between md:justify-center px-4 py-2 rounded-xl bg-slate-900/90 border border-slate-700 text-center shadow-inner">
        <div>
          <div class="text-2xl sm:text-3xl font-mono font-bold tracking-wider text-emerald-400 drop-shadow-[0_0_8px_rgba(52,211,153,0.3)]">
            {{ currentTimeString || '00:00:00 WIB' }}
          </div>
          <div class="text-[11px] font-medium text-slate-400 uppercase tracking-wider">
            {{ currentDateString || 'Memuat Waktu...' }}
          </div>
        </div>
      </div>

      <!-- System Connectivity Badges & Controls -->
      <div class="flex flex-wrap items-center gap-2.5">
        <!-- Hardware Scanner Status Badge -->
        <div
          class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-xl text-xs font-semibold border transition-all"
          :class="scannerSignalActive
            ? 'bg-amber-500/20 border-amber-400 text-amber-300 shadow-[0_0_12px_rgba(245,158,11,0.4)]'
            : 'bg-emerald-950/60 border-emerald-500/30 text-emerald-400'"
        >
          <span class="relative flex h-2 w-2">
            <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75"></span>
            <span class="relative inline-flex rounded-full h-2 w-2" :class="scannerSignalActive ? 'bg-amber-400' : 'bg-emerald-500'"></span>
          </span>
          <Barcode :size="14" />
          <span>{{ scannerSignalActive ? 'Data Scanner Diterima!' : 'Scanner Fisik (HID)' }}</span>
        </div>

        <!-- Camera Status Badge -->
        <div class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-xl text-xs font-semibold bg-slate-900 border border-slate-700 text-slate-300">
          <Camera :size="14" :class="isCameraActive ? 'text-emerald-400' : 'text-slate-500'" />
          <span>{{ isCameraActive ? 'Kamera Aktif' : 'Kamera Siaga' }}</span>
        </div>

        <!-- Audio Toggle Button -->
        <button
          type="button"
          class="p-2 sm:px-3 sm:py-2 rounded-xl text-xs font-medium border border-slate-700 bg-slate-900 hover:bg-slate-700/80 text-slate-200 transition-colors shadow-sm flex items-center gap-1.5"
          @click="soundEnabled = !soundEnabled"
          :title="soundEnabled ? 'Matikan Nada Konfirmasi' : 'Aktifkan Nada Konfirmasi'"
        >
          <Volume2 v-if="soundEnabled" :size="16" class="text-emerald-400" />
          <VolumeX v-else :size="16" class="text-slate-400" />
          <span class="hidden sm:inline">{{ soundEnabled ? 'Audio Aktif' : 'Mute' }}</span>
        </button>

        <!-- Fullscreen Kiosk Mode Toggle -->
        <button
          type="button"
          class="p-2 sm:px-3 sm:py-2 rounded-xl text-xs font-semibold border transition-all flex items-center gap-1.5 shadow-sm"
          :class="isFullscreen
            ? 'bg-emerald-600 border-emerald-500 text-white shadow-emerald-500/20'
            : 'bg-slate-900 border-slate-700 hover:bg-slate-700/80 text-slate-200'"
          @click="toggleFullscreen"
          title="Mode Layar Penuh (Kiosk Station)"
        >
          <Minimize v-if="isFullscreen" :size="16" />
          <Maximize v-else :size="16" />
          <span class="hidden sm:inline">{{ isFullscreen ? 'Keluar Fullscreen' : 'Layar Penuh' }}</span>
        </button>

        <!-- Back to Staff Dashboard -->
        <RouterLink
          to="/staff"
          class="p-2 sm:px-3 sm:py-2 rounded-xl text-xs font-semibold bg-slate-900 hover:bg-slate-700/80 border border-slate-700 text-slate-300 transition-colors flex items-center gap-1"
          title="Buka Dashboard Staf"
        >
          <ChevronLeft :size="15" />
          <span class="hidden sm:inline">Staf</span>
        </RouterLink>
      </div>
    </header>

    <!-- COOLDOWN SENSOR STATUS BANNER -->
    <section class="w-full max-w-7xl mx-auto mb-6">
      <div
        class="relative overflow-hidden rounded-2xl border transition-all duration-300 p-4 shadow-lg flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3"
        :class="isCooldown
          ? 'bg-gradient-to-r from-amber-950/80 via-slate-900 to-amber-950/80 border-amber-500/50 shadow-amber-950/50'
          : 'bg-gradient-to-r from-emerald-950/70 via-slate-900 to-teal-950/70 border-emerald-500/40 shadow-emerald-950/30'"
      >
        <!-- Cooldown Progress Track (animated smoothly) -->
        <div
          v-if="isCooldown"
          class="absolute bottom-0 left-0 h-1.5 bg-gradient-to-r from-amber-400 via-orange-400 to-amber-300 transition-all duration-75 ease-linear shadow-[0_0_10px_rgba(251,191,36,0.8)]"
          :style="{ width: `${cooldownProgress}%` }"
        ></div>

        <div class="flex items-center gap-3.5 min-w-0">
          <div
            class="w-10 h-10 rounded-xl flex items-center justify-center shrink-0 shadow-md transition-colors"
            :class="isCooldown
              ? 'bg-amber-500/20 text-amber-300 border border-amber-400/30 animate-pulse'
              : 'bg-emerald-500/20 text-emerald-300 border border-emerald-400/30'"
          >
            <Hourglass v-if="isCooldown" :size="20" class="animate-spin" />
            <Zap v-else :size="20" class="animate-bounce" />
          </div>

          <div class="min-w-0">
            <div class="flex items-center gap-2">
              <span
                class="text-xs font-black uppercase tracking-wider px-2 py-0.5 rounded-md"
                :class="isCooldown
                  ? 'bg-amber-500/20 text-amber-300 border border-amber-500/40'
                  : 'bg-emerald-500/20 text-emerald-300 border border-emerald-500/40'"
              >
                {{ isCooldown ? 'Jeda Sensor Aktif' : 'Sensor Siap Memindai' }}
              </span>
              <span class="text-xs text-slate-400 font-mono hidden sm:inline">
                Interval: 3 Detik Anti-Duplikasi
              </span>
            </div>

            <p class="text-sm font-semibold text-slate-100 truncate mt-0.5">
              <template v-if="isCooldown">
                Menyiapkan mesin untuk pasien berikutnya... Mohon tunggu
                <span class="font-mono font-bold text-amber-400 text-base">({{ cooldownRemaining }}s)</span>
              </template>
              <template v-else>
                Silakan tempelkan kode QR tiket Anda ke mesin scanner fisik atau hadapkan ke kamera kios.
              </template>
            </p>
          </div>
        </div>

        <!-- Sensor Status Indicator Pill -->
        <div class="flex items-center gap-3 self-end sm:self-auto shrink-0">
          <div
            class="px-3.5 py-1.5 rounded-xl text-xs font-mono font-bold flex items-center gap-2 border"
            :class="isCooldown
              ? 'bg-amber-900/40 border-amber-500/40 text-amber-300'
              : 'bg-emerald-900/40 border-emerald-500/40 text-emerald-300'"
          >
            <span class="w-2.5 h-2.5 rounded-full" :class="isCooldown ? 'bg-amber-400 animate-ping' : 'bg-emerald-400'"></span>
            {{ isCooldown ? `COOLDOWN (${cooldownRemaining}s)` : 'SENSOR READY' }}
          </div>
        </div>
      </div>
    </section>

    <!-- KIOSK NAVIGATION TABS -->
    <nav class="w-full max-w-7xl mx-auto mb-6 flex flex-wrap gap-2.5">
      <button
        v-for="tab in tabs"
        :key="tab.id"
        type="button"
        class="flex items-center gap-2 px-4 py-2.5 rounded-xl text-xs sm:text-sm font-bold transition-all shadow-md"
        :class="activeTab === tab.id
          ? 'bg-emerald-600 text-white shadow-emerald-600/30 ring-2 ring-emerald-400/30'
          : 'bg-slate-800/90 text-slate-300 hover:bg-slate-700/80 border border-slate-700 hover:text-white'"
        @click="switchTab(tab.id)"
      >
        <component :is="tab.icon" :size="17" />
        <span>{{ tab.label }}</span>
        <span
          v-if="tab.badge"
          class="ml-1 px-1.5 py-0.5 text-[10px] font-bold rounded-full uppercase tracking-wider"
          :class="activeTab === tab.id ? 'bg-white/20 text-white' : 'bg-emerald-500/20 text-emerald-300'"
        >
          {{ tab.badge }}
        </span>
      </button>
    </nav>

    <!-- MAIN TWO-COLUMN WORKSPACE -->
    <div class="w-full max-w-7xl mx-auto grid grid-cols-1 lg:grid-cols-12 gap-6 flex-1 pb-10">
      <!-- LEFT COLUMN: SCANNER INPUT ENGINE (7 COLS) -->
      <section class="lg:col-span-7 space-y-6">
        <!-- TAB 1: HARDWARE BARCODE SCANNER (HID PLUG & PLAY) -->
        <div
          v-show="activeTab === 'hardware'"
          class="bg-slate-800/90 rounded-2xl border border-slate-700/80 p-5 sm:p-6 shadow-xl space-y-5"
        >
          <div class="flex items-center justify-between border-b border-slate-700/70 pb-4">
            <div>
              <h2 class="text-lg font-bold text-white flex items-center gap-2">
                <Barcode :size="22" class="text-emerald-400" />
                Mesin Barcode / QR Scanner Fisik
              </h2>
              <p class="text-xs text-slate-400 mt-1">
                Mendukung scanner genggam atau barcode reader meja (Honeywell, Zebra, Sunmi, Datalogic, Eyoyo USB/Bluetooth HID).
              </p>
            </div>
            <span class="px-2.5 py-1 text-xs font-mono font-bold rounded-full bg-emerald-500/20 text-emerald-300 border border-emerald-500/30">
              Auto-Listen Aktif
            </span>
          </div>

          <!-- Kiosk Scanner Terminal Graphic -->
          <div class="relative overflow-hidden rounded-2xl bg-slate-950 border border-slate-800 p-8 flex flex-col items-center justify-center text-center gap-4">
            <!-- Animated Laser Scan Target -->
            <div class="relative w-48 h-48 rounded-2xl border-2 border-dashed border-emerald-500/40 flex items-center justify-center bg-slate-900/60 p-4">
              <div class="w-full h-full rounded-xl bg-slate-950 flex flex-col items-center justify-center gap-2 text-slate-500 border border-slate-800 relative overflow-hidden">
                <QrCode :size="64" class="text-slate-600 transition-colors" :class="{ '!text-emerald-400': scannerSignalActive }" />
                <span class="text-[11px] font-mono text-slate-400 uppercase tracking-wider">Arahkan Laser ke Sini</span>

                <!-- Animated Laser Line Sweep -->
                <div class="absolute inset-x-0 h-0.5 bg-gradient-to-r from-transparent via-emerald-400 to-transparent shadow-[0_0_12px_#34d399] animate-scanner-laser"></div>
              </div>

              <!-- Reticle Corner Accents -->
              <div class="absolute -top-1 -left-1 w-6 h-6 border-t-2 border-l-2 border-emerald-400 rounded-tl-lg"></div>
              <div class="absolute -top-1 -right-1 w-6 h-6 border-t-2 border-r-2 border-emerald-400 rounded-tr-lg"></div>
              <div class="absolute -bottom-1 -left-1 w-6 h-6 border-b-2 border-l-2 border-emerald-400 rounded-bl-lg"></div>
              <div class="absolute -bottom-1 -right-1 w-6 h-6 border-b-2 border-r-2 border-emerald-400 rounded-br-lg"></div>
            </div>

            <!-- Status Guidance Text -->
            <div class="space-y-1.5 max-w-md">
              <h3 class="text-base font-bold text-white flex items-center justify-center gap-1.5">
                <span class="w-2.5 h-2.5 rounded-full bg-emerald-400 animate-ping"></span>
                Mesin Pemindai Siap Menerima Scan
              </h3>
              <p class="text-xs text-slate-400 leading-relaxed">
                Pasien cukup mengarahkan lembar tiket atau layar ponsel ke sorotan sinar laser mesin scanner. Sistem akan membaca token dan memverifikasi secara instan tanpa perlu menekan tombol apapun.
              </p>
            </div>

            <!-- Signal Pulse Banner -->
            <div
              v-if="scannerSignalActive"
              class="w-full py-2 px-3 rounded-xl bg-amber-500/20 border border-amber-400/40 text-amber-300 text-xs font-mono font-bold flex items-center justify-center gap-2 animate-bounce"
            >
              <Zap :size="15" />
              Sinyal Scanner Fisik Terdeteksi: {{ lastDetectedScanString }}
            </div>
          </div>

          <!-- Hardware Configuration & Testing Input -->
          <div class="rounded-xl bg-slate-900/90 border border-slate-700/80 p-4 space-y-3">
            <div class="flex items-center justify-between text-xs">
              <span class="font-bold text-slate-300 flex items-center gap-1.5">
                <Cpu :size="15" class="text-emerald-400" /> Mode Komunikasi Mesin
              </span>
              <span class="font-mono text-emerald-400 bg-emerald-950/60 px-2 py-0.5 rounded border border-emerald-500/30">
                USB HID Keyboard Emulation
              </span>
            </div>
            <p class="text-xs text-slate-400 leading-relaxed">
              Scanner fisik Anda bertindak layaknya keyboard super-cepat. Kode yang terpindai otomatis diterima dari halaman manapun tanpa perlu meletakkan kursor pada kotak input.
            </p>

            <!-- Optional manual test input for scanner gun -->
            <div class="pt-2 border-t border-slate-800">
              <label for="hardware-test-input" class="block text-[11px] font-semibold text-slate-400 uppercase tracking-wider mb-1.5">
                Area Uji Scan / Tes Mesin Fisik:
              </label>
              <div class="relative">
                <input
                  id="hardware-test-input"
                  v-model="hardwareTestInput"
                  type="text"
                  placeholder="Klik di sini lalu tembakkan barcode scanner untuk pengujian..."
                  class="w-full rounded-xl bg-slate-950 border border-slate-700 py-2.5 pl-3.5 pr-20 text-xs font-mono text-emerald-300 placeholder-slate-600 focus:outline-none focus:ring-2 focus:ring-emerald-500"
                  @keydown.enter.prevent="handleHardwareTestSubmit"
                />
                <button
                  type="button"
                  class="absolute right-1.5 top-1/2 -translate-y-1/2 px-2.5 py-1 rounded-lg text-xs font-bold bg-emerald-600 hover:bg-emerald-500 text-white transition-colors"
                  :disabled="!hardwareTestInput.trim() || isCooldown"
                  @click="handleHardwareTestSubmit"
                >
                  Test Scan
                </button>
              </div>
            </div>
          </div>
        </div>

        <!-- TAB 2: LIVE CAMERA SCANNER (WEBCAM / BUILT-IN KIOSK CAMERA) -->
        <div
          v-show="activeTab === 'camera'"
          class="bg-slate-800/90 rounded-2xl border border-slate-700/80 p-5 sm:p-6 shadow-xl space-y-5"
        >
          <div class="flex items-center justify-between border-b border-slate-700/70 pb-4">
            <div>
              <h2 class="text-lg font-bold text-white flex items-center gap-2">
                <Camera :size="22" class="text-emerald-400" />
                Kamera Scanner Kios
              </h2>
              <p class="text-xs text-slate-400 mt-1">
                Kamera web kios otomatis membaca kode QR dari struk atau layar smartphone pasien.
              </p>
            </div>
            <span
              class="px-2.5 py-1 text-xs font-semibold rounded-full border"
              :class="isCameraActive
                ? 'bg-emerald-500/20 text-emerald-300 border-emerald-500/40'
                : 'bg-slate-700/60 text-slate-400 border-slate-600'"
            >
              {{ isCameraActive ? 'Kamera Aktif' : 'Kamera Nonaktif' }}
            </span>
          </div>

          <!-- Video Viewport -->
          <div class="relative overflow-hidden rounded-2xl bg-slate-950 aspect-video flex items-center justify-center border-2 border-slate-700/80 shadow-2xl">
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
                  ? 'border-amber-400/80 shadow-[0_0_25px_rgba(251,191,36,0.4)]'
                  : 'border-emerald-400/80 shadow-[0_0_25px_rgba(16,185,129,0.4)]'"
              >
                <!-- Reticle Corners -->
                <div class="absolute -top-1 -left-1 w-6 h-6 border-t-4 border-l-4 rounded-tl-lg" :class="isCooldown ? 'border-amber-400' : 'border-emerald-400'"></div>
                <div class="absolute -top-1 -right-1 w-6 h-6 border-t-4 border-r-4 rounded-tr-lg" :class="isCooldown ? 'border-amber-400' : 'border-emerald-400'"></div>
                <div class="absolute -bottom-1 -left-1 w-6 h-6 border-b-4 border-l-4 rounded-bl-lg" :class="isCooldown ? 'border-amber-400' : 'border-emerald-400'"></div>
                <div class="absolute -bottom-1 -right-1 w-6 h-6 border-b-4 border-r-4 rounded-br-lg" :class="isCooldown ? 'border-amber-400' : 'border-emerald-400'"></div>

                <!-- Laser Line Sweep -->
                <div
                  class="absolute inset-x-0 h-1 animate-scanner-laser shadow-[0_0_12px]"
                  :class="isCooldown
                    ? 'bg-gradient-to-r from-transparent via-amber-400 to-transparent shadow-amber-400'
                    : 'bg-gradient-to-r from-transparent via-emerald-400 to-transparent shadow-emerald-400'"
                ></div>

                <!-- Cooldown Overlay Inside Reticle -->
                <div
                  v-if="isCooldown"
                  class="absolute inset-0 bg-slate-950/70 backdrop-blur-xs rounded-xl flex flex-col items-center justify-center p-3 text-center text-amber-300 gap-1.5"
                >
                  <Hourglass :size="24" class="animate-spin text-amber-400" />
                  <span class="text-xs font-bold uppercase tracking-wider">Jeda Pemindaian</span>
                  <span class="text-lg font-mono font-black text-white">({{ cooldownRemaining }}s)</span>
                </div>
              </div>
            </div>

            <!-- Inactive Placeholder -->
            <div
              v-if="!isCameraActive"
              class="absolute inset-0 flex flex-col items-center justify-center p-6 text-center text-slate-400 gap-3 bg-slate-950/90"
            >
              <div class="p-4 rounded-2xl bg-slate-900 border border-slate-800 text-slate-500">
                <CameraOff :size="40" />
              </div>
              <div class="max-w-xs space-y-1">
                <p class="font-bold text-slate-200 text-sm">Kamera Belum Aktif</p>
                <p class="text-xs text-slate-400">
                  Nyalakan kamera untuk mengaktifkan pemindaian visual otomatis.
                </p>
              </div>
            </div>
          </div>

          <!-- Camera Controls & Device Picker -->
          <div class="space-y-3">
            <div v-if="cameraDevices.length > 0" class="flex items-center gap-2">
              <label for="camera-select" class="text-xs font-semibold text-slate-300 whitespace-nowrap flex items-center gap-1.5">
                <Camera :size="14" class="text-slate-400" /> Pilih Perangkat:
              </label>
              <select
                id="camera-select"
                v-model="selectedCameraId"
                :disabled="isCameraActive || cameraLoading"
                class="flex-1 rounded-xl border border-slate-700 py-2 px-3 text-xs bg-slate-900 text-slate-200 focus:ring-2 focus:ring-emerald-500 focus:outline-none"
                @change="onCameraChange"
              >
                <option v-for="dev in cameraDevices" :key="dev.id" :value="dev.id">
                  {{ dev.label || `Kamera ${dev.id}` }}
                </option>
              </select>
              <button
                type="button"
                title="Muat ulang daftar kamera"
                class="p-2 rounded-xl border border-slate-700 bg-slate-900 hover:bg-slate-700 text-slate-300 transition-colors"
                @click="refreshCameras"
              >
                <RefreshCw :size="14" />
              </button>
            </div>

            <div class="flex items-center gap-3">
              <button
                type="button"
                class="flex-1 py-3 px-4 rounded-xl font-bold text-sm transition-all shadow-lg flex items-center justify-center gap-2"
                :class="isCameraActive
                  ? 'bg-rose-600 hover:bg-rose-700 text-white shadow-rose-600/30'
                  : 'bg-emerald-600 hover:bg-emerald-500 text-white shadow-emerald-600/30'"
                :disabled="cameraLoading"
                @click="toggleCamera"
              >
                <Loader2 v-if="cameraLoading" :size="18" class="animate-spin" />
                <Square v-else-if="isCameraActive" :size="18" />
                <Play v-else :size="18" />
                {{ isCameraActive ? 'Matikan Kamera Scanner' : 'Nyalakan Kamera Scanner' }}
              </button>
            </div>

            <!-- Camera Error & Troubleshooting Card -->
            <div
              v-if="cameraError"
              class="rounded-xl border p-4 text-xs space-y-2.5"
              :class="isDeviceInUseError
                ? 'bg-amber-950/60 border-amber-500/40 text-amber-200'
                : 'bg-rose-950/60 border-rose-500/40 text-rose-200'"
            >
              <div class="flex items-start gap-2.5">
                <AlertTriangle v-if="isDeviceInUseError" :size="18" class="text-amber-400 shrink-0 mt-0.5" />
                <AlertCircle v-else :size="18" class="text-rose-400 shrink-0 mt-0.5" />
                <div class="space-y-1">
                  <p class="font-bold text-xs uppercase tracking-wide">
                    {{ isDeviceInUseError ? 'Perangkat Kamera Sedang Digunakan (Device In Use)' : 'Kendala Akses Kamera' }}
                  </p>
                  <p class="leading-relaxed text-slate-300">
                    {{ cameraError }}
                  </p>
                </div>
              </div>

              <div class="flex flex-wrap items-center gap-2 pt-2 border-t border-slate-700/60">
                <button
                  type="button"
                  class="py-1.5 px-3 rounded-lg font-semibold text-xs bg-slate-800 hover:bg-slate-700 text-slate-200 border border-slate-600 transition-colors flex items-center gap-1.5"
                  @click="forceResetCamera"
                >
                  <RotateCcw :size="13" /> Reset Kamera
                </button>
                <button
                  type="button"
                  class="py-1.5 px-3 rounded-lg font-semibold text-xs bg-emerald-600 hover:bg-emerald-500 text-white transition-colors flex items-center gap-1.5 ml-auto"
                  @click="switchTab('hardware')"
                >
                  <Barcode :size="13" /> Gunakan Mesin Barcode Scanner
                </button>
              </div>
            </div>
          </div>
        </div>

        <!-- TAB 3: 1-CLICK DEMO PATIENTS SIMULATION -->
        <div
          v-show="activeTab === 'demo'"
          class="bg-slate-800/90 rounded-2xl border border-slate-700/80 p-5 sm:p-6 shadow-xl space-y-5"
        >
          <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3 border-b border-slate-700/70 pb-4">
            <div>
              <h2 class="text-lg font-bold text-white flex items-center gap-2">
                <Sparkles :size="20" class="text-amber-400" />
                Simulasi Antrean Pasien (1-Click)
              </h2>
              <p class="text-xs text-slate-400 mt-1">
                Uji coba alur check-in langsung terhadap database Supabase PostgreSQL.
              </p>
            </div>

            <button
              type="button"
              class="inline-flex items-center gap-1.5 px-3.5 py-2 rounded-xl text-xs font-bold bg-slate-900 hover:bg-slate-700 text-slate-200 border border-slate-700 transition-colors"
              :disabled="resettingLoading"
              @click="resetAllDemoQueues"
            >
              <RefreshCw :size="14" :class="{ 'animate-spin': resettingLoading }" />
              Reset Status Semua ke 'Waiting'
            </button>
          </div>

          <!-- Demo Patients Grid -->
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-3.5">
            <div
              v-for="item in demoQueues"
              :key="item.qr_token"
              class="rounded-xl border p-4 flex flex-col justify-between transition-all bg-slate-900/90 hover:border-slate-600 shadow-md"
              :class="{
                'border-emerald-500/40 bg-emerald-950/20': item.status === 'checked_in',
                'border-blue-500/40 bg-blue-950/20': item.status === 'waiting',
                'border-amber-500/40 bg-amber-950/20': item.status === 'called',
                'border-purple-500/40 bg-purple-950/20': item.status === 'in_service'
              }"
            >
              <div class="flex items-start gap-3">
                <!-- QR Code Box -->
                <div class="w-16 h-16 shrink-0 rounded-xl bg-white p-1 border border-slate-200 flex items-center justify-center shadow-sm">
                  <img
                    v-if="qrCodeDataUrls[item.qr_token]"
                    :src="qrCodeDataUrls[item.qr_token]"
                    :alt="`QR ${item.patient_name}`"
                    class="w-full h-full object-contain"
                  />
                  <QrCode v-else :size="28" class="text-slate-800" />
                </div>

                <!-- Info -->
                <div class="min-w-0 flex-1">
                  <div class="flex items-center justify-between gap-1">
                    <span class="font-bold text-sm text-white truncate">{{ item.patient_name }}</span>
                    <span class="text-xs font-black px-1.5 py-0.5 rounded bg-slate-800 text-emerald-400 border border-slate-700">
                      {{ item.queue_code }}
                    </span>
                  </div>
                  <p class="text-xs text-slate-400 mt-0.5">{{ item.polyclinic_name }}</p>

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

              <!-- Action Buttons -->
              <div class="mt-3.5 pt-2.5 border-t border-slate-800 flex items-center gap-2">
                <button
                  type="button"
                  class="flex-1 py-2 px-3 rounded-xl text-xs font-bold transition-all flex items-center justify-center gap-1.5 shadow-md"
                  :class="isCooldown
                    ? 'bg-slate-800 text-slate-500 border border-slate-700 cursor-not-allowed'
                    : 'bg-emerald-600 hover:bg-emerald-500 text-white shadow-emerald-600/30'"
                  :disabled="verifyingToken === item.qr_token || isCooldown"
                  @click="handleCheckIn(item.qr_token, 'Simulasi ' + item.patient_name)"
                >
                  <Loader2 v-if="verifyingToken === item.qr_token" :size="13" class="animate-spin" />
                  <Check v-else :size="13" />
                  {{ isCooldown ? `Jeda (${cooldownRemaining}s)` : 'Uji Check-In' }}
                </button>

                <button
                  type="button"
                  class="p-2 rounded-xl text-xs font-medium text-slate-400 hover:text-white hover:bg-slate-800 border border-slate-700 transition-colors"
                  title="Kembalikan status pasien ini ke 'waiting'"
                  @click="resetSingleQueue(item.qr_token)"
                >
                  <RotateCcw :size="14" />
                </button>
              </div>
            </div>
          </div>
        </div>

        <!-- TAB 4: FILE UPLOAD QR IMAGE -->
        <div
          v-show="activeTab === 'upload'"
          class="bg-slate-800/90 rounded-2xl border border-slate-700/80 p-5 sm:p-6 shadow-xl space-y-5"
        >
          <div class="border-b border-slate-700/70 pb-4">
            <h2 class="text-lg font-bold text-white flex items-center gap-2">
              <FileImage :size="22" class="text-emerald-400" />
              Upload Berkas Gambar QR Tiket
            </h2>
            <p class="text-xs text-slate-400 mt-1">
              Unggah file gambar e-tiket atau struk antrean pasien untuk mendeteksi token QR secara otomatis.
            </p>
          </div>

          <div
            class="border-2 border-dashed rounded-2xl p-8 text-center transition-colors flex flex-col items-center justify-center gap-3 cursor-pointer bg-slate-900/60"
            :class="isDragging ? 'border-emerald-500 bg-emerald-950/20' : 'border-slate-700 hover:border-emerald-500/60'"
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
            <div class="p-3.5 rounded-full bg-emerald-500/20 text-emerald-400 border border-emerald-500/30">
              <UploadCloud :size="30" />
            </div>
            <div>
              <p class="text-sm font-bold text-white">
                Tarik file ke sini atau <span class="text-emerald-400 underline">pilih dari perangkat</span>
              </p>
              <p class="text-xs text-slate-400 mt-1">Mendukung file PNG, JPG, JPEG, WEBP (maksimal 5MB)</p>
            </div>
          </div>

          <div v-if="uploadedImagePreview" class="p-3 bg-slate-900 border border-slate-700 rounded-xl flex items-center justify-between gap-3">
            <div class="flex items-center gap-3 min-w-0">
              <img :src="uploadedImagePreview" alt="Preview" class="w-14 h-14 object-cover rounded-lg border border-slate-700 shrink-0 bg-white" />
              <div class="min-w-0">
                <p class="text-xs font-bold text-white truncate">{{ uploadedFileName }}</p>
                <p class="text-[11px] text-emerald-400 flex items-center gap-1 mt-0.5">
                  <CheckCircle2 :size="12" /> Berkas berhasil dianalisis dengan jsQR multi-pass
                </p>
              </div>
            </div>
            <button
              type="button"
              class="px-3 py-1.5 text-xs font-bold rounded-lg bg-slate-800 border border-slate-600 hover:bg-slate-700 text-slate-200 transition-colors"
              @click="triggerFileInput"
            >
              Ganti File
            </button>
          </div>

          <div v-if="fileScanning" class="flex items-center justify-center gap-2 py-4 text-emerald-400 text-xs font-semibold">
            <Loader2 :size="16" class="animate-spin" /> Menganalisis pola matriks QR...
          </div>
          <div v-else-if="fileError" class="text-xs text-rose-300 bg-rose-950/60 p-3 rounded-xl border border-rose-500/40 flex items-center gap-2">
            <AlertCircle :size="16" class="shrink-0 text-rose-400" /> {{ fileError }}
          </div>
        </div>

        <!-- TAB 5: MANUAL TOKEN INPUT -->
        <div
          v-show="activeTab === 'manual'"
          class="bg-slate-800/90 rounded-2xl border border-slate-700/80 p-5 sm:p-6 shadow-xl space-y-5"
        >
          <div class="border-b border-slate-700/70 pb-4">
            <h2 class="text-lg font-bold text-white flex items-center gap-2">
              <Keyboard :size="22" class="text-emerald-400" />
              Input Manual Kode Token QR
            </h2>
            <p class="text-xs text-slate-400 mt-1">
              Masukkan token string atau paste tautan check-in untuk pengujian langsung.
            </p>
          </div>

          <form @submit.prevent="handleManualSubmit" class="space-y-4">
            <div>
              <label for="manual-token" class="block text-xs font-semibold text-slate-300 mb-1.5 uppercase tracking-wider">
                String Token QR:
              </label>
              <div class="relative">
                <input
                  id="manual-token"
                  v-model="manualToken"
                  type="text"
                  required
                  placeholder="Contoh: qr-ahmad-fauzi-004 atau token UUID"
                  class="w-full rounded-xl bg-slate-950 border border-slate-700 py-2.5 pl-3.5 pr-10 text-sm font-mono text-emerald-300 focus:ring-2 focus:ring-emerald-500 focus:outline-none"
                />
                <button
                  v-if="manualToken"
                  type="button"
                  class="absolute right-3 top-1/2 -translate-y-1/2 text-slate-400 hover:text-slate-200"
                  @click="manualToken = ''"
                >
                  <X :size="16" />
                </button>
              </div>
            </div>

            <!-- Quick Chips -->
            <div class="flex flex-wrap items-center gap-1.5">
              <span class="text-xs text-slate-400">Pilih Cepat:</span>
              <button
                v-for="token in quickChips"
                :key="token.code"
                type="button"
                class="text-[11px] font-mono px-2 py-1 rounded-lg bg-slate-900 hover:bg-slate-700 text-slate-300 border border-slate-700 transition-colors"
                @click="manualToken = token.code"
              >
                {{ token.label }}
              </button>
            </div>

            <button
              type="submit"
              class="w-full py-3 px-4 rounded-xl font-bold text-sm transition-all shadow-lg flex items-center justify-center gap-2"
              :class="isCooldown
                ? 'bg-slate-800 text-slate-500 border border-slate-700 cursor-not-allowed'
                : 'bg-emerald-600 hover:bg-emerald-500 text-white shadow-emerald-600/30'"
              :disabled="verifyingToken === manualToken || !manualToken.trim() || isCooldown"
            >
              <Loader2 v-if="verifyingToken === manualToken" :size="16" class="animate-spin" />
              <CheckCircle2 v-else :size="16" />
              {{ isCooldown ? `Jeda Sensor (${cooldownRemaining}s)` : 'Verifikasi Check-In' }}
            </button>
          </form>
        </div>
      </section>

      <!-- RIGHT COLUMN: KIOSK RESULT CARD & INSPECTOR (5 COLS) -->
      <section class="lg:col-span-5 space-y-6">
        <!-- PRIMARY RESULT CARD -->
        <div class="bg-slate-800/90 rounded-2xl border border-slate-700/80 p-5 sm:p-6 shadow-xl space-y-4">
          <div class="flex items-center justify-between border-b border-slate-700/70 pb-3">
            <h3 class="font-bold text-white flex items-center gap-2 text-base">
              <Activity :size="18" class="text-emerald-400" />
              Hasil Pemindaian Kios
            </h3>
            <span v-if="lastResult" class="text-xs font-mono text-slate-400">
              {{ lastResultTime }} WIB
            </span>
          </div>

          <!-- STATE 1: STANDBY / WAITING FOR SCAN -->
          <div
            v-if="!lastResult && !verifyingToken"
            class="py-14 px-4 text-center text-slate-400 flex flex-col items-center justify-center gap-3"
          >
            <div class="p-4 rounded-2xl bg-slate-900 border border-slate-700 text-slate-500 relative">
              <ScanLine :size="40" class="text-emerald-400/80 animate-pulse" />
            </div>
            <div class="space-y-1">
              <p class="font-bold text-slate-200 text-base">Anjungan Siap Menerima Scan</p>
              <p class="text-xs text-slate-400 max-w-xs leading-relaxed">
                Arahkan lembar struk ke mesin barcode scanner USB atau kamera kios. Sistem akan merespon secara otomatis.
              </p>
            </div>
          </div>

          <!-- STATE 2: VERIFYING IN PROGRESS -->
          <div
            v-else-if="verifyingToken"
            class="py-14 px-4 text-center flex flex-col items-center justify-center gap-3 text-emerald-400"
          >
            <div class="relative">
              <div class="w-16 h-16 rounded-full border-4 border-emerald-950 border-t-emerald-400 animate-spin"></div>
              <div class="absolute inset-0 flex items-center justify-center">
                <QrCode :size="22" class="text-emerald-400" />
              </div>
            </div>
            <p class="font-bold text-white text-base">Memvalidasi Tiket Pasien...</p>
            <p class="text-xs text-slate-400 font-mono">Menghubungi database Supabase PostgreSQL...</p>
          </div>

          <!-- STATE 3A: SUCCESSFUL CHECK-IN (FIRST TIME) -->
          <div
            v-else-if="lastResult && lastResult.success && !lastResult.data?.already_checked_in"
            class="rounded-2xl border-2 border-emerald-500/80 bg-gradient-to-br from-emerald-950/70 via-slate-900 to-teal-950/70 p-5 sm:p-6 space-y-4 shadow-xl shadow-emerald-950/40"
          >
            <!-- Header Pill -->
            <div class="flex items-start gap-3">
              <div class="p-3 rounded-2xl bg-emerald-500 text-slate-950 shadow-lg shadow-emerald-500/30 shrink-0">
                <CheckCircle2 :size="28" />
              </div>
              <div class="min-w-0">
                <div class="inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full text-xs font-black bg-emerald-500 text-slate-950 uppercase tracking-wider">
                  Check-In Berhasil Diterima
                </div>
                <h4 class="text-xl font-black text-white mt-1 truncate">
                  {{ lastResultPatientName || 'Pasien Terdaftar' }}
                </h4>
              </div>
            </div>

            <!-- Big Ticket Number Box -->
            <div class="p-4 rounded-xl bg-slate-950/90 border border-emerald-500/40 space-y-3 shadow-inner">
              <div class="flex items-center justify-between">
                <div>
                  <p class="text-xs font-medium text-slate-400 uppercase tracking-wider">Nomor Antrean Anda</p>
                  <p class="text-4xl sm:text-5xl font-mono font-black text-emerald-400 tracking-tight drop-shadow-[0_0_12px_rgba(52,211,153,0.3)]">
                    {{ lastResult.data?.queue_code || 'A-001' }}
                  </p>
                </div>
                <div class="text-right">
                  <p class="text-xs font-medium text-slate-400 uppercase tracking-wider">Status Antrean</p>
                  <span class="inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-xs font-bold bg-emerald-500/20 text-emerald-300 border border-emerald-500/40">
                    <span class="w-2 h-2 rounded-full bg-emerald-400 animate-pulse"></span>
                    Checked-In
                  </span>
                </div>
              </div>

              <div class="grid grid-cols-2 gap-2 pt-3 border-t border-slate-800 text-xs">
                <div>
                  <span class="text-slate-400">Poliklinik Tujuan:</span>
                  <p class="font-bold text-white text-sm mt-0.5">{{ lastResultPolyclinic || 'Poli Umum' }}</p>
                </div>
                <div>
                  <span class="text-slate-400">Waktu Check-In:</span>
                  <p class="font-bold text-white text-sm mt-0.5">{{ lastResultCheckInTime }} WIB</p>
                </div>
              </div>

              <!-- Wait Time Estimate Info (EMA) -->
              <div v-if="queueEstimate" class="p-3 rounded-xl bg-emerald-950/60 border border-emerald-500/30 text-xs text-emerald-200 space-y-1">
                <div class="flex items-center justify-between font-bold">
                  <span>Estimasi Menunggu:</span>
                  <span class="text-emerald-400 font-mono text-sm">~{{ queueEstimate.estimate_minutes }} Menit</span>
                </div>
                <p class="text-[11px] text-emerald-300/80">
                  {{ queueEstimate.ahead_count }} pasien menunggu di depan antrean Anda.
                </p>
              </div>

              <!-- Next Steps Prompt -->
              <div class="p-2.5 rounded-lg bg-slate-900 border border-slate-800 text-[11px] text-slate-300 leading-relaxed">
                👉 <strong>Langkah Selanjutnya:</strong> Silakan simpan struk Anda dan langsung menuju ruang tunggu <strong>{{ lastResultPolyclinic }}</strong>.
              </div>
            </div>

            <!-- Action buttons -->
            <div class="flex gap-2">
              <RouterLink
                to="/staff"
                class="flex-1 py-2.5 px-3 rounded-xl text-xs font-bold bg-emerald-600 hover:bg-emerald-500 text-white transition-colors flex items-center justify-center gap-1.5 shadow-md shadow-emerald-600/20"
              >
                <Users :size="15" /> Lihat di Dashboard Staf
              </RouterLink>
              <button
                type="button"
                class="py-2.5 px-3 rounded-xl text-xs font-bold bg-slate-900 hover:bg-slate-700 text-slate-200 border border-slate-700 transition-colors"
                @click="resetLastResult"
              >
                Tutup
              </button>
            </div>
          </div>

          <!-- STATE 3B: ALREADY CHECKED-IN (PREVENT DOUBLE CHECK-IN) -->
          <div
            v-else-if="lastResult && lastResult.success && lastResult.data?.already_checked_in"
            class="rounded-2xl border-2 border-sky-400/80 bg-gradient-to-br from-sky-950/70 via-slate-900 to-indigo-950/70 p-5 sm:p-6 space-y-4 shadow-xl shadow-sky-950/40"
          >
            <!-- Header Pill -->
            <div class="flex items-start gap-3">
              <div class="p-3 rounded-2xl bg-sky-500 text-slate-950 shadow-lg shadow-sky-500/30 shrink-0">
                <CheckCheck :size="28" />
              </div>
              <div class="min-w-0">
                <div class="inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full text-xs font-black bg-sky-500 text-slate-950 uppercase tracking-wider">
                  Sudah Check-In Sebelumnya
                </div>
                <h4 class="text-xl font-black text-white mt-1 truncate">
                  {{ lastResultPatientName || 'Pasien AntriSehat' }}
                </h4>
              </div>
            </div>

            <!-- Big Ticket Number Box -->
            <div class="p-4 rounded-xl bg-slate-950/90 border border-sky-400/40 space-y-3 shadow-inner">
              <div class="flex items-center justify-between">
                <div>
                  <p class="text-xs font-medium text-slate-400 uppercase tracking-wider">Nomor Antrean Anda</p>
                  <p class="text-4xl sm:text-5xl font-mono font-black text-sky-400 tracking-tight drop-shadow-[0_0_12px_rgba(56,189,248,0.3)]">
                    {{ lastResult.data?.queue_code || 'A-001' }}
                  </p>
                </div>
                <div class="text-right">
                  <p class="text-xs font-medium text-slate-400 uppercase tracking-wider">Status Antrean</p>
                  <span class="inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-xs font-bold bg-sky-500/20 text-sky-300 border border-sky-400/40">
                    <span class="w-2 h-2 rounded-full bg-sky-400"></span>
                    {{ formatStatus(lastResult.data?.status) }}
                  </span>
                </div>
              </div>

              <div class="grid grid-cols-2 gap-2 pt-3 border-t border-slate-800 text-xs">
                <div>
                  <span class="text-slate-400">Poliklinik Tujuan:</span>
                  <p class="font-bold text-white text-sm mt-0.5">{{ lastResultPolyclinic || 'Poli Umum' }}</p>
                </div>
                <div>
                  <span class="text-slate-400">Waktu Check-In Awal:</span>
                  <p class="font-bold text-white text-sm mt-0.5">{{ lastResultCheckInTime }} WIB</p>
                </div>
              </div>

              <!-- Information Card on Duplicate Prevention -->
              <div class="p-3 rounded-xl bg-sky-950/60 border border-sky-400/30 text-xs text-sky-200 flex items-start gap-2.5">
                <Info :size="18" class="text-sky-400 shrink-0 mt-0.5" />
                <div class="leading-relaxed">
                  <strong class="font-bold text-white block">Mencegah Double Check-In:</strong>
                  Tiket ini telah divalidasi sebelumnya. Nomor antrean Anda tetap aktif dan tidak digandakan. Silakan langsung menuju ruang tunggu poli tujuan.
                </div>
              </div>
            </div>

            <div class="flex gap-2">
              <RouterLink
                to="/staff"
                class="flex-1 py-2.5 px-3 rounded-xl text-xs font-bold bg-sky-600 hover:bg-sky-500 text-white transition-colors flex items-center justify-center gap-1.5 shadow-md shadow-sky-600/20"
              >
                <Users :size="15" /> Buka di Dashboard Staf
              </RouterLink>
              <button
                type="button"
                class="py-2.5 px-3 rounded-xl text-xs font-bold bg-slate-900 hover:bg-slate-700 text-slate-200 border border-slate-700 transition-colors"
                @click="resetLastResult"
              >
                Selesai
              </button>
            </div>
          </div>

          <!-- STATE 4: CHECK-IN FAILED / RESTRICTED -->
          <div
            v-else-if="lastResult && !lastResult.success"
            class="rounded-2xl border-2 border-rose-500/80 bg-gradient-to-br from-rose-950/70 via-slate-900 to-orange-950/70 p-5 sm:p-6 space-y-4 shadow-xl shadow-rose-950/40"
          >
            <div class="flex items-start gap-3">
              <div class="p-3 rounded-2xl bg-rose-600 text-white shadow-lg shadow-rose-600/30 shrink-0">
                <XCircle :size="28" />
              </div>
              <div class="min-w-0">
                <span class="px-2.5 py-0.5 rounded-full text-xs font-bold bg-rose-600 text-white uppercase tracking-wider">
                  Check-In Tidak Berhasil
                </span>
                <h4 class="text-base font-bold text-white mt-1 leading-snug">
                  {{ lastResult.message || 'Verifikasi Tiket Gagal' }}
                </h4>
              </div>
            </div>

            <div class="p-4 rounded-xl bg-slate-950/90 border border-rose-500/40 text-xs space-y-2">
              <div class="flex items-center justify-between text-slate-400">
                <span>Kode Masalah:</span>
                <span class="font-mono font-bold text-rose-400">{{ lastResult.code || 'BAD_REQUEST' }}</span>
              </div>
              <div class="flex items-center justify-between text-slate-400">
                <span>Token Terpindai:</span>
                <span class="font-mono text-slate-300 truncate max-w-[200px]">{{ lastScannedToken }}</span>
              </div>

              <!-- Date Restriction Context -->
              <div
                v-if="lastResult.message?.includes('hari kunjungan') || lastResult.message?.includes('kedaluwarsa')"
                class="p-3 rounded-lg bg-amber-950/60 border border-amber-500/40 text-amber-200 text-xs leading-relaxed flex items-start gap-2.5 mt-2"
              >
                <CalendarDays :size="18" class="text-amber-400 shrink-0 mt-0.5" />
                <div>
                  <strong class="font-semibold block text-amber-300">Aturan Tanggal Tiket:</strong>
                  Kode QR tiket antrean dibuat khusus untuk sistem keamanan puskesmas dan <strong>hanya aktif pada hari kunjungan</strong> yang tertera di tiket.
                </div>
              </div>
            </div>

            <button
              type="button"
              class="w-full py-2.5 px-3 rounded-xl text-xs font-bold bg-slate-900 hover:bg-slate-700 text-slate-200 border border-slate-700 transition-colors"
              @click="resetLastResult"
            >
              Coba Token Lain
            </button>
          </div>
        </div>

        <!-- DEVELOPER JSON INSPECTOR -->
        <div class="bg-slate-800/90 rounded-2xl border border-slate-700/80 p-4 sm:p-5 shadow-xl space-y-3">
          <div class="flex items-center justify-between cursor-pointer" @click="showJsonInspector = !showJsonInspector">
            <h3 class="font-bold text-slate-200 flex items-center gap-2 text-xs uppercase tracking-wider">
              <Code :size="15" class="text-slate-400" />
              Developer Payload Inspector
            </h3>
            <span class="text-xs text-slate-400 flex items-center gap-1 font-mono">
              {{ showJsonInspector ? 'Sembunyikan' : 'Tampilkan' }}
              <ChevronDown :size="14" :class="{ 'rotate-180': showJsonInspector }" />
            </span>
          </div>

          <div v-show="showJsonInspector" class="space-y-2 pt-1">
            <div class="text-[11px] text-slate-400 flex justify-between font-mono">
              <span>Endpoint:</span>
              <span class="text-emerald-400">POST /api/v1/public/queues/checkin</span>
            </div>
            <pre class="bg-slate-950 text-emerald-400 p-3.5 rounded-xl text-[11px] font-mono overflow-x-auto max-h-52 leading-relaxed border border-slate-800">{{ formattedJsonResult }}</pre>
          </div>
        </div>

        <!-- SESSION AUDIT LOG -->
        <div class="bg-slate-800/90 rounded-2xl border border-slate-700/80 p-4 sm:p-5 shadow-xl space-y-3">
          <div class="flex items-center justify-between border-b border-slate-700/70 pb-2">
            <h3 class="font-bold text-slate-200 flex items-center gap-2 text-xs uppercase tracking-wider">
              <History :size="15" class="text-slate-400" />
              Riwayat Pemindaian Sesi Kios
            </h3>
            <span class="text-xs text-slate-400 font-mono">{{ scanHistory.length }} log</span>
          </div>

          <div v-if="scanHistory.length === 0" class="text-xs text-slate-500 text-center py-4">
            Belum ada riwayat pemindaian pada sesi ini.
          </div>
          <div v-else class="space-y-2 max-h-48 overflow-y-auto pr-1">
            <div
              v-for="(hist, idx) in scanHistory"
              :key="idx"
              class="flex items-center justify-between p-2 rounded-xl text-xs border"
              :class="hist.success
                ? (hist.alreadyCheckedIn ? 'bg-sky-950/40 border-sky-500/30' : 'bg-emerald-950/40 border-emerald-500/30')
                : 'bg-rose-950/40 border-rose-500/30'"
            >
              <div class="min-w-0 flex-1 pr-2">
                <div class="flex items-center gap-1.5">
                  <span
                    class="px-1.5 py-0.2 rounded text-[10px] font-bold uppercase"
                    :class="hist.success
                      ? (hist.alreadyCheckedIn ? 'bg-sky-600 text-white' : 'bg-emerald-600 text-white')
                      : 'bg-rose-600 text-white'"
                  >
                    {{ hist.success ? (hist.alreadyCheckedIn ? 'SUDAH CHECKIN' : 'BERHASIL') : 'GAGAL' }}
                  </span>
                  <span class="font-medium text-slate-200 truncate">{{ hist.source }}</span>
                </div>
                <p class="font-mono text-[10px] text-slate-400 truncate mt-0.5">{{ hist.token }}</p>
              </div>
              <div class="text-right text-[10px] text-slate-400 shrink-0 font-mono">
                <p>{{ hist.time }}</p>
                <p>{{ hist.durationMs }}ms</p>
              </div>
            </div>
          </div>
        </div>
      </section>
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
  AlertCircle, AlertTriangle, ScanLine, CalendarDays,
  Maximize, Minimize, Barcode, CheckCheck, Info, Zap, Hourglass, Cpu
} from "lucide-vue-next";
import { api } from "../../services/api";
import { useQueueStore } from "../../stores/queue";
import { decodeQrFromImageFile, extractCheckInToken } from "../../utils/qrReader";
import jsQR from "jsqr";

const router = useRouter();
const queueStore = useQueueStore();

// Polyclinic ID map for auto-registration
const polyclinicIdMap: Record<string, string> = {
  "Poli Umum": "11111111-1111-1111-1111-111111111111",
  "Poli Gigi": "22222222-2222-2222-2222-222222222222",
  "Poli KIA": "33333333-3333-3333-3333-333333333333",
  "Poli KIA / Anak": "33333333-3333-3333-3333-333333333333",
  "Poli Lansia": "44444444-4444-4444-4444-444444444444",
  "Poli Lansia / Geriatri": "44444444-4444-4444-4444-444444444444",
};

// Tabs Configuration
const tabs = [
  { id: "hardware", label: "Mesin Scan Fisik", icon: Barcode, badge: "USB / HID" },
  { id: "camera", label: "Kamera Kios", icon: Camera },
  { id: "demo", label: "Simulasi Pasien Demo", icon: Sparkles, badge: "1-Click" },
  { id: "upload", label: "Upload Berkas QR", icon: FileImage },
  { id: "manual", label: "Input Manual", icon: Keyboard },
];
const activeTab = ref("hardware");

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
const kioskContainerRef = ref<HTMLElement | null>(null);
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
const hardwareTestInput = ref("");
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

  // Barcode scanners type rapid keystrokes (< 80ms per char).
  // If the delay is more than 200ms, reset the buffer as it was human typing.
  if (timeSinceLastKey > 200) {
    barcodeKeystrokeBuffer = "";
  }

  // Active element check
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

  // Accumulate printable alphanumeric / symbol characters
  if (event.key.length === 1 && !event.ctrlKey && !event.altKey && !event.metaKey) {
    barcodeKeystrokeBuffer += event.key;
  }
}

function handleHardwareTestSubmit() {
  if (!hardwareTestInput.value.trim() || isCooldown.value) return;
  const token = hardwareTestInput.value.trim();
  hardwareTestInput.value = "";
  flashScannerSignal(token);
  handleCheckIn(token, "Uji Mesin Barcode Fisik");
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

// File Upload State
const fileInputRef = ref<HTMLInputElement | null>(null);
const isDragging = ref(false);
const fileScanning = ref(false);
const fileError = ref("");
const uploadedImagePreview = ref("");
const uploadedFileName = ref("");

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
  { label: "Ahmad Fauzi (Waiting)", code: "qr-ahmad-fauzi-004" },
  { label: "Dewi Lestari (Already Checked In)", code: "qr-dewi-lestari-003" },
  { label: "Budi Santoso (Called)", code: "qr-budi-santoso-002" },
  { label: "Siti Aminah (In Service)", code: "qr-siti-aminah-001" },
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
const showJsonInspector = ref(false);

// Session History Log
interface ScanHistoryItem {
  source: string;
  token: string;
  success: boolean;
  alreadyCheckedIn: boolean;
  time: string;
  durationMs: number;
}
const scanHistory = ref<ScanHistoryItem[]>([]);

const formattedJsonResult = computed(() => {
  if (!lastResult.value) return "// Belum ada respon scan.";
  return JSON.stringify(lastResult.value, null, 2);
});

// Sound Generator using Web Audio API
function playSound(type: "success" | "already_checked_in" | "error" | "scanner_beep") {
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
      osc.frequency.setValueAtTime(880, ctx.currentTime + 0.09); // A5
      gain.gain.setValueAtTime(0.25, ctx.currentTime);
      gain.gain.exponentialRampToValueAtTime(0.01, ctx.currentTime + 0.35);
      osc.start(ctx.currentTime);
      osc.stop(ctx.currentTime + 0.35);
    } else if (type === "already_checked_in") {
      // Friendly informational notification chord
      osc.type = "triangle";
      osc.frequency.setValueAtTime(440, ctx.currentTime); // A4
      osc.frequency.setValueAtTime(659.25, ctx.currentTime + 0.1); // E5
      gain.gain.setValueAtTime(0.2, ctx.currentTime);
      gain.gain.exponentialRampToValueAtTime(0.01, ctx.currentTime + 0.35);
      osc.start(ctx.currentTime);
      osc.stop(ctx.currentTime + 0.35);
    } else if (type === "scanner_beep") {
      // Fast crisp hardware beep
      osc.type = "sine";
      osc.frequency.setValueAtTime(1200, ctx.currentTime);
      gain.gain.setValueAtTime(0.15, ctx.currentTime);
      gain.gain.exponentialRampToValueAtTime(0.01, ctx.currentTime + 0.06);
      osc.start(ctx.currentTime);
      osc.stop(ctx.currentTime + 0.06);
    } else {
      // Low alert buzzer
      osc.type = "sawtooth";
      osc.frequency.setValueAtTime(220, ctx.currentTime);
      osc.frequency.setValueAtTime(160, ctx.currentTime + 0.1);
      gain.gain.setValueAtTime(0.25, ctx.currentTime);
      gain.gain.exponentialRampToValueAtTime(0.01, ctx.currentTime + 0.25);
      osc.start(ctx.currentTime);
      osc.stop(ctx.currentTime + 0.25);
    }
  } catch {
    // Audio context may require user interaction
  }
}

// Formatters
function formatStatus(status: string) {
  switch (status) {
    case "waiting": return "Menunggu (Waiting)";
    case "checked_in": return "Sudah Check-In";
    case "called": return "Dipanggil ke Ruang Poli";
    case "in_service": return "Sedang Dilayani Dokter";
    case "done": return "Pemeriksaan Selesai";
    default: return status || "Checked-In";
  }
}

function getStatusBadgeClass(status: string) {
  switch (status) {
    case "waiting": return "bg-blue-500/20 text-blue-300 border border-blue-500/30";
    case "checked_in": return "bg-emerald-500/20 text-emerald-300 border border-emerald-500/30";
    case "called": return "bg-amber-500/20 text-amber-300 border border-amber-500/30";
    case "in_service": return "bg-purple-500/20 text-purple-300 border border-purple-500/30";
    default: return "bg-slate-800 text-slate-300 border border-slate-700";
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
  const startTime = performance.now();

  try {
    let res = await api.checkInByQR(token);
    const duration = Math.round(performance.now() - startTime);

    // If check-in fails (e.g. 404 NOT_FOUND), check if this is an appointment in local store
    if (!res || !res.success) {
      const localPatient = queueStore.patients.find(
        (p) => p.qr_token === token || p.id === token || p.number === token,
      );

      if (localPatient) {
        // Auto-register to backend database on the fly!
        try {
          const polyId =
            polyclinicIdMap[localPatient.poli] ||
            "11111111-1111-1111-1111-111111111111";
          const guarantor =
            localPatient.insurance === "BPJS" ? "bpjs" : "umum";
          const gender =
            localPatient.gender === "Perempuan" ? "female" : "male";

          const regRes = await api.createAppointment({
            nik: localPatient.nik || "3201000000000001",
            full_name: localPatient.name || "Pasien",
            birth_date: localPatient.birth || "1995-01-01",
            gender,
            phone: localPatient.phone || "081234567890",
            guarantor,
            polyclinic_id: polyId,
            visit_date:
              localPatient.date || new Date().toISOString().slice(0, 10),
            chief_complaint:
              localPatient.complaint || "Pemeriksaan kesehatan",
            symptoms:
              localPatient.symptoms && localPatient.symptoms.length > 0
                ? localPatient.symptoms
                : ["Pemeriksaan"],
            duration_days: 1,
          });

          if (regRes && regRes.success && regRes.data) {
            localPatient.qr_token = regRes.data.qr_token;
            localPatient.number = regRes.data.queue.queue_code || localPatient.number;
            token = regRes.data.qr_token;
            res = await api.checkInByQR(token);
          }
        } catch (e) {
          console.warn("Auto-sync local patient failed:", e);
        }
      }
    }

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

      const demoMatched = demoQueues.value.find((q) => q.qr_token === token);
      const localMatched = queueStore.patients.find(
        (p) => p.qr_token === token || p.id === token || p.number === token,
      );

      // Also check localStorage
      let storedPatient: any = null;
      try {
        const storedStr = localStorage.getItem("antrisehat-patients");
        if (storedStr) {
          const list = JSON.parse(storedStr);
          storedPatient = list.find((p: any) => p.qr_token === token || p.id === token || p.number === token);
        }
      } catch (_) {}

      lastResultPatientName.value =
        demoMatched?.patient_name ||
        localMatched?.name ||
        storedPatient?.name ||
        (token.includes("576b758f") ? "dasdadas" : "Pasien Terverifikasi");

      lastResultPolyclinic.value =
        demoMatched?.polyclinic_name ||
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

      scanHistory.value.unshift({
        source,
        token,
        success: true,
        alreadyCheckedIn: isDuplicate,
        time: lastResultTime.value,
        durationMs: duration,
      });

      loadDemoQueues();
    } else {
      playSound("error");
      scanHistory.value.unshift({
        source,
        token,
        success: false,
        alreadyCheckedIn: false,
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
      alreadyCheckedIn: false,
      time: new Date().toLocaleTimeString("id-ID"),
      durationMs: duration,
    });
  } finally {
    verifyingToken.value = "";
    // ALWAYS start cooldown of 3 seconds after scan attempt to prevent double-scanning
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

      // IGNORE camera frames during cooldown to prevent double reads!
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

// ── Tab 4: Upload File QR ────────────────────────────────────
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
  uploadedFileName.value = file.name;

  try {
    uploadedImagePreview.value = URL.createObjectURL(file);
  } catch (_) {}

  try {
    let decodedText = await decodeQrFromImageFile(file);

    if (!decodedText) {
      try {
        if (!html5QrCode) {
          html5QrCode = new Html5Qrcode("qr-camera-reader");
        }
        decodedText = await html5QrCode.scanFile(file, true);
      } catch (_) {}
    }

    if (decodedText) {
      await handleCheckIn(decodedText, "Upload File (" + file.name + ")");
    } else {
      fileError.value =
        "Tidak menemukan kode QR yang jelas dalam gambar tersebut. Pastikan pola kotak QR terlihat kontras.";
      playSound("error");
    }
  } catch (err: any) {
    fileError.value =
      "Gagal memproses gambar: " + (err?.message || "Format file tidak didukung.");
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

// ── Tab 5: Manual Input ─────────────────────────────────────
function handleManualSubmit() {
  if (manualToken.value.trim() && !isCooldown.value) {
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
  updateClock();
  clockIntervalId = window.setInterval(updateClock, 1000);

  // Global Hardware Barcode Scanner listener
  window.addEventListener("keydown", handleGlobalKeyDown);
  document.addEventListener("fullscreenchange", onFullscreenChange);

  await loadDemoQueues();
  await refreshCameras();

  const params = new URLSearchParams(window.location.search);
  const urlToken = params.get("token") || params.get("checkin");
  if (urlToken) {
    switchTab("manual");
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
