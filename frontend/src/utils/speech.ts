// ============================================================
// MODUL AUDIO & TRANSKRIPSI PANGGILAN ANTREAN (BAHASA INDONESIA)
// ============================================================

/**
 * Konversi nomor antrean ke pembacaan bahasa Indonesia yang alami & jelas.
 * Contoh: "A-001" -> "A, nol nol satu"
 *         "B-012" -> "B, nol satu dua"
 *         "C-105" -> "C, satu nol lima"
 */
export function formatQueueForSpeech(queueNumber: string): string {
  if (!queueNumber) return "";
  const parts = queueNumber.split("-");
  if (parts.length === 2) {
    const letter = parts[0].trim().toUpperCase();
    const numPart = parts[1].trim();

    const digitWords: Record<string, string> = {
      "0": "nol",
      "1": "satu",
      "2": "dua",
      "3": "tiga",
      "4": "empat",
      "5": "lima",
      "6": "enam",
      "7": "tujuh",
      "8": "delapan",
      "9": "sembilan",
    };

    const digitsSpoken = numPart
      .split("")
      .map((d) => digitWords[d] || d)
      .join(" ");

    return `${letter}, ${digitsSpoken}`;
  }
  return queueNumber.replace("-", " ");
}

/**
 * Normalisasi nama poliklinik agar diucapkan jernih tanpa karakter slash/simbol.
 * Contoh: "Poli KIA / Anak" -> "Poli KIA Anak"
 */
export function cleanPoliNameForSpeech(poliName: string): string {
  if (!poliName) return "Poliklinik";
  return poliName
    .replace(/\//g, " ")
    .replace(/-/g, " ")
    .replace(/\s+/g, " ")
    .trim();
}

/**
 * Menghasilkan nada lonceng panggilan (Ding-Dong 2-nada) menggunakan Web Audio API.
 */
export function playClinicChime(): Promise<void> {
  return new Promise((resolve) => {
    try {
      if (typeof window === "undefined") {
        resolve();
        return;
      }
      const AudioCtx =
        window.AudioContext ||
        (window as unknown as { webkitAudioContext: typeof AudioContext }).webkitAudioContext;
      if (!AudioCtx) {
        resolve();
        return;
      }

      const ctx = new AudioCtx();
      if (ctx.state === "suspended") {
        ctx.resume().catch(() => {});
      }

      const now = ctx.currentTime;

      // Nada 1: Sol / G5 (784 Hz) - Ding
      const osc1 = ctx.createOscillator();
      const gain1 = ctx.createGain();
      osc1.type = "sine";
      osc1.frequency.setValueAtTime(784, now);
      gain1.gain.setValueAtTime(0.18, now);
      gain1.gain.exponentialRampToValueAtTime(0.001, now + 0.4);
      osc1.connect(gain1);
      gain1.connect(ctx.destination);
      osc1.start(now);
      osc1.stop(now + 0.42);

      // Nada 2: Mi / E5 (659.25 Hz) - Dong
      const osc2 = ctx.createOscillator();
      const gain2 = ctx.createGain();
      osc2.type = "sine";
      osc2.frequency.setValueAtTime(659.25, now + 0.22);
      gain2.gain.setValueAtTime(0.2, now + 0.22);
      gain2.gain.exponentialRampToValueAtTime(0.001, now + 0.65);
      osc2.connect(gain2);
      gain2.connect(ctx.destination);
      osc2.start(now + 0.22);
      osc2.stop(now + 0.68);

      setTimeout(() => {
        resolve();
      }, 600);
    } catch {
      resolve();
    }
  });
}

/**
 * Mencari suara Bahasa Indonesia resmi di browser pengguna.
 */
function getIndonesianVoice(): SpeechSynthesisVoice | null {
  if (typeof window === "undefined" || !("speechSynthesis" in window)) return null;
  const voices = window.speechSynthesis.getVoices();
  if (!voices || voices.length === 0) return null;

  // 1. Suara resmi id-ID atau id_ID
  const exact = voices.find(
    (v) =>
      v.lang === "id-ID" ||
      v.lang === "id_ID" ||
      v.lang.toLowerCase() === "id"
  );
  if (exact) return exact;

  // 2. Suara dengan awalan id-
  const prefix = voices.find(
    (v) =>
      v.lang.toLowerCase().startsWith("id-") ||
      v.lang.toLowerCase().startsWith("id_")
  );
  if (prefix) return prefix;

  // 3. Suara dengan nama mengandung Indonesia (misal: 'Google Bahasa Indonesia')
  const byName = voices.find(
    (v) =>
      v.name.toLowerCase().includes("indonesia") ||
      v.name.toLowerCase().includes("indonesian")
  );
  if (byName) return byName;

  return null;
}

// Inisialisasi daftar suara saat browser selesai memuat daftar TTS
if (typeof window !== "undefined" && "speechSynthesis" in window) {
  window.speechSynthesis.onvoiceschanged = () => {
    getIndonesianVoice();
  };
}

// Variable untuk melacak audio panggilan aktif
let activeCallAudio: HTMLAudioElement | null = null;

/**
 * Memutar audio TTS dengan logat Bahasa Indonesia asli (aksen alami Indonesia, tanpa aksen bule).
 */
function playNativeIndonesianAudio(text: string): Promise<boolean> {
  return new Promise((resolve) => {
    try {
      if (typeof window === "undefined") {
        resolve(false);
        return;
      }

      if (activeCallAudio) {
        activeCallAudio.pause();
        activeCallAudio.currentTime = 0;
        activeCallAudio = null;
      }

      const encoded = encodeURIComponent(text);
      const ttsUrl = `https://translate.google.com/translate_tts?ie=UTF-8&tl=id&client=tw-ob&q=${encoded}`;
      const audio = new Audio(ttsUrl);
      activeCallAudio = audio;

      let hasResolved = false;

      audio.onended = () => {
        activeCallAudio = null;
        if (!hasResolved) {
          hasResolved = true;
          resolve(true);
        }
      };

      audio.onerror = () => {
        activeCallAudio = null;
        if (!hasResolved) {
          hasResolved = true;
          resolve(false);
        }
      };

      // Timeout pengaman bila koneksi lambat
      const timer = setTimeout(() => {
        if (!hasResolved) {
          hasResolved = true;
          resolve(false);
        }
      }, 7000);

      audio.play().then(() => {
        clearTimeout(timer);
      }).catch(() => {
        clearTimeout(timer);
        if (!hasResolved) {
          hasResolved = true;
          resolve(false);
        }
      });
    } catch {
      resolve(false);
    }
  });
}

/**
 * Memutar panggilan suara antrean dalam logat dan pelafalan bahasa Indonesia asli.
 */
export async function speakQueueCall(queueNumber: string, poliName: string) {
  if (typeof window === "undefined") return;

  try {
    // 1. Hentikan panggilan suara sebelumnya jika masih berjalan
    if (activeCallAudio) {
      activeCallAudio.pause();
      activeCallAudio.currentTime = 0;
      activeCallAudio = null;
    }
    if ("speechSynthesis" in window) {
      window.speechSynthesis.cancel();
    }

    // 2. Mainkan nada lonceng panggilan klinik (Ding-Dong) terlebih dahulu
    await playClinicChime();

    const formattedNumber = formatQueueForSpeech(queueNumber);
    const cleanPoli = cleanPoliNameForSpeech(poliName);

    // Kalimat panggilan standar dengan jeda koma agar ritme bicara santun & jelas
    const textToSpeak = `Nomor antrean, ${formattedNumber}, silakan menuju, ${cleanPoli}`;

    // 3. Utama: Gunakan audio logat Indonesia asli (Native Indonesian Voice)
    const success = await playNativeIndonesianAudio(textToSpeak);

    // 4. Fallback: Jika offline atau diblokir peramban, gunakan Web Speech API dengan suara id-ID
    if (!success && "speechSynthesis" in window) {
      const speech = new SpeechSynthesisUtterance(textToSpeak);
      speech.lang = "id-ID";
      speech.rate = 0.88;
      speech.pitch = 1.0;

      const indonesianVoice = getIndonesianVoice();
      if (indonesianVoice) {
        speech.voice = indonesianVoice;
      }

      window.speechSynthesis.speak(speech);
    }
  } catch (err) {
    console.warn("Gagal memutar audio panggilan antrean:", err);
  }
}


