ALTER TABLE patients
ADD COLUMN guardian_name VARCHAR(150),
ADD COLUMN guardian_nik_encrypted BYTEA,
ADD COLUMN guardian_nik_hash VARCHAR(64),
ADD COLUMN guardian_relation VARCHAR(50),
ADD COLUMN guardian_phone VARCHAR(20);
