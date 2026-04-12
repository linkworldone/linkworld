-- LinkWorld Database Schema
-- MySQL 8.0+

CREATE DATABASE IF NOT EXISTS linkworld CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE linkworld;

-- 用户表
CREATE TABLE IF NOT EXISTS users (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    wallet_addr VARCHAR(42) NOT NULL COMMENT '钱包地址',
    email VARCHAR(255) NOT NULL COMMENT '邮箱',
    token_id BIGINT UNSIGNED DEFAULT 0 COMMENT 'NFT Token ID',
    is_active TINYINT(1) DEFAULT 1 COMMENT '是否活跃',
    registered_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '注册时间',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY idx_wallet (wallet_addr),
    KEY idx_email (email)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';

-- 运营商表
CREATE TABLE IF NOT EXISTS operators (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL COMMENT '运营商名称',
    region VARCHAR(100) NOT NULL COMMENT '地区',
    required_deposit VARCHAR(50) NOT NULL COMMENT '所需保证金',
    is_active TINYINT(1) DEFAULT 1 COMMENT '是否激活',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='运营商表';

-- 保证金表
CREATE TABLE IF NOT EXISTS deposits (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    amount VARCHAR(50) NOT NULL COMMENT '保证金金额',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    KEY idx_user (user_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='保证金表';

-- 账单表
CREATE TABLE IF NOT EXISTS bills (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    operator_id BIGINT UNSIGNED NOT NULL COMMENT '运营商ID',
    amount VARCHAR(50) NOT NULL COMMENT '金额',
    platform_fee VARCHAR(50) NOT NULL COMMENT '平台手续费',
    is_paid TINYINT(1) DEFAULT 0 COMMENT '是否已支付',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    paid_at TIMESTAMP NULL COMMENT '支付时间',
    KEY idx_user (user_id),
    KEY idx_operator (operator_id),
    KEY idx_is_paid (is_paid),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (operator_id) REFERENCES operators(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='账单表';

-- 用户服务表（虚拟号）
CREATE TABLE IF NOT EXISTS user_services (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    operator_id BIGINT UNSIGNED NOT NULL COMMENT '运营商ID',
    virtual_number VARCHAR(50) NOT NULL COMMENT '虚拟号码',
    password VARCHAR(255) NOT NULL COMMENT '密码',
    is_active TINYINT(1) DEFAULT 1 COMMENT '是否激活',
    activated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '激活时间',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY idx_user (user_id),
    KEY idx_operator (operator_id),
    KEY idx_is_active (is_active),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (operator_id) REFERENCES operators(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户服务表';

-- 用量数据表（预言机数据）
CREATE TABLE IF NOT EXISTS usage_data (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    operator_id BIGINT UNSIGNED NOT NULL COMMENT '运营商ID',
    data_usage BIGINT UNSIGNED DEFAULT 0 COMMENT '流量使用量',
    call_usage BIGINT UNSIGNED DEFAULT 0 COMMENT '通话使用量',
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '记录时间',
    signature TEXT COMMENT '数据签名',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    KEY idx_user (user_id),
    KEY idx_operator (operator_id),
    KEY idx_timestamp (timestamp),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (operator_id) REFERENCES operators(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用量数据表';

-- 初始化默认运营商数据
INSERT INTO operators (name, region, required_deposit) VALUES 
('T-Mobile US', 'United States', '0.01'),
('Vodafone UK', 'United Kingdom', '0.008');
