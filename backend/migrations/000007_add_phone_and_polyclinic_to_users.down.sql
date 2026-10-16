-- 000007_add_phone_and_polyclinic_to_users.down.sql
DROP INDEX IF EXISTS idx_users_polyclinic;
ALTER TABLE users DROP COLUMN IF EXISTS polyclinic_id;
ALTER TABLE users DROP COLUMN IF EXISTS phone;
