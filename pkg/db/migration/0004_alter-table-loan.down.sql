-- down.sql
ALTER TABLE loans ALTER COLUMN rate TYPE DECIMAL(5,2);
ALTER TABLE loans ALTER COLUMN roi TYPE DECIMAL(5,2);
