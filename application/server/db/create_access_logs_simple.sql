USE fabric_mims;

CREATE TABLE IF NOT EXISTS prescription_access_logs (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    log_id VARCHAR(100) UNIQUE NOT NULL,
    prescription_id VARCHAR(100) NOT NULL,
    prescription_no VARCHAR(100) NOT NULL,
    patient_id VARCHAR(100) NOT NULL,
    patient_name VARCHAR(100) NOT NULL,
    accessor_id VARCHAR(100) NOT NULL,
    accessor_name VARCHAR(100) NOT NULL,
    accessor_role VARCHAR(50) NOT NULL,
    accessor_organization VARCHAR(100) DEFAULT NULL,
    accessor_organization_name VARCHAR(200) DEFAULT NULL,
    access_type VARCHAR(50) NOT NULL,
    access_reason VARCHAR(500) DEFAULT NULL,
    ip_address VARCHAR(50) DEFAULT NULL,
    user_agent TEXT,
    tx_id VARCHAR(100) DEFAULT NULL,
    accessed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_access_prescription (prescription_id),
    INDEX idx_access_patient (patient_id),
    INDEX idx_access_accessor (accessor_id),
    INDEX idx_access_time (accessed_at),
    INDEX idx_access_type (access_type)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
