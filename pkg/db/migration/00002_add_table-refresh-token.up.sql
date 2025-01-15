-- Create the refresh_tokens table
CREATE TABLE IF NOT EXISTS refresh_tokens (
                                              id SERIAL PRIMARY KEY,                           -- Primary key
                                              user_id INT NOT NULL,                            -- FK ke users.id
                                              refresh_token TEXT NOT NULL UNIQUE,              -- Token unik untuk refresh
                                              expires_at TIMESTAMP NOT NULL,                   -- Waktu kedaluwarsa token
                                              created_at TIMESTAMP DEFAULT now(),              -- Waktu pembuatan token
    updated_at TIMESTAMP DEFAULT now(),              -- Waktu terakhir pembaruan
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
    );
