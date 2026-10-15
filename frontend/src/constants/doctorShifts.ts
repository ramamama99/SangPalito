// Master Data & Logic for Doctor Shifts (3 Shifts per Polyclinic)
// Formatted strictly according to Smart-Time Algorithm and Puskesmas Operational Rules

export interface ShiftTimeConfig {
  id: "pagi" | "siang" | "sore";
  label: "Pagi" | "Siang" | "Sore";
  doctorStartTime: string;  // Jam mulai dokter bertugas
  doctorEndTime: string;    // Jam selesai dokter bertugas
  patientStartTime: string; // Jam mulai pasien dilayani
  patientEndTime: string;   // Jam selesai pasien dilayani (1 jam buffer sebelum dokter selesai)
  doctorHoursText: string;  // e.g. "07.00 – 12.00"
  patientHoursText: string; // e.g. "07.00 – 11.00"
  effectivePatientMinutes: number; // 240 menit (4 jam)
}

export const SHIFT_TIMES: Record<"pagi" | "siang" | "sore", ShiftTimeConfig> = {
  pagi: {
    id: "pagi",
    label: "Pagi",
    doctorStartTime: "07:00",
    doctorEndTime: "12:00",
    patientStartTime: "07:00",
    patientEndTime: "11:00",
    doctorHoursText: "07.00 – 12.00",
    patientHoursText: "07.00 – 11.00",
    effectivePatientMinutes: 240, // 4 jam layanan pasien
  },
  siang: {
    id: "siang",
    label: "Siang",
    doctorStartTime: "12:00",
    doctorEndTime: "17:00",
    patientStartTime: "12:00",
    patientEndTime: "16:00",
    doctorHoursText: "12.00 – 17.00",
    patientHoursText: "12.00 – 16.00",
    effectivePatientMinutes: 240, // 4 jam layanan pasien
  },
  sore: {
    id: "sore",
    label: "Sore",
    doctorStartTime: "17:00",
    doctorEndTime: "22:00",
    patientStartTime: "17:00",
    patientEndTime: "21:00",
    doctorHoursText: "17.00 – 22.00",
    patientHoursText: "17.00 – 21.00",
    effectivePatientMinutes: 240, // 4 jam layanan pasien
  },
};

export interface DoctorShift {
  id: string; // e.g. "shift-umum-pagi"
  poliId: string;
  poliName: string;
  poliCode: string;
  shiftId: "pagi" | "siang" | "sore";
  doctorName: string;
  doctorSpecialty: string;
  room: string;
  smartAvgServiceMin: number; // Baseline smart time
  fixedQuota: number; // 240 min / smartAvgServiceMin (Fixed cap, no increment if finished early)
  isActive: boolean;
  notes?: string;
}

export const initialDoctorShifts: DoctorShift[] = [
  // ── 1. POLI UMUM (Avg Service: 10 menit -> Fixed Quota = 240/10 = 24 kuota) ──
  {
    id: "shift-umum-pagi",
    poliId: "poli-umum",
    poliName: "Poli Umum",
    poliCode: "A",
    shiftId: "pagi",
    doctorName: "dr. Rama Palito",
    doctorSpecialty: "Dokter Umum",
    room: "Ruang 01",
    smartAvgServiceMin: 10,
    fixedQuota: 24, // 240 / 10
    isActive: true,
    notes: "Shift utama pemeriksaan umum & skrining kesehatan",
  },
  {
    id: "shift-umum-siang",
    poliId: "poli-umum",
    poliName: "Poli Umum",
    poliCode: "A",
    shiftId: "siang",
    doctorName: "dr. Aditya Rahman",
    doctorSpecialty: "Dokter Umum",
    room: "Ruang 01",
    smartAvgServiceMin: 10,
    fixedQuota: 24,
    isActive: true,
    notes: "Pemeriksaan umum siang & penanganan rujukan",
  },
  {
    id: "shift-umum-sore",
    poliId: "poli-umum",
    poliName: "Poli Umum",
    poliCode: "A",
    shiftId: "sore",
    doctorName: "dr. Maya Putri",
    doctorSpecialty: "Dokter Umum",
    room: "Ruang 01",
    smartAvgServiceMin: 10,
    fixedQuota: 24,
    isActive: true,
    notes: "Pelayanan umum sore & konsultasi malam",
  },

  // ── 2. POLI GIGI (Avg Service: 15 menit -> Fixed Quota = 240/15 = 16 kuota) ──
  {
    id: "shift-gigi-pagi",
    poliId: "poli-gigi",
    poliName: "Poli Gigi",
    poliCode: "B",
    shiftId: "pagi",
    doctorName: "drg. Rani Prameswari",
    doctorSpecialty: "Dokter Gigi",
    room: "Ruang 02",
    smartAvgServiceMin: 15,
    fixedQuota: 16, // 240 / 15
    isActive: true,
    notes: "Pencabutan, penambalan & scaling gigi",
  },
  {
    id: "shift-gigi-siang",
    poliId: "poli-gigi",
    poliName: "Poli Gigi",
    poliCode: "B",
    shiftId: "siang",
    doctorName: "drg. Fajar Hidayat",
    doctorSpecialty: "Dokter Gigi",
    room: "Ruang 02",
    smartAvgServiceMin: 15,
    fixedQuota: 16,
    isActive: true,
    notes: "Perawatan saluran akar & konsultasi gigi",
  },
  {
    id: "shift-gigi-sore",
    poliId: "poli-gigi",
    poliName: "Poli Gigi",
    poliCode: "B",
    shiftId: "sore",
    doctorName: "drg. Rani Prameswari",
    doctorSpecialty: "Dokter Gigi",
    room: "Ruang 02",
    smartAvgServiceMin: 15,
    fixedQuota: 16,
    isActive: true,
    notes: "Penanganan nyeri gigi akut & darurat",
  },

  // ── 3. POLI KIA / ANAK (Avg Service: 12 menit -> Fixed Quota = 240/12 = 20 kuota) ──
  {
    id: "shift-kia-pagi",
    poliId: "poli-kia",
    poliName: "Poli KIA",
    poliCode: "C",
    shiftId: "pagi",
    doctorName: "dr. Siti Rahmawati, Sp.A",
    doctorSpecialty: "Dokter Spesialis Anak",
    room: "Ruang 03",
    smartAvgServiceMin: 12,
    fixedQuota: 20, // 240 / 12
    isActive: true,
    notes: "Imunisasi, tumbuh kembang & pemeriksaan balita",
  },
  {
    id: "shift-kia-siang",
    poliId: "poli-kia",
    poliName: "Poli KIA",
    poliCode: "C",
    shiftId: "siang",
    doctorName: "dr. Budi Setiawan, Sp.A",
    doctorSpecialty: "Dokter Spesialis Anak",
    room: "Ruang 03",
    smartAvgServiceMin: 12,
    fixedQuota: 20,
    isActive: true,
    notes: "Konsultasi laktasi, gizi anak & ibu hamil",
  },
  {
    id: "shift-kia-sore",
    poliId: "poli-kia",
    poliName: "Poli KIA",
    poliCode: "C",
    shiftId: "sore",
    doctorName: "dr. Siti Rahmawati, Sp.A",
    doctorSpecialty: "Dokter Spesialis Anak",
    room: "Ruang 03",
    smartAvgServiceMin: 12,
    fixedQuota: 20,
    isActive: true,
    notes: "Pemeriksaan anak demam & batuk malam",
  },

  // ── 4. POLI LANSIA / GERIATRI (Avg Service: 15 menit -> Fixed Quota = 240/15 = 16 kuota) ──
  {
    id: "shift-lansia-pagi",
    poliId: "poli-lansia",
    poliName: "Poli Lansia",
    poliCode: "D",
    shiftId: "pagi",
    doctorName: "dr. Hendra Prasetyo",
    doctorSpecialty: "Spesialis Penyakit Dalam / Geriatri",
    room: "Ruang 04",
    smartAvgServiceMin: 15,
    fixedQuota: 16, // 240 / 15
    isActive: true,
    notes: "Hipertensi, diabetes & skrining penyakit kronis",
  },
  {
    id: "shift-lansia-siang",
    poliId: "poli-lansia",
    poliName: "Poli Lansia",
    poliCode: "D",
    shiftId: "siang",
    doctorName: "dr. Maya Putri",
    doctorSpecialty: "Dokter Umum Geriatri",
    room: "Ruang 04",
    smartAvgServiceMin: 15,
    fixedQuota: 16,
    isActive: true,
    notes: "Konsultasi obat rutin & pemeriksaan sendi lansia",
  },
  {
    id: "shift-lansia-sore",
    poliId: "poli-lansia",
    poliName: "Poli Lansia",
    poliCode: "D",
    shiftId: "sore",
    doctorName: "dr. Hendra Prasetyo",
    doctorSpecialty: "Spesialis Penyakit Dalam",
    room: "Ruang 04",
    smartAvgServiceMin: 15,
    fixedQuota: 16,
    isActive: true,
    notes: "Kontrol tensi & evaluasi resep obat kronis",
  },
];

const STORAGE_KEY = "antrisehat-doctor-shifts-v1";

export function loadDoctorShifts(): DoctorShift[] {
  if (typeof window === "undefined") return initialDoctorShifts;
  try {
    const raw = localStorage.getItem(STORAGE_KEY);
    if (!raw) {
      localStorage.setItem(STORAGE_KEY, JSON.stringify(initialDoctorShifts));
      return initialDoctorShifts;
    }
    const parsed = JSON.parse(raw);
    if (Array.isArray(parsed) && parsed.length > 0) {
      return parsed;
    }
    return initialDoctorShifts;
  } catch {
    return initialDoctorShifts;
  }
}

export function saveDoctorShifts(shifts: DoctorShift[]): void {
  if (typeof window === "undefined") return;
  try {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(shifts));
  } catch (err) {
    console.error("Failed to save doctor shifts to localStorage", err);
  }
}

export function getCurrentShiftId(): "pagi" | "siang" | "sore" {
  const now = new Date();
  // Gunakan timezone WIB (UTC+7)
  const jakartaHours = (now.getUTCHours() + 7) % 24;

  if (jakartaHours >= 7 && jakartaHours < 12) {
    return "pagi";
  } else if (jakartaHours >= 12 && jakartaHours < 17) {
    return "siang";
  } else {
    return "sore";
  }
}

/**
 * Menghitung estimasi jam kedatangan pasien berdasarkan jam mulai shift dan nomor urut di shift tersebut.
 * Urutan 1 di Shift Pagi (Poli Umum) = 07.00 WIB
 * Urutan 4 di Shift Pagi (Poli Umum) = 07.30 WIB (07.00 + 3 * 10 mnt)
 */
export function calculateEstimatedArrival(
  shiftId: "pagi" | "siang" | "sore" | string,
  smartAvgMinutes: number,
  orderInShift: number = 1
): string {
  const shift = SHIFT_TIMES[(shiftId as "pagi" | "siang" | "sore")] || SHIFT_TIMES.pagi;
  const [startH, startM] = shift.patientStartTime.split(":").map(Number);
  
  // Total menit dari awal shift pasien
  const additionalMinutes = Math.max(0, orderInShift - 1) * smartAvgMinutes;
  const totalMinutes = startH * 60 + startM + additionalMinutes;
  
  const arrivalH = Math.floor(totalMinutes / 60);
  const arrivalM = totalMinutes % 60;
  
  const formattedH = String(arrivalH).padStart(2, "0");
  const formattedM = String(arrivalM).padStart(2, "0");
  
  return `${formattedH}.${formattedM} WIB`;
}

/**
 * Menghitung jam anjuran datang (15 menit sebelum estimasi kedatangan).
 */
export function calculateRecommendedArrivalTime(estimatedArrivalStr: string): string {
  const match = estimatedArrivalStr.match(/(\d{1,2})[.:](\d{2})/);
  if (!match) return "15 menit sebelum estimasi";
  
  let h = parseInt(match[1], 10);
  let m = parseInt(match[2], 10);
  
  let totalMin = h * 60 + m - 15;
  if (totalMin < 0) totalMin = 0;
  
  const recH = String(Math.floor(totalMin / 60)).padStart(2, "0");
  const recM = String(totalMin % 60).padStart(2, "0");
  
  return `${recH}.${recM} WIB`;
}
