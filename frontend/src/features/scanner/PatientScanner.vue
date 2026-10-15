<template>
  <main class="container py-8 sm:py-12">
    <div class="mx-auto max-w-2xl">
      <button
        class="back-link !inline-flex items-center gap-2"
        @click="goToStaff"
      >
        <ChevronRight class="rotate-180" :size="18" /> Kembali ke Dashboard Staf
      </button>

      <section class="dashboard-panel">
        <div class="mb-6 flex items-start gap-4">
          <span class="service-icon shrink-0">
            <QrCode :size="24" />
          </span>
          <div>
            <h2 class="!text-2xl !tracking-tight">Scan QR Pasien</h2>
            <p class="mt-2">
              Pindai QR pada e-tiket untuk membuka verifikasi pasien.
            </p>
          </div>
        </div>

        <div
          class="relative aspect-video overflow-hidden rounded-xl bg-[var(--deep)]"
        >
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
            <span>{{
              running ? "Menyiapkan kamera…" : "Kamera belum aktif"
            }}</span>
          </div>
          <div
            v-if="ready"
            class="pointer-events-none absolute inset-0 flex items-center justify-center"
          >
            <div class="aspect-square w-1/2 rounded-xl border-2 border-white" />
          </div>
        </div>

        <p class="mt-3 text-sm">
          Posisikan QR di tengah kamera. Kamera berhenti setelah QR valid
          terbaca atau saat Anda meninggalkan halaman.
        </p>

        <button
          type="button"
          class="button primary full-width"
          :disabled="running && !ready"
          @click="toggleCamera"
        >
          <X v-if="running" :size="18" />
          <Monitor v-else :size="18" />
          {{ running ? "Hentikan kamera" : "Aktifkan kamera" }}
        </button>

        <p v-if="error" class="error mt-4" role="alert">
          {{ error }}
        </p>

        <form
          class="mt-6 border-t border-[var(--line)] pt-5"
          @submit.prevent="handleManualSubmit"
        >
          <label>
            Atau masukkan kode e-tiket secara manual:
            <input
              required
              v-model="manualCode"
              placeholder="Contoh: demo-0 atau ID tiket"
            />
          </label>
          <button class="button outline full-width" type="submit">
            Verifikasi Kode
          </button>
        </form>
      </section>
    </div>
  </main>
</template>

<script setup lang="ts">
import { ref, onUnmounted } from "vue";
import { useRouter } from "vue-router";
import { ChevronRight, QrCode, Monitor, X } from "lucide-vue-next";
import { useQueueStore } from "../../stores/queue";
import { todayDate } from "../../constants/initialData";

type Detector = {
  detect: (video: HTMLVideoElement) => Promise<{ rawValue: string }[]>;
};
type DetectorConstructor = new (options: { formats: string[] }) => Detector;

const router = useRouter();
const queueStore = useQueueStore();

const videoRef = ref<HTMLVideoElement | null>(null);
const running = ref(false);
const ready = ref(false);
const error = ref("");
const manualCode = ref("");

let stopped = false;
let currentStream: MediaStream | null = null;
let detectTimer: ReturnType<typeof setTimeout> | null = null;
let lastScannedValue = "";

function stopCamera() {
  stopped = true;
  if (detectTimer) clearTimeout(detectTimer);
  if (currentStream) {
    currentStream.getTracks().forEach((track) => track.stop());
    currentStream = null;
  }
  ready.value = false;
  running.value = false;
}

async function startCamera() {
  stopped = false;
  error.value = "";
  running.value = true;

  try {
    const Constructor = (
      window as unknown as { BarcodeDetector?: DetectorConstructor }
    ).BarcodeDetector;

    if (!Constructor) {
      throw new Error(
        "Browser ini belum mendukung pemindaian QR kamera. Gunakan Chrome di Android atau masukkan kode e-tiket manual di bawah.",
      );
    }

    if (!navigator.mediaDevices?.getUserMedia) {
      throw new Error(
        "Kamera tidak tersedia. Gunakan input kode e-tiket di bawah.",
      );
    }

    const detector = new Constructor({ formats: ["qr_code"] });
    currentStream = await navigator.mediaDevices.getUserMedia({
      video: { facingMode: { ideal: "environment" } },
      audio: false,
    });

    if (stopped) {
      currentStream.getTracks().forEach((track) => track.stop());
      return;
    }

    if (!videoRef.value) {
      stopCamera();
      return;
    }

    videoRef.value.srcObject = currentStream;
    await videoRef.value.play();

    if (stopped) return;
    ready.value = true;

    const detectLoop = async () => {
      if (stopped || !videoRef.value) return;
      try {
        const results = await detector.detect(videoRef.value);
        if (stopped) return;
        const val = results[0]?.rawValue;
        if (val && val !== lastScannedValue) {
          lastScannedValue = val;
          const errMsg = processScannedValue(val);
          if (!errMsg) {
            stopCamera();
            return;
          }
          error.value = errMsg;
        }
        if (!val) lastScannedValue = "";
      } catch {
        if (!stopped) {
          error.value =
            "QR belum terbaca. Pastikan kode jelas dan pencahayaan cukup.";
        }
      }
      if (!stopped) {
        detectTimer = setTimeout(detectLoop, 350);
      }
    };

    void detectLoop();
  } catch (err: unknown) {
    stopCamera();
    if (err instanceof Error) {
      if (
        ["NotAllowedError", "NotFoundError", "NotReadableError"].includes(
          err.name,
        )
      ) {
        error.value =
          "Kamera tidak dapat diakses. Periksa izin kamera atau masukkan kode e-tiket.";
      } else {
        error.value = err.message;
      }
    } else {
      error.value = "Gagal mengakses kamera.";
    }
  }
}

function toggleCamera() {
  if (running.value) {
    stopCamera();
  } else {
    startCamera();
  }
}

function processScannedValue(value: string): string | null {
  let id = value.trim();
  try {
    const url = new URL(id);
    if (
      url.origin !== window.location.origin ||
      url.pathname !== window.location.pathname
    ) {
      return "QR bukan e-tiket AntriSehat dari aplikasi ini.";
    }
    id = url.searchParams.get("checkin") || "";
  } catch {
    if (!/^[a-zA-Z0-9-]+$/.test(id)) {
      return "Kode e-tiket tidak valid.";
    }
  }

  const patient = queueStore.patients.find((p) => p.id === id);
  if (!patient) {
    return "Pasien tidak ditemukan pada data perangkat ini.";
  }
  if (patient.date !== todayDate) {
    return "E-tiket ini bukan untuk kunjungan hari ini.";
  }
  if (!["Terdaftar", "Menunggu verifikasi"].includes(patient.status)) {
    return "Pasien sudah diverifikasi atau kunjungannya tidak lagi aktif.";
  }

  queueStore.verify = patient;
  router.push("/staff");
  return null;
}

function handleManualSubmit() {
  const errMsg = processScannedValue(manualCode.value);
  if (errMsg) {
    error.value = errMsg;
  } else {
    stopCamera();
  }
}

function goToStaff() {
  stopCamera();
  router.push("/staff");
}

onUnmounted(() => {
  stopCamera();
});
</script>
