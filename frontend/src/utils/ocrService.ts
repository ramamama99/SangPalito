import Tesseract, { type Worker } from 'tesseract.js'
import { processImageForOcr, type PreprocessOptions } from './imageProcessor'
import { parseKtpText, isAutoCatchReady, type KtpData } from './ktpParser'

let workerInstance: Worker | null = null
let isInitializing = false
let initPromise: Promise<Worker> | null = null

export type OcrProgressCallback = (progress: number, status: string) => void

export type EnhancementMethod = 'adaptive' | 'binarize' | 'contrast' | 'original'

export interface EnhancementMethodDef {
  id: EnhancementMethod
  label: string
  description: string
  options: PreprocessOptions
}

export const ENHANCEMENT_METHODS: EnhancementMethodDef[] = [
  {
    id: 'adaptive',
    label: 'Otomatis (Optimal)',
    description: 'Algoritma pembersih latar belakang dipadukan dengan filter penajaman teks',
    options: { mode: 'adaptive', contrast: 50, sharpen: true, targetWidth: 1400 }
  },
  {
    id: 'binarize',
    label: 'Binarisasi (Otsu)',
    description: 'Konversi Otsu Thresholding teks hitam latar belakang putih',
    options: { mode: 'binarize', sharpen: false, targetWidth: 1400 }
  },
  {
    id: 'contrast',
    label: 'Kontras Tinggi',
    description: 'Konversi grayscale dengan peningkatan tingkat kontras maksimal',
    options: { mode: 'contrast', contrast: 90, sharpen: false, targetWidth: 1400 }
  },
  {
    id: 'original',
    label: 'Foto Asli',
    description: 'Pemrosesan tangkapan kamera murni tanpa tambahan filter',
    options: { mode: 'original', targetWidth: 1400 }
  }
]

export interface EnhancementCandidate {
  method: EnhancementMethod
  methodLabel: string
  methodDescription: string
  dataUrl: string
  rawText: string
  ocrConfidence: number
  parsed: KtpData
  isReady: boolean
  score: number
}

export interface MultiEnhancementResult {
  winner: EnhancementCandidate
  candidates: EnhancementCandidate[]
  hasConfidenceAbove90: boolean
}

/**
 * Menginisialisasi Tesseract Worker secara singleton
 * Menggunakan model dwibahasa 'ind+eng' untuk akurasi maksimal pada karakter KTP
 */
export async function getOcrWorker(onProgress?: OcrProgressCallback): Promise<Worker> {
  if (workerInstance) {
    return workerInstance
  }

  if (isInitializing && initPromise) {
    return initPromise
  }

  isInitializing = true
  initPromise = (async () => {
    try {
      if (onProgress) onProgress(10, 'Mengunduh kamus OCR Bahasa Indonesia & Inggris...')

      // Inisialisasi Tesseract v5 worker
      const worker = await Tesseract.createWorker(['ind', 'eng'], 1, {
        logger: (m) => {
          if (onProgress && m.status) {
            let userStatus = m.status
            if (m.status === 'loading tesseract core') userStatus = 'Memuat modul AI OCR...'
            else if (m.status === 'initializing tesseract') userStatus = 'Menginisialisasi engine...'
            else if (m.status === 'loading language traineddata') userStatus = 'Memuat model bahasa KTP...'
            else if (m.status === 'recognizing text') userStatus = 'Mengenali teks KTP...'

            const pct = Math.round((m.progress || 0) * 100)
            onProgress(pct, userStatus)
          }
        }
      })

      // Konfigurasi Parameter untuk Teks KTP
      await worker.setParameters({
        tessedit_pageseg_mode: Tesseract.PSM.AUTO, // PSM 3: Otomatis deteksi tata letak
        preserve_interword_spaces: '1'
      })

      workerInstance = worker
      return worker
    } catch (error) {
      console.error('Gagal membuat Tesseract Worker:', error)
      throw error
    } finally {
      isInitializing = false
    }
  })()

  return initPromise
}

/**
 * Menjalankan OCR pada gambar (Data URL, Canvas, atau File)
 */
export async function recognizeKtp(
  imageSource: string | HTMLCanvasElement,
  onProgress?: OcrProgressCallback
): Promise<{ text: string; confidence: number }> {
  const worker = await getOcrWorker(onProgress)

  if (onProgress) onProgress(30, 'Memulai pengenalan karakter OCR...')

  const result = await worker.recognize(imageSource, {}, {
    text: true
  })

  return {
    text: result.data.text || '',
    confidence: result.data.confidence || 0
  }
}

/**
 * Menghitung skor kelayakan kandidat hasil ekstraksi
 */
function computeCandidateScore(parsed: KtpData, ocrConfidence: number, method: EnhancementMethod): number {
  let score = parsed.confidence * 10
  if (parsed.isValidNik && parsed.nik.length === 16) score += 400
  if (parsed.nama && parsed.nama.length >= 2) score += 200
  if (parsed.jenisKelamin) score += 150
  if (parsed.tglLahir) score += 150
  if (isAutoCatchReady(parsed)) score += 1000
  score += ocrConfidence * 0.5
  if (method === 'adaptive') score += 15
  if (method === 'binarize') score += 10
  return score
}

/**
 * 1. Optimasi Kualitas Gambar & Penajaman Teks Otomatis (Auto-Enhancement)
 * Sistem memproses gambar di latar belakang dengan membandingkan hasil bacaan dari 4 metode berikut:
 * a. Otomatis (Optimal): Algoritma pembersih latar belakang dipadukan dengan filter penajaman teks.
 * b. Binarisasi (Otsu): Konversi gambar menjadi format Otsu Thresholding teks hitam latar belakang putih.
 * c. Kontras Tinggi: Konversi gambar ke mode grayscale dengan peningkatan tingkat kontras maksimal.
 * d. Foto Asli: Pemrosesan menggunakan tangkapan kamera murni tanpa tambahan filter.
 * Lalu secara otomatis memilih hasil ekstraksi terbaik dengan confidence level di atas 90%.
 */
export async function processMultiEnhancementOcr(
  imageSource: HTMLCanvasElement | File | Blob | string,
  onProgress?: (progress: number, status: string) => void,
  options: {
    rotation?: number
    targetWidth?: number
    fastAutoScan?: boolean
  } = {}
): Promise<MultiEnhancementResult> {
  const { rotation = 0, targetWidth = 1400, fastAutoScan = false } = options
  const candidates: EnhancementCandidate[] = []

  for (let idx = 0; idx < ENHANCEMENT_METHODS.length; idx++) {
    const methodDef = ENHANCEMENT_METHODS[idx]
    const methodProgressStart = Math.round((idx / ENHANCEMENT_METHODS.length) * 100)

    if (onProgress) {
      onProgress(
        Math.max(10, methodProgressStart),
        `Mengevaluasi metode (${idx + 1}/4): ${methodDef.label}...`
      )
    }

    try {
      const processed = await processImageForOcr(imageSource, {
        ...methodDef.options,
        rotation,
        targetWidth
      })

      const ocrResult = await recognizeKtp(processed.dataUrl)
      const parsed = parseKtpText(ocrResult.text, ocrResult.confidence)
      const isReady = isAutoCatchReady(parsed)
      const score = computeCandidateScore(parsed, ocrResult.confidence, methodDef.id)

      const candidate: EnhancementCandidate = {
        method: methodDef.id,
        methodLabel: methodDef.label,
        methodDescription: methodDef.description,
        dataUrl: processed.dataUrl,
        rawText: ocrResult.text,
        ocrConfidence: ocrResult.confidence,
        parsed,
        isReady,
        score
      }

      candidates.push(candidate)

      // Pada mode real-time auto-scan: jika metode pertama (adaptive) sudah mencapai
      // confidence >= 90% dan memenuhi seluruh 5 kriteria stop condition,
      // return langsung tanpa jeda demi performa live camera yang instan!
      if (fastAutoScan && isReady && parsed.confidence >= 90) {
        return {
          winner: candidate,
          candidates,
          hasConfidenceAbove90: true
        }
      }
    } catch (err) {
      console.warn(`Gagal memproses metode OCR ${methodDef.id}:`, err)
    }
  }

  if (candidates.length === 0) {
    throw new Error('Gagal memproses gambar dengan seluruh metode auto-enhancement.')
  }

  // Urutkan kandidat berdasarkan skor tertinggi
  candidates.sort((a, b) => b.score - a.score)
  const winner = candidates[0]
  const hasConfidenceAbove90 = winner.parsed.confidence >= 90

  if (onProgress) {
    onProgress(100, `Ekstraksi terbaik: ${winner.methodLabel} (${winner.parsed.confidence}%)`)
  }

  return {
    winner,
    candidates,
    hasConfidenceAbove90
  }
}

/**
 * Hentikan worker ketika tidak dibutuhkan lagi
 */
export async function terminateOcrWorker(): Promise<void> {
  if (workerInstance) {
    await workerInstance.terminate()
    workerInstance = null
    initPromise = null
  }
}

