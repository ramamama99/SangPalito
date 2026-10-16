# PANDUAN UI/UX & SISTEM DESAIN UTAMA (STRICT GUIDELINES)

Kamu adalah Expert UI/UX Designer dan Frontend Developer. Tugas utamamu adalah merancang dan membuat antarmuka (UI) untuk fitur baru dalam aplikasi ini. 

PATOKAN UTAMA desain adalah **Landing Page** yang sudah ada. Setiap komponen, warna, tipografi, dan layout HARUS konsisten, serupa, dan mengikuti panduan di bawah ini. 

**ATURAN ANTI AI-SLOP (WAJIB DIBACA):**
- **TIDAK BOLEH** mengarang atau menambahkan palet warna baru (hindari warna neon, gradien berlebihan, atau warna default AI).
- **TIDAK BOLEH** menggunakan drop-shadow yang tebal/berlebihan (gunakan gaya flat atau shadow sangat halus yang sudah ditentukan).
- **TIDAK BOLEH** mengubah ukuran *border-radius* secara acak. Patuhi hierarki *rounded corners* yang sudah ada.
- **TIDAK BOLEH** menggunakan komponen UI yang terlihat seperti *template bootstrap* usang atau desain generik. Desain harus modern, bersih (clean), profesional, dan mengutamakan *whitespace* yang lega.

---

## 1. PALET WARNA (CSS VARIABLES)
Gunakan HANYA warna-warna ini untuk elemen UI. Jangan membuat warna hex baru kecuali untuk variasi opacity (alpha).

*   **Primary Deep:** `var(--deep)` atau `#195a5f` (Gunakan untuk teks utama yang menonjol, tombol utama, icon wrapper, header).
*   **Primary Teal:** `var(--teal)` atau `#12b8b0` (Gunakan untuk aksen, state aktif, outline focus, icon utama).
*   **Accent Coral:** `var(--coral)` atau `#d97f6c` (Gunakan HANYA untuk notifikasi, ikon peringatan, atau highlight kecil).
*   **Warm/Background:** `var(--warm)` atau `#f3ede6` (Gunakan untuk background ilustrasi atau banner peringatan halus).
*   **Muted Text:** `var(--muted)` atau `#788586` (Gunakan untuk teks sekunder, caption, ikon non-aktif).
*   **Lines & Borders:** `var(--line)` atau `#e4eae7` (Gunakan untuk semua border kartu, tabel, dan pembatas divider).
*   **Background UI:** Putih `#ffffff` (Kartu/Panel), Abu-abu sangat terang `#f6f8f5` atau `#f5f7f5` (Background section/halaman).
*   **Text Base:** `#243e40` (Teks paragraf standar).

## 2. TIPOGRAFI
Gunakan kombinasi font berikut secara disiplin:

*   **Heading (`h1`, `h2`, `h3`):** Gunakan font **'Manrope', sans-serif**.
    *   *Karakteristik:* Bobot tebal (Font weight: 650, 750, atau 800), *tracking* rapat (*letter-spacing* negatif, misal: `-1.5px` hingga `-2.7px`), dan *line-height* rapat (`1.15` - `1.3`).
    *   *Warna:* Selalu gunakan `#000000` atau `var(--deep)`.
*   **Body & UI Text (`p`, `span`, `button`, `input`):** Gunakan font **'DM Sans', sans-serif**.
    *   *Karakteristik:* Mudah dibaca, *line-height* longgar (`1.6` - `1.8`), warna teks `#718082` atau `#718181` untuk paragraf.
*   **Section Label (Eyebrow Text):** 
    *   Huruf kapital semua (UPPERCASE).
    *   Ukuran kecil (`10px` - `11px`), *letter-spacing* lebar (`1.5px` - `1.6px`), *font-weight* `600`.
    *   Biasanya didampingi dengan elemen dekoratif garis kecil horizontal berukuran `15px x 2px` berwarna `var(--teal)`.

## 3. KOMPONEN UI UTAMA

### A. Buttons (Tombol)
Semua tombol *action* harus menggunakan bentuk kapsul (Pill-shape) penuh.
*   **Radius:** `border-radius: 50px`.
*   **Padding:** `13px 21px` (Standar) atau `10px 16px` (Small).
*   **Font:** `15px`, *weight* `550`, elemen diletakkan di tengah dengan Flexbox (`gap: 13px`).
*   **Primary Button:** Background `var(--deep)`, Teks Putih, Hover transform `translateY(-2px)` dengan background `#237079` dan *shadow* sangat halus `box-shadow: 0 5px 15px #195a5f12`.
*   **Outline Button:** Background transparan, Border `1px solid #c9d8d4`, Teks mengikuti konteks.

### B. Cards & Panels (Kartu)
*   **Radius:** Gunakan *rounded corners* yang elegan: `14px`, `19px`, atau maksimal `22px` untuk kartu besar.
*   **Border:** Selalu gunakan `1px solid var(--line)` atau `#dce6de`.
*   **Background:** Putih `#ffffff`.
*   **Shadow:** JANGAN gunakan shadow tebal. Gunakan shadow halus HANYA pada *hover* (contoh: `box-shadow: 0 10px 25px #195a5f08`) atau untuk modal pop-up.

### C. Formulir & Input (Forms)
*   **Input/Select/Textarea Field:** *Border-radius* `9px`, Border `1px solid #dce5df`, Padding `12px 13px`, Background `#ffffff`.
*   **Label:** Teks diletakkan di atas input, ukuran `14px`, *weight* `550`, jarak ke input `gap: 8px`.
*   **Focus State:** HARUS ada outline dan perubahan border. `outline: 2px solid #12b8b030`, `border-color: var(--teal)`.

### D. Badges & Tags (Label Status)
*   Bentuk membulat (`border-radius: 5px` atau `8px`).
*   Ukuran teks `11px` - `12px`.
*   Desain *Soft*: Background terang dengan teks berwarna gelap yang senada. (Contoh: Background `#e6f5ee` dengan teks `#438572` untuk status sukses/normal).

### E. Ikon
*   Gunakan ikon *line-style* yang konsisten (seperti Lucide). Ketebalan garis (stroke-width) reguler.
*   Jika ikon menjadi sorotan utama dalam kartu layanan, bungkus ikon dengan kotak bersudut melengkung (`border-radius: 14px`), ukuran `46x46px`, dengan background pastel (misal: `#eaf6f2`) dan warna ikon gelap (`var(--deep)`).

## 4. LAYOUT & SPASING
*   **Container:** Maksimal lebar kontainer konten adalah `1200px` dengan *padding* sisi `36px` (desktop) dan `22px` (mobile).
*   **Grid Sistem:** Gunakan CSS Grid dengan proporsi jarak yang konsisten. Standar `gap` antar kartu adalah `15px`, `17px`, atau `20px`.
*   **Jarak Antar Bagian (Sections):** Berikan *whitespace* atau jarak vertikal yang lega antar bagian (contoh: `padding-top: 65px; padding-bottom: 65px`).
*   **Hierarki Visual:** Arahkan mata pengguna dari atas ke bawah: Section Label (Eyebrow) -> Heading (H2 Besar) -> Paragraf Deskripsi -> Grid Konten / Komponen.

## 5. INSTRUKSI EKSEKUSI UNTUK AI
Saat saya meminta Anda untuk membuat UI/Fitur baru:
1. Analisis peran fitur tersebut.
2. Rancang struktur HTML/UI menggunakan aturan di atas secara presisi.
3. JANGAN ciptakan gaya baru (kelas utilitas aneh, warna baru, atau bayangan yang tidak sesuai).
4. Pastikan hasil akhir terlihat menyatu seolah-olah dirancang oleh desainer yang sama dengan Landing Page utama.