// API Service client for SangPalito AntriSehat
const BASE_URL = '/api/v1'

export interface ApiQueueItem {
  id: string
  queue_number: number
  queue_code: string
  status: string
  queue_date: string
  call_count: number
  checked_in_at?: string
  called_at?: string
  service_start?: string
  service_end?: string
  polyclinic: {
    id: string
    name: string
    queue_prefix: string
  }
  patient: {
    id: string
    full_name: string
    medical_record_no: string
  }
}

export interface ApiVisitItem {
  id: string
  visit_date: string
  doctor_name?: string
  doctor_specialty?: string
  polyclinic_name?: string
  symptoms?: string
  vitals?: {
    bp_sys?: number
    bp_dia?: number
    blood_pressure?: string
    hr?: number
    heart_rate?: number
    respiratory_rate?: number
    temp?: number
    temperature?: number
    spo2?: number
    weight?: number
    height?: number
  }
  physical_exam?: string
  diagnosis?: string
  icd_code?: string
  diagnoses?: Array<{
    id: string
    icd10_code: string
    description?: string
    is_primary: boolean
    notes?: string
  }>
  treatment_plan?: string
  prescription?: string
  created_at: string
}

const getHeaders = () => {
  const token = localStorage.getItem('antrisehat-token')
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
  }
  if (token) {
    headers['Authorization'] = `Bearer ${token}`
  }
  return headers
}

async function authFetch(url: string, options: RequestInit = {}): Promise<Response> {
  const headers = { ...getHeaders(), ...(options.headers as Record<string, string> || {}) }
  let res = await fetch(url, { ...options, headers })

  // Jika token expired / 401 karena inaktivitas, coba refresh secara otomatis
  if (res.status === 401) {
    const newToken = await api.refreshToken()
    if (newToken) {
      const retryHeaders = {
        ...getHeaders(),
        ...(options.headers as Record<string, string> || {}),
        'Authorization': `Bearer ${newToken}`,
      }
      res = await fetch(url, { ...options, headers: retryHeaders })
    } else {
      api.clearSession()
      // Jika berada di halaman staff, arahkan ke login
      if (window.location.pathname.startsWith('/staff')) {
        window.location.href = '/login'
      }
    }
  }

  return res
}

export const api = {
  // ── 1. Antrean (Queues) ──────────────────────────────────
  async getQueuesToday(): Promise<ApiQueueItem[]> {
    try {
      const res = await authFetch(`${BASE_URL}/queues`)
      if (!res.ok) return []
      const json = await res.json()
      return json.data || []
    } catch {
      return []
    }
  },

  async updateQueueStatus(queueId: string, status: string): Promise<boolean> {
    try {
      const res = await authFetch(`${BASE_URL}/queues/${queueId}/status`, {
        method: 'PATCH',
        body: JSON.stringify({ status }),
      })
      return res.ok
    } catch {
      return false
    }
  },

  async getQueueEstimate(queueId: string) {
    try {
      const res = await fetch(`${BASE_URL}/public/queues/${queueId}/estimate`)
      if (!res.ok) return null
      const json = await res.json()
      return json.data || null
    } catch {
      return null
    }
  },

  // ── 2. Pendaftaran Pasien & Appointment ──────────────────
  async createAppointment(data: {
    nik: string
    full_name: string
    birth_date: string
    gender: 'male' | 'female'
    phone?: string
    guarantor: 'bpjs' | 'umum' | 'asuransi' | 'lainnya'
    bpjs_no?: string
    polyclinic_id: string
    schedule_id?: string
    visit_date: string
    chief_complaint?: string
    symptoms?: string[]
    duration_days?: number
  }) {
    try {
      const res = await fetch(`${BASE_URL}/public/appointments`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(data),
      })
      const json = await res.json()
      return json
    } catch (err) {
      console.error('Failed to create appointment:', err)
      return { success: false, message: 'Gagal terhubung ke backend server' }
    }
  },

  async checkInByQR(qrToken: string) {
    try {
      const res = await fetch(`${BASE_URL}/public/queues/checkin`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ qr_token: qrToken }),
      })
      return await res.json()
    } catch (err) {
      console.error('Check-in failed:', err)
      return null
    }
  },

  async getDemoQueues() {
    try {
      const res = await fetch(`${BASE_URL}/public/queues/demo`)
      if (!res.ok) return []
      const json = await res.json()
      return json.data || []
    } catch (err) {
      console.error('Failed to get demo queues:', err)
      return []
    }
  },

  async resetDemoQueue(qrToken = 'all') {
    try {
      const res = await fetch(`${BASE_URL}/public/queues/reset-demo`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ qr_token: qrToken }),
      })
      return await res.json()
    } catch (err) {
      console.error('Failed to reset demo queue:', err)
      return null
    }
  },

  // ── 3. Rekam Medis (Visits & Clinical) ────────────────────
  async getPatientVisits(patientId: string): Promise<ApiVisitItem[]> {
    try {
      const res = await authFetch(`${BASE_URL}/patients/${patientId}/visits`)
      if (!res.ok) return []
      const json = await res.json()
      return json.data || []
    } catch {
      return []
    }
  },

  async createVisit(payload: {
    appointment_id: string
    patient_id: string
    doctor_id: string
    vitals: Record<string, any>
    subjective: string
    objective?: string
    assessment: string
    plan: string
    diagnoses: Array<{
      icd10_code: string
      is_primary: boolean
      notes?: string
    }>
  }) {
    try {
      const res = await authFetch(`${BASE_URL}/visits`, {
        method: 'POST',
        body: JSON.stringify(payload),
      })
      return await res.json()
    } catch (err) {
      console.error('Failed to record visit:', err)
      return null
    }
  },

  async searchICD10(query: string, limit = 20) {
    try {
      const res = await authFetch(`${BASE_URL}/icd10?q=${encodeURIComponent(query)}&limit=${limit}`)
      if (!res.ok) return []
      const json = await res.json()
      return json.data || []
    } catch {
      return []
    }
  },

  // ── 4. Autentikasi Staf ──────────────────────────────────
  async login(email: string, password: string) {
    try {
      const res = await fetch(`${BASE_URL}/auth/login`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email, password }),
      })
      const json = await res.json()
      if (res.ok && json.data?.access_token) {
        localStorage.setItem('antrisehat-token', json.data.access_token)
        if (json.data.refresh_token) {
          localStorage.setItem('antrisehat-refresh-token', json.data.refresh_token)
        }
        localStorage.setItem('antrisehat-user', JSON.stringify(json.data.user || {}))
      }
      return json
    } catch (err) {
      console.error('Login error:', err)
      return { success: false, error: 'Koneksi ke backend gagal' }
    }
  },

  async refreshToken(): Promise<string | null> {
    const refreshToken = localStorage.getItem('antrisehat-refresh-token')
    if (!refreshToken) return null
    try {
      const res = await fetch(`${BASE_URL}/auth/refresh`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ refresh_token: refreshToken }),
      })
      if (!res.ok) {
        this.clearSession()
        return null
      }
      const json = await res.json()
      if (json.data?.access_token) {
        localStorage.setItem('antrisehat-token', json.data.access_token)
        return json.data.access_token
      }
      return null
    } catch (err) {
      console.error('Refresh token error:', err)
      return null
    }
  },

  async logout() {
    try {
      await fetch(`${BASE_URL}/auth/logout`, {
        method: 'POST',
        headers: getHeaders(),
      })
    } catch (err) {
      console.error('Logout error:', err)
    } finally {
      this.clearSession()
    }
  },

  clearSession() {
    localStorage.removeItem('antrisehat-token')
    localStorage.removeItem('antrisehat-refresh-token')
    localStorage.removeItem('antrisehat-user')
  },

  getCurrentUser() {
    try {
      const str = localStorage.getItem('antrisehat-user')
      return str ? JSON.parse(str) : null
    } catch {
      return null
    }
  },

  isAuthenticated(): boolean {
    return !!localStorage.getItem('antrisehat-token')
  },
}
