export interface KtpData {
  nik: string
  nama: string
  tempatLahir: string      // Kota kelahiran (Hanya huruf, bersih tanpa simbol/karakter aneh)
  tglLahir: string         // Format: DD-MM-YYYY
  tempatTglLahir: string   // Format gabungan: "JAKARTA, 19-05-2006"
  jenisKelamin: string     // "LAKI-LAKI" | "PEREMPUAN"
  golDarah?: string        // "A" | "B" | "AB" | "O" | "-"
  alamatJalan: string      // a. Alamat (Jalan / Nomor / Blok)
  rt: string               // b. RT (3 digit)
  rw: string               // c. RW (3 digit)
  kelDesa: string          // d. Kel/Desa
  kecamatan: string        // e. Kecamatan
  confidence: number       // 0-100%
  isValidNik: boolean      // Tepat 16 digit angka valid
  validationErrors?: string[]
}

/**
 * Aturan Validasi Data NIK:
 * a. Wajib Angka (Numerik): Koreksi otomatis jika mengandung huruf seperti O/o menjadi 0 dan I/l menjadi 1.
 */
export function cleanCharForDigit(str: string): string {
  if (!str) return ''
  return str
    .replace(/[OoDQUC]/g, '0')
    .replace(/[Il|!\]\[iT]/g, '1')
    .replace(/[Zz]/g, '2')
    .replace(/[E]/g, '3')
    .replace(/[Ah+]/g, '4')
    .replace(/[Ss]/g, '5')
    .replace(/[Gb]/g, '6')
    .replace(/[B&]/g, '8')
    .replace(/[gq]/g, '9')
}

/**
 * Filter Pembersihan Teks:
 * Memastikan tidak ada karakter aneh atau simbol yang masuk ke dalam form,
 * kecuali huruf, angka, tanda titik (.), tanda garis miring (/), dan spasi.
 */
export function sanitizeCleanText(val: string): string {
  if (!val) return ''
  return val
    .replace(/[^A-Za-z0-9./\s]/g, '')
    .replace(/\s{2,}/g, ' ')
}

/**
 * Filter Pembersihan NIK:
 * Hanya mengizinkan digit angka numerik maksimal 16 digit.
 */
export function sanitizeNik(val: string): string {
  if (!val) return ''
  return cleanCharForDigit(val).replace(/\D/g, '').slice(0, 16)
}

/**
 * Pembersih karakter teks: hanya huruf A-Z dan spasi (tanpa karakter aneh/simbol/angka)
 */
export function cleanLettersOnly(val: string): string {
  if (!val) return ''
  return val
    .replace(/[^A-Za-z\s]/g, ' ')
    .replace(/\s+/g, ' ')
    .trim()
    .toUpperCase()
}

/**
 * Aturan Validasi NIK:
 * a. Wajib Angka (Numerik): Nilai NIK harus berupa angka secara keseluruhan dengan koreksi otomatis huruf.
 * b. Panjang Karakter: Validasi dinyatakan berhasil jika NIK berjumlah tepat 16 digit angka numerik.
 */
export function validateNik(nik: string): { isValid: boolean; message?: string } {
  if (!nik) {
    return { isValid: false, message: 'NIK tidak boleh kosong' }
  }
  const cleaned = cleanCharForDigit(nik).replace(/\D/g, '')
  if (cleaned.length !== 16) {
    return {
      isValid: false,
      message: `Panjang NIK saat ini ${cleaned.length} digit. Validasi berhasil jika NIK berjumlah tepat 16 digit angka.`
    }
  }
  if (!/^\d{16}$/.test(cleaned)) {
    return { isValid: false, message: 'Nilai NIK harus berupa angka numerik secara keseluruhan.' }
  }

  return { isValid: true }
}

/**
 * Normalisasi singkatan jalan (Jalan / Jl.)
 */
function normalizeStreetPrefix(val: string): string {
  if (!val) return ''
  let cleaned = val
    .replace(/^[:;=~\-_.,*'"|/\\+?^#@!\s]+/, '')
    .replace(/[:;=~\-_.,*'"|/\\+?^#@!\s]+$/, '')
    .trim()
  cleaned = cleaned.replace(/^(?:JUN|TLN|JIN|JLN|JALAN|JLAN|ILN|1LN|LLN|JL)\b[\s.:-]*/i, 'Jl. ')
  cleaned = cleaned.replace(/\s{2,}/g, ' ').trim()
  return sanitizeCleanText(cleaned)
}

/**
 * Parser Cerdas KTP dengan Pemetaan Kolom Berdasarkan Kata Kunci (Keyword Mapping Case-Insensitive):
 * a. Kolom Tempat/Tgl Lahir: Terdeteksi kata Tempat atau Lahir.
 * b. Kolom Jenis Kelamin: Terdeteksi kata Jenis atau Kelamin.
 * c. Kolom Alamat: Terdeteksi suku kata Ala atau mat.
 * d. Kolom RT/RW: Terdeteksi kata RT atau RW.
 * e. Kolom Kelurahan/Desa: Terdeteksi kata KEL atau DESA.
 * f. Kolom Kecamatan: Terdeteksi suku kata Keca atau matan.
 */
export function parseKtpText(rawText: string, ocrRawConfidence: number = 0): KtpData {
  const result: KtpData = {
    nik: '',
    nama: '',
    tempatLahir: '',
    tglLahir: '',
    tempatTglLahir: '',
    jenisKelamin: '',
    golDarah: '',
    alamatJalan: '',
    rt: '',
    rw: '',
    kelDesa: '',
    kecamatan: '',
    confidence: 0,
    isValidNik: false,
    validationErrors: []
  }

  if (!rawText || rawText.trim().length === 0) {
    return result
  }

  const lines = rawText.split(/\r?\n/).map((l) => l.trim()).filter(Boolean)

  for (let i = 0; i < lines.length; i++) {
    const rawLine = lines[i]

    // Deteksi jika satu baris memuat dua key-value sekaligus (misal: Jenis Kelamin & Gol. Darah)
    const segments: string[] = []
    const splitMulti = rawLine.split(/(?=\b(?:Gol\.?\s*Darah|Darah)\s*[:;=])/i)
    for (const seg of splitMulti) {
      if (seg.trim()) segments.push(seg.trim())
    }

    for (const seg of segments) {
      let key = ''
      let value = ''

      const colonIndex = seg.indexOf(':')
      if (colonIndex !== -1) {
        key = seg.substring(0, colonIndex).trim().toUpperCase()
        value = seg.substring(colonIndex + 1).trim()
      } else {
        const altMatch = seg.match(/^([^;=|!]+)[;=|!]\s*(.+)$/)
        if (altMatch) {
          key = altMatch[1].trim().toUpperCase()
          value = altMatch[2].trim()
        } else {
          // Cari pemisahan kata kunci utama di awal baris
          const knownMatch = seg.match(/^(NIK|NAMA|TEMPAT\/TGL\s*LAHIR|TEMPAT\s*LAHIR|TGL\s*LAHIR|JENIS\s*KELAMIN|ALAMAT|RT\/RW|KEL(?:\/DESA)?|DESA|KECAMATAN|GOL\.?\s*DARAH)\b[\s:;=|.-]*(.*)$/i)
          if (knownMatch) {
            key = knownMatch[1].trim().toUpperCase()
            value = (knownMatch[2] || '').trim()
          } else {
            const firstSpace = seg.indexOf(' ')
            if (firstSpace !== -1) {
              key = seg.substring(0, firstSpace).trim().toUpperCase()
              value = seg.substring(firstSpace + 1).trim()
            }
          }
        }
      }

      if (!value && i + 1 < lines.length && !lines[i + 1].includes(':')) {
        value = lines[i + 1].trim()
      }

      const keyUpper = key.toUpperCase()

      // --- 1. NIK: Wajib Angka & Tepat 16 Digit ---
      if (keyUpper.includes('NIK') || keyUpper.includes('N!K') || keyUpper.includes('N1K') || keyUpper.includes('NLK')) {
        const cleanedDigits = cleanCharForDigit(value).replace(/\D/g, '')
        const match16 = cleanedDigits.match(/\d{16}/)
        if (match16) {
          result.nik = match16[0]
        } else if (cleanedDigits.length >= 16) {
          result.nik = cleanedDigits.slice(0, 16)
        } else if (cleanedDigits.length > 0) {
          result.nik = cleanedDigits
        }
      }

      // --- 2. NAMA LENGAP ---
      else if (!keyUpper.includes('PROVINSI') && !keyUpper.includes('KABUPATEN') && !keyUpper.includes('KOTA') && keyUpper.includes('NAMA')) {
        result.nama = sanitizeCleanText(cleanLettersOnly(value))
      }

      // --- 3.a. Kolom Tempat/Tgl Lahir: Terdeteksi kata Tempat atau Lahir ---
      else if (keyUpper.includes('TEMPAT') || keyUpper.includes('LAHIR') || keyUpper.includes('TEMPAY') || keyUpper.includes('TGL')) {
        const dateMatch = value.match(/\b(\d{1,2})[-/. ](\d{1,2})[-/. ](\d{4})\b/)
        if (dateMatch) {
          result.tglLahir = `${dateMatch[1].padStart(2, '0')}/${dateMatch[2].padStart(2, '0')}/${dateMatch[3]}`
          const cityPart = value.substring(0, value.indexOf(dateMatch[0]))
          result.tempatLahir = sanitizeCleanText(cleanLettersOnly(cityPart))
        } else {
          if (/\d/.test(value)) {
            result.tglLahir = cleanCharForDigit(value).replace(/[^0-9/.]/g, '')
          } else {
            result.tempatLahir = sanitizeCleanText(cleanLettersOnly(value))
          }
        }
        if (result.tempatLahir && result.tglLahir) {
          result.tempatTglLahir = `${result.tempatLahir} / ${result.tglLahir}`
        }
      }

      // --- 3.b. Kolom Jenis Kelamin: Terdeteksi kata Jenis atau Kelamin ---
      else if (keyUpper.includes('JENIS') || keyUpper.includes('KELAMIN')) {
        const upVal = value.toUpperCase()
        if (upVal.includes('LAKI')) {
          result.jenisKelamin = 'LAKI-LAKI'
        } else if (upVal.includes('PEREMPUAN') || upVal.includes('WANITA')) {
          result.jenisKelamin = 'PEREMPUAN'
        }
      }

      // --- GOLONGAN DARAH ---
      else if (keyUpper.includes('DARAH') || keyUpper.includes('GOL')) {
        const upVal = value.toUpperCase()
        const bloodMatch = upVal.match(/\b(AB|A|B|O|-)[+]?\b/i)
        if (bloodMatch) {
          result.golDarah = bloodMatch[1].toUpperCase()
        }
      }

      // --- 3.f. Kolom Kecamatan: Terdeteksi suku kata Keca atau matan ---
      else if (keyUpper.includes('KECA') || keyUpper.includes('MATAN') || keyUpper.includes('KEC')) {
        result.kecamatan = sanitizeCleanText(cleanLettersOnly(value))
      }

      // --- 3.e. Kolom Kelurahan/Desa: Terdeteksi kata KEL atau DESA ---
      else if (keyUpper.includes('KEL') || keyUpper.includes('DESA')) {
        result.kelDesa = sanitizeCleanText(cleanLettersOnly(value))
      }

      // --- 3.d. Kolom RT/RW: Terdeteksi kata RT atau RW ---
      else if (keyUpper.includes('RT') || keyUpper.includes('RW')) {
        const cleanedRtRw = cleanCharForDigit(value)
        const matchSlash = cleanedRtRw.match(/(\d{1,4})\s*[/\\|:-]\s*(\d{1,4})/)
        if (matchSlash) {
          result.rt = matchSlash[1].padStart(3, '0').slice(-3)
          result.rw = matchSlash[2].padStart(3, '0').slice(-3)
        } else {
          const matchSpace = cleanedRtRw.match(/(\d{1,4})\s+(\d{1,4})/)
          if (matchSpace) {
            result.rt = matchSpace[1].padStart(3, '0').slice(-3)
            result.rw = matchSpace[2].padStart(3, '0').slice(-3)
          }
        }
      }

      // --- 3.c. Kolom Alamat: Terdeteksi suku kata Ala atau mat (bukan dari Tempat atau Kecamatan) ---
      else if (
        !keyUpper.includes('TEMPAT') &&
        !keyUpper.includes('LAHIR') &&
        !keyUpper.includes('KEC') &&
        !keyUpper.includes('MATAN') &&
        (keyUpper.includes('ALA') || keyUpper.includes('MAT') || keyUpper.includes('ALMT'))
      ) {
        result.alamatJalan = normalizeStreetPrefix(value)
      }
    }
  }

  // Fallback NIK jika tidak menggunakan format titik dua
  if (!result.nik || result.nik.length < 16) {
    for (const line of lines) {
      const cleaned = cleanCharForDigit(line).replace(/\D/g, '')
      const match16 = cleaned.match(/\d{16}/)
      if (match16) {
        result.nik = match16[0]
        break
      }
    }
  }

  // Validasi NIK (Wajib Numerik & Tepat 16 Digit)
  if (result.nik) {
    const nikCheck = validateNik(result.nik)
    result.isValidNik = nikCheck.isValid
    if (!nikCheck.isValid && nikCheck.message) {
      result.validationErrors?.push(nikCheck.message)
    }
  }

  // Fallback Tanggal Lahir dari NIK jika belum terdeteksi dari teks
  if (!result.tglLahir && result.nik && result.nik.length === 16) {
    const tgl = parseInt(result.nik.slice(6, 8), 10)
    const isFemale = tgl > 40
    const realTgl = isFemale ? tgl - 40 : tgl
    const bln = result.nik.slice(8, 10)
    const thn = result.nik.slice(10, 12)
    const fullYear = parseInt(thn, 10) > 30 ? `19${thn}` : `20${thn}`
    if (realTgl >= 1 && realTgl <= 31 && parseInt(bln, 10) >= 1 && parseInt(bln, 10) <= 12) {
      result.tglLahir = `${realTgl.toString().padStart(2, '0')}/${bln}/${fullYear}`
      if (result.tempatLahir) {
        result.tempatTglLahir = `${result.tempatLahir} / ${result.tglLahir}`
      } else {
        result.tempatTglLahir = result.tglLahir
      }
    }
  }

  // Fallback Jenis Kelamin dari NIK jika belum terisi
  if (!result.jenisKelamin && result.nik && result.nik.length === 16) {
    const tgl = parseInt(result.nik.slice(6, 8), 10)
    if (!isNaN(tgl)) {
      result.jenisKelamin = tgl > 40 ? 'PEREMPUAN' : 'LAKI-LAKI'
    }
  }

  // Hitung Skor Akurasi (Confidence Level)
  let score = 0

  // a. NIK valid 16 digit angka: +30%
  if (result.isValidNik && result.nik && result.nik.length === 16 && /^\d{16}$/.test(result.nik)) {
    score += 30
  }

  // b. Nama terisi: +25%
  if (result.nama && result.nama.length >= 2) {
    score += 25
  }

  // c. Jenis Kelamin terisi: +15%
  if (result.jenisKelamin && (result.jenisKelamin === 'LAKI-LAKI' || result.jenisKelamin === 'PEREMPUAN')) {
    score += 15
  }

  // d. Tanggal Lahir terisi (DD/MM/YYYY atau DD-MM-YYYY): +15%
  if (result.tglLahir && /^\d{2}[-/.]\d{2}[-/.]\d{4}$/.test(result.tglLahir)) {
    score += 15
  }

  // e. Tempat Lahir (Kota) terisi: +5%
  if (result.tempatLahir && result.tempatLahir.length >= 2) {
    score += 5
  }

  // f. Alamat Jalan terisi: +5%
  if (result.alamatJalan && result.alamatJalan.length >= 3) {
    score += 5
  }

  // g. RT & RW terisi: +3%
  if (result.rt && result.rw) {
    score += 3
  }

  // h. Kelurahan/Desa atau Kecamatan terisi: +2%
  if (result.kelDesa || result.kecamatan) {
    score += 2
  }

  // Evaluasi 4 Kriteria Utama: NIK 16 digit, Nama, Jenis Kelamin, Tanggal Lahir
  const hasCore4 =
    Boolean(result.isValidNik && result.nik && result.nik.length === 16) &&
    Boolean(result.nama && result.nama.trim().length >= 2) &&
    Boolean(result.jenisKelamin) &&
    Boolean(result.tglLahir && /^\d{2}[-/.]\d{2}[-/.]\d{4}$/.test(result.tglLahir))

  if (hasCore4) {
    // Apabila 4 kriteria utama terpenuhi, confidence dasar minimal 90%
    score = Math.max(score, 90)
    if (result.tempatLahir && result.alamatJalan) {
      score = Math.max(score, 95)
    }
    if (result.rt && result.rw && (result.kelDesa || result.kecamatan)) {
      score = 100
    }
  }

  result.confidence = Math.min(100, Math.max(0, score))

  return result
}

/**
 * 4. Kondisi Berhenti (Stop Scanning & Auto-Fill)
 * Hentikan proses scanning kamera dan segera masukkan data ke dalam form secara otomatis
 * jika seluruh kriteria berikut telah terpenuhi secara bersamaan dalam satu frame:
 * a. NIK terbaca dan tervalidasi tepat 16 digit angka.
 * b. Nama terbaca dan terisi.
 * c. Jenis Kelamin terbaca dan terisi.
 * d. Tanggal Lahir terbaca dan terisi.
 * e. Akurasi Scan (Confidence Level) berada di atas 90%.
 */
export function isAutoCatchReady(data: KtpData): boolean {
  // a. NIK terbaca dan tervalidasi tepat 16 digit angka numerik
  const hasValidNik = Boolean(data.isValidNik && data.nik && data.nik.length === 16 && /^\d{16}$/.test(data.nik))

  // b. Nama terbaca dan terisi
  const hasNama = Boolean(data.nama && data.nama.trim().length >= 2)

  // c. Jenis Kelamin terbaca dan terisi
  const hasJenisKelamin = Boolean(data.jenisKelamin && (data.jenisKelamin === 'LAKI-LAKI' || data.jenisKelamin === 'PEREMPUAN'))

  // d. Tanggal Lahir terbaca dan terisi
  const hasTglLahir = Boolean(data.tglLahir && /^\d{2}[-/.]\d{2}[-/.]\d{4}$/.test(data.tglLahir))

  // e. Akurasi Scan (Confidence Level) berada di atas 90% (>= 90%)
  const hasConfidenceAbove90 = data.confidence >= 90

  // Seluruh kriteria harus terpenuhi bersamaan dalam satu frame
  return hasValidNik && hasNama && hasJenisKelamin && hasTglLahir && hasConfidenceAbove90
}
