-- LinkWorld Database Schema
-- PostgreSQL

CREATE DATABASE IF NOT EXISTS linkworld;

\c linkworld;

-- 用户表
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    wallet_addr VARCHAR(42) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL,
    token_id BIGINT DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    registered_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_wallet ON users(wallet_addr);
CREATE INDEX idx_users_email ON users(email);

-- 运营商表
CREATE TABLE IF NOT EXISTS operators (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    region VARCHAR(100) NOT NULL,
    country_code VARCHAR(10) NOT NULL,
    required_deposit VARCHAR(50) NOT NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_operators_country ON operators(country_code);

-- 保证金表
CREATE TABLE IF NOT EXISTS deposits (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    amount VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_deposits_user ON deposits(user_id);

-- 账单表
CREATE TABLE IF NOT EXISTS bills (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    operator_id BIGINT NOT NULL REFERENCES operators(id) ON DELETE CASCADE,
    amount VARCHAR(50) NOT NULL,
    platform_fee VARCHAR(50) NOT NULL,
    is_paid BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    paid_at TIMESTAMP
);

CREATE INDEX idx_bills_user ON bills(user_id);
CREATE INDEX idx_bills_operator ON bills(operator_id);
CREATE INDEX idx_bills_is_paid ON bills(is_paid);

-- 用户服务表（虚拟号）
CREATE TABLE IF NOT EXISTS user_services (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    operator_id BIGINT NOT NULL REFERENCES operators(id) ON DELETE CASCADE,
    virtual_number VARCHAR(50) NOT NULL,
    password VARCHAR(255) NOT NULL,
    is_active BOOLEAN DEFAULT true,
    activated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_user_services_operator ON user_services(operator_id);
CREATE INDEX idx_user_services_is_active ON user_services(is_active);

-- 用量数据表（预言机数据）
CREATE TABLE IF NOT EXISTS usage_data (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    operator_id BIGINT NOT NULL REFERENCES operators(id) ON DELETE CASCADE,
    data_usage BIGINT DEFAULT 0,
    call_usage BIGINT DEFAULT 0,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    signature TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_usage_data_user ON usage_data(user_id);
CREATE INDEX idx_usage_data_operator ON usage_data(operator_id);
CREATE INDEX idx_usage_data_timestamp ON usage_data(timestamp);

-- 初始化运营商数据
INSERT INTO operators (name, region, country_code, required_deposit) VALUES 
('T-Mobile US', 'United States', 'US', '0.01'),
('Vodafone UK', 'United Kingdom', 'GB', '0.008'),
('Orange France', 'France', 'FR', '0.008'),
('MTS Russia', 'Russia', 'RU', '0.005'),
('SoftBank Japan', 'Japan', 'JP', '0.012'),
('Viettel Vietnam', 'Vietnam', 'VN', '0.003'),
('Unitel Laos', 'Laos', 'LA', '0.003'),
('Smart Cambodia', 'Cambodia', 'KH', '0.003'),
('AIS Thailand', 'Thailand', 'TH', '0.004'),
('Maxis Malaysia', 'Malaysia', 'MY', '0.004'),
('Globe Philippines', 'Philippines', 'PH', '0.003');