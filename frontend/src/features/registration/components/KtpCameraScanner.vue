<template>
  <div class="relative w-full rounded-2xl overflow-hidden bg-slate-950 flex flex-col items-center justify-center min-h-[340px] sm:min-h-[400px] select-none shadow-inner border border-slate-800">
    <!-- 1. Video Stream Element -->
    <video
      ref="videoEl"
      autoplay
      playsinline
      muted
      class="w-full h-full object-cover absolute inset-0 transition-opacity duration-300"
      :class="isStreaming ? 'opacity-100' : 'opacity-0'"
    ></video>

    <!-- 2. Loading State saat Memulai Kamera -->
    <div v-if="isLoading" class="absolute inset-0 z-20 flex flex-col items-center justify-center bg-slate-950/90 text-white p-6 text-center">
      <Loader2 class="w-10 h-10 text-[var(--teal,#12b8b0)] animate-spin mb-3" />
      <p class="text-sm font-semibold text-slate-200">Menghubungkan ke Kamera...</p>
      <p class="text-xs text-slate-400 mt-1 max-w-xs">Izinkan akses kamera pada peramban Anda untuk pemindaian e-KTP otomatis.</p>
    </div>

    <!-- 3. Error State jika Izin Ditolak / Kamera Tidak Tersedia -->
    <div v-if="errorMessage" class="absolute inset-0 z-20 flex flex-col items-center justify-center bg-slate-900/95 text-white p-6 text-center">
      <div class="w-12 h-12 rounded-full bg-rose-500/20 text-rose-400 flex items-center justify-center mb-3">
        <AlertTriangle class="w-6 h-6" />
      </div>
      <h3 class="text-sm font-bold text-slate-100 mb-1">Kamera Tidak Dapat Diakses</h3>
      <p class="text-xs text-slate-400 max-w-sm mb-4 leading-relaxed">{{ errorMessage }}</p>
      <div class="flex items-center gap-3">
        <button
          type="button"
          @click="initCamera"
          class="px-4 py-2 bg-[var(--teal,#12b8b0)] hover:brightness-110 text-white text-xs font-semibold rounded-xl transition flex items-center gap-1.5 cursor-pointer shadow-md"
        >
          <RefreshCw class="w-3.5 h-3.5" />
          Coba Lagi
        </button>
        <button
          type="button"
          @click="$emit('switch-to-upload')"
          class="px-4 py-2 bg-slate-800 hover:bg-slate-700 text-slate-200 text-xs font-semibold rounded-xl transition flex items-center gap-1.5 cursor-pointer"
        >
          <UploadCloud class="w-3.5 h-3.5" />
          Gunakan Upload File
        </button>
      </div>
    </div>

    <!-- 4. Overlay Mask & Viewfinder Box KTP -->
    <div v-if="isStreaming" class="absolute inset-0 z-10 pointer-events-none flex flex-col items-center justify-center p-3 sm:p-4">
      <!-- Vignette Gelap di luar Frame -->
      <div class="absolute inset-0 bg-black/45 backdrop-blur-[1px]"></div>

      <!-- KTP Card Frame (Aspect Ratio Standar ID-1: ~1.586) -->
      <div
        ref="frameEl"
        class="relative w-[92%] max-w-[400px] aspect-[85.6/54] rounded-2xl border-2 transition-all duration-300 overflow-hidden pointer-events-auto shadow-2xl"
        :class="autoCatchSuccess ? 'border-emerald-400 ring-4 ring-emerald-400/50 shadow-[0_0_40px_rgba(16,185,129,0.7)]' : isDetecting ? 'border-amber-400 shadow-[0_0_30px_rgba(245,158,11,0.5)]' : 'border-[var(--teal,#12b8b0)] shadow-[0_0_30px_rgba(18,184,176,0.4)]'"
      >
        <!-- Corner Markers (Sudut Penunjuk KTP) -->
        <div
          class="absolute top-0 left-0 w-6 h-6 border-t-4 border-l-4 rounded-tl-xl transition-colors"
          :class="autoCatchSuccess ? 'border-emerald-300' : 'border-[var(--teal,#12b8b0)]'"
        ></div>
        <div
          class="absolute top-0 right-0 w-6 h-6 border-t-4 border-r-4 rounded-tr-xl transition-colors"
          :class="autoCatchSuccess ? 'border-emerald-300' : 'border-[var(--teal,#12b8b0)]'"
        ></div>
        <div
          class="absolute bottom-0 left-0 w-6 h-6 border-b-4 border-l-4 rounded-bl-xl transition-colors"
          :class="autoCatchSuccess ? 'border-emerald-300' : 'border-[var(--teal,#12b8b0)]'"
        ></div>
        <div
          class="absolute bottom-0 right-0 w-6 h-6 border-b-4 border-r-4 rounded-br-xl transition-colors"
          :class="autoCatchSuccess ? 'border-emerald-300' : 'border-[var(--teal,#12b8b0)]'"
        ></div>

        <!-- Laser Scan Animation (Berhenti ketika autoCatchSuccess) -->
        <div
          v-if="!autoCatchSuccess"
          class="laser-scanner absolute left-0 right-0 h-1 bg-gradient-to-r from-transparent via-[var(--teal,#12b8b0)] to-transparent shadow-[0_0_12px_#12b8b0]"
        ></div>

        <!-- Guidance Watermark di dalam Frame KTP -->
        <div class="absolute inset-0 flex flex-col items-center justify-between p-2.5 sm:p-3 text-center pointer-events-none">
          <div class="flex items-center gap-1.5 bg-slate-950/75 backdrop-blur-md px-2.5 py-1 rounded-lg border border-slate-700/80">
            <span class="w-2 h-2 rounded-full animate-ping" :class="autoCatchSuccess ? 'bg-emerald-400' : 'bg-[var(--teal,#12b8b0)]'"></span>
            <span class="text-[11px] font-bold tracking-wide uppercase text-teal-300">
              {{ autoCatchSuccess ? 'KTP Terbaca Jelas' : 'Auto-Read Aktif' }}
            </span>
          </div>

          <!-- Status deteksi dinamis -->
          <div class="bg-slate-950/85 backdrop-blur-md px-3 py-1.5 rounded-full border border-slate-700/80 max-w-[92%]">
            <p v-if="autoCatchSuccess" class="text-xs font-bold text-emerald-300 flex items-center gap-1.5 justify-center">
              <CheckCircle2 class="w-4 h-4 text-emerald-400" />
              Berhasil Ditangkap Otomatis!
            </p>
            <p v-else-if="detectedNikPreview" class="text-xs font-semibold text-amber-300 flex items-center gap-1.5 justify-center">
              <Loader2 class="w-3.5 h-3.5 animate-spin text-amber-400" />
              NIK Terdeteksi: {{ detectedNikPreview }}... Menstabilkan
            </p>
            <p v-else-if="isDetecting" class="text-xs font-medium text-slate-200 flex items-center gap-1.5 justify-center">
              <Loader2 class="w-3.5 h-3.5 animate-spin text-[var(--teal,#12b8b0)]" />
              Membaca kolom KTP... Jaga kamera stabil
            </p>
            <p v-else class="text-[11px] font-medium text-slate-300">
              Posisikan e-KTP pas di dalam bingkai (Otomatis membaca tanpa klik)
            </p>
          </div>
        </div>
      </div>

      <!-- Hint teks di bawah bingkai -->
      <p class="relative z-10 mt-3 text-xs text-white/90 font-medium text-center bg-slate-950/70 px-4 py-1 rounded-full backdrop-blur-md border border-slate-800 max-w-sm">
        {{ autoCatchSuccess ? 'Memproses data pasien...' : 'Sistem akan otomatis menangkap foto saat seluruh data terbaca jelas' }}
      </p>
    </div>

    <!-- 5. Top Controls Bar: Torch & Switch Camera -->
    <div v-if="isStreaming" class="absolute top-3 inset-x-3 z-20 flex items-center justify-between pointer-events-auto">
      <div class="flex items-center gap-1.5 bg-slate-900/85 backdrop-blur-md px-3 py-1 rounded-full border border-slate-700 text-xs text-slate-300 shadow-md">
        <span class="w-2 h-2 rounded-full bg-[var(--teal,#12b8b0)] animate-pulse"></span>
        <span class="font-bold text-[11px]">Kamera Live</span>
      </div>

      <div class="flex items-center gap-2">
        <!-- Tombol Flash / Torch jika disupport perangkat -->
        <button
          v-if="hasTorch"
          type="button"
          @click="toggleTorch"
          class="p-2 rounded-full backdrop-blur-md border border-slate-700 transition cursor-pointer shadow-md"
          :class="isTorchOn ? 'bg-amber-400 text-slate-950 shadow-amber-400/30' : 'bg-slate-900/80 text-white hover:bg-slate-800'"
          title="Nyalakan Lampu Flash"
        >
          <Zap v-if="!isTorchOn" class="w-4 h-4" />
          <ZapOff v-else class="w-4 h-4" />
        </button>

        <!-- Tombol Switch Camera (Depan / Belakang) -->
        <button
          v-if="hasMultipleCameras"
          type="button"
          @click="switchCamera"
          class="p-2 rounded-full bg-slate-900/80 hover:bg-slate-800 text-white backdrop-blur-md border border-slate-700 transition cursor-pointer shadow-md"
          title="Ganti Kamera"
        >
          <RefreshCw class="w-4 h-4" />
        </button>
      </div>
    </div>

    <!-- 6. Bottom Controls Bar: Manual Shutter (Sebagai Cadangan) -->
    <div v-if="isStreaming && !autoCatchSuccess" class="absolute bottom-3 inset-x-0 z-20 flex flex-col items-center justify-center pointer-events-auto px-4">
      <button
        type="button"
        @click="triggerManualCapture"
        class="group flex items-center gap-2 px-4 py-2 bg-slate-900/85 hover:bg-slate-800 active:scale-95 text-white/90 hover:text-white font-semibold text-xs rounded-full border border-slate-700 backdrop-blur-md shadow-lg transition cursor-pointer"
      >
        <Camera class="w-3.5 h-3.5 text-[var(--teal,#12b8b0)] group-hover:scale-110 transition" />
        <span>Ambil Manual</span>
      </button>
    </div>

    <!-- Shutter Flash Effect saat Sukses Ditangkap -->
    <div
      v-if="showFlash"
      class="absolute inset-0 z-30 bg-white transition-opacity duration-200 pointer-events-none"
    ></div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, nextTick } from 'vue'
import {
  Camera,
  RefreshCw,
  Zap,
  ZapOff,
  AlertTriangle,
  Loader2,
  UploadCloud,
  CheckCircle2
} from 'lucide-vue-next'
import { cropSourceToCanvas } from '@/utils/imageProcessor'
import { isAutoCatchReady, type KtpData } from '@/utils/ktpParser'
import { processMultiEnhancementOcr } from '@/utils/ocrService'

const emit = defineEmits<{
  (e: 'auto-captured', payload: {
    canvas: HTMLCanvasElement
    parsed: KtpData
    rawText: string
    winningMethod?: string
    winningMethodLabel?: string
    processedDataUrl?: string
  }): void
  (e: 'switch-to-upload'): void
}>()

const videoEl = ref<HTMLVideoElement | null>(null)
const frameEl = ref<HTMLDivElement | null>(null)

const isStreaming = ref(false)
const isLoading = ref(true)
const errorMessage = ref('')
const showFlash = ref(false)

const currentStream = ref<MediaStream | null>(null)
const facingMode = ref<'environment' | 'user'>('environment')
const hasMultipleCameras = ref(false)
const hasTorch = ref(false)
const isTorchOn = ref(false)

// State Auto-Scan
const isDetecting = ref(false)
const autoCatchSuccess = ref(false)
const detectedNikPreview = ref('')
let autoScanTimer: ReturnType<typeof setInterval> | null = null
let isFrameProcessing = false

/**
 * Suara Beep Sukses (Web Audio API Synthesizer - tanpa file eksternal)
 */
const playSuccessSound = () => {
  try {
    const AudioCtx = window.AudioContext || (window as unknown as { webkitAudioContext: typeof AudioContext }).webkitAudioContext
    if (!AudioCtx) return
    const ctx = new AudioCtx()
    const osc = ctx.createOscillator()
    const gain = ctx.createGain()

    osc.type = 'sine'
    osc.frequency.setValueAtTime(587.33, ctx.currentTime) // D5
    osc.frequency.setValueAtTime(880, ctx.currentTime + 0.1) // A5

    gain.gain.setValueAtTime(0.2, ctx.currentTime)
    gain.gain.exponentialRampToValueAtTime(0.001, ctx.currentTime + 0.35)

    osc.connect(gain)
    gain.connect(ctx.destination)

    osc.start()
    osc.stop(ctx.currentTime + 0.35)
  } catch (e) {
    console.warn('Audio feedback error:', e)
  }
}

/**
 * Mengambil cuplikan hanya pada area kotak kartu KTP (Viewfinder Frame)
 */
const getCroppedCardCanvas = (): HTMLCanvasElement | null => {
  if (!videoEl.value || !isStreaming.value) return null

  const video = videoEl.value
  const frame = frameEl.value

  if (!frame) {
    const canvas = document.createElement('canvas')
    canvas.width = video.videoWidth
    canvas.height = video.videoHeight
    const ctx = canvas.getContext('2d')
    if (ctx) ctx.drawImage(video, 0, 0)
    return canvas
  }

  const videoRect = video.getBoundingClientRect()
  const frameRect = frame.getBoundingClientRect()

  const xPercent = Math.max(0, (frameRect.left - videoRect.left) / videoRect.width)
  const yPercent = Math.max(0, (frameRect.top - videoRect.top) / videoRect.height)
  const widthPercent = Math.min(1 - xPercent, frameRect.width / videoRect.width)
  const heightPercent = Math.min(1 - yPercent, frameRect.height / videoRect.height)

  return cropSourceToCanvas(video, {
    xPercent,
    yPercent,
    widthPercent,
    heightPercent
  })
}

/**
 * Loop Pemindai Otomatis (Auto-Scan)
 */
const startAutoScanLoop = () => {
  stopAutoScanLoop()

  autoScanTimer = setInterval(async () => {
    if (isFrameProcessing || !isStreaming.value || autoCatchSuccess.value) return

    const croppedCanvas = getCroppedCardCanvas()
    if (!croppedCanvas || croppedCanvas.width < 100 || croppedCanvas.height < 100) return

    isFrameProcessing = true
    isDetecting.value = true

    try {
      const multiResult = await processMultiEnhancementOcr(croppedCanvas, undefined, {
        targetWidth: 1200,
        fastAutoScan: true
      })

      const winner = multiResult.winner
      const parsed = winner.parsed

      if (parsed.nik && parsed.nik.length >= 8) {
        detectedNikPreview.value = `${parsed.nik.slice(0, 10)}... (${winner.methodLabel})`
      } else {
        detectedNikPreview.value = ''
      }

      if (isAutoCatchReady(parsed)) {
        autoCatchSuccess.value = true
        stopAutoScanLoop()

        showFlash.value = true
        playSuccessSound()
        setTimeout(() => {
          showFlash.value = false
        }, 150)

        setTimeout(() => {
          emit('auto-captured', {
            canvas: croppedCanvas,
            parsed,
            rawText: winner.rawText,
            winningMethod: winner.method,
            winningMethodLabel: winner.methodLabel,
            processedDataUrl: winner.dataUrl
          })
        }, 350)
      }
    } catch (e) {
      console.warn('Auto-scan frame error:', e)
    } finally {
      isFrameProcessing = false
      isDetecting.value = false
    }
  }, 1000)
}

const stopAutoScanLoop = () => {
  if (autoScanTimer) {
    clearInterval(autoScanTimer)
    autoScanTimer = null
  }
}

/**
 * Trigger manual jika user memilih klik tombol
 */
const triggerManualCapture = async () => {
  if (isFrameProcessing || autoCatchSuccess.value) return

  const croppedCanvas = getCroppedCardCanvas()
  if (!croppedCanvas) return

  stopAutoScanLoop()
  autoCatchSuccess.value = true
  showFlash.value = true
  playSuccessSound()
  setTimeout(() => {
    showFlash.value = false
  }, 150)

  try {
    const multiResult = await processMultiEnhancementOcr(croppedCanvas, undefined, {
      targetWidth: 1400,
      fastAutoScan: false
    })

    const winner = multiResult.winner

    emit('auto-captured', {
      canvas: croppedCanvas,
      parsed: winner.parsed,
      rawText: winner.rawText,
      winningMethod: winner.method,
      winningMethodLabel: winner.methodLabel,
      processedDataUrl: winner.dataUrl
    })
  } catch (e) {
    console.error('Manual capture error:', e)
  }
}

/**
 * Inisialisasi Akses Kamera
 */
const initCamera = async () => {
  isLoading.value = true
  errorMessage.value = ''
  autoCatchSuccess.value = false
  detectedNikPreview.value = ''
  stopCamera()

  try {
    if (!navigator.mediaDevices || !navigator.mediaDevices.getUserMedia) {
      throw new Error('Browser Anda tidak mendukung akses kamera langsung. Silakan gunakan opsi Upload Berkas.')
    }

    try {
      const devices = await navigator.mediaDevices.enumerateDevices()
      const videoDevices = devices.filter((d) => d.kind === 'videoinput')
      hasMultipleCameras.value = videoDevices.length > 1
    } catch {
      hasMultipleCameras.value = false
    }

    let stream: MediaStream
    try {
      stream = await navigator.mediaDevices.getUserMedia({
        video: {
          facingMode: { ideal: facingMode.value },
          width: { ideal: 1920, min: 1280 },
          height: { ideal: 1080, min: 720 }
        },
        audio: false
      })
    } catch {
      stream = await navigator.mediaDevices.getUserMedia({
        video: true,
        audio: false
      })
    }

    currentStream.value = stream

    await nextTick()
    if (videoEl.value) {
      videoEl.value.srcObject = stream
      await videoEl.value.play()
      isStreaming.value = true

      const track = stream.getVideoTracks()[0]
      if (track && typeof track.getCapabilities === 'function') {
        const capabilities = track.getCapabilities() as { torch?: boolean }
        hasTorch.value = !!capabilities.torch
      }

      startAutoScanLoop()
    }
  } catch (err: unknown) {
    console.error('Camera init error:', err)
    const errObj = err as { name?: string; message?: string }
    if (errObj.name === 'NotAllowedError' || errObj.name === 'PermissionDeniedError') {
      errorMessage.value = 'Izin kamera ditolak. Berikan izin kamera di browser Anda untuk menggunakan fitur ini.'
    } else if (errObj.name === 'NotFoundError' || errObj.name === 'DevicesNotFoundError') {
      errorMessage.value = 'Tidak ditemukan perangkat kamera yang terpasang.'
    } else {
      errorMessage.value = errObj.message || 'Gagal menyalakan kamera. Silakan periksa koneksi kamera Anda.'
    }
    isStreaming.value = false
  } finally {
    isLoading.value = false
  }
}

const toggleTorch = async () => {
  if (!currentStream.value) return
  const track = currentStream.value.getVideoTracks()[0]
  if (!track) return

  try {
    const nextState = !isTorchOn.value
    await track.applyConstraints({
      advanced: [{ torch: nextState } as MediaTrackConstraintSet]
    })
    isTorchOn.value = nextState
  } catch (e) {
    console.warn('Torch not supported:', e)
  }
}

const switchCamera = async () => {
  facingMode.value = facingMode.value === 'environment' ? 'user' : 'environment'
  await initCamera()
}

const stopCamera = () => {
  stopAutoScanLoop()
  if (currentStream.value) {
    currentStream.value.getTracks().forEach((track) => track.stop())
    currentStream.value = null
  }
  if (videoEl.value) {
    videoEl.value.srcObject = null
  }
  isStreaming.value = false
  isTorchOn.value = false
}

onMounted(() => {
  initCamera()
})

onUnmounted(() => {
  stopCamera()
})
</script>

<style scoped>
@keyframes scanMove {
  0% {
    top: 6%;
    opacity: 0.85;
  }
  50% {
    top: 90%;
    opacity: 1;
  }
  100% {
    top: 6%;
    opacity: 0.85;
  }
}

.laser-scanner {
  animation: scanMove 2.2s ease-in-out infinite;
}
</style>
