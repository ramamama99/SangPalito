import { defineStore } from "pinia";
import { ref, computed, watch } from "vue";
import type { Patient, RegistrationForm, SoapData } from "../types";

export interface CallAnnouncement {
  id: string;
  number: string;
  name: string;
  poli: string;
  transcription: string;
  timestamp: number;
}

import {
  initialPatients,
  blankRegistrationForm,
  todayDate,
} from "../constants/initialData";
import {
  loadDoctorShifts,
  saveDoctorShifts,
  SHIFT_TIMES,
  calculateEstimatedArrival,
  calculateRecommendedArrivalTime,
  type DoctorShift,
} from "../constants/doctorShifts";
import { generateQrDataUrl } from "../utils/qr";
import { speakQueueCall } from "../utils/speech";
import { downloadTicketImage } from "../utils/ticketCanvas";
import { api } from "../services/api";

export const useQueueStore = defineStore("queue", () => {
  function getStoredPatients(): Patient[] {
    try {
      const stored = localStorage.getItem("antrisehat-patients");
      if (stored) {
        const parsed = JSON.parse(stored);
        if (Array.isArray(parsed)) {
          // Bersihkan seluruh data pasien dummy (demo-*)
          const filtered = parsed.filter(
            (p: Patient) => p && p.id && !p.id.startsWith("demo-"),
          );
          localStorage.setItem("antrisehat-patients", JSON.stringify(filtered));
          return filtered;
        }
      }
      return initialPatients;
    } catch {
      return initialPatients;
    }
  }

  function getStoredTicket(): string {
    try {
      const t = localStorage.getItem("antrisehat-ticket") || "";
      if (t && t.startsWith("demo-")) {
        localStorage.removeItem("antrisehat-ticket");
        return "";
      }
      return t;
    } catch {
      return "";
    }
  }

  // State
  const patients = ref<Patient[]>(getStoredPatients());
  const ticket = ref<string>(getStoredTicket());

  const qr = ref<string>("");
  const notice = ref<string>("");
  let noticeTimer: ReturnType<typeof setTimeout> | null = null;
  const lastCallAnnouncement = ref<CallAnnouncement | null>(null);

  const role = ref<"loket" | "medis" | "database_pasien" | "database_obat" | "shift_dokter" | "user_management" | "audit_log" | "analytics">("loket");
  const selected = ref<string>("");
  const verify = ref<Patient | null>(null);
  const search = ref<string>("");
  const staffPoli = ref<string>("Poli Umum");
  const walkin = ref<boolean>(false);

  const lookup = ref<string>("");
  const lookupSearch = ref<string>("");
  const lookupOpen = ref<boolean>(false);

  const form = ref<RegistrationForm>({ ...blankRegistrationForm });
  const step = ref<number>(1);
  const human = ref<boolean>(false);

  const soap = ref<SoapData>({
    subjective: "",
    objective: "",
    diagnosis: "",
    plan: "",
    temperature: "",
    pressure: "",
    pulse: "",
  });

  // Doctor Shifts Management State (3 shifts per polyclinic)
  const doctorShifts = ref<DoctorShift[]>(loadDoctorShifts());

  // Watchers for persistence
  watch(
    doctorShifts,
    (newVal) => {
      saveDoctorShifts(newVal);
    },
    { deep: true },
  );
  watch(
    patients,
    (newVal) => {
      localStorage.setItem("antrisehat-patients", JSON.stringify(newVal));
    },
    { deep: true },
  );

  watch(ticket, async (newVal) => {
    if (newVal) {
      localStorage.setItem("antrisehat-ticket", newVal);
      try {
        const found = patients.value.find(
          (p) => p.id === newVal || p.qr_token === newVal,
        );
        const tokenToEncode = found?.qr_token || newVal;
        qr.value = await generateQrDataUrl(tokenToEncode);
      } catch (err) {
        console.error("Failed to generate QR code", err);
      }
    } else {
      localStorage.removeItem("antrisehat-ticket");
      qr.value = "";
    }
  });

  // Initial QR code generation if ticket exists
  if (ticket.value) {
    const found = patients.value.find(
      (p) => p.id === ticket.value || p.qr_token === ticket.value,
    );
    generateQrDataUrl(found?.qr_token || ticket.value)
      .then((res) => {
        qr.value = res;
      })
      .catch(console.error);
  }

  // Cross-tab synchronization
  if (typeof window !== "undefined") {
    window.addEventListener("storage", (e: StorageEvent) => {
      if (e.key === "antrisehat-patients" && e.newValue) {
        try {
          patients.value = JSON.parse(e.newValue);
        } catch {
          // ignore parsing error
        }
      }
      if (e.key === "antrisehat-ticket" && e.newValue) {
        ticket.value = e.newValue;
      }
      if (e.key === "antrisehat-last-call" && e.newValue) {
        try {
          lastCallAnnouncement.value = JSON.parse(e.newValue);
        } catch {
          // ignore parsing error
        }
      }
    });

    // URL parameter checkin handler
    const params = new URLSearchParams(window.location.search);
    const checkinId = params.get("checkin");
    if (checkinId) {
      ticket.value = checkinId;
    }
  }

  // Getters / Computed
  const currentTicket = computed(() =>
    patients.value.find(
      (p) => p.id === ticket.value || p.qr_token === ticket.value,
    ),
  );
  const active = computed(() =>
    patients.value.find((p) => p.status === "Dipanggil"),
  );
  const chosen = computed(() =>
    patients.value.find((p) => p.id === selected.value),
  );

  // Actions
  function setNotice(message: string) {
    notice.value = message;
    if (noticeTimer) clearTimeout(noticeTimer);
    noticeTimer = setTimeout(() => {
      notice.value = "";
    }, 5000);
  }

  function clearNotice() {
    notice.value = "";
    if (noticeTimer) clearTimeout(noticeTimer);
  }

  const statusMapToBackend: Record<string, string> = {
    Hadir: "checked_in",
    Dipanggil: "called",
    "Sedang dilayani": "in_service",
    Selesai: "done",
    Batal: "cancelled",
    "Tidak Hadir": "no_show",
    "Menunggu verifikasi": "waiting",
  };

  const statusMapFromBackend: Record<string, string> = {
    waiting: "Menunggu verifikasi",
    checked_in: "Hadir",
    called: "Dipanggil",
    in_service: "Sedang dilayani",
    done: "Selesai",
    cancelled: "Batal",
    no_show: "Tidak Hadir",
  };

  async function fetchQueuesFromApi() {
    try {
      const data = await api.getQueuesToday();
      if (data && data.length > 0) {
        const mappedPatients: Patient[] = data.map((q) => ({
          id: q.id,
          number: q.queue_code,
          name: q.patient?.full_name || "Pasien",
          nik: q.patient?.medical_record_no || "",
          phone: "",
          birth: "",
          gender: "Laki-laki",
          insurance: "Umum",
          bpjs: "",
          date: q.queue_date || todayDate,
          poli: q.polyclinic?.name || "Poli Umum",
          complaint: "",
          allergy: "",
          symptoms: [],
          status: statusMapFromBackend[q.status] || q.status,
          source: "Online",
        }));
        patients.value = mappedPatients;
      }
    } catch {
      // Keep existing demo data if backend offline
    }
  }

  // Auto fetch from API on load
  fetchQueuesFromApi();

  function changeStatus(id: string, status: string) {
    patients.value = patients.value.map((p) =>
      p.id === id ? { ...p, status } : p,
    );
    const backendStatus = statusMapToBackend[status] || status;
    api.updateQueueStatus(id, backendStatus).catch(() => {});
  }

  function isSamePoli(p1: Patient, p2: Patient): boolean {
    if (p1.id === p2.id) return true;

    // 1. Cek dari prefix kode antrean (misal A-001 vs A-002)
    const prefix1 = p1.number ? p1.number.split("-")[0]?.trim().toUpperCase() : "";
    const prefix2 = p2.number ? p2.number.split("-")[0]?.trim().toUpperCase() : "";
    if (prefix1 && prefix2 && prefix1 === prefix2) {
      return true;
    }

    // 2. Cek dari nama poliklinik
    if (p1.poli && p2.poli) {
      const s1 = p1.poli.toLowerCase().trim();
      const s2 = p2.poli.toLowerCase().trim();
      if (s1 === s2) return true;

      const getPoliCat = (str: string): string => {
        if (str.includes("umum")) return "umum";
        if (str.includes("gigi")) return "gigi";
        if (str.includes("kia") || str.includes("anak")) return "kia";
        if (str.includes("lansia") || str.includes("geriatri")) return "lansia";
        return str;
      };

      const c1 = getPoliCat(s1);
      const c2 = getPoliCat(s2);
      if (c1 && c2 && c1 === c2) {
        return true;
      }
    }

    return false;
  }

  function call(p: Patient) {
    patients.value = patients.value.map((item) => {
      // 1. Pasien yang baru saja ditekan tombol 'Panggil': status menjadi 'Dipanggil'
      if (item.id === p.id) {
        return { ...item, status: "Dipanggil" };
      }
      // 2. Hanya pasien sebelumnya DI POLI YANG SAMA yang status 'Dipanggil'-nya dikembalikan ke 'Hadir'
      if (isSamePoli(item, p) && item.status === "Dipanggil") {
        return { ...item, status: "Hadir" };
      }
      // 3. Pasien di POLI LAIN (misal Poli A saat Poli C dipanggil) TETAP DIPERTAHANKAN dan tidak terhapus!
      return item;
    });

    const poliName = p.poli || "Poliklinik";
    const transcription = `Nomor antrean ${p.number}, silakan menuju ${poliName}`;
    setNotice(transcription);

    const callPayload: CallAnnouncement = {
      id: p.id,
      number: p.number,
      name: p.name,
      poli: poliName,
      transcription,
      timestamp: Date.now(),
    };
    lastCallAnnouncement.value = callPayload;
    try {
      localStorage.setItem("antrisehat-last-call", JSON.stringify(callPayload));
    } catch {
      // ignore storage error
    }

    speakQueueCall(p.number, poliName);
    api.updateQueueStatus(p.id, "called").catch(() => {});
  }


  function updateForm<K extends keyof RegistrationForm>(
    key: K,
    value: RegistrationForm[K],
  ) {
    form.value[key] = value;
  }

  function startRegistration(poli = "Poli Umum", isWalkin = false) {
    form.value = { ...blankRegistrationForm, poli };
    step.value = 1;
    human.value = false;
    walkin.value = isWalkin;
  }

  const isRegistering = ref(false);

  const polyclinicIdMap: Record<string, string> = {
    "Poli Umum": "11111111-1111-1111-1111-111111111111",
    "Poli Gigi": "22222222-2222-2222-2222-222222222222",
    "Poli KIA": "33333333-3333-3333-3333-333333333333",
    "Poli KIA / Anak": "33333333-3333-3333-3333-333333333333",
    "Poli Lansia": "44444444-4444-4444-4444-444444444444",
    "Poli Lansia / Geriatri": "44444444-4444-4444-4444-444444444444",
  };

  function getPoliCode(poliName: string): string {
    if (!poliName) return "A";
    if (poliName.includes("Gigi")) return "B";
    if (poliName.includes("KIA") || poliName.includes("Anak")) return "C";
    if (poliName.includes("Lansia") || poliName.includes("Geriatri")) return "D";
    return "A";
  }

  function getShiftQuotaInfo(
    poliName: string,
    shiftId: "pagi" | "siang" | "sore" | string,
    targetDate?: string,
  ) {
    const checkDate = targetDate || form.value.date || todayDate;
    const code = getPoliCode(poliName);

    // Find configured doctor shift
    const shiftConfig = doctorShifts.value.find(
      (s) => (s.poliCode === code || s.poliName === poliName) && s.shiftId === shiftId,
    );

    const maxQuota = shiftConfig ? shiftConfig.fixedQuota : 24;
    const smartAvg = shiftConfig ? shiftConfig.smartAvgServiceMin : 10;
    const doctorName = shiftConfig ? shiftConfig.doctorName : "Dokter Bertugas";
    const doctorSpecialty = shiftConfig ? shiftConfig.doctorSpecialty : "";
    const room = shiftConfig ? shiftConfig.room : "Ruang 01";

    // Count currently booked patients in this shift for that polyclinic and date
    const bookedCount = patients.value.filter((p) => {
      const pCode = getPoliCode(p.poli);
      const matchesPoli = pCode === code;
      const pShift = p.shift || "pagi";
      const pDate = p.date || todayDate;
      const isNotCancelled = p.status !== "Batal" && p.status !== "cancelled";
      return matchesPoli && pShift === shiftId && pDate === checkDate && isNotCancelled;
    }).length;

    const remainingQuota = Math.max(0, maxQuota - bookedCount);
    const isFull = remainingQuota <= 0;
    const nextQueueOrder = bookedCount + 1;
    const estimatedArrival = calculateEstimatedArrival(shiftId, smartAvg, nextQueueOrder);
    const recommendedArrival = calculateRecommendedArrivalTime(estimatedArrival);
    const shiftTimeInfo = SHIFT_TIMES[shiftId as "pagi" | "siang" | "sore"] || SHIFT_TIMES.pagi;

    return {
      shiftConfig,
      shiftTimeInfo,
      maxQuota,
      bookedCount,
      remainingQuota,
      isFull,
      doctorName,
      doctorSpecialty,
      room,
      smartAvg,
      nextQueueOrder,
      estimatedArrival,
      recommendedArrival,
    };
  }

  function updateDoctorShift(updated: DoctorShift) {
    doctorShifts.value = doctorShifts.value.map((s) => (s.id === updated.id ? { ...updated } : s));
    saveDoctorShifts(doctorShifts.value);
    setNotice(`Jadwal shift ${updated.doctorName} (${updated.poliName}) berhasil diperbarui.`);
  }

  async function registerPatient() {
    isRegistering.value = true;
    let id: string = crypto.randomUUID();
    let qrToken = "";
    let queueCode = "";
    let appointmentToken = "";

    const polyclinicId =
      polyclinicIdMap[form.value.poli] || "11111111-1111-1111-1111-111111111111";
    const guarantor =
      form.value.insurance === "BPJS"
        ? "bpjs"
        : form.value.insurance === "Asuransi Swasta"
          ? "asuransi"
          : "umum";
    const gender = form.value.gender === "Perempuan" ? "female" : "male";
    const visitDate = form.value.date || todayDate;
    const birthDate = form.value.birth || "1995-01-01";

    try {
      const res = await api.createAppointment({
        nik: form.value.nik || "3201000000000001",
        full_name: form.value.name || "Pasien Baru",
        birth_date: birthDate,
        gender,
        phone: form.value.phone || "081234567890",
        guarantor,
        bpjs_no: form.value.bpjs || undefined,
        polyclinic_id: polyclinicId,
        visit_date: visitDate,
        chief_complaint: form.value.complaint || "Pemeriksaan kesehatan umum",
        symptoms:
          form.value.symptoms && form.value.symptoms.length > 0
            ? form.value.symptoms
            : ["Pemeriksaan Umum"],
        duration_days: 1,
      });

      if (res && res.success && res.data) {
        qrToken = res.data.qr_token;
        queueCode = res.data.queue?.queue_code || "";
        appointmentToken = res.data.appointment_token || "";
        id = qrToken; // Set ID directly to the database QR token
      }
    } catch (err) {
      console.error("API createAppointment failed:", err);
    } finally {
      isRegistering.value = false;
    }

    if (!queueCode) {
      const prefix =
        form.value.poli === "Poli Gigi"
          ? "B"
          : form.value.poli === "Poli KIA"
            ? "C"
            : form.value.poli === "Poli Lansia"
              ? "D"
              : "A";

      const existingPoliNumbers = patients.value
        .filter((p) => p.number.startsWith(prefix))
        .map((p) => Number(p.number.split("-")[1]));

      const maxNum = Math.max(11, ...existingPoliNumbers);
      queueCode = `${prefix}-${String(maxNum + 1).padStart(3, "0")}`;
    }

    if (!qrToken) {
      qrToken = id;
    }

    const selectedShift = form.value.shift || "pagi";
    const quotaInfo = getShiftQuotaInfo(form.value.poli, selectedShift, form.value.date);
    const estimatedArrival = quotaInfo.estimatedArrival;
    const doctorName = quotaInfo.doctorName;
    const shiftTimeInfo = SHIFT_TIMES[selectedShift as "pagi" | "siang" | "sore"] || SHIFT_TIMES.pagi;
    const shiftLabel = `Shift ${shiftTimeInfo.label}`;

    form.value.estimatedArrival = estimatedArrival;
    form.value.shiftLabel = shiftLabel;

    const newPatient: Patient = {
      ...form.value,
      id,
      number: queueCode,
      qr_token: qrToken,
      appointment_token: appointmentToken,
      source: walkin.value ? "Langsung" : "Online",
      status: walkin.value ? "Hadir" : "Terdaftar",
      shift: selectedShift,
      shiftLabel,
      doctorName,
      estimatedArrival,
    };

    patients.value.push(newPatient);
    ticket.value = id;

    // Immediately generate real QR data URL for the ticket
    try {
      qr.value = await generateQrDataUrl(qrToken);
    } catch (e) {
      console.error("Failed to generate QR data url:", e);
    }

    if (walkin.value) {
      setNotice("Pasien walk-in berhasil didaftarkan.");
    }
    return { id, isWalkin: walkin.value, qr_token: qrToken, queue_code: queueCode };
  }

  function downloadTicket() {
    if (currentTicket.value && qr.value) {
      downloadTicketImage(currentTicket.value, qr.value);
    }
  }

  return {
    patients,
    ticket,
    qr,
    notice,
    role,
    selected,
    verify,
    search,
    staffPoli,
    walkin,
    lookup,
    lookupSearch,
    lookupOpen,
    form,
    step,
    human,
    soap,
    doctorShifts,
    currentTicket,
    active,
    chosen,
    lastCallAnnouncement,
    isSamePoli,
    setNotice,
    clearNotice,
    changeStatus,
    call,
    updateForm,
    startRegistration,
    registerPatient,
    isRegistering,
    downloadTicket,
    fetchQueuesFromApi,
    getShiftQuotaInfo,
    updateDoctorShift,
    getPoliCode,
  };
});

