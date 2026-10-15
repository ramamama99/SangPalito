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
  soap?: string;
  qr_token?: string;
  appointment_token?: string;
};

export type Service = {
  name: string;
  desc: string;
  iconName: string;
  quota: number;
};

export type SoapData = {
  subjective: string;
  objective: string;
  diagnosis: string;
  plan: string;
  temperature: string;
  pressure: string;
  pulse: string;
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
  complaint: string;
  allergy: string;
  symptoms: string[];
};
