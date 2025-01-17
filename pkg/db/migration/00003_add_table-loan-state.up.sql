-- Create the loan_state_histories table
CREATE TABLE loan_state_histories (
                                      id SERIAL PRIMARY KEY,                          -- Primary key
                                      loan_id INT NOT NULL,                           -- Foreign key to loans.id
                                      previous_state VARCHAR(20) NOT NULL,           -- State before the change
                                      new_state VARCHAR(20) NOT NULL,                -- State after the change
                                      action_by INT NOT NULL,                        -- Foreign key to users.id
                                      action_at TIMESTAMP DEFAULT now(),             -- Timestamp of the action
                                      remarks TEXT,                                  -- Optional remarks for the action
                                      CONSTRAINT fk_loan FOREIGN KEY (loan_id) REFERENCES loans (id) ON DELETE CASCADE,
                                      CONSTRAINT fk_user FOREIGN KEY (action_by) REFERENCES users (id) ON DELETE SET NULL
);
