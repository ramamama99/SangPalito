# 🧠 Implementasi Smart Wait-Time (EMA + Kompleksitas) & Rekam Medis Pasien

Dokumen ini menjelaskan implementasi teknis dua fitur utama yang telah diintegrasikan ke dalam backend **SangPalito (AntriSehat)**:
1. **Smart Wait-Time Engine**: Estimasi waktu tunggu adaptif berbasis **Exponential Moving Average (EMA)** dan **Faktor Kompleksitas Klinis Pasien**.
2. **Log / Track Record Rekam Medis Pasien**: Riwayat klinis lengkap (gejala awal, tanda vital, SOAP, diagnosa ICD-10, resep obat, dan dokter pemeriksa).

---

## ⏱️ 1. Smart Wait-Time Engine (EMA + Faktor Kompleksitas)

### 1.1 Masalah pada Sistem Antrean Konvensional
Pada sistem antrean puskesmas biasa, estimasi waktu tunggu umumnya hanya berupa perkalian statis:
> `Waktu Tunggu = Jumlah Antrean di Depan × Rata-Rata Waktu Standar (misal 10 menit)`

Pendekatan statis ini sering kali meleset jauh karena:
- **Ritme dokter berubah-ubah**: Dokter bisa lebih cepat atau lebih lambat tergantung kelelahan atau beban kerja.
- **Kondisi pasien tidak seragam**: Pasien lansia, balita, atau pasien gawat darurat membutuhkan waktu konsultasi jauh lebih lama dibanding pasien dewasa yang hanya kontrol rutin.

---

### 1.2 Formulasi & Logika Algoritma

Algoritma Smart Wait-Time menggabungkan dua komponen utama:
1. **EMA Durasi Pelayanan Dokter** (mengikuti ritme aktual hari ini).
2. **Koefisien Kompleksitas Klinis** (mengikuti kondisi spesifik tiap pasien).

---

#### A. Exponential Moving Average (EMA) Durasi Pelayanan

Sistem menghitung rata-rata bergerak (*moving average*) durasi pelayanan dokter hari ini berdasarkan pasien yang telah selesai diperiksa (`status = 'done'`).

##### 1. Hitung Durasi Aktual Pasien Selesai:
```text
Durasi_Aktual = Waktu_Selesai_Layanan - Waktu_Mulai_Layanan (dalam menit)
```

##### 2. Batas Nilai (Pencegahan Outlier):
Untuk mencegah anomali (misal dokter lupa menutup antrean di sistem atau sesi konsultasi terlalu singkat), durasi aktual dibatasi antara **3 menit hingga 45 menit**:
```text
Durasi_Clamped = min(45.0, max(3.0, Durasi_Aktual))
```

##### 3. Rumus Pembaruan EMA:
```text
EMA_baru = (α × Durasi_Clamped) + ((1 - α) × EMA_sebelumnya)
```

- **Faktor Penghalus (Smoothing Factor α)**: `α = 0.30` (memberikan bobot 30% ke pasien terakhir dan 70% ke riwayat sebelumnya).
- **Nilai Awal (Baseline)**: Jika belum ada pasien selesai hari ini, sistem menggunakan nilai standar poliklinik (`polyclinics.avg_service_min`, default 10–15 menit).

---

#### B. Faktor Kompleksitas Klinis Pasien (Cp)

Setiap pasien memiliki koefisien pengali kompleksitas individual `Cp` (nilai dasar = `1.0`) berdasarkan data usia, hasil triase, dan keluhan awal:

| Parameter | Kategori / Kondisi | Tambahan Koefisien (+Δ) | Dampak Waktu |
| :--- | :--- | :--- | :--- |
| **Usia Pasien** | Geriatri / Lansia (≥ 60 tahun) | +0.30 | +30% lebih lama |
| **Usia Pasien** | Pediatri / Balita (< 5 tahun) | +0.20 | +20% lebih lama |
| **Usia Pasien** | Anak-anak (5 – 12 tahun) | +0.10 | +10% lebih lama |
| **Usia Pasien** | Dewasa (13 – 59 tahun) | +0.00 | Standar |
| **Triase & Prioritas** | Emergency / Triase Merah | +0.40 | +40% lebih lama |
| **Triase & Prioritas** | Urgent / Triase Kuning | +0.25 | +25% lebih lama |
| **Triase & Prioritas** | Normal / Triase Hijau | +0.00 | Standar |
| **Jumlah Keluhan** | Keluhan Tinggi (≥ 4 gejala) | +0.25 | +25% lebih lama |
| **Jumlah Keluhan** | Multi-gejala (≥ 2 gejala / deskripsi keluhan panjang) | +0.15 | +15% lebih lama |

> **Batas Pengali (Clamp Range)**: Koefisien total dibatasi pada rentang aman **0.7 ≤ Cp ≤ 2.5**.

---

#### C. Kalkulasi Total Waktu Tunggu

Untuk menghitung estimasi waktu tunggu pasien target, sistem mengalikan nilai EMA dengan koefisien kompleksitas masing-masing pasien:

```text
Estimasi Durasi Pasien = EMA × Faktor_Kompleksitas (Cp)
```

##### 1. Pasien yang Sedang Diperiksa (`in_service`):
- Estimasi durasi total: `Durasi_Aktif = EMA × Cp_aktif`
- Waktu yang sudah berjalan: `Waktu_Berjalan = Waktu_Sekarang - Waktu_Mulai_Layanan`
- Sisa waktu konsultasi: `Sisa_Waktu = max(1.0, Durasi_Aktif - Waktu_Berjalan)`

##### 2. Pasien-Pasien Antrean di Depan (`waiting`, `checked_in`, `called`):
Untuk setiap antrean pasien di depan, dihitung estimasi durasinya:
```text
Durasi_Antrean = EMA × Cp
```

##### 3. Total Waktu Tunggu & Jam Panggilan:
```text
Total_Menit_Tunggu = Sisa_Waktu_Pasien_Aktif + Total(Durasi_Semua_Antrean_di_Depan)
Estimasi_Jam_Panggil = Waktu_Sekarang + Total_Menit_Tunggu
```

---

#### D. Contoh Simulasi Perhitungan Riil

Misalkan:
- Nilai **EMA Dokter** saat ini = **10.0 menit**.
- **Pasien A** sedang diperiksa: Lansia (`Cp = 1.3`), konsultasi sudah berjalan 6 menit.
  - Estimasi durasi = 10.0 × 1.3 = 13.0 menit.
  - Sisa waktu = 13.0 - 6.0 = **7.0 menit**.
- **Pasien B** di depan Anda: Balita (`Cp = 1.2`).
  - Estimasi durasi = 10.0 × 1.2 = **12.0 menit**.
- **Pasien C** (Anda): Dewasa, keluhan tunggal (`Cp = 1.0`).
- **Hasil Akhir**:
  - Total waktu tunggu Anda = 7.0 menit (sisa Pasien A) + 12.0 menit (Pasien B) = **19 menit**.

---

### 1.3 Endpoint API Estimasi Antrean

- **Metode**: `GET`
- **URL**: `/api/v1/public/queues/:id/estimate`
- **Contoh Response JSON**:
```json
{
  "success": true,
  "data": {
    "queue_id": "7fa12bf2-2b63-4bf8-b7ae-17ddc24cba19",
    "queue_number": 8,
    "queue_code": "A-008",
    "current_serving": 5,
    "ahead_count": 2,
    "estimate_minutes": 27,
    "ema_service_duration_minutes": 11.4,
    "complexity_factor": 1.3,
    "complexity_notes": [
      "Pasien Geriatri/Lansia (≥60 thn): +30%"
    ],
    "estimated_call_time": "2026-10-15T09:42:00+07:00"
  }
}
```

---

## 📋 2. Log & Track Record Rekam Medis Pasien

Riwayat pemeriksaan pasien disimpan secara terstruktur menggunakan standar medis (SOAP & ICD-10) sehingga dokter dapat meninjau riwayat penyakit terdahulu secara berkesinambungan.

### 2.1 Struktur Data Rekam Medis
Setiap kunjungan medis mencatat informasi menyeluruh:
- **Tanggal Kunjungan**: `visit_date`, `created_at`
- **Tenaga Medis**: `doctor_name`, `doctor_specialty`, `polyclinic_name`
- **Subjective (S)**: Gejala & keluhan awal yang dirasakan pasien
- **Objective (O)**: 
  - Hasil pemeriksaan fisik dokter
  - Tanda vital (tekanan darah, detak jantung, laju pernapasan, suhu, berat badan, tinggi badan)
- **Assessment (A)**: 
  - Diagnosa klinis dokter
  - Daftar kode ICD-10 (kode standar, deskripsi penyakit, status primer/sekunder)
- **Plan (P)**: Edukasi, rencana tindakan medis, dan resep obat

---

### 2.2 Daftar Endpoint Rekam Medis

| Metode | URL Endpoint | Hak Akses | Deskripsi & Kegunaan |
| :--- | :--- | :--- | :--- |
| `POST` | `/api/v1/clinical/visits` | Dokter (`medic`) | Mencatat hasil pemeriksaan, tanda vital, diagnosa ICD-10, resep obat, serta otomatis menyelesaikan antrean (`status: done`) |
| `GET` | `/api/v1/patients/:id/medical-records` | `medic`, `registrar` | Mengambil seluruh rekam medis pasien secara kronologis dari masa ke masa |
| `GET` | `/api/v1/patients/:id/visits` | `medic`, `registrar` | Riwayat kunjungan pasien lengkap dengan detail klinis |
| `GET` | `/api/v1/clinical/visits/:id` | `medic`, `registrar` | Mengambil detail 1 sesi kunjungan/rekam medis tertentu |
| `GET` | `/api/v1/public/appointments/:token/medical-record` | Publik (Token Pasien) | Pasien dapat melihat hasil diagnosa dan resep obatnya secara mandiri |
| `GET` | `/api/v1/clinical/icd10?q=ISPA` | `medic` | Pencarian master data kode ICD-10 untuk fitur autokomplit diagnosis dokter |

---

### 2.3 Contoh Payload: Input Rekam Medis (`POST /api/v1/clinical/visits`)

```json
{
  "appointment_id": "9b1deb4d-3b7d-4bad-9bdd-2b0d7b3dcb6d",
  "patient_id": "3a1deb4d-3b7d-4bad-9bdd-2b0d7b3dcb11",
  "vitals": {
    "blood_pressure": "120/80",
    "heart_rate": 84,
    "respiratory_rate": 20,
    "temperature": 38.2,
    "weight": 65.5,
    "height": 170.0
  },
  "subjective": "Pasien mengeluhkan demam sejak 3 hari yang lalu, sakit menelan, batuk kering, dan sakit kepala berdenyut.",
  "objective": "Keadaan umum: compos mentis, faring tampak hiperemis, pembesaran tonsil T2/T1 tanpa detritus.",
  "assessment": "Faringitis Akut (ISPA Bakterial)",
  "plan": "Edukasi hidrasi cairan hangat, istirahat cukup, terapi simtomatik dan antibiotik oral.",
  "prescription": "1. Paracetamol tab 500 mg No. X (3 x 1 tablet prn)\n2. Amoxicillin tab 500 mg No. XV (3 x 1 tablet habiskan)",
  "diagnoses": [
    {
      "icd10_code": "J02.9",
      "description": "Faringitis Akut Tidak Spesifik",
      "is_primary": true,
      "notes": "Diagnosis primer"
    },
    {
      "icd10_code": "R50.9",
      "description": "Demam Tidak Spesifik (Febris)",
      "is_primary": false
    }
  ],
  "complete_queue": true
}
```

---

### 2.4 Contoh Response: Riwayat Pasien (`GET /api/v1/patients/:id/medical-records`)

```json
{
  "success": true,
  "data": [
    {
      "id": "e4f8c6b2-6c9f-4f3b-8f1b-5e6a7b8c9d0e",
      "visit_date": "2026-10-15",
      "doctor_name": "dr. Budi Santoso, Sp.PD",
      "doctor_specialty": "Penyakit Dalam",
      "polyclinic_name": "Poli Umum",
      "symptoms": "Demam sejak 3 hari yang lalu, sakit menelan, batuk kering",
      "vitals": {
        "blood_pressure": "120/80",
        "heart_rate": 84,
        "respiratory_rate": 20,
        "temperature": 38.2,
        "weight": 65.5,
        "height": 170.0
      },
      "physical_exam": "Faring tampak hiperemis, tonsil T2/T1",
      "diagnosis": "Faringitis Akut (ISPA Bakterial)",
      "diagnoses": [
        {
          "id": "a1b2c3d4-1111-2222-3333-444455556666",
          "icd10_code": "J02.9",
          "description": "Faringitis Akut Tidak Spesifik",
          "is_primary": true
        }
      ],
      "treatment_plan": "Edukasi hidrasi cairan hangat, istirahat cukup",
      "prescription": "Paracetamol 500mg 3x1, Amoxicillin 500mg 3x1",
      "created_at": "2026-10-15T09:30:00Z"
    }
  ]
}
```

---

## 🛡️ 3. Pengujian & Validasi

- **Unit Testing Backend**:
  - `TestCalculateEMA`: Memvalidasi formula smoothing α = 0.3, penyesuaian durasi aktual, dan batas outlier (3 – 45 menit).
  - `TestCalculateComplexityFactor`: Memvalidasi koefisien usia (lansia ≥ 60 thn +30%, balita +20%), triase merah/kuning (+40%/+25%), serta multi-gejala.
  - `TestClinicalService_RecordVisitAndRetrieveHistory`: Memvalidasi pencatatan rekam medis, otomatisasi penutupan status antrean, pencarian ICD-10, dan query riwayat medis.
- **Hasil Test**: Status **PASS** (100% lulus seluruh test case).
- **Status Build**: Kompilasi binary API sukses tanpa error (`go build -v -o bin/api.exe ./cmd/api`).
