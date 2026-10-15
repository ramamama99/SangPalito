import type { Patient, RegistrationForm } from "../types";

export const todayDate = new Date().toISOString().slice(0, 10);

export const initialPatients: Patient[] = [
  "Siti Aminah",
  "Budi Santoso",
  "Dewi Lestari",
  "Ahmad Fauzi",
].map((name, i) => ({
  id: `demo-${i}`,
  number: `A-${String(i + 8).padStart(3, "0")}`,
  name,
  nik: `327301120190000${i}`,
  phone: "081234567890",
  birth: "1990-01-12",
  gender: i % 2 ? "Laki-laki" : "Perempuan",
  insurance: "BPJS",
  bpjs: "0001234567890",
  date: todayDate,
  poli: "Poli Umum",
  complaint: "Demam dan sakit kepala sejak kemarin",
  allergy: "Tidak ada",
  symptoms: ["Demam", "Pusing"],
  status: i === 0 ? "Dipanggil" : i === 1 ? "Hadir" : "Menunggu verifikasi",
  source: i === 3 ? "Langsung" : "Online",
}));

export const blankRegistrationForm: RegistrationForm = {
  name: "",
  nik: "",
  phone: "",
  birth: "",
  gender: "",
  insurance: "Umum",
  bpjs: "",
  date: todayDate,
  poli: "Poli Umum",
  complaint: "",
  allergy: "",
  symptoms: [],
};
