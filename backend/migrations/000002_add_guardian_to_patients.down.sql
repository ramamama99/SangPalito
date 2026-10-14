-- Remove guardian columns from patients
ALTER TABLE patients
    DROP COLUMN IF EXISTS guardian_name,
    DROP COLUMN IF EXISTS guardian_nik_hash,
    DROP COLUMN IF EXISTS guardian_nik_encrypted,
    DROP COLUMN IF EXISTS guardian_phone,
    DROP COLUMN IF EXISTS guardian_relation;
