import { defineStore } from "pinia";
import { ref, computed, watch } from "vue";
import type { Patient, RegistrationForm, SoapData } from "../types";
import {
  initialPatients,
  blankRegistrationForm,
} from "../constants/initialData";
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
        const hasOtherPoli = parsed.some((p: Patient) => p.poli && p.poli !== "Poli Umum");
        if (hasOtherPoli) {
          return parsed;
        }
      }
      return initialPatients;
    } catch {
      return initialPatients;
    }
  }

  function getStoredTicket(): string {
    try {
      return localStorage.getItem("antrisehat-ticket") || "";
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

  const role = ref<"loket" | "medis">("loket");
  const selected = ref<string>("demo-1");
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

  // Watchers for persistence
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
          date: q.queue_date || new Date().toISOString().slice(0, 10),
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

  function call(p: Patient) {
    patients.value = patients.value.map((item) =>
      item.id === p.id
        ? { ...item, status: "Dipanggil" }
        : item.status === "Dipanggil"
          ? { ...item, status: "Hadir" }
          : item,
    );
    setNotice(`${p.number} dipanggil ke ${p.poli}`);
    speakQueueCall(p.number, p.poli);
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
    const visitDate = form.value.date || new Date().toISOString().slice(0, 10);
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

    const newPatient: Patient = {
      ...form.value,
      id,
      number: queueCode,
      qr_token: qrToken,
      appointment_token: appointmentToken,
      source: walkin.value ? "Langsung" : "Online",
      status: walkin.value ? "Hadir" : "Terdaftar",
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
    currentTicket,
    active,
    chosen,
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
  };
});
