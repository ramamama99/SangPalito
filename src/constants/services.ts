import { Stethoscope, Smile, Baby, Heart } from 'lucide-vue-next'
import type { Component } from 'vue'

export type ServiceItem = {
  name: string
  desc: string
  icon: Component
  quota: number
}

export const services: ServiceItem[] = [
  {
    name: 'Poli Umum',
    desc: 'Pemeriksaan kesehatan & konsultasi umum',
    icon: Stethoscope,
    quota: 24,
  },
  {
    name: 'Poli Gigi',
    desc: 'Perawatan gigi dan kesehatan mulut',
    icon: Smile,
    quota: 12,
  },
  {
    name: 'Poli KIA',
    desc: 'Kesehatan ibu, anak, dan tumbuh kembang',
    icon: Baby,
    quota: 18,
  },
  {
    name: 'Poli Lansia',
    desc: 'Pelayanan khusus untuk usia lanjut',
    icon: Heart,
    quota: 15,
  },
]
