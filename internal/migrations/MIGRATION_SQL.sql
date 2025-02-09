CREATE TABLE IF NOT EXISTS borrowers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    agreement_letter TEXT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp NULL
    );


INSERT INTO borrowers (name, email) VALUES
('Alice Johnson', 'test@gmail.com'),
('Bob Smith', 'test01@gmail.com');

CREATE TABLE IF NOT EXISTS loans (
    id SERIAL PRIMARY KEY,
    borrower_id INT REFERENCES borrowers(id) ON DELETE CASCADE,
    principal_amount DECIMAL(15,2) NOT NULL,
    rate DECIMAL(5,2) NOT NULL,
    total_interest DECIMAL(15,2) NULL,
    loan_term INT not null,
    agreement_letter TEXT NULL,
    field_validator_picture_proof TEXT NULL,
    field_validator_id INT NULL,
    field_officer_id INT null,
    approval_date TIMESTAMP null,
    disbursment_date TIMESTAMP null,
    status VARCHAR(20) NOT NULL CHECK (status IN ('proposed', 'approved', 'invested', 'disbursed')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );


CREATE TABLE IF NOT EXISTS investments (
    id SERIAL PRIMARY KEY,
    investor_id INT REFERENCES investors(id) ON DELETE CASCADE,
    loan_id INT,
    amount DECIMAL(15,2) NOT NULL,
    roi DECIMAL(5,2) NOT NULL,
    total_gain DECIMAL(15,2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

CREATE TABLE IF NOT EXISTS investors (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    agreement_letter TEXT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

INSERT INTO investors (name, email) VALUES
('invest01', 'invest01@gmail.com'),
('invest02', 'invest02@gmail.com');