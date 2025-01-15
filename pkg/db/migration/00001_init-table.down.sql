-- ==============================================
-- Migration DOWN: Drop Tables and Triggers
-- ==============================================

BEGIN;

-- 1. Hapus trigger setiap tabel (jika Anda membuat trigger).
DROP TRIGGER IF EXISTS trigger_users_set_updated_at ON users;
DROP TRIGGER IF EXISTS trigger_loans_set_updated_at ON loans;
DROP TRIGGER IF EXISTS trigger_loan_approval_set_updated_at ON loan_approval_details;
DROP TRIGGER IF EXISTS trigger_loan_investors_set_updated_at ON loan_investors;
DROP TRIGGER IF EXISTS trigger_loan_disbursement_details_set_updated_at ON loan_disbursement_details;

-- 2. Hapus function set_updated_at
DROP FUNCTION IF EXISTS set_updated_at() CASCADE;

-- 3. Drop tabel dengan urutan child -> parent
--    (Karena foreign key constraints)
DROP TABLE IF EXISTS loan_disbursement_details;
DROP TABLE IF EXISTS loan_investors;
DROP TABLE IF EXISTS loan_approval_details;
DROP TABLE IF EXISTS loans;
DROP TABLE IF EXISTS users;

COMMIT;
