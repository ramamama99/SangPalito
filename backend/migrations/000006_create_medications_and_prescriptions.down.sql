-- 000006_create_medications_and_prescriptions.down.sql

ALTER TABLE visits DROP COLUMN IF EXISTS prescription_no;
ALTER TABLE visits DROP COLUMN IF EXISTS medication_plan;

DROP TABLE IF EXISTS prescription_items;
DROP TABLE IF EXISTS prescriptions;
DROP TABLE IF EXISTS medications;
