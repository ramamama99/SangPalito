-- ============================================
-- AntriSehat Initial Schema
-- Migration: 000001_init_schema.up.sql
-- ============================================

-- Enable required extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- ============================================
-- 1. Facilities (Fasilitas Kesehatan)
-- ============================================
CREATE TABLE facilities (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name        VARCHAR(200) NOT NULL,
    code        VARCHAR(20) NOT NULL UNIQUE,
    address     TEXT,
    phone       VARCHAR(20),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================
-- 2. Users (Staf / Pengguna Sistem)
-- ============================================
CREATE TABLE users (
    id                   UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    facility_id          UUID NOT NULL REFERENCES facilities(id),
    email                VARCHAR(150) NOT NULL,
    password_hash        VARCHAR(255) NOT NULL,
    full_name            VARCHAR(150) NOT NULL,
    role                 VARCHAR(20) NOT NULL CHECK (role IN ('admin', 'registrar', 'medic', 'head')),
    is_active            BOOLEAN NOT NULL DEFAULT true,
    last_login_at        TIMESTAMPTZ,
    refresh_token_hash   VARCHAR(64),
    created_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(email, facility_id)
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_facility_id ON users(facility_id);

-- ============================================
-- 3. Patients (Pasien)
-- ============================================
CREATE TABLE patients (
    id                   UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    facility_id          UUID NOT NULL REFERENCES facilities(id),
    medical_record_no    VARCHAR(20) NOT NULL,
    nik_hash             VARCHAR(64) NOT NULL,
    nik_encrypted        BYTEA NOT NULL,
    full_name            VARCHAR(150) NOT NULL,
    birth_date           DATE NOT NULL,
    gender               VARCHAR(10) NOT NULL CHECK (gender IN ('male', 'female')),
    phone                VARCHAR(20),
    address              TEXT,
    blood_type           VARCHAR(5),
    guarantor            VARCHAR(20) NOT NULL DEFAULT 'umum' CHECK (guarantor IN ('bpjs', 'umum', 'asuransi', 'lainnya')),
    bpjs_no              VARCHAR(20),
    created_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at           TIMESTAMPTZ,
    UNIQUE(nik_hash, facility_id)
);

CREATE INDEX idx_patients_nik_hash ON patients(nik_hash);
CREATE INDEX idx_patients_facility ON patients(facility_id);
CREATE INDEX idx_patients_medical_record ON patients(medical_record_no);
CREATE INDEX idx_patients_fullname_trgm ON patients USING gin (full_name gin_trgm_ops);

-- ============================================
-- 4. Polyclinics (Poliklinik)
-- ============================================
CREATE TABLE polyclinics (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    facility_id UUID NOT NULL REFERENCES facilities(id),
    name        VARCHAR(100) NOT NULL,
    code        VARCHAR(20) NOT NULL,
    is_active   BOOLEAN NOT NULL DEFAULT true,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(code, facility_id)
);

-- ============================================
-- 5. Doctors (Dokter)
-- ============================================
CREATE TABLE doctors (
    id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    facility_id   UUID NOT NULL REFERENCES facilities(id),
    user_id       UUID REFERENCES users(id),
    full_name     VARCHAR(150) NOT NULL,
    specialization VARCHAR(100),
    sip_no        VARCHAR(50),
    is_active     BOOLEAN NOT NULL DEFAULT true,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================
-- 6. Schedules (Jadwal Praktik)
-- ============================================
CREATE TABLE schedules (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    facility_id     UUID NOT NULL REFERENCES facilities(id),
    polyclinic_id   UUID NOT NULL REFERENCES polyclinics(id),
    doctor_id       UUID NOT NULL REFERENCES doctors(id),
    day_of_week     SMALLINT NOT NULL CHECK (day_of_week BETWEEN 0 AND 6),
    start_time      TIME NOT NULL,
    end_time        TIME NOT NULL,
    quota           INT NOT NULL DEFAULT 30,
    is_active       BOOLEAN NOT NULL DEFAULT true,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_schedules_poly ON schedules(polyclinic_id);
CREATE INDEX idx_schedules_doctor ON schedules(doctor_id);

-- ============================================
-- 7. Appointments (Kunjungan / Pendaftaran)
-- ============================================
CREATE TABLE appointments (
    id                  UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    facility_id         UUID NOT NULL REFERENCES facilities(id),
    patient_id          UUID NOT NULL REFERENCES patients(id),
    polyclinic_id       UUID NOT NULL REFERENCES polyclinics(id),
    schedule_id         UUID NOT NULL REFERENCES schedules(id),
    doctor_id           UUID NOT NULL REFERENCES doctors(id),
    visit_date          DATE NOT NULL,
    appointment_token   UUID NOT NULL DEFAULT uuid_generate_v4() UNIQUE,
    qr_token            UUID NOT NULL DEFAULT uuid_generate_v4() UNIQUE,
    chief_complaint     TEXT,
    status              VARCHAR(20) NOT NULL DEFAULT 'registered' CHECK (status IN ('registered', 'checked_in', 'completed', 'cancelled')),
    registered_by       UUID REFERENCES users(id),
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(patient_id, polyclinic_id, visit_date)
);

CREATE INDEX idx_appointments_patient ON appointments(patient_id);
CREATE INDEX idx_appointments_date ON appointments(visit_date);
CREATE INDEX idx_appointments_token ON appointments(appointment_token);
CREATE INDEX idx_appointments_qr ON appointments(qr_token);

-- ============================================
-- 8. Queues (Antrean)
-- ============================================
CREATE TABLE queues (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    facility_id     UUID NOT NULL REFERENCES facilities(id),
    appointment_id  UUID NOT NULL REFERENCES appointments(id) UNIQUE,
    polyclinic_id   UUID NOT NULL REFERENCES polyclinics(id),
    queue_number    INT NOT NULL,
    status          VARCHAR(20) NOT NULL DEFAULT 'waiting' CHECK (status IN ('waiting', 'checked_in', 'called', 'in_service', 'done', 'cancelled', 'no_show')),
    call_count      INT NOT NULL DEFAULT 0,
    called_at       TIMESTAMPTZ,
    service_start   TIMESTAMPTZ,
    service_end     TIMESTAMPTZ,
    queue_date      DATE NOT NULL DEFAULT CURRENT_DATE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_queues_poly_date ON queues(polyclinic_id, queue_date);
CREATE INDEX idx_queues_status ON queues(status);

-- ============================================
-- 9. Visits (Kunjungan Klinis)
-- ============================================
CREATE TABLE visits (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    facility_id     UUID NOT NULL REFERENCES facilities(id),
    appointment_id  UUID NOT NULL REFERENCES appointments(id),
    patient_id      UUID NOT NULL REFERENCES patients(id),
    doctor_id       UUID NOT NULL REFERENCES doctors(id),
    polyclinic_id   UUID NOT NULL REFERENCES polyclinics(id),
    visit_date      DATE NOT NULL,
    diagnosis       TEXT,
    icd_code        VARCHAR(10),
    treatment_plan  TEXT,
    prescription    TEXT,
    notes           TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_visits_patient ON visits(patient_id);
CREATE INDEX idx_visits_date ON visits(visit_date);

-- ============================================
-- 10. Vital Signs (Tanda-tanda Vital)
-- ============================================
CREATE TABLE vital_signs (
    id                  UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    visit_id            UUID NOT NULL REFERENCES visits(id),
    systolic_bp         INT,
    diastolic_bp        INT,
    heart_rate          INT,
    temperature         DECIMAL(4,1),
    respiratory_rate    INT,
    oxygen_saturation   DECIMAL(4,1),
    weight              DECIMAL(5,1),
    height              DECIMAL(5,1),
    recorded_by         UUID REFERENCES users(id),
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================
-- 11. Pre-Anamnesis
-- ============================================
CREATE TABLE pre_anamnesis (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    visit_id        UUID NOT NULL REFERENCES visits(id),
    allergies       TEXT,
    current_meds    TEXT,
    medical_history TEXT,
    chief_complaint TEXT,
    recorded_by     UUID REFERENCES users(id),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================
-- 12. Audit Logs
-- ============================================
CREATE TABLE audit_logs (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id     UUID,
    action      VARCHAR(50) NOT NULL,
    entity      VARCHAR(50) NOT NULL,
    entity_id   UUID,
    ip_address  VARCHAR(45),
    meta        JSONB,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_audit_logs_user ON audit_logs(user_id);
CREATE INDEX idx_audit_logs_entity ON audit_logs(entity, entity_id);
CREATE INDEX idx_audit_logs_created ON audit_logs(created_at);

-- ============================================
-- Seed default facility
-- ============================================
INSERT INTO facilities (id, name, code, address) VALUES
    ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Puskesmas SangPalito', 'PSK-SANGPALITO', 'Jl. Kesehatan No. 1');
