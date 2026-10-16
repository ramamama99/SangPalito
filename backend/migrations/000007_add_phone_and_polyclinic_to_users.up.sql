-- 000007_add_phone_and_polyclinic_to_users.up.sql
ALTER TABLE users ADD COLUMN IF NOT EXISTS phone VARCHAR(30);
ALTER TABLE users ADD COLUMN IF NOT EXISTS polyclinic_id UUID REFERENCES polyclinics(id) ON DELETE SET NULL;
CREATE INDEX IF NOT EXISTS idx_users_polyclinic ON users (polyclinic_id);
