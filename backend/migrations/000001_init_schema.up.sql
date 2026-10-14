-- Extension wajib
CREATE EXTENSION IF NOT EXISTS "pgcrypto";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";

-- ENUM Types
CREATE TYPE gender_t AS ENUM ('male', 'female');
CREATE TYPE guarantor_t AS ENUM ('bpjs', 'umum', 'asuransi', 'lainnya');
CREATE TYPE queue_status_t AS ENUM ('waiting', 'checked_in', 'called', 'in_service', 'done', 'cancelled', 'no_show');
CREATE TYPE appt_source_t AS ENUM ('online', 'walk_in');
CREATE TYPE priority_t AS ENUM ('normal', 'elderly', 'pregnant', 'disabled', 'emergency');
CREATE TYPE notif_channel_t AS ENUM ('whatsapp', 'sms', 'email');
CREATE TYPE notif_status_t AS ENUM ('pending', 'sent', 'failed');

-- 1. FASILITAS (PUSKESMAS / KLINIK)
CREATE TABLE facilities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(150) NOT NULL,
    type VARCHAR(30) NOT NULL, -- puskesmas | klinik
    code VARCHAR(30) UNIQUE,
    address TEXT,
    phone VARCHAR(20),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- 2. ROLE & PENGGUNA (STAF)
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
    refresh_token_hash VARCHAR(64),
    is_active BOOLEAN NOT NULL DEFAULT true,
    last_login_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE user_roles (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role_id SMALLINT NOT NULL REFERENCES roles (id),
    PRIMARY KEY (user_id, role_id)
);

-- 3. DATA PASIEN
CREATE TABLE patients (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    facility_id UUID NOT NULL REFERENCES facilities (id),
    medical_record_no VARCHAR(20) NOT NULL,
    nik_encrypted BYTEA,
    nik_hash VARCHAR(64),
    full_name VARCHAR(150) NOT NULL,
    birth_date DATE NOT NULL,
    gender gender_t NOT NULL,
    phone VARCHAR(20),
    address TEXT,
    guarantor guarantor_t NOT NULL DEFAULT 'umum',
    bpjs_no VARCHAR(20),
    consent_given_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    UNIQUE (facility_id, medical_record_no),
    UNIQUE (facility_id, nik_hash)
);

CREATE INDEX idx_patients_name_trgm ON patients USING gin (full_name gin_trgm_ops);
CREATE INDEX idx_patients_phone ON patients (phone);

-- 4. POLIKLINIK & DOKTER
CREATE TABLE polyclinics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    facility_id UUID NOT NULL REFERENCES facilities (id),
    name VARCHAR(80) NOT NULL,
    queue_prefix CHAR(2) NOT NULL,
    avg_service_min SMALLINT NOT NULL DEFAULT 10,
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

-- 5. PENDAFTARAN KUNJUNGAN & ANTREAN
CREATE TABLE appointments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    facility_id UUID NOT NULL REFERENCES facilities(id),
    patient_id UUID NOT NULL REFERENCES patients (id),
    polyclinic_id UUID NOT NULL REFERENCES polyclinics(id),
    schedule_id UUID REFERENCES schedules(id),
    visit_date DATE NOT NULL,
    source appt_source_t NOT NULL DEFAULT 'online',
    priority priority_t NOT NULL DEFAULT 'normal',
    appointment_token UUID NOT NULL DEFAULT gen_random_uuid(),
    qr_token VARCHAR(64) UNIQUE NOT NULL,
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_appt_date_poli ON appointments (visit_date, polyclinic_id);

CREATE TABLE queues (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    appointment_id UUID UNIQUE NOT NULL REFERENCES appointments (id) ON DELETE CASCADE,
    polyclinic_id UUID NOT NULL REFERENCES polyclinics(id),
    queue_date DATE NOT NULL,
    queue_number INT NOT NULL,
    queue_code VARCHAR(10) NOT NULL,
    status queue_status_t NOT NULL DEFAULT 'waiting',
    checked_in_at TIMESTAMPTZ,
    called_at TIMESTAMPTZ,
    call_count SMALLINT NOT NULL DEFAULT 0,
    service_start TIMESTAMPTZ,
    service_end TIMESTAMPTZ,
    UNIQUE (polyclinic_id, queue_date, queue_number)
);

CREATE INDEX idx_queue_live ON queues (polyclinic_id, queue_date, status);

-- 6. PRA-ANAMNESIS & PEMERIKSAAN MEDIS
CREATE TABLE pre_anamneses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    appointment_id UUID UNIQUE NOT NULL REFERENCES appointments (id) ON DELETE CASCADE,
    chief_complaint TEXT NOT NULL,
    symptoms JSONB,
    duration_days SMALLINT,
    allergies TEXT,
    triage_flag VARCHAR(20),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE visits (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    appointment_id UUID UNIQUE NOT NULL REFERENCES appointments (id),
    patient_id UUID NOT NULL REFERENCES patients (id),
    doctor_id UUID REFERENCES doctors(id),
    vitals JSONB,
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

CREATE INDEX idx_icd10_desc ON icd10_codes USING gin(to_tsvector('indonesian', description));

CREATE TABLE diagnoses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    visit_id UUID NOT NULL REFERENCES visits (id) ON DELETE CASCADE,
    icd10_code VARCHAR(10) NOT NULL REFERENCES icd10_codes (code),
    is_primary BOOLEAN NOT NULL DEFAULT false,
    notes TEXT
);

-- 7. NOTIFIKASI & AUDIT LOG
CREATE TABLE notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    appointment_id UUID REFERENCES appointments (id),
    channel notif_channel_t NOT NULL DEFAULT 'whatsapp',
    recipient VARCHAR(50) NOT NULL,
    message TEXT NOT NULL,
    status notif_status_t NOT NULL DEFAULT 'pending',
    sent_at TIMESTAMPTZ,
    error_msg TEXT,
    retry_count SMALLINT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE audit_logs (
    id BIGSERIAL PRIMARY KEY,
    user_id UUID REFERENCES users (id),
    action VARCHAR(50) NOT NULL,
    entity VARCHAR(50) NOT NULL,
    entity_id UUID,
    ip_address INET,
    metadata JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_audit_entity ON audit_logs (entity, entity_id, created_at DESC);

-- Default Facility Seed
INSERT INTO facilities (id, name, type, code, address, phone)
VALUES ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Puskesmas SangPalito', 'puskesmas', 'PSK-SANGPALITO', 'Jl. Kesehatan No. 1', '08123456789')
ON CONFLICT (id) DO NOTHING;
