<script setup lang="ts">
import { ref, watch, onUnmounted } from 'vue'
import { ChevronRight, QrCode, Monitor, X } from 'lucide-vue-next'

type Detector = { detect: (video: HTMLVideoElement) => Promise<{ rawValue: string }[]> }
type DetectorConstructor = new (options: { formats: string[] }) => Detector

const props = defineProps<{
  onScan: (value: string) => string | null
}>()

const emit = defineEmits<{
  (e: 'back'): void
}>()

const videoRef = ref<HTMLVideoElement | null>(null)
const running = ref(false)
const ready = ref(false)
const error = ref('')
const code = ref('')

let stopped = false
let stream: MediaStream | undefined
let timer: ReturnType<typeof setTimeout> | undefined
let lastValue = ''

const stopCamera = () => {
  stopped = true
  if (timer) clearTimeout(timer)
  stream?.getTracks().forEach((track) => track.stop())
  stream = undefined
}

const startCamera = async () => {
  try {
    stopped = false
    const Constructor = (window as unknown as { BarcodeDetector?: DetectorConstructor })
      .BarcodeDetector

    if (!Constructor) {
      throw new Error(
        'Browser ini belum mendukung pemindaian QR. Gunakan Chrome di Android atau masukkan kode e-tiket di bawah.',
      )
    }
    if (!navigator.mediaDevices?.getUserMedia) {
      throw new Error('Kamera tidak tersedia. Gunakan input kode e-tiket di bawah.')
    }

    const detector = new Constructor({ formats: ['qr_code'] })
    stream = await navigator.mediaDevices.getUserMedia({
      video: { facingMode: { ideal: 'environment' } },
      audio: false,
    })

    if (stopped) {
      stream.getTracks().forEach((track) => track.stop())
      return
    }

    if (!videoRef.value) {
      stopCamera()
      return
    }

    videoRef.value.srcObject = stream
    await videoRef.value.play()

    if (stopped) return
    ready.value = true

    const detect = async () => {
      if (stopped || !videoRef.value) return
      try {
        const results = await detector.detect(videoRef.value)
        if (stopped) return
        const value = results[0]?.rawValue
        if (value && value !== lastValue) {
          lastValue = value
          const message = props.onScan(value)
          if (!message) {
            stopCamera()
            running.value = false
            ready.value = false
            return
          }
          error.value = message
        }
        if (!value) lastValue = ''
      } catch {
        if (!stopped) error.value = 'QR belum terbaca. Pastikan kode jelas dan pencahayaan cukup.'
      }
      if (!stopped) timer = setTimeout(detect, 350)
    }

    void detect()
  } catch (cause) {
    stopCamera()
    ready.value = false
    running.value = false
    error.value =
      cause instanceof Error &&
      !['NotAllowedError', 'NotFoundError', 'NotReadableError'].includes(cause.name)
        ? cause.message
        : 'Kamera tidak dapat diakses. Periksa izin kamera atau masukkan kode e-tiket.'
  }
}

watch(running, (isRunning) => {
  if (isRunning) {
    startCamera()
  } else {
    stopCamera()
    ready.value = false
  }
})

onUnmounted(() => {
  stopCamera()
})

const toggleCamera = () => {
  error.value = ''
  ready.value = false
  running.value = !running.value
}

const handleManualSubmit = () => {
  const message = props.onScan(code.value)
  if (message) {
    error.value = message
  } else {
    stopCamera()
    running.value = false
    ready.value = false
  }
}
</script>

<template>
  <main class="container py-8 sm:py-12">
    <div class="mx-auto max-w-2xl">
      <button class="back-link !inline-flex items-center gap-2" @click="emit('back')">
        <ChevronRight class="rotate-180" :size="18" /> Kembali ke Dashboard Staf
      </button>

      <section class="dashboard-panel">
        <div class="mb-6 flex items-start gap-4">
          <span class="service-icon shrink-0"><QrCode :size="24" /></span>
          <div>
            <h2 class="!text-2xl !tracking-tight">Scan QR Pasien</h2>
            <p class="mt-2">Pindai QR pada e-tiket untuk membuka verifikasi pasien.</p>
          </div>
        </div>

        <div class="relative aspect-video overflow-hidden rounded-xl bg-[var(--deep)]">
          <video
            ref="videoRef"
            autoplay
            muted
            playsinline
            class="h-full w-full object-cover"
            aria-label="Pratinjau kamera pemindai"
          />
          <div
            v-if="!ready"
            class="absolute inset-0 flex flex-col items-center justify-center gap-3 p-5 text-center text-white"
          >
            <QrCode :size="48" />
            <span>{{ running ? 'Menyiapkan kamera…' : 'Kamera belum aktif' }}</span>
          </div>
          <div
            v-if="ready"
            class="pointer-events-none absolute inset-0 flex items-center justify-center"
          >
            <div class="aspect-square w-1/2 rounded-xl border-2 border-white" />
          </div>
        </div>

        <p class="mt-3 text-sm">
          Posisikan QR di tengah kamera. Kamera berhenti setelah QR valid terbaca atau saat Anda
          meninggalkan halaman.
        </p>

        <button
          type="button"
          class="button primary full-width"
          :disabled="running && !ready"
          @click="toggleCamera"
        >
          <X v-if="running" :size="18" />
          <Monitor v-else :size="18" />
          {{ running ? 'Hentikan kamera' : 'Aktifkan kamera' }}
        </button>

        <p v-if="error" class="error mt-4" role="alert">{{ error }}</p>

        <form class="mt-6 border-t border-[var(--line)] pt-5" @submit.prevent="handleManualSubmit">
          <label>
            Kode e-tiket / tautan QR
            <input
              required
              v-model="code"
              placeholder="Masukkan kode e-tiket atau tautan QR"
              autocomplete="off"
            />
          </label>
          <small class="mt-2 block text-[var(--muted)]">
            Alternatif jika kamera tidak tersedia.
          </small>
          <button class="button outline full-width" type="submit">
            <QrCode :size="18" /> Cari pasien
          </button>
        </form>
      </section>
    </div>
  </main>
</template>
