export type Medication = {
  id: string;
  code: string;
  name: string;
  strength?: string;
  generic_name: string;
  category: string;
  dosage_form: string;
  route: string; // oral_minum, oral_kunyah, topikal_oles, tetes_mata, tetes_telinga, kumur, inhalasi
  default_frequency: string;
  default_timing: string;
  usage_instruction: string;
  default_instruction?: string;
  indication: string;
  stock: number;
  unit: string;
  is_prescription_required: boolean;
};

export type PrescriptionItem = {
  id?: string;
  medication_id?: string;
  medication_name: string;
  dosage_form: string;
  quantity: number;
  unit: string;
  route: string;
  frequency: string;
  timing: string;
  instruction: string;
  notes?: string;
};

export type Prescription = {
  id: string;
  prescription_no: string;
  appointment_id?: string;
  patient_id?: string;
  patient_name?: string;
  doctor_name?: string;
  polyclinic?: string;
  status: "pending_dispense" | "dispensed" | "cancelled" | string;
  status_label?: string;
  notes?: string;
  created_at?: string;
  dispensed_at?: string;
  items: PrescriptionItem[];
};

export type Patient = {
  id: string;
  number: string;
  name: string;
  nik: string;
  phone: string;
  birth: string;
  birthPlace?: string;
  gender: string;
  bloodType?: string;
  address?: string;
  rt?: string;
  rw?: string;
  kelurahan?: string;
  kecamatan?: string;
  insurance: string;
  bpjs: string;
  insuranceName?: string;
  date: string;
  poli: string;
  complaint: string;
  allergy: string;
  symptoms: string[];
  status: string;
  source: string;
  shift?: "pagi" | "siang" | "sore" | string;
  shiftLabel?: string;
  estimatedArrival?: string;
  doctorName?: string;
  soap?: string;
  medication_plan?: string;
  prescription_no?: string;
  prescription?: any;
  qr_token?: string;
  appointment_token?: string;
  createdAt?: string;
};

export type Service = {
  name: string;
  desc: string;
  iconName: string;
  quota: number;
  avgServiceMin?: number;
};

export type SoapData = {
  subjective: string;
  objective: string;
  diagnosis: string;
  plan: string;
  medication_plan?: string;
  temperature: string;
  pressure: string;
  pulse: string;
  prescriptions?: PrescriptionItem[];
  issuePrescription?: boolean;
};

export type RegistrationForm = {
  name: string;
  nik: string;
  phone: string;
  birth: string;
  birthPlace: string;
  gender: string;
  bloodType: string;
  address: string;
  rt: string;
  rw: string;
  kelurahan: string;
  kecamatan: string;
  insurance: string;
  bpjs: string;
  insuranceName?: string;
  date: string;
  poli: string;
  shift: "pagi" | "siang" | "sore";
  shiftLabel?: string;
  estimatedArrival?: string;
  complaint: string;
  allergy: string;
  symptoms: string[];
};
