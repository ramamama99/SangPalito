<div align="center">

# 🏥 SangPalito (AntriSehat)
### *Pelita Baru Layanan Kesehatan: Ekosistem Cerdas Antrean Real-Time, Smart Wait-Time Engine & Rekam Medis Terpadu Puskesmas*

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://golang.org)
[![Vue Version](https://img.shields.io/badge/Vue.js-3.5-4FC08D?style=for-the-badge&logo=vue.js&logoColor=white)](https://vuejs.org)
[![TypeScript](https://img.shields.io/badge/TypeScript-5.x-3178C6?style=for-the-badge&logo=typescript&logoColor=white)](https://www.typescriptlang.org)
[![Tailwind CSS](https://img.shields.io/badge/Tailwind_CSS-v4-06B6D4?style=for-the-badge&logo=tailwindcss&logoColor=white)](https://tailwindcss.com)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-4169E1?style=for-the-badge&logo=postgresql&logoColor=white)](https://www.postgresql.org)
[![WebSocket](https://img.shields.io/badge/Realtime-WebSocket-010101?style=for-the-badge&logo=socketdotio&logoColor=white)](https://developer.mozilla.org/en-US/docs/Web/API/WebSockets_API)

<br />

<!-- PLACEHOLDER HERO BANNER / MOCKUP -->
```
+--------------------------------------------------------------------------------------------------------+
|                                                                                                        |
|                                [MASUKKAN BANNER / HERO MOCKUP DI SINI]                                 |
|                                       Contoh: docs/assets/hero-banner.png                              |
|                                                                                                        |
+--------------------------------------------------------------------------------------------------------+
```
*(Ganti placeholder di atas dengan screenshot mockup aplikasi: `![SangPalito Banner](path/to/hero-banner.png)`)*

<br />

[🌐 Live Demo](#-link-penting) • [🎬 Video Demo](#-link-penting) • [📊 Pitch Deck](#-link-penting) • [📑 Dokumentasi Sistem](#-dokumentasi--arsitektur) • [🚀 Panduan Instalasi](#-panduan-instalasi--menjalankan-lokal)

</div>

---

## 📌 Informasi Proyek

- **Nama Proyek:** SangPalito *(AntriSehat)*
  > *Filosofi: "Palito" (Bahasa Minangkabau) bermakna pelita atau lentera penerang. SangPalito hadir sebagai pelita yang menerangi dan menyederhanakan akses layanan kesehatan publik di fasilitas kesehatan primer.*
- **Tagline:** *"Lebih Sedikit Menunggu, Lebih Banyak Hidup Sehat — Pangkas Antrean Faskes dengan Estimasi Cerdas & Rekam Medis Terpadu."*
- **Nama / Tema Hackathon:** `[Masukkan Nama Hackathon, contoh: Hackathon Inovasi Layanan Kesehatan Digital 2026]`
- **Kategori:** *HealthTech / Public Service Digitalization / Smart Healthcare*

---

## 💡 Latar Belakang & Masalah

Fasilitas Kesehatan Tingkat Pertama (Puskesmas dan Klinik Pratama) merupakan garda terdepan layanan kesehatan masyarakat. Namun, faskes primer di Indonesia masih dihadapkan pada kendala operasional yang krusial:

```
           MASALAH UTAMA FASKES SAAT INI
┌─────────────────────────────────────────────────────────────┐
│ 1. Antrean Liar & Waktu Tunggu Tidak Pasti (2-4 Jam)        │
│    Pasien menumpuk di ruang tunggu -> Risiko infeksi silang │
├─────────────────────────────────────────────────────────────┤
│ 2. Estimasi Antrean Statis yang Sering Meleset              │
│    (Hanya mengalikan jumlah antrean x 10 menit tanpa tahu   │
│    kecepatan dokter hari itu & kompleksitas pasien)         │
├─────────────────────────────────────────────────────────────┤
│ 3. Beban Loket Manual & Input Berulang                      │
│    Pencatatan kertas, ketik NIK manual lambat, salah ketik  │
├─────────────────────────────────────────────────────────────┤
│ 4. Rekam Medis Tercecer & Belum Terstandar                  │
│    Riwayat keluhan, diagnosa ICD-10, & obat sulit dilacak   │
└─────────────────────────────────────────────────────────────┘
```

| Tantangan Konvensional | Dampak pada Pasien & Faskes | Solusi SangPalito |
| :--- | :--- | :--- |
| **Waktu Tunggu Buta** | Pasien datang pagi-pagi buta hanya untuk ambil kertas antrean | **Pendaftaran Online + E-Ticket QR + Real-Time Live Queue** |
| **Estimasi Waktu Meleset** | Ketidakpuasan pasien tinggi saat antrean lambat bergerak | **Smart Wait-Time Engine (EMA + Koefisien Klinis)** |
| **Verifikasi Loket Lambat** | Antrean mengular hanya untuk verifikasi berkas KTP | **Client-side OCR Scan KTP Otomatis (Tesseract.js)** |
| **Panggilan Manual Hilang** | Pasien tidak dengar panggilan saat keluar ruang tunggu | **Display TV + Voice Announcer + Notifikasi Real-time** |
| **Rekam Medis Terpisah** | Dokter kesulitan mengevaluasi riwayat kunjungan berkala | **EMR Terintegrasi (SOAP, Tanda Vital, ICD-10, E-Resep)** |

---

## ✨ Solusi: Apa itu SangPalito?

**SangPalito** adalah platform ekosistem faskes modern yang mendigitalkan seluruh siklus kunjungan pasien secara *end-to-end* — mulai dari pendaftaran mandiri di rumah, verifikasi kehadiran instan berbasis QR, estimasi giliran adaptif dengan kecerdasan buatan, hingga rekam medis digital komprehensif untuk dokter dan analitik operasional bagi pimpinan puskesmas.

```
       ALUR PENGGUNA TERINTEGRASI (END-TO-END FLOW)
       
[ Pasien di Rumah ]  ──(Daftar Online / OCR KTP)──> [ E-Ticket & Estimasi Adaptif ]
                                                                 │
                                                                 ▼
[ Display TV Poli ]  <──(WebSocket Realtime)──  [ Datang & Scan QR Check-in ]
  + Suara Panggil                                                │
                                                                 ▼
[ Dashboard Manajemen ] <──(Audit Log & ICD-10)── [ Dokter Periksa (SOAP + Resep) ]
```

---

## 🚀 Fitur Unggulan

### 1. 🧠 Smart Wait-Time Engine *(Algoritma EMA + Kompleksitas Klinis)*
Bukan sekadar perkalian statis! SangPalito menerapkan algoritma adaptif:
- **Exponential Moving Average (EMA, $\alpha = 0.30$):** Menyesuaikan ritme aktual dokter hari ini berdasarkan durasi riil pasien yang baru selesai (dengan batas aman *clamping* 3–45 menit).
- **Faktor Kompleksitas Klinis Pasien ($C_p$):** Mengkalkulasi koefisien khusus berdasarkan kelompok usia (lansia geriatri $+0.30$, balita $+0.20$), tingkat keparahan triase (*emergency* $+0.40$, *urgent* $+0.25$), dan banyaknya keluhan pra-anamnesis.
- Pasien tahu **jam panggilan yang sangat akurat**, sehingga tidak perlu berkerumun di ruang tunggu.

### 2. 🪪 Pendaftaran Cerdas & OCR Scan KTP Mandiri
- Pasien atau petugas cukup mengunggah/memfoto KTP; sistem mengekstrak NIK, Nama, dan Tanggal Lahir secara otomatis via **Tesseract.js OCR** di sisi klien (cepat, hemat bandwidth, privasi terjaga).
- Deteksi duplikasi pasien otomatis dengan pencocokan identitas tunggal.

### 3. 🎫 E-Ticket Interaktif & Scan QR Check-in Instan
- Pasien menerima tiket digital beranimasi dilengkapi **Dynamic QR Code**.
- Saat tiba di faskes, pasien cukup memindai QR di kiosk/kamera scanner petugas untuk mengubah status menjadi `checked_in` dalam hitungan detik.

### 4. 📺 Live Queue Tracking & TV Display Ruang Tunggu dengan Panggilan Suara
- **Live Tracking:** Pasien memantau pergerakan antrean langsung dari browser ponsel via WebSocket.
- **TV Display Mode:** Tampilan layar besar resolusi tinggi untuk ruang tunggu faskes, lengkap dengan **fitur text-to-speech audio announcer** berbahasa Indonesia (memanggil otomatis: *"Nomor Antrean A-005, silakan menuju Poli Umum"*).

### 5. 🩺 Dashboard Staf & Rekam Medis Elektronik (EMR) Terpadu
- **Multi-Role Access Control (RBAC):** Loket Pendaftaran, Tenaga Medis (Perawat/Dokter), Farmasi/Apotek, dan Kepala Puskesmas.
- **Pencatatan Klinis SOAP Lengkap:** Subjektif (keluhan & pra-anamnesis), Objektif (tanda vital: tensi, nadi, suhu, SpO2, BMI otomatis), Asesmen (pencarian kode standar **ICD-10**), dan Plan (resep obat digital terhubung ke katalog farmasi).

### 6. 📅 Manajemen Shift Dokter Dinamis *(Drag & Drop)*
- Antarmuka visual intuitif untuk mengatur jadwal praktik dokter, alokasi ruangan poli, dan batasan kuota harian pasien.

### 7. 🔒 Keamanan Medis & Kepatuhan Regulasi (UU PDP)
- Enkripsi data sensitif (NIK) pada basis data dengan algoritma **AES-256** (*at-rest encryption*).
- Password ter-hash menggunakan **Bcrypt** ($cost = 12$).
- Sesi JWT aman dengan perpanjangan *refresh window* dan **Audit Log** menyeluruh terhadap seluruh aktivitas akses data medis.

---

## 📸 Tangkapan Layar Aplikasi (Screenshots)

<div align="center">

### 📱 Antarmuka Pasien (Mobile-First)

| Pendaftaran & OCR KTP | E-Ticket & Estimasi Cerdas | Pantau Antrean Real-Time |
| :---: | :---: | :---: |
| *(Placeholder Screenshot)*<br>`docs/screenshots/register-ocr.png` | *(Placeholder Screenshot)*<br>`docs/screenshots/ticket-view.png` | *(Placeholder Screenshot)*<br>`docs/screenshots/live-queue.png` |

<br />

### 💻 Dashboard Faskes & Display Ruang Tunggu

| TV Ruang Tunggu + Panggilan Suara | Dashboard Dokter (Pemeriksaan SOAP) | Manajemen Shift Dokter |
| :---: | :---: | :---: |
| *(Placeholder Screenshot)*<br>`docs/screenshots/tv-display.png` | *(Placeholder Screenshot)*<br>`docs/screenshots/emr-soap.png` | *(Placeholder Screenshot)*<br>`docs/screenshots/doctor-shifts.png` |

</div>

> *(Catatan: Simpan gambar tangkapan layar ke dalam folder `docs/screenshots/` dan tautkan ke tabel di atas).*

---

## 🛠️ Tech Stack

SangPalito dibangun dengan arsitektur modern berorientasi performa tinggi, modularitas, dan skalabilitas tinggi:

```
┌────────────────────────────────────────────────────────────────────────┐
│                               TECH STACK                               │
├─────────────────┬──────────────────────────────────────────────────────┤
│ Frontend        │ Vue.js 3.5 (Composition API), Vite 8, TypeScript     │
│ Styling & UI    │ Tailwind CSS v4, Lucide Icons, Custom Design System  │
│ State & Routing │ Pinia Store, Vue Router 5                            │
│ Client Vision   │ Tesseract.js (KTP OCR), HTML5-QRCode, jsQR           │
├─────────────────┼──────────────────────────────────────────────────────┤
│ Backend API     │ Go 1.22 (Golang), Gin Web Framework                  │
│ Realtime Engine │ Gorilla WebSocket (Hub Pub/Sub Architecture)         │
│ Keamanan & Auth │ Golang-JWT v5, AES-256 Crypto, Bcrypt, Validator v10 │
├─────────────────┼──────────────────────────────────────────────────────┤
│ Database        │ PostgreSQL 16 (pgx/v5 connection pool, pgcrypto)     │
│ Cloud & DevOps  │ Docker, Docker Compose, Supabase                     │
└─────────────────┴──────────────────────────────────────────────────────┘
```

---

## 🏗️ Arsitektur Sistem

```
                      +-----------------------------+
                      |   Pasien & Display TV Faskes|
                      |   (PWA / Mobile / Kiosk)    |
                      +--------------+--------------+
                                     |
                          HTTPS / WSS (WebSocket)
                                     |
                                     v
                      +-----------------------------+
                      |       Gin API Router        |
                      +--------------+--------------+
                                     |
            +------------------------+------------------------+
            |                        |                        |
            v                        v                        v
+-----------------------+ +--------------------+ +-----------------------+
|  Queue & Wait-Time    | |  Medical Record    | |    WebSocket Hub      |
|  Service (EMA Engine) | |  Service (SOAP)    | |  (Broadcast Event)    |
+-----------+-----------+ +---------+----------+ +-----------+-----------+
            |                       |                        |
            +-----------------------+------------------------+
                                    |
                                    v
                      +-----------------------------+
                      |   PostgreSQL 16 Database    |
                      |   (Encrypted NIK, ICD-10)   |
                      +-----------------------------+
```

---

## ⚡ Panduan Instalasi & Menjalankan Lokal

Ikuti langkah-langkah berikut untuk menjalankan seluruh ekosistem SangPalito di lingkungan lokal:

### 📋 Prasyarat
- **Node.js** v20.x atau lebih baru & **pnpm** (atau npm)
- **Go** v1.22 atau lebih baru
- **PostgreSQL 16** (bisa menggunakan Docker Compose bawaan atau PostgreSQL lokal)
- **Git**

---

### 1. Kloning Repositori
```bash
git clone https://github.com/ramamama99/SangPalito.git
cd SangPalito
```

---

### 2. Menjalankan Backend (Go)

1. Masuk ke direktori backend:
   ```bash
   cd backend
   ```

2. Konfigurasi Environment (`.env`):
   Salin atau buat file `.env` di dalam folder `backend`:
   ```env
   APP_ENV=development
   APP_PORT=8080
   APP_SECRET_KEY=f4a8c2e1b7d5a9f3e6c0d8b2a4f6e1d3c5b7a9f2e4d6c8b0a1f3e5d7c9b1a3e5
   NIK_ENCRYPTION_KEY=7e2b9f4a1c8d5e0b3a6c9d2f5e8b1a4d7c0e3f6a9b2d5e8c1f4a7d0b3e6a9c2f
   BCRYPT_COST=12
   JWT_ACCESS_EXPIRE=15m
   JWT_REFRESH_EXPIRE=168h

   # Konfigurasi Database (Gunakan PostgreSQL lokal atau Cloud)
   DB_HOST=localhost
   DB_PORT=5432
   DB_NAME=antri_sehat
   DB_USER=postgres
   DB_PASSWORD=postgres
   DB_SSLMODE=disable
   DB_MAX_CONNS=20
   DB_MIN_CONNS=2
   ```

3. Jalankan Database via Docker (Opsional jika belum punya Postgres):
   ```bash
   docker-compose up -d db
   ```

4. Jalankan Server Backend Go:
   ```bash
   go mod download
   go run cmd/api/main.go
   ```
   *Backend akan aktif di `http://localhost:8080` (API & WebSocket Hub).*

---

### 3. Menjalankan Frontend (Vue 3 + Vite)

1. Buka terminal baru dan masuk ke direktori frontend:
   ```bash
   cd frontend
   ```

2. Instal dependensi:
   ```bash
   pnpm install
   # atau gunakan: npm install
   ```

3. Jalankan server pengembang:
   ```bash
   pnpm dev
   ```
   *Frontend akan berjalan di `http://localhost:5173`.*

---

## 🔑 Kredensial Akun Demo (Untuk Pengujian Juri)

Untuk memudahkan peninjauan juri dan penguji, sistem telah dilengkapi dengan data *seed* siap pakai:

| Role / Akses | Email Login | Password | Fungsi Akses |
| :--- | :--- | :--- | :--- |
| **Admin Sistem** | `admin@antrisehat.id` | `password123` | Manajemen Staf, Master Poli, Audit Log |
| **Dokter Pemeriksa** | `rama@antrisehat.id` | `password123` | Panggil Pasien, Rekam Medis SOAP, ICD-10, Resep |
| **Petugas Loket** | `loket@antrisehat.id` | `password123` | Loket Antrean, Pendaftaran Walk-in, Scan QR Kiosk |
| **Kepala Puskesmas** | `kepala@antrisehat.id` | `password123` | Dashboard Analitik Eksekutif, Beban Poli |

> **Akses Cepat TV Ruang Tunggu:** Buka menu TV di URL `http://localhost:5173/tv` untuk melihat layar antrean publik dengan audio panggil suara otomatis.

---

## 🧗‍♂️ Tantangan Selama Hackathon & Solusi

1. **Akurasi Estimasi Waktu Tunggu Tanpa Lonjakan Ekstrem:**
   - *Tantangan:* Ketika seorang dokter memeriksa pasien darurat yang memakan waktu lama atau sebaliknya menyelesaikan konsultasi kilat, rumus estimasi waktu tunggu konvensional akan melonjak kacau.
   - *Solusi:* Menerapkan Exponential Moving Average (EMA) dengan *smoothing factor* $\alpha=0.30$ dikombinasikan dengan batasan ambang (*clamping*) durasi 3–45 menit dan pengali koefisien klinis per pasien.

2. **OCR KTP yang Cepat Tanpa Membebani Server:**
   - *Tantangan:* Pengiriman foto KTP berukuran besar ke server berpotensi membuat *bottleneck* dan memicu risiko privasi data.
   - *Solusi:* Mengimplementasikan **Tesseract.js** langsung pada sisi browser (*client-side*). Gambar diproses di perangkat lokal pengguna, dan hanya teks terverifikasi yang dikirim ke backend.

3. **Sinkronisasi Multi-Layar Realtime Zero-Latency:**
   - *Tantangan:* Memastikan pemanggilan antrean dari ruang dokter langsung merefleksikan perubahan status di ponsel pasien dan layar TV ruang tunggu secara instan.
   - *Solusi:* Merancang arsitektur **Gorilla WebSocket Hub** dengan *room-based event broadcast* per poliklinik di backend Go, menghasilkan *delay* di bawah 200 ms.

---

## 🔮 Rencana Pengembangan (What's Next)

- [ ] **Integrasi SATUSEHAT Kemenkes RI:** Mengadopsi standar HL7 FHIR untuk pertukaran rekam medis elektronik nasional secara legal dan terstandar.
- [ ] **AI-Powered Triase NLP & Pra-Anamnesis:** Menggunakan model LLM ringan untuk menganalisis keluhan bahasa sehari-hari pasien dan merekomendasikan poli yang tepat secara otomatis.
- [ ] **WhatsApp Bot Gateway Dua Arah:** Pengiriman barcode tiket dan pengingat proaktif *"3 antrean lagi giliran Anda"* via WhatsApp resmi.
- [ ] **Mode PWA Offline-First:** Sinkronisasi *background* ketika internet faskes di daerah 3T mengalami pemadaman mendadak.

---

## 🔗 Link Penting

- 🌐 **Live Demo Web:** `[Masukkan URL Demo, contoh: https://sangpalito.vercel.app]`
- 🎬 **Video Demo Pitch (YouTube):** `[Masukkan Link Video YouTube]`
- 📊 **Pitch Deck Presentation:** `[Masukkan Link Canva / Google Slides]`
- 📑 **Dokumentasi Flowchart & Database:** [`SangPalito_Flowchart_dan_Database.pdf`](./SangPalito_Flowchart_dan_Database.pdf)

---

## 👥 Anggota Tim & Kontributor

| Nama | Peran | Tanggung Jawab Utama | Kontak / Profil |
| :--- | :--- | :--- | :--- |
| **[Nama Anggota 1]** | **Lead Developer & Backend** | Arsitektur Go API, WebSocket Engine, Database & Smart Wait-Time Algorithm | [![GitHub](https://img.shields.io/badge/GitHub-Profile-181717?style=flat-square&logo=github)](https://github.com/ramamama99) |
| **[Nama Anggota 2]** | **Frontend Engineer** | Vue 3 UI/UX, OCR Tesseract.js, E-Ticket & TV Display Announcer | [![GitHub](https://img.shields.io/badge/GitHub-Profile-181717?style=flat-square&logo=github)](https://github.com) |
| **[Nama Anggota 3]** | **UI/UX & Product Specialist** | Wireframing, Clinical Workflow Design, Dokumentasi PRD & Pitching | [![LinkedIn](https://img.shields.io/badge/LinkedIn-Profile-0A66C2?style=flat-square&logo=linkedin)](https://linkedin.com) |

---

<div align="center">

Dibuat dengan ❤️ dan dedikasi untuk memajukan transformasi kesehatan digital Indonesia 🇮🇩  
**SangPalito — Lebih Sedikit Menunggu, Lebih Banyak Hidup Sehat.**

</div>
