-- 000004_seed_demo_data.up.sql
-- Seed data realistis untuk antrean puskesmas, staf, poliklinik, pasien, dan rekam medis

-- 1. FASILITAS KESEHATAN (FASKES)
INSERT INTO facilities (id, name, type, code, address, phone)
VALUES (
    'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11',
    'Puskesmas Kebon Jeruk',
    'puskesmas',
    'PKM-KBJ-01',
    'Jl. Raya Kebon Jeruk No. 25, Jakarta Barat',
    '021-5301234'
)
ON CONFLICT (id) DO UPDATE SET 
    name = EXCLUDED.name,
    address = EXCLUDED.address,
    phone = EXCLUDED.phone;

-- 2. ROLES
INSERT INTO roles (id, name) VALUES
(1, 'admin'),
(2, 'registrar'),
(3, 'medic'),
(4, 'head')
ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name;

-- 3. AKUN PENGGUNA (USERS) - Password default: password123 ($2a$12$CTqJY92PEYtNokRy84lYfey60f4rlLlr009A2ksRUtE4cRj/2.S4K)
INSERT INTO users (id, facility_id, full_name, email, password_hash, is_active) VALUES
('b1c2d3e4-f5a6-4b7c-8d9e-0f1a2b3c4d5e', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Administrator Sistem', 'admin@antrisehat.id', '$2a$12$CTqJY92PEYtNokRy84lYfey60f4rlLlr009A2ksRUtE4cRj/2.S4K', true),
('b2c3d4e5-f6a7-4b8c-9d0e-1f2a3b4c5d6e', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'dr. Rama Palito', 'rama@antrisehat.id', '$2a$12$CTqJY92PEYtNokRy84lYfey60f4rlLlr009A2ksRUtE4cRj/2.S4K', true),
('b3c4d5e6-f7a8-4b9c-0d1e-2f3a4b5c6d7e', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Petugas Loket 1', 'loket@antrisehat.id', '$2a$12$CTqJY92PEYtNokRy84lYfey60f4rlLlr009A2ksRUtE4cRj/2.S4K', true),
('b4c5d6e7-f8a9-4b0c-1d2e-3f4a5b6c7d8e', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'dr. Hendra (Kepala Puskesmas)', 'kepala@antrisehat.id', '$2a$12$CTqJY92PEYtNokRy84lYfey60f4rlLlr009A2ksRUtE4cRj/2.S4K', true)
ON CONFLICT (email) DO UPDATE SET full_name = EXCLUDED.full_name, password_hash = EXCLUDED.password_hash;

-- 4. ROLE MAPPINGS
INSERT INTO user_roles (user_id, role_id) VALUES
('b1c2d3e4-f5a6-4b7c-8d9e-0f1a2b3c4d5e', 1),
('b2c3d4e5-f6a7-4b8c-9d0e-1f2a3b4c5d6e', 3),
('b3c4d5e6-f7a8-4b9c-0d1e-2f3a4b5c6d7e', 2),
('b4c5d6e7-f8a9-4b0c-1d2e-3f4a5b6c7d8e', 4)
ON CONFLICT DO NOTHING;

-- 5. POLIKLINIK
INSERT INTO polyclinics (id, facility_id, name, queue_prefix, avg_service_min, is_active) VALUES
('11111111-1111-1111-1111-111111111111', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Poli Umum', 'A', 10, true),
('22222222-2222-2222-2222-222222222222', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Poli Gigi', 'B', 15, true),
('33333333-3333-3333-3333-333333333333', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Poli KIA / Anak', 'C', 12, true),
('44444444-4444-4444-4444-444444444444', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Poli Lansia / Geriatri', 'D', 15, true)
ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name, queue_prefix = EXCLUDED.queue_prefix;

-- 6. DOKTER
INSERT INTO doctors (id, facility_id, user_id, full_name, specialty, str_no, is_active) VALUES
('d1111111-1111-1111-1111-111111111111', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'b2c3d4e5-f6a7-4b8c-9d0e-1f2a3b4c5d6e', 'dr. Rama Palito', 'Dokter Umum', 'STR-3201-9921', true),
('d2222222-2222-2222-2222-222222222222', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', NULL, 'dr. Maya Putri', 'Dokter Umum', 'STR-3201-8842', true),
('d3333333-3333-3333-3333-333333333333', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', NULL, 'drg. Rani Prameswari', 'Dokter Gigi', 'STR-3201-7715', true),
('d4444444-4444-4444-4444-444444444444', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', NULL, 'dr. Siti Rahmawati, Sp.A', 'Dokter Anak', 'STR-3201-6631', true)
ON CONFLICT (id) DO UPDATE SET full_name = EXCLUDED.full_name, specialty = EXCLUDED.specialty;

-- 7. JADWAL DOKTER
INSERT INTO schedules (id, polyclinic_id, doctor_id, day_of_week, start_time, end_time, quota) VALUES
('eb111111-1111-1111-1111-111111111111', '11111111-1111-1111-1111-111111111111', 'd1111111-1111-1111-1111-111111111111', 1, '08:00', '12:00', 30),
('eb222222-2222-2222-2222-222222222222', '11111111-1111-1111-1111-111111111111', 'd2222222-2222-2222-2222-222222222222', 1, '12:00', '15:00', 25),
('eb333333-3333-3333-3333-333333333333', '22222222-2222-2222-2222-222222222222', 'd3333333-3333-3333-3333-333333333333', 1, '08:30', '13:00', 20),
('eb444444-4444-4444-4444-444444444444', '33333333-3333-3333-3333-333333333333', 'd4444444-4444-4444-4444-444444444444', 1, '08:00', '12:00', 25)
ON CONFLICT (id) DO NOTHING;

-- 8. PASIEN & WALI
INSERT INTO patients (id, facility_id, medical_record_no, full_name, birth_date, gender, phone, address, guarantor, bpjs_no) VALUES
('bb111111-1111-1111-1111-111111111111', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'RM-2026-0001', 'Siti Aminah', '1990-01-12', 'female', '081234567890', 'Jl. Sudirman No. 12', 'bpjs', '0001234567890'),
('bb222222-2222-2222-2222-222222222222', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'RM-2026-0002', 'Budi Santoso', '1961-05-15', 'male', '081298765432', 'Jl. Merdeka No. 10', 'umum', NULL),
('bb333333-3333-3333-3333-333333333333', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'RM-2026-0003', 'Dewi Lestari', '1995-09-20', 'female', '081345678901', 'Jl. Melati No. 4', 'bpjs', '0002345678901'),
('bb444444-4444-4444-4444-444444444444', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'RM-2026-0004', 'Ahmad Fauzi', '2022-03-10', 'male', '081456789012', 'Jl. Mawar No. 7', 'umum', NULL),
('bb555555-5555-5555-5555-555555555555', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'RM-2026-0005', 'Hj. Rohayah', '1955-11-28', 'female', '081567890123', 'Jl. Anggrek No. 15', 'bpjs', '0003456789012')
ON CONFLICT (id) DO UPDATE SET full_name = EXCLUDED.full_name, phone = EXCLUDED.phone;

-- Data wali untuk pasien anak Ahmad Fauzi
UPDATE patients 
SET guardian_name = 'Hendri Fauzi (Ayah)', 
    guardian_relation = 'Ayah Kandung', 
    guardian_phone = '081456789012' 
WHERE id = 'bb444444-4444-4444-4444-444444444444';

-- 9. ANTREAN HARI INI (CURRENT DATE)
-- Appointment 1: Siti Aminah (Sedang diperiksa / in_service)
INSERT INTO appointments (id, facility_id, patient_id, polyclinic_id, schedule_id, visit_date, source, priority, appointment_token, qr_token, created_at)
VALUES (
    'a1111111-1111-1111-1111-111111111111',
    'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11',
    'bb111111-1111-1111-1111-111111111111',
    '11111111-1111-1111-1111-111111111111',
    'eb111111-1111-1111-1111-111111111111',
    CURRENT_DATE,
    'online',
    'normal',
    'f1111111-1111-1111-1111-111111111111',
    'qr-siti-aminah-001',
    CURRENT_DATE + TIME '07:30'
) ON CONFLICT (id) DO NOTHING;

INSERT INTO pre_anamneses (id, appointment_id, chief_complaint, symptoms, duration_days, allergies, triage_flag)
VALUES (
    'cb111111-1111-1111-1111-111111111111',
    'a1111111-1111-1111-1111-111111111111',
    'Demam dan sakit kepala sejak 2 hari yang lalu',
    '["Demam", "Sakit Kepala", "Batuk"]'::jsonb,
    2,
    'Tidak ada',
    'normal'
) ON CONFLICT (id) DO NOTHING;

INSERT INTO queues (id, appointment_id, polyclinic_id, queue_date, queue_number, queue_code, status, checked_in_at, called_at, call_count, service_start)
VALUES (
    'fb111111-1111-1111-1111-111111111111',
    'a1111111-1111-1111-1111-111111111111',
    '11111111-1111-1111-1111-111111111111',
    CURRENT_DATE,
    1,
    'A-001',
    'in_service',
    CURRENT_DATE + TIME '08:00',
    CURRENT_DATE + TIME '08:10',
    1,
    CURRENT_DATE + TIME '08:12'
) ON CONFLICT (id) DO NOTHING;

-- Appointment 2: Budi Santoso (Dipanggil / called)
INSERT INTO appointments (id, facility_id, patient_id, polyclinic_id, schedule_id, visit_date, source, priority, appointment_token, qr_token, created_at)
VALUES (
    'a2222222-2222-2222-2222-222222222222',
    'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11',
    'bb222222-2222-2222-2222-222222222222',
    '11111111-1111-1111-1111-111111111111',
    'eb111111-1111-1111-1111-111111111111',
    CURRENT_DATE,
    'walk_in',
    'elderly',
    'f2222222-2222-2222-2222-222222222222',
    'qr-budi-santoso-002',
    CURRENT_DATE + TIME '07:45'
) ON CONFLICT (id) DO NOTHING;

INSERT INTO pre_anamneses (id, appointment_id, chief_complaint, symptoms, duration_days, allergies, triage_flag)
VALUES (
    'cb222222-2222-2222-2222-222222222222',
    'a2222222-2222-2222-2222-222222222222',
    'Kontrol tekanan darah rutin, leher belakang tegang',
    '["Pusing", "Kaku Leher"]'::jsonb,
    4,
    'Alergi Amoxicillin',
    'attention'
) ON CONFLICT (id) DO NOTHING;

INSERT INTO queues (id, appointment_id, polyclinic_id, queue_date, queue_number, queue_code, status, checked_in_at, called_at, call_count)
VALUES (
    'fb222222-2222-2222-2222-222222222222',
    'a2222222-2222-2222-2222-222222222222',
    '11111111-1111-1111-1111-111111111111',
    CURRENT_DATE,
    2,
    'A-002',
    'called',
    CURRENT_DATE + TIME '08:05',
    CURRENT_DATE + TIME '08:20',
    2
) ON CONFLICT (id) DO NOTHING;

-- Appointment 3: Dewi Lestari (Sudah check in / checked_in)
INSERT INTO appointments (id, facility_id, patient_id, polyclinic_id, schedule_id, visit_date, source, priority, appointment_token, qr_token, created_at)
VALUES (
    'a3333333-3333-3333-3333-333333333333',
    'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11',
    'bb333333-3333-3333-3333-333333333333',
    '11111111-1111-1111-1111-111111111111',
    'eb111111-1111-1111-1111-111111111111',
    CURRENT_DATE,
    'online',
    'pregnant',
    'f3333333-3333-3333-3333-333333333333',
    'qr-dewi-lestari-003',
    CURRENT_DATE + TIME '07:50'
) ON CONFLICT (id) DO NOTHING;

INSERT INTO queues (id, appointment_id, polyclinic_id, queue_date, queue_number, queue_code, status, checked_in_at)
VALUES (
    'fb333333-3333-3333-3333-333333333333',
    'a3333333-3333-3333-3333-333333333333',
    '11111111-1111-1111-1111-111111111111',
    CURRENT_DATE,
    3,
    'A-003',
    'checked_in',
    CURRENT_DATE + TIME '08:15'
) ON CONFLICT (id) DO NOTHING;

-- Appointment 4: Ahmad Fauzi (Menunggu / waiting)
INSERT INTO appointments (id, facility_id, patient_id, polyclinic_id, schedule_id, visit_date, source, priority, appointment_token, qr_token, created_at)
VALUES (
    'a4444444-4444-4444-4444-444444444444',
    'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11',
    'bb444444-4444-4444-4444-444444444444',
    '11111111-1111-1111-1111-111111111111',
    'eb111111-1111-1111-1111-111111111111',
    CURRENT_DATE,
    'online',
    'normal',
    'f4444444-4444-4444-4444-444444444444',
    'qr-ahmad-fauzi-004',
    CURRENT_DATE + TIME '08:00'
) ON CONFLICT (id) DO NOTHING;

INSERT INTO queues (id, appointment_id, polyclinic_id, queue_date, queue_number, queue_code, status)
VALUES (
    'fb444444-4444-4444-4444-444444444444',
    'a4444444-4444-4444-4444-444444444444',
    '11111111-1111-1111-1111-111111111111',
    CURRENT_DATE,
    4,
    'A-004',
    'waiting'
) ON CONFLICT (id) DO NOTHING;

-- 10. HISTORI REKAM MEDIS TERDAHULU (VISITS & DIAGNOSES)
-- Kunjungan 1: Siti Aminah (Minggu lalu)
INSERT INTO visits (
    id, appointment_id, patient_id, doctor_id, vitals,
    subjective, objective, assessment, plan, created_by, created_at, updated_at
) VALUES (
    'ee111111-1111-1111-1111-111111111111',
    'a1111111-1111-1111-1111-111111111111',
    'bb111111-1111-1111-1111-111111111111',
    'd1111111-1111-1111-1111-111111111111',
    '{"bp_sys": 120, "bp_dia": 80, "hr": 82, "temp": 38.2, "spo2": 98, "weight": 55.0, "height": 160.0}'::jsonb,
    'Pasien mengeluh demam, radang tenggorokan saat menelan, dan batuk kering sejak 3 hari.',
    'Keadaan umum: compos mentis, faring tampak hiperemis, pembesaran tonsil T2/T1.',
    'Faringitis Akut (ISPA)',
    'Edukasi minum air hangat, istirahat cukup, terapi simtomatik dan antibiotik oral.',
    'b2c3d4e5-f6a7-4b8c-9d0e-1f2a3b4c5d6e',
    CURRENT_DATE - INTERVAL '7 days',
    CURRENT_DATE - INTERVAL '7 days'
) ON CONFLICT (id) DO NOTHING;

INSERT INTO diagnoses (id, visit_id, icd10_code, is_primary, notes) VALUES
('dd111111-1111-1111-1111-111111111111', 'ee111111-1111-1111-1111-111111111111', 'J02.9', true, 'Diagnosis primer faringitis akut'),
('dd222222-2222-2222-2222-222222222222', 'ee111111-1111-1111-1111-111111111111', 'R50.9', false, 'Febris')
ON CONFLICT (id) DO NOTHING;

-- Kunjungan 2: Budi Santoso (Bulan lalu)
INSERT INTO visits (
    id, appointment_id, patient_id, doctor_id, vitals,
    subjective, objective, assessment, plan, created_by, created_at, updated_at
) VALUES (
    'ee222222-2222-2222-2222-222222222222',
    'a2222222-2222-2222-2222-222222222222',
    'bb222222-2222-2222-2222-222222222222',
    'd1111111-1111-1111-1111-111111111111',
    '{"bp_sys": 145, "bp_dia": 95, "hr": 78, "temp": 36.6, "spo2": 97, "weight": 70.0, "height": 168.0}'::jsonb,
    'Pusing bagian belakang leher, riwayat darah tinggi 2 tahun terakhir tidak rutin minum obat.',
    'KU: Baik, TD 145/95 mmHg, bising jantung (-), ronki (-).',
    'Hipertensi Esensial Stage 1',
    'Edukasi diet rendah garam, olahraga ringan rutin, obat antihipertensi Amlodipine 5mg 1x1.',
    'b2c3d4e5-f6a7-4b8c-9d0e-1f2a3b4c5d6e',
    CURRENT_DATE - INTERVAL '30 days',
    CURRENT_DATE - INTERVAL '30 days'
) ON CONFLICT (id) DO NOTHING;

INSERT INTO diagnoses (id, visit_id, icd10_code, is_primary, notes) VALUES
('dd333333-3333-3333-3333-333333333333', 'ee222222-2222-2222-2222-222222222222', 'I10', true, 'Hipertensi esensial primer')
ON CONFLICT (id) DO NOTHING;
