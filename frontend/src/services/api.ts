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
  prescription?: any
  prescription_no?: string
  medication_plan?: string
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
    medication_plan?: string
    diagnoses: Array<{
      icd10_code: string
      is_primary: boolean
      notes?: string
    }>
    [key: string]: any
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

  async recordVisitAuditLog(payload: {
    visit_id?: string
    appointment_id?: string
    patient_id?: string
    patient_name?: string
    polyclinic?: string
    doctor_name?: string
    diagnosis?: string
    prescription_no?: string
    vitals?: Record<string, any>
    subjective?: string
    objective?: string
    plan?: string
    medication_plan?: string
  }) {
    try {
      const res = await authFetch(`${BASE_URL}/visits/audit-log`, {
        method: 'POST',
        body: JSON.stringify(payload),
      })
      return await res.json()
    } catch (err) {
      console.error('Failed to record visit audit log:', err)
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

  // ── 4. Resep & Obat (Prescriptions & Medications) ───────
  async getMedications(search = '', category = '', route = '') {
    try {
      const params = new URLSearchParams()
      if (search) params.append('search', search)
      if (category) params.append('category', category)
      if (route) params.append('route', route)
      const res = await fetch(`${BASE_URL}/medications?${params.toString()}`)
      if (!res.ok) return []
      const json = await res.json()
      return json.data || []
    } catch {
      return []
    }
  },

  async createPrescription(payload: any) {
    try {
      const res = await fetch(`${BASE_URL}/prescriptions`, {
        method: 'POST',
        headers: getHeaders(),
        body: JSON.stringify(payload),
      })
      return await res.json()
    } catch (err) {
      console.error('Failed to create prescription:', err)
      return null
    }
  },

  async getPrescriptionByNo(prescriptionNo: string) {
    try {
      const res = await fetch(`${BASE_URL}/public/prescriptions/${encodeURIComponent(prescriptionNo)}`)
      if (!res.ok) return null
      const json = await res.json()
      return json.data || null
    } catch {
      return null
    }
  },

  async getPrescriptionByAppointmentToken(token: string) {
    try {
      const res = await fetch(`${BASE_URL}/public/appointments/${token}/prescription`)
      if (!res.ok) return null
      const json = await res.json()
      return json.data || null
    } catch {
      return null
    }
  },

  async dispensePrescription(id: string) {
    try {
      const res = await fetch(`${BASE_URL}/prescriptions/${id}/dispense`, {
        method: 'PATCH',
        headers: getHeaders(),
      })
      return await res.json()
    } catch (err) {
      console.error('Failed to dispense prescription:', err)
      return null
    }
  },

  async createMedication(payload: any) {
    try {
      const res = await fetch(`${BASE_URL}/medications`, {
        method: 'POST',
        headers: getHeaders(),
        body: JSON.stringify(payload),
      })
      return await res.json()
    } catch (err) {
      console.error('Failed to create medication:', err)
      return null
    }
  },

  async updateMedication(id: string, payload: any) {
    try {
      const res = await fetch(`${BASE_URL}/medications/${id}`, {
        method: 'PUT',
        headers: getHeaders(),
        body: JSON.stringify(payload),
      })
      return await res.json()
    } catch (err) {
      console.error('Failed to update medication:', err)
      return null
    }
  },

  async deleteMedication(id: string) {
    try {
      const res = await fetch(`${BASE_URL}/medications/${id}`, {
        method: 'DELETE',
        headers: getHeaders(),
      })
      return await res.json()
    } catch (err) {
      console.error('Failed to delete medication:', err)
      return null
    }
  },

  async searchPatients(query = '', page = 1, limit = 50) {
    try {
      const res = await fetch(`${BASE_URL}/patients?q=${encodeURIComponent(query)}&page=${page}&limit=${limit}`, {
        headers: getHeaders(),
      })
      if (!res.ok) return []
      const json = await res.json()
      return json.data || []
    } catch {
      return []
    }
  },

  async createPatient(payload: any) {
    try {
      const res = await fetch(`${BASE_URL}/patients`, {
        method: 'POST',
        headers: getHeaders(),
        body: JSON.stringify(payload),
      })
      return await res.json()
    } catch (err) {
      console.error('Failed to create patient:', err)
      return null
    }
  },

  async updatePatient(id: string, payload: any) {
    try {
      const res = await fetch(`${BASE_URL}/patients/${id}`, {
        method: 'PATCH',
        headers: getHeaders(),
        body: JSON.stringify(payload),
      })
      return await res.json()
    } catch (err) {
      console.error('Failed to update patient:', err)
      return null
    }
  },

  async deletePatient(id: string) {
    try {
      const res = await fetch(`${BASE_URL}/patients/${id}`, {
        method: 'DELETE',
        headers: getHeaders(),
      })
      return await res.json()
    } catch (err) {
      console.error('Failed to delete patient:', err)
      return null
    }
  },

  // ── 5. Autentikasi Staf ──────────────────────────────────
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

  async forgotPassword(email: string) {
    try {
      const res = await fetch(`${BASE_URL}/auth/forgot-password`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email }),
      })
      return await res.json()
    } catch (err) {
      console.error('ForgotPassword error:', err)
      return { success: false, error: 'Koneksi ke backend gagal' }
    }
  },

  async resetPasswordWithToken(token: string, newPassword: string) {
    try {
      const res = await fetch(`${BASE_URL}/auth/reset-password`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ token, new_password: newPassword }),
      })
      return await res.json()
    } catch (err) {
      console.error('ResetPassword error:', err)
      return { success: false, error: 'Koneksi ke backend gagal' }
    }
  },

  async completeFirstLogin(data: { new_password: string; phone?: string; str_no?: string }) {
    try {
      const res = await authFetch(`${BASE_URL}/auth/first-login`, {
        method: 'POST',
        body: JSON.stringify(data),
      })
      return await res.json()
    } catch (err) {
      console.error('CompleteFirstLogin error:', err)
      return { success: false, error: 'Koneksi ke backend gagal' }
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

  // ── 6. Master Poliklinik & Dokter ──────────────────────────
  async getPolyclinics() {
    try {
      const res = await authFetch(`${BASE_URL}/polyclinics`)
      return await res.json()
    } catch (err) {
      console.error('Failed to get polyclinics:', err)
      return { success: false, data: [] }
    }
  },

  async getDoctors() {
    try {
      const res = await authFetch(`${BASE_URL}/doctors`)
      return await res.json()
    } catch (err) {
      console.error('Failed to get doctors:', err)
      return { success: false, data: [] }
    }
  },

  async createDoctor(data: {
    full_name: string
    email: string
    phone: string
    specialty?: string
    str_no?: string
    polyclinic_id?: string
    password?: string
  }) {
    try {
      const res = await authFetch(`${BASE_URL}/doctors`, {
        method: 'POST',
        body: JSON.stringify(data),
      })
      return await res.json()
    } catch (err) {
      console.error('Failed to create doctor:', err)
      return { success: false, error: 'Koneksi ke backend gagal' }
    }
  },

  async updateDoctor(id: string, data: {
    full_name: string
    email?: string
    phone?: string
    specialty?: string
    str_no?: string
    polyclinic_id?: string
    is_active: boolean
  }) {
    try {
      const res = await authFetch(`${BASE_URL}/doctors/${id}`, {
        method: 'PUT',
        body: JSON.stringify(data),
      })
      return await res.json()
    } catch (err) {
      console.error('Failed to update doctor:', err)
      return { success: false, error: 'Koneksi ke backend gagal' }
    }
  },

  async deleteDoctor(id: string) {
    try {
      const res = await authFetch(`${BASE_URL}/doctors/${id}`, {
        method: 'DELETE',
      })
      return await res.json()
    } catch (err) {
      console.error('Failed to delete doctor:', err)
      return { success: false, error: 'Koneksi ke backend gagal' }
    }
  },

  async getSchedules(polyclinicId?: string, date = 'all') {
    try {
      const params = new URLSearchParams()
      if (polyclinicId) params.append('polyclinic_id', polyclinicId)
      if (date) params.append('date', date)
      const res = await fetch(`${BASE_URL}/schedules?${params.toString()}`)
      if (!res.ok) return { success: false, data: [] }
      return await res.json()
    } catch (err) {
      console.error('Failed to get schedules:', err)
      return { success: false, data: [] }
    }
  },

  async assignShiftSchedule(data: {
    id?: string
    polyclinic_id: string
    doctor_id: string
    day_of_week: number
    shift_id?: string
    start_time?: string
    end_time?: string
    quota?: number
  }) {
    try {
      const res = await authFetch(`${BASE_URL}/schedules/assign`, {
        method: 'POST',
        body: JSON.stringify(data),
      })
      return await res.json()
    } catch (err) {
      console.error('Failed to assign schedule:', err)
      return { success: false, error: 'Koneksi ke backend gagal' }
    }
  },

  async unassignShiftSchedule(data: {
    schedule_id?: string
    polyclinic_id?: string
    day_of_week?: number
    shift_id?: string
    start_time?: string
  }) {
    try {
      const res = await authFetch(`${BASE_URL}/schedules/unassign`, {
        method: 'POST',
        body: JSON.stringify(data),
      })
      return await res.json()
    } catch (err) {
      console.error('Failed to unassign schedule:', err)
      return { success: false, error: 'Koneksi ke backend gagal' }
    }
  },

  async deleteSchedule(id: string) {
    try {
      const res = await authFetch(`${BASE_URL}/schedules/${id}`, {
        method: 'DELETE',
      })
      return await res.json()
    } catch (err) {
      console.error('Failed to delete schedule:', err)
      return { success: false, error: 'Koneksi ke backend gagal' }
    }
  },

  // ── 7. Manajemen User (Admin) ────────────────────────────
  async getUsers(page = 1, limit = 50) {
    try {
      const res = await authFetch(`${BASE_URL}/users?page=${page}&limit=${limit}`)
      return await res.json()
    } catch (err) {
      console.error('Failed to get users:', err)
      return { success: false, data: [] }
    }
  },

  async createUser(payload: {
    email: string
    full_name: string
    phone?: string
    password: string
    role: string
    polyclinic_id?: string
  }) {
    try {
      const res = await authFetch(`${BASE_URL}/users`, {
        method: 'POST',
        body: JSON.stringify(payload),
      })
      return await res.json()
    } catch (err) {
      console.error('Failed to create user:', err)
      return { success: false, error: 'Koneksi ke backend gagal' }
    }
  },

  async updateUser(id: string, payload: {
    full_name?: string
    phone?: string
    role?: string
    polyclinic_id?: string
    is_active?: boolean
  }) {
    try {
      const res = await authFetch(`${BASE_URL}/users/${id}`, {
        method: 'PATCH',
        body: JSON.stringify(payload),
      })
      return await res.json()
    } catch (err) {
      console.error('Failed to update user:', err)
      return { success: false, error: 'Koneksi ke backend gagal' }
    }
  },

  async deleteUser(id: string) {
    try {
      const res = await authFetch(`${BASE_URL}/users/${id}`, {
        method: 'DELETE',
      })
      return await res.json()
    } catch (err) {
      console.error('Failed to delete user:', err)
      return { success: false, error: 'Koneksi ke backend gagal' }
    }
  },

  async resetUserPassword(id: string, newPassword: string) {
    try {
      const res = await authFetch(`${BASE_URL}/users/${id}/password`, {
        method: 'PATCH',
        body: JSON.stringify({ new_password: newPassword }),
      })
      return await res.json()
    } catch (err) {
      console.error('Failed to reset password:', err)
      return { success: false, error: 'Koneksi ke backend gagal' }
    }
  },

  // ── 8. Log Audit Realtime (Admin & Kepala Puskesmas) ──────
  async getAuditLogs(params: { entity?: string; entity_id?: string; page?: number; limit?: number } = {}) {
    try {
      const queryParams = new URLSearchParams()
      if (params.entity) queryParams.append('entity', params.entity)
      if (params.entity_id) queryParams.append('entity_id', params.entity_id)
      if (params.page) queryParams.append('page', String(params.page))
      if (params.limit) queryParams.append('limit', String(params.limit || 50))

      const res = await authFetch(`${BASE_URL}/audit-logs?${queryParams.toString()}`)
      return await res.json()
    } catch (err) {
      console.error('Failed to get audit logs:', err)
      return { success: false, data: [] }
    }
  },

  // ── 9. Dashboard Analitik (Kepala Puskesmas & Admin) ──────
  async getWeeklyPatients() {
    try {
      const res = await authFetch(`${BASE_URL}/analytics/weekly-patients`)
      return await res.json()
    } catch (err) {
      console.error('Failed to get weekly patients:', err)
      return { success: false, data: null }
    }
  },

  async getDoctorPerformance() {
    try {
      const res = await authFetch(`${BASE_URL}/analytics/doctor-performance`)
      return await res.json()
    } catch (err) {
      console.error('Failed to get doctor performance:', err)
      return { success: false, data: [] }
    }
  },

  async getAnalyticsSummary(dateFrom?: string, dateTo?: string) {
    try {
      const query = new URLSearchParams()
      if (dateFrom) query.append('date_from', dateFrom)
      if (dateTo) query.append('date_to', dateTo)
      const res = await authFetch(`${BASE_URL}/analytics/summary?${query.toString()}`)
      return await res.json()
    } catch (err) {
      console.error('Failed to get analytics summary:', err)
      return { success: false, data: null }
    }
  },

  async getBusyHours(polyclinicId?: string, dateFrom?: string, dateTo?: string) {
    try {
      const query = new URLSearchParams()
      if (polyclinicId) query.append('polyclinic_id', polyclinicId)
      if (dateFrom) query.append('date_from', dateFrom)
      if (dateTo) query.append('date_to', dateTo)
      const res = await authFetch(`${BASE_URL}/analytics/busy-hours?${query.toString()}`)
      return await res.json()
    } catch (err) {
      console.error('Failed to get busy hours:', err)
      return { success: false, data: [] }
    }
  },

  async getTopDiagnoses(limit = 10, dateFrom?: string, dateTo?: string) {
    try {
      const query = new URLSearchParams()
      query.append('limit', String(limit))
      if (dateFrom) query.append('date_from', dateFrom)
      if (dateTo) query.append('date_to', dateTo)
      const res = await authFetch(`${BASE_URL}/analytics/top-diagnoses?${query.toString()}`)
      return await res.json()
    } catch (err) {
      console.error('Failed to get top diagnoses:', err)
      return { success: false, data: [] }
    }
  },
}
