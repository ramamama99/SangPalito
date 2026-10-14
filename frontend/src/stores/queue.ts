import { defineStore } from 'pinia'
import { ref, computed, watch } from 'vue'
import QRCode from 'qrcode'
import type { Patient } from '../types'

export const today = new Date().toISOString().slice(0, 10)

const initialPatients: Patient[] = [
  'Siti Aminah',
  'Budi Santoso',
  'Dewi Lestari',
  'Ahmad Fauzi',
].map((name, i) => ({
  id: `demo-${i}`,
  number: `A-${String(i + 8).padStart(3, '0')}`,
  name,
  nik: `327301120190000${i}`,
  phone: '081234567890',
  birth: '1990-01-12',
  gender: i % 2 ? 'Laki-laki' : 'Perempuan',
  insurance: 'BPJS',
  bpjs: '0001234567890',
  date: today,
  poli: 'Poli Umum',
  complaint: 'Demam dan sakit kepala sejak kemarin',
  allergy: 'Tidak ada',
  symptoms: ['Demam', 'Pusing'],
  status: i === 0 ? 'Dipanggil' : i === 1 ? 'Hadir' : 'Menunggu verifikasi',
  source: i === 3 ? 'Langsung' : 'Online',
}))

const getInitialPatients = (): Patient[] => {
  try {
    return JSON.parse(localStorage.getItem('antrisehat-patients') || 'null') || initialPatients
  } catch {
    return initialPatients
  }
}

export const useQueueStore = defineStore('queue', () => {
  const patients = ref<Patient[]>(getInitialPatients())

  const ticket = ref<string>(localStorage.getItem('antrisehat-ticket') || '')
  const qr = ref<string>('')
  const notice = ref<string>('')
  let noticeTimer: ReturnType<typeof setTimeout> | undefined

  const role = ref<'loket' | 'medis'>('loket')
  const staffPoli = ref<string>('Poli Umum')
  const selectedPatientId = ref<string>('demo-1')
  const verifyPatient = ref<Patient | null>(null)
  const lookupOpen = ref<boolean>(false)
  const search = ref<string>('')

  // Sync patients to localStorage
  watch(
    patients,
    (newVal) => {
      localStorage.setItem('antrisehat-patients', JSON.stringify(newVal))
    },
    { deep: true },
  )

  // Listen for storage events across tabs
  if (typeof window !== 'undefined') {
    window.addEventListener('storage', (e: StorageEvent) => {
      if (e.key === 'antrisehat-patients' && e.newValue) {
        try {
          patients.value = JSON.parse(e.newValue)
        } catch {
          // ignore parsing error
        }
      }
    })
  }

  // QR code generation when ticket changes
  watch(
    ticket,
    (newTicket) => {
      if (newTicket) {
        localStorage.setItem('antrisehat-ticket', newTicket)
        if (typeof window !== 'undefined') {
          QRCode.toDataURL(
            `${window.location.origin}${window.location.pathname}?checkin=${encodeURIComponent(newTicket)}`,
            { width: 240, margin: 1, color: { dark: '#195A5F' } },
          ).then((url: string) => {
            qr.value = url
          })
        }
      } else {
        localStorage.removeItem('antrisehat-ticket')
        qr.value = ''
      }
    },
    { immediate: true },
  )

  const setNotice = (msg: string) => {
    notice.value = msg
    if (noticeTimer) clearTimeout(noticeTimer)
    noticeTimer = setTimeout(() => {
      notice.value = ''
    }, 5000)
  }

  const currentTicket = computed(() => patients.value.find((p) => p.id === ticket.value))

  const activePatient = computed(() => patients.value.find((p) => p.status === 'Dipanggil'))

  const chosenPatient = computed(() => patients.value.find((p) => p.id === selectedPatientId.value))

  const changeStatus = (id: string, status: string) => {
    patients.value = patients.value.map((p) => (p.id === id ? { ...p, status } : p))
  }

  const callPatient = (p: Patient) => {
    patients.value = patients.value.map((item) =>
      item.id === p.id
        ? { ...item, status: 'Dipanggil' }
        : item.status === 'Dipanggil'
          ? { ...item, status: 'Hadir' }
          : item,
    )
    setNotice(`${p.number} dipanggil ke ${p.poli}`)

    if (typeof window !== 'undefined' && 'speechSynthesis' in window) {
      const speech = new SpeechSynthesisUtterance(
        `Nomor antrean ${p.number.replace('-', ' ')}, silakan menuju ${p.poli}`,
      )
      speech.lang = 'id-ID'
      window.speechSynthesis.speak(speech)
    }
  }

  const registerPatient = (
    formData: Omit<Patient, 'id' | 'number' | 'status' | 'source'>,
    walkin = false,
  ) => {
    const id = crypto.randomUUID()
    const prefix =
      formData.poli === 'Poli Gigi'
        ? 'B'
        : formData.poli === 'Poli KIA'
          ? 'C'
          : formData.poli === 'Poli Lansia'
            ? 'D'
            : 'A'

    const existingNums = patients.value
      .filter((p) => p.number.startsWith(prefix))
      .map((p) => Number(p.number.split('-')[1]))

    const nextNum = Math.max(11, ...existingNums, 0) + 1
    const number = `${prefix}-${String(nextNum).padStart(3, '0')}`

    const newPatient: Patient = {
      ...formData,
      id,
      number,
      source: walkin ? 'Langsung' : 'Online',
      status: walkin ? 'Hadir' : 'Terdaftar',
    }

    patients.value = [...patients.value, newPatient]
    ticket.value = id
    return newPatient
  }

  const downloadTicketCanvas = () => {
    if (!currentTicket.value || !qr.value) return
    const canvas = document.createElement('canvas')
    canvas.width = 720
    canvas.height = 940
    const ctx = canvas.getContext('2d')
    if (!ctx) return
    ctx.fillStyle = '#F3EDE6'
    ctx.fillRect(0, 0, 720, 940)
    ctx.fillStyle = '#195A5F'
    ctx.textAlign = 'center'
    ctx.font = 'bold 36px sans-serif'
    ctx.fillText('AntriSehat · E-Tiket', 360, 80)
    ctx.font = 'bold 96px sans-serif'
    ctx.fillText(currentTicket.value.number, 360, 210)
    const img = new Image()
    img.onload = () => {
      ctx.drawImage(img, 240, 270, 240, 240)
      ctx.font = '24px sans-serif'
      ;[
        currentTicket.value!.name,
        currentTicket.value!.poli,
        currentTicket.value!.date,
        `Status: ${currentTicket.value!.status}`,
        'Bawa KTP dan kartu BPJS saat kunjungan',
        'Scan QR di lokasi untuk check-in',
      ].forEach((line, i) => ctx.fillText(line, 360, 575 + i * 48))
      const link = document.createElement('a')
      link.download = `AntriSehat-${currentTicket.value!.number}.png`
      link.href = canvas.toDataURL()
      link.click()
    }
    img.src = qr.value
  }

  return {
    patients,
    ticket,
    qr,
    notice,
    role,
    staffPoli,
    selectedPatientId,
    verifyPatient,
    lookupOpen,
    search,
    currentTicket,
    activePatient,
    chosenPatient,
    setNotice,
    changeStatus,
    callPatient,
    registerPatient,
    downloadTicketCanvas,
  }
})
