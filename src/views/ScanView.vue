<script setup lang="ts">
import { useRouter } from 'vue-router'
import PatientScanner from '../components/PatientScanner.vue'
import { useQueueStore, today } from '../stores/queue'

const router = useRouter()
const queueStore = useQueueStore()

const handleScan = (value: string) => {
  let id = value.trim()
  try {
    const url = new URL(id)
    if (url.origin !== window.location.origin || url.pathname !== window.location.pathname) {
      return 'QR bukan e-tiket AntriSehat dari aplikasi ini.'
    }
    id = url.searchParams.get('checkin') || ''
  } catch {
    if (!/^[a-zA-Z0-9-]+$/.test(id)) return 'Kode e-tiket tidak valid.'
  }

  const patient = queueStore.patients.find((p) => p.id === id)
  if (!patient) return 'Pasien tidak ditemukan pada data perangkat ini.'
  if (patient.date !== today) return 'E-tiket ini bukan untuk kunjungan hari ini.'
  if (!['Terdaftar', 'Menunggu verifikasi'].includes(patient.status))
    return 'Pasien sudah diverifikasi atau kunjungannya tidak lagi aktif.'

  queueStore.verifyPatient = patient
  router.push('/staff')
  return null
}
</script>

<template>
  <PatientScanner :on-scan="handleScan" @back="router.push('/staff')" />
</template>
