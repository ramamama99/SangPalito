import type { Service } from "../types";

export const services: Service[] = [
  {
    name: "Poli Umum",
    desc: "Pemeriksaan kesehatan & konsultasi umum",
    iconName: "Stethoscope",
    quota: 24,
  },
  {
    name: "Poli Gigi",
    desc: "Perawatan gigi dan kesehatan mulut",
    iconName: "Smile",
    quota: 12,
  },
  {
    name: "Poli KIA",
    desc: "Kesehatan ibu, anak, dan tumbuh kembang",
    iconName: "Baby",
    quota: 18,
  },
  {
    name: "Poli Lansia",
    desc: "Pelayanan khusus untuk usia lanjut",
    iconName: "Heart",
    quota: 15,
  },
];
