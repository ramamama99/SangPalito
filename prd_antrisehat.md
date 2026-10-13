*Referensi Dokumen Asli: PRD - Sistem Antrean dan Data Pasien Puskesmas Klinik.pdf*

# PRD: Sistem Antrean dan Manajemen Data Pasien Puskesmas/Klinik

**Nama produk (working title):** AntriSehat 
**Platform:** Web (Vue.js PWA) + REST API (Go) + PostgreSQL 
**Versi:** 1.0 (MVP)

---

## 1. Ringkasan
AntriSehat adalah platform web untuk puskesmas dan klinik yang mendigitalkan pendaftaran, antrean, pencatatan kunjungan, dan pengelolaan data pasien. Pasien mendaftar dan memantau antrean dari ponsel, sementara petugas bekerja dari dashboard terpusat. Pimpinan puskesmas mendapat analitik untuk keputusan operasional.

## 2. Latar Belakang dan Masalah

| Masalah | Dampak |
| :--- | :--- |
| Antrean panjang dan tidak transparan | Pasien menunggu lama di lokasi, risiko penularan, kepuasan rendah |
| Administrasi manual (formulir kertas, input ganda) | Petugas kewalahan, salah input, data sulit dicari |
| Data pasien terpencar dan belum optimal | Riwayat kunjungan sulit dilacak, pelaporan lambat, sulit integrasi SATUSEHAT |
| Tidak ada data operasional | Pimpinan tidak tahu jam sibuk dan beban tiap poli |

## 3. Tujuan dan Metrik Keberhasilan

| Tujuan | Metrik | Target (pilot) |
| :--- | :--- | :--- |
| Mengurangi waktu tunggu | Rata-rata waktu tunggu pasien | Turun minimal 30% |
| Mengurangi beban administrasi | Waktu pendaftaran per pasien | Dari ±5 menit menjadi kurang dari 2 menit |
| Meningkatkan akurasi data | Persentase data lengkap dan valid | Lebih dari 95% |
| Adopsi | Persentase pasien yang mendaftar online | Lebih dari 40% dalam 3 bulan |
| Keputusan berbasis data | Dashboard dipakai mingguan oleh pimpinan | Ya |

**Non-tujuan (MVP):** billing/klaim BPJS penuh, integrasi alat laboratorium, telemedisin video, rekam medis klinis lengkap (SOAP mendalam, resep elektronik penuh).

## 4. Pengguna dan Persona

| Peran | Deskripsi | Kebutuhan utama |
| :--- | :--- | :--- |
| **Pasien** | Warga yang berobat, termasuk lansia dan pengguna ponsel sederhana | Daftar mudah, tahu estimasi giliran, tidak menunggu di lokasi |
| **Petugas pendaftaran** | Admin loket | Cari pasien cepat, daftar kunjungan, kelola antrean |
| **Tenaga medis** (dokter, perawat, bidan) | Pemeriksa di poli | Lihat antrean poli, ringkasan keluhan, catat hasil pemeriksaan |
| **Kepala puskesmas/klinik** | Pengambil keputusan | Dashboard analitik dan laporan |
| **Admin sistem** | Pengelola teknis | Kelola pengguna, poli, jadwal, hak akses |

## 5. Ruang Lingkup Fitur

### 5.1 Fitur MVP (Prioritas PO)

| ID Fitur | Fitur | Deskripsi | Prioritas |
| :--- | :--- | :--- | :--- |
| F-01 | Autentikasi dan RBAC | Login staf (JWT + refresh token), hak akses per peran | PO |
| F-02 | Pendaftaran pasien baru | Input data dasar (NIK, nama, tanggal lahir, kontak, penjamin), deteksi duplikat via NIK | PO |
| F-03 | Pendaftaran kunjungan online | Pasien memilih poli, tanggal, dan jadwal dokter, lalu mendapat nomor antrean dan kode QR | PO |
| F-04 | Antrean real-time | Status antrean (menunggu, dipanggil, dilayani, selesai, batal), tampilan layar TV antrean, update via WebSocket/SSE | PO |
| F-05 | Estimasi waktu tunggu | Dihitung dari rata-rata durasi layanan poli dan jumlah antrean di depan | PO |
| F-06 | Notifikasi | WhatsApp/SMS saat tersisa 3 antrean dan saat dipanggil | PO |
| F-07 | Pra-anamnesis digital | Pasien mengisi keluhan utama dan gejala sebelum datang; petugas melihat ringkasan dan prioritas awal | PO |
| F-08 | Check-in QR | Pasien scan QR di lokasi untuk konfirmasi kehadiran | PO |
| F-09 | Pencatatan kunjungan | Tanda vital, keluhan, diagnosis (ICD-10), tindakan, catatan | PO |
| F-10 | Pencarian dan riwayat pasien | Cari berdasarkan NIK, nama, no. RM; lihat riwayat kunjungan | PO |
| F-11 | Dashboard analitik | Jam sibuk, kunjungan per poli, 10 penyakit terbanyak, rata-rata waktu tunggu dan layanan | PO |
| F-12 | Audit log | Rekam siapa mengakses/mengubah data pasien | PO |

### 5.2 Fitur Lanjutan (P1/P2)

| ID Fitur | Deskripsi | Prioritas |
| :--- | :--- | :--- |
| F-13 | Mode offline-first (PWA, antrean input lokal lalu sinkron) | P1 |
| F-14 | Ekspor/integrasi SATUSEHAT (FHIR) | P1 |
| F-15 | Triase otomatis berbasis aturan (red flag symptoms) | P1 |
| F-16 | Prediksi kepadatan dengan data historis | P2 |
| F-17 | Manajemen stok obat sederhana | P2 |
| F-18 | Multi-fasilitas (satu akun dinas kesehatan, banyak puskesmas) | P2 |

## 6. Alur Pengguna Utama

**Alur pasien**
1. Buka web, masuk dengan nomor HP (OTP) atau daftar baru.
2. Pilih poli dan jadwal, isi pra-anamnesis.
3. Terima nomor antrean, QR, dan estimasi waktu.
4. Datang saat mendekati giliran, scan QR untuk check-in.
5. Dipanggil, diperiksa, selesai.

**Alur petugas dan tenaga medis**
1. Petugas pendaftaran memverifikasi pasien (online maupun walk-in) dan menerbitkan antrean.
2. Perawat mencatat tanda vital.
3. Dokter memanggil pasien dari dashboard, membaca ringkasan pra-anamnesis, mengisi hasil pemeriksaan.
4. Kunjungan ditutup, data tersimpan dan masuk ke analitik.

## 7. Kebutuhan Fungsional Rinci (Ringkas)
* Satu NIK hanya boleh memiliki satu data pasien; pasien tanpa NIK (anak, lansia tertentu) dapat didaftarkan dengan identifier alternatif.
* Nomor antrean di-reset per poli per hari, format contoh A-012.
* Pasien dapat membatalkan atau menjadwal ulang maksimal H-1 jam sebelum slot.
* Status antrean berubah hanya mengikuti transisi valid (state machine): `waiting` → `checked_in` → `called` → `in_service` → `done`, serta `cancelled` dan `no_show`.
* Pasien yang tidak hadir saat dipanggil 3 kali ditandai `no_show`.
* Pasien prioritas (lansia, ibu hamil, disabilitas, gawat darurat) dapat diberi flag prioritas oleh petugas.

## 8. Kebutuhan Non-Fungsional

| Aspek | Kebutuhan |
| :--- | :--- |
| **Performa** | Respons API p95 di bawah 300 ms; update antrean di bawah 2 detik |
| **Ketersediaan** | Target 99,5%; backup harian dan point-in-time recovery |
| **Keamanan** | HTTPS, hash password argon2/bcrypt, enkripsi field sensitif (NIK) at-rest, RBAC, rate limiting, audit log |
| **Privasi** | Kepatuhan UU Pelindungan Data Pribadi (persetujuan, minimisasi data, hak akses/hapus) |
| **Kompatibilitas** | Responsif mobile-first, PWA, ringan di koneksi lambat |
| **Aksesibilitas** | Kontras tinggi, font besar, bahasa sederhana untuk lansia |
| **Skalabilitas** | Dirancang multi-tenant (kolom facility_id) |

## 9. Arsitektur Teknis

| Lapisan | Teknologi |
| :--- | :--- |
| **Frontend** | Vue 3 + Vite, Pinia, Vue Router, Tailwind CSS, PWA (vite-plugin-pwa), Chart.js/ECharts untuk dashboard |
| **Backend** | Go, router chi atau gin, sqlc atau pgx untuk query, golang-migrate untuk migrasi, JWT, validator |
| **Realtime** | WebSocket (gorilla/websocket) atau SSE untuk status antrean |
| **Database** | PostgreSQL 16 |
| **Cache/Queue** (opsional)| Redis untuk rate limit, sesi OTP, dan job notifikasi |
| **Notifikasi** | Gateway WhatsApp/SMS (mis. Fonnte, Twilio, atau Zenziva) via worker Go |
| **Deploy** | Docker Compose, Nginx reverse proxy, CI/CD GitHub Actions |

**Struktur backend (clean architecture ringkas)**
```text
/cmd/api            -> entrypoint
/internal/handler   -> HTTP handlers
/internal/service   -> business logic
/internal/repository-> query PostgreSQL
/internal/domain    -> entity dan interface
/internal/middleware-> auth, RBAC, logging
/internal/worker    -> notifikasi, job terjadwal
/migrations         -> file SQL migrasi
```

## 10. Struktur Database (PostgreSQL)

### 10.1 Ringkasan Entitas

| Tabel | Fungsi |
| :--- | :--- |
| `facilities` | Puskesmas/klinik (dasar multi-tenant) |
| `users` | Akun staf |
| `roles`, `user_roles` | Peran dan hak akses |
| `patients` | Data master pasien |
| `polyclinics` | Poli/unit layanan |
| `doctors` | Tenaga medis |
| `schedules` | Jadwal praktik per poli/dokter |
| `appointments` | Pendaftaran kunjungan |
| `queues` | Antrean harian |
| `pre_anamneses` | Keluhan awal dari pasien |
| `visits` | Kunjungan/pemeriksaan |
| `diagnoses` | Diagnosis per kunjungan (ICD-10) |
| `icd10_codes` | Master kode penyakit |
| `notifications` | Log notifikasi |
| `audit_logs` | Jejak akses/perubahan data |

### 10.2 Relasi Utama
* `facilities` 1-N `polyclinics`, `users`, `patients`, `doctors`
* `polyclinics` 1-N `schedules`; `doctors` 1-N `schedules`
* `patients` 1-N `appointments`; `appointments` 1-1 `queues`; `appointments` 1-0..1 `pre_anamneses`
* `appointments` 1-0..1 `visits`; `visits` 1-N `diagnoses`
* `users` 1-N `audit_logs`

### 10.3 DDL

```sql
CREATE EXTENSION IF NOT EXISTS "pgcrypto";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";

-- ENUM
CREATE TYPE gender_t AS ENUM ('male', 'female');
CREATE TYPE guarantor_t AS ENUM ('bpjs', 'umum', 'asuransi', 'lainnya');
CREATE TYPE queue_status_t AS ENUM ('waiting', 'checked_in', 'called', 'in_service', 'done', 'cancelled', 'no_show');
CREATE TYPE appt_source_t AS ENUM ('online', 'walk_in');
CREATE TYPE priority_t AS ENUM ('normal', 'elderly', 'pregnant', 'disabled', 'emergency');
CREATE TYPE notif_channel_t AS ENUM ('whatsapp', 'sms', 'email');
CREATE TYPE notif_status_t AS ENUM ('pending', 'sent', 'failed');

-- FASILITAS
CREATE TABLE facilities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(150) NOT NULL,
    type VARCHAR(30) NOT NULL, -- puskesmas | klinik
    code VARCHAR(30) UNIQUE, -- kode faskes
    address TEXT,
    phone VARCHAR(20),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- ROLE & USER
CREATE TABLE roles (
    id SMALLSERIAL PRIMARY KEY,
    name VARCHAR(30) UNIQUE NOT NULL -- admin, registrar, medic, head
);

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    facility_id UUID NOT NULL REFERENCES facilities (id),
    full_name VARCHAR(120) NOT NULL,
    email VARCHAR(150) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    last_login_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE user_roles (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role_id SMALLINT NOT NULL REFERENCES roles (id),
    PRIMARY KEY (user_id, role_id)
);

-- PASIEN
CREATE TABLE patients (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    facility_id UUID NOT NULL REFERENCES facilities (id),
    medical_record_no VARCHAR(20) NOT NULL,
    nik_encrypted BYTEA, -- NIK terenkripsi
    nik_hash VARCHAR(64), -- hash untuk pencarian/unik
    full_name VARCHAR(150) NOT NULL,
    birth_date DATE NOT NULL,
    gender gender_t NOT NULL,
    phone VARCHAR(20),
    address TEXT,
    guarantor guarantor_t NOT NULL DEFAULT 'umum',
    bpjs_no VARCHAR(20),
    consent_given_at TIMESTAMPTZ, -- persetujuan data (UU PDP)
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    UNIQUE (facility_id, medical_record_no),
    UNIQUE (facility_id, nik_hash)
);

CREATE INDEX idx_patients_name_trgm ON patients USING gin (full_name gin_trgm_ops);
CREATE INDEX idx_patients_phone ON patients (phone);

-- POLI & DOKTER
CREATE TABLE polyclinics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    facility_id UUID NOT NULL REFERENCES facilities (id),
    name VARCHAR(80) NOT NULL, -- Poli Umum, KIA, Gigi, dll
    queue_prefix CHAR(2) NOT NULL, -- A, B, C
    avg_service_min SMALLINT NOT NULL DEFAULT 10, -- default estimasi
    is_active BOOLEAN NOT NULL DEFAULT true,
    UNIQUE (facility_id, name)
);

CREATE TABLE doctors (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    facility_id UUID NOT NULL REFERENCES facilities (id),
    user_id UUID REFERENCES users (id),
    full_name VARCHAR(120) NOT NULL,
    specialty VARCHAR(80),
    str_no VARCHAR(40),
    is_active BOOLEAN NOT NULL DEFAULT true
);

CREATE TABLE schedules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    polyclinic_id UUID NOT NULL REFERENCES polyclinics(id),
    doctor_id UUID NOT NULL REFERENCES doctors(id),
    day_of_week SMALLINT NOT NULL CHECK (day_of_week BETWEEN 1 AND 7),
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    quota SMALLINT NOT NULL DEFAULT 30,
    CHECK (end_time > start_time)
);

-- PENDAFTARAN KUNJUNGAN
CREATE TABLE appointments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    facility_id UUID NOT NULL REFERENCES facilities(id),
    patient_id UUID NOT NULL REFERENCES patients (id),
    polyclinic_id UUID NOT NULL REFERENCES polyclinics(id),
    schedule_id UUID REFERENCES schedules(id),
    visit_date DATE NOT NULL,
    source appt_source_t NOT NULL DEFAULT 'online',
    priority priority_t NOT NULL DEFAULT 'normal',
    qr_token VARCHAR(64) UNIQUE NOT NULL,
    created_by UUID REFERENCES users(id), -- null jika pasien sendiri
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_appt_date_poli ON appointments (visit_date, polyclinic_id);

-- ANTREAN
CREATE TABLE queues (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    appointment_id UUID UNIQUE NOT NULL REFERENCES appointments (id) ON DELETE CASCADE,
    polyclinic_id UUID NOT NULL REFERENCES polyclinics(id),
    queue_date DATE NOT NULL,
    queue_number INT NOT NULL,
    queue_code VARCHAR(10) NOT NULL, -- A-012
    status queue_status_t NOT NULL DEFAULT 'waiting',
    checked_in_at TIMESTAMPTZ,
    called_at TIMESTAMPTZ,
    call_count SMALLINT NOT NULL DEFAULT 0,
    service_start TIMESTAMPTZ,
    service_end TIMESTAMPTZ,
    UNIQUE (polyclinic_id, queue_date, queue_number)
);

CREATE INDEX idx_queue_live ON queues (polyclinic_id, queue_date, status);

-- PRA-ANAMNESIS
CREATE TABLE pre_anamneses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    appointment_id UUID UNIQUE NOT NULL REFERENCES appointments (id) ON DELETE CASCADE,
    chief_complaint TEXT NOT NULL,
    symptoms JSONB, -- daftar gejala terstruktur
    duration_days SMALLINT,
    allergies TEXT,
    triage_flag VARCHAR(20), -- normal | attention | urgent
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- KUNJUNGAN / PEMERIKSAAN
CREATE TABLE visits (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    appointment_id UUID UNIQUE NOT NULL REFERENCES appointments (id),
    patient_id UUID NOT NULL REFERENCES patients (id),
    doctor_id UUID REFERENCES doctors(id),
    vitals JSONB, -- {bp_sys, bp_dia, hr, temp, spo2, weight, height}
    subjective TEXT,
    objective TEXT,
    assessment TEXT,
    plan TEXT,
    created_by UUID REFERENCES users (id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_visits_patient ON visits (patient_id, created_at DESC);

CREATE TABLE icd10_codes (
    code VARCHAR(10) PRIMARY KEY,
    description TEXT NOT NULL
);

CREATE TABLE diagnoses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    visit_id UUID NOT NULL REFERENCES visits (id) ON DELETE CASCADE,
    icd10_code VARCHAR(10) NOT NULL REFERENCES icd10_codes (code),
    is_primary BOOLEAN NOT NULL DEFAULT false
);

-- NOTIFIKASI
CREATE TABLE notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    patient_id UUID NOT NULL REFERENCES patients (id),
    queue_id UUID REFERENCES queues (id),
    channel notif_channel_t NOT NULL,
    message TEXT NOT NULL,
    status notif_status_t NOT NULL DEFAULT 'pending',
    sent_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- AUDIT LOG
CREATE TABLE audit_logs (
    id BIGSERIAL PRIMARY KEY,
    user_id UUID REFERENCES users(id),
    action VARCHAR(30) NOT NULL, -- read, create, update, delete, login
    entity VARCHAR(40) NOT NULL,
    entity_id UUID,
    ip_address INET,
    meta JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_audit_entity ON audit_logs (entity, entity_id, created_at DESC);
```

**Catatan desain database:**
* NIK disimpan terenkripsi (`nik_encrypted`) dan di-hash (`nik_hash`) untuk pencarian serta pencegahan duplikat tanpa membuka NIK asli.
* Nomor antrean dibuat di dalam transaksi dengan `SELECT ... FOR UPDATE` (atau sequence per poli per hari) agar tidak bentrok saat banyak pendaftar bersamaan.
* Estimasi waktu tunggu: jumlah antrean di depan × rata-rata durasi layanan (dihitung dari `service_end` - `service_start` hari-hari terakhir, fallback ke `avg_service_min`).
* Soft delete (`deleted_at`) pada pasien; data klinis tidak dihapus fisik.

## 11. Rancangan API (Ringkas)

| Method | Endpoint | Fungsi | Peran |
| :--- | :--- | :--- | :--- |
| `POST` | `/api/v1/auth/login` | Login staf | Publik |
| `POST` | `/api/v1/auth/otp/request`, `/verify` | OTP pasien | Publik |
| `GET` | `/api/v1/polyclinics` | Daftar poli | Publik |
| `GET` | `/api/v1/schedules?polyclinic_id=&date=` | Jadwal dan sisa kuota | Publik |
| `POST` | `/api/v1/patients` | Daftar pasien baru | Pasien/Registrar |
| `GET` | `/api/v1/patients?q=` | Cari pasien | Registrar/Medic |
| `GET` | `/api/v1/patients/{id}/visits` | Riwayat kunjungan | Medic |
| `POST` | `/api/v1/appointments` | Buat pendaftaran kunjungan | Pasien/Registrar |
| `POST` | `/api/v1/appointments/{id}/pre-anamnesis` | Isi pra-anamnesis | Pasien |
| `POST` | `/api/v1/queues/checkin` | Check-in via QR | Pasien/Registrar |
| `GET` | `/api/v1/queues?polyclinic_id=&date=` | Antrean aktif | Staf |
| `PATCH` | `/api/v1/queues/{id}/status` | Ubah status (panggil, layani, selesai) | Medic/Registrar |
| `GET` | `/api/v1/queues/{id}/estimate` | Estimasi tunggu | Pasien |
| `GET` | `/ws/queues/{polyclinic_id}` | Stream realtime antrean | Publik/Staf |
| `POST` | `/api/v1/visits` | Catat pemeriksaan | Medic |
| `GET` | `/api/v1/analytics/summary` | Dashboard analitik | Head/Admin |
| `GET` | `/api/v1/audit-logs` | Audit | Admin |

## 12. Rancangan Antarmuka (Halaman Utama)

| Area | Halaman |
| :--- | :--- |
| **Pasien** | Beranda, pilih poli dan jadwal, form pra-anamnesis, tiket antrean (QR + estimasi), riwayat kunjungan |
| **Petugas pendaftaran** | Cari/daftar pasien, buat kunjungan walk-in, panel antrean per poli, check-in |
| **Tenaga medis** | Daftar antrean poli, detail pasien dan pra-anamnesis, form pemeriksaan, input diagnosis ICD-10 |
| **Kepala/Admin** | Dashboard analitik, kelola poli/jadwal/pengguna, audit log |
| **Publik** | Layar TV antrean (tampilan besar nomor dipanggil) |

**Prinsip desain:** mobile-first untuk pasien, tampilan padat dan cepat untuk petugas, teks besar dan kontras tinggi, maksimal 3 langkah untuk mendaftar.

## 13. Keamanan dan Kepatuhan
* Autentikasi JWT berumur pendek + refresh token, OTP untuk pasien.
* RBAC di middleware Go; data klinis hanya untuk peran medis.
* Enkripsi TLS in-transit, enkripsi kolom NIK at-rest, backup terenkripsi.
* Audit log untuk setiap akses data pasien.
* Persetujuan pasien (consent) tersimpan; data minimum yang diperlukan saja.
* Rate limiting dan proteksi dari brute force pada login dan OTP.

## 14. Rencana Rilis

| Fase | Durasi | Cakupan |
| :--- | :--- | :--- |
| **Fase 1 (MVP)** | 4 minggu | F-01 s.d. F-05, F-08, F-10 (autentikasi, pasien, pendaftaran, antrean realtime, check-in, pencarian) |
| **Fase 2** | 3 minggu | F-06, F-07, F-09, F-11, F-12 (notifikasi, pra-anamnesis, pencatatan kunjungan, dashboard, audit) |
| **Fase 3** | 3 minggu | F-13, F-14, F-15 (offline PWA, SATUSEHAT, triase aturan) dan uji pilot |

## 15. Risiko dan Mitigasi

| Risiko | Mitigasi |
| :--- | :--- |
| Petugas menolak sistem baru | UI sederhana, pelatihan singkat, mode walk-in di loket |
| Internet tidak stabil | PWA offline-first, payload ringan |
| Kebocoran data pasien | Enkripsi, RBAC, audit log, pembatasan akses |
| Pasien lansia kesulitan memakai web | Pendaftaran dibantu petugas, tampilan sederhana, notifikasi WhatsApp |
| Biaya notifikasi | Batasi notifikasi pada momen kritis, gunakan template |

## 16. Pemetaan ke Kriteria Penilaian

| Kriteria | Pemenuhan dalam PRD |
| :--- | :--- |
| **Inovasi solusi** | Estimasi waktu tunggu dinamis, pra-anamnesis digital dengan triase awal, notifikasi WhatsApp, dashboard analitik |
| **Implementasi teknis** | Vue 3 PWA, Go REST + WebSocket, PostgreSQL ternormalisasi dengan enkripsi NIK, arsitektur modular dan siap multi-tenant |
| **Fungsionalitas dan desain** | Alur dua sisi (pasien dan staf) yang jelas, UI mobile-first, layar TV antrean, aksesibilitas lansia |
| **Dampak dan kelayakan** | Metrik terukur (waktu tunggu, waktu pendaftaran), rencana pilot di puskesmas, jalur integrasi SATUSEHAT, kepatuhan UU PDP |