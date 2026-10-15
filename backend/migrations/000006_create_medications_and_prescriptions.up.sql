-- 000006_create_medications_and_prescriptions.up.sql

-- 1. KATALOG OBAT (MEDICATIONS) LENGKAP DENGAN CARA PENGGUNAAN
CREATE TABLE IF NOT EXISTS medications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(150) NOT NULL,
    generic_name VARCHAR(150) NOT NULL,
    category VARCHAR(100) NOT NULL,
    dosage_form VARCHAR(50) NOT NULL,       -- Tablet, Tablet Kunyah, Kapsul, Sirup, Salep, Tetes Mata, Tetes Telinga, dll.
    route VARCHAR(50) NOT NULL,             -- oral_minum, oral_kunyah, topikal_oles, tetes_mata, tetes_telinga, kumur, inhalasi
    default_frequency VARCHAR(100) NOT NULL, -- misal: 3 x 1 tablet sehari, 2 x 1 sendok takar (5 ml), dll.
    default_timing VARCHAR(100) NOT NULL,    -- misal: Sesudah makan, Sebelum makan (perut kosong), Sebelum tidur
    usage_instruction TEXT NOT NULL,         -- Petunjuk detail: diminum dengan air, dikunyah halus, dioles tipis merata
    indication TEXT NOT NULL,                -- Digunakan untuk mengobati penyakit apa
    stock INT NOT NULL DEFAULT 100,
    unit VARCHAR(50) NOT NULL,               -- Tablet, Botol, Tube, Kapsul, Strip
    is_prescription_required BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_medications_name ON medications (name);
CREATE INDEX IF NOT EXISTS idx_medications_category ON medications (category);
CREATE INDEX IF NOT EXISTS idx_medications_route ON medications (route);

-- 2. RESEP ELEKTRONIK (PRESCRIPTIONS)
CREATE TABLE IF NOT EXISTS prescriptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    prescription_no VARCHAR(50) UNIQUE NOT NULL, -- e.g. RX-20261015-0001
    appointment_id UUID REFERENCES appointments (id) ON DELETE SET NULL,
    visit_id UUID REFERENCES visits (id) ON DELETE CASCADE,
    patient_id UUID NOT NULL REFERENCES patients (id) ON DELETE CASCADE,
    doctor_id UUID REFERENCES doctors (id) ON DELETE SET NULL,
    doctor_name VARCHAR(150),
    status VARCHAR(30) NOT NULL DEFAULT 'pending_dispense', -- pending_dispense, dispensed, cancelled
    notes TEXT,                                              -- Catatan dari dokter untuk apoteker / farmasi
    dispensed_at TIMESTAMPTZ,
    dispensed_by UUID REFERENCES users (id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_prescriptions_patient ON prescriptions (patient_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_prescriptions_visit ON prescriptions (visit_id);
CREATE INDEX IF NOT EXISTS idx_prescriptions_no ON prescriptions (prescription_no);
CREATE INDEX IF NOT EXISTS idx_prescriptions_status ON prescriptions (status);

-- 3. ITEM RESEP OBAT (PRESCRIPTION ITEMS)
CREATE TABLE IF NOT EXISTS prescription_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    prescription_id UUID NOT NULL REFERENCES prescriptions (id) ON DELETE CASCADE,
    medication_id UUID REFERENCES medications (id) ON DELETE SET NULL,
    medication_name VARCHAR(150) NOT NULL,
    dosage_form VARCHAR(50) NOT NULL,
    quantity INT NOT NULL DEFAULT 1,
    unit VARCHAR(50) NOT NULL,
    route VARCHAR(50) NOT NULL,              -- oral_minum, oral_kunyah, topikal_oles, tetes_mata, tetes_telinga, kumur, inhalasi
    frequency VARCHAR(100) NOT NULL,         -- misal: 3 x 1 tablet sehari
    timing VARCHAR(100) NOT NULL,            -- Sesudah makan, Sebelum makan, Sebelum tidur
    instruction TEXT NOT NULL,               -- Petunjuk spesifik: Wajib dikunyah hingga halus, dioles tipis 2x sehari, dll.
    notes TEXT,                              -- Catatan khusus (contoh: "Habiskan", "Bila demam saja")
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_prescription_items_prescription ON prescription_items (prescription_id);

-- 4. TAMBAHKAN KOLOM PERENCANAAN OBAT DI TABEL VISITS JIKA BELUM ADA
ALTER TABLE visits ADD COLUMN IF NOT EXISTS medication_plan TEXT;
ALTER TABLE visits ADD COLUMN IF NOT EXISTS prescription_no VARCHAR(50);

-- 5. SEED DATA KATALOG OBAT LENGKAP DENGAN CARA PENGGUNAAN BERVARIASI
INSERT INTO medications (code, name, generic_name, category, dosage_form, route, default_frequency, default_timing, usage_instruction, indication, stock, unit, is_prescription_required)
VALUES
-- A. OBAT DIMINUM (ORAL TABLET / KAPSUL)
('MED-PCT-500', 'Paracetamol 500 mg', 'Paracetamol', 'Analgesik & Antipiretik', 'Tablet', 'oral_minum', '3 x 1 tablet sehari', 'Sesudah makan', 'Diminum utuh dengan 1 gelas air putih hangat setelah makan bila demam/nyeri.', 'Menurunkan demam, meredakan sakit kepala, sakit gigi, dan nyeri ringan.', 500, 'Tablet', false),
('MED-AMX-500', 'Amoxicillin 500 mg', 'Amoxicillin Trihydrate', 'Antibiotik', 'Kapsul', 'oral_minum', '3 x 1 kapsul sehari (tiap 8 jam)', 'Sesudah makan', 'Diminum teratur tiap 8 jam sesudah makan. Wajib dihabiskan walau gejala sudah membaik agar kuman tidak resisten.', 'Mengobati infeksi bakteri saluran napas, kulit, dan infeksi telinga/saluran kemih.', 300, 'Kapsul', true),
('MED-CTZ-010', 'Cetirizine 10 mg', 'Cetirizine HCl', 'Antihistamin & Anti-alergi', 'Tablet', 'oral_minum', '1 x 1 tablet sehari', 'Sebelum tidur malam', 'Diminum malam hari sebelum tidur karena dapat menyebabkan kantuk ringan.', 'Meredakan gejala alergi, hidung tersumbat, bersin-bersin, gatal biduran.', 250, 'Tablet', true),
('MED-CPT-025', 'Captopril 25 mg', 'Captopril', 'Antihipertensi (ACE Inhibitor)', 'Tablet', 'oral_minum', '2 x 1/2 atau 1 tablet sehari', 'Sebelum makan (1 jam sebelum)', 'Diminum dalam kondisi perut kosong, 1 jam sebelum makan atau 2 jam sesudah makan untuk penyerapan optimal.', 'Menurunkan tekanan darah tinggi (hipertensi) dan gagal jantung.', 200, 'Tablet', true),
('MED-MET-500', 'Metformin 500 mg', 'Metformin HCl', 'Antidiabetes Oral', 'Tablet', 'oral_minum', '2 x 1 tablet sehari', 'Bersama atau sesudah makan', 'Diminum bersamaan dengan suapan pertama makanan untuk mengurangi efek mual atau sakit perut.', 'Mengontrol kadar gula darah pada pasien Diabetes Melitus Tipe 2.', 300, 'Tablet', true),

-- B. OBAT DIKUNYAH (ORAL CHEWABLE)
('MED-ANT-KNY', 'Antasida DOEN Kunyah', 'Aluminium Hidroksida + Magnesium Hidroksida', 'Antasida & Antiulkus', 'Tablet Kunyah', 'oral_kunyah', '3 x 1-2 tablet kunyah sehari', '1 jam sebelum makan atau saat maag kambuh', 'Wajib dikunyah hingga benar-benar hancur dan lumat di mulut sebelum ditelan, lalu bilas dengan sedikit air.', 'Meredakan mual, kembung, nyeri ulu hati, dan asam lambung berlebih / maag akut.', 400, 'Tablet', false),
('MED-VTC-KNY', 'Vitamin C 500 mg Kunyah', 'Asam Askorbat', 'Suplemen & Vitamin', 'Tablet Kunyah', 'oral_kunyah', '1 x 1 tablet kunyah sehari', 'Sesudah makan', 'Dikunyah perlahan di dalam mulut sampai habis dan larut.', 'Meningkatkan daya tahan tubuh dan mempercepat masa pemulihan setelah sakit.', 600, 'Tablet', false),

-- C. OBAT DIOLESKAN (TOPIKAL SALEP / KRIM / GEL)
('MED-HDC-001', 'Hidrokortison Krim 1%', 'Hydrocortisone', 'Kortikosteroid Topikal', 'Krim / Salep', 'topikal_oles', '2-3 x sehari', 'Sesudah membersihkan area kulit', 'Bersihkan dan keringkan area kulit yang gatal/meradang, lalu oleskan krim secara tipis dan merata. Jangan gunakan pada luka terbuka.', 'Meredakan radang, eksim, dermatitis, dan gatal kemerahan pada kulit.', 150, 'Tube', true),
('MED-GNT-SLP', 'Gentamicin Salep Kulit 0.1%', 'Gentamicin Sulfate', 'Antibiotik Topikal', 'Salep', 'topikal_oles', '2-3 x sehari', 'Sesudah luka dibersihkan', 'Oleskan tipis pada area luka/infeksi bakteri setelah dibersihkan dengan cairan antiseptik atau NaCl.', 'Mengobati infeksi bakteri pada kulit seperti impetigo, folikulitis, atau luka bernanah.', 120, 'Tube', true),
('MED-DIC-GEL', 'Natrium Diklofenak Emulgel 1%', 'Diclofenac Sodium', 'Antiinflamasi & Analgesik Topikal', 'Gel', 'topikal_oles', '3-4 x sehari', 'Saat nyeri sendi/otot terasa', 'Oleskan gel secukupnya pada area sendi atau otot yang nyeri, pijat perlahan sampai gel terserap sempurna.', 'Meredakan nyeri otot, pegal linu, memar, terkilir, dan osteoartritis sendi.', 180, 'Tube', false),

-- D. OBAT DITETESKAN (TETES MATA & TELINGA)
('MED-CXT-TTS', 'Cendo Xitrol Tetes Mata 5 ml', 'Dexamethasone + Neomycin + Polymyxin B', 'Antibiotik & Antiinflamasi Mata', 'Tetes Mata', 'tetes_mata', '3 x 1-2 tetes pada mata yang sakit', 'Kapan saja sesuai jadwal', 'Tengadahkan kepala, tarik kelopak mata bawah perlahan, teteskan 1-2 tetes tanpa menyentuhkan ujung botol ke mata, tutup mata 1-2 menit.', 'Mengobati infeksi mata bagian luar yang disertai radang/pembengkakan (konjungtivitis bakteri).', 100, 'Botol', true),
('MED-OTP-TTS', 'Otopain Tetes Telinga 8 ml', 'Polymyxin B + Neomycin + Fludrocortisone + Lidocaine', 'Antibiotik & Analgesik Telinga', 'Tetes Telinga', 'tetes_telinga', '3 x 2-4 tetes pada liang telinga', 'Pagi, siang, dan malam', 'Miringkan kepala ke samping, teteskan ke lubang telinga, pertahankan posisi miring selama 2-3 menit agar obat masuk sempurna.', 'Mengobati otitis eksterna (infeksi telinga luar) dan nyeri telinga berair.', 80, 'Botol', true),

-- E. OBAT SIRUP & LARUTAN (ANAK / DEWASA)
('MED-PCT-SYR', 'Paracetamol Sirup 120 mg/5 ml', 'Paracetamol', 'Analgesik & Antipiretik Anak', 'Sirup', 'oral_minum', '3-4 x sehari 1 sendok takar (5 ml)', 'Sesudah makan', 'Kocok botol terlebih dahulu sebelum dituangkan. Gunakan sendok takar obat resmi 5 ml yang disediakan.', 'Menurunkan demam dan meredakan nyeri pada bayi dan anak-anak.', 200, 'Botol', false),
('MED-OBH-SYR', 'OBH Batuk & Flu Sirup 100 ml', 'Succus Liquiritae + Paracetamol + Ammonium Chloride + CTM', 'Ekspektoran & Antitusif', 'Sirup', 'oral_minum', '3 x sehari 1 sendok makan (15 ml)', 'Sesudah makan', 'Kocok botol dahulu. Minum sesudah makan. Dapat menyebabkan rasa kantuk.', 'Meredakan batuk berdahak, pilek, bersin, dan sakit tenggorokan.', 150, 'Botol', false),

-- F. OBAT DIKUMUR (MOUTHWASH / GARGLE)
('MED-BTD-KMR', 'Povidone Iodine Gargle 1% 100 ml', 'Povidone Iodine', 'Antiseptik Rongga Mulut', 'Larutan Kumur', 'kumur', '3 x sehari 15 ml', 'Sesudah gosok gigi', 'Tuangkan 15 ml ke cangkir, kumur-kumur dan tengadahkan kepala untuk kumur di tenggorokan (gargle) selama 30-60 detik, lalu buang. JANGAN DITELAN.', 'Mengatasi sakit tenggorokan, sariawan, gusi bengkak, dan bau mulut.', 120, 'Botol', false)
ON CONFLICT (code) DO NOTHING;
