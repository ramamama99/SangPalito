import type { Component } from 'vue'

export type Patient = {
  id: string
  number: string
  name: string
  nik: string
  phone: string
  birth: string
  gender: string
  insurance: string
  bpjs: string
  date: string
  poli: string
  complaint: string
  allergy: string
  symptoms: string[]
  status: string
  source: string
  soap?: string
}

export type Service = {
  name: string
  desc: string
  icon: Component
  quota: number
}

export type SoapFormData = {
  bp_sys?: string
  bp_dia?: string
  hr?: string
  temp?: string
  spo2?: string
  weight?: string
  height?: string
  doctor?: string
  subjective?: string
  objective?: string
  diagnosis?: string
  plan?: string
  [key: string]: unknown
}
