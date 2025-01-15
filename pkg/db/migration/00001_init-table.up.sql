-- ==============================================
-- Migration UP: Create Tables and Triggers
-- ==============================================

-- Pastikan menjalankan dalam satu transaction
BEGIN;

-- ==========================
-- 1. Tabel 'users'
-- ==========================
CREATE TABLE IF NOT EXISTS users (
                                     id              SERIAL PRIMARY KEY,
                                     email           VARCHAR(100) NOT NULL UNIQUE,
    password_hash   TEXT         NOT NULL,
    role            VARCHAR(50)  NOT NULL,   -- contoh: 'BORROWER', 'INVESTOR', 'STAFF', 'ADMIN'
    name            VARCHAR(100),
    phone           VARCHAR(20),
    created_at      TIMESTAMP    DEFAULT now(),
    updated_at      TIMESTAMP    DEFAULT now()
    );

-- ==========================
-- 2. Tabel 'loans'
-- ==========================
CREATE TABLE IF NOT EXISTS loans (
                                     id                      SERIAL PRIMARY KEY,
                                     borrower_id             INT NOT NULL,
                                     principal_amount        DECIMAL(20,2) NOT NULL,
    rate                    DECIMAL(5,2) NOT NULL,
    roi                     DECIMAL(5,2) NOT NULL,
    state                   VARCHAR(20)   NOT NULL,  -- 'proposed','approved','invested','disbursed'
    agreement_letter_link   TEXT,
    created_at              TIMESTAMP DEFAULT now(),
    updated_at              TIMESTAMP DEFAULT now(),
    CONSTRAINT fk_borrower
    FOREIGN KEY (borrower_id) REFERENCES users(id)
    );

-- ==========================
-- 3. Tabel 'loan_approval_details'
-- ==========================
CREATE TABLE IF NOT EXISTS loan_approval_details (
                                                     id              SERIAL PRIMARY KEY,
                                                     loan_id         INT NOT NULL,
                                                     staff_id        INT NOT NULL,  -- user dengan role=STAFF/ADMIN
                                                     photo_proof     TEXT,
                                                     approved_date   TIMESTAMP,
                                                     created_at      TIMESTAMP DEFAULT now(),
    updated_at      TIMESTAMP DEFAULT now(),
    CONSTRAINT fk_loan_approval
    FOREIGN KEY (loan_id) REFERENCES loans(id),
    CONSTRAINT fk_staff_approval
    FOREIGN KEY (staff_id) REFERENCES users(id)
    );

-- ==========================
-- 4. Tabel 'loan_investors'
-- ==========================
CREATE TABLE IF NOT EXISTS loan_investors (
                                              id              SERIAL PRIMARY KEY,
                                              loan_id         INT NOT NULL,
                                              investor_id     INT NOT NULL,  -- user dengan role=INVESTOR
                                              amount_invested DECIMAL(20,2) NOT NULL,
    created_at      TIMESTAMP DEFAULT now(),
    updated_at      TIMESTAMP DEFAULT now(),
    CONSTRAINT fk_loan_investors
    FOREIGN KEY (loan_id) REFERENCES loans(id),
    CONSTRAINT fk_investor
    FOREIGN KEY (investor_id) REFERENCES users(id)
    );

-- ==========================
-- 5. Tabel 'loan_disbursement_details'
-- ==========================
CREATE TABLE IF NOT EXISTS loan_disbursement_details (
                                                         id                  SERIAL PRIMARY KEY,
                                                         loan_id             INT NOT NULL,
                                                         staff_id            INT NOT NULL,  -- user dengan role=STAFF
                                                         signed_agreement_doc TEXT,
                                                         disbursed_date      TIMESTAMP,
                                                         created_at          TIMESTAMP DEFAULT now(),
    updated_at          TIMESTAMP DEFAULT now(),
    CONSTRAINT fk_loan_disbursement
    FOREIGN KEY (loan_id) REFERENCES loans(id),
    CONSTRAINT fk_staff_disbursement
    FOREIGN KEY (staff_id) REFERENCES users(id)
    );

-- ==========================
-- 6. Function & Trigger for auto-updating 'updated_at'
--    (Opsional, buat jika ingin updated_at selalu auto-update)
-- ==========================

-- a) Create or Replace Function
CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- b) Triggers (untuk setiap tabel yang ingin auto-update)
CREATE TRIGGER trigger_users_set_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE PROCEDURE set_updated_at();

CREATE TRIGGER trigger_loans_set_updated_at
    BEFORE UPDATE ON loans
    FOR EACH ROW
    EXECUTE PROCEDURE set_updated_at();

CREATE TRIGGER trigger_loan_approval_set_updated_at
    BEFORE UPDATE ON loan_approval_details
    FOR EACH ROW
    EXECUTE PROCEDURE set_updated_at();

CREATE TRIGGER trigger_loan_investors_set_updated_at
    BEFORE UPDATE ON loan_investors
    FOR EACH ROW
    EXECUTE PROCEDURE set_updated_at();

CREATE TRIGGER trigger_loan_disbursement_details_set_updated_at
    BEFORE UPDATE ON loan_disbursement_details
    FOR EACH ROW
    EXECUTE PROCEDURE set_updated_at();

COMMIT;
