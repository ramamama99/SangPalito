-- Migration 000008: Tambah fitur dokter terhubung akun user, first-time login, dan reset password

-- 1. Tambah kolom reset password dan first-time login ke tabel users
ALTER TABLE users ADD COLUMN IF NOT EXISTS is_first_login BOOLEAN NOT NULL DEFAULT false;
ALTER TABLE users ADD COLUMN IF NOT EXISTS reset_token VARCHAR(120);
ALTER TABLE users ADD COLUMN IF NOT EXISTS reset_expires_at TIMESTAMPTZ;
CREATE INDEX IF NOT EXISTS idx_users_reset_token ON users (reset_token);

-- 2. Tambah kolom relasi dan profil ke tabel doctors
ALTER TABLE doctors ADD COLUMN IF NOT EXISTS polyclinic_id UUID REFERENCES polyclinics(id) ON DELETE SET NULL;
ALTER TABLE doctors ADD COLUMN IF NOT EXISTS phone VARCHAR(30);
ALTER TABLE doctors ADD COLUMN IF NOT EXISTS email VARCHAR(150);
ALTER TABLE doctors ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ DEFAULT now();
ALTER TABLE doctors ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ DEFAULT now();
CREATE INDEX IF NOT EXISTS idx_doctors_user_id ON doctors (user_id);
CREATE INDEX IF NOT EXISTS idx_doctors_polyclinic ON doctors (polyclinic_id);

-- 3. Mapping data dokter ke poliklinik masing-masing
UPDATE doctors SET polyclinic_id = '11111111-1111-1111-1111-111111111111', email = 'rama@antrisehat.id', phone = '081234567890' WHERE id = 'd1111111-1111-1111-1111-111111111111';
UPDATE doctors SET polyclinic_id = '11111111-1111-1111-1111-111111111111', email = 'maya@antrisehat.id', phone = '081287654321' WHERE id = 'd2222222-2222-2222-2222-222222222222';
UPDATE doctors SET polyclinic_id = '22222222-2222-2222-2222-222222222222', email = 'rani@antrisehat.id', phone = '081901234567' WHERE id = 'd3333333-3333-3333-3333-333333333333';
UPDATE doctors SET polyclinic_id = '33333333-3333-3333-3333-333333333333', email = 'siti@antrisehat.id', phone = '081223344556' WHERE id = 'd4444444-4444-4444-4444-444444444444';

-- Dokter tambahan agar setiap poli memiliki beberapa pilihan dokter
INSERT INTO doctors (id, facility_id, full_name, specialty, str_no, is_active, polyclinic_id, email, phone)
VALUES
('d5555555-5555-5555-5555-555555555555', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'dr. Aditya Rahman', 'Dokter Umum', 'STR-3201-9944', true, '11111111-1111-1111-1111-111111111111', 'aditya@antrisehat.id', '081398765432'),
('d6666666-6666-6666-6666-666666666666', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'drg. Fajar Hidayat', 'Dokter Gigi', 'STR-3201-7788', true, '22222222-2222-2222-2222-222222222222', 'fajar@antrisehat.id', '081812345678'),
('d7777777-7777-7777-7777-777777777777', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'dr. Hendra Prasetyo', 'Spesialis Penyakit Dalam', 'STR-3201-5522', true, '44444444-4444-4444-4444-444444444444', 'hendra@antrisehat.id', '081334455667')
ON CONFLICT (id) DO UPDATE SET polyclinic_id = EXCLUDED.polyclinic_id, email = EXCLUDED.email, phone = EXCLUDED.phone;

-- Fallback umum
UPDATE doctors SET polyclinic_id = '11111111-1111-1111-1111-111111111111' WHERE specialty ILIKE '%umum%' AND polyclinic_id IS NULL;
UPDATE doctors SET polyclinic_id = '22222222-2222-2222-2222-222222222222' WHERE (specialty ILIKE '%gigi%' OR full_name ILIKE '%drg%') AND polyclinic_id IS NULL;
UPDATE doctors SET polyclinic_id = '33333333-3333-3333-3333-333333333333' WHERE (specialty ILIKE '%anak%' OR specialty ILIKE '%kia%' OR specialty ILIKE '%sp.a%') AND polyclinic_id IS NULL;
UPDATE doctors SET polyclinic_id = '44444444-4444-4444-4444-444444444444' WHERE (specialty ILIKE '%lansia%' OR specialty ILIKE '%geriatri%' OR specialty ILIKE '%dalam%') AND polyclinic_id IS NULL;
