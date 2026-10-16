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

[🌐 Live Demo (Vercel)](https://sangpalito-77777.vercel.app/) • [⚡ Backend API (Railway)](https://sangpalito-production.up.railway.app/) • [🎬 Video Demo (YouTube)](https://youtu.be/2Qlyf4NXTF4) • [📖 Panduan Penggunaan](#-panduan-penggunaan-cara-pakai) • [🚀 Instalasi Lokal](#-panduan-instalasi--menjalankan-lokal)

</div>

---

## 📌 Informasi Proyek

- **Nama Proyek:** SangPalito *(AntriSehat)*
  > *Filosofi: "Palito" (Bahasa Minangkabau) bermakna pelita atau lentera penerang. SangPalito hadir sebagai pelita yang menerangi dan menyederhanakan akses layanan kesehatan publik di fasilitas kesehatan primer.*
- **Tagline:** *"Lebih Sedikit Menunggu, Lebih Banyak Hidup Sehat — Pangkas Antrean Faskes dengan Estimasi Cerdas & Rekam Medis Terpadu."*
- **Kategori:** *HealthTech / Public Service Digitalization / Smart Healthcare*
- **Akses Live Demo:**
  - 🌐 **Frontend Web (Vercel):** [https://sangpalito-77777.vercel.app/](https://sangpalito-77777.vercel.app/)
  - ⚡ **Backend API & WebSocket (Railway):** [https://sangpalito-production.up.railway.app/](https://sangpalito-production.up.railway.app/)
  - 🎬 **Video Presentasi & Demo (YouTube):** [https://youtu.be/2Qlyf4NXTF4](https://youtu.be/2Qlyf4NXTF4)

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
- Deteksi duplikasi pasien otomatis dengan pencocokan identitas tunggal terenkripsi.

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

## 📖 Panduan Penggunaan (Cara Pakai)

Berikut adalah panduan langkah demi langkah mencoba seluruh fitur SangPalito, baik dari sisi **Pasien** maupun **Staf Faskes/Dokter**:

### A. 📱 Alur Pasien (Pendaftaran Mandiri & Pemantauan Antrean)

1. **Akses Web Pasien:**
   - Kunjungi URL [https://sangpalito-77777.vercel.app/](https://sangpalito-77777.vercel.app/) melalui browser ponsel atau laptop.
2. **Pilih Layanan & Poliklinik:**
   - Klik tombol **"Ambil Antrean"** di halaman beranda.
   - Pilih poli yang diinginkan (contoh: *Poli Umum*, *Poli Gigi*, *Poli KIA / Anak*, atau *Poli Lansia*).
3. **Pengisian Data (Dukungan OCR KTP Cerdas):**
   - **Gunakan OCR KTP:** Unggah foto KTP Anda, dan sistem akan membaca NIK serta Nama secara otomatis tanpa perlu diketik manual.
   - Atau isi data manual jika tidak membawa KTP fisik.
   - Isi keluhan utama dan gejala (pra-anamnesis) agar dokter mendapatkan ringkasan awal kondisi Anda.
4. **Mendapatkan E-Ticket Digital:**
   - Pasien langsung mendapatkan E-Ticket interaktif yang berisi:
     - Nomor antrean (misal: `A-002`)
     - **Dynamic QR Code** untuk check-in kehadiran di lokasi
     - **Estimasi Waktu Tunggu & Jam Panggilan Cerdas** yang selalu diperbarui secara real-time.
5. **Check-in di Faskes:**
   - Saat tiba di puskesmas/klinik, arahkan QR Code pada tiket ke kamera scanner kiosk atau loket faskes.
   - Status antrean Anda otomatis berubah menjadi **Checked In (Siap Dipanggil)**.
6. **Pantau Antrean:**
   - Buka menu **"Cek Antrean Saya"** untuk melihat sisa antrean di depan Anda, atau perhatikan layar TV antrean di ruang tunggu.

---

### B. 💻 Alur Staf & Dokter (Pelayanan Faskes)

1. **Login ke Sistem Staf:**
   - Buka halaman login di [https://sangpalito-77777.vercel.app/login](https://sangpalito-77777.vercel.app/login)
   - Masukkan kredensial akun:
     - **Email:** `admin@antrisehat.id`
     - **Password:** `password123`
2. **Dashboard Staf & Dokter (`/staff`):**
   - **Petugas Loket:** Dapat memverifikasi pasien *walk-in*, mengelola antrean loket, atau menggunakan menu scanner QR (`/scan`) untuk check-in pasien.
   - **Dokter Pemeriksa:**
     - Pilih poli aktif Anda, lihat antrean yang berstatus *Checked In*.
     - Klik tombol **"Panggil Pasien"**. Sistem akan mengubah status menjadi `called` dan secara bersamaan memicu **suara panggil otomatis** pada TV Display ruang tunggu.
     - Klik **"Mulai Layanan"** saat pasien masuk ke ruang periksa.
3. **Pengisian Rekam Medis Elektronik (SOAP & ICD-10):**
   - **Subjektif:** Tinjau keluhan awal yang sudah diisi pasien dari rumah.
   - **Objektif:** Masukkan tanda vital (Tensi Darah, Nadi, Suhu Tubuh, Laju Pernapasan).
   - **Asesmen:** Cari dan tetapkan diagnosis penyakit berdasarkan standar internasional **ICD-10** (contoh: *J00 - Nasopharyngitis Akut*).
   - **Plan:** Buat resep obat digital yang terintegrasi langsung dengan database katalog farmasi.
4. **Selesai Layanan & Pembaruan EMA:**
   - Klik tombol **"Selesai Pelayanan"**.
   - Sistem secara otomatis menghitung durasi konsultasi aktual dan menggunakannya untuk memperbarui formula *Exponential Moving Average (EMA)*, membuat estimasi antrean pasien berikutnya semakin akurat!

---

### C. 📺 Alur TV Display Ruang Tunggu (Voice Announcer)

1. Buka URL [https://sangpalito-77777.vercel.app/tv](https://sangpalito-77777.vercel.app/tv) pada monitor/layar TV ruang tunggu puskesmas.
2. Layar akan menampilkan daftar antrean aktif per poliklinik secara *full-screen*.
3. Setiap kali dokter menekan tombol panggil di ruang periksa, layar TV akan otomatis meng-highlight nomor antrean tersebut dan membunyikan **audio voice synthesizer bahasa Indonesia**:
   > *"Nomor antrean A-002, silakan menuju Ruang Poli Umum."*

---

## 🔑 Kredensial Akun Demo (Untuk Pengujian Juri)

Untuk memudahkan peninjauan juri dan penguji, sistem telah dilengkapi dengan data *seed* siap pakai:

| Role / Akses | Email Login | Password | Fungsi Akses |
| :--- | :--- | :--- | :--- |
| **Admin Sistem / Staf Medis** | `admin@antrisehat.id` | `password123` | Akses Dashboard Lengkap, Panggil Pasien, Rekam Medis SOAP & ICD-10, Manajemen Shift, Audit Log |

> **Akses Cepat TV Ruang Tunggu:** Kunjungi [https://sangpalito-77777.vercel.app/tv](https://sangpalito-77777.vercel.app/tv) untuk mencoba langsung antarmuka monitor ruang tunggu berfitur audio panggil real-time.

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
│ Frontend Host   │ Vercel (https://sangpalito-77777.vercel.app/)        │
├─────────────────┼──────────────────────────────────────────────────────┤
│ Backend API     │ Go 1.22 (Golang), Gin Web Framework                  │
│ Realtime Engine │ Gorilla WebSocket (Hub Pub/Sub Architecture)         │
│ Keamanan & Auth │ Golang-JWT v5, AES-256 Crypto, Bcrypt, Validator v10 │
│ Backend Host    │ Railway (https://sangpalito-production.up.railway.app)│
├─────────────────┼──────────────────────────────────────────────────────┤
│ Database        │ PostgreSQL 16 (pgx/v5 connection pool, pgcrypto)     │
│ Cloud & DevOps  │ Docker, Docker Compose, Supabase Cloud Database      │
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

Jika Anda ingin menjalankan seluruh ekosistem SangPalito di lingkungan lokal:

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
   PORT=5000
    DATABASE_URL=postgresql://user:password@localhost:5432/dbname
    SUPABASE_URL=https://your-project.supabase.co
    SUPABASE_KEY=your_public_anon_key_here
    OPENAI_API_KEY=your_openai_api_key_here
   (ini contoh)
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

---

## 🔮 Rencana Pengembangan (What's Next)

- [ ] **Integrasi SATUSEHAT Kemenkes RI:** Mengadopsi standar HL7 FHIR untuk pertukaran rekam medis elektronik nasional secara legal dan terstandar.
- [ ] **AI-Powered Triase NLP & Pra-Anamnesis:** Menggunakan model LLM ringan untuk menganalisis keluhan bahasa sehari-hari pasien dan merekomendasikan poli yang tepat secara otomatis.
- [ ] **WhatsApp Bot Gateway Dua Arah:** Pengiriman barcode tiket dan pengingat proaktif *"3 antrean lagi giliran Anda"* via WhatsApp resmi.
- [ ] **Mode PWA Offline-First:** Sinkronisasi *background* ketika internet faskes di daerah 3T mengalami pemadaman mendadak.

---

## 🔗 Link Penting

- 🌐 **Live Web Application (Vercel):** [https://sangpalito-77777.vercel.app/](https://sangpalito-77777.vercel.app/)
- ⚡ **Backend API & WebSocket Server (Railway):** [https://sangpalito-production.up.railway.app/](https://sangpalito-production.up.railway.app/)
- 🎬 **Video Presentasi & Demo (YouTube):** [https://youtu.be/2Qlyf4NXTF4](https://youtu.be/2Qlyf4NXTF4)

---

<div align="center">

Dibuat dengan ❤️ dan dedikasi untuk memajukan transformasi kesehatan digital Indonesia 🇮🇩  
**SangPalito — Lebih Sedikit Menunggu, Lebih Banyak Hidup Sehat.**

</div>
