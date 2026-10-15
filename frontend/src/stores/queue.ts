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
      return stored ? JSON.parse(stored) : initialPatients;
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
        qr.value = await generateQrDataUrl(newVal);
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
    generateQrDataUrl(ticket.value)
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
    patients.value.find((p) => p.id === ticket.value),
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

  function registerPatient() {
    const id = crypto.randomUUID();
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
    const number = `${prefix}-${String(maxNum + 1).padStart(3, "0")}`;

    const newPatient: Patient = {
      ...form.value,
      id,
      number,
      source: walkin.value ? "Langsung" : "Online",
      status: walkin.value ? "Hadir" : "Terdaftar",
    };

    patients.value.push(newPatient);
    ticket.value = id;

    if (walkin.value) {
      setNotice("Pasien walk-in berhasil didaftarkan.");
    }
    return { id, isWalkin: walkin.value };
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
    downloadTicket,
    fetchQueuesFromApi,
  };
});
